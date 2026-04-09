package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLegacyCorpus tokenizes every .sql file in testdata/legacy/ and asserts
// no lex errors and no tokInvalid tokens. This is the regression baseline:
// any failure here means the omni snowflake lexer cannot tokenize something
// the legacy ANTLR4 parser previously handled.
func TestLegacyCorpus(t *testing.T) {
	corpusDir := "testdata/legacy"
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("failed to read corpus dir %s: %v", corpusDir, err)
	}

	sqlCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		sqlCount++

		t.Run(entry.Name(), func(t *testing.T) {
			path := filepath.Join(corpusDir, entry.Name())
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			tokens, errs := Tokenize(string(data))

			if len(errs) > 0 {
				t.Errorf("%d lex errors:", len(errs))
				for _, e := range errs {
					t.Errorf("  loc=%+v msg=%q", e.Loc, e.Msg)
				}
			}

			invalidCount := 0
			for _, tok := range tokens {
				if tok.Type == tokInvalid {
					invalidCount++
				}
			}
			if invalidCount > 0 {
				t.Errorf("%d tokInvalid tokens emitted", invalidCount)
			}

			// Sanity check: stream must end with EOF.
			if len(tokens) == 0 || tokens[len(tokens)-1].Type != tokEOF {
				t.Errorf("stream did not end with tokEOF")
			}

			// Sanity check: the last non-EOF token's End should be at or
			// beyond the file's last semicolon byte position. This catches
			// the lexer halting prematurely on some content.
			lastSemi := strings.LastIndexByte(string(data), ';')
			if lastSemi >= 0 && len(tokens) >= 2 {
				lastNonEOF := tokens[len(tokens)-2]
				if lastNonEOF.Loc.End < lastSemi {
					t.Errorf("last non-EOF token Loc.End=%d < last semicolon position %d (lexer stopped early)",
						lastNonEOF.Loc.End, lastSemi)
				}
			}
		})
	}

	if sqlCount == 0 {
		t.Errorf("found 0 .sql files in %s — corpus missing?", corpusDir)
	}
	if sqlCount != 27 {
		t.Logf("note: found %d .sql files (expected 27 from the legacy lift)", sqlCount)
	}
}
