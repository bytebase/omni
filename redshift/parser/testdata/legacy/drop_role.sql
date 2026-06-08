-- DROP ROLE Test Cases
-- Based on AWS Redshift documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_ROLE.html

-- Basic DROP ROLE command (from AWS documentation example)
DROP ROLE sample_role1;

-- DROP ROLE with FORCE option (from AWS documentation example)
DROP ROLE sample_role FORCE;

-- DROP ROLE with RESTRICT option (default behavior)
DROP ROLE sample_role1 RESTRICT;

-- Basic role names
DROP ROLE analytics_role;
DROP ROLE data_reader_role;
DROP ROLE etl_processor_role;
DROP ROLE bi_analyst_role;
DROP ROLE temp_role;

-- Role names with underscores
DROP ROLE user_management_role;
DROP ROLE data_warehouse_admin;
DROP ROLE customer_analytics_reader;
DROP ROLE reporting_service_account;

-- Quoted role names for special characters
DROP ROLE "MyRole";
DROP ROLE "role-with-dash";
DROP ROLE "role.with.dots";
DROP ROLE "role with spaces";
DROP ROLE "UPPERCASE_ROLE";
DROP ROLE "Role_With_Mixed_Case";

-- Numeric role names (quoted)
DROP ROLE "123role";
DROP ROLE "role123";
DROP ROLE "role_v2";

-- Single character role names
DROP ROLE a;
DROP ROLE "A";
DROP ROLE z;

-- DROP ROLE with FORCE for roles with dependencies
DROP ROLE federated_role FORCE;
DROP ROLE sso_role FORCE;
DROP ROLE aws_role FORCE;
DROP ROLE azure_role FORCE;
DROP ROLE okta_role FORCE;

-- DROP ROLE with RESTRICT (explicit)
DROP ROLE special_char_role RESTRICT;
DROP ROLE unicode_role RESTRICT;
DROP ROLE long_external_id_role RESTRICT;

-- Real-world use case examples with FORCE
DROP ROLE etl_service_role FORCE;
DROP ROLE reporting_dashboard_role FORCE;
DROP ROLE data_scientist_role FORCE;
DROP ROLE readonly_analyst_role FORCE;
DROP ROLE emergency_access_role FORCE;

-- Service account patterns
DROP ROLE airflow_dag_role FORCE;
DROP ROLE kubernetes_role FORCE;
DROP ROLE lambda_function_role FORCE;

-- Development environment roles
DROP ROLE dev_readonly_role FORCE;
DROP ROLE staging_admin_role FORCE;
DROP ROLE test_data_role FORCE;

-- Compliance and audit roles
DROP ROLE audit_reader_role FORCE;
DROP ROLE security_monitor_role FORCE;
DROP ROLE compliance_officer_role FORCE;

-- Regional and department-specific roles
DROP ROLE north_america_analyst_role FORCE;
DROP ROLE europe_reader_role FORCE;
DROP ROLE finance_dept_role FORCE;
DROP ROLE marketing_analytics_role FORCE;

-- Temporary and contractor roles
DROP ROLE temp_contractor_role FORCE;
DROP ROLE consultant_access_role FORCE;
DROP ROLE vendor_readonly_role FORCE;

-- Machine learning and data science
DROP ROLE ml_model_role FORCE;
DROP ROLE data_pipeline_role FORCE;
DROP ROLE feature_store_role FORCE;

-- Case sensitivity tests
DROP ROLE lowercase_role;
DROP ROLE UPPERCASE_ROLE;
DROP ROLE "CamelCaseRole";

-- Case variations (should work the same)
drop role sample_role1;
Drop Role sample_role1;
DROP role sample_role1 FORCE;
drop ROLE sample_role1 restrict;

-- DROP ROLE with semicolon termination
DROP ROLE sample_role1;
DROP ROLE sample_role1 FORCE;
DROP ROLE sample_role1 RESTRICT;

-- DROP ROLE without semicolon (should also work)
DROP ROLE sample_role1
DROP ROLE sample_role1 FORCE
DROP ROLE sample_role1 RESTRICT

-- Edge cases and boundary testing
DROP ROLE a FORCE;  -- Minimal role name with FORCE
DROP ROLE longest_possible_role_name_within_identifier_limits FORCE;

-- Test SQL keywords as quoted identifiers
DROP ROLE "SELECT" FORCE;
DROP ROLE "CREATE" FORCE;
DROP ROLE "TABLE" FORCE;
DROP ROLE "USER" FORCE;
DROP ROLE "ROLE" FORCE;

-- Complex scenarios based on AWS documentation examples
-- These would be used after creating roles with dependencies
DROP ROLE sample_role1 FORCE;  -- Removes all role assignments if any exist
DROP ROLE sample_role2 FORCE;  -- For roles granted to other roles

-- Default RESTRICT behavior (these would fail if role has dependencies)
DROP ROLE independent_role;  -- Should work for roles without dependencies
DROP ROLE standalone_role RESTRICT;  -- Explicit RESTRICT