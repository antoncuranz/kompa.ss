package response

import (
	"kompass/internal/entity"

	"cloud.google.com/go/civil"
)

type Transportation struct {
	ID                   int32                     `json:"id"`
	TripID               int32                     `json:"tripId"`
	Type                 entity.TransportationType `json:"type"`
	Origin               entity.Location           `json:"origin"`
	Destination          entity.Location           `json:"destination"`
	DepartureDateTime    civil.DateTime            `json:"departureDateTime"`
	ArrivalDateTime      civil.DateTime            `json:"arrivalDateTime"`
	GeoJson              *string                   `json:"geoJson" extensions:"nullable"`
	Price                *int32                    `json:"price" extensions:"nullable"`
	TransportationDetail any                       `json:"transportationDetail" oneOf:"FlightDetail,TrainDetail" extensions:"nullable"`
}

type FlightDetail struct {
	Legs []entity.FlightLeg `json:"legs"`
	PNRs []entity.PNR       `json:"pnrs"`
}

type TrainDetail struct {
	TrainNumber string
}
