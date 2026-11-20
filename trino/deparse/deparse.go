package deparse

import (
	"fmt"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/analysis"
)

// DatabaseDefinitionFromProto renders a Trino database schema proto as SDL.
func DatabaseDefinitionFromProto(metadata *metadatapb.DatabaseSchemaMetadata) (string, error) {
	return GetDatabaseDefinition(fromProtoMetadata(metadata))
}

func fromProtoMetadata(metadata *metadatapb.DatabaseSchemaMetadata) *DatabaseSchemaMetadata {
	if metadata == nil {
		return nil
	}
	out := &DatabaseSchemaMetadata{
		Name:    metadata.Name,
		Schemas: make([]*SchemaMetadata, 0, len(metadata.Schemas)),
	}
	for _, s := range metadata.Schemas {
		if s == nil || s.SkipDump {
			continue
		}
		schemaMeta := &SchemaMetadata{
			Name:   s.Name,
			Tables: make([]*TableMetadata, 0, len(s.Tables)),
			Views:  make([]*ViewMetadata, 0, len(s.Views)),
		}
		for _, t := range s.Tables {
			if t == nil || t.SkipDump {
				continue
			}
			tableMeta := &TableMetadata{
				Name:    t.Name,
				Columns: make([]*ColumnMetadata, 0, len(t.Columns)),
			}
			for _, c := range t.Columns {
				if c == nil {
					continue
				}
				tableMeta.Columns = append(tableMeta.Columns, &ColumnMetadata{
					Name:     c.Name,
					Type:     c.Type,
					Nullable: c.Nullable,
				})
			}
			schemaMeta.Tables = append(schemaMeta.Tables, tableMeta)
		}
		for _, v := range s.Views {
			if v == nil || v.SkipDump {
				continue
			}
			schemaMeta.Views = append(schemaMeta.Views, &ViewMetadata{
				Name:       v.Name,
				Definition: v.Definition,
			})
		}
		out.Schemas = append(out.Schemas, schemaMeta)
	}
	return out
}

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
// schema. Views are ordered by their definition's references to other views, so
// a view that reads another is written second. A nil metadata yields the empty
// string. A view with no definition is an error: writing it would silently lose
// the view.
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
	for _, view := range orderedViews(metadata) {
		if err := writeCreateView(&buf, view.schema, view.view); err != nil {
			return "", err
		}
		buf.WriteString("\n\n")
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

// orderedView is one view and the schema that owns it, in snapshot order.
type orderedView struct {
	schema string
	view   *ViewMetadata
}

// orderedViews returns metadata's views in dependency order: every view comes
// after the views its definition reads. Independent views keep snapshot order.
// A definition the local parser cannot read contributes no edges rather than
// failing the dump; Trino accepted it, and a dependency the local parser cannot
// see is no reason to lose the statement.
func orderedViews(metadata *DatabaseSchemaMetadata) []orderedView {
	var views []orderedView
	for _, schema := range metadata.Schemas {
		if schema == nil {
			continue
		}
		for _, view := range schema.Views {
			if view == nil {
				continue
			}
			views = append(views, orderedView{schema: schema.Name, view: view})
		}
	}
	if len(views) < 2 {
		return views
	}

	index := make(map[string]int, len(views))
	for i, view := range views {
		index[viewKey(view.schema, view.view.Name)] = i
	}
	dependencies := make([]map[int]struct{}, len(views))
	for i, view := range views {
		span, err := analysis.GetQuerySpan(view.view.Definition)
		if err != nil {
			continue
		}
		for _, access := range span.AccessTables {
			if access.Catalog != "" && metadata.Name != "" && access.Catalog != metadata.Name {
				continue
			}
			schema := access.Schema
			if schema == "" {
				schema = view.schema
			}
			dependency, ok := index[viewKey(schema, access.Table)]
			if !ok || dependency == i {
				continue
			}
			if dependencies[i] == nil {
				dependencies[i] = map[int]struct{}{}
			}
			dependencies[i][dependency] = struct{}{}
		}
	}
	return sortViews(views, dependencies)
}

// viewKey identifies a view by its owning schema and its name.
func viewKey(schema, name string) string {
	return schema + "\x00" + name
}

// sortViews orders views so each follows its dependencies, keeping the input
// order among views that do not depend on each other. A cycle, which Trino does
// not let a CREATE VIEW make, leaves the remaining views in input order.
func sortViews(views []orderedView, dependencies []map[int]struct{}) []orderedView {
	remaining := make([]int, len(views))
	for i, deps := range dependencies {
		remaining[i] = len(deps)
	}
	ordered := make([]orderedView, 0, len(views))
	done := make([]bool, len(views))
	for len(ordered) < len(views) {
		progressed := false
		for i, view := range views {
			if done[i] || remaining[i] != 0 {
				continue
			}
			done[i] = true
			ordered = append(ordered, view)
			progressed = true
			for j, deps := range dependencies {
				if _, ok := deps[i]; ok && !done[j] {
					remaining[j]--
				}
			}
		}
		if !progressed {
			for i, view := range views {
				if !done[i] {
					done[i] = true
					ordered = append(ordered, view)
				}
			}
		}
	}
	return ordered
}
