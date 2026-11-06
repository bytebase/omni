package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// The words omni reserved after MySQL that MariaDB 11.8.8 does not, every verdict measured with
// PREPARE on mariadb:11.8.8 on 2026-11-05: a 1064 is the parser's refusal, and any other answer
// means the text parsed.
//
// keywordCategories followed MySQL, so 33 words were refused as names in every position
// (mariadbUnreservedWords). The function words had the same problem by another route: the parser
// reads count, now and the like as calls by token, and refused them without a `(`, though the table
// already called them unreserved.
//
// WINDOW is the one with a rule of its own: it names a table, a column or a select item, and it is
// no table's alias, because a WINDOW clause can follow one. `SELECT a window FROM t` is an alias,
// and `SELECT * FROM t window` is a 1064.
//
// Labels are any name but a short list, which keyword_sp_decl and a few others make, and
// mariadbNotLabelWords marks those as ambiguous_2. A label comes ahead of what its word starts, so
// `after: LOOP` is a label. A column's ON UPDATE takes no bare name.
//
// Not in this table, because omni still disagrees. A word called as a function that MariaDB will not
// call, as in `SELECT function(1)`, because omni takes any unreserved word there. `SELECT cast (a AS
// CHAR)` and `SELECT extract (day FROM a)`, whose `(` has to touch. `SELECT DEFAULT FROM t`. Forms omni reads and MariaDB does not: `LATERAL (...)`, `DECLARE UNDO HANDLER`, `GRANT
// OPTION ON`, `CHANGE MASTER TO MASTER_BIND`, `OVER (GROUPS ...)`. Forms MariaDB reads and omni does
// not: `DECLARE c SENSITIVE CURSOR`, `WITH CUBE`, `JSON_TABLE` with `ON EMPTY`, and `FOR SYSTEM_TIME AS
// OF TIMESTAMP now()`.

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
	"create table %s (a int)",
	"create table t (%s int)",
	"update t set %s = 1",
	"insert into t (%s) values (1)",
	"create procedure p() %s: loop leave %s; end loop",
}

