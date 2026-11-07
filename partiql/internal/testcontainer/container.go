// Package testcontainer owns the database image and lifecycle used by this engine's tests.
package testcontainer

import (
	"context"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const Image = "amazon/dynamodb-local:3.3.1"

type Container = testcontainers.Container

var server struct {
	once      sync.Once
	container Container
	err       error
}

// Start returns the process's shared native test server, starting it on first use.
func Start(ctx context.Context) (Container, error) {
	server.once.Do(func() {
		server.container, server.err = testcontainers.Run(ctx, Image,
			testcontainers.WithCmdArgs("-jar", "DynamoDBLocal.jar", "-inMemory", "-sharedDb"),
			testcontainers.WithExposedPorts("8000/tcp"),
			testcontainers.WithWaitStrategy(wait.ForListeningPort("8000/tcp").WithStartupTimeout(time.Minute)),
		)
	})
	return server.container, server.err
}

// TerminateContainer releases a server at the end of the test process.
func TerminateContainer(c Container) error { return testcontainers.TerminateContainer(c) }
