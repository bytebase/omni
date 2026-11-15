package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/catalog"
)

// Dump writes db as a mysqldump-style OceanBase DDL script, wrapped in the
// statements that disable and restore unique and foreign key checks.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := catalog.DatabaseDefinition(db, catalog.DefinitionOptions{PrintHeader: true})
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}
