-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestOracleClassification default#1
> exec
| SELECT 1
< ok - -
> exec
| SELECT name FROM tpch.sf1.nation
< ok - -
> exec
| SHOW CATALOGS
< ok - -
> exec
| EXPLAIN SELECT 1
< ok - -
> exec
| CREATE TABLE memory.default.omni_oracle_probe (x integer)
< ok - -
> exec
| SELECT * FROM tpch.sf1.orders MATCH_RECOGNIZE (PARTITION BY orderkey ORDER BY orderdate MEASURES A.totalprice AS p PATTERN (A+) DEFINE A AS true)
< ok - -
> exec
| SELECT nationkey, count(*) FROM tpch.sf1.nation GROUP BY ROLLUP (nationkey)
< ok - -
> exec
| SELECT filter(ARRAY[1,2,3], x -> x > 1)
< ok - -
> exec
| SELECT * FROM no_such_table_xyz
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.no_such_table_xyz' does not exist
> exec
| SELECT * FROM nocat.nosch.notbl
< error
| CATALOG_NOT_FOUND 44 1:15
| line 1:15: Catalog 'nocat' not found
> exec
| SELECT bogus_col FROM tpch.sf1.nation
< error
| COLUMN_NOT_FOUND 47 1:8
| line 1:8: Column 'bogus_col' cannot be resolved
> exec
| SELCT 1
< error
| SYNTAX_ERROR 1 1:1
| line 1:1: mismatched input 'SELCT'. Expecting: 'ALTER', 'ANALYZE', 'CALL', 'COMMENT', 'COMMIT', 'CREATE', 'DEALLOCATE', 'DELETE', 'DENY', 'DESC', 'DESCRIBE', 'DROP', 'EXECUTE', 'EXPLAIN', 'GRANT', 'INSERT', 'MERGE', 'PREPARE', 'REFRESH', 'RESET', 'REVOKE', 'ROLLBACK', 'SET', 'SHOW', 'START', 'TRUNCATE', 'UPDATE', 'USE', 'WITH', <query>
> exec
| SELECT * FROM
< error
| SYNTAX_ERROR 1 1:14
| line 1:14: mismatched input '<EOF>'. Expecting: '(', 'JSON_TABLE', 'LATERAL', 'NEAREST', 'TABLE', 'UNNEST', <identifier>
> exec
| SELECT (1 +
< error
| SYNTAX_ERROR 1 1:12
| line 1:12: mismatched input '<EOF>'. Expecting: <expression>
> exec
| CREATE TABLE memory.default.t (x intGAR-BAGE !!!)
< error
| SYNTAX_ERROR 1 1:40
| line 1:40: mismatched input '-'. Expecting: ')', ','
