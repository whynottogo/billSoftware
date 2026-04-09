package handler

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
)

const (
	recordTypeIncome  = "income"
	recordTypeExpense = "expense"
)

type UserLedgerHandler struct {
	engine *xorm.Engine
}

type createLedgerRequest struct {
	Type       string  `json:"type"`
	CategoryID uint64  `json:"category_id"`
	Amount     float64 `json:"amount"`
	Remark     string  `json:"remark"`
	Note       string  `json:"note"`
	RecordDate string  `json:"record_date"`
	ImageURL   string  `json:"image_url"`
}

type createCategoryRequest struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Badge string `json:"badge"`
}

type ledgerRecordRow struct {
	ID           uint64    `xorm:"id"`
	RecordType   string    `xorm:"record_type"`
	Amount       float64   `xorm:"amount"`
	Remark       string    `xorm:"remark"`
	RecordDate   time.Time `xorm:"record_date"`
	ImageURL     string    `xorm:"image_url"`
	CreatedAt    time.Time `xorm:"created_at"`
	CategoryName string    `xorm:"category_name"`
}

type categoryListItem struct {
	ID        uint64 `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Badge     string `json:"badge"`
	IsDefault bool   `json:"isDefault"`
}

type monthSummary struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type monthLedgerItem struct {
	ID        uint64  `json:"id"`
	Badge     string  `json:"badge"`
	Category  string  `json:"category"`
	Time      string  `json:"time"`
	Note      string  `json:"note"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
	ImageName string  `json:"imageName"`
}

type monthLedgerGroup struct {
	Date         string           `json:"date"`
	Weekday      string           `json:"weekday"`
	TotalIncome  float64          `json:"totalIncome"`
	TotalExpense float64          `json:"totalExpense"`
	Items        []monthLedgerItem `json:"items"`
}

type monthCategoryUsage struct {
	Name  string `json:"name"`
	Badge string `json:"badge"`
	Count int    `json:"count"`
}

type monthOverview struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Progress int    `json:"progress"`
}

func NewUserLedgerHandler(engine *xorm.Engine) *UserLedgerHandler {
	return &UserLedgerHandler{engine: engine}
}

func (h *UserLedgerHandler) GetLedger(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	monthStart, monthEnd, monthKey, err := resolveMonthRange(c.Query("month"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "month must be YYYY-MM")
		return
	}

	rows := make([]ledgerRecordRow, 0)
	if err := h.engine.SQL(`
		SELECT
			lr.id,
			lr.record_type,
			lr.amount,
			lr.remark,
			lr.record_date,
			lr.image_url,
			lr.created_at,
			COALESCE(uc.name, '') AS category_name
		FROM ledger_records lr
		LEFT JOIN user_categories uc ON uc.id = lr.category_id AND uc.user_id = lr.user_id
		WHERE lr.user_id = ? AND lr.record_date >= ? AND lr.record_date < ?
		ORDER BY lr.record_date DESC, lr.created_at DESC, lr.id DESC
	`, userID, monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02")).Find(&rows); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query ledger records failed")
		return
	}

	categories, err := h.loadUserCategories(userID, "")
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query categories failed")
		return
	}

	response.Success(c, buildMonthLedgerResponse(monthKey, rows, categories))
}

func (h *UserLedgerHandler) CreateLedger(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	var req createLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid ledger payload")
		return
	}

	recordType, ok := normalizeRecordType(req.Type)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "type must be income or expense")
		return
	}

	if req.CategoryID == 0 {
		response.Fail(c, http.StatusBadRequest, "category_id is required")
		return
	}

	if req.Amount <= 0 {
		response.Fail(c, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	recordDate, err := resolveRecordDate(req.RecordDate)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "record_date must be YYYY-MM-DD")
		return
	}

	category := &model.UserCategory{}
	has, err := h.engine.Where("id = ? AND user_id = ?", req.CategoryID, userID).Get(category)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query category failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusBadRequest, "category does not exist")
		return
	}
	if category.CategoryType != recordType {
		response.Fail(c, http.StatusBadRequest, "category type does not match ledger type")
		return
	}

	remark := strings.TrimSpace(req.Remark)
	if remark == "" {
		remark = strings.TrimSpace(req.Note)
	}

	ledger := &model.LedgerRecord{
		UserID:     userID,
		RecordType: recordType,
		CategoryID: req.CategoryID,
		Amount:     round2(req.Amount),
		Remark:     remark,
		RecordDate: recordDate,
		ImageURL:   strings.TrimSpace(req.ImageURL),
	}

	if _, err := h.engine.Insert(ledger); err != nil {
		response.Fail(c, http.StatusInternalServerError, "create ledger record failed")
		return
	}

	response.Success(c, gin.H{
		"record_id": ledger.ID,
		"month":     recordDate.Format("2006-01"),
	})
}

