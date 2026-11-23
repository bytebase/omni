package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Checks of generated columns, CHECK constraints and expression keys that the
// ALTER TABLE path lacked, read on v8.5.5 over simulate.tsv's
// features@generated fixture.

// columnInAnyIndex reports whether a key of the table reads the column.
func columnInAnyIndex(tbl *Table, name string) bool {
	for _, idx := range tbl.Indexes {
		for _, ic := range idx.Columns {
			if strings.EqualFold(ic.Name, name) {
				return true
			}
		}
	}
	return false
}

// checkMultiCmdGenerated holds the commands of an ALTER TABLE to the
// generated-column rules TiDB judges on the table the statement began with,
// before any command runs: an ADD COLUMN of a STORED generated column is 3106
// ahead of a key the statement adds over it (ADD COLUMN g ... STORED, ADD
// INDEX idx_g (g)), and a DROP COLUMN of a column a generated column reads is
// 3108 even when the statement drops that generated column too (DROP COLUMN
// full_name, DROP COLUMN name_len, DROP COLUMN first_name).
func checkMultiCmdGenerated(tbl *Table, cmds []*nodes.AlterTableCmd) error {
	for _, cmd := range cmds {
		if cmd == nil || cmd.Type != nodes.ATAddColumn {
			continue
		}
		defs := cmd.Columns
		if cmd.Column != nil {
			defs = append([]*nodes.ColumnDef{cmd.Column}, defs...)
		}
		for _, def := range defs {
			if def != nil && def.Generated != nil && def.Generated.Stored {
				return errAlterStoredGeneratedColumn()
			}
		}
	}
	for _, cmd := range cmds {
		if cmd == nil || cmd.Type != nodes.ATDropColumn || tbl.GetColumn(cmd.Name) == nil {
			continue
		}
		for _, gc := range tbl.Columns {
			if gc.Generated != nil && !strings.EqualFold(gc.Name, cmd.Name) && generatedExprReferencesColumn(gc.Generated.Expr, cmd.Name) {
				return errDependentByGeneratedColumn(cmd.Name, gc.Name, tbl.Name)
			}
		}
	}
	return nil
}

// validateCheckConstraint holds a CHECK a statement adds to what v8.5.5 takes:
// a name no CHECK of the schema holds (3822: CHECK constraint names are the
// schema's, so CREATE TABLE chk3 (..., CONSTRAINT chk_age CHECK (a > 0)) beside
// people's chk_age is refused), columns the table has (1054 in 'check
// constraint <name> expression'), no function whose value changes between rows
// (3814), and no AUTO_INCREMENT column (3818). tbl holds the statement's
// columns.
func validateCheckConstraint(tbl *Table, name string, expr nodes.ExprNode) error {
	if db := tbl.Database; db != nil {
		for _, other := range db.Tables {
			if other == tbl {
				continue
			}
			for _, con := range other.Constraints {
				if con.Type == ConCheck && strings.EqualFold(con.Name, name) {
					return errDupCheckName(name)
				}
			}
		}
	}
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck && strings.EqualFold(con.Name, name) {
			return errDupCheckName(name)
		}
	}
	var refused error
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if refused != nil {
			return false
		}
		switch e := n.(type) {
		case *nodes.ColumnRef:
			if e.Column == "" || e.Star {
				return true
			}
			col := tbl.GetColumn(e.Column)
			switch {
			case col == nil:
				refused = &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
					Message: fmt.Sprintf("Unknown column '%s' in 'check constraint %s expression'", e.Column, name)}
			case col.AutoIncrement:
				refused = &Error{Code: 3818, SQLState: "HY000",
					Message: fmt.Sprintf("Check constraint '%s' cannot refer to an auto-increment column.", name)}
			}
		case *nodes.FuncCallExpr:
			if disallowedGeneratedColumnFuncs[strings.ToLower(e.Name)] {
				refused = &Error{Code: 3814, SQLState: "HY000",
					Message: fmt.Sprintf("An expression of a check constraint '%s' contains disallowed function: %s.", name, strings.ToLower(e.Name))}
			}
		}
		return refused == nil
	})
	return refused
}

func errDupCheckName(name string) error {
	return &Error{Code: 3822, SQLState: "HY000", Message: fmt.Sprintf("Duplicate check constraint name '%s'.", name)}
}

