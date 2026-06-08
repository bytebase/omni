-- SHOW GRANTS test cases based on AWS Redshift documentation
-- https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_GRANTS.html

-- ==============================================
-- 1. Show grants on database objects
-- ==============================================

-- Show grants on a database
SHOW GRANTS ON DATABASE dev;
SHOW GRANTS ON DATABASE mydb;
SHOW GRANTS ON DATABASE sales_db;

-- Show grants on a database for a specific user
SHOW GRANTS ON DATABASE dev FOR alice;
SHOW GRANTS ON DATABASE mydb FOR bob;
SHOW GRANTS ON DATABASE sales_db FOR ROLE analyst_role;
SHOW GRANTS ON DATABASE test_db FOR PUBLIC;

-- Show grants on a database with limit
SHOW GRANTS ON DATABASE dev LIMIT 10;
SHOW GRANTS ON DATABASE mydb FOR alice LIMIT 100;
SHOW GRANTS ON DATABASE sales_db FOR ROLE admin_role LIMIT 5000;

-- ==============================================
-- 2. Show grants on schemas
-- ==============================================

-- Show grants on a schema (simple name)
SHOW GRANTS ON SCHEMA demo;
SHOW GRANTS ON SCHEMA public;
SHOW GRANTS ON SCHEMA sales;

-- Show grants on a schema (qualified with database)
SHOW GRANTS ON SCHEMA mydb.sales;
SHOW GRANTS ON SCHEMA dev.analytics;
SHOW GRANTS ON SCHEMA prod_db.reporting;

-- Show grants on a schema for specific principals
SHOW GRANTS ON SCHEMA demo FOR alice;
SHOW GRANTS ON SCHEMA public FOR ROLE reader_role;
SHOW GRANTS ON SCHEMA mydb.sales FOR bob;
SHOW GRANTS ON SCHEMA dev.analytics FOR PUBLIC;

-- Show grants on a schema with limit
SHOW GRANTS ON SCHEMA demo LIMIT 50;
SHOW GRANTS ON SCHEMA mydb.sales FOR alice LIMIT 1000;

-- ==============================================
-- 3. Show grants on tables
-- ==============================================

-- Show grants on a table (simple name)
SHOW GRANTS ON t4;
SHOW GRANTS ON customers;
SHOW GRANTS ON orders;

-- Show grants on a table with TABLE keyword
SHOW GRANTS ON TABLE t4;
SHOW GRANTS ON TABLE public.t4;
SHOW GRANTS ON TABLE customers;

-- Show grants on a table (two-part notation)
SHOW GRANTS ON TABLE demo_schema.t3;
SHOW GRANTS ON TABLE sales.customers;
SHOW GRANTS ON TABLE analytics.events;
SHOW GRANTS ON demo_schema.t3 FOR alice;
SHOW GRANTS ON sales.customers FOR ROLE sales_role;

-- Show grants on a table (three-part notation)
SHOW GRANTS ON TABLE demo_db.demo_schema.t3 FOR alice;
SHOW GRANTS ON TABLE mydb.sales.orders FOR bob;
SHOW GRANTS ON TABLE prod_db.analytics.metrics FOR ROLE analyst_role;
SHOW GRANTS ON demo_db.demo_schema.t3 FOR PUBLIC;

-- Show grants on a table with limit
SHOW GRANTS ON TABLE customers LIMIT 20;
SHOW GRANTS ON TABLE sales.orders FOR alice LIMIT 100;
SHOW GRANTS ON TABLE mydb.sales.products FOR ROLE viewer_role LIMIT 500;

-- ==============================================
-- 4. Show grants on functions
-- ==============================================

-- Show grants on a function (without parameters)
SHOW GRANTS ON FUNCTION calculate_total;
SHOW GRANTS ON FUNCTION public.get_user_count;
SHOW GRANTS ON FUNCTION sales.compute_revenue;
SHOW GRANTS ON FUNCTION mydb.analytics.process_data;

