package video

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	appredis "pilipili-go/backend/internal/redis"

	goredis "github.com/redis/go-redis/v9"
)

const (
	videoDetailCacheTTL = 30 * time.Second
	hotFeedCacheTTL     = 15 * time.Second
	hotFeedVersionKey   = "feed:hot:version"
)

type Cache interface {
	GetDetail(ctx context.Context, videoID uint64) (DetailResponse, bool, error)
	SetDetail(ctx context.Context, detail DetailResponse) error
	GetHotFeed(ctx context.Context, cursor string, limit int) (FeedResponse, bool, error)
	SetHotFeed(ctx context.Context, cursor string, limit int, feed FeedResponse) error
	InvalidateVideo(ctx context.Context, videoID uint64) error
	InvalidateHotFeed(ctx context.Context) error
}

type RedisCache struct {
	client *appredis.Client
}

func NewRedisCache(client *appredis.Client) *RedisCache {
	if client == nil || client.Raw() == nil {
		return nil
	}
	return &RedisCache{client: client}
}

func (c *RedisCache) GetDetail(ctx context.Context, videoID uint64) (DetailResponse, bool, error) {
	raw := c.raw()
	if raw == nil {
		return DetailResponse{}, false, nil
	}

	payload, err := raw.Get(ctx, detailCacheKey(videoID)).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return DetailResponse{}, false, nil
		}
		return DetailResponse{}, false, err
	}

	var detail DetailResponse
	if err := json.Unmarshal(payload, &detail); err != nil {
		return DetailResponse{}, false, err
	}
	return detail, true, nil
}

func (c *RedisCache) SetDetail(ctx context.Context, detail DetailResponse) error {
	raw := c.raw()
	if raw == nil {
		return nil
	}

	payload, err := json.Marshal(detail)
	if err != nil {
		return err
	}
	return raw.Set(ctx, detailCacheKey(detail.ID), payload, videoDetailCacheTTL).Err()
}

func (c *RedisCache) GetHotFeed(ctx context.Context, cursor string, limit int) (FeedResponse, bool, error) {
	raw := c.raw()
	if raw == nil {
		return FeedResponse{}, false, nil
	}

	key, err := c.hotFeedCacheKey(ctx, cursor, limit)
	if err != nil {
		return FeedResponse{}, false, err
	}

	payload, err := raw.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return FeedResponse{}, false, nil
		}
		return FeedResponse{}, false, err
	}

	var feed FeedResponse
	if err := json.Unmarshal(payload, &feed); err != nil {
		return FeedResponse{}, false, err
	}
	return feed, true, nil
}

func (c *RedisCache) SetHotFeed(ctx context.Context, cursor string, limit int, feed FeedResponse) error {
	raw := c.raw()
	if raw == nil {
		return nil
	}

	key, err := c.hotFeedCacheKey(ctx, cursor, limit)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(feed)
	if err != nil {
		return err
	}
	return raw.Set(ctx, key, payload, hotFeedCacheTTL).Err()
}

func (c *RedisCache) InvalidateVideo(ctx context.Context, videoID uint64) error {
	raw := c.raw()
	if raw == nil {
		return nil
	}
	return raw.Del(ctx, detailCacheKey(videoID)).Err()
}

func (c *RedisCache) InvalidateHotFeed(ctx context.Context) error {
	raw := c.raw()
	if raw == nil {
		return nil
	}
	return raw.Incr(ctx, hotFeedVersionKey).Err()
}

func (c *RedisCache) hotFeedCacheKey(ctx context.Context, cursor string, limit int) (string, error) {
	version, err := c.hotFeedVersion(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"feed:hot:v%s:cursor:%s:limit:%d",
		version,
		url.QueryEscape(cursor),
		limit,
	), nil
}

func (c *RedisCache) hotFeedVersion(ctx context.Context) (string, error) {
	raw := c.raw()
	if raw == nil {
		return "1", nil
	}

	version, err := raw.Get(ctx, hotFeedVersionKey).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "1", nil
		}
		return "", err
	}
	if _, err := strconv.ParseInt(version, 10, 64); err != nil {
		return "1", nil
	}
	return version, nil
}

func (c *RedisCache) raw() *goredis.Client {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Raw()
}

func detailCacheKey(videoID uint64) string {
	return fmt.Sprintf("video:detail:%d", videoID)
}
