package catalog

import (
	"fmt"
	"strings"
)

// defaultCollationForCharset returns the default collation for common MySQL charsets.
var defaultCollationForCharset = map[string]string{
	"utf8mb4":  "utf8mb4_0900_ai_ci",
	"utf8mb3":  "utf8mb3_general_ci",
	"utf8":     "utf8mb3_general_ci",
	"latin1":   "latin1_swedish_ci",
	"ascii":    "ascii_general_ci",
	"binary":   "binary",
	"gbk":      "gbk_chinese_ci",
	"big5":     "big5_chinese_ci",
	"euckr":    "euckr_korean_ci",
	"gb2312":   "gb2312_chinese_ci",
	"sjis":     "sjis_japanese_ci",
	"cp1252":   "cp1252_general_ci",
	"ucs2":     "ucs2_general_ci",
	"utf16":    "utf16_general_ci",
	"utf16le":  "utf16le_general_ci",
	"utf32":    "utf32_general_ci",
	"cp932":    "cp932_japanese_ci",
	"eucjpms":  "eucjpms_japanese_ci",
	"gb18030":  "gb18030_chinese_ci",
	"geostd8":  "geostd8_general_ci",
	"tis620":   "tis620_thai_ci",
	"hebrew":   "hebrew_general_ci",
	"greek":    "greek_general_ci",
	"armscii8": "armscii8_general_ci",
}

// ShowCreateTable produces MySQL 8.0-compatible SHOW CREATE TABLE output.
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
		parts = append(parts, showColumnWithTable(col, tbl))
	}

	// Indexes.
	for _, idx := range tbl.Indexes {
		parts = append(parts, showIndex(idx))
	}

	// Constraints (FK and CHECK only — PK/UNIQUE are shown via indexes).
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey || con.Type == ConCheck {
			parts = append(parts, showConstraint(con))
		}
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

	return b.String()
}

func showColumn(col *Column) string {
	return showColumnWithTable(col, nil)
}

func showColumnWithTable(col *Column, tbl *Table) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("`%s` %s", col.Name, col.ColumnType))

	// CHARACTER SET and COLLATE — MySQL 8.0 display rules:
	// - CHARACTER SET shown when column charset differs from table charset
	// - COLLATE shown when column collation differs from the charset's default collation
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
		// Show COLLATE when the column's collation differs from its charset's default.
		collationNonDefault := col.Collation != "" && !eqFoldStr(col.Collation, colCharsetDefault)

		// Determine the table's effective collation for comparison.
		tableCollation := ""
		if tbl != nil {
			tableCollation = tbl.Collation
		}
		// Column collation differs from table collation (= explicitly set on column).
		collationDiffersFromTable := col.Collation != "" && !eqFoldStr(col.Collation, tableCollation)

		if charsetDiffers {
			// When charset differs from table, show CHARACTER SET and always COLLATE.
			b.WriteString(fmt.Sprintf(" CHARACTER SET %s", col.Charset))
			collation := col.Collation
			if collation == "" {
				collation = colCharsetDefault
			}
			if collation != "" {
				b.WriteString(fmt.Sprintf(" COLLATE %s", collation))
			}
		} else if collationNonDefault && collationDiffersFromTable {
			// Collation explicitly set on column (differs from both charset default and table).
			// MySQL shows both CHARACTER SET and COLLATE.
			if col.Charset != "" {
				b.WriteString(fmt.Sprintf(" CHARACTER SET %s", col.Charset))
			}
			b.WriteString(fmt.Sprintf(" COLLATE %s", col.Collation))
		} else if collationNonDefault {
			// Collation inherited from table but non-default for charset.
			// MySQL shows only COLLATE (no CHARACTER SET).
			b.WriteString(fmt.Sprintf(" COLLATE %s", col.Collation))
		}
	}

	// Generated column.
	if col.Generated != nil {
		mode := "VIRTUAL"
		if col.Generated.Stored {
			mode = "STORED"
		}
		b.WriteString(fmt.Sprintf(" GENERATED ALWAYS AS (%s) %s", col.Generated.Expr, mode))
		if !col.Nullable {
			b.WriteString(" NOT NULL")
		}
		if col.Comment != "" {
			b.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeComment(col.Comment)))
		}
		if col.Invisible {
			b.WriteString(" /*!80023 INVISIBLE */")
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

	// DEFAULT.
	if col.Default != nil {
		b.WriteString(" DEFAULT ")
		b.WriteString(formatDefault(*col.Default, col))
	} else if col.Nullable && !col.AutoIncrement && !isTextBlobType(col.DataType) {
		b.WriteString(" DEFAULT NULL")
	}

	// AUTO_INCREMENT.
	if col.AutoIncrement {
		b.WriteString(" AUTO_INCREMENT")
	}

	// ON UPDATE.
	if col.OnUpdate != "" {
		b.WriteString(fmt.Sprintf(" ON UPDATE %s", formatOnUpdate(col.OnUpdate)))
	}

	// COMMENT.
	if col.Comment != "" {
		b.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeComment(col.Comment)))
	}

	// INVISIBLE.
	if col.Invisible {
		b.WriteString(" /*!80023 INVISIBLE */")
	}

	return b.String()
}

