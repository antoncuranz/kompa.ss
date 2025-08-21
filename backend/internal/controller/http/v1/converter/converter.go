package converter

import (
	"kompass/internal/controller/http/v1/response"
	"kompass/internal/entity"
)

// goverter:converter
type TripConverter interface {
	ConvertTrips(trips []entity.Trip) []response.Trip
	ConvertTrip(trip entity.Trip) response.Trip
}
