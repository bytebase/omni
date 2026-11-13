-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestEventTriggerEventsAgainstPostgres default#1
> exec
| CREATE TABLE ev_log (event text, tag text)
< ok 0 -
> exec
| CREATE FUNCTION ev_note() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN INSERT INTO ev_log VALUES (lower(tg_event), tg_tag); END $$
< ok 0 -
> exec
| CREATE FOREIGN DATA WRAPPER ev_fdw
< ok 0 -
> exec
| CREATE SERVER ev_srv FOREIGN DATA WRAPPER ev_fdw
< ok 0 -
> exec
| CREATE EVENT TRIGGER ev_end ON ddl_command_end EXECUTE FUNCTION ev_note()
< ok 0 -
> exec
| CREATE EVENT TRIGGER ev_drop ON sql_drop EXECUTE FUNCTION ev_note()
< ok 0 -
> exec
| CREATE EVENT TRIGGER ev_rewrite ON table_rewrite EXECUTE FUNCTION ev_note()
< ok 0 -
> exec
| CREATE EVENT TRIGGER ev_login ON login EXECUTE FUNCTION ev_note()
< ok 0 -
> exec
| CREATE TABLE ev_p (id int PRIMARY KEY, code text, n int DEFAULT 1)
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE TABLE
> exec
| CREATE UNIQUE INDEX ev_p_code ON ev_p (code)
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE INDEX
> exec
| CREATE TABLE ev_t (id int PRIMARY KEY, p_id int, x int)
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE TABLE
> exec
| CREATE VIEW ev_v AS SELECT 1 AS a
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE VIEW
> exec
| ALTER TABLE ev_t ADD COLUMN a int
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ADD COLUMN b int DEFAULT random()::int
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| table_rewrite
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ADD COLUMN c serial
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| table_rewrite
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ADD CONSTRAINT ev_fk FOREIGN KEY (p_id) REFERENCES ev_p
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ADD CONSTRAINT ev_chk CHECK (x > 0)
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| CREATE UNIQUE INDEX ev_x ON ev_t (x)
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE INDEX
> exec
| ALTER TABLE ev_t ADD CONSTRAINT ev_x_key UNIQUE USING INDEX ev_x
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ALTER COLUMN x SET NOT NULL
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ALTER COLUMN x DROP NOT NULL
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ALTER COLUMN x SET DEFAULT 2
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ALTER COLUMN x SET DEFAULT 3, ADD COLUMN d int
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ALTER COLUMN x DROP DEFAULT
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| sql_drop
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ALTER COLUMN a TYPE bigint
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| table_rewrite
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ALTER COLUMN x SET STATISTICS 100
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ADD COLUMN e int DEFAULT 5
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| CREATE DOMAIN ev_pos AS int CHECK (VALUE > 0)
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE DOMAIN
> exec
| ALTER TABLE ev_t ADD COLUMN f ev_pos
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| table_rewrite
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ADD COLUMN g int GENERATED ALWAYS AS (x * 2) STORED
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| table_rewrite
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ADD COLUMN h int GENERATED ALWAYS AS IDENTITY
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| table_rewrite
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ADD COLUMN i text NOT NULL, ADD COLUMN j varchar(10) COLLATE "C", ADD COLUMN k int[], ADD COLUMN l pg_catalog.int8
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t DISABLE TRIGGER ALL, ENABLE TRIGGER USER
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t SET UNLOGGED
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| table_rewrite
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t SET LOGGED
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| table_rewrite
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t OWNER TO CURRENT_USER
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t VALIDATE CONSTRAINT ev_chk
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t ENABLE ROW LEVEL SECURITY
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t DISABLE ROW LEVEL SECURITY
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t RENAME COLUMN a TO aa
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_t RENAME TO ev_t2
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| CREATE SCHEMA ev_s
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE SCHEMA
> exec
| ALTER TABLE ev_t2 SET SCHEMA ev_s
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER SEQUENCE ev_s.ev_t_c_seq OWNED BY NONE
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER SEQUENCE
> exec
| CREATE OR REPLACE VIEW ev_v AS SELECT 1 AS a, 2 AS b
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE VIEW
> exec
| COMMENT ON TABLE ev_p IS 'x'
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| COMMENT
> exec
| GRANT SELECT ON ev_p TO PUBLIC
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| GRANT
> exec
| REVOKE SELECT ON ev_p FROM PUBLIC
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| REVOKE
> exec
| CREATE TYPE ev_mood AS ENUM ('a')
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE TYPE
> exec
| ALTER TYPE ev_mood ADD VALUE 'b'
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER TYPE
> exec
| CREATE DOMAIN ev_d AS int
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE DOMAIN
> exec
| CREATE SEQUENCE ev_q
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE SEQUENCE
> exec
| CREATE FUNCTION ev_f() RETURNS int LANGUAGE sql AS 'SELECT 1'
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE FUNCTION
> exec
| CREATE OR REPLACE FUNCTION ev_f() RETURNS int LANGUAGE sql AS 'SELECT 2'
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE FUNCTION
> exec
| CREATE PROCEDURE ev_pr() LANGUAGE sql AS 'SELECT 1'
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE PROCEDURE
> exec
| DROP PROCEDURE ev_pr()
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| sql_drop
: string
| DROP PROCEDURE
< row
: string
| ddl_command_end
: string
| DROP PROCEDURE
> exec
| CREATE MATERIALIZED VIEW ev_mv AS SELECT 1 AS a
< ok 1 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE MATERIALIZED VIEW
> exec
| ALTER VIEW ev_v RENAME TO ev_v2
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER VIEW
> exec
| ALTER VIEW ev_v2 RENAME TO ev_v
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER VIEW
> exec
| ALTER INDEX ev_p_code RENAME TO ev_p_code2
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| ALTER INDEX
> exec
| CREATE FOREIGN TABLE ev_ft (a int) SERVER ev_srv
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE FOREIGN TABLE
> exec
| CREATE TABLE ev_c AS SELECT 1 AS a
< ok 1 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE TABLE AS
> exec
| TRUNCATE ev_c
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
> exec
| SELECT 1 AS a INTO ev_c2
< ok 1 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| SELECT INTO
> exec
| CREATE INDEX ON ev_p (n)
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE INDEX
> exec
| CREATE SCHEMA ev_s2 CREATE TABLE x (id int UNIQUE) CREATE VIEW y AS SELECT 1 AS a
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| ddl_command_end
: string
| CREATE SCHEMA
> exec
| INSERT INTO ev_p VALUES (1, 'a', 1)
< ok 1 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
> exec
| UPDATE ev_p SET n = 2
< ok 1 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
> exec
| SELECT * FROM ev_p
< ok 1 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
> exec
| EXPLAIN SELECT * FROM ev_p
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
> exec
| SET search_path = public
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
> exec
| ALTER TABLE ev_s.ev_t2 DROP CONSTRAINT ev_chk
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| sql_drop
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| ALTER TABLE ev_s.ev_t2 DROP COLUMN aa
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| sql_drop
: string
| ALTER TABLE
< row
: string
| ddl_command_end
: string
| ALTER TABLE
> exec
| DROP DOMAIN ev_d
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| sql_drop
: string
| DROP DOMAIN
< row
: string
| ddl_command_end
: string
| DROP DOMAIN
> exec
| DROP VIEW ev_v
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| sql_drop
: string
| DROP VIEW
< row
: string
| ddl_command_end
: string
| DROP VIEW
> exec
| DROP TABLE ev_c, ev_c2
< ok 0 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
< row
: string
| sql_drop
: string
| DROP TABLE
< row
: string
| ddl_command_end
: string
| DROP TABLE
> exec
| DELETE FROM ev_p
< ok 1 -
> query
| DELETE FROM ev_log RETURNING event, tag
< result
< column "event" "TEXT"
< column "tag" "TEXT"
