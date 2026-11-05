package parser

import (
	"sort"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// The words omni reserved after MySQL that TiDB v8.5.5 does not, every verdict measured with
// PREPARE on pingcap/tidb:v8.5.5 on 2026-11-05: a 1064 is the parser's refusal, and any other
// answer means the text parsed.
//
// keywordCategories followed MySQL, so 32 words were refused as names in every position
// (tidbUnreservedWords). Most are not keywords upstream at all, because TiDB has no LOOP,
// SIGNAL or GET DIAGNOSTICS. The function words had the same problem by another route: the
// parser reads count, now and the like as calls by token, and refused them without a `(`,
// though the table already called them unreserved.
//
// Not in this table, because omni still disagrees. Statements and forms TiDB does not have,
// which omni reads: `GET DIAGNOSTICS`, `PURGE`, `RESET PERSIST`, `RETURN` and `DECLARE UNDO
// HANDLER` in a body, `LATERAL (...)`, `JSON_TABLE(...)` and `START SLAVE UNTIL`. A word
// called as a function that TiDB will not call, as in `SELECT function(1)` and 381 others,
// because a call takes the bare identifier token upstream and omni takes any unreserved word.
// `SELECT now (1)`, `SELECT cast (a AS CHAR)` and `SELECT extract (day FROM a)`, because
// upstream's builtin tokens need the `(` to touch them. `SELECT DEFAULT FROM t`. A column
// DEFAULT beyond a bare name: `DEFAULT (b)`, `DEFAULT 1 + 1` and `DEFAULT NEXT VALUE FOR s`.

// reservedWordPositions are the positions the oracle sweep tries every keyword in.
var reservedWordPositions = []string{
	"select %s from t",
	"select %s.a from t",
	"select a.%s from t",
	"select a %s from t",
	"select a as %s from t",
	"select * from t %s",
	"select * from t as %s",
	"select * from %s",
	"select * from %s.a",
	"select * from a.%s",
	"create table %s (a int)",
	"create table t (%s int)",
	"update t set %s = 1",
	"insert into t (%s) values (1)",
}

// reservedWordSweepGap holds the one statement in the sweep where omni still disagrees. The
// oracle sweep fails when it starts to agree, so the entry cannot outlive the gap.
var reservedWordSweepGap = map[string]bool{"select default from t": true}

// registeredKeywords returns every word the lexer holds as a keyword, in order.
func registeredKeywords() []string {
	words := make([]string, 0, len(keywords))
	for word := range keywords {
		words = append(words, word)
	}
	sort.Strings(words)
	return words
}

// TestReservedWordForms holds omni to the recorded verdicts without a server, so a run with no
// Docker still reports a widened or narrowed grammar.
func TestReservedWordForms(t *testing.T) {
	for _, c := range reservedWordCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestTiDBUnreservedWordsAreClassified pins the table itself, which every name position reads.
func TestTiDBUnreservedWordsAreClassified(t *testing.T) {
	for word := range tidbUnreservedWords {
		if IsReservedKeyword(word) {
			t.Errorf("%q is reserved, and TiDB takes it as a name", word)
		}
	}
	for _, word := range []string{"array", "until", "select", "rank", "if", "row_number"} {
		if !IsReservedKeyword(word) {
			t.Errorf("%q is unreserved, and TiDB refuses it as a name", word)
		}
	}
}

// TestFunctionWordTrees pins what a function word turned into, because a verdict cannot tell a
// column from a call.
func TestFunctionWordTrees(t *testing.T) {
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
		{"select count from t", "", "count"},
		{"select count.a from t", "count", "a"},
		{"select now", "", "now"},
		{"select cast from t", "", "cast"},
		{"select extract.a from t", "extract", "a"},
		{"select variance from t", "", "variance"},
		{"select grouping from t", "", "grouping"},
		{"select return from t", "", "return"},
		{"select sql_no_cache.a from t", "sql_no_cache", "a"},
		{"select sql_buffer_result.a from t", "sql_buffer_result", "a"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			col, ok := first(t, c.sql).(*nodes.ColumnRef)
			if !ok {
				t.Fatalf("value is %T", first(t, c.sql))
			}
			if col.Table != c.table || col.Column != c.column {
				t.Errorf("table=%q column=%q, want %q %q", col.Table, col.Column, c.table, c.column)
			}
		})
	}
	for sql, want := range map[string]string{
		"select count(*) from t":                      "COUNT",
		"select count (a) from t":                     "COUNT",
		"select now()":                                "NOW",
		"select variance(a) from t":                   "VARIANCE",
		"select group_concat(a separator ',') from t": "GROUP_CONCAT",
	} {
		fn, ok := first(t, sql).(*nodes.FuncCallExpr)
		if !ok {
			t.Errorf("%q: value is %T", sql, first(t, sql))
			continue
		}
		if fn.Name != want {
			t.Errorf("%q: name = %q, want %q", sql, fn.Name, want)
		}
	}
	if _, ok := first(t, "select cast(a as char) from t").(*nodes.CastExpr); !ok {
		t.Errorf("CAST with `(` is %T", first(t, "select cast(a as char) from t"))
	}
	if _, ok := first(t, "select extract(day from a) from t").(*nodes.ExtractExpr); !ok {
		t.Errorf("EXTRACT with `(` is %T", first(t, "select extract(day from a) from t"))
	}

	// The select modifiers still read as modifiers where no dot follows.
	list, err := Parse("select sql_no_cache sql_buffer_result a from t")
	if err != nil {
		t.Fatal(err)
	}
	if s := list.Items[0].(*nodes.SelectStmt); !s.NoCache || !s.BufferResult {
		t.Errorf("modifiers lost: NoCache=%v BufferResult=%v", s.NoCache, s.BufferResult)
	}
}

