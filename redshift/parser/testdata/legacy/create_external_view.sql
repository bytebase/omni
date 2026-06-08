-- CREATE EXTERNAL VIEW test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_VIEW.html

-- Basic external view creation
CREATE EXTERNAL VIEW sample_schema.glue_data_catalog_view
AS SELECT * FROM sample_database.remote_table;

-- External view with IF NOT EXISTS
CREATE EXTERNAL VIEW sample_schema.conditional_view
IF NOT EXISTS
AS SELECT * FROM sample_database.remote_table;

-- Protected external view (from AWS documentation example)
CREATE EXTERNAL PROTECTED VIEW sample_schema.glue_data_catalog_view
IF NOT EXISTS
AS SELECT * FROM sample_database.remote_table;

-- External view with quoted table name (handling special characters)
CREATE EXTERNAL PROTECTED VIEW sample_schema.special_table_view
IF NOT EXISTS
AS SELECT * FROM sample_database."remote-table-name";

-- External view with awsdatacatalog prefix
CREATE EXTERNAL VIEW awsdatacatalog.analytics_db.customer_view
AS SELECT customer_id, customer_name, email FROM customer_data.customers;

-- External view with catalog_name.schema_name.view_name format
CREATE EXTERNAL VIEW my_catalog.analytics_schema.sales_view
AS SELECT * FROM sales_data.transactions;

-- External view with complex SELECT statement
CREATE EXTERNAL VIEW spectrum_schema.aggregated_sales_view
AS SELECT
    region,
    product_category,
    SUM(sales_amount) as total_sales,
    COUNT(*) as transaction_count,
    AVG(sales_amount) as avg_sale_amount
FROM sales_data.transactions
WHERE sale_date >= '2023-01-01'
GROUP BY region, product_category
HAVING SUM(sales_amount) > 1000
ORDER BY total_sales DESC;

-- External view with JOIN operations
CREATE EXTERNAL VIEW analytics_schema.customer_orders_view
AS SELECT
    c.customer_id,
    c.customer_name,
    c.email,
    o.order_id,
    o.order_date,
    o.total_amount
FROM customer_data.customers c
JOIN order_data.orders o ON c.customer_id = o.customer_id;

-- External view with multiple JOINs
CREATE EXTERNAL VIEW reporting_schema.detailed_order_view
AS SELECT
    c.customer_name,
    o.order_id,
    o.order_date,
    p.product_name,
    oi.quantity,
    oi.unit_price,
    (oi.quantity * oi.unit_price) as line_total
FROM customer_data.customers c
JOIN order_data.orders o ON c.customer_id = o.customer_id
JOIN order_data.order_items oi ON o.order_id = oi.order_id
JOIN product_data.products p ON oi.product_id = p.product_id;

-- External view with subquery
CREATE EXTERNAL VIEW analytics_schema.high_value_customers_view
AS SELECT
    customer_id,
    customer_name,
    email,
    total_orders,
    total_spent
FROM (
    SELECT
        c.customer_id,
        c.customer_name,
        c.email,
        COUNT(o.order_id) as total_orders,
        SUM(o.total_amount) as total_spent
    FROM customer_data.customers c
    LEFT JOIN order_data.orders o ON c.customer_id = o.customer_id
    GROUP BY c.customer_id, c.customer_name, c.email
) customer_stats
WHERE total_spent > 5000;

-- External view with CASE statements
CREATE EXTERNAL VIEW analytics_schema.customer_segmentation_view
AS SELECT
    customer_id,
    customer_name,
    total_spent,
    CASE
        WHEN total_spent >= 10000 THEN 'Premium'
        WHEN total_spent >= 5000 THEN 'Gold'
        WHEN total_spent >= 1000 THEN 'Silver'
        ELSE 'Bronze'
    END as customer_tier
