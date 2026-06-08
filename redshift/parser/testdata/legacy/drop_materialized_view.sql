-- DROP MATERIALIZED VIEW command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/materialized-view-drop-sql-command.html

-- Basic DROP MATERIALIZED VIEW syntax from AWS documentation
DROP MATERIALIZED VIEW tickets_mv;

-- DROP MATERIALIZED VIEW with IF EXISTS clause
DROP MATERIALIZED VIEW IF EXISTS tickets_mv;

-- Basic DROP MATERIALIZED VIEW without schema qualification
DROP MATERIALIZED VIEW simple_mv;
DROP MATERIALIZED VIEW sales_summary;
DROP MATERIALIZED VIEW customer_metrics;
DROP MATERIALIZED VIEW order_analytics;
DROP MATERIALIZED VIEW product_inventory;

-- DROP MATERIALIZED VIEW with IF EXISTS for simple views
DROP MATERIALIZED VIEW IF EXISTS simple_mv;
DROP MATERIALIZED VIEW IF EXISTS sales_summary;
DROP MATERIALIZED VIEW IF EXISTS customer_metrics;
DROP MATERIALIZED VIEW IF EXISTS order_analytics;
DROP MATERIALIZED VIEW IF EXISTS product_inventory;

-- DROP MATERIALIZED VIEW with schema qualification
DROP MATERIALIZED VIEW analytics.sales_summary_mv;
DROP MATERIALIZED VIEW reporting.customer_metrics_mv;
DROP MATERIALIZED VIEW warehouse.inventory_snapshot_mv;
DROP MATERIALIZED VIEW finance.quarterly_revenue_mv;
DROP MATERIALIZED VIEW operations.daily_orders_mv;

-- DROP MATERIALIZED VIEW with IF EXISTS and schema qualification
DROP MATERIALIZED VIEW IF EXISTS analytics.sales_summary_mv;
DROP MATERIALIZED VIEW IF EXISTS reporting.customer_metrics_mv;
DROP MATERIALIZED VIEW IF EXISTS warehouse.inventory_snapshot_mv;
DROP MATERIALIZED VIEW IF EXISTS finance.quarterly_revenue_mv;
DROP MATERIALIZED VIEW IF EXISTS operations.daily_orders_mv;

-- DROP MATERIALIZED VIEW with CASCADE option
DROP MATERIALIZED VIEW tickets_mv CASCADE;
DROP MATERIALIZED VIEW sales_summary CASCADE;
DROP MATERIALIZED VIEW analytics.customer_metrics_mv CASCADE;
DROP MATERIALIZED VIEW reporting.monthly_sales_mv CASCADE;

-- DROP MATERIALIZED VIEW with RESTRICT option (explicit default)
DROP MATERIALIZED VIEW tickets_mv RESTRICT;
DROP MATERIALIZED VIEW sales_summary RESTRICT;
DROP MATERIALIZED VIEW analytics.customer_metrics_mv RESTRICT;
DROP MATERIALIZED VIEW reporting.monthly_sales_mv RESTRICT;

-- DROP MATERIALIZED VIEW with IF EXISTS and CASCADE
DROP MATERIALIZED VIEW IF EXISTS tickets_mv CASCADE;
DROP MATERIALIZED VIEW IF EXISTS sales_summary CASCADE;
DROP MATERIALIZED VIEW IF EXISTS analytics.customer_metrics_mv CASCADE;
DROP MATERIALIZED VIEW IF EXISTS reporting.monthly_sales_mv CASCADE;

-- DROP MATERIALIZED VIEW with IF EXISTS and RESTRICT
DROP MATERIALIZED VIEW IF EXISTS tickets_mv RESTRICT;
DROP MATERIALIZED VIEW IF EXISTS sales_summary RESTRICT;
DROP MATERIALIZED VIEW IF EXISTS analytics.customer_metrics_mv RESTRICT;
DROP MATERIALIZED VIEW IF EXISTS reporting.monthly_sales_mv RESTRICT;

-- DROP MATERIALIZED VIEW with multiple views (comma-separated)
DROP MATERIALIZED VIEW view1, view2, view3;
DROP MATERIALIZED VIEW sales_mv, orders_mv, customers_mv;
DROP MATERIALIZED VIEW analytics.sales_summary, analytics.customer_metrics, analytics.product_performance;

