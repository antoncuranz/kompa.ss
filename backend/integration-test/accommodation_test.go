package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestPostAccommodation() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	_, err := suite.api.PostAccommodation(suite.T().Context(), &api.RequestAccommodation{
		Name:          "My Accommodation",
		DepartureDate: "2025-01-01",
		ArrivalDate:   "2025-01-04",
		Description:   api.NewNilString("Description"),
		CheckInTime:   api.NilString{Null: true},
		CheckOutTime:  api.NilString{Null: true},
		Address:       api.NilString{Null: true},
		Location:      api.NilEntityLocation{Null: true},
		Price:         api.NewNilInt(25000),
	}, api.PostAccommodationParams{TripID: tripID})
	suite.NoError(err)

	// when
	getAll, err := suite.api.GetAllAccommodation(suite.T().Context(), api.GetAllAccommodationParams{TripID: tripID})
	suite.NoError(err)

	// then
	allAccommodation := getAll.(*api.GetAllAccommodationOKApplicationJSON)
	suite.Len(*allAccommodation, 1)

	accommodation := (*allAccommodation)[0]
	suite.Equal("My Accommodation", accommodation.Name)
	suite.Equal(25000, accommodation.Price.Value)

	// when (forbiddenUser)
	getAll, err = suite.userApi(ForbiddenUser).GetAllAccommodation(suite.T().Context(), api.GetAllAccommodationParams{TripID: tripID})
	suite.NoError(err)
	getByID, err := suite.userApi(ForbiddenUser).GetAccommodationByID(suite.T().Context(), api.GetAccommodationByIDParams{TripID: tripID, AccommodationID: accommodation.ID})
	suite.NoError(err)

	// then (forbiddenUser)
	suite.IsType(&api.GetAllAccommodationForbidden{}, getAll)
	suite.IsType(&api.GetAccommodationByIDForbidden{}, getByID)
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

	suite.IsNotType(&api.PostAccommodationNoContent{}, response)

	// when
	res, err := suite.api.GetActivities(suite.T().Context(), api.GetActivitiesParams{TripID: tripID})
	suite.NoError(err)

	// then
	activities := res.(*api.GetActivitiesOKApplicationJSON)
	suite.Empty(*activities)
}
