package review

import (
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/snowflake/ast"
)

// checkDisallowDropObject reports every DROP: the DROP statement of every
// object kind the parser reads (table, view, materialized view, dynamic
// table, external table, stream, task, sequence, stage, file format,
// function, procedure, pipe, tag, role, warehouse, network rule, alert,
// connection, failover group, integration, managed account, the policies,
// replication group, resource monitor, share, and user), DROP DATABASE,
// DROP SCHEMA, and ALTER TABLE ... DROP COLUMN, which anchors each column.
// IF EXISTS does not exempt a statement, and neither does a Snowflake
// Scripting block, a task's body, or an alert's action around it. UNDROP
// restores an object and is not a drop. Dropping a constraint belongs to
// DisallowDropConstraint; dropping a clustering key, a search optimization,
// or a row access policy from a table removes no object.
func checkDisallowDropObject(sql string, s *statement, r *reporter) {
	report := func(rng review.Range, kind, name string) {
		r.report(review.DisallowDropObject, s.index, rng, "drops "+kind+" "+name)
	}
	ast.Inspect(s.node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.DropStmt:
			name := object(v.Name)
			// An overloaded function or procedure is named by its signature.
			if v.HasArgs {
				args := make([]string, len(v.ArgTypes))
				for i, t := range v.ArgTypes {
					if t != nil && t.Loc.Start >= 0 && t.Loc.End <= len(sql) && t.Loc.Start < t.Loc.End {
						args[i] = oneLine(sql[t.Loc.Start:t.Loc.End])
					} else if t != nil {
						args[i] = t.Name
					}
				}
				name += "(" + strings.Join(args, ", ") + ")"
			}
			report(anchor(s, v, v.Loc), strings.ToLower(v.Kind.String()), name)
		case *ast.DropDatabaseStmt:
			report(anchor(s, v, v.Loc), "database", object(v.Name))
		case *ast.DropSchemaStmt:
			report(anchor(s, v, v.Loc), "schema", object(v.Name))
		case *ast.AlterTableStmt:
			for _, a := range v.Actions {
				if a == nil || a.Kind != ast.AlterTableDropColumn {
					continue
				}
				for _, col := range a.DropColumnNames {
					report(rangeOf(col.Loc), "column", ident(col)+" of "+object(v.Name))
				}
			}
		}
		return true
	})
}
