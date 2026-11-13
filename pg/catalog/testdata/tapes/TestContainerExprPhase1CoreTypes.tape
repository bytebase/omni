-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestContainerExprPhase1CoreTypes default#1
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
| CREATE TABLE t (id int PRIMARY KEY, tags text[]);
< ok 0 -
> exec
| SET search_path TO "test_1", public;
| ALTER TABLE t ADD COLUMN first_tag text GENERATED ALWAYS AS (tags[1]) STORED
< ok 0 -
> exec
| SET search_path TO "test_2", public;
| CREATE TABLE t (id int PRIMARY KEY, tags text[], first_tag text GENERATED ALWAYS AS (tags[1]) STORED);
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
| t
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
| t
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
| t
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
| t
: string
| tags
: string
| ARRAY
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| t
: string
| first_tag
: string
| text
: string
| YES
: null
: null
: string
| tags[1]
: int64 3
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
| t
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
| t
: string
| tags
: string
| ARRAY
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| t
: string
| first_tag
: string
| text
: string
| YES
: null
: null
: string
| tags[1]
: int64 3
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
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_1.t USING btree (id)
> query ^
: string
| test_2
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_2.t USING btree (id)
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
| t_pkey
: string
| p
: string
| t
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
| t_pkey
: string
| p
: string
| t
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
> query ^
: string
| test_2
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
> query ^
: string
| test_2
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
| CREATE TABLE t (id int PRIMARY KEY, tags text[]);
< ok 0 -
> exec
| SET search_path TO "test_3", public;
| ALTER TABLE t ADD CONSTRAINT chk CHECK ((tags[1] IS NOT NULL))
< ok 0 -
> exec
| SET search_path TO "test_4", public;
| CREATE TABLE t (id int PRIMARY KEY, tags text[], CONSTRAINT chk CHECK (tags[1] IS NOT NULL));
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
| t
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
| t
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
| t
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
| t
: string
| tags
: string
| ARRAY
: string
| YES
: null
: null
: null
: int64 2
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
| t
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
| t
: string
| tags
: string
| ARRAY
: string
| YES
: null
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
| test_3
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_3.t USING btree (id)
> query ^
: string
| test_4
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_4.t USING btree (id)
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
| chk
: string
| c
: string
| t
: string
| CHECK ((tags[1] IS NOT NULL))
< row
: string
| t_pkey
: string
| p
: string
| t
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
| chk
: string
| c
: string
| t
: string
| CHECK ((tags[1] IS NOT NULL))
< row
: string
| t_pkey
: string
| p
: string
| t
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
> query ^
: string
| test_4
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
> query ^
: string
| test_4
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
| CREATE TABLE t (id int PRIMARY KEY, matrix int[][]);
< ok 0 -
> exec
| SET search_path TO "test_5", public;
| ALTER TABLE t ADD CONSTRAINT chk CHECK ((matrix[1][1] > 0))
< ok 0 -
> exec
| SET search_path TO "test_6", public;
| CREATE TABLE t (id int PRIMARY KEY, matrix int[][], CONSTRAINT chk CHECK (matrix[1][1] > 0));
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
| t
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
| t
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
| t
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
| t
: string
| matrix
: string
| ARRAY
: string
| YES
: null
: null
: null
: int64 2
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
| t
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
| t
: string
| matrix
: string
| ARRAY
: string
| YES
: null
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
| test_5
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_5.t USING btree (id)
> query ^
: string
| test_6
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_6.t USING btree (id)
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
| chk
: string
| c
: string
| t
: string
| CHECK ((matrix[1][1] > 0))
< row
: string
| t_pkey
: string
| p
: string
| t
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
| chk
: string
| c
: string
| t
: string
| CHECK ((matrix[1][1] > 0))
< row
: string
| t_pkey
: string
| p
: string
| t
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
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, c int, d int);
< ok 0 -
> exec
| SET search_path TO "test_7", public;
| ALTER TABLE t ADD CONSTRAINT chk CHECK (((a, b) < (c, d)))
< ok 0 -
> exec
| SET search_path TO "test_8", public;
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, c int, d int, CONSTRAINT chk CHECK ((a, b) < (c, d)));
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
| t
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
| t
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
| t
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
| t
: string
| a
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| t
: string
| b
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| t
: string
| c
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
| t
: string
| d
: string
| integer
: string
| YES
: null
: null
: null
: int64 5
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
| t
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
| t
: string
| a
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| t
: string
| b
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| t
: string
| c
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
| t
: string
| d
: string
| integer
: string
| YES
: null
: null
: null
: int64 5
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
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_7.t USING btree (id)
> query ^
: string
| test_8
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_8.t USING btree (id)
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
| chk
: string
| c
: string
| t
: string
| CHECK ((ROW(a, b) < ROW(c, d)))
< row
: string
| t_pkey
: string
| p
: string
| t
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
| chk
: string
| c
: string
| t
: string
| CHECK ((ROW(a, b) < ROW(c, d)))
< row
: string
| t_pkey
: string
| p
: string
| t
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
> query ^
: string
| test_8
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
> query ^
: string
| test_8
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
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, c int, d int);
< ok 0 -
> exec
| SET search_path TO "test_9", public;
| ALTER TABLE t ADD CONSTRAINT chk CHECK (((a, b) = (c, d)))
< ok 0 -
> exec
| SET search_path TO "test_10", public;
| CREATE TABLE t (id int PRIMARY KEY, a int, b int, c int, d int, CONSTRAINT chk CHECK ((a, b) = (c, d)));
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
| t
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
| t
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
| t
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
| t
: string
| a
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| t
: string
| b
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| t
: string
| c
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
| t
: string
| d
: string
| integer
: string
| YES
: null
: null
: null
: int64 5
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
| t
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
| t
: string
| a
: string
| integer
: string
| YES
: null
: null
: null
: int64 2
< row
: string
| t
: string
| b
: string
| integer
: string
| YES
: null
: null
: null
: int64 3
< row
: string
| t
: string
| c
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
| t
: string
| d
: string
| integer
: string
| YES
: null
: null
: null
: int64 5
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
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_9.t USING btree (id)
> query ^
: string
| test_10
< result
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< row
: string
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON test_10.t USING btree (id)
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
| chk
: string
| c
: string
| t
: string
| CHECK (((a = c) AND (b = d)))
< row
: string
| t_pkey
: string
| p
: string
| t
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
| chk
: string
| c
: string
| t
: string
| CHECK (((a = c) AND (b = d)))
< row
: string
| t_pkey
: string
| p
: string
| t
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
> query ^
: string
| test_10
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
> query ^
: string
| test_10
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
| CREATE FUNCTION "is_positive"(integer)
| RETURNS boolean
| LANGUAGE sql
| IMMUTABLE
| AS $$SELECT $1 > 0$$
< ok 0 -
> exec
| SET search_path TO "test_12", public;
| CREATE FUNCTION is_positive(int) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
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
> query ^
: string
| test_12
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
> query ^
: string
| test_12
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
| test_11
< result
< column "conname" "NAME"
< column "contype" "TEXT"
< column "relname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query ^
: string
| test_12
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
| is_positive
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
| is_positive
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
> query ^
: string
| test_12
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
> query ^
: string
| test_12
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
