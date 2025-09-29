package request

import (
	"cloud.google.com/go/civil"
	"kompass/internal/entity"
)

type FlightLeg struct {
	Date          civil.Date `json:"date"          example:"2026-01-30"`
	FlightNumber  string     `json:"flightNumber"  example:"EK412"`
	OriginAirport *string    `json:"originAirport" extensions:"nullable" example:"SYD"`
}

type Flight struct {
	Legs  []FlightLeg  `json:"legs"`
	PNRs  []entity.PNR `json:"pnrs"`
	Price *int32       `json:"price" extensions:"nullable"`
}
