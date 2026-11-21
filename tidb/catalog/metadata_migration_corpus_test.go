package catalog

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/parser"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type metadataMigrationCorpusCase struct {
	Description string `yaml:"description"`
	OldSchema   string `yaml:"oldSchema"`
	NewSchema   string `yaml:"newSchema"`
	Expected    string `yaml:"expected"`
}

// TestMetadataMigrationCorpus pins LoadMetadataInto -> Diff -> GenerateMigration
// across the TiDB migration corpus relocated from Bytebase.
func TestMetadataMigrationCorpus(t *testing.T) {
	const (
		record      = false
		filename    = "testdata/metadata_migration.yaml"
		syntheticDB = "bbcatalog"
	)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	var cases []metadataMigrationCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for i, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			oldMeta := extractTiDBMetadataForTest(t, tc.OldSchema)
			newMeta := extractTiDBMetadataForTest(t, tc.NewSchema)

			from := New()
			_, err := from.LoadMetadataInto(context.Background(), oldMeta, syntheticDB)
			require.NoError(t, err)
			to := New()
			_, err = to.LoadMetadataInto(context.Background(), newMeta, syntheticDB)
			require.NoError(t, err)

			var migration string
			if diff := Diff(from, to); !diff.IsEmpty() {
				migration = strings.ReplaceAll(GenerateMigration(from, to, diff).SQL(), "`"+syntheticDB+"`.", "")
			}

			if strings.TrimSpace(migration) != "" {
				_, err = parser.Parse(migration)
				require.NoError(t, err, "generated migration should parse:\n%s", migration)
			}

			if record {
				cases[i].Expected = migration
				return
			}
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
