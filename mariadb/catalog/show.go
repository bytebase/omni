package catalog

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// charsetForCollation derives the charset name from a collation name.
// MySQL collation names are prefixed with the charset name (e.g. latin1_swedish_ci → latin1).
func charsetForCollation(collation string) string {
	collation = toLower(collation)
	// Try known charsets from longest to shortest to handle prefixes like utf8mb4 vs utf8.
	knownCharsets := []string{
		"utf8mb4", "utf8mb3", "utf16le", "geostd8", "armscii8", "macroman",
		"eucjpms", "keybcs2", "gb2312", "euckr", "utf16", "utf32", "macce",
		"latin1", "latin2", "latin5", "latin7", "cp1250", "cp1251", "cp1256",
		"cp1257", "cp850", "cp852", "cp866", "cp932", "koi8r", "koi8u",
		"ascii", "binary", "tis620", "hebrew", "greek", "sjis", "swe7",
		"big5", "ucs2", "ujis", "dec8", "hp8", "utf8", "gbk",
	}
	for _, cs := range knownCharsets {
		if strings.HasPrefix(collation, cs+"_") || collation == cs {
			return cs
		}
	}
	// Fallback: take prefix before first underscore.
	if idx := strings.IndexByte(collation, '_'); idx > 0 {
		return collation[:idx]
	}
	return ""
}

// defaultCollationForCharset returns the default collation for each charset
// MariaDB has. Measured on mariadb:11.8.8 through
// information_schema.CHARACTER_SETS: the Unicode charsets default to the
// UCA 14.0.0 collations, where MySQL 8.0 defaults utf8mb4 to
// utf8mb4_0900_ai_ci and the rest to <charset>_general_ci. MariaDB has no
// gb18030.
var defaultCollationForCharset = map[string]string{
	"utf8mb4":  "utf8mb4_uca1400_ai_ci",
	"utf8mb3":  "utf8mb3_uca1400_ai_ci",
	"utf8":     "utf8mb3_uca1400_ai_ci",
	"ucs2":     "ucs2_uca1400_ai_ci",
	"utf16":    "utf16_uca1400_ai_ci",
	"utf32":    "utf32_uca1400_ai_ci",
	"utf16le":  "utf16le_general_ci",
	"latin1":   "latin1_swedish_ci",
	"latin2":   "latin2_general_ci",
	"latin5":   "latin5_turkish_ci",
	"latin7":   "latin7_general_ci",
	"ascii":    "ascii_general_ci",
	"binary":   "binary",
	"gbk":      "gbk_chinese_ci",
	"big5":     "big5_chinese_ci",
	"euckr":    "euckr_korean_ci",
	"gb2312":   "gb2312_chinese_ci",
	"sjis":     "sjis_japanese_ci",
	"ujis":     "ujis_japanese_ci",
	"cp932":    "cp932_japanese_ci",
	"eucjpms":  "eucjpms_japanese_ci",
	"cp1250":   "cp1250_general_ci",
	"cp1251":   "cp1251_general_ci",
	"cp1256":   "cp1256_general_ci",
	"cp1257":   "cp1257_general_ci",
	"cp850":    "cp850_general_ci",
	"cp852":    "cp852_general_ci",
	"cp866":    "cp866_general_ci",
	"dec8":     "dec8_swedish_ci",
	"hp8":      "hp8_english_ci",
	"keybcs2":  "keybcs2_general_ci",
	"koi8r":    "koi8r_general_ci",
	"koi8u":    "koi8u_general_ci",
	"macce":    "macce_general_ci",
	"macroman": "macroman_general_ci",
	"swe7":     "swe7_swedish_ci",
	"geostd8":  "geostd8_general_ci",
	"tis620":   "tis620_thai_ci",
	"hebrew":   "hebrew_general_ci",
	"greek":    "greek_general_ci",
	"armscii8": "armscii8_general_ci",
}

// ShowCreateTable produces MariaDB-compatible SHOW CREATE TABLE output.
// Returns "" if the database or table does not exist.
func (c *Catalog) ShowCreateTable(database, table string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	tbl := db.GetTable(table)
	if tbl == nil {
		return ""
	}

	var b strings.Builder
	if tbl.Temporary {
		b.WriteString(fmt.Sprintf("CREATE TEMPORARY TABLE `%s` (\n", tbl.Name))
	} else {
		b.WriteString(fmt.Sprintf("CREATE TABLE `%s` (\n", tbl.Name))
	}

	// Columns.
	parts := make([]string, 0, len(tbl.Columns)+len(tbl.Indexes)+len(tbl.Constraints))
	for _, col := range tbl.Columns {
		if col.Hidden == ColumnHiddenSystem {
			continue
		}
		parts = append(parts, showColumnWithTable(col, tbl))
	}

	// Application-time period — rendered before SYSTEM_TIME; the name is
	// backtick-quoted (unlike the SYSTEM_TIME keyword).
	if tbl.AppPeriodName != "" {
		parts = append(parts, fmt.Sprintf("PERIOD FOR `%s` (`%s`, `%s`)", tbl.AppPeriodName, tbl.AppPeriodStartCol, tbl.AppPeriodEndCol))
	}

	// PERIOD FOR SYSTEM_TIME — after columns, before keys/constraints.
	if tbl.PeriodStartCol != "" {
		parts = append(parts, fmt.Sprintf("PERIOD FOR SYSTEM_TIME (`%s`, `%s`)", tbl.PeriodStartCol, tbl.PeriodEndCol))
	}

	// Indexes — MySQL 8.0 orders them in groups:
	// 1. PRIMARY KEY
	// 2. UNIQUE KEYs (creation order)
	// 3. Regular + SPATIAL KEYs, non-expression (creation order)
	// 4. Expression-based KEYs (creation order)
	// 5. FULLTEXT KEYs (creation order)
	var idxPrimary, idxUnique, idxRegular, idxExpr, idxFulltext, idxVector []*Index
	for _, idx := range tbl.Indexes {
		switch {
		case idx.Primary:
			idxPrimary = append(idxPrimary, idx)
		case idx.Unique:
			idxUnique = append(idxUnique, idx)
		case idx.Fulltext:
			idxFulltext = append(idxFulltext, idx)
		case idx.Vector:
			// After every other key, as MariaDB 11.8.8 prints them.
			idxVector = append(idxVector, idx)
		case isExpressionIndex(idx):
			idxExpr = append(idxExpr, idx)
		default:
			idxRegular = append(idxRegular, idx)
		}
	}
	for _, group := range [][]*Index{idxPrimary, idxUnique, idxRegular, idxExpr, idxFulltext, idxVector} {
		for _, idx := range group {
			parts = append(parts, showIndex(idx))
		}
	}

	// Constraints (FK and CHECK only — PK/UNIQUE are shown via indexes).
	// A CHECK that belongs to a column was written on that column, so it is
	// left out here. The question is asked of this very constraint and not of
	// its name: a table-level check carrying a column's name is still the
	// table's, and two checks of one name are written one in each place.
	var fkConstraints, chkConstraints []*Constraint
	for _, con := range tbl.Constraints {
		switch con.Type {
		case ConForeignKey:
			fkConstraints = append(fkConstraints, con)
		case ConCheck:
			if !writtenOnColumn(tbl, con) {
				chkConstraints = append(chkConstraints, con)
			}
		}
	}
	sort.Slice(fkConstraints, func(i, j int) bool {
		return fkConstraints[i].Name < fkConstraints[j].Name
	})
	// The table-level checks keep the order they were declared in, which is
	// what the server writes: ck1 above reports CONSTRAINT_1, named,
	// CONSTRAINT_2. Measured on mariadb:11.8.8.
	for _, con := range fkConstraints {
		parts = append(parts, showConstraint(con))
	}
	for _, con := range chkConstraints {
		parts = append(parts, showConstraint(con))
	}

	b.WriteString("  ")
	b.WriteString(strings.Join(parts, ",\n  "))
	b.WriteString("\n)")

	// Table options.
	opts := showTableOptions(tbl)
	if opts != "" {
		b.WriteString(" ")
		b.WriteString(opts)
	}

	// System versioning suffix.
	if tbl.SystemVersioned {
		b.WriteString(" WITH SYSTEM VERSIONING")
	}

	// Partition clause.
	if tbl.Partitioning != nil {
		b.WriteString("\n")
		b.WriteString(showPartitioning(tbl.Partitioning))
	}

	return b.String()
}

