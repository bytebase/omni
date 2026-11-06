package catalog

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytebase/omni/tidb/parser"
)

// defaultCollationForCharset maps each charset TiDB has to its default
// collation, as SHOW CHARACTER SET on v8.5.5 lists them (read on 2026-11-06):
// the binary collation for every charset but gbk. A charset missing here is
// one TiDB answers 1115 for.
var defaultCollationForCharset = map[string]string{
	"ascii":   "ascii_bin",
	"binary":  "binary",
	"gbk":     "gbk_chinese_ci",
	"latin1":  "latin1_bin",
	"utf8":    "utf8_bin",
	"utf8mb4": "utf8mb4_bin",
}

// binModifierCollation returns the collation the standalone BINARY modifier
// resolves to for the given charset: {charset}_bin, except the `binary` charset
// whose binary collation is `binary` (not `binary_bin`). Matches TiDB.
func binModifierCollation(charset string) string {
	if toLower(charset) == "binary" {
		return "binary"
	}
	return toLower(charset) + "_bin"
}

// ShowCreateTable produces what TiDB v8.5.5's SHOW CREATE TABLE writes for the
// table: constructResultOfShowCreateTable in pkg/executor/show.go and
// AppendPartitionInfo in pkg/ddl/partition.go, checked against the server
// (tidb/catalog's differential tests). Returns "" if the database or table does
// not exist.
func (c *Catalog) ShowCreateTable(database, table string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	tbl := db.GetTable(table)
	if tbl == nil {
		return ""
	}

	tblCharset := tbl.Charset
	if tblCharset == "" {
		tblCharset = "utf8mb4"
	}
	tblCollate := tbl.Collation
	if tblCollate == "" {
		tblCollate = defaultCollationForCharset[toLower(tblCharset)]
	}

	var b strings.Builder
	switch {
	case tbl.GlobalTemporary:
		fmt.Fprintf(&b, "CREATE GLOBAL TEMPORARY TABLE %s (\n", quoteIdent(tbl.Name))
	case tbl.Temporary:
		fmt.Fprintf(&b, "CREATE TEMPORARY TABLE %s (\n", quoteIdent(tbl.Name))
	default:
		fmt.Fprintf(&b, "CREATE TABLE %s (\n", quoteIdent(tbl.Name))
	}

	parts := make([]string, 0, len(tbl.Columns)+len(tbl.Indexes)+len(tbl.Constraints))
	for _, col := range tbl.Columns {
		parts = append(parts, showColumn(col, tblCharset, tblCollate))
	}

	// A clustered primary key over one integer column is the row handle, which
	// TiDB keeps on the column rather than among the indexes, so it is written
	// right after the columns. Every other key is written in the order the
	// table got it.
	var handle *Index
	for _, idx := range tbl.Indexes {
		if idx.Primary && primaryKeyClustered(tbl, idx) && isRowHandle(tbl, idx) {
			handle = idx
			parts = append(parts, showIndex(tbl, idx))
		}
	}
	for _, idx := range tbl.Indexes {
		if idx == handle || idx.Fulltext || idx.Spatial {
			continue
		}
		parts = append(parts, showIndex(tbl, idx))
	}

	// Foreign keys in the order they were added, then CHECK constraints, which
	// exist only while tidb_enable_check_constraint is on.
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			parts = append(parts, showForeignKey(tbl, con))
		}
	}
	for _, con := range tbl.Constraints {
		if con.Type == ConCheck {
			s := fmt.Sprintf("CONSTRAINT %s CHECK ((%s))", quoteIdent(con.Name), stripOuterParens(con.CheckExpr))
			if con.NotEnforced {
				s += " /*!80016 NOT ENFORCED */"
			}
			parts = append(parts, s)
		}
	}

	b.WriteString("  ")
	b.WriteString(strings.Join(parts, ",\n  "))
	b.WriteString("\n) ")
	b.WriteString(showTableOptions(tbl, tblCharset, tblCollate))

	if tbl.Partitioning != nil {
		b.WriteString(showPartitioning(tbl.Partitioning))
	}
	if tbl.TTLColumn != "" {
		b.WriteString(showTTL(tbl))
	}
	// The ON COMMIT clause of a global temporary table goes last, after the
	// partition clause, because that is the one place the grammar takes it:
	// `... partition by hash(a) partitions 10 on commit delete rows` parses on
	// v8.5.5 and `... on commit delete rows comment 'x'` does not. With no
	// partitions it lands where the server writes it, last on the options
	// line. The server refuses a partitioned temporary table (1562), so that
	// state is one only this catalog reaches; the order still has to parse.
	if tbl.OnCommit != "" {
		b.WriteString(" ON COMMIT ")
		b.WriteString(tbl.OnCommit)
	}
	return b.String()
}

