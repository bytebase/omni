package v1review

import (
	"maps"
	"slices"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// TableRequirePK reports every table the statements leave without a primary
// key, as TablePrimaryKeyMissing; Name is the table name as written, without
// its database, and statements nested in a routine body are not judged. It
// follows the statements in order, keeping each table's primary key columns
// by its name as written, compared exactly:
//
//   - CREATE TABLE starts the table with no key, then takes the key of each
//     column declared PRIMARY KEY and of each PRIMARY KEY constraint, the last
//     one winning. CREATE TABLE ... LIKE and ... AS SELECT therefore leave it
//     with none.
//   - DROP TABLE forgets the table.
//   - ALTER TABLE ... ADD PRIMARY KEY sets the key, and a column added,
//     changed or modified with PRIMARY KEY makes it that one column. DROP
//     PRIMARY KEY, a DROP CONSTRAINT with no name, and DROP INDEX named
//     PRIMARY in any letter case empty it.
//   - ALTER TABLE ... DROP COLUMN, and CHANGE COLUMN under its old name,
//     remove the column from the key. A table no earlier statement keyed
//     takes its key from original's PRIMARY index first, the table and index
//     names compared as original says; one with no such index is left alone.
//     CHANGE COLUMN then adds the new name if the old one was in the key.
//
// Each finding anchors where the V1 advisor took its line: the CREATE TABLE
// that started the table, or the last ALTER TABLE that emptied its key or
// dropped a key column, Statement naming which of stmts it is. A table
// left keyless by neither, through a primary key of no plain columns in
// original or added by ALTER TABLE, has Statement -1 and no range. Findings are in table name order;
// the V1 advisor's order came from map iteration and varied from run to run.
// Ported from Bytebase's V1 rule TABLE_REQUIRE_PK.
//
// A nil original panics at the first column dropped from a table no earlier
// statement keyed, as the V1 advisor did with no original metadata.
func TableRequirePK(stmts []ast.Node, original *v1review.Schema) []v1review.Finding {
	type anchor struct{ statement, start, end int }
	keys := map[string]map[string]bool{}
	anchors := map[string]anchor{}
	keyOf := func(columns []string) map[string]bool {
		key := make(map[string]bool)
		for _, column := range columns {
			key[column] = true
		}
		return key
	}
	dropColumn := func(table, column string) bool {
		if _, ok := keys[table]; !ok {
			pk := original.TableIndex("", table, "PRIMARY", false)
			if pk == nil {
				return false
			}
			keys[table] = keyOf(pk.GetExpressions())
		}
		_, inKey := keys[table][column]
		delete(keys[table], column)
		return inKey
	}
	keyColumn := func(table string, column *ast.ColumnDef) {
		if column == nil {
			return
		}
		for _, c := range column.Constraints {
			if c.Type == ast.ColConstrPrimaryKey {
				keys[table] = keyOf([]string{column.Name})
			}
		}
	}

	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			table := n.Table.Name
			keys[table] = make(map[string]bool)
			start, end := span(n.Loc)
			anchors[table] = anchor{i, start, end}
			for _, column := range n.Columns {
				keyColumn(table, column)
			}
			for _, constraint := range n.Constraints {
				if constraint.Type == ast.ConstrPrimaryKey {
					keys[table] = keyOf(constraint.Columns)
				}
			}
		case *ast.DropTableStmt:
			for _, ref := range n.Tables {
				if ref != nil {
					delete(keys, ref.Name)
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			table := n.Table.Name
			start, end := span(n.Loc)
			here := anchor{i, start, end}
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATAddConstraint:
					if cmd.Constraint != nil && cmd.Constraint.Type == ast.ConstrPrimaryKey {
						keys[table] = keyOf(cmd.Constraint.Columns)
					}
				case ast.ATDropConstraint:
					if strings.EqualFold(cmd.Name, "PRIMARY") || cmd.Name == "" {
						keys[table] = make(map[string]bool)
						anchors[table] = here
					}
				case ast.ATDropIndex:
					if strings.EqualFold(cmd.Name, "PRIMARY") {
						keys[table] = make(map[string]bool)
						anchors[table] = here
					}
				case ast.ATAddColumn:
					if cmd.Column != nil {
						keyColumn(table, cmd.Column)
					}
					for _, column := range cmd.Columns {
						keyColumn(table, column)
					}
				case ast.ATChangeColumn:
					if cmd.Column != nil {
						if dropColumn(table, cmd.Name) {
							keys[table][cmd.Column.Name] = true
							anchors[table] = here
						}
						keyColumn(table, cmd.Column)
					}
				case ast.ATModifyColumn:
					if cmd.Column != nil {
						keyColumn(table, cmd.Column)
					}
				case ast.ATDropColumn:
					if dropColumn(table, cmd.Name) {
						anchors[table] = here
					}
				default:
				}
			}
		default:
		}
	}

	var findings []v1review.Finding
	for _, table := range slices.Sorted(maps.Keys(keys)) {
		if len(keys[table]) != 0 {
			continue
		}
		a, ok := anchors[table]
		if !ok {
			a = anchor{statement: -1}
		}
		findings = append(findings, v1review.Finding{Code: v1review.TablePrimaryKeyMissing, Name: table, Statement: a.statement, Start: a.start, End: a.end})
	}
	return findings
}
