package v1review

import (
	"slices"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowDML judges the tables every top-level INSERT, REPLACE,
// UPDATE and DELETE writes; statements nested in a routine body are not
// judged. A name that equals an entry of disallow exactly
// (case-sensitively) is TableDisallowed, with Name the name. INSERT and
// REPLACE judge their target table, not the tables of their SELECT. UPDATE
// judges every table its table references hold, in source order: joined
// tables, and tables read by a derived table or a join condition's
// subquery, though not the WHERE clause's. DELETE judges the tables it
// deletes from, not the other tables of a multi-table FROM or USING.
// Names are without their database; an alias is not judged in a table's
// place. A SELECT ... INTO OUTFILE judges the file name as if it were a
// table, as V1 did. Every finding anchors the whole statement. Ported from
// Bytebase's V1 rule TABLE_DISALLOW_DML.
func TableDisallowDML(stmts []ast.Node, disallow []string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var names []string
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.DeleteStmt:
			loc = n.Loc
			names = tableDisallowDMLNames(n.Tables)
		case *ast.InsertStmt:
			loc = n.Loc
			if n.Table != nil {
				names = append(names, n.Table.Name)
			}
		case *ast.SelectStmt:
			loc = n.Loc
			if n.Into != nil && n.Into.Outfile != "" {
				names = append(names, n.Into.Outfile)
			}
		case *ast.UpdateStmt:
			loc = n.Loc
			names = tableDisallowDMLNames(n.Tables)
		default:
			continue
		}
		start, end := span(loc)
		for _, name := range names {
			if slices.Contains(disallow, name) {
				findings = append(findings, v1review.Finding{Code: v1review.TableDisallowed, Name: name, Start: start, End: end})
			}
		}
	}
	return findings
}

// tableDisallowDMLNames returns the name of every table reference under
// exprs, in walk order, not descending into a table reference.
func tableDisallowDMLNames(exprs []ast.TableExpr) []string {
	var names []string
	for _, expr := range exprs {
		ast.Inspect(expr, func(n ast.Node) bool {
			t, ok := n.(*ast.TableRef)
			if !ok {
				return true
			}
			names = append(names, t.Name)
			return false
		})
	}
	return names
}
