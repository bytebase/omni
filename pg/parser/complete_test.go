package parser

import "testing"

func TestCollectStmtKeywords(t *testing.T) {
	cs := Collect("", 0)
	if cs == nil {
		t.Fatal("Collect returned nil")
	}

	// At the start of a statement, the parser should offer statement-starting keywords.
	want := []int{SELECT, INSERT, CREATE, ALTER, DROP, UPDATE, DELETE_P, WITH,
		SET, SHOW, GRANT, REVOKE, TRUNCATE, BEGIN_P, COMMIT, ROLLBACK}

	for _, tok := range want {
		if !cs.HasToken(tok) {
			t.Errorf("expected token %d in candidates, but not found", tok)
		}
	}
}

func TestCollectBasic(t *testing.T) {
	cs := newCandidateSet()

	cs.addToken(SELECT)
	cs.addToken(INSERT)
	cs.addToken(SELECT) // duplicate

	if len(cs.Tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(cs.Tokens))
	}
	if !cs.HasToken(SELECT) {
		t.Error("expected HasToken(SELECT) = true")
	}
	if cs.HasToken(UPDATE) {
		t.Error("expected HasToken(UPDATE) = false")
	}

	cs.addRule("table_ref")
	cs.addRule("table_ref") // duplicate
	cs.addRule("column_ref")

	if len(cs.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(cs.Rules))
	}
	if !cs.HasRule("table_ref") {
		t.Error("expected HasRule(table_ref) = true")
	}
	if cs.HasRule("expr") {
		t.Error("expected HasRule(expr) = false")
	}
}
