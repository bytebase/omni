-- Basic ATTACH RLS POLICY examples based on AWS documentation

-- Example 1: Attach policy to a single table for specific roles
ATTACH RLS POLICY policy_concerts
ON tickit_category_redshift
TO ROLE analyst, ROLE dbadmin;

-- Example 2: Attach policy to a table for specific users
ATTACH RLS POLICY sales_policy
ON sales_data
TO john_doe, jane_smith;

-- Example 3: Attach policy with TABLE keyword
ATTACH RLS POLICY regional_policy
ON TABLE customer_orders
TO ROLE regional_manager;

-- Example 4: Attach policy to multiple tables
ATTACH RLS POLICY data_protection_policy
ON orders, customers, transactions
TO ROLE data_analyst, ROLE compliance_officer;

-- Example 5: Attach policy to PUBLIC (all users)
ATTACH RLS POLICY basic_filter_policy
ON employee_directory
TO PUBLIC;

-- Example 6: Mixed targets - users, roles, and PUBLIC
ATTACH RLS POLICY mixed_access_policy
ON inventory_table
TO admin_user, ROLE warehouse_manager, PUBLIC;

-- Example 7: Schema-qualified table names
ATTACH RLS POLICY schema_policy
ON sales.orders, sales.customers
TO ROLE sales_team;

-- Example 8: Complex multi-table, multi-role attachment
ATTACH RLS POLICY comprehensive_policy
ON TABLE finance.transactions, finance.accounts, finance.budgets
TO ROLE finance_analyst, ROLE finance_manager, ROLE auditor;

-- Example 9: Single character identifiers (edge case)
ATTACH RLS POLICY p
ON t
TO ROLE r;

-- Example 10: Longer identifier names
ATTACH RLS POLICY very_long_policy_name_for_testing_limits
ON table_with_a_very_long_name_to_test_parser
TO ROLE role_with_long_name, user_with_long_name;

-- Example 11: Multiple users without roles
ATTACH RLS POLICY user_only_policy
ON products
TO user1, user2, user3, user4;

-- Example 12: Multiple roles without users
ATTACH RLS POLICY role_only_policy
ON categories
TO ROLE role1, ROLE role2, ROLE role3;

-- Example 13: Views and materialized views (as mentioned in docs)
ATTACH RLS POLICY view_policy
ON my_view, my_materialized_view
TO ROLE view_reader;

-- Example 14: Database.schema.table format
ATTACH RLS POLICY cross_schema_policy
ON mydb.public.orders, mydb.sales.customers
TO ROLE cross_schema_user;

-- Example 15: Policy with numbers in names
ATTACH RLS POLICY policy_2024_q1
ON sales_2024_q1
TO ROLE analyst_2024;