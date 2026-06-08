-- DELETE command examples from AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DELETE.html

-- Example 1: Delete entire table (from AWS docs)
DELETE FROM category;

-- Example 2: Delete with optional FROM keyword
DELETE category;

-- Example 3: Delete rows with specific condition (from AWS docs)
DELETE FROM category WHERE catid BETWEEN 0 AND 9;

-- Example 4: Delete rows based on subquery (from AWS docs)
DELETE FROM listing
WHERE listing.sellerid NOT IN (SELECT sales.sellerid FROM sales);

-- Example 5: Delete using join condition (from AWS docs)
DELETE FROM category
USING event
WHERE event.catid = category.catid AND category.catid = 9;

-- Example 6: Delete from materialized view (from AWS docs)
DELETE FROM mv_cities;

-- Example 7: Delete with simple WHERE clause
DELETE FROM users WHERE age > 65;

-- Example 8: Delete with multiple conditions
DELETE FROM orders WHERE status = 'cancelled' AND order_date < '2023-01-01';

-- Example 9: Delete with USING multiple tables
DELETE FROM inventory
USING products, categories
WHERE inventory.product_id = products.id
  AND products.category_id = categories.id
  AND categories.name = 'discontinued';

-- Example 10: Delete with schema-qualified table names
DELETE FROM sales.orders WHERE customer_id IS NULL;

-- Example 11: Delete with complex WHERE conditions
DELETE FROM customer_data
WHERE last_login < CURRENT_DATE - INTERVAL '2 years'
  AND account_status = 'inactive';

-- Example 12: Delete with IN clause
DELETE FROM temp_table WHERE id IN (1, 2, 3, 4, 5);

-- Example 13: Delete with NOT IN subquery
DELETE FROM products WHERE category_id NOT IN (
    SELECT id FROM active_categories WHERE status = 'active'
);

-- Example 14: Delete with EXISTS subquery
DELETE FROM users WHERE EXISTS (
    SELECT 1 FROM banned_users WHERE banned_users.user_id = users.id
);

-- Example 15: Delete with NOT EXISTS
DELETE FROM products WHERE NOT EXISTS (
    SELECT 1 FROM sales WHERE sales.product_id = products.id
);

-- Example 16: Delete with LIKE pattern matching
DELETE FROM log_entries WHERE message LIKE '%error%' OR message LIKE '%warning%';

-- Example 17: Delete with BETWEEN clause
DELETE FROM transactions WHERE amount BETWEEN 0.01 AND 10.00;

-- Example 18: Delete with IS NULL condition
DELETE FROM customer_profiles WHERE email IS NULL;

-- Example 19: Delete with IS NOT NULL condition
DELETE FROM temp_data WHERE processed_date IS NOT NULL;

-- Example 20: Delete with mathematical expressions
DELETE FROM inventory WHERE quantity * unit_price < 100;

-- Example 21: Delete with function calls in WHERE clause
DELETE FROM users WHERE EXTRACT(YEAR FROM created_date) < 2020;

-- Example 22: Delete with case-sensitive string comparison
DELETE FROM products WHERE UPPER(name) = 'DISCONTINUED ITEM';

-- Example 23: Delete with date comparison
DELETE FROM sessions WHERE session_date < '2023-01-01'::date;

-- Example 24: Delete with timestamp comparison
DELETE FROM audit_log WHERE created_at < '2023-01-01 00:00:00'::timestamp;

-- Example 25: Delete with complex join using USING
DELETE FROM order_items
USING orders, customers
WHERE order_items.order_id = orders.id
  AND orders.customer_id = customers.id
  AND customers.status = 'deleted';

-- Example 26: Delete with self-referencing subquery
DELETE FROM employees
WHERE department_id IN (
    SELECT id FROM departments WHERE budget < 50000
);

-- Example 27: Delete with nested subqueries
DELETE FROM products
WHERE category_id IN (
    SELECT category_id FROM product_categories
    WHERE parent_id IN (SELECT id FROM categories WHERE active = false)
);

-- Example 28: Delete with OR conditions
DELETE FROM temp_results
WHERE status = 'error'
   OR status = 'timeout'
   OR created_date < CURRENT_DATE - 7;

-- Example 29: Delete with AND/OR combination
DELETE FROM user_sessions
WHERE (expired = true OR last_activity < CURRENT_DATE - 1)
  AND session_type = 'guest';

-- Example 30: Delete with quoted identifiers
DELETE FROM "Special Table" WHERE "Column Name" = 'value';

