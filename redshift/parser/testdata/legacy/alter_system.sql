-- ALTER SYSTEM command test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_SYSTEM.html

-- Test 1: Enable AWS Glue Data Catalog automounting using TRUE
ALTER SYSTEM SET data_catalog_auto_mount = true;

-- Test 2: Enable AWS Glue Data Catalog automounting using T
ALTER SYSTEM SET data_catalog_auto_mount = t;

-- Test 3: Enable AWS Glue Data Catalog automounting using ON
ALTER SYSTEM SET data_catalog_auto_mount = on;

-- Test 4: Disable AWS Glue Data Catalog automounting using FALSE
ALTER SYSTEM SET data_catalog_auto_mount = false;

-- Test 5: Disable AWS Glue Data Catalog automounting using F
ALTER SYSTEM SET data_catalog_auto_mount = f;

-- Test 6: Disable AWS Glue Data Catalog automounting using OFF
ALTER SYSTEM SET data_catalog_auto_mount = off;

-- Test 7: Enable metadata security using TRUE
ALTER SYSTEM SET metadata_security = true;

-- Test 8: Enable metadata security using T
ALTER SYSTEM SET metadata_security = t;

-- Test 9: Enable metadata security using ON
ALTER SYSTEM SET metadata_security = on;

-- Test 10: Disable metadata security using FALSE
ALTER SYSTEM SET metadata_security = false;

-- Test 11: Disable metadata security using F
ALTER SYSTEM SET metadata_security = f;

-- Test 12: Disable metadata security using OFF
ALTER SYSTEM SET metadata_security = off;

-- Test 13: Set default identity namespace with simple name
ALTER SYSTEM SET default_identity_namespace = 'MYCO';

-- Test 14: Set default identity namespace with longer name
ALTER SYSTEM SET default_identity_namespace = 'MyCompanyIdentity';

-- Test 15: Set default identity namespace with underscore
ALTER SYSTEM SET default_identity_namespace = 'MY_COMPANY';

-- Test 16: Set default identity namespace with numbers
ALTER SYSTEM SET default_identity_namespace = 'Company123';

-- Test 17: Mixed case boolean values (case-insensitive)
ALTER SYSTEM SET data_catalog_auto_mount = TRUE;
ALTER SYSTEM SET data_catalog_auto_mount = True;
ALTER SYSTEM SET data_catalog_auto_mount = ON;
ALTER SYSTEM SET data_catalog_auto_mount = On;

-- Test 18: Mixed case for FALSE values
ALTER SYSTEM SET metadata_security = FALSE;
ALTER SYSTEM SET metadata_security = False;
ALTER SYSTEM SET metadata_security = OFF;
ALTER SYSTEM SET metadata_security = Off;

-- Test 19: Single character boolean values with different cases
ALTER SYSTEM SET data_catalog_auto_mount = T;
ALTER SYSTEM SET data_catalog_auto_mount = F;
ALTER SYSTEM SET metadata_security = t;
ALTER SYSTEM SET metadata_security = f;

-- Test 20: Real-world scenario - Enable catalog automounting and metadata security
ALTER SYSTEM SET data_catalog_auto_mount = true;
ALTER SYSTEM SET metadata_security = true;

-- Test 21: Real-world scenario - Set up identity namespace and enable security
ALTER SYSTEM SET default_identity_namespace = 'ACME_CORP';
ALTER SYSTEM SET metadata_security = on;

-- Test 22: Empty string for identity namespace (edge case)
ALTER SYSTEM SET default_identity_namespace = '';

-- Test 23: Identity namespace with special characters
ALTER SYSTEM SET default_identity_namespace = 'Company-Name';
ALTER SYSTEM SET default_identity_namespace = 'Company.Name';

-- Note: The following are expected to be invalid and should be tested for error handling:
-- ALTER SYSTEM SET invalid_parameter = true;
-- ALTER SYSTEM SET data_catalog_auto_mount = 'invalid';
-- ALTER SYSTEM SET default_identity_namespace = true;
-- ALTER SYSTEM SET metadata_security = 'MYCO';