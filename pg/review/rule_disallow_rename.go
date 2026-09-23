package review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowRename reports ALTER TABLE ... RENAME TO and the rename of
// a column of any relation kind, at top level or inside a SQL-standard
// routine body. Renaming an index, a constraint, a schema, or another
// object is not in scope.
func checkDisallowRename(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		v, ok := n.(*ast.RenameStmt)
		if !ok {
			return true
		}
		switch v.RenameType {
		case ast.OBJECT_TABLE:
			r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames table "+relation(v.Relation)+" to "+ident(v.Newname))
		case ast.OBJECT_COLUMN:
			r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames column "+ident(v.Subname)+" of "+relation(v.Relation)+" to "+ident(v.Newname))
		}
		return true
	})
}
