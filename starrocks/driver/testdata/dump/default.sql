--
-- Table structure for `aggregate_keys`
--
CREATE TABLE `aggregate_keys` (
  `k` int(11) NOT NULL COMMENT "",
  `s` int(11) SUM NULL COMMENT "",
  `m` int(11) MAX NULL COMMENT "",
  `r` varchar(10) REPLACE NULL COMMENT ""
) ENGINE=OLAP 
AGGREGATE KEY(`k`)
DISTRIBUTED BY HASH(`k`) BUCKETS 2 
PROPERTIES (
"compression" = "LZ4",
"fast_schema_evolution" = "true",
"replicated_storage" = "true",
"replication_num" = "1"
);

--
-- Table structure for `duplicate_keys`
--
CREATE TABLE `duplicate_keys` (
  `k` int(11) NOT NULL COMMENT "key column",
  `v` varchar(20) NULL COMMENT "a value",
  `d` datetime NULL COMMENT "",
  `w` bigint(20) NULL AS `k` + 1 COMMENT "",
  INDEX duplicate_keys_v (`v`) USING NGRAMBF("bloom_filter_fpp" = "0.05", "case_sensitive" = "true", "gram_num" = "2") COMMENT ''
) ENGINE=OLAP 
DUPLICATE KEY(`k`)
COMMENT "duplicate key table"
DISTRIBUTED BY HASH(`k`) BUCKETS 3 
PROPERTIES (
"compression" = "LZ4",
"fast_schema_evolution" = "true",
"replicated_storage" = "true",
"replication_num" = "1"
);

--
-- Table structure for `dynamic_partitions`
--
CREATE TABLE `dynamic_partitions` (
  `d` date NOT NULL COMMENT "",
  `v` int(11) NULL COMMENT ""
) ENGINE=OLAP 
UNIQUE KEY(`d`)
PARTITION BY RANGE(`d`)
(PARTITION p20261117 VALUES [("2026-11-17"), ("2026-11-18")),
PARTITION p20261118 VALUES [("2026-11-18"), ("2026-11-19")),
PARTITION p20261119 VALUES [("2026-11-19"), ("2026-11-20")))
DISTRIBUTED BY HASH(`d`)
PROPERTIES (
"compression" = "LZ4",
"dynamic_partition.buckets" = "2",
"dynamic_partition.enable" = "true",
"dynamic_partition.end" = "2",
"dynamic_partition.history_partition_num" = "0",
"dynamic_partition.prefix" = "p",
"dynamic_partition.start" = "-2",
"dynamic_partition.time_unit" = "DAY",
"dynamic_partition.time_zone" = "Etc/UTC",
"fast_schema_evolution" = "true",
"replicated_storage" = "true",
"replication_num" = "1"
);

--
-- Table structure for `list_partitions`
--
CREATE TABLE `list_partitions` (
  `k` int(11) NOT NULL COMMENT "",
  `v` int(11) NULL COMMENT ""
) ENGINE=OLAP 
DUPLICATE KEY(`k`)
PARTITION BY LIST(`k`)(
  PARTITION p1 VALUES IN ('1', '2'),
  PARTITION p2 VALUES IN ('3')
)
DISTRIBUTED BY HASH(`k`) BUCKETS 1 
PROPERTIES (
"compression" = "LZ4",
"fast_schema_evolution" = "true",
"replicated_storage" = "true",
"replication_num" = "1"
);

--
-- Table structure for `materialized_view_base`
--
CREATE TABLE `materialized_view_base` (
  `k` int(11) NOT NULL COMMENT "",
  `v` int(11) NULL COMMENT ""
) ENGINE=OLAP 
DUPLICATE KEY(`k`)
DISTRIBUTED BY HASH(`k`) BUCKETS 1 
PROPERTIES (
"compression" = "LZ4",
"fast_schema_evolution" = "true",
"replicated_storage" = "true",
"replication_num" = "1"
);

--
-- Table structure for `range_partitions`
--
CREATE TABLE `range_partitions` (
  `d` date NOT NULL COMMENT "",
  `v` int(11) NULL COMMENT ""
) ENGINE=OLAP 
DUPLICATE KEY(`d`)
PARTITION BY RANGE(`d`)
(PARTITION p2025 VALUES [("2025-01-01"), ("2026-01-01")),
PARTITION p2026 VALUES [("2026-01-01"), ("2027-01-01")))
DISTRIBUTED BY HASH(`d`) BUCKETS 1 
PROPERTIES (
"compression" = "LZ4",
"fast_schema_evolution" = "true",
"replicated_storage" = "true",
"replication_num" = "1"
);

