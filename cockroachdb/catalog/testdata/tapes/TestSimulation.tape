-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestSimulation default#1
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SCHEMA s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SCHEMA IF NOT EXISTS s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SCHEMA s2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s2;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SCHEMA "Mixed Case"
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA "Mixed Case";
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SCHEMA s RENAME TO s2
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SCHEMA s RENAME TO public
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SCHEMA missing RENAME TO s2
< error
| 3F000 unknown schema "missing"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SCHEMA s
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SCHEMA s CASCADE
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SCHEMA missing
< error
| 3F000 unknown schema "missing"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SCHEMA public
< error
| 3F000 cannot drop schema "public"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON SCHEMA s IS 'a schema'
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON SCHEMA public IS 'the default'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
| COMMENT ON SCHEMA public IS 'the default';
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE t (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS t (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, v STRING NOT NULL)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	v STRING NOT NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (x INT8, y STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NULL,
| 	y STRING NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT n_pkey PRIMARY KEY (rowid ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (x INT8 NOT NULL, y INT8 NOT NULL, PRIMARY KEY (x, y DESC))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NOT NULL,
| 	y INT8 NOT NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (x ASC, y DESC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8, CONSTRAINT n_key PRIMARY KEY (id))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	CONSTRAINT n_key PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 DEFAULT 7, b STRING DEFAULT 'x', c TIMESTAMPTZ DEFAULT now())
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL DEFAULT 7:::INT8,
| 	b STRING NULL DEFAULT 'x':::STRING,
| 	c TIMESTAMPTZ NULL DEFAULT now():::TIMESTAMPTZ,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 CHECK (a > 0), b INT8, CHECK (b < 10), CONSTRAINT n_ab CHECK (a < b))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL,
| 	b INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT check_b CHECK (b < 10:::INT8),
| 	CONSTRAINT n_ab CHECK (a < b),
| 	CONSTRAINT check_a CHECK (a > 0:::INT8)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, e STRING UNIQUE, f STRING, UNIQUE (f))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	e STRING NULL,
| 	f STRING NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX n_e_key (e ASC),
| 	UNIQUE INDEX n_f_key (f ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8, INDEX (a), INDEX n_named (a DESC))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	INDEX n_a_idx (a ASC),
| 	INDEX n_named (a DESC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, price DECIMAL(10, 2), cents INT8 AS ((price * 100)::INT8) STORED, twice DECIMAL AS (price * 2) VIRTUAL)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	price DECIMAL(10,2) NULL,
| 	cents INT8 NULL AS ((price * 100:::DECIMAL)::INT8) STORED,
| 	twice DECIMAL NULL AS (price * 2:::DECIMAL) VIRTUAL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, secret STRING NOT VISIBLE, tags STRING[], doc JSONB, name VARCHAR(20) COLLATE "en_US")
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	secret STRING NOT VISIBLE NULL,
| 	tags STRING[] NULL,
| 	doc JSONB NULL,
| 	name VARCHAR(20) COLLATE en_US NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id SERIAL PRIMARY KEY, v STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL DEFAULT unique_rowid(),
| 	v STRING NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES t (id))
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES t)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8, CONSTRAINT n_t FOREIGN KEY (t_id) REFERENCES t (id) ON DELETE CASCADE ON UPDATE SET NULL)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES missing (id))
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, k_x INT8 REFERENCES k (x))
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, m mood DEFAULT 'sad', p pair)
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, m missing_type)
< error
| 42704 type "missing_type" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, v INT8 DEFAULT nextval('sq'))
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, id INT8)
< error
| 42701 relation "n" (205): duplicate column name: "id"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 CHECK (missing > 0))
< error
| 42703 column "missing" does not exist, referenced in "missing > 0"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, INDEX (missing))
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY) WITH (ttl_expire_after = '1 day')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	crdb_internal_expiration TIMESTAMPTZ NOT VISIBLE NOT NULL DEFAULT current_timestamp():::TIMESTAMPTZ + '1 day':::INTERVAL ON UPDATE current_timestamp():::TIMESTAMPTZ + '1 day':::INTERVAL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| ) WITH (ttl = 'on', ttl_expire_after = '1 day':::INTERVAL);
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE s.n (id INT8 PRIMARY KEY)
< error
| 3F000 cannot create "s.n" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE s.u (id INT8 PRIMARY KEY)
< error
| 3F000 cannot create "s.u" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE missing.n (id INT8 PRIMARY KEY)
< error
| 3F000 cannot create "missing.n" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE v (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.v (
| 	id INT8 NOT NULL,
| 	CONSTRAINT v_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE sq (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.sq (
| 	id INT8 NOT NULL,
| 	CONSTRAINT sq_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE "Mixed" ("Id" INT8 PRIMARY KEY, "Some Value" STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public."Mixed" (
| 	"Id" INT8 NOT NULL,
| 	"Some Value" STRING NULL,
| 	CONSTRAINT "Mixed_pkey" PRIMARY KEY ("Id" ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TABLE n AS SELECT 1 AS x
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT n_pkey PRIMARY KEY (rowid ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TEMP TABLE n (id INT8 PRIMARY KEY)
< error
| XCEXF temporary tables are only supported experimentally
| hint: See: https://go.crdb.dev/issue-v/46260/v25.2\n--\nYou can enable temporary tables by running `SET experimental_enable_temp_tables = 'on'`.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME TO t2
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME TO c
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE IF EXISTS t RENAME TO t2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE missing RENAME TO t2
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE IF EXISTS missing RENAME TO t2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t SET SCHEMA s
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE k SET SCHEMA s
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE s.u RENAME TO u2
< error
| 42P01 relation "s.u" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE t CASCADE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE IF EXISTS t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE c
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE c, t
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE k
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE s.u
< error
| 42P01 relation "s.u" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TABLE v
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| TRUNCATE t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON TABLE t IS 'changed'
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON TABLE t IS NULL
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON TABLE missing IS 'x'
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD n INT8
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT NULL DEFAULT 5
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT NULL
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n STRING DEFAULT 'x' CHECK (n <> '')
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n STRING UNIQUE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 REFERENCES t (id)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 AS (a * 2) STORED
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 AS (a * 2) VIRTUAL
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n mood
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT VISIBLE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 DEFAULT nextval('sq')
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN a INT8
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN IF NOT EXISTS a INT8
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN IF NOT EXISTS n INT8
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8, ADD COLUMN o STRING
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE missing ADD COLUMN n INT8
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE k ADD COLUMN n INT8
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE v ADD COLUMN n INT8
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP b
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a CASCADE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN m
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN missing
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS missing
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS b
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN t_id
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN note
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE k DROP COLUMN y
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 'x'
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT nextval('missing')
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a DROP DEFAULT
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b DROP DEFAULT
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing SET DEFAULT 1
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET NOT NULL
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a DROP NOT NULL
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN id DROP NOT NULL
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing SET NOT NULL
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b TYPE STRING(100)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET DATA TYPE VARCHAR(50)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a TYPE INT4
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a TYPE STRING
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing TYPE STRING
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET NOT VISIBLE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET VISIBLE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN b TO z
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME b TO z
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN a TO z
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN id TO z
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN b TO a
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN missing TO z
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE c RENAME COLUMN t_id TO parent
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON COLUMN t.b IS 'bee'
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON COLUMN t.a IS NULL
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON COLUMN t.missing IS 'x'
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_pos CHECK (a > -1)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CHECK (a < 1000)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_pos CHECK (a > -1) NOT VALID
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_a_check CHECK (a > 1)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad CHECK (missing > 0)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_b_unique UNIQUE (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD UNIQUE (a, b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_b_key UNIQUE (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad UNIQUE (missing)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_self FOREIGN KEY (a) REFERENCES t (id)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES t (id) ON DELETE CASCADE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad FOREIGN KEY (a) REFERENCES missing (id)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad FOREIGN KEY (a) REFERENCES t (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE c ADD CONSTRAINT c_note CHECK (note <> '')
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE k ADD PRIMARY KEY (x)
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ADD PRIMARY KEY (a)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_a_check
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT IF EXISTS t_a_check
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT IF EXISTS missing
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT missing
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_b_key
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_pkey
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE c DROP CONSTRAINT c_t_id_fkey
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_a_check TO t_a_nonneg
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_b_key TO t_b_uniq
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_pkey TO t_key
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT missing TO z
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t VALIDATE CONSTRAINT t_a_check
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE c VALIDATE CONSTRAINT c_t_id_fkey
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t ALTER PRIMARY KEY USING COLUMNS (a)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE k ALTER PRIMARY KEY USING COLUMNS (x)
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON CONSTRAINT t_a_check ON t IS 'non-negative'
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON CONSTRAINT t_pkey ON t IS 'the key'
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON CONSTRAINT missing ON t IS 'x'
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX ON t (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_b_idx ON t (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_ab ON t (a DESC, b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE UNIQUE INDEX t_b_uniq ON t (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_b_store ON t (b) STORING (a)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_b_part ON t (b) WHERE a > 0
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_lower ON t (lower(b))
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_hash ON t (b) USING HASH
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_hidden ON t (b) NOT VISIBLE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INVERTED INDEX t_doc ON t (doc)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_a_idx ON t (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX IF NOT EXISTS t_a_idx ON t (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX IF NOT EXISTS t_b_idx ON t (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX t_bad ON t (missing)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX n_idx ON missing (a)
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX ON k (x)
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX ON s.u (v)
< error
| 42P01 relation "s.u" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX mv_b ON mv (b)
< error
| 42P01 relation "mv" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX v_a ON v (a)
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE INDEX CONCURRENTLY t_b_conc ON t (b)
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx RENAME TO t_a_idx2
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER INDEX t_a_idx RENAME TO t_a_idx2
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER INDEX t_a_idx RENAME TO t_b_key
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER INDEX missing RENAME TO x
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER INDEX IF EXISTS missing RENAME TO x
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx NOT VISIBLE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx VISIBLE
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TABLE t RENAME INDEX t_a_idx TO t_a_idx2
< error
| 42601 at or near "t_a_idx": syntax error
| detail: source SQL:\nALTER TABLE t RENAME INDEX t_a_idx TO t_a_idx2\n                           ^
| hint: try \\h ALTER TABLE
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX t_a_idx
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX t@t_a_idx
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX IF EXISTS t_a_idx
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX missing
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX t_b_key
< error
| 42704 index "t_b_key" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX t_b_key CASCADE
< error
| 42704 index "t_b_key" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX t@t_pkey
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP INDEX CONCURRENTLY t_a_idx
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON INDEX t@t_a_idx IS 'by a'
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON INDEX t_a_idx IS 'by a'
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON INDEX missing IS 'x'
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE sq
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE IF NOT EXISTS sq
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE sq2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE sq2 AS INT4 INCREMENT BY 2 MINVALUE 5 MAXVALUE 500 START WITH 7 CACHE 10
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq2 AS INT4 MINVALUE 5 MAXVALUE 500 INCREMENT 2 START 7 CACHE 10;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE sq2 INCREMENT BY -1
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE -9223372036854775808 MAXVALUE -1 INCREMENT -1 START -1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE sq2 CYCLE
< error
| 0A000 unimplemented: CYCLE option is not supported
| hint: You have attempted to use a feature that is not yet implemented.\nSee: https://go.crdb.dev/issue-v/20961/v25.2
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE sq2 OWNED BY t.a
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE sq2 OWNED BY t.missing
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE s.sq2
< error
| 3F000 cannot create "s.sq2" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.t MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE SEQUENCE sq2 MINVALUE 10 MAXVALUE 5
< error
| 22023 START value (10) cannot be greater than MAXVALUE (5)
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq INCREMENT BY 5
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq MINVALUE 2 MAXVALUE 100
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq NO MAXVALUE
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq MINVALUE 20 MAXVALUE 5
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq START WITH 50
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq RESTART WITH 50
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq CACHE 20
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq CYCLE
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq AS INT4
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq OWNED BY t.a
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq OWNED BY NONE
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE missing INCREMENT BY 5
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE IF EXISTS missing INCREMENT BY 5
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq RENAME TO sq2
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq RENAME TO t
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER SEQUENCE sq SET SCHEMA s
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SEQUENCE sq
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SEQUENCE sq CASCADE
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SEQUENCE IF EXISTS sq
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SEQUENCE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SEQUENCE missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP SEQUENCE t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v AS SELECT id FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW IF NOT EXISTS v AS SELECT id FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id, a FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT * FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 (x, y) AS SELECT id, a FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id, a + 1 AS a1, count(*) OVER () AS n FROM t WHERE a > 0
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id FROM v
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT t.id, c.note FROM t JOIN c ON c.t_id = t.id
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT missing FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT 1 AS one
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE VIEW public.v2 (
| 	one
| ) AS SELECT 1 AS one;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id FROM missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW s.v2 AS SELECT id FROM t
< error
| 3F000 cannot create "s.v2" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE VIEW t AS SELECT 1 AS one
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE VIEW public.t (
| 	one
| ) AS SELECT 1 AS one;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id, a FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id, a, b FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT a, id FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v2 AS SELECT id FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT * FROM k
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv2 AS SELECT id, a FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv AS SELECT id FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT id FROM t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv2 AS SELECT id FROM t WITH NO DATA
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER VIEW v RENAME TO v2
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER VIEW v RENAME TO t
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER VIEW missing RENAME TO v2
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER VIEW IF EXISTS missing RENAME TO v2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER VIEW v SET SCHEMA s
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER MATERIALIZED VIEW mv RENAME TO mv2
< error
| 42P01 relation "mv" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP VIEW v
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP VIEW v CASCADE
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP VIEW IF EXISTS v
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP VIEW IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP VIEW missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP VIEW t
< error
| 42P01 relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP VIEW mv
< error
| 42P01 relation "mv" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP MATERIALIZED VIEW mv
< error
| 42P01 relation "mv" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP MATERIALIZED VIEW v
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP MATERIALIZED VIEW IF EXISTS mv
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON VIEW v IS 'x'
< error
| 42601 at or near "view": syntax error
| detail: source SQL:\nCOMMENT ON VIEW v IS 'x'\n           ^
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE mood AS ENUM ('a')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('a');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE IF NOT EXISTS mood AS ENUM ('a')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('a');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ('red', 'green', 'blue')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.color AS ENUM ('red', 'green', 'blue');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ()
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.color AS ENUM ();
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ('red', 'red')
< error
| 42P17 enum definition contains duplicate value "red"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE s.color AS ENUM ('red')
< error
| 3F000 cannot create "s.color" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE missing.color AS ENUM ('red')
< error
| 3F000 cannot create "missing.color" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE point2 AS (x FLOAT8, y FLOAT8)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.point2 AS (x FLOAT8, y FLOAT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE point2 AS (m mood, label STRING)
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE pair AS (x INT8)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.pair AS (x INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| CREATE TYPE t AS ENUM ('a')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.t AS ENUM ('a');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' BEFORE 'sad'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' AFTER 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE IF NOT EXISTS 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' BEFORE 'missing'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'missing' TO 'blue'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'sad' TO 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'sad'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'missing'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood RENAME TO feeling
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood RENAME TO pair
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE mood SET SCHEMA s
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE missing ADD VALUE 'x'
< error
| 42704 type "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| ALTER TYPE pair ADD VALUE 'x'
< error
| 42704 type "pair" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TYPE mood
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TYPE pair
< error
| 42704 type "pair" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TYPE IF EXISTS mood
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TYPE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TYPE missing
< error
| 42704 type "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TYPE mood, pair
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| DROP TYPE t
< error
| 42704 type "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON TYPE mood IS 'feelings'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON TYPE pair IS 'two'
< error
| 42704 type "pair" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
|
< ok 0 -
> exec
| COMMENT ON TYPE missing IS 'x'
< error
| 42704 type "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SCHEMA s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SCHEMA IF NOT EXISTS s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SCHEMA s2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s2;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SCHEMA "Mixed Case"
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA "Mixed Case";
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SCHEMA s RENAME TO s2
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SCHEMA s RENAME TO public
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SCHEMA missing RENAME TO s2
< error
| 3F000 unknown schema "missing"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SCHEMA s
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SCHEMA s CASCADE
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SCHEMA missing
< error
| 3F000 unknown schema "missing"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SCHEMA public
< error
| 3F000 cannot drop schema "public"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON SCHEMA s IS 'a schema'
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON SCHEMA public IS 'the default'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
| COMMENT ON SCHEMA public IS 'the default';
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE t (id INT8 PRIMARY KEY)
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS t (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, v STRING NOT NULL)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	v STRING NOT NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (x INT8, y STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NULL,
| 	y STRING NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT n_pkey PRIMARY KEY (rowid ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (x INT8 NOT NULL, y INT8 NOT NULL, PRIMARY KEY (x, y DESC))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NOT NULL,
| 	y INT8 NOT NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (x ASC, y DESC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8, CONSTRAINT n_key PRIMARY KEY (id))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	CONSTRAINT n_key PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 DEFAULT 7, b STRING DEFAULT 'x', c TIMESTAMPTZ DEFAULT now())
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL DEFAULT 7:::INT8,
| 	b STRING NULL DEFAULT 'x':::STRING,
| 	c TIMESTAMPTZ NULL DEFAULT now():::TIMESTAMPTZ,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 CHECK (a > 0), b INT8, CHECK (b < 10), CONSTRAINT n_ab CHECK (a < b))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL,
| 	b INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT check_b CHECK (b < 10:::INT8),
| 	CONSTRAINT n_ab CHECK (a < b),
| 	CONSTRAINT check_a CHECK (a > 0:::INT8)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, e STRING UNIQUE, f STRING, UNIQUE (f))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	e STRING NULL,
| 	f STRING NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX n_e_key (e ASC),
| 	UNIQUE INDEX n_f_key (f ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8, INDEX (a), INDEX n_named (a DESC))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	INDEX n_a_idx (a ASC),
| 	INDEX n_named (a DESC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, price DECIMAL(10, 2), cents INT8 AS ((price * 100)::INT8) STORED, twice DECIMAL AS (price * 2) VIRTUAL)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	price DECIMAL(10,2) NULL,
| 	cents INT8 NULL AS ((price * 100:::DECIMAL)::INT8) STORED,
| 	twice DECIMAL NULL AS (price * 2:::DECIMAL) VIRTUAL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, secret STRING NOT VISIBLE, tags STRING[], doc JSONB, name VARCHAR(20) COLLATE "en_US")
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	secret STRING NOT VISIBLE NULL,
| 	tags STRING[] NULL,
| 	doc JSONB NULL,
| 	name VARCHAR(20) COLLATE en_US NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id SERIAL PRIMARY KEY, v STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL DEFAULT unique_rowid(),
| 	v STRING NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES t (id))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	t_id INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.n ADD CONSTRAINT n_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.n VALIDATE CONSTRAINT n_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES t)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	t_id INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.n ADD CONSTRAINT n_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.n VALIDATE CONSTRAINT n_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8, CONSTRAINT n_t FOREIGN KEY (t_id) REFERENCES t (id) ON DELETE CASCADE ON UPDATE SET NULL)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	t_id INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.n ADD CONSTRAINT n_t FOREIGN KEY (t_id) REFERENCES public.t(id) ON DELETE CASCADE ON UPDATE SET NULL;
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.n VALIDATE CONSTRAINT n_t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES missing (id))
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, k_x INT8 REFERENCES k (x))
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, m mood DEFAULT 'sad', p pair)
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, m missing_type)
< error
| 42704 type "missing_type" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, v INT8 DEFAULT nextval('sq'))
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, id INT8)
< error
| 42701 relation "n" (925): duplicate column name: "id"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 CHECK (missing > 0))
< error
| 42703 column "missing" does not exist, referenced in "missing > 0"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, INDEX (missing))
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY) WITH (ttl_expire_after = '1 day')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	crdb_internal_expiration TIMESTAMPTZ NOT VISIBLE NOT NULL DEFAULT current_timestamp():::TIMESTAMPTZ + '1 day':::INTERVAL ON UPDATE current_timestamp():::TIMESTAMPTZ + '1 day':::INTERVAL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| ) WITH (ttl = 'on', ttl_expire_after = '1 day':::INTERVAL);
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE s.n (id INT8 PRIMARY KEY)
< error
| 3F000 cannot create "s.n" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE s.u (id INT8 PRIMARY KEY)
< error
| 3F000 cannot create "s.u" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE missing.n (id INT8 PRIMARY KEY)
< error
| 3F000 cannot create "missing.n" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE v (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.v (
| 	id INT8 NOT NULL,
| 	CONSTRAINT v_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE sq (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.sq (
| 	id INT8 NOT NULL,
| 	CONSTRAINT sq_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE "Mixed" ("Id" INT8 PRIMARY KEY, "Some Value" STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public."Mixed" (
| 	"Id" INT8 NOT NULL,
| 	"Some Value" STRING NULL,
| 	CONSTRAINT "Mixed_pkey" PRIMARY KEY ("Id" ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TABLE n AS SELECT 1 AS x
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT n_pkey PRIMARY KEY (rowid ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TEMP TABLE n (id INT8 PRIMARY KEY)
< error
| XCEXF temporary tables are only supported experimentally
| hint: See: https://go.crdb.dev/issue-v/46260/v25.2\n--\nYou can enable temporary tables by running `SET experimental_enable_temp_tables = 'on'`.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME TO t2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t2 (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME TO c
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE IF EXISTS t RENAME TO t2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t2 (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE missing RENAME TO t2
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE IF EXISTS missing RENAME TO t2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t SET SCHEMA s
< error
| 3F000 unknown schema "s"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE k SET SCHEMA s
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE s.u RENAME TO u2
< error
| 42P01 relation "s.u" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE t CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE IF EXISTS t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE c
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE c, t
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE k
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE s.u
< error
| 42P01 relation "s.u" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TABLE v
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| TRUNCATE t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON TABLE t IS 'changed'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
| COMMENT ON TABLE public.t IS 'changed';
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON TABLE t IS NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON TABLE missing IS 'x'
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD n INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT NULL DEFAULT 5
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NOT NULL DEFAULT 5:::INT8,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NOT NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n STRING DEFAULT 'x' CHECK (n <> '')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n STRING NULL DEFAULT 'x':::STRING,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT check_n CHECK (n != '':::STRING)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n STRING UNIQUE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_n_key (n ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 REFERENCES t (id)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.t ADD CONSTRAINT t_n_fkey FOREIGN KEY (n) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.t VALIDATE CONSTRAINT t_n_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 AS (a * 2) STORED
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NULL AS (a * 2:::INT8) STORED,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 AS (a * 2) VIRTUAL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NULL AS (a * 2:::INT8) VIRTUAL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n mood
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NOT VISIBLE NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 DEFAULT nextval('sq')
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN a INT8
< error
| 42701 column "a" of relation "t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN IF NOT EXISTS a INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN IF NOT EXISTS n INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8, ADD COLUMN o STRING
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	n INT8 NULL,
| 	o STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE missing ADD COLUMN n INT8
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE k ADD COLUMN n INT8
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE v ADD COLUMN n INT8
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP b
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN m
< error
| 42703 column "m" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< error
| 42P10 column "id" is referenced by the primary key
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN missing
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS b
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN t_id
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN note
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE k DROP COLUMN y
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 9:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 'x'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL DEFAULT 'x':::STRING,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT nextval('missing')
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a DROP DEFAULT
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b DROP DEFAULT
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing SET DEFAULT 1
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET NOT NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NOT NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a DROP NOT NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN id DROP NOT NULL
< error
| 42P16 column "id" is in a primary index
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing SET NOT NULL
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b TYPE STRING(100)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING(100) NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET DATA TYPE VARCHAR(50)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b VARCHAR(50) NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a TYPE INT4
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT4 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a TYPE STRING
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a STRING NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing TYPE STRING
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET NOT VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NOT VISIBLE NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN b TO z
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	z STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME b TO z
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	z STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN a TO z
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	z INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN id TO z
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	z INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (z ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN b TO a
< error
| 42701 column "a" of relation "t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN missing TO z
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE c RENAME COLUMN t_id TO parent
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON COLUMN t.b IS 'bee'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
| COMMENT ON COLUMN public.t.b IS 'bee';
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON COLUMN t.a IS NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON COLUMN t.missing IS 'x'
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_pos CHECK (a > -1)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT t_pos CHECK (a > (-1):::INT8)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CHECK (a < 1000)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT check_a CHECK (a < 1000:::INT8)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_pos CHECK (a > -1) NOT VALID
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT t_pos CHECK (a > (-1):::INT8) NOT VALID
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_a_check CHECK (a > 1)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT t_a_check CHECK (a > 1:::INT8)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad CHECK (missing > 0)
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_b_unique UNIQUE (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_unique (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD UNIQUE (a, b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_a_b_key (a ASC, b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_b_key UNIQUE (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad UNIQUE (missing)
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_self FOREIGN KEY (a) REFERENCES t (id)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.t ADD CONSTRAINT t_self FOREIGN KEY (a) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.t VALIDATE CONSTRAINT t_self;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES t (id) ON DELETE CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.t ADD CONSTRAINT t_a_fkey FOREIGN KEY (a) REFERENCES public.t(id) ON DELETE CASCADE;
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.t VALIDATE CONSTRAINT t_a_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad FOREIGN KEY (a) REFERENCES missing (id)
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad FOREIGN KEY (a) REFERENCES t (b)
< error
| 42804 type of "a" (int) does not match foreign key "t"."b" (string)
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE c ADD CONSTRAINT c_note CHECK (note <> '')
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE k ADD PRIMARY KEY (x)
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ADD PRIMARY KEY (a)
< error
| 42611 multiple primary keys for table "t" are not allowed
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_a_check
< error
| 42704 constraint "t_a_check" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT IF EXISTS t_a_check
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT missing
< error
| 42704 constraint "missing" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_b_key
< error
| 42704 constraint "t_b_key" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_pkey
< error
| 0A000 relation "t" (1271): unimplemented: primary key dropped without subsequent addition of new primary key in same transaction
| hint: You have attempted to use a feature that is not yet implemented.\nSee: https://go.crdb.dev/issue-v/48026/v25.2
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE c DROP CONSTRAINT c_t_id_fkey
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_a_check TO t_a_nonneg
< error
| 42704 constraint "t_a_check" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_b_key TO t_b_uniq
< error
| 42704 constraint "t_b_key" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_pkey TO t_key
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_key PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT missing TO z
< error
| 42704 constraint "missing" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t VALIDATE CONSTRAINT t_a_check
< error
| 42704 constraint "t_a_check" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE c VALIDATE CONSTRAINT c_t_id_fkey
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t ALTER PRIMARY KEY USING COLUMNS (a)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (a ASC),
| 	UNIQUE INDEX t_id_key (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE k ALTER PRIMARY KEY USING COLUMNS (x)
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON CONSTRAINT t_a_check ON t IS 'non-negative'
< error
| 42704 constraint "t_a_check" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON CONSTRAINT t_pkey ON t IS 'the key'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
| COMMENT ON CONSTRAINT t_pkey ON public.t IS 'the key';
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON CONSTRAINT missing ON t IS 'x'
< error
| 42704 constraint "missing" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_b_idx (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_b_idx ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_b_idx (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_ab ON t (a DESC, b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_ab (a DESC, b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE UNIQUE INDEX t_b_uniq ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_uniq (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_b_store ON t (b) STORING (a)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_b_store (b ASC) STORING (a)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_b_part ON t (b) WHERE a > 0
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_b_part (b ASC) WHERE a > 0:::INT8
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_lower ON t (lower(b))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_lower (lower(b) ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_hash ON t (b) USING HASH
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	crdb_internal_b_shard_16 INT8 NOT VISIBLE NOT NULL AS (mod(fnv32(md5(crdb_internal.datums_to_bytes(b))), 16:::INT8)) VIRTUAL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_hash (b ASC) USING HASH WITH (bucket_count=16)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_hidden ON t (b) NOT VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_hidden (b ASC) NOT VISIBLE
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INVERTED INDEX t_doc ON t (doc)
< error
| 42703 column "doc" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_a_idx ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_a_idx (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX IF NOT EXISTS t_a_idx ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_a_idx (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX IF NOT EXISTS t_b_idx ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_b_idx (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX t_bad ON t (missing)
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX n_idx ON missing (a)
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX ON k (x)
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX ON s.u (v)
< error
| 42P01 relation "s.u" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX mv_b ON mv (b)
< error
| 42P01 relation "mv" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX v_a ON v (a)
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE INDEX CONCURRENTLY t_b_conc ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_b_conc (b ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx RENAME TO t_a_idx2
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER INDEX t_a_idx RENAME TO t_a_idx2
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER INDEX t_a_idx RENAME TO t_b_key
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER INDEX missing RENAME TO x
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER INDEX IF EXISTS missing RENAME TO x
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx NOT VISIBLE
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx VISIBLE
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TABLE t RENAME INDEX t_a_idx TO t_a_idx2
< error
| 42601 at or near "t_a_idx": syntax error
| detail: source SQL:\nALTER TABLE t RENAME INDEX t_a_idx TO t_a_idx2\n                           ^
| hint: try \\h ALTER TABLE
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX t_a_idx
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX t@t_a_idx
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX IF EXISTS t_a_idx
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX missing
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX t_b_key
< error
| 42704 index "t_b_key" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX t_b_key CASCADE
< error
| 42704 index "t_b_key" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX t@t_pkey
< error
| 0A000 cannot drop the primary index of a table using DROP INDEX
| hint: instead, use ALTER TABLE ... ALTER PRIMARY KEY oruse DROP CONSTRAINT ... PRIMARY KEY followed by ADD CONSTRAINT ... PRIMARY KEY in a transaction
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP INDEX CONCURRENTLY t_a_idx
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON INDEX t@t_a_idx IS 'by a'
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON INDEX t_a_idx IS 'by a'
< error
| 42704 index "t_a_idx" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON INDEX missing IS 'x'
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE sq
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE IF NOT EXISTS sq
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE sq2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE sq2 AS INT4 INCREMENT BY 2 MINVALUE 5 MAXVALUE 500 START WITH 7 CACHE 10
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE SEQUENCE public.sq2 AS INT4 MINVALUE 5 MAXVALUE 500 INCREMENT 2 START 7 CACHE 10;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE sq2 INCREMENT BY -1
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE -9223372036854775808 MAXVALUE -1 INCREMENT -1 START -1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE sq2 CYCLE
< error
| 0A000 unimplemented: CYCLE option is not supported
| hint: You have attempted to use a feature that is not yet implemented.\nSee: https://go.crdb.dev/issue-v/20961/v25.2
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE sq2 OWNED BY t.a
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE sq2 OWNED BY t.missing
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE s.sq2
< error
| 3F000 cannot create "s.sq2" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE t
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE SEQUENCE sq2 MINVALUE 10 MAXVALUE 5
< error
| 22023 START value (10) cannot be greater than MAXVALUE (5)
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq INCREMENT BY 5
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq MINVALUE 2 MAXVALUE 100
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq NO MAXVALUE
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq MINVALUE 20 MAXVALUE 5
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq START WITH 50
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq RESTART WITH 50
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq CACHE 20
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq CYCLE
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq AS INT4
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq OWNED BY t.a
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq OWNED BY NONE
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE missing INCREMENT BY 5
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE IF EXISTS missing INCREMENT BY 5
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq RENAME TO sq2
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq RENAME TO t
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER SEQUENCE sq SET SCHEMA s
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SEQUENCE sq
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SEQUENCE sq CASCADE
< error
| 42P01 relation "sq" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SEQUENCE IF EXISTS sq
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SEQUENCE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SEQUENCE missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP SEQUENCE t
< error
| 42809 "t" is not a sequence
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id
| ) AS SELECT id FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW IF NOT EXISTS v AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id
| ) AS SELECT id FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id, a FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v2 (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT * FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v2 (
| 	id,
| 	a,
| 	b
| ) AS SELECT t.id, t.a, t.b FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 (x, y) AS SELECT id, a FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v2 (
| 	x,
| 	y
| ) AS SELECT id, a FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id, a + 1 AS a1, count(*) OVER () AS n FROM t WHERE a > 0
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v2 (
| 	id,
| 	a1,
| 	n
| ) AS SELECT id, a + 1 AS a1, count(*) OVER () AS n FROM sim.public.t WHERE a > 0;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id FROM v
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT t.id, c.note FROM t JOIN c ON c.t_id = t.id
< error
| 42P01 relation "c" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT missing FROM t
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT 1 AS one
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v2 (
| 	one
| ) AS SELECT 1 AS one;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id FROM missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW s.v2 AS SELECT id FROM t
< error
| 3F000 cannot create "s.v2" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE VIEW t AS SELECT 1 AS one
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id, a FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id, a, b FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a,
| 	b
| ) AS SELECT id, a, b FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id
| ) AS SELECT id FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT a, id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	a,
| 	id
| ) AS SELECT a, id FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v2 AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v2 (
| 	id
| ) AS SELECT id FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT * FROM k
< error
| 42P01 relation "k" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv2 AS SELECT id, a FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE MATERIALIZED VIEW public.mv2 (
| 	id,
| 	a,
| 	rowid
| ) AS SELECT id, a FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	id,
| 	rowid
| ) AS SELECT id FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	id,
| 	rowid
| ) AS SELECT id FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv2 AS SELECT id FROM t WITH NO DATA
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE MATERIALIZED VIEW public.mv2 (
| 	id,
| 	rowid
| ) AS SELECT id FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER VIEW v RENAME TO v2
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER VIEW v RENAME TO t
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER VIEW missing RENAME TO v2
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER VIEW IF EXISTS missing RENAME TO v2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER VIEW v SET SCHEMA s
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER MATERIALIZED VIEW mv RENAME TO mv2
< error
| 42P01 relation "mv" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP VIEW v
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP VIEW v CASCADE
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP VIEW IF EXISTS v
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP VIEW IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP VIEW missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP VIEW t
< error
| 42809 "t" is not a view
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP VIEW mv
< error
| 42P01 relation "mv" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP MATERIALIZED VIEW mv
< error
| 42P01 relation "mv" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP MATERIALIZED VIEW v
< error
| 42P01 relation "v" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP MATERIALIZED VIEW IF EXISTS mv
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON VIEW v IS 'x'
< error
| 42601 at or near "view": syntax error
| detail: source SQL:\nCOMMENT ON VIEW v IS 'x'\n           ^
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE mood AS ENUM ('a')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('a');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE IF NOT EXISTS mood AS ENUM ('a')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('a');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ('red', 'green', 'blue')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.color AS ENUM ('red', 'green', 'blue');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ()
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.color AS ENUM ();
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ('red', 'red')
< error
| 42P17 enum definition contains duplicate value "red"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE s.color AS ENUM ('red')
< error
| 3F000 cannot create "s.color" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE missing.color AS ENUM ('red')
< error
| 3F000 cannot create "missing.color" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE point2 AS (x FLOAT8, y FLOAT8)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.point2 AS (x FLOAT8, y FLOAT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE point2 AS (m mood, label STRING)
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE pair AS (x INT8)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.pair AS (x INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| CREATE TYPE t AS ENUM ('a')
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' BEFORE 'sad'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' AFTER 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE IF NOT EXISTS 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' BEFORE 'missing'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'missing' TO 'blue'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'sad' TO 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'sad'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'happy'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'missing'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood RENAME TO feeling
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood RENAME TO pair
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE mood SET SCHEMA s
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE missing ADD VALUE 'x'
< error
| 42704 type "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| ALTER TYPE pair ADD VALUE 'x'
< error
| 42704 type "pair" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TYPE mood
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TYPE pair
< error
| 42704 type "pair" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TYPE IF EXISTS mood
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TYPE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TYPE missing
< error
| 42704 type "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TYPE mood, pair
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| DROP TYPE t
< error
| 2BP01 cannot modify table record type "t"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON TYPE mood IS 'feelings'
< error
| 42704 type "mood" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON TYPE pair IS 'two'
< error
| 42704 type "pair" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING);
< ok 0 -
> exec
| COMMENT ON TYPE missing IS 'x'
< error
| 42704 type "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SCHEMA s
< error
| 42P06 schema "s" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SCHEMA IF NOT EXISTS s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SCHEMA s2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
< row
: string
| CREATE SCHEMA s2;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SCHEMA "Mixed Case"
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA "Mixed Case";
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SCHEMA s RENAME TO s2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s2;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s2.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SCHEMA s RENAME TO public
< error
| 42P06 schema "public" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SCHEMA missing RENAME TO s2
< error
| 3F000 unknown schema "missing"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SCHEMA s
< error
| 2BP01 schema "s" is not empty and CASCADE was not specified
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SCHEMA s CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS s
< error
| 2BP01 schema "s" is not empty and CASCADE was not specified
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SCHEMA missing
< error
| 3F000 unknown schema "missing"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SCHEMA public
< error
| 3F000 cannot drop schema "public"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON SCHEMA s IS 'a schema'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
| COMMENT ON SCHEMA s IS 'a schema';
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON SCHEMA public IS 'the default'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
| COMMENT ON SCHEMA public IS 'the default';
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE t (id INT8 PRIMARY KEY)
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS t (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, v STRING NOT NULL)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	v STRING NOT NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (x INT8, y STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NULL,
| 	y STRING NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT n_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (x INT8 NOT NULL, y INT8 NOT NULL, PRIMARY KEY (x, y DESC))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NOT NULL,
| 	y INT8 NOT NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (x ASC, y DESC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8, CONSTRAINT n_key PRIMARY KEY (id))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	CONSTRAINT n_key PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 DEFAULT 7, b STRING DEFAULT 'x', c TIMESTAMPTZ DEFAULT now())
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL DEFAULT 7:::INT8,
| 	b STRING NULL DEFAULT 'x':::STRING,
| 	c TIMESTAMPTZ NULL DEFAULT now():::TIMESTAMPTZ,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 CHECK (a > 0), b INT8, CHECK (b < 10), CONSTRAINT n_ab CHECK (a < b))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL,
| 	b INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT check_b CHECK (b < 10:::INT8),
| 	CONSTRAINT n_ab CHECK (a < b),
| 	CONSTRAINT check_a CHECK (a > 0:::INT8)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, e STRING UNIQUE, f STRING, UNIQUE (f))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	e STRING NULL,
| 	f STRING NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX n_e_key (e ASC),
| 	UNIQUE INDEX n_f_key (f ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8, INDEX (a), INDEX n_named (a DESC))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	a INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC),
| 	INDEX n_a_idx (a ASC),
| 	INDEX n_named (a DESC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, price DECIMAL(10, 2), cents INT8 AS ((price * 100)::INT8) STORED, twice DECIMAL AS (price * 2) VIRTUAL)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	price DECIMAL(10,2) NULL,
| 	cents INT8 NULL AS ((price * 100:::DECIMAL)::INT8) STORED,
| 	twice DECIMAL NULL AS (price * 2:::DECIMAL) VIRTUAL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, secret STRING NOT VISIBLE, tags STRING[], doc JSONB, name VARCHAR(20) COLLATE "en_US")
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	secret STRING NOT VISIBLE NULL,
| 	tags STRING[] NULL,
| 	doc JSONB NULL,
| 	name VARCHAR(20) COLLATE en_US NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id SERIAL PRIMARY KEY, v STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL DEFAULT unique_rowid(),
| 	v STRING NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES t (id))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	t_id INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| ALTER TABLE public.n ADD CONSTRAINT n_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
< row
: string
| ALTER TABLE public.n VALIDATE CONSTRAINT n_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES t)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	t_id INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| ALTER TABLE public.n ADD CONSTRAINT n_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
< row
: string
| ALTER TABLE public.n VALIDATE CONSTRAINT n_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8, CONSTRAINT n_t FOREIGN KEY (t_id) REFERENCES t (id) ON DELETE CASCADE ON UPDATE SET NULL)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	t_id INT8 NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| ALTER TABLE public.n ADD CONSTRAINT n_t FOREIGN KEY (t_id) REFERENCES public.t(id) ON DELETE CASCADE ON UPDATE SET NULL;
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
< row
: string
| ALTER TABLE public.n VALIDATE CONSTRAINT n_t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, t_id INT8 REFERENCES missing (id))
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, k_x INT8 REFERENCES k (x))
< error
| 23503 there is no unique constraint matching given keys for referenced table k
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, m mood DEFAULT 'sad', p pair)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	m public.mood NULL DEFAULT 'sad':::public.mood,
| 	p public.pair NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, m missing_type)
< error
| 42704 type "missing_type" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, v INT8 DEFAULT nextval('sq'))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	v INT8 NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, id INT8)
< error
| 42701 relation "n" (2338): duplicate column name: "id"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, a INT8 CHECK (missing > 0))
< error
| 42703 column "missing" does not exist, referenced in "missing > 0"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY, INDEX (missing))
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n (id INT8 PRIMARY KEY) WITH (ttl_expire_after = '1 day')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	id INT8 NOT NULL,
| 	crdb_internal_expiration TIMESTAMPTZ NOT VISIBLE NOT NULL DEFAULT current_timestamp():::TIMESTAMPTZ + '1 day':::INTERVAL ON UPDATE current_timestamp():::TIMESTAMPTZ + '1 day':::INTERVAL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| ) WITH (ttl = 'on', ttl_expire_after = '1 day':::INTERVAL);
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE s.n (id INT8 PRIMARY KEY)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE s.n (
| 	id INT8 NOT NULL,
| 	CONSTRAINT n_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE s.u (id INT8 PRIMARY KEY)
< error
| 42P07 relation "sim.s.u" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE missing.n (id INT8 PRIMARY KEY)
< error
| 3F000 cannot create "missing.n" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE v (id INT8 PRIMARY KEY)
< error
| 42P07 relation "sim.public.v" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE sq (id INT8 PRIMARY KEY)
< error
| 42P07 relation "sim.public.sq" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE "Mixed" ("Id" INT8 PRIMARY KEY, "Some Value" STRING)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public."Mixed" (
| 	"Id" INT8 NOT NULL,
| 	"Some Value" STRING NULL,
| 	CONSTRAINT "Mixed_pkey" PRIMARY KEY ("Id" ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TABLE n AS SELECT 1 AS x
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE TABLE public.n (
| 	x INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT n_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TEMP TABLE n (id INT8 PRIMARY KEY)
< error
| XCEXF temporary tables are only supported experimentally
| hint: See: https://go.crdb.dev/issue-v/46260/v25.2\n--\nYou can enable temporary tables by running `SET experimental_enable_temp_tables = 'on'`.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME TO t2
< error
| 2BP01 cannot rename relation "t" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME TO c
< error
| 2BP01 cannot rename relation "t" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE IF EXISTS t RENAME TO t2
< error
| 2BP01 cannot rename relation "t" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE missing RENAME TO t2
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE IF EXISTS missing RENAME TO t2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t SET SCHEMA s
< error
| 2BP01 cannot set schema on relation "t" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE k SET SCHEMA s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE s.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE s.u RENAME TO u2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u2 (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE t
< error
| 2BP01 cannot drop relation "t" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE t CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE IF EXISTS t
< error
| 2BP01 cannot drop relation "t" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE c
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE c, t
< error
| 2BP01 cannot drop relation "t" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE k
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE s.u
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TABLE v
< error
| 42809 "v" is not a table
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| TRUNCATE t
< error
| XXUUU "t" is referenced by foreign key from table "c"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON TABLE t IS 'changed'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'changed';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON TABLE t IS NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON TABLE missing IS 'x'
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD n INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT NULL DEFAULT 5
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NOT NULL DEFAULT 5:::INT8,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NOT NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n STRING DEFAULT 'x' CHECK (n <> '')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n STRING NULL DEFAULT 'x':::STRING,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8),
| 	CONSTRAINT check_n CHECK (n != '':::STRING)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n STRING UNIQUE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	UNIQUE INDEX t_n_key (n ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 REFERENCES t (id)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.t ADD CONSTRAINT t_n_fkey FOREIGN KEY (n) REFERENCES public.t(id);
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.t VALIDATE CONSTRAINT t_n_fkey;
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 AS (a * 2) STORED
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NULL AS (a * 2:::INT8) STORED,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 AS (a * 2) VIRTUAL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NULL AS (a * 2:::INT8) VIRTUAL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n mood
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n public.mood NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 NOT VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NOT VISIBLE NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8 DEFAULT nextval('sq')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN a INT8
< error
| 42701 column "a" of relation "t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN IF NOT EXISTS a INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN IF NOT EXISTS n INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD COLUMN n INT8, ADD COLUMN o STRING
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	n INT8 NULL,
| 	o STRING NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE missing ADD COLUMN n INT8
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE k ADD COLUMN n INT8
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	n INT8 NULL,
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE v ADD COLUMN n INT8
< error
| 42809 "v" is not a table
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< error
| 2BP01 cannot drop column "b" because view "mv" depends on it
| hint: consider dropping "mv" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP b
< error
| 2BP01 cannot drop column "b" because view "mv" depends on it
| hint: consider dropping "mv" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 2BP01 cannot drop column "a" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC)
| );
| COMMENT ON TABLE public.t IS 'the table';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN m
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< error
| 42P10 column "id" is referenced by the primary key
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN missing
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN IF EXISTS b
< error
| 2BP01 cannot drop column "b" because view "mv" depends on it
| hint: consider dropping "mv" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN t_id
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN note
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE k DROP COLUMN y
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 9:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 'x'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL DEFAULT 'x':::STRING,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT nextval('missing')
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a DROP DEFAULT
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b DROP DEFAULT
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing SET DEFAULT 1
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET NOT NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NOT NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a DROP NOT NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN id DROP NOT NULL
< error
| 42P16 column "id" is in a primary index
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing SET NOT NULL
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b TYPE STRING(100)
< error
| 2BP01 cannot alter type of column "b" because view "mv" depends on it
| hint: consider dropping "mv" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET DATA TYPE VARCHAR(50)
< error
| 2BP01 cannot alter type of column "b" because view "mv" depends on it
| hint: consider dropping "mv" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a TYPE INT4
< error
| 2BP01 cannot alter type of column "a" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN a TYPE STRING
< error
| 2BP01 cannot alter type of column "a" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN missing TYPE STRING
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET NOT VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NOT VISIBLE NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN b SET VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN b TO z
< error
| 2BP01 cannot rename column "b" because view "mv" depends on it
| hint: consider dropping "mv" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME b TO z
< error
| 2BP01 cannot rename column "b" because view "mv" depends on it
| hint: consider dropping "mv" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN a TO z
< error
| 2BP01 cannot rename column "a" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN id TO z
< error
| 2BP01 cannot rename column "id" because view "v" depends on it
| hint: consider dropping "v" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN b TO a
< error
| 2BP01 cannot rename column "b" because view "mv" depends on it
| hint: consider dropping "mv" first.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME COLUMN missing TO z
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE c RENAME COLUMN t_id TO parent
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	parent INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (parent) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON COLUMN t.b IS 'bee'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
| COMMENT ON COLUMN public.t.b IS 'bee';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON COLUMN t.a IS NULL
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON COLUMN t.missing IS 'x'
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_pos CHECK (a > -1)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8),
| 	CONSTRAINT t_pos CHECK (a > (-1):::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CHECK (a < 1000)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8),
| 	CONSTRAINT check_a CHECK (a < 1000:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_pos CHECK (a > -1) NOT VALID
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8),
| 	CONSTRAINT t_pos CHECK (a > (-1):::INT8) NOT VALID
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_a_check CHECK (a > 1)
< error
| 42710 duplicate constraint name: "t_a_check"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad CHECK (missing > 0)
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_b_unique UNIQUE (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	UNIQUE INDEX t_b_unique (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD UNIQUE (a, b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	UNIQUE INDEX t_a_b_key (a ASC, b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_b_key UNIQUE (b)
< error
| 42P07 index with name "t_b_key" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad UNIQUE (missing)
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_self FOREIGN KEY (a) REFERENCES t (id)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.t ADD CONSTRAINT t_self FOREIGN KEY (a) REFERENCES public.t(id);
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.t VALIDATE CONSTRAINT t_self;
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES t (id) ON DELETE CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.t ADD CONSTRAINT t_a_fkey FOREIGN KEY (a) REFERENCES public.t(id) ON DELETE CASCADE;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.t VALIDATE CONSTRAINT t_a_fkey;
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad FOREIGN KEY (a) REFERENCES missing (id)
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_bad FOREIGN KEY (a) REFERENCES t (b)
< error
| 42804 type of "a" (int) does not match foreign key "t"."b" (string)
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE c ADD CONSTRAINT c_note CHECK (note <> '')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC),
| 	CONSTRAINT c_note CHECK (note != '':::STRING)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE k ADD PRIMARY KEY (x)
< error
| 42P15 cannot use nullable column "x" in primary key
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ADD PRIMARY KEY (a)
< error
| 42611 multiple primary keys for table "t" are not allowed
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_a_check
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT IF EXISTS t_a_check
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT missing
< error
| 42704 constraint "missing" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_b_key
< error
| 0A000 unimplemented: cannot drop UNIQUE constraint "t_b_key" using ALTER TABLE DROP CONSTRAINT, use DROP INDEX CASCADE instead
| hint: You have attempted to use a feature that is not yet implemented.\nSee: https://go.crdb.dev/issue-v/42840/v25.2
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_pkey
< error
| 0A000 relation "t" (3921): unimplemented: primary key dropped without subsequent addition of new primary key in same transaction
| hint: You have attempted to use a feature that is not yet implemented.\nSee: https://go.crdb.dev/issue-v/48026/v25.2
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE c DROP CONSTRAINT c_t_id_fkey
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_a_check TO t_a_nonneg
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_nonneg CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_b_key TO t_b_uniq
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_uniq (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_pkey TO t_key
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_key PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT missing TO z
< error
| 42704 constraint "missing" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t VALIDATE CONSTRAINT t_a_check
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE c VALIDATE CONSTRAINT c_t_id_fkey
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t ALTER PRIMARY KEY USING COLUMNS (a)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (a ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	UNIQUE INDEX t_id_key (id ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE k ALTER PRIMARY KEY USING COLUMNS (x)
< error
| 42P15 cannot use nullable column "x" in primary key
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON CONSTRAINT t_a_check ON t IS 'non-negative'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
| COMMENT ON CONSTRAINT t_a_check ON public.t IS 'non-negative';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON CONSTRAINT t_pkey ON t IS 'the key'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
| COMMENT ON CONSTRAINT t_pkey ON public.t IS 'the key';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON CONSTRAINT missing ON t IS 'x'
< error
| 42704 constraint "missing" of relation "t" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_b_idx (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_b_idx ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_b_idx (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_ab ON t (a DESC, b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_ab (a DESC, b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE UNIQUE INDEX t_b_uniq ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	UNIQUE INDEX t_b_uniq (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_b_store ON t (b) STORING (a)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_b_store (b ASC) STORING (a),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_b_part ON t (b) WHERE a > 0
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_b_part (b ASC) WHERE a > 0:::INT8,
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_lower ON t (lower(b))
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_lower (lower(b) ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_hash ON t (b) USING HASH
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	crdb_internal_b_shard_16 INT8 NOT VISIBLE NOT NULL AS (mod(fnv32(md5(crdb_internal.datums_to_bytes(b))), 16:::INT8)) VIRTUAL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_hash (b ASC) USING HASH WITH (bucket_count=16),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_hidden ON t (b) NOT VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_hidden (b ASC) NOT VISIBLE,
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INVERTED INDEX t_doc ON t (doc)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INVERTED INDEX t_doc (doc),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_a_idx ON t (b)
< error
| 42P07 index with name "t_a_idx" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX IF NOT EXISTS t_a_idx ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX IF NOT EXISTS t_b_idx ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_b_idx (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX t_bad ON t (missing)
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX n_idx ON missing (a)
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX ON k (x)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC),
| 	INDEX k_x_idx (x ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX ON s.u (v)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC),
| 	INDEX u_v_idx (v ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX mv_b ON mv (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX v_a ON v (a)
< error
| 42809 "v" is not a table or materialized view
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE INDEX CONCURRENTLY t_b_conc ON t (b)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	INDEX t_b_conc (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx RENAME TO t_a_idx2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx2 (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER INDEX t_a_idx RENAME TO t_a_idx2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx2 (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER INDEX t_a_idx RENAME TO t_b_key
< error
| 42P07 index name "t_b_key" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER INDEX missing RENAME TO x
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER INDEX IF EXISTS missing RENAME TO x
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx NOT VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC) NOT VISIBLE,
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER INDEX t@t_a_idx VISIBLE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TABLE t RENAME INDEX t_a_idx TO t_a_idx2
< error
| 42601 at or near "t_a_idx": syntax error
| detail: source SQL:\nALTER TABLE t RENAME INDEX t_a_idx TO t_a_idx2\n                           ^
| hint: try \\h ALTER TABLE
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX t_a_idx
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX t@t_a_idx
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX IF EXISTS t_a_idx
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX missing
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX t_b_key
< error
| 2BP01 index "t_b_key" is in use as unique constraint
| hint: use CASCADE if you really want to drop it.
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX t_b_key CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX t@t_pkey
< error
| 0A000 cannot drop the primary index of a table using DROP INDEX
| hint: instead, use ALTER TABLE ... ALTER PRIMARY KEY oruse DROP CONSTRAINT ... PRIMARY KEY followed by ADD CONSTRAINT ... PRIMARY KEY in a transaction
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP INDEX CONCURRENTLY t_a_idx
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON INDEX t@t_a_idx IS 'by a'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
| COMMENT ON INDEX public.t@t_a_idx IS 'by a';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON INDEX t_a_idx IS 'by a'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
| COMMENT ON INDEX public.t@t_a_idx IS 'by a';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON INDEX missing IS 'x'
< error
| 42704 index "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE sq
< error
| 42P07 relation "sim.public.sq" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE IF NOT EXISTS sq
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE sq2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE sq2 AS INT4 INCREMENT BY 2 MINVALUE 5 MAXVALUE 500 START WITH 7 CACHE 10
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE SEQUENCE public.sq2 AS INT4 MINVALUE 5 MAXVALUE 500 INCREMENT 2 START 7 CACHE 10;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE sq2 INCREMENT BY -1
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE -9223372036854775808 MAXVALUE -1 INCREMENT -1 START -1;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE sq2 CYCLE
< error
| 0A000 unimplemented: CYCLE option is not supported
| hint: You have attempted to use a feature that is not yet implemented.\nSee: https://go.crdb.dev/issue-v/20961/v25.2
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE sq2 OWNED BY t.a
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE sq2 OWNED BY t.missing
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE s.sq2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE SEQUENCE s.sq2 MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 1;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE t
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE SEQUENCE sq2 MINVALUE 10 MAXVALUE 5
< error
| 22023 START value (10) cannot be greater than MAXVALUE (5)
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq INCREMENT BY 5
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 5 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq MINVALUE 2 MAXVALUE 100
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 2 MAXVALUE 100 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq NO MAXVALUE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq MINVALUE 20 MAXVALUE 5
< error
| 22023 START value (10) cannot be greater than MAXVALUE (5)
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq START WITH 50
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 50;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq RESTART WITH 50
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq CACHE 20
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10 CACHE 20;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq CYCLE
< error
| 0A000 unimplemented: CYCLE option is not supported
| hint: You have attempted to use a feature that is not yet implemented.\nSee: https://go.crdb.dev/issue-v/20961/v25.2
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq AS INT4
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq AS INT4 MINVALUE 1 MAXVALUE 2147483647 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq OWNED BY t.a
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq OWNED BY NONE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE missing INCREMENT BY 5
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE IF EXISTS missing INCREMENT BY 5
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq RENAME TO sq2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq2 MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq2'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq RENAME TO t
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER SEQUENCE sq SET SCHEMA s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE s.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('s.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SEQUENCE sq
< error
| 2BP01 cannot drop sequence sq because other objects depend on it
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SEQUENCE sq CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL,
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SEQUENCE IF EXISTS sq
< error
| 2BP01 cannot drop sequence sq because other objects depend on it
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SEQUENCE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SEQUENCE missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP SEQUENCE t
< error
| 42809 "t" is not a sequence
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v AS SELECT id FROM t
< error
| 42P07 relation "sim.public.v" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW IF NOT EXISTS v AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id, a FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW public.v2 (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT * FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW public.v2 (
| 	id,
| 	a,
| 	b,
| 	m,
| 	doc
| ) AS SELECT t.id, t.a, t.b, t.m, t.doc FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 (x, y) AS SELECT id, a FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW public.v2 (
| 	x,
| 	y
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id, a + 1 AS a1, count(*) OVER () AS n FROM t WHERE a > 0
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW public.v2 (
| 	id,
| 	a1,
| 	n
| ) AS SELECT id, a + 1 AS a1, count(*) OVER () AS n FROM sim.public.t WHERE a > 0;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id FROM v
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW public.v2 (
| 	id
| ) AS SELECT id FROM sim.public.v;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT t.id, c.note FROM t JOIN c ON c.t_id = t.id
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW public.v2 (
| 	id,
| 	note
| ) AS SELECT t.id, c.note FROM sim.public.t JOIN sim.public.c ON c.t_id = t.id;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT missing FROM t
< error
| 42703 column "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT 1 AS one
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW public.v2 (
| 	one
| ) AS SELECT 1 AS one;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW v2 AS SELECT id FROM missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW s.v2 AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW s.v2 (
| 	id
| ) AS SELECT id FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE VIEW t AS SELECT 1 AS one
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id, a FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id, a, b FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a,
| 	b
| ) AS SELECT id, a, b FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT id FROM t
< error
| 42P16 cannot drop columns from view
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT a, id FROM t
< error
| 42P16 cannot change name of view column "id" to "a"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v2 AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE VIEW public.v2 (
| 	id
| ) AS SELECT id FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE OR REPLACE VIEW v AS SELECT * FROM k
< error
| 42P16 cannot change name of view column "id" to "x"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv2 AS SELECT id, a FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv2 (
| 	id,
| 	a,
| 	rowid
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv AS SELECT id FROM t
< error
| 42P07 relation "sim.public.mv" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW IF NOT EXISTS mv AS SELECT id FROM t
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv2 AS SELECT id FROM t WITH NO DATA
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv2 (
| 	id,
| 	rowid
| ) AS SELECT id FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER VIEW v RENAME TO v2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v2 (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER VIEW v RENAME TO t
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER VIEW missing RENAME TO v2
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER VIEW IF EXISTS missing RENAME TO v2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER VIEW v SET SCHEMA s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW s.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER MATERIALIZED VIEW mv RENAME TO mv2
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv2 (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP VIEW v
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP VIEW v CASCADE
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP VIEW IF EXISTS v
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP VIEW IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP VIEW missing
< error
| 42P01 relation "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP VIEW t
< error
| 42809 "t" is not a view
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP VIEW mv
< error
| 42809 "mv" is a materialized view
| hint: use the corresponding MATERIALIZED VIEW command
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP MATERIALIZED VIEW mv
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP MATERIALIZED VIEW v
< error
| 42809 "v" is not a materialized view
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP MATERIALIZED VIEW IF EXISTS mv
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON VIEW v IS 'x'
< error
| 42601 at or near "view": syntax error
| detail: source SQL:\nCOMMENT ON VIEW v IS 'x'\n           ^
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE mood AS ENUM ('a')
< error
| 42710 type "sim.public.mood" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE IF NOT EXISTS mood AS ENUM ('a')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ('red', 'green', 'blue')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
< row
: string
| CREATE TYPE public.color AS ENUM ('red', 'green', 'blue');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ()
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
< row
: string
| CREATE TYPE public.color AS ENUM ();
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE color AS ENUM ('red', 'red')
< error
| 42P17 enum definition contains duplicate value "red"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE s.color AS ENUM ('red')
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
< row
: string
| CREATE TYPE s.color AS ENUM ('red');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE missing.color AS ENUM ('red')
< error
| 3F000 cannot create "missing.color" because the target database or schema does not exist
| hint: verify that the current database and search_path are valid and/or the target database exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE point2 AS (x FLOAT8, y FLOAT8)
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
< row
: string
| CREATE TYPE public.point2 AS (x FLOAT8, y FLOAT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE point2 AS (m mood, label STRING)
< error
| 0A000 unimplemented: composite types that reference user-defined types not yet supported
| hint: You have attempted to use a feature that is not yet implemented.\nSee: https://go.crdb.dev/issue-v/91779/v25.2
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE pair AS (x INT8)
< error
| 42710 type "sim.public.pair" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| CREATE TYPE t AS ENUM ('a')
< error
| 42P07 relation "sim.public.t" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad', 'meh');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' BEFORE 'sad'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'meh', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' AFTER 'happy'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'meh', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'happy'
< error
| 42710 enum value "happy" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE IF NOT EXISTS 'happy'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood ADD VALUE 'meh' BEFORE 'missing'
< error
| 22023 "missing" is not an existing enum value
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'sad' TO 'blue'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'blue');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'missing' TO 'blue'
< error
| 22023 missing is not an existing enum value
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood RENAME VALUE 'sad' TO 'happy'
< error
| 42710 enum value happy already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'sad'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'happy'
< error
| 2BP01 could not remove enum value "happy" as it is being used in a default expresion of "t"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood DROP VALUE 'missing'
< error
| 42704 enum value "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood RENAME TO feeling
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.feeling AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.feeling NULL DEFAULT 'happy':::public.feeling,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood RENAME TO pair
< error
| 42710 type "pair" already exists
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE mood SET SCHEMA s
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE s.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m s.mood NULL DEFAULT 'happy':::s.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE missing ADD VALUE 'x'
< error
| 42704 type "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| ALTER TYPE pair ADD VALUE 'x'
< error
| 42809 "pair" is not an enum
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TYPE mood
< error
| 2BP01 cannot drop type "mood" because other objects ([sim.public.t]) still depend on it
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TYPE pair
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TYPE IF EXISTS mood
< error
| 2BP01 cannot drop type "mood" because other objects ([sim.public.t]) still depend on it
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TYPE IF EXISTS missing
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TYPE missing
< error
| 42704 type "missing" does not exist
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TYPE mood, pair
< error
| 2BP01 cannot drop type "mood" because other objects ([sim.public.t]) still depend on it
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| DROP TYPE t
< error
| 2BP01 cannot modify table record type "t"
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON TYPE mood IS 'feelings'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
| COMMENT ON TYPE public.mood IS 'feelings';
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON TYPE pair IS 'two'
< ok 0 -
> query
| SELECT create_statement FROM [SHOW CREATE ALL SCHEMAS]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SCHEMA public;
< row
: string
| CREATE SCHEMA s;
> query
| SELECT create_statement FROM [SHOW CREATE ALL TYPES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE TYPE public.mood AS ENUM ('happy', 'sad');
< row
: string
| CREATE TYPE public.pair AS (x INT8, y INT8);
| COMMENT ON TYPE public.pair IS 'two';
> query
| SELECT create_statement FROM [SHOW CREATE ALL TABLES]
< result
< column "create_statement" "TEXT"
< row
: string
| CREATE SEQUENCE public.sq MINVALUE 1 MAXVALUE 9223372036854775807 INCREMENT 1 START 10;
< row
: string
| CREATE TABLE public.t (
| 	id INT8 NOT NULL,
| 	a INT8 NOT NULL DEFAULT 0:::INT8,
| 	b STRING NULL,
| 	m public.mood NULL DEFAULT 'happy':::public.mood,
| 	doc JSONB NULL,
| 	CONSTRAINT t_pkey PRIMARY KEY (id ASC),
| 	UNIQUE INDEX t_b_key (b ASC),
| 	INDEX t_a_idx (a ASC),
| 	CONSTRAINT t_a_check CHECK (a >= 0:::INT8)
| );
| COMMENT ON TABLE public.t IS 'the table';
| COMMENT ON COLUMN public.t.a IS 'a count';
< row
: string
| CREATE TABLE public.c (
| 	id INT8 NOT NULL DEFAULT nextval('public.sq'::REGCLASS),
| 	t_id INT8 NULL,
| 	note STRING NULL,
| 	CONSTRAINT c_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE TABLE public.k (
| 	x INT8 NULL,
| 	y INT8 NULL,
| 	rowid INT8 NOT VISIBLE NOT NULL DEFAULT unique_rowid(),
| 	CONSTRAINT k_pkey PRIMARY KEY (rowid ASC)
| );
< row
: string
| CREATE TABLE s.u (
| 	id INT8 NOT NULL,
| 	v STRING NULL,
| 	CONSTRAINT u_pkey PRIMARY KEY (id ASC)
| );
< row
: string
| CREATE VIEW public.v (
| 	id,
| 	a
| ) AS SELECT id, a FROM sim.public.t;
< row
: string
| CREATE MATERIALIZED VIEW public.mv (
| 	b,
| 	rowid
| ) AS SELECT b FROM sim.public.t;
< row
: string
| ALTER TABLE public.c ADD CONSTRAINT c_t_id_fkey FOREIGN KEY (t_id) REFERENCES public.t(id);
< row
: string
| -- Validate foreign key constraints. These can fail if there was unvalidated data during the SHOW CREATE ALL TABLES
< row
: string
| ALTER TABLE public.c VALIDATE CONSTRAINT c_t_id_fkey;
> exec
| USE defaultdb;
| DROP DATABASE IF EXISTS sim CASCADE;
| CREATE DATABASE sim;
| USE sim;
| CREATE SCHEMA s;
| CREATE TYPE mood AS ENUM ('happy', 'sad');
| CREATE TYPE pair AS (x INT8, y INT8);
| CREATE SEQUENCE sq START WITH 10;
| CREATE TABLE t (id INT8 PRIMARY KEY, a INT8 NOT NULL DEFAULT 0, b STRING, m mood DEFAULT 'happy', doc JSONB, CONSTRAINT t_a_check CHECK (a >= 0), CONSTRAINT t_b_key UNIQUE (b));
| CREATE INDEX t_a_idx ON t (a);
| COMMENT ON TABLE t IS 'the table';
| COMMENT ON COLUMN t.a IS 'a count';
| CREATE TABLE c (id INT8 PRIMARY KEY DEFAULT nextval('sq'), t_id INT8 REFERENCES t (id), note STRING);
| CREATE TABLE k (x INT8, y INT8);
| CREATE TABLE s.u (id INT8 PRIMARY KEY, v STRING);
| CREATE VIEW v AS SELECT id, a FROM t;
| CREATE MATERIALIZED VIEW mv AS SELECT b FROM t;
< ok 0 -
> exec
| COMMENT ON TYPE missing IS 'x'
< error
| 42704 type "missing" does not exist
