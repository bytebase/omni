package cassandra

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const cqlExamplesDir = "/Users/rebeliceyang/Github/parser/cql/examples"

var expectedFailures = map[string]string{
	"applyBatch.cql": "standalone APPLY BATCH is not valid CQL without BEGIN BATCH",
}

func TestCompatibilityHarness(t *testing.T) {
	entries, err := os.ReadDir(cqlExamplesDir)
	if err != nil {
		t.Skipf("CQL examples directory not found: %v", err)
	}

	var totalStmts, passedStmts int
	var failures []string

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".cql") {
			continue
		}
		t.Run(entry.Name(), func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(cqlExamplesDir, entry.Name()))
			if err != nil {
				t.Fatal(err)
			}

			content := string(data)
			stmts, err := Parse(content)
			if err != nil {
				if reason, ok := expectedFailures[entry.Name()]; ok {
					t.Skipf("expected failure: %s (%v)", reason, err)
					return
				}
				failures = append(failures, entry.Name()+": "+err.Error())
				t.Errorf("Parse failed: %v", err)
				return
			}

			totalStmts += len(stmts)
			passedStmts += len(stmts)

			for i, s := range stmts {
				if s.AST == nil {
					t.Errorf("statement %d has nil AST", i)
				}
				loc := s.AST.GetLoc()
				if loc.Start < 0 || loc.End <= loc.Start {
					t.Errorf("statement %d has invalid Loc: %+v", i, loc)
				}
			}

			violations := CheckLocations(t, content)
			for _, v := range violations {
				t.Errorf("Loc violation: %s", v)
			}
		})
	}

	t.Logf("Compatibility: %d/%d statements parsed from %d files", passedStmts, totalStmts, len(entries))
	if len(failures) > 0 {
		t.Logf("Failures:\n  %s", strings.Join(failures, "\n  "))
	}
}
