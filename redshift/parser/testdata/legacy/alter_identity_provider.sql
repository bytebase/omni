-- ALTER IDENTITY PROVIDER test cases
-- Based on AWS Redshift documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_IDENTITY_PROVIDER.html

-- Test Case 1: Basic ALTER with PARAMETERS for Azure AD OAuth
ALTER IDENTITY PROVIDER oauth_standard
    PARAMETERS '{"issuer":"https://sts.windows.net/2sdfdsf-d475-420d-b5ac-667adad7c702/",
                 "client_id":"87f4aa26-78b7-410e-bf29-57b39929ef9a",
                 "client_secret":"BUAH~ewrqewrqwerUUY^%tHe1oNZShoiU7",
                 "audience":["https://analysis.windows.net/powerbi/connector/AmazonRedshift"]}';

-- Test Case 2: ALTER with NAMESPACE for IAM Identity Center
ALTER IDENTITY PROVIDER "my-redshift-idc-application"
    NAMESPACE 'MYCO';

-- Test Case 3: ALTER with IAM_ROLE
ALTER IDENTITY PROVIDER awsidc_provider
    IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftIdentityProviderRole';

-- Test Case 4: Enable AUTO_CREATE_ROLES with default settings
ALTER IDENTITY PROVIDER oauth_provider
    AUTO_CREATE_ROLES TRUE;

-- Test Case 5: Enable AUTO_CREATE_ROLES with INCLUDE filter
ALTER IDENTITY PROVIDER oauth_provider
    AUTO_CREATE_ROLES TRUE INCLUDE GROUPS LIKE 'admin_%';

-- Test Case 6: Enable AUTO_CREATE_ROLES with EXCLUDE filter
ALTER IDENTITY PROVIDER azure_ad_provider
    AUTO_CREATE_ROLES TRUE EXCLUDE GROUPS LIKE 'temp_%';

-- Test Case 7: Disable AUTO_CREATE_ROLES
ALTER IDENTITY PROVIDER oauth_provider
    AUTO_CREATE_ROLES FALSE;

-- Test Case 8: DISABLE identity provider
ALTER IDENTITY PROVIDER "redshift-idc-app" DISABLE;

-- Test Case 9: ENABLE identity provider
ALTER IDENTITY PROVIDER "redshift-idc-app" ENABLE;

-- Test Case 10: Multiple parameters in one ALTER statement
ALTER IDENTITY PROVIDER oauth_standard
    PARAMETERS '{"issuer":"https://login.microsoftonline.com/...",
                 "client_id":"app-id-123",
                 "client_secret":"secret-value"}'
    NAMESPACE 'PRODUCTION'
    AUTO_CREATE_ROLES TRUE INCLUDE GROUPS LIKE 'prod_%';

-- Test Case 11: ALTER with quoted identifier containing special characters
ALTER IDENTITY PROVIDER "my-provider-123"
    PARAMETERS '{"config":"value"}';

-- Test Case 12: ALTER with simple unquoted identifier
ALTER IDENTITY PROVIDER simple_provider
    NAMESPACE 'DEFAULT';

-- Test Case 13: Complex AUTO_CREATE_ROLES pattern
ALTER IDENTITY PROVIDER oauth_provider
    AUTO_CREATE_ROLES TRUE EXCLUDE GROUPS LIKE '%_test_%';

-- Test Case 14: Minimal ALTER - just disable
ALTER IDENTITY PROVIDER minimal_provider DISABLE;

-- Test Case 15: IAM Identity Center with all applicable options
ALTER IDENTITY PROVIDER "idc-application"
    NAMESPACE 'MYCOMPANY'
    IAM_ROLE 'arn:aws:iam::987654321098:role/MyIdentityRole'
    AUTO_CREATE_ROLES TRUE INCLUDE GROUPS LIKE 'aws_%'
    ENABLE;

-- Test Case 16: Parameters with nested JSON
ALTER IDENTITY PROVIDER complex_oauth
    PARAMETERS '{"issuer":"https://auth.example.com",
                 "client_id":"complex-app",
                 "client_secret":"complex-secret",
                 "audience":["aud1", "aud2", "aud3"],
                 "additional_config":{"timeout":300, "retry":true}}';

-- Test Case 17: Multiple actions without PARAMETERS
ALTER IDENTITY PROVIDER no_params_provider
    NAMESPACE 'TEST'
    AUTO_CREATE_ROLES FALSE
    DISABLE;

-- Test Case 18: Pattern with wildcards on both sides
ALTER IDENTITY PROVIDER pattern_provider
    AUTO_CREATE_ROLES TRUE INCLUDE GROUPS LIKE '%admin%';

-- Test Case 19: Empty namespace
ALTER IDENTITY PROVIDER empty_namespace_provider
    NAMESPACE '';

-- Test Case 20: Long IAM role ARN
ALTER IDENTITY PROVIDER long_arn_provider
    IAM_ROLE 'arn:aws:iam::123456789012:role/service-role/RedshiftIdentityProviderServiceRole-us-east-1-production';