-- The shapes a StarRocks schema dump must carry. Each table states a key model,
-- a distribution, and properties that information_schema does not hold; the
-- views and materialized views state definitions whose order matters.

CREATE TABLE duplicate_keys (
	k INT NOT NULL COMMENT 'key column',
	v VARCHAR(20) NULL COMMENT 'a value',
	d DATETIME NULL,
	w BIGINT AS (k + 1),
	INDEX duplicate_keys_v (v) USING NGRAMBF
) ENGINE=OLAP
DUPLICATE KEY(k)
COMMENT 'duplicate key table'
DISTRIBUTED BY HASH(k) BUCKETS 3
PROPERTIES("replication_num" = "1");

CREATE TABLE unique_keys (
	k INT NOT NULL,
	v INT NULL,
	g BIGINT NOT NULL AUTO_INCREMENT
) ENGINE=OLAP
UNIQUE KEY(k)
DISTRIBUTED BY HASH(k) BUCKETS 2
PROPERTIES("replication_num" = "1");

CREATE TABLE aggregate_keys (
	k INT NOT NULL,
	s INT SUM NULL,
	m INT MAX NULL,
	r VARCHAR(10) REPLACE NULL
) ENGINE=OLAP
AGGREGATE KEY(k)
DISTRIBUTED BY HASH(k) BUCKETS 2
PROPERTIES("replication_num" = "1");

-- The server keeps the partition bounds as written, so the dump reads them from
-- SHOW CREATE TABLE rather than from information_schema.
CREATE TABLE range_partitions (
	d DATE NOT NULL,
	v INT NULL
) ENGINE=OLAP
DUPLICATE KEY(d)
PARTITION BY RANGE(d) (
	PARTITION p2025 VALUES [('2025-01-01'), ('2026-01-01')),
	PARTITION p2026 VALUES [('2026-01-01'), ('2027-01-01'))
)
DISTRIBUTED BY HASH(d) BUCKETS 1
PROPERTIES("replication_num" = "1");

CREATE TABLE list_partitions (
	k INT NOT NULL,
	v INT NULL
) ENGINE=OLAP
DUPLICATE KEY(k)
PARTITION BY LIST(k) (PARTITION p1 VALUES IN ('1', '2'), PARTITION p2 VALUES IN ('3'))
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1");

-- Dynamic partitioning is a set of properties, not a partition list: the dump
-- must carry the properties that maintain the window.
CREATE TABLE dynamic_partitions (
	d DATE NOT NULL,
	v INT NULL
) ENGINE=OLAP
UNIQUE KEY(d)
PARTITION BY RANGE(d) ()
DISTRIBUTED BY HASH(d)
PROPERTIES("replication_num" = "1", "dynamic_partition.enable" = "true", "dynamic_partition.time_unit" = "DAY", "dynamic_partition.start" = "-2", "dynamic_partition.end" = "2", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "2");

CREATE TABLE materialized_view_base (
	k INT NOT NULL,
	v INT NULL
) ENGINE=OLAP
DUPLICATE KEY(k)
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1");

-- The names invert the dependency order: a_child_mv reads z_parent_mv, and
-- a_child_view reads z_parent_view, so a dump that wrote them in name order
-- would not restore.
CREATE MATERIALIZED VIEW z_parent_mv
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1")
AS SELECT k, sum(v) AS s FROM materialized_view_base GROUP BY k;

CREATE MATERIALIZED VIEW a_child_mv
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1")
AS SELECT k, sum(s) AS s2 FROM z_parent_mv GROUP BY k;

CREATE VIEW z_parent_view (k, v) AS SELECT k, v FROM materialized_view_base;

CREATE VIEW a_child_view (k) AS SELECT k FROM z_parent_view;

-- A view reads an asynchronous materialized view, and a materialized view reads
-- a view: either kind may follow the other.
CREATE VIEW b_mv_view (k, s) AS SELECT k, s FROM z_parent_mv;

CREATE MATERIALIZED VIEW mv_on_view
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1")
AS SELECT k, sum(v) AS s FROM z_parent_view GROUP BY k;

-- A view body is a query, not a bare SELECT: a set operation, a WITH clause and
-- a scalar subquery each read z_parent_view, and each is alphabetically before
-- it, so the dump must see the reference inside the form to order them.
CREATE VIEW aa_union_view (k) AS (SELECT k FROM z_parent_view) UNION ALL (SELECT k FROM z_parent_view);

CREATE VIEW ab_cte_view (k) AS WITH c AS (SELECT k FROM z_parent_view) SELECT k FROM c;

CREATE VIEW ac_subquery_view (k) AS SELECT (SELECT max(k) FROM z_parent_view) AS k;

-- A CTE that shares a name with a dumped view is not a dependency on it:
-- cyc_parent's WITH clause defines cyc_child, and the view cyc_child reads
-- cyc_parent. Taking the CTE for an object makes the pair a reference cycle,
-- which falls back to name order, cyc_child first, and cannot restore.
CREATE VIEW cyc_parent (k) AS WITH cyc_child AS (SELECT 1 AS k) SELECT k FROM cyc_child;

CREATE VIEW cyc_child (k) AS SELECT k FROM cyc_parent;
