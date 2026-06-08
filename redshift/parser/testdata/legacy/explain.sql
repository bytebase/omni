-- EXPLAIN command test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_EXPLAIN.html

-- Basic EXPLAIN statements
EXPLAIN SELECT * FROM table1;
explain select * from table1;

-- EXPLAIN with VERBOSE option
EXPLAIN VERBOSE SELECT * FROM table1;
explain verbose select * from table1;

-- Example from AWS documentation - basic EXPLAIN
EXPLAIN SELECT eventid, eventname, event.venueid, venuename
FROM event, venue
WHERE event.venueid = venue.venueid;

-- Example from AWS documentation - VERBOSE EXPLAIN
EXPLAIN VERBOSE SELECT eventid, eventname, event.venueid, venuename
FROM event, venue
WHERE event.venueid = venue.venueid;

-- EXPLAIN SELECT with various clauses
EXPLAIN SELECT customer_id, customer_name, order_count
FROM customers c
JOIN (
    SELECT customer_id, COUNT(*) as order_count
    FROM orders
    GROUP BY customer_id
) o ON c.customer_id = o.customer_id
WHERE c.registration_date >= '2023-01-01'
ORDER BY order_count DESC
LIMIT 100;

-- EXPLAIN VERBOSE with complex SELECT
EXPLAIN VERBOSE
SELECT
    category,
    AVG(price) as avg_price,
    COUNT(*) as product_count
FROM products
WHERE price > 100
GROUP BY category
HAVING COUNT(*) > 5
ORDER BY avg_price DESC;

-- EXPLAIN with window functions
EXPLAIN SELECT
    employee_id,
    department_id,
    salary,
    ROW_NUMBER() OVER (PARTITION BY department_id ORDER BY salary DESC) as rank_in_dept
FROM employees
WHERE hire_date > '2020-01-01';

-- EXPLAIN with CTE (Common Table Expression)
EXPLAIN WITH sales_summary AS (
    SELECT
        product_id,
        SUM(quantity) as total_quantity,
        SUM(amount) as total_amount
    FROM sales
    WHERE sale_date >= '2023-01-01'
    GROUP BY product_id
)
SELECT p.product_name, s.total_quantity, s.total_amount
FROM sales_summary s
JOIN products p ON s.product_id = p.product_id
ORDER BY s.total_amount DESC;

-- EXPLAIN VERBOSE with CTE
EXPLAIN VERBOSE WITH regional_sales AS (
    SELECT region, SUM(sales_amount) as total_sales
    FROM sales_data
    GROUP BY region
),
top_regions AS (
    SELECT region, total_sales
    FROM regional_sales
    WHERE total_sales > 1000000
)
SELECT * FROM top_regions ORDER BY total_sales DESC;

-- EXPLAIN with UNION
EXPLAIN SELECT customer_id, 'current' as status FROM current_customers
UNION ALL
SELECT customer_id, 'former' as status FROM former_customers
ORDER BY customer_id;

-- EXPLAIN with subqueries
EXPLAIN SELECT *
FROM products p
WHERE p.category_id IN (
    SELECT category_id
    FROM categories
    WHERE category_name IN ('Electronics', 'Books')
)
AND p.price > (
    SELECT AVG(price)
    FROM products
    WHERE category_id = p.category_id
);

-- EXPLAIN CREATE TABLE AS (CTAS) - Example from AWS documentation
EXPLAIN CREATE TABLE venue_nonulls
AS SELECT * FROM venue
WHERE venueseats IS NOT NULL;

-- EXPLAIN VERBOSE with CREATE TABLE AS
EXPLAIN VERBOSE CREATE TABLE sales_summary_2023
AS SELECT
    product_id,
    SUM(quantity) as total_quantity,
    SUM(amount) as total_amount,
    AVG(amount) as avg_amount
FROM sales
WHERE EXTRACT(YEAR FROM sale_date) = 2023
GROUP BY product_id;

-- EXPLAIN CREATE TEMPORARY TABLE AS
EXPLAIN CREATE TEMP TABLE temp_analysis
AS SELECT
    customer_id,
    COUNT(*) as order_count,
    SUM(total_amount) as total_spent
FROM orders
WHERE order_date >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY customer_id;

