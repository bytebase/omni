-- CREATE TABLE AS test cases for Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_TABLE_AS.html

-- ====================================
-- Basic CREATE TABLE AS (CTAS)
-- ====================================

-- Simple CTAS with SELECT statement
CREATE TABLE customer_summary AS
SELECT customer_id, COUNT(*) as order_count, SUM(total_amount) as total_spent
FROM orders
GROUP BY customer_id;

-- CTAS with column list specification
CREATE TABLE monthly_sales (year_month, total_sales, order_count) AS
SELECT
    DATE_TRUNC('month', order_date) as year_month,
    SUM(amount) as total_sales,
    COUNT(*) as order_count
FROM sales
GROUP BY DATE_TRUNC('month', order_date);

CREATE TABLE daily_aggregates AS
SELECT
    order_date,
    COUNT(*) as daily_orders,
    AVG(total_amount) as avg_order_value,
    SUM(total_amount) as daily_revenue
FROM orders
GROUP BY order_date;

-- ====================================
-- Temporary tables with CTAS
-- ====================================

-- Temporary table with CTAS
CREATE TEMPORARY TABLE temp_user_stats AS
SELECT
    user_id,
    COUNT(*) as login_count,
    MAX(last_login) as last_activity
FROM user_sessions
GROUP BY user_id;

-- Local temporary table with CTAS
CREATE LOCAL TEMPORARY TABLE local_temp_analysis AS
SELECT
    product_id,
    AVG(rating) as avg_rating,
    COUNT(review_id) as review_count
FROM product_reviews
WHERE created_date >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY product_id;

-- Short form temporary table
CREATE TEMP TABLE quick_summary AS
SELECT category, COUNT(*) as product_count
FROM products
WHERE status = 'active'
GROUP BY category;

-- Table name with # prefix (automatically temporary)
CREATE TABLE #session_data AS
SELECT session_id, user_id, duration_minutes
FROM user_sessions
WHERE session_date = CURRENT_DATE;

-- ====================================
-- BACKUP option
-- ====================================

-- CTAS with BACKUP YES
CREATE TABLE important_customer_data
BACKUP YES
AS
SELECT
    customer_id,
    registration_date,
    total_purchases,
    loyalty_tier
FROM customers
WHERE total_purchases > 1000;

-- CTAS with BACKUP NO
CREATE TABLE temp_calculations
BACKUP NO
AS
SELECT
    product_id,
    cost_price * 1.3 as suggested_price,
    inventory_count * cost_price as inventory_value
FROM product_inventory;

-- ====================================
-- Distribution styles (DISTSTYLE)
-- ====================================

-- DISTSTYLE ALL - replicate to all nodes
CREATE TABLE lookup_categories
DISTSTYLE ALL
AS
SELECT category_id, category_name, parent_category
FROM product_categories
WHERE is_active = true;

-- DISTSTYLE EVEN - distribute evenly
CREATE TABLE log_analysis
DISTSTYLE EVEN
AS
SELECT
    DATE_TRUNC('hour', timestamp) as hour_bucket,
    COUNT(*) as event_count,
    COUNT(DISTINCT user_id) as unique_users
FROM event_logs
GROUP BY DATE_TRUNC('hour', timestamp);

-- DISTSTYLE KEY with DISTKEY
CREATE TABLE customer_orders
DISTSTYLE KEY
DISTKEY(customer_id)
AS
SELECT
    customer_id,
    order_id,
    order_date,
    total_amount,
    order_status
FROM orders
WHERE order_date >= '2023-01-01';

-- DISTSTYLE AUTO - let Redshift choose
CREATE TABLE smart_distribution
DISTSTYLE AUTO
AS
SELECT
    region,
    sale_date,
    SUM(amount) as regional_sales,
    COUNT(*) as transaction_count
FROM sales_transactions
GROUP BY region, sale_date;

-- ====================================
-- Sort keys (SORTKEY)
-- ====================================

-- COMPOUND SORTKEY (default)
CREATE TABLE time_series_data
SORTKEY(timestamp, sensor_id)
AS
SELECT
    timestamp,
    sensor_id,
    temperature,
    humidity,
    pressure
