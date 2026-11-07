// Package testcontainer owns the Spanner emulator used by this harness.
package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Image matches googlesql/internal/testcontainer.Image.
const Image = "gcr.io/cloud-spanner-emulator/emulator:1.5.54"

func Start(ctx context.Context) (testcontainers.Container, string, error) {
	c, err := testcontainers.Run(ctx, Image,
		testcontainers.WithExposedPorts("9010/tcp"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("9010/tcp").WithStartupTimeout(2*time.Minute)),
	)
	if err != nil {
		return c, "", fmt.Errorf("start Spanner emulator: %w", err)
	}
	host, err := c.Host(ctx)
	if err != nil {
		return c, "", err
	}
	port, err := c.MappedPort(ctx, "9010/tcp")
	if err != nil {
		return c, "", err
	}
	return c, host + ":" + port.Port(), nil
}

func TerminateContainer(c testcontainers.Container) error {
	return testcontainers.TerminateContainer(c)
}
