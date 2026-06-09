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
// Multi-line single-quoted routine bodies (e.g. create-function example_01 /
// example_03, a Java / JavaScript handler) used to defeat F3's Split because the
// omni lexer's scanString could not span a newline, so a ';' on a continuation
// line inside the body read as a statement terminator. Split now runs its lexer
// in multilineString mode (split.go), so those segment correctly and parse via
// the normal Split path; corpusWholeFile is consequently empty.
//
// Coverage snapshot at authoring time: 651/657 files parse clean (all via
// Split), 6 files skipped (categorized below).
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

// corpusWholeFile lists files that are a SINGLE valid statement whose body
// defeats F3's Split and so must be parsed whole-file via parseSingle.
//
// It is now EMPTY: Split runs its lexer in multilineString mode (split.go), so a
// multi-line single-quoted routine body — the Java / JavaScript handler bodies
// in create-function example_01 / example_03 that used to live here — segments
// correctly and parses end-to-end through the normal Split path. The map and its
// branch are retained as a safety hatch for any future single-statement file
// whose body genuinely cannot be segmented.
var corpusWholeFile = map[string]string{}

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
	// --- DEPENDENCY GAP: $N:path semi-structured ref + stage path as a table ref (T5) ---
	// (Bare $N as a table ref and the SHOW ... ->> ... FROM $1 result-pipe now
	// parse via gap-from-values; these two remain because they additionally need
	// a stage-path source (@s1/, @my_stage) and the $N:path semi-structured cast.)
	"official/create-external-table/example_01.sql": "DEPENDENCY GAP: SELECT metadata$filename FROM @s1/ — stage-path table ref + metadata$col not parsed",
	"official/create-table/example_06.sql":          "DEPENDENCY GAP: CTAS AS SELECT $1:o_custkey::number FROM @my_stage — $N:path semi-structured + stage table ref",

	// --- RESIDUAL GAP: Snowflake Scripting blocks (DECLARE..BEGIN..END, CALL ... INTO) ---
	// The PROCEDURE AS DECLARE..BEGIN..END body now routes to the scripting parser
	// (gap-scripting-split), so the block itself parses; these two remain only for
	// an EXECUTE IMMEDIATE argument-grammar gap inside the body, NOT a split/body
	// gap: example_01 uses `EXECUTE IMMEDIATE <var> || '...'` (a string EXPRESSION
	// argument, not a bare variable) and example_04 uses `rs := (EXECUTE IMMEDIATE
	// :bind USING (...))` (a parenthesized EXECUTE IMMEDIATE assigned to a RESULTSET
	// + a `:`-bind argument). Both are owned by the EXECUTE IMMEDIATE node.
	"official/execute-immediate/example_01.sql": "RESIDUAL GAP: EXECUTE IMMEDIATE <var> || '...' string-expression argument (block body now parses; EXECUTE IMMEDIATE expr-arg not yet supported)",
	"official/execute-immediate/example_04.sql": "RESIDUAL GAP: rs := (EXECUTE IMMEDIATE :bind USING (...)) parenthesized-statement assignment + :bind arg (block body now parses; this EXECUTE IMMEDIATE form not yet supported)",
	"official/call/example_07.sql":              "RESIDUAL GAP: DECLARE ... BEGIN CALL p(...) INTO :var ... END scripting block not parsed",

	// --- RESIDUAL GAP / MALFORMED: legacy select.sql expression corpus ---
	// LEFT(...)/ILIKE(...) used as function names, lowercase `union` inside a
	// derived-table subquery, and to_date(SELECT ...) (a scalar subquery written
	// without parentheses — a legacy-corpus malformation). Mixed residual-gap and
	// docs-typo; tracked as one entry.
	"legacy/select.sql": "RESIDUAL GAP/MALFORMED: LEFT()/ILIKE() as function names, lowercase `union` in derived table, and to_date(SELECT...) (unparenthesized scalar subquery, a legacy malformation)",
}
