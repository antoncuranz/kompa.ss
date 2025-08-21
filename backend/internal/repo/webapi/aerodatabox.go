package webapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kompass/config"
	"kompass/internal/entity"
	"kompass/internal/repo/webapi/response"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/civil"
)

type AerodataboxWebAPI struct {
	baseURL string
	apiKey  string
}

func New(config config.WebApi) *AerodataboxWebAPI {
	return &AerodataboxWebAPI{
		baseURL: config.AerodataboxBaseURL,
		apiKey:  config.AerodataboxApiKey,
	}
}

func (a *AerodataboxWebAPI) RetrieveFlightLeg(ctx context.Context, date string, flightNumber string, origin *string) (entity.FlightLeg, error) {
	urlFormat := "%s/flights/number/%s/%s?dateLocalRole=Departure"
	url := fmt.Sprintf(urlFormat, a.baseURL, flightNumber, date)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Add("x-rapidapi-key", a.apiKey)
	req.Header.Add("x-rapidapi-host", "aerodatabox.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("do http request: %w", err)
	} else if res.StatusCode != 200 {
		return entity.FlightLeg{}, fmt.Errorf("http status code %d: %w", res.StatusCode, err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("read response body: %w", err)
	}

	var flightContracts []response.FlightContract
	if err := json.Unmarshal(body, &flightContracts); err != nil {
		return entity.FlightLeg{}, fmt.Errorf("unmarshall JSON: %w", err)
	}

	if len(flightContracts) == 0 {
		return entity.FlightLeg{}, fmt.Errorf("no flight contract found")
	}

	if origin == nil {
		return mapFlightLeg(flightContracts[0])
	}

	for _, flight := range flightContracts {
		if *flight.Departure.Airport.Iata == *origin {
			return mapFlightLeg(flight)
		}
	}

	return entity.FlightLeg{}, fmt.Errorf("no flight contract found with origin: %s", *origin)
}

func mapFlightLeg(flightContract response.FlightContract) (entity.FlightLeg, error) {
	if flightContract.Departure.ScheduledTime == nil || flightContract.Arrival.ScheduledTime == nil {
		return entity.FlightLeg{}, fmt.Errorf("flight leg does not have scheduled time")
	}

	departureDateTime, err1 := parseLocalDateTime(flightContract.Departure.ScheduledTime.Local)
	arrivalDateTime, err2 := parseLocalDateTime(flightContract.Arrival.ScheduledTime.Local)
	if err := errors.Join(err1, err2); err != nil {
		return entity.FlightLeg{}, fmt.Errorf("parse timestamps: %w", err)
	}

	durationInMinutes, err := calculateDuration(flightContract.Departure.ScheduledTime.UTC, flightContract.Arrival.ScheduledTime.UTC)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("calculate duration: %w", err)
	}

	return entity.FlightLeg{
		Origin:            mapAirport(flightContract.Departure.Airport),
		Destination:       mapAirport(flightContract.Arrival.Airport),
		Airline:           flightContract.Airline.Name, // FIXME
		FlightNumber:      flightContract.Number,
		DepartureDateTime: departureDateTime,
		ArrivalDateTime:   arrivalDateTime,
		DurationInMinutes: durationInMinutes,
		Aircraft:          flightContract.Aircraft.Model, // FIXME
	}, nil
}

func calculateDuration(fromUtcTimestamp string, toUtcTimestamp string) (int32, error) {
	from, err1 := parseTimestamp(fromUtcTimestamp)
	to, err2 := parseTimestamp(toUtcTimestamp)
	if err := errors.Join(err1, err2); err != nil {
		return 0, fmt.Errorf("parse timestamps: %w", err)
	}

	return int32(to.Sub(from).Minutes()), nil
}

func parseTimestamp(timestamp string) (time.Time, error) {
	format := "2006-01-02 15:04Z"
	return time.Parse(format, timestamp) // error already contains timestamp and format
}

func parseLocalDateTime(timestamp string) (civil.DateTime, error) {
	// Remove timezone and add missing seconds
	timestamp = timestamp[0:16] + ":00"
	// Separator between date and time needs to be a T
	timestamp = strings.Replace(timestamp, " ", "T", 1)
	return civil.ParseDateTime(timestamp) // error already contains timestamp and format
}

func mapAirport(airport response.ListingAirportContract) entity.Airport {
	return entity.Airport{
		Iata:         *airport.Iata, // FIXME
		Name:         airport.Name,
		Municipality: *airport.MunicipalityName, // FIXME
		Location: &entity.Location{
			Latitude:  airport.Location.Lat,
			Longitude: airport.Location.Lon,
		},
	}
}
