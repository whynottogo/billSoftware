package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
)

const familyInviteFallbackBaseURL = "http://localhost:9000"

var familyInviteCodePattern = regexp.MustCompile(`(?i)FAM-[A-Z0-9-]+`)

type UserFamilyHandler struct {
	engine *xorm.Engine
}

type familyBaseRow struct {
	ID              uint64    `xorm:"id"`
	FamilyUID       string    `xorm:"family_uid"`
	FamilyName      string    `xorm:"family_name"`
	CreatorUserID   uint64    `xorm:"creator_user_id"`
	CreatedAt       time.Time `xorm:"created_at"`
	CreatorNickname string    `xorm:"creator_nickname"`
	CreatorUsername string    `xorm:"creator_username"`
}

type familyMemberRow struct {
	UserID   uint64    `xorm:"user_id"`
	Nickname string    `xorm:"nickname"`
	Username string    `xorm:"username"`
	JoinedAt time.Time `xorm:"joined_at"`
}

type familyAggregateRow struct {
	PeriodKey string  `xorm:"period_key"`
	UserID    uint64  `xorm:"user_id"`
	Income    float64 `xorm:"income"`
	Expense   float64 `xorm:"expense"`
}

type familyScope struct {
	ID            uint64
	UID           string
	Name          string
	CreatorUserID uint64
	CreatorName   string
	CreatedAt     time.Time
	Members       []familyMemberInfo
}

type familyMemberInfo struct {
	UserID   uint64 `json:"userId"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Color    string `json:"color"`
	JoinedAt string `json:"joinedAt"`
}

type familyPeriodOption struct {
	Key     string  `json:"key"`
	Label   string  `json:"label"`
	Note    string  `json:"note"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type familyShareRow struct {
	Name    string  `json:"name"`
	Role    string  `json:"role"`
	Color   string  `json:"color"`
	Value   float64 `json:"value"`
	Percent float64 `json:"percent"`
}

func NewUserFamilyHandler(engine *xorm.Engine) *UserFamilyHandler {
	return &UserFamilyHandler{engine: engine}
}

func (h *UserFamilyHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	baseRows := make([]familyBaseRow, 0)
	if err := h.engine.SQL(`
		SELECT
			f.id,
			f.family_uid,
			f.family_name,
			f.creator_user_id,
			f.created_at,
			cu.nickname AS creator_nickname,
			cu.username AS creator_username
		FROM families f
		INNER JOIN family_members self_fm ON self_fm.family_id = f.id AND self_fm.user_id = ?
		INNER JOIN users cu ON cu.id = f.creator_user_id
		ORDER BY f.created_at DESC, f.id DESC
	`, userID).Find(&baseRows); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query families failed")
		return
	}

	list := make([]gin.H, 0, len(baseRows))
	totalMembers := 0
	for _, base := range baseRows {
		scope, err := h.loadFamilyScopeByID(base.ID)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "query family members failed")
			return
		}

		monthOptions, yearOptions, err := h.loadFamilyPeriodOptions(scope.ID)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "query family summary failed")
			return
		}

		currentMonth := monthOptions[0]
		currentYear := yearOptions[0]
		totalMembers += len(scope.Members)

		list = append(list, gin.H{
			"id":           scope.UID,
			"name":         scope.Name,
			"slogan":       buildFamilySlogan(scope.Name, scope.CreatedAt),
			"creator":      scope.CreatorName,
			"createdAt":    scope.CreatedAt.Format("2006-01-02"),
			"inviteCode":   scope.UID,
			"inviteLink":   buildFamilyInviteLink(c, scope.UID),
			"memberCount":  len(scope.Members),
			"monthIncome":  currentMonth.Income,
			"monthExpense": currentMonth.Expense,
			"monthBalance": currentMonth.Balance,
			"yearIncome":   currentYear.Income,
			"yearExpense":  currentYear.Expense,
			"yearBalance":  currentYear.Balance,
			"members":      scope.Members,
		})
	}

	response.Success(c, gin.H{
		"list": list,
		"overview": gin.H{
			"familyCount":  len(list),
			"totalMembers": totalMembers,
			"joinedCount":  len(list),
		},
	})
}

