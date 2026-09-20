package review

import (
	"strings"
	"testing"
	"unicode/utf8"
)

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
	index := Index(sql)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, column := Position(sql, tt.offset)
			if line != tt.line || column != tt.column {
				t.Fatalf("Position(%d) = %d:%d, want %d:%d", tt.offset, line, column, tt.line, tt.column)
			}
			if l, c := index.Position(tt.offset); l != line || c != column {
				t.Fatalf("Index.Position(%d) = %d:%d, Position = %d:%d", tt.offset, l, c, line, column)
			}
		})
	}
}

func TestPositionEmpty(t *testing.T) {
	if line, column := Position("", 0); line != 1 || column != 1 {
		t.Fatalf("Position(\"\", 0) = %d:%d, want 1:1", line, column)
	}
}

// TestIndexMatchesScan checks the indexed conversion against a single
// forward scan at every byte offset: many short lines with blank lines
// and CRLF, then one line far longer than a checkpoint interval whose
// 9-byte period puts checkpoint thresholds inside multi-byte runes.
func TestIndexMatchesScan(t *testing.T) {
	sql := strings.Repeat("a\n\nб𝄞\r\nSELECT '你好';\n", 20) + "\n\n" +
		strings.Repeat("你好abc", 4*checkpointBytes/9) + "\n" +
		strings.Repeat("𝄞", checkpointBytes)
	if n := 4 * checkpointBytes / 9 * 9; n <= 3*checkpointBytes {
		t.Fatalf("long line is %d bytes, not longer than three checkpoints", n)
	}
	index := Index(sql)
	if got, min := len(index.checkpointAt), len(sql)/(checkpointBytes+utf8.UTFMax); got < min {
		t.Fatalf("%d checkpoints for %d bytes, want at least %d", got, len(sql), min)
	}
	for i := 1; i < len(index.checkpointAt); i++ {
		if gap := index.checkpointAt[i] - index.checkpointAt[i-1]; gap < checkpointBytes || gap > checkpointBytes+utf8.UTFMax {
			t.Fatalf("checkpoint gap %d at %d", gap, i)
		}
	}
	assertScan(t, sql, index)
}

func assertScan(t *testing.T, sql string, index *Text) {
	t.Helper()
	line, column := 1, 1
	for offset := 0; offset <= len(sql); {
		if l, c := index.Position(offset); l != line || c != column {
			t.Fatalf("offset %d: index gives %d:%d, scan gives %d:%d", offset, l, c, line, column)
		}
		if offset == len(sql) {
			break
		}
		r, size := utf8.DecodeRuneInString(sql[offset:])
		for k := 1; k < size; k++ {
			if l, c := index.Position(offset + k); l != line || c != column+1 {
				t.Fatalf("offset %d inside rune at %d: got %d:%d, want %d:%d", offset+k, offset, l, c, line, column+1)
			}
		}
		offset += size
		if r == '\n' {
			line, column = line+1, 1
		} else {
			column++
		}
	}
}
