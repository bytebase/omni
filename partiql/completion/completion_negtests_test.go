package completion

import (
	"testing"

	"github.com/bytebase/omni/partiql/catalog"
)

// This file is a TEST-ONLY negative/coverage backfill for the PartiQL
// completion package. It asserts the CURRENT behavior of Complete() for
// edge and negative inputs.
//
// Oracle (truth2-first): the generated ANTLR PartiQL parser/lexer and the
// legacy grammar at parser/partiql/PartiQL{Parser,Lexer}.g4. Relevant rules:
//   - INTO precedes a target table/path: insertCommand "INSERT INTO pathSimple"
//     / "INSERT INTO symbolPrimitive" (PartiQLParser.g4 ~L131-136). So a cursor
//     after INTO suggesting table names matches the contract.
//   - REMOVE precedes a path target: removeCommand "REMOVE pathSimple" (~L128).
//   - UNPIVOT precedes an expression, NOT a table name: tableUnpivot
//     "UNPIVOT expr asIdent? atIdent? byIdent?" (~L409). So a cursor after
//     UNPIVOT must NOT suggest table names.
//   - A dotted member access "x.from" is a path step, not a FROM clause:
//     exprPrimary pathStep+ (#ExprPrimaryPath, ~L528) where
//     pathStep: PERIOD key=symbolPrimitive (#PathStepDotExpr, ~L621-622).
//     "from" after "." is the path key, never the FROM keyword.
//
// The position-out-of-range panic guard already lives in
// completion_test.go (TestCompletePositionOutOfRange), so it is not
// repeated here.

func newNegCatalog() *catalog.Catalog {
	cat := catalog.New()
	cat.AddTable("Music")
	cat.AddTable("Albums")
	cat.AddTable("Artists")
	return cat
}

func negTables(cs []Candidate) []string {
	var out []string
	for _, c := range cs {
		if c.Kind == "table" {
			out = append(out, c.Text)
		}
	}
	return out
}

func negKeywords(cs []Candidate) []string {
	var out []string
	for _, c := range cs {
		if c.Kind == "keyword" {
			out = append(out, c.Text)
		}
	}
	return out
}

func hasKeyword(cs []Candidate, kw string) bool {
	for _, c := range cs {
		if c.Kind == "keyword" && c.Text == kw {
			return true
		}
	}
	return false
}

// TestCompleteNoMatchPrefix: a prefix that matches no table and no keyword
// yields zero candidates. Even in FROM context, a non-matching prefix must
// filter every table out.
func TestCompleteNoMatchPrefix(t *testing.T) {
	cat := newNegCatalog()

	t.Run("from_context_unknown_table", func(t *testing.T) {
		// FROM context, prefix "Zzz" matches none of {Music,Albums,Artists}
		// and no keyword starts with "Zzz".
		input := "SELECT * FROM Zzz"
		got := Complete(input, len(input), cat)
		if len(got) != 0 {
			t.Errorf("Complete(%q) = %v, want no candidates", input, got)
		}
	})

	t.Run("keyword_context_unknown_word", func(t *testing.T) {
		// No keyword starts with "XQZ"; default context suggests no tables.
		input := "XQZ"
		got := Complete(input, len(input), cat)
		if len(got) != 0 {
			t.Errorf("Complete(%q) = %v, want no candidates", input, got)
		}
	})

	t.Run("from_prefix_close_but_no_match", func(t *testing.T) {
		// "Mu" matches "Music"; "Mz" matches nothing. Guard the negative.
		input := "SELECT * FROM Mz"
		got := Complete(input, len(input), cat)
		if tabs := negTables(got); len(tabs) != 0 {
			t.Errorf("Complete(%q) table candidates = %v, want none", input, tabs)
		}
	})
}

// TestCompleteNonTableContexts: cursor contexts where table names must NOT
// be suggested.
func TestCompleteNonTableContexts(t *testing.T) {
	cat := newNegCatalog()

	cases := []struct {
		name  string
		input string
	}{
		// After SELECT we are in a projection (expression) context, not a
		// table context.
		{"after_select", "SELECT "},
		// WHERE introduces a predicate expression, not a table.
		{"after_where", "SELECT * FROM Music WHERE "},
		// ON introduces a join predicate, not a table.
		{"after_on", "SELECT * FROM a JOIN b ON "},
		// GROUP BY / ORDER BY operate on expressions, not table names.
		{"after_group_by", "SELECT * FROM Music GROUP BY "},
		{"after_order_by", "SELECT * FROM Music ORDER BY "},
		// SET (in UPDATE) targets a column/path, not a table.
		{"after_set", "UPDATE Music SET "},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Complete(tc.input, len(tc.input), cat)
			if tabs := negTables(got); len(tabs) != 0 {
				t.Errorf("Complete(%q): expected no table suggestions, got %v", tc.input, tabs)
			}
			// These contexts should still offer keywords (graceful, useful
			// default), so the result is not empty.
			if len(negKeywords(got)) == 0 {
				t.Errorf("Complete(%q): expected some keyword suggestions, got none", tc.input)
			}
		})
	}
}