func (h *UserFamilyHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	payload, err := parseJSONMap(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid family payload")
		return
	}

	familyName := strings.TrimSpace(readString(payload, "name", "family_name", "familyName"))
	if familyName == "" {
		response.Fail(c, http.StatusBadRequest, "family name is required")
		return
	}

	session := h.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		response.Fail(c, http.StatusInternalServerError, "begin transaction failed")
		return
	}

	family := &model.Family{
		FamilyName:    familyName,
		CreatorUserID: userID,
	}

	inserted := false
	for attempt := 0; attempt < 3; attempt += 1 {
		family.FamilyUID = generateFamilyUID(userID, attempt)
		if _, err := session.Insert(family); err != nil {
			if isDuplicateEntry(err) {
				continue
			}
			_ = session.Rollback()
			response.Fail(c, http.StatusInternalServerError, "create family failed")
			return
		}
		inserted = true
		break
	}

	if !inserted {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "generate family id failed")
		return
	}

	if _, err := session.Insert(&model.FamilyMember{
		FamilyID: family.ID,
		UserID:   userID,
	}); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "create family member failed")
		return
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "commit transaction failed")
		return
	}

	item, err := h.buildFamilyListItemByUID(c, family.FamilyUID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family detail failed")
		return
	}

	response.Success(c, gin.H{
		"family": item,
	})
}

func (h *UserFamilyHandler) Join(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	payload, err := parseJSONMap(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid join payload")
		return
	}

	familyUID := normalizeFamilyUID(readString(payload, "family_id", "familyId", "family_uid", "id"))
	if familyUID == "" {
		response.Fail(c, http.StatusBadRequest, "family_id is required")
		return
	}

	item, err := h.joinFamily(c, userID, familyUID)
	if err != nil {
		h.failWithFamilyError(c, err)
		return
	}

	response.Success(c, gin.H{
		"family": item,
	})
}

func (h *UserFamilyHandler) JoinByLink(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	payload, err := parseJSONMap(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid join payload")
		return
	}

	inviteLink := strings.TrimSpace(readString(payload, "invite_link", "inviteLink", "link", "url"))
	familyUID := extractFamilyUIDFromInvite(inviteLink)
	if familyUID == "" {
		response.Fail(c, http.StatusBadRequest, "invite link is invalid")
		return
	}

	item, err := h.joinFamily(c, userID, familyUID)
	if err != nil {
		h.failWithFamilyError(c, err)
		return
	}

	response.Success(c, gin.H{
		"family": item,
	})
}

func (h *UserFamilyHandler) Leave(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	familyUID := normalizeFamilyUID(c.Param("familyId"))
	if familyUID == "" {
		response.Fail(c, http.StatusBadRequest, "family id is invalid")
		return
	}

	scope, has, err := h.loadAccessibleFamily(userID, familyUID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "family not found")
		return
	}

	session := h.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		response.Fail(c, http.StatusInternalServerError, "begin transaction failed")
		return
	}

	if scope.CreatorUserID == userID {
		if len(scope.Members) == 1 {
			if _, err := session.ID(scope.ID).Delete(new(model.Family)); err != nil {
				_ = session.Rollback()
				response.Fail(c, http.StatusInternalServerError, "delete family failed")
				return
			}
		} else {
			nextCreator := &model.FamilyMember{}
			hasNextCreator, err := session.Where("family_id = ? AND user_id <> ?", scope.ID, userID).
				Asc("joined_at").
				Asc("user_id").
				Get(nextCreator)
			if err != nil {
				_ = session.Rollback()
				response.Fail(c, http.StatusInternalServerError, "query next family creator failed")
				return
			}
			if !hasNextCreator {
				_ = session.Rollback()
				response.Fail(c, http.StatusInternalServerError, "next family creator not found")
				return
			}

			if _, err := session.Table(new(model.Family)).
				Where("id = ?", scope.ID).
				Update(map[string]any{
					"creator_user_id": nextCreator.UserID,
					"updated_at":      time.Now(),
				}); err != nil {
				_ = session.Rollback()
				response.Fail(c, http.StatusInternalServerError, "transfer family creator failed")
				return
			}

			if _, err := session.Where("family_id = ? AND user_id = ?", scope.ID, userID).Delete(new(model.FamilyMember)); err != nil {
				_ = session.Rollback()
				response.Fail(c, http.StatusInternalServerError, "leave family failed")
				return
			}
		}
	} else {
		if _, err := session.Where("family_id = ? AND user_id = ?", scope.ID, userID).Delete(new(model.FamilyMember)); err != nil {
			_ = session.Rollback()
			response.Fail(c, http.StatusInternalServerError, "leave family failed")
			return
		}
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "commit transaction failed")
		return
	}

	response.Success(c, gin.H{
		"familyId": scope.UID,
		"left":     true,
		"deleted":  scope.CreatorUserID == userID && len(scope.Members) == 1,
	})
}

