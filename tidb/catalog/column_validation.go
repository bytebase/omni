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

// errInvalidOnUpdate is TiDB 1294 for an ON UPDATE CURRENT_TIMESTAMP clause the
// column cannot carry.
func errInvalidOnUpdate(col string) error {
	return &Error{Code: ErrInvalidOnUpdate, SQLState: sqlState(ErrInvalidOnUpdate),
		Message: fmt.Sprintf("Invalid ON UPDATE clause for '%s' column", col)}
}

// checkTemporalBareDefault rejects a bare (non-parenthesized) function DEFAULT on
// a DATETIME or TIMESTAMP column when the function is not one TiDB v8.5.5 permits
// there. Only CURRENT_TIMESTAMP/NOW and the synonyms LOCALTIME/LOCALTIMESTAMP are
// valid as a datetime default; SYSDATE() and UTC_TIMESTAMP are 1067. A
// parenthesized expression default is checked by checkDefaultExprFunction, and
// the NOW()-family precision match by checkTemporalDefaultFsp.
func checkTemporalBareDefault(colName string, dt *nodes.DataType, defaultExpr nodes.ExprNode) error {
	if dt == nil || defaultExpr == nil {
		return nil
	}
	switch strings.ToLower(dt.Name) {
	case "datetime", "timestamp":
	default:
		return nil
	}
	fn, ok := defaultExpr.(*nodes.FuncCallExpr)
	if !ok {
		return nil
	}
	if nowFamilyDefaultFuncs[strings.ToLower(fn.Name)] {
		return nil
	}
	return errInvalidDefault(colName)
}

