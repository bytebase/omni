package driver

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

// TestExecute writes the corpus and reads it back: Execute answers 0 affected
// rows whatever the requests did, and the sync sees the indices the corpus
// created.
func TestExecute(t *testing.T) {
	d := openDriver(t)
	affected, err := d.Execute(context.Background(), `
PUT /exec_items
{"settings": {"number_of_shards": 1, "number_of_replicas": 0}}

POST /exec_items/_doc/1
{"note": "a"}

POST /exec_items/_refresh`, driver.ExecuteOptions{
		Transaction: driver.Transactional,
		Isolation:   driver.Serializable,
	})
	require.NoError(t, err)
	require.Zero(t, affected)

	schema, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	require.Len(t, schema.Schemas[0].Tables, 1)
	require.Equal(t, "exec_items", schema.Schemas[0].Tables[0].Name)
	require.Equal(t, int64(1), schema.Schemas[0].Tables[0].RowCount)
}

// TestExecuteStopsAtAFailingRequest holds a script's failure to its own error,
// and the requests before it stay: each REST request stands alone.
func TestExecuteStopsAtAFailingRequest(t *testing.T) {
	d := openDriver(t)
	_, err := d.Execute(context.Background(), `
PUT /kept
{"settings": {"number_of_shards": 1, "number_of_replicas": 0}}

GET /missing-index/_search
{"query": {"match_all": {}}}`, driver.ExecuteOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "request failed with status code 404")

	schema, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	require.Len(t, schema.Schemas[0].Tables, 1)
	require.Equal(t, "kept", schema.Schemas[0].Tables[0].Name)
}

// TestExecuteCorpus documents that the corpus every test records is itself an
// Execute call: one request per Kibana console block.
func TestExecuteCorpus(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))

	schema, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	require.Len(t, schema.Schemas[0].Tables, 2)
}
