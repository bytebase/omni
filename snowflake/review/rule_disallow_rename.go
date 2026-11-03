package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/snowflake/ast"
)

// checkDisallowRename reports the rename of a relation a query reads by
// name and of a table's column, at top level or inside a Snowflake
// Scripting block, a task's body, or an alert's action:
//
//   - ALTER TABLE ... RENAME TO, and ALTER TABLE ... SWAP WITH, which
//     "renames two tables in a single transaction";
//   - ALTER DYNAMIC TABLE ... RENAME TO and SWAP WITH;
//   - ALTER VIEW, ALTER MATERIALIZED VIEW, and ALTER SEMANTIC VIEW ...
//     RENAME TO;
//   - ALTER TABLE ... RENAME COLUMN.
//
// Renaming another object (a database, a schema, a sequence, a stage, a
// function, a constraint, and the rest) is not in scope: the rule is about
// the tables, views, and columns queries name.
//
// sql: https://docs.snowflake.com/en/sql-reference/sql/alter-table
func checkDisallowRename(s *statement, r *reporter) {
	report := func(rng review.Range, msg string) {
		r.report(review.DisallowRename, s.index, rng, msg)
	}
	renames := func(kind string, from, to *ast.ObjectName) string {
		return "Renames " + kind + " " + review.Code(object(from)) + " to " + review.Code(object(to)) + "."
	}
	swaps := func(kind string, a, b *ast.ObjectName) string {
		return "Swaps the names of " + kind + " " + review.Names(object(a), object(b)) + "."
	}
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.ExplainStmt:
			// EXPLAIN compiles its statement without running it.
			return false
		case *ast.AlterTableStmt:
			for _, a := range v.Actions {
				if a == nil {
					continue
				}
				switch a.Kind {
				case ast.AlterTableRename:
					report(rangeOf(a.Loc), renames("table", v.Name, a.NewName))
				case ast.AlterTableSwapWith:
					report(rangeOf(a.Loc), swaps("tables", v.Name, a.NewName))
				case ast.AlterTableRenameColumn:
					report(rangeOf(a.Loc), "Renames column "+review.Code(ident(a.OldName))+" of "+review.Code(object(v.Name))+" to "+review.Code(ident(a.NewColName))+".")
				}
			}
		case *ast.AlterDynamicTableStmt:
			switch v.Action {
			case ast.AlterDynamicTableRename:
				report(anchor(s, v, v.Loc), renames("dynamic table", v.Name, v.NewName))
			case ast.AlterDynamicTableSwap:
				report(anchor(s, v, v.Loc), swaps("dynamic tables", v.Name, v.NewName))
			}
		case *ast.AlterViewStmt:
			if v.Action == ast.AlterViewRename {
				report(anchor(s, v, v.Loc), renames("view", v.Name, v.NewName))
			}
		case *ast.AlterMaterializedViewStmt:
			if v.Action == ast.AlterMVRename {
				report(anchor(s, v, v.Loc), renames("materialized view", v.Name, v.NewName))
			}
		case *ast.AlterSemanticViewStmt:
			if v.Action == ast.AlterSemanticViewRename {
				report(anchor(s, v, v.Loc), renames("semantic view", v.Name, v.NewName))
			}
		}
		return true
	})
}
