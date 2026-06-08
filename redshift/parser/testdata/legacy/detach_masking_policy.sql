-- DETACH MASKING POLICY test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DETACH_MASKING_POLICY.html

-- Basic detach from a single column for a user
DETACH MASKING POLICY credit_card_mask
ON customer_data
(credit_card_number)
FROM john_doe;

-- Detach policy from multiple columns for a user
DETACH MASKING POLICY pii_mask
ON employee_table
(ssn, phone_number, email)
FROM alice_smith;

-- Detach policy from a role
DETACH MASKING POLICY salary_mask
ON hr_data
(salary, bonus)
FROM ROLE finance_team;

-- Detach policy from PUBLIC (all users)
DETACH MASKING POLICY default_mask
ON sensitive_table
(personal_info)
FROM PUBLIC;

-- Detach policy from SUPER data type paths
DETACH MASKING POLICY json_mask
ON customer_json
(person.name.first_name, person.name.last_name, person.contact.email)
FROM david_wong;

-- Detach policy with complex SUPER paths
DETACH MASKING POLICY super_conditional_mask
ON json_data
(sensitive_data.credit_info)
FROM ROLE data_analyst;

-- Detach policy with mixed column types
DETACH MASKING POLICY mixed_mask
ON hybrid_table
(regular_column, super_data.nested.field)
FROM PUBLIC;

-- Detach policy with quoted identifiers
DETACH MASKING POLICY "Special-Mask-Policy"
ON "Customer-Table"
("Credit-Card", "SSN-Number")
FROM "user-with-dash";

-- Detach policy with quoted identifier and role
DETACH MASKING POLICY policy_123
ON "sales_2023"
("total_revenue", "customer_id")
FROM ROLE "Sales-Team-Role";

-- Complex example with multiple columns and SUPER paths
DETACH MASKING POLICY comprehensive_mask
ON financial_records
(account_balance, routing_number, account_holder.personal.ssn)
FROM ROLE executive_team;

-- Detach multiple policies from different users on same table
DETACH MASKING POLICY basic_mask
ON employee_records
(salary)
FROM PUBLIC;

DETACH MASKING POLICY detailed_mask
ON employee_records
(salary, bonus, stock_options)
FROM ROLE hr_manager;

-- Edge cases with special characters in identifiers
DETACH MASKING POLICY "mask@policy#1"
ON "table$special%chars"
("column&name*1", "column^name!2")
FROM "user@company.com";

-- Maximum length identifiers (63 characters)
DETACH MASKING POLICY very_long_policy_name_that_reaches_maximum_allowed_length_12345
ON table_with_very_long_name_that_also_reaches_maximum_length_123
(column_with_extremely_long_name_reaching_sixty_three_chars_12345)
FROM user_with_long_name_reaching_maximum_allowed_identifier_length;

-- Detach from nested SUPER columns with deep paths
DETACH MASKING POLICY deep_nested_mask
ON analytics_data
(user_profile.personal.contact.emergency.phone, user_profile.financial.accounts.primary.number)
FROM ROLE analytics_team;

-- Detach policy from a single deeply nested column
DETACH MASKING POLICY single_nested_mask
ON order_data
(customer.billing.address.street)
FROM order_processor;

-- Detach from columns with numeric suffixes
DETACH MASKING POLICY numeric_column_mask
ON data_table_v2
(column1, column2, column3, column_final)
FROM ROLE data_team;

-- Detach using case variations (should be case-insensitive)
detach masking policy lowercase_test
on test_table
(test_column)
from test_user;

DETACH MASKING POLICY UPPERCASE_TEST
ON TEST_TABLE
(TEST_COLUMN)
FROM ROLE TEST_ROLE;

-- Detach with mixed case
Detach Masking Policy MixedCase_Test
On Mixed_Table
(Mixed_Column)
From Public;

-- Comments in statements
-- This removes the financial data mask
DETACH MASKING POLICY financial_mask
ON bank_accounts  -- Main banking table
(account_number, balance)  -- Previously masked columns
FROM ROLE teller;  -- Bank teller role

/* Multi-line comment example
   This detaches a comprehensive masking policy
   that was used for GDPR compliance */
DETACH MASKING POLICY gdpr_mask
ON eu_customers
(name, email, phone, address)
FROM PUBLIC;

-- Real-world scenarios based on AWS documentation patterns

-- Detaching HR-related masking policies
DETACH MASKING POLICY employee_ssn_mask
ON employees
(social_security_number)
FROM ROLE hr_administrator;

-- Detaching financial masking policies
DETACH MASKING POLICY account_mask
ON customer_accounts
(account_number, routing_number)
FROM financial_analyst;

-- Detaching healthcare data masking
DETACH MASKING POLICY patient_data_mask
ON medical_records
(patient_id, diagnosis, treatment_notes)
FROM ROLE healthcare_provider;

-- Detaching from multiple related columns
DETACH MASKING POLICY contact_info_mask
ON customer_contacts
(phone_home, phone_work, phone_mobile, email_personal, email_work)
FROM customer_service_rep;

-- Detaching from temporal data
DETACH MASKING POLICY temporal_mask
ON transaction_history
(transaction_date, amount, merchant_info)
FROM ROLE audit_team;

-- Detaching cross-schema references (if supported)
DETACH MASKING POLICY cross_schema_mask
ON production.customer_data
(sensitive_field)
FROM analytics_user;

-- Detaching with schema-qualified table names
DETACH MASKING POLICY schema_qualified_mask
ON finance.account_details
(balance, credit_limit)
FROM ROLE finance_viewer;

-- Edge case: single character identifiers
DETACH MASKING POLICY a
ON b
(c)
FROM d;

-- Edge case: identifiers with numbers
DETACH MASKING POLICY mask123
ON table456
(col789)
FROM user000;

-- Edge case: underscore variations
DETACH MASKING POLICY _mask_policy_
ON _table_name_
(_column_name_)
FROM _user_name_;

-- Whitespace tolerance test
DETACH   MASKING   POLICY   spaced_policy
ON   spaced_table
(   spaced_column   )
FROM   spaced_user   ;

-- Testing optional semicolon
DETACH MASKING POLICY no_semicolon_policy
ON test_table
(test_column)
FROM test_user

-- With semicolon for comparison
DETACH MASKING POLICY with_semicolon_policy
ON test_table
(test_column)
FROM test_user;