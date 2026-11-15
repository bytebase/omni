-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12
-- pin: pg12 postgres:12.22-alpine sha256:7c8f4870583184ebadf7f17a6513620aac5f365a7938dc6a6911c1d5df2f481a REL_12_22
-- pin: pg13 postgres:13.23-alpine sha256:fb9065b6e3e213bdc07edd372a5b2a26245840b7fb65d1fd8b6700106d51805c REL_13_23
-- pin: pg14 postgres:14.25-alpine sha256:14078c4aa0db8e2e8492866138cafd796fc39afd5ecb5b1c6f26a912b07529f5 REL_14_25
-- pin: pg15 postgres:15.20-alpine sha256:5cb339462bde5e6b617acedf6c2a3008f3aadc44a3c0507dd22226ee1279aee8 REL_15_20
-- pin: pg16 postgres:16.16-alpine sha256:0adbf25a9e96d38823bb813f3ac0ed41fadb665c88f1df7fe6dbfb060cbdc106 REL_16_16
-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestQuery/default default#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE SCHEMA "app""schema";
| 				CREATE TABLE public.lookup_precedence (marker text);
| 				CREATE TABLE app.lookup_precedence (marker text);
| 				CREATE TABLE "app""schema".lookup_precedence (marker text);
| 				INSERT INTO public.lookup_precedence VALUES ('public');
| 				INSERT INTO app.lookup_precedence VALUES ('app');
| 				INSERT INTO "app""schema".lookup_precedence VALUES ('quoted');
| 				CREATE FUNCTION public.helper_echo(rel text) RETURNS text LANGUAGE sql AS $$ SELECT rel $$;
| 				CREATE FUNCTION public.call_unqualified_helper(rel text) RETURNS text LANGUAGE plpgsql AS $$
| 				BEGIN
| 					RETURN helper_echo(rel);
| 				END; $$;
| 				CREATE TABLE plan_target (id int PRIMARY KEY, note text);
| 				INSERT INTO plan_target SELECT g, 'n' || g FROM generate_series(1, 200) g;
| 				ANALYZE plan_target;
| 			
< ok 0 -
> query
| SELECT
| 				true AS b,
| 				42::int4 AS i4,
| 				9007199254740993::int8 AS i8,
| 				3.5::float8 AS f8,
| 				12.340::numeric(6,3) AS num,
| 				'hello'::text AS txt,
| 				'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid AS u,
| 				'\xdeadbeef'::bytea AS bin,
| 				''::bytea AS empty_bin,
| 				'2026-03-04'::date AS dt,
| 				'infinity'::date AS dt_inf,
| 				'2026-03-04 05:06:07.891234'::timestamp AS ts,
| 				'-infinity'::timestamp AS ts_neginf,
| 				'2026-03-04 05:06:07.891234+00'::timestamptz AS tstz,
| 				'infinity'::timestamptz AS tstz_inf,
| 				'12:34:56.12'::time AS tm,
| 				'-00:02:45.25'::interval AS iv,
| 				'{"k": 1}'::jsonb AS jb,
| 				NULL::text AS nul;
< result
< column "b" "BOOL"
< column "i4" "INT4"
< column "i8" "INT8"
< column "f8" "FLOAT8"
< column "num" "NUMERIC"
< column "txt" "TEXT"
< column "u" "UUID"
< column "bin" "BYTEA"
< column "empty_bin" "BYTEA"
< column "dt" "DATE"
< column "dt_inf" "DATE"
< column "ts" "TIMESTAMP"
< column "ts_neginf" "TIMESTAMP"
< column "tstz" "TIMESTAMPTZ"
< column "tstz_inf" "TIMESTAMPTZ"
< column "tm" "TIME"
< column "iv" "INTERVAL"
< column "jb" "JSONB"
< column "nul" "TEXT"
< row
: bool true
: int64 42
: int64 9007199254740993
: float64 3.5
: string
| 12.340
: string
| hello
: string
| a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
: bytes
|q "ޭ\xbe\xef"
: bytes
|
: time 2026-03-04T00:00:00Z
: string
| infinity
: time 2026-03-04T05:06:07.891234Z
: string
| -infinity
: time 2026-03-04T05:06:07.891234Z
: string
| infinity
: string
| 12:34:56.12
: string
| -00:02:45.25
: bytes
| {"k": 1}
: null
> exec
|
| 				CREATE TEMP TABLE sess_items (id int PRIMARY KEY, note text);
< ok 0 -
> exec
|
| 				INSERT INTO sess_items SELECT g, repeat('x', 40) FROM generate_series(1, 10) g;
< ok 10 -
> query
|
| 				SELECT FROM sess_items;
< result
< row
< row
< row
< row
< row
< row
< row
< row
< row
< row
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 2;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
> query
| WITH result AS (
| SELECT id FROM sess_items ORDER BY id LIMIT ALL
| ) SELECT * FROM result LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT id, note FROM sess_items ORDER BY id;
< result
< column "id" "INT4"
< column "note" "TEXT"
< row
: int64 1
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 4
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 5
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 6
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 7
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 8
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 9
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 10
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7 AS after_trunc;
< result
< column "after_trunc" "INT4"
< row
: int64 7
> query
| SELECT 1;
< result
< column "?column?" "INT4"
< row
: int64 1
> query
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT public.call_unqualified_helper('ok'::text);
< result
< column "call_unqualified_helper" "TEXT"
< row
: string
| ok
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| app
> exec
| SET search_path TO "app""schema", public;
< ok 0 -
> query
| SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| quoted
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on plan_target  (cost=0.00..3.50 rows=150 width=4)
< row
: string
|   Filter: (id > 50)
> query
| EXPLAIN (FORMAT JSON) SELECT note FROM plan_target WHERE id > 50
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "plan_target",
|       "Alias": "plan_target",
|       "Startup Cost": 0.00,
|       "Total Cost": 3.50,
|       "Plan Rows": 150,
|       "Plan Width": 4,
|       "Filter": "(id > 50)"
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT YAML) SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| - Plan: 
|     Node Type: "Seq Scan"
|     Parallel Aware: false
|     Async Capable: false
|     Relation Name: "plan_target"
|     Alias: "plan_target"
|     Startup Cost: 0.00
|     Total Cost: 3.50
|     Plan Rows: 150
|     Plan Width: 4
|     Filter: "(id > 50)"

