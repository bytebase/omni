CREATE MATERIALIZED VIEW complete_mv (orderdate COMMENT 'Order date', orderkey COMMENT 'Order key', partkey COMMENT 'Part key')
BUILD IMMEDIATE
REFRESH AUTO
ON SCHEDULE EVERY 1 DAY STARTS '2024-12-01 20:30:00'
DISTRIBUTED BY HASH (orderkey) BUCKETS 2
PROPERTIES("replication_num" = "1")
AS
SELECT
    o_orderdate,
    l_orderkey,
    l_partkey
FROM orders
LEFT JOIN lineitem ON l_orderkey = o_orderkey
LEFT JOIN partsupp ON ps_partkey = l_partkey AND l_suppkey = ps_suppkey;

CREATE MATERIALIZED VIEW partition_mv
BUILD IMMEDIATE
REFRESH AUTO
ON SCHEDULE EVERY 1 DAY STARTS '2024-12-01 20:30:00'
PARTITION BY (DATE_TRUNC(o_orderdate, 'MONTH'))
DISTRIBUTED BY HASH (l_orderkey) BUCKETS 2
PROPERTIES("replication_num" = "3")
AS
SELECT
    o_orderdate,
    l_orderkey,
    l_partkey
FROM orders
LEFT JOIN lineitem ON l_orderkey = o_orderkey
LEFT JOIN partsupp ON ps_partkey = l_partkey AND l_suppkey = ps_suppkey;

ALTER MATERIALIZED VIEW partition_mv SET ("grace_period" = "10", "excluded_trigger_tables" = "lineitem,partsupp");
ALTER MATERIALIZED VIEW partition_mv REFRESH COMPLETE;
DROP MATERIALIZED VIEW mv1;
DROP MATERIALIZED VIEW IF EXISTS db1.mv1;
REFRESH MATERIALIZED VIEW mv1 AUTO;
REFRESH MATERIALIZED VIEW mv1 PARTITIONS(p_19950801_19950901, p_19950901_19951001);
REFRESH MATERIALIZED VIEW mv1 COMPLETE
