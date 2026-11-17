--
-- Table structure for `aggregate_keys`
--
CREATE TABLE `aggregate_keys` (
  `k` int NOT NULL,
  `s` int SUM NULL,
  `m` int MAX NULL,
  `r` varchar(10) REPLACE NULL
) ENGINE=OLAP
AGGREGATE KEY(`k`)
DISTRIBUTED BY HASH(`k`) BUCKETS 2
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
);

--
-- Table structure for `auto_partitions`
--
CREATE TABLE `auto_partitions` (
  `d` date NOT NULL,
  `v` int NULL
) ENGINE=OLAP
DUPLICATE KEY(`d`)
AUTO PARTITION BY RANGE (date_trunc(`d`, 'day'))
()
DISTRIBUTED BY HASH(`d`) BUCKETS AUTO
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
);

--
-- Table structure for `duplicate_keys`
--
CREATE TABLE `duplicate_keys` (
  `k` int NOT NULL COMMENT "key column",
  `v` varchar(20) NULL COMMENT "a value",
  `d` datetime NULL,
  `w` int AS ((`k` + 1)) NULL,
  INDEX duplicate_keys_v (`v`) USING INVERTED
) ENGINE=OLAP
DUPLICATE KEY(`k`)
COMMENT 'duplicate key table'
DISTRIBUTED BY HASH(`k`) BUCKETS 3
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
);

--
-- Table structure for `dynamic_partitions`
--
CREATE TABLE `dynamic_partitions` (
  `d` date NOT NULL,
  `v` int NULL
) ENGINE=OLAP
UNIQUE KEY(`d`)
PARTITION BY RANGE(`d`)
(PARTITION p20261117 VALUES [('2026-11-17'), ('2026-11-18')),
PARTITION p20261118 VALUES [('2026-11-18'), ('2026-11-19')),
PARTITION p20261119 VALUES [('2026-11-19'), ('2026-11-20')))
DISTRIBUTED BY HASH(`d`) BUCKETS AUTO
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"dynamic_partition.enable" = "true",
"dynamic_partition.time_unit" = "DAY",
"dynamic_partition.time_zone" = "Etc/UTC",
"dynamic_partition.start" = "-2",
"dynamic_partition.end" = "2",
"dynamic_partition.prefix" = "p",
"dynamic_partition.replication_allocation" = "tag.location.default: 1",
"dynamic_partition.buckets" = "2",
"dynamic_partition.create_history_partition" = "false",
"dynamic_partition.history_partition_num" = "-1",
"dynamic_partition.hot_partition_num" = "0",
"dynamic_partition.reserved_history_periods" = "NULL",
"dynamic_partition.storage_policy" = "",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"enable_unique_key_merge_on_write" = "true",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728",
"enable_mow_light_delete" = "false"
);

--
-- Table structure for `list_partitions`
--
CREATE TABLE `list_partitions` (
  `k` int NOT NULL,
  `v` int NULL
) ENGINE=OLAP
DUPLICATE KEY(`k`)
PARTITION BY LIST (`k`)
(PARTITION p1 VALUES IN ("1","2"),
PARTITION p2 VALUES IN ("3"))
DISTRIBUTED BY HASH(`k`) BUCKETS 1
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
);

--
-- Table structure for `materialized_view_base`
--
CREATE TABLE `materialized_view_base` (
  `k` int NOT NULL,
  `v` int NULL
) ENGINE=OLAP
DUPLICATE KEY(`k`)
DISTRIBUTED BY RANDOM BUCKETS 2
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
);

--
-- Table structure for `select_mv_base`
--
CREATE TABLE `select_mv_base` (
  `k` int NOT NULL,
  `v` int NULL
) ENGINE=OLAP
DUPLICATE KEY(`k`)
DISTRIBUTED BY HASH(`k`) BUCKETS 1
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
);

--
-- Table structure for `sync_mv_base`
--
CREATE TABLE `sync_mv_base` (
  `k` int NOT NULL,
  `v` int NULL
) ENGINE=OLAP
DUPLICATE KEY(`k`)
DISTRIBUTED BY HASH(`k`) BUCKETS 1
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
);

--
-- Table structure for `unique_keys`
--
CREATE TABLE `unique_keys` (
  `k` int NOT NULL,
  `v` int NULL,
  `g` bigint NOT NULL AUTO_INCREMENT(1)
) ENGINE=OLAP
UNIQUE KEY(`k`)
DISTRIBUTED BY HASH(`k`) BUCKETS 2
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"enable_unique_key_merge_on_write" = "true",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728",
"enable_mow_light_delete" = "false"
);

CREATE MATERIALIZED VIEW `select` AS SELECT k, max(v) AS m FROM select_mv_base GROUP BY k;

CREATE MATERIALIZED VIEW smv_total AS SELECT k, sum(v) AS s FROM sync_mv_base GROUP BY k;

--
-- Function structure for alias_abs(int)
--
CREATE ALIAS FUNCTION alias_abs(int) WITH PARAMETER(a) AS abs(`a`);

