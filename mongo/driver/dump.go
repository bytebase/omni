package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
	mongocatalog "github.com/bytebase/omni/mongo/catalog"
)

// Dump writes db as a mongosh script that recreates its collections and their
// indexes. A view is not restorable from the synced metadata, which carries no
// pipeline, so the script names it in a comment instead.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := mongocatalog.DatabaseDefinition(db)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}
