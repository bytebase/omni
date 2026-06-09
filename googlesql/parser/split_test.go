package parser

import (
	"reflect"
	"strings"
	"testing"
)

// runSplitCases asserts Split(input) returns exactly the expected Segment list.
func runSplitCases(t *testing.T, split func(string) []Segment, cases []struct {
	name  string
	input string
	want  []Segment
}) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := split(c.input)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("split(%q) mismatch\n got: %+v\nwant: %+v", c.input, got, c.want)
			}
		})
	}
}

func TestSplit_Empty(t *testing.T) {
	if got := Split(""); got != nil {
		t.Errorf("Split(\"\") = %+v, want nil", got)
	}
	if got := SplitFlat(""); got != nil {
		t.Errorf("SplitFlat(\"\") = %+v, want nil", got)
	}
}

func TestSplit_Basic(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{
			name:  "single_no_semicolon",
			input: "SELECT 1",
			want:  []Segment{{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8}},
		},
		{
			name:  "single_trailing_semicolon",
			input: "SELECT 1;",
			want:  []Segment{{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8}},
		},
		{
			name:  "two_statements",
			input: "SELECT 1; SELECT 2",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8},
				{Text: " SELECT 2", ByteStart: 9, ByteEnd: 18},
			},
		},
		{
			name:  "lone_semicolons_filtered",
			input: ";;;",
			want:  nil,
		},
		{
			name:  "blank_between_filtered",
			input: "SELECT 1;;SELECT 2",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8},
				{Text: "SELECT 2", ByteStart: 10, ByteEnd: 18},
			},
		},
	}
	t.Run("Split", func(t *testing.T) { runSplitCases(t, Split, cases) })
	t.Run("SplitFlat", func(t *testing.T) { runSplitCases(t, SplitFlat, cases) })
}

// Semicolons hidden inside strings/bytes/backtick-identifiers/comments must NOT
// split. Both variants share the lexer, so both behave identically here.
func TestSplit_SemicolonInHiddenContent(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []Segment
	}{
		{
			name:  "semicolon_in_single_quote",
			input: "SELECT ';';SELECT 2",
			want: []Segment{
				{Text: "SELECT ';'", ByteStart: 0, ByteEnd: 10},
				{Text: "SELECT 2", ByteStart: 11, ByteEnd: 19},
			},
		},
		{
			name:  "semicolon_in_triple_quote",
			input: "SELECT '''a;b''';SELECT 2",
			want: []Segment{
				{Text: "SELECT '''a;b'''", ByteStart: 0, ByteEnd: 16},
				{Text: "SELECT 2", ByteStart: 17, ByteEnd: 25},
			},
		},
		{
			name:  "semicolon_in_backtick_ident",
			input: "SELECT `a;b`;SELECT 2",
			want: []Segment{
				{Text: "SELECT `a;b`", ByteStart: 0, ByteEnd: 12},
				{Text: "SELECT 2", ByteStart: 13, ByteEnd: 21},
			},
		},
		{
			name:  "semicolon_in_dash_comment",
			input: "SELECT 1 -- a;b\n;SELECT 2",
			want: []Segment{
				{Text: "SELECT 1 -- a;b\n", ByteStart: 0, ByteEnd: 16},
				{Text: "SELECT 2", ByteStart: 17, ByteEnd: 25},
			},
		},
		{
			name:  "semicolon_in_block_comment",
			input: "SELECT 1 /* a;b */;SELECT 2",
			want: []Segment{
				{Text: "SELECT 1 /* a;b */", ByteStart: 0, ByteEnd: 18},
				{Text: "SELECT 2", ByteStart: 19, ByteEnd: 27},
			},
		},
		{
			name:  "semicolon_in_pound_comment",
			input: "SELECT 1 # a;b\n;SELECT 2",
			want: []Segment{
				{Text: "SELECT 1 # a;b\n", ByteStart: 0, ByteEnd: 15},
				{Text: "SELECT 2", ByteStart: 16, ByteEnd: 24},
			},
		},
	}
	t.Run("Split", func(t *testing.T) { runSplitCases(t, Split, cases) })
	t.Run("SplitFlat", func(t *testing.T) { runSplitCases(t, SplitFlat, cases) })
}

