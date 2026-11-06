package catalog

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/deparse"
)

// What TiDB v8.5.5 does with RANGE ... INTERVAL when it runs the statement, read from its answers
// and its SHOW CREATE TABLE on 2026-11-05. It keeps no interval: it generates ordinary partitions,
// one per step from FIRST to LAST, so that is what the catalog holds.
//
// A bound or a step this package cannot work out without the server, such as a function call,
// generates nothing. The catalog then keeps the clause as written (PartitionInfo.IntervalText), so
// SHOW CREATE TABLE stays text the server takes, though not the text it would print.

// maxIntervalPartitions is TiDB's limit on a table's partitions, past which it answers 1499.
const maxIntervalPartitions = 8192

func errIntervalUnsupported(format string, args ...any) error {
	return &Error{Code: 8200, SQLState: "HY000", Message: fmt.Sprintf(format, args...)}
}

// applyPartitionInterval fills in the partitions an INTERVAL clause generates. Written partitions
// must be the generated ones in number and in bound, under names of their own.
func applyPartitionInterval(tbl *Table, pc *nodes.PartitionClause, pi *PartitionInfo) error {
	iv := pc.Interval
	if iv.First == nil || iv.Last == nil {
		return errIntervalUnsupported("Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined")
	}
	if len(pc.Columns) > 1 {
		return errIntervalUnsupported("Unsupported INTERVAL partitioning, does not allow RANGE COLUMNS with more than one column")
	}
	var col *Column
	if len(pc.Columns) == 1 {
		if col = tbl.GetColumn(pc.Columns[0]); col == nil {
			return &Error{Code: 1488, SQLState: "HY000", Message: "Field in list of fields for partition function not found in table"}
		}
	}
	gen, known, err := generateIntervalPartitions(iv, col, tbl, pc)
	if err != nil {
		return err
	}
	// The server keeps neither of these: PARTITIONS 9 and SUBPARTITION BY are read and dropped.
	pi.NumParts = 0
	pi.SubType, pi.SubExpr, pi.SubColumns, pi.SubAlgo, pi.NumSubParts = "", "", nil, 0, 0
	if !known {
		pi.IntervalText = writtenPartitionInterval(iv)
		// The server generates the partitions, so the empty ones PARTITIONS n made are none of
		// them; written ones stay, to be read back beside the clause.
		if len(pc.Partitions) == 0 {
			pi.Partitions = nil
		}
		return nil
	}
	if len(pc.Partitions) > 0 {
		if len(pi.Partitions) != len(gen) {
			return errIntervalUnsupported("Unsupported number of partitions generated != partition defined (%d != %d)", len(gen), len(pi.Partitions))
		}
		for i, pd := range pi.Partitions {
			if pd.ValueExpr != gen[i].ValueExpr {
				return errIntervalUnsupported("Unsupported partition %s differs between generated and defined for expression", pd.Name)
			}
		}
	} else {
		pi.Partitions = gen
	}
	return nil
}

// generateIntervalPartitions returns the partitions iv generates for a RANGE COLUMNS column, or
// for a RANGE expression when col is nil. known is false when it cannot work them out.
func generateIntervalPartitions(iv *nodes.PartitionInterval, col *Column, tbl *Table, pc *nodes.PartitionClause) (gen []*PartitionDefInfo, known bool, err error) {
	kind := "int"
	if col != nil {
		switch strings.ToLower(col.DataType) {
		case "date":
			kind = "date"
		case "datetime":
			kind = "datetime"
		case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		default:
			return nil, false, errIntervalUnsupported("Unsupported INTERVAL partitioning, only supports Date, Datetime and INT types")
		}
	}
	var bounds []string
	if kind == "int" {
		bounds, known, err = intIntervalBounds(iv)
	} else {
		bounds, known, err = timeIntervalBounds(iv, kind)
	}
	if err != nil || !known {
		return nil, known, err
	}
	// NULL PARTITION and MAXVALUE PARTITION count toward the limit: 8191 bounds and both are 1499.
	total := len(bounds)
	if iv.NullPart {
		total++
	}
	if iv.MaxValPart {
		total++
	}
	if total > maxIntervalPartitions {
		return nil, false, &Error{Code: 1499, SQLState: "HY000", Message: "Too many partitions (including subpartitions) were defined"}
	}
	seen := make(map[string]bool, len(bounds))
	if iv.NullPart {
		gen = append(gen, &PartitionDefInfo{Name: "P_NULL", ValueExpr: intervalNullBound(col, kind, tbl, pc)})
	}
	for _, bound := range bounds {
		name := "P_LT_" + strings.Trim(bound, "'")
		if seen[name] {
			return nil, false, &Error{Code: 1517, SQLState: "HY000", Message: "Duplicate partition name " + name}
		}
		seen[name] = true
		gen = append(gen, &PartitionDefInfo{Name: name, ValueExpr: bound})
	}
	if iv.MaxValPart {
		gen = append(gen, &PartitionDefInfo{Name: "P_MAXVALUE", ValueExpr: "MAXVALUE"})
	}
	return gen, true, nil
}

