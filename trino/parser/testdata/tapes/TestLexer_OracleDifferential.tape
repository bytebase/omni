-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestLexer_OracleDifferential default#1
> exec
| SELECT 1
< ok - -
> exec
| SELECT 1, 2.5, 1.5e10, -3, +4 FROM t
< error
| TABLE_NOT_FOUND 46 1:36
| line 1:36: Table 'memory.default.t' does not exist
> exec
| SELECT * FROM tpch.sf1.nation
< ok - -
> exec
| SELECT a, b AS c, t.d, t.* FROM t
< error
| TABLE_NOT_FOUND 46 1:33
| line 1:33: Table 'memory.default.t' does not exist
> exec
| SELECT 'a string', 'it''s escaped', U&'\0041', X'00ff'
< ok - -
> exec
| SELECT a || b AS concatenated FROM t
< error
| TABLE_NOT_FOUND 46 1:36
| line 1:36: Table 'memory.default.t' does not exist
> exec
| SELECT CASE WHEN a > 1 THEN 'hi' ELSE 'lo' END FROM t
< error
| TABLE_NOT_FOUND 46 1:53
| line 1:53: Table 'memory.default.t' does not exist
> exec
| SELECT CAST(a AS varchar), TRY_CAST(b AS bigint) FROM t
< error
| TABLE_NOT_FOUND 46 1:55
| line 1:55: Table 'memory.default.t' does not exist
> exec
| SELECT a <> b, a != b, a <= b, a >= b, a < b, a > b FROM t
< error
| TABLE_NOT_FOUND 46 1:58
| line 1:58: Table 'memory.default.t' does not exist
> exec
| SELECT a + b - c * d / e % f FROM t
< error
| TABLE_NOT_FOUND 46 1:35
| line 1:35: Table 'memory.default.t' does not exist
> exec
| SELECT count(*) FILTER (WHERE a > 0) FROM t
< error
| TABLE_NOT_FOUND 46 1:43
| line 1:43: Table 'memory.default.t' does not exist
> exec
| SELECT "quoted col" FROM "quoted tbl"
< error
| TABLE_NOT_FOUND 46 1:26
| line 1:26: Table 'memory.default."quoted tbl"' does not exist
> exec
| SELECT timestamp '2012-10-31 01:00 UTC' AT TIME ZONE 'America/Los_Angeles'
< ok - -
> exec
| SELECT ARRAY[1, 2, 3], ROW(1, 'a')
< ok - -
> exec
| SELECT filter(ARRAY[1,2,3], x -> x > 1)
< ok - -
> exec
| SELECT transform(a, (k, v) -> k + v) FROM t
< error
| TABLE_NOT_FOUND 46 1:43
| line 1:43: Table 'memory.default.t' does not exist
> exec
| SELECT JSON_OBJECT('k' VALUE 1 RETURNING varchar FORMAT JSON ENCODING UTF8)
< error
| TYPE_MISMATCH 58 1:8
| line 1:8: Cannot output JSON value as varchar using formatting JSON ENCODING UTF8
> exec
| SELECT JSON_QUERY(a, 'lax $.b' RETURNING varchar FORMAT JSON ENCODING UTF16) FROM t
< error
| TABLE_NOT_FOUND 46 1:83
| line 1:83: Table 'memory.default.t' does not exist
> exec
| SELECT a, count(*) FROM t GROUP BY ROLLUP (a)
< error
| TABLE_NOT_FOUND 46 1:25
| line 1:25: Table 'memory.default.t' does not exist
> exec
| SELECT a, count(*) FROM t GROUP BY GROUPING SETS ((a), ())
< error
| TABLE_NOT_FOUND 46 1:25
| line 1:25: Table 'memory.default.t' does not exist
> exec
| SELECT sum(x) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t
< error
| TABLE_NOT_FOUND 46 1:102
| line 1:102: Table 'memory.default.t' does not exist
> exec
| SELECT a FROM t UNION ALL SELECT a FROM u
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.t' does not exist
> exec
| SELECT a FROM t INTERSECT SELECT a FROM u
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.t' does not exist
> exec
| WITH x AS (SELECT 1 AS a) SELECT a FROM x
< ok - -
> exec
| WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE n < 5) SELECT n FROM r
< ok - -
> exec
| SELECT * FROM a JOIN b ON a.id = b.id
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.a' does not exist
> exec
| SELECT * FROM a LEFT OUTER JOIN b USING (id)
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.a' does not exist
> exec
| SELECT * FROM a CROSS JOIN UNNEST(b) WITH ORDINALITY
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.a' does not exist
> exec
| SELECT * FROM TABLE(sequence(1, 10))
< ok - -
> exec
| SELECT * FROM t FOR TIMESTAMP AS OF TIMESTAMP '2023-01-01 00:00:00'
< error
| NOT_SUPPORTED 13 0:0
| This connector does not support versioned tables
> exec
| SELECT * FROM t MATCH_RECOGNIZE (PARTITION BY a ORDER BY b MEASURES c AS m PATTERN (^ A+ B* $) DEFINE A AS true)
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.t' does not exist
> exec
| CREATE TABLE memory.default.t (x integer, y varchar)
< ok - -
> exec
| CREATE TABLE t (x bigint) WITH (format = 'ORC')
< error
| TABLE_ALREADY_EXISTS 50 1:1
| line 1:1: Table 'memory.default.t' already exists
> exec
| CREATE TABLE t AS SELECT 1 AS x
< error
| TABLE_ALREADY_EXISTS 50 1:1
| line 1:1: Destination table 'memory.default.t' already exists
> exec
| DROP TABLE IF EXISTS t
< ok - -
> exec
| ALTER TABLE t ADD COLUMN z double
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.t' does not exist
> exec
| CREATE VIEW v AS SELECT 1 AS a
< ok - -
> exec
| CREATE MATERIALIZED VIEW mv GRACE PERIOD INTERVAL '1' HOUR AS SELECT 1 AS a
< error
| NOT_SUPPORTED 13 1:1
| line 1:1: Catalog 'memory' does not support GRACE PERIOD
> exec
| COMMENT ON TABLE t IS 'a comment'
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table does not exist: memory.default.t
> exec
| INSERT INTO t VALUES (1, 'a'), (2, 'b')
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.t' does not exist
> exec
| DELETE FROM t WHERE a > 1
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.t' does not exist
> exec
| UPDATE t SET a = a + 1 WHERE b = 'x'
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.t' does not exist
> exec
| MERGE INTO t USING u ON t.id = u.id WHEN MATCHED THEN UPDATE SET a = u.a
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.t' does not exist
> exec
| TRUNCATE TABLE t
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.t' does not exist
> exec
| USE catalog.schema
< error
| CATALOG_NOT_FOUND 44 0:0
| Catalog 'catalog' not found
> exec
| SET SESSION query_max_run_time = '1h'
< ok - -
> exec
| SET TIME ZONE LOCAL
< ok - -
> exec
| SHOW TABLES FROM s LIKE 'pattern%'
< error
| SCHEMA_NOT_FOUND 45 1:1
| line 1:1: Schema 's' does not exist
> exec
| EXPLAIN ANALYZE VERBOSE SELECT 1
< ok - -
> exec
| PREPARE my_query FROM SELECT ? FROM t
< ok - -
> exec
| EXECUTE my_query USING 1
< error
| NOT_FOUND 5 0:0
| Prepared statement not found: my_query
> exec
| START TRANSACTION ISOLATION LEVEL SERIALIZABLE
< error
| INCOMPATIBLE_CLIENT 27 0:0
| Client does not support transactions
> exec
| GRANT SELECT ON t TO USER alice WITH GRANT OPTION
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.t' does not exist
> exec
| CALL system.runtime.kill_query(query_id => '20210101_000000_00000_abcde')
< error
| NOT_FOUND 5 0:0
| Target query not found: 20210101_000000_00000_abcde
> exec
| SELECT 1 -- trailing comment
| FROM t
< error
| TABLE_NOT_FOUND 46 2:6
| line 2:6: Table 'memory.default.t' does not exist
> exec
| SELECT /* inline */ 1 FROM t
< error
| TABLE_NOT_FOUND 46 1:28
| line 1:28: Table 'memory.default.t' does not exist
> exec
| SELECT /* a /* not nested */ 1
< ok - -
