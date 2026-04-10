CREATE INTERACTIVE MATERIALIZED VIEW IF NOT EXISTS mv_summary
    AS
    SELECT SUM(quantity) AS total_quantity, SUM(net_paid) AS total_net_paid
    FROM my_interactive_table
    WHERE call_center_id = 52;

ALTER WAREHOUSE interactive_wh ADD TABLES (mv_summary, my_interactive_table);
