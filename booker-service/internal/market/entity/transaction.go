package entity

import "time"

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Price        float64
	Total        float64
	DateTime     time.Time
	Shares       int
}
