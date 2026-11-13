-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestAnalyzeDDLMatchesPostgres pg18#1
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
> exec
| CREATE VIEW v AS SELECT c FROM cc UNION ALL SELECT p FROM cp
< error
| 42P22 could not determine which collation to use for view column "c"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_1 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_2; SET search_path TO analyze_2
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
> exec
| CREATE VIEW v AS SELECT c AS renamed FROM cc UNION ALL SELECT p FROM cp
< error
| 42P22 could not determine which collation to use for view column "renamed"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_2 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_3; SET search_path TO analyze_3
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
> exec
| CREATE VIEW v (renamed) AS SELECT c FROM cc UNION ALL SELECT p FROM cp
< error
| 42P22 could not determine which collation to use for view column "renamed"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_3 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_4; SET search_path TO analyze_4
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
> exec
| CREATE VIEW v AS SELECT c FROM cc UNION ALL SELECT d FROM cc
< ok 0 -
> exec
| DROP SCHEMA analyze_4 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_5; SET search_path TO analyze_5
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
> exec
| CREATE VIEW v AS SELECT c FROM (SELECT c FROM cc UNION ALL SELECT p FROM cp) s
< error
| 42P22 could not determine which collation to use for view column "c"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_5 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_6; SET search_path TO analyze_6
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
> exec
| CREATE VIEW v AS WITH s AS (SELECT c FROM cc UNION ALL SELECT p FROM cp) SELECT c FROM s
< error
| 42P22 could not determine which collation to use for view column "c"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_6 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_7; SET search_path TO analyze_7
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
> exec
| CREATE VIEW v AS WITH s AS (SELECT c FROM cc UNION ALL SELECT p FROM cp) SELECT c FROM (SELECT c FROM s) x
< error
| 42P22 could not determine which collation to use for view column "c"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_7 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_8; SET search_path TO analyze_8
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
> exec
| CREATE VIEW v AS WITH s AS (SELECT c FROM cc UNION ALL SELECT p FROM cp) SELECT c FROM (WITH s AS (SELECT d AS c FROM cc) SELECT c FROM s) x
< ok 0 -
> exec
| DROP SCHEMA analyze_8 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_9; SET search_path TO analyze_9
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
> exec
| CREATE VIEW v AS SELECT s.c FROM (SELECT c FROM cc UNION ALL SELECT p FROM cp) s JOIN t ON true
< error
| 42P22 could not determine which collation to use for view column "c"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_9 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_10; SET search_path TO analyze_10
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
> exec
| CREATE VIEW v AS SELECT upper(c) AS u FROM (SELECT c FROM cc UNION ALL SELECT p FROM cp) s
< ok 0 -
> exec
| DROP SCHEMA analyze_10 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_11; SET search_path TO analyze_11
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
> exec
| CREATE VIEW v AS SELECT c || '' AS w FROM (SELECT c FROM cc UNION ALL SELECT p FROM cp) s
< ok 0 -
> exec
| DROP SCHEMA analyze_11 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_12; SET search_path TO analyze_12
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
> exec
| CREATE VIEW v AS SELECT (SELECT c FROM cc UNION ALL SELECT p FROM cp LIMIT 1) AS x
< error
| 42P22 could not determine which collation to use for view column "x"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_12 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_13; SET search_path TO analyze_13
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
> exec
| CREATE VIEW v AS SELECT upper((SELECT c FROM cc UNION ALL SELECT p FROM cp LIMIT 1)) AS x
< error
| 42P22 could not determine which collation to use for view column "x"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_13 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_14; SET search_path TO analyze_14
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
> exec
| CREATE VIEW v AS SELECT (SELECT c FROM cc UNION ALL SELECT p FROM cp LIMIT 1) || 'a' AS x
< error
| 42P22 could not determine which collation to use for view column "x"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_14 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_15; SET search_path TO analyze_15
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
> exec
| CREATE VIEW v AS SELECT (SELECT c FROM cc UNION ALL SELECT p FROM cp LIMIT 1) COLLATE "C" AS x
< ok 0 -
> exec
| DROP SCHEMA analyze_15 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_16; SET search_path TO analyze_16
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
> exec
| CREATE VIEW v AS SELECT (SELECT c FROM cc UNION ALL SELECT p FROM cp LIMIT 1) = 'a' AS x
< ok 0 -
> exec
| DROP SCHEMA analyze_16 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_17; SET search_path TO analyze_17
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
> exec
| CREATE VIEW v AS SELECT (SELECT c FROM cc UNION SELECT d FROM cc LIMIT 1) AS x
< ok 0 -
> exec
| DROP SCHEMA analyze_17 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_18; SET search_path TO analyze_18
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
> exec
| CREATE VIEW v AS SELECT (SELECT s.c) AS x FROM (SELECT c FROM cc UNION ALL SELECT p FROM cp) s
< error
| 42P22 could not determine which collation to use for view column "x"
| hint: Use the COLLATE clause to set the collation explicitly.
> exec
| DROP SCHEMA analyze_18 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_19; SET search_path TO analyze_19
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
> exec
| CREATE VIEW v AS SELECT (SELECT upper(s.c)) AS x FROM (SELECT c FROM cc UNION ALL SELECT p FROM cp) s
< ok 0 -
> exec
| DROP SCHEMA analyze_19 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_20; SET search_path TO analyze_20
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
> exec
| CREATE OR REPLACE FUNCTION poly_def(a int, b anyelement DEFAULT NULL::text) RETURNS anyelement LANGUAGE sql AS 'select $2'
< error
| 42P13 cannot change data type of existing parameter default value
| hint: Use DROP FUNCTION poly_def(integer,anyelement) first.
> exec
| DROP SCHEMA analyze_20 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_21; SET search_path TO analyze_21
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
> exec
| CREATE OR REPLACE FUNCTION f2(a int, b int) RETURNS int LANGUAGE sql AS 'select $1'
< error
| 42P13 cannot remove parameter defaults from existing function
| hint: Use DROP FUNCTION f2(integer,integer) first.
> exec
| DROP SCHEMA analyze_21 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_22; SET search_path TO analyze_22
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
> exec
| CREATE OR REPLACE FUNCTION f2(a int, b int DEFAULT 2) RETURNS int LANGUAGE sql AS 'select $1'
< ok 0 -
> exec
| DROP SCHEMA analyze_22 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_23; SET search_path TO analyze_23
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
> exec
| CREATE OR REPLACE FUNCTION f2(a int DEFAULT 0, b int DEFAULT 2) RETURNS int LANGUAGE sql AS 'select $1'; CREATE VIEW v AS SELECT f2() AS r
< ok 0 -
> exec
| DROP SCHEMA analyze_23 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_24; SET search_path TO analyze_24
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
> exec
| CREATE OR REPLACE FUNCTION f3(a int, b int DEFAULT 5) RETURNS int LANGUAGE sql AS 'select $1'; CREATE VIEW v AS SELECT f3(1) AS r
< ok 0 -
> exec
| DROP SCHEMA analyze_24 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_25; SET search_path TO analyze_25
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
> exec
| CREATE AGGREGATE my_bad(int) (sfunc = two_sfunc, stype = int, hypothetical)
< error
| 42P13 only ordered-set aggregates can be hypothetical
> exec
| DROP SCHEMA analyze_25 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_26; SET search_path TO analyze_26
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
> exec
| CREATE AGGREGATE my_hyp(int ORDER BY int) (sfunc = os_sfunc, stype = int, hypothetical)
< ok 0 -
> exec
| DROP SCHEMA analyze_26 CASCADE
< ok 0 -
