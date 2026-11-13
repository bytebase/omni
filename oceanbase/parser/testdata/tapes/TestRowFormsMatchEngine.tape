-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestRowFormsMatchEngine default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_sweep
< ok 1 0
> exec
| USE omni_sweep
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM ?
: string
| select row(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select row(1,2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select row(1,2,3)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select row((1))
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select row(row(1,2),3)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select row(1) = row(1)
< error
| 1064 42000 …near ') = row(1)' at line 1
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
| 1064 42000 …near '))' at line 1
> exec ^
: string
| select row(1) in (row(1))
< error
| 1064 42000 …near ') in (row(1))' at line 1
> exec ^
: string
| select row(1,2) = (1,2)
< ok 0 0
> exec ^
: string
| select * from t where row(1) = 1
< error
| 1064 42000 …near ') = 1' at line 1
> exec ^
: string
| select * from t where row(a,b) = (1,2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t where row(a) in (select 1)
< error
| 1064 42000 …near ') in (select 1)' at line 1
> exec ^
: string
| select row(1) from t
< error
| 1064 42000 …near ') from t' at line 1
> exec ^
: string
| do row(1)
< error
| 1064 42000 …near 'row(1)' at line 1
> exec ^
: string
| set @a = row(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| update t set a = row(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values (row(1))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| insert into t values (row(1,2))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(row(1))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| insert into t values row(row(1,2))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t set a = row(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| create table t (a int default (row(1)))
< error
| 1064 42000 …near ')))' at line 1
> exec ^
: string
| create table t (a int check (row(1)))
< error
| 1064 42000 …near ')))' at line 1
> exec ^
: string
| call p(row(1))
< error
| 1064 42000 …near 'row(1)' at line 1
> exec ^
: string
| select f(row(1))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| select count(row(1))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| select if(row(1), 1, 2)
< error
| 1064 42000 …near '), 1, 2)' at line 1
> exec ^
: string
| select case row(1) when 1 then 1 end
< error
| 1064 42000 …near ') when 1 then 1 end' at line 1
> exec ^
: string
| select row(1) is null
< error
| 1064 42000 …near ') is null' at line 1
> exec ^
: string
| select not row(1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select row (1,2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| create table row (a int)
< ok 0 0
> exec ^
: string
| select row(1) + 1
< error
| 1064 42000 …near ') + 1' at line 1
> exec ^
: string
| select row (1)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| create procedure row() select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create function row() returns int return 1
< ok 0 0
> exec ^
: string
| select row(default)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values (row(default, 1))
< error
| 1064 42000 …near ', 1))' at line 1
> exec ^
: string
| select row(1, 2) from t where row(1, 2) in (select 1, 2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select row(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| select row(1,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select row(,1)
< error
| 1064 42000 …near ',1)' at line 1
> exec ^
: string
| select row(1,,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| select row(t.*, 1)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec ^
: string
| select row(1, t.*)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec ^
: string
| select row(1, default(a))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select (1) = row(1)
< error
| 1064 42000 …near ')' at line 1
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
| select ((1,2))
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select (1, 2), (3, 4)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select (a, b) in ((1, 2), (3, 4)) from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select (1,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select (,1)
< error
| 1064 42000 …near ',1)' at line 1
> exec ^
: string
| select ((select 1), 2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select (1, (select 2))
< error
| 1241 21000 Operand should contain 1 column(s)
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
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| select (*, 1)
< error
| 1064 42000 …near '*, 1)' at line 1
> exec ^
: string
| select (1, default)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select (default, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| select (1, t.*)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec ^
: string
| select (t.*, 1)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec ^
: string
| select (1, default(a))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select * from t where (a, b) = (1, 2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t where (a, b) in (select 1, 2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t where (a, b) > (1, 2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| Select (1, 1) > (1, 1)
< ok 0 0
> exec ^
: string
| update t set a = (1, 2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values ((1, 2))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select ((1, 2), 3)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select ((1, 2))
< error
| 1241 21000 Operand should contain 1 column(s)
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
| select row
< error
| 1054 42S22 Unknown column 'row' in 'field list'
> exec ^
: string
| select row from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select t.row from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select row.a from row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| select `row` from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select 1 row
< ok 0 0
> exec ^
: string
| select 1 as row
< ok 0 0
> exec ^
: string
| select a row from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| select * from t row
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t as row
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| create table t (row int)
< ok 0 0
> exec ^
: string
| create table t (row int, a int)
< ok 0 0
> exec ^
: string
| alter table t add column row int
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| alter table t drop column row
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| alter table t rename column row to b
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| create index row on t (a)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| create index i on t (row)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (row) values (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a, row) values (1, 2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t (row) values (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t set row = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| update t set row = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| update t set row := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| update t set t.row = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| delete from row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| delete from t where row = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t order by row
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t group by row
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select row, a from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select a, row from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select row + 1
< error
| 1054 42S22 Unknown column 'row' in 'field list'
> exec ^
: string
| select row = 1
< error
| 1054 42S22 Unknown column 'row' in 'field list'
> exec ^
: string
| create view row as select 1
< ok 0 0
> exec ^
: string
| drop table row
< ok 0 0
> exec ^
: string
| truncate table row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| use row
< error
| 1049 42000 Unknown database 'row'
> exec ^
: string
| create database row
< ok 0 0
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
| 1064 42000 …near 'declare' at line 1
> exec ^
: string
| load data infile 'f' into table t (row)
< error
| 1017 HY000 File not exist
> exec ^
: string
| show columns from row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| grant select on row to u
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| create table t (a int) partition by hash(a) (partition row)
< ok 0 0
> exec ^
: string
| select row.row from row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| select db.row.a from row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
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
< ok 0 0
> exec ^
: string
| create table t (row int) engine = innodb
< ok 0 0
> exec ^
: string
| select a.row from a
< error
| 1146 42S02 Table 'omni_sweep.a' doesn't exist
> exec ^
: string
| select row.row.row from row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| set row = 1
< ok 0 0
> exec ^
: string
| insert into t values (row)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (row, 1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (row(1,2), row)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select row(row, 1)
< error
| 1054 42S22 Unknown column 'row' in 'field list'
> exec ^
: string
| select row(1, 2) as row
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| select row(row(row, row), row)
< error
| 1054 42S22 Unknown column 'row' in 'field list'
> exec ^
: string
| select row as row from row as row
< error
| 1146 42S02 Table 'omni_sweep.row' doesn't exist
> exec ^
: string
| insert into t (row) value (row)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| update t set row = row + 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select rows from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| create table t (rows int)
< ok 0 0
> exec ^
: string
| select * from t rows
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select 1 rows
< ok 0 0
> exec ^
: string
| create table rows (a int)
< ok 0 0
> exec ^
: string
| select row, rows from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (ignore)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t set a = ignore
< error
| 1064 42000 …near 'ignore' at line 1
> exec ^
: string
| update t set a = ignore
< error
| 1064 42000 …near 'ignore' at line 1
> exec ^
: string
| select ignore
< error
| 1064 42000 …near 'ignore' at line 1
> exec ^
: string
| select * from t ignore index (a)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