FROM (
    SELECT
        c.customer_id,
        c.customer_name,
        COALESCE(SUM(o.total_amount), 0) as total_spent
    FROM customer_data.customers c
    LEFT JOIN order_data.orders o ON c.customer_id = o.customer_id
    GROUP BY c.customer_id, c.customer_name
);

-- External view with window functions
CREATE EXTERNAL VIEW analytics_schema.sales_ranking_view
AS SELECT
    employee_id,
    employee_name,
    department,
    sales_amount,
    RANK() OVER (PARTITION BY department ORDER BY sales_amount DESC) as dept_rank,
    ROW_NUMBER() OVER (ORDER BY sales_amount DESC) as overall_rank
FROM sales_data.employee_sales;

-- External view with date functions
CREATE EXTERNAL VIEW reporting_schema.monthly_sales_view
AS SELECT
    EXTRACT(YEAR FROM order_date) as sales_year,
    EXTRACT(MONTH FROM order_date) as sales_month,
    COUNT(*) as order_count,
    SUM(total_amount) as monthly_revenue
FROM order_data.orders
GROUP BY EXTRACT(YEAR FROM order_date), EXTRACT(MONTH FROM order_date)
ORDER BY sales_year, sales_month;

-- External view with string functions
CREATE EXTERNAL VIEW analytics_schema.customer_contact_view
AS SELECT
    customer_id,
    UPPER(customer_name) as customer_name_upper,
    LOWER(email) as email_lower,
    SUBSTRING(phone_number, 1, 3) as area_code,
    LENGTH(customer_name) as name_length
FROM customer_data.customers
WHERE email IS NOT NULL;

-- External view with UNION
CREATE EXTERNAL VIEW reporting_schema.all_transactions_view
AS SELECT
    'ONLINE' as channel,
    transaction_id,
    customer_id,
    amount,
    transaction_date
FROM online_sales.transactions
UNION ALL
SELECT
    'RETAIL' as channel,
    transaction_id,
    customer_id,
    amount,
    transaction_date
FROM retail_sales.transactions;

-- External view with EXISTS clause
CREATE EXTERNAL VIEW analytics_schema.customers_with_orders_view
AS SELECT
    customer_id,
    customer_name,
    email,
    registration_date
FROM customer_data.customers c
WHERE EXISTS (
    SELECT 1
    FROM order_data.orders o
    WHERE o.customer_id = c.customer_id
);

-- External view with NOT EXISTS clause
CREATE EXTERNAL VIEW analytics_schema.customers_without_orders_view
AS SELECT
    customer_id,
    customer_name,
    email,
    registration_date
FROM customer_data.customers c
WHERE NOT EXISTS (
    SELECT 1
    FROM order_data.orders o
    WHERE o.customer_id = c.customer_id
);

-- External view with IN clause
CREATE EXTERNAL VIEW analytics_schema.premium_product_sales_view
AS SELECT
    order_id,
    product_id,
    quantity,
    unit_price,
    (quantity * unit_price) as line_total
FROM order_data.order_items
WHERE product_id IN ('PROD001', 'PROD002', 'PROD003');

-- External view with BETWEEN clause
CREATE EXTERNAL VIEW reporting_schema.recent_orders_view
AS SELECT
    order_id,
    customer_id,
    order_date,
    total_amount,
    status
FROM order_data.orders
WHERE order_date BETWEEN '2023-01-01' AND '2023-12-31';

-- External view with LIKE pattern matching
CREATE EXTERNAL VIEW analytics_schema.tech_customers_view
AS SELECT
    customer_id,
    customer_name,
    email,
    company_name
FROM customer_data.customers
WHERE company_name LIKE '%Tech%' OR company_name LIKE '%Software%';

-- External view with NULL handling
CREATE EXTERNAL VIEW analytics_schema.complete_customer_profiles_view
AS SELECT
    customer_id,
    customer_name,
    COALESCE(email, 'No Email') as email_status,
    COALESCE(phone_number, 'No Phone') as phone_status,
    CASE
        WHEN email IS NOT NULL AND phone_number IS NOT NULL THEN 'Complete'
        WHEN email IS NOT NULL OR phone_number IS NOT NULL THEN 'Partial'
        ELSE 'Incomplete'
    END as profile_completeness
