package catalog

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestSimulationPG18 holds the catalog to testdata/pg18/simulate.tsv:
// DDL statements across the core one_table and objects fixtures plus curated
// PG 18 DDL/builtin feature cases, recorded against PostgreSQL 18 (the pg18
// pin in pg/testdata/pins.tsv) and replayed offline from
// testdata/tapes/TestSimulationPG18.tape.
func TestSimulationPG18(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), pg18SimulationFixtures, simulated)
	if err != nil {
		t.Fatal(err)
	}
	rows = append(rows, pg18FeatureRows...)
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "pg18", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.PG18,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		MadeUp:   nil,
		Omni:     simulationOmniWithVersion(PG18),
		Engine:   simulationEngineWithRole(testcontainer.PG18, simulationFormsQuery),
	})
}

// pg18SimulationFixtures covers the core table/column/index/object families
// where PG 17 vs PG 18 catalog semantics differ (partitioned tables, regular
// tables, foreign keys, views, materialized views, sequences, and domains)
// while keeping the offline tape size compact.
var pg18SimulationFixtures = []string{
	"one_table@missing",
	"one_table@table",
	"one_table@column",
	"objects@view",
	"objects@matview",
	"objects@sequence",
	"objects@partitioned",
	"objects@fk",
	"objects@domain",
}

