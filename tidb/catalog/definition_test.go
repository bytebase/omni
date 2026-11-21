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

// TestViewDefinitionStatesColumnList pins where a dump states a view's explicit
// column list: the engine stores the body under the select's own names and keeps
// the list beside it, so only the list makes the restored body yield the view's
// columns. A body that names them itself needs none.
func TestViewDefinitionStatesColumnList(t *testing.T) {
	for _, tt := range []struct {
		name       string
		definition string
		columns    []string
		want       string
	}{
		{
			name:       "other columns",
			definition: "select `t`.`id` AS `id` from `t`",
			columns:    []string{"renamed_id"},
			want:       " (`renamed_id`)",
		},
		{
			name:       "the same columns",
			definition: "select `t`.`id` AS `id`,`t`.`v` AS `v` from `t`",
			columns:    []string{"id", "v"},
		},
		{
			name:       "an unaliased column",
			definition: "SELECT id FROM t",
			columns:    []string{"id"},
		},
		{
			name:       "a name's case",
			definition: "select `t`.`id` AS `id` from `t`",
			columns:    []string{"ID"},
			want:       " (`ID`)",
		},
		{
			name:       "an unnamed expression",
			definition: "select (`t`.`id` + 1) AS `id + 1` from `t`",
			columns:    []string{"total"},
			want:       " (`total`)",
		},
		{
			name:       "a set operation",
			definition: "select `t`.`id` AS `id` from `t` union all select `u`.`id` AS `id` from `u`",
			columns:    []string{"a"},
			want:       " (`a`)",
		},
		{
			name:       "a star",
			definition: "SELECT * FROM t",
			columns:    []string{"a", "b"},
			want:       " (`a`,`b`)",
		},
		{
			name:       "WITH CHECK OPTION",
			definition: "select `t`.`id` AS `id` from `t` where (`t`.`id` > 0) WITH CASCADED CHECK OPTION",
			columns:    []string{"c"},
			want:       " (`c`)",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			view := &metadata.ViewMetadata{Name: "v", Definition: tt.definition}
			for _, name := range tt.columns {
				view.Columns = append(view.Columns, &metadata.ColumnMetadata{Name: name})
			}
			got, err := ViewDefinition(view)
			if err != nil {
				t.Fatal(err)
			}
			if tt.want == "" {
				if strings.Contains(got, "`v` (") {
					t.Errorf("a body that names its own columns states no list:\n%s", got)
				}
				return
			}
			if !strings.Contains(got, "`v`"+tt.want+" AS ") {
				t.Errorf("the list %q is missing:\n%s", tt.want, got)
			}
		})
	}
}

func TestDatabaseDefinitionAppliesOptionsWithoutSnapshotName(t *testing.T) {
	got, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
		Name:                  "source_db",
		PlacementPolicy:       "p_fast",
		TiflashReplica:        2,
		TiflashLocationLabels: []string{"zone", "host"},
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name:    "t",
				Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int"}},
			}},
		}},
	}, DefinitionOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "source_db") {
		t.Fatalf("DatabaseDefinition targeted snapshot name instead of selected restore db:\n%s", got)
	}
	if !strings.Contains(got, "ALTER DATABASE PLACEMENT POLICY=`p_fast`;") {
		t.Fatalf("missing nameless placement ALTER DATABASE:\n%s", got)
	}
	if !strings.Contains(got, "ALTER DATABASE SET TIFLASH REPLICA 2 LOCATION LABELS 'zone', 'host';") {
		t.Fatalf("missing nameless tiflash ALTER DATABASE:\n%s", got)
	}

	c := New()
	for _, stmt := range []string{
		"CREATE PLACEMENT POLICY `p_fast` FOLLOWERS=2",
		"CREATE DATABASE `target_db`",
		"USE `target_db`",
	} {
		if _, err := c.Exec(stmt, nil); err != nil {
			t.Fatal(err)
		}
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
	db := c.GetDatabase("target_db")
	if db == nil {
		t.Fatal("target_db missing")
	}
	if db.PlacementPolicy != "p_fast" || db.TiFlashReplica != 2 || len(db.TiFlashLocationLabels) != 2 {
		t.Fatalf("options not applied to selected target_db: policy=%q replica=%d labels=%v", db.PlacementPolicy, db.TiFlashReplica, db.TiFlashLocationLabels)
	}
}
