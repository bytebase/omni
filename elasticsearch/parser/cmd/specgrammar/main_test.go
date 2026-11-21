package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/grammar"
)

func TestConvertMinimalSchema(t *testing.T) {
	raw := []byte(`{
  "endpoints": [
    {
      "name": "info",
      "urls": [{"methods": ["GET"], "path": "/"}]
    },
    {
      "name": "ping",
      "urls": [{"methods": ["HEAD"], "path": "/"}]
    },
    {
      "name": "search",
      "request": {"namespace": "search", "name": "Request"},
      "urls": [
        {"methods": ["GET", "POST"], "path": "/_search"},
        {"methods": ["GET", "POST"], "path": "/{index}/_search"}
      ]
    },
    {
      "name": "_internal.delete_balance",
      "urls": [{"methods": ["DELETE"], "path": "/_internal/balance/{version}"}]
    }
  ],
  "types": [
    {
      "name": {"namespace": "search", "name": "Request"},
      "query": [
        {"name": "size", "type": {"kind": "instance_of", "type": {"namespace": "_types", "name": "integer"}}},
        {"name": "timeout", "type": {"kind": "instance_of", "type": {"namespace": "_types", "name": "Duration"}}}
      ]
    }
  ]
}`)

	g4, stats, err := convert(raw, "deadbeef", "8.19", "sha256:0000")
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	if stats.endpoints != 4 {
		t.Fatalf("endpoints = %d, want 4", stats.endpoints)
	}
	if stats.distinctPaths != 4 {
		t.Fatalf("distinctPaths = %d, want 4", stats.distinctPaths)
	}
	if !strings.Contains(g4, "internalPath") {
		t.Fatalf("expected _internal namespace to become valid ANTLR rule name internalPath:\n%s", g4)
	}
	if !strings.Contains(g4, "'/_internal/balance/1'") {
		t.Fatalf("expected {version} to sample as 1:\n%s", g4)
	}
	if !strings.Contains(g4, "'/my_index/_search'") {
		t.Fatalf("expected {index} to sample as my_index:\n%s", g4)
	}
	if !strings.Contains(g4, "HEAD '/'") {
		t.Fatalf("expected HEAD '/' in headEndpoint:\n%s", g4)
	}

	g, err := grammar.ReadANTLR(g4, &grammar.Tokens{})
	if err != nil {
		t.Fatalf("ReadANTLR: %v", err)
	}
	rows, uncovered := g.Rows()
	if len(uncovered) > 0 {
		t.Fatalf("uncovered rules: %v", uncovered)
	}
	rec := g.Recognizer()
	for _, wantAccept := range []string{
		"GET /",
		"HEAD /",
		"POST /_search",
		"POST /my_index/_search",
		"DELETE /_internal/balance/1",
		"GET /_search?size=10&timeout=30s",
		"GET _search",
		"PATCH /",
	} {
		if !rec.Accepts(wantAccept) {
			t.Errorf("recognizer rejected %q", wantAccept)
		}
	}
	for _, wantReject := range []string{
		"GET",
		"/_search",
		"GET /_search zzz",
		"OPTIONS /",
		"GET zzz /",
	} {
		if rec.Accepts(wantReject) {
			t.Errorf("recognizer accepted %q", wantReject)
		}
	}
	if len(rows) == 0 {
		t.Fatalf("expected generated rows")
	}
}

func TestWritePinPreservesDriverRow(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "pins.tsv")
	initial := "# header\ndriver\tdocker.elastic.co/elasticsearch/elasticsearch:8.19.25\tsha256:aaa\t8.19.25\n"
	if err := os.WriteFile(p, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writePin(p, "abc123", "sha256:bbb", "8.19"); err != nil {
		t.Fatalf("writePin insert: %v", err)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	want := "# header\ndefault\tgithub.com/elastic/elasticsearch-specification@abc123\tsha256:bbb\t8.19\ndriver\tdocker.elastic.co/elasticsearch/elasticsearch:8.19.25\tsha256:aaa\t8.19.25\n"
	if string(got) != want {
		t.Fatalf("after insert:\n%s\nwant:\n%s", got, want)
	}
	if err := writePin(p, "def456", "sha256:ccc", "8.19"); err != nil {
		t.Fatalf("writePin replace: %v", err)
	}
	got, err = os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	want2 := "# header\ndefault\tgithub.com/elastic/elasticsearch-specification@def456\tsha256:ccc\t8.19\ndriver\tdocker.elastic.co/elasticsearch/elasticsearch:8.19.25\tsha256:aaa\t8.19.25\n"
	if string(got) != want2 {
		t.Fatalf("after replace:\n%s\nwant:\n%s", got, want2)
	}
}
