package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/cosmosdb/internal/testcontainer"
)

// Microsoft's dated release tag pins the NoSQL emulator used as the SQL oracle.
// https://github.com/Azure/azure-cosmos-db-emulator-docker/releases/tag/EN20261008
var cosmosOracle struct {
	container testcontainer.Container
	endpoint  string
	err       error
}

func TestMain(m *testing.M) {
	code := m.Run()
	_ = testcontainer.TerminateContainer(cosmosOracle.container)
	os.Exit(code)
}

func cosmosRequest(t *testing.T, path string, body any, query bool) (int, []byte) {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cosmosOracle.endpoint+path, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ms-version", "2018-12-31")
	req.Header.Set("x-ms-date", time.Now().UTC().Format(http.TimeFormat))
	if query {
		req.Header.Set("Content-Type", "application/query+json")
		req.Header.Set("x-ms-documentdb-isquery", "true")
		req.Header.Set("x-ms-documentdb-query-enablecrosspartition", "true")
	}
	// The official vNext emulator does not enforce authentication.
	// https://learn.microsoft.com/en-us/azure/cosmos-db/emulator-linux#authentication
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	result, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return response.StatusCode, result
}

func TestCosmosDBQueryOracle(t *testing.T) {
	cosmosOracle.container, cosmosOracle.err = testcontainer.Start(context.Background())
	if cosmosOracle.err == nil {
		cosmosOracle.endpoint, cosmosOracle.err = cosmosOracle.container.PortEndpoint(context.Background(), "8081/tcp", "http")
	}

	if cosmosOracle.err != nil {
		t.Fatalf("Cosmos DB emulator: %v", cosmosOracle.err)
	}
	for _, setup := range []struct {
		path string
		body any
	}{
		{"/dbs", map[string]any{"id": "omni"}},
		{"/dbs/omni/colls", map[string]any{"id": "docs", "partitionKey": map[string]any{"paths": []string{"/id"}, "kind": "Hash"}}},
	} {
		status, body := cosmosRequest(t, setup.path, setup.body, false)
		if status != http.StatusCreated && status != http.StatusConflict {
			t.Fatalf("create %s: HTTP %d: %s", setup.path, status, body)
		}
	}
	for _, sql := range []string{
		"SELECT * FROM c", "SELECT VALUE c.id FROM c", "SELECT TOP 1 c.id FROM c",
		"SELECT c.id FROM c WHERE c.id = 'one'", "SELECT VALUE COUNT(1) FROM c",
		"SELECT c.id FROM c ORDER BY c.id", "SELECT VALUE IS_DEFINED(c.id) FROM c",
		"SELECT FROM", "SELECT * FROM", "SELECT VALUE FROM c", "SELECT * FROM c WHERE", "SELECT (1 + ) FROM c",
	} {
		t.Run(sql, func(t *testing.T) {
			_, parseErr := Parse(sql)
			status, body := cosmosRequest(t, "/dbs/omni/colls/docs/docs", map[string]any{"query": sql, "parameters": []any{}}, true)
			if status != http.StatusOK && status != http.StatusBadRequest {
				t.Fatalf("query %q: HTTP %d: %s", sql, status, body)
			}
			if status == http.StatusBadRequest && !strings.Contains(string(body), "Syntax error") {
				t.Fatalf("query %q has no syntax verdict: HTTP %d: %s", sql, status, body)
			}
			if (parseErr == nil) != (status == http.StatusOK) {
				t.Errorf("query %q: parser error %v, emulator HTTP %d: %s", sql, parseErr, status, body)
			}
		})
	}
}