-- Show grants on a function (with empty parameter list)
SHOW GRANTS ON FUNCTION calculate_total();
SHOW GRANTS ON FUNCTION public.get_user_count();
SHOW GRANTS ON FUNCTION sales.compute_revenue();

-- Show grants on a function (with parameters)
SHOW GRANTS ON FUNCTION add_numbers(integer, integer);
SHOW GRANTS ON FUNCTION public.concat_strings(text, text);
SHOW GRANTS ON FUNCTION sales.calculate_discount(numeric, numeric);
SHOW GRANTS ON FUNCTION mydb.util.format_date(date);

-- Show grants on a function with named parameters
SHOW GRANTS ON FUNCTION calculate_tax(amount numeric, rate numeric);
SHOW GRANTS ON FUNCTION public.process_order(order_id integer, user_id integer);
SHOW GRANTS ON FUNCTION analytics.aggregate_metrics(start_date date, end_date date, metric_type varchar);

-- Show grants on a function for specific principals
SHOW GRANTS ON FUNCTION calculate_total FOR alice;
SHOW GRANTS ON FUNCTION public.get_user_count() FOR ROLE app_role;
SHOW GRANTS ON FUNCTION sales.compute_revenue(date, date) FOR bob;
SHOW GRANTS ON FUNCTION mydb.analytics.process_data FOR PUBLIC;

-- Show grants on a function with limit
SHOW GRANTS ON FUNCTION calculate_total LIMIT 10;
SHOW GRANTS ON FUNCTION public.get_user_count() FOR alice LIMIT 50;

-- ==============================================
-- 5. Show grants for users and roles
-- ==============================================

-- Show grants for a user
SHOW GRANTS FOR alice;
SHOW GRANTS FOR bob;
SHOW GRANTS FOR admin_user;
SHOW GRANTS FOR app_service_account;

-- Show grants for a role
SHOW GRANTS FOR ROLE analyst_role;
SHOW GRANTS FOR ROLE admin_role;
SHOW GRANTS FOR ROLE reader_role;
SHOW GRANTS FOR ROLE developer_role;

-- Show grants with limit
SHOW GRANTS FOR alice LIMIT 100;
SHOW GRANTS FOR bob LIMIT 500;
SHOW GRANTS FOR ROLE analyst_role LIMIT 1000;
SHOW GRANTS FOR ROLE admin_role LIMIT 10000;

-- ==============================================
-- 6. Edge cases and complex examples
-- ==============================================

-- Maximum limit value
SHOW GRANTS ON DATABASE dev LIMIT 10000;
SHOW GRANTS FOR alice LIMIT 10000;

-- Minimum limit value
SHOW GRANTS ON SCHEMA public LIMIT 0;
SHOW GRANTS FOR ROLE test_role LIMIT 1;

-- Mixed case identifiers
SHOW GRANTS ON DATABASE Dev;
SHOW GRANTS ON SCHEMA Public;
SHOW GRANTS ON TABLE Customers;
SHOW GRANTS FOR Alice;
SHOW GRANTS FOR ROLE Admin_Role;

-- Identifiers with underscores and numbers
SHOW GRANTS ON DATABASE test_db_2023;
SHOW GRANTS ON SCHEMA schema_v2;
SHOW GRANTS ON TABLE orders_2023_q4;
SHOW GRANTS ON FUNCTION calc_metrics_v3(integer, date);
SHOW GRANTS FOR user_123;
SHOW GRANTS FOR ROLE role_test_456;

-- Complex function signatures
SHOW GRANTS ON FUNCTION mydb.analytics.complex_calc(p1 integer, p2 numeric, p3 varchar, p4 date, p5 timestamp);
SHOW GRANTS ON FUNCTION process_array(arr integer[], delimiter text);
SHOW GRANTS ON FUNCTION validate_json(data json, schema text);