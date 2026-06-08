-- MERGE command examples from AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_MERGE.html

-- Example 1: Basic MERGE with UPDATE and INSERT (from AWS docs)
MERGE INTO target_table
USING source_table
ON target_table.id = source_table.id
WHEN MATCHED THEN UPDATE SET name = source_table.name, value = source_table.value
WHEN NOT MATCHED THEN INSERT VALUES (source_table.id, source_table.name, source_table.value);

-- Example 2: MERGE with DELETE action (from AWS docs)
MERGE INTO customer_table
USING temp_customer_data
ON customer_table.customer_id = temp_customer_data.customer_id
WHEN MATCHED THEN DELETE
WHEN NOT MATCHED THEN INSERT VALUES (temp_customer_data.customer_id, temp_customer_data.name, temp_customer_data.email);

-- Example 3: MERGE with REMOVE DUPLICATES (from AWS docs)
MERGE INTO sales_data
USING new_sales_data
ON sales_data.sale_id = new_sales_data.sale_id
REMOVE DUPLICATES;

-- Example 4: MERGE with table aliases
MERGE INTO inventory i
USING product_updates p
ON i.product_id = p.product_id
WHEN MATCHED THEN UPDATE SET quantity = p.quantity, price = p.price
WHEN NOT MATCHED THEN INSERT VALUES (p.product_id, p.quantity, p.price);

-- Example 5: MERGE with schema-qualified table names
MERGE INTO sales.customer_data
USING staging.new_customers
ON sales.customer_data.customer_id = staging.new_customers.customer_id
WHEN MATCHED THEN UPDATE SET
    name = staging.new_customers.name,
    email = staging.new_customers.email,
    last_updated = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN INSERT VALUES (
    staging.new_customers.customer_id,
    staging.new_customers.name,
    staging.new_customers.email,
    CURRENT_TIMESTAMP
);

-- Example 6: MERGE with subquery as source
MERGE INTO product_summary
USING (
    SELECT product_id, SUM(quantity) as total_qty, AVG(price) as avg_price
    FROM product_transactions
    WHERE transaction_date >= '2023-01-01'
    GROUP BY product_id
) t
ON product_summary.product_id = t.product_id
WHEN MATCHED THEN UPDATE SET total_quantity = t.total_qty, average_price = t.avg_price
WHEN NOT MATCHED THEN INSERT VALUES (t.product_id, t.total_qty, t.avg_price);

-- Example 7: MERGE with complex match condition
MERGE INTO employee_data
USING hr_updates
ON employee_data.emp_id = hr_updates.emp_id AND employee_data.department = hr_updates.department
WHEN MATCHED THEN UPDATE SET
    salary = hr_updates.salary,
    title = hr_updates.title,
    manager_id = hr_updates.manager_id
WHEN NOT MATCHED THEN INSERT VALUES (
    hr_updates.emp_id,
    hr_updates.department,
    hr_updates.salary,
    hr_updates.title,
    hr_updates.manager_id
);

-- Example 8: MERGE with multiple SET clauses
MERGE INTO user_profile
USING user_updates
ON user_profile.user_id = user_updates.user_id
WHEN MATCHED THEN UPDATE SET
    first_name = user_updates.first_name,
    last_name = user_updates.last_name,
    email = user_updates.email,
    phone = user_updates.phone,
    address = user_updates.address,
    last_modified = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN INSERT VALUES (
    user_updates.user_id,
    user_updates.first_name,
    user_updates.last_name,
    user_updates.email,
    user_updates.phone,
    user_updates.address,
    CURRENT_TIMESTAMP
);

-- Example 9: MERGE with WHERE clause in subquery source
MERGE INTO active_customers
USING (
    SELECT customer_id, name, email, phone
    FROM customer_staging
    WHERE status = 'active' AND created_date >= '2023-01-01'
) active_data
ON active_customers.customer_id = active_data.customer_id
WHEN MATCHED THEN UPDATE SET
    name = active_data.name,
    email = active_data.email,
    phone = active_data.phone
WHEN NOT MATCHED THEN INSERT VALUES (
    active_data.customer_id,
    active_data.name,
    active_data.email,
    active_data.phone
);

