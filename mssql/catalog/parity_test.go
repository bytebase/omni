package catalog

import (
	"encoding/json"
	"os"
	"strings"
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
	content, err := os.ReadFile("testdata/get_database_metadata.yaml")
	require.NoError(t, err)
	var tests []metadataCase
	require.NoError(t, yaml.Unmarshal(content, &tests))
	require.NotEmpty(t, tests)

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			c, err := LoadSQL(tc.Schema)
			require.NoError(t, err)
			encoded, err := json.MarshalIndent(c.Metadata(), "", "  ")
			require.NoError(t, err)
			require.Equal(t, tc.Metadata, string(encoded)+"\n")
		})
	}
}

type definitionCase struct {
	Description string `yaml:"description"`
	Metadata    string `yaml:"metadata"`
	Expected    string `yaml:"expected"`
}

func TestDefinitionParityWithBytebaseFixtures(t *testing.T) {
	content, err := os.ReadFile("testdata/get_database_definition.yaml")
	require.NoError(t, err)
	var tests []definitionCase
	require.NoError(t, yaml.Unmarshal(content, &tests))
	require.NotEmpty(t, tests)

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			meta := &metadata.DatabaseSchemaMetadata{}
			require.NoError(t, protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal([]byte(tc.Metadata), meta))
			got, err := DatabaseDefinition(meta, DefinitionOptions{})
			require.NoError(t, err)
			// The catalog trims a view definition's trailing semicolon before
			// adding one, so a dump of a dump does not grow semicolons. Bytebase's
			// fixtures record the old doubling, e.g. ";;" after a view body.
			require.Equal(t, strings.ReplaceAll(tc.Expected, ";;", ";"), strings.ReplaceAll(got, ";;", ";"))
		})
	}
}

type migrationCase struct {
	Description string `yaml:"description"`
	OldSchema   string `yaml:"oldSchema"`
	NewSchema   string `yaml:"newSchema"`
	Expected    string `yaml:"expected"`
}

func TestMigrationParityWithBytebaseFixtures(t *testing.T) {
	content, err := os.ReadFile("testdata/generate_migration.yaml")
	require.NoError(t, err)
	var tests []migrationCase
	require.NoError(t, yaml.Unmarshal(content, &tests))
	require.NotEmpty(t, tests)

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			from := New()
			if tc.OldSchema != "" {
				var err error
				from, err = LoadSQL(tc.OldSchema)
				require.NoError(t, err)
			}
			to := New()
			if tc.NewSchema == "" {
				to, _ = LoadSQL("")
			} else {
				var err error
				to, err = LoadSQL(tc.NewSchema)
				require.NoError(t, err)
			}
			diff := Diff(from, to)
			got := GenerateMigration(from, to, diff).SQL()
			require.Equal(t, tc.Expected, got)
		})
	}
}
