package parser

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// legacyCorpusRoot returns the absolute path to testdata/legacy, the committed
// mirror of the legacy ZetaSQL/GoogleSQL example .sql files (Corpus A, lifted
// from bytebase/parser/googlesql/examples; see testdata/legacy/README.md for
// provenance). These are the canonical ZetaSQL parser testdata that the legacy
// ANTLR grammar was validated against. Resolved relative to the test working
// directory (go test runs each package's tests from the package directory),
// not via runtime.Caller, so it also works under -trimpath where source paths
// are rewritten to import paths. Fails loudly if the corpus is missing: the
// corpus is part of the repo, so its absence is a real error, never a skip.
func legacyCorpusRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("testdata", "legacy"))
	if err != nil {
		t.Fatalf("resolving legacy corpus path: %v", err)
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("legacy corpus not found at %s: %v", root, err)
	}
	return root
}

// legacyCorpusFiles returns every .sql file under the legacy corpus, sorted.
// It fails if the corpus holds no .sql files at all.
func legacyCorpusFiles(t *testing.T) (root string, files []string) {
	t.Helper()
	root = legacyCorpusRoot(t)
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
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
		t.Fatalf("no .sql files found under %s", root)
	}
	sort.Strings(files)
	return root, files
}

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
	_, files := legacyCorpusFiles(t)

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

// TestLegacyCorpusParses (googlesql/corpus-closure, legacy half) is the
// whole-corpus PARSE-verification gate for the 72 legacy ZetaSQL example .sql
// files (Corpus A). For every file it runs Parse on the entire file (Parse
// splits into statements and parses each — the bytebase consumer path) and
// asserts ZERO parse errors, EXCEPT for files listed in legacyCorpusParseSkips.
//
// Relationship to TestLegacyCorpusTokenizes above: that test proves the LEXER
// cleanly tokenizes all 72 files (no lex errors); this test proves the PARSER
// accepts every file's statements that fall inside the migrated grammar's scope.
// The two together are the legacy-corpus closure gate.
//
// The legacy ZetaSQL corpus is a SUPERSET of what the legacy bytebase ANTLR
// grammar — and therefore the omni port — ever implemented: it includes
// pure-ZetaSQL statements (IMPORT MODULE / MODULE, DEFINE TABLE, CREATE
// CONSTANT, CREATE/DROP EXTERNAL …, CREATE MODEL with INPUT/OUTPUT/TRANSFORM,
// SHOW, ALTER ROW ACCESS POLICY, ALTER MATERIALIZED/APPROX VIEW / MODEL …) that
// no migration node owns. Files containing such a statement are skipped with a
// precise reason. The skip-list self-prunes: a skipped file that now parses
// clean fails with a "REMOVE me" note, so it cannot rot as the grammar grows.
//
// Coverage at authoring time: 60/72 files parse clean, 12 skipped.
func TestLegacyCorpusParses(t *testing.T) {
	root, files := legacyCorpusFiles(t)

	// Drift guard: the skip-list keys are paths relative to the corpus root.
	if len(files) != 72 {
		t.Errorf("found %d legacy .sql files, expected 72 — corpus changed; re-baseline legacyCorpusParseSkips", len(files))
	}

	var clean, skipped int
	for _, f := range files {
		rel, _ := filepath.Rel(root, f)
		rel = filepath.ToSlash(rel)
		t.Run(rel, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatalf("reading %s: %v", rel, err)
			}
			_, errs := Parse(string(data))

			if reason, ok := legacyCorpusParseSkips[rel]; ok {
				skipped++
				if len(errs) == 0 {
					t.Errorf("SKIP-LIST STALE: %s now parses clean — REMOVE it from legacyCorpusParseSkips (was: %s)", rel, reason)
				}
				return
			}
			if len(errs) > 0 {
				t.Errorf("%s produced %d parse error(s) (first: %q)", rel, len(errs), errs[0].Msg)
			}
			clean++
		})
	}

	t.Logf("legacy corpus parse: %d files — %d clean + %d skipped", len(files), clean, skipped)
}

