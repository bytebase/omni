package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowTruncate reports what top-level statements truncate: each
// table a TRUNCATE TABLE lists, and each partition an ALTER TABLE ...
// TRUNCATE PARTITION names, or "ALL PARTITIONS" for TRUNCATE PARTITION ALL.
// Statements nested in a routine body are not judged. Table is the table as
// written, qualified as db.table when the statement qualified it. A
// TRUNCATE TABLE finding anchors the statement; a partition finding anchors
// its ALTER TABLE command. Codes: Truncate, TruncatePartition (Name is the
// partition). Ported from Bytebase's V1 rule STATEMENT_DISALLOW_TRUNCATE.
func StatementDisallowTruncate(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.TruncateStmt:
			start, end := anchor(n.Loc)
			for _, table := range n.Tables {
				findings = append(findings, v1review.Finding{Code: v1review.Truncate, Table: qualifiedTableName(table), Start: start, End: end})
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			table := qualifiedTableName(n.Table)
			for _, cmd := range n.Commands {
				if cmd.Type != ast.ATTruncatePartition {
					continue
				}
				start, end := anchor(cmd.Loc)
				if cmd.AllPartitions {
					findings = append(findings, v1review.Finding{Code: v1review.TruncatePartition, Name: "ALL PARTITIONS", Table: table, Start: start, End: end})
					continue
				}
				for _, partition := range cmd.PartitionNames {
					findings = append(findings, v1review.Finding{Code: v1review.TruncatePartition, Name: partition, Table: table, Start: start, End: end})
				}
			}
		}
	}
	return findings
}

// qualifiedTableName returns the table as written: db.table when qualified.
func qualifiedTableName(table *ast.TableRef) string {
	if table.Schema != "" {
		return table.Schema + "." + table.Name
	}
	return table.Name
}
