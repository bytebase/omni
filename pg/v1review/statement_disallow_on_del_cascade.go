package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowOnDelCascade reports a top-level CREATE TABLE or ALTER
// TABLE that declares a foreign key with ON DELETE CASCADE, at most one
// OnDeleteCascade finding per statement. In CREATE TABLE both table and
// column constraints count; in ALTER TABLE only ADD CONSTRAINT does, not a
// foreign key on an added column. Findings anchor nothing: Bytebase's V1
// advisor placed them on the statement's first line. Ported from
// Bytebase's V1 rule STATEMENT_DISALLOW_ON_DEL_CASCADE.
func StatementDisallowOnDelCascade(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var cascade bool
		switch n := unwrapRaw(stmt).(type) {
		case *ast.CreateStmt:
			cascade = createDeletesCascade(n)
		case *ast.AlterTableStmt:
			cascade = alterDeletesCascade(n)
		}
		if cascade {
			findings = append(findings, v1review.Finding{Code: v1review.OnDeleteCascade})
		}
	}
	return findings
}

func deletesCascade(c *ast.Constraint) bool {
	return c.Contype == ast.CONSTR_FOREIGN && c.FkDelaction == 'c'
}

func createDeletesCascade(create *ast.CreateStmt) bool {
	if create.TableElts == nil {
		return false
	}
	for _, item := range create.TableElts.Items {
		switch n := item.(type) {
		case *ast.Constraint:
			if deletesCascade(n) {
				return true
			}
		case *ast.ColumnDef:
			for _, c := range columnConstraints(n) {
				if deletesCascade(c) {
					return true
				}
			}
		}
	}
	return false
}

func alterDeletesCascade(alter *ast.AlterTableStmt) bool {
	for _, cmd := range alterTableCmds(alter) {
		if ast.AlterTableType(cmd.Subtype) != ast.AT_AddConstraint {
			continue
		}
		if c, ok := cmd.Def.(*ast.Constraint); ok && deletesCascade(c) {
			return true
		}
	}
	return false
}