func (h *UserLedgerHandler) ListCategories(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	recordType, ok := normalizeRecordType(c.Query("type"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "type must be income or expense")
		return
	}

	categories, err := h.loadUserCategories(userID, recordType)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query categories failed")
		return
	}

	list := make([]categoryListItem, 0, len(categories))
	for _, item := range categories {
		list = append(list, categoryListItem{
			ID:        item.ID,
			Type:      item.CategoryType,
			Name:      item.Name,
			Badge:     badgeFromName(item.Name),
			IsDefault: item.IsSystem == 1,
		})
	}

	response.Success(c, gin.H{
		"type": recordType,
		"list": list,
	})
}

func (h *UserLedgerHandler) CreateCategory(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid category payload")
		return
	}

	recordType, ok := normalizeRecordType(req.Type)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "type must be income or expense")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		response.Fail(c, http.StatusBadRequest, "name is required")
		return
	}

	exists, err := h.engine.Where("user_id = ? AND category_type = ? AND name = ?", userID, recordType, name).Exist(&model.UserCategory{})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "check category name failed")
		return
	}
	if exists {
		response.Fail(c, http.StatusConflict, "category name already exists")
		return
	}

	var maxSort struct {
		Value int `xorm:"max_sort"`
	}
	if _, err := h.engine.SQL(`
		SELECT COALESCE(MAX(sort_order), 0) AS max_sort
		FROM user_categories
		WHERE user_id = ? AND category_type = ?
	`, userID, recordType).Get(&maxSort); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query category sort order failed")
		return
	}

	category := &model.UserCategory{
		UserID:       userID,
		CategoryType: recordType,
		Name:         name,
		SortOrder:    maxSort.Value + 1,
		IsSystem:     0,
	}

	if _, err := h.engine.Insert(category); err != nil {
		if isDuplicateEntry(err) {
			response.Fail(c, http.StatusConflict, "category name already exists")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "create category failed")
		return
	}

	response.Success(c, gin.H{
		"category": categoryListItem{
			ID:        category.ID,
			Type:      category.CategoryType,
			Name:      category.Name,
			Badge:     normalizeBadge(req.Badge, category.Name),
			IsDefault: false,
		},
	})
}

func (h *UserLedgerHandler) DeleteCategory(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid category id")
		return
	}

	category := &model.UserCategory{}
	has, err := h.engine.Where("id = ? AND user_id = ?", id, userID).Get(category)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query category failed")
		return
	}
	if !has {
		response.Fail(c, http.StatusNotFound, "category not found")
		return
	}
	if category.IsSystem == 1 {
		response.Fail(c, http.StatusBadRequest, "default category cannot be deleted")
		return
	}

	used, err := h.engine.Where("user_id = ? AND category_id = ?", userID, id).Exist(&model.LedgerRecord{})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "check category usage failed")
		return
	}
	if used {
		response.Fail(c, http.StatusBadRequest, "category already has ledger records")
		return
	}

	if _, err := h.engine.Where("id = ? AND user_id = ?", id, userID).Delete(&model.UserCategory{}); err != nil {
		response.Fail(c, http.StatusInternalServerError, "delete category failed")
		return
	}

	response.Success(c, gin.H{
		"deleted_id": id,
	})
}

