-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestRowFormsMatchEngine default#1
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1064 42000 …near ')' at line 1
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| insert into t values row(row(1))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| insert into t values row(row(1,2))
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1064 42000 …near '))' at line 1
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
< error
| 1064 42000 …near 'row (a int)' at line 1
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
| 1064 42000 …near 'row() select 1' at line 1
> exec ^
: string
| create function row() returns int return 1
< error
| 1064 42000 …near 'row() returns int return 1' at line 1
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1064 42000 …near '*, 1)' at line 1
> exec ^
: string
| select row(1, t.*)
< error
| 1064 42000 …near '*)' at line 1
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| select (t.*, 1)
< error
| 1064 42000 …near '*, 1)' at line 1
> exec ^
: string
| select (1, default(a))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select * from t where (a, b) = (1, 2)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t where (a, b) in (select 1, 2)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t where (a, b) > (1, 2)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| Select (1, 1) > (1, 1)
< ok 0 0
> exec ^
: string
| update t set a = (1, 2)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| insert into t values ((1, 2))
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1064 42000 …near '' at line 1
> exec ^
: string
| select row from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select t.row from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select row.a from row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| select `row` from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select 1 row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| select 1 as row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| select a row from t
< error
| 1064 42000 …near 'row from t' at line 1
> exec ^
: string
| select * from row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| select * from t row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| select * from t as row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| create table t (row int)
< error
| 1064 42000 …near 'row int)' at line 1
> exec ^
: string
| create table t (row int, a int)
< error
| 1064 42000 …near 'row int, a int)' at line 1
> exec ^
: string
| alter table t add column row int
< error
| 1064 42000 …near 'row int' at line 1
> exec ^
: string
| alter table t drop column row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| alter table t rename column row to b
< error
| 1064 42000 …near 'row to b' at line 1
> exec ^
: string
| create index row on t (a)
< error
| 1064 42000 …near 'row on t (a)' at line 1
> exec ^
: string
| create index i on t (row)
< error
| 1064 42000 …near 'row)' at line 1
> exec ^
: string
| insert into t (row) values (1)
< error
| 1064 42000 …near 'row) values (1)' at line 1
> exec ^
: string
| insert into t (a, row) values (1, 2)
< error
| 1064 42000 …near 'row) values (1, 2)' at line 1
> exec ^
: string
| replace into t (row) values (1)
< error
| 1064 42000 …near 'row) values (1)' at line 1
> exec ^
: string
| insert into t set row = 1
< error
| 1064 42000 …near 'row = 1' at line 1
> exec ^
: string
| update t set row = 1
< error
| 1064 42000 …near 'row = 1' at line 1
> exec ^
: string
| update t set row := 1
< error
| 1064 42000 …near 'row := 1' at line 1
> exec ^
: string
| update t set t.row = 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| delete from row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| delete from t where row = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| select * from t order by row
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select * from t group by row
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select row, a from t
< error
| 1064 42000 …near ', a from t' at line 1
> exec ^
: string
| select a, row from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select row + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| select row = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| create view row as select 1
< error
| 1064 42000 …near 'row as select 1' at line 1
> exec ^
: string
| drop table row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| truncate table row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| use row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| create database row
< error
| 1064 42000 …near 'row' at line 1
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
| 1064 42000 …near 'declare row int' at line 1
> exec ^
: string
| load data infile 'f' into table t (row)
< error
| 1064 42000 …near 'row)' at line 1
> exec ^
: string
| show columns from row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| grant select on row to u
< error
| 1064 42000 …near 'row to u' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) (partition row)
< error
| 1064 42000 …near 'row)' at line 1
> exec ^
: string
| select row.row from row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| select db.row.a from row
< error
| 1064 42000 …near 'row' at line 1
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
| 1064 42000 …near 'row varchar(10))' at line 1
> exec ^
: string
| create table t (row int) engine = innodb
< error
| 1064 42000 …near 'row int) engine = innodb' at line 1
> exec ^
: string
| select a.row from a
< error
| 1146 42S02 Table 'omni_t_mysql_parser.a' doesn't exist
> exec ^
: string
| select row.row.row from row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| set row = 1
< error
| 1064 42000 …near 'row = 1' at line 1
> exec ^
: string
| insert into t values (row)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values (row, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| insert into t values (row(1,2), row)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select row(row, 1)
< error
| 1064 42000 …near ', 1)' at line 1
> exec ^
: string
| select row(1, 2) as row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| select row(row(row, row), row)
< error
| 1064 42000 …near ', row), row)' at line 1
> exec ^
: string
| select row as row from row as row
< error
| 1064 42000 …near 'as row from row as row' at line 1
> exec ^
: string
| insert into t (row) value (row)
< error
| 1064 42000 …near 'row) value (row)' at line 1
> exec ^
: string
| update t set row = row + 1
< error
| 1064 42000 …near 'row = row + 1' at line 1
> exec ^
: string
| select rows from t
< error
| 1064 42000 …near 'rows from t' at line 1
> exec ^
: string
| create table t (rows int)
< error
| 1064 42000 …near 'rows int)' at line 1
> exec ^
: string
| select * from t rows
< error
| 1064 42000 …near 'rows' at line 1
> exec ^
: string
| select 1 rows
< error
| 1064 42000 …near 'rows' at line 1
> exec ^
: string
| create table rows (a int)
< error
| 1064 42000 …near 'rows (a int)' at line 1
> exec ^
: string
| select row, rows from t
< error
| 1064 42000 …near ', rows from t' at line 1
> exec ^
: string
| insert into t values (ignore)
< error
| 1064 42000 …near 'ignore)' at line 1
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
