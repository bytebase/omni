package analysis

import (
	metadatapb "github.com/bytebase/omni/metadata"
)

// buildMockDatabaseMetadataGetter mirrors the ResolveOptions contract: a
// database that is not synced returns nil metadata and no error.
func buildMockDatabaseMetadataGetter(databaseMetadata []*metadatapb.DatabaseSchemaMetadata) (func(string) (*metadatapb.DatabaseSchemaMetadata, error), func() ([]string, error)) {
	return func(databaseName string) (*metadatapb.DatabaseSchemaMetadata, error) {
			for _, meta := range databaseMetadata {
				if meta.GetName() == databaseName {
					return meta, nil
				}
			}
			return nil, nil
		}, func() ([]string, error) {
			var names []string
			for _, meta := range databaseMetadata {
				names = append(names, meta.GetName())
			}
			return names, nil
		}
}
