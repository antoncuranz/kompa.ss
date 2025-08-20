package request

import (
	"kompass/internal/entity"

	"cloud.google.com/go/civil"
)

type Activity struct {
	Name        string           `json:"name"`
	Date        civil.Date       `json:"date"`
	Time        *civil.Time      `json:"time"`
	Description *string          `json:"description"`
	Address     *string          `json:"address"`
	Location    *entity.Location `json:"location"`
	Price       *int32           `json:"price"`
}
