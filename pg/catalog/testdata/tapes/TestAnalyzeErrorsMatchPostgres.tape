-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestAnalyzeErrorsMatchPostgres pg18#1
> exec
| CREATE SCHEMA analyze_1; SET search_path TO analyze_1
< ok 0 -
> exec
|
| 	CREATE TABLE t (id int, masked text, x text);
| 	CREATE TABLE u (id int, y text);
| 	CREATE TABLE w (id int, x text, z int);
| 	CREATE TABLE p (a int, b text, c int);
| 	CREATE TABLE q (c int, b text, d int);
| 	CREATE VIEW cc AS SELECT x COLLATE "C" AS c, y AS d FROM t, u;
| 	CREATE VIEW cp AS SELECT y COLLATE "POSIX" AS p FROM u;
| 	CREATE FUNCTION os_sfunc(int, int) RETURNS int LANGUAGE sql AS 'select 1';
| 	CREATE AGGREGATE my_os(float8 ORDER BY int) (sfunc = os_sfunc, stype = int);
| 	CREATE FUNCTION two_sfunc(int, int, int) RETURNS int LANGUAGE sql AS 'select 1';
| 	CREATE AGGREGATE two_arg(int, int) (sfunc = two_sfunc, stype = int);
| 	CREATE FUNCTION poly_def(a int, b anyelement DEFAULT NULL::int) RETURNS anyelement LANGUAGE sql AS 'select $2';
| 	CREATE FUNCTION f2(a int, b int DEFAULT 1) RETURNS int LANGUAGE sql AS 'select $1';
| 	CREATE FUNCTION f3(a int, b int) RETURNS int LANGUAGE sql AS 'select $1';
| 	CREATE TABLE tt (id text);
| 	CREATE TABLE b8 (id bigint, q text);
| 	CREATE TABLE js1 (j json);
| 	CREATE TABLE js2 (j json);
| 	CREATE TABLE vc (v varchar);
| 	CREATE TABLE bx (b box);
| 	CREATE TABLE v1 (k varchar(10), a int);
| 	CREATE TABLE v2 (k varchar(10), b int);
| 	CREATE TABLE v3 (k varchar(20), c int);
|
< ok 0 -
> query
| EXPLAIN SELECT lower(x) FILTER (WHERE true) FROM t
< error
| 42809 FILTER specified, but lower is not an aggregate function
| position: 16
> query
| EXPLAIN SELECT lower(x ORDER BY x) FROM t
< error
| 42809 ORDER BY specified, but lower is not an aggregate function
| position: 16
> query
| EXPLAIN SELECT lower(DISTINCT x) FROM t
< error
| 42809 DISTINCT specified, but lower is not an aggregate function
| position: 16
> query
| EXPLAIN SELECT lower(x) OVER () FROM t
< error
| 42809 OVER specified, but lower is not a window function nor an aggregate function
| position: 16
> query
| EXPLAIN SELECT lower(x) WITHIN GROUP (ORDER BY x) FROM t
< error
| 42883 function lower(text, text) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT ntile() WITHIN GROUP (ORDER BY 3) OVER () FROM t
< error
| 42809 window function ntile cannot have WITHIN GROUP
| position: 16
> query
| EXPLAIN SELECT row_number() FROM t
< error
| 42809 window function row_number requires an OVER clause
| position: 16
> query
| EXPLAIN SELECT row_number() FILTER (WHERE true) OVER () FROM t
< error
| 0A000 FILTER is not implemented for non-aggregate window functions
| position: 16
> query
| EXPLAIN SELECT sum(id ORDER BY id) OVER () FROM t
< error
| 0A000 aggregate ORDER BY is not implemented for window functions
| position: 16
> query
| EXPLAIN SELECT sum(id) FILTER (WHERE id > 0) OVER () FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| WindowAgg  (cost=31.21..31.25 rows=850 width=8)
< row
: string
|   Window: w1 AS ()
< row
: string
|   ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
> query
| EXPLAIN SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY id) OVER () FROM t
< error
| 0A000 OVER is not supported for ordered-set aggregate percentile_cont
| position: 16
> query
| EXPLAIN SELECT percentile_cont(0.5, 1.0) FROM t
< error
| 42809 WITHIN GROUP is required for ordered-set aggregate percentile_cont
| position: 16
> query
| EXPLAIN SELECT string_agg(x) WITHIN GROUP (ORDER BY ',') FROM t
< error
| 42809 string_agg is not an ordered-set aggregate, so it cannot have WITHIN GROUP
| position: 16
> query
| EXPLAIN SELECT count(*) WITHIN GROUP (ORDER BY id) FROM t
< error
| 42809 count is not an ordered-set aggregate, so it cannot have WITHIN GROUP
| position: 16
> query
| EXPLAIN SELECT jsonb_path_exists('{}'::jsonb, '$.a', '{}', false, 1)
< error
| 42883 function jsonb_path_exists(jsonb, unknown, unknown, boolean, integer) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT * FROM t JOIN u USING (id, id)
< error
| 42701 column name "id" appears more than once in USING clause
> query
| EXPLAIN SELECT * FROM t JOIN u USING (nope)
< error
| 42703 column "nope" specified in USING clause does not exist in left table
> query
| EXPLAIN SELECT * FROM t JOIN u USING (masked)
< error
| 42703 column "masked" specified in USING clause does not exist in right table
> query
| EXPLAIN SELECT * FROM (t JOIN u ON true) NATURAL JOIN u AS u2
< error
| 42702 common column name "id" appears more than once in left table
> query
| EXPLAIN SELECT * FROM u AS u2 NATURAL JOIN (t JOIN u ON true)
< error
| 42702 common column name "id" appears more than once in right table
> query
| EXPLAIN SELECT * FROM (t JOIN u ON true) JOIN u AS u2 USING (y)
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Nested Loop  (cost=176.34..86004.29 rows=6854825 width=108)
< row
: string
|   ->  Merge Join  (cost=176.34..303.67 rows=8064 width=40)
< row
: string
|         Merge Cond: (u.y = u2.y)
< row
: string
|         ->  Sort  (cost=88.17..91.35 rows=1270 width=36)
< row
: string
|               Sort Key: u.y
< row
: string
|               ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=36)
< row
: string
|         ->  Sort  (cost=88.17..91.35 rows=1270 width=36)
< row
: string
|               Sort Key: u2.y
< row
: string
|               ->  Seq Scan on u u2  (cost=0.00..22.70 rows=1270 width=36)
< row
: string
|   ->  Materialize  (cost=0.00..22.75 rows=850 width=68)
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=68)
> query
| EXPLAIN SELECT id FROM t JOIN u ON t.id = u.id
< error
| 42702 column reference "id" is ambiguous
| position: 16
> query
| EXPLAIN SELECT id FROM t NATURAL JOIN u
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Merge Join  (cost=148.03..233.25 rows=5398 width=4)
< row
: string
|   Merge Cond: (t.id = u.id)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=4)
< row
: string
|         Sort Key: t.id
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
< row
: string
|   ->  Sort  (cost=88.17..91.35 rows=1270 width=4)
< row
: string
|         Sort Key: u.id
< row
: string
|         ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=4)
> query
| EXPLAIN SELECT masked FROM (t JOIN u USING (id)) JOIN w USING (id)
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Merge Join  (cost=231.40..804.28 rows=32388 width=32)
< row
: string
|   Merge Cond: (u.id = t.id)
< row
: string
|   ->  Sort  (cost=88.17..91.35 rows=1270 width=4)
< row
: string
|         Sort Key: u.id
< row
: string
|         ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=4)
< row
: string
|   ->  Materialize  (cost=143.23..236.73 rows=5100 width=40)
< row
: string
|         ->  Merge Join  (cost=143.23..223.98 rows=5100 width=40)
< row
: string
|               Merge Cond: (t.id = w.id)
< row
: string
|               ->  Sort  (cost=59.86..61.98 rows=850 width=36)
< row
: string
|                     Sort Key: t.id
< row
: string
|                     ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=36)
< row
: string
|               ->  Sort  (cost=83.37..86.37 rows=1200 width=4)
< row
: string
|                     Sort Key: w.id
< row
: string
|                     ->  Seq Scan on w  (cost=0.00..22.00 rows=1200 width=4)
> query
| EXPLAIN SELECT 1 FROM t JOIN u ON t.id = u.id JOIN w ON y = w.x
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Merge Join  (cost=651.25..1143.07 rows=32388 width=4)
< row
: string
|   Merge Cond: (w.x = u.y)
< row
: string
|   ->  Sort  (cost=83.37..86.37 rows=1200 width=32)
< row
: string
|         Sort Key: w.x
< row
: string
|         ->  Seq Scan on w  (cost=0.00..22.00 rows=1200 width=32)
< row
: string
|   ->  Sort  (cost=567.88..581.37 rows=5398 width=32)
< row
: string
|         Sort Key: u.y
< row
: string
|         ->  Merge Join  (cost=148.03..233.25 rows=5398 width=32)
< row
: string
|               Merge Cond: (t.id = u.id)
< row
: string
|               ->  Sort  (cost=59.86..61.98 rows=850 width=4)
< row
: string
|                     Sort Key: t.id
< row
: string
|                     ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
< row
: string
|               ->  Sort  (cost=88.17..91.35 rows=1270 width=36)
< row
: string
|                     Sort Key: u.id
< row
: string
|                     ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=36)
> query
| EXPLAIN SELECT 1 FROM (t JOIN u USING (id)) JOIN w ON masked = w.x
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Merge Join  (cost=626.22..1118.34 rows=32388 width=4)
< row
: string
|   Merge Cond: (u.id = t.id)
< row
: string
|   ->  Sort  (cost=88.17..91.35 rows=1270 width=4)
< row
: string
|         Sort Key: u.id
< row
: string
|         ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=4)
< row
: string
|   ->  Sort  (cost=538.05..550.80 rows=5100 width=4)
< row
: string
|         Sort Key: t.id
< row
: string
|         ->  Merge Join  (cost=143.23..223.98 rows=5100 width=4)
< row
: string
|               Merge Cond: (t.masked = w.x)
< row
: string
|               ->  Sort  (cost=59.86..61.98 rows=850 width=36)
< row
: string
|                     Sort Key: t.masked
< row
: string
|                     ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=36)
< row
: string
|               ->  Sort  (cost=83.37..86.37 rows=1200 width=32)
< row
: string
|                     Sort Key: w.x
< row
: string
|                     ->  Seq Scan on w  (cost=0.00..22.00 rows=1200 width=32)
> query
| EXPLAIN SELECT 1 FROM (t JOIN u USING (id)) JOIN w ON id = w.z
< error
| 42702 column reference "id" is ambiguous
| position: 55
> query
| EXPLAIN SELECT 1 FROM t NATURAL JOIN u JOIN w ON y = w.x AND z > 0
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Merge Join  (cost=348.28..514.45 rows=10796 width=4)
< row
: string
|   Merge Cond: (t.id = u.id)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=4)
< row
: string
|         Sort Key: t.id
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
< row
: string
|   ->  Sort  (cost=288.42..294.77 rows=2540 width=4)
< row
: string
|         Sort Key: u.id
< row
: string
|         ->  Hash Join  (cost=30.00..144.78 rows=2540 width=4)
< row
: string
|               Hash Cond: (u.y = w.x)
< row
: string
|               ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=36)
< row
: string
|               ->  Hash  (cost=25.00..25.00 rows=400 width=32)
< row
: string
|                     ->  Seq Scan on w  (cost=0.00..25.00 rows=400 width=32)
< row
: string
|                           Filter: (z > 0)
> query
| EXPLAIN SELECT 1 FROM t, u JOIN w ON y = w.x
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Nested Loop  (cost=171.55..81274.97 rows=6477000 width=4)
< row
: string
|   ->  Merge Join  (cost=171.55..291.85 rows=7620 width=0)
< row
: string
|         Merge Cond: (w.x = u.y)
< row
: string
|         ->  Sort  (cost=83.37..86.37 rows=1200 width=32)
< row
: string
|               Sort Key: w.x
< row
: string
|               ->  Seq Scan on w  (cost=0.00..22.00 rows=1200 width=32)
< row
: string
|         ->  Sort  (cost=88.17..91.35 rows=1270 width=32)
< row
: string
|               Sort Key: u.y
< row
: string
|               ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=32)
< row
: string
|   ->  Materialize  (cost=0.00..22.75 rows=850 width=0)
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=0)
> query
| EXPLAIN SELECT 1 FROM t JOIN u ON y = (SELECT max(x) FROM w WHERE w.id = u.id)
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Nested Loop  (cost=0.00..31889.89 rows=5100 width=4)
< row
: string
|   ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=0)
< row
: string
|   ->  Materialize  (cost=0.00..31807.66 rows=6 width=0)
< row
: string
|         ->  Seq Scan on u  (cost=0.00..31807.63 rows=6 width=0)
< row
: string
|               Filter: (y = (SubPlan 1))
< row
: string
|               SubPlan 1
< row
: string
|                 ->  Aggregate  (cost=25.02..25.03 rows=1 width=32)
< row
: string
|                       ->  Seq Scan on w  (cost=0.00..25.00 rows=6 width=32)
< row
: string
|                             Filter: (id = u.id)
> query
| EXPLAIN SELECT 1 FROM w f, t JOIN u ON f.z = t.id
< error
| 42P01 invalid reference to FROM-clause entry for table "f"
| position: 40
| detail: There is an entry for table "f", but it cannot be referenced from this part of the query.
> query
| EXPLAIN SELECT 1 FROM w, t JOIN u ON z = t.id
< error
| 42703 column "z" does not exist
| position: 38
| detail: There is a column named "z" in table "w", but it cannot be referenced from this part of the query.
> query
| EXPLAIN SELECT 1 FROM w f, t JOIN u ON t.id = u.id WHERE f.z = 1
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Nested Loop  (cost=148.03..663.12 rows=32385 width=4)
< row
: string
|   ->  Merge Join  (cost=148.03..233.25 rows=5398 width=0)
< row
: string
|         Merge Cond: (t.id = u.id)
< row
: string
|         ->  Sort  (cost=59.86..61.98 rows=850 width=4)
< row
: string
|               Sort Key: t.id
< row
: string
|               ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
< row
: string
|         ->  Sort  (cost=88.17..91.35 rows=1270 width=4)
< row
: string
|               Sort Key: u.id
< row
: string
|               ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=4)
< row
: string
|   ->  Materialize  (cost=0.00..25.03 rows=6 width=0)
< row
: string
|         ->  Seq Scan on w f  (cost=0.00..25.00 rows=6 width=0)
< row
: string
|               Filter: (z = 1)
> query
| EXPLAIN SELECT 1 FROM t JOIN u ON (SELECT w.z FROM w WHERE w.id = t.id) = u.id
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Merge Join  (cost=148.03..156433.25 rows=5398 width=4)
< row
: string
|   Merge Cond: (((SubPlan 1)) = u.id)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=4)
< row
: string
|         Sort Key: ((SubPlan 1))
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
< row
: string
|               SubPlan 1
< row
: string
|                 ->  Seq Scan on w  (cost=0.00..25.00 rows=6 width=4)
< row
: string
|                       Filter: (id = t.id)
< row
: string
|   ->  Sort  (cost=88.17..91.35 rows=1270 width=4)
< row
: string
|         Sort Key: u.id
< row
: string
|         ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=4)
< row
: string
| JIT:
< row
: string
|   Functions: 13
< row
: string
|   Options: Inlining false, Optimization false, Expressions true, Deforming true
> query
| EXPLAIN SELECT (SELECT f.masked FROM u f) FROM t f
< error
| 42703 column f.masked does not exist
| position: 24
| detail: There is a column named "masked" in table "f", but it cannot be referenced from this part of the query.
> query
| EXPLAIN SELECT (SELECT f.y FROM u f) FROM t f
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on t f  (cost=22.70..41.20 rows=850 width=32)
< row
: string
|   InitPlan 1
< row
: string
|     ->  Seq Scan on u f_1  (cost=0.00..22.70 rows=1270 width=32)
> query
| EXPLAIN SELECT (SELECT f.masked FROM u) FROM t f
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Seq Scan on t f  (cost=0.00..19313.50 rows=850 width=32)
< row
: string
|   SubPlan 1
< row
: string
|     ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=32)
> query
| EXPLAIN SELECT s.a FROM (SELECT 1 AS a, 2 AS a) s
< error
| 42702 column reference "a" is ambiguous
| position: 16
> query
| EXPLAIN SELECT s.b FROM (SELECT 1 AS a, 2 AS a, 3 AS b) s
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Result  (cost=0.00..0.01 rows=1 width=4)
> query
| EXPLAIN SELECT f.a FROM json_to_record('{}') AS f(a int, a int)
< error
| 42701 column name "a" specified more than once
> query
| EXPLAIN SELECT pg_logical_slot_get_changes('slot', NULL, NULL)
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| ProjectSet  (cost=0.00..7.51 rows=1000 width=32)
< row
: string
|   ->  Result  (cost=0.00..0.01 rows=1 width=0)
> query
| EXPLAIN SELECT pg_logical_slot_get_changes('slot', NULL)
< error
| 42883 function pg_logical_slot_get_changes(unknown, unknown) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT poly_def(1)
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Result  (cost=0.00..0.01 rows=1 width=4)
> query
| EXPLAIN SELECT poly_def(1, 'x'::text)
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Result  (cost=0.00..0.01 rows=1 width=32)
> query
| EXPLAIN SELECT jsonb_path_exists('{}'::jsonb, '$.a', silent => true)
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Result  (cost=0.00..0.01 rows=1 width=1)
> query
| EXPLAIN SELECT jsonb_path_exists(target => '{}'::jsonb, path => '$.a')
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Result  (cost=0.00..0.01 rows=1 width=1)
> query
| EXPLAIN SELECT jsonb_path_exists('{}'::jsonb, silent => true)
< error
| 42883 function jsonb_path_exists(jsonb, silent => boolean) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT jsonb_path_exists('{}'::jsonb, '$.a', target => '{}')
< error
| 42883 function jsonb_path_exists(jsonb, unknown, target => unknown) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT make_interval(days => 1, days => 2)
< error
| 42601 argument name "days" used more than once
| position: 41
> query
| EXPLAIN SELECT make_interval(days => 1, 2)
< error
| 42601 positional argument cannot follow named argument
| position: 41
> query
| EXPLAIN SELECT make_interval(nope => 1)
< error
| 42883 function make_interval(nope => integer) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT format(formatstr => 'x')
< error
| 42883 function format(formatstr => unknown) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT percentile_cont(fraction => 0.5) WITHIN GROUP (ORDER BY id) FROM t
< error
| 42601 positional argument cannot follow named argument
| position: 72
> query
| EXPLAIN SELECT my_os(0.5, id) FROM t
< error
| 42809 WITHIN GROUP is required for ordered-set aggregate my_os
| position: 16
> query
| EXPLAIN SELECT two_arg(id) WITHIN GROUP (ORDER BY id) FROM t
< error
| 42809 two_arg is not an ordered-set aggregate, so it cannot have WITHIN GROUP
| position: 16
> query
| EXPLAIN SELECT two_arg(id, id) FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=231.00..231.01 rows=1 width=4)
< row
: string
|   ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
> query
| EXPLAIN SELECT c FROM cc UNION SELECT p FROM cp
< error
| 42P21 collation mismatch between implicit collations "C" and "POSIX"
| position: 39
| hint: You can choose the collation by applying the COLLATE clause to one or both expressions.
> query
| EXPLAIN SELECT c FROM cc EXCEPT ALL SELECT p FROM cp
< error
| 42P21 collation mismatch between implicit collations "C" and "POSIX"
| position: 44
| hint: You can choose the collation by applying the COLLATE clause to one or both expressions.
> query
| EXPLAIN SELECT c FROM cc INTERSECT SELECT d FROM cc
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| HashSetOp Intersect  (cost=43266.65..43267.15 rows=200 width=32)
< row
: string
|   ->  Nested Loop  (cost=0.00..13537.08 rows=1079500 width=32)
< row
: string
|         ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=0)
< row
: string
|         ->  Materialize  (cost=0.00..22.75 rows=850 width=32)
< row
: string
|               ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=32)
< row
: string
|   ->  Subquery Scan on "*SELECT* 2"  (cost=0.00..24332.08 rows=1079500 width=32)
< row
: string
|         ->  Nested Loop  (cost=0.00..13537.08 rows=1079500 width=32)
< row
: string
|               ->  Seq Scan on u u_1  (cost=0.00..22.70 rows=1270 width=32)
< row
: string
|               ->  Materialize  (cost=0.00..22.75 rows=850 width=0)
< row
: string
|                     ->  Seq Scan on t t_1  (cost=0.00..18.50 rows=850 width=0)
> query
| EXPLAIN SELECT c FROM cc UNION ALL SELECT p FROM cp
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Append  (cost=0.00..18963.62 rows=1080770 width=32)
< row
: string
|   ->  Nested Loop  (cost=0.00..13537.08 rows=1079500 width=32)
< row
: string
|         ->  Seq Scan on u u_1  (cost=0.00..22.70 rows=1270 width=0)
< row
: string
|         ->  Materialize  (cost=0.00..22.75 rows=850 width=32)
< row
: string
|               ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=32)
< row
: string
|   ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=32)
> query
| EXPLAIN SELECT c COLLATE "POSIX" FROM cc UNION ALL SELECT d COLLATE "C" FROM cc
< error
| 42P21 collation mismatch between explicit collations "POSIX" and "C"
| position: 61
> query
| EXPLAIN SELECT c COLLATE "C" FROM cc UNION SELECT p FROM cp
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| HashAggregate  (cost=102398.26..127982.11 rows=1080770 width=32)
< row
: string
|   Group Key: ((t.x)::text)
< row
: string
|   Planned Partitions: 16
< row
: string
|   ->  Append  (cost=0.00..18976.33 rows=1080770 width=32)
< row
: string
|         ->  Nested Loop  (cost=0.00..13537.08 rows=1079500 width=32)
< row
: string
|               ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=0)
< row
: string
|               ->  Materialize  (cost=0.00..22.75 rows=850 width=32)
< row
: string
|                     ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=32)
< row
: string
|         ->  Subquery Scan on "*SELECT* 2"  (cost=0.00..35.40 rows=1270 width=32)
< row
: string
|               ->  Seq Scan on u u_1  (cost=0.00..22.70 rows=1270 width=32)
< row
: string
| JIT:
< row
: string
|   Functions: 9
< row
: string
|   Options: Inlining false, Optimization false, Expressions true, Deforming true
> query
| EXPLAIN SELECT (d COLLATE "C") || '' FROM cc UNION SELECT p FROM cp
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| HashAggregate  (cost=105097.01..130680.86 rows=1080770 width=32)
< row
: string
|   Group Key: (((u.y)::text || ''::text))
< row
: string
|   Planned Partitions: 16
< row
: string
|   ->  Append  (cost=0.00..21675.08 rows=1080770 width=32)
< row
: string
|         ->  Nested Loop  (cost=0.00..16235.83 rows=1079500 width=32)
< row
: string
|               ->  Seq Scan on u  (cost=0.00..22.70 rows=1270 width=32)
< row
: string
|               ->  Materialize  (cost=0.00..22.75 rows=850 width=0)
< row
: string
|                     ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=0)
< row
: string
|         ->  Subquery Scan on "*SELECT* 2"  (cost=0.00..35.40 rows=1270 width=32)
< row
: string
|               ->  Seq Scan on u u_1  (cost=0.00..22.70 rows=1270 width=32)
< row
: string
| JIT:
< row
: string
|   Functions: 8
< row
: string
|   Options: Inlining false, Optimization false, Expressions true, Deforming true
> query
| EXPLAIN SELECT (d COLLATE "C") || '' FROM cc UNION ALL SELECT (p COLLATE "POSIX") || '' FROM cp
< error
| 42P21 collation mismatch between explicit collations "C" and "POSIX"
| position: 66
> query
| EXPLAIN SELECT upper(c) FROM cc UNION SELECT p FROM cp
< error
| 42P21 collation mismatch between implicit collations "C" and "POSIX"
| position: 46
| hint: You can choose the collation by applying the COLLATE clause to one or both expressions.
> query
| EXPLAIN SELECT * FROM t NATURAL JOIN tt
< error
| 42883 operator does not exist: integer = text
| hint: No operator matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT * FROM t JOIN tt USING (id)
< error
| 42883 operator does not exist: integer = text
| hint: No operator matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT * FROM js1 NATURAL JOIN js2
< error
| 42883 operator does not exist: json = json
| hint: No operator matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT * FROM t NATURAL JOIN b8
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Merge Join  (cost=143.23..223.98 rows=5100 width=104)
< row
: string
|   Merge Cond: (t.id = b8.id)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=68)
< row
: string
|         Sort Key: t.id
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=68)
< row
: string
|   ->  Sort  (cost=83.37..86.37 rows=1200 width=40)
< row
: string
|         Sort Key: b8.id
< row
: string
|         ->  Seq Scan on b8  (cost=0.00..22.00 rows=1200 width=40)
> query
| EXPLAIN SELECT my_os() WITHIN GROUP (ORDER BY 0.5::float8, id) FROM t
< error
| 42883 function my_os(double precision, integer) does not exist
| position: 16
| hint: There is an ordered-set aggregate my_os, but it requires 1 direct argument, not 0.
> query
| EXPLAIN SELECT percentile_cont() WITHIN GROUP (ORDER BY 0.5::float8, 1.0::float8)
< error
| 42883 function percentile_cont(double precision, double precision) does not exist
| position: 16
| hint: There is an ordered-set aggregate percentile_cont, but it requires 1 direct argument, not 0.
> query
| EXPLAIN SELECT mode(1) WITHIN GROUP (ORDER BY id) FROM t
< error
| 42883 function mode(integer, integer) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT rank(1, 2) WITHIN GROUP (ORDER BY id) FROM t
< error
| 42883 function rank(integer, integer, integer) does not exist
| position: 16
| hint: To use the hypothetical-set aggregate rank, the number of hypothetical direct arguments (here 2) must match the number of ordering columns (here 1).
> query
| EXPLAIN SELECT rank(1) WITHIN GROUP (ORDER BY id, id) FROM t
< error
| 42883 function rank(integer, integer, integer) does not exist
| position: 16
| hint: To use the hypothetical-set aggregate rank, the number of hypothetical direct arguments (here 1) must match the number of ordering columns (here 2).
> query
| EXPLAIN SELECT rank(1, 2) WITHIN GROUP (ORDER BY id, id) FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=20.63..20.64 rows=1 width=8)
< row
: string
|   ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
> query
| EXPLAIN SELECT array_agg(DISTINCT x ORDER BY masked) FROM t
< error
| 42P10 in an aggregate with DISTINCT, ORDER BY expressions must appear in argument list
| position: 46
> query
| EXPLAIN SELECT string_agg(DISTINCT x, ',' ORDER BY ',') FROM t
< error
| 42P10 in an aggregate with DISTINCT, ORDER BY expressions must appear in argument list
| position: 52
> query
| EXPLAIN SELECT array_agg(DISTINCT x ORDER BY x DESC) FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=64.11..64.12 rows=1 width=32)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=32)
< row
: string
|         Sort Key: x DESC
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=32)
> query
| EXPLAIN SELECT string_agg(DISTINCT v, ',' ORDER BY v) FROM vc
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=101.19..101.20 rows=1 width=32)
< row
: string
|   ->  Sort  (cost=94.38..97.78 rows=1360 width=32)
< row
: string
|         Sort Key: v
< row
: string
|         ->  Seq Scan on vc  (cost=0.00..23.60 rows=1360 width=32)
> query
| EXPLAIN SELECT array_agg(DISTINCT id::text ORDER BY id::text) FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=68.36..68.37 rows=1 width=32)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=4)
< row
: string
|         Sort Key: ((id)::text)
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
> query
| EXPLAIN SELECT array_agg(id ORDER BY id USING OPERATOR(pg_catalog.>)) FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=64.11..64.12 rows=1 width=32)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=4)
< row
: string
|         Sort Key: id DESC
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
> query
| EXPLAIN SELECT sum(id) FILTER (WHERE id) FROM t
< error
| 42804 argument of FILTER must be type boolean, not type integer
| position: 38
> query
| EXPLAIN SELECT sum(id) FILTER (WHERE x) FROM t
< error
| 42804 argument of FILTER must be type boolean, not type text
| position: 38
> query
| EXPLAIN SELECT sum(id) FILTER (WHERE 'true') FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=20.62..20.64 rows=1 width=8)
< row
: string
|   ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=4)
> query
| EXPLAIN SELECT count(*) FILTER (WHERE id) OVER () FROM t
< error
| 42804 argument of FILTER must be type boolean, not type integer
| position: 39
> query
| EXPLAIN SELECT count(*) FILTER (WHERE id) FROM t
< error
| 42804 argument of FILTER must be type boolean, not type integer
| position: 39
> query
| EXPLAIN SELECT array_agg(1 ORDER BY j) FROM js1
< error
| 42883 could not identify an ordering operator for type json
| position: 37
| hint: Use an explicit ordering operator or modify the query.
> query
| EXPLAIN SELECT array_agg(1 ORDER BY j DESC) FROM js1
< error
| 42883 could not identify an ordering operator for type json
| position: 37
| hint: Use an explicit ordering operator or modify the query.
> query
| EXPLAIN SELECT array_agg(1 ORDER BY j USING <) FROM js1
< error
| 42883 operator does not exist: json < json
| position: 45
| hint: No operator matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT array_agg(1 ORDER BY ARRAY[j]) FROM js1
< error
| 42883 could not identify an ordering operator for type json[]
| position: 37
| hint: Use an explicit ordering operator or modify the query.
> query
| EXPLAIN SELECT array_agg(1 ORDER BY ARRAY[x]) FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=64.11..64.12 rows=1 width=32)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=32)
< row
: string
|         Sort Key: (ARRAY[x])
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=32)
> query
| EXPLAIN SELECT array_agg(1 ORDER BY x::varchar) FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=64.11..64.12 rows=1 width=32)
< row
: string
|   ->  Sort  (cost=59.86..61.98 rows=850 width=32)
< row
: string
|         Sort Key: x
< row
: string
|         ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=32)
> query
| EXPLAIN SELECT array_agg(1 ORDER BY 'a') FROM t
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Aggregate  (cost=20.63..20.64 rows=1 width=32)
< row
: string
|   ->  Seq Scan on t  (cost=0.00..18.50 rows=850 width=0)
> query
| EXPLAIN SELECT mode() WITHIN GROUP (ORDER BY j) FROM js1
< error
| 42883 could not identify an ordering operator for type json
| position: 46
| hint: Use an explicit ordering operator or modify the query.
> query
| EXPLAIN SELECT num_nulls()
< error
| 42883 function num_nulls() does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT concat()
< error
| 42883 function concat() does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT json_build_array()
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Result  (cost=0.00..0.01 rows=1 width=32)
> query
| EXPLAIN SELECT json_build_object()
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Result  (cost=0.00..0.01 rows=1 width=32)
> query
| EXPLAIN SELECT json_object_agg(key => x, value => id) FROM t
< error
| 0A000 aggregates cannot use named arguments
| position: 16
> query
| EXPLAIN SELECT string_agg(x, delimiter => ',') FROM t
< error
| 0A000 aggregates cannot use named arguments
| position: 16
> query
| EXPLAIN SELECT ntile(num_buckets => 2) OVER () FROM t
< error
| 42883 function ntile(num_buckets => integer) does not exist
| position: 16
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> query
| EXPLAIN SELECT array_agg(b ORDER BY b) FROM bx
< error
| 42883 could not identify an ordering operator for type box
| position: 37
| hint: Use an explicit ordering operator or modify the query.
> query
| EXPLAIN SELECT array_agg(1 ORDER BY ARRAY[b]) FROM bx
< error
| 42883 could not identify an ordering operator for type box[]
| position: 37
| hint: Use an explicit ordering operator or modify the query.
> exec
| DROP SCHEMA analyze_1 CASCADE
< ok 0 -
