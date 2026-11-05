package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// A word after a dot, as TiDB v8.5.5 reads it, every verdict measured with PREPARE on
// pingcap/tidb:v8.5.5 on 2026-11-05: a 1064 is the parser's refusal, and any other answer
// means the text parsed.
//
// The rule is isTokenIdentifier in upstream's scanner: a word right before a dot, or after
// a dot and nothing but spaces, is not a keyword. So `SELECT a . select FROM t` names a
// column, and `SELECT o. FROM t` reads FROM as a column and t as its alias. Parse applies
// it. Tokenize and completion do not, because they read incomplete text at a caret, where
// `o. FROM t` means FROM is the next clause (see dotNamesWord).
//
// Not in this table, because omni still disagrees and the rule is not why: `SELECT NEXT
// VALUE FOR a` and `CREATE TABLE t (a.b INT)`, which omni does not read however they are
// spaced. `SHOW INDEX FROM a.b FROM c`, which it accepts. `SELECT @@a.1`, which it refuses
// because the lexer reads .1 as a number. `SELECT INTERVAL 1. DAY + NOW()`, which it accepts
// because it reads an interval's unit by spelling.

// TestSpacedDotForms holds omni to the recorded verdicts without a server, so a run with no
// Docker still reports a widened or narrowed grammar.
func TestSpacedDotForms(t *testing.T) {
	for _, c := range spacedDotCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// firstTarget returns the first select item of sql, with its alias.
func firstTarget(t *testing.T, sql string) (nodes.ExprNode, string) {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("got %T", list.Items[0])
	}
	if res, ok := stmt.TargetList[0].(*nodes.ResTarget); ok {
		return res.Val, res.Name
	}
	return stmt.TargetList[0], ""
}

// TestSpacedDotTrees pins what the word turned into, because a verdict cannot tell a column
// named FROM from a FROM clause.
func TestSpacedDotTrees(t *testing.T) {
	for _, c := range []struct{ sql, table, column, alias string }{
		{"select a . select from t", "a", "select", ""},
		{"select a.  select from t", "a", "select", ""},
		{"select o. from t", "o", "from", "t"},
		{"select o . from t", "o", "from", "t"},
		{"select a . from from t", "a", "from", ""},
		{"select a . as b from t", "a", "as", "b"},
		{"select a . b from t", "a", "b", ""},
	} {
		t.Run(c.sql, func(t *testing.T) {
			val, alias := firstTarget(t, c.sql)
			col, ok := val.(*nodes.ColumnRef)
			if !ok {
				t.Fatalf("value is %T", val)
			}
			if col.Table != c.table || col.Column != c.column || alias != c.alias {
				t.Errorf("table=%q column=%q alias=%q, want %q %q %q", col.Table, col.Column, alias, c.table, c.column, c.alias)
			}
		})
	}

	t.Run("select 1. select from t", func(t *testing.T) {
		if _, alias := firstTarget(t, "select 1. select from t"); alias != "select" {
			t.Errorf("alias = %q, want select", alias)
		}
	})

	for _, c := range []struct{ sql, name, scope, alias string }{
		{"select @@a.b", "a.b", "", ""},
		{"select @@a.b c", "a.b", "", "c"},
		{"select @@a. b", "a.", "", "b"},
		{"select @@global.a", "a", "GLOBAL", ""},
		{"select @@global.a.b", "a.b", "GLOBAL", ""},
		{"select @@`a`", "a", "", ""},
	} {
		t.Run(c.sql, func(t *testing.T) {
			val, alias := firstTarget(t, c.sql)
			ref, ok := val.(*nodes.VariableRef)
			if !ok {
				t.Fatalf("value is %T", val)
			}
			if ref.Name != c.name || ref.Scope != c.scope || alias != c.alias {
				t.Errorf("name=%q scope=%q alias=%q, want %q %q %q", ref.Name, ref.Scope, alias, c.name, c.scope, c.alias)
			}
		})
	}

	// A table's name, in the positions that hold one.
	list, err := Parse("drop table a . table")
	if err != nil {
		t.Fatal(err)
	}
	if ref := list.Items[0].(*nodes.DropTableStmt).Tables[0]; ref.Schema != "a" || ref.Name != "table" {
		t.Errorf("schema=%q name=%q", ref.Schema, ref.Name)
	}

	// A SHOW's database stays where its readers look.
	for sql, want := range map[string]string{
		"show tables from a":          "a",
		"show full tables in a":       "a",
		"show table status from a":    "a",
		"show create database a":      "a",
		"show open tables from a":     "a",
		"show events from `a`":        "a",
		"show triggers from a":        "a",
		"show tables from a like 'x'": "a",
	} {
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("%q: %v", sql, err)
		}
		from := list.Items[0].(*nodes.ShowStmt).From
		if from == nil || from.Name != want || from.Schema != "" {
			t.Errorf("%q: From = %+v, want %q", sql, from, want)
		}
	}
	list, err = Parse("show columns from t from a")
	if err != nil {
		t.Fatal(err)
	}
	if from := list.Items[0].(*nodes.ShowStmt).From; from.Schema != "a" || from.Name != "t" {
		t.Errorf("SHOW COLUMNS From = %+v", from)
	}
}

