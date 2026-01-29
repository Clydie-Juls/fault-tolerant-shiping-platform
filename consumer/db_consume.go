package main

import (
	"context"
	"rabbitmq/internal/db"
	"rabbitmq/models"
	"rabbitmq/utils"
)

type ShipmentDB struct {
	db.DBConn
}

func NewShipmentDB(conn *db.DBConn) *ShipmentDB {
	return &ShipmentDB{
		DBConn: *conn,
	}
}

func (s *ShipmentDB) AddShipmentToInventory(ctx context.Context, shipment *models.Shipment) {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO inventory (product_name, price)
		VALUES (?, ?);`,
		shipment.ProductName,
		shipment.SellAmount,
	)
	utils.FailOnError(err, "unable to add shipment to inventory")
}
