package review

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
)

func walkSchema() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{
		Name:       "db",
		SearchPath: `"$user", public`,
		Schemas: []*metadata.SchemaMetadata{{
			Name:   "public",
			Tables: []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "a"}}}},
		}},
	}
}

func walkFindings(t *testing.T, sql string, targets ...review.Target) []review.Finding {
	t.Helper()
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{review.WalkThrough}}, targets)
	if err != nil {
		t.Fatal(err)
	}
	return result.Findings
}

func TestWalkThroughReportsTheFirstFailure(t *testing.T) {
	got := walkFindings(t, "CREATE TABLE x (a int);\nALTER TABLE public.nosuch ADD COLUMN b int;\nCREATE TABLE t (a int);", review.Target{Schema: walkSchema(), SessionUser: "u"})
	if len(got) != 1 || got[0].Statement != 1 || got[0].Message != "relation public.nosuch does not exist" {
		t.Fatalf("findings = %+v", got)
	}
	// The finding anchors at the name as written.
	if want := (review.Range{Start: 36, End: 49}); got[0].Range != want {
		t.Errorf("range = %+v, want %+v", got[0].Range, want)
	}
	got = walkFindings(t, "CREATE TABLE public.\"t\" (a int);", review.Target{Schema: walkSchema(), SessionUser: "u"})
	if want := (review.Range{Start: 13, End: 23}); len(got) != 1 || got[0].Range != want {
		t.Errorf("findings = %+v, want one at %+v", got, want)
	}
}

// PostgreSQL truncates a catalog name like any other, so a longer
// spelling of a 63-byte database name is that database (verified on
// PostgreSQL 17).
func TestWalkThroughLongDatabaseName(t *testing.T) {
	schema := walkSchema()
	schema.Name = strings.Repeat("d", 63)
	long := `"` + schema.Name + `xyz"`
	got := walkFindings(t, "CREATE TABLE "+long+".public.x (a int);\nCREATE TABLE "+long+".public.x (a int);", review.Target{Schema: schema, SessionUser: "u"})
	if len(got) != 1 || got[0].Statement != 1 || got[0].Message != "relation "+schema.Name+".public.x already exists" {
		t.Errorf("findings = %+v, want the second CREATE's", got)
	}
	got = walkFindings(t, "CREATE TABLE other.public.x (a int);", review.Target{Schema: schema, SessionUser: "u"})
	if len(got) != 1 || got[0].Message != "cross-database references are not implemented: other.public.x" {
		t.Errorf("findings = %+v, want the cross-database reference", got)
	}
}

