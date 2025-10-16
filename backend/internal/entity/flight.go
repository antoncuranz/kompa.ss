package entity

import (
	"cloud.google.com/go/civil"
	"fmt"
)

type FlightDetail struct {
	Legs []FlightLeg `json:"legs"`
	PNRs []PNR       `json:"pnrs"`
}

type Airport struct {
	Iata         string   `json:"iata"`
	Name         string   `json:"name"`
	Municipality string   `json:"municipality"`
	Location     Location `json:"location"`
}

type AirportWithTimezone struct {
	Airport
	Timezone string `json:"timezone"`
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
	Aircraft          *string        `json:"aircraft" extensions:"nullable"`
}

type PNR struct {
	ID      int32  `json:"id"`
	Airline string `json:"airline" example:"LH"`
	PNR     string `json:"pnr"     example:"123456"`
}

type ErrAmbiguousFlightRequest map[string][]AmbiguousFlightChoice

type AmbiguousFlightChoice struct {
	OriginIata        string         `json:"originIata"`
	DestinationIata   string         `json:"destinationIata"`
	DepartureDateTime civil.DateTime `json:"departureDateTime"`
}

func (e ErrAmbiguousFlightRequest) Error() string {
	return fmt.Sprint(e)
}
