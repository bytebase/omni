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

	// Explicit CONSTRAINT <name> NOT NULL on a PRIMARY KEY column preserves the explicit name.
	if got, c := execCodeVersion(t, PG18, `CREATE TABLE t_pk_nn (id int CONSTRAINT id_nn NOT NULL PRIMARY KEY);`,
		"COMMENT ON CONSTRAINT id_nn ON t_pk_nn IS 'explicit pk nn'"); got != "" {
		t.Fatalf("PG18 explicit NOT NULL on PK column: got %q, want ok", got)
	} else if res, err := c.Exec("ALTER TABLE t_pk_nn RENAME CONSTRAINT id_nn TO id_nn_renamed", nil); err != nil {
		t.Fatal(err)
	} else if len(res) > 0 && res[0].Error != nil {
		t.Fatalf("PG18 RENAME CONSTRAINT id_nn: %v", res[0].Error)
	}

	// Dropping a NOT NULL constraint on a partitioned table propagates to descendant partitions.
	partNNSetup := `
CREATE TABLE p (x int NOT NULL, k int) PARTITION BY RANGE (k);
CREATE TABLE p1 PARTITION OF p FOR VALUES FROM (0) TO (100);
`
	if got, c := execCodeVersion(t, PG18, partNNSetup, "ALTER TABLE p DROP CONSTRAINT p_x_not_null"); got != "" {
		t.Fatalf("PG18 DROP CONSTRAINT p_x_not_null on partitioned table: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p")
		p1 := c.GetRelation("public", "p1")
		if p.Columns[0].NotNull {
			t.Fatalf("expected parent column x to become nullable")
		}
		if p1.Columns[0].NotNull {
			t.Fatalf("expected partition column x to become nullable")
		}
		for _, con := range c.ConstraintsOf(p1.OID) {
			if con.Type == ConstraintNotNull {
				t.Fatalf("expected partition NOT NULL constraint to be dropped, still have %q", con.Name)
			}
		}
	}

	// Failed ADD COLUMN with duplicate explicit NOT NULL constraint name leaves the table unchanged.
	if got, c := execCodeVersion(t, PG18, `CREATE TABLE t_dup_nn (a int CONSTRAINT nn NOT NULL);`,
		"ALTER TABLE t_dup_nn ADD COLUMN b int CONSTRAINT nn NOT NULL"); got != CodeDuplicateObject {
		t.Fatalf("PG18 ADD COLUMN duplicate NOT NULL name: got %q, want %q", got, CodeDuplicateObject)
	} else {
		rel := c.GetRelation("public", "t_dup_nn")
		if len(rel.Columns) != 1 {
			t.Fatalf("expected failed ADD COLUMN to leave 1 column, got %d", len(rel.Columns))
		}
	}

	// Failed ADD COLUMN default analysis rolls back the registered NOT NULL constraint so retry succeeds.
	if got, c := execCodeVersion(t, PG18, `CREATE TABLE t_rb_nn (a int);`,
		"ALTER TABLE t_rb_nn ADD COLUMN b int CONSTRAINT nn NOT NULL DEFAULT no_such_function()"); got != CodeUndefinedFunction {
		t.Fatalf("PG18 ADD COLUMN bad default: got %q, want %q", got, CodeUndefinedFunction)
	} else if res, err := c.Exec("ALTER TABLE t_rb_nn ADD COLUMN b int CONSTRAINT nn NOT NULL", nil); err != nil {
		t.Fatal(err)
	} else if len(res) > 0 && res[0].Error != nil {
		t.Fatalf("PG18 retry ADD COLUMN after rollback failed: %v", res[0].Error)
	}

	// Generated NOT NULL constraint names stay within 63 bytes (NAMEDATALEN-1) and remain addressable via Exec.
	longTbl := strings.Repeat("t", 40)
	longCol := strings.Repeat("c", 30) + "é"
	longCreate := "CREATE TABLE \"" + longTbl + "\" (\"" + longCol + "\" int NOT NULL);"
	if c, err := LoadSQLWithVersion(longCreate, PG18); err != nil {
		t.Fatalf("LoadSQLWithVersion long ident: %v", err)
	} else {
		rel := c.GetRelation("public", longTbl)
		var nn *Constraint
		for _, con := range c.ConstraintsOf(rel.OID) {
			if con.Type == ConstraintNotNull {
				nn = con
			}
		}
		if nn == nil {
			t.Fatalf("expected NOT NULL constraint on %q", longTbl)
		}
		if len(nn.Name) > 63 || !strings.HasSuffix(nn.Name, "_not_null") {
			t.Fatalf("generated NOT NULL constraint name %q (len=%d): want <=63 bytes ending in _not_null", nn.Name, len(nn.Name))
		}
		if res, err := c.Exec(`COMMENT ON CONSTRAINT "`+nn.Name+`" ON "`+longTbl+`" IS 'ok'`, nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("COMMENT ON CONSTRAINT %q failed: %v", nn.Name, res[0].Error)
		}
	}

	// ALTER TABLE <partitioned> ADD PRIMARY KEY registers cataloged NOT NULL constraints on partitions too.
	partNullableSetup := `
CREATE TABLE p_pk (id int, k int) PARTITION BY RANGE (k);
CREATE TABLE p_pk_1 PARTITION OF p_pk FOR VALUES FROM (0) TO (100);
`
	if got, c := execCodeVersion(t, PG18, partNullableSetup, "ALTER TABLE p_pk ADD PRIMARY KEY (id, k)"); got != "" {
		t.Fatalf("PG18 ADD PRIMARY KEY on partitioned table: got %q, want ok", got)
	} else {
		if res, err := c.Exec("COMMENT ON CONSTRAINT p_pk_1_id_not_null ON p_pk_1 IS 'part pk nn'", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 COMMENT ON CONSTRAINT p_pk_1_id_not_null: %v", res[0].Error)
		}
		if res, err := c.Exec("ALTER TABLE p_pk_1 RENAME CONSTRAINT p_pk_1_k_not_null TO p_pk_1_k_nn", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 RENAME CONSTRAINT p_pk_1_k_not_null: %v", res[0].Error)
		}
	}

	// An explicit PK named <rel>_<col>_not_null reserves that name; the implicit NOT NULL picks <rel>_<col>_not_null1.
	if got, c := execCodeVersion(t, PG18, `CREATE TABLE t_pk_res (id int);`,
		"ALTER TABLE t_pk_res ADD CONSTRAINT t_pk_res_id_not_null PRIMARY KEY (id)"); got != "" {
		t.Fatalf("PG18 ADD CONSTRAINT t_pk_res_id_not_null PRIMARY KEY: got %q, want ok", got)
	} else {
		rel := c.GetRelation("public", "t_pk_res")
		var nnName string
		for _, con := range c.ConstraintsOf(rel.OID) {
			if con.Type == ConstraintNotNull {
				nnName = con.Name
			}
		}
		if nnName != "t_pk_res_id_not_null1" {
			t.Fatalf("implicit NOT NULL name = %q, want %q", nnName, "t_pk_res_id_not_null1")
		}
		if res, err := c.Exec("ALTER TABLE t_pk_res DROP CONSTRAINT t_pk_res_id_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP CONSTRAINT t_pk_res_id_not_null (PK): %v", res[0].Error)
		}
	}

	// Dropping a NOT NULL constraint on an ordinary inheritance parent propagates to children by column name.
	inhNNSetup := `
CREATE TABLE p_inh (x int NOT NULL);
CREATE TABLE c_inh () INHERITS (p_inh);
`
	if got, c := execCodeVersion(t, PG18, inhNNSetup, "ALTER TABLE p_inh DROP CONSTRAINT p_inh_x_not_null"); got != "" {
		t.Fatalf("PG18 DROP CONSTRAINT p_inh_x_not_null on inheritance parent: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_inh")
		ch := c.GetRelation("public", "c_inh")
		if p.Columns[0].NotNull || ch.Columns[0].NotNull {
			t.Fatalf("expected both parent and child x to become nullable (p=%v, ch=%v)", p.Columns[0].NotNull, ch.Columns[0].NotNull)
		}
		for _, con := range c.ConstraintsOf(ch.OID) {
			if con.Type == ConstraintNotNull {
				t.Fatalf("expected child NOT NULL constraint to be dropped, still have %q", con.Name)
			}
		}
	}

	// Explicit NOT NULL names are reserved before unnamed columns pick implicit names, regardless of column order.
	if c, err := LoadSQLWithVersion(`CREATE TABLE t_ord (a int NOT NULL, b int CONSTRAINT t_ord_a_not_null NOT NULL);`, PG18); err != nil {
		t.Fatalf("PG18 CREATE TABLE with colliding explicit NOT NULL after implicit column: %v", err)
	} else {
		rel := c.GetRelation("public", "t_ord")
		got := map[int16]string{}
		for _, con := range c.ConstraintsOf(rel.OID) {
			if con.Type == ConstraintNotNull && len(con.Columns) == 1 {
				got[con.Columns[0]] = con.Name
			}
		}
		if got[2] != "t_ord_a_not_null" || got[1] != "t_ord_a_not_null1" {
			t.Fatalf("NOT NULL names = %v, want col1=t_ord_a_not_null1 col2=t_ord_a_not_null", got)
		}
	}

	// Dropping parent NOT NULL preserves child's own locally declared NOT NULL constraint.
	localChildSetup := `
CREATE TABLE p_loc (x int NOT NULL);
CREATE TABLE c_loc (x int CONSTRAINT child_nn NOT NULL) INHERITS (p_loc);
`
	if got, c := execCodeVersion(t, PG18, localChildSetup, "ALTER TABLE p_loc DROP CONSTRAINT p_loc_x_not_null"); got != "" {
		t.Fatalf("PG18 DROP CONSTRAINT p_loc_x_not_null: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_loc")
		ch := c.GetRelation("public", "c_loc")
		if p.Columns[0].NotNull {
			t.Fatalf("expected parent p_loc.x to become nullable")
		}
		if !ch.Columns[0].NotNull {
			t.Fatalf("expected child c_loc.x to stay NOT NULL due to local child_nn")
		}
		var found bool
		for _, con := range c.ConstraintsOf(ch.OID) {
			if con.Type == ConstraintNotNull && con.Name == "child_nn" && con.ConIsLocal && con.ConInhCount == 0 {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected child_nn (ConIsLocal=true, ConInhCount=0) to survive, got %+v", c.ConstraintsOf(ch.OID))
		}
	}

	// Multiple inheritance preserves child NOT NULL until all enforcing parents drop theirs.
	multiParentSetup := `
CREATE TABLE p1_mi (x int NOT NULL);
CREATE TABLE p2_mi (x int NOT NULL);
CREATE TABLE c_mi () INHERITS (p1_mi, p2_mi);
`
	if got, c := execCodeVersion(t, PG18, multiParentSetup, "ALTER TABLE p1_mi DROP CONSTRAINT p1_mi_x_not_null"); got != "" {
		t.Fatalf("PG18 DROP CONSTRAINT p1_mi_x_not_null: got %q, want ok", got)
	} else {
		ch := c.GetRelation("public", "c_mi")
		if !ch.Columns[0].NotNull {
			t.Fatalf("expected c_mi.x to remain NOT NULL while p2_mi still requires it")
		}
		if res, err := c.Exec("ALTER TABLE p2_mi DROP CONSTRAINT p2_mi_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP CONSTRAINT p2_mi_x_not_null: %v", res[0].Error)
		}
		if ch.Columns[0].NotNull {
			t.Fatalf("expected c_mi.x to become nullable after both parents dropped NOT NULL")
		}
	}

	// Recursive SET NOT NULL on a partitioned table records inherited NOT NULL on partitions
	// so a subsequent DROP CONSTRAINT on the parent also drops the partition's NOT NULL.
	partSetDropSetup := `
CREATE TABLE p_sd (x int, k int) PARTITION BY RANGE (k);
CREATE TABLE p_sd_1 PARTITION OF p_sd FOR VALUES FROM (0) TO (100);
`
	if got, c := execCodeVersion(t, PG18, partSetDropSetup, "ALTER TABLE p_sd ALTER COLUMN x SET NOT NULL"); got != "" {
		t.Fatalf("PG18 SET NOT NULL on partitioned table: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_sd")
		p1 := c.GetRelation("public", "p_sd_1")
		if !p.Columns[0].NotNull || !p1.Columns[0].NotNull {
			t.Fatalf("expected both parent and partition x to be NOT NULL")
		}
		if res, err := c.Exec("ALTER TABLE p_sd DROP CONSTRAINT p_sd_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP CONSTRAINT p_sd_x_not_null: %v", res[0].Error)
		}
		if p.Columns[0].NotNull || p1.Columns[0].NotNull {
			t.Fatalf("expected both parent and partition x to become nullable (p=%v, p1=%v)", p.Columns[0].NotNull, p1.Columns[0].NotNull)
		}
	}

	// NO INHERIT updates NOT NULL constraint inheritance metadata so the detached child can drop it.
	noInhSetup := `
CREATE TABLE p_ni (x int NOT NULL);
CREATE TABLE c_ni () INHERITS (p_ni);
`
	if got, c := execCodeVersion(t, PG18, noInhSetup, "ALTER TABLE c_ni NO INHERIT p_ni"); got != "" {
		t.Fatalf("PG18 NO INHERIT: got %q, want ok", got)
	} else {
		if res, err := c.Exec("ALTER TABLE c_ni DROP CONSTRAINT p_ni_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP CONSTRAINT p_ni_x_not_null after NO INHERIT: %v", res[0].Error)
		}
		ch := c.GetRelation("public", "c_ni")
		if ch.Columns[0].NotNull {
			t.Fatalf("expected c_ni.x to become nullable after NO INHERIT + DROP CONSTRAINT")
		}
	}

	// INHERIT increments NOT NULL inheritance metadata; ADD COLUMN NOT NULL on parent marks child inherited.
	addInhSetup := `
CREATE TABLE p_ai (x int NOT NULL);
CREATE TABLE c_ai (x int NOT NULL);
`
	if got, c := execCodeVersion(t, PG18, addInhSetup, "ALTER TABLE c_ai INHERIT p_ai"); got != "" {
		t.Fatalf("PG18 INHERIT: got %q, want ok", got)
	} else {
		if res, err := c.Exec("ALTER TABLE c_ai DROP CONSTRAINT c_ai_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) == 0 || res[0].Error == nil {
			t.Fatalf("expected DROP CONSTRAINT c_ai_x_not_null to fail while inheriting from p_ai")
		}
		if res, err := c.Exec("ALTER TABLE p_ai ADD COLUMN z int NOT NULL", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 ADD COLUMN z NOT NULL: %v", res[0].Error)
		}
		if res, err := c.Exec("ALTER TABLE p_ai DROP CONSTRAINT p_ai_z_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP CONSTRAINT p_ai_z_not_null: %v", res[0].Error)
		}
		ch := c.GetRelation("public", "c_ai")
		if ch.Columns[1].NotNull {
			t.Fatalf("expected inherited column c_ai.z to become nullable after parent drops p_ai_z_not_null")
		}
	}

	// Non-recursive DROP CONSTRAINT (ALTER TABLE ONLY) retains child's NOT NULL
	// while marking it local so the child can later drop it independently.
	onlyDropSetup := `
CREATE TABLE p_only (x int NOT NULL);
CREATE TABLE c_only () INHERITS (p_only);
`
	if got, c := execCodeVersion(t, PG18, onlyDropSetup, "ALTER TABLE ONLY p_only DROP CONSTRAINT p_only_x_not_null"); got != "" {
		t.Fatalf("PG18 ALTER TABLE ONLY DROP CONSTRAINT: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_only")
		ch := c.GetRelation("public", "c_only")
		if p.Columns[0].NotNull {
			t.Fatalf("expected parent p_only.x to become nullable")
		}
		if !ch.Columns[0].NotNull {
			t.Fatalf("expected child c_only.x to retain NOT NULL after ONLY drop on parent")
		}
		if res, err := c.Exec("ALTER TABLE c_only DROP CONSTRAINT p_only_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 child DROP CONSTRAINT after parent ONLY drop: %v", res[0].Error)
		}
		if ch.Columns[0].NotNull {
			t.Fatalf("expected c_only.x to become nullable after child drops p_only_x_not_null")
		}
	}

	// Later explicit CHECK constraint name is reserved before implicit PK NOT NULL picks a name.
	if c, err := LoadSQLWithVersion(`CREATE TABLE t_pk_chk (a int PRIMARY KEY, b int CONSTRAINT t_pk_chk_a_not_null CHECK (b > 0));`, PG18); err != nil {
		t.Fatalf("PG18 CREATE TABLE with PK and colliding later CHECK: %v", err)
	} else {
		rel := c.GetRelation("public", "t_pk_chk")
		var nnName string
		for _, con := range c.ConstraintsOf(rel.OID) {
			if con.Type == ConstraintNotNull {
				nnName = con.Name
			}
		}
		if nnName != "t_pk_chk_a_not_null1" {
			t.Fatalf("implicit PK NOT NULL name = %q, want %q", nnName, "t_pk_chk_a_not_null1")
		}
		if res, err := c.Exec("ALTER TABLE t_pk_chk DROP CONSTRAINT t_pk_chk_a_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP CONSTRAINT t_pk_chk_a_not_null (CHECK): %v", res[0].Error)
		}
	}

	// ALTER TABLE ONLY on a partitioned table drops parent NOT NULL while keeping partition NOT NULL local.
	partOnlySetup := `
CREATE TABLE p_po (x int NOT NULL, k int) PARTITION BY RANGE (k);
CREATE TABLE p_po_1 PARTITION OF p_po FOR VALUES FROM (0) TO (100);
`
	if got, c := execCodeVersion(t, PG18, partOnlySetup, "ALTER TABLE ONLY p_po ALTER COLUMN x DROP NOT NULL"); got != "" {
		t.Fatalf("PG18 ALTER TABLE ONLY partitioned DROP NOT NULL: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_po")
		p1 := c.GetRelation("public", "p_po_1")
		if p.Columns[0].NotNull {
			t.Fatalf("expected parent p_po.x to become nullable")
		}
		if !p1.Columns[0].NotNull {
			t.Fatalf("expected partition p_po_1.x to retain NOT NULL after ONLY drop on parent")
		}
		if res, err := c.Exec("ALTER TABLE p_po_1 DROP CONSTRAINT p_po_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 partition DROP CONSTRAINT after parent ONLY drop: %v", res[0].Error)
		}
		if p1.Columns[0].NotNull {
			t.Fatalf("expected p_po_1.x to become nullable after partition drops p_po_x_not_null")
		}
	}

	// Typed table (CREATE TABLE OF type) preserves explicit CONSTRAINT <name> NOT NULL.
	typedSetup := `
CREATE TYPE ty_nn AS (x int);
CREATE TABLE t_ty OF ty_nn (x WITH OPTIONS CONSTRAINT nn NOT NULL);
`
	if got, c := execCodeVersion(t, PG18, typedSetup, "COMMENT ON CONSTRAINT nn ON t_ty IS 'typed nn'"); got != "" {
		t.Fatalf("PG18 COMMENT ON CONSTRAINT nn on typed table: got %q, want ok", got)
	} else if res, err := c.Exec("ALTER TABLE t_ty DROP CONSTRAINT nn", nil); err != nil {
		t.Fatal(err)
	} else if len(res) > 0 && res[0].Error != nil {
		t.Fatalf("PG18 DROP CONSTRAINT nn on typed table: %v", res[0].Error)
	}

	// ADD COLUMN inline CHECK constraint name is reserved before implicit NOT NULL picks a name.
	if got, c := execCodeVersion(t, PG18, `CREATE TABLE t_add_chk (a int);`,
		"ALTER TABLE t_add_chk ADD COLUMN x int NOT NULL CONSTRAINT t_add_chk_x_not_null CHECK (x > 0)"); got != "" {
		t.Fatalf("PG18 ADD COLUMN NOT NULL + colliding CHECK: got %q, want ok", got)
	} else {
		rel := c.GetRelation("public", "t_add_chk")
		var nnName string
		for _, con := range c.ConstraintsOf(rel.OID) {
			if con.Type == ConstraintNotNull {
				nnName = con.Name
			}
		}
		if nnName != "t_add_chk_x_not_null1" {
			t.Fatalf("implicit ADD COLUMN NOT NULL name = %q, want %q", nnName, "t_add_chk_x_not_null1")
		}
		if res, err := c.Exec("ALTER TABLE t_add_chk DROP CONSTRAINT t_add_chk_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP CONSTRAINT t_add_chk_x_not_null (CHECK): %v", res[0].Error)
		}
		if !rel.Columns[1].NotNull {
			t.Fatalf("expected x to remain NOT NULL after dropping CHECK t_add_chk_x_not_null")
		}
	}

	// No-op DROP NOT NULL on a nullable parent does not disturb children's NOT NULL inheritance count.
	noopDropSetup := `
CREATE TABLE p1_np (x int);
CREATE TABLE p2_np (x int NOT NULL);
CREATE TABLE c_np () INHERITS (p1_np, p2_np);
`
	if got, c := execCodeVersion(t, PG18, noopDropSetup, "ALTER TABLE ONLY p1_np ALTER x DROP NOT NULL"); got != "" {
		t.Fatalf("PG18 no-op ONLY DROP NOT NULL: got %q, want ok", got)
	} else {
		if res, err := c.Exec("ALTER TABLE c_np DROP CONSTRAINT p2_np_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) == 0 || res[0].Error == nil {
			t.Fatalf("expected child DROP CONSTRAINT to fail while p2_np still enforces NOT NULL")
		} else {
			assertErrorCode(t, res[0].Error, CodeInvalidTableDefinition)
		}
		if res, err := c.Exec("ALTER TABLE p2_np DROP CONSTRAINT p2_np_x_not_null", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 DROP CONSTRAINT p2_np_x_not_null: %v", res[0].Error)
		}
		ch := c.GetRelation("public", "c_np")
		if ch.Columns[0].NotNull {
			t.Fatalf("expected c_np.x to become nullable after p2_np drops p2_np_x_not_null")
		}
	}

	// Explicit table constraint colliding with an explicit NOT NULL constraint name is rejected.
	if got, _ := execCodeVersion(t, PG18, "", "CREATE TABLE t_dup_chk (x int CONSTRAINT nn NOT NULL, CONSTRAINT nn CHECK (x > 0))"); got != CodeDuplicateObject {
		t.Fatalf("PG18 CREATE TABLE duplicate explicit NOT NULL and CHECK name: got %q, want %q", got, CodeDuplicateObject)
	}

	// ALTER TABLE ONLY DROP COLUMN on parent updates child's NOT NULL inheritance metadata so child can drop NOT NULL.
	onlyDropColSetup := `
CREATE TABLE p_dc (x int NOT NULL, y int);
CREATE TABLE c_dc () INHERITS (p_dc);
`
	if got, c := execCodeVersion(t, PG18, onlyDropColSetup, "ALTER TABLE ONLY p_dc DROP COLUMN x"); got != "" {
		t.Fatalf("PG18 ALTER TABLE ONLY DROP COLUMN: got %q, want ok", got)
	} else {
		if res, err := c.Exec("ALTER TABLE c_dc ALTER COLUMN x DROP NOT NULL", nil); err != nil {
			t.Fatal(err)
		} else if len(res) > 0 && res[0].Error != nil {
			t.Fatalf("PG18 child DROP NOT NULL after parent ONLY DROP COLUMN: %v", res[0].Error)
		}
		ch := c.GetRelation("public", "c_dc")
		if ch.Columns[0].NotNull {
			t.Fatalf("expected c_dc.x to become nullable after parent ONLY DROP COLUMN + child DROP NOT NULL")
		}
	}

	// Recursive DROP COLUMN on one parent updates surviving child's NOT NULL constraint.
	recDropColSetup := `
CREATE TABLE p1_dc (x int NOT NULL, y int);
CREATE TABLE p2_dc (x int);
CREATE TABLE c_rdc () INHERITS (p1_dc, p2_dc);
`
	if got, c := execCodeVersion(t, PG18, recDropColSetup, "ALTER TABLE p1_dc DROP COLUMN x"); got != "" {
		t.Fatalf("PG18 recursive DROP COLUMN on p1_dc: got %q, want ok", got)
	} else {
		ch := c.GetRelation("public", "c_rdc")
		if ch.Columns[0].NotNull {
			t.Fatalf("expected surviving c_rdc.x to become nullable after sole NOT NULL parent p1_dc drops x")
		}
	}

	// ALTER TABLE ONLY SET NOT NULL on a nullable partitioned table with partitions fails with 42P16
	// without mutating either relation.
	partOnlySetSetup := `
CREATE TABLE p_os (x int, k int) PARTITION BY RANGE (k);
CREATE TABLE p_os_1 PARTITION OF p_os FOR VALUES FROM (0) TO (100);
`
	if got, c := execCodeVersion(t, PG18, partOnlySetSetup, "ALTER TABLE ONLY p_os ALTER COLUMN x SET NOT NULL"); got != CodeInvalidTableDefinition {
		t.Fatalf("PG18 ALTER TABLE ONLY partitioned SET NOT NULL: got %q, want %q", got, CodeInvalidTableDefinition)
	} else {
		p := c.GetRelation("public", "p_os")
		p1 := c.GetRelation("public", "p_os_1")
		if p.Columns[0].NotNull || p1.Columns[0].NotNull {
			t.Fatalf("expected failed ONLY SET NOT NULL to leave both parent and partition nullable")
		}
	}

	// ALTER TABLE ONLY SET NOT NULL on an ordinary inheritance parent with children marks the constraint
	// NO INHERIT; a subsequent recursive SET NOT NULL is rejected with 42P16.
	inhOnlySetSetup := `
CREATE TABLE p_ois (x int);
CREATE TABLE c_ois () INHERITS (p_ois);
`
	if got, c := execCodeVersion(t, PG18, inhOnlySetSetup, "ALTER TABLE ONLY p_ois ALTER COLUMN x SET NOT NULL"); got != "" {
		t.Fatalf("PG18 ALTER TABLE ONLY inheritance SET NOT NULL: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_ois")
		ch := c.GetRelation("public", "c_ois")
		if !p.Columns[0].NotNull || ch.Columns[0].NotNull {
			t.Fatalf("expected p_ois.x NOT NULL and c_ois.x nullable (p=%v, ch=%v)", p.Columns[0].NotNull, ch.Columns[0].NotNull)
		}
		var noInh bool
		for _, con := range c.ConstraintsOf(p.OID) {
			if con.Type == ConstraintNotNull && con.ConNoInherit {
				noInh = true
			}
		}
		if !noInh {
			t.Fatalf("expected ConNoInherit=true on p_ois_x_not_null")
		}
		if res, err := c.Exec("ALTER TABLE p_ois ALTER COLUMN x SET NOT NULL", nil); err != nil {
			t.Fatal(err)
		} else if len(res) == 0 || res[0].Error == nil {
			t.Fatalf("expected recursive SET NOT NULL on NO INHERIT constraint to fail")
		} else {
			assertErrorCode(t, res[0].Error, CodeInvalidTableDefinition)
		}
	}

	// Non-recursive DROP COLUMN on a parent whose NOT NULL is NO INHERIT does not decrement
	// a child's NOT NULL inheritance count inherited from another parent.
	onlyDropNoInhSetup := `
CREATE TABLE p_odni (x int);
CREATE TABLE q_odni (x int NOT NULL);
CREATE TABLE c_odni () INHERITS (p_odni, q_odni);
ALTER TABLE ONLY p_odni ALTER COLUMN x SET NOT NULL;
`
	if got, c := execCodeVersion(t, PG18, onlyDropNoInhSetup, "ALTER TABLE ONLY p_odni DROP COLUMN x"); got != "" {
		t.Fatalf("PG18 ALTER TABLE ONLY DROP COLUMN with NO INHERIT NOT NULL: got %q, want ok", got)
	} else {
		if res, err := c.Exec("ALTER TABLE c_odni ALTER COLUMN x DROP NOT NULL", nil); err != nil {
			t.Fatal(err)
		} else if len(res) == 0 || res[0].Error == nil {
			t.Fatalf("expected c_odni DROP NOT NULL to fail while q_odni still requires NOT NULL")
		} else {
			assertErrorCode(t, res[0].Error, CodeInvalidTableDefinition)
		}
	}

	// Recursive SET NOT NULL on an ancestor fails with 42P16 and rolls back when a descendant
	// holds a NO INHERIT NOT NULL constraint on the same column.
	descNoInhSetup := `
CREATE TABLE p_dni (x int);
CREATE TABLE c_dni () INHERITS (p_dni);
CREATE TABLE g_dni () INHERITS (c_dni);
ALTER TABLE ONLY c_dni ALTER COLUMN x SET NOT NULL;
`
	if got, c := execCodeVersion(t, PG18, descNoInhSetup, "ALTER TABLE p_dni ALTER COLUMN x SET NOT NULL"); got != CodeInvalidTableDefinition {
		t.Fatalf("PG18 recursive SET NOT NULL with descendant NO INHERIT conflict: got %q, want %q", got, CodeInvalidTableDefinition)
	} else {
		p := c.GetRelation("public", "p_dni")
		ch := c.GetRelation("public", "c_dni")
		g := c.GetRelation("public", "g_dni")
		if p.Columns[0].NotNull || !ch.Columns[0].NotNull || g.Columns[0].NotNull {
			t.Fatalf("expected rollback on descendant NO INHERIT conflict (p=%v, c=%v, g=%v)", p.Columns[0].NotNull, ch.Columns[0].NotNull, g.Columns[0].NotNull)
		}
	}

	// Unnamed foreign key avoids colliding with an explicit NOT NULL constraint named <rel>_<col>_fkey.
	fkCollSetup := `
CREATE TABLE p_fkc (id int PRIMARY KEY);
CREATE TABLE t_fkc (x int CONSTRAINT t_fkc_x_fkey NOT NULL REFERENCES p_fkc);
`
	if got, c := execCodeVersion(t, PG18, fkCollSetup, "ALTER TABLE t_fkc DROP CONSTRAINT t_fkc_x_fkey"); got != "" {
		t.Fatalf("PG18 drop explicit NOT NULL named t_fkc_x_fkey: got %q, want ok", got)
	} else {
		rel := c.GetRelation("public", "t_fkc")
		if rel.Columns[0].NotNull {
			t.Fatalf("expected t_fkc.x to become nullable after dropping t_fkc_x_fkey")
		}
		var fkName string
		for _, con := range c.ConstraintsOf(rel.OID) {
			if con.Type == ConstraintFK {
				fkName = con.Name
			}
		}
		if fkName != "t_fkc_x_fkey1" {
			t.Fatalf("expected unnamed FK to receive name t_fkc_x_fkey1, got %q", fkName)
		}
	}

	// LIKE INCLUDING ALL generates a fresh PK constraint name without colliding
	// with the destination's PG 18 NOT NULL constraint or copied CHECK constraints.
	likeCollSetup := `
CREATE TABLE src_lc (x int CONSTRAINT t_lc_x_not_null PRIMARY KEY);
CREATE TABLE t_lc (LIKE src_lc INCLUDING ALL);
`
	if got, c := execCodeVersion(t, PG18, likeCollSetup, "ALTER TABLE t_lc DROP CONSTRAINT t_lc_pkey"); got != "" {
		t.Fatalf("PG18 LIKE INCLUDING ALL drop cloned PK t_lc_pkey: got %q, want ok", got)
	} else {
		rel := c.GetRelation("public", "t_lc")
		seen := make(map[string]int)
		for _, con := range c.ConstraintsOf(rel.OID) {
			seen[con.Name]++
		}
		for name, count := range seen {
			if count > 1 {
				t.Fatalf("duplicate constraint name %q (%d times) on t_lc", name, count)
			}
		}
	}

	// ADD PRIMARY KEY rejects a column backed by a NO INHERIT NOT NULL constraint with 55000
	// without mutating the catalog.
	pkNoInhSetup := `
CREATE TABLE p_pkni (x int);
CREATE TABLE c_pkni () INHERITS (p_pkni);
ALTER TABLE ONLY p_pkni ALTER COLUMN x SET NOT NULL;
`
	if got, c := execCodeVersion(t, PG18, pkNoInhSetup, "ALTER TABLE ONLY p_pkni ADD PRIMARY KEY (x)"); got != CodeObjectNotInPrerequisiteState {
		t.Fatalf("PG18 ADD PRIMARY KEY on NO INHERIT NOT NULL: got %q, want %q", got, CodeObjectNotInPrerequisiteState)
	} else {
		p := c.GetRelation("public", "p_pkni")
		if c.findPKConstraint(p.OID) != nil {
			t.Fatalf("expected failed ADD PRIMARY KEY not to register a PK constraint")
		}
		if len(c.IndexesOf(p.OID)) != 0 {
			t.Fatalf("expected failed ADD PRIMARY KEY not to register a PK index")
		}
	}

	// Migrating a PG 18 NOT NULL column with a constraint comment to nullable
	// suppresses COMMENT ON CONSTRAINT ... IS NULL (constraint removed by DROP NOT NULL),
	// while removing only the comment on a kept NOT NULL constraint still emits IS NULL.
	fromDropNN, err := LoadSQLWithVersion(`
CREATE TABLE t_nn_cmt (x int NOT NULL);
COMMENT ON CONSTRAINT t_nn_cmt_x_not_null ON t_nn_cmt IS 'required';
`, PG18)
	if err != nil {
		t.Fatalf("LoadSQLWithVersion fromDropNN: %v", err)
	}
	toNullable, err := LoadSQLWithVersion(`CREATE TABLE t_nn_cmt (x int);`, PG18)
	if err != nil {
		t.Fatalf("LoadSQLWithVersion toNullable: %v", err)
	}
	planDropNN := GenerateMigration(fromDropNN, toNullable, Diff(fromDropNN, toNullable))
	for _, op := range planDropNN.Ops {
		if strings.Contains(op.SQL, "COMMENT ON CONSTRAINT") {
			t.Fatalf("unexpected COMMENT ON CONSTRAINT for dropped NOT NULL constraint: %s", op.SQL)
		}
	}

	toKeptNN, err := LoadSQLWithVersion(`CREATE TABLE t_nn_cmt (x int NOT NULL);`, PG18)
	if err != nil {
		t.Fatalf("LoadSQLWithVersion toKeptNN: %v", err)
	}
	planKeptNN := GenerateMigration(fromDropNN, toKeptNN, Diff(fromDropNN, toKeptNN))
	var sawCommentDrop bool
	for _, op := range planKeptNN.Ops {
		if strings.Contains(op.SQL, "COMMENT ON CONSTRAINT") && strings.Contains(op.SQL, "IS NULL") {
			sawCommentDrop = true
		}
	}
	if !sawCommentDrop {
		t.Fatalf("expected COMMENT ON CONSTRAINT ... IS NULL when NOT NULL constraint is kept")
	}

	// Recursive SET NOT NULL preserves the parent's chosen NOT NULL constraint name on children,
	// avoiding collisions across the inheritance hierarchy.
	recSetNameSetup := `
CREATE TABLE p_rsn (x int);
CREATE TABLE c_rsn () INHERITS (p_rsn);
ALTER TABLE p_rsn ALTER COLUMN x SET NOT NULL;
`
	if got, _ := execCodeVersion(t, PG18, recSetNameSetup, "COMMENT ON CONSTRAINT p_rsn_x_not_null ON c_rsn IS 'required'"); got != "" {
		t.Fatalf("PG18 COMMENT ON inherited NOT NULL constraint p_rsn_x_not_null on c_rsn: got %q, want ok", got)
	}

	recSetCollSetup := `
CREATE TABLE p_rsc (x int);
CREATE TABLE c_rsc (CONSTRAINT p_rsc_x_not_null CHECK (x > 0)) INHERITS (p_rsc);
`
	if got, c := execCodeVersion(t, PG18, recSetCollSetup, "ALTER TABLE p_rsc ALTER COLUMN x SET NOT NULL"); got != "" {
		t.Fatalf("PG18 recursive SET NOT NULL with child name collision: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_rsc")
		ch := c.GetRelation("public", "c_rsc")
		if gotP := c.notNullConstraintName(p, 1); gotP != "p_rsc_x_not_null1" {
			t.Fatalf("p_rsc NOT NULL name = %q, want p_rsc_x_not_null1", gotP)
		}
		if gotC := c.notNullConstraintName(ch, 1); gotC != "p_rsc_x_not_null1" {
			t.Fatalf("c_rsc NOT NULL name = %q, want p_rsc_x_not_null1", gotC)
		}
	}

	// Generated NOT NULL and FK names avoid collisions across all constraints in the schema
	// (including other tables' constraints and domain constraints), matching PG ChooseConstraintName.
	schemaCollSetup := `
CREATE TABLE other_sc (x int CONSTRAINT t_sc_x_not_null CHECK (x > 0));
CREATE DOMAIN d_sc AS int CONSTRAINT t_sc2_x_not_null CHECK (VALUE > 0);
CREATE TABLE t_sc (x int NOT NULL);
CREATE TABLE t_sc2 (x int NOT NULL);
CREATE TABLE p_sc_rec (x int);
CREATE TABLE c_sc_rec () INHERITS (p_sc_rec);
CREATE TABLE other_sc_rec (x int CONSTRAINT p_sc_rec_x_not_null CHECK (x > 0));
ALTER TABLE p_sc_rec ALTER COLUMN x SET NOT NULL;
`
	if got, _ := execCodeVersion(t, PG18, schemaCollSetup, `
COMMENT ON CONSTRAINT t_sc_x_not_null1 ON t_sc IS 'ok1';
COMMENT ON CONSTRAINT t_sc2_x_not_null1 ON t_sc2 IS 'ok2';
COMMENT ON CONSTRAINT p_sc_rec_x_not_null1 ON c_sc_rec IS 'ok3';
`); got != "" {
		t.Fatalf("PG18 schema-wide ChooseConstraintName collision avoidance: got %q, want ok", got)
	}

	// ATTACH PARTITION clears ConIsLocal on the merged NOT NULL constraint so a
	// subsequent recursive DROP NOT NULL on the parent also drops it on the partition.
	attachDropSetup := `
CREATE TABLE p_ad (x int NOT NULL, k int) PARTITION BY RANGE (k);
CREATE TABLE c_ad (x int NOT NULL, k int);
ALTER TABLE p_ad ATTACH PARTITION c_ad FOR VALUES FROM (0) TO (100);
`
	if got, c := execCodeVersion(t, PG18, attachDropSetup, "ALTER TABLE p_ad ALTER COLUMN x DROP NOT NULL"); got != "" {
		t.Fatalf("PG18 DROP NOT NULL after ATTACH PARTITION: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_ad")
		part := c.GetRelation("public", "c_ad")
		if p.Columns[0].NotNull || part.Columns[0].NotNull {
			t.Fatalf("expected both p_ad.x and c_ad.x to become nullable (p=%v, part=%v)", p.Columns[0].NotNull, part.Columns[0].NotNull)
		}
	}

	// ALTER TABLE c INHERIT p rejects merging an inheritable parent NOT NULL into
	// a child's NO INHERIT NOT NULL constraint with 42P16 without mutating c.
	inhConflictSetup := `
CREATE TABLE p_ic (x int NOT NULL);
CREATE TABLE c_ic (x int);
CREATE TABLE g_ic () INHERITS (c_ic);
ALTER TABLE ONLY c_ic ALTER COLUMN x SET NOT NULL;
`
	if got, c := execCodeVersion(t, PG18, inhConflictSetup, "ALTER TABLE c_ic INHERIT p_ic"); got != CodeInvalidTableDefinition {
		t.Fatalf("PG18 INHERIT with NO INHERIT NOT NULL conflict: got %q, want %q", got, CodeInvalidTableDefinition)
	} else {
		ch := c.GetRelation("public", "c_ic")
		if len(ch.InhParents) != 0 {
			t.Fatalf("expected failed INHERIT not to attach p_ic as parent, got %v", ch.InhParents)
		}
	}

	// Attaching a partition to a parent with a PRIMARY KEY counts NOT NULL inheritance only once,
	// so dropping the PK and then dropping NOT NULL on the parent makes the partition nullable.
	attachPKDropSetup := `
CREATE TABLE p_apk (x int PRIMARY KEY) PARTITION BY RANGE (x);
CREATE TABLE c_apk (x int NOT NULL);
ALTER TABLE p_apk ATTACH PARTITION c_apk FOR VALUES FROM (0) TO (100);
ALTER TABLE p_apk DROP CONSTRAINT p_apk_pkey;
`
	if got, c := execCodeVersion(t, PG18, attachPKDropSetup, "ALTER TABLE p_apk ALTER COLUMN x DROP NOT NULL"); got != "" {
		t.Fatalf("PG18 DROP NOT NULL after ATTACH PARTITION with PK: got %q, want ok", got)
	} else {
		p := c.GetRelation("public", "p_apk")
		part := c.GetRelation("public", "c_apk")
		if p.Columns[0].NotNull || part.Columns[0].NotNull {
			t.Fatalf("expected both p_apk.x and c_apk.x to become nullable (p=%v, part=%v)", p.Columns[0].NotNull, part.Columns[0].NotNull)
		}
	}

	// Migrating a PG 18 nullable column to a named NOT NULL constraint with a comment
	// emits SET NOT NULL + RENAME CONSTRAINT before COMMENT ON CONSTRAINT, and executes cleanly.
	fromNamedNN, err := LoadSQLWithVersion(`CREATE TABLE t_nnn (x int);`, PG18)
	if err != nil {
		t.Fatalf("LoadSQLWithVersion fromNamedNN: %v", err)
	}
	toNamedNN, err := LoadSQLWithVersion(`
CREATE TABLE t_nnn (x int CONSTRAINT nn NOT NULL);
COMMENT ON CONSTRAINT nn ON t_nnn IS 'required';
`, PG18)
	if err != nil {
		t.Fatalf("LoadSQLWithVersion toNamedNN: %v", err)
	}
	planNamedNN := GenerateMigration(fromNamedNN, toNamedNN, Diff(fromNamedNN, toNamedNN))
	if len(planNamedNN.Ops) == 0 {
		t.Fatalf("expected non-empty migration plan for named NOT NULL + comment")
	}
	for _, op := range planNamedNN.Ops {
		res, err := fromNamedNN.Exec(op.SQL, nil)
		if err != nil {
			t.Fatalf("exec migration op %q: %v", op.SQL, err)
		}
		for _, r := range res {
			if r.Error != nil {
				t.Fatalf("migration op %q failed: %v", op.SQL, r.Error)
			}
		}
	}
	if rem := GenerateMigration(fromNamedNN, toNamedNN, Diff(fromNamedNN, toNamedNN)); len(rem.Ops) != 0 {
		t.Fatalf("expected empty diff after applying named NOT NULL migration, still have: %+v", rem.Ops)
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

	// Attaching a partition that already has a CHECK named <part>_pkey picks <part>_pkey1
	// for the cloned PK index and constraint instead of creating a duplicate constraint name.
	pkCollisionSetup := `
CREATE TABLE p_pkc (id int PRIMARY KEY) PARTITION BY RANGE (id);
CREATE TABLE c_pkc (id int NOT NULL, CONSTRAINT c_pkc_pkey CHECK (id > 0));
`
	if got, c := execCodeVersion(t, PG18, pkCollisionSetup, "ALTER TABLE p_pkc ATTACH PARTITION c_pkc FOR VALUES FROM (0) TO (100)"); got != "" {
		t.Fatalf("ATTACH PARTITION with c_pkc_pkey collision: got %q, want ok", got)
	} else {
		part := c.GetRelation("public", "c_pkc")
		pk := c.findPKConstraint(part.OID)
		if pk == nil || pk.Name != "c_pkc_pkey1" {
			t.Fatalf("attached partition PK constraint = %+v, want Name=%q", pk, "c_pkc_pkey1")
		}
		if res, err := c.Exec("COMMENT ON CONSTRAINT c_pkc_pkey1 ON c_pkc IS 'pk'", nil); err != nil {
			t.Fatal(err)
		} else if len(res) == 0 || res[0].Error != nil {
			t.Fatalf("COMMENT ON CONSTRAINT c_pkc_pkey1 failed: %v", res[0].Error)
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
