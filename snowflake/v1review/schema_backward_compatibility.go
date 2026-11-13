package v1review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// SchemaBackwardCompatibility judges which top-level statements may break
// the existing data or code. It takes every statement of the change in one
// call, since an object the change itself created earlier may be dropped
// freely, and sets Statement on each finding. currentDatabase is the
// database the change runs against, used as written.
//
// Names are normalized as Snowflake resolves them: an unquoted part
// upper-cased, a quoted part kept. A table is remembered as
// database.schema.table, the database defaulting to currentDatabase and the
// schema to PUBLIC; a schema as database.schema; a database by its name.
// Each CREATE TABLE, CREATE SCHEMA and CREATE DATABASE, AS SELECT and OR
// REPLACE included, remembers whether it lacked IF NOT EXISTS, the last
// such statement for a name winning. A drop does not forget the object.
//
// Reported, in statement order, with Name holding the text the V1 advice
// printed, each normalized name in Go-quoted form (strconv.Quote):
//   - IncompatibleDropTable, IncompatibleDropSchema,
//     IncompatibleDropDatabase: a DROP TABLE, DROP SCHEMA or DROP DATABASE
//     of an object the change did not create. When the change created it
//     with IF NOT EXISTS the code is IncompatibleDropMaybeNewTable,
//     IncompatibleDropMaybeNewSchema or IncompatibleDropMaybeNewDatabase;
//     when it created it without, the drop is not reported.
//   - IncompatibleDropColumn: each DROP COLUMN action of an ALTER TABLE,
//     whether or not the change created the table, Name holding its
//     non-empty columns joined by ",".
//
// Every finding anchors the statement node, whose range for DROP SCHEMA,
// DROP DATABASE and ALTER TABLE starts after the leading DROP or ALTER.
// Ported from Bytebase's V1 rule SCHEMA_BACKWARD_COMPATIBILITY.
func SchemaBackwardCompatibility(stmts []ast.Node, currentDatabase string) []v1review.Finding {
	// Each map holds, per normalized name, whether the change created the
	// object without IF NOT EXISTS.
	newTables := map[string]bool{}
	newSchemas := map[string]bool{}
	newDatabases := map[string]bool{}
	var findings []v1review.Finding
	add := func(code v1review.Code, name string, loc ast.Loc, i int) {
		start, end := statementSpan(loc)
		findings = append(findings, v1review.Finding{Code: code, Name: name, Start: start, End: end, Statement: i})
	}
	// dropped reports a drop of name unless the change created it without
	// IF NOT EXISTS.
	dropped := func(created map[string]bool, name string, code, maybeNew v1review.Code, loc ast.Loc, i int) {
		mustNew, ok := created[name]
		switch {
		case ok && mustNew:
		case ok:
			add(maybeNew, strconv.Quote(name), loc, i)
		default:
			add(code, strconv.Quote(name), loc, i)
		}
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			newTables[normalizeTableName(n.Name, currentDatabase)] = !n.IfNotExists
		case *ast.CreateSchemaStmt:
			newSchemas[normalizeSchemaName(n.Name, currentDatabase)] = !n.IfNotExists
		case *ast.CreateDatabaseStmt:
			newDatabases[normalizeDatabaseName(n.Name)] = !n.IfNotExists
		case *ast.DropStmt:
			if n.Kind == ast.DropTable {
				dropped(newTables, normalizeTableName(n.Name, currentDatabase), v1review.IncompatibleDropTable, v1review.IncompatibleDropMaybeNewTable, n.Loc, i)
			}
		case *ast.DropSchemaStmt:
			dropped(newSchemas, normalizeSchemaName(n.Name, currentDatabase), v1review.IncompatibleDropSchema, v1review.IncompatibleDropMaybeNewSchema, n.Loc, i)
		case *ast.DropDatabaseStmt:
			dropped(newDatabases, normalizeDatabaseName(n.Name), v1review.IncompatibleDropDatabase, v1review.IncompatibleDropMaybeNewDatabase, n.Loc, i)
		case *ast.AlterTableStmt:
			for _, action := range n.Actions {
				if action == nil || action.Kind != ast.AlterTableDropColumn {
					continue
				}
				columns := make([]string, 0, len(action.DropColumnNames))
				for _, column := range action.DropColumnNames {
					if !column.IsEmpty() {
						columns = append(columns, strconv.Quote(column.Normalize()))
					}
				}
				add(v1review.IncompatibleDropColumn, strings.Join(columns, ","), n.Loc, i)
			}
		}
	}
	return findings
}

// statementSpan is the range of a statement's location, kept even when its
// end is unknown, since V1 read only where the statement starts.
func statementSpan(loc ast.Loc) (int, int) {
	if loc.Start < 0 {
		return 0, 0
	}
	return loc.Start, max(loc.End, loc.Start)
}

// normalizeTableName returns database.schema.table for name, the database
// defaulting to currentDatabase and the schema to PUBLIC, the table part
// left out when empty.
func normalizeTableName(name *ast.ObjectName, currentDatabase string) string {
	database, schema, table := currentDatabase, "PUBLIC", ""
	if name != nil {
		if d := name.Database.Normalize(); d != "" {
			database = d
		}
		if s := name.Schema.Normalize(); s != "" {
			schema = s
		}
		table = name.Name.Normalize()
	}
	parts := []string{database, schema}
	if table != "" {
		parts = append(parts, table)
	}
	return strings.Join(parts, ".")
}

// normalizeSchemaName returns database.schema for a schema's name, which
// parses with its database in the Schema part; the database defaults to
// currentDatabase.
func normalizeSchemaName(name *ast.ObjectName, currentDatabase string) string {
	database, schema := currentDatabase, ""
	if name != nil {
		if d := name.Schema.Normalize(); d != "" {
			database = d
		}
		schema = name.Name.Normalize()
	}
	return database + "." + schema
}

// normalizeDatabaseName returns a database's name, "" when it has none.
func normalizeDatabaseName(name *ast.ObjectName) string {
	if name == nil {
		return ""
	}
	return name.Name.Normalize()
}
