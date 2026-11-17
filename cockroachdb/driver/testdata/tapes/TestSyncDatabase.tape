-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestSyncDatabase default#1
> exec
| -- The schema TestSyncDatabase syncs and the parity harness compares with
| -- Bytebase's driver. It covers each object kind CockroachDB's catalogs report,
| -- in the forms that change how the sync reads them.
|
| CREATE SCHEMA app;
| COMMENT ON SCHEMA app IS 'the application';
| CREATE SCHEMA "Mixed Case";
| CREATE SCHEMA audit;
|
| CREATE TYPE app.mood AS ENUM ('happy', 'sad', 'it''s complicated');
| COMMENT ON TYPE app.mood IS 'how it feels';
| CREATE TYPE app.address AS (street STRING, city VARCHAR(50), zip INT8);
|
| CREATE SEQUENCE app.counter AS INT8 INCREMENT BY 5 MINVALUE -100 MAXVALUE 1000 START WITH 10 CACHE 3;
| CREATE SEQUENCE audit.big_seq AS INT8;
|
| CREATE FUNCTION app.add(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT a + b $$;
|
| CREATE TABLE app.customer (
|     id INT8 PRIMARY KEY DEFAULT unique_rowid(),
|     alt_id INT8 UNIQUE,
|     email STRING NOT NULL,
|     name VARCHAR(50) COLLATE "en_US" NOT NULL DEFAULT 'anon',
|     mood app.mood DEFAULT 'happy',
|     home app.address,
|     tags STRING[] DEFAULT ARRAY[]::STRING[],
|     balance DECIMAL(12, 2) DEFAULT 0 CHECK (balance >= 0),
|     ratio FLOAT8,
|     small INT2,
|     created TIMESTAMPTZ DEFAULT now(),
|     CONSTRAINT customer_email_check CHECK (email != '')
| );
| COMMENT ON TABLE app.customer IS 'a customer';
| COMMENT ON COLUMN app.customer.email IS 'login';
| CREATE INDEX customer_name_idx ON app.customer (name DESC) STORING (balance);
| COMMENT ON INDEX app.customer@customer_name_idx IS 'names';
|
| CREATE TABLE app.orders (
|     id INT8 PRIMARY KEY,
|     customer_id INT8 NOT NULL REFERENCES app.customer(id) ON DELETE CASCADE ON UPDATE RESTRICT,
|     amount DECIMAL(10, 2) NOT NULL,
|     placed_at TIMESTAMP,
|     CONSTRAINT orders_amount_positive CHECK (amount > 0)
| );
| COMMENT ON TABLE app.orders IS 'an order';
| CREATE INDEX orders_customer_idx ON app.orders (customer_id);
|
| CREATE VIEW app.customer_orders AS
|     SELECT c.id, c.name, o.id AS order_id, o.amount
|     FROM app.customer AS c JOIN app.orders AS o ON o.customer_id = c.id;
|
| CREATE MATERIALIZED VIEW app.order_totals AS
|     SELECT customer_id, sum(amount) AS total FROM app.orders GROUP BY customer_id;
|
| INSERT INTO app.customer (id, alt_id, email, name, mood, balance, ratio, small, created) VALUES
|     (1, 11, 'one@example.com', 'one', 'happy', 10.50, 1.5, 1, '2024-01-01 00:00:00+00'),
|     (2, 12, 'two@example.com', 'two', 'sad', 20.25, 2.5, 2, '2024-01-02 00:00:00+00');
| INSERT INTO app.orders (id, customer_id, amount, placed_at) VALUES
|     (1, 1, 10.50, '2024-02-01 12:00:00'),
|     (2, 1, 4.50, '2024-02-02 12:00:00'),
|     (3, 2, 20.25, '2024-02-03 12:00:00');
|
| ANALYZE app.customer;
| ANALYZE app.orders;
|
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
| omni_testsyncdatabase
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
| Mixed Case
< row
: string
| app
< row
: string
| audit
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
| 	cols.ordinal_position,
| 	cols.column_default,
| 	cols.is_nullable,
| 	cols.collation_name,
| 	cols.udt_schema,
| 	cols.udt_name,
| 	pg_catalog.col_description(c.oid, cols.ordinal_position::int) as column_comment
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
< column "ordinal_position" "INT8"
< column "column_default" "TEXT"
< column "is_nullable" "TEXT"
< column "collation_name" "TEXT"
< column "udt_schema" "TEXT"
< column "udt_name" "TEXT"
< column "column_comment" "TEXT"
< row
: string
| app
: string
| customer
: string
| id
: string
| bigint
: null
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
< row
: string
| app
: string
| customer
: string
| alt_id
: string
| bigint
: null
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
< row
: string
| app
: string
| customer
: string
| email
: string
| text
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
< row
: string
| app
: string
| customer
: string
| name
: string
| character varying
: int64 50
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
< row
: string
| app
: string
| customer
: string
| mood
: string
| USER-DEFINED
: null
: int64 5
: string
| 'happy'
: string
| YES
: null
: string
| app
: string
| mood
: null
< row
: string
| app
: string
| customer
: string
| home
: string
| USER-DEFINED
: null
: int64 6
: null
: string
| YES
: null
: string
| app
: string
| address
: null
< row
: string
| app
: string
| customer
: string
| tags
: string
| ARRAY
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
< row
: string
| app
: string
| customer
: string
| balance
: string
| numeric
: null
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
< row
: string
| app
: string
| customer
: string
| ratio
: string
| double precision
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
< row
: string
| app
: string
| customer
: string
| small
: string
| smallint
: null
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
< row
: string
| app
: string
| customer
: string
| created
: string
| timestamp with time zone
: null
: int64 11
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
< row
: string
| app
: string
| customer_orders
: string
| id
: string
| bigint
: null
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
| customer_orders
: string
| name
: string
| character varying
: int64 50
: int64 2
: null
: string
| YES
: string
| en_US
: string
| pg_catalog
: string
| varchar
: null
< row
: string
| app
: string
| customer_orders
: string
| order_id
: string
| bigint
: null
: int64 3
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
| customer_orders
: string
| amount
: string
| numeric
: null
: int64 4
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
| order_totals
: string
| customer_id
: string
| bigint
: null
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
| order_totals
: string
| total
: string
| numeric
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
< row
: string
| app
: string
| order_totals
: string
| rowid
: string
| bigint
: null
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
< row
: string
| app
: string
| orders
: string
| id
: string
| bigint
: null
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
| orders
: string
| customer_id
: string
| bigint
: null
: int64 2
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
| orders
: string
| amount
: string
| numeric
: null
: int64 3
: null
: string
| NO
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
| orders
: string
| placed_at
: string
| timestamp without time zone
: null
: int64 4
: null
: string
| YES
: null
: string
| pg_catalog
: string
| timestamp
: null
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
| app
: string
| counter
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testsyncdatabase.app.counter USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| app
: string
| customer
: string
| customer_alt_id_key
: string
| CREATE UNIQUE INDEX customer_alt_id_key ON omni_testsyncdatabase.app.customer USING btree (alt_id ASC)
: null
: null
< row
: string
| app
: string
| customer
: string
| customer_name_idx
: string
| CREATE INDEX customer_name_idx ON omni_testsyncdatabase.app.customer USING btree (name DESC) STORING (balance)
: null
: string
| names
< row
: string
| app
: string
| customer
: string
| customer_pkey
: string
| CREATE UNIQUE INDEX customer_pkey ON omni_testsyncdatabase.app.customer USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| app
: string
| order_totals
: string
| order_totals_pkey
: string
| CREATE UNIQUE INDEX order_totals_pkey ON omni_testsyncdatabase.app.order_totals USING btree (rowid ASC)
: string
| PRIMARY KEY
: null
< row
: string
| app
: string
| orders
: string
| orders_customer_idx
: string
| CREATE INDEX orders_customer_idx ON omni_testsyncdatabase.app.orders USING btree (customer_id ASC)
: null
: null
< row
: string
| app
: string
| orders
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON omni_testsyncdatabase.app.orders USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| audit
: string
| big_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testsyncdatabase.audit.big_seq USING btree (value ASC)
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
| app
: string
| customer
: string
| check_balance
: string
| CHECK ((balance >= 0))
: null
< row
: string
| app
: string
| customer
: string
| customer_email_check
: string
| CHECK ((email != ''::STRING))
: null
< row
: string
| app
: string
| orders
: string
| orders_amount_positive
: string
| CHECK ((amount > 0))
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
| app
: string
| orders
: string
| orders_customer_id_fkey
: string
| app
: string
| customer
: string
| c
: string
| r
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES app.customer(id) ON DELETE CASCADE ON UPDATE RESTRICT
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
| app
: string
| customer
: int64 0
: string
| a customer
< row
: string
| app
: string
| orders
: int64 0
: string
| an order
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
| app
: string
| customer_orders
: string
| SELECT c.id, c.name, o.id AS order_id, o.amount FROM omni_testsyncdatabase.app.customer AS c JOIN omni_testsyncdatabase.app.orders AS o ON o.customer_id = c.id
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
| app
: string
| customer_orders
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| app
: string
| customer
: string
| id
< row
: string
| app
: string
| customer
: string
| name
< row
: string
| app
: string
| orders
: string
| amount
< row
: string
| app
: string
| orders
: string
| customer_id
< row
: string
| app
: string
| orders
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
< row
: string
| app
: string
| order_totals
: string
| SELECT customer_id, sum(amount) AS total FROM omni_testsyncdatabase.app.orders GROUP BY customer_id
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
| app
: string
| order_totals
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| app
: string
| orders
: string
| amount
< row
: string
| app
: string
| orders
: string
| customer_id
> query
|
| select n.nspname as function_schema,
| 	p.proname as function_name,
| 	pg_get_function_identity_arguments(p.oid) as function_signature,
| 	case when l.lanname = 'internal' then p.prosrc
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
| app
: string
| add
: string
| int8, int8
: string
| CREATE FUNCTION app.add(a INT8, b INT8)
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT a + b;
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
| app
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
| audit
: string
| big_seq
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
< row
: string
| app
: string
| mood
: string
| happy
< row
: string
| app
: string
| mood
: string
| sad
< row
: string
| app
: string
| mood
: string
| it's complicated
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
