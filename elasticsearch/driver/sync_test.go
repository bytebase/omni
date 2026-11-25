package driver

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/elasticsearch/internal/testcontainer"
	metadatapb "github.com/bytebase/omni/metadata"
)

var update = flag.Bool("update", false, "rewrite the golden files under testdata/sync from this run")

// openDriver returns a driver on an endpoint emptied of the corpus of any test
// the binary ran before, so a tape holds one test's requests.
func openDriver(t *testing.T) *Driver {
	t.Helper()
	testcontainer.Reset(t)
	return &Driver{client: testcontainer.Open(t)}
}

func executeCorpus(t *testing.T, d *Driver, script string) {
	t.Helper()
	if _, err := d.Execute(context.Background(), script, driver.ExecuteOptions{}); err != nil {
		t.Fatalf("execute corpus: %v", err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// TestSyncDatabase holds the indices testdata/corpus.http creates to
// testdata/sync/database.json: their names, document counts and data sizes. The
// hidden index the corpus creates is not listed by the cat API, so the sync
// neither reports nor filters it.
func TestSyncDatabase(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))

	got, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	writeOrCompareGolden(t, filepath.Join("testdata", "sync", "database.json"), got)
}

// TestSyncInstance holds the pinned server's version, its one database and its
// users to testdata/sync/instance.json.
func TestSyncInstance(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, mustRead(t, filepath.Join("testdata", "corpus.http")))

	instance, err := d.SyncInstance(context.Background())
	require.NoError(t, err)

	databases := make([]string, 0, len(instance.Databases))
	for _, database := range instance.Databases {
		databases = append(databases, database.Name)
	}
	roles := make([]map[string]string, 0, len(instance.Roles))
	for _, role := range instance.Roles {
		attribute := ""
		if role.Attribute != nil {
			attribute = *role.Attribute
		}
		roles = append(roles, map[string]string{"name": role.Name, "attribute": attribute})
	}
	b, err := json.MarshalIndent(map[string]any{
		"version":   instance.Version,
		"databases": databases,
		"roles":     roles,
	}, "", "  ")
	require.NoError(t, err)
	writeOrCompareBytes(t, filepath.Join("testdata", "sync", "instance.json"), append(b, '\n'))
}

// TestSyncDatabaseAliases holds an alias as a synonym per index it names, so
// a request that reads an index by its alias reaches the index's masking
// policy. An alias of a hidden index is left out with the index.
func TestSyncDatabaseAliases(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, `
PUT /alpha
{"settings": {"number_of_shards": 1, "number_of_replicas": 0}}

PUT /beta
{"settings": {"number_of_shards": 1, "number_of_replicas": 0}}

PUT /hidden-one
{"settings": {"number_of_shards": 1, "number_of_replicas": 0, "index.hidden": true}}

POST /_aliases
{"actions": [{"add": {"index": "alpha", "alias": "both"}}, {"add": {"index": "beta", "alias": "both"}}, {"add": {"index": "beta", "alias": "beta-recent", "filter": {"range": {"n": {"gte": 3}}}}}, {"add": {"index": "hidden-one", "alias": "hidden-alias"}}]}
`)
	got, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	var synonyms []string
	for _, s := range got.Schemas[0].Synonyms {
		synonyms = append(synonyms, s.Name+"->"+s.TargetName)
	}
	require.Equal(t, []string{"beta-recent->beta", "both->alpha", "both->beta"}, synonyms)
}

// TestSyncClosedAndHiddenIndices pins how SyncDatabase handles a closed index
// (whose docs.count and store.size come back null from _cat/indices) alongside
// open and hidden indices against the pinned Elasticsearch server.
func TestSyncClosedAndHiddenIndices(t *testing.T) {
	d := openDriver(t)
	executeCorpus(t, d, `
PUT /beta
{"settings": {"number_of_shards": 1, "number_of_replicas": 0}}

PUT /alpha
{"settings": {"number_of_shards": 1, "number_of_replicas": 0}}

PUT /hidden-one
{"settings": {"number_of_shards": 1, "number_of_replicas": 0, "index.hidden": true}}

PUT /closed-index
{"settings": {"number_of_shards": 1, "number_of_replicas": 0}}

POST /closed-index/_close
`)

	got, err := d.SyncDatabase(context.Background())
	require.NoError(t, err)
	require.Equal(t, "node", got.Name)
	require.Len(t, got.Schemas, 1)

	tables := got.Schemas[0].Tables
	names := make([]string, len(tables))
	byName := make(map[string]*metadatapb.TableMetadata, len(tables))
	for i, table := range tables {
		names[i] = table.Name
		byName[table.Name] = table
	}
	require.Equal(t, []string{"alpha", "beta", "closed-index"}, names)
	require.Zero(t, byName["closed-index"].RowCount)
	require.Zero(t, byName["closed-index"].DataSize)
}

func writeOrCompareGolden(t *testing.T, path string, got *metadatapb.DatabaseSchemaMetadata) {
	t.Helper()
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		b, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(b, want); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("SyncDatabase differs from %s (-want +got):\n%s", path, diff)
	}
}

func writeOrCompareBytes(t *testing.T, path string, got []byte) {
	t.Helper()
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("%s differs (-want +got):\n%s", path, diff)
	}
}
