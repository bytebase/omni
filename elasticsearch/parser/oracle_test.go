package parser

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/elasticsearch/internal/testcontainer"
)

var elasticOracle testcontainer.Container

func TestMain(m *testing.M) {
	code := m.Run()
	_ = testcontainer.TerminateContainer(elasticOracle)
	os.Exit(code)
}

func TestElasticsearchRESTOracle(t *testing.T) {
	var err error
	elasticOracle, err = testcontainer.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	endpoint, err := elasticOracle.PortEndpoint(context.Background(), "9200/tcp", "http")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, endpoint+"/omni", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Fatalf("fixture reset: HTTP %d", resp.StatusCode)
	}

	for _, script := range []string{
		"GET /",
		"PUT /omni\n{\"mappings\":{\"properties\":{\"value\":{\"type\":\"keyword\"}}}}",
		"PUT /omni/_doc/1?refresh=true\n{\"value\":\"one\"}",
		"POST /omni/_search\n{\"query\":{\"term\":{\"value\":\"one\"}}}",
		"GET /omni/_doc/1",
	} {
		t.Run(script, func(t *testing.T) {
			parsed, err := Parse(script)
			if err != nil {
				t.Fatal(err)
			}
			if len(parsed.Errors) != 0 || len(parsed.Requests) != 1 {
				t.Fatalf("parse result: %+v", parsed)
			}
			lines := NewLines(script)
			request := lines.EditorRequest(lines.AdjustedParsedRequest(parsed.Requests[0], nil))
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, request.Method, endpoint+request.URL, strings.NewReader(strings.Join(request.Data, "\n")))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				t.Fatalf("HTTP %d: %s", resp.StatusCode, body)
			}
		})
	}
}
