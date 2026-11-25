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
	lt := NewLineTable("SELECT")
	cases := []struct {
		offset            int
		wantLine, wantCol int
	}{
		{0, 1, 1}, // 'S'
		{1, 1, 2}, // 'E'
		{5, 1, 6}, // 'T'
		{6, 1, 7}, // just past end
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)", c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_LF(t *testing.T) {
	// "a\nb\nc": a=0 \n=1 b=2 \n=3 c=4
	lt := NewLineTable("a\nb\nc")
	cases := []struct {
		offset            int
		wantLine, wantCol int
	}{
		{0, 1, 1}, // 'a'
		{1, 1, 2}, // '\n' (still line 1)
		{2, 2, 1}, // 'b'
		{3, 2, 2}, // '\n'
		{4, 3, 1}, // 'c'
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)", c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_CRLF(t *testing.T) {
	// "a\r\nb\r\nc": a=0 \r=1 \n=2 b=3 \r=4 \n=5 c=6 — CRLF is ONE break.
	lt := NewLineTable("a\r\nb\r\nc")
	cases := []struct {
		offset            int
		wantLine, wantCol int
	}{
		{0, 1, 1}, {1, 1, 2}, {2, 1, 3}, // line 1: a \r \n
		{3, 2, 1}, {4, 2, 2}, {5, 2, 3}, // line 2: b \r \n
		{6, 3, 1}, // line 3: c
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)", c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_BareCR(t *testing.T) {
	// "a\rb\rc": a=0 \r=1 b=2 \r=3 c=4 — bare CR is ONE break.
	lt := NewLineTable("a\rb\rc")
	cases := []struct {
		offset            int
		wantLine, wantCol int
	}{
		{0, 1, 1}, {1, 1, 2}, {2, 2, 1}, {3, 2, 2}, {4, 3, 1},
	}
	for _, c := range cases {
		line, col := lt.Position(c.offset)
		if line != c.wantLine || col != c.wantCol {
			t.Errorf("Position(%d) = (%d, %d), want (%d, %d)", c.offset, line, col, c.wantLine, c.wantCol)
		}
	}
}

func TestLineTable_NegativeOffset(t *testing.T) {
	lt := NewLineTable("abc")
	line, col := lt.Position(-1)
	if line != 1 || col != 1 {
		t.Errorf("Position(-1) = (%d, %d), want (1, 1)", line, col)
	}
}

func TestLineTable_BeyondEOF(t *testing.T) {
	// Column continues past end of the last line, monotonic and unbounded.
	lt := NewLineTable("abc")
	line, col := lt.Position(100)
	if line != 1 || col != 101 {
		t.Errorf("Position(100) = (%d, %d), want (1, 101)", line, col)
	}
}

func TestLineTable_BeyondEOF_MultiLine(t *testing.T) {
	// "ab\ncd": last line starts at offset 3; an offset of 10 is on line 2.
	lt := NewLineTable("ab\ncd")
	line, col := lt.Position(10)
	if line != 2 {
		t.Errorf("Position(10) line = %d, want 2", line)
	}
	if col != 10-3+1 {
		t.Errorf("Position(10) col = %d, want %d", col, 10-3+1)
	}
}

func TestLineTable_MultiByteUTF8_ColumnsAreBytes(t *testing.T) {
	// "é" is 2 bytes (0xC3 0xA9); the following 'x' is at byte offset 2.
	lt := NewLineTable("éx")
	line, col := lt.Position(2)
	if line != 1 || col != 3 {
		t.Errorf("Position(2) = (%d, %d), want (1, 3) — columns are bytes", line, col)
	}
}
