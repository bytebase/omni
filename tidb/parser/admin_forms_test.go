package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's ADMIN family, as v8.5.5 reads it.
//
// The corpus in testdata pins the forms pingcap's own parser test writes: 55
// of them, almost all acceptances. A family of thirty small subcommands is the
// easiest place in a grammar to be sloppy -- read ADMIN, read one word, skip
// to the end of the statement -- and a corpus of acceptances cannot see that
// at all. So this file carries the rejections: an unknown subcommand, a tail
// that belongs to a different subcommand, a missing table name, a LIMIT with
// no number, a quoted word where the grammar wants a keyword, and a trailing
// token after the statement should have ended.
//
// Every case below was measured against TiDB v8.5.5 and the expectation is
// what that server answered; TestAdminFormsMatchTiDB holds this file to it, so
// an image upgrade that changes a verdict fails there instead of passing
// silently.
//
// Ref: AdminStmt in TiDB v8.5.5 parser/parser.y:11204-11501.

// adminFormCase is a statement and whether the engine's parser accepts it.
type adminFormCase struct {
	sql    string
	accept bool
}

var adminFormCases = []adminFormCase{
	// ADMIN on its own, and a subcommand TiDB does not have. ADMIN SHOW
	// TELEMETRY and ADMIN CREATE WORKLOAD SNAPSHOT are real TiDB syntax in
	// other releases and not in this one.
	{"ADMIN", false},
	{"ADMIN foo", false},
	{"ADMIN SHOW", false},
	{"ADMIN SHOW TELEMETRY", false},
	{"ADMIN CREATE WORKLOAD SNAPSHOT", false},
	{"admin show ddl", true},
	{"Admin Show Ddl Jobs", true},

	// SHOW DDL and its JOBS tail. The count is unsigned and an int64, and
	// WHERE follows it rather than replacing it.
	{"ADMIN SHOW DDL", true},
	{"ADMIN SHOW DDL JOBS", true},
	{"ADMIN SHOW DDL JOBS 20", true},
	{"ADMIN SHOW DDL JOBS 0", true},
	{"ADMIN SHOW DDL JOBS 9223372036854775807", true},
	{"ADMIN SHOW DDL JOBS 9223372036854775808", false},
	{"ADMIN SHOW DDL JOBS -1", false},
	{"ADMIN SHOW DDL JOBS 1.0", false},
	{"ADMIN SHOW DDL JOBS 1e3", false},
	{"ADMIN SHOW DDL JOBS WHERE id > 0", true},
	{"ADMIN SHOW DDL JOBS 20 WHERE id = 0", true},
	{"ADMIN SHOW DDL JOBS 20 WHERE id = 0 AND state = 'x'", true},
	{"ADMIN SHOW DDL JOBS WHERE id > 0 AND state IN ('a', 'b')", true},
	{"ADMIN SHOW DDL JOBS 20 WHERE id IN (SELECT 1)", true},
	{"ADMIN SHOW DDL JOBS 20 WHERE 1", true},
	{"ADMIN SHOW DDL JOBS WHERE", false},
	{"ADMIN SHOW DDL extra", false},
	{"ADMIN SHOW DDL `jobs`", false},
	// LIMIT and OFFSET are SHOW DDL JOB QUERIES' tail, not SHOW DDL JOBS'.
	{"ADMIN SHOW DDL JOBS 20 LIMIT 5", false},
	{"ADMIN SHOW DDL JOBS 20 OFFSET 2", false},
	{"ADMIN SHOW DDL JOBS 20 WHERE id = 0 LIMIT 1", false},

	// SHOW DDL JOB QUERIES takes an id list or a LIMIT tail, and the two
	// are separate productions, so neither takes the other's tail.
	{"ADMIN SHOW DDL JOB QUERIES 1", true},
	{"ADMIN SHOW DDL JOB QUERIES 1, 2, 3, 4", true},
	{"ADMIN SHOW DDL JOB QUERIES", false},
	{"ADMIN SHOW DDL JOB QUERIES 1,", false},
	{"ADMIN SHOW DDL JOB QUERIES 9223372036854775808", false},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 5", true},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 0", true},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 5, 10", true},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 3 OFFSET 2", true},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 22 OFFSET 0", true},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 18446744073709551615", true},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 18446744073709551616", false},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT", false},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT -1", false},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 5 OFFSET", false},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 5 OFFSET 1 extra", false},
	{"ADMIN SHOW DDL JOB QUERIES 1 LIMIT 2", false},
	{"ADMIN SHOW DDL JOB QUERIES LIMIT 2 WHERE id = 1", false},
	{"ADMIN SHOW DDL JOB 1", false},

	// SHOW ... NEXT_ROW_ID, the one ADMIN SHOW form that takes a table.
	// DDL, SLOW and BDR open their own subcommands, and are table names
	// too, which one token of lookahead settles.
	{"ADMIN SHOW t1 NEXT_ROW_ID", true},
	{"ADMIN SHOW test.t1 NEXT_ROW_ID", true},
	{"ADMIN SHOW ddl NEXT_ROW_ID", true},
	{"ADMIN SHOW slow NEXT_ROW_ID", true},
	{"ADMIN SHOW bdr NEXT_ROW_ID", true},
	{"ADMIN SHOW jobs NEXT_ROW_ID", true},
	{"ADMIN SHOW ddl.t NEXT_ROW_ID", true},
	{"ADMIN SHOW slow.t NEXT_ROW_ID", true},
	{"ADMIN SHOW bdr.t NEXT_ROW_ID", true},
	{"ADMIN SHOW `ddl` NEXT_ROW_ID", true},
	{"ADMIN SHOW `next_row_id` NEXT_ROW_ID", true},
	{"ADMIN SHOW next_row_id NEXT_ROW_ID", true},
	{"ADMIN SHOW t1", false},
	{"ADMIN SHOW NEXT_ROW_ID", false},
	{"ADMIN SHOW t1 NEXT_ROW_ID extra", false},
	{"ADMIN SHOW `ddl`", false},
	{"ADMIN SHOW ddl `next_row_id`", false},
	{"ADMIN SHOW DDL JOBS NEXT_ROW_ID", false},
	{"ADMIN SHOW t.next_row_id", false},
	{"ADMIN SHOW SLOW TOP ALL NEXT_ROW_ID", false},

	// SHOW SLOW. Both forms end in a count, so neither is optional.
	{"ADMIN SHOW SLOW RECENT 11", true},
	{"ADMIN SHOW SLOW TOP 3", true},
	{"ADMIN SHOW SLOW TOP INTERNAL 7", true},
	{"ADMIN SHOW SLOW TOP ALL 9", true},
	{"ADMIN SHOW SLOW", false},
	{"ADMIN SHOW SLOW TOP", false},
	{"ADMIN SHOW SLOW RECENT", false},
	{"ADMIN SHOW SLOW RECENT -1", false},
	{"ADMIN SHOW SLOW TOP foo 3", false},
	{"ADMIN SHOW SLOW TOP 3 extra", false},
	{"ADMIN SHOW `slow` RECENT 11", false},
	{"ADMIN SHOW SLOW TOP ALL 9 WHERE id = 1", false},

	// BDR ROLE: SHOW reads it, SET writes it, UNSET clears it. The role is
	// one of two keywords and takes no quoted name.
	{"ADMIN SHOW BDR ROLE", true},
	{"ADMIN SET BDR ROLE PRIMARY", true},
	{"ADMIN SET BDR ROLE primary", true},
	{"ADMIN SET BDR ROLE SECONDARY", true},
	{"ADMIN UNSET BDR ROLE", true},
	{"ADMIN SHOW BDR ROLE extra", false},
	{"ADMIN SET BDR ROLE", false},
	{"ADMIN SET BDR ROLE foo", false},
	{"ADMIN SET BDR ROLE `primary`", false},
	{"ADMIN SET BDR PRIMARY", false},
	{"ADMIN UNSET BDR", false},

	// CHECK TABLE and CHECKSUM TABLE take a name list, and neither takes
	// the modifiers MySQL's own CHECK TABLE and CHECKSUM TABLE have.
	{"ADMIN CHECK TABLE t1", true},
	{"ADMIN CHECK TABLE t1, t2", true},
	{"ADMIN CHECK TABLE test.t1", true},
	{"ADMIN CHECK TABLE `table`", true},
	{"ADMIN CHECK TABLE", false},
	{"ADMIN CHECK TABLES t1", false},
	{"ADMIN CHECK TABLE t1 QUICK", false},
	{"ADMIN CHECK TABLE t1 FOR UPGRADE", false},
	{"ADMIN CHECKSUM TABLE t1", true},
	{"ADMIN CHECKSUM TABLE t1, t2", true},
	{"ADMIN CHECKSUM TABLE", false},
	{"ADMIN CHECKSUM t1", false},
	{"ADMIN CHECKSUM TABLE t1 QUICK", false},

	// The ADMIN prefix is mandatory: TiDB has no bare CHECK TABLE /
	// CHECKSUM TABLE statement the way MySQL does (parser.y:11235, 11275
	// give only the ADMIN-prefixed productions), so the ADMIN-less forms
	// are rejected.
	{"CHECK TABLE t1", false},
	{"CHECKSUM TABLE t1", false},

	// CHECK, RECOVER and CLEANUP INDEX each take a table and an index
	// name, and only CHECK INDEX takes handle ranges after them.
	{"ADMIN CHECK INDEX t idx", true},
	{"ADMIN CHECK INDEX t.a idx", true},
	{"ADMIN CHECK INDEX `t` `idx`", true},
	{"ADMIN CHECK INDEX t `select`", true},
	{"ADMIN CHECK INDEX next_row_id idx", true},
	{"ADMIN CHECK INDEX t select", false},
	{"ADMIN CHECK INDEX t", false},
	{"ADMIN CHECK INDEX t idx (1, 2)", true},
	{"ADMIN CHECK INDEX t idx (1, 2), (4, 5)", true},
	{"ADMIN CHECK INDEX t idx (1)", false},
	{"ADMIN CHECK INDEX t idx ()", false},
	{"ADMIN CHECK INDEX t idx 1", false},
	{"ADMIN CHECK INDEX t idx (-1, 2)", false},
	{"ADMIN CHECK INDEX t idx (1, 2),", false},
	{"ADMIN RECOVER INDEX t1 idx_a", true},
	{"ADMIN RECOVER INDEX t1", false},
	{"ADMIN RECOVER TABLE t1", false},
	{"ADMIN RECOVER INDEX t1 idx_a (1, 2)", false},
	{"ADMIN CLEANUP INDEX t1 idx_a", true},
	{"ADMIN CLEANUP INDEX t1", false},
	{"ADMIN CLEANUP INDEX t1 idx_a (1, 2)", false},

	// CLEANUP TABLE LOCK, which is CLEANUP's other tail.
	{"ADMIN CLEANUP TABLE LOCK t", true},
	{"ADMIN CLEANUP TABLE LOCK t1,t2", true},
	{"ADMIN CLEANUP TABLE LOCK db.t", true},
	{"ADMIN CLEANUP TABLE LOCK", false},
	{"ADMIN CLEANUP TABLE t", false},

	// CANCEL, PAUSE and RESUME DDL JOBS take an id list of integers.
	{"ADMIN CANCEL DDL JOBS 1", true},
	{"ADMIN CANCEL DDL JOBS 1, 2", true},
	{"ADMIN CANCEL DDL JOBS 0", true},
	{"ADMIN PAUSE DDL JOBS 5", true},
	{"ADMIN PAUSE DDL JOBS 1, 3", true},
	{"ADMIN RESUME DDL JOBS 3", true},
	{"ADMIN RESUME DDL JOBS 1, 2", true},
	{"ADMIN CANCEL DDL JOBS", false},
	{"ADMIN CANCEL DDL JOBS 1,", false},
	{"ADMIN CANCEL DDL JOBS str_not_num", false},
	{"ADMIN CANCEL DDL JOBS 0x1", false},
	{"ADMIN PAUSE DDL JOBS", false},
	{"ADMIN PAUSE DDL JOBS str_not_num", false},
	{"ADMIN RESUME DDL JOBS", false},
	{"ADMIN RESUME DDL JOBS str_not_num", false},
	{"ADMIN CANCEL DDL JOB 1", false},

	// ALTER DDL JOBS takes one id, then at least one `name = value`,
	// comma-separated. The name is any identifier and the value any signed
	// literal, so what TiDB's executor does with an unknown name is a
	// runtime answer (1105), not a parse error.
	{"ADMIN ALTER DDL JOBS 1 THREAD = 2", true},
	{"ADMIN ALTER DDL JOBS 1 thread = 2", true},
	{"ADMIN ALTER DDL JOBS 1 BATCH_SIZE = 3", true},
	{"ADMIN ALTER DDL JOBS 1 MAX_WRITE_SPEED = 4", true},
	{"ADMIN ALTER DDL JOBS 1 MAX_WRITE_SPEED = _UTF8MB4'4MiB'", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = 2, BATCH_SIZE = 3", true},
	{"ADMIN ALTER DDL JOBS 1 `thread` = 2", true},
	{"ADMIN ALTER DDL JOBS 1 foo = 1", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = -2", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = +2", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = 1.5", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = 'x'", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = 'a' 'b'", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = N'1'", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = 0xFF", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = 0b101", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = x'41'", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = b'1'", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = _utf8mb4 x'41'", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = TRUE", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = FALSE", true},
	{"ADMIN ALTER DDL JOBS 1 THREAD = NULL", true},
	{"ADMIN ALTER DDL JOBS 1", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD =", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD = 2 BATCH_SIZE = 3", false},
	{"ADMIN ALTER DDL JOBS 1, 2 THREAD = 2", false},
	{"ADMIN ALTER DDL JOB 1 THREAD = 2", false},
	// A signed literal is the whole grammar here: no expression, no
	// variable, no parentheses, no temporal literal, and the sign reaches
	// a number only.
	{"ADMIN ALTER DDL JOBS 1 THREAD = @a", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD = (1)", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD = 1+1", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD = DATE '2020-01-01'", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD = -'a'", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD = -0xFF", false},
	{"ADMIN ALTER DDL JOBS 1 THREAD = _utf8mb4'x' COLLATE utf8mb4_bin", false},

	// RELOAD's five targets, and nothing else.
	{"ADMIN RELOAD EXPR_PUSHDOWN_BLACKLIST", true},
	{"ADMIN RELOAD OPT_RULE_BLACKLIST", true},
	{"ADMIN RELOAD STATISTICS", true},
	{"ADMIN RELOAD STATS_EXTENDED", true},
	{"ADMIN RELOAD BINDINGS", true},
	{"ADMIN RELOAD", false},
	{"ADMIN RELOAD foo", false},
	{"ADMIN RELOAD EXPR_PUSHDOWN_BLACKLIST extra", false},

	// The binding subcommands. BINDINGS is plural and required; CREATE and
	// DROP BINDING are statements of their own, not ADMIN's.
	{"ADMIN FLUSH BINDINGS", true},
	{"ADMIN CAPTURE BINDINGS", true},
	{"ADMIN EVOLVE BINDINGS", true},
	{"ADMIN FLUSH BINDING", false},
	{"ADMIN CAPTURE BINDING", false},
	{"ADMIN EVOLVE BINDING", false},
	{"ADMIN FLUSH BINDINGS extra", false},

	// FLUSH's other tail. The scope is optional and SESSION when left out.
	{"ADMIN FLUSH PLAN_CACHE", true},
	{"ADMIN FLUSH SESSION PLAN_CACHE", true},
	{"ADMIN FLUSH GLOBAL PLAN_CACHE", true},
	{"ADMIN FLUSH INSTANCE PLAN_CACHE", true},
	{"ADMIN FLUSH LOCAL PLAN_CACHE", false},
	{"ADMIN FLUSH SESSION", false},
	{"ADMIN FLUSH INSTANCE PLAN_CACHE extra", false},

	// PLUGINS takes ENABLE or DISABLE and a comma-separated name list.
	{"ADMIN PLUGINS ENABLE audit", true},
	{"ADMIN PLUGINS ENABLE audit, whitelist", true},
	{"ADMIN PLUGINS DISABLE audit, whitelist", true},
	{"ADMIN PLUGINS ENABLE `audit`", true},
	{"ADMIN PLUGINS", false},
	{"ADMIN PLUGINS ENABLE", false},
	{"ADMIN PLUGINS TOGGLE audit", false},
	{"ADMIN PLUGINS ENABLE select", false},
	{"ADMIN PLUGINS ENABLE audit whitelist", false},
	{"ADMIN PLUGINS ENABLE audit,", false},

	// REPAIR TABLE carries a whole CREATE TABLE, in every form CREATE
	// TABLE takes, and nothing else CREATE opens.
	{"ADMIN REPAIR TABLE t CREATE TABLE t (a int)", true},
	{"ADMIN REPAIR TABLE t CREATE TABLE t (a char(1), b int)", true},
	{"ADMIN REPAIR TABLE t CREATE TABLE IF NOT EXISTS t (a int)", true},
	{"ADMIN REPAIR TABLE t CREATE TEMPORARY TABLE t (a int)", true},
	{"ADMIN REPAIR TABLE t CREATE TABLE t LIKE t2", true},
	{"ADMIN REPAIR TABLE t CREATE TABLE t AS SELECT 1", true},
	{"ADMIN REPAIR TABLE t CREATE TABLE t AS", true},
	{"ADMIN REPAIR TABLE t CREATE TABLE t (a int) AS", true},
	{"ADMIN REPAIR TABLE t CREATE TABLE t (a int) ENGINE=InnoDB", true},
	{"ADMIN REPAIR TABLE t CREATE TABLE t (a int) PARTITION BY HASH(a) PARTITIONS 2", true},
	{"ADMIN REPAIR TABLE t.a CREATE TABLE t (a int)", true},
	{"ADMIN REPAIR TABLE `t` CREATE TABLE t (a int)", true},
	{"ADMIN REPAIR TABLE t", false},
	{"ADMIN REPAIR TABLE t CREATE VIEW v AS SELECT 1", false},
	{"ADMIN REPAIR TABLE t CREATE OR REPLACE VIEW v AS SELECT 1", false},
	{"ADMIN REPAIR TABLE t CREATE INDEX i ON t (a)", false},
	{"ADMIN REPAIR TABLE t CREATE TABLE t (a int) /*T! AUTO_RANDOM */", false},
}

// adminWordsAsNames are statements that use ADMIN's subcommand words where an
// ordinary name belongs. TiDB spells those words as keywords; this parser
// registers none of them and reads each as a word only inside ADMIN, so a
// table, column, alias, or index of any of these names keeps working. A
// registered keyword would break every one of these, which is the failure
// adding SHOW's keywords caused in the family before this one.
var adminWordsAsNames = []string{
	"CREATE TABLE jobs (queries int, bindings int, thread int, bdr int)",
	"CREATE TABLE ddl (next_row_id int, plan_cache int, statistics int)",
	"SELECT next_row_id FROM t",
	"SELECT 1 AS admin",
	"SELECT 1 FROM t AS ddl",
	"SELECT ddl.id FROM t ddl",
	"DELETE FROM t ddl WHERE ddl.id = 1",
	"DELETE FROM t queries WHERE queries.id = 1",
	"UPDATE t jobs SET jobs.a = 1",
	"CREATE INDEX queries ON t (a)",
	"CREATE INDEX next_row_id ON jobs (bdr)",
	"INSERT INTO bindings (capture, evolve) VALUES (1, 2)",
	"SELECT * FROM t WHERE cleanup = 1 AND top = 2 AND recent = 3",
	"SELECT thread, batch_size, max_write_speed FROM jobs",
}

func TestAdminFormsAccepted(t *testing.T) {
	for _, c := range adminFormCases {
		if !c.accept {
			continue
		}
		t.Run(c.sql, func(t *testing.T) {
			if _, err := Parse(c.sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", c.sql, err)
			}
		})
	}
}

func TestAdminFormsRejected(t *testing.T) {
	for _, c := range adminFormCases {
		if c.accept {
			continue
		}
		t.Run(c.sql, func(t *testing.T) {
			if _, err := Parse(c.sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but TiDB %s rejects it", c.sql, tidbRelease)
			}
		})
	}
}

// TestAdminWordsStayNames holds the words ADMIN reads to being names
// everywhere else.
func TestAdminWordsStayNames(t *testing.T) {
	for _, sql := range adminWordsAsNames {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

// parseOneAdmin parses sql and returns the AdminStmt it is.
func parseOneAdmin(t *testing.T, sql string) *nodes.AdminStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*nodes.AdminStmt)
	if !ok {
		t.Fatalf("Parse(%q) returned %T, want *ast.AdminStmt", sql, list.Items[0])
	}
	return stmt
}

// TestAdminFormsAST asserts the tree each form parses into, not only that it
// parses. Thirty subcommands behind one node is exactly where a permissive
// tail hides: a parser that read ADMIN, then a word, then skipped to the end
// of the statement would turn every corpus row green and record nothing a
// caller can use. Each case below names the kind and every operand that kind
// fills.
func TestAdminFormsAST(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{"ADMIN SHOW DDL", `{ADMIN :kind SHOW DDL :loc 0}`},
		{"ADMIN SHOW DDL JOBS", `{ADMIN :kind SHOW DDL JOBS :loc 0}`},
		{"ADMIN SHOW DDL JOBS 20", `{ADMIN :kind SHOW DDL JOBS :loc 0 :jobnumber 20}`},
		{"ADMIN SHOW DDL JOBS 0", `{ADMIN :kind SHOW DDL JOBS :loc 0 :jobnumber 0}`},
		{"ADMIN SHOW DDL JOBS 20 WHERE id = 0", `{ADMIN :kind SHOW DDL JOBS :loc 0 :jobnumber 20 :where {BINEXPR :op = :loc 29 :left {COLREF :loc 29 :col id} :right {INT_LIT :val 0 :loc 34}}}`},
		{"ADMIN SHOW DDL JOBS WHERE id > 0", `{ADMIN :kind SHOW DDL JOBS :loc 0 :where {BINEXPR :op > :loc 26 :left {COLREF :loc 26 :col id} :right {INT_LIT :val 0 :loc 31}}}`},
		{"ADMIN SHOW DDL JOB QUERIES 1", `{ADMIN :kind SHOW DDL JOB QUERIES :loc 0 :jobids 1}`},
		{"ADMIN SHOW DDL JOB QUERIES 1, 2, 3, 4", `{ADMIN :kind SHOW DDL JOB QUERIES :loc 0 :jobids 1 2 3 4}`},
		{"ADMIN SHOW DDL JOB QUERIES LIMIT 5", `{ADMIN :kind SHOW DDL JOB QUERIES LIMIT :loc 0 :limit (count 5 offset 0)}`},
		{"ADMIN SHOW DDL JOB QUERIES LIMIT 5, 10", `{ADMIN :kind SHOW DDL JOB QUERIES LIMIT :loc 0 :limit (count 10 offset 5)}`},
		{"ADMIN SHOW DDL JOB QUERIES LIMIT 3 OFFSET 2", `{ADMIN :kind SHOW DDL JOB QUERIES LIMIT :loc 0 :limit (count 3 offset 2)}`},
		{"ADMIN SHOW t1 NEXT_ROW_ID", `{ADMIN :kind SHOW NEXT_ROW_ID :loc 0 :tables {TABLEREF :loc 11 :name t1}}`},
		{"ADMIN SHOW db.t1 NEXT_ROW_ID", `{ADMIN :kind SHOW NEXT_ROW_ID :loc 0 :tables {TABLEREF :loc 11 :schema db :name t1}}`},
		{"ADMIN SHOW ddl NEXT_ROW_ID", `{ADMIN :kind SHOW NEXT_ROW_ID :loc 0 :tables {TABLEREF :loc 11 :name ddl}}`},
		{"ADMIN SHOW SLOW RECENT 11", `{ADMIN :kind SHOW SLOW :loc 0 :slow (recent 11)}`},
		{"ADMIN SHOW SLOW TOP 3", `{ADMIN :kind SHOW SLOW :loc 0 :slow (top 3)}`},
		{"ADMIN SHOW SLOW TOP INTERNAL 7", `{ADMIN :kind SHOW SLOW :loc 0 :slow (top INTERNAL 7)}`},
		{"ADMIN SHOW SLOW TOP ALL 9", `{ADMIN :kind SHOW SLOW :loc 0 :slow (top ALL 9)}`},
		{"ADMIN SHOW BDR ROLE", `{ADMIN :kind SHOW BDR ROLE :loc 0}`},
		{"ADMIN CHECK TABLE t1, t2", `{ADMIN :kind CHECK TABLE :loc 0 :tables {TABLEREF :loc 18 :name t1} {TABLEREF :loc 22 :name t2}}`},
		{"ADMIN CHECK INDEX t idx", `{ADMIN :kind CHECK INDEX :loc 0 :tables {TABLEREF :loc 18 :name t} :index idx}`},
		{"ADMIN CHECK INDEX t idx (1, 2), (4, 5)", `{ADMIN :kind CHECK INDEX RANGE :loc 0 :tables {TABLEREF :loc 18 :name t} :index idx :range (1 2) :range (4 5)}`},
		{"ADMIN CHECKSUM TABLE t1, t2", `{ADMIN :kind CHECKSUM TABLE :loc 0 :tables {TABLEREF :loc 21 :name t1} {TABLEREF :loc 25 :name t2}}`},
		{"ADMIN RECOVER INDEX t1 idx_a", `{ADMIN :kind RECOVER INDEX :loc 0 :tables {TABLEREF :loc 20 :name t1} :index idx_a}`},
		{"ADMIN CLEANUP INDEX t1 idx_a", `{ADMIN :kind CLEANUP INDEX :loc 0 :tables {TABLEREF :loc 20 :name t1} :index idx_a}`},
		{"ADMIN CLEANUP TABLE LOCK t1, t2", `{ADMIN :kind CLEANUP TABLE LOCK :loc 0 :tables {TABLEREF :loc 25 :name t1} {TABLEREF :loc 29 :name t2}}`},
		{"ADMIN REPAIR TABLE t CREATE TABLE t (a int)", `{ADMIN :kind REPAIR TABLE :loc 0 :tables {TABLEREF :loc 19 :name t} :create {CREATE_TABLE :loc 21 :table {TABLEREF :loc 34 :name t} :columns {COLDEF :loc 37 :name a :type {DATATYPE :loc 39 :name INT}}}}`},
		{"ADMIN CANCEL DDL JOBS 1, 2", `{ADMIN :kind CANCEL DDL JOBS :loc 0 :jobids 1 2}`},
		{"ADMIN PAUSE DDL JOBS 5", `{ADMIN :kind PAUSE DDL JOBS :loc 0 :jobids 5}`},
		{"ADMIN RESUME DDL JOBS 1, 2", `{ADMIN :kind RESUME DDL JOBS :loc 0 :jobids 1 2}`},
		{"ADMIN ALTER DDL JOBS 1 THREAD = 2", `{ADMIN :kind ALTER DDL JOBS :loc 0 :jobnumber 1 :option (thread {INT_LIT :val 2 :loc 32})}`},
		{"ADMIN ALTER DDL JOBS 1 MAX_WRITE_SPEED = _UTF8MB4'4MiB'", `{ADMIN :kind ALTER DDL JOBS :loc 0 :jobnumber 1 :option (max_write_speed {STRING_LIT :val "4MiB" :loc 41 :charset _UTF8MB4})}`},
		{"ADMIN ALTER DDL JOBS 7 THREAD = -2, BATCH_SIZE = 3", `{ADMIN :kind ALTER DDL JOBS :loc 0 :jobnumber 7 :option (thread {UNARY :op - :loc 32 :operand {INT_LIT :val 2 :loc 33}}) :option (batch_size {INT_LIT :val 3 :loc 49})}`},
		{"ADMIN RELOAD EXPR_PUSHDOWN_BLACKLIST", `{ADMIN :kind RELOAD EXPR_PUSHDOWN_BLACKLIST :loc 0}`},
		{"ADMIN RELOAD OPT_RULE_BLACKLIST", `{ADMIN :kind RELOAD OPT_RULE_BLACKLIST :loc 0}`},
		// TiDB gives STATISTICS and STATS_EXTENDED the same kind.
		{"ADMIN RELOAD STATISTICS", `{ADMIN :kind RELOAD STATISTICS :loc 0}`},
		{"ADMIN RELOAD STATS_EXTENDED", `{ADMIN :kind RELOAD STATISTICS :loc 0}`},
		{"ADMIN RELOAD BINDINGS", `{ADMIN :kind RELOAD BINDINGS :loc 0}`},
		{"ADMIN FLUSH BINDINGS", `{ADMIN :kind FLUSH BINDINGS :loc 0}`},
		{"ADMIN CAPTURE BINDINGS", `{ADMIN :kind CAPTURE BINDINGS :loc 0}`},
		{"ADMIN EVOLVE BINDINGS", `{ADMIN :kind EVOLVE BINDINGS :loc 0}`},
		{"ADMIN FLUSH PLAN_CACHE", `{ADMIN :kind FLUSH PLAN_CACHE :loc 0 :scope SESSION}`},
		{"ADMIN FLUSH SESSION PLAN_CACHE", `{ADMIN :kind FLUSH PLAN_CACHE :loc 0 :scope SESSION}`},
		{"ADMIN FLUSH GLOBAL PLAN_CACHE", `{ADMIN :kind FLUSH PLAN_CACHE :loc 0 :scope GLOBAL}`},
		{"ADMIN FLUSH INSTANCE PLAN_CACHE", `{ADMIN :kind FLUSH PLAN_CACHE :loc 0 :scope INSTANCE}`},
		{"ADMIN PLUGINS ENABLE audit, whitelist", `{ADMIN :kind PLUGINS ENABLE :loc 0 :plugins audit, whitelist}`},
		{"ADMIN PLUGINS DISABLE audit", `{ADMIN :kind PLUGINS DISABLE :loc 0 :plugins audit}`},
		{"ADMIN SET BDR ROLE PRIMARY", `{ADMIN :kind SET BDR ROLE :loc 0 :bdrrole PRIMARY}`},
		{"ADMIN SET BDR ROLE SECONDARY", `{ADMIN :kind SET BDR ROLE :loc 0 :bdrrole SECONDARY}`},
		{"ADMIN UNSET BDR ROLE", `{ADMIN :kind UNSET BDR ROLE :loc 0}`},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseOneAdmin(t, tc.sql)
			if got := nodes.NodeToString(stmt); got != tc.want {
				t.Errorf("Parse(%q) =\n  %s\nwant\n  %s", tc.sql, got, tc.want)
			}
		})
	}
}

