-- ALTER SCHEMA test cases
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_SCHEMA.html

-- Basic rename schema
ALTER SCHEMA sales RENAME TO us_sales;

-- Rename with quoted identifiers
ALTER SCHEMA "my_schema" RENAME TO "new_schema";

-- Rename with mixed case identifiers
ALTER SCHEMA MySchema RENAME TO YourSchema;

-- Change schema owner
ALTER SCHEMA us_sales OWNER TO dwuser;

-- Change owner with quoted identifier
ALTER SCHEMA "production" OWNER TO "app_user";

-- Set schema quota with GB (default unit)
ALTER SCHEMA us_sales QUOTA 300 GB;

-- Set schema quota with MB
ALTER SCHEMA test_schema QUOTA 500 MB;

-- Set schema quota with TB
ALTER SCHEMA warehouse QUOTA 2 TB;

-- Set schema quota without unit (defaults to GB)
ALTER SCHEMA analytics QUOTA 100;

-- Set unlimited quota
ALTER SCHEMA us_sales QUOTA UNLIMITED;

-- Complex examples with different identifier types
ALTER SCHEMA public OWNER TO admin_user;
ALTER SCHEMA "My-Schema-123" RENAME TO "My_Schema_456";
ALTER SCHEMA data_warehouse QUOTA 1500 GB;

-- Edge cases
ALTER SCHEMA a RENAME TO b;
ALTER SCHEMA schema123 OWNER TO user456;
ALTER SCHEMA _underscore_schema QUOTA 50 MB;

-- Comments in statements
ALTER SCHEMA old_schema -- This schema is being renamed
    RENAME TO new_schema;

ALTER SCHEMA production /* Change owner for maintenance */
    OWNER TO maintenance_user;

-- Multi-line formatting
ALTER SCHEMA
    sales_data
RENAME TO
    historical_sales_data;

ALTER SCHEMA
    reporting
QUOTA
    750 GB;

-- Schema names with numbers
ALTER SCHEMA schema2023 RENAME TO schema2024;
ALTER SCHEMA data_v1 OWNER TO data_admin;
ALTER SCHEMA temp_123 QUOTA UNLIMITED;

-- Reserved keywords as identifiers (when quoted)
ALTER SCHEMA "alter" RENAME TO "update";
ALTER SCHEMA "schema" OWNER TO "user";
ALTER SCHEMA "quota" QUOTA 100 GB;