package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/deparse"
)

// Dump writes db as the CREATE TABLE statements that recreate it. The
// deparser emits names and types verbatim in double quotes, so a name holding a
// double quote does not restore; see TestDumpGaps.
func (*Driver) Dump(_ context.Context, w io.Writer, db *metadatapb.DatabaseSchemaMetadata) error {
	text, err := deparse.GetDatabaseDefinition(convertDatabaseMetadata(db))
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, text)
	return err
}

// convertDatabaseMetadata maps the synced metadata onto the deparser's own
// model, which keeps the fields the Trino SDL dump reads (schemas, tables and
// column Name/Type/Nullable) and no others.
func convertDatabaseMetadata(metadata *metadatapb.DatabaseSchemaMetadata) *deparse.DatabaseSchemaMetadata {
	if metadata == nil {
		return nil
	}
	out := &deparse.DatabaseSchemaMetadata{
		Name:    metadata.Name,
		Schemas: make([]*deparse.SchemaMetadata, 0, len(metadata.Schemas)),
	}
	for _, s := range metadata.Schemas {
		if s == nil {
			continue
		}
		schema := &deparse.SchemaMetadata{
			Name:   s.Name,
			Tables: make([]*deparse.TableMetadata, 0, len(s.Tables)),
		}
		for _, t := range s.Tables {
			if t == nil {
				continue
			}
			table := &deparse.TableMetadata{
				Name:    t.Name,
				Columns: make([]*deparse.ColumnMetadata, 0, len(t.Columns)),
			}
			for _, c := range t.Columns {
				if c == nil {
					continue
				}
				table.Columns = append(table.Columns, &deparse.ColumnMetadata{
					Name:     c.Name,
					Type:     c.Type,
					Nullable: c.Nullable,
				})
			}
			schema.Tables = append(schema.Tables, table)
		}
		out.Schemas = append(out.Schemas, schema)
	}
	return out
}
