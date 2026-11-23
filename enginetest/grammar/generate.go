package grammar

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"
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
	// site is set for the row an EverySite rule adds to an alternative that
	// names it: the body tokens the rule's clause spans, [start, end).
	// Its near-misses drop only those, since the alternative's own row
	// breaks the rest.
	site *[2]int
}

// full returns the whole statement's tokens.
func (a builtAlt) full() []string {
	return join(a.prefix, a.body, a.suffix)
}

// Rows returns the generated rows: for each alternative of each rule, in
// source order, the shortest statement that exercises it, followed by its
// near-misses (see nearMisses), and after it a row for each EverySite rule
// the alternative names, with that rule's clause written out. uncovered lists the alternatives no
// statement could be generated for: unreachable from the start symbol,
// never terminating, or needing an unreachable token. Each needs a
// hand-added row.
func (g *Grammar) Rows() (rows []Row, uncovered []string) {
	alts, uncovered := g.buildAlts()
	kw := g.keywordSet()

	positive := map[string]bool{}
	for _, a := range alts {
		if a.site == nil {
			positive[g.render(a.full())] = true
		}
	}
	// A clause's row at a site whose statement another row already holds,
	// as the row of the rule's own alternative does, adds nothing.
	kept := alts[:0]
	for _, a := range alts {
		if a.site != nil {
			stmt := g.render(a.full())
			if positive[stmt] {
				continue
			}
			positive[stmt] = true
		}
		kept = append(kept, a)
	}
	seen := map[string]bool{}
	for _, a := range kept {
		rows = append(rows, Row{Rule: a.label, Statement: g.render(a.full())})
		for _, nm := range nearMisses(a, kw, g.render) {
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
	var ne map[string][]string
	if len(g.everySite) > 0 {
		ne = g.shortestNonEmpty(sd)
	}

	used := map[string]int{}
	for _, r := range g.rules {
		cxN := cx[r.name]
		for _, alt := range r.alts {
			label := g.altLabel(r.name, alt)
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
			alts = append(alts, g.siteAlts(label, alt, sd, ne, cxN)...)
		}
	}
	return alts, uncovered
}

// siteAlts returns one builtAlt for each symbol of alt that names an
// EverySite rule, labelled <label>+<rule> (and #2, #3 for a rule the
// alternative names again): the alternative's shortest statement with that
// rule written as its shortest form that is not nothing.
func (g *Grammar) siteAlts(label string, alt []symbol, sd map[string]deriv, ne map[string][]string, cx ctxt) []builtAlt {
	var out []builtAlt
	named := map[string]int{}
	for k, sym := range alt {
		if !sym.nonterm || !g.everySite[sym.text] {
			continue
		}
		clause := ne[sym.text]
		if len(clause) == 0 {
			continue
		}
		siteLabel := label + "+" + sym.text
		if named[sym.text]++; named[sym.text] > 1 {
			siteLabel = fmt.Sprintf("%s#%d", siteLabel, named[sym.text])
		}
		before, ok := g.altTokens(alt[:k], sd)
		if !ok {
			continue
		}
		after, ok := g.altTokens(alt[k+1:], sd)
		if !ok {
			continue
		}
		out = append(out, builtAlt{
			label:  siteLabel,
			prefix: cx.prefix,
			body:   join(before, clause, after),
			suffix: cx.suffix,
			site:   &[2]int{len(before), len(before) + len(clause)},
		})
	}
	return out
}

// shortestNonEmpty returns, for each rule, the fewest tokens other than
// none it derives, through a nullable rule it names as well: the shortest
// derivation of every symbol but one, which derives its own shortest form
// that is not nothing. Among equally short forms it keeps the first found,
// walking the rules and their alternatives in source order.
func (g *Grammar) shortestNonEmpty(sd map[string]deriv) map[string][]string {
	ne := map[string][]string{}
	for name, text := range g.ruleText {
		if text != "" {
			ne[name] = []string{text}
		}
	}
	for changed := true; changed; {
		changed = false
		for _, r := range g.rules {
			if _, sampled := g.ruleText[r.name]; sampled {
				continue
			}
			for _, alt := range r.alts {
				ds := make([]deriv, len(alt))
				ok := true
				for k, sym := range alt {
					if ds[k] = g.symDeriv(sym, sd); !ds[k].ok {
						ok = false
						break
					}
				}
				if !ok {
					continue
				}
				var pre []string
				for k, sym := range alt {
					var mid []string
					if sym.nonterm {
						mid = ne[sym.text]
					} else {
						mid = ds[k].toks
					}
					if len(mid) > 0 {
						toks := join(pre, mid, sibSuffix(ds, k))
						if cur, known := ne[r.name]; !known || len(toks) < len(cur) {
							ne[r.name] = toks
							changed = true
						}
					}
					pre = concat(pre, ds[k].toks)
				}
			}
		}
	}
	return ne
}

// altLabel names an alternative by its rule and its symbols as the grammar
// writes them, a literal by its content, and %empty for no symbols. Where
// the grammar caps labels, a longer list of symbols is cut and ends in a
// hash of the whole.
func (g *Grammar) altLabel(rule string, alt []symbol) string {
	if len(alt) == 0 {
		return rule + "/%empty"
	}
	parts := make([]string, len(alt))
	for i, sym := range alt {
		parts[i] = sym.text
		if sym.label != "" {
			parts[i] = sym.label
		}
	}
	return rule + "/" + capLabel(strings.Join(parts, "_"), g.labelCap())
}

// labelCap is the longest an alternative's symbols may be in its label, or
// 0 for no limit.
func (g *Grammar) labelCap() int {
	if g.antlr {
		return maxAltLabel
	}
	return 0
}

// capLabel cuts text longer than max (0 for no limit) and ends it in ~ and
// eight hex digits of the whole text's SHA-256, so the label stays as
// stable as the text.
func capLabel(text string, max int) string {
	if max == 0 || len(text) <= max {
		return text
	}
	sum := sha256.Sum256([]byte(text))
	cut := text[:max-9]
	for !utf8.ValidString(cut) {
		cut = cut[:len(cut)-1]
	}
	return cut + "~" + hex.EncodeToString(sum[:4])
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
// it derives. Among equally short alternatives it takes the one the grammar
// lists first, which is its plainest form: a name before a keyword that may
// stand for one, the first data type before the rest. A rule the token
// table samples derives that text, as one token.
func (g *Grammar) shortestDerivations() map[string]deriv {
	// Costs first: relax each rule to the cheapest of its alternatives until
	// none improves. The cost of a shortest derivation is unique; which
	// alternative reaches it is not, and is chosen below.
	cost := make(map[string]int, len(g.rules))
	known := make(map[string]bool, len(g.rules))
	symCost := func(sym symbol) (int, bool) {
		if sym.nonterm {
			return cost[sym.text], known[sym.text]
		}
		// A token written as nothing still counts as a token consumed.
		return 1, sym.literal || !g.unreachable[sym.text]
	}
	altCost := func(alt []symbol) (int, bool) {
		total := 0
		for _, sym := range alt {
			c, ok := symCost(sym)
			if !ok {
				return 0, false
			}
			total += c
		}
		return total, true
	}
	for name := range g.ruleText {
		cost[name], known[name] = 1, true
	}
	for changed := true; changed; {
		changed = false
		for _, r := range g.rules {
			if _, sampled := g.ruleText[r.name]; sampled {
				continue
			}
			for _, alt := range r.alts {
				if c, ok := altCost(alt); ok && (!known[r.name] || c < cost[r.name]) {
					cost[r.name], known[r.name] = c, true
					changed = true
				}
			}
		}
	}

	// Then the text: the first alternative that costs the minimum, expanded
	// the same way. A nonterminal being expanded further up cannot be used
	// again, since that derivation would never end; the alternative that
	// needs it is passed over for the next one.
	sd := make(map[string]deriv, len(g.rules))
	for _, r := range g.rules {
		sd[r.name] = deriv{ok: false}
	}
	expanding := map[string]bool{}
	var expand func(name string) ([]string, bool)
	expand = func(name string) ([]string, bool) {
		if d := sd[name]; d.ok {
			return d.toks, true
		}
		if text, sampled := g.ruleText[name]; sampled {
			sd[name] = deriv{cost: 1, toks: []string{text}, ok: true}
			return sd[name].toks, true
		}
		if !known[name] || expanding[name] {
			return nil, false
		}
		expanding[name] = true
		defer delete(expanding, name)
	alternatives:
		for _, alt := range g.byName[name].alts {
			if c, ok := altCost(alt); !ok || c != cost[name] {
				continue
			}
			var toks []string
			for _, sym := range alt {
				if sym.nonterm {
					sub, ok := expand(sym.text)
					if !ok {
						continue alternatives
					}
					toks = append(toks, sub...)
				} else if text := g.sample(sym); text != "" {
					toks = append(toks, text)
				}
			}
			sd[name] = deriv{cost: cost[name], toks: toks, ok: true}
			return toks, true
		}
		return nil, false
	}
	for _, r := range g.rules {
		expand(r.name)
	}
	return sd
}

// shortestContexts computes, for every nonterminal, the shortest
// surrounding text in a derivation from the start symbol: the siblings
// around it filled with their shortest derivations, inside its parent's own
// shortest context. Among equally short contexts it takes the one found
// first walking the grammar from the start symbol, alternatives in order
// and depth first, so a nonterminal that several statements reach equally
// cheaply sits in the statement the grammar lists first.
func (g *Grammar) shortestContexts(sd map[string]deriv) map[string]ctxt {
	// derivs returns the shortest derivation of each symbol of an
	// alternative and their total cost, or false when one never terminates.
	derivs := func(alt []symbol) ([]deriv, int, bool) {
		ds := make([]deriv, len(alt))
		total := 0
		for k, sym := range alt {
			ds[k] = g.symDeriv(sym, sd)
			if !ds[k].ok {
				return nil, 0, false
			}
			total += ds[k].cost
		}
		return ds, total, true
	}

	// Costs first, as in shortestDerivations.
	cost := map[string]int{g.start: 0}
	for changed := true; changed; {
		changed = false
		for _, r := range g.rules {
			base, ok := cost[r.name]
			if !ok {
				continue
			}
			for _, alt := range r.alts {
				ds, total, ok := derivs(alt)
				if !ok {
					continue
				}
				for k, sym := range alt {
					if !sym.nonterm {
						continue
					}
					c := base + total - ds[k].cost
					if cur, seen := cost[sym.text]; !seen || c < cur {
						cost[sym.text] = c
						changed = true
					}
				}
			}
		}
	}

	cx := make(map[string]ctxt, len(g.rules))
	for _, r := range g.rules {
		cx[r.name] = ctxt{ok: false}
	}
	cx[g.start] = ctxt{cost: 0, ok: true}
	var walk func(name string)
	walk = func(name string) {
		cm := cx[name]
		for _, alt := range g.byName[name].alts {
			ds, total, ok := derivs(alt)
			if !ok {
				continue
			}
			var pre []string
			for k, sym := range alt {
				if sym.nonterm && !cx[sym.text].ok && cm.cost+total-ds[k].cost == cost[sym.text] {
					cx[sym.text] = ctxt{
						cost:   cost[sym.text],
						prefix: concat(cm.prefix, pre),
						suffix: concat(sibSuffix(ds, k), cm.suffix),
						ok:     true,
					}
					walk(sym.text)
				}
				pre = concat(pre, ds[k].toks)
			}
		}
	}
	walk(g.start)
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

// render writes a token list as SQL, inserting a space between tokens
// except where one would be wrong: before ',', ';', ')', '.', '(' and after
// '(' or '.'. Keeping function names and parentheses tight, and dotted
// names unbroken, keeps a generated sentence close to how it is written. A
// token the engine's table marks as spaced always stands apart.
func (g *Grammar) render(toks []string) string {
	var b strings.Builder
	prev := ""
	for _, t := range toks {
		if t == "" {
			continue
		}
		if prev != "" && (g.spaced[prev] || g.spaced[t] || needSpace(prev, t)) {
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