// quoteIdent backquotes a name, doubling any backquote in it.
func quoteIdent(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// outputFormat escapes a string literal's body as v8.5.5's SHOW CREATE TABLE
// does: a quote doubled, a backslash doubled, and NUL, LF and CR as escapes
// (`COMMENT 'it”s a \\test'` for it's a \test, read on 2026-11-06).
func outputFormat(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case 0:
			b.WriteString(`\0`)
		case '\'':
			b.WriteString("''")
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// stripOuterParens removes one pair of parentheses that encloses the whole of
// s, so a caller that adds its own does not double them.
func stripOuterParens(s string) string {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '(' || s[len(s)-1] != ')' {
		return s
	}
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i != len(s)-1 {
				return s
			}
		}
	}
	return s[1 : len(s)-1]
}

func showColumn(col *Column, tblCharset, tblCollate string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s", quoteIdent(col.Name), col.ColumnType)

	// A column with a charset (a string type, which CHARACTER SET binary turns
	// into a binary one except for ENUM and SET) names it when it differs from
	// the table's, and names its collation when that differs from the table's
	// or, matching it, is not its charset's default.
	if isStringType(col.DataType) && (isEnumSetType(col.DataType) || !eqFoldStr(col.Charset, "binary")) {
		if col.Charset != "" && !eqFoldStr(col.Charset, tblCharset) {
			fmt.Fprintf(&b, " CHARACTER SET %s", col.Charset)
		}
		if col.Collation != "" {
			if !eqFoldStr(col.Collation, tblCollate) ||
				!eqFoldStr(col.Collation, defaultCollationForCharset[toLower(col.Charset)]) {
				fmt.Fprintf(&b, " COLLATE %s", col.Collation)
			}
		}
	}

	if col.Generated != nil {
		mode := "VIRTUAL"
		if col.Generated.Stored {
			mode = "STORED"
		}
		fmt.Fprintf(&b, " GENERATED ALWAYS AS (%s) %s", stripOuterParens(col.Generated.Expr), mode)
	}

	if col.AutoIncrement {
		b.WriteString(" NOT NULL AUTO_INCREMENT")
	} else {
		if !col.Nullable {
			b.WriteString(" NOT NULL")
		}
		// A generated column has no default, and a NOT NULL column without one
		// writes none.
		if col.Generated == nil {
			switch {
			case col.Default != nil && !strings.EqualFold(*col.Default, "NULL"):
				b.WriteString(" DEFAULT ")
				b.WriteString(formatDefault(*col.Default, col))
			case col.Nullable && !col.DefaultDropped:
				if isTimestampType(col.DataType) {
					b.WriteString(" NULL")
				}
				b.WriteString(" DEFAULT NULL")
			}
		}
		if col.OnUpdate != "" {
			fmt.Fprintf(&b, " ON UPDATE %s", formatOnUpdate(col.OnUpdate, col))
		}
	}

	if col.AutoRandom {
		shard := col.AutoRandomShardBits
		if shard == 0 {
			shard = 5
		}
		if r := col.AutoRandomRangeBits; r == 0 || r == 64 {
			fmt.Fprintf(&b, " /*T![auto_rand] AUTO_RANDOM(%d) */", shard)
		} else {
			fmt.Fprintf(&b, " /*T![auto_rand] AUTO_RANDOM(%d, %d) */", shard, r)
		}
	}

	if col.Comment != "" {
		fmt.Fprintf(&b, " COMMENT '%s'", outputFormat(col.Comment))
	}
	return b.String()
}

