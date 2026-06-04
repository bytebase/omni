package deparse

import (
	"fmt"
	"strings"
)

// GetDatabaseDefinition renders the schema metadata as the CREATE TABLE SDL
// that recreates it. It is the omni-side implementation of the bytebase
// schema.GetDatabaseDefinition contract for the Trino engine.
//
// The output is one CREATE TABLE statement per table, in schema-then-table
// snapshot order, each separated by a blank line and terminated with the same
// trailing "\n\n" the bytebase source emits (so the generated dump is
// byte-identical to the legacy implementation). Names are written verbatim,
// always double-quoted; column types are emitted as-is. Non-nullable columns
// get a trailing " NOT NULL".
//
// A nil metadata, or a metadata with no tables, yields the empty string. The
// (string, error) signature mirrors the bytebase schema.GetDatabaseDefinition
// contract the bytebase-switch node adapts; the error is always nil today, but
// callers must still check it.
func GetDatabaseDefinition(metadata *DatabaseSchemaMetadata) (string, error) {
	if metadata == nil {
		return "", nil
	}

	var buf strings.Builder
	for _, schema := range metadata.Schemas {
		if schema == nil {
			continue
		}
		for _, table := range schema.Tables {
			if table == nil {
				continue
			}
			writeCreateTable(&buf, schema.Name, table)
			buf.WriteString("\n\n")
		}
	}

	return buf.String(), nil
}

// writeCreateTable appends a single CREATE TABLE statement for table (qualified
// by schemaName) to buf. It does not append the inter-statement blank line; the
// caller does. The emitted form is:
//
//	CREATE TABLE IF NOT EXISTS "schema"."table" (
//	    "col1" type1 NOT NULL,
//	    "col2" type2
//	);
//
// matching the legacy bytebase plugin/schema/trino output exactly, including the
// four-space column indent and the leading newline before the closing paren
// (which, for a table with no columns, leaves a blank line inside the parens).
func writeCreateTable(buf *strings.Builder, schemaName string, table *TableMetadata) {
	fmt.Fprintf(buf, "CREATE TABLE IF NOT EXISTS \"%s\".\"%s\" (\n", schemaName, table.Name)

	columnDefs := make([]string, 0, len(table.Columns))
	for _, column := range table.Columns {
		if column == nil {
			continue
		}
		nullable := ""
		if !column.Nullable {
			nullable = " NOT NULL"
		}
		columnDefs = append(columnDefs, fmt.Sprintf("    \"%s\" %s%s", column.Name, column.Type, nullable))
	}

	buf.WriteString(strings.Join(columnDefs, ",\n"))
	buf.WriteString("\n);")
}
