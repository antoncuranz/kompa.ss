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
	Address       *string     `json:"address"`
	Location      *Location   `json:"location"`
	Price         *int32      `json:"price"`
}