FROM sensor_readings
WHERE timestamp >= CURRENT_DATE - INTERVAL '7 days';

-- Explicit COMPOUND SORTKEY
CREATE TABLE sales_analysis
COMPOUND SORTKEY(sale_date, region)
AS
SELECT
    sale_date,
    region,
    product_category,
    SUM(amount) as daily_sales,
    COUNT(*) as transaction_count
FROM sales_data
GROUP BY sale_date, region, product_category;

-- INTERLEAVED SORTKEY
CREATE TABLE customer_events
INTERLEAVED SORTKEY(customer_id, event_date)
AS
SELECT
    customer_id,
    event_date,
    event_type,
    session_id,
    page_url
FROM web_events
WHERE event_date >= '2023-01-01';

-- ====================================
-- Multiple table attributes
-- ====================================

-- Combining DISTSTYLE, DISTKEY, SORTKEY, and BACKUP
CREATE TABLE comprehensive_sales
BACKUP YES
DISTSTYLE KEY
DISTKEY(customer_id)
SORTKEY(order_date, order_id)
AS
SELECT
    customer_id,
    order_id,
    order_date,
    product_id,
    quantity,
    unit_price,
    total_amount,
    region
FROM order_details od
JOIN orders o ON od.order_id = o.order_id
WHERE o.order_date >= '2023-01-01';

-- Complex table with INTERLEAVED SORTKEY and DISTSTYLE AUTO
CREATE TABLE user_activity_summary
BACKUP NO
DISTSTYLE AUTO
INTERLEAVED SORTKEY(user_id, activity_date)
AS
SELECT
    user_id,
    activity_date,
    page_views,
    session_count,
    total_time_spent,
    purchases_made
FROM (
    SELECT
        user_id,
        DATE(activity_timestamp) as activity_date,
        COUNT(DISTINCT page_id) as page_views,
        COUNT(DISTINCT session_id) as session_count,
        SUM(time_on_page) as total_time_spent,
        COUNT(purchase_id) as purchases_made
    FROM user_activity_log
    WHERE activity_timestamp >= CURRENT_DATE - INTERVAL '90 days'
    GROUP BY user_id, DATE(activity_timestamp)
) subquery;

-- ====================================
-- ENCODE AUTO
-- ====================================

-- Combining ENCODE AUTO with other attributes
CREATE TABLE optimized_table
BACKUP YES
DISTSTYLE KEY
DISTKEY(region)
SORTKEY(created_date)
AS
SELECT
    region,
    customer_segment,
    created_date,
    revenue,
    customer_count
FROM regional_summary;

-- ====================================
-- Complex SELECT statements
-- ====================================

-- CTAS with JOINs
CREATE TABLE customer_order_summary AS
SELECT
    c.customer_id,
    c.customer_name,
    c.email,
    COUNT(o.order_id) as total_orders,
    SUM(o.total_amount) as lifetime_value,
    AVG(o.total_amount) as avg_order_value,
    MIN(o.order_date) as first_order,
    MAX(o.order_date) as last_order
FROM customers c
LEFT JOIN orders o ON c.customer_id = o.customer_id
GROUP BY c.customer_id, c.customer_name, c.email;

-- CTAS with multiple JOINs and WHERE clause
CREATE TABLE product_performance AS
SELECT
    p.product_id,
    p.product_name,
    p.category,
    COUNT(oi.order_item_id) as times_ordered,
    SUM(oi.quantity) as total_quantity_sold,
    SUM(oi.quantity * oi.unit_price) as total_revenue,
    AVG(pr.rating) as avg_rating,
    COUNT(pr.review_id) as review_count
FROM products p
LEFT JOIN order_items oi ON p.product_id = oi.product_id
LEFT JOIN product_reviews pr ON p.product_id = pr.product_id
WHERE p.status = 'active'
GROUP BY p.product_id, p.product_name, p.category
HAVING COUNT(oi.order_item_id) > 0;

