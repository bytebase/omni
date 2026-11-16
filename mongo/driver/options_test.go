package driver

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCollectionOptionsKeepTheServersDocument(t *testing.T) {
	got, err := collectionOptions(bson.M{
		"name": "c",
		"options": bson.D{
			{Key: "timeseries", Value: bson.D{{Key: "timeField", Value: "at"}, {Key: "metaField", Value: "sensor"}}},
		},
	})
	require.NoError(t, err)
	require.Equal(t, `{"timeseries":{"timeField":"at","metaField":"sensor"}}`, got)
}

func TestCollectionOptionsAreEmptyWhenThereAreNone(t *testing.T) {
	got, err := collectionOptions(bson.M{"name": "c", "options": bson.D{}})
	require.NoError(t, err)
	require.Empty(t, got)
}

// TestIndexDefinitionLeavesTheKeyNameAndUniquenessToTheirFields holds the
// split: Expressions carries the key, Name and Unique have their own fields,
// and Definition carries the other options listIndexes reports.
func TestIndexDefinitionLeavesTheKeyNameAndUniquenessToTheirFields(t *testing.T) {
	got, err := indexDefinition(bson.M{
		"key":                bson.D{{Key: "a", Value: int32(1)}},
		"name":               "a_1",
		"unique":             true,
		"v":                  int32(2),
		"expireAfterSeconds": int32(3600),
	})
	require.NoError(t, err)
	require.Equal(t, `{"expireAfterSeconds":{"$numberInt":"3600"},"v":{"$numberInt":"2"}}`, got)
}
