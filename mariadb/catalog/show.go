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
		"utf8mb4", "utf8mb3", "utf16le", "gb18030", "geostd8", "armscii8",
		"eucjpms", "cp1252", "gb2312", "euckr", "utf16", "utf32",
		"latin1", "ascii", "binary", "cp932", "tis620", "hebrew",
		"greek", "sjis", "big5", "ucs2", "utf8", "gbk",
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
		if col.Hidden == ColumnHiddenSystem {
			continue
		}
		if col.GeneratedInvisiblePrimaryKey && !c.showGIPK {
			continue
		}
		parts = append(parts, showColumnWithTable(col, tbl))
	}

	// Indexes — MySQL 8.0 orders them in groups:
	// 1. PRIMARY KEY
	// 2. UNIQUE KEYs (creation order)
	// 3. Regular + SPATIAL KEYs, non-expression (creation order)
	// 4. Expression-based KEYs (creation order)
	// 5. FULLTEXT KEYs (creation order)
	var idxPrimary, idxUnique, idxRegular, idxExpr, idxFulltext []*Index
	for _, idx := range tbl.Indexes {
		if isGeneratedInvisiblePrimaryKeyIndex(idx) && !c.showGIPK {
			continue
		}
		switch {
		case idx.Primary:
			idxPrimary = append(idxPrimary, idx)
		case idx.Unique:
			idxUnique = append(idxUnique, idx)
		case idx.Fulltext:
			idxFulltext = append(idxFulltext, idx)
		case isExpressionIndex(idx):
			idxExpr = append(idxExpr, idx)
		default:
			idxRegular = append(idxRegular, idx)
		}
	}
	for _, group := range [][]*Index{idxPrimary, idxUnique, idxRegular, idxExpr, idxFulltext} {
		for _, idx := range group {
			parts = append(parts, showIndex(idx))
		}
	}

	// Constraints (FK and CHECK only — PK/UNIQUE are shown via indexes).
	// MySQL 8.0 sorts FKs alphabetically by name, then CHECKs alphabetically by name.
	var fkConstraints, chkConstraints []*Constraint
	for _, con := range tbl.Constraints {
		switch con.Type {
		case ConForeignKey:
			fkConstraints = append(fkConstraints, con)
		case ConCheck:
			chkConstraints = append(chkConstraints, con)
		}
	}
	sort.Slice(fkConstraints, func(i, j int) bool {
		return fkConstraints[i].Name < fkConstraints[j].Name
	})
	sort.Slice(chkConstraints, func(i, j int) bool {
		return chkConstraints[i].Name < chkConstraints[j].Name
	})
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

	// Partition clause.
	if tbl.Partitioning != nil {
		b.WriteString("\n")
		b.WriteString(showPartitioning(tbl.Partitioning))
	}

	return b.String()
}

