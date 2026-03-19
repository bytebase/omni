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

// Tokenize runs the lexer on sql and returns all non-EOF tokens with
// parser-mapped token types. Useful for walking the token stream without
// a full parse.
func Tokenize(sql string) []Token {
	lex := NewLexer(sql)
	var tokens []Token
	for {
		tok := lex.NextToken()
		if tok.Type == 0 { // EOF
			break
		}
		tok.Type = mapTokenType(tok.Type)
		tokens = append(tokens, tok)
	}
	return tokens
}

// IsReservedKeyword returns true if name (case-insensitive) is a PostgreSQL
// reserved keyword.
func IsReservedKeyword(name string) bool {
	kw := LookupKeyword(name)
	return kw != nil && kw.Category == ReservedKeyword
}

// IsIdentifierTokenType reports whether tokenType can act as an identifier
// (IDENT, or any keyword that is not strictly reserved).
func IsIdentifierTokenType(tokenType int) bool {
	if tokenType == IDENT {
		return true
	}
	for i := range Keywords {
		if Keywords[i].Token == tokenType {
			return Keywords[i].Category != ReservedKeyword
		}
	}
	return false
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

	// CTEPositions holds the byte offsets of WITH clause starts encountered
	// before the cursor. Bytebase uses these to re-parse CTE definitions
	// and extract virtual table names/columns for completion.
	CTEPositions []int

	// SelectAliasPositions holds the byte offsets of SELECT item alias
	// positions encountered before the cursor. Bytebase uses these to
	// extract alias names for ORDER BY / GROUP BY completion.
	SelectAliasPositions []int
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
		source:     sql,
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

	// Recover from panics in the parser — incomplete SQL can trigger nil
	// dereferences in grammar rules that expect more tokens.
	defer func() {
		recover() //nolint:errcheck
	}()

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

// addCTEPosition records a WITH clause byte offset in the candidate set.
func (p *Parser) addCTEPosition(pos int) {
	if p.candidates != nil {
		p.candidates.CTEPositions = append(p.candidates.CTEPositions, pos)
	}
}

// addSelectAliasPosition records a SELECT alias byte offset in the candidate set.
func (p *Parser) addSelectAliasPosition(pos int) {
	if p.candidates != nil {
		p.candidates.SelectAliasPositions = append(p.candidates.SelectAliasPositions, pos)
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

// snapshot returns a copy of the current candidate set state.
func (cs *CandidateSet) snapshot() *CandidateSet {
	s := &CandidateSet{
		Tokens: make([]int, len(cs.Tokens)),
		Rules:  make([]RuleCandidate, len(cs.Rules)),
		seen:   make(map[int]bool, len(cs.seen)),
		seenR:  make(map[string]bool, len(cs.seenR)),
	}
	copy(s.Tokens, cs.Tokens)
	copy(s.Rules, cs.Rules)
	for k, v := range cs.seen {
		s.seen[k] = v
	}
	for k, v := range cs.seenR {
		s.seenR[k] = v
	}
	return s
}

// diff returns candidates in cs that are not in before.
func (cs *CandidateSet) diff(before *CandidateSet) *CandidateSet {
	d := newCandidateSet()
	for _, tok := range cs.Tokens {
		if !before.seen[tok] {
			d.addToken(tok)
		}
	}
	for _, rule := range cs.Rules {
		if !before.seenR[rule.Rule] {
			d.addRule(rule.Rule)
		}
	}
	if len(d.Tokens) == 0 && len(d.Rules) == 0 {
		return nil
	}
	return d
}
