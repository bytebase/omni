-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestDDLFormsMatchPG default#1
> exec
| SELECT 1/0;
| CREATE TABLE t (a int
< error
| 42601 syntax error at end of input
| position: 34
> exec
| SELECT 1/0;
| CREATE TABLE t (a int) INHERITS (p
< error
| 42601 syntax error at end of input
| position: 47
> exec
| SELECT 1/0;
| CREATE TABLE t (a, ) AS SELECT 1, 2
< error
| 42601 syntax error at or near ")"
| position: 32
> exec
| SELECT 1/0;
| CREATE TABLE t (a int REFERENCES p ON UPDATE CASCADE ON UPDATE CASCADE)
< error
| 42601 syntax error at or near "UPDATE"
| position: 69
> exec
| SELECT 1/0;
| CREATE TABLE t (a int REFERENCES p ON DELETE SET NULL (a) ON UPDATE SET DEFAULT)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a int REFERENCES p ON UPDATE SET NULL (a))
< error
| 0A000 a column list with SET NULL is only supported for ON DELETE actions
| position: 48
> exec
| SELECT 1/0;
| CREATE TABLE t (a int REFERENCES p MATCH PARTIAL)
< error
| 0A000 MATCH PARTIAL not yet implemented
| position: 48
> exec
| SELECT 1/0;
| CREATE TABLE t (a int) PARTITION BY LIST (a)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TABLE t (a int) PARTITION BY foo (a)
< error
| 22023 unrecognized partitioning strategy "foo"
> exec
| SELECT 1/0;
| CREATE TABLE t (LIKE p INCLUDING)
< error
| 42601 syntax error at or near ")"
| position: 45
> exec
| SELECT 1/0;
| CREATE TABLE t (a int UNIQUE NULLS)
< error
| 42601 syntax error at or near ")"
| position: 47
> exec
| SELECT 1/0;
| CREATE TABLE t (a int GENERATED ALWAYS AS IDENTITY ())
< error
| 42601 syntax error at or near ")"
| position: 65
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, EXCLUDE USING (a WITH =))
< error
| 42601 syntax error at or near "("
| position: 50
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, UNIQUE (a) USING INDEX TABLESPACE)
< error
| 42601 syntax error at or near ")"
| position: 69
> exec
| SELECT 1/0;
| CREATE INDEX ON t (a COLLATE)
< error
| 42601 syntax error at or near ")"
| position: 41
> exec
| SELECT 1/0;
| CREATE SEQUENCE s CACHE
< error
| 42601 syntax error at end of input
| position: 36
> exec
| SELECT 1/0;
| CREATE SEQUENCE s NO
< error
| 42601 syntax error at end of input
| position: 33
> exec
| SELECT 1/0;
| ALTER SEQUENCE s CACHE 1 NO INHERIT
< error
| 42601 syntax error at or near "INHERIT"
| position: 41
> exec
| SELECT 1/0;
| ALTER SEQUENCE s RESTART WITH
< error
| 42601 syntax error at end of input
| position: 42
> exec
| SELECT 1/0;
| ALTER SEQUENCE s RESTART CACHE 1
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER TABLE t ALTER a RESTART WITH
< error
| 42601 syntax error at end of input
| position: 47
> exec
| SELECT 1/0;
| CREATE FUNCTION f(a int,) RETURNS int AS 'x' LANGUAGE sql
< error
| 42601 syntax error at or near ")"
| position: 37
> exec
| SELECT 1/0;
| CREATE FUNCTION f() RETURNS TABLE (a) AS 'x' LANGUAGE sql
< error
| 42601 syntax error at or near ")"
| position: 49
> exec
| SELECT 1/0;
| DROP AGGREGATE a
< error
| 42601 syntax error at end of input
| position: 29
> exec
| SELECT 1/0;
| COMMENT ON AGGREGATE a(ORDER BY int) IS NULL
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| COMMENT ON OPERATOR ~(int) IS NULL
< error
| 42601 missing argument
| position: 38
| hint: Use NONE to denote the missing argument of a unary operator.
> exec
| SELECT 1/0;
| CREATE POLICY p ON t AS foo
< error
| 42601 unrecognized row security option "foo"
| position: 37
| hint: Only PERMISSIVE or RESTRICTIVE policies are supported currently.
> exec
| SELECT 1/0;
| CREATE POLICY p ON t AS "permissive" FOR SELECT
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE POLICY p ON t FOR
< error
| 42601 syntax error at end of input
| position: 37
> exec
| SELECT 1/0;
| ALTER POLICY p ON t USING ()
< error
| 42601 syntax error at or near ")"
| position: 40
> exec
| SELECT 1/0;
| ALTER ROLE r SET LOCAL x = 1
< error
| 42601 syntax error at or near "x"
| position: 36
> exec
| SELECT 1/0;
| ALTER ROLE r IN DATABASE d
< error
| 42601 syntax error at end of input
| position: 39
> exec
| SELECT 1/0;
| ALTER ROLE ALL IN DATABASE d RESET ALL
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| ALTER GROUP g RENAME TO
< error
| 42601 syntax error at end of input
| position: 36
> exec
| SELECT 1/0;
| DROP DATABASE d WITH
< error
| 42601 syntax error at end of input
| position: 33
> exec
| SELECT 1/0;
| DROP DATABASE d (FORCE, FORCE)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| DROP DATABASE d (foo)
< error
| 42601 syntax error at or near "foo"
| position: 30
> exec
| SELECT 1/0;
| SET CATALOG 'x'
< error
| 0A000 current database cannot be changed
| position: 25
> exec
| SELECT 1/0;
| SET SESSION CHARACTERISTICS AS TRANSACTION
< error
| 42601 syntax error at end of input
| position: 55
> exec
| SELECT 1/0;
| BEGIN ISOLATION LEVEL
< error
| 42601 syntax error at end of input
| position: 34
> exec
| SELECT 1/0;
| ALTER DOMAIN d SET DEFAULT
< error
| 42601 syntax error at end of input
| position: 39
> exec
| SELECT 1/0;
| ALTER TABLE t ALTER a TYPE int USING
< error
| 42601 syntax error at end of input
| position: 49
> exec
| SELECT 1/0;
| CREATE DOMAIN d int INITIALLY
< error
| 42601 syntax error at end of input
| position: 42
