-- pin: default starrocks/allin1-ubuntu:3.4.10 sha256:8bc803ea7a5bb490c77117a456b32b792ccae698b1f0e47dfcafdac41542566d 3.4.10

== TestDumpGaps/function-definition default#1
> exec
| CREATE FUNCTION dump_gap_fn(INT) RETURNS INT PROPERTIES("type" = "Python", "file" = "file:///data/deploy/starrocks/fe/conf/fe.conf", "symbol" = "f", "input_type" = "INT")
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_gap_function
< result
< column "SCHEMA_NAME" "VARCHAR"
< row
: bytes
| omni_dump_gap_function
> query
| SELECT TABLE_NAME FROM information_schema.materialized_views WHERE TABLE_SCHEMA = ? AND IFNULL(REFRESH_TYPE, '') <> 'ROLLUP' ORDER BY TABLE_NAME
: string
| omni_dump_gap_function
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_gap_function
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
> query
| SHOW CREATE FUNCTION dump_gap_fn(INT)
< error
| 1064 HY000 Getting syntax error at line 1, column 12. Detail message: No viable statement for input 'SHOW CREATE FUNCTION'.
> exec
| DROP FUNCTION dump_gap_fn(INT)
< ok 0 0

== TestDumpGaps/global-function-definition default#1
> exec
| CREATE GLOBAL FUNCTION dump_gap_global(INT) RETURNS INT PROPERTIES("type" = "Python", "file" = "file:///data/deploy/starrocks/fe/conf/fe.conf", "symbol" = "f", "input_type" = "INT");
< ok 0 0
> exec
|
| CREATE VIEW calls_global (g) AS SELECT dump_gap_global(1)
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_gap_global_function
< result
< column "SCHEMA_NAME" "VARCHAR"
< row
: bytes
| omni_dump_gap_global_function
> query
| SELECT TABLE_NAME FROM information_schema.materialized_views WHERE TABLE_SCHEMA = ? AND IFNULL(REFRESH_TYPE, '') <> 'ROLLUP' ORDER BY TABLE_NAME
: string
| omni_dump_gap_global_function
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_gap_global_function
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| calls_global
: bytes
| VIEW
> query
| SHOW CREATE VIEW `omni_dump_gap_global_function`.`calls_global`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| calls_global
: bytes
| CREATE VIEW `calls_global` (`g`) SECURITY NONE AS SELECT dump_gap_global(1) AS `dump_gap_global(1)`;
: bytes
| utf8
: bytes
| utf8_general_ci
> query
| SHOW CREATE GLOBAL FUNCTION dump_gap_global(INT)
< error
| 1064 HY000 Getting syntax error at line 1, column 12. Detail message: No viable statement for input 'SHOW CREATE GLOBAL'.

