package request

import "kompass/internal/entity"

type FlightLeg struct {
	Date          string  `json:"date"          example:"2026-01-30"`
	FlightNumber  string  `json:"flightNumber"  example:"EK412"`
	OriginAirport *string `json:"originAirport" example:"SYD"`
}

type Flight struct {
	Legs  []FlightLeg  `json:"legs"`
	PNRs  []entity.PNR `json:"pnrs"`
	Price *int32       `json:"price"`
}
