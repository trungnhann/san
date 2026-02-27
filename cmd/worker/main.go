package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	conf "san/internal/config"
	"san/internal/worker"
	"san/pkg/logger"
	"san/pkg/mail"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "dev" {
		env = "development"
	}

	config := conf.LoadConfig(env, "./env")
	log := buildLogger(env)

	log.Info("starting worker")

	amqpConn, err := amqp.Dial(config.RabbitMQSource)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	gmailSender := mail.NewSmtpSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword, config.EmailSenderHost, config.EmailSenderPort, config.EmailSenderUsername)
	taskProcessor := worker.NewRabbitMQTaskProcessor(amqpConn, gmailSender, log)

	log.Info("starting task processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatalf("failed to start task processor: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Info("shutting down worker")
	taskProcessor.Shutdown()
	log.Info("worker stopped")
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

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	l, _ := config.Build()
	return l.Sugar()
}