// columnFsp returns the fractional-seconds precision of a temporal column's
// type, as in datetime(3).
func columnFsp(col *Column) int {
	ct := col.ColumnType
	open := strings.IndexByte(ct, '(')
	if open < 0 {
		return 0
	}
	n := 0
	for _, r := range ct[open+1:] {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return n
}

// formatDefault writes a column default as TiDB does: CURRENT_TIMESTAMP with
// the column's precision, CURRENT_DATE and other expressions in parentheses, a
// bit literal as written, and every other value quoted.
func formatDefault(val string, col *Column) string {
	upper := strings.ToUpper(strings.TrimSpace(val))
	switch {
	case isCurrentTimestamp(upper):
		if fsp := columnFsp(col); fsp > 0 && isTemporalWithFsp(col.DataType) {
			return fmt.Sprintf("CURRENT_TIMESTAMP(%d)", fsp)
		}
		return "CURRENT_TIMESTAMP"
	case upper == "CURRENT_DATE" || upper == "CURRENT_DATE()" || upper == "CURDATE()" ||
		upper == "(CURRENT_DATE)" || upper == "(CURRENT_DATE())" || upper == "(CURDATE())":
		return "(CURRENT_DATE)"
	}
	if strings.EqualFold(col.DataType, "bit") {
		if lit, ok := bitDefaultLiteral(val); ok {
			return lit
		}
	} else if b, ok := binaryLiteralBytes(val); ok {
		// A hex or bit literal is the bytes it spells: a string column takes
		// them as its text and a numeric column as an unsigned integer, so
		// VARCHAR DEFAULT 0x41 is 'A' and INT DEFAULT 0x41 is '65'
		// (getDefaultValue in TiDB's pkg/ddl/add_column.go). BINARY and
		// VARBINARY are the string field types with the binary charset there,
		// so they take the bytes too.
		if isStringType(col.DataType) || eqFoldStr(col.DataType, "binary") || eqFoldStr(col.DataType, "varbinary") {
			return "'" + outputFormat(string(b)) + "'"
		}
		var n uint64
		for _, c := range b {
			n = n<<8 | uint64(c)
		}
		return "'" + strconv.FormatUint(n, 10) + "'"
	}
	// An expression default is written in parentheses.
	if len(val) >= 2 && val[0] == '(' && val[len(val)-1] == ')' {
		return val
	}
	if len(val) >= 2 && val[0] == '\'' && val[len(val)-1] == '\'' {
		return val
	}
	return "'" + outputFormat(val) + "'"
}

// binaryLiteralBytes decodes a hex literal (0x41, X'41') or a bit literal
// (0b1000001, b'1000001') into the bytes it spells.
func binaryLiteralBytes(v string) ([]byte, bool) {
	v = strings.TrimSpace(v)
	var digits string
	base := 16
	switch {
	case len(v) > 2 && (v[:2] == "0x" || v[:2] == "0X"):
		digits = v[2:]
	case len(v) >= 3 && (v[0] == 'x' || v[0] == 'X') && v[1] == '\'' && v[len(v)-1] == '\'':
		digits = v[2 : len(v)-1]
	case len(v) > 2 && (v[:2] == "0b" || v[:2] == "0B"):
		digits, base = v[2:], 2
	case len(v) >= 3 && (v[0] == 'b' || v[0] == 'B') && v[1] == '\'' && v[len(v)-1] == '\'':
		digits, base = v[2:len(v)-1], 2
	default:
		return nil, false
	}
	var out []byte
	if base == 16 {
		if len(digits)%2 == 1 {
			digits = "0" + digits
		}
		for i := 0; i < len(digits); i += 2 {
			n, err := strconv.ParseUint(digits[i:i+2], 16, 8)
			if err != nil {
				return nil, false
			}
			out = append(out, byte(n))
		}
		return out, true
	}
	for len(digits)%8 != 0 {
		digits = "0" + digits
	}
	for i := 0; i < len(digits); i += 8 {
		n, err := strconv.ParseUint(digits[i:i+8], 2, 8)
		if err != nil {
			return nil, false
		}
		out = append(out, byte(n))
	}
	return out, true
}

// bitDefaultLiteral writes a BIT column's default as TiDB does, the value's
// binary digits with no leading zeros: b'1111' for 15, 0x0f and b'00001111'.
func bitDefaultLiteral(val string) (string, bool) {
	v := strings.TrimSpace(val)
	var n uint64
	var err error
	switch {
	case len(v) > 3 && (v[0] == 'b' || v[0] == 'B') && v[1] == '\'':
		n, err = strconv.ParseUint(strings.TrimSuffix(v[2:], "'"), 2, 64)
	case strings.HasPrefix(v, "0b") || strings.HasPrefix(v, "0B"):
		n, err = strconv.ParseUint(v[2:], 2, 64)
	case strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X"):
		n, err = strconv.ParseUint(v[2:], 16, 64)
	case len(v) > 3 && (v[0] == 'x' || v[0] == 'X') && v[1] == '\'':
		n, err = strconv.ParseUint(strings.TrimSuffix(v[2:], "'"), 16, 64)
	default:
		n, err = strconv.ParseUint(strings.Trim(v, "'"), 10, 64)
	}
	if err != nil {
		return "", false
	}
	return "b'" + strconv.FormatUint(n, 2) + "'", true
}

func isCurrentTimestamp(upper string) bool {
	for _, name := range []string{"CURRENT_TIMESTAMP", "NOW", "LOCALTIME", "LOCALTIMESTAMP"} {
		if upper == name || upper == name+"()" || strings.HasPrefix(upper, name+"(") {
			return true
		}
	}
	return false
}

func isTemporalWithFsp(dt string) bool {
	switch toLower(dt) {
	case "datetime", "timestamp", "time":
		return true
	}
	return false
}

// formatOnUpdate writes ON UPDATE CURRENT_TIMESTAMP with the column's
// precision, which is the one value TiDB takes there.
func formatOnUpdate(val string, col *Column) string {
	if isCurrentTimestamp(strings.ToUpper(strings.TrimSpace(val))) {
		if fsp := columnFsp(col); fsp > 0 {
			return fmt.Sprintf("CURRENT_TIMESTAMP(%d)", fsp)
		}
		return "CURRENT_TIMESTAMP"
	}
	return val
}

func isTimestampType(dt string) bool {
	return strings.EqualFold(dt, "timestamp")
}

// showIndex writes one key. TiDB writes no DESC, no USING clause and no
// KEY_BLOCK_SIZE, and marks every primary key CLUSTERED or NONCLUSTERED.
func showIndex(tbl *Table, idx *Index) string {
	var b strings.Builder
	switch {
	case idx.Primary:
		b.WriteString("PRIMARY KEY ")
	case idx.Unique:
		fmt.Fprintf(&b, "UNIQUE KEY %s ", quoteIdent(idx.Name))
	case idx.Vector:
		fmt.Fprintf(&b, "VECTOR INDEX %s", quoteIdent(idx.Name))
	default:
		fmt.Fprintf(&b, "KEY %s ", quoteIdent(idx.Name))
	}

	cols := make([]string, 0, len(idx.Columns))
	for _, ic := range idx.Columns {
		if ic.Expr != "" {
			expr := stripOuterParens(ic.Expr)
			if idx.Vector {
				expr = upperFuncName(expr)
			}
			cols = append(cols, "("+expr+")")
			continue
		}
		s := quoteIdent(ic.Name)
		if ic.Length > 0 {
			s += fmt.Sprintf("(%d)", ic.Length)
		}
		cols = append(cols, s)
	}
	// A vector index's one key part is its distance function, which TiDB
	// writes right after the name and in capitals, with no USING clause:
	// VECTOR INDEX `i`((VEC_COSINE_DISTANCE(`v`))) (constructResultOfShowCreateTable
	// in pkg/executor/show.go; the test server has no TiFlash and refuses
	// every vector index with 8200, so this one is read from the source).
	fmt.Fprintf(&b, "(%s)", strings.Join(cols, ","))

	if !idx.Visible {
		b.WriteString(" /*!80000 INVISIBLE */")
	}
	if idx.Comment != "" {
		fmt.Fprintf(&b, " COMMENT '%s'", outputFormat(idx.Comment))
	}
	if idx.Primary {
		if primaryKeyClustered(tbl, idx) {
			b.WriteString(" /*T![clustered_index] CLUSTERED */")
		} else {
			b.WriteString(" /*T![clustered_index] NONCLUSTERED */")
		}
	}
	if idx.Global {
		b.WriteString(" /*T![global_index] GLOBAL */")
	}
	return b.String()
}

// showForeignKey writes a foreign key as TiDB does: columns joined without
// spaces, the referenced schema only when it is not the table's, and ON DELETE
// and ON UPDATE whenever the statement gave them, NO ACTION included.
func showForeignKey(tbl *Table, con *Constraint) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CONSTRAINT %s FOREIGN KEY ", quoteIdent(con.Name))
	b.WriteString("(" + joinQuoted(con.Columns) + ")")
	if con.RefDatabase != "" && (tbl.Database == nil || !eqFoldStr(con.RefDatabase, tbl.Database.Name)) {
		fmt.Fprintf(&b, " REFERENCES %s.%s ", quoteIdent(con.RefDatabase), quoteIdent(con.RefTable))
	} else {
		fmt.Fprintf(&b, " REFERENCES %s ", quoteIdent(con.RefTable))
	}
	b.WriteString("(" + joinQuoted(con.RefColumns) + ")")
	if con.OnDelete != "" {
		fmt.Fprintf(&b, " ON DELETE %s", strings.ToUpper(con.OnDelete))
	}
	if con.OnUpdate != "" {
		fmt.Fprintf(&b, " ON UPDATE %s", strings.ToUpper(con.OnUpdate))
	}
	return b.String()
}

