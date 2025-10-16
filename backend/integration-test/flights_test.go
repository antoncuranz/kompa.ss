package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestCrudFlight() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	// when (post)
	flight := suite.postAndRetrieveFlightDetail(tripID, "2026-02-01", "LH717", api.NilString{Null: true})

	// then (post)
	suite.Equal(api.EntityTransportationType("FLIGHT"), flight.Type)

	flightDetail, ok := flight.FlightDetail.Get()
	suite.True(ok)
	suite.Len(flightDetail.Legs, 1)
	suite.Equal("Lufthansa", flightDetail.Legs[0].Airline)
	suite.Equal("LH 717", flightDetail.Legs[0].FlightNumber)
	suite.Equal("Boeing 747-8i", flightDetail.Legs[0].Aircraft.Value)
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

func (suite *IntegrationTestSuite) TestFLightEk412() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	date := "2026-01-30"
	flightNumber := "EK412"

	// when (no origin)
	//noOrigin := suite.postAndRetrieveFlightDetail(tripID, date, flightNumber, api.NilString{Null: true})
	noOrigin, err := suite.api.PostFlight(suite.T().Context(), &api.RequestFlight{
		Legs: []api.RequestFlightLeg{{
			Date:          date,
			FlightNumber:  flightNumber,
			OriginAirport: api.NilString{Null: true},
		}},
	}, api.PostFlightParams{TripID: tripID})

	// then (no origin)
	suite.NoError(err)
	suite.IsType(&api.EntityErrAmbiguousFlightRequest{}, noOrigin)

	// when (DXB origin)
	dxbOrigin := suite.postAndRetrieveFlightDetail(tripID, date, flightNumber, api.NewNilString("DXB"))

	// then (DXB origin)
	suite.Equal(api.EntityTransportationType("FLIGHT"), dxbOrigin.Type)

	dxbOriginDetail, ok := dxbOrigin.FlightDetail.Get()
	suite.True(ok)
	suite.Len(dxbOriginDetail.Legs, 1)
	suite.Equal("EK 412", dxbOriginDetail.Legs[0].FlightNumber)
	suite.Equal("Emirates", dxbOriginDetail.Legs[0].Airline)
	suite.Equal("Airbus A380", dxbOriginDetail.Legs[0].Aircraft.Value)
	suite.Equal("2026-01-30T10:15:00", dxbOriginDetail.Legs[0].DepartureDateTime)
	suite.Equal("2026-01-31T07:00:00", dxbOriginDetail.Legs[0].ArrivalDateTime)

	// when (SYD origin)
	sydOrigin := suite.postAndRetrieveFlightDetail(tripID, date, flightNumber, api.NewNilString("SYD"))

	// then (DXB origin)
	suite.Equal(api.EntityTransportationType("FLIGHT"), sydOrigin.Type)

	sydOriginDetail, ok := sydOrigin.FlightDetail.Get()
	suite.True(ok)
	suite.Len(sydOriginDetail.Legs, 1)
	suite.Equal("EK 412", sydOriginDetail.Legs[0].FlightNumber)
	suite.Equal("Emirates", sydOriginDetail.Legs[0].Airline)
	suite.Equal("Airbus A380", sydOriginDetail.Legs[0].Aircraft.Value)
	suite.Equal("2026-01-31T08:45:00", sydOriginDetail.Legs[0].DepartureDateTime)
	suite.Equal("2026-01-31T14:00:00", sydOriginDetail.Legs[0].ArrivalDateTime)
}

func (suite *IntegrationTestSuite) postAndRetrieveFlightDetail(tripID int, date string, flightNumber string, origin api.NilString) api.EntityTransportation {
	postRes, err := suite.api.PostFlight(suite.T().Context(), &api.RequestFlight{
		Legs: []api.RequestFlightLeg{{
			Date:          date,
			FlightNumber:  flightNumber,
			OriginAirport: origin,
		}},
	}, api.PostFlightParams{TripID: tripID})
	suite.NoError(err)
	suite.IsType(&api.EntityTransportation{}, postRes)
	return *postRes.(*api.EntityTransportation)
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
