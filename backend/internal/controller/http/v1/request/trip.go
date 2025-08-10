package request

import (
	"cloud.google.com/go/civil"
)

type Trip struct {
	Name        string     `json:"name"         validate:"required"`
	StartDate   civil.Date `json:"startDate"    validate:"required"`
	EndDate     civil.Date `json:"endDate"      validate:"required"`
	Description *string    `json:"description"`
	ImageUrl    *string    `json:"imageUrl"`
}
