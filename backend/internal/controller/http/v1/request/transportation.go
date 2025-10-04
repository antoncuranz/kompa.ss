package request

import (
	"cloud.google.com/go/civil"
	"kompass/internal/entity"
)

type Transportation struct {
	Name               string                    `json:"name"`
	TripID             int32                     `json:"tripId"`
	Type               entity.TransportationType `json:"type"`
	Origin             entity.Location           `json:"origin"`
	Destination        entity.Location           `json:"destination"`
	DepartureDateTime  civil.DateTime            `json:"departureDateTime"`
	ArrivalDateTime    civil.DateTime            `json:"arrivalDateTime"`
	OriginAddress      *string                   `json:"originAddress" extensions:"nullable"`
	DestinationAddress *string                   `json:"destinationAddress" extensions:"nullable"`
	Price              *int32                    `json:"price" extensions:"nullable"`
}