-- DROP MATERIALIZED VIEW with multiple views and IF EXISTS
DROP MATERIALIZED VIEW IF EXISTS view1, view2, view3;
DROP MATERIALIZED VIEW IF EXISTS sales_mv, orders_mv, customers_mv;
DROP MATERIALIZED VIEW IF EXISTS analytics.sales_summary, analytics.customer_metrics, analytics.product_performance;

-- DROP MATERIALIZED VIEW with multiple views and CASCADE
DROP MATERIALIZED VIEW view1, view2, view3 CASCADE;
DROP MATERIALIZED VIEW sales_mv, orders_mv, customers_mv CASCADE;
DROP MATERIALIZED VIEW analytics.sales_summary, analytics.customer_metrics, analytics.product_performance CASCADE;

-- DROP MATERIALIZED VIEW with multiple views and RESTRICT
DROP MATERIALIZED VIEW view1, view2, view3 RESTRICT;
DROP MATERIALIZED VIEW sales_mv, orders_mv, customers_mv RESTRICT;
DROP MATERIALIZED VIEW analytics.sales_summary, analytics.customer_metrics, analytics.product_performance RESTRICT;

-- DROP MATERIALIZED VIEW with multiple views, IF EXISTS, and CASCADE
DROP MATERIALIZED VIEW IF EXISTS view1, view2, view3 CASCADE;
DROP MATERIALIZED VIEW IF EXISTS sales_mv, orders_mv, customers_mv CASCADE;
DROP MATERIALIZED VIEW IF EXISTS analytics.sales_summary, analytics.customer_metrics, analytics.product_performance CASCADE;

-- DROP MATERIALIZED VIEW with multiple views, IF EXISTS, and RESTRICT
DROP MATERIALIZED VIEW IF EXISTS view1, view2, view3 RESTRICT;
DROP MATERIALIZED VIEW IF EXISTS sales_mv, orders_mv, customers_mv RESTRICT;
DROP MATERIALIZED VIEW IF EXISTS analytics.sales_summary, analytics.customer_metrics, analytics.product_performance RESTRICT;

-- DROP MATERIALIZED VIEW with quoted identifiers
DROP MATERIALIZED VIEW "SalesAnalytics"."MonthlySummary";
DROP MATERIALIZED VIEW "my-schema"."view-with-hyphens";
DROP MATERIALIZED VIEW "Schema With Spaces"."View With Spaces";
DROP MATERIALIZED VIEW "UPPERCASE_SCHEMA"."UPPERCASE_VIEW";

-- DROP MATERIALIZED VIEW with mixed quoted and unquoted identifiers
DROP MATERIALIZED VIEW analytics."Customer Metrics";
DROP MATERIALIZED VIEW "reporting-schema".sales_view;
DROP MATERIALIZED VIEW "My Analytics"."Daily Reports";

-- DROP MATERIALIZED VIEW with IF EXISTS and quoted identifiers
DROP MATERIALIZED VIEW IF EXISTS "SalesAnalytics"."MonthlySummary";
DROP MATERIALIZED VIEW IF EXISTS "my-schema"."view-with-hyphens";
DROP MATERIALIZED VIEW IF EXISTS "Schema With Spaces"."View With Spaces";

-- DROP MATERIALIZED VIEW with quoted identifiers and CASCADE/RESTRICT
DROP MATERIALIZED VIEW "SalesAnalytics"."MonthlySummary" CASCADE;
DROP MATERIALIZED VIEW "my-schema"."view-with-hyphens" RESTRICT;
DROP MATERIALIZED VIEW IF EXISTS "Schema With Spaces"."View With Spaces" CASCADE;

-- DROP MATERIALIZED VIEW with business domain naming patterns
DROP MATERIALIZED VIEW sales_analytics.regional_performance_mv;
DROP MATERIALIZED VIEW customer_insights.behavior_analysis_mv;
DROP MATERIALIZED VIEW financial_data.revenue_tracking_mv;
DROP MATERIALIZED VIEW operations.supply_chain_mv;
DROP MATERIALIZED VIEW marketing.campaign_effectiveness_mv;
DROP MATERIALIZED VIEW hr_analytics.employee_performance_mv;

-- DROP MATERIALIZED VIEW with time-based naming patterns
DROP MATERIALIZED VIEW reporting.daily_sales_summary;
DROP MATERIALIZED VIEW analytics.weekly_customer_metrics;
DROP MATERIALIZED VIEW warehouse.monthly_inventory_snapshot;
DROP MATERIALIZED VIEW finance.quarterly_revenue_analysis;
DROP MATERIALIZED VIEW operations.yearly_performance_review;