// reservedWordSweepGap holds the statement in the sweep where omni still disagrees. The oracle
// sweep fails when it starts to agree, so the entry cannot outlive the gap.
var reservedWordSweepGap = map[string]bool{
	"select default from t": true,
	// A 1314 is the parser's action refusing the statement before it reads on, so it says nothing
	// of the label.
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

// TestMariaDBUnreservedWordsAreClassified pins the table itself, which every name position reads.
func TestMariaDBUnreservedWordsAreClassified(t *testing.T) {
	for word := range mariadbUnreservedWords {
		if IsReservedKeyword(word) {
			t.Errorf("%q is reserved, and MariaDB takes it as a name", word)
		}
	}
	for _, word := range []string{"select", "from", "if", "offset", "returning"} {
		if !IsReservedKeyword(word) {
			t.Errorf("%q is unreserved, and MariaDB refuses it as a name", word)
		}
	}
}

// TestFunctionWordTrees pins what a word turned into, because a verdict cannot tell a column from a
// call, or a window clause from an alias.
func TestFunctionWordTrees(t *testing.T) {
	selectOf := func(t *testing.T, sql string) *nodes.SelectStmt {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		return list.Items[0].(*nodes.SelectStmt)
	}
	first := func(t *testing.T, sql string) (nodes.ExprNode, string) {
		t.Helper()
		item := selectOf(t, sql).TargetList[0]
		if res, ok := item.(*nodes.ResTarget); ok {
			return res.Val, res.Name
		}
		return item, ""
	}
	for _, c := range []struct{ sql, column string }{
		{"select count from t", "count"},
		{"select now", "now"},
		{"select cast from t", "cast"},
		{"select extract from t", "extract"},
		{"select rank from t", "rank"},
		{"select variance from t", "variance"},
		{"select function from t", "function"},
		{"select window from t", "window"},
	} {
		val, _ := first(t, c.sql)
		if col, ok := val.(*nodes.ColumnRef); !ok || col.Column != c.column {
			t.Errorf("%q: value is %#v", c.sql, val)
		}
	}
	for sql, want := range map[string]string{
		"select count(*) from t":    "COUNT",
		"select now()":              "NOW",
		"select variance(a) from t": "VARIANCE",
		"select rank() over ()":     "RANK",
	} {
		val, _ := first(t, sql)
		if fn, ok := val.(*nodes.FuncCallExpr); !ok || fn.Name != want {
			t.Errorf("%q: value is %#v", sql, val)
		}
	}
	if val, _ := first(t, "select cast(a as char) from t"); val == nil {
		t.Error("CAST with `(` read nothing")
	} else if _, ok := val.(*nodes.CastExpr); !ok {
		t.Errorf("CAST with `(` is %T", val)
	}
	if val, _ := first(t, "select extract(day from a) from t"); val == nil {
		t.Error("EXTRACT with `(` read nothing")
	} else if _, ok := val.(*nodes.ExtractExpr); !ok {
		t.Errorf("EXTRACT with `(` is %T", val)
	}

	// WINDOW: an alias on a select item, and a clause after a table.
	if _, alias := first(t, "select a window from t"); alias != "window" {
		t.Errorf("select item alias = %q, want window", alias)
	}
	stmt := selectOf(t, "select sum(a) over w from t window w as (partition by b)")
	if len(stmt.WindowClause) != 1 {
		t.Errorf("WINDOW clause holds %d definitions, want 1", len(stmt.WindowClause))
	}
	if ref, ok := stmt.From[0].(*nodes.TableRef); !ok || ref.Alias != "" {
		t.Errorf("the table took an alias: %#v", stmt.From[0])
	}
}

var reservedWordCases = []asOfCase{
	// Labels, which any name but a short list takes, ahead of what the word starts.
	{"create procedure p() after: loop leave after; end loop", true},
	{"create procedure p() absent: loop leave absent; end loop", true},
	{"create procedure p() begin after: loop leave after; end loop; end", true},
	{"create procedure p() database: loop leave database; end loop", true},
	{"create procedure p() clone: loop leave clone; end loop", true},
	{"create procedure p() format: loop leave format; end loop", false},
	{"create procedure p() close: loop leave close; end loop", false},
	{"create procedure p() get: loop leave get; end loop", false},
	{"create procedure p() window: loop leave window; end loop", false},
	{"create procedure p() select: loop leave select; end loop", false},
	{"create procedure p() begin declare c cursor for select 1; open c; close c; end", true},
	{"create procedure p() l: begin select 1; end l", true},

	// A column's ON UPDATE takes no bare name, though its DEFAULT does, as an expression.
	{"create table x (a timestamp on update now)", false},
	{"create table x (a int on update b)", false},
	{"create table x (a timestamp on update now())", true},
	{"create table x (a timestamp on update current_timestamp)", true},
	{"create table x (a timestamp default now())", true},
	{"create table x (a timestamp default current_timestamp)", true},

	{"select a into from t", false},
	{"select * from t into", false},
	{"select a into @b, from t", false},
	{"select a into @b, @c from t", true},
	{"select a from t into @b", true},
	{"update t for portion of p from '2020-01-01' to '2021-01-01' a set b = 1", true},
	{"update t for portion of p from '2020-01-01' to '2021-01-01' as a set b = 1", true},
	{"update t for portion of p from '2020-01-01' to '2021-01-01' window set b = 1", false},
	{"update t for portion of p from '2020-01-01' to '2021-01-01' as window set b = 1", false},
	{"alter schema d character set utf8", true},
	{"alter table t add column a int after before", false},
	{"alter table t add column a int first", true},
	{"alter table t add column a int generated always as (1) stored", true},
	{"alter table t add column before int", false},
	{"alter table t add system versioning", true},
	{"alter table t drop system versioning", true},
	{"alter table t modify a int after loop", false},
	{"change master to master_ssl_verify_server_cert = 1", true},
	{"create database d", true},
	{"create database get", true},
	{"create function f() returns int deterministic reads sql data return 1", true},
	{"create function f() returns int not deterministic modifies sql data return 1", true},
	{"create function f() returns int return 1", true},
	{"create function function() returns int return 1", true},
	{"create function return() returns int return 1", false},
	{"create index before on t (a)", false},
	{"create procedure p() begin a: loop select 1; end loop a; end", true},
	{"create procedure p() begin declare a condition for 1; end", true},
	{"create procedure p() begin declare a int; end", true},
	{"create procedure p() begin declare exit handler for sqlexception set @a = 1; end", true},
	{"create procedure p() begin declare loop int; end", false},
	{"create procedure p() begin get diagnostics @a = number; end", true},
	{"create procedure p() begin loop = 1; end", false},
	{"create procedure p() begin loop select 1; end loop; end", true},
	{"create procedure p() begin loop: loop select 1; end loop; end", false},
	{"create procedure p() begin resignal; end", true},
	{"create function f() returns int begin return 1; end", true},
	{"create procedure p() begin return 1; end", false},
	{"create procedure p() begin select loop; end", false},
	{"create procedure p() begin set loop = 1; end", false},
	{"create procedure p() begin signal sqlstate '45000'; end", true},
	{"create schema d", true},
	{"create table function (a int)", true},
	{"create table generated (a int)", true},
	{"create table loop (a int)", false},
	{"create table option (a int)", true},
	{"create table return (signal int, resignal int, get int, purge int, undo int, declare int)", false},
	{"create table schema (a int)", true},
	{"create table system (a int)", true},
	{"create table t (a int as (1) stored)", true},
	{"create table t (a int as (1) virtual)", true},
	{"create table t (a int as (1))", true},
	{"create table t (a int generated always as (1) stored)", true},
	{"create table t (a int generated always as (1) virtual)", true},
	{"create table t (a int generated always as (1))", true},
	{"create table t (a int generated)", false},
	{"create table t (a int without system versioning)", true},
	{"create table t (a int) with system versioning", true},
	{"create table t (a int) without system versioning", false},
	{"create table t (a int, constraint signal unique (a))", false},
	{"create table t (a int, index return (a))", false},
	{"create table t (generated int, stored int, virtual int)", true},
	{"create table t (window int)", true},
	{"create trigger before before insert on a for each row set @a = 1", false},
	{"create trigger each after insert on a for each row set @a = 1", false},
	{"create trigger t before insert on a for each row set @a = 1", true},
	{"create user declare", false},
	{"create view loop as select 1", false},
	{"deallocate prepare signal", false},
	{"delete from declare", false},
	{"delete from t as window", false},
	{"delete from t window", false},
	{"delete from window", true},
	{"drop database d", true},
	{"drop function f", true},
	{"drop function function", true},
	{"drop function if exists f", true},
	{"drop function if exists return", false},
	{"drop function return", false},
	{"drop schema d", true},
	{"execute signal", false},
	{"get diagnostics @a = number", true},
	{"grant execute on function f to u", true},
	{"grant execute on function function to u", true},
	{"grant select on function a to u", true},
	{"grant select on t to u with grant option", true},
	{"insert into return values (1)", false},
	{"insert into t (return) values (1)", false},
	{"insert into t values (1) return", false},
	{"insert into window values (1)", true},
	{"prepare signal from 'select 1'", false},
	{"revoke execute on function f from u", true},
	{"rollback to loop", false},
	{"rollback to savepoint return", false},
	{"savepoint loop", false},
	{"select 'x' return", false},
	{"select (1) signal", false},
	{"select * from (select 1) a window w as ()", true},
	{"select * from (select 1) as window", false},
	{"select * from (select 1) return", false},
	{"select * from (select 1) window", false},
	{"select * from a.window", true},
	{"select * from full", true},
	{"select * from json_table", true},
	{"select * from json_table a", true},
	{"select * from json_table('[]', '$' columns (a int path '$')) t", true},
	{"select * from json_table.a", true},
	{"select * from lateral", true},
	{"select * from lateral a", true},
	{"select * from lateral.a", true},
	{"select * from t a join loop on a.b = loop.c", false},
	{"select * from t a join return using (b)", false},
	{"select * from t a straight_join signal on 1", false},
	{"select * from t a window w as ()", true},
	{"select * from t as rank", true},
	{"select * from t as window", false},
	{"select * from t for system_time all", true},
	{"select * from t for system_time as of 1", true},
	{"select * from t for update of empty", false},
	{"select * from t for update of of", false},
	{"select * from t for update of t", false},
	{"select * from t for update return", false},
	{"select * from t full", true},
	{"select * from t full join b", true},
	{"select * from t function", true},
	{"select * from t get", true},
	{"select * from t group by count", true},
	{"select * from t group by loop", false},
	{"select * from t groups", true},
	{"select * from t having get = 1", true},
	{"select * from t lateral", true},
	{"select * from t limit 1 loop", false},
	{"select * from t of", true},
	{"select * from t order by count", true},
	{"select * from t order by return", false},
	{"select * from t rank", true},
	{"select * from t system", true},
	{"select * from t where a = 1 return", false},
	{"select * from t where count > 1", true},
	{"select * from t where each = 1", false},
	{"select * from t where function = 1", true},
	{"select * from t window", false},
	{"select * from t window w as ()", true},
	{"select * from t, lateral", true},
	{"select * from window", true},
	{"select * from window w", true},
	{"select * from window, t", true},
	{"select * from window, window", false},
	{"select 1 + 1 get", true},
	{"select 1 + count from t", true},
	{"select 1 from t tablesample system", false},
	{"select 1 from t tablesample system (1)", false},
	{"select 1 loop", false},
	{"select 1 signal, 2 resignal", false},
	{"select 1 window w as ()", false},
	{"select @a := return", false},
	{"select a as `return` from t", true},
	{"select a as window from t", true},
	{"select a before from t", false},
	{"select a count from t", true},
	{"select a from t a full join b on 1", false},
	{"select a from t full", true},
	{"select a from t full join b", true},
	{"select a from t full outer join b", false},
	{"select a from t group by cube(a)", true},
	{"select a from t into @b", true},
	{"select a from t join b as window", false},
	{"select a from t join b window", false},
	{"select a from t window w as ()", true},
	{"select a full", true},
	{"select a full from t", true},
	{"select a full, b from t", true},
	{"select a function from t", true},
	{"select a get from t", true},
	{"select a into @b from t", true},
	{"select a loop from t", false},
	{"select a rank from t", true},
	{"select a return from t", false},
	{"select a signal from t", false},
	{"select a system from t", true},
	{"select a window from t", true},
	{"select a window w as ()", false},
	{"select a.return from a", true},
	{"select accessible from t", false},
	{"select asensitive from t", false},
	{"select cast a from t", true},
	{"select cast from t", true},
	{"select cast(a as char) from t", true},
	{"select cast(a as signed) return from t", false},
	{"select cast.a from t", true},
	{"select char from t", false},
	{"select condition from t", false},
	{"select count (1) from t", true},
	{"select count (a) from t", true},
	{"select count + 1 from t", true},
	{"select count = 1 from t", true},
	{"select count a from t", true},
	{"select count count from t", true},
	{"select count from t", true},
	{"select count(*) from t", true},
	{"select count(*) loop from t", false},
	{"select count(a) from t", true},
	{"select count(return) from t", false},
	{"select count.a from t", true},
	{"select cube from t", true},
	{"select cube(1)", true},
	{"select cube(a) from t", true},
	{"select cume_dist() over () from t", true},
	{"select database", true},
	{"select database from t", true},
	{"select database()", true},
	{"select declare from t", false},
	{"select dense_rank() over (order by a) from t", true},
	{"select distinct sql_no_cache.a from t", true},
	{"select empty from t", true},
	{"select extract from t", true},
	{"select extract(day from a) from t", true},
	{"select extract.a from t", true},
	{"select first_value(a) over () from t", true},
	{"select full from t", true},
	{"select function from t", true},
	{"select function.a from t", true},
	{"select generated from t", true},
	{"select get", true},
	{"select get from t", true},
	{"select group_concat from t", true},
	{"select group_concat(a order by a separator ',') from t", true},
	{"select group_concat(a separator ',') from t", true},
	{"select group_concat(a) from t", true},
	{"select grouping from t", true},
	{"select grouping(a) from t group by a with rollup", true},
	{"select grouping.a from t", true},
	{"select groups from t", true},
	{"select if from t", false},
	{"select if(1,2,3)", true},
	{"select if(return, 1, 2)", false},
	{"select insensitive from t", false},
	{"select insert from t", false},
	{"select into from t", false},
	{"select io_after_gtids from t", true},
	{"select io_before_gtids from t", true},
	{"select json_table from t", true},
	{"select lag(a) over () from t", true},
	{"select lateral from t", true},
	{"select lead(a, 1) over (order by a) from t", true},
	{"select left from t", false},
	{"select loop(1)", false},
	{"select master_bind from t", true},
	{"select master_ssl_verify_server_cert from t", true},
	{"select max from t", true},
	{"select max.max from max", true},
	{"select min(a) from t", true},
	{"select mod from t", false},
	{"select modifies from t", false},
	{"select now", true},
	{"select now (1)", true},
	{"select now (6) from t", true},
	{"select now()", true},
	{"select now() function", true},
	{"select now(), system", true},
	{"select now.a", true},
	{"select nth_value(a, 1) over () from t", true},
	{"select ntile(2) over () from t", true},
	{"select of from t", true},
	{"select optimizer_costs from t", true},
	{"select option from t", true},
	{"select percent_rank() over () from t", true},
	{"select position from t", true},
	{"select position(a in b) from t", true},
	{"select rank from t", true},
	{"select rank() over ()", true},
	{"select rank() over () from t", true},
	{"select rank.a from t", true},
	{"select read_write from t", false},
	{"select reads from t", false},
	{"select repeat from t", false},
	{"select replace from t", false},
	{"select return(1)", false},
	{"select return.a from return", false},
	{"select right from t", false},
	{"select row_number from t", false},
	{"select row_number() over () from t", true},
	{"select schema", true},
	{"select schema()", true},
	{"select sensitive from t", false},
	{"select signal(1)", false},
	{"select specific from t", false},
	{"select sql_buffer_result a from t", true},
	{"select sql_buffer_result from t", false},
	{"select sql_buffer_result.a from t", true},
	{"select sql_no_cache", false},
	{"select sql_no_cache a from t", true},
	{"select sql_no_cache from t", false},
	{"select sql_no_cache sql_no_cache.a from t", true},
	{"select sql_no_cache.a from t", true},
	{"select stored from t", true},
	{"select substring from t", true},
	{"select substring(a from 1) from t", true},
	{"select substring(a, 1) from t", true},
	{"select sum from t", true},
	{"select sum(a) from t", true},
	{"select sum(a) over (partition by loop) from t", false},
	{"select sum(a) over (rows between 1 preceding and current row) from t", true},
	{"select sum(a) over return from t", false},
	{"select sum(a) over w from t window return as (order by a)", false},
	{"select sum(a) over w from t window w as (partition by b)", true},
	{"select sum(a) over window from t", true},
	{"select system", true},
	{"select system from t", true},
	{"select system(1)", true},
	{"select trim from t", true},
	{"select trim(a) from t", true},
	{"select trim(both 'a' from b) from t", true},
	{"select variance (1) from t", true},
	{"select variance from t", true},
	{"select variance from t group by variance", true},
	{"select variance(a) from t", true},
	{"select variance.a from t", true},
	{"select virtual from t", true},
	{"select window from t", true},
	{"select window.a from window", true},
	{"set @a = system", true},
	{"set autocommit = system", true},
	{"set autocommit = system + 1", true},
	{"set read_write = 1", false},
	{"set transaction read write", true},
	{"show create database d", true},
	{"show create function f", true},
	{"show create function function", true},
	{"show create schema d", true},
	{"show databases", true},
	{"show function code f", true},
	{"show function status", true},
	{"show schemas", true},
	{"start transaction read write", true},
	{"update return set loop = 1", false},
	{"update t as window set a = 1", false},
	{"update t window set a = 1", false},
	{"use database", true},
	{"use schema", true},
	{"use system", true},
	{"with get as (select 1) select * from get", true},
	{"with window as (select 1) select * from window", true},
}
