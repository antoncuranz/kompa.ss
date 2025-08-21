package converter

import (
	"kompass/internal/entity"
	"kompass/pkg/sqlc"
)

type ConvertFlightParams struct {
	Flight sqlc.Flight
	Legs   []sqlc.GetFlightLegsByFlightIDRow
	PNRs   []sqlc.Pnr
}

// goverter:converter
// goverter:extend ConvertFlightLeg
type FlightConverter interface {
	// goverter:map Flight.ID ID
	// goverter:map Flight.TripID TripID
	// goverter:map Flight.Price Price
	ConvertFlight(source ConvertFlightParams) entity.Flight

	ConvertFlightLegs(legs []sqlc.GetFlightLegsByFlightIDRow) []entity.FlightLeg

	ConvertPnrs(pnrs []sqlc.Pnr) []entity.PNR

	// goverter:map Pnr PNR
	ConvertPnr(pnr sqlc.Pnr) entity.PNR

	ConvertLocation(location sqlc.Location) entity.Location
}

func ConvertFlightLeg(c FlightConverter, leg sqlc.GetFlightLegsByFlightIDRow) entity.FlightLeg {
	return entity.FlightLeg{
		ID:                leg.FlightLeg.ID,
		Origin:            ConvertAirport(c, leg.Airport, leg.Location),
		Destination:       ConvertAirport(c, leg.Airport_2, leg.Location_2),
		Airline:           leg.FlightLeg.Airline,
		FlightNumber:      leg.FlightLeg.FlightNumber,
		DepartureDateTime: leg.FlightLeg.DepartureTime,
		ArrivalDateTime:   leg.FlightLeg.ArrivalTime,
		DurationInMinutes: leg.FlightLeg.DurationInMinutes,
		Aircraft:          leg.FlightLeg.Aircraft,
	}
}

func ConvertAirport(c FlightConverter, airport sqlc.Airport, location sqlc.Location) entity.Airport {
	mappedLocation := c.ConvertLocation(location)
	return entity.Airport{
		Iata:         airport.Iata,
		Name:         airport.Name,
		Municipality: airport.Municipality,
		Location:     &mappedLocation,
	}
}
