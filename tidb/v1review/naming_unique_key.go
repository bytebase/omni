package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingUniqueKey judges the names of the unique keys top-level statements
// define or rename against cfg: the table-level UNIQUE constraints of a
// CREATE TABLE, those an ALTER TABLE adds, the index a CREATE UNIQUE INDEX
// creates, and the new name of an ALTER TABLE ... RENAME INDEX. A UNIQUE
// written on a column is not judged, and neither are statements nested in
// a routine body or carried by EXPLAIN or PLAN REPLAYER.
//
// Name is the CONSTRAINT symbol, or the name written after UNIQUE when
// there is no symbol, or empty when there is neither: an unnamed key is
// judged as the empty name. Table is the statement's table without its
// qualifier. The pattern is cfg expanded with that table and the key's
// column names joined with "_"; a key part that is an expression gives no
// name. A key whose pattern does not compile gets one PatternInvalid
// finding and is not judged further. Length is in bytes.
//
// A renamed index is looked up by its old name in the unnamed schema of
// schema, across every table and materialized view and not only the
// statement's table, whatever database the statement names; index names
// compare without case. The rename is judged only when that index exists
// and is unique, which a primary key's index recorded as unique is too, as
// its new name against the found index's expressions and the statement's
// table. Where two tables hold an index of that name, V1 took whichever
// its map yielded and this takes the first in the snapshot's order.
//
// A CREATE TABLE finding anchors the constraint. Every other finding has
// no range: V1 anchors it at the first character of the statement's text
// that is not white space, which a leading comment moves off the
// statement. Ported from Bytebase's V1 rule NAMING_INDEX_UK.
func NamingUniqueKey(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	return judgeKeys(stmts, cfg, schema, true)
}
