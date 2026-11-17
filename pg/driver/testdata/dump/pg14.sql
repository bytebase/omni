
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

CREATE SCHEMA IF NOT EXISTS "Mixed Case";

ALTER SCHEMA "Mixed Case" OWNER TO "postgres";

CREATE SCHEMA IF NOT EXISTS "app";

ALTER SCHEMA "app" OWNER TO "corpus_owner";

COMMENT ON SCHEMA "app" IS 'the application';

CREATE SCHEMA IF NOT EXISTS "audit";

ALTER SCHEMA "audit" OWNER TO "postgres";

COMMENT ON SCHEMA "public" IS 'standard public schema';

CREATE EXTENSION IF NOT EXISTS "btree_gist" WITH SCHEMA "public" VERSION '1.6';

COMMENT ON EXTENSION "btree_gist" IS 'support for indexing common datatypes in GiST';

CREATE EXTENSION IF NOT EXISTS "citext" WITH SCHEMA "public" VERSION '1.6';

COMMENT ON EXTENSION "citext" IS 'data type for case-insensitive character strings';

CREATE EXTENSION IF NOT EXISTS "hstore" WITH SCHEMA "public" VERSION '1.8';

COMMENT ON EXTENSION "hstore" IS 'key/value pairs';

CREATE EXTENSION IF NOT EXISTS "postgres_fdw" WITH SCHEMA "public" VERSION '1.1';

COMMENT ON EXTENSION "postgres_fdw" IS 'foreign-data wrapper for remote PostgreSQL servers';

CREATE TYPE "app"."mood" AS ENUM (
    'happy',
    'sad',
    'it''s complicated'
);

COMMENT ON TYPE "app"."mood" IS 'how it feels';

CREATE TYPE "app"."address" AS (
    "street" text COLLATE "C",
    "city" character varying(50),
    "zip" integer
);

COMMENT ON TYPE "app"."address" IS 'postal address';

COMMENT ON COLUMN "app"."address"."street" IS 'street line';

CREATE TYPE "app"."addresses" AS (
    "home" app.address,
    "others" app.address[]
);

CREATE TYPE "app"."empty_type" AS (

);

CREATE SEQUENCE "app"."counter"
    AS integer
	START WITH 10
	INCREMENT BY 5
	MINVALUE -100
	MAXVALUE 1000
	CACHE 3
	CYCLE;

COMMENT ON SEQUENCE "app"."counter" IS 'counts by five';

CREATE SEQUENCE "app"."orders_id_seq"
    AS bigint
	START WITH 1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	CACHE 1
	NO CYCLE;

CREATE SEQUENCE "app"."owned_seq"
    AS bigint
	START WITH 1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	CACHE 1
	NO CYCLE;

CREATE SEQUENCE "app"."seq_owner_id_seq"
    AS integer
	START WITH 1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	CACHE 1
	NO CYCLE;

CREATE SEQUENCE "audit"."big_seq"
    AS bigint
	START WITH 1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 9223372036854775807
	CACHE 1
	NO CYCLE;

CREATE SEQUENCE "public"."plain_a_seq"
    AS integer
	START WITH 1
	INCREMENT BY 1
	MINVALUE 1
	MAXVALUE 2147483647
	CACHE 1
	NO CYCLE;

SET default_tablespace = '';

CREATE TABLE "Mixed Case"."Table" (
    "Id" integer NOT NULL,
    "Value" text
);

ALTER TABLE "Mixed Case"."Table" OWNER TO "postgres";

ALTER TABLE ONLY "Mixed Case"."Table" ADD CONSTRAINT "Table_pkey" PRIMARY KEY ("Id");

CREATE OR REPLACE FUNCTION app.add(a integer, b integer DEFAULT 1)
 RETURNS integer
 LANGUAGE sql
 IMMUTABLE STRICT
AS $function$ SELECT a + b $function$
;

COMMENT ON FUNCTION "app".add(integer, integer) IS 'adds integers';

