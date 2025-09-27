package entity

import (
	"cloud.google.com/go/civil"
)

type Trip struct {
	ID          int32      `json:"id"`
	OwnerID     int32      `json:"owner_id"`
	Name        string     `json:"name"`
	StartDate   civil.Date `json:"startDate"`
	EndDate     civil.Date `json:"endDate"`
	Description *string    `json:"description" extensions:"nullable"`
	ImageUrl    *string    `json:"imageUrl" extensions:"nullable"`
}
