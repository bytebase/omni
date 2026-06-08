package redshift

import "testing"

func TestDiagnose(t *testing.T) {
	if diagnostics := Diagnose("CREATE TABLE t (id INT) DISTSTYLE EVEN;"); len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics, got %#v", diagnostics)
	}

	diagnostics := Diagnose("SELECTT FROM")
	if len(diagnostics) != 1 {
		t.Fatalf("expected one diagnostic, got %d", len(diagnostics))
	}
	if diagnostics[0].Message == "" {
		t.Fatalf("expected diagnostic message")
	}
	if diagnostics[0].Range.Start.Line != 0 || diagnostics[0].Range.Start.Character != 0 {
		t.Fatalf("expected diagnostic at start, got %#v", diagnostics[0].Range.Start)
	}
}
