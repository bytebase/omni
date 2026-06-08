-- ALTER EXTERNAL VIEW test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_EXTERNAL_VIEW.html

-- Basic ALTER EXTERNAL VIEW with REMOVE DEFINITION
ALTER EXTERNAL VIEW sample_schema.glue_data_catalog_view REMOVE DEFINITION;

-- ALTER EXTERNAL VIEW with FORCE and REMOVE DEFINITION
ALTER EXTERNAL VIEW sample_schema.glue_data_catalog_view FORCE REMOVE DEFINITION;

-- ALTER EXTERNAL VIEW with new query definition
ALTER EXTERNAL VIEW external_schema.my_external_view AS (
    SELECT
        column1,
        column2,
        column3
    FROM external_schema.external_table
    WHERE condition = 'value'
);

-- ALTER EXTERNAL VIEW with catalog notation
ALTER EXTERNAL VIEW my_catalog.my_schema.my_view AS (
    SELECT * FROM my_catalog.my_schema.base_table
);

-- ALTER EXTERNAL VIEW with awsdatacatalog notation
ALTER EXTERNAL VIEW awsdatacatalog.mydb.customer_view AS (
    SELECT
        customer_id,
        customer_name,
        customer_email
    FROM awsdatacatalog.mydb.customers
);

-- ALTER EXTERNAL VIEW with FORCE and new complex query
ALTER EXTERNAL VIEW external_schema.sales_summary_view FORCE AS (
    SELECT
        region,
        product_category,
        SUM(sales_amount) as total_sales,
        COUNT(*) as transaction_count,
        AVG(sales_amount) as avg_sale
    FROM external_schema.sales_data
    WHERE sale_date >= '2023-01-01'
    GROUP BY region, product_category
    ORDER BY total_sales DESC
);

-- ALTER EXTERNAL VIEW with join query
ALTER EXTERNAL VIEW external_schema.customer_orders_view AS (
    SELECT
        c.customer_id,
        c.customer_name,
        o.order_id,
        o.order_date,
        o.total_amount
    FROM external_schema.customers c
    JOIN external_schema.orders o ON c.customer_id = o.customer_id
    WHERE o.order_status = 'completed'
);

-- ALTER EXTERNAL VIEW with subquery
ALTER EXTERNAL VIEW external_schema.top_products_view AS (
    SELECT
        product_id,
        product_name,
        total_revenue
    FROM (
        SELECT
            p.product_id,
            p.product_name,
            SUM(oi.quantity * oi.price) as total_revenue
        FROM external_schema.products p
        JOIN external_schema.order_items oi ON p.product_id = oi.product_id
        GROUP BY p.product_id, p.product_name
    ) AS product_revenues
    WHERE total_revenue > 10000
);

-- ALTER EXTERNAL VIEW with CTE
ALTER EXTERNAL VIEW external_schema.monthly_sales_view AS (
    WITH monthly_totals AS (
        SELECT
            DATE_TRUNC('month', sale_date) as sale_month,
            SUM(amount) as monthly_total
        FROM external_schema.sales
        GROUP BY DATE_TRUNC('month', sale_date)
    )
    SELECT
        sale_month,
        monthly_total,
        LAG(monthly_total) OVER (ORDER BY sale_month) as prev_month_total,
        monthly_total - LAG(monthly_total) OVER (ORDER BY sale_month) as month_over_month_change
    FROM monthly_totals
);

-- ALTER EXTERNAL VIEW with UNION
ALTER EXTERNAL VIEW external_schema.all_transactions_view AS (
    SELECT transaction_id, transaction_date, amount, 'online' as source
    FROM external_schema.online_transactions
    UNION ALL
    SELECT transaction_id, transaction_date, amount, 'store' as source
    FROM external_schema.store_transactions
);

-- ALTER EXTERNAL VIEW - remove definition for different schema notations
ALTER EXTERNAL VIEW my_external_schema.my_view REMOVE DEFINITION;
ALTER EXTERNAL VIEW my_catalog.my_schema.my_view REMOVE DEFINITION;
ALTER EXTERNAL VIEW awsdatacatalog.database_name.view_name REMOVE DEFINITION;

-- ALTER EXTERNAL VIEW with FORCE - remove definition for different views
ALTER EXTERNAL VIEW external_schema.stale_view FORCE REMOVE DEFINITION;
ALTER EXTERNAL VIEW catalog1.schema1.outdated_view FORCE REMOVE DEFINITION;

-- Complex example with multiple joins and aggregations
ALTER EXTERNAL VIEW external_schema.comprehensive_sales_report FORCE AS (
    SELECT
        c.customer_segment,
        p.product_category,
        r.region_name,
        DATE_TRUNC('quarter', s.sale_date) as sale_quarter,
        COUNT(DISTINCT s.transaction_id) as transaction_count,
        COUNT(DISTINCT c.customer_id) as unique_customers,
        SUM(s.quantity) as total_units_sold,
        SUM(s.amount) as total_revenue,
        AVG(s.amount) as avg_transaction_value,
        MAX(s.amount) as max_transaction_value
    FROM external_schema.sales s
    JOIN external_schema.customers c ON s.customer_id = c.customer_id
    JOIN external_schema.products p ON s.product_id = p.product_id
    JOIN external_schema.regions r ON s.region_id = r.region_id
    WHERE s.sale_date >= DATEADD(year, -2, CURRENT_DATE)
    GROUP BY
        c.customer_segment,
        p.product_category,
        r.region_name,
        DATE_TRUNC('quarter', s.sale_date)
    HAVING SUM(s.amount) > 1000
    ORDER BY sale_quarter DESC, total_revenue DESC
);