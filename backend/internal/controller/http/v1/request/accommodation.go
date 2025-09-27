package request

import (
	"kompass/internal/entity"

	"cloud.google.com/go/civil"
)

type Accommodation struct {
	Name          string           `json:"name"`
	ArrivalDate   civil.Date       `json:"arrivalDate"`
	DepartureDate civil.Date       `json:"departureDate"`
	CheckInTime   *civil.Time      `json:"checkInTime" extensions:"nullable"`
	CheckOutTime  *civil.Time      `json:"checkOutTime" extensions:"nullable"`
	Description   *string          `json:"description" extensions:"nullable"`
	Address       *string          `json:"address" extensions:"nullable"`
	Location      *entity.Location `json:"location" extensions:"nullable"`
	Price         *int32           `json:"price" extensions:"nullable"`
}
