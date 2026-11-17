package driver

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDumpWritesNothing pins the no-op: a caller that dumps a synced keyspace
// and restores the result restores nothing, and Dump reports neither a write
// nor an error.
func TestDumpWritesNothing(t *testing.T) {
	db := &metadatapb.DatabaseSchemaMetadata{
		Name: "omni_driver",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "omni_events",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "group_id", Type: "int"},
							{Name: "event_id", Type: "int"},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{
								Name:        "PRIMARY KEY",
								Expressions: []string{"group_id", "event_id"},
								Primary:     true,
								Definition:  "(group_id,event_id)",
							},
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
