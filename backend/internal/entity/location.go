package entity

type Location struct {
	ID        int32   `json:"id"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
