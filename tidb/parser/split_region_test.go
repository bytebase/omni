package parser

import (
	"strconv"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's SPLIT and DISTRIBUTE statements, which no other MySQL-family engine
// has.
//
//	SPLIT [REGION FOR] [PARTITION] TABLE tbl [PARTITION (p, ...)]
//	    [INDEX idx] {BY (v, ...), ... | BETWEEN (v, ...) AND (v, ...) REGIONS n}
//	DISTRIBUTE TABLE tbl [PARTITION (p, ...)]
//	    RULE [=] 'rule' ENGINE [=] 'engine' [TIMEOUT [=] 'timeout']
//
// https://docs.pingcap.com/tidb/stable/sql-statement-split-region/
//
// The corpus in testdata carries the spellings pingcap's own parser test
// writes, all of them accepted. What it cannot carry is the shape of the
// rejections, because pingcap never writes them, and the line between a
// syntax error and a later failure is not where it looks: TiDB's parser takes
// `BY ()`, `BETWEEN () AND () REGIONS 10`, `REGIONS 0`, and rows of unequal
// width, and only the executor rejects them (1136, 1105, 1265). It is the
// sign on the REGIONS count, the missing AND, and the two options at once
// that are 1064. Every verdict below is the engine's:
// TestSplitDistributeMatchTiDB asks v8.5.5 about each one.

var splitDistributeAccepted = []string{
	// The enumerated form, on a table and on an index.
	`SPLIT TABLE t1 BY ('a'),('b'),('c')`,
	`SPLIT TABLE t1 BY (1)`,
	`SPLIT TABLE t1 INDEX idx1 BY ('abc',123), ('xyz'), ('yz', 1000)`,
	// The ranged form, on a table and on an index.
	`SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS 10`,
	`SPLIT TABLE t1 INDEX idx1 BETWEEN ('a',1) AND ('z',2) REGIONS 10`,
	// All four SplitSyntaxOption spellings, in the one order the grammar
	// has: nothing, REGION FOR, PARTITION, REGION FOR PARTITION.
	`SPLIT REGION FOR TABLE t1 BY ('a')`,
	`SPLIT PARTITION TABLE t1 BY ('a')`,
	`SPLIT REGION FOR PARTITION TABLE t1 BY ('a')`,
	`SPLIT REGION FOR PARTITION TABLE t1 PARTITION (p0,p1) INDEX idx1 BETWEEN ('a') AND ('z') REGIONS 10`,
	// An empty value row is a syntax-valid row: the executor rejects it
	// (1136 for BY, 1105 for BETWEEN), the parser does not.
	`SPLIT TABLE t1 BY ()`,
	`SPLIT TABLE t1 BY ('a'), ()`,
	`SPLIT TABLE t1 BETWEEN () AND () REGIONS 10`,
	`SPLIT TABLE t1 INDEX idx1 BETWEEN () AND ('z') REGIONS 10`,
	// So are rows of unequal width and a zero region count: both fail when
	// the statement runs (1136, 1265), not when it parses.
	`SPLIT TABLE t1 BY ('a'),(1,2,3)`,
	`SPLIT TABLE t1 BETWEEN ('a','b') AND ('z') REGIONS 10`,
	`SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS 0`,
	// A value is an expression, not just a literal. Two values the grammar
	// also allows are left out: the DEFAULT keyword, because PREPARE of it
	// panics v8.5.5's executor and drops the connection
	// (TestSplitByDefaultValueParses covers it on this side), and the `?`
	// placeholder, which this parser reads nowhere -- `SELECT ?` is a
	// syntax error here too, so it is not this family's gap to close.
	`SPLIT TABLE t1 BY (1+1)`,
	`SPLIT TABLE t1 BY (now())`,
	`SPLIT TABLE t1 BY (-1)`,
	`SPLIT TABLE t1 BY (null)`,
	`SPLIT TABLE t1 BY (b)`,
	// A qualified and a quoted name, and an index name that is one of the
	// words this family registers as a keyword. All five are unreserved
	// upstream, so each is still an ordinary identifier.
	"SPLIT TABLE probe.t1 BY ('a')",
	"SPLIT TABLE `t1` INDEX `idx1` BY ('a')",
	`SPLIT TABLE t1 INDEX split BY ('a')`,
	`SPLIT TABLE t1 INDEX region BY ('a')`,
	`SPLIT TABLE t1 INDEX rule BY ('a')`,
	`SPLIT TABLE t1 INDEX timeout BY ('a')`,
	`SPLIT TABLE t1 INDEX distribute BY ('a')`,
	`SPLIT TABLE split BY ('a')`,
	`SPLIT TABLE t1 PARTITION (split, rule) BY ('a')`,
	"SPLIT TABLE t1 INDEX `primary` BY (1)",

	// DISTRIBUTE: RULE then ENGINE, both mandatory, TIMEOUT last and
	// optional, each `=` optional, each value a string literal.
	`DISTRIBUTE TABLE t1 RULE = 'leader-scatter' ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1 RULE = "leader-scatter" ENGINE = "tikv"`,
	`DISTRIBUTE TABLE t1 RULE 'leader-scatter' ENGINE 'tikv'`,
	`DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = '30m'`,
	`DISTRIBUTE TABLE t1 RULE 'x' ENGINE 'tikv' TIMEOUT '30m'`,
	`DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = ''`,
	`DISTRIBUTE TABLE t1 PARTITION(p0,p1) RULE = 'learner-scatter' ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE probe.t1 RULE = 'x' ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE split RULE = 'x' ENGINE = 'tikv'`,
}

var splitDistributeRejected = []string{
	// The split option is mandatory and comes once.
	`SPLIT TABLE t1`,
	`SPLIT TABLE`,
	`SPLIT TABLE t1 INDEX idx1`,
	`SPLIT TABLE t1 REGIONS 10`,
	`SPLIT TABLE t1 INDEX idx1 REGIONS 10`,
	`SPLIT TABLE t1 BY ('a') BETWEEN ('a') AND ('z') REGIONS 10`,
	`SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS 10 BY ('a')`,
	// BY takes one or more rows: none and a trailing comma are both 1064,
	// where an empty row is not.
	`SPLIT TABLE t1 BY`,
	`SPLIT TABLE t1 INDEX idx1 BY`,
	`SPLIT TABLE t1 BY ('a'),`,
	// The ROW(...) constructor an INSERT row takes is not a RowValue.
	`SPLIT TABLE t1 BY row('a')`,
	// BETWEEN needs its AND and its REGIONS count.
	`SPLIT TABLE t1 BETWEEN ('a') AND ('z')`,
	`SPLIT TABLE t1 BETWEEN ('a') REGIONS 10`,
	`SPLIT TABLE t1 BETWEEN ('a') AND ('z') REGIONS`,
	// REGIONS reads Int64Num: a bare unsigned integer literal that fits
	// int64. A sign, a decimal, an exponent, a hex literal, a string, a
	// placeholder, and an out-of-range value are all syntax errors, so
	// REGIONS -1 is a 1064 where REGIONS 0 is not.
	`SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS -1`,
	`SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS +10`,
	`SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 10.5`,
	`SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 1e3`,
	`SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 0x10`,
	`SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 'a'`,
	`SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS ?`,
	`SPLIT TABLE t1 BETWEEN (1) AND (2) REGIONS 99999999999999999999`,
	// One table and one index, never a list.
	`SPLIT TABLE t1, t2 BY ('a')`,
	`SPLIT TABLE t1 INDEX idx1, idx2 BY ('a')`,
	// INDEX names an index, and the name is an Identifier: PRIMARY is
	// reserved, so it needs quoting.
	`SPLIT TABLE t1 INDEX BY (1)`,
	`SPLIT TABLE t1 INDEX primary BY (1)`,
	// The PARTITION (p, ...) list is non-empty, is never the bare keyword,
	// and comes before INDEX.
	`SPLIT TABLE t1 PARTITION INDEX idx1 BY ('a')`,
	`SPLIT TABLE t1 PARTITION () INDEX idx1 BY ('a')`,
	`SPLIT TABLE t1 INDEX idx1 PARTITION (p0) BY ('a')`,
	// SplitSyntaxOption has four spellings and no others: REGION is always
	// followed by FOR, TABLE always follows, and PARTITION comes after
	// REGION FOR, never before it.
	`SPLIT REGION FOR t1 BY ('a')`,
	`SPLIT REGION TABLE t1 BY ('a')`,
	`SPLIT REGION PARTITION TABLE t1 BY ('a')`,
	`SPLIT PARTITION REGION FOR TABLE t1 BY ('a')`,
	`SPLIT PARTITION t1 BY ('a')`,
	`SPLIT REGION FOR PARTITION t1 BY ('a')`,
	`SPLIT t1 BY ('a')`,
	// MAXVALUE is not an expression, so it is not a value row's value.
	`SPLIT TABLE t1 BY (maxvalue)`,

	// DISTRIBUTE: RULE and ENGINE are both mandatory, fixed in that order,
	// each value a string literal, with no comma between the options.
	`DISTRIBUTE TABLE t1`,
	`DISTRIBUTE TABLE t1 PARTITION(p0)`,
	`DISTRIBUTE TABLE t1 RULE = 'leader-scatter'`,
	`DISTRIBUTE TABLE t1 ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1 TIMEOUT = '30m'`,
	`DISTRIBUTE TABLE t1 ENGINE = 'tikv' RULE = 'leader-scatter'`,
	`DISTRIBUTE TABLE t1 TIMEOUT = '30m' RULE = 'x' ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = '30m' RULE = 'y'`,
	`DISTRIBUTE TABLE t1 RULE = leader ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1 RULE = 1 ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = tikv`,
	`DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = 30`,
	`DISTRIBUTE TABLE t1 RULE='x',ENGINE='tikv'`,
	`DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' EXTRA = 'y'`,
	// No index clause, one table, and the same non-empty partition list.
	`DISTRIBUTE TABLE t1 INDEX idx1 RULE = 'x' ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1, t2 RULE = 'x' ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1 PARTITION() RULE = 'x' ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1 PARTITION RULE = 'x' ENGINE = 'tikv'`,
	`DISTRIBUTE t1 RULE = 'x' ENGINE = 'tikv'`,
	`DISTRIBUTE TABLE t1 RULE = 'x' ENGINE = 'tikv' TIMEOUT = '1m' TIMEOUT = '2m'`,
}

// splitKeywordIdentifiers are the words this family registers as keywords,
// used where TiDB reads an identifier. All five are unreserved upstream --
// SPLIT, DISTRIBUTE, and REGION are TiDBKeyword (parser.y:7390, 7366, 7395),
// RULE and TIMEOUT UnReservedKeyword (parser.y:7062, 7079) -- so registering
// them must not cost any of these. The implicit table alias in DELETE is
// here because it read tokIDENT alone and so lost every keyword this package
// registers.
var splitKeywordIdentifiers = []string{
	`SELECT 1 AS split`,
	`SELECT 1 AS distribute`,
	`SELECT 1 AS region`,
	`SELECT 1 AS rule`,
	`SELECT 1 AS timeout`,
	`SELECT split, distribute, region, rule, timeout FROM t1`,
	`SELECT a FROM t1 AS rule`,
	`SELECT rule.a FROM t1 rule`,
	`DELETE FROM t1 split WHERE split.a = 'x'`,
	`DELETE FROM t1 distribute WHERE distribute.a = 'x'`,
	`DELETE FROM t1 region WHERE region.a = 'x'`,
	`DELETE FROM t1 rule WHERE rule.a = 'x'`,
	`DELETE FROM t1 timeout WHERE timeout.a = 'x'`,
	`CREATE TABLE split (split int, distribute int, region int, rule int, timeout int)`,
	`CREATE INDEX split ON t1 (a)`,
	`UPDATE t1 SET rule = 1 WHERE timeout = 2`,
	`INSERT INTO t1 (rule, timeout) VALUES (1,2)`,
	`SAVEPOINT split`,
	`ANALYZE TABLE split`,
	`TRUNCATE TABLE region`,
}

func TestSplitDistributeAccepted(t *testing.T) {
	for _, sql := range splitDistributeAccepted {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

func TestSplitDistributeRejected(t *testing.T) {
	for _, sql := range splitDistributeRejected {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but TiDB v8.5.5 rejects it", sql)
			}
		})
	}
}

