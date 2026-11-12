package grammar

import (
	"fmt"
	"regexp"
	"sort"
)

// Near-misses break an alternative's shortest statement in exactly one
// place: drop a token, cut the end off, insert a stray word, or swap a
// keyword for another. Each is one step from valid, where a lenient
// recursive-descent parser is likeliest to accept what the engine rejects.
// Some land on valid text anyway; the judge decides every row.
//
// The mutation targets the body, the alternative's own tokens, so the break
// is in the production the row names and not in the context that reaches it.

// mutationSuffix matches what nearMisses appends to an alternative's label.
// A label can hold a ! of its own, from a literal such as '!' or "!=", so a
// near-miss is known by the whole suffix.
var mutationSuffix = regexp.MustCompile(`!(drop-[0-9]+|truncate|insert|swap)$`)

// IsNearMiss reports whether a generated row's label names a near-miss
// rather than an alternative.
func IsNearMiss(label string) bool { return mutationSuffix.MatchString(label) }

// strayToken is the word an insert mutation adds. A bare identifier where
// the grammar expects none is a clean break, and the lexer never merges it
// into a neighbouring token.
const strayToken = "zzz"

// keywordText matches a terminal sample that is a keyword or a builtin
// function name: an all-caps word. Literal samples (ident, 's', 1, 0x1) and
// operators are never all-caps, so only keywords are swapped.
var keywordText = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// keywordSet returns the sorted set of distinct keyword samples, so a swap
// picks a deterministic neighbour that still lexes as a keyword rather than
// degrading into an identifier.
func (g *Grammar) keywordSet() []string {
	seen := map[string]bool{}
	for _, t := range g.termText {
		if keywordText.MatchString(t) {
			seen[t] = true
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// neighbourKeyword returns the keyword after k in the sorted set,
// cyclically, and never k itself. It returns "" when k is not a keyword.
func neighbourKeyword(set []string, k string) string {
	i := sort.SearchStrings(set, k)
	if i >= len(set) || set[i] != k || len(set) < 2 {
		return ""
	}
	return set[(i+1)%len(set)]
}

// nearMisses returns the near-misses of one alternative in a fixed order,
// labelled <alternative>!<mutation>. render writes a token list as text.
func nearMisses(a builtAlt, kw []string, render func([]string) string) []Row {
	var out []Row
	emit := func(mut string, toks []string) {
		out = append(out, Row{Rule: a.label + "!" + mut, Statement: render(toks)})
	}

	// Drop each body token in turn.
	for i := range a.body {
		dropped := make([]string, 0, len(a.body)-1)
		dropped = append(dropped, a.body[:i]...)
		dropped = append(dropped, a.body[i+1:]...)
		emit(fmt.Sprintf("drop-%d", i+1), join(a.prefix, dropped, a.suffix))
	}

	// Cut the statement before the body's last token, dropping the suffix
	// too, so it ends early.
	if len(a.body) >= 1 {
		emit("truncate", join(a.prefix, a.body[:len(a.body)-1], nil))
	}

	// Insert a stray word where the production should be complete.
	emit("insert", join(a.prefix, append(append([]string(nil), a.body...), strayToken), a.suffix))

	// Swap the body's first keyword for the next one in the set.
	for i, t := range a.body {
		if n := neighbourKeyword(kw, t); n != "" {
			swapped := append([]string(nil), a.body...)
			swapped[i] = n
			emit("swap", join(a.prefix, swapped, a.suffix))
			break
		}
	}
	return out
}
