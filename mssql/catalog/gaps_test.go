package catalog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDefinitionGaps pins the objects a definition cannot restore yet. Each
// subtest fails when its gap closes, so the case moves into the fidelity tests
// instead of staying here forever.
func TestDefinitionGaps(t *testing.T) {
	t.Run("indexed view", func(t *testing.T) {
		// The catalog rejects an index on a view because the catalog holds
		// metadata.DatabaseSchemaMetadata, and ViewMetadata has no field for
		// indexes: the driver's sync cannot report one and a definition
		// cannot render one. Closing this needs `repeated IndexMetadata
		// indexes` on ViewMetadata in proto/omni/metadata/database.proto,
		// which Bytebase stores, so it waits for that proto change.
		_, err := LoadSQL(`
CREATE TABLE [dbo].[t] ([id] int NOT NULL);
GO
CREATE VIEW [dbo].[v] WITH SCHEMABINDING AS SELECT [id] FROM [dbo].[t];
GO
CREATE UNIQUE CLUSTERED INDEX [UX_v_id] ON [dbo].[v] ([id]);
`)
		require.Error(t, err)
		require.Contains(t, err.Error(), "indexes on views are not supported")
	})
}
