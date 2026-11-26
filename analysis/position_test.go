package analysis

import "testing"

func TestPositionAt(t *testing.T) {
	script := "ab\r\nc\rdé😀x\nz"
	for _, tc := range []struct {
		offset int
		want   Position
	}{
		{-1, Position{1, 1}},
		{0, Position{1, 1}},
		{2, Position{1, 3}},
		{3, Position{2, 1}}, // between \r and \n
		{4, Position{2, 1}},
		{6, Position{3, 1}},
		{8, Position{3, 2}},  // inside é
		{9, Position{3, 3}},  // after é
		{11, Position{3, 3}}, // inside 😀
		{13, Position{3, 4}}, // after 😀
		{len(script), Position{4, 2}},
		{len(script) + 5, Position{4, 2}},
	} {
		if got := PositionAt(script, tc.offset); got != tc.want {
			t.Errorf("PositionAt(%d) = %+v, want %+v", tc.offset, got, tc.want)
		}
	}
	if got := (Position{1, 3}).Shift(Position{4, 5}); got != (Position{4, 7}) {
		t.Errorf("Shift first line = %+v", got)
	}
	if got := (Position{2, 3}).Shift(Position{4, 5}); got != (Position{5, 3}) {
		t.Errorf("Shift later line = %+v", got)
	}
}

func TestLSPRanges(t *testing.T) {
	script := "a😀\nbc"
	got := LSPRanges(script, []ByteRange{{0, 5}, {6, 8}, {1, 3}, {5, 2}, {0, 99}})
	want := []LSPRange{
		{Start: LSPPosition{0, 0}, End: LSPPosition{0, 3}},
		{Start: LSPPosition{1, 0}, End: LSPPosition{1, 2}},
		{Start: LSPPosition{0, 1}, End: LSPPosition{}},
	}
	if len(got) != len(want) {
		t.Fatalf("LSPRanges = %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("range %d = %+v, want %+v", i, got[i], want[i])
		}
	}
	if LSPRanges(script, nil) != nil {
		t.Error("nil ranges should stay nil")
	}
}

func TestTrimStatement(t *testing.T) {
	if got := TrimStatement("\n  SELECT 1 ;; \n"); got != "SELECT 1" {
		t.Errorf("TrimStatement = %q", got)
	}
}
