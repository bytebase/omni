ALTER MATERIALIZED VIEW partition_mv SET ("grace_period" = "10", "excluded_trigger_tables" = "lineitem,partsupp");
SHOW CREATE MATERIALIZED VIEW partition_mv
