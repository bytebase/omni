package parser

import "testing"

func TestLineTable_Empty(t *testing.T) {
	lt := NewLineTable("")
	line, col := lt.Position(0)
	if line != 1 || col != 1 {
		t.Errorf("Position(0) on empty input = (%d, %d), want (1, 1)", line, col)
	}
}

func TestLineTable_SingleLine(t *testing.T) {
	lt := NewLineTable("abc")
	cases := []struct {
		offset   int
		wantLine int
		wantCol  int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // 'b'
		{2, 1, 3}, // 'c'
		{3, 1, 4}, // just past 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)",
				c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_LFLineEndings(t *testing.T) {
	// "a\nb\nc"  positions: a=0 \n=1 b=2 \n=3 c=4
	lt := NewLineTable("a\nb\nc")
	cases := []struct {
		offset   int
		wantLine int
		wantCol  int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // '\n' (still line 1)
		{2, 2, 1}, // 'b'
		{3, 2, 2}, // '\n' (still line 2)
		{4, 3, 1}, // 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)",
				c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_CRLFLineEndings(t *testing.T) {
	// "a\r\nb\r\nc"  positions: a=0 \r=1 \n=2 b=3 \r=4 \n=5 c=6
	// Line 1: starts at 0, contains a, \r, \n (through offset 2)
	// Line 2: starts at 3, contains b, \r, \n (through offset 5)
	// Line 3: starts at 6, contains c
	lt := NewLineTable("a\r\nb\r\nc")
	cases := []struct {
		offset   int
		wantLine int
		wantCol  int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // '\r'
		{2, 1, 3}, // '\n' (still line 1 — CRLF is ONE break)
		{3, 2, 1}, // 'b'
		{4, 2, 2}, // '\r'
		{5, 2, 3}, // '\n'
		{6, 3, 1}, // 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)",
				c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_CROnly(t *testing.T) {
	// "a\rb\rc"  positions: a=0 \r=1 b=2 \r=3 c=4
	lt := NewLineTable("a\rb\rc")
	cases := []struct {
		offset   int
		wantLine int
		wantCol  int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // '\r' (still line 1)
		{2, 2, 1}, // 'b'
		{3, 2, 2}, // '\r'
		{4, 3, 1}, // 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)",
				c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_OffsetPastEnd(t *testing.T) {
	lt := NewLineTable("abc")
	line, col := lt.Position(100)
	// Should clamp to the last known line (line 1), with col computed
	// from the last line start.
	if line != 1 {
		t.Errorf("Position(100) line = %d, want 1", line)
	}
	// col should be at least 4 (one past 'c')
	if col < 4 {
		t.Errorf("Position(100) col = %d, want >= 4", col)
	}
}

func TestLineTable_OffsetNegative(t *testing.T) {
	lt := NewLineTable("abc")
	line, col := lt.Position(-5)
	if line != 1 || col != 1 {
		t.Errorf("Position(-5) = (%d, %d), want (1, 1)", line, col)
	}
}
