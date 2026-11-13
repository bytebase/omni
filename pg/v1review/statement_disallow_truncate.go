package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowTruncate reports each table a top-level TRUNCATE lists,
// in the order written; a TRUNCATE nested in a function body is not judged,
// and a list item that is not a plain table reference is skipped. Table is
// the table as written, schema-qualified when the statement qualified it.
// Detail lists the options that change what the TRUNCATE does: "CASCADE"
// and "RESTART IDENTITY", comma-separated in that order, empty when neither
// is written. Each finding anchors its own table reference. Truncate is the
// only code. Ported from Bytebase's V1 rule STATEMENT_DISALLOW_TRUNCATE.
func StatementDisallowTruncate(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		truncate, ok := unwrapRaw(stmt).(*ast.TruncateStmt)
		if !ok || truncate.Relations == nil {
			continue
		}
		var options []string
		if truncate.Behavior == ast.DROP_CASCADE {
			options = append(options, "CASCADE")
		}
		if truncate.RestartSeqs {
			options = append(options, "RESTART IDENTITY")
		}
		for _, item := range truncate.Relations.Items {
			table, ok := item.(*ast.RangeVar)
			if !ok {
				continue
			}
			name := table.Relname
			if table.Schemaname != "" {
				name = table.Schemaname + "." + table.Relname
			}
			start, end := anchor(table.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.Truncate, Table: name, Detail: strings.Join(options, ","), Start: start, End: end})
		}
	}
	return findings
}
