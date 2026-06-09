package parser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestCorpusClosure (M2 — corpus-closure) is the whole-corpus parse-verification
// gate for the Snowflake engine. It walks EVERY .sql file under
// testdata/legacy/ (the legacy ANTLR lift) and testdata/official/ (the official
// docs scrape, one directory per Snowflake command) and, for each file, splits
// it with F3's Split and parses every segment with ParseBestEffort/parseSingle,
// asserting ZERO parse errors for every owned, complete statement.
//
// Files that exercise a genuine residual grammar gap (a clause/object the engine
// does not parse yet) or carry context-only / illustrative fragments are listed
// in corpusSkips with a precise reason. The skip-list self-prunes: a skipped
// file that starts parsing clean fails the test with a "REMOVE me" note, so the
// list cannot silently rot as the grammar grows.
//
// Two files (create-function example_01 / example_03) carry a multi-line
// single-quoted routine body. F3's Split mis-segments those because the omni
// lexer's scanString cannot span a newline, so a ';' that follows a newline
// inside the body looks like a statement terminator. They are NOT a parser gap —
// the whole file parses clean when handed to parseSingle directly — so they are
// driven whole-file (corpusWholeFile) instead of being skipped.
//
// Coverage snapshot at authoring time: 576/657 files parse clean
// (574 via Split + 2 via whole-file), 81 files skipped (categorized below).
func TestCorpusClosure(t *testing.T) {
	files := collectCorpusFiles(t)
	if len(files) != 657 {
		t.Errorf("found %d corpus .sql files, expected 657 (27 legacy + 630 official) — corpus changed; re-baseline the skip-list", len(files))
	}

	var (
		cleanSplit int // files that parse clean via Split + parseSingle
		cleanWhole int // single-quoted-body files that parse clean whole-file
		skipped    int // files matched by corpusSkips
	)

	for _, rel := range files {
		rel := rel
		t.Run(rel, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", rel))
			if err != nil {
				t.Fatalf("read %s: %v", rel, err)
			}
			src := string(data)

			// Whole-file single-statement files: Split mis-segments the embedded
			// single-quoted body, but the statement itself is valid and parses
			// clean via parseSingle on the whole text.
			if reason, ok := corpusWholeFile[rel]; ok {
				text := strings.TrimRight(strings.TrimSpace(src), ";")
				if _, errs := parseSingle(text, 0); len(errs) > 0 {
					t.Errorf("whole-file %s (%s) produced %d error(s): %v", rel, reason, len(errs), errs)
				}
				cleanWhole++
				return
			}

			// Skip-listed files: must STILL fail to parse, else the gap was
			// closed and the entry should be removed.
			if reason, ok := corpusSkips[rel]; ok {
				skipped++
				if corpusFileParsesClean(src) {
					t.Errorf("SKIP-LIST STALE: %s now parses clean — REMOVE it from corpusSkips (was skipped for: %s)", rel, reason)
				}
				return
			}

			// The contract: every segment of every non-skipped file parses with
			// zero errors.
			for _, seg := range Split(src) {
				if strings.TrimSpace(seg.Text) == "" {
					continue
				}
				if _, errs := parseSingle(seg.Text, seg.ByteStart); len(errs) > 0 {
					text := strings.TrimSpace(seg.Text)
					t.Errorf("statement produced %d error(s): %v\n  stmt: %s", len(errs), errs, text)
				}
			}
			cleanSplit++
		})
	}

	t.Logf("corpus closure: %d files total — %d clean (Split) + %d clean (whole-file) + %d skipped",
		len(files), cleanSplit, cleanWhole, skipped)
}

// corpusFileParsesClean reports whether every non-empty Split segment of src
// parses with zero errors. Used by the skip-list self-prune check so that the
// "still fails" verdict is computed with the exact path the main assertion uses.
func corpusFileParsesClean(src string) bool {
	any := false
	for _, seg := range Split(src) {
		if strings.TrimSpace(seg.Text) == "" {
			continue
		}
		any = true
		if _, errs := parseSingle(seg.Text, seg.ByteStart); len(errs) > 0 {
			return false
		}
	}
	return any
}

