package util

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/wiremock/go-wiremock"
	wiremockTc "github.com/wiremock/wiremock-testcontainers-go"
)

const (
	dbName     = "postgres"
	dbUser     = "postgres"
	dbPassword = "postgres"
)

func StartAllContainers(t testing.TB, wiremockRule *wiremock.StubRule) (*wiremock.Client, string, string) {
	t.Helper()

	wiremockClient, wiremockUrl := startWiremockContainer(t)
	err := wiremockClient.StubFor(wiremockRule)
	assert.NoError(t, err)

	dbConnectionString := startPostgresContainer(t)
	return wiremockClient, wiremockUrl, dbConnectionString
}

func startPostgresContainer(t testing.TB) string {
	t.Helper()

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
		testcontainers.WithLogger(tcLogger),
		//testcontainers.WithLogConsumers(testcontainers.LogConsumer(&containerLogger)),
	)
	assert.NoErrorf(t, err, "Failed to start postgres container: %v", err)

	cleanupContainer(t, postgresContainer)

	mappedPort, err := postgresContainer.MappedPort(t.Context(), "5432")
	assert.NoError(t, err, "Failed to get postgres mapped port: %v", err)

	return fmt.Sprintf("postgres://postgres:postgres@localhost:%s/postgres", mappedPort.Port())
}

func startWiremockContainer(t testing.TB) (*wiremock.Client, string) {
	t.Helper()

	tcLogger := TcLogger{t}
	containerLogger := ContainerLogger{
		containerName: "wiremock",
		colorPrefix:   ansiYellow,
	}

	wiremockContainer, err := wiremockTc.RunContainer(t.Context(),
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

	mappedPort, err := wiremockContainer.MappedPort(t.Context(), "8080")
	assert.NoError(t, err)

	wiremockURL := fmt.Sprintf("http://localhost:%s", mappedPort.Port())
	return wiremockContainer.Client, wiremockURL
}

func cleanupContainer(t testing.TB, container testcontainers.Container) {
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(context.Background()))
	})
}
