package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// What the statistics statements leave in the tree.
//
// The corpus records only whether a statement parses, which is why a family
// can close every one of its rows and still ship a defect: an option that
// parses into the wrong field, a list whose boundaries are lost, a choice
// that reads as its neighbour. `ANALYZE TABLE t WITH 4 CMSKETCH WIDTH` and
// `... CMSKETCH DEPTH` are the case in point -- both parse, and nothing in
// the corpus says they must not land identically.

// parseOneStmt parses sql and returns its single statement.
func parseOneStmt(t *testing.T, sql string) ast.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) gave %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}

func parseOneAnalyze(t *testing.T, sql string) *ast.AnalyzeTableStmt {
	t.Helper()
	stmt, ok := parseOneStmt(t, sql).(*ast.AnalyzeTableStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave a statement of another kind", sql)
	}
	return stmt
}

// tableNames renders a table list as "db.name" entries, so a list that lost
// or reordered an entry reads differently.
func tableNames(refs []*ast.TableRef) []string {
	out := make([]string, 0, len(refs))
	for _, r := range refs {
		if r.Schema != "" {
			out = append(out, r.Schema+"."+r.Name)
			continue
		}
		out = append(out, r.Name)
	}
	return out
}

func optionStrings(opts []*ast.AnalyzeOption) []string {
	out := make([]string, 0, len(opts))
	for _, o := range opts {
		out = append(out, fmt.Sprintf("%s %s", o.Value, o.Type))
	}
	return out
}

func join(ss []string) string { return "[" + strings.Join(ss, "|") + "]" }

func wantList(t *testing.T, field string, got, want []string) {
	t.Helper()
	if join(got) != join(want) {
		t.Errorf("%s = %s, want %s", field, join(got), join(want))
	}
}

// TestAnalyzeTreeTargets reads back the table, partition, index and column
// lists of every arm. A statement that named one table and one index must not
// come out with an empty index list, which is what distinguishes `INDEX` from
// `INDEX a` and is the whole of the bare-INDEX form.
func TestAnalyzeTreeTargets(t *testing.T) {
	for _, tc := range []struct {
		sql        string
		tables     []string
		partitions []string
		hasIndex   bool
		indexes    []string
		choice     ast.AnalyzeColumnChoice
		columns    []string
	}{
		{sql: "analyze table t", tables: []string{"t"}, choice: ast.AnalyzeColumnsUnset},
		{sql: "analyze table t1, t2, t3", tables: []string{"t1", "t2", "t3"}, choice: ast.AnalyzeColumnsUnset},
		{sql: "analyze table db1.t1", tables: []string{"db1.t1"}, choice: ast.AnalyzeColumnsUnset},
		// A bare INDEX is the arm with an empty list, not the absence of
		// the arm: HasIndex is what tells them apart.
		{sql: "analyze table t index", tables: []string{"t"}, hasIndex: true, indexes: nil},
		{sql: "analyze table t index a", tables: []string{"t"}, hasIndex: true, indexes: []string{"a"}},
		{sql: "analyze table t index a,b", tables: []string{"t"}, hasIndex: true, indexes: []string{"a", "b"}},
		// TiDB's leading comma is the empty base case of its list, so the
		// names after it are the whole list and the comma adds none.
		{sql: "analyze table t index ,a", tables: []string{"t"}, hasIndex: true, indexes: []string{"a"}},
		{sql: "analyze table t index ,a,b", tables: []string{"t"}, hasIndex: true, indexes: []string{"a", "b"}},
		{sql: "analyze table t index `a`", tables: []string{"t"}, hasIndex: true, indexes: []string{"a"}},
		{sql: "analyze table t partition p0", tables: []string{"t"}, partitions: []string{"p0"}},
		{sql: "analyze table t partition p0,p1", tables: []string{"t"}, partitions: []string{"p0", "p1"}},
		{sql: "analyze table db.t partition p0 index b", tables: []string{"db.t"}, partitions: []string{"p0"}, hasIndex: true, indexes: []string{"b"}},
		{sql: "analyze table t partition p0 index", tables: []string{"t"}, partitions: []string{"p0"}, hasIndex: true, indexes: nil},
		// The three column choices are three values, not one flag.
		{sql: "analyze table t all columns", tables: []string{"t"}, choice: ast.AnalyzeAllColumns},
		{sql: "analyze table t predicate columns", tables: []string{"t"}, choice: ast.AnalyzePredicateColumns},
		{sql: "analyze table t columns c1,c2", tables: []string{"t"}, choice: ast.AnalyzeColumnList, columns: []string{"c1", "c2"}},
		{sql: "analyze table t columns `c1`, `c2`", tables: []string{"t"}, choice: ast.AnalyzeColumnList, columns: []string{"c1", "c2"}},
		{sql: "analyze table t1,t2 all columns", tables: []string{"t1", "t2"}, choice: ast.AnalyzeAllColumns},
		{sql: "analyze table t1,t2 predicate columns", tables: []string{"t1", "t2"}, choice: ast.AnalyzePredicateColumns},
		{sql: "analyze table t partition p0 columns c1,c2", tables: []string{"t"}, partitions: []string{"p0"}, choice: ast.AnalyzeColumnList, columns: []string{"c1", "c2"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneAnalyze(t, tc.sql)
			wantList(t, "Tables", tableNames(stmt.Tables), tc.tables)
			wantList(t, "Partitions", stmt.Partitions, tc.partitions)
			if stmt.HasIndex != tc.hasIndex {
				t.Errorf("HasIndex = %v, want %v", stmt.HasIndex, tc.hasIndex)
			}
			wantList(t, "Indexes", stmt.Indexes, tc.indexes)
			if stmt.ColumnChoice != tc.choice {
				t.Errorf("ColumnChoice = %d, want %d", stmt.ColumnChoice, tc.choice)
			}
			wantList(t, "Columns", stmt.Columns, tc.columns)
		})
	}
}