func (h *UserLedgerHandler) loadUserCategories(userID uint64, recordType string) ([]model.UserCategory, error) {
	categories := make([]model.UserCategory, 0)
	session := h.engine.Where("user_id = ?", userID)
	if recordType != "" {
		session = session.And("category_type = ?", recordType)
		session = session.Asc("sort_order", "id")
	} else {
		session = session.Asc("category_type", "sort_order", "id")
	}
	if err := session.Find(&categories); err != nil {
		return nil, err
	}
	return categories, nil
}

func buildMonthLedgerResponse(monthKey string, rows []ledgerRecordRow, allCategories []model.UserCategory) gin.H {
	summary := monthSummary{
		Income:  0,
		Expense: 0,
		Balance: 0,
	}

	groupMap := make(map[string]*monthLedgerGroup)
	dateOrder := make([]string, 0)

	type usageAccumulator struct {
		Name  string
		Badge string
		Count int
	}
	usageMap := make(map[string]*usageAccumulator)
	expenseCount := 0

	for _, row := range rows {
		amount := round2(row.Amount)
		if row.RecordType == recordTypeIncome {
			summary.Income += amount
		} else {
			summary.Expense += amount
			expenseCount += 1
		}

		dateKey := row.RecordDate.Format("2006-01-02")
		group, exists := groupMap[dateKey]
		if !exists {
			group = &monthLedgerGroup{
				Date:         dateKey,
				Weekday:      weekdayLabel(dateKey),
				TotalIncome:  0,
				TotalExpense: 0,
				Items:        make([]monthLedgerItem, 0),
			}
			groupMap[dateKey] = group
			dateOrder = append(dateOrder, dateKey)
		}

		categoryName := strings.TrimSpace(row.CategoryName)
		if categoryName == "" {
			categoryName = "未分类"
		}

		badge := badgeFromName(categoryName)
		note := strings.TrimSpace(row.Remark)
		if note == "" {
			if row.RecordType == recordTypeIncome {
				note = "新增收入"
			} else {
				note = "新增支出"
			}
		}

		timeText := "00:00"
		if !row.CreatedAt.IsZero() {
			timeText = row.CreatedAt.Format("15:04")
		}

		group.Items = append(group.Items, monthLedgerItem{
			ID:        row.ID,
			Badge:     badge,
			Category:  categoryName,
			Time:      timeText,
			Note:      note,
			Amount:    amount,
			Type:      row.RecordType,
			ImageName: strings.TrimSpace(row.ImageURL),
		})

		if row.RecordType == recordTypeIncome {
			group.TotalIncome += amount
		} else {
			group.TotalExpense += amount
		}

		usage := usageMap[categoryName]
		if usage == nil {
			usage = &usageAccumulator{
				Name:  categoryName,
				Badge: badge,
				Count: 0,
			}
			usageMap[categoryName] = usage
		}
		usage.Count += 1
	}

	for _, item := range allCategories {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		if usageMap[name] == nil {
			usageMap[name] = &usageAccumulator{
				Name:  name,
				Badge: badgeFromName(name),
				Count: 0,
			}
		}
	}

	summary.Income = round2(summary.Income)
	summary.Expense = round2(summary.Expense)
	summary.Balance = round2(summary.Income - summary.Expense)

	sort.Slice(dateOrder, func(i, j int) bool {
		return dateOrder[i] > dateOrder[j]
	})

	groups := make([]monthLedgerGroup, 0, len(dateOrder))
	for _, dateKey := range dateOrder {
		group := groupMap[dateKey]
		sort.Slice(group.Items, func(i, j int) bool {
			return group.Items[i].Time > group.Items[j].Time
		})
		group.TotalIncome = round2(group.TotalIncome)
		group.TotalExpense = round2(group.TotalExpense)
		groups = append(groups, *group)
	}

	categoryUsage := make([]monthCategoryUsage, 0, len(usageMap))
	for _, item := range usageMap {
		categoryUsage = append(categoryUsage, monthCategoryUsage{
			Name:  item.Name,
			Badge: item.Badge,
			Count: item.Count,
		})
	}
	sort.Slice(categoryUsage, func(i, j int) bool {
		if categoryUsage[i].Count != categoryUsage[j].Count {
			return categoryUsage[i].Count > categoryUsage[j].Count
		}
		return categoryUsage[i].Name < categoryUsage[j].Name
	})
	if len(categoryUsage) > 6 {
		categoryUsage = categoryUsage[:6]
	}

	overview := buildMonthOverview(summary, len(dateOrder), expenseCount, categoryUsage)

	return gin.H{
		"key":        monthKey,
		"label":      monthLabel(monthKey),
		"summary":    summary,
		"groups":     groups,
		"categories": categoryUsage,
		"overview":   overview,
	}
}

