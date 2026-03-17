package parser

import "strings"

// TokenName returns the SQL keyword string for a token type, or "" if not a keyword.
func TokenName(tokenType int) string {
	// Single-char tokens
	if tokenType > 0 && tokenType < 256 {
		return string(rune(tokenType))
	}
	// Search keyword table
	for i := range Keywords {
		if Keywords[i].Token == tokenType {
			return strings.ToUpper(Keywords[i].Name)
		}
	}
	// Special non-keyword tokens
	switch tokenType {
	case IDENT:
		return "" // identifier — not a keyword
	case ICONST, FCONST, SCONST, BCONST, XCONST, PARAM, Op:
		return "" // literals/operators — not completable keywords
	}
	return ""
}

// MapTokenType maps lexer-internal token types to parser token constants.
func MapTokenType(typ int) int {
	return mapTokenType(typ)
}

// RuleCandidate represents a grammar rule that is a completion candidate.
type RuleCandidate struct {
	Rule string
}

// CandidateSet holds the token and rule candidates collected during a
// completion-mode parse.
type CandidateSet struct {
	Tokens []int            // token type candidates
	Rules  []RuleCandidate  // grammar rule candidates
	seen   map[int]bool     // dedup tokens
	seenR  map[string]bool  // dedup rules
}

// newCandidateSet creates an empty CandidateSet.
func newCandidateSet() *CandidateSet {
	return &CandidateSet{
		seen:  make(map[int]bool),
		seenR: make(map[string]bool),
	}
}

// addToken adds a token type to the candidate set (deduped).
func (cs *CandidateSet) addToken(t int) {
	if cs.seen[t] {
		return
	}
	cs.seen[t] = true
	cs.Tokens = append(cs.Tokens, t)
}

// addRule adds a rule name to the candidate set (deduped).
func (cs *CandidateSet) addRule(r string) {
	if cs.seenR[r] {
		return
	}
	cs.seenR[r] = true
	cs.Rules = append(cs.Rules, RuleCandidate{Rule: r})
}

// HasToken reports whether the candidate set contains the given token type.
func (cs *CandidateSet) HasToken(t int) bool {
	return cs.seen[t]
}

// HasRule reports whether the candidate set contains the given rule name.
func (cs *CandidateSet) HasRule(r string) bool {
	return cs.seenR[r]
}

// errCollecting is a sentinel error used to unwind the parser stack when
// collecting completion candidates.
var errCollecting = &ParseError{Message: "collecting"}

// Collect runs the parser in completion mode and returns the set of token
// and rule candidates at the given cursor offset.
func Collect(sql string, cursorOffset int) *CandidateSet {
	cs := newCandidateSet()
	p := &Parser{
		lexer:      NewLexer(sql),
		completing: true,
		cursorOff:  cursorOffset,
		candidates: cs,
		maxCollect: 100,
	}
	p.advance()

	// If we're already at EOF and cursor is at offset 0, trigger collection.
	if p.cur.Type == 0 && cursorOffset <= 0 {
		p.collecting = true
	}

	p.parseStmt()
	return cs
}

// collectMode reports whether the parser is in active collection mode
// (completing is true and the cursor position has been reached).
func (p *Parser) collectMode() bool {
	return p.completing && p.collecting
}

// checkCursor checks whether the current token is at or past the cursor
// offset, and if so, enables collection mode.
func (p *Parser) checkCursor() {
	if p.cur.Loc >= p.cursorOff {
		p.collecting = true
	}
}

// addTokenCandidate adds a token type to the candidate set.
func (p *Parser) addTokenCandidate(t int) {
	if p.candidates != nil {
		p.candidates.addToken(t)
	}
}

// addRuleCandidate adds a rule name to the candidate set.
func (p *Parser) addRuleCandidate(r string) {
	if p.candidates != nil {
		p.candidates.addRule(r)
	}
}

// addInfixCandidates adds all a_expr infix/postfix operator tokens
// with precedence >= minPrec as completion candidates.
func (p *Parser) addInfixCandidates(minPrec int) {
	type precTok struct {
		prec int
		tok  int
	}
	ops := []precTok{
		{precOr, OR}, {precAnd, AND},
		{precIs, IS}, {precIs, ISNULL}, {precIs, NOTNULL},
		{precComparison, '<'}, {precComparison, '>'}, {precComparison, '='},
		{precComparison, LESS_EQUALS}, {precComparison, GREATER_EQUALS}, {precComparison, NOT_EQUALS},
		{precIn, BETWEEN}, {precIn, IN_P}, {precIn, LIKE}, {precIn, ILIKE}, {precIn, SIMILAR},
		{precIn, NOT_LA},
		{precOp, Op},
		{precAdd, '+'}, {precAdd, '-'},
		{precMul, '*'}, {precMul, '/'}, {precMul, '%'},
		{precExp, '^'},
		{precAt, AT},
		{precCollate, COLLATE},
		{precTypecast, TYPECAST}, {precTypecast, '['},
	}
	for _, o := range ops {
		if o.prec >= minPrec {
			p.addTokenCandidate(o.tok)
		}
	}
}

// addKeywordsByCategory adds all keywords matching the given categories as candidates.
func (p *Parser) addKeywordsByCategory(categories ...KeywordCategory) {
	for i := range Keywords {
		for _, cat := range categories {
			if Keywords[i].Category == cat {
				p.addTokenCandidate(Keywords[i].Token)
				break
			}
		}
	}
}
