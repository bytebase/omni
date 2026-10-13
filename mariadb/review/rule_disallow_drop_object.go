package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
)

// checkDisallowDropObject reports every DROP: the DROP statement of any
// object kind, and ALTER TABLE ... DROP COLUMN, DROP INDEX, DROP
// PARTITION, DROP SYSTEM VERSIONING, which drops the history rows, and
// DROP PERIOD, which anchor the subcommand. IF EXISTS does not exempt a
// statement, and neither does a routine body around it: the statement
// tree is walked, so a DROP inside a procedure is reported where it is.
// Dropping a primary key, a foreign key, a check constraint, or DROP
// CONSTRAINT belongs to DisallowDropConstraint; a unique index dropped as
// an index is an index here, and a unique key there.
func checkDisallowDropObject(sql string, s *statement, r *reporter) {
	ast.Inspect(s.node, func(n ast.Node) bool {
		// report anchors a DROP statement; reportAt a subcommand.
		reportAt := func(rng review.Range, kind string, names ...string) {
			msg := "drops " + kind
			if len(names) > 0 {
				msg += " " + strings.Join(names, ", ")
			}
			r.report(review.DisallowDropObject, s.index, rng, msg)
		}
		report := func(loc ast.Loc, kind string, names ...string) {
			reportAt(anchor(sql, s, n, loc), kind, names...)
		}
		switch v := n.(type) {
		case *ast.DropTableStmt:
			kind := "table"
			if v.Temporary {
				kind = "temporary table"
			}
			report(v.Loc, kind, tables(v.Tables))
		case *ast.DropViewStmt:
			report(v.Loc, "view", tables(v.Views))
		case *ast.DropSequenceStmt:
			report(v.Loc, "sequence", tables(v.Sequences))
		case *ast.DropDatabaseStmt:
			report(v.Loc, "database", ident(v.Name))
		case *ast.DropIndexStmt:
			if sameName(v.Name, "PRIMARY") {
				return true
			}
			report(v.Loc, "index", ident(v.Name)+" of "+table(v.Table))
		case *ast.DropRoutineStmt:
			kind := "function"
			if v.IsProcedure {
				kind = "procedure"
			}
			report(v.Loc, kind, table(v.Name))
		case *ast.DropTriggerStmt:
			report(v.Loc, "trigger", table(v.Name))
		case *ast.DropEventStmt:
			report(v.Loc, "event", ident(v.Name))
		case *ast.DropUserStmt:
			report(v.Loc, "user", v.Users...)
		case *ast.DropRoleStmt:
			report(v.Loc, "role", v.Roles...)
		case *ast.DropTablespaceStmt:
			kind := "tablespace"
			if v.Undo {
				kind = "undo tablespace"
			}
			report(v.Loc, kind, ident(v.Name))
		case *ast.DropServerStmt:
			report(v.Loc, "server", ident(v.Name))
		case *ast.DropLogfileGroupStmt:
			report(v.Loc, "logfile group", ident(v.Name))
		case *ast.DropSpatialRefSysStmt:
			report(v.Loc, "spatial reference system", strconv.FormatInt(v.SRID, 10))
		case *ast.DropResourceGroupStmt:
			report(v.Loc, "resource group", ident(v.Name))
		case *ast.AlterTableStmt:
			for _, cmd := range v.Commands {
				if cmd == nil {
					continue
				}
				switch dropKindOf(sql, cmd) {
				case dropColumn:
					reportAt(rangeOf(cmd.Loc), "column", ident(cmd.Name)+" of "+table(v.Table))
				case dropIndex:
					reportAt(rangeOf(cmd.Loc), "index", ident(cmd.Name)+" of "+table(v.Table))
				case dropPartition:
					names := make([]string, len(cmd.PartitionNames))
					for i, name := range cmd.PartitionNames {
						names[i] = ident(name)
					}
					reportAt(rangeOf(cmd.Loc), "partition", strings.Join(names, ", ")+" of "+table(v.Table))
				case dropSystemVersioning:
					reportAt(rangeOf(cmd.Loc), "system versioning", "of "+table(v.Table))
				case dropPeriod:
					reportAt(rangeOf(cmd.Loc), "period", ident(cmd.PeriodName)+" of "+table(v.Table))
				}
			}
		}
		return true
	})
}
