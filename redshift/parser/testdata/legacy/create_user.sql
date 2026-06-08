-- CREATE USER test cases based on AWS Redshift documentation
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_USER.html

-- Basic user creation (from AWS documentation)
CREATE USER dbuser WITH PASSWORD 'abcD1234';

-- User creation without WITH keyword
CREATE USER dbuser PASSWORD 'abcD1234';

-- User with connection limits and database creation privileges
CREATE USER dbuser WITH PASSWORD 'abcD1234' CREATEDB CONNECTION LIMIT 30;

-- User with session timeout
CREATE USER dbuser PASSWORD 'abcD1234' SESSION TIMEOUT 120;

-- External identity provider user (from AWS documentation)
CREATE USER "myco_aad:bob" EXTERNALID "ABC123" PASSWORD DISABLE;

-- User with multiple privileges
CREATE USER admin_user WITH PASSWORD 'SecurePass123' CREATEDB CREATEUSER;

-- User with syslog access
CREATE USER syslog_user WITH PASSWORD 'LogPass456' SYSLOG ACCESS RESTRICTED;
CREATE USER admin_syslog WITH PASSWORD 'AdminLog789' SYSLOG ACCESS UNRESTRICTED;

-- User added to group
CREATE USER group_user WITH PASSWORD 'GroupPass321' IN GROUP analysts;

-- User with expiration date
CREATE USER temp_user WITH PASSWORD 'TempPass654' VALID UNTIL '2024-12-31 23:59:59';

-- User with unlimited connections
CREATE USER unlimited_user WITH PASSWORD 'UnlimitedPass987' CONNECTION LIMIT UNLIMITED;

-- MD5 password hash (from AWS documentation example)
CREATE USER md5_user WITH PASSWORD 'md5f6fdffe48c908deb0f4c3bd36c032e72d';

-- SHA-256 password hash
CREATE USER sha256_user WITH PASSWORD 'sha256|8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92';

-- Password disabled for IAM authentication
CREATE USER iam_user PASSWORD DISABLE;
CREATE USER federated_user EXTERNALID "arn:aws:iam::123456789012:role/RedshiftUser" PASSWORD DISABLE;

-- Complex user with multiple options
CREATE USER complex_user WITH PASSWORD 'ComplexPass123'
    CREATEDB
    NOCREATEUSER
    SYSLOG ACCESS RESTRICTED
    CONNECTION LIMIT 10
    SESSION TIMEOUT 600
    VALID UNTIL '2025-01-01';

-- User names with special characters (quoted)
CREATE USER "user-with-dash" WITH PASSWORD 'DashPass123';
CREATE USER "user.with.dots" WITH PASSWORD 'DotPass456';
CREATE USER "user with spaces" WITH PASSWORD 'SpacePass789';
CREATE USER "UPPERCASE_USER" WITH PASSWORD 'UpperPass321';

-- External provider format usernames (quoted for special characters)
CREATE USER "okta:john.doe" WITH PASSWORD 'OktaPass123';
CREATE USER "azure_ad:jane.smith" WITH PASSWORD 'AzurePass456';
CREATE USER "google:bob.wilson" WITH PASSWORD 'GooglePass789';

-- Numeric usernames (quoted)
CREATE USER "123user" WITH PASSWORD 'NumericUser123';
CREATE USER "user123" WITH PASSWORD 'UserNumeric456';

-- Single character usernames
CREATE USER a WITH PASSWORD 'SingleChar123';
CREATE USER "A" WITH PASSWORD 'SingleCharUpper456';

-- Real-world service account patterns
CREATE USER etl_service WITH PASSWORD 'ETLService123' NOCREATEDB NOCREATEUSER CONNECTION LIMIT 5;
CREATE USER reporting_svc WITH PASSWORD 'ReportSvc456' SYSLOG ACCESS RESTRICTED SESSION TIMEOUT 1800;
CREATE USER analytics_user WITH PASSWORD 'Analytics789' IN GROUP data_analysts CREATEDB;

