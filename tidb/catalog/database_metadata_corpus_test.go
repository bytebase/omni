package catalog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type databaseMetadataCorpusCase struct {
	Description string `yaml:"description"`
	Schema      string `yaml:"schema"`
	Metadata    string `yaml:"metadata"`
}

func extractTiDBMetadataForTest(t *testing.T, schemaText string) *metadatapb.DatabaseSchemaMetadata {
	t.Helper()
	const scratchDB = "_bytebase_schema_scratch"
	if strings.TrimSpace(schemaText) == "" {
		return &metadatapb.DatabaseSchemaMetadata{Schemas: []*metadatapb.SchemaMetadata{{}}}
	}
	c := New()
	_, err := c.Exec(fmt.Sprintf("SET foreign_key_checks = 0;\nCREATE DATABASE IF NOT EXISTS `%s`;\nUSE `%s`;", scratchDB, scratchDB), nil)
	require.NoError(t, err)
	results, err := c.Exec(schemaText, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "statement %d (%s)", r.Index, r.SQL)
	}
	dbName := c.CurrentDatabase()
	meta := c.DatabaseMetadata(dbName)
	if dbName == scratchDB {
		meta.Name = ""
		for _, s := range meta.GetSchemas() {
			for _, v := range s.GetViews() {
				v.Definition = strings.ReplaceAll(v.Definition, fmt.Sprintf("`%s`.", dbName), "")
			}
		}
	}
	return meta
}

// TestDatabaseMetadataCorpus pins what the catalog extracts from each TiDB DDL
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
			meta := extractTiDBMetadataForTest(t, tc.Schema)
			encoded, err := json.MarshalIndent(meta, "", "  ")
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
