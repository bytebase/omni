-- REVOKE test cases based on AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_REVOKE.html

-- Basic table privileges
REVOKE SELECT ON TABLE sales FROM dbuser;
REVOKE INSERT, UPDATE, DELETE ON TABLE sales FROM dbuser;
REVOKE ALL ON TABLE sales FROM dbuser;
REVOKE ALL PRIVILEGES ON TABLE sales FROM dbuser;
REVOKE ALTER, DROP, TRUNCATE, REFERENCES ON TABLE sales FROM admin_user;

-- Table privileges from multiple users
REVOKE SELECT ON TABLE sales FROM user1, user2, user3;
REVOKE INSERT ON sales FROM dbuser1, dbuser2;

-- Revoke GRANT OPTION only (users keep the privilege but can't grant it)
REVOKE GRANT OPTION FOR SELECT ON TABLE sales FROM dbuser;
REVOKE GRANT OPTION FOR ALL ON TABLE sales FROM dbuser;

-- Column-level privileges
REVOKE SELECT (col1) ON TABLE t1 FROM PUBLIC;
REVOKE UPDATE (col1, col2) ON TABLE t1 FROM user1;
REVOKE SELECT (customer_id, customer_name) ON TABLE customers FROM analyst;
REVOKE ALL (order_id, order_date, order_total) ON orders FROM report_user;

-- Revoke from PUBLIC
REVOKE SELECT ON TABLE public_data FROM PUBLIC;
REVOKE USAGE ON SCHEMA public FROM PUBLIC;

-- Revoke from ROLE and GROUP
REVOKE SELECT ON TABLE sales FROM ROLE analyst_role;
REVOKE INSERT ON TABLE logs FROM GROUP etl_group;
REVOKE ALL ON TABLE temp_data FROM ROLE temp_role, GROUP admin_group;

-- All tables in schema
REVOKE SELECT ON ALL TABLES IN SCHEMA myschema FROM dbuser;
REVOKE ALL ON ALL TABLES IN SCHEMA reporting FROM ROLE report_role;

-- Database privileges
REVOKE CREATE ON DATABASE mydb FROM dbuser;
REVOKE USAGE ON DATABASE prod_db FROM readonly_user;
REVOKE TEMPORARY ON DATABASE session_db FROM temp_user;
REVOKE TEMP ON DATABASE session_db FROM temp_user;
REVOKE ALTER ON DATABASE mydb FROM admin_user;
REVOKE ALL ON DATABASE mydb FROM dba_role;
REVOKE ALL PRIVILEGES ON DATABASE mydb FROM admin_group;

-- Database privileges with CASCADE/RESTRICT
REVOKE CREATE ON DATABASE mydb FROM dbuser CASCADE;
REVOKE ALL ON DATABASE mydb FROM admin RESTRICT;

-- Schema privileges
REVOKE CREATE ON SCHEMA myschema FROM dbuser;
REVOKE USAGE ON SCHEMA reporting FROM analyst;
REVOKE ALTER ON SCHEMA dynamic_schema FROM schema_admin;
REVOKE DROP ON SCHEMA temp_schema FROM cleanup_role;
REVOKE ALL ON SCHEMA myschema FROM schema_owner;
REVOKE ALL PRIVILEGES ON SCHEMA myschema FROM schema_admin;

-- Schema privileges with CASCADE/RESTRICT
REVOKE USAGE ON SCHEMA reporting FROM analyst CASCADE;
REVOKE ALL ON SCHEMA myschema FROM admin RESTRICT;

-- Function privileges
REVOKE EXECUTE ON FUNCTION f1() FROM dbuser;
REVOKE EXECUTE ON FUNCTION calculate_tax(numeric, numeric) FROM tax_calculator;
REVOKE ALL ON FUNCTION utility_func() FROM utility_user;
REVOKE ALL PRIVILEGES ON FUNCTION math_func(integer) FROM math_user;

-- Function privileges with parameters
REVOKE EXECUTE ON FUNCTION format_currency(numeric, varchar) FROM report_user;
REVOKE EXECUTE ON FUNCTION validate_email(varchar(255)) FROM validation_service;
REVOKE EXECUTE ON FUNCTION calculate_discount(decimal(10,2), integer) FROM pricing_engine;

-- All functions in schema
REVOKE EXECUTE ON ALL FUNCTIONS IN SCHEMA analytics FROM analyst_role;
REVOKE ALL ON ALL FUNCTIONS IN SCHEMA utilities FROM utility_group;

-- Procedure privileges
REVOKE EXECUTE ON PROCEDURE proc1() FROM dbuser;
REVOKE EXECUTE ON PROCEDURE update_inventory(integer, varchar) FROM inventory_user;
REVOKE ALL ON PROCEDURE maintenance_proc() FROM admin_role;

-- All procedures in schema
REVOKE EXECUTE ON ALL PROCEDURES IN SCHEMA maintenance FROM maintenance_role;
REVOKE ALL ON ALL PROCEDURES IN SCHEMA etl FROM etl_group;

-- Language privileges
REVOKE USAGE ON LANGUAGE plpythonu FROM untrusted_user;
REVOKE USAGE ON LANGUAGE plpythonu FROM PUBLIC;

-- Datashare privileges
REVOKE ALTER ON DATASHARE salesshare FROM myuser;
REVOKE SHARE ON DATASHARE marketingshare FROM marketing_admin;
REVOKE USAGE ON DATASHARE salesshare FROM namespace 'a3f3ae8c-14e8-45ba-9eaa-b42b9b7ae635';
REVOKE USAGE ON DATASHARE customerdata FROM ACCOUNT '123456789012';
REVOKE USAGE ON DATASHARE productdata FROM ACCOUNT '987654321098' VIA DATA_CATALOG;

-- Datashare database and schema usage
REVOKE USAGE ON DATABASE shared_db FROM consumer_user;
REVOKE USAGE ON SCHEMA shared_schema FROM consumer_role;

-- Role privileges
REVOKE ROLE sample_role1 FROM reguser;
REVOKE ROLE admin_role FROM user1;
REVOKE ROLE analyst_role FROM ROLE senior_analyst;
REVOKE ADMIN OPTION FOR ROLE admin_role FROM user1;

-- Multiple role revocations
REVOKE ROLE read_role FROM user1, user2, user3;
REVOKE ROLE admin_role FROM ROLE super_admin, ROLE department_head;
REVOKE ROLE junior_analyst, ROLE data_entry FROM new_hire;

-- System permissions from roles
REVOKE CREATE USER FROM ROLE user_admin;
REVOKE DROP USER FROM ROLE user_admin;
REVOKE ALTER USER FROM ROLE limited_admin;
REVOKE CREATE SCHEMA FROM ROLE schema_creator;
REVOKE DROP SCHEMA FROM ROLE schema_admin;
REVOKE ALTER DEFAULT PRIVILEGES FROM ROLE privilege_admin;
REVOKE ACCESS CATALOG FROM ROLE catalog_user;
REVOKE CREATE TABLE FROM ROLE table_creator;
REVOKE DROP TABLE FROM ROLE table_admin;
REVOKE ALTER TABLE FROM ROLE table_modifier;
REVOKE CREATE OR REPLACE FUNCTION FROM ROLE function_developer;
REVOKE CREATE OR REPLACE EXTERNAL FUNCTION FROM ROLE external_developer;
REVOKE DROP FUNCTION FROM ROLE function_admin;
REVOKE CREATE OR REPLACE PROCEDURE FROM ROLE procedure_developer;
REVOKE DROP PROCEDURE FROM ROLE procedure_admin;
REVOKE CREATE OR REPLACE VIEW FROM ROLE view_developer;
REVOKE DROP VIEW FROM ROLE view_admin;
REVOKE CREATE MODEL FROM ROLE ml_developer;
REVOKE DROP MODEL FROM ROLE ml_admin;
REVOKE CREATE DATASHARE FROM ROLE datashare_creator;
REVOKE ALTER DATASHARE FROM ROLE datashare_admin;
REVOKE DROP DATASHARE FROM ROLE datashare_owner;
REVOKE CREATE LIBRARY FROM ROLE library_manager;
REVOKE DROP LIBRARY FROM ROLE library_admin;
REVOKE CREATE ROLE FROM ROLE role_admin;
REVOKE DROP ROLE FROM ROLE role_manager;
REVOKE TRUNCATE TABLE FROM ROLE data_manager;
REVOKE VACUUM FROM ROLE maintenance_user;
REVOKE ANALYZE FROM ROLE maintenance_user;
REVOKE CANCEL FROM ROLE operations_user;
REVOKE IGNORE RLS FROM ROLE bypass_user;
REVOKE EXPLAIN RLS FROM ROLE security_analyst;
REVOKE EXPLAIN MASKING FROM ROLE compliance_officer;
REVOKE ALL FROM ROLE super_admin;
REVOKE ALL PRIVILEGES FROM ROLE restricted_user;

-- Model privileges (Amazon Redshift ML)
REVOKE EXECUTE ON MODEL customer_churn_model FROM ml_user;
REVOKE ALL ON MODEL sales_forecast_model FROM ml_team;
REVOKE ALL PRIVILEGES ON MODEL recommendation_model FROM data_scientist;
REVOKE CREATE MODEL FROM untrusted_user;

-- Row-level security permissions
REVOKE EXPLAIN RLS FROM ROLE analyst;
REVOKE IGNORE RLS FROM ROLE privileged_user;
REVOKE SELECT ON TABLE sensitive_data FROM RLS POLICY high_security_policy;

-- Assume role permissions
REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' FROM sales_user FOR COPY;
REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' FROM etl_user FOR UNLOAD;
REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' FROM ml_user FOR CREATE MODEL;
REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' FROM lambda_user FOR EXTERNAL FUNCTION;
REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' FROM analytics_team FOR ALL;
REVOKE ASSUMEROLE DEFAULT FROM legacy_user FOR COPY;
REVOKE ASSUMEROLE ALL FROM untrusted_user FOR ALL;

-- Multiple IAM roles
REVOKE ASSUMEROLE 'arn:aws:iam::123456789012:role/Role1', 'arn:aws:iam::123456789012:role/Role2' FROM multi_role_user FOR COPY;

-- External schema permissions (Spectrum integration)
REVOKE CREATE ON EXTERNAL SCHEMA spectrum_schema FROM IAM_ROLE 'spectrum_user';
REVOKE ALTER ON EXTERNAL SCHEMA spectrum_schema FROM IAM_ROLE 'schema_modifier';
REVOKE DROP ON EXTERNAL SCHEMA spectrum_schema FROM IAM_ROLE 'schema_admin';
REVOKE ALL ON EXTERNAL SCHEMA spectrum_schema FROM IAM_ROLE 'spectrum_admin';
REVOKE CREATE, ALTER, DROP ON EXTERNAL SCHEMA glue_schema FROM IAM_ROLE 'glue_user';

-- IAM role permissions for external objects
REVOKE SELECT ON EXTERNAL TABLE spectrum_schema.data FROM IAM_ROLE 'arn:aws:iam::123456789012:role/SpectrumRole';
REVOKE ALL ON EXTERNAL SCHEMA glue_schema FROM IAM_ROLE 'arn:aws:iam::123456789012:role/GlueRole';

-- Copy job permissions
REVOKE ALTER ON COPY JOB my_copy_job FROM job_modifier;
REVOKE DROP ON COPY JOB old_copy_job FROM job_cleaner;
REVOKE ALL ON COPY JOB critical_job FROM PUBLIC;
REVOKE ALTER, DROP ON COPY JOB job1, job2, job3 FROM job_admin;

-- Scoped permissions (for future objects)
REVOKE CREATE FOR SCHEMAS IN DATABASE mydb FROM schema_creator;
REVOKE SELECT FOR TABLES IN SCHEMA reporting FROM analyst_role;
REVOKE SELECT FOR TABLES IN SCHEMA finance DATABASE prod_db FROM finance_reader;
REVOKE EXECUTE FOR FUNCTIONS IN SCHEMA utilities FROM function_user;
REVOKE EXECUTE FOR FUNCTIONS IN SCHEMA math_lib DATABASE analytics_db FROM calculator_app;
REVOKE EXECUTE FOR PROCEDURES IN SCHEMA etl FROM etl_runner;
REVOKE EXECUTE FOR PROCEDURES IN SCHEMA maintenance DATABASE ops_db FROM maintenance_bot;
REVOKE USAGE FOR LANGUAGES IN DATABASE dev_db FROM developer;
REVOKE ALTER FOR COPY JOBS IN SCHEMA data_import FROM job_modifier;
REVOKE DROP FOR COPY JOBS IN SCHEMA data_import DATABASE warehouse FROM job_admin;

-- Complex combinations
REVOKE SELECT, INSERT, UPDATE ON TABLE orders FROM ROLE order_manager, GROUP sales_team;
REVOKE ALL ON TABLE customers FROM user1, ROLE customer_admin, GROUP support_team CASCADE;

-- Cross-object privileges
REVOKE CREATE ON DATABASE prod_db FROM schema_creator;
REVOKE CREATE ON SCHEMA new_schema FROM table_creator;
REVOKE SELECT ON TABLE new_schema.products FROM product_analyst;

-- Real-world examples
REVOKE SELECT ON TABLE fact_sales FROM ROLE analyst_role CASCADE;
REVOKE INSERT, UPDATE ON TABLE staging_orders FROM ROLE etl_role RESTRICT;
REVOKE ALL ON TABLE dim_customer FROM data_engineer CASCADE;
REVOKE USAGE ON SCHEMA analytics FROM ROLE business_analyst RESTRICT;
REVOKE CREATE ON DATABASE warehouse FROM ROLE data_architect CASCADE;

-- Service account patterns
REVOKE SELECT ON ALL TABLES IN SCHEMA reporting FROM airflow_service;
REVOKE EXECUTE ON ALL FUNCTIONS IN SCHEMA utilities FROM lambda_function;
REVOKE USAGE ON SCHEMA staging FROM kubernetes_pod;

-- Quoted identifiers
REVOKE SELECT ON TABLE "Order Details" FROM "Sales Manager";
REVOKE USAGE ON SCHEMA "Finance Reports" FROM "Finance Team";
REVOKE EXECUTE ON FUNCTION "Calculate Tax"() FROM "Tax Service";

-- Case variations
revoke select on table sales from analyst;
REVOKE SELECT ON TABLE SALES FROM ANALYST;
Revoke Select On Table Sales From Analyst;

-- Privilege combinations
REVOKE SELECT, INSERT ON TABLE logs FROM logger_service;
REVOKE UPDATE, DELETE, TRUNCATE ON TABLE temp_data FROM cleanup_service;
REVOKE DROP, ALTER ON TABLE migration_table FROM migration_tool;
REVOKE REFERENCES, SELECT ON TABLE lookup_table FROM referencing_table_owner;

-- Edge cases
REVOKE ALL ON TABLE single_column_table FROM single_user CASCADE;
REVOKE SELECT (id) ON TABLE wide_table FROM id_reader RESTRICT;
REVOKE USAGE ON SCHEMA "" FROM edge_case_user;

-- Multiple grantees with different types
REVOKE SELECT ON TABLE mixed_access_table FROM user1, ROLE reader_role, GROUP viewer_group, PUBLIC;
REVOKE USAGE ON SCHEMA shared_schema FROM service_user, ROLE service_role, GROUP service_group;

-- Nested role revocations
REVOKE ROLE junior_analyst FROM ROLE senior_analyst;
REVOKE ROLE team_lead FROM ROLE department_manager;
REVOKE ROLE read_only_role FROM ROLE power_user_role;

-- Function overloading scenarios
REVOKE EXECUTE ON FUNCTION calculate(integer) FROM math_user1;
REVOKE EXECUTE ON FUNCTION calculate(numeric) FROM math_user2;
REVOKE EXECUTE ON FUNCTION calculate(integer, integer) FROM math_user3;

-- Data type variations in function parameters
REVOKE EXECUTE ON FUNCTION process_string(varchar) FROM string_processor;
REVOKE EXECUTE ON FUNCTION process_string(text) FROM text_processor;
REVOKE EXECUTE ON FUNCTION handle_timestamp(timestamp) FROM time_handler;
REVOKE EXECUTE ON FUNCTION handle_numbers(numeric(10,2)) FROM number_handler;

-- Schema-qualified object names
REVOKE SELECT ON TABLE public.users FROM public_reader;
REVOKE EXECUTE ON FUNCTION analytics.calculate_metrics() FROM metrics_user;
REVOKE USAGE ON SCHEMA finance FROM finance_analyst;

-- Long identifier names
REVOKE SELECT ON TABLE very_long_table_name_for_comprehensive_sales_data FROM very_long_user_name_for_sales_analysis;
REVOKE EXECUTE ON FUNCTION extremely_long_function_name_for_complex_calculations() FROM user_with_very_long_descriptive_name;

-- Cleanup after grants with GRANT OPTION
REVOKE ALL PRIVILEGES ON TABLE revokable_table FROM revokable_user CASCADE;
REVOKE ALL ON ALL TABLES IN SCHEMA revokable_schema FROM revokable_role CASCADE;

-- Temporary and session-based revocations
REVOKE TEMPORARY ON DATABASE session_db FROM session_user;
REVOKE TEMP ON DATABASE temp_workspace FROM temp_worker;

-- GRANT OPTION revocations (keep privilege but remove ability to grant)
REVOKE GRANT OPTION FOR SELECT ON TABLE sensitive_data FROM trusted_user;
REVOKE GRANT OPTION FOR ALL ON DATABASE prod_db FROM admin_user;
REVOKE GRANT OPTION FOR USAGE ON SCHEMA reporting FROM lead_analyst;
REVOKE GRANT OPTION FOR EXECUTE ON FUNCTION calc() FROM senior_developer;
REVOKE GRANT OPTION FOR ALTER ON DATASHARE salesshare FROM share_admin;

-- Multiple privileges in one statement
REVOKE SELECT, INSERT, UPDATE, DELETE, TRUNCATE ON TABLE transactions FROM transaction_processor;
REVOKE CREATE, USAGE, ALTER, DROP ON SCHEMA dynamic_schema FROM schema_manager;
REVOKE CREATE, TEMPORARY, ALTER ON DATABASE test_db FROM test_user;

-- Revoke from users that were granted through groups/roles
REVOKE SELECT ON TABLE inherited_access FROM user_via_group CASCADE;
REVOKE USAGE ON SCHEMA inherited_schema FROM user_via_role RESTRICT;

-- Special cases with system tables and catalog access
REVOKE ACCESS CATALOG FROM ROLE limited_viewer;

-- Revoking combinations of column and table privileges
REVOKE SELECT ON TABLE orders FROM order_viewer;
REVOKE UPDATE (status, modified_date) ON TABLE orders FROM order_updater;
REVOKE ALL (sensitive_col1, sensitive_col2) ON TABLE secure_table FROM PUBLIC;

-- Conditional revocations based on object existence
REVOKE SELECT ON TABLE IF EXISTS optional_table FROM optional_user;
REVOKE USAGE ON SCHEMA IF EXISTS optional_schema FROM optional_role;