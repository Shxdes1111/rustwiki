package main

import (
    "os"
	"syscall"
	"os/signal"
	"net/http"
	"fmt"
	"context"
	"time"
    "database/sql"
    "encoding/json"
    _ "github.com/lib/pq"
)

type Env struct {
    db *sql.DB
}

// This is now a "method" of Env
func (en *Env) getWeapons(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

    rows, err := en.db.Query("SELECT id, name, type FROM weapon_item")
    if err != nil {
        http.Error(w, "Database error", 500)
        return
    }
    defer rows.Close()

    var weapons []map[string]interface{}
    for rows.Next() {
        var id int
        var name, wType string
        rows.Scan(&id, &name, &wType)
        weapons = append(weapons, map[string]interface{}{
            "id": id, "name": name, "type": wType,
        })
    }
    json.NewEncoder(w).Encode(weapons)
}

func main() {
    // 1. Setup DB
    dbURL := os.Getenv("DB_URL")
    db, _ := sql.Open("postgres", dbURL)
    
    // 2. Initialize your environment
    env := &Env{db: db}

    // 3. Register routes using the env methods
    http.HandleFunc("/api/weapons", env.getWeapons)
    // 1. Create a channel to listen for OS signals
    stop := make(chan os.Signal, 1)

    // 2. Relay Ctrl+C (SIGINT) and termination (SIGTERM) to our channel
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    server := &http.Server{Addr: ":8080"}

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