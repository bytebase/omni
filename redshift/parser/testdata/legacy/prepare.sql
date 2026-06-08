-- PREPARE command test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_PREPARE.html

-- Basic PREPARE statement without parameter types
PREPARE simple_plan AS SELECT * FROM users;

-- PREPARE with single parameter type
PREPARE prep_select_plan (int) AS SELECT * FROM prep1 WHERE c1 = $1;

-- PREPARE with multiple parameter types
PREPARE prep_insert_plan (int, char) AS INSERT INTO prep1 VALUES ($1, $2);

-- Examples from AWS documentation:

-- Example 1: Create a table and prepare an insert statement
PREPARE prep_insert_plan (int, char) AS INSERT INTO prep1 VALUES ($1, $2);

-- Example 2: Prepare a select statement with parameter
PREPARE prep_select_plan (int) AS SELECT * FROM prep1 WHERE c1 = $1;

-- Data type variations:

-- Character data types
PREPARE char_plan (char) AS INSERT INTO test VALUES ($1);
PREPARE char_precision (char(10)) AS INSERT INTO test VALUES ($1);
PREPARE varchar_plan (varchar) AS INSERT INTO test VALUES ($1);
PREPARE varchar_precision (varchar(255)) AS INSERT INTO test VALUES ($1);
PREPARE text_plan (text) AS INSERT INTO test VALUES ($1);

-- Numeric data types
PREPARE smallint_plan (smallint) AS INSERT INTO test VALUES ($1);
PREPARE integer_plan (integer) AS INSERT INTO test VALUES ($1);
PREPARE int_plan (int) AS INSERT INTO test VALUES ($1);
PREPARE bigint_plan (bigint) AS INSERT INTO test VALUES ($1);
PREPARE decimal_plan (decimal) AS INSERT INTO test VALUES ($1);
PREPARE decimal_precision (decimal(10,2)) AS INSERT INTO test VALUES ($1);
PREPARE numeric_plan (numeric) AS INSERT INTO test VALUES ($1);
PREPARE numeric_precision (numeric(18,4)) AS INSERT INTO test VALUES ($1);
PREPARE real_plan (real) AS INSERT INTO test VALUES ($1);
PREPARE double_precision_plan (double precision) AS INSERT INTO test VALUES ($1);
PREPARE float_plan (float) AS INSERT INTO test VALUES ($1);
PREPARE float_precision (float(24)) AS INSERT INTO test VALUES ($1);

-- Date/time data types
PREPARE date_plan (date) AS INSERT INTO test VALUES ($1);
PREPARE time_plan (time) AS INSERT INTO test VALUES ($1);
PREPARE time_precision (time(6)) AS INSERT INTO test VALUES ($1);
PREPARE timestamp_plan (timestamp) AS INSERT INTO test VALUES ($1);
PREPARE timestamp_precision (timestamp(6)) AS INSERT INTO test VALUES ($1);
PREPARE timestamptz_plan (timestamptz) AS INSERT INTO test VALUES ($1);

-- Boolean data type
PREPARE boolean_plan (boolean) AS INSERT INTO test VALUES ($1);
PREPARE bool_plan (bool) AS INSERT INTO test VALUES ($1);

-- Multiple parameter types
PREPARE multi_type_plan (int, varchar(50), date) AS
    INSERT INTO users (id, name, created_date) VALUES ($1, $2, $3);

PREPARE complex_types (bigint, decimal(15,2), timestamp, boolean, text) AS
    INSERT INTO orders (order_id, amount, order_date, is_active, notes)
    VALUES ($1, $2, $3, $4, $5);

-- SELECT statement variations:

-- Basic SELECT with WHERE clause
PREPARE select_by_id (int) AS SELECT * FROM customers WHERE customer_id = $1;

-- SELECT with multiple conditions
PREPARE select_by_criteria (varchar(50), date, boolean) AS
    SELECT * FROM orders WHERE status = $1 AND order_date >= $2 AND is_processed = $3;

-- SELECT with JOIN
PREPARE select_with_join (int) AS
    SELECT c.name, o.total FROM customers c JOIN orders o ON c.id = o.customer_id
    WHERE c.id = $1;

-- SELECT with GROUP BY
PREPARE select_aggregated (date) AS
    SELECT customer_id, COUNT(*), SUM(total) FROM orders
    WHERE order_date >= $1 GROUP BY customer_id;

