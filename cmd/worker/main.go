package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fankemp/GameMatch/internal/config"
	redisclient "github.com/Fankemp/GameMatch/internal/redis"
	"github.com/Fankemp/GameMatch/internal/worker"
)

func main() {
	cfg := config.NewConfig()

	redisClient, err := redisclient.NewClient(cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("worker started, waiting for events...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	notificationWorker := worker.NewNotificationWorker(redisClient)
	go notificationWorker.Start(ctx)

	<-quit
	log.Println("shutting down worker...")
	cancel()
}
