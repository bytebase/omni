package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/parsertest"
)

// TestOfficialCorpusParses (googlesql/corpus-closure, official half) is the
// whole-corpus PARSE-verification gate for the 350 documented GoogleSQL forms
// scraped into docs/migration/googlesql/truth1/{bigquery,spanner}/*.md (Corpus
// B). It lifts every fenced ```sql block out of those docs and feeds each block
// verbatim to Parse — the exact entry point bytebase's Diagnose / GetQuerySpan /
// SplitSQL consumers call — asserting ZERO parse errors for every block that is
// in scope.
//
// Why parseForTest(block) and not parseSingle(segment): Parse runs the block-aware
// Split (terminator stripping + procedural BEGIN/END handling) and then parses
// each segment, which is what the consumers do. parseSingle rejects a bare
// trailing ';' and a procedural block that still carries its terminator, so a
// per-segment loop would mis-report consumer-valid input (verified: parseSingle
// rejects "SELECT 1;" while Parse accepts it). Block granularity also matches
// how the docs present a form — a single ```sql block per documented construct,
// sometimes several `;`-separated statements illustrating one feature.
//
// officialCorpusSkips is the explicit, categorized skip-list of blocks that do
// NOT parse clean, each with a precise reason. It self-prunes: a skipped block
// that starts parsing clean fails the test with a "REMOVE me" note, so the list
// cannot silently rot as the grammar grows. Buckets:
//
//   - FRAGMENT — the block is illustrative reference material (bare identifiers,
//     a lone `col TYPE` column-def snippet, a literal/EBNF example), not one or
//     more complete statements. Not a parser gap; will never parse as a stmt.
//   - OUT-OF-SCOPE — a pure-ZetaSQL / dialect-extension form the legacy bytebase
//     ANTLR grammar never implemented and no migration node owns (CREATE
//     EXTERNAL …, UNDROP, CREATE/ALTER/DROP MODEL, BigQuery-extended typed
//     literals). omni's reject preserves legacy parity.
//   - DOC-REJECTED — the block contains a statement the live Spanner oracle
//     itself REJECTS (a doc that's aspirational or shows a BigQuery-only form on
//     a Spanner page). omni's reject MATCHES the oracle — parity holds; the block
//     is skipped only because one of its statements is not universally valid.
//   - RESIDUAL GAP — a genuine union-grammar gap: the live oracle ACCEPTS the
//     form but omni rejects it. These are mirrored as flagged divergences in the
//     v2 migration store and are the to-do list for a future grammar fix. Each is
//     annotated with the oracle verdict that proves it is a real gap.
//
// Coverage: 310/350 blocks parse clean, 40 skipped (see the counts on each
// bucket below). (The vector-array column-type block spanner/ddl.md#51 graduated
// from the skip-list when googlesql/maint-vector-array-type wired the ARRAY
// vector-length parameter; spanner/expressions.md#6 stays skipped as DOC-REJECTED
// — its constructor statement is rejected by the oracle itself.)
func TestOfficialCorpusParses(t *testing.T) {
	all := collectOfficialCorpus(t)

	// Drift guard: if the truth1 corpus gains/loses a ```sql block the index
	// keys in officialCorpusSkips shift and must be re-baselined. 350 = the
	// authoring-time block count across both dialects.
	if len(all) != 350 {
		t.Errorf("found %d official ```sql blocks, expected 350 — truth1 corpus changed; re-baseline officialCorpusSkips (keys are <file>#<index>)", len(all))
	}
	// googlesql parses Spanner's dialect; the BigQuery pages are bigquery's.
	var blocks []parsertest.SQLBlock
	for _, b := range all {
		if strings.HasPrefix(b.File, "spanner/") {
			blocks = append(blocks, b)
		}
	}

	var clean, skipped int
	for _, b := range blocks {
		b := b
		t.Run(b.Label(), func(t *testing.T) {
			if reason, ok := officialCorpusSkips[b.Key()]; ok {
				skipped++
				// Self-prune: a skipped block that now parses clean must be removed.
				if _, errs := parseForTest(b.Text); len(errs) == 0 {
					t.Errorf("SKIP-LIST STALE: %s now parses clean — REMOVE it from officialCorpusSkips (was: %s)", b.Label(), reason)
				}
				return
			}
			if _, errs := parseForTest(b.Text); len(errs) > 0 {
				t.Errorf("%s produced %d parse error(s): %v\n  block:\n%s",
					b.Label(), len(errs), errs, indentBlock(b.Text))
			}
			clean++
		})
	}

	t.Logf("official corpus: %d blocks — %d clean + %d skipped", len(blocks), clean, skipped)
}