-- Example 10: MERGE with functions in expressions
MERGE INTO customer_metrics
USING daily_activity
ON customer_metrics.customer_id = daily_activity.customer_id
WHEN MATCHED THEN UPDATE SET
    total_orders = customer_metrics.total_orders + daily_activity.order_count,
    last_activity = GREATEST(customer_metrics.last_activity, daily_activity.activity_date),
    average_order_value = (customer_metrics.total_revenue + daily_activity.revenue) /
                         (customer_metrics.total_orders + daily_activity.order_count)
WHEN NOT MATCHED THEN INSERT VALUES (
    daily_activity.customer_id,
    daily_activity.order_count,
    daily_activity.activity_date,
    daily_activity.revenue / daily_activity.order_count
);

-- Example 11: MERGE with quoted identifiers
MERGE INTO "Customer Data"
USING "New Customer Info" AS nci
ON "Customer Data"."Customer ID" = nci."Customer ID"
WHEN MATCHED THEN UPDATE SET
    "Customer Name" = nci."Customer Name",
    "Email Address" = nci."Email Address"
WHEN NOT MATCHED THEN INSERT VALUES (
    nci."Customer ID",
    nci."Customer Name",
    nci."Email Address"
);

-- Example 12: MERGE with CASE expressions
MERGE INTO product_pricing
USING market_data
ON product_pricing.product_id = market_data.product_id
WHEN MATCHED THEN UPDATE SET
    price = CASE
        WHEN market_data.competitor_price > product_pricing.price THEN market_data.competitor_price * 0.95
        WHEN market_data.competitor_price < product_pricing.price THEN market_data.competitor_price * 1.05
        ELSE product_pricing.price
    END,
    last_updated = CURRENT_DATE
WHEN NOT MATCHED THEN INSERT VALUES (
    market_data.product_id,
    market_data.competitor_price,
    CURRENT_DATE
);

-- Example 13: MERGE with NULL handling
MERGE INTO contact_info
USING contact_updates
ON contact_info.contact_id = contact_updates.contact_id
WHEN MATCHED THEN UPDATE SET
    email = COALESCE(contact_updates.email, contact_info.email),
    phone = COALESCE(contact_updates.phone, contact_info.phone),
    address = COALESCE(contact_updates.address, contact_info.address)
WHEN NOT MATCHED THEN INSERT VALUES (
    contact_updates.contact_id,
    contact_updates.email,
    contact_updates.phone,
    contact_updates.address
);

-- Example 14: MERGE with string functions
MERGE INTO user_accounts
USING user_imports
ON UPPER(user_accounts.username) = UPPER(user_imports.username)
WHEN MATCHED THEN UPDATE SET
    email = LOWER(user_imports.email),
    full_name = INITCAP(user_imports.full_name),
    last_login = user_imports.last_login
WHEN NOT MATCHED THEN INSERT VALUES (
    LOWER(user_imports.username),
    LOWER(user_imports.email),
    INITCAP(user_imports.full_name),
    user_imports.last_login
);

-- Example 15: MERGE with date/time functions
MERGE INTO sales_summary
USING daily_sales
ON sales_summary.sale_date = daily_sales.sale_date AND
   sales_summary.region = daily_sales.region
WHEN MATCHED THEN UPDATE SET
    total_amount = sales_summary.total_amount + daily_sales.amount,
    transaction_count = sales_summary.transaction_count + daily_sales.count,
    last_updated = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN INSERT VALUES (
    daily_sales.sale_date,
    daily_sales.region,
    daily_sales.amount,
    daily_sales.count,
    CURRENT_TIMESTAMP
);

-- Example 16: MERGE with mathematical operations
MERGE INTO inventory_levels
USING inventory_adjustments
ON inventory_levels.item_id = inventory_adjustments.item_id
WHEN MATCHED THEN UPDATE SET
    quantity = inventory_levels.quantity + inventory_adjustments.adjustment,
    value = (inventory_levels.quantity + inventory_adjustments.adjustment) * inventory_levels.unit_cost,
    last_adjustment = CURRENT_DATE
WHEN NOT MATCHED THEN INSERT VALUES (
    inventory_adjustments.item_id,
    inventory_adjustments.adjustment,
    inventory_adjustments.adjustment * inventory_adjustments.unit_cost,
    inventory_adjustments.unit_cost,
    CURRENT_DATE
);

