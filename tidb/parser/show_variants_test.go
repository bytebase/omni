package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/deparse"
)

// The shape TiDB's own SHOW statements parse into. The corpus in this package
// records only whether a statement parses, so it cannot see a SHOW whose
// target or filter landed in the wrong field: that parses clean and misreports
// for as long as it stands. These assertions are the other half, and they run
// without a container.

// showShape is the part of a ShowStmt a form is expected to fill. A zero field
// is expected to stay zero.
type showShape struct {
	typ        string
	schema     string
	name       string
	partitions []string
	indexName  string
	partition  string
	policyName string
	groupName  string
	jobID      int64
	like       string // the text of Like, deparsed, or "" for none
	where      string // the text of Where, deparsed, or "" for none
}

func TestShowTiDBVariantShapes(t *testing.T) {
	cases := []struct {
		sql  string
		want showShape
	}{
		// SHOW TABLE ... REGIONS. The partition list rides on the table
		// reference, as it does upstream, and REGIONS takes a WHERE and no
		// LIKE.
		{"show table t1 regions", showShape{typ: "TABLE REGIONS", name: "t1"}},
		{"show table test.t1 regions", showShape{typ: "TABLE REGIONS", schema: "test", name: "t1"}},
		{"show table t1 regions where a=1", showShape{typ: "TABLE REGIONS", name: "t1", where: "(`a` = 1)"}},
		{"show table t1 partition (p0) regions", showShape{typ: "TABLE REGIONS", name: "t1", partitions: []string{"p0"}}},
		{
			"show table t1 partition (p0,p1) regions where a=1",
			showShape{typ: "TABLE REGIONS", name: "t1", partitions: []string{"p0", "p1"}, where: "(`a` = 1)"},
		},
		// The index name is its own field, not the table's.
		{"show table t1 index idx1 regions", showShape{typ: "TABLE INDEX REGIONS", name: "t1", indexName: "idx1"}},
		{
			"show table t1 partition (p0,p1) index idx1 regions where a=2",
			showShape{typ: "TABLE INDEX REGIONS", name: "t1", partitions: []string{"p0", "p1"}, indexName: "idx1", where: "(`a` = 2)"},
		},
		// SHOW TABLE ... DISTRIBUTIONS, same shape as REGIONS.
		{"show table t1 distributions", showShape{typ: "TABLE DISTRIBUTIONS", name: "t1"}},
		{"show table t1 distributions where a=1", showShape{typ: "TABLE DISTRIBUTIONS", name: "t1", where: "(`a` = 1)"}},
		{
			"show table t1 partition (p0,p1) distributions",
			showShape{typ: "TABLE DISTRIBUTIONS", name: "t1", partitions: []string{"p0", "p1"}},
		},
		// SHOW TABLE ... NEXT_ROW_ID, which takes neither a partition list
		// nor a filter.
		{"show table t1 next_row_id", showShape{typ: "TABLE NEXT_ROW_ID", name: "t1"}},
		{"show table t1.t1 next_row_id", showShape{typ: "TABLE NEXT_ROW_ID", schema: "t1", name: "t1"}},

		// SHOW DISTRIBUTION JOB{,S}. The id is an integer literal of its own
		// and the plural form is filterable while the singular is not.
		{"show distribution jobs", showShape{typ: "DISTRIBUTION JOBS"}},
		{"show distribution jobs where id > 0", showShape{typ: "DISTRIBUTION JOBS", where: "(`id` > 0)"}},
		{"show distribution job 1", showShape{typ: "DISTRIBUTION JOB", jobID: 1}},
		{"show distribution job 42", showShape{typ: "DISTRIBUTION JOB", jobID: 42}},

		// SHOW PLACEMENT and its FOR targets. The database name lands on the
		// table reference's Name, as SHOW CREATE DATABASE's does; the single
		// partition of the FOR TABLE form lands on Partition, not on the
		// reference's partition list.
		{"SHOW PLACEMENT", showShape{typ: "PLACEMENT"}},
		{"SHOW PLACEMENT LIKE 'POLICY foo%'", showShape{typ: "PLACEMENT", like: "POLICY foo%"}},
		{"SHOW PLACEMENT WHERE Target='TABLE test.t1'", showShape{typ: "PLACEMENT", where: "(`Target` = 'TABLE test.t1')"}},
		{"SHOW PLACEMENT LABELS", showShape{typ: "PLACEMENT LABELS"}},
		{"SHOW PLACEMENT LABELS LIKE '%zone%'", showShape{typ: "PLACEMENT LABELS", like: "%zone%"}},
		{"SHOW PLACEMENT LABELS WHERE label='l123'", showShape{typ: "PLACEMENT LABELS", where: "(`label` = 'l123')"}},
		{"SHOW PLACEMENT FOR DATABASE db1", showShape{typ: "PLACEMENT FOR DATABASE", name: "db1"}},
		{"SHOW PLACEMENT FOR SCHEMA db1", showShape{typ: "PLACEMENT FOR DATABASE", name: "db1"}},
		{"SHOW PLACEMENT FOR TABLE tb1", showShape{typ: "PLACEMENT FOR TABLE", name: "tb1"}},
		{"SHOW PLACEMENT FOR TABLE db1.tb1", showShape{typ: "PLACEMENT FOR TABLE", schema: "db1", name: "tb1"}},
		{
			"SHOW PLACEMENT FOR TABLE tb1 PARTITION p1",
			showShape{typ: "PLACEMENT FOR PARTITION", name: "tb1", partition: "p1"},
		},
		{
			"SHOW PLACEMENT FOR TABLE db1.tb1 PARTITION p1",
			showShape{typ: "PLACEMENT FOR PARTITION", schema: "db1", name: "tb1", partition: "p1"},
		},

		// SHOW CREATE SEQUENCE, PLACEMENT POLICY and RESOURCE GROUP. The
		// policy and group names are fields of their own: neither is a table.
		{"show create sequence seq", showShape{typ: "CREATE SEQUENCE", name: "seq"}},
		{"show create sequence test.seq", showShape{typ: "CREATE SEQUENCE", schema: "test", name: "seq"}},
		{"show create placement policy x", showShape{typ: "CREATE PLACEMENT POLICY", policyName: "x"}},
		{"show create placement policy `placement`", showShape{typ: "CREATE PLACEMENT POLICY", policyName: "placement"}},
		{"SHOW CREATE RESOURCE GROUP rg", showShape{typ: "CREATE RESOURCE GROUP", groupName: "rg"}},
		{"SHOW CREATE RESOURCE GROUP default", showShape{typ: "CREATE RESOURCE GROUP", groupName: "DEFAULT"}},

		// The statistics targets, all filterable.
		{"show stats_extended", showShape{typ: "STATS_EXTENDED"}},
		{"show stats_extended where table_name = 't'", showShape{typ: "STATS_EXTENDED", where: "(`table_name` = 't')"}},
		{"show stats_meta", showShape{typ: "STATS_META"}},
		{"show stats_meta where table_name = 't'", showShape{typ: "STATS_META", where: "(`table_name` = 't')"}},
		{"show stats_locked", showShape{typ: "STATS_LOCKED"}},
		{"show stats_histograms", showShape{typ: "STATS_HISTOGRAMS"}},
		{"show stats_histograms where col_name = 'a'", showShape{typ: "STATS_HISTOGRAMS", where: "(`col_name` = 'a')"}},
		{"show stats_buckets", showShape{typ: "STATS_BUCKETS"}},
		{"show stats_healthy", showShape{typ: "STATS_HEALTHY"}},
		{"show stats_topn", showShape{typ: "STATS_TOPN"}},
		{"show histograms_in_flight", showShape{typ: "HISTOGRAMS_IN_FLIGHT"}},
		{"show column_stats_usage", showShape{typ: "COLUMN_STATS_USAGE"}},
		{"show column_stats_usage where table_name = 't'", showShape{typ: "COLUMN_STATS_USAGE", where: "(`table_name` = 't')"}},
		{"show analyze status", showShape{typ: "ANALYZE STATUS"}},
		{"show analyze status where table_name = 't'", showShape{typ: "ANALYZE STATUS", where: "(`table_name` = 't')"}},
		{"show analyze status where table_name like '%'", showShape{typ: "ANALYZE STATUS", where: "(`table_name` like '%')"}},

		// The SQL-binding targets, and the scope prefix they read.
		{"show bindings", showShape{typ: "BINDINGS"}},
		{"show global bindings", showShape{typ: "GLOBAL BINDINGS"}},
		{"show session bindings", showShape{typ: "SESSION BINDINGS"}},
		{"show binding_cache status", showShape{typ: "BINDING_CACHE STATUS"}},

		// The rest of TiDB's filterable targets, plus the one that is not.
		{"show session_states", showShape{typ: "SESSION_STATES"}},
		{"show config", showShape{typ: "CONFIG"}},
		{"show config where type='tidb'", showShape{typ: "CONFIG", where: "(`type` = 'tidb')"}},
		{"show backups", showShape{typ: "BACKUPS"}},
		{"show restores like 'r0001'", showShape{typ: "RESTORES", like: "r0001"}},
		{"show builtins", showShape{typ: "BUILTINS"}},

		// SHOW BINARY LOG STATUS, MySQL 8.2's name for SHOW MASTER STATUS.
		{"SHOW BINARY LOG STATUS", showShape{typ: "BINARY LOG STATUS"}},
	}

	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneShow(t, tc.sql)
			assertShowShape(t, stmt, tc.want)
		})
	}
}

