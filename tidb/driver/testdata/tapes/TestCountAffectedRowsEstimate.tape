-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestCountAffectedRowsEstimate default#1
> exec
| USE `omni_estimate_test`
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| CREATE TABLE archive (id INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| INSERT INTO t
| 			WITH RECURSIVE seq (n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM seq WHERE n < 1000)
| 			SELECT n, n % 10 FROM seq
< ok 1000 0
> exec
| ANALYZE TABLE t ALL COLUMNS
< ok 0 0
> query
| EXPLAIN UPDATE t SET c = 0 WHERE c = 1
< result
< column "id" "CHAR"
< column "estRows" "CHAR"
< column "task" "CHAR"
< column "access object" "CHAR"
< column "operator info" "CHAR"
< row
: bytes
| Update_4
: bytes
| N/A
: bytes
| root
: bytes
|
: bytes
| N/A
< row
: bytes
| └─TableReader_8
: bytes
| 100.00
: bytes
| root
: bytes
|
: bytes
| data:Selection_7
< row
: bytes
|   └─Selection_7
: bytes
| 100.00
: bytes
| cop[tikv]
: bytes
|
: bytes
| eq(omni_estimate_test.t.c, 1)
< row
: bytes
|     └─TableFullScan_6
: bytes
| 1000.00
: bytes
| cop[tikv]
: bytes
| table:t
: bytes
| keep order:false, stats:partial[id:allEvicted]
> query
| EXPLAIN DELETE FROM t WHERE c = 1
< result
< column "id" "CHAR"
< column "estRows" "CHAR"
< column "task" "CHAR"
< column "access object" "CHAR"
< column "operator info" "CHAR"
< row
: bytes
| Delete_3
: bytes
| N/A
: bytes
| root
: bytes
|
: bytes
| N/A
< row
: bytes
| └─TableReader_12
: bytes
| 100.00
: bytes
| root
: bytes
|
: bytes
| data:Projection_6
< row
: bytes
|   └─Projection_6
: bytes
| 100.00
: bytes
| cop[tikv]
: bytes
|
: bytes
| omni_estimate_test.t.id
< row
: bytes
|     └─Selection_11
: bytes
| 100.00
: bytes
| cop[tikv]
: bytes
|
: bytes
| eq(omni_estimate_test.t.c, 1)
< row
: bytes
|       └─TableFullScan_10
: bytes
| 1000.00
: bytes
| cop[tikv]
: bytes
| table:t
: bytes
| keep order:false, stats:partial[id:allEvicted]
> query
| EXPLAIN INSERT INTO archive SELECT * FROM t WHERE c = 1
< result
< column "id" "CHAR"
< column "estRows" "CHAR"
< column "task" "CHAR"
< column "access object" "CHAR"
< column "operator info" "CHAR"
< row
: bytes
| Insert_1
: bytes
| N/A
: bytes
| root
: bytes
|
: bytes
| N/A
< row
: bytes
| └─TableReader_9
: bytes
| 100.00
: bytes
| root
: bytes
|
: bytes
| data:Selection_8
< row
: bytes
|   └─Selection_8
: bytes
| 100.00
: bytes
| cop[tikv]
: bytes
|
: bytes
| eq(omni_estimate_test.t.c, 1)
< row
: bytes
|     └─TableFullScan_7
: bytes
| 1000.00
: bytes
| cop[tikv]
: bytes
| table:t
: bytes
| keep order:false, stats:partial[id:allEvicted]
