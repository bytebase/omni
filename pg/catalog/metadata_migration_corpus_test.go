package catalog

import (
	"bytes"
	"os"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

type metadataMigrationCorpusCase struct {
	Description string `yaml:"description"`
	OldMetadata string `yaml:"oldMetadata"`
	NewMetadata string `yaml:"newMetadata"`
	Expected    string `yaml:"expected,omitempty"`
	Error       string `yaml:"error,omitempty"`
}

// TestMetadataMigrationCorpus pins MetadataMigration across the PostgreSQL and
// CockroachDB-shaped metadata diff scenarios relocated from Bytebase.
func TestMetadataMigrationCorpus(t *testing.T) {
	const (
		record   = false
		filename = "testdata/metadata_migration.yaml"
	)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	var cases []metadataMigrationCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	parseMeta := func(t *testing.T, encoded string) *metadatapb.DatabaseSchemaMetadata {
		t.Helper()
		meta := &metadatapb.DatabaseSchemaMetadata{}
		require.NoError(t, unmarshaler.Unmarshal([]byte(encoded), meta))
		return meta
	}

	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			oldMeta := parseMeta(t, tc.OldMetadata)
			newMeta := parseMeta(t, tc.NewMetadata)

			migration, err := testMetadataMigration(oldMeta, newMeta)
			if record {
				cases[i].Expected = migration
				cases[i].Error = ""
				if err != nil {
					cases[i].Expected = ""
					cases[i].Error = err.Error()
				}
				return
			}
			if tc.Error != "" {
				require.ErrorContains(t, err, tc.Error)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.Expected, migration)
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