func showColumn(col *Column) string {
	return showColumnWithTable(col, nil)
}

func showColumnWithTable(col *Column, tbl *Table) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("`%s` %s", col.Name, col.ColumnType))

	// CHARACTER SET and COLLATE. MariaDB writes them together, and only where
	// the column's collation differs from the table's: a column that inherits
	// the table's collation reports neither, whatever the charset's own default
	// is, and one that does not reports both. MySQL 8.0 instead compares
	// against the charset's default collation. Measured on mariadb:11.8.8.
	if isStringType(col.DataType) || isEnumSetType(col.DataType) {
		tableCharset := ""
		if tbl != nil {
			tableCharset = tbl.Charset
		}
		// Resolve the default collation for the column's charset.
		colCharsetDefault := ""
		if col.Charset != "" {
			if dc, ok := defaultCollationForCharset[toLower(col.Charset)]; ok {
				colCharsetDefault = dc
			}
		}
		charsetDiffers := col.Charset != "" && !eqFoldStr(col.Charset, tableCharset)

		// The table's effective collation, which the column inherits when it
		// names none of its own.
		tableCollation := ""
		if tbl != nil {
			tableCollation = tbl.Collation
			if tableCollation == "" {
				tableCollation = defaultCollationForCharset[toLower(tbl.Charset)]
			}
		}
		collation := col.Collation
		if collation == "" {
			collation = colCharsetDefault
		}
		collationDiffers := collation != "" && !eqFoldStr(collation, tableCollation)

		if charsetDiffers || collationDiffers {
			if col.Charset != "" {
				b.WriteString(fmt.Sprintf(" CHARACTER SET %s", col.Charset))
			}
			if collation != "" {
				b.WriteString(fmt.Sprintf(" COLLATE %s", collation))
			}
		}
	}

	// Generated column.
	if col.Generated != nil {
		if col.Generated.RowBound != "" {
			// System-versioning period column: GENERATED ALWAYS AS ROW START/END
			// (no expression, no VIRTUAL/STORED mode).
			b.WriteString(fmt.Sprintf(" GENERATED ALWAYS AS %s", col.Generated.RowBound))
			if col.Comment != "" {
				b.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeComment(col.Comment)))
			}
			if col.Invisible {
				b.WriteString(" INVISIBLE")
			}
			return b.String()
		}
		mode := "VIRTUAL"
		if col.Generated.Stored {
			mode = "STORED"
		}
		// MariaDB writes the expression inside one pair of parentheses:
		// `(`a` + `b`)` for a + b, and `((`a` + `b`) * `c`)` for (a+b)*c, where
		// the inner pair is the one the precedence needs. The stored form
		// carries a pair of its own around a top-level operator, so one level
		// comes off here. Measured on mariadb:11.8.8.
		b.WriteString(fmt.Sprintf(" GENERATED ALWAYS AS (%s) %s", stripOuterParens(col.Generated.Expr), mode))
		if !col.Nullable {
			b.WriteString(" NOT NULL")
		}
		if col.Comment != "" {
			b.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeComment(col.Comment)))
		}
		if col.Invisible {
			b.WriteString(" INVISIBLE")
		}
		return b.String()
	}

	// NOT NULL / NULL.
	if !col.Nullable {
		b.WriteString(" NOT NULL")
	} else if isTimestampType(col.DataType) {
		// MySQL 8.0 explicitly shows NULL for TIMESTAMP columns.
		b.WriteString(" NULL")
	}

	// SRID for spatial types — MySQL 8.0 places it after NOT NULL, before DEFAULT.
	if col.SRID != 0 {
		b.WriteString(fmt.Sprintf(" /*!80003 SRID %d */", col.SRID))
	}

	// DEFAULT.
	if col.Default != nil {
		b.WriteString(" DEFAULT ")
		b.WriteString(formatDefault(*col.Default, col))
	} else if col.Nullable && !col.AutoIncrement {
		// MariaDB writes DEFAULT NULL for every nullable column — the text,
		// blob, JSON and spatial types included, and one whose default was
		// dropped — where MySQL 8.0 leaves it off the types that cannot carry
		// a literal default and off a column after DROP DEFAULT. Measured on
		// mariadb:11.8.8 through SHOW CREATE TABLE.
		b.WriteString(" DEFAULT NULL")
	}

	// AUTO_INCREMENT.
	if col.AutoIncrement {
		b.WriteString(" AUTO_INCREMENT")
	}

	// ON UPDATE.
	if col.OnUpdate != "" {
		b.WriteString(fmt.Sprintf(" ON UPDATE %s", formatOnUpdate(col.OnUpdate, col.OnUpdateKind)))
	}

	// COMMENT.
	if col.Comment != "" {
		b.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeComment(col.Comment)))
	}

	// INVISIBLE. MariaDB writes the keyword bare; MySQL 8.0 hides it in a
	// version comment. Measured on mariadb:11.8.8.
	if col.Invisible {
		b.WriteString(" INVISIBLE")
	}

	// Column-level WITHOUT SYSTEM VERSIONING (WITH is the implicit default).
	if col.SystemVersioning == "WITHOUT" {
		b.WriteString(" WITHOUT SYSTEM VERSIONING")
	}

	// The column's own CHECK, which MariaDB writes on the column and names
	// after it, last in the definition and after any COMMENT:
	// `a` int(11) DEFAULT 5 COMMENT 'hi' CHECK (`a` > 0). MySQL 8.0 writes
	// every check as a table-level constraint. Measured on mariadb:11.8.8.
	if tbl != nil {
		if con := columnCheckOf(tbl, col.Name); con != nil {
			b.WriteString(fmt.Sprintf(" CHECK (%s)", con.CheckExpr))
		}
	}

	return b.String()
}

// columnCheckOf returns the CHECK constraint that belongs to the named column
// and is written on it, or nil.
//
// The answer comes from each check's own record of the column it was declared
// on, since a name match is a different question: a check declared at table
// level is written at table level even when it carries a column's name.
// Measured on mariadb:11.8.8 —
// `CREATE TABLE t (k INT, a INT, CONSTRAINT k CHECK (a > 0))` writes
// CONSTRAINT `k` at table level, beside an untouched column `k`.
//
// A generated column does not carry one: MariaDB writes
// `CREATE TABLE c8 (a INT, b INT AS (a+1), CONSTRAINT b CHECK (a>0))` back with
// the check as a table-level CONSTRAINT `b`, not on `b`. Measured on
// mariadb:11.8.8. Treating it as the column's would drop it from the output
// altogether, since the generated-column branch writes no check.
func columnCheckOf(tbl *Table, name string) *Constraint {
	if tbl == nil || name == "" {
		return nil
	}
	var owner *Column
	for _, col := range tbl.Columns {
		if eqFoldStr(col.Name, name) {
			owner = col
			break
		}
	}
	if owner == nil || owner.Generated != nil {
		return nil
	}
	for _, con := range tbl.Constraints {
		if isColumnCheck(con) && eqFoldStr(con.Column, name) {
			return con
		}
	}
	return nil
}

// writtenOnColumn reports whether this check is the one SHOW CREATE TABLE
// writes on a column of the table, which is where it then appears instead of
// among the table's constraints.
func writtenOnColumn(tbl *Table, con *Constraint) bool {
	return isColumnCheck(con) && columnCheckOf(tbl, con.Column) == con
}

// formatExpressionDefault writes an expression default the way MariaDB stores
// it: bare when it is a single function call, and otherwise inside the
// parentheses the DEFAULT grammar requires for an operator expression.
// Measured on mariadb:11.8.8 — see formatDefault.
func formatExpressionDefault(val string) string {
	inner := stripOuterParens(val)
	if isBareFunctionCall(inner) {
		return inner
	}
	return "(" + inner + ")"
}