CREATE OR REPLACE FUNCTION app.add(a numeric, b numeric)
 RETURNS numeric
 LANGUAGE sql
 IMMUTABLE
AS $function$ SELECT a + b $function$
;

CREATE TABLE "app"."animal" (
    "name" text,
    "legs" integer,
    CONSTRAINT "animal_legs_check" CHECK (legs >= 0)
);

ALTER TABLE "app"."animal" OWNER TO "postgres";

CREATE TABLE "app"."booking" (
    "room" integer,
    "during" tstzrange,
    CONSTRAINT "booking_no_overlap" EXCLUDE USING gist (room WITH =, during WITH &&)
);

ALTER TABLE "app"."booking" OWNER TO "postgres";

CREATE TABLE "app"."customer" (
    "id" bigint NOT NULL,
    "alt_id" integer NOT NULL,
    "email" public.citext NOT NULL,
    "name" text COLLATE "C" DEFAULT 'anon'::text NOT NULL,
    "nick" text COLLATE "ci",
    "mood" app.mood DEFAULT 'happy'::app.mood,
    "home" app.address,
    "tags" _text DEFAULT '{}'::text[],
    "matrix" _int4,
    "attrs" public.hstore,
    "doc" jsonb DEFAULT '{"a": [1, 2]}'::jsonb,
    "balance" numeric(12,2) DEFAULT 0,
    "ratio" double precision,
    "score" real,
    "small" smallint,
    "flag" boolean DEFAULT true,
    "ts" timestamp(3) without time zone DEFAULT now(),
    "tstz" timestamp(6) with time zone,
    "born" date,
    "alarm" time(2) with time zone,
    "span" interval,
    "bits" bit(3),
    "vbits" bit varying(8),
    "ip" inet,
    "mac" macaddr,
    "pt" point,
    "words" tsvector,
    "uid" uuid,
    "raw" bytea DEFAULT '\xdeadbeef'::bytea,
    "counter" integer DEFAULT nextval('app.counter'::regclass),
    "owned" bigint DEFAULT nextval('app.owned_seq'::regclass),
    "full_name" text GENERATED ALWAYS AS ((((name || ' <'::text) || (email)::text) || '>'::text)) STORED,
    "Quoted Col" text,
    "updated_at" timestamp(6) with time zone,
    "moods" app._mood,
    CONSTRAINT "customer_balance_check" CHECK (balance >= 0::numeric),
    CONSTRAINT "customer_name_check" CHECK (length(name) > 0)
);

ALTER SEQUENCE "app"."owned_seq" OWNED BY "app"."customer"."owned";

ALTER TABLE "app"."customer" OWNER TO "corpus_owner";

COMMENT ON TABLE "app"."customer" IS 'who buys';

COMMENT ON COLUMN "app"."customer"."email" IS 'unique email';

COMMENT ON CONSTRAINT "customer_name_check" ON "app"."customer" IS 'the name a customer is known by';

ALTER TABLE "app"."customer" ALTER COLUMN "alt_id" ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME "app"."customer_alt_id_seq"
    START WITH 1
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1
);

ALTER TABLE "app"."customer" ALTER COLUMN "id" ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME "app"."customer_id_seq"
    START WITH 100
    INCREMENT BY 2
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1
);

ALTER TABLE ONLY "app"."customer" ADD CONSTRAINT "customer_pkey" PRIMARY KEY (id);

ALTER TABLE ONLY "app"."customer" ADD CONSTRAINT "customer_email_key" UNIQUE (email);

CREATE INDEX "customer_brin" ON "app"."customer" USING BRIN (ts) WITH (pages_per_range='16');

CREATE INDEX "customer_include" ON "app"."customer" (alt_id) INCLUDE (name, mood);

CREATE INDEX "customer_lower_name" ON "app"."customer" (lower(name) DESC NULLS LAST);

CREATE INDEX "customer_opclass" ON "app"."customer" (name text_pattern_ops, nick);

CREATE INDEX "customer_partial" ON "app"."customer" (mood, id) WHERE flag;

COMMENT ON INDEX "app"."customer_partial" IS 'happy flags';

