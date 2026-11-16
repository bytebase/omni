package catalog

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

func TestDatabaseDefinition(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name: "users",
					Indexes: []*metadata.IndexMetadata{
						{Name: "_id_", Expressions: []string{`{"_id":1}`}},
						{Name: "name_1", Unique: true, Expressions: []string{`{"name":1}`}},
					},
				},
				{Name: "empty"},
			},
			Views: []*metadata.ViewMetadata{{Name: "v_adults"}},
		}},
	}
	got, err := DatabaseDefinition(meta)
	require.NoError(t, err)
	require.Equal(t, `// MongoDB schema dump of database "app"
db.createCollection("users");
db.getCollection("users").createIndex({"_id":1}, {"name": "_id_"});
db.getCollection("users").createIndex({"name":1}, {"name": "name_1", "unique": true});
db.createCollection("empty");
// View "v_adults" is not restorable: the synced metadata carries no pipeline.
`, got)
}

func TestDatabaseDefinitionRefusesAnIndexItCannotRender(t *testing.T) {
	for _, tc := range []struct {
		name   string
		index  *metadata.IndexMetadata
		errMsg string
	}{
		{
			name:   "no expression",
			index:  &metadata.IndexMetadata{Name: "idx"},
			errMsg: "0 expressions",
		},
		{
			name:   "two expressions",
			index:  &metadata.IndexMetadata{Name: "idx", Expressions: []string{`{"a":1}`, `{"b":1}`}},
			errMsg: "2 expressions",
		},
		{
			name:   "not a document",
			index:  &metadata.IndexMetadata{Name: "idx", Expressions: []string{`not json`}},
			errMsg: "not a document",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
				Schemas: []*metadata.SchemaMetadata{{
					Tables: []*metadata.TableMetadata{{Name: "c", Indexes: []*metadata.IndexMetadata{tc.index}}},
				}},
			})
			require.ErrorContains(t, err, tc.errMsg)
		})
	}
}
