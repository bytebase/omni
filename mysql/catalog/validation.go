package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

func validateExpressionSemantics(node nodes.Node) error {
	var err error
	nodes.Inspect(node, func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		fn, ok := n.(*nodes.FuncCallExpr)
		if !ok {
			return true
		}
		name := strings.ToLower(fn.Name)
		// An argument alias parses on any generic call, and MySQL refuses it on
		// every function but a loadable UDF, which the catalog does not hold: 1583
		// for a native function, 1584 for a stored one, named or not.
		for _, alias := range fn.ArgAliases {
			if alias == nil {
				continue
			}
			if fn.Schema == "" && nativeFunctions[name] {
				err = errWrongParametersToNativeFct(name)
			} else {
				err = errWrongParametersToStoredFct(fn.Name)
			}
			return false
		}
		switch name {
		case "now", "current_timestamp", "localtime", "localtimestamp", "sysdate",
			"utc_timestamp", "curtime", "current_time", "utc_time":
			if fsp, ok := firstIntArg(fn); ok && fsp > 6 {
				err = errTooBigPrecision(fsp, name)
			}
		case "curdate", "current_date", "utc_date":
			if len(fn.Args) > 0 {
				err = errWrongArguments(strings.ToUpper(name))
			}
		}
		return err == nil
	})
	return err
}

func validateColumnDefSemantics(colDef *nodes.ColumnDef) error {
	if colDef == nil || colDef.TypeName == nil {
		return nil
	}
	typeName := normalizedColumnDataType(colDef.TypeName)
	switch typeName {
	case "datetime", "timestamp", "time":
		if colDef.TypeName.Length > 6 {
			return errTooBigPrecision(colDef.TypeName.Length, typeName)
		}
	case "year":
		if colDef.TypeName.Length != 0 && colDef.TypeName.Length != 4 {
			return errInvalidYearColumnLength()
		}
	case "decimal", "numeric", "dec", "fixed":
		if colDef.TypeName.Length > 65 {
			return errTooBigPrecision(colDef.TypeName.Length, typeName)
		}
		if colDef.TypeName.Scale > 30 {
			return &Error{Code: 1425, SQLState: "42000", Message: "Too big scale specified for column"}
		}
		if colDef.TypeName.Scale > colDef.TypeName.Length && colDef.TypeName.Length > 0 {
			return &Error{Code: 1427, SQLState: "42000", Message: "For float(M,D), double(M,D) or decimal(M,D), M must be >= D"}
		}
	case "varchar":
		if colDef.TypeName.Length == 0 {
			return &Error{Code: 1064, SQLState: "42000", Message: "VARCHAR requires a length"}
		}
	}
	if colDef.Generated != nil && (colDef.DefaultValue != nil || hasColumnDefaultConstraint(colDef)) {
		return errInvalidDefault(colDef.Name)
	}
	if colDef.Generated != nil {
		if err := validateGeneratedExpr(colDef.Name, colDef.Generated.Expr); err != nil {
			return err
		}
		if !colDef.Generated.Stored && columnDefHasPrimaryKey(colDef) {
			return errPrimaryKeyOnGeneratedColumn(colDef.Name)
		}
	}
	if hasExplicitNullPrimaryKey(colDef) {
		return &Error{Code: 1171, SQLState: "42000", Message: "All parts of a PRIMARY KEY must be NOT NULL"}
	}
	if !columnDefNullable(colDef) && hasNullDefault(colDef) {
		return errInvalidDefault(colDef.Name)
	}
	if isLiteralDefaultForbiddenType(typeName) {
		if isLiteralDefault(colDef.DefaultValue) || hasLiteralColumnDefaultConstraint(colDef) {
			return &Error{Code: 1101, SQLState: "42000", Message: fmtCantHaveDefault(colDef.Name)}
		}
	}

	if err := validateTemporalExprForColumn(colDef.Name, typeName, colDef.TypeName.Length, colDef.DefaultValue, true); err != nil {
		return err
	}
	if err := validateTemporalExprForColumn(colDef.Name, typeName, colDef.TypeName.Length, colDef.OnUpdate, false); err != nil {
		return err
	}
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault {
			if err := validateTemporalExprForColumn(colDef.Name, typeName, colDef.TypeName.Length, cc.Expr, true); err != nil {
				return err
			}
		}
	}
	// A spaced gated keyword is a generic stored-function call, which an
	// expression default forbids: mysql:8.0 answers 3770 for
	// `DEFAULT (position (1))`. (Other disallowed functions in a general default
	// expression are not validated here yet.)
	for _, expr := range columnDefDefaultExprs(colDef) {
		if exprHasSpacedFuncCall(expr) {
			return &Error{Code: 3770, SQLState: "HY000",
				Message: "Default value expression of column '" + colDef.Name + "' contains a disallowed function."}
		}
	}
	return nil
}

