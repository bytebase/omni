-- CREATE IDENTITY PROVIDER test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_IDENTITY_PROVIDER.html

-- Example 1: Azure OAuth Identity Provider with complete parameters
CREATE IDENTITY PROVIDER oauth_standard
TYPE azure
NAMESPACE 'aad'
PARAMETERS '{"issuer":"https://sts.windows.net/2sdfdsf-d475-420d-b5ac-667adad7c702/", "client_id":"87f4aa26-78b7-410e-bf29-57b39929ef9a", "client_secret":"BUAH~ewrqewrqwerUUY^%tHe1oNZShoiU7", "audience":["https://analysis.windows.net/powerbi/connector/AmazonRedshift"] }';

-- Example 2: AWS IAM Identity Center (AWSIDC) with quoted provider name
CREATE IDENTITY PROVIDER "redshift-idc-app"
TYPE AWSIDC
NAMESPACE 'awsidc'
APPLICATION_ARN 'arn:aws:sso::123456789012:application/ssoins-12345f67fe123d4/apl-a0b0a12dc123b1a4'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole';

-- Example 3: IAM Identity Center with explicit type name
CREATE IDENTITY PROVIDER my_idc_provider
TYPE IAM_IDENTITY_CENTER
NAMESPACE 'corporate_sso'
APPLICATION_ARN 'arn:aws:sso::987654321098:application/ssoins-98765a43bc987c2/apl-z9y8x7w6v5u4t3s2'
IAM_ROLE 'arn:aws:iam::987654321098:role/RedshiftFederatedRole';

-- Example 4: OAuth2 provider with minimal configuration
CREATE IDENTITY PROVIDER simple_oauth
TYPE OAUTH2
NAMESPACE 'oauth'
PARAMETERS '{"issuer":"https://login.example.com", "client_id":"client123"}';

-- Example 5: Provider with AUTO_CREATE_ROLES set to TRUE
CREATE IDENTITY PROVIDER auto_roles_provider
TYPE AWSIDC
NAMESPACE 'auto_sso'
APPLICATION_ARN 'arn:aws:sso::111111111111:application/ssoins-11111a11aa111a1/apl-1a1a1a1a1a1a1a1a'
IAM_ROLE 'arn:aws:iam::111111111111:role/AutoRole'
AUTO_CREATE_ROLES TRUE;

-- Example 6: Provider with AUTO_CREATE_ROLES TRUE and INCLUDE GROUPS filter
CREATE IDENTITY PROVIDER filtered_provider
TYPE AWSIDC
NAMESPACE 'filtered_sso'
APPLICATION_ARN 'arn:aws:sso::222222222222:application/ssoins-22222b22bb222b2/apl-2b2b2b2b2b2b2b2b'
IAM_ROLE 'arn:aws:iam::222222222222:role/FilteredRole'
AUTO_CREATE_ROLES TRUE INCLUDE GROUPS LIKE 'dev_%';

-- Example 7: Provider with AUTO_CREATE_ROLES TRUE and EXCLUDE GROUPS filter
CREATE IDENTITY PROVIDER exclusive_provider
TYPE AWSIDC
NAMESPACE 'exclusive_sso'
APPLICATION_ARN 'arn:aws:sso::333333333333:application/ssoins-33333c33cc333c3/apl-3c3c3c3c3c3c3c3c'
IAM_ROLE 'arn:aws:iam::333333333333:role/ExclusiveRole'
AUTO_CREATE_ROLES TRUE EXCLUDE GROUPS LIKE 'temp_%';

-- Example 8: Provider with AUTO_CREATE_ROLES set to FALSE
CREATE IDENTITY PROVIDER manual_provider
TYPE AWSIDC
NAMESPACE 'manual_sso'
APPLICATION_ARN 'arn:aws:sso::444444444444:application/ssoins-44444d44dd444d4/apl-4d4d4d4d4d4d4d4d'
IAM_ROLE 'arn:aws:iam::444444444444:role/ManualRole'
AUTO_CREATE_ROLES FALSE;

-- Example 9: Azure provider with complex JSON parameters
CREATE IDENTITY PROVIDER complex_azure
TYPE azure
NAMESPACE 'complex_aad'
PARAMETERS '{"issuer":"https://sts.windows.net/abcd1234-5678-90ef-ghij-klmnopqrstuv/", "client_id":"12345678-90ab-cdef-1234-567890abcdef", "client_secret":"secretvalue~with~special~characters", "audience":["https://analysis.windows.net/powerbi/connector/AmazonRedshift", "https://custom.app.example.com"], "additional_claims":{"custom_claim":"value"}}';

-- Example 10: Provider with underscores in name (unquoted)
CREATE IDENTITY PROVIDER my_company_provider
TYPE OAUTH2
NAMESPACE 'company_oauth'
PARAMETERS '{"issuer":"https://auth.company.com", "client_id":"company_client_123", "client_secret":"company_secret_456"}';

-- Example 11: Provider with hyphenated name (quoted)
CREATE IDENTITY PROVIDER "my-company-provider"
TYPE azure
NAMESPACE 'company-aad'
PARAMETERS '{"issuer":"https://login.microsoftonline.com/tenant-id/", "client_id":"azure-client-id"}';

-- Example 12: Minimal AWSIDC provider
CREATE IDENTITY PROVIDER minimal_idc
TYPE AWSIDC
NAMESPACE 'min'
APPLICATION_ARN 'arn:aws:sso::555555555555:application/ssoins-55555e55ee555e5/apl-5e5e5e5e5e5e5e5e'
IAM_ROLE 'arn:aws:iam::555555555555:role/MinimalRole';

-- Example 13: Provider with special characters in namespace
CREATE IDENTITY PROVIDER special_chars
TYPE OAUTH2
NAMESPACE 'oauth-2.0_provider'
PARAMETERS '{"issuer":"https://special-auth.example.com/oauth2/", "client_id":"special_client"}';

-- Example 14: Provider with empty parameters JSON
CREATE IDENTITY PROVIDER empty_params
TYPE OAUTH2
NAMESPACE 'empty'
PARAMETERS '{}';

-- Example 15: Provider with very long ARN and role
CREATE IDENTITY PROVIDER long_arn_provider
TYPE IAM_IDENTITY_CENTER
NAMESPACE 'long_namespace_for_testing_purposes'
APPLICATION_ARN 'arn:aws:sso::999999999999:application/ssoins-99999z99zz999z9/apl-9z9z9z9z9z9z9z9z'
IAM_ROLE 'arn:aws:iam::999999999999:role/VeryLongRoleNameForTestingPurposesWithManyCharacters';