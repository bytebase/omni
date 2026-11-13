package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDMLDryRunProbes returns, in order, the indexes of the statements
// the V1 rule STATEMENT_DML_DRY_RUN dry-runs: every INSERT, UPDATE, DELETE
// and MERGE, and every SELECT whose WITH holds one, judged through an
// EXPLAIN ANALYZE to the statement it runs. Bytebase EXPLAINs each one and
// passes the results to StatementDMLDryRun.
func StatementDMLDryRunProbes(stmts []ast.Node) []int {
	var probes []int
	for i, stmt := range stmts {
		switch n := liveTarget(stmt).(type) {
		case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
			probes = append(probes, i)
		case *ast.SelectStmt:
			if hasDataModifyingCTE(n.WithClause) {
				probes = append(probes, i)
			}
		}
	}
	return probes
}

// StatementDMLDryRun reports each probed statement whose dry run failed, in
// statement order: DryRunFailed, with Detail the probe's error and Statement
// its index. A statement with no probe result was not dry-run and is not
// reported. Findings anchor the statement (Start = End = 0). Ported from
// Bytebase's V1 rule STATEMENT_DML_DRY_RUN.
func StatementDMLDryRun(stmts []ast.Node, probes map[int]v1review.Probe) []v1review.Finding {
	var findings []v1review.Finding
	for _, i := range StatementDMLDryRunProbes(stmts) {
		probe, ok := probes[i]
		if !ok || probe.Error == "" {
			continue
		}
		findings = append(findings, v1review.Finding{Code: v1review.DryRunFailed, Detail: probe.Error, Statement: i})
	}
	return findings
}
