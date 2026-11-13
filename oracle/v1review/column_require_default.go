package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnRequireDefault judges the top-level statements together for
// columns without a DEFAULT. Each column is keyed table.column, the table
// as schema.name with the schema defaulting to currentDatabase. A column a
// CREATE TABLE defines without a DEFAULT is recorded, anchored at its
// definition, and one with a DEFAULT is forgotten. In an ALTER TABLE, a
// column any action defines (its column list, then its single column) is
// forgotten when it has a DEFAULT; otherwise an ADD action records it and
// any other action, MODIFY included, leaves it as it was. Nothing else
// forgets a column, a DROP COLUMN or DROP TABLE included; identity and
// virtual columns are judged like any other.
//
// Each finding is one recorded column, in key order compared bytewise,
// with code ColumnNoDefault, Name the key's text after its last dot (so a
// column name holding a dot prints only its tail), Statement and the range
// the last definition that recorded it. Ported from Bytebase's V1 rule
// COLUMN_REQUIRE_DEFAULT.
func ColumnRequireDefault(stmts []ast.Node, currentDatabase string) []v1review.Finding {
	noDefault := make(map[string]columnAnchor)
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			table := objectName(n.Name, currentDatabase)
			for _, col := range columnDefs(n.Columns) {
				key := table + "." + col.Name
				if col.Default == nil {
					noDefault[key] = columnAnchor{statement: i, loc: col.Loc}
				} else {
					delete(noDefault, key)
				}
			}
		case *ast.AlterTableStmt:
			table := objectName(n.Name, currentDatabase)
			for _, cmd := range alterTableCmds(n) {
				for _, col := range actionColumns(cmd) {
					if col == nil {
						continue
					}
					key := table + "." + col.Name
					if col.Default != nil {
						delete(noDefault, key)
					} else if cmd.Action == ast.AT_ADD_COLUMN {
						noDefault[key] = columnAnchor{statement: i, loc: col.Loc}
					}
				}
			}
		}
	}
	return keyedColumnFindings(noDefault, v1review.ColumnNoDefault)
}
