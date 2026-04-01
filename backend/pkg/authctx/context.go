package authctx

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	currentUserIDKey = "current_user_id"
	currentUserKey   = "current_user"
)

type currentUserContextKey struct{}

type CurrentUser struct {
	ID        uint64
	Username  string
	Role      string
	AvatarURL string
}

func SetCurrentUser(c *gin.Context, user CurrentUser) {
	c.Set(currentUserIDKey, user.ID)
	c.Set(currentUserKey, user)
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), currentUserContextKey{}, user))
}

func SetCurrentUserID(c *gin.Context, userID uint64) {
	c.Set(currentUserIDKey, userID)
}

func GetCurrentUserID(c *gin.Context) (uint64, bool) {
	value, exists := c.Get(currentUserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint64)
	return userID, ok
}

func GetCurrentUser(c *gin.Context) (CurrentUser, bool) {
	value, exists := c.Get(currentUserKey)
	if !exists {
		return CurrentUser{}, false
	}

	user, ok := value.(CurrentUser)
	return user, ok
}

func GetCurrentUserFromContext(ctx context.Context) (CurrentUser, bool) {
	user, ok := ctx.Value(currentUserContextKey{}).(CurrentUser)
	return user, ok
}
