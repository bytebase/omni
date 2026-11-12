package grammar

import (
	"fmt"
	"strings"
)

// A Row is one generated statement and the label it is known by: the rule
// and the symbols of the alternative it exercises, as in
// AlterTableStmt/ALTER_IgnoreOptional_TABLE_TableName, followed by !<mutation> for a
// near-miss. Labels name symbols rather than positions, so a release that
// adds an alternative changes no other row's label.
type Row struct {
	Rule      string
	Statement string
}

// A deriv is the shortest string a symbol can derive: its cost in tokens and
// the tokens themselves. ok is false until a terminating derivation is found.
type deriv struct {
	cost int
	toks []string
	ok   bool
}

// A ctxt is the shortest pair of token strings around a nonterminal in a
// derivation from the start symbol: start =>* prefix N suffix. cost is the
// token count of prefix and suffix together.
type ctxt struct {
	cost   int
	prefix []string
	suffix []string
	ok     bool
}

// A builtAlt is one reachable alternative with the tokens that reach it from
// the start symbol (prefix, suffix) kept apart from the alternative's own
// shortest tokens (body). A near-miss breaks the body, the production the
// row is about, and not the context that carries it there.
type builtAlt struct {
	label  string
	prefix []string
	body   []string
	suffix []string
}

// full renders the whole statement.
func (a builtAlt) full() string {
	return spaceJoin(join(a.prefix, a.body, a.suffix))
}

// Rows returns the generated rows: for each alternative of each rule, in
// source order, the shortest statement that exercises it, followed by its
// near-misses (see nearMisses). uncovered lists the alternatives no
// statement could be generated for: unreachable from the start symbol,
// never terminating, or needing an unreachable token. Each needs a
// hand-added row.
func (g *Grammar) Rows() (rows []Row, uncovered []string) {
	alts, uncovered := g.buildAlts()
	kw := g.keywordSet()

	positive := map[string]bool{}
	for _, a := range alts {
		positive[a.full()] = true
	}
	seen := map[string]bool{}
	for _, a := range alts {
		rows = append(rows, Row{Rule: a.label, Statement: a.full()})
		for _, nm := range nearMisses(a, kw) {
			// A mutation that lands on a valid sentence, or on another
			// near-miss, adds nothing.
			if nm.Statement == "" || positive[nm.Statement] || seen[nm.Statement] {
				continue
			}
			seen[nm.Statement] = true
			rows = append(rows, nm)
		}
	}
	return rows, uncovered
}

// buildAlts returns one builtAlt per alternative a statement can be
// generated for, and the labels of the rest.
func (g *Grammar) buildAlts() (alts []builtAlt, uncovered []string) {
	sd := g.shortestDerivations()
	cx := g.shortestContexts(sd)

	used := map[string]int{}
	for _, r := range g.rules {
		cxN := cx[r.name]
		for _, alt := range r.alts {
			label := altLabel(r.name, alt)
			if used[label]++; used[label] > 1 {
				// Two alternatives with the same symbols differ only in a
				// %prec or an action; number the later ones.
				label = fmt.Sprintf("%s#%d", label, used[label])
			}
			body, ok := g.altTokens(alt, sd)
			if !ok || !cxN.ok || g.needsUnreachable(alt) {
				uncovered = append(uncovered, label)
				continue
			}
			alts = append(alts, builtAlt{label: label, prefix: cxN.prefix, body: body, suffix: cxN.suffix})
		}
	}
	return alts, uncovered
}

// altLabel names an alternative by its rule and its symbols as the grammar
// writes them, a literal by its content, and %empty for no symbols.
func altLabel(rule string, alt []symbol) string {
	if len(alt) == 0 {
		return rule + "/%empty"
	}
	parts := make([]string, len(alt))
	for i, sym := range alt {
		parts[i] = sym.text
	}
	return rule + "/" + strings.Join(parts, "_")
}

// needsUnreachable reports whether an alternative names a token the lexer
// never produces in the judge's mode. Generating it anyway would exercise
// whatever production the substitute text lexes as instead.
func (g *Grammar) needsUnreachable(alt []symbol) bool {
	for _, sym := range alt {
		if !sym.nonterm && !sym.literal && g.unreachable[sym.text] {
			return true
		}
	}
	return false
}

// symDeriv returns the shortest derivation of one symbol: a terminal
// derives its own sample text in one token, a nonterminal what sd records.
func (g *Grammar) symDeriv(sym symbol, sd map[string]deriv) deriv {
	if sym.nonterm {
		return sd[sym.text]
	}
	if !sym.literal && g.unreachable[sym.text] {
		return deriv{ok: false}
	}
	text := g.sample(sym)
	if text == "" {
		// A token written as nothing still counts as a token consumed.
		return deriv{cost: 1, toks: nil, ok: true}
	}
	return deriv{cost: 1, toks: []string{text}, ok: true}
}