// scanOutsideQuotes calls f for each byte of expr that is not inside a quoted
// string or identifier, and reports whether every quote was closed. MariaDB
// quotes strings with '\” or '"' and identifiers with '`'; a quote inside one
// is written doubled, or, in a string, escaped with a backslash.
func scanOutsideQuotes(expr string, f func(i int, c byte)) bool {
	var quote byte
	for i := 0; i < len(expr); i++ {
		c := expr[i]
		if quote == 0 {
			switch c {
			case '\'', '"', '`':
				quote = c
			default:
				f(i, c)
			}
			continue
		}
		switch {
		case c == '\\' && quote != '`' && i+1 < len(expr):
			i++ // an escaped byte, whatever it is
		case c == quote && i+1 < len(expr) && expr[i+1] == quote:
			i++ // a doubled quote stands for one
		case c == quote:
			quote = 0
		}
	}
	return quote == 0
}

// isNameByte reports whether c can appear in an unquoted function name.
func isNameByte(c byte) bool {
	return c == '_' || c == '$' || c == '.' ||
		(c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}

// isBareFunctionCall reports whether expr is one function call and nothing
// else: a name, then a balanced argument list that ends the text.
//
// A parenthesis inside a string or a quoted identifier is text rather than
// structure, so it is not counted: `length('(') + length(')')` is two calls
// either side of an operator and keeps its parentheses, while `length('()')`
// is the one call and loses them. Measured on mariadb:11.8.8.
func isBareFunctionCall(expr string) bool {
	nameEnd := -1
	nameOK := true
	type paren struct {
		at   int
		open bool
	}
	var parens []paren
	balanced := scanOutsideQuotes(expr, func(i int, c byte) {
		switch c {
		case '(':
			if nameEnd < 0 {
				nameEnd = i
			}
			parens = append(parens, paren{i, true})
		case ')':
			parens = append(parens, paren{i, false})
		default:
			if nameEnd < 0 && !isNameByte(c) {
				nameOK = false
			}
		}
	})
	if !balanced || !nameOK || nameEnd <= 0 || len(parens) == 0 {
		return false
	}
	// The name's own argument list has to be the one that closes the text, so
	// `f(1) + 1` and `f(1) + g(2)` are not one call.
	depth := 0
	for _, p := range parens {
		if p.open {
			depth++
			continue
		}
		depth--
		if depth == 0 && p.at != len(expr)-1 {
			return false
		}
	}
	return depth == 0 && parens[len(parens)-1].at == len(expr)-1
}

// stripOuterParens removes one balanced outer pair of parentheses, and leaves
// the text alone when the pair does not wrap the whole expression.
func stripOuterParens(expr string) string {
	if len(expr) < 2 || expr[0] != '(' || expr[len(expr)-1] != ')' {
		return expr
	}
	// Quoted parentheses are text, so `('(')` is wrapped while `'('` is not.
	depth, wraps := 0, true
	balanced := scanOutsideQuotes(expr, func(i int, c byte) {
		switch c {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i != len(expr)-1 {
				wraps = false
			}
		}
	})
	if !balanced || !wraps || depth != 0 {
		return expr
	}
	return expr[1 : len(expr)-1]
}

// formatDefault formats a default value for SHOW CREATE TABLE output.
//
// MariaDB writes a numeric column's default unquoted whatever it was written
// as: `INT DEFAULT 0` and `INT DEFAULT '7'` come back as `DEFAULT 0` and
// `DEFAULT 7`, and `FLOAT DEFAULT 3.14` as `DEFAULT 3.14`. A string, date or
// enum default keeps its quotes. MySQL 8.0 quotes every literal, numerics
// included. Measured on mariadb:11.8.8 through SHOW CREATE TABLE.
func formatDefault(val string, col *Column) string {
	if strings.EqualFold(val, "NULL") {
		return "NULL"
	}
	switch col.DefaultKind {
	case ColumnDefaultCurrentTimestamp:
		return formatCurrentTimestamp(val)
	case ColumnDefaultExpression:
		// A parenthesized literal is still a literal, however many layers it
		// carries, and is written in the column's own notation without them.
		// Measured on mariadb:11.8.8: ('x') is 'x', ((7)) is 7, (-1) is -1,
		// (1e2) is 100, ('a' 'b') is 'ab' and BINARY(1) DEFAULT (X'FF') is
		// x'ff', while a real expression keeps its parentheses, (1 + 1)
		// staying (1 + 1).
		if inner := stripAllOuterParens(val); isOneLiteral(inner) {
			return formatLiteralValue(inner, col)
		}
		// A default that is one function call is written bare; anything else
		// keeps the parentheses, which the DEFAULT grammar needs for an
		// operator expression. Measured on mariadb:11.8.8:
		// `DEFAULT (uuid())` → DEFAULT uuid(),
		// `DEFAULT (floor(rand() * 100))` → DEFAULT floor(rand() * 100),
		// `DEFAULT (1+1)` → DEFAULT (1 + 1). MySQL 8.0 parenthesizes all of
		// them.
		return formatExpressionDefault(val)
	}
	if _, ok := currentTimestampSQL(val); ok {
		return formatCurrentTimestamp(val)
	}
	// Expression defaults — not quoted, and written without the parentheses
	// they may carry where the grammar does not need them.
	if len(val) >= 2 && val[0] == '(' && val[len(val)-1] == ')' {
		if inner := stripAllOuterParens(val); isOneLiteral(inner) {
			return formatLiteralValue(inner, col)
		}
		return formatExpressionDefault(val)
	}
	return formatLiteralValue(val, col)
}

// formatLiteralValue writes a literal default in the notation the column
// carries it in.
func formatLiteralValue(val string, col *Column) string {
	if out, ok := formatBinaryDefault(val, col); ok {
		return out
	}
	// b'...' and 0x... bit/hex literals — not quoted.
	if strings.HasPrefix(val, "b'") || strings.HasPrefix(val, "B'") ||
		strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		return val
	}
	// A numeric column holds the number the literal spells, not the text, so
	// the default is written from that number: the quotes go, and so does
	// everything else the text carried around the digits.
	if isNumericColumnType(col.DataType) {
		text := val
		if raw, ok := unquoteStringLiteral(val); ok {
			text = string(raw)
		}
		if out, ok := scaleNumberDefault(text, col); ok {
			return out
		}
		// The text spells no number, which the server refuses outright with
		// 1067, so there is no value to write it as. It is written back as it
		// stands, which at least parses.
		return val
	}
	// Already single-quoted string — return as-is.
	if len(val) >= 2 && val[0] == '\'' && val[len(val)-1] == '\'' {
		return val
	}
	return "'" + val + "'"
}

// isNumericColumnType reports whether the type's default is written unquoted.
func isNumericColumnType(dt string) bool {
	switch normalizedTypeName(dt) {
	case "tinyint", "smallint", "mediumint", "int", "bigint",
		"decimal", "float", "double", "boolean", "year":
		return true
	}
	return false
}

// formatOnUpdate normalizes ON UPDATE values to the form MariaDB writes.
func formatOnUpdate(val string, kind ColumnDefaultKind) string {
	if kind == ColumnDefaultCurrentTimestamp {
		return formatCurrentTimestamp(val)
	}
	if _, ok := currentTimestampSQL(val); ok {
		return formatCurrentTimestamp(val)
	}
	return val
}

// formatCurrentTimestamp renders a CURRENT_TIMESTAMP default or ON UPDATE the
// way MariaDB writes it: lowercase and called, `current_timestamp()`, and
// `current_timestamp(3)` where the column carries a precision. MySQL 8.0 writes
// the bare keyword. Measured on mariadb:11.8.8 through SHOW CREATE TABLE.
func formatCurrentTimestamp(val string) string {
	ts, ok := currentTimestampSQL(val)
	if !ok {
		return "current_timestamp()"
	}
	if args, found := strings.CutPrefix(ts, "CURRENT_TIMESTAMP("); found {
		return "current_timestamp(" + args
	}
	return "current_timestamp()"
}

func currentTimestampSQL(val string) (string, bool) {
	upper := strings.ToUpper(strings.TrimSpace(val))
	if upper == "CURRENT_TIMESTAMP" || upper == "CURRENT_TIMESTAMP()" || upper == "NOW()" {
		return "CURRENT_TIMESTAMP", true
	}
	for _, prefix := range []string{"CURRENT_TIMESTAMP(", "NOW("} {
		if strings.HasPrefix(upper, prefix) && strings.HasSuffix(upper, ")") {
			return "CURRENT_TIMESTAMP" + upper[len(prefix)-1:], true
		}
	}
	return "", false
}

