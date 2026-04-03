package parser

import (
	"testing"
)

func TestSplitSimple(t *testing.T) {
	tests := []struct {
		sql  string
		want []string // expected segment texts (non-empty only)
	}{
		{"SELECT 1", []string{"SELECT 1"}},
		{"SELECT 1;", []string{"SELECT 1"}},
		{"SELECT 1; SELECT 2;", []string{"SELECT 1", " SELECT 2"}},
		{"SELECT 1;  ", []string{"SELECT 1"}},
		{"", nil},
		{";;;", nil},
		{" ; ; ", nil},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			segs := Split(tt.sql)
			var got []string
			for _, s := range segs {
				got = append(got, s.Text)
			}
			if len(got) == 0 {
				got = nil
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Split(%q) = %v, want %v", tt.sql, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Split(%q)[%d] = %q, want %q", tt.sql, i, got[i], tt.want[i])
				}
				// Verify byte offset identity.
				seg := segs[i]
				if tt.sql[seg.ByteStart:seg.ByteEnd] != seg.Text {
					t.Errorf("Split(%q)[%d]: byte offset identity failed: sql[%d:%d] = %q, seg.Text = %q",
						tt.sql, i, seg.ByteStart, seg.ByteEnd, tt.sql[seg.ByteStart:seg.ByteEnd], seg.Text)
				}
			}
		})
	}
}
