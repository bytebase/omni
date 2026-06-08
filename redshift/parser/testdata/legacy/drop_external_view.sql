-- DROP EXTERNAL VIEW command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_EXTERNAL_VIEW.html

-- Basic DROP EXTERNAL VIEW syntax from AWS documentation
DROP EXTERNAL VIEW sample_schema.glue_data_catalog_view;

-- DROP EXTERNAL VIEW with IF EXISTS clause from AWS documentation
DROP EXTERNAL VIEW IF EXISTS sample_schema.glue_data_catalog_view;

-- Basic DROP EXTERNAL VIEW without schema qualification
DROP EXTERNAL VIEW simple_view;

-- DROP EXTERNAL VIEW with IF EXISTS for simple view
DROP EXTERNAL VIEW IF EXISTS simple_view;

-- DROP EXTERNAL VIEW with two-part names (schema.view)
DROP EXTERNAL VIEW analytics_schema.customer_view;
DROP EXTERNAL VIEW reporting_schema.sales_view;
DROP EXTERNAL VIEW data_warehouse.inventory_view;
DROP EXTERNAL VIEW spectrum_schema.product_view;

-- DROP EXTERNAL VIEW with IF EXISTS for two-part names
DROP EXTERNAL VIEW IF EXISTS analytics_schema.customer_view;
DROP EXTERNAL VIEW IF EXISTS reporting_schema.sales_view;
DROP EXTERNAL VIEW IF EXISTS data_warehouse.inventory_view;
DROP EXTERNAL VIEW IF EXISTS spectrum_schema.product_view;

-- DROP EXTERNAL VIEW with three-part names (catalog.schema.view)
DROP EXTERNAL VIEW my_catalog.analytics_schema.sales_view;
DROP EXTERNAL VIEW "data_catalog".reporting.monthly_metrics;
DROP EXTERNAL VIEW external_catalog.public.customer_data;

-- DROP EXTERNAL VIEW with awsdatacatalog prefix (three-part AWS format)
DROP EXTERNAL VIEW awsdatacatalog.analytics_db.customer_view;
DROP EXTERNAL VIEW awsdatacatalog.sales_database.transaction_view;
DROP EXTERNAL VIEW awsdatacatalog.inventory_db.stock_view;

-- DROP EXTERNAL VIEW with IF EXISTS for three-part names
DROP EXTERNAL VIEW IF EXISTS my_catalog.analytics_schema.sales_view;
DROP EXTERNAL VIEW IF EXISTS awsdatacatalog.analytics_db.customer_view;
DROP EXTERNAL VIEW IF EXISTS "data_catalog".reporting.monthly_metrics;

-- DROP EXTERNAL VIEW with quoted identifiers
DROP EXTERNAL VIEW "SampleSchema"."GlueDataCatalogView";
DROP EXTERNAL VIEW "my-schema"."view-with-hyphens";
DROP EXTERNAL VIEW "Schema With Spaces"."View With Spaces";

-- DROP EXTERNAL VIEW with mixed quoted and unquoted identifiers
DROP EXTERNAL VIEW analytics_schema."Customer View";
DROP EXTERNAL VIEW "reporting-schema".sales_view;
DROP EXTERNAL VIEW awsdatacatalog."my-database"."special-view";

-- DROP EXTERNAL VIEW with IF EXISTS and quoted identifiers
DROP EXTERNAL VIEW IF EXISTS "SampleSchema"."GlueDataCatalogView";
DROP EXTERNAL VIEW IF EXISTS "my-schema"."view-with-hyphens";
DROP EXTERNAL VIEW IF EXISTS "Schema With Spaces"."View With Spaces";

-- DROP EXTERNAL VIEW with common external schema names
DROP EXTERNAL VIEW spectrum_schema.customer_orders_view;
DROP EXTERNAL VIEW redshift_spectrum.sales_data_view;
DROP EXTERNAL VIEW external_analytics.user_behavior_view;
DROP EXTERNAL VIEW data_lake.financial_metrics_view;

