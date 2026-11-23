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
	checkCodes(t, result.Findings)
	return result.Findings
}

func TestWalkThroughReportsTheFirstFailure(t *testing.T) {
	got := walkFindings(t, "CREATE TABLE x (a int);\nALTER TABLE public.nosuch ADD COLUMN b int;\nCREATE TABLE t (a int);", review.Target{Schema: walkSchema(), SessionUser: "u"})
	if len(got) != 1 || got[0].Statement != 1 || got[0].Message != "`public.nosuch` doesn't exist." {
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
	if len(got) != 1 || got[0].Statement != 1 || got[0].Message != "Table `"+schema.Name+".public.x` already exists." {
		t.Errorf("findings = %+v, want the second CREATE's", got)
	}
	got = walkFindings(t, "CREATE TABLE other.public.x (a int);", review.Target{Schema: schema, SessionUser: "u"})
	if len(got) != 1 || got[0].Message != "`other.public.x` refers to another database, which PostgreSQL doesn't allow." {
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
	if len(got) != 1 || got[0].Message != "PostgreSQL doesn't accept `CREATE INDEX CONCURRENTLY` on partitioned table `parted`." {
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
		"DROP INDEX CONCURRENTLY IF EXISTS x, y;":      "PostgreSQL doesn't accept `DROP INDEX CONCURRENTLY` of several indexes.",
		"DROP INDEX CONCURRENTLY IF EXISTS x CASCADE;": "PostgreSQL doesn't accept `DROP INDEX CONCURRENTLY` with `CASCADE`.",
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
		{"DROP INDEX CONCURRENTLY t_a_idx;\nCREATE INDEX t_a_idx ON t (a);", 1, "Index `t_a_idx` already exists."},
		{"CREATE INDEX CONCURRENTLY i ON t (a);\nCREATE INDEX i ON t (a);\nCREATE INDEX CONCURRENTLY i2 ON t (a);", 0, "Index `i` already exists."},
		// Matched anywhere in its text, as Bytebase matches it.
		{"COMMENT ON TABLE t2 IS 'no DROP DATABASE here';\nCREATE TABLE t2 (a int);", -1, ""},
		// First, in auto-commit mode.
		{"-- txn-mode = off\nCREATE TABLE t_a_idx (a int);\nDROP INDEX CONCURRENTLY t_a_idx;", -1, ""},
		{"-- txn-mode = off\nCREATE TABLE i (a int);\nCREATE INDEX CONCURRENTLY i ON t (a);", 0, "The name `i` is already taken."},
		{"-- txn-isolation = SERIALIZABLE\n--   TXN-MODE=Off\nCREATE TABLE t_a_idx (a int);\nDROP INDEX CONCURRENTLY t_a_idx;", -1, ""},
		{"-- txn-mode = off\n-- txn-mode = on\nCREATE TABLE t_a_idx (a int);\nDROP INDEX CONCURRENTLY t_a_idx;", 0, "The name `t_a_idx` is already taken."},
		// A directive below the first statement is no directive.
		{"CREATE TABLE t_a_idx (a int);\n-- txn-mode = off\nDROP INDEX CONCURRENTLY t_a_idx;", 0, "The name `t_a_idx` is already taken."},
		// A CONCURRENTLY the patterns miss fails in the transaction: the
		// walk ends before it.
		{"CREATE /* x */ INDEX CONCURRENTLY t_a_idx ON t (a);\nDROP TABLE nosuch;", -1, ""},
		{"VACUUM t;\nDROP TABLE nosuch;", 1, "Table `nosuch` doesn't exist."},
		{"/* x */ VACUUM t;\nDROP TABLE nosuch;", -1, ""},
		{"-- txn-mode = off\n/* x */ VACUUM t;\nDROP TABLE nosuch;", 1, "Table `nosuch` doesn't exist."},
		{"BEGIN;\nREINDEX TABLE t;\nDROP TABLE nosuch;", -1, ""},
		// A second run of what is no SET ends the walk.
		{"COMMENT ON TABLE t IS 'SET ROLE';\nCREATE INDEX CONCURRENTLY i ON t (a);\nDROP TABLE nosuch;", 2, "Table `nosuch` doesn't exist."},
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

// ALTER TYPE's attribute forms alter a composite type, and only one:
// PostgreSQL refuses them on a table or view with 42809 "is not a
// composite type", and ALTER TABLE on a composite type with "is a
// composite type" (verified on PostgreSQL 17).
func TestWalkThroughAlterTypeTakesACompositeType(t *testing.T) {
	schema := walkSchema()
	schema.Schemas[0].CompositeTypes = []*metadata.CompositeTypeMetadata{{Name: "pair"}}
	target := review.Target{Schema: schema, SessionUser: "u"}
	for _, sql := range []string{
		"ALTER TYPE pair ADD ATTRIBUTE b int;",
		"CREATE TYPE p2 AS (a int);\nALTER TYPE p2 ADD ATTRIBUTE b int, ALTER ATTRIBUTE a TYPE text;",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}
	for sql, want := range map[string]string{
		"ALTER TYPE t ADD ATTRIBUTE b int;":                                  "`t` is a table, not a composite type.",
		"CREATE VIEW v AS SELECT 1 AS a;\nALTER TYPE v ADD ATTRIBUTE b int;": "`v` is a view, not a composite type.",
		"ALTER TABLE pair ADD COLUMN b int;":                                 "`pair` is a composite type, not a table.",
	} {
		got := walkFindings(t, sql, target)
		if len(got) != 1 || got[0].Message != want || !slices.Equal(got[0].Codes, []review.FindingCode{review.WrongObjectKind}) {
			t.Errorf("%q: findings = %+v, want one %s saying %q", sql, got, review.WrongObjectKind, want)
		}
	}
}

func TestWalkThroughTypesDomainsEnumsAndRanges(t *testing.T) {
	schema := walkSchema()
	schema.Schemas[0].Views = []*metadata.ViewMetadata{{Name: "v"}}
	schema.Schemas[0].EnumTypes = []*metadata.EnumTypeMetadata{{Name: "mood"}}
	schema.Schemas[0].CompositeTypes = []*metadata.CompositeTypeMetadata{{Name: "pair"}}
	schema.Schemas = append(schema.Schemas, &metadata.SchemaMetadata{Name: "s"})
	target := review.Target{Schema: schema, SessionUser: "u"}

	// Positive findings: built-in drops, row/array drops, wrong type kinds,
	// shell/base/range/enum/composite errors, sequence-type collisions.
	for sql, wantCode := range map[string]review.FindingCode{
		"DROP TYPE int;":                                       review.UnsupportedForm,
		"DROP TYPE IF EXISTS integer;":                         review.UnsupportedForm,
		"DROP TYPE character(1);":                              review.UnsupportedForm,
		"DROP TYPE bit(a);":                                    review.UnsupportedForm,
		"DROP TYPE t;":                                         review.UnsupportedForm,
		"DROP TYPE IF EXISTS t CASCADE;":                       review.UnsupportedForm,
		"DROP TYPE mood[];":                                    review.UnsupportedForm,
		"DROP TYPE mood(1);":                                   review.UnsupportedForm,
		"DROP DOMAIN int;":                                     review.WrongObjectKind,
		"DROP DOMAIN IF EXISTS t;":                             review.WrongObjectKind,
		"DROP DOMAIN mood;":                                    review.WrongObjectKind,
		"ALTER DOMAIN t SET NOT NULL;":                         review.WrongObjectKind,
		"ALTER DOMAIN mood DROP DEFAULT;":                      review.WrongObjectKind,
		"ALTER DOMAIN t RENAME TO t2;":                         review.WrongObjectKind,
		"ALTER DOMAIN t RENAME CONSTRAINT c TO c2;":            review.WrongObjectKind,
		"ALTER DOMAIN t SET SCHEMA s;":                         review.WrongObjectKind,
		"COMMENT ON DOMAIN t IS 'x';":                          review.WrongObjectKind,
		"COMMENT ON CONSTRAINT c ON DOMAIN t IS 'x';":          review.WrongObjectKind,
		"ALTER TYPE t ADD VALUE 'x';":                          review.WrongObjectKind,
		"ALTER TYPE pair RENAME VALUE 'a' TO 'b';":             review.WrongObjectKind,
		"CREATE DOMAIN d AS int;\nALTER TYPE d ADD VALUE 'x';": review.WrongObjectKind,
		"ALTER TYPE mood SET (bogus = 1);":                     review.UnsupportedForm,
		"ALTER TYPE t RENAME TO t2;":                           review.WrongObjectKind,
		"ALTER TYPE t SET SCHEMA s;":                           review.WrongObjectKind,
		"CREATE TABLE x OF mood;":                              review.WrongObjectKind,
		"CREATE TABLE x OF t;":                                 review.WrongObjectKind,
		"CREATE TYPE mood AS ENUM ('a');":                      review.NameTaken,
		"CREATE TYPE e2 AS ENUM ('a', 'a');":                   review.DuplicateName,
		"CREATE TYPE c2 AS (x int, x text);":                   review.DuplicateName,
		"CREATE TYPE mood AS RANGE (SUBTYPE = int);":           review.NameTaken,
		"CREATE TYPE r2 AS RANGE (COLLATION = \"C\");":         review.UnsupportedForm,
		"CREATE TYPE mood;":                                    review.NameTaken,
		"CREATE TYPE other.public.bt;":                         review.CrossDatabaseReference,
		"CREATE TYPE nosuch.bt;":                               review.ObjectNotFound,
		"CREATE TYPE bt (INPUT = f_in);":                       review.UnsupportedForm,
		"CREATE SEQUENCE mood;":                                review.NameTaken,
		"CREATE SEQUENCE IF NOT EXISTS mood;":                  review.NameTaken,
		"CREATE DOMAIN d AS int;\nDROP DOMAIN d;\nALTER DOMAIN d SET NOT NULL;": review.ObjectNotFound,
	} {
		got := walkFindings(t, sql, target)
		if len(got) != 1 || !slices.Equal(got[0].Codes, []review.FindingCode{wantCode}) {
			t.Errorf("%q: findings = %+v, want one %s", sql, got, wantCode)
		}
	}

	// Control cases (0 false positives): unsynced domains/ranges/base types,
	// valid enum/composite/range/shell type operations, search_path shadowing.
	for _, sql := range []string{
		"DROP TYPE unsynced_domain;",
		"DROP TYPE IF EXISTS unsynced_domain;",
		"DROP DOMAIN unsynced_domain;",
		"DROP DOMAIN IF EXISTS unsynced_domain;",
		"ALTER DOMAIN unsynced_domain SET NOT NULL;",
		"ALTER TYPE unsynced_enum ADD VALUE 'x';",
		"ALTER TYPE mood ADD VALUE 'c';",
		"ALTER TYPE mood RENAME TO mood2;",
		"ALTER TYPE mood SET SCHEMA s;",
		"ALTER TYPE pair RENAME TO pair2;",
		"ALTER TYPE pair SET SCHEMA s;",
		"DROP TYPE pair;",
		"DROP TYPE mood;",
		"CREATE TABLE typed OF pair;",
		"CREATE TYPE mybase;\nCREATE TYPE mybase (INPUT = f_in, OUTPUT = f_out);",
		"CREATE TYPE public.bt;\nCREATE TYPE public.bt (INPUT = f_in, OUTPUT = f_out, TYPMOD_IN = f_modin);\nDROP TYPE public.bt(3) CASCADE;",
		"CREATE TYPE public.bt2;\nCREATE TYPE public.bt2 (INPUT = f_in, OUTPUT = f_out);\nALTER TYPE public.bt2 SET (TYPMOD_IN = f_modin);\nDROP TYPE public.bt2(3) CASCADE;",
		"CREATE TYPE public.myrange;\nCREATE TYPE public.myrange AS RANGE (SUBTYPE = int4);",
		"CREATE TYPE public.r1;\nALTER TYPE public.r1 RENAME TO r2;\nCREATE TYPE public.r2 AS RANGE (SUBTYPE = int4);",
		"CREATE TYPE public.r3;\nALTER TYPE public.r3 SET SCHEMA s;\nCREATE TYPE s.r3 AS RANGE (SUBTYPE = int4);",
		"CREATE SCHEMA s_mr;\nCREATE TYPE s_mr.r AS RANGE (SUBTYPE = int4, MULTIRANGE_TYPE_NAME = s_mr.m);\nCOMMENT ON TYPE s_mr.m IS 'ok';",
		"CREATE SCHEMA s_mr2;\nCREATE TYPE s_mr2.r AS RANGE (SUBTYPE = int4);\nCOMMENT ON TYPE s_mr2.r_multirange IS 'ok';",
		"CREATE TYPE public.e_move AS ENUM ('a');\nCREATE SCHEMA dest_mv;\nALTER TYPE public.e_move SET SCHEMA dest_mv;\nCOMMENT ON TYPE dest_mv._e_move IS 'ok';",
		"CREATE DOMAIN public.d_move AS int;\nCREATE SCHEMA dest_dmv;\nALTER DOMAIN public.d_move SET SCHEMA dest_dmv;\nCOMMENT ON TYPE dest_dmv._d_move IS 'ok';",
		"CREATE SCHEMA s_arr;\nSET search_path = s_arr;\nCOMMENT ON TYPE _int4 IS 'ok';",
		"ALTER DOMAIN unsynced_domain ADD CONSTRAINT ck CHECK (public.zap(VALUE));\nCREATE TABLE t (a int);",
		"ALTER DOMAIN unsynced_domain VALIDATE CONSTRAINT ck;\nCREATE TABLE t (a int);",
		"ALTER TYPE unsynced_base SET (storage = extended);",
		"CREATE TYPE myrange AS RANGE (SUBTYPE = int4);",
		"CREATE DOMAIN d AS int;\nALTER DOMAIN d SET NOT NULL;\nALTER DOMAIN d RENAME TO d2;\nALTER DOMAIN d2 SET SCHEMA s;\nDROP DOMAIN s.d2;",
		// An earlier schema on search_path may hold an unsynced domain shadowing public.t.
		"SET search_path = s, public;\nDROP DOMAIN t;\nALTER DOMAIN t SET NOT NULL;",
		// Explicit public before pg_catalog may shadow a built-in type with an unsynced domain.
		"SET search_path = public, pg_catalog;\nDROP TYPE text;",
		"SET search_path = public, pg_catalog;\nDROP TYPE \"json\";",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}
}

func TestWalkThroughFunctionsTriggersPoliciesAndAlterSubcommands(t *testing.T) {
	schema := walkSchema()
	schema.Schemas[0].Tables[0].Indexes = []*metadata.IndexMetadata{{Name: "t_a_idx"}}
	schema.Schemas[0].Tables[0].Triggers = []*metadata.TriggerMetadata{{Name: "tr_existing"}}
	schema.Schemas[0].Views = []*metadata.ViewMetadata{{Name: "v"}}
	schema.Schemas[0].MaterializedViews = []*metadata.MaterializedViewMetadata{{Name: "mv"}}
	schema.Schemas[0].Sequences = []*metadata.SequenceMetadata{{Name: "q"}}
	schema.Schemas[0].Functions = []*metadata.FunctionMetadata{{Name: "unsynced_fn"}}
	target := review.Target{Schema: schema, SessionUser: "u"}

	for sql, wantCode := range map[string]review.FindingCode{
		// Functions / procedures.
		"CREATE FUNCTION f();":                                                           review.UnsupportedForm,
		"CREATE FUNCTION f() RETURNS int;":                                               review.UnsupportedForm,
		"CREATE FUNCTION f() RETURNS int AS 'SELECT 1';":                                 review.UnsupportedForm,
		"CREATE FUNCTION f() LANGUAGE sql AS 'SELECT 1';":                                review.UnsupportedForm,
		"CREATE FUNCTION f() RETURNS int LANGUAGE sql WINDOW WINDOW AS 'SELECT 1';":      review.UnsupportedForm,
		"CREATE PROCEDURE p() LANGUAGE sql STRICT AS 'SELECT 1';":                        review.UnsupportedForm,
		"CREATE FUNCTION f(SETOF int) RETURNS int LANGUAGE sql AS 'SELECT 1';":           review.UnsupportedForm,
		"CREATE FUNCTION f(VARIADIC int[], int) RETURNS int LANGUAGE sql AS 'SELECT 1';": review.UnsupportedForm,
		"CREATE FUNCTION nosuch.f() RETURNS int LANGUAGE sql AS 'SELECT 1';":             review.ObjectNotFound,
		"ALTER FUNCTION f() STRICT STRICT;":                                              review.UnsupportedForm,
		"ALTER FUNCTION f() PARALLEL bogus;":                                             review.UnsupportedForm,
		"DROP FUNCTION nosuch.f();":                                                      review.ObjectNotFound,
		// Triggers.
		"CREATE TRIGGER tr INSTEAD OF INSERT ON t FOR EACH ROW EXECUTE FUNCTION f();":                  review.WrongObjectKind,
		"CREATE TRIGGER tr BEFORE INSERT ON v FOR EACH ROW EXECUTE FUNCTION f();":                      review.WrongObjectKind,
		"CREATE TRIGGER tr BEFORE INSERT ON q FOR EACH ROW EXECUTE FUNCTION f();":                      review.WrongObjectKind,
		"CREATE TRIGGER tr BEFORE INSERT ON t_a_idx FOR EACH ROW EXECUTE FUNCTION f();":                review.WrongObjectKind,
		"CREATE TRIGGER tr BEFORE TRUNCATE ON t FOR EACH ROW EXECUTE FUNCTION f();":                    review.UnsupportedForm,
		"CREATE TRIGGER tr INSTEAD OF INSERT ON v FOR EACH STATEMENT EXECUTE FUNCTION f();":            review.UnsupportedForm,
		"CREATE TRIGGER tr BEFORE INSERT ON t REFERENCING NEW ROW AS n EXECUTE FUNCTION f();":          review.UnsupportedForm,
		"CREATE TRIGGER tr BEFORE INSERT ON t REFERENCING NEW TABLE AS n EXECUTE FUNCTION f();":        review.UnsupportedForm,
		"CREATE TRIGGER tr_existing BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION f();":             review.NameTaken,
		"CREATE OR REPLACE CONSTRAINT TRIGGER tr AFTER INSERT ON t FOR EACH ROW EXECUTE FUNCTION f();": review.UnsupportedForm,
		"DROP TRIGGER tr ON nosuch;":           review.ObjectNotFound,
		"DROP TRIGGER tr ON t_a_idx;":          review.WrongObjectKind,
		"DROP TRIGGER tr ON mv;":               review.ObjectNotFound,
		"ALTER TRIGGER tr ON q RENAME TO tr2;": review.WrongObjectKind,
		// Policies.
		"CREATE POLICY p ON v;":                              review.WrongObjectKind,
		"CREATE POLICY p ON t FOR SELECT WITH CHECK (true);": review.UnsupportedForm,
		"CREATE POLICY p ON t FOR INSERT USING (true);":      review.UnsupportedForm,
		"CREATE POLICY p ON t;\nCREATE POLICY p ON t;":       review.NameTaken,
		"ALTER POLICY p ON nosuch;":                          review.ObjectNotFound,
		"ALTER POLICY p ON t_a_idx;":                         review.WrongObjectKind,
		"DROP POLICY p ON v;":                                review.ObjectNotFound,
		"DROP POLICY IF EXISTS p ON t_a_idx;":                review.WrongObjectKind,
		"ALTER POLICY p ON v RENAME TO p2;":                  review.WrongObjectKind,
		"COMMENT ON POLICY p ON t_a_idx IS 'x';":             review.WrongObjectKind,
		"COMMENT ON POLICY p ON v IS 'x';":                   review.ObjectNotFound,
		// ALTER TABLE per-subcommand relkinds & targets.
		"ALTER TABLE v ADD COLUMN b int;":                 review.WrongObjectKind,
		"ALTER TABLE q ADD COLUMN b int;":                 review.WrongObjectKind,
		"ALTER TABLE t_a_idx DROP COLUMN a;":              review.WrongObjectKind,
		"ALTER TABLE t OPTIONS (foo 'bar');":              review.WrongObjectKind,
		"ALTER TABLE t CLUSTER ON t;":                     review.WrongObjectKind,
		"ALTER TABLE t CLUSTER ON nosuch_idx;":            review.ObjectNotFound,
		"ALTER TABLE t REPLICA IDENTITY USING INDEX t;":   review.WrongObjectKind,
		"ALTER TABLE t OF t;":                             review.WrongObjectKind,
		"ALTER TABLE t INHERIT t;":                        review.UnsupportedForm,
		"ALTER TABLE t INHERIT v;":                        review.WrongObjectKind,
		"ALTER TABLE t ATTACH PARTITION t DEFAULT;":       review.UnsupportedForm,
		"ALTER TABLE t DETACH PARTITION t;":               review.UnsupportedForm,
		"ALTER TABLE t ADD COLUMN self t;":                review.UnsupportedForm,
		"ALTER TABLE t ALTER COLUMN a TYPE t;":            review.UnsupportedForm,
		"ALTER TABLE t ADD UNIQUE (nosuch);":              review.ObjectNotFound,
		"ALTER TABLE t ALTER 1 SET STATISTICS 10;":        review.UnsupportedForm,
		"ALTER TABLE t ALTER COLUMN a SET STORAGE bogus;": review.UnsupportedForm,
		"ALTER TABLE q RENAME COLUMN a TO b;":             review.WrongObjectKind,
		"ALTER VIEW v RENAME COLUMN a TO a;":              review.DuplicateName,
		"ALTER TABLE v RENAME CONSTRAINT c TO c2;":        review.WrongObjectKind,
		// Indexes, views, CTAS, sequences, comments, depends.
		"CREATE INDEX ON t (nosuch);":                       review.ObjectNotFound,
		"CREATE INDEX ON t ((CURRENT_DATE));":               review.UnsupportedForm,
		"CREATE INDEX ON t (a) INCLUDE ((a + 1));":          review.UnsupportedForm,
		"CREATE VIEW v2 AS SELECT 1 WITH CHECK OPTION;":     review.UnsupportedForm,
		"CREATE VIEW v2 (x, y) AS SELECT 1;":                review.UnsupportedForm,
		"CREATE TABLE t2 ON COMMIT DROP AS SELECT 1;":       review.UnsupportedForm,
		"CREATE TABLE t AS EXECUTE prep;":                   review.NameTaken,
		"ALTER SEQUENCE q CYCLE CYCLE;":                     review.UnsupportedForm,
		"ALTER SEQUENCE q INCREMENT 1.5;":                   review.UnsupportedForm,
		"COMMENT ON COLUMN a IS 'x';":                       review.UnsupportedForm,
		"COMMENT ON CONSTRAINT c ON nosuch IS 'x';":         review.ObjectNotFound,
		"COMMENT ON CONSTRAINT c ON v IS 'x';":              review.ObjectNotFound,
		"COMMENT ON CONSTRAINT c ON t_a_idx IS 'x';":        review.WrongObjectKind,
		"ALTER MATERIALIZED VIEW t DEPENDS ON EXTENSION e;": review.WrongObjectKind,
		"ALTER INDEX t DEPENDS ON EXTENSION e;":             review.WrongObjectKind,
	} {
		got := walkFindings(t, sql, target)
		if len(got) != 1 || !slices.Equal(got[0].Codes, []review.FindingCode{wantCode}) {
			t.Errorf("%q: findings = %+v, want one %s", sql, got, wantCode)
		}
	}

	// Control cases (0 false positives).
	for _, sql := range []string{
		"CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT 1';",
		"CREATE FUNCTION f(OUT x int) LANGUAGE sql AS 'SELECT 1';",
		"CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT 1; END;",
		"CREATE PROCEDURE p() LANGUAGE sql AS 'SELECT 1';",
		"CREATE FUNCTION f() RETURNS int LANGUAGE plv8 AS 'return 1;';",
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION unsynced_fn();",
		"CREATE OR REPLACE TRIGGER tr_existing BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION unsynced_fn();",
		"CREATE TRIGGER tr INSTEAD OF INSERT ON v FOR EACH ROW EXECUTE FUNCTION unsynced_fn();",
		"CREATE TRIGGER tr AFTER INSERT ON t REFERENCING NEW TABLE AS nt FOR EACH STATEMENT EXECUTE FUNCTION unsynced_fn();",
		"CREATE CONSTRAINT TRIGGER ctr AFTER INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION unsynced_fn();\nALTER TABLE public.t DROP CONSTRAINT ctr;\nCREATE CONSTRAINT TRIGGER ctr AFTER INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION unsynced_fn();",
		"ALTER TRIGGER tr_existing ON t RENAME TO tr_existing;",
		"DROP TRIGGER IF EXISTS tr ON mv;",
		"DROP POLICY IF EXISTS p ON v;",
		"DROP POLICY IF EXISTS p ON mv;",
		"DROP POLICY IF EXISTS p ON q;",
		"DROP POLICY unsynced_pol ON t;",
		"ALTER POLICY unsynced_pol ON t USING (true);",
		"CREATE INDEX IF NOT EXISTS t_a_idx ON t (nosuch);",
		"CREATE INDEX IF NOT EXISTS t_a_idx ON t ((CURRENT_DATE));",
		"CREATE INDEX ON t ((CASE WHEN false THEN CURRENT_DATE ELSE DATE '2000-01-01' END));",
		"CREATE INDEX ON t (a);",
		"CREATE VIEW v2 AS SELECT a FROM t WITH CHECK OPTION;",
		"CREATE VIEW v2 (x, y) AS SELECT * FROM t;",
		"CREATE TEMP TABLE tmp ON COMMIT DROP AS SELECT 1;",
		"CREATE TABLE t2 AS SELECT 1;",
		"ALTER TABLE t ADD COLUMN self_arr t[];",
		"ALTER TABLE t ADD COLUMN b int, ADD UNIQUE (b);",
		"ALTER TABLE t ADD CONSTRAINT uq UNIQUE (b), ADD COLUMN b int;",
		"CREATE TABLE pt (k int) PARTITION BY RANGE (k);\nCREATE TABLE pt_1 PARTITION OF pt FOR VALUES FROM (0) TO (10);\nCREATE TRIGGER tr BEFORE INSERT ON pt FOR EACH ROW EXECUTE FUNCTION unsynced_fn();\nCOMMENT ON TRIGGER tr ON pt_1 IS 'ok';",
		"CREATE TABLE pt (k int) PARTITION BY RANGE (k);\nCREATE TRIGGER tr BEFORE INSERT ON pt FOR EACH ROW EXECUTE FUNCTION unsynced_fn();\nCREATE TABLE pt_1 PARTITION OF pt FOR VALUES FROM (0) TO (10);\nCOMMENT ON TRIGGER tr ON pt_1 IS 'ok';",
		"ALTER TABLE t CLUSTER ON t_a_idx;",
		"ALTER TABLE t ALTER COLUMN a SET STORAGE extended;",
		"ALTER INDEX t_a_idx ALTER COLUMN 1 SET STATISTICS 10;",
		"CREATE TABLE pt_idx (k int) PARTITION BY RANGE (k);\nCREATE TABLE pt_idx_1 PARTITION OF pt_idx FOR VALUES FROM (0) TO (10);\nCREATE INDEX parent_idx ON ONLY pt_idx (k);\nCREATE INDEX child_idx ON pt_idx_1 (k);\nALTER INDEX public.parent_idx ATTACH PARTITION public.child_idx;",
		"ALTER SEQUENCE q INCREMENT 2 NO CYCLE;",
		"ALTER MATERIALIZED VIEW mv DEPENDS ON EXTENSION e;",
		"ALTER INDEX t_a_idx DEPENDS ON EXTENSION e;",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}

	// Unresolved targets (when SessionUser is omitted so "$user" is unknown)
	// weaken policy and trigger facts rather than causing false duplicates.
	noUser := review.Target{Schema: schema}
	for _, sql := range []string{
		"CREATE POLICY p ON public.t;\nDROP POLICY p ON t;\nCREATE POLICY p ON public.t;",
		"CREATE TRIGGER tr BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION unsynced_fn();\nDROP TRIGGER tr ON t;\nCREATE TRIGGER tr BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION unsynced_fn();",
		"CREATE POLICY p ON public.t;\nALTER POLICY p ON t RENAME TO p2;\nCREATE POLICY p ON public.t;",
		"CREATE TRIGGER tr BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION unsynced_fn();\nALTER TRIGGER tr ON t RENAME TO tr2;\nCREATE TRIGGER tr BEFORE INSERT ON public.t FOR EACH ROW EXECUTE FUNCTION unsynced_fn();",
		"CREATE SCHEMA dest;\nCREATE TYPE public.c AS (a int);\nALTER TYPE c SET SCHEMA dest;\nCREATE TYPE public.c AS (a int);",
		"CREATE TYPE public.bt;\nCREATE TYPE public.bt (INPUT = unsynced_in, OUTPUT = unsynced_out);\nALTER TYPE bt SET (TYPMOD_IN = unsynced_modin);\nDROP TYPE public.bt(1);",
	} {
		if got := walkFindings(t, sql, noUser); len(got) != 0 {
			t.Errorf("%q (no SessionUser): findings = %+v, want none", sql, got)
		}
	}
}

func TestWalkThroughRoutinesExcludedMiddleAndOptions(t *testing.T) {
	schema := walkSchema()
	schema.Schemas[0].Tables[0].Indexes = []*metadata.IndexMetadata{{Name: "t_a_idx"}}
	schema.Schemas[0].Tables[0].Triggers = []*metadata.TriggerMetadata{{Name: "tr_existing"}}
	schema.Schemas[0].Tables = append(schema.Schemas[0].Tables, &metadata.TableMetadata{
		Name:       "parted",
		Columns:    []*metadata.ColumnMetadata{{Name: "k"}},
		Partitions: []*metadata.TablePartitionMetadata{{Name: "parted_1"}},
	})
	schema.Schemas[0].Views = []*metadata.ViewMetadata{{
		Name:    "v",
		Columns: []*metadata.ColumnMetadata{{Name: "a"}, {Name: "b"}},
	}}
	schema.Schemas[0].MaterializedViews = []*metadata.MaterializedViewMetadata{{Name: "mv"}}
	schema.Schemas[0].Sequences = []*metadata.SequenceMetadata{{Name: "q"}}
	schema.Schemas[0].CompositeTypes = []*metadata.CompositeTypeMetadata{{Name: "pair"}}
	schema.Schemas[0].EnumTypes = []*metadata.EnumTypeMetadata{{Name: "mood"}}
	schema.Schemas[0].Functions = []*metadata.FunctionMetadata{{Name: "fn_existing"}, {Name: "trig_fn"}}
	schema.Schemas[0].Procedures = []*metadata.ProcedureMetadata{{Name: "proc_existing"}}
	target := review.Target{Schema: schema, SessionUser: "u"}

	for sql, wantCode := range map[string]review.FindingCode{
		// Missing routines across ALTER/DROP/COMMENT/OWNER/DEPENDS/TRIGGER.
		"ALTER FUNCTION missing_fn STRICT;":                                             review.ObjectNotFound,
		"ALTER PROCEDURE missing_proc STRICT;":                                          review.ObjectNotFound,
		"ALTER ROUTINE missing_rtn STRICT;":                                             review.ObjectNotFound,
		"ALTER PROCEDURE fn_existing STRICT;":                                           review.ObjectNotFound,
		"ALTER FUNCTION proc_existing STRICT;":                                          review.ObjectNotFound,
		"DROP FUNCTION missing_fn;":                                                     review.ObjectNotFound,
		"DROP PROCEDURE missing_proc;":                                                  review.ObjectNotFound,
		"DROP ROUTINE missing_rtn;":                                                     review.ObjectNotFound,
		"COMMENT ON FUNCTION missing_fn IS 'x';":                                        review.ObjectNotFound,
		"COMMENT ON PROCEDURE missing_proc IS 'x';":                                     review.ObjectNotFound,
		"COMMENT ON ROUTINE missing_rtn IS 'x';":                                        review.ObjectNotFound,
		"ALTER FUNCTION missing_fn RENAME TO fn2;":                                      review.ObjectNotFound,
		"ALTER PROCEDURE missing_proc SET SCHEMA public;":                               review.ObjectNotFound,
		"ALTER FUNCTION fn_existing SET SCHEMA nosuch;":                                 review.ObjectNotFound,
		"ALTER FUNCTION missing_fn OWNER TO bob;":                                       review.ObjectNotFound,
		"ALTER PROCEDURE missing_proc OWNER TO bob;":                                    review.ObjectNotFound,
		"ALTER FUNCTION missing_fn DEPENDS ON EXTENSION e;":                             review.ObjectNotFound,
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION no_fn();":   review.ObjectNotFound,
		"CREATE TRIGGER tr BEFORE INSERT ON t WHEN (a > 0) EXECUTE FUNCTION trig_fn();": review.UnsupportedForm,
		"ALTER TABLE t ENABLE TRIGGER missing_trig;":                                    review.ObjectNotFound,
		"ALTER TABLE t DISABLE TRIGGER missing_trig;":                                   review.ObjectNotFound,
		"DROP TRIGGER missing_trig ON t;":                                               review.ObjectNotFound,
		"ALTER TRIGGER missing_trig ON t RENAME TO tr2;":                                review.ObjectNotFound,
		"ALTER TRIGGER missing_trig ON t DEPENDS ON EXTENSION e;":                       review.ObjectNotFound,
		// Self-referencing creations (law of excluded middle).
		"CREATE DOMAIN d AS d;":               review.UnsupportedForm,
		"CREATE DOMAIN public.d AS public.d;": review.UnsupportedForm,
		"CREATE TABLE x (c x);":               review.UnsupportedForm,
		"CREATE TABLE x OF x;":                review.UnsupportedForm,
		"CREATE TYPE c AS (a c);":             review.UnsupportedForm,
		// Self-renames and enum self-values.
		"ALTER DOMAIN unsynced_dom RENAME TO unsynced_dom;":             review.DuplicateName,
		"ALTER DOMAIN unsynced_dom RENAME CONSTRAINT c TO c;":           review.DuplicateName,
		"ALTER TYPE unsynced_typ RENAME TO unsynced_typ;":               review.DuplicateName,
		"ALTER TABLE t RENAME CONSTRAINT unsynced_con TO unsynced_con;": review.DuplicateName,
		"ALTER POLICY unsynced_pol ON t RENAME TO unsynced_pol;":        review.DuplicateName,
		"ALTER TYPE pair RENAME ATTRIBUTE x TO x;":                      review.DuplicateName,
		"ALTER TYPE q RENAME ATTRIBUTE x TO y;":                         review.WrongObjectKind,
		"ALTER TYPE t_a_idx RENAME ATTRIBUTE x TO x;":                   review.DuplicateName,
		"ALTER TABLE t_a_idx RENAME COLUMN x TO x;":                     review.DuplicateName,
		"ALTER TYPE unsynced_enum ADD VALUE 'x' BEFORE 'x';":            review.UnsupportedForm,
		"ALTER TYPE unsynced_enum ADD VALUE 'x' AFTER 'x';":             review.UnsupportedForm,
		"ALTER TYPE unsynced_enum RENAME VALUE 'x' TO 'x';":             review.UnsupportedForm,
		// DROP TYPE array without IF EXISTS.
		"DROP TYPE unsynced_dom[];":     review.UnsupportedForm,
		"DROP TYPE unsynced_dom ARRAY;": review.UnsupportedForm,
		// FROM-less SELECT with unbound column, OR REPLACE VIEW column drop, EXECUTE.
		"CREATE VIEW v2 AS SELECT zzz;":               review.ObjectNotFound,
		"CREATE TABLE t2 AS SELECT zzz;":              review.ObjectNotFound,
		"CREATE TABLE IF NOT EXISTS t AS SELECT zzz;": review.ObjectNotFound,
		"CREATE MATERIALIZED VIEW mv2 AS SELECT zzz;": review.ObjectNotFound,
		"CREATE OR REPLACE VIEW v AS SELECT 1;":       review.UnsupportedForm,
		"CREATE TABLE t2 AS EXECUTE missing_plan;":    review.ObjectNotFound,
		// ALTER TABLE / SEQUENCE options, compression, relkinds, identity, EXCLUDE.
		"ALTER TABLE t ALTER COLUMN a SET COMPRESSION bogus;":   review.UnsupportedForm,
		"ALTER TABLE mv ALTER COLUMN a SET COMPRESSION bogus;":  review.UnsupportedForm,
		"ALTER TABLE mv ALTER COLUMN a SET STORAGE bogus;":      review.UnsupportedForm,
		"ALTER TABLE mv ALTER COLUMN 1 SET STATISTICS 10;":      review.UnsupportedForm,
		"ALTER TABLE t ALTER COLUMN a SET (bogus = 1);":         review.UnsupportedForm,
		"ALTER TABLE t ALTER COLUMN a RESTART RESTART;":         review.UnsupportedForm,
		"ALTER TABLE t ALTER COLUMN a RESTART 1.5;":             review.UnsupportedForm,
		"ALTER TABLE v ALTER COLUMN a RESTART;":                 review.UnsupportedForm,
		"ALTER TABLE v ALTER COLUMN a DROP IDENTITY;":           review.UnsupportedForm,
		"ALTER TABLE v SET (bogus = true);":                     review.UnsupportedForm,
		"ALTER TABLE t SET (fillfactor = 70, fillfactor = 80);": review.UnsupportedForm,
		"ALTER TABLE parted SET (fillfactor = 70);":             review.WrongObjectKind,
		"ALTER TABLE parted SET WITHOUT CLUSTER;":               review.UnsupportedForm,
		"ALTER TABLE parted NOT OF;":                            review.WrongObjectKind,
		"ALTER TABLE v SET WITHOUT OIDS;":                       review.WrongObjectKind,
		"ALTER TABLE v ENABLE RULE r;":                          review.WrongObjectKind,
		"ALTER TABLE t NO INHERIT t;":                           review.UnsupportedForm,
		"ALTER TABLE t ADD EXCLUDE (missing_col WITH =);":       review.ObjectNotFound,
		"ALTER TYPE pair DROP ATTRIBUTE x, DROP ATTRIBUTE x;":   review.DuplicateName,
		"ALTER SEQUENCE q OWNED BY col_only;":                   review.UnsupportedForm,
		"ALTER SEQUENCE q SEQUENCE NAME s2;":                    review.UnsupportedForm,
		"ALTER SEQUENCE q LOGGED;":                              review.UnsupportedForm,
		"ALTER SEQUENCE q UNLOGGED;":                            review.UnsupportedForm,
		"ALTER SEQUENCE q AS text;":                             review.UnsupportedForm,
	} {
		got := walkFindings(t, sql, target)
		if len(got) != 1 || !slices.Equal(got[0].Codes, []review.FindingCode{wantCode}) {
			t.Errorf("%q: findings = %+v, want one %s", sql, got, wantCode)
		}
	}

	// Control cases (0 false positives): valid routines, overloads, built-in
	// pg_catalog / system-excluded routines, search_path shadowing for domains,
	// valid reloptions/compression/sequence/trigger/view forms.
	for _, sql := range []string{
		// Synced routines & routine creation/overloads.
		"ALTER FUNCTION fn_existing STRICT;",
		"ALTER PROCEDURE proc_existing SET search_path = public;",
		"ALTER ROUTINE fn_existing STRICT;",
		"ALTER ROUTINE proc_existing SET search_path = public;",
		"DROP FUNCTION fn_existing;",
		"DROP FUNCTION IF EXISTS missing_fn;",
		"DROP PROCEDURE IF EXISTS missing_proc;",
		"DROP ROUTINE IF EXISTS missing_rtn;",
		"CREATE FUNCTION f_ov(int) RETURNS int LANGUAGE sql AS 'SELECT 1';\nCREATE FUNCTION f_ov(text) RETURNS int LANGUAGE sql AS 'SELECT 1';\nDROP FUNCTION f_ov(int);\nDROP FUNCTION f_ov(text);",
		"CREATE FUNCTION pg_temp.tf() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;\nCREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION tf();",
		// Built-in pg_catalog & Bytebase-excluded routine names (never falsely absent).
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION suppress_redundant_updates_trigger();",
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION pg_custom_trig();",
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXECUTE FUNCTION hypopg_trig();",
		"ALTER FUNCTION lower(text) IMMUTABLE;",
		"ALTER FUNCTION pg_catalog.abs(int) IMMUTABLE;",
		"DROP FUNCTION pg_my_ext_fn();",
		"DROP FUNCTION hypopg_reset();",
		// Trigger WHEN valid forms (OLD/NEW/subquery) and RI internal trigger.
		"CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW WHEN (OLD.a IS DISTINCT FROM NEW.a) EXECUTE FUNCTION trig_fn();",
		"CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW WHEN (OLD IS DISTINCT FROM NEW) EXECUTE FUNCTION trig_fn();",
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW WHEN (EXISTS (SELECT 1 FROM t WHERE a = 1)) EXECUTE FUNCTION trig_fn();",
		"ALTER TABLE t ENABLE TRIGGER tr_existing;",
		"ALTER TABLE t ENABLE TRIGGER ALL;",
		"ALTER TABLE t ENABLE TRIGGER USER;",
		"ALTER TABLE t ENABLE TRIGGER \"RI_ConstraintTrigger_c_1234\";",
		// Valid self-named types against pg_catalog built-in or another schema on search_path.
		"CREATE DOMAIN text AS text;",
		"CREATE TABLE text (a text);",
		"CREATE SCHEMA s_dom;\nSET search_path = s_dom, public;\nCREATE DOMAIN d AS d;",
		"CREATE TYPE r_drop AS RANGE (SUBTYPE = int4);\nDROP TYPE r_drop;",
		// Enum ADD VALUE IF NOT EXISTS with same neighbor (valid no-op when value exists).
		"ALTER TYPE mood ADD VALUE IF NOT EXISTS 'a' BEFORE 'a';",
		// DROP TYPE IF EXISTS array / SETOF.
		"DROP TYPE IF EXISTS unsynced_dom[];",
		"DROP TYPE SETOF mood;",
		// Valid view, CTAS, EXECUTE IF NOT EXISTS on existing table.
		"CREATE VIEW v_ok AS SELECT 1 AS x, (SELECT a FROM t LIMIT 1);",
		"CREATE VIEW v_ord AS SELECT 1 AS x ORDER BY x;",
		"CREATE OR REPLACE VIEW v AS SELECT 1 AS a, 2 AS b;",
		"CREATE OR REPLACE VIEW v AS SELECT 1 AS a, 2 AS b, 3 AS c;",
		"CREATE TABLE IF NOT EXISTS t AS EXECUTE existing_plan;",
		// Valid column/view/table/sequence options.
		"ALTER TABLE t ALTER COLUMN a SET COMPRESSION lz4;",
		"ALTER TABLE t ALTER COLUMN a SET COMPRESSION pglz;",
		"ALTER TABLE t ALTER COLUMN a SET COMPRESSION default;",
		"ALTER TABLE t ALTER COLUMN a SET (n_distinct = 100);",
		"ALTER TABLE v SET (security_barrier = true, security_invoker = false);",
		"ALTER TABLE v ALTER COLUMN a DROP IDENTITY IF EXISTS;",
		"ALTER TABLE t SET (fillfactor = 80);",
		"ALTER TABLE t ADD EXCLUDE (a WITH =);",
		"ALTER SEQUENCE q SET LOGGED;",
		"ALTER SEQUENCE q SET UNLOGGED;",
		"ALTER SEQUENCE q OWNED BY NONE;",
		"ALTER SEQUENCE q OWNED BY t.a;",
		"ALTER SEQUENCE q AS bigint;",
	} {
		if got := walkFindings(t, sql, target); len(got) != 0 {
			t.Errorf("%q: findings = %+v, want none", sql, got)
		}
	}
}
