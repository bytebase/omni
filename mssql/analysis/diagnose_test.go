package analysis

import "testing"

func TestDiagnose(t *testing.T) {
	if got := Diagnose("SELECT 1\nSELECT 2"); got != nil {
		t.Errorf("valid script: got %+v", got)
	}
	got := Diagnose("SELECT 1;\nSELECT FROM")
	if got == nil || got.Text != "SELECT 1;\nSELECT FROM;" || got.Offset != 17 || got.Message == "" {
		t.Errorf("syntax error: got %+v", got)
	}
	got = Diagnose("SELECT N'你好';\nSELECT FROM")
	if got == nil || got.Offset != len("SELECT N'你好';\nSELECT ") || got.Message == "" {
		t.Errorf("non-ascii syntax error: got %+v", got)
	}
}
