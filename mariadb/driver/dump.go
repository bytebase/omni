package driver

import (
	"context"
	"io"

	"github.com/bytebase/omni/mariadb/catalog"
	metadatapb "github.com/bytebase/omni/metadata"
)

// Dump writes db as a MariaDB DDL script, wrapped in the statements that
// disable and restore unique and foreign key checks.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := catalog.DatabaseDefinition(db, catalog.DefinitionOptions{PrintHeader: true})
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}
