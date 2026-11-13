-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestAnalyzeViewDependenciesMatchPostgres pg18#1
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
| CREATE VIEW v AS SELECT sum(id) FILTER (WHERE masked = 'a') FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN masked
< error
| 2BP01 cannot drop column masked of table t because other objects depend on it
| detail: view v depends on column masked of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
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
| CREATE VIEW v AS SELECT string_agg(x, ',' ORDER BY masked) FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN masked
< error
| 2BP01 cannot drop column masked of table t because other objects depend on it
| detail: view v depends on column masked of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
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
| CREATE VIEW v AS SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY id) FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< error
| 2BP01 cannot drop column id of table t because other objects depend on it
| detail: view v depends on column id of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
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
| CREATE VIEW v AS SELECT y FROM t NATURAL JOIN u;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN id
< error
| 2BP01 cannot drop column id of table t because other objects depend on it
| detail: view v depends on column id of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
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
| CREATE VIEW v AS SELECT y FROM t NATURAL JOIN u;
< ok 0 -
> exec
| ALTER TABLE u DROP COLUMN id
< error
| 2BP01 cannot drop column id of table u because other objects depend on it
| detail: view v depends on column id of table u
| hint: Use DROP ... CASCADE to drop the dependent objects too.
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
| CREATE VIEW v AS SELECT id FROM t FULL JOIN u USING (id);
< ok 0 -
> exec
| ALTER TABLE u DROP COLUMN id
< error
| 2BP01 cannot drop column id of table u because other objects depend on it
| detail: view v depends on column id of table u
| hint: Use DROP ... CASCADE to drop the dependent objects too.
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
| CREATE VIEW v AS SELECT sum(id) FILTER (WHERE masked = 'a') FROM t;
< ok 0 -
> exec
| ALTER TABLE t DROP COLUMN x
< error
| 2BP01 cannot drop column x of table t because other objects depend on it
| detail: view cc depends on column x of table t
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| DROP SCHEMA analyze_7 CASCADE
< ok 0 -
