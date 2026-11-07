package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// What a SET list holds besides names and scopes, and the word before a dot, as MySQL 8.0 reads them,
// every verdict measured with PREPARE on mysql:8.0 (8.0.46) on 2026-11-05: a 1064 is the parser's refusal, and
// any other answer is a statement that parsed. A 1193 or 4081 is an unknown variable, raised while
// parsing, so a row holds one only where no text follows the assignment that raised it.
//
// NAMES and the charset forms are option_value_no_option_type, so they may sit anywhere in the list.
// NAMES takes COLLATE only after a name, and a collation is a string, a name or BINARY.
// set_expr_or_default is expr, DEFAULT, ON, ALL, BINARY, ROW or SYSTEM.
//
// A word running into a dot and a name is a name in every statement, whatever the word, so `SELECT
// select.a FROM t` names a column. See beforeDotName.
//
// Not in this table, because omni still disagrees and each is its own change. `SELECT DEFAULT`, because parseExpr reads it everywhere but SET. `SET a = DEFAULT` for a
// procedure's local variable and for NEW.a in a trigger, refused in an action that sees the routine. `SET a.1 = 1`, where the lexer reads .1 as a number. `SET PASSWORD = PASSWORD('x')`,
// which omni accepts.

// TestSetListForms holds omni to the recorded verdicts without a server, so a run with no Docker
// still reports a widened or narrowed grammar.
func TestSetListForms(t *testing.T) {
	for _, c := range setListCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// firstSetAssignment returns the first assignment of a SET.
func firstSetAssignment(t *testing.T, sql string) *nodes.Assignment {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("%q: %v", sql, err)
	}
	return list.Items[0].(*nodes.SetStmt).Assignments[0]
}

// TestSetListTrees pins what a list becomes, because a verdict cannot tell a word from a value.
func TestSetListTrees(t *testing.T) {
	columns := func(t *testing.T, sql string) []string {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		stmt, ok := list.Items[0].(*nodes.SetStmt)
		if !ok {
			t.Fatalf("got %T", list.Items[0])
		}
		if stmt.Loc.Start != 0 || stmt.Loc.End != len(sql) {
			t.Errorf("statement Loc = %+v, want {0 %d}", stmt.Loc, len(sql))
		}
		var got []string
		for _, a := range stmt.Assignments {
			got = append(got, a.Column.Column)
		}
		return got
	}
	for sql, want := range map[string][]string{
		"set autocommit = 1, names utf8":         {"autocommit", "NAMES"},
		"set names utf8, autocommit = 1":         {"NAMES", "autocommit"},
		"set @a = 1, charset utf8":               {"@a", "CHARACTER SET"},
		"set names utf8 collate utf8_bin, a = 1": {"NAMES", "COLLATE", "a"},
		"set names 'utf8'":                       {"NAMES"},
		"set names binary":                       {"NAMES"},
		"set char set utf8":                      {"CHARACTER SET"},
		"set character set 'utf8mb4'":            {"CHARACTER SET"},
		"set charset = 1":                        {"charset"},
	} {
		got := columns(t, sql)
		if len(got) != len(want) {
			t.Errorf("%q: columns = %v, want %v", sql, got, want)
			continue
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("%q: columns = %v, want %v", sql, got, want)
			}
		}
	}

	if v, ok := firstSetAssignment(t, "set names 'utf8'").Value.(*nodes.StringLit); !ok || v.Value != "utf8" {
		t.Errorf("the charset is %#v", firstSetAssignment(t, "set names 'utf8'").Value)
	}

	// The values. ON stays the true the review walk reads, and the other words are text.
	if v, ok := firstSetAssignment(t, "set a = on").Value.(*nodes.BoolLit); !ok || !v.Value {
		t.Errorf("ON is %#v", firstSetAssignment(t, "set a = on").Value)
	}
	for sql, want := range map[string]string{"set a = binary": "BINARY", "set a = all": "ALL", "set a = row": "ROW", "set a = system": "SYSTEM", "set @@a = binary": "BINARY"} {
		if v, ok := firstSetAssignment(t, sql).Value.(*nodes.StringLit); !ok || v.Value != want {
			t.Errorf("%q: value is %#v", sql, firstSetAssignment(t, sql).Value)
		}
	}
	if _, ok := firstSetAssignment(t, "set a = default").Value.(*nodes.DefaultExpr); !ok {
		t.Errorf("DEFAULT is %T", firstSetAssignment(t, "set a = default").Value)
	}
	if d, ok := firstSetAssignment(t, "set a = default(b)").Value.(*nodes.DefaultExpr); !ok || d.Column != "b" {
		t.Errorf("DEFAULT(b) is %#v", firstSetAssignment(t, "set a = default(b)").Value)
	}
	if _, ok := firstSetAssignment(t, "set a = binary 'x'").Value.(*nodes.StringLit); ok {
		t.Errorf("BINARY 'x' read as the word BINARY")
	}
}

