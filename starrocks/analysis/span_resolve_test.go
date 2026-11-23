package analysis

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/starrocks/internal/testcontainer"
)

// spanFixture is a small stored schema for the resolver's unit tests.
var spanFixture = testcontainer.SpanFixture{Databases: []testcontainer.SpanDatabase{{
	Name: "shop",
	Tables: []testcontainer.SpanRelation{
		{Name: "customers", Columns: []string{"id", "name", "info"}},
		{Name: "orders", Columns: []string{"id", "customer_id", "amount"}},
	},
	Views: []testcontainer.SpanRelation{
		{Name: "broken", Columns: []string{"id"}, Definition: "SELECT id FROM orders WHERE"},
		{Name: "big", Columns: []string{"id"}, Definition: "SELECT id FROM shop.orders WHERE amount > 5"},
	},
	MaterializedViews: []testcontainer.SpanRelation{
		{Name: "totals", Definition: "CREATE MATERIALIZED VIEW `totals` (`cid`, `total`)\nREFRESH MANUAL\nAS SELECT customer_id, sum(amount) FROM shop.orders GROUP BY customer_id"},
	},
}}}

// A materialized view's columns follow its stored CREATE statement's query
// and column list, and reading it reads the tables behind it.
func TestResolveQuerySpanMaterializedView(t *testing.T) {
	span, err := resolveWithCatalog(t, "SELECT cid, total FROM totals")
	if err != nil {
		t.Fatal(err)
	}
	if span.NotFoundError != nil {
		t.Fatalf("NotFoundError %v", span.NotFoundError)
	}
	if len(span.Results) != 2 || !span.Results[1].SourceColumns[ColumnResource{Database: "shop", Table: "orders", Column: "amount"}] {
		t.Errorf("results %+v, want total to read orders.amount", span.Results)
	}
	if !span.SourceColumns[ColumnResource{Database: "shop", Table: "orders"}] || !span.SourceColumns[ColumnResource{Database: "shop", Table: "totals"}] {
		t.Errorf("tables %v, want totals and orders", span.SourceColumns)
	}
}

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
		{"SELECT * FROM orders o RIGHT SEMI JOIN customers c ON o.id = c.id", []string{"id", "name", "info"}},
		{"SELECT * FROM orders o RIGHT ANTI JOIN customers c ON o.id = c.id", []string{"id", "name", "info"}},
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

// A star names a table of another catalog by its catalog too; the
// metadata cannot expand it.
func TestResolveQuerySpanExternalStar(t *testing.T) {
	for _, statement := range []string{
		"SELECT * FROM ext.shop.orders",
		"SELECT ext.shop.orders.* FROM ext.shop.orders",
		"SELECT o.* FROM ext.shop.orders o",
	} {
		span, err := resolveWithCatalog(t, statement)
		if err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
		if span.NotFoundError != nil {
			t.Fatalf("%s: NotFoundError %v", statement, span.NotFoundError)
		}
		if len(span.Results) != 1 || !span.Results[0].SelectAsterisk {
			t.Errorf("%s: results %+v, want one unexpanded star", statement, span.Results)
		}
	}
}

// SelectAsterisk marks a star, never a column a star names.
func TestResolveQuerySpanSelectAsterisk(t *testing.T) {
	for _, statement := range []string{"SELECT '*'", "SELECT 1 AS '*'", "SELECT id AS 'o.*' FROM orders"} {
		for _, opts := range []ResolveOptions{
			{DefaultDatabase: "shop", Catalog: spanCatalog(spanFixture)},
			{DefaultDatabase: "shop"},
		} {
			span, err := ResolveQuerySpan(statement, opts)
			if err != nil {
				t.Fatalf("%s: %v", statement, err)
			}
			if len(span.Results) != 1 || span.Results[0].SelectAsterisk {
				t.Errorf("%s (catalog %t): results %+v, want one column", statement, opts.Catalog.GetDatabase != nil, span.Results)
			}
		}
	}
}

