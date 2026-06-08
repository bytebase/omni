-- SHOW DATASHARES test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_DATASHARES.html

-- Basic SHOW DATASHARES command
-- Shows all datashares in the cluster
SHOW DATASHARES;

-- With LIKE pattern - exact match
SHOW DATASHARES LIKE 'salesshare';

-- With LIKE pattern - wildcard at end
-- Shows all datashares starting with 'sales'
SHOW DATASHARES LIKE 'sales%';

-- With LIKE pattern - wildcard at beginning
-- Shows all datashares ending with 'share'
SHOW DATASHARES LIKE '%share';

-- With LIKE pattern - wildcard in middle
-- Shows all datashares containing 'data'
SHOW DATASHARES LIKE '%data%';

-- With LIKE pattern - single character wildcard
-- Shows datashares matching pattern like 'share1', 'share2', etc.
SHOW DATASHARES LIKE 'share_';

-- Complex LIKE patterns
SHOW DATASHARES LIKE 'prod_%_share';
SHOW DATASHARES LIKE 'test_datashare_%';
SHOW DATASHARES LIKE '_____share'; -- Exactly 5 characters before 'share'

-- Case variations (SQL keywords are case-insensitive)
show datashares;
Show Datashares;
SHOW datashares;
show DATASHARES;

-- With LIKE in different cases
show datashares like 'myshare%';
SHOW DATASHARES Like 'test%';
Show Datashares LIKE 'prod%';

-- Edge cases with special characters in patterns
SHOW DATASHARES LIKE 'share$%';
SHOW DATASHARES LIKE 'share@test%';
SHOW DATASHARES LIKE 'share-prod%';
SHOW DATASHARES LIKE 'share.dev%';

-- Empty pattern
SHOW DATASHARES LIKE '';

-- Pattern with only wildcards
SHOW DATASHARES LIKE '%';
SHOW DATASHARES LIKE '%%';
SHOW DATASHARES LIKE '_%';

-- Patterns with escaped characters
SHOW DATASHARES LIKE 'share\_prod%'; -- Looking for literal underscore
SHOW DATASHARES LIKE 'share\%test';  -- Looking for literal percent

-- Multi-line formatting
SHOW
    DATASHARES;

SHOW DATASHARES
    LIKE 'sales%';

SHOW
    DATASHARES
    LIKE
    'prod%';

-- With comments
SHOW DATASHARES; -- Show all datashares
SHOW DATASHARES LIKE 'test%'; -- Show test datashares
/* Show production datashares */ SHOW DATASHARES LIKE 'prod%';

-- Common naming patterns
SHOW DATASHARES LIKE 'dev_%';
SHOW DATASHARES LIKE 'staging_%';
SHOW DATASHARES LIKE 'prod_%';
SHOW DATASHARES LIKE '%_backup';
SHOW DATASHARES LIKE '%_archive';
SHOW DATASHARES LIKE '%_replica';

-- Date-based patterns
SHOW DATASHARES LIKE '%2024%';
SHOW DATASHARES LIKE '%202401%';
SHOW DATASHARES LIKE '%20240115%';

-- Department/team patterns
SHOW DATASHARES LIKE 'finance_%';
SHOW DATASHARES LIKE 'marketing_%';
SHOW DATASHARES LIKE 'analytics_%';
SHOW DATASHARES LIKE 'engineering_%';

-- Region-based patterns
SHOW DATASHARES LIKE '%_us_east_%';
SHOW DATASHARES LIKE '%_eu_west_%';
SHOW DATASHARES LIKE '%_ap_south_%';

-- Version patterns
SHOW DATASHARES LIKE '%_v1';
SHOW DATASHARES LIKE '%_v2';
SHOW DATASHARES LIKE '%_v1.0';
SHOW DATASHARES LIKE '%_v2.1';

-- Complex real-world examples
SHOW DATASHARES LIKE 'sales_prod_us_east_%';
SHOW DATASHARES LIKE 'analytics_staging_v2_%';
SHOW DATASHARES LIKE 'finance_archive_2024%';
SHOW DATASHARES LIKE 'marketing_dev_%_backup';