CREATE INDEX "customer_tags" ON "app"."customer" USING GIN (tags);

CREATE OR REPLACE FUNCTION app.customer_count()
 RETURNS bigint
 LANGUAGE sql
 STABLE
AS $function$SELECT count(*) FROM app.customer$function$
;

CREATE TABLE "app"."dog" (
    "name" text,
    "legs" integer,
    "breed" text,
    CONSTRAINT "animal_legs_check" CHECK (legs >= 0)
);

ALTER TABLE "app"."dog" OWNER TO "postgres";

multirange_constructor0;

multirange_constructor2;

multirange_constructor1;

range_constructor2;

range_constructor3;

CREATE OR REPLACE FUNCTION app.items_named(p text)
 RETURNS SETOF text
 LANGUAGE sql
 STABLE SECURITY DEFINER
 SET search_path TO 'app'
AS $function$ SELECT p $function$
;

CREATE TABLE "app"."measurement" (
    "city_id" integer NOT NULL,
    "logdate" date NOT NULL,
    "region" text NOT NULL,
    "peak" integer,
    CONSTRAINT "measurement_peak_check" CHECK (peak > '-100'::integer)
) PARTITION BY RANGE (logdate);

ALTER TABLE "app"."measurement" OWNER TO "postgres";

CREATE TABLE "app"."measurement_2025" (
    "city_id" integer NOT NULL,
    "logdate" date NOT NULL,
    "region" text NOT NULL,
    "peak" integer,
    CONSTRAINT "measurement_peak_check" CHECK (peak > '-100'::integer)
) PARTITION BY LIST (region);

CREATE TABLE "app"."measurement_2025_eu" (
    "city_id" integer NOT NULL,
    "logdate" date NOT NULL,
    "region" text NOT NULL,
    "peak" integer,
    CONSTRAINT "measurement_peak_check" CHECK (peak > '-100'::integer)
);

CREATE TABLE "app"."measurement_2025_other" (
    "city_id" integer NOT NULL,
    "logdate" date NOT NULL,
    "region" text NOT NULL,
    "peak" integer,
    CONSTRAINT "measurement_peak_check" CHECK (peak > '-100'::integer)
);

ALTER TABLE ONLY "app"."measurement_2025" ATTACH PARTITION "app"."measurement_2025_eu" FOR VALUES IN ('eu');

ALTER TABLE ONLY "app"."measurement_2025" ATTACH PARTITION "app"."measurement_2025_other" DEFAULT;

CREATE TABLE "app"."measurement_2026" (
    "city_id" integer NOT NULL,
    "logdate" date NOT NULL,
    "region" text NOT NULL,
    "peak" integer,
    CONSTRAINT "measurement_peak_check" CHECK (peak > '-100'::integer)
);

ALTER TABLE ONLY "app"."measurement" ATTACH PARTITION "app"."measurement_2025" FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

ALTER TABLE ONLY "app"."measurement" ATTACH PARTITION "app"."measurement_2026" FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE INDEX "measurement_city" ON ONLY "app"."measurement" (city_id, logdate);

CREATE UNIQUE INDEX "measurement_key" ON ONLY "app"."measurement" (city_id, logdate, region);

CREATE INDEX "measurement_2025_city_id_logdate_idx" ON ONLY "app"."measurement_2025" (city_id, logdate);

CREATE UNIQUE INDEX "measurement_2025_city_id_logdate_region_idx" ON ONLY "app"."measurement_2025" (city_id, logdate, region);

CREATE INDEX "measurement_2025_eu_city_id_logdate_idx" ON "app"."measurement_2025_eu" (city_id, logdate);

CREATE UNIQUE INDEX "measurement_2025_eu_city_id_logdate_region_idx" ON "app"."measurement_2025_eu" (city_id, logdate, region);

CREATE INDEX "measurement_2025_other_city_id_logdate_idx" ON "app"."measurement_2025_other" (city_id, logdate);

