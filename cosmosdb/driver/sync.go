package driver

import (
	"context"

	metadatapb "github.com/bytebase/omni/metadata"

	"github.com/bytebase/omni/driver"
)

// SyncInstance syncs the instance meta: the account's databases, without their
// containers. Cosmos DB reports no version and no roles.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	names, err := d.client.ListDatabases(ctx)
	if err != nil {
		return nil, err
	}
	var databases []*metadatapb.DatabaseSchemaMetadata
	for _, name := range names {
		databases = append(databases, &metadatapb.DatabaseSchemaMetadata{
			Name: name,
		})
	}
	return &driver.Instance{
		Databases: databases,
	}, nil
}

// SyncDatabase syncs the database schema: one schema holding one table per
// container, with no columns, because a container's items are schemaless.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	names, err := d.client.ListContainers(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	var containers []*metadatapb.TableMetadata
	for _, name := range names {
		containers = append(containers, &metadatapb.TableMetadata{
			Name: name,
		})
	}
	return &metadatapb.DatabaseSchemaMetadata{
		Name: d.databaseName,
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Tables: containers,
			},
		},
	}, nil
}
