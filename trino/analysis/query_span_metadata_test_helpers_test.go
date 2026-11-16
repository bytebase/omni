package analysis

import (
	"fmt"

	metadatapb "github.com/bytebase/omni/metadata"
)

func buildMockDatabaseMetadataGetter(databaseMetadata []*metadatapb.DatabaseSchemaMetadata) (func(string) (*metadatapb.DatabaseSchemaMetadata, error), func() ([]string, error)) {
	return func(databaseName string) (*metadatapb.DatabaseSchemaMetadata, error) {
			for _, meta := range databaseMetadata {
				if meta.GetName() == databaseName {
					return meta, nil
				}
			}
			return nil, fmt.Errorf("database %q not found", databaseName)
		}, func() ([]string, error) {
			var names []string
			for _, meta := range databaseMetadata {
				names = append(names, meta.GetName())
			}
			return names, nil
		}
}