// TestAdminRepairTableKeepsItsCreate asserts that the CREATE TABLE a REPAIR
// carries reaches the AdminStmt as its own node, where a caller can read it,
// and that the repaired name is kept apart from the name the definition
// gives: ADMIN REPAIR TABLE t CREATE TABLE t2 (...) parses, and the two names
// are different.
func TestAdminRepairTableKeepsItsCreate(t *testing.T) {
	stmt := parseOneAdmin(t, "ADMIN REPAIR TABLE t CREATE TABLE t2 (a int, b int)")
	if stmt.Kind != nodes.AdminRepairTable {
		t.Fatalf("Kind = %s, want %s", stmt.Kind, nodes.AdminRepairTable)
	}
	if len(stmt.Tables) != 1 || stmt.Tables[0].Name != "t" {
		t.Fatalf("Tables = %v, want the one name t", stmt.Tables)
	}
	if stmt.CreateTable == nil {
		t.Fatal("CreateTable is nil, want the CREATE TABLE the statement carries")
	}
	if stmt.CreateTable.Table.Name != "t2" {
		t.Errorf("CreateTable.Table.Name = %q, want %q", stmt.CreateTable.Table.Name, "t2")
	}
	if len(stmt.CreateTable.Columns) != 2 {
		t.Errorf("CreateTable has %d columns, want 2", len(stmt.CreateTable.Columns))
	}
}

