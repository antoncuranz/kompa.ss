package request

import (
	"cloud.google.com/go/civil"
	"travel-planner/internal/entity"
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
