-- DROP DATASHARE command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_DATASHARE.html

-- Basic DROP DATASHARE syntax
DROP DATASHARE salesshare;

-- DROP DATASHARE with IF EXISTS clause
DROP DATASHARE IF EXISTS salesshare;

-- DROP DATASHARE with various datashare names
DROP DATASHARE myshare;
DROP DATASHARE customer_data;
DROP DATASHARE product_catalog;
DROP DATASHARE analytics_share;

-- DROP DATASHARE with quoted identifiers
DROP DATASHARE "SalesShare";
DROP DATASHARE "my-share";
DROP DATASHARE "Share With Spaces";

-- DROP DATASHARE with IF EXISTS and quoted identifiers
DROP DATASHARE IF EXISTS "SalesShare";
DROP DATASHARE IF EXISTS "my-share";
DROP DATASHARE IF EXISTS "Share With Spaces";

-- DROP DATASHARE with common naming patterns
DROP DATASHARE sales_2023;
DROP DATASHARE customer_360_view;
DROP DATASHARE financial_reporting;
DROP DATASHARE hr_analytics;

-- DROP DATASHARE with IF EXISTS for common patterns
DROP DATASHARE IF EXISTS sales_2023;
DROP DATASHARE IF EXISTS customer_360_view;
DROP DATASHARE IF EXISTS financial_reporting;
DROP DATASHARE IF EXISTS hr_analytics;

-- DROP DATASHARE with mixed case names
DROP DATASHARE SalesData;
DROP DATASHARE CustomerInfo;
DROP DATASHARE ProductCatalog;

-- DROP DATASHARE with underscored names
DROP DATASHARE sales_data_share;
DROP DATASHARE customer_info_share;
DROP DATASHARE product_catalog_share;

-- DROP DATASHARE with IF EXISTS for mixed case and underscored names
DROP DATASHARE IF EXISTS SalesData;
DROP DATASHARE IF EXISTS sales_data_share;

-- DROP DATASHARE examples from AWS documentation
-- Note: In practice, AWS Data Exchange datashares require special session variable:
-- SET datashare_break_glass_session_var to '620c871f890c49';
DROP DATASHARE salesshare;

-- Additional realistic datashare names
DROP DATASHARE market_data;
DROP DATASHARE financial_metrics;
DROP DATASHARE user_behavior;
DROP DATASHARE inventory_tracking;
DROP DATASHARE compliance_reports;

-- With IF EXISTS for safety
DROP DATASHARE IF EXISTS market_data;
DROP DATASHARE IF EXISTS financial_metrics;
DROP DATASHARE IF EXISTS user_behavior;
DROP DATASHARE IF EXISTS inventory_tracking;
DROP DATASHARE IF EXISTS compliance_reports;

-- Edge cases with special characters in quoted names
DROP DATASHARE "share-with-hyphens";
DROP DATASHARE "share.with.dots";
DROP DATASHARE "share$with$dollars";

-- DROP DATASHARE statements without semicolons (should also be valid)
DROP DATASHARE test_share
DROP DATASHARE IF EXISTS test_share