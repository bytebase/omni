package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableRequireCharset judges every top-level CREATE TABLE with a table
// name, including CREATE TABLE ... LIKE and ... AS SELECT; statements
// nested in a routine body and ALTER TABLE are not judged. A statement
// with no [DEFAULT] CHARSET or CHARACTER SET option (compared
// case-insensitively) is TableCharsetMissing, with Name the table name
// without its database; a COLLATE option does not satisfy it. Each finding
// anchors the whole statement. Ported from Bytebase's V1 rule
// TABLE_REQUIRE_CHARSET.
func TableRequireCharset(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.CreateTableStmt)
		if !ok || n.Table == nil || n.Table.Name == "" {
			continue
		}
		hasCharset := false
		for _, opt := range n.Options {
			if opt != nil && (strings.EqualFold(opt.Name, "CHARSET") || strings.EqualFold(opt.Name, "CHARACTER SET") || strings.EqualFold(opt.Name, "DEFAULT CHARSET") || strings.EqualFold(opt.Name, "DEFAULT CHARACTER SET")) {
				hasCharset = true
				break
			}
		}
		if !hasCharset {
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.TableCharsetMissing, Name: n.Table.Name, Start: start, End: end})
		}
	}
	return findings
}
