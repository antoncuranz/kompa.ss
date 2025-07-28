// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"github.com/guregu/null/v6"
	"time"
)

type Airport struct {
	Iata         string      `json:"iata"`
	Name         string      `json:"name"`
	Municipality string      `json:"municipality"`
	Location     null.String `json:"location"`
}

type FlightLeg struct {
	ID            int32       `json:"id"`
	Origin        Airport     `json:"origin"`
	Destination   Airport     `json:"destination"`
	Airline       string      `json:"airline"`
	FlightNumber  string      `json:"flightNumber"`
	DepartureTime time.Time   `json:"departureTime"`
	ArrivalTime   time.Time   `json:"arrivalTime"`
	Aircraft      null.String `json:"aircraft"`
}

type PNR struct {
	ID      int32  `json:"id"`
	Airline string `json:"airline"`
	PNR     string `json:"pnr"`
}

type Flight struct {
	ID     int32       `json:"id"`
	TripID int32       `json:"tripId"`
	Price  null.Int32  `json:"price"`
	Legs   []FlightLeg `json:"legs"`
	PNRs   []PNR       `json:"pnrs"`
}
