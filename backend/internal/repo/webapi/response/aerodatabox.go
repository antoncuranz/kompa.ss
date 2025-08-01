package response

type GeoCoordinatesContract struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type ListingAirportContract struct {
	Name             string                  `json:"name"`
	Iata             *string                 `json:"iata"`
	MunicipalityName *string                 `json:"municipalityName"`
	Location         *GeoCoordinatesContract `json:"location"`
}

type DateTimeContract struct {
	Local string `json:"local"`
	UTC   string `json:"utc"`
}

type FlightAirportMovementContract struct {
	Airport       ListingAirportContract `json:"airport"`
	ScheduledTime *DateTimeContract      `json:"scheduledTime"`
}

type FlightAircraftContract struct {
	Model *string `json:"model"`
}

type FlightAirlineContract struct {
	Name string `json:"name"`
}

type FlightContract struct {
	Number    string                        `json:"number"`
	Departure FlightAirportMovementContract `json:"departure"`
	Arrival   FlightAirportMovementContract `json:"arrival"`
	Aircraft  *FlightAircraftContract       `json:"aircraft"`
	Airline   *FlightAirlineContract        `json:"airline"`
}
