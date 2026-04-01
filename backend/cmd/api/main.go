package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"pilipili-go/backend/internal/account"
	"pilipili-go/backend/internal/admin"
	"pilipili-go/backend/internal/area"
	"pilipili-go/backend/internal/comment"
	appconfig "pilipili-go/backend/internal/config"
	"pilipili-go/backend/internal/db"
	"pilipili-go/backend/internal/history"
	apphttp "pilipili-go/backend/internal/http"
	"pilipili-go/backend/internal/notice"
	appredis "pilipili-go/backend/internal/redis"
	"pilipili-go/backend/internal/social"
	"pilipili-go/backend/internal/video"
)

// main 是后端启动总入口，对应文档《KEY_CODE_IMPLEMENTATION.md》中的“程序入口与生命周期管理”。
// 走读时建议从这里顺着“配置 -> 数据库/Redis -> 路由 -> 优雅停机”这条主链路往下看。
func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := appconfig.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gormDB, err := db.New(cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("get sql db: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("close sql db: %v", err)
		}
	}()

	redisClient, err := appredis.New(cfg.Redis)
	if err != nil {
		log.Fatalf("connect redis: %v", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("close redis: %v", err)
		}
	}()

	if err := account.NewRepository(gormDB).AutoMigrate(); err != nil {
		log.Fatalf("auto migrate users: %v", err)
	}
	areaRepo := area.NewRepository(gormDB)
	if err := areaRepo.AutoMigrate(); err != nil {
		log.Fatalf("auto migrate areas: %v", err)
	}
	videoRepo := video.NewRepository(gormDB)
	if err := videoRepo.AutoMigrate(); err != nil {
		log.Fatalf("auto migrate videos: %v", err)
	}
	if err := comment.NewRepository(gormDB).AutoMigrate(); err != nil {
		log.Fatalf("auto migrate comments: %v", err)
	}
	if err := social.NewRepository(gormDB).AutoMigrate(); err != nil {
		log.Fatalf("auto migrate follows: %v", err)
	}
	if err := history.NewRepository(gormDB).AutoMigrate(); err != nil {
		log.Fatalf("auto migrate view histories: %v", err)
	}
	if err := notice.NewRepository(gormDB).AutoMigrate(); err != nil {
		log.Fatalf("auto migrate notices: %v", err)
	}
	if err := admin.NewRepository(gormDB).AutoMigrate(); err != nil {
		log.Fatalf("auto migrate video reviews: %v", err)
	}
	if err := areaRepo.SeedDefaults(context.Background()); err != nil {
		log.Fatalf("seed default areas: %v", err)
	}

	router := apphttp.NewRouter(cfg, gormDB, redisClient)
	server := &http.Server{
		Addr:              cfg.Server.Addr(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("pilipili-go api listening on %s", cfg.Server.Addr())
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run server: %v", err)
		}
	}()

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-signalCtx.Done()
	log.Println("shutdown signal received, draining in-flight requests")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		if err := server.Close(); err != nil {
			log.Printf("force close server: %v", err)
		}
	}
}
