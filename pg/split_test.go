package pg

import (
	"os"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Segment
		empties []bool // expected Empty() for each segment
	}{
		// Basic cases.
		{
			name:  "single statement with semicolon",
			input: "SELECT 1;",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9},
			},
			empties: []bool{false},
		},
		{
			name:  "single statement without semicolon",
			input: "SELECT 1",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8},
			},
			empties: []bool{false},
		},
		{
			name:  "multiple statements",
			input: "SELECT 1; SELECT 2; SELECT 3",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9},
				{Text: " SELECT 2;", ByteStart: 9, ByteEnd: 19},
				{Text: " SELECT 3", ByteStart: 19, ByteEnd: 28},
			},
			empties: []bool{false, false, false},
		},
		{
			name:  "empty input",
			input: "",
			want:  nil,
		},
		{
			name:  "whitespace only",
			input: "  \t\n  ",
			want: []Segment{
				{Text: "  \t\n  ", ByteStart: 0, ByteEnd: 6},
			},
			empties: []bool{true},
		},

		// String/identifier handling.
		{
			name:  "semicolon inside single-quoted string",
			input: "SELECT 'a;b'; SELECT 2",
			want: []Segment{
				{Text: "SELECT 'a;b';", ByteStart: 0, ByteEnd: 13},
				{Text: " SELECT 2", ByteStart: 13, ByteEnd: 22},
			},
			empties: []bool{false, false},
		},
		{
			name:  "semicolon inside double-quoted identifier",
			input: `SELECT "a;b"; SELECT 2`,
			want: []Segment{
				{Text: `SELECT "a;b";`, ByteStart: 0, ByteEnd: 13},
				{Text: " SELECT 2", ByteStart: 13, ByteEnd: 22},
			},
			empties: []bool{false, false},
		},
		{
			name:  "semicolon inside dollar-quoted string",
			input: "SELECT $$a;b$$; SELECT 2",
			want: []Segment{
				{Text: "SELECT $$a;b$$;", ByteStart: 0, ByteEnd: 15},
				{Text: " SELECT 2", ByteStart: 15, ByteEnd: 24},
			},
			empties: []bool{false, false},
		},
		{
			name:  "tagged dollar-quote",
			input: "SELECT $tag$a;b$tag$; SELECT 2",
			want: []Segment{
				{Text: "SELECT $tag$a;b$tag$;", ByteStart: 0, ByteEnd: 21},
				{Text: " SELECT 2", ByteStart: 21, ByteEnd: 30},
			},
			empties: []bool{false, false},
		},

		// Comment handling.
		{
			name:  "semicolon inside block comment",
			input: "SELECT /* ; */ 1; SELECT 2",
			want: []Segment{
				{Text: "SELECT /* ; */ 1;", ByteStart: 0, ByteEnd: 17},
				{Text: " SELECT 2", ByteStart: 17, ByteEnd: 26},
			},
			empties: []bool{false, false},
		},
		{
			name:  "semicolon inside line comment",
			input: "SELECT 1 -- ;\n; SELECT 2",
			want: []Segment{
				{Text: "SELECT 1 -- ;\n;", ByteStart: 0, ByteEnd: 15},
				{Text: " SELECT 2", ByteStart: 15, ByteEnd: 24},
			},
			empties: []bool{false, false},
		},
		{
			name:  "nested block comments",
			input: "SELECT /* /* ; */ */ 1; SELECT 2",
			want: []Segment{
				{Text: "SELECT /* /* ; */ */ 1;", ByteStart: 0, ByteEnd: 23},
				{Text: " SELECT 2", ByteStart: 23, ByteEnd: 32},
			},
			empties: []bool{false, false},
		},
		{
			name:  "comments only",
			input: "-- comment\n/* block */",
			want: []Segment{
				{Text: "-- comment\n/* block */", ByteStart: 0, ByteEnd: 22},
			},
			empties: []bool{true},
		},

		// BEGIN ATOMIC.
		{
			name:  "simple BEGIN ATOMIC END",
			input: "CREATE FUNCTION f() BEGIN ATOMIC SELECT 1; SELECT 2; END; SELECT 3",
			want: []Segment{
				{Text: "CREATE FUNCTION f() BEGIN ATOMIC SELECT 1; SELECT 2; END;", ByteStart: 0, ByteEnd: 57},
				{Text: " SELECT 3", ByteStart: 57, ByteEnd: 66},
			},
			empties: []bool{false, false},
		},
		{
			name:  "BEGIN ATOMIC with CASE END inside",
			input: "BEGIN ATOMIC CASE WHEN true THEN 1 END; END; SELECT 1",
			want: []Segment{
				{Text: "BEGIN ATOMIC CASE WHEN true THEN 1 END; END;", ByteStart: 0, ByteEnd: 44},
				{Text: " SELECT 1", ByteStart: 44, ByteEnd: 53},
			},
			empties: []bool{false, false},
		},
		{
			name:  "BEGIN ATOMIC with nested CASE",
			input: "BEGIN ATOMIC CASE WHEN x THEN CASE WHEN y THEN 1 END END; END; SELECT 1",
			want: []Segment{
				{Text: "BEGIN ATOMIC CASE WHEN x THEN CASE WHEN y THEN 1 END END; END;", ByteStart: 0, ByteEnd: 62},
				{Text: " SELECT 1", ByteStart: 62, ByteEnd: 71},
			},
			empties: []bool{false, false},
		},
		{
			name:  "case insensitive begin atomic",
			input: "begin atomic select 1; end; SELECT 2",
			want: []Segment{
				{Text: "begin atomic select 1; end;", ByteStart: 0, ByteEnd: 27},
				{Text: " SELECT 2", ByteStart: 27, ByteEnd: 36},
			},
			empties: []bool{false, false},
		},
		{
			name:  "BEGIN without ATOMIC is normal transaction",
			input: "BEGIN; SELECT 1; COMMIT;",
			want: []Segment{
				{Text: "BEGIN;", ByteStart: 0, ByteEnd: 6},
				{Text: " SELECT 1;", ByteStart: 6, ByteEnd: 16},
				{Text: " COMMIT;", ByteStart: 16, ByteEnd: 24},
			},
			empties: []bool{false, false, false},
		},
		{
			name:  "BEGIN comment ATOMIC enters block mode",
			input: "BEGIN /* comment */ ATOMIC SELECT 1; END; SELECT 2",
			want: []Segment{
				{Text: "BEGIN /* comment */ ATOMIC SELECT 1; END;", ByteStart: 0, ByteEnd: 41},
				{Text: " SELECT 2", ByteStart: 41, ByteEnd: 50},
			},
			empties: []bool{false, false},
		},
		{
			name:  "WEEKEND not matched as END",
			input: "BEGIN ATOMIC WEEKEND; END; SELECT 1",
			want: []Segment{
				{Text: "BEGIN ATOMIC WEEKEND; END;", ByteStart: 0, ByteEnd: 26},
				{Text: " SELECT 1", ByteStart: 26, ByteEnd: 35},
			},
			empties: []bool{false, false},
		},

		// PL/pgSQL blocks inside dollar-quotes.
		// These validate that BEGIN/END, IF/END IF, LOOP/END LOOP, CASE/END CASE
		// inside dollar-quoted function bodies do NOT cause incorrect splitting.
		{
			name: "CREATE FUNCTION with BEGIN END in dollar-quote",
			input: `CREATE FUNCTION foo() RETURNS void AS $$
BEGIN
  INSERT INTO t VALUES (1);
  UPDATE t SET a = 2;
END;
$$ LANGUAGE plpgsql; SELECT 1;`,
			want: []Segment{
				{Text: "CREATE FUNCTION foo() RETURNS void AS $$\nBEGIN\n  INSERT INTO t VALUES (1);\n  UPDATE t SET a = 2;\nEND;\n$$ LANGUAGE plpgsql;", ByteStart: 0, ByteEnd: 122},
				{Text: " SELECT 1;", ByteStart: 122, ByteEnd: 132},
			},
			empties: []bool{false, false},
		},
		{
			name: "CREATE FUNCTION with IF END IF in dollar-quote",
			input: `CREATE FUNCTION bar() RETURNS void AS $$
BEGIN
  IF x > 0 THEN
    INSERT INTO t VALUES (1);
  END IF;
END;
$$ LANGUAGE plpgsql; SELECT 2;`,
			want: []Segment{
				{Text: "CREATE FUNCTION bar() RETURNS void AS $$\nBEGIN\n  IF x > 0 THEN\n    INSERT INTO t VALUES (1);\n  END IF;\nEND;\n$$ LANGUAGE plpgsql;", ByteStart: 0, ByteEnd: 128},
				{Text: " SELECT 2;", ByteStart: 128, ByteEnd: 138},
			},
			empties: []bool{false, false},
		},
		{
			name: "CREATE FUNCTION with LOOP END LOOP in dollar-quote",
			input: `CREATE FUNCTION baz() RETURNS void AS $$
BEGIN
  LOOP
    EXIT WHEN done;
    INSERT INTO t VALUES (1);
  END LOOP;
END;
$$ LANGUAGE plpgsql; SELECT 3;`,
			want: []Segment{
				{Text: "CREATE FUNCTION baz() RETURNS void AS $$\nBEGIN\n  LOOP\n    EXIT WHEN done;\n    INSERT INTO t VALUES (1);\n  END LOOP;\nEND;\n$$ LANGUAGE plpgsql;", ByteStart: 0, ByteEnd: 141},
				{Text: " SELECT 3;", ByteStart: 141, ByteEnd: 151},
			},
			empties: []bool{false, false},
		},
		{
			name: "CREATE FUNCTION with CASE END in dollar-quote",
			input: `CREATE FUNCTION qux() RETURNS int AS $$
BEGIN
  RETURN CASE WHEN x > 0 THEN 1 ELSE 2 END;
END;
$$ LANGUAGE plpgsql; SELECT 4;`,
			want: []Segment{
				{Text: "CREATE FUNCTION qux() RETURNS int AS $$\nBEGIN\n  RETURN CASE WHEN x > 0 THEN 1 ELSE 2 END;\nEND;\n$$ LANGUAGE plpgsql;", ByteStart: 0, ByteEnd: 115},
				{Text: " SELECT 4;", ByteStart: 115, ByteEnd: 125},
			},
			empties: []bool{false, false},
		},
		{
			name: "CREATE FUNCTION with tagged dollar-quote and nested blocks",
			input: `CREATE FUNCTION complex() RETURNS void AS $fn$
BEGIN
  IF x > 0 THEN
    LOOP
      EXIT WHEN y < 0;
      INSERT INTO t VALUES (CASE WHEN a THEN 1 ELSE 2 END);
    END LOOP;
  END IF;
END;
$fn$ LANGUAGE plpgsql; SELECT 5;`,
			want: []Segment{
				{Text: "CREATE FUNCTION complex() RETURNS void AS $fn$\nBEGIN\n  IF x > 0 THEN\n    LOOP\n      EXIT WHEN y < 0;\n      INSERT INTO t VALUES (CASE WHEN a THEN 1 ELSE 2 END);\n    END LOOP;\n  END IF;\nEND;\n$fn$ LANGUAGE plpgsql;", ByteStart: 0, ByteEnd: 212},
				{Text: " SELECT 5;", ByteStart: 212, ByteEnd: 222},
			},
			empties: []bool{false, false},
		},
		{
			name: "DO block with dollar-quote",
			input: `DO $$
BEGIN
  INSERT INTO t VALUES (1);
  INSERT INTO t VALUES (2);
END;
$$; SELECT 1;`,
			want: []Segment{
				{Text: "DO $$\nBEGIN\n  INSERT INTO t VALUES (1);\n  INSERT INTO t VALUES (2);\nEND;\n$$;", ByteStart: 0, ByteEnd: 76},
				{Text: " SELECT 1;", ByteStart: 76, ByteEnd: 86},
			},
			empties: []bool{false, false},
		},
		{
			name:  "BEGIN TRANSACTION is normal split",
			input: "BEGIN TRANSACTION; SELECT 1; COMMIT;",
			want: []Segment{
				{Text: "BEGIN TRANSACTION;", ByteStart: 0, ByteEnd: 18},
				{Text: " SELECT 1;", ByteStart: 18, ByteEnd: 28},
				{Text: " COMMIT;", ByteStart: 28, ByteEnd: 36},
			},
			empties: []bool{false, false, false},
		},
		{
			name:  "BEGIN WORK is normal split",
			input: "BEGIN WORK; SELECT 1; END WORK;",
			want: []Segment{
				{Text: "BEGIN WORK;", ByteStart: 0, ByteEnd: 11},
				{Text: " SELECT 1;", ByteStart: 11, ByteEnd: 21},
				{Text: " END WORK;", ByteStart: 21, ByteEnd: 31},
			},
			empties: []bool{false, false, false},
		},

		// Unterminated constructs.
		{
			name:  "unterminated single quote",
			input: "SELECT 'abc; SELECT 2",
			want: []Segment{
				{Text: "SELECT 'abc; SELECT 2", ByteStart: 0, ByteEnd: 21},
			},
			empties: []bool{false},
		},
		{
			name:  "unterminated double quote",
			input: `SELECT "abc; SELECT 2`,
			want: []Segment{
				{Text: `SELECT "abc; SELECT 2`, ByteStart: 0, ByteEnd: 21},
			},
			empties: []bool{false},
		},
		{
			name:  "unterminated dollar quote",
			input: "SELECT $$abc; SELECT 2",
			want: []Segment{
				{Text: "SELECT $$abc; SELECT 2", ByteStart: 0, ByteEnd: 22},
			},
			empties: []bool{false},
		},
		{
			name:  "unterminated block comment",
			input: "SELECT /* abc; SELECT 2",
			want: []Segment{
				{Text: "SELECT /* abc; SELECT 2", ByteStart: 0, ByteEnd: 23},
			},
			empties: []bool{false},
		},

		// Edge cases.
		{
			name:  "trailing content after last semicolon",
			input: "SELECT 1;  ",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9},
				{Text: "  ", ByteStart: 9, ByteEnd: 11},
			},
			empties: []bool{false, true},
		},
		{
			name:  "multiple semicolons in a row",
			input: ";;;",
			want: []Segment{
				{Text: ";", ByteStart: 0, ByteEnd: 1},
				{Text: ";", ByteStart: 1, ByteEnd: 2},
				{Text: ";", ByteStart: 2, ByteEnd: 3},
			},
			empties: []bool{true, true, true},
		},
		{
			name:  "CRLF line endings",
			input: "SELECT 1;\r\nSELECT 2;",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9},
				{Text: "\r\nSELECT 2;", ByteStart: 9, ByteEnd: 20},
			},
			empties: []bool{false, false},
		},
		{
			name:  "escaped single quotes",
			input: "SELECT 'it''s'; SELECT 2",
			want: []Segment{
				{Text: "SELECT 'it''s';", ByteStart: 0, ByteEnd: 15},
				{Text: " SELECT 2", ByteStart: 15, ByteEnd: 24},
			},
			empties: []bool{false, false},
		},
		{
			name:  "escaped double quotes",
			input: `SELECT "a""b"; SELECT 2`,
			want: []Segment{
				{Text: `SELECT "a""b";`, ByteStart: 0, ByteEnd: 14},
				{Text: " SELECT 2", ByteStart: 14, ByteEnd: 23},
			},
			empties: []bool{false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("Split(%q) returned %d segments, want %d\ngot:  %+v", tt.input, len(got), len(tt.want), got)
			}

			for i, g := range got {
				w := tt.want[i]
				if g.Text != w.Text || g.ByteStart != w.ByteStart || g.ByteEnd != w.ByteEnd {
					t.Errorf("segment[%d] = %+v, want %+v", i, g, w)
				}
			}

			for i, g := range got {
				if i < len(tt.empties) {
					if g.Empty() != tt.empties[i] {
						t.Errorf("segment[%d].Empty() = %v, want %v (text=%q)", i, g.Empty(), tt.empties[i], g.Text)
					}
				}
			}
		})
	}
}