-- CTAS with window functions
CREATE TABLE sales_rankings AS
SELECT
    salesperson_id,
    sales_month,
    monthly_sales,
    RANK() OVER (PARTITION BY sales_month ORDER BY monthly_sales DESC) as monthly_rank,
    LAG(monthly_sales, 1) OVER (PARTITION BY salesperson_id ORDER BY sales_month) as prev_month_sales,
    SUM(monthly_sales) OVER (PARTITION BY salesperson_id ORDER BY sales_month ROWS UNBOUNDED PRECEDING) as running_total
FROM (
    SELECT
        salesperson_id,
        DATE_TRUNC('month', sale_date) as sales_month,
        SUM(amount) as monthly_sales
    FROM sales
    GROUP BY salesperson_id, DATE_TRUNC('month', sale_date)
) monthly_data;

-- CTAS with CTE (Common Table Expression)
CREATE TABLE customer_segmentation AS
WITH customer_stats AS (
    SELECT
        customer_id,
        COUNT(*) as order_count,
        SUM(total_amount) as total_spent,
        AVG(total_amount) as avg_order_value,
        MAX(order_date) as last_order_date,
        MIN(order_date) as first_order_date
    FROM orders
    WHERE order_date >= '2022-01-01'
    GROUP BY customer_id
),
customer_segments AS (
    SELECT
        customer_id,
        order_count,
        total_spent,
        avg_order_value,
        DATEDIFF(day, last_order_date, CURRENT_DATE) as days_since_last_order,
        CASE
            WHEN total_spent >= 10000 THEN 'VIP'
            WHEN total_spent >= 5000 THEN 'Premium'
            WHEN total_spent >= 1000 THEN 'Standard'
            ELSE 'Basic'
        END as spending_tier,
        CASE
            WHEN days_since_last_order <= 30 THEN 'Active'
            WHEN days_since_last_order <= 90 THEN 'At Risk'
            ELSE 'Inactive'
        END as activity_status
    FROM customer_stats
)
SELECT * FROM customer_segments;

-- ====================================
-- UNION and set operations
-- ====================================

-- CTAS with UNION
CREATE TABLE all_transactions AS
SELECT
    transaction_id,
    customer_id,
    'PURCHASE' as transaction_type,
    amount,
    transaction_date
FROM purchases
WHERE transaction_date >= '2023-01-01'
UNION ALL
SELECT
    refund_id as transaction_id,
    customer_id,
    'REFUND' as transaction_type,
    -amount as amount,
    refund_date as transaction_date
FROM refunds
WHERE refund_date >= '2023-01-01';

-- CTAS with INTERSECT
CREATE TABLE common_customers AS
SELECT customer_id, customer_name
FROM online_customers
INTERSECT
SELECT customer_id, customer_name
FROM store_customers;

-- ====================================
-- Subqueries and correlated subqueries
-- ====================================

-- CTAS with subquery in SELECT
CREATE TABLE customer_insights AS
SELECT
    c.customer_id,
    c.customer_name,
    (SELECT COUNT(*) FROM orders o WHERE o.customer_id = c.customer_id) as order_count,
    (SELECT MAX(order_date) FROM orders o WHERE o.customer_id = c.customer_id) as last_order_date,
    (SELECT AVG(total_amount) FROM orders o WHERE o.customer_id = c.customer_id) as avg_order_value
FROM customers c
WHERE c.registration_date >= '2022-01-01';

-- CTAS with EXISTS subquery
CREATE TABLE active_products AS
SELECT
    product_id,
    product_name,
    price,
    category
FROM products p
WHERE EXISTS (
    SELECT 1
    FROM order_items oi
    WHERE oi.product_id = p.product_id
    AND oi.order_date >= CURRENT_DATE - INTERVAL '6 months'
);

-- ====================================
-- Aggregate functions and grouping
-- ====================================

-- CTAS with advanced aggregations
CREATE TABLE sales_metrics AS
SELECT
    region,
    product_category,
    sale_month,
    COUNT(*) as transaction_count,
    COUNT(DISTINCT customer_id) as unique_customers,
    SUM(amount) as total_sales,
    AVG(amount) as avg_transaction,
    MEDIAN(amount) as median_transaction,
    STDDEV(amount) as sales_volatility,
    PERCENTILE_CONT(0.75) WITHIN GROUP (ORDER BY amount) as p75_amount,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY amount) as p95_amount
