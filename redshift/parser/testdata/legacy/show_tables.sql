-- SHOW TABLES test cases based on AWS documentation
-- https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_TABLES.html

-- Basic syntax: Show all tables in a schema
SHOW TABLES FROM SCHEMA dev.public;

-- Show tables in different schemas
SHOW TABLES FROM SCHEMA mydb.myschema;
SHOW TABLES FROM SCHEMA production.sales;
SHOW TABLES FROM SCHEMA analytics.reports;

-- Show tables in AWS Glue Data Catalog
SHOW TABLES FROM SCHEMA awsdatacatalog.batman;
SHOW TABLES FROM SCHEMA awsdatacatalog.customer_data;
SHOW TABLES FROM SCHEMA awsdatacatalog.product_catalog;

-- Using LIKE pattern matching with wildcards
-- % matches zero or more characters
SHOW TABLES FROM SCHEMA dev.public LIKE 'user%';
SHOW TABLES FROM SCHEMA dev.public LIKE '%_temp';
SHOW TABLES FROM SCHEMA dev.public LIKE 'log_%';
SHOW TABLES FROM SCHEMA dev.public LIKE '%customer%';

-- Using LIKE with underscore wildcard
-- _ matches exactly one character
SHOW TABLES FROM SCHEMA dev.public LIKE 'user_';
SHOW TABLES FROM SCHEMA dev.public LIKE 't_mp';
SHOW TABLES FROM SCHEMA dev.public LIKE '___data';

-- Complex LIKE patterns
SHOW TABLES FROM SCHEMA dev.public LIKE 'fact_%_2024';
SHOW TABLES FROM SCHEMA dev.public LIKE 'dim_%_v%';
SHOW TABLES FROM SCHEMA dev.public LIKE '%_backup_%';

-- Using LIMIT clause
SHOW TABLES FROM SCHEMA dev.public LIMIT 10;
SHOW TABLES FROM SCHEMA dev.public LIMIT 100;
SHOW TABLES FROM SCHEMA dev.public LIMIT 1000;
SHOW TABLES FROM SCHEMA dev.public LIMIT 10000;

-- Combining LIKE and LIMIT
SHOW TABLES FROM SCHEMA dev.public LIKE 'user%' LIMIT 50;
SHOW TABLES FROM SCHEMA dev.public LIKE '%_temp' LIMIT 20;
SHOW TABLES FROM SCHEMA dev.public LIKE 'fact_%' LIMIT 100;

-- AWS Glue Data Catalog with filters
SHOW TABLES FROM SCHEMA awsdatacatalog.sales_db LIKE 'orders%' LIMIT 25;
SHOW TABLES FROM SCHEMA awsdatacatalog.analytics LIKE '%_daily' LIMIT 50;

-- Edge cases with special characters in identifiers
SHOW TABLES FROM SCHEMA "my-db"."my-schema";
SHOW TABLES FROM SCHEMA "DB_2024"."SCHEMA_PROD";

-- Mixed case identifiers
SHOW TABLES FROM SCHEMA MyDatabase.MySchema;
SHOW TABLES FROM SCHEMA MYDB.MYSCHEMA;

-- Pattern matching edge cases
SHOW TABLES FROM SCHEMA dev.public LIKE '';  -- Empty pattern
SHOW TABLES FROM SCHEMA dev.public LIKE '%';  -- Match all
SHOW TABLES FROM SCHEMA dev.public LIKE '_%';  -- At least one character
SHOW TABLES FROM SCHEMA dev.public LIKE '%_%';  -- Contains at least one character

-- Numeric limits
SHOW TABLES FROM SCHEMA dev.public LIMIT 0;  -- Minimum limit
SHOW TABLES FROM SCHEMA dev.public LIMIT 1;  -- Single result

-- Complex real-world examples
SHOW TABLES FROM SCHEMA production.analytics LIKE 'fact_%_daily' LIMIT 100;
SHOW TABLES FROM SCHEMA staging.etl LIKE 'tmp_%_2024%' LIMIT 50;
SHOW TABLES FROM SCHEMA warehouse.reporting LIKE 'rpt_%_v%' LIMIT 200;

-- Comments in queries
-- Show user tables
SHOW TABLES FROM SCHEMA dev.public LIKE 'user%'; -- User-related tables

/* Multi-line comment
   Showing temporary tables */
SHOW TABLES FROM SCHEMA dev.public LIKE '%_temp' LIMIT 10;