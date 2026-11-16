package catalog

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/parser"
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

// TestSequenceDefinitionWritesTiDBsForm pins the CREATE SEQUENCE a dump writes
// for a synced sequence: TiDB's own SHOW CREATE SEQUENCE form, which both the
// engine and omni's parser accept, so a restored database keeps the sequence.
func TestSequenceDefinitionWritesTiDBsForm(t *testing.T) {
	sequence := &metadata.SequenceMetadata{
		Name:      "seq_orders",
		DataType:  "BIGINT",
		Start:     "10",
		MinValue:  "1",
		MaxValue:  "10000",
		Increment: "2",
		Cycle:     true,
		CacheSize: "50",
		Comment:   "order's seq",
	}
	got, err := SequenceDefinition(sequence)
	if err != nil {
		t.Fatal(err)
	}
	want := "--\n-- Sequence structure for `seq_orders`\n--\n" +
		"CREATE SEQUENCE `seq_orders` start with 10 minvalue 1 maxvalue 10000 increment by 2 cache 50 cycle ENGINE=InnoDB COMMENT='order''s seq';\n\n"
	if got != want {
		t.Fatalf("SequenceDefinition = %q, want %q", got, want)
	}
	if _, err := parser.Parse(strings.TrimSuffix(strings.TrimPrefix(got, "--\n-- Sequence structure for `seq_orders`\n--\n"), ";\n\n")); err != nil {
		t.Fatalf("TiDB's parser rejects the dumped sequence: %v", err)
	}
}

// TestDatabaseDefinitionWritesSequencesBeforeTables pins the order: a table's
// default may read a sequence, so the dump creates them first.
func TestDatabaseDefinitionWritesSequencesBeforeTables(t *testing.T) {
	got, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
		Name: "d",
		Schemas: []*metadata.SchemaMetadata{{
			Tables:    []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}}}},
			Sequences: []*metadata.SequenceMetadata{{Name: "s", Start: "1", Increment: "1"}},
		}},
	}, DefinitionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	seq, table := strings.Index(got, "CREATE SEQUENCE `s`"), strings.Index(got, "CREATE TABLE `t`")
	if seq < 0 || table < 0 || seq > table {
		t.Fatalf("sequence at %d, table at %d:\n%s", seq, table, got)
	}
}
