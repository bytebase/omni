package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/deparse"
)

// ============================================================================
// The function-argument reader: what goes inside a call's parentheses.
//
// funcArgCases is the recorded expectation, and
// TestTiDBFunctionArgumentOracle holds every row of it to a live TiDB v8.5.5.
// The accept rows alone would not pin much: the reject rows are what keep a
// quantifier from being read where the engine has none, and a prefix the
// engine reads from being dropped.
//
// The tree tests below carry the rest of the weight. A quantifier or a unit
// that parses and is then dropped changes what the query means, and an
// accept/reject gate cannot see that.
// ============================================================================

var funcArgCases = []struct {
	sql    string
	accept bool
}{
	// --- the quantifier prefix, per aggregate ---
	//
	// BuggyDefaultFalseDistinctOpt: DISTINCT, DISTINCTROW, ALL and the
	// DISTINCT ALL pair (pingcap/tidb pkg/parser/parser.y).
	{"SELECT SUM(c1) FROM t", true},
	{"SELECT SUM(DISTINCT c1) FROM t", true},
	{"SELECT SUM(DISTINCTROW c1) FROM t", true},
	{"SELECT SUM(ALL c1) FROM t", true},
	{"SELECT SUM(DISTINCT ALL c1) FROM t", true},
	{"SELECT SUM(DISTINCTROW ALL c1) FROM t", true},
	{"SELECT AVG(DISTINCTROW ALL c1) FROM t", true},
	{"SELECT MAX(DISTINCT ALL c1) FROM t", true},
	{"SELECT MIN(DISTINCTROW c1) FROM t", true},
	// The statistical aggregates take the same prefix set, STD and STDDEV
	// being spellings of STDDEV_POP and VARIANCE of VAR_POP.
	{"SELECT STD(ALL c1) FROM t", true},
	{"SELECT STD(DISTINCT c1) FROM t", true},
	{"SELECT STD(DISTINCTROW ALL c1) FROM t", true},
	{"SELECT STDDEV(DISTINCT c1) FROM t", true},
	{"SELECT STDDEV_POP(ALL c1) FROM t", true},
	{"SELECT STDDEV_SAMP(DISTINCT ALL c1) FROM t", true},
	{"SELECT VARIANCE(ALL c1) FROM t", true},
	{"SELECT VAR_POP(DISTINCTROW c1) FROM t", true},
	{"SELECT VAR_SAMP(ALL c1) FROM t", true},
	// GROUP_CONCAT takes the prefix set and an argument list, with its own
	// trailing clauses and a window after it.
	{"SELECT GROUP_CONCAT(DISTINCTROW c2, c1) FROM t GROUP BY c1", true},
	{"SELECT GROUP_CONCAT(ALL c2, c1) FROM t GROUP BY c1", true},
	{"SELECT GROUP_CONCAT(DISTINCT ALL c2 ORDER BY c1 SEPARATOR '-') FROM t", true},
	{"SELECT GROUP_CONCAT(c1) OVER () FROM t", true},
	// The bitwise aggregates take a bare ALL and no DISTINCT.
	{"SELECT BIT_AND(ALL c1) FROM t", true},
	{"SELECT BIT_OR(ALL c1) FROM t", true},
	{"SELECT BIT_XOR(ALL c1) FROM t", true},
	{"SELECT BIT_AND(DISTINCT c1) FROM t", false},
	{"SELECT BIT_OR(DISTINCTROW c1) FROM t", false},
	{"SELECT BIT_XOR(DISTINCT ALL c1) FROM t", false},
	// So do the JSON aggregates, and JSON_OBJECTAGG takes it before either of
	// its two arguments.
	{"SELECT JSON_ARRAYAGG(ALL c1) FROM t", true},
	{"SELECT JSON_ARRAYAGG(DISTINCT c1) FROM t", false},
	{"SELECT JSON_OBJECTAGG(c1, c2) FROM t", true},
	{"SELECT JSON_OBJECTAGG(ALL c1, c2) FROM t", true},
	{"SELECT JSON_OBJECTAGG(c1, ALL c2) FROM t", true},
	{"SELECT JSON_OBJECTAGG(ALL c1, ALL c2) FROM t", true},
	{"SELECT JSON_OBJECTAGG(DISTINCT c1, c2) FROM t", false},
	{"SELECT JSON_OBJECTAGG(ALL ALL c1, c2) FROM t", false},
	// COUNT's four productions: a star with no prefix, one argument, ALL and
	// one argument, and a DISTINCT with a list. A de-duplicating COUNT takes
	// no window, where SUM does.
	{"SELECT COUNT(*) FROM t", true},
	{"SELECT COUNT(c1) FROM t", true},
	{"SELECT COUNT(ALL c1) FROM t", true},
	{"SELECT COUNT(DISTINCT c1, c2) FROM t", true},
	{"SELECT COUNT(DISTINCTROW c1, c2) FROM t", true},
	{"SELECT COUNT(ALL profit) OVER () FROM sales", true},
	{"SELECT SUM(DISTINCT c1) OVER () FROM t", true},
	{"SELECT COUNT(DISTINCT c1) OVER () FROM t", false},
	{"SELECT COUNT(DISTINCT *) FROM t", false},
	{"SELECT COUNT(ALL *) FROM t", false},
	{"SELECT COUNT(c1, c2) FROM t", false},
	// A prefix twice over, or in the wrong order, is a syntax error.
	{"SELECT SUM(ALL ALL c1) FROM t", false},
	{"SELECT SUM(ALL DISTINCT c1) FROM t", false},
	{"SELECT SUM(DISTINCT DISTINCT c1) FROM t", false},
	{"SELECT SUM(c1 ALL) FROM t", false},
	// A builtin that is not an aggregate takes no prefix at all, and neither
	// does an aggregate under a quoted or qualified name, which is an
	// ordinary call.
	{"SELECT ABS(DISTINCT c1) FROM t", false},
	{"SELECT ANY_VALUE(DISTINCT c1) FROM t", false},
	{"SELECT APPROX_COUNT_DISTINCT(DISTINCT c1) FROM t", false},
	{"SELECT `sum`(DISTINCT c1) FROM t", false},
	{"SELECT `count`(DISTINCT c1) FROM t", false},
	// The argument count is the grammar's, so a wrong one is a syntax error
	// and not the 1582 an ordinary builtin answers.
	{"SELECT AVG() FROM t", false},
	{"SELECT AVG(c1, c2) FROM t", false},
	{"SELECT BIT_AND() FROM t", false},
	{"SELECT STD(c1, c2) FROM t", false},
	{"SELECT GROUP_CONCAT() FROM t", false},
	{"SELECT JSON_ARRAYAGG(c1, c2) FROM t", false},
	{"SELECT JSON_OBJECTAGG(c1, c2, c1) FROM t", false},
	// Only COUNT's own production has a star.
	{"SELECT SUM(*) FROM t", false},
	{"SELECT SUM(t.*) FROM t", false},
	{"SELECT COUNT(t.*) FROM t", false},
	{"SELECT ABS(*) FROM t", false},
	{"SELECT GROUP_CONCAT(*) FROM t", false},

	// --- EXTRACT's unit, the full interval-unit set ---
	{"SELECT EXTRACT(MICROSECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(SECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(MINUTE FROM a) FROM t", true},
	{"SELECT EXTRACT(HOUR FROM a) FROM t", true},
	{"SELECT EXTRACT(DAY FROM a) FROM t", true},
	{"SELECT EXTRACT(WEEK FROM a) FROM t", true},
	{"SELECT EXTRACT(MONTH FROM a) FROM t", true},
	{"SELECT EXTRACT(QUARTER FROM a) FROM t", true},
	{"SELECT EXTRACT(YEAR FROM a) FROM t", true},
	{"SELECT EXTRACT(SECOND_MICROSECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(MINUTE_MICROSECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(MINUTE_SECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(HOUR_MICROSECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(HOUR_SECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(HOUR_MINUTE FROM a) FROM t", true},
	{"SELECT EXTRACT(DAY_MICROSECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(DAY_SECOND FROM a) FROM t", true},
	{"SELECT EXTRACT(DAY_MINUTE FROM a) FROM t", true},
	{"SELECT EXTRACT(DAY_HOUR FROM a) FROM t", true},
	{"SELECT EXTRACT(YEAR_MONTH FROM a) FROM t", true},
	{"SELECT EXTRACT(SQL_TSI_DAY FROM a) FROM t", true},
	// The unit is a keyword token, never an identifier or a string, and only
	// a unit at all. A compound unit is reserved, so it is no column either.
	{"SELECT EXTRACT(`DAY_HOUR` FROM a) FROM t", false},
	{"SELECT EXTRACT('DAY_HOUR' FROM a) FROM t", false},
	{"SELECT EXTRACT(c1 FROM a) FROM t", false},
	{"SELECT EXTRACT(NOTAUNIT FROM a) FROM t", false},
	{"SELECT EXTRACT(SQL_TSI_MICROSECOND FROM a) FROM t", false},
	{"SELECT EXTRACT(DAY_HOUR, a) FROM t", false},
	{"SELECT EXTRACT(DAY_HOUR FROM a) OVER () FROM t", false},
	{"SELECT `extract`(DAY_HOUR FROM a) FROM t", false},

	// --- INTERVAL as a function, beside the interval literal ---
	{"SELECT INTERVAL(1, 0, 1, 2)", true},
	{"SELECT (INTERVAL(1, 0, 1, 2)+5)*7+INTERVAL(1, 0, 1, 2)/2", true},
	{"SELECT INTERVAL(0, (1*5)/2)+INTERVAL(5, 4, 3)", true},
	{"SELECT INTERVAL((SELECT 1), 2)", true},
	{"SELECT INTERVAL(c1, 1, 2) FROM t", true},
	{"SELECT 1 + INTERVAL(1, 2)", true},
	{"SELECT INTERVAL 1 DAY + INTERVAL(1, 2)", true},
	{"INSERT INTO t (c1) VALUES (INTERVAL(1, 2))", true},
	// TiDB reads the name through FunctionNameConflict, so a '(' after
	// INTERVAL always opens an argument list and the arity is the function's
	// own business, past parsing.
	{"SELECT INTERVAL(1)", true},
	{"SELECT INTERVAL()", true},
	{"SELECT a + INTERVAL (1) DAY FROM t", true},
	// The literal's value may be the INTERVAL function, which is what this
	// row needs. The engine takes a whole expression there, where this reader
	// takes a primary one: TestTiDBIntervalValueStaysUnread records why.
	{"SELECT DATE_ADD('2008-01-02', INTERVAL INTERVAL(1, 0, 1) DAY)", true},
	// The function's arguments are ordinary expressions: no quantifier, no
	// bare subquery, and no window after the call.
	{"SELECT INTERVAL(ALL 1, 2)", false},
	{"SELECT INTERVAL(SELECT 1, 2)", false},
	{"SELECT INTERVAL(c1) OVER () FROM t", false},

	// --- CAST ... ARRAY, the multi-valued index form ---
	{"SELECT CAST(b AS SIGNED ARRAY) FROM t", true},
	{"SELECT CAST(b AS UNSIGNED ARRAY) FROM t", true},
	{"SELECT CAST(b AS SIGNED INTEGER ARRAY) FROM t", true},
	{"SELECT CAST(b AS CHAR ARRAY) FROM t", true},
	{"SELECT CAST(b AS CHAR(10) ARRAY) FROM t", true},
	{"SELECT CAST(b AS DECIMAL(10,2) ARRAY) FROM t", true},
	{"SELECT CAST(b AS DATE ARRAY) FROM t", true},
	{"SELECT CAST(b AS DATETIME(3) ARRAY) FROM t", true},
	{"SELECT CAST(b AS TIME ARRAY) FROM t", true},
	{"SELECT CAST(b AS BINARY(10) ARRAY) FROM t", true},
	{"SELECT CAST(b AS DOUBLE ARRAY) FROM t", true},
	{"SELECT CAST(b AS CHAR(10) CHARACTER SET utf8mb4 ARRAY) FROM t", true},
	{"CREATE INDEX idx ON t ((CAST(b AS SIGNED ARRAY)))", true},
	// ARRAY is the cast's own tail: once, on a type, and never on CONVERT.
	{"SELECT CAST(b AS SIGNED ARRAY ARRAY) FROM t", false},
	{"SELECT CAST(b AS ARRAY) FROM t", false},
	{"SELECT CONVERT(b, SIGNED ARRAY) FROM t", false},

	// --- POSITION(substr IN str) ---
	{"SELECT POSITION('bar' IN 'foobarbar')", true},
	{"SELECT POSITION(c1 IN c2) FROM t", true},
	{"SELECT POSITION('bar' IN 'foobarbar') + 1", true},
	{"SELECT POSITION('a' + 1 IN 'b')", true},
	{"SELECT POSITION('a' IN 'b' = 1)", true},
	// The keyword is the only separator, the first operand stops short of a
	// comparison, and the call takes no window. A quoted name is an ordinary
	// call, which takes the comma form and not this one.
	{"SELECT POSITION('bar', 'foobarbar')", false},
	{"SELECT POSITION(1 = 1 IN 'x')", false},
	{"SELECT POSITION(DISTINCT 'a' IN 'b')", false},
	{"SELECT POSITION('a' IN 'b') OVER ()", false},
	{"SELECT POSITION()", false},
	{"SELECT POSITION('a')", false},
	{"SELECT POSITION(* IN 'a')", false},
	{"SELECT `position`('bar', 'foo')", true},
	{"SELECT `position`('bar' IN 'foo')", false},
}

