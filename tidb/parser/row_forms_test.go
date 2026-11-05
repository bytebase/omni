package parser

import "testing"

// The row constructor on TiDB v8.5.5, where ROW is reserved and IGNORE is no value.
//
// The verdicts were measured against pingcap/tidb:v8.5.5 through PREPARE on 2026-11-05, and the
// recorded answer is what that server gave. TestRowFormsMatchEngine holds them to the
// server, and row_tree_test.go holds the tree.

// rowFormCase is a statement and whether the engine's parser accepts it.
type rowFormCase struct {
	sql    string
	accept bool
}

var rowFormCases = []rowFormCase{
	// ROW(...) has at least two elements, and they are expressions.
	{`select row(1)`, false},
	{`select row(1,2)`, true},
	{`select row(1,2,3)`, true},
	{`select row((1))`, false},
	{`select row(row(1,2),3)`, true},
	{`select row(1) = row(1)`, false},
	{`select row(1,2) = row(1,2)`, true},
	{`select row(1,2) in (row(1,2))`, true},
	{`select 1 in (row(1))`, false},
	{`select row(1) in (row(1))`, false},
	{`select row(1,2) = (1,2)`, true},
	{`select * from t where row(1) = 1`, false},
	{`select * from t where row(a,b) = (1,2)`, true},
	{`select * from t where row(a) in (select 1)`, false},
	{`select row(1) from t`, false},
	{`do row(1)`, false},
	{`set @a = row(1)`, false},
	{`update t set a = row(1)`, false},
	{`insert into t values (row(1))`, false},
	{`insert into t values (row(1,2))`, true},
	{`insert into t values row(row(1))`, false},
	{`insert into t values row(row(1,2))`, true},
	{`insert into t set a = row(1)`, false},
	{`create table t (a int default (row(1)))`, false},
	{`create table t (a int check (row(1)))`, false},
	{`call p(row(1))`, false},
	{`select f(row(1))`, false},
	{`select count(row(1))`, false},
	{`select if(row(1), 1, 2)`, false},
	{`select case row(1) when 1 then 1 end`, false},
	{`select row(1) is null`, false},
	{`select not row(1)`, false},
	{`select row (1,2)`, true},
	{`create table row (a int)`, false},
	{`select row(1) + 1`, false},
	{`select row (1)`, false},
	{`create procedure row() select 1`, false},
	{`create function row() returns int return 1`, false},
	{`select row(default)`, false},
	{`insert into t values (row(default, 1))`, false},
	{`select row(1, 2) from t where row(1, 2) in (select 1, 2)`, true},
	{`select row(*)`, false},
	{`select row(1,)`, false},
	{`select row(,1)`, false},
	{`select row(1,,2)`, false},
	{`select row(t.*, 1)`, false},
	{`select row(1, t.*)`, false},
	{`select row(1, default(a))`, true},

	// A parenthesized row reads the same way.
	{`select (1) = row(1)`, false},
	{`select (1,2) = (1,2)`, true},
	{`select (1,2)`, true},
	{`select ((1,2))`, true},
	{`select (1, 2), (3, 4)`, true},
	{`select (a, b) in ((1, 2), (3, 4)) from t`, true},
	{`select (1,)`, false},
	{`select (,1)`, false},
	{`select ((select 1), 2)`, true},
	{`select (1, (select 2))`, true},
	{`select (select 1, 2) = (1, 2)`, true},
	{`select (1, 2) in (select 1, 2)`, true},
	{`select (1, *)`, false},
	{`select (*, 1)`, false},
	{`select (1, default)`, false},
	{`select (default, 1)`, false},
	{`select (1, t.*)`, false},
	{`select (t.*, 1)`, false},
	{`select (1, default(a))`, true},
	{`select * from t where (a, b) = (1, 2)`, true},
	{`select * from t where (a, b) in (select 1, 2)`, true},
	{`select * from t where (a, b) > (1, 2)`, true},
	{`Select (1, 1) > (1, 1)`, true},
	{`update t set a = (1, 2)`, true},
	{`insert into t values ((1, 2))`, true},
	{`select ((1, 2), 3)`, true},
	{`select ((1, 2))`, true},
	{`select ((1, 2)) = (1, 2)`, true},
	{`select (((1, 2)))`, true},

	// ROW as a name.
	{`select row`, false},
	{`select row from t`, false},
	{`select t.row from t`, true},
	{`select row.a from row`, false},
	{"select `row` from t", true},
	{`select 1 row`, false},
	{`select 1 as row`, false},
	{`select a row from t`, false},
	{`select * from row`, false},
	{`select * from t row`, false},
	{`select * from t as row`, false},
	{`create table t (row int)`, false},
	{`create table t (row int, a int)`, false},
	{`alter table t add column row int`, false},
	{`alter table t drop column row`, false},
	{`alter table t rename column row to b`, false},
	{`create index row on t (a)`, false},
	{`create index i on t (row)`, false},
	{`insert into t (row) values (1)`, false},
	{`insert into t (a, row) values (1, 2)`, false},
	{`replace into t (row) values (1)`, false},
	{`insert into t set row = 1`, false},
	{`update t set row = 1`, false},
	{`update t set row := 1`, false},
	{`update t set t.row = 1`, true},
	{`delete from row`, false},
	{`delete from t where row = 1`, false},
	{`select * from t order by row`, false},
	{`select * from t group by row`, false},
	{`select row, a from t`, false},
	{`select a, row from t`, false},
	{`select row + 1`, false},
	{`select row = 1`, false},
	{`create view row as select 1`, false},
	{`drop table row`, false},
	{`truncate table row`, false},
	{`use row`, false},
	{`create database row`, false},
	{`set @row = 1`, true},
	{`select @row`, true},
	{`declare row int`, false},
	{`load data infile 'f' into table t (row)`, false},
	{`show columns from row`, false},
	{`grant select on row to u`, false},
	{`create table t (a int) partition by hash(a) (partition row)`, false},
	{`select row.row from row`, false},
	{`select db.row.a from row`, false},
	{`select count(*) over (rows current row)`, true},
	{`select count(*) over (rows between current row and unbounded following)`, true},
	{`create table t (row varchar(10))`, false},
	{`create table t (row int) engine = innodb`, false},
	{`select a.row from a`, true},
	{`select row.row.row from row`, false},
	{`set row = 1`, false},
	{`insert into t values (row)`, false},
	{`insert into t values (row, 1)`, false},
	{`insert into t values (row(1,2), row)`, false},
	{`select row(row, 1)`, false},
	{`select row(1, 2) as row`, false},
	{`select row(row(row, row), row)`, false},
	{`select row as row from row as row`, false},
	{`insert into t (row) value (row)`, false},
	{`update t set row = row + 1`, false},

	// ROWS as a name.
	{`select rows from t`, false},
	{`create table t (rows int)`, false},
	{`select * from t rows`, false},
	{`select 1 rows`, false},
	{`create table rows (a int)`, false},
	{`select row, rows from t`, false},

	// IGNORE as a value, where the grammar reads expr_or_ignore, and nowhere else.
	{`insert into t values (ignore)`, false},
	{`insert into t set a = ignore`, false},
	{`update t set a = ignore`, false},
	{`select ignore`, false},

	// Window frames keep reading ROW and ROWS as keywords.
	{`select * from t ignore index (a)`, true},
}

// TestRowForms holds the parser to the verdicts above. An accepted statement must read
// as one: Parse takes statements with no separator between them, so a tail the grammar
// should refuse could otherwise pass as a second one.
func TestRowForms(t *testing.T) {
	for _, c := range rowFormCases {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if got := err == nil; got != c.accept {
				t.Fatalf("Parse(%q) accepts = %v, the engine answered %v (err: %v)", c.sql, got, c.accept, err)
			}
			if err == nil && len(list.Items) != 1 {
				t.Errorf("Parse(%q) read %d statements, want 1", c.sql, len(list.Items))
			}
		})
	}
}
