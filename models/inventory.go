package models

type Inventory struct {
	Id          int     `json:"id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
}
