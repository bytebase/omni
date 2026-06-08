-- CREATE RLS POLICY test cases
-- Based on AWS Redshift documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_RLS_POLICY.html

-- Basic CREATE RLS POLICY with simple predicate
CREATE RLS POLICY policy_concerts
WITH (catgroup VARCHAR(10))
USING (catgroup = 'Concerts');

-- CREATE RLS POLICY with multiple column definitions
CREATE RLS POLICY policy_sales_region
WITH (region VARCHAR(20), year INTEGER)
USING (region = 'US' AND year >= 2020);

-- CREATE RLS POLICY with relation alias
CREATE RLS POLICY policy_user_access
WITH (user_id INTEGER) AS t
USING (user_id = current_user_id());

-- CREATE RLS POLICY with AS keyword explicit
CREATE RLS POLICY policy_department
WITH (dept_id INTEGER, dept_name VARCHAR(50)) AS dept_table
USING (dept_id IN (1, 2, 3));

-- CREATE RLS POLICY without WITH clause (minimal syntax)
CREATE RLS POLICY policy_simple
USING (status = 'active');

-- CREATE RLS POLICY with complex predicate using OR
CREATE RLS POLICY policy_multi_condition
WITH (category VARCHAR(30), price DECIMAL(10,2))
USING (category = 'Electronics' OR price < 100.00);

-- CREATE RLS POLICY with NULL checks
CREATE RLS POLICY policy_null_check
WITH (manager_id INTEGER)
USING (manager_id IS NOT NULL);

-- CREATE RLS POLICY with LIKE operator
CREATE RLS POLICY policy_name_pattern
WITH (product_name VARCHAR(100))
USING (product_name LIKE 'Premium%');

-- CREATE RLS POLICY with NOT conditions
CREATE RLS POLICY policy_exclude
WITH (status VARCHAR(20))
USING (status NOT IN ('deleted', 'archived'));

-- CREATE RLS POLICY with nested conditions
CREATE RLS POLICY policy_complex
WITH (user_role VARCHAR(20), department VARCHAR(30), level INTEGER)
USING ((user_role = 'admin' OR user_role = 'manager') AND department = 'finance' AND level >= 5);

-- CREATE RLS POLICY with date/timestamp columns
CREATE RLS POLICY policy_temporal
WITH (created_date DATE, updated_timestamp TIMESTAMP)
USING (created_date >= '2023-01-01' AND updated_timestamp IS NOT NULL);

-- CREATE RLS POLICY with different data types
CREATE RLS POLICY policy_datatypes
WITH (
    id BIGINT,
    name VARCHAR(255),
    active BOOLEAN,
    score REAL,
    precise_value DOUBLE PRECISION,
    created_time TIMESTAMPTZ
)
USING (active = TRUE AND score > 0.5);

-- CREATE RLS POLICY with CHAR and CHARACTER types
CREATE RLS POLICY policy_char_types
WITH (code CHAR(5), description CHARACTER(100))
USING (code = 'ADMIN');

-- CREATE RLS POLICY with TIME types
CREATE RLS POLICY policy_time_types
WITH (
    event_time TIME,
    event_timetz TIMETZ,
    log_timestamp TIMESTAMP WITH TIME ZONE,
    backup_time TIMESTAMP WITHOUT TIME ZONE
)
USING (event_time > '09:00:00');

-- CREATE RLS POLICY with function calls
CREATE RLS POLICY policy_functions
WITH (user_email VARCHAR(255))
USING (LOWER(user_email) = current_user());

-- CREATE RLS POLICY with quoted identifiers
CREATE RLS POLICY "Policy With Spaces"
WITH ("Column Name" VARCHAR(50), "Another Column" INTEGER)
USING ("Column Name" = 'test');

-- CREATE RLS POLICY with numeric literals and comparisons
CREATE RLS POLICY policy_numeric
WITH (quantity INTEGER, price NUMERIC(10,2))
USING (quantity >= 10 AND price <= 999.99);

-- CREATE RLS POLICY with boolean literals
CREATE RLS POLICY policy_boolean
WITH (is_active BOOLEAN, is_premium BOOLEAN)
USING (is_active = TRUE AND is_premium = FALSE);

-- CREATE RLS POLICY with scientific notation
CREATE RLS POLICY policy_scientific
WITH (measurement REAL)
USING (measurement > 1.5e-3);

-- CREATE RLS POLICY testing edge cases with parentheses
CREATE RLS POLICY policy_parentheses
WITH (a INTEGER, b INTEGER, c INTEGER)
USING ((a > 0 AND b > 0) OR (c IS NULL));

-- CREATE RLS POLICY with complex IN clause
CREATE RLS POLICY policy_in_clause
WITH (category VARCHAR(50))
USING (category IN ('Electronics', 'Books', 'Clothing', 'Home & Garden'));

-- CREATE RLS POLICY with string concatenation function
CREATE RLS POLICY policy_concat
WITH (first_name VARCHAR(50), last_name VARCHAR(50))
USING (CONCAT(first_name, ' ', last_name) = 'John Doe');