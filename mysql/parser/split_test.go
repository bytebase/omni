package parser

import (
	"testing"
)

func TestSplitCompoundBlocks(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{name: "BEGIN END", sql: "CREATE PROCEDURE p() BEGIN SELECT 1; SELECT 2; END; SELECT 3;",
			want: []string{"CREATE PROCEDURE p() BEGIN SELECT 1; SELECT 2; END", " SELECT 3"}},
		{name: "nested BEGIN", sql: "CREATE PROCEDURE p() BEGIN BEGIN SELECT 1; END; END; SELECT 2;",
			want: []string{"CREATE PROCEDURE p() BEGIN BEGIN SELECT 1; END; END", " SELECT 2"}},
		{name: "BEGIN WORK is not compound", sql: "BEGIN WORK; SELECT 1;",
			want: []string{"BEGIN WORK", " SELECT 1"}},
		{name: "BEGIN alone is transaction", sql: "BEGIN; SELECT 1;",
			want: []string{"BEGIN", " SELECT 1"}},
		{name: "XA BEGIN is not compound", sql: "XA BEGIN 'xid'; SELECT 1;",
			want: []string{"XA BEGIN 'xid'", " SELECT 1"}},
		{name: "IF END IF", sql: "CREATE FUNCTION f() RETURNS INT BEGIN IF x > 0 THEN SELECT 1; END IF; RETURN 0; END; SELECT 2;",
			want: []string{"CREATE FUNCTION f() RETURNS INT BEGIN IF x > 0 THEN SELECT 1; END IF; RETURN 0; END", " SELECT 2"}},
		{name: "IF EXISTS is not compound", sql: "DROP TABLE IF EXISTS t; SELECT 1;",
			want: []string{"DROP TABLE IF EXISTS t", " SELECT 1"}},
		{name: "IF NOT EXISTS is not compound", sql: "CREATE TABLE IF NOT EXISTS t (id INT); SELECT 1;",
			want: []string{"CREATE TABLE IF NOT EXISTS t (id INT)", " SELECT 1"}},
		{name: "CREATE DATABASE IF NOT EXISTS", sql: "CREATE DATABASE IF NOT EXISTS mydb; USE mydb;",
			want: []string{"CREATE DATABASE IF NOT EXISTS mydb", " USE mydb"}},
		{name: "CASE END", sql: "CREATE PROCEDURE p() BEGIN CASE x WHEN 1 THEN SELECT 1; END CASE; END; SELECT 2;",
			want: []string{"CREATE PROCEDURE p() BEGIN CASE x WHEN 1 THEN SELECT 1; END CASE; END", " SELECT 2"}},
		{name: "WHILE END WHILE", sql: "CREATE PROCEDURE p() BEGIN WHILE x > 0 DO SET x = x - 1; END WHILE; END; SELECT 1;",
			want: []string{"CREATE PROCEDURE p() BEGIN WHILE x > 0 DO SET x = x - 1; END WHILE; END", " SELECT 1"}},
		{name: "REPEAT END REPEAT", sql: "CREATE PROCEDURE p() BEGIN REPEAT SET x = x + 1; UNTIL x > 5 END REPEAT; END; SELECT 1;",
			want: []string{"CREATE PROCEDURE p() BEGIN REPEAT SET x = x + 1; UNTIL x > 5 END REPEAT; END", " SELECT 1"}},
		{name: "LOOP END LOOP", sql: "CREATE PROCEDURE p() BEGIN LOOP SELECT 1; END LOOP; END; SELECT 2;",
			want: []string{"CREATE PROCEDURE p() BEGIN LOOP SELECT 1; END LOOP; END", " SELECT 2"}},
		{name: "REPEAT function is not compound", sql: "SELECT REPEAT('x', 3); SELECT 1;",
			want: []string{"SELECT REPEAT('x', 3)", " SELECT 1"}},
		{name: "IF function is not compound", sql: "SELECT IF(1, 'a', 'b'); SELECT 1;",
			want: []string{"SELECT IF(1, 'a', 'b')", " SELECT 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			}
		})
	}
}