// TestTiDBFunctionArgumentForms holds the parser to the recorded verdicts.
func TestTiDBFunctionArgumentForms(t *testing.T) {
	for _, c := range funcArgCases {
		_, err := Parse(c.sql)
		if (err == nil) != c.accept {
			t.Errorf("Parse(%q) accepted=%v, want %v (err=%v)", c.sql, err == nil, c.accept, err)
		}
	}
}

// TestTiDBAggQuantifierTree asserts the quantifier reaches the tree as written
// and that RemovesDuplicates reads it the way the engine does. A dropped
// DISTINCT would leave the statement parsing and meaning something else.
func TestTiDBAggQuantifierTree(t *testing.T) {
	cases := []struct {
		sql      string
		want     nodes.AggQuantifier
		distinct bool
		args     int
	}{
		{"SELECT SUM(c1) FROM t", nodes.AggQuantNone, false, 1},
		{"SELECT SUM(DISTINCT c1) FROM t", nodes.AggQuantDistinct, true, 1},
		{"SELECT SUM(DISTINCTROW c1) FROM t", nodes.AggQuantDistinctRow, true, 1},
		{"SELECT SUM(ALL c1) FROM t", nodes.AggQuantAll, false, 1},
		{"SELECT SUM(DISTINCT ALL c1) FROM t", nodes.AggQuantDistinctAll, true, 1},
		{"SELECT SUM(DISTINCTROW ALL c1) FROM t", nodes.AggQuantDistinctRowAll, true, 1},
		{"SELECT COUNT(DISTINCTROW c1, c2) FROM t", nodes.AggQuantDistinctRow, true, 2},
		{"SELECT COUNT(ALL c1) FROM t", nodes.AggQuantAll, false, 1},
		{"SELECT GROUP_CONCAT(DISTINCTROW c2, c1) FROM t", nodes.AggQuantDistinctRow, true, 2},
		{"SELECT GROUP_CONCAT(ALL c2) FROM t", nodes.AggQuantAll, false, 1},
		{"SELECT BIT_AND(ALL c1) FROM t", nodes.AggQuantAll, false, 1},
		{"SELECT STD(DISTINCT ALL c1) FROM t", nodes.AggQuantDistinctAll, true, 1},
		{"SELECT JSON_ARRAYAGG(ALL c1) FROM t", nodes.AggQuantAll, false, 1},
		// ALL on either argument is the one quantifier JSON_OBJECTAGG takes,
		// and it is the default: the engine keeps no record of which argument
		// carried it, so the call records one ALL wherever it appeared.
		{"SELECT JSON_OBJECTAGG(ALL c1, c2) FROM t", nodes.AggQuantAll, false, 2},
		{"SELECT JSON_OBJECTAGG(c1, ALL c2) FROM t", nodes.AggQuantAll, false, 2},
		{"SELECT JSON_OBJECTAGG(ALL c1, ALL c2) FROM t", nodes.AggQuantAll, false, 2},
		{"SELECT JSON_OBJECTAGG(c1, c2) FROM t", nodes.AggQuantNone, false, 2},
	}
	for _, c := range cases {
		fc := firstFuncCall(t, c.sql)
		if fc.Quantifier != c.want {
			t.Errorf("%q: Quantifier = %v (%q), want %v", c.sql, uint8(fc.Quantifier), fc.Quantifier, uint8(c.want))
		}
		if fc.Quantifier.RemovesDuplicates() != c.distinct {
			t.Errorf("%q: RemovesDuplicates = %v, want %v", c.sql, fc.Quantifier.RemovesDuplicates(), c.distinct)
		}
		if len(fc.Args) != c.args {
			t.Errorf("%q: %d arguments, want %d", c.sql, len(fc.Args), c.args)
		}
	}
}

