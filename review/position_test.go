package review

import "testing"

func TestPosition(t *testing.T) {
	const sql = "SELECT 1;\n-- 你好 comment\r\nDELETE FROM t;"
	tests := []struct {
		name   string
		offset int
		line   int
		column int
	}{
		{"start", 0, 1, 1},
		{"inside first line", 7, 1, 8},
		{"newline byte itself", 9, 1, 10},
		{"first byte of second line", 10, 2, 1},
		{"before the first multi-byte rune", 13, 2, 4},
		{"second byte of a multi-byte rune counts it once", 14, 2, 5},
		{"third byte of the same rune gives the same column", 15, 2, 5},
		{"after the first multi-byte rune", 16, 2, 5},
		{"after two multi-byte runes", 19, 2, 6},
		{"carriage return is a column", 28, 2, 15},
		{"third line", 29, 3, 1},
		{"end of input", len(sql), 3, 15},
		{"past the end clamps", len(sql) + 10, 3, 15},
		{"negative clamps", -1, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, column := Position(sql, tt.offset)
			if line != tt.line || column != tt.column {
				t.Fatalf("Position(%d) = %d:%d, want %d:%d", tt.offset, line, column, tt.line, tt.column)
			}
		})
	}
}

func TestPositionEmpty(t *testing.T) {
	if line, column := Position("", 0); line != 1 || column != 1 {
		t.Fatalf("Position(\"\", 0) = %d:%d, want 1:1", line, column)
	}
}
