package deparse

import (
	"math/rand"
	"testing"
)

// TestBufferMatchesCutting runs random writes and nested take-backs against
// buffer and against the plain byte slice it replaced, where a take-back cut
// the byte out, and requires every read to agree after every step. Fragments
// of spaces alone and take-backs over nothing reach the dead-byte skipping
// that deparse output seldom does.
func TestBufferMatchesCutting(t *testing.T) {
	fragments := []string{" ", "  ", "x", " y", "(", "[", ", ", "AB", ""}
	rng := rand.New(rand.NewSource(1))
	for run := 0; run < 2000; run++ {
		var got buffer
		var want []byte
		type scope struct{ mark, start int }
		var open []scope
		for step := 0; step < 60; step++ {
			switch op := rng.Intn(10); {
			case op < 4:
				open = append(open, scope{got.mark(), len(want)})
			case op < 8 && len(open) > 0:
				s := open[len(open)-1]
				open = open[:len(open)-1]
				got.dropLeadingSpace(s.mark)
				if len(want) > s.start && want[s.start] == ' ' {
					want = append(want[:s.start], want[s.start+1:]...)
				}
			default:
				f := fragments[rng.Intn(len(fragments))]
				if len(f) == 1 && rng.Intn(2) == 0 {
					got.WriteByte(f[0])
				} else {
					got.WriteString(f)
				}
				want = append(want, f...)
			}
			if got.String() != string(want) || got.Len() != len(want) {
				t.Fatalf("run %d step %d: buffer holds %q (Len %d), cutting gives %q", run, step, got.String(), got.Len(), want)
			}
			last, ok := got.lastByte()
			if ok != (len(want) > 0) || ok && last != want[len(want)-1] {
				t.Fatalf("run %d step %d: lastByte = %q, %v for %q", run, step, last, ok, want)
			}
		}
	}
}