// TestAdminStmtEndsAtItsLastToken holds AdminStmt to the Loc contract: the
// span ends at the statement's last token, not at the delimiter or at
// whatever follows it in the script.
func TestAdminStmtEndsAtItsLastToken(t *testing.T) {
	const sql = "ADMIN CHECK TABLE t1;  SELECT 1"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	stmt, ok := list.Items[0].(*nodes.AdminStmt)
	if !ok {
		t.Fatalf("first statement is %T, want *ast.AdminStmt", list.Items[0])
	}
	if got := sql[stmt.Loc.Start:stmt.Loc.End]; got != "ADMIN CHECK TABLE t1" {
		t.Errorf("Loc spans %q, want %q", got, "ADMIN CHECK TABLE t1")
	}
}

// TestAdminIsNotARoutineStatement holds ADMIN out of a routine body. TiDB's
// procedure grammar takes a narrower statement list than its top level, so
// every form here is a 1064 on v8.5.5 while CREATE PROCEDURE p() SELECT 1
// parses and fails at 8108. The corpus writes none of these, and a statement
// family that quietly became legal inside a routine would be an
// over-acceptance nothing else in this package measures.
func TestAdminIsNotARoutineStatement(t *testing.T) {
	for _, sql := range adminRoutineBodyCases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("Parse(%q) was accepted, but TiDB %s rejects it", sql, tidbRelease)
			}
		})
	}
}

// adminRoutineBodyCases are ADMIN statements in the places a routine body
// reaches. TiDB rejects every one.
var adminRoutineBodyCases = []string{
	"CREATE PROCEDURE p() ADMIN SHOW DDL",
	"CREATE PROCEDURE p() ADMIN FLUSH BINDINGS",
	"CREATE PROCEDURE p() BEGIN ADMIN SHOW DDL; END",
	"CREATE PROCEDURE p() BEGIN ADMIN CHECK TABLE t; END",
	"CREATE PROCEDURE p() BEGIN WHILE 1 DO ADMIN SHOW DDL; END WHILE; END",
	"CREATE PROCEDURE p() BEGIN IF 1 THEN ADMIN SHOW DDL; END IF; END",
}
