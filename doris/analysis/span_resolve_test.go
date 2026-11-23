package analysis

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/doris/internal/testcontainer"
)

// spanFixture is a small stored schema for the resolver's unit tests.
var spanFixture = testcontainer.SpanFixture{Databases: []testcontainer.SpanDatabase{{
	Name: "shop",
	Tables: []testcontainer.SpanRelation{
		{Name: "customers", Columns: []string{"id", "name"}},
		{Name: "orders", Columns: []string{"id", "customer_id", "amount"}},
	},
	Views: []testcontainer.SpanRelation{
		{Name: "broken", Columns: []string{"id"}, Definition: "SELECT id FROM orders WHERE"},
	},
}}}

func resolveWithCatalog(t *testing.T, statement string) (*ResolvedQuerySpan, error) {
	t.Helper()
	return ResolveQuerySpan(statement, ResolveOptions{DefaultDatabase: "shop", Catalog: spanCatalog(spanFixture)})
}

// A view whose stored definition does not parse fails the span: resolving
// it without its definition would leave out the tables it reads.
func TestResolveQuerySpanViewDefinitionFails(t *testing.T) {
	if span, err := resolveWithCatalog(t, "SELECT * FROM broken"); err == nil {
		t.Fatalf("ResolveQuerySpan = %+v, want an error", span)
	}
}

// A semi or anti join returns one side's columns, and reads both sides.
func TestResolveQuerySpanSemiJoin(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      []string
	}{
		{"SELECT * FROM orders o LEFT SEMI JOIN customers c ON o.id = c.id", []string{"id", "customer_id", "amount"}},
		{"SELECT * FROM orders o LEFT ANTI JOIN customers c ON o.id = c.id", []string{"id", "customer_id", "amount"}},
		{"SELECT * FROM orders o RIGHT SEMI JOIN customers c ON o.id = c.id", []string{"id", "name"}},
		{"SELECT * FROM orders o RIGHT ANTI JOIN customers c ON o.id = c.id", []string{"id", "name"}},
	} {
		span, err := resolveWithCatalog(t, tc.statement)
		if err != nil {
			t.Fatalf("%s: %v", tc.statement, err)
		}
		var names []string
		for _, r := range span.Results {
			names = append(names, r.Name)
		}
		if !slices.Equal(names, tc.want) {
			t.Errorf("%s: results %v, want %v", tc.statement, names, tc.want)
		}
		if len(span.SourceColumns) != 2 {
			t.Errorf("%s: tables %v, want both", tc.statement, span.SourceColumns)
		}
	}
}

// A leading VALUES arm's WITH scopes over every arm of the set operation.
func TestResolveQuerySpanValuesWithScopesSetOp(t *testing.T) {
	span, err := resolveWithCatalog(t, "WITH c AS (SELECT id FROM orders) VALUES (0) UNION ALL SELECT id FROM c")
	if err != nil {
		t.Fatal(err)
	}
	if span.NotFoundError != nil {
		t.Fatalf("NotFoundError %v", span.NotFoundError)
	}
	if len(span.Results) != 1 || !span.Results[0].SourceColumns[ColumnResource{Database: "shop", Table: "orders", Column: "id"}] {
		t.Errorf("results %+v, want orders.id", span.Results)
	}
}
