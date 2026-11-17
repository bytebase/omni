package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/catalog"
)

// Dump writes db as the DDL that recreates its objects, in dependency order:
// sequences before the tables whose defaults call them, tables before the views
// that read them, and functions before the indexes that call them. Everything
// comes from db and no server is read; what the snapshot cannot carry is left to
// the catalog's documented gaps, which docs/migration/oracle/catalog.md names.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := catalog.DatabaseDefinition(db)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}
