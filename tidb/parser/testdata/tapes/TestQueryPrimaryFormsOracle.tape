-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestQueryPrimaryFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| TABLE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t ORDER BY a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t ORDER BY a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t LIMIT 3
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t LIMIT 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t ORDER BY a LIMIT 2,3
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t ORDER BY a LIMIT 2,3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t ORDER BY a LIMIT 1 FOR UPDATE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t ORDER BY a LIMIT 1 FOR UPDATE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t FOR SHARE
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| TABLE t FOR SHARE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t LOCK IN SHARE MODE
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| TABLE t LOCK IN SHARE MODE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t FOR UPDATE NOWAIT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t FOR UPDATE NOWAIT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t FOR UPDATE SKIP LOCKED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t FOR UPDATE SKIP LOCKED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t FOR UPDATE WAIT 5
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t FOR UPDATE WAIT 5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t FOR UPDATE OF t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t FOR UPDATE OF t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/o1.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/o1.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t ORDER BY a INTO OUTFILE '/tmp/o2.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| TABLE t ORDER BY a INTO OUTFILE '/tmp/o2.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t FOR UPDATE INTO OUTFILE '/tmp/a1.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| TABLE t FOR UPDATE INTO OUTFILE '/tmp/a1.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/a3.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| TABLE t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/a3.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/o3.txt' ORDER BY a
< error
| 1064 42000 …line 1 column 40 near "ORDER BY a" 
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/a2.txt' FOR UPDATE
< error
| 1064 42000 …line 1 column 38 near "FOR UPDATE" 
> exec ^
: string
| TABLE t INTO OUTFILE '/tmp/pj.txt' LIMIT 1
< error
| 1064 42000 …line 1 column 40 near "LIMIT 1" 
> exec ^
: string
| TABLE t ORDER BY a ORDER BY b
< error
| 1064 42000 …line 1 column 24 near "ORDER BY b" 
> exec ^
: string
| TABLE t LIMIT 1 LIMIT 2
< error
| 1064 42000 …line 1 column 21 near "LIMIT 2" 
> exec ^
: string
| TABLE t1, t2
< error
| 1064 42000 …line 1 column 9 near ", t2" 
> exec ^
: string
| TABLE t (a)
< error
| 1064 42000 …line 1 column 9 near "(a)" 
> exec ^
: string
| TABLE t WHERE a = 1
< error
| 1064 42000 …line 1 column 13 near "WHERE a = 1" 
> exec ^
: string
| SELECT * FROM TABLE t
< error
| 1064 42000 …line 1 column 19 near "TABLE t" 
> exec ^
: string
| VALUES ROW(1,2)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW()
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1, default)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1, default)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2), ROW(3)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2), ROW(3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) ORDER BY a
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) ORDER BY a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) LIMIT 1
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) ORDER BY a LIMIT 1
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) ORDER BY a LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1) FOR UPDATE
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1) FOR UPDATE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1) FOR UPDATE NOWAIT
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1) FOR UPDATE NOWAIT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1) FOR UPDATE WAIT 5
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1) FOR UPDATE WAIT 5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1) LOCK IN SHARE MODE
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| VALUES ROW(1) LOCK IN SHARE MODE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1) FOR SHARE
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| VALUES ROW(1) FOR SHARE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,-2,3), ROW(5,7,9) INTO OUTFILE '/tmp/v1.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| VALUES ROW(1,-2,3), ROW(5,7,9) INTO OUTFILE '/tmp/v1.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,-2,3), ROW(5,7,9) ORDER BY a INTO OUTFILE '/tmp/v2.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| VALUES ROW(1,-2,3), ROW(5,7,9) ORDER BY a INTO OUTFILE '/tmp/v2.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1) FOR UPDATE INTO OUTFILE '/tmp/a4.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| VALUES ROW(1) FOR UPDATE INTO OUTFILE '/tmp/a4.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) INTO OUTFILE '/tmp/v3.txt' ORDER BY a
< error
| 1064 42000 …line 1 column 48 near "ORDER BY a" 
> exec ^
: string
| VALUES ROW(1) INTO OUTFILE '/tmp/a5.txt' FOR UPDATE
< error
| 1064 42000 …line 1 column 44 near "FOR UPDATE" 
> exec ^
: string
| VALUES ROW(1) ORDER BY a ORDER BY b
< error
| 1064 42000 …line 1 column 30 near "ORDER BY b" 
> exec ^
: string
| VALUES ROW(1) LIMIT 1 LIMIT 2
< error
| 1064 42000 …line 1 column 27 near "LIMIT 2" 
> exec ^
: string
| VALUES (1,2)
< error
| 1064 42000 …line 1 column 8 near "(1,2)" 
> exec ^
: string
| VALUES (1,2), (3,4)
< error
| 1064 42000 …line 1 column 8 near "(1,2), (3,4)" 
> exec ^
: string
| VALUES ROW(1), (2)
< error
| 1064 42000 …line 1 column 16 near "(2)" 
> exec ^
: string
| TABLE t1 UNION TABLE t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION TABLE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION ALL TABLE t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION ALL TABLE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION (TABLE t2)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION (TABLE t2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION SELECT * FROM t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION SELECT * FROM t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION (SELECT c2 FROM t2) ORDER BY c1 LIMIT 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION (SELECT c2 FROM t2) ORDER BY c1 LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 EXCEPT TABLE t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 EXCEPT TABLE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 INTERSECT TABLE t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 INTERSECT TABLE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION TABLE t2 UNION TABLE t3
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION TABLE t2 UNION TABLE t3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION TABLE t2 INTERSECT TABLE t3
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION TABLE t2 INTERSECT TABLE t3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION TABLE t2 ORDER BY c1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION TABLE t2 ORDER BY c1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION TABLE t2 FOR UPDATE
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION TABLE t2 FOR UPDATE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t UNION TABLE t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/a8.txt'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TABLE t UNION TABLE t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/a8.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (TABLE t1) UNION SELECT c2 FROM t2 UNION (TABLE t3) ORDER BY c1 LIMIT 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| (TABLE t1) UNION SELECT c2 FROM t2 UNION (TABLE t3) ORDER BY c1 LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT c1 FROM t1 UNION (TABLE t2) ORDER BY c1 LIMIT 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT c1 FROM t1 UNION (TABLE t2) ORDER BY c1 LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) UNION VALUES ROW(3,4)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) UNION VALUES ROW(3,4)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) UNION ALL VALUES ROW(3,4)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) UNION ALL VALUES ROW(3,4)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) UNION (VALUES ROW(3,4))
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) UNION (VALUES ROW(3,4))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) UNION SELECT 1,2
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) UNION SELECT 1,2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) UNION TABLE t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| VALUES ROW(1,2) UNION TABLE t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TABLE t1 UNION VALUES ROW(1,2)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| TABLE t1 UNION VALUES ROW(1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) EXCEPT VALUES ROW(3,4)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) EXCEPT VALUES ROW(3,4)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1,2) INTERSECT VALUES ROW(3,4)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1,2) INTERSECT VALUES ROW(3,4)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| VALUES ROW(1) UNION VALUES ROW(2) ORDER BY 1 LIMIT 1
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| VALUES ROW(1) UNION VALUES ROW(2) ORDER BY 1 LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (TABLE t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| (TABLE t)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ((TABLE t))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ((TABLE t))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (TABLE t) ORDER BY a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| (TABLE t) ORDER BY a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (TABLE t) LIMIT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| (TABLE t) LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (TABLE t) UNION (TABLE t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| (TABLE t) UNION (TABLE t)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (TABLE t) FOR UPDATE
< error
| 1064 42000 …line 1 column 13 near "FOR UPDATE" 
> exec ^
: string
| (VALUES ROW(1,2))
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| (VALUES ROW(1,2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT (TABLE t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT (TABLE t)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT EXISTS (TABLE t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT EXISTS (TABLE t)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM (TABLE t) x
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM (TABLE t) x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM (VALUES ROW(1,2)) x
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| SELECT * FROM (VALUES ROW(1,2)) x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM ((TABLE t)) x
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM ((TABLE t)) x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM ((VALUES ROW(1))) x
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| SELECT * FROM ((VALUES ROW(1))) x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t1 WHERE c1 IN (TABLE t2)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT * FROM t1 WHERE c1 IN (TABLE t2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t1 WHERE c1 IN (VALUES ROW(1))
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT * FROM t1 WHERE c1 IN (VALUES ROW(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t1 WHERE a = (TABLE t2)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT * FROM t1 WHERE a = (TABLE t2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 WHERE 1 IN ((TABLE t) FOR UPDATE)
< error
| 1064 42000 …line 1 column 34 near "FOR UPDATE)" 
> exec ^
: string
| UPDATE t SET a=1 WHERE b IN (TABLE t1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a=1 WHERE b IN (TABLE t1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t WHERE b IN (TABLE t1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t WHERE b IN (TABLE t1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| WITH cte AS (TABLE t) TABLE cte
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| WITH cte AS (TABLE t) TABLE cte
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| WITH RECURSIVE cte AS (TABLE t) TABLE cte
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| WITH RECURSIVE cte AS (TABLE t) TABLE cte
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| WITH cte AS (TABLE t) SELECT * FROM cte
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| WITH cte AS (TABLE t) SELECT * FROM cte
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| WITH cte AS (SELECT 1) TABLE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| WITH cte AS (SELECT 1) TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| WITH cte AS (SELECT 1) VALUES ROW(1)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| WITH cte AS (SELECT 1) VALUES ROW(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| WITH cte AS (TABLE t) VALUES ROW(1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| WITH cte AS (TABLE t) VALUES ROW(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VIEW v1 AS (TABLE t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VIEW v1 AS (TABLE t)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VIEW v2 AS TABLE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE VIEW v2 AS TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VIEW v3 AS VALUES ROW(1,2)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| CREATE VIEW v3 AS VALUES ROW(1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE VIEW v9 AS TABLE t1 UNION TABLE t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| CREATE VIEW v9 AS TABLE t1 UNION TABLE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE tt AS TABLE t
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| CREATE TABLE tt AS TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE tt2 AS VALUES ROW(1)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| CREATE TABLE tt2 AS VALUES ROW(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE tt3 AS TABLE t UNION TABLE t
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| CREATE TABLE tt3 AS TABLE t UNION TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t1 TABLE t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| INSERT INTO t1 TABLE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| REPLACE INTO t1 TABLE t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| REPLACE INTO t1 TABLE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t1 TABLE t2 ORDER BY a
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| INSERT INTO t1 TABLE t2 ORDER BY a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t1 TABLE t2 LIMIT 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| INSERT INTO t1 TABLE t2 LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t1 TABLE t2 UNION TABLE t3
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| INSERT INTO t1 TABLE t2 UNION TABLE t3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN TABLE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN VALUES ROW(1)
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| EXPLAIN VALUES ROW(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN TABLE t1 UNION TABLE t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| EXPLAIN TABLE t1 UNION TABLE t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR SHARE
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| SELECT * FROM t FOR SHARE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE NOWAIT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE NOWAIT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP LOCKED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP LOCKED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 0
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t, t1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t, t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t.t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF db.t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF db.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF `t`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF `t`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t t1 FOR UPDATE OF t1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t t1 FOR UPDATE OF t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t NOWAIT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t NOWAIT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t SKIP LOCKED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t SKIP LOCKED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t WAIT 5
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t WAIT 5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t, t1 WAIT 5
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t, t1 WAIT 5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR SHARE NOWAIT
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| SELECT * FROM t FOR SHARE NOWAIT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR SHARE SKIP LOCKED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR SHARE SKIP LOCKED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t, t1
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t, t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t NOWAIT
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t NOWAIT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t SKIP LOCKED
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t SKIP LOCKED
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t LIMIT 1 FOR UPDATE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t LIMIT 1 FOR UPDATE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t LIMIT 1 FOR UPDATE WAIT 11
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t LIMIT 1 FOR UPDATE WAIT 11
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t LIMIT 1 FOR UPDATE OF t WAIT 11
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t LIMIT 1 FOR UPDATE OF t WAIT 11
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 UNION SELECT 2 FOR UPDATE
< ok 0 0
> exec ^
: string
| (SELECT 1) FOR UPDATE
< error
| 1064 42000 …line 1 column 14 near "FOR UPDATE" 
> exec ^
: string
| SELECT * FROM t FOR SHARE WAIT 5
< error
| 1064 42000 …line 1 column 30 near "WAIT 5" 
> exec ^
: string
| SELECT * FROM t FOR SHARE OF t WAIT 5
< error
| 1064 42000 …line 1 column 35 near "WAIT 5" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 'a'
< error
| 1064 42000 …line 1 column 35 near "'a'" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT -1
< error
| 1064 42000 …line 1 column 33 near "-1" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5.5
< error
| 1064 42000 …line 1 column 35 near "5.5" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 1e3
< error
| 1064 42000 …line 1 column 35 near "1e3" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 0x10
< error
| 1064 42000 …line 1 column 36 near "0x10" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 010
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 010
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 18446744073709551616
< error
| 1064 42000 …line 1 column 52 near "18446744073709551616" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 99999999999999999999999999
< error
| 1064 42000 …line 1 column 58 near "99999999999999999999999999" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT @v
< error
| 1064 42000 …line 1 column 34 near "@v" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT ?
< error
| 1064 42000 …line 1 column 33 near "?" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT (5)
< error
| 1064 42000 …line 1 column 33 near "(5)" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT OF t
< error
| 1064 42000 …line 1 column 34 near "OF t" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5 OF t
< error
| 1064 42000 …line 1 column 36 near "OF t" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE NOWAIT WAIT 5
< error
| 1064 42000 …line 1 column 38 near "WAIT 5" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5 NOWAIT
< error
| 1064 42000 …line 1 column 40 near "NOWAIT" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP LOCKED WAIT 5
< error
| 1064 42000 …line 1 column 43 near "WAIT 5" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5 SKIP LOCKED
< error
| 1064 42000 …line 1 column 38 near "SKIP LOCKED" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 1 WAIT 2
< error
| 1064 42000 …line 1 column 38 near "WAIT 2" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE NOWAIT NOWAIT
< error
| 1064 42000 …line 1 column 40 near "NOWAIT" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE SKIP LOCKED SKIP LOCKED
< error
| 1064 42000 …line 1 column 43 near "SKIP LOCKED" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE FOR UPDATE
< error
| 1064 42000 …line 1 column 30 near "FOR UPDATE" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE FOR SHARE
< error
| 1064 42000 …line 1 column 30 near "FOR SHARE" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE LOCK IN SHARE MODE
< error
| 1064 42000 …line 1 column 31 near "LOCK IN SHARE MODE" 
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE FOR UPDATE
< error
| 1064 42000 …line 1 column 38 near "FOR UPDATE" 
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE NOWAIT
< error
| 1064 42000 …line 1 column 41 near "NOWAIT" 
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE SKIP LOCKED
< error
| 1064 42000 …line 1 column 39 near "SKIP LOCKED" 
> exec ^
: string
| SELECT * FROM t LOCK IN SHARE MODE WAIT 5
< error
| 1064 42000 …line 1 column 39 near "WAIT 5" 
> exec ^
: string
| SELECT 1 INTO OUTFILE '/tmp/f3.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT 1 INTO OUTFILE '/tmp/f3.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/f3.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/f3.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t ORDER BY a INTO OUTFILE '/tmp/f4.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT * FROM t ORDER BY a INTO OUTFILE '/tmp/f4.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p8.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p8.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/pa.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/pa.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a=1 INTO OUTFILE '/tmp/p2.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT a FROM t WHERE a=1 INTO OUTFILE '/tmp/p2.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t GROUP BY a INTO OUTFILE '/tmp/p4.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT a FROM t GROUP BY a INTO OUTFILE '/tmp/p4.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t HAVING a>1 INTO OUTFILE '/tmp/p6.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT a FROM t HAVING a>1 INTO OUTFILE '/tmp/p6.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t LIMIT 1 INTO OUTFILE '/tmp/pc.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT a FROM t LIMIT 1 INTO OUTFILE '/tmp/pc.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5 INTO OUTFILE '/tmp/f1.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT * FROM t FOR UPDATE WAIT 5 INTO OUTFILE '/tmp/f1.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t WAIT 5 INTO OUTFILE '/tmp/f6.txt'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT * FROM t FOR UPDATE OF t WAIT 5 INTO OUTFILE '/tmp/f6.txt'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 UNION SELECT 2 INTO OUTFILE '/tmp/pg.txt'
< ok 0 0
> exec ^
: string
| SELECT 1 UNION SELECT 2 ORDER BY 1 INTO OUTFILE '/tmp/ph.txt'
< ok 0 0
> exec ^
: string
| SELECT * FROM t INTO OUTFILE '/tmp/f5.txt' ORDER BY a
< error
| 1064 42000 …line 1 column 48 near "ORDER BY a" 
> exec ^
: string
| SELECT a FROM t INTO OUTFILE '/tmp/p7.txt' LIMIT 1
< error
| 1064 42000 …line 1 column 48 near "LIMIT 1" 
> exec ^
: string
| SELECT * FROM t INTO OUTFILE '/tmp/f2.txt' FOR UPDATE
< error
| 1064 42000 …line 1 column 46 near "FOR UPDATE" 
> exec ^
: string
| SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p9.txt' FOR UPDATE
< error
| 1064 42000 …line 1 column 65 near "FOR UPDATE" 
> exec ^
: string
| SELECT a FROM t ORDER BY a INTO OUTFILE '/tmp/pb.txt' LIMIT 1
< error
| 1064 42000 …line 1 column 59 near "LIMIT 1" 
> exec ^
: string
| SELECT a FROM t INTO OUTFILE '/tmp/pe.txt' INTO OUTFILE '/tmp/pf.txt'
< error
| 1064 42000 …line 1 column 47 near "INTO OUTFILE '/tmp/pf.txt'" 
> exec ^
: string
| SELECT a FROM t INTO OUTFILE '/tmp/p3.txt' WHERE a=1
< error
| 1064 42000 …line 1 column 48 near "WHERE a=1" 
> exec ^
: string
| SELECT a FROM t INTO DUMPFILE '/tmp/pd.txt'
< error
| 1064 42000 …line 1 column 29 near "DUMPFILE '/tmp/pd.txt'" 
> exec ^
: string
| SELECT a FROM t INTO @v
< error
| 1064 42000 …line 1 column 23 near "@v" 
> exec ^
: string
| (SELECT 1) INTO OUTFILE '/tmp/pi.txt'
< error
| 1064 42000 …line 1 column 15 near "INTO OUTFILE '/tmp/pi.txt'" 
> exec ^
: string
| SELECT * FROM t wait
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t wait
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT wait FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT wait FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS wait
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS wait
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t wait WHERE wait.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t wait WHERE wait.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t wait WHERE wait.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t wait WHERE wait.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t wait SET wait.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t wait SET wait.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS wait
< ok 0 0
> exec ^
: string
| SELECT * FROM t locked
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t locked
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t nowait
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t nowait
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS locked WHERE locked.a=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t AS locked WHERE locked.a=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t locked WHERE locked.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t locked WHERE locked.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t nowait WHERE nowait.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t nowait WHERE nowait.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE nowait (locked int, wait int, `row` int)
< ok 0 0
