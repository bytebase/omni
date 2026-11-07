package catalog

import (
	"fmt"
	"sort"
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
	// A CHECK named after a column of the table is that column's own, and was
	// written on the column.
	var fkConstraints, chkConstraints []*Constraint
	for _, con := range tbl.Constraints {
		switch con.Type {
		case ConForeignKey:
			fkConstraints = append(fkConstraints, con)
		case ConCheck:
			if columnCheckOf(tbl, con.Name) == nil {
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

// columnCheckOf returns the CHECK constraint that belongs to the named column,
// which is the one carrying the column's name, or nil.
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
		if con.Type == ConCheck && eqFoldStr(con.Name, name) {
			return con
		}
	}
	return nil
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
		// A parenthesized literal is still a literal: `BINARY(1) DEFAULT
		// (X'FF')` comes back `DEFAULT x'ff'`, the same as without the
		// parentheses. Measured on mariadb:11.8.8.
		if out, ok := formatBinaryDefault(stripOuterParens(val), col); ok {
			return out
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
	if out, ok := formatBinaryDefault(val, col); ok {
		return out
	}
	// b'...' and 0x... bit/hex literals — not quoted.
	if strings.HasPrefix(val, "b'") || strings.HasPrefix(val, "B'") ||
		strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		return val
	}
	// Expression defaults — not quoted, and written without the parentheses
	// they may carry where the grammar does not need them.
	if len(val) >= 2 && val[0] == '(' && val[len(val)-1] == ')' {
		return formatExpressionDefault(val)
	}
	// A numeric column's default is unquoted, and loses the quotes it was
	// written with.
	if isNumericColumnType(col.DataType) {
		return strings.Trim(val, "'")
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
	fixed, variable, blob := isBinaryColumnType(col.DataType)
	if !fixed && !variable && !blob {
		return "", false
	}
	raw, ok := binaryLiteralBytes(val)
	if !ok {
		return "", false
	}
	if blob {
		switch {
		case hasPrefixFold(val, "x'"):
			return val[:2] + strings.ToLower(val[2:]), true
		case hasPrefixFold(val, "0x"):
			return "0x" + strings.ToLower(val[2:]), true
		case hasPrefixFold(val, "b'"):
			return "0x" + hexOf(raw), true
		default:
			return val, true // a string default keeps its quotes
		}
	}
	if width := declaredTypeLength(col.ColumnType); fixed && width > len(raw) {
		raw = append(raw, make([]byte, width-len(raw))...)
	}
	return "x'" + hexOf(raw) + "'", true
}

// isBinaryColumnType classifies a column type as fixed binary, variable
// binary, or one of the blob types.
func isBinaryColumnType(dt string) (fixed, variable, blob bool) {
	switch normalizedTypeName(dt) {
	case "binary":
		return true, false, false
	case "varbinary":
		return false, true, false
	case "tinyblob", "blob", "mediumblob", "longblob":
		return false, false, true
	}
	return false, false, false
}

// binaryLiteralBytes decodes a literal default to the bytes it stands for: a
// hex literal in either spelling, a bit literal, or a single-quoted string.
// It reports false for anything else, an expression included.
func binaryLiteralBytes(val string) ([]byte, bool) {
	switch {
	case hasPrefixFold(val, "x'") && strings.HasSuffix(val, "'"):
		return decodeHex(val[2 : len(val)-1])
	case hasPrefixFold(val, "0x"):
		return decodeHex(val[2:])
	case hasPrefixFold(val, "b'") && strings.HasSuffix(val, "'"):
		return decodeBits(val[2 : len(val)-1])
	case len(val) >= 2 && val[0] == '\'' && val[len(val)-1] == '\'':
		// The quotes come off and a doubled quote stands for one. A string
		// holding an escape is left to the caller, which keeps it as written.
		inner := val[1 : len(val)-1]
		if strings.ContainsRune(inner, '\\') {
			return nil, false
		}
		return []byte(strings.ReplaceAll(inner, "''", "'")), true
	}
	return nil, false
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
