package main

import (
	"os"
	"syscall"
	"os/signal"
	"net/http"
	"fmt"
	"context"
	"time"
)

func main() {
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