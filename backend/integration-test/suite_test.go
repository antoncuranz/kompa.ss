package integration_test

import (
	"fmt"
	"kompass/integration-test/client/api"
	"kompass/integration-test/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/wiremock/go-wiremock"
)

type IntegrationTestSuite struct {
	suite.Suite
	api      *api.Client
	wiremock *wiremock.Client
}

func (suite *IntegrationTestSuite) SetupSuite() {
	port := "8081"
	wiremockClient := util.StartAllContainers(suite.T(), port)

	server := fmt.Sprintf("http://127.0.0.1:%s/api/v1", port)
	app, err := api.NewClient(server)
	assert.NoError(suite.T(), err)

	suite.api = app
	suite.wiremock = wiremockClient
}

func (suite *IntegrationTestSuite) CreateTrip() int {
	res, err := suite.api.PostTrip(suite.T().Context(), &api.RequestTrip{
		Name:        "Test Trip",
		Description: api.NewNilString("This is a test"),
		StartDate:   "2025-01-01",
		EndDate:     "2026-01-01",
		ImageUrl:    api.NilString{},
	})
	suite.NoError(err)

	trip := res.(*api.ResponseTrip)
	return trip.ID
}

func (suite *IntegrationTestSuite) DeleteTrip(tripID int) {
	_, err := suite.api.DeleteTrip(suite.T().Context(), api.DeleteTripParams{
		TripID: tripID,
	})
	suite.NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
