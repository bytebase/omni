-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestDatabaseDefinitionCorpusRunsOnPostgres/basic_tables_with_constraints default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."posts_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| CREATE SEQUENCE "public"."users_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."posts" (
|     "id" integer DEFAULT nextval('public.posts_id_seq'::regclass) NOT NULL,
|     "user_id" integer NOT NULL,
|     "title" character varying(200) NOT NULL,
|     "content" text,
|     "published" boolean DEFAULT false
| );
|
| ALTER SEQUENCE "public"."posts_id_seq" OWNED BY "public"."posts"."id";
|
| ALTER TABLE "public"."posts" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."posts" ADD CONSTRAINT "posts_pkey" PRIMARY KEY (id);
|
| CREATE INDEX "idx_posts_published" ON "public"."posts" (published) WHERE published = true;
|
| CREATE INDEX "idx_posts_user_id" ON "public"."posts" (user_id);
|
| CREATE TABLE "public"."users" (
|     "id" integer DEFAULT nextval('public.users_id_seq'::regclass) NOT NULL,
|     "username" character varying(50) NOT NULL,
|     "email" character varying(100) NOT NULL,
|     "created_at" timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP,
|     "updated_at" timestamp(6) without time zone
| );
|
| ALTER SEQUENCE "public"."users_id_seq" OWNED BY "public"."users"."id";
|
| ALTER TABLE "public"."users" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY (id);
|
| ALTER TABLE ONLY "public"."users" ADD CONSTRAINT "users_username_key" UNIQUE (username);
|
| ALTER TABLE "public"."posts"
|     ADD CONSTRAINT "fk_user" FOREIGN KEY ("user_id")
|     REFERENCES "public"."users" ("id")
|     ON DELETE CASCADE;
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/check_constraints_and_comments default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."products_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."products" (
|     "id" integer DEFAULT nextval('public.products_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "price" numeric(10,2) NOT NULL,
|     "quantity" integer NOT NULL,
|     CONSTRAINT "chk_price" CHECK (price > 0::numeric),
|     CONSTRAINT "chk_quantity" CHECK (quantity >= 0)
| );
|
| ALTER SEQUENCE "public"."products_id_seq" OWNED BY "public"."products"."id";
|
| ALTER TABLE "public"."products" OWNER TO "postgres";
|
| COMMENT ON TABLE "public"."products" IS 'Product catalog table';
|
| COMMENT ON COLUMN "public"."products"."name" IS 'Product name';
|
| COMMENT ON COLUMN "public"."products"."price" IS 'Product price in USD';
|
| ALTER TABLE ONLY "public"."products" ADD CONSTRAINT "products_pkey" PRIMARY KEY (id);
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/composite_types default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| CREATE SCHEMA IF NOT EXISTS "geo";
|
| ALTER SCHEMA "geo" OWNER TO "postgres";
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE TYPE "public"."order_status" AS ENUM (
|     'pending',
|     'shipped'
| );
|
| CREATE TYPE "geo"."point2" AS (
|     "lat" numeric(9,6),
|     "lng" numeric(9,6)
| );
|
| CREATE TYPE "public"."empty_type" AS (
|
| );
|
| CREATE TYPE "public"."order_info" AS (
|     "status" public.order_status,
|     "note" text
| );
|
| CREATE TYPE "public"."wrapper" AS (
|     "p" geo.point2
| );
|
| CREATE TYPE "public"."zz_base" AS (
|     "street" text COLLATE "C",
|     "city" character varying(50)
| );
|
| COMMENT ON TYPE "public"."zz_base" IS 'base address type';
|
| COMMENT ON COLUMN "public"."zz_base"."street" IS 'street line';
|
| CREATE TYPE "public"."aa_nested" AS (
|     "home" public.zz_base,
|     "homes" public.zz_base[],
|     "tags" text[]
| );
|
| CREATE SEQUENCE "public"."orders_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."orders" (
|     "id" integer DEFAULT nextval('public.orders_id_seq'::regclass) NOT NULL,
|     "info" public.order_info,
|     "addresses" public.aa_nested,
|     "location" geo.point2
| );
|
| ALTER SEQUENCE "public"."orders_id_seq" OWNED BY "public"."orders"."id";
|
| ALTER TABLE "public"."orders" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."orders" ADD CONSTRAINT "orders_pkey" PRIMARY KEY (id);
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/comprehensive_views default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."customers_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| CREATE SEQUENCE "public"."orders_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."customers" (
|     "id" integer DEFAULT nextval('public.customers_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "email" character varying(100),
|     "active" boolean DEFAULT true,
|     "created_at" timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP
| );
|
| ALTER SEQUENCE "public"."customers_id_seq" OWNED BY "public"."customers"."id";
|
| ALTER TABLE "public"."customers" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."customers" ADD CONSTRAINT "customers_pkey" PRIMARY KEY (id);
|
| ALTER TABLE ONLY "public"."customers" ADD CONSTRAINT "customers_email_key" UNIQUE (email);
|
| CREATE TABLE "public"."orders" (
|     "id" integer DEFAULT nextval('public.orders_id_seq'::regclass) NOT NULL,
|     "customer_id" integer,
|     "total_amount" numeric(10,2) NOT NULL,
|     "order_date" timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP,
|     "status" character varying(20) DEFAULT 'pending'::character varying
| );
|
| ALTER SEQUENCE "public"."orders_id_seq" OWNED BY "public"."orders"."id";
|
| ALTER TABLE "public"."orders" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."orders" ADD CONSTRAINT "orders_pkey" PRIMARY KEY (id);
|
| CREATE VIEW "public"."active_customer_orders" AS 
|  SELECT c.id AS customer_id,
|     c.name,
|     c.email,
|     o.id AS order_id,
|     o.total_amount,
|     o.order_date
|    FROM (public.customers c
|      JOIN public.orders o ON ((c.id = o.customer_id)))
|   WHERE ((c.active = true) AND ((o.status)::text = 'completed'::text));
|
| CREATE VIEW "public"."customer_summary" AS 
|  SELECT id,
|     name,
|     email,
|     ( SELECT count(*) AS count
|            FROM public.orders
|           WHERE (orders.customer_id = customers.id)) AS order_count
|    FROM public.customers
|   WHERE (active = true);
|
| CREATE MATERIALIZED VIEW "public"."monthly_sales_summary" AS 
|  SELECT date_trunc('month'::text, order_date) AS month,
|     count(*) AS order_count,
|     sum(total_amount) AS total_revenue,
|     avg(total_amount) AS avg_order_value,
|     count(DISTINCT customer_id) AS unique_customers
|    FROM public.orders
|   WHERE ((status)::text = 'completed'::text)
|   GROUP BY (date_trunc('month'::text, order_date))
|   ORDER BY (date_trunc('month'::text, order_date))
|   WITH NO DATA;
|
| CREATE UNIQUE INDEX "idx_monthly_sales_month" ON "public"."monthly_sales_summary" (month);
|
| ALTER TABLE "public"."orders"
|     ADD CONSTRAINT "orders_customer_id_fkey" FOREIGN KEY ("customer_id")
|     REFERENCES "public"."customers" ("id");
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/extensions_and_advanced_features default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA "public" VERSION '1.1';
|
| COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."documents" (
|     "id" uuid DEFAULT public.uuid_generate_v4() NOT NULL,
|     "title" character varying(200) NOT NULL,
|     "content" text,
|     "tags" _text,
|     "metadata" jsonb,
|     "created_at" timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP
| );
|
| ALTER TABLE "public"."documents" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."documents" ADD CONSTRAINT "documents_pkey" PRIMARY KEY (id);
|
| CREATE INDEX "idx_documents_metadata" ON "public"."documents" USING GIN (metadata);
|
| CREATE INDEX "idx_documents_tags" ON "public"."documents" USING GIN (tags);
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/foreign_key_no_action default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."cities_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| CREATE SEQUENCE "public"."countries_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| CREATE SEQUENCE "public"."states_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."cities" (
|     "id" integer DEFAULT nextval('public.cities_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "state_id" integer
| );
|
| ALTER SEQUENCE "public"."cities_id_seq" OWNED BY "public"."cities"."id";
|
| ALTER TABLE "public"."cities" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."cities" ADD CONSTRAINT "cities_pkey" PRIMARY KEY (id);
|
| CREATE TABLE "public"."countries" (
|     "id" integer DEFAULT nextval('public.countries_id_seq'::regclass) NOT NULL,
|     "code" character varying(3) NOT NULL,
|     "name" character varying(100) NOT NULL
| );
|
| ALTER SEQUENCE "public"."countries_id_seq" OWNED BY "public"."countries"."id";
|
| ALTER TABLE "public"."countries" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."countries" ADD CONSTRAINT "countries_pkey" PRIMARY KEY (id);
|
| ALTER TABLE ONLY "public"."countries" ADD CONSTRAINT "countries_code_key" UNIQUE (code);
|
| CREATE TABLE "public"."states" (
|     "id" integer DEFAULT nextval('public.states_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "country_id" integer
| );
|
| ALTER SEQUENCE "public"."states_id_seq" OWNED BY "public"."states"."id";
|
| ALTER TABLE "public"."states" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."states" ADD CONSTRAINT "states_pkey" PRIMARY KEY (id);
|
| ALTER TABLE "public"."cities"
|     ADD CONSTRAINT "fk_state_implicit" FOREIGN KEY ("state_id")
|     REFERENCES "public"."states" ("id");
|
| ALTER TABLE "public"."states"
|     ADD CONSTRAINT "fk_country_explicit" FOREIGN KEY ("country_id")
|     REFERENCES "public"."countries" ("id");
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/foreign_key_on_delete_cascade default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."categories_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| CREATE SEQUENCE "public"."orders_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| CREATE SEQUENCE "public"."products_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."categories" (
|     "id" integer DEFAULT nextval('public.categories_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL
| );
|
| ALTER SEQUENCE "public"."categories_id_seq" OWNED BY "public"."categories"."id";
|
| ALTER TABLE "public"."categories" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."categories" ADD CONSTRAINT "categories_pkey" PRIMARY KEY (id);
|
| CREATE TABLE "public"."orders" (
|     "id" integer DEFAULT nextval('public.orders_id_seq'::regclass) NOT NULL,
|     "product_id" integer,
|     "quantity" integer DEFAULT 1
| );
|
| ALTER SEQUENCE "public"."orders_id_seq" OWNED BY "public"."orders"."id";
|
| ALTER TABLE "public"."orders" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."orders" ADD CONSTRAINT "orders_pkey" PRIMARY KEY (id);
|
| CREATE TABLE "public"."products" (
|     "id" integer DEFAULT nextval('public.products_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "category_id" integer
| );
|
| ALTER SEQUENCE "public"."products_id_seq" OWNED BY "public"."products"."id";
|
| ALTER TABLE "public"."products" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."products" ADD CONSTRAINT "products_pkey" PRIMARY KEY (id);
|
| ALTER TABLE "public"."orders"
|     ADD CONSTRAINT "fk_product_cascade" FOREIGN KEY ("product_id")
|     REFERENCES "public"."products" ("id")
|     ON DELETE SET NULL;
|
| ALTER TABLE "public"."products"
|     ADD CONSTRAINT "fk_category" FOREIGN KEY ("category_id")
|     REFERENCES "public"."categories" ("id")
|     ON DELETE CASCADE;
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/foreign_key_on_update_actions default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."employees_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| CREATE SEQUENCE "public"."projects_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."departments" (
|     "code" character varying(10) NOT NULL,
|     "name" character varying(100) NOT NULL
| );
|
| ALTER TABLE "public"."departments" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."departments" ADD CONSTRAINT "departments_pkey" PRIMARY KEY (code);
|
| CREATE TABLE "public"."employees" (
|     "id" integer DEFAULT nextval('public.employees_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "dept_code" character varying(10),
|     "manager_id" integer
| );
|
| ALTER SEQUENCE "public"."employees_id_seq" OWNED BY "public"."employees"."id";
|
| ALTER TABLE "public"."employees" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."employees" ADD CONSTRAINT "employees_pkey" PRIMARY KEY (id);
|
| CREATE TABLE "public"."projects" (
|     "id" integer DEFAULT nextval('public.projects_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "lead_emp_id" integer,
|     "dept_code" character varying(10)
| );
|
| ALTER SEQUENCE "public"."projects_id_seq" OWNED BY "public"."projects"."id";
|
| ALTER TABLE "public"."projects" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."projects" ADD CONSTRAINT "projects_pkey" PRIMARY KEY (id);
|
| ALTER TABLE "public"."employees"
|     ADD CONSTRAINT "fk_dept_restrict" FOREIGN KEY ("dept_code")
|     REFERENCES "public"."departments" ("code")
|     ON DELETE SET NULL
|     ON UPDATE RESTRICT;
|
| ALTER TABLE "public"."employees"
|     ADD CONSTRAINT "fk_manager_cascade" FOREIGN KEY ("manager_id")
|     REFERENCES "public"."employees" ("id")
|     ON DELETE CASCADE
|     ON UPDATE CASCADE;
|
| ALTER TABLE "public"."projects"
|     ADD CONSTRAINT "fk_lead_setnull" FOREIGN KEY ("lead_emp_id")
|     REFERENCES "public"."employees" ("id")
|     ON DELETE SET NULL
|     ON UPDATE SET NULL;
|
| ALTER TABLE "public"."projects"
|     ADD CONSTRAINT "fk_proj_dept" FOREIGN KEY ("dept_code")
|     REFERENCES "public"."departments" ("code")
|     ON DELETE RESTRICT
|     ON UPDATE RESTRICT;
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/materialized_views default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."raw_data_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."raw_data" (
|     "id" integer DEFAULT nextval('public.raw_data_id_seq'::regclass) NOT NULL,
|     "data_value" numeric NOT NULL,
|     "category" character varying(50) NOT NULL,
|     "created_at" timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP
| );
|
| ALTER SEQUENCE "public"."raw_data_id_seq" OWNED BY "public"."raw_data"."id";
|
| ALTER TABLE "public"."raw_data" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."raw_data" ADD CONSTRAINT "raw_data_pkey" PRIMARY KEY (id);
|
| CREATE MATERIALIZED VIEW "public"."data_summary" AS 
|  SELECT category,
|     count(*) AS count,
|     avg(data_value) AS avg_value,
|     max(data_value) AS max_value,
|     min(data_value) AS min_value
|    FROM public.raw_data
|   GROUP BY category
|   WITH NO DATA;
|
| CREATE INDEX "idx_data_summary_category" ON "public"."data_summary" (category);
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/owned_sequence_dependency_order default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."history_cost_changes_id_seq"
|     AS bigint
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 9223372036854775807
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."history_audit" (
|     "audit_id" bigint DEFAULT nextval('public.history_cost_changes_id_seq'::regclass) NOT NULL,
|     "change_type" character varying(50)
| );
|
| ALTER TABLE "public"."history_audit" OWNER TO "postgres";
|
| CREATE TABLE "public"."history_cost_changes" (
|     "id" bigint DEFAULT nextval('public.history_cost_changes_id_seq'::regclass) NOT NULL,
|     "amount" numeric(10,2),
|     "description" text,
|     "created_at" timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP
| );
|
| ALTER SEQUENCE "public"."history_cost_changes_id_seq" OWNED BY "public"."history_cost_changes"."id";
|
| ALTER TABLE "public"."history_cost_changes" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."history_cost_changes" ADD CONSTRAINT "history_cost_changes_pkey" PRIMARY KEY (id);
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/partitioned_tables default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."sales_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."sales" (
|     "id" integer DEFAULT nextval('public.sales_id_seq'::regclass) NOT NULL,
|     "sale_date" date NOT NULL,
|     "amount" numeric(10,2) NOT NULL,
|     "region" character varying(50) NOT NULL
| ) PARTITION BY RANGE (sale_date);
|
| ALTER SEQUENCE "public"."sales_id_seq" OWNED BY "public"."sales"."id";
|
| ALTER TABLE "public"."sales" OWNER TO "postgres";
|
| CREATE TABLE "public"."sales_2023_q1" (
|     "id" integer DEFAULT nextval('public.sales_id_seq'::regclass) NOT NULL,
|     "sale_date" date NOT NULL,
|     "amount" numeric(10,2) NOT NULL,
|     "region" character varying(50) NOT NULL
| );
|
| CREATE TABLE "public"."sales_2023_q2" (
|     "id" integer DEFAULT nextval('public.sales_id_seq'::regclass) NOT NULL,
|     "sale_date" date NOT NULL,
|     "amount" numeric(10,2) NOT NULL,
|     "region" character varying(50) NOT NULL
| );
|
| CREATE TABLE "public"."sales_2023_q3" (
|     "id" integer DEFAULT nextval('public.sales_id_seq'::regclass) NOT NULL,
|     "sale_date" date NOT NULL,
|     "amount" numeric(10,2) NOT NULL,
|     "region" character varying(50) NOT NULL
| );
|
| CREATE TABLE "public"."sales_2023_q4" (
|     "id" integer DEFAULT nextval('public.sales_id_seq'::regclass) NOT NULL,
|     "sale_date" date NOT NULL,
|     "amount" numeric(10,2) NOT NULL,
|     "region" character varying(50) NOT NULL
| );
|
| ALTER TABLE ONLY "public"."sales" ATTACH PARTITION "public"."sales_2023_q1" FOR VALUES FROM ('2023-01-01') TO ('2023-04-01');
|
| ALTER TABLE ONLY "public"."sales" ATTACH PARTITION "public"."sales_2023_q2" FOR VALUES FROM ('2023-04-01') TO ('2023-07-01');
|
| ALTER TABLE ONLY "public"."sales" ATTACH PARTITION "public"."sales_2023_q3" FOR VALUES FROM ('2023-07-01') TO ('2023-10-01');
|
| ALTER TABLE ONLY "public"."sales" ATTACH PARTITION "public"."sales_2023_q4" FOR VALUES FROM ('2023-10-01') TO ('2024-01-01');
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/sequences_and_custom_types default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE TYPE "public"."mood" AS ENUM (
|     'happy',
|     'sad',
|     'neutral'
| );
|
| CREATE TYPE "public"."status_enum" AS ENUM (
|     'pending',
|     'active',
|     'inactive',
|     'deleted'
| );
|
| CREATE SEQUENCE "public"."custom_id_seq"
|     AS bigint
| 	START WITH 1000
| 	INCREMENT BY 10
| 	MINVALUE 1
| 	MAXVALUE 9223372036854775807
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."items" (
|     "id" integer DEFAULT nextval('public.custom_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "status" public.status_enum DEFAULT 'pending'::public.status_enum,
|     "user_mood" public.mood
| );
|
| ALTER TABLE "public"."items" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."items" ADD CONSTRAINT "items_pkey" PRIMARY KEY (id);
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/triggers default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."accounts_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| CREATE SEQUENCE "public"."audit_log_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE TABLE "public"."accounts" (
|     "id" integer DEFAULT nextval('public.accounts_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "balance" numeric(10,2) DEFAULT 0,
|     "last_modified" timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP
| );
|
| ALTER SEQUENCE "public"."accounts_id_seq" OWNED BY "public"."accounts"."id";
|
| ALTER TABLE "public"."accounts" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."accounts" ADD CONSTRAINT "accounts_pkey" PRIMARY KEY (id);
|
| CREATE TABLE "public"."audit_log" (
|     "id" integer DEFAULT nextval('public.audit_log_id_seq'::regclass) NOT NULL,
|     "table_name" character varying(50) NOT NULL,
|     "operation" character varying(10) NOT NULL,
|     "user_name" character varying(50),
|     "timestamp" timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP
| );
|
| ALTER SEQUENCE "public"."audit_log_id_seq" OWNED BY "public"."audit_log"."id";
|
| ALTER TABLE "public"."audit_log" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."audit_log" ADD CONSTRAINT "audit_log_pkey" PRIMARY KEY (id);
|
| CREATE OR REPLACE FUNCTION public.log_account_changes()
|  RETURNS trigger
|  LANGUAGE plpgsql
| AS $function$
| BEGIN
|     INSERT INTO audit_log (table_name, operation, user_name)
|     VALUES ('accounts', TG_OP, current_user);
|     RETURN NEW;
| END;
| $function$
| ;
|
| CREATE OR REPLACE FUNCTION public.update_last_modified()
|  RETURNS trigger
|  LANGUAGE plpgsql
| AS $function$
| BEGIN
|     NEW.last_modified = CURRENT_TIMESTAMP;
|     RETURN NEW;
| END;
| $function$
| ;
|
| CREATE TRIGGER trg_audit_accounts AFTER INSERT OR DELETE OR UPDATE ON public.accounts FOR EACH ROW EXECUTE FUNCTION public.log_account_changes();
|
| CREATE TRIGGER trg_update_last_modified BEFORE UPDATE ON public.accounts FOR EACH ROW EXECUTE FUNCTION public.update_last_modified();
|
|
< ok 0 -

== TestDatabaseDefinitionCorpusRunsOnPostgres/views_and_functions default#1
> exec
|
| SET statement_timeout = 0;
| SET lock_timeout = 0;
| SET idle_in_transaction_session_timeout = 0;
| SET client_encoding = 'UTF8';
| SET standard_conforming_strings = on;
| SELECT pg_catalog.set_config('search_path', '', false);
| SET check_function_bodies = false;
| SET xmloption = content;
| SET client_min_messages = warning;
| SET row_security = off;
|
| COMMENT ON SCHEMA "public" IS 'standard public schema';
|
| CREATE SEQUENCE "public"."employees_id_seq"
|     AS integer
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 2147483647
| 	CACHE 1
| 	NO CYCLE;
|
| SET default_tablespace = '';
|
| CREATE OR REPLACE FUNCTION public.calculate_bonus(emp_id integer)
|  RETURNS numeric
|  LANGUAGE plpgsql
| AS $function$
| DECLARE
|     emp_salary DECIMAL;
| BEGIN
|     SELECT salary INTO emp_salary FROM employees WHERE id = emp_id;
|     RETURN emp_salary * 0.1;
| END;
| $function$
| ;
|
| CREATE TABLE "public"."employees" (
|     "id" integer DEFAULT nextval('public.employees_id_seq'::regclass) NOT NULL,
|     "name" character varying(100) NOT NULL,
|     "department" character varying(50),
|     "salary" numeric(10,2)
| );
|
| ALTER SEQUENCE "public"."employees_id_seq" OWNED BY "public"."employees"."id";
|
| ALTER TABLE "public"."employees" OWNER TO "postgres";
|
| ALTER TABLE ONLY "public"."employees" ADD CONSTRAINT "employees_pkey" PRIMARY KEY (id);
|
| CREATE OR REPLACE FUNCTION public.get_employee_count(dept character varying)
|  RETURNS integer
|  LANGUAGE plpgsql
| AS $function$
| BEGIN
|     RETURN (SELECT COUNT(*) FROM employees WHERE department = dept);
| END;
| $function$
| ;
|
| CREATE VIEW "public"."active_employees" AS 
|  SELECT id,
|     name,
|     department
|    FROM public.employees
|   WHERE (department IS NOT NULL);
|
|
< ok 0 -