-- Example 17: MERGE with subquery in match condition
MERGE INTO premium_customers
USING customer_analysis
ON premium_customers.customer_id = customer_analysis.customer_id
WHEN MATCHED THEN UPDATE SET
    tier_level = customer_analysis.calculated_tier,
    annual_revenue = customer_analysis.total_revenue
WHEN NOT MATCHED AND customer_analysis.total_revenue > (
    SELECT AVG(total_revenue) * 2 FROM customer_analysis
) THEN INSERT VALUES (
    customer_analysis.customer_id,
    customer_analysis.calculated_tier,
    customer_analysis.total_revenue
);

-- Example 18: MERGE with aggregate functions in source
MERGE INTO department_stats
USING (
    SELECT
        department_id,
        COUNT(*) as employee_count,
        AVG(salary) as avg_salary,
        MAX(hire_date) as latest_hire,
        MIN(hire_date) as earliest_hire
    FROM employees
    GROUP BY department_id
) dept_agg
ON department_stats.department_id = dept_agg.department_id
WHEN MATCHED THEN UPDATE SET
    total_employees = dept_agg.employee_count,
    average_salary = dept_agg.avg_salary,
    newest_hire = dept_agg.latest_hire,
    oldest_hire = dept_agg.earliest_hire
WHEN NOT MATCHED THEN INSERT VALUES (
    dept_agg.department_id,
    dept_agg.employee_count,
    dept_agg.avg_salary,
    dept_agg.latest_hire,
    dept_agg.earliest_hire
);

-- Example 19: MERGE with JOIN in source subquery
MERGE INTO customer_orders_summary
USING (
    SELECT
        c.customer_id,
        c.customer_name,
        COUNT(o.order_id) as order_count,
        SUM(o.total_amount) as total_spent,
        MAX(o.order_date) as last_order_date
    FROM customers c
    LEFT JOIN orders o ON c.customer_id = o.customer_id
    WHERE c.status = 'active'
    GROUP BY c.customer_id, c.customer_name
) summary_data
ON customer_orders_summary.customer_id = summary_data.customer_id
WHEN MATCHED THEN UPDATE SET
    customer_name = summary_data.customer_name,
    total_orders = summary_data.order_count,
    lifetime_value = summary_data.total_spent,
    last_purchase = summary_data.last_order_date
WHEN NOT MATCHED THEN INSERT VALUES (
    summary_data.customer_id,
    summary_data.customer_name,
    summary_data.order_count,
    summary_data.total_spent,
    summary_data.last_order_date
);

-- Example 20: MERGE with complex boolean expressions
MERGE INTO product_alerts
USING inventory_check
ON product_alerts.product_id = inventory_check.product_id
WHEN MATCHED AND (inventory_check.current_stock < inventory_check.reorder_point
                  OR inventory_check.days_until_expiry < 30) THEN UPDATE SET
    alert_type = CASE
        WHEN inventory_check.current_stock < inventory_check.reorder_point AND inventory_check.days_until_expiry < 30
        THEN 'LOW_STOCK_AND_EXPIRING'
        WHEN inventory_check.current_stock < inventory_check.reorder_point
        THEN 'LOW_STOCK'
        ELSE 'EXPIRING_SOON'
    END,
    alert_date = CURRENT_DATE,
    current_quantity = inventory_check.current_stock
WHEN NOT MATCHED AND (inventory_check.current_stock < inventory_check.reorder_point
                      OR inventory_check.days_until_expiry < 30) THEN INSERT VALUES (
    inventory_check.product_id,
    CASE
        WHEN inventory_check.current_stock < inventory_check.reorder_point AND inventory_check.days_until_expiry < 30
        THEN 'LOW_STOCK_AND_EXPIRING'
        WHEN inventory_check.current_stock < inventory_check.reorder_point
        THEN 'LOW_STOCK'
        ELSE 'EXPIRING_SOON'
    END,
    CURRENT_DATE,
    inventory_check.current_stock
);

