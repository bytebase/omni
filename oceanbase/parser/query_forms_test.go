package parser

import "testing"

// OceanBase's select grammar is narrower than MySQL 8.0's in two places this
// parser used to read the MySQL way.
//
// It has no TABLE query. select_clause and select_clause_set_body in OceanBase
// CE v4.3.5_CE_BP1 src/sql/parser/sql_parser_mysql_mode.y (lines 11088 and
// 11162) take a SELECT, a parenthesized query and VALUES, and nothing else in
// the grammar reads TABLE as a query, so "TABLE t" is 1064 wherever a query
// goes. tableQueryFormCases puts it in each of those places, next to a VALUES
// or SELECT in the same place as a control.
//
// It takes no locking clause right after a lone parenthesized query.
// select_no_parens (line 11016) puts opt_lock_type after a select_clause or a
// set operation, and a parenthesized query is a select_clause only with an
// ORDER BY or a LIMIT after it (select_with_parens_with_order_and_limit, line
// 11259). A lone one is a select_stmt by itself (line 10969), so "(SELECT 1)
// FOR UPDATE" is 1064 while "(SELECT 1) ORDER BY 1 FOR UPDATE" and "(SELECT 1)
// UNION (SELECT 2) FOR UPDATE" parse. parenLockFormCases holds both sides, in
// every place a query goes.
//
// Every case was measured against OceanBase CE 4.3.5.6
// (oceanbase/oceanbase-ce:4.3.5-lts) through PREPARE on 2026-11-05.
// TestQueryFormsOceanBaseOracle keeps it a recording: an image upgrade that
// changes an answer fails there rather than leaving a stale expectation behind.
// MySQL 8.0 accepts every TABLE form and every locking form here.
//
// That oracle runs on the package's own OceanBase (see
// sharedserver.PrivateOceanBase):
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1

// queryFormCase is one statement text and whether OceanBase's parser accepts
// it.
type queryFormCase struct {
	sql    string
	accept bool
}

var tableQueryFormCases = []queryFormCase{
	// A statement.
	{"table t", false},
	{"table t order by a", false},
	{"table t limit 1", false},
	{"table t into @x", false},
	{"values row(1)", true},
	// A parenthesized query, a subquery, a derived table, a CTE.
	{"(table t)", false},
	{"((table t))", false},
	{"(values row(1))", true},
	{"select (table t)", false},
	{"select ((table t) = 1)", false},
	{"select exists (table t)", false},
	{"select * from t where a in (table t)", false},
	{"select * from t where a = any (table t)", false},
	{"select * from t where a in (select a from t)", true},
	{"select * from (table t) d", false},
	{"select * from t join (table u) d on 1", false},
	{"select * from (values row(1)) d", true},
	{"with c as (table t) select * from c", false},
	{"with c as (select 1) table c", false},
	{"with c as (values row(1)) select * from c", true},
	// A set-operation operand.
	{"select 1 union table t", false},
	{"table t union select 1", false},
	{"(select 1) union (table t)", false},
	{"values row(1) union table t", false},
	{"select 1 union values row(1)", true},
	{"(values row(1)) union (select 1)", true},
	// The source of INSERT, REPLACE, EXPLAIN, CREATE VIEW, and CREATE TABLE.
	{"insert into t table u", false},
	{"insert into t (table u)", false},
	{"insert into t select 1 union table u", false},
	{"insert into t select 1 union select 2", true},
	{"replace into t table u", false},
	{"explain table t", false},
	{"explain (table t)", false},
	{"explain select 1 union table t", false},
	{"explain select 1 union select 2", true},
	{"create view v as table t", false},
	{"create view v as select 1 union table t", false},
	{"create view v as select 1 union select 2", true},
	{"create table n as table t", false},
	{"create table n (table t)", false},
	{"create table n as values row(1)", true},
	// A stored program body.
	{"create procedure p() begin table t; end", false},
	{"create procedure p() begin declare c cursor for table t; end", false},
	{"create procedure p() begin declare c cursor for select 1; end", true},
}

var parenLockFormCases = []queryFormCase{
	// A lone parenthesized query takes no locking clause, at any depth and
	// in any spelling.
	{"(select 1) for update", false},
	{"(select 1) lock in share mode", false},
	{"(select 1) for update nowait", false},
	{"(select 1) for update wait 1", false},
	{"(select 1) for update skip locked", false},
	{"((select 1)) for update", false},
	{"((select 1) order by 1) for update", false},
	{"((select 1) for update)", false},
	{"(select 1 for update) for update", false},
	{"((select 1) union (select 2)) for update", false},
	{"(select 1 union select 2) for update", false},
	{"with c as (select 1) (select 1) for update", false},
	// It takes one inside the parentheses, after an ORDER BY or a LIMIT, and
	// after a set operation.
	{"select 1 for update", true},
	{"(select 1 for update)", true},
	{"(select 1) order by 1 for update", true},
	{"(select 1) limit 1 for update", true},
	{"(select 1) order by 1 limit 1 for update", true},
	{"(select 1) order by 1 lock in share mode", true},
	{"(select 1) order by 1 for update nowait", true},
	{"(select 1) limit 1 for update wait 1", true},
	{"(select 1) order by 1 for update skip locked", true},
	{"((select 1) order by 1 for update)", true},
	{"(select 1 for update) order by 1 for update", true},
	{"(select 1) union (select 2) for update", true},
	{"(select 1) union all (select 2) for update", true},
	{"(select 1) union (select 2) lock in share mode", true},
	{"select 1 union (select 2) for update", true},
	{"(select 1) union select 2 for update", true},
	{"with c as (select 1) (select 1) order by 1 for update", true},
	// The same, wherever a query goes.
	{"select * from ((select 1) for update) d", false},
	{"select * from ((select 1) order by 1 for update) d", true},
	{"select * from ((select 1) union (select 2) for update) d", true},
	{"select ((select 1) for update)", false},
	{"select (select 1 for update)", true},
	{"select 1 in ((select 1) for update)", false},
	{"select exists ((select 1) for update)", false},
	{"select exists ((select 1) order by 1 for update)", true},
	{"insert into t (select 1) for update", false},
	{"insert into t (select 1) order by 1 for update", true},
	{"insert into t select 1 for update", true},
	{"explain (select 1) for update", false},
	{"explain (select 1) order by 1 for update", true},
	{"create view v as (select 1) for update", false},
	{"create view v as (select 1) order by 1 for update", true},
	{"create view v as select 1 for update", true},
	{"create table n (select 1) for update", false},
	{"create table n (select 1) order by 1 for update", true},
	{"create table n as (select 1) for update", false},
	{"create table n as (select 1) limit 1 for update", true},
	{"create procedure p() begin (select 1) for update; end", false},
	{"create procedure p() begin (select 1) order by 1 for update; end", true},
	{"create procedure p() begin declare c cursor for (select 1) for update; end", false},
	{"create procedure p() begin declare c cursor for (select 1) limit 1 for update; end", true},
}

func TestQueryForms(t *testing.T) {
	for _, cases := range [][]queryFormCase{tableQueryFormCases, parenLockFormCases} {
		for _, c := range cases {
			t.Run(c.sql, func(t *testing.T) {
				_, err := Parse(c.sql)
				if accepts := err == nil; accepts != c.accept {
					t.Errorf("omni accepts=%v, want %v for %q: %v", accepts, c.accept, c.sql, err)
				}
			})
		}
	}
}
