-- REFRESH MATERIALIZED VIEW command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/materialized-view-refresh-sql-command.html

-- Basic REFRESH MATERIALIZED VIEW statements
REFRESH MATERIALIZED VIEW tickets_mv;
refresh materialized view tickets_mv;

-- With schema qualification
REFRESH MATERIALIZED VIEW myschema.products_mv;
refresh materialized view sales_schema.monthly_sales_mv;

-- With quoted identifiers
REFRESH MATERIALIZED VIEW "MyMaterializedView";
REFRESH MATERIALIZED VIEW "my schema"."my view";

-- REFRESH with RESTRICT option (default behavior)
REFRESH MATERIALIZED VIEW tickets_mv RESTRICT;
refresh materialized view products_mv restrict;
REFRESH MATERIALIZED VIEW myschema.sales_mv RESTRICT;

-- REFRESH with CASCADE option
REFRESH MATERIALIZED VIEW products_mv CASCADE;
refresh materialized view tickets_mv cascade;
REFRESH MATERIALIZED VIEW sales_schema.base_mv CASCADE;

-- Mixed case variations
Refresh Materialized View my_mv;
REFRESH materialized VIEW my_mv;
refresh MATERIALIZED view my_mv;
REFRESH MATERIALIZED VIEW My_Mv;

-- Examples from AWS documentation
-- Basic refresh for a tickets materialized view
REFRESH MATERIALIZED VIEW tickets_mv;

-- Cascading refresh for dependent materialized views
REFRESH MATERIALIZED VIEW products_mv CASCADE;

-- Real-world scenarios
-- Refresh a materialized view that aggregates sales data
REFRESH MATERIALIZED VIEW daily_sales_summary;

-- Refresh with cascade for nested materialized views
-- where sales_summary_mv depends on other MVs
REFRESH MATERIALIZED VIEW sales_summary_mv CASCADE;

-- Refresh materialized views in different schemas
REFRESH MATERIALIZED VIEW analytics.customer_metrics;
REFRESH MATERIALIZED VIEW reporting.revenue_summary;
REFRESH MATERIALIZED VIEW dwh.fact_sales_mv;

-- Refresh with RESTRICT to avoid refreshing dependent views
REFRESH MATERIALIZED VIEW base_customers_mv RESTRICT;
REFRESH MATERIALIZED VIEW raw_data.source_mv RESTRICT;

-- Complex schema and view names
REFRESH MATERIALIZED VIEW "2024_sales"."Q1_summary";
REFRESH MATERIALIZED VIEW "sales-data"."monthly_metrics_v2";
REFRESH MATERIALIZED VIEW prod_schema."USER_ACTIVITY_MV";

-- Edge cases with special characters in names
REFRESH MATERIALIZED VIEW "my-mv-with-dashes";
REFRESH MATERIALIZED VIEW "MV_With_Mixed_Case";
REFRESH MATERIALIZED VIEW "mv with spaces";

-- Multiple refresh statements in sequence
REFRESH MATERIALIZED VIEW base_table_mv RESTRICT;
REFRESH MATERIALIZED VIEW dependent_mv_1 RESTRICT;
REFRESH MATERIALIZED VIEW dependent_mv_2 RESTRICT;
REFRESH MATERIALIZED VIEW top_level_mv CASCADE;

-- Refresh materialized views for data lake integrations
REFRESH MATERIALIZED VIEW lake_house.external_data_mv;
REFRESH MATERIALIZED VIEW zero_etl.integration_mv CASCADE;

-- Refresh materialized views for different use cases
REFRESH MATERIALIZED VIEW inventory.stock_levels;
REFRESH MATERIALIZED VIEW marketing.campaign_performance;
REFRESH MATERIALIZED VIEW finance.quarterly_reports CASCADE;
REFRESH MATERIALIZED VIEW operations.kpi_dashboard RESTRICT;