package parser

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

// execCloseShapes are inputs whose executable comments, scanned one at a time
// for their closes, rescan the same text once per opener. Each takes
// milliseconds when the scans are linear in total and seconds otherwise.
var execCloseShapes = []struct {
	name string
	sql  func(n int) string
}{
	{"unclosed openers", func(n int) string { return "SELECT 1 " + strings.Repeat("/*! ", n) }},
	{"unclosed versioned openers", func(n int) string { return "SELECT 1 " + strings.Repeat("/*!50000 ", n) }},
	{"unclosed openers holding strings", func(n int) string { return "SELECT 1 " + strings.Repeat("/*! 'x' ", n) }},
	{"openers sharing one close", func(n int) string {
		return "SELECT 1 " + strings.Repeat("/*! ", n) + "*/" + strings.Repeat("x ", n)
	}},
	{"nested openers", func(n int) string {
		return "SELECT 1 " + strings.Repeat("/*! ", n) + "1" + strings.Repeat(" */", n)
	}},
	{"closed comments", func(n int) string { return "SELECT 1 " + strings.Repeat("/*! + 1 */ ", n) }},
	{"unclosed TiDB openers", func(n int) string { return "SELECT 1 " + strings.Repeat("/*T! ", n) }},
	{"unclosed TiDB feature gates", func(n int) string {
		return "SELECT 1 " + strings.Repeat("/*T![clustered_index] ", n)
	}},
	{"unclosed openers between statements", func(n int) string {
		return strings.Repeat("SELECT 1 /*! ; ", n)
	}},
}

// TestExecCommentCloseStaysLinear guards the memo behind execCommentClose and
// the linked execStack: without them, 20,000 executable comments that never
// close took hundreds of milliseconds to lex and twice as many four times as
// long, and 20,000 nested ones took seconds.
// The budget below is CPU time, not clock time.
func TestExecCommentCloseStaysLinear(t *testing.T) {
	const n = 20000
	for _, tc := range execCloseShapes {
		sql := tc.sql(n)
		start := cpuTime(t)
		_, _ = Parse(sql)
		_ = Tokenize(sql)
		_ = Split(sql)
		if elapsed := cpuSince(t, start); elapsed > 2*time.Second {
			t.Errorf("%s: Parse, Tokenize and Split over %d comments took %v of CPU time", tc.name, n, elapsed)
		}
	}
}

// TestExecCommentCloseMemoAgrees holds the memo to the scan it replaces: from
// every offset of random comment-heavy input, asked in input order (as the
// lexer asks) and in random order, the memoized scan returns what the plain
// scan returns.
func TestExecCommentCloseMemoAgrees(t *testing.T) {
	pieces := []string{
		"/*!", "/*!50000", "/*T!", "/*T![clustered_index]", "/*", "*/", "*", "/", "'", "''", `"`, "`", `\`,
		"--", "-- ", "#", "\n", " ", "x", ";", "1",
	}
	rng := rand.New(rand.NewSource(1))
	for iter := 0; iter < 20000; iter++ {
		var sb strings.Builder
		for k := rng.Intn(40); k >= 0; k-- {
			sb.WriteString(pieces[rng.Intn(len(pieces))])
		}
		input := sb.String()
		want := make([]int, len(input))
		for p := range input {
			want[p] = scanExecClose(input, p, nil)
		}
		inOrder := &Lexer{input: input, execMemo: newExecCloseMemo(len(input))}
		for p := range input {
			if got := inOrder.execCommentClose(p); got != want[p] {
				t.Fatalf("%q from %d in order: got %d, want %d", input, p, got, want[p])
			}
		}
		shuffled := &Lexer{input: input, execMemo: newExecCloseMemo(len(input))}
		for _, p := range rng.Perm(len(input)) {
			if got := shuffled.execCommentClose(p); got != want[p] {
				t.Fatalf("%q from %d shuffled: got %d, want %d", input, p, got, want[p])
			}
		}
	}
}

// TestExecCommentCloseMemoAfterWalk: the memo is built only once the scans
// over one input have walked more than twice its length, so input whose
// comments close keeps the plain scan.
func TestExecCommentCloseMemoAfterWalk(t *testing.T) {
	l := NewLexer("SELECT 1 " + strings.Repeat("/*! + 1 */ ", 1000))
	for l.NextToken().Type != tokEOF {
	}
	if l.execMemo != nil {
		t.Errorf("closed comments built a memo after walking %d bytes of %d", l.execWalked, len(l.input))
	}
	l = NewLexer("SELECT 1 " + strings.Repeat("/*! ", 1000))
	for l.NextToken().Type != tokEOF {
	}
	if l.execMemo == nil {
		t.Errorf("unclosed comments walked %d bytes of %d without a memo", l.execWalked, len(l.input))
	}
}

// TestSplitExecCommentMemoAgrees holds Split's memo to the scan it replaces:
// from every block comment opener of random comment-heavy input, asked in
// input order and in random order, the memoized scan finds the close the plain
// scan finds.
func TestSplitExecCommentMemoAgrees(t *testing.T) {
	pieces := []string{
		"/*!", "/*!50000", "/*T!", "/*T![clustered_index]", "/*", "*/", "*", "/",
		"'", "''", `"`, "`", `\`, "--", "-- ", "#", "\n", " ", "x", ";", "1",
	}
	rng := rand.New(rand.NewSource(1))
	for iter := 0; iter < 20000; iter++ {
		var sb strings.Builder
		for k := rng.Intn(40); k >= 0; k-- {
			sb.WriteString(pieces[rng.Intn(len(pieces))])
		}
		input := sb.String()
		var openers []int
		for p := 0; p+1 < len(input); p++ {
			if input[p] == '/' && input[p+1] == '*' {
				openers = append(openers, p)
			}
		}
		type verdict struct {
			end    int
			closed bool
		}
		want := map[int]verdict{}
		for _, p := range openers {
			end, closed := scanExecutableComment(input, p, nil)
			if !closed {
				end = len(input)
			}
			want[p] = verdict{end, closed}
		}
		for _, order := range [][]int{openers, shuffled(rng, openers)} {
			closes := &splitCloses{memo: newExecCloseMemo(len(input))}
			for _, p := range order {
				end, closed := skipExecutableComment(input, p, closes)
				if !closed {
					end = len(input)
				}
				if got := (verdict{end, closed}); got != want[p] {
					t.Fatalf("%q from %d: got %+v, want %+v", input, p, got, want[p])
				}
			}
		}
	}
}

func shuffled(rng *rand.Rand, xs []int) []int {
	out := make([]int, len(xs))
	for i, j := range rng.Perm(len(xs)) {
		out[i] = xs[j]
	}
	return out
}