func (h *UserFamilyHandler) GetDetail(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	familyUID := normalizeFamilyUID(c.Param("familyId"))
	if familyUID == "" {
		response.Fail(c, http.StatusBadRequest, "family id is invalid")
		return
	}

	scope, has, err := h.loadAccessibleFamily(userID, familyUID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "family not found")
		return
	}

	monthOptions, yearOptions, err := h.loadFamilyPeriodOptions(scope.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family summary failed")
		return
	}

	response.Success(c, gin.H{
		"id":           scope.UID,
		"name":         scope.Name,
		"slogan":       buildFamilySlogan(scope.Name, scope.CreatedAt),
		"creator":      scope.CreatorName,
		"createdAt":    scope.CreatedAt.Format("2006-01-02"),
		"inviteCode":   scope.UID,
		"inviteLink":   buildFamilyInviteLink(c, scope.UID),
		"members":      scope.Members,
		"monthOptions": monthOptions,
		"yearOptions":  yearOptions,
	})
}

func (h *UserFamilyHandler) GetMemberShare(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	familyUID := normalizeFamilyUID(c.Param("familyId"))
	if familyUID == "" {
		response.Fail(c, http.StatusBadRequest, "family id is invalid")
		return
	}

	scope, has, err := h.loadAccessibleFamily(userID, familyUID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "family not found")
		return
	}

	periodType := normalizeFamilyPeriodType(c.Query("periodType"), c.Query("period_type"))
	if periodType == "" {
		response.Fail(c, http.StatusBadRequest, "periodType must be month or year")
		return
	}

	metricType := normalizeRecordTypeValue(c.Query("metricType"), c.Query("metric_type"))
	if metricType == "" {
		response.Fail(c, http.StatusBadRequest, "metricType must be income or expense")
		return
	}

	periodKey := strings.TrimSpace(firstNonEmpty(c.Query("periodKey"), c.Query("period_key")))
	if periodKey == "" {
		now := time.Now()
		if periodType == "year" {
			periodKey = now.Format("2006")
		} else {
			periodKey = now.Format("2006-01")
		}
	}

	rows, total, err := h.loadFamilyShareRows(scope, periodType, periodKey, metricType)
	if err != nil {
		h.failWithFamilyError(c, err)
		return
	}

	response.Success(c, gin.H{
		"title": buildFamilyShareTitle(periodKey, periodType, metricType),
		"total": total,
		"rows":  rows,
	})
}