// columnDefDefaultExprs returns the column's default expressions, from the field
// and from any DEFAULT column constraint.
func columnDefDefaultExprs(colDef *nodes.ColumnDef) []nodes.ExprNode {
	exprs := []nodes.ExprNode{colDef.DefaultValue}
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault {
			exprs = append(exprs, cc.Expr)
		}
	}
	return exprs
}

// exprHasSpacedFuncCall reports whether expr holds a function call written with a
// space before its '(' — a generic stored-function call (see SpacedParen).
func exprHasSpacedFuncCall(expr nodes.ExprNode) bool {
	found := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if found {
			return false
		}
		if fn, ok := n.(*nodes.FuncCallExpr); ok && fn.SpacedParen {
			found = true
			return false
		}
		return true
	})
	return found
}

func columnDefHasPrimaryKey(colDef *nodes.ColumnDef) bool {
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrPrimaryKey {
			return true
		}
	}
	return false
}

func hasExplicitNullPrimaryKey(colDef *nodes.ColumnDef) bool {
	hasNull := false
	hasPK := false
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrNull:
			hasNull = true
		case nodes.ColConstrPrimaryKey:
			hasPK = true
		}
	}
	return hasNull && hasPK
}

func columnDefNullable(colDef *nodes.ColumnDef) bool {
	nullable := true
	if colDef.AutoIncrement {
		nullable = false
	}
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrNotNull, nodes.ColConstrPrimaryKey, nodes.ColConstrAutoIncrement:
			nullable = false
		case nodes.ColConstrNull:
			nullable = true
		}
	}
	return nullable
}

func validateGeneratedExpr(colName string, expr nodes.ExprNode) error {
	var invalid bool
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if invalid {
			return false
		}
		switch x := n.(type) {
		case *nodes.SubqueryExpr:
			invalid = true
		case *nodes.VariableRef:
			invalid = true
		case *nodes.FuncCallExpr:
			// A spaced gated keyword is a generic stored-function call, which a
			// generated expression forbids: mysql:8.0 answers 3763 for
			// `b INT AS (position (a))`, where the adjacent builtin is a syntax
			// error at parse. Admitting it would also store the space-dropped
			// `position(a)`, an invalid builtin call.
			if isNondeterministicFunc(x.Name) || x.SpacedParen {
				invalid = true
			}
		}
		return !invalid
	})
	if invalid {
		return &Error{Code: 3763, SQLState: "HY000", Message: "Expression of generated column '" + colName + "' contains a disallowed function"}
	}
	return nil
}

func errPrimaryKeyOnGeneratedColumn(colName string) error {
	return &Error{Code: 3106, SQLState: "HY000", Message: "Defining a virtual generated column as primary key is not supported"}
}

func validatePrimaryKeyColumns(tbl *Table, cols []string) error {
	for _, colName := range cols {
		col := tbl.GetColumn(colName)
		if col != nil && col.Generated != nil && !col.Generated.Stored {
			return errPrimaryKeyOnGeneratedColumn(colName)
		}
	}
	return nil
}

// maxDisplayWidth is the widest length MySQL keeps in a type's length argument;
// past it the statement is refused with 1439 before anything else about the
// column or the table is looked at.
const maxDisplayWidth = 4294967295

// validateColumnDisplayWidth refuses a column whose declared length is past the
// widest MySQL holds for its type: 1439. Measured on mysql:8.0.46: an integer's
// display width stops at 255 and BIT's at 64, and every other type at
// 4294967295. The type's own limit is a later check, so `CHAR(4294967296)`
// answers 1439 where `CHAR(65536)` answers 1074, and `DECIMAL(4294967296)` 1439
// where `DECIMAL(66)` answers 1426. The width is checked before the table is
// opened and before the column's name: `ALTER TABLE nosuchdb.a ADD a
// INT(256)` and `ALTER TABLE a ADD existing_name INT(256)` are both 1439.
func validateColumnDisplayWidth(colDef *nodes.ColumnDef) error {
	if colDef == nil || colDef.TypeName == nil {
		return nil
	}
	max := maxDisplayWidth
	switch normalizedTypeName(colDef.TypeName.Name) {
	case "tinyint", "smallint", "mediumint", "int", "bigint":
		max = 255
	case "bit":
		max = 64
	}
	if colDef.TypeName.Length > max {
		return errTooBigDisplayWidth(colDef.Name, max)
	}
	return nil
}