func TestSplitKeywordsStayIdentifiers(t *testing.T) {
	for _, sql := range splitKeywordIdentifiers {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

// TestSplitByDefaultValueParses covers the one value the grammar allows that
// the oracle table cannot: ExprOrDefault's DEFAULT keyword
// (parser.y:7712-7717). PREPARE of `SPLIT TABLE t BY (default)` panics
// v8.5.5's executor with a nil dereference and drops the connection, so it is
// not asked of the shared container.
func TestSplitByDefaultValueParses(t *testing.T) {
	for _, sql := range []string{
		`SPLIT TABLE t1 BY (default)`,
		`SPLIT TABLE t1 BETWEEN (default) AND (default) REGIONS 10`,
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) error: %v", sql, err)
		}
	}
}

// splitStmt parses one SPLIT statement and returns its node. The corpus
// records only whether a statement parses, so the shape of the tree is
// asserted here: value rows flattened into one list, or an index name landing
// as the table name, would pass the corpus gate unseen.
func splitStmt(t *testing.T, sql string) *nodes.SplitRegionStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if list == nil || len(list.Items) != 1 {
		t.Fatalf("Parse(%q) gave %d statements, want 1", sql, list.Len())
	}
	stmt, ok := list.Items[0].(*nodes.SplitRegionStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave %T, want *nodes.SplitRegionStmt", sql, list.Items[0])
	}
	return stmt
}

