package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	appconfig "pilipili-go/backend/internal/config"
	"pilipili-go/backend/internal/db"
	"pilipili-go/backend/internal/mq/rocketmq"
	appredis "pilipili-go/backend/internal/redis"
	"pilipili-go/backend/internal/worker"
)

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

	mqClient, err := rocketmq.New(cfg.RocketMQ)
	if err != nil {
		log.Fatalf("connect rocketmq: %v", err)
	}

	runner := worker.NewRunner(gormDB, redisClient, mqClient)

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("pilipili-go worker starting")
	if err := runner.Run(signalCtx); err != nil {
		log.Fatalf("run worker: %v", err)
	}
	log.Printf("pilipili-go worker stopped")
}