func (h *UserFamilyHandler) joinFamily(c *gin.Context, userID uint64, familyUID string) (gin.H, error) {
	family := &model.Family{}
	has, err := h.engine.Where("family_uid = ?", familyUID).Get(family)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, familyHTTPError{status: http.StatusNotFound, message: "family not found"}
	}

	session := h.engine.NewSession()
	defer func() {
		_ = session.Close()
	}()

	if err := session.Begin(); err != nil {
		return nil, err
	}

	exists, err := session.Where("family_id = ? AND user_id = ?", family.ID, userID).Exist(new(model.FamilyMember))
	if err != nil {
		_ = session.Rollback()
		return nil, err
	}
	if exists {
		_ = session.Rollback()
		return nil, familyHTTPError{status: http.StatusConflict, message: "user already joined this family"}
	}

	if _, err := session.Insert(&model.FamilyMember{
		FamilyID: family.ID,
		UserID:   userID,
	}); err != nil {
		_ = session.Rollback()
		if isDuplicateEntry(err) {
			return nil, familyHTTPError{status: http.StatusConflict, message: "user already joined this family"}
		}
		return nil, err
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		return nil, err
	}

	return h.buildFamilyListItemByUID(c, familyUID)
}

func (h *UserFamilyHandler) buildFamilyListItemByUID(c *gin.Context, familyUID string) (gin.H, error) {
	base := &familyBaseRow{}
	has, err := h.engine.SQL(`
		SELECT
			f.id,
			f.family_uid,
			f.family_name,
			f.creator_user_id,
			f.created_at,
			cu.nickname AS creator_nickname,
			cu.username AS creator_username
		FROM families f
		INNER JOIN users cu ON cu.id = f.creator_user_id
		WHERE f.family_uid = ?
		LIMIT 1
	`, familyUID).Get(base)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, familyHTTPError{status: http.StatusNotFound, message: "family not found"}
	}

	scope, err := h.loadFamilyScopeByID(base.ID)
	if err != nil {
		return nil, err
	}

	monthOptions, yearOptions, err := h.loadFamilyPeriodOptions(scope.ID)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"id":           scope.UID,
		"name":         scope.Name,
		"slogan":       buildFamilySlogan(scope.Name, scope.CreatedAt),
		"creator":      scope.CreatorName,
		"createdAt":    scope.CreatedAt.Format("2006-01-02"),
		"inviteCode":   scope.UID,
		"inviteLink":   buildFamilyInviteLink(c, scope.UID),
		"memberCount":  len(scope.Members),
		"monthIncome":  monthOptions[0].Income,
		"monthExpense": monthOptions[0].Expense,
		"monthBalance": monthOptions[0].Balance,
		"yearIncome":   yearOptions[0].Income,
		"yearExpense":  yearOptions[0].Expense,
		"yearBalance":  yearOptions[0].Balance,
		"members":      scope.Members,
	}, nil
}

func (h *UserFamilyHandler) loadAccessibleFamily(userID uint64, familyUID string) (familyScope, bool, error) {
	row := &familyBaseRow{}
	has, err := h.engine.SQL(`
		SELECT
			f.id,
			f.family_uid,
			f.family_name,
			f.creator_user_id,
			f.created_at,
			cu.nickname AS creator_nickname,
			cu.username AS creator_username
		FROM families f
		INNER JOIN family_members self_fm ON self_fm.family_id = f.id AND self_fm.user_id = ?
		INNER JOIN users cu ON cu.id = f.creator_user_id
		WHERE f.family_uid = ?
		LIMIT 1
	`, userID, familyUID).Get(row)
	if err != nil {
		return familyScope{}, false, err
	}
	if !has {
		return familyScope{}, false, nil
	}

	scope, err := h.loadFamilyScopeByID(row.ID)
	if err != nil {
		return familyScope{}, false, err
	}
	return scope, true, nil
}

