package catalog

import (
	"fmt"
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
//   - a column the partitioning function names that the table does not have:
//     1054 "Unknown column '<name>' in 'partition function'" for the
//     expression form, 1488 "Field in list of fields for partition function
//     not found in table" for a COLUMNS or KEY column list (measured on
//     v8.5.5, 2026-11-19);
//   - a bare MAXVALUE bound before the last RANGE partition (1481).
//
// COLUMNS partitioning accepts string and temporal columns, so its types are
// left alone, and a form the catalog cannot judge (an expression that is
// neither a lone column nor a flagged function) is accepted, so no valid scheme
// is refused.
func validatePartitionExpr(tbl *Table, pc *nodes.PartitionClause) error {
	if pc == nil {
		return nil
	}
	if err := checkPartitionColumnsExist(tbl, pc); err != nil {
		return err
	}
	if len(pc.Columns) == 0 {
		switch pc.Type {
		case nodes.PartitionRange, nodes.PartitionList, nodes.PartitionHash:
			if err := checkPartitionExprType(tbl, pc.Expr); err != nil {
				return err
			}
		}
	}
	// KEY reads any column but a BLOB, TEXT or JSON one, which v8.5.5 answers
	// with 1659 (CREATE TABLE h9 (a TEXT) PARTITION BY KEY (a)).
	if pc.Type == nodes.PartitionKey {
		for _, name := range pc.Columns {
			col := tbl.GetColumn(name)
			if col == nil {
				continue
			}
			switch strings.ToLower(col.DataType) {
			case "tinyblob", "blob", "mediumblob", "longblob", "tinytext", "text", "mediumtext", "longtext", "json":
				return &Error{Code: 1659, SQLState: "HY000",
					Message: "Field '" + col.Name + "' is of a not allowed type for this type of partitioning"}
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

// checkPartitionColumnsExist resolves the columns a partition clause names
// against the table: an expression's column references answer 1054, a COLUMNS
// or KEY column list answers 1488 (both read on v8.5.5, 2026-11-19).
func checkPartitionColumnsExist(tbl *Table, pc *nodes.PartitionClause) error {
	if pc == nil {
		return nil
	}
	if len(pc.Columns) > 0 {
		for _, name := range pc.Columns {
			if tbl.GetColumn(name) == nil {
				return &Error{Code: ErrPartitionFieldNotFound, SQLState: sqlState(ErrPartitionFieldNotFound),
					Message: "Field in list of fields for partition function not found in table"}
			}
		}
		return nil
	}
	if pc.Expr == nil {
		return nil
	}
	var missing string
	nodes.Inspect(pc.Expr, func(n nodes.Node) bool {
		if missing != "" {
			return false
		}
		if ref, ok := n.(*nodes.ColumnRef); ok && ref.Column != "" {
			if tbl.GetColumn(ref.Column) == nil {
				missing = ref.Column
				return false
			}
		}
		return true
	})
	if missing != "" {
		return &Error{Code: ErrNoSuchColumn, SQLState: sqlState(ErrNoSuchColumn),
			Message: fmt.Sprintf("Unknown column '%s' in 'partition function'", missing)}
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
	// partitioning function, so a column of a type the engine does not allow
	// there answers 1659.
	if ref := bareColumnRef(expr); ref != nil {
		if col := tbl.GetColumn(ref.Column); col != nil && !isPartitionColumnType(strings.ToLower(col.DataType)) {
			return &Error{
				Code:     1659,
				SQLState: "HY000",
				Message:  "Field '" + ref.Column + "' is of a not allowed type for this type of partitioning",
			}
		}
	}
	return nil
}

// isPartitionColumnType reports whether a column of this type may stand alone in
// a non-COLUMNS partitioning function. TiDB v8.5.5 allows the integer types plus
// BIT, YEAR and BOOL (a tinyint(1)), and rejects every other type (DECIMAL,
// DOUBLE, ENUM, SET, the string types and the bare temporal types) with 1659,
// read from the engine on 2026-11-09. It is wider than isIntegerType, the
// storage-type helper, which BIT and YEAR are not.
func isPartitionColumnType(dt string) bool {
	switch dt {
	case "bit", "year":
		return true
	}
	return isIntegerType(dt)
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
