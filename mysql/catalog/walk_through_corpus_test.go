package catalog

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
	"gopkg.in/yaml.v3"
)

type walkThroughCorpusCase struct {
	Statement    string `yaml:"statement"`
	Want         string `yaml:"want,omitempty"`
	ErrorCode    int    `yaml:"errorCode,omitempty"`
	ErrorMessage string `yaml:"errorMessage,omitempty"`
}

// TestWalkThroughCorpus pins DDL simulation (LoadMetadata -> Exec -> DatabaseMetadata)
// across the MySQL walk-through corpus relocated from Bytebase.
func TestWalkThroughCorpus(t *testing.T) {
	const (
		record   = false
		filename = "testdata/walk_through.yaml"
	)

	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	var cases []walkThroughCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("case_%02d", i+1), func(t *testing.T) {
			c := New()
			_, err := c.Exec("SET foreign_key_checks = 0;", nil)
			require.NoError(t, err)
			_, err = c.LoadMetadata(context.Background(), &metadatapb.DatabaseSchemaMetadata{Name: "test"})
			require.NoError(t, err)

			results, execErr := c.Exec(tc.Statement, &ExecOptions{ContinueOnError: true})
			var firstErr error
			if execErr != nil {
				firstErr = execErr
			} else {
				for _, r := range results {
					if r.Error != nil {
						firstErr = r.Error
						break
					}
				}
			}

			if record {
				if firstErr != nil {
					var catErr *Error
					require.ErrorAs(t, firstErr, &catErr)
					cases[i].ErrorCode = catErr.Code
					cases[i].ErrorMessage = catErr.Message
					cases[i].Want = ""
				}
				return
			}

			if tc.ErrorCode != 0 || tc.ErrorMessage != "" {
				var catErr *Error
				require.ErrorAs(t, firstErr, &catErr, "statement %q should fail with *catalog.Error", tc.Statement)
				require.Equal(t, tc.ErrorCode, catErr.Code, "statement %q", tc.Statement)
				require.Equal(t, tc.ErrorMessage, catErr.Message, "statement %q", tc.Statement)
				return
			}
			require.NoError(t, firstErr, "statement %q", tc.Statement)
			if tc.Want == "" {
				return
			}

			want := &metadatapb.DatabaseSchemaMetadata{}
			require.NoError(t, unmarshaler.Unmarshal([]byte(tc.Want), want))
			got := c.DatabaseMetadata("test")
			diff := cmp.Diff(want, got, protocmp.Transform(),
				protocmp.SortRepeatedFields(&metadatapb.DatabaseSchemaMetadata{}, "schemas"),
				protocmp.SortRepeatedFields(&metadatapb.SchemaMetadata{}, "tables", "views"),
				protocmp.SortRepeatedFields(&metadatapb.TableMetadata{}, "indexes", "columns"),
			)
			require.Empty(t, diff, "statement %q", tc.Statement)
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
