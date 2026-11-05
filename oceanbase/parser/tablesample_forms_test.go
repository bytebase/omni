package parser

import "testing"

// OceanBase 4.3.5 has no TABLESAMPLE. Every verdict was measured with PREPARE on
// oceanbase/oceanbase-ce:4.3.5-lts (4.3.5.6) on 2026-11-05: a 1064 is the parser's
// refusal, and any other answer means the text parsed. Every form of the clause is a
// 1064 and the word is an ordinary name, so `t tablesample` reads an alias and the
// parentheses after it are the error. Its own sampling clause is SAMPLE, not this.
//
// Not in this table, because they do not depend on this word: `LOCK TABLES t
// tablesample READ` is a 1064 as `LOCK TABLES t a READ` is, `SELECT @@ tablesample` as
// `SELECT @@ a` is, and `SELECT @@tablesample.a` as `SELECT @@a.b` is. omni accepts them.

// tableSampleCase is a statement and whether the engine's parser accepts it.
type tableSampleCase struct {
	sql    string
	accept bool
}

// parsesAsOneStatement is the acceptance these tables mean. Parse needs no
// separator between statements, so a clause the grammar refuses could otherwise
// pass by being read as a second statement.
func parsesAsOneStatement(sql string) (bool, error) {
	list, err := Parse(sql)
	if err != nil {
		return false, err
	}
	return len(list.Items) == 1, nil
}