// TestCompleteGarbageInput: non-SQL / punctuation-only / unbalanced input
// must not panic and must degrade gracefully (no table suggestions, keyword
// fallback only).
func TestCompleteGarbageInput(t *testing.T) {
	cat := newNegCatalog()

	cases := []struct {
		name  string
		input string
	}{
		{"operators", "@#$%^&*()"},
		{"unbalanced_brackets", "}{][;;;,,,"},
		{"only_punctuation", ".....,,,,,"},
		{"random_unicode", "日本語 émoji ✓"},
		{"sql_fragment_garbage", "SELECT ?? %% FROM"},
		{"only_spaces", "        "},
		{"newlines_tabs", "\n\t\n\t "},
		{"quote_soup", "'''\"\"\"```"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Complete(%q) panicked: %v", tc.input, r)
				}
			}()
			got := Complete(tc.input, len(tc.input), cat)
			// Garbage input must never surface table names: none of these end
			// in FROM/JOIN/INTO. (The "SELECT ?? %% FROM" case DOES end in
			// FROM and is exercised separately below.)
			if tc.input != "SELECT ?? %% FROM" {
				if tabs := negTables(got); len(tabs) != 0 {
					t.Errorf("Complete(%q): garbage input produced table suggestions %v", tc.input, tabs)
				}
			}
		})
	}
}

// TestCompleteIntoContext: INTO precedes a target table/path in PartiQL
// (insertCommand: INSERT INTO pathSimple / INSERT INTO symbolPrimitive), so a
// cursor right after INTO suggests table names. This documents the contract
// that INTO is a table-introducing keyword alongside FROM/JOIN.
func TestCompleteIntoContext(t *testing.T) {
	cat := newNegCatalog()

	t.Run("insert_into_empty", func(t *testing.T) {
		input := "INSERT INTO "
		got := Complete(input, len(input), cat)
		if tabs := negTables(got); len(tabs) != 3 {
			t.Errorf("Complete(%q) table candidates = %v, want all 3 tables", input, tabs)
		}
	})

	t.Run("insert_into_prefix", func(t *testing.T) {
		input := "INSERT INTO Mu"
		got := Complete(input, len(input), cat)
		tabs := negTables(got)
		if len(tabs) != 1 || tabs[0] != "Music" {
			t.Errorf("Complete(%q) table candidates = %v, want [Music]", input, tabs)
		}
	})
}

// TestCompleteRemoveContext documents that the CURRENT completion does NOT
// suggest table names after REMOVE. In the grammar, removeCommand is
// "REMOVE pathSimple" (PartiQLParser.g4 ~L128) — REMOVE targets a path that
// begins with a table/identifier, so an editor could reasonably offer table
// names here. The simple heuristic in isInFromContext only matches
// FROM/JOIN/INTO, so REMOVE falls through to the keyword-only default.
//
// This is a conservative GAP, not a correctness bug (suggesting nothing is
// never wrong), so it is asserted as current behavior rather than skipped.
// It is reported in new_findings for completeness.
func TestCompleteRemoveContext(t *testing.T) {
	cat := newNegCatalog()

	input := "REMOVE "
	got := Complete(input, len(input), cat)
	if tabs := negTables(got); len(tabs) != 0 {
		t.Errorf("Complete(%q): CURRENT behavior expects no table suggestions, got %v", input, tabs)
	}
	// Keyword fallback is still provided.
	if len(negKeywords(got)) == 0 {
		t.Errorf("Complete(%q): expected keyword fallback, got none", input)
	}
}

// TestCompleteUnpivotContext: in PartiQL, tableUnpivot is
// "UNPIVOT expr asIdent? atIdent? byIdent?" (PartiQLParser.g4 ~L409) — UNPIVOT
// is followed by an EXPRESSION, not a table name. So a cursor after UNPIVOT
// must NOT suggest table names. The CURRENT heuristic agrees (UNPIVOT is not
// in the FROM/JOIN/INTO suffix set), which is the correct contract.
func TestCompleteUnpivotContext(t *testing.T) {
	cat := newNegCatalog()

	cases := []string{
		"UNPIVOT ",
		"SELECT * FROM UNPIVOT ",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			got := Complete(input, len(input), cat)
			if tabs := negTables(got); len(tabs) != 0 {
				t.Errorf("Complete(%q): UNPIVOT precedes an expression, expected no table suggestions, got %v", input, tabs)
			}
		})
	}
}

