package analysis

import (
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestEngineStatementRanges(t *testing.T) {
	if ranges := Engine.StatementRanges("// comment\n"); ranges != nil {
		t.Fatalf("StatementRanges of no statement = %#v", ranges)
	}
	ranges := Engine.StatementRanges("db.a.find({});\nfoo bar\ndb.c.find({})")
	if len(ranges) != 2 || ranges[0].End != (analysis.LSPPosition{Character: 14}) || ranges[1].Start != (analysis.LSPPosition{Line: 2}) {
		t.Fatalf("StatementRanges = %+v", ranges)
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("db.a.find({})"); diagnostics != nil {
		t.Fatalf("Diagnose of a valid script = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("db.a.find(\ndb.b.find({})\ndb.c.find(")
	if len(diagnostics) != 2 || diagnostics[1].Start != (analysis.Position{Line: 3, Column: 11}) {
		t.Fatalf("Diagnose reported %d errors: %+v", len(diagnostics), diagnostics)
	}
	for _, d := range diagnostics {
		if d.Start.Line < 1 || d.Start.Column < 1 || d.End != (analysis.Position{}) || d.Severity != analysis.SeverityError || d.Message == "" {
			t.Fatalf("diagnostic = %+v", d)
		}
	}
}

func TestEngineExplainIgnoresFormat(t *testing.T) {
	for _, format := range []analysis.ExplainFormat{analysis.ExplainFormatDefault, analysis.ExplainFormatJSON, analysis.ExplainFormatYAML} {
		got, err := Engine.Explain("db.a.find({})", format)
		if err != nil || got != `db.a.find({}).explain("queryPlanner")` {
			t.Fatalf("Explain(%d) = %q, %v", format, got, err)
		}
	}
	if _, err := Engine.Explain("db.a.insertOne({})", analysis.ExplainFormatDefault); err == nil {
		t.Fatal("Explain(insertOne) succeeded")
	}
}