// TestTiDBAggQuantifierInvariant asserts every quantifier spelling knows
// whether it de-duplicates, so a new one cannot be read and then counted as
// the default.
func TestTiDBAggQuantifierInvariant(t *testing.T) {
	want := map[nodes.AggQuantifier]struct {
		text     string
		distinct bool
	}{
		nodes.AggQuantNone:           {"", false},
		nodes.AggQuantDistinct:       {"DISTINCT", true},
		nodes.AggQuantDistinctRow:    {"DISTINCTROW", true},
		nodes.AggQuantAll:            {"ALL", false},
		nodes.AggQuantDistinctAll:    {"DISTINCT ALL", true},
		nodes.AggQuantDistinctRowAll: {"DISTINCTROW ALL", true},
	}
	for q, w := range want {
		if q.String() != w.text {
			t.Errorf("quantifier %d: String() = %q, want %q", uint8(q), q.String(), w.text)
		}
		if q.RemovesDuplicates() != w.distinct {
			t.Errorf("quantifier %q: RemovesDuplicates() = %v, want %v", q, q.RemovesDuplicates(), w.distinct)
		}
	}
	// Every spelling the parser reads must be one of them: a quantifier with
	// no text would deparse to a call with its prefix silently gone.
	for _, c := range funcArgCases {
		if !c.accept {
			continue
		}
		list, err := Parse(c.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.sql, err)
		}
		nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
			fc, ok := n.(*nodes.FuncCallExpr)
			if !ok {
				return true
			}
			if _, known := want[fc.Quantifier]; !known {
				t.Errorf("%q: %s carries unknown quantifier %d", c.sql, fc.Name, uint8(fc.Quantifier))
			}
			return true
		})
	}
}

