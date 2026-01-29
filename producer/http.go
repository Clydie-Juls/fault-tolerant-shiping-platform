package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rabbitmq/utils"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func sendMessage(ch *amqp.Channel, w http.ResponseWriter, r *http.Request) {
	shipmentByte, err := json.Marshal(r.Body)
	utils.FailOnError(err, "unable to send message")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		EXCHANGE_NAME,
		utils.SeverityFrom(os.Args),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(string(shipmentByte)),
		})

	utils.FailOnError(err, "unable to publish message")
	log.Printf("[x] Sent %s", string(shipmentByte))
	w.Write([]byte("message sent"))
}

func createHTTPServer(ch *amqp.Channel, connections ...connection) {
	mux := http.NewServeMux()
	mux.HandleFunc("/shipment/sell", func(w http.ResponseWriter, r *http.Request) {
		sendMessage(ch, w, r)
	})
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	serverError := make(chan error, 1)
	go func() {
		log.Println("starting server at: ", addr)
		serverError <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		utils.FailOnError(err, "something's wrong with the api server")
	case signal := <-shutdown:
		log.Printf("Shutting down: %v", signal)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			utils.FailOnError(err, "unable to shutdown gracefully")
		}

		for _, conn := range connections {
			conn.Close()
		}
	}
}
