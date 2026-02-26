package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	conf "san/internal/config"
	"san/internal/db"
	"san/internal/handler"
	"san/internal/server"
	"san/internal/service"
	storage_service "san/internal/service/storage"
	"san/internal/storage"
	"san/pkg/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title           San API
// @version         1.0
// @description     This is a sample server for San application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3001
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "dev" {
		env = "development"
	}

	config := conf.LoadConfig(env, "./env")
	log := buildLogger(env)

	database := db.NewDatabase(config)
	defer database.Close()

	database.Migrate()

	fileStorage, err := storage.NewStorage(config, log)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	activeStorageService := storage_service.NewActiveStorageService(database.Queries, fileStorage, log)

	userService := service.NewUserService(database.Queries, activeStorageService, log)

	userHandler := handler.NewUserHandler(userService)

	srv := server.NewServer(config, database, userHandler, log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Errorf("server shutdown error: %v", err)
	}

	log.Info("server stopped gracefully")
}

func buildLogger(env string) logger.Logger {
	if env == "test" {
		atom := zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.Level = atom
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		l, _ := zapConfig.Build()
		return l.Sugar()
	}

	if env == "development" {
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		l, _ := zapConfig.Build()
		return l.Sugar()
	}

	l, _ := zap.NewProduction()
	return l.Sugar()
}
