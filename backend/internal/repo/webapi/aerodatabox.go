package webapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guregu/null/v6"
	"io"
	"net/http"
	"travel-planner/internal/entity"
	"travel-planner/internal/repo/webapi/response"
)

type AerodataboxWebAPI struct {
	apiKey string
}

func New(apiKey string) *AerodataboxWebAPI {
	return &AerodataboxWebAPI{
		apiKey: apiKey,
	}
}

func (a *AerodataboxWebAPI) RetrieveFlightLeg(ctx context.Context, date string, flightNumber string, origin string) (entity.FlightLeg, error) {
	urlFormat := "https://aerodatabox.p.rapidapi.com/flights/number/%s/%s?dateLocalRole=Departure"
	url := fmt.Sprintf(urlFormat, flightNumber, date)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Add("x-rapidapi-key", a.apiKey)
	req.Header.Add("x-rapidapi-host", "aerodatabox.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var data []response.FlightContract
	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return entity.FlightLeg{}, err
	}

	if len(data) == 0 {
		return entity.FlightLeg{}, errors.New("AerodataboxWebAPI - RetrieveFlightLeg - no results")
	}

	// TODO: if len(data) > 1 search for flight with same origin (if available)
	return mapFlightLeg(data[0]), nil
}

func mapFlightLeg(flightContract response.FlightContract) entity.FlightLeg {
	//format := "2006-01-02 15:04-07:00"
	//departureTime, err := time.Parse(format, flightContract.Departure.ScheduledTime.Local)
	//if err != nil {
	//	fmt.Printf("Error parsing timestamp: %s\n", flightContract.Departure.ScheduledTime.Local)
	//}
	//arrivalTime, err := time.Parse(format, flightContract.Arrival.ScheduledTime.Local)
	//if err != nil {
	//	fmt.Printf("Error parsing timestamp: %s\n", flightContract.Arrival.ScheduledTime.Local)
	//}

	return entity.FlightLeg{
		Origin:        mapAirport(flightContract.Departure.Airport),
		Destination:   mapAirport(flightContract.Arrival.Airport),
		Airline:       flightContract.Airline.Name,
		FlightNumber:  flightContract.Number,
		DepartureTime: flightContract.Departure.ScheduledTime.Local,
		ArrivalTime:   flightContract.Arrival.ScheduledTime.Local,
		Aircraft:      null.StringFrom(flightContract.Aircraft.Model),
	}
}

func mapAirport(airport response.ListingAirportContract) entity.Airport {
	return entity.Airport{
		Iata:         airport.Iata,
		Name:         airport.Name,
		Municipality: airport.MunicipalityName,
	}
}
