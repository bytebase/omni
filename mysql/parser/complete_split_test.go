package parser

import "testing"

func TestCollectWithDelimiter(t *testing.T) {
	// Cursor is at the end, inside the second statement after DELIMITER is restored.
	sql := "DELIMITER ;;\nSELECT 1;;\nDELIMITER ;\nSELECT "
	cursor := len(sql)

	cs := Collect(sql, cursor)
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	if len(cs.Tokens) == 0 && len(cs.Rules) == 0 {
		t.Error("expected non-empty candidates")
	}
}

func TestCollectMultiStatement(t *testing.T) {
	// Cursor is in the second statement — should get candidates.
	sql := "SELECT 1; SELECT "
	cursor := len(sql)

	cs := Collect(sql, cursor)
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	if len(cs.Tokens) == 0 && len(cs.Rules) == 0 {
		t.Error("expected non-empty candidates")
	}
}
