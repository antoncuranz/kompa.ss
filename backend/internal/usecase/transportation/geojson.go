package transportation

import (
	"context"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"kompass/internal/entity"
)

func (uc *UseCase) GetAllGeoJson(ctx context.Context, tripID int32) ([]geojson.FeatureCollection, error) {
	return uc.repo.GetAllGeoJson(ctx, tripID)
}

func (uc *UseCase) saveGeoJson(ctx context.Context, transportation entity.Transportation) error {
	featureCollection, err := uc.ors.LookupDirections(ctx, transportation.Origin, transportation.Destination, transportation.Type)
	if err != nil {
		return fmt.Errorf("lookup directions: %w", err)
	}

	featureCollection.ExtraMembers = map[string]interface{}{"transportationType": transportation.Type}
	featureCollection.Append(featureWithProperties(transportation.Origin, transportation))
	featureCollection.Append(featureWithProperties(transportation.Destination, transportation))

	err = uc.repo.SaveGeoJson(ctx, transportation.ID, featureCollection)
	if err != nil {
		return fmt.Errorf("save geojson: %w", err)
	}
	return nil
}

func featureWithProperties(location entity.Location, transportation entity.Transportation) *geojson.Feature {
	feature := geojson.NewFeature(locationToPoint(location))

	feature.Properties["type"] = transportation.Type
	feature.Properties["name"] = transportation.GenericDetail.Name
	feature.Properties["departureDateTime"] = transportation.DepartureDateTime
	feature.Properties["arrivalDateTime"] = transportation.ArrivalDateTime

	return feature
}

func locationToPoint(location entity.Location) orb.Point {
	return orb.Point{
		float64(location.Longitude),
		float64(location.Latitude),
	}
}
