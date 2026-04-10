package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/response"
)

type UserBudgetHandler struct {
	core *budgetChartCore
}

type budgetCategoryInput struct {
	CategoryID uint64  `json:"category_id"`
	Amount     float64 `json:"amount"`
}

type updateBudgetRequest struct {
	TotalAmount *float64               `json:"total_amount"`
	Categories  *[]budgetCategoryInput `json:"categories"`
	Items       *[]budgetCategoryInput `json:"items"`
}

func NewUserBudgetHandler(engine *xorm.Engine) *UserBudgetHandler {
	return &UserBudgetHandler{
		core: newBudgetChartCore(engine),
	}
}

func (h *UserBudgetHandler) GetCurrentMonth(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	data, err := h.core.GetMonthBudget(userID, time.Now())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query current month budget failed")
		return
	}

	response.Success(c, data)
}

func (h *UserBudgetHandler) UpdateCurrentMonth(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	var req updateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid month budget payload")
		return
	}
	if req.TotalAmount == nil && req.Categories == nil && req.Items == nil {
		response.Fail(c, http.StatusBadRequest, "at least one budget field is required")
		return
	}

	if req.TotalAmount != nil && *req.TotalAmount < 0 {
		response.Fail(c, http.StatusBadRequest, "total_amount must be greater than or equal to 0")
		return
	}

	items := req.Categories
	if items == nil {
		items = req.Items
	}

	if err := h.core.UpsertCurrentMonthBudget(userID, req.TotalAmount, items); err != nil {
		switch err {
		case errBudgetInvalidCategory, errBudgetInvalidAmount:
			response.Fail(c, http.StatusBadRequest, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, "update current month budget failed")
		}
		return
	}

	data, err := h.core.GetMonthBudget(userID, time.Now())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query current month budget failed")
		return
	}

	response.Success(c, data)
}

func (h *UserBudgetHandler) ListYearOptions(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	years, err := h.core.ListBudgetYears(userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query budget year options failed")
		return
	}

	response.Success(c, gin.H{
		"years":       years,
		"currentYear": time.Now().Year(),
	})
}

func (h *UserBudgetHandler) GetYear(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	year, err := strconv.Atoi(strings.TrimSpace(c.Param("year")))
	if err != nil || year <= 0 {
		response.Fail(c, http.StatusBadRequest, "year must be YYYY")
		return
	}

	data, err := h.core.GetYearBudget(userID, year)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query budget year detail failed")
		return
	}

	response.Success(c, data)
}

func (h *UserBudgetHandler) UpdateCurrentYear(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	var req updateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid year budget payload")
		return
	}
	if req.TotalAmount == nil && req.Categories == nil && req.Items == nil {
		response.Fail(c, http.StatusBadRequest, "at least one budget field is required")
		return
	}

	if req.TotalAmount != nil && *req.TotalAmount < 0 {
		response.Fail(c, http.StatusBadRequest, "total_amount must be greater than or equal to 0")
		return
	}

	items := req.Categories
	if items == nil {
		items = req.Items
	}

	if err := h.core.UpsertCurrentYearBudget(userID, req.TotalAmount, items); err != nil {
		switch err {
		case errBudgetInvalidCategory, errBudgetInvalidAmount:
			response.Fail(c, http.StatusBadRequest, err.Error())
		default:
			response.Fail(c, http.StatusInternalServerError, "update current year budget failed")
		}
		return
	}

	data, err := h.core.GetYearBudget(userID, time.Now().Year())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query current year budget failed")
		return
	}

	response.Success(c, data)
}
