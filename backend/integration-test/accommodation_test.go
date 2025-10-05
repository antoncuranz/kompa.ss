package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestCrudAccommodation() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	// when (post)
	accommodation := suite.postAndRetrieveAccommodation(tripID)

	// then (post)
	suite.Equal("My Accommodation", accommodation.Name)
	suite.Equal(25000, accommodation.Price.Value)

	// when (put)
	updatedAccommodation := suite.putAndRetrieveAccommodation(tripID, accommodation.ID)

	// then (put)
	suite.Equal("Updated Accommodation", updatedAccommodation.Name)
	suite.True(updatedAccommodation.Price.Null)

	// when (delete)
	deleteByID, err := suite.api.DeleteAccommodation(suite.T().Context(), api.DeleteAccommodationParams{TripID: tripID, AccommodationID: accommodation.ID})

	// then (delete)
	suite.NoError(err)
	suite.IsType(&api.DeleteAccommodationNoContent{}, deleteByID)

	getByID, err := suite.api.GetAccommodationByID(suite.T().Context(), api.GetAccommodationByIDParams{TripID: tripID, AccommodationID: accommodation.ID})
	suite.NoError(err)
	suite.IsType(&api.GetAccommodationByIDNotFound{}, getByID)
}

func (suite *IntegrationTestSuite) TestCreateAccommodationOutsideOfTripDates() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	response, err := suite.api.PostAccommodation(suite.T().Context(), &api.RequestAccommodation{
		Name:          "My Accommodation",
		DepartureDate: "2020-01-01",
		ArrivalDate:   "2020-01-04",
		Description:   api.NewNilString("Description"),
		CheckInTime:   api.NilString{Null: true},
		CheckOutTime:  api.NilString{Null: true},
		Address:       api.NilString{Null: true},
		Location:      api.NilEntityLocation{Null: true},
		Price:         api.NewNilInt(25000),
	}, api.PostAccommodationParams{TripID: tripID})
	suite.NoError(err)

	suite.IsType(&api.PostAccommodationBadRequest{}, response)

	// when
	res, err := suite.api.GetActivities(suite.T().Context(), api.GetActivitiesParams{TripID: tripID})
	suite.NoError(err)

	// then
	activities := res.(*api.GetActivitiesOKApplicationJSON)
	suite.Empty(*activities)
}

func (suite *IntegrationTestSuite) postAndRetrieveAccommodation(tripID int) api.EntityAccommodation {
	_, err := suite.api.PostAccommodation(suite.T().Context(), &api.RequestAccommodation{
		Name:          "My Accommodation",
		DepartureDate: "2025-01-01",
		ArrivalDate:   "2025-01-04",
		Description:   api.NewNilString("Description"),
		CheckInTime:   api.NewNilString("12:34:56"),
		CheckOutTime:  api.NewNilString("12:34:56"),
		Address:       api.NewNilString("Address 123"),
		Location: api.NewNilEntityLocation(api.EntityLocation{
			Latitude:  12.34,
			Longitude: 43.21,
		}),
		Price: api.NewNilInt(25000),
	}, api.PostAccommodationParams{TripID: tripID})
	suite.NoError(err)

	getAll, err := suite.api.GetAllAccommodation(suite.T().Context(), api.GetAllAccommodationParams{TripID: tripID})
	suite.NoError(err)

	allAccommodation := getAll.(*api.GetAllAccommodationOKApplicationJSON)
	suite.Len(*allAccommodation, 1)

	return (*allAccommodation)[0]
}

func (suite *IntegrationTestSuite) putAndRetrieveAccommodation(tripID int, AccommodationID int) api.EntityAccommodation {
	putRes, err := suite.api.PutAccommodation(suite.T().Context(), &api.RequestAccommodation{
		Name:          "Updated Accommodation",
		DepartureDate: "2025-01-01",
		ArrivalDate:   "2025-01-04",
		Description:   api.NewNilString("Description"),
		CheckInTime:   api.NilString{Null: true},
		CheckOutTime:  api.NilString{Null: true},
		Address:       api.NilString{Null: true},
		Location:      api.NilEntityLocation{Null: true},
		Price:         api.NilInt{Null: true},
	}, api.PutAccommodationParams{TripID: tripID, AccommodationID: AccommodationID})
	suite.NoError(err)
	suite.IsType(&api.PutAccommodationNoContent{}, putRes)

	getRes, err := suite.api.GetAccommodationByID(suite.T().Context(), api.GetAccommodationByIDParams{TripID: tripID, AccommodationID: AccommodationID})
	suite.NoError(err)
	accommodation, ok := getRes.(*api.EntityAccommodation)
	suite.True(ok)

	return *accommodation
}
