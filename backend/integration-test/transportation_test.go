package integration_test

import (
	"context"
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestCrudTransportation() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	// when (post)
	transportation := suite.postAndRetrieveTransportation(tripID)

	// then (post)
	suite.Equal("My Transportation", transportation.GenericDetail.Value.Name)
	suite.Equal(123, transportation.Price.Value)

	// when (put)
	updatedTransportation := suite.putAndRetrieveTransportation(tripID, transportation.ID)

	// then (put)
	suite.Equal("Updated Transportation", updatedTransportation.GenericDetail.Value.Name)
	suite.True(updatedTransportation.Price.Null)

	// when (delete)
	deleteByID, err := suite.api.DeleteTransportation(suite.T().Context(), api.DeleteTransportationParams{TripID: tripID, TransportationID: transportation.ID})

	// then (delete)
	suite.NoError(err)
	suite.IsType(&api.DeleteTransportationNoContent{}, deleteByID)

	getByID, err := suite.api.GetTransportation(suite.T().Context(), api.GetTransportationParams{TripID: tripID, TransportationID: transportation.ID})
	suite.NoError(err)
	suite.IsType(&api.GetTransportationNotFound{}, getByID)
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
		{"otherTripID getByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.GetTransportation(ctx, api.GetTransportationParams{TripID: otherTripID, TransportationID: flight.ID})
			},
			&api.GetTransportationNotFound{},
		},
		{"otherTripID putByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.PutFlight(ctx, api.PutFlightParams{TripID: otherTripID, FlightID: flight.ID})
			},
			&api.PutFlightNotFound{},
		},
		{"otherTripID deleteByID",
			func(ctx context.Context, client *api.Client) (interface{}, error) {
				return client.DeleteTransportation(ctx, api.DeleteTransportationParams{TripID: otherTripID, TransportationID: flight.ID})
			},
			&api.DeleteTransportationNotFound{},
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

	// check that transportation was neither modified nor deleted by forbidden user
	res, err := suite.api.GetTransportation(suite.T().Context(), api.GetTransportationParams{TripID: actualTripID, TransportationID: flight.ID})
	newFlight := res.(*api.EntityTransportation)
	suite.NoError(err)
	suite.Equal(flight, *newFlight)
}

func (suite *IntegrationTestSuite) postAndRetrieveTransportation(tripID int) api.EntityTransportation {
	_, err := suite.api.PostTransportation(suite.T().Context(), &api.RequestTransportation{
		Name:              "My Transportation",
		Type:              "FERRY",
		DepartureDateTime: "2025-10-06T12:34:00.000000",
		ArrivalDateTime:   "2025-10-06T18:47:00.000000",
		Origin: api.NewNilEntityLocation(api.EntityLocation{
			Latitude:  -41.29439,
			Longitude: 174.00670,
		}),
		OriginAddress: api.NewNilString("Origin Address"),
		Destination: api.NewNilEntityLocation(api.EntityLocation{
			Latitude:  -41.299795,
			Longitude: 174.77919,
		}),
		DestinationAddress: api.NewNilString("Destination Address"),
		Price:              api.NewNilInt(123),
	}, api.PostTransportationParams{TripID: tripID})
	suite.NoError(err)

	getAll, err := suite.api.GetAllTransportation(suite.T().Context(), api.GetAllTransportationParams{TripID: tripID})
	suite.NoError(err)

	allTransportation := getAll.(*api.GetAllTransportationOKApplicationJSON)
	suite.Len(*allTransportation, 1)

	return (*allTransportation)[0]
}

func (suite *IntegrationTestSuite) putAndRetrieveTransportation(tripID int, TransportationID int) api.EntityTransportation {
	putRes, err := suite.api.PutTransportation(suite.T().Context(), &api.RequestTransportation{
		Name:              "Updated Transportation",
		Type:              "FERRY",
		DepartureDateTime: "2025-10-07T12:34:00.000000",
		ArrivalDateTime:   "2025-10-07T18:47:00.000000",
		Origin: api.NewNilEntityLocation(api.EntityLocation{
			Latitude:  -41.29439,
			Longitude: 174.00670,
		}),
		OriginAddress: api.NilString{Null: true},
		Destination: api.NewNilEntityLocation(api.EntityLocation{
			Latitude:  -41.299795,
			Longitude: 174.77919,
		}),
		DestinationAddress: api.NilString{Null: true},
		Price:              api.NilInt{Null: true},
	}, api.PutTransportationParams{TripID: tripID, TransportationID: TransportationID})
	suite.NoError(err)
	suite.IsType(&api.PutTransportationNoContent{}, putRes)

	getRes, err := suite.api.GetTransportation(suite.T().Context(), api.GetTransportationParams{TripID: tripID, TransportationID: TransportationID})
	suite.NoError(err)
	Transportation, ok := getRes.(*api.EntityTransportation)
	suite.True(ok)

	return *Transportation
}
