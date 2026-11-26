package analysis

import (
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
)

// TestResolveQuerySpanWildcardTable holds a wildcard table's lineage to the
// tables it matches, which is what Bytebase masks.
func TestResolveQuerySpanWildcardTable(t *testing.T) {
	cols := func(names ...string) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, n := range names {
			out = append(out, &metadata.ColumnMetadata{Name: n})
		}
		return out
	}
	meta := &metadata.DatabaseSchemaMetadata{Name: "ds", Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{
			{Name: "events_2024", Columns: cols("EventId", "Kind")},
			{Name: "events_2025", Columns: cols("EventId", "Kind")},
			{Name: "eventsx", Columns: cols("EventId", "Kind")},
			{Name: "logs_a", Columns: cols("Id", "Msg")},
			{Name: "logs_b", Columns: cols("Id", "Msg", "Extra")},
		},
		Views: []*metadata.ViewMetadata{{Name: "events_view", Columns: cols("EventId", "Kind")}},
	}}}
	opts := bigqueryResolveOptions()
	opts.DefaultDatabase = "ds"
	opts.GetDatabaseMetadata = corpusGetter(meta)

	both := func(column string) string {
		return "ds.events_2024." + column + ",ds.events_2025." + column
	}
	for _, tc := range []struct {
		query string
		want  []string // per result: sorted sources joined by ","
	}{
		{"SELECT EventId FROM `events_*`", []string{both("EventId")}},
		{"SELECT * FROM `ds.events_*`", []string{both("EventId"), both("Kind")}},
		{"SELECT e.Kind FROM `events_20*` AS e WHERE e._TABLE_SUFFIX > '2024'", []string{both("Kind")}},
		{"SELECT _TABLE_SUFFIX AS y, Kind FROM `events_*`", []string{"", both("Kind")}},
		{"FROM `events_*` |> SELECT Kind", []string{both("Kind")}},
	} {
		span, err := ResolveQuerySpan(tc.query, opts)
		if err != nil {
			t.Fatalf("%s: %v", tc.query, err)
		}
		if span.NotFoundError != nil {
			t.Fatalf("%s: %v", tc.query, span.NotFoundError)
		}
		var got []string
		for _, r := range span.Results {
			var sources []string
			for c := range r.SourceColumns {
				sources = append(sources, c.Database+"."+c.Table+"."+c.Column)
			}
			slices.Sort(sources)
			got = append(got, strings.Join(sources, ","))
		}
		if !slices.Equal(got, tc.want) {
			t.Errorf("%s:\n got %q\nwant %q", tc.query, got, tc.want)
		}
	}

	// Tables whose columns differ leave the wildcard's columns unknown, since
	// BigQuery takes them from the newest table, so the read is not found.
	span, err := ResolveQuerySpan("SELECT Id FROM `logs_*`", opts)
	if err != nil {
		t.Fatal(err)
	}
	if span.NotFoundError == nil {
		t.Errorf("SELECT Id FROM `logs_*`: want a NotFoundError, got %v", span.Results)
	}
}