func (h *UserFamilyHandler) loadFamilyScopeByID(familyID uint64) (familyScope, error) {
	row := &familyBaseRow{}
	has, err := h.engine.SQL(`
		SELECT
			f.id,
			f.family_uid,
			f.family_name,
			f.creator_user_id,
			f.created_at,
			cu.nickname AS creator_nickname,
			cu.username AS creator_username
		FROM families f
		INNER JOIN users cu ON cu.id = f.creator_user_id
		WHERE f.id = ?
		LIMIT 1
	`, familyID).Get(row)
	if err != nil {
		return familyScope{}, err
	}
	if !has {
		return familyScope{}, familyHTTPError{status: http.StatusNotFound, message: "family not found"}
	}

	memberRows := make([]familyMemberRow, 0)
	if err := h.engine.SQL(`
		SELECT
			u.id AS user_id,
			u.nickname,
			u.username,
			fm.joined_at
		FROM family_members fm
		INNER JOIN users u ON u.id = fm.user_id
		WHERE fm.family_id = ?
		ORDER BY fm.joined_at ASC, fm.id ASC
	`, familyID).Find(&memberRows); err != nil {
		return familyScope{}, err
	}

	members := make([]familyMemberInfo, 0, len(memberRows))
	for _, item := range memberRows {
		members = append(members, familyMemberInfo{
			UserID:   item.UserID,
			Name:     displayUserName(item.Nickname, item.Username, item.UserID),
			Role:     familyMemberRole(item.UserID == row.CreatorUserID),
			Color:    familyMemberColor(item.UserID, row.FamilyUID),
			JoinedAt: item.JoinedAt.Format("2006-01-02"),
		})
	}

	sort.SliceStable(members, func(i, j int) bool {
		if members[i].Role != members[j].Role {
			return members[i].Role == "创建人"
		}
		if members[i].JoinedAt != members[j].JoinedAt {
			return members[i].JoinedAt < members[j].JoinedAt
		}
		return members[i].UserID < members[j].UserID
	})

	return familyScope{
		ID:            row.ID,
		UID:           row.FamilyUID,
		Name:          row.FamilyName,
		CreatorUserID: row.CreatorUserID,
		CreatorName:   displayUserName(row.CreatorNickname, row.CreatorUsername, row.CreatorUserID),
		CreatedAt:     row.CreatedAt,
		Members:       members,
	}, nil
}

func (h *UserFamilyHandler) loadFamilyPeriodOptions(familyID uint64) ([]familyPeriodOption, []familyPeriodOption, error) {
	monthRows := make([]familyAggregateRow, 0)
	if err := h.engine.SQL(`
		SELECT
			DATE_FORMAT(lr.record_date, '%Y-%m') AS period_key,
			lr.user_id,
			SUM(CASE WHEN lr.record_type = 'income' THEN lr.amount ELSE 0 END) AS income,
			SUM(CASE WHEN lr.record_type = 'expense' THEN lr.amount ELSE 0 END) AS expense
		FROM family_members fm
		INNER JOIN ledger_records lr ON lr.user_id = fm.user_id
		WHERE fm.family_id = ?
		GROUP BY DATE_FORMAT(lr.record_date, '%Y-%m'), lr.user_id
	`, familyID).Find(&monthRows); err != nil {
		return nil, nil, err
	}

	yearRows := make([]familyAggregateRow, 0)
	if err := h.engine.SQL(`
		SELECT
			DATE_FORMAT(lr.record_date, '%Y') AS period_key,
			lr.user_id,
			SUM(CASE WHEN lr.record_type = 'income' THEN lr.amount ELSE 0 END) AS income,
			SUM(CASE WHEN lr.record_type = 'expense' THEN lr.amount ELSE 0 END) AS expense
		FROM family_members fm
		INNER JOIN ledger_records lr ON lr.user_id = fm.user_id
		WHERE fm.family_id = ?
		GROUP BY DATE_FORMAT(lr.record_date, '%Y'), lr.user_id
	`, familyID).Find(&yearRows); err != nil {
		return nil, nil, err
	}

	now := time.Now()
	monthOptions := buildFamilyPeriodOptions(monthRows, now.Format("2006-01"), "month")
	yearOptions := buildFamilyPeriodOptions(yearRows, now.Format("2006"), "year")
	return monthOptions, yearOptions, nil
}

