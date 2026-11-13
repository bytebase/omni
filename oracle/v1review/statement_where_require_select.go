package v1review

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireSelect reports, in walk order, every SELECT with a
// FROM clause and no WHERE clause anywhere in the statements: subqueries,
// set-operation branches, INSERT ... SELECT, and PL/SQL bodies included. A
// SELECT whose FROM is DUAL alone, in any letter case, is not judged, nor is
// one without FROM. A finding anchors the SELECT. WhereMissing is the only
// code; Name is not set. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_SELECT.
func StatementWhereRequireSelect(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ast.Inspect(stmt, func(n ast.Node) bool {
			sel, ok := n.(*ast.SelectStmt)
			if !ok || sel.FromClause == nil || len(sel.FromClause.Items) == 0 {
				return true
			}
			if len(sel.FromClause.Items) == 1 {
				if table, ok := sel.FromClause.Items[0].(*ast.TableRef); ok && table.Name != nil && strings.EqualFold(table.Name.Name, "DUAL") {
					return true
				}
			}
			if sel.WhereClause == nil {
				start, end := span(sel.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Start: start, End: end})
			}
			return true
		})
	}
	return findings
}
