-- Basic ALTER ROLE statements

-- Rename a role
ALTER ROLE sample_role1 RENAME TO sample_role2;
ALTER ROLE old_analytics_role RENAME TO new_analytics_role;
ALTER ROLE temp_role RENAME TO permanent_role;

-- Change role owner
ALTER ROLE sample_role1 OWNER TO user1;
ALTER ROLE sample_role1 WITH OWNER TO user1;
ALTER ROLE analytics_role OWNER TO admin_user;
ALTER ROLE data_role WITH OWNER TO data_admin;

-- Update external ID for identity provider
ALTER ROLE sample_role1 EXTERNALID TO "XYZ456";
ALTER ROLE sample_role1 WITH EXTERNALID TO "XYZ456";
ALTER ROLE federated_role EXTERNALID TO "ABC123DEF456";
ALTER ROLE sso_role WITH EXTERNALID TO "SSO_ID_12345";

-- Multiple operations require separate statements (per AWS docs)
-- Cannot combine RENAME TO with OWNER TO in a single statement
ALTER ROLE sample_role1 RENAME TO sample_role2;
ALTER ROLE sample_role2 OWNER TO user2;

-- Can combine OWNER TO with EXTERNALID
ALTER ROLE role1 OWNER TO admin;
ALTER ROLE role1 EXTERNALID TO "EXT123";

-- Sequential operations
ALTER ROLE test_role RENAME TO prod_role;
ALTER ROLE prod_role OWNER TO prod_admin;
ALTER ROLE prod_role EXTERNALID TO "PROD_ID";

-- Complex role names
ALTER ROLE "MyRole" RENAME TO "YourRole";
ALTER ROLE role_with_underscore RENAME TO role_without_underscore;
ALTER ROLE "role-with-dash" RENAME TO "role_with_underscore";

-- Various user names in OWNER TO
ALTER ROLE role1 OWNER TO "john.doe";
ALTER ROLE role1 OWNER TO user_with_underscore;
ALTER ROLE role1 OWNER TO "user-with-dash";

-- External IDs with different formats
ALTER ROLE role1 EXTERNALID TO "simple";
ALTER ROLE role1 EXTERNALID TO "with spaces and special chars!@#$%";
ALTER ROLE role1 EXTERNALID TO "123456789";
ALTER ROLE role1 EXTERNALID TO "arn:aws:iam::123456789012:role/MyRole";
ALTER ROLE role1 EXTERNALID TO "";  -- Empty external ID

-- WITH keyword usage variations
ALTER ROLE role1 WITH RENAME TO role2;
ALTER ROLE role1 RENAME TO role2;  -- WITH is optional

-- Real-world examples
ALTER ROLE etl_reader_role RENAME TO etl_full_access_role;
ALTER ROLE bi_analyst_role OWNER TO bi_admin_user;
ALTER ROLE aws_federated_role EXTERNALID TO "arn:aws:iam::987654321098:role/RedshiftAccess";

-- Edge cases
ALTER ROLE a RENAME TO b;  -- Single character role names
ALTER ROLE role123 RENAME TO role456;  -- Numeric suffixes
ALTER ROLE ROLE RENAME TO NEWROLE;  -- Keyword as identifier (if allowed)
ALTER ROLE "SELECT" RENAME TO "INSERT";  -- SQL keywords as quoted identifiers