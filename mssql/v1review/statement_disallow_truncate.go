package v1review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowTruncate reports each top-level TRUNCATE TABLE; one
// nested in a procedure body or a control-flow block is not judged, nor is
// one whose table name starts with #, a local or global temporary table.
// Table is the table as written, schema.table when the statement qualified
// it with a schema; a database or server qualifier is dropped. A finding
// anchors the statement. Truncate is the only code. Ported from Bytebase's
// V1 rule STATEMENT_DISALLOW_TRUNCATE.
func StatementDisallowTruncate(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.TruncateStmt)
		if !ok || n.Table == nil || strings.HasPrefix(n.Table.Object, "#") {
			continue
		}
		table := n.Table.Object
		if n.Table.Schema != "" {
			table = n.Table.Schema + "." + n.Table.Object
		}
		start, end := stmtAnchor(n.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.Truncate, Table: table, Start: start, End: end})
	}
	return findings
}
