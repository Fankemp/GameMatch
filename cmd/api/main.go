package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Fankemp/GameMatch/internal/config"
	"github.com/Fankemp/GameMatch/internal/db_conn"
	"github.com/Fankemp/GameMatch/internal/handler"
	"github.com/Fankemp/GameMatch/internal/repository"
	"github.com/Fankemp/GameMatch/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// Repositories
	userRepo := repository.NewUserRepository(db.DB)
	profileRepo := repository.NewProfileRepository(db.DB)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWT.Secret)
	profileSvc := service.NewProfileService(profileRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	profileHandler := handler.NewProfileHandler(profileSvc)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(handler.JWTMiddleware(cfg.JWT.Secret))

			// Auth
			r.Get("/auth/me", authHandler.Me)

			// Profile
			r.Post("/profile", profileHandler.CreateProfile)
			r.Put("/profile", profileHandler.UpdateProfile)
			r.Get("/profile", profileHandler.GetMyProfile)
			r.Get("/profiles/{id}", profileHandler.GetProfileByID)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("server starting on %s", addr)
	if err = http.ListenAndServe(addr, r); err != nil {
		log.Fatalln(err)
	}
}