func distributeStmt(t *testing.T, sql string) *nodes.DistributeTableStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if list == nil || len(list.Items) != 1 {
		t.Fatalf("Parse(%q) gave %d statements, want 1", sql, list.Len())
	}
	stmt, ok := list.Items[0].(*nodes.DistributeTableStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave %T, want *nodes.DistributeTableStmt", sql, list.Items[0])
	}
	return stmt
}

// rowStrings renders each value row as the literal texts it holds, so a test
// can assert that the rows kept their own boundaries.
func rowStrings(t *testing.T, rows [][]nodes.ExprNode) []string {
	t.Helper()
	out := make([]string, len(rows))
	for i, row := range rows {
		s := "("
		for j, e := range row {
			if j > 0 {
				s += ","
			}
			switch v := e.(type) {
			case *nodes.StringLit:
				s += "'" + v.Value + "'"
			case *nodes.IntLit:
				s += strconv.FormatInt(v.Value, 10)
			default:
				t.Fatalf("row %d value %d is %T, want a literal", i, j, e)
			}
		}
		out[i] = s + ")"
	}
	return out
}

func eqStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestSplitRegionTree asserts the parse of the enumerated form: the node
// type, the table, the index, the partition list, and the full nested value
// rows. `BY ('abc',123), ('xyz'), ('yz', 1000)` is three rows of 2, 1, and 2
// values, not one list of five.
func TestSplitRegionTree(t *testing.T) {
	stmt := splitStmt(t, `SPLIT REGION FOR PARTITION TABLE db1.t1 PARTITION (p0, p1) INDEX idx1 BY ('abc',123), ('xyz'), ('yz', 1000)`)

	if !stmt.RegionFor {
		t.Error("RegionFor = false, want true: the statement spelled REGION FOR")
	}
	if !stmt.Partition {
		t.Error("Partition = false, want true: the statement spelled PARTITION before TABLE")
	}
	if stmt.Table == nil || stmt.Table.Schema != "db1" || stmt.Table.Name != "t1" {
		t.Errorf("Table = %+v, want db1.t1", stmt.Table)
	}
	if got := stmt.Table.Partitions; !eqStrings(got, []string{"p0", "p1"}) {
		t.Errorf("Table.Partitions = %v, want [p0 p1]", got)
	}
	if stmt.Index != "idx1" {
		t.Errorf("Index = %q, want idx1: the index name is not the table's", stmt.Index)
	}
	if stmt.Opt == nil {
		t.Fatal("Opt = nil, want the BY option")
	}
	if stmt.Opt.Ranged {
		t.Error("Opt.Ranged = true, want false: the statement used BY")
	}
	if got := rowStrings(t, stmt.Opt.ValueLists); !eqStrings(got, []string{"('abc',123)", "('xyz')", "('yz',1000)"}) {
		t.Errorf("Opt.ValueLists = %v, want three rows of 2, 1, and 2 values", got)
	}
	if stmt.Opt.Lower != nil || stmt.Opt.Upper != nil || stmt.Opt.Regions != 0 {
		t.Errorf("the BY form carries ranged fields: lower=%v upper=%v regions=%d",
			stmt.Opt.Lower, stmt.Opt.Upper, stmt.Opt.Regions)
	}
	// The index name's own range covers idx1 and nothing else.
	if got := sqlAt(`SPLIT REGION FOR PARTITION TABLE db1.t1 PARTITION (p0, p1) INDEX idx1 BY ('abc',123), ('xyz'), ('yz', 1000)`, stmt.IndexLoc); got != "idx1" {
		t.Errorf("IndexLoc covers %q, want idx1", got)
	}
}