// A star the metadata cannot expand stays one through derived tables,
// CTEs and set operations.
func TestResolveQuerySpanStarThroughQueries(t *testing.T) {
	for _, statement := range []string{
		"SELECT * FROM (SELECT * FROM ext.shop.orders) d",
		"WITH c AS (SELECT * FROM ext.shop.orders) SELECT * FROM c",
		"SELECT * FROM ext.shop.orders UNION ALL SELECT * FROM ext.shop.orders",
		"SELECT * FROM ext.shop.orders UNION ALL SELECT 1",
	} {
		span, err := resolveWithCatalog(t, statement)
		if err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
		if span.NotFoundError != nil {
			t.Fatalf("%s: NotFoundError %v", statement, span.NotFoundError)
		}
		if len(span.Results) != 1 || !span.Results[0].SelectAsterisk {
			t.Errorf("%s: results %+v, want one unexpanded star", statement, span.Results)
		}
	}
	span, err := resolveWithCatalog(t, "SELECT d.anything FROM (SELECT * FROM ext.shop.orders) d")
	if err != nil || span.NotFoundError != nil {
		t.Fatalf("a column of a derived table over an unexpanded star: %v %v", err, span.NotFoundError)
	}
}

// A struct field path reads its column: s.f and t.s.f name column s.
func TestResolveQuerySpanStructField(t *testing.T) {
	info := ColumnResource{Database: "shop", Table: "customers", Column: "info"}
	for _, statement := range []string{
		"SELECT info.city FROM customers",
		"SELECT c.info.city FROM customers c",
		"SELECT customers.info.city.zip FROM customers",
	} {
		span, err := resolveWithCatalog(t, statement)
		if err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
		if span.NotFoundError != nil {
			t.Fatalf("%s: NotFoundError %v", statement, span.NotFoundError)
		}
		if len(span.Results) != 1 || len(span.Results[0].SourceColumns) != 1 || !span.Results[0].SourceColumns[info] {
			t.Errorf("%s: results %+v, want customers.info", statement, span.Results)
		}
	}
	span, err := resolveWithCatalog(t, "SELECT nope.city FROM customers")
	if err != nil || span.NotFoundError == nil {
		t.Errorf("a path on no column: %v, NotFoundError %v; want one", err, span.NotFoundError)
	}
}

// A table function's columns read its arguments, through a star and the
// derived table over it.
func TestResolveQuerySpanTableFunctionThroughStar(t *testing.T) {
	amount := ColumnResource{Database: "shop", Table: "orders", Column: "amount"}
	for _, statement := range []string{
		"SELECT u.unnest FROM orders o, unnest(o.amount) u",
		"SELECT d.unnest FROM (SELECT u.* FROM orders o, unnest(o.amount) u) d",
		"SELECT u.* FROM orders o, unnest(o.amount) u",
	} {
		span, err := resolveWithCatalog(t, statement)
		if err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
		if span.NotFoundError != nil {
			t.Fatalf("%s: NotFoundError %v", statement, span.NotFoundError)
		}
		if len(span.Results) != 1 || !span.Results[0].SourceColumns[amount] {
			t.Errorf("%s: results %+v, want orders.amount", statement, span.Results)
		}
	}
}

// A field path on a lambda parameter, x -> x.f, is no column: the array the
// lambda maps is what the result reads.
func TestResolveQuerySpanLambdaFieldPath(t *testing.T) {
	span, err := resolveWithCatalog(t, "SELECT array_map(x -> x.city, info) AS cities FROM customers")
	if err != nil {
		t.Fatal(err)
	}
	if span.NotFoundError != nil {
		t.Fatalf("NotFoundError %v", span.NotFoundError)
	}
	info := ColumnResource{Database: "shop", Table: "customers", Column: "info"}
	if len(span.Results) != 1 || len(span.Results[0].SourceColumns) != 1 || !span.Results[0].SourceColumns[info] {
		t.Errorf("results %+v, want customers.info", span.Results)
	}
}

// An EXPLAIN reads what its query reads, views expanded, and answers with
// a plan rather than the query's columns.
func TestResolveQuerySpanExplain(t *testing.T) {
	span, err := resolveWithCatalog(t, "EXPLAIN SELECT * FROM big")
	if err != nil {
		t.Fatal(err)
	}
	if !span.SourceColumns[ColumnResource{Database: "shop", Table: "orders"}] || !span.SourceColumns[ColumnResource{Database: "shop", Table: "big"}] {
		t.Errorf("tables %v, want big and orders", span.SourceColumns)
	}
	if len(span.Results) != 0 {
		t.Errorf("results %+v, want none", span.Results)
	}
	span, err = resolveWithCatalog(t, "EXPLAIN SELECT nope FROM orders")
	if err != nil || span.NotFoundError == nil {
		t.Errorf("EXPLAIN of a missing column: %v, NotFoundError %v; want one", err, span.NotFoundError)
	}
}
