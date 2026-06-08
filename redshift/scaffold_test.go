package redshift

import "testing"

func TestParsePostgresCompatibleSelect(t *testing.T) {
	stmts, err := Parse("SELECT 1")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Parse returned %d statements, want 1", len(stmts))
	}
	if stmts[0].Text != "SELECT 1" {
		t.Fatalf("statement text = %q, want %q", stmts[0].Text, "SELECT 1")
	}
}