-- SELECT with ORDER BY and LIMIT
PREPARE select_top_customers (int) AS
    SELECT customer_id, SUM(total) as total_spent FROM orders
    GROUP BY customer_id ORDER BY total_spent DESC LIMIT $1;

-- INSERT statement variations:

-- Basic INSERT with VALUES
PREPARE insert_customer (varchar(100), varchar(50)) AS
    INSERT INTO customers (name, email) VALUES ($1, $2);

-- INSERT with multiple value sets (using single parameter per value)
PREPARE insert_multiple_products (varchar(50), decimal(10,2)) AS
    INSERT INTO products (product_name, price) VALUES ($1, $2);

-- INSERT with SELECT
PREPARE insert_from_select (date) AS
    INSERT INTO monthly_summary (customer_id, total_orders, total_amount)
    SELECT customer_id, COUNT(*), SUM(amount) FROM orders
    WHERE order_date >= $1 GROUP BY customer_id;

-- UPDATE statement variations:

-- Basic UPDATE
PREPARE update_customer_email (varchar(100), int) AS
    UPDATE customers SET email = $1 WHERE customer_id = $2;

-- UPDATE with multiple columns
PREPARE update_order_status (varchar(20), timestamp, int) AS
    UPDATE orders SET status = $1, updated_at = $2 WHERE order_id = $3;

-- UPDATE with condition using multiple parameters
PREPARE update_product_price (decimal(10,2), varchar(50), boolean) AS
    UPDATE products SET price = $1 WHERE category = $2 AND is_active = $3;

-- DELETE statement variations:

-- Basic DELETE
PREPARE delete_customer (int) AS DELETE FROM customers WHERE customer_id = $1;

-- DELETE with multiple conditions
PREPARE delete_old_orders (date, varchar(20)) AS
    DELETE FROM orders WHERE order_date < $1 AND status = $2;

-- DELETE with subquery condition
PREPARE delete_inactive_customers (date) AS
    DELETE FROM customers WHERE customer_id NOT IN
    (SELECT DISTINCT customer_id FROM orders WHERE order_date >= $1);

-- Complex query examples:

-- Window function query
PREPARE sales_ranking (date, date) AS
    SELECT salesperson, SUM(amount) as total_sales,
           RANK() OVER (ORDER BY SUM(amount) DESC) as rank
    FROM sales WHERE sale_date BETWEEN $1 AND $2
    GROUP BY salesperson;

-- CTE (Common Table Expression) query
PREPARE monthly_report (date) AS
    WITH monthly_sales AS (
        SELECT customer_id, SUM(amount) as total
        FROM orders WHERE order_date >= $1
        GROUP BY customer_id
    )
    SELECT c.name, ms.total FROM customers c
    JOIN monthly_sales ms ON c.id = ms.customer_id
    ORDER BY ms.total DESC;

-- Subquery examples
PREPARE customers_above_average (date) AS
    SELECT * FROM customers WHERE customer_id IN (
        SELECT customer_id FROM orders
        WHERE order_date >= $1
        GROUP BY customer_id
        HAVING SUM(amount) > (SELECT AVG(total_amount) FROM (
            SELECT customer_id, SUM(amount) as total_amount
            FROM orders GROUP BY customer_id
        ) avg_calc)
    );

-- Parameter placeholder variations (up to the documented limit):

-- Multiple parameters in different positions
PREPARE param_positions (int, varchar(50), date, decimal(10,2), boolean) AS
    SELECT * FROM transactions
    WHERE user_id = $1 AND description LIKE $2 AND trans_date >= $3
    AND amount >= $4 AND is_verified = $5;

-- Parameters in different clauses
PREPARE multi_clause_params (varchar(50), int, decimal(10,2)) AS
    SELECT category, COUNT(*) as count, AVG(price) as avg_price
    FROM products
    WHERE category = $1
    GROUP BY category
    HAVING COUNT(*) >= $2 AND AVG(price) <= $3;

-- Case sensitivity tests:
prepare lowercase_plan (int) as select * from test where id = $1;
PREPARE uppercase_plan (INT) AS SELECT * FROM TEST WHERE ID = $1;
Prepare mixed_case_plan (Integer) As Select * From test Where id = $1;

-- Plan name variations:
PREPARE simple_name (int) AS SELECT $1;
PREPARE plan_with_underscore (int) AS SELECT $1;
PREPARE plan123 (int) AS SELECT $1;
PREPARE very_long_plan_name_with_many_words (int) AS SELECT $1;
PREPARE a (int) AS SELECT $1;
PREPARE z (int) AS SELECT $1;
PREPARE plan_1 (int) AS SELECT $1;
PREPARE my_plan_v2 (int) AS SELECT $1;