-- Example 21: MERGE with window functions in source (via subquery)
MERGE INTO sales_rankings
USING (
    SELECT
        salesperson_id,
        sales_amount,
        RANK() OVER (ORDER BY sales_amount DESC) as sales_rank,
        PERCENT_RANK() OVER (ORDER BY sales_amount DESC) as percentile_rank
    FROM (
        SELECT salesperson_id, SUM(amount) as sales_amount
        FROM sales_transactions
        WHERE transaction_date >= DATE_TRUNC('month', CURRENT_DATE)
        GROUP BY salesperson_id
    ) monthly_sales
) rankings
ON sales_rankings.salesperson_id = rankings.salesperson_id
WHEN MATCHED THEN UPDATE SET
    monthly_sales = rankings.sales_amount,
    rank_position = rankings.sales_rank,
    performance_percentile = rankings.percentile_rank,
    last_updated = CURRENT_DATE
WHEN NOT MATCHED THEN INSERT VALUES (
    rankings.salesperson_id,
    rankings.sales_amount,
    rankings.sales_rank,
    rankings.performance_percentile,
    CURRENT_DATE
);

-- Example 22: MERGE with type casting
MERGE INTO financial_metrics
USING raw_financial_data
ON financial_metrics.account_id = raw_financial_data.account_id::INTEGER
WHEN MATCHED THEN UPDATE SET
    balance = raw_financial_data.balance_text::DECIMAL(15,2),
    interest_rate = raw_financial_data.rate_text::NUMERIC(5,4),
    last_transaction = raw_financial_data.transaction_date::DATE
WHEN NOT MATCHED THEN INSERT VALUES (
    raw_financial_data.account_id::INTEGER,
    raw_financial_data.balance_text::DECIMAL(15,2),
    raw_financial_data.rate_text::NUMERIC(5,4),
    raw_financial_data.transaction_date::DATE
);

-- Example 23: MERGE with BETWEEN conditions
MERGE INTO age_demographics
USING customer_ages
ON age_demographics.age_group =
    CASE
        WHEN customer_ages.age BETWEEN 18 AND 25 THEN '18-25'
        WHEN customer_ages.age BETWEEN 26 AND 35 THEN '26-35'
        WHEN customer_ages.age BETWEEN 36 AND 50 THEN '36-50'
        WHEN customer_ages.age BETWEEN 51 AND 65 THEN '51-65'
        ELSE '65+'
    END
WHEN MATCHED THEN UPDATE SET
    customer_count = age_demographics.customer_count + 1,
    total_revenue = age_demographics.total_revenue + customer_ages.annual_spending
WHEN NOT MATCHED THEN INSERT VALUES (
    CASE
        WHEN customer_ages.age BETWEEN 18 AND 25 THEN '18-25'
        WHEN customer_ages.age BETWEEN 26 AND 35 THEN '26-35'
        WHEN customer_ages.age BETWEEN 36 AND 50 THEN '36-50'
        WHEN customer_ages.age BETWEEN 51 AND 65 THEN '51-65'
        ELSE '65+'
    END,
    1,
    customer_ages.annual_spending
);

-- Example 24: MERGE with IN clause
MERGE INTO priority_customers
USING customer_data
ON priority_customers.customer_id = customer_data.customer_id
WHEN MATCHED AND customer_data.tier IN ('gold', 'platinum', 'diamond') THEN UPDATE SET
    tier_level = customer_data.tier,
    priority_score = CASE customer_data.tier
        WHEN 'diamond' THEN 100
        WHEN 'platinum' THEN 75
        WHEN 'gold' THEN 50
        ELSE 25
    END,
    last_updated = CURRENT_TIMESTAMP
WHEN NOT MATCHED AND customer_data.tier IN ('gold', 'platinum', 'diamond') THEN INSERT VALUES (
    customer_data.customer_id,
    customer_data.tier,
    CASE customer_data.tier
        WHEN 'diamond' THEN 100
        WHEN 'platinum' THEN 75
        WHEN 'gold' THEN 50
        ELSE 25
    END,
    CURRENT_TIMESTAMP
);

-- Example 25: MERGE with LIKE pattern matching
MERGE INTO email_domains
USING (
    SELECT DISTINCT
        SUBSTRING(email FROM '@(.+)$') as domain,
        COUNT(*) as user_count
    FROM user_registrations
    WHERE email LIKE '%@%'
    GROUP BY SUBSTRING(email FROM '@(.+)$')
) domain_stats
ON email_domains.domain_name = domain_stats.domain
WHEN MATCHED THEN UPDATE SET
    registration_count = email_domains.registration_count + domain_stats.user_count,
    last_updated = CURRENT_DATE
