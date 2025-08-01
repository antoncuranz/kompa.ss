// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"cloud.google.com/go/civil"
)

type Airport struct {
	Iata         string  `json:"iata"`
	Name         string  `json:"name"`
	Municipality string  `json:"municipality"`
	Location     *string `json:"location"`
}

type FlightLeg struct {
	ID                int32          `json:"id"`
	Origin            Airport        `json:"origin"`
	Destination       Airport        `json:"destination"`
	Airline           string         `json:"airline"`
	FlightNumber      string         `json:"flightNumber"`
	DepartureDateTime civil.DateTime `json:"departureDateTime"`
	ArrivalDateTime   civil.DateTime `json:"arrivalDateTime"`
	DurationInMinutes int32          `json:"durationInMinutes"`
	Aircraft          *string        `json:"aircraft"`
}

type PNR struct {
	ID      int32  `json:"id"`
	Airline string `json:"airline" example:"LH"`
	PNR     string `json:"pnr"     example:"123456"`
}

type Flight struct {
	ID     int32       `json:"id"`
	TripID int32       `json:"tripId"`
	Legs   []FlightLeg `json:"legs"`
	PNRs   []PNR       `json:"pnrs"`
	Price  *int32      `json:"price"`
}
