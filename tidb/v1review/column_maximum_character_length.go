package v1review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnMaximumCharacterLength reports, at most once per top-level CREATE
// TABLE or ALTER TABLE statement, a CHAR or BINARY column longer than
// maximum; a CHAR or BINARY written without a length has length 1. The
// caller passes a positive maximum. Code is ColumnCharTooLong for both
// types, Name the column and Value its length.
//
// In a CREATE TABLE the first such column is reported, anchored at its
// definition. In an ALTER TABLE the actions are read in order and the first
// that writes such a column decides: for ADD COLUMN the last such column it
// adds, for CHANGE or MODIFY COLUMN the column it writes, anchored at the
// whole statement. A statement naming no table is skipped. Ported from
// Bytebase's V1 rule COLUMN_MAXIMUM_CHARACTER_LENGTH.
func ColumnMaximumCharacterLength(stmts []ast.Node, maximum int) []v1review.Finding {
	var findings []v1review.Finding
	tooLong := func(col *ast.ColumnDef, loc ast.Loc) v1review.Finding {
		start, end := span(loc)
		return v1review.Finding{Code: v1review.ColumnCharTooLong, Name: col.Name, Value: strconv.Itoa(charLength(col.TypeName)), Start: start, End: end}
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, col := range n.Columns {
				if col != nil && charLength(col.TypeName) > maximum {
					findings = append(findings, tooLong(col, col.Loc))
					break
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				var last *ast.ColumnDef
				for _, col := range alteredColumns(cmd) {
					if col != nil && charLength(col.TypeName) > maximum {
						last = col
					}
				}
				if last != nil {
					findings = append(findings, tooLong(last, n.Loc))
					break
				}
			}
		}
	}
	return findings
}

// charLength returns the length of a CHAR or BINARY type, 1 when written
// without one, and 0 for any other type.
func charLength(dt *ast.DataType) int {
	if dt == nil {
		return 0
	}
	switch strings.ToUpper(dt.Name) {
	case "CHAR", "BINARY":
		if dt.Length == 0 {
			return 1
		}
		return dt.Length
	default:
		return 0
	}
}
