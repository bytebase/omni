package parsertest

import (
	"testing"

	"github.com/bytebase/omni/elasticsearch"
)

func TestDiagnose(t *testing.T) {
	tests := []struct {
		name            string
		statement       string
		wantCount       int
		wantFirstLine   int
		wantFirstColumn int
	}{
		{
			name:      "no errors — valid single request",
			statement: "GET _search",
			wantCount: 0,
		},
		{
			name:      "no errors — valid request with body",
			statement: "POST /my-index/_search\n{\"query\":{\"match_all\":{}}}",
			wantCount: 0,
		},
		{
			name:      "no errors — multiple valid requests",
			statement: "GET _cluster/health\n\nPOST /my-index/_doc\n{\"field\":\"value\"}",
			wantCount: 0,
		},
		{
			// A bare non-HTTP-verb token before a valid request triggers one
			// recoverable syntax error; the parser skips ahead to the next verb.
			name:            "single error — bad token followed by valid request",
			statement:       "3333\nGET _search",
			wantCount:       1,
			wantFirstLine:   0,
			wantFirstColumn: 2,
		},
		{
			// Bad token sandwiched between two valid requests.
			name:            "single error — bad token between valid requests",
			statement:       "GET _a\n3333\nPUT /index/_doc",
			wantCount:       1,
			wantFirstLine:   1,
			wantFirstColumn: 2,
		},
		{
			// Two separate bad tokens, each followed by a valid request, produce
			// two distinct diagnostics.
			name:            "multiple errors — two bad tokens interleaved with valid requests",
			statement:       "3333\nGET _a\n4444\nGET _b",
			wantCount:       2,
			wantFirstLine:   0,
			wantFirstColumn: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags, err := elasticsearch.Diagnose(tc.statement)
			if err != nil {
				t.Fatalf("Diagnose returned unexpected error: %v", err)
			}
			if len(diags) != tc.wantCount {
				t.Fatalf("diagnostic count: got %d, want %d\ndiags: %+v", len(diags), tc.wantCount, diags)
			}
			if tc.wantCount == 0 {
				return
			}
			first := diags[0]
			if first.Range.Start.Line != tc.wantFirstLine || first.Range.Start.Column != tc.wantFirstColumn {
				t.Errorf("first diagnostic position: got line=%d col=%d, want line=%d col=%d",
					first.Range.Start.Line, first.Range.Start.Column,
					tc.wantFirstLine, tc.wantFirstColumn)
			}
			if first.Severity != elasticsearch.SeverityError {
				t.Errorf("severity: got %d, want SeverityError", first.Severity)
			}
			if first.Message == "" {
				t.Error("first diagnostic message is empty")
			}
			// End column must be exactly one past start (non-empty range).
			if first.Range.End.Column != first.Range.Start.Column+1 {
				t.Errorf("end column: got %d, want start+1=%d", first.Range.End.Column, first.Range.Start.Column+1)
			}
			if first.Range.End.Line != first.Range.Start.Line {
				t.Errorf("end line: got %d, want same as start=%d", first.Range.End.Line, first.Range.Start.Line)
			}
		})
	}
}
