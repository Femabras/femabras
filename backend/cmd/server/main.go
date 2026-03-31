// femabras/backend/cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Femabras/femabras/internal/auth/provider"
	"github.com/Femabras/femabras/internal/auth/repository"
	"github.com/Femabras/femabras/internal/auth/service"
	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/database"
	"github.com/Femabras/femabras/internal/middleware"
	"github.com/Femabras/femabras/internal/routes"
	"github.com/Femabras/femabras/internal/services"
	"github.com/Femabras/femabras/internal/worker"
	"github.com/hibiken/asynq"

	"github.com/gin-gonic/gin"
)

func main() {

	redisURL := os.Getenv("REDIS_URL") // From Railway
	redisOpt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		panic("Failed to parse Redis URI: " + err.Error())
	}

	// 2. Initialize Asynq Client (for enqueuing tasks)
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	// 3. Initialize Asynq Server (the background worker)
	asynqServer := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10, // Process up to 10 emails at the exact same time
	})

	// 4. Map the tasks to their handler functions
	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TypeSendVerificationEmail, worker.HandleVerificationEmailTask)

	// 5. Start the worker in a separate Goroutine so it doesn't block the API!
	go func() {
		if err := asynqServer.Run(mux); err != nil {
			fmt.Printf("could not start asynq server: %v\n", err)
		}
	}()
	cfg := config.Load()

	if cfg.JWTSecret == "" || cfg.DatabaseURL == "" {
		log.Fatal("Missing required environment variables (JWTSecret, DatabaseURL, etc.)")
	}

	// Initialize Infrastructure
	service.InitGoogleOAuth(&cfg)
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	services.InitRedis(cfg.RedisURL)

	// Initialize Auth Dependencies (Moved out of routes to be accessible here)
	authRepo := repository.NewAuthRepository(db)
	authFactory := provider.NewFactory(&cfg)
	authSvc := service.NewAuthService(authRepo, authFactory, &cfg, asynqClient)
	// Initialize Challenges
	_, err = services.CreateOrGetTodayChallenge(db)
	if err != nil {
		log.Panicf("Failed to initialize daily challenge: %v", err)
		return
	}

	// Setup Router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.FrontendURL))
	// r.Use(middleware.RateLimit())

	routes.Setup(r, db, &cfg, authSvc)

	// Server Management
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Server starting on :%s (Release Mode)", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful Shutdown Logic
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited gracefully")
}
