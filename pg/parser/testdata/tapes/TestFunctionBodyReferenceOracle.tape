-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestFunctionBodyReferenceOracle default#1
> exec
| SET check_function_bodies = on
< ok 0 -
> exec
| CREATE TABLE t(id int PRIMARY KEY, v int)
< ok 0 -
> exec
| CREATE TABLE archive(id int, v int)
< ok 0 -
> exec
| CREATE SCHEMA api
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv AS SELECT 1 AS id
< ok 1 -
> exec
| CREATE FUNCTION ref_sql_select(x int) RETURNS int LANGUAGE sql AS $$ SELECT x + 1 $$
< ok 0 -
> exec
| CREATE FUNCTION ref_sql_insert() RETURNS int LANGUAGE sql AS $$ INSERT INTO t(id, v) VALUES (1, 1) RETURNING id $$
< ok 0 -
> exec
| CREATE FUNCTION ref_sql_bad_if() RETURNS int LANGUAGE sql AS $$ IF true THEN SELECT 1; END IF $$
< error
| 42601 syntax error at or near "IF"
| position: 65
> exec
| CREATE FUNCTION ref_plpgsql_returning() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x int; BEGIN INSERT INTO t(id, v) VALUES (2, 2) RETURNING id INTO x; RETURN x; END $$
< ok 0 -
> exec
| CREATE FUNCTION ref_plpgsql_utility() RETURNS void LANGUAGE plpgsql AS $$ BEGIN REFRESH MATERIALIZED VIEW mv; GRANT USAGE ON SCHEMA api TO public; END $$
< ok 0 -
> exec
| CREATE FUNCTION ref_plpgsql_bad_into() RETURNS int LANGUAGE plpgsql AS $$ BEGIN SELECT 1 INTO; END $$
< error
| 42601 syntax error at or near ";"
| position: 94
> exec
| CREATE FUNCTION ref_sql_standard_return() RETURNS int LANGUAGE sql RETURN 1
< ok 0 -
> exec
| CREATE FUNCTION ref_sql_standard_bad() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT FROM; END
< error
| 42601 syntax error at or near ";"
| position: 89
> exec
| CREATE FUNCTION ref_sql_standard_missing_atomic() RETURNS int LANGUAGE sql BEGIN RETURN 1; END
< error
| 42601 syntax error at or near "RETURN"
| position: 82
> exec
| CREATE FUNCTION ref_sql_standard_missing_end() RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN 1;
< error
| 42601 syntax error at end of input
| position: 95
> exec
| CREATE FUNCTION ref_bad_as() RETURNS int AS LANGUAGE sql
< error
| 42601 syntax error at or near "LANGUAGE"
| position: 45
> exec
| CREATE FUNCTION ref_bad_language() RETURNS int LANGUAGE
< error
| 42601 syntax error at end of input
| position: 56
> exec
| CREATE FUNCTION ref_bad_transform() RETURNS int TRANSFORM int LANGUAGE sql AS 'SELECT 1'
< error
| 42601 syntax error at or near "int"
| position: 59
> exec
| CREATE FUNCTION ref_bad_security() RETURNS int SECURITY LANGUAGE sql AS 'SELECT 1'
< error
| 42601 syntax error at or near "LANGUAGE"
| position: 57
> exec
| CREATE FUNCTION ref_bad_set_time() RETURNS int SET TIME 'UTC' LANGUAGE sql AS 'SELECT 1'
< error
| 42601 syntax error at or near "'UTC'"
| position: 57
> exec
| CREATE FUNCTION ref_bad_set_session() RETURNS int SET SESSION ROLE LANGUAGE sql AS 'SELECT 1'
< error
| 42601 syntax error at or near "ROLE"
| position: 63
> exec
| CREATE FUNCTION ref_bad_set_from() RETURNS int SET search_path FROM 'public' LANGUAGE sql AS 'SELECT 1'
< error
| 42601 syntax error at or near "'public'"
| position: 69