-- Development environment users
CREATE USER dev_user WITH PASSWORD 'DevPass123' CREATEDB CONNECTION LIMIT 20;
CREATE USER test_user WITH PASSWORD 'TestPass456' VALID UNTIL '2024-06-30';
CREATE USER staging_user WITH PASSWORD 'StagingPass789' SESSION TIMEOUT 3600;

-- Administrative users
CREATE USER dba_user WITH PASSWORD 'DBAPass123' CREATEDB CREATEUSER SYSLOG ACCESS UNRESTRICTED;
CREATE USER backup_user WITH PASSWORD 'BackupPass456' NOCREATEDB NOCREATEUSER CONNECTION LIMIT 2;
CREATE USER monitor_user WITH PASSWORD 'MonitorPass789' SYSLOG ACCESS RESTRICTED;

-- External identity variations
CREATE USER sso_user EXTERNALID "sso-provider-id-12345" PASSWORD DISABLE;
CREATE USER saml_user EXTERNALID "urn:amazon:webservices:govcloud" PASSWORD DISABLE;
CREATE USER oidc_user EXTERNALID "https://example.com/.well-known/openid_configuration" PASSWORD DISABLE;

-- AWS IAM role patterns
CREATE USER aws_lambda EXTERNALID "arn:aws:iam::123456789012:role/LambdaExecutionRole" PASSWORD DISABLE;
CREATE USER cross_account EXTERNALID "arn:aws:iam::987654321098:role/CrossAccountRole" PASSWORD DISABLE;
CREATE USER service_role EXTERNALID "arn:aws:iam::555666777888:role/service-role/MyServiceRole" PASSWORD DISABLE;

-- Azure AD patterns
CREATE USER azure_user EXTERNALID "azure-ad-group-12345" PASSWORD DISABLE;
CREATE USER azure_app EXTERNALID "00000000-0000-0000-0000-000000000000" PASSWORD DISABLE;

-- Google Cloud patterns
CREATE USER gcp_user EXTERNALID "gcp-service-account@project.iam.gserviceaccount.com" PASSWORD DISABLE;

-- Okta patterns
CREATE USER okta_group EXTERNALID "00g1emaKYMLWCQXVHDSN" PASSWORD DISABLE;
CREATE USER okta_app EXTERNALID "0oa2hm8bXXXXXXXXXX" PASSWORD DISABLE;

-- Mixed case external IDs
CREATE USER mixed_case EXTERNALID "MixedCaseExternalId123" PASSWORD DISABLE;
CREATE USER special_chars EXTERNALID "external-id_with.special@chars.com" PASSWORD DISABLE;

-- Empty and minimal values
CREATE USER min_user PASSWORD 'MinPass12';
CREATE USER "a" PASSWORD 'SingleA12';

-- Complex passwords meeting requirements
CREATE USER secure_user1 WITH PASSWORD 'Abcd1234!@#$';
CREATE USER secure_user2 WITH PASSWORD 'MySecureP@ssw0rd2024';
CREATE USER secure_user3 WITH PASSWORD 'C0mpl3x&P@ssw0rd!';

-- Connection limit variations
CREATE USER limit_1 WITH PASSWORD 'LimitOne123' CONNECTION LIMIT 1;
CREATE USER limit_100 WITH PASSWORD 'LimitHund456' CONNECTION LIMIT 100;
CREATE USER no_limit WITH PASSWORD 'NoLimit789' CONNECTION LIMIT UNLIMITED;

-- Session timeout variations
CREATE USER timeout_60 WITH PASSWORD 'Timeout60sec' SESSION TIMEOUT 60;
CREATE USER timeout_3600 WITH PASSWORD 'Timeout1hr' SESSION TIMEOUT 3600;
CREATE USER timeout_86400 WITH PASSWORD 'Timeout1day' SESSION TIMEOUT 86400;