// PostgreSQL refuses CREATE INDEX CONCURRENTLY on a partitioned table
// (0A000, verified on PostgreSQL 17 and 18 outside a transaction, as
// Bytebase runs it; the oracle runs inside one, where 25001 comes first).
func TestWalkThroughConcurrentIndexOnPartitionedTable(t *testing.T) {
	schema := walkSchema()
	schema.Schemas[0].Tables = append(schema.Schemas[0].Tables, &metadata.TableMetadata{
		Name:       "parted",
		Columns:    []*metadata.ColumnMetadata{{Name: "k"}},
		Partitions: []*metadata.TablePartitionMetadata{{Name: "parted_1"}},
	})
	got := walkFindings(t, "CREATE INDEX CONCURRENTLY parted_k ON parted (k);", review.Target{Schema: schema, SessionUser: "u"})
	if len(got) != 1 || got[0].Message != "cannot create index on partitioned table parted concurrently" {
		t.Errorf("findings = %+v", got)
	}
	if got := walkFindings(t, "CREATE INDEX CONCURRENTLY t_a ON t (a);", review.Target{Schema: schema, SessionUser: "u"}); len(got) != 0 {
		t.Errorf("findings on a plain table = %+v", got)
	}
	// The sync shows a partitioned table without partitions as an ordinary
	// one; a partition the change creates proves it partitioned. Its
	// index then takes the partition's matching index, which goes with it
	// (verified on PostgreSQL 17 and 18).
	empty := walkSchema()
	empty.Schemas[0].Tables = append(empty.Schemas[0].Tables, &metadata.TableMetadata{Name: "ep", Columns: []*metadata.ColumnMetadata{{Name: "k"}}})
	alone := review.Target{Schema: empty, SessionUser: "u"}
	if got := walkFindings(t, "CREATE TABLE ep1 PARTITION OF ep FOR VALUES FROM (0) TO (10);\nCREATE INDEX ci ON ep1 (k);\nCREATE INDEX pi ON ep (k);\nDROP INDEX pi;\nCREATE INDEX ci ON ep1 (k);", alone); len(got) != 0 {
		t.Errorf("findings = %+v, want none", got)
	}
	if got := walkFindings(t, "CREATE TABLE ep1 PARTITION OF ep FOR VALUES FROM (0) TO (10);\nCREATE INDEX CONCURRENTLY pk ON ep (k);", alone); len(got) != 1 || got[0].Statement != 1 {
		t.Errorf("findings = %+v, want the CREATE INDEX", got)
	}

	// DROP INDEX CONCURRENTLY takes one index and no CASCADE (0A000).
	for sql, want := range map[string]string{
		"DROP INDEX CONCURRENTLY IF EXISTS x, y;":      "DROP INDEX CONCURRENTLY does not support dropping multiple objects",
		"DROP INDEX CONCURRENTLY IF EXISTS x CASCADE;": "DROP INDEX CONCURRENTLY does not support CASCADE",
	} {
		if got := walkFindings(t, sql, review.Target{Schema: schema, SessionUser: "u"}); len(got) != 1 || got[0].Message != want {
			t.Errorf("%s: findings = %+v, want %q", sql, got, want)
		}
	}
	if got := walkFindings(t, "DROP INDEX CONCURRENTLY IF EXISTS x;", review.Target{Schema: schema, SessionUser: "u"}); len(got) != 0 {
		t.Errorf("findings on one index = %+v", got)
	}
}

func TestWalkThroughWithoutASchema(t *testing.T) {
	if got := walkFindings(t, "CREATE TABLE x (a int);\nCREATE TABLE x (a int);", review.Target{}); len(got) != 0 {
		t.Errorf("findings without a synced schema: %+v", got)
	}
}

