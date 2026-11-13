package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementAddCheckNotValid reports each CHECK constraint a top-level ALTER
// TABLE adds without NOT VALID, one CheckValidated finding per ADD
// CONSTRAINT. A CHECK in CREATE TABLE or on an added column is not judged.
// Findings anchor nothing: Bytebase's V1 advisor placed them on the
// statement's first line. Ported from Bytebase's V1 rule
// STATEMENT_ADD_CHECK_NOT_VALID.
func StatementAddCheckNotValid(stmts []ast.Node) []v1review.Finding {
	return addedConstraintsValidated(stmts, ast.CONSTR_CHECK, v1review.CheckValidated)
}

// addedConstraintsValidated reports each constraint of type contype that a
// top-level ALTER TABLE ... ADD CONSTRAINT adds without NOT VALID.
func addedConstraintsValidated(stmts []ast.Node, contype ast.ConstrType, code v1review.Code) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		alter, ok := unwrapRaw(stmt).(*ast.AlterTableStmt)
		if !ok {
			continue
		}
		for _, cmd := range alterTableCmds(alter) {
			if ast.AlterTableType(cmd.Subtype) != ast.AT_AddConstraint {
				continue
			}
			constraint, ok := cmd.Def.(*ast.Constraint)
			if !ok || constraint.Contype != contype {
				continue
			}
			if !constraint.SkipValidation {
				findings = append(findings, v1review.Finding{Code: code})
			}
		}
	}
	return findings
}
