package mssql

import "testing"

func TestMSSQLPublicParseLocPrecisionMatrix(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantTexts  []string
		wantStarts []Position
		wantEnds   []Position
	}{
		{
			name:       "single statement text and positions",
			sql:        "SELECT 1",
			wantTexts:  []string{"SELECT 1"},
			wantStarts: []Position{{Line: 1, Column: 1}},
			wantEnds:   []Position{{Line: 1, Column: 9}},
		},
		{
			name:       "leading whitespace start position",
			sql:        " \n\tSELECT 1",
			wantTexts:  []string{" \n\tSELECT 1"},
			wantStarts: []Position{{Line: 2, Column: 2}},
			wantEnds:   []Position{{Line: 2, Column: 10}},
		},
		{
			name:       "multi statement text",
			sql:        "SELECT 1; SELECT 2",
			wantTexts:  []string{"SELECT 1;", " SELECT 2"},
			wantStarts: []Position{{Line: 1, Column: 1}, {Line: 1, Column: 11}},
			wantEnds:   []Position{{Line: 1, Column: 10}, {Line: 1, Column: 19}},
		},
		{
			name:       "tabs and unicode use byte columns",
			sql:        "\tSELECT N'你好'",
			wantTexts:  []string{"\tSELECT N'你好'"},
			wantStarts: []Position{{Line: 1, Column: 2}},
			wantEnds:   []Position{{Line: 1, Column: 18}},
		},
		{
			name:       "go batch separator",
			sql:        "SELECT 1\nGO\nSELECT 2",
			wantTexts:  []string{"SELECT 1", "\nGO", "\nSELECT 2"},
			wantStarts: []Position{{Line: 1, Column: 1}, {Line: 2, Column: 1}, {Line: 3, Column: 1}},
			wantEnds:   []Position{{Line: 1, Column: 9}, {Line: 2, Column: 3}, {Line: 3, Column: 9}},
		},
		{
			name:       "empty statements skipped",
			sql:        ";;SELECT 1",
			wantTexts:  []string{";;SELECT 1"},
			wantStarts: []Position{{Line: 1, Column: 3}},
			wantEnds:   []Position{{Line: 1, Column: 11}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
			if len(stmts) != len(tt.wantTexts) {
				t.Fatalf("statement count = %d, want %d: %#v", len(stmts), len(tt.wantTexts), stmts)
			}
			for i, stmt := range stmts {
				if stmt.Text != tt.wantTexts[i] {
					t.Fatalf("stmt[%d].Text = %q, want %q", i, stmt.Text, tt.wantTexts[i])
				}
				if stmt.ByteStart < 0 || stmt.ByteEnd < stmt.ByteStart || stmt.ByteEnd > len(tt.sql) {
					t.Fatalf("stmt[%d] byte range invalid: %d:%d", i, stmt.ByteStart, stmt.ByteEnd)
				}
				if stmt.Text != tt.sql[stmt.ByteStart:stmt.ByteEnd] {
					t.Fatalf("stmt[%d].Text does not match byte range %d:%d", i, stmt.ByteStart, stmt.ByteEnd)
				}
				if stmt.Start != tt.wantStarts[i] {
					t.Fatalf("stmt[%d].Start = %+v, want %+v", i, stmt.Start, tt.wantStarts[i])
				}
				if stmt.End != tt.wantEnds[i] {
					t.Fatalf("stmt[%d].End = %+v, want %+v", i, stmt.End, tt.wantEnds[i])
				}
			}
		})
	}
}
