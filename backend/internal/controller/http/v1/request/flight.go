package request

type FlightLeg struct {
	Date          string `json:"date"          validate:"required" example:"2026-01-30"`
	FlightNumber  string `json:"flightNumber"  validate:"required" example:"EK412"`
	OriginAirport string `json:"originAirport" example:"SYD"`
}

type PNR struct {
	Airline string `json:"pnr"     validate:"required" example:"LH"`
	PNR     string `json:"airline" validate:"required" example:"123456"`
}

type Flight struct {
	TripID int32       `json:"tripId" validate:"required" example:"1"`
	Legs   []FlightLeg `json:"legs"   validate:"required"`
	PNRs   []PNR       `json:"pnrs"   validate:"required"`
	Price  int32       `json:"price"`
}