// Block-aware Split keeps procedural BEGIN/END (and IF/CASE/LOOP/WHILE/REPEAT/
// FOR) blocks whole; SplitFlat (BigQuery lexer-split) does not.
// TestSplit_BlockWithFunctionKeyword pins the fix for the block-interior
// function-keyword bug (Codex review): a control-flow KEYWORD used as a scalar
// function inside a block body — IF(...) — must NOT be counted as a nested block
// opener. Counting it inflated the block depth so the ';' after END was swallowed
// and the following statement became trailing junk (Split returned 1, not 2).
func TestSplit_BlockWithFunctionKeyword(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  int
	}{
		{"if_function_in_block", "BEGIN SELECT IF(TRUE, 1, 0); END; SELECT 2;", 2},
		{"case_expr_with_if_function", "BEGIN SELECT CASE WHEN x THEN IF(a,b,c) ELSE 0 END; END; SELECT 2;", 2},
		{"real_nested_if_statement", "BEGIN IF x THEN SELECT 1; END IF; END; SELECT 2;", 2},
		{"plain_block", "BEGIN SELECT 1; END; SELECT 2;", 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Split(c.input); len(got) != c.want {
				t.Fatalf("Split(%q) = %d segments, want %d: %+v", c.input, len(got), c.want, got)
			}
		})
	}
}

func TestSplit_BeginEndBlock(t *testing.T) {
	// A stored-procedure body with internal ';' separators.
	const proc = "CREATE PROCEDURE p() BEGIN SELECT 1; SELECT 2; END"

	t.Run("Split_keeps_block_whole", func(t *testing.T) {
		got := Split(proc)
		if len(got) != 1 {
			t.Fatalf("Split kept %d segments, want 1 (block whole): %+v", len(got), got)
		}
		if got[0].Text != proc {
			t.Errorf("Split segment = %q, want whole proc", got[0].Text)
		}
	})

	t.Run("Split_block_then_next_statement", func(t *testing.T) {
		input := proc + "; SELECT 99"
		got := Split(input)
		if len(got) != 2 {
			t.Fatalf("Split kept %d segments, want 2: %+v", len(got), got)
		}
		if got[0].Text != proc {
			t.Errorf("segment[0] = %q, want whole proc", got[0].Text)
		}
		if strings.TrimSpace(got[1].Text) != "SELECT 99" {
			t.Errorf("segment[1] = %q, want ' SELECT 99'", got[1].Text)
		}
	})

	t.Run("SplitFlat_splits_inside_block", func(t *testing.T) {
		// BigQuery lexer-split splits on every ';', including those inside the
		// procedure body.
		got := SplitFlat(proc)
		if len(got) != 3 {
			t.Fatalf("SplitFlat kept %d segments, want 3 (splits on inner ';'): %+v", len(got), got)
		}
	})

	t.Run("nested_blocks", func(t *testing.T) {
		// IF nested inside BEGIN; depth must balance back to 0 only at the
		// outer END.
		input := "CREATE PROCEDURE p() BEGIN IF x THEN SELECT 1; END IF; SELECT 2; END; SELECT 3"
		got := Split(input)
		if len(got) != 2 {
			t.Fatalf("Split kept %d segments, want 2: %+v", len(got), got)
		}
		if !strings.Contains(got[0].Text, "END IF") || !strings.HasSuffix(strings.TrimSpace(got[0].Text), "END") {
			t.Errorf("segment[0] = %q, want whole nested-block proc", got[0].Text)
		}
		if strings.TrimSpace(got[1].Text) != "SELECT 3" {
			t.Errorf("segment[1] = %q, want ' SELECT 3'", got[1].Text)
		}
	})
}

