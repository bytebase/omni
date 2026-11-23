package catalog

import (
	"slices"
	"strings"

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
		refs, parsed := checkColumnRefs(con.CheckExpr)
		first, kept := "", false
		for _, col := range tbl.Columns {
			key := toLower(col.Name)
			reads := refs[key]
			if !parsed {
				reads = generatedExprReferencesColumn(con.CheckExpr, col.Name)
			}
			delete(refs, key)
			if !reads {
				continue
			}
			if !dropped[key] {
				kept = true
			} else if first == "" {
				first = col.Name
			}
		}
		// A name the table does not have is one the catalog failed to follow;
		// keep the check rather than take it for one that reads no column.
		kept = kept || len(refs) > 0
		switch {
		case first != "" && kept:
			return nil, &Error{Code: 1054, SQLState: "42S22", Message: "Unknown column '" + first + "' in 'CHECK'"}
		case !kept:
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

// checkColumnRefs returns the lowered names of the columns a stored CHECK
// expression reads, parsing it once, so a string holding a name does not
// count. ok is false when the parser does not take the expression.
func checkColumnRefs(expr string) (refs map[string]bool, ok bool) {
	const prefix = "SELECT ("
	list, err := mariadbparser.Parse(prefix + expr + ")")
	if err != nil || len(list.Items) != 1 {
		return map[string]bool{}, false
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
// that reads it, as MariaDB rewrites the stored expression (rename_fields_
// processor in sql/sql_table.cc): measured on mariadb:11.8.8 with
// features@check, RENAME COLUMN guests TO people leaves the column's check
// reading `people`, and RENAME COLUMN note TO memo the table's chk_note
// reading `memo`. An expression the parser does not take is left as it is.
func renameCheckColumnRefs(tbl *Table, oldName, newName string) {
	const prefix = "SELECT ("
	oldKey := toLower(oldName)
	for _, con := range tbl.Constraints {
		if con.Type != ConCheck || con.Implicit {
			continue
		}
		list, err := mariadbparser.Parse(prefix + con.CheckExpr + ")")
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
		if len(spans) > 0 {
			b.WriteString(con.CheckExpr[prev:])
			con.CheckExpr = b.String()
		}
	}
}
