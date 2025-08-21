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

	// when
	res, err := suite.api.GetFlights(suite.T().Context(), api.GetFlightsParams{TripID: tripID})
	suite.NoError(err)

	// then
	flights := res.(*api.GetFlightsOKApplicationJSON)
	suite.Len(*flights, 1)

	//// when
	//res, err := suite.api.GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	//suite.NoError(err)
	//
	//// then
	//transportation := res.(*api.GetAllTransportationOKApplicationJSON)
	//suite.Empty(*transportation)
	//suite.Len(*transportation, 1)
}
