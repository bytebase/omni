-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestSimulationCausesOracle default#1
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a DROP a
< error
| 42809 ALTER action DROP COLUMN cannot be performed on relation "a"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| DROP TABLE a
< error
| 42809 "a" is not a table
| hint: Use DROP INDEX to remove an index.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| TRUNCATE a
< error
| 42809 "a" is not a table
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| COMMENT ON TABLE a IS 's'
< error
| 42809 "a" is not a table
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE INDEX ON a(b)
< error
| 42809 cannot open relation "a"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| REFRESH MATERIALIZED VIEW a
< error
| 42809 "a" is not a table or materialized view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER VIEW a DROP a
< error
| 42809 "a" is not a view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| DROP VIEW a
< error
| 42809 "a" is not a view
| hint: Use DROP INDEX to remove an index.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| DROP MATERIALIZED VIEW a
< error
| 42809 "a" is not a materialized view
| hint: Use DROP INDEX to remove an index.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| DROP TABLE IF EXISTS a
< error
| 42809 "a" is not a table
| hint: Use DROP INDEX to remove an index.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a SET SCHEMA a
< error
| 42809 cannot change schema of index "a"
| hint: Change the schema of the table instead.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER SEQUENCE a CYCLE
< error
| 42809 cannot open relation "a"
| detail: This operation is not supported for tables.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER SEQUENCE a CYCLE
< error
| 42809 cannot open relation "a"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER SEQUENCE a CYCLE
< error
| 42P01 relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER SEQUENCE IF EXISTS a CYCLE
< error
| 42809 cannot open relation "a"
| detail: This operation is not supported for tables.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER SEQUENCE IF EXISTS a CYCLE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| DROP SEQUENCE a
< error
| 42809 "a" is not a sequence
| hint: Use DROP TABLE to remove a table.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| DROP SEQUENCE a
< error
| 42809 "a" is not a sequence
| hint: Use DROP INDEX to remove an index.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP SEQUENCE a
< error
| 42P01 sequence "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER SEQUENCE a SET SCHEMA a
< error
| 42809 "a" is not a sequence
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER INDEX a DROP a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS a DROP a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS a DROP a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| DROP INDEX a
< error
| 42809 "a" is not an index
| hint: Use DROP TABLE to remove a table.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| DROP INDEX IF EXISTS a
< error
| 42809 "a" is not an index
| hint: Use DROP TABLE to remove a table.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP INDEX IF EXISTS a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a CLUSTER ON a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a REPLICA IDENTITY USING INDEX a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD UNIQUE USING INDEX a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD PRIMARY KEY USING INDEX a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE IF EXISTS a SET SCHEMA a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER VIEW a SET SCHEMA a
< error
| 42809 "a" is not a view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| DROP VIEW a
< error
| 42809 "a" is not a view
| hint: Use DROP TABLE to remove a table.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| DROP MATERIALIZED VIEW a
< error
| 42809 "a" is not a materialized view
| hint: Use DROP TABLE to remove a table.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a(a a)
< error
| 42704 type "a" does not exist
| position: 18
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a()
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a(CHECK(a))
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a PARTITION OF a DEFAULT
< error
| 42809 cannot open relation "a"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a() INHERITS(a)
< error
| 42809 cannot open relation "a"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a AS SELECT
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW a AS SELECT
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE TABLE a AS SELECT zzz
< error
| 42703 column "zzz" does not exist
| position: 26
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS a AS SELECT zzz
< error
| 42703 column "zzz" does not exist
| position: 40
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE VIEW a AS SELECT
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE OR REPLACE VIEW a AS SELECT
< error
| 42809 "a" is not a view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE TABLE a OF a
< error
| 42809 type a is not a composite type
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS a OF a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a INHERIT a
< error
| 42P07 circular inheritance not allowed
| detail: "a" is already a child of "a".
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE a SET TABLESPACE pg_default
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE pg_default
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s TABLESPACE a AS SELECT 1
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s USING a AS SELECT 1
< error
| 42704 access method "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s TABLESPACE pg_default AS SELECT 1
< ok 1 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s USING btree AS SELECT 1
< error
| 55000 access method "btree" is not of type TABLE
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a USING heap(b)
< error
| XX000 index access method handler function 3 did not return an IndexAmRoutine struct
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP TABLESPACE pg_default
< error
| 42501 permission denied for tablespace pg_default
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP TABLESPACE pg_global
< error
| 42501 permission denied for tablespace pg_global
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) USING a
< error
| 42704 access method "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) USING btree
< error
| 55000 access method "btree" is not of type TABLE
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX i ON a(b) TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS no_such.i SET (fillfactor = 80)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS no_such.i RENAME TO j
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER INDEX a_d_idx SET TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| DROP TABLE IF EXISTS t, a
< error
| 42809 "a" is not a table
| hint: Use DROP INDEX to remove an index.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS a()
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS a PARTITION OF a DEFAULT
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) TABLESPACE pg_global
< error
| 22023 only shared relations can be placed in pg_global tablespace
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX i ON a(b) TABLESPACE pg_global
< error
| 22023 only shared relations can be placed in pg_global tablespace
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE pg_global
< error
| 22023 cannot move relations in to or out of pg_global tablespace
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE OPERATOR FAMILY f USING heap
< error
| 55000 access method "heap" is not of type INDEX
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE OPERATOR CLASS c FOR TYPE int USING heap AS STORAGE int
< error
| 55000 access method "heap" is not of type INDEX
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER OPERATOR FAMILY f USING heap ADD OPERATOR 1 < (int, int)
< error
| 55000 access method "heap" is not of type INDEX
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLESPACE pg_fast LOCATION '/tmp'
< error
| 42939 unacceptable tablespace name "pg_fast"
| detail: The prefix "pg_" is reserved for system tablespaces.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX i ON a USING heap(b) TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE pg_default OWNED BY a SET TABLESPACE pg_default
< error
| 42704 role "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a OWNER TO missing_role
< error
| 42704 role "missing_role" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER INDEX a RENAME TO a
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLESPACE missing RENAME TO pg_bad
< error
| 42704 tablespace "missing" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER ROLE missing RENAME TO pg_alice
< error
| 42704 role "missing" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER ROLE postgres RENAME TO alice
< error
| 0A000 session user cannot be renamed
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| SET default_tablespace = missing
< error
| 22023 invalid value for parameter "default_tablespace": "missing"
| detail: Tablespace "missing" does not exist.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| SET default_tablespace = pg_default;
| CREATE TABLE p (b int) PARTITION BY RANGE (b)
< error
| 0A000 cannot specify default tablespace for partitioned relations
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP ROLE IF EXISTS alice;
| CREATE ROLE alice;
| CREATE TABLE s_tbl (id serial);
| ALTER SEQUENCE s_tbl_id_seq OWNER TO alice
< error
| 0A000 cannot change owner of sequence "s_tbl_id_seq"
| detail: Sequence "s_tbl_id_seq" is linked to table "s_tbl".
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE p (b int) PARTITION BY RANGE (b);
| CREATE INDEX p_idx ON p (b);
| ALTER TABLE p SET TABLESPACE pg_default;
| ALTER INDEX p_idx SET TABLESPACE pg_default
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) WITH (a)
< error
| 22023 unrecognized parameter "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s WITH (a) AS SELECT 1
< error
| 22023 unrecognized parameter "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET (a)
< error
| 22023 unrecognized parameter "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET (fillfactor = 80, fillfactor = 90)
< error
| 22023 parameter "fillfactor" specified more than once
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET (toast.a = 1)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (t text) WITH (toast.a = 1)
< error
| 22023 unrecognized parameter "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET (a.a = 1)
< error
| 22023 unrecognized parameter namespace "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a SET (a)
< error
| 22023 unrecognized parameter "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a RESET (a)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER a SET (a)
< error
| 42703 column "a" of relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER b SET (a)
< error
| 22023 unrecognized parameter "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER a RESET (a)
< error
| 42703 column "a" of relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER b SET (n_distinct = -0.5)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD a a COMPRESSION a
< error
| 22023 invalid compression method "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD a a COMPRESSION DEFAULT
< error
| 42P16 composite type a cannot be made a member of itself
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD a a STORAGE a
< error
| 22023 invalid storage type "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD a a STORAGE DEFAULT
< error
| 42P16 composite type a cannot be made a member of itself
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER a SET STORAGE a
< error
| 42703 column "a" of relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER b SET STORAGE a
< error
| 22023 invalid storage type "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER b SET STORAGE extended
< error
| 0A000 column data type integer can only have storage PLAIN
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER b SET COMPRESSION a
< error
| 0A000 column data type integer does not support compression
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD x a[]
< error
| 42P16 composite type a cannot be made a member of itself
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER c TYPE a
< error
| 42P16 composite type a cannot be made a member of itself
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER c TYPE a COLLATE a
< error
| 42704 collation "a" for encoding "UTF8" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER c TYPE int COLLATE "C"
< error
| 42804 collations are not supported by type integer
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s ON COMMIT DROP AS SELECT 1
< error
| 42P16 ON COMMIT can only be used on temporary tables
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s AS EXECUTE a
< error
| 26000 prepared statement "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS s AS EXECUTE a
< error
| 26000 prepared statement "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE SEQUENCE a
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE TABLE a (a a, a a)
< error
| 42701 column "a" specified more than once
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| COMMENT ON COLUMN a IS 's'
< error
| 42601 column name must be qualified
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE VIEW a AS SELECT 1 WITH CHECK OPTION
< error
| 0A000 WITH CHECK OPTION is supported only on automatically updatable views
| hint: Views that do not select from a single table or view are not automatically updatable.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE VIEW a AS SELECT 1 WITH CASCADED CHECK OPTION
< error
| 0A000 WITH CHECK OPTION is supported only on automatically updatable views
| hint: Views that do not select from a single table or view are not automatically updatable.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT b FROM a WITH LOCAL CHECK OPTION
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT DISTINCT b FROM a WITH CHECK OPTION
< error
| 0A000 WITH CHECK OPTION is supported only on automatically updatable views
| hint: Views containing DISTINCT are not automatically updatable.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE VIEW v AS SELECT 1;
| ALTER VIEW v SET (check_option = local)
< error
| 0A000 WITH CHECK OPTION is supported only on automatically updatable views
| hint: Views that do not select from a single table or view are not automatically updatable.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| REFRESH MATERIALIZED VIEW a
< error
| 0A000 "a" is not a materialized view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE RECURSIVE VIEW a (a) AS SELECT
< error
| 42P10 WITH query "a" has 0 columns available but 1 columns specified
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE VIEW v (x, y) AS SELECT 1
< error
| 42601 CREATE VIEW specifies more column names than columns
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE IF EXISTS a DROP a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER VIEW IF EXISTS a DROP a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER MATERIALIZED VIEW IF EXISTS a DROP a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX a DROP a
< error
| 42P01 relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX a ATTACH PARTITION a
< error
| 42P01 relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX a RENAME TO a
< error
| 42P01 relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER INDEX a DROP a
< error
| 42809 ALTER action DROP COLUMN cannot be performed on relation "a"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS a DROP a
< error
| 42809 ALTER action DROP COLUMN cannot be performed on relation "a"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER INDEX a ATTACH PARTITION a
< error
| 42809 ALTER action ATTACH PARTITION cannot be performed on relation "a"
| detail: This operation is not supported for indexes.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a RENAME TO a
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a RENAME a TO b
< error
| 42703 column "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER VIEW a RENAME a TO b
< error
| 42703 column "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER MATERIALIZED VIEW a RENAME a TO b
< error
| 42703 column "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a RENAME d TO d2
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a RENAME CONSTRAINT a TO b
< error
| 42704 constraint "a" for table "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE IF EXISTS a RENAME CONSTRAINT a TO b
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER 1 SET STATISTICS 1
< error
| 0A000 cannot refer to non-index column by number
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER 1 SET STATISTICS 1
< error
| 0A000 cannot alter statistics on non-expression column "d" of index "a"
| hint: Alter statistics on table column instead.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER a SET STATISTICS 1
< error
| 42703 column "a" of relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER a SET STATISTICS 1
< error
| 42703 column "a" of relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER CONSTRAINT a
< error
| 42704 constraint "a" of relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ENABLE RULE a
< error
| 42704 rule "a" for relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a DISABLE RULE a
< error
| 42704 rule "a" for relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET WITHOUT OIDS
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET ACCESS METHOD DEFAULT
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET ACCESS METHOD a
< error
| 42704 access method "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET ACCESS METHOD btree
< error
| 55000 access method "btree" is not of type TABLE
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER c SET EXPRESSION AS (b + 1)
< error
| 55000 column "c" of relation "a" is not a generated column
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER c DROP EXPRESSION
< error
| 55000 column "c" of relation "a" is not a stored generated column
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a NO INHERIT a
< error
| 42P01 relation "a" is not a parent of relation "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a DETACH PARTITION a FINALIZE
< error
| 42P01 relation "a" is not a parent of relation "a"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a OF a
< error
| 42809 type a is not a composite type
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a NOT OF
< error
| 42809 "a" is not a typed table
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a OPTIONS (a 's')
< error
| 42809 ALTER action OPTIONS cannot be performed on relation "a"
| detail: This operation is not supported for tables.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ALTER c OPTIONS (a 's')
< error
| 42809 ALTER action ALTER COLUMN ... OPTIONS cannot be performed on relation "a"
| detail: This operation is not supported for tables.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE a (CHECK (a))
< error
| 42804 argument of CHECK must be type boolean, not type a
| position: 24
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD CHECK (a)
< error
| 42804 argument of CHECK must be type boolean, not type a
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD CHECK (b)
< error
| 42804 argument of CHECK must be type boolean, not type integer
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD EXCLUDE (b WITH ~)
< error
| 42883 operator does not exist: integer ~ integer
| hint: No operator matches the given name and argument types. You might need to add explicit type casts.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD EXCLUDE (b WITH OPERATOR(pg_catalog.~))
< error
| 42883 operator does not exist: integer pg_catalog.~ integer
| hint: No operator matches the given name and argument types. You might need to add explicit type casts.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a ((a))
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a ((a) zzz)
< error
| 42704 operator class "zzz" does not exist for access method "btree"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (a(a))
< error
| 42883 function a(a) does not exist
| position: 20
| hint: No function matches the given name and argument types. You might need to add explicit type casts.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (CURRENT_DATE)
< error
| 42P17 functions in index expression must be marked IMMUTABLE
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (CURRENT_DATE zzz)
< error
| 42P17 functions in index expression must be marked IMMUTABLE
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (CURRENT_ROLE)
< error
| 42P17 functions in index expression must be marked IMMUTABLE
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (b COLLATE a)
< error
| 42704 collation "a" for encoding "UTF8" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (b COLLATE "C")
< error
| 42804 collations are not supported by type integer
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (b zzz)
< error
| 42704 operator class "zzz" does not exist for access method "btree"
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (b text_ops)
< error
| 42804 operator class "text_ops" does not accept data type integer
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (t varchar(20));
| CREATE INDEX ON s (t varchar_pattern_ops)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (b) INCLUDE (c zzz)
< error
| 42P17 including column does not support an operator class
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (b, b)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a (b) INCLUDE (b, b)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD UNIQUE (b) INCLUDE (b)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TYPE s_enum AS ENUM ('a');
| CREATE SEQUENCE s_enum
< error
| 42710 type "s_enum" already exists
| hint: A relation has an associated type of the same name, so you must use a name that doesn't conflict with any existing type.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TYPE s_enum AS ENUM ('a');
| CREATE SEQUENCE IF NOT EXISTS s_enum
< error
| 42710 type "s_enum" already exists
| hint: A relation has an associated type of the same name, so you must use a name that doesn't conflict with any existing type.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE DOMAIN s_dom AS int;
| CREATE SEQUENCE s_dom
< error
| 42710 type "s_dom" already exists
| hint: A relation has an associated type of the same name, so you must use a name that doesn't conflict with any existing type.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TYPE s_comp AS (a int);
| CREATE SEQUENCE s_comp
< error
| 42P07 relation "s_comp" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TYPE s_comp AS (a int);
| CREATE SEQUENCE IF NOT EXISTS s_comp
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv_u AS SELECT 1 AS id WITH NO DATA;
| REFRESH MATERIALIZED VIEW CONCURRENTLY mv_u
< error
| 0A000 CONCURRENTLY cannot be used when the materialized view is not populated
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv_p AS SELECT 1 AS id;
| REFRESH MATERIALIZED VIEW CONCURRENTLY mv_p
< error
| 55000 cannot refresh materialized view "public.mv_p" concurrently
| hint: Create a unique index with no WHERE clause on one or more columns of the materialized view.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW mv_ok AS SELECT 1 AS id;
| CREATE UNIQUE INDEX mv_ok_idx ON mv_ok (id);
| REFRESH MATERIALIZED VIEW CONCURRENTLY mv_ok
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE COLLATION c_dep (LOCALE = 'C');
| CREATE TABLE t_dep (a text COLLATE c_dep);
| DROP COLLATION c_dep
< error
| 2BP01 cannot drop collation c_dep because other objects depend on it
| detail: column a of table t_dep depends on collation c_dep
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE COLLATION c_dep (LOCALE = 'C');
| CREATE TABLE t_dep (a int, b text COLLATE c_dep);
| DROP COLLATION c_dep CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| PREPARE q0 AS SELECT 1 AS x;
| CREATE TABLE t0 AS EXECUTE q0(1)
< ok 1 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| PREPARE q_inf AS SELECT $1::int AS x, $2::text AS y;
| CREATE TABLE t_inf AS EXECUTE q_inf(1)
< error
| 42601 wrong number of parameters for prepared statement "q_inf"
| detail: Expected 2 parameters but got 1.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| PREPARE q_dup AS SELECT 1 AS a;
| PREPARE q_dup AS SELECT 2 AS b
< error
| 42P05 prepared statement "q_dup" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE COLLATION c_ret (LOCALE = 'C');
| CREATE TABLE t_ret (a text COLLATE c_ret);
| ALTER TABLE t_ret ALTER a TYPE varchar;
| DROP COLLATION c_ret
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE COLLATION c_def (LOCALE = 'C');
| CREATE TABLE t_def (a text DEFAULT ('x' COLLATE c_def));
| DROP COLLATION c_def
< error
| 2BP01 cannot drop collation c_def because other objects depend on it
| detail: default value for column a of table t_def depends on collation c_def
| hint: Use DROP ... CASCADE to drop the dependent objects too.
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE COLLATION c_def (LOCALE = 'C');
| CREATE TABLE t_def (a text DEFAULT ('x' COLLATE c_def));
| DROP COLLATION c_def CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DEALLOCATE missing_dealloc_stmt
< error
| 26000 prepared statement "missing_dealloc_stmt" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DEALLOCATE ALL;
| DEALLOCATE ALL
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| PREPARE q_r4_inf_add AS SELECT $1 + 1 AS x;
| CREATE TABLE t_r4_inf_add AS EXECUTE q_r4_inf_add(1::bigint)
< ok 1 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TYPE t_u_enum_a_key AS ENUM ('x');
| CREATE TABLE t_u_enum (a int UNIQUE)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE t_ir (a int);
| CREATE INDEX idx_ir ON t_ir ((a + 1));
| ALTER TABLE idx_ir RENAME COLUMN expr TO plus_one
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE t_scd (a text STORAGE EXTERNAL COMPRESSION lz4);
| ALTER TABLE t_scd ALTER COLUMN a SET STORAGE DEFAULT;
| ALTER TABLE t_scd ALTER COLUMN a SET COMPRESSION DEFAULT
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| RESET default_tablespace;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE COLLATION c_same (LOCALE = 'C');
| CREATE TABLE t_same (a text COLLATE c_same DEFAULT ('x' COLLATE c_same));
| ALTER TABLE t_same ALTER COLUMN a DROP DEFAULT;
| ALTER TABLE t_same ALTER COLUMN a TYPE varchar;
| DROP COLLATION c_same
< ok 0 -
