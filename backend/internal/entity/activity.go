package entity

import (
	"cloud.google.com/go/civil"
)

type Activity struct {
	ID          int32       `json:"id"`
	TripID      int32       `json:"tripId"`
	Name        string      `json:"name"`
	Date        civil.Date  `json:"date"`
	Time        *civil.Time `json:"time"`
	Description *string     `json:"description"`
	Address     *string     `json:"address"`
	Location    *Location   `json:"location"`
	Price       *int32      `json:"price"`
}
