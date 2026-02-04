package main

import (
	"log"
	"rabbitmq/utils"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	warehouseState = strings.ToLower(utils.GetEnvString("WAREHOUSE_STATE", "colorado"))
	warehouseCity  = strings.ToLower(utils.GetEnvString("WAREHOUSE_CITY", "denver"))
	longitude      = utils.GetEnvFloat("WAREHOUSE_LONGITUDE", -40.0076)
	latitude       = utils.GetEnvFloat("WAREHOUSE_LATITUDE", -105.2659)
)

const EXCHANGE_NAME = "inventory"

var (
	addr      = utils.GetEnvString("CONSUMER_ADDR", ":8080")
	keyStatus = func(status string) string {
		return utils.GetTopicRoutingKey(warehouseState, warehouseCity, status)
	}
)

type connection interface {
	Close() error
}

func ConsumeResponses(ch *amqp.Channel, q *amqp.Queue) {
	utils.BindInfoTypeQueues(ch, q, keyStatus(utils.STATUS_ACCEPTED))
	utils.BindInfoTypeQueues(ch, q, keyStatus(utils.STATUS_ARRIVED_PICKUP))
	utils.BindInfoTypeQueues(ch, q, keyStatus(utils.STATUS_BUY_ACCEPTED))
	utils.BindInfoTypeQueues(ch, q, keyStatus(utils.STATUS_DELIVERED))

	listedMsgs, err := ch.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "unable to consume message")

	for msg := range listedMsgs {
		log.Printf("[x] %s\n", string(msg.Body))
	}
}

func main() {
	conn, err := utils.CreateAMQPServer()
	utils.FailOnError(err, "unable to create server")

	ch, err := conn.Channel()
	utils.FailOnError(err, "unable to create channel")

	err = utils.CreateExchange(ch, EXCHANGE_NAME, "topic")
	utils.FailOnError(err, "unable to create message")

	q, err := utils.CreateQueue(ch, "responses", true)
	utils.FailOnError(err, "unable to create queue")

	go ConsumeResponses(ch, &q)
	createHTTPServer(ch, conn)
}