-- DROP MATERIALIZED VIEW with aggregation naming patterns
DROP MATERIALIZED VIEW analytics.customer_order_totals_mv;
DROP MATERIALIZED VIEW reporting.product_sales_aggregates_mv;
DROP MATERIALIZED VIEW warehouse.inventory_level_summaries_mv;
DROP MATERIALIZED VIEW finance.revenue_by_region_mv;
DROP MATERIALIZED VIEW marketing.campaign_conversion_rates_mv;

-- DROP MATERIALIZED VIEW with common business scenarios
DROP MATERIALIZED VIEW ecommerce.customer_lifetime_value_mv;
DROP MATERIALIZED VIEW retail.product_performance_summary_mv;
DROP MATERIALIZED VIEW logistics.shipping_efficiency_metrics_mv;
DROP MATERIALIZED VIEW banking.transaction_risk_analysis_mv;
DROP MATERIALIZED VIEW healthcare.patient_outcome_metrics_mv;
DROP MATERIALIZED VIEW manufacturing.production_quality_summary_mv;

-- DROP MATERIALIZED VIEW with IF EXISTS for business scenarios
DROP MATERIALIZED VIEW IF EXISTS ecommerce.customer_lifetime_value_mv;
DROP MATERIALIZED VIEW IF EXISTS retail.product_performance_summary_mv;
DROP MATERIALIZED VIEW IF EXISTS logistics.shipping_efficiency_metrics_mv;
DROP MATERIALIZED VIEW IF EXISTS banking.transaction_risk_analysis_mv;
DROP MATERIALIZED VIEW IF EXISTS healthcare.patient_outcome_metrics_mv;
DROP MATERIALIZED VIEW IF EXISTS manufacturing.production_quality_summary_mv;

-- DROP MATERIALIZED VIEW with CASCADE for business scenarios
DROP MATERIALIZED VIEW ecommerce.customer_lifetime_value_mv CASCADE;
DROP MATERIALIZED VIEW retail.product_performance_summary_mv CASCADE;
DROP MATERIALIZED VIEW logistics.shipping_efficiency_metrics_mv CASCADE;
DROP MATERIALIZED VIEW banking.transaction_risk_analysis_mv CASCADE;

-- DROP MATERIALIZED VIEW with schema names that might be keywords
DROP MATERIALIZED VIEW data.sales_summary_mv;
DROP MATERIALIZED VIEW analytics.customer_data_mv;
DROP MATERIALIZED VIEW reporting.warehouse_metrics_mv;
DROP MATERIALIZED VIEW schema.view_analytics_mv;
DROP MATERIALIZED VIEW warehouse.summary_data_mv;

-- DROP MATERIALIZED VIEW with view names that might be keywords
DROP MATERIALIZED VIEW sales_schema.data;
DROP MATERIALIZED VIEW analytics_schema.view;
DROP MATERIALIZED VIEW reporting_schema.schema;
DROP MATERIALIZED VIEW warehouse_schema.analytics;
DROP MATERIALIZED VIEW finance_schema.reporting;
DROP MATERIALIZED VIEW operations_schema.warehouse;
DROP MATERIALIZED VIEW marketing_schema.summary;
DROP MATERIALIZED VIEW hr_schema.metrics;

-- DROP MATERIALIZED VIEW statements without semicolons (should also be valid)
DROP MATERIALIZED VIEW test_schema.test_mv
DROP MATERIALIZED VIEW IF EXISTS test_schema.conditional_mv
DROP MATERIALIZED VIEW simple_test_mv CASCADE
DROP MATERIALIZED VIEW IF EXISTS another_test_mv RESTRICT

-- DROP MATERIALIZED VIEW with complex multiple view combinations
DROP MATERIALIZED VIEW
    sales.daily_summary,
    sales.weekly_summary,
    sales.monthly_summary,
    sales.quarterly_summary CASCADE;

DROP MATERIALIZED VIEW IF EXISTS
    analytics.customer_segments,
    analytics.product_categories,
    analytics.regional_analysis RESTRICT;

-- DROP MATERIALIZED VIEW with underscored naming conventions
DROP MATERIALIZED VIEW customer_analytics.lifetime_value_mv;
DROP MATERIALIZED VIEW sales_reporting.monthly_summary_mv;
DROP MATERIALIZED VIEW inventory_management.stock_alerts_mv;
DROP MATERIALIZED VIEW financial_reporting.quarterly_results_mv;
DROP MATERIALIZED VIEW operational_metrics.daily_kpi_mv;
DROP MATERIALIZED VIEW marketing_analytics.campaign_roi_mv;

