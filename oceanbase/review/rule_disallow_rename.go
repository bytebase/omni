package review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowRename reports the rename of a table or view, by RENAME
// TABLE or ALTER TABLE ... RENAME, and the rename of a column, by RENAME
// COLUMN or a CHANGE COLUMN that gives the column another name, at top
// level or inside a routine body. Renaming an index is not in scope.
func checkDisallowRename(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.RenameTableStmt:
			for _, pair := range v.Pairs {
				if pair == nil {
					continue
				}
				r.report(review.DisallowRename, s.index, rangeOf(pair.Loc), "renames table "+table(pair.Old)+" to "+table(pair.New))
			}
		case *ast.AlterTableStmt:
			for _, cmd := range v.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATRenameTable:
					to := ident(cmd.NewName)
					if cmd.NewTable != nil {
						to = table(cmd.NewTable)
					}
					r.report(review.DisallowRename, s.index, rangeOf(cmd.Loc), "renames table "+table(v.Table)+" to "+to)
				case ast.ATRenameColumn:
					r.report(review.DisallowRename, s.index, rangeOf(cmd.Loc), "renames column "+ident(cmd.Name)+" of "+table(v.Table)+" to "+ident(cmd.NewName))
				case ast.ATChangeColumn:
					if cmd.Column != nil && !sameName(cmd.Name, cmd.Column.Name) {
						r.report(review.DisallowRename, s.index, rangeOf(cmd.Loc), "renames column "+ident(cmd.Name)+" of "+table(v.Table)+" to "+ident(cmd.Column.Name))
					}
				}
			}
		}
		return true
	})
}
