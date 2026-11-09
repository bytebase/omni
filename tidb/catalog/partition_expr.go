package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// allowedPartitionFuncs is the set of functions TiDB v8.5.5 permits in a
// RANGE/LIST/HASH partitioning expression (the non-COLUMNS form), read from the
// engine on 2026-11-09: a call to a function outside it answers 1564. EXTRACT is
// allowed at the function level; a week-based unit (EXTRACT(WEEK FROM ...))
// answers 1486 instead, a separate check this does not make. CEILING is allowed
// but its CEIL spelling is not, and WEEK is not allowed at all.
var allowedPartitionFuncs = map[string]bool{
	"abs":            true,
	"ceiling":        true,
	"datediff":       true,
	"day":            true,
	"dayofmonth":     true,
	"dayofweek":      true,
	"dayofyear":      true,
	"extract":        true,
	"floor":          true,
	"hour":           true,
	"microsecond":    true,
	"minute":         true,
	"mod":            true,
	"month":          true,
	"quarter":        true,
	"second":         true,
	"time_to_sec":    true,
	"to_days":        true,
	"to_seconds":     true,
	"unix_timestamp": true,
	"weekday":        true,
	"year":           true,
	"yearweek":       true,
}

// validatePartitionExpr refuses the partition schemes TiDB v8.5.5 refuses that
// the catalog would otherwise keep, measured against the engine on 2026-11-09:
//
//   - a RANGE/LIST/HASH expression (the non-COLUMNS form) that is not an
//     integer-valued expression the engine allows — a function outside
//     allowedPartitionFuncs (1564) or a bare column of a non-integer type
//     (1659);
//   - a bare MAXVALUE bound before the last RANGE partition (1481).
//
// COLUMNS partitioning accepts string and temporal columns, so it is left
// alone, and a form the catalog cannot judge (an unknown column, an expression
// that is neither a lone column nor a flagged function) is accepted, so no valid
// scheme is refused.
func validatePartitionExpr(tbl *Table, pc *nodes.PartitionClause) error {
	if pc == nil {
		return nil
	}
	if len(pc.Columns) == 0 {
		switch pc.Type {
		case nodes.PartitionRange, nodes.PartitionList, nodes.PartitionHash:
			if err := checkPartitionExprType(tbl, pc.Expr); err != nil {
				return err
			}
		}
	}
	if pc.Type == nodes.PartitionRange {
		if err := checkRangeMaxValuePosition(pc.Partitions); err != nil {
			return err
		}
	}
	return nil
}

// checkPartitionExprType refuses a non-COLUMNS RANGE/LIST/HASH expression that
// cannot yield an integer: a disallowed function or real division anywhere in
// the tree (1564), or a lone column of a non-integer type (1659).
func checkPartitionExprType(tbl *Table, expr nodes.ExprNode) error {
	if expr == nil {
		return nil
	}
	// A function call outside the allow-list, or a real-division operator, at
	// any depth in the expression tree answers 1564. Real division (`id / 2`)
	// is the only arithmetic operator the engine rejects; the rest
	// (`+ - * DIV MOD`) stay, so only it is flagged and no valid operator is
	// refused.
	var bad error
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if bad != nil {
			return false
		}
		switch e := n.(type) {
		case *nodes.FuncCallExpr:
			if !allowedPartitionFuncs[strings.ToLower(e.Name)] {
				bad = errPartitionFuncNotAllowed()
				return false
			}
		case *nodes.BinaryExpr:
			if e.Op == nodes.BinOpDiv {
				bad = errPartitionFuncNotAllowed()
				return false
			}
		}
		return true
	})
	if bad != nil {
		return bad
	}
	// A lone column reference carries the column's type straight into the
	// partitioning function, so a non-integer column answers 1659.
	if ref := bareColumnRef(expr); ref != nil {
		if col := tbl.GetColumn(ref.Column); col != nil && !isIntegerType(strings.ToLower(col.DataType)) {
			return &Error{
				Code:     1659,
				SQLState: "HY000",
				Message:  "Field '" + ref.Column + "' is of a not allowed type for this type of partitioning",
			}
		}
	}
	return nil
}

// bareColumnRef returns the column reference an expression is, through any
// parentheses, or nil when it is anything more than a lone column.
func bareColumnRef(expr nodes.ExprNode) *nodes.ColumnRef {
	for {
		switch e := expr.(type) {
		case *nodes.ParenExpr:
			expr = e.Expr
		case *nodes.ColumnRef:
			return e
		default:
			return nil
		}
	}
}

func errPartitionFuncNotAllowed() error {
	return &Error{Code: 1564, SQLState: "HY000", Message: "This partition function is not allowed"}
}

// checkRangeMaxValuePosition answers 1481 when a RANGE partition before the last
// has VALUES LESS THAN MAXVALUE: MAXVALUE is the open top bound, so it belongs to
// the last partition only. A multi-column RANGE COLUMNS tuple carries MAXVALUE
// inside a value list, not as a lone bound, so it is not matched here (its
// ordering is a separate 1493 the catalog leaves to the engine).
func checkRangeMaxValuePosition(defs []*nodes.PartitionDef) error {
	for i, pd := range defs {
		if i != len(defs)-1 && isBareMaxValue(pd.Values) {
			return &Error{
				Code:     1481,
				SQLState: "HY000",
				Message:  "MAXVALUE can only be used in last partition definition",
			}
		}
	}
	return nil
}

func isBareMaxValue(v nodes.Node) bool {
	s, ok := v.(*nodes.String)
	return ok && s.Str == "MAXVALUE"
}