// Port of the legacy spanner split_begin_end_test.go: a CASE *expression* (no
// internal ';') must split identically under both variants because the
// surrounding query's terminating ';' is the only delimiter.
func TestSplit_CaseExpressionNoInternalSemicolon(t *testing.T) {
	statement := "SELECT\n  CASE\n    WHEN status = 'active' THEN 1\n    WHEN status = 'inactive' THEN 0\n  END\nFROM users;\nSELECT * FROM orders;"

	for _, variant := range []struct {
		name  string
		split func(string) []Segment
	}{
		{"Split", Split},
		{"SplitFlat", SplitFlat},
	} {
		t.Run(variant.name, func(t *testing.T) {
			list := variant.split(statement)
			if len(list) != 2 {
				t.Fatalf("%s produced %d segments, want 2: %+v", variant.name, len(list), list)
			}
			if !strings.Contains(list[0].Text, "CASE") || !strings.Contains(list[0].Text, "END") {
				t.Errorf("segment[0] = %q, want to contain CASE..END", list[0].Text)
			}
			if !strings.Contains(list[1].Text, "SELECT * FROM orders") {
				t.Errorf("segment[1] = %q, want 'SELECT * FROM orders'", list[1].Text)
			}
		})
	}
}

// A top-level BEGIN that starts a transaction (BEGIN; / BEGIN TRANSACTION) must
// NOT be treated as a block opener — otherwise the splitter would swallow the
// rest of the script.
func TestSplit_BeginTransactionIsNotBlock(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  int
	}{
		{"begin_semicolon", "BEGIN; SELECT 1; COMMIT", 3},
		{"begin_transaction", "BEGIN TRANSACTION; SELECT 1; COMMIT", 3},
		// begin_statement: begin_transaction_keywords transaction_mode_list?,
		// where a transaction_mode begins with READ or ISOLATION and the optional
		// TRANSACTION keyword may be omitted. So `BEGIN READ ONLY` / `BEGIN
		// ISOLATION LEVEL ...` are flat transaction starts, NOT procedural block
		// openers (oracle-confirmed: Spanner parses each as BeginStatement). If
		// the splitter misread these as block openers it would swallow the rest
		// of the script into one segment.
		{"begin_read_only", "BEGIN READ ONLY; SELECT 1; COMMIT", 3},
		{"begin_read_write", "BEGIN READ WRITE; SELECT 1; COMMIT", 3},
		{"begin_isolation_level", "BEGIN ISOLATION LEVEL SERIALIZABLE; SELECT 1; COMMIT", 3},
		{"begin_transaction_read_only", "BEGIN TRANSACTION READ ONLY; SELECT 1; COMMIT", 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Split(c.input); len(got) != c.want {
				t.Errorf("Split(%q) = %d segments, want %d: %+v", c.input, len(got), c.want, got)
			}
		})
	}
}

// A typed block closer (END IF / END WHILE / …) at the TOP level — i.e. closing
// the outermost block back to depth 0 — must be recognized as the closer's
// suffix, NOT mistaken for a new top-level block opener. Otherwise the trailing
// statement-terminating ';' would be swallowed into the block's segment. (These
// bare top-level procedural forms are not themselves valid GoogleSQL, but the
// splitter is best-effort and infallible, so its boundary handling must stay
// consistent regardless of nesting depth.)
func TestSplit_TopLevelTypedCloserDoesNotReopen(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"if", "IF x THEN SELECT 1; END IF; SELECT 2"},
		{"while", "WHILE x DO SELECT 1; END WHILE; SELECT 2"},
		{"loop", "LOOP SELECT 1; END LOOP; SELECT 2"},
		{"repeat", "REPEAT SELECT 1; UNTIL x END REPEAT; SELECT 2"},
		{"for", "FOR r IN (SELECT 1) DO SELECT r; END FOR; SELECT 2"},
		{"case", "CASE WHEN x THEN SELECT 1; END CASE; SELECT 2"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Split(c.input)
			if len(got) != 2 {
				t.Fatalf("Split(%q) = %d segments, want 2 (typed closer reopened a block): %+v", c.input, len(got), got)
			}
			if strings.TrimSpace(got[1].Text) != "SELECT 2" {
				t.Errorf("segment[1] = %q, want ' SELECT 2'", got[1].Text)
			}
		})
	}
}