// isTimestampType returns true for TIMESTAMP/DATETIME types.
func isTimestampType(dt string) bool {
	switch strings.ToLower(dt) {
	case "timestamp":
		return true
	}
	return false
}

// isTextBlobType returns true for the types an index must give a prefix length
// for.
func isTextBlobType(dt string) bool {
	switch strings.ToLower(dt) {
	case "text", "tinytext", "mediumtext", "longtext",
		"blob", "tinyblob", "mediumblob", "longblob":
		return true
	}
	return false
}

func showIndex(idx *Index) string {
	var b strings.Builder

	if idx.Primary {
		b.WriteString("PRIMARY KEY (")
	} else if idx.Unique {
		b.WriteString(fmt.Sprintf("UNIQUE KEY `%s` (", idx.Name))
	} else if idx.Fulltext {
		b.WriteString(fmt.Sprintf("FULLTEXT KEY `%s` (", idx.Name))
	} else if idx.Spatial {
		b.WriteString(fmt.Sprintf("SPATIAL KEY `%s` (", idx.Name))
	} else if idx.Vector {
		b.WriteString(fmt.Sprintf("VECTOR KEY `%s` (", idx.Name))
	} else {
		b.WriteString(fmt.Sprintf("KEY `%s` (", idx.Name))
	}

	cols := make([]string, 0, len(idx.Columns))
	for _, ic := range idx.Columns {
		cols = append(cols, showIndexColumn(ic))
	}
	b.WriteString(strings.Join(cols, ","))
	b.WriteString(")")

	// USING clause: shown when explicitly specified, not for PRIMARY/FULLTEXT/SPATIAL.
	if !idx.Primary && !idx.Fulltext && !idx.Spatial && !idx.Vector && idx.IndexType != "" {
		b.WriteString(fmt.Sprintf(" USING %s", strings.ToUpper(idx.IndexType)))
	}

	// KEY_BLOCK_SIZE, which MariaDB writes on the index and before its
	// comment, whichever order they were given in:
	// KEY `idx1` (`a`) KEY_BLOCK_SIZE=4 COMMENT 'c'. MySQL 8.0 parses the
	// option but writes only the table-level one. Measured on mariadb:11.8.8.
	if idx.KeyBlockSize > 0 {
		b.WriteString(fmt.Sprintf(" KEY_BLOCK_SIZE=%d", idx.KeyBlockSize))
	}

	// Comment.
	if idx.Comment != "" {
		b.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeComment(idx.Comment)))
	}

	// Engine-defined options follow the comment and precede IGNORED, as MariaDB 11.8.8
	// prints them: VECTOR KEY `i` (`v`) COMMENT 'c' `M`=6 `DISTANCE`=cosine.
	for _, opt := range idx.EngineOptions {
		b.WriteString(" " + opt)
	}

	// Ignored, last, as MariaDB 11.8.8 prints it: KEY `i` (`a`) COMMENT 'c' IGNORED.
	if !idx.Visible {
		b.WriteString(" IGNORED")
	}

	return b.String()
}

func showIndexColumn(ic *IndexColumn) string {
	var b strings.Builder
	if ic.Expr != "" {
		b.WriteString(fmt.Sprintf("(%s)", ic.Expr))
	} else {
		b.WriteString(fmt.Sprintf("`%s`", ic.Name))
		if ic.Length > 0 {
			b.WriteString(fmt.Sprintf("(%d)", ic.Length))
		}
	}
	if ic.Descending {
		b.WriteString(" DESC")
	}
	if ic.WithoutOverlaps {
		b.WriteString(" WITHOUT OVERLAPS")
	}
	return b.String()
}

func showConstraint(con *Constraint) string {
	var b strings.Builder

	switch con.Type {
	case ConForeignKey:
		b.WriteString(fmt.Sprintf("CONSTRAINT `%s` FOREIGN KEY (", con.Name))
		cols := make([]string, 0, len(con.Columns))
		for _, c := range con.Columns {
			cols = append(cols, fmt.Sprintf("`%s`", c))
		}
		b.WriteString(strings.Join(cols, ", "))
		if con.RefDatabase != "" {
			b.WriteString(fmt.Sprintf(") REFERENCES `%s`.`%s` (", con.RefDatabase, con.RefTable))
		} else {
			b.WriteString(fmt.Sprintf(") REFERENCES `%s` (", con.RefTable))
		}
		refCols := make([]string, 0, len(con.RefColumns))
		for _, c := range con.RefColumns {
			refCols = append(refCols, fmt.Sprintf("`%s`", c))
		}
		b.WriteString(strings.Join(refCols, ", "))
		b.WriteString(")")

		// ON DELETE — omit if RESTRICT or NO ACTION (MySQL defaults).
		if con.OnDelete != "" && !isFKDefault(con.OnDelete) {
			b.WriteString(fmt.Sprintf(" ON DELETE %s", strings.ToUpper(con.OnDelete)))
		}
		// ON UPDATE — omit if RESTRICT or NO ACTION (MySQL defaults).
		if con.OnUpdate != "" && !isFKDefault(con.OnUpdate) {
			b.WriteString(fmt.Sprintf(" ON UPDATE %s", strings.ToUpper(con.OnUpdate)))
		}

	case ConCheck:
		b.WriteString(fmt.Sprintf("CONSTRAINT `%s` CHECK (%s)", con.Name, con.CheckExpr))
		if con.NotEnforced {
			b.WriteString(" /*!80016 NOT ENFORCED */")
		}
	}

	return b.String()
}

func showTableOptions(tbl *Table) string {
	var parts []string

	engine := tbl.Engine
	if engine == "" {
		engine = "InnoDB"
	}
	if engine != "" {
		parts = append(parts, fmt.Sprintf("ENGINE=%s", engine))
	}

	// AUTO_INCREMENT — shown only when > 1.
	if tbl.AutoIncrement > 1 {
		parts = append(parts, fmt.Sprintf("AUTO_INCREMENT=%d", tbl.AutoIncrement))
	}

	if tbl.Charset != "" {
		parts = append(parts, fmt.Sprintf("DEFAULT CHARSET=%s", tbl.Charset))
	}

	// MariaDB always writes the table's collation, the charset's own default
	// included: a latin1 table reports COLLATE=latin1_swedish_ci. MySQL 8.0
	// writes it only for utf8mb4 or a non-default collation. Measured on
	// mariadb:11.8.8.
	if tbl.Charset != "" {
		effectiveCollation := tbl.Collation
		if effectiveCollation == "" {
			effectiveCollation = defaultCollationForCharset[toLower(tbl.Charset)]
		}
		if effectiveCollation != "" {
			parts = append(parts, fmt.Sprintf("COLLATE=%s", effectiveCollation))
		}
	}

	// KEY_BLOCK_SIZE — shown when non-zero.
	if tbl.KeyBlockSize > 0 {
		parts = append(parts, fmt.Sprintf("KEY_BLOCK_SIZE=%d", tbl.KeyBlockSize))
	}

	// ROW_FORMAT — shown when explicitly set.
	if tbl.RowFormat != "" {
		parts = append(parts, fmt.Sprintf("ROW_FORMAT=%s", strings.ToUpper(tbl.RowFormat)))
	}

	if tbl.Comment != "" {
		parts = append(parts, fmt.Sprintf("COMMENT='%s'", escapeComment(tbl.Comment)))
	}

	if tbl.Compression != "" {
		parts = append(parts, fmt.Sprintf("COMPRESSION='%s'", escapeComment(unquoteTableOption(tbl.Compression))))
	}
	if tbl.Encryption != "" {
		parts = append(parts, fmt.Sprintf("ENCRYPTION='%s'", escapeComment(unquoteTableOption(tbl.Encryption))))
	}
	if tbl.StatsPersistent != "" && !eqFoldStr(tbl.StatsPersistent, "DEFAULT") {
		parts = append(parts, fmt.Sprintf("STATS_PERSISTENT=%s", tbl.StatsPersistent))
	}
	if tbl.StatsAutoRecalc != "" && !eqFoldStr(tbl.StatsAutoRecalc, "DEFAULT") {
		parts = append(parts, fmt.Sprintf("STATS_AUTO_RECALC=%s", tbl.StatsAutoRecalc))
	}
	if tbl.StatsSamplePages != "" && !eqFoldStr(tbl.StatsSamplePages, "DEFAULT") {
		parts = append(parts, fmt.Sprintf("STATS_SAMPLE_PAGES=%s", tbl.StatsSamplePages))
	}
	if tbl.MinRows != "" && tbl.MinRows != "0" {
		parts = append(parts, fmt.Sprintf("MIN_ROWS=%s", tbl.MinRows))
	}
	if tbl.MaxRows != "" && tbl.MaxRows != "0" {
		parts = append(parts, fmt.Sprintf("MAX_ROWS=%s", tbl.MaxRows))
	}
	if tbl.AvgRowLength != "" && tbl.AvgRowLength != "0" {
		parts = append(parts, fmt.Sprintf("AVG_ROW_LENGTH=%s", tbl.AvgRowLength))
	}
	if tbl.Tablespace != "" {
		parts = append(parts, fmt.Sprintf("/*!50100 TABLESPACE `%s` */", unquoteTableOption(tbl.Tablespace)))
	}
	if tbl.PackKeys != "" && !eqFoldStr(tbl.PackKeys, "DEFAULT") {
		parts = append(parts, fmt.Sprintf("PACK_KEYS=%s", tbl.PackKeys))
	}
	if tbl.Checksum != "" && tbl.Checksum != "0" {
		parts = append(parts, fmt.Sprintf("CHECKSUM=%s", tbl.Checksum))
	}
	if tbl.DelayKeyWrite != "" && tbl.DelayKeyWrite != "0" {
		parts = append(parts, fmt.Sprintf("DELAY_KEY_WRITE=%s", tbl.DelayKeyWrite))
	}

	return strings.Join(parts, " ")
}

