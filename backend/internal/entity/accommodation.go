// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"cloud.google.com/go/civil"
)

type Accommodation struct {
	ID            int32       `json:"id"`
	TripID        int32       `json:"tripId"`
	Name          string      `json:"name"`
	ArrivalDate   civil.Date  `json:"arrivalDate"`
	DepartureDate civil.Date  `json:"departureDate"`
	CheckInTime   *civil.Time `json:"checkInTime"`
	CheckOutTime  *civil.Time `json:"checkOutTime"`
	Description   *string     `json:"description"`
	Location      *Location   `json:"location"`
	Price         *int32      `json:"price"`
}
