// mq/rabbitmq.go
package mq

import (
	"context"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

var conn *amqp.Connection

func InitRabbitMQ() (*amqp.Connection, error) {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Log.Error("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	return conn, nil
}

func SendToQueue(ctx context.Context, conn *amqp.Connection, queueName string, message []byte) error {
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
		Headers: amqp.Table{
			"trace_id": spanContext.TraceID().String(),
			"span_id":  spanContext.SpanID().String(),
		},
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		msg,
	)
	return err
}