// collectCorpusFiles returns every .sql file under testdata/legacy and
// testdata/official as a slice of paths RELATIVE to testdata/ (e.g.
// "official/select/example_01.sql"), sorted for deterministic ordering.
func collectCorpusFiles(t *testing.T) []string {
	t.Helper()
	var out []string
	for _, dir := range []string{"legacy", "official"} {
		root := filepath.Join("testdata", dir)
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".sql") {
				return nil
			}
			rel, rerr := filepath.Rel("testdata", path)
			if rerr != nil {
				return rerr
			}
			out = append(out, filepath.ToSlash(rel))
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", root, err)
		}
	}
	sort.Strings(out)
	return out
}

// corpusWholeFile lists files that are a SINGLE valid statement whose multi-line
// single-quoted body defeats F3's Split (a Split-layer limitation tracked as a
// flagged divergence, not a parser gap). They are parsed whole-file via
// parseSingle and must parse clean.
var corpusWholeFile = map[string]string{
	"official/create-function/example_01.sql": "single CREATE FUNCTION with a multi-line single-quoted Java body (embedded ';' defeats Split)",
	"official/create-function/example_03.sql": "single CREATE FUNCTION with a multi-line single-quoted JavaScript body (embedded ';' defeats Split)",
}

// corpusSkips is the explicit, categorized skip-list. Each key is a corpus file
// (relative to testdata/) that contains at least one statement the engine cannot
// yet parse; the value is the reason. Skips fall into three buckets:
//
//   - RESIDUAL GAP — a real clause/object the grammar does not implement yet.
//     These are mirrored as flagged divergences in the v2 migration store and
//     are the to-do list for the remaining Snowflake grammar nodes.
//   - DEPENDENCY GAP — the statement is owned by a built node but trips a shared
//     sub-parser (expression / table-reference) that lacks a form. Closing the
//     dependency clears the skip with no change to the owning node.
//   - CONTEXT — the file's statements are entirely owned by an unbuilt node or
//     are illustrative-only on a docs page whose primary command lives in sibling
//     example files.
//
// Grouped by root cause. Counts in comments are at authoring time.
var corpusSkips = map[string]string{
	// --- RESIDUAL GAP: SHOW ->> result-pipe + $N table-ref (other nodes) ---
	// CREATE/ALTER WAREHOUSE now parse (gap-warehouse); this file remains skipped
	// only because its two `SHOW WAREHOUSES ... ->> SELECT ... FROM $1` statements
	// exercise the ->> result-pipe and $N table-ref gaps owned by other nodes. The
	// CREATE OR ALTER WAREHOUSE and DROP WAREHOUSE statements in it parse clean.
	"official/create-warehouse/example_07.sql": "RESIDUAL GAP: SHOW WAREHOUSES ->> SELECT ... FROM $1 — ->> result-pipe + $N table-ref owned by other nodes (CREATE OR ALTER WAREHOUSE parses)",

	// --- RESIDUAL GAP: CREATE INTERACTIVE MATERIALIZED VIEW (object node not built) ---
	// The ALTER WAREHOUSE ADD TABLES statement in this file now parses
	// (gap-warehouse); it remains skipped only for the CREATE INTERACTIVE
	// MATERIALIZED VIEW (INTERACTIVE keyword) statement owned by another node.
	"official/create-materialized-view/example_02.sql": "RESIDUAL GAP: CREATE INTERACTIVE MATERIALIZED VIEW (INTERACTIVE keyword) — object node not built (ALTER WAREHOUSE ADD TABLES parses)",

	// --- RESIDUAL GAP: CREATE/ALTER/DROP ALERT (object node not built) ---
	"legacy/alerts.sql": "RESIDUAL GAP: CREATE/ALTER/DROP ALERT — alert object node not built",

	// --- RESIDUAL GAP: DROP object variants (parser-drop covers only some objects) ---
	"legacy/drop.sql": "RESIDUAL GAP: DROP {CONNECTION,FAILOVER GROUP,INTEGRATION,MASKING POLICY,NETWORK POLICY,REPLICATION GROUP,RESOURCE MONITOR,ROW ACCESS POLICY,SESSION POLICY,SHARE,USER,...} not all built",

	// --- RESIDUAL GAP: dynamic-table clauses (named-arg => and INTERVAL/refresh-mode) ---
	"official/create-dynamic-table/example_04.sql": "RESIDUAL GAP: CLONE ... AT (TIMESTAMP => TO_TIMESTAMP_TZ(...)) named-arg => in time-travel clause",
	"official/create-dynamic-table/example_07.sql": "RESIDUAL GAP: IMMUTABLE WHERE (... - INTERVAL '1 day') refresh-mode + interval literal not parsed",

	// --- RESIDUAL GAP: named function arguments  name => value  (expr/func-call) ---
	"official/copy-into-table/example_06.sql":       "RESIDUAL GAP: TABLE(DATA_SOURCE(TYPE => 'STREAMING')) named-arg => in function call",
	"official/create-external-table/example_09.sql": "RESIDUAL GAP: INFER_SCHEMA(LOCATION=>'...', FILE_FORMAT=>'...') named-arg => in function call",
	"official/create-external-table/example_10.sql": "RESIDUAL GAP: INFER_SCHEMA(...) named-arg => in function call",
	"official/create-pipe/example_08.sql":           "RESIDUAL GAP: TABLE(DATA_SOURCE(TYPE => 'STREAMING')) named-arg => in function call",
	"official/create-pipe/example_09.sql":           "RESIDUAL GAP: named-arg => in function call (pipe streaming source)",
	"official/create-pipe/example_10.sql":           "RESIDUAL GAP: named-arg => in function call (pipe streaming source)",
	"official/join-lateral/example_05.sql":          "RESIDUAL GAP: LATERAL FLATTEN(INPUT => ..., PATH => ...) named-arg => in function call",
	"official/truncate-table/example_02.sql":        "RESIDUAL GAP: table(generator(rowcount=>20)) named-arg => in function call",

	// --- DEPENDENCY GAP: $N / $N:path result-set ref + stage path as a table ref (T5) ---
	"official/create-external-table/example_01.sql":  "DEPENDENCY GAP: SELECT metadata$filename FROM @s1/ — stage-path table ref + metadata$col not parsed",
	"official/create-external-volume/example_05.sql": "DEPENDENCY GAP: SHOW ... ->> SELECT * FROM $1 — $N result-set table ref (T5) not parsed",
	"official/create-table/example_06.sql":           "DEPENDENCY GAP: CTAS AS SELECT $1:o_custkey::number FROM @my_stage — $N:path semi-structured + stage table ref",
	"official/show-tables/example_07.sql":            "DEPENDENCY GAP: SHOW ... ->> SELECT * FROM $1 — $N result-set table ref (T5) not parsed",
	"official/show-tables/example_08.sql":            "DEPENDENCY GAP: SHOW ... ->> SELECT ... FROM $1 — $N result-set table ref (T5) not parsed",
	"official/show-tables/example_09.sql":            "DEPENDENCY GAP: SHOW ... ->> SELECT ... FROM $1 — $N result-set table ref (T5) not parsed",

	// --- RESIDUAL GAP: VALUES as a table source  FROM (VALUES (...), ...) ---
	"official/values/example_01.sql": "RESIDUAL GAP: SELECT * FROM (VALUES (...), ...) — VALUES table source not parsed",
	"official/values/example_02.sql": "RESIDUAL GAP: FROM (VALUES ...) with positional $N column refs — VALUES table source not parsed",
	"official/values/example_03.sql": "RESIDUAL GAP: FROM (VALUES ...) AS v(...) join — VALUES table source not parsed",
	"official/values/example_04.sql": "RESIDUAL GAP: FROM (VALUES ...) AS v(c1, c2) — VALUES table source + derived column list not parsed",

	// --- RESIDUAL GAP: ORDER BY ALL ---
	"official/order-by/example_01.sql": "RESIDUAL GAP: ORDER BY ALL not parsed",
	"official/order-by/example_11.sql": "RESIDUAL GAP: ORDER BY ALL not parsed",
	"official/order-by/example_12.sql": "RESIDUAL GAP: ORDER BY ALL ASC not parsed",
	"official/order-by/example_13.sql": "RESIDUAL GAP: ORDER BY ALL not parsed (ALTER SESSION setup now parses)",
	"official/order-by/example_14.sql": "RESIDUAL GAP: ORDER BY ALL NULLS FIRST not parsed",
	"official/order-by/example_15.sql": "RESIDUAL GAP: ORDER BY ALL NULLS LAST not parsed",
	"official/merge/example_09.sql":    "RESIDUAL GAP: WHEN MATCHED ... DELETE / ORDER BY ALL form not parsed",

	// --- RESIDUAL GAP: SELECT * EXCLUDE / RENAME column transforms ---
	"official/select/example_07.sql": "RESIDUAL GAP: SELECT * EXCLUDE col not parsed",
	"official/select/example_11.sql": "RESIDUAL GAP: SELECT * EXCLUDE col RENAME (...) not parsed",
	"official/select/example_16.sql": "RESIDUAL GAP: SELECT tbl.* EXCLUDE / RENAME column transforms not parsed",

	// --- RESIDUAL GAP: trailing comma in SELECT list ---
	"official/select/example_24.sql": "RESIDUAL GAP: trailing comma before FROM in SELECT list (Snowflake permits it) not parsed",

	// --- RESIDUAL GAP: (+) Oracle-style outer-join operator ---
	"official/where/example_09.sql": "RESIDUAL GAP: WHERE t1.c1 = t2.c2(+) Oracle-style outer join operator not parsed",
	"official/where/example_10.sql": "RESIDUAL GAP: WHERE ... (+) Oracle-style outer join operator not parsed",

	// --- RESIDUAL GAP: Snowflake Scripting blocks (DECLARE..BEGIN..END, CALL ... INTO) ---
	"official/execute-immediate/example_01.sql": "RESIDUAL GAP: CREATE PROCEDURE ... AS DECLARE ... BEGIN ... END scripting body (non-$$) not parsed",
	"official/execute-immediate/example_04.sql": "RESIDUAL GAP: CREATE PROCEDURE ... AS DECLARE ... BEGIN ... END scripting body (non-$$) not parsed",
	"official/call/example_07.sql":              "RESIDUAL GAP: DECLARE ... BEGIN CALL p(...) INTO :var ... END scripting block not parsed",

	// --- RESIDUAL GAP: parenthesized view body with leading WITH (CTE) ---
	"official/create-view/example_05.sql": "RESIDUAL GAP: CREATE VIEW ... AS ( WITH ... SELECT ... ) — parenthesized CTE view body not parsed",

	// --- RESIDUAL GAP / MALFORMED: legacy select.sql expression corpus ---
	// LEFT(...)/ILIKE(...) used as function names, lowercase `union` inside a
	// derived-table subquery, and to_date(SELECT ...) (a scalar subquery written
	// without parentheses — a legacy-corpus malformation). Mixed residual-gap and
	// docs-typo; tracked as one entry.
	"legacy/select.sql": "RESIDUAL GAP/MALFORMED: LEFT()/ILIKE() as function names, lowercase `union` in derived table, and to_date(SELECT...) (unparenthesized scalar subquery, a legacy malformation)",
}
