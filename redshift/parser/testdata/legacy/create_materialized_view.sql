-- CREATE MATERIALIZED VIEW test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/materialized-view-create-sql-command.html

-- Basic CREATE MATERIALIZED VIEW
CREATE MATERIALIZED VIEW tickets_mv AS
SELECT catgroup, SUM(qtysold) as sold
FROM category c, event e, sales s
WHERE c.catid = e.catid AND e.eventid = s.eventid
GROUP BY catgroup;

-- With schema qualification
CREATE MATERIALIZED VIEW sales.revenue_mv AS
SELECT product_id, SUM(amount) as total_revenue
FROM sales.transactions
GROUP BY product_id;

-- With quoted identifiers
CREATE MATERIALIZED VIEW "Customer-Summary" AS
SELECT "Customer ID", COUNT(*) as order_count
FROM "Order-History"
GROUP BY "Customer ID";

-- BACKUP clause examples
CREATE MATERIALIZED VIEW backup_mv BACKUP YES AS
SELECT customer_id, COUNT(*) as orders
FROM customer_orders
GROUP BY customer_id;

CREATE MATERIALIZED VIEW no_backup_mv BACKUP NO AS
SELECT product_id, AVG(rating) as avg_rating
FROM product_reviews
GROUP BY product_id;

-- DISTSTYLE examples
CREATE MATERIALIZED VIEW even_dist_mv DISTSTYLE EVEN AS
SELECT region, SUM(sales) as total_sales
FROM regional_sales
GROUP BY region;

CREATE MATERIALIZED VIEW all_dist_mv DISTSTYLE ALL AS
SELECT category, COUNT(*) as product_count
FROM products
GROUP BY category;

CREATE MATERIALIZED VIEW key_dist_mv DISTSTYLE KEY AS
SELECT customer_id, SUM(amount) as total_spent
FROM purchases
GROUP BY customer_id;

-- DISTKEY examples
CREATE MATERIALIZED VIEW customer_dist_mv DISTKEY(customer_id) AS
SELECT customer_id, order_date, amount
FROM orders;

CREATE MATERIALIZED VIEW product_dist_mv DISTKEY(product_id) AS
SELECT product_id, sale_date, quantity
FROM sales;

-- SORTKEY examples
CREATE MATERIALIZED VIEW date_sorted_mv SORTKEY(order_date) AS
SELECT order_date, customer_id, amount
FROM orders;

CREATE MATERIALIZED VIEW multi_sorted_mv SORTKEY(region, sale_date) AS
SELECT region, sale_date, product_id, quantity
FROM regional_sales;

CREATE MATERIALIZED VIEW compound_sorted_mv COMPOUND SORTKEY(customer_id, order_date) AS
SELECT customer_id, order_date, product_id, amount
FROM customer_orders;

CREATE MATERIALIZED VIEW interleaved_sorted_mv INTERLEAVED SORTKEY(product_id, sale_date) AS
SELECT product_id, sale_date, customer_id, quantity
FROM product_sales;

-- AUTO REFRESH examples
CREATE MATERIALIZED VIEW auto_refresh_mv AUTO REFRESH YES AS
SELECT category, AVG(price) as avg_price
FROM products
GROUP BY category;

CREATE MATERIALIZED VIEW manual_refresh_mv AUTO REFRESH NO AS
SELECT warehouse_id, SUM(inventory) as total_inventory
FROM warehouse_stock
GROUP BY warehouse_id;

-- Combined attributes examples
CREATE MATERIALIZED VIEW full_featured_mv
BACKUP YES
DISTSTYLE KEY
DISTKEY(customer_id)
SORTKEY(order_date, product_id)
AUTO REFRESH YES
AS
SELECT customer_id, order_date, product_id, SUM(amount) as total_amount
FROM order_details
GROUP BY customer_id, order_date, product_id;

CREATE MATERIALIZED VIEW complex_mv
BACKUP NO
DISTSTYLE ALL
COMPOUND SORTKEY(region, category)
AUTO REFRESH NO
AS
SELECT region, category, COUNT(*) as product_count, AVG(price) as avg_price
FROM product_catalog
GROUP BY region, category;

