package completion

import (
	"slices"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
)

// probePrefix returns a shorter text the grammar is in the same state after,
// for keyword probing to parse many times. Each probe parses the whole
// prefix, so its length is what a completion costs.
//
// Two reductions keep the LALR parser's state, and with it the set of tokens
// that can follow:
//
//   - Inside "(SELECT" or "(TABLE" the grammar reads the next token as it
//     does in a statement that starts there.
//   - A comma list after SELECT, FROM, BY or VALUES reduces to the same list
//     nonterminal whatever its middle elements are, so "SELECT a, b, c, d"
//     leaves the parser where "SELECT a, d" does. A chain of AND, or of OR,
//     is left-associative and reduces the same way.
//
// The prefix must parse up to its end; otherwise it is returned unchanged, so
// a reduction never hides an error before the caret.
func probePrefix(prefix string) string {
	reduced := collapseLists(innermostQuery(prefix))
	if reduced == prefix || !viable(prefix) || !viable(reduced) {
		return prefix
	}
	return reduced
}

// viable reports whether text is a prefix of some statement: the parser
// fails, if at all, at the end of the text.
func viable(text string) bool {
	_, err := parseOne(text)
	return err == nil || errorOffset(text, err) >= len(text)
}

// queryStarts are the reserved keywords that open a parenthesized query
// body the way they open a statement. VALUES is left out because CockroachDB
// also reads it as a name, and WITH because a statement's WITH may lead to
// INSERT, UPDATE or DELETE where a parenthesized one may not.
var queryStarts = map[string]bool{"select": true, "table": true}

// innermostQuery returns the body of the innermost parenthesized or
// bracketed query the end of prefix is inside, or prefix.
func innermostQuery(prefix string) string {
	toks := tokenize(prefix)
	if n := len(toks); n > 0 && toks[n-1].id == lexbase.ERROR {
		return prefix
	}
	var open []int
	for i, t := range toks {
		switch t.id {
		case '(', '[', '{':
			open = append(open, i)
		case ')', ']', '}':
			if len(open) > 0 {
				open = open[:len(open)-1]
			}
		}
	}
	for j := len(open) - 1; j >= 0; j-- {
		if toks[open[j]].id == '{' {
			continue
		}
		i := open[j] + 1
		if i < len(toks) && isKeywordToken(toks[i]) && queryStarts[toks[i].str] {
			return prefix[toks[i].start:]
		}
	}
	return prefix
}

// listStarts are the keywords a left-recursive comma list follows:
// target_list, from_list, sortby_list and group_by_list, and values_clause.
var listStarts = map[string]bool{"select": true, "from": true, "by": true, "values": true}

// elementKeywords are the reserved keywords that occur inside a list element
// (an expression, a table reference or a sort key). Any other reserved
// keyword at the list's depth may start another clause, so it ends the run
// of commas a collapse may remove.
var elementKeywords = map[string]bool{
	"all": true, "and": true, "any": true, "array": true, "as": true, "asc": true,
	"asymmetric": true, "case": true, "cast": true, "collate": true,
	"current_catalog": true, "current_date": true, "current_role": true,
	"current_schema": true, "current_time": true, "current_timestamp": true,
	"current_user": true, "desc": true, "distinct": true, "else": true, "end": true,
	"false": true, "in": true, "lateral": true, "localtime": true,
	"localtimestamp": true, "not": true, "null": true, "only": true, "or": true,
	"session_user": true, "some": true, "symmetric": true, "then": true,
	"true": true, "user": true, "when": true,
}

// collapseLists removes the middle elements of the runs in text that
// reduce left-recursively to one nonterminal: from just after the first
// separator of a run to just after its last one. A run is a comma list after
// a listStarts keyword, or a chain of AND or of OR at one bracket depth.
func collapseLists(text string) string {
	toks := tokenize(text)
	if n := len(toks); n > 0 && toks[n-1].id == lexbase.ERROR {
		return text
	}
	type cut struct{ from, to int }
	var cuts []cut
	// run is one open run of separators at one bracket depth.
	type run struct {
		active            bool
		firstEnd, lastEnd int
		separators        int
	}
	finish := func(r *run) {
		if r.active && r.separators >= 2 {
			cuts = append(cuts, cut{r.firstEnd, r.lastEnd})
		}
		*r = run{}
	}
	separate := func(r *run, t token) {
		r.active = true
		if r.separators == 0 {
			r.firstEnd = t.end
		}
		r.lastEnd = t.end
		r.separators++
	}
	// level is the runs open at one bracket depth. betweenAnd is set after
	// BETWEEN, whose own AND is not a conjunction.
	type level struct {
		comma, and, or run
		betweenAnd     bool
	}
	stack := []level{{}}
	finishAll := func(l *level) {
		finish(&l.comma)
		finish(&l.and)
		finish(&l.or)
	}
	for _, t := range toks {
		l := &stack[len(stack)-1]
		keyword := isKeywordToken(t)
		switch {
		case t.id == '(' || t.id == '[' || t.id == '{':
			stack = append(stack, level{})
		case t.id == ')' || t.id == ']' || t.id == '}':
			finishAll(l)
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
		case t.id == ',':
			finish(&l.and)
			finish(&l.or)
			if l.comma.active {
				separate(&l.comma, t)
			}
		case t.id == ';':
			finishAll(l)
		case keyword && t.str == "and":
			if l.betweenAnd {
				l.betweenAnd = false
				finish(&l.and)
				break
			}
			separate(&l.and, t)
		case keyword && t.str == "or":
			// AND binds tighter, so an OR ends a chain of AND.
			finish(&l.and)
			separate(&l.or, t)
		case keyword && t.str == "between":
			finish(&l.and)
			l.betweenAnd = true
		case keyword && listStarts[t.str]:
			finishAll(l)
			l.comma.active = true
		case keyword && lexbase.KeywordsCategories[t.str] == "R" && !elementKeywords[t.str]:
			finishAll(l)
		case keyword && operandEnds[t.str]:
			finish(&l.and)
			finish(&l.or)
		}
	}
	for i := range stack {
		finishAll(&stack[i])
	}
	// Runs are disjoint or nested; a run inside a removed span goes with it.
	slices.SortFunc(cuts, func(a, b cut) int { return a.from - b.from })
	out := make([]byte, 0, len(text))
	pos := 0
	for _, c := range cuts {
		if c.from < pos {
			continue
		}
		out = append(out, text[pos:c.from]...)
		pos = c.to
	}
	return string(append(out, text[pos:]...))
}

// operandEnds are the element keywords that end a chain of AND or OR: the
// branches of CASE hold expressions of their own.
var operandEnds = map[string]bool{"case": true, "when": true, "then": true, "else": true, "end": true}

// isKeywordToken reports whether t is a keyword, not an identifier spelled
// like one.
func isKeywordToken(t token) bool {
	return t.id != lexbase.IDENT && lexbase.GetKeywordID(t.str) == t.id
}
