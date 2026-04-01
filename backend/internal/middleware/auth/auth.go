package auth

import (
	"strings"

	"pilipili-go/backend/internal/account"
	appauth "pilipili-go/backend/internal/auth"
	"pilipili-go/backend/pkg/authctx"
	"pilipili-go/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func Require(tokenManager *appauth.TokenManager, repo *account.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, 401, 4001, "missing bearer token")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := tokenManager.ParseAccessToken(tokenString)
		if err != nil {
			response.Error(c, 401, 4002, "invalid access token")
			c.Abort()
			return
		}

		user, err := repo.FindByID(c.Request.Context(), claims.UserID)
		if err != nil || user.Status != account.StatusActive || user.TokenVersion != claims.TokenVersion {
			response.Error(c, 401, 4003, "token is expired or revoked")
			c.Abort()
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
