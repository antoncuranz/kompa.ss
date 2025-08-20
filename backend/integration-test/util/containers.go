package util

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/wiremock/go-wiremock"
	wiremockTc "github.com/wiremock/wiremock-testcontainers-go"
)

const (
	startupTimeout    = 5 * time.Second
	wiremockAlias     = "wiremock"
	postgresAlias     = "postgres"
	dbName            = "postgres"
	dbUser            = "postgres"
	dbPassword        = "postgres"
	aerodataboxApiKey = "todo"
)

func StartAllContainers(t testing.TB, applicationPort string) *wiremock.Client {
	t.Helper()

	net, err := network.New(t.Context())
	assert.NoError(t, err)

	dbConnectionString := startPostgresContainer(t, net)
	startApplicationContainer(t, dbConnectionString, applicationPort, net)
	return startWiremockContainer(t, net)
}

func startPostgresContainer(t testing.TB, net *testcontainers.DockerNetwork) string {
	logger := log.TestLogger(t)

	postgresContainer, err := postgres.Run(t.Context(), "postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		network.WithNetwork([]string{postgresAlias}, net),
		testcontainers.WithLogger(logger),
	)
	assert.NoError(t, err)

	cleanupContainer(t, postgresContainer)

	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPassword, postgresAlias, dbName)
}

func startWiremockContainer(t testing.TB, net *testcontainers.DockerNetwork) *wiremock.Client {
	logger := log.TestLogger(t)

	wiremockContainer, err := wiremockTc.RunContainer(t.Context(),
		network.WithNetwork([]string{wiremockAlias}, net),
		testcontainers.WithLogger(logger),
	)
	assert.NoError(t, err)

	cleanupContainer(t, wiremockContainer)

	return wiremockContainer.Client
}

func startApplicationContainer(t testing.TB, dbConnectionString string, port string, net *testcontainers.DockerNetwork) {
	logger := log.TestLogger(t)

	req := testcontainers.ContainerRequest{
		Image:        "kompa.ss/backend:latest",
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, "8080")},
		WaitingFor:   wait.ForListeningPort("8080").WithStartupTimeout(startupTimeout),
		Env: map[string]string{
			"PG_URL":              dbConnectionString,
			"AERODATABOX_API_KEY": aerodataboxApiKey,
			"AERODATABOX_URL":     "http://" + wiremockAlias + ":8080",
		},
		Networks: []string{net.Name},
	}
	applicationContainer, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           logger,
	})
	assert.NoError(t, err)

	cleanupContainer(t, applicationContainer)
}

func cleanupContainer(t testing.TB, container testcontainers.Container) {
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(context.Background()))
	})
}
