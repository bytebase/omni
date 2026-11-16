-- The schema trino/driver syncs, dumps and restores. Every memory-catalog table
-- lives in one schema so a test can drop and recreate just that schema; the
-- image's information_schema and default schemas are in the synced metadata
-- too, and a dump leaves them alone (CREATE TABLE IF NOT EXISTS).
CREATE SCHEMA memory.omni_driver;

CREATE TABLE memory.omni_driver.types (
    c_boolean boolean,
    c_tinyint tinyint,
    c_smallint smallint,
    c_integer integer,
    c_bigint bigint,
    c_real real,
    c_double double,
    c_decimal decimal(10, 2),
    c_varchar varchar,
    c_varchar_length varchar(20),
    c_char char(5),
    c_varbinary varbinary,
    c_date date,
    c_time time(3),
    c_time_tz time(3) with time zone,
    c_timestamp timestamp(3),
    c_timestamp_tz timestamp(3) with time zone,
    c_json json,
    c_uuid uuid,
    c_ipaddress ipaddress,
    c_interval interval day to second,
    c_array array(integer),
    c_map map(varchar, integer),
    c_row row(x integer, y varchar)
);

CREATE TABLE memory.omni_driver.not_null (
    id bigint NOT NULL,
    note varchar(64) NOT NULL,
    maybe varchar
);

CREATE TABLE memory.omni_driver."MixedCase" (
    "Id" integer NOT NULL,
    "Weird Name" varchar(10),
    normal varchar
);

CREATE TABLE memory.omni_driver."a""b" (
    "c""d" varchar,
    id integer
);

CREATE VIEW memory.omni_driver.v_types AS
    SELECT c_integer, c_varchar FROM memory.omni_driver.types;

-- Sorts before v_types, so a dump must order it after its dependency.
CREATE VIEW memory.omni_driver.v_derived AS
    SELECT c_integer FROM memory.omni_driver.v_types;
