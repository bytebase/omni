package testcontainer

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTheTapeHoldsWhatTheDriverReads holds each answer's narrowing to the
// fields the driver reads: the values a container makes up (a cluster UUID, a
// node id, a timing) never reach a tape, so two recordings answer the same
// bytes.
func TestTheTapeHoldsWhatTheDriverReads(t *testing.T) {
	t.Run("root document", func(t *testing.T) {
		body := `{"name":"node-1","cluster_name":"docker-cluster","cluster_uuid":"abc","version":{"number":"8.19.25","build_hash":"def"}}`
		narrowed, ok := narrow("/", "application/json", []byte(body))
		require.True(t, ok)
		require.Equal(t, `{"version":{"number":"8.19.25"}}`, string(narrowed))
	})

	t.Run("index settings", func(t *testing.T) {
		body := `{"alpha":{"settings":{"index":{"hidden":"false","number_of_shards":"1"}}},"hidden-one":{"settings":{"index":{"hidden":"true"}}}}`
		narrowed, ok := narrow("/*/_settings?filter_path=**.hidden&expand_wildcards=open,hidden", "application/json", []byte(body))
		require.True(t, ok)
		require.Equal(t, `{"hidden-one":{"settings":{"index":{"hidden":"true"}}}}`, string(narrowed))
	})

	t.Run("cat indices", func(t *testing.T) {
		body := `[{"health":"green","status":"open","index":"beta","uuid":"one","pri":"1","rep":"0","docs.count":"5","docs.deleted":"0","store.size":"5.1kb","pri.store.size":"5.1kb","dataset.size":"5.1kb"},` +
			`{"health":"green","status":"open","index":"alpha","uuid":"two","pri":"1","rep":"0","docs.count":"3","docs.deleted":"0","store.size":"4.9kb","pri.store.size":"4.9kb","dataset.size":"4.9kb"}]`
		narrowed, ok := narrow("/_cat/indices?format=json&pretty", "application/json", []byte(body))
		require.True(t, ok)
		require.Equal(t, `[{"index":"alpha","docs.count":"3","store.size":"4.9kb"},{"index":"beta","docs.count":"5","store.size":"5.1kb"}]`, string(narrowed))
	})

	t.Run("search", func(t *testing.T) {
		body := `{"took":7,"timed_out":false,"_shards":{"total":1},"hits":{"total":{"value":3}},"profile":{"shards":[{"id":"[node][alpha][0]","node_id":"node","searches":[{"rewrite_time":1234,"query":[{"time_in_nanos":99,"breakdown":{"build_scorer":5,"build_scorer_count":2}}],"collector":[{"time_in_nanos":3}]}]}]}}`
		narrowed, ok := narrow("/alpha/_search", "application/json", []byte(body))
		require.True(t, ok)
		var got map[string]any
		decoder := json.NewDecoder(bytes.NewReader(narrowed))
		decoder.UseNumber()
		require.NoError(t, decoder.Decode(&got))
		require.Equal(t, json.Number("0"), got["took"])
		profile := got["profile"].(map[string]any)["shards"].([]any)[0].(map[string]any)
		require.Equal(t, "0", profile["id"])
		require.Equal(t, "0", profile["node_id"])
		search := profile["searches"].([]any)[0].(map[string]any)
		require.Equal(t, json.Number("0"), search["rewrite_time"])
		query := search["query"].([]any)[0].(map[string]any)
		require.Equal(t, json.Number("0"), query["time_in_nanos"])
		// A count inside the breakdown is a number too, and it is zeroed with
		// the measurements: the driver states the profile, not the run.
		require.Equal(t, json.Number("0"), query["breakdown"].(map[string]any)["build_scorer_count"])
		// What the answer said about the search itself stays.
		require.Equal(t, map[string]any{"total": json.Number("1")}, got["_shards"])
	})

	t.Run("a document keeps its numbers", func(t *testing.T) {
		body := `{"_index":"alpha","_id":"1","found":true,"_source":{"id":7,"n":9223372036854775807,"took":12,"profile":{"id":"customer","time_in_nanos":42}}}`
		narrowed, ok := narrow("/alpha/_doc/1", "application/json", []byte(body))
		require.True(t, ok)
		require.Equal(t, `{"_id":"1","_index":"alpha","_source":{"id":7,"n":9223372036854775807,"profile":{"id":"customer","time_in_nanos":42},"took":12},"found":true}`, string(narrowed))
	})

	t.Run("a hit's document keeps its fields", func(t *testing.T) {
		body := `{"took":7,"hits":{"hits":[{"_id":"1","_source":{"took":"original","profile":{"id":"customer-id","time_in_nanos":42}}}]}}`
		narrowed, ok := narrow("/alpha/_search", "application/json", []byte(body))
		require.True(t, ok)
		require.Equal(t, `{"hits":{"hits":[{"_id":"1","_source":{"profile":{"id":"customer-id","time_in_nanos":42},"took":"original"}}]},"took":0}`, string(narrowed))
	})

	t.Run("plain text", func(t *testing.T) {
		_, ok := narrow("/_cat/indices", "text/plain; charset=UTF-8", []byte("alpha\nbeta\n"))
		require.False(t, ok, "a plain-text answer is held as sent")
	})
}
