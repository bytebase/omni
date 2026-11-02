package parser

import (
	"math"
	"strings"
)

// execCloseMemo keeps the scans for executable comments' closes linear in
// total, without changing what any scan returns.
//
// execCommentClose is a pure function of where it starts: it steps over one
// token at a time and stops at the first "*/" outside a token, or at the end
// of the input. A scan that stands between tokens on an offset an earlier scan
// stood on therefore ends exactly where that one did. Without a memo, input in
// which no comment closes, such as "SELECT 1 " followed by "/*! " many times,
// rescans to the end of the input once per opener, which is quadratic.
//
// The memo is built only after the scans over one input have walked more than
// twice its length (execMemoAfter), which ordinary SQL never does, so it costs
// nothing there.
type execCloseMemo struct {
	// from[p], when non-zero, is the result+2 of the scan from offset p: 1 for
	// a comment that never closes, otherwise the offset of its "*/" plus 2.
	from []int32
	// trail holds the offsets the current scan stood on between tokens, to be
	// recorded once it ends.
	trail []int32
	// No "*/" starts in [rawFrom, rawAt), and rawAt is where the next one
	// starts, or -1 when none does before the end of the input.
	rawFrom, rawAt int
}

// execMemoAfter reports whether scans that have walked walked bytes of an
// input of length n should switch to a memo, and whether one can index it.
func execMemoAfter(walked, n int) bool {
	return walked > 2*n && n < math.MaxInt32-2
}

func newExecCloseMemo(n int) *execCloseMemo {
	return &execCloseMemo{from: make([]int32, n), rawFrom: math.MaxInt}
}

// begin starts a scan.
func (m *execCloseMemo) begin() {
	m.trail = m.trail[:0]
}

// visit is called when a scan stands on pos between tokens. It returns the
// recorded result of the scan from pos and true when there is one; otherwise
// it adds pos to the trail.
func (m *execCloseMemo) visit(pos int) (int, bool) {
	if v := m.from[pos]; v != 0 {
		return int(v) - 2, true
	}
	m.trail = append(m.trail, int32(pos))
	return 0, false
}

// end records res as the result of the scan from every offset on the trail.
func (m *execCloseMemo) end(res int) {
	v := int32(res + 2)
	for _, p := range m.trail {
		m.from[p] = v
	}
}

// rawClose returns the offset of the first "*/" at or after from, or -1. The
// scan for an ordinary comment nested in an executable one ends there.
func (m *execCloseMemo) rawClose(input string, from int) int {
	if from >= m.rawFrom && (m.rawAt < 0 || from <= m.rawAt) {
		return m.rawAt
	}
	at := strings.Index(input[from:], "*/")
	if at >= 0 {
		at += from
	}
	m.rawFrom, m.rawAt = from, at
	return at
}
