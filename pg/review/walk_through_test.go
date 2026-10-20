package review

import (
	"context"
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

	// The same schema with and without other specs is two inputs.
	got = walkFindings(t, "CREATE TABLE t (a int);", alone, other, alone)
	if len(got) != 1 || len(got[0].Targets) != 2 || got[0].Targets[0] != 0 || got[0].Targets[1] != 2 {
		t.Errorf("findings = %+v, want one on targets 0 and 2", got)
	}
}
