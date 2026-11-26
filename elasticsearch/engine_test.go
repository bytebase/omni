package elasticsearch

import (
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestEngineStatementRanges(t *testing.T) {
	for _, script := range []string{"", "# comment\n"} {
		if ranges := Engine.StatementRanges(script); ranges == nil || len(ranges) != 0 {
			t.Fatalf("StatementRanges(%q) = %#v", script, ranges)
		}
	}
	ranges := Engine.StatementRanges("GET /a/_search\n{}\n\nGET /😀/_search\n")
	want := analysis.LSPRange{Start: analysis.LSPPosition{Line: 3}, End: analysis.LSPPosition{Line: 3, Character: 15}}
	if len(ranges) != 2 || ranges[1] != want {
		t.Fatalf("StatementRanges = %+v", ranges)
	}
}
