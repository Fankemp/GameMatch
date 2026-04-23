package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Fankemp/GameMatch/internal/config"
	"github.com/redis/go-redis/v9"
)

const (
	FeedCacheTTL      = 5 * time.Minute
	MatchEventsStream = "match_events"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(cfg *config.RedisConfig) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

// Feed cache

func feedKey(userID int64, gameID string, limit, offset int) string {
	return fmt.Sprintf("feed:%d:%s:%d:%d", userID, gameID, limit, offset)
}

func (c *Client) GetFeedCache(ctx context.Context, userID int64, gameID string, limit, offset int) ([]byte, error) {
	key := feedKey(userID, gameID, limit, offset)
	val, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *Client) SetFeedCache(ctx context.Context, userID int64, gameID string, limit, offset int, data any) error {
	key := feedKey(userID, gameID, limit, offset)
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal feed cache: %w", err)
	}
	return c.rdb.Set(ctx, key, b, FeedCacheTTL).Err()
}

func (c *Client) InvalidateFeedCache(ctx context.Context, userID int64, gameID string) error {
	pattern := fmt.Sprintf("feed:%d:%s:*", userID, gameID)
	keys, err := c.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return c.rdb.Del(ctx, keys...).Err()
}

// Match events stream

type MatchEvent struct {
	UserAID int64  `json:"user_a_id"`
	UserBID int64  `json:"user_b_id"`
	GameID  string `json:"game_id"`
}

func (c *Client) PublishMatchEvent(ctx context.Context, event MatchEvent) error {
	return c.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: MatchEventsStream,
		Values: map[string]any{
			"user_a_id": event.UserAID,
			"user_b_id": event.UserBID,
			"game_id":   event.GameID,
		},
	}).Err()
}

func (c *Client) ReadMatchEvents(ctx context.Context, lastID string) ([]redis.XMessage, error) {
	streams, err := c.rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{MatchEventsStream, lastID},
		Count:   10,
		Block:   2 * time.Second,
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(streams) == 0 {
		return nil, nil
	}
	return streams[0].Messages, nil
}
