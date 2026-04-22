package worker

import (
	"context"
	"fmt"
	"log"

	"pilipili-go/backend/internal/mq/rocketmq"
	appredis "pilipili-go/backend/internal/redis"

	"gorm.io/gorm"
)

type Runner struct {
	db    *gorm.DB
	redis *appredis.Client
	mq    *rocketmq.Client
}

func NewRunner(db *gorm.DB, redis *appredis.Client, mq *rocketmq.Client) *Runner {
	return &Runner{
		db:    db,
		redis: redis,
		mq:    mq,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	if r == nil {
		return fmt.Errorf("worker runner is nil")
	}
	if r.db == nil {
		return fmt.Errorf("worker database is nil")
	}
	if r.redis == nil {
		return fmt.Errorf("worker redis is nil")
	}

	log.Printf("worker dependencies ready: mysql=ok redis=ok rocketmq=%s", dependencyState(r.mq))

	if r.mq != nil && r.mq.Enabled() {
		log.Printf("rocketmq worker topics prepared: %s, %s, %s",
			r.mq.Topic(TopicMediaTranscode),
			r.mq.Topic(TopicNoticeDispatch),
			r.mq.Topic(TopicPopularitySync),
		)
		log.Printf("rocketmq consumer groups prepared: %s, %s, %s",
			r.mq.ConsumerGroup("media"),
			r.mq.ConsumerGroup("notice"),
			r.mq.ConsumerGroup("popularity"),
		)
	} else {
		log.Printf("rocketmq is disabled; worker is running in scaffold mode")
	}

	<-ctx.Done()
	return nil
}

func dependencyState(client *rocketmq.Client) string {
	if client == nil || !client.Enabled() {
		return "disabled"
	}
	return "ok"
}
