package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"

	"github.com/bytebase/omni/metadata"
)

type metadataCase struct {
	Description string `yaml:"description"`
	Schema      string `yaml:"schema"`
	Metadata    string `yaml:"metadata"`
}

func TestMetadataParityWithBytebaseFixtures(t *testing.T) {
	const (
		record   = false
		filename = "testdata/get_database_metadata.yaml"
	)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	var tests []metadataCase
	require.NoError(t, yaml.Unmarshal(content, &tests))
	require.NotEmpty(t, tests)

	for i, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			got, err := GetDatabaseMetadata(tc.Schema)
			require.NoError(t, err)
			encoded, err := json.MarshalIndent(got, "", "  ")
			require.NoError(t, err)
			result := string(encoded) + "\n"
			if record {
				tests[i].Metadata = result
				return
			}
			require.Equal(t, tc.Metadata, result)
		})
	}

	if record {
		writeYAMLCorpus(t, filename, tests)
	}
}

type definitionCase struct {
	Description string `yaml:"description"`
	Metadata    string `yaml:"metadata"`
	Expected    string `yaml:"expected"`
}

func TestDefinitionParityWithBytebaseFixtures(t *testing.T) {
	const (
		record   = false
		filename = "testdata/get_database_definition.yaml"
	)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	var tests []definitionCase
	require.NoError(t, yaml.Unmarshal(content, &tests))
	require.NotEmpty(t, tests)

	for i, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			meta := &metadata.DatabaseSchemaMetadata{}
			require.NoError(t, protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal([]byte(tc.Metadata), meta))
			got, err := DatabaseDefinition(meta, DefinitionOptions{})
			require.NoError(t, err)
			if record {
				tests[i].Expected = got
				return
			}
			require.Equal(t, tc.Expected, got)
		})
	}

	if record {
		writeYAMLCorpus(t, filename, tests)
	}
}

type migrationCase struct {
	Description string `yaml:"description"`
	OldSchema   string `yaml:"oldSchema"`
	NewSchema   string `yaml:"newSchema"`
	Expected    string `yaml:"expected"`
}

func TestMigrationParityWithBytebaseFixtures(t *testing.T) {
	const (
		record   = false
		filename = "testdata/generate_migration.yaml"
	)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	var tests []migrationCase
	require.NoError(t, yaml.Unmarshal(content, &tests))
	require.NotEmpty(t, tests)

	ctx := context.Background()
	for i, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			from := New()
			if tc.OldSchema != "" {
				oldMeta, err := GetDatabaseMetadata(tc.OldSchema)
				require.NoError(t, err)
				require.NoError(t, from.LoadMetadata(ctx, oldMeta))
			}
			newMeta, err := GetDatabaseMetadata(tc.NewSchema)
			require.NoError(t, err)
			to := New()
			require.NoError(t, to.LoadMetadata(ctx, newMeta))

			diff := Diff(from, to)
			got := GenerateMigration(from, to, diff).SQL()
			if record {
				tests[i].Expected = got
				return
			}
			require.Equal(t, tc.Expected, got)
		})
	}

	if record {
		writeYAMLCorpus(t, filename, tests)
	}
}

func writeYAMLCorpus(t *testing.T, filename string, cases any) {
	t.Helper()
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	require.NoError(t, enc.Encode(cases))
	require.NoError(t, enc.Close())
	require.NoError(t, os.WriteFile(filename, buf.Bytes(), 0o644))
}
