package catalog

import (
	"testing"
)

func TestOracleInfra(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	t.Run("container_starts", func(t *testing.T) {
		// Scenario 1: PG testcontainer starts and connects successfully
		err := oracle.db.PingContext(oracle.ctx)
		if err != nil {
			t.Fatalf("ping failed: %v", err)
		}
	})

	t.Run("schema_isolation", func(t *testing.T) {
		// Scenario 2: Schema-level isolation works
		schema := oracle.freshSchema(t)
		oracle.execInSchema(t, schema, `CREATE TABLE iso_t (id int)`)
		// Verify table exists in the schema
		var count int
		err := oracle.db.QueryRowContext(oracle.ctx,
			`SELECT count(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name = 'iso_t'`, schema).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("expected table to exist, got count=%d", count)
		}
		// freshSchema registers cleanup, so after test the schema is dropped
	})

	t.Run("multi_statement_ddl", func(t *testing.T) {
		// Scenario 3: SQL execution helper executes multi-statement DDL on PG
		schema := oracle.freshSchema(t)
		oracle.execInSchema(t, schema, `
			CREATE TABLE t1 (id int PRIMARY KEY);
			CREATE TABLE t2 (id int, t1_id int REFERENCES t1(id));
			CREATE FUNCTION my_func() RETURNS trigger LANGUAGE plpgsql AS $$
			BEGIN RETURN NEW; END;
			$$;
		`)
		// Verify all objects created
		var tableCount int
		err := oracle.db.QueryRowContext(oracle.ctx,
			`SELECT count(*) FROM information_schema.tables WHERE table_schema = $1`, schema).Scan(&tableCount)
		if err != nil {
			t.Fatal(err)
		}
		if tableCount < 2 {
			t.Errorf("expected >= 2 tables, got %d", tableCount)
		}
		funcs := oracle.queryFunctions(t, schema)
		if len(funcs) == 0 {
			t.Error("expected function my_func to exist")
		}
	})

	t.Run("comparison_queries_all_types", func(t *testing.T) {
		// Scenario 4: Schema comparison queries return structured results for all object types
		schema := oracle.freshSchema(t)
		oracle.execInSchema(t, schema, `
			CREATE TYPE status AS ENUM ('active', 'inactive');
			CREATE SEQUENCE my_seq INCREMENT 5;
			CREATE TABLE users (
				id int PRIMARY KEY,
				name text NOT NULL DEFAULT 'anon',
				s status DEFAULT 'active'
			);
			CREATE INDEX idx_name ON users(name);
			CREATE VIEW active_users AS SELECT * FROM users WHERE s = 'active';
			CREATE FUNCTION greet(text) RETURNS text LANGUAGE sql AS 'SELECT ''hi '' || $1';
			COMMENT ON TABLE users IS 'User table';
		`)
		// Call each query function and verify non-empty results
		tables := oracle.queryTables(t, schema)
		if len(tables) == 0 {
			t.Error("no tables found")
		}
		cols := oracle.queryColumns(t, schema)
		if len(cols) == 0 {
			t.Error("no columns found")
		}
		idxs := oracle.queryIndexes(t, schema)
		if len(idxs) == 0 {
			t.Error("no indexes found")
		}
		cons := oracle.queryConstraints(t, schema)
		if len(cons) == 0 {
			t.Error("no constraints found")
		}
		funcs := oracle.queryFunctions(t, schema)
		if len(funcs) == 0 {
			t.Error("no functions found")
		}
		views := oracle.queryViews(t, schema)
		if len(views) == 0 {
			t.Error("no views found")
		}
		seqs := oracle.querySequences(t, schema)
		if len(seqs) == 0 {
			t.Error("no sequences found")
		}
		enums := oracle.queryEnumTypes(t, schema)
		if len(enums) == 0 {
			t.Error("no enum types found")
		}
		comments := oracle.queryComments(t, schema)
		if len(comments) == 0 {
			t.Error("no comments found")
		}
	})

	t.Run("comparison_identical", func(t *testing.T) {
		// Scenario 5: Comparison function detects identical schemas as equal
		ddl := `CREATE TABLE t (id int PRIMARY KEY, name text NOT NULL);`
		schemaA := oracle.freshSchema(t)
		schemaB := oracle.freshSchema(t)
		oracle.execInSchema(t, schemaA, ddl)
		oracle.execInSchema(t, schemaB, ddl)
		oracle.assertSchemasEqual(t, schemaA, schemaB)
	})

	t.Run("comparison_detects_diff", func(t *testing.T) {
		// Scenario 6: Comparison function detects schema differences correctly
		schemaA := oracle.freshSchema(t)
		schemaB := oracle.freshSchema(t)
		oracle.execInSchema(t, schemaA, `CREATE TABLE t (id int);`)
		oracle.execInSchema(t, schemaB, `CREATE TABLE t (id int, name text);`)
		diffs := oracle.compareSchemas(t, schemaA, schemaB)
		if len(diffs) == 0 {
			t.Error("expected differences, got none")
		}
	})

	t.Run("roundtrip_helper_smoke", func(t *testing.T) {
		// Scenario 7: Migration roundtrip helper works end-to-end
		assertOracleRoundtrip(t, oracle,
			`CREATE TABLE t (id int);`,
			`CREATE TABLE t (id int, name text DEFAULT 'x');`,
		)
	})
}
