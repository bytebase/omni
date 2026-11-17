-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestDumpGaps/function-definition default#1
> exec
| CREATE ALIAS FUNCTION dump_gap_abs(INT) WITH PARAMETER(a) AS abs(a)
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_gap_function
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_gap_function
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_gap_function") ORDER BY Name
< result
< column "Name" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_gap_function
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_gap_function
< result
< column "SPECIFIC_NAME" "CHAR"
> query
| SHOW CREATE FUNCTION dump_gap_abs(INT)
< result
< column "Function Signature" "CHAR"
< column "Create Function" "CHAR"
< row
: bytes
| dump_gap_abs(int)
: bytes
| CREATE ALIAS FUNCTION dump_gap_abs(int) WITH PARAMETER(a) AS abs(`a`);
> query
| SHOW FUNCTIONS LIKE 'dump_gap_abs'
< result
< column "Function Name" "CHAR"
< row
: bytes
| dump_gap_abs
> query
| SELECT COUNT(*) FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND SPECIFIC_NAME = ?
: string
| omni_dump_gap_function
: string
| dump_gap_abs
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0

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
| SHOW CREATE TABLE `omni_dump_gap_rollup`.r
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| r
: bytes
| CREATE TABLE `r` (
|   `k` int NOT NULL,
|   `v` int NULL,
|   `d` datetime NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| );
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_gap_rollup
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_gap_rollup
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_gap_rollup") ORDER BY Name
< result
< column "Name" "TEXT"
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_gap_rollup
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< row
: bytes
| r
: bytes
| BASE TABLE
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_gap_rollup
< result
< column "SPECIFIC_NAME" "CHAR"
> query
| SHOW CREATE TABLE `omni_dump_gap_rollup`.`r`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| r
: bytes
| CREATE TABLE `r` (
|   `k` int NOT NULL,
|   `v` int NULL,
|   `d` datetime NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| );

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
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
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
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_gap_relocate
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_gap_relocate") ORDER BY Name
< result
< column "Name" "TEXT"
< row
: bytes
| mv
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_gap_relocate
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< row
: bytes
| mv
: bytes
| BASE TABLE
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
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_gap_relocate
< result
< column "SPECIFIC_NAME" "CHAR"
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_gap_relocate`.`mv`
< result
< column "Materialized View" "CHAR"
< column "Create Materialized View" "CHAR"
< row
: bytes
| mv
: bytes
| CREATE MATERIALIZED VIEW mv
| (k)
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| )
| AS SELECT `internal`.`omni_dump_gap_relocate`.`t`.`k` FROM `internal`.`omni_dump_gap_relocate`.`t`
> query
| SHOW CREATE TABLE `omni_dump_gap_relocate`.`t`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `k` int NOT NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728"
| );
> query
| SHOW CREATE VIEW `omni_dump_gap_relocate`.`v`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| v
: bytes
| CREATE VIEW `v`
| (k)
|  AS SELECT `internal`.`omni_dump_gap_relocate`.`t`.`k` AS `k` FROM `internal`.`omni_dump_gap_relocate`.`t`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