== TestQuery/pg12 pg12#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE SCHEMA "app""schema";
| 				CREATE TABLE public.lookup_precedence (marker text);
| 				CREATE TABLE app.lookup_precedence (marker text);
| 				CREATE TABLE "app""schema".lookup_precedence (marker text);
| 				INSERT INTO public.lookup_precedence VALUES ('public');
| 				INSERT INTO app.lookup_precedence VALUES ('app');
| 				INSERT INTO "app""schema".lookup_precedence VALUES ('quoted');
| 				CREATE FUNCTION public.helper_echo(rel text) RETURNS text LANGUAGE sql AS $$ SELECT rel $$;
| 				CREATE FUNCTION public.call_unqualified_helper(rel text) RETURNS text LANGUAGE plpgsql AS $$
| 				BEGIN
| 					RETURN helper_echo(rel);
| 				END; $$;
| 				CREATE TABLE plan_target (id int PRIMARY KEY, note text);
| 				INSERT INTO plan_target SELECT g, 'n' || g FROM generate_series(1, 200) g;
| 				ANALYZE plan_target;
| 			
< ok 0 -
> query
| SELECT
| 				true AS b,
| 				42::int4 AS i4,
| 				9007199254740993::int8 AS i8,
| 				3.5::float8 AS f8,
| 				12.340::numeric(6,3) AS num,
| 				'hello'::text AS txt,
| 				'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid AS u,
| 				'\xdeadbeef'::bytea AS bin,
| 				''::bytea AS empty_bin,
| 				'2026-03-04'::date AS dt,
| 				'infinity'::date AS dt_inf,
| 				'2026-03-04 05:06:07.891234'::timestamp AS ts,
| 				'-infinity'::timestamp AS ts_neginf,
| 				'2026-03-04 05:06:07.891234+00'::timestamptz AS tstz,
| 				'infinity'::timestamptz AS tstz_inf,
| 				'12:34:56.12'::time AS tm,
| 				'-00:02:45.25'::interval AS iv,
| 				'{"k": 1}'::jsonb AS jb,
| 				NULL::text AS nul;
< result
< column "b" "BOOL"
< column "i4" "INT4"
< column "i8" "INT8"
< column "f8" "FLOAT8"
< column "num" "NUMERIC"
< column "txt" "TEXT"
< column "u" "UUID"
< column "bin" "BYTEA"
< column "empty_bin" "BYTEA"
< column "dt" "DATE"
< column "dt_inf" "DATE"
< column "ts" "TIMESTAMP"
< column "ts_neginf" "TIMESTAMP"
< column "tstz" "TIMESTAMPTZ"
< column "tstz_inf" "TIMESTAMPTZ"
< column "tm" "TIME"
< column "iv" "INTERVAL"
< column "jb" "JSONB"
< column "nul" "TEXT"
< row
: bool true
: int64 42
: int64 9007199254740993
: float64 3.5
: string
| 12.340
: string
| hello
: string
| a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
: bytes
|q "ޭ\xbe\xef"
: bytes
|
: time 2026-03-04T00:00:00Z
: string
| infinity
: time 2026-03-04T05:06:07.891234Z
: string
| -infinity
: time 2026-03-04T05:06:07.891234Z
: string
| infinity
: string
| 12:34:56.12
: string
| -00:02:45.25
: bytes
| {"k": 1}
: null
> exec
|
| 				CREATE TEMP TABLE sess_items (id int PRIMARY KEY, note text);
< ok 0 -
> exec
|
| 				INSERT INTO sess_items SELECT g, repeat('x', 40) FROM generate_series(1, 10) g;
< ok 10 -
> query
|
| 				SELECT FROM sess_items;
< result
< row
< row
< row
< row
< row
< row
< row
< row
< row
< row
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 2;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
> query
| WITH result AS (
| SELECT id FROM sess_items ORDER BY id LIMIT ALL
| ) SELECT * FROM result LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT id, note FROM sess_items ORDER BY id;
< result
< column "id" "INT4"
< column "note" "TEXT"
< row
: int64 1
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 4
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 5
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 6
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 7
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 8
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 9
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 10
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7 AS after_trunc;
< result
< column "after_trunc" "INT4"
< row
: int64 7
> query
| SELECT 1;
< result
< column "?column?" "INT4"
< row
: int64 1
> query
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT public.call_unqualified_helper('ok'::text);
< result
< column "call_unqualified_helper" "TEXT"
< row
: string
| ok
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| app
> exec
| SET search_path TO "app""schema", public;
< ok 0 -
> query
| SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| quoted
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on plan_target  (cost=0.00..3.50 rows=150 width=4)
< row
: string
|   Filter: (id > 50)
> query
| EXPLAIN (FORMAT JSON) SELECT note FROM plan_target WHERE id > 50
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Relation Name": "plan_target",
|       "Alias": "plan_target",
|       "Startup Cost": 0.00,
|       "Total Cost": 3.50,
|       "Plan Rows": 150,
|       "Plan Width": 4,
|       "Filter": "(id > 50)"
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT YAML) SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| - Plan: 
|     Node Type: "Seq Scan"
|     Parallel Aware: false
|     Relation Name: "plan_target"
|     Alias: "plan_target"
|     Startup Cost: 0.00
|     Total Cost: 3.50
|     Plan Rows: 150
|     Plan Width: 4
|     Filter: "(id > 50)"

