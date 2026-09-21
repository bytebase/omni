package review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowRename reports ALTER TABLE ... RENAME TO and the rename of
// a column of any relation kind. Renaming an index, a constraint, a
// schema, or another object is not in scope.
func checkDisallowRename(s *statement, r *reporter) {
	v, ok := s.node.(*ast.RenameStmt)
	if !ok {
		return
	}
	switch v.RenameType {
	case ast.OBJECT_TABLE:
		r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames table "+relation(v.Relation)+" to "+ident(v.Newname))
	case ast.OBJECT_COLUMN:
		r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames column "+ident(v.Subname)+" of "+relation(v.Relation)+" to "+ident(v.Newname))
	}
}
