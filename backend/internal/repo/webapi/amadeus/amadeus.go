package amadeus

import (
	"cloud.google.com/go/civil"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	goiso8601duration "github.com/xnacly/go-iso8601-duration"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"kompass/config"
	"kompass/internal/entity"
	"net/http"
	"net/url"
	"strings"
)

type AmadeusWebAPI struct {
	baseURL   string
	apiKey    string
	apiSecret string
}

func New(config config.WebApi) *AmadeusWebAPI {
	return &AmadeusWebAPI{
		baseURL:   config.AmadeusBaseURL,
		apiKey:    config.AmadeusApiKey,
		apiSecret: config.AmadeusApiSecret,
	}
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

func (a *AmadeusWebAPI) RetrieveFlightLeg(ctx context.Context, date civil.Date, flightNumber string, unsupported *string) (entity.FlightLeg, error) {
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

	if len(flightStatusResponse.Data) != 1 || len(flightStatusResponse.Data[0].Legs) != 1 {
		return entity.FlightLeg{}, fmt.Errorf("not found or too complex")
	}

	datedFlight := flightStatusResponse.Data[0]
	leg := datedFlight.Legs[0]

	return a.mapDatedFlight(ctx, datedFlight, leg)
}

func (a *AmadeusWebAPI) mapDatedFlight(ctx context.Context, datedFlight DatedFlight, leg Leg) (entity.FlightLeg, error) {
	duration, err := goiso8601duration.From(leg.ScheduledLegDuration)
	if err != nil {
		return entity.FlightLeg{}, fmt.Errorf("parse duration: %w", err)
	}

	origin, err1 := findFlightPointByIata(datedFlight, leg.BoardPointIataCode)
	destination, err2 := findFlightPointByIata(datedFlight, leg.OffPointIataCode)
	if err := errors.Join(err1, err2); err != nil {
		return entity.FlightLeg{}, fmt.Errorf("no flight point found")
	}

	departureDateTime, err1 := findAndParseTimestamp(origin.Departure.Timings, "STD")
	arrivalDateTime, err2 := findAndParseTimestamp(destination.Arrival.Timings, "STA")
	if err := errors.Join(err1, err2); err != nil {
		return entity.FlightLeg{}, fmt.Errorf("parse timestamps: %w", err)
	}

	originAirport, err1 := a.retrieveAirport(ctx, origin.IataCode)
	destinationAirport, err1 := a.retrieveAirport(ctx, destination.IataCode)
	if err := errors.Join(err1, err2); err != nil {
		return entity.FlightLeg{}, fmt.Errorf("retrieve airports: %w", err)
	}

	return entity.FlightLeg{
		Origin:            originAirport,
		Destination:       destinationAirport,
		Airline:           datedFlight.FlightDesignator.CarrierCode,
		FlightNumber:      fmt.Sprintf("%s %d", datedFlight.FlightDesignator.CarrierCode, datedFlight.FlightDesignator.FlightNumber),
		DepartureDateTime: departureDateTime,
		ArrivalDateTime:   arrivalDateTime,
		DurationInMinutes: int32(duration.Duration().Minutes()),
		Aircraft:          &leg.AircraftEquipment.AircraftType,
	}, nil
}

func (a *AmadeusWebAPI) retrieveAirport(ctx context.Context, iata string) (entity.Airport, error) {
	urlFormat := "%s/v1/reference-data/locations?subType=AIRPORT&keyword=%s"
	locationsUrl := fmt.Sprintf(urlFormat, a.baseURL, iata)

	req, err := http.NewRequestWithContext(ctx, "GET", locationsUrl, nil)
	if err != nil {
		return entity.Airport{}, fmt.Errorf("create http request: %w", err)
	}

	accessToken, err := a.getAccessToken(ctx)
	if err != nil {
		return entity.Airport{}, fmt.Errorf("get access token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.Airport{}, fmt.Errorf("do http request: %w", err)
	} else if res.StatusCode != 200 {
		return entity.Airport{}, fmt.Errorf("http status code %d: %w", res.StatusCode, err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return entity.Airport{}, fmt.Errorf("read response body: %w", err)
	}

	var airportStatusResponse AirportResponse
	if err := json.Unmarshal(body, &airportStatusResponse); err != nil {
		return entity.Airport{}, fmt.Errorf("unmarshall JSON: %w", err)
	}

	if len(airportStatusResponse.Data) == 0 {
		return entity.Airport{}, fmt.Errorf("aiport not found")
	}

	airport := airportStatusResponse.Data[0]

	caser := cases.Title(language.English)
	name := airport.Name
	municipality := airport.Address.CityName
	if !strings.HasPrefix(name, municipality) {
		name = fmt.Sprintf("%s %s", municipality, name)
	}

	return entity.Airport{
		Iata:         airport.IataCode,
		Name:         caser.String(name),
		Municipality: caser.String(municipality),
		Location: entity.Location{
			Latitude:  airport.GeoCode.Latitude,
			Longitude: airport.GeoCode.Longitude,
		},
	}, nil
}

func findFlightPointByIata(flightContract DatedFlight, iata string) (FlightPoint, error) {
	for _, point := range flightContract.FlightPoints {
		if point.IataCode == iata {
			return point, nil
		}
	}
	return FlightPoint{}, fmt.Errorf("flight point '%s' not found", iata)
}

func findAndParseTimestamp(timings []Timing, preferredQualifier string) (civil.DateTime, error) {
	if len(timings) == 0 {
		return civil.DateTime{}, fmt.Errorf("timings empty")
	}

	for _, x := range timings {
		if x.Qualifier == preferredQualifier {
			return parseLocalDateTime(x.Value)
		}
	}

	return parseLocalDateTime(timings[0].Value)
}

func parseLocalDateTime(timestamp string) (civil.DateTime, error) {
	// Remove timezone and add missing seconds
	timestamp = timestamp[0:16] + ":00"
	return civil.ParseDateTime(timestamp)
}
