package catalog

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/internal/testcontainer"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// TestSimulationCausesOracle holds the relation-namespace refusals
// testdata/simulate.tsv records to postgres:17.12-alpine cause by cause.
// PostgreSQL keeps tables, views, materialized views, sequences and indexes
// in one namespace (pg_class): a name a statement resolves to another kind is
// 42809, not a missing relation, and a name no relation holds is 42P01 for a
// sequence as for a table. IF EXISTS skips a missing name only; it still
// refuses one that holds another kind. Each case is answered by the pinned
// server and by the catalog, and both must give the code the case states.
func TestSimulationCausesOracle(t *testing.T) {
	conn := testcontainer.Open(t)
	ctx := context.Background()

	fixtures := map[string]string{
		// an index named a on a table t, as one_table@index.
		"index": "CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE INDEX a ON t (d)",
		// a table named a, with an index of another name, as one_table@table.
		"table": "CREATE TABLE a (b int PRIMARY KEY, c int, d int); CREATE INDEX a_d_idx ON a (d)",
		// neither, as one_table@missing.
		"missing": "CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE INDEX t_d_idx ON t (d)",
	}
	cases := []struct {
		name    string
		fixture string
		stmt    string
		code    string
	}{
		{"ALTER TABLE on an index", "index", "ALTER TABLE a DROP a", CodeWrongObjectType},
		{"DROP TABLE on an index", "index", "DROP TABLE a", CodeWrongObjectType},
		{"TRUNCATE on an index", "index", "TRUNCATE a", CodeWrongObjectType},
		{"COMMENT ON TABLE on an index", "index", "COMMENT ON TABLE a IS 's'", CodeWrongObjectType},
		{"CREATE INDEX on an index", "index", "CREATE INDEX ON a(b)", CodeWrongObjectType},
		{"REFRESH MATERIALIZED VIEW on an index", "index", "REFRESH MATERIALIZED VIEW a", CodeWrongObjectType},
		{"ALTER VIEW on an index", "index", "ALTER VIEW a DROP a", CodeWrongObjectType},
		{"DROP VIEW on an index", "index", "DROP VIEW a", CodeWrongObjectType},
		{"DROP MATERIALIZED VIEW on an index", "index", "DROP MATERIALIZED VIEW a", CodeWrongObjectType},
		{"DROP TABLE IF EXISTS on an index", "index", "DROP TABLE IF EXISTS a", CodeWrongObjectType},
		{"ALTER TABLE SET SCHEMA on an index", "index", "ALTER TABLE a SET SCHEMA a", CodeWrongObjectType},

		{"ALTER SEQUENCE on a table", "table", "ALTER SEQUENCE a CYCLE", CodeWrongObjectType},
		{"ALTER SEQUENCE on an index", "index", "ALTER SEQUENCE a CYCLE", CodeWrongObjectType},
		{"ALTER SEQUENCE on nothing", "missing", "ALTER SEQUENCE a CYCLE", CodeUndefinedTable},
		{"ALTER SEQUENCE IF EXISTS on a table", "table", "ALTER SEQUENCE IF EXISTS a CYCLE", CodeWrongObjectType},
		{"ALTER SEQUENCE IF EXISTS on nothing", "missing", "ALTER SEQUENCE IF EXISTS a CYCLE", "0"},
		{"DROP SEQUENCE on a table", "table", "DROP SEQUENCE a", CodeWrongObjectType},
		{"DROP SEQUENCE on an index", "index", "DROP SEQUENCE a", CodeWrongObjectType},
		{"DROP SEQUENCE on nothing", "missing", "DROP SEQUENCE a", CodeUndefinedTable},
		{"ALTER SEQUENCE SET SCHEMA on a table", "table", "ALTER SEQUENCE a SET SCHEMA a", CodeWrongObjectType},

		{"ALTER INDEX on a table", "table", "ALTER INDEX a DROP a", CodeWrongObjectType},
		{"ALTER INDEX IF EXISTS on a table", "table", "ALTER INDEX IF EXISTS a DROP a", CodeWrongObjectType},
		{"ALTER INDEX IF EXISTS on nothing", "missing", "ALTER INDEX IF EXISTS a DROP a", "0"},
		{"DROP INDEX on a table", "table", "DROP INDEX a", CodeWrongObjectType},
		{"DROP INDEX IF EXISTS on a table", "table", "DROP INDEX IF EXISTS a", CodeWrongObjectType},
		{"DROP INDEX IF EXISTS on nothing", "missing", "DROP INDEX IF EXISTS a", "0"},

		{"ALTER TABLE CLUSTER ON a name a table holds", "table", "ALTER TABLE a CLUSTER ON a", CodeWrongObjectType},
		{"ALTER TABLE REPLICA IDENTITY USING INDEX a name a table holds", "table", "ALTER TABLE a REPLICA IDENTITY USING INDEX a", CodeWrongObjectType},
		{"ALTER TABLE ADD UNIQUE USING INDEX a name a table holds", "table", "ALTER TABLE a ADD UNIQUE USING INDEX a", CodeWrongObjectType},
		{"ALTER TABLE ADD PRIMARY KEY USING INDEX a name a table holds", "table", "ALTER TABLE a ADD PRIMARY KEY USING INDEX a", CodeWrongObjectType},

		{"ALTER TABLE IF EXISTS SET SCHEMA on nothing", "missing", "ALTER TABLE IF EXISTS a SET SCHEMA a", "0"},
		{"ALTER VIEW SET SCHEMA on a table", "table", "ALTER VIEW a SET SCHEMA a", CodeWrongObjectType},
		{"DROP VIEW on a table", "table", "DROP VIEW a", CodeWrongObjectType},
		{"DROP MATERIALIZED VIEW on a table", "table", "DROP MATERIALIZED VIEW a", CodeWrongObjectType},

		// A creation reads the body and the column types before it checks the
		// target name, and a name an index holds is not free.
		{"CREATE TABLE with a column of a missing type on an index", "index", "CREATE TABLE a(a a)", CodeUndefinedObject},
		{"CREATE TABLE with no columns on an index", "index", "CREATE TABLE a()", CodeDuplicateTable},
		{"CREATE TABLE with a check on an index", "index", "CREATE TABLE a(CHECK(a))", CodeDuplicateTable},
		{"CREATE TABLE PARTITION OF an index", "index", "CREATE TABLE a PARTITION OF a DEFAULT", CodeWrongObjectType},
		{"CREATE TABLE INHERITS an index", "index", "CREATE TABLE a() INHERITS(a)", CodeWrongObjectType},
		{"CREATE TABLE AS SELECT on an index", "index", "CREATE TABLE a AS SELECT", CodeDuplicateTable},
		{"CREATE MATERIALIZED VIEW AS SELECT on an index", "index", "CREATE MATERIALIZED VIEW a AS SELECT", CodeDuplicateTable},
		{"CREATE TABLE AS SELECT over a missing column on a table", "table", "CREATE TABLE a AS SELECT zzz", CodeUndefinedColumn},
		{"CREATE TABLE IF NOT EXISTS AS SELECT over a missing column on a table", "table", "CREATE TABLE IF NOT EXISTS a AS SELECT zzz", CodeUndefinedColumn},
		{"CREATE VIEW on an index", "index", "CREATE VIEW a AS SELECT", CodeDuplicateTable},
		{"CREATE OR REPLACE VIEW on an index", "index", "CREATE OR REPLACE VIEW a AS SELECT", CodeWrongObjectType},
		{"CREATE TABLE OF a table's name", "table", "CREATE TABLE a OF a", CodeWrongObjectType},
		{"CREATE TABLE IF NOT EXISTS OF a table's name", "table", "CREATE TABLE IF NOT EXISTS a OF a", "0"},
		{"ALTER TABLE INHERIT itself", "table", "ALTER TABLE a INHERIT a", CodeDuplicateTable},

		// A tablespace or table access method PostgreSQL does not hold is
		// 42704, and CREATE TABLESPACE makes one the next statement may name.
		{"ALTER TABLE ALL IN TABLESPACE with a missing target", "missing", "ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE a", CodeUndefinedObject},
		{"ALTER TABLE ALL IN TABLESPACE with a missing source", "missing", "ALTER TABLE ALL IN TABLESPACE a SET TABLESPACE pg_default", CodeUndefinedObject},
		{"ALTER TABLE ALL IN TABLESPACE between the built-in tablespaces", "missing", "ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE pg_default", "0"},
		{"ALTER TABLE SET TABLESPACE to a missing tablespace", "table", "ALTER TABLE a SET TABLESPACE a", CodeUndefinedObject},
		{"CREATE TABLE AS SELECT into a missing tablespace", "missing", "CREATE TABLE s TABLESPACE a AS SELECT 1", CodeUndefinedObject},
		{"CREATE TABLE AS SELECT using a missing access method", "missing", "CREATE TABLE s USING a AS SELECT 1", CodeUndefinedObject},
		{"CREATE TABLE AS SELECT into the default tablespace", "missing", "CREATE TABLE s TABLESPACE pg_default AS SELECT 1", "0"},
		{"CREATE TABLE AS SELECT using an index access method", "missing", "CREATE TABLE s USING btree AS SELECT 1", CodeObjectNotInPrerequisiteState},
		{"CREATE INDEX using a table access method", "table", "CREATE INDEX ON a USING heap(b)", CodeInternalError},
		{"DROP the default tablespace", "missing", "DROP TABLESPACE pg_default", CodeInsufficientPrivilege},
		{"DROP the global tablespace", "missing", "DROP TABLESPACE pg_global", CodeInsufficientPrivilege},
		{"CREATE TABLE into a missing tablespace", "missing", "CREATE TABLE s (b int) TABLESPACE a", CodeUndefinedObject},
		{"CREATE TABLE using a missing access method", "missing", "CREATE TABLE s (b int) USING a", CodeUndefinedObject},
		{"CREATE TABLE using an index access method", "missing", "CREATE TABLE s (b int) USING btree", CodeObjectNotInPrerequisiteState},
		{"CREATE INDEX into a missing tablespace", "table", "CREATE INDEX i ON a(b) TABLESPACE a", CodeUndefinedObject},
		{"ALTER INDEX IF EXISTS in a missing schema", "missing", "ALTER INDEX IF EXISTS no_such.i SET (fillfactor = 80)", "0"},
		{"ALTER INDEX IF EXISTS RENAME in a missing schema", "missing", "ALTER INDEX IF EXISTS no_such.i RENAME TO j", "0"},
		{"ALTER INDEX SET TABLESPACE to a missing tablespace", "table", "ALTER INDEX a_d_idx SET TABLESPACE a", CodeUndefinedObject},
		{"DROP TABLE IF EXISTS a valid table and an index", "index", "DROP TABLE IF EXISTS t, a", CodeWrongObjectType},
		{"CREATE TABLE IF NOT EXISTS on an index", "index", "CREATE TABLE IF NOT EXISTS a()", "0"},
		{"CREATE TABLE IF NOT EXISTS PARTITION OF an index", "index", "CREATE TABLE IF NOT EXISTS a PARTITION OF a DEFAULT", "0"},
		{"CREATE TABLE in pg_global", "missing", "CREATE TABLE s (b int) TABLESPACE pg_global", CodeInvalidParameterValue},
		{"CREATE INDEX in pg_global", "table", "CREATE INDEX i ON a(b) TABLESPACE pg_global", CodeInvalidParameterValue},
		{"ALTER TABLE ALL IN TABLESPACE to pg_global", "missing", "ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE pg_global", CodeInvalidParameterValue},
		{"CREATE OPERATOR FAMILY using a table access method", "missing", "CREATE OPERATOR FAMILY f USING heap", CodeObjectNotInPrerequisiteState},
		{"CREATE OPERATOR CLASS using a table access method", "missing", "CREATE OPERATOR CLASS c FOR TYPE int USING heap AS STORAGE int", CodeObjectNotInPrerequisiteState},
		{"ALTER OPERATOR FAMILY using a table access method", "missing", "ALTER OPERATOR FAMILY f USING heap ADD OPERATOR 1 < (int, int)", CodeObjectNotInPrerequisiteState},
		{"CREATE TABLESPACE with a reserved name", "missing", "CREATE TABLESPACE pg_fast LOCATION '/tmp'", CodeReservedName},
		{"CREATE INDEX with a missing tablespace and a table access method", "table", "CREATE INDEX i ON a USING heap(b) TABLESPACE a", CodeUndefinedObject},
		{"ALTER TABLE ALL IN TABLESPACE OWNED BY a missing role", "missing", "ALTER TABLE ALL IN TABLESPACE pg_default OWNED BY a SET TABLESPACE pg_default", CodeUndefinedObject},
		{"ALTER TABLE OWNER TO a missing role", "table", "ALTER TABLE a OWNER TO missing_role", CodeUndefinedObject},
		{"ALTER INDEX RENAME to the table's own name", "table", "ALTER INDEX a RENAME TO a", CodeDuplicateTable},
		{"ALTER TABLESPACE missing RENAME TO a reserved name", "missing", "ALTER TABLESPACE missing RENAME TO pg_bad", CodeUndefinedObject},
		{"ALTER ROLE missing RENAME TO a reserved name", "missing", "ALTER ROLE missing RENAME TO pg_alice", CodeUndefinedObject},
		{"ALTER ROLE current user RENAME TO another name", "missing", "ALTER ROLE postgres RENAME TO alice", CodeFeatureNotSupported},
		{"SET default_tablespace to a missing tablespace", "missing", "SET default_tablespace = missing", CodeInvalidParameterValue},
		{"CREATE partitioned table under SET default_tablespace = pg_default", "missing", "SET default_tablespace = pg_default;\nCREATE TABLE p (b int) PARTITION BY RANGE (b)", CodeFeatureNotSupported},
		{"ALTER SEQUENCE OWNER TO another role on a serial sequence", "missing", "DROP ROLE IF EXISTS alice;\nCREATE ROLE alice;\nCREATE TABLE s_tbl (id serial);\nALTER SEQUENCE s_tbl_id_seq OWNER TO alice", CodeFeatureNotSupported},
		{"ALTER partitioned table and index SET TABLESPACE pg_default", "missing", "CREATE TABLE p (b int) PARTITION BY RANGE (b);\nCREATE INDEX p_idx ON p (b);\nALTER TABLE p SET TABLESPACE pg_default;\nALTER INDEX p_idx SET TABLESPACE pg_default", "0"},

		// Reloptions, storage, compression, collations, and composite self-references.
		{"CREATE TABLE with unknown reloption", "missing", "CREATE TABLE s (b int) WITH (a)", CodeInvalidParameterValue},
		{"CREATE TABLE AS SELECT with unknown reloption", "missing", "CREATE TABLE s WITH (a) AS SELECT 1", CodeInvalidParameterValue},
		{"ALTER TABLE SET unknown reloption", "table", "ALTER TABLE a SET (a)", CodeInvalidParameterValue},
		{"ALTER TABLE SET duplicate reloption", "table", "ALTER TABLE a SET (fillfactor = 80, fillfactor = 90)", CodeInvalidParameterValue},
		{"ALTER TABLE SET toast reloption ignored without toast table", "table", "ALTER TABLE a SET (toast.a = 1)", "0"},
		{"CREATE TABLE with toastable column and unknown toast reloption", "missing", "CREATE TABLE s (t text) WITH (toast.a = 1)", CodeInvalidParameterValue},
		{"ALTER TABLE SET unknown reloption namespace", "table", "ALTER TABLE a SET (a.a = 1)", CodeInvalidParameterValue},
		{"ALTER TABLE SET unknown reloption on index", "index", "ALTER TABLE a SET (a)", CodeInvalidParameterValue},
		{"ALTER TABLE RESET unknown reloption on index", "index", "ALTER TABLE a RESET (a)", "0"},
		{"ALTER TABLE ALTER missing column SET reloption", "table", "ALTER TABLE a ALTER a SET (a)", CodeUndefinedColumn},
		{"ALTER TABLE ALTER column SET unknown reloption", "table", "ALTER TABLE a ALTER b SET (a)", CodeInvalidParameterValue},
		{"ALTER TABLE ALTER missing column RESET reloption", "table", "ALTER TABLE a ALTER a RESET (a)", CodeUndefinedColumn},
		{"ALTER TABLE ALTER column SET valid n_distinct", "table", "ALTER TABLE a ALTER b SET (n_distinct = -0.5)", "0"},
		{"ALTER TABLE ADD column with invalid compression method", "table", "ALTER TABLE a ADD a a COMPRESSION a", CodeInvalidParameterValue},
		{"ALTER TABLE ADD self-typed column with DEFAULT compression", "table", "ALTER TABLE a ADD a a COMPRESSION DEFAULT", CodeInvalidTableDefinition},
		{"ALTER TABLE ADD column with invalid storage", "table", "ALTER TABLE a ADD a a STORAGE a", CodeInvalidParameterValue},
		{"ALTER TABLE ADD self-typed column with DEFAULT storage", "table", "ALTER TABLE a ADD a a STORAGE DEFAULT", CodeInvalidTableDefinition},
		{"ALTER TABLE ALTER missing column SET STORAGE", "table", "ALTER TABLE a ALTER a SET STORAGE a", CodeUndefinedColumn},
		{"ALTER TABLE ALTER column SET invalid storage", "table", "ALTER TABLE a ALTER b SET STORAGE a", CodeInvalidParameterValue},
		{"ALTER TABLE ALTER plain-storage column SET STORAGE extended", "table", "ALTER TABLE a ALTER b SET STORAGE extended", CodeFeatureNotSupported},
		{"ALTER TABLE ALTER integer column SET COMPRESSION", "table", "ALTER TABLE a ALTER b SET COMPRESSION a", CodeFeatureNotSupported},
		{"ALTER TABLE ADD self-typed array column", "table", "ALTER TABLE a ADD x a[]", CodeInvalidTableDefinition},
		{"ALTER TABLE ALTER column TYPE to own row type", "table", "ALTER TABLE a ALTER c TYPE a", CodeInvalidTableDefinition},
		{"ALTER TABLE ALTER column TYPE with missing collation", "table", "ALTER TABLE a ALTER c TYPE a COLLATE a", CodeUndefinedObject},
		{"ALTER TABLE ALTER column TYPE integer with C collation", "table", "ALTER TABLE a ALTER c TYPE int COLLATE \"C\"", CodeDatatypeMismatch},

		// CTAS ON COMMIT / EXECUTE, CREATE SEQUENCE on index, duplicate column, COMMENT ON COLUMN.
		{"CREATE TABLE ON COMMIT DROP AS SELECT on non-temp table", "missing", "CREATE TABLE s ON COMMIT DROP AS SELECT 1", CodeInvalidTableDefinition},
		{"CREATE TABLE AS EXECUTE missing prepared statement", "missing", "CREATE TABLE s AS EXECUTE a", CodeInvalidSQLStatementName},
		{"CREATE TABLE IF NOT EXISTS AS EXECUTE missing prepared statement", "missing", "CREATE TABLE IF NOT EXISTS s AS EXECUTE a", CodeInvalidSQLStatementName},
		{"CREATE SEQUENCE colliding with an index", "index", "CREATE SEQUENCE a", CodeDuplicateTable},
		{"CREATE TABLE with duplicate column names on existing table name", "table", "CREATE TABLE a (a a, a a)", CodeDuplicateColumn},
		{"COMMENT ON unqualified COLUMN", "table", "COMMENT ON COLUMN a IS 's'", CodeSyntaxError},

		// Views and materialized views.
		{"CREATE VIEW AS SELECT WITH CHECK OPTION without base table", "missing", "CREATE VIEW a AS SELECT 1 WITH CHECK OPTION", CodeFeatureNotSupported},
		{"CREATE VIEW AS SELECT WITH CASCADED CHECK OPTION colliding with table", "table", "CREATE VIEW a AS SELECT 1 WITH CASCADED CHECK OPTION", CodeFeatureNotSupported},
		{"CREATE VIEW over single table WITH LOCAL CHECK OPTION", "table", "CREATE VIEW v AS SELECT b FROM a WITH LOCAL CHECK OPTION", "0"},
		{"CREATE VIEW with DISTINCT and WITH CHECK OPTION", "table", "CREATE VIEW v AS SELECT DISTINCT b FROM a WITH CHECK OPTION", CodeFeatureNotSupported},
		{"ALTER VIEW SET check_option on non-updatable view", "missing", "CREATE VIEW v AS SELECT 1;\nALTER VIEW v SET (check_option = local)", CodeFeatureNotSupported},
		{"REFRESH MATERIALIZED VIEW on a regular table", "table", "REFRESH MATERIALIZED VIEW a", CodeFeatureNotSupported},
		{"CREATE RECURSIVE VIEW with more column aliases than target list", "missing", "CREATE RECURSIVE VIEW a (a) AS SELECT", CodeInvalidColumnReference},
		{"CREATE VIEW with more column names than query columns", "missing", "CREATE VIEW v (x, y) AS SELECT 1", CodeSyntaxError},

		// ALTER / RENAME / IF EXISTS / subcommand kinds and prerequisites.
		{"ALTER TABLE IF EXISTS on missing relation", "missing", "ALTER TABLE IF EXISTS a DROP a", "0"},
		{"ALTER VIEW IF EXISTS on missing relation", "missing", "ALTER VIEW IF EXISTS a DROP a", "0"},
		{"ALTER MATERIALIZED VIEW IF EXISTS on missing relation", "missing", "ALTER MATERIALIZED VIEW IF EXISTS a DROP a", "0"},
		{"ALTER INDEX DROP on missing relation", "missing", "ALTER INDEX a DROP a", CodeUndefinedTable},
		{"ALTER INDEX ATTACH PARTITION on missing relation", "missing", "ALTER INDEX a ATTACH PARTITION a", CodeUndefinedTable},
		{"ALTER INDEX RENAME on missing relation", "missing", "ALTER INDEX a RENAME TO a", CodeUndefinedTable},
		{"ALTER INDEX DROP column on an index", "index", "ALTER INDEX a DROP a", CodeWrongObjectType},
		{"ALTER INDEX IF EXISTS DROP column on an index", "index", "ALTER INDEX IF EXISTS a DROP a", CodeWrongObjectType},
		{"ALTER INDEX ATTACH PARTITION on a non-partitioned index", "index", "ALTER INDEX a ATTACH PARTITION a", CodeWrongObjectType},
		{"ALTER TABLE RENAME TO same name on an index", "index", "ALTER TABLE a RENAME TO a", CodeDuplicateTable},
		{"ALTER TABLE RENAME missing column on an index", "index", "ALTER TABLE a RENAME a TO b", CodeUndefinedColumn},
		{"ALTER VIEW RENAME missing column on an index", "index", "ALTER VIEW a RENAME a TO b", CodeUndefinedColumn},
		{"ALTER MATERIALIZED VIEW RENAME missing column on an index", "index", "ALTER MATERIALIZED VIEW a RENAME a TO b", CodeUndefinedColumn},
		{"ALTER TABLE RENAME existing column on an index", "index", "ALTER TABLE a RENAME d TO d2", "0"},
		{"ALTER TABLE RENAME CONSTRAINT on an index", "index", "ALTER TABLE a RENAME CONSTRAINT a TO b", CodeUndefinedObject},
		{"ALTER TABLE IF EXISTS RENAME CONSTRAINT on missing relation", "missing", "ALTER TABLE IF EXISTS a RENAME CONSTRAINT a TO b", "0"},
		{"ALTER TABLE ALTER column number SET STATISTICS on a table", "table", "ALTER TABLE a ALTER 1 SET STATISTICS 1", CodeFeatureNotSupported},
		{"ALTER TABLE ALTER non-expression column number SET STATISTICS on an index", "index", "ALTER TABLE a ALTER 1 SET STATISTICS 1", CodeFeatureNotSupported},
		{"ALTER TABLE ALTER missing column SET STATISTICS on a table", "table", "ALTER TABLE a ALTER a SET STATISTICS 1", CodeUndefinedColumn},
		{"ALTER TABLE ALTER missing column SET STATISTICS on an index", "index", "ALTER TABLE a ALTER a SET STATISTICS 1", CodeUndefinedColumn},
		{"ALTER TABLE ALTER missing CONSTRAINT", "table", "ALTER TABLE a ALTER CONSTRAINT a", CodeUndefinedObject},
		{"ALTER TABLE ENABLE missing RULE", "table", "ALTER TABLE a ENABLE RULE a", CodeUndefinedObject},
		{"ALTER TABLE DISABLE missing RULE", "table", "ALTER TABLE a DISABLE RULE a", CodeUndefinedObject},
		{"ALTER TABLE SET WITHOUT OIDS", "table", "ALTER TABLE a SET WITHOUT OIDS", "0"},
		{"ALTER TABLE SET ACCESS METHOD DEFAULT", "table", "ALTER TABLE a SET ACCESS METHOD DEFAULT", "0"},
		{"ALTER TABLE SET ACCESS METHOD missing", "table", "ALTER TABLE a SET ACCESS METHOD a", CodeUndefinedObject},
		{"ALTER TABLE SET ACCESS METHOD index AM", "table", "ALTER TABLE a SET ACCESS METHOD btree", CodeObjectNotInPrerequisiteState},
		{"ALTER TABLE ALTER non-generated column SET EXPRESSION", "table", "ALTER TABLE a ALTER c SET EXPRESSION AS (b + 1)", CodeObjectNotInPrerequisiteState},
		{"ALTER TABLE ALTER non-generated column DROP EXPRESSION", "table", "ALTER TABLE a ALTER c DROP EXPRESSION", CodeObjectNotInPrerequisiteState},
		{"ALTER TABLE NO INHERIT non-parent", "table", "ALTER TABLE a NO INHERIT a", CodeUndefinedTable},
		{"ALTER TABLE DETACH PARTITION FINALIZE non-partition", "table", "ALTER TABLE a DETACH PARTITION a FINALIZE", CodeUndefinedTable},
		{"ALTER TABLE OF a table row type", "table", "ALTER TABLE a OF a", CodeWrongObjectType},
		{"ALTER TABLE NOT OF untyped table", "table", "ALTER TABLE a NOT OF", CodeWrongObjectType},
		{"ALTER TABLE OPTIONS on non-foreign table", "table", "ALTER TABLE a OPTIONS (a 's')", CodeWrongObjectType},
		{"ALTER TABLE ALTER column OPTIONS on non-foreign table", "table", "ALTER TABLE a ALTER c OPTIONS (a 's')", CodeWrongObjectType},

		// Whole-row CHECK, EXCLUDE operators, and CREATE INDEX checks/naming.
		{"CREATE TABLE with whole-row CHECK", "missing", "CREATE TABLE a (CHECK (a))", CodeDatatypeMismatch},
		{"ALTER TABLE ADD whole-row CHECK", "table", "ALTER TABLE a ADD CHECK (a)", CodeDatatypeMismatch},
		{"ALTER TABLE ADD integer column CHECK", "table", "ALTER TABLE a ADD CHECK (b)", CodeDatatypeMismatch},
		{"ALTER TABLE ADD EXCLUDE with undefined operator for integer", "table", "ALTER TABLE a ADD EXCLUDE (b WITH ~)", CodeUndefinedFunction},
		{"ALTER TABLE ADD EXCLUDE with qualified undefined operator for integer", "table", "ALTER TABLE a ADD EXCLUDE (b WITH OPERATOR(pg_catalog.~))", CodeUndefinedFunction},
		{"CREATE INDEX on parenthesized whole-row var", "table", "CREATE INDEX ON a ((a))", "0"},
		{"CREATE INDEX on parenthesized whole-row var with missing opclass", "table", "CREATE INDEX ON a ((a) zzz)", CodeUndefinedObject},
		{"CREATE INDEX on function call over whole-row var", "table", "CREATE INDEX ON a (a(a))", CodeUndefinedFunction},
		{"CREATE INDEX on CURRENT_DATE", "table", "CREATE INDEX ON a (CURRENT_DATE)", CodeInvalidObjectDefinition},
		{"CREATE INDEX on CURRENT_DATE with missing opclass", "table", "CREATE INDEX ON a (CURRENT_DATE zzz)", CodeInvalidObjectDefinition},
		{"CREATE INDEX on CURRENT_ROLE", "table", "CREATE INDEX ON a (CURRENT_ROLE)", CodeInvalidObjectDefinition},
		{"CREATE INDEX with missing collation", "table", "CREATE INDEX ON a (b COLLATE a)", CodeUndefinedObject},
		{"CREATE INDEX with collation on non-collatable column", "table", "CREATE INDEX ON a (b COLLATE \"C\")", CodeDatatypeMismatch},
		{"CREATE INDEX with missing opclass", "table", "CREATE INDEX ON a (b zzz)", CodeUndefinedObject},
		{"CREATE INDEX with incompatible opclass", "table", "CREATE INDEX ON a (b text_ops)", CodeDatatypeMismatch},
		{"CREATE INDEX with valid builtin pattern opclass", "missing", "CREATE TABLE s (t varchar(20));\nCREATE INDEX ON s (t varchar_pattern_ops)", "0"},
		{"CREATE INDEX INCLUDE with opclass", "table", "CREATE INDEX ON a (b) INCLUDE (c zzz)", CodeInvalidObjectDefinition},
		{"CREATE INDEX with duplicate key columns", "table", "CREATE INDEX ON a (b, b)", "0"},
		{"CREATE INDEX with duplicate INCLUDE columns", "table", "CREATE INDEX ON a (b) INCLUDE (b, b)", "0"},
		{"ALTER TABLE ADD UNIQUE with duplicate INCLUDE column", "table", "ALTER TABLE a ADD UNIQUE (b) INCLUDE (b)", "0"},

		// Sequence vs standalone data type collisions, concurrent matview refresh, and collation drops.
		{"CREATE SEQUENCE colliding with enum type", "missing", "CREATE TYPE s_enum AS ENUM ('a');\nCREATE SEQUENCE s_enum", CodeDuplicateObject},
		{"CREATE SEQUENCE IF NOT EXISTS colliding with enum type", "missing", "CREATE TYPE s_enum AS ENUM ('a');\nCREATE SEQUENCE IF NOT EXISTS s_enum", CodeDuplicateObject},
		{"CREATE SEQUENCE colliding with domain type", "missing", "CREATE DOMAIN s_dom AS int;\nCREATE SEQUENCE s_dom", CodeDuplicateObject},
		{"CREATE SEQUENCE colliding with composite type", "missing", "CREATE TYPE s_comp AS (a int);\nCREATE SEQUENCE s_comp", CodeDuplicateTable},
		{"CREATE SEQUENCE IF NOT EXISTS colliding with composite type", "missing", "CREATE TYPE s_comp AS (a int);\nCREATE SEQUENCE IF NOT EXISTS s_comp", "0"},
		{"REFRESH MATERIALIZED VIEW CONCURRENTLY when unpopulated", "missing", "CREATE MATERIALIZED VIEW mv_u AS SELECT 1 AS id WITH NO DATA;\nREFRESH MATERIALIZED VIEW CONCURRENTLY mv_u", CodeFeatureNotSupported},
		{"REFRESH MATERIALIZED VIEW CONCURRENTLY without unique index", "missing", "CREATE MATERIALIZED VIEW mv_p AS SELECT 1 AS id;\nREFRESH MATERIALIZED VIEW CONCURRENTLY mv_p", CodeObjectNotInPrerequisiteState},
		{"REFRESH MATERIALIZED VIEW CONCURRENTLY with qualifying unique index", "missing", "CREATE MATERIALIZED VIEW mv_ok AS SELECT 1 AS id;\nCREATE UNIQUE INDEX mv_ok_idx ON mv_ok (id);\nREFRESH MATERIALIZED VIEW CONCURRENTLY mv_ok", "0"},
		{"DROP COLLATION restricted by table column", "missing", "CREATE COLLATION c_dep (LOCALE = 'C');\nCREATE TABLE t_dep (a text COLLATE c_dep);\nDROP COLLATION c_dep", CodeDependentObjects},
		{"DROP COLLATION CASCADE drops dependent column", "missing", "CREATE COLLATION c_dep (LOCALE = 'C');\nCREATE TABLE t_dep (a int, b text COLLATE c_dep);\nDROP COLLATION c_dep CASCADE", "0"},
		{"CTAS EXECUTE zero-parameter prepared statement with ignored argument", "missing", "PREPARE q0 AS SELECT 1 AS x;\nCREATE TABLE t0 AS EXECUTE q0(1)", "0"},
		{"CTAS EXECUTE inferred-parameter prepared statement with wrong arity", "missing", "PREPARE q_inf AS SELECT $1::int AS x, $2::text AS y;\nCREATE TABLE t_inf AS EXECUTE q_inf(1)", CodeSyntaxError},
		{"duplicate PREPARE rejected with 42P05", "missing", "PREPARE q_dup AS SELECT 1 AS a;\nPREPARE q_dup AS SELECT 2 AS b", CodeDuplicatePstatement},
		{"ALTER COLUMN TYPE without COLLATE resets collation to type default", "missing", "CREATE COLLATION c_ret (LOCALE = 'C');\nCREATE TABLE t_ret (a text COLLATE c_ret);\nALTER TABLE t_ret ALTER a TYPE varchar;\nDROP COLLATION c_ret", "0"},
		{"DROP COLLATION restricted by column default expression", "missing", "CREATE COLLATION c_def (LOCALE = 'C');\nCREATE TABLE t_def (a text DEFAULT ('x' COLLATE c_def));\nDROP COLLATION c_def", CodeDependentObjects},
		{"DROP COLLATION CASCADE drops column default expression", "missing", "CREATE COLLATION c_def (LOCALE = 'C');\nCREATE TABLE t_def (a text DEFAULT ('x' COLLATE c_def));\nDROP COLLATION c_def CASCADE", "0"},
		{"DEALLOCATE missing prepared statement rejected with 26000", "missing", "DEALLOCATE missing_dealloc_stmt", CodeInvalidSQLStatementName},
		{"DEALLOCATE ALL succeeds when no prepared statements exist", "missing", "DEALLOCATE ALL;\nDEALLOCATE ALL", "0"},
		{"CTAS EXECUTE assignment-coerces bigint into inferred integer parameter", "missing", "PREPARE q_r4_inf_add AS SELECT $1 + 1 AS x;\nCREATE TABLE t_r4_inf_add AS EXECUTE q_r4_inf_add(1::bigint)", "0"},
		{"UNIQUE constraint index reuses standalone enum type name", "missing", "CREATE TYPE t_u_enum_a_key AS ENUM ('x');\nCREATE TABLE t_u_enum (a int UNIQUE)", "0"},
		{"ALTER TABLE on an index RENAME COLUMN renames expression attribute", "missing", "CREATE TABLE t_ir (a int);\nCREATE INDEX idx_ir ON t_ir ((a + 1));\nALTER TABLE idx_ir RENAME COLUMN expr TO plus_one", "0"},
		{"ALTER COLUMN SET STORAGE DEFAULT and SET COMPRESSION DEFAULT", "missing", "CREATE TABLE t_scd (a text STORAGE EXTERNAL COMPRESSION lz4);\nALTER TABLE t_scd ALTER COLUMN a SET STORAGE DEFAULT;\nALTER TABLE t_scd ALTER COLUMN a SET COMPRESSION DEFAULT", "0"},
		{"DROP DEFAULT sharing column explicit collation keeps single dependency", "missing", "CREATE COLLATION c_same (LOCALE = 'C');\nCREATE TABLE t_same (a text COLLATE c_same DEFAULT ('x' COLLATE c_same));\nALTER TABLE t_same ALTER COLUMN a DROP DEFAULT;\nALTER TABLE t_same ALTER COLUMN a TYPE varchar;\nDROP COLLATION c_same", "0"},
		{"CREATE COLLATION FROM missing source rejected with 42704", "missing", "CREATE COLLATION c_bad FROM missing_coll", CodeUndefinedObject},
		{"CREATE COLLATION FROM default rejected with 42P16", "missing", "CREATE COLLATION c_bad_dfl FROM \"default\"", CodeInvalidObjectDefinition},
		{"CREATE COLLATION without locale rejected with 42P16", "missing", "CREATE COLLATION c_no_loc (PROVIDER = icu)", CodeInvalidObjectDefinition},
	}

	// A failed CREATE TABLE must not leave the implicit sequence a SERIAL or
	// identity column made behind.
	t.Run("a failed CREATE TABLE leaves no implicit sequence", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, fixtures["index"]); code != "0" {
			t.Fatalf("catalog refused fixture index with %s", code)
		}
		if code := simulationCatalogCode(t, c, "CREATE TABLE a (id serial)"); code != CodeDuplicateTable {
			t.Fatalf("catalog answered %s, want %s", code, CodeDuplicateTable)
		}
		if c.schemaByName["public"].Sequences["a_id_seq"] != nil {
			t.Fatal("a failed CREATE TABLE left a_id_seq behind")
		}
	})

	// A failed CREATE TABLE must not leave inheritance rows for a child that
	// does not exist.
	t.Run("a failed CREATE TABLE leaves no inheritance row", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, fixtures["index"]); code != "0" {
			t.Fatalf("catalog refused fixture index with %s", code)
		}
		if code := simulationCatalogCode(t, c, "CREATE TABLE a (id int) INHERITS (t)"); code != CodeDuplicateTable {
			t.Fatalf("catalog answered %s, want %s", code, CodeDuplicateTable)
		}
		if len(c.inhEntries) != 0 {
			t.Fatalf("a failed CREATE TABLE left %d inheritance rows", len(c.inhEntries))
		}
	})

	// A DROP whose second target is not a table leaves the first target in
	// place.
	t.Run("a failed DROP leaves the relation before the failure", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, fixtures["index"]); code != "0" {
			t.Fatalf("catalog refused fixture index with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP TABLE IF EXISTS t, a"); code != CodeWrongObjectType {
			t.Fatalf("catalog answered %s, want %s", code, CodeWrongObjectType)
		}
		if c.schemaByName["public"].Relations["t"] == nil {
			t.Fatal("a failed DROP removed t")
		}
	})

	// A table owner transfer reaches the table's indexes, and a new index
	// takes the table's owner.
	t.Run("indexes follow the table owner", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE TABLE a (b int)",
			"CREATE INDEX a_b_idx ON a(b)",
			"ALTER TABLE a OWNER TO alice",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if idx := c.schemaByName["public"].Indexes["a_b_idx"]; idx == nil || idx.OwnerName != "alice" {
			t.Fatal("the table owner transfer did not reach the index")
		}
		if code := simulationCatalogCode(t, c, "CREATE INDEX a_b_idx2 ON a(b)"); code != "0" {
			t.Fatalf("catalog refused the second index with %s", code)
		}
		if idx := c.schemaByName["public"].Indexes["a_b_idx2"]; idx == nil || idx.OwnerName != "alice" {
			t.Fatal("a new index did not take the table's owner")
		}
	})

	// CURRENT_USER in OWNER TO is the effective role, not the marker.
	t.Run("OWNER TO CURRENT_USER records the effective role", func(t *testing.T) {
		c := New()
		c.SetSessionUser("alice")
		for _, stmt := range []string{
			"CREATE TABLE a (b int)",
			"ALTER TABLE a OWNER TO CURRENT_USER",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.OwnerName != "alice" {
			t.Fatalf("OWNER TO CURRENT_USER recorded %q", rel.OwnerName)
		}
	})

	// default_tablespace reaches indexes cloned by LIKE ... INCLUDING INDEXES.
	t.Run("default_tablespace reaches cloned indexes", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"SET default_tablespace = ts",
			"CREATE TABLE src (b int)",
			"CREATE INDEX src_b_idx ON src(b)",
			"CREATE TABLE dst (LIKE src INCLUDING INDEXES)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		dst := c.schemaByName["public"].Relations["dst"]
		if dst == nil {
			t.Fatal("dst was not created")
		}
		for _, idx := range c.indexesByRel[dst.OID] {
			if idx.Tablespace != "ts" {
				t.Fatalf("cloned index %s has tablespace %q, want ts", idx.Name, idx.Tablespace)
			}
		}
	})

	// default_tablespace places a relation with no TABLESPACE clause.
	t.Run("default_tablespace places new relations", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"SET default_tablespace = ts",
			"CREATE TABLE a (b int)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.Tablespace != "ts" {
			t.Fatal("the table did not take default_tablespace")
		}
		if code := simulationCatalogCode(t, c, "DROP TABLESPACE ts"); code != CodeObjectNotInPrerequisiteState {
			t.Fatalf("catalog answered %s, want %s", code, CodeObjectNotInPrerequisiteState)
		}
	})

	// A catalog loaded from a snapshot does not know the server's custom
	// tablespaces, so it does not reject them.
	t.Run("a snapshot catalog allows an unobserved tablespace", func(t *testing.T) {
		c := New()
		c.snapshot = true
		if code := simulationCatalogCode(t, c, "CREATE TABLE a (b int) TABLESPACE fast"); code != "0" {
			t.Fatalf("snapshot catalog refused an unobserved tablespace with %s", code)
		}
	})

	// OWNED BY limits a bulk move to the named owners.
	t.Run("a bulk move honors OWNED BY", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLESPACE ts2 LOCATION '/tmp/ts2'",
			"CREATE TABLE a (b int) TABLESPACE ts",
			"CREATE TABLE b (b int) TABLESPACE ts",
			"ALTER TABLE a OWNER TO alice",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "ALTER TABLE ALL IN TABLESPACE ts OWNED BY alice SET TABLESPACE ts2"); code != "0" {
			t.Fatalf("catalog refused the owned move with %s", code)
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.Tablespace != "ts2" {
			t.Fatal("the owned move missed the owner's table")
		}
		if rel := c.schemaByName["public"].Relations["b"]; rel == nil || rel.Tablespace != "ts" {
			t.Fatal("the owned move moved another owner's table")
		}
	})

	// A partition with no tablespace of its own takes its parent's.
	t.Run("a partition inherits its parent's tablespace", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLE p (b int) PARTITION BY RANGE (b) TABLESPACE ts",
			"CREATE TABLE c PARTITION OF p FOR VALUES FROM (0) TO (10)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if rel := c.schemaByName["public"].Relations["c"]; rel == nil || rel.Tablespace != "ts" {
			t.Fatal("the partition did not inherit the parent's tablespace")
		}
	})

	// A constraint's tablespace is checked before a SERIAL sequence is made.
	t.Run("a constraint tablespace error leaves no implicit sequence", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, "CREATE TABLESPACE ts LOCATION '/tmp/ts'"); code != "0" {
			t.Fatalf("catalog refused CREATE TABLESPACE with %s", code)
		}
		if code := simulationCatalogCode(t, c, "CREATE TABLE a (id serial, UNIQUE(id) USING INDEX TABLESPACE missing)"); code != CodeUndefinedObject {
			t.Fatalf("catalog answered %s, want %s", code, CodeUndefinedObject)
		}
		if c.schemaByName["public"].Sequences["a_id_seq"] != nil {
			t.Fatal("a failed CREATE TABLE left a_id_seq behind")
		}
	})

	// User-defined type ownership reaches DROP ROLE, DROP OWNED and REASSIGN.
	t.Run("type ownership follows the role lifecycle", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"SET ROLE alice",
			"CREATE TYPE mood AS ENUM ('sad')",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE alice"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c, "DROP OWNED BY alice"); code != "0" {
			t.Fatalf("catalog refused DROP OWNED with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE alice"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE after DROP OWNED with %s", code)
		}
	})

	// REASSIGN OWNED and ALTER ROLE RENAME move a type's owner too.
	t.Run("type ownership follows REASSIGN and RENAME", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"SET ROLE alice",
			"CREATE TYPE mood AS ENUM ('sad')",
			"RESET ROLE",
			"REASSIGN OWNED BY alice TO bob",
			"ALTER ROLE bob RENAME TO carol",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP OWNED BY carol"); code != "0" {
			t.Fatalf("catalog refused DROP OWNED BY carol with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE carol"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE carol with %s", code)
		}
	})

	// REASSIGN OWNED and ALTER ROLE RENAME move a function's owner too.
	t.Run("function ownership follows REASSIGN and RENAME", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"SET ROLE alice",
			"CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT 1'",
			"RESET ROLE",
			"REASSIGN OWNED BY alice TO bob",
			"ALTER ROLE bob RENAME TO carol",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP OWNED BY carol"); code != "0" {
			t.Fatalf("catalog refused DROP OWNED BY carol with %s", code)
		}
		if len(c.userProcs) != 0 {
			t.Fatal("DROP OWNED kept the function after REASSIGN and RENAME")
		}
	})

	// DROP OWNED drops a function the role owns.
	t.Run("DROP OWNED drops an owned function", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"SET ROLE alice",
			"CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT 1'",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP OWNED BY alice"); code != "0" {
			t.Fatalf("catalog refused DROP OWNED with %s", code)
		}
		if len(c.userProcs) != 0 {
			t.Fatal("DROP OWNED kept the owned function")
		}
	})

	// A composite type does not occupy a tablespace.
	t.Run("a composite type does not occupy a tablespace", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"SET default_tablespace = ts",
			"CREATE TYPE ct AS (a int)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP TABLESPACE ts"); code != "0" {
			t.Fatalf("catalog refused DROP TABLESPACE ts with %s", code)
		}
	})

	// An invalid tablespace is refused before an earlier subcommand runs.
	t.Run("an invalid tablespace leaves no earlier ALTER mutation", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, "CREATE TABLE a (b int)"); code != "0" {
			t.Fatalf("catalog refused CREATE TABLE with %s", code)
		}
		if code := simulationCatalogCode(t, c, "ALTER TABLE a ADD COLUMN c int, SET TABLESPACE missing"); code != CodeUndefinedObject {
			t.Fatalf("catalog answered %s, want %s", code, CodeUndefinedObject)
		}
		rel := c.schemaByName["public"].Relations["a"]
		if rel == nil {
			t.Fatal("the table is gone")
		}
		if _, exists := rel.colByName["c"]; exists {
			t.Fatal("a failed ALTER TABLE left column c behind")
		}
	})

	// A failed ALTER TABLE restores the owner as well as the tablespace.
	t.Run("a failed ALTER TABLE restores the owner", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLE a (b int)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		before := c.schemaByName["public"].Relations["a"].OwnerName
		if code := simulationCatalogCode(t, c, "ALTER TABLE a OWNER TO alice, CLUSTER ON missing_idx"); code != CodeUndefinedObject {
			t.Fatalf("catalog answered %s, want %s", code, CodeUndefinedObject)
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.OwnerName != before {
			t.Fatalf("a failed ALTER TABLE kept owner %q, want %q", rel.OwnerName, before)
		}
	})

	// A view records its owner, so DROP OWNED can find it.
	t.Run("a view records its owner", func(t *testing.T) {
		c := New()
		c.SetSessionUser("alice")
		for _, stmt := range []string{
			"CREATE TABLE a (b int)",
			"CREATE VIEW v AS SELECT b FROM a",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if rel := c.schemaByName["public"].Relations["v"]; rel == nil || rel.OwnerName != "alice" {
			t.Fatal("the view did not record its owner")
		}
	})

	// A snapshot catalog accepts lifecycle statements for tablespaces it
	// cannot see.
	t.Run("a snapshot catalog allows an unobserved tablespace lifecycle", func(t *testing.T) {
		c := New()
		c.snapshot = true
		for _, stmt := range []string{
			"ALTER TABLESPACE fast RENAME TO slow",
			"DROP TABLESPACE slow",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("snapshot catalog refused %s with %s", stmt, code)
			}
		}
	})

	// A later failing subcommand restores the tablespace of the relation.
	t.Run("a failed ALTER TABLE restores the tablespace", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLE a (b int)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "ALTER TABLE a SET TABLESPACE ts, CLUSTER ON missing_idx"); code != CodeUndefinedObject {
			t.Fatalf("catalog answered %s, want %s", code, CodeUndefinedObject)
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.Tablespace != "" {
			t.Fatal("a failed ALTER TABLE kept the new tablespace")
		}
	})

	// A cloned LIKE index keeps the source's custom tablespace.
	t.Run("LIKE INCLUDING INDEXES keeps the source tablespace", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE old LOCATION '/tmp/old'",
			"CREATE TABLESPACE new LOCATION '/tmp/new'",
			"SET default_tablespace = old",
			"CREATE TABLE src (b int)",
			"CREATE INDEX src_b_idx ON src(b)",
			"SET default_tablespace = new",
			"CREATE TABLE dst (LIKE src INCLUDING INDEXES)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		dst := c.schemaByName["public"].Relations["dst"]
		if dst == nil {
			t.Fatal("dst was not created")
		}
		for _, idx := range c.indexesByRel[dst.OID] {
			if idx.Tablespace != "old" {
				t.Fatalf("cloned index %s has tablespace %q, want old", idx.Name, idx.Tablespace)
			}
		}
	})

	// RESET ALL clears default_tablespace.
	t.Run("RESET ALL clears default_tablespace", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"SET default_tablespace = ts",
			"RESET ALL",
			"CREATE TABLE a (b int)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.Tablespace != "" {
			t.Fatal("RESET ALL left default_tablespace in place")
		}
	})

	// Bulk moves leave temporary relations alone.
	t.Run("a bulk move leaves temporary relations alone", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLESPACE ts2 LOCATION '/tmp/ts2'",
			"CREATE TEMP TABLE a (b int) TABLESPACE ts",
			"ALTER TABLE ALL IN TABLESPACE ts SET TABLESPACE ts2",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		rel := c.GetRelation("pg_temp", "a")
		if rel == nil {
			rel = c.GetRelation("", "a")
		}
		if rel == nil || rel.Tablespace != "ts" {
			t.Fatal("a bulk move changed a temporary relation's tablespace")
		}
	})

	// A temporary CTAS relation is left out of a bulk move too.
	t.Run("a bulk move leaves a temporary CTAS relation alone", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLESPACE ts2 LOCATION '/tmp/ts2'",
			"CREATE TEMP TABLE a TABLESPACE ts AS SELECT 1",
			"ALTER TABLE ALL IN TABLESPACE ts SET TABLESPACE ts2",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		rel := c.GetRelation("pg_temp", "a")
		if rel == nil {
			rel = c.GetRelation("", "a")
		}
		if rel == nil || rel.Persistence != 't' || rel.Tablespace != "ts" {
			t.Fatal("a bulk move changed a temporary CTAS relation")
		}
	})

	// An index with no tablespace uses the database default, not its table's.
	t.Run("an index uses the database default tablespace", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLE a (b int) TABLESPACE ts",
			"CREATE INDEX a_b_idx ON a(b)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if idx := c.schemaByName["public"].Indexes["a_b_idx"]; idx == nil || idx.Tablespace != "" {
			t.Fatal("the index did not use the database default tablespace")
		}
	})

	// A partition's cloned index keeps the parent index's tablespace.
	t.Run("a partition index keeps the parent index's tablespace", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLE p (b int) PARTITION BY RANGE (b) TABLESPACE ts",
			"CREATE INDEX p_b_idx ON p(b) TABLESPACE ts",
			"CREATE TABLE c PARTITION OF p FOR VALUES FROM (0) TO (10)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		part := c.schemaByName["public"].Relations["c"]
		if part == nil {
			t.Fatal("partition c was not created")
		}
		for _, idx := range c.indexesByRel[part.OID] {
			if idx.Tablespace != "ts" {
				t.Fatalf("partition index %s has tablespace %q, want ts", idx.Name, idx.Tablespace)
			}
		}
	})

	// A role that still owns a relation cannot be dropped.
	t.Run("a role that owns a relation cannot be dropped", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE TABLE a (b int)",
			"ALTER TABLE a OWNER TO alice",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE alice"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s", code, CodeDependentObjects)
		}
		if !c.roles["alice"] {
			t.Fatal("a failed DROP ROLE removed alice")
		}
	})

	// DROP OWNED drops the owned relation, REASSIGN moves it, and a role
	// rename follows both the registry and the objects it owns.
	t.Run("role ownership changes follow DROP OWNED, REASSIGN and RENAME", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"CREATE TABLE a (b int)",
			"ALTER TABLE a OWNER TO alice",
			"REASSIGN OWNED BY alice TO bob",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.OwnerName != "bob" {
			t.Fatal("REASSIGN OWNED did not move the table")
		}
		if code := simulationCatalogCode(t, c, "ALTER ROLE bob RENAME TO carol"); code != "0" {
			t.Fatalf("catalog refused ALTER ROLE RENAME with %s", code)
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.OwnerName != "carol" {
			t.Fatal("ALTER ROLE RENAME did not move the table")
		}
		if code := simulationCatalogCode(t, c, "DROP OWNED BY carol"); code != "0" {
			t.Fatalf("catalog refused DROP OWNED with %s", code)
		}
		if c.schemaByName["public"].Relations["a"] != nil {
			t.Fatal("DROP OWNED kept the owned table")
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE carol"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE after DROP OWNED with %s", code)
		}
	})

	t.Run("a reserved role name is refused", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, "CREATE ROLE pg_app"); code != CodeReservedName {
			t.Fatalf("catalog answered %s, want %s", code, CodeReservedName)
		}
	})

	t.Run("the current role cannot be dropped", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, "DROP ROLE postgres"); code != CodeObjectNotInPrerequisiteState {
			t.Fatalf("catalog answered %s, want %s", code, CodeObjectNotInPrerequisiteState)
		}
	})

	// A multi-role DROP resolves every role before removing any.
	t.Run("a failed DROP ROLE leaves the earlier roles", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, "CREATE ROLE alice"); code != "0" {
			t.Fatalf("catalog refused CREATE ROLE with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE alice, missing"); code != CodeUndefinedObject {
			t.Fatalf("catalog answered %s, want %s", code, CodeUndefinedObject)
		}
		if !c.roles["alice"] {
			t.Fatal("a failed DROP ROLE removed alice")
		}
	})

	// A failed constraint tablespace check leaves the column as it was.
	t.Run("a failed primary key tablespace check keeps the column nullable", func(t *testing.T) {
		c := New()
		if code := simulationCatalogCode(t, c, "CREATE TABLE a (b int)"); code != "0" {
			t.Fatalf("catalog refused CREATE TABLE with %s", code)
		}
		if code := simulationCatalogCode(t, c, "ALTER TABLE a ADD PRIMARY KEY (b) USING INDEX TABLESPACE missing"); code != CodeUndefinedObject {
			t.Fatalf("catalog answered %s, want %s", code, CodeUndefinedObject)
		}
		rel := c.schemaByName["public"].Relations["a"]
		if rel == nil || rel.Columns[rel.colByName["b"]].NotNull {
			t.Fatal("a failed ALTER TABLE left b NOT NULL")
		}
	})

	// SET TABLESPACE is only for tables, materialized views and indexes.
	t.Run("SET TABLESPACE on a view", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLE a (b int)",
			"CREATE VIEW v AS SELECT b FROM a",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "ALTER VIEW v SET TABLESPACE pg_default"); code != CodeWrongObjectType {
			t.Fatalf("catalog answered %s, want %s", code, CodeWrongObjectType)
		}
	})

	// A bulk move touches only the object class it names, and treats the
	// empty tablespace as pg_default.
	t.Run("bulk moves touch only their object class", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLESPACE ts2 LOCATION '/tmp/ts2'",
			"CREATE TABLE a (b int) TABLESPACE ts",
			"CREATE INDEX a_b_idx ON a(b) TABLESPACE ts",
			"CREATE TABLE d (b int)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "ALTER INDEX ALL IN TABLESPACE ts SET TABLESPACE pg_default"); code != "0" {
			t.Fatalf("catalog refused the index move with %s", code)
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.Tablespace != "ts" {
			t.Fatalf("the index move changed the table's tablespace")
		}
		if idx := c.schemaByName["public"].Indexes["a_b_idx"]; idx == nil || idx.Tablespace != "" {
			t.Fatalf("the index move did not set the default tablespace")
		}
		if code := simulationCatalogCode(t, c, "ALTER TABLE ALL IN TABLESPACE pg_default SET TABLESPACE ts2"); code != "0" {
			t.Fatalf("catalog refused the table move with %s", code)
		}
		if rel := c.schemaByName["public"].Relations["d"]; rel == nil || rel.Tablespace != "ts2" {
			t.Fatalf("the table move missed a table in the default tablespace")
		}
	})

	// A constraint's backing index keeps the tablespace the constraint named.
	t.Run("a constraint index keeps its tablespace", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLE a (b int, UNIQUE (b) USING INDEX TABLESPACE ts)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		found := false
		for _, idx := range c.schemaByName["public"].Indexes {
			if idx.Tablespace == "ts" {
				found = true
			}
		}
		if !found {
			t.Fatal("the constraint index did not keep its tablespace")
		}
		if code := simulationCatalogCode(t, c, "DROP TABLESPACE ts"); code != CodeObjectNotInPrerequisiteState {
			t.Fatalf("catalog answered %s, want %s", code, CodeObjectNotInPrerequisiteState)
		}
	})

	// Renaming a tablespace updates the objects stored in it.
	t.Run("renaming a tablespace updates its objects", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE old LOCATION '/tmp/old'",
			"CREATE TABLE a (b int) TABLESPACE old",
			"ALTER TABLESPACE old RENAME TO new",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if rel := c.schemaByName["public"].Relations["a"]; rel == nil || rel.Tablespace != "new" {
			t.Fatal("the tablespace rename did not move the table with it")
		}
		if code := simulationCatalogCode(t, c, "DROP TABLESPACE new"); code != CodeObjectNotInPrerequisiteState {
			t.Fatalf("catalog answered %s, want %s", code, CodeObjectNotInPrerequisiteState)
		}
	})

	// A tablespace that still holds a table cannot be dropped.
	t.Run("a non-empty tablespace cannot be dropped", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"CREATE TABLE a (b int) TABLESPACE ts",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP TABLESPACE ts"); code != CodeObjectNotInPrerequisiteState {
			t.Fatalf("catalog answered %s, want %s", code, CodeObjectNotInPrerequisiteState)
		}
		if code := simulationCatalogCode(t, c, "DROP TABLE a"); code != "0" {
			t.Fatalf("catalog refused DROP TABLE a with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP TABLESPACE ts"); code != "0" {
			t.Fatalf("catalog refused the now-empty DROP TABLESPACE ts with %s", code)
		}
	})

	// DROP OWNED preflights the entire deletion set before mutating: internal
	// dependencies among owned tables, views, FK constraints, functions, and
	// types succeed regardless of order, while an external dependency aborts
	// without deleting unrelated owned objects.
	t.Run("DROP OWNED preflights the complete owned set", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"SET ROLE alice",
			"CREATE TYPE mood AS ENUM ('ok')",
			"CREATE FUNCTION check_ok(mood) RETURNS bool LANGUAGE sql AS 'SELECT true'",
			"CREATE DOMAIN valid_mood AS mood CHECK (check_ok(VALUE))",
			"CREATE TABLE parent_t (id int PRIMARY KEY, m valid_mood)",
			"CREATE TABLE child_t (id int REFERENCES parent_t(id))",
			"CREATE VIEW v_owned AS SELECT id FROM parent_t",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP OWNED BY alice RESTRICT"); code != "0" {
			t.Fatalf("catalog refused intra-set DROP OWNED with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE alice"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE alice with %s", code)
		}

		c2 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"SET ROLE alice",
			"CREATE TABLE unrelated_t (id int)",
			"CREATE FUNCTION f_ext() RETURNS int LANGUAGE sql AS 'SELECT 1'",
			"SET ROLE bob",
			"CREATE VIEW v_ext AS SELECT f_ext()",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c2, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c2, "DROP OWNED BY alice RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s", code, CodeDependentObjects)
		}
		if c2.schemaByName["public"].Relations["unrelated_t"] == nil {
			t.Fatal("failed DROP OWNED deleted unrelated_t")
		}
	})

	// Ordinary DROP TYPE clears ownership for enums, domains, and ranges.
	t.Run("DROP TYPE clears enum, domain and range ownership", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"SET ROLE alice",
			"CREATE TYPE e AS ENUM ('a')",
			"CREATE DOMAIN d AS int",
			"CREATE TYPE r AS RANGE (subtype = int4)",
			"RESET ROLE",
			"DROP TYPE e",
			"DROP DOMAIN d",
			"DROP TYPE r",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE alice"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE alice after DROP TYPE with %s", code)
		}
	})

	// Grants, policies, schemas, sequences, types, functions, and tablespaces
	// participate in DROP ROLE, ALTER ... OWNER TO, REASSIGN OWNED, DROP OWNED,
	// and ALTER ROLE RENAME.
	t.Run("role dependencies cover grants, policies, sequences, owners and tablespaces", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"CREATE TABLE t (a int)",
			"GRANT SELECT ON t TO alice",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		// #103 & #112: grant blocks DROP ROLE, and REASSIGN OWNED leaves ACLs on alice.
		if code := simulationCatalogCode(t, c, "REASSIGN OWNED BY alice TO bob"); code != "0" {
			t.Fatalf("catalog refused REASSIGN OWNED with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE alice"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s while grant remains", code, CodeDependentObjects)
		}
		// #104: renaming alice rewrites ACL grantees so carol retains the grant.
		if code := simulationCatalogCode(t, c, "ALTER ROLE alice RENAME TO carol"); code != "0" {
			t.Fatalf("catalog refused ALTER ROLE RENAME with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE carol"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s after rename", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c, "DROP OWNED BY carol"); code != "0" {
			t.Fatalf("catalog refused DROP OWNED BY carol with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE carol"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE carol after DROP OWNED with %s", code)
		}

		// #108, #113, #114: tablespace ownership is preserved on DROP OWNED,
		// updated by ALTER TABLESPACE OWNER TO, and moved by REASSIGN OWNED.
		for _, stmt := range []string{
			"CREATE ROLE ts_owner",
			"CREATE TABLESPACE ts LOCATION '/tmp/ts'",
			"ALTER TABLESPACE ts OWNER TO ts_owner",
			"DROP OWNED BY ts_owner",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if !c.tablespaces["ts"] {
			t.Fatal("DROP OWNED removed tablespace ts")
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE ts_owner"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s while owning tablespace", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c, "REASSIGN OWNED BY ts_owner TO bob"); code != "0" {
			t.Fatalf("catalog refused REASSIGN OWNED for tablespace with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE ts_owner"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE ts_owner after REASSIGN with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE bob"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s for reassigned tablespace owner", code, CodeDependentObjects)
		}

		// #105 & #110: ALTER SEQUENCE OWNER TO updates sequence owner, and
		// renaming to a reserved role name is rejected with 42939.
		for _, stmt := range []string{
			"CREATE ROLE seq_role",
			"CREATE SEQUENCE s",
			"ALTER SEQUENCE s OWNER TO seq_role",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "ALTER ROLE seq_role RENAME TO pg_bad"); code != CodeReservedName {
			t.Fatalf("catalog answered %s, want %s", code, CodeReservedName)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE seq_role"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s while owning sequence", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c, "DROP OWNED BY seq_role"); code != "0" {
			t.Fatalf("catalog refused DROP OWNED BY seq_role with %s", code)
		}
		if code := simulationCatalogCode(t, c, "DROP ROLE seq_role"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE seq_role with %s", code)
		}
	})

	// #106, #107, #109, #116: partition tablespace precedence, storage-free
	// composite types, dropped default_tablespace fallback, and temp indexes.
	t.Run("tablespace fallbacks for partitions, composite types and temp indexes", func(t *testing.T) {
		c := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE parent_ts LOCATION '/tmp/parent_ts'",
			"CREATE TABLESPACE child_ts LOCATION '/tmp/child_ts'",
			"CREATE TABLESPACE temp_ts LOCATION '/tmp/temp_ts'",
			"CREATE TABLE p_with_ts (a int) PARTITION BY RANGE (a) TABLESPACE parent_ts",
			"CREATE TABLE p_no_ts (a int) PARTITION BY RANGE (a)",
			"SET default_tablespace = child_ts",
			"CREATE TABLE c_from_parent PARTITION OF p_with_ts FOR VALUES FROM (0) TO (10)",
			"CREATE TABLE c_from_default PARTITION OF p_no_ts FOR VALUES FROM (0) TO (10)",
			"SET default_tablespace = pg_global",
			"CREATE TABLE c_parent_beats_global PARTITION OF p_with_ts FOR VALUES FROM (10) TO (20)",
			"CREATE TYPE ct AS (a int)",
			"SET default_tablespace = temp_ts",
			"DROP TABLESPACE temp_ts",
			"CREATE TABLE after_drop (a int)",
			"CREATE TEMP TABLE tt (a int) TABLESPACE child_ts",
			"CREATE INDEX tt_a_idx ON tt (a)",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if got := c.schemaByName["public"].Relations["c_from_parent"].Tablespace; got != "parent_ts" {
			t.Fatalf("c_from_parent tablespace = %q, want parent_ts", got)
		}
		if got := c.schemaByName["public"].Relations["c_from_default"].Tablespace; got != "child_ts" {
			t.Fatalf("c_from_default tablespace = %q, want child_ts", got)
		}
		if got := c.schemaByName["public"].Relations["c_parent_beats_global"].Tablespace; got != "parent_ts" {
			t.Fatalf("c_parent_beats_global tablespace = %q, want parent_ts", got)
		}
		if got := c.schemaByName["public"].Relations["after_drop"].Tablespace; got != "" {
			t.Fatalf("after_drop tablespace = %q, want empty default", got)
		}
		tt := c.GetRelation("pg_temp", "tt")
		if tt == nil {
			tt = c.GetRelation("", "tt")
		}
		if tt == nil {
			t.Fatal("temp table tt not found")
		}
		for _, idx := range c.indexesByRel[tt.OID] {
			if idx.Tablespace != "" {
				t.Fatalf("temp index %s tablespace = %q, want empty default", idx.Name, idx.Tablespace)
			}
		}
	})

	// #117, #118, #119, #120, #121: table-owned sequence owners, public and
	// implicit-owner schemas in DROP OWNED, owned table row-type dependents,
	// and snapshot-discovered owner roles.
	t.Run("table-owned sequences, schema ownership, row-type deps and snapshot owners", func(t *testing.T) {
		// #117: table-owned sequences reject independent owner changes (0A000)
		// while allowing a no-op owner change to the existing owner.
		c := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE TABLE t (id serial)",
			"ALTER SEQUENCE t_id_seq OWNER TO CURRENT_USER",
		} {
			if code := simulationCatalogCode(t, c, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c, "ALTER SEQUENCE t_id_seq OWNER TO alice"); code != CodeFeatureNotSupported {
			t.Fatalf("catalog answered %s, want %s", code, CodeFeatureNotSupported)
		}
		if code := simulationCatalogCode(t, c, "ALTER TABLE t_id_seq OWNER TO alice"); code != CodeFeatureNotSupported {
			t.Fatalf("catalog answered %s, want %s", code, CodeFeatureNotSupported)
		}

		// #118 & #119: schemas created without AUTHORIZATION record the current
		// role, and DROP OWNED drops owned schemas including public.
		c2 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"SET ROLE alice",
			"CREATE SCHEMA s_implicit",
			"RESET ROLE",
			"ALTER SCHEMA public OWNER TO alice",
		} {
			if code := simulationCatalogCode(t, c2, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c2, "DROP ROLE alice"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s while owning schemas", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c2, "DROP OWNED BY alice"); code != "0" {
			t.Fatalf("catalog refused DROP OWNED BY alice with %s", code)
		}
		if c2.schemaByName["s_implicit"] != nil || c2.schemaByName["public"] != nil {
			t.Fatal("DROP OWNED BY alice kept owned schema")
		}
		if code := simulationCatalogCode(t, c2, "DROP ROLE alice"); code != "0" {
			t.Fatalf("catalog refused DROP ROLE alice after dropping owned schemas with %s", code)
		}

		// #120: external tables using an owned ordinary table's composite row
		// type block DROP OWNED RESTRICT without partial deletion.
		c3 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"SET ROLE alice",
			"CREATE TABLE keep_me (id int)",
			"CREATE TABLE a_row (id int)",
			"SET ROLE bob",
			"CREATE TABLE b_uses_a (x a_row)",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c3, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c3, "DROP OWNED BY alice RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("catalog answered %s, want %s when external table uses row type", code, CodeDependentObjects)
		}
		if c3.schemaByName["public"].Relations["a_row"] == nil || c3.schemaByName["public"].Relations["keep_me"] == nil {
			t.Fatal("failed DROP OWNED deleted owned tables despite external row-type dependency")
		}

		// #121: snapshot-loaded owners register in c.roles so CREATE ROLE sees
		// duplicates and DROP ROLE checks dependencies.
		c4 := New()
		_, err := c4.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
			Schemas: []*metadata.SchemaMetadata{{
				Name: "public",
				Tables: []*metadata.TableMetadata{{
					Name:    "snap_t",
					Owner:   "alice",
					Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "integer"}},
				}},
			}},
		}, LoadMetadataOptions{Full: true})
		if err != nil {
			t.Fatalf("LoadMetadata: %v", err)
		}
		if code := simulationCatalogCode(t, c4, "CREATE ROLE alice"); code != CodeDuplicateObject {
			t.Fatalf("snapshot CREATE ROLE alice answered %s, want %s", code, CodeDuplicateObject)
		}
		if code := simulationCatalogCode(t, c4, "DROP ROLE alice"); code != CodeDependentObjects {
			t.Fatalf("snapshot DROP ROLE alice answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c4, "DROP OWNED BY alice"); code != "0" {
			t.Fatalf("snapshot DROP OWNED BY alice answered %s, want 0", code)
		}
		if code := simulationCatalogCode(t, c4, "DROP ROLE alice"); code != "0" {
			t.Fatalf("snapshot DROP ROLE alice after DROP OWNED answered %s, want 0", code)
		}

		// #122: snapshot database owner blocks DROP ROLE, survives DROP OWNED,
		// and follows RENAME and REASSIGN OWNED.
		c5 := New()
		_, err = c5.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
			Name:    "appdb",
			Owner:   "db_alice",
			Schemas: []*metadata.SchemaMetadata{{Name: "public"}},
		}, LoadMetadataOptions{Full: true})
		if err != nil {
			t.Fatalf("LoadMetadata db owner: %v", err)
		}
		if code := simulationCatalogCode(t, c5, "CREATE ROLE bob"); code != "0" {
			t.Fatalf("CREATE ROLE bob: %s", code)
		}
		if code := simulationCatalogCode(t, c5, "DROP OWNED BY db_alice"); code != "0" {
			t.Fatalf("DROP OWNED BY db_alice: %s", code)
		}
		if code := simulationCatalogCode(t, c5, "DROP ROLE db_alice"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE db_alice answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c5, "ALTER ROLE db_alice RENAME TO db_carol"); code != "0" {
			t.Fatalf("RENAME db_alice: %s", code)
		}
		if code := simulationCatalogCode(t, c5, "DROP ROLE db_carol"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE db_carol answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c5, "REASSIGN OWNED BY db_carol TO bob"); code != "0" {
			t.Fatalf("REASSIGN db_carol: %s", code)
		}
		if code := simulationCatalogCode(t, c5, "DROP ROLE db_carol"); code != "0" {
			t.Fatalf("DROP ROLE db_carol after REASSIGN answered %s, want 0", code)
		}
		if code := simulationCatalogCode(t, c5, "DROP ROLE bob"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE bob answered %s, want %s", code, CodeDependentObjects)
		}

		// #125: non-index CHECK constraints on storage-free foreign tables do
		// not validate default_tablespace.
		c6 := New()
		for _, stmt := range []string{
			"SET default_tablespace = pg_global",
			"CREATE FOREIGN TABLE ft_chk (a int CHECK (a > 0)) SERVER srv",
		} {
			if code := simulationCatalogCode(t, c6, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}

		// #126: DROP TABLE ... CASCADE follows multi-hop row-type dependents.
		c7 := New()
		for _, stmt := range []string{
			"CREATE TABLE hop_a (x int)",
			"CREATE TABLE hop_b (x hop_a)",
			"CREATE TABLE hop_c (x hop_b)",
			"DROP TABLE hop_a CASCADE",
		} {
			if code := simulationCatalogCode(t, c7, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if c7.schemaByName["public"].Relations["hop_b"] != nil || c7.schemaByName["public"].Relations["hop_c"] != nil {
			t.Fatal("DROP TABLE hop_a CASCADE left dependent row-type table behind")
		}

		// #127: CREATE SCHEMA ... AUTHORIZATION registers the named role in
		// c.roles so CREATE ROLE, REASSIGN OWNED, and DROP OWNED recognize it.
		c8 := New()
		if code := simulationCatalogCode(t, c8, "CREATE SCHEMA auth_s AUTHORIZATION auth_role"); code != "0" {
			t.Fatalf("CREATE SCHEMA AUTHORIZATION: %s", code)
		}
		if code := simulationCatalogCode(t, c8, "CREATE ROLE auth_role"); code != CodeDuplicateObject {
			t.Fatalf("CREATE ROLE auth_role answered %s, want %s", code, CodeDuplicateObject)
		}
		if code := simulationCatalogCode(t, c8, "DROP OWNED BY auth_role"); code != "0" {
			t.Fatalf("DROP OWNED BY auth_role: %s", code)
		}
		if code := simulationCatalogCode(t, c8, "DROP ROLE auth_role"); code != "0" {
			t.Fatalf("DROP ROLE auth_role: %s", code)
		}

		// #128: LoadMetadata after SetSessionUser does not invent session-user
		// ownership for snapshot objects whose Owner field is empty.
		c9 := New()
		c9.SetSessionUser("alice")
		_, err = c9.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
			Schemas: []*metadata.SchemaMetadata{{
				Name: "public",
				Tables: []*metadata.TableMetadata{{
					Name:    "unowned_snap_t",
					Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "integer"}},
				}},
			}},
		}, LoadMetadataOptions{Full: true})
		if err != nil {
			t.Fatalf("LoadMetadata unowned: %v", err)
		}
		if code := simulationCatalogCode(t, c9, "DROP OWNED BY alice"); code != "0" {
			t.Fatalf("DROP OWNED BY alice on unowned snapshot: %s", code)
		}
		if c9.schemaByName["public"].Relations["unowned_snap_t"] == nil {
			t.Fatal("DROP OWNED BY alice deleted snapshot table with empty Owner")
		}

		// #129: extensions track ownership and drop with their members as one
		// unit during DROP OWNED.
		c10 := New()
		for _, stmt := range []string{
			"CREATE ROLE ext_role",
			"SET ROLE ext_role",
			"CREATE EXTENSION citext",
			"RESET ROLE",
			"CREATE TABLE ext_dep_t (c citext)",
		} {
			if code := simulationCatalogCode(t, c10, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c10, "DROP OWNED BY ext_role RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED on in-use extension answered %s, want %s", code, CodeDependentObjects)
		}
		if c10.extByName["citext"] == nil {
			t.Fatal("failed DROP OWNED unregistered citext")
		}
		if code := simulationCatalogCode(t, c10, "DROP TABLE ext_dep_t"); code != "0" {
			t.Fatalf("DROP TABLE ext_dep_t: %s", code)
		}
		if code := simulationCatalogCode(t, c10, "DROP OWNED BY ext_role RESTRICT"); code != "0" {
			t.Fatalf("DROP OWNED BY ext_role: %s", code)
		}
		if c10.extByName["citext"] != nil {
			t.Fatal("DROP OWNED left citext registered in extByName")
		}
		if code := simulationCatalogCode(t, c10, "DROP ROLE ext_role"); code != "0" {
			t.Fatalf("DROP ROLE ext_role: %s", code)
		}

		// #132: an external index on an extension operator class is a RESTRICT
		// dependent, and CASCADE drops it.
		c13 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION pg_trgm",
			"CREATE TABLE trgm_t (t text)",
			"CREATE INDEX trgm_idx ON trgm_t USING gin (t gin_trgm_ops)",
		} {
			if code := simulationCatalogCode(t, c13, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c13, "DROP EXTENSION pg_trgm"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION on an index-using opclass answered %s, want %s", code, CodeDependentObjects)
		}
		if c13.extByName["pg_trgm"] == nil || c13.schemaByName["public"].Indexes["trgm_idx"] == nil {
			t.Fatal("failed DROP EXTENSION pg_trgm changed the catalog")
		}
		if code := simulationCatalogCode(t, c13, "DROP EXTENSION pg_trgm CASCADE"); code != "0" {
			t.Fatalf("DROP EXTENSION pg_trgm CASCADE: %s", code)
		}
		if c13.extByName["pg_trgm"] != nil || c13.schemaByName["public"].Indexes["trgm_idx"] != nil {
			t.Fatal("DROP EXTENSION pg_trgm CASCADE left the extension or its opclass index")
		}

		// #133: extension-owned operators leave the catalog with the extension.
		c14 := New()
		if code := simulationCatalogCode(t, c14, "CREATE EXTENSION pg_trgm"); code != "0" {
			t.Fatalf("CREATE EXTENSION pg_trgm: %s", code)
		}
		pctBefore := len(c14.operByKey[operKey{name: "%", left: TEXTOID, right: TEXTOID}])
		if pctBefore == 0 {
			t.Fatal("pg_trgm did not register its % operator")
		}
		if code := simulationCatalogCode(t, c14, "DROP EXTENSION pg_trgm"); code != "0" {
			t.Fatalf("DROP EXTENSION pg_trgm: %s", code)
		}
		if ops := c14.operByKey[operKey{name: "%", left: TEXTOID, right: TEXTOID}]; len(ops) != 0 {
			t.Fatalf("DROP EXTENSION left %d %% operators in operByKey", len(ops))
		}
		if code := simulationCatalogCode(t, c14, "CREATE EXTENSION pg_trgm"); code != "0" {
			t.Fatalf("recreating pg_trgm: %s", code)
		}
		if ops := c14.operByKey[operKey{name: "%", left: TEXTOID, right: TEXTOID}]; len(ops) != pctBefore {
			t.Fatalf("recreated pg_trgm registered %d %% operators, want %d", len(ops), pctBefore)
		}

		// #134: DROP OWNED of an extension owner preflights indexes that use
		// the extension's operator class, and CASCADE drops them.
		c15 := New()
		for _, stmt := range []string{
			"CREATE ROLE ext_role",
			"SET ROLE ext_role",
			"CREATE EXTENSION pg_trgm",
			"RESET ROLE",
			"CREATE TABLE owned_trgm_t (t text)",
			"CREATE INDEX owned_trgm_idx ON owned_trgm_t USING gin (t gin_trgm_ops)",
		} {
			if code := simulationCatalogCode(t, c15, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c15, "DROP OWNED BY ext_role RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED on an extension with an opclass index answered %s, want %s", code, CodeDependentObjects)
		}
		if c15.extByName["pg_trgm"] == nil || c15.schemaByName["public"].Indexes["owned_trgm_idx"] == nil {
			t.Fatal("failed DROP OWNED changed the extension or its index")
		}
		if code := simulationCatalogCode(t, c15, "DROP OWNED BY ext_role CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED BY ext_role CASCADE: %s", code)
		}
		if c15.extByName["pg_trgm"] != nil || c15.schemaByName["public"].Indexes["owned_trgm_idx"] != nil {
			t.Fatal("DROP OWNED CASCADE left the extension or its opclass index")
		}

		// #143 & #139: an index that names no opclass depends on the default one
		// its access method resolves for the column type, and an unqualified
		// name resolves through the search path.
		c16 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION btree_gin",
			"CREATE TABLE opc_t (t text)",
			"CREATE INDEX opc_idx ON opc_t USING gin (t)",
		} {
			if code := simulationCatalogCode(t, c16, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c16, "DROP EXTENSION btree_gin"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION btree_gin over a default opclass answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c16, "DROP EXTENSION btree_gin CASCADE"); code != "0" {
			t.Fatalf("DROP EXTENSION btree_gin CASCADE: %s", code)
		}
		if c16.schemaByName["public"].Indexes["opc_idx"] != nil {
			t.Fatal("CASCADE left the index that used the extension's default opclass")
		}

		// #136: ALTER DATABASE naming a database a snapshot does not hold is
		// 3D000, not a silent reuse of the loaded database's owner.
		c17 := New()
		_, err = c17.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
			Owner:   "db_alice",
			Schemas: []*metadata.SchemaMetadata{{Name: "public"}},
		}, LoadMetadataOptions{Full: true})
		if err != nil {
			t.Fatalf("LoadMetadata: %v", err)
		}
		if code := simulationCatalogCode(t, c17, "CREATE ROLE bob"); code != "0" {
			t.Fatalf("CREATE ROLE bob: %s", code)
		}
		if code := simulationCatalogCode(t, c17, "ALTER DATABASE missing OWNER TO bob"); code != CodeInvalidCatalogName {
			t.Fatalf("ALTER DATABASE missing answered %s, want %s", code, CodeInvalidCatalogName)
		}
		if c17.databaseOwners["\x00current"] != "db_alice" {
			t.Fatalf("failed ALTER DATABASE changed the loaded owner to %q", c17.databaseOwners["\x00current"])
		}

		// #137: dropping an extension from a clone leaves the source intact.
		c18 := New()
		if code := simulationCatalogCode(t, c18, "CREATE EXTENSION pg_trgm"); code != "0" {
			t.Fatalf("CREATE EXTENSION pg_trgm: %s", code)
		}
		clone := c18.Clone()
		if code := simulationCatalogCode(t, clone, "DROP EXTENSION pg_trgm"); code != "0" {
			t.Fatalf("clone DROP EXTENSION pg_trgm: %s", code)
		}
		if ops := c18.operByKey[operKey{name: "%", left: TEXTOID, right: TEXTOID}]; len(ops) == 0 {
			t.Fatal("clone DROP EXTENSION removed the source catalog's operators")
		}

		// #138: event-trigger ownership participates in the role lifecycle.
		c19 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"CREATE FUNCTION evt_fn() RETURNS event_trigger LANGUAGE plpgsql AS 'BEGIN END'",
			"SET ROLE alice",
			"CREATE EVENT TRIGGER evt_owned ON ddl_command_start EXECUTE FUNCTION evt_fn()",
			"RESET ROLE",
			"ALTER FUNCTION evt_fn() OWNER TO bob",
		} {
			if code := simulationCatalogCode(t, c19, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c19, "DROP ROLE alice"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE alice with an event trigger answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c19, "DROP OWNED BY alice"); code != "0" {
			t.Fatalf("DROP OWNED BY alice: %s", code)
		}
		if c19.eventTriggerByName["evt_owned"] != nil {
			t.Fatal("DROP OWNED left the owned event trigger")
		}
		if code := simulationCatalogCode(t, c19, "DROP ROLE alice"); code != "0" {
			t.Fatalf("DROP ROLE alice after DROP OWNED answered %s, want 0", code)
		}

		// #140: a user-defined operator records its owner and participates in
		// the role lifecycle.
		c20 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE FUNCTION op_cmp(int, int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"SET ROLE alice",
			"CREATE OPERATOR === (LEFTARG = int, RIGHTARG = int, PROCEDURE = op_cmp)",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c20, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c20, "DROP ROLE alice"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE alice with an operator answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c20, "DROP OWNED BY alice"); code != "0" {
			t.Fatalf("DROP OWNED BY alice: %s", code)
		}
		if ops := c20.operByKey[operKey{name: "===", left: INT4OID, right: INT4OID}]; len(ops) != 0 {
			t.Fatalf("DROP OWNED left %d owned operators", len(ops))
		}
		if code := simulationCatalogCode(t, c20, "DROP ROLE alice"); code != "0" {
			t.Fatalf("DROP ROLE alice after DROP OWNED answered %s, want 0", code)
		}

		// #141: DROP SCHEMA ... CASCADE routes a contained extension through the
		// extension's CASCADE path, so an external index on its opclass goes too.
		c21 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA ext_s",
			"CREATE EXTENSION pg_trgm SCHEMA ext_s",
			"CREATE TABLE schema_cascade_t (t text)",
			"CREATE INDEX schema_cascade_idx ON schema_cascade_t USING gist (t ext_s.gist_trgm_ops)",
			"DROP SCHEMA ext_s CASCADE",
		} {
			if code := simulationCatalogCode(t, c21, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if c21.extByName["pg_trgm"] != nil {
			t.Fatal("DROP SCHEMA CASCADE left the extension registered")
		}
		if c21.schemaByName["public"].Indexes["schema_cascade_idx"] != nil {
			t.Fatal("DROP SCHEMA CASCADE left the index that used the extension's opclass")
		}

		// #142: casts an extension created leave castIndex with the extension.
		c22 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION citext",
			"CREATE TABLE cast_t (c citext)",
		} {
			if code := simulationCatalogCode(t, c22, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		castsBefore := len(c22.castIndex)
		if castsBefore == 0 {
			t.Fatal("citext registered no casts")
		}
		if code := simulationCatalogCode(t, c22, "DROP TABLE cast_t"); code != "0" {
			t.Fatalf("DROP TABLE cast_t: %s", code)
		}
		if code := simulationCatalogCode(t, c22, "DROP EXTENSION citext"); code != "0" {
			t.Fatalf("DROP EXTENSION citext: %s", code)
		}
		if len(c22.castIndex) != len(BuiltinCasts) {
			t.Fatalf("DROP EXTENSION left %d casts, want %d built-in", len(c22.castIndex), len(BuiltinCasts))
		}
		if code := simulationCatalogCode(t, c22, "CREATE EXTENSION citext"); code != "0" {
			t.Fatalf("recreating citext: %s", code)
		}

		// #144: ALTER OPERATOR / ALTER EVENT TRIGGER OWNER TO move tracked
		// ownership, so the old owner can be dropped and the new one cannot.
		c23 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"CREATE FUNCTION own_cmp(int, int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"CREATE FUNCTION own_evt() RETURNS event_trigger LANGUAGE plpgsql AS 'BEGIN END'",
			"SET ROLE alice",
			"CREATE OPERATOR ==>> (LEFTARG = int, RIGHTARG = int, PROCEDURE = own_cmp)",
			"CREATE EVENT TRIGGER own_evt_trig ON ddl_command_start EXECUTE FUNCTION own_evt()",
			"RESET ROLE",
			"ALTER OPERATOR ==>> (int, int) OWNER TO bob",
			"ALTER EVENT TRIGGER own_evt_trig OWNER TO bob",
		} {
			if code := simulationCatalogCode(t, c23, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c23, "DROP ROLE alice"); code != "0" {
			t.Fatalf("DROP ROLE alice after ownership transfer answered %s, want 0", code)
		}
		if code := simulationCatalogCode(t, c23, "DROP ROLE bob"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE bob after ownership transfer answered %s, want %s", code, CodeDependentObjects)
		}

		// #145: an unqualified type or routine owner target resolves through the
		// search path, not to the first schema.
		c24 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"CREATE SCHEMA s1",
			"CREATE SCHEMA s2",
			"SET search_path = s1, s2",
			"CREATE TYPE s2.mood2 AS ENUM ('ok')",
			"CREATE FUNCTION s2.sp_fn() RETURNS int LANGUAGE sql AS 'SELECT 1'",
			"SET ROLE alice",
		} {
			if code := simulationCatalogCode(t, c24, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		// Give alice ownership through ALTER ... OWNER TO in the later schema.
		if code := simulationCatalogCode(t, c24, "ALTER TYPE mood2 OWNER TO alice"); code != "0" {
			t.Fatalf("ALTER TYPE through the search path answered %s", code)
		}
		if code := simulationCatalogCode(t, c24, "ALTER FUNCTION sp_fn() OWNER TO alice"); code != "0" {
			t.Fatalf("ALTER FUNCTION through the search path answered %s", code)
		}

		// #146: SET ROLE refuses a role the catalog does not hold.
		c25 := New()
		if code := simulationCatalogCode(t, c25, "SET ROLE ghost"); code != CodeUndefinedObject {
			t.Fatalf("SET ROLE ghost answered %s, want %s", code, CodeUndefinedObject)
		}
		if c25.currentUser != "" {
			t.Fatalf("failed SET ROLE changed the current user to %q", c25.currentUser)
		}

		// #147: DROP OWNED removes an owned event trigger together with the
		// function it calls under RESTRICT.
		c26 := New()
		for _, stmt := range []string{
			"CREATE ROLE evt_role",
			"SET ROLE evt_role",
			"CREATE FUNCTION evt_role_fn() RETURNS event_trigger LANGUAGE plpgsql AS 'BEGIN END'",
			"CREATE EVENT TRIGGER evt_role_trig ON ddl_command_start EXECUTE FUNCTION evt_role_fn()",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c26, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c26, "DROP OWNED BY evt_role RESTRICT"); code != "0" {
			t.Fatalf("DROP OWNED on an owned event trigger answered %s, want 0", code)
		}
		if c26.eventTriggerByName["evt_role_trig"] != nil || len(c26.userProcs) != 0 {
			t.Fatal("DROP OWNED left the owned event trigger or its function")
		}

		// #148: a grant's grantor is a role dependency and DROP OWNED revokes it.
		c27 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"CREATE ROLE carol",
			"CREATE TABLE acl_t (id int)",
			"SET ROLE alice",
			"GRANT SELECT ON acl_t TO carol",
			"RESET ROLE",
			"ALTER TABLE acl_t OWNER TO bob",
		} {
			if code := simulationCatalogCode(t, c27, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c27, "DROP ROLE alice"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE alice as grantor answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c27, "DROP OWNED BY alice"); code != "0" {
			t.Fatalf("DROP OWNED BY alice: %s", code)
		}
		if len(c27.grants) != 0 {
			t.Fatalf("DROP OWNED left %d grants issued by alice", len(c27.grants))
		}
		if code := simulationCatalogCode(t, c27, "DROP ROLE alice"); code != "0" {
			t.Fatalf("DROP ROLE alice after DROP OWNED answered %s, want 0", code)
		}

		// #149: an exclusion constraint's backing index depends on the operator
		// class its elements name.
		c28 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION pg_trgm",
			"CREATE TABLE excl_t (txt text)",
			"ALTER TABLE excl_t ADD CONSTRAINT excl_trgm EXCLUDE USING gist (txt gist_trgm_ops WITH %)",
		} {
			if code := simulationCatalogCode(t, c28, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c28, "DROP EXTENSION pg_trgm"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION over an exclusion index answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c28, "DROP EXTENSION pg_trgm CASCADE"); code != "0" {
			t.Fatalf("DROP EXTENSION pg_trgm CASCADE over an exclusion index answered %s", code)
		}
		if len(c28.schemaByName["public"].Indexes) != 0 {
			t.Fatal("CASCADE left the exclusion constraint's backing index")
		}

		// #152: a failed ALTER INDEX subcommand rolls back an earlier SET
		// TABLESPACE in the same statement.
		c31 := New()
		for _, stmt := range []string{
			"CREATE TABLESPACE alt_ts LOCATION '/tmp/alt_ts'",
			"CREATE TABLE alt_p (a int) PARTITION BY RANGE (a)",
			"CREATE INDEX alt_p_idx ON alt_p (a)",
		} {
			if code := simulationCatalogCode(t, c31, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		// The parser accepts ATTACH PARTITION only as a whole statement, so the
		// two-subcommand form is built here to exercise the catalog's rollback.
		stmt := &nodes.AlterTableStmt{
			Relation: &nodes.RangeVar{Relname: "alt_p_idx"},
			ObjType:  int(nodes.OBJECT_INDEX),
			Cmds: &nodes.List{Items: []nodes.Node{
				&nodes.AlterTableCmd{Subtype: int(nodes.AT_SetTableSpace), Name: "alt_ts"},
				&nodes.AlterTableCmd{Subtype: int(nodes.AT_AttachPartition), Def: &nodes.PartitionCmd{
					Name: &nodes.RangeVar{Relname: "missing_idx"},
				}},
			}},
		}
		if err := c31.ProcessUtility(stmt); err == nil {
			t.Fatal("ALTER INDEX with a missing partition index was accepted")
		}
		if idx := c31.schemaByName["public"].Indexes["alt_p_idx"]; idx == nil || idx.Tablespace != "" {
			t.Fatal("a failed ALTER INDEX kept the new tablespace")
		}

		// #154 & #158: GRANTED BY records the named grantor, and REVOKE removes
		// only that grantor's grant.
		c33 := New()
		for _, stmt := range []string{
			"CREATE ROLE grantor_a",
			"CREATE ROLE grantor_b",
			"CREATE ROLE grantee_c",
			"CREATE TABLE g_t (id int)",
		} {
			if code := simulationCatalogCode(t, c33, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c33, "GRANT SELECT ON g_t TO grantee_c GRANTED BY grantor_a"); code != "0" {
			t.Fatalf("GRANT ... GRANTED BY answered %s", code)
		}
		if code := simulationCatalogCode(t, c33, "DROP ROLE grantor_a"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE grantor_a answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c33, "REVOKE SELECT ON g_t FROM grantee_c GRANTED BY grantor_b"); code != "0" {
			t.Fatalf("REVOKE by another grantor answered %s", code)
		}
		if len(c33.grants) != 1 {
			t.Fatalf("REVOKE by another grantor left %d grants, want 1", len(c33.grants))
		}
		if code := simulationCatalogCode(t, c33, "REVOKE SELECT ON g_t FROM grantee_c GRANTED BY grantor_a"); code != "0" {
			t.Fatalf("REVOKE by the grantor answered %s", code)
		}
		if len(c33.grants) != 0 {
			t.Fatal("REVOKE by the grantor left its grant")
		}

		// #157: a loaded catalog still checks sequence/table owner agreement
		// when both owners are known.
		c34 := New()
		_, err = c34.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{
			Schemas: []*metadata.SchemaMetadata{{
				Name: "public",
				Tables: []*metadata.TableMetadata{{
					Name:    "loaded_t",
					Owner:   "loaded_alice",
					Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "integer"}},
				}},
			}},
		}, LoadMetadataOptions{Full: true})
		if err != nil {
			t.Fatalf("LoadMetadata: %v", err)
		}
		for _, stmt := range []string{"CREATE ROLE loaded_bob", "CREATE SEQUENCE loaded_seq"} {
			if code := simulationCatalogCode(t, c34, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c34, "ALTER SEQUENCE loaded_seq OWNER TO loaded_bob"); code != "0" {
			t.Fatalf("ALTER SEQUENCE OWNER TO loaded_bob answered %s", code)
		}
		if code := simulationCatalogCode(t, c34, "ALTER SEQUENCE loaded_seq OWNED BY loaded_t.id"); code != CodeObjectNotInPrerequisiteState {
			t.Fatalf("OWNED BY across owners answered %s, want %s", code, CodeObjectNotInPrerequisiteState)
		}

		// #160: an external operator implemented by an extension function keeps
		// the extension in use.
		c35 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION pg_trgm",
			"CREATE OPERATOR ==== (LEFTARG = text, RIGHTARG = text, PROCEDURE = similarity_op)",
		} {
			if code := simulationCatalogCode(t, c35, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c35, "DROP EXTENSION pg_trgm"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION over an operator dependent answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c35, "DROP EXTENSION pg_trgm CASCADE"); code != "0" {
			t.Fatalf("DROP EXTENSION pg_trgm CASCADE answered %s", code)
		}

		// #163: CREATE OPERATOR in a schema that does not exist is refused
		// before the operator is registered.
		c36 := New()
		if code := simulationCatalogCode(t, c36, "CREATE OPERATOR missing.@@ (LEFTARG = int, RIGHTARG = int, PROCEDURE = int4eq)"); code != CodeUndefinedSchema {
			t.Fatalf("CREATE OPERATOR in a missing schema answered %s, want %s", code, CodeUndefinedSchema)
		}
		if len(c36.operByOID) != len(BuiltinOperators) {
			t.Fatalf("failed CREATE OPERATOR registered %d user operators", len(c36.operByOID)-len(BuiltinOperators))
		}

		// #164 & #166: DROP OWNED removes an operator's dependency edges, and
		// DROP FUNCTION ... CASCADE removes a dependent operator.
		c37 := New()
		for _, stmt := range []string{
			"CREATE ROLE op_role",
			"CREATE FUNCTION op_impl(int, int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"SET ROLE op_role",
			"CREATE OPERATOR <<>> (LEFTARG = int, RIGHTARG = int, PROCEDURE = op_impl)",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c37, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c37, "DROP OWNED BY op_role"); code != "0" {
			t.Fatalf("DROP OWNED BY op_role: %s", code)
		}
		if len(c37.operByKey[operKey{name: "<<>>", left: INT4OID, right: INT4OID}]) != 0 {
			t.Fatal("DROP OWNED left the owned operator")
		}
		for _, d := range c37.deps {
			if d.ObjType == 'O' || d.RefType == 'O' {
				t.Fatalf("DROP OWNED left an operator dependency row: %+v", d)
			}
		}
		if code := simulationCatalogCode(t, c37, "CREATE OPERATOR <<>> (LEFTARG = int, RIGHTARG = int, PROCEDURE = op_impl)"); code != "0" {
			t.Fatalf("recreating the operator answered %s", code)
		}
		if code := simulationCatalogCode(t, c37, "DROP FUNCTION op_impl(int, int) CASCADE"); code != "0" {
			t.Fatalf("DROP FUNCTION ... CASCADE answered %s", code)
		}
		if len(c37.operByKey[operKey{name: "<<>>", left: INT4OID, right: INT4OID}]) != 0 {
			t.Fatal("DROP FUNCTION ... CASCADE left the dependent operator")
		}

		// #165: an unqualified ALTER OPERATOR picks the first search-path schema.
		c38 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"CREATE SCHEMA s1",
			"CREATE SCHEMA s2",
			"CREATE FUNCTION s1.cmp_a(int, int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"CREATE FUNCTION s2.cmp_b(int, int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"CREATE OPERATOR s2.%% (LEFTARG = int, RIGHTARG = int, PROCEDURE = s2.cmp_b)",
			"CREATE OPERATOR s1.%% (LEFTARG = int, RIGHTARG = int, PROCEDURE = s1.cmp_a)",
			"SET search_path = s1, s2",
			"ALTER OPERATOR %% (int, int) OWNER TO bob",
		} {
			if code := simulationCatalogCode(t, c38, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c38, "DROP ROLE bob"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE bob answered %s, want %s", code, CodeDependentObjects)
		}

		// #167: DROP OWNED over an extension that defines an access method keeps
		// external users of it (an index built ON the method) and takes them
		// with a CASCADE.
		RegisterExtensionSQL("test_am_ext", "CREATE ACCESS METHOD test_am TYPE INDEX HANDLER bthandler;")
		defer delete(extensionScripts, "test_am_ext")
		c39 := New()
		for _, stmt := range []string{
			"CREATE ROLE am_role",
			"SET ROLE am_role",
			"CREATE EXTENSION test_am_ext",
			"RESET ROLE",
			"CREATE TABLE am_t (a int)",
			"CREATE INDEX am_idx ON am_t USING test_am (a)",
		} {
			if code := simulationCatalogCode(t, c39, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c39, "DROP OWNED BY am_role RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED over an AM-dependent index answered %s, want %s", code, CodeDependentObjects)
		}
		if c39.schemaByName["public"].Indexes["am_idx"] == nil {
			t.Fatal("failed DROP OWNED removed the AM-dependent index")
		}
		if code := simulationCatalogCode(t, c39, "DROP OWNED BY am_role CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED BY am_role CASCADE: %s", code)
		}
		if c39.schemaByName["public"].Indexes["am_idx"] != nil || c39.accessMethodByName["test_am"] != nil {
			t.Fatal("DROP OWNED CASCADE left the index or the extension access method")
		}

		// #170: ALTER EXTENSION ... SET SCHEMA moves the extension's schema.
		c40 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA ext_move",
			"CREATE EXTENSION citext",
			"ALTER EXTENSION citext SET SCHEMA ext_move",
		} {
			if code := simulationCatalogCode(t, c40, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if ext := c40.extByName["citext"]; ext == nil || c40.schemas[ext.SchemaOID] == nil || c40.schemas[ext.SchemaOID].Name != "ext_move" {
			t.Fatal("ALTER EXTENSION SET SCHEMA did not move the extension")
		}
		if code := simulationCatalogCode(t, c40, "DROP SCHEMA ext_move CASCADE"); code != "0" {
			t.Fatalf("DROP SCHEMA ext_move CASCADE: %s", code)
		}
		if c40.extByName["citext"] != nil {
			t.Fatal("DROP SCHEMA CASCADE left the relocated extension")
		}

		// #171/#172/#173: ALTER EXTENSION SET SCHEMA moves only that
		// extension's members, keeps their namespaces in step, and refuses a
		// target that already holds a member's name.
		RegisterExtensionSQL("test_move_ext", "CREATE TABLE move_member (a int);")
		defer delete(extensionScripts, "test_move_ext")
		c41 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA move_a",
			"CREATE SCHEMA move_b",
			"CREATE EXTENSION test_move_ext SCHEMA move_a",
			"CREATE SCHEMA blocker",
			"CREATE TABLE blocker.move_member (a int)",
		} {
			if code := simulationCatalogCode(t, c41, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c41, "ALTER EXTENSION test_move_ext SET SCHEMA blocker"); code != CodeDuplicateTable {
			t.Fatalf("SET SCHEMA onto a taken name answered %s, want %s", code, CodeDuplicateTable)
		}
		if c41.schemaByName["move_a"].Relations["move_member"] == nil {
			t.Fatal("a refused SET SCHEMA moved the member out of its schema")
		}
		if code := simulationCatalogCode(t, c41, "ALTER EXTENSION test_move_ext SET SCHEMA move_b"); code != "0" {
			t.Fatalf("SET SCHEMA into a free schema answered %s", code)
		}
		if c41.schemaByName["move_a"].Relations["move_member"] != nil || c41.schemaByName["move_b"].Relations["move_member"] == nil {
			t.Fatal("SET SCHEMA did not move the member relation")
		}

		// #174: DROP EXTENSION removes an event trigger its script created.
		RegisterExtensionSQL("test_evt_ext", "CREATE FUNCTION evt_ext_fn() RETURNS event_trigger LANGUAGE plpgsql AS 'BEGIN END'; CREATE EVENT TRIGGER evt_ext_trig ON ddl_command_start EXECUTE FUNCTION evt_ext_fn();")
		defer delete(extensionScripts, "test_evt_ext")
		c42 := New()
		for _, stmt := range []string{"CREATE EXTENSION test_evt_ext", "DROP EXTENSION test_evt_ext"} {
			if code := simulationCatalogCode(t, c42, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if len(c42.eventTriggers) != 0 {
			t.Fatalf("DROP EXTENSION left %d event triggers", len(c42.eventTriggers))
		}

		// #175 & #176: CREATE TABLE AS and LIKE-cloned indexes record their
		// access-method dependencies.
		RegisterExtensionSQL("test_am_ext2", "CREATE ACCESS METHOD test_am2 TYPE INDEX HANDLER bthandler;")
		defer delete(extensionScripts, "test_am_ext2")
		c43 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION test_am_ext2",
			"CREATE TABLE am_src (a int)",
			"CREATE INDEX am_src_idx ON am_src USING test_am2 (a)",
			"CREATE TABLE am_dst (LIKE am_src INCLUDING INDEXES)",
			"DROP TABLE am_src",
		} {
			if code := simulationCatalogCode(t, c43, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c43, "DROP EXTENSION test_am_ext2"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION with a LIKE-cloned AM index answered %s, want %s", code, CodeDependentObjects)
		}

		// #177: a failed GRANT does not leave its grantee recorded.
		c44 := New()
		if code := simulationCatalogCode(t, c44, "GRANT SELECT ON missing_table TO ghost"); code == "0" {
			t.Fatal("GRANT on a missing table was accepted")
		}
		if c44.roles["ghost"] {
			t.Fatal("a failed GRANT recorded its grantee")
		}

		// #178: a schema holding only a user operator is not empty.
		c45 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA op_s",
			"CREATE FUNCTION op_s.cmp(int, int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"CREATE OPERATOR op_s.=== (LEFTARG = int, RIGHTARG = int, PROCEDURE = op_s.cmp)",
		} {
			if code := simulationCatalogCode(t, c45, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c45, "DROP SCHEMA op_s"); code != CodeSchemaNotEmpty {
			t.Fatalf("DROP SCHEMA over an operator answered %s, want %s", code, CodeSchemaNotEmpty)
		}
		if code := simulationCatalogCode(t, c45, "DROP SCHEMA op_s CASCADE"); code != "0" {
			t.Fatalf("DROP SCHEMA op_s CASCADE answered %s", code)
		}
		if c45.operatorSchemas != nil {
			for oid, ns := range c45.operatorSchemas {
				if ns == c45.schemaByName["op_s"].OID {
					t.Fatalf("CASCADE left operator %d in the dropped schema", oid)
				}
			}
		}

		// #179: a routine signature collision blocks ALTER EXTENSION SET SCHEMA.
		RegisterExtensionSQL("test_routine_ext", "CREATE FUNCTION shared_fn(int) RETURNS int LANGUAGE sql AS 'SELECT 1';")
		defer delete(extensionScripts, "test_routine_ext")
		c46 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA rsrc",
			"CREATE SCHEMA rdst",
			"CREATE EXTENSION test_routine_ext SCHEMA rsrc",
			"CREATE FUNCTION rdst.shared_fn(int) RETURNS int LANGUAGE sql AS 'SELECT 2'",
		} {
			if code := simulationCatalogCode(t, c46, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c46, "ALTER EXTENSION test_routine_ext SET SCHEMA rdst"); code != CodeDuplicateFunction {
			t.Fatalf("SET SCHEMA onto a taken signature answered %s, want %s", code, CodeDuplicateFunction)
		}
		if ext := c46.extByName["test_routine_ext"]; ext == nil || c46.schemas[ext.SchemaOID].Name != "rsrc" {
			t.Fatal("a refused SET SCHEMA moved the extension")
		}

		// #180: DROP OWNED of a schema owner takes operators inside it, and
		// RESTRICT reports one owned by another role.
		c47 := New()
		for _, stmt := range []string{
			"CREATE ROLE schema_owner",
			"CREATE ROLE op_owner",
			"CREATE SCHEMA own_s",
			"ALTER SCHEMA own_s OWNER TO schema_owner",
			"CREATE FUNCTION own_s.f(int, int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"SET ROLE op_owner",
			"CREATE OPERATOR own_s.!! (LEFTARG = int, RIGHTARG = int, PROCEDURE = own_s.f)",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c47, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c47, "DROP OWNED BY schema_owner RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED over a foreign-owned operator answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c47, "DROP OWNED BY schema_owner CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED BY schema_owner CASCADE: %s", code)
		}
		for oid := range c47.operByOID {
			if ns, ok := c47.operatorSchemas[oid]; ok && c47.schemas[ns] == nil {
				t.Fatalf("DROP OWNED CASCADE left operator %d in a dropped schema", oid)
			}
		}

		// #181: an access method whose handler comes from an extension keeps it
		// in use.
		RegisterExtensionSQL("test_handler_ext", "CREATE FUNCTION handler_fn(internal) RETURNS table_am_handler LANGUAGE internal AS 'bthandler';")
		defer delete(extensionScripts, "test_handler_ext")
		c48 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION test_handler_ext",
			"CREATE ACCESS METHOD ext_am TYPE TABLE HANDLER handler_fn",
		} {
			if code := simulationCatalogCode(t, c48, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c48, "DROP EXTENSION test_handler_ext"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION with a dependent AM answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c48, "DROP EXTENSION test_handler_ext CASCADE"); code != "0" {
			t.Fatalf("DROP EXTENSION CASCADE with a dependent AM answered %s", code)
		}
		if c48.accessMethodByName["ext_am"] != nil {
			t.Fatal("CASCADE left the AM that used the extension handler")
		}

		// #182: an operator class on an extension access method keeps it in use.
		c49 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION btree_gist",
			"CREATE OPERATOR FAMILY fam1 USING gist",
		} {
			if code := simulationCatalogCode(t, c49, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}

		// #183: default_table_access_method reaches CREATE TABLE and CTAS.
		RegisterExtensionSQL("test_tam_ext", "CREATE ACCESS METHOD test_tam TYPE TABLE HANDLER heap_tableam_handler;")
		defer delete(extensionScripts, "test_tam_ext")
		c50 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION test_tam_ext",
			"SET default_table_access_method = test_tam",
			"CREATE TABLE tam_t (a int)",
			"CREATE TABLE tam_ctas AS SELECT 1",
		} {
			if code := simulationCatalogCode(t, c50, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c50, "DROP EXTENSION test_tam_ext"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION with default-AM relations answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c50, "SET default_table_access_method = missing_am"); code != CodeInvalidParameterValue {
			t.Fatalf("SET default_table_access_method = missing_am answered %s, want %s", code, CodeInvalidParameterValue)
		}

		// #184: a failed START crosscheck leaves the previous OWNED BY in place.
		c51 := New()
		for _, stmt := range []string{
			"CREATE TABLE own_a (id int)",
			"CREATE TABLE own_b (id int)",
			"CREATE SEQUENCE own_seq OWNED BY own_a.id",
		} {
			if code := simulationCatalogCode(t, c51, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c51, "ALTER SEQUENCE own_seq OWNED BY own_b.id START 1 MINVALUE 100"); code != CodeInvalidParameterValue {
			t.Fatalf("ALTER SEQUENCE with a bad START answered %s, want %s", code, CodeInvalidParameterValue)
		}
		if code := simulationCatalogCode(t, c51, "DROP TABLE own_a"); code != "0" {
			t.Fatalf("DROP TABLE own_a: %s", code)
		}
		if c51.schemaByName["public"].Sequences["own_seq"] != nil {
			t.Fatal("a failed ALTER SEQUENCE moved the OWNED BY link")
		}

		// #185: a failed CREATE SCHEMA does not record its AUTHORIZATION role.
		c52 := New()
		// The second element repeats the first table's name, so the whole
		// CREATE SCHEMA fails.
		if code := simulationCatalogCode(t, c52, "CREATE SCHEMA failed_s AUTHORIZATION ghost_role CREATE TABLE t1 (a int) CREATE TABLE t1 (a int)"); code == "0" {
			t.Fatal("CREATE SCHEMA with a bad element was accepted")
		}
		if c52.roles["ghost_role"] {
			t.Fatal("a failed CREATE SCHEMA recorded its AUTHORIZATION role")
		}

		// #186: a foreign-owned operator inside an owned schema blocks RESTRICT
		// and only goes with CASCADE.
		c53 := New()
		for _, stmt := range []string{
			"CREATE ROLE schema_owner",
			"CREATE ROLE other_role",
			"CREATE SCHEMA gate_s",
			"ALTER SCHEMA gate_s OWNER TO schema_owner",
			"CREATE FUNCTION gate_s.g(int, int) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"SET ROLE other_role",
			"CREATE OPERATOR gate_s.!! (LEFTARG = int, RIGHTARG = int, PROCEDURE = gate_s.g)",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c53, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c53, "DROP OWNED BY schema_owner RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED over a foreign operator answered %s, want %s", code, CodeDependentObjects)
		}
		if len(c53.operByKey[operKey{name: "!!", left: INT4OID, right: INT4OID}]) == 0 {
			t.Fatal("failed DROP OWNED removed the foreign operator")
		}
		if code := simulationCatalogCode(t, c53, "DROP OWNED BY schema_owner CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED CASCADE: %s", code)
		}

		// #187 & #189: opclasses and families in an owned schema, and users of
		// an owned extension's access method, follow DROP OWNED.
		c54 := New()
		for _, stmt := range []string{
			"CREATE ROLE fam_owner",
			"CREATE SCHEMA fam_s",
			"ALTER SCHEMA fam_s OWNER TO fam_owner",
			"CREATE OPERATOR FAMILY fam_s.fam1 USING gist",
		} {
			if code := simulationCatalogCode(t, c54, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c54, "DROP OWNED BY fam_owner RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED over an owned-schema operator family answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c54, "DROP OWNED BY fam_owner CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED CASCADE: %s", code)
		}

		// #190: DROP FUNCTION ... CASCADE removes an access method using it.
		c55 := New()
		for _, stmt := range []string{
			"CREATE FUNCTION am_handler(internal) RETURNS table_am_handler LANGUAGE internal AS 'bthandler'",
			"CREATE ACCESS METHOD cascade_am TYPE TABLE HANDLER am_handler",
			"DROP FUNCTION am_handler(internal) CASCADE",
		} {
			if code := simulationCatalogCode(t, c55, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if c55.accessMethodByName["cascade_am"] != nil {
			t.Fatal("DROP FUNCTION ... CASCADE left the access method")
		}

		// #192: a procedure and a function with the same signature collide.
		RegisterExtensionSQL("test_proc_ext", "CREATE PROCEDURE shared_sig(int) LANGUAGE sql AS 'SELECT 1';")
		defer delete(extensionScripts, "test_proc_ext")
		c56 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA psrc",
			"CREATE SCHEMA pdst",
			"CREATE EXTENSION test_proc_ext SCHEMA psrc",
			"CREATE FUNCTION pdst.shared_sig(int) RETURNS int LANGUAGE sql AS 'SELECT 2'",
		} {
			if code := simulationCatalogCode(t, c56, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c56, "ALTER EXTENSION test_proc_ext SET SCHEMA pdst"); code != CodeDuplicateFunction {
			t.Fatalf("cross-kind SET SCHEMA collision answered %s, want %s", code, CodeDuplicateFunction)
		}

		// #193: a schema holding only an operator class or family is not empty.
		c57 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA only_opc",
			"CREATE OPERATOR FAMILY only_opc.only_fam USING gist",
		} {
			if code := simulationCatalogCode(t, c57, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c57, "DROP SCHEMA only_opc"); code != CodeSchemaNotEmpty {
			t.Fatalf("DROP SCHEMA over an operator family answered %s, want %s", code, CodeSchemaNotEmpty)
		}
		if code := simulationCatalogCode(t, c57, "DROP SCHEMA only_opc CASCADE"); code != "0" {
			t.Fatalf("DROP SCHEMA only_opc CASCADE: %s", code)
		}
		for _, fam := range c57.opFamilies {
			if fam.Schema == c57.schemaByName["only_opc"].OID {
				t.Fatal("CASCADE left an operator family in the dropped schema")
			}
		}

		// #194: an extension's own opclass of its own access method does not
		// block DROP OWNED of the extension owner.
		RegisterExtensionSQL("test_self_am", "CREATE ACCESS METHOD self_am TYPE INDEX HANDLER bthandler; CREATE OPERATOR FAMILY self_fam USING self_am;")
		defer delete(extensionScripts, "test_self_am")
		c58 := New()
		for _, stmt := range []string{
			"CREATE ROLE self_owner",
			"SET ROLE self_owner",
			"CREATE EXTENSION test_self_am",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c58, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c58, "DROP OWNED BY self_owner RESTRICT"); code != "0" {
			t.Fatalf("DROP OWNED over a self-contained extension answered %s, want 0", code)
		}

		// #195: an operator class/family name collision blocks SET SCHEMA.
		RegisterExtensionSQL("test_opc_ext", "CREATE OPERATOR FAMILY shared_fam USING gist;")
		defer delete(extensionScripts, "test_opc_ext")
		c59 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA osrc",
			"CREATE SCHEMA odst",
			"CREATE EXTENSION test_opc_ext SCHEMA osrc",
			"CREATE OPERATOR FAMILY odst.shared_fam USING gist",
		} {
			if code := simulationCatalogCode(t, c59, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c59, "ALTER EXTENSION test_opc_ext SET SCHEMA odst"); code != CodeDuplicateObject {
			t.Fatalf("SET SCHEMA with an opfamily collision answered %s, want %s", code, CodeDuplicateObject)
		}

		// #198: a qualified ALTER OPERATOR names a missing schema with 3F000.
		c60 := New()
		if code := simulationCatalogCode(t, c60, "ALTER OPERATOR missing.@@ (int, int) OWNER TO postgres"); code != CodeUndefinedSchema {
			t.Fatalf("qualified ALTER OPERATOR on a missing schema answered %s, want %s", code, CodeUndefinedSchema)
		}

		// #199: an operator on an extension type keeps the extension in use.
		c61 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION citext",
			"CREATE FUNCTION ct_cmp(citext, citext) RETURNS int LANGUAGE sql IMMUTABLE AS 'SELECT 1'",
			"CREATE OPERATOR &&&& (LEFTARG = citext, RIGHTARG = citext, PROCEDURE = ct_cmp)",
		} {
			if code := simulationCatalogCode(t, c61, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c61, "DROP EXTENSION citext"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION over an operator on its type answered %s, want %s", code, CodeDependentObjects)
		}

		// #200: dropping a relation clears its array type's comments and deps.
		c62 := New()
		for _, stmt := range []string{
			"CREATE TABLE arr_t (a int)",
			"COMMENT ON TYPE _arr_t IS 'array'",
			"DROP TABLE arr_t",
		} {
			if code := simulationCatalogCode(t, c62, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		for key := range c62.comments {
			if key.ObjType == 't' {
				t.Fatalf("DROP TABLE left a type comment: %+v", key)
			}
		}
		for _, d := range c62.deps {
			if d.ObjType == 't' || d.RefType == 't' {
				t.Fatalf("DROP TABLE left a type dependency: %+v", d)
			}
		}

		// #201: an operator class on an extension type keeps the extension in
		// use, and #202 an external index on a schema's opclass goes with the
		// schema's CASCADE.
		c63 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION citext",
			"CREATE OPERATOR CLASS ct_ops FOR TYPE citext USING btree AS OPERATOR 1 <, FUNCTION 1 bttextcmp(text, text)",
		} {
			if code := simulationCatalogCode(t, c63, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c63, "DROP EXTENSION citext"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION over an opclass on its type answered %s, want %s", code, CodeDependentObjects)
		}

		c64 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA opc_s",
			"CREATE OPERATOR CLASS opc_s.sops FOR TYPE int USING btree AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int, int)",
			"CREATE TABLE opc_t (a int)",
			"CREATE INDEX opc_idx ON opc_t USING btree (a opc_s.sops)",
		} {
			if code := simulationCatalogCode(t, c64, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c64, "DROP SCHEMA opc_s CASCADE"); code != "0" {
			t.Fatalf("DROP SCHEMA opc_s CASCADE answered %s", code)
		}
		if c64.schemaByName["public"].Indexes["opc_idx"] != nil {
			t.Fatal("DROP SCHEMA CASCADE left the index that used the schema's opclass")
		}

		// #204: an extension member opclass inside an owned schema does not
		// block DROP OWNED of the schema owner.
		RegisterExtensionSQL("test_schema_ext", "CREATE OPERATOR FAMILY schema_fam USING gist;")
		defer delete(extensionScripts, "test_schema_ext")
		c65 := New()
		for _, stmt := range []string{
			"CREATE ROLE schema_ext_owner",
			"CREATE SCHEMA ext_own_s",
			"ALTER SCHEMA ext_own_s OWNER TO schema_ext_owner",
			"SET ROLE schema_ext_owner",
			"CREATE EXTENSION test_schema_ext SCHEMA ext_own_s",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c65, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c65, "DROP OWNED BY schema_ext_owner RESTRICT"); code != "0" {
			t.Fatalf("DROP OWNED over an owned schema's extension answered %s, want 0", code)
		}

		// #205: RESTRICT refuses revoking a grant the owner had passed on.
		c66 := New()
		for _, stmt := range []string{
			"CREATE ROLE g_alice",
			"CREATE ROLE g_bob",
			"CREATE TABLE g_chain (id int)",
			"GRANT SELECT ON g_chain TO g_alice WITH GRANT OPTION",
			"SET ROLE g_alice",
			"GRANT SELECT ON g_chain TO g_bob",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c66, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c66, "DROP OWNED BY g_alice RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED over a passed-on grant answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c66, "DROP OWNED BY g_alice CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED CASCADE over a passed-on grant answered %s", code)
		}

		// #206: ALTER TABLE OWNER TO moves the grants with the object.
		c67 := New()
		for _, stmt := range []string{
			"CREATE ROLE acl_old",
			"CREATE ROLE acl_new",
			"CREATE ROLE acl_user",
			"CREATE TABLE acl_move (id int)",
			"SET ROLE acl_old",
			"ALTER TABLE acl_move OWNER TO acl_old",
			"GRANT SELECT ON acl_move TO acl_user",
			"ALTER TABLE acl_move OWNER TO acl_new",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c67, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c67, "DROP ROLE acl_old"); code != "0" {
			t.Fatalf("DROP ROLE acl_old after the owner transfer answered %s, want 0", code)
		}

		// #207: the default table access method does not reach composite types.
		RegisterExtensionSQL("test_tam_ext2", "CREATE ACCESS METHOD test_tam2 TYPE TABLE HANDLER heap_tableam_handler;")
		defer delete(extensionScripts, "test_tam_ext2")
		c68 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION test_tam_ext2",
			"SET default_table_access_method = test_tam2",
			"CREATE TYPE ct_no_am AS (a int)",
		} {
			if code := simulationCatalogCode(t, c68, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c68, "DROP EXTENSION test_tam_ext2"); code != "0" {
			t.Fatalf("DROP EXTENSION after a composite type answered %s, want 0", code)
		}

		// #208: an operator class that joins an extension family keeps it in use.
		RegisterExtensionSQL("test_family_ext", "CREATE OPERATOR FAMILY ext_fam USING btree;")
		defer delete(extensionScripts, "test_family_ext")
		c69 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION test_family_ext",
			"CREATE OPERATOR CLASS own_ops FOR TYPE int USING btree FAMILY ext_fam AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int, int)",
		} {
			if code := simulationCatalogCode(t, c69, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c69, "DROP EXTENSION test_family_ext"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION over an opclass in its family answered %s, want %s", code, CodeDependentObjects)
		}

		// #209: DROP OWNED CASCADE of a schema's opclass drops the external index
		// that uses it.
		c70 := New()
		for _, stmt := range []string{
			"CREATE ROLE opc_schema_owner",
			"CREATE SCHEMA owned_opc_s",
			"ALTER SCHEMA owned_opc_s OWNER TO opc_schema_owner",
			"CREATE OPERATOR CLASS owned_opc_s.owned_ops FOR TYPE int USING btree AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int, int)",
			"CREATE TABLE owned_opc_t (a int)",
			"CREATE INDEX owned_opc_idx ON owned_opc_t USING btree (a owned_opc_s.owned_ops)",
		} {
			if code := simulationCatalogCode(t, c70, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c70, "DROP OWNED BY opc_schema_owner RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED with an external opclass index answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c70, "DROP OWNED BY opc_schema_owner CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED CASCADE: %s", code)
		}
		if c70.schemaByName["public"].Indexes["owned_opc_idx"] != nil {
			t.Fatal("DROP OWNED CASCADE left the external index")
		}

		// #210: a failed grant-chain preflight does not delete owned objects.
		c71 := New()
		for _, stmt := range []string{
			"CREATE ROLE chain_alice",
			"CREATE ROLE chain_bob",
			"SET ROLE chain_alice",
			"CREATE TABLE chain_owned (id int)",
			"RESET ROLE",
			"CREATE TABLE chain_obj (id int)",
			"GRANT SELECT ON chain_obj TO chain_alice WITH GRANT OPTION",
			"SET ROLE chain_alice",
			"GRANT SELECT ON chain_obj TO chain_bob",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c71, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c71, "DROP OWNED BY chain_alice RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED over a grant chain answered %s, want %s", code, CodeDependentObjects)
		}
		if c71.schemaByName["public"].Relations["chain_owned"] == nil {
			t.Fatal("a failed DROP OWNED deleted the owned table")
		}

		// #211: REASSIGN OWNED moves the ACL identities too.
		c72 := New()
		for _, stmt := range []string{
			"CREATE ROLE ra_alice",
			"CREATE ROLE ra_bob",
			"CREATE ROLE ra_user",
			"CREATE TABLE ra_t (id int)",
			"ALTER TABLE ra_t OWNER TO ra_alice",
			"SET ROLE ra_alice",
			"GRANT SELECT ON ra_t TO ra_user",
			"RESET ROLE",
			"REASSIGN OWNED BY ra_alice TO ra_bob",
		} {
			if code := simulationCatalogCode(t, c72, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c72, "DROP ROLE ra_alice"); code != "0" {
			t.Fatalf("DROP ROLE ra_alice after REASSIGN answered %s, want 0", code)
		}

		// #212: removing an AM user clears its dependency rows.
		c73 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA tmp_fam_s",
			"CREATE OPERATOR FAMILY tmp_fam_s.tmp_fam USING gist",
			"DROP SCHEMA tmp_fam_s CASCADE",
		} {
			if code := simulationCatalogCode(t, c73, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		for _, d := range c73.deps {
			if d.ObjType == 'Q' || d.RefType == 'Q' {
				t.Fatalf("dropped family left a dependency row: %+v", d)
			}
		}

		// #213: DROP SCHEMA CASCADE of an opclass removes the constraint that
		// backed its index.
		c74 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA excl_opc_s",
			"CREATE OPERATOR CLASS excl_opc_s.excl_ops FOR TYPE int USING gist AS OPERATOR 1 &&, FUNCTION 1 btint4cmp(int, int)",
			"CREATE TABLE excl_opc_t (a int)",
			"ALTER TABLE excl_opc_t ADD CONSTRAINT excl_opc_con EXCLUDE USING gist (a excl_opc_s.excl_ops WITH &&)",
			"DROP SCHEMA excl_opc_s CASCADE",
		} {
			if code := simulationCatalogCode(t, c74, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if rel := c74.schemaByName["public"].Relations["excl_opc_t"]; rel != nil && len(c74.consByRel[rel.OID]) != 0 {
			t.Fatal("DROP SCHEMA CASCADE left the exclusion constraint behind")
		}

		// #214: DROP EXTENSION removes member triggers and policies.
		RegisterExtensionSQL("test_member_ext", "CREATE FUNCTION member_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END'; CREATE TABLE pre_made (a int); CREATE TRIGGER member_trig AFTER INSERT ON pre_made FOR EACH ROW EXECUTE FUNCTION member_fn();")
		defer delete(extensionScripts, "test_member_ext")
		c75 := New()
		for _, stmt := range []string{"CREATE EXTENSION test_member_ext", "DROP EXTENSION test_member_ext"} {
			if code := simulationCatalogCode(t, c75, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if len(c75.triggers) != 0 {
			t.Fatalf("DROP EXTENSION left %d triggers", len(c75.triggers))
		}

		// #215: DROP SCHEMA CASCADE of a family takes external classes and their
		// indexes with it.
		c76 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA fam_casc_s",
			"CREATE OPERATOR FAMILY fam_casc_s.fam2 USING btree",
			"CREATE OPERATOR CLASS ext_cls FOR TYPE int USING btree FAMILY fam_casc_s.fam2 AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int, int)",
			"CREATE TABLE fam_casc_t (a int)",
			"CREATE INDEX fam_casc_idx ON fam_casc_t USING btree (a ext_cls)",
			"DROP SCHEMA fam_casc_s CASCADE",
		} {
			if code := simulationCatalogCode(t, c76, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if c76.schemaByName["public"].Indexes["fam_casc_idx"] != nil || len(c76.opClasses) != 0 {
			t.Fatal("DROP SCHEMA CASCADE left the external class or its index")
		}

		// #216: policy roles are observed, so CREATE ROLE reports the duplicate.
		c77 := New()
		for _, stmt := range []string{
			"CREATE TABLE pol_t (a int, owner text)",
			"CREATE POLICY pol_ghost ON pol_t TO ghost_role USING (true)",
		} {
			if code := simulationCatalogCode(t, c77, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c77, "CREATE ROLE ghost_role"); code != CodeDuplicateObject {
			t.Fatalf("CREATE ROLE ghost_role answered %s, want %s", code, CodeDuplicateObject)
		}

		// #130: failed ALTER SEQUENCE ... OWNED BY preserves the previous
		// automatic column dependency.
		c11 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"SET ROLE alice",
			"CREATE TABLE old_t (id int)",
			"CREATE SEQUENCE s_owned OWNED BY old_t.id",
			"SET ROLE bob",
			"CREATE TABLE new_t (id int)",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c11, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c11, "ALTER SEQUENCE s_owned OWNED BY new_t.id"); code != CodeObjectNotInPrerequisiteState {
			t.Fatalf("ALTER SEQUENCE OWNED BY mismatch answered %s, want %s", code, CodeObjectNotInPrerequisiteState)
		}
		if code := simulationCatalogCode(t, c11, "DROP TABLE old_t"); code != "0" {
			t.Fatalf("DROP TABLE old_t: %s", code)
		}
		if c11.schemaByName["public"].Sequences["s_owned"] != nil {
			t.Fatal("failed ALTER SEQUENCE OWNED BY lost previous DepAuto link to old_t")
		}

		// #131: completing a shell type under another role fails with 42501.
		c12 := New()
		for _, stmt := range []string{
			"CREATE ROLE alice",
			"CREATE ROLE bob",
			"SET ROLE alice",
			"CREATE TYPE sh_t",
			"CREATE FUNCTION sh_t_in(cstring) RETURNS sh_t LANGUAGE internal IMMUTABLE STRICT AS 'int4in'",
			"CREATE FUNCTION sh_t_out(sh_t) RETURNS cstring LANGUAGE internal IMMUTABLE STRICT AS 'int4out'",
			"SET ROLE bob",
		} {
			if code := simulationCatalogCode(t, c12, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c12, "CREATE TYPE sh_t (INPUT = sh_t_in, OUTPUT = sh_t_out)"); code != CodeInsufficientPrivilege {
			t.Fatalf("shell completion by non-owner answered %s, want %s", code, CodeInsufficientPrivilege)
		}

		// #218: operator class and family ownership is tracked across DROP ROLE,
		// ALTER ... OWNER TO, RENAME ROLE, REASSIGN OWNED, and DROP OWNED.
		c78 := New()
		for _, stmt := range []string{
			"CREATE ROLE opc_owner",
			"CREATE ROLE opc_new",
			"SET ROLE opc_owner",
			"CREATE OPERATOR FAMILY own_fam USING btree",
			"CREATE OPERATOR CLASS own_cls FOR TYPE int USING btree FAMILY own_fam AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int, int)",
			"RESET ROLE",
		} {
			if code := simulationCatalogCode(t, c78, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c78, "DROP ROLE opc_owner"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE owning opclass/opfamily answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c78, "ALTER OPERATOR CLASS own_cls USING btree OWNER TO opc_new"); code != "0" {
			t.Fatalf("ALTER OPERATOR CLASS OWNER TO: %s", code)
		}
		if code := simulationCatalogCode(t, c78, "DROP ROLE opc_owner"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE still owning opfamily answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c78, "ALTER OPERATOR FAMILY own_fam USING btree OWNER TO opc_new"); code != "0" {
			t.Fatalf("ALTER OPERATOR FAMILY OWNER TO: %s", code)
		}
		if code := simulationCatalogCode(t, c78, "DROP ROLE opc_owner"); code != "0" {
			t.Fatalf("DROP ROLE after transferring opclass/opfamily answered %s, want 0", code)
		}
		for _, stmt := range []string{
			"ALTER ROLE opc_new RENAME TO opc_renamed",
			"REASSIGN OWNED BY opc_renamed TO postgres",
			"DROP ROLE opc_renamed",
		} {
			if code := simulationCatalogCode(t, c78, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}

		c78b := New()
		for _, stmt := range []string{
			"CREATE ROLE fam_drop_role",
			"SET ROLE fam_drop_role",
			"CREATE SCHEMA fam_drop_s",
			"CREATE OPERATOR FAMILY fam_drop_s.f1 USING btree",
			"RESET ROLE",
			"DROP OWNED BY fam_drop_role RESTRICT",
			"DROP ROLE fam_drop_role",
		} {
			if code := simulationCatalogCode(t, c78b, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}

		// #219: DROP OWNED CASCADE on an owned schema's operator family cascades
		// external Q -> q operator classes and their o -> i indexes + constraints,
		// and clears dependency rows.
		c79 := New()
		for _, stmt := range []string{
			"CREATE ROLE fam_schema_owner",
			"CREATE SCHEMA owned_fam_s AUTHORIZATION fam_schema_owner",
			"CREATE OPERATOR FAMILY owned_fam_s.fam USING gist",
			"CREATE OPERATOR CLASS public.ext_gist_ops FOR TYPE int USING gist FAMILY owned_fam_s.fam AS OPERATOR 1 &&, FUNCTION 1 btint4cmp(int, int)",
			"CREATE TABLE public.fam_idx_t (a int, CONSTRAINT fam_excl EXCLUDE USING gist (a public.ext_gist_ops WITH &&))",
		} {
			if code := simulationCatalogCode(t, c79, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c79, "DROP OWNED BY fam_schema_owner RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED RESTRICT with external opclass in schema family answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c79, "DROP OWNED BY fam_schema_owner CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED CASCADE: %s", code)
		}
		if len(c79.opFamilies) != 0 || len(c79.opClasses) != 0 {
			t.Fatalf("DROP OWNED CASCADE left opfamilies=%d opclasses=%d", len(c79.opFamilies), len(c79.opClasses))
		}
		if rel := c79.schemaByName["public"].Relations["fam_idx_t"]; rel != nil && (len(c79.indexesByRel[rel.OID]) != 0 || len(c79.consByRel[rel.OID]) != 0) {
			t.Fatal("DROP OWNED CASCADE left index or backing exclusion constraint")
		}
		for _, d := range c79.deps {
			if d.ObjType == 'Q' || d.RefType == 'Q' || d.ObjType == 'q' || d.RefType == 'q' || d.RefType == 'o' {
				t.Fatalf("DROP OWNED CASCADE left stale opclass/opfamily dependency: %+v", d)
			}
		}

		// #220: DROP EXTENSION clears dependency rows for member opclasses,
		// opfamilies, and access methods.
		RegisterExtensionSQL("test_ext_dep_clean", `
			CREATE FUNCTION test_iam_h(internal) RETURNS index_am_handler LANGUAGE internal AS 'bthandler';
			CREATE ACCESS METHOD ext_iam TYPE INDEX HANDLER test_iam_h;
			CREATE OPERATOR FAMILY ext_iam_fam USING ext_iam;
			CREATE OPERATOR CLASS ext_iam_ops DEFAULT FOR TYPE int4 USING ext_iam FAMILY ext_iam_fam AS STORAGE int4;
		`)
		defer delete(extensionScripts, "test_ext_dep_clean")
		c80 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION test_ext_dep_clean",
			"DROP EXTENSION test_ext_dep_clean",
		} {
			if code := simulationCatalogCode(t, c80, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		for _, d := range c80.deps {
			if d.ObjType == 'q' || d.RefType == 'q' || d.ObjType == 'Q' || d.RefType == 'Q' || d.ObjType == 'm' || d.RefType == 'm' || d.RefType == 'o' {
				t.Fatalf("DROP EXTENSION left stale AM/opclass/opfamily dependency row: %+v", d)
			}
		}

		// #221: LIKE ... INCLUDING INDEXES records default (column-type-resolved)
		// operator class dependencies on cloned indexes when idx.OpClasses is nil.
		RegisterExtensionSQL("test_like_def_opc_ext", `
			CREATE TYPE like_dom AS (v int4);
			CREATE OPERATOR CLASS like_dom_ops DEFAULT FOR TYPE int4 USING gist AS OPERATOR 1 &&, FUNCTION 1 btint4cmp(int4, int4);
		`)
		defer delete(extensionScripts, "test_like_def_opc_ext")
		c81 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION test_like_def_opc_ext",
			"CREATE TABLE like_src (a int4)",
			"CREATE INDEX like_src_gist ON like_src USING gist (a)",
			"CREATE TABLE like_dst (LIKE like_src INCLUDING INDEXES)",
			"DROP TABLE like_src",
		} {
			if code := simulationCatalogCode(t, c81, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c81, "DROP EXTENSION test_like_def_opc_ext RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION RESTRICT over LIKE-cloned default-opclass index answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c81, "DROP EXTENSION test_like_def_opc_ext CASCADE"); code != "0" {
			t.Fatalf("DROP EXTENSION CASCADE: %s", code)
		}
		if dst := c81.schemaByName["public"].Relations["like_dst"]; dst != nil && len(c81.indexesByRel[dst.OID]) != 0 {
			t.Fatal("DROP EXTENSION CASCADE left LIKE-cloned index")
		}

		// Extension OpFamily external class dependent in DropOwned + 3-hop grant
		// chain + snapshot role/tablespace tombstones + SET SCHEMA preflights.
		RegisterExtensionSQL("test_owned_fam_ext", "CREATE OPERATOR FAMILY owned_ext_fam USING btree;")
		defer delete(extensionScripts, "test_owned_fam_ext")
		c82 := New()
		for _, stmt := range []string{
			"CREATE ROLE ext_fam_r",
			"SET ROLE ext_fam_r",
			"CREATE EXTENSION test_owned_fam_ext",
			"RESET ROLE",
			"CREATE OPERATOR CLASS ext_fam_user_opc FOR TYPE int4 USING btree FAMILY owned_ext_fam AS OPERATOR 1 <, FUNCTION 1 btint4cmp(int4, int4)",
		} {
			if code := simulationCatalogCode(t, c82, stmt); code != "0" {
				t.Fatalf("catalog refused %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c82, "DROP OWNED BY ext_fam_r RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED RESTRICT with external opclass in owned ext family answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c82, "DROP OWNED BY ext_fam_r CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED CASCADE: %s", code)
		}
		if len(c82.opClasses) != 0 {
			t.Fatal("DROP OWNED CASCADE left external opclass in owned ext family")
		}

		c83 := New()
		for _, stmt := range []string{
			"CREATE ROLE hop1",
			"CREATE ROLE hop2",
			"CREATE ROLE hop3",
			"CREATE TABLE hop_t (a int)",
			"GRANT SELECT ON hop_t TO hop1 WITH GRANT OPTION",
			"SET ROLE hop1",
			"GRANT SELECT ON hop_t TO hop2 WITH GRANT OPTION",
			"SET ROLE hop2",
			"GRANT SELECT ON hop_t TO hop3",
			"RESET ROLE",
			"DROP OWNED BY hop1 CASCADE",
			"DROP ROLE hop2",
		} {
			if code := simulationCatalogCode(t, c83, stmt); code != "0" {
				t.Fatalf("3-hop grant revocation failed on %s with %s", stmt, code)
			}
		}

		c84 := New()
		c84.snapshot = true
		for _, stmt := range []string{
			"CREATE ROLE snap_r",
			"DROP ROLE snap_r",
		} {
			if code := simulationCatalogCode(t, c84, stmt); code != "0" {
				t.Fatalf("snapshot role setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c84, "DROP ROLE snap_r"); code != CodeUndefinedObject {
			t.Fatalf("second DROP ROLE on snapshot answered %s, want %s", code, CodeUndefinedObject)
		}
		if code := simulationCatalogCode(t, c84, "ALTER ROLE unobs_r RENAME TO unobs_r"); code != CodeDuplicateObject {
			t.Fatalf("self RENAME ROLE on snapshot answered %s, want %s", code, CodeDuplicateObject)
		}

		c85 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA mv_s1",
			"CREATE SCHEMA mv_s2",
			"CREATE TYPE mv_s1.rng AS RANGE (SUBTYPE = int4)",
			"ALTER TYPE mv_s1.rng SET SCHEMA mv_s2",
		} {
			if code := simulationCatalogCode(t, c85, stmt); code != "0" {
				t.Fatalf("range SET SCHEMA setup failed on %s with %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c85, "DROP SCHEMA mv_s1 RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP SCHEMA mv_s1 while multirange remains answered %s, want %s", code, CodeDependentObjects)
		}
		for _, stmt := range []string{
			"ALTER TYPE mv_s1.rng_multirange SET SCHEMA mv_s2",
			"DROP SCHEMA mv_s1 RESTRICT",
		} {
			if code := simulationCatalogCode(t, c85, stmt); code != "0" {
				t.Fatalf("multirange SET SCHEMA + DROP SCHEMA failed on %s with %s", stmt, code)
			}
		}

		// Access method handler signature/return-type validation.
		c86 := New()
		for _, stmt := range []string{
			"CREATE FUNCTION bad_am_h0() RETURNS void LANGUAGE sql AS 'SELECT 1'",
			"CREATE FUNCTION bad_am_h1(internal) RETURNS int4 LANGUAGE internal AS 'bthandler'",
		} {
			if code := simulationCatalogCode(t, c86, stmt); code != "0" {
				t.Fatalf("AM handler test setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c86, "CREATE ACCESS METHOD bad_am0 TYPE TABLE HANDLER bad_am_h0"); code != CodeUndefinedFunction {
			t.Fatalf("AM handler wrong-arg answered %s, want %s", code, CodeUndefinedFunction)
		}
		if code := simulationCatalogCode(t, c86, "CREATE ACCESS METHOD bad_am1 TYPE TABLE HANDLER bad_am_h1"); code != CodeWrongObjectType {
			t.Fatalf("AM handler wrong-rettype answered %s, want %s", code, CodeWrongObjectType)
		}
		if code := simulationCatalogCode(t, c86, "DROP FUNCTION bad_am_h0()"); code != "0" {
			t.Fatalf("DROP FUNCTION bad_am_h0() answered %s, want 0", code)
		}

		// Partitioned FK clones record dependencies on all equality operator groups
		// (PFEqOp, PPEqOp, FFEqOp).
		RegisterExtensionSQL("test_fk_ffeq_ext", `
			CREATE FUNCTION fk_ff_eq(fk_ext_u, fk_ext_u) RETURNS bool LANGUAGE sql IMMUTABLE AS 'SELECT true';
			CREATE OPERATOR = (LEFTARG = fk_ext_u, RIGHTARG = fk_ext_u, FUNCTION = fk_ff_eq, COMMUTATOR = =, HASHES, MERGES);
		`)
		defer delete(extensionScripts, "test_fk_ffeq_ext")
		c87 := New()
		for _, stmt := range []string{
			"CREATE TYPE fk_ext_u AS (v int4)",
			"CREATE FUNCTION fk_pf_eq(int4, fk_ext_u) RETURNS bool LANGUAGE sql IMMUTABLE AS 'SELECT true'",
			"CREATE OPERATOR = (LEFTARG = int4, RIGHTARG = fk_ext_u, FUNCTION = fk_pf_eq)",
			"CREATE EXTENSION test_fk_ffeq_ext",
			"CREATE TABLE fk_ref (id int4 PRIMARY KEY)",
			"CREATE TABLE fk_par (p int4, r fk_ext_u, CONSTRAINT fk_par_c FOREIGN KEY (r) REFERENCES fk_ref(id)) PARTITION BY LIST (p)",
			"CREATE TABLE fk_ch PARTITION OF fk_par FOR VALUES IN (1)",
			"ALTER TABLE fk_par DETACH PARTITION fk_ch",
			"DROP TABLE fk_par",
		} {
			if code := simulationCatalogCode(t, c87, stmt); code != "0" {
				t.Fatalf("partition FK op setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c87, "DROP EXTENSION test_fk_ffeq_ext RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION RESTRICT over detached partition FFEqOp answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c87, "DROP EXTENSION test_fk_ffeq_ext CASCADE"); code != "0" {
			t.Fatalf("DROP EXTENSION CASCADE: %s", code)
		}

		// Partitioned trigger creation preflights existing child trigger collisions.
		c88 := New()
		for _, stmt := range []string{
			"CREATE FUNCTION pt_trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END'",
			"CREATE TABLE pt_par (a int) PARTITION BY LIST (a)",
			"CREATE TABLE pt_ch1 PARTITION OF pt_par FOR VALUES IN (1)",
			"CREATE TABLE pt_ch2 PARTITION OF pt_par FOR VALUES IN (2)",
			"CREATE TRIGGER dup_tr BEFORE INSERT ON pt_ch2 FOR EACH ROW EXECUTE FUNCTION pt_trig_fn()",
		} {
			if code := simulationCatalogCode(t, c88, stmt); code != "0" {
				t.Fatalf("partition trigger setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c88, "CREATE TRIGGER dup_tr BEFORE INSERT ON pt_par FOR EACH ROW EXECUTE FUNCTION pt_trig_fn()"); code != CodeDuplicateObject {
			t.Fatalf("CREATE TRIGGER colliding with child partition trigger answered %s, want %s", code, CodeDuplicateObject)
		}
		if par := c88.schemaByName["public"].Relations["pt_par"]; par != nil && len(c88.triggersByRel[par.OID]) != 0 {
			t.Fatal("failed parent CREATE TRIGGER registered trigger on parent")
		}
		if ch1 := c88.schemaByName["public"].Relations["pt_ch1"]; ch1 != nil && len(c88.triggersByRel[ch1.OID]) != 0 {
			t.Fatal("failed parent CREATE TRIGGER cloned trigger onto earlier sibling partition")
		}

		// CREATE OR REPLACE TRIGGER on a partitioned parent updates an existing
		// child trigger in place (preserving OID and comment), and DETACH PARTITION
		// removes cloned triggers so a subsequent reattach succeeds.
		c89 := New()
		for _, stmt := range []string{
			"CREATE FUNCTION pt_rep_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END'",
			"CREATE TABLE pt_rep (a int) PARTITION BY LIST (a)",
			"CREATE TABLE pt_rep_c PARTITION OF pt_rep FOR VALUES IN (1)",
			"CREATE TRIGGER rep_tr BEFORE INSERT ON pt_rep_c FOR EACH ROW EXECUTE FUNCTION pt_rep_fn()",
			"COMMENT ON TRIGGER rep_tr ON pt_rep_c IS 'keep me'",
		} {
			if code := simulationCatalogCode(t, c89, stmt); code != "0" {
				t.Fatalf("trigger replace setup %s: %s", stmt, code)
			}
		}
		repChild := c89.schemaByName["public"].Relations["pt_rep_c"]
		origChildTrigOID := c89.triggersByRel[repChild.OID][0].OID
		if code := simulationCatalogCode(t, c89, "CREATE OR REPLACE TRIGGER rep_tr AFTER INSERT ON pt_rep FOR EACH ROW EXECUTE FUNCTION pt_rep_fn()"); code != "0" {
			t.Fatalf("CREATE OR REPLACE TRIGGER on parent: %s", code)
		}
		if len(c89.triggersByRel[repChild.OID]) != 1 || c89.triggersByRel[repChild.OID][0].OID != origChildTrigOID {
			t.Fatal("recursive OR REPLACE TRIGGER did not preserve child trigger OID")
		}
		if cmt, ok := c89.GetComment('g', origChildTrigOID, 0); !ok || cmt != "keep me" {
			t.Fatalf("recursive OR REPLACE TRIGGER lost child trigger comment: %q (%v)", cmt, ok)
		}
		if code := simulationCatalogCode(t, c89, "CREATE OR REPLACE TRIGGER rep_tr BEFORE INSERT ON pt_rep_c FOR EACH ROW EXECUTE FUNCTION pt_rep_fn()"); code != CodeDuplicateObject {
			t.Fatalf("direct OR REPLACE TRIGGER on child partition answered %s, want %s", code, CodeDuplicateObject)
		}
		for _, stmt := range []string{
			"ALTER TABLE pt_rep DETACH PARTITION pt_rep_c",
			"ALTER TABLE pt_rep ATTACH PARTITION pt_rep_c FOR VALUES IN (1)",
		} {
			if code := simulationCatalogCode(t, c89, stmt); code != "0" {
				t.Fatalf("DETACH + re-ATTACH with cloned trigger failed on %s: %s", stmt, code)
			}
		}

		// DROP SCHEMA CASCADE removes stranded multirange companions and their
		// owning range type in another schema.
		c90 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA mr_old",
			"CREATE SCHEMA mr_new",
			"CREATE TYPE mr_old.r AS RANGE (SUBTYPE = int4)",
			"ALTER TYPE mr_old.r SET SCHEMA mr_new",
			"DROP SCHEMA mr_old CASCADE",
		} {
			if code := simulationCatalogCode(t, c90, stmt); code != "0" {
				t.Fatalf("stranded multirange DROP SCHEMA CASCADE %s: %s", stmt, code)
			}
		}
		if c90.typeByName[typeKey{ns: c90.schemaByName["mr_new"].OID, name: "r"}] != nil {
			t.Fatal("DROP SCHEMA CASCADE of stranded multirange left owning range type alive")
		}
		for _, bt := range c90.typeByOID {
			if bt.OID >= FirstNormalObjectId {
				t.Fatalf("DROP SCHEMA CASCADE left user type %s (OID %d)", bt.TypeName, bt.OID)
			}
		}

		// DROP TABLE RESTRICT honors user casts referencing the table's row type.
		c91 := New()
		for _, stmt := range []string{
			"CREATE TABLE cast_rel (a int)",
			"CREATE CAST (cast_rel AS text) WITH INOUT",
		} {
			if code := simulationCatalogCode(t, c91, stmt); code != "0" {
				t.Fatalf("table rowtype cast setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c91, "DROP TABLE cast_rel RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP TABLE RESTRICT with row-type cast answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c91, "DROP TABLE cast_rel CASCADE"); code != "0" {
			t.Fatalf("DROP TABLE CASCADE with row-type cast: %s", code)
		}

		// Failed ALTER POLICY preserves existing predicate and expression dependencies.
		c92 := New()
		for _, stmt := range []string{
			"CREATE FUNCTION pol_guard_fn(int) RETURNS bool LANGUAGE sql AS 'SELECT true'",
			"CREATE TABLE pol_guard_t (a int)",
			"CREATE POLICY pol_g ON pol_guard_t AS RESTRICTIVE USING (pol_guard_fn(a))",
		} {
			if code := simulationCatalogCode(t, c92, stmt); code != "0" {
				t.Fatalf("policy guard setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c92, "ALTER POLICY pol_g ON pol_guard_t USING (missing_col > 0)"); code != CodeUndefinedColumn {
			t.Fatalf("ALTER POLICY with missing column answered %s, want %s", code, CodeUndefinedColumn)
		}
		if code := simulationCatalogCode(t, c92, "DROP FUNCTION pol_guard_fn(int) RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("failed ALTER POLICY erased existing function dependency (got %s)", code)
		}

		// Identity sequences added via ALTER TABLE inherit the table owner rather
		// than the executing role, so DROP OWNED BY the executing role leaves them intact.
		c93 := New()
		for _, stmt := range []string{
			"CREATE ROLE id_owner",
			"CREATE ROLE id_admin",
			"SET ROLE id_owner",
			"CREATE TABLE id_tbl (x int NOT NULL)",
			"SET ROLE id_admin",
			"ALTER TABLE id_tbl ADD COLUMN id1 int GENERATED ALWAYS AS IDENTITY",
			"ALTER TABLE id_tbl ALTER COLUMN x ADD GENERATED BY DEFAULT AS IDENTITY",
			"DROP OWNED BY id_admin",
		} {
			if code := simulationCatalogCode(t, c93, stmt); code != "0" {
				t.Fatalf("identity sequence owner %s: %s", stmt, code)
			}
		}
		for _, seqName := range []string{"id_tbl_id1_seq", "id_tbl_x_seq"} {
			seq := c93.schemaByName["public"].Sequences[seqName]
			if seq == nil || seq.OwnerName != "id_owner" {
				t.Fatalf("sequence %s missing or wrong owner after DROP OWNED BY id_admin: %+v", seqName, seq)
			}
		}

		// ALTER OPERATOR (type, NONE) OWNER TO preserves the postfix unary signature
		// rather than mutating a same-named prefix unary operator.
		c94 := New()
		for _, stmt := range []string{
			"CREATE ROLE op_alice",
			"CREATE ROLE op_bob",
			"SET ROLE op_alice",
			"CREATE FUNCTION ufn(int) RETURNS int LANGUAGE sql AS 'SELECT $1'",
			"CREATE OPERATOR #@ (RIGHTARG = int, FUNCTION = ufn)",
		} {
			if code := simulationCatalogCode(t, c94, stmt); code != "0" {
				t.Fatalf("unary operator owner setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c94, "ALTER OPERATOR #@ (int, NONE) OWNER TO op_bob"); code != CodeUndefinedObject {
			t.Fatalf("ALTER OPERATOR postfix signature answered %s, want %s", code, CodeUndefinedObject)
		}
		prefOp := c94.operByKey[operKey{name: "#@", left: 0, right: INT4OID}][0]
		if c94.operatorOwners[prefOp.OID] != "op_alice" {
			t.Fatalf("postfix ALTER OPERATOR mutated prefix operator owner to %q", c94.operatorOwners[prefOp.OID])
		}
		if code := simulationCatalogCode(t, c94, "ALTER OPERATOR #@ (NONE, int) OWNER TO op_bob"); code != "0" {
			t.Fatalf("ALTER OPERATOR prefix signature: %s", code)
		}
		if c94.operatorOwners[prefOp.OID] != "op_bob" {
			t.Fatalf("prefix ALTER OPERATOR did not update owner: %q", c94.operatorOwners[prefOp.OID])
		}

		// GRANT enforces routine/domain object kinds on built-in and user objects.
		c95 := New()
		if code := simulationCatalogCode(t, c95, "GRANT EXECUTE ON PROCEDURE pg_catalog.abs(integer) TO PUBLIC"); code != CodeWrongObjectType {
			t.Fatalf("GRANT EXECUTE ON PROCEDURE pg_catalog.abs(integer) answered %s, want %s", code, CodeWrongObjectType)
		}
		if code := simulationCatalogCode(t, c95, "GRANT EXECUTE ON PROCEDURE pg_catalog.clock_timestamp TO PUBLIC"); code != CodeUndefinedFunction {
			t.Fatalf("GRANT EXECUTE ON PROCEDURE pg_catalog.clock_timestamp answered %s, want %s", code, CodeUndefinedFunction)
		}
		if code := simulationCatalogCode(t, c95, "GRANT EXECUTE ON FUNCTION pg_catalog.abs TO PUBLIC"); code != CodeAmbiguousFunction {
			t.Fatalf("GRANT EXECUTE ON FUNCTION pg_catalog.abs answered %s, want %s", code, CodeAmbiguousFunction)
		}
		if code := simulationCatalogCode(t, c95, "GRANT USAGE ON DOMAIN int4 TO PUBLIC"); code != CodeWrongObjectType {
			t.Fatalf("GRANT USAGE ON DOMAIN int4 answered %s, want %s", code, CodeWrongObjectType)
		}
		if code := simulationCatalogCode(t, c95, "GRANT EXECUTE ON FUNCTION pg_catalog.abs(integer) TO PUBLIC"); code != "0" {
			t.Fatalf("GRANT EXECUTE ON FUNCTION pg_catalog.abs(integer): %s", code)
		}

		// Extension-created schemas are pinned while the extension exists, honor
		// RESTRICT when non-member objects reside in them, and are removed on DROP EXTENSION.
		RegisterExtensionSQL("test_schema_ext", "CREATE SCHEMA ext_aux_ns; CREATE TABLE ext_aux_ns.ext_t (a int);")
		defer delete(extensionScripts, "test_schema_ext")
		c96 := New()
		if code := simulationCatalogCode(t, c96, "CREATE EXTENSION test_schema_ext"); code != "0" {
			t.Fatalf("CREATE EXTENSION test_schema_ext: %s", code)
		}
		if code := simulationCatalogCode(t, c96, "DROP SCHEMA ext_aux_ns CASCADE"); code != CodeDependentObjects {
			t.Fatalf("DROP SCHEMA on extension-owned schema answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c96, "CREATE TABLE ext_aux_ns.usr_t (a int)"); code != "0" {
			t.Fatalf("CREATE TABLE ext_aux_ns.usr_t: %s", code)
		}
		if code := simulationCatalogCode(t, c96, "DROP EXTENSION test_schema_ext RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION RESTRICT with user table in ext schema answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c96, "DROP EXTENSION test_schema_ext CASCADE"); code != "0" {
			t.Fatalf("DROP EXTENSION CASCADE with user table in ext schema: %s", code)
		}
		if c96.schemaByName["ext_aux_ns"] != nil {
			t.Fatal("DROP EXTENSION CASCADE left extension-created schema alive")
		}
		for _, stmt := range []string{
			"CREATE EXTENSION test_schema_ext",
			"DROP EXTENSION test_schema_ext",
		} {
			if code := simulationCatalogCode(t, c96, stmt); code != "0" {
				t.Fatalf("reinstall/drop test_schema_ext %s: %s", stmt, code)
			}
		}
		if c96.schemaByName["ext_aux_ns"] != nil {
			t.Fatal("DROP EXTENSION left extension-created schema alive")
		}

		// CREATE OPERATOR only resolves unqualified implementation functions on search_path.
		c97 := New()
		for _, stmt := range []string{
			"CREATE SCHEMA hidden_op_ns",
			"CREATE FUNCTION hidden_op_ns.hid_add(int, int) RETURNS int LANGUAGE sql AS 'SELECT $1 + $2'",
		} {
			if code := simulationCatalogCode(t, c97, stmt); code != "0" {
				t.Fatalf("hidden operator fn setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c97, "CREATE OPERATOR public.## (LEFTARG = int, RIGHTARG = int, FUNCTION = hid_add)"); code != CodeUndefinedFunction {
			t.Fatalf("CREATE OPERATOR with hidden function answered %s, want %s", code, CodeUndefinedFunction)
		}

		// Partitions inherit their parent's explicit table access method when USING is omitted.
		RegisterExtensionSQL("test_part_tam_ext", "CREATE ACCESS METHOD part_tam TYPE TABLE HANDLER heap_tableam_handler;")
		defer delete(extensionScripts, "test_part_tam_ext")
		c98 := New()
		for _, stmt := range []string{
			"CREATE EXTENSION test_part_tam_ext",
			"CREATE TABLE p_tam (a int) PARTITION BY LIST (a) USING part_tam",
			"RESET default_table_access_method",
			"CREATE TABLE c_tam PARTITION OF p_tam FOR VALUES IN (1)",
			"ALTER TABLE p_tam DETACH PARTITION c_tam",
			"DROP TABLE p_tam",
		} {
			if code := simulationCatalogCode(t, c98, stmt); code != "0" {
				t.Fatalf("partition table AM setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c98, "DROP EXTENSION test_part_tam_ext RESTRICT"); code != CodeDependentObjects {
			t.Fatalf("DROP EXTENSION RESTRICT with detached partition using inherited AM answered %s, want %s", code, CodeDependentObjects)
		}

		// Failed CREATE TABLE PARTITION BY with an extension table access method
		// leaves neither a relation nor a dependency on the access method behind.
		RegisterExtensionSQL("test_bad_part_tam_ext", "CREATE ACCESS METHOD bad_part_tam TYPE TABLE HANDLER heap_tableam_handler;")
		defer delete(extensionScripts, "test_bad_part_tam_ext")
		c99 := New()
		if code := simulationCatalogCode(t, c99, "CREATE EXTENSION test_bad_part_tam_ext"); code != "0" {
			t.Fatalf("CREATE EXTENSION test_bad_part_tam_ext: %s", code)
		}
		if code := simulationCatalogCode(t, c99, "CREATE TABLE bad_p (id serial, a int, b int) PARTITION BY LIST (a, b) USING bad_part_tam"); code != CodeInvalidObjectDefinition {
			t.Fatalf("invalid multi-column LIST partition answered %s, want %s", code, CodeInvalidObjectDefinition)
		}
		if c99.schemaByName["public"].Relations["bad_p"] != nil || len(c99.schemaByName["public"].Sequences) != 0 {
			t.Fatal("failed CREATE TABLE PARTITION BY leaked relation or serial sequence")
		}
		if code := simulationCatalogCode(t, c99, "DROP EXTENSION test_bad_part_tam_ext RESTRICT"); code != "0" {
			t.Fatalf("DROP EXTENSION RESTRICT after rolled-back CREATE TABLE answered %s, want 0", code)
		}

		// Failed ALTER TABLE ADD COLUMN with an invalid inline CHECK rolls back
		// both the added column and any identity sequence it created.
		c100 := New()
		if code := simulationCatalogCode(t, c100, "CREATE TABLE add_rb (id int)"); code != "0" {
			t.Fatalf("CREATE TABLE add_rb: %s", code)
		}
		if code := simulationCatalogCode(t, c100, "ALTER TABLE add_rb ADD COLUMN bad int GENERATED ALWAYS AS IDENTITY CHECK (bad)"); code != CodeDatatypeMismatch {
			t.Fatalf("ALTER TABLE ADD COLUMN with non-boolean CHECK answered %s, want %s", code, CodeDatatypeMismatch)
		}
		if rel := c100.schemaByName["public"].Relations["add_rb"]; rel == nil || len(rel.Columns) != 1 {
			t.Fatalf("failed ADD COLUMN left extra column on add_rb: %+v", rel)
		}
		if len(c100.schemaByName["public"].Sequences) != 0 {
			t.Fatalf("failed ADD COLUMN leaked identity sequence: %+v", c100.schemaByName["public"].Sequences)
		}

		// Index and UNIQUE INCLUDE naming deduplicates repeated column names
		// across key and INCLUDE lists, and REPLICA IDENTITY FULL/NOTHING/DEFAULT
		// round-trips through Diff.
		c101 := New()
		for _, stmt := range []string{
			"CREATE TABLE inc_t (a int, b int)",
			"CREATE INDEX ON inc_t (a, a)",
			"CREATE INDEX ON inc_t (a) INCLUDE (a, a)",
			"ALTER TABLE inc_t ADD UNIQUE (a) INCLUDE (a)",
		} {
			if code := simulationCatalogCode(t, c101, stmt); code != "0" {
				t.Fatalf("inc_t setup %s: %s", stmt, code)
			}
		}
		for _, wantIdx := range []string{"inc_t_a_a1_idx", "inc_t_a_a1_a2_idx", "inc_t_a_a1_key"} {
			if c101.schemaByName["public"].Indexes[wantIdx] == nil {
				t.Fatalf("expected auto-named index %q in schema, got %v", wantIdx, c101.schemaByName["public"].Indexes)
			}
		}

		// CREATE TABLE ... AS EXECUTE creates a relation from the prepared SELECT
		// query, and failed multi-subcommand ALTER TABLE restores REPLICA IDENTITY.
		c102 := New()
		for _, stmt := range []string{
			"PREPARE p_ctas AS SELECT 1 AS id, 'ok'::text AS label",
			"CREATE TABLE t_exec AS EXECUTE p_ctas",
		} {
			if code := simulationCatalogCode(t, c102, stmt); code != "0" {
				t.Fatalf("CTAS EXECUTE %s: %s", stmt, code)
			}
		}
		tExec := c102.schemaByName["public"].Relations["t_exec"]
		if tExec == nil || len(tExec.Columns) != 2 || tExec.Columns[0].Name != "id" || tExec.Columns[1].Name != "label" {
			t.Fatalf("unexpected columns on t_exec: %+v", tExec)
		}
		if code := simulationCatalogCode(t, c102, "ALTER TABLE t_exec REPLICA IDENTITY FULL, DROP COLUMN missing_col"); code != CodeUndefinedColumn {
			t.Fatalf("ALTER TABLE REPLICA IDENTITY + bad DROP answered %s, want %s", code, CodeUndefinedColumn)
		}
		if tExec.ReplicaIdentity == 'f' {
			t.Fatal("failed ALTER TABLE left ReplicaIdentity='f' mutated")
		}

		// Review round 1 regressions:
		// 1. EXCLUDE ... INCLUDE (...) migration rendering and round-trip.
		// 2. DISCARD PLANS retains prepared statements; DISCARD ALL drops them.
		// 3. Top-level whole-row Var deparses with ::typename cast so reparsing
		//    preserves the single composite column.
		// 4. Failed multi-action ALTER TABLE restores both ReplicaIdentity and
		//    index IsReplicaIdent flags.
		// 5. WITH CHECK OPTION over a non-updatable view checks only the outer
		//    query locally at definition time.
		c103Base := New()
		c103To := New()
		if code := simulationCatalogCode(t, c103Base, "CREATE TABLE ex_inc (a int4range, b text)"); code != "0" {
			t.Fatalf("c103Base setup: %s", code)
		}
		if code := simulationCatalogCode(t, c103To, "CREATE TABLE ex_inc (a int4range, b text, CONSTRAINT ex_inc_c EXCLUDE USING gist (a WITH &&) INCLUDE (b))"); code != "0" {
			t.Fatalf("c103To setup: %s", code)
		}
		migPlan := GenerateMigration(c103Base, c103To, Diff(c103Base, c103To))
		if !strings.Contains(migPlan.SQL(), "INCLUDE (b)") {
			t.Fatalf("expected EXCLUDE migration to emit INCLUDE (b), got:\n%s", migPlan.SQL())
		}
		if code := simulationCatalogCode(t, c103Base, migPlan.SQL()); code != "0" {
			t.Fatalf("applying EXCLUDE INCLUDE migration failed with %s:\n%s", code, migPlan.SQL())
		}
		if rem := GenerateMigration(c103Base, c103To, Diff(c103Base, c103To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining migration ops after applying EXCLUDE INCLUDE, got:\n%s", rem.SQL())
		}

		c103 := New()
		for _, stmt := range []string{
			"PREPARE p_keep AS SELECT 42 AS v",
			"DISCARD PLANS",
			"CREATE TABLE t_after_discard_plans AS EXECUTE p_keep",
			"DISCARD ALL",
		} {
			if code := simulationCatalogCode(t, c103, stmt); code != "0" {
				t.Fatalf("DISCARD PLANS / ALL %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c103, "CREATE TABLE t_after_discard_all AS EXECUTE p_keep"); code != CodeInvalidSQLStatementName {
			t.Fatalf("CTAS EXECUTE after DISCARD ALL answered %s, want %s", code, CodeInvalidSQLStatementName)
		}

		for _, stmt := range []string{
			"CREATE TABLE wr_t (a int NOT NULL, b text)",
			"CREATE UNIQUE INDEX wr_t_idx ON wr_t (a)",
			"ALTER TABLE wr_t REPLICA IDENTITY USING INDEX wr_t_idx",
			"CREATE VIEW wr_v AS SELECT wr_t, (wr_t IS NOT NULL) AS j FROM wr_t",
			"CREATE VIEW v_agg AS SELECT count(*) AS c FROM wr_t",
			"CREATE VIEW v_outer_casc AS SELECT c FROM v_agg WITH CHECK OPTION",
			"CREATE VIEW v_outer_loc AS SELECT c FROM v_agg WITH LOCAL CHECK OPTION",
		} {
			if code := simulationCatalogCode(t, c103, stmt); code != "0" {
				t.Fatalf("c103 view/replica setup %s: %s", stmt, code)
			}
		}
		wrDef, err := c103.GetViewDefinition("", "wr_v")
		if err != nil {
			t.Fatalf("GetViewDefinition(wr_v): %v", err)
		}
		if !strings.Contains(wrDef, "wr_t.*::wr_t AS wr_t") || !strings.Contains(wrDef, "wr_t.* IS NOT NULL") {
			t.Fatalf("unexpected whole-row view definition:\n%s", wrDef)
		}
		if code := simulationCatalogCode(t, c103, "CREATE VIEW wr_v_reparse AS "+wrDef); code != "0" {
			t.Fatalf("reparsing whole-row viewdef failed with %s:\n%s", code, wrDef)
		}
		if reparsed := c103.schemaByName["public"].Relations["wr_v_reparse"]; reparsed == nil || len(reparsed.Columns) != 2 {
			t.Fatalf("reparsed whole-row view should keep 2 columns, got %+v", reparsed)
		}
		if code := simulationCatalogCode(t, c103, "ALTER TABLE wr_t REPLICA IDENTITY FULL, ALTER COLUMN missing SET STATISTICS 10"); code != CodeUndefinedColumn {
			t.Fatalf("ALTER TABLE REPLICA IDENTITY + bad SET STATISTICS answered %s, want %s", code, CodeUndefinedColumn)
		}
		wrRel := c103.schemaByName["public"].Relations["wr_t"]
		wrIdx := c103.schemaByName["public"].Indexes["wr_t_idx"]
		if wrRel.ReplicaIdentity != 'i' || !wrIdx.IsReplicaIdent {
			t.Fatalf("failed ALTER TABLE did not restore index replica identity: rel=%q idx=%v", wrRel.ReplicaIdentity, wrIdx.IsReplicaIdent)
		}

		// Review round 2 regressions:
		// 1. PREPARE parameter types and EXECUTE argument types propagate to CTAS columns.
		// 2. Multi-action ALTER TABLE rolls back AT_SetStorage (and other column fields).
		// 3. Failed REFRESH MATERIALIZED VIEW CONCURRENTLY leaves IsPopulated=false.
		// 4. DROP COLLATION RESTRICT / CASCADE tracks column and index dependencies.
		c104 := New()
		for _, stmt := range []string{
			"PREPARE q_typed(int4) AS SELECT $1 AS x",
			"CREATE TABLE t_q_typed AS EXECUTE q_typed(1)",
			"PREPARE q_inferred AS SELECT $1 AS y",
			"CREATE TABLE t_q_inferred AS EXECUTE q_inferred('hello'::text)",
		} {
			if code := simulationCatalogCode(t, c104, stmt); code != "0" {
				t.Fatalf("CTAS prepared param %s: %s", stmt, code)
			}
		}
		if r := c104.schemaByName["public"].Relations["t_q_typed"]; r == nil || len(r.Columns) != 1 || r.Columns[0].TypeOID != INT4OID {
			t.Fatalf("expected t_q_typed.x to have INT4OID, got %+v", r)
		}
		if r := c104.schemaByName["public"].Relations["t_q_inferred"]; r == nil || len(r.Columns) != 1 || r.Columns[0].TypeOID != TEXTOID {
			t.Fatalf("expected t_q_inferred.y to have TEXTOID, got %+v", r)
		}

		if code := simulationCatalogCode(t, c104, "CREATE TABLE stor_t (a text)"); code != "0" {
			t.Fatalf("CREATE TABLE stor_t: %s", code)
		}
		if code := simulationCatalogCode(t, c104, "ALTER TABLE stor_t ALTER COLUMN a SET STORAGE PLAIN, ALTER COLUMN missing SET STATISTICS 10"); code != CodeUndefinedColumn {
			t.Fatalf("ALTER TABLE SET STORAGE + bad SET STATISTICS answered %s, want %s", code, CodeUndefinedColumn)
		}
		if storCol := c104.schemaByName["public"].Relations["stor_t"].Columns[0]; storCol.Storage != 'x' {
			t.Fatalf("failed ALTER TABLE left column storage=%q, want 'x'", storCol.Storage)
		}

		if code := simulationCatalogCode(t, c104, "CREATE MATERIALIZED VIEW mv_unpop_chk AS SELECT 1 AS id WITH NO DATA"); code != "0" {
			t.Fatalf("CREATE MATERIALIZED VIEW WITH NO DATA: %s", code)
		}
		if code := simulationCatalogCode(t, c104, "REFRESH MATERIALIZED VIEW CONCURRENTLY mv_unpop_chk"); code != CodeFeatureNotSupported {
			t.Fatalf("REFRESH CONCURRENTLY unpopulated answered %s, want %s", code, CodeFeatureNotSupported)
		}
		if c104.schemaByName["public"].Relations["mv_unpop_chk"].IsPopulated {
			t.Fatal("failed REFRESH CONCURRENTLY marked unpopulated matview as populated")
		}

		for _, stmt := range []string{
			"CREATE COLLATION c_use (LOCALE = 'C')",
			"CREATE TABLE t_c_use (id int, val text COLLATE c_use, plain text)",
			"CREATE INDEX idx_c_use ON t_c_use (plain COLLATE c_use)",
		} {
			if code := simulationCatalogCode(t, c104, stmt); code != "0" {
				t.Fatalf("collation dep setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c104, "DROP COLLATION c_use"); code != CodeDependentObjects {
			t.Fatalf("DROP COLLATION RESTRICT answered %s, want %s", code, CodeDependentObjects)
		}
		if c104.schemaByName["public"].Collations["c_use"] == nil {
			t.Fatal("failed DROP COLLATION removed c_use")
		}
		if code := simulationCatalogCode(t, c104, "DROP COLLATION c_use CASCADE"); code != "0" {
			t.Fatalf("DROP COLLATION CASCADE: %s", code)
		}
		if len(c104.schemaByName["public"].Relations["t_c_use"].Columns) != 2 || c104.schemaByName["public"].Indexes["idx_c_use"] != nil {
			t.Fatal("DROP COLLATION CASCADE did not drop dependent column and index")
		}

		// Review round 3 regressions:
		// 1. Duplicate PREPARE fails with 42P05 and preserves the original plan.
		// 2. Inferred-param PREPARE enforces parameter count on EXECUTE.
		// 3. Multi-action ALTER TABLE restores OfTypeOID on failure.
		// 4. ALTER COLUMN ... TYPE omitting COLLATE clears CollationName while
		//    retaining any DEFAULT-expression collation dependency.
		// 5. DROP COLLATION CASCADE on a DEFAULT expression drops only the default.
		// 6. Auto-named index does not skip names held by standalone enum/domain types.
		c105 := New()
		if code := simulationCatalogCode(t, c105, "PREPARE q_keep AS SELECT 1 AS a"); code != "0" {
			t.Fatalf("PREPARE q_keep: %s", code)
		}
		if code := simulationCatalogCode(t, c105, "PREPARE q_keep AS SELECT 2 AS b"); code != CodeDuplicatePstatement {
			t.Fatalf("duplicate PREPARE answered %s, want %s", code, CodeDuplicatePstatement)
		}
		if code := simulationCatalogCode(t, c105, "CREATE TABLE t_keep AS EXECUTE q_keep"); code != "0" {
			t.Fatalf("CTAS EXECUTE q_keep: %s", code)
		}
		if r := c105.schemaByName["public"].Relations["t_keep"]; r == nil || len(r.Columns) != 1 || r.Columns[0].Name != "a" {
			t.Fatalf("expected original prepared query column 'a', got %+v", r)
		}

		for _, stmt := range []string{
			"CREATE TYPE ty_of AS (a int)",
			"CREATE TABLE t_of (a int)",
		} {
			if code := simulationCatalogCode(t, c105, stmt); code != "0" {
				t.Fatalf("t_of setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c105, "ALTER TABLE t_of OF ty_of, ALTER COLUMN missing SET STATISTICS 1"); code != CodeUndefinedColumn {
			t.Fatalf("ALTER TABLE OF + bad SET STATISTICS answered %s, want %s", code, CodeUndefinedColumn)
		}
		if c105.schemaByName["public"].Relations["t_of"].OfTypeOID != 0 {
			t.Fatal("failed ALTER TABLE left OfTypeOID non-zero")
		}

		for _, stmt := range []string{
			"CREATE COLLATION c_col (LOCALE = 'C')",
			"CREATE COLLATION c_dfl (LOCALE = 'C')",
			"CREATE TABLE t_col_dfl (a text COLLATE c_col DEFAULT ('x' COLLATE c_dfl))",
			"ALTER TABLE t_col_dfl ALTER COLUMN a TYPE varchar",
			"DROP COLLATION c_col",
		} {
			if code := simulationCatalogCode(t, c105, stmt); code != "0" {
				t.Fatalf("c105 collation/default %s: %s", stmt, code)
			}
		}
		colA := c105.schemaByName["public"].Relations["t_col_dfl"].Columns[0]
		if colA.CollationName != "" {
			t.Fatalf("expected ALTER COLUMN TYPE without COLLATE to reset CollationName, got %q", colA.CollationName)
		}
		if code := simulationCatalogCode(t, c105, "DROP COLLATION c_dfl"); code != CodeDependentObjects {
			t.Fatalf("DROP COLLATION c_dfl restricted by default answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c105, "DROP COLLATION c_dfl CASCADE"); code != "0" {
			t.Fatalf("DROP COLLATION c_dfl CASCADE: %s", code)
		}
		if len(c105.schemaByName["public"].Relations["t_col_dfl"].Columns) != 1 || colA.HasDefault {
			t.Fatalf("expected DROP COLLATION CASCADE to drop only default while keeping column a: cols=%d hasDefault=%v",
				len(c105.schemaByName["public"].Relations["t_col_dfl"].Columns), colA.HasDefault)
		}

		for _, stmt := range []string{
			"CREATE TYPE t_enum_idx_a_idx AS ENUM ('x')",
			"CREATE TABLE t_enum_idx (a int)",
			"CREATE INDEX ON t_enum_idx (a)",
		} {
			if code := simulationCatalogCode(t, c105, stmt); code != "0" {
				t.Fatalf("auto-named index vs enum setup %s: %s", stmt, code)
			}
		}
		if c105.schemaByName["public"].Indexes["t_enum_idx_a_idx"] == nil {
			t.Fatalf("expected auto-named index t_enum_idx_a_idx, got %v", c105.schemaByName["public"].Indexes)
		}

		// Review round 4 regressions:
		// 1. Structural multi-action ALTER TABLE rollback (ADD COLUMN, DROP COLUMN,
		//    ADD CONSTRAINT) restores columns, colByName, and constraints on failure.
		// 2. PREPARE infers parameter types at PREPARE time and defaults unconstrained
		//    parameters to text; EXECUTE assignment-coerces arguments to inferred types.
		// 3. Constraint-backed index naming ignores standalone enum/domain types.
		// 4. Explicit column STORAGE and COMPRESSION diff and emit in migrations.
		// 5. User collations diff and emit in GenerateMigration in dependency order.
		// 6. Renamed index columns (ALTER INDEX ... RENAME COLUMN) diff and migrate.
		// 7. DEALLOCATE on a missing prepared statement fails with 26000.
		c106 := New()
		if code := simulationCatalogCode(t, c106, "CREATE TABLE t_rb_struct (a int, b int)"); code != "0" {
			t.Fatalf("CREATE TABLE t_rb_struct: %s", code)
		}
		if code := simulationCatalogCode(t, c106, "ALTER TABLE t_rb_struct ADD COLUMN c int, ALTER COLUMN missing SET STATISTICS 10"); code != CodeUndefinedColumn {
			t.Fatalf("ALTER TABLE ADD COLUMN + bad SET STATISTICS answered %s, want %s", code, CodeUndefinedColumn)
		}
		rbRel := c106.schemaByName["public"].Relations["t_rb_struct"]
		if _, hasC := rbRel.colByName["c"]; len(rbRel.Columns) != 2 || hasC {
			t.Fatalf("failed ALTER TABLE ADD COLUMN leaked column c: cols=%d hasC=%v", len(rbRel.Columns), hasC)
		}
		if code := simulationCatalogCode(t, c106, "ALTER TABLE t_rb_struct DROP COLUMN b, ADD UNIQUE (a), ALTER COLUMN missing SET STATISTICS 10"); code != CodeUndefinedColumn {
			t.Fatalf("ALTER TABLE DROP COLUMN + ADD UNIQUE + bad SET STATISTICS answered %s, want %s", code, CodeUndefinedColumn)
		}
		if bIdx, hasB := rbRel.colByName["b"]; len(rbRel.Columns) != 2 || !hasB || bIdx != 1 || len(c106.consByRel[rbRel.OID]) != 0 || len(c106.indexesByRel[rbRel.OID]) != 0 {
			t.Fatalf("failed ALTER TABLE DROP COLUMN/ADD UNIQUE did not restore structural state: cols=%d hasB=%v bIdx=%d cons=%d idxs=%d",
				len(rbRel.Columns), hasB, bIdx, len(c106.consByRel[rbRel.OID]), len(c106.indexesByRel[rbRel.OID]))
		}

		for _, stmt := range []string{
			"PREPARE q_r4_add AS SELECT $1 + 1 AS x",
			"CREATE TABLE t_r4_add AS EXECUTE q_r4_add(1::bigint)",
			"PREPARE q_r4_bare AS SELECT $1 AS x",
			"CREATE TABLE t_r4_bare AS EXECUTE q_r4_bare(1)",
		} {
			if code := simulationCatalogCode(t, c106, stmt); code != "0" {
				t.Fatalf("c106 PREPARE/EXECUTE %s: %s", stmt, code)
			}
		}
		if r := c106.schemaByName["public"].Relations["t_r4_add"]; r == nil || r.Columns[0].TypeOID != INT4OID {
			t.Fatalf("expected t_r4_add.x to have INT4OID (inferred $1 + 1), got %+v", r)
		}
		if r := c106.schemaByName["public"].Relations["t_r4_bare"]; r == nil || r.Columns[0].TypeOID != TEXTOID {
			t.Fatalf("expected t_r4_bare.x to have TEXTOID (unconstrained $1 defaults to text), got %+v", r)
		}
		if code := simulationCatalogCode(t, c106, "DEALLOCATE missing_r4_stmt"); code != CodeInvalidSQLStatementName {
			t.Fatalf("DEALLOCATE missing_r4_stmt answered %s, want %s", code, CodeInvalidSQLStatementName)
		}

		for _, stmt := range []string{
			"CREATE TYPE t_u_r4_a_key AS ENUM ('x')",
			"CREATE TABLE t_u_r4 (a int UNIQUE)",
		} {
			if code := simulationCatalogCode(t, c106, stmt); code != "0" {
				t.Fatalf("c106 constraint vs enum %s: %s", stmt, code)
			}
		}
		if c106.schemaByName["public"].Indexes["t_u_r4_a_key"] == nil {
			t.Fatalf("expected UNIQUE constraint index t_u_r4_a_key, got %v", c106.schemaByName["public"].Indexes)
		}

		c106From := New()
		c106To := New()
		for _, stmt := range []string{
			"CREATE TABLE t_sc_mod (a text, b text STORAGE EXTERNAL COMPRESSION lz4)",
			"CREATE TABLE t_idx_ren (a int)",
			"CREATE INDEX idx_ren ON t_idx_ren ((a + 1))",
		} {
			if code := simulationCatalogCode(t, c106From, stmt); code != "0" {
				t.Fatalf("c106From setup %s: %s", stmt, code)
			}
		}
		for _, stmt := range []string{
			"CREATE COLLATION coll_r4 (LOCALE = 'C')",
			"CREATE TABLE t_sc_new (a text COLLATE coll_r4 STORAGE EXTERNAL COMPRESSION lz4 NOT NULL)",
			"CREATE TABLE t_sc_mod (a text STORAGE MAIN COMPRESSION pglz, b text)",
			"CREATE TABLE t_idx_ren (a int)",
			"CREATE INDEX idx_ren ON t_idx_ren ((a + 1))",
			"ALTER TABLE idx_ren RENAME COLUMN expr TO plus_one",
		} {
			if code := simulationCatalogCode(t, c106To, stmt); code != "0" {
				t.Fatalf("c106To setup %s: %s", stmt, code)
			}
		}
		mig106 := GenerateMigration(c106From, c106To, Diff(c106From, c106To))
		migSQL := mig106.SQL()
		for _, wantFrag := range []string{
			"CREATE COLLATION",
			"STORAGE EXTERNAL",
			"COMPRESSION lz4",
			"SET STORAGE MAIN",
			"SET COMPRESSION pglz",
			"SET STORAGE DEFAULT",
			"SET COMPRESSION DEFAULT",
			"RENAME COLUMN",
		} {
			if !strings.Contains(migSQL, wantFrag) {
				t.Fatalf("expected migration SQL to contain %q, got:\n%s", wantFrag, migSQL)
			}
		}
		if code := simulationCatalogCode(t, c106From, migSQL); code != "0" {
			t.Fatalf("applying c106 migration failed with %s:\n%s", code, migSQL)
		}
		if rem := GenerateMigration(c106From, c106To, Diff(c106From, c106To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining migration ops after c106 round-trip, got:\n%s", rem.SQL())
		}

		// Review round 5 regressions:
		// 1. Dropping a DEFAULT that references the same user collation as the
		//    column's explicit COLLATE leaves exactly one collation dependency so
		//    subsequent ALTER COLUMN TYPE (without COLLATE) allows DROP COLLATION.
		// 2. Replacing a collation whose definition changed while retained columns,
		//    defaults, indexes, constraints, domains, and views reference it
		//    detaches/rebuilds dependents before DROP COLLATION and restores them
		//    after CREATE COLLATION.
		c107 := New()
		for _, stmt := range []string{
			"CREATE COLLATION c_same_r5 (LOCALE = 'C')",
			"CREATE DOMAIN d_keep_r5 AS text",
			"CREATE TABLE t_same_r5 (a text COLLATE c_same_r5 DEFAULT ('x' COLLATE c_same_r5), b d_keep_r5 DEFAULT ('y'::d_keep_r5))",
			"ALTER TABLE t_same_r5 ALTER COLUMN a DROP DEFAULT",
			"ALTER TABLE t_same_r5 ALTER COLUMN b DROP DEFAULT",
			"ALTER TABLE t_same_r5 ALTER COLUMN a TYPE varchar",
			"DROP COLLATION c_same_r5",
		} {
			if code := simulationCatalogCode(t, c107, stmt); code != "0" {
				t.Fatalf("c107 %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c107, "DROP DOMAIN d_keep_r5"); code != CodeDependentObjects {
			t.Fatalf("DROP DOMAIN d_keep_r5 after DROP DEFAULT answered %s, want %s", code, CodeDependentObjects)
		}

		c107From := New()
		c107To := New()
		for _, stmt := range []string{
			"CREATE COLLATION coll_repl (LOCALE = 'C')",
			"CREATE DOMAIN dom_repl AS text COLLATE coll_repl",
			"CREATE TABLE t_repl (a text COLLATE coll_repl STORAGE EXTERNAL DEFAULT ('d' COLLATE coll_repl), b text, d dom_repl, CONSTRAINT t_repl_chk CHECK (b COLLATE coll_repl <> ''))",
			"CREATE INDEX idx_repl ON t_repl (b COLLATE coll_repl)",
			"CREATE VIEW v_repl AS SELECT a, b FROM t_repl",
		} {
			if code := simulationCatalogCode(t, c107From, stmt); code != "0" {
				t.Fatalf("c107From setup %s: %s", stmt, code)
			}
		}
		for _, stmt := range []string{
			"CREATE COLLATION coll_repl (LOCALE = 'POSIX')",
			"CREATE DOMAIN dom_repl AS text COLLATE coll_repl",
			"CREATE TABLE t_repl (a text COLLATE coll_repl STORAGE EXTERNAL DEFAULT ('d' COLLATE coll_repl), b text, d dom_repl, CONSTRAINT t_repl_chk CHECK (b COLLATE coll_repl <> ''))",
			"CREATE INDEX idx_repl ON t_repl (b COLLATE coll_repl)",
			"CREATE VIEW v_repl AS SELECT a, b FROM t_repl",
		} {
			if code := simulationCatalogCode(t, c107To, stmt); code != "0" {
				t.Fatalf("c107To setup %s: %s", stmt, code)
			}
		}
		mig107 := GenerateMigration(c107From, c107To, Diff(c107From, c107To))
		if code := simulationCatalogCode(t, c107From, mig107.SQL()); code != "0" {
			t.Fatalf("applying c107 collation replacement migration failed with %s:\n%s", code, mig107.SQL())
		}
		if rem := GenerateMigration(c107From, c107To, Diff(c107From, c107To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining migration ops after c107 collation replacement, got:\n%s", rem.SQL())
		}

		// Review round 6 regressions:
		// 1. Index attribute rename cycles ((a, b) -> (b, a)) unwind without
		//    duplicate-column errors.
		// 2. CREATE COLLATION FROM missing fails with 42704 and leaves no
		//    registered collation; valid CREATE COLLATION FROM copies cleanly.
		// 3. ALTER TABLE ... SET ACCESS METHOD diffs and round-trips through
		//    GenerateMigration.
		c108 := New()
		if code := simulationCatalogCode(t, c108, "CREATE COLLATION c_bad_r6 FROM missing_r6"); code != CodeUndefinedObject {
			t.Fatalf("CREATE COLLATION FROM missing_r6 answered %s, want %s", code, CodeUndefinedObject)
		}
		if code := simulationCatalogCode(t, c108, "CREATE TABLE t_bad_r6 (a text COLLATE c_bad_r6)"); code != CodeUndefinedObject {
			t.Fatalf("column referencing uncreated collation answered %s, want %s", code, CodeUndefinedObject)
		}
		for _, stmt := range []string{
			"CREATE COLLATION c_src_r6 (LOCALE = 'C')",
			"CREATE COLLATION c_copy_r6 FROM c_src_r6",
			"DROP COLLATION c_src_r6",
			"CREATE TABLE t_copy_r6 (a text COLLATE c_copy_r6)",
		} {
			if code := simulationCatalogCode(t, c108, stmt); code != "0" {
				t.Fatalf("c108 valid collation copy %s: %s", stmt, code)
			}
		}

		c108From := New()
		c108To := New()
		tamSetup := []string{
			"CREATE FUNCTION tam_handler_r6(internal) RETURNS table_am_handler LANGUAGE internal AS 'heap_tableam_handler'",
			"CREATE ACCESS METHOD custom_tam_r6 TYPE TABLE HANDLER tam_handler_r6",
		}
		for _, stmt := range tamSetup {
			if code := simulationCatalogCode(t, c108From, stmt); code != "0" {
				t.Fatalf("c108From tam setup %s: %s", stmt, code)
			}
			if code := simulationCatalogCode(t, c108To, stmt); code != "0" {
				t.Fatalf("c108To tam setup %s: %s", stmt, code)
			}
		}
		for _, stmt := range []string{
			"CREATE TABLE t_cyc_r6 (x int, y int)",
			"CREATE INDEX idx_cyc_r6 ON t_cyc_r6 ((x + 1), (y + 1))",
			"ALTER TABLE idx_cyc_r6 RENAME COLUMN expr TO a",
			"ALTER TABLE idx_cyc_r6 RENAME COLUMN expr1 TO b",
			"CREATE TABLE t_am_r6 (id int)",
		} {
			if code := simulationCatalogCode(t, c108From, stmt); code != "0" {
				t.Fatalf("c108From setup %s: %s", stmt, code)
			}
		}
		for _, stmt := range []string{
			"CREATE TABLE t_cyc_r6 (x int, y int)",
			"CREATE INDEX idx_cyc_r6 ON t_cyc_r6 ((x + 1), (y + 1))",
			"ALTER TABLE idx_cyc_r6 RENAME COLUMN expr TO b",
			"ALTER TABLE idx_cyc_r6 RENAME COLUMN expr1 TO a",
			"CREATE TABLE t_am_r6 (id int)",
			"ALTER TABLE t_am_r6 SET ACCESS METHOD custom_tam_r6",
		} {
			if code := simulationCatalogCode(t, c108To, stmt); code != "0" {
				t.Fatalf("c108To setup %s: %s", stmt, code)
			}
		}
		mig108 := GenerateMigration(c108From, c108To, Diff(c108From, c108To))
		if !strings.Contains(mig108.SQL(), "SET ACCESS METHOD") {
			t.Fatalf("expected migration SQL to contain SET ACCESS METHOD, got:\n%s", mig108.SQL())
		}
		if code := simulationCatalogCode(t, c108From, mig108.SQL()); code != "0" {
			t.Fatalf("applying c108 migration failed with %s:\n%s", code, mig108.SQL())
		}
		if rem := GenerateMigration(c108From, c108To, Diff(c108From, c108To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining migration ops after c108 round-trip, got:\n%s", rem.SQL())
		}

		// Review round 7 regressions:
		// 1. Unique index on partitioned table putting partition key only in INCLUDE is rejected (0A000).
		// 2. Index COLLATE with quoted dots in schema/name round-trips without splitting on literal dots.
		// 3. Domain collation change replaces the domain safely around dependent columns, constraints, and views.
		// 4. Range type COLLATION is retained in RangeType, diffed, and emitted in GenerateMigration.
		// 5. Stored generated column rebuilt for a modified collation preserves unchanged indexes, constraints, inbound FKs, triggers, policies, and views.
		c109Part := New()
		if code := simulationCatalogCode(t, c109Part, "CREATE TABLE p_inc_r7 (k int, v int) PARTITION BY LIST (k)"); code != "0" {
			t.Fatalf("c109Part CREATE TABLE failed: %s", code)
		}
		if code := simulationCatalogCode(t, c109Part, "CREATE UNIQUE INDEX p_inc_r7_bad ON p_inc_r7 (v) INCLUDE (k)"); code != CodeFeatureNotSupported {
			t.Fatalf("unique index with partition key only in INCLUDE answered %s, want %s", code, CodeFeatureNotSupported)
		}
		if code := simulationCatalogCode(t, c109Part, "CREATE UNIQUE INDEX p_inc_r7_ok ON p_inc_r7 (k) INCLUDE (v)"); code != "0" {
			t.Fatalf("unique index with partition key in key columns answered %s, want 0", code)
		}

		c109From := New()
		c109To := New()
		for _, stmt := range []string{
			`CREATE COLLATION "a.b" (LOCALE = 'C')`,
			`CREATE COLLATION coll_gen_r7 (LOCALE = 'C')`,
			`CREATE COLLATION coll_dom_a_r7 (LOCALE = 'C')`,
			`CREATE COLLATION coll_dom_b_r7 (LOCALE = 'POSIX')`,
			`CREATE DOMAIN dom_chg_r7 AS text COLLATE coll_dom_a_r7`,
			`CREATE TABLE t_gen_r7 (id int PRIMARY KEY, src text, g text COLLATE coll_gen_r7 GENERATED ALWAYS AS (lower(src)) STORED, d dom_chg_r7 DEFAULT 'x'::dom_chg_r7, CONSTRAINT t_gen_r7_uk UNIQUE (g), CONSTRAINT t_gen_r7_chk CHECK (g <> '' AND d <> ''))`,
			`CREATE INDEX idx_gen_r7 ON t_gen_r7 (g)`,
			`CREATE INDEX idx_dot_r7 ON t_gen_r7 (src COLLATE public."a.b")`,
			`CREATE TABLE t_fk_r7 (id int PRIMARY KEY, g_ref text REFERENCES t_gen_r7(g))`,
			`CREATE VIEW v_gen_r7 AS SELECT id, g, d FROM t_gen_r7`,
		} {
			if code := simulationCatalogCode(t, c109From, stmt); code != "0" {
				t.Fatalf("c109From setup %s: %s", stmt, code)
			}
		}
		for _, stmt := range []string{
			`CREATE COLLATION "a.b" (LOCALE = 'C')`,
			`CREATE COLLATION coll_gen_r7 (LOCALE = 'POSIX')`,
			`CREATE COLLATION coll_dom_a_r7 (LOCALE = 'C')`,
			`CREATE COLLATION coll_dom_b_r7 (LOCALE = 'POSIX')`,
			`CREATE DOMAIN dom_chg_r7 AS text COLLATE coll_dom_b_r7`,
			`CREATE TYPE rng_r7 AS RANGE (SUBTYPE = text, COLLATION = public."a.b")`,
			`CREATE TABLE t_gen_r7 (id int PRIMARY KEY, src text, g text COLLATE coll_gen_r7 GENERATED ALWAYS AS (lower(src)) STORED, d dom_chg_r7 DEFAULT 'x'::dom_chg_r7, CONSTRAINT t_gen_r7_uk UNIQUE (g), CONSTRAINT t_gen_r7_chk CHECK (g <> '' AND d <> ''))`,
			`CREATE INDEX idx_gen_r7 ON t_gen_r7 (g)`,
			`CREATE INDEX idx_dot_r7 ON t_gen_r7 (src COLLATE public."a.b")`,
			`CREATE TABLE t_fk_r7 (id int PRIMARY KEY, g_ref text REFERENCES t_gen_r7(g))`,
			`CREATE VIEW v_gen_r7 AS SELECT id, g, d FROM t_gen_r7`,
		} {
			if code := simulationCatalogCode(t, c109To, stmt); code != "0" {
				t.Fatalf("c109To setup %s: %s", stmt, code)
			}
		}
		migEmptyDot := GenerateMigration(New(), c109To, Diff(New(), c109To))
		if !strings.Contains(migEmptyDot.SQL(), `COLLATE public."a.b"`) {
			t.Fatalf("expected index DDL to preserve quoted dot in collation name, got:\n%s", migEmptyDot.SQL())
		}
		if !strings.Contains(migEmptyDot.SQL(), `CREATE TYPE "public"."rng_r7" AS RANGE (SUBTYPE = text, COLLATION = "public"."a.b")`) {
			t.Fatalf("expected range DDL to retain COLLATION, got:\n%s", migEmptyDot.SQL())
		}
		c109Empty := New()
		if code := simulationCatalogCode(t, c109Empty, migEmptyDot.SQL()); code != "0" {
			t.Fatalf("applying empty->c109To migration failed with %s:\n%s", code, migEmptyDot.SQL())
		}
		if rem := GenerateMigration(c109Empty, c109To, Diff(c109Empty, c109To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining ops after empty->c109To, got:\n%s", rem.SQL())
		}

		mig109 := GenerateMigration(c109From, c109To, Diff(c109From, c109To))
		if code := simulationCatalogCode(t, c109From, mig109.SQL()); code != "0" {
			t.Fatalf("applying c109From->c109To migration failed with %s:\n%s", code, mig109.SQL())
		}
		if rem := GenerateMigration(c109From, c109To, Diff(c109From, c109To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining ops after c109From->c109To (generated column dependents, domain collation, range collation), got:\n%s", rem.SQL())
		}

		// Review round 8 regressions:
		// 1. Schema-qualified column collations outside search_path round-trip and diff across schemas.
		// 2. DISCARD ALL resets search_path, default_tablespace, default_table_access_method, and current role.
		// 3. CREATE RULE propagates missing-relation (42P01), wrong-relkind (42809/0A000), and duplicate-rule (42710) errors.
		// 4. CREATE COLLATION rejects RULES when provider != icu (42P17).
		// 5. ALTER COLUMN ... TYPE resets Column.Compression to 0.
		c110 := New()
		if code := simulationCatalogCode(t, c110, "CREATE RULE r_missing AS ON INSERT TO missing_r8 DO NOTHING"); code != CodeUndefinedTable {
			t.Fatalf("CREATE RULE on missing table answered %s, want %s", code, CodeUndefinedTable)
		}
		if code := simulationCatalogCode(t, c110, "CREATE SEQUENCE seq_r8; CREATE RULE r_seq AS ON INSERT TO seq_r8 DO NOTHING"); code != CodeWrongObjectType {
			t.Fatalf("CREATE RULE on sequence answered %s, want %s", code, CodeWrongObjectType)
		}
		if code := simulationCatalogCode(t, c110, "CREATE MATERIALIZED VIEW mv_r8 AS SELECT 1 AS a; CREATE RULE r_mv AS ON INSERT TO mv_r8 DO NOTHING"); code != CodeFeatureNotSupported {
			t.Fatalf("CREATE RULE on matview answered %s, want %s", code, CodeFeatureNotSupported)
		}
		if code := simulationCatalogCode(t, c110, "CREATE COLLATION c_bad_rules_r8 (PROVIDER = libc, LOCALE = 'C', RULES = '&a < b')"); code != CodeInvalidObjectDefinition {
			t.Fatalf("CREATE COLLATION with RULES and libc answered %s, want %s", code, CodeInvalidObjectDefinition)
		}
		if code := simulationCatalogCode(t, c110, "CREATE COLLATION c_icu_rules_r8 (PROVIDER = icu, LOCALE = 'en-x-icu', RULES = '&a < b')"); code != "0" {
			t.Fatalf("CREATE COLLATION with RULES and icu answered %s, want 0", code)
		}
		for _, stmt := range []string{
			"CREATE SCHEMA s_disc_r8",
			"CREATE ROLE role_disc_r8",
			"CREATE TABLESPACE ts_disc_r8 LOCATION '/tmp/ts_disc_r8'",
			"CREATE FUNCTION tam_disc_r8(internal) RETURNS table_am_handler LANGUAGE internal AS 'heap_tableam_handler'",
			"CREATE ACCESS METHOD am_disc_r8 TYPE TABLE HANDLER tam_disc_r8",
			"SET search_path = s_disc_r8",
			"SET ROLE role_disc_r8",
			"SET default_tablespace = ts_disc_r8",
			"SET default_table_access_method = am_disc_r8",
			"PREPARE p_disc_r8 AS SELECT 1",
			"DISCARD ALL",
			"CREATE TABLE t_after_disc_r8 (a text COMPRESSION pglz)",
		} {
			if code := simulationCatalogCode(t, c110, stmt); code != "0" {
				t.Fatalf("c110 %s: %s", stmt, code)
			}
		}
		relAfter := c110.GetRelation("public", "t_after_disc_r8")
		if relAfter == nil || relAfter.Tablespace != "" || relAfter.AccessMethod != "" || relAfter.OwnerName != "postgres" {
			t.Fatalf("expected DISCARD ALL to reset schema/tablespace/am/role, got rel=%+v", relAfter)
		}
		if relAfter.Columns[0].Compression != 'p' {
			t.Fatalf("expected initial Compression 'p', got %q", relAfter.Columns[0].Compression)
		}
		if code := simulationCatalogCode(t, c110, "ALTER TABLE t_after_disc_r8 ALTER COLUMN a TYPE varchar"); code != "0" {
			t.Fatalf("ALTER COLUMN TYPE answered %s, want 0", code)
		}
		if relAfter.Columns[0].Compression != 0 {
			t.Fatalf("expected ALTER COLUMN TYPE to reset Compression to 0, got %q", relAfter.Columns[0].Compression)
		}

		c110From := New()
		c110To := New()
		for _, stmt := range []string{
			"CREATE SCHEMA s1_r8",
			"CREATE SCHEMA s2_r8",
			"CREATE COLLATION s1_r8.c_shared (LOCALE = 'C')",
			"CREATE COLLATION s2_r8.c_shared (LOCALE = 'POSIX')",
			"CREATE TABLE public.t_qcoll_r8 (a text COLLATE s1_r8.c_shared)",
		} {
			if code := simulationCatalogCode(t, c110From, stmt); code != "0" {
				t.Fatalf("c110From %s: %s", stmt, code)
			}
		}
		for _, stmt := range []string{
			"CREATE SCHEMA s1_r8",
			"CREATE SCHEMA s2_r8",
			"CREATE COLLATION s1_r8.c_shared (LOCALE = 'C')",
			"CREATE COLLATION s2_r8.c_shared (LOCALE = 'POSIX')",
			"CREATE TABLE public.t_qcoll_r8 (a text COLLATE s2_r8.c_shared)",
		} {
			if code := simulationCatalogCode(t, c110To, stmt); code != "0" {
				t.Fatalf("c110To %s: %s", stmt, code)
			}
		}
		migEmpty110 := GenerateMigration(New(), c110To, Diff(New(), c110To))
		c110Empty := New()
		if code := simulationCatalogCode(t, c110Empty, migEmpty110.SQL()); code != "0" {
			t.Fatalf("applying empty->c110To failed with %s:\n%s", code, migEmpty110.SQL())
		}
		if rem := GenerateMigration(c110Empty, c110To, Diff(c110Empty, c110To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining ops after empty->c110To, got:\n%s", rem.SQL())
		}
		mig110 := GenerateMigration(c110From, c110To, Diff(c110From, c110To))
		if len(mig110.Ops) == 0 {
			t.Fatal("expected migration ops when changing column collation from s1_r8.c_shared to s2_r8.c_shared")
		}
		if code := simulationCatalogCode(t, c110From, mig110.SQL()); code != "0" {
			t.Fatalf("applying c110From->c110To failed with %s:\n%s", code, mig110.SQL())
		}
		if rem := GenerateMigration(c110From, c110To, Diff(c110From, c110To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining ops after c110From->c110To, got:\n%s", rem.SQL())
		}

		// Review round 9 regressions:
		// 1. Explicit column collations populate Column.Collation so conflicting
		//    implicit collations in views fail with 42P22.
		// 2. RESET ALL preserves the active SET ROLE (unlike DISCARD ALL).
		// 3. Dropped collations referenced only by a retained column's default
		//    expression defer DROP COLLATION to PhaseLateDrop after DROP DEFAULT.
		c111 := New()
		for _, stmt := range []string{
			"CREATE COLLATION c1_r9 (LOCALE = 'C')",
			"CREATE COLLATION c2_r9 (LOCALE = 'POSIX')",
			"CREATE TABLE t_coll_r9 (a text COLLATE c1_r9, b text COLLATE c2_r9)",
		} {
			if code := simulationCatalogCode(t, c111, stmt); code != "0" {
				t.Fatalf("c111 %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c111, "CREATE VIEW v_coll_u_r9 AS SELECT a FROM t_coll_r9 UNION SELECT b FROM t_coll_r9"); code != CodeCollationMismatch {
			t.Fatalf("UNION view over conflicting explicit column collations answered %s, want %s", code, CodeCollationMismatch)
		}
		if code := simulationCatalogCode(t, c111, "CREATE VIEW v_coll_ua_r9 AS SELECT a FROM t_coll_r9 UNION ALL SELECT b FROM t_coll_r9"); code != CodeIndeterminateCollation {
			t.Fatalf("UNION ALL view over conflicting explicit column collations answered %s, want %s", code, CodeIndeterminateCollation)
		}
		if code := simulationCatalogCode(t, c111, "CREATE VIEW v_coll_rep_r9 AS SELECT a AS x FROM t_coll_r9; CREATE OR REPLACE VIEW v_coll_rep_r9 AS SELECT b AS x FROM t_coll_r9"); code != CodeInvalidObjectDefinition {
			t.Fatalf("CREATE OR REPLACE VIEW changing column collation answered %s, want %s", code, CodeInvalidObjectDefinition)
		}
		for _, stmt := range []string{
			"CREATE ROLE app_role_r9",
			"SET ROLE app_role_r9",
			"RESET ALL",
			"CREATE TABLE t_role_r9 (id int)",
		} {
			if code := simulationCatalogCode(t, c111, stmt); code != "0" {
				t.Fatalf("c111 role %s: %s", stmt, code)
			}
		}
		if rel := c111.GetRelation("public", "t_role_r9"); rel == nil || rel.OwnerName != "app_role_r9" {
			t.Fatalf("expected RESET ALL to preserve active SET ROLE app_role_r9, got %+v", rel)
		}

		c111From := New()
		c111To := New()
		for _, stmt := range []string{
			"CREATE COLLATION c_dfl_r9 (LOCALE = 'C')",
			"CREATE TABLE t_dfl_r9 (a text DEFAULT ('x' COLLATE c_dfl_r9))",
		} {
			if code := simulationCatalogCode(t, c111From, stmt); code != "0" {
				t.Fatalf("c111From %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c111To, "CREATE TABLE t_dfl_r9 (a text)"); code != "0" {
			t.Fatalf("c111To CREATE TABLE: %s", code)
		}
		mig111 := GenerateMigration(c111From, c111To, Diff(c111From, c111To))
		if code := simulationCatalogCode(t, c111From, mig111.SQL()); code != "0" {
			t.Fatalf("applying c111From->c111To (dropping default + collation) failed with %s:\n%s", code, mig111.SQL())
		}
		if rem := GenerateMigration(c111From, c111To, Diff(c111From, c111To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining ops after c111From->c111To, got:\n%s", rem.SQL())
		}

		// Review round 10 regressions:
		// 1. Bounded varlena columns below TOAST_TUPLE_THRESHOLD (2032) do not
		//    need a TOAST table (so ALTER TABLE ... SET (toast.*) ignores
		//    toast.* on varchar(1)/varchar(501) and rejects unknown toast.* on
		//    varchar(502)/text/domains), while CREATE TABLE validates toast.*
		//    unconditionally like PostgreSQL's DefineRelation.
		// 2. REPLICA IDENTITY carries ReplicaIdentityStmt in the AST so AST-only
		//    ProcessUtility callers apply FULL, NOTHING, DEFAULT, and USING INDEX.
		// 3. DROP COLLATION ... CASCADE drops inherited and partitioned column
		//    dependents recursively, and rolls back cleanly when a dependent
		//    column (such as a partition key) cannot be dropped.
		// 4. CREATE RULE definitions and firing modes are stored on Relation,
		//    compared in compareRelation and compareView, and emitted in
		//    GenerateMigration (including rule renames, mode changes, and view
		//    rebuilds).
		c112 := New()
		for _, stmt := range []string{
			"CREATE TABLE t_vc1_r10 (v varchar(1))",
			"CREATE TABLE t_vc501_r10 (v varchar(501))",
			"CREATE TABLE t_vc502_r10 (v varchar(502))",
			"CREATE DOMAIN dom_vc1_r10 AS varchar(1)",
			"CREATE TABLE t_dom_vc1_r10 (v dom_vc1_r10)",
		} {
			if code := simulationCatalogCode(t, c112, stmt); code != "0" {
				t.Fatalf("c112 toast setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c112, "ALTER TABLE t_vc1_r10 SET (toast.a = 1)"); code != "0" {
			t.Fatalf("ALTER TABLE varchar(1) SET (toast.a = 1) answered %s, want 0", code)
		}
		if code := simulationCatalogCode(t, c112, "ALTER TABLE t_vc501_r10 SET (toast.a = 1)"); code != "0" {
			t.Fatalf("ALTER TABLE varchar(501) SET (toast.a = 1) answered %s, want 0", code)
		}
		if code := simulationCatalogCode(t, c112, "ALTER TABLE t_vc502_r10 SET (toast.a = 1)"); code != CodeInvalidParameterValue {
			t.Fatalf("ALTER TABLE varchar(502) SET (toast.a = 1) answered %s, want %s", code, CodeInvalidParameterValue)
		}
		if code := simulationCatalogCode(t, c112, "ALTER TABLE t_dom_vc1_r10 SET (toast.a = 1)"); code != CodeInvalidParameterValue {
			t.Fatalf("ALTER TABLE domain(varchar(1)) SET (toast.a = 1) answered %s, want %s", code, CodeInvalidParameterValue)
		}
		if code := simulationCatalogCode(t, c112, "CREATE TABLE t_ct_vc1_r10 (v varchar(1)) WITH (toast.a = 1)"); code != CodeInvalidParameterValue {
			t.Fatalf("CREATE TABLE varchar(1) WITH (toast.a = 1) answered %s, want %s", code, CodeInvalidParameterValue)
		}

		// 2. AST-only ProcessUtility for REPLICA IDENTITY (no Exec / raw SQL).
		if code := simulationCatalogCode(t, c112, "CREATE TABLE t_ri_ast_r10 (id int NOT NULL); CREATE UNIQUE INDEX t_ri_ast_r10_uk ON t_ri_ast_r10 (id)"); code != "0" {
			t.Fatalf("c112 ri setup: %s", code)
		}
		for _, tcRI := range []struct {
			sql  string
			want byte
		}{
			{"ALTER TABLE t_ri_ast_r10 REPLICA IDENTITY FULL", 'f'},
			{"ALTER TABLE t_ri_ast_r10 REPLICA IDENTITY NOTHING", 'n'},
			{"ALTER TABLE t_ri_ast_r10 REPLICA IDENTITY USING INDEX t_ri_ast_r10_uk", 'i'},
			{"ALTER TABLE t_ri_ast_r10 REPLICA IDENTITY DEFAULT", 'd'},
		} {
			parsed, err := pgparser.Parse(tcRI.sql)
			if err != nil || parsed == nil || len(parsed.Items) != 1 {
				t.Fatalf("Parse(%q): %v", tcRI.sql, err)
			}
			stmt := parsed.Items[0]
			if raw, ok := stmt.(*nodes.RawStmt); ok {
				stmt = raw.Stmt
			}
			if err := c112.ProcessUtility(stmt); err != nil {
				t.Fatalf("ProcessUtility(%q): %v", tcRI.sql, err)
			}
			if got := c112.GetRelation("public", "t_ri_ast_r10").ReplicaIdentity; got != tcRI.want {
				t.Fatalf("ProcessUtility(%q): ReplicaIdentity = %q, want %q", tcRI.sql, got, tcRI.want)
			}
		}

		// 3. DROP COLLATION ... CASCADE on inherited, merged, partitioned, and partition-key columns.
		for _, stmt := range []string{
			"CREATE COLLATION c_inh_r10 (LOCALE = 'C')",
			"CREATE TABLE p_inh_r10 (id int, a text COLLATE c_inh_r10, b text COLLATE c_inh_r10)",
			"CREATE TABLE ch_inh_r10 () INHERITS (p_inh_r10)",
			"CREATE TABLE ch_merge_r10 (a text COLLATE c_inh_r10) INHERITS (p_inh_r10)",
			"CREATE TABLE gch_inh_r10 () INHERITS (ch_merge_r10)",
			"CREATE TABLE p_part_r10 (k int, a text COLLATE c_inh_r10) PARTITION BY LIST (k)",
			"CREATE TABLE ch_part_r10 PARTITION OF p_part_r10 FOR VALUES IN (1)",
			"DROP COLLATION c_inh_r10 CASCADE",
		} {
			if code := simulationCatalogCode(t, c112, stmt); code != "0" {
				t.Fatalf("c112 inh cascade %s: %s", stmt, code)
			}
		}
		for _, tbl := range []string{"p_inh_r10", "ch_inh_r10", "ch_merge_r10", "gch_inh_r10", "p_part_r10", "ch_part_r10"} {
			r := c112.GetRelation("public", tbl)
			if r == nil {
				t.Fatalf("expected %s to remain after DROP COLLATION CASCADE", tbl)
			}
			for _, col := range r.Columns {
				if col.Name == "a" || col.Name == "b" {
					t.Fatalf("expected column %s.%s to be dropped by DROP COLLATION CASCADE", tbl, col.Name)
				}
			}
		}
		for _, stmt := range []string{
			"CREATE COLLATION c_pkey_r10 (LOCALE = 'C')",
			"CREATE TABLE p_pkey_r10 (a text COLLATE c_pkey_r10) PARTITION BY LIST (a)",
			"CREATE TABLE ch_pkey_r10 PARTITION OF p_pkey_r10 FOR VALUES IN ('x')",
		} {
			if code := simulationCatalogCode(t, c112, stmt); code != "0" {
				t.Fatalf("c112 pkey setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c112, "DROP COLLATION c_pkey_r10 CASCADE"); code != CodeInvalidTableDefinition {
			t.Fatalf("DROP COLLATION CASCADE on partition-key column answered %s, want %s", code, CodeInvalidTableDefinition)
		}
		if c112.schemaByName["public"].Collations["c_pkey_r10"] == nil {
			t.Fatal("failed DROP COLLATION CASCADE deleted c_pkey_r10 despite partition-key error")
		}
		if len(c112.GetRelation("public", "ch_pkey_r10").Columns) != 1 {
			t.Fatal("failed DROP COLLATION CASCADE mutated child partition columns")
		}

		// 4. Rules diff and migration on tables and views (create, modify, rename, firing mode, drop, view rebuild).
		c112From := New()
		c112To := New()
		for _, stmt := range []string{
			"CREATE TABLE t_rule_r10 (id int, val varchar(10))",
			"CREATE TABLE t_rule_log_r10 (id int, val text)",
			"CREATE RULE r_keep_mode_r10 AS ON INSERT TO t_rule_r10 WHERE NEW.id > 0 DO ALSO INSERT INTO t_rule_log_r10 (id, val) VALUES (NEW.id, NEW.val)",
			"CREATE RULE r_mod_r10 AS ON DELETE TO t_rule_r10 DO NOTHING",
			"CREATE RULE r_drop_r10 AS ON UPDATE TO t_rule_r10 DO NOTHING",
			"CREATE VIEW v_rule_r10 AS SELECT id, val FROM t_rule_r10",
			"CREATE RULE r_v_ins_r10 AS ON INSERT TO v_rule_r10 DO INSTEAD INSERT INTO t_rule_r10 (id, val) VALUES (NEW.id, NEW.val)",
		} {
			if code := simulationCatalogCode(t, c112From, stmt); code != "0" {
				t.Fatalf("c112From %s: %s", stmt, code)
			}
		}
		for _, stmt := range []string{
			"CREATE TABLE t_rule_r10 (id int, val varchar(50))",
			"CREATE TABLE t_rule_log_r10 (id int, val text)",
			"CREATE RULE r_keep_mode_r10 AS ON INSERT TO t_rule_r10 WHERE NEW.id > 0 DO ALSO INSERT INTO t_rule_log_r10 (id, val) VALUES (NEW.id, NEW.val)",
			"ALTER TABLE t_rule_r10 DISABLE RULE r_keep_mode_r10",
			"CREATE RULE r_mod_r10 AS ON DELETE TO t_rule_r10 WHERE OLD.id > 10 DO ALSO (DELETE FROM t_rule_log_r10 WHERE id = OLD.id; NOTIFY rule_ch_r10, 'del')",
			"ALTER TABLE t_rule_r10 ENABLE ALWAYS RULE r_mod_r10",
			"CREATE RULE r_oldname_r10 AS ON UPDATE TO t_rule_r10 DO NOTHING",
			"ALTER RULE r_oldname_r10 ON t_rule_r10 RENAME TO r_added_r10",
			"CREATE VIEW v_rule_r10 AS SELECT id, val FROM t_rule_r10",
			"CREATE RULE r_v_ins_r10 AS ON INSERT TO v_rule_r10 DO INSTEAD INSERT INTO t_rule_r10 (id, val) VALUES (NEW.id, NEW.val)",
			"CREATE RULE r_v_del_r10 AS ON DELETE TO v_rule_r10 DO INSTEAD DELETE FROM t_rule_r10 WHERE id = OLD.id",
		} {
			if code := simulationCatalogCode(t, c112To, stmt); code != "0" {
				t.Fatalf("c112To %s: %s", stmt, code)
			}
		}
		migEmpty112 := GenerateMigration(New(), c112To, Diff(New(), c112To))
		c112Empty := New()
		if code := simulationCatalogCode(t, c112Empty, migEmpty112.SQL()); code != "0" {
			t.Fatalf("applying empty->c112To failed with %s:\n%s", code, migEmpty112.SQL())
		}
		if rem := GenerateMigration(c112Empty, c112To, Diff(c112Empty, c112To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining ops after empty->c112To, got:\n%s", rem.SQL())
		}
		mig112 := GenerateMigration(c112From, c112To, Diff(c112From, c112To))
		if !strings.Contains(mig112.SQL(), `DISABLE RULE "r_keep_mode_r10"`) || !strings.Contains(mig112.SQL(), `ENABLE ALWAYS RULE "r_mod_r10"`) {
			t.Fatalf("expected rule firing mode DDL in migration, got:\n%s", mig112.SQL())
		}
		if code := simulationCatalogCode(t, c112From, mig112.SQL()); code != "0" {
			t.Fatalf("applying c112From->c112To failed with %s:\n%s", code, mig112.SQL())
		}
		if rem := GenerateMigration(c112From, c112To, Diff(c112From, c112To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining ops after c112From->c112To, got:\n%s", rem.SQL())
		}

		// Review round 11 regressions:
		// 1. Rule actions refresh referenced table and column names after renames.
		// 2. Window function OVER clauses (inline frame specs and WINDOW clauses)
		//    are preserved when deparsing rule actions.
		// 3. Unpopulated materialized views (WITH NO DATA) emit WITH NO DATA in
		//    migration DDL and diff against populated materialized views.
		// 4. User collations record OwnerName and participate in DROP ROLE,
		//    RENAME ROLE, REASSIGN OWNED, DROP OWNED, and ALTER COLLATION OWNER/RENAME.
		c113To := New()
		for _, stmt := range []string{
			"CREATE TABLE src_r11 (a int)",
			"CREATE TABLE dst_r11 (x int, rn bigint)",
			"CREATE RULE r_ren_r11 AS ON INSERT TO src_r11 WHERE NEW.a > 0 DO ALSO INSERT INTO dst_r11 (x, rn) SELECT NEW.a, row_number() OVER (PARTITION BY NEW.a ORDER BY NEW.a DESC ROWS BETWEEN 1 PRECEDING AND CURRENT ROW EXCLUDE CURRENT ROW)",
			"CREATE RULE r_win_r11 AS ON DELETE TO src_r11 DO ALSO INSERT INTO dst_r11 (x, rn) SELECT OLD.a, count(*) OVER w FROM src_r11 WINDOW w AS (ORDER BY a)",
			"ALTER TABLE dst_r11 RENAME TO dst2_r11",
			"ALTER TABLE dst2_r11 RENAME COLUMN x TO x2",
			"ALTER TABLE src_r11 RENAME COLUMN a TO a2",
			"CREATE MATERIALIZED VIEW mv_nodata_r11 AS SELECT a2 FROM src_r11 WITH NO DATA",
		} {
			if code := simulationCatalogCode(t, c113To, stmt); code != "0" {
				t.Fatalf("c113To %s: %s", stmt, code)
			}
		}
		migEmpty113 := GenerateMigration(New(), c113To, Diff(New(), c113To))
		if !strings.Contains(migEmpty113.SQL(), "OVER (PARTITION BY") || !strings.Contains(migEmpty113.SQL(), "OVER w") {
			t.Fatalf("expected OVER clauses in rule migration DDL, got:\n%s", migEmpty113.SQL())
		}
		if !strings.Contains(migEmpty113.SQL(), `CREATE MATERIALIZED VIEW "public"."mv_nodata_r11" AS`) || !strings.Contains(migEmpty113.SQL(), "WITH NO DATA") {
			t.Fatalf("expected WITH NO DATA in unpopulated matview DDL, got:\n%s", migEmpty113.SQL())
		}
		c113Applied := New()
		if code := simulationCatalogCode(t, c113Applied, migEmpty113.SQL()); code != "0" {
			t.Fatalf("applying empty->c113To failed with %s:\n%s", code, migEmpty113.SQL())
		}
		if mv := c113Applied.GetRelation("public", "mv_nodata_r11"); mv == nil || mv.IsPopulated {
			t.Fatalf("expected migrated mv_nodata_r11 to remain unpopulated, got %+v", mv)
		}
		if rem := GenerateMigration(c113Applied, c113To, Diff(c113Applied, c113To)); len(rem.Ops) != 0 {
			t.Fatalf("expected zero remaining ops after empty->c113To, got:\n%s", rem.SQL())
		}

		// Collation role lifecycle checks.
		c113Role := New()
		for _, stmt := range []string{
			"CREATE ROLE alice_r11",
			"CREATE ROLE bob_r11",
			"SET ROLE alice_r11",
			"CREATE COLLATION c_owned_r11 (LOCALE = 'C')",
			"RESET ROLE",
			"CREATE TABLE t_uses_c_r11 (v text COLLATE c_owned_r11)",
		} {
			if code := simulationCatalogCode(t, c113Role, stmt); code != "0" {
				t.Fatalf("c113Role setup %s: %s", stmt, code)
			}
		}
		if code := simulationCatalogCode(t, c113Role, "DROP ROLE alice_r11"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE owning collation answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c113Role, "ALTER ROLE alice_r11 RENAME TO alice2_r11"); code != "0" {
			t.Fatalf("ALTER ROLE RENAME answered %s, want 0", code)
		}
		if code := simulationCatalogCode(t, c113Role, "DROP ROLE alice2_r11"); code != CodeDependentObjects {
			t.Fatalf("DROP ROLE after rename answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c113Role, "REASSIGN OWNED BY alice2_r11 TO bob_r11; DROP ROLE alice2_r11"); code != "0" {
			t.Fatalf("REASSIGN OWNED + DROP ROLE answered %s, want 0", code)
		}
		if code := simulationCatalogCode(t, c113Role, "DROP OWNED BY bob_r11"); code != CodeDependentObjects {
			t.Fatalf("DROP OWNED RESTRICT with external dependent column answered %s, want %s", code, CodeDependentObjects)
		}
		if code := simulationCatalogCode(t, c113Role, "DROP OWNED BY bob_r11 CASCADE"); code != "0" {
			t.Fatalf("DROP OWNED CASCADE answered %s, want 0", code)
		}
		if c113Role.schemaByName["public"].Collations["c_owned_r11"] != nil {
			t.Fatal("expected c_owned_r11 to be dropped by DROP OWNED BY bob_r11 CASCADE")
		}
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := conn.ExecContext(ctx, "DROP SCHEMA IF EXISTS public CASCADE;\nCREATE SCHEMA public;\nRESET default_tablespace;\n"+fixtures[tc.fixture]); err != nil {
				t.Fatalf("build fixture %s: %v", tc.fixture, err)
			}
			_, err := conn.ExecContext(ctx, tc.stmt)
			engineCode := simulationServerCode(t, tc.stmt, err)
			if engineCode != tc.code {
				t.Fatalf("postgres: %s answered %s, want %s", tc.stmt, engineCode, tc.code)
			}

			c := New()
			if code := simulationCatalogCode(t, c, fixtures[tc.fixture]); code != "0" {
				t.Fatalf("catalog refused fixture %s with %s", tc.fixture, code)
			}
			if code := simulationCatalogCode(t, c, tc.stmt); code != tc.code {
				t.Fatalf("catalog: %s answered %s, want %s", tc.stmt, code, tc.code)
			}
		})
	}
}

// simulationServerCode reads PostgreSQL's answer as a SQLSTATE, "0" when it
// took the statement.
func simulationServerCode(t *testing.T, stmt string, err error) string {
	t.Helper()
	if err == nil {
		return "0"
	}
	var refusal *pgconn.PgError
	if !errors.As(err, &refusal) {
		t.Fatalf("%s: not a PostgreSQL error: %v", stmt, err)
	}
	return refusal.Code
}

// simulationCatalogCode reads the catalog's answer as a SQLSTATE, "0" when it
// took every statement.
func simulationCatalogCode(t *testing.T, c *Catalog, sql string) string {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", sql, err)
	}
	code := "0"
	for _, result := range results {
		if result.Error == nil {
			continue
		}
		var refusal *Error
		if !errors.As(result.Error, &refusal) {
			t.Fatalf("%s: not a catalog error: %v", sql, result.Error)
		}
		code = refusal.Code
	}
	return code
}
