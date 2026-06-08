-- GRANT test cases based on AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_GRANT.html

-- Basic table privileges (from AWS documentation)
GRANT SELECT ON TABLE sales TO dbuser;
GRANT INSERT, UPDATE, DELETE ON TABLE sales TO dbuser;
GRANT ALL ON TABLE sales TO dbuser;
GRANT ALL PRIVILEGES ON TABLE sales TO dbuser;

-- Table privileges to multiple users
GRANT SELECT ON TABLE sales TO user1, user2, user3;
GRANT INSERT ON sales TO dbuser1, dbuser2;

-- Table privileges with WITH GRANT OPTION
GRANT SELECT ON TABLE sales TO dbuser WITH GRANT OPTION;
GRANT ALL ON TABLE sales TO dbuser WITH GRANT OPTION;

-- Column-level privileges (from AWS documentation)
GRANT SELECT (col1) ON TABLE t1 TO PUBLIC;
GRANT UPDATE (col1, col2) ON TABLE t1 TO user1;

-- Grant to PUBLIC
GRANT SELECT ON TABLE public_data TO PUBLIC;
GRANT USAGE ON SCHEMA public TO PUBLIC;

-- Grant to ROLE and GROUP
GRANT SELECT ON TABLE sales TO ROLE analyst_role;
GRANT INSERT ON TABLE logs TO GROUP etl_group;
GRANT ALL ON TABLE temp_data TO ROLE temp_role, GROUP admin_group;

-- All tables in schema (from AWS documentation)
GRANT SELECT ON ALL TABLES IN SCHEMA myschema TO dbuser;
GRANT ALL ON ALL TABLES IN SCHEMA reporting TO ROLE report_role;

-- Database privileges (from AWS documentation)
GRANT CREATE ON DATABASE mydb TO dbuser;
GRANT USAGE ON DATABASE prod_db TO readonly_user;
GRANT TEMPORARY ON DATABASE session_db TO temp_user;
GRANT TEMP ON DATABASE session_db TO temp_user;
GRANT ALTER ON DATABASE mydb TO admin_user;
GRANT ALL ON DATABASE mydb TO dba_role;
GRANT ALL PRIVILEGES ON DATABASE mydb TO admin_group;

-- Database privileges with WITH GRANT OPTION
GRANT CREATE ON DATABASE mydb TO dbuser WITH GRANT OPTION;
GRANT ALL ON DATABASE mydb TO admin WITH GRANT OPTION;

-- Schema privileges (from AWS documentation)
GRANT CREATE ON SCHEMA myschema TO dbuser;
GRANT USAGE ON SCHEMA reporting TO analyst;
GRANT ALTER ON SCHEMA dynamic_schema TO schema_admin;
GRANT DROP ON SCHEMA temp_schema TO cleanup_role;
GRANT ALL ON SCHEMA myschema TO schema_owner;
GRANT ALL PRIVILEGES ON SCHEMA myschema TO schema_admin;

-- Schema privileges with WITH GRANT OPTION
GRANT USAGE ON SCHEMA reporting TO analyst WITH GRANT OPTION;
GRANT ALL ON SCHEMA myschema TO admin WITH GRANT OPTION;

-- Function privileges (from AWS documentation)
GRANT EXECUTE ON FUNCTION f1() TO dbuser;
GRANT EXECUTE ON FUNCTION calculate_tax(numeric, numeric) TO tax_calculator;
GRANT ALL ON FUNCTION utility_func() TO utility_user;
GRANT ALL PRIVILEGES ON FUNCTION math_func(integer) TO math_user;

-- Function privileges with parameters
GRANT EXECUTE ON FUNCTION format_currency(numeric, varchar) TO report_user;
GRANT EXECUTE ON FUNCTION validate_email(varchar(255)) TO validation_service;
GRANT EXECUTE ON FUNCTION calculate_discount(decimal(10,2), integer) TO pricing_engine;

-- All functions in schema
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA analytics TO analyst_role;
GRANT ALL ON ALL FUNCTIONS IN SCHEMA utilities TO utility_group;

-- Procedure privileges
GRANT EXECUTE ON PROCEDURE proc1() TO dbuser;
GRANT EXECUTE ON PROCEDURE update_inventory(integer, varchar) TO inventory_user;
GRANT ALL ON PROCEDURE maintenance_proc() TO admin_role;

-- All procedures in schema
GRANT EXECUTE ON ALL PROCEDURES IN SCHEMA maintenance TO maintenance_role;
GRANT ALL ON ALL PROCEDURES IN SCHEMA etl TO etl_group;

-- Function and procedure privileges with WITH GRANT OPTION
GRANT EXECUTE ON FUNCTION f1() TO dbuser WITH GRANT OPTION;
GRANT EXECUTE ON PROCEDURE proc1() TO admin WITH GRANT OPTION;

-- Datashare privileges (from AWS documentation)
GRANT ALTER ON DATASHARE salesshare TO myuser WITH GRANT OPTION;
GRANT USAGE ON DATASHARE salesshare TO namespace 'a3f3ae8c-14e8-45ba-9eaa-b42b9b7ae635';

-- Role privileges (from AWS documentation)
GRANT ROLE sample_role1 TO reguser;
GRANT ROLE admin_role TO user1 WITH ADMIN OPTION;
GRANT ROLE analyst_role TO ROLE senior_analyst;