== TestQuery/pg13 pg13#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE SCHEMA "app""schema";
| 				CREATE TABLE public.lookup_precedence (marker text);
| 				CREATE TABLE app.lookup_precedence (marker text);
| 				CREATE TABLE "app""schema".lookup_precedence (marker text);
| 				INSERT INTO public.lookup_precedence VALUES ('public');
| 				INSERT INTO app.lookup_precedence VALUES ('app');
| 				INSERT INTO "app""schema".lookup_precedence VALUES ('quoted');
| 				CREATE FUNCTION public.helper_echo(rel text) RETURNS text LANGUAGE sql AS $$ SELECT rel $$;
| 				CREATE FUNCTION public.call_unqualified_helper(rel text) RETURNS text LANGUAGE plpgsql AS $$
| 				BEGIN
| 					RETURN helper_echo(rel);
| 				END; $$;
| 				CREATE TABLE plan_target (id int PRIMARY KEY, note text);
| 				INSERT INTO plan_target SELECT g, 'n' || g FROM generate_series(1, 200) g;
| 				ANALYZE plan_target;
| 			
< ok 0 -
> query
| SELECT
| 				true AS b,
| 				42::int4 AS i4,
| 				9007199254740993::int8 AS i8,
| 				3.5::float8 AS f8,
| 				12.340::numeric(6,3) AS num,
| 				'hello'::text AS txt,
| 				'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid AS u,
| 				'\xdeadbeef'::bytea AS bin,
| 				''::bytea AS empty_bin,
| 				'2026-03-04'::date AS dt,
| 				'infinity'::date AS dt_inf,
| 				'2026-03-04 05:06:07.891234'::timestamp AS ts,
| 				'-infinity'::timestamp AS ts_neginf,
| 				'2026-03-04 05:06:07.891234+00'::timestamptz AS tstz,
| 				'infinity'::timestamptz AS tstz_inf,
| 				'12:34:56.12'::time AS tm,
| 				'-00:02:45.25'::interval AS iv,
| 				'{"k": 1}'::jsonb AS jb,
| 				NULL::text AS nul;
< result
< column "b" "BOOL"
< column "i4" "INT4"
< column "i8" "INT8"
< column "f8" "FLOAT8"
< column "num" "NUMERIC"
< column "txt" "TEXT"
< column "u" "UUID"
< column "bin" "BYTEA"
< column "empty_bin" "BYTEA"
< column "dt" "DATE"
< column "dt_inf" "DATE"
< column "ts" "TIMESTAMP"
< column "ts_neginf" "TIMESTAMP"
< column "tstz" "TIMESTAMPTZ"
< column "tstz_inf" "TIMESTAMPTZ"
< column "tm" "TIME"
< column "iv" "INTERVAL"
< column "jb" "JSONB"
< column "nul" "TEXT"
< row
: bool true
: int64 42
: int64 9007199254740993
: float64 3.5
: string
| 12.340
: string
| hello
: string
| a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
: bytes
|q "ޭ\xbe\xef"
: bytes
|
: time 2026-03-04T00:00:00Z
: string
| infinity
: time 2026-03-04T05:06:07.891234Z
: string
| -infinity
: time 2026-03-04T05:06:07.891234Z
: string
| infinity
: string
| 12:34:56.12
: string
| -00:02:45.25
: bytes
| {"k": 1}
: null
> exec
|
| 				CREATE TEMP TABLE sess_items (id int PRIMARY KEY, note text);
< ok 0 -
> exec
|
| 				INSERT INTO sess_items SELECT g, repeat('x', 40) FROM generate_series(1, 10) g;
< ok 10 -
> query
|
| 				SELECT FROM sess_items;
< result
< row
< row
< row
< row
< row
< row
< row
< row
< row
< row
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 2;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
> query
| WITH result AS (
| SELECT id FROM sess_items ORDER BY id LIMIT ALL
| ) SELECT * FROM result LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT id, note FROM sess_items ORDER BY id;
< result
< column "id" "INT4"
< column "note" "TEXT"
< row
: int64 1
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 4
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 5
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 6
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 7
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 8
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 9
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 10
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7 AS after_trunc;
< result
< column "after_trunc" "INT4"
< row
: int64 7
> query
| SELECT 1;
< result
< column "?column?" "INT4"
< row
: int64 1
> query
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT public.call_unqualified_helper('ok'::text);
< result
< column "call_unqualified_helper" "TEXT"
< row
: string
| ok
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| app
> exec
| SET search_path TO "app""schema", public;
< ok 0 -
> query
| SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| quoted
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on plan_target  (cost=0.00..3.50 rows=150 width=4)
< row
: string
|   Filter: (id > 50)
> query
| EXPLAIN (FORMAT JSON) SELECT note FROM plan_target WHERE id > 50
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Relation Name": "plan_target",
|       "Alias": "plan_target",
|       "Startup Cost": 0.00,
|       "Total Cost": 3.50,
|       "Plan Rows": 150,
|       "Plan Width": 4,
|       "Filter": "(id > 50)"
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT YAML) SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| - Plan: 
|     Node Type: "Seq Scan"
|     Parallel Aware: false
|     Relation Name: "plan_target"
|     Alias: "plan_target"
|     Startup Cost: 0.00
|     Total Cost: 3.50
|     Plan Rows: 150
|     Plan Width: 4
|     Filter: "(id > 50)"

