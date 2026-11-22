-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestProbeAndDryRun default#1
> exec
| USE `omni_qe_probe`
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, note VARCHAR(64))
< ok 0 0
> begin
< ok
> query
| EXPLAIN DELETE FROM t WHERE id = 1
< result
< column "id" "CHAR"
< column "estRows" "CHAR"
< column "task" "CHAR"
< column "access object" "CHAR"
< column "operator info" "CHAR"
< row
: bytes
| Delete_2
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
| └─Point_Get_1
: bytes
| 1.00
: bytes
| root
: bytes
| table:t
: bytes
| handle:1, lock
> rollback
< ok
> begin
< ok
> query
| EXPLAIN DELETE FROM t WHERE missing_column = 1
< error
| 1054 42S22 Unknown column 'missing_column' in 'where clause'
> rollback
< ok
> query
| BATCH ON id LIMIT 2 DRY RUN DELETE FROM t WHERE id > 0
< result
< column "split statement examples" "CHAR"
> exec
| START TRANSACTION
< ok 0 0
> exec
| BATCH ON id LIMIT 2 DRY RUN DELETE FROM t WHERE id > 0
< error
| 1105 HY000 non-transactional DML can only run in auto-commit mode. auto-commit:true, inTxn:true
> exec
| ROLLBACK
< ok 0 0
> exec
| START TRANSACTION
< ok 0 0
> exec
| INSERT INTO t VALUES (42, 'probe')
< ok 1 0
> query
| SELECT COUNT(*) FROM t
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 1
> exec
| ROLLBACK
< ok 0 0
> query
| SELECT COUNT(*) FROM t
< result
< column "COUNT(*)" "BIGINT" "int64"
< row
: int64 0
