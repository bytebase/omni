package v1review

import (
	"strconv"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableLimitSize reports every table a top-level ALTER TABLE, TRUNCATE or
// DROP TABLE names, in statement order and once per mention, whose row count
// in snapshot is maxRows or more, as TableTooLarge; statements nested in a
// routine body are not judged. Name is the table name as written, without
// its database, and Value its row count in decimal. The row count is that of
// the first table of exactly that name in the snapshot's first schema, and 0
// when there is none, so with maxRows 0 or less a table the snapshot lacks is
// reported too. A nil snapshot or one with no schema reports nothing.
//
// Every finding names the last of stmts as its Statement, with no range: the
// V1 advisor placed all of them on the first non-blank line of the last
// statement it was given, whatever statement named the table, and the caller
// keeps that computation. Ported from Bytebase's V1 rule TABLE_LIMIT_SIZE.
func TableLimitSize(stmts []ast.Node, snapshot *metadata.DatabaseSchemaMetadata, maxRows int) []v1review.Finding {
	var names []string
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.AlterTableStmt:
			if n.Table != nil {
				names = append(names, n.Table.Name)
			}
		case *ast.TruncateStmt:
			for _, ref := range n.Tables {
				if ref != nil {
					names = append(names, ref.Name)
				}
			}
		case *ast.DropTableStmt:
			for _, ref := range n.Tables {
				if ref != nil {
					names = append(names, ref.Name)
				}
			}
		default:
		}
	}
	if snapshot == nil || len(snapshot.Schemas) == 0 || len(stmts) == 0 {
		return nil
	}
	last := len(stmts) - 1
	var findings []v1review.Finding
	for _, name := range names {
		var rows int64
		for _, table := range snapshot.Schemas[0].Tables {
			if table.Name == name {
				rows = table.RowCount
				break
			}
		}
		if rows >= int64(maxRows) {
			findings = append(findings, v1review.Finding{Code: v1review.TableTooLarge, Name: name, Value: strconv.FormatInt(rows, 10), Statement: last})
		}
	}
	return findings
}
