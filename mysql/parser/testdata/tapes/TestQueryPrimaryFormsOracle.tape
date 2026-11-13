-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestQueryPrimaryFormsOracle default#1
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, b INT, c1 INT, c2 INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT, c1 INT, c2 INT)
< ok 0 0
> exec
| CREATE TABLE t3 (a INT, b INT, c1 INT, c2 INT)
< ok 0 0
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| TABLE t
< ok 0 0
> exec ^
: string
| TABLE t ORDER BY a LIMIT 2,3
< ok 0 0
> exec ^
: string
| TABLE t ORDER BY a LIMIT 1 FOR UPDATE
< ok 0 0
> exec ^
: string
| TABLE t FOR SHARE
< ok 0 0
> exec ^
: string
| TABLE t LOCK IN SHARE MODE
< ok 0 0
> exec ^
: string
| TABLE t FOR UPDATE NOWAIT
< ok 0 0
> exec ^
: string
| TABLE t FOR UPDATE SKIP LOCKED
< ok 0 0
> exec ^
: string
| TABLE t FOR UPDATE OF t
< ok 0 0
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/o1.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/o1.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| TABLE t FOR UPDATE INTO OUTFILE '/tmp/a1.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| TABLE t FOR UPDATE INTO OUTFILE '/tmp/a1.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/a2.txt' FOR UPDATE
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/a2.txt' FOR UPDATE
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/o3.txt' ORDER BY a
< error
| 1064 42000 …near 'ORDER BY a' at line 1
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/pj.txt' LIMIT 1
< error
| 1064 42000 …near 'LIMIT 1' at line 1
> exec ^
: string
| TABLE t ORDER BY a ORDER BY b
< error
| 1064 42000 …near 'ORDER BY b' at line 1
> exec ^
: string
| TABLE t LIMIT 1 LIMIT 2
< error
| 1064 42000 …near 'LIMIT 2' at line 1
> exec ^
: string
| TABLE t1, t2
< error
| 1064 42000 …near ', t2' at line 1
> exec ^
: string
| TABLE t (a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| VALUES ROW(1,2)
< ok 0 0
> exec ^
: string
| VALUES ROW(1,2) ORDER BY a LIMIT 1
< error
| 1054 42S22 Unknown column 'a' in 'order clause'
> exec ^
: string
| VALUES ROW(1,2) ORDER BY a LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| VALUES ROW(1) FOR UPDATE
< ok 0 0
> exec ^
: string
| VALUES ROW(1) FOR UPDATE NOWAIT
< ok 0 0
> exec ^
: string
| VALUES ROW(1) LOCK IN SHARE MODE
< ok 0 0
> exec ^
: string
| VALUES ROW(1) FOR SHARE
< ok 0 0
> exec ^
: string
| VALUES ROW(1,-2,3), ROW(5,7,9) INTO OUTFILE '/tmp/v1.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| VALUES ROW(1,-2,3), ROW(5,7,9) INTO OUTFILE '/tmp/v1.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| VALUES ROW(1,-2,3), ROW(5,7,9) ORDER BY a INTO OUTFILE '/tmp/v2.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| VALUES ROW(1,-2,3), ROW(5,7,9) ORDER BY a INTO OUTFILE '/tmp/v2.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| VALUES ROW(1) FOR UPDATE INTO OUTFILE '/tmp/a4.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| VALUES ROW(1) FOR UPDATE INTO OUTFILE '/tmp/a4.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| VALUES ROW(1,2) INTO OUTFILE '/tmp/v3.txt' ORDER BY a
< error
| 1064 42000 …near 'ORDER BY a' at line 1
> exec ^
: string
| VALUES ROW(1) ORDER BY a ORDER BY b
< error
| 1064 42000 …near 'ORDER BY b' at line 1
> exec ^
: string
| VALUES (1,2)
< error
| 1064 42000 …near '(1,2)' at line 1
> exec ^
: string
| VALUES ROW(1), (2)
< error
| 1064 42000 …near '(2)' at line 1
> exec ^
: string
| TABLE t1 UNION TABLE t2
< ok 0 0
> exec ^
: string
| TABLE t1 UNION ALL TABLE t2
< ok 0 0
> exec ^
: string
| TABLE t1 UNION (TABLE t2)
< ok 0 0
> exec ^
: string
| TABLE t1 UNION SELECT * FROM t2
< ok 0 0
> exec ^
: string
| TABLE t1 EXCEPT TABLE t2
< ok 0 0
> exec ^
: string
| TABLE t1 INTERSECT TABLE t2
< ok 0 0
> exec ^
: string
| TABLE t1 UNION TABLE t2 ORDER BY c1
< ok 0 0
> exec ^
: string
| TABLE t1 UNION TABLE t2 FOR UPDATE
< ok 0 0
> exec ^
: string
| VALUES ROW(1,2) UNION VALUES ROW(3,4)
< ok 0 0
> exec ^
: string
| VALUES ROW(1,2) UNION SELECT 1,2
< ok 0 0
> exec ^
: string
| VALUES ROW(1,2) UNION TABLE t1
< error
| 1222 21000 The used SELECT statements have a different number of columns
> exec ^
: string
| VALUES ROW(1,2) UNION TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| TABLE t1 UNION VALUES ROW(1,2)
< error
| 1222 21000 The used SELECT statements have a different number of columns
> exec ^
: string
| TABLE t1 UNION VALUES ROW(1,2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| VALUES ROW(1,2) EXCEPT VALUES ROW(3,4)
< ok 0 0
> exec ^
: string
| VALUES ROW(1,2) INTERSECT VALUES ROW(3,4)
< ok 0 0
> exec ^
: string
| WITH cte AS (SELECT 1) TABLE t
< ok 0 0
> exec ^
: string
| WITH cte AS (SELECT 1) VALUES ROW(1)
< ok 0 0
> exec ^
: string
| WITH cte AS (TABLE t) TABLE cte
< ok 0 0
> exec ^
: string
| WITH RECURSIVE cte AS (TABLE t) TABLE cte
< ok 0 0
> exec ^
: string
| EXPLAIN TABLE t1 UNION TABLE t2
< ok 0 0
> exec ^
: string
| INSERT INTO t1 TABLE t2 ORDER BY a
< ok 0 0
> exec ^
: string
| INSERT INTO t1 TABLE t2 UNION TABLE t3
< ok 0 0
> exec ^
: string
| CREATE TABLE tt3 AS TABLE t UNION TABLE t
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR SHARE
< ok 0 0
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE NOWAIT
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP LOCKED
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t, t1
< error
| 3568 HY000 Unresolved table name `t1` in locking clause.
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t, t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t NOWAIT
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t SKIP LOCKED
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5
< error
| 1064 42000 …near 'WAIT 5' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t WAIT 5
< error
| 1064 42000 …near 'WAIT 5' at line 1
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE SKIP LOCKED
< error
| 1064 42000 …near 'SKIP LOCKED' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE NOWAIT NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP LOCKED SKIP LOCKED
< error
| 1064 42000 …near 'SKIP LOCKED' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/tmp/f3.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT 1 INTO OUTFILE '/tmp/f3.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/f3.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/f3.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p8.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p8.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p9.txt' FOR UPDATE
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p9.txt' FOR UPDATE
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/pa.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/pa.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t GROUP BY a INTO OUTFILE '/tmp/p4.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT a FROM t GROUP BY a INTO OUTFILE '/tmp/p4.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a INTO OUTFILE '/tmp/p1.txt' FROM t
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT a INTO OUTFILE '/tmp/p1.txt' FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a INTO @v FROM t
< ok 0 0
> exec ^
: string
| SELECT a FROM t INTO @v
< ok 0 0
> exec ^
: string
| SELECT a FROM t INTO DUMPFILE '/tmp/pd.txt'
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| SELECT a FROM t INTO DUMPFILE '/tmp/pd.txt'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM t INTO OUTFILE '/tmp/f5.txt' ORDER BY a
< error
| 1064 42000 …near 'ORDER BY a' at line 1
> exec ^
: string
| SELECT a FROM t INTO OUTFILE '/tmp/p7.txt' LIMIT 1
< error
| 1064 42000 …near 'LIMIT 1' at line 1
> exec ^
: string
| SELECT a FROM t ORDER BY a INTO OUTFILE '/tmp/pb.txt' LIMIT 1
< error
| 1064 42000 …near 'LIMIT 1' at line 1
> exec ^
: string
| SELECT a FROM t INTO OUTFILE '/tmp/pe.txt' INTO OUTFILE '/tmp/pf.txt'
< error
| 1064 42000 …near 'INTO OUTFILE '/tmp/pf.txt'' at line 1
> exec ^
: string
| SELECT * FROM t wait
< ok 0 0
> exec ^
: string
| SELECT * FROM t AS wait
< ok 0 0
> exec ^
: string
| DELETE FROM t wait WHERE wait.a = 1
< ok 0 0
> exec ^
: string
| DELETE FROM t locked WHERE locked.a = 1
< ok 0 0
> exec ^
: string
| DELETE FROM t nowait WHERE nowait.a = 1
< ok 0 0
> exec ^
: string
| UPDATE t wait SET wait.a = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE nowait (locked int, wait int, `row` int)
< ok 0 0
