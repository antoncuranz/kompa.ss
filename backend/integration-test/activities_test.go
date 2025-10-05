package integration_test

import (
	"kompass/integration-test/client/api"
)

func (suite *IntegrationTestSuite) TestCrudActivity() {
	// given
	tripID := suite.CreateTrip()
	defer suite.DeleteTrip(tripID)

	// when (post)
	activity := suite.postAndRetrieveActivity(tripID)

	// then (post)
	suite.Equal("My Activity", activity.Name)
	suite.Equal(100, activity.Price.Value)

	// when (put)
	updatedActivity := suite.putAndRetrieveActivity(tripID, activity.ID)

	// then (put)
	suite.Equal("Updated Activity", updatedActivity.Name)
	suite.True(updatedActivity.Price.Null)

	// when (delete)
	deleteByID, err := suite.api.DeleteActivity(suite.T().Context(), api.DeleteActivityParams{TripID: tripID, ActivityID: activity.ID})

	// then (delete)
	suite.NoError(err)
	suite.IsType(&api.DeleteActivityNoContent{}, deleteByID)

	getByID, err := suite.api.GetActivity(suite.T().Context(), api.GetActivityParams{TripID: tripID, ActivityID: activity.ID})
	suite.NoError(err)
	suite.IsType(&api.GetActivityNotFound{}, getByID)
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

	suite.IsType(&api.PostActivityBadRequest{}, response)

	// when
	res, err := suite.api.GetActivities(suite.T().Context(), api.GetActivitiesParams{TripID: tripID})
	suite.NoError(err)

	// then
	activities := res.(*api.GetActivitiesOKApplicationJSON)
	suite.Empty(*activities)
}

func (suite *IntegrationTestSuite) postAndRetrieveActivity(tripID int) api.EntityActivity {
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

	getAll, err := suite.api.GetActivities(suite.T().Context(), api.GetActivitiesParams{TripID: tripID})
	suite.NoError(err)

	activities := getAll.(*api.GetActivitiesOKApplicationJSON)
	suite.Len(*activities, 1)

	return (*activities)[0]
}

func (suite *IntegrationTestSuite) putAndRetrieveActivity(tripID int, activityID int) api.EntityActivity {
	putRes, err := suite.api.PutActivity(suite.T().Context(), &api.RequestActivity{
		Name:        "Updated Activity",
		Date:        "2025-02-02",
		Description: api.NilString{Null: true},
		Time:        api.NewNilString("12:34:56"),
		Address:     api.NilString{Null: true},
		Location:    api.NilEntityLocation{Null: true},
		Price:       api.NilInt{Null: true},
	}, api.PutActivityParams{TripID: tripID, ActivityID: activityID})
	suite.NoError(err)
	suite.IsType(&api.PutActivityNoContent{}, putRes)

	getRes, err := suite.api.GetActivity(suite.T().Context(), api.GetActivityParams{TripID: tripID, ActivityID: activityID})
	suite.NoError(err)
	activity, ok := getRes.(*api.EntityActivity)
	suite.True(ok)

	return *activity
}
