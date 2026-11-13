-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestContainerTable_ColumnChanges default#1
> exec
| DROP SCHEMA IF EXISTS "test_1" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_1"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_2" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_2"
< ok 0 -
> exec
| SET search_path TO "test_1", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_1", public;
| ALTER TABLE t1 ALTER COLUMN name DROP DEFAULT;
| ALTER TABLE t1 ALTER COLUMN name TYPE character varying(200);
| ALTER TABLE t1 ALTER COLUMN name SET DEFAULT 'anon'::character varying
< ok 0 -
> exec
| SET search_path TO "test_2", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(200) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| t1
: string
| r
> query ^
: string
| test_2
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_1
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_1.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query ^
: string
| test_2
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
| idx_name
: string
| CREATE INDEX idx_name ON test_1.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_1.t1 USING btree (id)
> query ^
: string
| test_2
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_2.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_2.t1 USING btree (id)
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
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_2
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_1
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_2
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| test_1
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_2
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
| test_1
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
| test_2
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
| test_1
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query ^
: string
| test_2
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_1
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_2
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
| test_1
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
> query ^
: string
| test_2
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
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
| test_1
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_2
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
| test_1
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
| test_2
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_2" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_1" CASCADE
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
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_3", public;
| ALTER TABLE t1 ALTER COLUMN age TYPE bigint
< ok 0 -
> exec
| SET search_path TO "test_4", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age bigint,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| t1
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
| t1
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| bigint
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_3.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| bigint
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
| idx_name
: string
| CREATE INDEX idx_name ON test_3.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_3.t1 USING btree (id)
> query ^
: string
| test_4
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_4.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_4.t1 USING btree (id)
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
| t1_pkey
: string
| p
: string
| t1
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
| t1_pkey
: string
| p
: string
| t1
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
< row
: string
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
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
< row
: string
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| status
: string
| active,inactive
> query ^
: string
| test_4
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
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
| t1
: string
| Test table
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
| t1
: string
| Test table
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
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status text DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_5", public;
| ALTER TABLE t1 ALTER COLUMN status DROP DEFAULT;
| ALTER TABLE t1 ALTER COLUMN status TYPE status USING status::status;
| ALTER TABLE t1 ALTER COLUMN status SET DEFAULT 'active'::status
< ok 0 -
> exec
| SET search_path TO "test_6", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| test_5
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_6
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
< row
: string
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_5.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
< row
: string
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_5.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_5.t1 USING btree (id)
> query ^
: string
| test_6
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_6.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_6.t1 USING btree (id)
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
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_6
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
< row
: string
| status
: string
| active,inactive
> query ^
: string
| test_6
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
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
< row
: string
| TABLE
: string
| t1
: string
| Test table
> query ^
: string
| test_6
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
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
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_7", public;
| ALTER TABLE t1 ALTER COLUMN email SET NOT NULL
< ok 0 -
> exec
| SET search_path TO "test_8", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text NOT NULL,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| t1
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
| t1
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
: string
| email
: string
| text
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_7.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
: string
| email
: string
| text
: string
| NO
: null
: null
: null
: int64 3
< row
: string
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
| idx_name
: string
| CREATE INDEX idx_name ON test_7.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_7.t1 USING btree (id)
> query ^
: string
| test_8
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_8.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_8.t1 USING btree (id)
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
| t1_pkey
: string
| p
: string
| t1
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
| t1_pkey
: string
| p
: string
| t1
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
< row
: string
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
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
< row
: string
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| status
: string
| active,inactive
> query ^
: string
| test_8
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
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
| t1
: string
| Test table
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
| t1
: string
| Test table
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
| DROP SCHEMA IF EXISTS "test_10" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_10"
< ok 0 -
> exec
| SET search_path TO "test_9", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_9", public;
| ALTER TABLE t1 ALTER COLUMN name DROP NOT NULL
< ok 0 -
> exec
| SET search_path TO "test_10", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| test_9
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_10
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_9
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| YES
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_9.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query ^
: string
| test_10
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| YES
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_9
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_9.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_9.t1 USING btree (id)
> query ^
: string
| test_10
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_10.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_10.t1 USING btree (id)
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
| test_9
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_10
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| test_9
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_10
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
| test_9
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
| test_10
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
| test_9
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query ^
: string
| test_10
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_9
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_10
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
| test_9
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
> query ^
: string
| test_10
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
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
| test_9
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_10
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
| test_9
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
| test_10
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_10" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_9" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_11" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_11"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_12" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_12"
< ok 0 -
> exec
| SET search_path TO "test_11", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_11", public;
| ALTER TABLE t1 ALTER COLUMN email SET DEFAULT 'unknown@example.com'::text
< ok 0 -
> exec
| SET search_path TO "test_12", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text DEFAULT 'unknown@example.com',
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| test_11
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_12
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_11
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
: string
| email
: string
| text
: string
| YES
: string
| 'unknown@example.com'::text
: null
: null
: int64 3
< row
: string
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_11.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query ^
: string
| test_12
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
: string
| email
: string
| text
: string
| YES
: string
| 'unknown@example.com'::text
: null
: null
: int64 3
< row
: string
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_11
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_11.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_11.t1 USING btree (id)
> query ^
: string
| test_12
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_12.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_12.t1 USING btree (id)
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
| test_11
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_12
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_11
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_12
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| test_11
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_12
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
| test_11
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
| test_12
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
| test_11
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query ^
: string
| test_12
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
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
> query ^
: string
| test_12
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
| test_11
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
> query ^
: string
| test_12
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
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
| test_11
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_12
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
| test_11
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
| test_12
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_12" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_11" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_13" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_13"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_14" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_14"
< ok 0 -
> exec
| SET search_path TO "test_13", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_13", public;
| ALTER TABLE t1 ALTER COLUMN name SET DEFAULT 'unknown'::character varying
< ok 0 -
> exec
| SET search_path TO "test_14", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'unknown',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| test_13
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_14
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_13
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'unknown'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_13.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query ^
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'unknown'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_13
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_13.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_13.t1 USING btree (id)
> query ^
: string
| test_14
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_14.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_14.t1 USING btree (id)
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
| test_13
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_14
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_13
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| test_13
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_14
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
| test_13
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
| test_14
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
| test_13
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query ^
: string
| test_14
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_13
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_14
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
| test_13
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
> query ^
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
| t1
: string
| Test table
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
| test_13
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_14
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
| test_13
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
| test_14
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_14" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_13" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_15" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_15"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_16" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_16"
< ok 0 -
> exec
| SET search_path TO "test_15", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_15", public;
| ALTER TABLE t1 ALTER COLUMN name DROP DEFAULT
< ok 0 -
> exec
| SET search_path TO "test_16", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL,
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| test_15
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_16
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_15.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query ^
: string
| test_16
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: null
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_15
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_15.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_15.t1 USING btree (id)
> query ^
: string
| test_16
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_16.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_16.t1 USING btree (id)
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
| test_15
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_16
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_16
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| test_15
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_16
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
| test_15
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
| test_16
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
| test_15
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query ^
: string
| test_16
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_15
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_16
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
| test_15
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
> query ^
: string
| test_16
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
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
| test_15
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_16
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
| test_15
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
| test_16
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_16" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_15" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_17" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_17"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_18" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_18"
< ok 0 -
> exec
| SET search_path TO "test_17", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_17", public;
| ALTER TABLE t1 ADD COLUMN bio text NOT NULL DEFAULT ''::text
< ok 0 -
> exec
| SET search_path TO "test_18", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now(),
|     bio text NOT NULL DEFAULT ''
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| test_17
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_18
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_17
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_17.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
| t1
: string
| bio
: string
| text
: string
| NO
: string
| ''::text
: null
: null
: int64 9
> query ^
: string
| test_18
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| age
: string
| integer
: string
| YES
: null
: null
: null
: int64 4
< row
: string
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
| t1
: string
| bio
: string
| text
: string
| NO
: string
| ''::text
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
| test_17
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_17.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_17.t1 USING btree (id)
> query ^
: string
| test_18
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_18.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_18.t1 USING btree (id)
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
| test_17
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_18
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_17
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_18
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| test_17
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_18
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
| test_17
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
| test_18
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
| test_17
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query ^
: string
| test_18
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_17
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_18
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
| test_17
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
> query ^
: string
| test_18
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
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
| test_17
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_18
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
| test_17
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
| test_18
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_18" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_17" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_19" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_19"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_20" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_20"
< ok 0 -
> exec
| SET search_path TO "test_19", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     age integer,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
|
< ok 0 -
> exec
| SET search_path TO "test_19", public;
| ALTER TABLE t1 DROP COLUMN age
< ok 0 -
> exec
| SET search_path TO "test_20", public;
|
| CREATE TYPE status AS ENUM ('active', 'inactive');
| CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon',
|     email text,
|     status status DEFAULT 'active',
|     score numeric(10,2) DEFAULT 0.0,
|     active boolean DEFAULT true,
|     created_at timestamptz DEFAULT now()
| );
| CREATE INDEX idx_name ON t1(name);
| COMMENT ON TABLE t1 IS 'Test table';
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
| test_19
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_20
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_19
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::test_19.status
: null
: null
: int64 5
< row
: string
| t1
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
: int64 6
< row
: string
| t1
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
: int64 7
< row
: string
| t1
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
> query ^
: string
| test_20
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
< row
: string
| t1
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
| t1
: string
| status
: string
| USER-DEFINED
: string
| YES
: string
| 'active'::status
: null
: null
: int64 4
< row
: string
| t1
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
| t1
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
| t1
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
: int64 7
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_19
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_19.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_19.t1 USING btree (id)
> query ^
: string
| test_20
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| idx_name
: string
| CREATE INDEX idx_name ON test_20.t1 USING btree (name)
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_20.t1 USING btree (id)
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
| test_19
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_20
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_19
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
: bool false
: string
| invoker
: bool false
: string
| u
> query ^
: string
| test_20
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
| check_positive
: string
| boolean
: string
| integer
: string
| SELECT $1 > 0
: string
| sql
: string
| i
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
| test_19
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_20
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
| test_19
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
| test_20
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
| test_19
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query ^
: string
| test_20
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
< row
: string
| status
: string
| active,inactive
> query
|
| 		SELECT viewname, definition, NULL::text
| 		FROM pg_views
| 		WHERE schemaname = $1
| 		ORDER BY viewname
: string
| test_19
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_20
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
| test_19
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
> query ^
: string
| test_20
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
< row
: string
| TABLE
: string
| t1
: string
| Test table
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
| test_19
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_20
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
| test_19
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
| test_20
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_20" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_19" CASCADE
< ok 0 -
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
|
| CREATE TABLE t1 (
|     id integer,
|     name varchar(100) NOT NULL DEFAULT 'anon'
| );
|
< ok 0 -
> exec
| SET search_path TO "test_21", public;
| ALTER TABLE t1 ALTER COLUMN id SET NOT NULL;
| ALTER TABLE t1 ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY
< ok 0 -
> exec
| SET search_path TO "test_22", public;
|
| CREATE TABLE t1 (
|     id integer GENERATED ALWAYS AS IDENTITY,
|     name varchar(100) NOT NULL DEFAULT 'anon'
| );
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
| test_21
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_22
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
< row
: string
| t1
: string
| id
: string
| integer
: string
| NO
: null
: string
| ALWAYS
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
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
< row
: string
| t1
: string
| id
: string
| integer
: string
| NO
: null
: string
| ALWAYS
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
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
| t1_id_seq
: string
| integer
: int64 1
: int64 1
: int64 1
: int64 2147483647
: bool false
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
| t1_id_seq
: string
| integer
: int64 1
: int64 1
: int64 1
: int64 2147483647
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
| test_21
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_22
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
> exec
| DROP SCHEMA IF EXISTS "test_23" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_23"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_24" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_24"
< ok 0 -
> exec
| SET search_path TO "test_23", public;
|
| CREATE TABLE t1 (
|     id integer GENERATED ALWAYS AS IDENTITY,
|     name varchar(100) NOT NULL DEFAULT 'anon'
| );
|
< ok 0 -
> exec
| SET search_path TO "test_23", public;
| ALTER TABLE t1 ALTER COLUMN id DROP IDENTITY;
| ALTER TABLE t1 ALTER COLUMN id DROP NOT NULL
< ok 0 -
> exec
| SET search_path TO "test_24", public;
|
| CREATE TABLE t1 (
|     id integer,
|     name varchar(100) NOT NULL DEFAULT 'anon'
| );
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
| test_23
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_24
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_23
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
| t1
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
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
> query ^
: string
| test_24
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
| t1
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
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_23
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
> query ^
: string
| test_24
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
| test_23
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query ^
: string
| test_24
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
| test_23
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
| test_24
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
| test_23
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_24
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
| test_23
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
| test_24
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
| test_23
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_24
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
| test_23
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_24
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
| test_23
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_24
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
| test_23
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_24
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
| test_23
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
| test_24
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_24" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_23" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_25" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_25"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_26" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_26"
< ok 0 -
> exec
| SET search_path TO "test_25", public;
|
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name text COLLATE "C" NOT NULL DEFAULT 'anon'
| );
|
< ok 0 -
> exec
| SET search_path TO "test_26", public;
|
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name text COLLATE "POSIX" NOT NULL DEFAULT 'anon'
| );
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
| test_25
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_26
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_25
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| text
: string
| NO
: string
| 'anon'::text
: null
: null
: int64 2
> query ^
: string
| test_26
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| text
: string
| NO
: string
| 'anon'::text
: null
: null
: int64 2
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_25
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_25.t1 USING btree (id)
> query ^
: string
| test_26
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_26.t1 USING btree (id)
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
| test_25
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_26
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_25
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
| test_26
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
| test_25
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_26
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
| test_25
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
| test_26
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
| test_25
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_26
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
| test_25
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_26
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
| test_25
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_26
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
| test_25
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_26
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
| test_25
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
| test_26
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_26" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_25" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_27" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_27"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_28" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_28"
< ok 0 -
> exec
| SET search_path TO "test_27", public;
|
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     price numeric(10,2) DEFAULT 0.0,
|     tax numeric(10,2) DEFAULT 0.0
| );
|
< ok 0 -
> exec
| SET search_path TO "test_27", public;
| ALTER TABLE t1 ADD COLUMN total numeric GENERATED ALWAYS AS ((price + tax)) STORED
< ok 0 -
> exec
| SET search_path TO "test_28", public;
|
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     price numeric(10,2) DEFAULT 0.0,
|     tax numeric(10,2) DEFAULT 0.0,
|     total numeric GENERATED ALWAYS AS (price + tax) STORED
| );
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
| test_27
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_28
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_27
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| price
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 2
< row
: string
| t1
: string
| tax
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 3
< row
: string
| t1
: string
| total
: string
| numeric
: string
| YES
: null
: null
: string
| (price + tax)
: int64 4
> query ^
: string
| test_28
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| price
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 2
< row
: string
| t1
: string
| tax
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 3
< row
: string
| t1
: string
| total
: string
| numeric
: string
| YES
: null
: null
: string
| (price + tax)
: int64 4
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_27
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_27.t1 USING btree (id)
> query ^
: string
| test_28
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_28.t1 USING btree (id)
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
| test_27
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_28
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_27
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
| test_28
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
| test_27
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_28
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
| test_27
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
| test_28
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
| test_27
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_28
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
| test_27
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_28
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
| test_27
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_28
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
| test_27
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_28
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
| test_27
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
| test_28
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_28" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_27" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_29" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_29"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_30" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_30"
< ok 0 -
> exec
| SET search_path TO "test_29", public;
|
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     price numeric(10,2) DEFAULT 0.0,
|     tax numeric(10,2) DEFAULT 0.0,
|     total numeric GENERATED ALWAYS AS (price + tax) STORED
| );
|
< ok 0 -
> exec
| SET search_path TO "test_29", public;
| ALTER TABLE t1 DROP COLUMN total
< ok 0 -
> exec
| SET search_path TO "test_30", public;
|
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     price numeric(10,2) DEFAULT 0.0,
|     tax numeric(10,2) DEFAULT 0.0
| );
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
| test_29
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_30
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_29
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| price
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 2
< row
: string
| t1
: string
| tax
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 3
> query ^
: string
| test_30
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| price
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 2
< row
: string
| t1
: string
| tax
: string
| numeric
: string
| YES
: string
| 0.0
: null
: null
: int64 3
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_29
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_29.t1 USING btree (id)
> query ^
: string
| test_30
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_30.t1 USING btree (id)
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
| test_29
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_30
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_29
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
| test_30
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
| test_29
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_30
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
| test_29
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
| test_30
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
| test_29
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_30
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
| test_29
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_30
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
| test_29
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_30
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
| test_29
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_30
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
| test_29
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
| test_30
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_30" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_29" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_31" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_31"
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_32" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_32"
< ok 0 -
> exec
| SET search_path TO "test_31", public;
|
| CREATE UNLOGGED TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon'
| );
|
< ok 0 -
> exec
| SET search_path TO "test_32", public;
|
| CREATE TABLE t1 (
|     id integer PRIMARY KEY,
|     name varchar(100) NOT NULL DEFAULT 'anon'
| );
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
| test_31
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
: string
| r
> query ^
: string
| test_32
< result
< column "relname" "NAME"
< column "relkind" "TEXT"
< row
: string
| t1
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
| test_31
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
> query ^
: string
| test_32
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
| t1
: string
| id
: string
| integer
: string
| NO
: null
: null
: null
: int64 1
< row
: string
| t1
: string
| name
: string
| character varying
: string
| NO
: string
| 'anon'::character varying
: null
: null
: int64 2
> query
|
| 		SELECT indexname, indexdef
| 		FROM pg_indexes
| 		WHERE schemaname = $1
| 		ORDER BY indexname
: string
| test_31
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_31.t1 USING btree (id)
> query ^
: string
| test_32
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t1_pkey
: string
| CREATE UNIQUE INDEX t1_pkey ON test_32.t1 USING btree (id)
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
| test_31
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
: string
| PRIMARY KEY (id)
> query ^
: string
| test_32
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< row
: string
| t1_pkey
: string
| p
: string
| t1
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
| test_31
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
| test_32
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
| test_31
< result
< column "tgname" "NAME"
< column "relname" "NAME"
< column "pg_get_triggerdef" "TEXT"
> query ^
: string
| test_32
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
| test_31
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
| test_32
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
| test_31
< result
< column "typname" "NAME"
< column "string_agg" "TEXT"
> query ^
: string
| test_32
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
| test_31
< result
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "text" "TEXT"
> query ^
: string
| test_32
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
| test_31
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_32
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
| test_31
< result
< column "object_type" "TEXT"
< column "object_name" "TEXT"
< column "description" "TEXT"
> query ^
: string
| test_32
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
| test_31
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
| test_32
< result
< column "polname" "NAME"
< column "relname" "NAME"
< column "polcmd" "TEXT"
< column "case" "TEXT"
< column "array_to_string" "TEXT"
< column "pg_get_expr" "TEXT"
< column "pg_get_expr" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "test_32" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_31" CASCADE
< ok 0 -