// A control-flow keyword (IF / CASE / LOOP / WHILE / REPEAT / FOR) that appears
// MID-statement — as a function call (IF(...)), a query clause (FOR UPDATE), a
// guard (IF NOT EXISTS), or an expression (CASE … END) — must NOT be mistaken
// for a procedural block opener. Only a STATEMENT-LEADING occurrence opens a
// block. Otherwise the splitter flips into block mode and swallows the next
// statement's terminating ';'. These are all common valid forms and SplitSQL is
// P0 for both BigQuery and Spanner, so a misread here mis-splits real input.
// (Oracle-confirmed against the Spanner emulator: `SELECT IF(x,1,2)`,
// `SELECT a FROM t FOR UPDATE`, `CREATE SCHEMA IF NOT EXISTS s`, and
// `SELECT CASE WHEN x THEN 1 ELSE 2 END` each parse as ONE statement.)
func TestSplit_MidStatementControlKeywordIsNotBlock(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"if_function_call", "SELECT IF(x, 1, 2) AS c; SELECT 1"},
		{"for_update_clause", "SELECT a FROM t FOR UPDATE; SELECT 1"},
		{"if_not_exists_guard", "CREATE SCHEMA IF NOT EXISTS s; SELECT 1"},
		{"case_expression", "SELECT CASE WHEN x THEN 1 ELSE 2 END; SELECT 1"},
		{"while_as_alias", "SELECT 1 AS while; SELECT 2"},
		{"if_in_where", "SELECT 1 WHERE IF(a, b, c); SELECT 2"},
		// A CASE-expression branch may itself contain a control-flow keyword used
		// as an expression (an IF(...) call, a nested CASE). Neither the inner
		// keyword nor the CASE's own END is statement-leading, so the whole CASE
		// expression stays one statement. (Oracle-confirmed: both parse as a
		// single accepted query.) These guard against over-eagerly treating a
		// CASE branch's THEN/ELSE as a procedural statement-list introducer.
		{"if_call_in_case_branch", "SELECT CASE WHEN x THEN IF(a, b, c) ELSE 0 END; SELECT 1"},
		{"nested_case_expression", "SELECT CASE WHEN x THEN CASE WHEN y THEN 1 END END; SELECT 1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Split(c.input)
			if len(got) != 2 {
				t.Fatalf("Split(%q) = %d segments, want 2 (mid-statement control keyword opened a block): %+v", c.input, len(got), got)
			}
			if strings.TrimSpace(got[1].Text) != "SELECT 1" && strings.TrimSpace(got[1].Text) != "SELECT 2" {
				t.Errorf("segment[1] = %q, want the trailing SELECT", got[1].Text)
			}
			// Split must agree with the flat lexer split here: there is no real
			// procedural block, so both see exactly two top-level statements.
			if flat := SplitFlat(c.input); len(flat) != 2 {
				t.Errorf("SplitFlat(%q) = %d segments, want 2: %+v", c.input, len(flat), flat)
			}
		})
	}
}

// A stray top-level typed closer (`END IF` / `END WHILE` / …) that opens no
// preceding block must NOT itself reopen a block and swallow the next
// statement's ';'. (Invalid GoogleSQL, but the splitter is documented
// best-effort/infallible and must never drop a following statement.) This is the
// statement-position rule applied to the closer-suffix keyword: `IF` after a
// stray top-level `END` is not statement-leading, so it cannot open a block.
func TestSplit_StrayTopLevelCloserDoesNotSwallow(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"end_if", "END IF; SELECT 1"},
		{"end_while", "END WHILE; SELECT 1"},
		{"end_loop", "END LOOP; SELECT 1"},
		{"end_for", "END FOR; SELECT 1"},
		{"bare_end", "END; SELECT 1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Split(c.input)
			if len(got) != 2 {
				t.Fatalf("Split(%q) = %d segments, want 2 (stray closer swallowed the next statement): %+v", c.input, len(got), got)
			}
			if strings.TrimSpace(got[1].Text) != "SELECT 1" {
				t.Errorf("segment[1] = %q, want ' SELECT 1'", got[1].Text)
			}
		})
	}
}

