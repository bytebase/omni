package completion

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bytebase/omni/pg/yacc"
	"gopkg.in/yaml.v3"
)

// testCase matches the YAML test oracle format.
type testCase struct {
	Input string              `yaml:"input"`
	Want  []expectedCandidate `yaml:"want"`
}

type expectedCandidate struct {
	Text       string `yaml:"text"`
	Type       string `yaml:"type"`
	Definition string `yaml:"definition"`
	Comment    string `yaml:"comment"`
}

// mockMetadata implements MetadataProvider for testing.
type mockMetadata struct{}

func (m *mockMetadata) GetSchemaNames(_ context.Context) []string {
	return []string{"public", "test"}
}

func (m *mockMetadata) GetTables(_ context.Context, schema string) []TableInfo {
	switch schema {
	case "public":
		return []TableInfo{{Name: "t1"}, {Name: "t2"}}
	case "test":
		return []TableInfo{{Name: "auto"}, {Name: "users"}}
	}
	return nil
}

func (m *mockMetadata) GetViews(_ context.Context, schema string) []string {
	switch schema {
	case "public":
		return []string{"v1"}
	}
	return nil
}

func (m *mockMetadata) GetSequences(_ context.Context, schema string) []string {
	switch schema {
	case "public":
		return []string{"seq1", "user_id_seq"}
	case "test":
		return []string{"order_id_seq"}
	}
	return nil
}

func (m *mockMetadata) GetColumns(_ context.Context, schema, table string) []ColumnInfo {
	key := schema + "." + table
	switch key {
	case "public.t1":
		return []ColumnInfo{
			{Name: "c1", Type: "int", NotNull: true, Definition: "public.t1 | int, NOT NULL"},
		}
	case "public.t2":
		return []ColumnInfo{
			{Name: "c1", Type: "int", NotNull: true, Definition: "public.t2 | int, NOT NULL"},
			{Name: "c2", Type: "int", NotNull: true, Definition: "public.t2 | int, NOT NULL"},
		}
	case "test.auto":
		return []ColumnInfo{
			{Name: "id", Type: "int", NotNull: true, Definition: "test.auto | int, NOT NULL"},
			{Name: "name", Type: "varchar", NotNull: true, Definition: "test.auto | varchar, NOT NULL"},
		}
	case "test.users":
		return []ColumnInfo{
			{Name: "user_id", Type: "int", NotNull: true, Definition: "test.users | int, NOT NULL"},
			{Name: "username", Type: "varchar", NotNull: true, Definition: "test.users | varchar, NOT NULL"},
		}
	}
	return nil
}

func TestContextInference(t *testing.T) {
	tests := []struct {
		sql  string
		want []GrammarContext
	}{
		// After FROM → should infer relation context
		{"SELECT * FROM ", []GrammarContext{CtxRelationRef}},
		// After WHERE → should infer column context
		{"SELECT * FROM t1 WHERE ", []GrammarContext{CtxColumnRef}},
		// After SELECT → should infer column context (select list position)
		{"SELECT ", []GrammarContext{CtxColumnRef}},
		// After JOIN → should infer relation context
		{"SELECT * FROM t1 JOIN ", []GrammarContext{CtxRelationRef}},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			pp := preprocess(tt.sql+"|", -1)
			state := simulateParse(pp.tokens)
			contexts := inferContexts(state.stack)
			t.Logf("SQL: %q → contexts: %v", tt.sql, contexts)

			for _, want := range tt.want {
				found := false
				for _, got := range contexts {
					if got == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected context %d in %v", want, contexts)
				}
			}
		})
	}
}

func TestSimulator(t *testing.T) {
	tests := []struct {
		sql       string
		minStates int // minimum expected stack depth
	}{
		{"SELECT 1", 2},
		{"SELECT * FROM t1 WHERE", 3},
		{"SELECT", 1},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			tokens := tokenize(tt.sql)
			var tokenInfos []tokenInfo
			for _, tok := range tokens {
				tokenInfos = append(tokenInfos, tok)
			}
			state := simulateParse(tokenInfos)
			if len(state.stack) < tt.minStates {
				t.Errorf("stack depth %d < minimum %d for %q", len(state.stack), tt.minStates, tt.sql)
			}
		})
	}
}

func TestCollector(t *testing.T) {
	// "SELECT * FROM " — should have IDENT as a valid next token
	sql := "SELECT * FROM "
	tokens := tokenize(sql)
	state := simulateParse(tokens)
	validTokens := collectValidTokens(state.stack)

	// IDENT should be valid (for table names)
	toknames := testTokNames()
	identTok := findTokenID(toknames, "IDENT")
	if !validTokens[identTok] {
		t.Error("IDENT should be valid after SELECT * FROM")
	}

	// Dump some valid tokens for debugging
	var names []string
	for tok := range validTokens {
		name := TokenNameForID(tok)
		if name != "" {
			names = append(names, name)
		}
	}
	t.Logf("Valid tokens after 'SELECT * FROM': %v", names)
}

