// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
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
}