-- Multiple role grants
GRANT ROLE read_role TO user1, user2, user3;
GRANT ROLE admin_role TO ROLE super_admin, ROLE department_head;

-- Model privileges (Amazon Redshift ML)
GRANT EXECUTE ON MODEL customer_churn_model TO ml_user;
GRANT ALL ON MODEL sales_forecast_model TO ml_team;
GRANT ALL PRIVILEGES ON MODEL recommendation_model TO data_scientist;

-- Model privileges with WITH GRANT OPTION
GRANT EXECUTE ON MODEL fraud_detection_model TO security_analyst WITH GRANT OPTION;
GRANT ALL ON MODEL price_optimization_model TO pricing_team WITH GRANT OPTION;

-- Complex combinations
GRANT SELECT, INSERT, UPDATE ON TABLE orders TO ROLE order_manager, GROUP sales_team;
GRANT ALL ON TABLE customers TO user1 WITH GRANT OPTION, ROLE customer_admin, GROUP support_team;

-- Cross-object privileges
GRANT CREATE ON DATABASE prod_db TO schema_creator;
GRANT CREATE ON SCHEMA new_schema TO table_creator;
GRANT SELECT ON TABLE new_schema.products TO product_analyst;

-- Real-world examples
GRANT SELECT ON TABLE fact_sales TO ROLE analyst_role;
GRANT INSERT, UPDATE ON TABLE staging_orders TO ROLE etl_role;
GRANT ALL ON TABLE dim_customer TO data_engineer WITH GRANT OPTION;
GRANT USAGE ON SCHEMA analytics TO ROLE business_analyst;
GRANT CREATE ON DATABASE warehouse TO ROLE data_architect;

-- Service account patterns
GRANT SELECT ON ALL TABLES IN SCHEMA reporting TO airflow_service;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA utilities TO lambda_function;
GRANT USAGE ON SCHEMA staging TO kubernetes_pod;

-- Quoted identifiers
GRANT SELECT ON TABLE "Order Details" TO "Sales Manager";
GRANT USAGE ON SCHEMA "Finance Reports" TO "Finance Team";
GRANT EXECUTE ON FUNCTION "Calculate Tax"() TO "Tax Service";

-- Case variations
grant select on table sales to analyst;
GRANT SELECT ON TABLE SALES TO ANALYST;
Grant Select On Table Sales To Analyst;

-- Privilege combinations
GRANT SELECT, INSERT ON TABLE logs TO logger_service;
GRANT UPDATE, DELETE, TRUNCATE ON TABLE temp_data TO cleanup_service;
GRANT DROP, ALTER ON TABLE migration_table TO migration_tool;
GRANT REFERENCES, SELECT ON TABLE lookup_table TO referencing_table_owner;

-- Edge cases
GRANT ALL ON TABLE single_column_table TO single_user;
GRANT SELECT (id) ON TABLE wide_table TO id_reader;
GRANT USAGE ON SCHEMA "" TO edge_case_user;

-- Multiple grantees with different types
GRANT SELECT ON TABLE mixed_access_table TO user1, ROLE reader_role, GROUP viewer_group, PUBLIC;
GRANT USAGE ON SCHEMA shared_schema TO service_user, ROLE service_role, GROUP service_group;

-- Nested role grants
GRANT ROLE junior_analyst TO ROLE senior_analyst;
GRANT ROLE team_lead TO ROLE department_manager;
GRANT ROLE read_only_role TO ROLE power_user_role;

-- Function overloading scenarios
GRANT EXECUTE ON FUNCTION calculate(integer) TO math_user1;
GRANT EXECUTE ON FUNCTION calculate(numeric) TO math_user2;
GRANT EXECUTE ON FUNCTION calculate(integer, integer) TO math_user3;

-- Data type variations in function parameters
GRANT EXECUTE ON FUNCTION process_string(varchar) TO string_processor;
GRANT EXECUTE ON FUNCTION process_string(text) TO text_processor;
GRANT EXECUTE ON FUNCTION handle_timestamp(timestamp) TO time_handler;
GRANT EXECUTE ON FUNCTION handle_numbers(numeric(10,2)) TO number_handler;

-- Schema-qualified object names
GRANT SELECT ON TABLE public.users TO public_reader;
GRANT EXECUTE ON FUNCTION analytics.calculate_metrics() TO metrics_user;
GRANT USAGE ON SCHEMA finance TO finance_analyst;

-- Long identifier names
GRANT SELECT ON TABLE very_long_table_name_for_comprehensive_sales_data TO very_long_user_name_for_sales_analysis;
GRANT EXECUTE ON FUNCTION extremely_long_function_name_for_complex_calculations() TO user_with_very_long_descriptive_name;

-- Privilege revocation preparation (understanding what can be granted)
GRANT ALL PRIVILEGES ON TABLE revokable_table TO revokable_user WITH GRANT OPTION;
GRANT ALL ON ALL TABLES IN SCHEMA revokable_schema TO revokable_role WITH GRANT OPTION;

-- Temporary and session-based grants
GRANT TEMPORARY ON DATABASE session_db TO session_user;
GRANT TEMP ON DATABASE temp_workspace TO temp_worker;