// Every block kind's typed closer (END IF / END WHILE / END LOOP / END REPEAT /
// END FOR / bare END) must be recognized so the closer's trailing keyword is not
// miscounted as a NEW block opener (which would desync depth and swallow the
// statement-terminating ';' after the procedure body). Each procedure has an
// internal ';' that must stay inside the single segment.
func TestSplit_TypedClosersBalanceDepth(t *testing.T) {
	cases := []struct {
		name string
		body string // the control block inside CREATE PROCEDURE p() <body>
	}{
		{"if_end_if", "BEGIN IF x THEN SELECT 1; END IF; END"},
		{"while_end_while", "BEGIN WHILE x DO SELECT 1; END WHILE; END"},
		{"loop_end_loop", "BEGIN LOOP SELECT 1; LEAVE; END LOOP; END"},
		{"repeat_end_repeat", "BEGIN REPEAT SELECT 1; UNTIL x END REPEAT; END"},
		{"for_end_for", "BEGIN FOR r IN (SELECT 1) DO SELECT r; END FOR; END"},
		{"nested_if_in_while", "BEGIN WHILE x DO IF y THEN SELECT 1; END IF; END WHILE; END"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// The whole proc is one statement; a trailing ';' then a SELECT is the
			// second. If a typed closer were miscounted, the depth would never
			// return to 0 and the trailing ';' + SELECT would be swallowed.
			input := "CREATE PROCEDURE p() " + c.body + "; SELECT 99"
			got := Split(input)
			if len(got) != 2 {
				t.Fatalf("Split(%q) = %d segments, want 2: %+v", input, len(got), got)
			}
			if strings.TrimSpace(got[1].Text) != "SELECT 99" {
				t.Errorf("segment[1] = %q, want ' SELECT 99' (typed closer desynced depth)", got[1].Text)
			}
		})
	}
}

// On a transaction-control script (BEGIN [TRANSACTION] is flat, not a block),
// Split and SplitFlat must agree segment-for-segment: there is no procedural
// block to keep whole, so the block-aware path degenerates to the flat path.
// Oracle-confirmed: the Spanner emulator parses `BEGIN`, `SELECT 1`, and
// `COMMIT` as three independent statements.
func TestSplit_TCLScriptMatchesFlat(t *testing.T) {
	inputs := []string{
		"BEGIN; SELECT 1; COMMIT",
		"BEGIN TRANSACTION; UPDATE t SET x=1; COMMIT;",
		"START TRANSACTION; INSERT INTO t VALUES (1); ROLLBACK",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			block := Split(in)
			flat := SplitFlat(in)
			if !reflect.DeepEqual(block, flat) {
				t.Errorf("Split and SplitFlat disagree on a TCL script\n input: %q\n Split:     %+v\n SplitFlat: %+v", in, block, flat)
			}
		})
	}
}

// Round-trip: re-joining segment texts with the original delimiters must
// reconstruct the input (the offsets are correct and lossless).
func TestSplit_OffsetsAreLossless(t *testing.T) {
	inputs := []string{
		"SELECT 1; SELECT 2; SELECT 3",
		"  SELECT 1 ;\n SELECT 2 ;",
		"SELECT ';' ; SELECT `a;b`",
		"CREATE PROCEDURE p() BEGIN SELECT 1; END; SELECT 2",
	}
	for _, in := range inputs {
		for _, variant := range []struct {
			name  string
			split func(string) []Segment
		}{{"Split", Split}, {"SplitFlat", SplitFlat}} {
			t.Run(variant.name+"/"+in, func(t *testing.T) {
				for _, seg := range variant.split(in) {
					// Each segment's Text must equal the input slice it claims.
					if seg.Text != in[seg.ByteStart:seg.ByteEnd] {
						t.Errorf("segment Text %q != input[%d:%d] %q",
							seg.Text, seg.ByteStart, seg.ByteEnd, in[seg.ByteStart:seg.ByteEnd])
					}
					// ByteEnd must point at a ';' or len(input) (the delimiter
					// is never part of Text).
					if seg.ByteEnd < len(in) && in[seg.ByteEnd] != ';' {
						t.Errorf("segment ByteEnd %d points at %q, want ';' or EOF", seg.ByteEnd, in[seg.ByteEnd])
					}
				}
			})
		}
	}
}

func TestSegment_Empty(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{"", true},
		{"   ", true},
		{"-- only a comment\n", true},
		{"/* block */", true},
		{"# pound\n", true},
		{"SELECT 1", false},
		{"  x  ", false},
	}
	for _, c := range cases {
		if got := (Segment{Text: c.text}).Empty(); got != c.want {
			t.Errorf("Segment{%q}.Empty() = %v, want %v", c.text, got, c.want)
		}
	}
}
