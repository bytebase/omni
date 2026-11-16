CREATE SCHEMA IF NOT EXISTS "default";

CREATE SCHEMA IF NOT EXISTS "information_schema";

CREATE SCHEMA IF NOT EXISTS "omni_driver";

CREATE TABLE IF NOT EXISTS "information_schema"."applicable_roles" (
    "grantee" varchar,
    "grantee_type" varchar,
    "role_name" varchar,
    "is_grantable" varchar
);

CREATE TABLE IF NOT EXISTS "information_schema"."columns" (
    "table_catalog" varchar,
    "table_schema" varchar,
    "table_name" varchar,
    "column_name" varchar,
    "ordinal_position" bigint,
    "column_default" varchar,
    "is_nullable" varchar,
    "data_type" varchar
);

CREATE TABLE IF NOT EXISTS "information_schema"."enabled_roles" (
    "role_name" varchar
);

CREATE TABLE IF NOT EXISTS "information_schema"."roles" (
    "role_name" varchar
);

CREATE TABLE IF NOT EXISTS "information_schema"."schemata" (
    "catalog_name" varchar,
    "schema_name" varchar
);

CREATE TABLE IF NOT EXISTS "information_schema"."table_privileges" (
    "grantor" varchar,
    "grantor_type" varchar,
    "grantee" varchar,
    "grantee_type" varchar,
    "table_catalog" varchar,
    "table_schema" varchar,
    "table_name" varchar,
    "privilege_type" varchar,
    "is_grantable" varchar,
    "with_hierarchy" varchar
);

CREATE TABLE IF NOT EXISTS "information_schema"."tables" (
    "table_catalog" varchar,
    "table_schema" varchar,
    "table_name" varchar,
    "table_type" varchar
);

CREATE TABLE IF NOT EXISTS "information_schema"."views" (
    "table_catalog" varchar,
    "table_schema" varchar,
    "table_name" varchar,
    "view_definition" varchar
);

CREATE TABLE IF NOT EXISTS "omni_driver"."a""b" (
    "c""d" varchar,
    "id" integer
);

CREATE TABLE IF NOT EXISTS "omni_driver"."mixedcase" (
    "id" integer NOT NULL,
    "weird name" varchar(10),
    "normal" varchar
);

CREATE TABLE IF NOT EXISTS "omni_driver"."not_null" (
    "id" bigint NOT NULL,
    "note" varchar(64) NOT NULL,
    "maybe" varchar
);

CREATE TABLE IF NOT EXISTS "omni_driver"."types" (
    "c_boolean" boolean,
    "c_tinyint" tinyint,
    "c_smallint" smallint,
    "c_integer" integer,
    "c_bigint" bigint,
    "c_real" real,
    "c_double" double,
    "c_decimal" decimal(10,2),
    "c_varchar" varchar,
    "c_varchar_length" varchar(20),
    "c_char" char(5),
    "c_varbinary" varbinary,
    "c_date" date,
    "c_time" time(3),
    "c_time_tz" time(3) with time zone,
    "c_timestamp" timestamp(3),
    "c_timestamp_tz" timestamp(3) with time zone,
    "c_json" json,
    "c_uuid" uuid,
    "c_ipaddress" ipaddress,
    "c_interval" interval day to second,
    "c_array" array(integer),
    "c_map" map(varchar, integer),
    "c_row" row("x" integer, "y" varchar)
);

CREATE OR REPLACE VIEW "omni_driver"."v_types" AS
SELECT
  c_integer
, c_varchar
FROM
  memory.omni_driver.types;

CREATE OR REPLACE VIEW "omni_driver"."v_derived" AS
SELECT c_integer
FROM
  memory.omni_driver.v_types;

