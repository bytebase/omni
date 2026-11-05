package parser

import (
	"slices"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// SET as MariaDB 11.8.8 reads it, every verdict measured with PREPARE on mariadb:11.8.8 on
// 2026-11-05. A 1064 is the parser's refusal, and any other answer means the text parsed. That
// includes 1193 and 4081, an unknown variable, which the server raises while parsing: omni holds
// no list of system variables, so here they mean the word was read as a name.
//
// Measure a later word against a real variable. `SET a = 1, x = 2` answers 1193 for a before it
// reaches x, whatever x is, which is why these rows say autocommit.
//
// Not in this table, because omni still disagrees and the change is not about them. A bare DEFAULT
// inside an expression, as in `SET a = 1 + DEFAULT`, `SET a = (DEFAULT)` and `SELECT DEFAULT`, and
// `SET @@ a = 1`, because expressions read them everywhere. `SET STATEMENT a = 1 FOR ...` and `SET
// PASSWORD = OLD_PASSWORD('x')`, which omni does not read. `a VECTOR(3)` and `CREATE TABLE t (a.b
// INT)`, which it does not read either.
// `RESET PERSIST`, a statement MariaDB does not have. The words of setNameGap, whose fix belongs to
// keywordCategories, among them SYSTEM as in `SET a = SYSTEM`.

// setNameGap holds the words MariaDB takes as a name and omni's keyword table reserves after
// MySQL. Every name position refuses them for that reason, not only SET. The oracle sweep fails
// when one of them starts to agree, so this list cannot outlive the gap.
var setNameGap = []string{
	"cube", "cume_dist", "database", "dense_rank", "empty", "first_value", "function", "generated",
	"get", "grouping", "groups", "io_after_gtids", "io_before_gtids", "json_table", "lag",
	"last_value", "lateral", "lead", "master_bind", "master_ssl_verify_server_cert", "nth_value",
	"ntile", "of", "optimizer_costs", "option", "percent_rank", "rank", "schema", "stored", "system",
	"variance", "virtual", "window",
}

// setNamePositions are the positions the oracle sweep tries every keyword in. autocommit is real,
// so the word after it is what the server answers about.
var setNamePositions = []string{
	"set %s = 1",
	"set %s := 1",
	"set session %s = 1",
	"set global %s = 1",
	"set local %s = 1",
	"set %s.a = 1",
	"set a.%s = 1",
	"set autocommit = 1, %s = 2",
	"set %s = 1, autocommit = 2",
}

type setNameCase struct {
	sql    string
	accept bool
}

func setNameParses(sql string) (bool, error) {
	list, err := Parse(sql)
	if err != nil {
		return false, err
	}
	return len(list.Items) == 1, nil
}

// TestSetNameForms holds omni to the recorded verdicts without a server, so a run with no Docker
// still reports a widened or narrowed grammar.
func TestSetNameForms(t *testing.T) {
	for _, c := range setNameCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := setNameParses(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestSetAssignmentTrees pins the assignments a list becomes. Scope carries to the unprefixed
// names after it, past @x, @@x and NAMES, and is empty on those: measured on 11.8 by running each
// list and reading @@global and @@session back.
func TestSetAssignmentTrees(t *testing.T) {
	type want struct{ scope, column string }
	for _, c := range []struct {
		sql       string
		stmtScope string
		want      []want
	}{
		{"set a = 1", "", []want{{"", "a"}}},
		{"set global a = 1, b = 2", "GLOBAL", []want{{"GLOBAL", "a"}, {"GLOBAL", "b"}}},
		{"set session a = 1, b = 2", "SESSION", []want{{"SESSION", "a"}, {"SESSION", "b"}}},
		{"set a = 1, global b = 2", "", []want{{"", "a"}, {"GLOBAL", "b"}}},
		{"set global a = 1, session b = 2, c = 3", "GLOBAL", []want{{"GLOBAL", "a"}, {"SESSION", "b"}, {"SESSION", "c"}}},
		{"set session a = 1, global b = 2, c = 3", "SESSION", []want{{"SESSION", "a"}, {"GLOBAL", "b"}, {"GLOBAL", "c"}}},
		{"set global a = 1, @@session.b = 2, c = 3", "GLOBAL", []want{{"GLOBAL", "a"}, {"", "@@SESSION.b"}, {"GLOBAL", "c"}}},
		{"set global a = 1, @x = 2, b = 3", "GLOBAL", []want{{"GLOBAL", "a"}, {"", "@x"}, {"GLOBAL", "b"}}},
		{"set global a = 1, names utf8, b = 3", "GLOBAL", []want{{"GLOBAL", "a"}, {"", "NAMES"}, {"GLOBAL", "b"}}},
		{"set transaction = 1", "", []want{{"", "transaction"}}},
		{"set session transaction = 1", "SESSION", []want{{"SESSION", "transaction"}}},
		{"set global = 1", "", []want{{"", "global"}}},
		{"set session global = 1", "SESSION", []want{{"SESSION", "global"}}},
		{"set persist = 1", "", []want{{"", "persist"}}},
		{"set timestamp = 1", "", []want{{"", "timestamp"}}},
		{"set offset = 1", "", []want{{"", "offset"}}},
		{"set new = 1", "", []want{{"", "new"}}},
		{"set charset = 1", "", []want{{"", "charset"}}},
		{"set role = 1", "", []want{{"", "role"}}},
		{"set session names = 1", "SESSION", []want{{"SESSION", "names"}}},
		{"set names 'utf8'", "", []want{{"", "NAMES"}}},
		{"set names utf8 collate default", "", []want{{"", "NAMES"}, {"", "COLLATE"}}},
		{"set names default collate binary", "", []want{{"", "NAMES"}, {"", "COLLATE"}}},
		{"set char set utf8", "", []want{{"", "CHARACTER SET"}}},
		{"set a = 1, charset latin1", "", []want{{"", "a"}, {"", "CHARACTER SET"}}},
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
			if stmt.Scope != c.stmtScope {
				t.Errorf("statement scope = %q, want %q", stmt.Scope, c.stmtScope)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(c.sql) {
				t.Errorf("statement Loc = %+v, want {0 %d}", stmt.Loc, len(c.sql))
			}
			got := []want{}
			for _, a := range stmt.Assignments {
				got = append(got, want{a.Scope, a.Column.Column})
			}
			if !slices.Equal(got, c.want) {
				t.Errorf("assignments = %v, want %v", got, c.want)
			}
		})
	}

	first := func(t *testing.T, sql string) *nodes.Assignment {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		return list.Items[0].(*nodes.SetStmt).Assignments[0]
	}

	// An assignment's Loc takes in its scope word, and the printed tree holds the scope.
	const sql = "set a = 1, global b = 2, c = on"
	list, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	stmt := list.Items[0].(*nodes.SetStmt)
	got := []string{}
	for _, a := range stmt.Assignments {
		got = append(got, sql[a.Loc.Start:a.Loc.End])
	}
	if want := []string{"a = 1", "global b = 2", "c = on"}; !slices.Equal(got, want) {
		t.Errorf("assignment text = %q, want %q", got, want)
	}
	if printed := nodes.NodeToString(stmt); !strings.Contains(printed, ":loc 11 :scope GLOBAL :col") {
		t.Errorf("printed tree has no scope on the second assignment: %s", printed)
	}

	// The values. ON stays the true the review walk reads, which `SET a = 'ON'` is not.
	if v, ok := first(t, "set a = on").Value.(*nodes.BoolLit); !ok || !v.Value {
		t.Errorf("ON is %#v", first(t, "set a = on").Value)
	}
	if v, ok := first(t, "set a = off").Value.(*nodes.BoolLit); !ok || v.Value {
		t.Errorf("OFF is %#v", first(t, "set a = off").Value)
	}
	for sql, want := range map[string]string{"set a = binary": "BINARY", "set a = all": "ALL", "set @@a = binary": "BINARY"} {
		if v, ok := first(t, sql).Value.(*nodes.StringLit); !ok || v.Value != want {
			t.Errorf("%q: value is %#v", sql, first(t, sql).Value)
		}
	}
	if _, ok := first(t, "set a = default").Value.(*nodes.DefaultExpr); !ok {
		t.Errorf("DEFAULT is %T", first(t, "set a = default").Value)
	}
	if d, ok := first(t, "set a = default(b)").Value.(*nodes.DefaultExpr); !ok || d.Column != "b" {
		t.Errorf("DEFAULT(b) is %#v", first(t, "set a = default(b)").Value)
	}
	if _, ok := first(t, "set a = binary 'x'").Value.(*nodes.StringLit); ok {
		t.Errorf("BINARY 'x' read as the word BINARY")
	}

	// A trigger's qualifier stays distinct from a table's, and a word before a dot is a table's.
	if col := first(t, "set new.a = 1").Column; col.Table != "NEW" || col.Column != "a" {
		t.Errorf("NEW.a is table=%q column=%q", col.Table, col.Column)
	}
	if col := first(t, "set global.a = 1").Column; col.Table != "global" || col.Column != "a" {
		t.Errorf("global.a is table=%q column=%q", col.Table, col.Column)
	}
}

// TestWordBeforeDot pins what the lexer's rule turns the word into, against the readers of
// @@scope and of a trigger's qualifier, which need the keyword to stay.
func TestWordBeforeDot(t *testing.T) {
	list, err := Parse("select select.a, from.b.c from t")
	if err != nil {
		t.Fatal(err)
	}
	item := list.Items[0].(*nodes.SelectStmt).TargetList[0]
	if res, ok := item.(*nodes.ResTarget); ok {
		item = res.Val
	}
	if col, ok := item.(*nodes.ColumnRef); !ok || col.Table != "select" || col.Column != "a" {
		t.Errorf("select.a read as %#v", item)
	}
	list, err = Parse("select @@global.a")
	if err != nil {
		t.Fatal(err)
	}
	item = list.Items[0].(*nodes.SelectStmt).TargetList[0]
	if res, ok := item.(*nodes.ResTarget); ok {
		item = res.Val
	}
	if ref, ok := item.(*nodes.VariableRef); !ok || ref.Scope != "GLOBAL" || ref.Name != "a" {
		t.Errorf("@@global.a read as %#v", item)
	}
	for sql, wantKeyword := range map[string]bool{"select.a": false, "select .a": true, "select. a": true, "select.": true, "select.*": true} {
		if got := NewLexer(sql).NextToken().Type != tokIDENT; got != wantKeyword {
			t.Errorf("%q: keyword=%v, want %v", sql, got, wantKeyword)
		}
	}
}

var setNameCases = []setNameCase{
	// Every assignment takes a scope, and a written scope takes a name and nothing else.
	{"set", false},
	{"set , autocommit = 1", false},
	{"set @@autocommit = 1", true},
	{"set @a = 1", true},
	{"set @a = 1, autocommit = 1", true},
	{"set `global` = 1", true},
	{"set `select` = 1", true},
	{"set autocommit", false},
	{"set autocommit := 1", true},
	{"set autocommit =", false},
	{"set autocommit = 1", true},
	{"set autocommit = 1,", false},
	{"set autocommit = 1, @@autocommit = 1", true},
	{"set autocommit = 1, @a = 1", true},
	{"set autocommit = 1, autocommit = 0", true},
	{"set autocommit = 1, global autocommit = 1", true},
	{"set autocommit = 1, local autocommit = 1", true},
	{"set autocommit = 1, session autocommit = 1", true},
	{"set autocommit = 1, timestamp = 1", true},
	{"set autocommit = off", true},
	{"set binding = 1", true},
	{"set date = 1", true},
	{"set from = 1", false},
	{"set global := 1", true},
	{"set global = 1", true},
	{"set global @@autocommit = 1", false},
	{"set global @a = 1", false},
	{"set global autocommit = 1", true},
	{"set global autocommit = 1, @a = 1, autocommit = 1", true},
	{"set global autocommit = 1, autocommit = 0", true},
	{"set global autocommit = 1, session autocommit = 0", true},
	{"set global autocommit = 1, session autocommit = 0, autocommit = 1", true},
	{"set global global = 1", true},
	{"set global timestamp = 1", true},
	{"set local = 1", true},
	{"set local autocommit = 1", true},
	{"set resource = 1", true},
	{"set resource group rg", false},
	{"set select = 1", false},
	{"set session = 1", true},
	{"set session @a = 1", false},
	{"set session autocommit = 1", true},
	{"set session global = 1", true},
	{"set session session = 1", true},
	{"set session timestamp = 1", true},
	{"set status = 1", true},
	{"set time = 1", true},
	{"set timestamp = 1", true},
	{"set user = 1", true},

	// TRANSACTION and the words around it, as names against the statement they can start.
	{"set `transaction` = 1", true},
	{"set autocommit = 1, transaction = 1", true},
	{"set global transaction = 1", true},
	{"set global transaction read only", true},
	{"set isolation = 1", true},
	{"set level = 1", true},
	{"set local transaction = 1", true},
	{"set local transaction read only", true},
	{"set only = 1", true},
	{"set read = 1", false},
	{"set session transaction = 1", true},
	{"set session transaction read only", true},
	{"set transaction := 1", true},
	{"set transaction = 1", true},
	{"set transaction isolation level serializable", true},
	{"set transaction read only", true},
	{"set transaction read write", true},
	{"set write = 1", false},

	// There is no PERSIST: the word is only a name.
	{"set autocommit = 1, persist = 1", true},
	{"set persist = 1", true},
	{"set persist a = 1", false},
	{"set persist autocommit = 1", false},
	{"set persist_only = 1", true},
	{"set persist_only autocommit = 1", false},

	// set_expr_or_default for a name or @@name. @name takes an expression.
	{"set @@autocommit = all", true},
	{"set @@autocommit = binary", true},
	{"set @@autocommit = default", true},
	{"set @@autocommit = on", true},
	{"set @@autocommit = row", true},
	{"set @a = all", false},
	{"set @a = binary", false},
	{"set @a = binary 'x'", true},
	{"set @a = default", false},
	{"set @a = default(b)", true},
	{"set @a = on", false},
	{"set @a = row", true},
	{"set @a = row(1, 2)", true},
	{"set autocommit = 'on'", true},
	{"set autocommit = (on)", false},
	{"set autocommit = 1 + 1", true},
	{"set autocommit = all", true},
	{"set autocommit = all + 1", false},
	{"set autocommit = all, autocommit = 1", true},
	{"set autocommit = binary", true},
	{"set autocommit = binary 'x'", true},
	{"set autocommit = binary 1 + 1", true},
	{"set autocommit = binary, autocommit = 1", true},
	{"set autocommit = default", true},
	{"set autocommit = default + 1", false},
	{"set autocommit = default(a)", true},
	{"set autocommit = default, autocommit = 1", true},
	{"set autocommit = on", true},
	{"set autocommit = on + 1", false},
	{"set autocommit = on, autocommit = 1", true},
	{"set autocommit = row", true},
	{"set autocommit = row + 1", true},
	{"set autocommit = row(1, 2)", true},
	{"set autocommit = row, autocommit = 1", true},

	{"set autocommit = default(b) + 1", true},
	{"set @a = default(b) + 1", true},

	// NAMES and the charset forms, anywhere in the list. NAMES has no reading as a name without
	// a scope.
	{"set autocommit = 1, character set utf8", true},
	{"set autocommit = 1, charset utf8", true},
	{"set autocommit = 1, names = 2", false},
	{"set autocommit = 1, names utf8", true},
	{"set autocommit = 1, session names = 2", true},
	{"set char = 1", false},
	{"set char set 'utf8'", true},
	{"set char set default", true},
	{"set char set utf8", true},
	{"set character = 1", false},
	{"set character set 'utf8'", true},
	{"set character set binary", true},
	{"set character set default", true},
	{"set character set utf8", true},
	{"set character set utf8, autocommit = 1", true},
	{"set charset 'utf8'", true},
	{"set charset = 1", true},
	{"set charset binary", true},
	{"set charset default", true},
	{"set charset utf8", true},
	{"set charset utf8 collate utf8_bin", false},
	{"set charset utf8, autocommit = 1", true},
	{"set charset.a = 1", true},
	{"set global autocommit = 1, names utf8", true},
	{"set global charset = 1", true},
	{"set global names utf8", false},
	{"set names 'utf8'", true},
	{"set names := 1", false},
	{"set names = 1", false},
	{"set names @a", false},
	{"set names binary", true},
	{"set names binary collate binary", true},
	{"set names binary collate default", true},
	{"set names default", true},
	{"set names default collate utf8_bin", true},
	{"set names utf8", true},
	{"set names utf8 collate 'utf8_bin'", true},
	{"set names utf8 collate binary", true},
	{"set names utf8 collate default", true},
	{"set names utf8 collate utf8_bin", true},
	{"set names utf8, autocommit = 1", true},
	{"set names.a = 1", true},
	{"set session names = 1", true},

	// The other forms, which are whole statements.
	{"set autocommit = 1, password := 2", false},
	{"set autocommit = 1, password = 1", false},
	{"set autocommit = 1, session password = 1", true},
	{"set default = 1", false},
	{"set default role a", true},
	{"set default role a for b", true},
	{"set default role none", true},
	{"set new = 1", true},
	{"set old = 1", true},
	{"set password := 1", false},
	{"set password = 'x'", true},
	{"set password = 1", false},
	{"set password = password('x')", true},
	{"set password for u = 'x'", true},
	{"set role := 1", true},
	{"set role = 1", true},
	{"set role a", true},
	{"set role none", true},
	{"set session password = 1", true},
	{"set statement = 1", true},

	// A word touching a dot and a name is a name, whatever the word.
	{"set @@a . b = 1", true},
	{"set @@global.autocommit = 1", true},
	{"set @@local.autocommit = 1", true},
	{"set @@persist.autocommit = 1", true},
	{"set @@persist_only.autocommit = 1", true},
	{"set @@select.a = 1", false},
	{"set @@session.autocommit = 1", true},
	{"set a . select = 1", false},
	{"set a.@b = 1", false},
	{"set a.new = 1", true},
	{"set a.timestamp = 1", true},
	{"set autocommit . select = 1", false},
	{"set autocommit = 1, global.a = 1", true},
	{"set autocommit = 1, password.a = 2", true},
	{"set default.a = 1", true},
	{"set from.a = 1", true},
	{"set global.a = 1", true},
	{"set global.autocommit = 1", true},
	{"set local.a = 1", true},
	{"set new.a = 1", true},
	{"set new.a.b = 1", false},
	{"set old.a = 1", true},
	{"set password.a = 1", true},
	{"set persist.a = 1", true},
	{"set role.a = 1", true},
	{"set select .a = 1", false},
	{"set select. a = 1", false},
	{"set select.a = 1", true},
	{"set session.a = 1", true},
	{"set timestamp.a = 1", true},
	{"set transaction.a = 1", true},

	// The same rule outside SET, and what stays of ON and the words around it.
	{"alter table t lock = none", true},
	{"alter table t page_compressed = on", false},
	{"alter table t stats_persistent = 0", true},
	{"change master to master_host = 'a'", true},
	{"create table offset (a int)", false},
	{"create table select.a (x int)", true},
	{"create table select.from (x int)", true},
	{"create table t (a int) page_compressed = on", false},
	{"create table t (a int) sequence = on", false},
	{"create table t (a int) stats_auto_recalc = 1", true},
	{"create table t (a int) stats_persistent = 1", true},
	{"create table t (a int) stats_persistent = on", false},
	{"create table t (a int) stats_sample_pages = 10", true},
	{"create table t (a int) with system versioning", true},
	{"create table t (vector int)", false},
	{"create table vector (a int)", false},
	{"delete from select.a", true},
	{"drop table select.a", true},
	{"fetch c into timestamp", false},
	{"insert into select.a values (1)", true},
	{"reset persist timestamp", false},
	{"select * from from.select", true},
	{"select * from select.a", true},
	{"select * from t for system_time all", true},
	{"select * from t limit 1 offset 2", true},
	{"select * from t limit 1, 2", true},
	{"select * from t offset", false},
	{"select 1 from select.a", true},
	{"select @@global.a", true},
	{"select @@global.autocommit", true},
	{"select @@local.autocommit", true},
	{"select @@select", false},
	{"select @@select.a", false},
	{"select @@session.autocommit", true},
	{"select `select`.a from t", true},
	{"select a offset from t", false},
	{"select a.select from t", true},
	{"select all", false},
	{"select binary", false},
	{"select case.a from t", true},
	{"select count.a from t", true},
	{"select date.a from t", true},
	{"select from.a from t", true},
	{"select global.a", true},
	{"select if.a from t", true},
	{"select interval.a from t", true},
	{"select new.a", true},
	{"select not.a from t", true},
	{"select now.a from t", true},
	{"select null.a from t", true},
	{"select off", true},
	{"select off from t", true},
	{"select offset from t", false},
	{"select old.a", true},
	{"select on", false},
	{"select row", true},
	{"select row from t", true},
	{"select select .a from t", false},
	{"select select.", false},
	{"select select. a from t", false},
	{"select select.`a` from t", false},
	{"select select.a from t", true},
	{"select select.a.b from t", true},
	{"select session.a", true},
	{"select table.* from table.a", false},
	{"select true.a from t", true},
	{"show columns from select.a", true},
	{"show tables from select", false},
	{"start slave until master_log_file = 'a', master_log_pos = 1", true},
	{"start slave until sql_after_gtids = 'a'", true},
	{"update select.a set b = 1", true},

	// In a routine or a trigger.
	{"create procedure p() begin declare a int; set a = 1, autocommit = 1; end", true},
	{"create procedure p() begin declare a int; set a = 1, global autocommit = 1; end", true},
	{"create procedure p() begin declare a int; set a = 1; end", true},
	{"create procedure p() begin declare a int; set a = binary; end", true},
	{"create procedure p() begin declare a int; set a = on; end", true},
	{"create procedure p() begin declare a int; set global autocommit = 1, a = 1; end", true},
	{"create procedure p() begin declare timestamp int; set timestamp = 1; end", true},
	{"create procedure p() begin declare transaction int; set transaction = 1; end", true},
	{"create procedure p() begin set @a = on; end", false},
	{"create procedure p() begin set autocommit = 1; end", true},
	{"create trigger t before insert on a for each row set new.a = 1", true},
	{"create trigger t before insert on a for each row set new.a = 1, old.b = 2", true},
	{"create trigger t before insert on a for each row set new.a = binary", true},
	{"create trigger t before insert on a for each row set new.a = default", true},
	{"create trigger t before insert on a for each row set new.a = on", true},
	{"create trigger t before insert on a for each row set timestamp = 1", true},
}
