package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDMLDryRunProbes returns, in order, the indexes of the top-level
// statements the V1 rule STATEMENT_DML_DRY_RUN dry-runs: every INSERT,
// UPDATE and DELETE. Statements nested in a routine body are not probed.
// Bytebase EXPLAINs each one and passes the results to StatementDMLDryRun.
func StatementDMLDryRunProbes(stmts []ast.Node) []int {
	var probes []int
	for i, stmt := range stmts {
		switch stmt.(type) {
		case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
			probes = append(probes, i)
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
