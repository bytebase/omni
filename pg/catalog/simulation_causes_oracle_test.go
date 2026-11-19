package catalog

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/pg/internal/testcontainer"
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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := conn.ExecContext(ctx, "DROP SCHEMA IF EXISTS public CASCADE;\nCREATE SCHEMA public;\n"+fixtures[tc.fixture]); err != nil {
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
