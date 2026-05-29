package parser

import "testing"

func TestOracleCompletionAPIStatementStarters(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		pos  int
	}{
		{name: "empty", sql: "", pos: 0},
		{name: "whitespace", sql: " \n\t", pos: 3},
		{name: "after semicolon", sql: "SELECT 1 FROM dual; ", pos: len("SELECT 1 FROM dual; ")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, tt.pos)
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, tok := range []int{SELECT, INSERT, UPDATE, DELETE, MERGE, CREATE, ALTER, DROP} {
				if !cs.HasToken(tok) {
					t.Fatalf("missing statement starter %q; got tokens=%v rules=%v", TokenName(tok), tokenNamesForTest(cs.Tokens), cs.Rules)
				}
			}
		})
	}
}

func TestOracleCompletionAPITokenizeAndTokenName(t *testing.T) {
	tokens := Tokenize("SELECT * FROM dual")
	if len(tokens) != 4 {
		t.Fatalf("Tokenize returned %d tokens, want 4: %#v", len(tokens), tokens)
	}
	if tokens[0].Type != SELECT || tokens[0].Loc != 0 || tokens[0].End != len("SELECT") {
		t.Fatalf("first token = %#v, want SELECT at [0,%d)", tokens[0], len("SELECT"))
	}
	if got := TokenName(SELECT); got != "SELECT" {
		t.Fatalf("TokenName(SELECT) = %q, want SELECT", got)
	}
	if got := TokenName(tokIDENT); got != "" {
		t.Fatalf("TokenName(tokIDENT) = %q, want empty", got)
	}
}

func TestOracleCompletionAPIPrefixRetry(t *testing.T) {
	sql := "SEL"
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	if !cs.HasToken(SELECT) {
		t.Fatalf("missing SELECT for partial keyword; got tokens=%v rules=%v", tokenNamesForTest(cs.Tokens), cs.Rules)
	}
}

func tokenNamesForTest(tokens []int) []string {
	result := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		name := TokenName(tok)
		if name == "" && tok > 0 && tok < 256 {
			name = string(rune(tok))
		}
		result = append(result, name)
	}
	return result
}
