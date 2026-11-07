// Package testcontainer owns this engine's live-test server configuration.
package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const Image = "dyrnq/doris:4.0.8"

type Container = testcontainers.Container

func Start(ctx context.Context) (Container, error) { return StartImage(ctx, Image) }
func StartImage(ctx context.Context, image string) (Container, error) {
	if image != Image {
		return nil, fmt.Errorf("test database image must be %s, got %q", Image, image)
	}
	return testcontainers.Run(ctx, Image,
		testcontainers.WithExposedPorts("9030/tcp"), testcontainers.WithEnv(map[string]string{"RUN_MODE": "STANDALONE", "SKIP_CHECK_ULIMIT": "true"}), testcontainers.WithWaitStrategy(wait.ForListeningPort("9030/tcp").WithStartupTimeout(5*time.Minute)),
	)
}
func TerminateContainer(c Container) error { return testcontainers.TerminateContainer(c) }
