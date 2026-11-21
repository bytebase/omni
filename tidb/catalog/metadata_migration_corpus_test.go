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

// TestMetadataMigrationCorpusRunsOnTiDB applies every generated migration in
// testdata/metadata_migration.yaml on the pinned TiDB server and verifies the
// migrated database converges to newSchema. The tape replays it in CI; record
// with OMNI_RECORD=1.
func TestMetadataMigrationCorpusRunsOnTiDB(t *testing.T) {
	content, err := os.ReadFile("testdata/metadata_migration.yaml")
	require.NoError(t, err)
	var cases []metadataMigrationCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	ctr, done := startContainer(t)
	defer done()
	exec := func(t *testing.T, stmt string) {
		t.Helper()
		_, err := ctr.db.ExecContext(ctr.ctx, stmt)
		require.NoError(t, err, "TiDB rejected statement:\n%s", stmt)
	}
	defer exec(t, "USE test")

	setUp := func(t *testing.T, db, script string) {
		t.Helper()
		exec(t, "DROP DATABASE IF EXISTS "+db)
		exec(t, "CREATE DATABASE "+db)
		exec(t, "USE "+db)
		exec(t, "SET foreign_key_checks = 0")
		defer exec(t, "SET foreign_key_checks = 1")
		for _, stmt := range splitStatements(script) {
			if stmt = strings.TrimSpace(stmt); stmt != "" {
				exec(t, stmt)
			}
		}
	}

	const applied, target = "omni_corpus_a", "omni_corpus_b"
	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			setUp(t, applied, tc.OldSchema)
			setUp(t, target, tc.NewSchema)
			want := storedSchema(t, ctr, target)

			oldMeta := extractTiDBMetadataForTest(t, tc.OldSchema)
			newMeta := extractTiDBMetadataForTest(t, tc.NewSchema)
			from := New()
			_, err := from.LoadMetadataInto(context.Background(), oldMeta, applied)
			require.NoError(t, err)
			to := New()
			_, err = to.LoadMetadataInto(context.Background(), newMeta, applied)
			require.NoError(t, err)

			plan := GenerateMigration(from, to, Diff(from, to))
			exec(t, "USE "+applied)
			for _, op := range plan.Ops {
				exec(t, op.SQL)
			}
			got := storedSchema(t, ctr, applied)
			require.Equal(t, want, got, "migrated schema differs\nplan:\n%s", plan.SQL())
		})
	}
}
