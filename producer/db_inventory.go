package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"rabbitmq/models"
	"rabbitmq/utils"
	"time"
)

type InventoryDB struct {
	db *sql.DB
}

func NewInventoryDB(db *sql.DB) *InventoryDB {
	return &InventoryDB{
		db: db,
	}
}

func (i *InventoryDB) GetAllInventory(w http.ResponseWriter, r *http.Request) {
	warehouseID := i.GetWarehouseID(warehouseState, warehouseCity)
	inventory := i.GetInventory(warehouseID)
	err := json.NewEncoder(w).Encode(inventory)
	utils.FailOnError(err, "unable to encode inventory")
}

func (i *InventoryDB) GetWarehouseID(warehouseState, warehouseCity string) int {
	const q = `
		SELECT id
		FROM warehouse
		WHERE warehouse_state = $1
		  AND warehouse_city  = $2
		LIMIT 1;
	`

	var id int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.db.QueryRowContext(ctx, q, warehouseState, warehouseCity).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.FailOnError(err, "warehouse does not exist")
			return -1
		}
		utils.FailOnError(err, "unable to get warehouse id")
		return -1
	}

	return id
}

func (i *InventoryDB) GetInventory(warehouse_id int) []models.Inventory {
	inventory := []models.Inventory{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := i.db.QueryContext(ctx, "SELECT id, product_name, price FROM inventory WHERE warehouse_id = $1;", warehouse_id)
	utils.FailOnError(err, "unable to query inventory")

	for rows.Next() {
		inv := models.Inventory{}
		if err := rows.Scan(&inv.Id, &inv.ProductName, &inv.Price); err != nil {
			utils.FailOnError(err, "unable to scan row")
		}

		inventory = append(inventory, inv)
	}

	log.Printf("no of inventory items: %d", len(inventory))
	return inventory
}
