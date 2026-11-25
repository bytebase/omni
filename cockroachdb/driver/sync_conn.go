package driver

import (
	"context"
	"database/sql"

	metadatapb "github.com/bytebase/omni/metadata"
)

// SyncDatabase reads the schema of databaseName through db, a connection to it
// the caller opened. The Driver method is this on its own pool; a caller that
// owns a connection, such as a catalog fidelity test, calls this.
func SyncDatabase(ctx context.Context, db *sql.DB, databaseName string) (*metadatapb.DatabaseSchemaMetadata, error) {
	d := &Driver{db: db, databaseName: databaseName}
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, err
	}
	d.engineVersion = version
	return d.SyncDatabase(ctx)
}
