package service

import (
	"context"
	"log"

	redisclient "github.com/Fankemp/GameMatch/internal/redis"
)

type NotificationService interface {
	NotifyMatch(ctx context.Context, userAID, userBID int64, gameID string)
}

type notificationService struct {
	redis *redisclient.Client
}

func NewNotificationService(redis *redisclient.Client) NotificationService {
	return &notificationService{redis: redis}
}

func (s *notificationService) NotifyMatch(ctx context.Context, userAID, userBID int64, gameID string) {
	if s.redis == nil {
		log.Printf("[MATCH] user %d and user %d matched in game %s (redis unavailable)", userAID, userBID, gameID)
		return
	}
	event := redisclient.MatchEvent{
		UserAID: userAID,
		UserBID: userBID,
		GameID:  gameID,
	}
	if err := s.redis.PublishMatchEvent(ctx, event); err != nil {
		log.Printf("failed to publish match event: %v", err)
	}
}
