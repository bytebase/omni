package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementRequireAlgorithmOption reports each top-level ALTER TABLE with no
// ALGORITHM clause; statements in a routine body are not judged. Findings
// anchor nothing (the V1 advisor took the statement's first content line).
// OptionMissing is the only code; Name is "ALGORITHM". Ported from
// Bytebase's V1 rule STATEMENT_REQUIRE_ALGORITHM_OPTION.
func StatementRequireAlgorithmOption(stmts []ast.Node) []v1review.Finding {
	return alterTablesWithout(stmts, ast.ATAlgorithm, "ALGORITHM")
}

// StatementRequireLockOption reports each top-level ALTER TABLE with no LOCK
// clause; statements in a routine body are not judged. Findings anchor
// nothing (the V1 advisor took the statement's first content line).
// OptionMissing is the only code; Name is "LOCK". Ported from Bytebase's V1
// rule STATEMENT_REQUIRE_LOCK_OPTION.
func StatementRequireLockOption(stmts []ast.Node) []v1review.Finding {
	return alterTablesWithout(stmts, ast.ATLock, "LOCK")
}

func alterTablesWithout(stmts []ast.Node, option ast.AlterTableCmdType, name string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		alter, ok := stmt.(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		hasOption := false
		for _, cmd := range alter.Commands {
			if cmd.Type == option {
				hasOption = true
				break
			}
		}
		if !hasOption {
			findings = append(findings, v1review.Finding{Code: v1review.OptionMissing, Name: name})
		}
	}
	return findings
}