// collectOfficialCorpus resolves the repo-committed truth1 directory relative to
// this test file and returns every ```sql block from every per-dialect markdown
// page, in deterministic order. The corpus is committed under
// docs/migration/googlesql/truth1 (on main), so this gate is self-contained and
// runs in CI with no external checkout.
func collectOfficialCorpus(t *testing.T) []parsertest.SQLBlock {
	t.Helper()
	root := truth1Root(t)
	files, err := parsertest.WalkMarkdownFiles(root)
	if err != nil {
		t.Fatalf("walking truth1 corpus at %s: %v", root, err)
	}
	if len(files) == 0 {
		t.Fatalf("no truth1 markdown files under %s", root)
	}
	var blocks []parsertest.SQLBlock
	for _, rel := range files {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("reading %s: %v", rel, err)
		}
		blocks = append(blocks, parsertest.ExtractSQLBlocks(rel, string(data))...)
	}
	return blocks
}

// truth1Root returns the absolute path to docs/migration/googlesql/truth1,
// resolved from this source file's location (googlesql/parser/ → repo root is
// ../..). Fails loudly if the committed corpus is missing — unlike the legacy
// corpus, the official corpus is part of the repo and its absence is a real
// error, not a skip.
func truth1Root(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed; cannot locate truth1 corpus")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	root := filepath.Join(repoRoot, "docs", "migration", "googlesql", "truth1")
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("truth1 corpus not found at %s: %v", root, err)
	}
	return root
}

func indentBlock(s string) string {
	return "    " + strings.ReplaceAll(strings.TrimRight(s, "\n"), "\n", "\n    ")
}

