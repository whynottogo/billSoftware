package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
)

type AdminDashboardHandler struct {
	engine *xorm.Engine
}

type AdminDashboardTrendItem struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type AdminDashboardRecentUser struct {
	ID             uint64    `json:"id"`
	Username       string    `json:"username"`
	Nickname       string    `json:"nickname"`
	Email          string    `json:"email"`
	Status         int       `json:"status"`
	ApprovalStatus string    `json:"approval_status"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewAdminDashboardHandler(engine *xorm.Engine) *AdminDashboardHandler {
	return &AdminDashboardHandler{engine: engine}
}

func (h *AdminDashboardHandler) Overview(c *gin.Context) {
	now := time.Now()
	today := now.Format("2006-01-02")

	totalUsers, err := h.engine.Count(new(model.User))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "count users failed")
		return
	}

	pendingUsers, err := h.engine.Where("approval_status = ?", model.UserApprovalPending).Count(new(model.User))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "count pending users failed")
		return
	}

	todayBills, err := h.engine.Where("record_date = ?", today).Count(new(model.LedgerRecord))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "count today bills failed")
		return
	}

	families, err := h.engine.Count(new(model.Family))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "count families failed")
		return
	}

	userGrowth, err := h.queryUserGrowth(now)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query user growth failed")
		return
	}

	billTrend, err := h.queryBillTrend(now)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query bill trend failed")
		return
	}

	recentUserRows := make([]model.User, 0)
	if err := h.engine.Desc("created_at").Limit(5).Find(&recentUserRows); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query recent users failed")
		return
	}
	recentUsers := make([]AdminDashboardRecentUser, 0, len(recentUserRows))
	for _, user := range recentUserRows {
		recentUsers = append(recentUsers, AdminDashboardRecentUser{
			ID:             user.ID,
			Username:       user.Username,
			Nickname:       user.Nickname,
			Email:          user.Email,
			Status:         user.Status,
			ApprovalStatus: user.ApprovalStatus,
			CreatedAt:      user.CreatedAt,
		})
	}

	response.Success(c, gin.H{
		"total_users":   totalUsers,
		"pending_users": pendingUsers,
		"today_bills":   todayBills,
		"families":      families,
		"user_growth":   userGrowth,
		"bill_trend":    billTrend,
		"recent_users":  recentUsers,
	})
}

func (h *AdminDashboardHandler) queryUserGrowth(now time.Time) ([]AdminDashboardTrendItem, error) {
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -5, 0)
	monthCounts := make(map[string]int64)

	type row struct {
		MonthLabel string `xorm:"month_label"`
		Total      int64  `xorm:"total"`
	}

	rows := make([]row, 0)
	if err := h.engine.SQL(`
		SELECT DATE_FORMAT(created_at, '%Y-%m') AS month_label, COUNT(*) AS total
		FROM users
		WHERE created_at >= ?
		GROUP BY DATE_FORMAT(created_at, '%Y-%m')
	`, monthStart).Find(&rows); err != nil {
		return nil, err
	}

	for _, item := range rows {
		monthCounts[item.MonthLabel] = item.Total
	}

	trend := make([]AdminDashboardTrendItem, 0, 6)
	for index := 0; index < 6; index++ {
		month := monthStart.AddDate(0, index, 0)
		key := month.Format("2006-01")
		trend = append(trend, AdminDashboardTrendItem{
			Label: month.Format("01月"),
			Value: monthCounts[key],
		})
	}

	return trend, nil
}

func (h *AdminDashboardHandler) queryBillTrend(now time.Time) ([]AdminDashboardTrendItem, error) {
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)
	dayCounts := make(map[string]int64)

	type row struct {
		DayLabel string `xorm:"day_label"`
		Total    int64  `xorm:"total"`
	}

	rows := make([]row, 0)
	if err := h.engine.SQL(`
		SELECT DATE_FORMAT(record_date, '%Y-%m-%d') AS day_label, COUNT(*) AS total
		FROM ledger_records
		WHERE record_date >= ?
		GROUP BY DATE_FORMAT(record_date, '%Y-%m-%d')
	`, startDate.Format("2006-01-02")).Find(&rows); err != nil {
		return nil, err
	}

	for _, item := range rows {
		dayCounts[item.DayLabel] = item.Total
	}

	trend := make([]AdminDashboardTrendItem, 0, 7)
	for index := 0; index < 7; index++ {
		day := startDate.AddDate(0, 0, index)
		key := day.Format("2006-01-02")
		trend = append(trend, AdminDashboardTrendItem{
			Label: day.Format("01-02"),
			Value: dayCounts[key],
		})
	}

	return trend, nil
}