-- Real-world business logic examples:

-- Customer management
PREPARE get_customer_orders (int, date) AS
    SELECT o.order_id, o.order_date, o.total, o.status
    FROM orders o WHERE o.customer_id = $1 AND o.order_date >= $2
    ORDER BY o.order_date DESC;

-- Inventory management
PREPARE update_stock_level (int, int) AS
    UPDATE products SET stock_quantity = stock_quantity - $2 WHERE product_id = $1;

-- Financial reporting
PREPARE revenue_by_period (date, date) AS
    SELECT DATE_TRUNC('month', order_date) as month, SUM(total) as revenue
    FROM orders WHERE order_date BETWEEN $1 AND $2
    GROUP BY DATE_TRUNC('month', order_date)
    ORDER BY month;

-- User authentication
PREPARE validate_user (varchar(100), varchar(255)) AS
    SELECT user_id, role, last_login FROM users
    WHERE email = $1 AND password_hash = $2 AND is_active = true;

-- Analytics queries
PREPARE top_selling_products (int, date) AS
    SELECT p.product_name, SUM(oi.quantity) as total_sold, SUM(oi.quantity * oi.price) as revenue
    FROM products p
    JOIN order_items oi ON p.product_id = oi.product_id
    JOIN orders o ON oi.order_id = o.order_id
    WHERE o.order_date >= $2
    GROUP BY p.product_id, p.product_name
    ORDER BY total_sold DESC
    LIMIT $1;

-- Data cleanup operations
PREPARE archive_old_data (date) AS
    INSERT INTO archived_orders SELECT * FROM orders WHERE order_date < $1;

PREPARE cleanup_temp_data (timestamp) AS
    DELETE FROM temp_processing_log WHERE created_at < $1;

-- Redshift-specific features in prepared statements:

-- UNLOAD command preparation (simplified representation)
PREPARE export_customer_data (date) AS
    SELECT customer_id, name, email FROM customers WHERE created_date >= $1;

-- Advanced analytics with window functions
PREPARE customer_lifetime_value (int) AS
    SELECT customer_id,
           SUM(total) OVER (PARTITION BY customer_id ORDER BY order_date
                           ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) as running_total,
           LAG(order_date, 1) OVER (PARTITION BY customer_id ORDER BY order_date) as prev_order_date
    FROM orders WHERE customer_id = $1 ORDER BY order_date;

-- JSON operations (if supported)
PREPARE json_extract_query (varchar(100)) AS
    SELECT customer_id, JSON_EXTRACT_PATH_TEXT(profile_data, $1) as extracted_value
    FROM customer_profiles WHERE profile_data IS NOT NULL;

-- Working with arrays (if supported)
PREPARE array_operations (int) AS
    SELECT order_id, item_ids[1:$1] as first_n_items
    FROM order_summary WHERE ARRAY_LENGTH(item_ids, 1) >= $1;

-- Error handling scenarios (these should still parse correctly):

-- Minimum parameters (1 parameter)
PREPARE single_param (int) AS SELECT $1;

-- Edge case with parameter at end
PREPARE param_at_end (varchar(50)) AS
    SELECT * FROM logs WHERE message LIKE '%error%' ORDER BY timestamp DESC LIMIT $1;

-- Multiple tables with parameters
PREPARE cross_table_query (int, varchar(50), date) AS
    SELECT u.name, p.title, o.order_date
    FROM users u
    JOIN profiles p ON u.id = p.user_id
    JOIN orders o ON u.id = o.customer_id
    WHERE u.id = $1 AND p.department = $2 AND o.order_date >= $3;

-- Nested function calls with parameters
PREPARE nested_functions (varchar(50), int) AS
    SELECT UPPER(SUBSTRING($1, 1, $2)) as formatted_text;

-- Mathematical operations with parameters
PREPARE math_operations (decimal(10,2), decimal(10,2), int) AS
    SELECT $1 * $2 as product, POWER($1, $3) as power_result;

-- String operations with parameters
PREPARE string_operations (varchar(100), varchar(10)) AS
    SELECT CONCAT($1, ' - ', $2) as combined,
           LENGTH($1) as text_length,
           POSITION($2 IN $1) as substring_position;