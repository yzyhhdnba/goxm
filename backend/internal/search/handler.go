package search

import (
	"errors"
	stdhttp "net/http"

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

func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	searchGroup := api.Group("/search")
	searchGroup.GET("/videos", h.searchVideos)
	searchGroup.GET("/users", h.searchUsers)
}

func (h *Handler) searchVideos(c *gin.Context) {
	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3601, "invalid search videos pagination")
		return
	}

	result, err := h.service.SearchVideos(c.Request.Context(), c.Query("keyword"), pagination.Page, pagination.PageSize)
	if err != nil {
		if errors.Is(err, ErrInvalidKeyword) {
			response.Error(c, stdhttp.StatusBadRequest, 3602, "keyword is required")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3901, "search videos failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) searchUsers(c *gin.Context) {
	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3603, "invalid search users pagination")
		return
	}

	result, err := h.service.SearchUsers(c.Request.Context(), c.Query("keyword"), pagination.Page, pagination.PageSize)
	if err != nil {
		if errors.Is(err, ErrInvalidKeyword) {
			response.Error(c, stdhttp.StatusBadRequest, 3604, "keyword is required")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3902, "search users failed")
		return
	}

	response.Success(c, result)
}
