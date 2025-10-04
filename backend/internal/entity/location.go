package entity

type Location struct {
	ID        int32   `json:"id"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type GeocodeLocation struct {
	Label     string  `json:"label"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
