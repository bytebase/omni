package v1review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// defaultSchema is the schema an index's table is looked up in when the
// statement does not name one.
const defaultSchema = "dbo"

// IndexNotRedundant reports each top-level CREATE INDEX whose first key
// column is the first expression of an index the schema snapshot already
// holds on the same table. Statements nested in a routine body or an IF or
// BEGIN block are not judged, nor is a CREATE INDEX on a table qualified
// with a database other than currentDatabase, compared without regard to
// case. Only the first key column is compared, and only with an existing
// index's first expression, as text and case sensitively; the rest of
// either column list, uniqueness, clustering, INCLUDE columns and filters
// are ignored. A first key column that is an expression is not judged.
//
// The table is looked up in schema.Metadata directly, not through Schema's
// lookups: the schema is the statement's, dbo when it names none, and the
// schema and table names must equal the snapshot's exactly, whatever
// ObjectCaseSensitive says. A schema or table with an empty name in the
// snapshot is never found, a table's partitions do not find it, and where
// several schemas and tables share both names the last table in the
// snapshot's order is used. A nil Metadata finds no table. As in the V1
// advisor, every schema and table name is read before any statement is,
// without nil checks, so a nil schema, or a nil table in a named schema,
// panics whatever the statements; a nil index panics when a matching
// statement reaches it.
//
// IndexRedundant is the code, one finding per statement at most: Name is
// the statement's index name, Other the first existing index that matches,
// in the table's order, and Schema and Table the names looked up. The
// finding anchors the whole statement and sets Statement, since the
// snapshot is read once for the whole script. Ported from Bytebase's V1 rule
// INDEX_NOT_REDUNDANT.
func IndexNotRedundant(stmts []ast.Node, schema v1review.Schema, currentDatabase string) []v1review.Finding {
	tables := tablesByName(schema.Metadata)
	var findings []v1review.Finding
	for i, stmt := range stmts {
		ci, ok := stmt.(*ast.CreateIndexStmt)
		if !ok {
			continue
		}
		schemaName, tableName := defaultSchema, ""
		if ci.Table != nil {
			if ci.Table.Database != "" && !strings.EqualFold(ci.Table.Database, currentDatabase) {
				continue
			}
			if ci.Table.Schema != "" {
				schemaName = ci.Table.Schema
			}
			tableName = ci.Table.Object
		}
		table, ok := tables[tableKey{schemaName, tableName}]
		if !ok {
			continue
		}
		if ci.Columns == nil || ci.Columns.Len() == 0 {
			continue
		}
		firstColumn := ""
		if column, ok := ci.Columns.Items[0].(*ast.IndexColumn); ok && column != nil {
			firstColumn = column.Name
		}
		if firstColumn == "" {
			continue
		}
		for _, index := range table.Indexes {
			if len(index.Expressions) > 0 && index.Expressions[0] == firstColumn {
				start, end := span(ci.Loc)
				findings = append(findings, v1review.Finding{
					Code:      v1review.IndexRedundant,
					Name:      ci.Name,
					Other:     index.Name,
					Schema:    schemaName,
					Table:     tableName,
					Start:     start,
					End:       end,
					Statement: i,
				})
				break
			}
		}
	}
	return findings
}

// tableKey is a table's schema and table names, as written.
type tableKey struct{ schema, table string }

// tablesByName returns the snapshot's tables keyed by their exact names,
// skipping a schema or table with an empty name; a later table replaces an
// earlier one of the same key.
func tablesByName(m *metadata.DatabaseSchemaMetadata) map[tableKey]*metadata.TableMetadata {
	tables := make(map[tableKey]*metadata.TableMetadata)
	if m == nil {
		return tables
	}
	for _, schema := range m.Schemas {
		if schema.Name == "" {
			continue
		}
		for _, table := range schema.Tables {
			if table.Name != "" {
				tables[tableKey{schema.Name, table.Name}] = table
			}
		}
	}
	return tables
}
