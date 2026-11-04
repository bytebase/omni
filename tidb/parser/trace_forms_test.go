package parser

import "testing"

// The forms TRACE takes, and -- the point of this file -- the forms it does
// not.
//
// The corpus in testdata/pingcap_parser_corpus.tsv pins what pingcap's own
// parser test writes, and for TRACE that is fifteen accepting rows and no
// rejection at all. A reader that took any word after TRACE and handed the
// rest to the statement dispatcher would turn all fifteen green while
// accepting TRACE ALTER TABLE, TRACE SHOW TABLES and TRACE FORMAT = row,
// none of which v8.5.5 parses, so the rejections are written here by hand:
// every statement outside upstream's TraceableStmt, the bare names FORMAT and
// TARGET take and must not, the two modifiers in combination, and TRACE with
// nothing after it.
//
// Every case below was measured against pingcap/tidb:v8.5.5 through PREPARE
// on 2026-11-04, and the recorded verdict is what that server answered.
// TestTraceFormsMatchTiDB is what keeps it a recording: an image upgrade that
// changes an answer fails there rather than leaving a stale expectation
// behind.
//
// Ref: TraceStmt and TraceableStmt in pingcap/tidb pkg/parser/parser.y.

// traceFormCase is a statement and whether TiDB's parser accepts it.
type traceFormCase struct {
	sql    string
	accept bool
}

var traceFormCases = []traceFormCase{
	// --- the four arms ---
	{"trace select 1", true},
	{"trace format = 'row' select 1", true},
	{"trace format = 'json' select 1", true},
	{"trace plan select 1", true},
	{"trace plan target = 'estimation' select 1", true},
	// The target name is checked when the statement runs, not when it
	// parses, so any text parses.
	{"trace plan target = 'arandomstring' select 1", true},
	{"trace format = 'nosuchformat' select 1", true},
	{"trace format = '' select 1", true},

	// --- the modifiers do not combine, in either order, and TARGET
	// --- belongs to PLAN alone ---
	{"trace plan format = 'row' select 1", false},
	{"trace format = 'row' plan select 1", false},
	{"trace target = 'estimation' select 1", false},
	{"trace format = 'row' target = 'estimation' select 1", false},

	// --- both values are string literals, which is where TRACE parts from
	// --- EXPLAIN: EXPLAIN FORMAT = ROW is EXPLAIN's bare-name arm, and
	// --- TRACE has no such arm ---
	{"trace format = row select 1", false},
	{"trace format = json select 1", false},
	{"trace plan target = estimation select 1", false},
	{"trace format = `row` select 1", false},
	{"trace format = 1 select 1", false},
	// The = is not optional and neither is the value.
	{"trace format row select 1", false},
	{"trace format = select 1", false},
	{"trace format = 'row'", false},
	{"trace plan target 'estimation' select 1", false},
	{"trace plan target = select 1", false},

	// --- a statement is required ---
	{"trace", false},
	{"trace plan", false},
	{"trace trace select 1", false},
	{"trace `select` 1", false},

	// --- TraceableStmt: the queries ---
	{"trace select c1 from t1", true},
	{"trace select 1 union select 2", true},
	{"trace select c1 from t1 union (select c2 from t2) limit 1, 1", true},
	{"trace (select 1)", true},
	{"trace (select 1) union (select 2)", true},
	{"trace with c as (select 1) select * from c", true},
	{"trace with recursive c as (select 1) select * from c", true},
	{"trace table t", true},
	{"trace table t order by a limit 1", true},
	{"trace values row(1)", true},
	{"trace values row(1), row(2)", true},
	{"trace select 1 for update", true},

	// --- TraceableStmt: the DML ---
	{"trace insert into t values (1), (2), (3)", true},
	{"trace insert into t select 1", true},
	{"trace replace into t values (1)", true},
	{"trace update t set id = id + 1 order by id desc", true},
	{"trace update t, t2 set t.a = 1", true},
	{"trace delete from t", true},
	{"trace delete from t where a = 1 limit 1", true},
	{"trace delete t1, t2 from t1 inner join t2 where t1.id=t2.id", true},

	// --- TraceableStmt: the transaction statements. BEGIN OPTIMISTIC and
	// BEGIN PESSIMISTIC are not here: this parser has no BEGIN modifier at
	// all, so `begin optimistic` is a gap row in the corpus on its own and
	// `trace begin optimistic` fails for that reason rather than for
	// anything TRACE does. The gap belongs to the transaction grammar.
	{"trace begin", true},
	{"trace start transaction", true},
	{"trace start transaction read only", true},
	{"trace commit", true},
	{"trace commit and chain", true},
	{"trace rollback", true},
	{"trace rollback to savepoint sp1", true},
	{"trace savepoint sp1", true},
	{"trace release savepoint sp1", true},
	// START reaches BeginTransactionStmt and nothing else START opens.
	{"trace start replica", false},

	// --- TraceableStmt: SET, which is upstream's SetStmt production ---
	{"trace set a = 1", true},
	{"trace set global tidb_mem_quota_query = 1", true},
	{"trace set session autocommit = 1", true},
	{"trace set @v = 1", true},
	{"trace set @@global.autocommit = 1", true},
	{"trace set names utf8mb4", true},
	{"trace set character set utf8mb4", true},
	{"trace set transaction isolation level read committed", true},
	{"trace set password = 'x'", true},
	{"trace set resource group rg1", true},
	// SET ROLE and SET DEFAULT ROLE are productions of their own and are
	// not in SetStmt, so neither is traceable.
	{"trace set role r1", false},
	{"trace set default role r1 to u1", false},

	// --- TraceableStmt: LOAD DATA, and only LOAD DATA ---
	{"trace load data infile '/tmp/a.csv' into table t", true},
	{"trace load data local infile '/tmp/a.csv' into table t", true},
	{"trace load xml infile '/tmp/a.xml' into table t", false},
	{"trace load stats 'x'", false},

	// --- TraceableStmt: ANALYZE TABLE ---
	{"trace analyze table t", true},
	{"trace analyze table t index idx", true},
	{"trace analyze table t all columns", true},
	{"trace analyze incremental table t index idx", true},

	// --- the DDL, which TraceableStmt names none of. This is what parts
	// --- TRACE from EXPLAIN: ALTER TABLE is explainable on TiDB and is not
	// --- traceable, so the catalog never has a traced DDL to apply ---
	{"trace alter table t add column c int", false},
	{"trace alter table t drop column c", false},
	{"trace alter database d charset utf8mb4", false},
	{"trace create table t2 (a int)", false},
	{"trace create index i on t (a)", false},
	{"trace drop table t2", false},
	{"trace truncate table t", false},
	{"trace rename table t to t9", false},

	// --- and the rest of the statements that are not traceable ---
	{"trace explain select 1", false},
	{"trace show tables", false},
	{"trace admin show ddl", false},
	{"trace import into t from '/tmp/a.csv'", false},
	{"trace do 1", false},
	{"trace prepare s from 'select 1'", false},
	{"trace execute s", false},
	{"trace kill 1", false},
	{"trace flush tables", false},
	{"trace use d", false},
	{"trace lock tables t read", false},
	{"trace grant select on *.* to u", false},
}

