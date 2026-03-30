package parsertest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	mongo "github.com/bytebase/omni/mongo"
)

// TestANTLRExampleFiles reads every .js file from the ANTLR example directory,
// parses it through mongo.Parse(), and validates that:
//  1. No error is returned.
//  2. At least one statement is produced.
//  3. Every statement has valid position tracking.
func TestANTLRExampleFiles(t *testing.T) {
	dir := "/Users/h3n4l/OpenSource/parser/mongodb/examples"

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("cannot read examples directory: %v", err)
	}

	var jsFiles []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".js") {
			jsFiles = append(jsFiles, e)
		}
	}
	if len(jsFiles) == 0 {
		t.Fatal("no .js files found in examples directory")
	}

	t.Logf("Found %d .js example files", len(jsFiles))

	passed := 0
	failed := 0

	for _, f := range jsFiles {
		name := f.Name()
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("cannot read file: %v", err)
			}
			content := string(data)

			stmts, err := mongo.Parse(content)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			if len(stmts) == 0 {
				t.Fatal("expected at least one statement, got 0")
			}

			for i, s := range stmts {
				if s.ByteStart < 0 {
					t.Errorf("stmt[%d] ByteStart = %d, want >= 0", i, s.ByteStart)
				}
				if s.ByteEnd <= s.ByteStart {
					t.Errorf("stmt[%d] ByteEnd = %d <= ByteStart = %d", i, s.ByteEnd, s.ByteStart)
				}
				if s.Start.Line < 1 {
					t.Errorf("stmt[%d] Start.Line = %d, want >= 1", i, s.Start.Line)
				}
				if s.Start.Column < 1 {
					t.Errorf("stmt[%d] Start.Column = %d, want >= 1", i, s.Start.Column)
				}
			}
		})
	}

	t.Logf("Results: %d passed, %d failed out of %d total", passed, failed, len(jsFiles))
}
