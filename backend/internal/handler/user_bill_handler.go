package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/response"
)

var errBillUserNotFound = errors.New("bill user not found")

type UserBillHandler struct {
	core *billStatsCore
}

type AdminUserBillHandler struct {
	core *billStatsCore
}

func NewUserBillHandler(engine *xorm.Engine) *UserBillHandler {
	return &UserBillHandler{
		core: newBillStatsCore(engine),
	}
}

func NewAdminUserBillHandler(engine *xorm.Engine) *AdminUserBillHandler {
	return &AdminUserBillHandler{
		core: newBillStatsCore(engine),
	}
}

func (h *UserBillHandler) ListYears(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	years, history, err := h.core.ListBillYears(userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query bill years failed")
		return
	}

	response.Success(c, gin.H{
		"years":   years,
		"history": history,
	})
}

func (h *UserBillHandler) GetYearDetail(c *gin.Context) {
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

	data, err := h.core.GetBillYearDetail(userID, year)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query bill year detail failed")
		return
	}

	response.Success(c, data)
}

func (h *UserBillHandler) GetMonthDetail(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "user context is invalid")
		return
	}

	monthKey := strings.TrimSpace(c.Param("month"))
	if _, _, _, err := resolveMonthRange(monthKey); err != nil {
		response.Fail(c, http.StatusBadRequest, "month must be YYYY-MM")
		return
	}

	data, err := h.core.GetBillMonthDetail(userID, monthKey)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query bill month detail failed")
		return
	}

	response.Success(c, data)
}

func (h *AdminUserBillHandler) GetOverview(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	data, err := h.core.GetAdminUserBillOverview(userID)
	if err != nil {
		if errors.Is(err, errBillUserNotFound) {
			response.Fail(c, http.StatusNotFound, "user not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "query admin user bill overview failed")
		return
	}

	response.Success(c, data)
}
