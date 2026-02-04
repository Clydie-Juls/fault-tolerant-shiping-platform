package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rabbitmq/models"
	"rabbitmq/utils"
	"time"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

type InventoryHandlerRepo interface {
	SendShipment(ctx context.Context, shipment *models.Shipment)
	SendOrder(ctx context.Context, order *models.Order)
}

type InventoryHandler struct {
	ch *amqp.Channel
}

func NewInventoryHandler(ch *amqp.Channel) *InventoryHandler {
	return &InventoryHandler{ch: ch}
}

func (ih *InventoryHandler) SendShipment(ctx context.Context, shipment *models.Shipment) {
	shipmentByte, err := json.Marshal(shipment)
	utils.FailOnError(err, "unable to convert shipment to byte")

	err = ih.ch.PublishWithContext(
		ctx,
		EXCHANGE_NAME,
		keyStatus(utils.STATUS_LISTED),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        shipmentByte,
		})

	utils.FailOnError(err, "unable to publish message")
	log.Printf("[x] Sent %s\n routing key: %s", string(shipmentByte), keyStatus(utils.STATUS_LISTED))
}

func (ih *InventoryHandler) SendOrder(ctx context.Context, order *models.Order) {
	orderByte, err := json.Marshal(order)
	utils.FailOnError(err, "unable to convert shipment to byte")

	err = ih.ch.PublishWithContext(
		ctx,
		EXCHANGE_NAME,
		keyStatus(utils.STATUS_BUY_REQUEST),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        orderByte,
		})

	utils.FailOnError(err, "unable to publish message")
	log.Printf("[x] Sent %s\n routing key: %s", string(orderByte), keyStatus(utils.STATUS_BUY_REQUEST))
}

type WsInventoryHandler struct {
	ih InventoryHandlerRepo
}

func NewWSHandler(ih InventoryHandlerRepo) *WsInventoryHandler {
	return &WsInventoryHandler{
		ih: ih,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ws *WsInventoryHandler) wsHandleMessage(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.FailOnError(err, "unable to upgrade http to websocket")

	for {
		shipment := models.Shipment{}
		err := conn.ReadJSON(&shipment)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws unexpected close: %v", err)
			} else {
				log.Printf("ws closed: %v", err)
			}
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		ws.ih.SendShipment(ctx, &shipment)
		cancel()
	}
}

func (ws *WsInventoryHandler) wsHandleOrder(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.FailOnError(err, "unable to upgrade http to websocket")

	for {
		order := models.Order{}
		err := conn.ReadJSON(&order)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws unexpected close: %v", err)
			} else {
				log.Printf("ws closed: %v", err)
			}
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		ws.ih.SendOrder(ctx, &order)
		cancel()
	}
}