// legacyCorpusParseSkips maps a legacy .sql file (path relative to
// testdata/legacy, slash-separated) to the reason its whole-file Parse does not
// succeed. All entries are one of:
//
//   - OUT-OF-SCOPE — the file contains a pure-ZetaSQL / dialect-extension
//     statement no migration node owns (the legacy bytebase ANTLR grammar may
//     have a rule for it, but it is P1 long-tail outside the cutover DAG). omni's
//     reject preserves the implemented scope.
//   - RESIDUAL GAP — the file is mostly in-scope but trips one genuine grammar
//     gap; annotated with the failing form. A future grammar fix clears it.
//
// Each reason names the statement(s) that defeat the file so the skip is
// auditable and the self-prune check stays meaningful.
var legacyCorpusParseSkips = map[string]string{
	// --- OUT-OF-SCOPE: pure-ZetaSQL statements, no migration node ---
	"zetasql/parser/testdata/modules.sql":                 "OUT-OF-SCOPE: IMPORT MODULE / MODULE (ZetaSQL module system) — no migration node",
	"zetasql/parser/testdata/define_table.sql":            "OUT-OF-SCOPE: DEFINE TABLE (ZetaSQL) — no migration node",
	"zetasql/parser/testdata/create_constant.sql":         "OUT-OF-SCOPE: CREATE CONSTANT (ZetaSQL) — no migration node",
	"zetasql/parser/testdata/create_model.sql":            "OUT-OF-SCOPE: CREATE MODEL … INPUT/OUTPUT/TRANSFORM (ZetaSQL ML) — no migration node",
	"zetasql/parser/testdata/create_external_table.sql":   "OUT-OF-SCOPE: CREATE [PRIVATE] EXTERNAL TABLE (ZetaSQL) — no migration node",
	"zetasql/parser/testdata/create_external_schema.sql":  "OUT-OF-SCOPE: CREATE EXTERNAL SCHEMA (ZetaSQL) — no migration node",
	"zetasql/parser/testdata/show.sql":                    "OUT-OF-SCOPE: SHOW TABLES/COLUMNS/INDEXES/STATUS/VARIABLES (ZetaSQL) — no migration node",
	"zetasql/parser/testdata/drop.sql":                    "OUT-OF-SCOPE: DROP EXTERNAL TABLE / DROP CONNECTION (ZetaSQL long-tail) defeat the file; the plain DROP TABLE/MATERIALIZED VIEW/SNAPSHOT TABLE statements parse",
	"zetasql/parser/testdata/alter_row_access_policy.sql": "OUT-OF-SCOPE: ALTER [ALL] ROW ACCESS POLICY/POLICIES (only CREATE ROW ACCESS POLICY is in scope) — no migration node for the ALTER actions",
	"zetasql/parser/testdata/alter_set_options.sql":       "OUT-OF-SCOPE: ALTER MATERIALIZED VIEW / APPROX VIEW / MODEL … SET OPTIONS (multi-action, BigQuery long-tail) — no migration node",

	// --- RESIDUAL GAP: mostly in-scope, one form trips a grammar gap ---
	"zetasql/parser/testdata/create_index.sql":      "RESIDUAL GAP: ZetaSQL CREATE INDEX … AS foo UNNEST(...) (search-index expression syntax) and CREATE [SEARCH|VECTOR] INDEX OPTIONS forms defeat the file; the plain CREATE [UNIQUE] INDEX … [STORING] statements parse",
	"zetasql/parser/testdata/alter_column_type.sql": "RESIDUAL GAP: BigQuery `ALTER COLUMN … SET DATA TYPE <type> NOT NULL OPTIONS(...)` (trailing OPTIONS after a type+NOT NULL) is the sole failing statement of 16 — Spanner oracle rejects it too (BigQuery-only form); the SET DATA TYPE / SET OPTIONS / DROP NOT NULL / DROP GENERATED / COLLATE statements parse",
}
