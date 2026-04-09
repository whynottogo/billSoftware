package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
)

type UserHandler struct {
	engine *xorm.Engine
}

func NewUserHandler(engine *xorm.Engine) *UserHandler {
	return &UserHandler{engine: engine}
}

func (h *UserHandler) List(c *gin.Context) {
	users := make([]model.User, 0)
	if err := h.engine.Desc("id").Find(&users); err != nil {
		response.Fail(c, http.StatusInternalServerError, "query users failed")
		return
	}

	response.Success(c, gin.H{
		"list": users,
	})
}

func (h *UserHandler) ChangeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var payload struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid status payload")
		return
	}

	if payload.Status != 0 && payload.Status != 1 {
		response.Fail(c, http.StatusBadRequest, "status must be 0 or 1")
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

	user := &model.User{}
	has, err := session.ID(id).Get(user)
	if err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "query user failed")
		return
	}
	if !has {
		_ = session.Rollback()
		response.Fail(c, http.StatusNotFound, "user not found")
		return
	}

	if _, err := session.ID(id).Cols("status").Update(&model.User{
		Status: payload.Status,
	}); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "update user status failed")
		return
	}

	if payload.Status == 0 {
		if _, err := session.Where("user_id = ? AND is_active = ?", id, 1).Cols("is_active").Update(&model.UserSession{
			IsActive: 0,
		}); err != nil {
			_ = session.Rollback()
			response.Fail(c, http.StatusInternalServerError, "invalidate user sessions failed")
			return
		}
	}

	if err := session.Commit(); err != nil {
		_ = session.Rollback()
		response.Fail(c, http.StatusInternalServerError, "commit transaction failed")
		return
	}

	response.Success(c, gin.H{
		"user_id": id,
		"status":  payload.Status,
	})
}

func (h *UserHandler) Summary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	type Summary struct {
		Income  string `json:"income"`
		Expense string `json:"expense"`
		Balance string `json:"balance"`
	}

	result := Summary{
		Income:  "0.00",
		Expense: "0.00",
		Balance: "0.00",
	}

	_, err = h.engine.SQL(`
		SELECT
			COALESCE(SUM(CASE WHEN record_type = 'income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN record_type = 'expense' THEN amount ELSE 0 END), 0) AS expense,
			COALESCE(SUM(CASE WHEN record_type = 'income' THEN amount ELSE -amount END), 0) AS balance
		FROM ledger_records
		WHERE user_id = ?
	`, id).Get(&result)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "query user summary failed")
		return
	}

	response.Success(c, result)
}
