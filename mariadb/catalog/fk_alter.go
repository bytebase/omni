package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// What an ALTER TABLE may not do to the columns of a foreign key while
// foreign_key_checks is on, read on mariadb:11.8.8 over features@fk, whose
// orders.customer_id references customers.id ON DELETE CASCADE ON UPDATE SET
// NULL:
//
//   - MODIFY orders.customer_id BIGINT is 1832, a column of the table's own
//     foreign key that changes its type;
//   - MODIFY customers.id BIGINT is 1833, a column another table's foreign key
//     references;
//   - MODIFY orders.customer_id INT NOT NULL is 1830, a column a SET NULL action
//     writes;
//   - DROP COLUMN customers.id is 1829, a column another table's foreign key
//     references.
//
// A column redefined with the type it had (a rename, a new default) passes.

// fkColumnSnapshot is a table's columns before an ALTER TABLE, by lowered name.
type fkColumnSnapshot map[string]*Column

func snapshotColumns(tbl *Table) fkColumnSnapshot {
	cols := make(fkColumnSnapshot, len(tbl.Columns))
	for _, col := range tbl.Columns {
		cols[toLower(col.Name)] = col
	}
	return cols
}

// checkForeignKeyColumnChanges holds the table an ALTER TABLE left to the
// rules above. before holds its columns before the statement, renames the
// columns it renamed (alterColumnRenames), and existing its constraints
// before it; the table was oldName in oldDB.
func (c *Catalog) checkForeignKeyColumnChanges(tbl *Table, oldDB *Database, oldName string, before fkColumnSnapshot, renames map[string]string, existing map[*Constraint]bool) error {
	if !c.foreignKeyChecks {
		return nil
	}
	oldNameOf := make(map[string]string, len(renames))
	for old, now := range renames {
		oldNameOf[toLower(now)] = old
	}
	// was returns the column name had before the statement, and is the one it
	// names now: name is as the statement left it.
	was := func(name string) *Column {
		if old, ok := oldNameOf[toLower(name)]; ok {
			return before[old]
		}
		if _, renamedAway := renames[toLower(name)]; renamedAway {
			return nil
		}
		return before[toLower(name)]
	}
	for _, con := range tbl.Constraints {
		if con.Type != ConForeignKey || !existing[con] {
			continue
		}
		for _, name := range con.Columns {
			old, now := was(name), tbl.GetColumn(name)
			if old == nil || now == nil || old == now {
				continue
			}
			if columnTypeChanged(tbl, old, now) {
				return &Error{Code: 1832, SQLState: "HY000", Message: fmt.Sprintf(
					"Cannot change column '%s': used in a foreign key constraint '%s'", now.Name, con.Name)}
			}
			if old.Nullable && !now.Nullable && (strings.EqualFold(con.OnDelete, "SET NULL") || strings.EqualFold(con.OnUpdate, "SET NULL")) {
				return &Error{Code: 1830, SQLState: "HY000", Message: fmt.Sprintf(
					"Column '%s' cannot be NOT NULL: needed in a foreign key constraint '%s' SET NULL", now.Name, con.Name)}
			}
		}
	}
	for _, ref := range c.referencingForeignKeysOf(oldDB.Name, oldName) {
		if ref.child == tbl {
			continue
		}
		for _, name := range ref.con.RefColumns {
			old := before[toLower(name)]
			if old == nil {
				continue
			}
			nowName := name
			if renamed, ok := renames[toLower(name)]; ok {
				nowName = renamed
			}
			now := tbl.GetColumn(nowName)
			if now == nil {
				return &Error{Code: 1829, SQLState: "HY000", Message: fmt.Sprintf(
					"Cannot drop column '%s': needed in a foreign key constraint '%s/%s' of table %s.%s",
					old.Name, ref.child.Database.Name, ref.con.Name, backtickIDQ(ref.child.Database.Name), backtickIDQ(ref.child.Name))}
			}
			if now != old && columnTypeChanged(tbl, old, now) {
				return &Error{Code: 1833, SQLState: "HY000", Message: fmt.Sprintf(
					"Cannot change column '%s': used in a foreign key constraint '%s' of table '%s.%s'",
					now.Name, ref.con.Name, ref.child.Database.Name, ref.child.Name)}
			}
		}
	}
	return nil
}

// columnTypeChanged reports whether a redefinition changed what a column
// stores: its type, or a string column's character set or collation.
func columnTypeChanged(tbl *Table, old, now *Column) bool {
	if !strings.EqualFold(old.ColumnType, now.ColumnType) {
		return true
	}
	if !isStringType(old.DataType) {
		return false
	}
	charset := func(col *Column) string {
		if col.Charset != "" {
			return toLower(col.Charset)
		}
		return toLower(tbl.Charset)
	}
	collation := func(col *Column) string {
		if col.Collation != "" {
			return toLower(col.Collation)
		}
		return toLower(tbl.Collation)
	}
	return charset(old) != charset(now) || collation(old) != collation(now)
}

// alterLeavesEngine returns the engine an ALTER TABLE leaves the table in.
func alterLeavesEngine(tbl *Table, commands []*nodes.AlterTableCmd) string {
	engine := tbl.Engine
	for _, cmd := range commands {
		if cmd != nil && cmd.Type == nodes.ATTableOption && cmd.Option != nil &&
			(eqFoldStr(cmd.Option.Name, "ENGINE") || eqFoldStr(cmd.Option.Name, "STORAGE ENGINE")) {
			engine = canonicalEngineName(cmd.Option.Value)
		}
	}
	return engine
}

// alterMovesForeignKeyTable refuses moving a table that has a foreign key, or
// that another table's references, off InnoDB while foreign_key_checks is on:
// mariadb:11.8.8 answers 1217 for ALTER TABLE orders ENGINE = MyISAM and for
// customers, the table orders references.
func (c *Catalog) alterMovesForeignKeyTable(tbl *Table, commands []*nodes.AlterTableCmd) error {
	if !c.foreignKeyChecks {
		return nil
	}
	from, to := alterLeavesEngine(tbl, nil), alterLeavesEngine(tbl, commands)
	if !isInnoDB(from) || isInnoDB(to) || eqFoldStr(to, "DEFAULT") {
		return nil
	}
	participates := len(c.referencingForeignKeys(tbl)) > 0
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			participates = true
		}
	}
	if participates {
		return errRowIsReferencedOld()
	}
	return nil
}

// checkInvisibleColumns holds a table a statement leaves to the rules of
// invisible columns, measured on mariadb:11.8.8: a column the statement
// defined (before does not hold it) that is INVISIBLE and NOT NULL takes a
// default, an AUTO_INCREMENT or a generated value (4108), and the table keeps
// a column that is not invisible (1113, also when a DROP COLUMN leaves only
// invisible ones).
func checkInvisibleColumns(tbl *Table, before map[*Column]bool) error {
	for _, col := range tbl.Columns {
		if before[col] || !col.Invisible || col.Nullable || col.AutoIncrement || col.Generated != nil {
			continue
		}
		if col.Default == nil {
			return &Error{Code: 4108, SQLState: "HY000", Message: fmt.Sprintf("Invisible column `%s` must have a default value", col.Name)}
		}
	}
	for _, col := range tbl.Columns {
		if !col.Invisible && col.Hidden == ColumnHiddenNone {
			return nil
		}
	}
	return errTableMustHaveColumns()
}
