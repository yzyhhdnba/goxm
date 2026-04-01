package admin

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

func (h *Handler) RegisterRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	adminGroup := api.Group("/admin")
	adminGroup.Use(authMiddleware)
	adminGroup.GET("/videos", h.listVideos)
	adminGroup.GET("/videos/pending", h.listPendingVideos)
	adminGroup.POST("/videos/:id/approve", h.approve)
	adminGroup.POST("/videos/:id/reject", h.reject)
	adminGroup.GET("/stats/today", h.todayStats)
	adminGroup.GET("/stats/area", h.areaStats)
}

func (h *Handler) listVideos(c *gin.Context) {
	h.listWithStatus(c, c.Query("review_status"))
}

func (h *Handler) listPendingVideos(c *gin.Context) {
	h.listWithStatus(c, ReviewStatusPending)
}

func (h *Handler) listWithStatus(c *gin.Context, reviewStatus string) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3801, "invalid admin videos pagination")
		return
	}

	result, err := h.service.ListVideos(c.Request.Context(), userID, reviewStatus, pagination.Page, pagination.PageSize)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			response.Error(c, stdhttp.StatusForbidden, 3802, "forbidden")
		case errors.Is(err, ErrInvalidStatus):
			response.Error(c, stdhttp.StatusBadRequest, 3803, "invalid review status")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3905, "list admin videos failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) approve(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || videoID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 3804, "invalid video id")
		return
	}

	result, err := h.service.Approve(c.Request.Context(), userID, videoID)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			response.Error(c, stdhttp.StatusForbidden, 3802, "forbidden")
		case errors.Is(err, ErrInvalidPayload):
			response.Error(c, stdhttp.StatusBadRequest, 3805, "invalid approve input")
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3806, "video not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3906, "approve video failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) reject(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	videoID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || videoID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 3804, "invalid video id")
		return
	}

	var req ReviewInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3807, "invalid reject payload")
		return
	}

	result, err := h.service.Reject(c.Request.Context(), userID, videoID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			response.Error(c, stdhttp.StatusForbidden, 3802, "forbidden")
		case errors.Is(err, ErrInvalidPayload):
			response.Error(c, stdhttp.StatusBadRequest, 3808, "invalid reject input")
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3806, "video not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3907, "reject video failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) todayStats(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	result, err := h.service.GetTodayStats(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Error(c, stdhttp.StatusForbidden, 3802, "forbidden")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3908, "get today stats failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) areaStats(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	result, err := h.service.GetAreaStats(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Error(c, stdhttp.StatusForbidden, 3802, "forbidden")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3909, "get area stats failed")
		return
	}

	response.Success(c, result)
}