// With another spec of the plan on the same database, the synced schema
// says nothing: only the change's own conflicts are reported.
func TestWalkThroughOtherSpecs(t *testing.T) {
	other := review.Target{Schema: walkSchema(), SessionUser: "u", OtherSpecs: true}
	if got := walkFindings(t, "CREATE TABLE t (a int);\nALTER TABLE t ADD COLUMN a int;", other); len(got) != 0 {
		t.Errorf("findings from the synced schema: %+v", got)
	}
	// A table created in a schema the SQL names is there afterwards; one
	// the search path places may be anywhere, as the path's schemas are
	// unknown.
	got := walkFindings(t, "CREATE TABLE public.x (a int);\nCREATE TABLE public.x (b int);", other)
	if len(got) != 1 || got[0].Statement != 1 {
		t.Errorf("findings = %+v, want the second CREATE", got)
	}
	if got := walkFindings(t, "CREATE TABLE x (a int);\nCREATE TABLE x (b int);", other); len(got) != 0 {
		t.Errorf("findings = %+v, want none", got)
	}

	// The other spec may have added a trigger that DML fires.
	if got := walkFindings(t, "CREATE TABLE public.x (a int);\nINSERT INTO public.x VALUES (1);\nCREATE TABLE public.x (b int);", other); len(got) != 0 {
		t.Errorf("findings after DML = %+v, want none", got)
	}
	alone := review.Target{Schema: other.Schema, SessionUser: "u"}
	if got := walkFindings(t, "CREATE TABLE public.x (a int);\nINSERT INTO public.x VALUES (1);\nCREATE TABLE public.x (b int);", alone); len(got) != 1 || got[0].Statement != 2 {
		t.Errorf("findings = %+v, want the second CREATE", got)
	}

	// IF NOT EXISTS on a table that may exist may leave it an ordinary
	// table, which takes CONCURRENTLY (verified on PostgreSQL 17 and 18).
	if got := walkFindings(t, "CREATE TABLE IF NOT EXISTS public.y (a int) PARTITION BY RANGE (a);\nCREATE INDEX CONCURRENTLY y_a ON public.y (a);", other); len(got) != 0 {
		t.Errorf("findings = %+v, want none", got)
	}
	if got := walkFindings(t, "CREATE TABLE public.z (a int) PARTITION BY RANGE (a);\nCREATE INDEX CONCURRENTLY z_a ON public.z (a);", other); len(got) != 1 || got[0].Statement != 1 {
		t.Errorf("findings = %+v, want the CREATE INDEX", got)
	}

	// The other spec may have replaced a synced partitioned table with an
	// ordinary one.
	parted := walkSchema()
	parted.Schemas[0].Tables = append(parted.Schemas[0].Tables, &metadata.TableMetadata{
		Name:       "p",
		Columns:    []*metadata.ColumnMetadata{{Name: "k"}},
		Partitions: []*metadata.TablePartitionMetadata{{Name: "p_1"}},
	})
	if got := walkFindings(t, "CREATE TABLE IF NOT EXISTS public.p (k int);\nCREATE INDEX CONCURRENTLY p_k ON public.p (k);", review.Target{Schema: parted, SessionUser: "u", OtherSpecs: true}); len(got) != 0 {
		t.Errorf("findings = %+v, want none", got)
	}

	// The same schema with and without other specs is two inputs.
	got = walkFindings(t, "CREATE TABLE t (a int);", alone, other, alone)
	if len(got) != 1 || len(got[0].Targets) != 2 || got[0].Targets[0] != 0 || got[0].Targets[1] != 2 {
		t.Errorf("findings = %+v, want one on targets 0 and 2", got)
	}
}

// Bytebase runs CREATE and DROP INDEX CONCURRENTLY, VACUUM, DROP DATABASE,
// and a second run of SET ROLE after its transaction commits, and first
// of all in auto-commit mode; the walk follows that order.
func TestWalkThroughRunsAsBytebase(t *testing.T) {
	schema := walkSchema()
	schema.Schemas[0].Tables[0].Indexes = []*metadata.IndexMetadata{{Name: "t_a_idx", Definition: "CREATE INDEX t_a_idx ON public.t USING btree (a);"}}
	target := review.Target{Schema: schema, SessionUser: "u"}
	for _, tc := range []struct {
		sql  string
		stmt int // -1 for none
		msg  string
	}{
		// After the transaction: the table exists, the index does not yet
		// go.
		{"CREATE INDEX CONCURRENTLY i ON t2 (a);\nCREATE TABLE t2 (a int);", -1, ""},
		{"DROP INDEX CONCURRENTLY t_a_idx;\nCREATE INDEX t_a_idx ON t (a);", 1, "relation t_a_idx already exists"},
		{"CREATE INDEX CONCURRENTLY i ON t (a);\nCREATE INDEX i ON t (a);\nCREATE INDEX CONCURRENTLY i2 ON t (a);", 0, "relation i already exists"},
		// Matched anywhere in its text, as Bytebase matches it.
		{"COMMENT ON TABLE t2 IS 'no DROP DATABASE here';\nCREATE TABLE t2 (a int);", -1, ""},
		// First, in auto-commit mode.
		{"-- txn-mode = off\nCREATE TABLE t_a_idx (a int);\nDROP INDEX CONCURRENTLY t_a_idx;", -1, ""},
		{"-- txn-mode = off\nCREATE TABLE i (a int);\nCREATE INDEX CONCURRENTLY i ON t (a);", 0, "relation i already exists"},
		{"-- txn-isolation = SERIALIZABLE\n--   TXN-MODE=Off\nCREATE TABLE t_a_idx (a int);\nDROP INDEX CONCURRENTLY t_a_idx;", -1, ""},
		{"-- txn-mode = off\n-- txn-mode = on\nCREATE TABLE t_a_idx (a int);\nDROP INDEX CONCURRENTLY t_a_idx;", 0, "relation t_a_idx already exists"},
		// A directive below the first statement is no directive.
		{"CREATE TABLE t_a_idx (a int);\n-- txn-mode = off\nDROP INDEX CONCURRENTLY t_a_idx;", 0, "relation t_a_idx already exists"},
		// A CONCURRENTLY the patterns miss fails in the transaction: the
		// walk ends before it.
		{"CREATE /* x */ INDEX CONCURRENTLY t_a_idx ON t (a);\nDROP TABLE nosuch;", -1, ""},
		{"VACUUM t;\nDROP TABLE nosuch;", 1, "table nosuch does not exist"},
		{"/* x */ VACUUM t;\nDROP TABLE nosuch;", -1, ""},
		{"-- txn-mode = off\n/* x */ VACUUM t;\nDROP TABLE nosuch;", 1, "table nosuch does not exist"},
		{"BEGIN;\nREINDEX TABLE t;\nDROP TABLE nosuch;", -1, ""},
		// A second run of what is no SET ends the walk.
		{"COMMENT ON TABLE t IS 'SET ROLE';\nCREATE INDEX CONCURRENTLY i ON t (a);\nDROP TABLE nosuch;", 2, "table nosuch does not exist"},
		{"COMMENT ON TABLE t IS 'SET ROLE';\nDROP INDEX CONCURRENTLY t_a_idx;\nDROP INDEX t_a_idx;", -1, ""},
	} {
		got := walkFindings(t, tc.sql, target)
		switch {
		case tc.stmt < 0 && len(got) != 0:
			t.Errorf("%q: findings = %+v, want none", tc.sql, got)
		case tc.stmt >= 0 && (len(got) != 1 || got[0].Statement != tc.stmt || got[0].Message != tc.msg):
			t.Errorf("%q: findings = %+v, want %q on statement %d", tc.sql, got, tc.msg, tc.stmt)
		}
	}
}

