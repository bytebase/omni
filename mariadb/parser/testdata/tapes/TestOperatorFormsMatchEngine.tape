-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOperatorFormsMatchEngine default#1
> exec
| SET @omni_sweep = ?
: string
| select 1 || 2
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
| select a || b from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select a || b from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select 1||2
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
| select 1 ||2
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
| select 1 ||/* c */ 2
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
| select 1 && 2
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
| select a && b from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select a && b from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select 1&&2
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
| select 'a' || 'b'
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
| select concat('a') || 'b'
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
| select (1 || 2)
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
| select f(1 || 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1305 42000 FUNCTION test.f does not exist
> exec
| SET @omni_sweep = ?
: string
| select f(1 || 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where a || b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where a || b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where a=1 || b=2 && c=3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where a=1 || b=2 && c=3
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where a && (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where a && (select 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where a || exists (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where a || exists (select 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t group by a || b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t group by a || b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t having a || b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t having a || b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t order by a && b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t order by a && b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t join u on a || b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t join u on a || b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select case when a || b then 1 end
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select case when a || b then 1 end
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select if(a || b, 1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select if(a || b, 1, 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select count(a || b)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select count(a || b)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 || 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 || 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| delete from t where a || b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| delete from t where a || b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into foo values (1 || 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.foo' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into foo values (1 || 2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into foo values (false || true)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.foo' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into foo values (false || true)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| create view v as select a || b from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| create view v as select a || b from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| create table t (a int, b int as (a || 1))
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
| create table t (a int check (a || 1))
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
| create procedure p() select 1 || 2
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
| show databases where 1 || 2
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
| select @a := 1 || 0
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
| set @a := 1 || 2
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
| select 1 | | 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '| 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 |/* c */| 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '| 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 ||| 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '| 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 || | 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '| 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 ||
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| select || 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 & & 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '& 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 &&& 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '& 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 && & 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '& 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 &&
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| select && 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '&& 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select a between 1 || 2 and 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 2 and 3' at line 1
> exec
| SET @omni_sweep = ?
: string
| select a between 1 && 2 and 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '&& 2 and 3' at line 1
> exec
| SET @omni_sweep = ?
: string
| select a || b over ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'over ()' at line 1
> exec
| SET @omni_sweep = ?
: string
| select position(1 || 2 in 3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 2 in 3)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select position(1 in 2 || 3)
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
| show databases like 'a' || 'b'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 'b'' at line 1
> exec
| SET @omni_sweep = ?
: string
| show databases like 'a' && 'b'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '&& 'b'' at line 1
> exec
| SET @omni_sweep = ?
: string
| select not 1
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
| select not not 1
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
| select not 1 = 2
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
| select not 1 in (1)
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
| select not a is null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select not a is null
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select not exists (select 1)
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
| select 1 and not 1
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
| select 1 or not 1
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
| select 1 xor not 1
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
| select 1 || not 1
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
| select 1 && not 1
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
| select 1 + not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 * not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 | not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 = not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 < not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 like not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 in (1) = not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select - not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select + not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select ~ not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select ! not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select binary not 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select + not exists (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not exists (select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select - not exists (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not exists (select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select a between not 1 and 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1 and 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select a between 1 and not 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select not
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 and not
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| select ! 1
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
| select ! 1 = 2
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
| select ! 1 + 1
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
| select 1 + ! 1
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
| select 1 = ! 1
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
| select ! exists (select 1)
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
| select 1 in (1) in (1)
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
| select 1 in (1) not in (1)
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
| select 1 in (1) like 2
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
| select 1 in (1) between 1 and 2
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
| select 1 in (1) regexp 2
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
| select 1 like 2 like 3
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
| select 1 like 2 not like 3
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
| select 1 like 2 in (1)
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
| select 1 regexp 2 regexp 3
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
| select 1 regexp 2 like 3
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
| select 1 between 2 and 3 in (1)
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
| select 1 in (1) + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 in (1) * 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '* 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 in (1) | 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '| 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 in (1) & 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '& 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 in (1) >> 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '>> 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 like 2 + 3
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
| select 1 like 2 | 3
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
| select 1 not like 2 * 3
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
| select 1 regexp 2 + 3
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
| select 1 regexp 2 | 3
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
| select 1 between 2 and 3 + 4
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
| select 1 between 2 and 3 | 4
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
| select 1 between 2 | 3 and 4
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
| select 1 between 2 + 3 and 4 + 5
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
| select 1 between 2 and 3 between 4 and 5
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
| select 1 between 2 and 3 like 4
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
| select 1 + 1 in (1)
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
| select 1 | 1 in (1)
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
| select 1 + 1 like 2
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
| select 1 + 1 between 1 and 2
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
| select 1 in (1) = 1
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
| select 1 in (1) is null
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
| select 1 in (1) and 2
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
| select 1 in (1) or 2
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
| select 1 like 2 = 3
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
| select 1 between 2 and 3 = 4
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
| select 1 = 1 in (1)
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
| select 1 = 1 like 2
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
| select 1 = 1 between 1 and 2
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
| select 1 = 1 regexp 2
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
| select 1 < 2 = 3
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
| select 1 is null is null
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
| select 1 is null = 1
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
| select 1 is null is true
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
| select 1 is not null is false
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
| select 1 is true is null
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
| select 1 is true is true
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
| select 1 is true = 1
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
| select 1 is false < 1
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
| select 1 is unknown <> 1
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
| select 1 is true and 2
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
| select 1 is true or 2
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
| select 1 is true xor 2
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
| select 1 is true || 2
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
| select 1 is null in (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'in (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 is null + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 is true + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 = 2 is null
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
| select 1 = 2 is true
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
| select not 1 is true
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
| select @a := 1
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
| select @a := @b := 1
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
| select @a := 1 + 1
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
| select @a := 1 = 1
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
| select @a := 1 or 1
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
| select @a := 1 between 0 and 2
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
| select @a := 1 is null
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
| select @a := (select 1)
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
| select @a := 1, @b := 2
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
| select @a.b := 1
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
| select @`a` := 1
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
| select a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 := 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select t.a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select (@a) := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select @@autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select @@session.autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select @a is null := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select @a collate utf8mb4_bin := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select @a + 1 := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select @a := 
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a = b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| set @a = b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| set autocommit = b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| set autocommit = @b := 1
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
| set autocommit := @b := 1
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
| set @a = @b := 1
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
| set @a := @b := 1
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
| load data infile 'x' into table t set a = b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select not @a := 1
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
| select 1 or @a := 1
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
| select 1 xor @a := 1
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
| select 1 and @a := 1
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
| select 1 || @a := 1
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
| select 1 && @a := 1
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
| select 1 = @a := 1
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
| select 1 < @a := 1
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
| select 1 <=> @a := 1
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
| select @b = @a := 1
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
| select @@autocommit = @b := 1
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
| select 1 + @a := 2
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
| select 1 * @a := 2
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
| select 1 | @a := 2
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
| select 1 & @a := 2
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
| select 1 << @a := 2
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
| select 1 ^ @a := 2
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
| select 1 div @a := 2
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
| select - @a := 1
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
| select + @a := 1
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
| select ~ @a := 1
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
| select ! @a := 1
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
| select binary @a := 1
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
| select 1 like @a := 1
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
| select 1 regexp @a := 1
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
| select 1 between 0 and @a := 1
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
| select 1 between @a := 1 and 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 in (@a := 1)
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
| select f(@a := 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1305 42000 FUNCTION test.f does not exist
> exec
| SET @omni_sweep = ?
: string
| select f(@a := 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select (@a := 1)
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
| select case when @a := 1 then @b := 2 else @c := 3 end
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
| select * from t where @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where @a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t where a = 1 and @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t where a = 1 and @a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t order by @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t order by @a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t group by @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t group by @a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (@a := 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (@a := 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = @b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = @b := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = @b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = @b := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = @b := 1 + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = @b := 1 + 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select 1 member of ('[1]')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of ('[1]')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member ('[1]')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '('[1]')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member
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
| select 1 member from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select 1 member from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select a member from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select a member from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select 1 member of ((1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of ((1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member of (- 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of (- 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member of (1 + 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of (1 + 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member of (1 or 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of (1 or 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member of (1 collate utf8mb4_bin)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of (1 collate utf8mb4_bin)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 + 1 member of ('[1]')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of ('[1]')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 = 1 member of ('[1]')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of ('[1]')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member of ('[1]') member of ('[1]')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of ('[1]') member of ('[1]')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member of ('[1]') = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of ('[1]') = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 member of ('[1]') + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'of ('[1]') + 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 not member of ('[1]')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'member of ('[1]')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 'a' like 'b' escape 'c'
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
| select 'a' like 'b' escape 'c' + 1
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
| select 'a' like 'b' escape 1 + 1
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
| select 'a' like 'b' escape concat('c')
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
| select 'a' like 'b' escape - 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec
| SET @omni_sweep = ?
: string
| select 'a' like 'b' escape - 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select 'a' like 'b' escape ('c')
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
| select 'a' like - 1
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
| select 'a' like ! 1
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
| select 'a' like ~ 1
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
| select 'a' like binary 'b'
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
| select 'a' like 'b' collate utf8mb4_bin
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
| select 'a' regexp - 1
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
| select 'a' regexp 'b' || 1
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
| select 'a' regexp 'b' + 1
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
| select 'a' rlike 'b' | 1
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
| select 'a' sounds like 'b'
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
| select 'a' sounds like 'b' + 1
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
| select 'a' sounds like 'b' sounds like 'c'
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
| select 'a' sounds like 'b' = 1
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
| select 1 = 'a' sounds like 'b'
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
| select 'a' not sounds like 'b'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'sounds like 'b'' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 'a' sounds
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| select match (a) against ('x') from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select match (a) against ('x') from t
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select match (a) against (1 or 1) from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 1) from t' at line 1
> exec
| SET @omni_sweep = ?
: string
| select match (a) against (1 in (1)) from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '(1)) from t' at line 1
> exec
| SET @omni_sweep = ?
: string
| select match (a) against (1 = 1) from t
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '= 1) from t' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 = any (select 1)
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
| select 1 = some (select 1)
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
| select 1 = all (select 1)
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
| select 1 <> any (select 1)
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
| select 1 != all (select 1)
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
| select 1 > any (select 1)
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
| select 1 >= some (select 1)
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
| select 1 < all (select 1)
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
| select 1 <= all (select 1)
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
| select 1 <=> any (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '(select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 <=> all (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'all (select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 + any (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '(select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 in any (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'any (select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 like any (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '(select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select any (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '(select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select all (select 1)
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
| select 1 = any
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1054 42S22 Unknown column 'any' in 'SELECT'
> exec
| SET @omni_sweep = ?
: string
| select 1 = any
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select 1 = any ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 = any ((select 1))
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
| select 1 = any (select 1) = 1
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
| select 1 = any (select 1) is null
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
| select 1 = any (select 1) + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 = any (select 1) in (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'in (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select 1 = any (with x as (select 1) select * from x)
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
| select 1 = all ((select 1) union (select 2))
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
| select not 1 = any (select 1)
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
| select 1 from dual where 1 < any (select 2) order by 1
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
| create table t (a int) partition by hash (a)
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
| create table t (a int) partition by hash (a + 1)
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
| create table t (a int) partition by hash (a | 1)
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
| create table t (a int) partition by hash (a & 1)
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
| create table t (a int) partition by hash (a << 1)
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
| create table t (a int) partition by hash (a ^ 1)
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
| create table t (a int) partition by hash (-a)
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
| create table t (a int) partition by hash (~a)
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
| create table t (a int) partition by hash (! a)
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
| create table t (a int) partition by hash (year(a))
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
| create table t (a int) partition by hash ((a or 1))
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
| create table t (a int) partition by hash ((not a))
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
| create table t (a int) partition by hash (a or 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a || 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a and 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'and 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a && 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '&& 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a xor 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'xor 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (not a)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not a)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a is null)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'is null)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a = 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '= 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a < 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '< 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a in (1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'in (1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a between 1 and 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'between 1 and 2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a like 'x')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'like 'x')' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by hash (a regexp 'x')
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'regexp 'x')' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a or 1) (partition p0 values less than (10))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 1) (partition p0 values less than (10))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a + 1) (partition p0 values less than (10))
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
| create table t (a int) partition by range (not a) (partition p0 values less than (10))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not a) (partition p0 values less than (10))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a or 1) (partition p0 values in (10))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 1) (partition p0 values in (10))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a + 1) (partition p0 values in (10))
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
| create table t (a int) partition by range (a) (partition p0 values less than (1 or 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (1 || 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (1 && 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '&& 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (not 1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (1 = 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '= 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (1 is null))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'is null))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (1 + 2))
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
| create table t (a int) partition by range (a) (partition p0 values less than (1 | 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1564 HY000 This partition function is not allowed
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (1 | 2))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1564 HY000 This partition function is not allowed
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than ((1 or 2)))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1564 HY000 This partition function is not allowed
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than ((1 or 2)))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1564 HY000 This partition function is not allowed
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (maxvalue))
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
| create table t (a int) partition by range (a) (partition p0 values less than maxvalue)
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
| create table t (a int) partition by list (a) (partition p0 values in (1 or 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a) (partition p0 values in (1 || 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a) (partition p0 values in (not 1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a) (partition p0 values in (1 + 2))
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
| create table t (a int) partition by range columns (a) (partition p0 values less than (1 or 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range columns (a) (partition p0 values less than (not 1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not 1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range columns (a) (partition p0 values less than (1 is null))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'is null))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range columns (a) (partition p0 values less than (3 in (3, 4)))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'in (3, 4)))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range columns (a) (partition p0 values less than (1 + 2))
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
| create table t (a int) partition by list columns (a) (partition p0 values in (1 or 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int, b int) partition by range (a) subpartition by hash (b or 1) subpartitions 2 (partition p0 values less than (10))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 1) subpartitions 2 (partition p0 values less than (10))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int, b int) partition by range (a) subpartition by hash (b + 1) subpartitions 2 (partition p0 values less than (10))
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
| create table t (a int, b int) partition by range (a) subpartition by hash (not b) subpartitions 2 (partition p0 values less than (10))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'not b) subpartitions 2 (partition p0 values less than (10))' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t partition by hash (a or 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t partition by hash (a + 1)
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
| alter table t add partition (partition p1 values less than (1 or 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t add partition (partition p1 values less than (1 + 2))
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
| alter table t add partition (partition p1 values in (1 or 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| alter table t reorganize partition p0 into (partition p1 values less than (1 or 2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less than)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by range (a) (partition p0 values less (1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '(1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a) (partition p0 values in (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a) (partition p0 values in 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a) (partition p0 values in)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| create table t (a int) partition by list (a) (partition p0 values in ())
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '))' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1, 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1, 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 offset 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 offset 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 5, 18446744073709551615
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 5, 18446744073709551615
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select * from t limit a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| select * from t limit a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| select * from t limit a, b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| select * from t limit a, b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| select * from t limit a offset b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| select * from t limit a offset b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| select * from t limit ``
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: 
> exec
| SET @omni_sweep = ?
: string
| select * from t limit ``
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: 
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1, ``
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: 
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1, ``
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: 
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 offset ``
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: 
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 offset ``
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: 
> exec
| SET @omni_sweep = ?
: string
| select * from t limit t.a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: t
> exec
| SET @omni_sweep = ?
: string
| select * from t limit t.a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: t
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit -1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '-1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 or 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 || 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 || 2, 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 2, 3' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1, 2 || 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 3' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1 offset 2 || 3
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '|| 3' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit @a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 0x10
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '0x10' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1.5
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1.5' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 1e2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1e2' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit 'a'
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ''a'' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from t limit null
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'null' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 limit 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 limit 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 limit a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 limit a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 limit 1 + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1 limit 1 or 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2' at line 1
> exec
| SET @omni_sweep = ?
: string
| delete from t limit 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| delete from t limit 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| delete from t limit a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| delete from t limit a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1327 42000 Undeclared variable: a
> exec
| SET @omni_sweep = ?
: string
| delete from t limit 1 + 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| delete from t limit 1 or 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'or 2' at line 1