FROM sales_data
WHERE sale_date >= '2023-01-01'
GROUP BY region, product_category, DATE_TRUNC('month', sale_date);

-- CTAS with ROLLUP
CREATE TABLE sales_hierarchy AS
SELECT
    region,
    state,
    city,
    SUM(sales_amount) as total_sales,
    COUNT(*) as transaction_count
FROM regional_sales
GROUP BY ROLLUP(region, state, city);

-- CTAS with CUBE
CREATE TABLE multi_dimensional_sales AS
SELECT
    product_category,
    customer_segment,
    sales_channel,
    SUM(revenue) as total_revenue,
    COUNT(*) as sale_count
FROM sales_fact
GROUP BY CUBE(product_category, customer_segment, sales_channel);

-- ====================================
-- Date and time functions
-- ====================================

-- CTAS with date manipulation
CREATE TABLE time_based_analysis AS
SELECT
    EXTRACT(YEAR FROM order_date) as order_year,
    EXTRACT(QUARTER FROM order_date) as order_quarter,
    EXTRACT(MONTH FROM order_date) as order_month,
    EXTRACT(DAYOFWEEK FROM order_date) as day_of_week,
    DATE_TRUNC('week', order_date) as week_start,
    COUNT(*) as order_count,
    SUM(total_amount) as weekly_revenue
FROM orders
WHERE order_date >= '2022-01-01'
GROUP BY
    EXTRACT(YEAR FROM order_date),
    EXTRACT(QUARTER FROM order_date),
    EXTRACT(MONTH FROM order_date),
    EXTRACT(DAYOFWEEK FROM order_date),
    DATE_TRUNC('week', order_date);

-- ====================================
-- JSON and SUPER data type operations
-- ====================================

-- CTAS with JSON operations (Redshift SUPER type)
CREATE TABLE user_preferences AS
SELECT
    user_id,
    preferences.email_notifications as email_prefs,
    preferences.language as preferred_language,
    preferences.theme as ui_theme,
    JSON_EXTRACT_PATH_TEXT(preferences, 'notification_frequency') as notification_freq
FROM user_profiles
WHERE preferences IS NOT NULL;

-- ====================================
-- String and text functions
-- ====================================

-- CTAS with string manipulation
CREATE TABLE processed_names AS
SELECT
    customer_id,
    UPPER(first_name) as first_name_upper,
    LOWER(last_name) as last_name_lower,
    INITCAP(TRIM(first_name || ' ' || last_name)) as full_name_proper,
    LENGTH(email) as email_length,
    SPLIT_PART(email, '@', 2) as email_domain,
    REGEXP_REPLACE(phone, '[^0-9]', '') as phone_digits_only
FROM customers
WHERE email IS NOT NULL;

-- ====================================
-- Mathematical and statistical functions
-- ====================================

-- CTAS with mathematical operations
CREATE TABLE financial_calculations AS
SELECT
    account_id,
    balance,
    interest_rate,
    ROUND(balance * interest_rate / 100, 2) as annual_interest,
    POWER(1 + interest_rate/100, 12) as compound_factor,
    LN(balance) as log_balance,
    CASE
        WHEN balance > 100000 THEN 'High'
        WHEN balance > 10000 THEN 'Medium'
        ELSE 'Low'
    END as balance_tier
FROM account_balances
WHERE balance > 0;

-- ====================================
-- Qualified table names
-- ====================================

-- CTAS with schema-qualified source tables
CREATE TABLE public.customer_analysis AS
SELECT
    customer_id,
    order_count,
    total_revenue
FROM analytics.customer_metrics
WHERE last_updated >= CURRENT_DATE - INTERVAL '1 day';

-- CTAS with database.schema.table naming
CREATE TABLE tickit.public.event_summary AS
SELECT
    event_id,
    event_name,
    venue_name,
    event_date,
    total_tickets,
    tickets_sold