-- Example 31: Delete with table alias in USING
DELETE FROM a
USING products b, categories c
WHERE a.product_id = b.id
  AND b.category_id = c.id
  AND c.discontinued = true;

-- Example 32: Delete with ESCAPE clause
DELETE FROM search_terms WHERE term LIKE '%\_%' ESCAPE '\';

-- Example 33: Delete with multiple mathematical operations
DELETE FROM financial_data
WHERE (revenue - expenses) / revenue < 0.1;

-- Example 34: Delete with string concatenation
DELETE FROM logs WHERE CONCAT(level, ': ', message) LIKE '%ERROR:%';

-- Example 35: Delete with window function results (using subquery)
DELETE FROM sales
WHERE sale_id IN (
    SELECT sale_id FROM (
        SELECT sale_id, ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY sale_date DESC) as rn
        FROM sales
    ) ranked
    WHERE rn > 10
);

-- Example 36: Delete with complex boolean logic
DELETE FROM user_preferences
WHERE NOT (
    (notification_email = true AND email IS NOT NULL)
    OR
    (notification_sms = true AND phone IS NOT NULL)
);

-- Example 37: Delete with type casting
DELETE FROM measurements WHERE reading::NUMERIC > 100.5;

-- Example 38: Delete from schema-qualified materialized view
DELETE FROM analytics.customer_summary WHERE last_updated < '2023-01-01';

-- Example 39: Delete with DISTINCT in subquery
DELETE FROM duplicates
WHERE id NOT IN (
    SELECT DISTINCT MIN(id) FROM duplicates GROUP BY email
);

-- Example 40: Delete with ALL in subquery comparison
DELETE FROM products
WHERE price > ALL (
    SELECT price FROM competitor_products WHERE category = 'electronics'
);

-- Example 41: Delete with ANY in subquery comparison
DELETE FROM inventory
WHERE quantity < ANY (
    SELECT min_threshold FROM reorder_rules WHERE product_type = 'perishable'
);

-- Example 42: Delete using CTE (Common Table Expression)
WITH inactive_users AS (
    SELECT user_id FROM user_activity
    WHERE last_login < CURRENT_DATE - INTERVAL '1 year'
)
DELETE FROM user_profiles
WHERE user_id IN (SELECT user_id FROM inactive_users);

-- Example 43: Delete using recursive CTE
WITH RECURSIVE category_tree AS (
    SELECT id, parent_id, name FROM categories WHERE parent_id IS NULL
    UNION ALL
    SELECT c.id, c.parent_id, c.name
    FROM categories c
    JOIN category_tree ct ON c.parent_id = ct.id
)
DELETE FROM products
WHERE category_id IN (
    SELECT id FROM category_tree WHERE name = 'obsolete'
);

-- Example 44: Delete with multiple CTEs
WITH
    expired_sessions AS (
        SELECT session_id FROM sessions WHERE expires_at < CURRENT_TIMESTAMP
    ),
    inactive_users AS (
        SELECT user_id FROM users WHERE last_activity < CURRENT_DATE - 30
    )
DELETE FROM user_sessions
WHERE session_id IN (SELECT session_id FROM expired_sessions)
   OR user_id IN (SELECT user_id FROM inactive_users);

-- Example 45: Delete with nested parentheses in conditions
DELETE FROM complex_data
WHERE ((status = 'pending' AND priority < 3) OR (status = 'failed' AND retry_count > 5))
  AND (created_date BETWEEN '2023-01-01' AND '2023-12-31');

-- Example 46: Delete with function calls and expressions
DELETE FROM metrics
WHERE ABS(value - target_value) > (target_value * 0.2);

-- Example 47: Delete with CASE-like logic using subqueries
DELETE FROM orders
WHERE status = 'pending'
  AND order_date < (
      CASE
          WHEN priority = 'high' THEN CURRENT_DATE - 1
          WHEN priority = 'normal' THEN CURRENT_DATE - 7
          ELSE CURRENT_DATE - 30
      END
  );

-- Example 48: Delete with JSON operations (if supported)
DELETE FROM user_data WHERE preferences::JSON->>'theme' = 'dark';

-- Example 49: Delete with array operations
DELETE FROM tags WHERE 'deprecated' = ANY(tag_array);

-- Example 50: Delete with complex multi-table USING clause
DELETE FROM order_line_items
USING orders o, customers c, products p
WHERE order_line_items.order_id = o.id
  AND o.customer_id = c.id
  AND order_line_items.product_id = p.id
  AND c.status = 'inactive'
  AND p.discontinued = true
  AND o.order_date < '2022-01-01';