func unquoteTableOption(s string) string {
	return strings.Trim(s, "'\"`")
}

// isFKDefault returns true if the action is a MySQL FK default that should not be shown.
// MySQL 8.0 hides NO ACTION (the implicit default) but shows RESTRICT when explicitly specified.
// isFKDefault reports whether MariaDB leaves a referential action out of SHOW
// CREATE TABLE. It leaves out RESTRICT, which is what an unwritten action and a
// written one both mean, and writes NO ACTION; MySQL 8.0 leaves out NO ACTION
// and RESTRICT. Measured on mariadb:11.8.8.
func isFKDefault(action string) bool {
	return strings.EqualFold(action, "RESTRICT")
}

func isEnumSetType(dt string) bool {
	switch strings.ToLower(dt) {
	case "enum", "set":
		return true
	}
	return false
}

// isExpressionIndex returns true if any column in the index is an expression.
func isExpressionIndex(idx *Index) bool {
	for _, ic := range idx.Columns {
		if ic.Expr != "" {
			return true
		}
	}
	return false
}

func eqFoldStr(a, b string) bool {
	return strings.EqualFold(a, b)
}

func escapeComment(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "''")
	return s
}

// showPartitioning renders the partition clause for SHOW CREATE TABLE.
//
// MariaDB writes the clause plainly, with one leading space and no version
// comment, and backticks every partition name and every COLUMNS column:
//
//	 PARTITION BY RANGE (`val`)
//	(PARTITION `p0` VALUES LESS THAN (10) ENGINE = InnoDB)
//
// MySQL 8.0 wraps the whole clause in /*!50100 ... */ (/*!50500 for the COLUMNS
// forms) and leaves those names bare. Measured on mariadb:11.8.8.
func showPartitioning(pi *PartitionInfo) string {
	var b strings.Builder

	b.WriteString(" PARTITION BY ")
	if pi.Linear {
		b.WriteString("LINEAR ")
	}
	switch pi.Type {
	case "RANGE":
		b.WriteString(fmt.Sprintf("RANGE (%s)", pi.Expr))
	case "RANGE COLUMNS":
		// Double space before COLUMNS, as the server writes it.
		b.WriteString(fmt.Sprintf("RANGE  COLUMNS(%s)", formatPartitionColumns(pi.Columns)))
	case "LIST":
		b.WriteString(fmt.Sprintf("LIST (%s)", pi.Expr))
	case "LIST COLUMNS":
		b.WriteString(fmt.Sprintf("LIST  COLUMNS(%s)", formatPartitionColumns(pi.Columns)))
	case "HASH":
		b.WriteString(fmt.Sprintf("HASH (%s)", pi.Expr))
	case "KEY":
		if pi.Algorithm > 0 {
			b.WriteString(fmt.Sprintf("KEY ALGORITHM = %d (%s)", pi.Algorithm, formatPartitionColumns(pi.Columns)))
		} else {
			b.WriteString(fmt.Sprintf("KEY (%s)", formatPartitionColumns(pi.Columns)))
		}
	case "SYSTEM_TIME":
		b.WriteString("SYSTEM_TIME")
		if pi.IntervalValue != "" {
			b.WriteString(fmt.Sprintf(" INTERVAL %s %s", pi.IntervalValue, pi.IntervalUnit))
		} else if pi.Limit > 0 {
			b.WriteString(fmt.Sprintf(" LIMIT %d", pi.Limit))
		}
		if pi.Starts != "" {
			b.WriteString(fmt.Sprintf(" STARTS %s", pi.Starts))
		}
		if pi.Auto {
			b.WriteString(" AUTO")
		}
	}

	// Subpartition clause.
	if pi.SubType != "" {
		b.WriteString("\n")
		b.WriteString("SUBPARTITION BY ")
		if pi.SubLinear {
			b.WriteString("LINEAR ")
		}
		switch pi.SubType {
		case "HASH":
			b.WriteString(fmt.Sprintf("HASH (%s)", pi.SubExpr))
		case "KEY":
			if pi.SubAlgo > 0 {
				b.WriteString(fmt.Sprintf("KEY ALGORITHM = %d (%s)", pi.SubAlgo, formatPartitionColumns(pi.SubColumns)))
			} else {
				b.WriteString(fmt.Sprintf("KEY (%s)", formatPartitionColumns(pi.SubColumns)))
			}
		}
		if pi.UseDefaultSubpartitions && !pi.UseDefaultNumSubpartitions && pi.NumSubParts > 0 {
			b.WriteString(fmt.Sprintf("\nSUBPARTITIONS %d", pi.NumSubParts))
		}
	}

	// Partition definitions.
	// For HASH/KEY partitions with NumParts > 0, auto-generated partition defs
	// are rendered as "PARTITIONS N" (matching MySQL 8.0's SHOW CREATE TABLE).
	autoGenParts := pi.NumParts > 0 && (pi.Type == "HASH" || pi.Type == "KEY" ||
		(pi.Type == "SYSTEM_TIME" && len(pi.Partitions) == 0))
	if autoGenParts {
		b.WriteString(fmt.Sprintf("\nPARTITIONS %d", pi.NumParts))
	} else if len(pi.Partitions) > 0 {
		b.WriteString("\n(")
		for i, pd := range pi.Partitions {
			if i > 0 {
				b.WriteString(",\n ")
			}
			b.WriteString("PARTITION ")
			b.WriteString(backtickIDQ(pd.Name))
			if pd.SystemTime != "" {
				b.WriteString(" " + pd.SystemTime)
			}
			if pd.ValueExpr != "" {
				switch {
				case strings.HasPrefix(pi.Type, "RANGE"):
					if pd.ValueExpr == "MAXVALUE" {
						if pi.Type == "RANGE COLUMNS" {
							// RANGE COLUMNS uses parenthesized MAXVALUE
							b.WriteString(" VALUES LESS THAN (MAXVALUE)")
						} else {
							b.WriteString(" VALUES LESS THAN MAXVALUE")
						}
					} else {
						b.WriteString(fmt.Sprintf(" VALUES LESS THAN (%s)", pd.ValueExpr))
					}
				case strings.HasPrefix(pi.Type, "LIST"):
					b.WriteString(fmt.Sprintf(" VALUES IN (%s)", pd.ValueExpr))
				}
			}
			if pi.SubType == "" || pi.UseDefaultSubpartitions {
				// MariaDB writes a partition's comment before its engine:
				// PARTITION `p0` VALUES LESS THAN (10) COMMENT = 'pc0'
				// ENGINE = InnoDB. Measured on mariadb:11.8.8.
				if pd.Comment != "" {
					b.WriteString(fmt.Sprintf(" COMMENT = '%s'", escapeComment(pd.Comment)))
				}
				b.WriteString(fmt.Sprintf(" ENGINE = %s", partitionEngine(pd.Engine)))
			}
			// Explicit subpartition definitions are rendered instead of parent options.
			if len(pd.SubPartitions) > 0 && !pi.UseDefaultSubpartitions {
				b.WriteString("\n (")
				for j, spd := range pd.SubPartitions {
					if j > 0 {
						b.WriteString(",\n  ")
					}
					b.WriteString("SUBPARTITION ")
					b.WriteString(backtickIDQ(spd.Name))
					if spd.Comment != "" {
						b.WriteString(fmt.Sprintf(" COMMENT = '%s'", escapeComment(spd.Comment)))
					}
					b.WriteString(fmt.Sprintf(" ENGINE = %s", partitionEngine(spd.Engine)))
				}
				b.WriteString(")")
			}
		}
		b.WriteString(")")
	}

	return b.String()
}