-- DROP EXTERNAL VIEW with underscored naming patterns
DROP EXTERNAL VIEW customer_analytics.lifetime_value_view;
DROP EXTERNAL VIEW sales_reporting.monthly_summary_view;
DROP EXTERNAL VIEW inventory_management.stock_alerts_view;
DROP EXTERNAL VIEW financial_reporting.quarterly_results_view;

-- DROP EXTERNAL VIEW with mixed case naming
DROP EXTERNAL VIEW DataLake.CustomerAnalytics;
DROP EXTERNAL VIEW SalesSchema.ProductPerformance;
DROP EXTERNAL VIEW AnalyticsWarehouse.CustomerSegmentation;

-- DROP EXTERNAL VIEW with IF EXISTS for various naming patterns
DROP EXTERNAL VIEW IF EXISTS spectrum_schema.customer_orders_view;
DROP EXTERNAL VIEW IF EXISTS customer_analytics.lifetime_value_view;
DROP EXTERNAL VIEW IF EXISTS DataLake.CustomerAnalytics;
DROP EXTERNAL VIEW IF EXISTS external_analytics.user_behavior_view;

-- DROP EXTERNAL VIEW with special characters in quoted names
DROP EXTERNAL VIEW "view-with-hyphens";
DROP EXTERNAL VIEW "view.with.dots";
DROP EXTERNAL VIEW "view$with$dollars";
DROP EXTERNAL VIEW "view@with@symbols";

-- DROP EXTERNAL VIEW with three-part quoted names
DROP EXTERNAL VIEW "my-catalog"."analytics-schema"."sales-view";
DROP EXTERNAL VIEW "data.catalog"."reporting.schema"."monthly.metrics";
DROP EXTERNAL VIEW awsdatacatalog."database-name"."view-name";

-- DROP EXTERNAL VIEW with complex external schema configurations
DROP EXTERNAL VIEW external_s3_data.customer_transactions;
DROP EXTERNAL VIEW glue_catalog_schema.product_inventory;
DROP EXTERNAL VIEW redshift_spectrum_db.order_analytics;
DROP EXTERNAL VIEW cross_platform_data.unified_customer_view;

-- DROP EXTERNAL VIEW with IF EXISTS for complex configurations
DROP EXTERNAL VIEW IF EXISTS external_s3_data.customer_transactions;
DROP EXTERNAL VIEW IF EXISTS glue_catalog_schema.product_inventory;
DROP EXTERNAL VIEW IF EXISTS redshift_spectrum_db.order_analytics;
DROP EXTERNAL VIEW IF EXISTS cross_platform_data.unified_customer_view;

-- DROP EXTERNAL VIEW with realistic business scenarios
DROP EXTERNAL VIEW sales_analytics.regional_performance_view;
DROP EXTERNAL VIEW customer_insights.behavior_analysis_view;
DROP EXTERNAL VIEW financial_data.revenue_tracking_view;
DROP EXTERNAL VIEW operations.supply_chain_view;
DROP EXTERNAL VIEW marketing.campaign_effectiveness_view;
DROP EXTERNAL VIEW hr_analytics.employee_performance_view;

-- DROP EXTERNAL VIEW with database-specific naming
DROP EXTERNAL VIEW awsdatacatalog.ecommerce_db.customer_orders;
DROP EXTERNAL VIEW awsdatacatalog.analytics_warehouse.sales_metrics;
DROP EXTERNAL VIEW awsdatacatalog.operational_data.inventory_levels;
DROP EXTERNAL VIEW awsdatacatalog.financial_reporting.quarterly_summary;

-- DROP EXTERNAL VIEW with IF EXISTS for business scenarios
DROP EXTERNAL VIEW IF EXISTS sales_analytics.regional_performance_view;
DROP EXTERNAL VIEW IF EXISTS awsdatacatalog.ecommerce_db.customer_orders;
DROP EXTERNAL VIEW IF EXISTS customer_insights.behavior_analysis_view;
DROP EXTERNAL VIEW IF EXISTS financial_data.revenue_tracking_view;

-- DROP EXTERNAL VIEW with schema names that might be keywords
DROP EXTERNAL VIEW external.view_name;
DROP EXTERNAL VIEW schema.external_view;
DROP EXTERNAL VIEW catalog.schema_view;
DROP EXTERNAL VIEW database.table_view;
DROP EXTERNAL VIEW analytics.data_view;
DROP EXTERNAL VIEW reporting.warehouse_view;

