-- LOCK command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_LOCK.html

-- Basic LOCK statement (defaults to ACCESS EXCLUSIVE mode)
LOCK event;
LOCK table1;
lock users;

-- LOCK with TABLE keyword
LOCK TABLE event;
LOCK TABLE sales;
lock table customers;

-- Multiple tables (from AWS documentation example)
LOCK event, sales;
LOCK TABLE event, sales;
lock table users, orders, products;

-- Schema-qualified table names
LOCK public.event;
LOCK TABLE schema1.table1;
LOCK myschema.users, public.orders;

-- LOCK with explicit lock modes
-- ACCESS SHARE (least restrictive)
LOCK event IN ACCESS SHARE MODE;
LOCK TABLE sales IN ACCESS SHARE MODE;
lock table users in access share mode;

-- ROW SHARE
LOCK event IN ROW SHARE MODE;
LOCK TABLE sales IN ROW SHARE MODE;
lock table users in row share mode;

-- ROW EXCLUSIVE
LOCK event IN ROW EXCLUSIVE MODE;
LOCK TABLE sales IN ROW EXCLUSIVE MODE;
lock table users in row exclusive mode;

-- SHARE UPDATE EXCLUSIVE
LOCK event IN SHARE UPDATE EXCLUSIVE MODE;
LOCK TABLE sales IN SHARE UPDATE EXCLUSIVE MODE;
lock table users in share update exclusive mode;

-- SHARE
LOCK event IN SHARE MODE;
LOCK TABLE sales IN SHARE MODE;
lock table users in share mode;

-- SHARE ROW EXCLUSIVE
LOCK event IN SHARE ROW EXCLUSIVE MODE;
LOCK TABLE sales IN SHARE ROW EXCLUSIVE MODE;
lock table users in share row exclusive mode;

-- EXCLUSIVE
LOCK event IN EXCLUSIVE MODE;
LOCK TABLE sales IN EXCLUSIVE MODE;
lock table users in exclusive mode;

-- ACCESS EXCLUSIVE (most restrictive, default mode)
LOCK event IN ACCESS EXCLUSIVE MODE;
LOCK TABLE sales IN ACCESS EXCLUSIVE MODE;
lock table users in access exclusive mode;

-- LOCK with NOWAIT option
LOCK event NOWAIT;
LOCK TABLE sales NOWAIT;
lock table users nowait;

-- LOCK with explicit mode and NOWAIT
LOCK event IN ACCESS SHARE MODE NOWAIT;
LOCK TABLE sales IN ROW EXCLUSIVE MODE NOWAIT;
LOCK users IN EXCLUSIVE MODE NOWAIT;
lock table orders in share mode nowait;

-- Multiple tables with lock mode
LOCK event, sales IN ACCESS SHARE MODE;
LOCK TABLE users, orders IN ROW EXCLUSIVE MODE;
lock products, categories in exclusive mode;

-- Multiple tables with lock mode and NOWAIT
LOCK event, sales IN ACCESS SHARE MODE NOWAIT;
LOCK TABLE users, orders IN ROW EXCLUSIVE MODE NOWAIT;
lock products, categories in exclusive mode nowait;

-- Schema-qualified tables with modes
LOCK public.event IN SHARE MODE;
LOCK TABLE schema1.users, schema2.orders IN ROW EXCLUSIVE MODE;
LOCK mydb.public.sales IN ACCESS EXCLUSIVE MODE NOWAIT;

-- Mixed case variations
Lock Event;
LOCK table Sales;
lock TABLE Users in ACCESS share MODE;
Lock Table orders IN row EXCLUSIVE mode NOWAIT;

-- Real-world examples from AWS documentation context
-- Lock tables before performing operations that need consistency
LOCK event, sales;
LOCK TABLE venue IN ACCESS EXCLUSIVE MODE;

-- Lock with different modes for different scenarios
-- For SELECT operations that need consistent view
LOCK reporting_table IN ACCESS SHARE MODE;

-- For UPDATE/INSERT/DELETE operations
LOCK transactional_table IN ROW EXCLUSIVE MODE;

-- For schema changes
LOCK target_table IN ACCESS EXCLUSIVE MODE;

-- Complex scenarios with multiple schema-qualified tables
LOCK public.fact_sales, staging.temp_data IN SHARE MODE NOWAIT;
LOCK TABLE analytics.user_metrics, public.session_data IN ROW EXCLUSIVE MODE;

-- Edge cases with identifier names
LOCK "table with spaces";
LOCK TABLE "Mixed_Case_Table" IN ACCESS SHARE MODE;
LOCK "schema"."table" IN EXCLUSIVE MODE NOWAIT;