package area

import (
	stdhttp "net/http"

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
	api.GET("/areas", h.list)
}

func (h *Handler) list(c *gin.Context) {
	areas, err := h.service.List(c.Request.Context())
	if err != nil {
		response.Error(c, stdhttp.StatusInternalServerError, 3001, "list areas failed")
		return
	}

	response.Success(c, areas)
}
