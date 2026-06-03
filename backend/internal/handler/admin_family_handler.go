package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
)

const adminFamilyActiveBillDays = 30

type AdminFamilyHandler struct {
	engine *xorm.Engine
}

type adminFamilyListRow struct {
	ID              uint64    `xorm:"id"`
	FamilyUID       string    `xorm:"family_uid"`
	FamilyName      string    `xorm:"family_name"`
	CreatorUserID   uint64    `xorm:"creator_user_id"`
	CreatorNickname string    `xorm:"creator_nickname"`
	CreatorUsername string    `xorm:"creator_username"`
	CreatedAt       time.Time `xorm:"created_at"`
	MemberCount     int64     `xorm:"member_count"`
	BillCount       int64     `xorm:"bill_count"`
	TotalAssets     float64   `xorm:"total_assets"`
}

type adminFamilyPeriodSummary struct {
	Income  float64 `xorm:"income" json:"income"`
	Expense float64 `xorm:"expense" json:"expense"`
	Balance float64 `xorm:"balance" json:"balance"`
	Records int64   `xorm:"records" json:"records"`
}

type adminFamilyAssetSummary struct {
	TotalAssets      float64 `xorm:"total_assets" json:"total_assets"`
	TotalLiabilities float64 `xorm:"total_liabilities" json:"total_liabilities"`
	NetAssets        float64 `json:"net_assets"`
	AccountCount     int64   `xorm:"account_count" json:"account_count"`
}

type adminFamilyRecentBill struct {
	ID         uint64    `xorm:"id"`
	UserID     uint64    `xorm:"user_id"`
	MemberName string    `xorm:"member_name"`
	RecordType string    `xorm:"record_type"`
	Amount     float64   `xorm:"amount"`
	Remark     string    `xorm:"remark"`
	RecordDate time.Time `xorm:"record_date"`
	CreatedAt  time.Time `xorm:"created_at"`
}

func NewAdminFamilyHandler(engine *xorm.Engine) *AdminFamilyHandler {
	return &AdminFamilyHandler{engine: engine}
}

func (h *AdminFamilyHandler) Summary(c *gin.Context) {
	totalFamilies, err := h.engine.Count(new(model.Family))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "count families failed")
		return
	}

	totalMembers, err := h.engine.Count(new(model.FamilyMember))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "count family members failed")
		return
	}

	activeFamilies, err := h.countActiveFamilies()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "count active families failed")
		return
	}

	averageMembers := 0.0
	if totalFamilies > 0 {
		averageMembers = round2(float64(totalMembers) / float64(totalFamilies))
	}

	response.Success(c, gin.H{
		"total_families":  totalFamilies,
		"total_members":   totalMembers,
		"average_members": averageMembers,
		"active_families": activeFamilies,
		"active_days":     adminFamilyActiveBillDays,
	})
}

