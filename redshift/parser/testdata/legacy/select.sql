-- Basic SELECT statements
-- Select all columns
SELECT * FROM users;

-- Select specific columns
SELECT username, email, created_date FROM users;

-- Select with column aliases
SELECT
    username AS user_name,
    email AS user_email,
    created_date AS signup_date
FROM users;

-- Select with expressions
SELECT
    first_name || ' ' || last_name AS full_name,
    age * 12 AS age_in_months,
    salary / 12 AS monthly_salary
FROM employees;

-- TOP clause
SELECT TOP 10 * FROM sales ORDER BY sale_amount DESC;
SELECT TOP 100 product_id, SUM(quantity) FROM sales GROUP BY product_id;

-- DISTINCT
SELECT DISTINCT category FROM products;
SELECT DISTINCT department, location FROM offices;

-- WHERE clause
SELECT * FROM orders WHERE order_date >= '2025-01-01';
SELECT * FROM products WHERE price BETWEEN 10 AND 100;
SELECT * FROM customers WHERE country IN ('US', 'CA', 'MX');
SELECT * FROM users WHERE email LIKE '%@example.com';
SELECT * FROM employees WHERE manager_id IS NOT NULL;

-- Joins
-- INNER JOIN
SELECT o.order_id, o.order_date, c.customer_name
FROM orders o
INNER JOIN customers c ON o.customer_id = c.customer_id;

-- LEFT JOIN
SELECT c.customer_name, o.order_id
FROM customers c
LEFT JOIN orders o ON c.customer_id = o.customer_id;

-- RIGHT JOIN
SELECT o.order_id, c.customer_name
FROM orders o
RIGHT JOIN customers c ON o.customer_id = c.customer_id;

-- FULL OUTER JOIN
SELECT c.customer_name, o.order_id
FROM customers c
FULL OUTER JOIN orders o ON c.customer_id = o.customer_id;

-- CROSS JOIN
SELECT p.product_name, c.category_name
FROM products p
CROSS JOIN categories c;

-- Multiple joins
SELECT
    o.order_id,
    c.customer_name,
    p.product_name,
    oi.quantity,
    oi.unit_price
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
JOIN order_items oi ON o.order_id = oi.order_id
JOIN products p ON oi.product_id = p.product_id;

-- Self join
SELECT
    e1.employee_name AS employee,
    e2.employee_name AS manager
FROM employees e1
LEFT JOIN employees e2 ON e1.manager_id = e2.employee_id;

-- GROUP BY and aggregates
SELECT category, COUNT(*) AS product_count
FROM products
GROUP BY category;

SELECT
    department,
    AVG(salary) AS avg_salary,
    MIN(salary) AS min_salary,
    MAX(salary) AS max_salary,
    SUM(salary) AS total_salary,
    COUNT(*) AS employee_count
FROM employees
GROUP BY department;

-- GROUP BY with HAVING
SELECT category, COUNT(*) AS product_count
FROM products
GROUP BY category
HAVING COUNT(*) > 10;

SELECT department, AVG(salary) AS avg_salary
FROM employees
GROUP BY department
HAVING AVG(salary) > 50000;

-- ROLLUP
SELECT
    region,
    country,
    SUM(sales_amount) AS total_sales
FROM sales_data
GROUP BY ROLLUP(region, country);

-- CUBE
SELECT
    product_category,
    store_location,
    SUM(quantity) AS total_quantity
FROM inventory
GROUP BY CUBE(product_category, store_location);

-- GROUPING SETS
SELECT
    department,
    job_title,
    COUNT(*) AS employee_count
FROM employees
GROUP BY GROUPING SETS (
    (department),
    (job_title),
    (department, job_title),
    ()
);

-- Common Table Expressions (WITH clause)
WITH regional_sales AS (
    SELECT region, SUM(sales_amount) AS total_sales
    FROM orders
    GROUP BY region
)
SELECT * FROM regional_sales WHERE total_sales > 1000000;

-- Multiple CTEs
WITH
customer_orders AS (
    SELECT customer_id, COUNT(*) AS order_count
    FROM orders
    GROUP BY customer_id
),
high_value_customers AS (
    SELECT customer_id
    FROM customer_orders
    WHERE order_count > 10
)
SELECT c.*
FROM customers c
JOIN high_value_customers hvc ON c.customer_id = hvc.customer_id;

-- Recursive CTE
WITH RECURSIVE employee_hierarchy AS (
    -- Base case: top-level employees
    SELECT
        employee_id,
        employee_name,
        manager_id,
        1 AS level
    FROM employees
    WHERE manager_id IS NULL

    UNION ALL

    -- Recursive case: employees with managers
    SELECT
        e.employee_id,
        e.employee_name,
        e.manager_id,
        eh.level + 1
    FROM employees e
    JOIN employee_hierarchy eh ON e.manager_id = eh.employee_id
    WHERE eh.level < 5
)
SELECT * FROM employee_hierarchy ORDER BY level, employee_name;

-- Window functions
SELECT
    employee_name,
    department,
    salary,
    AVG(salary) OVER (PARTITION BY department) AS dept_avg_salary,
    salary - AVG(salary) OVER (PARTITION BY department) AS salary_diff
FROM employees;

SELECT
    order_date,
    order_amount,
    SUM(order_amount) OVER (ORDER BY order_date) AS running_total,
    ROW_NUMBER() OVER (ORDER BY order_amount DESC) AS order_rank
FROM orders;

