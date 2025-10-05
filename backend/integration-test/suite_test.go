package integration_test

import (
	"crypto/rsa"
	"fmt"
	"kompass/integration-test/client/api"
	"kompass/integration-test/util"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/wiremock/go-wiremock"
)

type IntegrationTestSuite struct {
	suite.Suite
	server     string
	api        *api.Client
	wiremock   *wiremock.Client
	privateKey *rsa.PrivateKey
}

const DefaultUser util.UserName = "Anton"
const ForbiddenUser util.UserName = "Forbidden"
const ReadingUser util.UserName = "Reader"
const WritingUser util.UserName = "Writer"

func (suite *IntegrationTestSuite) SetupSuite() {
	privateKey, jwkSetRule := util.GeneratePrivateKeyAndJwkStub(suite.T())
	suite.privateKey = privateKey

	port := "8081"
	wiremockClient := util.StartAllContainers(suite.T(), port, jwkSetRule)
	suite.wiremock = wiremockClient

	suite.server = fmt.Sprintf("http://127.0.0.1:%s/api/v1", port)
	suite.api = suite.userApi(DefaultUser)
}

func (suite *IntegrationTestSuite) SetupTest() {
	err := suite.wiremock.Reset()
	suite.NoError(err)
}

func (suite *IntegrationTestSuite) userApi(user util.UserName) *api.Client {
	app, err := api.NewClient(suite.server, util.GenerateJwtForUser(suite.T(), user, suite.privateKey))
	suite.NoError(err)
	return app
}

func (suite *IntegrationTestSuite) CreateTripUser(user util.UserName) int {
	res, err := suite.userApi(user).PostTrip(suite.T().Context(), &api.RequestTrip{
		Name:        "Test Trip",
		Description: api.NewNilString("This is a test"),
		StartDate:   "2025-01-01",
		EndDate:     "2026-01-01",
		ImageUrl:    api.NilString{},
	})
	suite.NoError(err)

	trip := res.(*api.EntityTrip)
	return trip.ID
}

func (suite *IntegrationTestSuite) CreateTrip() int {
	return suite.CreateTripUser(DefaultUser)
}

func (suite *IntegrationTestSuite) DeleteTripUser(user util.UserName, tripID int) {
	_, err := suite.userApi(user).DeleteTrip(suite.T().Context(), api.DeleteTripParams{
		TripID: tripID,
	})
	suite.NoError(err)
}

func (suite *IntegrationTestSuite) DeleteTrip(tripID int) {
	suite.DeleteTripUser(DefaultUser, tripID)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
