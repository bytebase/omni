// Package testcontainer owns this engine's live-test server configuration.
package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const Image = "starrocks/allin1-ubuntu:4.1.7"

type Container = testcontainers.Container

func Start(ctx context.Context) (Container, error) { return StartImage(ctx, Image) }
func StartImage(ctx context.Context, image string) (Container, error) {
	if image != Image {
		return nil, fmt.Errorf("test database image must be %s, got %q", Image, image)
	}
	return testcontainers.Run(ctx, Image,
		testcontainers.WithImagePlatform("linux/amd64"), testcontainers.WithExposedPorts("9030/tcp"), testcontainers.WithWaitStrategy(wait.ForListeningPort("9030/tcp").WithStartupTimeout(4*time.Minute)),
	)
}
func TerminateContainer(c Container) error { return testcontainers.TerminateContainer(c) }
