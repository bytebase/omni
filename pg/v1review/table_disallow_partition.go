package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowPartition judges top-level statements that partition a
// table: a CREATE TABLE with a PARTITION BY clause, and an ALTER statement
// with an ATTACH PARTITION command, are TablePartitioned, once per
// statement. ALTER INDEX ... ATTACH PARTITION counts, as it did in V1; a
// CREATE TABLE ... PARTITION OF without its own PARTITION BY, DETACH
// PARTITION, and statements under EXPLAIN or in a routine body are not
// judged. Each finding anchors the whole statement. A statement may be
// passed bare or wrapped in an *ast.RawStmt. Ported from Bytebase's V1 rule
// TABLE_DISALLOW_PARTITION.
func TableDisallowPartition(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if n.Partspec == nil {
				continue
			}
			loc = n.Loc
		case *ast.AlterTableStmt:
			if !alterTableHas(n, func(cmd *ast.AlterTableCmd) bool {
				return ast.AlterTableType(cmd.Subtype) == ast.AT_AttachPartition
			}) {
				continue
			}
			loc = n.Loc
		default:
			continue
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: v1review.TablePartitioned, Start: start, End: end})
	}
	return findings
}

// alterTableHas reports whether any of an ALTER statement's commands
// satisfies match.
func alterTableHas(n *ast.AlterTableStmt, match func(*ast.AlterTableCmd) bool) bool {
	if n.Cmds == nil {
		return false
	}
	for _, item := range n.Cmds.Items {
		if cmd, ok := item.(*ast.AlterTableCmd); ok && match(cmd) {
			return true
		}
	}
	return false
}
