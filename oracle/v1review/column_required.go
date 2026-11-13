package v1review

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequired judges each top-level statement on its own, carrying
// nothing across statements, for required columns it leaves a table
// without. A CREATE TABLE lacks every required column its column list does
// not define, so one created AS SELECT lacks them all. An ALTER TABLE lacks
// each required column one of its DROP COLUMN actions drops, by the action's
// single column name only, or one of its RENAME COLUMN actions renames to a
// different name; what an action adds is not judged. Names are compared as
// the parser stores them, an unquoted name upper-cased, against the entries
// as given.
//
// Each finding is one missing column: Name is the column, Table the table's
// name without its schema, and the anchor is the statement. A statement's
// findings come in the order of their names quoted as Go string literals,
// since V1 sorted the quoted text it printed, compared bytewise. Ported
// from Bytebase's V1 rule COLUMN_REQUIRED.
func ColumnRequired(stmts []ast.Node, required []string) []v1review.Finding {
	need := make(map[string]bool, len(required))
	for _, column := range required {
		need[column] = true
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		missing := make(map[string]bool)
		var table string
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			for column := range need {
				missing[column] = true
			}
			for _, col := range columnDefs(n.Columns) {
				delete(missing, col.Name)
			}
			table, loc = lastObjectName(n.Name), n.Loc
		case *ast.AlterTableStmt:
			for _, cmd := range alterTableCmds(n) {
				switch cmd.Action {
				case ast.AT_DROP_COLUMN:
					if need[cmd.ColumnName] {
						missing[cmd.ColumnName] = true
					}
				case ast.AT_RENAME_COLUMN:
					if cmd.ColumnName != cmd.NewName && need[cmd.ColumnName] {
						missing[cmd.ColumnName] = true
					}
				default:
				}
			}
			table, loc = lastObjectName(n.Name), n.Loc
		default:
			continue
		}
		columns := make([]string, 0, len(missing))
		for column := range missing {
			columns = append(columns, column)
		}
		slices.SortFunc(columns, func(a, b string) int {
			return strings.Compare(strconv.Quote(a), strconv.Quote(b))
		})
		start, end := span(loc)
		for _, column := range columns {
			findings = append(findings, v1review.Finding{Code: v1review.ColumnMissing, Name: column, Table: table, Start: start, End: end})
		}
	}
	return findings
}

func lastObjectName(name *ast.ObjectName) string {
	if name == nil {
		return ""
	}
	return name.Name
}
