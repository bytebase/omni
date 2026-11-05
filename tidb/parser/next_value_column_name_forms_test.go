package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// NEXT VALUE FOR and qualified column names, as TiDB v8.5.5 reads them, every verdict measured
// with PREPARE on pingcap/tidb:v8.5.5 on 2026-11-05: a 1064 is the parser's refusal, and any other
// answer means the text parsed.
//
// NextValueForSequence, parser.y:4120, is NEXT VALUE FOR TableName, and NEXTVAL '(' TableName ')' is
// the same tree. VALUE after NEXT commits to the form, so `SELECT next value FROM t` is a 1064 where
// `SELECT next a FROM t` is a column with an alias.
//
// ColumnName, parser.y:3511, takes up to two qualifiers, and it is what a column definition, a key
// part, MODIFY, CHANGE, ALTER, DROP, AFTER and REFERENCES all read. RENAME COLUMN reads a bare name.
// What a qualifier must match is not the parser's to say: the catalog holds that.
//
// Not in this table, because omni still disagrees: `SELECT next(1)`, `SELECT lastval(a, 1)` and
// `SELECT nextval(a) OVER ()`, which it accepts, because it lets any unreserved word be called and
// reads OVER after any call. `ALTER TABLE t ALTER t.a SET NOT NULL` and `ALTER TABLE t ALTER a SET INVISIBLE`, which it accepts, the second a form
// MySQL has.

