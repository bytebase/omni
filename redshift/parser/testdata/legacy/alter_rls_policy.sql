-- ALTER RLS POLICY test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_RLS_POLICY.html

-- Note: ALTER RLS POLICY in Redshift only supports modifying the USING clause.
-- RENAME TO, OWNER TO, and WITH CHECK are not supported based on AWS documentation.

-- Basic ALTER RLS POLICY with simple predicate
ALTER RLS POLICY policy_concerts USING (catgroup = 'concerts');

-- ALTER RLS POLICY with more specific filter
ALTER RLS POLICY policy_concerts USING (catgroup = 'piano concerts');

-- ALTER RLS POLICY with multiple conditions using AND
ALTER RLS POLICY sales_policy USING (region = 'US' AND year = 2024);

-- ALTER RLS POLICY with OR condition
ALTER RLS POLICY user_access_policy USING (username = 'alice' OR username = 'bob');

-- ALTER RLS POLICY with complex predicate using current_user
ALTER RLS POLICY user_data_policy USING (owner = current_user AND status = 'active');

-- ALTER RLS POLICY with IN clause
ALTER RLS POLICY regional_policy USING (region IN ('US', 'CA', 'MX'));

-- ALTER RLS POLICY with NOT IN clause
ALTER RLS POLICY exclusion_policy USING (category NOT IN ('restricted', 'confidential'));

-- ALTER RLS POLICY with numeric comparison
ALTER RLS POLICY price_policy USING (price > 100 AND price <= 1000);

-- ALTER RLS POLICY with BETWEEN clause
ALTER RLS POLICY date_range_policy USING (event_date BETWEEN '2024-01-01' AND '2024-12-31');

-- ALTER RLS POLICY with IS NULL check
ALTER RLS POLICY null_check_policy USING (deleted_at IS NULL);

-- ALTER RLS POLICY with IS NOT NULL check
ALTER RLS POLICY required_field_policy USING (customer_id IS NOT NULL);

-- ALTER RLS POLICY with function call
ALTER RLS POLICY time_based_policy USING (created_at > dateadd(day, -30, current_date));

-- ALTER RLS POLICY with nested conditions
ALTER RLS POLICY complex_policy USING ((status = 'active' AND region = 'US') OR (status = 'premium' AND region = 'EU'));

-- ALTER RLS POLICY with table alias reference
ALTER RLS POLICY qualified_policy USING (orders.customer_id = current_user_id);

-- ALTER RLS POLICY with mixed data types
ALTER RLS POLICY mixed_type_policy USING (is_active = true AND priority >= 5);

-- ALTER RLS POLICY with NOT condition
ALTER RLS POLICY negation_policy USING (NOT (category = 'internal' OR category = 'system'));

-- ALTER RLS POLICY with case-sensitive string comparison
ALTER RLS POLICY case_sensitive_policy USING (department = 'Engineering');

-- ALTER RLS POLICY with quoted identifier
ALTER RLS POLICY "Special-Policy-Name" USING (status = 'approved');

-- ALTER RLS POLICY with numeric literals
ALTER RLS POLICY numeric_policy USING (amount > 1000.50 AND quantity < 100);

-- ALTER RLS POLICY - comprehensive example combining multiple features
ALTER RLS POLICY comprehensive_access_policy USING (
    (user_role IN ('admin', 'manager') OR user_id = current_user_id)
    AND status NOT IN ('deleted', 'archived')
    AND created_date >= '2024-01-01'
    AND (region = 'US' OR global_access = true)
);

-- ALTER RLS POLICY with semicolon terminator
ALTER RLS POLICY terminated_policy USING (active = true);

-- ALTER RLS POLICY with single item IN list
ALTER RLS POLICY single_in_policy USING (role IN ('viewer'));

-- ALTER RLS POLICY with deeply nested parentheses
ALTER RLS POLICY nested_parens_policy USING (((((active = true)))));

-- ALTER RLS POLICY with all comparison operators
ALTER RLS POLICY comparison_policy USING (
    col1 = 'value1'
    AND col2 <> 'value2'
    AND col3 < 100
    AND col4 > 50
    AND col5 <= 200
    AND col6 >= 10
);

-- ALTER RLS POLICY with != operator (alternative to <>)
ALTER RLS POLICY not_equal_policy USING (status != 'inactive');