-- Group membership examples
CREATE USER analyst1 WITH PASSWORD 'Analyst123' IN GROUP readonly_users;
CREATE USER analyst2 WITH PASSWORD 'Analyst456' IN GROUP data_scientists;
CREATE USER admin1 WITH PASSWORD 'Admin123' IN GROUP administrators;

-- Expiration date variations
CREATE USER expire_2024 WITH PASSWORD 'Expire2024' VALID UNTIL '2024-12-31';
CREATE USER expire_specific WITH PASSWORD 'ExpireSpec' VALID UNTIL '2024-06-15 14:30:00';
CREATE USER expire_eoy WITH PASSWORD 'ExpireEOY' VALID UNTIL '2024-12-31 23:59:59';

-- Multiple options in different orders
CREATE USER ordered1 WITH PASSWORD 'Order123' CREATEDB CONNECTION LIMIT 15 SESSION TIMEOUT 1200;
CREATE USER ordered2 WITH PASSWORD 'Order456' SESSION TIMEOUT 900 CREATEUSER SYSLOG ACCESS RESTRICTED;
CREATE USER ordered3 WITH PASSWORD 'Order789' VALID UNTIL '2025-01-01' IN GROUP power_users NOCREATEDB;

-- External ID with various quote styles and special characters
CREATE USER quote_test1 EXTERNALID "external-id-with-quotes" PASSWORD DISABLE;
CREATE USER quote_test2 EXTERNALID "external 'id' with internal quotes" PASSWORD DISABLE;
CREATE USER unicode_test EXTERNALID "üñíçødé_éxtérnål_íd" PASSWORD DISABLE;

-- Realistic scenario combinations
CREATE USER data_engineer WITH PASSWORD 'DataEng123'
    CREATEDB
    IN GROUP engineers
    CONNECTION LIMIT 25
    SESSION TIMEOUT 7200
    SYSLOG ACCESS RESTRICTED;

CREATE USER readonly_analyst WITH PASSWORD 'ReadOnly456'
    NOCREATEDB
    NOCREATEUSER
    IN GROUP analysts
    CONNECTION LIMIT 5
    SESSION TIMEOUT 3600;

CREATE USER temp_contractor WITH PASSWORD 'TempWork789'
    NOCREATEDB
    NOCREATEUSER
    CONNECTION LIMIT 3
    VALID UNTIL '2024-03-31 23:59:59';

CREATE USER service_account WITH PASSWORD 'ServiceAcc321'
    NOCREATEDB
    NOCREATEUSER
    CONNECTION LIMIT 10
    SESSION TIMEOUT 1800;

-- Edge cases and boundary testing
CREATE USER edge_case1 WITH PASSWORD 'EdgeCase1' CONNECTION LIMIT 0;  -- Zero connections
CREATE USER edge_case2 WITH PASSWORD 'EdgeCase2' SESSION TIMEOUT 1;   -- One second timeout

-- SQL keywords as quoted identifiers
CREATE USER "SELECT" WITH PASSWORD 'SelectUser123';
CREATE USER "CREATE" WITH PASSWORD 'CreateUser456';
CREATE USER "TABLE" WITH PASSWORD 'TableUser789';
CREATE USER "DATABASE" WITH PASSWORD 'DatabaseUser321';
CREATE USER "PASSWORD" WITH PASSWORD 'PasswordUser654';

-- Case sensitivity tests for usernames
CREATE USER lowercase_user WITH PASSWORD 'LowerCase123';
CREATE USER UPPERCASE_USER WITH PASSWORD 'UpperCase456';
CREATE USER "CamelCaseUser" WITH PASSWORD 'CamelCase789';

-- Testing all password options in isolation
CREATE USER clear_pass WITH PASSWORD 'ClearPassword123';
CREATE USER md5_pass WITH PASSWORD 'md5098f6bcd4621d373cade4e832627b4f6';
CREATE USER sha256_pass WITH PASSWORD 'sha256|e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855';
CREATE USER disabled_pass PASSWORD DISABLE;