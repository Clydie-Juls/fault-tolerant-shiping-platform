package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	bindInfoTypeQueues(ch, q, utils.STATUS_LISTED)

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

	for msg := range msgs {
		go processDelivery(ch, &msg, *warehouse, db)
	}
	<-forever
}

func processDelivery(ch *amqp.Channel, msg *amqp.Delivery, warehouse models.Warehouse, db *ShipmentDB) {
	var shipment models.Shipment
	json.Unmarshal(msg.Body, &shipment)
	estimateTime := warehouse.EstimateTimeToDestination(shipment.Location)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ch.PublishWithContext(ctx,
		utils.EXCHANGE_NAME,
		getTopicRoutingKey(utils.STATUS_BUY_ACCEPTED),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("reaching destination in about %.2f seconds", estimateTime)),
		},
	)

	// reaching destination
	time.Sleep(estimateTime)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ch.PublishWithContext(ctx,

		utils.EXCHANGE_NAME,
		getTopicRoutingKey(utils.STATUS_ARRIVED_PICKUP),
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("destination reached"),
		},
	)

	// going back to the warehouse
	time.Sleep(estimateTime)
	msg.Ack(false)
}
