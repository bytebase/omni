package parser

import (
	"context"
	"path/filepath"
	"slices"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"
)

func TestExpectedKeywords(t *testing.T) {
	for _, tc := range []struct {
		prefix  string
		has     []string
		hasNone []string
	}{
		{prefix: "", has: []string{"SELECT", "WITH", "INSERT", "CREATE", "ALTER", "DROP", "SYSTEM", "SHOW", "EXPLAIN"}, hasNone: []string{"WHERE", "JOIN"}},
		{prefix: "SELECT ", has: []string{"DISTINCT", "ALL", "TOP", "NOT", "CASE", "NULL"}, hasNone: []string{"FROM", "WHERE", "SELECT"}},
		{prefix: "SELECT * ", has: []string{"FROM", "WHERE", "GROUP", "ORDER", "LIMIT", "SETTINGS", "FORMAT", "UNION", "INTO"}, hasNone: []string{"JOIN", "ON"}},
		{prefix: "SELECT * FROM t ", has: []string{"FINAL", "SAMPLE", "PREWHERE", "WHERE", "AS", "JOIN", "LEFT", "INNER", "CROSS", "ARRAY", "GLOBAL", "GROUP", "ORDER"}, hasNone: []string{"SELECT", "ON", "BY"}},
		{prefix: "SELECT * FROM t LEFT ", has: []string{"JOIN", "OUTER", "SEMI", "ANTI", "ANY", "ARRAY"}, hasNone: []string{"WHERE"}},
		{prefix: "SELECT a FROM t WHERE a ", has: []string{"AND", "OR", "IN", "NOT", "LIKE", "ILIKE", "BETWEEN", "IS", "GROUP"}},
		{prefix: "SELECT a FROM t GROUP ", has: []string{"BY"}, hasNone: []string{"WHERE"}},
		{prefix: "SELECT * FROM t ORDER BY a ", has: []string{"ASC", "DESC", "NULLS", "COLLATE", "LIMIT"}},
		{prefix: "SELECT count() OVER (ORDER BY a ", has: []string{"ROWS", "RANGE", "GROUPS"}},
		{prefix: "CREATE ", has: []string{"TABLE", "VIEW", "MATERIALIZED", "DATABASE", "DICTIONARY", "OR"}},
		{prefix: "CREATE TABLE t (a UInt64) ", has: []string{"ENGINE", "AS", "COMMENT", "ORDER", "PRIMARY", "PARTITION", "TTL", "SETTINGS"}},
		{prefix: "ALTER TABLE t ", has: []string{"ADD", "DROP", "MODIFY", "UPDATE", "DELETE", "ON"}},
		{prefix: "INSERT INTO t ", has: []string{"VALUES", "FORMAT", "SELECT", "SETTINGS"}},
		{prefix: "GRANT SELECT ON *.* TO u WITH GRANT ", has: []string{"OPTION"}},
		{prefix: "SYSTEM ", has: []string{"FLUSH", "RELOAD", "START", "STOP"}},
		// No keyword is a name's continuation, and nothing follows a bad
		// prefix or one ending inside a word, a string or a comment.
		{prefix: "SELECT * FROM db."},
		{prefix: "SELECT * FROM t WHERE )"},
		{prefix: "SELECT * FROM t WH"},
		{prefix: "SELECT 'open "},
		{prefix: "SELECT 1 -- note "},
		{prefix: "SELECT 1 /* open "},
	} {
		t.Run(tc.prefix, func(t *testing.T) {
			got := ExpectedKeywords(context.Background(), tc.prefix)
			if !slices.IsSorted(got) {
				t.Errorf("not sorted: %v", got)
			}
			if tc.has == nil && len(got) > 0 {
				t.Errorf("ExpectedKeywords(%q) = %v, want none", tc.prefix, got)
			}
			for _, kw := range tc.has {
				if !slices.Contains(got, kw) {
					t.Errorf("ExpectedKeywords(%q) lacks %s: %v", tc.prefix, kw, got)
				}
			}
			for _, kw := range tc.hasNone {
				if slices.Contains(got, kw) {
					t.Errorf("ExpectedKeywords(%q) has %s: %v", tc.prefix, kw, got)
				}
			}
		})
	}
}

func TestExpectedKeywordsStopsWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if got := ExpectedKeywords(ctx, "SELECT * FROM t "); len(got) > 0 {
		t.Errorf("ExpectedKeywords with a done context = %v, want none", got)
	}
}

// expectedKeywordGaps are keywords ExpectedKeywords misses where ClickHouse
// takes them, each with why. TestExpectedKeywordsFollowEngine fails an
// entry that no longer misses.
var expectedKeywordGaps = map[string]string{
	// The parser takes PROJECTION in a column list only before a name, (
	// and SELECT or WITH; probing tries at most two tokens after a keyword.
	"PROJECTION": "needs three followers",
}

// TestExpectedKeywordsFollowEngine cuts every statement ClickHouse accepts
// before each keyword the parser reads as a keyword, one the statement fails
// without, and requires ExpectedKeywords to offer it there.
func TestExpectedKeywordsFollowEngine(t *testing.T) {
	f, err := compat.ReadFile(filepath.Join("testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	missedGap := map[string]bool{}
	for _, row := range f.Rows {
		if row.Expected != compat.Accept || row.Omni != compat.Accept {
			continue
		}
		toks, err := Tokenize(row.Statement)
		if err != nil {
			continue
		}
		for _, tok := range toks {
			if !tok.IsKeyword() {
				continue
			}
			if _, err := Parse(row.Statement[:tok.Pos] + "omni_name" + row.Statement[tok.End:]); err == nil {
				// The keyword may be read as a name there.
				continue
			}
			checked++
			prefix, kw := row.Statement[:tok.Pos], keywordOf(tok)
			if got := ExpectedKeywords(context.Background(), prefix); slices.Contains(got, kw) {
				continue
			}
			if _, ok := expectedKeywordGaps[kw]; ok {
				missedGap[kw] = true
				continue
			}
			t.Errorf("ExpectedKeywords(%q) lacks %s, which ClickHouse accepts in %q", prefix, kw, row.Statement)
		}
	}
	if checked < 1000 {
		t.Errorf("checked %d keywords, want at least 1000", checked)
	}
	for kw, why := range expectedKeywordGaps {
		if !missedGap[kw] {
			t.Errorf("ExpectedKeywords no longer misses %s (%s): remove the gap", kw, why)
		}
	}
}

// keywordOf returns the keyword table's spelling of a keyword token.
func keywordOf(tok Token) string {
	for word, tt := range keywords {
		if tt == tok.Type && len(word) == len(tok.Raw) {
			return word
		}
	}
	return ""
}