// officialCorpusSkips maps "<file>#<index>" (the block's stable
// parsertest.SQLBlock.Key) to the reason it is excluded from the clean-parse
// assertion. See the four buckets documented on TestOfficialCorpusParses. Every
// RESIDUAL GAP carries the live-oracle verdict that proves omni (not the oracle)
// is wrong; every DOC-REJECTED carries the oracle's matching reject.
var officialCorpusSkips = map[string]string{
	// ============================ FRAGMENT (18) ============================
	// Reference-only blocks: identifier/path examples and bare `col TYPE`
	// column-definition snippets from the lexical & datatype docs. These are not
	// statements and never will be.
	"spanner/datatypes.md#0":  "FRAGMENT: DT-001 bare `col1 BOOL` column-def snippet, not a statement",
	"spanner/datatypes.md#1":  "FRAGMENT: DT-002 bare `col1 INT64` column-def snippet, not a statement",
	"spanner/datatypes.md#2":  "FRAGMENT: DT-003 bare `col1 FLOAT32` column-def snippet, not a statement",
	"spanner/datatypes.md#3":  "FRAGMENT: DT-004 bare `col1 FLOAT64` column-def snippet, not a statement",
	"spanner/datatypes.md#4":  "FRAGMENT: DT-005 bare `col1 NUMERIC` column-def snippet, not a statement",
	"spanner/datatypes.md#5":  "FRAGMENT: DT-006 bare `col1 STRING(100)` column-def snippet, not a statement",
	"spanner/datatypes.md#6":  "FRAGMENT: DT-007 bare `col1 BYTES(256)` column-def snippet, not a statement",
	"spanner/datatypes.md#7":  "FRAGMENT: DT-008 bare `col1 JSON` column-def snippet, not a statement",
	"spanner/datatypes.md#8":  "FRAGMENT: DT-009 bare `col1 DATE` column-def snippet, not a statement",
	"spanner/datatypes.md#9":  "FRAGMENT: DT-010 bare `col1 TIMESTAMP` column-def snippet, not a statement",
	"spanner/datatypes.md#10": "FRAGMENT: DT-011 bare `col1 ARRAY<...>` column-def snippet, not a statement",
	"spanner/datatypes.md#11": "FRAGMENT: DT-012 bare `col1 STRUCT<...>` column-def snippet, not a statement",
	"spanner/datatypes.md#12": "FRAGMENT: DT-013 TOKENLIST prose/comment-only block, not a statement",
	"spanner/datatypes.md#13": "FRAGMENT: DT-014 bare proto-typed `col1` column-def snippet, not a statement",
	"spanner/datatypes.md#14": "FRAGMENT: DT-015 bare INTERVAL literal example, not a statement",
	"spanner/datatypes.md#15": "FRAGMENT: DT-016 bare proto-enum `col1` column-def snippet, not a statement",

	// ========================== OUT-OF-SCOPE (2) ==========================
	// Pure-ZetaSQL / dialect-extension forms the legacy bytebase ANTLR grammar
	// never implemented and no migration node owns. omni's reject preserves
	// legacy parity.
	"spanner/ddl.md#38": "OUT-OF-SCOPE: ALTER MODEL … SET OPTIONS — ML model DDL, no migration node",

	// ========================== DOC-REJECTED (6) ==========================
	// One statement in the block is REJECTED by the live Spanner oracle too, so
	// omni's reject matches the oracle (parity holds). The block is skipped only
	// because that statement is not universally valid GoogleSQL.
	"spanner/ddl.md#12":         "DOC-REJECTED: DDL block — `ALTER TABLE … ALTER COLUMN … SET NOT NULL` rejected by Spanner oracle (\"ALTER COLUMN SET NOT NULL not supported without a column type\"); omni rejects too (parity). Other ALTER COLUMN statements in the block parse.",
	"spanner/ddl.md#15":         "DOC-REJECTED: DDL block — `ALTER TABLE … ALTER ROW DELETION POLICY (…)` rejected by Spanner oracle (\"Expecting 'EOF' but found 'POLICY'\" — only ADD/DROP ROW DELETION POLICY exist); omni rejects too (parity). ADD/DROP forms in the block parse.",
	"spanner/ddl.md#53":         "DOC-REJECTED: DDL-… `DEFAULT (NEXT VALUE FOR Seq)` in a column default rejected by Spanner oracle (Error parsing Spanner DDL statement); omni rejects too (parity).",
	"spanner/expressions.md#36": "DOC-REJECTED: EXPR block — `VALUES (NEXT VALUE FOR Seq, …)` rejected by Spanner oracle (Syntax error: keyword VALUE); omni rejects too (parity). GET_NEXT_SEQUENCE_VALUE(SEQUENCE …) parses, but the block also has the rejected form.",
	"spanner/expressions.md#38": "DOC-REJECTED: EXPR-… proto construction `NEW \\`pkg.Msg\\`(field: value)` rejected by Spanner oracle (Syntax error: got ':'); omni rejects too (parity).",
	"spanner/expressions.md#34": "DOC-REJECTED: EXPR block — `SELECT TIMESTAMP '…' AT TIME ZONE '…'` (bare AT TIME ZONE select item) rejected by Spanner oracle (\"Expected end of input but got AT\"); omni rejects too (parity). The EXTRACT(… AT TIME ZONE …) form in the same block parses on both.",

	// =========================== RESIDUAL GAP ==========================
	// Genuine grammar gaps: the live oracle ACCEPTS the form, omni rejects.
	// NOTE: spanner/expressions.md#33 (quantified comparison `expr {>|=|…}
	// {ALL|ANY|SOME} (subquery)`, divergence #201) was CLOSED by the
	// googlesql/maint-quantified-comparison node — it now parses clean and is
	// covered by TestExpr_QuantifiedComparison + the live differential, so its
	// skip entry was removed.
	// NOTE: spanner/ddl.md#48 (DDL-049 `ALTER SEARCH INDEX … {ADD|DROP} STORED
	// COLUMN`) was likewise removed — it is now WIRED (alter_table.go,
	// AlterSearchIndex). It is a documented Spanner form the union parser accepts;
	// the live emulator non-authoritatively rejects it (proven in the
	// the suite:bq-ddl rows of testdata/compat.tsv triangulation guard). Divergence #203.
	// DOC-REJECTED: the block's `SELECT ARRAY<FLOAT32>(vector_length => 128)`
	// statement is rejected by the live Spanner oracle itself — a value
	// constructor `ARRAY<scalar>(...)` is not a query expression form (oracle @
	// sha256:caf1bd24: "Syntax error: Expected \"[\" but got \"(\"", same reject
	// for the bare `SELECT ARRAY<FLOAT32>()`). omni's reject ("syntax error at or
	// near =>") MATCHES the oracle, so parity holds; the block is skipped only
	// because that one statement is not valid. (The vector_length parameter is a
	// COLUMN-SCHEMA feature, not an expression — its DDL form parses clean and is
	// covered by TestParseType_ArrayVectorLength + the live differential. This was
	// re-classified from a "RESIDUAL GAP" after the maint-vector-array-type node
	// re-probed the oracle: the prior divergence-#202 evidence claiming the
	// constructor ACCEPTS was stale.) The TOKENIZE_NGRAMS(... =>) statement in the
	// same block parses fine.
	"spanner/expressions.md#6": "DOC-REJECTED: `SELECT ARRAY<FLOAT32>(vector_length => 128)` — the live Spanner oracle REJECTS a typed-ARRAY value constructor (\"Expected '[' but got '('\"); omni's reject (\"syntax error at or near =>\") matches. vector_length is a column-schema parameter, not an expression form. (The first statement, TOKENIZE_NGRAMS named args, parses.)",

	// ============ RESIDUAL GAP: Spanner forms omni rejects (3) ============
	// The live Spanner oracle ACCEPTS these and omni rejects them. The union
	// parser accepted them through BigQuery's grammar (generic entities, the ML
	// OPTIONS clause), which googlesql no longer parses; Spanner's own forms need
	// a parser of their own.
	"spanner/ddl.md#37": "RESIDUAL GAP: CREATE MODEL … REMOTE OPTIONS (endpoints = […]) — oracle accepts; omni rejects the OPTIONS clause after REMOTE. Fix: parse a non-empty OPTIONS list in parseCreateModel.",
	"spanner/ddl.md#43": "RESIDUAL GAP: CREATE PLACEMENT us_placement OPTIONS (instance_partition = '…') — oracle accepts; omni has no CREATE PLACEMENT. Fix: parse CREATE PLACEMENT name [OPTIONS (…)].",
	"spanner/ddl.md#44": "RESIDUAL GAP: DROP PLACEMENT us_placement — oracle accepts (semantic: Placement not found); omni has no DROP PLACEMENT. Fix: parse DROP PLACEMENT name.",

	// =========== DOC-REJECTED: QUALIFY (1) ===========
	"spanner/query.md#22": "DOC-REJECTED: SELECT … QUALIFY rn = 1 — Spanner's docs show QUALIFY, but the pinned emulator leaves FEATURE_QUALIFY off and GoogleSQL's grammar refuses it (\"QUALIFY is not supported\", googlesql.tm:8461); omni rejects too (parity).",
}