// pg18FeatureRows exercises PostgreSQL 18-specific DDL behavior and catalog
// features beyond the generated compat.tsv statements (partitioned persistence
// and attach/detach, named NOT NULL constraints, domain NOT NULL NO INHERIT,
// generated expressions, identity on partitioned tables, NULLS NOT DISTINCT,
// SQL/JSON constructors, and PG 18 built-in functions in DDL).
var pg18FeatureRows = []simulate.Generated{
	// PG 18 domain NOT NULL NO INHERIT rejection (named and unnamed).
	{Fixture: "objects@domain", Rule: "pg18/domain_named_notnull_noinherit", Statement: "ALTER DOMAIN a ADD CONSTRAINT a_nn NOT NULL NO INHERIT"},
	{Fixture: "one_table@missing", Rule: "pg18/domain_named_notnull_noinherit", Statement: "ALTER DOMAIN a ADD CONSTRAINT a_nn NOT NULL NO INHERIT"},
	{Fixture: "one_table@missing", Rule: "pg18/create_domain_named_notnull", Statement: "CREATE DOMAIN d_nn AS int CONSTRAINT d_nn_c NOT NULL"},
	{Fixture: "objects@domain", Rule: "pg18/domain_add_check_not_valid", Statement: "ALTER DOMAIN a ADD CONSTRAINT c_nv CHECK (VALUE > 0) NOT VALID"},

	// PG 18 partitioned table DDL and persistence rules.
	{Fixture: "objects@partitioned", Rule: "pg18/detach_partition", Statement: "ALTER TABLE a DETACH PARTITION a_p1"},
	{Fixture: "objects@partitioned", Rule: "pg18/partition_add_identity", Statement: "ALTER TABLE a ALTER d ADD GENERATED ALWAYS AS IDENTITY"},
	{Fixture: "objects@partitioned", Rule: "pg18/partition_set_not_null", Statement: "ALTER TABLE a ALTER d SET NOT NULL"},
	{Fixture: "one_table@missing", Rule: "pg18/create_partitioned_identity", Statement: "CREATE TABLE part_id (id int GENERATED ALWAYS AS IDENTITY, k int) PARTITION BY RANGE (k)"},

	// PG 18 table NOT NULL constraints, inheritance, and LIKE.
	{Fixture: "one_table@missing", Rule: "pg18/create_table_named_notnull", Statement: "CREATE TABLE nn_t (id int CONSTRAINT nn_id NOT NULL, v int NOT NULL)"},
	{Fixture: "one_table@column", Rule: "pg18/create_table_inherits", Statement: "CREATE TABLE child () INHERITS (a)"},
	{Fixture: "one_table@column", Rule: "pg18/drop_pk_not_null", Statement: "ALTER TABLE a ALTER b DROP NOT NULL"},
	{Fixture: "one_table@column", Rule: "pg18/set_col_not_null", Statement: "ALTER TABLE a ALTER a SET NOT NULL"},
	{Fixture: "one_table@column", Rule: "pg18/drop_pk_not_null_constraint", Statement: "ALTER TABLE a DROP CONSTRAINT a_b_not_null"},
	{Fixture: "one_table@column", Rule: "pg18/rename_pk_not_null_constraint", Statement: "ALTER TABLE a RENAME CONSTRAINT a_b_not_null TO pk_b_nn"},
	{Fixture: "one_table@column", Rule: "pg18/comment_pk_not_null_constraint", Statement: "COMMENT ON CONSTRAINT a_b_not_null ON a IS 'pk not null'"},
	{Fixture: "one_table@column", Rule: "pg18/create_table_like_all", Statement: "CREATE TABLE nn_like (LIKE a INCLUDING ALL)"},
	{Fixture: "objects@fk", Rule: "pg18/create_table_like_all", Statement: "CREATE TABLE nn_like (LIKE a INCLUDING ALL)"},

	// Generated columns and SET/DROP EXPRESSION.
	{Fixture: "one_table@missing", Rule: "pg18/create_table_generated_stored", Statement: "CREATE TABLE gen_t (x int, y int GENERATED ALWAYS AS (x * 2) STORED)"},
	{Fixture: "one_table@column", Rule: "pg18/alter_add_generated_stored", Statement: "ALTER TABLE a ADD COLUMN g int GENERATED ALWAYS AS (b + 1) STORED"},
	{Fixture: "one_table@column", Rule: "pg18/alter_set_expression_non_gen", Statement: "ALTER TABLE a ALTER COLUMN a SET EXPRESSION AS (b + 1)"},
	{Fixture: "one_table@column", Rule: "pg18/alter_drop_expression_non_gen", Statement: "ALTER TABLE a ALTER COLUMN a DROP EXPRESSION"},
	{Fixture: "one_table@column", Rule: "pg18/alter_drop_expression_if_exists", Statement: "ALTER TABLE a ALTER COLUMN a DROP EXPRESSION IF EXISTS"},

	// PG 18 built-in functions in DDL defaults, checks, generated columns, indexes, and views.
	{Fixture: "one_table@missing", Rule: "pg18/builtin_uuidv7_casefold_table", Statement: "CREATE TABLE u7 (id uuid PRIMARY KEY DEFAULT uuidv7(), name text CHECK (casefold(name) <> ''))"},
	{Fixture: "one_table@missing", Rule: "pg18/builtin_pg18_view", Statement: "CREATE VIEW v_pg18 AS SELECT uuidv7() AS id, casefold('Hello') AS cf, array_sort(ARRAY[3, 1, 2]) AS sorted"},
	{Fixture: "one_table@missing", Rule: "pg18/builtin_array_sort_generated", Statement: "CREATE TABLE t_arr (vals int[], sorted int[] GENERATED ALWAYS AS (array_sort(vals)) STORED)"},
	{Fixture: "one_table@column", Rule: "pg18/builtin_casefold_index", Statement: "CREATE INDEX idx_casefold ON a (casefold(b::text))"},

	// SQL/JSON, NULLS NOT DISTINCT, UNLOGGED SEQUENCE, SECURITY INVOKER views, SQL body functions, FK SET NULL column list.
	{Fixture: "one_table@missing", Rule: "pg18/sql_json_view", Statement: "CREATE VIEW v_json AS SELECT JSON_VALUE('{\"k\": 1}'::jsonb, '$.\"k\"' RETURNING int) AS v, JSON_EXISTS('{\"k\": 1}'::jsonb, '$.\"k\"') AS e, JSON_QUERY('{\"k\": [1, 2]}'::jsonb, '$.\"k\"') AS q"},
	{Fixture: "one_table@missing", Rule: "pg18/unique_nulls_not_distinct", Statement: "CREATE TABLE t_nnd (id int PRIMARY KEY, k int, UNIQUE NULLS NOT DISTINCT (k))"},
	{Fixture: "one_table@column", Rule: "pg18/index_nulls_not_distinct", Statement: "CREATE UNIQUE INDEX idx_nnd ON a (d) NULLS NOT DISTINCT"},
	{Fixture: "one_table@missing", Rule: "pg18/unlogged_sequence", Statement: "CREATE UNLOGGED SEQUENCE s_unlog START WITH 10 INCREMENT BY 2"},
	{Fixture: "one_table@column", Rule: "pg18/security_invoker_view", Statement: "CREATE VIEW v_inv WITH (security_invoker = true, security_barrier = true) AS SELECT b, a FROM a"},
	{Fixture: "one_table@missing", Rule: "pg18/sql_body_function", Statement: "CREATE FUNCTION f_atomic(x int) RETURNS int LANGUAGE sql IMMUTABLE RETURN x + 1"},
	{Fixture: "one_table@column", Rule: "pg18/fk_on_delete_set_null_cols", Statement: "CREATE TABLE t_fk_setnull (id int PRIMARY KEY, ab int, d int, FOREIGN KEY (ab) REFERENCES a (b) ON DELETE SET NULL (ab))"},
}

