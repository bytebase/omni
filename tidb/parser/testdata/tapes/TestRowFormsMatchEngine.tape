-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestRowFormsMatchEngine default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select row(1)
< error
| 1064 42000 …line 1 column 13 near ")" 
> exec ^
: string
| select row(1,2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select row(1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select row(1,2,3)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select row(1,2,3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select row((1))
< error
| 1064 42000 …line 1 column 15 near ")" 
> exec ^
: string
| select row(row(1,2),3)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select row(row(1,2),3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select row(1) = row(1)
< error
| 1064 42000 …line 1 column 13 near ") = row(1)" 
> exec ^
: string
| select row(1,2) = row(1,2)
< ok 0 0
> exec ^
: string
| select row(1,2) in (row(1,2))
< ok 0 0
> exec ^
: string
| select 1 in (row(1))
< error
| 1064 42000 …line 1 column 19 near "))" 
> exec ^
: string
| select row(1) in (row(1))
< error
| 1064 42000 …line 1 column 13 near ") in (row(1))" 
> exec ^
: string
| select row(1,2) = (1,2)
< ok 0 0
> exec ^
: string
| select * from t where row(1) = 1
< error
| 1064 42000 …line 1 column 28 near ") = 1" 
> exec ^
: string
| select * from t where row(a,b) = (1,2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where row(a,b) = (1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where row(a) in (select 1)
< error
| 1064 42000 …line 1 column 28 near ") in (select 1)" 
> exec ^
: string
| select row(1) from t
< error
| 1064 42000 …line 1 column 13 near ") from t" 
> exec ^
: string
| do row(1)
< error
| 1064 42000 …line 1 column 9 near ")" 
> exec ^
: string
| set @a = row(1)
< error
| 1064 42000 …line 1 column 15 near ")" 
> exec ^
: string
| update t set a = row(1)
< error
| 1064 42000 …line 1 column 23 near ")" 
> exec ^
: string
| insert into t values (row(1))
< error
| 1064 42000 …line 1 column 28 near "))" 
> exec ^
: string
| insert into t values (row(1,2))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (row(1,2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(row(1))
< error
| 1064 42000 …line 1 column 31 near "))" 
> exec ^
: string
| insert into t values row(row(1,2))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(row(1,2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a = row(1)
< error
| 1064 42000 …line 1 column 28 near ")" 
> exec ^
: string
| create table t (a int default (row(1)))
< error
| 1064 42000 …line 1 column 34 near "row(1)))" 
> exec ^
: string
| create table t (a int check (row(1)))
< error
| 1064 42000 …line 1 column 35 near ")))" 
> exec ^
: string
| call p(row(1))
< error
| 1064 42000 …line 1 column 13 near "))" 
> exec ^
: string
| select f(row(1))
< error
| 1064 42000 …line 1 column 15 near "))" 
> exec ^
: string
| select count(row(1))
< error
| 1064 42000 …line 1 column 19 near "))" 
> exec ^
: string
| select if(row(1), 1, 2)
< error
| 1064 42000 …line 1 column 16 near "), 1, 2)" 
> exec ^
: string
| select case row(1) when 1 then 1 end
< error
| 1064 42000 …line 1 column 18 near ") when 1 then 1 end" 
> exec ^
: string
| select row(1) is null
< error
| 1064 42000 …line 1 column 13 near ") is null" 
> exec ^
: string
| select not row(1)
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec ^
: string
| select row (1,2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select row (1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table row (a int)
< error
| 1064 42000 …line 1 column 16 near "row (a int)" 
> exec ^
: string
| select row(1) + 1
< error
| 1064 42000 …line 1 column 13 near ") + 1" 
> exec ^
: string
| select row (1)
< error
| 1064 42000 …line 1 column 14 near ")" 
> exec ^
: string
| create procedure row() select 1
< error
| 1064 42000 …line 1 column 20 near "row() select 1" 
> exec ^
: string
| create function row() returns int return 1
< error
| 1064 42000 …line 1 column 15 near "function row() returns int return 1" 
> exec ^
: string
| select row(default)
< error
| 1064 42000 …line 1 column 19 near ")" 
> exec ^
: string
| insert into t values (row(default, 1))
< error
| 1064 42000 …line 1 column 34 near ", 1))" 
> exec ^
: string
| select row(1, 2) from t where row(1, 2) in (select 1, 2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select row(1, 2) from t where row(1, 2) in (select 1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select row(*)
< error
| 1064 42000 …line 1 column 12 near "*)" 
> exec ^
: string
| select row(1,)
< error
| 1064 42000 …line 1 column 14 near ")" 
> exec ^
: string
| select row(,1)
< error
| 1064 42000 …line 1 column 12 near ",1)" 
> exec ^
: string
| select row(1,,2)
< error
| 1064 42000 …line 1 column 14 near ",2)" 
> exec ^
: string
| select row(t.*, 1)
< error
| 1064 42000 …line 1 column 14 near "*, 1)" 
> exec ^
: string
| select row(1, t.*)
< error
| 1064 42000 …line 1 column 17 near "*)" 
> exec ^
: string
| select row(1, default(a))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select row(1, default(a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (1) = row(1)
< error
| 1064 42000 …line 1 column 19 near ")" 
> exec ^
: string
| select (1,2) = (1,2)
< ok 0 0
> exec ^
: string
| select (1,2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select (1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select ((1,2))
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select ((1,2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (1, 2), (3, 4)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select (1, 2), (3, 4)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (a, b) in ((1, 2), (3, 4)) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select (a, b) in ((1, 2), (3, 4)) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (1,)
< error
| 1064 42000 …line 1 column 11 near ")" 
> exec ^
: string
| select (,1)
< error
| 1064 42000 …line 1 column 9 near ",1)" 
> exec ^
: string
| select ((select 1), 2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select ((select 1), 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (1, (select 2))
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select (1, (select 2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (select 1, 2) = (1, 2)
< ok 0 0
> exec ^
: string
| select (1, 2) in (select 1, 2)
< ok 0 0
> exec ^
: string
| select (1, *)
< error
| 1064 42000 …line 1 column 12 near "*)" 
> exec ^
: string
| select (*, 1)
< error
| 1064 42000 …line 1 column 9 near "*, 1)" 
> exec ^
: string
| select (1, default)
< error
| 1064 42000 …line 1 column 19 near ")" 
> exec ^
: string
| select (default, 1)
< error
| 1064 42000 …line 1 column 16 near ", 1)" 
> exec ^
: string
| select (1, t.*)
< error
| 1064 42000 …line 1 column 14 near "*)" 
> exec ^
: string
| select (t.*, 1)
< error
| 1064 42000 …line 1 column 11 near "*, 1)" 
> exec ^
: string
| select (1, default(a))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select (1, default(a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where (a, b) = (1, 2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where (a, b) = (1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where (a, b) in (select 1, 2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where (a, b) in (select 1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where (a, b) > (1, 2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where (a, b) > (1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| Select (1, 1) > (1, 1)
< ok 0 0
> exec ^
: string
| update t set a = (1, 2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a = (1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values ((1, 2))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values ((1, 2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select ((1, 2), 3)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select ((1, 2), 3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select ((1, 2))
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select ((1, 2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select ((1, 2)) = (1, 2)
< ok 0 0
> exec ^
: string
| select (((1, 2)))
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select (((1, 2)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select row
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| select row from t
< error
| 1064 42000 …line 1 column 15 near "from t" 
> exec ^
: string
| select t.row from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.row from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select row.a from row
< error
| 1064 42000 …line 1 column 21 near "row" 
> exec ^
: string
| select `row` from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `row` from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 row
< error
| 1064 42000 …line 1 column 12 near "row" 
> exec ^
: string
| select 1 as row
< error
| 1064 42000 …line 1 column 15 near "row" 
> exec ^
: string
| select a row from t
< error
| 1064 42000 …line 1 column 12 near "row from t" 
> exec ^
: string
| select * from row
< error
| 1064 42000 …line 1 column 17 near "row" 
> exec ^
: string
| select * from t row
< error
| 1064 42000 …line 1 column 19 near "row" 
> exec ^
: string
| select * from t as row
< error
| 1064 42000 …line 1 column 22 near "row" 
> exec ^
: string
| create table t (row int)
< error
| 1064 42000 …line 1 column 19 near "row int)" 
> exec ^
: string
| create table t (row int, a int)
< error
| 1064 42000 …line 1 column 19 near "row int, a int)" 
> exec ^
: string
| alter table t add column row int
< error
| 1064 42000 …line 1 column 28 near "row int" 
> exec ^
: string
| alter table t drop column row
< error
| 1064 42000 …line 1 column 29 near "row" 
> exec ^
: string
| alter table t rename column row to b
< error
| 1064 42000 …line 1 column 31 near "row to b" 
> exec ^
: string
| create index row on t (a)
< error
| 1064 42000 …line 1 column 16 near "row on t (a)" 
> exec ^
: string
| create index i on t (row)
< error
| 1064 42000 …line 1 column 24 near "row)" 
> exec ^
: string
| insert into t (row) values (1)
< error
| 1064 42000 …line 1 column 18 near "row) values (1)" 
> exec ^
: string
| insert into t (a, row) values (1, 2)
< error
| 1064 42000 …line 1 column 21 near "row) values (1, 2)" 
> exec ^
: string
| replace into t (row) values (1)
< error
| 1064 42000 …line 1 column 19 near "row) values (1)" 
> exec ^
: string
| insert into t set row = 1
< error
| 1064 42000 …line 1 column 21 near "row = 1" 
> exec ^
: string
| update t set row = 1
< error
| 1064 42000 …line 1 column 16 near "row = 1" 
> exec ^
: string
| update t set row := 1
< error
| 1064 42000 …line 1 column 16 near "row := 1" 
> exec ^
: string
| update t set t.row = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set t.row = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from row
< error
| 1064 42000 …line 1 column 15 near "row" 
> exec ^
: string
| delete from t where row = 1
< error
| 1064 42000 …line 1 column 25 near "= 1" 
> exec ^
: string
| select * from t order by row
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| select * from t group by row
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| select row, a from t
< error
| 1064 42000 …line 1 column 11 near ", a from t" 
> exec ^
: string
| select a, row from t
< error
| 1064 42000 …line 1 column 18 near "from t" 
> exec ^
: string
| select row + 1
< error
| 1064 42000 …line 1 column 12 near "+ 1" 
> exec ^
: string
| select row = 1
< error
| 1064 42000 …line 1 column 12 near "= 1" 
> exec ^
: string
| create view row as select 1
< error
| 1064 42000 …line 1 column 15 near "row as select 1" 
> exec ^
: string
| drop table row
< error
| 1064 42000 …line 1 column 14 near "row" 
> exec ^
: string
| truncate table row
< error
| 1064 42000 …line 1 column 18 near "row" 
> exec ^
: string
| use row
< error
| 1064 42000 …line 1 column 7 near "row" 
> exec ^
: string
| create database row
< error
| 1064 42000 …line 1 column 19 near "row" 
> exec ^
: string
| set @row = 1
< ok 0 0
> exec ^
: string
| select @row
< ok 0 0
> exec ^
: string
| declare row int
< error
| 1064 42000 …line 1 column 7 near "declare row int" 
> exec ^
: string
| load data infile 'f' into table t (row)
< error
| 1064 42000 …line 1 column 38 near "row)" 
> exec ^
: string
| show columns from row
< error
| 1064 42000 …line 1 column 21 near "row" 
> exec ^
: string
| grant select on row to u
< error
| 1064 42000 …line 1 column 19 near "row to u" 
> exec ^
: string
| create table t (a int) partition by hash(a) (partition row)
< error
| 1064 42000 …line 1 column 58 near "row)" 
> exec ^
: string
| select row.row from row
< error
| 1064 42000 …line 1 column 23 near "row" 
> exec ^
: string
| select db.row.a from row
< error
| 1064 42000 …line 1 column 24 near "row" 
> exec ^
: string
| select count(*) over (rows current row)
< ok 0 0
> exec ^
: string
| select count(*) over (rows between current row and unbounded following)
< ok 0 0
> exec ^
: string
| create table t (row varchar(10))
< error
| 1064 42000 …line 1 column 19 near "row varchar(10))" 
> exec ^
: string
| create table t (row int) engine = innodb
< error
| 1064 42000 …line 1 column 19 near "row int) engine = innodb" 
> exec ^
: string
| select a.row from a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select a.row from a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select row.row.row from row
< error
| 1064 42000 …line 1 column 27 near "row" 
> exec ^
: string
| set row = 1
< error
| 1064 42000 …line 1 column 7 near "row = 1" 
> exec ^
: string
| insert into t values (row)
< error
| 1064 42000 …line 1 column 26 near ")" 
> exec ^
: string
| insert into t values (row, 1)
< error
| 1064 42000 …line 1 column 26 near ", 1)" 
> exec ^
: string
| insert into t values (row(1,2), row)
< error
| 1064 42000 …line 1 column 36 near ")" 
> exec ^
: string
| select row(row, 1)
< error
| 1064 42000 …line 1 column 15 near ", 1)" 
> exec ^
: string
| select row(1, 2) as row
< error
| 1064 42000 …line 1 column 23 near "row" 
> exec ^
: string
| select row(row(row, row), row)
< error
| 1064 42000 …line 1 column 19 near ", row), row)" 
> exec ^
: string
| select row as row from row as row
< error
| 1064 42000 …line 1 column 13 near "as row from row as row" 
> exec ^
: string
| insert into t (row) value (row)
< error
| 1064 42000 …line 1 column 18 near "row) value (row)" 
> exec ^
: string
| update t set row = row + 1
< error
| 1064 42000 …line 1 column 16 near "row = row + 1" 
> exec ^
: string
| select rows from t
< error
| 1064 42000 …line 1 column 11 near "rows from t" 
> exec ^
: string
| create table t (rows int)
< error
| 1064 42000 …line 1 column 20 near "rows int)" 
> exec ^
: string
| select * from t rows
< error
| 1064 42000 …line 1 column 20 near "rows" 
> exec ^
: string
| select 1 rows
< error
| 1064 42000 …line 1 column 13 near "rows" 
> exec ^
: string
| create table rows (a int)
< error
| 1064 42000 …line 1 column 17 near "rows (a int)" 
> exec ^
: string
| select row, rows from t
< error
| 1064 42000 …line 1 column 11 near ", rows from t" 
> exec ^
: string
| insert into t values (ignore)
< error
| 1064 42000 …line 1 column 28 near "ignore)" 
> exec ^
: string
| insert into t set a = ignore
< error
| 1064 42000 …line 1 column 28 near "ignore" 
> exec ^
: string
| update t set a = ignore
< error
| 1064 42000 …line 1 column 23 near "ignore" 
> exec ^
: string
| select ignore
< error
| 1064 42000 …line 1 column 13 near "ignore" 
> exec ^
: string
| select * from t ignore index (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t ignore index (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
