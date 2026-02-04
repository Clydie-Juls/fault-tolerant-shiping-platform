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

func ConsumeInventoryMessages(ch *amqp.Channel, q *amqp.Queue, db *ShipmentDB) {
	location := models.Location{
		Longitude: longitude,
		Latitude:  latitude,
	}
	warehouse := models.NewWarehouse(location)
	utils.BindInfoTypeQueues(ch, q, keyStatus(utils.STATUS_LISTED))

	forever := make(chan int)

	msgs, err := ch.Consume(q.Name,
		"",
		false,
		false, false,
		false,
		nil,
	)
	utils.FailOnError(err, "unable to consume message")

	log.Printf("waiting to consume inventory messages at routing key: %s\n", keyStatus(utils.STATUS_LISTED))
	for msg := range msgs {
		go processInventory(ch, &msg, *warehouse, db)
	}
	<-forever
}

func processInventory(ch *amqp.Channel, msg *amqp.Delivery, warehouse models.Warehouse, db *ShipmentDB) {
	var shipment models.Shipment
	err := json.Unmarshal(msg.Body, &shipment)
	if err != nil {
		log.Printf("unable to unmarshal shipment json: %s", err)
		msg.Nack(false, false)
		return
	}
	estimateTime := warehouse.EstimateTimeToDestination(shipment.Location)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,
		utils.EXCHANGE_NAME,
		keyStatus(utils.STATUS_ACCEPTED),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("pickup: reaching destination in about %.f seconds", estimateTime.Seconds())),
		},
	)
	if err != nil {
		log.Printf("unable to publish message: %s", err)
		return
	}
	log.Println("status accepted request notif sent")

	// reaching destination
	time.Sleep(estimateTime)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(ctx,

		utils.EXCHANGE_NAME,
		keyStatus(utils.STATUS_ARRIVED_PICKUP),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("arrived pickup: destination reached"),
		},
	)
	if err != nil {
		log.Printf("unable to publish message: %s", err)
		return
	}
	log.Println("status arrived pickup notif sent")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// store in inventory
	db.AddShipmentToInventory(ctx, &shipment)

	// going back to the warehouse
	time.Sleep(estimateTime)
	msg.Ack(false)
}