func execCodeVersion(t *testing.T, v Version, setupSQL, stmt string) (string, *Catalog) {
	t.Helper()
	c, err := LoadSQLWithVersion(setupSQL, v)
	if err != nil {
		t.Fatalf("LoadSQLWithVersion(%d, %q): %v", v, setupSQL, err)
	}
	results, err := c.Exec(stmt, nil)
	if err != nil {
		var ce *Error
		if errors.As(err, &ce) {
			return ce.Code, c
		}
		t.Fatalf("Exec(%q): unexpected non-catalog error %v", stmt, err)
	}
	for _, r := range results {
		if r.Error != nil {
			var ce *Error
			if errors.As(r.Error, &ce) {
				return ce.Code, c
			}
			t.Fatalf("Exec(%q): unexpected non-catalog error %v", stmt, r.Error)
		}
	}
	return "", c
}

func TestPG18CatalogVersionSemantics(t *testing.T) {
	t.Run("version_default_and_clone", func(t *testing.T) {
		c := New()
		if c.Version() != PG17 {
			t.Fatalf("New().Version() = %d, want %d", c.Version(), PG17)
		}
		c18 := NewWithVersion(PG18)
		if c18.Version() != PG18 {
			t.Fatalf("NewWithVersion(PG18).Version() = %d, want %d", c18.Version(), PG18)
		}
		if cloned := c18.Clone(); cloned.Version() != PG18 {
			t.Fatalf("Clone().Version() = %d, want %d", cloned.Version(), PG18)
		}
		c18.SetVersion(0)
		if c18.Version() != PG17 {
			t.Fatalf("SetVersion(0).Version() = %d, want %d", c18.Version(), PG17)
		}
	})

	plainSetup := `CREATE TABLE a (b int PRIMARY KEY, a int, d int);`
	partSetup := `
CREATE TABLE a (b int, a int, d int, PRIMARY KEY (b, a)) PARTITION BY RANGE (a);
CREATE TABLE a_p1 PARTITION OF a FOR VALUES FROM (0) TO (100);
`
	domainSetup := `CREATE DOMAIN a AS int DEFAULT 1 NOT NULL;`

	cases := []struct {
		name   string
		setup  string
		stmt   string
		want17 string
		want18 string
	}{
		{"domain_add_notnull_noinherit_missing", "", "ALTER DOMAIN a ADD NOT NULL NO INHERIT", CodeUndefinedObject, CodeFeatureNotSupported},
		{"domain_add_notnull_noinherit_domain", domainSetup, "ALTER DOMAIN a ADD NOT NULL NO INHERIT", "", CodeFeatureNotSupported},
		{"domain_add_named_notnull_noinherit", domainSetup, "ALTER DOMAIN a ADD CONSTRAINT a_nn NOT NULL NO INHERIT", "", CodeFeatureNotSupported},
		{"attach_partition_on_plain_table", plainSetup, "ALTER TABLE a ATTACH PARTITION a DEFAULT", CodeInvalidObjectDefinition, CodeWrongObjectType},
		{"detach_partition_on_plain_table", plainSetup, "ALTER TABLE a DETACH PARTITION a", CodeInvalidObjectDefinition, CodeWrongObjectType},
		{"detach_partition_finalize_on_plain_table", plainSetup, "ALTER TABLE a DETACH PARTITION a FINALIZE", CodeUndefinedTable, CodeWrongObjectType},
		{"set_logged_on_partitioned", partSetup, "ALTER TABLE a SET LOGGED", "", CodeWrongObjectType},
		{"set_unlogged_on_partitioned", partSetup, "ALTER TABLE a SET UNLOGGED", "", CodeWrongObjectType},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got17, _ := execCodeVersion(t, PG17, tc.setup, tc.stmt); got17 != tc.want17 {
				t.Errorf("PG17 %s: got code %q, want %q", tc.stmt, got17, tc.want17)
			}
			if got18, _ := execCodeVersion(t, PG18, tc.setup, tc.stmt); got18 != tc.want18 {
				t.Errorf("PG18 %s: got code %q, want %q", tc.stmt, got18, tc.want18)
			}
		})
	}
}

