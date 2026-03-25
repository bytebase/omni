package catalog

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Base DDL for table change tests
// ---------------------------------------------------------------------------

const tableTestBase = `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`

// ---------------------------------------------------------------------------
// 2.1 Column Changes
// ---------------------------------------------------------------------------

func TestOracleTable_ColumnChanges(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// --- Change column type: varchar(100) → varchar(200) ---
	t.Run("change_varchar_length", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(200) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change column type: integer → bigint ---
	t.Run("change_integer_to_bigint", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age bigint,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change column type: text → user enum (with USING) ---
	t.Run("change_text_to_enum", func(t *testing.T) {
		before := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status text DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add NOT NULL to existing column ---
	t.Run("add_not_null", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text NOT NULL,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop NOT NULL from existing column ---
	t.Run("drop_not_null", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add DEFAULT to column ---
	t.Run("add_default", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text DEFAULT 'unknown@example.com',
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change DEFAULT value ---
	t.Run("change_default", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'unknown',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop DEFAULT from column ---
	t.Run("drop_default", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL,
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add new column with all attributes ---
	t.Run("add_column", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now(),
    bio text NOT NULL DEFAULT ''
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop existing column ---
	t.Run("drop_column", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add GENERATED ALWAYS AS IDENTITY ---
	t.Run("add_identity", func(t *testing.T) {
		// Use a simpler base to avoid conflict with existing PK
		before := `
CREATE TABLE t1 (
    id integer,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer GENERATED ALWAYS AS IDENTITY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop identity from column ---
	t.Run("drop_identity", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer GENERATED ALWAYS AS IDENTITY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change column collation ---
	t.Run("change_collation", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name text COLLATE "C" NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name text COLLATE "POSIX" NOT NULL DEFAULT 'anon'
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add GENERATED ALWAYS AS (expr) STORED column ---
	t.Run("add_generated_stored", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    price numeric(10,2) DEFAULT 0.0,
    tax numeric(10,2) DEFAULT 0.0
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    price numeric(10,2) DEFAULT 0.0,
    tax numeric(10,2) DEFAULT 0.0,
    total numeric GENERATED ALWAYS AS (price + tax) STORED
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop generated column ---
	t.Run("drop_generated_stored", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    price numeric(10,2) DEFAULT 0.0,
    tax numeric(10,2) DEFAULT 0.0,
    total numeric GENERATED ALWAYS AS (price + tax) STORED
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    price numeric(10,2) DEFAULT 0.0,
    tax numeric(10,2) DEFAULT 0.0
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Table UNLOGGED → permanent (persistence change) ---
	t.Run("unlogged_to_permanent", func(t *testing.T) {
		before := `
CREATE UNLOGGED TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})
}

// ---------------------------------------------------------------------------
// 2.2 Constraint Changes
// ---------------------------------------------------------------------------

func TestOracleTable_ConstraintChanges(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// --- Add PRIMARY KEY ---
	t.Run("add_primary_key", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer NOT NULL,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text
);
`
		after := `
CREATE TABLE t1 (
    id integer NOT NULL,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    CONSTRAINT t1_pkey PRIMARY KEY (id)
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop PRIMARY KEY ---
	t.Run("drop_primary_key", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer NOT NULL,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    CONSTRAINT t1_pkey PRIMARY KEY (id)
);
`
		after := `
CREATE TABLE t1 (
    id integer NOT NULL,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add UNIQUE constraint ---
	t.Run("add_unique", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text UNIQUE,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop UNIQUE constraint ---
	t.Run("drop_unique", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text UNIQUE
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add CHECK constraint (with function call) ---
	t.Run("add_check_with_function", func(t *testing.T) {
		before := `
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    age integer
);
`
		after := `
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    age integer,
    CONSTRAINT age_positive CHECK (check_positive(age))
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop CHECK constraint ---
	t.Run("drop_check", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    age integer,
    CONSTRAINT age_positive CHECK (age > 0)
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    age integer
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change CHECK expression ---
	t.Run("change_check_expression", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    age integer,
    CONSTRAINT age_valid CHECK (age > 0)
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    age integer,
    CONSTRAINT age_valid CHECK (age > 0 AND age < 200)
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add FOREIGN KEY ---
	t.Run("add_foreign_key", func(t *testing.T) {
		before := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer
);
`
		after := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id)
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop FOREIGN KEY ---
	t.Run("drop_foreign_key", func(t *testing.T) {
		before := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id)
);
`
		after := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change FK ON DELETE action (CASCADE → SET NULL) ---
	t.Run("change_fk_on_delete", func(t *testing.T) {
		before := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE CASCADE
);
`
		after := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE SET NULL
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add DEFERRABLE constraint ---
	t.Run("add_deferrable", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    CONSTRAINT name_unique UNIQUE (name) DEFERRABLE INITIALLY DEFERRED
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add EXCLUDE constraint ---
	t.Run("add_exclude", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    range_col int4range
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    range_col int4range,
    CONSTRAINT t1_range_excl EXCLUDE USING gist (range_col WITH &&)
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop EXCLUDE constraint ---
	t.Run("drop_exclude", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    range_col int4range,
    CONSTRAINT t1_range_excl EXCLUDE USING gist (range_col WITH &&)
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    range_col int4range
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change FK ON UPDATE action ---
	t.Run("change_fk_on_update", func(t *testing.T) {
		before := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON UPDATE CASCADE
);
`
		after := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON UPDATE SET NULL
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})
}

// ---------------------------------------------------------------------------
// 2.3 Attached Object Changes
// ---------------------------------------------------------------------------

func TestOracleTable_AttachedObjectChanges(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// --- Add standalone index ---
	t.Run("add_index", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
CREATE INDEX idx_email ON t1(email);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop standalone index ---
	t.Run("drop_index", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add partial index (WHERE clause) ---
	t.Run("add_partial_index", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
CREATE INDEX idx_active_score ON t1(score) WHERE active = true;
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add expression index ---
	t.Run("add_expression_index", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
CREATE INDEX idx_lower_email ON t1(lower(email));
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add index with INCLUDE columns ---
	t.Run("add_index_with_include", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
CREATE INDEX idx_name_incl ON t1(name) INCLUDE (email);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add index with DESC/NULLS FIRST ---
	t.Run("add_index_desc_nulls_first", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name);
CREATE INDEX idx_score_desc ON t1(score DESC NULLS FIRST);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add GIN index ---
	t.Run("add_gin_index", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    tags text[],
    data jsonb
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    tags text[],
    data jsonb
);
CREATE INDEX idx_data_gin ON t1 USING gin (data);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change index (columns changed) → DROP + CREATE ---
	t.Run("change_index_columns", func(t *testing.T) {
		before := tableTestBase
		after := `
CREATE TYPE status AS ENUM ('active', 'inactive');
CREATE FUNCTION check_positive(integer) RETURNS boolean LANGUAGE sql IMMUTABLE AS 'SELECT $1 > 0';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text,
    age integer,
    status status DEFAULT 'active',
    score numeric(10,2) DEFAULT 0.0,
    active boolean DEFAULT true,
    created_at timestamptz DEFAULT now()
);
CREATE INDEX idx_name ON t1(name, email);
COMMENT ON TABLE t1 IS 'Test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add trigger with UPDATE OF columns ---
	t.Run("add_trigger", func(t *testing.T) {
		before := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text
);
`
		after := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text
);
CREATE TRIGGER t1_before_update
    BEFORE UPDATE ON t1
    FOR EACH ROW
    EXECUTE FUNCTION trig_fn();
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop trigger ---
	t.Run("drop_trigger", func(t *testing.T) {
		before := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text
);
CREATE TRIGGER t1_before_update
    BEFORE UPDATE ON t1
    FOR EACH ROW
    EXECUTE FUNCTION trig_fn();
`
		after := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    email text
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change trigger (WHEN clause changed) → DROP + CREATE ---
	t.Run("change_trigger_when", func(t *testing.T) {
		before := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
CREATE TRIGGER t1_before_update
    BEFORE UPDATE ON t1
    FOR EACH ROW
    EXECUTE FUNCTION trig_fn();
`
		after := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
CREATE TRIGGER t1_before_update
    BEFORE UPDATE ON t1
    FOR EACH ROW
    WHEN (OLD.name IS DISTINCT FROM NEW.name)
    EXECUTE FUNCTION trig_fn();
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add RLS policy with USING and WITH CHECK ---
	t.Run("add_policy", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_select ON t1 FOR SELECT USING (active = true);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop policy ---
	t.Run("drop_policy", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_select ON t1 FOR SELECT USING (active = true);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Enable ROW LEVEL SECURITY ---
	t.Run("enable_rls", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Disable ROW LEVEL SECURITY ---
	t.Run("disable_rls", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    active boolean DEFAULT true
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add table comment ---
	t.Run("add_table_comment", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
COMMENT ON TABLE t1 IS 'A test table';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Add column comment ---
	t.Run("add_column_comment", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
COMMENT ON COLUMN t1.name IS 'Display name';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Change table comment ---
	t.Run("change_table_comment", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
COMMENT ON TABLE t1 IS 'Old comment';
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
COMMENT ON TABLE t1 IS 'New comment';
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Drop table comment (set NULL) ---
	t.Run("drop_table_comment", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
COMMENT ON TABLE t1 IS 'A comment';
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- ALTER TABLE REPLICA IDENTITY FULL ---
	t.Run("replica_identity_full", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
ALTER TABLE t1 REPLICA IDENTITY FULL;
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- FK ON DELETE RESTRICT ---
	t.Run("fk_on_delete_restrict", func(t *testing.T) {
		before := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE CASCADE
);
`
		after := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE RESTRICT
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- FK ON DELETE SET DEFAULT ---
	t.Run("fk_on_delete_set_default", func(t *testing.T) {
		before := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer DEFAULT 1,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id)
);
`
		after := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer DEFAULT 1,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON DELETE SET DEFAULT
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- FK ON UPDATE RESTRICT ---
	t.Run("fk_on_update_restrict", func(t *testing.T) {
		before := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON UPDATE CASCADE
);
`
		after := `
CREATE TABLE parent (
    id integer PRIMARY KEY
);
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    parent_id integer,
    CONSTRAINT t1_parent_fk FOREIGN KEY (parent_id) REFERENCES parent(id) ON UPDATE RESTRICT
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- DEFERRABLE INITIALLY IMMEDIATE ---
	t.Run("deferrable_initially_immediate", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    CONSTRAINT name_unique UNIQUE (name)
);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon',
    CONSTRAINT name_unique UNIQUE (name) DEFERRABLE INITIALLY IMMEDIATE
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- GENERATED BY DEFAULT AS IDENTITY ---
	t.Run("add_identity_by_default", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE TABLE t1 (
    id integer GENERATED BY DEFAULT AS IDENTITY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})
}

// ---------------------------------------------------------------------------
// 3.1 Trigger Coverage Gaps
// ---------------------------------------------------------------------------

func TestOracleTable_TriggerGaps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// --- Trigger BEFORE DELETE ---
	t.Run("trigger_before_delete", func(t *testing.T) {
		before := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN OLD; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN OLD; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
CREATE TRIGGER t1_before_delete
    BEFORE DELETE ON t1
    FOR EACH ROW
    EXECUTE FUNCTION trig_fn();
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Trigger FOR EACH STATEMENT ---
	t.Run("trigger_for_each_statement", func(t *testing.T) {
		before := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NULL; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NULL; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
CREATE TRIGGER t1_after_insert_stmt
    AFTER INSERT ON t1
    FOR EACH STATEMENT
    EXECUTE FUNCTION trig_fn();
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- INSTEAD OF trigger on view ---
	t.Run("trigger_instead_of_on_view", func(t *testing.T) {
		before := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
CREATE VIEW v1 AS SELECT id, name FROM t1;
`
		after := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
CREATE VIEW v1 AS SELECT id, name FROM t1;
CREATE TRIGGER v1_instead_of_insert
    INSTEAD OF INSERT ON v1
    FOR EACH ROW
    EXECUTE FUNCTION trig_fn();
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Trigger with REFERENCING OLD TABLE AS ---
	t.Run("trigger_referencing_old_table", func(t *testing.T) {
		before := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NULL; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
`
		after := `
CREATE FUNCTION trig_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NULL; END';
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name varchar(100) NOT NULL DEFAULT 'anon'
);
CREATE TRIGGER t1_after_delete_ref
    AFTER DELETE ON t1
    REFERENCING OLD TABLE AS old_rows
    FOR EACH STATEMENT
    EXECUTE FUNCTION trig_fn();
`
		assertOracleRoundtrip(t, oracle, before, after)
	})
}

// ---------------------------------------------------------------------------
// 3.2 Policy Coverage Gaps
// ---------------------------------------------------------------------------

func TestOracleTable_PolicyGaps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	// --- Policy AS RESTRICTIVE ---
	t.Run("policy_restrictive", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_restrict ON t1 AS RESTRICTIVE FOR SELECT USING (true);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Policy FOR INSERT ---
	t.Run("policy_for_insert", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    col integer DEFAULT 0
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    col integer DEFAULT 0
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_ins ON t1 FOR INSERT WITH CHECK (col > 0);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Policy FOR UPDATE ---
	t.Run("policy_for_update", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    col integer DEFAULT 0
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    col integer DEFAULT 0
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_upd ON t1 FOR UPDATE USING (true) WITH CHECK (col > 0);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Policy FOR DELETE ---
	t.Run("policy_for_delete", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    col integer DEFAULT 0
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    col integer DEFAULT 0
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_del ON t1 FOR DELETE USING (col > 0);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Policy WITH CHECK expression (complex) ---
	t.Run("policy_with_check_expression", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    col integer DEFAULT 0,
    name text
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    col integer DEFAULT 0,
    name text
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_validate ON t1 FOR INSERT WITH CHECK (col > 0 AND name IS NOT NULL);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	// --- Policy with specific roles (PUBLIC) ---
	t.Run("policy_to_public", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_pub ON t1 FOR SELECT TO PUBLIC USING (true);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})
}

// ---------------------------------------------------------------------------
// Group 2: Policy ALTER Variants
// ---------------------------------------------------------------------------

func TestOracleTable_PolicyAlterVariants(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	t.Run("alter_policy_change_using_expression", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true,
    verified boolean DEFAULT false
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_sel ON t1 FOR SELECT USING (active);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true,
    verified boolean DEFAULT false
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_sel ON t1 FOR SELECT USING (active AND verified);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("alter_policy_change_with_check", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    score integer DEFAULT 0
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_ins ON t1 FOR INSERT WITH CHECK (true);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    score integer DEFAULT 0
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_ins ON t1 FOR INSERT WITH CHECK (score > 0);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("alter_policy_change_command_type", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_pol ON t1 FOR SELECT USING (active);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_pol ON t1 FOR ALL USING (active);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("alter_policy_change_using_expression_different", func(t *testing.T) {
		// Tests altering USING expression without changing roles
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_sel ON t1 FOR SELECT TO PUBLIC USING (true);
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
CREATE POLICY t1_sel ON t1 FOR SELECT TO PUBLIC USING (active);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("force_row_level_security", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
ALTER TABLE t1 FORCE ROW LEVEL SECURITY;
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("no_force_row_level_security", func(t *testing.T) {
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
ALTER TABLE t1 FORCE ROW LEVEL SECURITY;
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    active boolean DEFAULT true
);
ALTER TABLE t1 ENABLE ROW LEVEL SECURITY;
`
		assertOracleRoundtrip(t, oracle, before, after)
	})
}

// ---------------------------------------------------------------------------
// Group 5: Constraint Trigger and TABLE LIKE
// ---------------------------------------------------------------------------

func TestOracleTable_ConstraintTriggerAndLike(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping oracle test: requires Docker")
	}
	oracle := startPGOracle(t)

	t.Run("constraint_trigger_deferrable", func(t *testing.T) {
		// BUG: Constraint triggers are not properly supported by the migration
		// engine. It generates "ALTER TABLE ADD CONSTRAINT ... /* unsupported type t */"
		// instead of CREATE CONSTRAINT TRIGGER ... DEFERRABLE INITIALLY DEFERRED.
		before := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    val integer
);
CREATE FUNCTION check_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
`
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    val integer
);
CREATE FUNCTION check_fn() RETURNS trigger LANGUAGE plpgsql AS 'BEGIN RETURN NEW; END';
CREATE CONSTRAINT TRIGGER t1_check AFTER INSERT ON t1
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW EXECUTE FUNCTION check_fn();
`
		assertOracleRoundtrip(t, oracle, before, after)
	})

	t.Run("table_like_including_all", func(t *testing.T) {
		// BUG: TABLE LIKE INCLUDING ALL is CREATE TABLE syntax sugar. LoadSQL
		// doesn't expand it, so constraints from t1 don't carry over to t2
		// in the catalog. The migration creates t2 without constraints.
		before := ``
		after := `
CREATE TABLE t1 (
    id integer PRIMARY KEY,
    name text NOT NULL
);
CREATE TABLE t2 (LIKE t1 INCLUDING ALL);
`
		assertOracleRoundtrip(t, oracle, before, after)
	})
}
