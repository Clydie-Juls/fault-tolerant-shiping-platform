package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	AMQP_DEFAULT_USER = GetEnvString("AMQP_USER", "guest")
	AMQP_HOST         = GetEnvString("AMQP_HOST", "127.0.0.1")
	AMQP_PORT         = GetEnvString("AMQP_PORT", "5672")
	RABBITMQ_SECRET   = ReadSecret("/run/secrets/rabbitmq_pass")
	AMQP_DEFAULT_PASS = GetSecretString(RABBITMQ_SECRET, "AMQP_PASS", "guest")
	AMQP_URL          = fmt.Sprintf("amqp://%s:%s@%s:%s",
		AMQP_DEFAULT_USER,
		AMQP_DEFAULT_PASS,
		AMQP_HOST,
		AMQP_PORT,
	)
)

const (
	EXCHANGE_NAME         = "inventory"
	STATUS_LISTED         = "listed"
	STATUS_BUY_REQUEST    = "buy_request"
	STATUS_ACCEPTED       = "accepted"
	STATUS_ARRIVED_PICKUP = "arrived_pickup"
	STATUS_BUY_ACCEPTED   = "buy_accepted"
	STATUS_DELIVERED      = "delivered"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func CreateAMQPServer() (*amqp.Connection, error) {
	conn, err := amqp.Dial(AMQP_URL)
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

func GetTopicRoutingKey(warehouseState, warehouseCity, status string) string {
	return fmt.Sprintf("%s.%s.%s", warehouseState, warehouseCity, status)
}

func BindInfoTypeQueues(ch *amqp.Channel, q *amqp.Queue, key string) error {
	return ch.QueueBind(
		q.Name,
		key,
		EXCHANGE_NAME,
		false,
		nil,
	)
}

func GetOrCreateWarehouseID(
	ctx context.Context,
	db *sql.DB,
	state string,
	city string,
) int {
	const query = `
		INSERT INTO warehouse (warehouse_state, warehouse_city)
		VALUES ($1, $2)
		ON CONFLICT (warehouse_state, warehouse_city)
		DO UPDATE SET warehouse_city = EXCLUDED.warehouse_city
		RETURNING id;
	`

	var warehouseID int
	err := db.QueryRowContext(ctx, query, state, city).Scan(&warehouseID)
	FailOnError(err, "unable to create or get warehouse id")

	return warehouseID
}
