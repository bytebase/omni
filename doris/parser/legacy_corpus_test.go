package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLegacyCorpus_AllFiles runs the omni Doris parser against every SQL
// file lifted from the legacy ANTLR4 parser's example directory.
//
// Goal: verify omni parser compatibility with what the legacy parser
// supports. Every file is split into statements; each statement should
// parse without producing a ParseError. Statements producing a
// "not yet supported" error are tracked separately as compatibility gaps.
func TestLegacyCorpus_AllFiles(t *testing.T) {
	files, err := filepath.Glob("testdata/legacy/*.sql")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	regression, err := filepath.Glob("testdata/legacy/regression/*.sql")
	if err != nil {
		t.Fatalf("glob regression: %v", err)
	}
	files = append(files, regression...)

	if len(files) == 0 {
		t.Fatal("no legacy corpus files found")
	}

	type stat struct {
		file        string
		statements  int
		successful  int
		unsupported []string // collected "not yet supported" stmt names
		errors      []string // genuine parse errors (syntax issues)
	}

	var stats []stat
	var totalStmts, totalOK, totalUnsupported, totalErrors int

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		input := string(data)
		file, errs := Parse(input)

		s := stat{file: filepath.Base(path)}
		segs := Split(input)
		s.statements = len(segs)
		s.successful = len(file.Stmts)

		for _, e := range errs {
			if strings.Contains(e.Msg, "not yet supported") {
				s.unsupported = append(s.unsupported, e.Msg)
			} else {
				s.errors = append(s.errors, e.Msg)
			}
		}
		stats = append(stats, s)
		totalStmts += s.statements
		totalOK += s.successful
		totalUnsupported += len(s.unsupported)
		totalErrors += len(s.errors)
	}

	// Report — print a compatibility summary even on success
	t.Logf("=== Legacy Corpus Compatibility Report ===")
	t.Logf("Files scanned: %d", len(files))
	t.Logf("Total statements (via Split): %d", totalStmts)
	t.Logf("Successfully parsed: %d", totalOK)
	t.Logf("Unsupported (stub dispatch): %d", totalUnsupported)
	t.Logf("Genuine parse errors: %d", totalErrors)

	// Per-file detail for any file with errors or unsupported stubs
	for _, s := range stats {
		if len(s.errors) == 0 && len(s.unsupported) == 0 {
			continue
		}
		t.Logf("--- %s: %d errors, %d unsupported, %d stmts ---", s.file, len(s.errors), len(s.unsupported), s.statements)
		for _, e := range s.errors {
			t.Logf("  ERROR: %s", e)
		}
		for _, u := range s.unsupported {
			t.Logf("  UNSUP: %s", u)
		}
	}

	// Genuine parse errors are failures; unsupported stubs are accepted (they
	// mean dispatch routed to a still-stubbed handler — to be filled in later).
	if totalErrors > 0 {
		t.Errorf("legacy corpus has %d genuine parse errors", totalErrors)
	}
}

// TestLegacyCorpus_PerFile is the per-file version. Useful for seeing which
// specific files have which specific issues during debugging.
func TestLegacyCorpus_PerFile(t *testing.T) {
	files, err := filepath.Glob("testdata/legacy/*.sql")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	regression, _ := filepath.Glob("testdata/legacy/regression/*.sql")
	files = append(files, regression...)

	for _, path := range files {
		t.Run(filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			input := string(data)
			_, errs := Parse(input)

			var genuineErrors []string
			for _, e := range errs {
				if !strings.Contains(e.Msg, "not yet supported") {
					genuineErrors = append(genuineErrors, e.Msg)
				}
			}
			if len(genuineErrors) > 0 {
				for _, e := range genuineErrors {
					t.Errorf("parse error: %s", e)
				}
			}
		})
	}
}
