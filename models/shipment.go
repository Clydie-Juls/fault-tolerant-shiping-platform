package models

type person struct {
	ID   int64
	name string
	Age  int
}

type Location struct {
	Longitude float64
	Latitude  float64
}

type Shipment struct {
	ID          int
	Person      person
	ProductName string
	SellAmount  float64
	Location    Location
	Inventory   Inventory
}
