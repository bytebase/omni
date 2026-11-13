package v1review

import (
	"maps"
	"slices"
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// TableRequirePK reports every table the statements leave without a primary
// key, as TablePrimaryKeyMissing; Name is the table name as written, without
// its database, and statements nested in a routine body are not judged. It
// follows the statements in order, keeping each table's primary key columns
// by its name as written, compared exactly:
//
//   - CREATE TABLE starts the table with no key, then takes the key of the
//     first column declared PRIMARY KEY and of each PRIMARY KEY constraint, the
//     last one winning. CREATE TABLE ... LIKE takes the key of the table it
//     copies: that table's key so far if an earlier statement keyed it, else
//     its primary key in original, the name compared as original says, else
//     none. CREATE TABLE ... AS SELECT is left with no key.
//   - DROP TABLE forgets the table.
//   - ALTER TABLE ... ADD PRIMARY KEY sets the key, and the first column
//     added, changed or modified with PRIMARY KEY makes it that one column.
//     DROP PRIMARY KEY, a DROP CONSTRAINT named PRIMARY in any letter case or
//     with no name, and a DROP INDEX whose name upper-cases to PRIMARY empty
//     it.
//   - ALTER TABLE ... DROP COLUMN, and CHANGE COLUMN under its old name,
//     remove the column from the key. A table no earlier statement keyed
//     takes its key from original's PRIMARY index first, the table and index
//     names compared as original says; one with no such index is left alone.
//     CHANGE COLUMN then adds the new name if the old one was in the key.
//
// Each finding anchors where the V1 advisor took its line: the CREATE TABLE
// that started the table, or the last ALTER TABLE that emptied its key or
// dropped a key column, Statement naming which of stmts it is. A table left
// keyless by neither, through a primary key of no plain columns in original
// or added by ALTER TABLE, has Statement -1 and no range. Findings are in
// table name order, as the V1 advisor sorted them. Ported from Bytebase's V1
// rule TABLE_REQUIRE_PK.
//
// A nil original panics at the first lookup in it, as the V1 advisor did
// with no original metadata.
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
	keyFirstColumn := func(table string, columns []*ast.ColumnDef) {
		for _, column := range columns {
			if column == nil {
				continue
			}
			for _, c := range column.Constraints {
				if c != nil && c.Type == ast.ColConstrPrimaryKey {
					keys[table] = keyOf([]string{column.Name})
					return
				}
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
			if n.Like != nil {
				if key, ok := keys[n.Like.Name]; ok {
					keys[table] = keyOf(slices.Collect(maps.Keys(key)))
				} else if pk := original.TablePrimaryKey("", n.Like.Name, false); pk != nil {
					keys[table] = keyOf(pk.GetExpressions())
				}
			} else {
				keyFirstColumn(table, n.Columns)
				for _, constraint := range n.Constraints {
					if constraint != nil && constraint.Type == ast.ConstrPrimaryKey {
						keys[table] = keyOf(constraint.Columns)
					}
				}
			}
			start, end := span(n.Loc)
			anchors[table] = anchor{i, start, end}
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
					if strings.ToUpper(cmd.Name) == "PRIMARY" {
						keys[table] = make(map[string]bool)
						anchors[table] = here
					}
				case ast.ATAddColumn:
					if len(cmd.Columns) > 0 {
						keyFirstColumn(table, cmd.Columns)
					} else if cmd.Column != nil {
						keyFirstColumn(table, []*ast.ColumnDef{cmd.Column})
					}
				case ast.ATChangeColumn:
					if cmd.Column == nil {
						continue
					}
					if dropColumn(table, cmd.Name) {
						keys[table][cmd.Column.Name] = true
						anchors[table] = here
					}
					keyFirstColumn(table, []*ast.ColumnDef{cmd.Column})
				case ast.ATModifyColumn:
					if cmd.Column != nil {
						keyFirstColumn(table, []*ast.ColumnDef{cmd.Column})
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
