package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/catalog"
)

// Dump writes db as the DDL that creates it, after a header that says which
// tool wrote it.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := catalog.DatabaseDefinition(db, catalog.DefinitionOptions{PrintHeader: true})
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}
