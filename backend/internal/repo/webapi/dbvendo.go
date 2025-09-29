package webapi

import (
	"context"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"kompass/config"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo/webapi/converter"
	"kompass/internal/repo/webapi/response"
	"strings"
)

type DbVendoWebAPI struct {
	baseURL string
	c       converter.TrainConverter
}

func NewDbVendoWebAPI(config config.WebApi) *DbVendoWebAPI {
	return &DbVendoWebAPI{
		baseURL: config.DbVendoBaseURL,
		c:       &converter.TrainConverterImpl{},
	}
}

func (a *DbVendoWebAPI) RetrieveLocation(ctx context.Context, query string) (entity.TrainStation, error) {
	urlFormat := "%s/locations?query=%s&poi=false"
	url := fmt.Sprintf(urlFormat, a.baseURL, query)

	results, err := RequestAndParseJsonBody[[]response.StationOrStop](ctx, "GET", url, nil)
	if err != nil {
		return entity.TrainStation{}, fmt.Errorf("requestAndParseJsonBody: %w", err)
	}

	if len(*results) == 0 {
		return entity.TrainStation{}, fmt.Errorf("no train stations found")
	}

	return a.c.ConvertStation((*results)[0]), nil
}

func (a *DbVendoWebAPI) RetrievePolylines(ctx context.Context, refreshToken string) ([]geojson.FeatureCollection, error) {
	urlFormat := "%s/journeys/%s?polylines=true"
	url := fmt.Sprintf(urlFormat, a.baseURL, refreshToken)

	rsp, err := RequestAndParseJsonBody[response.JourneyResponse](ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("requestAndParseJsonBody: %w", err)
	}

	featureCollections := []geojson.FeatureCollection{}

	for _, leg := range rsp.Journey.Legs {
		if leg.Polyline != nil {
			featureCollections = append(featureCollections, *leg.Polyline)
		}
	}

	return featureCollections, nil
}

const JourneyUrlFormatBase = "%s/journeys?from=%s&to=%s&transfers=%d&results=10"
const MaxRetries = 10

func (a *DbVendoWebAPI) RetrieveJourney(ctx context.Context, request request.TrainJourney) (entity.TrainDetail, error) {
	journeys, err := a.retrieveJourneysInitial(ctx, request)
	if err != nil {
		return entity.TrainDetail{}, fmt.Errorf("retrieveJourneysInitial: %w", err)
	}

	journey, ok := checkJourneys(journeys.Journeys, request)
	if ok {
		return a.convertJourney(journey)
	}

	for range MaxRetries {
		journeys, err = a.retrieveJourneysLaterThan(ctx, request, journeys.LaterRef)
		if err != nil {
			return entity.TrainDetail{}, fmt.Errorf("retrieveJourneysLaterThan: %w", err)
		}

		journey, ok := checkJourneys(journeys.Journeys, request)
		if ok {
			return a.convertJourney(journey)
		}
	}

	return entity.TrainDetail{}, fmt.Errorf("no journeys found after %d tries", MaxRetries)
}

func (a *DbVendoWebAPI) retrieveJourneysInitial(ctx context.Context, journey request.TrainJourney) (*response.JourneysResponse, error) {
	urlFormat := JourneyUrlFormatBase + "&departure=%s"
	url := fmt.Sprintf(urlFormat, a.baseURL, journey.FromStationID, journey.ToStationID, len(journey.TrainNumbers), journey.DepartureDate)

	return RequestAndParseJsonBody[response.JourneysResponse](ctx, "GET", url, nil)
}

func (a *DbVendoWebAPI) retrieveJourneysLaterThan(ctx context.Context, journey request.TrainJourney, laterRef string) (*response.JourneysResponse, error) {
	urlFormat := JourneyUrlFormatBase + "&laterThan=%s"
	url := fmt.Sprintf(urlFormat, a.baseURL, journey.FromStationID, journey.ToStationID, len(journey.TrainNumbers), laterRef)

	return RequestAndParseJsonBody[response.JourneysResponse](ctx, "GET", url, nil)
}

func (a *DbVendoWebAPI) convertJourney(source response.Journey) (entity.TrainDetail, error) {
	legs := []entity.TrainLeg{}

	for _, leg := range source.Legs {
		if leg.Line == nil {
			continue
		}

		convertedLeg, err := a.c.ConvertLeg(leg)
		if err != nil {
			return entity.TrainDetail{}, err
		}
		legs = append(legs, convertedLeg)
	}

	return entity.TrainDetail{
		RefreshToken: source.RefreshToken,
		Legs:         legs,
	}, nil
}

func checkJourneys(journeys []response.Journey, request request.TrainJourney) (response.Journey, bool) {
journeyLoop:
	for _, journey := range journeys {
		requestLegIdx := 0

		for _, leg := range journey.Legs {
			if leg.Line == nil {
				continue
			}

			if !equalIgnoringWhitespaceAndCase(leg.Line.Name, request.TrainNumbers[requestLegIdx]) {
				continue journeyLoop
			}

			if requestLegIdx = requestLegIdx + 1; requestLegIdx == len(request.TrainNumbers) {
				return journey, true
			}
		}

		return journey, true
	}

	return response.Journey{}, false
}

func equalIgnoringWhitespaceAndCase(s, t string) bool {
	sWithoutWhitespace := strings.ReplaceAll(s, " ", "")
	tWithoutWhitespace := strings.ReplaceAll(t, " ", "")
	return strings.EqualFold(sWithoutWhitespace, tWithoutWhitespace)
}
