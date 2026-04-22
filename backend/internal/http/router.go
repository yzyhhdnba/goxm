package http

import (
	"context"
	stdhttp "net/http"
	"strings"
	"time"

	"pilipili-go/backend/internal/account"
	"pilipili-go/backend/internal/admin"
	"pilipili-go/backend/internal/area"
	appauth "pilipili-go/backend/internal/auth"
	"pilipili-go/backend/internal/comment"
	appconfig "pilipili-go/backend/internal/config"
	"pilipili-go/backend/internal/history"
	"pilipili-go/backend/internal/media"
	authmw "pilipili-go/backend/internal/middleware/auth"
	appmq "pilipili-go/backend/internal/mq/rocketmq"
	"pilipili-go/backend/internal/notice"
	appredis "pilipili-go/backend/internal/redis"
	"pilipili-go/backend/internal/search"
	"pilipili-go/backend/internal/social"
	"pilipili-go/backend/internal/video"
	"pilipili-go/backend/pkg/response"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter 把依赖组装、中间件分层和路由注册收口到一个入口。
// 对应文档《key-code-implementation.md》中的“路由装配与中间件分层”。
func NewRouter(cfg *appconfig.Config, gormDB *gorm.DB, redisClient *appredis.Client, mqClient *appmq.Client) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	mediaRootDir := cfg.Media.RootDir
	if strings.TrimSpace(mediaRootDir) == "" {
		mediaRootDir = "storage"
	}
	mediaPublicBaseURL := cfg.Media.PublicBaseURL
	if strings.TrimSpace(mediaPublicBaseURL) == "" {
		mediaPublicBaseURL = "/uploads"
	}

	accountRepo := account.NewRepository(gormDB)
	socialRepo := social.NewRepository(gormDB)
	tokenManager := appauth.NewTokenManager(cfg.JWT)
	mediaStorage := media.NewLocalStorage(mediaRootDir, mediaPublicBaseURL)
	videoCache := video.NewRedisCache(redisClient)
	accountService := account.NewService(accountRepo, tokenManager)
	accountHandler := account.NewHandler(accountService)
	areaHandler := area.NewHandler(area.NewService(area.NewRepository(gormDB)))
	videoHandler := video.NewHandler(video.NewService(video.NewRepository(gormDB), socialRepo, mediaStorage, videoCache))
	commentHandler := comment.NewHandler(comment.NewService(comment.NewRepository(gormDB), videoCache))
	socialHandler := social.NewHandler(social.NewService(socialRepo))
	searchHandler := search.NewHandler(search.NewService(search.NewRepository(gormDB)))
	historyHandler := history.NewHandler(history.NewService(history.NewRepository(gormDB)))
	noticeHandler := notice.NewHandler(notice.NewService(notice.NewRepository(gormDB)))
	adminHandler := admin.NewHandler(admin.NewService(admin.NewRepository(gormDB), accountRepo, videoCache))

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}))
	router.Static(mediaPublicBaseURL, mediaRootDir)

	router.GET("/healthz", func(c *gin.Context) {
		sqlDB, err := gormDB.DB()
		if err != nil {
			response.Error(c, stdhttp.StatusServiceUnavailable, 1001, "database handle unavailable")
			return
		}

		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			response.Error(c, stdhttp.StatusServiceUnavailable, 1002, "database unavailable")
			return
		}

		redisStatus := "disabled"
		if redisClient != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
			defer cancel()

			if err := redisClient.Ping(ctx); err != nil {
				response.Error(c, stdhttp.StatusServiceUnavailable, 1003, "redis unavailable")
				return
			}
			redisStatus = "ok"
		}

		rocketMQStatus := "disabled"
		if mqClient != nil && mqClient.Enabled() {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
			defer cancel()

			if err := mqClient.Ping(ctx); err != nil {
				response.Error(c, stdhttp.StatusServiceUnavailable, 1004, "rocketmq unavailable")
				return
			}
			rocketMQStatus = "ok"
		}

		response.Success(c, gin.H{
			"status":    "ok",
			"service":   "pilipili-go-api",
			"timestamp": time.Now().Format(time.RFC3339),
			"dependencies": gin.H{
				"mysql":    "ok",
				"redis":    redisStatus,
				"rocketmq": rocketMQStatus,
			},
		})
	})

	api := router.Group("/api/v1")
	api.GET("/ping", func(c *gin.Context) {
		response.Success(c, gin.H{
			"message":   "pong",
			"service":   "pilipili-go-api",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	requiredAuth := authmw.Require(tokenManager, accountRepo)
	optionalAuth := authmw.Optional(tokenManager, accountRepo)
	accountHandler.RegisterRoutes(api, requiredAuth, optionalAuth)
	areaHandler.RegisterRoutes(api)
	videoHandler.RegisterRoutes(api, optionalAuth, requiredAuth)
	commentHandler.RegisterRoutes(api, optionalAuth, requiredAuth)
	socialHandler.RegisterRoutes(api, requiredAuth)
	searchHandler.RegisterRoutes(api)
	historyHandler.RegisterRoutes(api, requiredAuth)
	noticeHandler.RegisterRoutes(api, requiredAuth)
	adminHandler.RegisterRoutes(api, requiredAuth)

	return router
}
