package main

import (
    "os"
	"syscall"
	"os/signal"
	"net/http"
	"fmt"
	"context"
	"time"
    "backend/internal/config"
    "backend/internal/logger"
    "backend/internal/database"
    "backend/internal/handlers"
    "backend/internal/repository"
    "backend/internal/seed"
    _ "github.com/lib/pq"
)

func main() {
    // 1. Setup DB
    cfg := config.Load()
    log := logger.New(&cfg.Logger)
    
    db, err := database.Connect(&cfg.Database, log)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    if err := seed.Seed(db.DB, log); err != nil {
        log.Fatalf("Failed to seed database: %v", err)
    }

    // 2. Initialize your environment
    weaponRepo := repository.NewWeaponRepository(db.DB, log)
    modRepo := repository.NewModRepository(db.DB, log)
	ammoRepo := repository.NewAmmoRepository(db.DB, log)
	ingredientRepo := repository.NewIngredientRepository(db.DB, log)

    weaponHandler := handlers.NewWeaponHandler(weaponRepo, log)
    modHandler := handlers.NewModHandler(modRepo, log)
	ammoHandler := handlers.NewAmmoHandler(ammoRepo, log)
	ingredientHandler := handlers.NewIngredientHandler(ingredientRepo, log)

    // 3. Register routes using the env methods
    http.HandleFunc("GET /api/weapons", weaponHandler.GetWeapons)
    http.HandleFunc("GET /api/weapons/{id}", weaponHandler.GetWeapon)
    http.HandleFunc("GET /api/mods", modHandler.GetModList)
    http.HandleFunc("GET /api/mods/{id}", modHandler.GetMod)
    http.HandleFunc("GET /api/ammo", ammoHandler.GetAmmoList)
    http.HandleFunc("GET /api/ammo/{id}", ammoHandler.GetAmmo)
    http.HandleFunc("GET /api/ingredients", ingredientHandler.GetIngredientList)
    http.HandleFunc("POST /api/weapons", weaponHandler.CreateWeapon)
    http.HandleFunc("POST /api/upload", weaponHandler.UploadIcon)
    http.HandleFunc("DELETE /api/weapons/{id}", weaponHandler.DeleteWeapon)
    http.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
    // 1. Create a channel to listen for OS signals
    stop := make(chan os.Signal, 1)

    // 2. Relay Ctrl+C (SIGINT) and termination (SIGTERM) to our channel
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    server := &http.Server{Addr: ":8080", Handler: handlers.CORSMiddleware(http.DefaultServeMux)}

    // 3. Run the server in a goroutine so it doesn't block the main thread
    go func() {
        fmt.Println("Server starting on :8080...")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Printf("Error: %v\n", err)
        }
    }()

    // 4. BLOCK HERE: The program stops and waits for the channel to receive a signal
    <-stop

    // 5. Cleanup (Graceful Shutdown)
    fmt.Println("\nShutting down gracefully...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        fmt.Printf("Server Shutdown Failed:%+v", err)
    }
    fmt.Println("Server stopped.")
}