// TestAnalyzeTreeOptions reads back the WITH list: every unit on its own,
// every unit together, and the value each carries. CMSKETCH WIDTH and
// CMSKETCH DEPTH are two units and must not land as one.
func TestAnalyzeTreeOptions(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want []string
	}{
		{"analyze table t with 4 buckets", []string{"4 BUCKETS"}},
		{"analyze table t with 4 topn", []string{"4 TOPN"}},
		{"analyze table t with 4 cmsketch width", []string{"4 CMSKETCH WIDTH"}},
		{"analyze table t with 4 cmsketch depth", []string{"4 CMSKETCH DEPTH"}},
		{"analyze table t with 4 samples", []string{"4 SAMPLES"}},
		{"analyze table t with 10 samplerate", []string{"10 SAMPLERATE"}},
		// Spelling is not the value: the unit is read case-insensitively
		// and the value is kept as written.
		{"analyze table t with 4 BUCKETS", []string{"4 BUCKETS"}},
		{"analyze table t with 4 CMSketch Width", []string{"4 CMSKETCH WIDTH"}},
		// SAMPLERATE's number keeps its spelling, so a consumer reads the
		// rate the statement wrote rather than a rounded one.
		{"analyze table t with 0.1 samplerate", []string{"0.1 SAMPLERATE"}},
		{"analyze table t with .5 samplerate", []string{".5 SAMPLERATE"}},
		{"analyze table t with 1e3 samplerate", []string{"1e3 SAMPLERATE"}},
		{"analyze table t with 1.0 samplerate", []string{"1.0 SAMPLERATE"}},
		{"analyze table t with 18446744073709551615 buckets", []string{"18446744073709551615 BUCKETS"}},
		// The whole list, in the order written, with every unit distinct.
		{"analyze table t with 4 buckets, 4 topn, 4 cmsketch width, 4 cmsketch depth, 4 samples",
			[]string{"4 BUCKETS", "4 TOPN", "4 CMSKETCH WIDTH", "4 CMSKETCH DEPTH", "4 SAMPLES"}},
		{"analyze table t with 0.1 samplerate, 4 buckets", []string{"0.1 SAMPLERATE", "4 BUCKETS"}},
		{"analyze table t with 4 buckets, 0.5 samplerate", []string{"4 BUCKETS", "0.5 SAMPLERATE"}},
		// A repeat is kept, since TiDB's parser takes it.
		{"analyze table t with 4 buckets, 8 buckets", []string{"4 BUCKETS", "8 BUCKETS"}},
		// Every arm carries the list, and none of them loses it.
		{"analyze table t index a with 4 buckets, 4 topn", []string{"4 BUCKETS", "4 TOPN"}},
		{"analyze table t partition p0 index a with 4 buckets", []string{"4 BUCKETS"}},
		{"analyze table t all columns with 1024 buckets", []string{"1024 BUCKETS"}},
		{"analyze table t columns c1 with 4 buckets, 4 topn", []string{"4 BUCKETS", "4 TOPN"}},
		{"analyze table t update histogram on b with 4 buckets, 4 topn", []string{"4 BUCKETS", "4 TOPN"}},
		// No WITH is an empty list, not a zero-valued option.
		{"analyze table t", nil},
		{"analyze table t drop histogram on b", nil},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneAnalyze(t, tc.sql)
			wantList(t, "Options", optionStrings(stmt.Options), tc.want)
		})
	}
}

