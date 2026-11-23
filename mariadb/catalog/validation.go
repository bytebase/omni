package catalog

import (
	"fmt"
	"math"
	"strconv"
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
		if !ok || fn.Schema != "" || fn.SpacedParen || fn.NameQuoted {
			return true
		}
		name := strings.ToLower(fn.Name)
		switch name {
		case "now", "current_timestamp", "localtime", "localtimestamp", "sysdate",
			"utc_timestamp", "curtime", "current_time", "utc_time":
			if fsp, ok := optTimePrecisionArg(fn); ok && fsp > 6 {
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
	// A ROW START or ROW END column has no expression and holds a value of its
	// own, so a key may cover it: MariaDB 11.8.8 makes `s TIMESTAMP(6) AS ROW
	// START KEY` the primary key.
	if colDef.Generated != nil && colDef.Generated.RowBound == nodes.RowBoundNone {
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

// errUninitializedField is 4029, an expression that reads a column the statement has
// not built yet: the column it defines itself, or a generated column defined after it.
func errUninitializedField(field, uninitialized string) error {
	return &Error{Code: 4029, SQLState: "01000",
		Message: fmt.Sprintf("Expression for field `%s` is referring to uninitialized field `%s`", field, uninitialized)}
}

// columnExprRefs returns the columns of targetDB.targetTable a generated or
// default expression reads. The server reads both the same way here, whichever
// spelling wrote the default: `ADD a INT DEFAULT a`, `DEFAULT(a)` and `AS(a)`
// all answer 4029.
func columnExprRefs(targetDB, targetTable string, colDef *nodes.ColumnDef) []string {
	var refs []string
	add := func(expr nodes.ExprNode) {
		refs = append(refs, tableColumnRefs(targetDB, targetTable, expr)...)
	}
	add(colDef.DefaultValue)
	if colDef.Generated != nil {
		add(colDef.Generated.Expr)
	}
	for _, cc := range colDef.Constraints {
		if cc.Type == nodes.ColConstrDefault {
			add(cc.Expr)
		}
	}
	return refs
}

// validateColumnSelfReference refuses a column whose own generated or default
// expression reads it: 4029, "Expression for field `a` is referring to uninitialized
// field `a`", measured on mariadb:11.8.8 for `CREATE TABLE z (a INT AS(a))`,
// `ALTER TABLE t ADD a INT DEFAULT a` and `ALTER TABLE t ALTER a SET DEFAULT a`. A
// column name the table already holds is 1060 first, and a name the table does not
// hold at all is 1054, so this runs where the column is built.
func validateColumnSelfReference(targetDB, targetTable string, colDef *nodes.ColumnDef) error {
	key := toLower(colDef.Name)
	for _, ref := range columnExprRefs(targetDB, targetTable, colDef) {
		if toLower(ref) == key {
			return errUninitializedField(colDef.Name, colDef.Name)
		}
	}
	return nil
}

// validateColumnGeneratedOrder refuses a column whose generated expression reads a
// generated column defined after it in the same CREATE TABLE: 4029 naming the field
// read, measured on mariadb:11.8.8 for `CREATE TABLE z (a INT AS(b), b INT AS(a))`.
// A plain column defined later is initialized by the time the value is computed, so
// `CREATE TABLE z (a INT AS(b), b INT)` runs, as does `CREATE TABLE z (a INT DEFAULT
// b, b INT)`: the order is judged for a generated expression only.
func validateColumnGeneratedOrder(targetDB, targetTable string, colDef *nodes.ColumnDef, later []*nodes.ColumnDef) error {
	if colDef.Generated == nil {
		return nil
	}
	for _, ref := range tableColumnRefs(targetDB, targetTable, colDef.Generated.Expr) {
		for _, other := range later {
			if other.Generated != nil && toLower(other.Name) == toLower(ref) {
				return errUninitializedField(colDef.Name, other.Name)
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
		if col != nil && col.Generated != nil && !col.Generated.Stored && col.Generated.RowBound == "" {
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

// validateExprColumns refuses an expression that names a column the table does
// not have (1054) or qualifies a column with another table or database (1901).
// MariaDB reads a bare word, or a reference qualified with the statement's own
// table, as a column of the table (`ALTER TABLE a ADD CHECK(a)` and
// `CHECK(a.missing)` answer 1054, `DEFAULT (b + 1)` and `CHECK(a.b > 0)` run),
// and answers 1901 when the qualifier names another table (`CHECK(other.b > 0)`,
// `DEFAULT (other.b)`, `AS (other.b)`). Measured on mariadb:11.8.8.
func validateExprColumns(tbl *Table, targetDB, targetTable string, extraCols map[string]bool, expr nodes.ExprNode, context, owner string) error {
	if expr == nil {
		return nil
	}
	sequences := sequenceArgColumnRefs(expr)
	var exprErr error
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if exprErr != nil {
			return false
		}
		cr, ok := n.(*nodes.ColumnRef)
		if !ok || cr.Column == "" || sequences[cr] {
			return true
		}
		if cr.Table != "" || cr.Schema != "" {
			if (cr.Schema == "" || eqFoldStr(cr.Schema, targetDB)) && eqFoldStr(cr.Table, targetTable) {
				if tbl.GetColumn(cr.Column) == nil && !extraCols[toLower(cr.Column)] {
					exprErr = errNoSuchColumn(cr.Column, context)
				}
				return true
			}
			exprErr = &Error{
				Code:     1901,
				SQLState: "HY000",
				Message:  fmt.Sprintf("Function or expression '%s' cannot be used in the %s clause of `%s`", formatQualifiedColumnRef(cr), context, owner),
			}
			return false
		}
		if tbl.GetColumn(cr.Column) == nil && !extraCols[toLower(cr.Column)] {
			exprErr = errNoSuchColumn(cr.Column, context)
		}
		return true
	})
	return exprErr
}

func formatQualifiedColumnRef(cr *nodes.ColumnRef) string {
	if cr.Schema != "" {
		return cr.Schema + "." + cr.Table + "." + cr.Column
	}
	if cr.Table != "" {
		return cr.Table + "." + cr.Column
	}
	return cr.Column
}

// sequenceArgColumnRefs returns the ColumnRef nodes that name sequences in
// NEXTVAL(s), LASTVAL(s) and SETVAL(s, ...), parenthesized or not, so column
// checks do not treat a sequence's name as a column of the table.
func sequenceArgColumnRefs(expr nodes.ExprNode) map[*nodes.ColumnRef]bool {
	sequences := map[*nodes.ColumnRef]bool{}
	if expr == nil {
		return sequences
	}
	nodes.Inspect(expr, func(n nodes.Node) bool {
		fc, ok := n.(*nodes.FuncCallExpr)
		if !ok || fc.Schema != "" || len(fc.Args) == 0 {
			return true
		}
		switch strings.ToUpper(fc.Name) {
		case "NEXTVAL", "LASTVAL", "SETVAL":
			if cr, ok := fc.Args[0].(*nodes.ColumnRef); ok {
				sequences[cr] = true
			}
			if cr, ok := unwrapParenColumnRef(fc.Args[0]); ok {
				sequences[cr] = true
			}
		}
		return true
	})
	return sequences
}

// unwrapParenColumnRef returns the column reference a parenthesized argument
// holds, as `NEXTVAL((s))` writes it.
func unwrapParenColumnRef(expr nodes.ExprNode) (*nodes.ColumnRef, bool) {
	if p, ok := expr.(*nodes.ParenExpr); ok {
		return unwrapParenColumnRef(p.Expr)
	}
	cr, ok := expr.(*nodes.ColumnRef)
	return cr, ok
}

// validateAlteredColumnExprColumns holds the DEFAULT, GENERATED ALWAYS AS, and
// CHECK expressions an ALTER adds, modifies or changes to the table the
// statement leaves: measured on mariadb:11.8.8,
// `ALTER TABLE t ADD a INT CHECK (a > b), ADD b INT CHECK (b < a)` and
// `ALTER TABLE t ADD CHECK (b > 0), ADD b INT` run over a table with neither
// column, where a name no column of the result has is 1054. A same-ALTER
// rename rewrites only expressions that existed before the statement began:
// `ALTER TABLE a ADD e INT CHECK (c > 0), RENAME COLUMN c TO z` answers 1054
// for c while `ADD e INT CHECK (z > 0), RENAME COLUMN c TO z` runs (recorded in
// TestSimulationCausesOracle).
func validateAlteredColumnExprColumns(tbl *Table, targetDB, targetTable string, cmds []*nodes.AlterTableCmd, skippedColumns map[*nodes.ColumnDef]bool, skippedKeys map[*nodes.Constraint]bool, checkNames map[*nodes.Constraint]string) error {
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		switch cmd.Type {
		case nodes.ATAlterColumnDefault:
			owner := cmd.Name
			if col := tbl.GetColumn(cmd.Name); col != nil {
				owner = col.Name
			}
			if err := validateExprColumns(tbl, targetDB, targetTable, nil, cmd.DefaultExpr, "DEFAULT", owner); err != nil {
				return err
			}
		case nodes.ATAddColumn, nodes.ATModifyColumn, nodes.ATChangeColumn:
			for _, colDef := range append([]*nodes.ColumnDef{cmd.Column}, cmd.Columns...) {
				// A column IF NOT EXISTS finds is a note whose definition is never
				// read.
				if colDef == nil || skippedColumns[colDef] {
					continue
				}
				if err := validateColumnExprColumns(tbl, targetDB, targetTable, nil, colDef); err != nil {
					return err
				}
			}
		}
		for _, con := range append([]*nodes.Constraint{cmd.Constraint}, cmd.Constraints...) {
			if con == nil || con.Type != nodes.ConstrCheck || skippedKeys[con] {
				continue
			}
			owner := checkNames[con]
			if owner == "" {
				owner = con.Name
			}
			if err := validateExprColumns(tbl, targetDB, targetTable, nil, con.Expr, "CHECK", owner); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateColumnExprColumns holds the expressions a column definition carries to
// validateExprColumns: the column's own DEFAULT, its GENERATED ALWAYS AS
// expression, and each CHECK written beside it. A column-level CHECK may name
// the column it is declared on, which the table holds once the column is in.
func validateColumnExprColumns(tbl *Table, targetDB, targetTable string, extraCols map[string]bool, colDef *nodes.ColumnDef) error {
	if colDef == nil {
		return nil
	}
	if err := validateExprColumns(tbl, targetDB, targetTable, extraCols, colDef.DefaultValue, "DEFAULT", colDef.Name); err != nil {
		return err
	}
	if colDef.Generated != nil && colDef.Generated.RowBound == nodes.RowBoundNone {
		if err := validateExprColumns(tbl, targetDB, targetTable, extraCols, colDef.Generated.Expr, "GENERATED ALWAYS AS", colDef.Name); err != nil {
			return err
		}
	}
	for _, cc := range colDef.Constraints {
		switch cc.Type {
		case nodes.ColConstrDefault:
			if err := validateExprColumns(tbl, targetDB, targetTable, extraCols, cc.Expr, "DEFAULT", colDef.Name); err != nil {
				return err
			}
		case nodes.ColConstrCheck:
			owner := cc.Name
			if owner == "" {
				owner = colDef.Name
			}
			if err := validateExprColumns(tbl, targetDB, targetTable, extraCols, cc.Expr, "CHECK", owner); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateEnumSetValues refuses an ENUM or SET column that lists one value
// twice under strict sql_mode: 1291, which MariaDB compares with the column's
// collation and answers with the first trimmed spelling. Outside strict mode
// (`SET sql_mode = ”` or `'ANSI'`) MariaDB keeps the definition with a 1291
// note. Measured on mariadb:11.8.8: `CREATE TABLE t (e ENUM('s','S'))` and
// `SET ('a','A')` answer 1291 under DEFAULT, TRADITIONAL and
// STRICT_TRANS_TABLES, and `ENUM('a ','a')` names `'a'`.
func (c *Catalog) validateEnumSetValues(colDef *nodes.ColumnDef, collation string) error {
	if colDef == nil || colDef.TypeName == nil {
		return nil
	}
	kind := strings.ToUpper(colDef.TypeName.Name)
	if kind != "ENUM" && kind != "SET" {
		return nil
	}
	return c.validateEnumSetList(colDef.Name, kind, colDef.TypeName.EnumValues, collation)
}

func (c *Catalog) validateEnumSetList(colName, kind string, values []string, collation string) error {
	if !c.isStrictSQLMode() {
		return nil
	}
	// MariaDB compares the values with the column's collation, so a
	// case-sensitive one lists `('s','S')` twice over: measured on
	// mariadb:11.8.8, `ENUM('s','S') COLLATE utf8mb4_bin`, `BINARY` and a
	// table's `COLLATE utf8mb4_bin` are all accepted, where the default
	// utf8mb4_uca1400_ai_ci answers 1291.
	fold := !collationCaseSensitive(collation)
	seen := make(map[string]string, len(values))
	for _, value := range values {
		// A trailing space is the collation's: MariaDB stores `ENUM('a ','b ')`
		// as enum('a','b') and answers 1291 for `('a','a ')`, whatever the
		// collation, naming the trimmed value ('a'). Measured on mariadb:11.8.8,
		// a tab being kept.
		trimmed := strings.TrimRight(value, " ")
		key := trimmed
		if fold {
			key = foldEnumValue(key, collation)
		}
		if first, ok := seen[key]; ok {
			// The message names the first spelling: `ENUM('s','S')` over a
			// case-insensitive collation reports "duplicated value 's'".
			return errDupEnumValue(colName, first, kind)
		}
		seen[key] = trimmed
	}
	return nil
}

// isStrictSQLMode reports whether session sql_mode promotes warnings such as
// 1291 to errors: STRICT_TRANS_TABLES (part of DEFAULT on MariaDB 11.8.8),
// STRICT_ALL_TABLES, or TRADITIONAL.
func (c *Catalog) isStrictSQLMode() bool {
	for _, mode := range strings.Split(strings.ToUpper(c.session.SQLMode), ",") {
		switch strings.TrimSpace(mode) {
		case "DEFAULT", "TRADITIONAL", "STRICT_TRANS_TABLES", "STRICT_ALL_TABLES":
			return true
		}
	}
	return false
}

// parseColumnEnumSetValues reads the member list back from a stored ColumnType
// (`enum('a','b”c')`), where escapeEnumValue doubled single quotes.
func parseColumnEnumSetValues(columnType string) []string {
	open := strings.IndexByte(columnType, '(')
	close := strings.LastIndexByte(columnType, ')')
	if open < 0 || close <= open {
		return nil
	}
	s := columnType[open+1 : close]
	var values []string
	for i := 0; i < len(s); {
		if s[i] != '\'' {
			i++
			continue
		}
		i++ // skip opening '\''
		var b strings.Builder
		for i < len(s) {
			if s[i] == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					b.WriteByte('\'')
					i += 2
					continue
				}
				i++ // skip closing '\''
				break
			}
			b.WriteByte(s[i])
			i++
		}
		values = append(values, b.String())
	}
	return values
}

// effectiveColumnCollation is the collation a column definition's ENUM or SET
// values are compared with: the one it names (DEFAULT resolved to the column
// charset's default), the binary one BINARY or CHARACTER SET binary gives it,
// the default of the CHARACTER SET it names, and otherwise the table's. A named
// character set drops the table's collation even where they share one: measured
// on mariadb:11.8.8, `ENUM('s','S') CHARACTER SET utf8mb4` and `COLLATE
// DEFAULT` under a table `COLLATE utf8mb4_bin` both answer 1291.
// declaredColumnCollate returns the collation a column definition writes, on
// its type (`VARCHAR(10) COLLATE ...`) or among its column attributes
// (`VARCHAR(10) NOT NULL COLLATE ...`).
func declaredColumnCollate(colDef *nodes.ColumnDef) string {
	if colDef == nil {
		return ""
	}
	if colDef.TypeName != nil && colDef.TypeName.Collate != "" {
		return colDef.TypeName.Collate
	}
	collate := ""
	for _, cc := range colDef.Constraints {
		if cc == nil || cc.Type != nodes.ColConstrCollate {
			continue
		}
		if cc.Name != "" {
			collate = cc.Name
		} else if s, ok := cc.Expr.(*nodes.StringLit); ok && s.Value != "" {
			collate = s.Value
		}
	}
	return collate
}

func (c *Catalog) effectiveColumnCollation(db *Database, colDef *nodes.ColumnDef, tableCharset, tableCollation string) string {
	dt := colDef.TypeName
	if dt == nil {
		return tableCollation
	}
	collate := declaredColumnCollate(colDef)
	colCharset := ""
	if isNationalStringType(dt.Name) {
		colCharset = "utf8mb3"
	}
	if dt.Charset != "" {
		colCharset = normalizeCharsetName(dt.Charset)
	}
	if collate != "" && !eqFoldStr(collate, "DEFAULT") && colCharset == "" {
		colCharset = normalizeCharsetName(charsetForCollation(collate))
	}
	if colCharset == "" {
		colCharset = tableCharset
	}
	if collate != "" {
		if eqFoldStr(collate, "DEFAULT") {
			if dt.Binary || eqFoldStr(colCharset, "binary") {
				return "binary"
			}
			if db != nil {
				return versionCollation(db, colCharset)
			}
			return c.charsetCollation(colCharset)
		}
		return collate
	}
	if dt.Binary || eqFoldStr(colCharset, "binary") {
		return "binary"
	}
	if dt.Charset != "" || isNationalStringType(dt.Name) {
		if db != nil {
			if dc := versionCollation(db, colCharset); dc != "" {
				return dc
			}
		} else if dc := c.charsetCollation(colCharset); dc != "" {
			return dc
		}
	}
	return tableCollation
}

// foldEnumValue folds one value the way a case-insensitive collation reads it.
// Case is all it folds: a collation's weights also make `a` and `á` one value
// under an accent-insensitive one, which needs the collation's table and is not
// modeled here.
// Turkish is its own case mapping: "I" and "i" are different letters there, the
// pairs being I with ı (dotless) and İ with i (dotted). Measured on
// mariadb:11.8.8 under utf8mb4_turkish_ci: `ENUM('I','i')` runs, `('i','İ')` and
// `('s','S')` answer 1291.
func foldEnumValue(value, collation string) string {
	if strings.Contains(strings.ToLower(collation), "turkish") {
		value = strings.ReplaceAll(value, "I", "ı")
		value = strings.ReplaceAll(value, "İ", "i")
	}
	return strings.ToLower(value)
}

// collationCaseSensitive reports whether a collation tells two spellings apart
// by case: MariaDB's binary ones name it, `binary`, `_bin` and `_cs`.
func collationCaseSensitive(collation string) bool {
	name := strings.ToLower(collation)
	return name == "binary" || strings.HasSuffix(name, "_bin") || strings.HasSuffix(name, "_cs")
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

// validateCheckExpr refuses what a CHECK cannot read, as mariadb:11.8.8
// answers it: a subquery is 1970, and a variable or a function whose value
// changes between rows 1901, naming the expression and the constraint (the
// column, for a CHECK written on one): CHECK (a > RAND()) is "Function or
// expression 'rand()' cannot be used in the CHECK clause of `c1`", NOW() is
// 'current_timestamp()' and @x is '@x'.
func validateCheckExpr(name string, expr nodes.ExprNode) error {
	exprSQL := strings.ToLower(nodeToSQL(expr))
	if strings.Contains(exprSQL, "select ") || strings.Contains(exprSQL, "(select") {
		return errCheckSubquery()
	}
	var refused error
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if refused != nil {
			return false
		}
		switch x := n.(type) {
		case *nodes.SubqueryExpr:
			refused = errCheckSubquery()
		case *nodes.VariableRef:
			text := "@" + x.Name
			if x.System {
				text = "@@" + x.Name
			}
			refused = errUnusableInCheck(text, name)
		case *nodes.FuncCallExpr:
			if fn := nondeterministicFuncName(x.Name); fn != "" {
				refused = errUnusableInCheck(fn+"()", name)
			}
		}
		return refused == nil
	})
	return refused
}

func errCheckSubquery() error {
	return &Error{Code: 1970, SQLState: "42000", Message: "CHECK does not support subqueries or stored functions"}
}

func errUnusableInCheck(expr, name string) error {
	return &Error{Code: 1901, SQLState: "HY000", Message: fmt.Sprintf("Function or expression '%s' cannot be used in the CHECK clause of `%s`", expr, name)}
}

// checkAutoIncrementInChecks refuses a CHECK a CREATE TABLE writes that reads
// an AUTO_INCREMENT column, whose value is not known when the row is checked:
// measured on mariadb:11.8.8, CREATE TABLE (a INT AUTO_INCREMENT PRIMARY KEY,
// CHECK (a > 0)) answers 1901 with 'AUTO_INCREMENT' and the column's name, and
// so does a CHECK reading it beside another column. It reads the statement's
// own expressions, once every column is in the table.
func checkAutoIncrementInChecks(tbl *Table, stmt *nodes.CreateTableStmt) error {
	for _, def := range stmt.Columns {
		for _, cc := range def.Constraints {
			if cc != nil && cc.Type == nodes.ColConstrCheck {
				if err := checkAutoIncrementInCheckExpr(tbl, cc.Expr); err != nil {
					return err
				}
			}
		}
	}
	for _, con := range stmt.Constraints {
		if con != nil && con.Type == nodes.ConstrCheck {
			if err := checkAutoIncrementInCheckExpr(tbl, con.Expr); err != nil {
				return err
			}
		}
	}
	return nil
}

// checkAutoIncrementInTableChecks holds the table an ALTER TABLE leaves to
// checkAutoIncrementInChecks' rule, whichever clause brought the CHECK or the
// AUTO_INCREMENT in: measured on mariadb:11.8.8 over (id INT AUTO_INCREMENT
// PRIMARY KEY, b INT, CHECK (b > 0)), MODIFY id INT AUTO_INCREMENT CHECK (id >
// 0), ADD COLUMN x INT CHECK (x > id) and MODIFY id INT, MODIFY b INT
// AUTO_INCREMENT UNIQUE each answer 1901. A CHECK's columns are the ones it was
// declared reading (Constraint.checkRefs); one whose columns the catalog did
// not read is passed over.
func checkAutoIncrementInTableChecks(tbl *Table) error {
	var auto []*Column
	for _, col := range tbl.Columns {
		if col.AutoIncrement {
			auto = append(auto, col)
		}
	}
	if len(auto) == 0 {
		return nil
	}
	for _, con := range tbl.Constraints {
		if con.Type != ConCheck || con.Implicit {
			continue
		}
		for _, col := range auto {
			if con.checkRefs[toLower(col.Name)] {
				return errUnusableInCheck("AUTO_INCREMENT", col.Name)
			}
		}
	}
	return nil
}

// checkAutoIncrementInCheckExpr refuses one CHECK expression that reads an
// AUTO_INCREMENT column of tbl (checkAutoIncrementInChecks).
func checkAutoIncrementInCheckExpr(tbl *Table, expr nodes.ExprNode) error {
	for _, name := range generatedColumnRefs(expr) {
		if col := tbl.GetColumn(name); col != nil && col.AutoIncrement {
			return errUnusableInCheck("AUTO_INCREMENT", col.Name)
		}
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
// order it reads them, skipping sequence names in NEXTVAL/LASTVAL/SETVAL.
func generatedColumnRefs(expr nodes.ExprNode) []string {
	sequences := sequenceArgColumnRefs(expr)
	var out []string
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if cr, ok := n.(*nodes.ColumnRef); ok && cr.Column != "" && !sequences[cr] {
			out = append(out, cr.Column)
		}
		return true
	})
	return out
}

// tableColumnRefs returns the columns of targetDB.targetTable an expression
// reads, skipping sequence arguments and references qualified with another
// table or database.
func tableColumnRefs(targetDB, targetTable string, expr nodes.ExprNode) []string {
	sequences := sequenceArgColumnRefs(expr)
	var out []string
	nodes.Inspect(expr, func(n nodes.Node) bool {
		cr, ok := n.(*nodes.ColumnRef)
		if !ok || cr.Column == "" || sequences[cr] {
			return true
		}
		if cr.Table != "" || cr.Schema != "" {
			if targetTable == "" || !eqFoldStr(cr.Table, targetTable) || (cr.Schema != "" && !eqFoldStr(cr.Schema, targetDB)) {
				return true
			}
		}
		out = append(out, cr.Column)
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
func nondeterministicDependency(tbl *Table, col *Column, carries map[string]bool) string {
	if tbl == nil || col == nil || col.Generated == nil {
		return ""
	}
	if carries == nil {
		carries = map[string]bool{}
	}
	for _, ref := range col.Generated.Refs {
		dep := tbl.GetColumn(ref)
		if dep == nil || dep.Generated == nil {
			continue
		}
		if carriesNondeterminism(tbl, dep, carries, nil) {
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

// carriesNondeterminism reports whether a generated column's value is
// non-deterministic, by its own expression or through the columns it reads.
// Each answer is kept in carries for as long as the caller holds it, which is
// one check of one table: a chain of columns each reading the one before is
// then walked once rather than once per column that depends on it, where a
// table of a few thousand took seconds of the latter.
//
// path holds the columns the walk is inside, which stops it going round a
// cycle. Nothing is lost by that: the edge it cuts leads back to a column
// already being walked, and a non-deterministic column anywhere below is
// still reached by the edges that are not cut.
func carriesNondeterminism(tbl *Table, col *Column, carries, path map[string]bool) bool {
	name := toLower(col.Name)
	if answer, ok := carries[name]; ok {
		return answer
	}
	if path[name] {
		return false
	}
	if col.Generated.Nondeterministic != "" {
		carries[name] = true
		return true
	}
	if path == nil {
		path = map[string]bool{}
	}
	path[name] = true
	defer delete(path, name)
	for _, ref := range col.Generated.Refs {
		dep := tbl.GetColumn(ref)
		if dep == nil || dep.Generated == nil {
			continue
		}
		if carriesNondeterminism(tbl, dep, carries, path) {
			carries[name] = true
			return true
		}
	}
	carries[name] = false
	return false
}

// validateMaterializedGeneratedColumns refuses a generated column whose value
// the table materializes — a STORED one, or a VIRTUAL one a key covers — and
// which takes that value from a non-deterministic generated column. An
// unmaterialized VIRTUAL column may take it, which is why this asks about the
// table and not about the expression alone.
func validateMaterializedGeneratedColumns(tbl *Table, indexed map[string]bool) error {
	// One map for the whole table, so a column's answer is reached once
	// however many columns depend on it.
	carries := map[string]bool{}
	for _, col := range tbl.Columns {
		if col.Generated == nil || col.Generated.RowBound != "" {
			continue
		}
		if !col.Generated.Stored && !indexed[toLower(col.Name)] {
			continue
		}
		if dep := nondeterministicDependency(tbl, col, carries); dep != "" {
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

// optTimePrecisionArg returns the effective fractional seconds precision that
// MariaDB 11.8.8 computes for temporal function opt_time_precision arguments
// (CURTIME/CURRENT_TIME/NOW/CURRENT_TIMESTAMP/LOCALTIME/LOCALTIMESTAMP/SYSDATE/
// UTC_TIME/UTC_TIMESTAMP). MariaDB converts real_ulong_num through uint16
// truncation for decimal integer literals (so 65536, 2147483648, and
// 9223372036854775808 wrap to 0, while 2147483649 wraps to 1 and 65535 stays
// 65535 -> 1426). For 0x hex literals, values exceeding MaxInt64 saturate at
// MaxInt64 whose low 16 bits are 65535 -> 1426.
func optTimePrecisionArg(fn *nodes.FuncCallExpr) (int, bool) {
	if fn == nil || len(fn.Args) == 0 {
		return 0, false
	}
	switch lit := fn.Args[0].(type) {
	case *nodes.IntLit:
		if lit.Text != "" {
			if v, err := strconv.ParseUint(lit.Text, 10, 64); err == nil {
				return int(uint16(v)), true
			}
			return 0xFFFF, true
		}
		if lit.Value >= 0 {
			return int(uint16(lit.Value)), true
		}
		return int(lit.Value), true
	case *nodes.HexLit:
		hexStr := lit.Value
		if strings.HasPrefix(hexStr, "0x") || strings.HasPrefix(hexStr, "0X") {
			hexStr = hexStr[2:]
		}
		v, err := strconv.ParseUint(hexStr, 16, 64)
		if err != nil || v > math.MaxInt64 {
			return 0xFFFF, true
		}
		return int(uint16(v)), true
	}
	return 0, false
}

func temporalFuncFSP(fn *nodes.FuncCallExpr) (int, bool) {
	if fn == nil || fn.Schema != "" || fn.SpacedParen || fn.NameQuoted {
		return 0, false
	}
	name := strings.ToLower(fn.Name)
	switch name {
	case "now", "current_timestamp", "localtime", "localtimestamp", "sysdate",
		"utc_timestamp":
		if fsp, ok := optTimePrecisionArg(fn); ok {
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
		if e := checkWindowRef(call.Over, declared); e != nil {
			err = e
			return false
		}
		if call.Schema == "" && !call.SpacedParen && windowFrameDisallowedFunc(call.Name) && windowDefHasFrame(call.Over, declared) {
			err = errWindowFrameNotAllowed(strings.ToLower(call.Name))
			return false
		}
		return true
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

// windowFrameDisallowedFunc reports window functions that MariaDB 11.8 rejects
// with 4017 ("Window frame is not allowed with '<name>'") when paired with any
// window frame clause.
func windowFrameDisallowedFunc(name string) bool {
	switch strings.ToLower(name) {
	case "row_number", "rank", "dense_rank", "percent_rank", "cume_dist",
		"ntile", "lag", "lead", "median", "percentile_cont", "percentile_disc":
		return true
	}
	return false
}

// windowDefHasFrame reports whether wd or any named window it references
// specifies a window frame clause.
func windowDefHasFrame(wd *nodes.WindowDef, declared map[string]*nodes.WindowDef) bool {
	seen := map[string]bool{}
	for cur := wd; cur != nil; {
		if cur.Frame != nil {
			return true
		}
		if cur.RefName == "" {
			return false
		}
		key := toLower(cur.RefName)
		if seen[key] {
			return false
		}
		seen[key] = true
		cur = declared[key]
	}
	return false
}

// checkExprRefs returns the lowered names of the columns a CHECK expression
// reads.
func checkExprRefs(expr nodes.ExprNode) map[string]bool {
	refs := map[string]bool{}
	for _, name := range generatedColumnRefs(expr) {
		refs[toLower(name)] = true
	}
	return refs
}