// TestTiDBFunctionArgumentTree asserts the keyword arguments reach the tree:
// EXTRACT's unit, CAST's ARRAY and POSITION's two operands.
func TestTiDBFunctionArgumentTree(t *testing.T) {
	for _, c := range []struct {
		sql  string
		unit string
	}{
		{"SELECT EXTRACT(DAY FROM a) FROM t", "DAY"},
		{"SELECT EXTRACT(DAY_HOUR FROM a) FROM t", "DAY_HOUR"},
		{"SELECT EXTRACT(second_microsecond FROM a) FROM t", "SECOND_MICROSECOND"},
		{"SELECT EXTRACT(YEAR_MONTH FROM a) FROM t", "YEAR_MONTH"},
		// A SQL_TSI_ synonym folds to its canonical unit, which is the form
		// the engine stores.
		{"SELECT EXTRACT(SQL_TSI_DAY FROM a) FROM t", "DAY"},
	} {
		ex, ok := firstExpr(t, c.sql).(*nodes.ExtractExpr)
		if !ok {
			t.Fatalf("%q: not an ExtractExpr", c.sql)
		}
		if ex.Unit != c.unit {
			t.Errorf("%q: Unit = %q, want %q", c.sql, ex.Unit, c.unit)
		}
		if ex.Expr == nil {
			t.Errorf("%q: no operand", c.sql)
		}
	}

	for _, c := range []struct {
		sql      string
		typeName string
		array    bool
	}{
		{"SELECT CAST(b AS SIGNED ARRAY) FROM t", "SIGNED", true},
		{"SELECT CAST(b AS SIGNED) FROM t", "SIGNED", false},
		{"SELECT CAST(b AS CHAR(10) ARRAY) FROM t", "CHAR", true},
		{"SELECT CAST(b AS CHAR(10)) FROM t", "CHAR", false},
	} {
		ce, ok := firstExpr(t, c.sql).(*nodes.CastExpr)
		if !ok {
			t.Fatalf("%q: not a CastExpr", c.sql)
		}
		if ce.Array != c.array {
			t.Errorf("%q: Array = %v, want %v", c.sql, ce.Array, c.array)
		}
		if ce.TypeName == nil || !strings.EqualFold(ce.TypeName.Name, c.typeName) {
			t.Errorf("%q: type = %v, want %s", c.sql, ce.TypeName, c.typeName)
		}
	}

	pe, ok := firstExpr(t, "SELECT POSITION('bar' IN 'foobar')").(*nodes.PositionExpr)
	if !ok {
		t.Fatal("POSITION: not a PositionExpr")
	}
	substr, _ := pe.Substr.(*nodes.StringLit)
	str, _ := pe.Str.(*nodes.StringLit)
	if substr == nil || substr.Value != "bar" || str == nil || str.Value != "foobar" {
		t.Errorf("POSITION operands = %#v, %#v, want 'bar' in 'foobar'", pe.Substr, pe.Str)
	}
}