// TestTableSampleForms holds omni to the recorded verdicts without a server, so a
// run with no Docker still reports a widened or narrowed grammar.
func TestTableSampleForms(t *testing.T) {
	for _, c := range tableSampleCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

var tableSampleCases = []tableSampleCase{
	// The clause.
	{"select * from tbl tablesample system (50)", false},
	{"select * from tbl tablesample bernoulli (50)", false},
	{"select * from tbl tablesample regions ()", false},
	{"select * from tbl tablesample regions (50)", false},
	{"select * from tbl tablesample ()", false},
	{"select * from tbl tablesample system ()", false},
	{"select * from tbl tablesample bernoulli ()", false},
	{"select * from tbl tablesample (50)", false},
	{"select * from tbl tablesample (50 percent)", false},
	{"select * from tbl tablesample (50 rows)", false},
	{"select * from tbl tablesample (50 percent rows)", false},
	{"select * from tbl tablesample (50 rows percent)", false},
	{"select * from tbl tablesample (percent)", false},
	{"select * from tbl tablesample (percent percent)", false},
	{"select * from tbl tablesample (rows)", false},
	{"select * from tbl tablesample (1+1)", false},
	{"select * from tbl tablesample (a)", false},
	{"select * from tbl tablesample (a + b rows)", false},
	{"select * from tbl tablesample ((select 1))", false},
	{"select * from tbl tablesample ('34')", false},
	{"select * from tbl tablesample ('34' rows)", false},
	{"select * from tbl tablesample (50, 50)", false},
	{"select * from tbl tablesample (50) repeatable (1)", false},
	{"select * from tbl tablesample (50) repeatable (1+1)", false},
	{"select * from tbl tablesample (50) repeatable ('a')", false},
	{"select * from tbl tablesample (50) repeatable ()", false},
	{"select * from tbl tablesample (50) repeatable", false},
	{"select * from tbl tablesample () repeatable (1)", false},
	{"select * from tbl tablesample regions () repeatable (1)", false},
	{"select * from tbl tablesample (50) repeatable (1) repeatable (2)", false},
	{"select * from tbl tablesample (50) tablesample (50)", false},
	{"select * from tbl tablesample (50) a", false},
	{"select * from tbl tablesample (50) as a", false},
	{"select * from tbl tablesample", true},
	{"select * from tbl tablesample system", false},
	{"select * from tbl tablesample (", false},
	{"select * from tbl tablesample (50", false},
	{"select * from tbl tablesample foo (50)", false},
	{"select * from tbl tablesample 50", false},
	{"select * from tbl tablesample system 50", false},

	// Where the clause may stand.
	{"select * from db.tbl tablesample (50)", false},
	{"select * from tbl a tablesample (50)", false},
	{"select * from tbl as a tablesample (50)", false},
	{"select * from tbl `tablesample` tablesample (50)", false},
	{"select * from tbl partition (p0) tablesample (50)", false},
	{"select * from tbl partition (p0) a tablesample (50)", false},
	{"select * from tbl tablesample (50) partition (p0)", false},
	{"select * from tbl use index (i) tablesample (50)", false},
	{"select * from tbl a use index (i) tablesample (50)", false},
	{"select * from tbl tablesample (50) use index (i)", false},
	{"select * from tbl partition (p0) a use index (i) tablesample (50) repeatable (1)", false},
	{"select * from tbl as of timestamp now() tablesample (50)", false},
	{"select * from (tbl tablesample (50))", false},
	{"select * from (tbl tablesample (50), t2)", false},
	{"select * from (select 1) a tablesample (50)", false},
	{"select * from (t1 join t2) tablesample (50)", false},
	{"select * from dual tablesample (50)", false},
	{"select * from json_table('[]', '$' columns (a int path '$')) a tablesample (50)", false},
	{"select * from t1 tablesample (10), t2 tablesample (20)", false},
	{"select * from t1 a tablesample (10) join t2 b tablesample (20) on a.id = b.id", false},
	{"select * from t1 tablesample (10) left join t2 tablesample (20) on t1.id = t2.id", false},
	{"select * from t1 tablesample (10) straight_join t2 tablesample (20)", false},
	{"select * from t1 tablesample (10) natural join t2 tablesample (20)", false},
	{"select * from t1 tablesample (10) cross join t2 tablesample (20)", false},
	{"select * from t1 join t2 tablesample (20) using (id)", false},
	{"select * from tbl where id > 20 tablesample (50)", false},
	{"select * from tbl tablesample (50) where id > 20", false},
	{"select * from tbl tablesample (50) group by a", false},
	{"select * from tbl tablesample (50) having a > 1", false},
	{"select * from tbl tablesample (50) order by a", false},
	{"select * from tbl tablesample (50) limit 1", false},
	{"select * from tbl tablesample (50) for update", false},
	{"select * from tbl tablesample (50) union select * from t2 tablesample (50)", false},
	{"select * from tbl tablesample bernoulli(50) limit 1 into outfile '/tmp/sample.csv'", false},
	{"select * from tbl tablesample bernoulli(50) order by a, b into outfile '/tmp/sample.csv'", false},
	{"select (select 1 from tbl tablesample (50))", false},
	{"with c as (select * from tbl tablesample (50)) select * from c", false},
	{"with c as (select 1) select * from c tablesample (50)", false},

	// Other statements that read table references.
	{"insert into t select * from tbl tablesample (50)", false},
	{"create table t2 as select * from tbl tablesample (50)", false},
	{"create view v as select * from tbl tablesample (50)", false},
	{"replace into t select * from tbl tablesample (50)", false},
	{"explain select * from tbl tablesample (50)", false},
	{"update tbl tablesample (50) set a = 1", false},
	{"update tbl a tablesample (50) set a = 1", false},
	{"update tbl tablesample (50) join t2 on 1 = 1 set a = 1", false},
	{"delete from tbl tablesample (50)", false},
	{"delete from tbl a tablesample (50) where 1 = 1", false},
	{"delete t from tbl t tablesample (50)", false},
	{"delete from t using tbl t tablesample (50)", false},
	{"delete t from t1 t join t2 tablesample (50)", false},
	{"lock tables tbl tablesample (50) read", false},
	{"lock tables tbl read tablesample (50)", false},
	{"table tbl tablesample (50)", false},
	{"insert into tbl tablesample (50) values (1)", false},
	{"analyze table tbl tablesample (50)", false},
	{"admin check table tbl tablesample (50)", false},

	// The unquoted word as a name.
	{"select * from tablesample", true},
	{"select * from tablesample a", true},
	{"select * from db.tablesample", true},
	{"select tablesample.a from tablesample", true},
	{"select tablesample from t", true},
	{"select tablesample.a from t", true},
	{"select t.tablesample from t", true},
	{"select db.t.tablesample from t", true},
	{"select tablesample, a from t", true},
	{"select a from t where tablesample = 1", true},
	{"select a from t order by tablesample", true},
	{"select a from t group by tablesample", true},
	{"select tablesample()", true},
	{"select tablesample(1)", true},
	{"select 1 tablesample", true},
	{"select 1 as tablesample", true},
	{"select * from t tablesample", true},
	{"select * from t as tablesample", true},
	{"select * from t tablesample where a = 1", true},
	{"select * from t tablesample, t2", true},
	{"select * from t tablesample join t2", true},
	{"select * from t tablesample order by 1", true},
	{"select * from t tablesample limit 1", true},
	{"select * from t tablesample group by a", true},
	{"select * from t tablesample use index (i)", true},
	{"select * from t partition (p0) tablesample", true},
	{"select * from (select 1) tablesample", true},
	{"select * from (select 1) as tablesample", true},
	{"select * from t1 join t2 tablesample on 1 = 1", true},
	{"with tablesample as (select 1) select * from tablesample", true},
	{"update tablesample set a = 1", true},
	{"update t tablesample set a = 1", true},
	{"update t as tablesample set a = 1", true},
	{"update t set tablesample = 1", true},
	{"delete from tablesample", true},
	{"delete from db.tablesample", true},
	{"delete tablesample from t", true},
	{"delete tablesample from tablesample", true},
	{"delete from t tablesample", true},
	{"delete from t tablesample where a = 1", true},
	{"delete from t as tablesample where a = 1", true},
	{"delete from t tablesample order by a", true},
	{"delete from t tablesample limit 1", true},
	{"delete from t tablesample using t", false},
	{"delete tablesample.* from t tablesample", true},
	{"delete t.* from t tablesample", true},
	{"insert into tablesample values (1)", true},
	{"insert into t (tablesample) values (1)", true},
	{"create table tablesample (a int)", true},
	{"create table t (tablesample int)", true},
	{"create table db.tablesample (a int)", true},
	{"create index tablesample on t (a)", true},
	{"create index i on t (tablesample)", true},
	{"create database tablesample", true},
	{"create view tablesample as select 1", true},
	{"drop table tablesample", true},
	{"alter table t add column tablesample int", true},
	{"alter table t rename to tablesample", true},
	{"alter table tablesample add column a int", true},
	{"lock tables tablesample read", true},
	{"lock tables t as tablesample read", true},
	{"truncate table tablesample", true},
	{"select @tablesample", true},
	{"select @@tablesample", true},
	{"select 1 from t where a in (select tablesample from t2)", true},
	{"select * from t tablesample force index (i)", true},
	{"select * from t as tablesample tablesample (50)", false},
	{"select * from tablesample tablesample (50)", false},
	{"select tablesample from tablesample tablesample tablesample (50)", false},

	// More of the clause and the word.
	{"select * from tablesample.t", true},
	{"select * from tablesample.tablesample", true},
	{"select tablesample.t.a from t", true},
	{"select * from tbl tablesample (select 1)", false},
	{"select * from tbl tablesample `system` (50)", false},
	{"select * from tbl tablesample `bernoulli` (50)", false},
	{"select * from tbl tablesample `regions` ()", false},
	{"select * from tbl tablesample (50 `percent`)", false},
	{"select * from tbl tablesample (50 `rows`)", false},
	{"select * from tbl tablesample (50) `repeatable` (1)", false},
	{"select * from tbl tablesample (@a := 1)", false},
	{"select * from tbl tablesample (50) repeatable (@a := 1)", false},
	{"select * from tbl TABLESAMPLE SYSTEM (50 PERCENT) REPEATABLE (1)", false},
	{"select * from tbl tablesample region ()", false},
	{"select * from tbl tablesample (50) (select 1)", false},
	{"select * from t1 tablesample, t2 tablesample", true},
	{"select * from tbl a tablesample", false},
	{"select 1 from t where tablesample", true},
	{"select count(tablesample) from t", true},
	{"select `tablesample` from `tablesample` `tablesample`", true},
	{"delete from t `tablesample` where a = 1", true},
	{"delete `tablesample` from t `tablesample`", true},
	{"select * from t1 tablesample (10) join t2 tablesample (20) join t3 tablesample (30)", false},
	{"select * from (t1 tablesample (10) join t2) join t3 tablesample (30) on 1 = 1", false},
	{"select * from t1 join (t2 tablesample (20)) on 1 = 1", false},
	{"select @@global.tablesample", true},
	{"select @@session.tablesample", true},
	{"select tablesample .a from t", true},
	{"select * from tablesample .t", true},
	{"select tablesample. a from t", true},
	{"select t . tablesample from t", true},
	{"select t. tablesample from t", true},
	{"select t .tablesample from t", true},
	{"select * from db . tablesample", true},
	{"select * from db .tablesample", true},
	{"create table tablesample.t (a int)", true},
	{"drop table tablesample.t", true},
	{"delete from tablesample.t", true},
	{"update tablesample.t set a = 1", true},
	{"insert into tablesample.t values (1)", true},
	{"delete tablesample.t from tablesample.t", true},
	{"select tablesample.* from t", true},
	{"select tablesample . * from t", true},
	{"set @@tablesample = 1", true},
	{"set @@session.tablesample = 1", true},
	{"set session tablesample = 1", true},
	{"set tablesample = 1", true},
	{"set global tablesample = 1", true},
	{"select 1. tablesample from t", true},
	{"select 1.tablesample from t", true},
	{"select 1. tablesample", true},
	{"select 1.5 tablesample", true},
	{"select .5 tablesample", true},
	{"select t.a. tablesample from t", true},
	{"select a from t  .  tablesample", true},
	{"select tablesample\t.a from t", true},
	{"select t.\ttablesample from t", true},

	// Whitespace around the dot.
	{"select tablesample\t.a from t", true},
	{"select t.\ttablesample from t", true},
	{"select t.\ntablesample from t", true},
	{"select tablesample\n.a from t", true},
	{"select t./**/tablesample from t", true},
	{"select tablesample/**/.a from t", true},
	{"select tablesample  .a from t", true},
	{"select t  .  tablesample from t", true},
	{"select * from db  .  tablesample", true},
}
