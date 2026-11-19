package catalog

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/metadata"
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
