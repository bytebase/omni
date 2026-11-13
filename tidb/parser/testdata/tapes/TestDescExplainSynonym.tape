-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDescExplainSynonym default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| EXPLAIN SELECT 1
< ok 0 0
> exec ^
: string
| DESC SELECT 1
< ok 0 0
> exec ^
: string
| DESCRIBE SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT = 'brief' SELECT max(a) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN FORMAT = 'brief' SELECT max(a) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC FORMAT = 'brief' SELECT max(a) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC FORMAT = 'brief' SELECT max(a) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE FORMAT = 'brief' SELECT max(a) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE FORMAT = 'brief' SELECT max(a) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ANALYZE SELECT * FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN ANALYZE SELECT * FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC ANALYZE SELECT * FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC ANALYZE SELECT * FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE ANALYZE SELECT * FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE ANALYZE SELECT * FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN UPDATE t SET a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN UPDATE t SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC UPDATE t SET a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC UPDATE t SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE UPDATE t SET a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE UPDATE t SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN DELETE FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN DELETE FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC DELETE FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC DELETE FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE DELETE FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE DELETE FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN INSERT INTO t VALUES (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN INSERT INTO t VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC INSERT INTO t VALUES (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC INSERT INTO t VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE INSERT INTO t VALUES (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE INSERT INTO t VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN REPLACE INTO t VALUES (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN REPLACE INTO t VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC REPLACE INTO t VALUES (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC REPLACE INTO t VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE REPLACE INTO t VALUES (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE REPLACE INTO t VALUES (1)
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
| DESC TABLE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE TABLE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN WITH cte AS (SELECT 1) SELECT * FROM cte
< ok 0 0
> exec ^
: string
| DESC WITH cte AS (SELECT 1) SELECT * FROM cte
< ok 0 0
> exec ^
: string
| DESCRIBE WITH cte AS (SELECT 1) SELECT * FROM cte
< ok 0 0
> exec ^
: string
| EXPLAIN ALTER TABLE t ADD COLUMN c int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN ALTER TABLE t ADD COLUMN c int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC ALTER TABLE t ADD COLUMN c int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC ALTER TABLE t ADD COLUMN c int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE ALTER TABLE t ADD COLUMN c int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE ALTER TABLE t ADD COLUMN c int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC FOR CONNECTION 1
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| DESC FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE FOR CONNECTION 1
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| DESCRIBE FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN t c
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| EXPLAIN t c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESC t c
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESC t c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DESCRIBE t c
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE t c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN ANALYZE t
< error
| 1064 42000 …line 1 column 17 near "t" 
> exec ^
: string
| DESC ANALYZE t
< error
| 1064 42000 …line 1 column 14 near "t" 
> exec ^
: string
| DESCRIBE ANALYZE t
< error
| 1064 42000 …line 1 column 18 near "t" 
> exec ^
: string
| EXPLAIN FORMAT = 'row' t
< error
| 1064 42000 …line 1 column 24 near "t" 
> exec ^
: string
| DESC FORMAT = 'row' t
< error
| 1064 42000 …line 1 column 21 near "t" 
> exec ^
: string
| DESCRIBE FORMAT = 'row' t
< error
| 1064 42000 …line 1 column 25 near "t" 
> exec ^
: string
| EXPLAIN CREATE TABLE x (a int)
< error
| 1064 42000 …line 1 column 14 near "CREATE TABLE x (a int)" 
> exec ^
: string
| DESC CREATE TABLE x (a int)
< error
| 1064 42000 …line 1 column 11 near "CREATE TABLE x (a int)" 
> exec ^
: string
| DESCRIBE CREATE TABLE x (a int)
< error
| 1064 42000 …line 1 column 15 near "CREATE TABLE x (a int)" 
