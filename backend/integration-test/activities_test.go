package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestPostActivity() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	_, err := suite.api.PostActivity(suite.T().Context(), &api.RequestActivity{
		Name:        "My Activity",
		Date:        "2025-02-01",
		Description: api.NewNilString("Description"),
		Time:        api.NilString{Null: true},
		Address:     api.NewNilString("Some Address 1"),
		Location: api.NewNilEntityLocation(api.EntityLocation{
			Latitude:  12.34,
			Longitude: 43.21,
		}),
		Price: api.NewNilInt(100),
	}, api.PostActivityParams{TripID: tripID})
	suite.NoError(err)

	// when
	getAll, err := suite.api.GetActivities(suite.T().Context(), api.GetActivitiesParams{TripID: tripID})
	suite.NoError(err)

	// then
	activities := getAll.(*api.GetActivitiesOKApplicationJSON)
	suite.Len(*activities, 1)

	activity := (*activities)[0]
	suite.Equal("My Activity", activity.Name)
	suite.Equal(100, activity.Price.Value)

	// when (forbiddenUser)
	getAll, _ = suite.userApi(ForbiddenUser).GetActivities(suite.T().Context(), api.GetActivitiesParams{TripID: tripID})
	getByID, _ := suite.userApi(ForbiddenUser).GetActivity(suite.T().Context(), api.GetActivityParams{TripID: tripID, ActivityID: activity.ID})

	// then (forbiddenUser)
	_, ok := getAll.(*api.GetActivitiesForbidden)
	suite.True(ok)

	_, ok = getByID.(*api.GetActivityForbidden)
	suite.True(ok)
}

func (suite *IntegrationTestSuite) TestCreateActivityOutsideOfTripDates() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	response, err := suite.api.PostActivity(suite.T().Context(), &api.RequestActivity{
		Name:        "My Activity",
		Date:        "2000-01-01",
		Description: api.NilString{Null: true},
		Time:        api.NilString{Null: true},
		Address:     api.NilString{Null: true},
		Location:    api.NilEntityLocation{Null: true},
		Price:       api.NilInt{Null: true},
	}, api.PostActivityParams{TripID: tripID})
	suite.NoError(err)

	_, ok := response.(*api.PostActivityNoContent)
	suite.False(ok)

	// when
	res, err := suite.api.GetActivities(suite.T().Context(), api.GetActivitiesParams{TripID: tripID})
	suite.NoError(err)

	// then
	activities := res.(*api.GetActivitiesOKApplicationJSON)
	suite.Empty(*activities)
}
