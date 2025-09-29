package converter

import (
	"cloud.google.com/go/civil"
	"kompass/internal/entity"
	"kompass/internal/repo/webapi/response"
	"strings"
)

// goverter:converter
// goverter:extend ParseTimestamp
type TrainConverter interface {
	ConvertJourney(source response.Journey) (entity.TrainDetail, error)
	// goverter:ignore ID
	// goverter:map PlannedDeparture DepartureDateTime
	// goverter:map PlannedArrival ArrivalDateTime
	// goverter:map Line.Name LineName
	// goverter:useZeroValueOnPointerInconsistency
	// TODO!
	// goverter:ignore DurationInMinutes
	ConvertLeg(source response.Leg) (entity.TrainLeg, error)

	ConvertStation(source response.StationOrStop) entity.TrainStation

	// goverter:ignore ID
	ConvertLocation(source response.Location) entity.Location
}

func ParseTimestamp(timestamp string) (civil.DateTime, error) {
	parts := strings.Split(timestamp, "+")
	return civil.ParseDateTime(parts[0])
}