// TestTiDBFunctionArgumentDeparseRoundTrip asserts the stored text of every
// accepted form parses back to the same text, so nothing inside the
// parentheses is lost or gained in one pass.
func TestTiDBFunctionArgumentDeparseRoundTrip(t *testing.T) {
	for _, c := range funcArgCases {
		if !c.accept || !strings.HasPrefix(strings.ToUpper(c.sql), "SELECT ") {
			continue
		}
		// A backtick-quoted builtin is an ordinary call, and the tree keeps no
		// record of the quoting, so its text comes back out unquoted and reads
		// as the builtin's own syntax. That predates POSITION joining the
		// builtins read this way — `trim`(b, 'x') deparses the same — and it
		// is not this reader's to fix.
		if strings.Contains(c.sql, "`") {
			continue
		}
		list, err := Parse(c.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.sql, err)
		}
		sel, ok := list.Items[0].(*nodes.SelectStmt)
		if !ok || len(sel.TargetList) == 0 {
			continue
		}
		first := deparse.Deparse(targetExpr(sel.TargetList[0]))
		again, err := Parse("SELECT " + first)
		if err != nil {
			t.Errorf("%q: deparsed to %q, which does not parse: %v", c.sql, first, err)
			continue
		}
		sel2, ok := again.Items[0].(*nodes.SelectStmt)
		if !ok || len(sel2.TargetList) == 0 {
			t.Errorf("%q: deparsed to %q, which is not a select item", c.sql, first)
			continue
		}
		second := deparse.Deparse(targetExpr(sel2.TargetList[0]))
		if first != second {
			t.Errorf("%q: deparse is not stable: %q then %q", c.sql, first, second)
		}
	}
}

