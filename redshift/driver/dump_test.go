package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/redshift/catalog"
	"github.com/bytebase/omni/redshift/parser"
)

// TestDump renders the recorded cluster snapshot testdata/sync/dev.json and
// verifies that the emitted DDL parses cleanly and executes in redshift/catalog
// without mutating the input metadata.
func TestDump(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("testdata", "sync", "dev.json"))
	if err != nil {
		t.Fatalf("read synced golden: %v", err)
	}
	m := &metadatapb.DatabaseSchemaMetadata{}
	if err := protojson.Unmarshal(body, m); err != nil {
		t.Fatalf("unmarshal synced golden: %v", err)
	}

	before := proto.Clone(m)
	var b strings.Builder
	if err := (&Driver{}).Dump(context.Background(), &b, m); err != nil {
		t.Fatalf("Dump: %v", err)
	}
	if !proto.Equal(before, m) {
		t.Errorf("Dump mutated input metadata")
	}

	ddl := b.String()
	if ddl == "" {
		t.Fatal("Dump returned empty DDL for dev.json")
	}

	// Verify FK dependency ordering across schemas ("app"."facts" references "public"."lookup").
	lookupPos := strings.Index(ddl, `CREATE TABLE "public"."lookup"`)
	factsPos := strings.Index(ddl, `CREATE TABLE "app"."facts"`)
	if lookupPos < 0 || factsPos < 0 || lookupPos > factsPos {
		t.Fatalf("expected public.lookup before app.facts:\n%s", ddl)
	}

	// Verify Redshift parser accepts the full rendered script.
	if _, err := parser.Parse(ddl); err != nil {
		t.Fatalf("parser.Parse(Dump): %v\n%s", err, ddl)
	}

	// Verify simulated execution in redshift/catalog.
	c := catalog.New()
	results, err := c.Exec(ddl, nil)
	if err != nil {
		t.Fatalf("catalog.Exec(Dump): %v\n%s", err, ddl)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("catalog.Exec statement %d failed: %v\n%s", r.Index, r.Error, ddl)
		}
	}
}
