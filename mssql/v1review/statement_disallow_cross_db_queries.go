package v1review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowCrossDBQueries reports, in statement then walk order,
// each table reference anywhere in the statements whose database qualifier
// differs, case-insensitively, from the current database. The current
// database starts as currentDatabase, as given, and each top-level USE
// with a database name replaces it with that name lower-cased; a USE is not
// otherwise judged, and one nested in a block does not change it. Name is
// the qualifier as written and Detail the current database at that point.
// A finding anchors the table reference, and Statement is the index of the
// statement holding it. CrossDatabaseQuery is the only code. Ported from
// Bytebase's V1 rule STATEMENT_DISALLOW_CROSS_DB_QUERIES.
func StatementDisallowCrossDBQueries(stmts []ast.Node, currentDatabase string) []v1review.Finding {
	var findings []v1review.Finding
	current := currentDatabase
	for i, stmt := range stmts {
		if use, ok := stmt.(*ast.UseStmt); ok {
			if use.Database != "" {
				current = strings.ToLower(use.Database)
			}
			continue
		}
		ast.Inspect(stmt, func(n ast.Node) bool {
			ref, ok := n.(*ast.TableRef)
			if !ok || ref == nil {
				return true
			}
			if ref.Database != "" && !strings.EqualFold(ref.Database, current) {
				start, end := stmtAnchor(ref.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.CrossDatabaseQuery, Name: ref.Database, Detail: current, Start: start, End: end, Statement: i})
			}
			return true
		})
	}
	return findings
}
