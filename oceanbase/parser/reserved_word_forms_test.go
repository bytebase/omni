package parser

import (
	"sort"
	"strings"
	"testing"
)

// Every registered keyword as a name, and the statements the words start beside it, every verdict
// measured with PREPARE on OceanBase CE 4.3.5.6 (oceanbase/oceanbase-ce:4.3.5-lts) on 2026-11-05: a
// 1064 is the parser's refusal, and any other answer is a statement that parsed.
//
// keywordCategories began as MySQL's, and these are where OceanBase differs. It leaves unreserved
// the words oceanBaseUnreservedWords lists, among them FUNCTION, SYSTEM, RANGE and the window
// functions, and it reserves SKIP (unreserved_keyword_normal, sql_parser_mysql_mode.y:25818). A word
// the lexer knows only as a function is an identifier unless `(` follows it.
// unreserved_keyword_ambiguous_roles are names everywhere but a role and an unscoped SET. A stored
// program is read by another grammar, whose labels isLabelKeyword follows.
//
// reservedWordPositions leaves out `SELECT word(1)`, where omni reads any word as a function, and a
// DEFAULT or ON UPDATE of each word, which takes function calls and literals the server refuses.
// Each is its own change.

// TestReservedWordForms holds omni to the recorded verdicts without a server.
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

// registeredKeywords lists the lexer's words in a stable order.
func registeredKeywords() []string {
	words := make([]string, 0, len(keywords))
	for w := range keywords {
		words = append(words, w)
	}
	sort.Strings(words)
	return words
}

// checkReservedWordSweep tries every registered keyword in each name position, which is what keeps
// keywordCategories this engine's. A word added to the lexer with the wrong category fails here.
// accepts is the server's verdict.
func checkReservedWordSweep(t *testing.T, accepts func(sql string) bool) {
	t.Helper()
	for _, position := range reservedWordPositions {
		for _, word := range registeredKeywords() {
			sql := strings.ReplaceAll(position, "%s", word)
			server := accepts(sql)
			omni, err := parsesAsOneStatement(sql)
			if reservedWordSweepGap[sql] {
				if server == omni {
					t.Errorf("%q now agrees: drop it from reservedWordSweepGap", sql)
				}
				continue
			}
			if server != omni {
				t.Errorf("DIVERGENCE omni accepts=%v, server accepts=%v for %q: %v", omni, server, sql, err)
			}
		}
	}
}

// reservedWordPositions are the name positions the sweep tries every keyword in.
var reservedWordPositions = []string{
	"select %s from t",
	"select t.%s from t",
	"select a %s from t",
	"select a as %s from t",
	"select * from %s",
	"select * from t %s",
	"select * from t as %s",
	"create table %s (a int)",
	"create table x (%s int)",
	"alter table t add column %s int",
	"create index %s on t (a)",
	"create procedure p() %s: loop leave %s; end loop",
	"create role %s",
	"set %s = 1",
	"select %s.a from t",
	"create table x (a int, key %s (a))",
	"select 1 from t where a = %s",
	"update t set a = %s",
	"insert into t values (%s)",
	"select %s + 1 from t",
}

// reservedWordSweepGap holds the forms where omni and the server still disagree. A 1314 is the
// parser's action refusing a statement before it reads on, so it says nothing of the label. The sweep
// fails when one starts to agree, so this list cannot outlive the gap.
var reservedWordSweepGap = map[string]bool{
	"create procedure p() comment: loop leave comment; end loop":   true,
	"create procedure p() contains: loop leave contains; end loop": true,
	"select default from t":                                        true,
	"select 1 from t where a = default":                            true,
	"select default + 1 from t":                                    true,
	"create procedure p() end: loop leave end; end loop":           true,
	"create procedure p() language: loop leave language; end loop": true,
	"select unique + 1 from t":                                     true,
}

