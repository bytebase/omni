package parser

import "testing"

// The forms PLAN REPLAYER takes, and -- the point of this file -- the forms
// it does not.
//
// The corpus in testdata/pingcap_parser_corpus.tsv pins what pingcap's own
// parser test writes, and for PLAN REPLAYER that is sixteen accepting rows
// and no rejection at all. The statement has four shapes behind one opener,
// and a reader that took any word after DUMP EXPLAIN would turn all sixteen
// green while accepting PLAN REPLAYER CAPTURE with one digest, LOAD with two
// files, and a DUMP of a statement TiDB will not explain, so the rejections
// are written here by hand: each arm's arity, the statements outside
// upstream's ExplainableStmt, the WITH STATS clause written short, and a file
// name that is not quoted.
//
// Every case below was measured against pingcap/tidb:v8.5.5 through PREPARE
// on 2026-11-04, and the recorded verdict is what that server answered.
// TestPlanReplayerFormsMatchTiDB is what keeps it a recording: an image
// upgrade that changes an answer fails there rather than leaving a stale
// expectation behind.
//
// Ref: PlanReplayerStmt, PlanReplayerDumpOpt and ExplainableStmt in
// pingcap/tidb pkg/parser/parser.y.

// planReplayerFormCase is a statement and whether TiDB's parser accepts it.
type planReplayerFormCase struct {
	sql    string
	accept bool
}

var planReplayerFormCases = []planReplayerFormCase{
	// --- DUMP EXPLAIN <stmt>, with and without ANALYZE ---
	{"plan replayer dump explain select 1", true},
	{"plan replayer dump explain analyze select 1", true},
	{"plan replayer dump explain select a from t", true},
	{"plan replayer dump explain select * from t where a > 10", true},

	// --- the explainable statements, which are EXPLAIN's set exactly ---
	{"plan replayer dump explain table t", true},
	{"plan replayer dump explain values row(1)", true},
	// A TABLE or VALUES target is the query expression it begins, so a set
	// operation over two of them is one explainable statement.
	{"plan replayer dump explain table t1 union table t2", true},
	{"plan replayer dump explain values row(1) union values row(2)", true},
	{"plan replayer dump explain table t1 union select 1", true},
	{"plan replayer dump explain analyze table t1 union table t2", true},
	{"plan replayer dump explain with c as (select 1) select * from c", true},
	{"plan replayer dump explain (select 1)", true},
	{"plan replayer dump explain select 1 union select 2", true},
	{"plan replayer dump explain insert into t values (1)", true},
	{"plan replayer dump explain replace into t values (1)", true},
	{"plan replayer dump explain update t set a = 1", true},
	{"plan replayer dump explain delete from t", true},
	// ALTER TABLE is the one DDL in the set, and the only reason the
	// catalog has to read what a dump carries.
	{"plan replayer dump explain alter table t add column c int", true},
	{"plan replayer dump explain analyze alter table t add column c int", true},
	// IMPORT INTO is the other half of the set, and is read now that this
	// parser has the statement.
	{"plan replayer dump explain import into t from '/tmp/a.csv'", true},
	{"plan replayer dump explain analyze import into t from '/tmp/a.csv'", true},
	{"plan replayer dump explain import into t from table t1", true},
	// Which is TRACE's set almost inverted: TraceableStmt holds the
	// transaction statements, SET and LOAD DATA and holds no DDL, and
	// ExplainableStmt holds none of those three and does hold ALTER TABLE.
	{"plan replayer dump explain begin", false},
	{"plan replayer dump explain set a = 1", false},
	{"plan replayer dump explain load data infile '/tmp/a.csv' into table t", false},
	{"plan replayer dump explain create table t2 (a int)", false},
	{"plan replayer dump explain truncate table t", false},
	{"plan replayer dump explain alter database d charset utf8mb4", false},
	{"plan replayer dump explain drop table t", false},
	{"plan replayer dump explain show tables", false},
	{"plan replayer dump explain explain select 1", false},
	{"plan replayer dump explain t", false},
	// ANALYZE TABLE is not explainable; this parses as EXPLAIN ANALYZE of
	// TABLE t, which is the arm that makes the two readings collide.
	{"plan replayer dump explain analyze table t", true},

	// --- DUMP EXPLAIN [ANALYZE] SLOW QUERY, with its three optional tails ---
	{"plan replayer dump explain slow query", true},
	{"plan replayer dump explain analyze slow query", true},
	{"plan replayer dump explain slow query limit 10", true},
	{"plan replayer dump explain slow query limit 1, 10", true},
	{"plan replayer dump explain slow query where a > 10 limit 10 offset 5", true},
	{"plan replayer dump explain slow query where a > 10 and t < 1 limit 10", true},
	{"plan replayer dump explain slow query where a > 10 order by t", true},
	{"plan replayer dump explain slow query order by t", true},
	{"plan replayer dump explain slow query where a > 10 and t < 1 order by t limit 10", true},
	{"plan replayer dump explain analyze slow query where a > 10 and t < 1 order by t limit 10", true},
	// SLOW is followed by QUERY and nothing else, and the three tails come
	// in this order.
	{"plan replayer dump explain slow", false},
	{"plan replayer dump explain slow queries", false},
	{"plan replayer dump explain slow query limit 10 where a > 10", false},
	{"plan replayer dump explain slow query order by t where a > 10", false},
	{"plan replayer dump explain slow query where", false},

	// --- DUMP EXPLAIN [ANALYZE] 'file' ---
	{"plan replayer dump explain 'sql.txt'", true},
	{"plan replayer dump explain analyze 'sql.txt'", true},
	// The file name is a string literal: a backquoted name is no part of
	// this grammar, and the arm takes one name.
	{"plan replayer dump explain `sql.txt`", false},
	{"plan replayer dump explain analyze 'a' 'b'", false},

	// --- DUMP WITH STATS AS OF TIMESTAMP <expr>, which stands in front of
	// --- EXPLAIN and reaches all three of its arms ---
	{"plan replayer dump with stats as of timestamp '2024-01-01 00:00:00' explain select 1", true},
	{"plan replayer dump with stats as of timestamp now() explain analyze select 1", true},
	{"plan replayer dump with stats as of timestamp '2024-01-01' explain slow query limit 10", true},
	{"plan replayer dump with stats as of timestamp '2024-01-01' explain 'sql.txt'", true},
	// The clause is WITH STATS AS OF TIMESTAMP expr whole; no part of it is
	// optional once WITH is written.
	{"plan replayer dump with stats explain select 1", false},
	{"plan replayer dump with stats as of explain select 1", false},
	{"plan replayer dump with stats as of timestamp explain select 1", false},
	{"plan replayer dump with explain select 1", false},
	{"plan replayer dump as of timestamp '2024-01-01' explain select 1", false},
	// And it belongs to DUMP, not to LOAD or CAPTURE.
	{"plan replayer with stats as of timestamp '2024-01-01' dump explain select 1", false},

	// --- LOAD 'file' ---
	{"plan replayer load '/tmp/sdfaalskdjf.zip'", true},
	{"plan replayer load", false},
	{"plan replayer load 'a.zip' 'b.zip'", false},
	{"plan replayer load `a.zip`", false},

	// --- CAPTURE and CAPTURE REMOVE, each taking exactly two digests ---
	{"plan replayer capture '123' '123'", true},
	{"plan replayer capture remove '123' '123'", true},
	{"plan replayer capture '123'", false},
	{"plan replayer capture", false},
	{"plan replayer capture '123' '123' '123'", false},
	{"plan replayer capture remove '123'", false},
	{"plan replayer capture remove", false},
	{"plan replayer capture 123 123", false},

	// --- the opener itself ---
	{"plan replayer", false},
	{"plan replayer dump", false},
	{"plan replayer dump explain", false},
	{"plan replayer dump select 1", false},
	{"plan replayer explain select 1", false},
	{"plan replayer dumping explain select 1", false},
	{"plan replayers dump explain select 1", false},
	{"plan dump explain select 1", false},
	{"plan", false},
	{"replayer dump explain select 1", false},
	// Each word of the opener is read bare: a backquoted one is no part of
	// this grammar.
	{"plan `replayer` dump explain select 1", false},
	{"`plan` replayer dump explain select 1", false},
	{"plan replayer `dump` explain select 1", false},
	{"plan replayer `capture` '123' '123'", false},
}

