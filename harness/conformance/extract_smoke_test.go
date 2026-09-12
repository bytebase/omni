package main

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestSmokeRealCorpus guards extraction health against the real TiDB corpus,
// fetching it (a 2 s sparse clone) when it is not there yet.
func TestSmokeRealCorpus(t *testing.T) {
	if _, err := os.Stat("corpus/tidb"); err != nil {
		if out, err := exec.Command("./fetch_corpus.sh").CombinedOutput(); err != nil {
			t.Fatalf("fetch corpus: %v\n%s", err, out)
		}
	}
	entries, err := extractTiDBCorpus("corpus")
	if err != nil {
		t.Fatal(err)
	}
	var accepts, rejects, skips int
	for _, e := range entries {
		switch {
		case e.SkipReason != "":
			skips++
		case e.Expected == VerdictAccept:
			accepts++
		default:
			rejects++
		}
	}
	t.Logf("extracted %d entries: %d accept / %d reject / %d skip", len(entries), accepts, rejects, skips)
	if len(entries) < 3000 {
		t.Fatalf("suspiciously few entries: %d", len(entries))
	}
	if skips > len(entries)/10 {
		t.Fatalf("skip rate over 10%%: %d/%d", skips, len(entries))
	}
	if n := countTestCaseElements(t); n != len(entries) {
		t.Fatalf("zero-loss cross-check failed: %d []testCase elements in corpus, %d extracted entries", n, len(entries))
	}
}

// countTestCaseElements independently counts every []testCase{...} slice
// element plus every bare testCase{...} literal (append-form rows) across the
// corpus test files with a flat whole-file walk — unlike the extractor's
// FuncDecl scope, it would also see package-level tables — so extraction can
// never silently lose rows on a future corpus tag bump.
func countTestCaseElements(t *testing.T) int {
	t.Helper()
	files, err := filepath.Glob(filepath.Join("corpus", "tidb", "pkg", "parser", "*_test.go"))
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		f, err := goparser.ParseFile(token.NewFileSet(), path, src, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			cl, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			// Mirror the extractor's arms; return false on all so slice
			// elements are never double-counted as bare literals.
			switch {
			case isTestCaseSlice(cl), isTestErrMsgCaseSlice(cl):
				count += len(cl.Elts)
				return false
			case isTestCaseLit(cl):
				count++
				return false
			}
			return true
		})
	}
	return count
}
