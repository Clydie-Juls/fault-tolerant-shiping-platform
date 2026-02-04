package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rabbitmq/models"
	"rabbitmq/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeDeliveryMessages(ch *amqp.Channel, q *amqp.Queue, db *ShipmentDB) {
	location := models.Location{
		Longitude: longitude,
		Latitude:  latitude,
	}
	warehouse := models.NewWarehouse(location)
	err := utils.BindInfoTypeQueues(ch, q, keyStatus(utils.STATUS_BUY_REQUEST))
	utils.FailOnError(err, "unable to bind queue")

	forever := make(chan int)

	msgs, err := ch.Consume(q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "unable to consume message")

	log.Println("waiting to consume delivery messages")
	for msg := range msgs {
		go processDelivery(ch, &msg, *warehouse, db)
	}
	<-forever
}

func processDelivery(ch *amqp.Channel, msg *amqp.Delivery, warehouse models.Warehouse, db *ShipmentDB) {
	var order models.Order
	err := json.Unmarshal(msg.Body, &order)
	utils.FailOnError(err, "unable to unmarshal order json")
	estimateTime := warehouse.EstimateTimeToDestination(order.Location)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	inv := db.ConsumeInventoryByID(ctx, int64(order.InventoryID))
	if inv == nil {
		log.Println("unable to get inventory")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,
		utils.EXCHANGE_NAME,
		keyStatus(utils.STATUS_BUY_ACCEPTED),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        fmt.Appendf(nil, "buy accepted: reaching destination in about %.f seconds", estimateTime.Seconds()),
		},
	)
	utils.FailOnError(err, "unable to publish message")

	// reaching destination
	time.Sleep(estimateTime)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	invByte, err := json.Marshal(inv)
	utils.FailOnError(err, "unable to to marshal inventory")
	err = ch.PublishWithContext(ctx,

		utils.EXCHANGE_NAME,
		keyStatus(utils.STATUS_DELIVERED),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        invByte,
		},
	)
	utils.FailOnError(err, "unable to publish message")

	// going back to the warehouse
	time.Sleep(estimateTime)
	msg.Ack(false)
}
