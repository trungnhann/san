package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Define task names
const (
	TaskSendVerifyEmail = "task:send_verify_email"
	// TaskSendWelcomeEmail = "task:send_welcome_email" // Example of a new task
)

// Define payloads
type PayloadSendVerifyEmail struct {
	Email string `json:"email"`
}

// Example of a generic interface if you want to support any payload,
// but for type safety, specific methods are better.
type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RabbitMQTaskDistributor struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQTaskDistributor(conn *amqp.Connection) (TaskDistributor, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQTaskDistributor{
		conn:    conn,
		channel: ch,
	}, nil
}

func (distributor *RabbitMQTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	return distributor.publish(ctx, TaskSendVerifyEmail, jsonPayload)
}

// Generic publish method to reuse logic
func (distributor *RabbitMQTaskDistributor) publish(ctx context.Context, queueName string, body []byte) error {
	// Declare queue on the fly or ensure it exists.
	// In production, queues are often pre-declared or declared by consumers.
	// But declaring here is safe and idempotent.
	_, err := distributor.channel.QueueDeclare(
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

	err = distributor.channel.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		})
	if err != nil {
		return fmt.Errorf("failed to publish message to %s: %w", queueName, err)
	}

	return nil
}