--
-- View structure for `af_function_view`
--
CREATE VIEW `af_function_view`
(k)
 AS SELECT `omni_dump_corpus`.`alias_abs`(`internal`.`omni_dump_corpus`.`materialized_view_base`.`k`) AS `k` FROM `internal`.`omni_dump_corpus`.`materialized_view_base`;

--
-- View structure for `cyc_parent`
--
CREATE VIEW `cyc_parent`
(k)
 AS WITH `cyc_child` AS (SELECT 1 AS `k`) SELECT `cyc_child`.`k` AS `k` FROM cyc_child;

--
-- Materialized view structure for `z_parent_mv`
--
CREATE MATERIALIZED VIEW z_parent_mv
(k,s)
BUILD DEFERRED REFRESH COMPLETE ON MANUAL
DUPLICATE KEY(`k`, `s`)
DISTRIBUTED BY HASH(`k`) BUCKETS 1
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
)
AS SELECT `internal`.`omni_dump_corpus`.`materialized_view_base`.`k`, sum(`internal`.`omni_dump_corpus`.`materialized_view_base`.`v`) AS `s` FROM `internal`.`omni_dump_corpus`.`materialized_view_base` GROUP BY `internal`.`omni_dump_corpus`.`materialized_view_base`.`k`;

--
-- View structure for `z_parent_view`
--
CREATE VIEW `z_parent_view`
(k,v)
 AS SELECT `internal`.`omni_dump_corpus`.`materialized_view_base`.`k` AS `k`, `internal`.`omni_dump_corpus`.`materialized_view_base`.`v` AS `v` FROM `internal`.`omni_dump_corpus`.`materialized_view_base`;

--
-- Materialized view structure for `a_child_mv`
--
CREATE MATERIALIZED VIEW a_child_mv
(k,s2)
BUILD DEFERRED REFRESH COMPLETE ON MANUAL
DUPLICATE KEY(`k`, `s2`)
DISTRIBUTED BY HASH(`k`) BUCKETS 1
PROPERTIES (
"replication_allocation" = "tag.location.default: 1",
"min_load_replica_num" = "-1",
"is_being_synced" = "false",
"storage_medium" = "hdd",
"storage_format" = "V2",
"inverted_index_storage_format" = "V2",
"light_schema_change" = "true",
"disable_auto_compaction" = "false",
"enable_single_replica_compaction" = "false",
"group_commit_interval_ms" = "10000",
"group_commit_data_bytes" = "134217728"
)
AS SELECT `internal`.`omni_dump_corpus`.`z_parent_mv`.`k`, sum(`internal`.`omni_dump_corpus`.`z_parent_mv`.`s`) AS `s2` FROM `internal`.`omni_dump_corpus`.`z_parent_mv` GROUP BY `internal`.`omni_dump_corpus`.`z_parent_mv`.`k`;

--
-- View structure for `a_child_view`
--
CREATE VIEW `a_child_view`
(k)
 AS SELECT `internal`.`omni_dump_corpus`.`z_parent_view`.`k` AS `k` FROM `internal`.`omni_dump_corpus`.`z_parent_view`;

--
-- View structure for `aa_union_view`
--
CREATE VIEW `aa_union_view`
(k)
 AS SELECT `internal`.`omni_dump_corpus`.`z_parent_view`.`k` AS `k` FROM `internal`.`omni_dump_corpus`.`z_parent_view` UNION ALL SELECT `internal`.`omni_dump_corpus`.`z_parent_view`.`k` FROM `internal`.`omni_dump_corpus`.`z_parent_view`;

--
-- View structure for `ab_cte_view`
--
CREATE VIEW `ab_cte_view`
(k)
 AS WITH `c` AS (SELECT `internal`.`omni_dump_corpus`.`z_parent_view`.`k` FROM `internal`.`omni_dump_corpus`.`z_parent_view`) SELECT `c`.`k` AS `k` FROM c;

--
-- View structure for `ac_subquery_view`
--
CREATE VIEW `ac_subquery_view`
(k)
 AS SELECT (SELECT max(k) FROM `internal`.`omni_dump_corpus`.`z_parent_view`) AS `k`;

--
-- View structure for `b_mv_view`
--
CREATE VIEW `b_mv_view`
(k,s)
 AS SELECT `internal`.`omni_dump_corpus`.`z_parent_mv`.`k` AS `k`, `internal`.`omni_dump_corpus`.`z_parent_mv`.`s` AS `s` FROM `internal`.`omni_dump_corpus`.`z_parent_mv`;

--
-- View structure for `cyc_child`
--
CREATE VIEW `cyc_child`
(k)
 AS SELECT `internal`.`omni_dump_corpus`.`cyc_parent`.`k` AS `k` FROM `internal`.`omni_dump_corpus`.`cyc_parent`;

--
-- Procedure structure for `PROC_WITH_BODY`
--
DELIMITER ;;
CREATE PROCEDURE proc_with_body(IN a INT) COMMENT 'doris procedure'
BEGIN
	SET @first = a + 1;
	SET @second = @first * 2;
END;;
DELIMITER ;

