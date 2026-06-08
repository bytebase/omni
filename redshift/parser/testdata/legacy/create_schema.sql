-- CREATE SCHEMA test cases based on AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_SCHEMA.html

-- Basic schema creation (from AWS documentation)
CREATE SCHEMA us_sales;

-- Schema with authorization (from AWS documentation)
CREATE SCHEMA us_sales AUTHORIZATION dwuser;

-- Schema with quota (from AWS documentation)
CREATE SCHEMA us_sales AUTHORIZATION dwuser QUOTA 50 GB;

-- Schema with IF NOT EXISTS (from AWS documentation)
CREATE SCHEMA IF NOT EXISTS us_sales;

-- Basic schema variations
CREATE SCHEMA sales;
CREATE SCHEMA marketing;
CREATE SCHEMA finance;
CREATE SCHEMA hr;
CREATE SCHEMA analytics;
CREATE SCHEMA reporting;
CREATE SCHEMA data_warehouse;
CREATE SCHEMA staging;
CREATE SCHEMA temp;
CREATE SCHEMA archive;

-- Schema names with underscores
CREATE SCHEMA customer_data;
CREATE SCHEMA product_catalog;
CREATE SCHEMA order_management;
CREATE SCHEMA user_analytics;
CREATE SCHEMA financial_reporting;
CREATE SCHEMA marketing_campaigns;
CREATE SCHEMA supply_chain;
CREATE SCHEMA inventory_management;
CREATE SCHEMA compliance_audit;
CREATE SCHEMA data_science;

-- Quoted schema names for special characters
CREATE SCHEMA "MySchema";
CREATE SCHEMA "schema-with-dash";
CREATE SCHEMA "schema.with.dots";
CREATE SCHEMA "schema with spaces";
CREATE SCHEMA "UPPERCASE_SCHEMA";
CREATE SCHEMA "Schema_With_Mixed_Case";
CREATE SCHEMA "123schema";
CREATE SCHEMA "schema123";
CREATE SCHEMA "schema_v2";

-- Single character schema names
CREATE SCHEMA a;
CREATE SCHEMA "A";
CREATE SCHEMA z;

-- IF NOT EXISTS variations
CREATE SCHEMA IF NOT EXISTS customer_data;
CREATE SCHEMA IF NOT EXISTS product_catalog;
CREATE SCHEMA IF NOT EXISTS "Special Schema";
CREATE SCHEMA IF NOT EXISTS analytics AUTHORIZATION admin_user;
CREATE SCHEMA IF NOT EXISTS reporting AUTHORIZATION report_user QUOTA 100 GB;
CREATE SCHEMA IF NOT EXISTS temp_schema QUOTA UNLIMITED;

-- Authorization variations
CREATE SCHEMA sales_dept AUTHORIZATION sales_manager;
CREATE SCHEMA marketing_dept AUTHORIZATION marketing_lead;
CREATE SCHEMA finance_dept AUTHORIZATION finance_admin;
CREATE SCHEMA analytics_team AUTHORIZATION data_scientist;
CREATE SCHEMA reporting_service AUTHORIZATION service_account;
CREATE SCHEMA external_data AUTHORIZATION external_user;
CREATE SCHEMA vendor_access AUTHORIZATION vendor_account;
CREATE SCHEMA contractor_space AUTHORIZATION temp_contractor;
CREATE SCHEMA audit_logs AUTHORIZATION compliance_officer;
CREATE SCHEMA ml_workspace AUTHORIZATION ml_engineer;

-- Quoted authorization users
CREATE SCHEMA special_project AUTHORIZATION "User With Spaces";
CREATE SCHEMA legacy_system AUTHORIZATION "legacy-service-account";
CREATE SCHEMA cross_region AUTHORIZATION "cross.region.user";
CREATE SCHEMA federated_access AUTHORIZATION "FEDERATED_USER";
CREATE SCHEMA service_integration AUTHORIZATION "service_account_123";

-- Quota variations with MB
CREATE SCHEMA small_schema QUOTA 100 MB;
CREATE SCHEMA test_schema QUOTA 500 MB;
CREATE SCHEMA temp_data QUOTA 1024 MB;
CREATE SCHEMA cache_schema QUOTA 2048 MB;
CREATE SCHEMA sample_data QUOTA 50 MB;

