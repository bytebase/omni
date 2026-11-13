package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowPartition judges every top-level CREATE TABLE and ALTER
// TABLE; statements nested in a routine body are not judged. A CREATE
// TABLE with a PARTITION BY clause is TablePartitioned, as is an ALTER
// TABLE with a PARTITION BY action, once however many it has. Other
// partition actions (ADD, DROP, REORGANIZE PARTITION, REMOVE PARTITIONING
// and the like) are not judged. Each finding anchors the whole statement
// and names nothing. Ported from Bytebase's V1 rule
// TABLE_DISALLOW_PARTITION.
func TableDisallowPartition(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Partitions != nil {
				start, end := span(n.Loc)
				findings = append(findings, v1review.Finding{Code: v1review.TablePartitioned, Start: start, End: end})
			}
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd != nil && cmd.Type == ast.ATPartitionBy && cmd.PartitionBy != nil {
					start, end := span(n.Loc)
					findings = append(findings, v1review.Finding{Code: v1review.TablePartitioned, Start: start, End: end})
					break
				}
			}
		default:
		}
	}
	return findings
}