func TestPreprocess(t *testing.T) {
	tests := []struct {
		sql       string
		wantCount int // expected number of tokens before cursor
	}{
		{"SELECT * FROM |", 3},           // SELECT, *, FROM
		{"SELECT 1; SELECT * FROM |", 3}, // only the second statement
		{"SELECT | FROM t1", 1},          // SELECT
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			pp := preprocess(tt.sql, -1)
			if len(pp.tokens) != tt.wantCount {
				var toks []string
				for _, t := range pp.tokens {
					toks = append(toks, fmt.Sprintf("%d(%s)", t.parserToken, t.tok.Str))
				}
				t.Errorf("got %d tokens %v, want %d", len(pp.tokens), toks, tt.wantCount)
			}
		})
	}
}

func TestCompletion(t *testing.T) {
	data, err := os.ReadFile("testdata/completion.yaml")
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var cases []testCase
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatalf("failed to parse test data: %v", err)
	}

	meta := &mockMetadata{}
	ctx := context.Background()

	passed := 0
	for i, tc := range cases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			candidates, err := Complete(ctx, tc.Input, -1, meta)
			if err != nil {
				t.Fatalf("completion error: %v", err)
			}

			// Filter out keywords and functions (test oracle excludes them)
			var filtered []Candidate
			for _, c := range candidates {
				if c.Type == CandidateKeyword || c.Type == CandidateFunction {
					continue
				}
				filtered = append(filtered, c)
			}

			// Compare
			got := formatCandidates(filtered)
			want := formatExpected(tc.Want)

			if got == want {
				passed++
				return
			}

			t.Errorf("input: %s\n  want:\n%s\n  got:\n%s", tc.Input, want, got)
		})
	}

	t.Logf("Passed %d/%d test cases", passed, len(cases))
}

func testTokNames() []string {
	return yacc.TokNames()
}

