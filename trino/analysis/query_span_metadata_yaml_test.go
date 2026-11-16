package analysis

import (
	"os"
	"path/filepath"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

type querySpanYAMLType struct {
	Type struct {
		QueryType int `yaml:"querytype"`
	} `yaml:"type"`
}

type querySpanYAMLCase struct {
	Description         string             `yaml:"description"`
	Statement           string             `yaml:"statement"`
	DefaultDatabase     string             `yaml:"defaultDatabase"`
	IgnoreCaseSensitive bool               `yaml:"ignoreCaseSensitive"`
	Metadata            string             `yaml:"metadata"`
	QuerySpan           *querySpanYAMLType `yaml:"querySpan"`
}

func TestQuerySpanMetadataCorpus(t *testing.T) {
	paths, err := filepath.Glob("testdata/query-span/*.yaml")
	require.NoError(t, err)
	require.NotEmpty(t, paths)
	for _, path := range paths {
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		var cases []querySpanYAMLCase
		require.NoError(t, yaml.Unmarshal(content, &cases), path)
		for _, tc := range cases {
			t.Run(tc.Description, func(t *testing.T) {
				meta := &metadatapb.DatabaseSchemaMetadata{}
				require.NoError(t, protojson.Unmarshal([]byte(tc.Metadata), meta))
				getter, lister := buildMockDatabaseMetadataGetter([]*metadatapb.DatabaseSchemaMetadata{meta})
				span, err := ResolveQuerySpan(tc.Statement, ResolveOptions{
					DefaultDatabase:     tc.DefaultDatabase,
					IgnoreCaseSensitive: tc.IgnoreCaseSensitive,
					ListDatabases:       lister,
					GetDatabaseMetadata: getter,
				})
				require.NoError(t, err)
				if tc.QuerySpan != nil {
					require.Equal(t, QueryType(tc.QuerySpan.Type.QueryType), span.Type)
				}
			})
		}
	}
}
