package driver

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDumpWritesNothing pins the no-op: a caller that dumps a synced table and
// restores the result restores nothing, and Dump reports neither a write nor an
// error.
func TestDumpWritesNothing(t *testing.T) {
	db := &metadatapb.DatabaseSchemaMetadata{
		Name: "123456789012-us-east-1",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Tables: []*metadatapb.TableMetadata{
					{
						Name:     "orders",
						RowCount: 2,
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "pk"},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{Name: "pk", Expressions: []string{"pk"}, Type: "HASH"},
						},
					},
				},
			},
		},
	}
	before := proto.Clone(db)

	var w strings.Builder
	require.NoError(t, (&Driver{}).Dump(context.Background(), &w, db))
	require.Zero(t, w.Len(), "Dump wrote %q", w.String())
	if !proto.Equal(before, db) {
		t.Errorf("Dump mutated its input metadata")
	}

	require.NoError(t, (&Driver{}).Dump(context.Background(), &w, nil), "Dump of nil metadata")
	require.Zero(t, w.Len())
}
