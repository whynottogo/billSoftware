package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"

	"billsoftware/backend/internal/model"
	"billsoftware/backend/internal/response"
)

func UserAuthRequired(engine *xorm.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer"))
		if token == "" {
			response.Fail(c, 401, "user token is required")
			c.Abort()
			return
		}

		session := &model.UserSession{}
		has, err := engine.Where("session_token = ? AND is_active = ?", token, 1).Get(session)
		if err != nil || !has {
			response.Fail(c, 401, "user session is invalid")
			c.Abort()
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}

