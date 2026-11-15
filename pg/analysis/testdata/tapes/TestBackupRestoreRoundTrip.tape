-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestBackupRestoreRoundTrip default#1
> exec
| CREATE SCHEMA bbdataarchive
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
