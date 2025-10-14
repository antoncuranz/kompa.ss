package util

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BackendProcess struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func StartBackendSubprocess(t testing.TB, dbConnectionString string, wiremockURL string, port string) *exec.Cmd {
	t.Helper()

	buildCmd := exec.Command("go", "build", "-o", "/tmp/kompass-backend", "../cmd/app")
	err := buildCmd.Run()
	require.NoErrorf(t, err, "Failed to build backend binary")

	cmd := exec.CommandContext(t.Context(), "/tmp/kompass-backend")
	cmd.Dir = "../" // Set working directory to backend root so migrations can be found
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("HTTP_PORT=%s", port),
		fmt.Sprintf("PG_URL=%s", dbConnectionString),
		fmt.Sprintf("AUTH_JWKS_URL=%s/auth/jwks.json", wiremockURL),
		fmt.Sprintf("AEDBX_URL=%s/aedbx", wiremockURL),
		fmt.Sprintf("DBVENDO_URL=%s/dbvendo", wiremockURL),
		fmt.Sprintf("ORS_URL=%s/ors", wiremockURL),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	require.NoErrorf(t, err, "Failed to start backend subprocess")

	backendURL := fmt.Sprintf("http://localhost:%s", port)
	waitForBackendReady(t, backendURL)

	t.Cleanup(func() {
		err := cmd.Process.Kill()
		assert.NoError(t, err)
		err = cmd.Wait()
		assert.NoError(t, err)
	})

	return cmd
}

func waitForBackendReady(t testing.TB, backendURL string) {
	t.Helper()

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Backend failed to start within timeout")
		case <-ticker.C:
			cmd := exec.Command("curl", "-f", "-s", fmt.Sprintf("%s/healthz", backendURL))
			if err := cmd.Run(); err == nil {
				return
			}
		}
	}
}