func TestPG18NotNullConstraints(t *testing.T) {
	setup := `
CREATE TABLE a (
	b int PRIMARY KEY,
	c int NOT NULL,
	d int CONSTRAINT d_nn NOT NULL,
	e int GENERATED ALWAYS AS IDENTITY
);
`
	// PG 17 does not create cataloged NOT NULL constraints.
	if got, _ := execCodeVersion(t, PG17, setup, "ALTER TABLE a DROP CONSTRAINT a_c_not_null"); got != CodeUndefinedObject {
		t.Fatalf("PG17 DROP CONSTRAINT a_c_not_null: got %q, want %q", got, CodeUndefinedObject)
	}

	// PG 18 creates cataloged NOT NULL constraints for PK, NOT NULL, named NOT NULL, and IDENTITY columns.
	if got, _ := execCodeVersion(t, PG18, setup, "ALTER TABLE a DROP CONSTRAINT a_b_not_null"); got != CodeInvalidTableDefinition {
		t.Fatalf("PG18 DROP CONSTRAINT on PK NOT NULL: got %q, want %q", got, CodeInvalidTableDefinition)
	}
	if got, _ := execCodeVersion(t, PG18, setup, "ALTER TABLE a DROP CONSTRAINT a_e_not_null"); got != CodeInvalidObjectDefinition {
		t.Fatalf("PG18 DROP CONSTRAINT on identity NOT NULL: got %q, want %q", got, CodeInvalidObjectDefinition)
	}

	// Dropping a normal NOT NULL constraint clears NotNull and removes the constraint.
	if got, c := execCodeVersion(t, PG18, setup, "ALTER TABLE a DROP CONSTRAINT a_c_not_null"); got != "" {
		t.Fatalf("PG18 DROP CONSTRAINT a_c_not_null: got %q, want ok", got)
	} else {
		rel := c.GetRelation("public", "a")
		if rel.Columns[1].NotNull {
			t.Fatalf("expected column c to become nullable after DROP CONSTRAINT a_c_not_null")
		}
		if res, err := c.Exec("ALTER TABLE a DROP CONSTRAINT a_c_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) == 0 || res[0].Error == nil {
			t.Fatalf("expected second DROP CONSTRAINT a_c_not_null to fail")
		}
	}

	// Explicit name, RENAME CONSTRAINT, and COMMENT ON CONSTRAINT all work in PG 18.
	if got, c := execCodeVersion(t, PG18, setup, "ALTER TABLE a RENAME CONSTRAINT d_nn TO d_not_null_renamed"); got != "" {
		t.Fatalf("PG18 RENAME CONSTRAINT d_nn: got %q, want ok", got)
	} else if res, err := c.Exec("COMMENT ON CONSTRAINT d_not_null_renamed ON a IS 'renamed nn'", nil); err != nil {
		t.Fatal(err)
	} else if len(res) > 0 && res[0].Error != nil {
		t.Fatalf("PG18 COMMENT ON CONSTRAINT d_not_null_renamed: %v", res[0].Error)
	}

	// SET NOT NULL creates a cataloged NOT NULL constraint; DROP NOT NULL removes it.
	if got, c := execCodeVersion(t, PG18, "CREATE TABLE t (x int);", "ALTER TABLE t ALTER x SET NOT NULL"); got != "" {
		t.Fatalf("PG18 SET NOT NULL: got %q, want ok", got)
	} else {
		if res, err := c.Exec("COMMENT ON CONSTRAINT t_x_not_null ON t IS 'added nn'", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 COMMENT ON CONSTRAINT t_x_not_null: %v", res[0].Error)
		}
		if res, err := c.Exec("ALTER TABLE t ALTER x DROP NOT NULL", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP NOT NULL: %v", res[0].Error)
		}
		if res, err := c.Exec("ALTER TABLE t DROP CONSTRAINT t_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) == 0 || res[0].Error == nil {
			t.Fatalf("expected t_x_not_null to be removed by DROP NOT NULL")
		}
	}
}