// formatPartitionColumns formats column names for partition clauses with backticks.
func formatPartitionColumns(cols []string) string {
	parts := make([]string, len(cols))
	for i, c := range cols {
		parts[i] = "`" + c + "`"
	}
	return strings.Join(parts, ",")
}

// partitionEngine returns the engine for a partition, defaulting to InnoDB.
func partitionEngine(engine string) string {
	if engine == "" {
		return "InnoDB"
	}
	return engine
}

// formatBinaryDefault renders a literal default on a binary column, and
// reports whether it applied. The two halves of the family answer differently,
// both measured on mariadb:11.8.8 through SHOW CREATE TABLE.
//
// BINARY and VARBINARY store the coerced value, so the literal is evaluated
// and written back as lowercase hex, padded with zero bytes to the declared
// width where the type is fixed:
//
//	BINARY(1) DEFAULT X'FF'  -> x'ff'        BINARY(4) DEFAULT X'FF' -> x'ff000000'
//	BINARY(4) DEFAULT 'a'    -> x'61000000'  BINARY(3) DEFAULT 0x4142 -> x'414200'
//	BINARY(1) DEFAULT X''    -> x'00'        VARBINARY(4) DEFAULT 'a' -> x'61'
//	VARBINARY(4) DEFAULT b'1000001' -> x'41'
//
// The BLOB family stores the expression instead, so the literal is kept as it
// was written and only normalized: hex keeps its own spelling with the digits
// lowercased, a bit literal becomes 0x, and a string stays a string:
//
//	BLOB DEFAULT X'FF' -> X'ff'   BLOB DEFAULT 0xFF -> 0xff
//	BLOB DEFAULT b'1000001' -> 0x41   BLOB DEFAULT 'a' -> 'a'
func formatBinaryDefault(val string, col *Column) (string, bool) {
	family := columnDefaultFamily(col.DataType)
	if family == familyOther {
		return "", false
	}
	kind, raw, ok := literalDefaultBytes(val)
	if !ok {
		return "", false
	}
	// A column whose default is held as an expression keeps the literal as it
	// was written, only normalized. Everything else holds a value, so the
	// literal is coerced and written in that type's own notation.
	if family == familyExpression {
		switch kind {
		case literalHexX:
			return val[:2] + strings.ToLower(val[2:]), true
		case literalHex0x:
			return "0x" + strings.ToLower(val[2:]), true
		case literalBits:
			return "0x" + hexOf(raw), true
		default:
			return val, true // a string default keeps its quotes
		}
	}
	if (family == familyString || family == familyNumber) && kind == literalString {
		// A plain string literal is already in the column's own notation. A
		// hex, bit or introduced one has to be converted: measured on
		// mariadb:11.8.8, `INT DEFAULT _utf8mb4'7'` comes back DEFAULT 7,
		// `DECIMAL(5,2) DEFAULT _utf8mb4'7'` DEFAULT 7.00, and
		// `VARCHAR(4) DEFAULT _latin1'x'` DEFAULT 'x' with the introducer
		// gone, while a TEXT column keeps the literal as written.
		return "", false
	}
	switch family {
	case familyFixedBinary, familyVarBinary:
		if width := declaredTypeLength(col.ColumnType); family == familyFixedBinary && width > len(raw) {
			raw = append(raw, make([]byte, width-len(raw))...)
		}
		return "x'" + hexOf(raw) + "'", true
	case familyString:
		return quoteStringDefault(raw), true
	case familyNumber:
		return numberDefaultOf(kind, raw, col)
	}
	return "", false
}

// columnDefaultFamily groups the column types by how a literal default is
// written back, each group measured on mariadb:11.8.8 through SHOW CREATE
// TABLE. See formatBinaryDefault.
type defaultFamily int

const (
	familyOther defaultFamily = iota
	// familyExpression holds its default as an expression, so the literal
	// survives as written: TEXT DEFAULT 0x41 stays 0x41 and TEXT DEFAULT
	// X'41' stays X'41', where a bit literal becomes 0x41.
	familyExpression
	// familyFixedBinary pads the value to the declared width:
	// BINARY(4) DEFAULT X'FF' is x'ff000000'.
	familyFixedBinary
	// familyVarBinary writes the value as lowercase hex, unpadded.
	familyVarBinary
	// familyString writes the value as a quoted string:
	// VARCHAR(4) DEFAULT 0x41 is 'A', and so is ENUM('A','B') DEFAULT 0x41.
	familyString
	// familyNumber reads the bytes as a number: INT DEFAULT 0xFF is 255,
	// DECIMAL(5,2) DEFAULT 0x41 is 65.00, YEAR DEFAULT 0x07D0 is 2000.
	familyNumber
)

func columnDefaultFamily(dt string) defaultFamily {
	switch normalizedTypeName(dt) {
	case "tinytext", "text", "mediumtext", "longtext",
		"tinyblob", "blob", "mediumblob", "longblob":
		return familyExpression
	case "binary":
		return familyFixedBinary
	case "varbinary":
		return familyVarBinary
	case "char", "varchar", "enum", "set":
		return familyString
	case "tinyint", "smallint", "mediumint", "int", "bigint",
		"decimal", "float", "double", "year":
		return familyNumber
	}
	return familyOther
}

