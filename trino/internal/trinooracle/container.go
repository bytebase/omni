package trinooracle

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// DefaultImage is the Trino image used by StartContainer. Pinned so
// differential results are reproducible — a floating `latest` tag lets new
// Trino releases silently shift the accept/reject oracle. Bump deliberately
// and re-run the differentials when moving to a new Trino version.
const DefaultImage = "trinodb/trino:482"

// StartContainer starts a disposable Trino server via testcontainers and returns
// an Oracle pointed at it plus a terminate func. Trino's JVM takes ~30-60s to
// become ready, so for iterative local testing prefer a long-lived container and
// Connect(URLFromEnv()); use StartContainer for self-contained CI runs.
func StartContainer(ctx context.Context, opts ...Option) (oracle *Oracle, terminate func(), err error) {
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        DefaultImage,
			ExposedPorts: []string{"8080/tcp"},
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, err
	}
	terminate = func() { _ = testcontainers.TerminateContainer(c) }

	host, err := c.Host(ctx)
	if err != nil {
		terminate()
		return nil, nil, err
	}
	port, err := c.MappedPort(ctx, "8080/tcp")
	if err != nil {
		terminate()
		return nil, nil, err
	}
	o := Connect(fmt.Sprintf("http://%s:%s", host, port.Port()), opts...)

	// GenericContainer returns once the process is up; poll until Trino reports
	// it has finished starting.
	deadline := time.Now().Add(3 * time.Minute)
	for {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, perr := o.Ping(pingCtx)
		cancel()
		if perr == nil {
			return o, terminate, nil
		}
		if time.Now().After(deadline) {
			terminate()
			return nil, nil, fmt.Errorf("trino did not become ready within timeout: %w", perr)
		}
		select {
		case <-ctx.Done():
			terminate()
			return nil, nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}