// TestWordBeforeDot pins what the lexer's rule turns a word into, against the readers of @@scope
// and of a trigger's qualifier, which need the keyword to stay.
func TestWordBeforeDot(t *testing.T) {
	target := func(t *testing.T, sql string) nodes.ExprNode {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		item := list.Items[0].(*nodes.SelectStmt).TargetList[0]
		if res, ok := item.(*nodes.ResTarget); ok {
			return res.Val
		}
		return item
	}
	if col, ok := target(t, "select select.a from t").(*nodes.ColumnRef); !ok || col.Table != "select" || col.Column != "a" {
		t.Errorf("select.a read as %#v", target(t, "select select.a from t"))
	}
	if col, ok := target(t, "select count.a from t").(*nodes.ColumnRef); !ok || col.Table != "count" {
		t.Errorf("count.a read as %#v", target(t, "select count.a from t"))
	}
	if ref, ok := target(t, "select @@global.a").(*nodes.VariableRef); !ok || ref.Scope != "GLOBAL" || ref.Name != "a" {
		t.Errorf("@@global.a read as %#v", target(t, "select @@global.a"))
	}
	for sql, wantKeyword := range map[string]bool{"select.a": false, "select .a": true, "select. a": true, "select.": true, "select.*": true, "new.a": true, "old.a": true} {
		if got := NewLexer(sql).NextToken().Type != tokIDENT; got != wantKeyword {
			t.Errorf("%q: keyword=%v, want %v", sql, got, wantKeyword)
		}
	}
}

