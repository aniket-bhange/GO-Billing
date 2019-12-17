// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// failOnError(err, "Failed to connect to RabbitMQ")
// defer conn.Close()
package amq

import (
	config "billing-gorilla/core"
	"log"

	"github.com/streadway/amqp"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}

type RbmqConfig struct {
	q       amqp.Queue
	ch      *amqp.Channel
	conn    *amqp.Connection
	rbmqErr error
}

func ConnectRBMQ() *RbmqConfig {
	rbmqConfig := &RbmqConfig{}
	config := config.New()

	rbmqConfig.conn, rbmqConfig.rbmqErr = amqp.Dial(config.AMQ.Url)
	handleError(rbmqConfig.rbmqErr, "Failed to connect")

	rbmqConfig.ch, rbmqConfig.rbmqErr = rbmqConfig.conn.Channel()
	handleError(rbmqConfig.rbmqErr, "Failed to open channel")

	rbmqConfig.q, rbmqConfig.rbmqErr = rbmqConfig.ch.QueueDeclare(
		"BILL_QUEUE",
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(rbmqConfig.rbmqErr, "Failed to declare a queue")

	return rbmqConfig
}

func SendMessage(config *RbmqConfig, msg string) {
	config.rbmqErr = config.ch.Publish(
		"",
		config.q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(msg),
		},
	)

	handleError(config.rbmqErr, "Failed to publish a message")
}

func ConsumeMessage(config *RbmqConfig) (<-chan amqp.Delivery, error) {
	msgs, err := config.ch.Consume(
		config.q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	handleError(err, "Failed to register a consumer")

	return msgs, err

}
