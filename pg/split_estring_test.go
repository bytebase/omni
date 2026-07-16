package pg

import "testing"

// TestSplitEscapeString covers E'...' escape-string handling in the
// splitter (audit defect D1): backslash escapes the next character
// inside E-strings, so a \' does not close the string. Plain '...'
// strings keep standard_conforming_strings=on semantics where
// backslash is a literal character.
func TestSplitEscapeString(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "escaped quote does not close E-string",
			sql:  `SELECT E'a\';b';`,
			want: []string{`SELECT E'a\';b';`},
		},
		{
			name: "E-string followed by second statement",
			sql:  `SELECT E'a\';b'; SELECT 2;`,
			want: []string{`SELECT E'a\';b';`, ` SELECT 2;`},
		},
		{
			name: "lowercase e prefix",
			sql:  `SELECT e'a\';b';`,
			want: []string{`SELECT e'a\';b';`},
		},
		{
			name: "even backslashes close normally",
			sql:  `SELECT E'a\\'; SELECT 2;`,
			want: []string{`SELECT E'a\\';`, ` SELECT 2;`},
		},
		{
			name: "doubled quote inside E-string",
			sql:  `SELECT E'it''s \' ok';`,
			want: []string{`SELECT E'it''s \' ok';`},
		},
		// Negative fences: these must NOT get escape-aware scanning.
		{
			name: "plain string backslash is literal (conforming on)",
			sql:  `SELECT 'a\'; SELECT 2;`,
			want: []string{`SELECT 'a\';`, ` SELECT 2;`},
		},
		{
			name: "E-string at input start (lookback boundary)",
			sql:  `E'a\';b'; SELECT 2;`,
			want: []string{`E'a\';b';`, ` SELECT 2;`},
		},
		{
			name: "E as identifier tail is not an escape prefix",
			sql:  `SELECT abcE'x\'; SELECT 2;`,
			want: []string{`SELECT abcE'x\';`, ` SELECT 2;`},
		},
		{
			name: "U& string is not escape-aware",
			sql:  `SELECT U&'d\0061t;a'; SELECT 2;`,
			want: []string{`SELECT U&'d\0061t;a';`, ` SELECT 2;`},
		},
		// BEGIN ATOMIC body path (second dispatch site).
		{
			name: "E-string inside BEGIN ATOMIC body",
			sql:  "CREATE FUNCTION f() RETURNS text LANGUAGE sql BEGIN ATOMIC SELECT E'a\\';b'; END; SELECT 2;",
			want: []string{"CREATE FUNCTION f() RETURNS text LANGUAGE sql BEGIN ATOMIC SELECT E'a\\';b'; END;", " SELECT 2;"},
		},
		// Resilience: unterminated constructs must not panic or hang.
		{
			name: "unterminated E-string at EOF",
			sql:  `SELECT E'abc\'`,
			want: []string{`SELECT E'abc\'`},
		},
		{
			name: "trailing backslash at EOF",
			sql:  `SELECT E'abc\`,
			want: []string{`SELECT E'abc\`},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			segs := Split(c.sql)
			if len(segs) != len(c.want) {
				t.Fatalf("got %d segments, want %d: %#v", len(segs), len(c.want), segs)
			}
			joined := ""
			for i, s := range segs {
				if s.Text != c.want[i] {
					t.Errorf("segment %d = %q, want %q", i, s.Text, c.want[i])
				}
				if s.Text != c.sql[s.ByteStart:s.ByteEnd] {
					t.Errorf("segment %d Range mismatch: Text=%q, sql[%d:%d]=%q", i, s.Text, s.ByteStart, s.ByteEnd, c.sql[s.ByteStart:s.ByteEnd])
				}
				joined += s.Text
			}
			if joined != c.sql {
				t.Errorf("lossless invariant violated: concat=%q, input=%q", joined, c.sql)
			}
		})
	}
}
