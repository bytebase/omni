-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestSyncForeignTablesWithoutTablePrivilege default#1
> exec
|
| CREATE EXTENSION postgres_fdw;
| CREATE SERVER loopback_server FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'postgres');
| CREATE SCHEMA fdw_schema;
| CREATE FOREIGN TABLE fdw_schema.remote_orders (
|     id integer NOT NULL,
|     note varchar(30)
| ) SERVER loopback_server OPTIONS (schema_name 'public', table_name 'orders');
| CREATE FOREIGN TABLE public.remote_items (
|     item_id bigint
| ) SERVER loopback_server OPTIONS (schema_name 'public', table_name 'items');
|
| CREATE ROLE limited LOGIN;
| GRANT pg_read_all_stats TO limited;
< ok 0 -
> exec
| SET ROLE limited
< ok 0 -
> query
| SELECT count(*) FROM information_schema.foreign_tables
< result
< column "count" "INT8"
< row
: int64 0
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "TEXT"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncforeigntableswithouttableprivilege
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template0
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template1
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
> query
| SHOW server_version_num
< result
< column "server_version_num" "TEXT"
< row
: string
| 170012
> query
| SHOW search_path
< result
< column "search_path" "TEXT"
< row
: string
| "$user", public
> begin
< ok
> exec
| SELECT pg_catalog.set_config('search_path', $1, true);
: string
|
< ok 1 -
> query
|
| SELECT
| 	objid
| FROM
| 	pg_depend
| WHERE
| 	deptype = 'e'
| ORDER BY objid
|
< result
< column "objid" "OID"
< row
: int64 13668
< row
: int64 13669
< row
: int64 13670
< row
: int64 13671
< row
: int64 4294967297
< row
: int64 4294967298
< row
: int64 4294967299
< row
: int64 4294967300
< row
: int64 4294967301
< row
: int64 4294967302
> query
|
| SELECT oid, nspname, pg_catalog.pg_get_userbyid(nspowner) as schema_owner,
|        obj_description(oid, 'pg_namespace') as schema_comment
| FROM pg_catalog.pg_namespace
| WHERE nspname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
|   AND nspname NOT LIKE 'pg_temp%'
|   AND nspname NOT LIKE 'pg_toast%'
| ORDER BY nspname;
|
< result
< column "oid" "OID"
< column "nspname" "NAME"
< column "schema_owner" "NAME"
< column "schema_comment" "TEXT"
< row
: int64 4294967303
: string
| fdw_schema
: string
| postgres
: null
< row
: int64 2200
: string
| public
: string
| pg_database_owner
: string
| standard public schema
> query
|
| SELECT
| 	nc.nspname AS table_schema,
| 	c.relname AS table_name,
| 	a.attname AS column_name,
| 	CASE
| 		WHEN t.typtype = 'd' THEN
| 			CASE
| 				WHEN bt.typelem <> 0 AND bt.typlen = -1 THEN 'ARRAY'
| 				WHEN nbt.nspname = 'pg_catalog' THEN pg_catalog.format_type(t.typbasetype, NULL)
| 				ELSE 'USER-DEFINED'
| 			END
| 		ELSE
| 			CASE
| 				WHEN t.typelem <> 0 AND t.typlen = -1 THEN 'ARRAY'
| 				WHEN nt.nspname = 'pg_catalog' THEN pg_catalog.format_type(a.atttypid, NULL)
| 				ELSE 'USER-DEFINED'
| 			END
| 	END AS data_type,
| 	information_schema._pg_char_max_length(information_schema._pg_truetypid(a.*, t.*), information_schema._pg_truetypmod(a.*, t.*)) AS character_maximum_length,
| 	information_schema._pg_numeric_precision(information_schema._pg_truetypid(a.*, t.*), information_schema._pg_truetypmod(a.*, t.*)) AS numeric_precision,
| 	information_schema._pg_numeric_scale(information_schema._pg_truetypid(a.*, t.*), information_schema._pg_truetypmod(a.*, t.*)) AS numeric_scale,
| 	information_schema._pg_datetime_precision(information_schema._pg_truetypid(a.*, t.*), information_schema._pg_truetypmod(a.*, t.*)) AS datetime_precision,
| 	a.attnum AS ordinal_position,
| 	CASE WHEN a.attgenerated = '' THEN pg_catalog.pg_get_expr(ad.adbin, ad.adrelid) END AS column_default,
| 	a.attgenerated AS generation_kind,
| 	CASE WHEN a.attgenerated <> '' THEN pg_catalog.pg_get_expr(ad.adbin, ad.adrelid) END AS generation_expression,
| 	CASE WHEN a.attnotnull OR (t.typtype = 'd' AND t.typnotnull) THEN 'NO' ELSE 'YES' END AS is_nullable,
| 	co.collname AS collation_name,
| 	COALESCE(nbt.nspname, nt.nspname) AS udt_schema,
| 	COALESCE(bt.typname, t.typname) AS udt_name,
| 	CASE a.attidentity WHEN 'a' THEN 'ALWAYS' WHEN 'd' THEN 'BY DEFAULT' END AS identity_generation,
| 	pg_catalog.col_description(c.oid, a.attnum) AS column_comment
| FROM pg_catalog.pg_attribute a
| 	LEFT JOIN pg_catalog.pg_attrdef ad ON a.attrelid = ad.adrelid AND a.attnum = ad.adnum
| 	JOIN pg_catalog.pg_class c ON a.attrelid = c.oid
| 	JOIN pg_catalog.pg_namespace nc ON c.relnamespace = nc.oid
| 	JOIN pg_catalog.pg_type t ON a.atttypid = t.oid
| 	JOIN pg_catalog.pg_namespace nt ON t.typnamespace = nt.oid
| 	LEFT JOIN (pg_catalog.pg_type bt JOIN pg_catalog.pg_namespace nbt ON bt.typnamespace = nbt.oid)
| 		ON t.typtype = 'd' AND t.typbasetype = bt.oid
| 	LEFT JOIN (pg_catalog.pg_collation co JOIN pg_catalog.pg_namespace nco ON co.collnamespace = nco.oid)
| 		ON a.attcollation = co.oid AND (nco.nspname <> 'pg_catalog' OR co.collname <> 'default')
| WHERE a.attnum > 0
| 	AND NOT a.attisdropped
| 	AND c.relkind IN ('r', 'v', 'f', 'p')
| 	AND nc.nspname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
| 	AND nc.nspname NOT LIKE 'pg_temp%'
| 	AND nc.nspname NOT LIKE 'pg_toast%'
| ORDER BY nc.nspname, c.relname, a.attnum;
< result
< column "table_schema" "NAME"
< column "table_name" "NAME"
< column "column_name" "NAME"
< column "data_type" "TEXT"
< column "character_maximum_length" "INT4"
< column "numeric_precision" "INT4"
< column "numeric_scale" "INT4"
< column "datetime_precision" "INT4"
< column "ordinal_position" "INT2"
< column "column_default" "TEXT"
< column "generation_kind" "CHAR"
< column "generation_expression" "TEXT"
< column "is_nullable" "TEXT"
< column "collation_name" "NAME"
< column "udt_schema" "NAME"
< column "udt_name" "NAME"
< column "identity_generation" "TEXT"
< column "column_comment" "TEXT"
< row
: string
| fdw_schema
: string
| remote_orders
: string
| id
: string
| integer
: null
: int64 32
: int64 0
: null
: int64 1
: null
: string
|
: null
: string
| NO
: null
: string
| pg_catalog
: string
| int4
: null
: null
< row
: string
| fdw_schema
: string
| remote_orders
: string
| note
: string
| character varying
: int64 30
: null
: null
: null
: int64 2
: null
: string
|
: null
: string
| YES
: null
: string
| pg_catalog
: string
| varchar
: null
: null
< row
: string
| public
: string
| remote_items
: string
| item_id
: string
| bigint
: null
: int64 64
: int64 0
: null
: int64 1
: null
: string
|
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int8
: null
: null
> query
|
| SELECT
|   sc.nspname,
|   cc.relname,
|   sp.nspname,
|   cp.relname
| FROM
|   pg_catalog.pg_inherits i
|   left JOIN pg_catalog.pg_class cp ON cp.oid = i.inhparent
|   left join pg_catalog.pg_class cc ON cc.oid = i.inhrelid
|   left join pg_catalog.pg_namespace sp on cp.relnamespace = sp.oid
|   left join pg_catalog.pg_namespace sc on cc.relnamespace = sc.oid
| WHERE (cp.relkind = 'i' or cp.relkind = 'I') and (cc.relkind = 'i' or cc.relkind = 'I')
| ORDER BY sc.nspname, cc.relname, sp.nspname, cp.relname
|
< result
< column "nspname" "NAME"
< column "relname" "NAME"
< column "nspname" "NAME"
< column "relname" "NAME"
> query
|
| SELECT 
|     n.nspname as schema_name,
|     t.relname as table_name,
|     i.relname as index_name,
|     am.amname as index_method,
|     ix.indisunique as is_unique,
|     ix.indisprimary as is_primary,
|     pg_get_indexdef(i.oid) as index_definition,
|     -- Get array of key expressions
|     (SELECT array_agg(pg_get_indexdef(i.oid, k + 1, true) ORDER BY k)
|      FROM generate_subscripts(ix.indkey, 1) as k
|      WHERE k < ix.indnkeyatts  -- Only key columns, not included columns
|     ) as key_expressions,
|     (SELECT array_agg(op.opcname ORDER BY k)
|      FROM generate_subscripts(ix.indkey, 1) as k
|      JOIN pg_opclass op ON ix.indclass[k] = op.oid
|      WHERE k < ix.indnkeyatts  -- Only key columns, not included columns
|     ) as key_opclass_names,
|     (SELECT array_agg(op.opcdefault ORDER BY k)
|      FROM generate_subscripts(ix.indkey, 1) as k
|      JOIN pg_opclass op ON ix.indclass[k] = op.oid
|      WHERE k < ix.indnkeyatts  -- Only key columns, not included columns
|     ) as key_opclass_defaults,
|     -- Check if it's a constraint (primary key, unique, or exclude constraint)
|     CASE
|         WHEN ix.indisprimary THEN true
|         WHEN EXISTS (
|             SELECT 1 FROM pg_constraint c
|             WHERE c.conindid = i.oid AND c.contype IN ('u', 'p', 'x')
|         ) THEN true
|         ELSE false
|     END as is_constraint,
|     obj_description(i.oid) as comment,
|     -- Extract indoption array for sort order information (ASC/DESC, NULLS FIRST/LAST)
|     ix.indoption as index_options
| FROM pg_class i
| JOIN pg_index ix ON i.oid = ix.indexrelid
| JOIN pg_class t ON ix.indrelid = t.oid
| JOIN pg_namespace n ON t.relnamespace = n.oid
| JOIN pg_am am ON i.relam = am.oid
| WHERE n.nspname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
|   AND n.nspname NOT LIKE 'pg_temp%'
|   AND n.nspname NOT LIKE 'pg_toast%'
| ORDER BY n.nspname, t.relname, i.relname;
< result
< column "schema_name" "NAME"
< column "table_name" "NAME"
< column "index_name" "NAME"
< column "index_method" "NAME"
< column "is_unique" "BOOL"
< column "is_primary" "BOOL"
< column "index_definition" "TEXT"
< column "key_expressions" "_TEXT"
< column "key_opclass_names" "_NAME"
< column "key_opclass_defaults" "_BOOL"
< column "is_constraint" "BOOL"
< column "comment" "TEXT"
< column "index_options" "-"
> query
|
| 	SELECT
| 		pt.oid,
| 		pn.nspname as schema_name,
| 		pc.relname as table_name,
| 		pt.tgname as trigger_name,
| 		pg_get_triggerdef(pt.oid) as trigger_def,
| 		obj_description(pt.oid) as trigger_comment
| 	FROM pg_trigger as pt
| 		LEFT JOIN pg_class as pc ON pc.oid = pt.tgrelid
| 		LEFT JOIN pg_namespace as pn ON pn.oid = pc.relnamespace
| 	WHERE pn.nspname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information') AND pt.tgisinternal = false
| 	ORDER BY pn.nspname, pc.relname, pt.tgname;
< result
< column "oid" "OID"
< column "schema_name" "NAME"
< column "table_name" "NAME"
< column "trigger_name" "NAME"
< column "trigger_def" "TEXT"
< column "trigger_comment" "TEXT"
> query
|
| SELECT nsp.nspname, rel.relname, con.conname, pg_get_constraintdef(con.oid, true)
|     FROM pg_catalog.pg_constraint con
|         INNER JOIN pg_catalog.pg_class rel ON rel.oid = con.conrelid
|         INNER JOIN pg_catalog.pg_namespace nsp ON nsp.oid = connamespace
|         WHERE contype = 'c' and nsp.nspname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
|         AND nsp.nspname NOT LIKE 'pg_temp%'
|         AND nsp.nspname NOT LIKE 'pg_toast%'
|         ORDER BY nsp.nspname, rel.relname, con.conname
< result
< column "nspname" "NAME"
< column "relname" "NAME"
< column "conname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query
|
| SELECT nsp.nspname, rel.relname, con.conname, pg_get_constraintdef(con.oid, true)
|     FROM pg_catalog.pg_constraint con
|         INNER JOIN pg_catalog.pg_class rel ON rel.oid = con.conrelid
|         INNER JOIN pg_catalog.pg_namespace nsp ON nsp.oid = connamespace
|         WHERE contype = 'x' and nsp.nspname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
|         AND nsp.nspname NOT LIKE 'pg_temp%'
|         AND nsp.nspname NOT LIKE 'pg_toast%'
|         ORDER BY nsp.nspname, rel.relname, con.conname
< result
< column "nspname" "NAME"
< column "relname" "NAME"
< column "conname" "NAME"
< column "pg_get_constraintdef" "TEXT"
> query
|
| SELECT
| 	n.nspname AS fk_schema,
| 	quote_ident(fk_nsp.nspname) || '.' || quote_ident(fk_cls.relname) AS fk_table,
| 	conname AS fk_name,
| 	ref_nsp.nspname AS fk_ref_schema,
| 	quote_ident(ref_nsp.nspname) || '.' || quote_ident(ref_cls.relname) AS fk_ref_table,
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
| 	n.nspname NOT IN('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
| 	AND n.nspname NOT LIKE 'pg_temp%'
| 	AND n.nspname NOT LIKE 'pg_toast%'
| 	AND c.contype = 'f'
| 	AND c.conparentid = 0
| ORDER BY fk_schema, fk_table, fk_name;
< result
< column "fk_schema" "NAME"
< column "fk_table" "TEXT"
< column "fk_name" "NAME"
< column "fk_ref_schema" "NAME"
< column "fk_ref_table" "TEXT"
< column "delete_option" "CHAR"
< column "update_option" "CHAR"
< column "match_option" "CHAR"
< column "fk_def" "TEXT"
> query
| SELECT
| 		n.nspname AS foreign_table_schema,
| 		c.relname AS foreign_table_name,
| 		current_database() AS foreign_server_catalog,
| 		s.srvname AS foreign_server_name
| 	FROM pg_catalog.pg_foreign_table t
| 		JOIN pg_catalog.pg_class c ON c.oid = t.ftrelid
| 		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_catalog.pg_foreign_server s ON s.oid = t.ftserver
| 	WHERE c.relkind = 'f'
| 	ORDER BY n.nspname, c.relname;
< result
< column "foreign_table_schema" "NAME"
< column "foreign_table_name" "NAME"
< column "foreign_server_catalog" "NAME"
< column "foreign_server_name" "NAME"
< row
: string
| fdw_schema
: string
| remote_orders
: string
| omni_testsyncforeigntableswithouttableprivilege
: string
| loopback_server
< row
: string
| public
: string
| remote_items
: string
| omni_testsyncforeigntableswithouttableprivilege
: string
| loopback_server
> query
|
| 	SELECT pc.oid, tbl.schemaname, tbl.tablename,
| 		pg_table_size(pc.oid),
| 		pg_indexes_size(pc.oid),
| 		GREATEST(pc.reltuples::bigint, 0::BIGINT) AS estimate,
| 		obj_description(pc.oid) AS comment,
| 		tbl.tableowner
| 	FROM pg_catalog.pg_tables tbl
| 	LEFT JOIN pg_namespace n ON n.nspname = tbl.schemaname
| 	LEFT JOIN pg_class pc ON pc.relname = tbl.tablename AND pc.relnamespace = n.oid AND pc.relkind IN ('r', 'p')
| 	WHERE tbl.schemaname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
| 	  AND tbl.schemaname NOT LIKE 'pg_temp%'
| 	  AND tbl.schemaname NOT LIKE 'pg_toast%' AND pc.relispartition IS FALSE
| 	ORDER BY tbl.schemaname, tbl.tablename;
< result
< column "oid" "OID"
< column "schemaname" "NAME"
< column "tablename" "NAME"
< column "pg_table_size" "INT8"
< column "pg_indexes_size" "INT8"
< column "estimate" "INT8"
< column "comment" "TEXT"
< column "tableowner" "NAME"
> query
|
| SELECT
| 	n.nspname AS schema_name,
| 	c.relname AS table_name,
| 	i2.nspname AS inh_schema_name,
| 	i2.relname AS inh_table_name,
| 	i2.partstrat AS partition_type,
| 	pg_get_expr(c.relpartbound, c.oid) AS rel_part_bound,
| 	pg_get_partkeydef(i2.inhparent) AS part_key_def
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
| 	AND n.nspname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
| 	AND n.nspname NOT LIKE 'pg_temp%'
| 	AND n.nspname NOT LIKE 'pg_toast%'
| ORDER BY c.oid;
< result
< column "schema_name" "NAME"
< column "table_name" "NAME"
< column "inh_schema_name" "NAME"
< column "inh_table_name" "NAME"
< column "partition_type" "CHAR"
< column "rel_part_bound" "TEXT"
< column "part_key_def" "TEXT"
> query
|
| SELECT pc.oid, v.schemaname, v.viewname, v.definition, obj_description(pc.oid)
| FROM pg_catalog.pg_views v
| 	LEFT JOIN pg_namespace n ON n.nspname = v.schemaname
| 	LEFT JOIN pg_class pc ON pc.relname = v.viewname AND pc.relnamespace = n.oid AND pc.relkind = 'v'
| WHERE v.schemaname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
|   AND v.schemaname NOT LIKE 'pg_temp%'
|   AND v.schemaname NOT LIKE 'pg_toast%'
| ORDER BY v.schemaname, v.viewname;
< result
< column "oid" "OID"
< column "schemaname" "NAME"
< column "viewname" "NAME"
< column "definition" "TEXT"
< column "obj_description" "TEXT"
> query
|
| 		WITH rule_data AS (
| 			SELECT
| 				n.nspname AS schema_name,
| 				c.relname AS table_name,
| 				r.rulename AS rule_name,
| 				r.ev_type,
| 				r.is_instead,
| 				r.ev_enabled,
| 				pg_get_ruledef(r.oid, true) AS definition
| 			FROM pg_rewrite r
| 			JOIN pg_class c ON c.oid = r.ev_class
| 			JOIN pg_namespace n ON n.oid = c.relnamespace
| 			WHERE r.rulename NOT IN ('_RETURN', '_NOTHING')
| 				AND n.nspname NOT IN ('pg_catalog', 'information_schema', 'pg_toast')
| 		)
| 		SELECT
| 			schema_name,
| 			table_name,
| 			rule_name,
| 			CASE ev_type
| 				WHEN '1' THEN 'SELECT'
| 				WHEN '2' THEN 'UPDATE'
| 				WHEN '3' THEN 'INSERT'
| 				WHEN '4' THEN 'DELETE'
| 			END AS event,
| 			is_instead,
| 			ev_enabled != 'D' AS is_enabled,
| 			CASE
| 				WHEN position(' WHERE ' IN definition) > 0
| 					AND position(' WHERE ' IN definition) < position(' DO ' IN definition) THEN
| 					trim(substring(
| 						definition
| 						FROM position(' WHERE ' IN definition) + 7
| 						FOR position(' DO ' IN definition) - position(' WHERE ' IN definition) - 7
| 					))
| 				ELSE NULL
| 			END AS condition,
| 			definition
| 		FROM rule_data
| 		ORDER BY schema_name, table_name, rule_name;
< result
< column "schema_name" "NAME"
< column "table_name" "NAME"
< column "rule_name" "NAME"
< column "event" "TEXT"
< column "is_instead" "BOOL"
< column "is_enabled" "BOOL"
< column "condition" "TEXT"
< column "definition" "TEXT"
> query
|
| SELECT pc.oid, m.schemaname, m.matviewname, m.definition, obj_description(pc.oid)
| FROM pg_catalog.pg_matviews m
| 	LEFT JOIN pg_namespace n ON n.nspname = m.schemaname
| 	LEFT JOIN pg_class pc ON pc.relname = m.matviewname AND pc.relnamespace = n.oid AND pc.relkind = 'm'
| WHERE m.schemaname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
|   AND m.schemaname NOT LIKE 'pg_temp%'
|   AND m.schemaname NOT LIKE 'pg_toast%'
| ORDER BY m.schemaname, m.matviewname;
< result
< column "oid" "OID"
< column "schemaname" "NAME"
< column "matviewname" "NAME"
< column "definition" "TEXT"
< column "obj_description" "TEXT"
> query
|
| select
| 	p.oid as function_oid,
| 	pt.typrelid as table_oid
| from pg_proc p
| 	left join pg_depend d on p.oid = d.objid
| 	left join pg_type pt on d.refobjid = pt.oid
| 	left join pg_namespace n on p.pronamespace = n.oid
| where n.nspname not in ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
|   AND n.nspname NOT LIKE 'pg_temp%'
|   AND n.nspname NOT LIKE 'pg_toast%'
|   AND p.prokind <> 'a'
|   AND pt.typrelid IS NOT NULL
| order by function_oid, table_oid
|
< result
< column "function_oid" "OID"
< column "table_oid" "OID"
> query
|
| select p.oid, n.nspname as function_schema,
| 	p.proname as function_name,
| 	pg_catalog.pg_get_function_identity_arguments(p.oid) as arguments,
| 	case when l.lanname = 'internal' then p.prosrc
| 			else pg_get_functiondef(p.oid)
| 			end as definition,
| 	pg_catalog.obj_description(p.oid) as comment
| from pg_proc p
| left join pg_namespace n on p.pronamespace = n.oid
| left join pg_language l on p.prolang = l.oid
| left join pg_type t on t.oid = p.prorettype 
| where n.nspname not in ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
|   AND n.nspname NOT LIKE 'pg_temp%'
|   AND n.nspname NOT LIKE 'pg_toast%'
|   AND p.prokind <> 'a'
| order by function_schema, function_name, arguments;
< result
< column "oid" "OID"
< column "function_schema" "NAME"
< column "function_name" "NAME"
< column "arguments" "TEXT"
< column "definition" "TEXT"
< column "comment" "TEXT"
< row
: int64 4294967301
: string
| public
: string
| postgres_fdw_disconnect
: string
| text
: string
| CREATE OR REPLACE FUNCTION public.postgres_fdw_disconnect(text)
|  RETURNS boolean
|  LANGUAGE c
|  PARALLEL RESTRICTED STRICT
| AS '$libdir/postgres_fdw', $function$postgres_fdw_disconnect$function$
|
: null
< row
: int64 4294967302
: string
| public
: string
| postgres_fdw_disconnect_all
: string
|
: string
| CREATE OR REPLACE FUNCTION public.postgres_fdw_disconnect_all()
|  RETURNS boolean
|  LANGUAGE c
|  PARALLEL RESTRICTED STRICT
| AS '$libdir/postgres_fdw', $function$postgres_fdw_disconnect_all$function$
|
: null
< row
: int64 4294967300
: string
| public
: string
| postgres_fdw_get_connections
: string
| OUT server_name text, OUT valid boolean
: string
| CREATE OR REPLACE FUNCTION public.postgres_fdw_get_connections(OUT server_name text, OUT valid boolean)
|  RETURNS SETOF record
|  LANGUAGE c
|  PARALLEL RESTRICTED STRICT
| AS '$libdir/postgres_fdw', $function$postgres_fdw_get_connections$function$
|
: null
< row
: int64 4294967297
: string
| public
: string
| postgres_fdw_handler
: string
|
: string
| CREATE OR REPLACE FUNCTION public.postgres_fdw_handler()
|  RETURNS fdw_handler
|  LANGUAGE c
|  STRICT
| AS '$libdir/postgres_fdw', $function$postgres_fdw_handler$function$
|
: null
< row
: int64 4294967298
: string
| public
: string
| postgres_fdw_validator
: string
| text[], oid
: string
| CREATE OR REPLACE FUNCTION public.postgres_fdw_validator(text[], oid)
|  RETURNS void
|  LANGUAGE c
|  STRICT
| AS '$libdir/postgres_fdw', $function$postgres_fdw_validator$function$
|
: null
> query
|
| 	SELECT
| 		c.oid,
| 		refobjid AS owning_tab,
| 		refobjsubid AS owning_col
| 	FROM pg_class c
|   		LEFT JOIN pg_depend d ON
|   			(c.relkind =  'S' AND
|                 d.classid = 'pg_class'::regclass AND d.objid = c.oid AND
|                 d.objsubid = 0 AND
|                 d.refclassid = 'pg_class'::regclass AND d.deptype IN ('a', 'i'))
| 	WHERE refobjid is NOT NULL and refobjsubid is NOT NULL
| 	ORDER BY c.oid;
< result
< column "oid" "OID"
< column "owning_tab" "OID"
< column "owning_col" "INT4"
> query
|
| 	SELECT
| 		pc.oid,
| 		s.schemaname,
| 		s.sequencename,
| 		s.data_type,
| 		s.start_value,
| 		s.min_value,
| 		s.max_value,
| 		s.increment_by,
| 		s.cycle,
| 		s.cache_size,
| 		s.last_value,
| 		pg_catalog.obj_description(pc.oid) as sequence_comment
| 	FROM pg_sequences s
| 		LEFT JOIN pg_namespace n ON n.nspname = s.schemaname
| 		LEFT JOIN pg_class pc ON pc.relname = s.sequencename AND pc.relnamespace = n.oid AND pc.relkind = 'S'
| 	WHERE s.schemaname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
| 	  AND s.schemaname NOT LIKE 'pg_temp%%'
| 	  AND s.schemaname NOT LIKE 'pg_toast%%'
| 	ORDER BY s.schemaname, s.sequencename;
< result
< column "oid" "OID"
< column "schemaname" "NAME"
< column "sequencename" "NAME"
< column "data_type" "-"
< column "start_value" "INT8"
< column "min_value" "INT8"
< column "max_value" "INT8"
< column "increment_by" "INT8"
< column "cycle" "BOOL"
< column "cache_size" "INT8"
< column "last_value" "INT8"
< column "sequence_comment" "TEXT"
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
< row
: string
| postgres_fdw
: string
| 1.1
: string
| public
: string
| foreign-data wrapper for remote PostgreSQL servers
> query
|
| 	SELECT
| 		pt.oid,
| 		pn.nspname as schema_name,
| 		pt.typname as enum_name,
| 		pe.enumlabel as enum_value,
| 		pg_catalog.obj_description(pt.oid) as enum_comment
| 	FROM pg_enum as pe
| 		LEFT JOIN pg_type as pt ON pe.enumtypid = pt.oid
| 		LEFT JOIN pg_namespace as pn ON pt.typnamespace = pn.oid
| 	WHERE pn.nspname NOT IN ('_timescaledb_cache','_timescaledb_catalog','_timescaledb_config','_timescaledb_internal','citus','citus_internal','columnar','columnar_internal','cron','information_schema','partman','pg_bitmapindex','pg_catalog','pg_toast','rw_catalog','squeeze','timescaledb_experimental','timescaledb_information')
| 	  AND pn.nspname NOT LIKE 'pg_temp%'
| 	  AND pn.nspname NOT LIKE 'pg_toast%'
| 	ORDER BY pn.nspname, pt.typname, pe.enumsortorder;
< result
< column "oid" "OID"
< column "schema_name" "NAME"
< column "enum_name" "NAME"
< column "enum_value" "NAME"
< column "enum_comment" "TEXT"
> query
|
| 	SELECT
| 		pt.oid,
| 		pc.oid,
| 		pn.nspname AS schema_name,
| 		pt.typname AS type_name,
| 		pg_catalog.obj_description(pt.oid) AS type_comment,
| 		a.attname AS attribute_name,
| 		CASE
| 			WHEN en.nspname IS NOT NULL
| 				AND en.nspname NOT IN ('pg_catalog', 'information_schema')
| 				AND position(quote_ident(en.nspname) || '.' IN pg_catalog.format_type(a.atttypid, a.atttypmod)) <> 1
| 			THEN quote_ident(en.nspname) || '.' || pg_catalog.format_type(a.atttypid, a.atttypmod)
| 			ELSE pg_catalog.format_type(a.atttypid, a.atttypmod)
| 		END AS attribute_type,
| 		CASE
| 			WHEN co.collname IS NULL THEN NULL
| 			WHEN nco.nspname <> 'pg_catalog' THEN quote_ident(nco.nspname) || '.' || quote_ident(co.collname)
| 			ELSE quote_ident(co.collname)
| 		END AS attribute_collation,
| 		pg_catalog.col_description(pc.oid, a.attnum) AS attribute_comment
| 	FROM pg_type pt
| 		JOIN pg_class pc ON pc.oid = pt.typrelid AND pc.relkind = 'c'
| 		JOIN pg_namespace pn ON pn.oid = pt.typnamespace
| 		LEFT JOIN pg_attribute a ON a.attrelid = pc.oid AND a.attnum > 0 AND NOT a.attisdropped
| 		LEFT JOIN pg_type at ON at.oid = a.atttypid
| 		LEFT JOIN pg_type et ON et.oid = CASE WHEN at.typlen = -1 AND at.typelem <> 0 THEN at.typelem ELSE at.oid END
| 		LEFT JOIN pg_namespace en ON en.oid = et.typnamespace
| 		LEFT JOIN (pg_catalog.pg_collation co JOIN pg_catalog.pg_namespace nco ON co.collnamespace = nco.oid)
| 			ON a.attcollation = co.oid AND (nco.nspname <> 'pg_catalog' OR co.collname <> 'default')
| 	WHERE pn.nspname NOT LIKE 'pg_temp%'
| 	  AND pn.nspname NOT LIKE 'pg_toast%'
| 	ORDER BY pn.nspname, pt.typname, a.attnum;
< result
< column "oid" "OID"
< column "oid" "OID"
< column "schema_name" "NAME"
< column "type_name" "NAME"
< column "type_comment" "TEXT"
< column "attribute_name" "NAME"
< column "attribute_type" "TEXT"
< column "attribute_collation" "TEXT"
< column "attribute_comment" "TEXT"
> query
|
| 	SELECT
| 		et.oid,
| 		et.evtname AS trigger_name,
| 		et.evtevent AS event_type,
| 		et.evttags AS tags,
| 		n.nspname AS function_schema,
| 		p.proname AS function_name,
| 		et.evtenabled AS enabled,
| 		obj_description(et.oid, 'pg_event_trigger') AS comment
| 	FROM pg_event_trigger et
| 	JOIN pg_proc p ON et.evtfoid = p.oid
| 	JOIN pg_namespace n ON p.pronamespace = n.oid
| 	ORDER BY et.evtname;
< result
< column "oid" "OID"
< column "trigger_name" "NAME"
< column "event_type" "NAME"
< column "tags" "_TEXT"
< column "function_schema" "NAME"
< column "function_name" "NAME"
< column "enabled" "CHAR"
< column "comment" "TEXT"
> commit
< ok
