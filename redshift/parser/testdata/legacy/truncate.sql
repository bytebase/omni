-- TRUNCATE command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_TRUNCATE.html

-- Basic TRUNCATE statements
TRUNCATE category;
truncate category;

-- TRUNCATE with TABLE keyword (optional)
TRUNCATE TABLE category;
truncate table category;

-- TRUNCATE materialized view
TRUNCATE my_materialized_view;
truncate my_materialized_view;

-- Examples from AWS documentation
-- Example 1: Truncate a table (CATEGORY table is empty after this command)
TRUNCATE category;

-- Example 2: Truncate a materialized view
TRUNCATE my_materialized_view;

-- Example 3: Demonstrating automatic commit behavior
-- The following transaction is automatically committed, and the rollback has no effect
BEGIN;
TRUNCATE date;
ROLLBACK;
-- DATE table remains empty after ROLLBACK

-- Mixed case variations
Truncate category;
TRUNCATE CATEGORY;
TrUnCaTe TaBlE category;
Truncate Table Category;

-- Quoted identifiers
TRUNCATE "category";
TRUNCATE TABLE "category";
TRUNCATE "my_table";
truncate table "my_table";

-- Tables with schema qualifiers
TRUNCATE public.category;
TRUNCATE TABLE public.category;
truncate sales.fact_sales;
TRUNCATE TABLE marketing.campaign_data;

-- Quoted schema and table names
TRUNCATE "public"."category";
TRUNCATE TABLE "public"."category";
truncate "sales"."fact_sales";

-- Mixed quoted and unquoted identifiers
TRUNCATE public."category";
TRUNCATE "public".category;
truncate table sales."fact_sales";

-- Special characters in quoted names
TRUNCATE "table-with-hyphens";
TRUNCATE TABLE "table.with.dots";
truncate "table with spaces";
TRUNCATE TABLE "CaseSensitiveTable";

-- Materialized views with various naming patterns
TRUNCATE mv_sales_summary;
TRUNCATE materialized_view_1;
truncate "mv_complex-name";
TRUNCATE public.mv_aggregated_data;
TRUNCATE "schema"."mv_quoted";

-- Real-world usage examples
-- Example: Truncating staging tables
TRUNCATE staging.raw_events;
TRUNCATE TABLE staging.temp_users;

-- Example: Truncating temporary tables
TRUNCATE temp_analysis;
TRUNCATE TABLE temp_calculations;

-- Example: Truncating fact tables for reload
TRUNCATE fact_sales;
TRUNCATE TABLE fact_inventory;

-- Example: Truncating dimension tables
TRUNCATE dim_product;
TRUNCATE TABLE dim_customer;

-- Example: Truncating error or log tables
TRUNCATE etl_errors;
TRUNCATE TABLE process_log;

-- Example: Transaction blocks showing automatic commit
BEGIN;
TRUNCATE staging_table;
-- Any subsequent commands in this transaction will see empty staging_table
SELECT COUNT(*) FROM staging_table; -- Returns 0
COMMIT; -- This COMMIT is effectively a no-op since TRUNCATE already committed

BEGIN TRANSACTION;
TRUNCATE audit_log;
-- Even if we try to rollback, the TRUNCATE has already committed
ROLLBACK; -- This has no effect on the TRUNCATE operation

-- Example: Using TRUNCATE in ETL processes
-- Step 1: Truncate staging
TRUNCATE staging.customer_data;

-- Step 2: Truncate target before reload
TRUNCATE fact_customer_transactions;

-- Example: Maintenance operations
TRUNCATE old_backup_table;
TRUNCATE TABLE archived_data_2020;