var setListCases = []asOfCase{
	{"alter table t lock = none", true},
	{"alter table t page_compressed = on", false},
	{"alter table t stats_persistent = 0", true},
	{"alter table where.a add column b int", true},
	{"change master to master_host = 'a'", false},
	{"create index i on where.a (b)", true},
	{"create procedure p() begin declare a int; set a = 1, autocommit = 1; end", true},
	{"create procedure p() begin declare a int; set a = 1, global autocommit = 1; end", true},
	{"create procedure p() begin declare a int; set a = binary; end", true},
	{"create procedure p() begin declare a int; set a = on; end", true},
	{"create procedure p() begin declare a int; set global autocommit = 1, a = 1; end", true},
	{"create procedure p() begin declare timestamp int; set timestamp = 1; end", true},
	{"create procedure p() begin set @a = 1, names utf8; end", true},
	{"create procedure p() begin set @a = default; end", false},
	{"create procedure p() begin set @a = on; end", false},
	{"create procedure p() begin set autocommit = 1; end", true},
	{"create procedure p() begin set autocommit = all; end", true},
	{"create procedure p() begin set autocommit = binary; end", true},
	{"create procedure p() begin set names utf8; end", true},
	{"create table offset (a int)", true},
	{"create table select.a (x int)", true},
	{"create table select.from (x int)", true},
	{"create table t (a int) page_compressed = on", false},
	{"create table t (a int) sequence = on", false},
	{"create table t (a int) stats_auto_recalc = 1", true},
	{"create table t (a int) stats_persistent = 1", true},
	{"create table t (a int) stats_persistent = on", false},
	{"create table t (a int) stats_sample_pages = 10", true},
	{"create table t (a int) with system versioning", false},
	{"create table t (a int, primary.b int)", false},
	{"create table t (a vector(3))", false},
	{"create table t (vector int)", true},
	{"create table vector (a int)", true},
	{"create trigger t before insert on a for each row set new.a = 1", true},
	{"create trigger t before insert on a for each row set new.a = 1, new.b = 2", true},
	{"create trigger t before insert on a for each row set new.a = 1, old.b = 2", false},
	{"create trigger t before insert on a for each row set new.a = binary", true},
	{"create trigger t before insert on a for each row set new.a = on", true},
	{"create trigger t before insert on a for each row set timestamp = 1", true},
	{"delete from select.a", true},
	{"delete from where.a", true},
	{"drop table select.a", true},
	{"drop table where.a", true},
	{"fetch c into timestamp", false},
	{"insert into select.a values (1)", true},
	{"insert into where.a values (1)", true},
	{"reset persist a", true},
	{"reset persist timestamp", true},
	{"select * from from.select", true},
	{"select * from select.a", true},
	{"select * from t for system_time all", false},
	{"select * from t limit 1 offset 2", true},
	{"select * from t limit 1, 2", true},
	{"select * from t offset", true},
	{"select 1 from select.a", true},
	{"select 1 from where.where", true},
	{"select 1. from t", true},
	{"select @@ autocommit", false},
	{"select @@`autocommit`", true},
	{"select @@`autocommit`.a", true},
	{"select @@autocommit . a", true},
	{"select @@global . autocommit", true},
	{"select @@global.a", true},
	{"select @@global.autocommit", true},
	{"select @@local.autocommit", true},
	{"select @@select", false},
	{"select @@select.a", false},
	{"select @@session.`autocommit`", true},
	{"select @@session.autocommit", true},
	{"select @a := on", false},
	{"select `select`.a from t", true},
	{"select a . select from t", false},
	{"select a from t where where.b = 1", true},
	{"select a from t where.b", false},
	{"select a offset from t", true},
	{"select a.select from t", true},
	{"select a.select.b from t", true},
	{"select all", false},
	{"select binary", false},
	{"select case.a from t", true},
	{"select count.a from t", true},
	{"select date.a from t", true},
	{"select default(a) from t", true},
	{"select from.a from t", true},
	{"select global.a", true},
	{"select if.a from t", true},
	{"select interval.a from t", true},
	{"select local.a", true},
	{"select new.a", true},
	{"select not.a from t", true},
	{"select now.a from t", true},
	{"select null.a from t", true},
	{"select off", true},
	{"select off from t", true},
	{"select offset from t", true},
	{"select old.a", true},
	{"select on", false},
	{"select persist.a", true},
	{"select row", false},
	{"select row from t", false},
	{"select select .a from t", false},
	{"select select.", false},
	{"select select. a from t", false},
	{"select select.`a` from t", false},
	{"select select.a from t", true},
	{"select select.a.b from t", true},
	{"select session.a", true},
	{"select show.* from t", false},
	{"select show.a from t", true},
	{"select system", false},
	{"select system from t", false},
	{"select t.where from t", true},
	{"select table.* from table.a", false},
	{"select true.a from t", true},
	{"select where.a from t", true},
	{"select where.where from where", false},
	{"set , autocommit = 1", false},
	{"set @@ autocommit = 1", false},
	{"set @@ global.autocommit = 1", false},
	{"set @@'a' = 1", false},
	{"set @@a . b = 1", true},
	{"set @@autocommit = 1", true},
	{"set @@autocommit = 1, names utf8", true},
	{"set @@autocommit = all", true},
	{"set @@autocommit = binary", true},
	{"set @@autocommit = default", true},
	{"set @@autocommit = on", true},
	{"set @@autocommit = row", true},
	{"set @@autocommit = system", true},
	{"set @@global . autocommit = 1", true},
	{"set @@global.autocommit = 1", true},
	{"set @@local.autocommit = 1", true},
	{"set @@persist.autocommit = 1", true},
	{"set @@persist_only.autocommit = 1", true},
	{"set @@select.a = 1", false},
	{"set @@session.autocommit = 1", true},
	{"set @a = 1, autocommit = 1", true},
	{"set @a = 1, character set utf8", true},
	{"set @a = 1, charset utf8", true},
	{"set @a = 1, names utf8", true},
	{"set @a = all", false},
	{"set @a = binary", false},
	{"set @a = binary 'x'", true},
	{"set @a = default", false},
	{"set @a = default(b)", true},
	{"set @a = off", true},
	{"set @a = on", false},
	{"set @a = row", false},
	{"set @a = row(1, 2)", true},
	{"set @a = system", false},
	{"set `select` = 1", true},
	{"set a = 1 , names utf8 collate utf8_bin, b = 2", true},
	{"set a = 1, char set utf8", true},
	{"set a.timestamp = 1", true},
	{"set autocommit", false},
	{"set autocommit . select = 1", false},
	{"set autocommit := 1", true},
	{"set autocommit =", false},
	{"set autocommit = 'off'", true},
	{"set autocommit = 'on'", true},
	{"set autocommit = (on)", false},
	{"set autocommit = 1", true},
	{"set autocommit = 1 + 1", true},
	{"set autocommit = 1,", false},
	{"set autocommit = 1, @@autocommit = 1", true},
	{"set autocommit = 1, @a = 1", true},
	{"set autocommit = 1, autocommit = 0", true},
	{"set autocommit = 1, character set utf8", true},
	{"set autocommit = 1, charset = 2", true},
	{"set autocommit = 1, charset utf8", true},
	{"set autocommit = 1, global autocommit = 1", true},
	{"set autocommit = 1, global.a = 1", false},
	{"set autocommit = 1, local autocommit = 1", true},
	{"set autocommit = 1, names = 2", false},
	{"set autocommit = 1, names utf8", true},
	{"set autocommit = 1, password := 2", true},
	{"set autocommit = 1, password = 1", true},
	{"set autocommit = 1, password.a = 2", true},
	{"set autocommit = 1, persist = 1", false},
	{"set autocommit = 1, session autocommit = 1", true},
	{"set autocommit = 1, session names = 2", true},
	{"set autocommit = 1, session password = 1", true},
	{"set autocommit = 1, timestamp = 1", true},
	{"set autocommit = 1, transaction = 1", true},
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
	{"set autocommit = off", true},
	{"set autocommit = on", true},
	{"set autocommit = on + 1", false},
	{"set autocommit = on, autocommit = 1", true},
	{"set autocommit = row", true},
	{"set autocommit = row + 1", false},
	{"set autocommit = row(1, 2)", true},
	{"set autocommit = row, autocommit = 1", true},
	{"set autocommit = system", true},
	{"set autocommit = system + 1", false},
	{"set autocommit = system, autocommit = 1", true},
	{"set binding = 1", true},
	{"set char set 'utf8'", true},
	{"set char set binary", true},
	{"set char set default", true},
	{"set char set utf8", true},
	{"set character set 'utf8'", true},
	{"set character set 'utf8mb4'", true},
	{"set character set binary", true},
	{"set character set default", true},
	{"set character set utf8", true},
	{"set character set utf8, autocommit = 1", true},
	{"set charset 'utf8'", true},
	{"set charset binary", true},
	{"set charset default", true},
	{"set charset utf8", true},
	{"set charset utf8 collate utf8_bin", false},
	{"set charset utf8, autocommit = 1", true},
	{"set date = 1", true},
	{"set default = 1", false},
	{"set default role a", false},
	{"set default role a for b", false},
	{"set default role none", false},
	{"set from = 1", false},
	{"set from.a = 1", true},
	{"set global @@autocommit = 1", false},
	{"set global autocommit = 1", true},
	{"set global autocommit = 1, @a = 1, autocommit = 1", true},
	{"set global autocommit = 1, autocommit = 0", true},
	{"set global autocommit = 1, names utf8", true},
	{"set global autocommit = 1, session autocommit = 0", true},
	{"set global autocommit = 1, session autocommit = 0, autocommit = 1", true},
	{"set global autocommit = all", true},
	{"set global autocommit = binary", true},
	{"set global names utf8", false},
	{"set global timestamp = 1", true},
	{"set global transaction read only", true},
	{"set global.autocommit = 1", false},
	{"set local autocommit = 1", true},
	{"set local transaction read only", true},
	{"set names 'utf8'", true},
	{"set names @a", false},
	{"set names binary", true},
	{"set names binary collate binary", true},
	{"set names binary collate default", false},
	{"set names default", true},
	{"set names default collate utf8_bin", false},
	{"set names utf8", true},
	{"set names utf8 collate 'utf8_bin'", true},
	{"set names utf8 collate binary", true},
	{"set names utf8 collate default", false},
	{"set names utf8 collate utf8_bin", true},
	{"set names utf8 collate utf8_bin, autocommit = 1", true},
	{"set names utf8, autocommit = 1", true},
	{"set names utf8, charset utf8", true},
	{"set names utf8, names latin1", true},
	{"set new.a.b = 1", false},
	{"set password = 'x'", true},
	{"set password = old_password('x')", false},
	{"set password for u = 'x'", true},
	{"set persist autocommit = 1", true},
	{"set persist_only autocommit = 1", true},
	{"set resource group rg", true},
	{"set role a", true},
	{"set role none", true},
	{"set select = 1", false},
	{"set select. a = 1", false},
	{"set select.a = 1", true},
	{"set session autocommit = 1", true},
	{"set session global = 1", false},
	{"set session names = 1", true},
	{"set session password = 1", true},
	{"set session session = 1", false},
	{"set session timestamp = 1", true},
	{"set session transaction read only", true},
	{"set statement = 1", true},
	{"set statement autocommit = 1 for select 1", false},
	{"set status = 1", true},
	{"set time = 1", true},
	{"set timestamp = 1", true},
	{"set timestamp.a = 1", true},
	{"set transaction isolation level serializable", true},
	{"set transaction read only", true},
	{"set transaction read write", true},
	{"set user = 1", true},
	{"show columns from select.a", true},
	{"show columns from where.a", true},
	{"show tables from select", false},
	{"start slave until master_log_file = 'a', master_log_pos = 1", false},
	{"start slave until sql_after_gtids = 'a'", false},
	{"truncate table where.a", true},
	{"update select.a set b = 1", true},
	{"update where.a set b = 1", true},
}
