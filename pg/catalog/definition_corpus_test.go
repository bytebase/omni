package catalog

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/internal/testcontainer"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

type definitionCorpusCase struct {
	Description string `yaml:"description"`
	Metadata    string `yaml:"metadata"`
	Expected    string `yaml:"expected"`
}

// TestDatabaseDefinitionCorpus pins the DDL generated for each metadata snapshot
// and reloads it into an omni catalog to prove it is valid PostgreSQL.
func TestDatabaseDefinitionCorpus(t *testing.T) {
	const (
		record   = false
		filename = "testdata/get_database_definition.yaml"
	)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	var cases []definitionCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			var meta metadatapb.DatabaseSchemaMetadata
			require.NoError(t, unmarshaler.Unmarshal([]byte(tc.Metadata), &meta))

			definition, err := DatabaseDefinition(&meta, DefinitionOptions{PrintHeader: true})
			require.NoError(t, err)
			require.NotEmpty(t, definition)

			results, err := New().Exec(definition, &ExecOptions{ContinueOnError: true})
			require.NoError(t, err, "generated definition should parse")
			for _, r := range results {
				require.NoError(t, r.Error, "statement %d (%s)", r.Index, r.SQL)
			}

			// A YAML block scalar cannot carry a leading newline, so trim the
			// blank line the header opens with.
			definition = strings.TrimPrefix(definition, "\n")

			if record {
				cases[i].Expected = definition
				return
			}
			require.Equal(t, tc.Expected, definition)
		})
	}

	if record {
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(2)
		require.NoError(t, enc.Encode(cases))
		require.NoError(t, os.WriteFile(filename, buf.Bytes(), 0o644))
	}
}

func TestWalkThroughQuotedIndexDrop(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name:    "postgres",
		Schemas: []*metadatapb.SchemaMetadata{{Name: "public"}},
	}
	cat := New()
	_, err := cat.LoadMetadata(context.Background(), meta, LoadMetadataOptions{Full: true})
	require.NoError(t, err)

	results, err := cat.Exec(`
		CREATE TABLE public.t(a int, "B" int);
		CREATE INDEX ON t(a, "B");
		CREATE UNIQUE INDEX ON t(a, "B");
		DROP INDEX IF EXISTS t_a_B_idx_11111;
	`, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "statement %d (%s)", r.Index, r.SQL)
	}
}

// TestDatabaseDefinitionCorpusRunsOnPostgres executes every generated definition
// in testdata/get_database_definition.yaml against the pinned PostgreSQL server.
// The tape replays it in CI; record with OMNI_RECORD=1.
func TestDatabaseDefinitionCorpusRunsOnPostgres(t *testing.T) {
	content, err := os.ReadFile("testdata/get_database_definition.yaml")
	require.NoError(t, err)
	var cases []definitionCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	ctx := context.Background()
	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			var meta metadatapb.DatabaseSchemaMetadata
			require.NoError(t, unmarshaler.Unmarshal([]byte(tc.Metadata), &meta))

			definition, err := DatabaseDefinition(&meta, DefinitionOptions{PrintHeader: true})
			require.NoError(t, err)

			db := testcontainer.OpenDatabase(t, testcontainer.Default, fmt.Sprintf("def_corpus_%d", i))
			_, err = db.ExecContext(ctx, definition)
			require.NoError(t, err, "PostgreSQL rejected generated definition:\n%s", definition)
		})
	}
}
