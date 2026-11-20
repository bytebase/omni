package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/deparse"
)

// Dump writes db as the DDL that recreates it.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := deparse.DatabaseDefinitionFromProto(db)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}
