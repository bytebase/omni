package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// checkTableCheckExpr refuses a table CHECK that names a column tbl lacks (3820),
// then one whose expression is not a condition (3812), before its name is compared
// with the table's other checks (3822): OceanBase answers ADD CONSTRAINT k CHECK (c),
// with k taken, with 3812 (TestOracle_OceanBaseCheckExpr).
func checkTableCheckExpr(tbl *Table, name string, expr nodes.ExprNode) error {
	for _, col := range exprColumnNames(expr) {
		if tbl.GetColumn(col) == nil {
			return &Error{Code: 3820, SQLState: "HY000",
				Message: fmt.Sprintf("Check constraint '%s' refers to non-existing column '%s'.", name, col)}
		}
	}
	if notCondition(expr) {
		return &Error{Code: 3812, SQLState: "HY000",
			Message: fmt.Sprintf("An expression of non-boolean type specified to a check constraint '%s'", name)}
	}
	return nil
}

// notCondition reports an expression OceanBase does not take as a CHECK's
// condition: a column, a string, or arithmetic. CHECK (1) is taken, and so is any
// form not listed here, whose type this does not work out.
func notCondition(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return notCondition(e.Expr)
	case *nodes.ColumnRef, *nodes.StringLit:
		return true
	case *nodes.BinaryExpr:
		switch e.Op {
		case nodes.BinOpAdd, nodes.BinOpSub, nodes.BinOpMul, nodes.BinOpDiv, nodes.BinOpMod, nodes.BinOpDivInt,
			nodes.BinOpBitAnd, nodes.BinOpBitOr, nodes.BinOpBitXor, nodes.BinOpShiftLeft, nodes.BinOpShiftRight:
			return true
		}
	}
	return false
}