// TestShowProfileLimitOffsetShape pins which of the two LIMIT numbers is the
// offset. SHOW PROFILE reads a whole limit clause, where `LIMIT a, b` means
// offset a and count b while `LIMIT a OFFSET b` means count a and offset b;
// reading them the same way round would have swapped the pair silently.
func TestShowProfileLimitOffsetShape(t *testing.T) {
	stmt := parseOneShow(t, "SHOW PROFILE CPU FOR QUERY 2 LIMIT 1,3")
	if stmt.Type != "PROFILE" {
		t.Fatalf("Type = %q, want PROFILE", stmt.Type)
	}
	if got := intLitValue(t, stmt.LimitOffset); got != 1 {
		t.Errorf("LimitOffset = %d, want 1", got)
	}
	if got := intLitValue(t, stmt.LimitCount); got != 3 {
		t.Errorf("LimitCount = %d, want 3", got)
	}

	stmt = parseOneShow(t, "SHOW PROFILE CPU LIMIT 1 OFFSET 3")
	if got := intLitValue(t, stmt.LimitCount); got != 1 {
		t.Errorf("LimitCount = %d, want 1", got)
	}
	if got := intLitValue(t, stmt.LimitOffset); got != 3 {
		t.Errorf("LimitOffset = %d, want 3", got)
	}
}

