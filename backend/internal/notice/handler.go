package notice

import (
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
	noticeGroup := api.Group("/notices", authMiddleware)
	{
		noticeGroup.GET("", h.list)
		noticeGroup.PATCH("/:id/read", h.read)
	}
}

func (h *Handler) list(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3601, "invalid notice pagination")
		return
	}

	result, err := h.service.List(c.Request.Context(), userID, pagination.Page, pagination.PageSize)
	if err != nil {
		response.Error(c, stdhttp.StatusInternalServerError, 3602, "list notices failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) read(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	noticeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || noticeID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 3603, "invalid notice id")
		return
	}

	result, err := h.service.MarkRead(c.Request.Context(), userID, noticeID)
	if err != nil {
		if err == ErrNoticeNotFound {
			response.Error(c, stdhttp.StatusNotFound, 3604, "notice not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3605, "mark notice read failed")
		return
	}

	response.Success(c, result)
}