// quoteStringDefault writes bytes as a quoted string default. MariaDB escapes
// five bytes and leaves the rest alone, control bytes and UTF-8 included:
// 0x27 is ”, 0x5C is \\, 0x0A is \n, 0x0D is \r and 0x00 is \0, while
// 0x08, 0x09, 0x0B, 0x0C, 0x1A, 0x22 and 0x7F come through raw. Measured on
// mariadb:11.8.8 one byte at a time.
func quoteStringDefault(raw []byte) string {
	var b strings.Builder
	b.Grow(len(raw) + 2)
	b.WriteByte('\'')
	for _, c := range raw {
		switch c {
		case '\'':
			b.WriteString("''")
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case 0:
			b.WriteString(`\0`)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte('\'')
	return b.String()
}

// numberDefaultOf reads the bytes as a big-endian unsigned number and writes
// it in the column's notation, which for DECIMAL carries the declared scale.
func numberDefaultOf(kind literalKind, raw []byte, col *Column) (string, bool) {
	// A hex or bit literal is the number its bytes spell, while a string one
	// is the number its text spells: INT DEFAULT 0x37 is 55 and
	// INT DEFAULT '7' is 7, the same byte read two ways. Measured on
	// mariadb:11.8.8.
	digits := ""
	switch kind {
	case literalString, literalIntroduced:
		digits = string(raw)
	default:
		if len(raw) > 8 {
			return "", false // wider than a number the column could hold
		}
		var n uint64
		for _, c := range raw {
			n = n<<8 | uint64(c)
		}
		digits = strconv.FormatUint(n, 10)
	}
	return scaleNumberDefault(digits, col)
}

// scaleNumberDefault writes the number a default's text spells the way the
// column carries it. A DECIMAL is shown to its declared scale however the
// default was written: DECIMAL(5,2) DEFAULT 7, '7', 7.5 and 0x41 come back
// 7.00, 7.00, 7.50 and 65.00, while DECIMAL(5,0) DEFAULT 7 stays 7. Measured
// on mariadb:11.8.8.
//
// It reports false where the text spells no number, as `INT DEFAULT 'abc'`
// does. The server refuses such a default outright with 1067, so there is no
// value for it to be written as; omni does not check that yet.
//
// YEAR is the one numeric type whose value is more than its digits: measured,
// `YEAR DEFAULT '99'` comes back 1999 and `YEAR DEFAULT '0'` 2000, while
// unquoted 0 stays 0000. That expansion is not done here.
func scaleNumberDefault(text string, col *Column) (string, bool) {
	digits, ok := normalizeNumberText(text)
	if !ok {
		return "", false
	}
	scale, declared := declaredScaleOf(col.ColumnType)
	switch normalizedTypeName(col.DataType) {
	case "decimal":
		// DECIMAL is decimal arithmetic, so a tie rounds away from zero:
		// measured, DECIMAL(5,2) DEFAULT 7.125 comes back 7.13, 0.005 0.01
		// and 9.999 10.00. A DECIMAL that declares no scale has scale 0, so
		// DECIMAL DEFAULT 7.9 comes back 8.
		return dropNegativeZero(roundDecimalText(digits, scale)), true
	case "float", "double":
		if !declared {
			// A scaleless float prints as the double it holds, so a trailing
			// run of zeros in the fraction is notation rather than value:
			// measured, DOUBLE DEFAULT 7.000 comes back 7 and 7.10 comes
			// back 7.1.
			return dropNegativeZero(trimFractionZeros(digits)), true
		}
		// A float that declares a scale is shown to it, but rounded as the
		// binary double it holds rather than as decimal digits: measured,
		// DOUBLE(10,2) DEFAULT 0.125 comes back 0.12 where DECIMAL(10,2)
		// DEFAULT 0.125 comes back 0.13. That is what FormatFloat does.
		if f, err := strconv.ParseFloat(digits, 64); err == nil {
			return dropNegativeZero(strconv.FormatFloat(f, 'f', scale, 64)), true
		}
		return digits, true
	default:
		// The integer types hold no fraction at all, and round away from zero
		// the same way: measured, INT DEFAULT 7.9 comes back 8, 7.5 comes
		// back 8 and -7.5 comes back -8.
		return dropNegativeZero(roundDecimalText(digits, 0)), true
	}
}

// normalizeNumberText reads the number a default's text spells and writes it
// as plain signed digits. Surrounding whitespace, a leading plus, leading
// zeros, a bare point and exponent notation are all notation rather than
// value, so the server keeps none of them: measured on mariadb:11.8.8,
// `DECIMAL(5,2) DEFAULT ' 7 '` comes back 7.00, `INT DEFAULT '007.00'` 7,
// `DECIMAL(5,2) DEFAULT '.7'` 0.70 and `DECIMAL(5,2) DEFAULT '7.'` 7.00.
//
// A number default is written as a plain decimal, so exponent notation is
// worked out rather than carried: measured, `DECIMAL(5,2) DEFAULT 1e2` comes
// back 100.00, `DECIMAL(6,2) DEFAULT 1.5e2` 150.00, `DECIMAL(5,2) DEFAULT
// 1e-1` 0.10, and on a column with no scale, `INT DEFAULT 1e2` and `DOUBLE
// DEFAULT 1e2` both come back 100. A double large enough that the server
// prints it back in exponent notation, as DOUBLE DEFAULT 1e15 does, is
// expanded here instead.
func normalizeNumberText(text string) (string, bool) {
	text = strings.Trim(text, " \t\n\r\v\f")
	sign := ""
	if text != "" && (text[0] == '+' || text[0] == '-') {
		if text[0] == '-' {
			sign = "-"
		}
		text = text[1:]
	}
	if plain, ok := expandExponent(text); ok {
		text = plain
	}
	whole, fraction, hasPoint := strings.Cut(text, ".")
	if whole+fraction == "" || !allDigits(whole) || !allDigits(fraction) {
		return "", false
	}
	if whole = strings.TrimLeft(whole, "0"); whole == "" {
		whole = "0"
	}
	if hasPoint && fraction != "" {
		return sign + whole + "." + fraction, true
	}
	return sign + whole, true
}

// roundDecimalText writes signed plain digits to the given scale, padding the
// fraction with zeros or rounding it away from zero at the first digit the
// scale cannot hold.
func roundDecimalText(digits string, scale int) string {
	sign := ""
	if strings.HasPrefix(digits, "-") {
		sign, digits = "-", digits[1:]
	}
	whole, fraction, _ := strings.Cut(digits, ".")
	if len(fraction) > scale {
		roundUp := fraction[scale] >= '5'
		fraction = fraction[:scale]
		if roundUp {
			whole, fraction = carryOne(whole, fraction)
		}
	}
	fraction += strings.Repeat("0", scale-len(fraction))
	if scale == 0 {
		return sign + whole
	}
	return sign + whole + "." + fraction
}

// carryOne adds one to the last digit of a number held as a whole part and a
// fraction of fixed width, carrying into the whole part and widening it where
// every digit was a nine.
func carryOne(whole, fraction string) (string, string) {
	digits := []byte(whole + fraction)
	i := len(digits) - 1
	for ; i >= 0; i-- {
		if digits[i] != '9' {
			digits[i]++
			break
		}
		digits[i] = '0'
	}
	if i < 0 {
		digits = append([]byte{'1'}, digits...)
	}
	return string(digits[:len(digits)-len(fraction)]), string(digits[len(digits)-len(fraction):])
}

// trimFractionZeros drops a fraction's trailing zeros, and the point with them
// where nothing of the fraction is left.
func trimFractionZeros(digits string) string {
	if !strings.Contains(digits, ".") {
		return digits
	}
	return strings.TrimSuffix(strings.TrimRight(digits, "0"), ".")
}

// dropNegativeZero drops the sign from a value that came out zero, as a
// rounded one can: measured on mariadb:11.8.8, DECIMAL(5,2) DEFAULT '-0.001'
// comes back 0.00 and DOUBLE DEFAULT '-0.0' comes back 0.
func dropNegativeZero(digits string) string {
	if strings.HasPrefix(digits, "-") && strings.Trim(digits[1:], "0.") == "" {
		return digits[1:]
	}
	return digits
}

// literalKind is the notation a literal default was written in.
type literalKind int

const (
	literalNone       literalKind = iota
	literalHexX                   // X'ff'
	literalHex0x                  // 0xff
	literalBits                   // b'101'
	literalString                 // 'a'
	literalIntroduced             // _utf8mb4'a'
)

// literalDefaultBytes decodes a literal default to the bytes it stands for and
// reports the notation it was written in. It reports false for anything else,
// an expression included.
func literalDefaultBytes(val string) (literalKind, []byte, bool) {
	switch {
	case hasPrefixFold(val, "x'") && strings.HasSuffix(val, "'"):
		raw, ok := decodeHex(val[2 : len(val)-1])
		return literalHexX, raw, ok
	case hasPrefixFold(val, "0x"):
		raw, ok := decodeHex(val[2:])
		return literalHex0x, raw, ok
	case hasPrefixFold(val, "b'") && strings.HasSuffix(val, "'"):
		raw, ok := decodeBits(val[2 : len(val)-1])
		return literalBits, raw, ok
	case len(val) >= 2 && val[0] == '\'' && val[len(val)-1] == '\'':
		raw, ok := unquoteStringLiteral(val)
		return literalString, raw, ok
	case val != "" && val[0] == '_':
		// A character-set introducer, `_utf8mb4'a'`, is another spelling of a
		// string literal: the charset says how to read the bytes and the value
		// is the quoted part.
		quote := strings.IndexByte(val, '\'')
		if quote < 2 {
			return literalNone, nil, false
		}
		for i := 1; i < quote; i++ {
			if !isNameByte(val[i]) {
				return literalNone, nil, false
			}
		}
		raw, ok := unquoteStringLiteral(val[quote:])
		return literalIntroduced, raw, ok
	}
	return literalNone, nil, false
}

// unquoteStringLiteral takes the quotes off a single-quoted literal, where a
// doubled quote stands for one. A literal holding a backslash escape is
// refused, and its caller keeps it as written.
func unquoteStringLiteral(val string) ([]byte, bool) {
	if !isOneStringLiteral(val) {
		return nil, false
	}
	inner := val[1 : len(val)-1]
	if strings.ContainsRune(inner, '\\') {
		return nil, false
	}
	return []byte(strings.ReplaceAll(inner, "''", "'")), true
}

// declaredScaleOf reads the scale out of a column type, as the 2 of
// decimal(5,2), and reports whether the type declares one at all: decimal(5),
// decimal and double do not, and a double that declares no scale is written
// back differently from one that declares 0.
func declaredScaleOf(columnType string) (int, bool) {
	open := strings.IndexByte(columnType, '(')
	if open < 0 {
		return 0, false
	}
	rest := columnType[open+1:]
	comma := strings.IndexByte(rest, ',')
	closing := strings.IndexByte(rest, ')')
	if comma < 0 || closing < 0 || comma > closing {
		return 0, false
	}
	scale := strings.TrimSpace(rest[comma+1 : closing])
	if scale == "" || !allDigits(scale) {
		return 0, false
	}
	n := 0
	for i := 0; i < len(scale); i++ {
		n = n*10 + int(scale[i]-'0')
	}
	return n, true
}

// decodeHex reads an even-length run of hex digits.
func decodeHex(digits string) ([]byte, bool) {
	if len(digits)%2 != 0 {
		digits = "0" + digits
	}
	out := make([]byte, 0, len(digits)/2)
	for i := 0; i < len(digits); i += 2 {
		hi, ok := hexNibble(digits[i])
		if !ok {
			return nil, false
		}
		lo, ok2 := hexNibble(digits[i+1])
		if !ok2 {
			return nil, false
		}
		out = append(out, hi<<4|lo)
	}
	return out, true
}

func hexNibble(c byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// decodeBits reads a bit literal as a big-endian number in as few bytes as
// hold it, which is how b'1000001' becomes one byte, 0x41.
func decodeBits(bits string) ([]byte, bool) {
	if bits == "" {
		return nil, true
	}
	for i := 0; i < len(bits); i++ {
		if bits[i] != '0' && bits[i] != '1' {
			return nil, false
		}
	}
	n := (len(bits) + 7) / 8
	out := make([]byte, n)
	for i := 0; i < len(bits); i++ {
		if bits[len(bits)-1-i] == '1' {
			out[n-1-i/8] |= 1 << (i % 8)
		}
	}
	return out, true
}

// hexOf renders bytes as lowercase hex digits.
func hexOf(raw []byte) string {
	const digits = "0123456789abcdef"
	var b strings.Builder
	b.Grow(len(raw) * 2)
	for _, c := range raw {
		b.WriteByte(digits[c>>4])
		b.WriteByte(digits[c&0xf])
	}
	return b.String()
}

// hasPrefixFold reports whether s starts with prefix, ignoring case.
func hasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

// declaredTypeLength reads the width out of a column type, as the 4 of
// binary(4), or 0 where the type carries none.
func declaredTypeLength(columnType string) int {
	open := strings.IndexByte(columnType, '(')
	if open < 0 {
		return 0
	}
	close := strings.IndexByte(columnType[open:], ')')
	if close < 0 {
		return 0
	}
	n := 0
	digits := columnType[open+1 : open+close]
	if digits == "" {
		return 0
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] < '0' || digits[i] > '9' {
			return 0
		}
		n = n*10 + int(digits[i]-'0')
	}
	return n
}

// expandExponent rewrites a number in exponent notation as a plain decimal,
// and reports whether it applied. The point is moved through the digits rather
// than the value passed through a float, so no precision is invented.
func expandExponent(text string) (string, bool) {
	mantissa, exponent, found := cutAnyFold(text, "e")
	if !found || mantissa == "" || exponent == "" {
		return "", false
	}
	sign := ""
	if mantissa[0] == '+' || mantissa[0] == '-' {
		if mantissa[0] == '-' {
			sign = "-"
		}
		mantissa = mantissa[1:]
	}
	exp, ok := parseSignedInt(exponent)
	if !ok {
		return "", false
	}
	whole, fraction, _ := strings.Cut(mantissa, ".")
	if whole == "" && fraction == "" {
		return "", false
	}
	if !allDigits(whole) || !allDigits(fraction) {
		return "", false
	}

	digits := whole + fraction
	point := len(whole) + exp
	var out string
	switch {
	case point <= 0:
		out = "0." + strings.Repeat("0", -point) + digits
	case point >= len(digits):
		out = digits + strings.Repeat("0", point-len(digits))
	default:
		out = digits[:point] + "." + digits[point:]
	}

	// A shifted mantissa can carry leading zeros, as 0.5e1 does; one digit
	// stays before the point.
	whole, fraction, hasPoint := strings.Cut(out, ".")
	whole = strings.TrimLeft(whole, "0")
	if whole == "" {
		whole = "0"
	}
	if hasPoint {
		return sign + whole + "." + fraction, true
	}
	return sign + whole, true
}

// cutAnyFold cuts text at the first occurrence of sep in either case.
func cutAnyFold(text, sep string) (before, after string, found bool) {
	i := strings.IndexAny(text, strings.ToLower(sep)+strings.ToUpper(sep))
	if i < 0 {
		return text, "", false
	}
	return text[:i], text[i+1:], true
}

// parseSignedInt reads an optionally signed run of digits.
func parseSignedInt(text string) (int, bool) {
	negative := false
	if text != "" && (text[0] == '+' || text[0] == '-') {
		negative = text[0] == '-'
		text = text[1:]
	}
	if text == "" || !allDigits(text) || len(text) > 4 {
		return 0, false // a wild exponent is left as written
	}
	n := 0
	for i := 0; i < len(text); i++ {
		n = n*10 + int(text[i]-'0')
	}
	if negative {
		return -n, true
	}
	return n, true
}

// allDigits reports whether text is only decimal digits; the empty string is.
func allDigits(text string) bool {
	for i := 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			return false
		}
	}
	return true
}