// TestShowVariantsReachNodeToString pins that the new fields are written by
// the AST printer, so a change that drops one of them shows up here rather
// than in a silent diff of someone else's golden file.
func TestShowVariantsReachNodeToString(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{
			"show table t1 partition (p0,p1) index idx1 regions where a=2",
			`{SHOW :loc 0 :type TABLE INDEX REGIONS :from {TABLEREF :loc 11 :name t1 :partitions (p0 p1)} :where {BINEXPR :op = :loc 57 :left {COLREF :loc 57 :col a} :right {INT_LIT :val 2 :loc 59}} :index_name idx1}`,
		},
		{
			"SHOW PLACEMENT FOR TABLE db1.tb1 PARTITION p1",
			`{SHOW :loc 0 :type PLACEMENT FOR PARTITION :from {TABLEREF :loc 25 :schema db1 :name tb1} :partition p1}`,
		},
		{
			"show create placement policy x",
			`{SHOW :loc 0 :type CREATE PLACEMENT POLICY :policy_name x}`,
		},
		{
			"SHOW CREATE RESOURCE GROUP rg",
			`{SHOW :loc 0 :type CREATE RESOURCE GROUP :resource_group_name rg}`,
		},
		{
			"show distribution job 7",
			`{SHOW :loc 0 :type DISTRIBUTION JOB :job_id {INT_LIT :val 7 :loc 22}}`,
		},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			ParseAndCompare(t, tc.sql, tc.want)
		})
	}
}

// TestShowVariantsAreWalked pins that ast.Walk reaches the expressions the new
// forms hold. The parser measures every statement's tree height with Walk
// before it returns, so a child Walk does not visit is a child that cannot be
// counted -- and the review rules that read a statement's WHERE find it
// through Walk too.
func TestShowVariantsAreWalked(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		want int // expressions Walk should reach
	}{
		{"show table t1 regions where a=1", 3},   // BinExpr, ColumnRef, IntLit
		{"show distribution job 7", 1},           // IntLit
		{"show config where type='tidb'", 3},     // BinExpr, ColumnRef, StringLit
		{"show restores like 'r0001'", 1},        // StringLit
		{"SHOW PROFILE CPU LIMIT 1 OFFSET 3", 2}, // two IntLits
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneShow(t, tc.sql)
			var exprs int
			nodes.Inspect(stmt, func(n nodes.Node) bool {
				if _, ok := n.(nodes.ExprNode); ok {
					exprs++
				}
				return true
			})
			if exprs != tc.want {
				t.Errorf("Walk reached %d expressions, want %d", exprs, tc.want)
			}
		})
	}
}