CREATE UNIQUE INDEX "measurement_2025_other_city_id_logdate_region_idx" ON "app"."measurement_2025_other" (city_id, logdate, region);

CREATE INDEX "measurement_2026_city_id_logdate_idx" ON "app"."measurement_2026" (city_id, logdate);

CREATE UNIQUE INDEX "measurement_2026_city_id_logdate_region_idx" ON "app"."measurement_2026" (city_id, logdate, region);

ALTER INDEX "app"."measurement_city" ATTACH PARTITION "app"."measurement_2025_city_id_logdate_idx";

ALTER INDEX "app"."measurement_key" ATTACH PARTITION "app"."measurement_2025_city_id_logdate_region_idx";

ALTER INDEX "app"."measurement_2025_city_id_logdate_idx" ATTACH PARTITION "app"."measurement_2025_eu_city_id_logdate_idx";

ALTER INDEX "app"."measurement_2025_city_id_logdate_region_idx" ATTACH PARTITION "app"."measurement_2025_eu_city_id_logdate_region_idx";

ALTER INDEX "app"."measurement_2025_city_id_logdate_idx" ATTACH PARTITION "app"."measurement_2025_other_city_id_logdate_idx";

ALTER INDEX "app"."measurement_2025_city_id_logdate_region_idx" ATTACH PARTITION "app"."measurement_2025_other_city_id_logdate_region_idx";

ALTER INDEX "app"."measurement_city" ATTACH PARTITION "app"."measurement_2026_city_id_logdate_idx";

ALTER INDEX "app"."measurement_key" ATTACH PARTITION "app"."measurement_2026_city_id_logdate_region_idx";

CREATE OR REPLACE FUNCTION app.names_insert()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN NEW;
END
$function$
;

CREATE TABLE "app"."orders" (
    "id" bigint DEFAULT nextval('app.orders_id_seq'::regclass) NOT NULL,
    "customer_id" bigint NOT NULL,
    "placed_on" date NOT NULL,
    "total" numeric
);

ALTER SEQUENCE "app"."orders_id_seq" OWNED BY "app"."orders"."id";

ALTER TABLE "app"."orders" OWNER TO "postgres";

ALTER TABLE ONLY "app"."orders" ADD CONSTRAINT "orders_pkey" PRIMARY KEY (id, placed_on);

CREATE TABLE "app"."ref_source" (
    "x" integer,
    "y" integer
);

ALTER TABLE "app"."ref_source" OWNER TO "postgres";

CREATE TABLE "app"."ref_target" (
    "a" integer,
    "b" integer
);

ALTER TABLE "app"."ref_target" OWNER TO "postgres";

ALTER TABLE ONLY "app"."ref_target" ADD CONSTRAINT "ref_target_a_b_key" UNIQUE (a, b);

CREATE OR REPLACE PROCEDURE app.reset(IN n integer)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    PERFORM setval('app.counter', n);
END
$procedure$
;

CREATE TABLE "app"."seq_owner" (
    "id" integer DEFAULT nextval('app.seq_owner_id_seq'::regclass) NOT NULL,
    "payload" text
);

ALTER SEQUENCE "app"."seq_owner_id_seq" OWNED BY "app"."seq_owner"."id";

ALTER TABLE "app"."seq_owner" OWNER TO "postgres";

CREATE TABLE "app"."sharded" (
    "id" integer,
    "v" text
) PARTITION BY HASH (id);

ALTER TABLE "app"."sharded" OWNER TO "postgres";

CREATE TABLE "app"."sharded_0" (
    "id" integer,
    "v" text
);

CREATE TABLE "app"."sharded_1" (
    "id" integer,
    "v" text
);

ALTER TABLE ONLY "app"."sharded" ATTACH PARTITION "app"."sharded_0" FOR VALUES WITH (modulus 2, remainder 0);

ALTER TABLE ONLY "app"."sharded" ATTACH PARTITION "app"."sharded_1" FOR VALUES WITH (modulus 2, remainder 1);

CREATE OR REPLACE FUNCTION app.touch()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END
$function$
;

