-- ALTER DATABASE Test Cases for Amazon Redshift
-- Based on: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_DATABASE.html

-- Test 1: RENAME TO
-- Rename a database from tickit_sandbox to tickit_test
ALTER DATABASE tickit_sandbox RENAME TO tickit_test;

-- Rename a database with mixed case
ALTER DATABASE MyDatabase RENAME TO NewDatabaseName;

-- Test 2: OWNER TO
-- Change database owner to dwuser
ALTER DATABASE tickit OWNER TO dwuser;

-- Change owner with mixed case username
ALTER DATABASE analytics_db OWNER TO DataAnalyst;

-- Test 3: CONNECTION LIMIT
-- Set connection limit to specific number
ALTER DATABASE sampledb CONNECTION LIMIT 100;

-- Set connection limit to 1 (minimum practical limit)
ALTER DATABASE test_db CONNECTION LIMIT 1;

-- Set connection limit to unlimited
ALTER DATABASE production_db CONNECTION LIMIT UNLIMITED;

-- Test 4: COLLATE options
-- Set collation to case sensitive (full name)
ALTER DATABASE sampledb COLLATE CASE_SENSITIVE;

-- Set collation to case sensitive (abbreviation)
ALTER DATABASE reports_db COLLATE CS;

-- Set collation to case insensitive (full name)
ALTER DATABASE analytics COLLATE CASE_INSENSITIVE;

-- Set collation to case insensitive (abbreviation)
ALTER DATABASE staging_db COLLATE CI;

-- Test 5: Complex identifiers
-- Database names with underscores
ALTER DATABASE my_test_database RENAME TO my_production_database;

-- Database names with numbers
ALTER DATABASE db2024 OWNER TO admin123;

-- Test 6: Edge cases from documentation
-- Cannot rename system databases (these should be parsed but would fail at execution)
-- ALTER DATABASE dev RENAME TO dev_new;  -- Would fail: cannot rename dev
-- ALTER DATABASE template0 RENAME TO template0_new;  -- Would fail: cannot rename template0
-- ALTER DATABASE template1 RENAME TO template1_new;  -- Would fail: cannot rename template1

-- Test 7: Mixed case in options
ALTER DATABASE TestDB CONNECTION LIMIT 50;
ALTER DATABASE TestDB COLLATE CASE_INSENSITIVE;

-- Test 8: Real-world scenarios
-- Development to production promotion
ALTER DATABASE app_dev RENAME TO app_staging;

-- Security update - change owner
ALTER DATABASE customer_data OWNER TO security_admin;

-- Performance tuning - limit connections
ALTER DATABASE reporting CONNECTION LIMIT 25;

-- Data warehouse configuration
ALTER DATABASE warehouse COLLATE CI;