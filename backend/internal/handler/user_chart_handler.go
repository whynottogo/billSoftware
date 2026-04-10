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

type UserChartHandler struct {
	core *budgetChartCore
}

func NewUserChartHandler(engine *xorm.Engine) *UserChartHandler {
	return &UserChartHandler{
		core: newBudgetChartCore(engine),
	}
}

func (h *UserChartHandler) ListYears(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	years, err := h.core.ListChartYears(userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query chart years failed")
		return
	}

	response.Success(c, gin.H{
		"years":       years,
		"currentYear": time.Now().Year(),
	})
}

func (h *UserChartHandler) GetExpenseYear(c *gin.Context) {
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

	data, err := h.core.GetExpenseChart(userID, year)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query expense chart failed")
		return
	}

	response.Success(c, data)
}

func (h *UserChartHandler) GetIncomeYear(c *gin.Context) {
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

	data, err := h.core.GetIncomeChart(userID, year)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query income chart failed")
		return
	}

	response.Success(c, data)
}