// TestAnalyzeTreeModifiers reads back the modifiers and the histogram arms,
// which the WITH list and the column choice must not disturb.
func TestAnalyzeTreeModifiers(t *testing.T) {
	for _, tc := range []struct {
		sql         string
		binlog      bool
		incremental bool
		histOp      string
		histCols    []string
	}{
		{sql: "analyze table t"},
		{sql: "analyze local table t", binlog: true},
		{sql: "analyze no_write_to_binlog table t", binlog: true},
		{sql: "analyze incremental table t index a", incremental: true},
		{sql: "analyze local incremental table t index a", binlog: true, incremental: true},
		{sql: "analyze no_write_to_binlog incremental table t partition p0 index a", binlog: true, incremental: true},
		{sql: "analyze table t update histogram on b", histOp: "UPDATE", histCols: []string{"b"}},
		{sql: "analyze table t update histogram on c1, c2", histOp: "UPDATE", histCols: []string{"c1", "c2"}},
		{sql: "analyze table t drop histogram on b", histOp: "DROP", histCols: []string{"b"}},
		{sql: "analyze table t drop histogram on c1, c2", histOp: "DROP", histCols: []string{"c1", "c2"}},
		{sql: "analyze no_write_to_binlog table t update histogram on b", binlog: true, histOp: "UPDATE", histCols: []string{"b"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneAnalyze(t, tc.sql)
			if stmt.NoWriteToBinlog != tc.binlog {
				t.Errorf("NoWriteToBinlog = %v, want %v", stmt.NoWriteToBinlog, tc.binlog)
			}
			if stmt.Incremental != tc.incremental {
				t.Errorf("Incremental = %v, want %v", stmt.Incremental, tc.incremental)
			}
			if stmt.HistogramOp != tc.histOp {
				t.Errorf("HistogramOp = %q, want %q", stmt.HistogramOp, tc.histOp)
			}
			wantList(t, "HistogramColumns", stmt.HistogramColumns, tc.histCols)
		})
	}
}

// TestDropStatsTree reads back DROP STATS: which arm a statement took is a
// field, so the partition arm and the GLOBAL arm cannot read as the plain
// list.
func TestDropStatsTree(t *testing.T) {
	for _, tc := range []struct {
		sql        string
		tables     []string
		partitions []string
		global     bool
	}{
		{sql: "drop stats t", tables: []string{"t"}},
		{sql: "drop stats db.t", tables: []string{"db.t"}},
		{sql: "drop stats t1, t2, t3", tables: []string{"t1", "t2", "t3"}},
		{sql: "drop stats t global", tables: []string{"t"}, global: true},
		{sql: "drop stats t partition p0", tables: []string{"t"}, partitions: []string{"p0"}},
		{sql: "drop stats t partition p0, p1, p2", tables: []string{"t"}, partitions: []string{"p0", "p1", "p2"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneStmt(t, tc.sql).(*ast.DropStatsStmt)
			if !ok {
				t.Fatalf("Parse(%q) gave a statement of another kind", tc.sql)
			}
			wantList(t, "Tables", tableNames(stmt.Tables), tc.tables)
			wantList(t, "Partitions", stmt.Partitions, tc.partitions)
			if stmt.Global != tc.global {
				t.Errorf("Global = %v, want %v", stmt.Global, tc.global)
			}
		})
	}
}

// TestLoadStatsTree reads back the path LOAD STATS names, unescaped, as the
// consumer that would open it needs it.
func TestLoadStatsTree(t *testing.T) {
	for _, tc := range []struct{ sql, path string }{
		{"load stats '/tmp/stats.json'", "/tmp/stats.json"},
		{`load stats "x"`, "x"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneStmt(t, tc.sql).(*ast.LoadStatsStmt)
			if !ok {
				t.Fatalf("Parse(%q) gave a statement of another kind", tc.sql)
			}
			if stmt.Path != tc.path {
				t.Errorf("Path = %q, want %q", stmt.Path, tc.path)
			}
		})
	}
}