== TestQuery/pg14 pg14#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE SCHEMA "app""schema";
| 				CREATE TABLE public.lookup_precedence (marker text);
| 				CREATE TABLE app.lookup_precedence (marker text);
| 				CREATE TABLE "app""schema".lookup_precedence (marker text);
| 				INSERT INTO public.lookup_precedence VALUES ('public');
| 				INSERT INTO app.lookup_precedence VALUES ('app');
| 				INSERT INTO "app""schema".lookup_precedence VALUES ('quoted');
| 				CREATE FUNCTION public.helper_echo(rel text) RETURNS text LANGUAGE sql AS $$ SELECT rel $$;
| 				CREATE FUNCTION public.call_unqualified_helper(rel text) RETURNS text LANGUAGE plpgsql AS $$
| 				BEGIN
| 					RETURN helper_echo(rel);
| 				END; $$;
| 				CREATE TABLE plan_target (id int PRIMARY KEY, note text);
| 				INSERT INTO plan_target SELECT g, 'n' || g FROM generate_series(1, 200) g;
| 				ANALYZE plan_target;
| 			
< ok 0 -
> query
| SELECT
| 				true AS b,
| 				42::int4 AS i4,
| 				9007199254740993::int8 AS i8,
| 				3.5::float8 AS f8,
| 				12.340::numeric(6,3) AS num,
| 				'hello'::text AS txt,
| 				'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid AS u,
| 				'\xdeadbeef'::bytea AS bin,
| 				''::bytea AS empty_bin,
| 				'2026-03-04'::date AS dt,
| 				'infinity'::date AS dt_inf,
| 				'2026-03-04 05:06:07.891234'::timestamp AS ts,
| 				'-infinity'::timestamp AS ts_neginf,
| 				'2026-03-04 05:06:07.891234+00'::timestamptz AS tstz,
| 				'infinity'::timestamptz AS tstz_inf,
| 				'12:34:56.12'::time AS tm,
| 				'-00:02:45.25'::interval AS iv,
| 				'{"k": 1}'::jsonb AS jb,
| 				NULL::text AS nul;
< result
< column "b" "BOOL"
< column "i4" "INT4"
< column "i8" "INT8"
< column "f8" "FLOAT8"
< column "num" "NUMERIC"
< column "txt" "TEXT"
< column "u" "UUID"
< column "bin" "BYTEA"
< column "empty_bin" "BYTEA"
< column "dt" "DATE"
< column "dt_inf" "DATE"
< column "ts" "TIMESTAMP"
< column "ts_neginf" "TIMESTAMP"
< column "tstz" "TIMESTAMPTZ"
< column "tstz_inf" "TIMESTAMPTZ"
< column "tm" "TIME"
< column "iv" "INTERVAL"
< column "jb" "JSONB"
< column "nul" "TEXT"
< row
: bool true
: int64 42
: int64 9007199254740993
: float64 3.5
: string
| 12.340
: string
| hello
: string
| a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
: bytes
|q "ޭ\xbe\xef"
: bytes
|
: time 2026-03-04T00:00:00Z
: string
| infinity
: time 2026-03-04T05:06:07.891234Z
: string
| -infinity
: time 2026-03-04T05:06:07.891234Z
: string
| infinity
: string
| 12:34:56.12
: string
| -00:02:45.25
: bytes
| {"k": 1}
: null
> exec
|
| 				CREATE TEMP TABLE sess_items (id int PRIMARY KEY, note text);
< ok 0 -
> exec
|
| 				INSERT INTO sess_items SELECT g, repeat('x', 40) FROM generate_series(1, 10) g;
< ok 10 -
> query
|
| 				SELECT FROM sess_items;
< result
< row
< row
< row
< row
< row
< row
< row
< row
< row
< row
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 2;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
> query
| WITH result AS (
| SELECT id FROM sess_items ORDER BY id LIMIT ALL
| ) SELECT * FROM result LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT id, note FROM sess_items ORDER BY id;
< result
< column "id" "INT4"
< column "note" "TEXT"
< row
: int64 1
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 4
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 5
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 6
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 7
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 8
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 9
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 10
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7 AS after_trunc;
< result
< column "after_trunc" "INT4"
< row
: int64 7
> query
| SELECT 1;
< result
< column "?column?" "INT4"
< row
: int64 1
> query
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT public.call_unqualified_helper('ok'::text);
< result
< column "call_unqualified_helper" "TEXT"
< row
: string
| ok
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| app
> exec
| SET search_path TO "app""schema", public;
< ok 0 -
> query
| SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| quoted
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on plan_target  (cost=0.00..3.50 rows=150 width=4)
< row
: string
|   Filter: (id > 50)
> query
| EXPLAIN (FORMAT JSON) SELECT note FROM plan_target WHERE id > 50
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "plan_target",
|       "Alias": "plan_target",
|       "Startup Cost": 0.00,
|       "Total Cost": 3.50,
|       "Plan Rows": 150,
|       "Plan Width": 4,
|       "Filter": "(id > 50)"
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT YAML) SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| - Plan: 
|     Node Type: "Seq Scan"
|     Parallel Aware: false
|     Async Capable: false
|     Relation Name: "plan_target"
|     Alias: "plan_target"
|     Startup Cost: 0.00
|     Total Cost: 3.50
|     Plan Rows: 150
|     Plan Width: 4
|     Filter: "(id > 50)"