-- Quota variations with GB
CREATE SCHEMA medium_schema QUOTA 1 GB;
CREATE SCHEMA user_workspace QUOTA 5 GB;
CREATE SCHEMA project_data QUOTA 10 GB;
CREATE SCHEMA department_analytics QUOTA 25 GB;
CREATE SCHEMA staging_area QUOTA 100 GB;
CREATE SCHEMA historical_data QUOTA 250 GB;
CREATE SCHEMA backup_schema QUOTA 500 GB;
CREATE SCHEMA data_lake QUOTA 1000 GB;

-- Quota variations with TB
CREATE SCHEMA enterprise_data QUOTA 1 TB;
CREATE SCHEMA big_analytics QUOTA 5 TB;
CREATE SCHEMA data_warehouse_main QUOTA 10 TB;
CREATE SCHEMA historical_archive QUOTA 25 TB;
CREATE SCHEMA raw_data_lake QUOTA 50 TB;
CREATE SCHEMA enterprise_backup QUOTA 100 TB;

-- Decimal quota values
CREATE SCHEMA fractional_quota_mb QUOTA 100.5 MB;
CREATE SCHEMA fractional_quota_gb QUOTA 1.5 GB;
CREATE SCHEMA fractional_quota_tb QUOTA 2.25 TB;
CREATE SCHEMA precise_allocation QUOTA 10.75 GB;
CREATE SCHEMA micro_allocation QUOTA 0.5 GB;

-- UNLIMITED quota
CREATE SCHEMA unlimited_schema QUOTA UNLIMITED;
CREATE SCHEMA no_limit_data QUOTA UNLIMITED;
CREATE SCHEMA enterprise_unlimited QUOTA UNLIMITED;

-- Full syntax combinations
CREATE SCHEMA IF NOT EXISTS sales_analytics AUTHORIZATION sales_analyst QUOTA 50 GB;
CREATE SCHEMA IF NOT EXISTS marketing_data AUTHORIZATION marketing_team QUOTA 25 GB;
CREATE SCHEMA IF NOT EXISTS finance_reports AUTHORIZATION finance_admin QUOTA 100 GB;
CREATE SCHEMA IF NOT EXISTS hr_analytics AUTHORIZATION hr_manager QUOTA 10 GB;
CREATE SCHEMA IF NOT EXISTS customer_insights AUTHORIZATION data_scientist QUOTA 200 GB;
CREATE SCHEMA IF NOT EXISTS product_metrics AUTHORIZATION product_manager QUOTA 75 GB;
CREATE SCHEMA IF NOT EXISTS operational_data AUTHORIZATION ops_team QUOTA 500 GB;
CREATE SCHEMA IF NOT EXISTS compliance_audit AUTHORIZATION compliance_officer QUOTA 50 GB;
CREATE SCHEMA IF NOT EXISTS vendor_data AUTHORIZATION vendor_manager QUOTA 25 GB;
CREATE SCHEMA IF NOT EXISTS research_lab AUTHORIZATION research_lead QUOTA UNLIMITED;

-- Complex quoted combinations
CREATE SCHEMA IF NOT EXISTS "Complex Schema Name" AUTHORIZATION "Complex User Name" QUOTA 100 GB;
CREATE SCHEMA IF NOT EXISTS "schema-with-special.chars" AUTHORIZATION "user-with-special.chars" QUOTA 50 GB;
CREATE SCHEMA IF NOT EXISTS "2024_Q4_Analysis" AUTHORIZATION "quarterly.analyst" QUOTA 200 GB;
CREATE SCHEMA IF NOT EXISTS "Multi Word Schema" AUTHORIZATION "Multi Word User" QUOTA 150 GB;

-- Real-world departmental schemas
CREATE SCHEMA customer_360 AUTHORIZATION customer_analytics_team QUOTA 500 GB;
CREATE SCHEMA fraud_detection AUTHORIZATION security_team QUOTA 100 GB;
CREATE SCHEMA recommendation_engine AUTHORIZATION ml_team QUOTA 250 GB;
CREATE SCHEMA supply_chain_optimization AUTHORIZATION operations_team QUOTA 300 GB;
CREATE SCHEMA financial_risk_modeling AUTHORIZATION risk_management QUOTA 200 GB;
CREATE SCHEMA marketing_attribution AUTHORIZATION marketing_analytics QUOTA 150 GB;
CREATE SCHEMA product_performance AUTHORIZATION product_analytics QUOTA 100 GB;
CREATE SCHEMA user_behavior_analysis AUTHORIZATION ux_research_team QUOTA 75 GB;
CREATE SCHEMA competitive_intelligence AUTHORIZATION strategy_team QUOTA 50 GB;
CREATE SCHEMA regulatory_reporting AUTHORIZATION compliance_team QUOTA 25 GB;