// planRun lays a change out as Bytebase runs it.
func TestWalkThroughPlanRun(t *testing.T) {
	for _, tc := range []struct {
		sql        string
		autocommit bool
		steps      []step
	}{
		{"SET ROLE bob;\nCREATE INDEX CONCURRENTLY i ON t (a);\nGRANT SELECT ON t TO bob;\nVACUUM t;", false, []step{
			{stmt: 0}, {stmt: 2, wrapped: true}, {stmt: 0, again: true, committed: true}, {stmt: 1}, {stmt: 3},
		}},
		{"-- txn-mode = off\nRESET ROLE;\nCREATE EXTENSION hstore;\nDROP INDEX CONCURRENTLY i;", true, []step{
			{stmt: 0}, {stmt: 2}, {stmt: 0, again: true}, {stmt: 1, wrapped: true},
		}},
		{"CREATE TABLE t (a int);", false, []step{{stmt: 0}}},
	} {
		stmts, syntax := parse(tc.sql, statementRanges(tc.sql))
		if syntax != nil {
			t.Fatalf("%q: %s", tc.sql, syntax.Message)
		}
		got := planRun(tc.sql, stmts)
		if got == nil || got.autocommit != tc.autocommit || !slices.Equal(got.steps, tc.steps) {
			t.Errorf("%q: planRun = %+v, want auto-commit %v, %+v", tc.sql, got, tc.autocommit, tc.steps)
		}
	}
}

