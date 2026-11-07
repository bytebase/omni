// Package testcontainer owns this engine's live-test server configuration.
package testcontainer

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const Image = "gvenzl/oracle-free:23.26.3-slim-faststart"

type Container = testcontainers.Container

func Start(ctx context.Context) (Container, error) { return StartImage(ctx, Image) }
func StartImage(ctx context.Context, image string) (Container, error) {
	if image != Image {
		return nil, fmt.Errorf("test database image must be %s, got %q", Image, image)
	}
	return testcontainers.Run(ctx, Image,
		testcontainers.WithExposedPorts("1521/tcp"), testcontainers.WithEnv(map[string]string{"ORACLE_PASSWORD": "testpass", "APP_USER": "testuser", "APP_USER_PASSWORD": "testpass"}), testcontainers.WithWaitStrategy(wait.ForLog("DATABASE IS READY TO USE!").WithStartupTimeout(5*time.Minute)),
	)
}
func TerminateContainer(c Container) error { return testcontainers.TerminateContainer(c) }
