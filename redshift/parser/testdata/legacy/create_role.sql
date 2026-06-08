-- CREATE ROLE test cases based on AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_ROLE.html

-- Basic role creation (from AWS documentation)
CREATE ROLE sample_role1;

-- Role creation with external ID (from AWS documentation)
CREATE ROLE sample_role1 EXTERNALID "ABC123";

-- Additional basic role creation examples
CREATE ROLE analytics_role;
CREATE ROLE data_reader_role;
CREATE ROLE etl_processor_role;
CREATE ROLE bi_analyst_role;
CREATE ROLE temp_role;

-- Role names with underscores
CREATE ROLE user_management_role;
CREATE ROLE data_warehouse_admin;
CREATE ROLE customer_analytics_reader;
CREATE ROLE reporting_service_account;

-- Quoted role names for special characters
CREATE ROLE "MyRole";
CREATE ROLE "role-with-dash";
CREATE ROLE "role.with.dots";
CREATE ROLE "role with spaces";
CREATE ROLE "UPPERCASE_ROLE";
CREATE ROLE "Role_With_Mixed_Case";

-- Numeric role names (quoted)
CREATE ROLE "123role";
CREATE ROLE "role123";
CREATE ROLE "role_v2";

-- Single character role names
CREATE ROLE a;
CREATE ROLE "A";
CREATE ROLE z;

-- External ID variations
CREATE ROLE federated_role EXTERNALID "XYZ456";
CREATE ROLE sso_role EXTERNALID "SSO_PROVIDER_ID_12345";
CREATE ROLE aws_role EXTERNALID "arn:aws:iam::123456789012:role/MyRedshiftRole";
CREATE ROLE azure_role EXTERNALID "azure-ad-group-id-abcd1234";
CREATE ROLE okta_role EXTERNALID "okta-group-00g1emaKYMLWCQXVHDSN";

-- External IDs with special characters
CREATE ROLE special_char_role EXTERNALID "ID_with-special.chars@domain.com";
CREATE ROLE unicode_role EXTERNALID "üñíçødé_íd";
CREATE ROLE long_external_id_role EXTERNALID "very_long_external_id_that_contains_many_characters_and_underscores";

-- External IDs with spaces and quotes
CREATE ROLE space_role EXTERNALID "external id with spaces";
CREATE ROLE quote_role EXTERNALID "external 'id' with quotes";
CREATE ROLE mixed_role EXTERNALID "Mixed Case External ID 123";

-- Numeric external IDs
CREATE ROLE numeric_external_role EXTERNALID "123456789";
CREATE ROLE decimal_external_role EXTERNALID "123.456.789";

-- Empty external ID (edge case)
CREATE ROLE empty_external_role EXTERNALID "";

-- Real-world use case examples
CREATE ROLE etl_service_role EXTERNALID "arn:aws:iam::987654321098:role/ETLServiceRole";
CREATE ROLE reporting_dashboard_role EXTERNALID "tableau-service-account-xyz";
CREATE ROLE data_scientist_role EXTERNALID "data-science-team-okta-group";
CREATE ROLE readonly_analyst_role EXTERNALID "readonly-analysts-ad-group";
CREATE ROLE emergency_access_role EXTERNALID "emergency-access-break-glass";

-- Service account patterns
CREATE ROLE airflow_dag_role EXTERNALID "airflow-production-service-account";
CREATE ROLE kubernetes_role EXTERNALID "k8s-redshift-operator-sa";
CREATE ROLE lambda_function_role EXTERNALID "arn:aws:iam::555666777888:role/LambdaRedshiftAccess";

-- Development environment roles
CREATE ROLE dev_readonly_role EXTERNALID "dev-team-readonly-access";
CREATE ROLE staging_admin_role EXTERNALID "staging-environment-admin";
CREATE ROLE test_data_role EXTERNALID "test-data-generator-service";

-- Compliance and audit roles
CREATE ROLE audit_reader_role EXTERNALID "compliance-audit-team-group";
CREATE ROLE security_monitor_role EXTERNALID "security-monitoring-service";
CREATE ROLE compliance_officer_role EXTERNALID "compliance-officer-access-group";

-- Regional and department-specific roles
CREATE ROLE north_america_analyst_role EXTERNALID "na-regional-analysts";
CREATE ROLE europe_reader_role EXTERNALID "eu-regional-readonly";
CREATE ROLE finance_dept_role EXTERNALID "finance-department-users";
CREATE ROLE marketing_analytics_role EXTERNALID "marketing-analytics-team";

-- Temporary and contractor roles
CREATE ROLE temp_contractor_role EXTERNALID "temporary-contractor-q4-2024";
CREATE ROLE consultant_access_role EXTERNALID "external-consultant-project-alpha";
CREATE ROLE vendor_readonly_role EXTERNALID "vendor-readonly-limited-access";

-- Machine learning and data science
CREATE ROLE ml_model_role EXTERNALID "ml-model-training-service";
CREATE ROLE data_pipeline_role EXTERNALID "automated-data-pipeline-v2";
CREATE ROLE feature_store_role EXTERNALID "feature-store-service-account";

-- Case sensitivity tests
CREATE ROLE lowercase_role EXTERNALID "lowercase_external_id";
CREATE ROLE UPPERCASE_ROLE EXTERNALID "UPPERCASE_EXTERNAL_ID";
CREATE ROLE "CamelCaseRole" EXTERNALID "CamelCaseExternalId";

-- Special AWS IAM patterns
CREATE ROLE cross_account_role EXTERNALID "arn:aws:iam::111222333444:role/CrossAccountAccess";
CREATE ROLE assume_role_pattern EXTERNALID "arn:aws:iam::123456789012:role/service-role/MyServiceRole";
CREATE ROLE federated_web_identity EXTERNALID "arn:aws:iam::123456789012:oidc-provider/example.com";

-- Edge cases and boundary testing
CREATE ROLE a EXTERNALID "a";  -- Minimal role and external ID
CREATE ROLE longest_possible_role_name_within_identifier_limits EXTERNALID "longest_possible_external_id_string_that_might_be_used";

-- Test SQL keywords as quoted identifiers (if supported)
CREATE ROLE "SELECT" EXTERNALID "select-service-account";
CREATE ROLE "CREATE" EXTERNALID "create-operation-role";
CREATE ROLE "TABLE" EXTERNALID "table-management-service";
CREATE ROLE "USER" EXTERNALID "user-management-system";
CREATE ROLE "ROLE" EXTERNALID "role-management-service";