func TestSplitDelimiter(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "basic DELIMITER change",
			sql:  "DELIMITER ;;\nSELECT 1;;\nDELIMITER ;\nSELECT 2;",
			want: []string{"SELECT 1", "SELECT 2"},
		},
		{
			name: "procedure with DELIMITER",
			sql:  "DELIMITER ;;\nCREATE PROCEDURE p()\nBEGIN\n  SELECT 1;\n  SELECT 2;\nEND;;\nDELIMITER ;\nSELECT 3;",
			want: []string{"CREATE PROCEDURE p()\nBEGIN\n  SELECT 1;\n  SELECT 2;\nEND", "SELECT 3"},
		},
		{
			name: "custom delimiter $$",
			sql:  "DELIMITER $$\nSELECT 1$$\nSELECT 2$$\nDELIMITER ;\nSELECT 3;",
			want: []string{"SELECT 1", "\nSELECT 2", "SELECT 3"},
		},
		{
			name: "delimiter with leading whitespace in segment",
			sql:  "DELIMITER ;;\n  SELECT 1 ;;\nDELIMITER ;",
			want: []string{"  SELECT 1 "},
		},
		{
			name: "DELIMITER is case-insensitive",
			sql:  "delimiter ;;\nSELECT 1;;\ndelimiter ;\nSELECT 2;",
			want: []string{"SELECT 1", "SELECT 2"},
		},
		{
			name: "DELIMITER restores normal splitting",
			sql:  "DELIMITER ;;\nSELECT 1;;\nDELIMITER ;\nSELECT 2; SELECT 3;",
			want: []string{"SELECT 1", "SELECT 2", " SELECT 3"},
		},
		{
			name: "DELIMITER // multichar",
			sql:  "DELIMITER //\nSELECT 1//\nDELIMITER ;\nSELECT 2;",
			want: []string{"SELECT 1", "SELECT 2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			var got []string
			for _, s := range segs {
				got = append(got, s.Text)
			}
			if len(got) == 0 {
				got = nil
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d segments, want %d\n  got:  %q\n  want: %q",
					len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("seg[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSplitQuotingAndComments(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "semicolon in single-quoted string",
			sql:  "SELECT 'a;b'; SELECT 1;",
			want: []string{"SELECT 'a;b'", " SELECT 1"},
		},
		{
			name: "semicolon in double-quoted string",
			sql:  `SELECT "a;b"; SELECT 1;`,
			want: []string{`SELECT "a;b"`, " SELECT 1"},
		},
		{
			name: "semicolon in backtick identifier",
			sql:  "SELECT `a;b`; SELECT 1;",
			want: []string{"SELECT `a;b`", " SELECT 1"},
		},
		{
			name: "escaped single quote with backslash",
			sql:  `SELECT 'it\'s'; SELECT 1;`,
			want: []string{`SELECT 'it\'s'`, " SELECT 1"},
		},
		{
			name: "escaped single quote with double",
			sql:  "SELECT 'it''s'; SELECT 1;",
			want: []string{"SELECT 'it''s'", " SELECT 1"},
		},
		{
			name: "semicolon in block comment",
			sql:  "SELECT /* ; */ 1; SELECT 2;",
			want: []string{"SELECT /* ; */ 1", " SELECT 2"},
		},
		{
			name: "semicolon in line comment --",
			sql:  "SELECT 1 -- ;\n; SELECT 2;",
			want: []string{"SELECT 1 -- ;\n", " SELECT 2"},
		},
		{
			name: "semicolon in hash comment",
			sql:  "SELECT 1 # ;\n; SELECT 2;",
			want: []string{"SELECT 1 # ;\n", " SELECT 2"},
		},
		{
			name: "nested block comments",
			sql:  "SELECT /* /* ; */ */ 1; SELECT 2;",
			want: []string{"SELECT /* /* ; */ */ 1", " SELECT 2"},
		},
		{
			name: "-- without space is not comment",
			sql:  "SELECT 1--2; SELECT 3;",
			want: []string{"SELECT 1--2", " SELECT 3"},
		},
		{
			name: "conditional comment",
			sql:  "SELECT /*!50000 1, */ 2; SELECT 3;",
			want: []string{"SELECT /*!50000 1, */ 2", " SELECT 3"},
		},
		{
			name: "CRLF line endings",
			sql:  "SELECT 1;\r\nSELECT 2;",
			want: []string{"SELECT 1", "\r\nSELECT 2"},
		},
		{
			name: "Unicode in identifiers",
			sql:  "SELECT * FROM 表名; SELECT 1;",
			want: []string{"SELECT * FROM 表名", " SELECT 1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			var got []string
			for _, s := range segs {
				got = append(got, s.Text)
			}
			if len(got) == 0 {
				got = nil
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d segments, want %d\n  got:  %q\n  want: %q",
					len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("seg[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSplitByteOffsets(t *testing.T) {
	sql := "SELECT 1; SELECT 2; SELECT 3;"
	segs := Split(sql)
	for _, s := range segs {
		got := sql[s.ByteStart:s.ByteEnd]
		if got != s.Text {
			t.Errorf("sql[%d:%d] = %q, but Text = %q", s.ByteStart, s.ByteEnd, got, s.Text)
		}
	}
}

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
