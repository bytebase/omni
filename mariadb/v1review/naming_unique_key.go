package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingUniqueKey judges the names of the unique keys top-level statements
// create or rename against cfg:
//
//   - the table-level UNIQUE definitions of a CREATE TABLE;
//   - those an ALTER TABLE adds with ADD [CONSTRAINT] UNIQUE;
//   - a CREATE UNIQUE INDEX;
//   - the new name of an ALTER TABLE ... RENAME {INDEX | KEY}, when schema
//     holds the old name on that table and marks the index unique. A rename
//     of an index schema does not hold is not judged.
//
// A column's own UNIQUE attribute, a unique key inside ALTER TABLE's
// parenthesized ADD form, and statements nested in a routine body are not
// judged. The name is the one written after UNIQUE [INDEX | KEY], else the
// CONSTRAINT symbol, else empty, and an empty name is judged like any
// other; its length is in bytes.
//
// The pattern, the schema lookup, the anchors and the order of findings are
// NamingIndex's. Ported from Bytebase's V1 rule NAMING_INDEX_UK.
func NamingUniqueKey(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	naming := &indexNaming{
		cfg:        cfg,
		schema:     schema,
		constraint: ast.ConstrUnique,
		created:    func(n *ast.CreateIndexStmt) bool { return n.Unique },
		unique:     true,
	}
	naming.judge(stmts)
	return naming.findings
}
