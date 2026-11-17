package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
)

// Dump writes nothing. googlesql has no catalog that renders a synced schema back
// into DDL, and a BigQuery dataset's DDL is not what the sync read: a synced table
// carries its columns, comment, row count and data size, and not the partitioning,
// clustering, labels, expiration or view query a CREATE TABLE or CREATE VIEW
// states. Bytebase's BigQuery Dump wrote nothing either, so a caller that dumps a
// BigQuery database writes an empty dump, which is what the product does today.
// VERIFY.md says what a dump would need.
func (*Driver) Dump(_ context.Context, _ io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	return nil
}