-- Window function with frame
SELECT
    sale_date,
    sale_amount,
    AVG(sale_amount) OVER (
        ORDER BY sale_date
        ROWS BETWEEN 2 PRECEDING AND 2 FOLLOWING
    ) AS moving_avg
FROM daily_sales;

-- QUALIFY clause
SELECT
    product_id,
    sale_date,
    sale_amount,
    ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY sale_amount DESC) AS rn
FROM sales
QUALIFY rn <= 3;

-- Set operations
-- UNION
SELECT customer_id, customer_name FROM customers_us
UNION
SELECT customer_id, customer_name FROM customers_ca;

-- UNION ALL
SELECT product_id FROM orders_2024
UNION ALL
SELECT product_id FROM orders_2025;

-- INTERSECT
SELECT customer_id FROM customers
INTERSECT
SELECT customer_id FROM orders;

-- EXCEPT
SELECT customer_id FROM customers
EXCEPT
SELECT customer_id FROM orders;

-- MINUS (synonym for EXCEPT)
SELECT product_id FROM all_products
MINUS
SELECT product_id FROM discontinued_products;

-- ORDER BY
SELECT * FROM products ORDER BY price DESC;
SELECT * FROM employees ORDER BY department, salary DESC;
SELECT * FROM orders ORDER BY order_date DESC NULLS LAST;

-- LIMIT and OFFSET
SELECT * FROM products LIMIT 10;
SELECT * FROM products LIMIT 10 OFFSET 20;
SELECT * FROM products ORDER BY price DESC LIMIT ALL;

-- Subqueries
-- Scalar subquery
SELECT
    product_name,
    price,
    (SELECT AVG(price) FROM products) AS avg_price
FROM products;

-- IN subquery
SELECT * FROM customers
WHERE customer_id IN (
    SELECT DISTINCT customer_id
    FROM orders
    WHERE order_date >= '2025-01-01'
);

-- EXISTS subquery
SELECT c.* FROM customers c
WHERE EXISTS (
    SELECT 1
    FROM orders o
    WHERE o.customer_id = c.customer_id
    AND o.order_amount > 1000
);

-- Correlated subquery
SELECT
    e1.employee_name,
    e1.salary,
    e1.department
FROM employees e1
WHERE e1.salary > (
    SELECT AVG(e2.salary)
    FROM employees e2
    WHERE e2.department = e1.department
);

-- CASE expressions
SELECT
    order_id,
    order_amount,
    CASE
        WHEN order_amount >= 1000 THEN 'High'
        WHEN order_amount >= 500 THEN 'Medium'
        ELSE 'Low'
    END AS order_category
FROM orders;

SELECT
    product_name,
    CASE product_category
        WHEN 'Electronics' THEN price * 0.9
        WHEN 'Clothing' THEN price * 0.8
        ELSE price
    END AS discounted_price
FROM products;

-- LISTAGG
SELECT
    department,
    LISTAGG(employee_name, ', ') WITHIN GROUP (ORDER BY employee_name) AS employees
FROM employees
GROUP BY department;

SELECT
    customer_id,
    LISTAGG(DISTINCT product_category, ' | ') AS purchased_categories
FROM customer_purchases
GROUP BY customer_id;

-- Complex query with multiple features
WITH monthly_sales AS (
    SELECT
        DATE_TRUNC('month', sale_date) AS month,
        product_category,
        SUM(sale_amount) AS total_sales,
        COUNT(*) AS transaction_count
    FROM sales
    WHERE sale_date >= '2024-01-01'
    GROUP BY DATE_TRUNC('month', sale_date), product_category
),
category_ranks AS (
    SELECT
        month,
        product_category,
        total_sales,
        transaction_count,
        RANK() OVER (PARTITION BY month ORDER BY total_sales DESC) AS sales_rank,
        total_sales / SUM(total_sales) OVER (PARTITION BY month) * 100 AS pct_of_month
    FROM monthly_sales
)
SELECT
    TO_CHAR(month, 'YYYY-MM') AS month_str,
    product_category,
    total_sales,
    transaction_count,
    sales_rank,
    ROUND(pct_of_month, 2) AS pct_of_month
FROM category_ranks
WHERE sales_rank <= 5
ORDER BY month, sales_rank;

-- EXCLUDE clause (Redshift specific)
SELECT * EXCLUDE (password, ssn) FROM users;
SELECT * EXCLUDE (internal_notes, debug_info) FROM orders;

-- TABLESAMPLE
SELECT * FROM large_table TABLESAMPLE BERNOULLI (10);
SELECT * FROM huge_table TABLESAMPLE SYSTEM (1) WHERE created_date >= '2025-01-01';

-- FOR UPDATE
SELECT * FROM inventory WHERE product_id = 123 FOR UPDATE;
SELECT * FROM accounts WHERE status = 'pending' FOR UPDATE SKIP LOCKED;

-- -- CONNECT BY (hierarchical queries)
-- SELECT employee_id, employee_name, manager_id, LEVEL
-- FROM employees
-- START WITH manager_id IS NULL
-- CONNECT BY PRIOR employee_id = manager_id;
--
-- -- Complex hierarchical query with CONNECT BY
-- SELECT
--     LPAD(' ', (LEVEL - 1) * 2) || employee_name AS org_chart,
--     employee_id,
--     manager_id,
--     LEVEL AS hierarchy_level
-- FROM employees
--          START WITH employee_id = 1
-- CONNECT BY PRIOR employee_id = manager_id
-- ORDER BY employee_name;