-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestAnalyzeViewsMatchPostgres pg18#1
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
| CREATE VIEW v AS SELECT sum(id) FILTER (WHERE masked = 'a') AS s, count(*) FILTER (WHERE x IS NULL) AS c FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT sum(id) FILTER (WHERE masked = 'a'::text) AS s,
|     count(*) FILTER (WHERE x IS NULL) AS c
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| s
: int64 20
: int64 -1
: int64 0
< row
: string
| c
: int64 20
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT string_agg(x, ',' ORDER BY masked DESC NULLS LAST) AS s, array_agg(DISTINCT x ORDER BY x) AS a FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT string_agg(x, ','::text ORDER BY masked DESC NULLS LAST) AS s,
|     array_agg(DISTINCT x ORDER BY x) AS a
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| s
: int64 25
: int64 -1
: int64 100
< row
: string
| a
: int64 1009
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT array_agg(id ORDER BY x NULLS FIRST, masked DESC NULLS FIRST) AS a FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT array_agg(id ORDER BY x NULLS FIRST, masked DESC) AS a
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 1007
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT array_agg(x ORDER BY masked USING >, id USING <) AS a FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT array_agg(x ORDER BY masked DESC, id) AS a
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 1009
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT array_agg(x ORDER BY masked USING OPERATOR(pg_catalog.>)) AS a FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT array_agg(x ORDER BY masked DESC) AS a
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 1009
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT string_agg(x, ',' ORDER BY lower(masked), id + 1) AS s FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT string_agg(x, ','::text ORDER BY (lower(masked)), (id + 1)) AS s
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| s
: int64 25
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT string_agg(x, ',' ORDER BY masked) FILTER (WHERE id > 0) AS s FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT string_agg(x, ','::text ORDER BY masked) FILTER (WHERE id > 0) AS s
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| s
: int64 25
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY id) AS p, mode() WITHIN GROUP (ORDER BY x) AS m FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT percentile_cont(0.5::double precision) WITHIN GROUP (ORDER BY (id::double precision)) AS p,
|     mode() WITHIN GROUP (ORDER BY x) AS m
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| p
: int64 701
: int64 -1
: int64 0
< row
: string
| m
: int64 25
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT percentile_disc(0.5) WITHIN GROUP (ORDER BY x DESC) FILTER (WHERE id > 0) AS p FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT percentile_disc(0.5::double precision) WITHIN GROUP (ORDER BY x DESC) FILTER (WHERE id > 0) AS p
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| p
: int64 25
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT rank(3) WITHIN GROUP (ORDER BY id) AS r FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT rank(3) WITHIN GROUP (ORDER BY id) AS r
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| r
: int64 20
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM t NATURAL JOIN u;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.id,
|     t.masked,
|     t.x,
|     u.y
|    FROM t
|      JOIN u USING (id);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| x
: int64 25
: int64 -1
: int64 100
< row
: string
| y
: int64 25
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT * FROM p NATURAL JOIN q;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT p.b,
|     p.c,
|     p.a,
|     q.d
|    FROM p
|      JOIN q USING (b, c);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| b
: int64 25
: int64 -1
: int64 100
< row
: string
| c
: int64 23
: int64 -1
: int64 0
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| d
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM t NATURAL LEFT JOIN w;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.id,
|     t.x,
|     t.masked,
|     w.z
|    FROM t
|      LEFT JOIN w USING (id, x);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| x
: int64 25
: int64 -1
: int64 100
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| z
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM t NATURAL RIGHT JOIN w;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT w.id,
|     w.x,
|     t.masked,
|     w.z
|    FROM t
|      RIGHT JOIN w USING (id, x);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| x
: int64 25
: int64 -1
: int64 100
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| z
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM t NATURAL FULL JOIN w;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     x,
|     t.masked,
|     w.z
|    FROM t
|      FULL JOIN w USING (id, x);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| x
: int64 25
: int64 -1
: int64 100
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| z
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM u NATURAL JOIN (SELECT 1 AS k) s;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT u.id,
|     u.y,
|     s.k
|    FROM u
|      CROSS JOIN ( SELECT 1 AS k) s;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| y
: int64 25
: int64 -1
: int64 100
< row
: string
| k
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM u NATURAL LEFT JOIN (SELECT 1 AS k) s;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT u.id,
|     u.y,
|     s.k
|    FROM u
|      LEFT JOIN ( SELECT 1 AS k) s ON TRUE;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| y
: int64 25
: int64 -1
: int64 100
< row
: string
| k
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM t NATURAL JOIN u NATURAL JOIN w;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.id,
|     t.x,
|     t.masked,
|     u.y,
|     w.z
|    FROM t
|      JOIN u USING (id)
|      JOIN w USING (id, x);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| x
: int64 25
: int64 -1
: int64 100
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| y
: int64 25
: int64 -1
: int64 100
< row
: string
| z
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM (p NATURAL JOIN q) NATURAL JOIN u;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT p.b,
|     p.c,
|     p.a,
|     q.d,
|     u.id,
|     u.y
|    FROM p
|      JOIN q USING (b, c)
|      CROSS JOIN u;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| b
: int64 25
: int64 -1
: int64 100
< row
: string
| c
: int64 23
: int64 -1
: int64 0
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| d
: int64 23
: int64 -1
: int64 0
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| y
: int64 25
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT id, masked, y FROM t NATURAL JOIN u WHERE id > 1;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.id,
|     t.masked,
|     u.y
|    FROM t
|      JOIN u USING (id)
|   WHERE t.id > 1;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| y
: int64 25
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT masked, y, z FROM t NATURAL JOIN u JOIN w ON y = w.x;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.masked,
|     u.y,
|     w.z
|    FROM t
|      JOIN u USING (id)
|      JOIN w ON u.y = w.x;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| y
: int64 25
: int64 -1
: int64 100
< row
: string
| z
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM p JOIN (SELECT b, d FROM q) s USING (b);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT p.b,
|     p.a,
|     p.c,
|     s.d
|    FROM p
|      JOIN ( SELECT q.b,
|             q.d
|            FROM q) s USING (b);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| b
: int64 25
: int64 -1
: int64 100
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| c
: int64 23
: int64 -1
: int64 0
< row
: string
| d
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT * FROM p FULL JOIN q USING (b, c);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT b,
|     c,
|     p.a,
|     q.d
|    FROM p
|      FULL JOIN q USING (b, c);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| b
: int64 25
: int64 -1
: int64 100
< row
: string
| c
: int64 23
: int64 -1
: int64 0
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| d
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT c, a, d FROM p RIGHT JOIN q USING (c);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT q.c,
|     p.a,
|     q.d
|    FROM p
|      RIGHT JOIN q USING (c);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| c
: int64 23
: int64 -1
: int64 0
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| d
: int64 23
: int64 -1
: int64 0
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
| CREATE VIEW v AS SELECT y, masked FROM t JOIN u ON t.id = u.id;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT u.y,
|     t.masked
|    FROM t
|      JOIN u ON t.id = u.id;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| y
: int64 25
: int64 -1
: int64 100
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
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
| CREATE VIEW v AS SELECT y, id FROM t JOIN u USING (id);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT u.y,
|     t.id
|    FROM t
|      JOIN u USING (id);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| y
: int64 25
: int64 -1
: int64 100
< row
: string
| id
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_26 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_27; SET search_path TO analyze_27
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
| CREATE VIEW v AS SELECT id, x FROM t UNION SELECT id, y FROM u;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.id,
|     t.x
|    FROM t
| UNION
|  SELECT u.id,
|     u.y AS x
|    FROM u;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 23
: int64 -1
: int64 0
< row
: string
| x
: int64 25
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_27 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_28; SET search_path TO analyze_28
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
| CREATE VIEW v AS SELECT jsonb_path_exists('{}'::jsonb, '$."a"') AS e, jsonb_set('{}'::jsonb, '{a}', '1') AS s, random_normal() AS r;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT jsonb_path_exists('{}'::jsonb, '$."a"'::jsonpath) AS e,
|     jsonb_set('{}'::jsonb, '{a}'::text[], '1'::jsonb) AS s,
|     random_normal() AS r;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| e
: int64 16
: int64 -1
: int64 0
< row
: string
| s
: int64 3802
: int64 -1
: int64 0
< row
: string
| r
: int64 701
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_28 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_29; SET search_path TO analyze_29
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
| CREATE VIEW v AS SELECT make_interval(days => 1) AS a, jsonb_path_exists('{}'::jsonb, '$."a"', silent => true) AS b, make_interval(1, mins => id) AS c FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT make_interval(days => 1) AS a,
|     jsonb_path_exists('{}'::jsonb, '$."a"'::jsonpath, silent => true) AS b,
|     make_interval(1, mins => id) AS c
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 1186
: int64 -1
: int64 0
< row
: string
| b
: int64 16
: int64 -1
: int64 0
< row
: string
| c
: int64 1186
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_29 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_30; SET search_path TO analyze_30
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
| CREATE VIEW v AS SELECT jsonb_path_exists(path => '$."a"', target => '{}'::jsonb) AS e;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT jsonb_path_exists(path => '$."a"'::jsonpath, target => '{}'::jsonb) AS e;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| e
: int64 16
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_30 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_31; SET search_path TO analyze_31
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
| CREATE VIEW v AS SELECT my_os(0.5) WITHIN GROUP (ORDER BY id) AS m FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT my_os(0.5::double precision) WITHIN GROUP (ORDER BY id) AS m
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| m
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_31 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_32; SET search_path TO analyze_32
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
| CREATE VIEW v AS SELECT c FROM cc UNION SELECT d FROM cc;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT cc.c
|    FROM cc
| UNION
|  SELECT cc.d AS c
|    FROM cc;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| c
: int64 25
: int64 -1
: int64 950
> exec
| DROP SCHEMA analyze_32 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_33; SET search_path TO analyze_33
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
| CREATE VIEW v AS SELECT d FROM cc UNION SELECT c FROM cc;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT cc.d
|    FROM cc
| UNION
|  SELECT cc.c AS d
|    FROM cc;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| d
: int64 25
: int64 -1
: int64 950
> exec
| DROP SCHEMA analyze_33 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_34; SET search_path TO analyze_34
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
| CREATE VIEW v AS SELECT c COLLATE "POSIX" AS c FROM cc UNION SELECT p FROM cp;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT cc.c COLLATE "POSIX" AS c
|    FROM cc
| UNION
|  SELECT cp.p AS c
|    FROM cp;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| c
: int64 25
: int64 -1
: int64 951
> exec
| DROP SCHEMA analyze_34 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_35; SET search_path TO analyze_35
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
| CREATE VIEW v AS (SELECT c FROM cc UNION ALL SELECT p FROM cp) UNION ALL SELECT d FROM cc;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT cc.c
|    FROM cc
| UNION ALL
|  SELECT cp.p AS c
|    FROM cp
| UNION ALL
|  SELECT cc.d AS c
|    FROM cc;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| c
: int64 25
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_35 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_36; SET search_path TO analyze_36
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
| CREATE VIEW v AS SELECT (d COLLATE "C") || '' AS d FROM cc UNION SELECT p FROM cp;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (cc.d COLLATE "C") || ''::text AS d
|    FROM cc
| UNION
|  SELECT cp.p AS d
|    FROM cp;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| d
: int64 25
: int64 -1
: int64 950
> exec
| DROP SCHEMA analyze_36 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_37; SET search_path TO analyze_37
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
| CREATE VIEW v AS SELECT d || '' AS d FROM cc UNION SELECT p FROM cp;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT cc.d || ''::text AS d
|    FROM cc
| UNION
|  SELECT cp.p AS d
|    FROM cp;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| d
: int64 25
: int64 -1
: int64 951
> exec
| DROP SCHEMA analyze_37 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_38; SET search_path TO analyze_38
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
| CREATE VIEW v AS SELECT * FROM t NATURAL JOIN b8;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT b8.id,
|     t.masked,
|     t.x,
|     b8.q
|    FROM t
|      JOIN b8 USING (id);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 20
: int64 -1
: int64 0
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| x
: int64 25
: int64 -1
: int64 100
< row
: string
| q
: int64 25
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_38 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_39; SET search_path TO analyze_39
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
| CREATE VIEW v AS SELECT * FROM t LEFT JOIN b8 USING (id);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     t.masked,
|     t.x,
|     b8.q
|    FROM t
|      LEFT JOIN b8 USING (id);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 20
: int64 -1
: int64 0
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| x
: int64 25
: int64 -1
: int64 100
< row
: string
| q
: int64 25
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_39 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_40; SET search_path TO analyze_40
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
| CREATE VIEW v AS SELECT * FROM b8 RIGHT JOIN t USING (id);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     b8.q,
|     t.masked,
|     t.x
|    FROM b8
|      RIGHT JOIN t USING (id);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 20
: int64 -1
: int64 0
< row
: string
| q
: int64 25
: int64 -1
: int64 100
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| x
: int64 25
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_40 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_41; SET search_path TO analyze_41
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
| CREATE VIEW v AS SELECT * FROM t FULL JOIN b8 USING (id);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     t.masked,
|     t.x,
|     b8.q
|    FROM t
|      FULL JOIN b8 USING (id);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| id
: int64 20
: int64 -1
: int64 0
< row
: string
| masked
: int64 25
: int64 -1
: int64 100
< row
: string
| x
: int64 25
: int64 -1
: int64 100
< row
: string
| q
: int64 25
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_41 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_42; SET search_path TO analyze_42
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
| CREATE VIEW v AS SELECT array_agg(DISTINCT x ORDER BY x DESC) AS a, array_agg(DISTINCT lower(x) ORDER BY lower(x)) AS l, array_agg(DISTINCT id::text ORDER BY id::text) AS i FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT array_agg(DISTINCT x ORDER BY x DESC) AS a,
|     array_agg(DISTINCT lower(x) ORDER BY (lower(x))) AS l,
|     array_agg(DISTINCT id::text ORDER BY (id::text)) AS i
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 1009
: int64 -1
: int64 100
< row
: string
| l
: int64 1009
: int64 -1
: int64 100
< row
: string
| i
: int64 1009
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_42 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_43; SET search_path TO analyze_43
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
| CREATE VIEW v AS SELECT string_agg(DISTINCT v, ',' ORDER BY v) AS s FROM vc;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT string_agg(DISTINCT v::text, ','::text ORDER BY (v::text)) AS s
|    FROM vc;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| s
: int64 25
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_43 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_44; SET search_path TO analyze_44
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
| CREATE VIEW v AS SELECT rank(1, 'a') WITHIN GROUP (ORDER BY id, x) AS r FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT rank(1, 'a'::text) WITHIN GROUP (ORDER BY id, x) AS r
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| r
: int64 20
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_44 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_45; SET search_path TO analyze_45
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
| CREATE VIEW v AS SELECT rank(1.5) WITHIN GROUP (ORDER BY id) AS r, percent_rank('5') WITHIN GROUP (ORDER BY x) AS p FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT rank(1.5) WITHIN GROUP (ORDER BY (id::numeric)) AS r,
|     percent_rank('5'::text) WITHIN GROUP (ORDER BY x) AS p
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| r
: int64 20
: int64 -1
: int64 0
< row
: string
| p
: int64 701
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_45 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_46; SET search_path TO analyze_46
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
| CREATE VIEW v AS SELECT * FROM json_to_record('{"a":1}') AS x(a int, b text);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM json_to_record('{"a":1}'::json) x(a integer, b text);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| b
: int64 25
: int64 -1
: int64 100
> exec
| DROP SCHEMA analyze_46 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_47; SET search_path TO analyze_47
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
| CREATE VIEW v AS SELECT x.a, x.b, x.c FROM json_to_recordset('[]') AS x(a int, b varchar(5), c numeric(10, 2));
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b,
|     c
|    FROM json_to_recordset('[]'::json) x(a integer, b character varying(5), c numeric(10,2));
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| b
: int64 1043
: int64 9
: int64 100
< row
: string
| c
: int64 1700
: int64 655366
: int64 0
> exec
| DROP SCHEMA analyze_47 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_48; SET search_path TO analyze_48
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
| CREATE VIEW v AS SELECT * FROM ROWS FROM (json_to_record('{}') AS (a int), generate_series(1, 2)) AS z(a, g);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     g
|    FROM ROWS FROM(json_to_record('{}'::json) AS (a integer), generate_series(1, 2)) z(a, g);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| g
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_48 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_49; SET search_path TO analyze_49
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
| CREATE VIEW v AS SELECT * FROM json_to_record('{}') AS x(a text COLLATE "C", b text COLLATE "default", c text COLLATE "POSIX", d int);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b,
|     c,
|     d
|    FROM json_to_record('{}'::json) x(a text COLLATE "C", b text, c text COLLATE "POSIX", d integer);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 25
: int64 -1
: int64 950
< row
: string
| b
: int64 25
: int64 -1
: int64 100
< row
: string
| c
: int64 25
: int64 -1
: int64 951
< row
: string
| d
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_49 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_50; SET search_path TO analyze_50
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
| CREATE VIEW v AS SELECT a FROM json_to_record('{"a":1}') AS (a int);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM json_to_record('{"a":1}'::json) json_to_record(a integer);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_50 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_51; SET search_path TO analyze_51
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
| CREATE VIEW v AS SELECT poly_def(1) AS r, poly_def(1, 'x'::text) AS s, poly_def(a => 2) AS n;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT poly_def(1) AS r,
|     poly_def(1, 'x'::text) AS s,
|     poly_def(a => 2) AS n;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| r
: int64 23
: int64 -1
: int64 0
< row
: string
| s
: int64 25
: int64 -1
: int64 100
< row
: string
| n
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_51 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_52; SET search_path TO analyze_52
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
| CREATE VIEW v AS SELECT * FROM v1 FULL JOIN v2 USING (k);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT k,
|     v1.a,
|     v2.b
|    FROM v1
|      FULL JOIN v2 USING (k);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| k
: int64 1043
: int64 14
: int64 100
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| b
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_52 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_53; SET search_path TO analyze_53
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
| CREATE VIEW v AS SELECT * FROM v1 FULL JOIN v3 USING (k);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT k,
|     v1.a,
|     v3.c
|    FROM v1
|      FULL JOIN v3 USING (k);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| k
: int64 1043
: int64 -1
: int64 100
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| c
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_53 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_54; SET search_path TO analyze_54
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
| CREATE VIEW v AS SELECT * FROM v1 JOIN v3 USING (k);
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT k,
|     v1.a,
|     v3.c
|    FROM v1
|      JOIN v3 USING (k);
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| k
: int64 1043
: int64 -1
: int64 100
< row
: string
| a
: int64 23
: int64 -1
: int64 0
< row
: string
| c
: int64 23
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_54 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_55; SET search_path TO analyze_55
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
| CREATE VIEW v AS SELECT coalesce(v1.k, v2.k) AS k FROM v1, v2;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT COALESCE(v1.k, v2.k) AS k
|    FROM v1,
|     v2;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| k
: int64 1043
: int64 14
: int64 100
> exec
| DROP SCHEMA analyze_55 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_56; SET search_path TO analyze_56
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
| CREATE VIEW v AS SELECT sum(id) FILTER (WHERE 'true') AS s, count(*) FILTER (WHERE NULL) AS c FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT sum(id) FILTER (WHERE true) AS s,
|     count(*) FILTER (WHERE NULL::boolean) AS c
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| s
: int64 20
: int64 -1
: int64 0
< row
: string
| c
: int64 20
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_56 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_57; SET search_path TO analyze_57
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
| CREATE VIEW v AS SELECT array_agg(1 ORDER BY 'a') AS a FROM t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT array_agg(1 ORDER BY 'a'::text) AS a
|    FROM t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| a
: int64 1007
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_57 CASCADE
< ok 0 -
> exec
| CREATE SCHEMA analyze_58; SET search_path TO analyze_58
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
| CREATE VIEW v AS SELECT string_agg(v, ',' ORDER BY v) AS s, string_agg(v, ',' ORDER BY v || '') AS s2, array_agg(id ORDER BY id + 1) AS a FROM vc, t;
< ok 0 -
> query
| SELECT pg_get_viewdef('v'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT string_agg(vc.v::text, ','::text ORDER BY (vc.v::text)) AS s,
|     string_agg(vc.v::text, ','::text ORDER BY (vc.v::text || ''::text)) AS s2,
|     array_agg(t.id ORDER BY (t.id + 1)) AS a
|    FROM vc,
|     t;
> query
| SELECT attname, atttypid, atttypmod, attcollation FROM pg_attribute
| 				WHERE attrelid = 'v'::regclass AND attnum > 0 ORDER BY attnum
< result
< column "attname" "NAME"
< column "atttypid" "OID"
< column "atttypmod" "INT4"
< column "attcollation" "OID"
< row
: string
| s
: int64 25
: int64 -1
: int64 100
< row
: string
| s2
: int64 25
: int64 -1
: int64 100
< row
: string
| a
: int64 1007
: int64 -1
: int64 0
> exec
| DROP SCHEMA analyze_58 CASCADE
< ok 0 -
