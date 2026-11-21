package analysis

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/parser"
)

// splitSpanStatements splits a script respecting Trino's statement boundaries
// and reattaches a trailing semicolon, as the Bytebase statement splitter does.
func splitSpanStatements(statement string) ([]spanStatement, error) {
	segs := parser.Split(statement)
	if len(segs) == 0 {
		return nil, nil
	}
	stmts := make([]spanStatement, 0, len(segs))
	for _, seg := range segs {
		text := seg.Text
		if seg.ByteEnd < len(statement) && statement[seg.ByteEnd] == ';' {
			text = statement[seg.ByteStart : seg.ByteEnd+1]
		}
		stmts = append(stmts, spanStatement{Text: text, Empty: seg.Empty()})
	}
	return stmts, nil
}

type spanStatement struct {
	Text  string
	Empty bool
}

func findSpanSchema(meta *metadata.DatabaseSchemaMetadata, name string) *metadata.SchemaMetadata {
	for _, schema := range meta.GetSchemas() {
		if schema.GetName() == name {
			return schema
		}
	}
	return nil
}

func listSpanSchemaNames(meta *metadata.DatabaseSchemaMetadata) []string {
	names := make([]string, 0, len(meta.GetSchemas()))
	for _, schema := range meta.GetSchemas() {
		names = append(names, schema.GetName())
	}
	return names
}

func findSpanTable(schema *metadata.SchemaMetadata, name string) *metadata.TableMetadata {
	for _, table := range schema.GetTables() {
		if table.GetName() == name {
			return table
		}
	}
	return nil
}

func listSpanTableNames(schema *metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schema.GetTables()))
	for _, table := range schema.GetTables() {
		names = append(names, table.GetName())
	}
	return names
}

func findSpanView(schema *metadata.SchemaMetadata, name string) *metadata.ViewMetadata {
	for _, view := range schema.GetViews() {
		if view.GetName() == name {
			return view
		}
	}
	return nil
}

func listSpanViewNames(schema *metadata.SchemaMetadata) []string {
	names := make([]string, 0, len(schema.GetViews()))
	for _, view := range schema.GetViews() {
		names = append(names, view.GetName())
	}
	return names
}

func viewDefinitions(meta *metadata.DatabaseSchemaMetadata) []string {
	var definitions []string
	for _, schema := range meta.GetSchemas() {
		for _, view := range schema.GetViews() {
			if def := view.GetDefinition(); def != "" {
				definitions = append(definitions, def)
			}
		}
	}
	return definitions
}

func catalogNeeded(normName, defaultDB, lowerStmt string) bool {
	if normName == "" {
		return false
	}
	if defaultDB != "" && normName == defaultDB {
		return true
	}
	if strings.Contains(normName, `"`) {
		return true
	}
	return strings.Contains(lowerStmt, strings.ToLower(normName))
}

func spanContainsSystemSchema(sql string) bool {
	lowerSQL := strings.ToLower(sql)
	for _, prefix := range []string{
		"system.",
		"information_schema.",
		"$system.",
		"catalog.",
		"metadata.",
	} {
		if strings.Contains(lowerSQL, prefix) {
			return true
		}
	}
	return false
}
