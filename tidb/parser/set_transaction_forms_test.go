package parser

import (
	"slices"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// SET TRANSACTION as TiDB v8.5.5 reads it, every verdict measured with PREPARE on
// pingcap/tidb:v8.5.5 on 2026-11-05: a 1064 is the parser's refusal, and any other
// answer means the text parsed.
//
// The rules are SetStmt, parser.y:10758-10770, and TransactionChars and TransactionChar,
// parser.y:10840 and 10859. GLOBAL and SESSION are the only scopes it names. The list
// may not be empty, each characteristic has to be whole, and one may repeat. READ ONLY
// alone takes AsOfClause, parser.y:6559, which START TRANSACTION also reads.
//
// TRANSACTION and ISOLATION are unreserved upstream, so they are also variable names,
// and so are the scope words and every other unreserved word: VariableName,
// parser.y:10924, takes any Identifier, parser.y:6774. The token after the word decides,
// since no form of SET continues its keyword with `=`, `:=` or `.` but SET PASSWORD,
// whose `=` makes `SET PASSWORD = 1` a 1064 where `SET PASSWORD := 1` parses.
//
// Not in this table, because omni still disagrees and each is its own change. A scope on
// a later assignment, as in `SET a = 1, GLOBAL b = 2` and `SET a = 1, SESSION transaction
// = 2`: the tree holds one scope per statement. The values ON and BINARY of SetExpr, as
// in `SET a = ON`. `SET NAMES 'utf8'`. A reserved word after a spaced dot, as in `SET a .
// select = 1`. The 32 words TiDB leaves unreserved and omni's table reserves after MySQL,
// which setVariableNameGap lists.

// TestSetTransactionForms holds omni to the recorded verdicts without a server, so a run
// with no Docker still reports a widened or narrowed grammar.
func TestSetTransactionForms(t *testing.T) {
	for _, c := range setTransactionCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

var setTransactionCases = []asOfCase{
	// The clause on SET TRANSACTION.
	{"set transaction read only as of timestamp now()", true},
	{"set transaction read only as of timestamp '2021-04-21 00:42:12'", true},
	{"SET TRANSACTION READ ONLY AS OF TIMESTAMP NOW()", true},
	{"set session transaction read only as of timestamp now()", true},
	{"set global transaction read only as of timestamp now()", true},
	{"set local transaction read only as of timestamp now()", false},
	{"set persist transaction read only as of timestamp now()", false},
	{"set persist_only transaction read only as of timestamp now()", false},
	{"set @@session transaction read only as of timestamp now()", false},
	{"set transaction read only as of timestamp now() - interval 3 second", true},
	{"set transaction read only as of timestamp tidb_bounded_staleness(now(), now())", true},
	{"set transaction read only as of timestamp (select now())", true},
	{"set transaction read only as of timestamp a", true},
	{"set transaction read only as of timestamp @a", true},
	{"set transaction read only as of timestamp 1 or 2", true},
	{"set transaction read only as /* c */ of timestamp now()", true},
	{"set transaction read only as\nof\ntimestamp now()", true},
	{"set transaction read only as of", false},
	{"set transaction read only as of timestamp", false},
	{"set transaction read only as of now()", false},
	{"set transaction read only as of time now()", false},
	{"set transaction read only as now()", false},
	{"set transaction read only as timestamp now()", false},
	{"set transaction read only of timestamp now()", false},
	{"set transaction read only timestamp now()", false},
	{"set transaction read only as `of` timestamp now()", false},
	{"set transaction read only as of `timestamp` now()", false},
	{"set transaction read only as of timestamp now() a", false},
	{"set transaction read only as of timestamp now() as of timestamp now()", false},
	{"set transaction read only as of timestamp now() read write", false},
	{"set transaction read only as of timestamp now() for update", false},
	{"set transaction read write as of timestamp now()", false},
	{"set transaction as of timestamp now()", false},
	{"set transaction isolation level serializable as of timestamp now()", false},
	{"set transaction read only, as of timestamp now()", false},

	// The clause in a list of characteristics.
	{"set transaction read only as of timestamp now(), isolation level serializable", true},
	{"set transaction isolation level serializable, read only as of timestamp now()", true},
	{"set transaction read only as of timestamp now(), read write", true},
	{"set transaction read write, read only as of timestamp now()", true},
	{"set transaction read only as of timestamp 1, read only as of timestamp 2", true},
	{"set transaction read only as of timestamp now(),", false},
	{"set transaction read only as of timestamp 1, 2", false},
	{"set session transaction isolation level read committed, read only as of timestamp now()", true},
	{"set global transaction read only as of timestamp now(), isolation level repeatable read", true},

	// The rest of the characteristics.
	{"set transaction", false},
	{"set session transaction", false},
	{"set global transaction", false},
	{"set transaction read", false},
	{"set transaction read only", true},
	{"set transaction read write", true},
	{"set transaction read committed", false},
	{"set transaction only", false},
	{"set transaction write", false},
	{"set transaction isolation", false},
	{"set transaction isolation level", false},
	{"set transaction level serializable", false},
	{"set transaction isolation serializable", false},
	{"set transaction isolation level serializable", true},
	{"set transaction isolation level repeatable read", true},
	{"set transaction isolation level repeatable", false},
	{"set transaction isolation level read committed", true},
	{"set transaction isolation level read uncommitted", true},
	{"set transaction isolation level read", false},
	{"set transaction isolation level read only", false},
	{"set transaction isolation level committed", false},
	{"set transaction isolation level snapshot", false},
	{"set transaction isolation level `serializable`", false},
	{"set transaction isolation level serializable, read only", true},
	{"set transaction read only, isolation level serializable", true},
	{"set transaction read only, read write", true},
	{"set transaction isolation level serializable, isolation level read committed", true},
	{"set transaction read only,", false},
	{"set transaction, read only", false},
	{"set transaction read only read write", false},
	{"set transaction isolation level serializable read only", false},
	{"set transaction read only foo", false},
	{"set local transaction read only", false},
	{"set persist transaction read only", false},
	{"set persist_only transaction read only", false},
	{"set session transaction isolation level serializable", true},
	{"set global transaction isolation level read committed", true},
	{"set local transaction isolation level serializable", false},
	{"set @@global transaction read only", false},
	{"set @@transaction = 1", true},
	{"set transaction_isolation = 'READ-COMMITTED'", true},

	// TRANSACTION as a variable name, against the statement it can also start.
	{"set persist transaction = 1", false},
	{"set a = 1, transaction read only", false},
	{"set transaction read only, a = 1", false},
	{"set transaction read only, @a = 1", false},
	{"set transaction read only, transaction read write", false},
	{"set transaction isolation level serializable = 1", false},
	{"set transaction read only as of timestamp a = 1", true},

	// TRANSACTION and ISOLATION as variable names: VariableName, parser.y:10924, takes any
	// Identifier, and the word after them is what parts them from the statement.
	{"set transaction = 1", true},
	{"set transaction := 1", true},
	{"set session transaction = 1", true},
	{"set session transaction := 1", true},
	{"set global transaction = 1", true},
	{"set global transaction := 1", true},
	{"set local transaction = 1", true},
	{"set local transaction := 1", true},
	{"set transaction.a = 1", true},
	{"set transaction .a = 1", true},
	{"set transaction. a = 1", true},
	{"set transaction . a = 1", true},
	{"set session transaction.a = 1", true},
	{"set global transaction.a = 1", true},
	{"set local transaction.a = 1", true},
	{"set transaction.transaction = 1", true},
	{"set transaction.isolation = 1", true},
	{"set transaction.select = 1", true},
	{"set transaction.global = 1", true},
	{"set a.transaction = 1", true},
	{"set transaction = 1, a = 2", true},
	{"set a = 1, transaction = 2", true},
	{"set transaction = 1, transaction = 2", true},
	{"set transaction = 1, isolation = 2", true},
	{"set session a = 1, transaction = 2", true},
	{"set global a = 1, transaction = 2", true},
	{"set transaction = 1, a := 2", true},
	{"set session transaction = 1, a = 2", true},
	{"set global transaction = 1, a = 2", true},
	{"set session transaction.a = 1, b = 2", true},
	{"set transaction = (select 1)", true},
	{"set transaction = default", true},
	{"set transaction = a", true},
	{"set transaction = @a", true},
	{"set transaction = 1 + 1", true},
	{"set transaction = 1;", true},
	{"set transaction /* c */ = 1", true},
	{"set /* c */ transaction /* c */ = /* c */ 1", true},
	{"set TRANSACTION = 1", true},
	{"set Transaction.A = 1", true},
	{"set `transaction` = 1", true},
	{"set `transaction`.a = 1", true},
	{"set session `transaction` = 1", true},
	{"set isolation = 1", true},
	{"set isolation := 1", true},
	{"set session isolation = 1", true},
	{"set global isolation = 1", true},
	{"set local isolation = 1", true},
	{"set isolation.a = 1", true},
	{"set a.isolation = 1", true},
	{"set `isolation` = 1", true},
	{"set level = 1", true},
	{"set only = 1", true},
	{"set @transaction = 1", true},
	{"set @@session.transaction = 1", true},
	{"set @@global.transaction = 1", true},
	{"set @@local.transaction = 1", true},
	{"set @isolation = 1", true},
	{"set @@isolation = 1", true},
	{"trace set transaction = 1", true},
	{"trace set transaction read only", true},

	// What stays a 1064 around those names.
	{"set transaction read only, transaction = 1", false},
	{"set transaction isolation level serializable, a = 1", false},
	{"set a = 1, transaction isolation level serializable", false},
	{"set a = 1, session transaction read only", false},
	{"set a = 1, global transaction read only", false},
	{"set transaction = 1, read only", false},
	{"set transaction = 1 read only", false},
	{"set transaction.a.b = 1", false},
	{"set transaction.a", false},
	{"set transaction.", false},
	{"set transaction.a read only", false},
	{"set transaction =", false},
	{"set transaction :=", false},
	{"set transaction == 1", false},
	{"set transaction != 1", false},
	{"set transaction < 1", false},
	{"set transaction 1", false},
	{"set transaction (1)", false},
	{"set transaction, a = 1", false},
	{"set transaction;", false},
	{"set `transaction` read only", false},
	{"set isolation level = 1", false},
	{"set isolation level serializable", false},
	{"set read = 1", false},
	{"set read only = 1", false},
	{"set write = 1", false},
	{"set session session transaction = 1", false},
	{"set global session transaction = 1", false},
	{"set session global transaction read only", false},

	// PERSIST and PERSIST_ONLY are not scopes on TiDB, only names. VariableAssignment,
	// parser.y:10942, has GLOBAL, SESSION and LOCAL.
	{"set persist transaction := 1", false},
	{"set persist_only transaction = 1", false},
	{"set persist transaction.a = 1", false},
	{"set persist isolation = 1", false},
	{"set persist a = 1", false},
	{"set persist_only a = 1", false},
	{"set a = 1, persist b = 2", false},
	{"set persist = 1", true},
	{"set persist := 1", true},
	{"set persist_only = 1", true},
	{"set persist_only := 1", true},
	{"set persist.a = 1", true},
	{"set persist_only.a = 1", true},
	{"set session persist = 1", true},
	{"set global persist = 1", true},
	{"set a = 1, persist = 2", true},
	{"set a = 1, persist_only = 2", true},
	{"set @@persist.a = 1", true},
	{"set @@persist_only.a = 1", true},

	// The scope words, read as names when `=`, `:=` or `.` follows.
	{"set global = 1", true},
	{"set session = 1", true},
	{"set local = 1", true},
	{"set global := 1", true},
	{"set session := 1", true},
	{"set global.a = 1", true},
	{"set session.a = 1", true},
	{"set local.a = 1", true},
	{"set global . a = 1", true},
	{"set session.transaction = 1", true},
	{"set global.transaction = 1", true},
	{"set session global = 1", true},
	{"set session session = 1", true},
	{"set session local = 1", true},
	{"set local session = 1", true},
	{"set global = 1, session = 2", true},
	{"set session = (select 1)", true},
	{"set a.global = 1", true},
	{"set a.session = 1", true},
	{"set a.local = 1", true},
	{"set a . global = 1", true},
	{"set global.a.b = 1", false},

	// The other words that start a form of SET.
	{"set names = 1", true},
	{"set names := 1", true},
	{"set names.a = 1", true},
	{"set names utf8", true},
	{"set names default", true},
	{"set charset = 1", true},
	{"set charset := 1", true},
	{"set charset.a = 1", true},
	{"set charset utf8", true},
	{"set charset default", true},
	{"set character = 1", false},
	{"set character set utf8", true},
	{"set role = 1", true},
	{"set role := 1", true},
	{"set role.a = 1", true},
	{"set role all", true},
	{"set role none", true},
	{"set resource = 1", true},
	{"set resource := 1", true},
	{"set resource.a = 1", true},
	{"set resource group rg", true},
	{"set binding = 1", true},
	{"set config = 1", true},
	{"set default = 1", false},
	{"set default.a = 1", true},
	{"set default role none to u", true},
	{"set password = 1", false},
	{"set password := 1", true},
	{"set password.a = 1", true},
	{"set password = 'x'", true},
	{"set password = password('x')", true},
	{"set password for u = 'x'", true},
	{"set session password = 1", true},

	// Unreserved words that had no category, which MySQL's lvalue_ident test refused.
	{"set timestamp = 1", true},
	{"set time = 1", true},
	{"set date = 1", true},
	{"set user = 1", true},
	{"set status = 1", true},
	{"set session timestamp = 1", true},
	{"set a = 1, timestamp = 2", true},
	{"set timestamp.a = 1", true},
	{"set a.timestamp = 1", true},
	{"set a.select = 1", true},
	{"set select.a = 1", true},
	{"set a . transaction = 1", true},

	// Nothing changes for the words outside SET: every other name position already read them,
	// and a role name takes no keyword at all, RoleNameString reading a bare identifier.
	{"select transaction from t", true},
	{"select isolation from t", true},
	{"select transaction.a from t", true},
	{"select a.transaction from t", true},
	{"select a transaction from t", true},
	{"select a as transaction from t", true},
	{"select * from transaction", true},
	{"select * from t transaction", true},
	{"select * from t as isolation", true},
	{"create table transaction (transaction int, isolation int)", true},
	{"update t set transaction = 1", true},
	{"insert into t (transaction) values (1)", true},
	{"create role transaction", false},
	{"drop role transaction", false},
	{"grant transaction to u", false},
	{"grant isolation to u", false},
	{"revoke transaction from u", false},
	{"set role transaction", false},
	{"set role isolation", false},
	{"set default role transaction to u", false},
	{"grant select on t to u as a with role transaction", false},
	{"grant select on t to u as a with role isolation", false},
	{"reset persist transaction", false},
	{"reset persist isolation", false},
	{"fetch c into transaction", false},

	// ARRAY and UNTIL are reserved on TiDB, parser.y:85 and 289.
	{"set array = 1", false},
	{"set until = 1", false},
	{"set session array = 1", false},
	{"set array.a = 1", true},
	{"set until.a = 1", true},
	{"set a.array = 1", true},
	{"set a.until = 1", true},
	{"select array from t", false},
	{"select until from t", false},
	{"select array.a from t", true},
	{"select until.a from t", true},
	{"select a.array from t", true},
	{"select a.until from t", true},
	{"create table array (a int)", false},
	{"create table until (a int)", false},
	{"select cast(a as unsigned array) from t", true},
	{"create table t (a json, index i ((cast(a as unsigned array))))", true},

	// START TRANSACTION, which shares the rule.
	{"start transaction read only as of timestamp now()", true},
	{"start transaction read only as of", false},
	{"start transaction read only as of timestamp", false},
	{"start transaction read only as now()", false},
	{"start transaction read only as timestamp now()", false},
	{"start transaction read only as `of` timestamp now()", false},
	{"start transaction read only of timestamp now()", false},
	{"start transaction read only timestamp now()", false},
	{"start transaction read only as of timestamp now(), read write", false},
}

// setVariableNameGap holds the words TiDB v8.5.5 takes as a variable name and omni
// refuses, because omni's keyword table reserves them after MySQL. Every name position
// refuses them for that reason, not only SET, so the fix belongs to the table. The oracle
// sweep fails when one of them starts to agree, so this list cannot outlive the gap.
var setVariableNameGap = []string{
	"accessible", "asensitive", "before", "condition", "cube", "declare", "deterministic",
	"each", "empty", "function", "get", "grouping", "insensitive", "io_after_gtids",
	"io_before_gtids", "json_table", "lateral", "loop", "modifies", "optimizer_costs",
	"purge", "read_write", "reads", "resignal", "return", "sensitive", "separator",
	"signal", "specific", "system", "undo", "variance",
}

// setVariableNamePositions are the positions the oracle sweep tries every keyword in.
var setVariableNamePositions = []string{
	"set %s = 1",
	"set %s := 1",
	"set session %s = 1",
	"set global %s = 1",
	"set local %s = 1",
	"set %s.a = 1",
	"set a.%s = 1",
	"set a = 1, %s = 2",
	"set %s = 1, a = 2",
}

// TestSetVariableNameTrees pins what the word turned into, because a verdict cannot tell
// a variable from the statement or the scope the same word can start.
func TestSetVariableNameTrees(t *testing.T) {
	for _, c := range []struct {
		sql, scope, table, column string
	}{
		{"set transaction = 1", "", "", "transaction"},
		{"set transaction := 1", "", "", "transaction"},
		{"set session transaction = 1", "SESSION", "", "transaction"},
		{"set global transaction = 1", "GLOBAL", "", "transaction"},
		{"set local transaction = 1", "LOCAL", "", "transaction"},
		{"set transaction.a = 1", "", "transaction", "a"},
		{"set session transaction.a = 1", "SESSION", "transaction", "a"},
		{"set a.transaction = 1", "", "a", "transaction"},
		{"set isolation = 1", "", "", "isolation"},
		{"set global = 1", "", "", "global"},
		{"set session = 1", "", "", "session"},
		{"set session.a = 1", "", "session", "a"},
		{"set session global = 1", "SESSION", "", "global"},
		{"set persist = 1", "", "", "persist"},
		{"set names = 1", "", "", "names"},
		{"set charset = 1", "", "", "charset"},
		{"set role = 1", "", "", "role"},
		{"set resource = 1", "", "", "resource"},
		{"set password := 1", "", "", "password"},
		{"set timestamp = 1", "", "", "timestamp"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatal(err)
			}
			stmt, ok := list.Items[0].(*nodes.SetStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.SetStmt", list.Items[0])
			}
			if len(stmt.Assignments) != 1 {
				t.Fatalf("got %d assignments, want 1", len(stmt.Assignments))
			}
			col := stmt.Assignments[0].Column
			if stmt.Scope != c.scope || col.Table != c.table || col.Column != c.column {
				t.Errorf("scope=%q table=%q column=%q, want scope=%q table=%q column=%q",
					stmt.Scope, col.Table, col.Column, c.scope, c.table, c.column)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(c.sql) {
				t.Errorf("statement Loc = %+v, want {0 %d}", stmt.Loc, len(c.sql))
			}
		})
	}

	// The list keeps both names, and the statements the same words start stay themselves.
	list, err := Parse("set a = 1, transaction = 2")
	if err != nil {
		t.Fatal(err)
	}
	got := []string{}
	for _, a := range list.Items[0].(*nodes.SetStmt).Assignments {
		got = append(got, a.Column.Column)
	}
	if !slices.Equal(got, []string{"a", "transaction"}) {
		t.Errorf("assignments = %v", got)
	}
	for _, sql := range []string{"set transaction read only", "set session transaction read only", "set global transaction isolation level serializable"} {
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := list.Items[0].(*nodes.SetTransactionStmt); !ok {
			t.Errorf("%q: got %T, want *ast.SetTransactionStmt", sql, list.Items[0])
		}
	}
}
