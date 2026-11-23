package catalog

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
	mariadbparser "github.com/bytebase/omni/mariadb/parser"
)

// checksDroppedWithColumns decides, before an ALTER TABLE runs, what becomes
// of each table-level CHECK when the statement drops columns, as MariaDB's
// mysql_prepare_alter_table does (sql/sql_table.cc, 11.8.8): a check the
// statement drops by name is left to that DROP CONSTRAINT; one that reads only
// columns the statement drops goes with them; and one that reads a dropped
// column beside a kept one is 1054, "Unknown column '<col>' in 'CHECK'",
// naming the first of its dropped columns in the table's order.
//
// MariaDB reads the columns the table had before the statement, so a rename
// or a re-added column changes nothing: measured on mariadb:11.8.8 with
// CHECK (starts <= ends), DROP COLUMN starts, DROP COLUMN ends drops the
// check, DROP COLUMN starts is 1054, and so are RENAME COLUMN starts TO
// begins, DROP COLUMN ends and DROP COLUMN starts, ADD COLUMN starts DATE. A
// check that reads no column is a subset of any dropped set, so it goes with
// whatever column the statement drops.
//
// A column's own check goes with its column (dropColumnCheck) and is not
// looked at here.
func checksDroppedWithColumns(tbl *Table, cmds []*nodes.AlterTableCmd) (map[*Constraint]bool, error) {
	dropped := map[string]bool{}
	named := map[string]bool{}
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case nodes.ATDropColumn:
			if tbl.GetColumn(cmd.Name) != nil {
				dropped[toLower(cmd.Name)] = true
			}
		case nodes.ATDropConstraint:
			named[toLower(cmd.Name)] = true
		}
	}
	if len(dropped) == 0 {
		return nil, nil
	}
	gone := map[*Constraint]bool{}
	for _, con := range tbl.Constraints {
		if con.Type != ConCheck || isColumnCheck(con) || named[toLower(con.Name)] {
			continue
		}
		first, kept := "", false
		for _, col := range tbl.Columns {
			if !checkReadsColumn(con.CheckExpr, col.Name) {
				continue
			}
			if !dropped[toLower(col.Name)] {
				kept = true
			} else if first == "" {
				first = col.Name
			}
		}
		switch {
		case first == "":
			if !kept {
				gone[con] = true
			}
		case kept:
			return nil, &Error{Code: 1054, SQLState: "42S22", Message: "Unknown column '" + first + "' in 'CHECK'"}
		default:
			gone[con] = true
		}
	}
	return gone, nil
}

// dropChecks removes the constraints in gone from the table.
func dropChecks(tbl *Table, gone map[*Constraint]bool) {
	if len(gone) == 0 {
		return
	}
	kept := tbl.Constraints[:0]
	for _, con := range tbl.Constraints {
		if !gone[con] {
			kept = append(kept, con)
		}
	}
	tbl.Constraints = kept
}

// checkReadsColumn reports whether a stored CHECK expression reads the named
// column. The expression is parsed, so a string holding the name does not
// count; one the parser does not take falls back to looking for the quoted
// name, as the catalog stores it.
func checkReadsColumn(expr, column string) bool {
	const prefix = "SELECT ("
	list, err := mariadbparser.Parse(prefix + expr + ")")
	if err != nil || len(list.Items) != 1 {
		return generatedExprReferencesColumn(expr, column)
	}
	found := false
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		if ref, ok := n.(*nodes.ColumnRef); ok && !ref.Star && toLower(ref.Column) == toLower(column) {
			found = true
		}
		return !found
	})
	return found
}