func (h *UserFamilyHandler) loadFamilyShareRows(scope familyScope, periodType string, periodKey string, metricType string) ([]familyShareRow, float64, error) {
	layout := "2006-01"
	startText := periodKey
	if periodType == "year" {
		layout = "2006"
	}

	start, err := time.ParseInLocation(layout, startText, time.Local)
	if err != nil {
		return nil, 0, familyHTTPError{status: http.StatusBadRequest, message: "periodKey is invalid"}
	}

	end := start.AddDate(0, 1, 0)
	if periodType == "year" {
		end = start.AddDate(1, 0, 0)
	}

	rows := make([]familyAggregateRow, 0)
	if err := h.engine.SQL(`
		SELECT
			? AS period_key,
			fm.user_id,
			SUM(CASE WHEN lr.record_type = 'income' THEN lr.amount ELSE 0 END) AS income,
			SUM(CASE WHEN lr.record_type = 'expense' THEN lr.amount ELSE 0 END) AS expense
		FROM family_members fm
		LEFT JOIN ledger_records lr
			ON lr.user_id = fm.user_id
			AND lr.record_date >= ?
			AND lr.record_date < ?
		WHERE fm.family_id = ?
		GROUP BY fm.user_id
	`, periodKey, start.Format("2006-01-02"), end.Format("2006-01-02"), scope.ID).Find(&rows); err != nil {
		return nil, 0, err
	}

	valueByUser := make(map[uint64]float64, len(rows))
	for _, row := range rows {
		if metricType == recordTypeIncome {
			valueByUser[row.UserID] = round2(row.Income)
			continue
		}
		valueByUser[row.UserID] = round2(row.Expense)
	}

	shareRows := make([]familyShareRow, 0, len(scope.Members))
	total := 0.0
	for _, member := range scope.Members {
		value := round2(valueByUser[member.UserID])
		total += value
		shareRows = append(shareRows, familyShareRow{
			Name:  member.Name,
			Role:  member.Role,
			Color: member.Color,
			Value: value,
		})
	}

	total = round2(total)
	for index := range shareRows {
		if total == 0 {
			shareRows[index].Percent = 0
			continue
		}
		shareRows[index].Percent = round2((shareRows[index].Value / total) * 100)
	}

	sort.SliceStable(shareRows, func(i, j int) bool {
		if shareRows[i].Value != shareRows[j].Value {
			return shareRows[i].Value > shareRows[j].Value
		}
		return shareRows[i].Name < shareRows[j].Name
	})

	return shareRows, total, nil
}

