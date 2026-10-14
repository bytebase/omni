package review

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/review"
)

// TestParserCorpus reviews every script of the parser's corpus, the
// documentation's examples among them, with every rule and every request:
// the review never fails, every range lies in the text, and a script has
// at most one Syntax finding.
func TestParserCorpus(t *testing.T) {
	var files []string
	err := filepath.WalkDir("../parser/testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no corpus scripts")
	}
	opts := review.Options{Change: review.Change{OnlineMigration: true, PriorBackup: true}}
	targets := []review.Target{{Schema: shop()}, {}}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		sql := string(data)
		result, err := Review(context.Background(), sql, opts, targets)
		if err != nil {
			t.Errorf("%s: Review: %v", f, err)
			continue
		}
		syntax := 0
		for _, rng := range result.Statements {
			if rng.Start < 0 || rng.End > len(sql) || rng.Start >= rng.End {
				t.Errorf("%s: statement range %v", f, rng)
			}
		}
		for _, fd := range result.Findings {
			if fd.Rule == review.Syntax {
				syntax++
			}
			if fd.Range.Start < 0 || fd.Range.End > len(sql) || fd.Range.Start > fd.Range.End {
				t.Errorf("%s: %s finding range %v", f, fd.Rule, fd.Range)
			}
			if fd.Statement >= len(result.Statements) {
				t.Errorf("%s: %s finding statement %d of %d", f, fd.Rule, fd.Statement, len(result.Statements))
			}
		}
		if syntax > 1 {
			t.Errorf("%s: %d syntax findings", f, syntax)
		}
	}
}
