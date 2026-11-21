package catalog

import (
	"bytes"
	"os"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

type definitionCorpusCase struct {
	Description string `yaml:"description"`
	Metadata    string `yaml:"metadata"`
	Expected    string `yaml:"expected"`
}

// requireDelimiterDiscipline checks what a mysql client relies on to restore a
// dump: a line ending in ";;" terminates a statement only while a DELIMITER
// directive has ";;" active.
func requireDelimiterDiscipline(t *testing.T, definition string) {
	t.Helper()

	active := ";"
	for i, line := range strings.Split(definition, "\n") {
		trimmed := strings.TrimSpace(line)
		if directive, ok := strings.CutPrefix(trimmed, "DELIMITER "); ok {
			active = strings.TrimSpace(directive)
			continue
		}
		if strings.HasSuffix(trimmed, ";;") && active != ";;" {
			t.Fatalf("line %d ends a statement with %q while %q is the active delimiter: %s",
				i+1, ";;", active, trimmed)
		}
	}
}

// TestDatabaseDefinitionCorpus pins the DDL generated for each metadata input
// and verifies it parses back through the catalog (or maintains DELIMITER
// discipline for routine scripts).
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

			if strings.Contains(definition, "DELIMITER") {
				requireDelimiterDiscipline(t, definition)
			} else {
				c := New()
				_, err = c.Exec("SET foreign_key_checks = 0;\nCREATE DATABASE IF NOT EXISTS `tmp`;\nUSE `tmp`;", nil)
				require.NoError(t, err)
				results, err := c.Exec(definition, &ExecOptions{ContinueOnError: true})
				require.NoError(t, err, "generated definition should parse")
				for _, r := range results {
					require.NoError(t, r.Error, "statement %d (%s)", r.Index, r.SQL)
				}
			}

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

// TestDatabaseMetadataParsesPrettyView proves the catalog parses the multi-line
// CREATE OR REPLACE VIEW emitted by SDLFormat and canonicalizes it back to one line.
func TestDatabaseMetadataParsesPrettyView(t *testing.T) {
	sdl, err := DatabaseDefinition(&metadatapb.DatabaseSchemaMetadata{
		Name: "testdb",
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{
				Name:    "t1",
				Engine:  "InnoDB",
				Columns: []*metadatapb.ColumnMetadata{{Name: "a", Type: "int"}, {Name: "b", Type: "varchar(10)", Nullable: true, Default: "NULL"}},
				Indexes: []*metadatapb.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Expressions: []string{"a"}, Visible: true}},
			}},
			Views: []*metadatapb.ViewMetadata{{
				Name:       "v_pretty",
				Definition: "select `testdb`.`t1`.`a` AS `a`,`testdb`.`t1`.`b` AS `b` from `testdb`.`t1` where (`testdb`.`t1`.`a` > 0)",
			}},
		}},
	}, DefinitionOptions{SDLFormat: true})
	require.NoError(t, err)
	require.Contains(t, sdl, "AS\nselect\n  `t1`.`a` AS `a`,")

	c := New()
	_, err = c.Exec("SET foreign_key_checks = 0;\nCREATE DATABASE IF NOT EXISTS `tmp`;\nUSE `tmp`;", nil)
	require.NoError(t, err)
	results, err := c.Exec(sdl, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	for _, r := range results {
		require.NoError(t, r.Error, "statement %d (%s)", r.Index, r.SQL)
	}
	meta := c.DatabaseMetadata("tmp")
	require.Len(t, meta.Schemas, 1)
	require.Len(t, meta.Schemas[0].Views, 1)
	require.Equal(t, "v_pretty", meta.Schemas[0].Views[0].Name)
	require.NotContains(t, meta.Schemas[0].Views[0].Definition, "\n")
}
