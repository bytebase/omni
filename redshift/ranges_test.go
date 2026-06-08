package redshift

import "testing"

func TestStatementRanges(t *testing.T) {
	ranges, err := StatementRanges("SELECT 1;\n  SELECT 2;")
	if err != nil {
		t.Fatalf("StatementRanges returned error: %v", err)
	}
	if len(ranges) != 2 {
		t.Fatalf("expected two ranges, got %d", len(ranges))
	}
	if ranges[0].Start != (LSPPosition{Line: 0, Character: 0}) || ranges[0].End != (LSPPosition{Line: 0, Character: 9}) {
		t.Fatalf("unexpected first range: %#v", ranges[0])
	}
	if ranges[1].Start != (LSPPosition{Line: 1, Character: 2}) || ranges[1].End != (LSPPosition{Line: 1, Character: 11}) {
		t.Fatalf("unexpected second range: %#v", ranges[1])
	}
}

func TestStatementRangesUTF16CommentPrefix(t *testing.T) {
	ranges, err := StatementRanges("/*🙂*/ SELECT 1;")
	if err != nil {
		t.Fatalf("StatementRanges returned error: %v", err)
	}
	if len(ranges) != 1 {
		t.Fatalf("expected one range, got %d", len(ranges))
	}
	if ranges[0].Start != (LSPPosition{Line: 0, Character: 7}) {
		t.Fatalf("expected UTF-16 start after comment at character 7, got %#v", ranges[0].Start)
	}
}
