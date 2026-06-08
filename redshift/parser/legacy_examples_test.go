package parser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

const legacyExamplesDir = "testdata/legacy"

var legacyExpectedFailures = map[string]string{
	"alter_role.sql":             `P1 Redshift role fixture edge case: zero-length EXTERNALID delimited identifier`,
	"create_role.sql":            `P1 Redshift role fixture edge case: zero-length EXTERNALID delimited identifier`,
	"drop_materialized_view.sql": `P1 Redshift materialized view drop variants: first failure near "DROP"`,
	"drop_role.sql":              `P1 Redshift DROP ROLE fixture edge case: multiple DROP ROLE statements without semicolon separators`,
	"grant.sql":                  `P1 Redshift GRANT fixture edge cases: WITH GRANT OPTION before additional grantees and zero-length schema identifier`,
	"revoke.sql":                 `P1 Redshift REVOKE fixture edge cases: zero-length schema identifier and IF EXISTS object targets`,
}

func TestLegacyExamples(t *testing.T) {
	entries, err := os.ReadDir(legacyExamplesDir)
	if err != nil {
		t.Fatalf("read legacy examples: %v", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	if len(files) != 115 {
		t.Fatalf("expected 115 legacy Redshift example files, got %d", len(files))
	}

	for _, file := range files {
		file := file
		t.Run(file, func(t *testing.T) {
			sql, err := os.ReadFile(filepath.Join(legacyExamplesDir, file))
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}
			_, parseErr := Parse(string(sql))
			reason, expectedFailure := legacyExpectedFailures[file]
			if expectedFailure {
				if parseErr == nil {
					t.Fatalf("legacy example now parses; remove expected failure %q and promote it to P0", reason)
				}
				t.Logf("expected failure: %s: %v", reason, parseErr)
				return
			}
			if parseErr != nil {
				t.Fatalf("legacy example should parse: %v", parseErr)
			}
		})
	}
}
