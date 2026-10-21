package review

import (
	"context"
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
