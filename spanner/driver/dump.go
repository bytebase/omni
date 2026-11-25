package driver

import (
	"context"
	"fmt"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
)

// Dump writes the database's DDL as the server states it, one statement per
// line, which is what Bytebase dumped. The metadata is not rendered and does not
// narrow the dump: spanner has no catalog that turns a synced schema back into
// DDL, and a Spanner database's own DDL is what the server holds.
func (d *Driver) Dump(ctx context.Context, out io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	statements, err := d.client.GetDatabaseDdl(ctx, d.databasePath(""))
	if err != nil {
		return err
	}
	for _, statement := range statements {
		if _, err := io.WriteString(out, fmt.Sprintf("%s;\n", statement)); err != nil {
			return err
		}
	}
	return nil
}