--
-- Table structure for `unique_keys`
--
CREATE TABLE `unique_keys` (
  `k` int(11) NOT NULL COMMENT "",
  `v` int(11) NULL COMMENT "",
  `g` bigint(20) NOT NULL AUTO_INCREMENT COMMENT ""
) ENGINE=OLAP 
UNIQUE KEY(`k`)
DISTRIBUTED BY HASH(`k`) BUCKETS 2 
PROPERTIES (
"compression" = "LZ4",
"fast_schema_evolution" = "true",
"replicated_storage" = "true",
"replication_num" = "1"
);

--
-- View structure for `cyc_parent`
--
CREATE VIEW `cyc_parent` (`k`) SECURITY NONE AS WITH `cyc_child` (`k`) AS (SELECT 1 AS `k`) SELECT `cyc_child`.`k`
FROM `cyc_child`;

--
-- Materialized view structure for `z_parent_mv`
--
CREATE MATERIALIZED VIEW `z_parent_mv` (`k`, `s`)
DISTRIBUTED BY HASH(`k`) BUCKETS 1 
REFRESH MANUAL
PROPERTIES (
"replicated_storage" = "true",
"replication_num" = "1",
"storage_medium" = "HDD"
)
AS SELECT `materialized_view_base`.`k`, sum(`materialized_view_base`.`v`) AS `s`
FROM `omni_dump_corpus`.`materialized_view_base`
GROUP BY `materialized_view_base`.`k`;

--
-- View structure for `z_parent_view`
--
CREATE VIEW `z_parent_view` (`k`, `v`) SECURITY NONE AS SELECT `omni_dump_corpus`.`materialized_view_base`.`k`, `omni_dump_corpus`.`materialized_view_base`.`v`
FROM `omni_dump_corpus`.`materialized_view_base`;

--
-- Materialized view structure for `a_child_mv`
--
CREATE MATERIALIZED VIEW `a_child_mv` (`k`, `s2`)
DISTRIBUTED BY HASH(`k`) BUCKETS 1 
REFRESH MANUAL
PROPERTIES (
"replicated_storage" = "true",
"replication_num" = "1",
"storage_medium" = "HDD"
)
AS SELECT `z_parent_mv`.`k`, sum(`z_parent_mv`.`s`) AS `s2`
FROM `omni_dump_corpus`.`z_parent_mv`
GROUP BY `z_parent_mv`.`k`;

--
-- View structure for `a_child_view`
--
CREATE VIEW `a_child_view` (`k`) SECURITY NONE AS SELECT `omni_dump_corpus`.`z_parent_view`.`k`
FROM `omni_dump_corpus`.`z_parent_view`;

--
-- View structure for `aa_union_view`
--
CREATE VIEW `aa_union_view` (`k`) SECURITY NONE AS (SELECT `omni_dump_corpus`.`z_parent_view`.`k`
FROM `omni_dump_corpus`.`z_parent_view`) UNION ALL (SELECT `omni_dump_corpus`.`z_parent_view`.`k`
FROM `omni_dump_corpus`.`z_parent_view`);

--
-- View structure for `ab_cte_view`
--
CREATE VIEW `ab_cte_view` (`k`) SECURITY NONE AS WITH `c` (`k`) AS (SELECT `omni_dump_corpus`.`z_parent_view`.`k`
FROM `omni_dump_corpus`.`z_parent_view`) SELECT `c`.`k`
FROM `c`;

--
-- View structure for `ac_subquery_view`
--
CREATE VIEW `ac_subquery_view` (`k`) SECURITY NONE AS SELECT (SELECT max(`omni_dump_corpus`.`z_parent_view`.`k`) AS `max(k)`
FROM `omni_dump_corpus`.`z_parent_view`) AS `k`;

--
-- View structure for `b_mv_view`
--
CREATE VIEW `b_mv_view` (`k`, `s`) SECURITY NONE AS SELECT `omni_dump_corpus`.`z_parent_mv`.`k`, `omni_dump_corpus`.`z_parent_mv`.`s`
FROM `omni_dump_corpus`.`z_parent_mv`;

--
-- View structure for `cyc_child`
--
CREATE VIEW `cyc_child` (`k`) SECURITY NONE AS SELECT `omni_dump_corpus`.`cyc_parent`.`k`
FROM `omni_dump_corpus`.`cyc_parent`;

--
-- Materialized view structure for `mv_on_view`
--
CREATE MATERIALIZED VIEW `mv_on_view` (`k`, `s`)
DISTRIBUTED BY HASH(`k`) BUCKETS 1 
REFRESH MANUAL
PROPERTIES (
"replicated_storage" = "true",
"replication_num" = "1",
"storage_medium" = "HDD"
)
AS SELECT `z_parent_view`.`k`, sum(`z_parent_view`.`v`) AS `s`
FROM `omni_dump_corpus`.`z_parent_view`
GROUP BY `z_parent_view`.`k`;

