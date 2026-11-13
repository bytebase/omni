package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingUniqueKey judges the names of the unique keys top-level statements
// create or rename against cfg: the UNIQUE definitions of a CREATE TABLE,
// those an ALTER TABLE adds, the new name an ALTER TABLE ... RENAME INDEX
// gives, and the index of a CREATE UNIQUE INDEX. A column's own UNIQUE
// attribute is not judged. Statements nested in a routine body are not
// judged.
//
// A renamed index is judged only when schema holds the old name on the
// statement's table, and holds it as unique. Whether it is also the primary
// key is not read, so renaming a primary key held as unique is judged here.
//
// The name is the one written after UNIQUE [INDEX | KEY], else the
// CONSTRAINT symbol, else empty: a unique key with no name is judged as the
// empty name. The schema lookup, the unit of length, the expansion of
// cfg.Format, the codes, and the anchors are NamingIndex's. Ported from
// Bytebase's V1 rule NAMING_INDEX_UK.
func NamingUniqueKey(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	return namingIndexes(stmts, cfg, schema, true)
}
