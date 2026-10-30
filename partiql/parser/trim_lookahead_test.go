package parser

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

// scanTrimArgs is the lookahead trimHasTopLevelFrom made before it ran the
// lookaheads of nested TRIMs in one pass: it reads one argument list, from
// p.cur, token by token.
func scanTrimArgs(p *Parser) bool {
	snapshot := p.save()
	defer p.restore(snapshot)
	depth := 0
	for {
		switch p.cur.Type {
		case tokEOF:
			return false
		case tokFROM:
			if depth == 0 {
				return true
			}
		case tokPAREN_LEFT, tokBRACKET_LEFT, tokBRACE_LEFT, tokANGLE_DOUBLE_LEFT:
			depth++
		case tokPAREN_RIGHT:
			if depth == 0 {
				return false
			}
			depth--
		case tokBRACKET_RIGHT, tokBRACE_RIGHT, tokANGLE_DOUBLE_RIGHT:
			if depth > 0 {
				depth--
			}
		}
		p.advance()
	}
}

// trimParens returns the offsets of the '(' that follow a TRIM in sql.
func trimParens(sql string) []int {
	var out []int
	lx := NewLexer(sql)
	prevTRIM := false
	for {
		tok := lx.Next()
		if tok.Type == tokEOF {
			return out
		}
		if tok.Type == tokPAREN_LEFT && prevTRIM {
			out = append(out, tok.Loc.Start)
		}
		prevTRIM = tok.Type == tokTRIM
	}
}

// parserAfter returns a parser of sql whose last token read is the '(' at
// offset open.
func parserAfter(t *testing.T, sql string, open int) *Parser {
	p := NewParser(sql)
	for p.prev.Type != tokPAREN_LEFT || p.prev.Loc.Start != open {
		if p.cur.Type == tokEOF {
			t.Fatalf("no '(' at %d in %q", open, sql)
		}
		p.advance()
	}
	return p
}

// TestTrimLookaheadMatchesScan: the lookahead that runs the lookaheads of
// nested TRIMs in one pass gives, for every TRIM, the result of reading that
// TRIM's argument list on its own, whichever TRIM's lookahead runs first.
// The inputs mix every delimiter the lookahead counts, unbalanced ones
// included, since a ] } or >> with nothing open is skipped by one lookahead
// and closes a group for another.
func TestTrimLookaheadMatchesScan(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	pieces := []string{"TRIM(", "TRIM(", "TRIM(", "TRIM", "(", ")", ")", "[", "]", "{", "}", "<<", ">>", " FROM ", "x", ",", "'a'", "/*"}
	for i := 0; i < 20000; i++ {
		var b strings.Builder
		for n := r.Intn(40); n > 0; n-- {
			b.WriteString(pieces[r.Intn(len(pieces))])
			b.WriteByte(' ')
		}
		sql := b.String()
		opens := trimParens(sql)
		want := make(map[int]bool, len(opens))
		for _, o := range opens {
			want[o] = scanTrimArgs(parserAfter(t, sql, o))
		}
		// In the order the parser reaches them, on one parser.
		if len(opens) > 0 {
			p := NewParser(sql)
			for _, o := range opens {
				for p.prev.Type != tokPAREN_LEFT || p.prev.Loc.Start != o {
					p.advance()
				}
				if got := p.trimHasTopLevelFrom(o); got != want[o] {
					t.Fatalf("%q: TRIM at %d: lookahead %v, want %v", sql, o, got, want[o])
				}
			}
		}
		// Innermost first, sharing the records.
		records := map[int]bool{}
		for k := len(opens) - 1; k >= 0; k-- {
			p := parserAfter(t, sql, opens[k])
			p.trimFroms = records
			if got := p.trimHasTopLevelFrom(opens[k]); got != want[opens[k]] {
				t.Fatalf("%q: TRIM at %d, innermost first: lookahead %v, want %v", sql, opens[k], got, want[opens[k]])
			}
		}
	}
}

// TestTrimLookaheadPastNestingLimit: the one pass runs the lookaheads of at
// most maxNestingDepth nested TRIMs at a time. One nested more deeply, which
// the parser cannot reach, runs its own when asked.
func TestTrimLookaheadPastNestingLimit(t *testing.T) {
	const n = maxNestingDepth + 5
	sql := strings.Repeat("TRIM(", n) + "x" + strings.Repeat(" FROM y)", n)
	opens := trimParens(sql)
	p := parserAfter(t, sql, opens[0])
	if !p.trimHasTopLevelFrom(opens[0]) {
		t.Fatalf("outermost TRIM: no FROM found")
	}
	if len(p.trimFroms) != maxNestingDepth {
		t.Errorf("%d results recorded, want %d", len(p.trimFroms), maxNestingDepth)
	}
	for _, o := range opens[maxNestingDepth-2:] {
		q := parserAfter(t, sql, o)
		q.trimFroms = p.trimFroms
		if got, want := q.trimHasTopLevelFrom(o), scanTrimArgs(parserAfter(t, sql, o)); got != want {
			t.Errorf("TRIM at %d: lookahead %v, want %v", o, got, want)
		}
	}
}

// TestNestedTrimStaysLinear: each TRIM looked ahead through its argument
// list, the TRIMs nested in it included, to find a FROM at its own level, so
// TRIM(TRIM(...TRIM(x)...)) took time quadratic in the depth: 3,000 levels
// took about a second, ten of them in a select list about ten. One pass now
// serves every nested TRIM. The shapes run against a deadline, as the
// quadratic ones would take far longer than it.
func TestNestedTrimStaysLinear(t *testing.T) {
	// Each level costs three levels of maxNestingDepth, a bracketed one
	// six.
	const depth, copies = 3000, 10
	nest := func(n int, open, mid, closing string) string {
		return strings.Repeat(open, n) + mid + strings.Repeat(closing, n)
	}
	list := func(item string, last string) string {
		items := make([]string, copies)
		for i := range items {
			items[i] = item
		}
		items[copies-1] = last
		return "SELECT " + strings.Join(items, ", ") + " FROM t"
	}
	for _, tc := range []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"targets", list(nest(depth, "TRIM(", "x", ")"), nest(depth, "TRIM(", "x", ")")), false},
		{"FROM forms", list(nest(depth, "TRIM(", "x", " FROM y)"), nest(depth, "TRIM(", "x", " FROM y)")), false},
		{"modifier forms", list(nest(depth, "TRIM(BOTH ", "x", " FROM y)"), nest(depth, "TRIM(BOTH ", "x", " FROM y)")), false},
		{"bracketed", list(nest(depth/2, "TRIM([", "x", "])"), nest(depth/2, "TRIM([", "x", "])")), false},
		{"failing innermost", list(nest(depth, "TRIM(", "x", ")"), nest(depth, "TRIM(", "x +", ")")), true},
		{"unclosed", "SELECT " + strings.Repeat("TRIM(", depth) + strings.Repeat("x, ", 100000), true},
	} {
		done := make(chan error, 1)
		go func() {
			_, err := Parse(tc.sql)
			done <- err
		}()
		select {
		case err := <-done:
			if (err != nil) != tc.wantErr || isNestingError(err) {
				t.Fatalf("%s: Parse error = %v, want error %v", tc.name, err, tc.wantErr)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: Parse took over 2s", tc.name)
		}
	}
}
