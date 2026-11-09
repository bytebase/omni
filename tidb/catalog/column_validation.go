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

// errTooBigFspPrecision is TiDB 1426 for a DATETIME/TIMESTAMP/TIME whose
// fractional-seconds precision exceeds 6.
func errTooBigFspPrecision(precision int, col string) error {
	return &Error{Code: ErrTooBigPrecision, SQLState: sqlState(ErrTooBigPrecision),
		Message: fmt.Sprintf("Too-big precision %d specified for '%s'. Maximum is 6.", precision, col)}
}

// errInvalidDefault is TiDB 1067 for an invalid column default, including a
// NOW()-family default whose precision does not match the column's.
func errInvalidDefault(col string) error {
	return &Error{Code: ErrInvalidDefault, SQLState: sqlState(ErrInvalidDefault),
		Message: fmt.Sprintf("Invalid default value for '%s'", col)}
}

// checkTemporalFspBounds enforces TiDB v8.5.5's 1426 ceiling on the fractional-
// seconds precision of DATETIME, TIMESTAMP and TIME: a precision above 6 is
// rejected, 6 and below are accepted. Measured against a live container.
func checkTemporalFspBounds(colName string, dt *nodes.DataType) error {
	if dt == nil {
		return nil
	}
	switch strings.ToLower(dt.Name) {
	case "datetime", "timestamp", "time":
		if dt.Length > 6 {
			return errTooBigFspPrecision(dt.Length, colName)
		}
	}
	return nil
}

// nowFamilyDefaultFuncs are the temporal DEFAULT functions whose fractional-
// seconds precision TiDB v8.5.5 holds to the column's own precision.
var nowFamilyDefaultFuncs = map[string]bool{
	"now":               true,
	"current_timestamp": true,
	"localtime":         true,
	"localtimestamp":    true,
}