// intIntervalBounds works out the bounds of an integer INTERVAL.
func intIntervalBounds(iv *nodes.PartitionInterval) ([]string, bool, error) {
	if iv.Unit != "" {
		return nil, false, &Error{Code: 1292, SQLState: "22007", Message: "Incorrect time value"}
	}
	if _, ok := iv.Expr.(*nodes.StringLit); ok {
		return nil, false, errIntervalUnsupported("Unsupported INTERVAL, should be a positive number")
	}
	for _, bound := range []nodes.ExprNode{iv.First, iv.Last} {
		if _, ok := bound.(*nodes.StringLit); ok {
			return nil, false, &Error{Code: 1697, SQLState: "HY000", Message: "VALUES value for partition must have type INT"}
		}
	}
	if !intervalStepWritten(iv.Expr) {
		return nil, false, errIntervalUnsupported("Unsupported INTERVAL, should be a positive number")
	}
	var vals [3]int64
	for i, e := range []nodes.ExprNode{iv.Expr, iv.First, iv.Last} {
		v, ok, err := constantInt(e)
		if err != nil {
			return nil, false, err
		}
		if !ok {
			return nil, false, nil
		}
		vals[i] = v
	}
	step, first, last := vals[0], vals[1], vals[2]
	if step <= 0 {
		return nil, false, errIntervalUnsupported("Unsupported INTERVAL, should be a positive number")
	}
	// Each bound is FIRST + i * INTERVAL, worked out the way v8.5.5 does: in BIGINT, so a step
	// past it answers 1690 even where LAST is off the steps, and up to LAST, so FIRST and LAST at
	// the top of BIGINT make one partition without stepping past it.
	var bounds []string
	for i := int64(0); ; i++ {
		if i >= maxIntervalPartitions {
			return nil, false, &Error{Code: 1499, SQLState: "HY000", Message: "Too many partitions (including subpartitions) were defined"}
		}
		offset, ok := mulInt64(i, step)
		if !ok {
			return nil, false, errBigintOutOfRange(fmt.Sprintf("(%d * %d)", i, step))
		}
		v, ok := addInt64(first, offset)
		if !ok {
			return nil, false, errBigintOutOfRange(fmt.Sprintf("(%d + %d)", first, offset))
		}
		if v > last {
			return nil, false, errIntervalUnsupported("Unsupported INTERVAL: expr (%d) not matching FIRST + n INTERVALs (%d + n * %d)", last, first, step)
		}
		bounds = append(bounds, strconv.FormatInt(v, 10))
		if v == last {
			return bounds, true, nil
		}
	}
}

func errBigintOutOfRange(expr string) error {
	return &Error{Code: 1690, SQLState: "22003", Message: "BIGINT value is out of range in '" + expr + "'"}
}

// mulInt64 and addInt64 report whether the result fits in an int64.
func mulInt64(a, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	c := a * b
	if c/b != a || (a == -1 && b == math.MinInt64) || (b == -1 && a == math.MinInt64) {
		return 0, false
	}
	return c, true
}

func addInt64(a, b int64) (int64, bool) {
	c := a + b
	if (b > 0 && c < a) || (b < 0 && c > a) {
		return 0, false
	}
	return c, true
}