func buildMonthOverview(summary monthSummary, recordedDays int, expenseCount int, categories []monthCategoryUsage) []monthOverview {
	averageExpense := 0.0
	if recordedDays > 0 {
		averageExpense = summary.Expense / float64(recordedDays)
	}

	maxCategory := monthCategoryUsage{
		Name:  "暂无",
		Badge: "-",
		Count: 0,
	}
	if len(categories) > 0 {
		maxCategory = categories[0]
	}

	categoryProgress := 20
	if maxCategory.Count > 0 {
		categoryProgress = minInt(100, maxCategory.Count*14)
	}

	categoryValue := maxCategory.Name
	if maxCategory.Count > 0 {
		categoryValue = fmt.Sprintf("%s · %d笔", maxCategory.Name, maxCategory.Count)
	}

	return []monthOverview{
		{
			Label:    "记账天数",
			Value:    fmt.Sprintf("%d天", recordedDays),
			Progress: minInt(100, recordedDays*12),
		},
		{
			Label:    "日均支出",
			Value:    formatMoney(averageExpense),
			Progress: minInt(100, expenseCount*8),
		},
		{
			Label:    "高频分类",
			Value:    categoryValue,
			Progress: categoryProgress,
		},
	}
}

func resolveMonthRange(raw string) (time.Time, time.Time, string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		return start, start.AddDate(0, 1, 0), start.Format("2006-01"), nil
	}

	month, err := time.ParseInLocation("2006-01", trimmed, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, "", err
	}

	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.Local)
	return start, start.AddDate(0, 1, 0), start.Format("2006-01"), nil
}

func resolveRecordDate(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local), nil
	}

	recordDate, err := time.ParseInLocation("2006-01-02", trimmed, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(recordDate.Year(), recordDate.Month(), recordDate.Day(), 0, 0, 0, 0, time.Local), nil
}

func normalizeRecordType(raw string) (string, bool) {
	recordType := strings.ToLower(strings.TrimSpace(raw))
	if recordType == recordTypeIncome || recordType == recordTypeExpense {
		return recordType, true
	}
	return "", false
}

func normalizeBadge(badge string, categoryName string) string {
	trimmed := strings.TrimSpace(badge)
	if trimmed == "" {
		return badgeFromName(categoryName)
	}
	runes := []rune(trimmed)
	if len(runes) > 2 {
		return string(runes[:2])
	}
	return string(runes)
}

func badgeFromName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "-"
	}

	runes := []rune(trimmed)
	if len(runes) == 0 {
		return "-"
	}
	return string(runes[0])
}

func weekdayLabel(dateText string) string {
	day, err := time.ParseInLocation("2006-01-02", dateText, time.Local)
	if err != nil {
		return ""
	}
	labels := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
	return labels[int(day.Weekday())]
}

func monthLabel(monthKey string) string {
	parts := strings.Split(monthKey, "-")
	if len(parts) != 2 {
		return monthKey
	}
	monthValue, err := strconv.Atoi(parts[1])
	if err != nil {
		return monthKey
	}
	return fmt.Sprintf("%s年%d月", parts[0], monthValue)
}

func formatMoney(value float64) string {
	rounded := round2(value)
	if math.Mod(rounded, 1) == 0 {
		return fmt.Sprintf("¥%.0f", rounded)
	}
	return fmt.Sprintf("¥%.2f", rounded)
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func currentUserID(c *gin.Context) (uint64, bool) {
	raw, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch value := raw.(type) {
	case uint64:
		return value, true
	case int64:
		if value < 0 {
			return 0, false
		}
		return uint64(value), true
	case int:
		if value < 0 {
			return 0, false
		}
		return uint64(value), true
	case float64:
		if value < 0 {
			return 0, false
		}
		return uint64(value), true
	case string:
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}
