-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestContainerFullyLoaded default#1
> exec
| DROP SCHEMA IF EXISTS "test_1" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_1"
< ok 0 -
> exec
| SET search_path TO "test_1", public;
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE DOMAIN positive_int AS integer NOT NULL CHECK (VALUE > 0);
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     age positive_int,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT email_valid CHECK (email LIKE '%@%')
| );
|
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON COLUMN users.name IS 'Display name';
| COMMENT ON COLUMN users.email IS 'Contact email';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_select ON users FOR SELECT USING (active = true);
|
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_1
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| users
: string
| r
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_1
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_users_active
: string
| CREATE INDEX idx_users_active ON test_1.users USING btree (active) WHERE (active = true)
< row
: string
| idx_users_email_lower
: string
| CREATE INDEX idx_users_email_lower ON test_1.users USING btree (lower(email))
< row
: string
| idx_users_name
: string
| CREATE INDEX idx_users_name ON test_1.users USING btree (name)
< row
: string
| users_email_key
: string
| CREATE UNIQUE INDEX users_email_key ON test_1.users USING btree (email)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_1.users USING btree (id)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_1
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| email_valid
: string
| c
: string
| users
: string
| CHECK ((email ~~ '%@%'::text))
< row
: string
| users_email_key
: string
| u
: string
| users
: string
| UNIQUE (email)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_1
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_1
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
< row
: string
| users_select
: string
| users
: string
| r
: string
| PERMISSIVE
: string
|
: string
| (active = true)
: null
> exec
| DROP SCHEMA IF EXISTS "test_1" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_2" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_2"
< ok 0 -
> exec
| SET search_path TO "test_2", public;
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE DOMAIN positive_int AS integer NOT NULL CHECK (VALUE > 0);
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     age positive_int,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT email_valid CHECK (email LIKE '%@%')
| );
|
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON COLUMN users.name IS 'Display name';
| COMMENT ON COLUMN users.email IS 'Contact email';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_select ON users FOR SELECT USING (active = true);
|
< ok 0 -
> exec
| SET search_path TO "test_2", public;
| DROP TABLE users CASCADE;
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_2
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_2" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_3" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_3"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_4" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_4"
< ok 0 -
> exec
| SET search_path TO "test_3", public;
| CREATE TYPE "status_type" AS ENUM ('active', 'inactive', 'pending');
| CREATE SEQUENCE "user_seq"
|     INCREMENT BY 1
|     MINVALUE 1
|     MAXVALUE 999999
|     START WITH 100
|     CACHE 10
|     NO CYCLE;
| CREATE TABLE "users" (
|     "id" integer NOT NULL DEFAULT nextval('user_seq'::regclass),
|     "name" character varying(100) NOT NULL DEFAULT 'anonymous'::character varying,
|     "email" text,
|     "status" status_type DEFAULT 'active'::status_type,
|     "score" numeric(10,2) DEFAULT 0.0,
|     "active" boolean DEFAULT true,
|     "metadata" jsonb DEFAULT '{}'::jsonb,
|     "created_at" timestamp with time zone DEFAULT now(),
|     "tags" text[] DEFAULT '{}'::text[],
|     CONSTRAINT "users_email_key" UNIQUE ("email"),
|     CONSTRAINT "users_pkey" PRIMARY KEY ("id")
| );
| CREATE INDEX idx_users_active ON users (active) WHERE active = true;
| CREATE INDEX idx_users_email_lower ON users (lower(email));
| CREATE INDEX idx_users_name ON users (name);
| COMMENT ON TABLE "users" IS 'Core user accounts'
< ok 0 -
> exec
| SET search_path TO "test_4", public;
|
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id)
| );
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
| COMMENT ON TABLE users IS 'Core user accounts';
|
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_3
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| users
: string
| r
> query ^
: string
| test_4
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| users
: string
| r
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_3
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_3.user_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| name
: string
| character varying
: string
| NO
: string
| 'anonymous'::character varying
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_3.status_type
: null
: null
: int64 4
< row
: string
| users
: string
| score
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 5
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| metadata
: string
| jsonb
: string
| YES
: string
| '{}'::jsonb
: null
: null
: int64 7
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 8
< row
: string
| users
: string
| tags
: string
| ARRAY
: string
| YES
: string
| '{}'::text[]
: null
: null
: int64 9
> query ^
: string
| test_4
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('user_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| name
: string
| character varying
: string
| NO
: string
| 'anonymous'::character varying
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status_type
: null
: null
: int64 4
< row
: string
| users
: string
| score
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 5
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| metadata
: string
| jsonb
: string
| YES
: string
| '{}'::jsonb
: null
: null
: int64 7
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 8
< row
: string
| users
: string
| tags
: string
| ARRAY
: string
| YES
: string
| '{}'::text[]
: null
: null
: int64 9
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_3
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_users_active
: string
| CREATE INDEX idx_users_active ON test_3.users USING btree (active) WHERE (active = true)
< row
: string
| idx_users_email_lower
: string
| CREATE INDEX idx_users_email_lower ON test_3.users USING btree (lower(email))
< row
: string
| idx_users_name
: string
| CREATE INDEX idx_users_name ON test_3.users USING btree (name)
< row
: string
| users_email_key
: string
| CREATE UNIQUE INDEX users_email_key ON test_3.users USING btree (email)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_3.users USING btree (id)
> query ^
: string
| test_4
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_users_active
: string
| CREATE INDEX idx_users_active ON test_4.users USING btree (active) WHERE (active = true)
< row
: string
| idx_users_email_lower
: string
| CREATE INDEX idx_users_email_lower ON test_4.users USING btree (lower(email))
< row
: string
| idx_users_name
: string
| CREATE INDEX idx_users_name ON test_4.users USING btree (name)
< row
: string
| users_email_key
: string
| CREATE UNIQUE INDEX users_email_key ON test_4.users USING btree (email)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_4.users USING btree (id)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_3
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| users_email_key
: string
| u
: string
| users
: string
| UNIQUE (email)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
> query ^
: string
| test_4
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| users_email_key
: string
| u
: string
| users
: string
| UNIQUE (email)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_3
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query ^
: string
| test_4
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_3
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_4
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_3
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| user_seq
: string
| bigint
: int64 100
: int64 1
: int64 1
: int64 999999
: bool false
> query ^
: string
| test_4
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| user_seq
: string
| bigint
: int64 100
: int64 1
: int64 1
: int64 999999
: bool false
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_3
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status_type
: string
| active,inactive,pending
> query ^
: string
| test_4
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status_type
: string
| active,inactive,pending
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_3
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_4
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_3
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
> query ^
: string
| test_4
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_3
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_4
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_3
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_4
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_4" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_3" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_5" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_5"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_6" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_6"
< ok 0 -
> exec
| SET search_path TO "test_5", public;
|
| CREATE TABLE users (
|     id integer NOT NULL,
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     active boolean DEFAULT true,
|     CONSTRAINT users_pkey PRIMARY KEY (id)
| );
| CREATE INDEX idx_users_name ON users(name);
|
< ok 0 -
> exec
| SET search_path TO "test_5", public;
| DROP TABLE "users" CASCADE
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_5
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
> query ^
: string
| test_6
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_5
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
> query ^
: string
| test_6
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_5
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
> query ^
: string
| test_6
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_5
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query ^
: string
| test_6
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_5
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query ^
: string
| test_6
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_5
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_6
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_5
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query ^
: string
| test_6
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_5
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_6
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_5
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_6
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_5
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_6
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_5
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_6
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_5
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_6
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_6" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_5" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_7" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_7"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_8" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_8"
< ok 0 -
> exec
| SET search_path TO "test_7", public;
| CREATE TYPE "status_type" AS ENUM ('active', 'inactive', 'pending');
| CREATE SEQUENCE "user_seq"
|     INCREMENT BY 1
|     MINVALUE 1
|     MAXVALUE 999999
|     START WITH 100
|     CACHE 10
|     NO CYCLE;
| CREATE TABLE "users" (
|     "id" integer NOT NULL DEFAULT nextval('user_seq'::regclass),
|     "name" character varying(100) NOT NULL DEFAULT 'anonymous'::character varying,
|     "email" text,
|     "status" status_type DEFAULT 'active'::status_type,
|     "score" numeric(10,2) DEFAULT 0.0,
|     "active" boolean DEFAULT true,
|     "metadata" jsonb DEFAULT '{}'::jsonb,
|     "created_at" timestamp with time zone DEFAULT now(),
|     "tags" text[] DEFAULT '{}'::text[],
|     CONSTRAINT "users_email_key" UNIQUE ("email"),
|     CONSTRAINT "users_pkey" PRIMARY KEY ("id")
| );
| CREATE INDEX idx_users_active ON users (active) WHERE active = true;
| CREATE INDEX idx_users_email_lower ON users (lower(email));
| CREATE INDEX idx_users_name ON users (name);
| COMMENT ON TABLE "users" IS 'Core user accounts'
< ok 0 -
> exec
| SET search_path TO "test_8", public;
|
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id)
| );
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
| COMMENT ON TABLE users IS 'Core user accounts';
|
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_7
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| users
: string
| r
> query ^
: string
| test_8
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| users
: string
| r
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_7
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_7.user_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| name
: string
| character varying
: string
| NO
: string
| 'anonymous'::character varying
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_7.status_type
: null
: null
: int64 4
< row
: string
| users
: string
| score
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 5
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| metadata
: string
| jsonb
: string
| YES
: string
| '{}'::jsonb
: null
: null
: int64 7
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 8
< row
: string
| users
: string
| tags
: string
| ARRAY
: string
| YES
: string
| '{}'::text[]
: null
: null
: int64 9
> query ^
: string
| test_8
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('user_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| name
: string
| character varying
: string
| NO
: string
| 'anonymous'::character varying
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status_type
: null
: null
: int64 4
< row
: string
| users
: string
| score
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 5
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| metadata
: string
| jsonb
: string
| YES
: string
| '{}'::jsonb
: null
: null
: int64 7
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 8
< row
: string
| users
: string
| tags
: string
| ARRAY
: string
| YES
: string
| '{}'::text[]
: null
: null
: int64 9
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_7
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_users_active
: string
| CREATE INDEX idx_users_active ON test_7.users USING btree (active) WHERE (active = true)
< row
: string
| idx_users_email_lower
: string
| CREATE INDEX idx_users_email_lower ON test_7.users USING btree (lower(email))
< row
: string
| idx_users_name
: string
| CREATE INDEX idx_users_name ON test_7.users USING btree (name)
< row
: string
| users_email_key
: string
| CREATE UNIQUE INDEX users_email_key ON test_7.users USING btree (email)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_7.users USING btree (id)
> query ^
: string
| test_8
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_users_active
: string
| CREATE INDEX idx_users_active ON test_8.users USING btree (active) WHERE (active = true)
< row
: string
| idx_users_email_lower
: string
| CREATE INDEX idx_users_email_lower ON test_8.users USING btree (lower(email))
< row
: string
| idx_users_name
: string
| CREATE INDEX idx_users_name ON test_8.users USING btree (name)
< row
: string
| users_email_key
: string
| CREATE UNIQUE INDEX users_email_key ON test_8.users USING btree (email)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_8.users USING btree (id)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_7
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| users_email_key
: string
| u
: string
| users
: string
| UNIQUE (email)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
> query ^
: string
| test_8
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| users_email_key
: string
| u
: string
| users
: string
| UNIQUE (email)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_7
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query ^
: string
| test_8
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_7
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_8
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_7
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| user_seq
: string
| bigint
: int64 100
: int64 1
: int64 1
: int64 999999
: bool false
> query ^
: string
| test_8
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| user_seq
: string
| bigint
: int64 100
: int64 1
: int64 1
: int64 999999
: bool false
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_7
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status_type
: string
| active,inactive,pending
> query ^
: string
| test_8
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status_type
: string
| active,inactive,pending
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_7
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_8
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_7
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
> query ^
: string
| test_8
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_7
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_8
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_7
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_8
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_8" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_7" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_9" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_9"
< ok 0 -
> exec
| SET search_path TO "test_9", public;
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE DOMAIN positive_int AS integer NOT NULL CHECK (VALUE > 0);
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     age positive_int,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT email_valid CHECK (email LIKE '%@%')
| );
|
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON COLUMN users.name IS 'Display name';
| COMMENT ON COLUMN users.email IS 'Contact email';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_select ON users FOR SELECT USING (active = true);
|
| CREATE FUNCTION get_user_count() RETURNS bigint
|     LANGUAGE sql
|     STABLE
|     STRICT
|     SECURITY DEFINER
|     PARALLEL SAFE
|     AS 'SELECT count(*) FROM users';
|
< ok 0 -
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_9
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_count
: string
| bigint
: string
|
: string
| SELECT count(*) FROM users
: string
| sql
: string
| s
: bool true
: string
| definer
: bool false
: string
| s
> exec
| DROP SCHEMA IF EXISTS "test_9" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_10" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_10"
< ok 0 -
> exec
| SET search_path TO "test_10", public;
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE DOMAIN positive_int AS integer NOT NULL CHECK (VALUE > 0);
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     age positive_int,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT email_valid CHECK (email LIKE '%@%')
| );
|
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON COLUMN users.name IS 'Display name';
| COMMENT ON COLUMN users.email IS 'Contact email';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_select ON users FOR SELECT USING (active = true);
|
| CREATE PROCEDURE cleanup_inactive()
|     LANGUAGE plpgsql
|     SECURITY DEFINER
|     AS $$BEGIN DELETE FROM users WHERE active = false; END;$$;
|
< ok 0 -
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_10
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| cleanup_inactive
: null
: string
|
: string
| BEGIN DELETE FROM users WHERE active = false; END;
: string
| plpgsql
: string
| v
: bool false
: string
| definer
: bool false
: string
| u
> exec
| DROP SCHEMA IF EXISTS "test_10" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_11" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_11"
< ok 0 -
> exec
| SET search_path TO "test_11", public;
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE DOMAIN positive_int AS integer NOT NULL CHECK (VALUE > 0);
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     age positive_int,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT email_valid CHECK (email LIKE '%@%')
| );
|
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON COLUMN users.name IS 'Display name';
| COMMENT ON COLUMN users.email IS 'Contact email';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_select ON users FOR SELECT USING (active = true);
|
| CREATE VIEW active_users AS
|     SELECT id, name, email, status
|     FROM users
|     WHERE active = true;
| COMMENT ON VIEW active_users IS 'Only active users';
|
< ok 0 -
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_11
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     name,
|     email,
|     status
|    FROM users
|   WHERE (active = true);
: null
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_11
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
< row
: string
| VIEW
: string
| active_users
: string
| Only active users
> exec
| DROP SCHEMA IF EXISTS "test_11" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_12" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_12"
< ok 0 -
> exec
| SET search_path TO "test_12", public;
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE DOMAIN positive_int AS integer NOT NULL CHECK (VALUE > 0);
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     age positive_int,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT email_valid CHECK (email LIKE '%@%')
| );
|
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON COLUMN users.name IS 'Display name';
| COMMENT ON COLUMN users.email IS 'Contact email';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_select ON users FOR SELECT USING (active = true);
|
| CREATE MATERIALIZED VIEW user_stats AS
|     SELECT status, count(*) AS cnt
|     FROM users
|     GROUP BY status
|     WITH DATA;
| CREATE UNIQUE INDEX idx_user_stats_status ON user_stats(status);
|
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_12
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| user_stats
: string
| m
< row
: string
| users
: string
| r
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_12
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_user_stats_status
: string
| CREATE UNIQUE INDEX idx_user_stats_status ON test_12.user_stats USING btree (status)
< row
: string
| idx_users_active
: string
| CREATE INDEX idx_users_active ON test_12.users USING btree (active) WHERE (active = true)
< row
: string
| idx_users_email_lower
: string
| CREATE INDEX idx_users_email_lower ON test_12.users USING btree (lower(email))
< row
: string
| idx_users_name
: string
| CREATE INDEX idx_users_name ON test_12.users USING btree (name)
< row
: string
| users_email_key
: string
| CREATE UNIQUE INDEX users_email_key ON test_12.users USING btree (email)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_12.users USING btree (id)
> exec
| DROP SCHEMA IF EXISTS "test_12" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_13" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_13"
< ok 0 -
> exec
| SET search_path TO "test_13", public;
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE DOMAIN positive_int AS integer NOT NULL CHECK (VALUE > 0);
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
|
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     age positive_int,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id),
|     CONSTRAINT email_valid CHECK (email LIKE '%@%')
| );
|
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
|
| COMMENT ON TABLE users IS 'Core user accounts';
| COMMENT ON COLUMN users.name IS 'Display name';
| COMMENT ON COLUMN users.email IS 'Contact email';
|
| ALTER TABLE users ENABLE ROW LEVEL SECURITY;
| CREATE POLICY users_select ON users FOR SELECT USING (active = true);
|
| CREATE FUNCTION update_timestamp() RETURNS trigger
|     LANGUAGE plpgsql AS $$BEGIN NEW.created_at = now(); RETURN NEW; END;$$;
|
| CREATE TRIGGER users_update_ts
|     BEFORE UPDATE OF name, email ON users
|     FOR EACH ROW
|     WHEN (OLD.* IS DISTINCT FROM NEW.*)
|     EXECUTE FUNCTION update_timestamp();
|
< ok 0 -
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_13
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| users_update_ts
: string
| users
: string
| CREATE TRIGGER users_update_ts BEFORE UPDATE OF name, email ON test_13.users FOR EACH ROW WHEN ((old.* IS DISTINCT FROM new.*)) EXECUTE FUNCTION update_timestamp()
> exec
| DROP SCHEMA IF EXISTS "test_13" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_14" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_14"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_15" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_15"
< ok 0 -
> exec
| SET search_path TO "test_14", public;
| CREATE TYPE "status_type" AS ENUM ('active', 'inactive', 'pending');
| CREATE SEQUENCE "user_seq"
|     INCREMENT BY 1
|     MINVALUE 1
|     MAXVALUE 999999
|     START WITH 100
|     CACHE 10
|     NO CYCLE;
| CREATE FUNCTION "update_timestamp"()
| RETURNS trigger
| LANGUAGE plpgsql
| AS $$BEGIN NEW.created_at = now(); RETURN NEW; END;$$;
| CREATE TABLE "users" (
|     "id" integer NOT NULL DEFAULT nextval('user_seq'::regclass),
|     "name" character varying(100) NOT NULL DEFAULT 'anonymous'::character varying,
|     "email" text,
|     "status" status_type DEFAULT 'active'::status_type,
|     "score" numeric(10,2) DEFAULT 0.0,
|     "active" boolean DEFAULT true,
|     "metadata" jsonb DEFAULT '{}'::jsonb,
|     "created_at" timestamp with time zone DEFAULT now(),
|     "tags" text[] DEFAULT '{}'::text[],
|     CONSTRAINT "users_email_key" UNIQUE ("email"),
|     CONSTRAINT "users_pkey" PRIMARY KEY ("id")
| );
| CREATE FUNCTION "get_user_count"()
| RETURNS bigint
| LANGUAGE sql
| STABLE
| STRICT
| SECURITY DEFINER
| PARALLEL SAFE
| AS $$SELECT count(*) FROM users$$;
| CREATE VIEW "active_users" AS  SELECT id,
|     name,
|     email,
|     status
|    FROM users
|   WHERE (active = true);
| CREATE INDEX idx_users_active ON users (active) WHERE active = true;
| CREATE INDEX idx_users_email_lower ON users (lower(email));
| CREATE INDEX idx_users_name ON users (name);
| CREATE TRIGGER "users_update_ts" BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE FUNCTION "update_timestamp"();
| COMMENT ON VIEW "active_users" IS 'Only active users';
| COMMENT ON TABLE "users" IS 'Core user accounts'
< ok 0 -
> exec
| SET search_path TO "test_15", public;
|
| CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
| CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
| CREATE TABLE users (
|     id integer NOT NULL DEFAULT nextval('user_seq'),
|     name varchar(100) NOT NULL DEFAULT 'anonymous',
|     email text UNIQUE,
|     status status_type DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     metadata jsonb DEFAULT '{}',
|     created_at timestamptz DEFAULT now(),
|     tags text[] DEFAULT '{}',
|     CONSTRAINT users_pkey PRIMARY KEY (id)
| );
| CREATE INDEX idx_users_name ON users(name);
| CREATE INDEX idx_users_email_lower ON users(lower(email));
| CREATE INDEX idx_users_active ON users(active) WHERE active = true;
| COMMENT ON TABLE users IS 'Core user accounts';
|
| CREATE FUNCTION get_user_count() RETURNS bigint
|     LANGUAGE sql
|     STABLE
|     STRICT
|     SECURITY DEFINER
|     PARALLEL SAFE
|     AS 'SELECT count(*) FROM users';
| CREATE VIEW active_users AS
|     SELECT id, name, email, status
|     FROM users
|     WHERE active = true;
| COMMENT ON VIEW active_users IS 'Only active users';
| CREATE FUNCTION update_timestamp() RETURNS trigger
|     LANGUAGE plpgsql AS $$BEGIN NEW.created_at = now(); RETURN NEW; END;$$;
| CREATE TRIGGER users_update_ts
|     BEFORE UPDATE ON users
|     FOR EACH ROW
|     EXECUTE FUNCTION update_timestamp();
|
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_14
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_users
: string
| v
< row
: string
| users
: string
| r
> query ^
: string
| test_15
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| active_users
: string
| v
< row
: string
| users
: string
| r
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_14
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_users
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_users
: string
| name
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| active_users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| active_users
: string
| status
: string
| USER-DEFINED
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('test_14.user_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| name
: string
| character varying
: string
| NO
: string
| 'anonymous'::character varying
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_14.status_type
: null
: null
: int64 4
< row
: string
| users
: string
| score
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 5
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| metadata
: string
| jsonb
: string
| YES
: string
| '{}'::jsonb
: null
: null
: int64 7
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 8
< row
: string
| users
: string
| tags
: string
| ARRAY
: string
| YES
: string
| '{}'::text[]
: null
: null
: int64 9
> query ^
: string
| test_15
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
< row
: string
| active_users
: string
| id
: string
| integer
: string
| YES
: null
: null
: null
: int64 1
< row
: string
| active_users
: string
| name
: string
| character varying
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| active_users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| active_users
: string
| status
: string
| USER-DEFINED
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| users
: string
| id
: string
| integer
: string
| NO
: string
| nextval('user_seq'::regclass)
: null
: null
: int64 1
< row
: string
| users
: string
| name
: string
| character varying
: string
| NO
: string
| 'anonymous'::character varying
: null
: null
: int64 2
< row
: string
| users
: string
| email
: string
| text
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| users
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status_type
: null
: null
: int64 4
< row
: string
| users
: string
| score
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 5
< row
: string
| users
: string
| active
: string
| boolean
: string
| YES
: string
| true
: null
: null
: int64 6
< row
: string
| users
: string
| metadata
: string
| jsonb
: string
| YES
: string
| '{}'::jsonb
: null
: null
: int64 7
< row
: string
| users
: string
| created_at
: string
| timestamp with time zone
: string
| YES
: string
| now()
: null
: null
: int64 8
< row
: string
| users
: string
| tags
: string
| ARRAY
: string
| YES
: string
| '{}'::text[]
: null
: null
: int64 9
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_14
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_users_active
: string
| CREATE INDEX idx_users_active ON test_14.users USING btree (active) WHERE (active = true)
< row
: string
| idx_users_email_lower
: string
| CREATE INDEX idx_users_email_lower ON test_14.users USING btree (lower(email))
< row
: string
| idx_users_name
: string
| CREATE INDEX idx_users_name ON test_14.users USING btree (name)
< row
: string
| users_email_key
: string
| CREATE UNIQUE INDEX users_email_key ON test_14.users USING btree (email)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_14.users USING btree (id)
> query ^
: string
| test_15
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_users_active
: string
| CREATE INDEX idx_users_active ON test_15.users USING btree (active) WHERE (active = true)
< row
: string
| idx_users_email_lower
: string
| CREATE INDEX idx_users_email_lower ON test_15.users USING btree (lower(email))
< row
: string
| idx_users_name
: string
| CREATE INDEX idx_users_name ON test_15.users USING btree (name)
< row
: string
| users_email_key
: string
| CREATE UNIQUE INDEX users_email_key ON test_15.users USING btree (email)
< row
: string
| users_pkey
: string
| CREATE UNIQUE INDEX users_pkey ON test_15.users USING btree (id)
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_14
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| users_email_key
: string
| u
: string
| users
: string
| UNIQUE (email)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
> query ^
: string
| test_15
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| users_email_key
: string
| u
: string
| users
: string
| UNIQUE (email)
< row
: string
| users_pkey
: string
| p
: string
| users
: string
| PRIMARY KEY (id)
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_14
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_count
: string
| bigint
: string
|
: string
| SELECT count(*) FROM users
: string
| sql
: string
| s
: bool true
: string
| definer
: bool false
: string
| s
< row
: string
| update_timestamp
: string
| trigger
: string
|
: string
| BEGIN NEW.created_at = now(); RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_15
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| get_user_count
: string
| bigint
: string
|
: string
| SELECT count(*) FROM users
: string
| sql
: string
| s
: bool true
: string
| definer
: bool false
: string
| s
< row
: string
| update_timestamp
: string
| trigger
: string
|
: string
| BEGIN NEW.created_at = now(); RETURN NEW; END;
: string
| plpgsql
: string
| v
: bool false
: string
| invoker
: bool false
: string
| u
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_14
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| users_update_ts
: string
| users
: string
| CREATE TRIGGER users_update_ts BEFORE UPDATE ON test_14.users FOR EACH ROW EXECUTE FUNCTION test_14.update_timestamp()
> query ^
: string
| test_15
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
< row
: string
| users_update_ts
: string
| users
: string
| CREATE TRIGGER users_update_ts BEFORE UPDATE ON test_15.users FOR EACH ROW EXECUTE FUNCTION update_timestamp()
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_14
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| user_seq
: string
| bigint
: int64 100
: int64 1
: int64 1
: int64 999999
: bool false
> query ^
: string
| test_15
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| user_seq
: string
| bigint
: int64 100
: int64 1
: int64 1
: int64 999999
: bool false
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_14
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status_type
: string
| active,inactive,pending
> query ^
: string
| test_15
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status_type
: string
| active,inactive,pending
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_14
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     name,
|     email,
|     status
|    FROM test_14.users
|   WHERE (active = true);
: null
> query ^
: string
| test_15
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
< row
: string
| active_users
: string
|  SELECT id,
|     name,
|     email,
|     status
|    FROM users
|   WHERE (active = true);
: null
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_14
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
< row
: string
| VIEW
: string
| active_users
: string
| Only active users
> query ^
: string
| test_15
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| users
: string
| Core user accounts
< row
: string
| VIEW
: string
| active_users
: string
| Only active users
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_14
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_15
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_14
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_15
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_15" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_14" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_16" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_16"
< ok 0 -
> exec
| SET search_path TO "test_16", public;
| CREATE TYPE priority AS ENUM ('low', 'medium', 'high', 'critical');
< ok 0 -
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_16
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority
: string
| low,medium,high,critical
> exec
| DROP SCHEMA IF EXISTS "test_16" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_17" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_17"
< ok 0 -
> exec
| SET search_path TO "test_17", public;
| CREATE DOMAIN email_addr AS text NOT NULL CHECK (VALUE LIKE '%@%');
< ok 0 -
> query
|
| 			SELECT count(*) FROM pg_type t
| 			JOIN pg_namespace n ON n.oid = t.typnamespace
| 			WHERE n.nspname = $1 AND t.typname = 'email_addr' AND t.typtype = 'd'
| 		
: string
| test_17
< result
< column "count" "INT8"
< row
: int64 1
> exec
| DROP SCHEMA IF EXISTS "test_17" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_18" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_18"
< ok 0 -
> exec
| SET search_path TO "test_18", public;
| CREATE TYPE address AS (street text, city text, zip varchar(10));
< ok 0 -
> query
|
| 			SELECT count(*) FROM pg_type t
| 			JOIN pg_namespace n ON n.oid = t.typnamespace
| 			WHERE n.nspname = $1 AND t.typname = 'address' AND t.typtype = 'c'
| 		
: string
| test_18
< result
< column "count" "INT8"
< row
: int64 1
> exec
| DROP SCHEMA IF EXISTS "test_18" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_19" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_19"
< ok 0 -
> exec
| SET search_path TO "test_19", public;
| CREATE TYPE float_range AS RANGE (SUBTYPE = float8);
< ok 0 -
> query
|
| 			SELECT count(*) FROM pg_type t
| 			JOIN pg_namespace n ON n.oid = t.typnamespace
| 			WHERE n.nspname = $1 AND t.typname = 'float_range' AND t.typtype = 'r'
| 		
: string
| test_19
< result
< column "count" "INT8"
< row
: int64 1
> exec
| DROP SCHEMA IF EXISTS "test_19" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_20" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_20"
< ok 0 -
> exec
| SET search_path TO "test_20", public;
| CREATE SEQUENCE global_seq INCREMENT 5 MINVALUE 0 MAXVALUE 1000000 START 0 CACHE 20 CYCLE;
< ok 0 -
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_20
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| global_seq
: string
| bigint
: int64 0
: int64 5
: int64 0
: int64 1000000
: bool true
> exec
| DROP SCHEMA IF EXISTS "test_20" CASCADE
< ok 0 -
> exec
| CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
< ok 0 -
> query
| SELECT count(*) FROM pg_extension WHERE extname = 'uuid-ossp'
< result
< column "count" "INT8"
< row
: int64 1
> exec
| DROP SCHEMA IF EXISTS "test_21" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_21"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_22" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_22"
< ok 0 -
> exec
| SET search_path TO "test_21", public;
| CREATE TYPE "priority" AS ENUM ('low', 'medium', 'high', 'critical');
| CREATE DOMAIN "email_addr" AS text NOT NULL CONSTRAINT "email_addr_check" CHECK (VALUE ~~ '%@%'::text);
| CREATE TYPE "float_range" AS RANGE (SUBTYPE = double precision);
| CREATE TYPE "address" AS ("street" text, "city" text, "zip" character varying(10));
| CREATE SEQUENCE "global_seq"
|     INCREMENT BY 5
|     MINVALUE 0
|     MAXVALUE 1000000
|     START WITH 0
|     CACHE 20
|     CYCLE
< ok 0 -
> exec
| SET search_path TO "test_22", public;
| CREATE TYPE priority AS ENUM ('low', 'medium', 'high', 'critical');
| CREATE DOMAIN email_addr AS text NOT NULL CHECK (VALUE LIKE '%@%');
| CREATE TYPE address AS (street text, city text, zip varchar(10));
| CREATE TYPE float_range AS RANGE (SUBTYPE = float8);
| CREATE SEQUENCE global_seq INCREMENT 5 MINVALUE 0 MAXVALUE 1000000 START 0 CACHE 20 CYCLE;
< ok 0 -
> query
|
| 		SELECT c.relname, c.relkind::text
| 		FROM pg_class c
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND c.relkind IN ('r', 'v', 'm', 'f', 'p')
| 		ORDER BY c.relname
: string
| test_21
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
> query ^
: string
| test_22
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
> query
|
| 		SELECT table_name, column_name, data_type, is_nullable,
| 		       column_default, identity_generation, generation_expression,
| 		       ordinal_position
| 		FROM information_schema.columns
| 		WHERE table_schema = $1
| 		ORDER BY table_name, ordinal_position
: string
| test_21
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
> query ^
: string
| test_22
< result
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "VARCHAR"
< column "is_nullable" "VARCHAR"
< column "column_default" "VARCHAR"
< column "identity_generation" "VARCHAR"
< column "generation_expression" "VARCHAR"
< column "ordinal_position" "INT4"
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_21
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
> query ^
: string
| test_22
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
> query
|
| 		SELECT con.conname, con.contype::text,
| 		       c.relname,
| 		       pg_get_constraintdef(con.oid)
| 		FROM pg_constraint con
| 		JOIN pg_class c ON c.oid = con.conrelid
| 		JOIN pg_namespace n ON n.oid = con.connamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, con.conname
: string
| test_21
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query ^
: string
| test_22
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query
|
| 		SELECT p.proname,
| 		       pg_get_function_result(p.oid),
| 		       pg_get_function_arguments(p.oid),
| 		       p.prosrc,
| 		       l.lanname,
| 		       p.provolatile::text,
| 		       p.proisstrict,
| 		       CASE WHEN p.prosecdef THEN 'definer' ELSE 'invoker' END,
| 		       p.proleakproof,
| 		       p.proparallel::text
| 		FROM pg_proc p
| 		JOIN pg_namespace n ON n.oid = p.pronamespace
| 		JOIN pg_language l ON l.oid = p.prolang
| 		WHERE n.nspname = $1
| 		ORDER BY p.proname, pg_get_function_arguments(p.oid)
: string
| test_21
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| float_multirange
: string
| test_21.float_multirange
: string
|
: string
| multirange_constructor0
: string
| internal
: string
| i
: bool true
: string
| invoker
: bool false
: string
| s
< row
: string
| float_multirange
: string
| test_21.float_multirange
: string
| VARIADIC test_21.float_range[]
: string
| multirange_constructor2
: string
| internal
: string
| i
: bool true
: string
| invoker
: bool false
: string
| s
< row
: string
| float_multirange
: string
| test_21.float_multirange
: string
| test_21.float_range
: string
| multirange_constructor1
: string
| internal
: string
| i
: bool true
: string
| invoker
: bool false
: string
| s
< row
: string
| float_range
: string
| test_21.float_range
: string
| double precision, double precision
: string
| range_constructor2
: string
| internal
: string
| i
: bool false
: string
| invoker
: bool false
: string
| s
< row
: string
| float_range
: string
| test_21.float_range
: string
| double precision, double precision, text
: string
| range_constructor3
: string
| internal
: string
| i
: bool false
: string
| invoker
: bool false
: string
| s
> query ^
: string
| test_22
< result
< column "proname" "NAME"
< column "pg_get_function_result" "TEXT"
< column "pg_get_function_arguments" "TEXT"
< column "prosrc" "TEXT"
< column "lanname" "NAME"
< column "provolatile" "TEXT"
< column "proisstrict" "BOOL"
< column "case" "TEXT"
< column "proleakproof" "BOOL"
< column "proparallel" "TEXT"
< row
: string
| float_multirange
: string
| float_multirange
: string
|
: string
| multirange_constructor0
: string
| internal
: string
| i
: bool true
: string
| invoker
: bool false
: string
| s
< row
: string
| float_multirange
: string
| float_multirange
: string
| VARIADIC float_range[]
: string
| multirange_constructor2
: string
| internal
: string
| i
: bool true
: string
| invoker
: bool false
: string
| s
< row
: string
| float_multirange
: string
| float_multirange
: string
| float_range
: string
| multirange_constructor1
: string
| internal
: string
| i
: bool true
: string
| invoker
: bool false
: string
| s
< row
: string
| float_range
: string
| float_range
: string
| double precision, double precision
: string
| range_constructor2
: string
| internal
: string
| i
: bool false
: string
| invoker
: bool false
: string
| s
< row
: string
| float_range
: string
| float_range
: string
| double precision, double precision, text
: string
| range_constructor3
: string
| internal
: string
| i
: bool false
: string
| invoker
: bool false
: string
| s
> query
|
| 		SELECT tg.tgname, c.relname,
| 		       pg_get_triggerdef(tg.oid)
| 		FROM pg_trigger tg
| 		JOIN pg_class c ON c.oid = tg.tgrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		  AND NOT tg.tgisinternal
| 		ORDER BY c.relname, tg.tgname
: string
| test_21
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_22
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query
|
| 		SELECT s.sequencename, s.data_type,
| 		       s.start_value, s.increment_by,
| 		       s.min_value, s.max_value, s.cycle
| 		FROM pg_sequences s
| 		WHERE s.schemaname = $1
| 		ORDER BY s.sequencename
: string
| test_21
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| global_seq
: string
| bigint
: int64 0
: int64 5
: int64 0
: int64 1000000
: bool true
> query ^
: string
| test_22
< result
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "increment_by" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "cycle" "BOOL"
< row
: string
| global_seq
: string
| bigint
: int64 0
: int64 5
: int64 0
: int64 1000000
: bool true
> query
|
| 		SELECT t.typname,
| 		       string_agg(e.enumlabel, ',' ORDER BY e.enumsortorder)
| 		FROM pg_type t
| 		JOIN pg_namespace n ON n.oid = t.typnamespace
| 		JOIN pg_enum e ON e.enumtypid = t.oid
| 		WHERE n.nspname = $1
| 		GROUP BY t.typname
| 		ORDER BY t.typname
: string
| test_21
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority
: string
| low,medium,high,critical
> query ^
: string
| test_22
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| priority
: string
| low,medium,high,critical
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_21
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_22
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query
|
| 		SELECT
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN
| 					CASE rel.relkind
| 						WHEN 'r' THEN 'TABLE'
| 						WHEN 'v' THEN 'VIEW'
| 						WHEN 'm' THEN 'MATERIALIZED VIEW'
| 						WHEN 'S' THEN 'SEQUENCE'
| 						WHEN 'i' THEN 'INDEX'
| 						ELSE 'RELATION'
| 					END
| 				WHEN 'pg_type'::regclass THEN 'TYPE'
| 				WHEN 'pg_proc'::regclass THEN 'FUNCTION'
| 				WHEN 'pg_namespace'::regclass THEN 'SCHEMA'
| 				ELSE c.classoid::regclass::text
| 			END AS object_type,
| 			CASE c.classoid
| 				WHEN 'pg_class'::regclass THEN rel.relname
| 				WHEN 'pg_type'::regclass THEN typ.typname
| 				WHEN 'pg_proc'::regclass THEN pro.proname
| 				WHEN 'pg_namespace'::regclass THEN nsp.nspname
| 				ELSE c.objoid::text
| 			END AS object_name,
| 			d.description
| 		FROM pg_description d
| 		JOIN (
| 			SELECT objoid, classoid,
| 			       CASE WHEN classoid = 'pg_class'::regclass THEN (SELECT relnamespace FROM pg_class WHERE oid = objoid)
| 			            WHEN classoid = 'pg_type'::regclass THEN (SELECT typnamespace FROM pg_type WHERE oid = objoid)
| 			            WHEN classoid = 'pg_proc'::regclass THEN (SELECT pronamespace FROM pg_proc WHERE oid = objoid)
| 			            WHEN classoid = 'pg_namespace'::regclass THEN objoid
| 			            ELSE 0
| 			       END AS nsoid
| 			FROM pg_description
| 			WHERE objsubid = 0
| 		) c ON c.objoid = d.objoid AND c.classoid = d.classoid
| 		LEFT JOIN pg_class rel ON d.classoid = 'pg_class'::regclass AND d.objoid = rel.oid
| 		LEFT JOIN pg_type typ ON d.classoid = 'pg_type'::regclass AND d.objoid = typ.oid
| 		LEFT JOIN pg_proc pro ON d.classoid = 'pg_proc'::regclass AND d.objoid = pro.oid
| 		LEFT JOIN pg_namespace nsp ON d.classoid = 'pg_namespace'::regclass AND d.objoid = nsp.oid
| 		JOIN pg_namespace fns ON fns.oid = c.nsoid
| 		WHERE fns.nspname = $1
| 		  AND d.objsubid = 0
| 		ORDER BY object_type, object_name
: string
| test_21
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_22
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT 'COLUMN' AS object_type,
| 		       c.relname || '.' || a.attname AS object_name,
| 		       d.description
| 		FROM pg_description d
| 		JOIN pg_class c ON d.objoid = c.oid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = d.objsubid
| 		WHERE n.nspname = $1
| 		  AND d.classoid = 'pg_class'::regclass
| 		  AND d.objsubid > 0
| 		ORDER BY c.relname, a.attnum
: string
| test_21
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_22
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query
|
| 		SELECT pol.polname, c.relname,
| 		       pol.polcmd::text,
| 		       CASE WHEN pol.polpermissive THEN 'PERMISSIVE' ELSE 'RESTRICTIVE' END,
| 		       array_to_string(
| 		           ARRAY(SELECT rolname FROM pg_roles WHERE oid = ANY(pol.polroles)), ','
| 		       ),
| 		       pg_get_expr(pol.polqual, pol.polrelid),
| 		       pg_get_expr(pol.polwithcheck, pol.polrelid)
| 		FROM pg_policy pol
| 		JOIN pg_class c ON c.oid = pol.polrelid
| 		JOIN pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1
| 		ORDER BY c.relname, pol.polname
: string
| test_21
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> query ^
: string
| test_22
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_22" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_21" CASCADE
< ok 0 -