// timeIntervalBounds works out the bounds of a DATE or DATETIME INTERVAL. Each is FIRST plus n
// steps, counted from FIRST, so 2020-01-31 by one month is 2020-02-29 and then 2020-03-31.
func timeIntervalBounds(iv *nodes.PartitionInterval, kind string) ([]string, bool, error) {
	if iv.Unit == "" {
		return nil, false, &Error{Code: 1292, SQLState: "22007", Message: "Truncated incorrect DOUBLE value"}
	}
	var months, seconds int64
	switch iv.Unit {
	case "YEAR":
		months = 12
	case "QUARTER":
		months = 3
	case "MONTH":
		months = 1
	case "WEEK":
		seconds = 7 * 24 * 3600
	case "DAY":
		seconds = 24 * 3600
	case "HOUR":
		seconds = 3600
	case "MINUTE":
		seconds = 60
	case "SECOND":
		seconds = 1
	default:
		return nil, false, errIntervalUnsupported("Unsupported INTERVAL partitioning, only supports YEAR, QUARTER, MONTH, WEEK, DAY, HOUR, MINUTE and SECOND as time unit")
	}
	if !intervalStepWritten(iv.Expr) {
		return nil, false, errIntervalUnsupported("Unsupported INTERVAL, should be a positive number")
	}
	step, ok, err := constantInt(iv.Expr)
	if err != nil || !ok {
		return nil, false, err
	}
	if step <= 0 {
		return nil, false, errIntervalUnsupported("Unsupported INTERVAL, should be a positive number")
	}
	first, ok1 := literalTime(iv.First)
	last, ok2 := literalTime(iv.Last)
	if !ok1 || !ok2 {
		return nil, false, nil
	}
	layout := "2006-01-02 15:04:05"
	if kind == "date" {
		layout = "2006-01-02"
	}
	var bounds []string
	for n := int64(0); ; n++ {
		if n >= maxIntervalPartitions {
			return nil, false, &Error{Code: 1499, SQLState: "HY000", Message: "Too many partitions (including subpartitions) were defined"}
		}
		var at time.Time
		if months > 0 {
			at, ok = addMonthsClamped(first, n, step, months)
		} else {
			at, ok = addSeconds(first, n, step, seconds)
		}
		if !ok {
			// A bound past 9999-12-31 is 1441 on v8.5.5, however far past.
			return nil, false, &Error{Code: 1441, SQLState: "22008", Message: "Datetime function: datetime field overflow"}
		}
		if at.After(last) {
			return nil, false, errIntervalUnsupported("Unsupported INTERVAL: expr (%s) not matching FIRST + n INTERVALs", last.Format(layout))
		}
		bounds = append(bounds, "'"+at.Format(layout)+"'")
		if at.Equal(last) {
			return bounds, true, nil
		}
	}
}

// maxIntervalTime is the first instant past the DATETIME range.
var maxIntervalTime = time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)

// addSeconds is t plus n steps of step units of unit seconds each, counted in whole seconds rather
// than a time.Duration, which stops at about 292 years. It reports false past the DATETIME range.
func addSeconds(t time.Time, n, step, unit int64) (time.Time, bool) {
	offset, ok := mulInt64(n, step)
	if ok {
		offset, ok = mulInt64(offset, unit)
	}
	if !ok {
		return time.Time{}, false
	}
	sec, ok := addInt64(t.Unix(), offset)
	if !ok || sec >= maxIntervalTime.Unix() {
		return time.Time{}, false
	}
	return time.Unix(sec, int64(t.Nanosecond())).UTC(), true
}

// addMonthsClamped is t plus n steps of step units of unit months each, the way DATE_ADD adds
// them, keeping to the last day of a shorter month. It reports false past the DATETIME range.
func addMonthsClamped(t time.Time, n, step, unit int64) (time.Time, bool) {
	months, ok := mulInt64(n, step)
	if ok {
		months, ok = mulInt64(months, unit)
	}
	if !ok {
		return time.Time{}, false
	}
	total, ok := addInt64(int64(t.Year())*12+int64(t.Month())-1, months)
	if !ok || total >= 10000*12 {
		return time.Time{}, false
	}
	year, month := int(total/12), time.Month(total%12+1)
	day := t.Day()
	if last := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day(); day > last {
		day = last
	}
	return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC), true
}

// intervalNullBound is the bound of NULL PARTITION: the lowest value of a RANGE COLUMNS column's
// type, and for a RANGE expression the lowest bigint, or 0 over an unsigned column. RANGE (a) over
// INT is -9223372036854775808 where RANGE COLUMNS (a) is -2147483648.
func intervalNullBound(col *Column, kind string, tbl *Table, pc *nodes.PartitionClause) string {
	switch kind {
	case "date":
		return "'0000-01-01'"
	case "datetime":
		return "'0000-01-01 00:00:00'"
	}
	if col == nil {
		if ref, ok := pc.Expr.(*nodes.ColumnRef); ok {
			if c := tbl.GetColumn(ref.Column); c != nil && strings.Contains(strings.ToLower(c.ColumnType), "unsigned") {
				return "0"
			}
		}
		return "-9223372036854775808"
	}
	if strings.Contains(strings.ToLower(col.ColumnType), "unsigned") {
		return "0"
	}
	switch strings.ToLower(col.DataType) {
	case "tinyint":
		return "-128"
	case "smallint":
		return "-32768"
	case "mediumint":
		return "-8388608"
	case "int", "integer":
		return "-2147483648"
	}
	return "-9223372036854775808"
}

