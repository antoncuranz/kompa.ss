package integration_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"kompass/integration-test/client"
	"kompass/integration-test/util"
	"testing"
)

func startApplicationAndBuildClient(t *testing.T) *client.ClientWithResponses {
	port := "8080"
	util.StartDockerServer(t, port)

	server := fmt.Sprintf("http://127.0.0.1:%s/api/v1", port)
	app, err := client.NewClientWithResponses(server)
	assert.NoError(t, err)

	return app
}

func TestGetTrips(t *testing.T) {
	// given
	app := startApplicationAndBuildClient(t)

	// when
	trips, err := app.GetTripsWithResponse(context.Background())

	// then
	assert.NoError(t, err)
	assert.Equal(t, 200, trips.HTTPResponse.StatusCode)
	assert.Len(t, *trips.JSON200, 1)
}
