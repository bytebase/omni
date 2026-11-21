package catalog

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type databaseMetadataCorpusCase struct {
	Description string `yaml:"description"`
	Schema      string `yaml:"schema"`
	Metadata    string `yaml:"metadata"`
}

// TestDatabaseMetadataCorpus pins what the catalog extracts from each MySQL DDL
// schema text into DatabaseSchemaMetadata.
func TestDatabaseMetadataCorpus(t *testing.T) {
	const (
		record   = false
		filename = "testdata/get_database_metadata.yaml"
	)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	var cases []databaseMetadataCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			c := New()
			_, err := c.Exec("SET foreign_key_checks = 0;\nCREATE DATABASE IF NOT EXISTS `tmp`;\nUSE `tmp`;", nil)
			require.NoError(t, err)
			results, err := c.Exec(tc.Schema, &ExecOptions{ContinueOnError: true})
			require.NoError(t, err)
			for _, r := range results {
				require.NoError(t, r.Error, "statement %d (%s)", r.Index, r.SQL)
			}

			encoded, err := json.MarshalIndent(c.DatabaseMetadata("tmp"), "", "  ")
			require.NoError(t, err)
			result := string(encoded) + "\n"

			if record {
				cases[i].Metadata = result
				return
			}
			require.Equal(t, tc.Metadata, result)
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
