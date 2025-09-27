package entity

type Attachment struct {
	ID     int32  `json:"id"`
	TripID int32  `json:"tripId"`
	Name   string `json:"name"`
	Blob   []byte `json:"blob"`
}
