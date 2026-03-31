package parser_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/cosmosdb/ast"
	"github.com/bytebase/omni/cosmosdb/parser"
)

var update = flag.Bool("update", true, "update golden files")

func TestParseExamples(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no test SQL files found in testdata/")
	}

	goldenDir := filepath.Join("testdata", "golden")
	if *update {
		if err := os.MkdirAll(goldenDir, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	for _, file := range files {
		name := filepath.Base(file)
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			sql := string(data)

			list, err := parser.Parse(sql)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			got := ast.NodeToString(list)

			goldenFile := filepath.Join(goldenDir, strings.TrimSuffix(name, ".sql")+".txt")

			if *update {
				if err := os.WriteFile(goldenFile, []byte(got+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}

			want, err := os.ReadFile(goldenFile)
			if err != nil {
				t.Fatalf("Golden file not found: %s (run with -update to create)", goldenFile)
			}

			if got+"\n" != string(want) {
				t.Errorf("AST output mismatch for %s.\nGot:\n%s\nWant:\n%s", name, got, string(want))
			}
		})
	}
}
