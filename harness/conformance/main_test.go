package main

import "testing"

// TestWriteCommittedBoard pins the board-destination decision: a label-only
// run (no adjudication, no explicit opt-in) must never touch the committed
// scoreboards dir — a quickstart `go run .` overwriting the committed
// adjudicated board is exactly the accident this guards against.
func TestWriteCommittedBoard(t *testing.T) {
	cases := []struct {
		name                    string
		adjudicated, writeOptIn bool
		wantCommitted           bool
	}{
		{"label-only quickstart -> out only", false, false, false},
		{"adjudicated -> committed (default unchanged)", true, false, true},
		{"adjudicated with redundant opt-in -> committed", true, true, true},
		{"explicit label-only baseline -> committed", false, true, true},
	}
	for _, c := range cases {
		if got := writeCommittedBoard(c.adjudicated, c.writeOptIn); got != c.wantCommitted {
			t.Errorf("%s: writeCommittedBoard(%v, %v) = %v, want %v",
				c.name, c.adjudicated, c.writeOptIn, got, c.wantCommitted)
		}
	}
}
