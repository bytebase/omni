-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestLockingReadFormsOracle default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT)
< ok 0 0
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t FOR UPDATE
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
| SELECT * FROM t FOR UPDATE WAIT 5
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 0
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5.5
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 0.5
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 1e3
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 010
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 0x10
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 0b10
< error
| 1064 42000 …near '0b10' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT X'10'
< error
| 1064 42000 …near 'X'10'' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 18446744073709551615
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 99999999999999999999999999
< ok 0 0
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT -5
< error
| 1064 42000 …near '-5' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT @v
< error
| 1064 42000 …near '@v' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT ?
< error
| 1064 42000 …near '?' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT (5)
< error
| 1064 42000 …near '(5)' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 1 WAIT 2
< error
| 1064 42000 …near 'WAIT 2' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE NOWAIT WAIT 5
< error
| 1064 42000 …near 'WAIT 5' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5 NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP LOCKED WAIT 5
< error
| 1064 42000 …near 'WAIT 5' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5 SKIP LOCKED
< error
| 1064 42000 …near 'SKIP LOCKED' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE NOWAIT NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| SELECT * FROM t FOR SHARE
< error
| 1064 42000 …near 'SHARE' at line 1
> exec ^
: string
| SELECT * FROM t FOR SHARE NOWAIT
< error
| 1064 42000 …near 'SHARE NOWAIT' at line 1
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t
< error
| 1064 42000 …near 'SHARE OF t' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t
< error
| 1064 42000 …near 'OF t' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t NOWAIT
< error
| 1064 42000 …near 'OF t NOWAIT' at line 1
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF
< error
| 1064 42000 …near 'OF' at line 1
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE NOWAIT
< ok 0 0
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE SKIP LOCKED
< ok 0 0
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE WAIT 5
< ok 0 0
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE WAIT 5.5
< ok 0 0
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE WAIT 0x10
< ok 0 0
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE WAIT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE NOWAIT NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| SELECT * FROM t LIMIT 1 FOR UPDATE
< ok 0 0
> exec ^
: string
| SELECT * FROM t LIMIT 1 FOR UPDATE WAIT 11
< ok 0 0
> exec ^
: string
| SELECT 1 UNION SELECT 2 FOR UPDATE
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/tmp/f3.txt'
< ok 0 0
> exec ^
: string
| SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/f3.txt'
< ok 0 0
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p8.txt'
< ok 0 0
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p9.txt' FOR UPDATE
< ok 0 0
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/pa.txt'
< ok 0 0
> exec ^
: string
| SELECT a FROM t GROUP BY a INTO OUTFILE '/tmp/p4.txt'
< ok 0 0
> exec ^
: string
| SELECT a INTO OUTFILE '/tmp/p1.txt' FROM t
< ok 0 0
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
< ok 0 0
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
| SELECT * FROM t FOR UPDATE WAIT 5 INTO OUTFILE '/tmp/f1.txt'
< ok 0 0
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
| SELECT * FROM t wait WHERE wait.a = 1
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
| SELECT 1 AS wait
< ok 0 0
> exec ^
: string
| CREATE TABLE nowait (locked int, wait int, `row` int)
< ok 0 0
> exec ^
: string
| TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec ^
: string
| TABLE t ORDER BY a
< error
| 1064 42000 …near 'TABLE t ORDER BY a' at line 1
> exec ^
: string
| VALUES ROW(1,2)
< error
| 1064 42000 …near 'ROW(1,2)' at line 1
> exec ^
: string
| VALUES (1,2)
< ok 0 0
> exec ^
: string
| VALUES (1,2), (3,4)
< ok 0 0
