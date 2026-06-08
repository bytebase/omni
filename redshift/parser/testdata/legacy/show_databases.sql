-- SHOW DATABASES test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_DATABASES.html

-- Basic Usage

-- Show all databases in the current data warehouse
SHOW DATABASES;

-- Show databases with a specific limit
SHOW DATABASES LIMIT 10;

-- Show databases with a smaller limit
SHOW DATABASES LIMIT 5;

-- Pattern Matching with LIKE

-- Show databases that start with 'dev'
SHOW DATABASES LIKE 'dev%';

-- Show databases that contain 'test'
SHOW DATABASES LIKE '%test%';

-- Show databases that end with '_db'
SHOW DATABASES LIKE '%_db';

-- Show databases with single character wildcard
SHOW DATABASES LIKE 'db_';

-- Show databases matching complex pattern
SHOW DATABASES LIKE 'prod_%_2023';

-- Combine LIKE with LIMIT
SHOW DATABASES LIKE 'dev%' LIMIT 5;

-- Data Catalog Usage

-- Show databases from data catalog with default settings
SHOW DATABASES FROM DATA CATALOG;

-- Show databases from data catalog with specific account
SHOW DATABASES FROM DATA CATALOG ACCOUNT '123456789012';

-- Show databases from multiple accounts
SHOW DATABASES FROM DATA CATALOG ACCOUNT '123456789012', '234567890123', '345678901234';

-- Show databases from data catalog with pattern matching
SHOW DATABASES FROM DATA CATALOG LIKE 'analytics%';

-- Show databases from data catalog with limit
SHOW DATABASES FROM DATA CATALOG LIMIT 20;

-- IAM Role Options

-- Use default IAM role
SHOW DATABASES FROM DATA CATALOG IAM_ROLE DEFAULT;

-- Use session credentials
SHOW DATABASES FROM DATA CATALOG IAM_ROLE SESSION;

-- Use specific IAM role ARN
SHOW DATABASES FROM DATA CATALOG IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftDataCatalogRole';

-- Complex Combinations

-- Data catalog with account and IAM role
SHOW DATABASES FROM DATA CATALOG ACCOUNT '123456789012' IAM_ROLE DEFAULT;

-- Data catalog with account, pattern, and limit
SHOW DATABASES FROM DATA CATALOG ACCOUNT '123456789012' LIKE 'prod%' LIMIT 10;

-- Data catalog with multiple accounts and pattern
SHOW DATABASES FROM DATA CATALOG ACCOUNT '123456789012', '234567890123' LIKE 'staging_%';

-- Full syntax with all options
SHOW DATABASES FROM DATA CATALOG ACCOUNT '123456789012', '234567890123' LIKE 'analytics%' IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftRole' LIMIT 50;

-- Edge Cases and Special Patterns

-- Show databases with exact name match (no wildcards)
SHOW DATABASES LIKE 'mydb';

-- Show all databases (equivalent to no LIKE clause)
SHOW DATABASES LIKE '%';

-- Empty result scenarios
SHOW DATABASES LIKE 'nonexistent%';

-- Maximum limit
SHOW DATABASES LIMIT 10000;

-- Minimum limit
SHOW DATABASES LIMIT 1;

-- Zero limit (should return no results)
SHOW DATABASES LIMIT 0;

-- Pattern with underscores (literal underscore in database name)
SHOW DATABASES LIKE 'my\_database';

-- Pattern with percentage sign (literal % in database name)
SHOW DATABASES LIKE 'db\%prod';

-- Multiple wildcards in pattern
SHOW DATABASES LIKE '%dev%test%';

-- Case sensitivity tests (Redshift is case-insensitive by default)
SHOW DATABASES LIKE 'DEV%';
SHOW DATABASES LIKE 'Dev%';
SHOW DATABASES LIKE 'dev%';

-- Quoted identifiers in LIKE patterns
SHOW DATABASES LIKE '"MyDatabase"';

-- Special characters in patterns
SHOW DATABASES LIKE 'db-prod-%';
SHOW DATABASES LIKE 'db.test.%';
SHOW DATABASES LIKE 'db$prod%';

-- IAM Role variations
SHOW DATABASES FROM DATA CATALOG IAM_ROLE 'SESSION';
SHOW DATABASES FROM DATA CATALOG IAM_ROLE 'default';
SHOW DATABASES FROM DATA CATALOG IAM_ROLE 'arn:aws:iam::123456789012:role/service-role/AmazonRedshift-CommandsAccessRole-20231201';

-- Multiple account IDs with different formats
SHOW DATABASES FROM DATA CATALOG ACCOUNT '123456789012', '000000000000', '999999999999';

-- All options combined in different orders
SHOW DATABASES FROM DATA CATALOG IAM_ROLE DEFAULT ACCOUNT '123456789012' LIKE 'prod%' LIMIT 25;
SHOW DATABASES FROM DATA CATALOG ACCOUNT '123456789012' IAM_ROLE 'arn:aws:iam::123456789012:role/RedshiftRole' LIKE 'test%' LIMIT 15;