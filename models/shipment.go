package models

type person struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Shipment struct {
	ID          int       `json:"id"`
	Person      person    `json:"person"`
	ProductName string    `json:"product_name"`
	SellAmount  float64   `json:"sell_amount"`
	Location    Location  `json:"location"`
	Inventory   Inventory `json:"inventory"`
}

type Order struct {
	ID          int      `json:"id"`
	InventoryID int      `json:"inventory_id"`
	Person      person   `json:"person"`
	BuyAmount   float64  `json:"sell_amount"`
	Location    Location `json:"location"`
}
