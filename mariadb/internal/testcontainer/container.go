// Package testcontainer owns this engine's live-test server configuration.
package testcontainer

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	module "github.com/testcontainers/testcontainers-go/modules/mariadb"
)

const Image = "mariadb:11.8.8"

type Container = testcontainers.Container
type DatabaseContainer = module.MariaDBContainer

// Start starts the pinned server with the requested fixture database.
func Start(ctx context.Context, database string, args ...string) (*DatabaseContainer, error) {
	return StartImage(ctx, Image, database, args...)
}

// StartImage validates sweep overrides against the engine's selected version.
func StartImage(ctx context.Context, image, database string, args ...string) (*DatabaseContainer, error) {
	if image != Image {
		return nil, fmt.Errorf("test database image must be %s, got %q", Image, image)
	}
	return module.Run(ctx, Image,
		module.WithDatabase(database), module.WithUsername("root"), module.WithPassword("test"), testcontainers.WithCmdArgs(args...),
	)
}

func TerminateContainer(c Container) error { return testcontainers.TerminateContainer(c) }
