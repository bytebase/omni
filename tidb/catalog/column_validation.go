package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// This file holds CREATE TABLE column-type/default validations that TiDB v8.5.5
// enforces and the catalog used to over-accept. Every rule was measured against
// a live TiDB v8.5.5 container (2026-11-09); each rejects only what TiDB
// rejects, so no valid DDL is turned away.

// errTooBigPrecision is TiDB 1426 for a DECIMAL whose precision exceeds 65.
func errTooBigPrecision(precision int, col string) error {
	return &Error{Code: ErrTooBigPrecision, SQLState: sqlState(ErrTooBigPrecision),
		Message: fmt.Sprintf("Too-big precision %d specified for '%s'. Maximum is 65.", precision, col)}
}

// errTooBigScale is TiDB 1425 for a DECIMAL whose scale exceeds 30.
func errTooBigScale(scale int, col string) error {
	return &Error{Code: ErrTooBigScale, SQLState: sqlState(ErrTooBigScale),
		Message: fmt.Sprintf("Too big scale %d specified for column '%s'. Maximum is 30.", scale, col)}
}

// errMBiggerThanD is TiDB 1427 for a DECIMAL whose scale exceeds its precision.
func errMBiggerThanD(col string) error {
	return &Error{Code: ErrMBiggerThanD, SQLState: sqlState(ErrMBiggerThanD),
		Message: fmt.Sprintf("For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '%s').", col)}
}

// errPrimaryKeyCantHaveNull is TiDB 1171 for a PRIMARY KEY part declared NULL.
func errPrimaryKeyCantHaveNull() error {
	return &Error{Code: ErrPrimaryCantHaveNull, SQLState: sqlState(ErrPrimaryCantHaveNull),
		Message: "All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead"}
}

// errDefaultFunctionNotAllowed is TiDB 3770 for a DEFAULT expression whose
// top-level function is not on the allow list.
func errDefaultFunctionNotAllowed(col, fn string) error {
	return &Error{Code: ErrDefValFunctionNotAllow, SQLState: sqlState(ErrDefValFunctionNotAllow),
		Message: fmt.Sprintf("Default value expression of column '%s' contains a disallowed function: `%s`.", col, strings.ToUpper(fn))}
}

// checkDecimalBounds enforces TiDB v8.5.5's DECIMAL/NUMERIC precision and scale
// limits: precision above 65 is 1426, scale above 30 is 1425, and scale above
// precision is 1427. DECIMAL(65,30), DECIMAL(30,30) and DECIMAL(0,0) are all
// accepted. FLOAT/DOUBLE carry a different precision ceiling (TiDB accepts
// FLOAT(66,30)) and are deliberately left to their own path.
func checkDecimalBounds(colName string, dt *nodes.DataType) error {
	if dt == nil {
		return nil
	}
	switch strings.ToLower(dt.Name) {
	case "decimal", "numeric", "dec", "fixed":
	default:
		return nil
	}
	if dt.Length > 65 {
		return errTooBigPrecision(dt.Length, colName)
	}
	if dt.Scale > 30 {
		return errTooBigScale(dt.Scale, colName)
	}
	if dt.Scale > dt.Length {
		return errMBiggerThanD(colName)
	}
	return nil
}

// columnDefExplicitlyNull reports whether a column definition carries an
// explicit NULL nullability. TiDB v8.5.5 rejects a PRIMARY KEY column declared
// explicitly NULL with 1171, while a column left implicitly nullable is silently
// promoted to NOT NULL. A later NOT NULL does not clear the rejection: TiDB
// rejects `a INT NULL NOT NULL PRIMARY KEY` and `a INT NOT NULL NULL PRIMARY
// KEY` alike (measured on v8.5.5), so any explicit NULL token is enough.
func columnDefExplicitlyNull(colDef *nodes.ColumnDef) bool {
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrNull {
			return true
		}
	}
	return false
}

// findColumnDef returns the column definition named name (case-insensitively),
// or nil when the statement declares no such column.
func findColumnDef(cols []*nodes.ColumnDef, name string) *nodes.ColumnDef {
	for _, cd := range cols {
		if strings.EqualFold(cd.Name, name) {
			return cd
		}
	}
	return nil
}

// allowedDefaultExprFuncs is the exact set of functions TiDB v8.5.5 permits as
// the top-level call of a parenthesized DEFAULT expression (pkg/ddl/add_column.go
// getFuncCallDefaultValue), plus the temporal synonyms its parser folds into
// CURRENT_TIMESTAMP / CURRENT_DATE. Anything else is 3770. The arg-conditional
// functions (date_format, replace, upper, str_to_date) are allowed outright:
// TiDB rejects only their disallowed-argument forms, so accepting them here can
// only over-accept, never turn away a form TiDB takes.
var allowedDefaultExprFuncs = map[string]bool{
	"current_timestamp": true,
	"now":               true,
	"localtime":         true,
	"localtimestamp":    true,
	"current_date":      true,
	"curdate":           true,
	"nextval":           true,
	"rand":              true,
	"uuid":              true,
	"uuid_to_bin":       true,
	"date_format":       true,
	"replace":           true,
	"upper":             true,
	"str_to_date":       true,
	"json_object":       true,
	"json_array":        true,
	"json_quote":        true,
	"vec_from_text":     true,
}

// checkDefaultExprFunction rejects a parenthesized DEFAULT whose top-level
// function is not on TiDB v8.5.5's allow list (3770). Only a parenthesized
// expression default is inspected: a bare DEFAULT NOW()/CURRENT_TIMESTAMP is a
// temporal default the catalog already handles, and a non-function expression
// default is left to the parser, which rejects column references itself.
func checkDefaultExprFunction(colName string, expr nodes.ExprNode) error {
	paren, ok := expr.(*nodes.ParenExpr)
	if !ok {
		return nil
	}
	inner := paren.Expr
	for {
		if p, ok := inner.(*nodes.ParenExpr); ok {
			inner = p.Expr
			continue
		}
		break
	}
	fn, ok := inner.(*nodes.FuncCallExpr)
	if !ok {
		return nil
	}
	if allowedDefaultExprFuncs[strings.ToLower(fn.Name)] {
		return nil
	}
	return errDefaultFunctionNotAllowed(colName, fn.Name)
}

// columnDefaultExpr returns the effective DEFAULT expression of a column
// definition: the column-constraint DEFAULT when present (it overrides), else
// the top-level one. nil when the column has no default.
func columnDefaultExpr(colDef *nodes.ColumnDef) nodes.ExprNode {
	expr := colDef.DefaultValue
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault && cc.Expr != nil {
			expr = cc.Expr
		}
	}
	return expr
}