// TestSplitDollarIdentAdjacency pins the scan.l rule that '$' continues an
// identifier: abc$x$y is one identifier, not "abc" + dollar-quote $x$. A
// mis-lex here swallows all input from the false opening tag to EOF (or, in
// the other direction, splits inside a real dollar-quote). Engine-verified
// on PostgreSQL 17: SELECT 1 AS abc$x$y, 名$x$y, _a$t$b all parse as plain
// aliases; $tag$ after a non-identifier byte opens a quote as usual.
func TestSplitDollarIdentAdjacency(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string // exact non-empty segment texts
	}{
		{
			name:  "identifier with dollar then real statement",
			input: `SELECT 1 AS abc$x$y; SELECT 2;`,
			want:  []string{`SELECT 1 AS abc$x$y;`, ` SELECT 2;`},
		},
		{
			name:  "identifier with multiple dollars",
			input: `SELECT 1 AS abc$x$y$z; SELECT 2;`,
			want:  []string{`SELECT 1 AS abc$x$y$z;`, ` SELECT 2;`},
		},
		{
			name:  "multibyte identifier tail dollar",
			input: "SELECT 1 AS 名$x$y; SELECT 2;",
			want:  []string{"SELECT 1 AS 名$x$y;", " SELECT 2;"},
		},
		{
			name:  "underscore identifier tail dollar",
			input: `SELECT 1 AS _a$t$b; SELECT 2;`,
			want:  []string{`SELECT 1 AS _a$t$b;`, ` SELECT 2;`},
		},
		{
			name:  "identifier absorbs double dollar (scan.l x$$a is one identifier)",
			input: `SELECT 1 AS x$$a;b$$; SELECT 2;`,
			want:  []string{`SELECT 1 AS x$$a;`, `b$$;`, ` SELECT 2;`},
		},
		{
			name:  "dollar quote still opens after operator",
			input: `SELECT $tag$a;b$tag$; SELECT 2;`,
			want:  []string{`SELECT $tag$a;b$tag$;`, ` SELECT 2;`},
		},
		{
			name:  "dollar quote at input start",
			input: `$t$;$t$; SELECT 2;`,
			want:  []string{`$t$;$t$;`, ` SELECT 2;`},
		},
		{
			name:  "dollar quote after open paren",
			input: `SELECT f($q$a;b$q$); SELECT 2;`,
			want:  []string{`SELECT f($q$a;b$q$);`, ` SELECT 2;`},
		},
		{
			// Documented divergence: PG lexes tokens, so 123$t$a;b$t$ is
			// number + string there (single grammar-INVALID statement, since
			// a literal cannot be adjacent to a number). The byte scanner
			// keeps '$' after a digit as identifier continuation and splits
			// at the inner semicolon. Both sides reject the input; only the
			// error boundary differs. Accepted in the splitter audit.
			name:  "digit-adjacent dollar stays byte-scanned (known divergence)",
			input: `SELECT 123$t$a;b$t$`,
			want:  []string{`SELECT 123$t$a;`, `b$t$`},
		},
		{
			name:  "dollar quote inside BEGIN ATOMIC after identifier",
			input: "CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT 1 AS a$x$b; END; SELECT 2;",
			want:  []string{"CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT 1 AS a$x$b; END;", " SELECT 2;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.input)

			// Lossless invariant: offsets must reconstruct the input.
			var rebuilt []byte
			for _, s := range segs {
				if s.Text != tt.input[s.ByteStart:s.ByteEnd] {
					t.Fatalf("segment text %q does not match input[%d:%d]", s.Text, s.ByteStart, s.ByteEnd)
				}
				rebuilt = append(rebuilt, s.Text...)
			}
			if string(rebuilt) != tt.input {
				t.Fatalf("segments do not reconstruct input:\ngot  %q\nwant %q", rebuilt, tt.input)
			}

			var got []string
			for _, s := range segs {
				if !s.Empty() {
					got = append(got, s.Text)
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d non-empty segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSplitParenDepth pins the rule that statement-separating semicolons
// only occur at parenthesis depth zero. PG's grammar nests semicolons
// inside parentheses (CREATE RULE multi-action lists), and psqlscan never
// splits inside parens. Engine-verified on PostgreSQL 17: the rule
// statement executes as ONE statement; a stray ')' splits normally
// (clamped); an unclosed '(' buffers the remainder to end of input.
func TestSplitParenDepth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string // exact non-empty segment texts
	}{
		{
			// Audit shape (3 wrong segments before the fix).
			name:  "rule multi-action then statement",
			input: `CREATE RULE r AS ON UPDATE TO t DO ALSO (SELECT 1; SELECT 2); SELECT 3;`,
			want:  []string{`CREATE RULE r AS ON UPDATE TO t DO ALSO (SELECT 1; SELECT 2);`, ` SELECT 3;`},
		},
		{
			// Cross-validation shape (2 wrong segments before the fix).
			name:  "rule multi-action alone",
			input: `CREATE RULE r AS ON DELETE TO t DO INSTEAD (DELETE FROM a; DELETE FROM b);`,
			want:  []string{`CREATE RULE r AS ON DELETE TO t DO INSTEAD (DELETE FROM a; DELETE FROM b);`},
		},
		{
			name:  "deeply nested semicolon",
			input: `CREATE RULE r AS ON UPDATE TO t DO ALSO (UPDATE a SET x = (1); DELETE FROM b); SELECT 2;`,
			want:  []string{`CREATE RULE r AS ON UPDATE TO t DO ALSO (UPDATE a SET x = (1); DELETE FROM b);`, ` SELECT 2;`},
		},
		{
			// Resilience: unclosed '(' leaves the remainder as one segment —
			// psql buffers to end of input the same way. Fenced so a future
			// "fix" does not silently change it.
			name:  "unclosed paren swallows remainder",
			input: `SELECT (1; SELECT 2;`,
			want:  []string{`SELECT (1; SELECT 2;`},
		},
		{
			// A stray ')' must not corrupt depth: clamp at zero and split
			// normally (engine: statement 1 errors, statement 2 runs).
			name:  "stray close paren splits normally",
			input: `SELECT 1); SELECT 2;`,
			want:  []string{`SELECT 1);`, ` SELECT 2;`},
		},
		{
			// Parens inside strings, comments, and dollar-quotes do not
			// count toward depth.
			name:  "paren-like bytes in literals do not count",
			input: `SELECT '(', $t$($t$ /* ( */; SELECT 2;`,
			want:  []string{`SELECT '(', $t$($t$ /* ( */;`, ` SELECT 2;`},
		},
		{
			name:  "function call args with semicolon in dollar quote",
			input: `SELECT f($q$a;b$q$, 1); SELECT 2;`,
			want:  []string{`SELECT f($q$a;b$q$, 1);`, ` SELECT 2;`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.input)

			var rebuilt []byte
			for _, s := range segs {
				if s.Text != tt.input[s.ByteStart:s.ByteEnd] {
					t.Fatalf("segment text %q does not match input[%d:%d]", s.Text, s.ByteStart, s.ByteEnd)
				}
				rebuilt = append(rebuilt, s.Text...)
			}
			if string(rebuilt) != tt.input {
				t.Fatalf("segments do not reconstruct input:\ngot  %q\nwant %q", rebuilt, tt.input)
			}

			var got []string
			for _, s := range segs {
				if !s.Empty() {
					got = append(got, s.Text)
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d non-empty segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSplitCopyFromStdin pins inline COPY data handling: psql scripts and
// pg_dump plain-format output carry the data lines in the SQL stream, the
// data may contain semicolons, and the block ends only at a line that is
// exactly "\.". The statement, its data, and the terminator stay one
// segment. Engine-verified on PostgreSQL 17 (lowercase copy, WITH options,
// column lists, semicolons in data all load correctly via psql).
func TestSplitCopyFromStdin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "pg_dump shape with semicolons in data",
			input: "COPY t (a, b) FROM stdin;\nx;y\t1\np;q;r\t2\n\\.\nSELECT 2;",
			want:  []string{"COPY t (a, b) FROM stdin;\nx;y\t1\np;q;r\t2\n\\.\n", "SELECT 2;"},
		},
		{
			name:  "lowercase with options",
			input: "copy t from stdin with (format text);\na;b\t1\n\\.\nSELECT 2;",
			want:  []string{"copy t from stdin with (format text);\na;b\t1\n\\.\n", "SELECT 2;"},
		},
		{
			name:  "two copy blocks",
			input: "COPY a FROM stdin;\n1;\n\\.\nCOPY b FROM stdin;\n2;\n\\.\n",
			want:  []string{"COPY a FROM stdin;\n1;\n\\.\n", "COPY b FROM stdin;\n2;\n\\.\n"},
		},
		{
			name:  "missing terminator swallows to EOF like psql",
			input: "COPY t FROM stdin;\na;b\t1\nSELECT 2;",
			want:  []string{"COPY t FROM stdin;\na;b\t1\nSELECT 2;"},
		},
		{
			name:  "terminator only at line start",
			input: "COPY t FROM stdin;\ndata \\. not end\n\\.\nSELECT 2;",
			want:  []string{"COPY t FROM stdin;\ndata \\. not end\n\\.\n", "SELECT 2;"},
		},
		{
			name:  "terminator line with more content is data",
			input: "COPY t FROM stdin;\n\\.x\n\\.\nSELECT 2;",
			want:  []string{"COPY t FROM stdin;\n\\.x\n\\.\n", "SELECT 2;"},
		},
		{
			name:  "terminator at EOF without newline",
			input: "COPY t FROM stdin;\na\t1\n\\.",
			want:  []string{"COPY t FROM stdin;\na\t1\n\\."},
		},
		{
			name:  "CRLF data and terminator",
			input: "COPY t FROM stdin;\r\na;b\t1\r\n\\.\r\nSELECT 2;",
			want:  []string{"COPY t FROM stdin;\r\na;b\t1\r\n\\.\r\n", "SELECT 2;"},
		},
		{
			name:  "copy from file does not enter data mode",
			input: "COPY t FROM 'x.csv'; SELECT 2;",
			want:  []string{"COPY t FROM 'x.csv';", " SELECT 2;"},
		},
		{
			name:  "copy to stdout does not enter data mode",
			input: "COPY t TO stdout; SELECT 2;",
			want:  []string{"COPY t TO stdout;", " SELECT 2;"},
		},
		{
			name:  "relation named stdin inside copy query form does not match",
			input: "COPY (SELECT a FROM stdin) TO stdout; SELECT 2;",
			want:  []string{"COPY (SELECT a FROM stdin) TO stdout;", " SELECT 2;"},
		},
		{
			name:  "non-copy statement mentioning from stdin in string",
			input: "SELECT 'COPY t FROM stdin;'; SELECT 2;",
			want:  []string{"SELECT 'COPY t FROM stdin;';", " SELECT 2;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.input)

			var rebuilt []byte
			for _, s := range segs {
				if s.Text != tt.input[s.ByteStart:s.ByteEnd] {
					t.Fatalf("segment text %q does not match input[%d:%d]", s.Text, s.ByteStart, s.ByteEnd)
				}
				rebuilt = append(rebuilt, s.Text...)
			}
			if string(rebuilt) != tt.input {
				t.Fatalf("segments do not reconstruct input:\ngot  %q\nwant %q", rebuilt, tt.input)
			}

			var got []string
			for _, s := range segs {
				if !s.Empty() {
					got = append(got, s.Text)
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d non-empty segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSplitBeginAtomicKeywordContext pins context-aware keyword counting
// inside BEGIN ATOMIC bodies: dot-qualified references (t.end, t . case),
// AS aliases (AS end), and bare unreserved BEGIN as a column name are
// ordinary identifiers, not block delimiters. All positive shapes
// engine-verified on PostgreSQL 17 (server-side; psql's own client
// scanner mis-splits the dot/AS shapes — see KnownBetterThanPsql).
func TestSplitBeginAtomicKeywordContext(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "dot-qualified end does not close block",
			input: "CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT t4.end FROM t4; END; SELECT 3;",
			want:  []string{"CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT t4.end FROM t4; END;", " SELECT 3;"},
		},
		{
			name:  "spaced dot-qualified end",
			input: "CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT t4 . end FROM t4; END; SELECT 3;",
			want:  []string{"CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT t4 . end FROM t4; END;", " SELECT 3;"},
		},
		{
			name:  "AS end alias does not close block",
			input: "CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT 1 AS end; END; SELECT 3;",
			want:  []string{"CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT 1 AS end; END;", " SELECT 3;"},
		},
		{
			name:  "dot-qualified case does not open block",
			input: "CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT t4.end + t4.case FROM t4; END; SELECT 3;",
			want:  []string{"CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT t4.end + t4.case FROM t4; END;", " SELECT 3;"},
		},
		{
			name:  "bare begin column does not open block",
			input: "CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT begin FROM t4; END; SELECT 3;",
			want:  []string{"CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT begin FROM t4; END;", " SELECT 3;"},
		},
		{
			name:  "case expression still counts",
			input: "CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT CASE WHEN t.end > 0 THEN 1 ELSE 2 END FROM t; END; SELECT 3;",
			want:  []string{"CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT CASE WHEN t.end > 0 THEN 1 ELSE 2 END FROM t; END;", " SELECT 3;"},
		},
		{
			// Lexical contract: nested BEGIN ATOMIC opens a nested block
			// (grammar validity aside); bare BEGIN no longer does.
			name:  "nested begin atomic still counts",
			input: "CREATE PROCEDURE p() BEGIN ATOMIC BEGIN ATOMIC SELECT 1; END; END; SELECT 3;",
			want:  []string{"CREATE PROCEDURE p() BEGIN ATOMIC BEGIN ATOMIC SELECT 1; END; END;", " SELECT 3;"},
		},
		{
			name:  "quoted end identifier does not close block",
			input: `CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT "end" FROM t4; END; SELECT 3;`,
			want:  []string{`CREATE FUNCTION f() RETURNS int BEGIN ATOMIC SELECT "end" FROM t4; END;`, ` SELECT 3;`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.input)

			var rebuilt []byte
			for _, s := range segs {
				if s.Text != tt.input[s.ByteStart:s.ByteEnd] {
					t.Fatalf("segment text %q does not match input[%d:%d]", s.Text, s.ByteStart, s.ByteEnd)
				}
				rebuilt = append(rebuilt, s.Text...)
			}
			if string(rebuilt) != tt.input {
				t.Fatalf("segments do not reconstruct input:\ngot  %q\nwant %q", rebuilt, tt.input)
			}

			var got []string
			for _, s := range segs {
				if !s.Empty() {
					got = append(got, s.Text)
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d non-empty segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSplitPsqlMetacommandLines pins psql metacommand handling: since the
// CVE-2025-8714 point releases (August 2025), pg_dump plain-format output
// brackets the dump in \restrict / \unrestrict lines, so metacommand lines
// are part of the supported input. A top-level line starting with a
// backslash and a letter (after a true newline) is consumed like a comment:
// its content never affects boundaries, and it rides along as trivia of the
// segment it sits in; Empty reports trivia-only segments statement-less.
// Inside strings/dollar-quotes the scanner is in string state first, so
// pseudo-metacommands there are untouched. A metacommand at byte zero of a
// script is glued (no preceding newline) and rescued by the parser instead.
func TestSplitPsqlMetacommandLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string // exact non-empty segment texts
	}{
		{
			name:  "restrict header attaches as leading trivia",
			input: "--\n\\restrict Abc123\nSELECT 1;",
			want:  []string{"--\n\\restrict Abc123\nSELECT 1;"},
		},
		{
			name:  "unrestrict after last statement is trailing trivia",
			input: "SELECT 1;\n\\unrestrict Abc123\n",
			want:  []string{"SELECT 1;"},
		},
		{
			name:  "metacommand between statements attaches forward",
			input: "SELECT 1;\n\\connect other\nSELECT 2;",
			want:  []string{"SELECT 1;", "\n\\connect other\nSELECT 2;"},
		},
		{
			name:  "metacommand content does not affect boundaries",
			input: "SELECT 1;\n\\copy (select 'a;b') to stdout;\nSELECT 2;",
			want:  []string{"SELECT 1;", "\n\\copy (select 'a;b') to stdout;\nSELECT 2;"},
		},
		{
			name:  "line-start backslash inside dollar quote is data",
			input: "SELECT $$x\n\\restrict y$$; SELECT 2;",
			want:  []string{"SELECT $$x\n\\restrict y$$;", " SELECT 2;"},
		},
		{
			name:  "line-start backslash inside E-string is data",
			input: "SELECT E'a\n\\restrict b'; SELECT 2;",
			want:  []string{"SELECT E'a\n\\restrict b';", " SELECT 2;"},
		},
		{
			name:  "copy terminator is not a metacommand",
			input: "COPY t FROM stdin;\n\\N\n\\.\nSELECT 2;",
			want:  []string{"COPY t FROM stdin;\n\\N\n\\.\n", "SELECT 2;"},
		},
		{
			// psql recognizes backslash commands mid-line too (engine-
			// verified); \gset is consumed to end of line as trivia. The
			// buffer-send semantics (\g sends the query) are NOT simulated —
			// the statement runs on to the semicolon, a documented loud
			// divergence for interactive idioms.
			name:  "mid-line metacommand is trivia to end of line",
			input: "SELECT 1 \\gset\n; SELECT 2;",
			want:  []string{"SELECT 1 \\gset\n;", " SELECT 2;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.input)

			var rebuilt []byte
			for _, s := range segs {
				if s.Text != tt.input[s.ByteStart:s.ByteEnd] {
					t.Fatalf("segment text %q does not match input[%d:%d]", s.Text, s.ByteStart, s.ByteEnd)
				}
				rebuilt = append(rebuilt, s.Text...)
			}
			if string(rebuilt) != tt.input {
				t.Fatalf("segments do not reconstruct input:\ngot  %q\nwant %q", rebuilt, tt.input)
			}

			var got []string
			for _, s := range segs {
				if !s.Empty() {
					got = append(got, s.Text)
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d non-empty segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSplitRealPgDump runs the full real pg_dump 17.7 output (generated by
// the tool inside a container, committed verbatim per corpus rule C19 —
// including the \restrict/\unrestrict bracket lines) through Split: the
// segments must reconstruct the input losslessly with every metacommand
// line in a statement-less segment.
func TestSplitRealPgDump(t *testing.T) {
	data, err := os.ReadFile("testdata/pg_dump_17_restrict.sql")
	if err != nil {
		t.Fatalf("read dump: %v", err)
	}
	sql := string(data)
	segs := Split(sql)
	var rebuilt []byte
	for _, s := range segs {
		rebuilt = append(rebuilt, s.Text...)
	}
	if string(rebuilt) != sql {
		t.Fatal("segments do not reconstruct the dump")
	}
	nonEmpty := 0
	for _, s := range segs {
		if !s.Empty() {
			nonEmpty++
		}
	}
	if nonEmpty != 18 {
		t.Fatalf("expected 18 statements in dump, got %d", nonEmpty)
	}
}

// TestSplitCopyAfterMetacommand pins the D5×D6 interaction: a COPY ... FROM
// STDIN statement whose segment carries leading metacommand trivia
// (pg_dumpall emits \connect right before each database's statements) must
// still enter data mode — the metacommand words must not shadow COPY as
// the statement's first word.
func TestSplitCopyAfterMetacommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "restrict then copy",
			input: "SELECT 0;\n\\restrict x\nCOPY t FROM stdin;\na;b\n\\.\n",
			want:  []string{"SELECT 0;", "\n\\restrict x\nCOPY t FROM stdin;\na;b\n\\.\n"},
		},
		{
			// pg_dumpall output has comment headers before \connect, so the
			// metacommand is always newline-preceded in practice.
			name:  "pg_dumpall shape connect then copy",
			input: "-- pg_dumpall\n\\connect mydb\nCOPY t FROM stdin;\np;q\n\\.\nSELECT 1;",
			want:  []string{"-- pg_dumpall\n\\connect mydb\nCOPY t FROM stdin;\np;q\n\\.\n", "SELECT 1;"},
		},
		{
			name:  "consecutive metacommands then copy",
			input: "--\n\\restrict A\n\\connect db\nCOPY t FROM stdin;\nx;y\n\\.\nSELECT 1;",
			want:  []string{"--\n\\restrict A\n\\connect db\nCOPY t FROM stdin;\nx;y\n\\.\n", "SELECT 1;"},
		},
		{
			name:  "metacommand-lookalike inside copy options is not trivia",
			input: "COPY t FROM stdin WITH (DELIMITER ';');\na;b\n\\.\nSELECT 1;",
			want:  []string{"COPY t FROM stdin WITH (DELIMITER ';');\na;b\n\\.\n", "SELECT 1;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.input)

			var rebuilt []byte
			for _, s := range segs {
				if s.Text != tt.input[s.ByteStart:s.ByteEnd] {
					t.Fatalf("segment text %q does not match input[%d:%d]", s.Text, s.ByteStart, s.ByteEnd)
				}
				rebuilt = append(rebuilt, s.Text...)
			}
			if string(rebuilt) != tt.input {
				t.Fatalf("segments do not reconstruct input:\ngot  %q\nwant %q", rebuilt, tt.input)
			}

			var got []string
			for _, s := range segs {
				if !s.Empty() {
					got = append(got, s.Text)
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d non-empty segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSplitCommentSweepFollowups pins the fixes from the post-merge review
// sweep of the D1-D6b series, each engine-verified on PostgreSQL 17.
func TestSplitCommentSweepFollowups(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			// psql strips a leading UTF-8 BOM (mainloop.c) and the COPY
			// data loads; the BOM must not shadow COPY as the first word.
			name:  "BOM then copy enters data mode",
			input: "\xEF\xBB\xBFCOPY t FROM stdin;\na;b\t1\n\\.\nSELECT 2;",
			want:  []string{"\xEF\xBB\xBFCOPY t FROM stdin;\na;b\t1\n\\.\n", "SELECT 2;"},
		},
		{
			// Engine-verified: the server maps FROM STDIN and FROM STDOUT
			// both to a NULL filename and runs copy-in, and psql loads the
			// inline data. Swallowing the data lines IS the engine
			// behavior — do not "fix" this back to a rejection.
			name:  "copy from stdout is copy-in like the engine",
			input: "COPY t FROM STDOUT;\nviastdout\t9\n\\.\nSELECT 2;",
			want:  []string{"COPY t FROM STDOUT;\nviastdout\t9\n\\.\n", "SELECT 2;"},
		},
		{
			// psql buffers same-line SQL after a COPY statement and runs it
			// after the data — unrepresentable in contiguous segments. Data
			// mode disengages so the remainder stays a real statement and
			// the data lines fail loudly instead of being silently
			// swallowed together with the trailing SQL.
			name:  "same-line SQL after copy disables data mode",
			input: "COPY t FROM stdin; SELECT 42 AS marker;\nsameline\t2\n\\.\n",
			want:  []string{"COPY t FROM stdin;", " SELECT 42 AS marker;", "\nsameline\t2\n\\.\n"},
		},
		{
			// psql executes metacommands mid-buffer and the function body
			// continues (engine-verified: \echo END inside BEGIN ATOMIC,
			// function created, body has both statements).
			name:  "metacommand inside begin atomic body is trivia",
			input: "CREATE FUNCTION fm() RETURNS int BEGIN ATOMIC\nSELECT 1;\n\\echo END\nSELECT 2;\nEND; SELECT 3;",
			want:  []string{"CREATE FUNCTION fm() RETURNS int BEGIN ATOMIC\nSELECT 1;\n\\echo END\nSELECT 2;\nEND;", " SELECT 3;"},
		},
		{
			// scan.l: x$BEGIN is one identifier — the BEGIN ATOMIC dispatch
			// must not fire on an identifier tail.
			name:  "identifier tail BEGIN does not open atomic block",
			input: "SELECT 1 AS x$BEGIN ATOMIC; SELECT 2;",
			want:  []string{"SELECT 1 AS x$BEGIN ATOMIC;", " SELECT 2;"},
		},
		{
			name:  "BOM only script",
			input: "\xEF\xBB\xBFSELECT 1;",
			want:  []string{"\xEF\xBB\xBFSELECT 1;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.input)

			var rebuilt []byte
			for _, s := range segs {
				if s.Text != tt.input[s.ByteStart:s.ByteEnd] {
					t.Fatalf("segment text %q does not match input[%d:%d]", s.Text, s.ByteStart, s.ByteEnd)
				}
				rebuilt = append(rebuilt, s.Text...)
			}
			if string(rebuilt) != tt.input {
				t.Fatalf("segments do not reconstruct input:\ngot  %q\nwant %q", rebuilt, tt.input)
			}

			var got []string
			for _, s := range segs {
				if !s.Empty() {
					got = append(got, s.Text)
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d non-empty segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestSplitMetacommandAnywhere pins the position-context-free metacommand
// rule (D7): psql recognizes backslash commands at ANY top-level position,
// not only line starts (engine-verified: SELECT 1; \echo MIDLINE executes,
// mid-line \g buffer-sends), and a top-level backslash is never valid SQL,
// so consuming backslash+letter to end of line can never eat a legal
// statement. Being context-free is what makes re-splitting stable: the two
// line-start rules each broke re-split idempotence at segment offset zero
// in opposite directions.
func TestSplitMetacommandAnywhere(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "byte-zero metacommand is trivia",
			input: "\\restrict abc\nSELECT 1;",
			want:  []string{"\\restrict abc\nSELECT 1;"},
		},
		{
			name:  "byte-zero metacommand then copy enters data mode",
			input: "\\connect db\nCOPY t FROM stdin;\na;b\t1\n\\.\nSELECT 2;",
			want:  []string{"\\connect db\nCOPY t FROM stdin;\na;b\t1\n\\.\n", "SELECT 2;"},
		},
		{
			name:  "mid-line metacommand after semicolon",
			input: "SELECT 1; \\echo mid\nSELECT 2;",
			want:  []string{"SELECT 1;", " \\echo mid\nSELECT 2;"},
		},
		{
			// The lenient line-start rule's counterexample: the quote inside
			// the metacommand line is trivia in BOTH passes — the line is
			// consumed to its newline, and the next line's semicolon splits
			// identically on re-split.
			name:  "mid-line metacommand with quote content",
			input: "SELECT 1;\\echo 'x\n;y';",
			want:  []string{"SELECT 1;", "y';"},
		},
		{
			// The strict rule's counterexample (the D7 case): a byte-zero
			// metacommand line with quotes and semicolons is one trivia-only
			// (statement-less) segment, stable under re-split.
			name:  "byte-zero metacommand with boundary-active content",
			input: "\\col U&'; ;' 1 x +; -- note;\n",
			want:  nil,
		},
		{
			name:  "backslash inside string untouched",
			input: "SELECT 'has \\backslash;inside'; SELECT 2;",
			want:  []string{"SELECT 'has \\backslash;inside';", " SELECT 2;"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.input)

			var rebuilt []byte
			for _, s := range segs {
				if s.Text != tt.input[s.ByteStart:s.ByteEnd] {
					t.Fatalf("segment text %q does not match input[%d:%d]", s.Text, s.ByteStart, s.ByteEnd)
				}
				rebuilt = append(rebuilt, s.Text...)
			}
			if string(rebuilt) != tt.input {
				t.Fatalf("segments do not reconstruct input:\ngot  %q\nwant %q", rebuilt, tt.input)
			}

			// Re-split stability: every segment re-splits to itself.
			for _, s := range segs {
				re := Split(s.Text)
				if len(re) != 1 || re[0].Text != s.Text {
					t.Fatalf("segment %q re-splits into %d segments", s.Text, len(re))
				}
			}

			var got []string
			for _, s := range segs {
				if !s.Empty() {
					got = append(got, s.Text)
				}
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d non-empty segments %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("segment[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestStripTopLevelMetacommands pins the strip walker's byte parity with
// Split: statement anchoring at depth-zero semicolons (a COPY that is not
// the first statement must still enter data mode, so backslash-letter
// lines inside its data are preserved), and only true top-level
// metacommands are removed.
func TestStripTopLevelMetacommands(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			// The review counterexample: COPY preceded by another statement.
			name: "copy after another statement keeps data lines",
			in:   "SELECT 1;\nCOPY t FROM stdin;\na\tb\n\\x in data\n\\.\nSELECT 2;",
			want: "SELECT 1;\nCOPY t FROM stdin;\na\tb\n\\x in data\n\\.\nSELECT 2;",
		},
		{
			// Consecutive COPY blocks: the second anchor depends on
			// stmtStart advancing past the first data block.
			name: "two consecutive copies keep both data blocks",
			in:   "COPY a FROM stdin;\n\\x1\n\\.\nCOPY b FROM stdin;\n\\x2\n\\.\n",
			want: "COPY a FROM stdin;\n\\x1\n\\.\nCOPY b FROM stdin;\n\\x2\n\\.\n",
		},
		{
			name: "metacommand between statements is stripped",
			in:   "SELECT 1;\n\\connect db\nCOPY t FROM stdin;\n\\N\n\\.\n",
			want: "SELECT 1;\nCOPY t FROM stdin;\n\\N\n\\.\n",
		},
		{
			name: "string content preserved",
			in:   "SELECT 'a\n\\restrict b';",
			want: "SELECT 'a\n\\restrict b';",
		},
		{
			name: "rule parens do not advance the anchor",
			in:   "CREATE RULE r AS ON UPDATE TO t DO ALSO (SELECT 1; SELECT 2);\nCOPY t FROM stdin;\n\\x\n\\.\n",
			want: "CREATE RULE r AS ON UPDATE TO t DO ALSO (SELECT 1; SELECT 2);\nCOPY t FROM stdin;\n\\x\n\\.\n",
		},
		{
			name: "top-level metacommand stripped mid-line",
			in:   "SELECT 1; \\echo mid\nSELECT 2;",
			want: "SELECT 1; SELECT 2;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripTopLevelMetacommands(tt.in); got != tt.want {
				t.Fatalf("strip mismatch:\ngot  %q\nwant %q", got, tt.want)
			}
		})
	}
}
