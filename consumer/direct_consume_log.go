package main

import (
	"fmt"
	"rabbitmq/internal/db"
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

func getTopicRoutingKey(status string) string {
	return fmt.Sprintf("%s.%s.%s", warehouseState, warehouseCity, status)
}

func bindInfoTypeQueues(ch *amqp.Channel, q *amqp.Queue, status string) {
	ch.QueueBind(q.Name,
		getTopicRoutingKey(status),
		utils.EXCHANGE_NAME,
		false,
		nil,
	)
}

func main() {
	conn, err := utils.CreateAMQPServer()
	utils.FailOnError(err, "unable to create server")

	ch, err := conn.Channel()
	utils.FailOnError(err, "unable to create channel")

	err = utils.CreateExchange(ch, utils.EXCHANGE_NAME, "topic")
	utils.FailOnError(err, "unable to create message")

	q, err := utils.CreateQueue(ch, "inventory", true)
	utils.FailOnError(err, "unable to create queue")

	noOfCars := utils.GetEnvInt("WAREHOUSE_CARS", 10)
	ch.Qos(noOfCars, 0, false)

	dbConn := db.NewDbConn()
	db := NewShipmentDB(dbConn)

	forerver := make(chan int, 1)
	go ConsumeDeliveryMessages(ch, &q, db)
	go ConsumeInventoryMessages(ch, &q, db)
	<-forerver
}
