package pg

import (
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
