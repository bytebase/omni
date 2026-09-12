// Package spannertest provides the Spanner emulator the googlesql differential
// tests run against.
package spannertest

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Image is the emulator build the recorded verdicts came from
// (sha256:caf1bd24…). Newer builds drift on a few forms, so it stays pinned.
const Image = "gcr.io/cloud-spanner-emulator/emulator:1.5.54"

var shared struct {
	once sync.Once
	host string
	err  error
}

// Host returns the emulator address: $SPANNER_EMULATOR_HOST when set, else an
// emulator testcontainer started once per test binary and exported through
// that variable so the harness subprocesses find it too. Docker is assumed,
// so a start failure fails the test.
func Host(t testing.TB) string {
	t.Helper()
	shared.once.Do(func() {
		if h := os.Getenv("SPANNER_EMULATOR_HOST"); h != "" {
			shared.host = h
			return
		}
		ctx := context.Background()
		c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        Image,
				ExposedPorts: []string{"9010/tcp"},
				WaitingFor:   wait.ForListeningPort("9010/tcp").WithStartupTimeout(2 * time.Minute),
			},
			Started: true,
		})
		if err != nil {
			shared.err = fmt.Errorf("start Spanner emulator: %w", err)
			return
		}
		host, err := c.Host(ctx)
		if err != nil {
			shared.err = err
			return
		}
		port, err := c.MappedPort(ctx, "9010/tcp")
		if err != nil {
			shared.err = err
			return
		}
		shared.host = fmt.Sprintf("%s:%s", host, port.Port())
		_ = os.Setenv("SPANNER_EMULATOR_HOST", shared.host)
	})
	if shared.err != nil {
		t.Fatalf("spanner emulator: %v", shared.err)
	}
	return shared.host
}
