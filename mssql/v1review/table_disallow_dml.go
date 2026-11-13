package v1review

import (
	"slices"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowDML judges the table a top-level INSERT, UPDATE, DELETE or
// MERGE targets and the table a top-level SELECT ... INTO creates;
// statements nested in a routine body or an IF or BEGIN block are not
// judged, nor is a target written as a table variable, and nor are the
// tables a statement reads. A target written as an alias is judged as that
// alias. The table is
// judged as TableDisallowDDL judges it, and a statement is TableDisallowed
// when that name equals an entry of disallow exactly (case-sensitively),
// with Name the joined name. Every finding anchors the whole statement.
// Ported from Bytebase's V1 rule TABLE_DISALLOW_DML.
func TableDisallowDML(stmts []ast.Node, disallow []string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var name string
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.InsertStmt:
			if ref, ok := n.Relation.(*ast.TableRef); ok {
				name, loc = joinTableRef(ref, ""), n.Loc
			}
		case *ast.UpdateStmt:
			if ref, ok := n.Relation.(*ast.TableRef); ok {
				name, loc = joinTableRef(ref, ""), n.Loc
			}
		case *ast.DeleteStmt:
			if ref, ok := n.Relation.(*ast.TableRef); ok {
				name, loc = joinTableRef(ref, ""), n.Loc
			}
		case *ast.MergeStmt:
			if ref, ok := n.Target.(*ast.TableRef); ok {
				name, loc = joinTableRef(ref, ""), n.Loc
			}
		case *ast.SelectStmt:
			if n.IntoTable != nil {
				name, loc = joinTableRef(n.IntoTable, ""), n.Loc
			}
		default:
		}
		if name == "" || !slices.Contains(disallow, name) {
			continue
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: v1review.TableDisallowed, Name: name, Start: start, End: end})
	}
	return findings
}
