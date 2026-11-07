package main

import (
	"fmt"
	"strings"
)

// A Row is one generated corpus case: the alternative it exercises, labelled
// Rule#index (1-based), and the shortest statement that reaches it from the
// start symbol.
type Row struct {
	label string
	stmt  string
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

// generate returns one Row per reachable alternative, and the labels of the
// alternatives it could not reach or terminate, which need a hand-written
// sentence instead.
func (g *Grammar) generate() (rows []Row, uncovered []string) {
	sd := g.shortestDerivations()
	cx := g.shortestContexts(sd)

	for _, rule := range g.rules {
		cxN := cx[rule.name]
		for i, alt := range rule.alts {
			label := fmt.Sprintf("%s#%d", rule.name, i+1)
			body, ok := g.altTokens(alt, sd)
			if !ok || !cxN.ok {
				uncovered = append(uncovered, label)
				continue
			}
			toks := make([]string, 0, len(cxN.prefix)+len(body)+len(cxN.suffix))
			toks = append(toks, cxN.prefix...)
			toks = append(toks, body...)
			toks = append(toks, cxN.suffix...)
			rows = append(rows, Row{label: label, stmt: spaceJoin(toks)})
		}
	}
	return rows, uncovered
}

// symDeriv returns the shortest derivation of one symbol: a terminal derives
// its own sample text in one token, a nonterminal derives what sd records.
func (g *Grammar) symDeriv(sym Symbol, sd map[string]deriv) deriv {
	if sym.nonterm {
		return sd[sym.text]
	}
	text := g.sample(sym)
	if text == "" {
		// A terminal with no spelling (the lexer-only error token, hints):
		// it still counts as a token consumed, but emits nothing.
		return deriv{cost: 1, toks: nil, ok: true}
	}
	return deriv{cost: 1, toks: []string{text}, ok: true}
}

// altTokens returns the tokens an alternative derives by taking the shortest
// derivation of each of its symbols, and whether every symbol terminated.
func (g *Grammar) altTokens(alt []Symbol, sd map[string]deriv) ([]string, bool) {
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

// shortestDerivations computes, for every nonterminal, the shortest string it
// derives, by relaxing each rule to the cheapest of its alternatives until no
// cost improves.
func (g *Grammar) shortestDerivations() map[string]deriv {
	sd := make(map[string]deriv, len(g.rules))
	for _, rule := range g.rules {
		sd[rule.name] = deriv{ok: false}
	}
	for changed := true; changed; {
		changed = false
		for _, rule := range g.rules {
			best := deriv{ok: false}
			for _, alt := range rule.alts {
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
				cur := sd[rule.name]
				if !cur.ok || best.cost < cur.cost {
					sd[rule.name] = best
					changed = true
				}
			}
		}
	}
	return sd
}

// shortestContexts computes, for every nonterminal, the shortest surrounding
// text in a derivation from the start symbol. It starts the start symbol at an
// empty context and relaxes each occurrence of a nonterminal in a rule body,
// filling the siblings around it with their shortest derivations, until no
// context cost improves.
func (g *Grammar) shortestContexts(sd map[string]deriv) map[string]ctxt {
	cx := make(map[string]ctxt, len(g.rules))
	for _, rule := range g.rules {
		cx[rule.name] = ctxt{ok: false}
	}
	start := g.startSymbol()
	cx[start] = ctxt{cost: 0, ok: true}

	for changed := true; changed; {
		changed = false
		for _, rule := range g.rules {
			cm := cx[rule.name]
			if !cm.ok {
				continue
			}
			for _, alt := range rule.alts {
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

// startSymbol returns the grammar's start symbol: the rule named Start if one
// exists, else the first rule goyacc would take as the start.
func (g *Grammar) startSymbol() string {
	if _, ok := g.byName["Start"]; ok {
		return "Start"
	}
	return g.rules[0].name
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

// spaceJoin renders a token list as SQL, inserting a space between tokens
// except where one would be wrong: before ',', ';', ')', '.', '(' and after
// '(' or '.'. Keeping function names and parentheses tight, and dotted names
// unbroken, keeps a generated sentence closer to how the statement is written.
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