// formatDefault formats a default value for SHOW CREATE TABLE output.
// MySQL 8.0 quotes numeric defaults as strings (e.g. DEFAULT '0').
func formatDefault(val string, col *Column) string {
	if strings.EqualFold(val, "NULL") {
		return "NULL"
	}
	// Normalize CURRENT_TIMESTAMP() → CURRENT_TIMESTAMP (MySQL 8.0 format).
	upper := strings.ToUpper(val)
	if upper == "CURRENT_TIMESTAMP" || upper == "CURRENT_TIMESTAMP()" {
		return "CURRENT_TIMESTAMP"
	}
	if strings.HasPrefix(upper, "CURRENT_TIMESTAMP(") {
		// CURRENT_TIMESTAMP(N) — keep precision, use uppercase.
		return upper
	}
	if upper == "NOW()" {
		return "CURRENT_TIMESTAMP"
	}
	// b'...' and 0x... bit/hex literals — not quoted.
	if strings.HasPrefix(val, "b'") || strings.HasPrefix(val, "B'") ||
		strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		return val
	}
	// Expression defaults: (expr) — not quoted, shown as-is.
	if len(val) >= 2 && val[0] == '(' && val[len(val)-1] == ')' {
		return val
	}
	// Already single-quoted string — return as-is.
	if len(val) >= 2 && val[0] == '\'' && val[len(val)-1] == '\'' {
		return val
	}
	// MySQL 8.0 quotes all literal defaults (including numerics).
	return "'" + val + "'"
}

// formatOnUpdate normalizes ON UPDATE values to MySQL 8.0 format.
func formatOnUpdate(val string) string {
	upper := strings.ToUpper(val)
	if upper == "CURRENT_TIMESTAMP" || upper == "CURRENT_TIMESTAMP()" {
		return "CURRENT_TIMESTAMP"
	}
	if strings.HasPrefix(upper, "CURRENT_TIMESTAMP(") {
		return upper
	}
	if upper == "NOW()" {
		return "CURRENT_TIMESTAMP"
	}
	return val
}

// isTimestampType returns true for TIMESTAMP/DATETIME types.
func isTimestampType(dt string) bool {
	switch strings.ToLower(dt) {
	case "timestamp":
		return true
	}
	return false
}

// isTextBlobType returns true for types where MySQL doesn't show DEFAULT NULL.
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
	if !idx.Primary && !idx.Fulltext && !idx.Spatial && idx.IndexType != "" {
		b.WriteString(fmt.Sprintf(" USING %s", strings.ToUpper(idx.IndexType)))
	}

	// KEY_BLOCK_SIZE.
	if idx.KeyBlockSize > 0 {
		b.WriteString(fmt.Sprintf(" KEY_BLOCK_SIZE=%d", idx.KeyBlockSize))
	}

	// Comment.
	if idx.Comment != "" {
		b.WriteString(fmt.Sprintf(" COMMENT '%s'", escapeComment(idx.Comment)))
	}

	// Invisible.
	if !idx.Visible {
		b.WriteString(" /*!80000 INVISIBLE */")
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
		b.WriteString(fmt.Sprintf(") REFERENCES `%s` (", con.RefTable))
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

	if tbl.Engine != "" {
		parts = append(parts, fmt.Sprintf("ENGINE=%s", tbl.Engine))
	}

	// AUTO_INCREMENT — shown only when > 1.
	if tbl.AutoIncrement > 1 {
		parts = append(parts, fmt.Sprintf("AUTO_INCREMENT=%d", tbl.AutoIncrement))
	}

	if tbl.Charset != "" {
		parts = append(parts, fmt.Sprintf("DEFAULT CHARSET=%s", tbl.Charset))
	}

	// MySQL 8.0 shows COLLATE when:
	// - The collation differs from the charset's default, OR
	// - The collation was explicitly specified, OR
	// - The charset is utf8mb4 (MySQL 8.0 always shows collation for utf8mb4)
	if tbl.Charset != "" {
		effectiveCollation := tbl.Collation
		if effectiveCollation == "" {
			effectiveCollation = defaultCollationForCharset[toLower(tbl.Charset)]
		}
		defColl := defaultCollationForCharset[toLower(tbl.Charset)]
		isNonDefaultCollation := effectiveCollation != "" && !eqFoldStr(effectiveCollation, defColl)
		isUtf8mb4 := eqFoldStr(tbl.Charset, "utf8mb4")
		if isNonDefaultCollation || isUtf8mb4 {
			if effectiveCollation != "" {
				parts = append(parts, fmt.Sprintf("COLLATE=%s", effectiveCollation))
			}
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

	return strings.Join(parts, " ")
}

// isFKDefault returns true if the action is a MySQL FK default that should not be shown.
// MySQL 8.0 hides NO ACTION (the implicit default) but shows RESTRICT when explicitly specified.
func isFKDefault(action string) bool {
	upper := strings.ToUpper(action)
	return upper == "NO ACTION"
}

func isEnumSetType(dt string) bool {
	switch strings.ToLower(dt) {
	case "enum", "set":
		return true
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
