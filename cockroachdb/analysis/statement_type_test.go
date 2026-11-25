package analysis

import (
	"testing"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser/statements"
	"github.com/google/go-cmp/cmp"
)

func TestStatementTypes(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      []StatementType
	}{
		{statement: `CREATE DATABASE d`, want: []StatementType{CreateDatabase}},
		{statement: `CREATE TABLE t (id INT PRIMARY KEY)`, want: []StatementType{CreateTable}},
		{statement: `CREATE TABLE t2 AS SELECT * FROM t`, want: []StatementType{CreateTable}},
		{statement: `CREATE VIEW v AS SELECT * FROM t`, want: []StatementType{CreateView}},
		{statement: `CREATE MATERIALIZED VIEW mv AS SELECT * FROM t`, want: []StatementType{CreateView}},
		{statement: `CREATE INDEX ON t (c)`, want: []StatementType{CreateIndex}},
		{statement: `CREATE SEQUENCE s`, want: []StatementType{CreateSequence}},
		{statement: `CREATE SCHEMA app`, want: []StatementType{CreateSchema}},
		{statement: `CREATE FUNCTION f() RETURNS INT LANGUAGE SQL AS 'SELECT 1'`, want: []StatementType{CreateFunction}},
		{statement: `CREATE PROCEDURE p() LANGUAGE SQL AS 'SELECT 1'`, want: []StatementType{CreateFunction}},
		{statement: `CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION f()`, want: []StatementType{CreateTrigger}},
		{statement: `CREATE EXTENSION IF NOT EXISTS pg_trgm`, want: []StatementType{CreateExtension}},
		{statement: `CREATE TYPE e AS ENUM ('a', 'b')`, want: []StatementType{CreateType}},
		{statement: `DROP DATABASE d`, want: []StatementType{DropDatabase}},
		{statement: `DROP TABLE t`, want: []StatementType{DropTable}},
		{statement: `DROP MATERIALIZED VIEW mv`, want: []StatementType{DropTable}},
		{statement: `DROP VIEW v`, want: []StatementType{DropView}},
		{statement: `DROP INDEX t@idx`, want: []StatementType{DropIndex}},
		{statement: `DROP SEQUENCE s`, want: []StatementType{DropSequence}},
		{statement: `DROP SCHEMA app`, want: []StatementType{DropSchema}},
		{statement: `DROP TYPE e`, want: []StatementType{DropType}},
		{statement: `DROP TRIGGER tr ON t`, want: []StatementType{DropTrigger}},
		{statement: `DROP FUNCTION f`, want: []StatementType{DropFunction}},
		{statement: `TRUNCATE t`, want: []StatementType{Truncate}},
		{statement: `ALTER DATABASE d RENAME TO d2`, want: []StatementType{AlterDatabase}},
		{statement: `ALTER DATABASE d OWNER TO u`, want: []StatementType{AlterDatabase}},
		{statement: `ALTER DATABASE d PRIMARY REGION "us-east1"`, want: []StatementType{AlterDatabase}},
		{statement: `ALTER DATABASE d SURVIVE REGION FAILURE`, want: []StatementType{AlterDatabase}},
		{statement: `ALTER DATABASE d CONFIGURE ZONE USING num_replicas = 5`, want: []StatementType{AlterDatabase}},
		{statement: `ALTER DATABASE d SET timezone = 'UTC'`, want: []StatementType{AlterDatabase}},
		{statement: `ALTER ROLE u SET timezone = 'UTC'`},
		{statement: `ALTER TABLE t CONFIGURE ZONE USING num_replicas = 5`, want: []StatementType{AlterTable}},
		{statement: `ALTER INDEX t@idx CONFIGURE ZONE USING num_replicas = 5`, want: []StatementType{AlterIndex}},
		{statement: `ALTER RANGE default CONFIGURE ZONE USING num_replicas = 5`},
		{statement: `ALTER TABLE t ADD COLUMN c2 INT`, want: []StatementType{AlterTable}},
		{statement: `ALTER TABLE t RENAME COLUMN c TO c2`, want: []StatementType{AlterTable}},
		{statement: `ALTER TABLE t SET LOCALITY REGIONAL BY ROW`, want: []StatementType{AlterTable}},
		{statement: `ALTER TABLE t SET SCHEMA app`, want: []StatementType{AlterTable}},
		{statement: `ALTER VIEW v OWNER TO u`, want: []StatementType{AlterView}},
		{statement: `ALTER SEQUENCE s OWNER TO u`, want: []StatementType{AlterSequence}},
		{statement: `ALTER INDEX t@idx NOT VISIBLE`, want: []StatementType{AlterIndex}},
		{statement: `ALTER SEQUENCE s RESTART WITH 100`, want: []StatementType{AlterSequence}},
		{statement: `ALTER TYPE e ADD VALUE 'c'`, want: []StatementType{AlterType}},
		{statement: `ALTER TABLE t RENAME TO t2`, want: []StatementType{AlterTable}},
		{statement: `ALTER VIEW v RENAME TO v2`, want: []StatementType{AlterView}},
		{statement: `ALTER MATERIALIZED VIEW mv RENAME TO mv2`, want: []StatementType{AlterTable}},
		{statement: `ALTER SEQUENCE s RENAME TO s2`, want: []StatementType{RenameSequence}},
		{statement: `ALTER INDEX t@idx RENAME TO idx2`, want: []StatementType{RenameIndex}},
		{statement: `ALTER SCHEMA app RENAME TO app2`, want: []StatementType{RenameSchema}},
		{statement: `COMMENT ON TABLE t IS 'comment'`, want: []StatementType{Comment}},
		{statement: `COMMENT ON COLUMN t.c IS 'comment'`, want: []StatementType{Comment}},
		{statement: `INSERT INTO t VALUES (1)`, want: []StatementType{Insert}},
		{statement: `UPSERT INTO t VALUES (1)`, want: []StatementType{Insert}},
		{statement: `UPDATE t SET c = 1`, want: []StatementType{Update}},
		{statement: `DELETE FROM t WHERE id = 1`, want: []StatementType{Delete}},
		{statement: `CALL archive_orders(30)`, want: []StatementType{Execute}},
		{statement: `DO $$ BEGIN DELETE FROM t; END $$`, want: []StatementType{Execute}},
		{statement: `EXECUTE purge_stmt(30)`, want: []StatementType{Execute}},
		{statement: `PREPARE purge_stmt (INT) AS DELETE FROM t WHERE id = $1`},
		{
			statement: `SET application_name = 'x'; CREATE TABLE t (id INT PRIMARY KEY); SELECT 1; INSERT INTO t VALUES (1)`,
			want:      []StatementType{CreateTable, Insert},
		},
		{statement: `WITH d AS (DELETE FROM t RETURNING id) SELECT count(*) FROM d`, want: []StatementType{Delete}},
		{statement: `SELECT * FROM [UPDATE t SET c = 1 RETURNING id]`, want: []StatementType{Update}},
		{statement: `SELECT * FROM [EXPLAIN ANALYZE DELETE FROM t]`, want: []StatementType{Delete}},
		{statement: `SELECT row_number() OVER (PARTITION BY (SELECT 1 FROM [DELETE FROM t RETURNING 1])) FROM u`, want: []StatementType{Delete}},
		{
			statement: `WITH d AS (DELETE FROM t RETURNING *) INSERT INTO t2 SELECT * FROM d`,
			want:      []StatementType{Insert, Delete},
		},
		{
			statement: `CREATE TABLE t_moved AS SELECT * FROM [DELETE FROM t RETURNING *]`,
			want:      []StatementType{CreateTable, Delete},
		},
		{statement: `EXPLAIN ANALYZE DELETE FROM t WHERE id = 1`, want: []StatementType{Delete}},
		{statement: `EXPLAIN ANALYZE (VERBOSE) UPSERT INTO t VALUES (1)`, want: []StatementType{Insert}},
		{statement: `EXPLAIN DELETE FROM t WHERE id = 1`},
		{statement: `SELECT * FROM t`},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			got := StatementTypes(mustParse(t, tc.statement))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StatementTypes() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func mustParse(t *testing.T, script string) statements.Statements {
	t.Helper()
	parsed, err := Parse(script)
	if err != nil {
		t.Fatalf("Parse(%q): %v", script, err)
	}
	return parsed
}

func TestStatementTypesFromParsed(t *testing.T) {
	parsed, err := ParseScript(`CREATE TABLE t (id INT); UPSERT INTO t VALUES (1);`)
	if err != nil {
		t.Fatalf("ParseScript: %v", err)
	}
	want := []StatementType{CreateTable, Insert}
	if diff := cmp.Diff(want, StatementTypesFromParsed(parsed)); diff != "" {
		t.Errorf("StatementTypesFromParsed() mismatch (-want +got):\n%s", diff)
	}
}
