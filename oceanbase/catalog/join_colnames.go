package catalog

import "strings"

// joinColNames holds the column names of a join under analysis in
// arr[lo:hi]. Every join RTE lists the columns of both of its inputs, so a
// chain of n joins lists n inputs' worth of names at its top, and copying the
// accumulated names at every level kept a quadratic number of them alive: 3.2
// GB at 20,000 joins. A join instead extends the buffer of its larger input in
// place, appending the other input's names behind it or putting them in front
// of it, and each join RTE keeps a view into the shared array.
//
// Only a join's own result is ever extended, and only by the one join it is an
// input of. Every view into an array is a range inside the range of the
// buffer that owns the array, so writing outside [lo, hi) changes no view
// handed out before. Base relations' names belong to their RTEs and are copied
// in, never extended.
type joinColNames struct {
	arr    []string
	lo, hi int
	// counts maps each lowered name in arr[lo:hi] to its number of
	// occurrences, for the USING fast path; nil until that path first asks.
	counts map[string]int
}

// newJoinColNames returns a buffer holding names, which the caller hands
// over. It has no room to spare until a join extends it.
func newJoinColNames(names []string) *joinColNames {
	return &joinColNames{arr: names, hi: len(names)}
}

// names returns the view a join RTE keeps. Its capacity ends at its length,
// so a caller that appends to it copies instead of writing into the buffer.
func (j *joinColNames) names() []string {
	return j.arr[j.lo:j.hi:j.hi]
}

// add appends names behind the buffer's names, or with front puts them in
// front.
func (j *joinColNames) add(names []string, front bool) {
	if front {
		if j.lo < len(names) {
			j.grow(len(names), 0)
		}
		j.lo -= len(names)
		copy(j.arr[j.lo:], names)
	} else {
		if len(j.arr)-j.hi < len(names) {
			j.grow(0, len(names))
		}
		copy(j.arr[j.hi:], names)
		j.hi += len(names)
	}
	if j.counts != nil {
		for _, name := range names {
			j.counts[strings.ToLower(name)]++
		}
	}
}

// grow moves the names to a new array with room for front more names before
// them and back more after, plus as much again split over both ends, so a run
// of additions to one end moves the names a logarithmic number of times. The
// old array stays behind for the views handed out before.
func (j *joinColNames) grow(front, back int) {
	n := j.hi - j.lo
	room := n + front + back
	arr := make([]string, 2*room+8)
	lo := (len(arr) - n) / 2
	if lo < front {
		lo = front
	}
	copy(arr[lo:], j.arr[j.lo:j.hi])
	j.arr, j.lo, j.hi = arr, lo, lo+n
}

// lowerCounts returns counts, filling it on first use.
func (j *joinColNames) lowerCounts() map[string]int {
	if j.counts == nil {
		j.counts = make(map[string]int, j.hi-j.lo)
		for _, name := range j.arr[j.lo:j.hi] {
			j.counts[strings.ToLower(name)]++
		}
	}
	return j.counts
}

// joinColumnNames returns the coalesced column names of a join whose inputs
// hold leftCols and rightCols, with the inputs' own buffers when they are
// joins themselves (nil for a base relation). The names are those of
// buildCoalescedColNames; this only decides whether an input's buffer can hold
// them without copying it.
func joinColumnNames(left *joinColNames, leftCols []string, right *joinColNames, rightCols []string, usingCols []string) *joinColNames {
	if len(usingCols) == 0 {
		// All the left names, then all the right names: extend the larger
		// input when it is a join.
		if right != nil && (left == nil || len(rightCols) > len(leftCols)) {
			right.add(leftCols, true)
			return right
		}
		if left == nil {
			names := make([]string, 0, len(leftCols)+len(rightCols))
			return newJoinColNames(append(append(names, leftCols...), rightCols...))
		}
		left.add(rightCols, false)
		return left
	}
	// USING columns first, then the other left names, then the other right
	// names. When the left input is a join that already begins with exactly
	// these USING columns, each named once in it, the first two parts are the
	// left names as they stand: the shape of every join in a chain
	// `t JOIN t USING (a) JOIN t USING (a) ...`.
	if left != nil && usingPrefixOf(left, leftCols, usingCols) {
		usingSet := make(map[string]bool, len(usingCols))
		for _, c := range usingCols {
			usingSet[strings.ToLower(c)] = true
		}
		rest := make([]string, 0, len(rightCols))
		for _, c := range rightCols {
			if !usingSet[strings.ToLower(c)] {
				rest = append(rest, c)
			}
		}
		left.add(rest, false)
		return left
	}
	return newJoinColNames(buildCoalescedColNames(leftCols, rightCols, usingCols))
}

// usingPrefixOf reports whether leftCols, the names in left, begin with
// usingCols spelled the same, and name each of them nowhere else.
func usingPrefixOf(left *joinColNames, leftCols, usingCols []string) bool {
	if len(usingCols) > len(leftCols) {
		return false
	}
	for i, c := range usingCols {
		if leftCols[i] != c {
			return false
		}
	}
	counts := left.lowerCounts()
	for _, c := range usingCols {
		if counts[strings.ToLower(c)] != 1 {
			return false
		}
	}
	return true
}
