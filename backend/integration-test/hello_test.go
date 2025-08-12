package integration_test

import (
	"fmt"
	"testing"
	"travel-planner/integration-test/util"
)

func TestGreeterServer(t *testing.T) {
	util.StartDockerServer(t, "8080")
	fmt.Println("Hello World")
}
