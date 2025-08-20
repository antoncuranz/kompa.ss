package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestGetTransportation() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	//suite.api.PostFlight(suite.T().Context(), &api.RequestFlight{
	//	Legs:  nil,
	//	Pnrs:  nil,
	//	Price: 0,
	//})

	// when
	res, err := suite.api.GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	suite.NoError(err)

	// then
	transportation := res.(*api.GetAllTransportationOKApplicationJSON)
	suite.Empty(*transportation)
}
