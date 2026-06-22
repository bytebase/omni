package cassandra

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func testdataDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "testdata", "cql", "examples")
}

var expectedFailures = map[string]string{
	"applyBatch.cql": "standalone APPLY BATCH is not valid CQL without BEGIN BATCH",
}

func TestCompatibilityHarness(t *testing.T) {
	dir := testdataDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("CQL examples corpus missing at %s: %v", dir, err)
	}

	var cqlFiles []os.DirEntry
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".cql") {
			cqlFiles = append(cqlFiles, e)
		}
	}
	if len(cqlFiles) == 0 {
		t.Fatal("CQL examples corpus is empty")
	}

	var (
		totalFiles           = len(cqlFiles)
		passedFiles          int
		expectedFailureFiles int
		totalStmts           int
	)

	for _, entry := range cqlFiles {
		t.Run(entry.Name(), func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				t.Fatal(err)
			}

			content := string(data)
			stmts, err := Parse(content)
			if err != nil {
				if reason, ok := expectedFailures[entry.Name()]; ok {
					expectedFailureFiles++
					t.Skipf("expected failure: %s (%v)", reason, err)
					return
				}
				t.Errorf("Parse failed: %v", err)
				return
			}

			passedFiles++
			totalStmts += len(stmts)

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

	t.Logf("Compatibility: %d/%d files passed, %d expected failures, %d total statements",
		passedFiles, totalFiles, expectedFailureFiles, totalStmts)
}
