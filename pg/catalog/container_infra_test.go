package catalog

import (
	"testing"
)

func TestContainerInfra(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test: requires Docker")
	}
	ctr := startPGContainer(t)

	t.Run("container_starts", func(t *testing.T) {
		// Scenario 1: PG testcontainer starts and connects successfully
		err := ctr.db.PingContext(ctr.ctx)
		if err != nil {
			t.Fatalf("ping failed: %v", err)
		}
	})

	t.Run("schema_isolation", func(t *testing.T) {
		// Scenario 2: Schema-level isolation works
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, `CREATE TABLE iso_t (id int)`)
		// Verify table exists in the schema
		var count int
		err := ctr.db.QueryRowContext(ctr.ctx,
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
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, `
			CREATE TABLE t1 (id int PRIMARY KEY);
			CREATE TABLE t2 (id int, t1_id int REFERENCES t1(id));
			CREATE FUNCTION my_func() RETURNS trigger LANGUAGE plpgsql AS $$
			BEGIN RETURN NEW; END;
			$$;
		`)
		// Verify all objects created
		var tableCount int
		err := ctr.db.QueryRowContext(ctr.ctx,
			`SELECT count(*) FROM information_schema.tables WHERE table_schema = $1`, schema).Scan(&tableCount)
		if err != nil {
			t.Fatal(err)
		}
		if tableCount < 2 {
			t.Errorf("expected >= 2 tables, got %d", tableCount)
		}
		funcs := ctr.queryFunctions(t, schema)
		if len(funcs) == 0 {
			t.Error("expected function my_func to exist")
		}
	})

	t.Run("comparison_queries_all_types", func(t *testing.T) {
		// Scenario 4: Schema comparison queries return structured results for all object types
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, `
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
		tables := ctr.queryTables(t, schema)
		if len(tables) == 0 {
			t.Error("no tables found")
		}
		cols := ctr.queryColumns(t, schema)
		if len(cols) == 0 {
			t.Error("no columns found")
		}
		idxs := ctr.queryIndexes(t, schema)
		if len(idxs) == 0 {
			t.Error("no indexes found")
		}
		cons := ctr.queryConstraints(t, schema)
		if len(cons) == 0 {
			t.Error("no constraints found")
		}
		funcs := ctr.queryFunctions(t, schema)
		if len(funcs) == 0 {
			t.Error("no functions found")
		}
		views := ctr.queryViews(t, schema)
		if len(views) == 0 {
			t.Error("no views found")
		}
		seqs := ctr.querySequences(t, schema)
		if len(seqs) == 0 {
			t.Error("no sequences found")
		}
		enums := ctr.queryEnumTypes(t, schema)
		if len(enums) == 0 {
			t.Error("no enum types found")
		}
		comments := ctr.queryComments(t, schema)
		if len(comments) == 0 {
			t.Error("no comments found")
		}
	})

	t.Run("comparison_identical", func(t *testing.T) {
		// Scenario 5: Comparison function detects identical schemas as equal
		ddl := `CREATE TABLE t (id int PRIMARY KEY, name text NOT NULL);`
		schemaA := ctr.freshSchema(t)
		schemaB := ctr.freshSchema(t)
		ctr.execInSchema(t, schemaA, ddl)
		ctr.execInSchema(t, schemaB, ddl)
		ctr.assertSchemasEqual(t, schemaA, schemaB)
	})

	t.Run("comparison_detects_diff", func(t *testing.T) {
		// Scenario 6: Comparison function detects schema differences correctly
		schemaA := ctr.freshSchema(t)
		schemaB := ctr.freshSchema(t)
		ctr.execInSchema(t, schemaA, `CREATE TABLE t (id int);`)
		ctr.execInSchema(t, schemaB, `CREATE TABLE t (id int, name text);`)
		diffs := ctr.compareSchemas(t, schemaA, schemaB)
		if len(diffs) == 0 {
			t.Error("expected differences, got none")
		}
	})

	t.Run("roundtrip_helper_smoke", func(t *testing.T) {
		// Scenario 7: Migration roundtrip helper works end-to-end
		assertOracleRoundtrip(t, ctr,
			`CREATE TABLE t (id int);`,
			`CREATE TABLE t (id int, name text DEFAULT 'x');`,
		)
	})
}
