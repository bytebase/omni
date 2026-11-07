// Package testcontainer owns this engine's live-test server configuration.
package testcontainer

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	module "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const Image = "postgres:18.6-alpine"

type Container = testcontainers.Container
type DatabaseContainer = module.PostgresContainer

// Start starts the pinned server with the requested fixture database.
func Start(ctx context.Context, database string) (*DatabaseContainer, error) {
	return StartImage(ctx, Image, database)
}

// StartImage validates sweep overrides against the engine's selected version.
func StartImage(ctx context.Context, image, database string) (*DatabaseContainer, error) {
	if image != Image {
		return nil, fmt.Errorf("test database image must be %s, got %q", Image, image)
	}
	return module.Run(ctx, Image,
		module.WithDatabase(database), module.WithUsername("postgres"), module.WithPassword("test"), testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
}

func TerminateContainer(c Container) error { return testcontainers.TerminateContainer(c) }