var reservedWordCases = []asOfCase{
	{"select a full from t", true},
	{"set maxvalue = 1", false},
	{"set range = 1", false},
	{"set global maxvalue = 1", true},
	{"set global range = 1", true},
	{"select count from t", true},
	{"select count(*) from t", true},
	{"select sum from t", true},
	{"select sum(a) from t", true},
	{"select now from t", true},
	{"select now() from t", true},
	{"select now.a from t", true},
	{"select a now from t", true},
	{"select a as now from t", true},
	{"select cast from t", true},
	{"select cast(a as char) from t", true},
	{"select extract from t", true},
	{"select extract(year from a) from t", true},
	{"select trim from t", true},
	{"select trim(a) from t", true},
	{"select substring from t", true},
	{"select substring(a, 1) from t", true},
	{"select position from t", true},
	{"select position('a' in a) from t", true},
	{"select group_concat from t", true},
	{"select group_concat(a) from t", true},
	{"select max from t", true},
	{"select max(a) from t", true},
	{"select variance from t", true},
	{"select variance(a) from t", true},
	{"select count + 1 from t", true},
	{"select 1 from t where count = 1", true},
	{"select count from t order by count", true},
	{"select t.count from t", true},
	{"select * from count", true},
	{"select if from t", false},
	{"select if(1, 2, 3)", true},
	{"select left from t", false},
	{"select left(a, 1) from t", true},
	{"select rank from t", true},
	{"select rank() over () from t", true},
	{"select lag from t", true},
	{"select lag(a) over () from t", true},
	{"select row_number from t", true},
	{"select row_number() over () from t", true},
	{"select first_value from t", true},
	{"select first_value(a) over () from t", true},
	{"select grouping from t", true},
	{"select a, grouping(a) from t group by a with rollup", true},
	{"select lateral from t", true},
	{"select * from t, lateral (select 1) x", true},
	{"select * from lateral", true},
	{"select json_table from t", true},
	{"select * from json_table('[]', '$' columns (a int path '$')) j", true},
	{"select * from json_table", true},
	{"select * from t window w as ()", true},
	{"select a, rank() over w from t window w as ()", true},
	{"select a window from t", false},
	{"select * from t window", false},
	{"select 1 from t except select 1", true},
	{"select 1 from t intersect select 1", true},
	{"select * from t except select 1", true},
	{"select * from t intersect select 1", true},
	{"select a except from t", false},
	{"select * from t except", false},
	{"select a over from t", true},
	{"select count(*) over () from t", true},
	{"create database d", true},
	{"drop database d", true},
	{"create schema d", true},
	{"drop schema d", true},
	{"show databases", true},
	{"show schemas", true},
	{"select database()", true},
	{"select database from t", false},
	{"select schema from t", false},
	{"create function f() returns int return 1", true},
	{"drop function f", true},
	{"create table function (a int)", true},
	{"create table x (function int)", true},
	{"get diagnostics @a = number", true},
	{"select get from t", false},
	{"create table get (a int)", false},
	{"grant select on t to a with grant option", true},
	{"select option from t", false},
	{"create table option (a int)", false},
	{"create view v as select 1 with check option", true},
	{"create table x (a int generated always as (1) stored)", true},
	{"create table x (a int generated always as (1) virtual)", true},
	{"create table x (a int as (1) stored)", true},
	{"create table x (a int as (1) virtual)", true},
	{"create table x (generated int)", false},
	{"create table x (stored int)", false},
	{"create table x (virtual int)", false},
	{"select system from t", true},
	{"select a from t for system_time as of now()", false},
	{"select * from t for update of t", false},
	{"select of from t", true},
	{"create table of (a int)", true},
	{"create table x (a dec)", true},
	{"create table x (a dec(3))", true},
	{"create table x (dec int)", true},
	{"create table x (a int) partition by range (a) (partition p values less than (1))", true},
	{"create table x (a int) partition by range (a) (partition p values less than maxvalue)", true},
	{"create table range (a int)", true},
	{"create table maxvalue (a int)", true},
	{"select cube from t", true},
	{"select a from t group by a with cube", false},
	{"select recursive from t", true},
	{"with recursive c as (select 1) select * from c", true},
	{"select empty from t", true},
	{"select * from json_table('[]', '$' columns (a int path '$' default '0' on empty)) j", true},
	{"create table x (a timestamp default now)", false},
	{"create table x (a timestamp default now())", true},
	{"create table x (a timestamp default current_timestamp)", true},
	{"create table x (a timestamp default current_timestamp(6))", true},
	{"create table x (a int default b)", false},
	{"create table x (a int default 1)", true},
	{"create table x (a int default (1))", true},
	{"create table x (a int default (b))", true},
	{"create table x (a timestamp on update now)", false},
	{"create table x (a timestamp on update now())", true},
	{"create table x (a timestamp on update current_timestamp)", true},
	{"create table x (a int on update b)", false},
	{"create table x (a timestamp default count)", false},
	{"select a into from t", false},
	{"select * from t into", false},
	{"select a into @b from t", true},
	{"select a from t into @b", true},
	{"select 1 into @a, @b", true},
	{"select a into outfile 'x' from t", true},
	{"select a into dumpfile 'x' from t", true},
	{"select into from t", false},
	{"create role event", false},
	{"create role `event`", true},
	{"create role 'event'", true},
	{"create role a, event", false},
	{"create role process", false},
	{"create role none", false},
	{"create role shutdown", false},
	{"create role super", false},
	{"create role a", true},
	{"create role a, b", true},
	{"create role if not exists event", false},
	{"create role if not exists a", true},
	{"create procedure p() close: loop leave close; end loop", true},
	{"create procedure p() open: loop leave open; end loop", true},
	{"create procedure p() begin: loop leave begin; end loop", true},
	{"create procedure p() add: loop leave add; end loop", true},
	{"create procedure p() select: loop leave select; end loop", false},
	{"create procedure p() l: loop leave l; end loop", true},
	{"create procedure p() begin close: loop leave close; end loop; end", true},
	{"create procedure p() l: begin select 1; end l", true},
	{"create procedure p() begin declare c cursor for select 1; open c; close c; end", true},
	{"create procedure p() begin begin select 1; end; end", true},
	{"create procedure p() begin end", true},
	{"create procedure p() begin declare a int; begin end; end", true},
	{"set event = 1", false},
	{"set global event = 1", true},
	{"set @a = 1, event = 2", false},
	{"set process = 1", false},
	{"set none = 1", false},
	{"set skip = 1", false},
	{"select skip from t", false},
	{"create table skip (a int)", false},
}