// upperFuncName capitalizes the function name a call starts with.
func upperFuncName(expr string) string {
	if i := strings.IndexByte(expr, '('); i > 0 {
		return strings.ToUpper(expr[:i]) + expr[i:]
	}
	return expr
}

func joinQuoted(names []string) string {
	q := make([]string, len(names))
	for i, n := range names {
		q[i] = quoteIdent(n)
	}
	return strings.Join(q, ",")
}

// showTableOptions writes the options TiDB keeps. ENGINE is always InnoDB, and
// ROW_FORMAT, KEY_BLOCK_SIZE and the STATS_ options are read and dropped.
// AUTO_INCREMENT shows only on a table with an AUTO_INCREMENT column.
func showTableOptions(tbl *Table, charset, collate string) string {
	var b strings.Builder
	b.WriteString("ENGINE=InnoDB")
	if collate == "" || eqFoldStr(collate, "binary") {
		fmt.Fprintf(&b, " DEFAULT CHARSET=%s", charset)
	} else {
		fmt.Fprintf(&b, " DEFAULT CHARSET=%s COLLATE=%s", charset, collate)
	}
	if tbl.AutoIncrement > 1 {
		for _, col := range tbl.Columns {
			if col.AutoIncrement {
				fmt.Fprintf(&b, " AUTO_INCREMENT=%d", tbl.AutoIncrement)
				break
			}
		}
	}
	if tbl.AutoIDCache != 0 {
		fmt.Fprintf(&b, " /*T![auto_id_cache] AUTO_ID_CACHE=%d */", tbl.AutoIDCache)
	}
	if tbl.AutoRandomBase > 1 {
		fmt.Fprintf(&b, " /*T![auto_rand_base] AUTO_RANDOM_BASE=%d */", tbl.AutoRandomBase)
	}
	if tbl.ShardRowIDBits > 0 {
		fmt.Fprintf(&b, " /*T! SHARD_ROW_ID_BITS=%d ", tbl.ShardRowIDBits)
		if tbl.PreSplitRegions > 0 {
			fmt.Fprintf(&b, "PRE_SPLIT_REGIONS=%d ", tbl.PreSplitRegions)
		}
		b.WriteString("*/")
	}
	if tbl.PreSplitRegions > 0 && tbl.ShardRowIDBits == 0 && hasAutoRandom(tbl) {
		fmt.Fprintf(&b, " /*T! PRE_SPLIT_REGIONS=%d */", tbl.PreSplitRegions)
	}
	if tbl.Comment != "" {
		fmt.Fprintf(&b, " COMMENT='%s'", outputFormat(tbl.Comment))
	}
	if tbl.PlacementPolicy != "" {
		fmt.Fprintf(&b, " /*T![placement] PLACEMENT POLICY=%s */", quoteIdent(tbl.PlacementPolicy))
	}
	return b.String()
}

