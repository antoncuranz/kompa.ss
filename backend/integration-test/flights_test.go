package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestCrudFlight() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	// when (post)
	flight := suite.postAndRetrieveFlight(tripID)

	// then (post)
	suite.Equal(api.EntityTransportationType("FLIGHT"), flight.Type)

	flightDetail, ok := flight.FlightDetail.Get()
	suite.True(ok)
	suite.Len(flightDetail.Legs, 1)
	suite.Equal("LH 717", flightDetail.Legs[0].FlightNumber)
	suite.Equal("Boeing 747-8", flightDetail.Legs[0].Aircraft.Value)
	suite.Equal("2026-02-01T12:35:00", flightDetail.Legs[0].DepartureDateTime)
	suite.Equal("2026-02-01T19:00:00", flightDetail.Legs[0].ArrivalDateTime)

	// when (put)
	updatedFlight := suite.putAndRetrieveFlight(tripID, flight.ID)

	// then (put)
	updatedDetail, ok := updatedFlight.FlightDetail.Get()
	suite.True(ok)
	suite.Len(flightDetail.Legs, 1)
	suite.Equal("Boeing 747-400", updatedDetail.Legs[0].Aircraft.Value)
	suite.Equal("2026-02-01T12:40:00", updatedFlight.DepartureDateTime)
	suite.Equal("2026-02-01T12:40:00", updatedDetail.Legs[0].DepartureDateTime)
	suite.Equal("2026-02-01T19:15:00", updatedFlight.ArrivalDateTime)
	suite.Equal("2026-02-01T19:15:00", updatedDetail.Legs[0].ArrivalDateTime)
}

func (suite *IntegrationTestSuite) postAndRetrieveFlight(tripID int) api.EntityTransportation {
	postRes, err := suite.api.PostFlight(suite.T().Context(), &api.RequestFlight{
		Legs: []api.RequestFlightLeg{{
			Date:          "2026-02-01",
			FlightNumber:  "LH717",
			OriginAirport: api.NilString{Null: true},
		}},
	}, api.PostFlightParams{TripID: tripID})
	suite.NoError(err)
	suite.IsType(&api.PostFlightNoContent{}, postRes)

	getRes, err := suite.api.GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	suite.NoError(err)
	allTransportation := getRes.(*api.GetAllTransportationOKApplicationJSON)
	suite.Len(*allTransportation, 1)

	return (*allTransportation)[0]
}

func (suite *IntegrationTestSuite) putAndRetrieveFlight(tripID int, flightID int) api.EntityTransportation {
	putRes, err := suite.api.PutFlight(suite.T().Context(), api.PutFlightParams{TripID: tripID, FlightID: flightID})
	suite.NoError(err)
	suite.IsType(&api.PutFlightNoContent{}, putRes)

	getRes, err := suite.api.GetTransportation(suite.T().Context(), api.GetTransportationParams{TripID: tripID, TransportationID: flightID})
	suite.NoError(err)
	flight, ok := getRes.(*api.EntityTransportation)
	suite.True(ok)

	return *flight
}