-- DROP MATERIALIZED VIEW with mixed case naming
DROP MATERIALIZED VIEW DataLake.CustomerAnalyticsMV;
DROP MATERIALIZED VIEW SalesSchema.ProductPerformanceMV;
DROP MATERIALIZED VIEW AnalyticsWarehouse.CustomerSegmentationMV;
DROP MATERIALIZED VIEW ReportingHub.FinancialSummaryMV;

-- DROP MATERIALIZED VIEW with special characters in quoted names
DROP MATERIALIZED VIEW "view-with-hyphens";
DROP MATERIALIZED VIEW "view.with.dots";
DROP MATERIALIZED VIEW "view$with$dollars";
DROP MATERIALIZED VIEW "view@with@symbols";
DROP MATERIALIZED VIEW "view_with_underscores";
DROP MATERIALIZED VIEW "view123with456numbers";

-- DROP MATERIALIZED VIEW with comprehensive quoted name combinations
DROP MATERIALIZED VIEW "my-analytics"."sales-summary" CASCADE;
DROP MATERIALIZED VIEW "data.warehouse"."customer.metrics" RESTRICT;
DROP MATERIALIZED VIEW IF EXISTS "schema-name"."view-name" CASCADE;
DROP MATERIALIZED VIEW IF EXISTS "UPPER_CASE_SCHEMA"."UPPER_CASE_VIEW" RESTRICT;

-- DROP MATERIALIZED VIEW examples that correspond to typical CREATE MATERIALIZED VIEW scenarios
-- These views could have been created with various aggregations and joins
DROP MATERIALIZED VIEW analytics.customer_order_summary_mv;
DROP MATERIALIZED VIEW reporting.sales_performance_by_region_mv;
DROP MATERIALIZED VIEW warehouse.inventory_turnover_analysis_mv;
DROP MATERIALIZED VIEW finance.monthly_revenue_breakdown_mv;
DROP MATERIALIZED VIEW operations.supplier_performance_metrics_mv;
DROP MATERIALIZED VIEW marketing.customer_acquisition_cost_mv;
DROP MATERIALIZED VIEW hr.employee_productivity_summary_mv;
DROP MATERIALIZED VIEW logistics.delivery_performance_tracking_mv;
DROP MATERIALIZED VIEW retail.product_recommendation_matrix_mv;
DROP MATERIALIZED VIEW banking.fraud_detection_summary_mv;

-- DROP MATERIALIZED VIEW with IF EXISTS for typical scenarios
DROP MATERIALIZED VIEW IF EXISTS analytics.customer_order_summary_mv;
DROP MATERIALIZED VIEW IF EXISTS reporting.sales_performance_by_region_mv;
DROP MATERIALIZED VIEW IF EXISTS warehouse.inventory_turnover_analysis_mv;
DROP MATERIALIZED VIEW IF EXISTS finance.monthly_revenue_breakdown_mv;
DROP MATERIALIZED VIEW IF EXISTS operations.supplier_performance_metrics_mv;

-- DROP MATERIALIZED VIEW with CASCADE for cleanup scenarios
DROP MATERIALIZED VIEW analytics.customer_order_summary_mv CASCADE;
DROP MATERIALIZED VIEW reporting.sales_performance_by_region_mv CASCADE;
DROP MATERIALIZED VIEW warehouse.inventory_turnover_analysis_mv CASCADE;
DROP MATERIALIZED VIEW finance.monthly_revenue_breakdown_mv CASCADE;

-- DROP MATERIALIZED VIEW for data warehouse cleanup scenarios
DROP MATERIALIZED VIEW dwh.fact_sales_summary_mv;
DROP MATERIALIZED VIEW dwh.dim_customer_attributes_mv;
DROP MATERIALIZED VIEW dwh.agg_product_performance_mv;
DROP MATERIALIZED VIEW dwh.mart_financial_kpis_mv;
DROP MATERIALIZED VIEW dwh.staging_order_metrics_mv;

-- DROP MATERIALIZED VIEW with comprehensive combinations
DROP MATERIALIZED VIEW IF EXISTS dwh.fact_sales_summary_mv, dwh.dim_customer_attributes_mv CASCADE;
DROP MATERIALIZED VIEW IF EXISTS analytics.customer_segments, reporting.monthly_metrics RESTRICT;