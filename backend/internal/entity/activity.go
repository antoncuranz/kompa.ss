package entity

import (
	"cloud.google.com/go/civil"
)

type Activity struct {
	ID          int32       `json:"id"`
	TripID      int32       `json:"tripId"`
	Name        string      `json:"name"`
	Date        civil.Date  `json:"date"`
	Time        *civil.Time `json:"time" extensions:"nullable"`
	Description *string     `json:"description" extensions:"nullable"`
	Address     *string     `json:"address" extensions:"nullable"`
	Location    *Location   `json:"location" extensions:"nullable"`
	Price       *int32      `json:"price" extensions:"nullable"`
}