FROM tickit.public.events e
JOIN tickit.public.venues v ON e.venue_id = v.venue_id
WHERE event_date >= CURRENT_DATE;

-- ====================================
-- Complex real-world examples
-- ====================================

-- E-commerce customer lifetime value calculation
CREATE TABLE customer_ltv_analysis
BACKUP YES
DISTSTYLE KEY
DISTKEY(customer_id)
SORTKEY(registration_date, customer_id)
AS
WITH customer_orders AS (
    SELECT
        customer_id,
        MIN(order_date) as first_order_date,
        MAX(order_date) as last_order_date,
        COUNT(*) as total_orders,
        SUM(total_amount) as total_spent,
        AVG(total_amount) as avg_order_value,
        AVG(DATEDIFF(day, LAG(order_date) OVER (PARTITION BY customer_id ORDER BY order_date), order_date)) as avg_days_between_orders
    FROM orders
    GROUP BY customer_id
),
customer_demographics AS (
    SELECT
        customer_id,
        registration_date,
        age,
        gender,
        city,
        state,
        country
    FROM customers
)
SELECT
    cd.customer_id,
    cd.registration_date,
    cd.age,
    cd.gender,
    cd.city,
    cd.state,
    co.total_orders,
    co.total_spent,
    co.avg_order_value,
    co.avg_days_between_orders,
    DATEDIFF(day, co.first_order_date, co.last_order_date) as customer_lifespan_days,
    CASE
        WHEN co.last_order_date >= CURRENT_DATE - INTERVAL '30 days' THEN 'Active'
        WHEN co.last_order_date >= CURRENT_DATE - INTERVAL '90 days' THEN 'At Risk'
        ELSE 'Churned'
    END as customer_status,
    co.total_spent / NULLIF(DATEDIFF(day, co.first_order_date, CURRENT_DATE), 0) * 365 as annualized_ltv
FROM customer_demographics cd
LEFT JOIN customer_orders co ON cd.customer_id = co.customer_id;

-- Financial risk assessment table
CREATE TABLE risk_assessment
BACKUP YES
DISTSTYLE EVEN
COMPOUND SORTKEY(assessment_date, risk_score)
AS
SELECT
    customer_id,
    assessment_date,
    credit_score,
    debt_to_income_ratio,
    payment_history_score,
    account_age_months,
    (credit_score * 0.4 +
     payment_history_score * 0.3 +
     (100 - debt_to_income_ratio) * 0.2 +
     LEAST(account_age_months / 12 * 10, 100) * 0.1) as risk_score,
    CASE
        WHEN credit_score >= 750 AND debt_to_income_ratio <= 30 THEN 'Low Risk'
        WHEN credit_score >= 650 AND debt_to_income_ratio <= 50 THEN 'Medium Risk'
        ELSE 'High Risk'
    END as risk_category
FROM credit_assessments
WHERE assessment_date >= CURRENT_DATE - INTERVAL '1 year';

-- ====================================
-- Error cases and edge cases
-- ====================================

-- CTAS with empty result set (valid but creates empty table)
CREATE TABLE empty_table AS
SELECT customer_id, order_date, total_amount
FROM orders
WHERE 1 = 0;

-- CTAS with NULL handling
CREATE TABLE null_safe_aggregates AS
SELECT
    region,
    COUNT(*) as total_records,
    COUNT(customer_id) as non_null_customers,
    COUNT(revenue) as non_null_revenue,
    COALESCE(SUM(revenue), 0) as total_revenue,
    COALESCE(AVG(revenue), 0) as avg_revenue,
    MAX(order_date) as latest_order
FROM sales_data
GROUP BY region;

-- CTAS with type casting
CREATE TABLE typed_conversions AS
SELECT
    customer_id,
    CAST(total_amount AS INTEGER) as amount_int,
    CAST(order_date AS VARCHAR(10)) as order_date_str,
    order_status::VARCHAR(20) as status_varchar,
    CASE WHEN is_premium = 'true' THEN 1 ELSE 0 END as premium_flag
FROM orders
WHERE total_amount IS NOT NULL;