-- DROP EXTERNAL VIEW with view names that might be keywords
DROP EXTERNAL VIEW sample_schema.external;
DROP EXTERNAL VIEW sample_schema.view;
DROP EXTERNAL VIEW sample_schema.schema;
DROP EXTERNAL VIEW sample_schema.catalog;
DROP EXTERNAL VIEW sample_schema.database;
DROP EXTERNAL VIEW sample_schema.analytics;
DROP EXTERNAL VIEW sample_schema.reporting;
DROP EXTERNAL VIEW sample_schema.warehouse;

-- DROP EXTERNAL VIEW statements without semicolons (should also be valid)
DROP EXTERNAL VIEW test_schema.test_view
DROP EXTERNAL VIEW IF EXISTS test_schema.conditional_view

-- DROP EXTERNAL VIEW with various AWS Data Catalog patterns
DROP EXTERNAL VIEW awsdatacatalog.default.customer_data;
DROP EXTERNAL VIEW awsdatacatalog.sampledb.sales_data;
DROP EXTERNAL VIEW awsdatacatalog.analytics.user_events;
DROP EXTERNAL VIEW awsdatacatalog.warehouse.inventory_snapshot;

-- DROP EXTERNAL VIEW examples that match CREATE EXTERNAL VIEW patterns
-- These should correspond to views that could have been created with CREATE EXTERNAL VIEW
DROP EXTERNAL VIEW sample_schema.glue_data_catalog_view;
DROP EXTERNAL VIEW awsdatacatalog.analytics_db.customer_view;
DROP EXTERNAL VIEW my_catalog.analytics_schema.sales_view;
DROP EXTERNAL VIEW spectrum_schema.aggregated_sales_view;
DROP EXTERNAL VIEW analytics_schema.customer_orders_view;
DROP EXTERNAL VIEW reporting_schema.detailed_order_view;
DROP EXTERNAL VIEW analytics_schema.high_value_customers_view;
DROP EXTERNAL VIEW analytics_schema.customer_segmentation_view;
DROP EXTERNAL VIEW analytics_schema.sales_ranking_view;
DROP EXTERNAL VIEW reporting_schema.monthly_sales_view;
DROP EXTERNAL VIEW analytics_schema.customer_contact_view;
DROP EXTERNAL VIEW reporting_schema.all_transactions_view;
DROP EXTERNAL VIEW analytics_schema.customers_with_orders_view;
DROP EXTERNAL VIEW analytics_schema.customers_without_orders_view;
DROP EXTERNAL VIEW analytics_schema.premium_product_sales_view;
DROP EXTERNAL VIEW reporting_schema.recent_orders_view;
DROP EXTERNAL VIEW analytics_schema.tech_customers_view;
DROP EXTERNAL VIEW analytics_schema.complete_customer_profiles_view;
DROP EXTERNAL VIEW analytics_schema.order_metrics_view;
DROP EXTERNAL VIEW reporting_schema.product_performance_view;
DROP EXTERNAL VIEW external_analytics.regional_sales_summary;
DROP EXTERNAL VIEW data_warehouse.customer_lifetime_value;
DROP EXTERNAL VIEW DataLake.CustomerAnalytics;
DROP EXTERNAL VIEW secure_analytics.financial_summary;
DROP EXTERNAL VIEW cross_platform.unified_customer_data;
DROP EXTERNAL VIEW spectrum_database.public.inventory_status;

-- DROP EXTERNAL VIEW with IF EXISTS for views that match CREATE patterns
DROP EXTERNAL VIEW IF EXISTS sample_schema.glue_data_catalog_view;
DROP EXTERNAL VIEW IF EXISTS awsdatacatalog.analytics_db.customer_view;
DROP EXTERNAL VIEW IF EXISTS spectrum_schema.aggregated_sales_view;
DROP EXTERNAL VIEW IF EXISTS analytics_schema.customer_segmentation_view;
DROP EXTERNAL VIEW IF EXISTS cross_platform.unified_customer_data;