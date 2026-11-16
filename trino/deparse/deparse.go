package deparse

import (
	"fmt"
	"strings"
)

// GetDatabaseDefinition renders the schema metadata as the DDL that recreates
// it: one CREATE SCHEMA per schema, then one CREATE TABLE per table, then one
// CREATE VIEW per view. It is the omni-side implementation of the bytebase
// schema.GetDatabaseDefinition contract for the Trino engine.
//
// The output is in snapshot order within each kind, each statement separated by
// a blank line and terminated with the same trailing "\n\n" the bytebase source
// emits. Names are written always double-quoted, with an embedded double quote
// escaped as two. Column types are emitted as-is. Non-nullable columns get a
// trailing " NOT NULL".
//
// Schemas come first because a restore starts in a catalog that may not have
// them; tables come before views because a view may read a table in another
// schema. A nil metadata yields the empty string. A view with no definition is
// an error: writing it would silently lose the view.
func GetDatabaseDefinition(metadata *DatabaseSchemaMetadata) (string, error) {
	if metadata == nil {
		return "", nil
	}

	var buf strings.Builder
	for _, schema := range metadata.Schemas {
		if schema == nil {
			continue
		}
		writeCreateSchema(&buf, schema.Name)
		buf.WriteString("\n\n")
	}
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
	for _, schema := range metadata.Schemas {
		if schema == nil {
			continue
		}
		for _, view := range schema.Views {
			if view == nil {
				continue
			}
			if err := writeCreateView(&buf, schema.Name, view); err != nil {
				return "", err
			}
			buf.WriteString("\n\n")
		}
	}

	return buf.String(), nil
}

// quoteIdentifier writes name as a double-quoted Trino identifier. A double
// quote inside name is doubled, as Trino's own lexer reads it back.
func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// writeCreateSchema appends the CREATE SCHEMA statement for schemaName. It does
// not append the inter-statement blank line; the caller does.
func writeCreateSchema(buf *strings.Builder, schemaName string) {
	fmt.Fprintf(buf, "CREATE SCHEMA IF NOT EXISTS %s;", quoteIdentifier(schemaName))
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
	fmt.Fprintf(buf, "CREATE TABLE IF NOT EXISTS %s.%s (\n", quoteIdentifier(schemaName), quoteIdentifier(table.Name))

	columnDefs := make([]string, 0, len(table.Columns))
	for _, column := range table.Columns {
		if column == nil {
			continue
		}
		nullable := ""
		if !column.Nullable {
			nullable = " NOT NULL"
		}
		columnDefs = append(columnDefs, fmt.Sprintf("    %s %s%s", quoteIdentifier(column.Name), column.Type, nullable))
	}

	buf.WriteString(strings.Join(columnDefs, ",\n"))
	buf.WriteString("\n);")
}

// writeCreateView appends a single CREATE OR REPLACE VIEW statement for view
// (qualified by schemaName) to buf. It does not append the inter-statement
// blank line; the caller does. OR REPLACE matches CREATE TABLE IF NOT EXISTS:
// a dump run twice restores rather than fails.
func writeCreateView(buf *strings.Builder, schemaName string, view *ViewMetadata) error {
	definition := strings.TrimSpace(view.Definition)
	if definition == "" {
		return fmt.Errorf("view %s.%s has no definition to restore", schemaName, view.Name)
	}
	fmt.Fprintf(buf, "CREATE OR REPLACE VIEW %s.%s AS\n%s;", quoteIdentifier(schemaName), quoteIdentifier(view.Name), definition)
	return nil
}
