-- DROP SCHEMA Test Cases for Amazon Redshift
-- Based on AWS documentation: https://docs.aws.amazon.com/redshift/latest/dg/r_DROP_SCHEMA.html

-- Basic DROP SCHEMA with RESTRICT (default behavior)
DROP SCHEMA s_sales RESTRICT;

-- DROP SCHEMA with CASCADE to drop all objects in schema
DROP SCHEMA s_sales CASCADE;

-- DROP SCHEMA with IF EXISTS to prevent error if schema doesn't exist
DROP SCHEMA IF EXISTS s_sales;

-- DROP SCHEMA with external database clause
DROP SCHEMA s_spectrum DROP EXTERNAL DATABASE RESTRICT;

-- DROP multiple schemas with external database and CASCADE
DROP SCHEMA s_sales, s_profit, s_revenue DROP EXTERNAL DATABASE CASCADE;

-- Basic DROP SCHEMA without explicit CASCADE/RESTRICT (uses RESTRICT by default)
DROP SCHEMA test_schema;

-- DROP SCHEMA IF EXISTS with CASCADE
DROP SCHEMA IF EXISTS old_schema CASCADE;

-- DROP SCHEMA IF EXISTS with RESTRICT
DROP SCHEMA IF EXISTS backup_schema RESTRICT;

-- DROP multiple schemas with IF EXISTS
DROP SCHEMA IF EXISTS schema1, schema2, schema3;

-- DROP single external schema
DROP SCHEMA external_data DROP EXTERNAL DATABASE;

-- DROP SCHEMA with quoted identifier
DROP SCHEMA "Special-Schema";

-- DROP multiple schemas with mixed quoted/unquoted identifiers
DROP SCHEMA regular_schema, "quoted-schema", another_schema;

-- DROP SCHEMA with IF EXISTS and external database
DROP SCHEMA IF EXISTS spectrum_schema DROP EXTERNAL DATABASE CASCADE;

-- Complex case: multiple schemas with all options
DROP SCHEMA IF EXISTS analytics_schema, reporting_schema DROP EXTERNAL DATABASE CASCADE;

-- DROP SCHEMA with underscores in name
DROP SCHEMA data_warehouse_schema RESTRICT;

-- DROP SCHEMA with numbers in name
DROP SCHEMA schema2023 CASCADE;

-- DROP SCHEMA with mixed case (should be case-insensitive)
drop schema MixedCaseSchema cascade;

-- DROP SCHEMA with multiple schemas and RESTRICT
DROP SCHEMA temp_schema1, temp_schema2, temp_schema3 RESTRICT;

-- DROP SCHEMA with IF EXISTS and multiple schemas with external database
DROP SCHEMA IF EXISTS external_schema1, external_schema2 DROP EXTERNAL DATABASE;