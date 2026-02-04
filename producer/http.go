package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rabbitmq/internal/db"
	"rabbitmq/utils"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func createHTTPServer(ch *amqp.Channel, connections ...connection) {
	mux := http.NewServeMux()
	ws := NewWSHandler(NewInventoryHandler(ch))
	mux.HandleFunc("/shipment/sell", ws.wsHandleMessage)
	mux.HandleFunc("/shipment/buy", ws.wsHandleOrder)
	dbConn := db.NewDbConn()
	invDb := NewInventoryDB(dbConn.DB)
	mux.HandleFunc("GET /inventory", invDb.GetAllInventory)
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
			utils.FailOnError(err, "unable to shutdown gracefully")
		}

		for _, conn := range connections {
			conn.Close()
		}
	}
}
