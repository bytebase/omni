package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementAddForeignKeyNotValid reports each foreign key a top-level ALTER
// TABLE adds as a table constraint without NOT VALID, one
// ForeignKeyValidated finding per ADD CONSTRAINT. A foreign key in CREATE
// TABLE or on an added column is not judged. Findings anchor nothing:
// Bytebase's V1 advisor placed them on the statement's first line. Ported
// from Bytebase's V1 rule STATEMENT_ADD_FOREIGN_KEY_NOT_VALID.
func StatementAddForeignKeyNotValid(stmts []ast.Node) []v1review.Finding {
	return addedConstraintsValidated(stmts, ast.CONSTR_FOREIGN, v1review.ForeignKeyValidated)
}
