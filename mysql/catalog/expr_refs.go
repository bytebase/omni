package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// exprRef is one expression a statement defines on a table: a generated
// column's (column set) or a CHECK constraint's (check set, the constraint's
// name).
type exprRef struct {
	column string
	check  string
	expr   nodes.ExprNode
}

// checkExprReferences holds the expressions a statement defined to the table
// it leaves, which is where MySQL reads their columns: ADD COLUMN z INT, ADD
// CHECK (z > 0) runs, and DROP COLUMN age, ADD CHECK (age > 1) is 1054. A
// generated column reads columns the table has (1054), none of them
// AUTO_INCREMENT (3109) and no generated column defined after it (3107); a
// CHECK constraint reads columns the table has (1054), none AUTO_INCREMENT
// (3818) (mysql:8.0.46).
func checkExprReferences(tbl *Table, refs []exprRef) error {
	for _, ref := range refs {
		for _, name := range exprColumnNames(ref.expr) {
			col := tbl.GetColumn(name)
			switch {
			case col == nil && ref.check != "":
				return &Error{Code: 1054, SQLState: "42S22", Message: fmt.Sprintf("Unknown column '%s' in 'check constraint %s expression'", name, ref.check)}
			case col == nil:
				return &Error{Code: 1054, SQLState: "42S22", Message: fmt.Sprintf("Unknown column '%s' in 'generated column function'", name)}
			case col.AutoIncrement && ref.check != "":
				return &Error{Code: 3818, SQLState: "HY000", Message: fmt.Sprintf("Check constraint '%s' cannot refer to an auto-increment column.", ref.check)}
			case col.AutoIncrement:
				return &Error{Code: 3109, SQLState: "HY000", Message: fmt.Sprintf("Generated column '%s' cannot refer to auto-increment column.", ref.column)}
			case ref.column != "" && col.Generated != nil && tbl.colByName[toLower(col.Name)] >= tbl.colByName[toLower(ref.column)]:
				return &Error{Code: 3107, SQLState: "HY000", Message: "Generated column can refer only to generated columns defined prior to it."}
			}
		}
	}
	return nil
}

// columnExprRefs returns the generated expressions and column CHECKs of the
// columns defs define, with checks, the CHECK constraints the table gained
// from them in order, whose names the column CHECKs take.
func columnExprRefs(defs []*nodes.ColumnDef) (generated []exprRef, checks []nodes.ExprNode) {
	for _, def := range defs {
		if def == nil {
			continue
		}
		if def.Generated != nil && def.Generated.Expr != nil {
			generated = append(generated, exprRef{column: def.Name, expr: def.Generated.Expr})
		}
		for _, cc := range def.Constraints {
			if cc != nil && cc.Type == nodes.ColConstrCheck && cc.Expr != nil {
				checks = append(checks, cc.Expr)
			}
		}
	}
	return generated, checks
}

// namedChecks pairs CHECK expressions with the constraints they made, in
// order; with a count that does not match, it pairs none.
func namedChecks(exprs []nodes.ExprNode, made []*Constraint) []exprRef {
	if len(exprs) != len(made) {
		return nil
	}
	refs := make([]exprRef, len(exprs))
	for i, e := range exprs {
		refs[i] = exprRef{check: made[i].Name, expr: e}
	}
	return refs
}