// Bytebase's sync drops the views and materialized views named like an
// extension's, and loses the partitions of a partitioned partition in
// another schema than its table's.
func TestWalkThroughSyncGaps(t *testing.T) {
	target := review.Target{Schema: walkSchema(), SessionUser: "u"}
	for _, sql := range []string{
		"GRANT SELECT ON hypopg_list_indexes TO PUBLIC;",
		"COMMENT ON MATERIALIZED VIEW g_agg_stat_x IS 'x';",
		"DROP VIEW public.g_columnar_x;",
		"DROP VIEW google_db_advisor_x;",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}

	parted := func(indexes ...string) review.Target {
		s := walkSchema()
		var idx []*metadata.IndexMetadata
		for i, def := range indexes {
			idx = append(idx, &metadata.IndexMetadata{Name: fmt.Sprintf("p1_%d", i), Definition: def})
		}
		s.Schemas[0].Tables = append(s.Schemas[0].Tables, &metadata.TableMetadata{
			Name:       "p",
			Partitions: []*metadata.TablePartitionMetadata{{Name: "p1", Indexes: idx}},
		})
		s.Schemas = append(s.Schemas, &metadata.SchemaMetadata{Name: "s2"})
		return review.Target{Schema: s, SessionUser: "u"}
	}
	const missing = "ALTER TABLE s2.nosuch ADD COLUMN b int;"
	for _, tc := range []struct {
		name   string
		target review.Target
		found  bool
	}{
		{"no index tells", parted(), false},
		{"a leaf", parted("CREATE INDEX p1_k_idx ON s2.p1 USING btree (k);"), true},
		{"partitioned in its table's schema", parted("CREATE INDEX p1_k_idx ON ONLY public.p1 USING btree (k);"), true},
		{"partitioned in another schema", parted("CREATE INDEX p1_k_idx ON ONLY s2.p1 USING btree (k);"), false},
		{"a quoted schema", parted(`CREATE INDEX p1_k_idx ON ONLY "S 2".p1 USING btree (k);`), false},
	} {
		if got := walkFindings(t, missing, tc.target); (len(got) == 1) != tc.found {
			t.Errorf("%s: findings = %+v, want found %v", tc.name, got, tc.found)
		}
	}
}

// In tenant mode Bytebase runs the change as the database owner by SET
// ROLE from an admin session the walk does not know: RESET ROLE returns
// to it, and Bytebase runs GRANT as it. The target names the owner as the
// change's user then, which is how the walk tells the mode may be on.
func TestWalkThroughTenantMode(t *testing.T) {
	schema := walkSchema()
	schema.Owner = "app"
	schema.Schemas = append(schema.Schemas,
		&metadata.SchemaMetadata{Name: "app", Tables: []*metadata.TableMetadata{{Name: "x"}}},
		&metadata.SchemaMetadata{Name: "admin", Tables: []*metadata.TableMetadata{{Name: "g"}}},
		&metadata.SchemaMetadata{Name: "bob", Tables: []*metadata.TableMetadata{{Name: "y"}}},
	)
	tenant := review.Target{Schema: schema, SessionUser: "app"}
	plain := review.Target{Schema: schema, SessionUser: "admin"}
	for _, tc := range []struct {
		sql    string
		target review.Target
		stmt   int
	}{
		{"RESET ROLE;\nCREATE TABLE x (a int);", tenant, -1},
		{"SET ROLE NONE;\nCREATE TABLE x (a int);", tenant, -1},
		{"CREATE TABLE x (a int);", tenant, 0},
		{"RESET ROLE;\nCREATE TABLE g (a int);", plain, 1},
		{"GRANT SELECT ON g TO PUBLIC;", tenant, -1},
		{"GRANT SELECT ON g TO PUBLIC;", plain, -1},
		{"GRANT SELECT ON x TO PUBLIC;\nCREATE TABLE x (a int);", tenant, 1},
		// After GRANT the role is the owner again, in tenant mode.
		{"SET ROLE bob;\nCREATE TABLE z (a int);\nGRANT SELECT ON t TO PUBLIC;\nCREATE TABLE y (a int);", tenant, -1},
		{"SET ROLE bob;\nCREATE TABLE y (a int);", tenant, 1},
		{"-- txn-mode = off\nSET ROLE admin;\nGRANT SELECT ON t TO PUBLIC;\nCREATE TABLE g (a int);", tenant, 2},
	} {
		got := walkFindings(t, tc.sql, tc.target)
		switch {
		case tc.stmt < 0 && len(got) != 0:
			t.Errorf("%q as %s: findings = %+v, want none", tc.sql, tc.target.SessionUser, got)
		case tc.stmt >= 0 && (len(got) != 1 || got[0].Statement != tc.stmt):
			t.Errorf("%q as %s: findings = %+v, want one on statement %d", tc.sql, tc.target.SessionUser, got, tc.stmt)
		}
	}
}