// checkFunctionalKeyColumns refuses an expression key part that names a
// column the table does not have: 1054 in 'expression'.
func checkFunctionalKeyColumns(tbl *Table, parts []*nodes.IndexColumn) error {
	for _, ic := range parts {
		if ic == nil || !ic.Functional {
			continue
		}
		var missing string
		nodes.Inspect(ic.Expr, func(n nodes.Node) bool {
			if ref, ok := n.(*nodes.ColumnRef); ok && missing == "" && ref.Column != "" && tbl.GetColumn(ref.Column) == nil {
				missing = ref.Column
			}
			return missing == ""
		})
		if missing != "" {
			return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
				Message: fmt.Sprintf("Unknown column '%s' in 'expression'", missing)}
		}
	}
	return nil
}

// checkGeneratedColumnOrder refuses a CREATE TABLE's generated column that
// reads a generated column defined after it (3107) or an AUTO_INCREMENT
// column (3109), as v8.5.5 does (gc3, gc4). It reads the column references of
// the statement's expressions, so a string literal naming a column reads none.
func checkGeneratedColumnOrder(tbl *Table, stmt *nodes.CreateTableStmt) error {
	position := map[string]int{}
	for i, def := range stmt.Columns {
		if def != nil {
			position[toLower(def.Name)] = i
		}
	}
	for i, def := range stmt.Columns {
		if def == nil || def.Generated == nil || def.Generated.Expr == nil {
			continue
		}
		var refused error
		nodes.Inspect(def.Generated.Expr, func(n nodes.Node) bool {
			ref, ok := n.(*nodes.ColumnRef)
			if !ok || refused != nil || ref.Column == "" {
				return refused == nil
			}
			other := tbl.GetColumn(ref.Column)
			if other == nil || strings.EqualFold(other.Name, def.Name) {
				return true
			}
			switch j := position[toLower(other.Name)]; {
			case other.AutoIncrement:
				refused = &Error{Code: 3109, SQLState: "HY000",
					Message: fmt.Sprintf("Generated column '%s' cannot refer to auto-increment column.", def.Name)}
			case other.Generated != nil && j > i:
				refused = &Error{Code: 3107, SQLState: "HY000",
					Message: "Generated column can refer only to generated columns defined prior to it."}
			}
			return refused == nil
		})
		if refused != nil {
			return refused
		}
	}
	return nil
}

// checkColumnLength refuses a string column longer than its type holds: 1074,
// measured on v8.5.5 under the default sql_mode, for CHAR over 255 characters,
// VARBINARY over 65535 bytes and VARCHAR over 65535 bytes in its character set
// (max = 16383 in utf8mb4). Outside a strict sql_mode a VARCHAR that long is
// made a TEXT instead (promoteDeclaredType), and is not refused.
func (c *Catalog) checkColumnLength(colDef *nodes.ColumnDef, charset string) error {
	dt := colDef.TypeName
	if dt == nil || dt.Length <= 0 {
		return nil
	}
	tooLong := func(max int) error {
		return &Error{Code: 1074, SQLState: "42000",
			Message: fmt.Sprintf("Column length too big for column '%s' (max = %d); use BLOB or TEXT instead", colDef.Name, max)}
	}
	switch strings.ToLower(dt.Name) {
	case "char", "character":
		if dt.Length > 255 {
			return tooLong(255)
		}
	case "varbinary":
		if dt.Length > 65535 && !c.nonStrict {
			return tooLong(65535)
		}
	case "varchar", "character varying":
		mb := charsetMaxLen(charset)
		if mb < 1 {
			mb = 1
		}
		if dt.Length*mb > 65535 && !c.nonStrict {
			return tooLong(65535 / mb)
		}
	}
	return nil
}

// alterAlgorithmIs reports whether an ALTER TABLE's last ALGORITHM clause is
// the one named.
func alterAlgorithmIs(cmds []*nodes.AlterTableCmd, name string) bool {
	algorithm := ""
	for _, cmd := range cmds {
		if cmd != nil && cmd.Type == nodes.ATAlgorithm {
			algorithm = cmd.Name
		}
	}
	return strings.EqualFold(strings.TrimSpace(algorithm), name)
}