-- EXPLAIN CREATE TABLE AS with table attributes
EXPLAIN CREATE TABLE dist_sales
DISTSTYLE KEY
DISTKEY (customer_id)
SORTKEY (sale_date)
AS SELECT customer_id, sale_date, amount, product_id
FROM sales
WHERE sale_date >= '2023-01-01';

-- EXPLAIN CREATE TABLE AS with BACKUP option
EXPLAIN CREATE TABLE backup_customers
BACKUP YES
AS SELECT * FROM customers WHERE status = 'active';

-- EXPLAIN INSERT statement
EXPLAIN INSERT INTO target_table (col1, col2, col3)
SELECT source_col1, source_col2, source_col3
FROM source_table
WHERE created_date >= '2023-01-01';

-- EXPLAIN INSERT with VALUES
EXPLAIN INSERT INTO products (product_name, price, category_id)
VALUES
    ('Product A', 99.99, 1),
    ('Product B', 149.99, 2),
    ('Product C', 79.99, 1);

-- EXPLAIN UPDATE statement
EXPLAIN UPDATE customers
SET status = 'inactive',
    last_updated = CURRENT_TIMESTAMP
WHERE last_order_date < CURRENT_DATE - INTERVAL '365 days';

-- EXPLAIN UPDATE with JOIN
EXPLAIN UPDATE inventory i
SET quantity = i.quantity - o.quantity
FROM order_items o
WHERE i.product_id = o.product_id
AND o.order_date = CURRENT_DATE;

-- EXPLAIN DELETE statement
EXPLAIN DELETE FROM log_entries
WHERE log_date < CURRENT_DATE - INTERVAL '90 days';

-- EXPLAIN DELETE with subquery
EXPLAIN DELETE FROM orders
WHERE customer_id IN (
    SELECT customer_id
    FROM customers
    WHERE status = 'deleted'
);

-- EXPLAIN with complex joins
EXPLAIN SELECT
    c.customer_name,
    o.order_id,
    p.product_name,
    oi.quantity,
    oi.price
FROM customers c
INNER JOIN orders o ON c.customer_id = o.customer_id
INNER JOIN order_items oi ON o.order_id = oi.order_id
INNER JOIN products p ON oi.product_id = p.product_id
LEFT JOIN product_reviews pr ON p.product_id = pr.product_id
WHERE o.order_date >= '2023-01-01'
AND c.country = 'USA';

-- EXPLAIN with outer joins
EXPLAIN VERBOSE SELECT
    d.department_name,
    e.employee_name,
    e.salary
FROM departments d
LEFT OUTER JOIN employees e ON d.department_id = e.department_id
RIGHT OUTER JOIN locations l ON d.location_id = l.location_id
FULL OUTER JOIN managers m ON d.department_id = m.department_id;

-- EXPLAIN with aggregation and grouping
EXPLAIN SELECT
    EXTRACT(YEAR FROM order_date) as year,
    EXTRACT(MONTH FROM order_date) as month,
    COUNT(*) as order_count,
    SUM(total_amount) as total_revenue,
    AVG(total_amount) as avg_order_value
FROM orders
WHERE order_date >= '2022-01-01'
GROUP BY 1, 2
ORDER BY 1, 2;

-- EXPLAIN with HAVING clause
EXPLAIN SELECT
    category_id,
    COUNT(*) as product_count,
    AVG(price) as avg_price
FROM products
GROUP BY category_id
HAVING COUNT(*) > 10
AND AVG(price) < 500;

-- EXPLAIN with CASE expressions
EXPLAIN SELECT
    customer_id,
    CASE
        WHEN total_spent >= 10000 THEN 'Premium'
        WHEN total_spent >= 5000 THEN 'Gold'
        WHEN total_spent >= 1000 THEN 'Silver'
        ELSE 'Bronze'
    END as customer_tier,
    total_spent
FROM (
    SELECT
        customer_id,
        SUM(order_total) as total_spent
    FROM orders
    GROUP BY customer_id
) customer_totals;

-- EXPLAIN with CAST operations
EXPLAIN SELECT
    product_id,
    CAST(price AS INTEGER) as price_int,
    CAST(created_date AS VARCHAR(10)) as created_date_str,
    price::DECIMAL(10,2) as price_decimal
FROM products
WHERE CAST(price AS INTEGER) > 100;

