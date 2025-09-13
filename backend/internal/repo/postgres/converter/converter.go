package converter

import (
	"kompass/internal/entity"
	"kompass/pkg/sqlc"
)

type ConvertTransportationParams struct {
	Transportation sqlc.Transportation
	Origin         sqlc.Location
	Destination    sqlc.Location
}

// goverter:converter
type TransportationConverter interface {
	// goverter:map Transportation.ID ID
	// goverter:map Transportation.TripID TripID
	// goverter:map Transportation.Type Type
	// goverter:map Transportation.DepartureTime DepartureDateTime
	// goverter:map Transportation.ArrivalTime ArrivalDateTime
	// goverter:map Transportation.Price Price
	// goverter:ignore FlightDetail TrainDetail
	// goverter:ignore GeoJson
	ConvertTransportation(source ConvertTransportationParams) entity.Transportation
}

// goverter:converter
// goverter:extend ConvertFlightLeg
type FlightConverter interface {
	ConvertFlightLegs(legs []sqlc.GetFlightLegsByTransportationIDRow) []entity.FlightLeg

	ConvertPnrs(pnrs []sqlc.FlightPnr) []entity.PNR

	// goverter:map Pnr PNR
	ConvertPnr(pnr sqlc.FlightPnr) entity.PNR

	ConvertLocation(location sqlc.Location) entity.Location
}

func ConvertFlightLeg(c FlightConverter, leg sqlc.GetFlightLegsByTransportationIDRow) entity.FlightLeg {
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
	return entity.Airport{
		Iata:         airport.Iata,
		Name:         airport.Name,
		Municipality: airport.Municipality,
		Location:     c.ConvertLocation(location),
	}
}