func buildFamilyPeriodOptions(rows []familyAggregateRow, fallbackKey string, periodType string) []familyPeriodOption {
	type totals struct {
		income  float64
		expense float64
	}

	byPeriod := make(map[string]*totals)
	for _, row := range rows {
		if strings.TrimSpace(row.PeriodKey) == "" {
			continue
		}
		bucket, exists := byPeriod[row.PeriodKey]
		if !exists {
			bucket = &totals{}
			byPeriod[row.PeriodKey] = bucket
		}
		bucket.income += row.Income
		bucket.expense += row.Expense
	}

	if _, exists := byPeriod[fallbackKey]; !exists {
		byPeriod[fallbackKey] = &totals{}
	}

	keys := make([]string, 0, len(byPeriod))
	for key := range byPeriod {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	options := make([]familyPeriodOption, 0, len(keys))
	for _, key := range keys {
		bucket := byPeriod[key]
		income := round2(bucket.income)
		expense := round2(bucket.expense)
		options = append(options, familyPeriodOption{
			Key:     key,
			Label:   buildFamilyPeriodLabel(periodType, key),
			Note:    buildFamilyPeriodNote(periodType, key),
			Income:  income,
			Expense: expense,
			Balance: round2(income - expense),
		})
	}
	return options
}

func buildFamilyPeriodLabel(periodType string, key string) string {
	if periodType == "year" {
		return key + "年"
	}

	parts := strings.Split(key, "-")
	if len(parts) != 2 {
		return key
	}
	month := strings.TrimLeft(parts[1], "0")
	if month == "" {
		month = "0"
	}
	return parts[0] + "年" + month + "月"
}

func buildFamilyPeriodNote(periodType string, key string) string {
	if periodType == "year" {
		return key + " 年家庭成员真实流水汇总。"
	}
	return key + " 家庭成员真实流水汇总。"
}

func buildFamilyShareTitle(periodKey string, periodType string, metricType string) string {
	metricText := "支出"
	if metricType == recordTypeIncome {
		metricText = "收入"
	}

	periodText := "月度"
	if periodType == "year" {
		periodText = "年度"
	}

	return fmt.Sprintf("%s %s%s成员占比", periodKey, periodText, metricText)
}

func buildFamilySlogan(familyName string, createdAt time.Time) string {
	return fmt.Sprintf("%s 创建于 %s，当前展示家庭成员真实记账汇总。", familyName, createdAt.Format("2006-01-02"))
}

func buildFamilyInviteLink(c *gin.Context, familyUID string) string {
	base := strings.TrimSpace(c.GetHeader("Origin"))
	if base == "" {
		base = familyInviteFallbackBaseURL
	}
	base = strings.TrimRight(base, "/")
	return base + "/user/families?invite=" + url.QueryEscape(familyUID)
}

func generateFamilyUID(userID uint64, attempt int) string {
	nowToken := strings.ToUpper(fmt.Sprintf("%X", time.Now().UnixNano()))
	if attempt == 0 {
		return fmt.Sprintf("FAM-%s-%d", nowToken, userID)
	}
	return fmt.Sprintf("FAM-%s-%d-%d", nowToken, userID, attempt)
}

func normalizeFamilyUID(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func extractFamilyUIDFromInvite(raw string) string {
	text := strings.TrimSpace(raw)
	if text == "" {
		return ""
	}

	if matched := familyInviteCodePattern.FindString(text); matched != "" {
		return normalizeFamilyUID(matched)
	}

	parsed, err := url.Parse(text)
	if err != nil {
		return ""
	}

	for _, key := range []string{"invite", "familyId", "family_id", "family_uid", "code"} {
		if value := normalizeFamilyUID(parsed.Query().Get(key)); value != "" {
			if matched := familyInviteCodePattern.FindString(value); matched != "" {
				return normalizeFamilyUID(matched)
			}
		}
	}

	segment := normalizeFamilyUID(pathLastSegment(parsed.Path))
	if matched := familyInviteCodePattern.FindString(segment); matched != "" {
		return normalizeFamilyUID(matched)
	}
	return ""
}

func pathLastSegment(path string) string {
	trimmed := strings.Trim(strings.TrimSpace(path), "/")
	if trimmed == "" {
		return ""
	}

	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}

func displayUserName(nickname string, username string, userID uint64) string {
	if strings.TrimSpace(nickname) != "" {
		return strings.TrimSpace(nickname)
	}
	if strings.TrimSpace(username) != "" {
		return strings.TrimSpace(username)
	}
	return fmt.Sprintf("用户%d", userID)
}

func familyMemberRole(isCreator bool) string {
	if isCreator {
		return "创建人"
	}
	return "成员"
}

func familyMemberColor(userID uint64, familyUID string) string {
	palette := []string{"#f6d34a", "#6bcf7c", "#4d96ff", "#ff8b8b", "#9b8cff", "#34d399", "#f97316"}
	seed := 0
	for _, char := range familyUID {
		seed += int(char)
	}
	seed += int(userID % uint64(len(palette)*17))
	return palette[seed%len(palette)]
}

func normalizeFamilyPeriodType(values ...string) string {
	for _, value := range values {
		switch strings.ToLower(strings.TrimSpace(value)) {
		case "month":
			return "month"
		case "year":
			return "year"
		}
	}
	return ""
}

func normalizeRecordTypeValue(values ...string) string {
	for _, value := range values {
		normalized, ok := normalizeRecordType(value)
		if ok {
			return normalized
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

type familyHTTPError struct {
	status  int
	message string
}

func (e familyHTTPError) Error() string {
	return e.message
}

func (h *UserFamilyHandler) failWithFamilyError(c *gin.Context, err error) {
	if familyErr, ok := err.(familyHTTPError); ok {
		response.Fail(c, familyErr.status, familyErr.message)
		return
	}
	response.Fail(c, http.StatusInternalServerError, "family operation failed")
}
