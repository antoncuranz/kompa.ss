package response

import (
	"kompass/internal/entity"

	"cloud.google.com/go/civil"
)

type Transportation struct {
	ID                int32                `json:"id"`
	TripID            int32                `json:"tripId"`
	Type              string               `json:"type"`
	Origin            entity.Location      `json:"origin"`
	Destination       entity.Location      `json:"destination"`
	DepartureDateTime civil.DateTime       `json:"departureDateTime"`
	ArrivalDateTime   civil.DateTime       `json:"arrivalDateTime"`
	GeoJson           *string              `json:"geoJson" extensions:"nullable"`
	Price             *int32               `json:"price" extensions:"nullable"`
	FlightDetail      *entity.FlightDetail `json:"flightDetail,omitempty" validate:"optional" extensions:"nullable"`
	TrainDetail       *entity.TrainDetail  `json:"trainDetail,omitempty" validate:"optional" extensions:"nullable"`
}