var reservedWordCases = []tableSampleCase{
	// As names, in the positions that take one.
	{"alter table t add column a int after before", true},
	{"alter table t add column a int first", true},
	{"alter table t add column before int", true},
	{"create index before on t (a)", true},
	{"create table function (a int)", true},
	{"create table loop (a int)", true},
	{"create table return (signal int, resignal int, get int, purge int, undo int, declare int)", true},
	{"create table t (a int, constraint signal unique (a))", true},
	{"create table t (a int, index return (a))", true},
	{"create user declare", true},
	{"create view loop as select 1", true},
	{"deallocate prepare signal", true},
	{"delete from declare", true},
	{"execute signal", true},
	{"prepare signal from 'select 1'", true},
	{"rollback to loop", true},
	{"rollback to savepoint return", true},
	{"savepoint loop", true},
	{"select 'x' return", true},
	{"select (1) signal", true},
	{"select * from (select 1) return", true},
	{"select * from full", true},
	{"select * from json_table a", true},
	{"select * from json_table", true},
	{"select * from json_table.a", true},
	{"select * from lateral a", true},
	{"select * from lateral", true},
	{"select * from lateral.a", true},
	{"select * from t a join loop on a.b = loop.c", true},
	{"select * from t a join return using (b)", true},
	{"select * from t a straight_join signal on 1", true},
	{"select * from t for update of empty", true},
	{"select * from t for update return", false},
	{"select * from t group by loop", true},
	{"select * from t having get = 1", true},
	{"select * from t limit 1 loop", false},
	{"select * from t order by return", true},
	{"select * from t where a = 1 return", false},
	{"select * from t where each = 1", true},
	{"select * from t where function = 1", true},
	{"select * from t, lateral", true},
	{"select 1 + 1 get", true},
	{"select 1 from t tablesample system (1)", true},
	{"select 1 from t tablesample system", false},
	{"select 1 loop", true},
	{"select 1 signal, 2 resignal", true},
	{"select @a := return", true},
	{"select a as `return` from t", true},
	{"select a before from t", true},
	{"select a from t a full join b on 1", false},
	{"select a from t full join b", true},
	{"select a from t full outer join b", false},
	{"select a from t full", true},
	{"select a from t group by a with cube", false},
	{"select a from t group by cube(a)", true},
	{"select a full from t", true},
	{"select a full", true},
	{"select a full, b from t", true},
	{"select a function from t", true},
	{"select a get from t", true},
	{"select a loop from t", true},
	{"select a return from t", true},
	{"select a signal from t", true},
	{"select a system from t", true},
	{"select a.return from a", true},
	{"select accessible from t", true},
	{"select asensitive from t", true},
	{"select condition from t", true},
	{"select cube(1)", true},
	{"select cube(a) from t", true},
	{"select declare from t", true},
	{"select distinct sql_no_cache.a from t", true},
	{"select empty from t", true},
	{"select full from t", true},
	{"select function from t", true},
	{"select function.a from t", true},
	{"select get from t", true},
	{"select get", true},
	{"select get(1)", true},
	{"select insensitive from t", true},
	{"select io_after_gtids from t", true},
	{"select io_before_gtids from t", true},
	{"select json_table from t", true},
	{"select lateral from t", true},
	{"select loop(1)", true},
	{"select optimizer_costs from t", true},
	{"select read_write from t", true},
	{"select reads from t", true},
	{"select return(1)", true},
	{"select return.a from return", true},
	{"select sensitive from t", true},
	{"select signal(1)", true},
	{"select sql_buffer_result a from t", true},
	{"select sql_buffer_result from t", false},
	{"select sql_buffer_result.a from t", true},
	{"select sql_no_cache a from t", true},
	{"select sql_no_cache from t", false},
	{"select sql_no_cache sql_no_cache.a from t", true},
	{"select sql_no_cache", false},
	{"select sql_no_cache.a from t", true},
	{"select system from t", true},
	{"select system", true},
	{"update return set loop = 1", true},
	{"use system", true},
	{"with get as (select 1) select * from get", true},

	// A function word is a name unless `(` follows. Those TiDB reserves stay refused.
	{"alter table t modify a int after loop", true},
	{"create database get", true},
	{"create function f() returns int not deterministic modifies sql data return 1", false},
	{"create trigger before before insert on a for each row set @a = 1", false},
	{"create trigger each after insert on a for each row set @a = 1", false},
	{"create trigger t before insert on a for each row set @a = 1", false},
	{"drop function if exists return", false},
	{"insert into return values (1)", true},
	{"insert into t (return) values (1)", true},
	{"insert into t values (1) return", false},
	{"select * from t group by count", true},
	{"select * from t order by count", true},
	{"select * from t where count > 1", true},
	{"select 1 + count from t", true},
	{"select a count from t", true},
	{"select cast a from t", true},
	{"select cast from t", true},
	{"select cast(a as char) from t", true},
	{"select cast(a as signed) return from t", true},
	{"select cast.a from t", true},
	{"select char from t", false},
	{"select count (1) from t", true},
	{"select count (a) from t", true},
	{"select count + 1 from t", true},
	{"select count = 1 from t", true},
	{"select count a from t", true},
	{"select count count from t", true},
	{"select count from t", true},
	{"select count(*) from t", true},
	{"select count(*) loop from t", true},
	{"select count(a) from t", true},
	{"select count(return) from t", true},
	{"select count.a from t", true},
	{"select database from t", false},
	{"select extract from t", true},
	{"select extract(day from a) from t", true},
	{"select extract.a from t", true},
	{"select group_concat from t", true},
	{"select group_concat(a order by a separator ',') from t", true},
	{"select group_concat(a separator ',') from t", true},
	{"select group_concat(a) from t", true},
	{"select grouping from t", true},
	{"select grouping(a) from t group by a with rollup", true},
	{"select grouping.a from t", true},
	{"select if from t", false},
	{"select if(1,2,3)", true},
	{"select if(return, 1, 2)", true},
	{"select insert from t", false},
	{"select left from t", false},
	{"select max.max from max", true},
	{"select mod from t", false},
	{"select modifies from t", true},
	{"select now", true},
	{"select now() function", true},
	{"select now()", true},
	{"select now(), system", true},
	{"select now.a", true},
	{"select position from t", true},
	{"select position(a in b) from t", true},
	{"select rank from t", false},
	{"select rank() over ()", true},
	{"select rank.a from t", true},
	{"select repeat from t", false},
	{"select replace from t", false},
	{"select right from t", false},
	{"select row_number from t", false},
	{"select specific from t", true},
	{"select substring from t", true},
	{"select substring(a from 1) from t", true},
	{"select substring(a, 1) from t", true},
	{"select sum from t", true},
	{"select sum(a) from t", true},
	{"select sum(a) over (partition by loop) from t", true},
	{"select sum(a) over return from t", true},
	{"select sum(a) over w from t window return as (order by a)", true},
	{"select trim from t", true},
	{"select trim(a) from t", true},
	{"select trim(both 'a' from b) from t", true},
	{"select variance from t group by variance", true},
	{"select variance from t", true},
	{"select variance(a) from t", true},
	{"select variance.a from t", true},

	// The words where each one starts something, and as a name beside it.
	{"create function f() returns int deterministic reads sql data return 1", false},
	{"create function f() returns int return 1", false},
	{"create function function() returns int return 1", false},
	{"create function return() returns int return 1", false},
	{"create procedure p() begin a: loop select 1; end loop a; end", false},
	{"create procedure p() begin declare a condition for 1; end", false},
	{"create procedure p() begin declare a int; end", true},
	{"create procedure p() begin declare c asensitive cursor for select 1; end", false},
	{"create procedure p() begin declare c cursor for select 1; declare c2 insensitive cursor for select 1; end", false},
	{"create procedure p() begin declare c sensitive cursor for select 1; end", false},
	{"create procedure p() begin declare exit handler for sqlexception set @a = 1; end", true},
	{"create procedure p() begin declare loop int; end", true},
	{"create procedure p() begin loop = 1; end", false},
	{"create procedure p() begin loop select 1; end loop; end", false},
	{"create procedure p() begin loop: loop select 1; end loop; end", false},
	{"create procedure p() begin resignal; end", false},
	{"create procedure p() begin select loop; end", true},
	{"create procedure p() begin set loop = 1; end", true},
	{"create procedure p() begin signal sqlstate '45000'; end", false},
	{"declare a int", false},
	{"declare c insensitive cursor for select 1", false},
	{"declare", false},
	{"drop function function", false},
	{"drop function return", false},
	{"grant execute on function f to u", true},
	{"grant execute on function function to u", true},
	{"grant select on function a to u", true},
	{"loop", false},
	{"resignal sqlstate '45000'", false},
	{"resignal", false},
	{"return 1", false},
	{"return", false},
	{"revoke execute on function f from u", true},
	{"set read_write = 1", true},
	{"set transaction read write", true},
	{"show create function f", false},
	{"show create function function", false},
	{"show function code f", false},
	{"show function status", true},
	{"signal = 1", false},
	{"signal sqlstate '45000'", false},
	{"signal", false},
	{"start transaction read write", true},

	// A column's DEFAULT and ON UPDATE take no bare name, DefaultValueExpr, parser.y:4052.
	{"alter table t add column a int default b", false},
	{"alter table t add column a timestamp default now", false},
	{"alter table t alter a set default (now())", true},
	{"alter table t change a b int default c", false},
	{"alter table t modify a int default b", false},
	{"create table t (a int default 'x')", true},
	{"create table t (a int default (now()))", true},
	{"create table t (a int default -1)", true},
	{"create table t (a int default 1)", true},
	{"create table t (a int default `b`)", false},
	{"create table t (a int default b)", false},
	{"create table t (a int default b.c)", false},
	{"create table t (a int default count)", false},
	{"create table t (a int default current_timestamp())", true},
	{"create table t (a int default now())", true},
	{"create table t (a int default now(1))", true},
	{"create table t (a int default now)", false},
	{"create table t (a int default null)", true},
	{"create table t (a int default true)", true},
	{"create table t (a int references b (c) on update cascade)", true},
	{"create table t (a int, foreign key (a) references b (c) on update cascade on delete set null)", true},
	{"create table t (a timestamp default now() on update now())", true},
	{"create table t (a timestamp default now() on update now)", false},
	{"create table t (a timestamp default now)", false},
	{"create table t (a timestamp on update b)", false},
	{"create table t (a timestamp on update current_timestamp())", true},
	{"create table t (a timestamp on update now())", true},
	{"create table t (a timestamp on update now)", false},
}
