package catalog

import (
	"fmt"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	obparser "github.com/bytebase/omni/oceanbase/parser"
)

// checksDroppedWithColumns decides, before an ALTER TABLE runs, what the
// CHECK constraints the table had become, as OceanBase 4.3.5.6 judges them in
// MySQL mode:
//
//   - a CHECK reading one column goes with that column
//     (ObDDLResolver::check_column_in_check_constraint);
//   - RENAME COLUMN of a column a CHECK reads is 3959, unless the statement
//     drops that CHECK (ObAlterTableResolver::check_mysql_rename_column);
//   - DROP COLUMN of a column a CHECK reading more columns reads is 3959,
//     naming the first such column in the statement's order, unless the
//     statement drops that CHECK, in any position
//     (ObDDLService::check_drop_column_with_drop_constraint);
//   - a CHECK that reads no column goes when the statement drops a column.
//
// Measured on 4.3.5.6 with features@check in testdata/simulate.tsv: on CHECK
// (starts <= ends), DROP COLUMN starts, DROP COLUMN ends is 3959 naming
// starts and DROP COLUMN starts, DROP COLUMN ends, DROP CHECK chk_dates runs;
// RENAME COLUMN guests TO people is 3959 on a column CHECK; DROP COLUMN
// guests drops that CHECK; and DROP COLUMN beds drops CHECK (1 = 1).
//
// It returns the checks the statement leaves behind once it has run.
func checksDroppedWithColumns(tbl *Table, cmds []*nodes.AlterTableCmd) (map[*Constraint]bool, error) {
	named := map[string]bool{}
	var drops []string
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case nodes.ATDropColumn:
			if col := tbl.GetColumn(cmd.Name); col != nil {
				drops = append(drops, col.Name)
			}
		case nodes.ATDropConstraint:
			// DROP CHECK (a, b) holds its names in ConstraintNames.
			for _, name := range append([]string{cmd.Name}, cmd.ConstraintNames...) {
				if name != "" {
					named[toLower(name)] = true
				}
			}
		}
	}
	type check struct {
		con   *Constraint
		reads map[string]bool
	}
	var checks []check
	for _, con := range tbl.Constraints {
		if con.Type != ConCheck {
			continue
		}
		reads, ok := checkColumnRefs(con.CheckExpr)
		if !ok {
			continue
		}
		checks = append(checks, check{con, reads})
	}

	gone := map[*Constraint]bool{}
	for _, ck := range checks {
		if named[toLower(ck.con.Name)] {
			continue
		}
		switch len(ck.reads) {
		case 0:
			if len(drops) > 0 {
				gone[ck.con] = true
			}
		case 1:
			for _, col := range drops {
				if ck.reads[toLower(col)] {
					gone[ck.con] = true
				}
			}
		}
	}
	kept := func(ck check) bool { return !named[toLower(ck.con.Name)] && !gone[ck.con] }

	for _, cmd := range cmds {
		if cmd == nil || cmd.Type != nodes.ATRenameColumn || tbl.GetColumn(cmd.Name) == nil {
			continue
		}
		for _, ck := range checks {
			if kept(ck) && ck.reads[toLower(cmd.Name)] {
				return nil, errCantDropColWithCheck(ck.con.Name, cmd.Name)
			}
		}
	}
	for _, col := range drops {
		for _, ck := range checks {
			if kept(ck) && ck.reads[toLower(col)] {
				return nil, errCantDropColWithCheck(ck.con.Name, col)
			}
		}
	}
	return gone, nil
}

// dropChecks removes the constraints in gone from the table.
func dropChecks(tbl *Table, gone map[*Constraint]bool) {
	if len(gone) == 0 {
		return
	}
	tbl.Constraints = slices.DeleteFunc(tbl.Constraints, func(con *Constraint) bool { return gone[con] })
}

// errCantDropColWithCheck is 3959, OB_ERR_DROP_COL_REFERENCED_MULTI_COLS_CONSTRAINT.
func errCantDropColWithCheck(constraint, column string) error {
	return &Error{Code: 3959, SQLState: "HY000", Message: fmt.Sprintf(
		"Check constraint '%s' uses column '%s', hence column cannot be dropped or renamed.", constraint, column)}
}

// checkColumnRefs returns the lowered names of the columns a stored CHECK
// expression reads, parsing it once, so a string holding a name does not
// count. ok is false when the parser does not take the expression.
func checkColumnRefs(expr string) (refs map[string]bool, ok bool) {
	const prefix = "SELECT ("
	list, err := obparser.Parse(prefix + expr + ")")
	if err != nil || len(list.Items) != 1 {
		return nil, false
	}
	refs = map[string]bool{}
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		if ref, ok := n.(*nodes.ColumnRef); ok && !ref.Star {
			refs[toLower(ref.Column)] = true
		}
		return true
	})
	return refs, true
}

// renameCheckColumnRefs follows a renamed column into the text of every CHECK
// that reads it, as ObDDLService::rebuild_constraint_check_expr does for an
// altered column: measured on 4.3.5.6, CHANGE COLUMN nights stays INT leaves
// CHECK (nights > 0) reading `stays`. RENAME COLUMN of a column a CHECK reads
// is refused (checksDroppedWithColumns), so only CHANGE reaches this with a
// CHECK to rewrite. An expression the parser does not take is left as it is.
func renameCheckColumnRefs(tbl *Table, oldName, newName string) {
	const prefix = "SELECT ("
	oldKey := toLower(oldName)
	if oldKey == toLower(newName) {
		return
	}
	for _, con := range tbl.Constraints {
		if con.Type != ConCheck {
			continue
		}
		list, err := obparser.Parse(prefix + con.CheckExpr + ")")
		if err != nil || len(list.Items) != 1 {
			continue
		}
		type span struct {
			start, end int
			name       string
		}
		var spans []span
		nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
			ref, ok := n.(*nodes.ColumnRef)
			if !ok || ref.Star || toLower(ref.Column) != oldKey {
				return true
			}
			start, end := ref.Loc.Start-len(prefix), ref.Loc.End-len(prefix)
			if start < 0 || end > len(con.CheckExpr) || start >= end {
				return true
			}
			name := backtickIDQ(newName)
			if ref.Table != "" {
				name = backtickIDQ(ref.Table) + "." + name
				if ref.Schema != "" {
					name = backtickIDQ(ref.Schema) + "." + name
				}
			}
			spans = append(spans, span{start, end, name})
			return true
		})
		if len(spans) == 0 {
			continue
		}
		slices.SortFunc(spans, func(a, b span) int { return a.start - b.start })
		var b strings.Builder
		prev := 0
		for _, sp := range spans {
			if sp.start < prev {
				break
			}
			b.WriteString(con.CheckExpr[prev:sp.start])
			b.WriteString(sp.name)
			prev = sp.end
		}
		b.WriteString(con.CheckExpr[prev:])
		con.CheckExpr = b.String()
	}
}
