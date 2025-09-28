package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestPostFlight() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	_, err := suite.api.PostFlight(suite.T().Context(), &api.RequestFlight{
		Legs: []api.RequestFlightLeg{{
			Date:          "2026-02-01",
			FlightNumber:  "LH717",
			OriginAirport: api.NilString{Null: true},
		}},
	}, api.PostFlightParams{TripID: tripID})
	suite.NoError(err)

	// when
	res, err := suite.api.GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	suite.NoError(err)

	// then
	allTransportation := res.(*api.GetAllTransportationOKApplicationJSON)
	suite.Len(*allTransportation, 1)

	flight := (*allTransportation)[0]
	suite.Equal(api.EntityTransportationType("PLANE"), flight.Type)

	flightDetail, ok := flight.FlightDetail.Get()
	suite.True(ok)
	suite.Len(flightDetail.Legs, 1)
	suite.Equal("LH 717", flightDetail.Legs[0].FlightNumber)

	// when (forbiddenUser)
	getAll, _ := suite.userApi(ForbiddenUser).GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	getByID, _ := suite.userApi(ForbiddenUser).GetTransportation(suite.T().Context(), api.GetTransportationParams{TripID: tripID, TransportationID: flight.ID})

	// then (forbiddenUser)
	_, ok = getAll.(*api.GetAllTransportationForbidden)
	suite.True(ok)

	_, ok = getByID.(*api.GetTransportationForbidden)
	suite.True(ok)
}

func (suite *IntegrationTestSuite) TestPostTrainJourney() {
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
