-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestDependencyPG17Oracle default#1
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_0" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_0"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_0"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, n int CHECK (n > 0))
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_0
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.n
< row
: string
| constraint t.t_n_check
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_0" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_1" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_1"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_1"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, n int CHECK (n > 0))
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN n
< ok 0 -
> exec
| ALTER TABLE t ADD CONSTRAINT t_n_check CHECK (id > 0)
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_1
< result
< column "case" "TEXT"
< row
: string
| index t_pkey
< row
: string
| relation t
< row
: string
| column t.id
< row
: string
| constraint t.t_pkey index t_pkey
< row
: string
| constraint t.t_n_check
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_1" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_2" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_2"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_2"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int, CHECK (a < b))
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_2
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.b
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_2" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_3" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_3"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_3"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, code text UNIQUE, n int)
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN code
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_3
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.n
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_3" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_4" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_4"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_4"
< ok 0 -
> exec
| CREATE TABLE t (a int, name text)
< ok 0 -
> exec
| CREATE INDEX t_lower ON t (lower(name))
< ok 0 -
> exec
| CREATE INDEX t_part ON t (a) WHERE name IS NOT NULL
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN name
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_4
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.a
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_4" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_5" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_5"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_5"
< ok 0 -
> exec
| CREATE TABLE t (a int, b serial)
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_5
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.a
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_5" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_6" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_6"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_6"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED)
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 2BP01 cannot drop column a of table t because other objects depend on it
| detail: column b of table t depends on column a of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_6
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.a
< row
: string
| column t.b
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_6" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_7" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_7"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_7"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED, c int)
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a CASCADE
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_7
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.c
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_7" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_8" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_8"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_8"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED)
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_8
< result
< column "case" "TEXT"
< row
: string
| relation t
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_8" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_9" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_9"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_9"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_9
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| relation v
< row
: string
| column t.a
< row
: string
| column v.a
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_9" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_10" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_10"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_10"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT a FROM t
< ok 0 -
> exec
| CREATE VIEW vv AS SELECT a FROM v
< ok 0 -
> exec
| CREATE VIEW w AS SELECT b FROM t
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 2BP01 cannot drop column a of table t because other objects depend on it
| detail: view v depends on column a of table t\nview vv depends on view v
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE t DROP COLUMN a CASCADE
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_10
< result
< column "case" "TEXT"
< row
: string
| relation w
< row
: string
| relation t
< row
: string
| column w.b
< row
: string
| column t.b
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_10" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_11" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_11"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_11"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT a FROM t ORDER BY b
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< error
| 2BP01 cannot drop column b of table t because other objects depend on it
| detail: view v depends on column b of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_11
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| relation v
< row
: string
| column t.a
< row
: string
| column t.b
< row
: string
| column v.a
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_11" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_12" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_12"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_12"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT a, count(*) FROM t GROUP BY a HAVING max(b) > 0
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< error
| 2BP01 cannot drop column b of table t because other objects depend on it
| detail: view v depends on column b of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_12
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| relation v
< row
: string
| column t.a
< row
: string
| column t.b
< row
: string
| column v.a
< row
: string
| column v.count
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_12" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_13" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_13"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_13"
< ok 0 -
> exec
| CREATE TABLE u (id int, name text)
< ok 0 -
> exec
| CREATE TABLE o (id int, uid int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT u.name FROM u JOIN o ON o.uid = u.id
< ok 0 -
> exec
| ALTER TABLE o DROP COLUMN uid
< error
| 2BP01 cannot drop column uid of table o because other objects depend on it
| detail: view v depends on column uid of table o
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE o DROP COLUMN id
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_13
< result
< column "case" "TEXT"
< row
: string
| relation u
< row
: string
| relation o
< row
: string
| relation v
< row
: string
| column u.id
< row
: string
| column u.name
< row
: string
| column o.uid
< row
: string
| column v.name
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_13" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_14" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_14"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_14"
< ok 0 -
> exec
| CREATE TABLE u (id int, name text)
< ok 0 -
> exec
| CREATE TABLE o (id int, uid int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT u.name FROM u WHERE EXISTS (SELECT 1 FROM o WHERE o.uid = u.id)
< ok 0 -
> exec
| ALTER TABLE u DROP COLUMN id
< error
| 2BP01 cannot drop column id of table u because other objects depend on it
| detail: view v depends on column id of table u
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE o DROP COLUMN id
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_14
< result
< column "case" "TEXT"
< row
: string
| relation u
< row
: string
| relation o
< row
: string
| relation v
< row
: string
| column u.id
< row
: string
| column u.name
< row
: string
| column o.uid
< row
: string
| column v.name
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_14" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_15" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_15"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_15"
< ok 0 -
> exec
| CREATE TABLE u (id int, name text)
< ok 0 -
> exec
| CREATE TABLE o (id int, uid int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT u.name FROM u WHERE u.id IN (SELECT uid FROM o)
< ok 0 -
> exec
| ALTER TABLE u DROP COLUMN id
< error
| 2BP01 cannot drop column id of table u because other objects depend on it
| detail: view v depends on column id of table u
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_15
< result
< column "case" "TEXT"
< row
: string
| relation u
< row
: string
| relation o
< row
: string
| relation v
< row
: string
| column u.id
< row
: string
| column u.name
< row
: string
| column o.id
< row
: string
| column o.uid
< row
: string
| column v.name
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_15" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_16" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_16"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_16"
< ok 0 -
> exec
| CREATE TABLE t (a int, x int)
< ok 0 -
> exec
| CREATE TABLE u (b int, x int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT a FROM t JOIN u ON true
< ok 0 -
> exec
| CREATE VIEW w AS SELECT x FROM t JOIN u USING (x)
< ok 0 -
> exec
| CREATE VIEW z AS SELECT b FROM (t JOIN u USING (x)) j
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 2BP01 cannot drop column a of table t because other objects depend on it
| detail: view v depends on column a of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE u DROP COLUMN b
< error
| 2BP01 cannot drop column b of table u because other objects depend on it
| detail: view z depends on column b of table u
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE t DROP COLUMN x
< error
| 2BP01 cannot drop column x of table t because other objects depend on it
| detail: view w depends on column x of table t\nview z depends on column x of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE u DROP COLUMN x
< error
| 2BP01 cannot drop column x of table u because other objects depend on it
| detail: view w depends on column x of table u\nview z depends on column x of table u
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE t ALTER COLUMN a TYPE bigint
< error
| 0A000 cannot alter type of a column used by a view or rule
| detail: rule _RETURN on view v depends on column "a"
> exec
| DROP VIEW v, z
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< ok 0 -
> exec
| ALTER TABLE u ALTER COLUMN b TYPE bigint
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_16
< result
< column "case" "TEXT"
< row
: string
| relation u
< row
: string
| relation t
< row
: string
| relation w
< row
: string
| column u.b
< row
: string
| column u.x
< row
: string
| column t.x
< row
: string
| column w.x
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_16" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_17" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_17"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_17"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT x FROM (SELECT a AS x, b AS y FROM t) s
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< error
| 2BP01 cannot drop column b of table t because other objects depend on it
| detail: view v depends on column b of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_17
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| relation v
< row
: string
| column t.a
< row
: string
| column t.b
< row
: string
| column v.x
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_17" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_18" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_18"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_18"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, n int)
< ok 0 -
> exec
| CREATE TABLE c (tid int REFERENCES t (id))
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< error
| 2BP01 cannot drop column id of table t because other objects depend on it
| detail: constraint c_tid_fkey on table c depends on column id of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE t DROP COLUMN n
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id CASCADE
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_18
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| relation c
< row
: string
| column c.tid
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_18" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_19" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_19"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_19"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, parent int REFERENCES t (id))
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< error
| 2BP01 cannot drop column id of table t because other objects depend on it
| detail: constraint t_parent_fkey on table t depends on column id of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE t DROP COLUMN parent
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_19
< result
< column "case" "TEXT"
< row
: string
| index t_pkey
< row
: string
| relation t
< row
: string
| column t.id
< row
: string
| constraint t.t_pkey index t_pkey
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_19" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_20" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_20"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_20"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int)
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$
< ok 0 -
> exec
| CREATE TRIGGER tg BEFORE UPDATE OF b ON t FOR EACH ROW EXECUTE FUNCTION f()
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< error
| 2BP01 cannot drop column b of table t because other objects depend on it
| detail: trigger tg on table t depends on column b of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE t DROP COLUMN a
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_20
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.b
< row
: string
| trigger t.tg
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_20" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_21" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_21"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_21"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int, c int, d int, e int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT c FROM t
< ok 0 -
> exec
| CREATE UNIQUE INDEX t_e ON t (e)
< ok 0 -
> exec
| COMMENT ON COLUMN t.b IS 'b'
< ok 0 -
> exec
| COMMENT ON COLUMN t.c IS 'c'
< ok 0 -
> exec
| COMMENT ON COLUMN t.e IS 'e'
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN a
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN b
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN c
< error
| 2BP01 cannot drop column c of table t because other objects depend on it
| detail: view v depends on column c of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE t DROP COLUMN d
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN e
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_21
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| relation v
< row
: string
| column t.c comment c
< row
: string
| column v.c
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_21" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_22" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_22"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_22"
< ok 0 -
> exec
| CREATE TABLE p (id int, x int, y int)
< ok 0 -
> exec
| CREATE TABLE c (y int) INHERITS (p)
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN x
< error
| 42P16 cannot drop inherited column "x"
> exec
| ALTER TABLE c DROP COLUMN y
< error
| 42P16 cannot drop inherited column "y"
> exec
| ALTER TABLE p DROP COLUMN x
< ok 0 -
> exec
| ALTER TABLE p DROP COLUMN y
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN y
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_22
< result
< column "case" "TEXT"
< row
: string
| relation c
< row
: string
| relation p
< row
: string
| column c.id
< row
: string
| column p.id
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_22" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_23" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_23"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_23"
< ok 0 -
> exec
| CREATE TABLE p (id int, x int, y int) PARTITION BY RANGE (id)
< ok 0 -
> exec
| CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10)
< ok 0 -
> exec
| CREATE TABLE p2 (id int, x int, y int)
< ok 0 -
> exec
| ALTER TABLE p ATTACH PARTITION p2 FOR VALUES FROM (10) TO (20)
< ok 0 -
> exec
| ALTER TABLE p1 DROP COLUMN x
< error
| 42P16 cannot drop inherited column "x"
> exec
| ALTER TABLE p DROP COLUMN x
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_23
< result
< column "case" "TEXT"
< row
: string
| relation p2
< row
: string
| relation p1
< row
: string
| relation p
< row
: string
| column p2.id
< row
: string
| column p2.y
< row
: string
| column p1.id
< row
: string
| column p1.y
< row
: string
| column p.id
< row
: string
| column p.y
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_23" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_24" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_24"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_24"
< ok 0 -
> exec
| CREATE TABLE p (id int, a int)
< ok 0 -
> exec
| CREATE TABLE c () INHERITS (p)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT a FROM c
< ok 0 -
> exec
| ALTER TABLE p DROP COLUMN a
< error
| 2BP01 cannot drop desired object(s) because other objects depend on them
| detail: view v depends on column a of table c
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_24
< result
< column "case" "TEXT"
< row
: string
| relation c
< row
: string
| relation p
< row
: string
| relation v
< row
: string
| column c.id
< row
: string
| column c.a
< row
: string
| column p.id
< row
: string
| column p.a
< row
: string
| column v.a
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_24" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_25" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_25"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_25"
< ok 0 -
> exec
| CREATE TABLE p (id int, a int) PARTITION BY RANGE (id)
< ok 0 -
> exec
| CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT a FROM p1
< ok 0 -
> exec
| ALTER TABLE p DROP COLUMN a
< error
| 2BP01 cannot drop desired object(s) because other objects depend on them
| detail: view v depends on column a of table p1
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE p DROP COLUMN a CASCADE
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_25
< result
< column "case" "TEXT"
< row
: string
| relation p1
< row
: string
| relation p
< row
: string
| column p1.id
< row
: string
| column p.id
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_25" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_26" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_26"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_26"
< ok 0 -
> exec
| CREATE TABLE p (id int, a int, b int) PARTITION BY RANGE (id)
< ok 0 -
> exec
| CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10)
< ok 0 -
> exec
| CREATE INDEX p_a ON p (a)
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$
< ok 0 -
> exec
| CREATE TRIGGER tg BEFORE UPDATE OF b ON p FOR EACH ROW EXECUTE FUNCTION f()
< ok 0 -
> exec
| ALTER TABLE p DROP COLUMN a
< ok 0 -
> exec
| ALTER TABLE p DROP COLUMN b
< error
| 2BP01 cannot drop desired object(s) because other objects depend on them
| detail: trigger tg on table p depends on column b of table p
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE p DROP COLUMN b CASCADE
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_26
< result
< column "case" "TEXT"
< row
: string
| relation p
< row
: string
| relation p1
< row
: string
| column p.id
< row
: string
| column p1.id
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_26" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_27" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_27"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_27"
< ok 0 -
> exec
| CREATE TABLE p (id int, a int)
< ok 0 -
> exec
| CREATE TABLE c () INHERITS (p)
< ok 0 -
> exec
| CREATE TABLE g () INHERITS (c)
< ok 0 -
> exec
| CREATE TABLE q (id int, a int)
< ok 0 -
> exec
| CREATE TABLE q1 () INHERITS (q)
< ok 0 -
> exec
| CREATE TABLE q2 () INHERITS (q1)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT a FROM q2
< ok 0 -
> exec
| ALTER TABLE q DROP COLUMN a
< error
| 2BP01 cannot drop desired object(s) because other objects depend on them
| detail: view v depends on column a of table q2
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE p DROP COLUMN a
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_27
< result
< column "case" "TEXT"
< row
: string
| relation p
< row
: string
| relation g
< row
: string
| relation c
< row
: string
| relation q
< row
: string
| relation q2
< row
: string
| relation q1
< row
: string
| relation v
< row
: string
| column p.id
< row
: string
| column g.id
< row
: string
| column c.id
< row
: string
| column q.id
< row
: string
| column q.a
< row
: string
| column q2.id
< row
: string
| column q2.a
< row
: string
| column q1.id
< row
: string
| column q1.a
< row
: string
| column v.a
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_27" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_28" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_28"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_28"
< ok 0 -
> exec
| CREATE TABLE p (id int, a int, b int) PARTITION BY RANGE (id)
< ok 0 -
> exec
| CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10) PARTITION BY RANGE (b)
< ok 0 -
> exec
| CREATE TABLE p11 PARTITION OF p1 FOR VALUES FROM (0) TO (5)
< ok 0 -
> exec
| ALTER TABLE p DROP COLUMN a
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_28
< result
< column "case" "TEXT"
< row
: string
| relation p
< row
: string
| relation p11
< row
: string
| relation p1
< row
: string
| column p.id
< row
: string
| column p.b
< row
: string
| column p11.id
< row
: string
| column p11.b
< row
: string
| column p1.id
< row
: string
| column p1.b
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_28" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_29" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_29"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_29"
< ok 0 -
> exec
| CREATE TABLE p (id int, a int, b int) PARTITION BY RANGE (id)
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$
< ok 0 -
> exec
| CREATE TRIGGER tg BEFORE INSERT ON p FOR EACH ROW WHEN (NEW.b > 0) EXECUTE FUNCTION f()
< ok 0 -
> exec
| CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10)
< ok 0 -
> exec
| ALTER TABLE p DROP COLUMN b
< error
| 2BP01 cannot drop desired object(s) because other objects depend on them
| detail: trigger tg on table p depends on column b of table p
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE p DROP COLUMN b CASCADE
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_29
< result
< column "case" "TEXT"
< row
: string
| relation p
< row
: string
| relation p1
< row
: string
| column p.id
< row
: string
| column p.a
< row
: string
| column p1.id
< row
: string
| column p1.a
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_29" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_30" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_30"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_30"
< ok 0 -
> exec
| CREATE TABLE p (id int, a int)
< ok 0 -
> exec
| CREATE TABLE c () INHERITS (p)
< ok 0 -
> exec
| CREATE TABLE p2 (id int, a int)
< ok 0 -
> exec
| CREATE TABLE c2 (a int) INHERITS (p2)
< ok 0 -
> exec
| ALTER TABLE ONLY p DROP COLUMN a
< ok 0 -
> exec
| ALTER TABLE ONLY p2 DROP COLUMN a
< ok 0 -
> exec
| ALTER TABLE c DROP COLUMN a
< ok 0 -
> exec
| ALTER TABLE c2 DROP COLUMN a
< ok 0 -
> exec
| CREATE TABLE q (id int, a int) PARTITION BY RANGE (id)
< ok 0 -
> exec
| CREATE TABLE q1 PARTITION OF q FOR VALUES FROM (0) TO (10)
< ok 0 -
> exec
| ALTER TABLE ONLY q DROP COLUMN a
< error
| 42P16 cannot drop column from only the partitioned table when partitions exist
| hint: Do not specify the ONLY keyword.
> exec
| CREATE TABLE r (id int, a int) PARTITION BY RANGE (id)
< ok 0 -
> exec
| ALTER TABLE ONLY r DROP COLUMN a
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_30
< result
< column "case" "TEXT"
< row
: string
| relation c
< row
: string
| relation p
< row
: string
| relation c2
< row
: string
| relation p2
< row
: string
| relation q1
< row
: string
| relation q
< row
: string
| relation r
< row
: string
| column c.id
< row
: string
| column p.id
< row
: string
| column c2.id
< row
: string
| column p2.id
< row
: string
| column q1.id
< row
: string
| column q1.a
< row
: string
| column q.id
< row
: string
| column q.a
< row
: string
| column r.id
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_30" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_31" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_31"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_31"
< ok 0 -
> exec
| CREATE TABLE p (a int, b int) PARTITION BY RANGE (a)
< ok 0 -
> exec
| CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10) PARTITION BY RANGE (b)
< ok 0 -
> exec
| ALTER TABLE p DROP COLUMN b
< error
| 42P16 cannot drop column "b" because it is part of the partition key of relation "p1"
> exec
| ALTER TABLE p DROP COLUMN b CASCADE
< error
| 42P16 cannot drop column "b" because it is part of the partition key of relation "p1"
> exec
| ALTER TABLE p DROP COLUMN a
< error
| 42P16 cannot drop column "a" because it is part of the partition key of relation "p"
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_31
< result
< column "case" "TEXT"
< row
: string
| relation p1
< row
: string
| relation p
< row
: string
| column p1.a
< row
: string
| column p1.b
< row
: string
| column p.a
< row
: string
| column p.b
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_31" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_32" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_32"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_32"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, k int)
< ok 0 -
> exec
| CREATE TABLE p (id int, tid int REFERENCES t (id)) PARTITION BY RANGE (id)
< ok 0 -
> exec
| CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10)
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< error
| 2BP01 cannot drop column id of table t because other objects depend on it
| detail: constraint p_tid_fkey on table p depends on column id of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| ALTER TABLE t DROP COLUMN id CASCADE
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_32
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| relation p
< row
: string
| relation p1
< row
: string
| column t.k
< row
: string
| column p.id
< row
: string
| column p.tid
< row
: string
| column p1.id
< row
: string
| column p1.tid
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_32" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_33" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_33"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_33"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, n int CHECK (n > 0), m int GENERATED ALWAYS AS (n * 2) STORED, parent int REFERENCES t (id))
< ok 0 -
> exec
| DROP TABLE t
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_33
< result
< column "case" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_33" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_34" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_34"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_34"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int)
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$
< ok 0 -
> exec
| CREATE TRIGGER tg BEFORE UPDATE OF b ON t FOR EACH ROW EXECUTE FUNCTION f()
< ok 0 -
> exec
| DROP TABLE t
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_34
< result
< column "case" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_34" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_35" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_35"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_35"
< ok 0 -
> exec
| CREATE TABLE t (a int, b int GENERATED ALWAYS AS (a * 2) STORED)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT b FROM t
< ok 0 -
> exec
| DROP TABLE t
< error
| 2BP01 cannot drop table t because other objects depend on it
| detail: view v depends on table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| DROP TABLE t CASCADE
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_35
< result
< column "case" "TEXT"
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_35" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_36" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_36"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_36"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY)
< ok 0 -
> exec
| CREATE TABLE c (tid int REFERENCES t (id))
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN id TYPE bigint
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_36
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| index t_pkey
< row
: string
| relation c
< row
: string
| column t.id
< row
: string
| column c.tid
< row
: string
| constraint c.c_tid_fkey
< row
: string
| constraint t.t_pkey index t_pkey
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_36" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_37" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_37"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_37"
< ok 0 -
> exec
| CREATE TABLE t (id int, n int CHECK (n > 0), m int, g int GENERATED ALWAYS AS (m * 2) STORED, k int, w int)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT k FROM t
< ok 0 -
> exec
| CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN RETURN NEW; END$$
< ok 0 -
> exec
| CREATE TRIGGER tg BEFORE UPDATE OF w ON t FOR EACH ROW EXECUTE FUNCTION f()
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN id TYPE bigint
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN n TYPE bigint
< ok 0 -
> exec
| ALTER TABLE t ALTER COLUMN m TYPE bigint
< error
| 0A000 cannot alter type of a column used by a generated column
| detail: Column "m" is used by generated column "g".
> exec
| ALTER TABLE t ALTER COLUMN k TYPE bigint
< error
| 0A000 cannot alter type of a column used by a view or rule
| detail: rule _RETURN on view v depends on column "k"
> exec
| ALTER TABLE t ALTER COLUMN w TYPE bigint
< error
| 0A000 cannot alter type of a column used in a trigger definition
| detail: trigger tg on table t depends on column "w"
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_37
< result
< column "case" "TEXT"
< row
: string
| relation v
< row
: string
| relation t
< row
: string
| column v.k
< row
: string
| column t.id
< row
: string
| column t.n
< row
: string
| column t.m
< row
: string
| column t.g
< row
: string
| column t.k
< row
: string
| column t.w
< row
: string
| trigger t.tg
< row
: string
| constraint t.t_n_check
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_37" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_38" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_38"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_38"
< ok 0 -
> exec
| CREATE TABLE t (c int CHECK (c > 0), d int, CHECK (d > 0), CHECK (c < d))
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_c_check
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_d_check
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_check
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_38
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| column t.c
< row
: string
| column t.d
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_38" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_39" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_39"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_39"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE)
< ok 0 -
> exec
| ALTER INDEX t_code_key RENAME TO t_code_uk
< ok 0 -
> exec
| ALTER INDEX t_pkey RENAME TO t_pk
< ok 0 -
> exec
| ALTER TABLE t DROP CONSTRAINT t_code_uk
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_39
< result
< column "case" "TEXT"
< row
: string
| index t_pk
< row
: string
| relation t
< row
: string
| column t.id
< row
: string
| column t.code
< row
: string
| constraint t.t_pk index t_pk
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_39" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_40" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_40"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_40"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE)
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_code_key TO t_code_uk
< ok 0 -
> exec
| DROP INDEX t_code_key
< error
| 42704 index "t_code_key" does not exist
> exec
| ALTER INDEX t_code_uk RENAME TO t_code_uk2
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_pkey TO t_pk
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_40
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| index t_code_uk2
< row
: string
| index t_pk
< row
: string
| column t.id
< row
: string
| column t.code
< row
: string
| constraint t.t_code_uk2 index t_code_uk2
< row
: string
| constraint t.t_pk index t_pk
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_40" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_41" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_41"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_41"
< ok 0 -
> exec
| CREATE TABLE t (r int4range, EXCLUDE USING gist (r WITH &&))
< ok 0 -
> exec
| ALTER INDEX t_r_excl RENAME TO t_ex
< ok 0 -
> exec
| ALTER TABLE t RENAME CONSTRAINT t_ex TO t_ex2
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_41
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| index t_ex2
< row
: string
| column t.r
< row
: string
| constraint t.t_ex2 index t_ex2
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_41" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_42" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_42"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_42"
< ok 0 -
> exec
| CREATE TABLE p (id int, PRIMARY KEY (id)) PARTITION BY RANGE (id)
< ok 0 -
> exec
| CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (10)
< ok 0 -
> exec
| ALTER INDEX p1_pkey RENAME TO p1_pk
< ok 0 -
> exec
| ALTER INDEX p_pkey RENAME TO p_pk
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_42
< result
< column "case" "TEXT"
< row
: string
| relation p1
< row
: string
| relation p
< row
: string
| index p1_pk
< row
: string
| index p_pk
< row
: string
| column p1.id
< row
: string
| column p.id
< row
: string
| constraint p1.p1_pk index p1_pk
< row
: string
| constraint p.p_pk index p_pk
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_42" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_43" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "dep_oracle_43"
< ok 0 -
> exec
| SET search_path TO "dep_oracle_43"
< ok 0 -
> exec
| CREATE TABLE t (id int PRIMARY KEY, code text CONSTRAINT t_code_key UNIQUE, CONSTRAINT foo CHECK (id > 0))
< ok 0 -
> exec
| CREATE INDEX t_id_idx ON t (id)
< ok 0 -
> exec
| CREATE SEQUENCE s
< ok 0 -
> exec
| CREATE TABLE u (x int)
< ok 0 -
> exec
| ALTER INDEX t_code_key RENAME TO foo
< error
| 42710 constraint "foo" for relation "t" already exists
> exec
| ALTER INDEX t_code_key RENAME TO s
< error
| 42P07 relation "s" already exists
> exec
| ALTER INDEX t_code_key RENAME TO t_id_idx
< error
| 42P07 relation "t_id_idx" already exists
> exec
| ALTER TABLE t RENAME CONSTRAINT t_code_key TO u
< error
| 42P07 relation "u" already exists
> exec
| ALTER TABLE t RENAME CONSTRAINT t_code_key TO t_pkey
< error
| 42P07 relation "t_pkey" already exists
> exec
| ALTER TABLE t RENAME CONSTRAINT t_code_key TO foo
< error
| 42710 constraint "foo" for relation "t" already exists
> exec
| ALTER TABLE t RENAME CONSTRAINT foo TO t_pkey
< error
| 42710 constraint "t_pkey" for relation "t" already exists
> exec
| ALTER TABLE t RENAME CONSTRAINT nope TO bar
< error
| 42704 constraint "nope" for table "t" does not exist
> exec
| ALTER INDEX t_id_idx RENAME TO t_id_idx2
< ok 0 -
> query
|
| SELECT CASE WHEN c.relkind IN ('i', 'I') THEN 'index ' || c.relname
|             WHEN c.relkind = 'S' THEN 'sequence ' || c.relname
|             ELSE 'relation ' || c.relname END
| FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm', 'i', 'I', 'S')
| UNION ALL
| SELECT 'column ' || c.relname || '.' || a.attname
|        || coalesce(' comment ' || col_description(c.oid, a.attnum), '')
| FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid
|      JOIN pg_namespace n ON n.oid = c.relnamespace
| WHERE n.nspname = $1 AND c.relkind IN ('r', 'p', 'v', 'm')
|       AND a.attnum > 0 AND NOT a.attisdropped
| UNION ALL
| SELECT 'trigger ' || r.relname || '.' || tg.tgname
| FROM pg_trigger tg JOIN pg_class r ON r.oid = tg.tgrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
| WHERE n.nspname = $1 AND NOT tg.tgisinternal
| UNION ALL
| SELECT 'constraint ' || r.relname || '.' || con.conname
|        || CASE WHEN con.contype IN ('p', 'u', 'x') THEN ' index ' || i.relname ELSE '' END
| FROM pg_constraint con JOIN pg_class r ON r.oid = con.conrelid
|      JOIN pg_namespace n ON n.oid = r.relnamespace
|      LEFT JOIN pg_class i ON i.oid = con.conindid
| WHERE n.nspname = $1
: string
| dep_oracle_43
< result
< column "case" "TEXT"
< row
: string
| relation t
< row
: string
| index t_pkey
< row
: string
| index t_code_key
< row
: string
| sequence s
< row
: string
| relation u
< row
: string
| index t_id_idx2
< row
: string
| column t.id
< row
: string
| column t.code
< row
: string
| column u.x
< row
: string
| constraint t.t_pkey index t_pkey
< row
: string
| constraint t.t_code_key index t_code_key
< row
: string
| constraint t.foo
> exec
| DROP SCHEMA IF EXISTS "dep_oracle_43" CASCADE
< ok 0 -
