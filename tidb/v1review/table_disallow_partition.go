package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowPartition judges top-level CREATE TABLE and ALTER TABLE
// statements for partitioning. A CREATE TABLE with a PARTITION BY clause is
// TablePartitioned, and so is an ALTER TABLE with a PARTITION BY command;
// partition maintenance (ADD, DROP, REORGANIZE PARTITION and the like) and
// statements nested in a routine body or carried by EXPLAIN are not judged.
// A statement has at most one finding. Table is the table name as parsed,
// without its database, empty when the statement names none; each finding
// anchors the whole statement. Ported from Bytebase's V1 rule
// TABLE_DISALLOW_PARTITION.
func TableDisallowPartition(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var partitioned bool
		var table *ast.TableRef
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			partitioned, table, loc = n.Partitions != nil, n.Table, n.Loc
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd != nil && cmd.Type == ast.ATPartitionBy && cmd.PartitionBy != nil {
					partitioned = true
					break
				}
			}
			table, loc = n.Table, n.Loc
		}
		if !partitioned {
			continue
		}
		name := ""
		if table != nil {
			name = table.Name
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: v1review.TablePartitioned, Table: name, Start: start, End: end})
	}
	return findings
}
