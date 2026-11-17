package driver

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDumpWritesNothing pins the no-op: a caller that dumps a synced container
// and restores the result restores nothing, and Dump reports neither a write nor
// an error.
func TestDumpWritesNothing(t *testing.T) {
	db := &metadatapb.DatabaseSchemaMetadata{
		Name: "omni_driver",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Tables: []*metadatapb.TableMetadata{
					{Name: "items"},
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
