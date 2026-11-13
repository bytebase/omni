-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestUtility_OracleDifferential default#1
> exec
| USE information_schema
< ok - -
> exec
| USE hive.finance
< error
| CATALOG_NOT_FOUND 44 0:0
| Catalog 'hive' not found
> exec
| USE "My Catalog"."My Schema"
< error
| CATALOG_NOT_FOUND 44 0:0
| Catalog 'my catalog' not found
> exec
| SET SESSION query_max_run_time = '10m'
< ok - -
> exec
| SET SESSION example.incremental_refresh_enabled = false
< error
| CATALOG_NOT_FOUND 44 1:51
| line 1:51: Catalog 'example' not found
> exec
| SET SESSION foo = 42
< error
| INVALID_SESSION_PROPERTY 14 1:19
| line 1:19: Session property 'foo' does not exist
> exec
| SET SESSION cat.prop = ARRAY[1, 2]
< error
| CATALOG_NOT_FOUND 44 1:24
| line 1:24: Catalog 'cat' not found
> exec
| SET SESSION a.b.c.d = 1
< error
| INVALID_SESSION_PROPERTY 14 0:0
| Invalid session property 'a.b.c.d'
> exec
| SET SESSION AUTHORIZATION = 1
< error
| INVALID_SESSION_PROPERTY 14 1:29
| line 1:29: Session property 'authorization' does not exist
> exec
| SET SESSION AUTHORIZATION.foo = 1
< error
| CATALOG_NOT_FOUND 44 1:33
| line 1:33: Catalog 'authorization' not found
> exec
| RESET SESSION query_max_run_time
< ok - -
> exec
| RESET SESSION hive.optimized_reader_enabled
< error
| CATALOG_NOT_FOUND 44 1:1
| line 1:1: Catalog 'hive' not found
> exec
| RESET SESSION a.b.c.d
< error
| INVALID_SESSION_PROPERTY 14 1:1
| line 1:1: Invalid session property 'a.b.c.d'
> exec
| RESET SESSION AUTHORIZATION.foo
< error
| CATALOG_NOT_FOUND 44 1:1
| line 1:1: Catalog 'authorization' not found
> exec
| SET SESSION AUTHORIZATION 'John'
< error
| NOT_SUPPORTED 13 0:0
| SET SESSION AUTHORIZATION not supported by client
> exec
| SET SESSION AUTHORIZATION "John"
< error
| NOT_SUPPORTED 13 0:0
| SET SESSION AUTHORIZATION not supported by client
> exec
| SET SESSION AUTHORIZATION John
< error
| NOT_SUPPORTED 13 0:0
| SET SESSION AUTHORIZATION not supported by client
> exec
| RESET SESSION AUTHORIZATION
< error
| NOT_SUPPORTED 13 0:0
| RESET SESSION AUTHORIZATION not supported by client
> exec
| SET ROLE admin
< error
| ROLE_NOT_FOUND 48 1:1
| line 1:1: Role 'admin' does not exist
> exec
| SET ROLE ALL
< ok - -
> exec
| SET ROLE NONE
< ok - -
> exec
| SET ROLE analyst IN hive
< error
| CATALOG_NOT_FOUND 44 1:1
| line 1:1: Catalog 'hive' not found
> exec
| SET PATH example.system
< error
| NOT_SUPPORTED 13 0:0
| SET PATH not supported by client
> exec
| SET PATH a, b.c, d
< error
| NOT_SUPPORTED 13 0:0
| SET PATH not supported by client
> exec
| SET TIME ZONE LOCAL
< ok - -
> exec
| SET TIME ZONE '-08:00'
< ok - -
> exec
| SET TIME ZONE INTERVAL '10' HOUR
< ok - -
> exec
| SET TIME ZONE INTERVAL -'08:00' HOUR TO MINUTE
< ok - -
> exec
| SET TIME ZONE 'America/Los_Angeles'
< ok - -
> exec
| SET TIME ZONE concat_ws('/', 'America', 'Los_Angeles')
< ok - -
> exec
| SHOW GRANTS
< ok - -
> exec
| SHOW GRANTS ON TABLE orders
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'Optional[orders]' does not exist
> exec
| SHOW GRANTS ON orders
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'Optional[orders]' does not exist
> exec
| SHOW CREATE TABLE sf1.orders
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.sf1.orders' does not exist
> exec
| SHOW CREATE SCHEMA hive.web
< error
| SCHEMA_NOT_FOUND 45 1:1
| line 1:1: Schema 'hive.web' does not exist
> exec
| SHOW CREATE VIEW my_view
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: View 'memory.default.my_view' does not exist
> exec
| SHOW CREATE MATERIALIZED VIEW my_mv
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Materialized view 'memory.default.my_mv' does not exist
> exec
| SHOW CREATE FUNCTION example.default.meaning_of_life
< error
| CATALOG_NOT_FOUND 44 0:0
| Catalog 'example' not found
> exec
| SHOW CREATE FUNCTION f
< error
| NOT_SUPPORTED 13 1:1
| line 1:1: Catalog and schema must be specified when function schema is not configured
> exec
| SHOW TABLES
< ok - -
> exec
| SHOW TABLES FROM tpch.tiny
< ok - -
> exec
| SHOW TABLES IN tpch.tiny
< ok - -
> exec
| SHOW TABLES FROM tpch.tiny LIKE 'p%'
< ok - -
> exec
| SHOW TABLES LIKE 'p%' ESCAPE '#'
< ok - -
> exec
| SHOW SCHEMAS
< ok - -
> exec
| SHOW SCHEMAS FROM tpch
< ok - -
> exec
| SHOW SCHEMAS IN tpch
< ok - -
> exec
| SHOW SCHEMAS FROM tpch LIKE '__3%'
< ok - -
> exec
| SHOW CATALOGS
< ok - -
> exec
| SHOW CATALOGS LIKE 't%'
< ok - -
> exec
| SHOW COLUMNS FROM nation
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.nation' does not exist
> exec
| SHOW COLUMNS IN nation
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.nation' does not exist
> exec
| SHOW COLUMNS FROM tpch.sf1.nation
< ok - -
> exec
| SHOW COLUMNS FROM nation LIKE '%key'
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.nation' does not exist
> exec
| SHOW STATS FOR nation
< error
| TABLE_NOT_FOUND 46 0:0
| Table 'memory.default.nation' does not exist
> exec
| SHOW STATS FOR (SELECT * FROM nation WHERE regionkey = 1)
< error
| TABLE_NOT_FOUND 46 1:31
| line 1:31: Table 'memory.default.nation' does not exist
> exec
| SHOW ROLES
< ok - -
> exec
| SHOW ROLES FROM hive
< error
| CATALOG_NOT_FOUND 44 1:1
| line 1:1: Catalog 'hive' not found
> exec
| SHOW CURRENT ROLES
< ok - -
> exec
| SHOW CURRENT ROLES FROM hive
< error
| CATALOG_NOT_FOUND 44 1:1
| line 1:1: Catalog 'hive' not found
> exec
| SHOW ROLE GRANTS
< ok - -
> exec
| SHOW ROLE GRANTS FROM hive
< error
| CATALOG_NOT_FOUND 44 1:1
| line 1:1: Catalog 'hive' not found
> exec
| SHOW FUNCTIONS
< ok - -
> exec
| SHOW FUNCTIONS FROM example.default
< ok - -
> exec
| SHOW FUNCTIONS LIKE 'array%'
< ok - -
> exec
| SHOW SESSION
< ok - -
> exec
| SHOW SESSION LIKE 'query%'
< ok - -
> exec
| DESCRIBE nation
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.nation' does not exist
> exec
| DESC nation
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.nation' does not exist
> exec
| DESCRIBE tpch.sf1.nation
< ok - -
> exec
| DESCRIBE input
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.input' does not exist
> exec
| DESCRIBE output
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.output' does not exist
> exec
| DESC input
< error
| TABLE_NOT_FOUND 46 1:1
| line 1:1: Table 'memory.default.input' does not exist
> exec
| CALL test(123, 'apple')
< error
| PROCEDURE_NOT_FOUND 29 0:0
| Procedure not registered: default.test
> exec
| CALL test(name => 'apple', id => 123)
< error
| PROCEDURE_NOT_FOUND 29 0:0
| Procedure not registered: default.test
> exec
| CALL p()
< error
| PROCEDURE_NOT_FOUND 29 0:0
| Procedure not registered: default.p
> exec
| CALL s.p()
< error
| PROCEDURE_NOT_FOUND 29 0:0
| Procedure not registered: s.p
> exec
| CALL c.s.p()
< error
| CATALOG_NOT_FOUND 44 1:1
| line 1:1: Catalog 'c' not found
> exec
| CALL system.runtime.kill_query(query_id => '20210101_000000_00000_abcde')
< error
| NOT_FOUND 5 0:0
| Target query not found: 20210101_000000_00000_abcde
> exec
| CALL f(name => 1, 2)
< error
| PROCEDURE_NOT_FOUND 29 0:0
| Procedure not registered: default.f
> exec
| EXPLAIN SHOW TABLES
< ok - -
> exec
| EXPLAIN SHOW SCHEMAS FROM tpch
< ok - -
> exec
| EXPLAIN ANALYZE SHOW TABLES
< ok - -
> exec
| EXPLAIN (TYPE LOGICAL) SHOW TABLES
< ok - -
> exec
| SHOW
< error
| SYNTAX_ERROR 1 1:5
| line 1:5: mismatched input '<EOF>'. Expecting: 'BRANCHES', 'CATALOGS', 'COLUMNS', 'CREATE', 'CURRENT', 'FUNCTIONS', 'GRANTS', 'ROLE', 'ROLES', 'SCHEMAS', 'SESSION', 'STATS', 'TABLES'
> exec
| SHOW TABLES FROM
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input '<EOF>'. Expecting: <identifier>
> exec
| SHOW COLUMNS
< error
| SYNTAX_ERROR 1 1:13
| line 1:13: mismatched input '<EOF>'. Expecting: 'FROM', 'IN'
> exec
| SHOW COLUMNS FROM
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input '<EOF>'. Expecting: <identifier>
> exec
| SHOW STATS FOR
< error
| SYNTAX_ERROR 1 1:15
| line 1:15: mismatched input '<EOF>'. Expecting: '(', <identifier>
> exec
| SHOW STATS FOR (nation)
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input 'nation'. Expecting: 'WITH', <query>
> exec
| SHOW STATS FOR (garbage tokens here)
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input 'garbage'. Expecting: 'WITH', <query>
> exec
| SHOW STATS FOR ()
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input ')'. Expecting: 'WITH', <query>
> exec
| SHOW SESSION LIKE
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input '<EOF>'. Expecting: <string>
> exec
| SHOW CREATE TABLE
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input '<EOF>'. Expecting: <identifier>
> exec
| SHOW CREATE
< error
| SYNTAX_ERROR 1 1:12
| line 1:12: mismatched input '<EOF>'. Expecting: 'FUNCTION', 'MATERIALIZED', 'SCHEMA', 'TABLE', 'VIEW'
> exec
| DESCRIBE
< error
| SYNTAX_ERROR 1 1:9
| line 1:9: mismatched input '<EOF>'. Expecting: 'INPUT', 'OUTPUT', <identifier>
> exec
| USE
< error
| SYNTAX_ERROR 1 1:4
| line 1:4: mismatched input '<EOF>'. Expecting: <identifier>
> exec
| SET
< error
| SYNTAX_ERROR 1 1:4
| line 1:4: mismatched input '<EOF>'. Expecting: 'PATH', 'ROLE', 'SESSION', 'TIME'
> exec
| SET SESSION = 'x'
< error
| SYNTAX_ERROR 1 1:13
| line 1:13: mismatched input '='. Expecting: 'AUTHORIZATION', <identifier>
> exec
| SET SESSION foo
< error
| SYNTAX_ERROR 1 1:16
| line 1:16: mismatched input '<EOF>'. Expecting: '.', '='
> exec
| SET TIME ZONE
< error
| SYNTAX_ERROR 1 1:14
| line 1:14: mismatched input '<EOF>'. Expecting: 'LOCAL', <expression>
> exec
| SET ROLE
< error
| SYNTAX_ERROR 1 1:9
| line 1:9: mismatched input '<EOF>'. Expecting: 'ALL', 'NONE', <identifier>
> exec
| SET PATH
< error
| SYNTAX_ERROR 1 1:9
| line 1:9: mismatched input '<EOF>'. Expecting: <identifier>
> exec
| RESET
< error
| SYNTAX_ERROR 1 1:6
| line 1:6: mismatched input '<EOF>'. Expecting: 'SESSION'
> exec
| RESET SESSION
< error
| SYNTAX_ERROR 1 1:14
| line 1:14: mismatched input '<EOF>'. Expecting: 'AUTHORIZATION', <identifier>
> exec
| CALL
< error
| SYNTAX_ERROR 1 1:5
| line 1:5: mismatched input '<EOF>'. Expecting: <identifier>
> exec
| CALL f
< error
| SYNTAX_ERROR 1 1:7
| line 1:7: mismatched input '<EOF>'. Expecting: '(', '.'
> exec
| SHOW ROLE
< error
| SYNTAX_ERROR 1 1:10
| line 1:10: mismatched input '<EOF>'. Expecting: 'GRANTS'
> exec
| SHOW CURRENT ROLE GRANTS
< error
| SYNTAX_ERROR 1 1:14
| line 1:14: mismatched input 'ROLE'. Expecting: 'ROLES'
> exec
| SHOW SCHEMAS FROM c.x
< error
| SYNTAX_ERROR 1 1:20
| line 1:20: mismatched input '.'. Expecting: 'LIKE', <EOF>
> exec
| CALL c.s.p.x()
< error
| SYNTAX_ERROR 1 0:0
| Too many dots in table name: c.s.p.x
> exec
| SHOW CREATE TABLE a.b.c.d
< error
| SYNTAX_ERROR 1 0:0
| Too many dots in table name: a.b.c.d
> exec
| SHOW CREATE FUNCTION a.b.c.d
< error
| SYNTAX_ERROR 1 1:1
| line 1:1: Too many dots in function name: a.b.c.d
> exec
| SHOW COLUMNS FROM a.b.c.d
< error
| SYNTAX_ERROR 1 0:0
| Too many dots in table name: a.b.c.d
> exec
| SHOW STATS FOR a.b.c.d
< error
| SYNTAX_ERROR 1 0:0
| Too many dots in table name: a.b.c.d
> exec
| SHOW GRANTS ON TABLE a.b.c.d
< error
| SYNTAX_ERROR 1 0:0
| Too many dots in table name: a.b.c.d
> exec
| DESCRIBE a.b.c.d
< error
| SYNTAX_ERROR 1 0:0
| Too many dots in table name: a.b.c.d
> exec
| SHOW TABLES FROM a.b.c
< error
| SYNTAX_ERROR 1 1:1
| line 1:1: Too many parts in schema name: a.b.c
> exec
| SHOW FUNCTIONS FROM a.b.c
< error
| SYNTAX_ERROR 1 1:1
| line 1:1: Too many parts in schema name: a.b.c
> exec
| SET PATH a.b.c
< error
| SYNTAX_ERROR 1 1:13
| line 1:13: mismatched input '.'. Expecting: ',', <EOF>
> exec
| DESC INPUT a
< error
| SYNTAX_ERROR 1 1:12
| line 1:12: mismatched input 'a'. Expecting: '.', <EOF>
> exec
| SHOW TABLES extra
< error
| SYNTAX_ERROR 1 1:13
| line 1:13: mismatched input 'extra'. Expecting: 'FROM', 'IN', 'LIKE', <EOF>
> exec
| USE a b
< error
| SYNTAX_ERROR 1 1:7
| line 1:7: mismatched input 'b'. Expecting: <EOF>
> exec
| SET ROLE ALL x
< error
| SYNTAX_ERROR 1 1:14
| line 1:14: mismatched input 'x'. Expecting: 'IN', <EOF>
> exec
| SHOW STATS FOR (SELECT 1) x
< error
| SYNTAX_ERROR 1 1:27
| line 1:27: mismatched input 'x'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| DESCRIBE a b
< error
| SYNTAX_ERROR 1 1:12
| line 1:12: mismatched input 'b'. Expecting: '.', <EOF>
> exec
| SET TIME ZONE LOCAL x
< error
| SYNTAX_ERROR 1 1:21
| line 1:21: mismatched input 'x'. Expecting: <EOF>
> exec
| CALL p() x
< error
| SYNTAX_ERROR 1 1:10
| line 1:10: mismatched input 'x'. Expecting: <EOF>
> exec
| EXPLAIN ()
< error
| SYNTAX_ERROR 1 1:10
| line 1:10: mismatched input ')'. Expecting: '(', 'FORMAT', 'SELECT', 'TABLE', 'TYPE', 'VALUES'
> exec
| EXPLAIN (TYPE) SELECT 1
< error
| SYNTAX_ERROR 1 1:14
| line 1:14: mismatched input ')'. Expecting: 'DISTRIBUTED', 'IO', 'LOGICAL', 'VALIDATE'
> exec
| EXPLAIN (TYPE FOO) SELECT 1
< error
| SYNTAX_ERROR 1 1:15
| line 1:15: mismatched input 'FOO'. Expecting: 'DISTRIBUTED', 'IO', 'LOGICAL', 'VALIDATE'
> exec
| EXPLAIN (FORMAT BAR) SELECT 1
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input 'BAR'. Expecting: 'GRAPHVIZ', 'JSON', 'TEXT'
> exec
| EXPLAIN x
< error
| SYNTAX_ERROR 1 1:9
| line 1:9: mismatched input 'x'. Expecting: '(', 'ALTER', 'ANALYZE', 'CALL', 'COMMENT', 'COMMIT', 'CREATE', 'DEALLOCATE', 'DELETE', 'DENY', 'DESC', 'DESCRIBE', 'DROP', 'EXECUTE', 'EXPLAIN', 'GRANT', 'INSERT', 'MERGE', 'PREPARE', 'REFRESH', 'RESET', 'REVOKE', 'ROLLBACK', 'SET', 'SHOW', 'START', 'TRUNCATE', 'UPDATE', 'USE', 'WITH', <query>
> exec
| EXPLAIN ANALYZE
< error
| SYNTAX_ERROR 1 1:16
| line 1:16: mismatched input '<EOF>'. Expecting: 'ALTER', 'ANALYZE', 'CALL', 'COMMENT', 'COMMIT', 'CREATE', 'DEALLOCATE', 'DELETE', 'DENY', 'DESC', 'DESCRIBE', 'DROP', 'EXECUTE', 'EXPLAIN', 'GRANT', 'INSERT', 'MERGE', 'PREPARE', 'REFRESH', 'RESET', 'REVOKE', 'ROLLBACK', 'SET', 'SHOW', 'START', 'TRUNCATE', 'UPDATE', 'USE', 'VERBOSE', 'WITH', <identifier>, <query>
> exec
| EXPLAIN SELECT 1
< ok - -
> exec
| EXPLAIN (TYPE LOGICAL) SELECT regionkey FROM nation
< error
| TABLE_NOT_FOUND 46 1:46
| line 1:46: Table 'memory.default.nation' does not exist
> exec
| EXPLAIN (TYPE LOGICAL, FORMAT JSON) SELECT 1
< ok - -
> exec
| EXPLAIN (TYPE DISTRIBUTED) SELECT 1
< ok - -
> exec
| EXPLAIN (TYPE VALIDATE) SELECT 1
< ok - -
> exec
| EXPLAIN (TYPE IO, FORMAT JSON) SELECT 1
< ok - -
> exec
| EXPLAIN (FORMAT TEXT) SELECT 1
< ok - -
> exec
| EXPLAIN (FORMAT GRAPHVIZ) SELECT 1
< ok - -
> exec
| EXPLAIN ANALYZE SELECT 1
< ok - -
> exec
| EXPLAIN ANALYZE VERBOSE SELECT 1
< ok - -
> exec
| EXPLAIN INSERT INTO t VALUES (1)
< error
| TABLE_NOT_FOUND 46 1:9
| line 1:9: Table 'memory.default.t' does not exist
> exec
| EXPLAIN CREATE TABLE t AS SELECT 1 AS x
< ok - -
