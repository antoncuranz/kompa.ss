package entity

import "github.com/google/uuid"

type User struct {
	ID     int32     `json:"id"`
	Name   string    `json:"name"`
	JwtSub uuid.UUID `json:"jwtSub"`
}
