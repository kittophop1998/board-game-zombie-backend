package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kittiphop/zombie_board_game_back/infrastructure/config"
)

func main() {
	// ===== Load environment variables =====
	loadEnv()

	// ===== Load configuration =====
	cfg := loadConfig()

	// ===== Initialize application =====
	router := initializeAppOrExit(cfg)

	// ===== Start server with graceful shutdown =====
	startServer(router, cfg.Server.Port)
}

// --------------- Helper Functions ---------------
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: No .env file found, using system environment")
	}
}

func loadConfig() *config.Config {
	env := os.Getenv("ENV_NAME")
	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}
	return cfg
}

func initializeAppOrExit(cfg *config.Config) http.Handler {
	router, err := initializeApp(cfg)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize app: %v", err)
	}
	return router
}

func startServer(handler http.Handler, port int) {
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		log.Printf("INFO: Starting server on port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: Server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("INFO: Shutting down server...")

	ctx := context.Background()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("FATAL: Failed to gracefully shutdown: %v", err)
	}

	log.Println("INFO: Server stopped")
}
