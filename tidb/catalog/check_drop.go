package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/tidb/ast"
	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// TiDB v8.5.5 ties a CHECK constraint to the columns it reads
// (ConstraintInfo.ConstraintCols), judged against the table as it was before
// the statement:
//
//   - DROP COLUMN of a column that a CHECK reading two or more columns reads
//     is 3959, whatever else the statement drops (IsColumnDroppableWithCheck
//     Constraint, pkg/ddl/constraint.go);
//   - RENAME COLUMN of a column any CHECK reads is 3959, even to its own name
//     (IsColumnRenameableWithCheckConstraint, called by RenameColumn);
//   - a CHECK that reads a column the table no longer has is removed when the
//     table is next loaded (removeInvalidCheckConstraintsInfo,
//     pkg/table/constraint.go): so goes a one-column CHECK with its column,
//     and a CHECK whose column CHANGE COLUMN renamed.
//
// Measured on v8.5.5 with tidb_enable_check_constraint on (features@check in
// testdata/simulate.tsv): on CHECK (starts <= ends), DROP COLUMN starts, DROP
// COLUMN ends is 3959 naming starts, and so is RENAME COLUMN starts TO begins;
// DROP COLUMN nights drops CHECK (nights > 0), and CHANGE COLUMN nights stays
// INT drops it too.

// checkReads returns the lowered names of the columns a stored CHECK
// expression reads, as ConstraintCols holds them. The expression is parsed,
// so a string holding a name does not count; one the parser does not take
// falls back to looking for each column's quoted name.
func checkReads(tbl *Table, con *Constraint) map[string]bool {
	reads := map[string]bool{}
	const prefix = "SELECT ("
	list, err := tidbparser.Parse(prefix + con.CheckExpr + ")")
	if err != nil || len(list.Items) != 1 {
		for _, col := range tbl.Columns {
			if generatedExprReferencesColumn(con.CheckExpr, col.Name) {
				reads[toLower(col.Name)] = true
			}
		}
		return reads
	}
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		if ref, ok := n.(*nodes.ColumnRef); ok && !ref.Star {
			reads[toLower(ref.Column)] = true
		}
		return true
	})
	return reads
}

// checkDropColumnWithCheck refuses dropping a column that a CHECK reading
// more than one column reads: 3959, naming the column as the statement wrote
// it.
func checkDropColumnWithCheck(tbl *Table, column string) error {
	key := toLower(column)
	for _, con := range tbl.Constraints {
		if con.Type != ConCheck {
			continue
		}
		if reads := checkReads(tbl, con); len(reads) > 1 && reads[key] {
			return errCantDropColWithCheck(con.Name, column)
		}
	}
	return nil
}

// checkRenameColumnWithCheck refuses renaming a column any CHECK reads: 3959,
// naming the column as the table holds it.
func checkRenameColumnWithCheck(tbl *Table, column string) error {
	key := toLower(column)
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck && checkReads(tbl, con)[key] {
			return errCantDropColWithCheck(con.Name, column)
		}
	}
	return nil
}

// dropChecksOfMissingColumns removes every CHECK that reads a column the
// table no longer has, as TiDB does when it next loads the table.
func dropChecksOfMissingColumns(tbl *Table) {
	kept := tbl.Constraints[:0]
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck && !checkColumnsExist(tbl, con) {
			continue
		}
		kept = append(kept, con)
	}
	tbl.Constraints = kept
}

func checkColumnsExist(tbl *Table, con *Constraint) bool {
	for name := range checkReads(tbl, con) {
		if _, ok := tbl.colByName[name]; !ok {
			return false
		}
	}
	return true
}

// errCantDropColWithCheck is 3959, ErrCantDropColWithCheckConstraint.
func errCantDropColWithCheck(constraint, column string) error {
	return &Error{Code: 3959, SQLState: "HY000", Message: fmt.Sprintf(
		"Check constraint '%s' uses column '%s', hence column cannot be dropped or renamed.", constraint, column)}
}

// errMultiSchemaDropCheck is 8200: TiDB runs DROP CHECK as a job of its own,
// so a multi-command ALTER that holds one is refused once the check is found.
func errMultiSchemaDropCheck() error {
	return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL),
		Message: "Unsupported multi schema change for drop check constraint"}
}

// isDropOfCheck reports whether cmd is DROP CHECK or DROP CONSTRAINT naming a
// CHECK the table has.
func isDropOfCheck(tbl *Table, cmd *nodes.AlterTableCmd) bool {
	if cmd == nil || cmd.Type != nodes.ATDropConstraint || cmd.PrimaryKey || cmd.ForeignKey {
		return false
	}
	key := toLower(cmd.Name)
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck && toLower(con.Name) == key {
			return true
		}
	}
	return false
}
