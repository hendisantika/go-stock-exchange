package entity

import "sync"

type Book struct {
	Order         []*Order
	Transactions  []*Transaction
	OrdersChan    chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}