// TestSplitRegionRangedTree asserts the parse of the ranged form: the two
// bound rows stay apart and keep their own values, and the REGIONS count is
// read as a number.
func TestSplitRegionRangedTree(t *testing.T) {
	sql := `SPLIT TABLE t1 INDEX idx2 BETWEEN ('a',1) AND ('z',2) REGIONS 10`
	stmt := splitStmt(t, sql)

	if stmt.RegionFor || stmt.Partition {
		t.Errorf("RegionFor=%v Partition=%v, want both false", stmt.RegionFor, stmt.Partition)
	}
	if stmt.Table == nil || stmt.Table.Schema != "" || stmt.Table.Name != "t1" || len(stmt.Table.Partitions) != 0 {
		t.Errorf("Table = %+v, want t1 with no partition list", stmt.Table)
	}
	if stmt.Index != "idx2" {
		t.Errorf("Index = %q, want idx2", stmt.Index)
	}
	if !stmt.Opt.Ranged {
		t.Fatal("Opt.Ranged = false, want true: the statement used BETWEEN")
	}
	if got := rowStrings(t, [][]nodes.ExprNode{stmt.Opt.Lower}); !eqStrings(got, []string{"('a',1)"}) {
		t.Errorf("Opt.Lower = %v, want ('a',1)", got)
	}
	if got := rowStrings(t, [][]nodes.ExprNode{stmt.Opt.Upper}); !eqStrings(got, []string{"('z',2)"}) {
		t.Errorf("Opt.Upper = %v, want ('z',2)", got)
	}
	if stmt.Opt.Regions != 10 {
		t.Errorf("Opt.Regions = %d, want 10", stmt.Opt.Regions)
	}
	if stmt.Opt.ValueLists != nil {
		t.Errorf("Opt.ValueLists = %v, want nil for the ranged form", stmt.Opt.ValueLists)
	}
	// The statement's own range ends at the count, not past it.
	if got := sqlAt(sql, stmt.Loc); got != sql {
		t.Errorf("Loc covers %q, want the whole statement", got)
	}
}

