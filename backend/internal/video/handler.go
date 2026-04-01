package video

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
	api.GET("/feed/recommend", h.recommend)
	api.GET("/feed/hot", h.hot)
	api.GET("/feed/following", requiredAuth, h.following)
	api.GET("/areas/:id/videos", h.areaVideos)
	api.GET("/videos/:id", optionalAuth, h.detail)
	api.GET("/users/:id/videos", h.authorVideos)
	api.GET("/creator/videos", requiredAuth, h.creatorVideos)
	api.POST("/videos", requiredAuth, h.createVideo)
	api.PATCH("/videos/:id", requiredAuth, h.updateVideo)
	api.POST("/videos/:id/source", requiredAuth, h.uploadSource)
	api.POST("/videos/:id/cover", requiredAuth, h.uploadCover)
	api.POST("/videos/:id/likes", requiredAuth, h.like)
	api.DELETE("/videos/:id/likes", requiredAuth, h.unlike)
	api.GET("/videos/:id/likes/me", requiredAuth, h.likeStatus)
	api.POST("/videos/:id/favorites", requiredAuth, h.favorite)
	api.DELETE("/videos/:id/favorites", requiredAuth, h.unfavorite)
	api.GET("/videos/:id/favorites/me", requiredAuth, h.favoriteStatus)
}

