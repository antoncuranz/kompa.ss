// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"cloud.google.com/go/civil"
)

type TransportationType int32

const (
	PLANE TransportationType = iota
	TRAIN
	BUS
	BOAT
	BIKE
	CAR
	FOOT
	OTHER
)

type Transportation struct {
	ID                int32              `json:"id"`
	TripID            int32              `json:"tripId"`
	Type              TransportationType `json:"type"`
	Origin            Location           `json:"origin"`
	Destination       Location           `json:"destination"`
	DepartureDateTime civil.DateTime     `json:"departureDateTime"`
	ArrivalDateTime   civil.DateTime     `json:"arrivalDateTime"`
	GeoJson           *string            `json:"geoJson"`
	Price             *int32             `json:"price"`
}
