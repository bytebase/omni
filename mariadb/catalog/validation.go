package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

func validateExpressionSemantics(node nodes.Node) error {
	var err error
	nodes.Inspect(node, func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		// A window name resolves in the query that declares it, and Inspect
		// reaches every query, so each is checked against its own clause.
		if sel, ok := n.(*nodes.SelectStmt); ok {
			if e := validateWindowReferences(sel); e != nil {
				err = e
				return false
			}
			return true
		}
		fn, ok := n.(*nodes.FuncCallExpr)
		if !ok {
			return true
		}
		name := strings.ToLower(fn.Name)
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

// validateColumnDefSemantics checks one column definition. indexed holds the
// lowercased names of the columns a key outside the definition covers, which
// decides whether a VIRTUAL generated column's value is materialized; it may
// be nil.
func validateColumnDefSemantics(colDef *nodes.ColumnDef, indexed map[string]bool) error {
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
		// MariaDB takes any YEAR width: it keeps 2, which it still has, and
		// normalizes every other to 4. MySQL 8.0 answers
		// ER_INVALID_YEAR_COLUMN_LENGTH (1818) for anything but 4. Measured on
		// mariadb:11.8.8: YEAR(2) → year(2), YEAR(3) and YEAR(5) → year(4).
	case "decimal", "numeric", "dec", "fixed":
		if colDef.TypeName.Length > 65 {
			return errTooBigPrecision(colDef.TypeName.Length, typeName)
		}
		// MariaDB's scale maximum is 38, where MySQL 8.0's is 30, and its
		// message names the limit rather than the value asked for. Measured on
		// mariadb:11.8.8: DECIMAL(40,31) is accepted and DECIMAL(40,39)
		// answers "Too big scale specified for 'd'. Maximum is 38".
		if colDef.TypeName.Scale > 38 {
			return &Error{Code: 1425, SQLState: "42000",
				Message: fmt.Sprintf("Too big scale specified for '%s'. Maximum is 38", colDef.Name)}
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
		// A key over the column makes a VIRTUAL value materialized, and the
		// server then refuses a non-deterministic expression for it.
		stored := colDef.Generated.Stored || columnDefHasKey(colDef) ||
			indexed[strings.ToLower(colDef.Name)]
		if err := validateGeneratedExpr(colDef.Name, colDef.Generated.Expr, stored); err != nil {
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
	return nil
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

// validateGeneratedExpr checks a generated column's expression. stored says
// whether the value is materialized — a STORED column, or a VIRTUAL one an
// index covers — which is what decides whether a non-deterministic function is
// allowed.
func validateGeneratedExpr(colName string, expr nodes.ExprNode, stored bool) error {
	var invalid bool
	var badFunc string
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
			// A non-deterministic function is allowed only in a VIRTUAL column
			// no index covers: MariaDB answers 1901, "Function or expression
			// 'rand()' cannot be used in the GENERATED ALWAYS AS clause of
			// `a`", for a STORED column and for a VIRTUAL one with a key over
			// it. Measured on mariadb:11.8.8. MySQL 8.0 refuses it everywhere.
			if stored && badFunc == "" {
				badFunc = nondeterministicFuncName(x.Name)
			}
		}
		return !invalid
	})
	if invalid {
		return &Error{Code: 3763, SQLState: "HY000", Message: "Expression of generated column '" + colName + "' contains a disallowed function"}
	}
	if badFunc != "" {
		return errNondeterministicGeneratedColumn(badFunc, colName)
	}
	return nil
}

// columnDefHasKey reports whether the column definition carries an index of
// its own: a PRIMARY KEY or UNIQUE written on the column. A table-level key
// over a generated column is checked where the index is added.
func columnDefHasKey(colDef *nodes.ColumnDef) bool {
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrPrimaryKey, nodes.ColConstrUnique:
			return true
		}
	}
	return false
}

// errNondeterministicGeneratedColumn reports 1901, which MariaDB answers for a
// non-deterministic expression in a materialized generated column.
func errNondeterministicGeneratedColumn(fn, colName string) error {
	return errUnusableInGeneratedColumn(fn+"()", colName)
}

