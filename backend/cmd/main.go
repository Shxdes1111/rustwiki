package main

import (
	"context"
	"fmt"
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

	db, err := database.Connect(&cfg.Database, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.AutoMigrate(db.DB, log); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := seed.Seed(db.DB, log); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	weaponRepo := repository.NewWeaponRepository(db.DB, log)
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

	auth := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Public routes
	http.HandleFunc("GET /api/weapons", weaponHandler.GetWeapons)
	http.HandleFunc("GET /api/weapons/{id}", weaponHandler.GetWeapon)
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
	http.HandleFunc("GET /api/suggestions", auth.Authenticate(middleware.RequireRole(suggestionHandler.List, "admin")))
	http.HandleFunc("GET /api/suggestions/{id}", auth.Authenticate(middleware.RequireRole(suggestionHandler.Get, "admin")))
	http.HandleFunc("PUT /api/suggestions/{id}/approve", auth.Authenticate(middleware.RequireRole(suggestionHandler.Approve, "admin")))
	http.HandleFunc("PUT /api/suggestions/{id}/reject", auth.Authenticate(middleware.RequireRole(suggestionHandler.Reject, "admin")))

	// Admin-only weapon routes
	http.HandleFunc("POST /api/weapons", auth.Authenticate(middleware.RequireRole(weaponHandler.CreateWeapon, "admin")))
	http.HandleFunc("DELETE /api/weapons/{id}", auth.Authenticate(middleware.RequireRole(weaponHandler.DeleteWeapon, "admin")))

	// Upload route (admin only)
	http.HandleFunc("POST /api/upload", auth.Authenticate(middleware.RequireRole(weaponHandler.UploadIcon, "admin")))

	http.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{Addr: ":8080", Handler: handlers.CORSMiddleware(http.DefaultServeMux)}

	go func() {
		fmt.Println("Server starting on :8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	<-stop

	fmt.Println("\nShutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown Failed:%+v", err)
	}
	fmt.Println("Server stopped.")
}
