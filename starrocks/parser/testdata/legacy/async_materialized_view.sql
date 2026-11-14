ALTER MATERIALIZED VIEW partition_mv SET ("grace_period" = "10", "excluded_trigger_tables" = "lineitem,partsupp");
DROP MATERIALIZED VIEW mv1;
DROP MATERIALIZED VIEW IF EXISTS db1.mv1;
