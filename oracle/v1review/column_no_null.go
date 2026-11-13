package v1review

import (
	"slices"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnNoNull judges the top-level statements together for nullable
// columns. Each column is keyed table.column, the table as schema.name with
// the schema defaulting to currentDatabase. A column a CREATE TABLE defines,
// or an ADD or MODIFY action of an ALTER TABLE defines (its column list,
// then its single column), is nullable unless it says NOT NULL or carries a
// NOT NULL or PRIMARY KEY constraint; a nullable one is recorded, anchored
// at that definition, and a non-nullable one is forgotten. A CREATE TABLE
// then forgets every column its PRIMARY KEY table constraints name. Nothing
// else forgets a column: a later DROP COLUMN, DROP TABLE or a NOT NULL
// table constraint does not, and a MODIFY that only changes the type
// records the column again.
//
// Each finding is one recorded column, in key order compared bytewise,
// with code ColumnNullable, Name the key's text after its last dot (so a
// column name holding a dot prints only its tail), Statement and the range
// the last definition that recorded it. Ported from Bytebase's V1 rule
// COLUMN_NO_NULL.
func ColumnNoNull(stmts []ast.Node, currentDatabase string) []v1review.Finding {
	nullable := make(map[string]columnAnchor)
	record := func(i int, table string, col *ast.ColumnDef) {
		key := table + "." + col.Name
		if col.NotNull || columnHasConstraint(col, ast.CONSTRAINT_NOT_NULL) || columnHasConstraint(col, ast.CONSTRAINT_PRIMARY) {
			delete(nullable, key)
			return
		}
		nullable[key] = columnAnchor{statement: i, loc: col.Loc}
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			table := objectName(n.Name, currentDatabase)
			for _, col := range columnDefs(n.Columns) {
				record(i, table, col)
			}
			for _, c := range tableConstraints(n.Constraints) {
				if c.Type == ast.CONSTRAINT_PRIMARY {
					for _, column := range listStrings(c.Columns) {
						delete(nullable, table+"."+column)
					}
				}
			}
		case *ast.AlterTableStmt:
			table := objectName(n.Name, currentDatabase)
			for _, cmd := range alterTableCmds(n) {
				if cmd.Action != ast.AT_MODIFY_COLUMN && cmd.Action != ast.AT_ADD_COLUMN {
					continue
				}
				for _, col := range actionColumns(cmd) {
					if col != nil {
						record(i, table, col)
					}
				}
			}
		}
	}
	return keyedColumnFindings(nullable, v1review.ColumnNullable)
}

// keyedColumnFindings reports each recorded column in key order, compared
// bytewise, named by the key's text after its last dot.
func keyedColumnFindings(columns map[string]columnAnchor, code v1review.Code) []v1review.Finding {
	keys := make([]string, 0, len(columns))
	for key := range columns {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	findings := make([]v1review.Finding, 0, len(keys))
	for _, key := range keys {
		anchor := columns[key]
		start, end := span(anchor.loc)
		findings = append(findings, v1review.Finding{Code: code, Name: lastIdentifier(key), Start: start, End: end, Statement: anchor.statement})
	}
	return findings
}
