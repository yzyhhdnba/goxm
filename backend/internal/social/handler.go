package social

import (
	"errors"
	stdhttp "net/http"
	"strconv"

	"pilipili-go/backend/pkg/authctx"
	"pilipili-go/backend/pkg/request"
	"pilipili-go/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(api *gin.RouterGroup, requiredAuth gin.HandlerFunc) {
	api.POST("/users/:id/follow", requiredAuth, h.follow)
	api.DELETE("/users/:id/follow", requiredAuth, h.unfollow)
	api.GET("/users/:id/follow-status", requiredAuth, h.followStatus)
	api.GET("/users/:id/followers", h.followers)
	api.GET("/users/:id/following", h.following)
}

func (h *Handler) follow(c *gin.Context) {
	targetID, userID, ok := currentAndTargetUserID(c)
	if !ok {
		return
	}

	if err := h.service.Follow(c.Request.Context(), userID, targetID); err != nil {
		switch {
		case errors.Is(err, ErrCannotFollowSelf):
			response.Error(c, stdhttp.StatusBadRequest, 5201, "cannot follow yourself")
		case errors.Is(err, ErrFollowTargetNotFound):
			response.Error(c, stdhttp.StatusNotFound, 5202, "follow target not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 5501, "follow user failed")
		}
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) unfollow(c *gin.Context) {
	targetID, userID, ok := currentAndTargetUserID(c)
	if !ok {
		return
	}

	if err := h.service.Unfollow(c.Request.Context(), userID, targetID); err != nil {
		switch {
		case errors.Is(err, ErrCannotFollowSelf):
			response.Error(c, stdhttp.StatusBadRequest, 5201, "cannot follow yourself")
		case errors.Is(err, ErrFollowTargetNotFound):
			response.Error(c, stdhttp.StatusNotFound, 5202, "follow target not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 5502, "unfollow user failed")
		}
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) followStatus(c *gin.Context) {
	targetID, userID, ok := currentAndTargetUserID(c)
	if !ok {
		return
	}

	result, err := h.service.Status(c.Request.Context(), userID, targetID)
	if err != nil {
		if errors.Is(err, ErrFollowTargetNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 5202, "follow target not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 5503, "get follow status failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) followers(c *gin.Context) {
	targetID, ok := targetUserID(c)
	if !ok {
		return
	}

	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 5203, "invalid followers pagination")
		return
	}

	result, err := h.service.ListFollowers(c.Request.Context(), targetID, pagination.Page, pagination.PageSize)
	if err != nil {
		if errors.Is(err, ErrFollowTargetNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 5202, "follow target not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 5504, "list followers failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) following(c *gin.Context) {
	targetID, ok := targetUserID(c)
	if !ok {
		return
	}

	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 5204, "invalid following pagination")
		return
	}

	result, err := h.service.ListFollowing(c.Request.Context(), targetID, pagination.Page, pagination.PageSize)
	if err != nil {
		if errors.Is(err, ErrFollowTargetNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 5202, "follow target not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 5505, "list following failed")
		return
	}

	response.Success(c, result)
}

func currentAndTargetUserID(c *gin.Context) (uint64, uint64, bool) {
	targetID, ok := targetUserID(c)
	if !ok {
		return 0, 0, false
	}
	userID, exists := authctx.GetCurrentUserID(c)
	if !exists {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return 0, 0, false
	}
	return targetID, userID, true
}

func targetUserID(c *gin.Context) (uint64, bool) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || targetID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 5200, "invalid user id")
		return 0, false
	}
	return targetID, true
}
