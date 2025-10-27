package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"modernfi-treasury-app/internal/database"
	"modernfi-treasury-app/internal/handlers"
	"modernfi-treasury-app/internal/services"
)

const (
	// Server configuration
	serverPort         = ":8080"
	serverReadTimeout  = 15 * time.Second
	serverWriteTimeout = 15 * time.Second
	serverIdleTimeout  = 60 * time.Second
	shutdownTimeout    = 30 * time.Second

	// CORS configuration
	corsMaxAge = 300
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Database connection
	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Create connection pool
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	config.MaxConns = 25
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Database connection established")

	// Initialize sqlc queries
	queries := database.New(pool)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(queries)

	// Initialize TreasuryService
	treasuryService := services.NewTreasuryService()

	// Start cache warming in background (non-blocking - returns immediately)
	// Pre-fetches historical yield data for all periods (1W through 30Y)
	// so subsequent user requests are served instantly from cache
	treasuryService.WarmCache()

	// Initialize YieldHandler with service
	yieldHandler := handlers.NewYieldHandler(treasuryService)

	// Initialize TransactionService and handlers
	txService := services.NewTransactionService(queries, pool)
	txHandlers := handlers.NewTransactionHandlers(txService, queries, treasuryService)

	// Initialize HoldingsHandlers
	holdingsHandlers := handlers.NewHoldingsHandlers(queries)

	// Create chi router
	r := chi.NewRouter()

	// Get allowed origins from environment or use permissive defaults
	// Since we're using nginx proxy in production, most requests come through same-origin
	// But we still support direct API access for development and testing
	allowedOrigins := []string{
		"http://localhost:5173", // Vite dev server (default port)
		"http://localhost:5174", // Vite dev server (alternate port)
		"http://localhost:80",   // Local Docker frontend
		"http://localhost",      // Local Docker frontend (alt)
		"http://localhost:3000", // Alternative dev port
	}

	// Allow additional origins from environment (comma-separated)
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		for _, origin := range strings.Split(envOrigins, ",") {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}

	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           corsMaxAge,
	}))

	// Register routes
	r.Get("/api/v1/users", userHandler.GetAllUsers)
	r.Get("/api/v1/users/{userId}/transactions", txHandlers.GetUserTransactions)
	r.Get("/api/v1/users/{id}/holdings", holdingsHandlers.GetUserHoldings)

	// Historical yield data endpoint (must be registered before /api/yields)
	r.Get("/api/yields/historical", yieldHandler.GetHistoricalYields)
	// Current yield snapshot endpoint
	r.Get("/api/yields", yieldHandler.GetYields)

	r.Post("/api/v1/fund", txHandlers.FundHandler)
	r.Post("/api/v1/withdraw", txHandlers.WithdrawHandler)
	r.Post("/api/v1/buy", txHandlers.BuyHandler)
	r.Post("/api/v1/sell", txHandlers.SellHandler)

	// Health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Configure server
	server := &http.Server{
		Addr:         serverPort,
		Handler:      r,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
