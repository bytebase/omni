-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestSyncViewDependencyQuoting default#1
> exec
|
| CREATE SCHEMA "app's";
| CREATE TABLE "app's".t (id INT8 PRIMARY KEY);
| CREATE VIEW "app's"."owner's_view" AS SELECT id FROM "app's".t;
| CREATE VIEW "app's".pg_indexes AS SELECT id FROM "app's".t;
| CREATE TABLE "app's".pg_class (id INT8 PRIMARY KEY);
< ok 0 -
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "TEXT"
< column "datcollate" "TEXT"
< row
: string
| defaultdb
: string
| UTF8
: string
| en_US.utf8
< row
: string
| omni_testsyncviewdependencyquoting
: string
| UTF8
: string
| en_US.utf8
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
< row
: string
| system
: string
| UTF8
: string
| en_US.utf8
> query
| SHOW search_path
< result
< column "search_path" "TEXT"
< row
: string
| "$user", public
> begin
< ok
> query
|
| SELECT nspname
| FROM pg_catalog.pg_namespace
| WHERE nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY nspname;
|
< result
< column "nspname" "NAME"
< row
: string
| app's
< row
: string
| public
> query
|
| SELECT
| 	cols.table_schema,
| 	cols.table_name,
| 	cols.column_name,
| 	cols.data_type,
| 	cols.character_maximum_length,
| 	cols.numeric_precision,
| 	cols.numeric_scale,
| 	cols.ordinal_position,
| 	cols.column_default,
| 	cols.is_nullable,
| 	cols.collation_name,
| 	cols.udt_schema,
| 	cols.udt_name,
| 	pg_catalog.col_description(c.oid, cols.ordinal_position::int) as column_comment,
| 	cols.is_hidden
| FROM INFORMATION_SCHEMA.COLUMNS AS cols
| LEFT JOIN pg_namespace n ON n.nspname = cols.table_schema
| LEFT JOIN pg_class c ON c.relname = cols.table_name AND c.relnamespace = n.oid AND c.relkind IN ('r', 'v', 'm', 'p', 'f')
| WHERE cols.table_schema NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY cols.table_schema, cols.table_name, cols.ordinal_position;
< result
< column "table_schema" "TEXT"
< column "table_name" "TEXT"
< column "column_name" "TEXT"
< column "data_type" "TEXT"
< column "character_maximum_length" "INT8"
< column "numeric_precision" "INT8"
< column "numeric_scale" "INT8"
< column "ordinal_position" "INT8"
< column "column_default" "TEXT"
< column "is_nullable" "TEXT"
< column "collation_name" "TEXT"
< column "udt_schema" "TEXT"
< column "udt_name" "TEXT"
< column "column_comment" "TEXT"
< column "is_hidden" "TEXT"
< row
: string
| app's
: string
| owner's_view
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 1
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int8
: null
: string
| NO
< row
: string
| app's
: string
| pg_class
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 1
: null
: string
| NO
: null
: string
| pg_catalog
: string
| int8
: null
: string
| NO
< row
: string
| app's
: string
| pg_indexes
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 1
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int8
: null
: string
| NO
< row
: string
| app's
: string
| t
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 1
: null
: string
| NO
: null
: string
| pg_catalog
: string
| int8
: null
: string
| NO
> query
|
| SELECT idx.schemaname, idx.tablename, idx.indexname, idx.indexdef, (SELECT constraint_type
| 	FROM information_schema.table_constraints
| 	WHERE constraint_schema = idx.schemaname
| 	AND constraint_name = idx.indexname
| 	AND table_schema = idx.schemaname
| 	AND table_name = idx.tablename
| 	AND constraint_type = 'PRIMARY KEY') AS constraint_type,
| 	obj_description(c.oid, 'pg_class') AS comment
| FROM pg_indexes AS idx
| LEFT JOIN pg_namespace n ON n.nspname = idx.schemaname
| LEFT JOIN pg_class c ON c.relname = idx.indexname AND c.relnamespace = n.oid AND c.relkind = 'i'
| WHERE idx.schemaname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY idx.schemaname, idx.tablename, idx.indexname;
< result
< column "schemaname" "NAME"
< column "tablename" "NAME"
< column "indexname" "NAME"
< column "indexdef" "TEXT"
< column "constraint_type" "TEXT"
< column "comment" "TEXT"
< row
: string
| app's
: string
| pg_class
: string
| pg_class_pkey
: string
| CREATE UNIQUE INDEX pg_class_pkey ON omni_testsyncviewdependencyquoting."app's".pg_class USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| app's
: string
| t
: string
| t_pkey
: string
| CREATE UNIQUE INDEX t_pkey ON omni_testsyncviewdependencyquoting."app's".t USING btree (id ASC)
: string
| PRIMARY KEY
: null
> query
|
| SELECT nsp.nspname, rel.relname, con.conname, pg_get_constraintdef(con.oid), obj_description(con.oid, 'pg_constraint')
| FROM pg_catalog.pg_constraint con
| JOIN pg_catalog.pg_class rel ON rel.oid = con.conrelid
| JOIN pg_catalog.pg_namespace nsp ON nsp.oid = con.connamespace
| WHERE con.contype = 'c'
| 	AND nsp.nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY nsp.nspname, rel.relname, con.conname;
< result
< column "nspname" "NAME"
< column "relname" "NAME"
< column "conname" "NAME"
< column "pg_get_constraintdef" "TEXT"
< column "obj_description" "TEXT"
> query
|
| SELECT
| 	n.nspname AS fk_schema,
| 	fk_cls.relname AS fk_table,
| 	conname AS fk_name,
| 	ref_nsp.nspname AS fk_ref_schema,
| 	ref_cls.relname AS fk_ref_table,
| 	confdeltype AS delete_option,
| 	confupdtype AS update_option,
| 	confmatchtype AS match_option,
| 	pg_get_constraintdef(c.oid) AS fk_def
| FROM
| 	pg_constraint c
| 	JOIN pg_namespace n ON n.oid = c.connamespace
| 	JOIN pg_class fk_cls ON fk_cls.oid = c.conrelid
| 	JOIN pg_namespace fk_nsp ON fk_nsp.oid = fk_cls.relnamespace
| 	JOIN pg_class ref_cls ON ref_cls.oid = c.confrelid
| 	JOIN pg_namespace ref_nsp ON ref_nsp.oid = ref_cls.relnamespace
| WHERE
| 	n.nspname NOT IN('crdb_internal','information_schema','pg_catalog','pg_extension')
| 	AND c.contype = 'f'
| ORDER BY fk_schema, fk_table, fk_name;
< result
< column "fk_schema" "NAME"
< column "fk_table" "NAME"
< column "fk_name" "NAME"
< column "fk_ref_schema" "NAME"
< column "fk_ref_table" "NAME"
< column "delete_option" "CHAR"
< column "update_option" "CHAR"
< column "match_option" "CHAR"
< column "fk_def" "TEXT"
> query
| SELECT
| 		foreign_table.foreign_table_schema,
| 		foreign_table.foreign_table_name,
| 		foreign_table.foreign_server_catalog,
| 		foreign_table.foreign_server_name
| 	FROM information_schema.foreign_tables AS foreign_table
| 	ORDER BY foreign_table.foreign_table_schema, foreign_table.foreign_table_name;
< result
< column "foreign_table_schema" "TEXT"
< column "foreign_table_name" "TEXT"
< column "foreign_server_catalog" "TEXT"
< column "foreign_server_name" "TEXT"
> query
|
| 	SELECT tbl.schemaname, tbl.tablename,
| 		GREATEST(pc.reltuples::bigint, 0::BIGINT) AS estimate,
| 		obj_description(pc.oid, 'pg_class') AS comment
| 	FROM pg_catalog.pg_tables tbl
| 	LEFT JOIN pg_namespace n ON n.nspname = tbl.schemaname
| 	LEFT JOIN pg_class pc ON pc.relname = tbl.tablename AND pc.relnamespace = n.oid AND pc.relkind IN ('r', 'p')
| 	WHERE tbl.schemaname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension') AND (pc.relispartition IS NULL OR pc.relispartition IS FALSE)
| 	ORDER BY tbl.schemaname, tbl.tablename;
< result
< column "schemaname" "NAME"
< column "tablename" "NAME"
< column "estimate" "INT8"
< column "comment" "TEXT"
< row
: string
| app's
: string
| pg_class
: int64 0
: null
< row
: string
| app's
: string
| t
: int64 0
: null
> query
|
| SELECT
| 	n.nspname AS schema_name,
| 	c.relname AS table_name,
| 	i2.nspname AS inh_schema_name,
| 	i2.relname AS inh_table_name,
| 	i2.partstrat AS partition_type,
| 	pg_get_expr(c.relpartbound, c.oid) AS rel_part_bound
| FROM
| 	pg_catalog.pg_class c
| 	LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
|     LEFT JOIN (
| 		pg_inherits i 
| 		INNER JOIN pg_class c2 ON i.inhparent = c2.oid 
| 		LEFT JOIN pg_namespace n2 ON n2.oid = c2.relnamespace
| 		LEFT JOIN pg_partitioned_table p ON p.partrelid = c2.oid
| 	) i2 ON i2.inhrelid = c.oid 
| WHERE
| 	((c.relkind = 'r'::"char") OR (c.relkind = 'f'::"char") OR (c.relkind = 'p'::"char"))
| 	AND c.relispartition IS TRUE 
| 	AND n.nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY c.oid;
< result
< column "schema_name" "NAME"
< column "table_name" "NAME"
< column "inh_schema_name" "NAME"
< column "inh_table_name" "NAME"
< column "partition_type" "CHAR"
< column "rel_part_bound" "TEXT"
> query
|
| SELECT v.schemaname, v.viewname, v.definition, obj_description(c.oid, 'pg_class')
| FROM pg_catalog.pg_views v
| LEFT JOIN pg_namespace n ON n.nspname = v.schemaname
| LEFT JOIN pg_class c ON c.relname = v.viewname AND c.relnamespace = n.oid AND c.relkind = 'v'
| WHERE v.schemaname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY v.schemaname, v.viewname;
< result
< column "schemaname" "NAME"
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "obj_description" "TEXT"
< row
: string
| app's
: string
| owner's_view
: string
| SELECT id FROM omni_testsyncviewdependencyquoting."app's".t
: null
< row
: string
| app's
: string
| pg_indexes
: string
| SELECT id FROM omni_testsyncviewdependencyquoting."app's".t
: null
> query
|
| 		SELECT source_ns.nspname as source_schema,
| 	  		source_table.relname as source_table,
| 	  		pg_attribute.attname as column_name
| 	  	FROM pg_depend 
| 	  		JOIN pg_rewrite ON pg_depend.objid = pg_rewrite.oid 
| 	  		JOIN pg_class as dependency_view ON pg_rewrite.ev_class = dependency_view.oid 
| 	  		JOIN pg_class as source_table ON pg_depend.refobjid = source_table.oid 
| 	  		JOIN pg_attribute ON pg_depend.refobjid = pg_attribute.attrelid 
| 	  		    AND pg_depend.refobjsubid = pg_attribute.attnum 
| 	  		JOIN pg_namespace dependency_ns ON dependency_ns.oid = dependency_view.relnamespace
| 	  		JOIN pg_namespace source_ns ON source_ns.oid = source_table.relnamespace
| 	  	WHERE 
| 	  		dependency_ns.nspname = $1
| 	  		AND dependency_view.relname = $2
| 	  		AND pg_attribute.attnum > 0
| 	  		-- Only consider SELECT rules (view definitions), not INSERT/UPDATE/DELETE rules
| 	  		AND pg_rewrite.ev_type = '1'
| 	  	ORDER BY 1,2,3;
| 	
: string
| app's
: string
| owner's_view
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| app's
: string
| t
: string
| id
> query ^
: string
| app's
: string
| pg_indexes
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| app's
: string
| t
: string
| id
> query
|
| SELECT m.schemaname, m.matviewname, m.definition, obj_description(c.oid, 'pg_class')
| FROM pg_catalog.pg_matviews m
| LEFT JOIN pg_namespace n ON n.nspname = m.schemaname
| LEFT JOIN pg_class c ON c.relname = m.matviewname AND c.relnamespace = n.oid AND c.relkind = 'm'
| WHERE m.schemaname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY m.schemaname, m.matviewname;
< result
< column "schemaname" "NAME"
< column "matviewname" "NAME"
< column "definition" "TEXT"
< column "obj_description" "TEXT"
> query
|
| select n.nspname as function_schema,
| 	p.proname as function_name,
| 	p.proname || '(' || pg_get_function_identity_arguments(p.oid) || ')' as function_signature,
| 	case when l.lanname = 'internal' then p.prosrc
| 			when p.prokind = 'p' then (
| 				select s.create_statement from crdb_internal.create_procedure_statements as s
| 				where s.procedure_id + 100000 = p.oid::INT8)
| 			else pg_get_functiondef(p.oid)
| 			end as definition
| from pg_proc p
| left join pg_namespace n on p.pronamespace = n.oid
| left join pg_language l on p.prolang = l.oid
| left join pg_type t on t.oid = p.prorettype 
| where n.nspname not in ('crdb_internal','information_schema','pg_catalog','pg_extension')
| order by function_schema, function_name;
< result
< column "function_schema" "NAME"
< column "function_name" "NAME"
< column "function_signature" "TEXT"
< column "definition" "TEXT"
> query
|
| 		SELECT sequence_schema, sequence_name, data_type, start_value, minimum_value, maximum_value, increment, cycle_option
| 		FROM information_schema.sequences
| 		ORDER BY sequence_schema, sequence_name;
< result
< column "sequence_schema" "TEXT"
< column "sequence_name" "TEXT"
< column "data_type" "TEXT"
< column "start_value" "TEXT"
< column "minimum_value" "TEXT"
< column "maximum_value" "TEXT"
< column "increment" "TEXT"
< column "cycle_option" "TEXT"
> query
|
| SELECT pn.nspname, pt.typname, pe.enumlabel
| FROM pg_enum pe
| JOIN pg_type pt ON pe.enumtypid = pt.oid
| JOIN pg_namespace pn ON pt.typnamespace = pn.oid
| WHERE pn.nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY pn.nspname, pt.typname, pe.enumsortorder;
< result
< column "nspname" "NAME"
< column "typname" "NAME"
< column "enumlabel" "TEXT"
> query
|
| SELECT pn.nspname, pt.typname, a.attname, pg_catalog.format_type(a.atttypid, a.atttypmod)
| FROM pg_type pt
| JOIN pg_class pc ON pc.oid = pt.typrelid AND pc.relkind = 'c'
| JOIN pg_namespace pn ON pn.oid = pt.typnamespace
| JOIN pg_attribute a ON a.attrelid = pc.oid AND a.attnum > 0 AND NOT a.attisdropped
| WHERE pn.nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY pn.nspname, pt.typname, a.attnum;
< result
< column "nspname" "NAME"
< column "typname" "NAME"
< column "attname" "NAME"
< column "format_type" "TEXT"
> query
|
| 		SELECT e.extname, e.extversion, n.nspname, c.description
| 		FROM pg_catalog.pg_extension e
| 		LEFT JOIN pg_catalog.pg_namespace n ON n.oid = e.extnamespace
| 		LEFT JOIN pg_catalog.pg_description c ON c.objoid = e.oid AND c.classoid = 'pg_catalog.pg_extension'::pg_catalog.regclass
| 		WHERE n.nspname != 'pg_catalog'
| 		ORDER BY e.extname;
< result
< column "extname" "NAME"
< column "extversion" "TEXT"
< column "nspname" "NAME"
< column "description" "TEXT"
> commit
< ok