func (h *AdminFamilyHandler) List(c *gin.Context) {
	page := parseAdminFamilyPositiveInt(c.Query("page"), 1)
	pageSize := parseAdminFamilyPositiveInt(c.Query("page_size"), 10)
	if pageSize > 50 {
		pageSize = 50
	}
	keyword := strings.TrimSpace(c.Query("keyword"))

	total, err := h.countFamiliesByKeyword(keyword)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "count families failed")
		return
	}

	rows := make([]adminFamilyListRow, 0)
	query, args := buildAdminFamilyListSQL(keyword, pageSize, (page-1)*pageSize)
	if err := h.engine.SQL(query, args...).Find(&rows); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query families failed")
		return
	}

	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, buildAdminFamilyListItem(row))
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	response.Success(c, gin.H{
		"list": items,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func (h *AdminFamilyHandler) Detail(c *gin.Context) {
	familyUID := normalizeFamilyUID(c.Param("familyId"))
	if familyUID == "" {
		response.Fail(c, http.StatusBadRequest, "family id is invalid")
		return
	}

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
		response.Fail(c, http.StatusInternalServerError, "query family failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "family not found")
		return
	}

	scope, err := h.loadAdminFamilyScope(base.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family members failed")
		return
	}

	incomeExpense, err := h.loadAdminFamilyIncomeExpense(base.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family income expense failed")
		return
	}

	assetSummary, err := h.loadAdminFamilyAssets(base.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family assets failed")
		return
	}

	recentBills, err := h.loadAdminFamilyRecentBills(base.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query family recent bills failed")
		return
	}

	response.Success(c, gin.H{
		"family": gin.H{
			"id":            scope.UID,
			"family_uid":    scope.UID,
			"name":          scope.Name,
			"family_name":   scope.Name,
			"creator":       scope.CreatorName,
			"creatorUserId": scope.CreatorUserID,
			"memberCount":   len(scope.Members),
			"createdAt":     scope.CreatedAt.Format("2006-01-02"),
			"status":        buildAdminFamilyStatus(len(scope.Members), incomeExpense.Records),
		},
		"members":        scope.Members,
		"income_expense": incomeExpense,
		"asset_summary":  assetSummary,
		"recent_bills":   recentBills,
	})
}

func (h *AdminFamilyHandler) countActiveFamilies() (int64, error) {
	startDate := time.Now().AddDate(0, 0, -adminFamilyActiveBillDays).Format("2006-01-02")
	type row struct {
		Total int64 `xorm:"total"`
	}
	result := &row{}
	_, err := h.engine.SQL(`
		SELECT COUNT(DISTINCT fm.family_id) AS total
		FROM family_members fm
		INNER JOIN ledger_records lr ON lr.user_id = fm.user_id
		WHERE lr.record_date >= ?
	`, startDate).Get(result)
	return result.Total, err
}

func (h *AdminFamilyHandler) countFamiliesByKeyword(keyword string) (int64, error) {
	query := `
		SELECT COUNT(DISTINCT f.id) AS total
		FROM families f
		INNER JOIN users cu ON cu.id = f.creator_user_id
	`
	args := make([]any, 0)
	if keyword != "" {
		query += `
			WHERE f.family_uid LIKE ?
				OR f.family_name LIKE ?
				OR cu.nickname LIKE ?
				OR cu.username LIKE ?
		`
		like := "%" + keyword + "%"
		args = append(args, like, like, like, like)
	}

	type row struct {
		Total int64 `xorm:"total"`
	}
	result := &row{}
	_, err := h.engine.SQL(query, args...).Get(result)
	return result.Total, err
}

func (h *AdminFamilyHandler) loadAdminFamilyScope(familyID uint64) (familyScope, error) {
	userFamilyHandler := NewUserFamilyHandler(h.engine)
	return userFamilyHandler.loadFamilyScopeByID(familyID)
}

func (h *AdminFamilyHandler) loadAdminFamilyIncomeExpense(familyID uint64) (adminFamilyPeriodSummary, error) {
	summary := adminFamilyPeriodSummary{}
	_, err := h.engine.SQL(`
		SELECT
			COALESCE(SUM(CASE WHEN lr.record_type = 'income' THEN lr.amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN lr.record_type = 'expense' THEN lr.amount ELSE 0 END), 0) AS expense,
			COALESCE(SUM(CASE WHEN lr.record_type = 'income' THEN lr.amount ELSE -lr.amount END), 0) AS balance,
			COUNT(lr.id) AS records
		FROM family_members fm
		LEFT JOIN ledger_records lr ON lr.user_id = fm.user_id
		WHERE fm.family_id = ?
	`, familyID).Get(&summary)
	summary.Income = round2(summary.Income)
	summary.Expense = round2(summary.Expense)
	summary.Balance = round2(summary.Balance)
	return summary, err
}

func (h *AdminFamilyHandler) loadAdminFamilyAssets(familyID uint64) (adminFamilyAssetSummary, error) {
	summary := adminFamilyAssetSummary{}
	_, err := h.engine.SQL(`
		SELECT
			COALESCE(SUM(CASE WHEN aa.balance_nature = 'liability' THEN 0 ELSE aa.balance END), 0) AS total_assets,
			COALESCE(SUM(CASE WHEN aa.balance_nature = 'liability' THEN aa.balance ELSE 0 END), 0) AS total_liabilities,
			COUNT(aa.id) AS account_count
		FROM family_members fm
		LEFT JOIN asset_accounts aa ON aa.user_id = fm.user_id
		WHERE fm.family_id = ?
	`, familyID).Get(&summary)
	summary.TotalAssets = round2(summary.TotalAssets)
	summary.TotalLiabilities = round2(summary.TotalLiabilities)
	summary.NetAssets = round2(summary.TotalAssets - summary.TotalLiabilities)
	return summary, err
}

func (h *AdminFamilyHandler) loadAdminFamilyRecentBills(familyID uint64) ([]gin.H, error) {
	rows := make([]adminFamilyRecentBill, 0)
	if err := h.engine.SQL(`
		SELECT
			lr.id,
			lr.user_id,
			COALESCE(NULLIF(u.nickname, ''), NULLIF(u.username, ''), CONCAT('用户', u.id)) AS member_name,
			lr.record_type,
			lr.amount,
			lr.remark,
			lr.record_date,
			lr.created_at
		FROM family_members fm
		INNER JOIN ledger_records lr ON lr.user_id = fm.user_id
		INNER JOIN users u ON u.id = lr.user_id
		WHERE fm.family_id = ?
		ORDER BY lr.record_date DESC, lr.id DESC
		LIMIT 8
	`, familyID).Find(&rows); err != nil {
		return nil, err
	}

	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{
			"id":         row.ID,
			"userId":     row.UserID,
			"memberName": row.MemberName,
			"recordType": row.RecordType,
			"amount":     round2(row.Amount),
			"remark":     row.Remark,
			"recordDate": row.RecordDate.Format("2006-01-02"),
			"createdAt":  row.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	return items, nil
}

func buildAdminFamilyListSQL(keyword string, limit int, offset int) (string, []any) {
	query := `
		SELECT
			f.id,
			f.family_uid,
			f.family_name,
			f.creator_user_id,
			f.created_at,
			cu.nickname AS creator_nickname,
			cu.username AS creator_username,
			COALESCE(member_stats.member_count, 0) AS member_count,
			COALESCE(bill_stats.bill_count, 0) AS bill_count,
			COALESCE(asset_stats.total_assets, 0) AS total_assets
		FROM families f
		INNER JOIN users cu ON cu.id = f.creator_user_id
		LEFT JOIN (
			SELECT family_id, COUNT(DISTINCT user_id) AS member_count
			FROM family_members
			GROUP BY family_id
		) member_stats ON member_stats.family_id = f.id
		LEFT JOIN (
			SELECT fm.family_id, COUNT(lr.id) AS bill_count
			FROM family_members fm
			INNER JOIN ledger_records lr ON lr.user_id = fm.user_id
			GROUP BY fm.family_id
		) bill_stats ON bill_stats.family_id = f.id
		LEFT JOIN (
			SELECT
				fm.family_id,
				SUM(CASE WHEN aa.balance_nature = 'liability' THEN -aa.balance ELSE aa.balance END) AS total_assets
			FROM family_members fm
			INNER JOIN asset_accounts aa ON aa.user_id = fm.user_id
			GROUP BY fm.family_id
		) asset_stats ON asset_stats.family_id = f.id
	`
	args := make([]any, 0)
	if keyword != "" {
		query += `
			WHERE f.family_uid LIKE ?
				OR f.family_name LIKE ?
				OR cu.nickname LIKE ?
				OR cu.username LIKE ?
		`
		like := "%" + keyword + "%"
		args = append(args, like, like, like, like)
	}
	query += `
		ORDER BY f.created_at DESC, f.id DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)
	return query, args
}

func buildAdminFamilyListItem(row adminFamilyListRow) gin.H {
	return gin.H{
		"id":            row.FamilyUID,
		"family_uid":    row.FamilyUID,
		"name":          row.FamilyName,
		"family_name":   row.FamilyName,
		"creator":       displayUserName(row.CreatorNickname, row.CreatorUsername, row.CreatorUserID),
		"creatorUserId": row.CreatorUserID,
		"members":       row.MemberCount,
		"memberCount":   row.MemberCount,
		"billCount":     row.BillCount,
		"totalAssets":   round2(row.TotalAssets),
		"createdAt":     row.CreatedAt.Format("2006-01-02"),
		"status":        buildAdminFamilyStatus(int(row.MemberCount), row.BillCount),
	}
}

func buildAdminFamilyStatus(memberCount int, billCount int64) string {
	if memberCount <= 0 {
		return "空家庭"
	}
	if billCount > 0 {
		return "活跃"
	}
	return "待记账"
}

func parseAdminFamilyPositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
