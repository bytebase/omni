package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementRequireAlgorithmOption reports each top-level ALTER TABLE with
// no ALGORITHM option. Name is "ALGORITHM". Findings anchor nothing (the V1
// advisor placed them on the statement's first line). OptionMissing is the
// only code. Ported from Bytebase's V1 rule
// STATEMENT_REQUIRE_ALGORITHM_OPTION.
func StatementRequireAlgorithmOption(stmts []ast.Node) []v1review.Finding {
	return requireAlterTableOption(stmts, ast.ATAlgorithm, "ALGORITHM")
}

// StatementRequireLockOption reports each top-level ALTER TABLE with no
// LOCK option. Name is "LOCK". Findings anchor nothing (the V1 advisor
// placed them on the statement's first line). OptionMissing is the only
// code. Ported from Bytebase's V1 rule STATEMENT_REQUIRE_LOCK_OPTION.
func StatementRequireLockOption(stmts []ast.Node) []v1review.Finding {
	return requireAlterTableOption(stmts, ast.ATLock, "LOCK")
}

func requireAlterTableOption(stmts []ast.Node, option ast.AlterTableCmdType, name string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		alter, ok := stmt.(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		has := false
		for _, cmd := range alter.Commands {
			if cmd.Type == option {
				has = true
				break
			}
		}
		if !has {
			findings = append(findings, v1review.Finding{Code: v1review.OptionMissing, Name: name})
		}
	}
	return findings
}
