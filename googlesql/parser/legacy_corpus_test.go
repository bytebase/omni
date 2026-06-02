package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// legacyCorpusRoot is the directory holding the legacy ZetaSQL/GoogleSQL
// example .sql files (Corpus A). These are the canonical ZetaSQL parser
// testdata that the legacy ANTLR grammar was validated against. If the
// directory is missing (CI without the legacy checkout), the test skips.
const legacyCorpusRoot = "/Users/h3n4l/OpenSource/parser/googlesql/examples"

// TestLegacyCorpusTokenizes is a coverage smoke test: it tokenizes every
// legacy .sql example end-to-end and asserts the lexer produces no lex errors
// and a non-empty token stream terminated by EOF.
//
// Rationale (correctness-protocol completeness gate): the ZetaSQL corpus
// exercises the full breadth of GoogleSQL surface syntax — every literal form
// (raw/triple/bytes), all three comment styles, backtick and dashed paths,
// the pipe operator, scripting, GQL, DDL/DML/query. A lexer that cleanly
// tokenizes the entire corpus with zero spurious errors demonstrates broad
// real-world coverage beyond the hand-written unit cases. (The corpus mixes
// valid and intentionally-invalid *parse* inputs, but the invalid ones are
// syntactically/semantically malformed at the PARSER level, not the lexer
// level — they still tokenize cleanly. Should a future corpus file contain a
// genuinely lex-invalid construct, this test will surface it for triage.)
func TestLegacyCorpusTokenizes(t *testing.T) {
	if _, err := os.Stat(legacyCorpusRoot); err != nil {
		t.Skipf("legacy corpus not available at %s: %v", legacyCorpusRoot, err)
	}

	var files []string
	err := filepath.Walk(legacyCorpusRoot, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(p) == ".sql" {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking corpus: %v", err)
	}
	if len(files) == 0 {
		t.Skipf("no .sql files found under %s", legacyCorpusRoot)
	}

	totalTokens := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Errorf("reading %s: %v", f, err)
			continue
		}
		tokens, errs := Tokenize(string(data))
		if len(errs) != 0 {
			t.Errorf("%s: lexer produced %d errors (first: %q at %d)",
				filepath.Base(f), len(errs), errs[0].Msg, errs[0].Loc.Start)
			continue
		}
		if len(tokens) == 0 || tokens[len(tokens)-1].Type != tokEOF {
			t.Errorf("%s: token stream not terminated by EOF", filepath.Base(f))
			continue
		}
		totalTokens += len(tokens)
	}
	t.Logf("tokenized %d legacy corpus files, %d tokens total, 0 lex errors", len(files), totalTokens)
}