// TestPlanReplayerForms holds this parser to the verdicts above. It needs no
// container.
func TestPlanReplayerForms(t *testing.T) {
	for _, c := range planReplayerFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni %ss %q, TiDB %s %ss it (err=%v)",
					verdictName(accepts), c.sql, tidbCorpusVersion, verdictName(c.accept), err)
			}
		})
	}
}

// planReplayerWordsAsNames are statements in which the words PLAN REPLAYER
// reads are ordinary names instead. REPLAYER, DUMP, CAPTURE and STATS are not
// registered in the lexer at all -- this family reads each with isBareWord --
// and these statements are what says so. PLAN is registered, for the SQL plan
// bindings, and is listed here because the dispatcher's new arm for it must
// not change that: PLAN starts a statement only when REPLAYER follows.
//
// Every statement here parses on v8.5.5 too, which
// TestPlanReplayerWordsStayNamesOnTiDB measures.
var planReplayerWordsAsNames = []string{
	// As a table name, bare and qualified.
	"SELECT * FROM replayer",
	"SELECT * FROM d.replayer",
	"SELECT * FROM dump",
	"SELECT * FROM capture",
	"SELECT * FROM plan",
	// As an implicit table alias, the site that broke when the SHOW
	// keywords were registered.
	"DELETE FROM t replayer WHERE replayer.id = 1",
	"DELETE FROM t dump WHERE dump.id = 1",
	"DELETE FROM t capture WHERE capture.id = 1",
	"DELETE FROM t plan WHERE plan.id = 1",
	"UPDATE t replayer SET replayer.a = 1",
	"SELECT dump.a FROM t dump",
	// As an explicit alias and a column name.
	"SELECT 1 AS replayer",
	"SELECT 1 AS dump",
	"SELECT 1 AS capture",
	"SELECT replayer, dump, capture FROM t",
	"SELECT t.plan FROM t",
	// As names in DDL, which is where a registered keyword hurts most.
	"CREATE TABLE replayer (dump int, capture int, plan int, stats int)",
	"CREATE TABLE t (replayer int, INDEX dump (replayer))",
	"CREATE INDEX replayer ON t (a)",
	"DROP TABLE replayer",
	"ALTER TABLE t ADD COLUMN replayer int",
	"ALTER TABLE plan RENAME TO replayer",
	// And beside the statement of this family.
	"PLAN REPLAYER DUMP EXPLAIN SELECT * FROM replayer",
	"PLAN REPLAYER DUMP EXPLAIN SELECT dump.a FROM t dump",
	"PLAN REPLAYER DUMP EXPLAIN DELETE FROM t plan WHERE plan.id = 1",
}

// TestPlanReplayerWordsStayNames holds the other half of the keyword
// decision: every word this family reads bare is still an ordinary name
// everywhere else, and so is PLAN. It needs no container.
func TestPlanReplayerWordsStayNames(t *testing.T) {
	for _, sql := range planReplayerWordsAsNames {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
}
