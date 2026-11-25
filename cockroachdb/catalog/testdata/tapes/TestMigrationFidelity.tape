-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestMigrationFidelity default#1
> exec
| CREATE SCHEMA fid_m;
< ok 0 -
> exec
| CREATE SCHEMA fid_new;
< ok 0 -
> exec
| COMMENT ON SCHEMA fid_m IS 'migrated';
< ok 0 -
> exec
| CREATE TYPE fid_m.status AS ENUM ('new', 'paid', 'done', 'void');
< ok 0 -
> exec
| CREATE TYPE fid_m.level AS ENUM ('low', 'mid', 'high');
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.ids START WITH 10 INCREMENT BY 2 CACHE 4;
< ok 0 -
> exec
| CREATE TABLE fid_m.customers (
|     id INT8 PRIMARY KEY,
|     email STRING,
|     nick STRING NOT VISIBLE,
|     age INT8 DEFAULT 18 NOT NULL,
|     joined TIMESTAMPTZ DEFAULT now(),
|     CONSTRAINT customers_age_check CHECK (age >= 18),
|     CONSTRAINT customers_email_check CHECK (email != '')
| );
< ok 0 -
> exec
| COMMENT ON TABLE fid_m.customers IS 'new';
< ok 0 -
> exec
| COMMENT ON COLUMN fid_m.customers.email IS 'login';
< ok 0 -
> exec
| CREATE INDEX customers_nick_idx ON fid_m.customers (nick DESC);
< ok 0 -
> exec
| CREATE UNIQUE INDEX customers_joined_key ON fid_m.customers (joined) STORING (email);
< ok 0 -
> exec
| COMMENT ON INDEX fid_m.customers@customers_joined_key IS 'joins';
< ok 0 -
> exec
| CREATE TABLE fid_m.orders (
|     id INT8 NOT NULL,
|     customer_id INT8,
|     status fid_m.status DEFAULT 'paid',
|     line INT8 NOT NULL DEFAULT 1,
|     CONSTRAINT orders_pkey PRIMARY KEY (id, line),
|     CONSTRAINT orders_customer_fk FOREIGN KEY (customer_id) REFERENCES fid_m.customers (id) ON DELETE CASCADE
| );
< ok 0 -
> exec
| CREATE TABLE fid_new.audit (at TIMESTAMPTZ, customer_id INT8 REFERENCES fid_m.customers (id), id INT8 DEFAULT nextval('fid_m.ids'));
< ok 0 -
> exec
| COMMENT ON TABLE fid_new.audit IS e'first;\n\nsecond';
< ok 0 -
> exec
| CREATE TYPE fid_m."Mood" AS ENUM ('ok');
< ok 0 -
> exec
| CREATE TABLE fid_m.moods (id INT8 PRIMARY KEY, m fid_m."Mood" DEFAULT 'ok');
< ok 0 -
> exec
| CREATE TABLE fid_m.rekey (replacement INT8 PRIMARY KEY);
< ok 0 -
> exec
| CREATE TABLE fid_m.tickets (id INT8 PRIMARY KEY, n INT8);
< ok 0 -
> exec
| CREATE TABLE fid_m.codes (id INT8 PRIMARY KEY, code STRING NOT NULL, label STRING, UNIQUE INDEX codes_code_key (code) STORING (label));
< ok 0 -
> exec
| CREATE TABLE fid_m.uses (id INT8 PRIMARY KEY, code STRING REFERENCES fid_m.codes (code));
< ok 0 -
> exec
| CREATE TABLE fid_m.amounts (id INT8 PRIMARY KEY, amount STRING, INDEX amounts_amount_idx (amount));
< ok 0 -
> exec
| CREATE TABLE fid_m.prices (id INT8 PRIMARY KEY, price FLOAT8, CONSTRAINT prices_price_check CHECK (price > 0));
< ok 0 -
> exec
| COMMENT ON CONSTRAINT prices_price_check ON fid_m.prices IS 'positive';
< ok 0 -
> exec
| CREATE FUNCTION fid_m.price_text() RETURNS STRING LANGUAGE SQL STABLE AS $$ SELECT max(price)::STRING FROM fid_m.prices $$;
< ok 0 -
> exec
| CREATE TABLE fid_m.owners (id INT8 PRIMARY KEY, b INT8);
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.owned_seq OWNED BY fid_m.owners.b;
< ok 0 -
> exec
| CREATE TABLE fid_m.calc (id INT8 PRIMARY KEY);
< ok 0 -
> exec
| CREATE TABLE fid_m.paints (id INT8 PRIMARY KEY, shade STRING);
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.cached_seq CACHE 8;
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.reset_seq;
< ok 0 -
> exec
| CREATE TABLE fid_m.rekey2 (id INT8, replacement INT8 PRIMARY KEY);
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
| omni_testmigrationfidelity
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
| SELECT nspname, obj_description(oid, 'pg_namespace')
| FROM pg_catalog.pg_namespace
| WHERE nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY nspname;
|
< result
< column "nspname" "NAME"
< column "obj_description" "TEXT"
< row
: string
| fid_m
: string
| migrated
< row
: string
| fid_new
: null
< row
: string
| public
: null
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
| 	cols.is_hidden,
| 	cols.generation_expression,
| 	(SELECT a.attgenerated FROM pg_catalog.pg_attribute a WHERE a.attrelid = c.oid AND a.attname = cols.column_name) AS attgenerated
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
< column "generation_expression" "TEXT"
< column "attgenerated" "CHAR"
< row
: string
| fid_m
: string
| amounts
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
: string
|
: string
|
< row
: string
| fid_m
: string
| amounts
: string
| amount
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
: string
|
: string
|
< row
: string
| fid_m
: string
| calc
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
: string
|
: string
|
< row
: string
| fid_m
: string
| codes
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
: string
|
: string
|
< row
: string
| fid_m
: string
| codes
: string
| code
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
: string
|
: string
|
< row
: string
| fid_m
: string
| codes
: string
| label
: string
| text
: null
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
| text
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
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
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| email
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
: string
| login
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| nick
: string
| text
: null
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
| text
: null
: string
| YES
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| age
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
: string
| 18
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
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| joined
: string
| timestamp with time zone
: null
: null
: null
: int64 5
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
: string
|
: string
|
< row
: string
| fid_m
: string
| moods
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
: string
|
: string
|
< row
: string
| fid_m
: string
| moods
: string
| m
: string
| USER-DEFINED
: null
: null
: null
: int64 2
: string
| 'ok'
: string
| YES
: null
: string
| fid_m
: string
| Mood
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
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
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
: string
| customer_id
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
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
: string
| status
: string
| USER-DEFINED
: null
: null
: null
: int64 3
: string
| 'paid'
: string
| YES
: null
: string
| fid_m
: string
| status
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
: string
| line
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
: string
| 1
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
: string
|
: string
|
< row
: string
| fid_m
: string
| owners
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
: string
|
: string
|
< row
: string
| fid_m
: string
| owners
: string
| b
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
: string
|
: string
|
< row
: string
| fid_m
: string
| paints
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
: string
|
: string
|
< row
: string
| fid_m
: string
| paints
: string
| shade
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
: string
|
: string
|
< row
: string
| fid_m
: string
| prices
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
: string
|
: string
|
< row
: string
| fid_m
: string
| prices
: string
| price
: string
| double precision
: null
: int64 53
: null
: int64 2
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
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey
: string
| replacement
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
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey2
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
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey2
: string
| replacement
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| tickets
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
: string
|
: string
|
< row
: string
| fid_m
: string
| tickets
: string
| n
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
: string
|
: string
|
< row
: string
| fid_m
: string
| uses
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
: string
|
: string
|
< row
: string
| fid_m
: string
| uses
: string
| code
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
: string
|
: string
|
< row
: string
| fid_new
: string
| audit
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
: string
|
: string
|
< row
: string
| fid_new
: string
| audit
: string
| customer_id
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
: string
|
: string
|
< row
: string
| fid_new
: string
| audit
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| nextval('fid_m.ids'::REGCLASS)
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
: string
|
: string
|
< row
: string
| fid_new
: string
| audit
: string
| rowid
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
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
: string
|
: string
|
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
| fid_m
: string
| amounts
: string
| amounts_amount_idx
: string
| CREATE INDEX amounts_amount_idx ON omni_testmigrationfidelity.fid_m.amounts USING btree (amount ASC)
: null
: null
< row
: string
| fid_m
: string
| amounts
: string
| amounts_pkey
: string
| CREATE UNIQUE INDEX amounts_pkey ON omni_testmigrationfidelity.fid_m.amounts USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| cached_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.cached_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| calc
: string
| calc_pkey
: string
| CREATE UNIQUE INDEX calc_pkey ON omni_testmigrationfidelity.fid_m.calc USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| codes
: string
| codes_code_key
: string
| CREATE UNIQUE INDEX codes_code_key ON omni_testmigrationfidelity.fid_m.codes USING btree (code ASC) STORING (label)
: null
: null
< row
: string
| fid_m
: string
| codes
: string
| codes_pkey
: string
| CREATE UNIQUE INDEX codes_pkey ON omni_testmigrationfidelity.fid_m.codes USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_joined_key
: string
| CREATE UNIQUE INDEX customers_joined_key ON omni_testmigrationfidelity.fid_m.customers USING btree (joined ASC) STORING (email)
: null
: string
| joins
< row
: string
| fid_m
: string
| customers
: string
| customers_nick_idx
: string
| CREATE INDEX customers_nick_idx ON omni_testmigrationfidelity.fid_m.customers USING btree (nick DESC)
: null
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_pkey
: string
| CREATE UNIQUE INDEX customers_pkey ON omni_testmigrationfidelity.fid_m.customers USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| ids
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.ids USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| moods
: string
| moods_pkey
: string
| CREATE UNIQUE INDEX moods_pkey ON omni_testmigrationfidelity.fid_m.moods USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| orders
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON omni_testmigrationfidelity.fid_m.orders USING btree (id ASC, line ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| owned_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.owned_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| owners
: string
| owners_pkey
: string
| CREATE UNIQUE INDEX owners_pkey ON omni_testmigrationfidelity.fid_m.owners USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| paints
: string
| paints_pkey
: string
| CREATE UNIQUE INDEX paints_pkey ON omni_testmigrationfidelity.fid_m.paints USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| prices
: string
| prices_pkey
: string
| CREATE UNIQUE INDEX prices_pkey ON omni_testmigrationfidelity.fid_m.prices USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| rekey
: string
| rekey_pkey
: string
| CREATE UNIQUE INDEX rekey_pkey ON omni_testmigrationfidelity.fid_m.rekey USING btree (replacement ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| rekey2
: string
| rekey2_pkey
: string
| CREATE UNIQUE INDEX rekey2_pkey ON omni_testmigrationfidelity.fid_m.rekey2 USING btree (replacement ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| reset_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.reset_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| tickets
: string
| tickets_pkey
: string
| CREATE UNIQUE INDEX tickets_pkey ON omni_testmigrationfidelity.fid_m.tickets USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| uses
: string
| uses_pkey
: string
| CREATE UNIQUE INDEX uses_pkey ON omni_testmigrationfidelity.fid_m.uses USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_new
: string
| audit
: string
| audit_pkey
: string
| CREATE UNIQUE INDEX audit_pkey ON omni_testmigrationfidelity.fid_new.audit USING btree (rowid ASC)
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
| fid_m
: string
| customers
: string
| customers_age_check
: string
| CHECK ((age >= 18))
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_email_check
: string
| CHECK ((email != ''::STRING))
: null
< row
: string
| fid_m
: string
| prices
: string
| prices_price_check
: string
| CHECK ((price > 0.0))
: string
| positive
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
| fid_m
: string
| orders
: string
| orders_customer_fk
: string
| fid_m
: string
| customers
: string
| c
: string
| a
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES fid_m.customers(id) ON DELETE CASCADE
< row
: string
| fid_m
: string
| uses
: string
| uses_code_fkey
: string
| fid_m
: string
| codes
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY (code) REFERENCES fid_m.codes(code)
< row
: string
| fid_new
: string
| audit
: string
| audit_customer_id_fkey
: string
| fid_m
: string
| customers
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES fid_m.customers(id)
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
| fid_m
: string
| amounts
: int64 0
: null
< row
: string
| fid_m
: string
| calc
: int64 0
: null
< row
: string
| fid_m
: string
| codes
: int64 0
: null
< row
: string
| fid_m
: string
| customers
: int64 0
: string
| new
< row
: string
| fid_m
: string
| moods
: int64 0
: null
< row
: string
| fid_m
: string
| orders
: int64 0
: null
< row
: string
| fid_m
: string
| owners
: int64 0
: null
< row
: string
| fid_m
: string
| paints
: int64 0
: null
< row
: string
| fid_m
: string
| prices
: int64 0
: null
< row
: string
| fid_m
: string
| rekey
: int64 0
: null
< row
: string
| fid_m
: string
| rekey2
: int64 0
: null
< row
: string
| fid_m
: string
| tickets
: int64 0
: null
< row
: string
| fid_m
: string
| uses
: int64 0
: null
< row
: string
| fid_new
: string
| audit
: int64 0
: string
| first;
|
| second
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
< row
: string
| fid_m
: string
| price_text
: string
| price_text()
: string
| CREATE FUNCTION fid_m.price_text()
| 	RETURNS STRING
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT max(price)::STRING FROM omni_testmigrationfidelity.fid_m.prices;
| $$
> query
|
| 		SELECT s.sequence_schema, s.sequence_name, s.data_type, s.start_value, s.minimum_value, s.maximum_value, s.increment, s.cycle_option,
| 			(SELECT t.relname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_class t ON t.oid = d.refobjid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT a.attname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT cs.create_statement FROM crdb_internal.create_statements cs WHERE cs.descriptor_id = c.oid::INT8)
| 		FROM information_schema.sequences s
| 		LEFT JOIN pg_catalog.pg_namespace n ON n.nspname = s.sequence_schema
| 		LEFT JOIN pg_catalog.pg_class c ON c.relname = s.sequence_name AND c.relnamespace = n.oid AND c.relkind = 'S'
| 		ORDER BY s.sequence_schema, s.sequence_name;
< result
< column "sequence_schema" "TEXT"
< column "sequence_name" "TEXT"
< column "data_type" "TEXT"
< column "start_value" "TEXT"
< column "minimum_value" "TEXT"
< column "maximum_value" "TEXT"
< column "increment" "TEXT"
< column "cycle_option" "TEXT"
< column "relname" "NAME"
< column "attname" "NAME"
< column "create_statement" "TEXT"
< row
: string
| fid_m
: string
| cached_seq
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
: null
: null
: string
| CREATE SEQUENCE fid_m.cached_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1 CACHE 8
< row
: string
| fid_m
: string
| ids
: string
| INT8
: string
| 10
: string
| 1
: string
| 9223372036854775807
: string
| 2
: string
| NO
: null
: null
: string
| CREATE SEQUENCE fid_m.ids MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 2 START 10 CACHE 4
< row
: string
| fid_m
: string
| owned_seq
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
: string
| owners
: string
| b
: string
| CREATE SEQUENCE fid_m.owned_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1
< row
: string
| fid_m
: string
| reset_seq
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
: null
: null
: string
| CREATE SEQUENCE fid_m.reset_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1
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
| fid_m
: string
| Mood
: string
| ok
< row
: string
| fid_m
: string
| level
: string
| low
< row
: string
| fid_m
: string
| level
: string
| mid
< row
: string
| fid_m
: string
| level
: string
| high
< row
: string
| fid_m
: string
| status
: string
| new
< row
: string
| fid_m
: string
| status
: string
| paid
< row
: string
| fid_m
: string
| status
: string
| done
< row
: string
| fid_m
: string
| status
: string
| void
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
| DROP SCHEMA IF EXISTS "fid_m" CASCADE;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "fid_new" CASCADE;
< ok 0 -
> exec
| CREATE SCHEMA fid_m;
< ok 0 -
> exec
| CREATE TYPE fid_m.status AS ENUM ('new', 'done');
< ok 0 -
> exec
| CREATE TYPE fid_m.level AS ENUM ('high');
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.ids START WITH 10;
< ok 0 -
> exec
| CREATE TABLE fid_m.customers (
|     id INT8 PRIMARY KEY,
|     email STRING NOT NULL,
|     nick STRING,
|     age INT4 DEFAULT 0,
|     legacy STRING,
|     CONSTRAINT customers_age_check CHECK (age >= 0)
| );
< ok 0 -
> exec
| COMMENT ON TABLE fid_m.customers IS 'old';
< ok 0 -
> exec
| COMMENT ON COLUMN fid_m.customers.legacy IS 'going';
< ok 0 -
> exec
| CREATE INDEX customers_nick_idx ON fid_m.customers (nick);
< ok 0 -
> exec
| CREATE INDEX customers_email_idx ON fid_m.customers (email);
< ok 0 -
> exec
| CREATE TABLE fid_m.orders (
|     id INT8 NOT NULL,
|     customer_id INT8 REFERENCES fid_m.customers (id),
|     status fid_m.status DEFAULT 'new',
|     CONSTRAINT orders_pkey PRIMARY KEY (id)
| );
< ok 0 -
> exec
| CREATE TABLE fid_m.gone (id INT8 PRIMARY KEY, customer_id INT8 REFERENCES fid_m.customers (id));
< ok 0 -
> exec
| CREATE TYPE fid_m."Mood" AS ENUM ('ok');
< ok 0 -
> exec
| CREATE TABLE fid_m.rekey (id INT8 PRIMARY KEY, replacement INT8 NOT NULL);
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.ticket_seq;
< ok 0 -
> exec
| CREATE TABLE fid_m.tickets (id INT8 PRIMARY KEY, n INT8 DEFAULT nextval('fid_m.ticket_seq'));
< ok 0 -
> exec
| CREATE TABLE fid_m.codes (id INT8 PRIMARY KEY, code STRING NOT NULL, label STRING, UNIQUE INDEX codes_code_key (code));
< ok 0 -
> exec
| CREATE TABLE fid_m.uses (id INT8 PRIMARY KEY, code STRING REFERENCES fid_m.codes (code));
< ok 0 -
> exec
| CREATE TABLE fid_m.amounts (id INT8 PRIMARY KEY, amount DECIMAL, INDEX amounts_amount_idx (amount));
< ok 0 -
> exec
| CREATE TABLE fid_m.prices (id INT8 PRIMARY KEY, price DECIMAL, CONSTRAINT prices_price_check CHECK (price > 0));
< ok 0 -
> exec
| COMMENT ON CONSTRAINT prices_price_check ON fid_m.prices IS 'positive';
< ok 0 -
> exec
| CREATE FUNCTION fid_m.price_text() RETURNS STRING LANGUAGE SQL STABLE AS $$ SELECT max(price)::STRING FROM fid_m.prices $$;
< ok 0 -
> exec
| CREATE TABLE fid_m.owners (id INT8 PRIMARY KEY, a INT8, b INT8);
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.owned_seq OWNED BY fid_m.owners.a;
< ok 0 -
> exec
| CREATE TABLE fid_m.gone_owner (id INT8 PRIMARY KEY, n INT8);
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.gone_owner_seq OWNED BY fid_m.gone_owner.n;
< ok 0 -
> exec
| CREATE TABLE fid_m.calc (id INT8 PRIMARY KEY, a INT8, b INT8 AS (a + 1) STORED);
< ok 0 -
> exec
| CREATE TYPE fid_m.shade AS ENUM ('light', 'dark');
< ok 0 -
> exec
| CREATE TABLE fid_m.paints (id INT8 PRIMARY KEY, shade fid_m.shade);
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.cached_seq CACHE 4;
< ok 0 -
> exec
| CREATE SEQUENCE fid_m.reset_seq CACHE 4;
< ok 0 -
> exec
| CREATE TABLE fid_m.rekey2 (id INT8 PRIMARY KEY, replacement INT8 NOT NULL);
< ok 0 -
> exec
| CREATE SCHEMA fid_gone;
< ok 0 -
> exec
| CREATE TYPE fid_gone.kind AS ENUM ('x');
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
| omni_testmigrationfidelity
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
| SELECT nspname, obj_description(oid, 'pg_namespace')
| FROM pg_catalog.pg_namespace
| WHERE nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY nspname;
|
< result
< column "nspname" "NAME"
< column "obj_description" "TEXT"
< row
: string
| fid_gone
: null
< row
: string
| fid_m
: null
< row
: string
| public
: null
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
| 	cols.is_hidden,
| 	cols.generation_expression,
| 	(SELECT a.attgenerated FROM pg_catalog.pg_attribute a WHERE a.attrelid = c.oid AND a.attname = cols.column_name) AS attgenerated
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
< column "generation_expression" "TEXT"
< column "attgenerated" "CHAR"
< row
: string
| fid_m
: string
| amounts
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
: string
|
: string
|
< row
: string
| fid_m
: string
| amounts
: string
| amount
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
: string
|
: string
|
< row
: string
| fid_m
: string
| calc
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
: string
|
: string
|
< row
: string
| fid_m
: string
| calc
: string
| a
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
: string
|
: string
|
< row
: string
| fid_m
: string
| calc
: string
| b
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
| a + 1
: string
| s
< row
: string
| fid_m
: string
| codes
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
: string
|
: string
|
< row
: string
| fid_m
: string
| codes
: string
| code
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
: string
|
: string
|
< row
: string
| fid_m
: string
| codes
: string
| label
: string
| text
: null
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
| text
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
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
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| email
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
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| nick
: string
| text
: null
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
| text
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| age
: string
| integer
: null
: int64 32
: int64 0
: int64 4
: string
| 0
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
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| legacy
: string
| text
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
| text
: string
| going
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| gone
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
: string
|
: string
|
< row
: string
| fid_m
: string
| gone
: string
| customer_id
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
: string
|
: string
|
< row
: string
| fid_m
: string
| gone_owner
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
: string
|
: string
|
< row
: string
| fid_m
: string
| gone_owner
: string
| n
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
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
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
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
: string
| customer_id
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
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
: string
| status
: string
| USER-DEFINED
: null
: null
: null
: int64 3
: string
| 'new'
: string
| YES
: null
: string
| fid_m
: string
| status
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| owners
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
: string
|
: string
|
< row
: string
| fid_m
: string
| owners
: string
| a
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
: string
|
: string
|
< row
: string
| fid_m
: string
| owners
: string
| b
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| paints
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
: string
|
: string
|
< row
: string
| fid_m
: string
| paints
: string
| shade
: string
| USER-DEFINED
: null
: null
: null
: int64 2
: null
: string
| YES
: null
: string
| fid_m
: string
| shade
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| prices
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
: string
|
: string
|
< row
: string
| fid_m
: string
| prices
: string
| price
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
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey
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
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey
: string
| replacement
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey2
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
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey2
: string
| replacement
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| tickets
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
: string
|
: string
|
< row
: string
| fid_m
: string
| tickets
: string
| n
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| nextval('fid_m.ticket_seq'::REGCLASS)
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
: string
|
: string
|
< row
: string
| fid_m
: string
| uses
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
: string
|
: string
|
< row
: string
| fid_m
: string
| uses
: string
| code
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
: string
|
: string
|
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
| fid_m
: string
| amounts
: string
| amounts_amount_idx
: string
| CREATE INDEX amounts_amount_idx ON omni_testmigrationfidelity.fid_m.amounts USING btree (amount ASC)
: null
: null
< row
: string
| fid_m
: string
| amounts
: string
| amounts_pkey
: string
| CREATE UNIQUE INDEX amounts_pkey ON omni_testmigrationfidelity.fid_m.amounts USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| cached_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.cached_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| calc
: string
| calc_pkey
: string
| CREATE UNIQUE INDEX calc_pkey ON omni_testmigrationfidelity.fid_m.calc USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| codes
: string
| codes_code_key
: string
| CREATE UNIQUE INDEX codes_code_key ON omni_testmigrationfidelity.fid_m.codes USING btree (code ASC)
: null
: null
< row
: string
| fid_m
: string
| codes
: string
| codes_pkey
: string
| CREATE UNIQUE INDEX codes_pkey ON omni_testmigrationfidelity.fid_m.codes USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_email_idx
: string
| CREATE INDEX customers_email_idx ON omni_testmigrationfidelity.fid_m.customers USING btree (email ASC)
: null
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_nick_idx
: string
| CREATE INDEX customers_nick_idx ON omni_testmigrationfidelity.fid_m.customers USING btree (nick ASC)
: null
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_pkey
: string
| CREATE UNIQUE INDEX customers_pkey ON omni_testmigrationfidelity.fid_m.customers USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| gone
: string
| gone_pkey
: string
| CREATE UNIQUE INDEX gone_pkey ON omni_testmigrationfidelity.fid_m.gone USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| gone_owner
: string
| gone_owner_pkey
: string
| CREATE UNIQUE INDEX gone_owner_pkey ON omni_testmigrationfidelity.fid_m.gone_owner USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| gone_owner_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.gone_owner_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| ids
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.ids USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| orders
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON omni_testmigrationfidelity.fid_m.orders USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| owned_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.owned_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| owners
: string
| owners_pkey
: string
| CREATE UNIQUE INDEX owners_pkey ON omni_testmigrationfidelity.fid_m.owners USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| paints
: string
| paints_pkey
: string
| CREATE UNIQUE INDEX paints_pkey ON omni_testmigrationfidelity.fid_m.paints USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| prices
: string
| prices_pkey
: string
| CREATE UNIQUE INDEX prices_pkey ON omni_testmigrationfidelity.fid_m.prices USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| rekey
: string
| rekey_pkey
: string
| CREATE UNIQUE INDEX rekey_pkey ON omni_testmigrationfidelity.fid_m.rekey USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| rekey2
: string
| rekey2_pkey
: string
| CREATE UNIQUE INDEX rekey2_pkey ON omni_testmigrationfidelity.fid_m.rekey2 USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| reset_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.reset_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| ticket_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.ticket_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| tickets
: string
| tickets_pkey
: string
| CREATE UNIQUE INDEX tickets_pkey ON omni_testmigrationfidelity.fid_m.tickets USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| uses
: string
| uses_pkey
: string
| CREATE UNIQUE INDEX uses_pkey ON omni_testmigrationfidelity.fid_m.uses USING btree (id ASC)
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
| fid_m
: string
| customers
: string
| customers_age_check
: string
| CHECK ((age >= 0))
: null
< row
: string
| fid_m
: string
| prices
: string
| prices_price_check
: string
| CHECK ((price > 0))
: string
| positive
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
| fid_m
: string
| gone
: string
| gone_customer_id_fkey
: string
| fid_m
: string
| customers
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES fid_m.customers(id)
< row
: string
| fid_m
: string
| orders
: string
| orders_customer_id_fkey
: string
| fid_m
: string
| customers
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES fid_m.customers(id)
< row
: string
| fid_m
: string
| uses
: string
| uses_code_fkey
: string
| fid_m
: string
| codes
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY (code) REFERENCES fid_m.codes(code)
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
| fid_m
: string
| amounts
: int64 0
: null
< row
: string
| fid_m
: string
| calc
: int64 0
: null
< row
: string
| fid_m
: string
| codes
: int64 0
: null
< row
: string
| fid_m
: string
| customers
: int64 0
: string
| old
< row
: string
| fid_m
: string
| gone
: int64 0
: null
< row
: string
| fid_m
: string
| gone_owner
: int64 0
: null
< row
: string
| fid_m
: string
| orders
: int64 0
: null
< row
: string
| fid_m
: string
| owners
: int64 0
: null
< row
: string
| fid_m
: string
| paints
: int64 0
: null
< row
: string
| fid_m
: string
| prices
: int64 0
: null
< row
: string
| fid_m
: string
| rekey
: int64 0
: null
< row
: string
| fid_m
: string
| rekey2
: int64 0
: null
< row
: string
| fid_m
: string
| tickets
: int64 0
: null
< row
: string
| fid_m
: string
| uses
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
< row
: string
| fid_m
: string
| price_text
: string
| price_text()
: string
| CREATE FUNCTION fid_m.price_text()
| 	RETURNS STRING
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT max(price)::STRING FROM omni_testmigrationfidelity.fid_m.prices;
| $$
> query
|
| 		SELECT s.sequence_schema, s.sequence_name, s.data_type, s.start_value, s.minimum_value, s.maximum_value, s.increment, s.cycle_option,
| 			(SELECT t.relname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_class t ON t.oid = d.refobjid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT a.attname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT cs.create_statement FROM crdb_internal.create_statements cs WHERE cs.descriptor_id = c.oid::INT8)
| 		FROM information_schema.sequences s
| 		LEFT JOIN pg_catalog.pg_namespace n ON n.nspname = s.sequence_schema
| 		LEFT JOIN pg_catalog.pg_class c ON c.relname = s.sequence_name AND c.relnamespace = n.oid AND c.relkind = 'S'
| 		ORDER BY s.sequence_schema, s.sequence_name;
< result
< column "sequence_schema" "TEXT"
< column "sequence_name" "TEXT"
< column "data_type" "TEXT"
< column "start_value" "TEXT"
< column "minimum_value" "TEXT"
< column "maximum_value" "TEXT"
< column "increment" "TEXT"
< column "cycle_option" "TEXT"
< column "relname" "NAME"
< column "attname" "NAME"
< column "create_statement" "TEXT"
< row
: string
| fid_m
: string
| cached_seq
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
: null
: null
: string
| CREATE SEQUENCE fid_m.cached_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1 CACHE 4
< row
: string
| fid_m
: string
| gone_owner_seq
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
: string
| gone_owner
: string
| n
: string
| CREATE SEQUENCE fid_m.gone_owner_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1
< row
: string
| fid_m
: string
| ids
: string
| INT8
: string
| 10
: string
| 1
: string
| 9223372036854775807
: string
| 1
: string
| NO
: null
: null
: string
| CREATE SEQUENCE fid_m.ids MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10
< row
: string
| fid_m
: string
| owned_seq
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
: string
| owners
: string
| a
: string
| CREATE SEQUENCE fid_m.owned_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1
< row
: string
| fid_m
: string
| reset_seq
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
: null
: null
: string
| CREATE SEQUENCE fid_m.reset_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1 CACHE 4
< row
: string
| fid_m
: string
| ticket_seq
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
: null
: null
: string
| CREATE SEQUENCE fid_m.ticket_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1
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
| fid_gone
: string
| kind
: string
| x
< row
: string
| fid_m
: string
| Mood
: string
| ok
< row
: string
| fid_m
: string
| level
: string
| high
< row
: string
| fid_m
: string
| shade
: string
| light
< row
: string
| fid_m
: string
| shade
: string
| dark
< row
: string
| fid_m
: string
| status
: string
| new
< row
: string
| fid_m
: string
| status
: string
| done
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
| ALTER SEQUENCE "fid_m"."owned_seq" OWNED BY NONE;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."tickets" ALTER COLUMN "n" DROP DEFAULT;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."gone" DROP CONSTRAINT "gone_customer_id_fkey";
< ok 0 -
> exec
| ALTER TABLE "fid_m"."orders" DROP CONSTRAINT "orders_customer_id_fkey";
< ok 0 -
> exec
| ALTER TABLE "fid_m"."uses" DROP CONSTRAINT "uses_code_fkey";
< ok 0 -
> exec
| DROP INDEX "fid_m"."amounts"@"amounts_amount_idx";
< ok 0 -
> exec
| DROP INDEX "fid_m"."codes"@"codes_code_key" CASCADE;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" DROP CONSTRAINT "customers_age_check";
< ok 0 -
> exec
| DROP INDEX "fid_m"."customers"@"customers_email_idx";
< ok 0 -
> exec
| DROP INDEX "fid_m"."customers"@"customers_nick_idx";
< ok 0 -
> exec
| ALTER TABLE "fid_m"."prices" DROP CONSTRAINT "prices_price_check";
< ok 0 -
> exec
| DROP TABLE "fid_m"."gone";
< ok 0 -
> exec
| DROP TABLE "fid_m"."gone_owner";
< ok 0 -
> exec
| DROP FUNCTION "fid_m"."price_text"();
< ok 0 -
> exec
| ALTER TABLE "fid_m"."calc" DROP COLUMN "b";
< ok 0 -
> exec
| ALTER TABLE "fid_m"."calc" DROP COLUMN "a";
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" DROP COLUMN "legacy";
< ok 0 -
> exec
| ALTER TABLE "fid_m"."owners" DROP COLUMN "a";
< ok 0 -
> exec
| DROP SEQUENCE IF EXISTS "fid_m"."gone_owner_seq";
< ok 0 -
> exec
| DROP SEQUENCE IF EXISTS "fid_m"."ticket_seq";
< ok 0 -
> exec
| CREATE SCHEMA "fid_new";
< ok 0 -
> exec
| ALTER TYPE "fid_m"."level" ADD VALUE 'low' BEFORE 'high';
< ok 0 -
> exec
| ALTER TYPE "fid_m"."level" ADD VALUE 'mid' AFTER 'low';
< ok 0 -
> exec
| ALTER TYPE "fid_m"."status" ADD VALUE 'paid' AFTER 'new';
< ok 0 -
> exec
| ALTER TYPE "fid_m"."status" ADD VALUE 'void' AFTER 'done';
< ok 0 -
> exec
| ALTER SEQUENCE "fid_m"."cached_seq" CACHE 8;
< ok 0 -
> exec
| ALTER SEQUENCE "fid_m"."ids" INCREMENT BY 2 CACHE 4;
< ok 0 -
> exec
| ALTER SEQUENCE "fid_m"."reset_seq" CACHE 1;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."amounts" ALTER COLUMN "amount" SET DATA TYPE text;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" ALTER COLUMN "email" DROP NOT NULL;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" ALTER COLUMN "nick" SET NOT VISIBLE;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" ALTER COLUMN "age" SET DATA TYPE bigint;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" ALTER COLUMN "age" SET DEFAULT 18;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" ALTER COLUMN "age" SET NOT NULL;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" ADD COLUMN "joined" timestamp with time zone DEFAULT now();
< ok 0 -
> exec
| ALTER TABLE "fid_m"."orders" ALTER COLUMN "status" SET DEFAULT 'paid';
< ok 0 -
> exec
| ALTER TABLE "fid_m"."orders" ADD COLUMN "line" bigint DEFAULT 1 NOT NULL;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."orders" DROP CONSTRAINT "orders_pkey", ADD CONSTRAINT "orders_pkey" PRIMARY KEY (id ASC, line ASC);
< ok 0 -
> exec
| ALTER TABLE "fid_m"."paints" ALTER COLUMN "shade" SET DATA TYPE text;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."prices" ALTER COLUMN "price" SET DATA TYPE double precision;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."rekey" DROP CONSTRAINT "rekey_pkey", ADD CONSTRAINT "rekey_pkey" PRIMARY KEY (replacement ASC);
< ok 0 -
> exec
| ALTER TABLE "fid_m"."rekey" DROP COLUMN "id";
< ok 0 -
> exec
| ALTER TABLE "fid_m"."rekey2" DROP CONSTRAINT "rekey2_pkey", ADD CONSTRAINT "rekey2_pkey" PRIMARY KEY (replacement ASC);
< ok 0 -
> exec
| ALTER TABLE "fid_m"."rekey2" ALTER COLUMN "id" DROP NOT NULL;
< ok 0 -
> exec
| CREATE TABLE "fid_m"."moods" (
|     "id" bigint NOT NULL,
|     "m" "fid_m"."Mood" DEFAULT 'ok',
|     CONSTRAINT "moods_pkey" PRIMARY KEY (id ASC)
| );
< ok 0 -
> exec
| CREATE FUNCTION fid_m.price_text()
| 	RETURNS STRING
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| SELECT max(price)::STRING FROM fid_m.prices;
| $$;
< ok 0 -
> exec
| CREATE TABLE "fid_new"."audit" (
|     "at" timestamp with time zone,
|     "customer_id" bigint,
|     "id" bigint DEFAULT nextval('fid_m.ids'::REGCLASS)
| );
< ok 0 -
> exec
| DROP TYPE "fid_m"."shade";
< ok 0 -
> exec
| DROP TYPE "fid_gone"."kind";
< ok 0 -
> exec
| DROP SCHEMA "fid_gone";
< ok 0 -
> exec
| ALTER SEQUENCE "fid_m"."owned_seq" OWNED BY "fid_m"."owners"."b";
< ok 0 -
> exec
| CREATE INDEX amounts_amount_idx ON fid_m.amounts (amount ASC);
< ok 0 -
> exec
| CREATE UNIQUE INDEX codes_code_key ON fid_m.codes (code ASC) STORING (label);
< ok 0 -
> exec
| CREATE UNIQUE INDEX customers_joined_key ON fid_m.customers (joined ASC) STORING (email);
< ok 0 -
> exec
| CREATE INDEX customers_nick_idx ON fid_m.customers (nick DESC);
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" ADD CONSTRAINT "customers_age_check" CHECK ((age >= 18));
< ok 0 -
> exec
| ALTER TABLE "fid_m"."customers" ADD CONSTRAINT "customers_email_check" CHECK ((email != ''::STRING));
< ok 0 -
> exec
| ALTER TABLE "fid_m"."prices" ADD CONSTRAINT "prices_price_check" CHECK ((price > 0.0));
< ok 0 -
> exec
| ALTER TABLE "fid_m"."orders"
|     ADD CONSTRAINT "orders_customer_fk" FOREIGN KEY ("customer_id")
|     REFERENCES "fid_m"."customers" ("id")
|     ON DELETE CASCADE;
< ok 0 -
> exec
| ALTER TABLE "fid_m"."uses"
|     ADD CONSTRAINT "uses_code_fkey" FOREIGN KEY ("code")
|     REFERENCES "fid_m"."codes" ("code");
< ok 0 -
> exec
| ALTER TABLE "fid_new"."audit"
|     ADD CONSTRAINT "audit_customer_id_fkey" FOREIGN KEY ("customer_id")
|     REFERENCES "fid_m"."customers" ("id");
< ok 0 -
> exec
| COMMENT ON SCHEMA "fid_m" IS 'migrated';
< ok 0 -
> exec
| COMMENT ON INDEX "fid_m"."customers"@"customers_joined_key" IS 'joins';
< ok 0 -
> exec
| COMMENT ON TABLE "fid_m"."customers" IS 'new';
< ok 0 -
> exec
| COMMENT ON COLUMN "fid_m"."customers"."email" IS 'login';
< ok 0 -
> exec
| COMMENT ON CONSTRAINT "prices_price_check" ON "fid_m"."prices" IS 'positive';
< ok 0 -
> exec
| COMMENT ON TABLE "fid_new"."audit" IS 'first;
|
| second';
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
| omni_testmigrationfidelity
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
| SELECT nspname, obj_description(oid, 'pg_namespace')
| FROM pg_catalog.pg_namespace
| WHERE nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY nspname;
|
< result
< column "nspname" "NAME"
< column "obj_description" "TEXT"
< row
: string
| fid_m
: string
| migrated
< row
: string
| fid_new
: null
< row
: string
| public
: null
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
| 	cols.is_hidden,
| 	cols.generation_expression,
| 	(SELECT a.attgenerated FROM pg_catalog.pg_attribute a WHERE a.attrelid = c.oid AND a.attname = cols.column_name) AS attgenerated
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
< column "generation_expression" "TEXT"
< column "attgenerated" "CHAR"
< row
: string
| fid_m
: string
| amounts
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
: string
|
: string
|
< row
: string
| fid_m
: string
| amounts
: string
| amount
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
: string
|
: string
|
< row
: string
| fid_m
: string
| calc
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
: string
|
: string
|
< row
: string
| fid_m
: string
| codes
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
: string
|
: string
|
< row
: string
| fid_m
: string
| codes
: string
| code
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
: string
|
: string
|
< row
: string
| fid_m
: string
| codes
: string
| label
: string
| text
: null
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
| text
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
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
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| email
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
: string
| login
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| nick
: string
| text
: null
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
| text
: null
: string
| YES
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| age
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
: string
| 18
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
: string
|
: string
|
< row
: string
| fid_m
: string
| customers
: string
| joined
: string
| timestamp with time zone
: null
: null
: null
: int64 6
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
: string
|
: string
|
< row
: string
| fid_m
: string
| moods
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
: string
|
: string
|
< row
: string
| fid_m
: string
| moods
: string
| m
: string
| USER-DEFINED
: null
: null
: null
: int64 2
: string
| 'ok'
: string
| YES
: null
: string
| fid_m
: string
| Mood
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
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
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
: string
| customer_id
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
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
: string
| status
: string
| USER-DEFINED
: null
: null
: null
: int64 3
: string
| 'paid'
: string
| YES
: null
: string
| fid_m
: string
| status
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| orders
: string
| line
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
: string
| 1
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
: string
|
: string
|
< row
: string
| fid_m
: string
| owners
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
: string
|
: string
|
< row
: string
| fid_m
: string
| owners
: string
| b
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| paints
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
: string
|
: string
|
< row
: string
| fid_m
: string
| paints
: string
| shade
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
: string
|
: string
|
< row
: string
| fid_m
: string
| prices
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
: string
|
: string
|
< row
: string
| fid_m
: string
| prices
: string
| price
: string
| double precision
: null
: int64 53
: null
: int64 2
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
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey
: string
| replacement
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey2
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
: string
|
: string
|
< row
: string
| fid_m
: string
| rekey2
: string
| replacement
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_m
: string
| tickets
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
: string
|
: string
|
< row
: string
| fid_m
: string
| tickets
: string
| n
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
: string
|
: string
|
< row
: string
| fid_m
: string
| uses
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
: string
|
: string
|
< row
: string
| fid_m
: string
| uses
: string
| code
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
: string
|
: string
|
< row
: string
| fid_new
: string
| audit
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
: string
|
: string
|
< row
: string
| fid_new
: string
| audit
: string
| customer_id
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
: string
|
: string
|
< row
: string
| fid_new
: string
| audit
: string
| id
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| nextval('fid_m.ids'::REGCLASS)
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
: string
|
: string
|
< row
: string
| fid_new
: string
| audit
: string
| rowid
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
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
: string
|
: string
|
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
| fid_m
: string
| amounts
: string
| amounts_amount_idx
: string
| CREATE INDEX amounts_amount_idx ON omni_testmigrationfidelity.fid_m.amounts USING btree (amount ASC)
: null
: null
< row
: string
| fid_m
: string
| amounts
: string
| amounts_pkey
: string
| CREATE UNIQUE INDEX amounts_pkey ON omni_testmigrationfidelity.fid_m.amounts USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| cached_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.cached_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| calc
: string
| calc_pkey
: string
| CREATE UNIQUE INDEX calc_pkey ON omni_testmigrationfidelity.fid_m.calc USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| codes
: string
| codes_code_key
: string
| CREATE UNIQUE INDEX codes_code_key ON omni_testmigrationfidelity.fid_m.codes USING btree (code ASC) STORING (label)
: null
: null
< row
: string
| fid_m
: string
| codes
: string
| codes_pkey
: string
| CREATE UNIQUE INDEX codes_pkey ON omni_testmigrationfidelity.fid_m.codes USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_joined_key
: string
| CREATE UNIQUE INDEX customers_joined_key ON omni_testmigrationfidelity.fid_m.customers USING btree (joined ASC) STORING (email)
: null
: string
| joins
< row
: string
| fid_m
: string
| customers
: string
| customers_nick_idx
: string
| CREATE INDEX customers_nick_idx ON omni_testmigrationfidelity.fid_m.customers USING btree (nick DESC)
: null
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_pkey
: string
| CREATE UNIQUE INDEX customers_pkey ON omni_testmigrationfidelity.fid_m.customers USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| ids
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.ids USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| moods
: string
| moods_pkey
: string
| CREATE UNIQUE INDEX moods_pkey ON omni_testmigrationfidelity.fid_m.moods USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| orders
: string
| orders_pkey
: string
| CREATE UNIQUE INDEX orders_pkey ON omni_testmigrationfidelity.fid_m.orders USING btree (id ASC, line ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| owned_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.owned_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| owners
: string
| owners_pkey
: string
| CREATE UNIQUE INDEX owners_pkey ON omni_testmigrationfidelity.fid_m.owners USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| paints
: string
| paints_pkey
: string
| CREATE UNIQUE INDEX paints_pkey ON omni_testmigrationfidelity.fid_m.paints USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| prices
: string
| prices_pkey
: string
| CREATE UNIQUE INDEX prices_pkey ON omni_testmigrationfidelity.fid_m.prices USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| rekey
: string
| rekey_pkey
: string
| CREATE UNIQUE INDEX rekey_pkey ON omni_testmigrationfidelity.fid_m.rekey USING btree (replacement ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| rekey2
: string
| rekey2_pkey
: string
| CREATE UNIQUE INDEX rekey2_pkey ON omni_testmigrationfidelity.fid_m.rekey2 USING btree (replacement ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| reset_seq
: string
| primary
: string
| CREATE UNIQUE INDEX "primary" ON omni_testmigrationfidelity.fid_m.reset_seq USING btree (value ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| tickets
: string
| tickets_pkey
: string
| CREATE UNIQUE INDEX tickets_pkey ON omni_testmigrationfidelity.fid_m.tickets USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_m
: string
| uses
: string
| uses_pkey
: string
| CREATE UNIQUE INDEX uses_pkey ON omni_testmigrationfidelity.fid_m.uses USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_new
: string
| audit
: string
| audit_pkey
: string
| CREATE UNIQUE INDEX audit_pkey ON omni_testmigrationfidelity.fid_new.audit USING btree (rowid ASC)
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
| fid_m
: string
| customers
: string
| customers_age_check
: string
| CHECK ((age >= 18))
: null
< row
: string
| fid_m
: string
| customers
: string
| customers_email_check
: string
| CHECK ((email != ''::STRING))
: null
< row
: string
| fid_m
: string
| prices
: string
| prices_price_check
: string
| CHECK ((price > 0.0))
: string
| positive
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
| fid_m
: string
| orders
: string
| orders_customer_fk
: string
| fid_m
: string
| customers
: string
| c
: string
| a
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES fid_m.customers(id) ON DELETE CASCADE
< row
: string
| fid_m
: string
| uses
: string
| uses_code_fkey
: string
| fid_m
: string
| codes
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY (code) REFERENCES fid_m.codes(code)
< row
: string
| fid_new
: string
| audit
: string
| audit_customer_id_fkey
: string
| fid_m
: string
| customers
: string
| a
: string
| a
: string
| s
: string
| FOREIGN KEY (customer_id) REFERENCES fid_m.customers(id)
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
| fid_m
: string
| amounts
: int64 0
: null
< row
: string
| fid_m
: string
| calc
: int64 0
: null
< row
: string
| fid_m
: string
| codes
: int64 0
: null
< row
: string
| fid_m
: string
| customers
: int64 0
: string
| new
< row
: string
| fid_m
: string
| moods
: int64 0
: null
< row
: string
| fid_m
: string
| orders
: int64 0
: null
< row
: string
| fid_m
: string
| owners
: int64 0
: null
< row
: string
| fid_m
: string
| paints
: int64 0
: null
< row
: string
| fid_m
: string
| prices
: int64 0
: null
< row
: string
| fid_m
: string
| rekey
: int64 0
: null
< row
: string
| fid_m
: string
| rekey2
: int64 0
: null
< row
: string
| fid_m
: string
| tickets
: int64 0
: null
< row
: string
| fid_m
: string
| uses
: int64 0
: null
< row
: string
| fid_new
: string
| audit
: int64 0
: string
| first;
|
| second
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
< row
: string
| fid_m
: string
| price_text
: string
| price_text()
: string
| CREATE FUNCTION fid_m.price_text()
| 	RETURNS STRING
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT max(price)::STRING FROM omni_testmigrationfidelity.fid_m.prices;
| $$
> query
|
| 		SELECT s.sequence_schema, s.sequence_name, s.data_type, s.start_value, s.minimum_value, s.maximum_value, s.increment, s.cycle_option,
| 			(SELECT t.relname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_class t ON t.oid = d.refobjid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT a.attname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT cs.create_statement FROM crdb_internal.create_statements cs WHERE cs.descriptor_id = c.oid::INT8)
| 		FROM information_schema.sequences s
| 		LEFT JOIN pg_catalog.pg_namespace n ON n.nspname = s.sequence_schema
| 		LEFT JOIN pg_catalog.pg_class c ON c.relname = s.sequence_name AND c.relnamespace = n.oid AND c.relkind = 'S'
| 		ORDER BY s.sequence_schema, s.sequence_name;
< result
< column "sequence_schema" "TEXT"
< column "sequence_name" "TEXT"
< column "data_type" "TEXT"
< column "start_value" "TEXT"
< column "minimum_value" "TEXT"
< column "maximum_value" "TEXT"
< column "increment" "TEXT"
< column "cycle_option" "TEXT"
< column "relname" "NAME"
< column "attname" "NAME"
< column "create_statement" "TEXT"
< row
: string
| fid_m
: string
| cached_seq
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
: null
: null
: string
| CREATE SEQUENCE fid_m.cached_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1 CACHE 8
< row
: string
| fid_m
: string
| ids
: string
| INT8
: string
| 10
: string
| 1
: string
| 9223372036854775807
: string
| 2
: string
| NO
: null
: null
: string
| CREATE SEQUENCE fid_m.ids MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 2 START 10 CACHE 4
< row
: string
| fid_m
: string
| owned_seq
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
: string
| owners
: string
| b
: string
| CREATE SEQUENCE fid_m.owned_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1
< row
: string
| fid_m
: string
| reset_seq
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
: null
: null
: string
| CREATE SEQUENCE fid_m.reset_seq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1
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
| fid_m
: string
| Mood
: string
| ok
< row
: string
| fid_m
: string
| level
: string
| low
< row
: string
| fid_m
: string
| level
: string
| mid
< row
: string
| fid_m
: string
| level
: string
| high
< row
: string
| fid_m
: string
| status
: string
| new
< row
: string
| fid_m
: string
| status
: string
| paid
< row
: string
| fid_m
: string
| status
: string
| done
< row
: string
| fid_m
: string
| status
: string
| void
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
| DROP SCHEMA IF EXISTS "fid_m" CASCADE;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "fid_new" CASCADE;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "fid_m" CASCADE;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "fid_new" CASCADE;
< ok 0 -
> exec
| CREATE SCHEMA fid_v;
< ok 0 -
> exec
| CREATE TABLE fid_v.items (id INT8 PRIMARY KEY, name STRING, price DECIMAL(10, 2), bonus INT8);
< ok 0 -
> exec
| CREATE VIEW fid_v.priced AS SELECT id, name FROM fid_v.items WHERE price > 1;
< ok 0 -
> exec
| CREATE VIEW fid_v.priced_names AS SELECT name FROM fid_v.priced;
< ok 0 -
> exec
| CREATE VIEW fid_v.unchanged AS SELECT id FROM fid_v.items;
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW fid_v.totals AS SELECT name, sum(price) AS total FROM fid_v.items GROUP BY name;
< ok 0 -
> exec
| CREATE INDEX totals_total_idx ON fid_v.totals (total);
< ok 0 -
> exec
| CREATE FUNCTION fid_v.add(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT a + b + 0 $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.count_items() RETURNS INT8 LANGUAGE SQL STABLE AS $$ SELECT count(*) FROM fid_v.priced $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.touch() RETURNS INT8 LANGUAGE PLPGSQL AS $$ BEGIN RETURN (SELECT count(*) FROM fid_v.items); END $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.bonus() RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT 5 $$;
< ok 0 -
> exec
| ALTER TABLE fid_v.items ALTER COLUMN bonus SET DEFAULT fid_v.bonus();
< ok 0 -
> exec
| CREATE TABLE fid_v.tagged (id INT8 PRIMARY KEY, v INT8);
< ok 0 -
> exec
| CREATE FUNCTION fid_v.count_old() RETURNS INT8 LANGUAGE SQL STABLE AS $$ SELECT count(*) FROM fid_v.items $$;
< ok 0 -
> exec
| CREATE VIEW fid_v.named (other) AS SELECT id FROM fid_v.items;
< ok 0 -
> exec
| CREATE TABLE fid_v.keep (id INT8 PRIMARY KEY);
< ok 0 -
> exec
| CREATE FUNCTION fid_v.callee() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 2 $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.caller() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_v.callee() $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.count_old_caller() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_v.count_old() $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.ab (id INT8 PRIMARY KEY, b INT8);
< ok 0 -
> exec
| CREATE FUNCTION fid_v.read_ab() RETURNS INT8 LANGUAGE SQL AS $$ SELECT max(b) FROM fid_v.ab $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.xa() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 7 $$;
< ok 0 -
> exec
| ALTER TABLE fid_v.ab ADD COLUMN x INT8;
< ok 0 -
> exec
| ALTER TABLE fid_v.ab ALTER COLUMN x SET DEFAULT fid_v.xa();
< ok 0 -
> exec
| CREATE FUNCTION fid_v.xb() RETURNS INT8 LANGUAGE SQL AS $$ SELECT max(x) FROM fid_v.ab $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.rt() RETURNS STRING LANGUAGE SQL AS $$ SELECT 'x' $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.pab (id INT8 PRIMARY KEY, b INT8);
< ok 0 -
> exec
| CREATE FUNCTION fid_v.seven() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 7 $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.eight() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 8 $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.grow (id INT8 PRIMARY KEY, s INT8 DEFAULT fid_v.seven(), e INT8 NOT NULL DEFAULT fid_v.eight());
< ok 0 -
> exec
| CREATE FUNCTION fid_v.pf() RETURNS INT8 LANGUAGE SQL AS $$ SELECT COALESCE(max(b), 0) FROM fid_v.pab $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.puse (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_v.pf(), CONSTRAINT puse_v_check CHECK (v >= fid_v.pf() - 1000));
< ok 0 -
> exec
| ALTER TABLE fid_v.items ADD COLUMN doubled DECIMAL AS (price * 2) VIRTUAL;
< ok 0 -
> exec
| CREATE TABLE fid_v.scored (id INT8 PRIMARY KEY, score INT8 DEFAULT fid_v.add(1, 2));
< ok 0 -
> exec
| CREATE VIEW fid_v.cheap AS SELECT p.id FROM fid_v.priced AS p JOIN fid_v.items AS i ON i.id = p.id;
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
| omni_testmigrationfidelity
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
| SELECT nspname, obj_description(oid, 'pg_namespace')
| FROM pg_catalog.pg_namespace
| WHERE nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY nspname;
|
< result
< column "nspname" "NAME"
< column "obj_description" "TEXT"
< row
: string
| fid_v
: null
< row
: string
| public
: null
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
| 	cols.is_hidden,
| 	cols.generation_expression,
| 	(SELECT a.attgenerated FROM pg_catalog.pg_attribute a WHERE a.attrelid = c.oid AND a.attname = cols.column_name) AS attgenerated
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
< column "generation_expression" "TEXT"
< column "attgenerated" "CHAR"
< row
: string
| fid_v
: string
| ab
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
: string
|
: string
|
< row
: string
| fid_v
: string
| ab
: string
| b
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
: string
|
: string
|
< row
: string
| fid_v
: string
| ab
: string
| x
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| fid_v.xa()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| cheap
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
: string
|
: string
|
< row
: string
| fid_v
: string
| grow
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
: string
|
: string
|
< row
: string
| fid_v
: string
| grow
: string
| s
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.seven()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| grow
: string
| e
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| fid_v.eight()
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
: string
|
: string
|
< row
: string
| fid_v
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
: string
|
: string
|
< row
: string
| fid_v
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
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_v
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
: string
|
: string
|
< row
: string
| fid_v
: string
| items
: string
| bonus
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
: string
| fid_v.bonus()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| items
: string
| doubled
: string
| numeric
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
| numeric
: null
: string
| NO
: string
| price * 2
: string
| v
< row
: string
| fid_v
: string
| keep
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
: string
|
: string
|
< row
: string
| fid_v
: string
| named
: string
| other
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
: string
|
: string
|
< row
: string
| fid_v
: string
| pab
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
: string
|
: string
|
< row
: string
| fid_v
: string
| pab
: string
| b
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
: string
|
: string
|
< row
: string
| fid_v
: string
| priced
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
: string
|
: string
|
< row
: string
| fid_v
: string
| priced
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
: string
|
: string
|
< row
: string
| fid_v
: string
| priced_names
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
: string
|
: string
|
< row
: string
| fid_v
: string
| puse
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
: string
|
: string
|
< row
: string
| fid_v
: string
| puse
: string
| v
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.pf()
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
: string
|
: string
|
< row
: string
| fid_v
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
: string
|
: string
|
< row
: string
| fid_v
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
| fid_v.add(1, 2)
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
: string
|
: string
|
< row
: string
| fid_v
: string
| tagged
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
: string
|
: string
|
< row
: string
| fid_v
: string
| tagged
: string
| v
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| unchanged
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
: string
|
: string
|
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
| fid_v
: string
| ab
: string
| ab_pkey
: string
| CREATE UNIQUE INDEX ab_pkey ON omni_testmigrationfidelity.fid_v.ab USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| grow
: string
| grow_pkey
: string
| CREATE UNIQUE INDEX grow_pkey ON omni_testmigrationfidelity.fid_v.grow USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| items
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON omni_testmigrationfidelity.fid_v.items USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| keep
: string
| keep_pkey
: string
| CREATE UNIQUE INDEX keep_pkey ON omni_testmigrationfidelity.fid_v.keep USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| pab
: string
| pab_pkey
: string
| CREATE UNIQUE INDEX pab_pkey ON omni_testmigrationfidelity.fid_v.pab USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| puse
: string
| puse_pkey
: string
| CREATE UNIQUE INDEX puse_pkey ON omni_testmigrationfidelity.fid_v.puse USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| scored
: string
| scored_pkey
: string
| CREATE UNIQUE INDEX scored_pkey ON omni_testmigrationfidelity.fid_v.scored USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| tagged
: string
| tagged_pkey
: string
| CREATE UNIQUE INDEX tagged_pkey ON omni_testmigrationfidelity.fid_v.tagged USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| totals
: string
| totals_pkey
: string
| CREATE UNIQUE INDEX totals_pkey ON omni_testmigrationfidelity.fid_v.totals USING btree (rowid ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| totals
: string
| totals_total_idx
: string
| CREATE INDEX totals_total_idx ON omni_testmigrationfidelity.fid_v.totals USING btree (total ASC)
: null
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
| fid_v
: string
| puse
: string
| puse_v_check
: string
| CHECK ((v >= (fid_v.pf() - 1000)))
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
| fid_v
: string
| ab
: int64 0
: null
< row
: string
| fid_v
: string
| grow
: int64 0
: null
< row
: string
| fid_v
: string
| items
: int64 0
: null
< row
: string
| fid_v
: string
| keep
: int64 0
: null
< row
: string
| fid_v
: string
| pab
: int64 0
: null
< row
: string
| fid_v
: string
| puse
: int64 0
: null
< row
: string
| fid_v
: string
| scored
: int64 0
: null
< row
: string
| fid_v
: string
| tagged
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
| fid_v
: string
| cheap
: string
| SELECT p.id FROM omni_testmigrationfidelity.fid_v.priced AS p JOIN omni_testmigrationfidelity.fid_v.items AS i ON i.id = p.id
: null
< row
: string
| fid_v
: string
| named
: string
| SELECT id FROM omni_testmigrationfidelity.fid_v.items
: null
< row
: string
| fid_v
: string
| priced
: string
| SELECT id, name FROM omni_testmigrationfidelity.fid_v.items WHERE price > 1
: null
< row
: string
| fid_v
: string
| priced_names
: string
| SELECT name FROM omni_testmigrationfidelity.fid_v.priced
: null
< row
: string
| fid_v
: string
| unchanged
: string
| SELECT id FROM omni_testmigrationfidelity.fid_v.items
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
| fid_v
: string
| cheap
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| id
< row
: string
| fid_v
: string
| priced
: string
| id
< row
: string
| fid_v
: string
| priced
: string
| name
> query ^
: string
| fid_v
: string
| named
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| id
> query ^
: string
| fid_v
: string
| priced
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| id
< row
: string
| fid_v
: string
| items
: string
| name
< row
: string
| fid_v
: string
| items
: string
| price
> query ^
: string
| fid_v
: string
| priced_names
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| priced
: string
| id
< row
: string
| fid_v
: string
| priced
: string
| name
> query ^
: string
| fid_v
: string
| unchanged
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
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
| fid_v
: string
| totals
: string
| SELECT name, sum(price) AS total FROM omni_testmigrationfidelity.fid_v.items GROUP BY name
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
| fid_v
: string
| totals
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| name
< row
: string
| fid_v
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
| fid_v
: string
| add
: string
| add(int8, int8)
: string
| CREATE FUNCTION fid_v.add(a INT8, b INT8)
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT (a + b) + 0;
| $$
< row
: string
| fid_v
: string
| bonus
: string
| bonus()
: string
| CREATE FUNCTION fid_v.bonus()
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 5;
| $$
< row
: string
| fid_v
: string
| callee
: string
| callee()
: string
| CREATE FUNCTION fid_v.callee()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 2;
| $$
< row
: string
| fid_v
: string
| caller
: string
| caller()
: string
| CREATE FUNCTION fid_v.caller()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT fid_v.callee();
| $$
< row
: string
| fid_v
: string
| count_items
: string
| count_items()
: string
| CREATE FUNCTION fid_v.count_items()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testmigrationfidelity.fid_v.priced;
| $$
< row
: string
| fid_v
: string
| count_old
: string
| count_old()
: string
| CREATE FUNCTION fid_v.count_old()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testmigrationfidelity.fid_v.items;
| $$
< row
: string
| fid_v
: string
| count_old_caller
: string
| count_old_caller()
: string
| CREATE FUNCTION fid_v.count_old_caller()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT fid_v.count_old();
| $$
< row
: string
| fid_v
: string
| eight
: string
| eight()
: string
| CREATE FUNCTION fid_v.eight()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 8;
| $$
< row
: string
| fid_v
: string
| pf
: string
| pf()
: string
| CREATE FUNCTION fid_v.pf()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT COALESCE(max(b), 0) FROM omni_testmigrationfidelity.fid_v.pab;
| $$
< row
: string
| fid_v
: string
| read_ab
: string
| read_ab()
: string
| CREATE FUNCTION fid_v.read_ab()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT max(b) FROM omni_testmigrationfidelity.fid_v.ab;
| $$
< row
: string
| fid_v
: string
| rt
: string
| rt()
: string
| CREATE FUNCTION fid_v.rt()
| 	RETURNS STRING
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 'x';
| $$
< row
: string
| fid_v
: string
| seven
: string
| seven()
: string
| CREATE FUNCTION fid_v.seven()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 7;
| $$
< row
: string
| fid_v
: string
| touch
: string
| touch()
: string
| CREATE FUNCTION fid_v.touch()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE plpgsql
| 	SECURITY INVOKER
| 	AS $$
| 	BEGIN
| 	RETURN (SELECT count(*) FROM omni_testmigrationfidelity.fid_v.items);
| 	END;
| $$
< row
: string
| fid_v
: string
| xa
: string
| xa()
: string
| CREATE FUNCTION fid_v.xa()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 7;
| $$
< row
: string
| fid_v
: string
| xb
: string
| xb()
: string
| CREATE FUNCTION fid_v.xb()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT max(x) FROM omni_testmigrationfidelity.fid_v.ab;
| $$
> query
|
| 		SELECT s.sequence_schema, s.sequence_name, s.data_type, s.start_value, s.minimum_value, s.maximum_value, s.increment, s.cycle_option,
| 			(SELECT t.relname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_class t ON t.oid = d.refobjid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT a.attname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT cs.create_statement FROM crdb_internal.create_statements cs WHERE cs.descriptor_id = c.oid::INT8)
| 		FROM information_schema.sequences s
| 		LEFT JOIN pg_catalog.pg_namespace n ON n.nspname = s.sequence_schema
| 		LEFT JOIN pg_catalog.pg_class c ON c.relname = s.sequence_name AND c.relnamespace = n.oid AND c.relkind = 'S'
| 		ORDER BY s.sequence_schema, s.sequence_name;
< result
< column "sequence_schema" "TEXT"
< column "sequence_name" "TEXT"
< column "data_type" "TEXT"
< column "start_value" "TEXT"
< column "minimum_value" "TEXT"
< column "maximum_value" "TEXT"
< column "increment" "TEXT"
< column "cycle_option" "TEXT"
< column "relname" "NAME"
< column "attname" "NAME"
< column "create_statement" "TEXT"
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
| DROP SCHEMA IF EXISTS "fid_v" CASCADE;
< ok 0 -
> exec
| CREATE SCHEMA fid_v;
< ok 0 -
> exec
| CREATE TABLE fid_v.items (id INT8 PRIMARY KEY, name STRING, price DECIMAL(10, 2), legacy INT8, bonus INT8);
< ok 0 -
> exec
| CREATE VIEW fid_v.priced AS SELECT id, name, legacy FROM fid_v.items WHERE price > 0;
< ok 0 -
> exec
| CREATE VIEW fid_v.priced_names AS SELECT name FROM fid_v.priced;
< ok 0 -
> exec
| CREATE VIEW fid_v.unchanged AS SELECT id FROM fid_v.items;
< ok 0 -
> exec
| CREATE VIEW fid_v.dropped AS SELECT name FROM fid_v.items;
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW fid_v.totals AS SELECT name, sum(price) AS total FROM fid_v.items GROUP BY name;
< ok 0 -
> exec
| CREATE INDEX totals_name_idx ON fid_v.totals (name);
< ok 0 -
> exec
| CREATE FUNCTION fid_v.add(a INT8, b INT8) RETURNS INT8 LANGUAGE SQL IMMUTABLE AS $$ SELECT a + b $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.gone() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.count_items() RETURNS INT8 LANGUAGE SQL STABLE AS $$ SELECT count(*) FROM fid_v.items $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.tagged (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_v.gone());
< ok 0 -
> exec
| CREATE TABLE fid_v.old_items (id INT8 PRIMARY KEY);
< ok 0 -
> exec
| CREATE FUNCTION fid_v.count_old() RETURNS INT8 LANGUAGE SQL STABLE AS $$ SELECT count(*) FROM fid_v.old_items $$;
< ok 0 -
> exec
| CREATE VIEW fid_v.named (renamed) AS SELECT id FROM fid_v.items;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.ra() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.rb() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_v.ra() $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.c_const() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.ct1 (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_v.c_const());
< ok 0 -
> exec
| CREATE FUNCTION fid_v.cz_read() RETURNS INT8 LANGUAGE SQL AS $$ SELECT count(*) FROM fid_v.ct1 $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.ct2 (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_v.cz_read());
< ok 0 -
> exec
| CREATE FUNCTION fid_v.colfn() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 3 $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.keep (id INT8 PRIMARY KEY, x INT8 DEFAULT fid_v.colfn());
< ok 0 -
> exec
| CREATE FUNCTION fid_v.callee() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.caller() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_v.callee() $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.count_old_caller() RETURNS INT8 LANGUAGE SQL AS $$ SELECT fid_v.count_old() $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.ab (id INT8 PRIMARY KEY, a INT8, b INT8);
< ok 0 -
> exec
| CREATE FUNCTION fid_v.read_ab() RETURNS INT8 LANGUAGE SQL AS $$ SELECT max(a) FROM fid_v.ab $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.rt() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.pab (id INT8 PRIMARY KEY, a INT8, b INT8);
< ok 0 -
> exec
| CREATE TABLE fid_v.grow (id INT8 PRIMARY KEY);
< ok 0 -
> exec
| INSERT INTO fid_v.grow VALUES (1), (2);
< ok 2 -
> exec
| CREATE FUNCTION fid_v.seven() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 7 $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.pf() RETURNS INT8 LANGUAGE SQL AS $$ SELECT COALESCE(max(a), 0) FROM fid_v.pab $$;
< ok 0 -
> exec
| CREATE TABLE fid_v.puse (id INT8 PRIMARY KEY, v INT8 DEFAULT fid_v.pf(), CONSTRAINT puse_v_check CHECK (v >= fid_v.pf() - 1000));
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
| omni_testmigrationfidelity
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
| SELECT nspname, obj_description(oid, 'pg_namespace')
| FROM pg_catalog.pg_namespace
| WHERE nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY nspname;
|
< result
< column "nspname" "NAME"
< column "obj_description" "TEXT"
< row
: string
| fid_v
: null
< row
: string
| public
: null
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
| 	cols.is_hidden,
| 	cols.generation_expression,
| 	(SELECT a.attgenerated FROM pg_catalog.pg_attribute a WHERE a.attrelid = c.oid AND a.attname = cols.column_name) AS attgenerated
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
< column "generation_expression" "TEXT"
< column "attgenerated" "CHAR"
< row
: string
| fid_v
: string
| ab
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
: string
|
: string
|
< row
: string
| fid_v
: string
| ab
: string
| a
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
: string
|
: string
|
< row
: string
| fid_v
: string
| ab
: string
| b
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_v
: string
| ct1
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
: string
|
: string
|
< row
: string
| fid_v
: string
| ct1
: string
| v
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.c_const()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| ct2
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
: string
|
: string
|
< row
: string
| fid_v
: string
| ct2
: string
| v
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.cz_read()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| dropped
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
: string
|
: string
|
< row
: string
| fid_v
: string
| grow
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
: string
|
: string
|
< row
: string
| fid_v
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
: string
|
: string
|
< row
: string
| fid_v
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
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_v
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
: string
|
: string
|
< row
: string
| fid_v
: string
| items
: string
| legacy
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
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
: string
|
: string
|
< row
: string
| fid_v
: string
| items
: string
| bonus
: string
| bigint
: null
: int64 64
: int64 0
: int64 5
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
: string
|
: string
|
< row
: string
| fid_v
: string
| keep
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
: string
|
: string
|
< row
: string
| fid_v
: string
| keep
: string
| x
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.colfn()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| named
: string
| renamed
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
: string
|
: string
|
< row
: string
| fid_v
: string
| old_items
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
: string
|
: string
|
< row
: string
| fid_v
: string
| pab
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
: string
|
: string
|
< row
: string
| fid_v
: string
| pab
: string
| a
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
: string
|
: string
|
< row
: string
| fid_v
: string
| pab
: string
| b
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_v
: string
| priced
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
: string
|
: string
|
< row
: string
| fid_v
: string
| priced
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
: string
|
: string
|
< row
: string
| fid_v
: string
| priced
: string
| legacy
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_v
: string
| priced_names
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
: string
|
: string
|
< row
: string
| fid_v
: string
| puse
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
: string
|
: string
|
< row
: string
| fid_v
: string
| puse
: string
| v
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.pf()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| tagged
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
: string
|
: string
|
< row
: string
| fid_v
: string
| tagged
: string
| v
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.gone()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| unchanged
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
: string
|
: string
|
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
| fid_v
: string
| ab
: string
| ab_pkey
: string
| CREATE UNIQUE INDEX ab_pkey ON omni_testmigrationfidelity.fid_v.ab USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| ct1
: string
| ct1_pkey
: string
| CREATE UNIQUE INDEX ct1_pkey ON omni_testmigrationfidelity.fid_v.ct1 USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| ct2
: string
| ct2_pkey
: string
| CREATE UNIQUE INDEX ct2_pkey ON omni_testmigrationfidelity.fid_v.ct2 USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| grow
: string
| grow_pkey
: string
| CREATE UNIQUE INDEX grow_pkey ON omni_testmigrationfidelity.fid_v.grow USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| items
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON omni_testmigrationfidelity.fid_v.items USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| keep
: string
| keep_pkey
: string
| CREATE UNIQUE INDEX keep_pkey ON omni_testmigrationfidelity.fid_v.keep USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| old_items
: string
| old_items_pkey
: string
| CREATE UNIQUE INDEX old_items_pkey ON omni_testmigrationfidelity.fid_v.old_items USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| pab
: string
| pab_pkey
: string
| CREATE UNIQUE INDEX pab_pkey ON omni_testmigrationfidelity.fid_v.pab USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| puse
: string
| puse_pkey
: string
| CREATE UNIQUE INDEX puse_pkey ON omni_testmigrationfidelity.fid_v.puse USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| tagged
: string
| tagged_pkey
: string
| CREATE UNIQUE INDEX tagged_pkey ON omni_testmigrationfidelity.fid_v.tagged USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| totals
: string
| totals_name_idx
: string
| CREATE INDEX totals_name_idx ON omni_testmigrationfidelity.fid_v.totals USING btree (name ASC)
: null
: null
< row
: string
| fid_v
: string
| totals
: string
| totals_pkey
: string
| CREATE UNIQUE INDEX totals_pkey ON omni_testmigrationfidelity.fid_v.totals USING btree (rowid ASC)
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
| fid_v
: string
| puse
: string
| puse_v_check
: string
| CHECK ((v >= (fid_v.pf() - 1000)))
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
| fid_v
: string
| ab
: int64 0
: null
< row
: string
| fid_v
: string
| ct1
: int64 0
: null
< row
: string
| fid_v
: string
| ct2
: int64 0
: null
< row
: string
| fid_v
: string
| grow
: int64 0
: null
< row
: string
| fid_v
: string
| items
: int64 0
: null
< row
: string
| fid_v
: string
| keep
: int64 0
: null
< row
: string
| fid_v
: string
| old_items
: int64 0
: null
< row
: string
| fid_v
: string
| pab
: int64 0
: null
< row
: string
| fid_v
: string
| puse
: int64 0
: null
< row
: string
| fid_v
: string
| tagged
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
| fid_v
: string
| dropped
: string
| SELECT name FROM omni_testmigrationfidelity.fid_v.items
: null
< row
: string
| fid_v
: string
| named
: string
| SELECT id FROM omni_testmigrationfidelity.fid_v.items
: null
< row
: string
| fid_v
: string
| priced
: string
| SELECT id, name, legacy FROM omni_testmigrationfidelity.fid_v.items WHERE price > 0
: null
< row
: string
| fid_v
: string
| priced_names
: string
| SELECT name FROM omni_testmigrationfidelity.fid_v.priced
: null
< row
: string
| fid_v
: string
| unchanged
: string
| SELECT id FROM omni_testmigrationfidelity.fid_v.items
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
| fid_v
: string
| dropped
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| name
> query ^
: string
| fid_v
: string
| named
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| id
> query ^
: string
| fid_v
: string
| priced
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| id
< row
: string
| fid_v
: string
| items
: string
| legacy
< row
: string
| fid_v
: string
| items
: string
| name
< row
: string
| fid_v
: string
| items
: string
| price
> query ^
: string
| fid_v
: string
| priced_names
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| priced
: string
| id
< row
: string
| fid_v
: string
| priced
: string
| legacy
< row
: string
| fid_v
: string
| priced
: string
| name
> query ^
: string
| fid_v
: string
| unchanged
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
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
| fid_v
: string
| totals
: string
| SELECT name, sum(price) AS total FROM omni_testmigrationfidelity.fid_v.items GROUP BY name
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
| fid_v
: string
| totals
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| name
< row
: string
| fid_v
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
| fid_v
: string
| add
: string
| add(int8, int8)
: string
| CREATE FUNCTION fid_v.add(a INT8, b INT8)
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
| fid_v
: string
| c_const
: string
| c_const()
: string
| CREATE FUNCTION fid_v.c_const()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 1;
| $$
< row
: string
| fid_v
: string
| callee
: string
| callee()
: string
| CREATE FUNCTION fid_v.callee()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 1;
| $$
< row
: string
| fid_v
: string
| caller
: string
| caller()
: string
| CREATE FUNCTION fid_v.caller()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT fid_v.callee();
| $$
< row
: string
| fid_v
: string
| colfn
: string
| colfn()
: string
| CREATE FUNCTION fid_v.colfn()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 3;
| $$
< row
: string
| fid_v
: string
| count_items
: string
| count_items()
: string
| CREATE FUNCTION fid_v.count_items()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testmigrationfidelity.fid_v.items;
| $$
< row
: string
| fid_v
: string
| count_old
: string
| count_old()
: string
| CREATE FUNCTION fid_v.count_old()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testmigrationfidelity.fid_v.old_items;
| $$
< row
: string
| fid_v
: string
| count_old_caller
: string
| count_old_caller()
: string
| CREATE FUNCTION fid_v.count_old_caller()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT fid_v.count_old();
| $$
< row
: string
| fid_v
: string
| cz_read
: string
| cz_read()
: string
| CREATE FUNCTION fid_v.cz_read()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testmigrationfidelity.fid_v.ct1;
| $$
< row
: string
| fid_v
: string
| gone
: string
| gone()
: string
| CREATE FUNCTION fid_v.gone()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 1;
| $$
< row
: string
| fid_v
: string
| pf
: string
| pf()
: string
| CREATE FUNCTION fid_v.pf()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT COALESCE(max(a), 0) FROM omni_testmigrationfidelity.fid_v.pab;
| $$
< row
: string
| fid_v
: string
| ra
: string
| ra()
: string
| CREATE FUNCTION fid_v.ra()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 1;
| $$
< row
: string
| fid_v
: string
| rb
: string
| rb()
: string
| CREATE FUNCTION fid_v.rb()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT fid_v.ra();
| $$
< row
: string
| fid_v
: string
| read_ab
: string
| read_ab()
: string
| CREATE FUNCTION fid_v.read_ab()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT max(a) FROM omni_testmigrationfidelity.fid_v.ab;
| $$
< row
: string
| fid_v
: string
| rt
: string
| rt()
: string
| CREATE FUNCTION fid_v.rt()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 1;
| $$
< row
: string
| fid_v
: string
| seven
: string
| seven()
: string
| CREATE FUNCTION fid_v.seven()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 7;
| $$
> query
|
| 		SELECT s.sequence_schema, s.sequence_name, s.data_type, s.start_value, s.minimum_value, s.maximum_value, s.increment, s.cycle_option,
| 			(SELECT t.relname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_class t ON t.oid = d.refobjid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT a.attname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT cs.create_statement FROM crdb_internal.create_statements cs WHERE cs.descriptor_id = c.oid::INT8)
| 		FROM information_schema.sequences s
| 		LEFT JOIN pg_catalog.pg_namespace n ON n.nspname = s.sequence_schema
| 		LEFT JOIN pg_catalog.pg_class c ON c.relname = s.sequence_name AND c.relnamespace = n.oid AND c.relkind = 'S'
| 		ORDER BY s.sequence_schema, s.sequence_name;
< result
< column "sequence_schema" "TEXT"
< column "sequence_name" "TEXT"
< column "data_type" "TEXT"
< column "start_value" "TEXT"
< column "minimum_value" "TEXT"
< column "maximum_value" "TEXT"
< column "increment" "TEXT"
< column "cycle_option" "TEXT"
< column "relname" "NAME"
< column "attname" "NAME"
< column "create_statement" "TEXT"
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
| ALTER TABLE "fid_v"."keep" ALTER COLUMN "x" DROP DEFAULT;
< ok 0 -
> exec
| ALTER TABLE "fid_v"."puse" ALTER COLUMN "v" DROP DEFAULT;
< ok 0 -
> exec
| ALTER TABLE "fid_v"."tagged" ALTER COLUMN "v" DROP DEFAULT;
< ok 0 -
> exec
| ALTER TABLE "fid_v"."puse" DROP CONSTRAINT "puse_v_check";
< ok 0 -
> exec
| DROP INDEX "fid_v"."totals"@"totals_name_idx";
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."colfn"();
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."count_old_caller"();
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."count_old"();
< ok 0 -
> exec
| DROP TABLE "fid_v"."ct2";
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."cz_read"();
< ok 0 -
> exec
| DROP TABLE "fid_v"."ct1";
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."c_const"();
< ok 0 -
> exec
| DROP VIEW "fid_v"."dropped";
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."gone"();
< ok 0 -
> exec
| DROP VIEW "fid_v"."named";
< ok 0 -
> exec
| DROP TABLE "fid_v"."old_items";
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."pf"();
< ok 0 -
> exec
| DROP VIEW "fid_v"."priced_names";
< ok 0 -
> exec
| DROP VIEW "fid_v"."priced";
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."rb"();
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."ra"();
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."read_ab"();
< ok 0 -
> exec
| DROP FUNCTION "fid_v"."rt"();
< ok 0 -
> exec
| ALTER TABLE "fid_v"."ab" DROP COLUMN "a";
< ok 0 -
> exec
| ALTER TABLE "fid_v"."items" DROP COLUMN "legacy";
< ok 0 -
> exec
| ALTER TABLE "fid_v"."keep" DROP COLUMN "x";
< ok 0 -
> exec
| ALTER TABLE "fid_v"."pab" DROP COLUMN "a";
< ok 0 -
> exec
| ALTER TABLE "fid_v"."ab" ADD COLUMN "x" bigint;
< ok 0 -
> exec
| ALTER TABLE "fid_v"."grow" ADD COLUMN "s" bigint;
< ok 0 -
> exec
| ALTER TABLE "fid_v"."grow" ALTER COLUMN "s" SET DEFAULT fid_v.seven();
< ok 0 -
> exec
| ALTER TABLE "fid_v"."grow" ADD COLUMN "e" bigint;
< ok 0 -
> exec
| ALTER TABLE "fid_v"."items" ADD COLUMN "doubled" numeric AS (price * 2) VIRTUAL;
< ok 0 -
> exec
| CREATE OR REPLACE FUNCTION fid_v.add(a INT8, b INT8)
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT (a + b) + 0;
| $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.bonus()
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 5;
| $$;
< ok 0 -
> exec
| CREATE OR REPLACE FUNCTION fid_v.callee()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 2;
| $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.eight()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 8;
| $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.count_old()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| SELECT count(*) FROM fid_v.items;
| $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.count_old_caller()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT fid_v.count_old();
| $$;
< ok 0 -
> exec
| CREATE VIEW "fid_v"."named" ("other") AS 
| SELECT id FROM fid_v.items;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.pf()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| SELECT COALESCE(max(b), 0) FROM fid_v.pab;
| $$;
< ok 0 -
> exec
| CREATE VIEW "fid_v"."priced" ("id", "name") AS 
| SELECT id, name FROM fid_v.items WHERE price > 1;
< ok 0 -
> exec
| CREATE VIEW "fid_v"."cheap" ("id") AS 
| SELECT p.id FROM fid_v.priced AS p JOIN fid_v.items AS i ON i.id = p.id;
< ok 0 -
> exec
| CREATE OR REPLACE FUNCTION fid_v.count_items()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| SELECT count(*) FROM fid_v.priced;
| $$;
< ok 0 -
> exec
| CREATE VIEW "fid_v"."priced_names" ("name") AS 
| SELECT name FROM fid_v.priced;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.read_ab()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| SELECT max(b) FROM fid_v.ab;
| $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.rt()
| 	RETURNS STRING
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 'x';
| $$;
< ok 0 -
> exec
| CREATE TABLE "fid_v"."scored" (
|     "id" bigint NOT NULL,
|     "score" bigint DEFAULT fid_v.add(1, 2),
|     CONSTRAINT "scored_pkey" PRIMARY KEY (id ASC)
| );
< ok 0 -
> exec
| CREATE FUNCTION fid_v.touch()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE plpgsql
| 	SECURITY INVOKER
| 	AS $$
| BEGIN
| RETURN (SELECT count(*) FROM fid_v.items);
| END;
|
| $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.xa()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 7;
| $$;
< ok 0 -
> exec
| CREATE FUNCTION fid_v.xb()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| SELECT max(x) FROM fid_v.ab;
| $$;
< ok 0 -
> exec
| ALTER TABLE "fid_v"."ab" ALTER COLUMN "x" SET DEFAULT fid_v.xa();
< ok 0 -
> exec
| ALTER TABLE "fid_v"."grow" ALTER COLUMN "e" SET DEFAULT fid_v.eight();
< ok 0 -
> exec
| UPDATE "fid_v"."grow" SET "e" = fid_v.eight() WHERE "e" IS NULL;
< ok 2 -
> exec
| ALTER TABLE "fid_v"."grow" ALTER COLUMN "e" SET NOT NULL;
< ok 0 -
> exec
| ALTER TABLE "fid_v"."items" ALTER COLUMN "bonus" SET DEFAULT fid_v.bonus();
< ok 0 -
> exec
| ALTER TABLE "fid_v"."puse" ALTER COLUMN "v" SET DEFAULT fid_v.pf();
< ok 0 -
> exec
| ALTER TABLE "fid_v"."puse" ADD CONSTRAINT "puse_v_check" CHECK ((v >= (fid_v.pf() - 1000)));
< ok 0 -
> exec
| CREATE INDEX totals_total_idx ON fid_v.totals (total ASC);
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
| omni_testmigrationfidelity
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
| SELECT nspname, obj_description(oid, 'pg_namespace')
| FROM pg_catalog.pg_namespace
| WHERE nspname NOT IN ('crdb_internal','information_schema','pg_catalog','pg_extension')
| ORDER BY nspname;
|
< result
< column "nspname" "NAME"
< column "obj_description" "TEXT"
< row
: string
| fid_v
: null
< row
: string
| public
: null
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
| 	cols.is_hidden,
| 	cols.generation_expression,
| 	(SELECT a.attgenerated FROM pg_catalog.pg_attribute a WHERE a.attrelid = c.oid AND a.attname = cols.column_name) AS attgenerated
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
< column "generation_expression" "TEXT"
< column "attgenerated" "CHAR"
< row
: string
| fid_v
: string
| ab
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
: string
|
: string
|
< row
: string
| fid_v
: string
| ab
: string
| b
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_v
: string
| ab
: string
| x
: string
| bigint
: null
: int64 64
: int64 0
: int64 4
: string
| fid_v.xa()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| cheap
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
: string
|
: string
|
< row
: string
| fid_v
: string
| grow
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
: string
|
: string
|
< row
: string
| fid_v
: string
| grow
: string
| s
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.seven()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| grow
: string
| e
: string
| bigint
: null
: int64 64
: int64 0
: int64 3
: string
| fid_v.eight()
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
: string
|
: string
|
< row
: string
| fid_v
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
: string
|
: string
|
< row
: string
| fid_v
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
| YES
: null
: string
| pg_catalog
: string
| text
: null
: string
| NO
: string
|
: string
|
< row
: string
| fid_v
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
: string
|
: string
|
< row
: string
| fid_v
: string
| items
: string
| bonus
: string
| bigint
: null
: int64 64
: int64 0
: int64 5
: string
| fid_v.bonus()
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
: string
|
: string
|
< row
: string
| fid_v
: string
| items
: string
| doubled
: string
| numeric
: null
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
| numeric
: null
: string
| NO
: string
| price * 2
: string
| v
< row
: string
| fid_v
: string
| keep
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
: string
|
: string
|
< row
: string
| fid_v
: string
| named
: string
| other
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
: string
|
: string
|
< row
: string
| fid_v
: string
| pab
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
: string
|
: string
|
< row
: string
| fid_v
: string
| pab
: string
| b
: string
| bigint
: null
: int64 64
: int64 0
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
: string
| NO
: string
|
: string
|
< row
: string
| fid_v
: string
| priced
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
: string
|
: string
|
< row
: string
| fid_v
: string
| priced
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
: string
|
: string
|
< row
: string
| fid_v
: string
| priced_names
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
: string
|
: string
|
< row
: string
| fid_v
: string
| puse
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
: string
|
: string
|
< row
: string
| fid_v
: string
| puse
: string
| v
: string
| bigint
: null
: int64 64
: int64 0
: int64 2
: string
| fid_v.pf()
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
: string
|
: string
|
< row
: string
| fid_v
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
: string
|
: string
|
< row
: string
| fid_v
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
| fid_v.add(1, 2)
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
: string
|
: string
|
< row
: string
| fid_v
: string
| tagged
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
: string
|
: string
|
< row
: string
| fid_v
: string
| tagged
: string
| v
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| totals
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
: string
|
: string
|
< row
: string
| fid_v
: string
| unchanged
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
: string
|
: string
|
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
| fid_v
: string
| ab
: string
| ab_pkey
: string
| CREATE UNIQUE INDEX ab_pkey ON omni_testmigrationfidelity.fid_v.ab USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| grow
: string
| grow_pkey
: string
| CREATE UNIQUE INDEX grow_pkey ON omni_testmigrationfidelity.fid_v.grow USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| items
: string
| items_pkey
: string
| CREATE UNIQUE INDEX items_pkey ON omni_testmigrationfidelity.fid_v.items USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| keep
: string
| keep_pkey
: string
| CREATE UNIQUE INDEX keep_pkey ON omni_testmigrationfidelity.fid_v.keep USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| pab
: string
| pab_pkey
: string
| CREATE UNIQUE INDEX pab_pkey ON omni_testmigrationfidelity.fid_v.pab USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| puse
: string
| puse_pkey
: string
| CREATE UNIQUE INDEX puse_pkey ON omni_testmigrationfidelity.fid_v.puse USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| scored
: string
| scored_pkey
: string
| CREATE UNIQUE INDEX scored_pkey ON omni_testmigrationfidelity.fid_v.scored USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| tagged
: string
| tagged_pkey
: string
| CREATE UNIQUE INDEX tagged_pkey ON omni_testmigrationfidelity.fid_v.tagged USING btree (id ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| totals
: string
| totals_pkey
: string
| CREATE UNIQUE INDEX totals_pkey ON omni_testmigrationfidelity.fid_v.totals USING btree (rowid ASC)
: string
| PRIMARY KEY
: null
< row
: string
| fid_v
: string
| totals
: string
| totals_total_idx
: string
| CREATE INDEX totals_total_idx ON omni_testmigrationfidelity.fid_v.totals USING btree (total ASC)
: null
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
| fid_v
: string
| puse
: string
| puse_v_check
: string
| CHECK ((v >= (fid_v.pf() - 1000)))
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
| fid_v
: string
| ab
: int64 0
: null
< row
: string
| fid_v
: string
| grow
: int64 0
: null
< row
: string
| fid_v
: string
| items
: int64 0
: null
< row
: string
| fid_v
: string
| keep
: int64 0
: null
< row
: string
| fid_v
: string
| pab
: int64 0
: null
< row
: string
| fid_v
: string
| puse
: int64 0
: null
< row
: string
| fid_v
: string
| scored
: int64 0
: null
< row
: string
| fid_v
: string
| tagged
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
| fid_v
: string
| cheap
: string
| SELECT p.id FROM omni_testmigrationfidelity.fid_v.priced AS p JOIN omni_testmigrationfidelity.fid_v.items AS i ON i.id = p.id
: null
< row
: string
| fid_v
: string
| named
: string
| SELECT id FROM omni_testmigrationfidelity.fid_v.items
: null
< row
: string
| fid_v
: string
| priced
: string
| SELECT id, name FROM omni_testmigrationfidelity.fid_v.items WHERE price > 1
: null
< row
: string
| fid_v
: string
| priced_names
: string
| SELECT name FROM omni_testmigrationfidelity.fid_v.priced
: null
< row
: string
| fid_v
: string
| unchanged
: string
| SELECT id FROM omni_testmigrationfidelity.fid_v.items
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
| fid_v
: string
| cheap
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| id
< row
: string
| fid_v
: string
| priced
: string
| id
< row
: string
| fid_v
: string
| priced
: string
| name
> query ^
: string
| fid_v
: string
| named
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| id
> query ^
: string
| fid_v
: string
| priced
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| id
< row
: string
| fid_v
: string
| items
: string
| name
< row
: string
| fid_v
: string
| items
: string
| price
> query ^
: string
| fid_v
: string
| priced_names
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| priced
: string
| id
< row
: string
| fid_v
: string
| priced
: string
| name
> query ^
: string
| fid_v
: string
| unchanged
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
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
| fid_v
: string
| totals
: string
| SELECT name, sum(price) AS total FROM omni_testmigrationfidelity.fid_v.items GROUP BY name
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
| fid_v
: string
| totals
< result
< column "source_schema" "NAME"
< column "source_table" "NAME"
< column "column_name" "NAME"
< row
: string
| fid_v
: string
| items
: string
| name
< row
: string
| fid_v
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
| fid_v
: string
| add
: string
| add(int8, int8)
: string
| CREATE FUNCTION fid_v.add(a INT8, b INT8)
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT (a + b) + 0;
| $$
< row
: string
| fid_v
: string
| bonus
: string
| bonus()
: string
| CREATE FUNCTION fid_v.bonus()
| 	RETURNS INT8
| 	IMMUTABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 5;
| $$
< row
: string
| fid_v
: string
| callee
: string
| callee()
: string
| CREATE FUNCTION fid_v.callee()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 2;
| $$
< row
: string
| fid_v
: string
| caller
: string
| caller()
: string
| CREATE FUNCTION fid_v.caller()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT fid_v.callee();
| $$
< row
: string
| fid_v
: string
| count_items
: string
| count_items()
: string
| CREATE FUNCTION fid_v.count_items()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testmigrationfidelity.fid_v.priced;
| $$
< row
: string
| fid_v
: string
| count_old
: string
| count_old()
: string
| CREATE FUNCTION fid_v.count_old()
| 	RETURNS INT8
| 	STABLE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT count(*) FROM omni_testmigrationfidelity.fid_v.items;
| $$
< row
: string
| fid_v
: string
| count_old_caller
: string
| count_old_caller()
: string
| CREATE FUNCTION fid_v.count_old_caller()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT fid_v.count_old();
| $$
< row
: string
| fid_v
: string
| eight
: string
| eight()
: string
| CREATE FUNCTION fid_v.eight()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 8;
| $$
< row
: string
| fid_v
: string
| pf
: string
| pf()
: string
| CREATE FUNCTION fid_v.pf()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT COALESCE(max(b), 0) FROM omni_testmigrationfidelity.fid_v.pab;
| $$
< row
: string
| fid_v
: string
| read_ab
: string
| read_ab()
: string
| CREATE FUNCTION fid_v.read_ab()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT max(b) FROM omni_testmigrationfidelity.fid_v.ab;
| $$
< row
: string
| fid_v
: string
| rt
: string
| rt()
: string
| CREATE FUNCTION fid_v.rt()
| 	RETURNS STRING
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 'x';
| $$
< row
: string
| fid_v
: string
| seven
: string
| seven()
: string
| CREATE FUNCTION fid_v.seven()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 7;
| $$
< row
: string
| fid_v
: string
| touch
: string
| touch()
: string
| CREATE FUNCTION fid_v.touch()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE plpgsql
| 	SECURITY INVOKER
| 	AS $$
| 	BEGIN
| 	RETURN (SELECT count(*) FROM omni_testmigrationfidelity.fid_v.items);
| 	END;
| $$
< row
: string
| fid_v
: string
| xa
: string
| xa()
: string
| CREATE FUNCTION fid_v.xa()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT 7;
| $$
< row
: string
| fid_v
: string
| xb
: string
| xb()
: string
| CREATE FUNCTION fid_v.xb()
| 	RETURNS INT8
| 	VOLATILE
| 	NOT LEAKPROOF
| 	CALLED ON NULL INPUT
| 	LANGUAGE SQL
| 	SECURITY INVOKER
| 	AS $$
| 	SELECT max(x) FROM omni_testmigrationfidelity.fid_v.ab;
| $$
> query
|
| 		SELECT s.sequence_schema, s.sequence_name, s.data_type, s.start_value, s.minimum_value, s.maximum_value, s.increment, s.cycle_option,
| 			(SELECT t.relname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_class t ON t.oid = d.refobjid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT a.attname FROM pg_catalog.pg_depend d JOIN pg_catalog.pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
| 				WHERE d.objid = c.oid AND d.deptype = 'a' AND d.refobjsubid > 0 LIMIT 1),
| 			(SELECT cs.create_statement FROM crdb_internal.create_statements cs WHERE cs.descriptor_id = c.oid::INT8)
| 		FROM information_schema.sequences s
| 		LEFT JOIN pg_catalog.pg_namespace n ON n.nspname = s.sequence_schema
| 		LEFT JOIN pg_catalog.pg_class c ON c.relname = s.sequence_name AND c.relnamespace = n.oid AND c.relkind = 'S'
| 		ORDER BY s.sequence_schema, s.sequence_name;
< result
< column "sequence_schema" "TEXT"
< column "sequence_name" "TEXT"
< column "data_type" "TEXT"
< column "start_value" "TEXT"
< column "minimum_value" "TEXT"
< column "maximum_value" "TEXT"
< column "increment" "TEXT"
< column "cycle_option" "TEXT"
< column "relname" "NAME"
< column "attname" "NAME"
< column "create_statement" "TEXT"
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
| DROP SCHEMA IF EXISTS "fid_v" CASCADE;
< ok 0 -
> exec ^
< ok 0 -