// checkTemporalDefaultFsp enforces TiDB v8.5.5's 1067 for a DATETIME or
// TIMESTAMP column whose NOW()-family DEFAULT names a precision other than the
// column's own: DATETIME(6) DEFAULT NOW() (0 vs 6) and DATETIME(3) DEFAULT
// NOW(6) both fail, while DATETIME DEFAULT NOW() and DATETIME(6) DEFAULT NOW(6)
// match and pass. A non-NOW default and a non-temporal column are left alone.
func checkTemporalDefaultFsp(colName string, dt *nodes.DataType, defaultExpr nodes.ExprNode) error {
	if dt == nil {
		return nil
	}
	switch strings.ToLower(dt.Name) {
	case "datetime", "timestamp":
	default:
		return nil
	}
	fn, ok := defaultExpr.(*nodes.FuncCallExpr)
	if !ok || !nowFamilyDefaultFuncs[strings.ToLower(fn.Name)] {
		return nil
	}
	funcFsp := 0
	if len(fn.Args) > 0 {
		if lit, ok := fn.Args[0].(*nodes.IntLit); ok {
			funcFsp = int(lit.Value)
		}
	}
	if funcFsp != dt.Length {
		return errInvalidDefault(colName)
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

// errFKColumnCannotBeNull is TiDB 1830 for a NOT NULL referencing column under
// ON DELETE/UPDATE SET NULL.
func errFKColumnCannotBeNull(col, constraint string) error {
	return &Error{Code: ErrFKColumnCannotBeNull, SQLState: sqlState(ErrFKColumnCannotBeNull),
		Message: fmt.Sprintf("Column '%s' cannot be NOT NULL: needed in a foreign key constraint '%s' SET NULL", col, constraint)}
}

// errFKCannotUseVirtualColumn is TiDB 3733 for a foreign key whose referencing
// column is a VIRTUAL generated column. A STORED generated column is allowed.
func errFKCannotUseVirtualColumn(constraint, col string) error {
	return &Error{Code: ErrFKCannotUseVirtualColumn, SQLState: sqlState(ErrFKCannotUseVirtualColumn),
		Message: fmt.Sprintf("Foreign key '%s' uses virtual column '%s' which is not supported.", constraint, col)}
}

// errVirtualGeneratedColumnPrimaryKey is TiDB 3106 for a VIRTUAL generated
// column named as a primary key (inline or in a table-level PRIMARY KEY). A
// STORED generated column may be a primary key.
func errVirtualGeneratedColumnPrimaryKey() error {
	return &Error{Code: ErrUnsupportedGeneratedStorageChange, SQLState: sqlState(ErrUnsupportedGeneratedStorageChange),
		Message: "'Defining a virtual generated column as primary key' is not supported for generated columns."}
}

// errGeneratedColumnFunction is TiDB 3102 for a generated-column expression that
// contains a disallowed (non-deterministic) function.
func errGeneratedColumnFunction(col string) error {
	return &Error{Code: ErrGeneratedColumnFunction, SQLState: sqlState(ErrGeneratedColumnFunction),
		Message: fmt.Sprintf("Expression of generated column '%s' contains a disallowed function.", col)}
}

// errCantRemoveAllFields is TiDB 1090 for dropping a table's only column.
func errCantRemoveAllFields(col, table string) error {
	return &Error{Code: ErrCantRemoveAllFields, SQLState: sqlState(ErrCantRemoveAllFields),
		Message: fmt.Sprintf("can't drop only column %s in table %s", col, table)}
}

// checkIndexName rejects a user-supplied secondary-index name of "PRIMARY"
// (case-insensitive), which TiDB v8.5.5 answers 1280 for — that name is
// reserved for the primary key. Applies to UNIQUE and plain KEY/INDEX in CREATE
// TABLE, CREATE INDEX and ALTER TABLE ADD; a FULLTEXT index named PRIMARY is
// accepted by TiDB and the catalog keeps no fulltext index anyway. An empty
// (auto-allocated) name is left to the catalog's own naming.
func checkIndexName(name string) error {
	if strings.EqualFold(name, "PRIMARY") {
		return errWrongNameForIndex(name)
	}
	return nil
}

// disallowedGeneratedColumnFuncs is the set of functions TiDB v8.5.5 refuses
// inside a generated-column expression (3102). Every entry was measured against
// a live container (2026-11-09); the list holds only names TiDB actually
// rejects, so a generated column built from deterministic functions still
// succeeds. Matched case-insensitively and at any depth of the expression.
var disallowedGeneratedColumnFuncs = map[string]bool{
	"now":               true,
	"sysdate":           true,
	"curdate":           true,
	"current_date":      true,
	"curtime":           true,
	"current_time":      true,
	"current_timestamp": true,
	"localtime":         true,
	"localtimestamp":    true,
	"utc_date":          true,
	"utc_time":          true,
	"utc_timestamp":     true,
	"unix_timestamp":    true,
	"rand":              true,
	"uuid":              true,
	"uuid_short":        true,
	"connection_id":     true,
	"last_insert_id":    true,
	"version":           true,
	"database":          true,
	"user":              true,
	"current_user":      true,
	"benchmark":         true,
	"sleep":             true,
	"found_rows":        true,
	"row_count":         true,
	"get_lock":          true,
	"release_lock":      true,
}

// checkGeneratedColumnFunctions rejects a generated-column expression that calls
// a disallowed function anywhere in its tree, which TiDB v8.5.5 answers 3102
// for. The whole expression is walked so a disallowed function nested inside a
// deterministic one (CONCAT('x', NOW())) is caught too.
func checkGeneratedColumnFunctions(colName string, expr nodes.ExprNode) error {
	if expr == nil {
		return nil
	}
	var bad bool
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if bad {
			return false
		}
		if fn, ok := n.(*nodes.FuncCallExpr); ok {
			if disallowedGeneratedColumnFuncs[strings.ToLower(fn.Name)] {
				bad = true
				return false
			}
		}
		return true
	})
	if bad {
		return errGeneratedColumnFunction(colName)
	}
	return nil
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
