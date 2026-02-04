package main

import (
	"context"
	"log"
	"rabbitmq/internal/db"
	"rabbitmq/models"
	"rabbitmq/utils"
)

type ShipmentDB struct {
	db.DBConn
	warehouseID int
}

func NewShipmentDB(conn *db.DBConn, warehouseID int) *ShipmentDB {
	return &ShipmentDB{
		DBConn:      *conn,
		warehouseID: warehouseID,
	}
}

func (s *ShipmentDB) AddShipmentToInventory(ctx context.Context, shipment *models.Shipment) {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO inventory (product_name, price, warehouse_id)
		VALUES ($1, $2, $3);`,
		shipment.ProductName,
		shipment.SellAmount,
		s.warehouseID,
	)
	utils.FailOnError(err, "unable to add shipment to inventory")
}

func (s *ShipmentDB) ConsumeInventoryByID(
	ctx context.Context,
	id int64,
) *models.Inventory {
	const q = `
		DELETE FROM inventory
		WHERE id = $1
		RETURNING
			id,
			product_name,
			price;
	`

	var inv models.Inventory

	err := s.DB.QueryRowContext(ctx, q, id).Scan(
		&inv.Id,
		&inv.ProductName,
		&inv.Price,
	)

	if err == nil {
		return &inv
	}
	log.Println(err)
	return nil
}
