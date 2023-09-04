// mq/rabbitmq.go
package mq

import (
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"github.com/streadway/amqp"
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

func SendToQueue(conn *amqp.Connection, queueName string, message []byte) error {
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

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	return err
}
