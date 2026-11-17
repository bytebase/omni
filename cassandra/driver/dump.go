package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
)

// Dump writes nothing and returns nil. The package doc says what a synced
// keyspace lacks that a restore would need.
func (*Driver) Dump(_ context.Context, _ io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	return nil
}
