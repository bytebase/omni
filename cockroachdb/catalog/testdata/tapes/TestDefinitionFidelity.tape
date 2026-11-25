-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestDefinitionFidelity default#1
> exec
| CREATE SCHEMA fid_cols;
< ok 0 -
> exec
| COMMENT ON SCHEMA fid_cols IS 'columns';
< ok 0 -
> exec
| CREATE TYPE fid_cols.mood AS ENUM ('happy', 'sad', 'it''s complicated');
< ok 0 -
> exec
| COMMENT ON TYPE fid_cols.mood IS 'how it feels';
< ok 0 -
> exec
| CREATE TYPE fid_cols.address AS (street STRING, city VARCHAR(50), zip INT8);
< ok 0 -
> exec
| CREATE TABLE fid_cols.customers (
|     id INT8 PRIMARY KEY DEFAULT unique_rowid(),
|     alt_id INT8 UNIQUE,
|     email STRING NOT NULL,
|     name VARCHAR(50) COLLATE "en_US" NOT NULL DEFAULT 'anon',
|     mood fid_cols.mood DEFAULT 'happy',
|     home fid_cols.address,
|     tags STRING[] DEFAULT ARRAY[]::STRING[],
|     balance DECIMAL(12, 2) DEFAULT 0 CHECK (balance >= 0),
|     ratio FLOAT8,
|     small INT2,
|     secret STRING NOT VISIBLE,
|     doc JSONB,
|     created TIMESTAMPTZ DEFAULT now(),
|     CONSTRAINT customers_email_check CHECK (email != '')
| );
< ok 0 -
> exec
| COMMENT ON TABLE fid_cols.customers IS 'a customer';
< ok 0 -
> exec
| COMMENT ON COLUMN fid_cols.customers.email IS 'login';
< ok 0 -
> exec
| COMMENT ON CONSTRAINT customers_email_check ON fid_cols.customers IS 'not empty';
< ok 0 -
> exec
| CREATE INDEX customers_name_idx ON fid_cols.customers (name DESC) STORING (balance);
< ok 0 -
> exec
| COMMENT ON INDEX fid_cols.customers@customers_name_idx IS 'names';
< ok 0 -
> exec
| CREATE INDEX customers_small_idx ON fid_cols.customers (small) WHERE small > 0;
< ok 0 -
> exec
| CREATE INDEX customers_lower_email_idx ON fid_cols.customers (lower(email));
< ok 0 -
> exec
| CREATE INVERTED INDEX customers_doc_idx ON fid_cols.customers (doc);
< ok 0 -
> exec
| CREATE INDEX customers_ratio_idx ON fid_cols.customers (ratio) USING HASH WITH (bucket_count = 8);
< ok 0 -
> exec
| CREATE INDEX customers_created_idx ON fid_cols.customers (created) NOT VISIBLE;
< ok 0 -
> exec
| CREATE TABLE fid_cols.order_lines (
|     order_id INT8 NOT NULL,
|     line INT4 NOT NULL,
|     customer_id INT8 NOT NULL REFERENCES fid_cols.customers (id) ON DELETE CASCADE ON UPDATE RESTRICT,
|     amount DECIMAL(10, 2) NOT NULL,
|     CONSTRAINT order_lines_pkey PRIMARY KEY (order_id ASC, line DESC),
|     CONSTRAINT order_lines_amount_positive CHECK (amount > 0)
| );
< ok 0 -
> exec
| CREATE TABLE fid_cols.events (at TIMESTAMPTZ, note STRING);
< ok 0 -
> exec
| CREATE TABLE fid_cols.sharded (id INT8 NOT NULL, small INT4, PRIMARY KEY (id) USING HASH WITH (bucket_count = 4));
< ok 0 -
> query
| SHOW CLUSTER SETTING version;
< result
< column "version" "TEXT"
< row
: string
| 25.2
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
| omni_testdefinitionfidelity
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
| fid_cols
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
| fid_cols
: string
| customers
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 1
: string
| unique_rowid()
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
| fid_cols
: string
| customers
: string
| alt_id
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
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
| fid_cols
: string
| customers
: string
| email
: string
| text
: null
: null
: null
: int64 3
: null
: string
| NO
: null
: string
| pg_catalog
: string
| text
: string
| login
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| name
: string
| character varying
: int64 50
: null
: null
: int64 4
: string
| 'anon' COLLATE en_US
: string
| NO
: string
| en_US
: string
| pg_catalog
: string
| varchar
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| mood
: string
| USER-DEFINED
: null
: null
: null
: int64 5
: string
| 'happy'
: string
| YES
: null
: string
| fid_cols
: string
| mood
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| home
: string
| USER-DEFINED
: null
: null
: null
: int64 6
: null
: string
| YES
: null
: string
| fid_cols
: string
| address
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| tags
: string
| ARRAY
: null
: null
: null
: int64 7
: string
| ARRAY[]
: string
| YES
: null
: string
| pg_catalog
: string
| _text
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| balance
: string
| numeric
: null
: int64 12
: int64 2
: int64 8
: string
| 0
: string
| YES
: null
: string
| pg_catalog
: string
| numeric
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| ratio
: string
| double precision
: null
: int64 53
: null
: int64 9
: null
: string
| YES
: null
: string
| pg_catalog
: string
| float8
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| small
: string
| smallint
: null
: int64 16
: int64 0
: int64 10
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int2
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| secret
: string
| text
: null
: null
: null
: int64 11
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| YES
< row
: string
| fid_cols
: string
| customers
: string
| doc
: string
| jsonb
: null
: null
: null
: int64 12
: null
: string
| YES
: null
: string
| pg_catalog
: string
| jsonb
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| created
: string
| timestamp with time zone
: null
: null
: null
: int64 13
: string
| now()
: string
| YES
: null
: string
| pg_catalog
: string
| timestamptz
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| crdb_internal_ratio_shard_8
: string
| bigint
: null
: int64 64
: int64 0
: int64 15
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
| YES
< row
: string
| fid_cols
: string
| events
: string
| at
: string
| timestamp with time zone
: null
: null
: null
: int64 1
: null
: string
| YES
: null
: string
| pg_catalog
: string
| timestamptz
: null
: string
| NO
< row
: string
| fid_cols
: string
| events
: string
| note
: string
| text
: null
: null
: null
: int64 2
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
< row
: string
| fid_cols
: string
| events
: string
| rowid
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| unique_rowid()
: string
| NO
: null
: string
| pg_catalog
: string
| int8
: null
: string
| YES
< row
: string
| fid_cols
: string
| order_lines
: string
| order_id
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
| fid_cols
: string
| order_lines
: string
| line
: string
| integer
: null
: int64 32
: int64 0
: int64 2
: null
: string
| NO
: null
: string
| pg_catalog
: string
| int4
: null
: string
| NO
< row
: string
| fid_cols
: string
| order_lines
: string
| customer_id
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
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
| fid_cols
: string
| order_lines
: string
| amount
: string
| numeric
: null
: int64 10
: int64 2
: int64 4
: null
: string
| NO
: null
: string
| pg_catalog
: string
| numeric
: null
: string
| NO
< row
: string
| fid_cols
: string
| sharded
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
| fid_cols
: string
| sharded
: string
| small
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
: string
| NO
< row
: string
| fid_cols
: string
| sharded
: string
| crdb_internal_id_shard_4
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
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
| YES
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
| fid_cols
: string
| customers
: string
| customers_alt_id_key
: string
| CREATE UNIQUE INDEX customers_alt_id_key ON omni_testdefinitionfidelity.fid_cols.customers USING btree (alt_id ASC)
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_created_idx
: string
| CREATE INDEX customers_created_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (created ASC) NOT VISIBLE
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_doc_idx
: string
| CREATE INDEX customers_doc_idx ON omni_testdefinitionfidelity.fid_cols.customers USING gin (doc)
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_lower_email_idx
: string
| CREATE INDEX customers_lower_email_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (lower(email) ASC)
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_name_idx
: string
| CREATE INDEX customers_name_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (name DESC) STORING (balance)
: null
: string
| names
< row
: string
| fid_cols
: string
| customers
: string
| customers_pkey
: string
| CREATE UNIQUE INDEX customers_pkey ON omni_testdefinitionfidelity.fid_cols.customers USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_ratio_idx
: string
| CREATE INDEX customers_ratio_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (ratio ASC) USING HASH WITH (bucket_count=8)
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_small_idx
: string
| CREATE INDEX customers_small_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (small ASC) WHERE (small > 0)
: null
: null
< row
: string
| fid_cols
: string
| events
: string
| events_pkey
: string
| CREATE UNIQUE INDEX events_pkey ON omni_testdefinitionfidelity.fid_cols.events USING btree (rowid ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_cols
: string
| order_lines
: string
| order_lines_pkey
: string
| CREATE UNIQUE INDEX order_lines_pkey ON omni_testdefinitionfidelity.fid_cols.order_lines USING btree (order_id ASC, line DESC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_cols
: string
| sharded
: string
| sharded_pkey
: string
| CREATE UNIQUE INDEX sharded_pkey ON omni_testdefinitionfidelity.fid_cols.sharded USING btree (id ASC) USING HASH WITH (bucket_count=4)
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
< row
: string
| fid_cols
: string
| customers
: string
| check_balance
: string
| CHECK ((balance >= 0))
: null
< row
: string
| fid_cols
: string
| customers
: string
| check_crdb_internal_ratio_shard_8
: string
| CHECK ((crdb_internal_ratio_shard_8 IN (0, 1, 2, 3, 4, 5, 6, 7)))
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_email_check
: string
| CHECK ((email != ''::STRING))
: string
| not empty
< row
: string
| fid_cols
: string
| order_lines
: string
| order_lines_amount_positive
: string
| CHECK ((amount > 0))
: null
< row
: string
| fid_cols
: string
| sharded
: string
| check_crdb_internal_id_shard_4
: string
| CHECK ((crdb_internal_id_shard_4 IN (0, 1, 2, 3)))
: null
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
< row
: string
| fid_cols
: string
| order_lines
: string
| order_lines_customer_id_fkey
: string
| fid_cols
: string
| customers
: string
| c
: string
| r
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES fid_cols.customers(id) ON DELETE CASCADE ON UPDATE RESTRICT
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
| fid_cols
: string
| customers
: int64 0
: string
| a customer
< row
: string
| fid_cols
: string
| events
: int64 0
: null
< row
: string
| fid_cols
: string
| order_lines
: int64 0
: null
< row
: string
| fid_cols
: string
| sharded
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
< row
: string
| fid_cols
: string
| mood
: string
| happy
< row
: string
| fid_cols
: string
| mood
: string
| sad
< row
: string
| fid_cols
: string
| mood
: string
| it's complicated
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
< row
: string
| fid_cols
: string
| address
: string
| street
: string
| text
< row
: string
| fid_cols
: string
| address
: string
| city
: string
| character varying(50)
< row
: string
| fid_cols
: string
| address
: string
| zip
: string
| bigint
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
> exec
| DROP SCHEMA IF EXISTS "fid_cols" CASCADE;
< ok 0 -
> exec
| SET statement_timeout = 0;
< ok 0 -
> exec
| SET lock_timeout = 0;
< ok 0 -
> exec
| SET idle_in_transaction_session_timeout = 0;
< ok 0 -
> exec
| SET client_encoding = 'UTF8';
< ok 0 -
> exec
| SET standard_conforming_strings = on;
< ok 0 -
> exec
| SELECT pg_catalog.set_config('search_path', '', false);
< ok 1 -
> exec
| SET check_function_bodies = false;
< ok 0 -
> exec
| SET xmloption = content;
< ok 0 -
> exec
| SET client_min_messages = warning;
< ok 0 -
> exec
| SET row_security = off;
< ok 0 -
> exec
| CREATE SCHEMA IF NOT EXISTS "fid_cols";
< ok 0 -
> exec
| CREATE TYPE "fid_cols"."address" AS (
|     "street" text,
|     "city" character varying(50),
|     "zip" bigint
| );
< ok 0 -
> exec
| CREATE TYPE "fid_cols"."mood" AS ENUM (
|     'happy',
|     'sad',
|     'it''s complicated'
| );
< ok 0 -
> exec
| SET default_tablespace = '';
< ok 0 -
> exec
| CREATE TABLE "fid_cols"."customers" (
|     "id" bigint DEFAULT unique_rowid() NOT NULL,
|     "alt_id" bigint,
|     "email" text NOT NULL,
|     "name" character varying(50) COLLATE "en_US" DEFAULT ('anon' COLLATE en_US) NOT NULL,
|     "mood" fid_cols.mood DEFAULT 'happy',
|     "home" fid_cols.address,
|     "tags" _text DEFAULT ARRAY[],
|     "balance" numeric(12,2) DEFAULT 0,
|     "ratio" double precision,
|     "small" smallint,
|     "secret" text NOT VISIBLE,
|     "doc" jsonb,
|     "created" timestamp with time zone DEFAULT now(),
|     CONSTRAINT "customers_pkey" PRIMARY KEY (id ASC),
|     CONSTRAINT "check_balance" CHECK ((balance >= 0)),
|     CONSTRAINT "customers_email_check" CHECK ((email != ''::STRING))
| );
< ok 0 -
> exec
| COMMENT ON TABLE "fid_cols"."customers" IS 'a customer';
< ok 0 -
> exec
| COMMENT ON COLUMN "fid_cols"."customers"."email" IS 'login';
< ok 0 -
> exec
| COMMENT ON CONSTRAINT "customers_email_check" ON "fid_cols"."customers" IS 'not empty';
< ok 0 -
> exec
| CREATE UNIQUE INDEX customers_alt_id_key ON fid_cols.customers (alt_id ASC);
< ok 0 -
> exec
| CREATE INDEX customers_created_idx ON fid_cols.customers (created ASC) NOT VISIBLE;
< ok 0 -
> exec
| CREATE INVERTED INDEX customers_doc_idx ON fid_cols.customers (doc);
< ok 0 -
> exec
| CREATE INDEX customers_lower_email_idx ON fid_cols.customers (lower(email) ASC);
< ok 0 -
> exec
| CREATE INDEX customers_name_idx ON fid_cols.customers (name DESC) STORING (balance);
< ok 0 -
> exec
| COMMENT ON INDEX "fid_cols"."customers"@"customers_name_idx" IS 'names';
< ok 0 -
> exec
| CREATE INDEX customers_ratio_idx ON fid_cols.customers (ratio ASC) USING HASH WITH ('bucket_count' = 8);
< ok 0 -
> exec
| CREATE INDEX customers_small_idx ON fid_cols.customers (small ASC) WHERE (small > 0);
< ok 0 -
> exec
| CREATE TABLE "fid_cols"."events" (
|     "at" timestamp with time zone,
|     "note" text
| );
< ok 0 -
> exec
| CREATE TABLE "fid_cols"."order_lines" (
|     "order_id" bigint NOT NULL,
|     "line" INT4 NOT NULL,
|     "customer_id" bigint NOT NULL,
|     "amount" numeric(10,2) NOT NULL,
|     CONSTRAINT "order_lines_pkey" PRIMARY KEY (order_id ASC, line DESC),
|     CONSTRAINT "order_lines_amount_positive" CHECK ((amount > 0))
| );
< ok 0 -
> exec
| CREATE TABLE "fid_cols"."sharded" (
|     "id" bigint NOT NULL,
|     "small" INT4,
|     CONSTRAINT "sharded_pkey" PRIMARY KEY (id ASC) USING HASH WITH ('bucket_count' = 4)
| );
< ok 0 -
> exec
| ALTER TABLE "fid_cols"."order_lines"
|     ADD CONSTRAINT "order_lines_customer_id_fkey" FOREIGN KEY ("customer_id")
|     REFERENCES "fid_cols"."customers" ("id")
|     ON DELETE CASCADE
|     ON UPDATE RESTRICT;
< ok 0 -
> query
| SHOW CLUSTER SETTING version;
< result
< column "version" "TEXT"
< row
: string
| 25.2
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
| omni_testdefinitionfidelity
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
|
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
| fid_cols
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
| fid_cols
: string
| customers
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 1
: string
| unique_rowid()
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
| fid_cols
: string
| customers
: string
| alt_id
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
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
| fid_cols
: string
| customers
: string
| email
: string
| text
: null
: null
: null
: int64 3
: null
: string
| NO
: null
: string
| pg_catalog
: string
| text
: string
| login
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| name
: string
| character varying
: int64 50
: null
: null
: int64 4
: string
| ('anon' COLLATE en_US)
: string
| NO
: string
| en_US
: string
| pg_catalog
: string
| varchar
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| mood
: string
| USER-DEFINED
: null
: null
: null
: int64 5
: string
| 'happy'
: string
| YES
: null
: string
| fid_cols
: string
| mood
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| home
: string
| USER-DEFINED
: null
: null
: null
: int64 6
: null
: string
| YES
: null
: string
| fid_cols
: string
| address
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| tags
: string
| ARRAY
: null
: null
: null
: int64 7
: string
| ARRAY[]
: string
| YES
: null
: string
| pg_catalog
: string
| _text
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| balance
: string
| numeric
: null
: int64 12
: int64 2
: int64 8
: string
| 0
: string
| YES
: null
: string
| pg_catalog
: string
| numeric
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| ratio
: string
| double precision
: null
: int64 53
: null
: int64 9
: null
: string
| YES
: null
: string
| pg_catalog
: string
| float8
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| small
: string
| smallint
: null
: int64 16
: int64 0
: int64 10
: null
: string
| YES
: null
: string
| pg_catalog
: string
| int2
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| secret
: string
| text
: null
: null
: null
: int64 11
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| YES
< row
: string
| fid_cols
: string
| customers
: string
| doc
: string
| jsonb
: null
: null
: null
: int64 12
: null
: string
| YES
: null
: string
| pg_catalog
: string
| jsonb
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| created
: string
| timestamp with time zone
: null
: null
: null
: int64 13
: string
| now()
: string
| YES
: null
: string
| pg_catalog
: string
| timestamptz
: null
: string
| NO
< row
: string
| fid_cols
: string
| customers
: string
| crdb_internal_ratio_shard_8
: string
| bigint
: null
: int64 64
: int64 0
: int64 15
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
| YES
< row
: string
| fid_cols
: string
| events
: string
| at
: string
| timestamp with time zone
: null
: null
: null
: int64 1
: null
: string
| YES
: null
: string
| pg_catalog
: string
| timestamptz
: null
: string
| NO
< row
: string
| fid_cols
: string
| events
: string
| note
: string
| text
: null
: null
: null
: int64 2
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
< row
: string
| fid_cols
: string
| events
: string
| rowid
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| unique_rowid()
: string
| NO
: null
: string
| pg_catalog
: string
| int8
: null
: string
| YES
< row
: string
| fid_cols
: string
| order_lines
: string
| order_id
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
| fid_cols
: string
| order_lines
: string
| line
: string
| integer
: null
: int64 32
: int64 0
: int64 2
: null
: string
| NO
: null
: string
| pg_catalog
: string
| int4
: null
: string
| NO
< row
: string
| fid_cols
: string
| order_lines
: string
| customer_id
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
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
| fid_cols
: string
| order_lines
: string
| amount
: string
| numeric
: null
: int64 10
: int64 2
: int64 4
: null
: string
| NO
: null
: string
| pg_catalog
: string
| numeric
: null
: string
| NO
< row
: string
| fid_cols
: string
| sharded
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
| fid_cols
: string
| sharded
: string
| small
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
: string
| NO
< row
: string
| fid_cols
: string
| sharded
: string
| crdb_internal_id_shard_4
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
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
| YES
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
| fid_cols
: string
| customers
: string
| customers_alt_id_key
: string
| CREATE UNIQUE INDEX customers_alt_id_key ON omni_testdefinitionfidelity.fid_cols.customers USING btree (alt_id ASC)
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_created_idx
: string
| CREATE INDEX customers_created_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (created ASC) NOT VISIBLE
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_doc_idx
: string
| CREATE INDEX customers_doc_idx ON omni_testdefinitionfidelity.fid_cols.customers USING gin (doc)
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_lower_email_idx
: string
| CREATE INDEX customers_lower_email_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (lower(email) ASC)
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_name_idx
: string
| CREATE INDEX customers_name_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (name DESC) STORING (balance)
: null
: string
| names
< row
: string
| fid_cols
: string
| customers
: string
| customers_pkey
: string
| CREATE UNIQUE INDEX customers_pkey ON omni_testdefinitionfidelity.fid_cols.customers USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_ratio_idx
: string
| CREATE INDEX customers_ratio_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (ratio ASC) USING HASH WITH (bucket_count=8)
: null
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_small_idx
: string
| CREATE INDEX customers_small_idx ON omni_testdefinitionfidelity.fid_cols.customers USING btree (small ASC) WHERE (small > 0)
: null
: null
< row
: string
| fid_cols
: string
| events
: string
| events_pkey
: string
| CREATE UNIQUE INDEX events_pkey ON omni_testdefinitionfidelity.fid_cols.events USING btree (rowid ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_cols
: string
| order_lines
: string
| order_lines_pkey
: string
| CREATE UNIQUE INDEX order_lines_pkey ON omni_testdefinitionfidelity.fid_cols.order_lines USING btree (order_id ASC, line DESC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_cols
: string
| sharded
: string
| sharded_pkey
: string
| CREATE UNIQUE INDEX sharded_pkey ON omni_testdefinitionfidelity.fid_cols.sharded USING btree (id ASC) USING HASH WITH (bucket_count=4)
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
< row
: string
| fid_cols
: string
| customers
: string
| check_balance
: string
| CHECK ((balance >= 0))
: null
< row
: string
| fid_cols
: string
| customers
: string
| check_crdb_internal_ratio_shard_8
: string
| CHECK ((crdb_internal_ratio_shard_8 IN (0, 1, 2, 3, 4, 5, 6, 7)))
: null
< row
: string
| fid_cols
: string
| customers
: string
| customers_email_check
: string
| CHECK ((email != ''::STRING))
: string
| not empty
< row
: string
| fid_cols
: string
| order_lines
: string
| order_lines_amount_positive
: string
| CHECK ((amount > 0))
: null
< row
: string
| fid_cols
: string
| sharded
: string
| check_crdb_internal_id_shard_4
: string
| CHECK ((crdb_internal_id_shard_4 IN (0, 1, 2, 3)))
: null
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
< row
: string
| fid_cols
: string
| order_lines
: string
| order_lines_customer_id_fkey
: string
| fid_cols
: string
| customers
: string
| c
: string
| r
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES fid_cols.customers(id) ON DELETE CASCADE ON UPDATE RESTRICT
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
| fid_cols
: string
| customers
: int64 0
: string
| a customer
< row
: string
| fid_cols
: string
| events
: int64 0
: null
< row
: string
| fid_cols
: string
| order_lines
: int64 0
: null
< row
: string
| fid_cols
: string
| sharded
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
< row
: string
| fid_cols
: string
| mood
: string
| happy
< row
: string
| fid_cols
: string
| mood
: string
| sad
< row
: string
| fid_cols
: string
| mood
: string
| it's complicated
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
< row
: string
| fid_cols
: string
| address
: string
| street
: string
| text
< row
: string
| fid_cols
: string
| address
: string
| city
: string
| character varying(50)
< row
: string
| fid_cols
: string
| address
: string
| zip
: string
| bigint
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
> exec
| DROP SCHEMA IF EXISTS "fid_cols" CASCADE;
< ok 0 -
> exec
| CREATE SCHEMA fid_routines;
< ok 0 -
> exec
| CREATE SEQUENCE fid_routines.counter AS INT8 INCREMENT BY 5 MINVALUE -100 MAXVALUE 1000 START WITH 10 CACHE 3;
< ok 0 -
> exec
| CREATE SEQUENCE fid_routines.plain;
< ok 0 -
> exec
| CREATE TABLE fid_routines.items (
|     id INT8 PRIMARY KEY DEFAULT nextval('fid_routines.plain'),
|     name STRING NOT NULL,
|     price DECIMAL(10, 2)
| );
< ok 0 -
> exec
| CREATE VIEW fid_routines.v_items AS SELECT id, name FROM fid_routines.items WHERE price > 0;
< ok 0 -
> exec
| CREATE VIEW fid_routines.v_item_names AS SELECT name FROM fid_routines.v_items;
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW fid_routines.item_totals AS SELECT name, sum(price) AS total FROM fid_routines.items GROUP BY name;
< ok 0 -
> exec
| CREATE INDEX item_totals_name_idx ON fid_routines.item_totals (name);
< ok 0 -
> exec
| CREATE FUNCTION fid_routines.add(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT a + b $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_routines.item_count() RETURNS INT8 LANGUAGE SQL STABLE AS $$ SELECT count(*) FROM fid_routines.v_items $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_routines.touch() RETURNS INT8 LANGUAGE PLPGSQL AS $$
| BEGIN
|     RETURN (SELECT count(*) FROM fid_routines.items);
| END
| $$;
< ok 0 -
> exec
| CREATE PROCEDURE fid_routines.reset_prices() LANGUAGE SQL AS $$ UPDATE fid_routines.items SET price = 0 $$;
< ok 0 -
> exec
| CREATE TABLE fid_routines.scored (id INT8 PRIMARY KEY, score INT8 DEFAULT fid_routines.add(1, 2));
< ok 0 -
> query
| SHOW CLUSTER SETTING version;
< result
< column "version" "TEXT"
< row
: string
| 25.2
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
| omni_testdefinitionfidelity
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
|
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
| fid_routines
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
| fid_routines
: string
| item_totals
: string
| name
: string
| text
: null
: null
: null
: int64 1
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
< row
: string
| fid_routines
: string
| item_totals
: string
| total
: string
| numeric
: null
: null
: null
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
: string
| NO
< row
: string
| fid_routines
: string
| item_totals
: string
| rowid
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| unique_rowid()
: string
| NO
: null
: string
| pg_catalog
: string
| int8
: null
: string
| YES
< row
: string
| fid_routines
: string
| items
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 1
: string
| nextval('fid_routines.plain'::REGCLASS)
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
| fid_routines
: string
| items
: string
| name
: string
| text
: null
: null
: null
: int64 2
: null
: string
| NO
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
< row
: string
| fid_routines
: string
| items
: string
| price
: string
| numeric
: null
: int64 10
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
: string
| NO
< row
: string
| fid_routines
: string
| scored
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
| fid_routines
: string
| scored
: string
| score
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_routines.add(1, 2)
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
| fid_routines
: string
| v_item_names
: string
| name
: string
| text
: null
: null
: null
: int64 1
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
< row
: string
| fid_routines
: string
| v_items
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
| fid_routines
: string
| v_items
: string
| name
: string
| text
: null
: null
: null
: int64 2
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
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
| fid_routines
: string
| counter
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testdefinitionfidelity.fid_routines.counter USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_routines
: string
| item_totals
: string
| item_totals_name_idx
: string
| CREATE INDEX item_totals_name_idx ON omni_testdefinitionfidelity.fid_routines.item_totals USING btree (name ASC)
: null
: null
< row
: string
| fid_routines
: string
| item_totals
: string
| item_totals_pkey
: string
| CREATE UNIQUE INDEX item_totals_pkey ON omni_testdefinitionfidelity.fid_routines.item_totals USING btree (rowid ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_routines
: string
| items
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON omni_testdefinitionfidelity.fid_routines.items USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_routines
: string
| plain
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testdefinitionfidelity.fid_routines.plain USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_routines
: string
| scored
: string
| scored_pkey
: string
| CREATE UNIQUE INDEX scored_pkey ON omni_testdefinitionfidelity.fid_routines.scored USING btree (id ASC)
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
| fid_routines
: string
| items
: int64 0
: null
< row
: string
| fid_routines
: string
| scored
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
| fid_routines
: string
| v_item_names
: string
| SELECT name FROM omni_testdefinitionfidelity.fid_routines.v_items
: null
< row
: string
| fid_routines
: string
| v_items
: string
| SELECT id, name FROM omni_testdefinitionfidelity.fid_routines.items WHERE price > 0
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
| fid_routines
: string
| v_item_names
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_routines
: string
| v_items
: string
| id
< row
: string
| fid_routines
: string
| v_items
: string
| name
> query ^
: string
| fid_routines
: string
| v_items
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_routines
: string
| items
: string
| id
< row
: string
| fid_routines
: string
| items
: string
| name
< row
: string
| fid_routines
: string
| items
: string
| price
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
< row
: string
| fid_routines
: string
| item_totals
: string
| SELECT name, sum(price) AS total FROM omni_testdefinitionfidelity.fid_routines.items GROUP BY name
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
| fid_routines
: string
| item_totals
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_routines
: string
| items
: string
| name
< row
: string
| fid_routines
: string
| items
: string
| price
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
< row
: string
| fid_routines
: string
| add
: string
| add(int8, int8)
: string
| CREATE FUNCTION fid_routines.add(a INT8, b INT8)
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT a + b;
| $$
< row
: string
| fid_routines
: string
| item_count
: string
| item_count()
: string
| CREATE FUNCTION fid_routines.item_count()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testdefinitionfidelity.fid_routines.v_items;
| $$
< row
: string
| fid_routines
: string
| reset_prices
: string
| reset_prices()
: string
| CREATE PROCEDURE fid_routines.reset_prices()
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	UPDATE omni_testdefinitionfidelity.fid_routines.items SET price = 0;
| $$
< row
: string
| fid_routines
: string
| touch
: string
| touch()
: string
| CREATE FUNCTION fid_routines.touch()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE plpgsql
| 	SECURITY INVOKER
| 	AS $$
| 	BEGIN
| 	RETURN (SELECT count(*) FROM omni_testdefinitionfidelity.fid_routines.items);
| 	END;
| $$
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
< row
: string
| fid_routines
: string
| counter
: string
| INT8
: string
| 10
: string
| -100
: string
| 1000
: string
| 5
: string
| NO
< row
: string
| fid_routines
: string
| plain
: string
| INT8
: string
| 1
: string
| 1
: string
| 9223372036854775807
: string
| 1
: string
| NO
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
> exec
| DROP SCHEMA IF EXISTS "fid_routines" CASCADE;
< ok 0 -
> exec
| SET statement_timeout = 0;
< ok 0 -
> exec
| SET lock_timeout = 0;
< ok 0 -
> exec
| SET idle_in_transaction_session_timeout = 0;
< ok 0 -
> exec
| SET client_encoding = 'UTF8';
< ok 0 -
> exec
| SET standard_conforming_strings = on;
< ok 0 -
> exec
| SELECT pg_catalog.set_config('search_path', '', false);
< ok 1 -
> exec
| SET check_function_bodies = false;
< ok 0 -
> exec
| SET xmloption = content;
< ok 0 -
> exec
| SET client_min_messages = warning;
< ok 0 -
> exec
| SET row_security = off;
< ok 0 -
> exec
| CREATE SCHEMA IF NOT EXISTS "fid_routines";
< ok 0 -
> exec
| CREATE SEQUENCE "fid_routines"."counter"
|     AS INT8
| 	START WITH 10
| 	INCREMENT BY 5
| 	MINVALUE -100
| 	MAXVALUE 1000
| 	NO CYCLE;
< ok 0 -
> exec
| CREATE SEQUENCE "fid_routines"."plain"
|     AS INT8
| 	START WITH 1
| 	INCREMENT BY 1
| 	MINVALUE 1
| 	MAXVALUE 9223372036854775807
| 	NO CYCLE;
< ok 0 -
> exec
| SET default_tablespace = '';
< ok 0 -
> exec
| CREATE TABLE "fid_routines"."items" (
|     "id" bigint DEFAULT nextval('fid_routines.plain'::REGCLASS) NOT NULL,
|     "name" text NOT NULL,
|     "price" numeric(10,2),
|     CONSTRAINT "items_pkey" PRIMARY KEY (id ASC)
| );
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW "fid_routines"."item_totals" AS 
| SELECT name, sum(price) AS total FROM fid_routines.items GROUP BY name;
< ok 0 -
> exec
| CREATE INDEX item_totals_name_idx ON fid_routines.item_totals (name ASC);
< ok 0 -
> exec
| CREATE VIEW "fid_routines"."v_items" AS 
| SELECT id, name FROM fid_routines.items WHERE price > 0;
< ok 0 -
> exec
| CREATE VIEW "fid_routines"."v_item_names" AS 
| SELECT name FROM fid_routines.v_items;
< ok 0 -
> exec
| CREATE FUNCTION fid_routines.add(a INT8, b INT8)
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT a + b;
| $$;
< ok 0 -
> exec
| CREATE TABLE "fid_routines"."scored" (
|     "id" bigint NOT NULL,
|     "score" bigint DEFAULT fid_routines.add(1, 2),
|     CONSTRAINT "scored_pkey" PRIMARY KEY (id ASC)
| );
< ok 0 -
> exec
| CREATE FUNCTION fid_routines.item_count()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| SELECT count(*) FROM fid_routines.v_items;
| $$;
< ok 0 -
> exec
| CREATE PROCEDURE fid_routines.reset_prices()
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| UPDATE fid_routines.items SET price = 0;
| $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_routines.touch()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE plpgsql
| 	SECURITY INVOKER
| 	AS $$
| BEGIN
| RETURN (SELECT count(*) FROM fid_routines.items);
| END;
|
| $$;
< ok 0 -
> query
| SHOW CLUSTER SETTING version;
< result
< column "version" "TEXT"
< row
: string
| 25.2
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
| omni_testdefinitionfidelity
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
|
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
| fid_routines
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
| fid_routines
: string
| item_totals
: string
| name
: string
| text
: null
: null
: null
: int64 1
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
< row
: string
| fid_routines
: string
| item_totals
: string
| total
: string
| numeric
: null
: null
: null
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
: string
| NO
< row
: string
| fid_routines
: string
| item_totals
: string
| rowid
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| unique_rowid()
: string
| NO
: null
: string
| pg_catalog
: string
| int8
: null
: string
| YES
< row
: string
| fid_routines
: string
| items
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 1
: string
| nextval('fid_routines.plain'::REGCLASS)
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
| fid_routines
: string
| items
: string
| name
: string
| text
: null
: null
: null
: int64 2
: null
: string
| NO
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
< row
: string
| fid_routines
: string
| items
: string
| price
: string
| numeric
: null
: int64 10
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
: string
| NO
< row
: string
| fid_routines
: string
| scored
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
| fid_routines
: string
| scored
: string
| score
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_routines.add(1, 2)
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
| fid_routines
: string
| v_item_names
: string
| name
: string
| text
: null
: null
: null
: int64 1
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
< row
: string
| fid_routines
: string
| v_items
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
| fid_routines
: string
| v_items
: string
| name
: string
| text
: null
: null
: null
: int64 2
: null
: string
| YES
: null
: string
| pg_catalog
: string
| text
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
| fid_routines
: string
| counter
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testdefinitionfidelity.fid_routines.counter USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_routines
: string
| item_totals
: string
| item_totals_name_idx
: string
| CREATE INDEX item_totals_name_idx ON omni_testdefinitionfidelity.fid_routines.item_totals USING btree (name ASC)
: null
: null
< row
: string
| fid_routines
: string
| item_totals
: string
| item_totals_pkey
: string
| CREATE UNIQUE INDEX item_totals_pkey ON omni_testdefinitionfidelity.fid_routines.item_totals USING btree (rowid ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_routines
: string
| items
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON omni_testdefinitionfidelity.fid_routines.items USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_routines
: string
| plain
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testdefinitionfidelity.fid_routines.plain USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_routines
: string
| scored
: string
| scored_pkey
: string
| CREATE UNIQUE INDEX scored_pkey ON omni_testdefinitionfidelity.fid_routines.scored USING btree (id ASC)
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
| fid_routines
: string
| items
: int64 0
: null
< row
: string
| fid_routines
: string
| scored
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
| fid_routines
: string
| v_item_names
: string
| SELECT name FROM omni_testdefinitionfidelity.fid_routines.v_items
: null
< row
: string
| fid_routines
: string
| v_items
: string
| SELECT id, name FROM omni_testdefinitionfidelity.fid_routines.items WHERE price > 0
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
| fid_routines
: string
| v_item_names
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_routines
: string
| v_items
: string
| id
< row
: string
| fid_routines
: string
| v_items
: string
| name
> query ^
: string
| fid_routines
: string
| v_items
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_routines
: string
| items
: string
| id
< row
: string
| fid_routines
: string
| items
: string
| name
< row
: string
| fid_routines
: string
| items
: string
| price
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
< row
: string
| fid_routines
: string
| item_totals
: string
| SELECT name, sum(price) AS total FROM omni_testdefinitionfidelity.fid_routines.items GROUP BY name
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
| fid_routines
: string
| item_totals
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_routines
: string
| items
: string
| name
< row
: string
| fid_routines
: string
| items
: string
| price
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
< row
: string
| fid_routines
: string
| add
: string
| add(int8, int8)
: string
| CREATE FUNCTION fid_routines.add(a INT8, b INT8)
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT a + b;
| $$
< row
: string
| fid_routines
: string
| item_count
: string
| item_count()
: string
| CREATE FUNCTION fid_routines.item_count()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testdefinitionfidelity.fid_routines.v_items;
| $$
< row
: string
| fid_routines
: string
| reset_prices
: string
| reset_prices()
: string
| CREATE PROCEDURE fid_routines.reset_prices()
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	UPDATE omni_testdefinitionfidelity.fid_routines.items SET price = 0;
| $$
< row
: string
| fid_routines
: string
| touch
: string
| touch()
: string
| CREATE FUNCTION fid_routines.touch()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE plpgsql
| 	SECURITY INVOKER
| 	AS $$
| 	BEGIN
| 	RETURN (SELECT count(*) FROM omni_testdefinitionfidelity.fid_routines.items);
| 	END;
| $$
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
< row
: string
| fid_routines
: string
| counter
: string
| INT8
: string
| 10
: string
| -100
: string
| 1000
: string
| 5
: string
| NO
< row
: string
| fid_routines
: string
| plain
: string
| INT8
: string
| 1
: string
| 1
: string
| 9223372036854775807
: string
| 1
: string
| NO
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
> exec
| DROP SCHEMA IF EXISTS "fid_routines" CASCADE;
< ok 0 -
> exec
| CREATE SCHEMA fid_a;
< ok 0 -
> exec
| CREATE SCHEMA "Fid Mixed";
< ok 0 -
> exec
| CREATE TABLE fid_a.parents (id INT8 PRIMARY KEY);
< ok 0 -
> exec
| CREATE TABLE "Fid Mixed"."Children" (
|     id INT8 PRIMARY KEY,
|     "Parent Id" INT8 REFERENCES fid_a.parents (id)
| );
< ok 0 -
> exec
| CREATE VIEW "Fid Mixed".family AS SELECT p.id, c.id AS child FROM fid_a.parents AS p JOIN "Fid Mixed"."Children" AS c ON c."Parent Id" = p.id;
< ok 0 -
> query
| SHOW CLUSTER SETTING version;
< result
< column "version" "TEXT"
< row
: string
| 25.2
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
| omni_testdefinitionfidelity
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
|
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
| Fid Mixed
< row
: string
| fid_a
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
| Fid Mixed
: string
| Children
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
| Fid Mixed
: string
| Children
: string
| Parent Id
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
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
| Fid Mixed
: string
| family
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
| Fid Mixed
: string
| family
: string
| child
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
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
| fid_a
: string
| parents
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
| Fid Mixed
: string
| Children
: string
| Children_pkey
: string
| CREATE UNIQUE INDEX "Children_pkey" ON omni_testdefinitionfidelity."Fid Mixed"."Children" USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_a
: string
| parents
: string
| parents_pkey
: string
| CREATE UNIQUE INDEX parents_pkey ON omni_testdefinitionfidelity.fid_a.parents USING btree (id ASC)
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
< row
: string
| Fid Mixed
: string
| Children
: string
| Children_Parent Id_fkey
: string
| fid_a
: string
| parents
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY ("Parent Id") REFERENCES fid_a.parents(id)
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
| Fid Mixed
: string
| Children
: int64 0
: null
< row
: string
| fid_a
: string
| parents
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
| Fid Mixed
: string
| family
: string
| SELECT p.id, c.id AS child FROM omni_testdefinitionfidelity.fid_a.parents AS p JOIN omni_testdefinitionfidelity."Fid Mixed"."Children" AS c ON c."Parent Id" = p.id
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
| Fid Mixed
: string
| family
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| Fid Mixed
: string
| Children
: string
| Parent Id
< row
: string
| Fid Mixed
: string
| Children
: string
| id
< row
: string
| fid_a
: string
| parents
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
> exec
| DROP SCHEMA IF EXISTS "Fid Mixed" CASCADE;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "fid_a" CASCADE;
< ok 0 -
> exec
| SET statement_timeout = 0;
< ok 0 -
> exec
| SET lock_timeout = 0;
< ok 0 -
> exec
| SET idle_in_transaction_session_timeout = 0;
< ok 0 -
> exec
| SET client_encoding = 'UTF8';
< ok 0 -
> exec
| SET standard_conforming_strings = on;
< ok 0 -
> exec
| SELECT pg_catalog.set_config('search_path', '', false);
< ok 1 -
> exec
| SET check_function_bodies = false;
< ok 0 -
> exec
| SET xmloption = content;
< ok 0 -
> exec
| SET client_min_messages = warning;
< ok 0 -
> exec
| SET row_security = off;
< ok 0 -
> exec
| CREATE SCHEMA IF NOT EXISTS "Fid Mixed";
< ok 0 -
> exec
| CREATE SCHEMA IF NOT EXISTS "fid_a";
< ok 0 -
> exec
| SET default_tablespace = '';
< ok 0 -
> exec
| CREATE TABLE "Fid Mixed"."Children" (
|     "id" bigint NOT NULL,
|     "Parent Id" bigint,
|     CONSTRAINT "Children_pkey" PRIMARY KEY (id ASC)
| );
< ok 0 -
> exec
| CREATE TABLE "fid_a"."parents" (
|     "id" bigint NOT NULL,
|     CONSTRAINT "parents_pkey" PRIMARY KEY (id ASC)
| );
< ok 0 -
> exec
| CREATE VIEW "Fid Mixed"."family" AS 
| SELECT p.id, c.id AS child FROM fid_a.parents AS p JOIN "Fid Mixed"."Children" AS c ON c."Parent Id" = p.id;
< ok 0 -
> exec
| ALTER TABLE "Fid Mixed"."Children"
|     ADD CONSTRAINT "Children_Parent Id_fkey" FOREIGN KEY ("Parent Id")
|     REFERENCES "fid_a"."parents" ("id");
< ok 0 -
> query
| SHOW CLUSTER SETTING version;
< result
< column "version" "TEXT"
< row
: string
| 25.2
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
| omni_testdefinitionfidelity
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
|
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
| Fid Mixed
< row
: string
| fid_a
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
| Fid Mixed
: string
| Children
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
| Fid Mixed
: string
| Children
: string
| Parent Id
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
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
| Fid Mixed
: string
| family
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
| Fid Mixed
: string
| family
: string
| child
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
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
| fid_a
: string
| parents
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
| Fid Mixed
: string
| Children
: string
| Children_pkey
: string
| CREATE UNIQUE INDEX "Children_pkey" ON omni_testdefinitionfidelity."Fid Mixed"."Children" USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_a
: string
| parents
: string
| parents_pkey
: string
| CREATE UNIQUE INDEX parents_pkey ON omni_testdefinitionfidelity.fid_a.parents USING btree (id ASC)
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
< row
: string
| Fid Mixed
: string
| Children
: string
| Children_Parent Id_fkey
: string
| fid_a
: string
| parents
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY ("Parent Id") REFERENCES fid_a.parents(id)
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
| Fid Mixed
: string
| Children
: int64 0
: null
< row
: string
| fid_a
: string
| parents
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
| Fid Mixed
: string
| family
: string
| SELECT p.id, c.id AS child FROM omni_testdefinitionfidelity.fid_a.parents AS p JOIN omni_testdefinitionfidelity."Fid Mixed"."Children" AS c ON c."Parent Id" = p.id
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
| Fid Mixed
: string
| family
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| Fid Mixed
: string
| Children
: string
| Parent Id
< row
: string
| Fid Mixed
: string
| Children
: string
| id
< row
: string
| fid_a
: string
| parents
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
> exec
| DROP SCHEMA IF EXISTS "Fid Mixed" CASCADE;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "fid_a" CASCADE;
< ok 0 -