== TestQuery/pg15 pg15#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE SCHEMA "app""schema";
| 				CREATE TABLE public.lookup_precedence (marker text);
| 				CREATE TABLE app.lookup_precedence (marker text);
| 				CREATE TABLE "app""schema".lookup_precedence (marker text);
| 				INSERT INTO public.lookup_precedence VALUES ('public');
| 				INSERT INTO app.lookup_precedence VALUES ('app');
| 				INSERT INTO "app""schema".lookup_precedence VALUES ('quoted');
| 				CREATE FUNCTION public.helper_echo(rel text) RETURNS text LANGUAGE sql AS $$ SELECT rel $$;
| 				CREATE FUNCTION public.call_unqualified_helper(rel text) RETURNS text LANGUAGE plpgsql AS $$
| 				BEGIN
| 					RETURN helper_echo(rel);
| 				END; $$;
| 				CREATE TABLE plan_target (id int PRIMARY KEY, note text);
| 				INSERT INTO plan_target SELECT g, 'n' || g FROM generate_series(1, 200) g;
| 				ANALYZE plan_target;
| 			
< ok 0 -
> query
| SELECT
| 				true AS b,
| 				42::int4 AS i4,
| 				9007199254740993::int8 AS i8,
| 				3.5::float8 AS f8,
| 				12.340::numeric(6,3) AS num,
| 				'hello'::text AS txt,
| 				'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid AS u,
| 				'\xdeadbeef'::bytea AS bin,
| 				''::bytea AS empty_bin,
| 				'2026-03-04'::date AS dt,
| 				'infinity'::date AS dt_inf,
| 				'2026-03-04 05:06:07.891234'::timestamp AS ts,
| 				'-infinity'::timestamp AS ts_neginf,
| 				'2026-03-04 05:06:07.891234+00'::timestamptz AS tstz,
| 				'infinity'::timestamptz AS tstz_inf,
| 				'12:34:56.12'::time AS tm,
| 				'-00:02:45.25'::interval AS iv,
| 				'{"k": 1}'::jsonb AS jb,
| 				NULL::text AS nul;
< result
< column "b" "BOOL"
< column "i4" "INT4"
< column "i8" "INT8"
< column "f8" "FLOAT8"
< column "num" "NUMERIC"
< column "txt" "TEXT"
< column "u" "UUID"
< column "bin" "BYTEA"
< column "empty_bin" "BYTEA"
< column "dt" "DATE"
< column "dt_inf" "DATE"
< column "ts" "TIMESTAMP"
< column "ts_neginf" "TIMESTAMP"
< column "tstz" "TIMESTAMPTZ"
< column "tstz_inf" "TIMESTAMPTZ"
< column "tm" "TIME"
< column "iv" "INTERVAL"
< column "jb" "JSONB"
< column "nul" "TEXT"
< row
: bool true
: int64 42
: int64 9007199254740993
: float64 3.5
: string
| 12.340
: string
| hello
: string
| a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
: bytes
|q "ޭ\xbe\xef"
: bytes
|
: time 2026-03-04T00:00:00Z
: string
| infinity
: time 2026-03-04T05:06:07.891234Z
: string
| -infinity
: time 2026-03-04T05:06:07.891234Z
: string
| infinity
: string
| 12:34:56.12
: string
| -00:02:45.25
: bytes
| {"k": 1}
: null
> exec
|
| 				CREATE TEMP TABLE sess_items (id int PRIMARY KEY, note text);
< ok 0 -
> exec
|
| 				INSERT INTO sess_items SELECT g, repeat('x', 40) FROM generate_series(1, 10) g;
< ok 10 -
> query
|
| 				SELECT FROM sess_items;
< result
< row
< row
< row
< row
< row
< row
< row
< row
< row
< row
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 2;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
> query
| WITH result AS (
| SELECT id FROM sess_items ORDER BY id LIMIT ALL
| ) SELECT * FROM result LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT id, note FROM sess_items ORDER BY id;
< result
< column "id" "INT4"
< column "note" "TEXT"
< row
: int64 1
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 4
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 5
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 6
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 7
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 8
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 9
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 10
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7 AS after_trunc;
< result
< column "after_trunc" "INT4"
< row
: int64 7
> query
| SELECT 1;
< result
< column "?column?" "INT4"
< row
: int64 1
> query
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT public.call_unqualified_helper('ok'::text);
< result
< column "call_unqualified_helper" "TEXT"
< row
: string
| ok
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| app
> exec
| SET search_path TO "app""schema", public;
< ok 0 -
> query
| SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| quoted
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on plan_target  (cost=0.00..3.50 rows=150 width=4)
< row
: string
|   Filter: (id > 50)
> query
| EXPLAIN (FORMAT JSON) SELECT note FROM plan_target WHERE id > 50
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "plan_target",
|       "Alias": "plan_target",
|       "Startup Cost": 0.00,
|       "Total Cost": 3.50,
|       "Plan Rows": 150,
|       "Plan Width": 4,
|       "Filter": "(id > 50)"
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT YAML) SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| - Plan: 
|     Node Type: "Seq Scan"
|     Parallel Aware: false
|     Async Capable: false
|     Relation Name: "plan_target"
|     Alias: "plan_target"
|     Startup Cost: 0.00
|     Total Cost: 3.50
|     Plan Rows: 150
|     Plan Width: 4
|     Filter: "(id > 50)"

