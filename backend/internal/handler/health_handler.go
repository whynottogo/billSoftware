package handler

import (
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/response"
)

type HealthHandler struct {
	engine *xorm.Engine
}

func NewHealthHandler(engine *xorm.Engine) *HealthHandler {
	return &HealthHandler{engine: engine}
}

func (h *HealthHandler) Ping(c *gin.Context) {
	err := h.engine.Ping()
	if err != nil {
		response.Fail(c, 500, "database ping failed")
		return
	}

	response.Success(c, gin.H{
		"service": "bill-software-backend",
		"status":  "ok",
	})
}

