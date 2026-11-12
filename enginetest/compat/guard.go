package compat

import "fmt"

// Guard compares a compat.tsv before and after a change and returns what
// the change may not do. A test run sees only the file it has; these rules
// need the file it replaced.
//
//   - An omni cell changes only from ?, or to its row's expected verdict:
//     a regression is never committed as omni's new verdict.
//   - An expected cell (and its source) changes only from ?, or when the
//     file was recorded at a new pin: the judge's answers are not edited to
//     pass.
//
// The pin is the file's own # pin: line, which a recording writes and
// TestCompatibility holds to the engine's pins.tsv. An edit to pins.tsv
// that leaves the judge's pin alone, such as a comment, unlocks nothing.
// Rows added or removed are not compared: the test regenerates the
// generated ones, and a new hand-added row is judged by a recording.
func Guard(before, after *File) []string {
	repinned := PinChanged(before, after)
	old := map[rowKey]Row{}
	for _, r := range before.Rows {
		old[r.key()] = r
	}
	var problems []string
	for _, r := range after.Rows {
		o, ok := old[r.key()]
		if !ok {
			continue
		}
		if r.Omni != o.Omni && o.Omni != Unjudged && r.Omni != r.Expected {
			problems = append(problems, fmt.Sprintf("%s: omni changed from %s to %s, which is not the expected %s: a regression, not a closed gap\n    %s",
				r.Rule, o.Omni, r.Omni, r.Expected, r.Statement))
		}
		if (r.Expected != o.Expected || r.Source != o.Source) && o.Expected != Unjudged && !repinned {
			problems = append(problems, fmt.Sprintf("%s: expected changed from %s (%s) to %s (%s) without a pin bump; only a recording at a new pin changes it\n    %s",
				r.Rule, o.Expected, o.Source, r.Expected, r.Source, r.Statement))
		}
	}
	return problems
}

// PinChanged reports whether after was recorded at a different pin than
// before: a pin bump.
func PinChanged(before, after *File) bool { return before.Pin != after.Pin }
