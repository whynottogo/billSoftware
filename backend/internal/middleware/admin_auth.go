package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"billsoftware/backend/internal/response"
	"billsoftware/backend/internal/service"
)

func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer"))
		expectedToken := service.BuildAdminAccessToken(
			service.ResolveAdminUsername(),
			service.ResolveAdminPasswordHash(),
		)

		if token == "" || token != expectedToken {
			response.Fail(c, 401, "admin token is invalid")
			c.Abort()
			return
		}

		c.Next()
	}
}

