package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

// These tests pin dialect behaviors the legacy-recorded corpora do not
// discriminate, so a change to either flag would otherwise pass unnoticed.

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
