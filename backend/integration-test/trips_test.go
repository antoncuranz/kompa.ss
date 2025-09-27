package integration_test

import (
	"fmt"
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestGetTripOK() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	// when
	res, err := suite.api.GetTrip(suite.T().Context(), api.GetTripParams{TripID: tripID})
	suite.NoError(err)

	// then
	trip := res.(*api.EntityTrip)
	fmt.Println("Trip found: ", trip)

	// when (forbiddenUser)
	res, err = suite.userApi(ForbiddenUser).GetTrip(suite.T().Context(), api.GetTripParams{TripID: tripID})
	suite.NoError(err)

	// then (forbiddenUser)
	notFound := res.(*api.GetTripNotFound)
	fmt.Println("Error message: ", notFound.Error)
}

func (suite *IntegrationTestSuite) TestGetTripNotFound() {
	// given

	// when
	res, err := suite.api.GetTrip(suite.T().Context(), api.GetTripParams{TripID: 404})
	suite.NoError(err)

	// then
	trip := res.(*api.GetTripNotFound)
	fmt.Println("Error message: ", trip.Error)
}
