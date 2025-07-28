// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"cloud.google.com/go/civil"
	"github.com/guregu/null/v6"
)

type Activity struct {
	ID          int32       `json:"id"`
	TripID      int32       `json:"tripId"`
	Name        string      `json:"name"`
	Description null.String `json:"description"`
	Date        civil.Date  `json:"date"`
}
