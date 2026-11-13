package v1review

import (
	"slices"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequired judges each top-level statement on its own for tables that
// lack a required column: a CREATE TABLE lacking one among the columns it
// defines, so one created by LIKE, or by OF or PARTITION OF without column
// options, lacks them all; an ALTER ... DROP COLUMN dropping one, whatever
// relation the ALTER names; and a RENAME COLUMN renaming one to another name,
// on any relation. Names are
// compared as the parser stores them, unquoted names folded to lower case,
// and a required name may be empty, which a CREATE TABLE always lacks. A
// statement whose relation has no name is skipped, and CREATE TABLE AS and a
// statement in a routine body are not judged.
//
// Each finding reports ColumnMissing, with Name the missing column and Table
// the relation's name without its schema. A CREATE TABLE reports its missing
// columns in bytewise name order, each anchored at the statement; a DROP
// COLUMN reports its column anchored at its action, once per action; a
// RENAME reports the old name anchored at the statement. Findings that share
// an anchor are one V1 advice; Bytebase's V1 advisor reported every one at
// the line its statement starts on. Ported from Bytebase's V1 rule
// COLUMN_REQUIRED.
func ColumnRequired(stmts []ast.Node, required []string) []v1review.Finding {
	need := make(map[string]bool, len(required))
	for _, column := range required {
		need[column] = true
	}
	var findings []v1review.Finding
	missing := func(loc ast.Loc, table, column string) {
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: v1review.ColumnMissing, Name: column, Table: table, Start: start, End: end})
	}
	for _, stmt := range stmts {
		switch n := statementNode(stmt).(type) {
		case *ast.CreateStmt:
			table := relname(n.Relation)
			if table == "" {
				continue
			}
			lacking := make(map[string]bool, len(need))
			for column := range need {
				lacking[column] = true
			}
			for _, col := range tableColumns(n) {
				delete(lacking, col.Colname)
			}
			columns := make([]string, 0, len(lacking))
			for column := range lacking {
				columns = append(columns, column)
			}
			slices.Sort(columns)
			for _, column := range columns {
				missing(n.Loc, table, column)
			}
		case *ast.AlterTableStmt:
			table := relname(n.Relation)
			if table == "" {
				continue
			}
			for _, cmd := range alterCommands(n) {
				if ast.AlterTableType(cmd.Subtype) == ast.AT_DropColumn && need[cmd.Name] {
					missing(cmd.Loc, table, cmd.Name)
				}
			}
		case *ast.RenameStmt:
			if n.RenameType != ast.OBJECT_COLUMN {
				continue
			}
			table := relname(n.Relation)
			if table == "" {
				continue
			}
			if need[n.Subname] && n.Subname != n.Newname {
				missing(n.Loc, table, n.Subname)
			}
		}
	}
	return findings
}
