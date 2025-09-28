package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestGetTripsOK() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	// when
	getAll, err := suite.api.GetTrips(suite.T().Context())
	suite.NoError(err)
	getByID, err := suite.api.GetTrip(suite.T().Context(), api.GetTripParams{TripID: tripID})
	suite.NoError(err)

	// then
	trips := getAll.(*api.GetTripsOKApplicationJSON)
	suite.Len(*trips, 1)

	_, ok := getByID.(*api.EntityTrip)
	suite.True(ok)

	// when (forbiddenUser)
	getAll, _ = suite.userApi(ForbiddenUser).GetTrips(suite.T().Context())
	getByID, _ = suite.userApi(ForbiddenUser).GetTrip(suite.T().Context(), api.GetTripParams{TripID: tripID})

	// then (forbiddenUser)
	trips = getAll.(*api.GetTripsOKApplicationJSON)
	suite.Empty(*trips)

	_, ok = getByID.(*api.EntityTrip)
	suite.False(ok)
}

func (suite *IntegrationTestSuite) TestGetTripNotFound() {
	// given

	// when
	res, err := suite.api.GetTrip(suite.T().Context(), api.GetTripParams{TripID: 404})
	suite.NoError(err)

	// then
	_, ok := res.(*api.GetTripForbidden)
	suite.True(ok)
}
