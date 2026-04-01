package history

import (
	"errors"
	stdhttp "net/http"

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
	historyGroup := api.Group("/histories")
	historyGroup.Use(authMiddleware)
	historyGroup.GET("", h.list)
	historyGroup.POST("", h.report)
}

func (h *Handler) list(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3701, "invalid histories pagination")
		return
	}

	result, err := h.service.List(c.Request.Context(), userID, pagination.Page, pagination.PageSize)
	if err != nil {
		response.Error(c, stdhttp.StatusInternalServerError, 3903, "list histories failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) report(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	var req ReportInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3702, "invalid history payload")
		return
	}

	if err := h.service.Report(c.Request.Context(), userID, req); err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, stdhttp.StatusBadRequest, 3703, "invalid history input")
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3704, "video not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3904, "report history failed")
		}
		return
	}

	response.Success(c, gin.H{})
}
