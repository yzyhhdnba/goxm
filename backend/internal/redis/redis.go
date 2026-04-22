package redis

import (
	"context"
	"fmt"
	"time"

	appconfig "pilipili-go/backend/internal/config"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	raw *goredis.Client
}

func New(cfg appconfig.RedisConfig) (*Client, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis addr is required")
	}

	raw := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	client := &Client{raw: raw}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		_ = raw.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.raw == nil {
		return fmt.Errorf("redis client is unavailable")
	}
	return c.raw.Ping(ctx).Err()
}

func (c *Client) Close() error {
	if c == nil || c.raw == nil {
		return nil
	}
	return c.raw.Close()
}

func (c *Client) Raw() *goredis.Client {
	if c == nil {
		return nil
	}
	return c.raw
}
