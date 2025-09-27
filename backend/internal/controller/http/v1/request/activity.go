package request

import (
	"kompass/internal/entity"

	"cloud.google.com/go/civil"
)

type Activity struct {
	Name        string           `json:"name"`
	Date        civil.Date       `json:"date"`
	Time        *civil.Time      `json:"time" extensions:"nullable"`
	Description *string          `json:"description" extensions:"nullable"`
	Address     *string          `json:"address" extensions:"nullable"`
	Location    *entity.Location `json:"location" extensions:"nullable"`
	Price       *int32           `json:"price" extensions:"nullable"`
}
