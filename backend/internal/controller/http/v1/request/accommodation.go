package request

import (
	"kompass/internal/entity"

	"cloud.google.com/go/civil"
)

type Accommodation struct {
	Name          string           `json:"name"`
	ArrivalDate   civil.Date       `json:"arrivalDate"`
	DepartureDate civil.Date       `json:"departureDate"`
	CheckInTime   *civil.Time      `json:"checkInTime"`
	CheckOutTime  *civil.Time      `json:"checkOutTime"`
	Description   *string          `json:"description"`
	Address       *string          `json:"address"`
	Location      *entity.Location `json:"location"`
	Price         *int32           `json:"price"`
}