// altTokens returns the tokens an alternative derives by taking the
// shortest derivation of each of its symbols, and whether every symbol
// terminated.
func (g *Grammar) altTokens(alt []symbol, sd map[string]deriv) ([]string, bool) {
	var toks []string
	for _, sym := range alt {
		d := g.symDeriv(sym, sd)
		if !d.ok {
			return nil, false
		}
		toks = append(toks, d.toks...)
	}
	return toks, true
}

// shortestDerivations computes, for every nonterminal, the shortest string
// it derives, by relaxing each rule to the cheapest of its alternatives
// until no cost improves.
func (g *Grammar) shortestDerivations() map[string]deriv {
	sd := make(map[string]deriv, len(g.rules))
	for _, r := range g.rules {
		sd[r.name] = deriv{ok: false}
	}
	for changed := true; changed; {
		changed = false
		for _, r := range g.rules {
			best := deriv{ok: false}
			for _, alt := range r.alts {
				cost, toks, ok := 0, []string(nil), true
				for _, sym := range alt {
					d := g.symDeriv(sym, sd)
					if !d.ok {
						ok = false
						break
					}
					cost += d.cost
					toks = append(toks, d.toks...)
				}
				if ok && (!best.ok || cost < best.cost) {
					best = deriv{cost: cost, toks: toks, ok: true}
				}
			}
			if best.ok {
				cur := sd[r.name]
				if !cur.ok || best.cost < cur.cost {
					sd[r.name] = best
					changed = true
				}
			}
		}
	}
	return sd
}

// shortestContexts computes, for every nonterminal, the shortest
// surrounding text in a derivation from the start symbol. It starts the
// start symbol at an empty context and relaxes each occurrence of a
// nonterminal in a rule body, filling the siblings around it with their
// shortest derivations, until no context cost improves.
func (g *Grammar) shortestContexts(sd map[string]deriv) map[string]ctxt {
	cx := make(map[string]ctxt, len(g.rules))
	for _, r := range g.rules {
		cx[r.name] = ctxt{ok: false}
	}
	cx[g.start] = ctxt{cost: 0, ok: true}

	for changed := true; changed; {
		changed = false
		for _, r := range g.rules {
			cm := cx[r.name]
			if !cm.ok {
				continue
			}
			for _, alt := range r.alts {
				derivs := make([]deriv, len(alt))
				total, ok := 0, true
				for k, sym := range alt {
					derivs[k] = g.symDeriv(sym, sd)
					if !derivs[k].ok {
						ok = false
						break
					}
					total += derivs[k].cost
				}
				if !ok {
					continue
				}
				var pre []string
				for k, sym := range alt {
					if sym.nonterm {
						cost := cm.cost + total - derivs[k].cost
						prefix := concat(cm.prefix, pre)
						suffix := concat(sibSuffix(derivs, k), cm.suffix)
						if relaxCtx(cx, sym.text, ctxt{cost: cost, prefix: prefix, suffix: suffix, ok: true}) {
							changed = true
						}
					}
					pre = concat(pre, derivs[k].toks)
				}
			}
		}
	}
	return cx
}

// sibSuffix returns the tokens the symbols after position k derive.
func sibSuffix(derivs []deriv, k int) []string {
	var suf []string
	for j := k + 1; j < len(derivs); j++ {
		suf = append(suf, derivs[j].toks...)
	}
	return suf
}

// relaxCtx lowers a nonterminal's recorded context when the candidate is
// cheaper, and reports whether it changed anything.
func relaxCtx(cx map[string]ctxt, name string, cand ctxt) bool {
	cur := cx[name]
	if cur.ok && cur.cost <= cand.cost {
		return false
	}
	cx[name] = cand
	return true
}

// concat returns a fresh slice with b appended to a, sharing neither's array.
func concat(a, b []string) []string {
	out := make([]string, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
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

// spaceJoin renders a token list as SQL, inserting a space between tokens
// except where one would be wrong: before ',', ';', ')', '.', '(' and after
// '(' or '.'. Keeping function names and parentheses tight, and dotted
// names unbroken, keeps a generated sentence close to how it is written.
func spaceJoin(toks []string) string {
	var b strings.Builder
	prev := ""
	for _, t := range toks {
		if t == "" {
			continue
		}
		if prev != "" && needSpace(prev, t) {
			b.WriteByte(' ')
		}
		b.WriteString(t)
		prev = t
	}
	return b.String()
}

func needSpace(prev, cur string) bool {
	// The user@host separator binds tightly on both sides: DROP USER 'u'@'h'.
	if cur == "@" || prev == "@" {
		return false
	}
	switch cur[0] {
	case ',', ';', ')', '.', '(':
		return false
	}
	switch prev[len(prev)-1] {
	case '(', '.':
		return false
	}
	return true
}