== TestQuery/pg16 pg16#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE SCHEMA "app""schema";
| 				CREATE TABLE public.lookup_precedence (marker text);
| 				CREATE TABLE app.lookup_precedence (marker text);
| 				CREATE TABLE "app""schema".lookup_precedence (marker text);
| 				INSERT INTO public.lookup_precedence VALUES ('public');
| 				INSERT INTO app.lookup_precedence VALUES ('app');
| 				INSERT INTO "app""schema".lookup_precedence VALUES ('quoted');
| 				CREATE FUNCTION public.helper_echo(rel text) RETURNS text LANGUAGE sql AS $$ SELECT rel $$;
| 				CREATE FUNCTION public.call_unqualified_helper(rel text) RETURNS text LANGUAGE plpgsql AS $$
| 				BEGIN
| 					RETURN helper_echo(rel);
| 				END; $$;
| 				CREATE TABLE plan_target (id int PRIMARY KEY, note text);
| 				INSERT INTO plan_target SELECT g, 'n' || g FROM generate_series(1, 200) g;
| 				ANALYZE plan_target;
| 			
< ok 0 -
> query
| SELECT
| 				true AS b,
| 				42::int4 AS i4,
| 				9007199254740993::int8 AS i8,
| 				3.5::float8 AS f8,
| 				12.340::numeric(6,3) AS num,
| 				'hello'::text AS txt,
| 				'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid AS u,
| 				'\xdeadbeef'::bytea AS bin,
| 				''::bytea AS empty_bin,
| 				'2026-03-04'::date AS dt,
| 				'infinity'::date AS dt_inf,
| 				'2026-03-04 05:06:07.891234'::timestamp AS ts,
| 				'-infinity'::timestamp AS ts_neginf,
| 				'2026-03-04 05:06:07.891234+00'::timestamptz AS tstz,
| 				'infinity'::timestamptz AS tstz_inf,
| 				'12:34:56.12'::time AS tm,
| 				'-00:02:45.25'::interval AS iv,
| 				'{"k": 1}'::jsonb AS jb,
| 				NULL::text AS nul;
< result
< column "b" "BOOL"
< column "i4" "INT4"
< column "i8" "INT8"
< column "f8" "FLOAT8"
< column "num" "NUMERIC"
< column "txt" "TEXT"
< column "u" "UUID"
< column "bin" "BYTEA"
< column "empty_bin" "BYTEA"
< column "dt" "DATE"
< column "dt_inf" "DATE"
< column "ts" "TIMESTAMP"
< column "ts_neginf" "TIMESTAMP"
< column "tstz" "TIMESTAMPTZ"
< column "tstz_inf" "TIMESTAMPTZ"
< column "tm" "TIME"
< column "iv" "INTERVAL"
< column "jb" "JSONB"
< column "nul" "TEXT"
< row
: bool true
: int64 42
: int64 9007199254740993
: float64 3.5
: string
| 12.340
: string
| hello
: string
| a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
: bytes
|q "ޭ\xbe\xef"
: bytes
|
: time 2026-03-04T00:00:00Z
: string
| infinity
: time 2026-03-04T05:06:07.891234Z
: string
| -infinity
: time 2026-03-04T05:06:07.891234Z
: string
| infinity
: string
| 12:34:56.12
: string
| -00:02:45.25
: bytes
| {"k": 1}
: null
> exec
|
| 				CREATE TEMP TABLE sess_items (id int PRIMARY KEY, note text);
< ok 0 -
> exec
|
| 				INSERT INTO sess_items SELECT g, repeat('x', 40) FROM generate_series(1, 10) g;
< ok 10 -
> query
|
| 				SELECT FROM sess_items;
< result
< row
< row
< row
< row
< row
< row
< row
< row
< row
< row
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 2;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
> query
| WITH result AS (
| SELECT id FROM sess_items ORDER BY id LIMIT ALL
| ) SELECT * FROM result LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT id, note FROM sess_items ORDER BY id;
< result
< column "id" "INT4"
< column "note" "TEXT"
< row
: int64 1
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 4
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 5
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 6
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 7
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 8
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 9
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 10
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7 AS after_trunc;
< result
< column "after_trunc" "INT4"
< row
: int64 7
> query
| SELECT 1;
< result
< column "?column?" "INT4"
< row
: int64 1
> query
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT public.call_unqualified_helper('ok'::text);
< result
< column "call_unqualified_helper" "TEXT"
< row
: string
| ok
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| app
> exec
| SET search_path TO "app""schema", public;
< ok 0 -
> query
| SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| quoted
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on plan_target  (cost=0.00..3.50 rows=150 width=4)
< row
: string
|   Filter: (id > 50)
> query
| EXPLAIN (FORMAT JSON) SELECT note FROM plan_target WHERE id > 50
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "plan_target",
|       "Alias": "plan_target",
|       "Startup Cost": 0.00,
|       "Total Cost": 3.50,
|       "Plan Rows": 150,
|       "Plan Width": 4,
|       "Filter": "(id > 50)"
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT YAML) SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| - Plan: 
|     Node Type: "Seq Scan"
|     Parallel Aware: false
|     Async Capable: false
|     Relation Name: "plan_target"
|     Alias: "plan_target"
|     Startup Cost: 0.00
|     Total Cost: 3.50
|     Plan Rows: 150
|     Plan Width: 4
|     Filter: "(id > 50)"