-- Geographic/regional schemas
CREATE SCHEMA north_america_sales AUTHORIZATION na_sales_manager QUOTA 500 GB;
CREATE SCHEMA europe_analytics AUTHORIZATION eu_analyst QUOTA 300 GB;
CREATE SCHEMA asia_pacific_data AUTHORIZATION apac_team QUOTA 400 GB;
CREATE SCHEMA latin_america_metrics AUTHORIZATION latam_manager QUOTA 200 GB;
CREATE SCHEMA global_consolidated AUTHORIZATION global_admin QUOTA 1 TB;

-- Time-based schemas
CREATE SCHEMA daily_operations AUTHORIZATION daily_ops QUOTA 100 GB;
CREATE SCHEMA weekly_reports AUTHORIZATION weekly_analyst QUOTA 50 GB;
CREATE SCHEMA monthly_analytics AUTHORIZATION monthly_reporter QUOTA 200 GB;
CREATE SCHEMA quarterly_insights AUTHORIZATION quarterly_analyst QUOTA 500 GB;
CREATE SCHEMA annual_archive AUTHORIZATION archive_manager QUOTA 2 TB;

-- Project-specific schemas
CREATE SCHEMA project_alpha AUTHORIZATION project_lead_alpha QUOTA 100 GB;
CREATE SCHEMA project_beta AUTHORIZATION project_lead_beta QUOTA 150 GB;
CREATE SCHEMA project_gamma AUTHORIZATION project_lead_gamma QUOTA 75 GB;
CREATE SCHEMA migration_workspace AUTHORIZATION migration_team QUOTA 500 GB;
CREATE SCHEMA integration_testing AUTHORIZATION qa_team QUOTA 200 GB;

-- Service and application schemas
CREATE SCHEMA web_analytics AUTHORIZATION web_team QUOTA 300 GB;
CREATE SCHEMA mobile_analytics AUTHORIZATION mobile_team QUOTA 200 GB;
CREATE SCHEMA api_metrics AUTHORIZATION api_team QUOTA 100 GB;
CREATE SCHEMA etl_staging AUTHORIZATION etl_service QUOTA 500 GB;
CREATE SCHEMA data_quality AUTHORIZATION dq_team QUOTA 50 GB;

-- Development lifecycle schemas
CREATE SCHEMA dev_sandbox AUTHORIZATION developer QUOTA 10 GB;
CREATE SCHEMA test_environment AUTHORIZATION test_team QUOTA 50 GB;
CREATE SCHEMA staging_validation AUTHORIZATION staging_team QUOTA 100 GB;
CREATE SCHEMA production_replica AUTHORIZATION prod_admin QUOTA 1 TB;
CREATE SCHEMA backup_recovery AUTHORIZATION backup_service QUOTA 2 TB;

-- Edge cases and boundary testing
CREATE SCHEMA a AUTHORIZATION b QUOTA 1 MB;  -- Minimal schema
CREATE SCHEMA longest_possible_schema_name_within_limits AUTHORIZATION longest_possible_user_name_within_limits QUOTA 999999 TB;

-- Case sensitivity tests
CREATE SCHEMA lowercase_schema AUTHORIZATION lowercase_user QUOTA 10 gb;
CREATE SCHEMA UPPERCASE_SCHEMA AUTHORIZATION UPPERCASE_USER QUOTA 10 GB;
CREATE SCHEMA "CamelCaseSchema" AUTHORIZATION "CamelCaseUser" QUOTA 10 GB;

-- SQL keywords as quoted identifiers (edge case testing)
CREATE SCHEMA "SELECT" AUTHORIZATION "CREATE" QUOTA 10 GB;
CREATE SCHEMA "TABLE" AUTHORIZATION "USER" QUOTA 5 GB;
CREATE SCHEMA "SCHEMA" AUTHORIZATION "AUTHORIZATION" QUOTA 1 GB;
CREATE SCHEMA "QUOTA" AUTHORIZATION "UNLIMITED" QUOTA UNLIMITED;