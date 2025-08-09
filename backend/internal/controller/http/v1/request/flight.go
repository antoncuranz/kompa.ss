package request

import "travel-planner/internal/entity"

type FlightLeg struct {
	Date          string  `json:"date"          validate:"required" example:"2026-01-30"`
	FlightNumber  string  `json:"flightNumber"  validate:"required" example:"EK412"`
	OriginAirport *string `json:"originAirport" example:"SYD"`
}

type Flight struct {
	Legs  []FlightLeg  `json:"legs"   validate:"required"`
	PNRs  []entity.PNR `json:"pnrs"   validate:"required"`
	Price *int32       `json:"price"`
}
