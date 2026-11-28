package catalog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	omnianalysis "github.com/bytebase/omni/analysis"
	metadatapb "github.com/bytebase/omni/metadata"
)

func TestCockroachDBEngineMetadataMigration(t *testing.T) {
	ctx := context.Background()
	table := func(schema string) *metadatapb.SchemaMetadata {
		return &metadatapb.SchemaMetadata{Name: schema, Tables: []*metadatapb.TableMetadata{{
			Name:    "t",
			Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "INT8", Position: 1}},
		}}}
	}
	to := &omnianalysis.Database{Metadata: &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{table("public"), table("archive")}}}

	migration, err := Engine.MetadataMigration(ctx, omnianalysis.MetadataMigrationRequest{To: to})
	require.NoError(t, err)
	want, err := MetadataMigration(nil, to.Metadata, MetadataMigrationOptions{})
	require.NoError(t, err)
	require.Equal(t, want, migration)
	require.Contains(t, migration, "archive")

	migration, err = Engine.MetadataMigration(ctx, omnianalysis.MetadataMigrationRequest{To: to, SkipSchemas: []string{"archive"}})
	require.NoError(t, err)
	require.NotContains(t, migration, "archive")

	migration, err = Engine.MetadataMigration(ctx, omnianalysis.MetadataMigrationRequest{From: to, To: to})
	require.NoError(t, err)
	require.Empty(t, migration)
}