func hasAutoRandom(tbl *Table) bool {
	for _, col := range tbl.Columns {
		if col.AutoRandom {
			return true
		}
	}
	return false
}

// showTTL writes the three TTL options, each in its own version comment, with
// TiDB's defaults when the statement left one out: TTL_ENABLE 'ON' and
// TTL_JOB_INTERVAL '1h'.
func showTTL(tbl *Table) string {
	enable := "ON"
	if !tbl.TTLEnable {
		enable = "OFF"
	}
	job := tbl.TTLJobInterval
	if job == "" {
		job = "1h"
	}
	return fmt.Sprintf(" /*T![ttl] TTL=%s + %s */ /*T![ttl] TTL_ENABLE='%s' */ /*T![ttl] TTL_JOB_INTERVAL='%s' */",
		quoteIdent(tbl.TTLColumn), tbl.TTLInterval, enable, job)
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

// escapeComment escapes a comment for the MySQL-shaped SHOW CREATE output of
// objects TiDB does not have (events, routines).
func escapeComment(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "''")
	return s
}

// showPartitioning writes the partition clause as AppendPartitionInfo does: no
// version comment, HASH and KEY partitions named p0, p1, ... as a count, and
// every other definition list in full. TiDB keeps no LINEAR, no KEY ALGORITHM,
// no partition ENGINE and no subpartitions: v8.5.5 reads each and writes
// `PARTITION BY HASH (a)` for `LINEAR HASH (a)` (read on 2026-11-06).
func showPartitioning(pi *PartitionInfo) string {
	var b strings.Builder
	typ := pi.Type

	if pi.Type == "HASH" || pi.Type == "KEY" {
		defaults := true
		for i, pd := range pi.Partitions {
			if pd.Name != fmt.Sprintf("p%d", i) || pd.Comment != "" {
				defaults = false
				break
			}
		}
		if defaults {
			n := pi.NumParts
			if n == 0 {
				n = len(pi.Partitions)
			}
			if n == 0 {
				n = 1
			}
			if pi.Type == "HASH" {
				fmt.Fprintf(&b, "\nPARTITION BY %s (%s) PARTITIONS %d", typ, pi.Expr, n)
			} else {
				fmt.Fprintf(&b, "\nPARTITION BY %s (%s) PARTITIONS %d", typ, joinQuoted(pi.Columns), n)
			}
			return b.String()
		}
	}

	switch pi.Type {
	case "RANGE COLUMNS", "LIST COLUMNS":
		fmt.Fprintf(&b, "\nPARTITION BY %s(%s)", typ, joinQuoted(pi.Columns))
	case "KEY":
		fmt.Fprintf(&b, "\nPARTITION BY %s (%s)", typ, joinQuoted(pi.Columns))
	default:
		fmt.Fprintf(&b, "\nPARTITION BY %s (%s)", typ, pi.Expr)
	}
	// An INTERVAL clause the catalog could not expand stays as written, so the
	// text still parses (partition_interval.go); TiDB itself prints the
	// partitions it generated.
	if pi.IntervalText != "" {
		b.WriteString(" " + pi.IntervalText)
		if len(pi.Partitions) == 0 {
			return b.String()
		}
	}
	b.WriteString("\n(")
	for i, pd := range pi.Partitions {
		if i > 0 {
			b.WriteString(",\n ")
		}
		fmt.Fprintf(&b, "PARTITION %s", quoteIdent(pd.Name))
		// The catalog keeps a partition's values bare, joined by bare commas
		// (partitionValueToString), so each list takes one pair of
		// parentheses here and a LIST COLUMNS tuple keeps its own.
		switch {
		case strings.HasPrefix(pi.Type, "RANGE"):
			fmt.Fprintf(&b, " VALUES LESS THAN (%s)", strings.TrimSpace(pd.ValueExpr))
		case strings.HasPrefix(pi.Type, "LIST"):
			if v := strings.TrimSpace(pd.ValueExpr); v == "" || strings.EqualFold(v, "DEFAULT") {
				b.WriteString(" DEFAULT")
			} else {
				fmt.Fprintf(&b, " VALUES IN (%s)", v)
			}
		}
		if pd.Comment != "" {
			fmt.Fprintf(&b, " COMMENT '%s'", outputFormat(pd.Comment))
		}
	}
	b.WriteString(")")
	return b.String()
}

// partitionNameInDDL quotes a partition name that needs it: one with other characters, as a
// generated P_LT_2020-01-01 00:00:00 has, and a reserved word, such as ilike.
func partitionNameInDDL(name string) string {
	if parser.IsReservedKeyword(name) {
		return "`" + name + "`"
	}
	for _, r := range name {
		if r != '_' && (r < '0' || r > '9') && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return "`" + name + "`"
		}
	}
	return name
}