// errUnusableInGeneratedColumn is 1901 as the server words it: what cannot be
// used, then the column whose expression used it. A function is named with its
// parentheses and a column the value came through by name alone, which is how
// the server names each: measured on mariadb:11.8.8, 'rand()' for the call and
// 'a' for the column. Measured on mariadb:11.8.8.
func errUnusableInGeneratedColumn(what, colName string) error {
	return &Error{Code: 1901, SQLState: "HY000",
		Message: fmt.Sprintf("Function or expression '%s' cannot be used in the GENERATED ALWAYS AS clause of `%s`", what, colName)}
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

// validateColumnCheckReferences checks a CHECK written on a column. MariaDB
// lets it name the table's other columns: `CREATE TABLE xr (a INT CHECK (a > b),
// b INT)` is accepted and reads back as `a` int(11) DEFAULT NULL
// CHECK (`a` > `b`). MySQL 8.0 answers ER_CHECK_CONSTRAINT_REFERS (3823).
// Measured on mariadb:11.8.8.
func validateColumnCheckReferences(colName, conName string, expr nodes.ExprNode) error {
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

// isLiteralDefaultForbiddenType reports whether a literal DEFAULT is refused
// for the type. MariaDB allows one on a blob, text or JSON column and stores
// it; only the spatial types refuse, and the refusal is about the value, not
// the type. MySQL 8.0 refuses every one of them with
// ER_BLOB_CANT_HAVE_DEFAULT (1101). Measured on mariadb:11.8.8.
func isLiteralDefaultForbiddenType(typeName string) bool {
	switch strings.ToLower(typeName) {
	case "geometry", "point", "linestring", "polygon",
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
			if isNondeterministicFunc(x.Name) {
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
	return nondeterministicFuncName(name) != ""
}

// nondeterministicFuncName returns the name MariaDB reports for a
// non-deterministic function in its 1901 message, or "" for a deterministic
// one. The server reports the name it folded the call to rather than the one
// written: measured on mariadb:11.8.8, NOW() and CURRENT_TIMESTAMP are both
// 'current_timestamp()', CURRENT_TIME and LOCALTIME() are both 'curtime()',
// CURRENT_DATE is 'curdate()', and LOCALTIMESTAMP() and SYSDATE() keep their
// own names.
func nondeterministicFuncName(name string) string {
	lower := strings.ToLower(name)
	switch lower {
	case "now", "current_timestamp":
		return "current_timestamp"
	case "curtime", "current_time", "localtime":
		return "curtime"
	case "curdate", "current_date":
		return "curdate"
	case "localtimestamp", "sysdate", "utc_timestamp", "utc_time", "utc_date",
		"rand", "uuid", "uuid_short", "connection_id":
		return lower
	}
	return ""
}

// generatedColumnRefs returns the columns a generated expression reads, in the
// order it reads them.
func generatedColumnRefs(expr nodes.ExprNode) []string {
	var out []string
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if cr, ok := n.(*nodes.ColumnRef); ok && cr.Column != "" {
			out = append(out, cr.Column)
		}
		return true
	})
	return out
}

// nondeterministicDependency returns the generated column col takes a
// non-deterministic value from, or "". MariaDB carries nondeterminism through
// the columns an expression reads, so a column that only names another is as
// non-deterministic as the one it names, however many steps away the function
// call is. The column reported is the one named directly, which is what the
// server reports: measured on mariadb:11.8.8,
// `CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) VIRTUAL,
// c DOUBLE AS (b) STORED)` answers 1901 naming b, not a.
func nondeterministicDependency(tbl *Table, col *Column, seen map[string]bool) string {
	if tbl == nil || col == nil || col.Generated == nil {
		return ""
	}
	if seen == nil {
		seen = map[string]bool{}
	}
	if seen[toLower(col.Name)] {
		return ""
	}
	seen[toLower(col.Name)] = true
	for _, ref := range col.Generated.Refs {
		dep := tbl.GetColumn(ref)
		if dep == nil || dep.Generated == nil {
			continue
		}
		if dep.Generated.Nondeterministic != "" ||
			nondeterministicDependency(tbl, dep, seen) != "" {
			// The column's own name, not the text of the reference: a CHANGE
			// renames the column the reference was written against, and the
			// server names it as it now stands. Measured on mariadb:11.8.8,
			// `ALTER TABLE t CHANGE a z DOUBLE AS (RAND()) VIRTUAL` with
			// `b DOUBLE AS (a) STORED` in place answers 1901 naming z.
			return dep.Name
		}
	}
	return ""
}

// validateMaterializedGeneratedColumns refuses a generated column whose value
// the table materializes — a STORED one, or a VIRTUAL one a key covers — and
// which takes that value from a non-deterministic generated column. An
// unmaterialized VIRTUAL column may take it, which is why this asks about the
// table and not about the expression alone.
func validateMaterializedGeneratedColumns(tbl *Table, indexed map[string]bool) error {
	for _, col := range tbl.Columns {
		if col.Generated == nil || col.Generated.RowBound != "" {
			continue
		}
		if !col.Generated.Stored && !indexed[toLower(col.Name)] {
			continue
		}
		if dep := nondeterministicDependency(tbl, col, nil); dep != "" {
			return errUnusableInGeneratedColumn(dep, col.Name)
		}
	}
	return nil
}

// nondeterministicFuncIn returns the name of the first non-deterministic
// function an expression calls, spelled as MariaDB spells it, or "".
func nondeterministicFuncIn(expr nodes.ExprNode) string {
	found := ""
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if found != "" {
			return false
		}
		if fn, ok := n.(*nodes.FuncCallExpr); ok {
			found = nondeterministicFuncName(fn.Name)
		}
		return found == ""
	})
	return found
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

// validateWindowReferences holds one query's window names and references to
// what MariaDB accepts. A reference may extend the window it names, but only
// with what that window leaves out, and the names themselves have to resolve.
//
// Measured on mariadb:11.8.8, where each refusal names the window referenced,
// not the specification doing the referencing:
//
//	4009 a name no WINDOW clause entry before it declares, which makes a
//	     forward reference an undefined one: WINDOW w1 AS (w2), w2 AS (...)
//	4010 the same name declared twice in one WINDOW clause
//	4011 a specification that references another and adds a partition list,
//	     whatever the referenced one holds
//	4012 one that adds an order list to a name that already has one
//	4013 one that references a name that has a frame, even adding nothing:
//	     `OVER (w)` is refused where `OVER w`, a bare reference with no
//	     specification of its own, inherits the frame instead
func validateWindowReferences(stmt *nodes.SelectStmt) error {
	// The clause is read in order, so a reference sees only the entries
	// before it and a name declared twice is caught on its second entry.
	declared := make(map[string]*nodes.WindowDef, len(stmt.WindowClause))
	for _, wd := range stmt.WindowClause {
		if wd == nil {
			continue
		}
		if err := checkWindowRef(wd, declared); err != nil {
			return err
		}
		key := toLower(wd.Name)
		if _, dup := declared[key]; dup {
			return errWindowDuplicateName(wd.Name)
		}
		declared[key] = wd
	}

	// Every OVER clause of this query, wherever it sits in an expression. The
	// walk stops at a nested query, which has its own clause and is reached
	// by the caller's own walk.
	var err error
	check := func(n nodes.Node) bool {
		if err != nil {
			return false
		}
		if nested, ok := n.(*nodes.SelectStmt); ok && nested != stmt {
			return false
		}
		call, ok := n.(*nodes.FuncCallExpr)
		if !ok || call.Over == nil {
			return true
		}
		err = checkWindowRef(call.Over, declared)
		return err == nil
	}
	for _, target := range stmt.TargetList {
		nodes.Inspect(target, check)
	}
	for _, src := range stmt.From {
		nodes.Inspect(src, check)
	}
	if stmt.Where != nil {
		nodes.Inspect(stmt.Where, check)
	}
	for _, item := range stmt.GroupBy {
		nodes.Inspect(item.Expr, check)
	}
	if stmt.Having != nil {
		nodes.Inspect(stmt.Having, check)
	}
	for _, item := range stmt.OrderBy {
		nodes.Inspect(item.Expr, check)
	}
	return err
}

// checkWindowRef holds one window specification to the rules above. declared
// holds the names in scope for it.
func checkWindowRef(wd *nodes.WindowDef, declared map[string]*nodes.WindowDef) error {
	if wd.RefName == "" {
		return nil
	}
	ref, ok := declared[toLower(wd.RefName)]
	if !ok {
		return errWindowNoSuchWindow(wd.RefName)
	}
	if wd.BareRef {
		// `OVER name` adds nothing, so it collides with nothing.
		return nil
	}
	if len(wd.PartitionBy) > 0 {
		return errWindowRefWithPartition(wd.RefName)
	}
	if len(wd.OrderBy) > 0 && len(ref.OrderBy) > 0 {
		return errWindowRefHasOrder(wd.RefName)
	}
	if ref.Frame != nil {
		return errWindowRefHasFrame(wd.RefName)
	}
	return nil
}
