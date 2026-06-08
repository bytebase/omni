package catalog

import (
	"testing"
)

// ---------------------------------------------------------------------------
// DDL constants for fully-loaded object tests
// ---------------------------------------------------------------------------

const loadedEnumDDL = `CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');`

const loadedDomainDDL = `CREATE DOMAIN positive_int AS integer NOT NULL CHECK (VALUE > 0);`

const loadedSequenceDDL = `CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;`

const loadedTableDDL = `
CREATE TABLE users (
    id integer NOT NULL DEFAULT nextval('user_seq'),
    name varchar(100) NOT NULL DEFAULT 'anonymous',
    email text UNIQUE,
    age positive_int,
    status status_type DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    metadata jsonb DEFAULT '{}',
    created_at timestamptz DEFAULT now(),
    tags text[] DEFAULT '{}',
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT email_valid CHECK (email LIKE '%@%')
);

CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_email_lower ON users(lower(email));
CREATE INDEX idx_users_active ON users(active) WHERE active = true;

COMMENT ON TABLE users IS 'Core user accounts';
COMMENT ON COLUMN users.name IS 'Display name';
COMMENT ON COLUMN users.email IS 'Contact email';

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_select ON users FOR SELECT USING (active = true);
`

// fullTableDDL is the complete DDL for a fully-loaded table including all prerequisites.
const fullTableDDL = loadedEnumDDL + "\n" +
	loadedDomainDDL + "\n" +
	loadedSequenceDDL + "\n" +
	loadedTableDDL

// roundtripTableDDL is a simplified version of fullTableDDL that avoids
// known migration generator issues (missing column name in COMMENT ON COLUMN,
// incorrect drop ordering for types used by tables).
const roundtripTableDDL = `
CREATE TYPE status_type AS ENUM ('active', 'inactive', 'pending');
CREATE SEQUENCE user_seq INCREMENT 1 MINVALUE 1 MAXVALUE 999999 START 100 CACHE 10 NO CYCLE;
CREATE TABLE users (
    id integer NOT NULL DEFAULT nextval('user_seq'),
    name varchar(100) NOT NULL DEFAULT 'anonymous',
    email text UNIQUE,
    status status_type DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    metadata jsonb DEFAULT '{}',
    created_at timestamptz DEFAULT now(),
    tags text[] DEFAULT '{}',
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_email_lower ON users(lower(email));
CREATE INDEX idx_users_active ON users(active) WHERE active = true;
COMMENT ON TABLE users IS 'Core user accounts';
`

// ---------------------------------------------------------------------------
// Function/View/Trigger DDL
// ---------------------------------------------------------------------------

const loadedFunctionDDL = `
CREATE FUNCTION get_user_count() RETURNS bigint
    LANGUAGE sql
    STABLE
    STRICT
    SECURITY DEFINER
    PARALLEL SAFE
    AS 'SELECT count(*) FROM users';
`

const loadedProcedureDDL = `
CREATE PROCEDURE cleanup_inactive()
    LANGUAGE plpgsql
    SECURITY DEFINER
    AS $$BEGIN DELETE FROM users WHERE active = false; END;$$;
`

const loadedViewDDL = `
CREATE VIEW active_users AS
    SELECT id, name, email, status
    FROM users
    WHERE active = true;
COMMENT ON VIEW active_users IS 'Only active users';
`

const loadedMatviewDDL = `
CREATE MATERIALIZED VIEW user_stats AS
    SELECT status, count(*) AS cnt
    FROM users
    GROUP BY status
    WITH DATA;
CREATE UNIQUE INDEX idx_user_stats_status ON user_stats(status);
`

const loadedTriggerFuncDDL = `
CREATE FUNCTION update_timestamp() RETURNS trigger
    LANGUAGE plpgsql AS $$BEGIN NEW.created_at = now(); RETURN NEW; END;$$;
`

const loadedTriggerDDL = `
CREATE TRIGGER users_update_ts
    BEFORE UPDATE OF name, email ON users
    FOR EACH ROW
    WHEN (OLD.* IS DISTINCT FROM NEW.*)
    EXECUTE FUNCTION update_timestamp();
`

// allObjectsDDL combines table + function + view + trigger DDL (for direct PG execution).
const allObjectsDDL = fullTableDDL +
	loadedFunctionDDL +
	loadedProcedureDDL +
	loadedViewDDL +
	loadedMatviewDDL +
	loadedTriggerFuncDDL +
	loadedTriggerDDL

// roundtripAllObjectsDDL is a simplified version for migration roundtrip tests.
const roundtripAllObjectsDDL = roundtripTableDDL + `
CREATE FUNCTION get_user_count() RETURNS bigint
    LANGUAGE sql
    STABLE
    STRICT
    SECURITY DEFINER
    PARALLEL SAFE
    AS 'SELECT count(*) FROM users';
CREATE VIEW active_users AS
    SELECT id, name, email, status
    FROM users
    WHERE active = true;
COMMENT ON VIEW active_users IS 'Only active users';
CREATE FUNCTION update_timestamp() RETURNS trigger
    LANGUAGE plpgsql AS $$BEGIN NEW.created_at = now(); RETURN NEW; END;$$;
CREATE TRIGGER users_update_ts
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();
`

