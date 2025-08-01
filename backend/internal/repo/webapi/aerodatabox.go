package webapi

import (
	"cloud.google.com/go/civil"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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

func (a *AerodataboxWebAPI) RetrieveFlightLeg(ctx context.Context, date string, flightNumber string, origin *string) (entity.FlightLeg, error) {
	urlFormat := "https://aerodatabox.p.rapidapi.com/flights/number/%s/%s?dateLocalRole=Departure"
	url := fmt.Sprintf(urlFormat, flightNumber, date)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return entity.FlightLeg{}, err
	}

	req.Header.Add("x-rapidapi-key", a.apiKey)
	req.Header.Add("x-rapidapi-host", "aerodatabox.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.FlightLeg{}, err
	} else if res.StatusCode < 200 || res.StatusCode >= 300 {
		return entity.FlightLeg{}, errors.New("Error calling Aerodatabox API, status: " + res.Status)
	} else if res.StatusCode != 200 {
		return entity.FlightLeg{}, errors.New("AerodataboxWebAPI - RetrieveFlightLeg - no results")
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return entity.FlightLeg{}, err
	}

	var data []response.FlightContract
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return entity.FlightLeg{}, err
	}

	if len(data) == 0 {
		return entity.FlightLeg{}, errors.New("AerodataboxWebAPI - RetrieveFlightLeg - no results")
	}

	if origin == nil {
		return mapFlightLeg(data[0])
	}

	for _, flight := range data {
		if *flight.Departure.Airport.Iata == *origin {
			return mapFlightLeg(flight)
		}
	}

	return entity.FlightLeg{}, errors.New("AerodataboxWebAPI - RetrieveFlightLeg - no results")
}

func mapFlightLeg(flightContract response.FlightContract) (entity.FlightLeg, error) {
	departureDateTime, err1 := parseLocalDateTime(flightContract.Departure.ScheduledTime.Local)
	arrivalDateTime, err2 := parseLocalDateTime(flightContract.Arrival.ScheduledTime.Local)
	if err := errors.Join(err1, err2); err != nil {
		fmt.Printf("Error parsing timestamp: departure=%s, arrival=%s\n", departureDateTime, arrivalDateTime)
		return entity.FlightLeg{}, err
	}

	durationInMinutes, err := calculateDuration(flightContract.Departure.ScheduledTime.UTC, flightContract.Arrival.ScheduledTime.UTC)
	if err != nil {
		fmt.Println("Error Calculating duration of flight leg.")
		return entity.FlightLeg{}, err
	}

	return entity.FlightLeg{
		Origin:            mapAirport(flightContract.Departure.Airport),
		Destination:       mapAirport(flightContract.Arrival.Airport),
		Airline:           flightContract.Airline.Name,
		FlightNumber:      flightContract.Number,
		DepartureDateTime: departureDateTime,
		ArrivalDateTime:   arrivalDateTime,
		DurationInMinutes: durationInMinutes,
		Aircraft:          flightContract.Aircraft.Model,
	}, nil
}

func calculateDuration(fromUtcTimestamp string, toUtcTimestamp string) (int32, error) {
	from, err1 := parseTimestamp(fromUtcTimestamp)
	to, err2 := parseTimestamp(toUtcTimestamp)
	if err := errors.Join(err1, err2); err != nil {
		fmt.Printf("Error parsing timestamp: from=%s, to=%s\n", from, to)
		return 0, err
	}

	return int32(to.Sub(from).Minutes()), nil
}

func parseTimestamp(timestamp string) (time.Time, error) {
	format := "2006-01-02 15:04Z"
	return time.Parse(format, timestamp)
}

func parseLocalDateTime(timestamp string) (civil.DateTime, error) {
	// Remove timezone and add missing seconds
	timestamp = timestamp[0:16] + ":00"
	// Separator between date and time needs to be a T
	timestamp = strings.Replace(timestamp, " ", "T", 1)
	return civil.ParseDateTime(timestamp)
}

func mapAirport(airport response.ListingAirportContract) entity.Airport {
	return entity.Airport{
		Iata:         *airport.Iata,
		Name:         airport.Name,
		Municipality: *airport.MunicipalityName,
	}
}
