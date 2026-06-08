-- DETACH RLS POLICY test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_DETACH_RLS_POLICY.html

-- Example from AWS documentation
DETACH RLS POLICY policy_concerts ON tickit_category_redshift
FROM ROLE analyst, ROLE dbadmin;

-- Basic detach from a single table for a specific user
DETACH RLS POLICY user_access_policy
ON employee_data
FROM john_doe;

-- Detach policy from multiple users
DETACH RLS POLICY data_filter_policy
ON sales_records
FROM alice_smith, bob_johnson, carol_white;

-- Detach policy from a single role
DETACH RLS POLICY finance_policy
ON financial_records
FROM ROLE finance_team;

-- Detach policy from multiple roles
DETACH RLS POLICY audit_policy
ON transaction_logs
FROM ROLE auditor, ROLE compliance_officer, ROLE security_admin;

-- Detach policy from PUBLIC (all users)
DETACH RLS POLICY public_access_policy
ON public_announcements
FROM PUBLIC;

-- Detach policy from mixed targets - users, roles, and PUBLIC
DETACH RLS POLICY mixed_access_policy
ON inventory_table
FROM admin_user, ROLE warehouse_manager, PUBLIC;

-- Detach policy with TABLE keyword
DETACH RLS POLICY regional_policy
ON TABLE customer_orders
FROM ROLE regional_manager;

-- Detach policy from multiple tables
DETACH RLS POLICY cross_table_policy
ON orders, customers, transactions
FROM ROLE data_analyst, ROLE business_user;

-- Detach policy from multiple tables with TABLE keyword
DETACH RLS POLICY data_protection_policy
ON TABLE sales_data, customer_profiles, order_history
FROM ROLE marketing_team;

-- Schema-qualified table names
DETACH RLS POLICY schema_policy
ON sales.orders, sales.customers
FROM ROLE sales_team;

-- Database.schema.table format
DETACH RLS POLICY cross_schema_policy
ON mydb.public.orders, mydb.sales.customers
FROM ROLE cross_schema_user;

-- Complex multi-table, multi-target detachment
DETACH RLS POLICY comprehensive_policy
ON TABLE finance.transactions, finance.accounts, finance.budgets
FROM ROLE finance_analyst, ROLE finance_manager, ROLE auditor, finance_admin;

-- Detach from views and materialized views
DETACH RLS POLICY view_policy
ON customer_view, sales_mv
FROM ROLE view_reader;

-- Policy with quoted identifiers
DETACH RLS POLICY "Special-Policy-Name"
ON "Customer-Table", "Sales-Data"
FROM "user-with-dashes", ROLE "Role-With-Dashes";

-- Policy with numbers in names
DETACH RLS POLICY policy_2024_q1
ON sales_2024_q1, revenue_2024_q1
FROM ROLE analyst_2024, quarterly_reviewer;

-- Single character identifiers (edge case)
DETACH RLS POLICY p
ON t
FROM ROLE r;

-- Long identifier names
DETACH RLS POLICY very_long_policy_name_for_testing_limits_and_boundaries
ON table_with_a_very_long_name_to_test_parser_capabilities
FROM ROLE role_with_extremely_long_name_for_testing_purposes;

-- Multiple users without roles
DETACH RLS POLICY user_only_policy
ON products
FROM user1, user2, user3, user4;

-- Multiple roles without users
DETACH RLS POLICY role_only_policy
ON categories
FROM ROLE role1, ROLE role2, ROLE role3;

-- Identifiers with underscores
DETACH RLS POLICY _policy_with_underscores_
ON _table_with_underscores_
FROM _user_with_underscores_, ROLE _role_with_underscores_;

-- Identifiers with numbers and mixed case
DETACH RLS POLICY Policy123
ON Table456, View789
FROM User000, ROLE Role999;

-- Case sensitivity variations (should be case-insensitive)
detach rls policy lowercase_test
on test_table
from test_user;

DETACH RLS POLICY UPPERCASE_TEST
ON TEST_TABLE
FROM ROLE TEST_ROLE;

Detach Rls Policy MixedCase_Test
On Mixed_Table
From Public;

-- Healthcare data example
DETACH RLS POLICY patient_privacy_policy
ON medical_records, patient_info
FROM ROLE doctor, ROLE nurse, medical_admin;

-- Financial services example
DETACH RLS POLICY account_security_policy
ON customer_accounts, transaction_history
FROM ROLE teller, ROLE loan_officer, bank_manager;

-- E-commerce example
DETACH RLS POLICY order_visibility_policy
ON orders, order_items, shipments
FROM ROLE customer_service, ROLE fulfillment_team;