// parserOptionName returns the parser a FULLTEXT key names WITH PARSER, and
// whether it names one at all.
func parserOptionName(opts []*nodes.IndexOption) (string, bool) {
	for _, opt := range opts {
		if opt == nil || !strings.EqualFold(opt.Name, "PARSER") {
			continue
		}
		if s, ok := opt.Value.(*nodes.StringLit); ok {
			return s.Value, true
		}
		return "", true
	}
	return "", false
}

// validateFulltextParserOption refuses a FULLTEXT key whose parser MySQL does
// not have: 1128, which it answers while preparing the statement, before it
// opens the table and before it judges the key's columns. `ngram` is the one
// fulltext parser its server ships. fulltext says whether the key is one: a
// parser named on any other key is not a parser at all.
func validateFulltextParserOption(fulltext bool, opts []*nodes.IndexOption) error {
	if !fulltext {
		return nil
	}
	name, ok := parserOptionName(opts)
	if !ok || strings.EqualFold(name, "ngram") {
		return nil
	}
	return errFunctionNotDefined(name)
}

// validateFulltextParser is validateFulltextParserOption for a constraint.
func validateFulltextParser(con *nodes.Constraint) error {
	if con == nil || con.Type != nodes.ConstrFulltextIndex {
		return nil
	}
	return validateFulltextParserOption(true, con.IndexOptions)
}

// validateKeyPartsExist refuses a key part that names a column the table does
// not have: 1072. A functional part is an expression, and names no column.
//
// The check is made where the command runs, so a key naming a column a later
// command of the same ALTER TABLE adds is refused; MySQL judges the columns of
// the table the statement leaves. No statement of simulate.tsv reaches that.
func validateKeyPartsExist(tbl *Table, parts []*nodes.IndexColumn) error {
	for _, ic := range parts {
		if ic == nil || ic.Functional {
			continue
		}
		ref, ok := ic.Expr.(*nodes.ColumnRef)
		if !ok {
			continue
		}
		if tbl.GetColumn(ref.Column) == nil {
			return errKeyColumnDoesNotExist(ref.Column)
		}
	}
	return nil
}

// validateKeyColumnsExist is validateKeyPartsExist for a constraint, whose key
// parts may also be written as a plain column list.
func validateKeyColumnsExist(tbl *Table, con *nodes.Constraint) error {
	if con == nil {
		return nil
	}
	if len(con.IndexColumns) == 0 {
		for _, name := range con.Columns {
			if tbl.GetColumn(name) == nil {
				return errKeyColumnDoesNotExist(name)
			}
		}
		return nil
	}
	return validateKeyPartsExist(tbl, con.IndexColumns)
}

func validateColumnCheckReferences(colName, conName string, expr nodes.ExprNode) error {
	bad := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if bad {
			return false
		}
		if cr, ok := n.(*nodes.ColumnRef); ok && cr.Column != "" && !strings.EqualFold(cr.Column, colName) {
			bad = true
			return false
		}
		return true
	})
	if bad {
		return &Error{Code: 3823, SQLState: "HY000", Message: "Column check constraint cannot reference other columns"}
	}
	return validateCheckExpr(conName, expr)
}

func hasColumnDefaultConstraint(colDef *nodes.ColumnDef) bool {
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault && cc.Expr != nil {
			return true
		}
	}
	return false
}

func hasNullDefault(colDef *nodes.ColumnDef) bool {
	if isNullDefault(colDef.DefaultValue) {
		return true
	}
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault && isNullDefault(cc.Expr) {
			return true
		}
	}
	return false
}

func hasLiteralColumnDefaultConstraint(colDef *nodes.ColumnDef) bool {
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault && isLiteralDefault(cc.Expr) {
			return true
		}
	}
	return false
}

