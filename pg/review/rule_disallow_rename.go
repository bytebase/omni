package review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowRename reports the RENAME TO of a relation a query reads
// by name (ALTER TABLE, ALTER FOREIGN TABLE, ALTER VIEW, and ALTER
// MATERIALIZED VIEW) and the rename of a column of any relation kind, at
// top level or inside a SQL-standard routine body. Renaming an index, a
// sequence, a constraint, a schema, or another object is not in scope.
func checkDisallowRename(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		v, ok := n.(*ast.RenameStmt)
		if !ok {
			return true
		}
		switch v.RenameType {
		case ast.OBJECT_TABLE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW:
			r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames "+objectKind(v.RenameType)+" "+relation(v.Relation)+" to "+ident(v.Newname))
		case ast.OBJECT_COLUMN:
			r.report(review.DisallowRename, s.index, rangeOf(v.Loc), "renames column "+ident(v.Subname)+" of "+relation(v.Relation)+" to "+ident(v.Newname))
		}
		return true
	})
}