// parseOneShow parses sql through the whole entry point and returns the one
// SHOW statement it holds. Going through Parse rather than parseShowStmt is
// deliberate: a statement that leaves tokens behind is a second statement in
// the list, which is the shape #629's bug took, and asserting on one item
// catches it.
func parseOneShow(t *testing.T, sql string) *nodes.ShowStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*nodes.ShowStmt)
	if !ok {
		t.Fatalf("Parse(%q) returned %T, want *ast.ShowStmt", sql, list.Items[0])
	}
	return stmt
}

func assertShowShape(t *testing.T, stmt *nodes.ShowStmt, want showShape) {
	t.Helper()
	if stmt.Type != want.typ {
		t.Errorf("Type = %q, want %q", stmt.Type, want.typ)
	}
	if want.name == "" {
		if stmt.From != nil {
			t.Errorf("From = %+v, want nil", stmt.From)
		}
	} else {
		switch {
		case stmt.From == nil:
			t.Errorf("From is nil, want name %q", want.name)
		default:
			if stmt.From.Name != want.name {
				t.Errorf("From.Name = %q, want %q", stmt.From.Name, want.name)
			}
			if stmt.From.Schema != want.schema {
				t.Errorf("From.Schema = %q, want %q", stmt.From.Schema, want.schema)
			}
			if strings.Join(stmt.From.Partitions, ",") != strings.Join(want.partitions, ",") {
				t.Errorf("From.Partitions = %v, want %v", stmt.From.Partitions, want.partitions)
			}
		}
	}
	if stmt.IndexName != want.indexName {
		t.Errorf("IndexName = %q, want %q", stmt.IndexName, want.indexName)
	}
	if stmt.Partition != want.partition {
		t.Errorf("Partition = %q, want %q", stmt.Partition, want.partition)
	}
	if stmt.PolicyName != want.policyName {
		t.Errorf("PolicyName = %q, want %q", stmt.PolicyName, want.policyName)
	}
	if stmt.ResourceGroupName != want.groupName {
		t.Errorf("ResourceGroupName = %q, want %q", stmt.ResourceGroupName, want.groupName)
	}
	if want.jobID == 0 {
		if stmt.JobID != nil {
			t.Errorf("JobID = %v, want nil", stmt.JobID)
		}
	} else if got := intLitValue(t, stmt.JobID); got != want.jobID {
		t.Errorf("JobID = %d, want %d", got, want.jobID)
	}
	if want.like == "" {
		if stmt.Like != nil {
			t.Errorf("Like = %v, want nil", stmt.Like)
		}
	} else {
		lit, ok := stmt.Like.(*nodes.StringLit)
		if !ok {
			t.Errorf("Like = %T, want *ast.StringLit", stmt.Like)
		} else if lit.Value != want.like {
			t.Errorf("Like = %q, want %q", lit.Value, want.like)
		}
	}
	if want.where == "" {
		if stmt.Where != nil {
			t.Errorf("Where = %v, want nil", stmt.Where)
		}
	} else {
		if stmt.Where == nil {
			t.Errorf("Where is nil, want %q", want.where)
		} else if got := deparse.Deparse(stmt.Where); got != want.where {
			t.Errorf("Where text = %q, want %q", got, want.where)
		}
	}
}

func intLitValue(t *testing.T, e nodes.ExprNode) int64 {
	t.Helper()
	lit, ok := e.(*nodes.IntLit)
	if !ok {
		t.Fatalf("expression is %T, want *ast.IntLit", e)
	}
	return lit.Value
}

