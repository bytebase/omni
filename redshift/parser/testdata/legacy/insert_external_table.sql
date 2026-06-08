-- Redshift INSERT into external tables test cases
-- Based on https://docs.aws.amazon.com/redshift/latest/dg/r_INSERT_external_table.html

-- Basic INSERT with SELECT from local table
INSERT INTO spectrum.lineitem SELECT * FROM local_lineitem;

-- INSERT into non-partitioned external table
INSERT INTO spectrum.customers
SELECT customer_id, customer_name, customer_email, signup_date
FROM staging.new_customers;

-- INSERT with WHERE clause filtering
INSERT INTO spectrum.orders
SELECT order_id, customer_id, order_date, total_amount
FROM local_orders
WHERE order_date >= '2024-01-01';

-- INSERT with JOIN between local tables
INSERT INTO spectrum.order_details
SELECT od.order_id, od.product_id, p.product_name, od.quantity, od.price
FROM local_order_details od
JOIN local_products p ON od.product_id = p.product_id;

-- Static partitioning - adding specific partition values
INSERT INTO spectrum.customer
SELECT name, age, gender, 'May', 28
FROM local_customer;

-- Static partitioning with multiple partition columns
INSERT INTO spectrum.sales_data
SELECT product_id, quantity, price, customer_id, '2024', '05', '15'
FROM daily_sales
WHERE sale_date = '2024-05-15';

-- Dynamic partitioning - partition values from source data
INSERT INTO spectrum.customer
SELECT name, age, gender, month, day
FROM local_customer;

-- Dynamic partitioning with computed partition values
INSERT INTO spectrum.web_logs
SELECT request_id, user_id, page_url, response_time,
       YEAR(log_timestamp), MONTH(log_timestamp), DAY(log_timestamp)
FROM raw_web_logs;

-- INSERT with aggregation
INSERT INTO spectrum.daily_summaries
SELECT product_category,
       COUNT(*) as total_orders,
       SUM(quantity) as total_quantity,
       AVG(price) as avg_price,
       '2024-05-15'
FROM local_orders
WHERE order_date = '2024-05-15'
GROUP BY product_category;

-- INSERT with complex SELECT including CTEs (Common Table Expressions)
INSERT INTO spectrum.customer_metrics
WITH customer_totals AS (
    SELECT customer_id,
           COUNT(DISTINCT order_id) as order_count,
           SUM(total_amount) as lifetime_value
    FROM orders
    GROUP BY customer_id
)
SELECT c.customer_id, c.customer_name, ct.order_count, ct.lifetime_value, CURRENT_DATE
FROM customers c
JOIN customer_totals ct ON c.customer_id = ct.customer_id;

-- INSERT with CASE expressions for partition columns
INSERT INTO spectrum.customer_segments
SELECT customer_id,
       customer_name,
       total_purchases,
       CASE
           WHEN total_purchases >= 1000 THEN 'premium'
           WHEN total_purchases >= 100 THEN 'regular'
           ELSE 'basic'
       END,
       YEAR(CURRENT_DATE),
       MONTH(CURRENT_DATE)
FROM customer_purchase_summary;

-- INSERT with multiple JOINs
INSERT INTO spectrum.product_performance
SELECT p.product_id,
       p.product_name,
       p.category,
       COUNT(DISTINCT o.order_id) as order_count,
       SUM(od.quantity) as total_quantity,
       SUM(od.quantity * od.unit_price) as total_revenue
FROM products p
JOIN order_details od ON p.product_id = od.product_id
JOIN orders o ON od.order_id = o.order_id
WHERE o.order_date >= DATEADD(month, -3, CURRENT_DATE)
GROUP BY p.product_id, p.product_name, p.category;

-- INSERT with window functions
INSERT INTO spectrum.customer_rankings
SELECT customer_id,
       customer_name,
       total_purchases,
       RANK() OVER (ORDER BY total_purchases DESC) as purchase_rank,
       NTILE(10) OVER (ORDER BY total_purchases DESC) as decile,
       CURRENT_DATE
FROM customer_summary;

-- INSERT with UNION ALL
INSERT INTO spectrum.all_events
SELECT event_id, 'order' as event_type, customer_id, event_timestamp
FROM order_events
UNION ALL
SELECT event_id, 'return' as event_type, customer_id, event_timestamp
FROM return_events;

-- INSERT with subquery in SELECT list
INSERT INTO spectrum.product_comparisons
SELECT p1.product_id,
       p1.product_name,
       p1.price,
       (SELECT AVG(price) FROM products p2 WHERE p2.category = p1.category) as category_avg_price,
       p1.category
FROM products p1;

-- INSERT with EXISTS clause
INSERT INTO spectrum.active_customers
SELECT c.customer_id, c.customer_name, c.email, CURRENT_DATE
FROM customers c
WHERE EXISTS (
    SELECT 1
    FROM orders o
    WHERE o.customer_id = c.customer_id
    AND o.order_date >= DATEADD(day, -30, CURRENT_DATE)
);

-- INSERT with NOT IN clause
INSERT INTO spectrum.inactive_products
SELECT product_id, product_name, category, last_modified_date
FROM products
WHERE product_id NOT IN (
    SELECT DISTINCT product_id
    FROM order_details
    WHERE order_date >= DATEADD(month, -6, CURRENT_DATE)
);

-- INSERT with multiple partition columns and complex expressions
INSERT INTO spectrum.hourly_metrics
SELECT server_id,
       COUNT(*) as request_count,
       AVG(response_time) as avg_response_time,
       MAX(response_time) as max_response_time,
       DATE_PART(year, request_timestamp),
       DATE_PART(month, request_timestamp),
       DATE_PART(day, request_timestamp),
       DATE_PART(hour, request_timestamp)
FROM server_logs
GROUP BY server_id,
         DATE_PART(year, request_timestamp),
         DATE_PART(month, request_timestamp),
         DATE_PART(day, request_timestamp),
         DATE_PART(hour, request_timestamp);