// TestNextValueColumnNameForms holds omni to the recorded verdicts without a server, so a run with
// no Docker still reports a widened or narrowed grammar.
func TestNextValueColumnNameForms(t *testing.T) {
	for _, c := range nextValueColumnNameCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestNextValueTrees pins that both spellings give one tree, which is what a consumer reads, and
// that NEXT stays a name where VALUE does not follow it.
func TestNextValueTrees(t *testing.T) {
	first := func(t *testing.T, sql string) nodes.ExprNode {
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
	for _, c := range []struct{ sql, table, column string }{
		{"select next value for s", "", "s"},
		{"select next value for d.s", "d", "s"},
		{"select next value for d . select", "d", "select"},
		{"select NEXT VALUE FOR `s`", "", "s"},
		{"select nextval(s)", "", "s"},
		{"select nextval(d.s)", "d", "s"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			fn, ok := first(t, c.sql).(*nodes.FuncCallExpr)
			if !ok {
				t.Fatalf("value is %T", first(t, c.sql))
			}
			if len(fn.Args) != 1 {
				t.Fatalf("args = %v", fn.Args)
			}
			arg, ok := fn.Args[0].(*nodes.ColumnRef)
			if !ok || arg.Table != c.table || arg.Column != c.column {
				t.Errorf("sequence is %#v, want table=%q column=%q", fn.Args[0], c.table, c.column)
			}
			if fn.Name != "NEXTVAL" && fn.Name != "nextval" {
				t.Errorf("name = %q", fn.Name)
			}
		})
	}
	const sql = "select 1 + next value for d.s"
	list, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	item := list.Items[0].(*nodes.SelectStmt).TargetList[0]
	if res, ok := item.(*nodes.ResTarget); ok {
		item = res.Val
	}
	fn, ok := item.(*nodes.BinaryExpr).Right.(*nodes.FuncCallExpr)
	if !ok {
		t.Fatalf("right side is %T", item.(*nodes.BinaryExpr).Right)
	}
	if got := sql[fn.Loc.Start:fn.Loc.End]; got != "next value for d.s" {
		t.Errorf("call text = %q", got)
	}
	if col, ok := first(t, "select next from t").(*nodes.ColumnRef); !ok || col.Column != "next" {
		t.Errorf("NEXT alone is %#v", first(t, "select next from t"))
	}
}

// TestColumnNameTrees pins where the qualifiers land, since the catalog reads them.
func TestColumnNameTrees(t *testing.T) {
	list, err := Parse("create table t (a int, t.b int, d.t.c int)")
	if err != nil {
		t.Fatal(err)
	}
	cols := list.Items[0].(*nodes.CreateTableStmt).Columns
	for i, want := range []struct{ schema, table, name string }{{"", "", "a"}, {"", "t", "b"}, {"d", "t", "c"}} {
		if cols[i].Schema != want.schema || cols[i].Table != want.table || cols[i].Name != want.name {
			t.Errorf("column %d = %q.%q.%q, want %+v", i, cols[i].Schema, cols[i].Table, cols[i].Name, want)
		}
	}
	if cols[2].Loc.Start != len("create table t (a int, t.b int, ") {
		t.Errorf("a column's Loc starts at %d, before its qualifier", cols[2].Loc.Start)
	}

	cmd := func(t *testing.T, sql string) *nodes.AlterTableCmd {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		return list.Items[0].(*nodes.AlterTableStmt).Commands[0]
	}
	if c := cmd(t, "alter table t change d.t.a t.b int"); c.Name != "a" || c.OldColumnSchema != "d" || c.OldColumnTable != "t" || c.Column.Table != "t" || c.Column.Name != "b" {
		t.Errorf("CHANGE = %+v, new %+v", c, c.Column)
	}
	if c := cmd(t, "alter table t drop column d.t.a"); c.Name != "a" {
		t.Errorf("DROP name = %q", c.Name)
	}
	if c := cmd(t, "alter table t alter t.a set default 1"); c.Name != "a" {
		t.Errorf("ALTER name = %q", c.Name)
	}
	if c := cmd(t, "alter table t add column a int after t.b"); c.After != "b" {
		t.Errorf("AFTER = %q", c.After)
	}
}

var nextValueColumnNameCases = []tableSampleCase{
	// NEXT VALUE FOR and NEXTVAL.
	{"alter table t add column a int default next value for s", true},
	{"alter table t alter a set default (next value for s)", true},
	{"alter table t alter a set default next value for s", false},
	{"alter table t alter a set default nextval(s)", false},
	{"alter table t modify a int default next value for s", true},
	{"create table next (a int)", true},
	{"create table t (a int check (a > next value for s))", true},
	{"create table t (a int default ((next value for s)))", true},
	{"create table t (a int default (next value for s))", true},
	{"create table t (a int default (nextval(s)))", true},
	{"create table t (a int default next value for a.s)", true},
	{"create table t (a int default next value for s on update now())", true},
	{"create table t (a int default next value for s primary key)", true},
	{"create table t (a int default next value for s)", true},
	{"create table t (a int default next value for)", false},
	{"create table t (a int default next value)", false},
	{"create table t (a int default nextval(s))", true},
	{"create table t (a int generated always as (next value for s))", true},
	{"create table t (a int primary key default next value for s, b int)", true},
	{"create table t (a timestamp on update next value for s)", false},
	{"create table t (a timestamp on update nextval(s))", false},
	{"create table t (next int, value int, `for` int)", true},
	{"create table t2 as select next value for s", true},
	{"create view v as select next value for s", true},
	{"delete from t where a = next value for s", true},
	{"explain select next value for s", true},
	{"insert into t (a) values (next value for s), (next value for s)", true},
	{"insert into t select next value for s", true},
	{"insert into t set a = next value for s", true},
	{"insert into t values (next value for s)", true},
	{"next value for s", false},
	{"select ((next value for a))", true},
	{"select (next value for a)", true},
	{"select * from t order by next value for s", true},
	{"select * from t where a = next value for s", true},
	{"select 1 + next value for a", true},
	{"select NEXT VALUE FOR a", true},
	{"select `next` from t", true},
	{"select `next` value from t", true},
	{"select count(next value for s) from t", true},
	{"select if(1, next value for a, 2)", true},
	{"select lastval(a)", true},
	{"select next", true},
	{"select next /* c */ value /* c */ for /* c */ a", true},
	{"select next VALUE For a", true},
	{"select next `value` from t", true},
	{"select next a", true},
	{"select next from t", true},
	{"select next value", false},
	{"select next value a", false},
	{"select next value for", false},
	{"select next value for 'a'", false},
	{"select next value for 1", false},
	{"select next value for `a`", true},
	{"select next value for `a`.`b`", true},
	{"select next value for a", true},
	{"select next value for a + 1", true},
	{"select next value for a . b", true},
	{"select next value for a . b . c", false},
	{"select next value for a . select", true},
	{"select next value for a as b", true},
	{"select next value for a b", true},
	{"select next value for a for update", true},
	{"select next value for a from t", true},
	{"select next value for a, 1", true},
	{"select next value for a.b", true},
	{"select next value for a.b from t", true},
	{"select next value for a.b.c", false},
	{"select next value for a.select", true},
	{"select next value for select", false},
	{"select next value for transaction", true},
	{"select next value from t", false},
	{"select next value(1)", false},
	{"select next value, a from t", false},
	{"select next, value from t", true},
	{"select next.a from t", true},
	{"select nextval", true},
	{"select nextval from t", true},
	{"select nextval('a')", false},
	{"select nextval()", false},
	{"select nextval(*)", false},
	{"select nextval(1 + 1)", false},
	{"select nextval(1)", false},
	{"select nextval(`a`)", true},
	{"select nextval(a + 1)", false},
	{"select nextval(a . select)", true},
	{"select nextval(a)", true},
	{"select nextval(a, b)", false},
	{"select nextval(a.*)", false},
	{"select nextval(a.b)", true},
	{"select nextval(a.b.c)", false},
	{"select nextval(distinct a)", false},
	{"select setval(a, 1)", true},
	{"select value from t", true},
	{"set @a = next value for s", true},
	{"set a = next value for s", true},
	{"update t set a = next value for s", true},

	// Qualified column names.
	{"alter table t add (t.a int)", true},
	{"alter table t add column (t.a int, t.b int)", true},
	{"alter table t add column a int after a.b.c.d", false},
	{"alter table t add column a int after b", true},
	{"alter table t add column a int after d.t.b", true},
	{"alter table t add column a int after t.b", true},
	{"alter table t add column a int first", true},
	{"alter table t add column a.b.c.d int", false},
	{"alter table t add column d.t.a int", true},
	{"alter table t add column if not exists t.a int", true},
	{"alter table t add column t.a int", true},
	{"alter table t add column x.a int", true},
	{"alter table t add index (t.a)", true},
	{"alter table t alter column t.a drop default", true},
	{"alter table t alter column t.a set default 1", true},
	{"alter table t alter t.a drop default", true},
	{"alter table t alter t.a set default 1", true},
	{"alter table t alter t.a set invisible", false},
	{"alter table t alter t.a set visible", false},
	{"alter table t change a b int after t.c", true},
	{"alter table t change a t.b int", true},
	{"alter table t change column d.t.a x.b int", true},
	{"alter table t change column if exists t.a b int", true},
	{"alter table t change t.a b int", true},
	{"alter table t change t.a t.b int", true},
	{"alter table t drop column if exists t.a", true},
	{"alter table t drop column t.a", true},
	{"alter table t drop t.a", true},
	{"alter table t modify a int after d.t.b", true},
	{"alter table t modify a int after t.b", true},
	{"alter table t modify column d.t.a int", true},
	{"alter table t modify column if exists t.a int", true},
	{"alter table t modify column x.a int", true},
	{"alter table t modify t.a int", true},
	{"alter table t rename column a to t.b", false},
	{"alter table t rename column t.a to b", false},
	{"create index i on t (d.t.a)", true},
	{"create index i on t (t.a)", true},
	{"create table d.t (d.t.a int)", true},
	{"create table d.t (d.x.a int)", true},
	{"create table d.t (t.a int)", true},
	{"create table d.t (x.t.a int)", true},
	{"create table t (.a int)", false},
	{"create table t (`t`.`a` int)", true},
	{"create table t (a . b int)", true},
	{"create table t (a . select int)", true},
	{"create table t (a .b int)", true},
	{"create table t (a int references p ())", false},
	{"create table t (a int references p (a.b.c.d))", false},
	{"create table t (a int references p (d.t.x))", true},
	{"create table t (a int references p (x, y.z))", true},
	{"create table t (a int references p (x.y))", true},
	{"create table t (a int, b int, foreign key (t.a) references p (x))", true},
	{"create table t (a int, foreign key (a) references p (x.y))", true},
	{"create table t (a int, index i (t.a))", true},
	{"create table t (a int, key ((t.a)))", true},
	{"create table t (a int, key (a.b.c.d))", false},
	{"create table t (a int, key (t.a desc))", true},
	{"create table t (a int, key (t.a(10)))", true},
	{"create table t (a int, key (t.a))", true},
	{"create table t (a int, primary key (d.t.a))", true},
	{"create table t (a int, unique (t.a))", true},
	{"create table t (a. b int)", true},
	{"create table t (a.b int)", true},
	{"create table t (a.b.c int)", true},
	{"create table t (a.b.c.d int)", false},
	{"create table t (a.global int)", true},
	{"create table t (a.select int)", true},
	{"create table t (b.a int)", true},
	{"create table t (d.t.a int)", true},
	{"create table t (global.a int)", true},
	{"create table t (select.a int)", true},
	{"create table t (t. int)", false},
	{"create table t (t.`a` int)", true},
	{"create table t (t.a int default 1 comment 'x')", true},
	{"create table t (t.a int primary key)", true},
	{"create table t (t.a int)", true},
	{"create table t (t.a int) select 1", true},
	{"create table t (t.a int, db.t.b int)", true},
	{"create table t (t.a int, primary key (a))", true},
	{"create table t (t.a int, primary key (t.a))", true},
	{"create table t (t.a)", false},
	{"create table t (t.t.t int)", true},
	{"create temporary table t (t.a int)", true},
}
