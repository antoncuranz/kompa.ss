package webapi

import (
	"context"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"kompass/config"
	"kompass/internal/entity"
	"net/url"
)

type OpenRouteServiceWebAPI struct {
	baseURL string
	apiKey  string
}

func NewOpenRouteServiceWebAPI(config config.WebApi) *OpenRouteServiceWebAPI {
	return &OpenRouteServiceWebAPI{
		baseURL: config.OpenRouteServiceBaseURL,
		apiKey:  config.OpenRouteServiceApiKey,
	}
}

func (a *OpenRouteServiceWebAPI) LookupLocation(ctx context.Context, query string) (entity.GeocodeLocation, error) {
	urlFormat := "%s/geocode/search?api_key=%s&size=1&text=%s"
	searchUrl := fmt.Sprintf(urlFormat, a.baseURL, a.apiKey, url.QueryEscape(query))

	result, err := RequestAndParseJsonBody[geojson.FeatureCollection](ctx, "GET", searchUrl, nil)
	if err != nil {
		return entity.GeocodeLocation{}, fmt.Errorf("requestAndParseJsonBody: %w", err)
	}

	if len(result.Features) == 0 {
		return entity.GeocodeLocation{}, fmt.Errorf("no train stations found")
	}
	feature := result.Features[0]
	point := feature.Geometry.(orb.Point)

	return entity.GeocodeLocation{
		Label:     feature.Properties["label"].(string),
		Latitude:  float32(point[1]),
		Longitude: float32(point[0]),
	}, nil
}
