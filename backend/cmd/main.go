package main

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/seed"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	log := logger.New(&cfg.Logger)

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	db, err := database.Connect(&cfg.Database, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	if err := database.AutoMigrate(db.DB, log); err != nil {
		log.WithError(err).Fatal("Failed to run migrations")
	}

	if err := seed.Seed(db.DB, log, cfg.Server.AdminPassword); err != nil {
		log.WithError(err).Fatal("Failed to seed database")
	}

	weaponRepo := repository.NewWeaponRepository(db.DB, log, cfg.Server.PublicURL)
	modRepo := repository.NewModRepository(db.DB, log)
	ammoRepo := repository.NewAmmoRepository(db.DB, log)
	ingredientRepo := repository.NewIngredientRepository(db.DB, log)
	userRepo := repository.NewUserRepository(db.DB, log)
	suggestionRepo := repository.NewSuggestionRepository(db.DB, log)

	weaponHandler := handlers.NewWeaponHandler(weaponRepo, log)
	modHandler := handlers.NewModHandler(modRepo, log)
	ammoHandler := handlers.NewAmmoHandler(ammoRepo, log)
	ingredientHandler := handlers.NewIngredientHandler(ingredientRepo, log)
	authHandler := handlers.NewAuthHandler(userRepo, log, cfg.JWTSecret)
	suggestionHandler := handlers.NewSuggestionHandler(suggestionRepo, weaponRepo, log)

	auth := middleware.NewAuthMiddleware(cfg.JWTSecret, log)

	// Public routes
	http.HandleFunc("GET /api/weapons", weaponHandler.GetWeapons)
	http.HandleFunc("GET /api/weapons/{id}", weaponHandler.GetWeapon)
	http.HandleFunc("GET /api/my-weapons", auth.Authenticate(weaponHandler.ListMyWeapons))
	http.HandleFunc("GET /api/mods", modHandler.GetModList)
	http.HandleFunc("GET /api/mods/{id}", modHandler.GetMod)
	http.HandleFunc("GET /api/ammo", ammoHandler.GetAmmoList)
	http.HandleFunc("GET /api/ammo/{id}", ammoHandler.GetAmmo)
	http.HandleFunc("GET /api/ingredients", ingredientHandler.GetIngredientList)

	// Auth routes
	http.HandleFunc("POST /api/register", authHandler.Register)
	http.HandleFunc("POST /api/login", authHandler.Login)
	http.HandleFunc("GET /api/users/me", auth.Authenticate(authHandler.Me))

	// Suggestion routes
	http.HandleFunc("POST /api/suggestions", auth.Authenticate(suggestionHandler.Create))
	http.HandleFunc("GET /api/suggestions", auth.Authenticate(auth.RequireRole(suggestionHandler.List, "admin")))
	http.HandleFunc("GET /api/suggestions/{id}", auth.Authenticate(auth.RequireRole(suggestionHandler.Get, "admin")))
	http.HandleFunc("GET /api/suggestions/my", auth.Authenticate(suggestionHandler.ListMy))
	http.HandleFunc("GET /api/suggestions/my/{id}", auth.Authenticate(suggestionHandler.GetMy))
	http.HandleFunc("PUT /api/suggestions/{id}/resubmit", auth.Authenticate(suggestionHandler.Resubmit))
	http.HandleFunc("PUT /api/suggestions/{id}/approve", auth.Authenticate(auth.RequireRole(suggestionHandler.Approve, "admin")))
	http.HandleFunc("PUT /api/suggestions/{id}/reject", auth.Authenticate(auth.RequireRole(suggestionHandler.Reject, "admin")))
	http.HandleFunc("DELETE /api/suggestions/{id}", auth.Authenticate(auth.RequireRole(suggestionHandler.Delete, "admin")))

	// Admin-only weapon routes
	http.HandleFunc("POST /api/weapons", auth.Authenticate(auth.RequireRole(weaponHandler.CreateWeapon, "admin")))
	http.HandleFunc("DELETE /api/weapons/{id}", auth.Authenticate(auth.RequireRole(weaponHandler.DeleteWeapon, "admin")))

	// Upload route (admin only)
	http.HandleFunc("POST /api/upload", auth.Authenticate(auth.RequireRole(weaponHandler.UploadIcon, "admin")))

	http.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	mux := handlers.RecoveryMiddleware(log)(handlers.CORSMiddleware(cfg.Server.AllowedOrigin)(handlers.LoggingMiddleware(log)(http.DefaultServeMux)))
	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		log.Info("Server starting on :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("Server error")
		}
	}()

	<-stop

	log.Info("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server shutdown failed")
	}
	log.Info("Server stopped.")
}