== TestQuery/pg18 pg18#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE SCHEMA "app""schema";
| 				CREATE TABLE public.lookup_precedence (marker text);
| 				CREATE TABLE app.lookup_precedence (marker text);
| 				CREATE TABLE "app""schema".lookup_precedence (marker text);
| 				INSERT INTO public.lookup_precedence VALUES ('public');
| 				INSERT INTO app.lookup_precedence VALUES ('app');
| 				INSERT INTO "app""schema".lookup_precedence VALUES ('quoted');
| 				CREATE FUNCTION public.helper_echo(rel text) RETURNS text LANGUAGE sql AS $$ SELECT rel $$;
| 				CREATE FUNCTION public.call_unqualified_helper(rel text) RETURNS text LANGUAGE plpgsql AS $$
| 				BEGIN
| 					RETURN helper_echo(rel);
| 				END; $$;
| 				CREATE TABLE plan_target (id int PRIMARY KEY, note text);
| 				INSERT INTO plan_target SELECT g, 'n' || g FROM generate_series(1, 200) g;
| 				ANALYZE plan_target;
| 			
< ok 0 -
> query
| SELECT
| 				true AS b,
| 				42::int4 AS i4,
| 				9007199254740993::int8 AS i8,
| 				3.5::float8 AS f8,
| 				12.340::numeric(6,3) AS num,
| 				'hello'::text AS txt,
| 				'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid AS u,
| 				'\xdeadbeef'::bytea AS bin,
| 				''::bytea AS empty_bin,
| 				'2026-03-04'::date AS dt,
| 				'infinity'::date AS dt_inf,
| 				'2026-03-04 05:06:07.891234'::timestamp AS ts,
| 				'-infinity'::timestamp AS ts_neginf,
| 				'2026-03-04 05:06:07.891234+00'::timestamptz AS tstz,
| 				'infinity'::timestamptz AS tstz_inf,
| 				'12:34:56.12'::time AS tm,
| 				'-00:02:45.25'::interval AS iv,
| 				'{"k": 1}'::jsonb AS jb,
| 				NULL::text AS nul;
< result
< column "b" "BOOL"
< column "i4" "INT4"
< column "i8" "INT8"
< column "f8" "FLOAT8"
< column "num" "NUMERIC"
< column "txt" "TEXT"
< column "u" "UUID"
< column "bin" "BYTEA"
< column "empty_bin" "BYTEA"
< column "dt" "DATE"
< column "dt_inf" "DATE"
< column "ts" "TIMESTAMP"
< column "ts_neginf" "TIMESTAMP"
< column "tstz" "TIMESTAMPTZ"
< column "tstz_inf" "TIMESTAMPTZ"
< column "tm" "TIME"
< column "iv" "INTERVAL"
< column "jb" "JSONB"
< column "nul" "TEXT"
< row
: bool true
: int64 42
: int64 9007199254740993
: float64 3.5
: string
| 12.340
: string
| hello
: string
| a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
: bytes
|q "ޭ\xbe\xef"
: bytes
|
: time 2026-03-04T00:00:00Z
: string
| infinity
: time 2026-03-04T05:06:07.891234Z
: string
| -infinity
: time 2026-03-04T05:06:07.891234Z
: string
| infinity
: string
| 12:34:56.12
: string
| -00:02:45.25
: bytes
| {"k": 1}
: null
> exec
|
| 				CREATE TEMP TABLE sess_items (id int PRIMARY KEY, note text);
< ok 0 -
> exec
|
| 				INSERT INTO sess_items SELECT g, repeat('x', 40) FROM generate_series(1, 10) g;
< ok 10 -
> query
|
| 				SELECT FROM sess_items;
< result
< row
< row
< row
< row
< row
< row
< row
< row
< row
< row
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
|
| 				SELECT id FROM sess_items ORDER BY id LIMIT 2;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
> query
| WITH result AS (
| SELECT id FROM sess_items ORDER BY id LIMIT ALL
| ) SELECT * FROM result LIMIT 4;
< result
< column "id" "INT4"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT id, note FROM sess_items ORDER BY id;
< result
< column "id" "INT4"
< column "note" "TEXT"
< row
: int64 1
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 2
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 3
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 4
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 5
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 6
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 7
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 8
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 9
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
< row
: int64 10
: string
| xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
> query
|
| SELECT 7 AS after_trunc;
< result
< column "after_trunc" "INT4"
< row
: int64 7
> query
| SELECT 1;
< result
< column "?column?" "INT4"
< row
: int64 1
> query
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT public.call_unqualified_helper('ok'::text);
< result
< column "call_unqualified_helper" "TEXT"
< row
: string
| ok
> exec
| SET search_path TO "app", public;
< ok 0 -
> query
|
| 				SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| app
> exec
| SET search_path TO "app""schema", public;
< ok 0 -
> query
| SELECT marker FROM lookup_precedence;
< result
< column "marker" "TEXT"
< row
: string
| quoted
> query
| EXPLAIN SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on plan_target  (cost=0.00..3.50 rows=150 width=4)
< row
: string
|   Filter: (id > 50)
> query
| EXPLAIN (FORMAT JSON) SELECT note FROM plan_target WHERE id > 50
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "plan_target",
|       "Alias": "plan_target",
|       "Startup Cost": 0.00,
|       "Total Cost": 3.50,
|       "Plan Rows": 150,
|       "Plan Width": 4,
|       "Disabled": false,
|       "Filter": "(id > 50)"
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT YAML) SELECT note FROM plan_target WHERE id > 50;
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| - Plan: 
|     Node Type: "Seq Scan"
|     Parallel Aware: false
|     Async Capable: false
|     Relation Name: "plan_target"
|     Alias: "plan_target"
|     Startup Cost: 0.00
|     Total Cost: 3.50
|     Plan Rows: 150
|     Plan Width: 4
|     Disabled: false
|     Filter: "(id > 50)"
