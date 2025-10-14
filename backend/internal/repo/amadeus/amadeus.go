package amadeus

import (
	"cloud.google.com/go/civil"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	goiso8601duration "github.com/xnacly/go-iso8601-duration"
	"io"
	"kompass/config"
	"kompass/internal/entity"
	"kompass/internal/repo"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AmadeusWebAPI struct {
	baseURL    string
	apiKey     string
	apiSecret  string
	iataLookup repo.IataLookup
}

func New(config config.WebApi, iataLookup repo.IataLookup) *AmadeusWebAPI {
	return &AmadeusWebAPI{
		baseURL:    config.AmadeusBaseURL,
		apiKey:     config.AmadeusApiKey,
		apiSecret:  config.AmadeusApiSecret,
		iataLookup: iataLookup,
	}
}

func (a *AmadeusWebAPI) RetrieveFlightLeg(ctx context.Context, date civil.Date, flightNumber string, requestedOrigin *string) (entity.FlightLeg, error) {
	urlFormat := "%s/v2/schedule/flights?carrierCode=%s&flightNumber=%s&scheduledDepartureDate=%s"
	scheduleUrl := fmt.Sprintf(urlFormat, a.baseURL, flightNumber[:2], strings.TrimSpace(flightNumber[2:]), date.String())

	req, err := http.NewRequestWithContext(ctx, "GET", scheduleUrl, nil)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("create http request: %w", err)
	}

	accessToken, err := a.getAccessToken(ctx)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("get access token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)

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

	var flightStatusResponse FlightStatusResponse
	if err := json.Unmarshal(body, &flightStatusResponse); err != nil {
		return entity.FlightLeg{}, fmt.Errorf("unmarshall JSON: %w", err)
	}

	if len(flightStatusResponse.Data) != 1 || len(flightStatusResponse.Data[0].Legs) > 2 {
		return entity.FlightLeg{}, fmt.Errorf("not found or too complex")
	}

	return a.mapDatedFlight(flightStatusResponse.Data[0], requestedOrigin)
}

func (a *AmadeusWebAPI) mapDatedFlight(datedFlight DatedFlight, requestedOrigin *string) (entity.FlightLeg, error) {
	leg, err := findLeg(datedFlight.Legs, requestedOrigin)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("find leg: %w", err)
	}

	originAirport, err1 := a.iataLookup.LookupAirport(leg.BoardPointIataCode)
	destinationAirport, err2 := a.iataLookup.LookupAirport(leg.OffPointIataCode)
	if err := errors.Join(err1, err2); err != nil {
		return entity.FlightLeg{}, fmt.Errorf("lookup airports: %w", err)
	}

	duration, err := goiso8601duration.From(leg.ScheduledLegDuration)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("parse duration: %w", err)
	}

	origin, originFound := findFlightPointByIata(datedFlight, leg.BoardPointIataCode)
	destination, destinationFound := findFlightPointByIata(datedFlight, leg.OffPointIataCode)
	if !originFound && !destinationFound {
		return entity.FlightLeg{}, fmt.Errorf("no flight point found")
	}

	originTz, err1 := time.LoadLocation(originAirport.Timezone)
	destinationTz, err2 := time.LoadLocation(destinationAirport.Timezone)
	if err := errors.Join(err1, err2); err != nil {
		return entity.FlightLeg{}, fmt.Errorf("load timezones: %w", err)
	}

	var departureDateTime time.Time
	if originFound {
		parsed, err := findAndParseTimestampInLocation(origin.Departure.Timings, "STD", originTz)
		if err != nil {
			return entity.FlightLeg{}, fmt.Errorf("parse timestamp: %w", err)
		}
		departureDateTime = parsed
	}

	var arrivalDateTime time.Time
	if destinationFound {
		parsed, err := findAndParseTimestampInLocation(destination.Arrival.Timings, "STA", destinationTz)
		if err != nil {
			return entity.FlightLeg{}, fmt.Errorf("parse timestamp: %w", err)
		}
		arrivalDateTime = parsed
	}

	if !originFound {
		departureDateTime = arrivalDateTime.Add(-duration.Duration()).In(originTz)
	}

	if !destinationFound {
		arrivalDateTime = departureDateTime.Add(duration.Duration()).In(destinationTz)
	}

	aircraftName, err := a.iataLookup.LookupAircraftName(leg.AircraftEquipment.AircraftType)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("lookup aircraft name: %w", err)
	}

	airlineName, err := a.iataLookup.LookupAirlineName(datedFlight.FlightDesignator.CarrierCode)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("lookup airline name: %w", err)
	}

	return entity.FlightLeg{
		Origin:            originAirport.Airport,
		Destination:       destinationAirport.Airport,
		Airline:           airlineName,
		FlightNumber:      fmt.Sprintf("%s %d", datedFlight.FlightDesignator.CarrierCode, datedFlight.FlightDesignator.FlightNumber),
		DepartureDateTime: civil.DateTimeOf(departureDateTime),
		ArrivalDateTime:   civil.DateTimeOf(arrivalDateTime),
		DurationInMinutes: int32(duration.Duration().Minutes()),
		Aircraft:          &aircraftName,
	}, nil
}

func findLeg(legs []Leg, requestedOrigin *string) (Leg, error) {
	if requestedOrigin == nil {
		return legs[0], nil
	}

	for _, leg := range legs {
		if leg.BoardPointIataCode == *requestedOrigin {
			return leg, nil
		}
	}

	return legs[0], nil
}

func findFlightPointByIata(flightContract DatedFlight, iata string) (FlightPoint, bool) {
	for _, point := range flightContract.FlightPoints {
		if point.IataCode == iata {
			return point, true
		}
	}
	return FlightPoint{}, false
}

func findAndParseTimestampInLocation(timings []Timing, preferredQualifier string, loc *time.Location) (time.Time, error) {
	if len(timings) == 0 {
		return time.Time{}, fmt.Errorf("timings empty")
	}

	layout := "2006-01-02T15:04-07:00"
	for _, x := range timings {
		if x.Qualifier == preferredQualifier {
			return time.ParseInLocation(layout, x.Value, loc)
		}
	}

	return time.ParseInLocation(layout, timings[0].Value, loc)
}

func (a *AmadeusWebAPI) getAccessToken(ctx context.Context) (AccessTokenResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/security/oauth2/token", a.baseURL)

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {a.apiKey},
		"client_secret": {a.apiSecret},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return AccessTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AccessTokenResponse{}, fmt.Errorf("do http request: %w", err)
	} else if res.StatusCode != 200 {
		return AccessTokenResponse{}, fmt.Errorf("http status code %d: %w", res.StatusCode, err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return AccessTokenResponse{}, fmt.Errorf("read response body: %w", err)
	}

	var respData AccessTokenResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return AccessTokenResponse{}, fmt.Errorf("unmarshal JSON: %w", err)
	}

	return respData, nil
}