func (h *Handler) recommend(c *gin.Context) {
	limit, err := parseLimit(c.Query("limit"))
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3101, "invalid recommend limit")
		return
	}

	result, err := h.service.ListRecommend(c.Request.Context(), c.Query("cursor"), limit)
	if err != nil {
		if errors.Is(err, ErrInvalidCursor) {
			response.Error(c, stdhttp.StatusBadRequest, 3102, "invalid recommend cursor")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3501, "list recommend feed failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) hot(c *gin.Context) {
	limit, err := parseLimit(c.Query("limit"))
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3105, "invalid hot limit")
		return
	}

	result, err := h.service.ListHot(c.Request.Context(), c.Query("cursor"), limit)
	if err != nil {
		if errors.Is(err, ErrInvalidCursor) {
			response.Error(c, stdhttp.StatusBadRequest, 3106, "invalid hot cursor")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3509, "list hot feed failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) following(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	limit, err := parseLimit(c.Query("limit"))
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3107, "invalid following limit")
		return
	}

	result, err := h.service.ListFollowing(c.Request.Context(), userID, c.Query("cursor"), limit)
	if err != nil {
		if errors.Is(err, ErrInvalidCursor) {
			response.Error(c, stdhttp.StatusBadRequest, 3108, "invalid following cursor")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3510, "list following feed failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) detail(c *gin.Context) {
	videoID, err := parseUintParam(c.Param("id"))
	if err != nil || videoID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 3103, "invalid video id")
		return
	}

	viewerID, _ := authctx.GetCurrentUserID(c)
	result, err := h.service.GetDetail(c.Request.Context(), videoID, viewerID)
	if err != nil {
		if errors.Is(err, ErrVideoNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3502, "get video detail failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) areaVideos(c *gin.Context) {
	areaID, err := parseUintParam(c.Param("id"))
	if err != nil || areaID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 3112, "invalid area id")
		return
	}

	limit, err := parseLimit(c.Query("limit"))
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3113, "invalid area videos limit")
		return
	}

	result, err := h.service.ListByArea(c.Request.Context(), areaID, c.Query("sort"), c.Query("cursor"), limit)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCursor):
			response.Error(c, stdhttp.StatusBadRequest, 3114, "invalid area videos cursor")
		case errors.Is(err, ErrInvalidSort):
			response.Error(c, stdhttp.StatusBadRequest, 3115, "invalid area videos sort")
		case errors.Is(err, ErrAreaNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3116, "area not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3512, "list area videos failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) authorVideos(c *gin.Context) {
	authorID, err := parseUintParam(c.Param("id"))
	if err != nil || authorID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 3109, "invalid author id")
		return
	}

	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3110, "invalid author videos pagination")
		return
	}

	result, err := h.service.ListByAuthor(c.Request.Context(), authorID, pagination.Page, pagination.PageSize)
	if err != nil {
		if errors.Is(err, ErrAuthorNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 3111, "author not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3511, "list author videos failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) like(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	if err := h.service.Like(c.Request.Context(), videoID, userID); err != nil {
		if errors.Is(err, ErrVideoNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3503, "like video failed")
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) unlike(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	if err := h.service.Unlike(c.Request.Context(), videoID, userID); err != nil {
		if errors.Is(err, ErrVideoNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3504, "unlike video failed")
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) likeStatus(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	result, err := h.service.LikeStatus(c.Request.Context(), videoID, userID)
	if err != nil {
		if errors.Is(err, ErrVideoNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3505, "get like status failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) favorite(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	if err := h.service.Favorite(c.Request.Context(), videoID, userID); err != nil {
		if errors.Is(err, ErrVideoNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3506, "favorite video failed")
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) unfavorite(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	if err := h.service.Unfavorite(c.Request.Context(), videoID, userID); err != nil {
		if errors.Is(err, ErrVideoNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3507, "unfavorite video failed")
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) favoriteStatus(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	result, err := h.service.FavoriteStatus(c.Request.Context(), videoID, userID)
	if err != nil {
		if errors.Is(err, ErrVideoNotFound) {
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
			return
		}
		response.Error(c, stdhttp.StatusInternalServerError, 3508, "get favorite status failed")
		return
	}

	response.Success(c, result)
}

func (h *Handler) creatorVideos(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	pagination, err := request.ParsePagination(c)
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3117, "invalid creator videos pagination")
		return
	}

	result, err := h.service.ListCreatorVideos(c.Request.Context(), userID, c.Query("review_status"), pagination.Page, pagination.PageSize)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidReviewStatus):
			response.Error(c, stdhttp.StatusBadRequest, 3118, "invalid creator videos review status")
		case errors.Is(err, ErrAuthorNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3111, "author not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3513, "list creator videos failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) createVideo(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	var req CreateVideoInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3119, "invalid create video payload")
		return
	}

	result, err := h.service.CreateVideo(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, stdhttp.StatusBadRequest, 3120, "invalid create video input")
		case errors.Is(err, ErrAreaNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3116, "area not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3514, "create video failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) updateVideo(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	var req UpdateVideoInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3125, "invalid update video payload")
		return
	}

	result, err := h.service.UpdateVideo(c.Request.Context(), videoID, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, stdhttp.StatusBadRequest, 3126, "invalid update video input")
		case errors.Is(err, ErrAreaNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3116, "area not found")
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3517, "update video failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) uploadSource(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3121, "video source file is required")
		return
	}

	result, err := h.service.UploadSource(c.Request.Context(), videoID, userID, fileHeader)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, stdhttp.StatusBadRequest, 3122, "invalid video source upload")
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3515, "upload video source failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) uploadCover(c *gin.Context) {
	videoID, ok := currentUserAndVideoID(c)
	if !ok {
		return
	}
	userID, _ := authctx.GetCurrentUserID(c)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 3123, "video cover file is required")
		return
	}

	result, err := h.service.UploadCover(c.Request.Context(), videoID, userID, fileHeader)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, stdhttp.StatusBadRequest, 3124, "invalid video cover upload")
		case errors.Is(err, ErrVideoNotFound):
			response.Error(c, stdhttp.StatusNotFound, 3104, "video not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 3516, "upload video cover failed")
		}
		return
	}

	response.Success(c, result)
}

func parseLimit(raw string) (int, error) {
	if raw == "" {
		return 0, nil
	}
	return strconv.Atoi(raw)
}

func parseUintParam(raw string) (uint64, error) {
	return strconv.ParseUint(raw, 10, 64)
}

func currentUserAndVideoID(c *gin.Context) (uint64, bool) {
	videoID, err := parseUintParam(c.Param("id"))
	if err != nil || videoID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 3103, "invalid video id")
		return 0, false
	}
	if _, ok := authctx.GetCurrentUserID(c); !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return 0, false
	}
	return videoID, true
}