func TestPG18PartitionDetachPKNameAndJsonExpr(t *testing.T) {
	// Partition cloned PK constraint uses <part>_pkey (not <part>_<parent>_pkey).
	partSetup := `
CREATE TABLE a (b int, a int, d int, PRIMARY KEY (b, a)) PARTITION BY RANGE (a);
CREATE TABLE a_p1 PARTITION OF a FOR VALUES FROM (0) TO (100);
`
	if got, c := execCodeVersion(t, PG18, partSetup, "ALTER TABLE a DETACH PARTITION a_p1"); got != "" {
		t.Fatalf("DETACH PARTITION a_p1: got %q, want ok", got)
	} else {
		p1 := c.GetRelation("public", "a_p1")
		pk := c.findPKConstraint(p1.OID)
		if pk == nil || pk.Name != "a_p1_pkey" {
			t.Fatalf("detached partition PK constraint = %+v, want Name=%q", pk, "a_p1_pkey")
		}
	}

	// JSON_VALUE RETURNING resolves column type and deparses RETURNING; JSON_QUERY deparses RETURNING/WRAPPER/QUOTES.
	c, err := LoadSQLWithVersion(`CREATE VIEW v_json AS SELECT
		JSON_VALUE('{"k": 1}'::jsonb, '$."k"' RETURNING int) AS v,
		JSON_QUERY('{"k": [1, 2]}'::jsonb, '$."k"') AS q,
		JSON_QUERY('{"k": [1, 2]}'::jsonb, '$."k"' RETURNING text WITH CONDITIONAL WRAPPER OMIT QUOTES) AS q2`, PG18)
	if err != nil {
		t.Fatalf("LoadSQLWithVersion: %v", err)
	}
	rel := c.GetRelation("public", "v_json")
	if rel == nil || len(rel.Columns) != 3 {
		t.Fatalf("expected 3 view columns, got %+v", rel)
	}
	if rel.Columns[0].TypeOID != INT4OID {
		t.Errorf("JSON_VALUE RETURNING int column type = %d, want %d (INT4OID)", rel.Columns[0].TypeOID, INT4OID)
	}
	if rel.Columns[1].TypeOID != JSONBOID {
		t.Errorf("JSON_QUERY default column type = %d, want %d (JSONBOID)", rel.Columns[1].TypeOID, JSONBOID)
	}
	if rel.Columns[2].TypeOID != TEXTOID {
		t.Errorf("JSON_QUERY RETURNING text column type = %d, want %d (TEXTOID)", rel.Columns[2].TypeOID, TEXTOID)
	}
	def, err := c.GetViewDefinition("public", "v_json")
	if err != nil {
		t.Fatalf("GetViewDefinition: %v", err)
	}
	for _, sub := range []string{
		"RETURNING integer",
		"RETURNING jsonb WITHOUT WRAPPER KEEP QUOTES",
		"RETURNING text WITH CONDITIONAL WRAPPER OMIT QUOTES",
	} {
		if !strings.Contains(def, sub) {
			t.Errorf("view def %q missing %q", def, sub)
		}
	}
}
