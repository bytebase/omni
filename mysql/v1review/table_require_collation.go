package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableRequireCollation judges every top-level CREATE TABLE with a table
// name, including CREATE TABLE ... LIKE and ... AS SELECT; statements
// nested in a routine body and ALTER TABLE are not judged. A statement
// with no [DEFAULT] COLLATE option (compared case-insensitively) is
// TableCollationMissing, with Name the table name without its database; a
// character set option does not satisfy it, though it implies a
// collation. Each finding anchors the whole statement. Ported from
// Bytebase's V1 rule TABLE_REQUIRE_COLLATION.
func TableRequireCollation(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.CreateTableStmt)
		if !ok || n.Table == nil || n.Table.Name == "" {
			continue
		}
		hasCollation := false
		for _, opt := range n.Options {
			if opt != nil && (strings.EqualFold(opt.Name, "COLLATE") || strings.EqualFold(opt.Name, "DEFAULT COLLATE")) {
				hasCollation = true
				break
			}
		}
		if !hasCollation {
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.TableCollationMissing, Name: n.Table.Name, Start: start, End: end})
		}
	}
	return findings
}
