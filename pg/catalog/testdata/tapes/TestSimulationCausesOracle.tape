-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestSimulationCausesOracle default#1
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| TRUNCATE a
< error
| 42809 "a" is not a table
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| COMMENT ON TABLE a IS 's'
< error
| 42809 "a" is not a table
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| REFRESH MATERIALIZED VIEW a
< error
| 42809 "a" is not a table or materialized view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| ALTER VIEW a DROP a
< error
| 42809 "a" is not a view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER SEQUENCE a CYCLE
< error
| 42P01 relation "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER SEQUENCE IF EXISTS a CYCLE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP SEQUENCE a
< error
| 42P01 sequence "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER SEQUENCE a SET SCHEMA a
< error
| 42809 "a" is not a sequence
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER INDEX a DROP a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS a DROP a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS a DROP a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP INDEX IF EXISTS a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a CLUSTER ON a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a REPLICA IDENTITY USING INDEX a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD UNIQUE USING INDEX a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a ADD PRIMARY KEY USING INDEX a
< error
| 42809 "a" is not an index
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE IF EXISTS a SET SCHEMA a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER VIEW a SET SCHEMA a
< error
| 42809 "a" is not a view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a()
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a(CHECK(a))
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE a AS SELECT
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE MATERIALIZED VIEW a AS SELECT
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE VIEW a AS SELECT
< error
| 42P07 relation "a" already exists
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE OR REPLACE VIEW a AS SELECT
< error
| 42809 "a" is not a view
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE TABLE a OF a
< error
| 42809 type a is not a composite type
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS a OF a
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE a SET TABLESPACE pg_default
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE pg_default
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER TABLE a SET TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s TABLESPACE a AS SELECT 1
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s USING a AS SELECT 1
< error
| 42704 access method "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s TABLESPACE pg_default AS SELECT 1
< ok 1 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s USING btree AS SELECT 1
< error
| 55000 access method "btree" is not of type TABLE
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX ON a USING heap(b)
< error
| XX000 index access method handler function 3 did not return an IndexAmRoutine struct
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP TABLESPACE pg_default
< error
| 42501 permission denied for tablespace pg_default
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| DROP TABLESPACE pg_global
< error
| 42501 permission denied for tablespace pg_global
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) USING a
< error
| 42704 access method "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) USING btree
< error
| 55000 access method "btree" is not of type TABLE
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX i ON a(b) TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS no_such.i SET (fillfactor = 80)
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER INDEX IF EXISTS no_such.i RENAME TO j
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| ALTER INDEX a_d_idx SET TABLESPACE a
< error
| 42704 tablespace "a" does not exist
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS a()
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS a PARTITION OF a DEFAULT
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE TABLE s (b int) TABLESPACE pg_global
< error
| 22023 only shared relations can be placed in pg_global tablespace
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX i ON a(b) TABLESPACE pg_global
< error
| 22023 only shared relations can be placed in pg_global tablespace
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE pg_global
< error
| 22023 cannot move relations in to or out of pg_global tablespace
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE OPERATOR FAMILY f USING heap
< error
| 55000 access method "heap" is not of type INDEX
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| CREATE OPERATOR CLASS c FOR TYPE int USING heap AS STORAGE int
< error
| 55000 access method "heap" is not of type INDEX
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
| CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)
< ok 0 -
> exec
| ALTER OPERATOR FAMILY f USING heap ADD OPERATOR 1 < (int, int)
< error
| 55000 access method "heap" is not of type INDEX
> exec
| DROP SCHEMA IF EXISTS public CASCADE;
| CREATE SCHEMA public;
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
| CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)
< ok 0 -
> exec
| CREATE INDEX i ON a USING heap(b) TABLESPACE a
< error
| 42704 tablespace "a" does not exist
