// femabras/backend/cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Femabras/femabras/backend/internal/auth/provider"
	"github.com/Femabras/femabras/backend/internal/auth/repository"
	"github.com/Femabras/femabras/backend/internal/auth/service"
	"github.com/Femabras/femabras/backend/internal/config"
	"github.com/Femabras/femabras/backend/internal/database"
	"github.com/Femabras/femabras/backend/internal/middleware"
	"github.com/Femabras/femabras/backend/internal/routes"
	"github.com/Femabras/femabras/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Unable to locate the .env path: %v", err)
		return
	}

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
	authSvc := service.NewAuthService(authRepo, authFactory, &cfg)

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
	r.Use(middleware.RateLimit())

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
