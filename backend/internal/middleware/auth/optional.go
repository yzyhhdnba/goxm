package auth

import (
	"strings"

	"pilipili-go/backend/internal/account"
	appauth "pilipili-go/backend/internal/auth"
	"pilipili-go/backend/pkg/authctx"

	"github.com/gin-gonic/gin"
)

func Optional(tokenManager *appauth.TokenManager, repo *account.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := tokenManager.ParseAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		user, err := repo.FindByID(c.Request.Context(), claims.UserID)
		if err != nil || user.Status != account.StatusActive || user.TokenVersion != claims.TokenVersion {
			c.Next()
			return
		}

		authctx.SetCurrentUser(c, authctx.CurrentUser{
			ID:        user.ID,
			Username:  user.Username,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
		})
		c.Next()
	}
}