-- Cross-database reference (Spectrum/Federated)
CREATE MATERIALIZED VIEW cross_db_mv AS
SELECT cityname, population
FROM external_db.public.cities
WHERE population > 100000;

-- Complex SELECT statements
CREATE MATERIALIZED VIEW join_mv AS
SELECT c.customer_name, p.product_name, SUM(o.quantity) as total_quantity
FROM customers c
JOIN orders o ON c.customer_id = o.customer_id
JOIN products p ON o.product_id = p.product_id
GROUP BY c.customer_name, p.product_name;

CREATE MATERIALIZED VIEW union_mv AS
SELECT 'Q1' as quarter, product_id, SUM(sales) as total_sales
FROM q1_sales
GROUP BY product_id
UNION ALL
SELECT 'Q2' as quarter, product_id, SUM(sales) as total_sales
FROM q2_sales
GROUP BY product_id;

-- Window function example (if supported)
CREATE MATERIALIZED VIEW window_mv AS
SELECT customer_id, order_date, amount,
       SUM(amount) OVER (PARTITION BY customer_id ORDER BY order_date) as running_total
FROM orders;

-- CASE expression example
CREATE MATERIALIZED VIEW case_mv AS
SELECT product_id,
       CASE
         WHEN price < 100 THEN 'Low'
         WHEN price BETWEEN 100 AND 500 THEN 'Medium'
         ELSE 'High'
       END as price_category,
       COUNT(*) as product_count
FROM products
GROUP BY product_id, price_category;

-- Aggregate functions
CREATE MATERIALIZED VIEW agg_mv AS
SELECT category,
       COUNT(*) as total_products,
       SUM(price) as total_value,
       AVG(price) as avg_price,
       MIN(price) as min_price,
       MAX(price) as max_price
FROM products
GROUP BY category;

-- DISTINCT example
CREATE MATERIALIZED VIEW distinct_mv AS
SELECT DISTINCT customer_id, region
FROM customer_locations;

-- Subquery example
CREATE MATERIALIZED VIEW subquery_mv AS
SELECT customer_id, total_orders
FROM (
  SELECT customer_id, COUNT(*) as total_orders
  FROM orders
  GROUP BY customer_id
) sub
WHERE total_orders > 5;

-- Multiple JOIN types
CREATE MATERIALIZED VIEW multi_join_mv AS
SELECT c.customer_name, o.order_id, p.product_name, od.quantity
FROM customers c
LEFT JOIN orders o ON c.customer_id = o.customer_id
INNER JOIN order_details od ON o.order_id = od.order_id
RIGHT JOIN products p ON od.product_id = p.product_id;

-- WITH clause (CTE) - if supported
CREATE MATERIALIZED VIEW cte_mv AS
SELECT region, total_sales, avg_sales
FROM (
  SELECT region,
         SUM(amount) as total_sales,
         AVG(amount) as avg_sales
  FROM sales
  GROUP BY region
) regional_stats
WHERE total_sales > 10000;

-- Function calls
CREATE MATERIALIZED VIEW function_mv AS
SELECT
  EXTRACT(YEAR FROM order_date) as order_year,
  EXTRACT(MONTH FROM order_date) as order_month,
  COUNT(*) as order_count,
  SUM(amount) as total_amount
FROM orders
GROUP BY order_year, order_month;

-- NULL handling
CREATE MATERIALIZED VIEW null_handling_mv AS
SELECT
  customer_id,
  CASE WHEN email IS NULL THEN 'No Email' ELSE 'Has Email' END as email_status,
  COUNT(*) as customer_count
FROM customers
GROUP BY customer_id, email_status;

-- Complex WHERE conditions
CREATE MATERIALIZED VIEW complex_where_mv AS
SELECT product_id, category, price
FROM products
WHERE (price > 100 AND category IN ('Electronics', 'Computers'))
   OR (price BETWEEN 50 AND 100 AND category LIKE 'Home%')
   AND NOT (product_id IN (SELECT discontinued_id FROM discontinued_products));

-- Multiple schema qualification
CREATE MATERIALIZED VIEW schema_qualified_mv AS
SELECT s.sale_id, c.customer_name, p.product_name
FROM sales_db.public.sales s
JOIN customer_db.public.customers c ON s.customer_id = c.customer_id
JOIN product_db.public.products p ON s.product_id = p.product_id;