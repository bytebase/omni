-- SHOW SCHEMAS test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_SCHEMAS.html

-- Basic Usage

-- Show all schemas from a specific database
SHOW SCHEMAS FROM DATABASE dev;

-- Show schemas from the default database
SHOW SCHEMAS FROM DATABASE mydb;

-- Show schemas from a test database
SHOW SCHEMAS FROM DATABASE test_db;

-- AWS Data Catalog Usage

-- Show schemas from AWS Glue Data Catalog
SHOW SCHEMAS FROM DATABASE awsdatacatalog;

-- Show schemas from AWS Glue Data Catalog with limit
SHOW SCHEMAS FROM DATABASE awsdatacatalog LIMIT 5;

-- Show schemas from AWS Glue Data Catalog with larger limit
SHOW SCHEMAS FROM DATABASE awsdatacatalog LIMIT 100;

-- Pattern Matching with LIKE

-- Show schemas that start with 'public'
SHOW SCHEMAS FROM DATABASE dev LIKE 'public%';

-- Show schemas that contain 'test'
SHOW SCHEMAS FROM DATABASE mydb LIKE '%test%';

-- Show schemas that end with '_schema'
SHOW SCHEMAS FROM DATABASE dev LIKE '%_schema';

-- Show schemas with single character wildcard
SHOW SCHEMAS FROM DATABASE test_db LIKE 'schema_';

-- Show schemas matching complex pattern
SHOW SCHEMAS FROM DATABASE dev LIKE 'app_%_v2';

-- Show schemas that start with 'pg_'
SHOW SCHEMAS FROM DATABASE dev LIKE 'pg_%';

-- Combining LIKE and LIMIT

-- Show up to 10 schemas that start with 'dev'
SHOW SCHEMAS FROM DATABASE mydb LIKE 'dev%' LIMIT 10;

-- Show up to 5 schemas containing 'analytics'
SHOW SCHEMAS FROM DATABASE awsdatacatalog LIKE '%analytics%' LIMIT 5;

-- Show up to 20 schemas ending with '_staging'
SHOW SCHEMAS FROM DATABASE dev LIKE '%_staging' LIMIT 20;

-- Edge Cases and Special Scenarios

-- Limit with zero rows (should return no results)
SHOW SCHEMAS FROM DATABASE dev LIMIT 0;

-- Maximum limit value
SHOW SCHEMAS FROM DATABASE dev LIMIT 10000;

-- Pattern with no wildcards (exact match)
SHOW SCHEMAS FROM DATABASE dev LIKE 'public';

-- Pattern with multiple wildcards
SHOW SCHEMAS FROM DATABASE dev LIKE '%test%schema%';

-- Pattern with underscore wildcard for single character
SHOW SCHEMAS FROM DATABASE dev LIKE 'schema_1';

-- Database names with special characters (quoted identifiers)
SHOW SCHEMAS FROM DATABASE "my-database";

-- Database names with mixed case (quoted identifiers)
SHOW SCHEMAS FROM DATABASE "MyDatabase";

-- Comments in queries

-- Show schemas with inline comment
SHOW SCHEMAS FROM DATABASE dev; -- This shows all schemas in dev

-- Show schemas with multi-line comment
/*
 * This query shows schemas from the AWS Data Catalog
 * with a limit of 10 rows
 */
SHOW SCHEMAS FROM DATABASE awsdatacatalog LIMIT 10;

-- Complex pattern matching examples

-- Show schemas with pattern containing both wildcards
SHOW SCHEMAS FROM DATABASE dev LIKE 'pg_%_temp';

-- Show schemas with escaped characters in pattern
SHOW SCHEMAS FROM DATABASE dev LIKE 'schema\_with\_underscores';

-- Show schemas with case-sensitive pattern matching
SHOW SCHEMAS FROM DATABASE dev LIKE 'Schema%';  -- Won't match 'schema...'

-- Additional test cases

-- Show schemas from catalog with specific patterns
SHOW SCHEMAS FROM DATABASE awsdatacatalog LIKE 'glue_%';

-- Show schemas with numeric patterns
SHOW SCHEMAS FROM DATABASE dev LIKE '%2023%';

-- Show schemas with very specific pattern
SHOW SCHEMAS FROM DATABASE dev LIKE 'app_v2_%_prod';

-- Limit edge cases
SHOW SCHEMAS FROM DATABASE dev LIMIT 1;
SHOW SCHEMAS FROM DATABASE dev LIMIT 9999;