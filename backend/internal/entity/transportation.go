// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"cloud.google.com/go/civil"
)

type TransportationType string

const (
	PLANE TransportationType = "PLANE"
	TRAIN TransportationType = "TRAIN"
	BUS   TransportationType = "BUS"
	BOAT  TransportationType = "BOAT"
	BIKE  TransportationType = "BIKE"
	CAR   TransportationType = "CAR"
	FOOT  TransportationType = "FOOT"
	OTHER TransportationType = "OTHER"
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
}
