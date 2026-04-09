package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"billsoftware/backend/internal/response"
	"billsoftware/backend/internal/service"
)

type AdminAuthHandler struct{}

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAdminAuthHandler() *AdminAuthHandler {
	return &AdminAuthHandler{}
}

func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid admin login payload")
		return
	}

	username := service.ResolveAdminUsername()
	passwordHash := service.ResolveAdminPasswordHash()
	requestHash := service.HashPassword(req.Password)

	if req.Username != username || requestHash != passwordHash {
		response.Fail(c, http.StatusUnauthorized, "admin username or password is incorrect")
		return
	}

	response.Success(c, gin.H{
		"token": service.BuildAdminAccessToken(username, passwordHash),
	})
}