-- EXPLAIN with date/time functions
EXPLAIN SELECT
    order_id,
    order_date,
    EXTRACT(YEAR FROM order_date) as order_year,
    EXTRACT(MONTH FROM order_date) as order_month,
    EXTRACT(DAY FROM order_date) as order_day,
    DATE_PART('dow', order_date) as day_of_week
FROM orders
WHERE order_date >= CURRENT_DATE - INTERVAL '1 year';

-- EXPLAIN with string functions and pattern matching
EXPLAIN SELECT
    customer_id,
    customer_name,
    email
FROM customers
WHERE customer_name LIKE 'John%'
AND email LIKE '%@gmail.com'
AND LENGTH(customer_name) > 5;

-- EXPLAIN with IN and EXISTS clauses
EXPLAIN SELECT *
FROM products p
WHERE p.category_id IN (1, 2, 3, 4, 5)
AND EXISTS (
    SELECT 1
    FROM order_items oi
    WHERE oi.product_id = p.product_id
    AND oi.order_date >= CURRENT_DATE - INTERVAL '30 days'
);

-- EXPLAIN with BETWEEN
EXPLAIN SELECT
    product_id,
    product_name,
    price
FROM products
WHERE price BETWEEN 50 AND 200
AND created_date BETWEEN '2023-01-01' AND '2023-12-31';

-- EXPLAIN with NULL handling
EXPLAIN SELECT
    customer_id,
    customer_name,
    COALESCE(phone, 'N/A') as contact_phone,
    CASE WHEN email IS NULL THEN 'No Email' ELSE email END as email_status
FROM customers
WHERE phone IS NOT NULL
OR email IS NOT NULL;

-- EXPLAIN with DISTINCT
EXPLAIN SELECT DISTINCT
    category_id,
    supplier_id
FROM products
WHERE price > 100;

-- EXPLAIN with GROUP BY ROLLUP
EXPLAIN SELECT
    category_id,
    supplier_id,
    COUNT(*) as product_count,
    AVG(price) as avg_price
FROM products
GROUP BY ROLLUP(category_id, supplier_id);

-- EXPLAIN with GROUP BY CUBE
EXPLAIN SELECT
    EXTRACT(YEAR FROM order_date) as year,
    EXTRACT(QUARTER FROM order_date) as quarter,
    SUM(total_amount) as revenue
FROM orders
GROUP BY CUBE(1, 2);

-- EXPLAIN with GROUPING SETS
EXPLAIN SELECT
    category_id,
    brand_id,
    COUNT(*) as product_count
FROM products
GROUP BY GROUPING SETS (
    (category_id),
    (brand_id),
    (category_id, brand_id),
    ()
);

-- EXPLAIN with LIMIT and OFFSET
EXPLAIN SELECT
    product_id,
    product_name,
    price
FROM products
ORDER BY price DESC
LIMIT 50 OFFSET 100;

-- Mixed case variations
Explain Select * From table1;
EXPLAIN verbose SELECT col1, col2 FROM table2 WHERE col1 > 100;
explain VERBOSE select count(*) from orders group by customer_id;

-- Complex query with multiple features
EXPLAIN VERBOSE
WITH monthly_sales AS (
    SELECT
        DATE_TRUNC('month', order_date) as month,
        customer_id,
        SUM(total_amount) as monthly_total
    FROM orders
    WHERE order_date >= '2023-01-01'
    GROUP BY 1, 2
),
customer_metrics AS (
    SELECT
        customer_id,
        COUNT(DISTINCT month) as active_months,
        AVG(monthly_total) as avg_monthly_total,
        MAX(monthly_total) as max_monthly_total
    FROM monthly_sales
    GROUP BY customer_id
)
SELECT
    c.customer_name,
    cm.active_months,
    cm.avg_monthly_total,
    cm.max_monthly_total,
    CASE
        WHEN cm.avg_monthly_total >= 1000 THEN 'High Value'
        WHEN cm.avg_monthly_total >= 500 THEN 'Medium Value'
        ELSE 'Low Value'
    END as customer_segment
FROM customer_metrics cm
JOIN customers c ON cm.customer_id = c.customer_id
WHERE cm.active_months >= 3
ORDER BY cm.avg_monthly_total DESC
LIMIT 100;