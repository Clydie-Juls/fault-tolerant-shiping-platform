package main

import (
	"rabbitmq/utils"
)

const EXCHANGE_NAME = "inventory"

var addr = utils.GetEnvString("CONSUMER_ADDR", ":8080")

type connection interface {
	Close() error
}

func main() {
	conn, err := utils.CreateAMQPServer()
	utils.FailOnError(err, "unable to create server")

	ch, err := conn.Channel()
	utils.FailOnError(err, "unable to create channel")

	err = utils.CreateExchange(ch, EXCHANGE_NAME, "topic")
	utils.FailOnError(err, "unable to create message")

	createHTTPServer(ch, conn)
}
