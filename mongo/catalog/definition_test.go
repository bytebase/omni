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
						{Name: "_id_", Expressions: []string{`{"_id":1}`}, Definition: `{"v":{"$numberInt":"2"}}`},
						{Name: "name_1", Unique: true, Expressions: []string{`{"name":1}`}, Definition: `{"v":{"$numberInt":"2"}}`},
					},
				},
				{Name: "empty"},
				{
					Name:          "telemetry",
					CreateOptions: `{"timeseries":{"timeField":"at","metaField":"sensor","granularity":"seconds","bucketMaxSpanSeconds":{"$numberInt":"3600"}}}`,
				},
			},
			Views: []*metadata.ViewMetadata{{
				Name:       "v_adults",
				Definition: `{"viewOn":"users","pipeline":[{"$match":{"age":{"$gte":{"$numberInt":"18"}}}}]}`,
			}},
		}},
	}
	got, err := DatabaseDefinition(meta)
	require.NoError(t, err)
	require.Equal(t, `// MongoDB schema dump of database "app"
db.createCollection("users");
db.runCommand({"createIndexes": "users", "indexes": [{"key": {"_id":1}, "name": "_id_", "v": 2}]});
db.runCommand({"createIndexes": "users", "indexes": [{"key": {"name":1}, "name": "name_1", "unique": true, "v": 2}]});
db.createCollection("empty");
db.runCommand({"create": "telemetry", "timeseries": {"timeField": "at", "metaField": "sensor", "granularity": "seconds", "bucketMaxSpanSeconds": 3600}});
db.runCommand({"create": "v_adults", "viewOn": "users", "pipeline": [{"$match": {"age": {"$gte": 18}}}]});
`, got)
}

// TestDatabaseDefinitionNamesAViewWithoutADefinition keeps the one case a
// dump cannot restore: metadata synced before a view carried its definition.
func TestDatabaseDefinitionNamesAViewWithoutADefinition(t *testing.T) {
	got, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
		Name:    "app",
		Schemas: []*metadata.SchemaMetadata{{Views: []*metadata.ViewMetadata{{Name: "v_adults"}}}},
	})
	require.NoError(t, err)
	require.Equal(t, `// MongoDB schema dump of database "app"
// View "v_adults" is not restorable: the synced metadata carries no definition.
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

func TestDatabaseDefinitionRefusesAnOptionItCannotRender(t *testing.T) {
	_, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:          "c",
				CreateOptions: `{"validator":{"$jsonSchema":{"bsonType":"binData"}},"indexOptionDefaults":{"storageEngine":{"wiredTiger":{"configString":{"$binary":{"base64":"AQI=","subType":"00"}}}}}}`,
			}},
		}},
	})
	require.ErrorContains(t, err, "no mongosh spelling")
}
