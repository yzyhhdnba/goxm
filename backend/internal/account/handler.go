package account

import (
	"errors"
	stdhttp "net/http"
	"strconv"
	"strings"

	"pilipili-go/backend/pkg/authctx"
	"pilipili-go/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(api *gin.RouterGroup, authMiddleware gin.HandlerFunc, optionalAuth gin.HandlerFunc) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", h.register)
		authGroup.POST("/login", h.login)
		authGroup.POST("/refresh", h.refresh)
		authGroup.POST("/logout", authMiddleware, h.logout)
		authGroup.GET("/check-username", h.checkUsername)
		authGroup.GET("/check-email", h.checkEmail)
	}

	usersGroup := api.Group("/users")
	{
		usersGroup.GET("/me", authMiddleware, h.me)
		usersGroup.GET("/me/dashboard", authMiddleware, h.dashboard)
		usersGroup.GET("/:id/profile", optionalAuth, h.profile)
	}
}

func (h *Handler) register(c *gin.Context) {
	var req RegisterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 2001, "invalid register payload")
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUsernameTaken):
			response.Error(c, stdhttp.StatusConflict, 2002, "username already exists")
		case errors.Is(err, ErrEmailTaken):
			response.Error(c, stdhttp.StatusConflict, 2003, "email already exists")
		case errors.Is(err, ErrInvalidInput):
			response.Error(c, stdhttp.StatusBadRequest, 2004, "invalid register input")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 2500, "register failed")
		}
		return
	}

	c.JSON(stdhttp.StatusCreated, response.Envelope{
		Code:    0,
		Message: "ok",
		Data: gin.H{
			"user": user,
		},
	})
}

func (h *Handler) login(c *gin.Context) {
	var req LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 2005, "invalid login payload")
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredential):
			response.Error(c, stdhttp.StatusUnauthorized, 2006, "invalid credentials")
		case errors.Is(err, ErrInactiveUser):
			response.Error(c, stdhttp.StatusForbidden, 2007, "user is inactive")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 2501, "login failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) refresh(c *gin.Context) {
	var req RefreshInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, stdhttp.StatusBadRequest, 2008, "invalid refresh payload")
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidToken):
			response.Error(c, stdhttp.StatusUnauthorized, 2009, "invalid refresh token")
		case errors.Is(err, ErrInactiveUser):
			response.Error(c, stdhttp.StatusForbidden, 2010, "user is inactive")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 2502, "refresh failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) logout(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	if err := h.service.Logout(c.Request.Context(), userID); err != nil {
		response.Error(c, stdhttp.StatusInternalServerError, 2503, "logout failed")
		return
	}

	response.Success(c, gin.H{})
}

func (h *Handler) me(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	user, err := h.service.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, stdhttp.StatusInternalServerError, 2504, "get current user failed")
		return
	}

	response.Success(c, user)
}

func (h *Handler) profile(c *gin.Context) {
	profileUserID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || profileUserID == 0 {
		response.Error(c, stdhttp.StatusBadRequest, 2014, "invalid user id")
		return
	}

	viewerID, _ := authctx.GetCurrentUserID(c)
	result, err := h.service.GetProfile(c.Request.Context(), profileUserID, viewerID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			response.Error(c, stdhttp.StatusNotFound, 2015, "user not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 2507, "get profile failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) dashboard(c *gin.Context) {
	userID, ok := authctx.GetCurrentUserID(c)
	if !ok {
		response.Error(c, stdhttp.StatusUnauthorized, 2011, "unauthorized")
		return
	}

	result, err := h.service.GetDashboard(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			response.Error(c, stdhttp.StatusNotFound, 2015, "user not found")
		default:
			response.Error(c, stdhttp.StatusInternalServerError, 2508, "get dashboard failed")
		}
		return
	}

	response.Success(c, result)
}

func (h *Handler) checkUsername(c *gin.Context) {
	username := strings.TrimSpace(c.Query("username"))
	if username == "" {
		response.Error(c, stdhttp.StatusBadRequest, 2012, "username is required")
		return
	}

	result, err := h.service.CheckUsername(c.Request.Context(), username)
	if err != nil {
		response.Error(c, stdhttp.StatusInternalServerError, 2505, "check username failed")
		return
	}
	response.Success(c, result)
}

func (h *Handler) checkEmail(c *gin.Context) {
	email := strings.TrimSpace(c.Query("email"))
	if email == "" {
		response.Error(c, stdhttp.StatusBadRequest, 2013, "email is required")
		return
	}

	result, err := h.service.CheckEmail(c.Request.Context(), email)
	if err != nil {
		response.Error(c, stdhttp.StatusInternalServerError, 2506, "check email failed")
		return
	}
	response.Success(c, result)
}
