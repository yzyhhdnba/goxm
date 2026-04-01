package comment

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

func (h *Handler) RegisterRoutes(api *gin.RouterGroup, optionalAuth gin.HandlerFunc, requiredAuth gin.HandlerFunc) {
	api.GET("/videos/:id/comments", optionalAuth, h.listComments)
	api.POST("/videos/:id/comments", requiredAuth, h.createComment)
	api.GET("/comments/:id/replies", optionalAuth, h.listReplies)
	api.POST("/comments/:id/replies", requiredAuth, h.createReply)
	api.POST("/comments/:id/likes", requiredAuth, h.like)
	api.DELETE("/comments/:id/likes", requiredAuth, h.unlike)
	api.GET("/comments/:id/likes/me", requiredAuth, h.likeStatus)
}

func (h *Handler) listComments(c *gin.Context) {
	videoID, ok := parseTargetID(c, 4100, "invalid video id")
	if !ok {
		return
	}
	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 4101, "invalid comment pagination")
		return
	}
	viewerID, _ := authctx.GetCurrentUserID(c)

	result, err := h.service.ListComments(c.Request.Context(), videoID, pagination.Page, pagination.PageSize, viewerID)
	if err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, stdhttp.StatusNotFound, 4102, "video not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 4501, "list comments failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) createComment(c *gin.Context) {
	videoID, ok := parseTargetID(c, 4100, "invalid video id")
	if !ok {
		return
	}
	userID, exists := authctx.GetCurrentUserID(c)
	if !exists {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	var req CreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 4103, "invalid comment payload")
		return
	}

	result, err := h.service.CreateComment(c.Request.Context(), videoID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, stdhttp.StatusNotFound, 4102, "video not found")
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, stdhttp.StatusBadRequest, 4104, "invalid comment input")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 4502, "create comment failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) listReplies(c *gin.Context) {
	commentID, ok := parseTargetID(c, 4105, "invalid comment id")
	if !ok {
		return
	}
	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 4106, "invalid reply pagination")
		return
	}
	viewerID, _ := authctx.GetCurrentUserID(c)

	result, err := h.service.ListReplies(c.Request.Context(), commentID, pagination.Page, pagination.PageSize, viewerID)
	if err != nil {
		switch {
		case errors.Is(err, ErrCommentNotFound):
			response.Error(c, stdhttp.StatusNotFound, 4107, "comment not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 4503, "list replies failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) createReply(c *gin.Context) {
	commentID, ok := parseTargetID(c, 4105, "invalid comment id")
	if !ok {
		return
	}
	userID, exists := authctx.GetCurrentUserID(c)
	if !exists {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	var req CreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 4108, "invalid reply payload")
		return
	}

	result, err := h.service.CreateReply(c.Request.Context(), commentID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrCommentNotFound):
			response.Error(c, stdhttp.StatusNotFound, 4107, "comment not found")
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, stdhttp.StatusBadRequest, 4109, "invalid reply input")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 4504, "create reply failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) like(c *gin.Context) {
	commentID, userID, ok := currentUserAndCommentID(c)
	if !ok {
		return
	}

	if err := h.service.Like(c.Request.Context(), commentID, userID); err != nil {
		if errors.Is(err, ErrCommentNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 4107, "comment not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 4505, "like comment failed")
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) unlike(c *gin.Context) {
	commentID, userID, ok := currentUserAndCommentID(c)
	if !ok {
		return
	}

	if err := h.service.Unlike(c.Request.Context(), commentID, userID); err != nil {
		if errors.Is(err, ErrCommentNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 4107, "comment not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 4506, "unlike comment failed")
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) likeStatus(c *gin.Context) {
	commentID, userID, ok := currentUserAndCommentID(c)
	if !ok {
		return
	}

	result, err := h.service.LikeStatus(c.Request.Context(), commentID, userID)
	if err != nil {
		response.Error(c, stdhttp.StatusInternalServerError, 4507, "get comment like status failed")
		return
	}

	response.Success(c, result)
}

func currentUserAndCommentID(c *gin.Context) (uint64, uint64, bool) {
	commentID, ok := parseTargetID(c, 4105, "invalid comment id")
	if !ok {
		return 0, 0, false
	}
	userID, exists := authctx.GetCurrentUserID(c)
	if !exists {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return 0, 0, false
	}
	return commentID, userID, true
}

func parseTargetID(c *gin.Context, code int, message string) (uint64, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || value == 0 {
		response.Error(c, stdhttp.StatusBadRequest, code, message)
		return 0, false
	}
	return value, true
}
