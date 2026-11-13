package v1review

import (
	"slices"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequired judges each top-level statement on its own for required
// columns it leaves a table without. Entries of required are compared as
// given with column names keyed as Snowflake resolves them, an unquoted name
// upper-cased and a quoted one as written, so an entry with a lower-case
// letter matches only a quoted column. A plain CREATE TABLE (not AS SELECT,
// LIKE or CLONE) lacks every entry, duplicates once, that no column it
// defines matches. An ALTER TABLE lacks each entry a DROP COLUMN action
// drops, once however often it is dropped. What an action adds or renames is
// not judged, nor are statements nested in a Snowflake Scripting block or
// under EXPLAIN.
//
// Each finding has code ColumnMissing, Name the entry, Table the table in its
// source form (quoted parts re-quoted, joined by dots), and a statement's
// findings come in name order, compared bytewise. A CREATE TABLE's findings
// anchor the last byte of the column or table constraint that ends last,
// or the statement when it has neither; an ALTER TABLE's anchor its first
// DROP COLUMN action, whether or not that action drops a required column.
// Only the range's start is meaningful. Ported from Bytebase's V1 rule
// COLUMN_REQUIRED.
func ColumnRequired(stmts []ast.Node, required []string) []v1review.Finding {
	need := make(map[string]bool, len(required))
	for _, column := range required {
		need[column] = true
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		missing := make(map[string]bool)
		var loc ast.Loc
		var name *ast.ObjectName
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.AsSelect != nil || n.Like != nil || n.Clone != nil {
				continue
			}
			for column := range need {
				missing[column] = true
			}
			for _, column := range n.Columns {
				delete(missing, column.Name.Normalize())
			}
			loc, name = n.Loc, n.Name
			lastItemEnd := -1
			for _, column := range n.Columns {
				lastItemEnd = max(lastItemEnd, column.Loc.End)
			}
			for _, constraint := range n.Constraints {
				lastItemEnd = max(lastItemEnd, constraint.Loc.End)
			}
			if lastItemEnd > 0 {
				loc = ast.Loc{Start: lastItemEnd - 1, End: lastItemEnd}
			}
		case *ast.AlterTableStmt:
			first := true
			for _, action := range n.Actions {
				if action.Kind != ast.AlterTableDropColumn {
					continue
				}
				if first {
					loc, first = action.Loc, false
				}
				for _, column := range action.DropColumnNames {
					if key := column.Normalize(); need[key] {
						missing[key] = true
					}
				}
			}
			name = n.Name
		default:
			continue
		}
		if len(missing) == 0 {
			continue
		}
		table := name.String()
		columns := make([]string, 0, len(missing))
		for column := range missing {
			columns = append(columns, column)
		}
		slices.Sort(columns)
		start, end := columnAnchor(loc)
		for _, column := range columns {
			findings = append(findings, v1review.Finding{Code: v1review.ColumnMissing, Name: column, Table: table, Start: start, End: end})
		}
	}
	return findings
}
