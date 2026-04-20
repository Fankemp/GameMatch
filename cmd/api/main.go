package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Fankemp/GameMatch/internal/config"
	"github.com/Fankemp/GameMatch/internal/db_conn"
	"github.com/Fankemp/GameMatch/internal/delivery/http"
	redisclient "github.com/Fankemp/GameMatch/internal/redis"
	"github.com/Fankemp/GameMatch/internal/repository"
	"github.com/Fankemp/GameMatch/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "run database migrations")
	migrateFlagDown := flag.Bool("migrate-down", false, "rollback last migrations")
	flag.Parse()

	cfg := config.NewConfig()

	db, err := db_conn.NewDB(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if *migrateFlag {
		if err = db.Migrate(); err != nil {
			log.Fatalln(err)
		}
		log.Println("migrations applied successfully")
		return
	}

	if *migrateFlagDown {
		if err = db.MigrateDown(); err != nil {
			log.Println("migration rollback failed:", err)
			return
		}
		log.Println("migration rolled back successfully")
		return
	}

	// Redis
	redisClient, err := redisclient.NewClient(cfg.Redis)
	if err != nil {
		log.Printf("warning: redis unavailable: %v", err)
		redisClient = nil
	} else {
		defer redisClient.Close()
		log.Println("redis connected")
	}

	// Repositories
	userRepo := repository.NewUserRepository(db.DB)
	profileRepo := repository.NewProfileRepository(db.DB)
	cardRepo := repository.NewCardRepository(db.DB)
	swipeRepo := repository.NewSwipeRepository(db.DB)
	matchRepo := repository.NewMatchRepository(db.DB)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWT.Secret)
	profileSvc := service.NewProfileService(profileRepo)
	cardSvc := service.NewCardService(cardRepo)
	feedSvc := service.NewFeedService(cardRepo, redisClient)
	notificationSvc := service.NewNotificationService(redisClient)
	swipeSvc := service.NewSwipeService(swipeRepo, cardRepo, matchRepo, notificationSvc)

	// Handlers
	authHandler := http.NewAuthHandler(authSvc)
	profileHandler := http.NewProfileHandler(profileSvc)
	cardHandler := http.NewCardHandler(cardSvc)
	feedHandler := http.NewFeedHandler(feedSvc)
	swipeHandler := http.NewSwipeHandler(swipeSvc)
	matchHandler := http.NewMatchHandler(matchRepo)

	// Router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/auth/register", authHandler.SignUp)
		api.POST("/auth/login", authHandler.SignIn)

		// Protected routes
		protected := api.Group("")
		protected.Use(http.JWTMiddleware(cfg.JWT.Secret))
		{
			// Auth
			protected.GET("/auth/me", authHandler.Me)

			// Profile
			protected.POST("/profile", profileHandler.CreateProfile)
			protected.PUT("/profile", profileHandler.UpdateProfile)
			protected.GET("/profile", profileHandler.GetMyProfile)
			protected.GET("/profiles/:id", profileHandler.GetProfileByID)

			// Cards
			protected.POST("/cards", cardHandler.Create)
			protected.GET("/cards", cardHandler.GetMyCards)
			protected.PUT("/cards/:id", cardHandler.Update)
			protected.DELETE("/cards/:id", cardHandler.Delete)

			// Feed
			protected.GET("/feed/:game_id", feedHandler.GetFeed)

			// Swipes
			protected.POST("/swipes", swipeHandler.Swipe)

			// Matches
			protected.GET("/matches", matchHandler.GetMatches)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("server starting on %s", addr)
	if err = r.Run(addr); err != nil {
		log.Fatalln(err)
	}
}
