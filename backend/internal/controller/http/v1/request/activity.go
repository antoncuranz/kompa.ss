package request

import (
	"cloud.google.com/go/civil"
	"kompass/internal/entity"
)

type Activity struct {
	Name        string           `json:"name"        validate:"required"`
	Date        civil.Date       `json:"date"        validate:"required"`
	Time        *civil.Time      `json:"time"`
	Description *string          `json:"description"`
	Address     *string          `json:"address"`
	Location    *entity.Location `json:"location"`
	Price       *int32           `json:"price"`
}