// TestSplitRegionEmptyRowsTree pins what the corpus row
// `split table t1 between () and () regions 10` means: two bound rows that
// exist and are empty, which is not the same tree as a statement with no
// bounds at all. TiDB's parser takes it and its executor rejects it (1105).
func TestSplitRegionEmptyRowsTree(t *testing.T) {
	ranged := splitStmt(t, `SPLIT TABLE t1 BETWEEN () AND () REGIONS 10`)
	if !ranged.Opt.Ranged {
		t.Error("Opt.Ranged = false, want true")
	}
	if ranged.Opt.Lower == nil || len(ranged.Opt.Lower) != 0 {
		t.Errorf("Opt.Lower = %v, want an empty, non-nil row", ranged.Opt.Lower)
	}
	if ranged.Opt.Upper == nil || len(ranged.Opt.Upper) != 0 {
		t.Errorf("Opt.Upper = %v, want an empty, non-nil row", ranged.Opt.Upper)
	}

	// `BY ()` is one empty row, and `BY ('a'), ()` is two rows of which the
	// second is empty: the rows are not collapsed away.
	by := splitStmt(t, `SPLIT TABLE t1 BY ()`)
	if by.Opt.Ranged {
		t.Error("Opt.Ranged = true, want false")
	}
	if len(by.Opt.ValueLists) != 1 || by.Opt.ValueLists[0] == nil || len(by.Opt.ValueLists[0]) != 0 {
		t.Errorf("Opt.ValueLists = %v, want one empty, non-nil row", by.Opt.ValueLists)
	}
	mixed := splitStmt(t, `SPLIT TABLE t1 BY ('a'), ()`)
	if len(mixed.Opt.ValueLists) != 2 {
		t.Fatalf("Opt.ValueLists = %v, want two rows", mixed.Opt.ValueLists)
	}
	if got := rowStrings(t, mixed.Opt.ValueLists[:1]); !eqStrings(got, []string{"('a')"}) {
		t.Errorf("first row = %v, want ('a')", got)
	}
	if len(mixed.Opt.ValueLists[1]) != 0 {
		t.Errorf("second row = %v, want empty", mixed.Opt.ValueLists[1])
	}
}

