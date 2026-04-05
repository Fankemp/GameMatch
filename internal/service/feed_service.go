package service

import (
	"context"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/repository"
)

type FeedService interface {
	GetFeed(ctx context.Context, userID int64, gameID string, limit, offset int) ([]*model.FeedCard, error)
}

type feedService struct {
	cardRepo repository.CardRepository
}

func NewFeedService(cardRepo repository.CardRepository) FeedService {
	return &feedService{cardRepo: cardRepo}
}

func (s *feedService) GetFeed(ctx context.Context, userID int64, gameID string, limit, offset int) ([]*model.FeedCard, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	return s.cardRepo.GetFeed(ctx, userID, gameID, limit, offset)
}
