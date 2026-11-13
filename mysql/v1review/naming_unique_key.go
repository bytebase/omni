package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingUniqueKey judges the names top-level statements give unique keys
// against cfg: the table-level UNIQUE definitions of a CREATE TABLE, those
// an ALTER TABLE adds with ADD [CONSTRAINT] UNIQUE, the index a CREATE
// UNIQUE INDEX creates, and the new name of an ALTER TABLE ... RENAME
// {INDEX | KEY} whose old index schema holds as unique, a primary key
// included. A column's own UNIQUE attribute is not judged. The name of a
// definition is the one written after UNIQUE [INDEX | KEY], else the
// CONSTRAINT symbol, else empty.
//
// Everything else is as NamingIndex judges an ordinary index: the name
// text, the rename lookup, the pattern, the anchors, and the order. Ported
// from Bytebase's V1 rule NAMING_INDEX_UK.
func NamingUniqueKey(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	return namingKeys(stmts, cfg, schema, true)
}
