package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/catalog"
)

// Dump writes db as the T-SQL script that recreates it. The catalog owns the
// rendering; the driver only passes the synced metadata through.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := catalog.DatabaseDefinition(db, catalog.DefinitionOptions{})
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}