CREATE OR REPLACE FUNCTION audit.on_ddl()
 RETURNS event_trigger
 LANGUAGE plpgsql
AS $function$
BEGIN
    RAISE NOTICE 'ddl: %', tg_tag;
END
$function$
;

CREATE TABLE "audit"."scratch" (
    "note" text
);

ALTER TABLE "audit"."scratch" OWNER TO "postgres";

CREATE TABLE "public"."plain" (
    "a" integer DEFAULT nextval('public.plain_a_seq'::regclass) NOT NULL,
    "b" text
);

ALTER SEQUENCE "public"."plain_a_seq" OWNED BY "public"."plain"."a";

ALTER TABLE "public"."plain" OWNER TO "postgres";

CREATE TABLE "public"."z_source" (
    "v" integer
);

ALTER TABLE "public"."z_source" OWNER TO "postgres";

CREATE VIEW "app"."customer_names" AS 
 SELECT customer.id,
    customer.name
   FROM app.customer;

CREATE VIEW "app"."owner's view" AS 
 SELECT customer.id
   FROM app.customer;

CREATE VIEW "app"."happy_customer" AS 
 SELECT c.id,
    c.name,
    o.total
   FROM (app.customer c
     JOIN app.orders o ON ((o.customer_id = c.id)))
  WHERE (c.mood = 'happy'::app.mood);

COMMENT ON VIEW "app"."happy_customer" IS 'happy ones';

CREATE MATERIALIZED VIEW "app"."order_totals" AS 
 SELECT orders.customer_id,
    sum(orders.total) AS total
   FROM app.orders
  GROUP BY orders.customer_id
  WITH NO DATA;

CREATE UNIQUE INDEX "order_totals_customer" ON "app"."order_totals" (customer_id);

COMMENT ON MATERIALIZED VIEW "app"."order_totals" IS 'per customer';

CREATE VIEW "public"."a_count" AS 
 SELECT count(*) AS n
   FROM public.z_source;

CREATE VIEW "app"."view_on_view" AS 
 SELECT happy_customer.name
   FROM app.happy_customer;

CREATE TRIGGER customer_touch BEFORE UPDATE ON app.customer FOR EACH ROW WHEN ((old.name IS DISTINCT FROM new.name)) EXECUTE FUNCTION app.touch();

COMMENT ON TRIGGER "customer_touch" ON "app"."customer" IS 'keeps updated_at';

CREATE TRIGGER names_insert INSTEAD OF INSERT ON app.customer_names FOR EACH ROW EXECUTE FUNCTION app.names_insert();

CREATE RULE plain_log AS
    ON INSERT TO public.plain
   WHERE new.b IS NOT NULL DO  INSERT INTO audit.scratch (note)
  VALUES (new.b);

CREATE RULE plain_nothing AS
    ON DELETE TO public.plain DO INSTEAD NOTHING;

ALTER TABLE "app"."orders"
    ADD CONSTRAINT "orders_customer_fk" FOREIGN KEY ("customer_id")
    REFERENCES "app"."customer" ("id")
    ON DELETE CASCADE
    ON UPDATE SET NULL;

COMMENT ON CONSTRAINT "orders_customer_fk" ON "app"."orders" IS 'an order belongs to a customer';

ALTER TABLE "app"."ref_source"
    ADD CONSTRAINT "ref_source_x_y_fkey" FOREIGN KEY ("x", "y")
    REFERENCES "app"."ref_target" ("a", "b")
    MATCH FULL
    ON DELETE SET DEFAULT;

CREATE EVENT TRIGGER "corpus_ddl" ON ddl_command_end
  WHEN TAG IN ('CREATE TABLE', 'ALTER TABLE')
  EXECUTE FUNCTION "audit"."on_ddl"();

COMMENT ON EVENT TRIGGER "corpus_ddl" IS 'logs DDL';

CREATE EVENT TRIGGER "corpus_drop" ON sql_drop
  EXECUTE FUNCTION "audit"."on_ddl"();

ALTER EVENT TRIGGER "corpus_drop" DISABLE;

