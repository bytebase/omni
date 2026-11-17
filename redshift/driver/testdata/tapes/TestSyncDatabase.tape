-- pin: driver aws:redshift:1.0.539507 sha256:afc4041fe6e35ce273970308e430eb8245ea9c4bbf597c497bb8eb8bb00810d2 2026-11-17

== TestSyncDatabase driver#1
> query
| SELECT consumer_database FROM SVV_DATASHARES WHERE share_type = 'INBOUND' ORDER BY consumer_database;
< result
< column "consumer_database" "VARCHAR" "string"
> query
|
| 		SELECT datname,
| 		pg_encoding_to_char(encoding)
| 		FROM pg_database
| 		ORDER BY datname;
| 	
< result
< column "datname" "NAME" "string"
< column "pg_encoding_to_char" "NAME" "string"
< row
: string
| awsdatacatalog
: string
| UNICODE
< row
: string
| dev
: string
| UNICODE
< row
: string
| padb_harvest
: string
| UNICODE
< row
: string
| sys:internal
: string
| UNICODE
< row
: string
| template0
: string
| UNICODE
< row
: string
| template1
: string
| UNICODE
> begin
< ok
> query
|
| 		SELECT
| 			schema_name
| 		FROM
| 			SVV_ALL_SCHEMAS
| 		WHERE
| 			database_name = $1
| 			AND substr(schema_name, 1, 3) <> 'pg_' AND schema_name <> 'information_schema'
| 		ORDER BY
| 			schema_name;
| 	
: string
| dev
< result
< column "schema_name" "VARCHAR" "string"
< row
: string
| app
< row
: string
| public
> query
|
| 	SELECT
| 		cols.table_schema,
| 		cols.table_name,
| 		cols.column_name,
| 		cols.data_type,
| 		cols.character_maximum_length,
| 		cols.numeric_precision,
| 		cols.numeric_scale,
| 		cols.ordinal_position,
| 		cols.column_default,
| 		cols.is_nullable,
| 		cols.collation_name,
| 		cols.udt_schema,
| 		cols.udt_name,
| 		pg_catalog.col_description(pc.oid, cols.ordinal_position::int) as column_comment
| 	FROM 
| 		INFORMATION_SCHEMA.COLUMNS AS cols
| 		JOIN pg_namespace AS pns ON pns.nspname = cols.table_schema
| 		LEFT JOIN pg_class AS pc ON pc.relname = cols.table_name AND pns.oid = pc.relnamespace
| 	WHERE substr(cols.table_schema, 1, 3) <> 'pg_' AND cols.table_schema <> 'information_schema'
| 	ORDER BY cols.table_schema, cols.table_name, cols.ordinal_position;
< result
< column "table_schema" "VARCHAR" "string"
< column "table_name" "VARCHAR" "string"
< column "column_name" "VARCHAR" "string"
< column "data_type" "VARCHAR" "string"
< column "character_maximum_length" "INT4" "int32"
< column "numeric_precision" "INT4" "int32"
< column "numeric_scale" "INT4" "int32"
< column "ordinal_position" "INT4" "int32"
< column "column_default" "VARCHAR" "string"
< column "is_nullable" "VARCHAR" "string"
< column "collation_name" "VARCHAR" "string"
< column "udt_schema" "VARCHAR" "string"
< column "udt_name" "VARCHAR" "string"
< column "column_comment" "VARCHAR" "string"
< row
: string
| app
: string
| facts
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
< row
: string
| app
: string
| facts
: string
| lookup_id
: string
| integer
: null
: int64 32
: int64 0
: int64 2
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int4
: null
< row
: string
| app
: string
| facts
: string
| amount
: string
| numeric
: null
: int64 12
: int64 2
: int64 3
: null
: string
| YES
: null
: string
| pg_catalog
: string
| numeric
: null
< row
: string
| app
: string
| facts
: string
| flag
: string
| boolean
: null
: null
: null
: int64 4
: null
: string
| YES
: null
: string
| pg_catalog
: string
| bool
: null
< row
: string
| app
: string
| facts
: string
| seen
: string
| timestamp without time zone
: null
: null
: null
: int64 5
: null
: string
| YES
: null
: string
| pg_catalog
: string
| timestamp
: null
< row
: string
| app
: string
| facts
: string
| notes
: string
| character varying
: int64 80
: null
: null
: int64 6
: null
: string
| YES
: null
: string
| pg_catalog
: string
| varchar
: null
< row
: string
| app
: string
| long_view
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
< row
: string
| app
: string
| long_view
: string
| amount
: string
| numeric
: null
: int64 12
: int64 2
: int64 2
: null
: string
| YES
: null
: string
| pg_catalog
: string
| numeric
: null
< row
: string
| app
: string
| long_view
: string
| notes
: string
| character varying
: int64 80
: null
: null
: int64 3
: null
: string
| YES
: null
: string
| pg_catalog
: string
| varchar
: null
< row
: string
| app
: string
| long_view
: string
| flag
: string
| boolean
: null
: null
: null
: int64 4
: null
: string
| YES
: null
: string
| pg_catalog
: string
| bool
: null
< row
: string
| app
: string
| long_view
: string
| seen
: string
| timestamp without time zone
: null
: null
: null
: int64 5
: null
: string
| YES
: null
: string
| pg_catalog
: string
| timestamp
: null
< row
: string
| app
: string
| long_view
: string
| lookup_id
: string
| integer
: null
: int64 32
: int64 0
: int64 6
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int4
: null
< row
: string
| app
: string
| long_view
: string
| next_id
: string
| bigint
: null
: int64 64
: int64 0
: int64 7
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int8
: null
< row
: string
| app
: string
| long_view
: string
| doubled
: string
| numeric
: null
: int64 14
: int64 2
: int64 8
: null
: string
| YES
: null
: string
| pg_catalog
: string
| numeric
: null
< row
: string
| app
: string
| recent
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
< row
: string
| app
: string
| recent
: string
| amount
: string
| numeric
: null
: int64 12
: int64 2
: int64 2
: null
: string
| YES
: null
: string
| pg_catalog
: string
| numeric
: null
< row
: string
| app
: string
| very_long_view
: string
| identifier_column_one
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
< row
: string
| app
: string
| very_long_view
: string
| amount_column_two
: string
| numeric
: null
: int64 12
: int64 2
: int64 2
: null
: string
| YES
: null
: string
| pg_catalog
: string
| numeric
: null
< row
: string
| app
: string
| very_long_view
: string
| notes_column_three
: string
| character varying
: int64 80
: null
: null
: int64 3
: null
: string
| YES
: null
: string
| pg_catalog
: string
| varchar
: null
< row
: string
| app
: string
| very_long_view
: string
| flag_column_four
: string
| boolean
: null
: null
: null
: int64 4
: null
: string
| YES
: null
: string
| pg_catalog
: string
| bool
: null
< row
: string
| app
: string
| very_long_view
: string
| seen_column_five
: string
| timestamp without time zone
: null
: null
: null
: int64 5
: null
: string
| YES
: null
: string
| pg_catalog
: string
| timestamp
: null
< row
: string
| app
: string
| very_long_view
: string
| lookup_column_six
: string
| integer
: null
: int64 32
: int64 0
: int64 6
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int4
: null
< row
: string
| app
: string
| very_long_view
: string
| next_identifier_seven
: string
| bigint
: null
: int64 64
: int64 0
: int64 7
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int8
: null
< row
: string
| app
: string
| very_long_view
: string
| doubled_amount_eight
: string
| numeric
: null
: int64 14
: int64 2
: int64 8
: null
: string
| YES
: null
: string
| pg_catalog
: string
| numeric
: null
< row
: string
| app
: string
| very_long_view
: string
| tripled_identifier_nine
: string
| bigint
: null
: int64 64
: int64 0
: int64 9
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int8
: null
< row
: string
| app
: string
| very_long_view
: string
| notes_with_suffix_ten
: string
| text
: null
: null
: null
: int64 10
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
< row
: string
| public
: string
| lookup
: string
| id
: string
| integer
: null
: int64 32
: int64 0
: int64 1
: null
: string
| NO
: null
: string
| pg_catalog
: string
| int4
: null
< row
: string
| public
: string
| lookup
: string
| label
: string
| character varying
: int64 20
: null
: null
: int64 2
: string
| 'x'::character varying
: string
| NO
: null
: string
| pg_catalog
: string
| varchar
: string
| a column comment
> query
|
| 	WITH index_positions AS (
| 		SELECT generate_series(0, 31) as pos
| 	)
| 	SELECT 
| 		pgidx.schemaname,
| 		pgidx.tablename, 
| 		pgidx.indexname,
| 		pgidx.indexdef,
| 		ix.indisunique,
| 		ix.indisprimary,
| 		p.pos + 1 as position,
| 		pg_get_indexdef(pc.oid, p.pos + 1, true) as column_expression,
| 		obj_description(pc.oid) AS comment
| 	FROM
| 		pg_indexes AS pgidx 
| 		JOIN pg_namespace AS pns ON pns.nspname = pgidx.schemaname
| 		JOIN pg_class AS pc ON pc.relname = pgidx.indexname AND pns.oid = pc.relnamespace
| 		JOIN pg_index AS ix ON ix.indexrelid = pc.oid
| 		CROSS JOIN index_positions p
| 	WHERE 
| 		substr(pgidx.schemaname, 1, 3) <> 'pg_' AND pgidx.schemaname <> 'information_schema'
| 		AND ix.indkey[p.pos] > 0  -- Only process positions that have actual columns
| 	ORDER BY
| 		pgidx.schemaname, pgidx.tablename, pgidx.indexname, p.pos;
< result
< column "schemaname" "NAME" "string"
< column "tablename" "NAME" "string"
< column "indexname" "NAME" "string"
< column "indexdef" "VARCHAR" "string"
< column "indisunique" "BOOL" "bool"
< column "indisprimary" "BOOL" "bool"
< column "position" "INT4" "int32"
< column "column_expression" "VARCHAR" "string"
< column "comment" "VARCHAR" "string"
< row
: string
| app
: string
| facts
: string
| facts_pkey
: string
| CREATE UNIQUE INDEX facts_pkey ON app.facts USING btree (id)
: bool true
: bool true
: int64 1
: string
| id
: null
< row
: string
| public
: string
| lookup
: string
| lookup_pkey
: string
| CREATE UNIQUE INDEX lookup_pkey ON lookup USING btree (id)
: bool true
: bool true
: int64 1
: string
| id
: null
> query
|
| 	SELECT
| 		n.nspname AS fk_schema,
| 		fk_cls.relname AS fk_table,
| 		conname AS fk_name,
| 		ref_nsp.nspname AS fk_ref_schema,
| 		ref_cls.relname AS fk_ref_table,
| 		confdeltype AS delete_option,
| 		confupdtype AS update_option,
| 		confmatchtype AS match_option,
| 		pg_get_constraintdef(c.oid) AS fk_def
| 	FROM
| 		pg_constraint c
| 		JOIN pg_namespace n ON n.oid = c.connamespace
| 		JOIN pg_class fk_cls ON fk_cls.oid = c.conrelid
| 		JOIN pg_namespace fk_nsp ON fk_nsp.oid = fk_cls.relnamespace
| 		JOIN pg_class ref_cls ON ref_cls.oid = c.confrelid
| 		JOIN pg_namespace ref_nsp ON ref_nsp.oid = ref_cls.relnamespace
| 	WHERE
| 		substr(n.nspname, 1, 3) <> 'pg_' AND n.nspname <> 'information_schema'
| 		AND c.contype = 'f'
| 	ORDER BY fk_schema, fk_table, fk_name;
| 	
< result
< column "fk_schema" "NAME" "string"
< column "fk_table" "NAME" "string"
< column "fk_name" "NAME" "string"
< column "fk_ref_schema" "NAME" "string"
< column "fk_ref_table" "NAME" "string"
< column "delete_option" "CHAR" "string"
< column "update_option" "CHAR" "string"
< column "match_option" "CHAR" "string"
< column "fk_def" "VARCHAR" "string"
< row
: string
| app
: string
| facts
: string
| facts_lookup_id_fkey
: string
| public
: string
| lookup
: string
| a
: string
| a
: string
| u
: string
| FOREIGN KEY (lookup_id) REFERENCES lookup(id)
> query
|
| 	SELECT
| 		ptbl.schemaname,
| 		ptbl.tablename,
| 		0, -- data size
| 		0, -- index size
| 		GREATEST(pc.reltuples::bigint, 0::bigint) AS estimate,
| 		obj_description(pc.oid) AS comment
| 	FROM pg_catalog.pg_tables AS ptbl
| 	JOIN pg_namespace AS pns ON pns.nspname = ptbl.schemaname
| 	LEFT JOIN pg_class AS pc ON pc.relname = ptbl.tablename AND pns.oid = pc.relnamespace
| 	WHERE substr(ptbl.schemaname, 1, 3) <> 'pg_' AND ptbl.schemaname <> 'information_schema'
| 	ORDER BY ptbl.schemaname, ptbl.tablename;
< result
< column "schemaname" "NAME" "string"
< column "tablename" "NAME" "string"
< column "?column?" "INT4" "int32"
< column "?column?" "INT4" "int32"
< column "estimate" "INT8" "int64"
< column "comment" "VARCHAR" "string"
< row
: string
| app
: string
| facts
: int64 0
: int64 0
: int64 0
: null
< row
: string
| public
: string
| lookup
: int64 0
: int64 0
: int64 0
: string
| a synced comment
> query
|
| SELECT
| 	pgv.schemaname,
| 	pgv.viewname,
| 	obj_description(pc.oid) AS comment
| FROM pg_catalog.pg_views AS pgv
| JOIN pg_namespace AS pns ON pns.nspname = pgv.schemaname
| JOIN pg_class AS pc ON pc.relname = pgv.viewname AND pns.oid = pc.relnamespace
| WHERE substr(pgv.schemaname, 1, 3) <> 'pg_' AND pgv.schemaname <> 'information_schema'
| ORDER BY pgv.schemaname, pgv.viewname;
< result
< column "schemaname" "NAME" "string"
< column "viewname" "NAME" "string"
< column "comment" "VARCHAR" "string"
< row
: string
| app
: string
| long_view
: null
< row
: string
| app
: string
| recent
: null
< row
: string
| app
: string
| very_long_view
: null
> query
| SHOW VIEW "app"."long_view";
< result
< column "Show View DDL statement" "VARCHAR" "string"
< row
: string
| SELECT f.id, f.amount, f.notes, f.flag, f.seen, f.lookup_id, (f.id + 1) AS next_id, (f.amount * (2)::numeric) AS doubled FROM app.facts f WHERE ((f.seen > '2025-01-01 00:00:00'::timestamp without time zone) AND (f.flag = true));
> query
| SHOW VIEW "app"."recent";
< result
< column "Show View DDL statement" "VARCHAR" "string"
< row
: string
| SELECT facts.id, facts.amount FROM app.facts WHERE (facts.seen > '2026-01-01 00:00:00'::timestamp without time zone);
> query
| SHOW VIEW "app"."very_long_view";
< result
< column "Show View DDL statement" "VARCHAR" "string"
< row
: string
| SELECT f.id AS identifier_column_one, f.amount AS amount_column_two, f.notes AS notes_column_three, f.flag AS flag_column_four, f.seen AS seen_column_five, f.lookup_id AS lookup_column_six, (f.id + 1) AS next_identifier_seven, (f.amount * (2)::numeric) AS doubled_amount_eight, (f.id * 3) AS tripled_identifier_nine, ((f.notes)::text || 'x'::text) AS notes_with_suffix_ten FROM app.facts f WHERE ((((f.seen > '2025-01-01 00:00:00'::timestamp without time zone) AND (f.flag = true)) AND (f.amount > (0)::numeric(12,2))) AND (f.notes IS NOT NULL));
> commit
< ok
