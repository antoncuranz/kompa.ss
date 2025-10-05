package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestCrudTrainJourney() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	_, err := suite.api.PostTrainJourney(suite.T().Context(), &api.RequestTrainJourney{
		DepartureDate: "2025-09-20",
		FromStationId: "8011113",
		ToStationId:   "8000261",
		TrainNumbers:  []string{"ICE707"},
	}, api.PostTrainJourneyParams{TripID: tripID})
	suite.NoError(err)

	// when
	res, err := suite.api.GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	suite.NoError(err)

	// then
	allTransportation := res.(*api.GetAllTransportationOKApplicationJSON)
	suite.Len(*allTransportation, 1)

	flight := (*allTransportation)[0]
	suite.Equal(api.EntityTransportationType("TRAIN"), flight.Type)

	trainDetail, ok := flight.TrainDetail.Get()
	suite.True(ok)
	suite.Len(trainDetail.Legs, 1)
	suite.Equal("ICE 707", trainDetail.Legs[0].LineName)
}
