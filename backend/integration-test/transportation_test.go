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
	suite.Equal(api.EntityTransportationType("PLANE"), flight.Type)

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
	suite.Equal("2026-02-01T12:40:00", updatedDetail.Legs[0].DepartureDateTime)
	suite.Equal("2026-02-01T19:15:00", updatedDetail.Legs[0].ArrivalDateTime)

	// when (forbiddenUser)
	getAll, _ := suite.userApi(ForbiddenUser).GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	getByID, _ := suite.userApi(ForbiddenUser).GetTransportation(suite.T().Context(), api.GetTransportationParams{TripID: tripID, TransportationID: flight.ID})
	updateByID, _ := suite.userApi(ForbiddenUser).PutFlight(suite.T().Context(), api.PutFlightParams{TripID: tripID, FlightID: flight.ID})
	deleteByID, _ := suite.userApi(ForbiddenUser).DeleteTransportation(suite.T().Context(), api.DeleteTransportationParams{TripID: tripID, TransportationID: flight.ID})

	// then (forbiddenUser)
	_, ok = getAll.(*api.GetAllTransportationForbidden)
	suite.True(ok)
	_, ok = getByID.(*api.GetTransportationForbidden)
	suite.True(ok)
	_, ok = updateByID.(*api.PutFlightForbidden)
	suite.True(ok)
	_, ok = deleteByID.(*api.DeleteTransportationForbidden)
	suite.True(ok)
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
	_, ok := postRes.(*api.PostFlightNoContent)
	suite.True(ok)

	getRes, err := suite.api.GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	suite.NoError(err)
	allTransportation := getRes.(*api.GetAllTransportationOKApplicationJSON)
	suite.Len(*allTransportation, 1)

	return (*allTransportation)[0]
}

func (suite *IntegrationTestSuite) putAndRetrieveFlight(tripID int, flightID int) api.EntityTransportation {
	putRes, err := suite.api.PutFlight(suite.T().Context(), api.PutFlightParams{TripID: tripID, FlightID: flightID})
	suite.NoError(err)
	_, ok := putRes.(*api.PutFlightNoContent)
	suite.True(ok)

	getRes, err := suite.api.GetTransportation(suite.T().Context(), api.GetTransportationParams{TripID: tripID, TransportationID: flightID})
	suite.NoError(err)
	flight, ok := getRes.(*api.EntityTransportation)
	suite.True(ok)

	return *flight
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
