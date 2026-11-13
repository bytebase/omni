package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// ColumnDisallowDropInIndex judges the statements together for DROP COLUMN
// actions of top-level ALTER TABLE statements that drop a column an index
// uses. Each such action is one finding, anchored at the action, with Name
// the column and Table the table as written, without a qualifier; Statement
// says which statement holds it. Findings are in statement order, then
// action order.
//
// A column is in an index when a top-level CREATE TABLE earlier in the
// statements lists it among the plain columns of a KEY or INDEX of its own
// (not UNIQUE, PRIMARY KEY, FULLTEXT or SPATIAL), or when any index schema
// holds for the table lists it among its expressions, read in the schema
// named "" at each DROP COLUMN. Columns are tracked per table name as
// written and compared bytewise, so a DROP COLUMN naming a column in another
// case, or a table written in another case than its CREATE TABLE, is not
// reported; the schema lookup itself follows schema.ObjectCaseSensitive. A
// CREATE or DROP TABLE does not clear what is known of a table. An ALTER or
// CREATE TABLE with no table is skipped, and so are statements nested in a
// routine body. Ported from Bytebase's V1 rule COLUMN_DISALLOW_DROP_IN_INDEX.
func ColumnDisallowDropInIndex(stmts []ast.Node, schema v1review.Schema) []v1review.Finding {
	indexed := make(map[string]map[string]bool)
	mark := func(table, column string) {
		if indexed[table] == nil {
			indexed[table] = make(map[string]bool)
		}
		indexed[table][column] = true
	}
	var findings []v1review.Finding
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, c := range n.Constraints {
				if c == nil || c.Type != ast.ConstrIndex {
					continue
				}
				for _, column := range c.Columns {
					mark(n.Table.Name, column)
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			table := n.Table.Name
			for _, cmd := range n.Commands {
				if cmd == nil || cmd.Type != ast.ATDropColumn {
					continue
				}
				for _, index := range schema.TableIndexesByName("", table, false) {
					for _, expression := range index.GetExpressions() {
						mark(table, expression)
					}
				}
				if indexed[table][cmd.Name] {
					start, end := span(cmd.Loc)
					findings = append(findings, v1review.Finding{
						Code:      v1review.ColumnInIndexDropped,
						Name:      cmd.Name,
						Table:     table,
						Start:     start,
						End:       end,
						Statement: i,
					})
				}
			}
		}
	}
	return findings
}
