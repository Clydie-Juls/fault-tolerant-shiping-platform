package main

import (
	"context"
	"log"
	"rabbitmq/internal/db"
	"rabbitmq/utils"
	"strings"
	"time"
)

var (
	warehouseState = strings.ToLower(utils.GetEnvString("WAREHOUSE_STATE", "colorado"))
	warehouseCity  = strings.ToLower(utils.GetEnvString("WAREHOUSE_CITY", "denver"))
	latitude       = utils.GetEnvFloat("WAREHOUSE_LATITUDE", -40.0076)
	longitude      = utils.GetEnvFloat("WAREHOUSE_LONGITUDE", -105.2659)
	keyStatus      = func(status string) string {
		return utils.GetTopicRoutingKey(warehouseState, warehouseCity, status)
	}
)

func main() {
	conn, err := utils.CreateAMQPServer()
	utils.FailOnError(err, "unable to create server")

	ch, err := conn.Channel()
	utils.FailOnError(err, "unable to create channel")

	err = utils.CreateExchange(ch, utils.EXCHANGE_NAME, "topic")
	utils.FailOnError(err, "unable to create message")

	inventoryQueue, err := utils.CreateQueue(ch, "inventory", true)
	utils.FailOnError(err, "unable to create queue")

	deliveryQueue, err := utils.CreateQueue(ch, "delivery", true)
	utils.FailOnError(err, "unable to create queue")

	noOfCars := utils.GetEnvInt("WAREHOUSE_CARS", 10)
	ch.Qos(noOfCars, 0, false)

	dbConn := db.NewDbConn()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	warehouseID := utils.GetOrCreateWarehouseID(ctx, dbConn.DB, warehouseState, warehouseCity)

	db := NewShipmentDB(dbConn, warehouseID)

	forerver := make(chan int, 1)
	go ConsumeDeliveryMessages(ch, &deliveryQueue, db)
	go ConsumeInventoryMessages(ch, &inventoryQueue, db)
	log.Println("Consumer Service starting")
	<-forerver
}
