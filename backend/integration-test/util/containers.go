package util

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/wiremock/go-wiremock"
	wiremockTc "github.com/wiremock/wiremock-testcontainers-go"
)

const (
	startupTimeout = 5 * time.Second
	wiremockAlias  = "wiremock"
	postgresAlias  = "postgres"
	dbName         = "postgres"
	dbUser         = "postgres"
	dbPassword     = "postgres"
)

func StartAllContainers(t testing.TB, applicationPort string, wiremockRule *wiremock.StubRule) *wiremock.Client {
	t.Helper()

	net, err := network.New(t.Context())
	assert.NoError(t, err)

	wiremockClient := startWiremockContainer(t, net)
	err = wiremockClient.StubFor(wiremockRule)
	assert.NoError(t, err)

	dbConnectionString := startPostgresContainer(t, net)
	startApplicationContainer(t, dbConnectionString, applicationPort, net)
	return wiremockClient
}

func startPostgresContainer(t testing.TB, net *testcontainers.DockerNetwork) string {
	tcLogger := TcLogger{t}
	//containerLogger := ContainerLogger{
	//	containerName: "postgres",
	//	colorPrefix:   ansiBlue,
	//}

	postgresContainer, err := postgres.Run(t.Context(), "postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		network.WithNetwork([]string{postgresAlias}, net),
		testcontainers.WithLogger(tcLogger),
		//testcontainers.WithLogConsumers(testcontainers.LogConsumer(&containerLogger)),
	)
	assert.NoError(t, err)

	cleanupContainer(t, postgresContainer)
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPassword, postgresAlias, dbName)
}

func startWiremockContainer(t testing.TB, net *testcontainers.DockerNetwork) *wiremock.Client {
	tcLogger := TcLogger{t}
	containerLogger := ContainerLogger{
		containerName: "wiremock",
		colorPrefix:   ansiYellow,
	}

	wiremockContainer, err := wiremockTc.RunContainer(t.Context(),
		network.WithNetwork([]string{wiremockAlias}, net),
		testcontainers.WithImage("docker.io/wiremock/wiremock:3.13.1"),
		testcontainers.WithLogger(tcLogger),
		testcontainers.WithLogConsumers(testcontainers.LogConsumer(&containerLogger)),
		testcontainers.WithEnv(map[string]string{
			"WIREMOCK_OPTIONS": "--disable-banner --verbose",
		}),
		testcontainers.WithFiles(testcontainers.ContainerFile{
			HostFilePath:      "wiremock/",
			ContainerFilePath: "/home/",
			FileMode:          0755,
		}),
	)
	assert.NoError(t, err)

	cleanupContainer(t, wiremockContainer)
	return wiremockContainer.Client
}

func startApplicationContainer(t testing.TB, dbConnectionString string, port string, net *testcontainers.DockerNetwork) {
	tcLogger := TcLogger{t}
	containerLogger := ContainerLogger{
		containerName: "kompass",
		colorPrefix:   ansiGreen,
	}

	wiremockUrl := "http://" + wiremockAlias + ":8080"

	req := testcontainers.ContainerRequest{
		Image:        "kompa.ss/backend:latest",
		ExposedPorts: []string{fmt.Sprintf("%s:%s", port, "8080")},
		WaitingFor:   wait.ForListeningPort("8080").WithStartupTimeout(startupTimeout),
		Env: map[string]string{
			"PG_URL":        dbConnectionString,
			"AUTH_JWKS_URL": wiremockUrl + "/auth/jwks.json",
			"AEDBX_URL":     wiremockUrl + "/aedbx",
			"DBVENDO_URL":   wiremockUrl + "/dbvendo",
		},
		Networks: []string{net.Name},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{&containerLogger},
		},
	}
	applicationContainer, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           tcLogger,
	})
	assert.NoError(t, err)

	cleanupContainer(t, applicationContainer)
}

func cleanupContainer(t testing.TB, container testcontainers.Container) {
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(context.Background()))
	})
}
