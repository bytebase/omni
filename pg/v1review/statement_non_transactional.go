package v1review

import (
	"regexp"
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// The patterns Bytebase's V1 advisor matched a statement's text against,
// copied from its PostgreSQL driver.
var (
	dropDatabasePattern = regexp.MustCompile(`(?i)DROP DATABASE`)
	createIndexPattern  = regexp.MustCompile(`(?i)CREATE(\s+(UNIQUE\s+)?)INDEX(\s+)CONCURRENTLY`)
	dropIndexPattern    = regexp.MustCompile(`(?i)DROP(\s+)INDEX(\s+)CONCURRENTLY`)
	vacuumPattern       = regexp.MustCompile(`(?i)^\s*VACUUM`)
)

// StatementNonTransactional reports each top-level DROP DATABASE, CREATE
// INDEX, DROP or VACUUM (which includes ANALYZE) statement whose text
// matches a pattern of a statement that cannot run in a transaction, one
// NonTransactional finding per statement. text is the text the statements
// were parsed from; Bytebase's V1 advisor matched each statement's own text,
// leading comments included, so it passes one statement at a time. The
// text, with surrounding white space and trailing semicolons removed, is
// matched anywhere for DROP DATABASE (with a single space), CREATE [UNIQUE]
// INDEX CONCURRENTLY and DROP INDEX CONCURRENTLY, and at its start for
// VACUUM, all in any letter case. So a match inside a comment, a string or
// a quoted name counts, a DROP  DATABASE with another gap does not, and
// VACUUM after a leading comment does not. Findings anchor nothing: the V1
// advisor placed them on the statement's first line. Ported from
// Bytebase's V1 rule STATEMENT_NON_TRANSACTIONAL.
func StatementNonTransactional(stmts []ast.Node, text string) []v1review.Finding {
	trimmed := strings.TrimRight(strings.TrimSpace(text), ";")
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch unwrapRaw(stmt).(type) {
		case *ast.DropdbStmt, *ast.IndexStmt, *ast.DropStmt, *ast.VacuumStmt:
		default:
			continue
		}
		if dropDatabasePattern.MatchString(trimmed) || createIndexPattern.MatchString(trimmed) ||
			dropIndexPattern.MatchString(trimmed) || vacuumPattern.MatchString(trimmed) {
			findings = append(findings, v1review.Finding{Code: v1review.NonTransactional})
		}
	}
	return findings
}
