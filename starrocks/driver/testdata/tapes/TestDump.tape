-- pin: default starrocks/allin1-ubuntu:3.4.10 sha256:8bc803ea7a5bb490c77117a456b32b792ccae698b1f0e47dfcafdac41542566d 3.4.10

== TestDump default#1
> exec
| -- The shapes a StarRocks schema dump must carry. Each table states a key model,
| -- a distribution, and properties that information_schema does not hold; the
| -- views and materialized views state definitions whose order matters.
|
| CREATE TABLE duplicate_keys (
| 	k INT NOT NULL COMMENT 'key column',
| 	v VARCHAR(20) NULL COMMENT 'a value',
| 	d DATETIME NULL,
| 	w BIGINT AS (k + 1),
| 	INDEX duplicate_keys_v (v) USING NGRAMBF
| ) ENGINE=OLAP
| DUPLICATE KEY(k)
| COMMENT 'duplicate key table'
| DISTRIBUTED BY HASH(k) BUCKETS 3
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
|
|
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
|
|
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
|
|
| -- The server keeps the partition bounds as written, so the dump reads them from
| -- SHOW CREATE TABLE rather than from information_schema.
| CREATE TABLE range_partitions (
| 	d DATE NOT NULL,
| 	v INT NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(d)
| PARTITION BY RANGE(d) (
| 	PARTITION p2025 VALUES [('2025-01-01'), ('2026-01-01')),
| 	PARTITION p2026 VALUES [('2026-01-01'), ('2027-01-01'))
| )
| DISTRIBUTED BY HASH(d) BUCKETS 1
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
|
|
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
|
|
| -- Dynamic partitioning is a set of properties, not a partition list: the dump
| -- must carry the properties that maintain the window.
| CREATE TABLE dynamic_partitions (
| 	d DATE NOT NULL,
| 	v INT NULL
| ) ENGINE=OLAP
| UNIQUE KEY(d)
| PARTITION BY RANGE(d) ()
| DISTRIBUTED BY HASH(d)
| PROPERTIES("replication_num" = "1", "dynamic_partition.enable" = "true", "dynamic_partition.time_unit" = "DAY", "dynamic_partition.start" = "-2", "dynamic_partition.end" = "2", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "2");
< ok 0 0
> exec
|
|
| CREATE TABLE materialized_view_base (
| 	k INT NOT NULL,
| 	v INT NULL
| ) ENGINE=OLAP
| DUPLICATE KEY(k)
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1");
< ok 0 0
> exec
|
|
| -- The names invert the dependency order: a_child_mv reads z_parent_mv, and
| -- a_child_view reads z_parent_view, so a dump that wrote them in name order
| -- would not restore.
| CREATE MATERIALIZED VIEW z_parent_mv
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1")
| AS SELECT k, sum(v) AS s FROM materialized_view_base GROUP BY k;
< ok 0 0
> exec
|
|
| CREATE MATERIALIZED VIEW a_child_mv
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1")
| AS SELECT k, sum(s) AS s2 FROM z_parent_mv GROUP BY k;
< ok 0 0
> exec
|
|
| CREATE VIEW z_parent_view (k, v) AS SELECT k, v FROM materialized_view_base;
< ok 0 0
> exec
|
|
| CREATE VIEW a_child_view (k) AS SELECT k FROM z_parent_view;
< ok 0 0
> exec
|
|
| -- A view reads an asynchronous materialized view, and a materialized view reads
| -- a view: either kind may follow the other.
| CREATE VIEW b_mv_view (k, s) AS SELECT k, s FROM z_parent_mv;
< ok 0 0
> exec
|
|
| CREATE MATERIALIZED VIEW mv_on_view
| DISTRIBUTED BY HASH(k) BUCKETS 1
| PROPERTIES("replication_num" = "1")
| AS SELECT k, sum(v) AS s FROM z_parent_view GROUP BY k;
< ok 0 0
> exec
|
|
| -- A view body is a query, not a bare SELECT: a set operation, a WITH clause and
| -- a scalar subquery each read z_parent_view, and each is alphabetically before
| -- it, so the dump must see the reference inside the form to order them.
| CREATE VIEW aa_union_view (k) AS (SELECT k FROM z_parent_view) UNION ALL (SELECT k FROM z_parent_view);
< ok 0 0
> exec
|
|
| CREATE VIEW ab_cte_view (k) AS WITH c AS (SELECT k FROM z_parent_view) SELECT k FROM c;
< ok 0 0
> exec
|
|
| CREATE VIEW ac_subquery_view (k) AS SELECT (SELECT max(k) FROM z_parent_view) AS k;
< ok 0 0
> exec
|
|
| -- A CTE that shares a name with a dumped view is not a dependency on it:
| -- cyc_parent's WITH clause defines cyc_child, and the view cyc_child reads
| -- cyc_parent. Taking the CTE for an object makes the pair a reference cycle,
| -- which falls back to name order, cyc_child first, and cannot restore.
| CREATE VIEW cyc_parent (k) AS WITH cyc_child AS (SELECT 1 AS k) SELECT k FROM cyc_child;
< ok 0 0
> exec
|
|
| CREATE VIEW cyc_child (k) AS SELECT k FROM cyc_parent;
< ok 0 0
> query
| SELECT SCHEMA_NAME FROM information_schema.schemata WHERE SCHEMA_NAME = ?
: string
| omni_dump_corpus
< result
< column "SCHEMA_NAME" "VARCHAR"
< row
: bytes
| omni_dump_corpus
> query
| SELECT TABLE_NAME FROM information_schema.materialized_views WHERE TABLE_SCHEMA = ? AND IFNULL(REFRESH_TYPE, '') <> 'ROLLUP' ORDER BY TABLE_NAME
: string
| omni_dump_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< row
: bytes
| a_child_mv
< row
: bytes
| mv_on_view
< row
: bytes
| z_parent_mv
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.tables WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME
: string
| omni_dump_corpus
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
< row
: bytes
| a_child_mv
: bytes
| VIEW
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
| mv_on_view
: bytes
| VIEW
< row
: bytes
| range_partitions
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
| VIEW
< row
: bytes
| z_parent_view
: bytes
| VIEW
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_corpus`.`a_child_mv`
< result
< column "Materialized View" "VARCHAR"
< column "Create Materialized View" "VARCHAR"
< row
: bytes
| a_child_mv
: bytes
| CREATE MATERIALIZED VIEW `a_child_mv` (`k`, `s2`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1 
| REFRESH MANUAL
| PROPERTIES (
| "replicated_storage" = "true",
| "replication_num" = "1",
| "storage_medium" = "HDD"
| )
| AS SELECT `z_parent_mv`.`k`, sum(`z_parent_mv`.`s`) AS `s2`
| FROM `omni_dump_corpus`.`z_parent_mv`
| GROUP BY `z_parent_mv`.`k`;
> query
| SHOW CREATE VIEW `omni_dump_corpus`.`a_child_view`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| a_child_view
: bytes
| CREATE VIEW `a_child_view` (`k`) SECURITY NONE AS SELECT `omni_dump_corpus`.`z_parent_view`.`k`
| FROM `omni_dump_corpus`.`z_parent_view`;
: bytes
| utf8
: bytes
| utf8_general_ci
> query
| SHOW CREATE VIEW `omni_dump_corpus`.`aa_union_view`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| aa_union_view
: bytes
| CREATE VIEW `aa_union_view` (`k`) SECURITY NONE AS (SELECT `omni_dump_corpus`.`z_parent_view`.`k`
| FROM `omni_dump_corpus`.`z_parent_view`) UNION ALL (SELECT `omni_dump_corpus`.`z_parent_view`.`k`
| FROM `omni_dump_corpus`.`z_parent_view`);
: bytes
| utf8
: bytes
| utf8_general_ci
> query
| SHOW CREATE VIEW `omni_dump_corpus`.`ab_cte_view`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ab_cte_view
: bytes
| CREATE VIEW `ab_cte_view` (`k`) SECURITY NONE AS WITH `c` (`k`) AS (SELECT `omni_dump_corpus`.`z_parent_view`.`k`
| FROM `omni_dump_corpus`.`z_parent_view`) SELECT `c`.`k`
| FROM `c`;
: bytes
| utf8
: bytes
| utf8_general_ci
> query
| SHOW CREATE VIEW `omni_dump_corpus`.`ac_subquery_view`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| ac_subquery_view
: bytes
| CREATE VIEW `ac_subquery_view` (`k`) SECURITY NONE AS SELECT (SELECT max(`omni_dump_corpus`.`z_parent_view`.`k`) AS `max(k)`
| FROM `omni_dump_corpus`.`z_parent_view`) AS `k`;
: bytes
| utf8
: bytes
| utf8_general_ci
> query
| SHOW CREATE TABLE `omni_dump_corpus`.`aggregate_keys`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| aggregate_keys
: bytes
| CREATE TABLE `aggregate_keys` (
|   `k` int(11) NOT NULL COMMENT "",
|   `s` int(11) SUM NULL COMMENT "",
|   `m` int(11) MAX NULL COMMENT "",
|   `r` varchar(10) REPLACE NULL COMMENT ""
| ) ENGINE=OLAP 
| AGGREGATE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 2 
| PROPERTIES (
| "compression" = "LZ4",
| "fast_schema_evolution" = "true",
| "replicated_storage" = "true",
| "replication_num" = "1"
| );
> query
| SHOW CREATE VIEW `omni_dump_corpus`.`b_mv_view`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| b_mv_view
: bytes
| CREATE VIEW `b_mv_view` (`k`, `s`) SECURITY NONE AS SELECT `omni_dump_corpus`.`z_parent_mv`.`k`, `omni_dump_corpus`.`z_parent_mv`.`s`
| FROM `omni_dump_corpus`.`z_parent_mv`;
: bytes
| utf8
: bytes
| utf8_general_ci
> query
| SHOW CREATE VIEW `omni_dump_corpus`.`cyc_child`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| cyc_child
: bytes
| CREATE VIEW `cyc_child` (`k`) SECURITY NONE AS SELECT `omni_dump_corpus`.`cyc_parent`.`k`
| FROM `omni_dump_corpus`.`cyc_parent`;
: bytes
| utf8
: bytes
| utf8_general_ci
> query
| SHOW CREATE VIEW `omni_dump_corpus`.`cyc_parent`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| cyc_parent
: bytes
| CREATE VIEW `cyc_parent` (`k`) SECURITY NONE AS WITH `cyc_child` (`k`) AS (SELECT 1 AS `k`) SELECT `cyc_child`.`k`
| FROM `cyc_child`;
: bytes
| utf8
: bytes
| utf8_general_ci
> query
| SHOW CREATE TABLE `omni_dump_corpus`.`duplicate_keys`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| duplicate_keys
: bytes
| CREATE TABLE `duplicate_keys` (
|   `k` int(11) NOT NULL COMMENT "key column",
|   `v` varchar(20) NULL COMMENT "a value",
|   `d` datetime NULL COMMENT "",
|   `w` bigint(20) NULL AS `k` + 1 COMMENT "",
|   INDEX duplicate_keys_v (`v`) USING NGRAMBF("bloom_filter_fpp" = "0.05", "case_sensitive" = "true", "gram_num" = "2") COMMENT ''
| ) ENGINE=OLAP 
| DUPLICATE KEY(`k`)
| COMMENT "duplicate key table"
| DISTRIBUTED BY HASH(`k`) BUCKETS 3 
| PROPERTIES (
| "compression" = "LZ4",
| "fast_schema_evolution" = "true",
| "replicated_storage" = "true",
| "replication_num" = "1"
| );
> query
| SHOW CREATE TABLE `omni_dump_corpus`.`dynamic_partitions`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| dynamic_partitions
: bytes
| CREATE TABLE `dynamic_partitions` (
|   `d` date NOT NULL COMMENT "",
|   `v` int(11) NULL COMMENT ""
| ) ENGINE=OLAP 
| UNIQUE KEY(`d`)
| PARTITION BY RANGE(`d`)
| (PARTITION p20261117 VALUES [("2026-11-17"), ("2026-11-18")),
| PARTITION p20261118 VALUES [("2026-11-18"), ("2026-11-19")),
| PARTITION p20261119 VALUES [("2026-11-19"), ("2026-11-20")))
| DISTRIBUTED BY HASH(`d`)
| PROPERTIES (
| "compression" = "LZ4",
| "dynamic_partition.buckets" = "2",
| "dynamic_partition.enable" = "true",
| "dynamic_partition.end" = "2",
| "dynamic_partition.history_partition_num" = "0",
| "dynamic_partition.prefix" = "p",
| "dynamic_partition.start" = "-2",
| "dynamic_partition.time_unit" = "DAY",
| "dynamic_partition.time_zone" = "Etc/UTC",
| "fast_schema_evolution" = "true",
| "replicated_storage" = "true",
| "replication_num" = "1"
| );
> query
| SHOW CREATE TABLE `omni_dump_corpus`.`list_partitions`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| list_partitions
: bytes
| CREATE TABLE `list_partitions` (
|   `k` int(11) NOT NULL COMMENT "",
|   `v` int(11) NULL COMMENT ""
| ) ENGINE=OLAP 
| DUPLICATE KEY(`k`)
| PARTITION BY LIST(`k`)(
|   PARTITION p1 VALUES IN ('1', '2'),
|   PARTITION p2 VALUES IN ('3')
| )
| DISTRIBUTED BY HASH(`k`) BUCKETS 1 
| PROPERTIES (
| "compression" = "LZ4",
| "fast_schema_evolution" = "true",
| "replicated_storage" = "true",
| "replication_num" = "1"
| );
> query
| SHOW CREATE TABLE `omni_dump_corpus`.`materialized_view_base`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| materialized_view_base
: bytes
| CREATE TABLE `materialized_view_base` (
|   `k` int(11) NOT NULL COMMENT "",
|   `v` int(11) NULL COMMENT ""
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
| SHOW CREATE MATERIALIZED VIEW `omni_dump_corpus`.`mv_on_view`
< result
< column "Materialized View" "VARCHAR"
< column "Create Materialized View" "VARCHAR"
< row
: bytes
| mv_on_view
: bytes
| CREATE MATERIALIZED VIEW `mv_on_view` (`k`, `s`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1 
| REFRESH MANUAL
| PROPERTIES (
| "replicated_storage" = "true",
| "replication_num" = "1",
| "storage_medium" = "HDD"
| )
| AS SELECT `z_parent_view`.`k`, sum(`z_parent_view`.`v`) AS `s`
| FROM `omni_dump_corpus`.`z_parent_view`
| GROUP BY `z_parent_view`.`k`;
> query
| SHOW CREATE TABLE `omni_dump_corpus`.`range_partitions`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| range_partitions
: bytes
| CREATE TABLE `range_partitions` (
|   `d` date NOT NULL COMMENT "",
|   `v` int(11) NULL COMMENT ""
| ) ENGINE=OLAP 
| DUPLICATE KEY(`d`)
| PARTITION BY RANGE(`d`)
| (PARTITION p2025 VALUES [("2025-01-01"), ("2026-01-01")),
| PARTITION p2026 VALUES [("2026-01-01"), ("2027-01-01")))
| DISTRIBUTED BY HASH(`d`) BUCKETS 1 
| PROPERTIES (
| "compression" = "LZ4",
| "fast_schema_evolution" = "true",
| "replicated_storage" = "true",
| "replication_num" = "1"
| );
> query
| SHOW CREATE TABLE `omni_dump_corpus`.`unique_keys`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| unique_keys
: bytes
| CREATE TABLE `unique_keys` (
|   `k` int(11) NOT NULL COMMENT "",
|   `v` int(11) NULL COMMENT "",
|   `g` bigint(20) NOT NULL AUTO_INCREMENT COMMENT ""
| ) ENGINE=OLAP 
| UNIQUE KEY(`k`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 2 
| PROPERTIES (
| "compression" = "LZ4",
| "fast_schema_evolution" = "true",
| "replicated_storage" = "true",
| "replication_num" = "1"
| );
> query
| SHOW CREATE MATERIALIZED VIEW `omni_dump_corpus`.`z_parent_mv`
< result
< column "Materialized View" "VARCHAR"
< column "Create Materialized View" "VARCHAR"
< row
: bytes
| z_parent_mv
: bytes
| CREATE MATERIALIZED VIEW `z_parent_mv` (`k`, `s`)
| DISTRIBUTED BY HASH(`k`) BUCKETS 1 
| REFRESH MANUAL
| PROPERTIES (
| "replicated_storage" = "true",
| "replication_num" = "1",
| "storage_medium" = "HDD"
| )
| AS SELECT `materialized_view_base`.`k`, sum(`materialized_view_base`.`v`) AS `s`
| FROM `omni_dump_corpus`.`materialized_view_base`
| GROUP BY `materialized_view_base`.`k`;
> query
| SHOW CREATE VIEW `omni_dump_corpus`.`z_parent_view`
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| z_parent_view
: bytes
| CREATE VIEW `z_parent_view` (`k`, `v`) SECURITY NONE AS SELECT `omni_dump_corpus`.`materialized_view_base`.`k`, `omni_dump_corpus`.`materialized_view_base`.`v`
| FROM `omni_dump_corpus`.`materialized_view_base`;
: bytes
| utf8
: bytes
| utf8_general_ci