// isOneStringLiteral reports whether val is a single quoted string and nothing
// else: the quote opening it closes at its last byte, where a doubled quote
// stands for one. `'1' + '2'` is two literals around an operator and so is not
// one, which keeps its source text from being read as a value.
func isOneStringLiteral(val string) bool {
	if len(val) < 2 || val[0] != '\'' || val[len(val)-1] != '\'' {
		return false
	}
	outside := false
	balanced := scanOutsideQuotes(val, func(int, byte) { outside = true })
	return balanced && !outside
}

// isOneLiteral reports whether text is a single literal: a quoted string, a
// hex or bit literal, one with a character-set introducer, or a number.
func isOneLiteral(text string) bool {
	if text == "" {
		return false
	}
	if kind, _, ok := literalDefaultBytes(text); ok {
		if kind != literalString {
			return true
		}
		return isOneStringLiteral(text)
	}
	return isNumberLiteral(text)
}

// isNumberLiteral reports whether text is a number and nothing else, in plain
// or exponent notation, with an optional sign.
func isNumberLiteral(text string) bool {
	if text != "" && (text[0] == '+' || text[0] == '-') {
		text = text[1:]
	}
	if mantissa, exponent, found := cutAnyFold(text, "e"); found {
		if _, ok := parseSignedInt(exponent); !ok {
			return false
		}
		text = mantissa
	}
	whole, fraction, _ := strings.Cut(text, ".")
	if whole == "" && fraction == "" {
		return false
	}
	return allDigits(whole) && allDigits(fraction)
}

// stripAllOuterParens takes off every layer of parentheses wrapping the whole
// expression, so ((7)) reduces to 7.
func stripAllOuterParens(expr string) string {
	for {
		inner := stripOuterParens(expr)
		if inner == expr {
			return expr
		}
		expr = strings.TrimSpace(inner)
	}
}
