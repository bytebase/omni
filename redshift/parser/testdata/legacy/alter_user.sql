-- ALTER USER test cases
-- Based on: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_USER.html

-- Basic ALTER USER with single option
ALTER USER admin CREATEDB;
ALTER USER dbuser NOCREATEDB;
ALTER USER admin CREATEUSER;
ALTER USER regular_user NOCREATEUSER;

-- SYSLOG ACCESS options
ALTER USER admin SYSLOG ACCESS UNRESTRICTED;
ALTER USER regular_user SYSLOG ACCESS RESTRICTED;

-- Password management
ALTER USER admin PASSWORD 'newSecurePass123!';
ALTER USER admin PASSWORD 'adminPass9' VALID UNTIL '2017-12-31 23:59';
ALTER USER dbuser PASSWORD DISABLE;

-- MD5 password hash
ALTER USER admin PASSWORD 'md5153c434b4b77c89e6b94f12c5393af5b';

-- SHA256 password hash with salt
ALTER USER admin PASSWORD 'sha256|4e2bf2c96c5cad3d65a85b5dc0d5236a933e8d1285f4beefabf769afe633b1b5|5df545aa401d27d9c5c2d2704299b8d1ba66f760a343fa4568fa70c584d05cee';

-- Rename user
ALTER USER admin RENAME TO sysadmin;

-- Connection limits
ALTER USER dbuser CONNECTION LIMIT 10;
ALTER USER admin CONNECTION LIMIT UNLIMITED;

-- Session timeout
ALTER USER dbuser SESSION TIMEOUT 300;
ALTER USER admin SESSION TIMEOUT 900;
ALTER USER dbuser RESET SESSION TIMEOUT;

-- SET and RESET parameters
ALTER USER admin SET search_path TO myschema, public;
ALTER USER dbuser SET statement_timeout = 60000;
ALTER USER admin SET enable_result_cache_for_session TO false;
ALTER USER dbuser RESET search_path;
ALTER USER admin RESET ALL;

-- External ID for identity provider
ALTER USER bob EXTERNALID "ABC123";
ALTER USER myco_aad:bob EXTERNALID "ABC123" PASSWORD DISABLE;
ALTER USER "mycompany.com:alice" EXTERNALID "XYZ456";

-- Multiple options with WITH keyword
ALTER USER admin WITH CREATEDB CREATEUSER;
ALTER USER dbuser WITH NOCREATEDB NOCREATEUSER SYSLOG ACCESS RESTRICTED;
ALTER USER newuser WITH PASSWORD 'tempPass123' VALID UNTIL '2024-12-31 23:59:59';

-- Multiple options without WITH keyword
ALTER USER admin CREATEDB CREATEUSER CONNECTION LIMIT UNLIMITED;
ALTER USER regular_user NOCREATEDB NOCREATEUSER CONNECTION LIMIT 5 SESSION TIMEOUT 600;
ALTER USER dbuser PASSWORD 'newPass456!' SYSLOG ACCESS RESTRICTED SET search_path = myschema;

-- Complex combinations
ALTER USER admin
    CREATEDB
    CREATEUSER
    SYSLOG ACCESS UNRESTRICTED
    PASSWORD 'SuperSecure!123'
    CONNECTION LIMIT UNLIMITED
    SESSION TIMEOUT 3600
    SET enable_result_cache_for_session = true;

ALTER USER readonly_user
    NOCREATEDB
    NOCREATEUSER
    SYSLOG ACCESS RESTRICTED
    CONNECTION LIMIT 3
    SESSION TIMEOUT 300
    SET statement_timeout TO 30000
    SET search_path TO public;

-- External identity provider with namespace
ALTER USER "mycompany.aad:john.doe" EXTERNALID "AAD-12345" PASSWORD DISABLE;
ALTER USER "aws:iam:123456789012:user/jane" EXTERNALID "IAM-67890";

-- Case sensitivity tests
alter user DBUSER createdb;
ALTER USER "MixedCase_User" NOCREATEDB;
ALTER USER "user.with.dots" PASSWORD 'test123';

-- Edge cases with special characters in usernames
ALTER USER "user-with-hyphens" CREATEDB;
ALTER USER "user_with_underscores" NOCREATEUSER;
ALTER USER "123_numeric_start" PASSWORD DISABLE;

-- Setting multiple parameters
ALTER USER admin
    SET search_path TO schema1, schema2, public
    SET statement_timeout = 120000
    SET enable_result_cache_for_session TO true;

-- Resetting multiple parameters
ALTER USER dbuser
    RESET search_path
    RESET statement_timeout
    PASSWORD 'newPassword123';

-- Valid until with different date formats
ALTER USER temp_user PASSWORD 'tempPass' VALID UNTIL '2024-12-31';
ALTER USER contractor PASSWORD 'contractPass' VALID UNTIL '2025-06-30 15:30:00';
ALTER USER seasonal_user PASSWORD 'seasonPass' VALID UNTIL '2024-08-31 23:59:59.999';

-- Connection limit edge cases
ALTER USER test_user CONNECTION LIMIT 0;
ALTER USER power_user CONNECTION LIMIT 1000;
ALTER USER special_user CONNECTION LIMIT 1;

-- Session timeout edge cases
ALTER USER quick_user SESSION TIMEOUT 60;
ALTER USER long_user SESSION TIMEOUT 86400;
ALTER USER standard_user SESSION TIMEOUT 1800;