func isGeneratedInvisiblePrimaryKeyIndex(idx *Index) bool {
	return idx != nil &&
		idx.Primary &&
		len(idx.Columns) == 1 &&
		strings.EqualFold(idx.Columns[0].Name, generatedInvisiblePrimaryKeyColumnName)
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

	// SRID for spatial types — MySQL 8.0 places it after NOT NULL, before DEFAULT.
	if col.SRID != 0 {
		b.WriteString(fmt.Sprintf(" /*!80003 SRID %d */", col.SRID))
	}

	// DEFAULT.
	if col.Default != nil {
		b.WriteString(" DEFAULT ")
		b.WriteString(formatDefault(*col.Default, col))
	} else if col.Nullable && !col.AutoIncrement && !isTextBlobType(col.DataType) && !col.DefaultDropped {
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
	switch col.DefaultKind {
	case ColumnDefaultCurrentTimestamp:
		return formatCurrentTimestamp(val)
	case ColumnDefaultExpression:
		return "(" + val + ")"
	}
	// Normalize CURRENT_TIMESTAMP() → CURRENT_TIMESTAMP (MySQL 8.0 format).
	if ts, ok := currentTimestampSQL(val); ok {
		return ts
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
func formatOnUpdate(val string, kind ColumnDefaultKind) string {
	if kind == ColumnDefaultCurrentTimestamp {
		return formatCurrentTimestamp(val)
	}
	if ts, ok := currentTimestampSQL(val); ok {
		return ts
	}
	return val
}

func formatCurrentTimestamp(val string) string {
	if ts, ok := currentTimestampSQL(val); ok {
		return ts
	}
	return "CURRENT_TIMESTAMP"
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

	// KEY_BLOCK_SIZE is intentionally NOT rendered here: MySQL 8.0 parses
	// the index-level KEY_BLOCK_SIZE option but does not include it in
	// SHOW CREATE TABLE output (only the table-level KEY_BLOCK_SIZE is
	// shown). The catalog still preserves the value for programmatic
	// access.

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
// MySQL 8.0 outputs partitioning after table options with specific formatting.
func showPartitioning(pi *PartitionInfo) string {
	var b strings.Builder

	// MySQL uses /*!50500 for RANGE COLUMNS and LIST COLUMNS, /*!50100 for others.
	versionComment := "50100"
	if pi.Type == "RANGE COLUMNS" || pi.Type == "LIST COLUMNS" {
		versionComment = "50500"
	}
	b.WriteString(fmt.Sprintf("/*!%s PARTITION BY ", versionComment))
	if pi.Linear {
		b.WriteString("LINEAR ")
	}
	switch pi.Type {
	case "RANGE":
		b.WriteString(fmt.Sprintf("RANGE (%s)", pi.Expr))
	case "RANGE COLUMNS":
		// MySQL uses double space before COLUMNS and no backticks on column names.
		b.WriteString(fmt.Sprintf("RANGE  COLUMNS(%s)", formatPartitionColumnsPlain(pi.Columns)))
	case "LIST":
		b.WriteString(fmt.Sprintf("LIST (%s)", pi.Expr))
	case "LIST COLUMNS":
		b.WriteString(fmt.Sprintf("LIST  COLUMNS(%s)", formatPartitionColumnsPlain(pi.Columns)))
	case "HASH":
		b.WriteString(fmt.Sprintf("HASH (%s)", pi.Expr))
	case "KEY":
		// MySQL does not backtick-quote KEY column names.
		if pi.Algorithm > 0 {
			b.WriteString(fmt.Sprintf("KEY ALGORITHM = %d (%s)", pi.Algorithm, formatPartitionColumnsPlain(pi.Columns)))
		} else {
			b.WriteString(fmt.Sprintf("KEY (%s)", formatPartitionColumnsPlain(pi.Columns)))
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
	hashKeyAutoGen := pi.NumParts > 0 && (pi.Type == "HASH" || pi.Type == "KEY")
	if hashKeyAutoGen {
		b.WriteString(fmt.Sprintf("\nPARTITIONS %d", pi.NumParts))
	} else if len(pi.Partitions) > 0 {
		b.WriteString("\n(")
		for i, pd := range pi.Partitions {
			if i > 0 {
				b.WriteString(",\n ")
			}
			b.WriteString("PARTITION ")
			b.WriteString(pd.Name)
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
				b.WriteString(fmt.Sprintf(" ENGINE = %s", partitionEngine(pd.Engine)))
				if pd.Comment != "" {
					b.WriteString(fmt.Sprintf(" COMMENT = '%s'", escapeComment(pd.Comment)))
				}
			}
			// Explicit subpartition definitions are rendered instead of parent options.
			if len(pd.SubPartitions) > 0 && !pi.UseDefaultSubpartitions {
				b.WriteString("\n (")
				for j, spd := range pd.SubPartitions {
					if j > 0 {
						b.WriteString(",\n  ")
					}
					b.WriteString("SUBPARTITION ")
					b.WriteString(spd.Name)
					b.WriteString(fmt.Sprintf(" ENGINE = %s", partitionEngine(spd.Engine)))
					if spd.Comment != "" {
						b.WriteString(fmt.Sprintf(" COMMENT = '%s'", escapeComment(spd.Comment)))
					}
				}
				b.WriteString(")")
			}
		}
		b.WriteString(")")
	}

	b.WriteString(" */")
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

// formatPartitionColumnsPlain formats column names without backticks (MySQL style for RANGE COLUMNS, LIST COLUMNS, KEY).
func formatPartitionColumnsPlain(cols []string) string {
	return strings.Join(cols, ",")
}

// partitionEngine returns the engine for a partition, defaulting to InnoDB.
func partitionEngine(engine string) string {
	if engine == "" {
		return "InnoDB"
	}
	return engine
}
