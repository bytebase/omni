package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/ast"
)

// checkDisallowRename reports the rename of a table or view, by RENAME
// TABLE or ALTER TABLE ... RENAME, and the rename of a column, by RENAME
// COLUMN or a CHANGE COLUMN that gives the column another name, at top
// level or inside a routine body. Renaming an index is not in scope.
func checkDisallowRename(s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExplainStmt:
			// A plain EXPLAIN reports a plan and changes nothing, so what it
			// carries is no change of this kind; ANALYZE runs it, and on TiDB
			// that reaches ALTER TABLE.
			return v.Analyze
		case *ast.RenameTableStmt:
			for _, pair := range v.Pairs {
				if pair == nil {
					continue
				}
				r.report(review.RenamesObject, s.index, rangeOf(pair.Loc), renamesTable(table(pair.Old), table(pair.New)))
			}
		case *ast.AlterTableStmt:
			for _, cmd := range v.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATRenameTable:
					// The new name may carry a database, so write it as a
					// table reference rather than as one identifier.
					to := ident(cmd.NewName)
					if cmd.NewTable != nil {
						to = table(cmd.NewTable)
					}
					r.report(review.RenamesObject, s.index, rangeOf(cmd.Loc), renamesTable(table(v.Table), to))
				case ast.ATRenameColumn:
					r.report(review.RenamesObject, s.index, rangeOf(cmd.Loc), renamesColumn(ident(cmd.Name), table(v.Table), ident(cmd.NewName)))
				case ast.ATChangeColumn:
					if cmd.Column != nil && !sameName(cmd.Name, cmd.Column.Name) {
						r.report(review.RenamesObject, s.index, rangeOf(cmd.Loc), renamesColumn(ident(cmd.Name), table(v.Table), ident(cmd.Column.Name)))
					}
				}
			}
		}
		return true
	})
}

// renamesTable writes the finding for the rename of a table.
func renamesTable(from, to string) string {
	return "Renames table " + review.Code(from) + " to " + review.Code(to) + "."
}

// renamesColumn writes the finding for the rename of a column of a table.
func renamesColumn(from, tbl, to string) string {
	return "Renames column " + review.Code(from) + " of " + review.Code(tbl) + " to " + review.Code(to) + "."
}
