package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

// These tests pin dialect behaviors the legacy-recorded corpora do not
// discriminate, so a change to either flag would otherwise pass unnoticed.

// TestUppercaseStarMergeName: a set-operation merge whose LEFT arm is a base
// star takes the output name from the expanded star column. BigQuery renders it
// upper-cased (the legacy expression-name rendering); Spanner passes the
// column's metadata case through.
func TestUppercaseStarMergeName(t *testing.T) {
	tables := pinTables(
		&metadata.TableMetadata{Name: "rt", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "secret"}}},
		&metadata.TableMetadata{Name: "lt", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "label"}}},
	)
	statement := "SELECT * FROM rt UNION ALL SELECT id, label FROM lt"

	bigquery := resolvePin(t, bigqueryResolveOptions(), statement, "ds1", tables)
	require.Equal(t, []string{"ID", "SECRET"}, []string{bigquery.Results[0].Name, bigquery.Results[1].Name},
		"BigQuery upper-cases a left-star merge name")

	spannerOpts := spannerResolveOptions()
	spanner := resolvePin(t, spannerOpts, statement, "ds1", tables)
	require.Equal(t, []string{"id", "secret"}, []string{spanner.Results[0].Name, spanner.Results[1].Name},
		"Spanner keeps the star column's metadata case")
	require.Equal(t, []string{"lt.id", "rt.id"}, sourcesOf(spanner.Results[0]))
}

// TestIgnoreWrittenDatabase: under the Spanner metadata model a 3-part
// db.schema.table reference records the session database and the written schema;
// BigQuery resolves the written dataset as the database and uses its single
// unnamed schema.
func TestIgnoreWrittenDatabase(t *testing.T) {
	spannerSpan := resolvePin(t, spannerResolveOptions(), "SELECT id FROM other.analytics.events", "db", []*metadata.SchemaMetadata{
		{Name: ""},
		{Name: "analytics", Tables: []*metadata.TableMetadata{
			{Name: "events", Columns: []*metadata.ColumnMetadata{{Name: "id"}}},
		}},
	})
	require.Equal(t, ColumnSet{
		{Database: "db", Schema: "analytics", Table: "events"}: true,
	}, spannerSpan.SourceColumns, "Spanner records the session database, not the written one")
	require.Equal(t, []string{"analytics.events.id"}, sourcesOf(spannerSpan.Results[0]))

	bigquerySpan := resolvePin(t, bigqueryResolveOptions(), "SELECT id FROM ds1.events", "ds1", pinTables(
		&metadata.TableMetadata{Name: "events", Columns: []*metadata.ColumnMetadata{{Name: "id"}}},
	))
	require.Equal(t, ColumnSet{
		{Database: "ds1", Table: "events"}: true,
	}, bigquerySpan.SourceColumns, "BigQuery resolves the dataset as the database")
	require.Equal(t, []string{"events.id"}, sourcesOf(bigquerySpan.Results[0]), "the dataset is the database, not a schema")
}

// TestAliasedColumnLineage: a column reached through a relation alias carries the
// physical table's lineage (aliases are resolved before lineage matching).
func TestAliasedColumnLineage(t *testing.T) {
	span := resolvePin(t, bigqueryResolveOptions(), "SELECT g.x AS y FROM galleries g", "ds1", pinTables(
		&metadata.TableMetadata{Name: "galleries", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "x"}}},
	))
	require.Len(t, span.Results, 1)
	require.Equal(t, "Y", span.Results[0].Name)
	require.Equal(t, []string{"galleries.x"}, sourcesOf(span.Results[0]))
}
