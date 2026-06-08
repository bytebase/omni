-- ALTER DEFAULT PRIVILEGES examples from AWS Redshift documentation
-- https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_DEFAULT_PRIVILEGES.html

-- Grant SELECT on all future tables to a user group
ALTER DEFAULT PRIVILEGES
FOR USER report_admin
GRANT SELECT ON TABLES
TO GROUP report_readers;

-- Grant SELECT on all future tables to all users
ALTER DEFAULT PRIVILEGES
GRANT SELECT ON TABLES
TO PUBLIC;

-- Grant INSERT on future tables in specific schema
ALTER DEFAULT PRIVILEGES
IN SCHEMA sales
GRANT INSERT ON TABLES
TO GROUP sales_admin;

-- Grant multiple privileges to multiple users
ALTER DEFAULT PRIVILEGES
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES
TO user1, user2, user3;

-- Grant privileges with GRANT OPTION
ALTER DEFAULT PRIVILEGES
GRANT SELECT ON TABLES
TO admin_user WITH GRANT OPTION;

-- Grant privileges to a role
ALTER DEFAULT PRIVILEGES
GRANT SELECT, INSERT ON TABLES
TO ROLE data_analyst;

-- Grant privileges for specific user's future objects
ALTER DEFAULT PRIVILEGES
FOR USER data_owner
GRANT ALL PRIVILEGES ON TABLES
TO GROUP data_team;

-- Grant privileges in multiple schemas
ALTER DEFAULT PRIVILEGES
IN SCHEMA sales, marketing, finance
GRANT SELECT ON TABLES
TO GROUP analysts;

-- Revoke privileges
ALTER DEFAULT PRIVILEGES
IN SCHEMA sales
REVOKE INSERT ON TABLES
FROM GROUP sales_admin;

-- Revoke GRANT OPTION
ALTER DEFAULT PRIVILEGES
REVOKE GRANT OPTION FOR SELECT ON TABLES
FROM admin_user;

-- Revoke all privileges
ALTER DEFAULT PRIVILEGES
REVOKE ALL PRIVILEGES ON TABLES
FROM PUBLIC;

-- Grant EXECUTE on functions
ALTER DEFAULT PRIVILEGES
GRANT EXECUTE ON FUNCTIONS
TO GROUP dev_test;

-- Revoke EXECUTE on functions from PUBLIC
ALTER DEFAULT PRIVILEGES
REVOKE EXECUTE ON FUNCTIONS
FROM PUBLIC;

-- Grant EXECUTE on procedures with specific user
ALTER DEFAULT PRIVILEGES
FOR USER proc_owner
GRANT EXECUTE ON PROCEDURES
TO GROUP app_users;

-- Grant ALL on procedures in schema
ALTER DEFAULT PRIVILEGES
IN SCHEMA app_schema
GRANT ALL PRIVILEGES ON PROCEDURES
TO ROLE app_admin;

-- Complex example with multiple options
ALTER DEFAULT PRIVILEGES
FOR USER schema_owner
IN SCHEMA production, staging
GRANT SELECT, INSERT, UPDATE ON TABLES
TO GROUP prod_users, GROUP staging_users;

-- Grant REFERENCES privilege
ALTER DEFAULT PRIVILEGES
GRANT REFERENCES ON TABLES
TO user_with_fk_perms;

-- Grant TRUNCATE privilege
ALTER DEFAULT PRIVILEGES
IN SCHEMA cleanup_schema
GRANT TRUNCATE ON TABLES
TO GROUP maintenance_team;

-- Multiple users in FOR USER clause
ALTER DEFAULT PRIVILEGES
FOR USER user1, user2, user3
GRANT SELECT ON TABLES
TO PUBLIC;

-- Revoke specific privileges while keeping others
ALTER DEFAULT PRIVILEGES
REVOKE UPDATE, DELETE ON TABLES
FROM GROUP read_only_users;