-- HR data example
DETACH RLS POLICY employee_confidentiality_policy
ON employee_records, salary_info, performance_reviews
FROM ROLE hr_manager, ROLE payroll_admin;

-- Multi-tenant application example
DETACH RLS POLICY tenant_isolation_policy
ON tenant_data, tenant_config
FROM ROLE tenant_admin, application_user;

-- Compliance and audit example
DETACH RLS POLICY gdpr_compliance_policy
ON eu_customer_data, privacy_logs
FROM ROLE compliance_officer, ROLE privacy_officer, PUBLIC;

-- Development environment example
DETACH RLS POLICY dev_access_policy
ON dev_tables, test_data
FROM ROLE developer, ROLE qa_tester, dev_admin;

-- Analytics and reporting example
DETACH RLS POLICY reporting_policy
ON analytics_data, business_metrics
FROM ROLE business_analyst, ROLE data_scientist, report_viewer;

-- Security and monitoring example
DETACH RLS POLICY security_monitoring_policy
ON security_logs, access_logs, audit_trail
FROM ROLE security_admin, ROLE system_monitor;

-- Data warehouse example
DETACH RLS POLICY warehouse_access_policy
ON fact_sales, dim_customer, dim_product
FROM ROLE data_engineer, ROLE warehouse_admin;

-- Cross-department access
DETACH RLS POLICY cross_dept_policy
ON shared_resources, common_lookup_tables
FROM ROLE marketing, ROLE sales, ROLE support, shared_user;

-- Temporal data access
DETACH RLS POLICY temporal_access_policy
ON historical_data, archived_records
FROM ROLE historian, ROLE archive_admin;

-- Geographic data access
DETACH RLS POLICY geographic_policy
ON location_data, regional_sales
FROM ROLE regional_manager_us, ROLE regional_manager_eu;

-- Special characters in quoted identifiers
DETACH RLS POLICY "policy@domain.com"
ON "table$special%chars", "view#with&symbols"
FROM "user@company.com", ROLE "role!with*special^chars";

-- Maximum complexity example
DETACH RLS POLICY enterprise_master_policy
ON TABLE production.sales.orders,
         production.sales.customers,
         production.inventory.products,
         staging.imports.raw_data,
         analytics.reports.daily_summary
FROM enterprise_admin,
     data_steward,
     business_owner,
     ROLE enterprise_viewer,
     ROLE data_analyst,
     ROLE business_intelligence,
     ROLE system_administrator,
     PUBLIC;

-- Edge case: detaching from everything
DETACH RLS POLICY universal_policy
ON all_tables, all_views, all_materialized_views
FROM all_users, ROLE all_roles, PUBLIC;

-- Whitespace tolerance test
DETACH   RLS   POLICY   spaced_policy
ON   spaced_table   ,   another_table
FROM   spaced_user   ,   ROLE   spaced_role   ,   PUBLIC   ;

-- Comments in statements
-- This removes the customer data access policy
DETACH RLS POLICY customer_data_policy
ON customers  -- Customer master table
FROM ROLE customer_service,  -- Customer service representatives
     ROLE sales_team;  -- Sales team members

/* Multi-line comment example
   This detaches a comprehensive data governance policy
   that was implemented for regulatory compliance */
DETACH RLS POLICY regulatory_compliance_policy
ON sensitive_financial_data, customer_pii, transaction_records
FROM ROLE compliance_team, ROLE legal_team, audit_manager;

-- Real-world enterprise scenarios

-- Detaching departmental policies
DETACH RLS POLICY hr_department_policy
ON employee_data, payroll_info, benefits_data
FROM ROLE hr_manager, ROLE hr_specialist, hr_admin;

-- Detaching project-based access
DETACH RLS POLICY project_alpha_policy
ON project_alpha_data, project_alpha_reports
FROM ROLE project_alpha_team, project_alpha_lead;

-- Detaching time-based access (end of quarter)
DETACH RLS POLICY q4_2024_policy
ON q4_sales_data, q4_financial_reports
FROM ROLE q4_analyst, ROLE quarterly_reviewer;

-- Detaching contractor access
DETACH RLS POLICY contractor_access_policy
ON limited_project_data, public_documentation
FROM contractor_user1, contractor_user2, ROLE contractor_team;

-- Emergency access removal
DETACH RLS POLICY emergency_access_policy
ON critical_systems_data, emergency_procedures
FROM ROLE emergency_responder, ROLE incident_manager;

-- Vendor access removal
DETACH RLS POLICY vendor_integration_policy
ON integration_logs, api_access_data
FROM vendor_integration_user, ROLE vendor_support;