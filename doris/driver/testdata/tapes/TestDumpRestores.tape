-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestDumpRestores default#1
> exec
| -- The shapes a Doris schema dump must carry. Each table states a key model,
| -- a distribution, and properties that information_schema does not hold; the
| -- views and materialized views state definitions whose order matters; and the
| -- procedure's body holds semicolons.
|
| CREATE TABLE duplicate_keys (
| 	k INT NOT NULL COMMENT 'key column',
| 	v VARCHAR(20) NULL COMMENT 'a value',
| 	d DATETIME NULL,
| 	w INT GENERATED ALWAYS AS (k + 1)
| ) ENGINE=OLAP
| DUPLICATE KEY(k)
| COMMENT 'duplicate key table'
| DISTRIBUTED BY HASH(k) BUCKETS 3
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
| -- A table with an inline index, added after creation so the dump reads it from
| -- SHOW CREATE TABLE rather than from the statement that made the table.
| CREATE INDEX duplicate_keys_v ON duplicate_keys (v) USING INVERTED;
< ok 0 0
> exec
| CREATE TABLE unique_keys (
| 	k INT NOT NULL,
| 	v INT NULL,
| 	g BIGINT NOT NULL AUTO_INCREMENT
| ) ENGINE=OLAP
| UNIQUE KEY(k)
| DISTRIBUTED BY HASH(k) BUCKETS 2
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
| CREATE TABLE aggregate_keys (
| 	k INT NOT NULL,
| 	s INT SUM NULL,
| 	m INT MAX NULL,
| 	r VARCHAR(10) REPLACE NULL
| ) ENGINE=OLAP
| AGGREGATE KEY(k)
| DISTRIBUTED BY HASH(k) BUCKETS 2
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
| -- Dynamic partitioning: the server materializes a window of partitions around
| -- the current date and states them with the properties that maintain them.
| CREATE TABLE dynamic_partitions (
| 	d DATE NOT NULL,
| 	v INT NULL
| ) ENGINE=OLAP
| UNIQUE KEY(d)
| PARTITION BY RANGE(d) ()
| DISTRIBUTED BY HASH(d) BUCKETS AUTO
| PROPERTIES("replication_num" = "1", "dynamic_partition.enable" = "true", "dynamic_partition.time_unit" = "DAY", "dynamic_partition.start" = "-2", "dynamic_partition.end" = "2", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "2");
< ok 0 0
> exec
| CREATE TABLE list_partitions (
| 	k INT NOT NULL,
| 	v INT NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(k)
| PARTITION BY LIST(k) (PARTITION p1 VALUES IN ('1', '2'), PARTITION p2 VALUES IN ('3'))
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
| -- A partition expression the server keeps as written.
| CREATE TABLE auto_partitions (
| 	d DATE NOT NULL,
| 	v INT NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(d)
| AUTO PARTITION BY RANGE (date_trunc(d, 'day')) ()
| DISTRIBUTED BY HASH(d) BUCKETS AUTO
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
| CREATE TABLE materialized_view_base (
| 	k INT NOT NULL,
| 	v INT NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(k)
| DISTRIBUTED BY RANDOM BUCKETS 2
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
| -- The names invert the dependency order: a_child_mv reads z_parent_mv, and
| -- a_child_view reads z_parent_view, so a dump that wrote them in name order
| -- would not restore.
| CREATE MATERIALIZED VIEW z_parent_mv
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1")
| AS SELECT k, sum(v) AS s FROM materialized_view_base GROUP BY k;
< ok 0 0
> exec
| CREATE MATERIALIZED VIEW a_child_mv
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1")
| AS SELECT k, sum(s) AS s2 FROM z_parent_mv GROUP BY k;
< ok 0 0
> exec
| CREATE VIEW z_parent_view (k, v) AS SELECT k, v FROM materialized_view_base;
< ok 0 0
> exec
| CREATE VIEW a_child_view (k) AS SELECT k FROM z_parent_view;
< ok 0 0
> exec
| -- A view reads a materialized view, which Doris allows in this direction only.
| CREATE VIEW b_mv_view (k, s) AS SELECT k, s FROM z_parent_mv;
< ok 0 0
> exec
| -- A view body is a query, not a bare SELECT: a set operation, a WITH clause and
| -- a scalar subquery each read z_parent_view, and each is alphabetically before
| -- it, so the dump must see the reference inside the form to order them.
| CREATE VIEW aa_union_view (k) AS SELECT k FROM z_parent_view UNION ALL SELECT k FROM z_parent_view;
< ok 0 0
> exec
| CREATE VIEW ab_cte_view (k) AS WITH c AS (SELECT k FROM z_parent_view) SELECT k FROM c;
< ok 0 0
> exec
| CREATE VIEW ac_subquery_view (k) AS SELECT (SELECT max(k) FROM z_parent_view) AS k;
< ok 0 0
> exec
| -- A CTE that shares a name with a dumped view is not a dependency on it:
| -- cyc_parent's WITH clause defines cyc_child, and the view cyc_child reads
| -- cyc_parent. Taking the CTE for an object makes the pair a reference cycle,
| -- which falls back to name order, cyc_child first, and cannot restore.
| CREATE VIEW cyc_parent (k) AS WITH cyc_child AS (SELECT 1 AS k) SELECT k FROM cyc_child;
< ok 0 0
> exec
| CREATE VIEW cyc_child (k) AS SELECT k FROM cyc_parent;
< ok 0 0
> exec
| CREATE PROCEDURE proc_with_body(IN a INT) COMMENT 'doris procedure'
| BEGIN
| 	SET @first = a + 1;
| 	SET @second = @first * 2;
| END
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_roundtrip
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_roundtrip
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_roundtrip") ORDER BY Name
< result
< column "Name" "TEXT"
< row
: bytes
| a_child_mv
< row
: bytes
| z_parent_mv
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< row
: bytes
| a_child_mv
: bytes
| BASE TABLE
< row
: bytes
| a_child_view
: bytes
| VIEW
< row
: bytes
| aa_union_view
: bytes
| VIEW
< row
: bytes
| ab_cte_view
: bytes
| VIEW
< row
: bytes
| ac_subquery_view
: bytes
| VIEW
< row
: bytes
| aggregate_keys
: bytes
| BASE TABLE
< row
: bytes
| auto_partitions
: bytes
| BASE TABLE
< row
: bytes
| b_mv_view
: bytes
| VIEW
< row
: bytes
| cyc_child
: bytes
| VIEW
< row
: bytes
| cyc_parent
: bytes
| VIEW
< row
: bytes
| duplicate_keys
: bytes
| BASE TABLE
< row
: bytes
| dynamic_partitions
: bytes
| BASE TABLE
< row
: bytes
| list_partitions
: bytes
| BASE TABLE
< row
: bytes
| materialized_view_base
: bytes
| BASE TABLE
< row
: bytes
| unique_keys
: bytes
| BASE TABLE
< row
: bytes
| z_parent_mv
: bytes
| BASE TABLE
< row
: bytes
| z_parent_view
: bytes
| VIEW
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_roundtrip
< result
< column "SPECIFIC_NAME" "CHAR"
< row
: bytes
| PROC_WITH_BODY
> query
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`PROC_WITH_BODY`
< result
< column "Procedure" "TEXT"
< column "Create Procedure" "TEXT"
< row
: bytes
| PROC_WITH_BODY
: bytes
| CREATE PROCEDURE proc_with_body(IN a INT) COMMENT 'doris procedure'
| BEGIN
| 	SET @first = a + 1;
| 	SET @second = @first * 2;
| END
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_roundtrip`.`a_child_mv`
< result
< column "Materialized View" "CHAR"
< column "Create Materialized View" "CHAR"
< row
: bytes
| a_child_mv
: bytes
| CREATE MATERIALIZED VIEW a_child_mv
| (k,s2)
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DUPLICATE KEY(`k`, `s2`)
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
| AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k`, sum(`internal`.`omni_dump_roundtrip`.`z_parent_mv`.`s`) AS `s2` FROM `internal`.`omni_dump_roundtrip`.`z_parent_mv` GROUP BY `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`a_child_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| a_child_view
: bytes
| CREATE VIEW `a_child_view`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`aa_union_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| aa_union_view
: bytes
| CREATE VIEW `aa_union_view`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view` UNION ALL SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`ab_cte_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| ab_cte_view
: bytes
| CREATE VIEW `ab_cte_view`
| (k)
|  AS WITH `c` AS (SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`) SELECT `c`.`k` AS `k` FROM c;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`ac_subquery_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| ac_subquery_view
: bytes
| CREATE VIEW `ac_subquery_view`
| (k)
|  AS SELECT (SELECT max(k) FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`) AS `k`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`aggregate_keys`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| aggregate_keys
: bytes
| CREATE TABLE `aggregate_keys` (
|   `k` int NOT NULL,
|   `s` int SUM NULL,
|   `m` int MAX NULL,
|   `r` varchar(10) REPLACE NULL
| ) ENGINE=OLAP
| AGGREGATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 2
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`auto_partitions`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| auto_partitions
: bytes
| CREATE TABLE `auto_partitions` (
|   `d` date NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`d`)
| AUTO PARTITION BY RANGE (date_trunc(`d`, 'day'))
| ()
| DISTRIBUTED BY HASH(`d`) BUCKETS AUTO
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
| SHOW CREATE VIEW `omni_dump_roundtrip`.`b_mv_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| b_mv_view
: bytes
| CREATE VIEW `b_mv_view`
| (k,s)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k` AS `k`, `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`s` AS `s` FROM `internal`.`omni_dump_roundtrip`.`z_parent_mv`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`cyc_child`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| cyc_child
: bytes
| CREATE VIEW `cyc_child`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`cyc_parent`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`cyc_parent`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`cyc_parent`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| cyc_parent
: bytes
| CREATE VIEW `cyc_parent`
| (k)
|  AS WITH `cyc_child` AS (SELECT 1 AS `k`) SELECT `cyc_child`.`k` AS `k` FROM cyc_child;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`duplicate_keys`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| duplicate_keys
: bytes
| CREATE TABLE `duplicate_keys` (
|   `k` int NOT NULL COMMENT "key column",
|   `v` varchar(20) NULL COMMENT "a value",
|   `d` datetime NULL,
|   `w` int AS ((`k` + 1)) NULL,
|   INDEX duplicate_keys_v (`v`) USING INVERTED
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| COMMENT 'duplicate key table'
| DISTRIBUTED BY HASH(`k`) BUCKETS 3
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`dynamic_partitions`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| dynamic_partitions
: bytes
| CREATE TABLE `dynamic_partitions` (
|   `d` date NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| UNIQUE KEY(`d`)
| PARTITION BY RANGE(`d`)
| (PARTITION p20261117 VALUES [('2026-11-17'), ('2026-11-18')),
| PARTITION p20261118 VALUES [('2026-11-18'), ('2026-11-19')),
| PARTITION p20261119 VALUES [('2026-11-19'), ('2026-11-20')))
| DISTRIBUTED BY HASH(`d`) BUCKETS AUTO
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "dynamic_partition.enable" = "true",
| "dynamic_partition.time_unit" = "DAY",
| "dynamic_partition.time_zone" = "Etc/UTC",
| "dynamic_partition.start" = "-2",
| "dynamic_partition.end" = "2",
| "dynamic_partition.prefix" = "p",
| "dynamic_partition.replication_allocation" = "tag.location.default: 1",
| "dynamic_partition.buckets" = "2",
| "dynamic_partition.create_history_partition" = "false",
| "dynamic_partition.history_partition_num" = "-1",
| "dynamic_partition.hot_partition_num" = "0",
| "dynamic_partition.reserved_history_periods" = "NULL",
| "dynamic_partition.storage_policy" = "",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "enable_unique_key_merge_on_write" = "true",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728",
| "enable_mow_light_delete" = "false"
| );
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`list_partitions`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| list_partitions
: bytes
| CREATE TABLE `list_partitions` (
|   `k` int NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| PARTITION BY LIST (`k`)
| (PARTITION p1 VALUES IN ("1","2"),
| PARTITION p2 VALUES IN ("3"))
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`materialized_view_base`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| materialized_view_base
: bytes
| CREATE TABLE `materialized_view_base` (
|   `k` int NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY RANDOM BUCKETS 2
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`unique_keys`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| unique_keys
: bytes
| CREATE TABLE `unique_keys` (
|   `k` int NOT NULL,
|   `v` int NULL,
|   `g` bigint NOT NULL AUTO_INCREMENT(1)
| ) ENGINE=OLAP
| UNIQUE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 2
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "enable_unique_key_merge_on_write" = "true",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728",
| "enable_mow_light_delete" = "false"
| );
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_roundtrip`.`z_parent_mv`
< result
< column "Materialized View" "CHAR"
< column "Create Materialized View" "CHAR"
< row
: bytes
| z_parent_mv
: bytes
| CREATE MATERIALIZED VIEW z_parent_mv
| (k,s)
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DUPLICATE KEY(`k`, `s`)
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
| AS SELECT `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k`, sum(`internal`.`omni_dump_roundtrip`.`materialized_view_base`.`v`) AS `s` FROM `internal`.`omni_dump_roundtrip`.`materialized_view_base` GROUP BY `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`z_parent_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| z_parent_view
: bytes
| CREATE VIEW `z_parent_view`
| (k,v)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k` AS `k`, `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`v` AS `v` FROM `internal`.`omni_dump_roundtrip`.`materialized_view_base`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> exec
| DROP DATABASE `omni_dump_roundtrip`
< ok 0 0
> exec
| CREATE DATABASE `omni_dump_roundtrip`
< ok 0 0
> exec
| USE `omni_dump_roundtrip`
< ok 0 0
> exec
| --
| -- Table structure for `aggregate_keys`
| --
| CREATE TABLE `aggregate_keys` (
|   `k` int NOT NULL,
|   `s` int SUM NULL,
|   `m` int MAX NULL,
|   `r` varchar(10) REPLACE NULL
| ) ENGINE=OLAP
| AGGREGATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 2
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
< ok 0 0
> exec
| --
| -- Table structure for `auto_partitions`
| --
| CREATE TABLE `auto_partitions` (
|   `d` date NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`d`)
| AUTO PARTITION BY RANGE (date_trunc(`d`, 'day'))
| ()
| DISTRIBUTED BY HASH(`d`) BUCKETS AUTO
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
< ok 0 0
> exec
| --
| -- Table structure for `duplicate_keys`
| --
| CREATE TABLE `duplicate_keys` (
|   `k` int NOT NULL COMMENT "key column",
|   `v` varchar(20) NULL COMMENT "a value",
|   `d` datetime NULL,
|   `w` int AS ((`k` + 1)) NULL,
|   INDEX duplicate_keys_v (`v`) USING INVERTED
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| COMMENT 'duplicate key table'
| DISTRIBUTED BY HASH(`k`) BUCKETS 3
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
< ok 0 0
> exec
| --
| -- Table structure for `dynamic_partitions`
| --
| CREATE TABLE `dynamic_partitions` (
|   `d` date NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| UNIQUE KEY(`d`)
| PARTITION BY RANGE(`d`)
| (PARTITION p20261117 VALUES [('2026-11-17'), ('2026-11-18')),
| PARTITION p20261118 VALUES [('2026-11-18'), ('2026-11-19')),
| PARTITION p20261119 VALUES [('2026-11-19'), ('2026-11-20')))
| DISTRIBUTED BY HASH(`d`) BUCKETS AUTO
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "dynamic_partition.enable" = "true",
| "dynamic_partition.time_unit" = "DAY",
| "dynamic_partition.time_zone" = "Etc/UTC",
| "dynamic_partition.start" = "-2",
| "dynamic_partition.end" = "2",
| "dynamic_partition.prefix" = "p",
| "dynamic_partition.replication_allocation" = "tag.location.default: 1",
| "dynamic_partition.buckets" = "2",
| "dynamic_partition.create_history_partition" = "false",
| "dynamic_partition.history_partition_num" = "-1",
| "dynamic_partition.hot_partition_num" = "0",
| "dynamic_partition.reserved_history_periods" = "NULL",
| "dynamic_partition.storage_policy" = "",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "enable_unique_key_merge_on_write" = "true",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728",
| "enable_mow_light_delete" = "false"
| );
< ok 0 0
> exec
| --
| -- Table structure for `list_partitions`
| --
| CREATE TABLE `list_partitions` (
|   `k` int NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| PARTITION BY LIST (`k`)
| (PARTITION p1 VALUES IN ("1","2"),
| PARTITION p2 VALUES IN ("3"))
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
< ok 0 0
> exec
| --
| -- Table structure for `materialized_view_base`
| --
| CREATE TABLE `materialized_view_base` (
|   `k` int NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY RANDOM BUCKETS 2
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
< ok 0 0
> exec
| --
| -- Table structure for `unique_keys`
| --
| CREATE TABLE `unique_keys` (
|   `k` int NOT NULL,
|   `v` int NULL,
|   `g` bigint NOT NULL AUTO_INCREMENT(1)
| ) ENGINE=OLAP
| UNIQUE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 2
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "enable_unique_key_merge_on_write" = "true",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728",
| "enable_mow_light_delete" = "false"
| );
< ok 0 0
> exec
| --
| -- View structure for `cyc_parent`
| --
| CREATE VIEW `cyc_parent`
| (k)
|  AS WITH `cyc_child` AS (SELECT 1 AS `k`) SELECT `cyc_child`.`k` AS `k` FROM cyc_child;
< ok 0 0
> exec
| --
| -- Materialized view structure for `z_parent_mv`
| --
| CREATE MATERIALIZED VIEW z_parent_mv
| (k,s)
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DUPLICATE KEY(`k`, `s`)
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
| AS SELECT `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k`, sum(`internal`.`omni_dump_roundtrip`.`materialized_view_base`.`v`) AS `s` FROM `internal`.`omni_dump_roundtrip`.`materialized_view_base` GROUP BY `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k`;
< ok 0 0
> exec
| --
| -- View structure for `z_parent_view`
| --
| CREATE VIEW `z_parent_view`
| (k,v)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k` AS `k`, `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`v` AS `v` FROM `internal`.`omni_dump_roundtrip`.`materialized_view_base`;
< ok 0 0
> exec
| --
| -- Materialized view structure for `a_child_mv`
| --
| CREATE MATERIALIZED VIEW a_child_mv
| (k,s2)
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DUPLICATE KEY(`k`, `s2`)
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
| AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k`, sum(`internal`.`omni_dump_roundtrip`.`z_parent_mv`.`s`) AS `s2` FROM `internal`.`omni_dump_roundtrip`.`z_parent_mv` GROUP BY `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k`;
< ok 0 0
> exec
| --
| -- View structure for `a_child_view`
| --
| CREATE VIEW `a_child_view`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`;
< ok 0 0
> exec
| --
| -- View structure for `aa_union_view`
| --
| CREATE VIEW `aa_union_view`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view` UNION ALL SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`;
< ok 0 0
> exec
| --
| -- View structure for `ab_cte_view`
| --
| CREATE VIEW `ab_cte_view`
| (k)
|  AS WITH `c` AS (SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`) SELECT `c`.`k` AS `k` FROM c;
< ok 0 0
> exec
| --
| -- View structure for `ac_subquery_view`
| --
| CREATE VIEW `ac_subquery_view`
| (k)
|  AS SELECT (SELECT max(k) FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`) AS `k`;
< ok 0 0
> exec
| --
| -- View structure for `b_mv_view`
| --
| CREATE VIEW `b_mv_view`
| (k,s)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k` AS `k`, `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`s` AS `s` FROM `internal`.`omni_dump_roundtrip`.`z_parent_mv`;
< ok 0 0
> exec
| --
| -- View structure for `cyc_child`
| --
| CREATE VIEW `cyc_child`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`cyc_parent`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`cyc_parent`;
< ok 0 0
> exec
| CREATE PROCEDURE proc_with_body(IN a INT) COMMENT 'doris procedure'
| BEGIN
| 	SET @first = a + 1;
| 	SET @second = @first * 2;
| END
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_roundtrip
< result
< column "SCHEMA_NAME" "CHAR"
< row
: bytes
| omni_dump_roundtrip
> query
| SELECT current_catalog()
< result
< column "current_catalog()" "CHAR"
< row
: bytes
| internal
> query
| SELECT Name FROM mv_infos("database"="omni_dump_roundtrip") ORDER BY Name
< result
< column "Name" "TEXT"
< row
: bytes
| a_child_mv
< row
: bytes
| z_parent_mv
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_roundtrip
< result
< column "TABLE_NAME" "CHAR"
< column "TABLE_TYPE" "CHAR"
< row
: bytes
| a_child_mv
: bytes
| BASE TABLE
< row
: bytes
| a_child_view
: bytes
| VIEW
< row
: bytes
| aa_union_view
: bytes
| VIEW
< row
: bytes
| ab_cte_view
: bytes
| VIEW
< row
: bytes
| ac_subquery_view
: bytes
| VIEW
< row
: bytes
| aggregate_keys
: bytes
| BASE TABLE
< row
: bytes
| auto_partitions
: bytes
| BASE TABLE
< row
: bytes
| b_mv_view
: bytes
| VIEW
< row
: bytes
| cyc_child
: bytes
| VIEW
< row
: bytes
| cyc_parent
: bytes
| VIEW
< row
: bytes
| duplicate_keys
: bytes
| BASE TABLE
< row
: bytes
| dynamic_partitions
: bytes
| BASE TABLE
< row
: bytes
| list_partitions
: bytes
| BASE TABLE
< row
: bytes
| materialized_view_base
: bytes
| BASE TABLE
< row
: bytes
| unique_keys
: bytes
| BASE TABLE
< row
: bytes
| z_parent_mv
: bytes
| BASE TABLE
< row
: bytes
| z_parent_view
: bytes
| VIEW
> query
| SELECT SPECIFIC_NAME FROM information_schema.routines WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE' ORDER BY SPECIFIC_NAME
: string
| omni_dump_roundtrip
< result
< column "SPECIFIC_NAME" "CHAR"
< row
: bytes
| PROC_WITH_BODY
> query
| SHOW CREATE PROCEDURE `omni_dump_roundtrip`.`PROC_WITH_BODY`
< result
< column "Procedure" "TEXT"
< column "Create Procedure" "TEXT"
< row
: bytes
| PROC_WITH_BODY
: bytes
| CREATE PROCEDURE proc_with_body(IN a INT) COMMENT 'doris procedure'
| BEGIN
| 	SET @first = a + 1;
| 	SET @second = @first * 2;
| END
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_roundtrip`.`a_child_mv`
< result
< column "Materialized View" "CHAR"
< column "Create Materialized View" "CHAR"
< row
: bytes
| a_child_mv
: bytes
| CREATE MATERIALIZED VIEW a_child_mv
| (k,s2)
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DUPLICATE KEY(`k`, `s2`)
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
| AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k`, sum(`internal`.`omni_dump_roundtrip`.`z_parent_mv`.`s`) AS `s2` FROM `internal`.`omni_dump_roundtrip`.`z_parent_mv` GROUP BY `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`a_child_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| a_child_view
: bytes
| CREATE VIEW `a_child_view`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`aa_union_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| aa_union_view
: bytes
| CREATE VIEW `aa_union_view`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view` UNION ALL SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`ab_cte_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| ab_cte_view
: bytes
| CREATE VIEW `ab_cte_view`
| (k)
|  AS WITH `c` AS (SELECT `internal`.`omni_dump_roundtrip`.`z_parent_view`.`k` FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`) SELECT `c`.`k` AS `k` FROM c;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`ac_subquery_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| ac_subquery_view
: bytes
| CREATE VIEW `ac_subquery_view`
| (k)
|  AS SELECT (SELECT max(k) FROM `internal`.`omni_dump_roundtrip`.`z_parent_view`) AS `k`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`aggregate_keys`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| aggregate_keys
: bytes
| CREATE TABLE `aggregate_keys` (
|   `k` int NOT NULL,
|   `s` int SUM NULL,
|   `m` int MAX NULL,
|   `r` varchar(10) REPLACE NULL
| ) ENGINE=OLAP
| AGGREGATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 2
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`auto_partitions`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| auto_partitions
: bytes
| CREATE TABLE `auto_partitions` (
|   `d` date NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`d`)
| AUTO PARTITION BY RANGE (date_trunc(`d`, 'day'))
| ()
| DISTRIBUTED BY HASH(`d`) BUCKETS AUTO
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
| SHOW CREATE VIEW `omni_dump_roundtrip`.`b_mv_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| b_mv_view
: bytes
| CREATE VIEW `b_mv_view`
| (k,s)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`k` AS `k`, `internal`.`omni_dump_roundtrip`.`z_parent_mv`.`s` AS `s` FROM `internal`.`omni_dump_roundtrip`.`z_parent_mv`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`cyc_child`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| cyc_child
: bytes
| CREATE VIEW `cyc_child`
| (k)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`cyc_parent`.`k` AS `k` FROM `internal`.`omni_dump_roundtrip`.`cyc_parent`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`cyc_parent`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| cyc_parent
: bytes
| CREATE VIEW `cyc_parent`
| (k)
|  AS WITH `cyc_child` AS (SELECT 1 AS `k`) SELECT `cyc_child`.`k` AS `k` FROM cyc_child;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`duplicate_keys`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| duplicate_keys
: bytes
| CREATE TABLE `duplicate_keys` (
|   `k` int NOT NULL COMMENT "key column",
|   `v` varchar(20) NULL COMMENT "a value",
|   `d` datetime NULL,
|   `w` int AS ((`k` + 1)) NULL,
|   INDEX duplicate_keys_v (`v`) USING INVERTED
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| COMMENT 'duplicate key table'
| DISTRIBUTED BY HASH(`k`) BUCKETS 3
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`dynamic_partitions`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| dynamic_partitions
: bytes
| CREATE TABLE `dynamic_partitions` (
|   `d` date NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| UNIQUE KEY(`d`)
| PARTITION BY RANGE(`d`)
| (PARTITION p20261117 VALUES [('2026-11-17'), ('2026-11-18')),
| PARTITION p20261118 VALUES [('2026-11-18'), ('2026-11-19')),
| PARTITION p20261119 VALUES [('2026-11-19'), ('2026-11-20')))
| DISTRIBUTED BY HASH(`d`) BUCKETS AUTO
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "dynamic_partition.enable" = "true",
| "dynamic_partition.time_unit" = "DAY",
| "dynamic_partition.time_zone" = "Etc/UTC",
| "dynamic_partition.start" = "-2",
| "dynamic_partition.end" = "2",
| "dynamic_partition.prefix" = "p",
| "dynamic_partition.replication_allocation" = "tag.location.default: 1",
| "dynamic_partition.buckets" = "2",
| "dynamic_partition.create_history_partition" = "false",
| "dynamic_partition.history_partition_num" = "-1",
| "dynamic_partition.hot_partition_num" = "0",
| "dynamic_partition.reserved_history_periods" = "NULL",
| "dynamic_partition.storage_policy" = "",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "enable_unique_key_merge_on_write" = "true",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728",
| "enable_mow_light_delete" = "false"
| );
> query
| SHOW CREATE TABLE `omni_dump_roundtrip`.`list_partitions`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| list_partitions
: bytes
| CREATE TABLE `list_partitions` (
|   `k` int NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| PARTITION BY LIST (`k`)
| (PARTITION p1 VALUES IN ("1","2"),
| PARTITION p2 VALUES IN ("3"))
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`materialized_view_base`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| materialized_view_base
: bytes
| CREATE TABLE `materialized_view_base` (
|   `k` int NOT NULL,
|   `v` int NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(`k`)
| DISTRIBUTED BY RANDOM BUCKETS 2
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
| SHOW CREATE TABLE `omni_dump_roundtrip`.`unique_keys`
< result
< column "Table" "CHAR"
< column "Create Table" "CHAR"
< row
: bytes
| unique_keys
: bytes
| CREATE TABLE `unique_keys` (
|   `k` int NOT NULL,
|   `v` int NULL,
|   `g` bigint NOT NULL AUTO_INCREMENT(1)
| ) ENGINE=OLAP
| UNIQUE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 2
| PROPERTIES (
| "replication_allocation" = "tag.location.default: 1",
| "min_load_replica_num" = "-1",
| "is_being_synced" = "false",
| "storage_medium" = "hdd",
| "storage_format" = "V2",
| "inverted_index_storage_format" = "V2",
| "enable_unique_key_merge_on_write" = "true",
| "light_schema_change" = "true",
| "disable_auto_compaction" = "false",
| "enable_single_replica_compaction" = "false",
| "group_commit_interval_ms" = "10000",
| "group_commit_data_bytes" = "134217728",
| "enable_mow_light_delete" = "false"
| );
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_roundtrip`.`z_parent_mv`
< result
< column "Materialized View" "CHAR"
< column "Create Materialized View" "CHAR"
< row
: bytes
| z_parent_mv
: bytes
| CREATE MATERIALIZED VIEW z_parent_mv
| (k,s)
| BUILD DEFERRED REFRESH COMPLETE ON MANUAL
| DUPLICATE KEY(`k`, `s`)
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
| AS SELECT `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k`, sum(`internal`.`omni_dump_roundtrip`.`materialized_view_base`.`v`) AS `s` FROM `internal`.`omni_dump_roundtrip`.`materialized_view_base` GROUP BY `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k`
> query
| SHOW CREATE VIEW `omni_dump_roundtrip`.`z_parent_view`
< result
< column "View" "CHAR"
< column "Create View" "CHAR"
< column "character_set_client" "CHAR"
< column "collation_connection" "CHAR"
< row
: bytes
| z_parent_view
: bytes
| CREATE VIEW `z_parent_view`
| (k,v)
|  AS SELECT `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`k` AS `k`, `internal`.`omni_dump_roundtrip`.`materialized_view_base`.`v` AS `v` FROM `internal`.`omni_dump_roundtrip`.`materialized_view_base`;
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