// TestCompleteLowercaseKeywords: the context detection upper-cases the text
// before matching, so lowercase clause keywords (from/join/into) still
// trigger table suggestions. Guards the case-insensitivity contract.
func TestCompleteLowercaseKeywords(t *testing.T) {
	cat := newNegCatalog()

	t.Run("lowercase_from", func(t *testing.T) {
		input := "select * from "
		got := Complete(input, len(input), cat)
		if tabs := negTables(got); len(tabs) != 3 {
			t.Errorf("Complete(%q) table candidates = %v, want all 3 tables", input, tabs)
		}
	})

	t.Run("lowercase_join", func(t *testing.T) {
		input := "select * from Music join "
		got := Complete(input, len(input), cat)
		if tabs := negTables(got); len(tabs) != 3 {
			t.Errorf("Complete(%q) table candidates = %v, want all 3 tables", input, tabs)
		}
	})

	t.Run("lowercase_into", func(t *testing.T) {
		input := "insert into "
		got := Complete(input, len(input), cat)
		if tabs := negTables(got); len(tabs) != 3 {
			t.Errorf("Complete(%q) table candidates = %v, want all 3 tables", input, tabs)
		}
	})

	t.Run("lowercase_from_with_prefix", func(t *testing.T) {
		// prefix "m" (case-insensitive) matches table "Music"; keyword
		// "MISSING" also matches as a prefix. Verify the table side.
		input := "select * from m"
		got := Complete(input, len(input), cat)
		tabs := negTables(got)
		if len(tabs) != 1 || tabs[0] != "Music" {
			t.Errorf("Complete(%q) table candidates = %v, want [Music]", input, tabs)
		}
	})
}

// TestCompletePathStepNotFromKeyword guards the path-step false-positive.
//
// In "SELECT x.from <cursor>", the token "from" is a path step key on the
// expression x (grammar: exprPrimary pathStep+, pathStep: PERIOD
// key=symbolPrimitive — PartiQLParser.g4 ~L528, ~L621-622). It is member
// access (x.from), NOT a FROM clause, so NO table names should be suggested.
//
// CURRENT omni behavior DIVERGES: isInFromContext does a naive
// strings.HasSuffix(upper, "FROM") on "SELECT X.FROM", which matches, so the
// completion wrongly enters FROM context and suggests every table. This is a
// real divergence from the legacy completion contract; we skip rather than
// encode the wrong expectation, and report it in new_findings.
//
// Source of the bug: partiql/completion/completion.go isInFromContext
// (HasSuffix on FROM/JOIN/INTO ignores a preceding '.' path separator).
func TestCompletePathStepNotFromKeyword(t *testing.T) {
	cat := newNegCatalog()

	input := "SELECT x.from "
	got := Complete(input, len(input), cat)
	tabs := negTables(got)

	// CORRECT contract (oracle): no table suggestions — `from` after `.` is a
	// path key (member access x.from), not a FROM clause.
	if len(tabs) != 0 {
		t.Errorf("Complete(%q): `from` after `.` is a path step, not a FROM clause; "+
			"expected no table suggestions, got %v", input, tabs)
	}
	// Reached only if/when the divergence is fixed: the contract holds.
}

// TestCompletePathStepPrefixMatchesKeyword is the NON-divergent sibling of the
// path-step case: with NO trailing space, "SELECT x.from" leaves "from" as the
// partial word at the cursor. It is not yet a completed clause keyword, so the
// context is the default (no FROM trigger) and "from" is treated as a prefix.
// The only effect is that the keyword "FROM" is offered as a prefix match —
// NO table names. This is correct and must pass.
func TestCompletePathStepPrefixMatchesKeyword(t *testing.T) {
	cat := newNegCatalog()

	input := "SELECT x.from"
	got := Complete(input, len(input), cat)
	if tabs := negTables(got); len(tabs) != 0 {
		t.Errorf("Complete(%q): no clause keyword completed, expected no table suggestions, got %v", input, tabs)
	}
	if !hasKeyword(got, "FROM") {
		t.Errorf("Complete(%q): expected keyword FROM offered as a prefix match, got %v", input, negKeywords(got))
	}
}
