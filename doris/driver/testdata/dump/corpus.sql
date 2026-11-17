-- The shapes a Doris schema dump must carry. Each table states a key model,
-- a distribution, and properties that information_schema does not hold; the
-- views and materialized views state definitions whose order matters; and the
-- procedure's body holds semicolons.

CREATE TABLE duplicate_keys (
	k INT NOT NULL COMMENT 'key column',
	v VARCHAR(20) NULL COMMENT 'a value',
	d DATETIME NULL,
	w INT GENERATED ALWAYS AS (k + 1)
) ENGINE=OLAP
DUPLICATE KEY(k)
COMMENT 'duplicate key table'
DISTRIBUTED BY HASH(k) BUCKETS 3
PROPERTIES("replication_num" = "1");

-- A table with an inline index, added after creation so the dump reads it from
-- SHOW CREATE TABLE rather than from the statement that made the table.
CREATE INDEX duplicate_keys_v ON duplicate_keys (v) USING INVERTED;

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

-- Dynamic partitioning: the server materializes a window of partitions around
-- the current date and states them with the properties that maintain them.
CREATE TABLE dynamic_partitions (
	d DATE NOT NULL,
	v INT NULL
) ENGINE=OLAP
UNIQUE KEY(d)
PARTITION BY RANGE(d) ()
DISTRIBUTED BY HASH(d) BUCKETS AUTO
PROPERTIES("replication_num" = "1", "dynamic_partition.enable" = "true", "dynamic_partition.time_unit" = "DAY", "dynamic_partition.start" = "-2", "dynamic_partition.end" = "2", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "2");

CREATE TABLE list_partitions (
	k INT NOT NULL,
	v INT NULL
) ENGINE=OLAP
DUPLICATE KEY(k)
PARTITION BY LIST(k) (PARTITION p1 VALUES IN ('1', '2'), PARTITION p2 VALUES IN ('3'))
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1");

-- A partition expression the server keeps as written.
CREATE TABLE auto_partitions (
	d DATE NOT NULL,
	v INT NULL
) ENGINE=OLAP
DUPLICATE KEY(d)
AUTO PARTITION BY RANGE (date_trunc(d, 'day')) ()
DISTRIBUTED BY HASH(d) BUCKETS AUTO
PROPERTIES("replication_num" = "1");

CREATE TABLE materialized_view_base (
	k INT NOT NULL,
	v INT NULL
) ENGINE=OLAP
DUPLICATE KEY(k)
DISTRIBUTED BY RANDOM BUCKETS 2
PROPERTIES("replication_num" = "1");

-- The names invert the dependency order: a_child_mv reads z_parent_mv, and
-- a_child_view reads z_parent_view, so a dump that wrote them in name order
-- would not restore.
CREATE MATERIALIZED VIEW z_parent_mv
BUILD DEFERRED REFRESH COMPLETE ON MANUAL
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1")
AS SELECT k, sum(v) AS s FROM materialized_view_base GROUP BY k;

CREATE MATERIALIZED VIEW a_child_mv
BUILD DEFERRED REFRESH COMPLETE ON MANUAL
DISTRIBUTED BY HASH(k) BUCKETS 1
PROPERTIES("replication_num" = "1")
AS SELECT k, sum(s) AS s2 FROM z_parent_mv GROUP BY k;

CREATE VIEW z_parent_view (k, v) AS SELECT k, v FROM materialized_view_base;

CREATE VIEW a_child_view (k) AS SELECT k FROM z_parent_view;

-- A view reads a materialized view, which Doris allows in this direction only.
CREATE VIEW b_mv_view (k, s) AS SELECT k, s FROM z_parent_mv;

-- A view body is a query, not a bare SELECT: a set operation, a WITH clause and
-- a scalar subquery each read z_parent_view, and each is alphabetically before
-- it, so the dump must see the reference inside the form to order them.
CREATE VIEW aa_union_view (k) AS SELECT k FROM z_parent_view UNION ALL SELECT k FROM z_parent_view;

CREATE VIEW ab_cte_view (k) AS WITH c AS (SELECT k FROM z_parent_view) SELECT k FROM c;

CREATE VIEW ac_subquery_view (k) AS SELECT (SELECT max(k) FROM z_parent_view) AS k;

-- A CTE that shares a name with a dumped view is not a dependency on it:
-- cyc_parent's WITH clause defines cyc_child, and the view cyc_child reads
-- cyc_parent. Taking the CTE for an object makes the pair a reference cycle,
-- which falls back to name order, cyc_child first, and cannot restore.
CREATE VIEW cyc_parent (k) AS WITH cyc_child AS (SELECT 1 AS k) SELECT k FROM cyc_child;

CREATE VIEW cyc_child (k) AS SELECT k FROM cyc_parent;

DELIMITER //
CREATE PROCEDURE proc_with_body(IN a INT) COMMENT 'doris procedure'
BEGIN
	SET @first = a + 1;
	SET @second = @first * 2;
END//
DELIMITER ;
