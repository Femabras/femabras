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

	"github.com/Femabras/femabras/internal/auth/repository"
	"github.com/Femabras/femabras/internal/auth/service"
	challengeRepo "github.com/Femabras/femabras/internal/challenge/repository"
	challengeSvc "github.com/Femabras/femabras/internal/challenge/service"
	"github.com/Femabras/femabras/internal/config"
	"github.com/Femabras/femabras/internal/database"
	"github.com/Femabras/femabras/internal/middleware"
	"github.com/Femabras/femabras/internal/routes"
	"github.com/Femabras/femabras/internal/services"
	"github.com/Femabras/femabras/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	if cfg.JWTSecret == "" || cfg.DatabaseURL == "" {
		log.Fatal("Missing required environment variables (JWTSecret, DatabaseURL, etc.)")
	}

	redisURL := os.Getenv("REDIS_URL")
	redisOpt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URI: %v", err)
	}

	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	asynqServer := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TypeSendVerificationEmail, worker.HandleVerificationEmailTask)

	go func() {
		if err := asynqServer.Run(mux); err != nil {
			fmt.Printf("could not start asynq server: %v\n", err)
		}
	}()

	service.InitGoogleOAuth(&cfg)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	services.InitRedis(cfg.RedisURL)

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo, &cfg, asynqClient)

	challengeRepository := challengeRepo.NewChallengeRepository(db)
	challengeService := challengeSvc.NewChallengeService(challengeRepository)

	if _, err = services.CreateOrGetTodayChallenge(db); err != nil {
		log.Fatalf("Failed to initialize daily challenge: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.FrontendURL))

	routes.Setup(r, db, &cfg, authSvc, challengeService)

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