== TestDumpGaps/rollup default#1
> exec
| CREATE TABLE r (
| 	k INT NOT NULL,
| 	v INT NULL,
| 	d DATETIME NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(k)
| DISTRIBUTED BY HASH(k) BUCKETS 1
| ROLLUP (r1 (v, k))
| PROPERTIES("replication_num" = "1")
< ok 0 0
> query
| DESC `omni_dump_gap_rollup`.`r` ALL
< result
< column "IndexName" "VARCHAR"
< column "IndexKeysType" "VARCHAR"
< column "Field" "VARCHAR"
< column "Type" "VARCHAR"
< column "Null" "VARCHAR"
< column "Key" "VARCHAR"
< column "Default" "VARCHAR"
< column "Extra" "VARCHAR"
< row
: bytes
| r
: bytes
| DUP_KEYS
: bytes
| k
: bytes
| int
: bytes
| NO
: bytes
| true
: null
: bytes
|
< row
: bytes
|
: bytes
|
: bytes
| v
: bytes
| int
: bytes
| YES
: bytes
| false
: null
: bytes
| NONE
< row
: bytes
|
: bytes
|
: bytes
| d
: bytes
| datetime
: bytes
| YES
: bytes
| false
: null
: bytes
| NONE
< row
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| r1
: bytes
| DUP_KEYS
: bytes
| v
: bytes
| int
: bytes
| YES
: bytes
| true
: null
: bytes
| NONE
< row
: bytes
|
: bytes
|
: bytes
| k
: bytes
| int
: bytes
| NO
: bytes
| true
: null
: bytes
|
> query ^
< result
< column "IndexName" "VARCHAR"
< column "IndexKeysType" "VARCHAR"
< column "Field" "VARCHAR"
< column "Type" "VARCHAR"
< column "Null" "VARCHAR"
< column "Key" "VARCHAR"
< column "Default" "VARCHAR"
< column "Extra" "VARCHAR"
< row
: bytes
| r
: bytes
| DUP_KEYS
: bytes
| k
: bytes
| int
: bytes
| NO
: bytes
| true
: null
: bytes
|
< row
: bytes
|
: bytes
|
: bytes
| v
: bytes
| int
: bytes
| YES
: bytes
| false
: null
: bytes
| NONE
< row
: bytes
|
: bytes
|
: bytes
| d
: bytes
| datetime
: bytes
| YES
: bytes
| false
: null
: bytes
| NONE
< row
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| r1
: bytes
| DUP_KEYS
: bytes
| v
: bytes
| int
: bytes
| YES
: bytes
| true
: null
: bytes
| NONE
< row
: bytes
|
: bytes
|
: bytes
| k
: bytes
| int
: bytes
| NO
: bytes
| true
: null
: bytes
|
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_gap_rollup`.`r1`
< error
| 1064 HY000 Getting analyzing error. Detail message: Table r1 is not found.
> query
| SHOW CREATE MATERIALIZED VIEW `r1` ON `omni_dump_gap_rollup`.`r`
< error
| 1064 HY000 Getting syntax error at line 1, column 35. Detail message: Unexpected input 'ON', the most similar input is {<EOF>, ';'}.
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_gap_rollup
< result
< column "SCHEMA_NAME" "VARCHAR"
< row
: bytes
| omni_dump_gap_rollup
> query
| SELECT TABLE_NAME FROM information_schema.materialized_views WHERE TABLE_SCHEMA = ? AND IFNULL(REFRESH_TYPE, '') <> 'ROLLUP' ORDER BY TABLE_NAME
: string
| omni_dump_gap_rollup
< result
< column "TABLE_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_gap_rollup
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| r
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `omni_dump_gap_rollup`.`r`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| r
: bytes
| CREATE TABLE `r` (
|   `k` int(11) NOT NULL COMMENT "",
|   `v` int(11) NULL COMMENT "",
|   `d` datetime NULL COMMENT ""
| ) ENGINE=OLAP 
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1 
| PROPERTIES (
| "compression" = "LZ4",
| "fast_schema_evolution" = "true",
| "replicated_storage" = "true",
| "replication_num" = "1"
| );
> exec
| CREATE DATABASE `omni_dump_gap_rollup_apply`
< ok 0 0
> exec
| USE `omni_dump_gap_rollup_apply`
< ok 0 0
> exec
| CREATE TABLE `omni_dump_gap_rollup_apply`.`r` (k INT NOT NULL, v INT NULL, d DATETIME NULL) ENGINE=OLAP DUPLICATE KEY(k) DISTRIBUTED BY HASH(k) BUCKETS 1 PROPERTIES("replication_num" = "1")
< ok 0 0
> exec
| create materialized view r1 as select v,k from r group by v,k
< ok 0 0
> exec
| USE `omni_dump_gap_rollup`
< ok 0 0
> query
| DESC `omni_dump_gap_rollup_apply`.`r` ALL
< result
< column "IndexName" "VARCHAR"
< column "IndexKeysType" "VARCHAR"
< column "Field" "VARCHAR"
< column "Type" "VARCHAR"
< column "Null" "VARCHAR"
< column "Key" "VARCHAR"
< column "Default" "VARCHAR"
< column "Extra" "VARCHAR"
< row
: bytes
| r
: bytes
| DUP_KEYS
: bytes
| k
: bytes
| int
: bytes
| NO
: bytes
| true
: null
: bytes
|
< row
: bytes
|
: bytes
|
: bytes
| v
: bytes
| int
: bytes
| YES
: bytes
| false
: null
: bytes
| NONE
< row
: bytes
|
: bytes
|
: bytes
| d
: bytes
| datetime
: bytes
| YES
: bytes
| false
: null
: bytes
| NONE
< row
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
: bytes
|
< row
: bytes
| r1
: bytes
| AGG_KEYS
: bytes
| v
: bytes
| int
: bytes
| YES
: bytes
| true
: null
: bytes
|
< row
: bytes
|
: bytes
|
: bytes
| k
: bytes
| int
: bytes
| NO
: bytes
| true
: null
: bytes
|
> exec
| DROP DATABASE `omni_dump_gap_rollup_apply`
< ok 0 0

== TestDumpGaps/view-definition-database-qualifier default#1
> exec
| CREATE TABLE t (
| 	k INT NOT NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(k)
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
|
| CREATE MATERIALIZED VIEW mv
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1")
| AS SELECT k FROM t;
< ok 0 0
> exec
|
| CREATE VIEW v (k) AS SELECT k FROM t
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_gap_relocate
< result
< column "SCHEMA_NAME" "VARCHAR"
< row
: bytes
| omni_dump_gap_relocate
> query
| SELECT TABLE_NAME FROM information_schema.materialized_views WHERE TABLE_SCHEMA = ? AND IFNULL(REFRESH_TYPE, '') <> 'ROLLUP' ORDER BY TABLE_NAME
: string
| omni_dump_gap_relocate
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| mv
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_gap_relocate
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| mv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| v
: bytes
| VIEW
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_gap_relocate`.`mv`
< result
< column "Materialized View" "VARCHAR"
< column "Create Materialized View" "VARCHAR"
< row
: bytes
| mv
: bytes
| CREATE MATERIALIZED VIEW `mv` (`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1 
| REFRESH MANUAL
| PROPERTIES (
| "replicated_storage" = "true",
| "replication_num" = "1",
| "storage_medium" = "HDD"
| )
| AS SELECT `t`.`k`
| FROM `omni_dump_gap_relocate`.`t`;
> query
| SHOW CREATE TABLE `omni_dump_gap_relocate`.`t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `k` int(11) NOT NULL COMMENT ""
| ) ENGINE=OLAP 
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1 
| PROPERTIES (
| "compression" = "LZ4",
| "fast_schema_evolution" = "true",
| "replicated_storage" = "true",
| "replication_num" = "1"
| );
> query
| SHOW CREATE VIEW `omni_dump_gap_relocate`.`v`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v` (`k`) SECURITY NONE AS SELECT `omni_dump_gap_relocate`.`t`.`k`
| FROM `omni_dump_gap_relocate`.`t`;
: bytes
| utf8
: bytes
| utf8_general_ci
