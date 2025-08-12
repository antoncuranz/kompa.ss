package util

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

const (
	startupTimeout    = 5 * time.Second
	dockerfileName    = "Dockerfile"
	postgresAlias     = "postgres"
	dbName            = "postgres"
	dbUser            = "postgres"
	dbPassword        = "postgres"
	aerodataboxApiKey = "todo"
)

func StartDockerServer(
	t testing.TB,
	port string,
) {
	t.Helper()

	ctx := context.Background()
	logger := log.TestLogger(t)

	net, err := network.New(ctx)
	postgresContainer, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		network.WithNetwork([]string{postgresAlias}, net),
		testcontainers.WithLogger(logger),
	)
	assert.NoError(t, err)

	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", dbUser, dbPassword, postgresAlias, dbName)

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../.",
			Dockerfile: dockerfileName,
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, port)},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(startupTimeout),
		Env:          map[string]string{"PG_URL": dbConnectionString, "AERODATABOX_API_KEY": aerodataboxApiKey},
		Networks:     []string{net.Name},
	}
	applicationContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           logger,
	})

	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, applicationContainer.Terminate(ctx))
		assert.NoError(t, postgresContainer.Terminate(ctx))
	})
}