func isLiteralDefault(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.StringLit, *nodes.IntLit, *nodes.FloatLit, *nodes.BitLit, *nodes.BoolLit:
		return true
	case *nodes.ParenExpr:
		return isLiteralDefault(e.Expr)
	default:
		return false
	}
}

func isNullDefault(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.NullLit:
		return true
	case *nodes.ParenExpr:
		return isNullDefault(e.Expr)
	default:
		return false
	}
}

func isLiteralDefaultForbiddenType(typeName string) bool {
	switch strings.ToLower(typeName) {
	case "blob", "tinyblob", "mediumblob", "longblob",
		"text", "tinytext", "mediumtext", "longtext",
		"json", "geometry", "point", "linestring", "polygon",
		"multipoint", "multilinestring", "multipolygon", "geometrycollection", "geomcollection":
		return true
	default:
		return false
	}
}

func fmtCantHaveDefault(colName string) string {
	return "BLOB, TEXT, GEOMETRY or JSON column '" + colName + "' can't have a default value"
}

func validateTemporalExprForColumn(colName, typeName string, colFsp int, expr nodes.ExprNode, isDefault bool) error {
	if expr == nil {
		return nil
	}
	if err := validateExpressionSemantics(expr); err != nil {
		return err
	}
	fn, ok := unwrapTemporalFunc(expr)
	if !ok {
		return nil
	}
	name := strings.ToLower(fn.Name)
	if isDefault && (name == "sysdate" || name == "utc_timestamp") {
		return errInvalidDefault(colName)
	}
	if !isDefault && typeName != "datetime" && typeName != "timestamp" {
		return errInvalidOnUpdate(colName)
	}
	if typeName == "datetime" || typeName == "timestamp" {
		if fsp, ok := temporalFuncFSP(fn); ok && fsp != colFsp {
			if !isDefault {
				return errInvalidOnUpdate(colName)
			}
			return errInvalidDefault(colName)
		}
	}
	return nil
}

func validateCheckExpr(name string, expr nodes.ExprNode) error {
	var invalid bool
	exprSQL := strings.ToLower(nodeToSQL(expr))
	if strings.Contains(exprSQL, "select ") || strings.Contains(exprSQL, "(select") {
		return errCheckConstraintNotAllowed(name)
	}
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if invalid {
			return false
		}
		switch x := n.(type) {
		case *nodes.SubqueryExpr:
			invalid = true
		case *nodes.VariableRef:
			invalid = true
		case *nodes.FuncCallExpr:
			// A spaced gated keyword is a generic stored-function call, which a
			// check constraint forbids: mysql:8.0 answers 3814 for
			// `CHECK (sum (a) > 0)`.
			if isNondeterministicFunc(x.Name) || x.SpacedParen {
				invalid = true
			}
		}
		return !invalid
	})
	if invalid {
		return errCheckConstraintNotAllowed(name)
	}
	return nil
}

func isNondeterministicFunc(name string) bool {
	switch strings.ToLower(name) {
	case "now", "current_timestamp", "localtime", "localtimestamp", "sysdate",
		"utc_timestamp", "curtime", "current_time", "utc_time", "curdate",
		"current_date", "utc_date", "rand", "uuid", "uuid_short", "connection_id":
		return true
	default:
		return false
	}
}

func firstIntArg(fn *nodes.FuncCallExpr) (int, bool) {
	if fn == nil || len(fn.Args) == 0 {
		return 0, false
	}
	if lit, ok := fn.Args[0].(*nodes.IntLit); ok {
		return int(lit.Value), true
	}
	return 0, false
}

func temporalFuncFSP(fn *nodes.FuncCallExpr) (int, bool) {
	if fn == nil {
		return 0, false
	}
	name := strings.ToLower(fn.Name)
	switch name {
	case "now", "current_timestamp", "localtime", "localtimestamp", "sysdate",
		"utc_timestamp":
		if fsp, ok := firstIntArg(fn); ok {
			return fsp, true
		}
		return 0, true
	default:
		return 0, false
	}
}

func unwrapTemporalFunc(expr nodes.ExprNode) (*nodes.FuncCallExpr, bool) {
	switch e := expr.(type) {
	case *nodes.FuncCallExpr:
		return e, true
	case *nodes.ParenExpr:
		return unwrapTemporalFunc(e.Expr)
	default:
		return nil, false
	}
}
