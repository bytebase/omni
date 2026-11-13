-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestRowFormsMatchEngine default#1
> exec
| SET @omni_sweep = ?
: string
| select row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select row(1,2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row(1,2,3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select row(1,2,3)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row((1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(row(1,2),3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select row(row(1,2),3)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row(1) = row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') = row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1,2) = row(1,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select row(1,2) in (row(1,2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select 1 in (row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '))' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1) in (row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') in (row(1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1,2) = (1,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select * from t where row(1) = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t where row(a,b) = (1,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where row(a,b) = (1,2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where row(a) in (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') in (select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1) from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') from t' at line 1
> exec
| SET @omni_sweep = ?
: string
| do row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| set @a = row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a = row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(1,2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(1,2))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(row(1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(row(1,2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(row(1,2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int default (row(1)))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int check (row(1)))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')))' at line 1
> exec
| SET @omni_sweep = ?
: string
| call p(row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '))' at line 1
> exec
| SET @omni_sweep = ?
: string
| select f(row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '))' at line 1
> exec
| SET @omni_sweep = ?
: string
| select count(row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '))' at line 1
> exec
| SET @omni_sweep = ?
: string
| select if(row(1), 1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '), 1, 2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select case row(1) when 1 then 1 end
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') when 1 then 1 end' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1) is null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') is null' at line 1
> exec
| SET @omni_sweep = ?
: string
| select not row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row (1,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select row (1,2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| create table row (a int)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select row(1) + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') + 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| create procedure row() select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create function row() returns int return 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select row(default)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(default, 1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ' 1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1, 2) from t where row(1, 2) in (select 1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select row(1, 2) from t where row(1, 2) in (select 1, 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row(*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1,,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(t.*, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*, 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1, t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(1, default(a))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select row(1, default(a))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select (1) = row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (1,2) = (1,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select (1,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select (1,2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select ((1,2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select ((1,2))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select (1, 2), (3, 4)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select (1, 2), (3, 4)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select (a, b) in ((1, 2), (3, 4)) from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select (a, b) in ((1, 2), (3, 4)) from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select (1,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select ((select 1), 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select ((select 1), 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select (1, (select 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select (1, (select 2))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select (select 1, 2) = (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select (1, 2) in (select 1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select (1, *)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (*, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*, 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (1, default)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (default, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ' 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (1, t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (t.*, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*, 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (1, default(a))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select (1, default(a))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where (a, b) = (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where (a, b) = (1, 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where (a, b) in (select 1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where (a, b) in (select 1, 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where (a, b) > (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where (a, b) > (1, 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| Select (1, 1) > (1, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| update t set a = (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = (1, 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values ((1, 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values ((1, 2))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select ((1, 2), 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select ((1, 2), 3)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select ((1, 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select ((1, 2))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select ((1, 2)) = (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select (((1, 2)))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select (((1, 2)))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'row' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select row from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select t.row from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select t.row from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row.a from row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.row' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select row.a from row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select `row` from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select `row` from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select 1 row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select 1 as row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select a row from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select a row from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.row' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t as row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t as row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| create table t (row int)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (row int, a int)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t add column row int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t drop column row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| alter table t rename column row to b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create index row on t (a)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create index i on t (row)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| insert into t (row) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (row) values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (a, row) values (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (a, row) values (1, 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t (row) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t (row) values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set row = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set row = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set row = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set row = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set row := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set row := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set t.row = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set t.row = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| delete from row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.row' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| delete from row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| delete from t where row = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| delete from t where row = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t order by row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t order by row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t group by row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t group by row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row, a from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select row, a from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select a, row from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select a, row from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'row' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select row + 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'row' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select row = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| create view row as select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| drop table row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| truncate table row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| use row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create database row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set @row = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select @row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| declare row int
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'declare row int' at line 1
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t (row)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t (row)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| show columns from row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| grant select on row to u
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash(a) (partition row)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select row.row from row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.row' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select row.row from row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select db.row.a from row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.row' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select db.row.a from row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select count(*) over (rows current row)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select count(*) over (rows between current row and unbounded following)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (row varchar(10))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| create table t (row int) engine = innodb
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| select a.row from a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select a.row from a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row.row.row from row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.row' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select row.row.row from row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| set row = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1193 HY000 Unknown system variable 'row'
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row, 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(1,2), row)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(1,2), row)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row(row, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'row' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select row(row, 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row(1, 2) as row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1241 21000 Operand should contain 1 column(s)
> exec
| SET @omni_sweep = ?
: string
| select row(1, 2) as row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row(row(row, row), row)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'row' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select row(row(row, row), row)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select row as row from row as row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.row' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select row as row from row as row
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (row) value (row)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (row) value (row)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set row = row + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set row = row + 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select rows from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'rows from t' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (rows int)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'rows int)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t rows
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'rows' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 rows
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'rows' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table rows (a int)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'rows (a int)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row, rows from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'rows from t' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore, 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1, ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1, ignore)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (default, ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (default, ignore)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t value (ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t value (ignore)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) values (ignore), (default)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) values (ignore), (default)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore) on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore) on duplicate key update a = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t values (ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t values (ignore)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = ignore, b = default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = ignore, b = default
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t set a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t set a = ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a = ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a := ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a := ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = ignore, b = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = ignore, b = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t, u set t.a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t, u set t.a = ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = ignore where b = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = ignore where b = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a = ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a := ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a := ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| execute s using ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| SET @omni_sweep = ?
: string
| execute s using ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| execute s using ignore, default, 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec
| SET @omni_sweep = ?
: string
| execute s using ignore, default, 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| values (ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 4172 HY000 'ignore' is not allowed in this context
> exec
| SET @omni_sweep = ?
: string
| values (ignore)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from (values (ignore)) v
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 4172 HY000 'ignore' is not allowed in this context
> exec
| SET @omni_sweep = ?
: string
| select * from (values (ignore)) v
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore) union values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore) union values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore' at line 1
> exec
| SET @omni_sweep = ?
: string
| select ignore from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore from t' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select ignore + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore + 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 + ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore' at line 1
> exec
| SET @omni_sweep = ?
: string
| select f(ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select row(ignore, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore, 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (ignore + 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1 + ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values ((ignore))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(ignore, 1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore, 1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a = ignore + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a = (ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore)' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 + ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set ignore = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set ignore = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (ignore) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore) values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int default ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table ignore (a int)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore (a int)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (ignore int)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore int)' at line 1
> exec
| SET @omni_sweep = ?
: string
| delete from t where a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore' at line 1
> exec
| SET @omni_sweep = ?
: string
| set @a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore' at line 1
> exec
| SET @omni_sweep = ?
: string
| set autocommit = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore' at line 1
> exec
| SET @omni_sweep = ?
: string
| call p(ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore)' at line 1
> exec
| SET @omni_sweep = ?
: string
| do ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert ignore into t values (ignore)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert ignore into t values (ignore)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update ignore t set a = ignore
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update ignore t set a = ignore
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t ignore index (a)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t ignore index (a)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
