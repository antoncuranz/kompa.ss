package integration_test

import (
	"context"
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestCrudTransportation() {
	// TODO
}

func (suite *IntegrationTestSuite) TestTransportationNotFound() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	// when
	res, err := suite.api.GetTransportation(suite.T().Context(), api.GetTransportationParams{
		TripID:           tripID,
		TransportationID: 404,
	})

	// then
	suite.NoError(err)
	suite.IsType(api.GetTransportationNotFound{}, res)
}

func (suite *IntegrationTestSuite) TestForbiddenUserPermissions() {
	// given
	actualTripID := suite.CreateTrip()
	defer suite.DeleteTrip(actualTripID)
	flight := suite.postAndRetrieveFlight(actualTripID)

	otherTripID := suite.CreateTripUser(ForbiddenUser)
	defer suite.DeleteTripUser(ForbiddenUser, otherTripID)

	var tests = []struct {
		name             string
		apiCall          func(ctx context.Context, client *api.Client) (interface{}, error)
		expectedResponse interface{}
	}{
		{"actualTripID getAll",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.GetAllTransportation(ctx, api.GetAllTransportationParams{TripID: actualTripID})
			},
			&api.GetAllTransportationForbidden{},
		},
		{"actualTripID getByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.GetTransportation(ctx, api.GetTransportationParams{TripID: actualTripID, TransportationID: flight.ID})
			},
			&api.GetTransportationForbidden{},
		},
		{"actualTripID putByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.PutFlight(ctx, api.PutFlightParams{TripID: actualTripID, FlightID: flight.ID})
			},
			&api.PutFlightForbidden{},
		},
		{"actualTripID deleteByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.DeleteTransportation(ctx, api.DeleteTransportationParams{TripID: actualTripID, TransportationID: flight.ID})
			},
			&api.DeleteTransportationForbidden{},
		},
		// TODO: in the following cases, a 404 should be ok
		{"otherTripID getByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.GetTransportation(ctx, api.GetTransportationParams{TripID: otherTripID, TransportationID: flight.ID})
			},
			&api.GetTransportationForbidden{},
		},
		{"otherTripID putByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.PutFlight(ctx, api.PutFlightParams{TripID: otherTripID, FlightID: flight.ID})
			},
			&api.PutFlightForbidden{},
		},
		{"otherTripID deleteByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.DeleteTransportation(ctx, api.DeleteTransportationParams{TripID: otherTripID, TransportationID: flight.ID})
			},
			&api.DeleteTransportationForbidden{},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			// when
			response, err := test.apiCall(suite.T().Context(), suite.userApi(ForbiddenUser))

			// then
			suite.NoError(err)
			suite.IsType(test.expectedResponse, response)
		})
	}
}