// TestSpacedDotKeepsParseAndCompletionApart holds the two readers to their own rule: the
// strict one reads a name and the one at a caret a clause.
func TestSpacedDotKeepsParseAndCompletionApart(t *testing.T) {
	if tokens := Tokenize("a. from"); len(tokens) != 3 || tokens[2].Type != kwFROM {
		t.Errorf("Tokenize read FROM after a spaced dot as %v", tokens)
	}
	if cs := Collect("SELECT t. FROM x", 9); cs == nil || !cs.HasRule("columnref") {
		t.Errorf("completion after t. offers no column: %+v", cs)
	}
	expr, err := ParseExpr("a . select")
	if err != nil {
		t.Fatalf("ParseExpr: %v", err)
	}
	if col, ok := expr.(*nodes.ColumnRef); !ok || col.Table != "a" || col.Column != "select" {
		t.Errorf("ParseExpr read %#v", expr)
	}
}

var spacedDotCases = []tableSampleCase{
	// The rule across statements, and what ends it: a tab or a comment before the word, and a
	// dot that is part of nothing.
	{"alter table a . select add column x int", true},
	{"call a . select()", true},
	{"create database a . select", false},
	{"create index a . select on t (a)", false},
	{"create procedure a . select() begin end", true},
	{"create sequence a . select", true},
	{"create table a . select (x int)", true},
	{"create table t (a int, b int as (t . select))", true},
	{"create table t (a int, primary key (a), foreign key (a) references b . select (c))", true},
	{"create trigger a . select before insert on t for each row set @a = 1", false},
	{"create view a . select as select 1", true},
	{"delete from a . select", true},
	{"drop table a . select", true},
	{"drop table a . table", true},
	{"drop view a . select", true},
	{"explain select a . select from t", true},
	{"grant all on a . from . select to u", false},
	{"grant select on a . * to u", true},
	{"grant select on a . select to u", true},
	{"insert into a . select values (1)", true},
	{"insert into t (a . select) values (1)", true},
	{"load data infile 'x' into table a . select", true},
	{"rename table a . select to b . from", true},
	{"select \"a.\" from t", true},
	{"select 'a.' from t", true},
	{"select * from a . b . select", false},
	{"select * from a . from", true},
	{"select * from a . from where 1", true},
	{"select * from a . select", true},
	{"select * from a . select . b", false},
	{"select * from db . of", true},
	{"select * from db . tablesample", true},
	{"select * from t a join b on a . join = 1", true},
	{"select * from t a join b on a . on = b . on", true},
	{"select * from t where a . from = 1", true},
	{"select * from t where a . in (1)", true},
	{"select * from t where a . in in (1)", true},
	{"select 1 . from t", false},
	{"select @a . select", false},
	{"select `a.` from t", true},
	{"select `a.` select from t", false},
	{"select a .", false},
	{"select a .\tselect from t", false},
	{"select a .  select from t", true},
	{"select a . * from t", true},
	{"select a . /* c */ b from t", true},
	{"select a . 1 from t", false},
	{"select a . as b from t", true},
	{"select a . b (1) from t", true},
	{"select a . b . * from t", true},
	{"select a . b . c from t", true},
	{"select a . b . select from t", true},
	{"select a . case from t", true},
	{"select a . count (1) from t", true},
	{"select a . current_user from t", true},
	{"select a . default from t", true},
	{"select a . from from t", true},
	{"select a . if (1, 2, 3) from t", true},
	{"select a . interval from t", true},
	{"select a . max (1) from t", true},
	{"select a . not from t", true},
	{"select a . now() from t", true},
	{"select a . null from t", true},
	{"select a . select (1) from t", true},
	{"select a . select . b from t", true},
	{"select a . select from t", true},
	{"select a . select(1) from t", true},
	{"select a . select, a . from from t", true},
	{"select a . true from t", true},
	{"select a . user () from t", true},
	{"select a. ", false},
	{"select a.  select from t", true},
	{"select a. * from t", true},
	{"select a. /* c */ from t", false},
	{"select a. as from t", true},
	{"select a. b. c from t", true},
	{"select a. select b from t", true},
	{"select a.b . select from t", true},
	{"select a/* . */ from t", true},
	{"select count(a . select) from t", true},
	{"select interval 1 foo + now()", false},
	{"select lastval(a . select)", true},
	{"select max(a . max) from t", true},
	{"select nextval(a . select)", true},
	{"select now() + interval 1 day", true},
	{"select o . from t", true},
	{"select o .from t", true},
	{"select o. a from t", true},
	{"select o. from from t", true},
	{"select o. from t", true},
	{"select o. select from t", true},
	{"select setval(a . select, 1)", true},
	{"select t . of from t", true},
	{"select t . tablesample from t", true},
	{"select t. * , b from t", true},
	{"select t. from from t", true},
	{"select t. of from t", true},
	{"select t.a. from t", true},
	{"truncate table a . select", true},
	{"update a . select set b = 1", true},
	{"update t set a . select = 1", true},
	{"use a . select", false},
	{"with a as (select 1) select a . select from a", true},

	// The dot that ends a number counts too, because the rule reads the bytes.
	{"create table t (a decimal(10,2) default 1. comment 'x')", false},
	{"create table t (a decimal(10,2) default 1. not null)", false},
	{"insert into t values (1.) on duplicate key update a = 1", true},
	{"select .5 from t", true},
	{"select 1. + 2", true},
	{"select 1. a", true},
	{"select 1. a from t", true},
	{"select 1. and 2", false},
	{"select 1. as from t", true},
	{"select 1. from", true},
	{"select 1. from t", false},
	{"select 1. select from t", true},
	{"select 1. tablesample", true},
	{"select 1. union select 2", false},
	{"select 1.0 from t", true},
	{"select 1.e1 from t", true},
	{"select a from t group by 1. having 1", false},
	{"select a from t limit 1. offset 2", false},
	{"select a from t order by 1. desc", false},
	{"select a from t where b = 1. and c = 2", false},
	{"select a from t where b = 1. or c = 2", false},
	{"select cast(1. as char)", false},
	{"select interval 1. foo + now()", false},

	// SET, whose names reach the same rule.
	{"set @a . select = 1", false},
	{"set a .  select = 1", true},
	{"set a . /* c */ select = 1", false},
	{"set a . from = 1", true},
	{"set a . select = 1", true},
	{"set a. select = 1", true},
	{"set global a . select = 1", true},
	{"set transaction . select = 1", true},
	{"set x = 1, a . select = 2", true},

	// @@ and its name are one token upstream, so a space ends it. A trailing dot stays in the
	// token, and a backquoted name is the whole of it.
	{"select 1 + @@a * 2", true},
	{"select @@ a", false},
	{"select @@ a.b", false},
	{"select @@ global.a", false},
	{"select @@/* c */a", false},
	{"select @@`a`", true},
	{"select @@`a`.`b`", false},
	{"select @@`a`.b", false},
	{"select @@`global`.a", false},
	{"select @@a", true},
	{"select @@a . b", false},
	{"select @@a . select", false},
	{"select @@a .b", false},
	{"select @@a /* c */ . b", false},
	{"select @@a as b", true},
	{"select @@a b", true},
	{"select @@a+1", true},
	{"select @@a.", true},
	{"select @@a. as b", false},
	{"select @@a. b", true},
	{"select @@a.*", false},
	{"select @@a., 1", true},
	{"select @@a.`b`", true},
	{"select @@a.b", true},
	{"select @@a.b c", true},
	{"select @@a.b+1", true},
	{"select @@a.b.c", true},
	{"select @@a.select", true},
	{"select @@a.select.b", true},
	{"select @@a=1", true},
	{"select @@global . a", false},
	{"select @@global .a", false},
	{"select @@global.", false},
	{"select @@global. a", false},
	{"select @@global.,1", false},
	{"select @@global.`a`", true},
	{"select @@global.a", true},
	{"select @@global.a . b", false},
	{"select @@global.a .b", false},
	{"select @@global.a, @@session.b", true},
	{"select @@global.select", true},
	{"select @@local.a", true},
	{"select @@persist . a", false},
	{"select @@persist.a", true},
	{"select @@select", true},
	{"select @@select.a", true},
	{"select @@session .a", false},
	{"select @@session.a.b", true},
	{"set @@ a = 1", false},
	{"set @@a . b = 1", false},
	{"set @@a .b = 1", false},
	{"set @@a. b = 1", false},
	{"set @@a.b = 1", true},
	{"set @@global . a = 1", false},
	{"set @@global .a = 1", false},
	{"set @@global. a = 1", false},
	{"set @@global.a . b = 1", false},
	{"set @@global.a = 1", true},
	{"set @@global.a.b = 1", true},

	// A SHOW's database is one Identifier, ShowDatabaseNameOpt, parser.y:12128, where its table
	// is a TableName.
	{"desc a . select b", true},
	{"describe a . select", true},
	{"explain a . select", true},
	{"show bindings from a.b", false},
	{"show columns from a . b", true},
	{"show columns from a . select", true},
	{"show columns from a.b", true},
	{"show columns from t from a", true},
	{"show columns from t from a . b", false},
	{"show columns from t from a.b", false},
	{"show create database a . b", false},
	{"show create sequence a . select", true},
	{"show create table a . select", true},
	{"show create user a . b", false},
	{"show create view a . select", true},
	{"show events from a", true},
	{"show events from a . b", false},
	{"show events from a.b", false},
	{"show full tables from a", true},
	{"show full tables from a . b", false},
	{"show full tables in a.b", false},
	{"show function status from a", false},
	{"show index from a . b", true},
	{"show index from t from a", true},
	{"show index from t from a.b", false},
	{"show keys from t in a", true},
	{"show keys from t in a.b", false},
	{"show open tables from a", true},
	{"show open tables from a.b", false},
	{"show placement for table a.b", true},
	{"show procedure status from a", false},
	{"show stats_histograms from a.b", false},
	{"show stats_meta from a", false},
	{"show stats_meta from a.b", false},
	{"show table a . select next_row_id", true},
	{"show table a . select regions", true},
	{"show table status from a", true},
	{"show table status from a . b", false},
	{"show table status from a.b", false},
	{"show table status in a like 'x'", true},
	{"show tables from 'a'", false},
	{"show tables from `a`", true},
	{"show tables from a", true},
	{"show tables from a . b", false},
	{"show tables from a . select", false},
	{"show tables from a .b", false},
	{"show tables from a as b", false},
	{"show tables from a b", false},
	{"show tables from a like 'x'", true},
	{"show tables from a where 1", true},
	{"show tables from a. b", false},
	{"show tables from a.b", false},
	{"show tables from select", false},
	{"show tables in a . b", false},
	{"show triggers from a", true},
	{"show triggers from a . b", false},
	{"show triggers from a.b", false},
}
