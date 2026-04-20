package worker

import (
	"context"
	"fmt"
	"log"

	redisclient "github.com/Fankemp/GameMatch/internal/redis"
	"github.com/redis/go-redis/v9"
)

type NotificationWorker struct {
	redis *redisclient.Client
}

func NewNotificationWorker(redis *redisclient.Client) *NotificationWorker {
	return &NotificationWorker{redis: redis}
}

func (w *NotificationWorker) Start(ctx context.Context) {
	log.Println("notification worker started")

	lastID := "$"

	for {
		select {
		case <-ctx.Done():
			log.Println("notification worker stopped")
			return
		default:
		}

		messages, err := w.redis.ReadMatchEvents(ctx, lastID)
		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Printf("read match events error: %v", err)
			continue
		}

		for _, msg := range messages {
			lastID = msg.ID
			w.processMatchEvent(msg)
		}
	}
}

func (w *NotificationWorker) processMatchEvent(msg redis.XMessage) {
	userAID := fmt.Sprintf("%v", msg.Values["user_a_id"])
	userBID := fmt.Sprintf("%v", msg.Values["user_b_id"])
	gameID := fmt.Sprintf("%v", msg.Values["game_id"])

	log.Printf("[MATCH] user %s and user %s matched in game %s", userAID, userBID, gameID)

	// TODO: отправить push-уведомление через FCM/APNs
}