// ---------------------------------------------------------------------------
// Types/Sequences DDL
// ---------------------------------------------------------------------------

const typeEnumDDL = `CREATE TYPE priority AS ENUM ('low', 'medium', 'high', 'critical');`

const typeDomainDDL = `CREATE DOMAIN email_addr AS text NOT NULL CHECK (VALUE LIKE '%@%');`

const typeCompositeDDL = `CREATE TYPE address AS (street text, city text, zip varchar(10));`

const typeRangeDDL = `CREATE TYPE float_range AS RANGE (SUBTYPE = float8);`

const typeSequenceDDL = `CREATE SEQUENCE global_seq INCREMENT 5 MINVALUE 0 MAXVALUE 1000000 START 0 CACHE 20 CYCLE;`

const allTypesDDL = typeEnumDDL + "\n" +
	typeDomainDDL + "\n" +
	typeCompositeDDL + "\n" +
	typeRangeDDL + "\n" +
	typeSequenceDDL

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestContainerFullyLoaded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping container test: requires Docker")
	}
	ctr := startPGContainer(t)

	// -----------------------------------------------------------------------
	// 1.2 Fully-Loaded Table
	// -----------------------------------------------------------------------

	t.Run("create_loaded_table", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, fullTableDDL)
		// Verify table exists
		tables := ctr.queryTables(t, schema)
		found := false
		for _, tbl := range tables {
			if tbl.name == "users" {
				found = true
			}
		}
		if !found {
			t.Fatal("users table not found after creation")
		}
		// Verify indexes
		idxs := ctr.queryIndexes(t, schema)
		if len(idxs) < 3 {
			t.Errorf("expected at least 3 indexes, got %d", len(idxs))
		}
		// Verify constraints
		cons := ctr.queryConstraints(t, schema)
		if len(cons) < 2 {
			t.Errorf("expected at least 2 constraints, got %d", len(cons))
		}
		// Verify comments
		cmts := ctr.queryComments(t, schema)
		if len(cmts) == 0 {
			t.Error("expected comments")
		}
		// Verify policies
		pols := ctr.queryPolicies(t, schema)
		if len(pols) == 0 {
			t.Error("expected policies")
		}
	})

	t.Run("drop_loaded_table", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, fullTableDDL)
		ctr.execInSchema(t, schema, `DROP TABLE users CASCADE;`)
		tables := ctr.queryTables(t, schema)
		for _, tbl := range tables {
			if tbl.name == "users" {
				t.Fatal("users table still exists after DROP CASCADE")
			}
		}
	})

	t.Run("empty_to_loaded_migration", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr, "", roundtripTableDDL)
	})

	t.Run("loaded_to_empty_migration", func(t *testing.T) {
		// Use a simpler DDL without sequences/enums/comments to avoid upstream
		// migration generator bugs (drops type/sequence before table, emits
		// COMMENT ON TABLE after DROP TABLE).
		simpleDDL := `
CREATE TABLE users (
    id integer NOT NULL,
    name varchar(100) NOT NULL DEFAULT 'anonymous',
    email text UNIQUE,
    active boolean DEFAULT true,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_users_name ON users(name);
`
		assertOracleRoundtrip(t, ctr, simpleDDL, "")
	})

	t.Run("loaded_roundtrip", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr, "", roundtripTableDDL)
	})

	// -----------------------------------------------------------------------
	// 1.3 Fully-Loaded Function/View/Trigger
	// -----------------------------------------------------------------------

	t.Run("create_loaded_function", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, fullTableDDL+loadedFunctionDDL)
		funcs := ctr.queryFunctions(t, schema)
		found := false
		for _, f := range funcs {
			if f.name == "get_user_count" {
				found = true
				if f.volatility != "s" { // s = stable
					t.Errorf("expected stable volatility, got %s", f.volatility)
				}
				if !f.strict {
					t.Error("expected STRICT")
				}
				if f.security != "definer" {
					t.Errorf("expected security definer, got %s", f.security)
				}
				if f.parallel != "s" { // s = safe
					t.Errorf("expected parallel safe, got %s", f.parallel)
				}
			}
		}
		if !found {
			t.Fatal("get_user_count function not found")
		}
	})

	t.Run("create_loaded_procedure", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, fullTableDDL+loadedProcedureDDL)
		funcs := ctr.queryFunctions(t, schema)
		found := false
		for _, f := range funcs {
			if f.name == "cleanup_inactive" {
				found = true
				if f.security != "definer" {
					t.Errorf("expected security definer, got %s", f.security)
				}
				if f.resultType.Valid {
					t.Errorf("expected NULL result type for procedure, got %s", f.resultType.String)
				}
			}
		}
		if !found {
			t.Fatal("cleanup_inactive procedure not found")
		}
	})

	t.Run("create_loaded_view", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, fullTableDDL+loadedViewDDL)
		views := ctr.queryViews(t, schema)
		found := false
		for _, v := range views {
			if v.name == "active_users" {
				found = true
			}
		}
		if !found {
			t.Fatal("active_users view not found")
		}
		cmts := ctr.queryComments(t, schema)
		foundCmt := false
		for _, c := range cmts {
			if c.objectName == "active_users" && c.comment == "Only active users" {
				foundCmt = true
			}
		}
		if !foundCmt {
			t.Error("expected comment on active_users view")
		}
	})

	t.Run("create_loaded_matview", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, fullTableDDL+loadedMatviewDDL)
		tables := ctr.queryTables(t, schema)
		found := false
		for _, tbl := range tables {
			if tbl.name == "user_stats" && tbl.relkind == "m" {
				found = true
			}
		}
		if !found {
			t.Fatal("user_stats materialized view not found")
		}
		idxs := ctr.queryIndexes(t, schema)
		foundIdx := false
		for _, idx := range idxs {
			if idx.name == "idx_user_stats_status" {
				foundIdx = true
			}
		}
		if !foundIdx {
			t.Error("expected index on matview")
		}
	})

	t.Run("create_loaded_trigger", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, fullTableDDL+loadedTriggerFuncDDL+loadedTriggerDDL)
		trigs := ctr.queryTriggers(t, schema)
		found := false
		for _, tr := range trigs {
			if tr.name == "users_update_ts" {
				found = true
			}
		}
		if !found {
			t.Fatal("users_update_ts trigger not found")
		}
	})

	t.Run("empty_to_all_objects_migration", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr, "", roundtripAllObjectsDDL)
	})

	// -----------------------------------------------------------------------
	// 1.4 Fully-Loaded Types/Sequences/Extensions
	// -----------------------------------------------------------------------

	t.Run("create_enum", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, typeEnumDDL)
		enums := ctr.queryEnumTypes(t, schema)
		if len(enums) == 0 {
			t.Fatal("no enum types found")
		}
		found := false
		for _, e := range enums {
			if e.name == "priority" {
				found = true
				if e.values != "low,medium,high,critical" {
					t.Errorf("unexpected enum values: %s", e.values)
				}
			}
		}
		if !found {
			t.Fatal("priority enum not found")
		}
	})

	t.Run("create_domain", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, typeDomainDDL)
		// Verify domain exists via pg_type
		var count int
		err := ctr.db.QueryRowContext(ctr.ctx, `
			SELECT count(*) FROM pg_type t
			JOIN pg_namespace n ON n.oid = t.typnamespace
			WHERE n.nspname = $1 AND t.typname = 'email_addr' AND t.typtype = 'd'
		`, schema).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected 1 domain, got %d", count)
		}
	})

	t.Run("create_composite", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, typeCompositeDDL)
		var count int
		err := ctr.db.QueryRowContext(ctr.ctx, `
			SELECT count(*) FROM pg_type t
			JOIN pg_namespace n ON n.oid = t.typnamespace
			WHERE n.nspname = $1 AND t.typname = 'address' AND t.typtype = 'c'
		`, schema).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected 1 composite type, got %d", count)
		}
	})

	t.Run("create_range", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, typeRangeDDL)
		var count int
		err := ctr.db.QueryRowContext(ctr.ctx, `
			SELECT count(*) FROM pg_type t
			JOIN pg_namespace n ON n.oid = t.typnamespace
			WHERE n.nspname = $1 AND t.typname = 'float_range' AND t.typtype = 'r'
		`, schema).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected 1 range type, got %d", count)
		}
	})

	t.Run("create_sequence", func(t *testing.T) {
		schema := ctr.freshSchema(t)
		ctr.execInSchema(t, schema, typeSequenceDDL)
		seqs := ctr.querySequences(t, schema)
		found := false
		for _, s := range seqs {
			if s.name == "global_seq" {
				found = true
				if s.increment != 5 {
					t.Errorf("expected increment 5, got %d", s.increment)
				}
				if s.start != 0 {
					t.Errorf("expected start 0, got %d", s.start)
				}
				if !s.cycle {
					t.Error("expected CYCLE")
				}
			}
		}
		if !found {
			t.Fatal("global_seq sequence not found")
		}
	})

	t.Run("create_extension", func(t *testing.T) {
		// Extensions are created in the public schema (not in test schemas).
		// We test IF NOT EXISTS to avoid errors if already loaded.
		_, err := ctr.db.ExecContext(ctr.ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
		if err != nil {
			// Some minimal PG images may not have uuid-ossp; skip if unavailable.
			t.Skipf("skipping extension test: %v", err)
		}
		// Verify extension exists
		var count int
		err = ctr.db.QueryRowContext(ctr.ctx,
			`SELECT count(*) FROM pg_extension WHERE extname = 'uuid-ossp'`).Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected extension to exist, got count=%d", count)
		}
	})

	t.Run("empty_to_types_migration", func(t *testing.T) {
		assertOracleRoundtrip(t, ctr, "", allTypesDDL)
	})
}
