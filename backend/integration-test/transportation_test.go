package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestGetTransportation() {
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

	//// when
	res, err := suite.api.GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	suite.NoError(err)

	// then
	allTransportation := res.(*api.GetAllTransportationOKApplicationJSON)
	suite.Len(*allTransportation, 1)

	flight := (*allTransportation)[0]
	suite.Equal("PLANE", flight.Type)

	flightDetail, ok := flight.FlightDetail.Get()
	suite.True(ok)
	suite.Len(flightDetail.Legs, 1)
	suite.Equal("LH 717", flightDetail.Legs[0].FlightNumber)
}