func formatCandidates(candidates []Candidate) string {
	var lines []string
	for _, c := range candidates {
		line := fmt.Sprintf("    - text: %s | type: %s", c.Text, c.Type)
		if c.Definition != "" {
			line += fmt.Sprintf(" | def: %s", c.Definition)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func formatExpected(expected []expectedCandidate) string {
	var lines []string
	for _, e := range expected {
		line := fmt.Sprintf("    - text: %s | type: %s", e.Text, CandidateType(e.Type))
		if e.Definition != "" {
			line += fmt.Sprintf(" | def: %s", e.Definition)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// --- Gap #1: Keyword candidates ---

func TestKeywordCompletion(t *testing.T) {
	meta := &mockMetadata{}
	ctx := context.Background()

	tests := []struct {
		sql          string
		wantKeywords []string // keywords that MUST be present
		noKeywords   []string // keywords that must NOT be present
	}{
		{
			sql:          "SELECT * FR|",
			wantKeywords: []string{"FROM"},
		},
		{
			sql:          "SELECT * FROM t1 WH|",
			wantKeywords: []string{"WHERE"},
		},
		{
			sql:          "SELECT * FROM t1 WHERE c1 = 1 OR|",
			wantKeywords: []string{"ORDER"},
		},
		{
			sql:          "SEL|",
			wantKeywords: []string{"SELECT"},
		},
		{
			// After FROM, SELECT should not be a valid keyword
			sql:        "SELECT * FROM |",
			noKeywords: []string{"WHERE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			candidates, err := Complete(ctx, tt.sql, -1, meta)
			if err != nil {
				t.Fatal(err)
			}

			kwMap := make(map[string]bool)
			for _, c := range candidates {
				if c.Type == CandidateKeyword {
					kwMap[c.Text] = true
				}
			}

			for _, want := range tt.wantKeywords {
				if !kwMap[want] {
					t.Errorf("expected keyword %q, got keywords: %v", want, mapKeys(kwMap))
				}
			}
			for _, no := range tt.noKeywords {
				if kwMap[no] {
					t.Errorf("unexpected keyword %q present", no)
				}
			}
		})
	}
}

// --- Gap #2: Function candidates ---

func TestFunctionCompletion(t *testing.T) {
	meta := &mockMetadata{}
	ctx := context.Background()

	// After SELECT, function names should be offered (CtxFuncName)
	candidates, err := Complete(ctx, "SELECT |", -1, meta)
	if err != nil {
		t.Fatal(err)
	}

	funcMap := make(map[string]bool)
	for _, c := range candidates {
		if c.Type == CandidateFunction {
			funcMap[c.Text] = true
		}
	}

	wantFuncs := []string{"count", "max", "min", "avg", "sum", "now", "coalesce"}
	for _, fn := range wantFuncs {
		if !funcMap[fn] {
			t.Errorf("expected function %q in SELECT list", fn)
		}
	}
}

// --- Gap #3: Type candidates ---

func TestTypeCompletion(t *testing.T) {
	meta := &mockMetadata{}
	ctx := context.Background()

	// After column name in CREATE TABLE, type names should appear
	candidates, err := Complete(ctx, "CREATE TABLE t1 (id |", -1, meta)
	if err != nil {
		t.Fatal(err)
	}

	kwMap := make(map[string]bool)
	for _, c := range candidates {
		kwMap[strings.ToLower(c.Text)] = true
	}

	// These should be present as keywords (type names are grammar keywords)
	wantTypes := []string{"integer", "varchar", "text", "boolean", "timestamp"}
	for _, tp := range wantTypes {
		if !kwMap[tp] {
			t.Errorf("expected type %q after column name in CREATE TABLE", tp)
		}
	}
}

// --- Gap #4: Grammar context replaces heuristic ---

func TestNoHeuristicFallback(t *testing.T) {
	// Verify that grammar context is always inferred for all 21 test cases
	// (i.e., the isInFromContext heuristic is never needed)
	data, err := os.ReadFile("testdata/completion.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var cases []testCase
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}

	for i, tc := range cases {
		cleanSQL := tc.Input
		actualOffset := -1
		if idx := findCursorMarker(cleanSQL); idx >= 0 {
			cleanSQL = cleanSQL[:idx] + cleanSQL[idx+1:]
			actualOffset = idx
		}
		if actualOffset < 0 {
			actualOffset = len(cleanSQL)
		}
		// Qualified (dot) cases don't need context
		q := extractQualifier(cleanSQL, actualOffset)
		if q.hasDot {
			continue
		}

		pp := preprocessFromClean(cleanSQL, actualOffset)
		state := simulateParse(pp.tokens)
		contexts := inferContexts(state.stack)

		hasAny := false
		for _, c := range contexts {
			if c == CtxColumnRef || c == CtxRelationRef {
				hasAny = true
				break
			}
		}
		if !hasAny {
			t.Errorf("case_%d (%s): no CtxColumnRef or CtxRelationRef inferred, grammar context insufficient", i, tc.Input)
		}
	}
}

// --- Gap #5: DDL/DML tests ---

func TestDDLDMLCompletion(t *testing.T) {
	meta := &mockMetadata{}
	ctx := context.Background()

	tests := []struct {
		name         string
		sql          string
		wantTypes    []CandidateType // at least one candidate of these types should exist
		wantKeywords []string        // specific keywords expected
		wantTexts    []string        // specific candidate texts expected (any type)
	}{
		{
			name:         "INSERT INTO table",
			sql:          "INSERT INTO |",
			wantTypes:    []CandidateType{CandidateTable, CandidateSchema},
			wantKeywords: []string{},
		},
		{
			name:      "INSERT INTO columns",
			sql:       "INSERT INTO t1 (|",
			wantTexts: []string{"c1"},
		},
		{
			name:         "UPDATE table",
			sql:          "UPDATE |",
			wantTypes:    []CandidateType{CandidateTable},
			wantKeywords: []string{},
		},
		{
			name:      "UPDATE SET column",
			sql:       "UPDATE t1 SET |",
			wantTexts: []string{"c1"},
		},
		{
			name:      "DELETE FROM table",
			sql:       "DELETE FROM |",
			wantTypes: []CandidateType{CandidateTable, CandidateSchema},
		},
		{
			name:      "DELETE WHERE column",
			sql:       "DELETE FROM t1 WHERE |",
			wantTexts: []string{"c1"},
		},
		{
			name:         "CREATE TABLE keyword",
			sql:          "CREATE |",
			wantKeywords: []string{"TABLE", "INDEX", "VIEW"},
		},
		{
			name:         "ALTER TABLE keyword",
			sql:          "ALTER |",
			wantKeywords: []string{"TABLE"},
		},
		{
			name:         "DROP keyword",
			sql:          "DROP |",
			wantKeywords: []string{"TABLE", "INDEX", "VIEW"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates, err := Complete(ctx, tt.sql, -1, meta)
			if err != nil {
				t.Fatalf("completion error: %v", err)
			}

			typeSet := make(map[CandidateType]bool)
			textSet := make(map[string]bool)
			kwSet := make(map[string]bool)
			for _, c := range candidates {
				typeSet[c.Type] = true
				textSet[strings.ToLower(c.Text)] = true
				if c.Type == CandidateKeyword {
					kwSet[c.Text] = true
				}
			}

			for _, wt := range tt.wantTypes {
				if !typeSet[wt] {
					t.Errorf("expected candidate type %s", wt)
				}
			}
			for _, wk := range tt.wantKeywords {
				if !kwSet[wk] {
					t.Errorf("expected keyword %q, got: %v", wk, mapKeys(kwSet))
				}
			}
			for _, wt := range tt.wantTexts {
				if !textSet[strings.ToLower(wt)] {
					t.Errorf("expected text %q in candidates", wt)
				}
			}
		})
	}
}

func mapKeys(m map[string]bool) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