// intervalStepWritten reports whether an INTERVAL step passes the check TiDB
// makes on its text before evaluating it (pkg/ddl/partition.go,
// GeneratePartDefsFromInterval): restored, and unquoted, it must start with a
// digit 1 to 9. So INTERVAL (5+5) is a step of 10, while (+10), ((10)) and
// (-1) answer 8200 (verified on v8.5.5).
func intervalStepWritten(e nodes.ExprNode) bool {
	text := deparse.Restore(e)
	if len(text) >= 2 && text[0] == '\'' && text[len(text)-1] == '\'' {
		text = text[1 : len(text)-1]
	}
	return text != "" && text[0] >= '1' && text[0] <= '9'
}

// constantInt evaluates integer literals under signs, parentheses, +, - and *, in BIGINT: a sum or
// product past it is 1690 on v8.5.5, naming the operation, as INTERVAL (9223372036854775807 * 2 + 3)
// answers out of range in '(9223372036854775807 * 2)'. A minus sign wraps, as the lowest BIGINT
// negated is no positive step there either (8200).
func constantInt(e nodes.ExprNode) (int64, bool, error) {
	switch n := e.(type) {
	case *nodes.IntLit:
		return n.Value, true, nil
	case *nodes.ParenExpr:
		return constantInt(n.Expr)
	case *nodes.UnaryExpr:
		v, ok, err := constantInt(n.Operand)
		if err != nil || !ok {
			return 0, false, err
		}
		switch n.Op {
		case nodes.UnaryMinus:
			return -v, true, nil
		case nodes.UnaryPlus:
			return v, true, nil
		}
	case *nodes.BinaryExpr:
		l, ok1, err := constantInt(n.Left)
		if err != nil {
			return 0, false, err
		}
		r, ok2, err := constantInt(n.Right)
		if err != nil || !ok1 || !ok2 {
			return 0, false, err
		}
		var v int64
		var ok bool
		var op string
		switch n.Op {
		case nodes.BinOpAdd:
			v, ok = addInt64(l, r)
			op = "+"
		case nodes.BinOpSub:
			v, ok = addInt64(l, -r)
			if r == math.MinInt64 {
				ok = false
			}
			op = "-"
		case nodes.BinOpMul:
			v, ok = mulInt64(l, r)
			op = "*"
		default:
			return 0, false, nil
		}
		if !ok {
			return 0, false, errBigintOutOfRange(fmt.Sprintf("(%d %s %d)", l, op, r))
		}
		return v, true, nil
	}
	return 0, false, nil
}

// literalTime reads a string literal as a date or a datetime.
func literalTime(e nodes.ExprNode) (time.Time, bool) {
	lit, ok := e.(*nodes.StringLit)
	if !ok {
		return time.Time{}, false
	}
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, lit.Value); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// writtenPartitionInterval is the clause as DDL, for an interval the catalog cannot expand.
func writtenPartitionInterval(iv *nodes.PartitionInterval) string {
	var b strings.Builder
	b.WriteString("INTERVAL (" + nodeToSQL(iv.Expr))
	if iv.Unit != "" {
		b.WriteString(" " + iv.Unit)
	}
	b.WriteString(")")
	b.WriteString(" FIRST PARTITION LESS THAN (" + nodeToSQL(iv.First) + ")")
	b.WriteString(" LAST PARTITION LESS THAN (" + nodeToSQL(iv.Last) + ")")
	if iv.NullPart {
		b.WriteString(" NULL PARTITION")
	}
	if iv.MaxValPart {
		b.WriteString(" MAXVALUE PARTITION")
	}
	return b.String()
}

// refuseHistoryPartitions refuses what v8.5.5 refuses for HISTORY and CURRENT on a table that is
// not SYSTEM_TIME partitioned, which none in the catalog is: ADD PARTITION answers 1480 on HASH and
// KEY and 4113 on RANGE and LIST, and REORGANIZE PARTITION answers 4113.
func refuseHistoryPartitions(tbl *Table, defs []*nodes.PartitionDef, adding bool) error {
	for _, pd := range defs {
		value, ok := pd.Values.(*nodes.String)
		if !ok || (value.Str != "HISTORY" && value.Str != "CURRENT") {
			continue
		}
		if adding && (tbl.Partitioning.Type == "HASH" || tbl.Partitioning.Type == "KEY") {
			return &Error{Code: 1480, SQLState: "HY000", Message: "Only SYSTEM_TIME PARTITIONING can use VALUES HISTORY in partition definition"}
		}
		return &Error{Code: 4113, SQLState: "HY000", Message: "Wrong partitioning type, expected type: `SYSTEM_TIME`"}
	}
	return nil
}
