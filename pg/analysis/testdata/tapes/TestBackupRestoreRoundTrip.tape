-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestBackupRestoreRoundTrip default#1
> exec
| CREATE SCHEMA bbdataarchive
< ok 0 -
> exec
| CREATE TABLE canary (id INT)
< ok 0 -
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, a INT); INSERT INTO plain VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT x::text FROM "public"."plain" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
> query
| SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = $1 AND tablename = $2)
: string
| public
: string
| plain
< result
< column "exists" "BOOL"
< row
: bool true
> exec
| CREATE TABLE "bbdataarchive"."_bb__plain_public" AS
|   SELECT "public"."plain".* FROM plain WHERE id = 1;
< ok 1 -
> exec
| DELETE FROM plain WHERE id = 1
< ok 1 -
> query
| SELECT x::text FROM "public"."plain" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (2,2)
> query
| SELECT a.attname, a.attgenerated <> '' FROM pg_catalog.pg_attribute a
| 		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum
: string
| public
: string
| plain
< result
< column "attname" "NAME"
< column "?column?" "BOOL"
< row
: string
| id
: bool false
< row
: string
| a
: bool false
> query
| SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL,
| 		(SELECT string_agg(a.attname, E'\x1f' ORDER BY k.n) FROM unnest(x.indkey) WITH ORDINALITY AS k(attnum, n)
| 			JOIN pg_catalog.pg_attribute a ON a.attrelid = x.indrelid AND a.attnum = k.attnum)
| 		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
| 		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
| 		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname
: string
| public
: string
| plain
< result
< column "relname" "NAME"
< column "indisprimary" "BOOL"
< column "indisunique" "BOOL"
< column "?column?" "BOOL"
< column "string_agg" "TEXT"
< row
: string
| plain_pkey
: bool true
: bool true
: bool true
: string
| id
> exec
| /*
| Original SQL:
| DELETE FROM plain WHERE id = 1
| */
| INSERT INTO "public"."plain" ("id", "a") SELECT "id", "a" FROM "bbdataarchive"."_bb__plain_public";
< ok 1 -
> query
| SELECT x::text FROM "public"."plain" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
> query
| SELECT to_regclass($1)::text
: string
| canary
< result
< column "to_regclass" "TEXT"
< row
: string
| canary
> exec
| CREATE TABLE plain_update (id INT PRIMARY KEY, a INT); INSERT INTO plain_update VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT x::text FROM "public"."plain_update" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
> query
| SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = $1 AND tablename = $2)
: string
| public
: string
| plain_update
< result
< column "exists" "BOOL"
< row
: bool true
> exec
| CREATE TABLE "bbdataarchive"."_bb__plain_update_public" AS
|   SELECT "public"."plain_update".* FROM plain_update WHERE id = 2;
< ok 1 -
> exec
| UPDATE plain_update SET a = 10 WHERE id = 2
< ok 1 -
> query
| SELECT x::text FROM "public"."plain_update" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,10)
> query
| SELECT a.attname, a.attgenerated <> '' FROM pg_catalog.pg_attribute a
| 		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum
: string
| public
: string
| plain_update
< result
< column "attname" "NAME"
< column "?column?" "BOOL"
< row
: string
| id
: bool false
< row
: string
| a
: bool false
> query
| SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL,
| 		(SELECT string_agg(a.attname, E'\x1f' ORDER BY k.n) FROM unnest(x.indkey) WITH ORDINALITY AS k(attnum, n)
| 			JOIN pg_catalog.pg_attribute a ON a.attrelid = x.indrelid AND a.attnum = k.attnum)
| 		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
| 		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
| 		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname
: string
| public
: string
| plain_update
< result
< column "relname" "NAME"
< column "indisprimary" "BOOL"
< column "indisunique" "BOOL"
< column "?column?" "BOOL"
< column "string_agg" "TEXT"
< row
: string
| plain_update_pkey
: bool true
: bool true
: bool true
: string
| id
> exec
| /*
| Original SQL:
| UPDATE plain_update SET a = 10 WHERE id = 2
| */
| INSERT INTO "public"."plain_update" ("id", "a") SELECT "id", "a" FROM "bbdataarchive"."_bb__plain_update_public" ON CONFLICT ON CONSTRAINT "plain_update_pkey" DO UPDATE SET "a" = EXCLUDED."a";
< ok 1 -
> query
| SELECT x::text FROM "public"."plain_update" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
> query
| SELECT to_regclass($1)::text
: string
| canary
< result
< column "to_regclass" "TEXT"
< row
: string
| canary
> exec
| CREATE TABLE cte_t (id INT PRIMARY KEY); INSERT INTO cte_t VALUES (1), (2), (3)
< ok 3 -
> query
| SELECT x::text FROM "public"."cte_t" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1)
< row
: string
| (2)
< row
: string
| (3)
> query
| SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = $1 AND tablename = $2)
: string
| public
: string
| cte_t
< result
< column "exists" "BOOL"
< row
: bool true
> query ^
: string
| public
: string
| cte_t
< result
< column "exists" "BOOL"
< row
: bool true
> exec
| CREATE TABLE "bbdataarchive"."_bb__cte_t_public" AS
|   (WITH x AS (SELECT 1 AS id)
|   SELECT "public"."cte_t".* FROM cte_t, x WHERE cte_t.id = x.id)
|   UNION
|   (WITH x AS (SELECT 2 AS id)
|   SELECT "public"."cte_t".* FROM cte_t, x WHERE cte_t.id = x.id);
< ok 2 -
> exec
| WITH x AS (SELECT 1 AS id) DELETE FROM cte_t USING x WHERE cte_t.id = x.id;WITH x AS (SELECT 2 AS id) DELETE FROM cte_t USING x WHERE cte_t.id = x.id
< ok 1 -
> query
| SELECT x::text FROM "public"."cte_t" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (3)
> query
| SELECT a.attname, a.attgenerated <> '' FROM pg_catalog.pg_attribute a
| 		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum
: string
| public
: string
| cte_t
< result
< column "attname" "NAME"
< column "?column?" "BOOL"
< row
: string
| id
: bool false
> query
| SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL,
| 		(SELECT string_agg(a.attname, E'\x1f' ORDER BY k.n) FROM unnest(x.indkey) WITH ORDINALITY AS k(attnum, n)
| 			JOIN pg_catalog.pg_attribute a ON a.attrelid = x.indrelid AND a.attnum = k.attnum)
| 		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
| 		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
| 		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname
: string
| public
: string
| cte_t
< result
< column "relname" "NAME"
< column "indisprimary" "BOOL"
< column "indisunique" "BOOL"
< column "?column?" "BOOL"
< column "string_agg" "TEXT"
< row
: string
| cte_t_pkey
: bool true
: bool true
: bool true
: string
| id
> exec
| /*
| Original SQL:
| WITH x AS (SELECT 1 AS id) DELETE FROM cte_t USING x WHERE cte_t.id = x.id;WITH x AS (SELECT 2 AS id) DELETE FROM cte_t USING x WHERE cte_t.id = x.id
| */
| INSERT INTO "public"."cte_t" ("id") SELECT "id" FROM "bbdataarchive"."_bb__cte_t_public";
< ok 2 -
> query
| SELECT x::text FROM "public"."cte_t" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1)
< row
: string
| (2)
< row
: string
| (3)
> query
| SELECT to_regclass($1)::text
: string
| canary
< result
< column "to_regclass" "TEXT"
< row
: string
| canary
> exec
| CREATE TABLE notes (id INT PRIMARY KEY, note TEXT); INSERT INTO notes VALUES (1, '*/'), (2, '/*'), (3, 'x')
< ok 3 -
> query
| SELECT x::text FROM "public"."notes" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,*/)
< row
: string
| (2,/*)
< row
: string
| (3,x)
> query
| SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = $1 AND tablename = $2)
: string
| public
: string
| notes
< result
< column "exists" "BOOL"
< row
: bool true
> exec
| CREATE TABLE "bbdataarchive"."_bb__notes_public" AS
|   SELECT "public"."notes".* FROM notes WHERE note IN ('*/ DROP TABLE canary; /*', '/*/') OR id < 3;
< ok 2 -
> exec
| DELETE FROM notes WHERE note IN ('*/ DROP TABLE canary; /*', '/*/') OR id < 3
< ok 2 -
> query
| SELECT x::text FROM "public"."notes" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (3,x)
> query
| SELECT a.attname, a.attgenerated <> '' FROM pg_catalog.pg_attribute a
| 		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum
: string
| public
: string
| notes
< result
< column "attname" "NAME"
< column "?column?" "BOOL"
< row
: string
| id
: bool false
< row
: string
| note
: bool false
> query
| SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL,
| 		(SELECT string_agg(a.attname, E'\x1f' ORDER BY k.n) FROM unnest(x.indkey) WITH ORDINALITY AS k(attnum, n)
| 			JOIN pg_catalog.pg_attribute a ON a.attrelid = x.indrelid AND a.attnum = k.attnum)
| 		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
| 		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
| 		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname
: string
| public
: string
| notes
< result
< column "relname" "NAME"
< column "indisprimary" "BOOL"
< column "indisunique" "BOOL"
< column "?column?" "BOOL"
< column "string_agg" "TEXT"
< row
: string
| notes_pkey
: bool true
: bool true
: bool true
: string
| id
> exec
| /*
| Original SQL:
| DELETE FROM notes WHERE note IN ('* / DROP TABLE canary; / *', '/ * /') OR id < 3
| */
| INSERT INTO "public"."notes" ("id", "note") SELECT "id", "note" FROM "bbdataarchive"."_bb__notes_public";
< ok 2 -
> query
| SELECT x::text FROM "public"."notes" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,*/)
< row
: string
| (2,/*)
< row
: string
| (3,x)
> query
| SELECT to_regclass($1)::text
: string
| canary
< result
< column "to_regclass" "TEXT"
< row
: string
| canary
> exec
| CREATE TABLE alias_t (id INT PRIMARY KEY, v INT); INSERT INTO alias_t VALUES (1, 0), (2, 0), (3, 0)
< ok 3 -
> query
| SELECT x::text FROM "public"."alias_t" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,0)
< row
: string
| (2,0)
< row
: string
| (3,0)
> query
| SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = $1 AND tablename = $2)
: string
| public
: string
| alias_t
< result
< column "exists" "BOOL"
< row
: bool true
> query ^
: string
| public
: string
| alias_t
< result
< column "exists" "BOOL"
< row
: bool true
> query ^
: string
| public
: string
| alias_t
< result
< column "exists" "BOOL"
< row
: bool true
> exec
| CREATE TABLE "bbdataarchive"."_bb__alias_t_public" AS
|   SELECT "a".* FROM alias_t AS a WHERE a.id = 1
|   UNION
|   SELECT "b".* FROM alias_t AS b WHERE b.id = 2
|   UNION
|   SELECT "public"."alias_t".* FROM alias_t WHERE id = 3;
< ok 3 -
> exec
| UPDATE alias_t AS a SET v = 1 WHERE a.id = 1; UPDATE alias_t AS b SET v = 2 WHERE b.id = 2; UPDATE alias_t SET v = 3 WHERE id = 3
< ok 1 -
> query
| SELECT x::text FROM "public"."alias_t" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
< row
: string
| (3,3)
> query
| SELECT a.attname, a.attgenerated <> '' FROM pg_catalog.pg_attribute a
| 		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum
: string
| public
: string
| alias_t
< result
< column "attname" "NAME"
< column "?column?" "BOOL"
< row
: string
| id
: bool false
< row
: string
| v
: bool false
> query
| SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL,
| 		(SELECT string_agg(a.attname, E'\x1f' ORDER BY k.n) FROM unnest(x.indkey) WITH ORDINALITY AS k(attnum, n)
| 			JOIN pg_catalog.pg_attribute a ON a.attrelid = x.indrelid AND a.attnum = k.attnum)
| 		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
| 		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
| 		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname
: string
| public
: string
| alias_t
< result
< column "relname" "NAME"
< column "indisprimary" "BOOL"
< column "indisunique" "BOOL"
< column "?column?" "BOOL"
< column "string_agg" "TEXT"
< row
: string
| alias_t_pkey
: bool true
: bool true
: bool true
: string
| id
> exec
| /*
| Original SQL:
| UPDATE alias_t AS a SET v = 1 WHERE a.id = 1; UPDATE alias_t AS b SET v = 2 WHERE b.id = 2; UPDATE alias_t SET v = 3 WHERE id = 3
| */
| INSERT INTO "public"."alias_t" ("id", "v") SELECT "id", "v" FROM "bbdataarchive"."_bb__alias_t_public" ON CONFLICT ON CONSTRAINT "alias_t_pkey" DO UPDATE SET "v" = EXCLUDED."v";
< ok 3 -
> query
| SELECT x::text FROM "public"."alias_t" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,0)
< row
: string
| (2,0)
< row
: string
| (3,0)
> query
| SELECT to_regclass($1)::text
: string
| canary
< result
< column "to_regclass" "TEXT"
< row
: string
| canary
> exec
| CREATE SCHEMA "s""q"; CREATE TABLE "s""q"."t""q" ("i""d" INT CONSTRAINT "p""k" PRIMARY KEY, "a""b" INT);INSERT INTO "s""q"."t""q" VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT x::text FROM "s""q"."t""q" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
> exec
| CREATE TABLE "bbdataarchive"."_bb__t""q_s""q" AS
|   SELECT "x""y".* FROM "s""q"."t""q" AS "x""y" WHERE "x""y"."i""d" = 1;
< ok 1 -
> exec
| DELETE FROM "s""q"."t""q" AS "x""y" WHERE "x""y"."i""d" = 1
< ok 1 -
> query
| SELECT x::text FROM "s""q"."t""q" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (2,2)
> query
| SELECT a.attname, a.attgenerated <> '' FROM pg_catalog.pg_attribute a
| 		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum
: string
| s"q
: string
| t"q
< result
< column "attname" "NAME"
< column "?column?" "BOOL"
< row
: string
| i"d
: bool false
< row
: string
| a"b
: bool false
> query
| SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL,
| 		(SELECT string_agg(a.attname, E'\x1f' ORDER BY k.n) FROM unnest(x.indkey) WITH ORDINALITY AS k(attnum, n)
| 			JOIN pg_catalog.pg_attribute a ON a.attrelid = x.indrelid AND a.attnum = k.attnum)
| 		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
| 		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
| 		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname
: string
| s"q
: string
| t"q
< result
< column "relname" "NAME"
< column "indisprimary" "BOOL"
< column "indisunique" "BOOL"
< column "?column?" "BOOL"
< column "string_agg" "TEXT"
< row
: string
| p"k
: bool true
: bool true
: bool true
: string
| i"d
> exec
| /*
| Original SQL:
| DELETE FROM "s""q"."t""q" AS "x""y" WHERE "x""y"."i""d" = 1
| */
| INSERT INTO "s""q"."t""q" ("i""d", "a""b") SELECT "i""d", "a""b" FROM "bbdataarchive"."_bb__t""q_s""q";
< ok 1 -
> query
| SELECT x::text FROM "s""q"."t""q" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
> query
| SELECT to_regclass($1)::text
: string
| canary
< result
< column "to_regclass" "TEXT"
< row
: string
| canary
> exec
| CREATE SCHEMA "u""q"; CREATE TABLE "u""q"."t""q" ("i""d" INT CONSTRAINT "p""k2" PRIMARY KEY, "a""b" INT);INSERT INTO "u""q"."t""q" VALUES (1, 1), (2, 2)
< ok 2 -
> query
| SELECT x::text FROM "u""q"."t""q" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
> exec
| CREATE TABLE "bbdataarchive"."_bb__t""q_u""q" AS
|   SELECT "u""q"."t""q".* FROM "u""q"."t""q" WHERE "i""d" = 2;
< ok 1 -
> exec
| UPDATE "u""q"."t""q" SET "a""b" = 10 WHERE "i""d" = 2
< ok 1 -
> query
| SELECT x::text FROM "u""q"."t""q" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,10)
> query
| SELECT a.attname, a.attgenerated <> '' FROM pg_catalog.pg_attribute a
| 		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		WHERE n.nspname = $1 AND c.relname = $2 AND a.attnum > 0 AND NOT a.attisdropped ORDER BY a.attnum
: string
| u"q
: string
| t"q
< result
< column "attname" "NAME"
< column "?column?" "BOOL"
< row
: string
| i"d
: bool false
< row
: string
| a"b
: bool false
> query
| SELECT i.relname, x.indisprimary, x.indisunique, con.oid IS NOT NULL,
| 		(SELECT string_agg(a.attname, E'\x1f' ORDER BY k.n) FROM unnest(x.indkey) WITH ORDINALITY AS k(attnum, n)
| 			JOIN pg_catalog.pg_attribute a ON a.attrelid = x.indrelid AND a.attnum = k.attnum)
| 		FROM pg_catalog.pg_index x JOIN pg_catalog.pg_class i ON i.oid = x.indexrelid
| 		JOIN pg_catalog.pg_class c ON c.oid = x.indrelid JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		LEFT JOIN pg_catalog.pg_constraint con ON con.conindid = x.indexrelid
| 		WHERE n.nspname = $1 AND c.relname = $2 ORDER BY i.relname
: string
| u"q
: string
| t"q
< result
< column "relname" "NAME"
< column "indisprimary" "BOOL"
< column "indisunique" "BOOL"
< column "?column?" "BOOL"
< column "string_agg" "TEXT"
< row
: string
| p"k2
: bool true
: bool true
: bool true
: string
| i"d
> exec
| /*
| Original SQL:
| UPDATE "u""q"."t""q" SET "a""b" = 10 WHERE "i""d" = 2
| */
| INSERT INTO "u""q"."t""q" ("i""d", "a""b") SELECT "i""d", "a""b" FROM "bbdataarchive"."_bb__t""q_u""q" ON CONFLICT ON CONSTRAINT "p""k2" DO UPDATE SET "a""b" = EXCLUDED."a""b";
< ok 1 -
> query
| SELECT x::text FROM "u""q"."t""q" AS x ORDER BY 1
< result
< column "x" "TEXT"
< row
: string
| (1,1)
< row
: string
| (2,2)
> query
| SELECT to_regclass($1)::text
: string
| canary
< result
< column "to_regclass" "TEXT"
< row
: string
| canary
