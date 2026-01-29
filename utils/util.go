package utils

import (
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

const FALLBACK_AMQP_URL = "amqp://guest:guest@127.0.0.1:5672/"
const (
	EXCHANGE_NAME         = "inventory"
	STATUS_LISTED         = "listed"
	STATUS_ACCEPTED       = "accepted"
	STATUS_ARRIVED_PICKUP = "arrived_pickup"
	STATUS_BUY_REQUEST    = "buy_request"
	STATUS_BUY_ACCEPTED   = "buy_accepted"
	STATUS_DELIVERED      = "delivered"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func CreateAMQPServer() (*amqp.Connection, error) {
	url := GetEnvString("AMQP_URL", FALLBACK_AMQP_URL)
	conn, err := amqp.Dial(url)
	return conn, err
}

func CreateQueue(ch *amqp.Channel, name string, isDurable bool) (amqp.Queue, error) {
	// create a queue to start sending messages
	q, err := ch.QueueDeclare(
		name,
		isDurable,
		false,
		false,
		false,
		nil,
	)

	return q, err
}

func CreateExchange(ch *amqp.Channel, exchangeName string, kind string) error {
	return ch.ExchangeDeclare(
		exchangeName,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func BindQueue(ch *amqp.Channel, queueName string, exchangeName string, routingKey string) error {
	return ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}

func SeverityFrom(args []string) string {
	var s string
	if len(args) < 2 || args[1] == "" {
		s = "*.orange.*"
	} else {
		s = args[1]
	}

	return s
}

func BodyFrom(args []string) string {
	var s string
	if len(args) < 3 || args[2] == "" {
		s = "yo mama"
	} else {
		s = strings.Join(args[2:], " | ")
	}

	return s
}
