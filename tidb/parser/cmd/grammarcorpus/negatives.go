package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// The negatives corpus is the grammar-guided half of the rejection suite
// (design doc §3.6). For each reachable alternative's shortest valid sentence,
// it derives near-misses that break the production in exactly one place: drop a
// required token, truncate the production, insert a stray token, or swap a
// keyword for a neighbour. Each is one step from valid, where a lenient
// recursive-descent parser is likeliest to over-accept. The expected verdict is
// reject/reject; a row omni accepts and the server rejects is an over-acceptance
// bug, to fix across the mysql/mariadb forks like a gap.
//
// The mutation targets the body — the alternative's own tokens — so the break
// is in the production the row names, not in the prefix/suffix that reaches it.

// strayToken is the extra identifier an "insert" mutation drops into a sentence.
// An identifier is the least likely token to be absorbed into a neighbouring one
// by the lexer, and a bareword where the grammar expects none is a clean break.
const strayToken = "zzz"

// keywordText matches a terminal sample that is a keyword or a builtin-function
// name: an all-caps word. The literal-class samples (ident, 's', 1, 0x1) and the
// operators are never all-caps, so only genuine keywords are swapped.
var keywordText = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// keywordSet returns the sorted set of distinct keyword samples in the grammar,
// so a swap can pick a deterministic neighbour that is still a real keyword (and
// so still lexes as one, rather than degrading into an identifier).
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

// neighbourKeyword returns the keyword after k in the sorted set, cyclically,
// and never k itself. With fewer than two keywords it returns "" (no swap).
func neighbourKeyword(set []string, k string) string {
	i := sort.SearchStrings(set, k)
	if i >= len(set) || set[i] != k || len(set) < 2 {
		return ""
	}
	return set[(i+1)%len(set)]
}

// generateNegatives derives the negative corpus from the same reachable
// alternatives the positive corpus covers. Rows are deduplicated by statement,
// and a mutation that reproduces a valid sentence verbatim is dropped: the
// negatives corpus holds only intended near-misses, and a measured verdict
// decides whether each truly rejects.
func (g *Grammar) generateNegatives() (rows []Row, uncovered []string) {
	alts, uncovered := g.buildAlts()
	kw := g.keywordSet()

	positive := map[string]bool{}
	for _, a := range alts {
		positive[a.full()] = true
	}

	seen := map[string]bool{}
	for _, a := range alts {
		for _, nr := range mutateAlt(a, kw) {
			if nr.stmt == "" || positive[nr.stmt] || seen[nr.stmt] {
				continue
			}
			seen[nr.stmt] = true
			rows = append(rows, nr)
		}
	}
	return rows, uncovered
}

// mutateAlt returns the near-misses of one alternative, in a deterministic
// order, labelled Rule#alt/mutation.
func mutateAlt(a builtAlt, kw []string) []Row {
	var out []Row
	emit := func(mut string, toks []string) {
		out = append(out, Row{label: a.label + "/" + mut, stmt: spaceJoin(toks)})
	}

	// drop a required token: remove each body token in turn.
	for i := range a.body {
		dropped := make([]string, 0, len(a.body)-1)
		dropped = append(dropped, a.body[:i]...)
		dropped = append(dropped, a.body[i+1:]...)
		emit(fmt.Sprintf("drop-%d", i+1), join(a.prefix, dropped, a.suffix))
	}

	// truncate: cut the sentence after the body's second-to-last token, dropping
	// the final body token and the whole suffix, so it ends early.
	if len(a.body) >= 1 {
		emit("truncate", join(a.prefix, a.body[:len(a.body)-1], nil))
	}

	// insert a stray token: an extra identifier at the body/suffix boundary,
	// where the grammar expects the production to be complete.
	emit("insert", join(a.prefix, append(cloneToks(a.body), strayToken), a.suffix))

	// swap a keyword for a neighbour: replace the first body keyword with the
	// next keyword in the set.
	for i, t := range a.body {
		if n := neighbourKeyword(kw, t); n != "" {
			swapped := cloneToks(a.body)
			swapped[i] = n
			emit("swap", join(a.prefix, swapped, a.suffix))
			break
		}
	}
	return out
}

// join concatenates prefix, mid, and suffix into one fresh token slice.
func join(prefix, mid, suffix []string) []string {
	out := make([]string, 0, len(prefix)+len(mid)+len(suffix))
	out = append(out, prefix...)
	out = append(out, mid...)
	out = append(out, suffix...)
	return out
}

func cloneToks(a []string) []string {
	return append([]string(nil), a...)
}

// negativesHeader is the provenance header of the negatives corpus.
func negativesHeader(source string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# TiDB parser.y negatives corpus: grammar-guided near-misses, one step from valid.\n")
	fmt.Fprintf(&b, "# Generated by tidb/parser/cmd/grammarcorpus -mode negatives from %s.\n", source)
	fmt.Fprintf(&b, "# Columns: omni verdict, engine verdict, Rule#alternative/mutation, statement.\n")
	fmt.Fprintf(&b, "# Expected reject/reject; an omni-accept/engine-reject row is an over-acceptance.\n")
	fmt.Fprintf(&b, "# Verdicts start as placeholders; fill them with\n")
	fmt.Fprintf(&b, "#   go test ./tidb/parser -run TestTiDBNegativesCorpusOmniVerdicts -update\n")
	fmt.Fprintf(&b, "#   go test ./tidb/parser -run TestTiDBNegativesCorpusEngineVerdicts -update\n")
	return b.String()
}
