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
// Why Parse(block) and not parseSingle(segment): Parse runs the block-aware
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
// Coverage at authoring time: 307/350 blocks parse clean, 43 skipped (see the
// counts on each bucket below).
func TestOfficialCorpusParses(t *testing.T) {
	blocks := collectOfficialCorpus(t)

	// Drift guard: if the truth1 corpus gains/loses a ```sql block the index
	// keys in officialCorpusSkips shift and must be re-baselined. 350 = the
	// authoring-time block count across both dialects.
	if len(blocks) != 350 {
		t.Errorf("found %d official ```sql blocks, expected 350 — truth1 corpus changed; re-baseline officialCorpusSkips (keys are <file>#<index>)", len(blocks))
	}

	var clean, skipped int
	for _, b := range blocks {
		b := b
		t.Run(b.Label(), func(t *testing.T) {
			if reason, ok := officialCorpusSkips[b.Key()]; ok {
				skipped++
				// Self-prune: a skipped block that now parses clean must be removed.
				if _, errs := Parse(b.Text); len(errs) == 0 {
					t.Errorf("SKIP-LIST STALE: %s now parses clean — REMOVE it from officialCorpusSkips (was: %s)", b.Label(), reason)
				}
				return
			}
			if _, errs := Parse(b.Text); len(errs) > 0 {
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
	"bigquery/lexical.md#0":   "FRAGMENT: LEX-001 bare unquoted-identifier examples (abc5, _5abc, dataField), not statements",
	"bigquery/lexical.md#2":   "FRAGMENT: LEX-003 path-expression examples (foo.bar, foo/bar:25), not statements",
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

	// ========================== OUT-OF-SCOPE (8) ==========================
	// Pure-ZetaSQL / dialect-extension forms the legacy bytebase ANTLR grammar
	// never implemented and no migration node owns. omni's reject preserves
	// legacy parity.
	"bigquery/datatypes.md#0": "OUT-OF-SCOPE: DT-001 BigQuery-extended typed literals (INT64 '42', BOOL 'TRUE', STRING 'hello', BYTES, FLOAT64) — legacy .g4 only has DATE/TIME/DATETIME/TIMESTAMP/NUMERIC/JSON typed literals; live Spanner oracle also REJECTS `SELECT INT64 '42'` (Syntax error). NUMERIC/DATE/TIMESTAMP literals in the same block parse; the INT64/BOOL/STRING ones do not.",
	"bigquery/ddl.md#12":      "OUT-OF-SCOPE: DDL-013 CREATE EXTERNAL SCHEMA — no migration node (legacy create_external_schema_statement is a ZetaSQL long-tail form not in the P0/P1 DAG)",
	"bigquery/ddl.md#13":      "OUT-OF-SCOPE: DDL-014 CREATE [OR REPLACE] EXTERNAL TABLE — no migration node (legacy create_external_table_statement, ZetaSQL long-tail)",
	"bigquery/ddl.md#42":      "OUT-OF-SCOPE: UNDROP SCHEMA — no migration node (legacy undrop_statement, not in the DAG)",
	"bigquery/ddl.md#45":      "OUT-OF-SCOPE: DROP EXTERNAL TABLE — no migration node (DROP external-object variant, ZetaSQL long-tail)",
	"spanner/ddl.md#37":       "OUT-OF-SCOPE: CREATE MODEL … INPUT/OUTPUT — ML model DDL, no migration node",
	"spanner/ddl.md#38":       "OUT-OF-SCOPE: ALTER MODEL … SET OPTIONS — ML model DDL, no migration node",
	"spanner/ddl.md#39":       "OUT-OF-SCOPE: DROP MODEL — ML model DDL, no migration node",

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

	// =========================== RESIDUAL GAP (4) ==========================
	// Genuine union-grammar gaps: the live oracle ACCEPTS the form, omni rejects.
	// Mirrored as flagged divergences in the v2 migration store; the to-do list
	// for a future grammar fix (out of corpus-closure's writes scope).
	// (11 RESIDUAL GAP total = these 4 + 6 pipe/FROM-first/MATCH_RECOGNIZE + 1 scripting Split.)
	"spanner/expressions.md#33": "RESIDUAL GAP: quantified comparison `expr {>|=|…} {ALL|ANY|SOME} (subquery)` — live Spanner oracle ACCEPTS `x > ALL (SELECT …)` / `x = ANY (SELECT …)` (Table not found, i.e. parsed), omni rejects (\"syntax error at or near ALL\"). Expression-grammar gap.",
	"spanner/expressions.md#6":  "RESIDUAL GAP: parameterized vector-array constructor `ARRAY<FLOAT32>(vector_length => N)` — live Spanner oracle ACCEPTS `SELECT ARRAY<FLOAT32>(vector_length => 128)`, omni rejects (\"syntax error at or near =>\"). Named args in a typed ARRAY constructor not parsed (plain named args f(a=>1) DO parse). TOKENIZE_NGRAMS(... =>) in the same block parses.",
	"spanner/ddl.md#51":         "RESIDUAL GAP: parameterized vector-array column type `Embedding ARRAY<FLOAT32>(vector_length=>128)` — live Spanner oracle ACCEPTS the CREATE TABLE, omni rejects (\"expected a type parameter\"). The CREATE VECTOR INDEX in the same block parses on omni.",
	"spanner/ddl.md#48":         "RESIDUAL GAP: `ALTER SEARCH INDEX … ADD/DROP STORED COLUMN` — live Spanner oracle ACCEPTS (semantic, parsed), omni rejects (\"ALTER statement parsing is not yet supported\" for SEARCH INDEX). ALTER-search-index action not wired.",

	// ===== RESIDUAL GAP / pipe & FROM-first & MATCH_RECOGNIZE (BigQuery) (6) =====
	// BigQuery-only query forms; the Spanner oracle is non-authoritative, so
	// these are triangulated against the legacy .g4 + truth1. The legacy grammar
	// leaves pipe `|>` lexed-but-unparsed (documented TODO gap) and only
	// error-marks FROM-first queries; MATCH_RECOGNIZE is in the .g4 but owned by
	// no migration node. All are P1 long-tail, not consumed by bytebase today.
	"bigquery/query.md#15": "RESIDUAL GAP: QUERY-… MATCH_RECOGNIZE clause — present in legacy .g4 (match_recognize_clause) but no migration node implements it; P1 long-tail. omni rejects (\"syntax error at or near MATCH_RECOGNIZE\").",
	"bigquery/query.md#16": "RESIDUAL GAP: QUERY-… FROM-first query (`FROM A JOIN B …` with no leading SELECT) — legacy .g4 only error-marks this (query_without_pipe_operators / bad_keyword_after_from_query); omni rejects (\"a query must begin with SELECT\"). FROM-first / pipe syntax not implemented (P1).",
	"bigquery/query.md#17": "RESIDUAL GAP: QUERY-… FROM-first query (FROM A CROSS JOIN B …) — same pipe/FROM-first gap as #16.",
	"bigquery/query.md#18": "RESIDUAL GAP: QUERY-… FROM-first query (FROM A FULL OUTER JOIN B …) — same pipe/FROM-first gap as #16.",
	"bigquery/query.md#19": "RESIDUAL GAP: QUERY-… FROM-first query (FROM A LEFT OUTER JOIN B …) — same pipe/FROM-first gap as #16.",
	"bigquery/query.md#20": "RESIDUAL GAP: QUERY-… FROM-first query (FROM A RIGHT OUTER JOIN B …) — same pipe/FROM-first gap as #16.",

	// ============== RESIDUAL GAP: scripting Split edge case (1) ==============
	"bigquery/scripting.md#15": "RESIDUAL GAP: a procedural `BEGIN … END;` block that contains a nested `BEGIN TRANSACTION;` AND a trailing terminator — the block-aware Split keeps the ';' on the block segment and parseSingle then rejects it (\"syntax error at or near ';'\"). Verified narrow: `BEGIN SELECT 1; END;` parses (Split strips the ';'); the nested transaction-statement defeats the block-depth tracking. Split/scripting edge case, not an expression gap.",
}