// TestTraceForms holds this parser to the verdicts above. It needs no
// container.
func TestTraceForms(t *testing.T) {
	for _, c := range traceFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni %ss %q, TiDB %s %ss it (err=%v)",
					verdictName(accepts), c.sql, tidbCorpusVersion, verdictName(c.accept), err)
			}
		})
	}
}

// traceWordsAsNames are statements in which TRACE, FORMAT's and TARGET's
// words are ordinary names instead: a table, a column, an alias, an index.
// TRACE and TARGET are not registered in the lexer at all -- the dispatcher
// reads TRACE with isBareWord -- and these statements are what says so.
// Registering a word breaks every site that reads a name, and the implicit
// alias is the one that broke when the SHOW keywords were registered: its
// branch takes tokIDENT and nothing else, so `DELETE FROM t trace WHERE
// trace.id = 1` would stop parsing the moment TRACE became a keyword.
//
// Every statement here parses on v8.5.5 too, which
// TestTraceWordsStayNamesOnTiDB measures.
var traceWordsAsNames = []string{
	// As a table name, bare and qualified.
	"SELECT * FROM trace",
	"SELECT * FROM d.trace",
	"SELECT * FROM target",
	// As an implicit table alias, the site that broke when the SHOW
	// keywords were registered.
	"DELETE FROM t trace WHERE trace.id = 1",
	"DELETE FROM t target WHERE target.id = 1",
	"UPDATE t trace SET trace.a = 1",
	"SELECT trace.a FROM t trace",
	// As an explicit alias and a column name.
	"SELECT 1 AS trace",
	"SELECT 1 AS target",
	"SELECT trace, target FROM t",
	"SELECT t.trace FROM t",
	// As names in DDL, which is where a registered keyword hurts most.
	"CREATE TABLE trace (target int, trace int)",
	"CREATE TABLE t (trace int, target int, INDEX trace (trace))",
	"CREATE INDEX trace ON t (a)",
	"DROP TABLE trace",
	"ALTER TABLE t ADD COLUMN trace int",
	"ALTER TABLE trace RENAME TO target",
	// And beside the statement of this family, which is what makes the
	// bare-word reading load-bearing rather than incidental.
	"TRACE SELECT * FROM trace",
	"TRACE SELECT trace.a FROM t trace",
	"TRACE DELETE FROM t trace WHERE trace.id = 1",
	"TRACE PLAN SELECT target FROM trace",
}

// TestTraceWordsStayNames holds the other half of the keyword decision: TRACE
// and TARGET are still ordinary names everywhere else. It needs no container.
func TestTraceWordsStayNames(t *testing.T) {
	for _, sql := range traceWordsAsNames {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
}
