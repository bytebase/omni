-- DROP RLS POLICY Test Cases
-- Based on AWS Redshift documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_RLS_POLICY.html

-- Basic DROP RLS POLICY command
DROP RLS POLICY policy_concerts;

-- DROP RLS POLICY with IF EXISTS clause
DROP RLS POLICY IF EXISTS policy_concerts;

-- DROP RLS POLICY with CASCADE option
DROP RLS POLICY policy_concerts CASCADE;

-- DROP RLS POLICY with RESTRICT option (default behavior)
DROP RLS POLICY policy_concerts RESTRICT;

-- DROP RLS POLICY with IF EXISTS and CASCADE
DROP RLS POLICY IF EXISTS policy_concerts CASCADE;

-- DROP RLS POLICY with IF EXISTS and RESTRICT
DROP RLS POLICY IF EXISTS policy_concerts RESTRICT;

-- DROP RLS POLICY with quoted policy name
DROP RLS POLICY "policy_concerts";

-- DROP RLS POLICY with quoted policy name and IF EXISTS
DROP RLS POLICY IF EXISTS "policy_concerts";

-- DROP RLS POLICY with quoted policy name and CASCADE
DROP RLS POLICY "policy_concerts" CASCADE;

-- DROP RLS POLICY with quoted policy name and RESTRICT
DROP RLS POLICY "policy_concerts" RESTRICT;

-- DROP RLS POLICY with complex policy names
DROP RLS POLICY policy_sales_2023;
DROP RLS POLICY policy_customer_data;
DROP RLS POLICY rls_finance_reports;

-- DROP RLS POLICY with underscores and numbers in policy name
DROP RLS POLICY policy_user_123;
DROP RLS POLICY rls_policy_v2;
DROP RLS POLICY data_access_policy_final;

-- DROP RLS POLICY with IF EXISTS for non-existent policies (should not error)
DROP RLS POLICY IF EXISTS non_existent_policy;
DROP RLS POLICY IF EXISTS "another_missing_policy";

-- Case variations (should work the same)
drop rls policy policy_concerts;
Drop RLS Policy policy_concerts;
DROP rls policy policy_concerts CASCADE;
drop RLS POLICY IF EXISTS policy_concerts restrict;

-- DROP RLS POLICY with semicolon termination
DROP RLS POLICY policy_concerts;
DROP RLS POLICY IF EXISTS policy_concerts;
DROP RLS POLICY policy_concerts CASCADE;
DROP RLS POLICY policy_concerts RESTRICT;

-- DROP RLS POLICY without semicolon (should also work)
DROP RLS POLICY policy_concerts
DROP RLS POLICY IF EXISTS policy_concerts
DROP RLS POLICY policy_concerts CASCADE
DROP RLS POLICY policy_concerts RESTRICT