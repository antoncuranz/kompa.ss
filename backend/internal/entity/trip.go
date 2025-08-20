// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import (
	"cloud.google.com/go/civil"
)

type Trip struct {
	ID          int32
	Name        string
	StartDate   civil.Date
	EndDate     civil.Date
	Description *string
	ImageUrl    *string
}
