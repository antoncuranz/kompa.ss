package entity

import (
	"cloud.google.com/go/civil"
)

type TransportationType string

const (
	FLIGHT TransportationType = "FLIGHT"
	TRAIN  TransportationType = "TRAIN"
	BUS    TransportationType = "BUS"
	CAR    TransportationType = "CAR"
	FERRY  TransportationType = "FERRY"
	BOAT   TransportationType = "BOAT"
	BIKE   TransportationType = "BIKE"
	HIKE   TransportationType = "HIKE"
	OTHER  TransportationType = "OTHER"
)

func (t TransportationType) String() string {
	return string(t)
}

type Transportation struct {
	ID                int32              `json:"id"`
	TripID            int32              `json:"tripId"`
	Type              TransportationType `json:"type"`
	Origin            Location           `json:"origin"`
	Destination       Location           `json:"destination"`
	DepartureDateTime civil.DateTime     `json:"departureDateTime"`
	ArrivalDateTime   civil.DateTime     `json:"arrivalDateTime"`
	Price             *int32             `json:"price" extensions:"nullable"`
	FlightDetail      *FlightDetail      `json:"flightDetail,omitempty" validate:"optional" extensions:"nullable"`
	TrainDetail       *TrainDetail       `json:"trainDetail,omitempty" validate:"optional" extensions:"nullable"`
	//GenericDetail     *GenericDetail     `json:"genericDetail,omitempty" validate:"optional" extensions:"nullable"`
}

type GenericDetail struct {
	Name               string `json:"name"`
	OriginAddress      string `json:"originAddress"`
	DestinationAddress string `json:"destinationAddress"`
}
