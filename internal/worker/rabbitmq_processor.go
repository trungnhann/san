package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"san/pkg/logger"
	"san/pkg/mail"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQTaskProcessor struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mailer  mail.EmailSender
	logger  logger.Logger
}

func NewRabbitMQTaskProcessor(conn *amqp.Connection, mailer mail.EmailSender, logger logger.Logger) *RabbitMQTaskProcessor {
	return &RabbitMQTaskProcessor{
		conn:   conn,
		mailer: mailer,
		logger: logger,
	}
}

func (processor *RabbitMQTaskProcessor) Start() error {
	ch, err := processor.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	processor.channel = ch

	// Setup consumers for each task type
	if err := processor.consumeTask(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail); err != nil {
		return err
	}

	processor.logger.Info(" [*] RabbitMQ Worker started")
	return nil
}

func (processor *RabbitMQTaskProcessor) consumeTask(queueName string, handler func(context.Context, []byte) error) error {
	// Ensure queue exists
	q, err := processor.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	err = processor.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := processor.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer for %s: %w", queueName, err)
	}

	go func() {
		for d := range msgs {
			processor.logger.Infof("Received a message from %s: %s", queueName, d.Body)

			err := handler(context.Background(), d.Body)

			if err != nil {
				processor.logger.Errorf("Error processing task %s: %v", queueName, err)
				// Requeue the message
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (processor *RabbitMQTaskProcessor) Shutdown() {
	if processor.channel != nil {
		processor.channel.Close()
	}
	if processor.conn != nil {
		processor.conn.Close()
	}
}

func (processor *RabbitMQTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, payloadBytes []byte) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	subject := "Welcome to San App"
	content := fmt.Sprintf(`
	<h1>Welcome!</h1>
	<p>Please verify your email: %s</p>
	`, payload.Email)
	to := []string{payload.Email}

	err := processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	processor.logger.Infof("processed task send verify email: email=%s", payload.Email)
	return nil
}
