package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowSetCharset reports each top-level CREATE TABLE or ALTER
// TABLE statement that gives a column a CHARACTER SET other than binary,
// compared without regard to case, once however many columns it gives one,
// anchored at the whole statement. The columns judged are every column a
// CREATE TABLE defines and the columns ALTER TABLE's ADD, CHANGE and MODIFY
// COLUMN write; the statement need not name a table. A COLLATE alone is not
// judged. Ported from Bytebase's V1 rule COLUMN_DISALLOW_SET_CHARSET.
func ColumnDisallowSetCharset(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var cols []*ast.ColumnDef
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			cols, loc = n.Columns, n.Loc
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd != nil {
					cols = append(cols, alteredColumns(cmd)...)
				}
			}
			loc = n.Loc
		default:
			continue
		}
		if slices.ContainsFunc(cols, setsCharset) {
			start, end := span(loc)
			findings = append(findings, v1review.Finding{Code: v1review.ColumnCharsetSet, Start: start, End: end})
		}
	}
	return findings
}

func setsCharset(col *ast.ColumnDef) bool {
	if col == nil || col.TypeName == nil {
		return false
	}
	switch strings.ToLower(col.TypeName.Charset) {
	case "", "binary":
		return false
	default:
		return true
	}
}
