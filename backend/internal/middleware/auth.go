package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"wrappedweekly/backend/internal/usecase"
	"wrappedweekly/backend/pkg/response"
)

const CookieName = "wrapped_weekly_token"
const ContextUserIDKey = "userID"

func RequireAuth(jwtManager *usecase.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.Error(c, 401, "autentikasi diperlukan")
			c.Abort()
			return
		}

		claims, err := jwtManager.Verify(token)
		if err != nil || claims == nil {
			response.Error(c, 401, "token tidak valid atau kedaluwarsa")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if cookie, err := c.Cookie(CookieName); err == nil && cookie != "" {
		return cookie
	}
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

func UserIDFromContext(c *gin.Context) string {
	v, exists := c.Get(ContextUserIDKey)
	if !exists {
		return ""
	}
	id, _ := v.(string)
	return id
}