// firstFuncCall returns the first FuncCallExpr of sql.
func firstFuncCall(t *testing.T, sql string) *nodes.FuncCallExpr {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	var found *nodes.FuncCallExpr
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		if fc, ok := n.(*nodes.FuncCallExpr); ok && found == nil {
			found = fc
		}
		return true
	})
	if found == nil {
		t.Fatalf("Parse(%q): no function call in the tree", sql)
	}
	return found
}

// firstExpr returns the first select item's expression.
func firstExpr(t *testing.T, sql string) nodes.ExprNode {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	sel, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok || len(sel.TargetList) == 0 {
		t.Fatalf("Parse(%q): no select item", sql)
	}
	return targetExpr(sel.TargetList[0])
}

// targetExpr unwraps a select item, which carries its expression directly
// unless it was given an alias.
func targetExpr(item nodes.ExprNode) nodes.ExprNode {
	if rt, ok := item.(*nodes.ResTarget); ok {
		return rt.Val
	}
	return item
}

// TestTiDBIntervalValueStaysUnread records the narrowing that remains in this
// reader. The engine takes a whole expression for an interval literal's value
// and stores it unparenthesized. omni's deparse parenthesizes a NOT, a CASE, an
// operator or an EXISTS (and drops a BINARY), as mysql/deparse does, and
// `interval (not(`n`)) day` is the INTERVAL function to this grammar — TiDB
// v8.5.5 answers 1582 for it — so reading one of these would deparse to a
// statement the engine rejects or reads differently. A signed value is read now
// (`INTERVAL -1 DAY` deparses as `-1`, no parentheses — see
// TestIntervalSignedValue); the forms below are not.
//
// Accepting any of them means making its deparse round-trip first;
// TestTiDBFunctionArgumentDeparseRoundTrip is what catches the omission.
func TestTiDBIntervalValueStaysUnread(t *testing.T) {
	for _, sql := range []string{
		"SELECT a + INTERVAL NOT c1 DAY FROM t",
		"SELECT a + INTERVAL CASE WHEN 1 THEN 1 END DAY FROM t",
		"SELECT a + INTERVAL 1 + 2 DAY FROM t",
		// A sign is read only on a bare numeric literal; on anything else it is
		// narrowed, since deparse drops a unary plus and `interval (1 + 2) day`
		// is the INTERVAL function to TiDB.
		"SELECT a + INTERVAL +(1 + 2) DAY FROM t",
		"SELECT a + INTERVAL -(1 + 2) DAY FROM t",
		"SELECT a + INTERVAL EXISTS (SELECT 1) DAY FROM t",
		// This one deparses to `interval 1 day`, the operator gone, which the
		// engine does parse and means something else.
		"SELECT a + INTERVAL BINARY 1 DAY FROM t",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted; its deparse has to round-trip first", sql)
		}
	}
}