// TestShowEndsWhereItsTargetEnds pins that a SHOW reads to the end of its
// statement. Split has already cut on the delimiter, so anything left in the
// segment is a second statement, and the statement loop takes one without
// needing a separator between them: every form here parsed as a SHOW plus a
// `TABLE t`, which is a 1064 on TiDB v8.5.5.
//
// The verdicts are in showTargetCases, measured against the server. What this
// test adds is the shape: a boolean accept/reject cannot tell a statement
// that swallowed its trailing text from one that left it for the loop, and it
// was the second. So the assertion is on the statement count.
func TestShowEndsWhereItsTargetEnds(t *testing.T) {
	for _, sql := range []string{
		"SHOW PLACEMENT FOR DATABASE db1 TABLE tb1",
		"SHOW PLACEMENT FOR SCHEMA db1 TABLE tb1",
		"SHOW PLACEMENT FOR TABLE tb1 TABLE tb2",
		"SHOW PLACEMENT TABLE t",
		"SHOW DATABASES TABLE t",
		"SHOW TABLES TABLE t",
		"SHOW CREATE TABLE t TABLE t2",
		"SHOW TABLE STATUS TABLE t",
		"SHOW TABLE t1 REGIONS TABLE t2",
		"SHOW TABLE t1 NEXT_ROW_ID TABLE t2",
		"SHOW BUILTINS TABLE t",
		"SHOW DISTRIBUTION JOB 1 TABLE t",
		"SHOW STATS_META TABLE t",
		"SHOW VARIABLES LIKE 'x%' TABLE t",
	} {
		t.Run(sql, func(t *testing.T) {
			list, err := Parse(sql)
			if err == nil {
				t.Errorf("Parse(%q) accepted as %d statements; TiDB v8.5.5 answers 1064", sql, len(list.Items))
			}
		})
	}
	// What the guard does not refuse: the delimiter Split leaves behind, and
	// a second statement that is separated from the first.
	for _, sql := range []string{
		"SHOW TABLES",
		"SHOW TABLES;",
		"SHOW TABLES ;",
		"SHOW TABLES /* trailing comment */",
		"CREATE PROCEDURE p() BEGIN SHOW TABLES; END",
		"CREATE PROCEDURE p() BEGIN IF 1 THEN SHOW TABLES; END IF; END",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
	// Two statements with a separator stay two statements.
	list, err := Parse("SHOW TABLES; SHOW DATABASES")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(list.Items) != 2 {
		t.Errorf("got %d statements, want 2", len(list.Items))
	}
}

// TestShowLikePatternShape pins what a SHOW filter's LIKE holds. TiDB reads a
// SimpleExpr there, so the pattern may carry a COLLATE suffix, and the
// collation is a CollationName: a quoted name and the keyword BINARY are as
// good as an identifier. The suffix wraps the pattern rather than replacing
// it, so a reader still finds the string underneath.
func TestShowLikePatternShape(t *testing.T) {
	for _, tc := range []struct{ sql, pattern, collation string }{
		{"SHOW DATABASES LIKE 't%'", "t%", ""},
		{"SHOW DATABASES LIKE 't%' COLLATE utf8mb4_bin", "t%", "utf8mb4_bin"},
		{"SHOW DATABASES LIKE 't%' COLLATE 'utf8mb4_bin'", "t%", "utf8mb4_bin"},
		{"SHOW DATABASES LIKE 't%' COLLATE binary", "t%", "binary"},
		{"SHOW TABLES LIKE 't%' COLLATE utf8mb4_bin", "t%", "utf8mb4_bin"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneShow(t, tc.sql)
			like := stmt.Like
			if tc.collation != "" {
				coll, ok := like.(*nodes.CollateExpr)
				if !ok {
					t.Fatalf("Like = %T, want *ast.CollateExpr", like)
				}
				if coll.Collation != tc.collation {
					t.Errorf("Collation = %q, want %q", coll.Collation, tc.collation)
				}
				like = coll.Expr
			}
			lit, ok := like.(*nodes.StringLit)
			if !ok {
				t.Fatalf("pattern is %T, want *ast.StringLit", like)
			}
			if lit.Value != tc.pattern {
				t.Errorf("pattern = %q, want %q", lit.Value, tc.pattern)
			}
			if stmt.Where != nil {
				t.Errorf("Where = %v, want nil", stmt.Where)
			}
		})
	}
	// What a SimpleExpr is not: the comparison, boolean, arithmetic, bitwise
	// and JSON operators of a whole Expression, and a COLLATE with no name.
	for _, sql := range []string{
		"SHOW DATABASES LIKE 't' LIKE 'u'",
		"SHOW DATABASES LIKE 't' AND 1",
		"SHOW DATABASES LIKE 'a' || 'b'",
		"SHOW DATABASES LIKE 'a' + 1",
		"SHOW DATABASES LIKE 'a' * 2",
		"SHOW DATABASES LIKE 'a' | 1",
		"SHOW DATABASES LIKE 'a' -> '$.x'",
		"SHOW DATABASES LIKE 't%' COLLATE",
		"SHOW DATABASES LIKE 't%' COLLATE 1",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("Parse(%q) accepted; TiDB v8.5.5 answers 1064", sql)
			}
		})
	}
}
