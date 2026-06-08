-- ATTACH MASKING POLICY test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ATTACH_MASKING_POLICY.html

-- Basic attach to a single column for a user
ATTACH MASKING POLICY credit_card_mask
ON customer_data
(credit_card_number)
TO john_doe;

-- Attach policy to multiple columns for a user
ATTACH MASKING POLICY pii_mask
ON employee_table
(ssn, phone_number, email)
TO alice_smith;

-- Attach policy with input columns
ATTACH MASKING POLICY conditional_mask
ON sales_data
(revenue, profit)
USING (region, department)
TO bob_johnson;

-- Attach policy to a role
ATTACH MASKING POLICY salary_mask
ON hr_data
(salary, bonus)
TO ROLE finance_team;

-- Attach policy to PUBLIC (all users)
ATTACH MASKING POLICY default_mask
ON sensitive_table
(personal_info)
TO PUBLIC;

-- Attach policy with priority
ATTACH MASKING POLICY high_priority_mask
ON customer_info
(account_number)
TO susan_lee
PRIORITY 10;

-- Attach policy with multiple input and output columns
ATTACH MASKING POLICY complex_mask
ON transaction_table
(card_number, cvv, expiry_date)
USING (transaction_type, user_role, access_level)
TO ROLE security_admin
PRIORITY 5;

-- Attach policy to SUPER data type paths
ATTACH MASKING POLICY json_mask
ON customer_json
(person.name.first_name, person.name.last_name, person.contact.email)
TO david_wong;

-- Attach policy with SUPER paths in USING clause
ATTACH MASKING POLICY super_conditional_mask
ON json_data
(sensitive_data.credit_info)
USING (user_profile.access_level, user_profile.department)
TO ROLE data_analyst;

-- Attach policy with mixed column types
ATTACH MASKING POLICY mixed_mask
ON hybrid_table
(regular_column, super_data.nested.field)
USING (access_type)
TO PUBLIC
PRIORITY 0;

-- Attach policy with quoted identifiers
ATTACH MASKING POLICY "Special-Mask-Policy"
ON "Customer-Table"
("Credit-Card", "SSN-Number")
TO "user-with-dash";

-- Attach policy with quoted identifier and role
ATTACH MASKING POLICY policy_123
ON "sales_2023"
("total_revenue", "customer_id")
TO ROLE "Sales-Team-Role"
PRIORITY 3;

-- Complex example with all features
ATTACH MASKING POLICY comprehensive_mask
ON financial_records
(account_balance, routing_number, account_holder.personal.ssn)
USING (user_clearance_level, department_code, access_timestamp)
TO ROLE executive_team
PRIORITY 99;

-- Multiple policies for different users on same table
ATTACH MASKING POLICY basic_mask
ON employee_records
(salary)
TO PUBLIC
PRIORITY 0;

ATTACH MASKING POLICY detailed_mask
ON employee_records
(salary, bonus, stock_options)
TO ROLE hr_manager
PRIORITY 10;

-- Edge cases with special characters in identifiers
ATTACH MASKING POLICY "mask@policy#1"
ON "table$special%chars"
("column&name*1", "column^name!2")
TO "user@company.com";

-- Maximum length identifiers (63 characters)
ATTACH MASKING POLICY very_long_policy_name_that_reaches_maximum_allowed_length_12345
ON table_with_very_long_name_that_also_reaches_maximum_length_123
(column_with_extremely_long_name_reaching_sixty_three_chars_12345)
TO user_with_long_name_reaching_maximum_allowed_identifier_length;

-- Comments in statements
-- This policy masks sensitive financial data
ATTACH MASKING POLICY financial_mask
ON bank_accounts  -- Main banking table
(account_number, balance)  -- Sensitive columns
TO ROLE teller  -- Bank teller role
PRIORITY 5;  -- Medium priority

/* Multi-line comment example
   This attaches a comprehensive masking policy
   for GDPR compliance */
ATTACH MASKING POLICY gdpr_mask
ON eu_customers
(name, email, phone, address)
TO PUBLIC
PRIORITY 1;