FROM customer_data.customers;

-- External view with mathematical operations
CREATE EXTERNAL VIEW analytics_schema.order_metrics_view
AS SELECT
    order_id,
    customer_id,
    subtotal,
    tax_rate,
    (subtotal * tax_rate / 100) as tax_amount,
    (subtotal + (subtotal * tax_rate / 100)) as total_amount,
    ROUND((subtotal * tax_rate / 100), 2) as rounded_tax
FROM order_data.orders
WHERE subtotal > 0;

-- External view with complex aggregations
CREATE EXTERNAL VIEW reporting_schema.product_performance_view
AS SELECT
    p.product_id,
    p.product_name,
    p.category,
    COUNT(oi.order_item_id) as times_ordered,
    SUM(oi.quantity) as total_quantity_sold,
    AVG(oi.unit_price) as avg_selling_price,
    MIN(oi.unit_price) as min_selling_price,
    MAX(oi.unit_price) as max_selling_price,
    SUM(oi.quantity * oi.unit_price) as total_revenue
FROM product_data.products p
LEFT JOIN order_data.order_items oi ON p.product_id = oi.product_id
GROUP BY p.product_id, p.product_name, p.category
ORDER BY total_revenue DESC;

-- External view with different schema naming patterns
CREATE EXTERNAL VIEW external_analytics.regional_sales_summary
AS SELECT
    region,
    COUNT(*) as total_orders,
    SUM(total_amount) as total_revenue
FROM sales_data.regional_orders
GROUP BY region;

-- External view with underscored naming
CREATE EXTERNAL VIEW data_warehouse.customer_lifetime_value
AS SELECT
    customer_id,
    first_order_date,
    last_order_date,
    total_orders,
    total_spent,
    (total_spent / total_orders) as avg_order_value
FROM customer_analytics.customer_summary;

-- External view with mixed case naming
CREATE EXTERNAL VIEW DataLake.CustomerAnalytics
AS SELECT
    CustomerId as customer_id,
    CustomerName as customer_name,
    TotalPurchases as total_purchases
FROM ExternalData.Customers;

-- Protected view with complex query
CREATE EXTERNAL PROTECTED VIEW secure_analytics.financial_summary
AS SELECT
    fiscal_year,
    quarter,
    revenue,
    expenses,
    (revenue - expenses) as profit,
    CASE
        WHEN (revenue - expenses) > 0 THEN 'Profitable'
        ELSE 'Loss'
    END as profit_status
FROM financial_data.quarterly_results
WHERE fiscal_year >= 2020;

-- External view with IF NOT EXISTS and complex naming
CREATE EXTERNAL VIEW cross_platform.unified_customer_data
IF NOT EXISTS
AS SELECT
    mobile_app.user_id as customer_id,
    mobile_app.app_usage_hours,
    web_platform.page_views,
    retail_pos.in_store_purchases
FROM mobile_data.app_usage mobile_app
FULL OUTER JOIN web_data.user_activity web_platform
    ON mobile_app.user_id = web_platform.user_id
FULL OUTER JOIN retail_data.pos_transactions retail_pos
    ON mobile_app.user_id = retail_pos.customer_id;

-- Simple view for testing basic functionality
CREATE EXTERNAL VIEW simple_view
AS SELECT id, name FROM basic_table;

-- View with fully qualified external schema reference
CREATE EXTERNAL VIEW spectrum_database.public.inventory_status
AS SELECT
    product_id,
    warehouse_location,
    quantity_on_hand,
    reorder_level,
    CASE
        WHEN quantity_on_hand <= reorder_level THEN 'REORDER'
        WHEN quantity_on_hand <= (reorder_level * 1.5) THEN 'LOW'
        ELSE 'NORMAL'
    END as stock_status
FROM inventory_data.warehouse_stock;