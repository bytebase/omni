package catalog

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/metadata"
)

// definitionCase is one fixture of testdata/definition.yaml: a snapshot written as
// protojson and the restore script DatabaseDefinition renders for it.
type definitionCase struct {
	Description string `yaml:"description"`
	Metadata    string `yaml:"metadata"`
	Expected    string `yaml:"expected"`
}

// TestDatabaseDefinition pins the restore script for each snapshot and executes it
// on a fresh catalog. To regenerate, set record, run this test, and set it back.
func TestDatabaseDefinition(t *testing.T) {
	const (
		record   = false
		filename = "testdata/definition.yaml"
	)

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var cases []definitionCase
	if err := yaml.Unmarshal(content, &cases); err != nil {
		t.Fatal(err)
	}

	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			meta := &metadata.DatabaseSchemaMetadata{}
			if err := protojson.Unmarshal([]byte(tc.Metadata), meta); err != nil {
				t.Fatal(err)
			}
			got, err := DatabaseDefinition(meta, DefinitionOptions{PrintHeader: true})
			if err != nil {
				t.Fatal(err)
			}

			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d;", nil); err != nil {
				t.Fatal(err)
			}
			results, err := c.Exec(got, &ExecOptions{ContinueOnError: true})
			if err != nil {
				t.Fatalf("definition does not parse: %v", err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("definition statement %d fails: %v\n%s", r.Index, r.Error, r.SQL)
				}
			}

			if record {
				cases[i].Expected = got
				return
			}
			if got != tc.Expected {
				t.Errorf("definition mismatch\n--- got ---\n%s\n--- want ---\n%s", got, tc.Expected)
			}
		})
	}

	if record {
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(2)
		if err := enc.Encode(cases); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestPrintIndexClauseWritesComment pins the index COMMENT a sync reads from
// INDEX_COMMENT: without it the dump loses the comment on every secondary index.
func TestPrintIndexClauseWritesComment(t *testing.T) {
	var buf strings.Builder
	index := &metadata.IndexMetadata{
		Name:        "idx_v",
		Type:        "BTREE",
		Visible:     true,
		Expressions: []string{"v"},
		KeyLength:   []int64{-1},
		Comment:     "owner's note",
	}
	if err := printIndexClause(&buf, index); err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), ",\n  KEY `idx_v` (`v`) COMMENT 'owner''s note'"; got != want {
		t.Errorf("clause = %q, want %q", got, want)
	}
}
