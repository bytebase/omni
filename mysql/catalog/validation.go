package catalog

import (
	"fmt"
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

func validateColumnDefSemantics(colDef *nodes.ColumnDef, ignoreSpace bool) error {
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
		if exprHasDisallowedDefaultFunc(expr, ignoreSpace) {
			return errDefaultValExprDisallowedFunc(colDef.Name)
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

// validateGeneratedColumnSelfReference refuses a generated column whose own
// expression names the column being added: MySQL 8.0.46 answers 3107, "Generated
// column can refer only to generated columns defined prior to it", even though
// the name is not yet a column. It runs after the duplicate-column check, which
// answers 1060 first when the table already holds the name.
func validateGeneratedColumnSelfReference(colDef *nodes.ColumnDef) error {
	if colDef == nil || colDef.Generated == nil || colDef.Generated.Expr == nil {
		return nil
	}
	self := false
	nodes.Inspect(colDef.Generated.Expr, func(n nodes.Node) bool {
		if cr, ok := n.(*nodes.ColumnRef); ok && cr.Table == "" && cr.Schema == "" &&
			cr.Column != "" && strings.EqualFold(cr.Column, colDef.Name) {
			self = true
			return false
		}
		return true
	})
	if self {
		return &Error{Code: 3107, SQLState: "HY000",
			Message: "Generated column can refer only to generated columns defined prior to it."}
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
	case "float", "double":
		// `FLOAT(256,2)` and `DOUBLE(256,2)` answer 1439 "max = 255"; without a
		// scale `FLOAT(256)` answers 1063, and only a width past 4294967295 is
		// 1439 there. The parser keeps no flag for a scale written as 0, so
		// `FLOAT(256,0)` reads as `FLOAT(256)` here.
		if colDef.TypeName.Scale > 0 {
			max = 255
		}
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
// leaves, when not nil, is the columns the statement leaves (alterFinalColumns),
// which MySQL judges the keys of an ALTER TABLE against: a key may name a column
// a later command adds, and may not name one the statement drops.
func validateKeyPartsExist(tbl *Table, leaves map[string]bool, parts []*nodes.IndexColumn) error {
	has := func(name string) bool {
		if leaves != nil {
			return leaves[toLower(name)]
		}
		return tbl.GetColumn(name) != nil
	}
	for _, ic := range parts {
		if ic == nil || ic.Functional {
			continue
		}
		ref, ok := ic.Expr.(*nodes.ColumnRef)
		if !ok {
			continue
		}
		if !has(ref.Column) {
			return errKeyColumnDoesNotExist(ref.Column)
		}
	}
	return nil
}

// validateKeyColumnsExist is validateKeyPartsExist for a constraint, whose key
// parts may also be written as a plain column list.
func validateKeyColumnsExist(tbl *Table, leaves map[string]bool, con *nodes.Constraint) error {
	if con == nil {
		return nil
	}
	if len(con.IndexColumns) == 0 {
		for _, name := range con.Columns {
			if leaves != nil {
				if !leaves[toLower(name)] {
					return errKeyColumnDoesNotExist(name)
				}
			} else if tbl.GetColumn(name) == nil {
				return errKeyColumnDoesNotExist(name)
			}
		}
		return nil
	}
	return validateKeyPartsExist(tbl, leaves, con.IndexColumns)
}

// storageEngines are the engines a table may keep on mysql:8.0.46, from its SHOW
// ENGINES: InnoDB, MyISAM, MEMORY, MRG_MYISAM, BLACKHOLE, CSV and ARCHIVE, with
// MEMORY's other name HEAP and MRG_MYISAM's MERGE. FEDERATED, NDB and ndbinfo
// are listed there but disabled, and naming one of them, or any engine the
// server does not have, is 1286. `default` is the sentinel MySQL takes quoted,
// `ENGINE = 'DEFAULT'`, for the server's configured default engine.
// PERFORMANCE_SCHEMA is listed as usable but holds the performance schema's own
// tables only: CREATE TABLE answers 1683 for it and ALTER TABLE 1031, so it is
// refused where those are, not while the statement is prepared.
var storageEngines = map[string]bool{
	"innodb":     true,
	"myisam":     true,
	"memory":     true,
	"heap":       true,
	"mrg_myisam": true,
	"merge":      true,
	"default":    true,
	"blackhole":  true,
	"csv":        true,
	"archive":    true,
}

// validateStorageEngine refuses an ENGINE option naming an engine the pinned
// server does not ship: 1286. MySQL answers it while preparing the statement, so
// it comes before the table is looked up and before the database is resolved —
// unless the session's sql_mode leaves NO_ENGINE_SUBSTITUTION out, and then it
// substitutes the default engine and takes the statement with a warning.
func validateStorageEngine(name string, refuseUnknown bool) error {
	key := strings.ToLower(strings.TrimSpace(name))
	if storageEngines[key] || key == performanceSchemaEngine {
		return nil
	}
	if refuseUnknown {
		return errUnknownStorageEngine(name)
	}
	return nil
}

// performanceSchemaEngine is the engine the performance schema's own tables use.
const performanceSchemaEngine = "performance_schema"

// canonicalEngineNames are the engines whose aliases MySQL stores and reports
// under another name: `ENGINE = HEAP` reads back MEMORY and `ENGINE = MERGE`
// MRG_MyISAM (measured on mysql:8.0.46).
var canonicalEngineNames = map[string]string{
	"heap":  "MEMORY",
	"merge": "MRG_MyISAM",
}

// storedEngineName is the engine a table or partition keeps for an ENGINE
// option: a name it may keep, with an alias written as MySQL stores it, and the
// catalog's empty default — InnoDB — for `DEFAULT` and for a name MySQL
// substitutes.
func storedEngineName(name string) string {
	key := strings.ToLower(strings.TrimSpace(name))
	if !storageEngines[key] || key == "default" {
		return ""
	}
	if canonical, ok := canonicalEngineNames[key]; ok {
		return canonical
	}
	return name
}

// compositeSQLModes are the sql_mode values MySQL expands into a list, as
// @@sql_mode reports them afterwards. Only NO_ENGINE_SUBSTITUTION matters here,
// and TRADITIONAL carries it where ANSI does not (both measured on
// mysql:8.0.46).
var compositeSQLModes = map[string][]string{
	"DEFAULT": {
		"ONLY_FULL_GROUP_BY", "STRICT_TRANS_TABLES", "NO_ZERO_IN_DATE",
		"NO_ZERO_DATE", "ERROR_FOR_DIVISION_BY_ZERO", "NO_ENGINE_SUBSTITUTION",
	},
	"TRADITIONAL": {
		"STRICT_TRANS_TABLES", "STRICT_ALL_TABLES", "NO_ZERO_IN_DATE",
		"NO_ZERO_DATE", "ERROR_FOR_DIVISION_BY_ZERO", "NO_ENGINE_SUBSTITUTION",
	},
	"ANSI": {"REAL_AS_FLOAT", "PIPES_AS_CONCAT", "ANSI_QUOTES", "IGNORE_SPACE", "ONLY_FULL_GROUP_BY"},
}

// partitionEngineName is storedEngineName for a partition of a CREATE TABLE
// statement, where a loader's snapshot keeps the engine its server reported.
func partitionEngineName(name string, declared bool) string {
	if declared {
		return storedEngineName(name)
	}
	return name
}

// noEngineSubstitution reports whether the session's sql_mode refuses an engine
// the server does not ship. The server's default mode, `DEFAULT`, has it; a list
// without it — an empty one included — makes MySQL substitute the default engine
// and take the statement.
func (c *Catalog) noEngineSubstitution() bool {
	return sqlModeHas(c.session.SQLMode, "NO_ENGINE_SUBSTITUTION")
}

func (c *Catalog) isStrictSQLMode() bool {
	return sqlModeHas(c.session.SQLMode, "STRICT_TRANS_TABLES") || sqlModeHas(c.session.SQLMode, "STRICT_ALL_TABLES")
}

func (c *Catalog) isIgnoreSpaceSQLMode() bool {
	return sqlModeHas(c.session.SQLMode, "IGNORE_SPACE")
}

func (c *Catalog) isNoZeroDateSQLMode() bool {
	return c.isStrictSQLMode() && sqlModeHas(c.session.SQLMode, "NO_ZERO_DATE")
}

func (c *Catalog) isNoZeroInDateSQLMode() bool {
	return c.isStrictSQLMode() && sqlModeHas(c.session.SQLMode, "NO_ZERO_IN_DATE")
}

func (c *Catalog) isAllowInvalidDatesSQLMode() bool {
	return sqlModeHas(c.session.SQLMode, "ALLOW_INVALID_DATES")
}

func (c *Catalog) isTimeTruncateFractionalSQLMode() bool {
	return sqlModeHas(c.session.SQLMode, "TIME_TRUNCATE_FRACTIONAL")
}

// sqlModeHas reports whether mode holds the flag, expanding a composite mode as
// MySQL does. `DEFAULT` expands to MySQL 8.0's default mode list.
func sqlModeHas(mode, flag string) bool {
	mode = strings.ToUpper(strings.TrimSpace(mode))
	for _, part := range strings.Split(mode, ",") {
		part = strings.ToUpper(strings.TrimSpace(part))
		if part == flag {
			return true
		}
		for _, expanded := range compositeSQLModes[part] {
			if expanded == flag {
				return true
			}
		}
	}
	return false
}

// engineOptionName reports whether a table option is the engine, as the table
// and its partitions spell it.
func engineOptionName(name string) bool {
	return strings.EqualFold(name, "engine") || strings.EqualFold(name, "storage engine")
}

// partitionBoundColumn returns the first column a partition definition's bound
// names, if it names one. A bound is a constant expression, so MySQL reads a
// name in it as a column of the partition function and answers 1054 "Unknown
// column 'x' in 'partition function'" for it — even for a name the table has.
func partitionBoundColumn(v nodes.Node) (string, bool) {
	found := ""
	var visit func(nodes.Node)
	visit = func(n nodes.Node) {
		if found != "" || n == nil {
			return
		}
		switch node := n.(type) {
		case *nodes.List:
			for _, item := range node.Items {
				visit(item)
			}
		case nodes.ExprNode:
			nodes.Inspect(node, func(n nodes.Node) bool {
				if found != "" {
					return false
				}
				if ref, ok := n.(*nodes.ColumnRef); ok && ref.Column != "" {
					found = ref.Column
					return false
				}
				return true
			})
		}
	}
	visit(v)
	return found, found != ""
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

// validateBareNowDefault rejects a bare NOW()-family DEFAULT on a column whose
// type is neither DATETIME nor TIMESTAMP: mysql:8.0.46 reads such a default as a
// current-timestamp value, which only those two types hold, and answers 1067
// for a TINYINT, INT, DATE, TIME, YEAR or VARCHAR column alike. Judged while the
// statement is prepared, before the database and the table are resolved and
// before a duplicate column: `ALTER TABLE nodb.a ADD x INT DEFAULT NOW()`, and
// the same over an existing table, are 1067 (read on mysql:8.0.46).
//
// Only the bare form is this rule's. A parenthesized `DEFAULT (NOW())` is a
// general expression default, which the same server takes on those types, so
// the parenthesis is not read through here.
func validateBareNowDefault(colDef *nodes.ColumnDef) error {
	if colDef == nil || colDef.TypeName == nil {
		return nil
	}
	typeName := normalizedColumnDataType(colDef.TypeName)
	for _, expr := range columnDefDefaultExprs(colDef) {
		fn, ok := expr.(*nodes.FuncCallExpr)
		if !ok {
			continue
		}
		switch strings.ToLower(fn.Name) {
		case "now", "current_timestamp", "localtime", "localtimestamp":
			if err := validateExpressionSemantics(fn); err != nil {
				return err
			}
			if typeName != "datetime" && typeName != "timestamp" {
				return errInvalidDefault(colDef.Name)
			}
			if fsp, ok := temporalFuncFSP(fn); ok && fsp != colDef.TypeName.Length {
				return errInvalidDefault(colDef.Name)
			}
		}
	}
	return nil
}

// validateColumnOnUpdate rejects an ON UPDATE CURRENT_TIMESTAMP clause whose
// column is neither DATETIME nor TIMESTAMP, or whose FSP mismatches the column,
// at prepare time ahead of database/table lookup and duplicate-column checks
// (measured on mysql:8.0.46).
func validateColumnOnUpdate(colDef *nodes.ColumnDef) error {
	if colDef == nil || colDef.TypeName == nil || colDef.OnUpdate == nil {
		return nil
	}
	typeName := normalizedColumnDataType(colDef.TypeName)
	return validateTemporalExprForColumn(colDef.Name, typeName, colDef.TypeName.Length, colDef.OnUpdate, false)
}

func validateCheckExpr(name string, expr nodes.ExprNode) error {
	if !checkExprIsBoolean(expr) {
		return errCheckNotBoolean(name)
	}
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

// checkExprIsBoolean reports whether a CHECK expression's top level is one
// MySQL takes. The server refuses a bare column, literal, arithmetic, bit
// operation, unary other than NOT, or CAST with 3812 while it prepares the
// statement, even when the column does not exist; it takes a comparison, a
// logical operator, IS/IN/BETWEEN/LIKE/REGEXP, EXISTS, and the boolean builtins
// below. An expression this does not recognize is left to the other checks, so
// a form the catalog does not model is not refused by this one.
func checkExprIsBoolean(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.ParenExpr:
		return checkExprIsBoolean(e.Expr)
	case *nodes.BinaryExpr:
		switch e.Op {
		case nodes.BinOpEq, nodes.BinOpNe, nodes.BinOpLt, nodes.BinOpGt, nodes.BinOpLe, nodes.BinOpGe,
			nodes.BinOpAnd, nodes.BinOpOr, nodes.BinOpXor, nodes.BinOpRegexp, nodes.BinOpLikeEscape,
			nodes.BinOpNullSafeEq, nodes.BinOpSoundsLike:
			return true
		}
		return false
	case *nodes.UnaryExpr:
		return e.Op == nodes.UnaryNot
	case *nodes.IsExpr, *nodes.InExpr, *nodes.BetweenExpr, *nodes.LikeExpr, *nodes.ExistsExpr,
		*nodes.MemberOfExpr, *nodes.MatchExpr:
		return true
	case *nodes.BoolLit:
		return true
	case *nodes.ColumnRef, *nodes.IntLit, *nodes.FloatLit, *nodes.StringLit, *nodes.NullLit,
		*nodes.BitLit, *nodes.CastExpr, *nodes.ConvertExpr, *nodes.CollateExpr,
		*nodes.WeightStringExpr, *nodes.IntervalExpr:
		return false
	case *nodes.FuncCallExpr:
		switch strings.ToLower(e.Name) {
		case "json_valid", "isnull", "strcmp":
			return true
		}
	}
	return true
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

// validateCreateSpatialRefSys refuses a CREATE SPATIAL REFERENCE SYSTEM the server
// will not take: SRID 0 (3714), a NAME it does not write (3708) or writes empty or
// padded (3710), a DEFINITION it does not write (3708) or writes empty (3517). A
// definition it writes in full is not read: the catalog does not parse the WKT, so
// it takes the statement. Measured on mysql:8.0.46.
func validateCreateSpatialRefSys(stmt *nodes.CreateSpatialRefSysStmt) error {
	if stmt.SRID == 0 {
		return &Error{Code: 3714, SQLState: "SR000", Message: "SRID 0 is not modifiable."}
	}
	if !stmt.HasName {
		return errSRSMandatoryAttribute("NAME")
	}
	if stmt.Name == "" || strings.TrimSpace(stmt.Name) != stmt.Name {
		return &Error{Code: 3710, SQLState: "SR006",
			Message: "The spatial reference system name can't be an empty string or start or end with whitespace."}
	}
	if !stmt.HasDefinition {
		return errSRSMandatoryAttribute("DEFINITION")
	}
	if stmt.Definition == "" {
		return &Error{Code: 3517, SQLState: "SR002",
			Message: fmt.Sprintf("Can't parse the spatial reference system definition of SRID %d.", stmt.SRID)}
	}
	return nil
}
