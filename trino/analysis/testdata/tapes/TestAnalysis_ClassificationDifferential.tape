-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestAnalysis_ClassificationDifferential default#1
> exec
| DROP TABLE IF EXISTS t
< ok - -
> exec
| DROP TABLE IF EXISTS t2
< ok - -
> exec
| DROP TABLE IF EXISTS t3
< ok - -
> exec
| CREATE TABLE t (a integer, b integer, id integer)
< ok - -
> exec
| CREATE TABLE t2 (a2 integer, b integer, id integer)
< ok - -
> exec
| SELECT a FROM t
< ok - -
> exec
| SELECT * FROM t1 JOIN t2 ON t1.id = t2.id
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.t1' does not exist
> exec
| WITH c AS (SELECT a FROM t) SELECT a FROM c
< ok - -
> exec
| SELECT a FROM t WHERE b IN (SELECT b FROM t2)
< ok - -
> exec
| SELECT a FROM t UNION SELECT a2 FROM t2
< ok - -
> exec
| SELECT * FROM information_schema.tables
< ok - -
> exec
| SHOW TABLES
< ok - -
> exec
| SHOW SCHEMAS
< ok - -
> exec
| DESCRIBE t
< ok - -
> exec
| EXPLAIN SELECT a FROM t
< ok - -
> exec
| EXPLAIN ANALYZE SELECT a FROM t
< ok - -
> exec
| EXPLAIN ANALYZE INSERT INTO t VALUES (1)
< error
| TYPE_MISMATCH 58 1:17
| line 1:17: Insert query has mismatched column types: Table: [integer, integer, integer], Query: [integer]
> exec
| EXPLAIN INSERT INTO t VALUES (1)
< error
| TYPE_MISMATCH 58 1:9
| line 1:9: Insert query has mismatched column types: Table: [integer, integer, integer], Query: [integer]
> exec
| INSERT INTO t VALUES (1)
< error
| TYPE_MISMATCH 58 1:1
| line 1:1: Insert query has mismatched column types: Table: [integer, integer, integer], Query: [integer]
> exec
| DELETE FROM t WHERE a = 1
< error
| NOT_SUPPORTED 13 0:0
| This connector does not support modifying table rows
> exec
| UPDATE t SET a = 2 WHERE a = 1
< error
| NOT_SUPPORTED 13 0:0
| This connector does not support modifying table rows
> exec
| CREATE TABLE t3 (a integer)
< ok - -
> exec
| DROP TABLE t3
< ok - -