WHEN NOT MATCHED THEN INSERT VALUES (
    domain_stats.domain,
    domain_stats.user_count,
    CURRENT_DATE
);

-- Example 26: MERGE with EXISTS subquery
MERGE INTO active_products
USING product_catalog
ON active_products.product_id = product_catalog.product_id
WHEN MATCHED AND EXISTS (
    SELECT 1 FROM recent_sales
    WHERE recent_sales.product_id = product_catalog.product_id
    AND sale_date >= CURRENT_DATE - 90
) THEN UPDATE SET
    status = 'active',
    last_sale_date = (
        SELECT MAX(sale_date) FROM recent_sales
        WHERE recent_sales.product_id = product_catalog.product_id
    ),
    total_sales = (
        SELECT COUNT(*) FROM recent_sales
        WHERE recent_sales.product_id = product_catalog.product_id
        AND sale_date >= CURRENT_DATE - 90
    )
WHEN NOT MATCHED AND EXISTS (
    SELECT 1 FROM recent_sales
    WHERE recent_sales.product_id = product_catalog.product_id
    AND sale_date >= CURRENT_DATE - 90
) THEN INSERT VALUES (
    product_catalog.product_id,
    'active',
    (SELECT MAX(sale_date) FROM recent_sales
     WHERE recent_sales.product_id = product_catalog.product_id),
    (SELECT COUNT(*) FROM recent_sales
     WHERE recent_sales.product_id = product_catalog.product_id
     AND sale_date >= CURRENT_DATE - 90)
);

-- Example 27: MERGE with NOT EXISTS
MERGE INTO inactive_customers
USING all_customers
ON inactive_customers.customer_id = all_customers.customer_id
WHEN MATCHED AND NOT EXISTS (
    SELECT 1 FROM recent_orders
    WHERE recent_orders.customer_id = all_customers.customer_id
    AND order_date >= CURRENT_DATE - 365
) THEN UPDATE SET
    days_inactive = CURRENT_DATE - (
        SELECT MAX(order_date) FROM customer_orders
        WHERE customer_orders.customer_id = all_customers.customer_id
    ),
    status = 'inactive'
WHEN NOT MATCHED AND NOT EXISTS (
    SELECT 1 FROM recent_orders
    WHERE recent_orders.customer_id = all_customers.customer_id
    AND order_date >= CURRENT_DATE - 365
) THEN INSERT VALUES (
    all_customers.customer_id,
    all_customers.customer_name,
    CURRENT_DATE - (
        SELECT MAX(order_date) FROM customer_orders
        WHERE customer_orders.customer_id = all_customers.customer_id
    ),
    'inactive'
);

-- Example 28: MERGE with IS NULL conditions
MERGE INTO customer_contacts
USING customer_updates
ON customer_contacts.customer_id = customer_updates.customer_id
WHEN MATCHED THEN UPDATE SET
    email = CASE WHEN customer_updates.email IS NOT NULL THEN customer_updates.email ELSE customer_contacts.email END,
    phone = CASE WHEN customer_updates.phone IS NOT NULL THEN customer_updates.phone ELSE customer_contacts.phone END,
    address = CASE WHEN customer_updates.address IS NOT NULL THEN customer_updates.address ELSE customer_contacts.address END,
    updated_fields = (
        CASE WHEN customer_updates.email IS NOT NULL THEN 'email,' ELSE '' END ||
        CASE WHEN customer_updates.phone IS NOT NULL THEN 'phone,' ELSE '' END ||
        CASE WHEN customer_updates.address IS NOT NULL THEN 'address,' ELSE '' END
    )
WHEN NOT MATCHED THEN INSERT VALUES (
    customer_updates.customer_id,
    customer_updates.email,
    customer_updates.phone,
    customer_updates.address,
    'new_record'
);

