package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
	mongocatalog "github.com/bytebase/omni/mongo/catalog"
)

// Dump writes db as a mongosh script that recreates its collections, their
// indexes and its views. A collection is created with the options its metadata
// holds, an index with its key and options, and a view with its definition; a
// view whose metadata predates a synced definition is named in a comment
// instead, since a view cannot be recreated without one.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := mongocatalog.DatabaseDefinition(db)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}