// TestSplitRegionNoIndexTree pins that the table form leaves Index empty, so
// a reader cannot mistake the table's name for an index's.
func TestSplitRegionNoIndexTree(t *testing.T) {
	stmt := splitStmt(t, `SPLIT TABLE t1 BY ('a')`)
	if stmt.Index != "" {
		t.Errorf("Index = %q, want empty: the statement names no index", stmt.Index)
	}
	if stmt.IndexLoc != nodes.NoLoc() {
		t.Errorf("IndexLoc = %+v, want NoLoc()", stmt.IndexLoc)
	}
	if stmt.Table == nil || stmt.Table.Name != "t1" {
		t.Errorf("Table = %+v, want t1", stmt.Table)
	}
}

// TestDistributeTableTree asserts the parse of DISTRIBUTE TABLE: the node
// type, the table with its partition list, and each of the three option
// values in its own field.
func TestDistributeTableTree(t *testing.T) {
	sql := `DISTRIBUTE TABLE db1.t1 PARTITION (p0, p1) RULE = 'peer-scatter' ENGINE = 'tiflash' TIMEOUT = '30m'`
	stmt := distributeStmt(t, sql)

	if stmt.Table == nil || stmt.Table.Schema != "db1" || stmt.Table.Name != "t1" {
		t.Errorf("Table = %+v, want db1.t1", stmt.Table)
	}
	if got := stmt.Table.Partitions; !eqStrings(got, []string{"p0", "p1"}) {
		t.Errorf("Table.Partitions = %v, want [p0 p1]", got)
	}
	if stmt.Rule != "peer-scatter" {
		t.Errorf("Rule = %q, want peer-scatter", stmt.Rule)
	}
	if stmt.Engine != "tiflash" {
		t.Errorf("Engine = %q, want tiflash", stmt.Engine)
	}
	if stmt.Timeout != "30m" {
		t.Errorf("Timeout = %q, want 30m", stmt.Timeout)
	}
	if got := sqlAt(sql, stmt.Loc); got != sql {
		t.Errorf("Loc covers %q, want the whole statement", got)
	}

	// Without TIMEOUT the field stays empty, and the equals signs are
	// optional without changing the values.
	bare := distributeStmt(t, `DISTRIBUTE TABLE t1 RULE 'leader-scatter' ENGINE 'tikv'`)
	if bare.Rule != "leader-scatter" || bare.Engine != "tikv" || bare.Timeout != "" {
		t.Errorf("rule=%q engine=%q timeout=%q, want leader-scatter/tikv/empty",
			bare.Rule, bare.Engine, bare.Timeout)
	}
	if len(bare.Table.Partitions) != 0 {
		t.Errorf("Table.Partitions = %v, want none", bare.Table.Partitions)
	}
}

// sqlAt returns the text a Loc covers.
func sqlAt(sql string, loc nodes.Loc) string {
	if loc.Start < 0 || loc.End < 0 || loc.End > len(sql) || loc.Start > loc.End {
		return ""
	}
	return sql[loc.Start:loc.End]
}

// TestSplitDistributeMatchTiDB holds the three tables above to the engine
// itself (container test).
func TestSplitDistributeMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for want, cases := range map[bool][][]string{
		true:  {splitDistributeAccepted, splitKeywordIdentifiers},
		false: {splitDistributeRejected},
	} {
		for _, group := range cases {
			for _, sql := range group {
				t.Run(sql, func(t *testing.T) {
					accepts, err := tc.prepareSyntax(sql)
					if err != nil {
						t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
					}
					if accepts != want {
						t.Errorf("TiDB %ss %q, this file records %s",
							verdictName(accepts), sql, verdictName(want))
					}
				})
			}
		}
	}
}
