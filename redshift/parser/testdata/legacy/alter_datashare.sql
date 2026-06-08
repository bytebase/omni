-- ALTER DATASHARE test cases
-- Based on: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_DATASHARE.html

-- Basic ALTER DATASHARE: Add a schema
ALTER DATASHARE salesshare ADD SCHEMA public;

-- Add a specific table
ALTER DATASHARE salesshare ADD TABLE public.tickit_sales_redshift;

-- Add multiple tables
ALTER DATASHARE salesshare ADD TABLE public.sales, public.customers, marketing.campaigns;

-- Add table without schema qualification (uses current schema)
ALTER DATASHARE salesshare ADD TABLE events;

-- Remove a table
ALTER DATASHARE salesshare REMOVE TABLE public.tickit_sales_redshift;

-- Remove multiple tables
ALTER DATASHARE salesshare REMOVE TABLE public.sales, public.customers;

-- Add all tables in a schema
ALTER DATASHARE salesshare ADD ALL TABLES IN SCHEMA public;

-- Add all tables in multiple schemas
ALTER DATASHARE salesshare ADD ALL TABLES IN SCHEMA public, marketing, sales;

-- Add a function with arguments
ALTER DATASHARE salesshare ADD FUNCTION public.compute_total(integer, decimal);

-- Add multiple functions
ALTER DATASHARE salesshare ADD FUNCTION public.get_user_name(integer), public.calculate_tax(decimal, decimal);

-- Add function without arguments
ALTER DATASHARE salesshare ADD FUNCTION public.get_current_timestamp();

-- Remove a function
ALTER DATASHARE salesshare REMOVE FUNCTION public.compute_total(integer, decimal);

-- Add all functions in a schema
ALTER DATASHARE salesshare ADD ALL FUNCTIONS IN SCHEMA public;

-- Add all functions in multiple schemas
ALTER DATASHARE salesshare ADD ALL FUNCTIONS IN SCHEMA public, utility, analytics;

-- Remove all functions in a schema
ALTER DATASHARE salesshare REMOVE ALL FUNCTIONS IN SCHEMA public;

-- Set PUBLICACCESSIBLE to TRUE
ALTER DATASHARE salesshare SET PUBLICACCESSIBLE TRUE;

-- Set PUBLICACCESSIBLE to FALSE
ALTER DATASHARE salesshare SET PUBLICACCESSIBLE FALSE;

-- Set PUBLICACCESSIBLE with equals sign
ALTER DATASHARE salesshare SET PUBLICACCESSIBLE = TRUE;

-- Set INCLUDENEW for a schema
ALTER DATASHARE salesshare SET INCLUDENEW TRUE FOR SCHEMA public;

-- Set INCLUDENEW to FALSE
ALTER DATASHARE salesshare SET INCLUDENEW FALSE FOR SCHEMA marketing;

-- Set INCLUDENEW with equals sign
ALTER DATASHARE salesshare SET INCLUDENEW = TRUE FOR SCHEMA sales;

-- Complex example: Add multiple schemas
ALTER DATASHARE analytics_share ADD SCHEMA finance, operations, hr;

-- Complex example: Remove multiple object types (multiple statements needed)
ALTER DATASHARE analytics_share REMOVE SCHEMA finance;
ALTER DATASHARE analytics_share REMOVE TABLE operations.inventory;
ALTER DATASHARE analytics_share REMOVE FUNCTION hr.calculate_bonus(integer, decimal);

-- Add tables with mixed schema qualification
ALTER DATASHARE myshare ADD TABLE schema1.table1, table2, schema2.table3;

-- Function with multiple argument types
ALTER DATASHARE funcshare ADD FUNCTION analytics.process_data(varchar, integer, date, boolean);

-- Function with array types
ALTER DATASHARE funcshare ADD FUNCTION public.aggregate_values(integer[]);

-- Multiple configuration changes (requires separate statements)
ALTER DATASHARE configshare SET PUBLICACCESSIBLE TRUE;
ALTER DATASHARE configshare SET INCLUDENEW TRUE FOR SCHEMA public;

-- Edge cases with special names
ALTER DATASHARE "My-DataShare" ADD SCHEMA "schema-with-dashes";
ALTER DATASHARE my_share ADD TABLE "Table With Spaces";
ALTER DATASHARE share123 ADD FUNCTION "schema"."function-name"(text);

-- Remove all tables from multiple schemas
ALTER DATASHARE cleanup_share REMOVE ALL TABLES IN SCHEMA temp, staging, archive;