// checkTemporalOnUpdate enforces TiDB v8.5.5's 1294 for an ON UPDATE
// CURRENT_TIMESTAMP clause: it is valid only on a TIMESTAMP or DATETIME column,
// and its fractional-seconds precision must match the column's own. DATE ON
// UPDATE NOW() and DATETIME(6) ON UPDATE NOW() (0 vs 6) both fail, while
// DATETIME(6) ON UPDATE NOW(6) and TIMESTAMP ON UPDATE CURRENT_TIMESTAMP pass.
func checkTemporalOnUpdate(colName string, dt *nodes.DataType, onUpdate nodes.ExprNode) error {
	if onUpdate == nil {
		return nil
	}
	if dt == nil {
		return nil
	}
	switch strings.ToLower(dt.Name) {
	case "timestamp", "datetime":
	default:
		return errInvalidOnUpdate(colName)
	}
	fn, ok := onUpdate.(*nodes.FuncCallExpr)
	if !ok {
		return nil
	}
	updFsp := 0
	if len(fn.Args) > 0 {
		if lit, ok := fn.Args[0].(*nodes.IntLit); ok {
			updFsp = int(lit.Value)
		}
	}
	if updFsp != dt.Length {
		return errInvalidOnUpdate(colName)
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

// errFunctionalIndexOnField is TiDB 3762 for an expression index over a bare
// column (INDEX ((a))).
func errFunctionalIndexOnField() error {
	return &Error{Code: ErrFunctionalIndexOnField, SQLState: sqlState(ErrFunctionalIndexOnField),
		Message: "Expression index on a column is not supported. Consider using a regular index instead"}
}

// errFunctionalIndexOnJSONOrGeo is TiDB 3753 for an expression index whose
// expression returns a JSON or GEOMETRY value (e.g. the -> operator).
func errFunctionalIndexOnJSONOrGeo() error {
	return &Error{Code: ErrFunctionalIndexOnJSONOrGeo, SQLState: sqlState(ErrFunctionalIndexOnJSONOrGeo),
		Message: "Cannot create an expression index on a function that returns a JSON or GEOMETRY value"}
}

// errFunctionalIndexOnBlob is TiDB 3757 for an expression index whose expression
// returns a BLOB or TEXT value (e.g. the ->> operator without a CAST).
func errFunctionalIndexOnBlob() error {
	return &Error{Code: ErrFunctionalIndexOnBlob, SQLState: sqlState(ErrFunctionalIndexOnBlob),
		Message: "Cannot create an expression index on an expression that returns a BLOB or TEXT. Please consider using CAST"}
}

// errFunctionalIndexFunctionNotAllowed is TiDB 3758 for an expression index that
// calls a non-deterministic / disallowed function.
func errFunctionalIndexFunctionNotAllowed(fn string) error {
	return &Error{Code: ErrFunctionalIndexFunctionIsNotAllowed, SQLState: sqlState(ErrFunctionalIndexFunctionIsNotAllowed),
		Message: fmt.Sprintf("Expression of expression index 'expression_index' contains a disallowed function: `%s`.", strings.ToUpper(fn))}
}

// validateFunctionalIndexColumns rejects functional (expression) key parts TiDB
// v8.5.5 refuses: a bare column in parentheses (3762), an expression returning
// JSON/GEOMETRY (3753, the -> operator), one returning BLOB/TEXT (3757, the ->>
// operator), or one that calls a non-deterministic / disallowed function (3758).
// Measured against a live container (2026-11-09).
func validateFunctionalIndexColumns(cols []*nodes.IndexColumn) error {
	for _, ic := range cols {
		if !ic.Functional {
			continue
		}
		if _, ok := ic.Expr.(*nodes.ColumnRef); ok {
			return errFunctionalIndexOnField()
		}
		// The expression's return type is set by its OUTERMOST operator: a
		// top-level -> returns JSON (3753) and a top-level ->> returns BLOB/TEXT
		// (3757). A CAST(... AS CHAR) wrapping a ->> changes the return type and
		// is accepted, so the operator is checked only at the top level, never
		// deep in the tree.
		if bin, ok := ic.Expr.(*nodes.BinaryExpr); ok {
			switch bin.Op {
			case nodes.BinOpJsonExtract:
				return errFunctionalIndexOnJSONOrGeo()
			case nodes.BinOpJsonUnquote:
				return errFunctionalIndexOnBlob()
			}
		}
		// A non-deterministic / disallowed function anywhere in the expression
		// is refused (3758).
		var bad error
		nodes.Inspect(ic.Expr, func(n nodes.Node) bool {
			if bad != nil {
				return false
			}
			if fn, ok := n.(*nodes.FuncCallExpr); ok {
				if disallowedGeneratedColumnFuncs[strings.ToLower(fn.Name)] {
					bad = errFunctionalIndexFunctionNotAllowed(fn.Name)
					return false
				}
			}
			return true
		})
		if bad != nil {
			return bad
		}
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

// errIncorrectAutoIncrementDefinition is TiDB 1075 for a table that declares
// more than one AUTO_INCREMENT column.
func errIncorrectAutoIncrementDefinition() error {
	return &Error{Code: ErrIncorrectTableDefinition, SQLState: sqlState(ErrIncorrectTableDefinition),
		Message: "Incorrect table definition; there can be only one auto column and it must be defined as a key"}
}

// errCreateTablePKColumnDoesNotExist is TiDB 1072 as its CREATE TABLE PRIMARY
// KEY path words it. The code is errKeyColumnDoesNotExist's; only the message
// differs, and v8.5.5 is inconsistent about which it uses — a table-level
// PRIMARY KEY in CREATE TABLE gets this wording, the same key added by ALTER
// gets "column does not exist".
func errCreateTablePKColumnDoesNotExist(col string) error {
	return &Error{Code: ErrKeyColumnDoesNotExist, SQLState: sqlState(ErrKeyColumnDoesNotExist),
		Message: fmt.Sprintf("Key column '%s' doesn't exist in table", col)}
}

// emptyStringDefaultForbiddenType reports whether TiDB v8.5.5 refuses a blank
// (empty or whitespace-only) string literal as the DEFAULT of a column of this
// type, with 1067. Measured against a live container on 2026-11-10: date,
// datetime, timestamp, year and every numeric type are refused, while char,
// varchar, binary, varbinary, bit and set take a blank default happily. TIME is
// left out deliberately — v8.5.5 answers 1292 ("Truncated incorrect time
// value") there, a different error this check does not model. BLOB, TEXT and
// JSON never reach here: checkLiteralDefaultAllowed refuses any literal default
// on them first (1101). ENUM is not a type question at all and goes through
// blankStringDefaultRefused.
//
// The type name is the catalog's normalized DataType, so BOOL arrives as
// tinyint and NUMERIC as decimal; the extra spellings are listed anyway so the
// predicate also reads correctly off a parser type name.
func emptyStringDefaultForbiddenType(dataType string) bool {
	switch strings.ToLower(dataType) {
	case "date", "datetime", "timestamp", "year",
		"tinyint", "smallint", "mediumint", "int", "integer", "bigint",
		"decimal", "numeric", "dec", "fixed", "float", "double", "real",
		"bool", "boolean":
		return true
	}
	return false
}

// blankStringDefaultRefused reports whether TiDB v8.5.5 refuses value — a blank
// string literal — as the DEFAULT of a column of this type.
//
// For an ENUM the question is membership, not emptiness: v8.5.5 takes an
// empty-string default on an ENUM that declares an empty member, and
// `ENUM(' ','x') DEFAULT ' '` likewise, because the blank string is a declared
// member; it refuses the same default on `ENUM('x','y')` and
// `ENUM('x') DEFAULT ' '`, where it is not (measured 2026-11-10). A blank
// string has no case, so the members are compared byte for byte and the
// column's collation does not come into it. An ENUM whose members could not be
// read is left alone rather than guessed at.
//
// A non-blank default that the ENUM does not declare (`ENUM('x') DEFAULT 'z'`,
// also 1067) is still over-accepted; it belongs with the other literal-
// conversion gap this file's checkColumnDefaultValue notes.
func blankStringDefaultRefused(dataType, value string, enumMembers []string) bool {
	if strings.EqualFold(dataType, "enum") {
		if len(enumMembers) == 0 {
			return false
		}
		for _, m := range enumMembers {
			if m == value {
				return false
			}
		}
		return true
	}
	return emptyStringDefaultForbiddenType(dataType)
}

// blankStringLit returns the value of expr when it is a string literal holding
// nothing but whitespace.
func blankStringLit(expr nodes.ExprNode) (string, bool) {
	s, ok := expr.(*nodes.StringLit)
	if !ok || strings.TrimSpace(s.Value) != "" {
		return "", false
	}
	return s.Value, true
}

// enumMembersFromColumnType reads the declared members out of a rendered ENUM
// or SET column type, as formatColumnType writes it, undoing escapeEnumValue's
// doubled quote so a member holding a quote reads back whole. It returns nil
// for a type that declares
// none, so a caller cannot tell "no members" from "could not read them" and
// must treat both as unknown.
func enumMembersFromColumnType(columnType string) []string {
	open := strings.IndexByte(columnType, '(')
	if open < 0 {
		return nil
	}
	var members []string
	var cur strings.Builder
	inQuote := false
	for i := open + 1; i < len(columnType); i++ {
		ch := columnType[i]
		switch {
		case inQuote && ch == '\'':
			// A doubled quote is one literal quote inside the member.
			if i+1 < len(columnType) && columnType[i+1] == '\'' {
				cur.WriteByte('\'')
				i++
				continue
			}
			inQuote = false
			members = append(members, cur.String())
			cur.Reset()
		case inQuote:
			cur.WriteByte(ch)
		case ch == '\'':
			inQuote = true
		case ch == ')':
			return members
		}
	}
	// An unterminated list is a rendering this parser does not understand.
	return nil
}

// isNullDefault reports whether expr is the NULL literal.
func isNullDefault(expr nodes.ExprNode) bool {
	_, ok := expr.(*nodes.NullLit)
	return ok
}

// columnDefNotNullable reports whether the column definition resolves to NOT
// NULL. NOT NULL, an inline PRIMARY KEY and AUTO_INCREMENT each make it so, and
// a later explicit NULL takes it back, which is the order TiDB reads them in.
func columnDefNotNullable(colDef *nodes.ColumnDef) bool {
	notNull := colDef.AutoIncrement
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrNotNull, nodes.ColConstrPrimaryKey, nodes.ColConstrAutoIncrement:
			notNull = true
		case nodes.ColConstrNull:
			notNull = false
		}
	}
	return notNull
}

// columnNamedInTableLevelPK reports whether a table-level PRIMARY KEY(...) in
// the statement names the column. TiDB answers 1171 rather than 1067 when such
// a column carries DEFAULT NULL, so the two cases are told apart.
func columnNamedInTableLevelPK(cons []*nodes.Constraint, colName string) bool {
	for _, con := range cons {
		if con.Type != nodes.ConstrPrimaryKey {
			continue
		}
		for _, name := range extractColumnNames(con) {
			if strings.EqualFold(name, colName) {
				return true
			}
		}
	}
	return false
}

// checkColumnDefaultValue rejects the two invalid-default forms TiDB v8.5.5
// refuses and this catalog used to accept: a NULL default on a column that is
// not nullable (1067, or 1171 when the column is a part of a table-level
// PRIMARY KEY), and a blank string default the column's type cannot hold
// (1067). dataType is the catalog's normalized type name; tableCons are the
// statement's table-level constraints, empty on an ALTER that declares none.
//
// Still over-accepted, and left to its own change: a non-blank literal that the
// type cannot hold (`a date DEFAULT 'abc'`, `a int DEFAULT 'abc'`,
// `a enum('x') DEFAULT 'z'`), which needs the engine's full literal conversion
// rather than a shape check.
func checkColumnDefaultValue(colDef *nodes.ColumnDef, dataType string, tableCons []*nodes.Constraint) error {
	expr := columnDefaultExpr(colDef)
	if expr == nil {
		return nil
	}
	if isNullDefault(expr) {
		if columnNamedInTableLevelPK(tableCons, colDef.Name) {
			return errPrimaryKeyCantHaveNull()
		}
		if columnDefNotNullable(colDef) {
			return errInvalidDefault(colDef.Name)
		}
		return nil
	}
	// The declared members come straight from the type, so an ENUM written in
	// this statement is judged on its own member list.
	var enumMembers []string
	if colDef.TypeName != nil {
		enumMembers = colDef.TypeName.EnumValues
	}
	if value, ok := blankStringLit(expr); ok && blankStringDefaultRefused(dataType, value, enumMembers) {
		return errInvalidDefault(colDef.Name)
	}
	return nil
}

// checkExistingColumnDefault holds ALTER TABLE ... ALTER COLUMN c SET DEFAULT
// to the same rules as the CREATE TABLE path, against the column as the catalog
// already holds it: a NULL default on a primary-key column is 1171, on any
// other non-nullable column 1067, and a blank string default the column's type
// cannot hold is 1067. An ENUM's members are read back out of its rendered
// column type, which is the only place this catalog keeps them.
func checkExistingColumnDefault(tbl *Table, col *Column, expr nodes.ExprNode) error {
	if expr == nil {
		return nil
	}
	if isNullDefault(expr) {
		if columnInPrimaryKey(tbl, col.Name) {
			return errPrimaryKeyCantHaveNull()
		}
		if !col.Nullable {
			return errInvalidDefault(col.Name)
		}
		return nil
	}
	value, ok := blankStringLit(expr)
	if !ok {
		return nil
	}
	var enumMembers []string
	if strings.EqualFold(col.DataType, "enum") {
		enumMembers = enumMembersFromColumnType(col.ColumnType)
	}
	if blankStringDefaultRefused(col.DataType, value, enumMembers) {
		return errInvalidDefault(col.Name)
	}
	return nil
}

// columnInPrimaryKey reports whether the table's primary key covers the column.
func columnInPrimaryKey(tbl *Table, colName string) bool {
	for _, idx := range tbl.Indexes {
		if !idx.Primary {
			continue
		}
		for _, ic := range idx.Columns {
			if strings.EqualFold(ic.Name, colName) {
				return true
			}
		}
	}
	return false
}

// checkSingleAutoIncrementColumn rejects a table that declares more than one
// AUTO_INCREMENT column, which TiDB v8.5.5 answers 1075 for. Unlike MySQL,
// v8.5.5 accepts a lone AUTO_INCREMENT column that backs no key (`CREATE TABLE
// t (a int auto_increment)` succeeds, measured 2026-11-10), so only the count
// is checked. An AUTO_RANDOM column is not an auto column for this rule:
// v8.5.5 takes AUTO_RANDOM and AUTO_INCREMENT on one table.
func checkSingleAutoIncrementColumn(cols []*Column) error {
	seen := 0
	for _, col := range cols {
		if col.AutoIncrement {
			seen++
			if seen > 1 {
				return errIncorrectAutoIncrementDefinition()
			}
		}
	}
	return nil
}