-- Example 29: MERGE with complex JOIN in source
MERGE INTO order_fulfillment_stats
USING (
    SELECT
        o.order_id,
        o.customer_id,
        o.order_date,
        s.shipment_date,
        COALESCE(s.shipment_date - o.order_date, 0) as processing_days,
        CASE WHEN s.shipment_date IS NULL THEN 'pending' ELSE 'shipped' END as status
    FROM orders o
    LEFT JOIN shipments s ON o.order_id = s.order_id
    WHERE o.order_date >= CURRENT_DATE - 30
) order_data
ON order_fulfillment_stats.order_id = order_data.order_id
WHEN MATCHED THEN UPDATE SET
    shipment_date = order_data.shipment_date,
    processing_time = order_data.processing_days,
    fulfillment_status = order_data.status,
    last_updated = CURRENT_TIMESTAMP
WHEN NOT MATCHED THEN INSERT VALUES (
    order_data.order_id,
    order_data.customer_id,
    order_data.order_date,
    order_data.shipment_date,
    order_data.processing_days,
    order_data.status,
    CURRENT_TIMESTAMP
);

-- Example 30: MERGE with only DELETE action (cleanup operation)
MERGE INTO temp_processing_data
USING cleanup_list
ON temp_processing_data.record_id = cleanup_list.record_id
WHEN MATCHED THEN DELETE;

-- Example 31: MERGE with only INSERT action (new records only)
MERGE INTO new_user_queue
USING user_registrations
ON new_user_queue.user_id = user_registrations.user_id
WHEN NOT MATCHED THEN INSERT VALUES (
    user_registrations.user_id,
    user_registrations.email,
    user_registrations.registration_date,
    'pending_verification'
);

-- Example 32: Simple REMOVE DUPLICATES with minimal syntax
MERGE INTO sales_transactions
USING sales_staging
ON sales_transactions.transaction_id = sales_staging.transaction_id
REMOVE DUPLICATES;

-- Example 33: MERGE with all three actions (UPDATE, DELETE, INSERT)
MERGE INTO product_inventory
USING product_updates
ON product_inventory.product_id = product_updates.product_id
WHEN MATCHED AND product_updates.action_type = 'delete' THEN DELETE
WHEN MATCHED AND product_updates.action_type = 'update' THEN UPDATE SET
    quantity = product_updates.quantity,
    price = product_updates.price,
    last_updated = CURRENT_TIMESTAMP
WHEN NOT MATCHED AND product_updates.action_type = 'insert' THEN INSERT VALUES (
    product_updates.product_id,
    product_updates.quantity,
    product_updates.price,
    CURRENT_TIMESTAMP
);

-- Example 34: MERGE with complex expressions in VALUES
MERGE INTO calculated_metrics
USING base_data
ON calculated_metrics.metric_id = base_data.metric_id
WHEN MATCHED THEN UPDATE SET
    calculated_value = base_data.raw_value * base_data.multiplier + base_data.offset,
    confidence_score = CASE
        WHEN base_data.sample_size >= 1000 THEN 0.95
        WHEN base_data.sample_size >= 100 THEN 0.85
        ELSE 0.70
    END
WHEN NOT MATCHED THEN INSERT VALUES (
    base_data.metric_id,
    base_data.raw_value * base_data.multiplier + base_data.offset,
    CASE
        WHEN base_data.sample_size >= 1000 THEN 0.95
        WHEN base_data.sample_size >= 100 THEN 0.85
        ELSE 0.70
    END,
    CURRENT_DATE
);

-- Example 35: MERGE with ORDER BY in source subquery
MERGE INTO top_performers
USING (
    SELECT employee_id, performance_score, performance_rank
    FROM (
        SELECT
            employee_id,
            performance_score,
            ROW_NUMBER() OVER (ORDER BY performance_score DESC) as performance_rank
        FROM employee_evaluations
        WHERE evaluation_date >= DATE_TRUNC('quarter', CURRENT_DATE)
    ) ranked_employees
    WHERE performance_rank <= 10
    ORDER BY performance_rank
) top_10
ON top_performers.employee_id = top_10.employee_id
WHEN MATCHED THEN UPDATE SET
    quarterly_score = top_10.performance_score,
    rank_position = top_10.performance_rank,
    quarter_year = DATE_TRUNC('quarter', CURRENT_DATE)
WHEN NOT MATCHED THEN INSERT VALUES (
    top_10.employee_id,
    top_10.performance_score,
    top_10.performance_rank,
    DATE_TRUNC('quarter', CURRENT_DATE)
);