// PostgreSQL 12 adds the constraints of ADD COLUMN IF NOT EXISTS of a
// column that exists; 17 and 18 do not (verified on all three), so their
// names may or may not be taken.
func TestWalkThroughAddColumnIfNotExistsConstraints(t *testing.T) {
	target := review.Target{Schema: walkSchema(), SessionUser: "u"}
	for _, sql := range []string{
		"ALTER TABLE t ADD COLUMN IF NOT EXISTS a int CONSTRAINT uq UNIQUE;\nALTER INDEX uq RENAME TO uq2;",
		"ALTER TABLE t ADD COLUMN b int, ADD COLUMN IF NOT EXISTS b int CONSTRAINT uq UNIQUE;\nDROP INDEX uq;",
		"ALTER TABLE t ADD COLUMN IF NOT EXISTS a serial;\nDROP SEQUENCE t_a_seq;",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}
}

// A finding about a column or a constraint anchors at it.
func TestWalkThroughAnchorsColumnsAndConstraints(t *testing.T) {
	schema := walkSchema()
	schema.Schemas[0].Tables[0].Indexes = []*metadata.IndexMetadata{{Name: "t_pkey", Primary: true, IsConstraint: true, Unique: true}}
	schema.Schemas[0].Tables[0].CheckConstraints = []*metadata.CheckConstraintMetadata{{Name: "ck"}}
	target := review.Target{Schema: schema, SessionUser: "u"}
	for sql, want := range map[string]string{
		"SELECT 1;\nCREATE TABLE x (a int, b int, a text);":                                      "a text",
		"SELECT 1;\nCREATE TABLE x (a int, ctid int);":                                           "ctid int",
		"SELECT 1;\nCREATE TABLE x (a int PRIMARY KEY, b int, PRIMARY KEY (b));":                 "PRIMARY KEY (b)",
		"SELECT 1;\nCREATE TABLE x (a int, CONSTRAINT c CHECK (a>0), CONSTRAINT c CHECK (a>1));": "CONSTRAINT c CHECK (a>1)",
		"SELECT 1;\nCREATE TABLE x (a int CONSTRAINT x UNIQUE);":                                 "CONSTRAINT x UNIQUE",
		"SELECT 1;\nALTER TABLE t ADD COLUMN a int;":                                             "a int",
		"SELECT 1;\nALTER TABLE t ADD PRIMARY KEY (a);":                                          "PRIMARY KEY (a)",
		"SELECT 1;\nALTER TABLE t ADD CONSTRAINT ck CHECK (a > 0);":                              "CONSTRAINT ck CHECK (a > 0)",
		"SELECT 1;\nALTER TABLE t ADD COLUMN k int PRIMARY KEY;":                                 "PRIMARY KEY",
		"SELECT 1;\nALTER TABLE t ADD COLUMN xmin int;":                                          "xmin int",
	} {
		got := walkFindings(t, sql, target)
		if len(got) != 1 || sql[got[0].Range.Start:got[0].Range.End] != want {
			t.Errorf("%q: findings = %+v, want one at %q", sql, got, want)
		}
	}
}

// A name outside ASCII truncates where the server's encoding says: two
// that collide in UTF8 do not in LATIN1 (verified on PostgreSQL 17). The
// walk never reads such a name, so it claims nothing either way.
func TestWalkThroughNonASCIINames(t *testing.T) {
	long := strings.Repeat("é", 40)
	target := review.Target{Schema: walkSchema(), SessionUser: "u"}
	for _, sql := range []string{
		"CREATE TABLE \"" + long + "\" (a int);\nCREATE TABLE \"" + long + "x\" (a int);",
		"DROP TABLE \"" + long + "\";",
		"CREATE TABLE \"" + long + "\" (a int);\nDROP TABLE \"" + long + "x\";\nDROP TABLE \"" + long + "\";",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}
}
