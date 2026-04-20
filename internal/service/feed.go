package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Fankemp/GameMatch/internal/model"
	redisclient "github.com/Fankemp/GameMatch/internal/redis"
	"github.com/Fankemp/GameMatch/internal/repository"
)

type feedService struct {
	cardRepo repository.CardRepository
	redis    *redisclient.Client
}

func NewFeedService(cardRepo repository.CardRepository, redis *redisclient.Client) FeedService {
	return &feedService{
		cardRepo: cardRepo,
		redis:    redis,
	}
}

func (s *feedService) GetFeed(ctx context.Context, userID int64, gameID string, limit, offset int) ([]*model.FeedCard, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// Try cache
	if s.redis != nil {
		if cached, err := s.redis.GetFeedCache(ctx, userID, gameID, limit, offset); err == nil {
			var cards []*model.FeedCard
			if err := json.Unmarshal(cached, &cards); err == nil {
				return cards, nil
			}
		}
	}

	// Cache miss — go to DB
	cards, err := s.cardRepo.GetFeed(ctx, userID, gameID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Save to cache (non-blocking)
	if s.redis != nil {
		go func() {
			if err := s.redis.SetFeedCache(ctx, userID, gameID, limit, offset, cards); err != nil {
				log.Printf("set feed cache error: %v", err)
			}
		}()
	}

	return cards, nil
}
