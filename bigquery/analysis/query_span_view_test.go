package analysis

import (
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

// viewSchemas is a dataset ds1 as bigquery/driver's sync stores it: two tables,
// a view over both, a view over that view with a derived column, and a view
// whose definition reads a dataset the metadata does not hold.
func viewSchemas() []*metadata.SchemaMetadata {
	columns := func(names ...string) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, n := range names {
			out = append(out, &metadata.ColumnMetadata{Name: n})
		}
		return out
	}
	return []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{
			{Name: "customers", Columns: columns("id", "email")},
			{Name: "orders", Columns: columns("customer_id", "amount")},
		},
		Views: []*metadata.ViewMetadata{
			{
				Name:       "customer_orders",
				Definition: "SELECT c.id, c.email, o.amount FROM ds1.customers AS c JOIN ds1.orders AS o ON c.id = o.customer_id",
				Columns:    columns("id", "email", "amount"),
			},
			{
				Name:       "big_orders",
				Definition: "SELECT co.id, UPPER(co.email) AS shout FROM ds1.customer_orders AS co WHERE co.amount > 100",
				Columns:    columns("id", "shout"),
			},
			{
				Name:       "elsewhere",
				Definition: "SELECT x FROM other.t",
				Columns:    columns("x"),
			},
		},
	}}
}

func viewTableNames(set ColumnSet) []string {
	var out []string
	for c := range set {
		out = append(out, c.Database+"."+c.Table)
	}
	slices.Sort(out)
	return out
}

// TestViewListsItsTables: a view read lists the tables and views its
// definition reads, through nested views, as BigQuery's dry run reports them.
func TestViewListsItsTables(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT * FROM big_orders", "ds1", viewSchemas())
	require.Equal(t, []string{"ds1.big_orders", "ds1.customer_orders", "ds1.customers", "ds1.orders"}, viewTableNames(span.SourceColumns))
	require.Nil(t, span.NotFoundError)
}

// TestViewLineageReachesBaseColumns: a column read through a view names the
// base table column it comes from, so a semantic type on the table masks it.
func TestViewLineageReachesBaseColumns(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT id, shout FROM big_orders", "ds1", viewSchemas())
	require.Len(t, span.Results, 2)
	require.Equal(t, []string{"big_orders.id", "customer_orders.id", "customers.id"}, sourcesOf(span.Results[0]))
	require.Equal(t, []string{"big_orders.shout", "customer_orders.email", "customers.email"}, sourcesOf(span.Results[1]))
}

// TestViewStarIsNoPlainField: a star over a base table is a plain field, but a
// star over a view is not, as a column resolved through a derived relation is
// not: the view's query may transform the column (UPPER(email) here), so a
// partial mask of the base column would not fit it.
func TestViewStarIsNoPlainField(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT * FROM customers", "ds1", viewSchemas())
	require.True(t, span.Results[0].IsPlainField)
	span = resolvePin(t, bigqueryResolveOptions(), "SELECT * FROM big_orders", "ds1", viewSchemas())
	require.Len(t, span.Results, 2)
	require.False(t, span.Results[0].IsPlainField)
	require.False(t, span.Results[1].IsPlainField)
	require.Equal(t, []string{"big_orders.shout", "customer_orders.email", "customers.email"}, sourcesOf(span.Results[1]))
}

// TestViewReadingAnUnheldDataset: a view whose definition reads a dataset the
// metadata does not hold lists the table it reads, as BigQuery does, and the
// query is not refused for a dataset it does not name.
func TestViewReadingAnUnheldDataset(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT x FROM elsewhere", "ds1", viewSchemas())
	require.Equal(t, []string{"ds1.elsewhere", "other.t"}, viewTableNames(span.SourceColumns))
	require.Nil(t, span.NotFoundError)
	require.Len(t, span.Results, 1)
	require.Equal(t, []string{"elsewhere.x"}, sourcesOf(span.Results[0]), "other.t holds no columns the metadata lists")
}

// TestViewMetadataFailureIsFatal: a metadata lookup that fails while a view's
// definition is resolved fails the span, as a failed lookup of a table the
// query names does: lineage silently cut at the view would leave a masked base
// column unmasked.
func TestViewMetadataFailureIsFatal(t *testing.T) {
	dbMetadata := &metadata.DatabaseSchemaMetadata{Name: "ds1", Schemas: []*metadata.SchemaMetadata{{
		Views: []*metadata.ViewMetadata{{
			Name:       "v",
			Definition: "SELECT x FROM ds2.t",
			Columns:    []*metadata.ColumnMetadata{{Name: "x"}},
		}},
	}}}
	outage := errors.New("store outage")
	opts := bigqueryResolveOptions()
	opts.DefaultDatabase = "ds1"
	opts.GetDatabaseMetadata = func(name string) (*metadata.DatabaseSchemaMetadata, error) {
		if name == "ds1" {
			return dbMetadata, nil
		}
		return nil, outage
	}
	_, err := ResolveQuerySpan("SELECT x FROM v", opts)
	require.ErrorIs(t, err, outage)
}

// TestSelfReadingViewStops: a definition that reaches its own view (which
// BigQuery refuses to store, but metadata can hold) stops there.
func TestSelfReadingViewStops(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT a FROM v", "ds1", []*metadata.SchemaMetadata{{
		Views: []*metadata.ViewMetadata{{
			Name:       "v",
			Definition: "SELECT a FROM ds1.v",
			Columns:    []*metadata.ColumnMetadata{{Name: "a"}},
		}},
	}})
	require.Equal(t, []string{"ds1.v"}, viewTableNames(span.SourceColumns))
}
