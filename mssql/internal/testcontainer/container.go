// Package testcontainer owns the SQL Server image and startup settings for tests.
package testcontainer

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	module "github.com/testcontainers/testcontainers-go/modules/mssql"
)

const Image = "mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04"

type Container = testcontainers.Container
type DatabaseContainer = module.MSSQLServerContainer

func Start(ctx context.Context) (*DatabaseContainer, error) {
	return module.Run(ctx, Image, module.WithAcceptEULA(), module.WithPassword("Str0ngPa$$w0rd!"))
}
func TerminateContainer(c Container) error { return testcontainers.TerminateContainer(c) }