// TestCreateStatisticsTree reads back the name, the type, the table and the
// column list. The three types are three values: a reader that took any word
// in the parentheses would make CARDINALITY and CORRELATION the same
// statistic.
func TestCreateStatisticsTree(t *testing.T) {
	for _, tc := range []struct {
		sql         string
		ifNotExists bool
		name        string
		typ         ast.StatisticsType
		table       string
		columns     []string
	}{
		{sql: "create statistics s1 (cardinality) on t(a,b,c)", name: "s1", typ: ast.StatisticsCardinality, table: "t", columns: []string{"a", "b", "c"}},
		{sql: "create statistics s1(cardinality) on t(a)", name: "s1", typ: ast.StatisticsCardinality, table: "t", columns: []string{"a"}},
		{sql: "create statistics s2 (dependency) on t(a,b)", name: "s2", typ: ast.StatisticsDependency, table: "t", columns: []string{"a", "b"}},
		{sql: "create statistics s3 (correlation) on t(a,b)", name: "s3", typ: ast.StatisticsCorrelation, table: "t", columns: []string{"a", "b"}},
		{sql: "create statistics s1 (CARDINALITY) on t(a)", name: "s1", typ: ast.StatisticsCardinality, table: "t", columns: []string{"a"}},
		{sql: "create statistics `s1` (cardinality) on t(a)", name: "s1", typ: ast.StatisticsCardinality, table: "t", columns: []string{"a"}},
		{sql: "create statistics if not exists s1 (cardinality) on t(a)", ifNotExists: true, name: "s1", typ: ast.StatisticsCardinality, table: "t", columns: []string{"a"}},
		{sql: "create statistics s1 (cardinality) on db.t(a)", name: "s1", typ: ast.StatisticsCardinality, table: "db.t", columns: []string{"a"}},
		// A qualified column keeps its qualifier, which ANALYZE's COLUMNS
		// list cannot carry at all.
		{sql: "create statistics s1 (cardinality) on t(t.a)", name: "s1", typ: ast.StatisticsCardinality, table: "t", columns: []string{"t.a"}},
		{sql: "create statistics s1 (cardinality) on t(d.t.a)", name: "s1", typ: ast.StatisticsCardinality, table: "t", columns: []string{"d.t.a"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneStmt(t, tc.sql).(*ast.CreateStatisticsStmt)
			if !ok {
				t.Fatalf("Parse(%q) gave a statement of another kind", tc.sql)
			}
			if stmt.IfNotExists != tc.ifNotExists {
				t.Errorf("IfNotExists = %v, want %v", stmt.IfNotExists, tc.ifNotExists)
			}
			if stmt.Name != tc.name {
				t.Errorf("Name = %q, want %q", stmt.Name, tc.name)
			}
			if stmt.Type != tc.typ {
				t.Errorf("Type = %s, want %s", stmt.Type, tc.typ)
			}
			if stmt.Table == nil {
				t.Fatal("Table is nil")
			}
			wantList(t, "Table", tableNames([]*ast.TableRef{stmt.Table}), []string{tc.table})
			var cols []string
			for _, c := range stmt.Columns {
				var b strings.Builder
				if c.Schema != "" {
					b.WriteString(c.Schema + ".")
				}
				if c.Table != "" {
					b.WriteString(c.Table + ".")
				}
				b.WriteString(c.Column)
				cols = append(cols, b.String())
			}
			wantList(t, "Columns", cols, tc.columns)
		})
	}
}

// TestDropStatisticsTree reads back the one name DROP STATISTICS takes.
func TestDropStatisticsTree(t *testing.T) {
	for _, tc := range []struct{ sql, name string }{
		{"drop statistics s1", "s1"},
		{"drop statistics `s1`", "s1"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt, ok := parseOneStmt(t, tc.sql).(*ast.DropStatisticsStmt)
			if !ok {
				t.Fatalf("Parse(%q) gave a statement of another kind", tc.sql)
			}
			if stmt.Name != tc.name {
				t.Errorf("Name = %q, want %q", stmt.Name, tc.name)
			}
		})
	}
}

// TestStatisticsStatementsAreNotTableStatements pins what the dispatcher
// decides: DROP STATS t is not a DROP TABLE of a table called stats, LOAD
// STATS is not a LOAD DATA, and `DROP TABLE stats` is still a DROP TABLE. A
// loose dispatch here is what shipped LOCK STATS t as a LOCK TABLES: the
// corpus went green and the statement was the wrong kind.
func TestStatisticsStatementsAreNotTableStatements(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want string
	}{
		{"drop stats t", "*ast.DropStatsStmt"},
		{"drop statistics s1", "*ast.DropStatisticsStmt"},
		{"drop table stats", "*ast.DropTableStmt"},
		{"drop table statistics", "*ast.DropTableStmt"},
		{"load stats 'x'", "*ast.LoadStatsStmt"},
		{"create statistics s1 (cardinality) on t(a)", "*ast.CreateStatisticsStmt"},
		{"create table statistics (a int)", "*ast.CreateTableStmt"},
		{"analyze table t index a", "*ast.AnalyzeTableStmt"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			if got := fmt.Sprintf("%T", parseOneStmt(t, tc.sql)); got != tc.want {
				t.Errorf("Parse(%q) gave %s, want %s", tc.sql, got, tc.want)
			}
		})
	}
}
