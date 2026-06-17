package parser

import "strings"

// TokenName returns the SQL keyword string for a token type, or "" if not a keyword.
func TokenName(tokenType int) string {
	// Single-char tokens (e.g., '(' ')' ',' ';')
	if tokenType > 0 && tokenType < 256 {
		return string(rune(tokenType))
	}
	// Search keyword table (reverse lookup)
	for name, tok := range keywords {
		if tok == tokenType {
			return strings.ToUpper(name)
		}
	}
	return ""
}

// Tokenize runs the lexer on sql and returns all non-EOF tokens.
// Useful for walking the token stream without a full parse.
func Tokenize(sql string) []Token {
	lex := NewLexer(sql)
	var tokens []Token
	for {
		tok := lex.NextToken()
		if tok.Type == tokEOF {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

// CandidateSet holds the token and rule candidates collected during a
// completion-mode parse.
type CandidateSet struct {
	Tokens []int           // token type candidates
	Rules  []RuleCandidate // grammar rule candidates
	seen   map[int]bool    // dedup tokens
	seenR  map[string]bool // dedup rules

	// CTEPositions holds the byte offsets of WITH clause starts encountered
	// before the cursor. Bytebase uses these to re-parse CTE definitions
	// and extract virtual table names/columns for completion.
	CTEPositions []int

	// SelectAliasPositions holds the byte offsets of SELECT item alias
	// positions encountered before the cursor. Bytebase uses these to
	// extract alias names for ORDER BY / GROUP BY completion.
	SelectAliasPositions []int
}

// RuleCandidate represents a grammar rule that is a completion candidate.
type RuleCandidate struct {
	Rule string
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

// Collect runs the parser in completion mode and returns the set of token
// and rule candidates at the given cursor offset.
func Collect(sql string, cursorOffset int) *CandidateSet {
	// Use Split to find the segment containing the cursor.
	segments := Split(sql)
	for _, seg := range segments {
		if seg.Empty() {
			continue
		}
		if cursorOffset >= seg.ByteStart && cursorOffset <= seg.ByteEnd {
			return collectSingle(seg.Text, cursorOffset-seg.ByteStart)
		}
	}
	// Cursor is past all segments — collect on trailing text after last segment.
	if len(segments) > 0 {
		last := segments[len(segments)-1]
		trailing := sql[last.ByteEnd:]
		if len(trailing) > 0 {
			localCursor := cursorOffset - last.ByteEnd
			if localCursor >= 0 {
				return collectSingle(trailing, localCursor)
			}
		}
	}
	// Fallback: collect on the whole input.
	return collectSingle(sql, cursorOffset)
}

// collectSingle runs completion on a single segment (the original Collect logic).
func collectSingle(sql string, cursorOffset int) *CandidateSet {
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
	if p.cur.Type == tokEOF && cursorOffset <= 0 {
		p.collecting = true
	}

	// Recover from panics in the parser — incomplete SQL can trigger nil
	// dereferences in grammar rules that expect more tokens.
	defer func() {
		recover() //nolint:errcheck
	}()

	// Parse statements in a loop so that multi-statement
	// input with semicolons works within a segment.
	for {
		if p.cur.Type == ';' {
			p.advance()
			if p.completing && !p.collecting && p.cur.Loc >= p.cursorOff {
				p.collecting = true
			}
			continue
		}
		if p.cur.Type == tokEOF && !p.collectMode() {
			break
		}
		prevLoc := p.cur.Loc
		p.parseStmt() //nolint:errcheck
		if p.cur.Type == tokEOF || p.collectMode() {
			break
		}
		if p.cur.Loc == prevLoc {
			p.advance()
		}
	}
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
	if p.completing && !p.collecting && p.cur.Loc >= p.cursorOff {
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

// IsIdentTokenType reports whether a token type can be used as an identifier
// (IDENT or non-reserved keyword).
func IsIdentTokenType(typ int) bool {
	return typ == tokIDENT || typ >= 700
}

// Exported token type constants for use by the completion package.
const (
	SELECT  = kwSELECT
	INSERT  = kwINSERT
	UPDATE  = kwUPDATE
	DELETE  = kwDELETE
	FROM    = kwFROM
	WHERE   = kwWHERE
	SET     = kwSET
	INTO    = kwINTO
	VALUES  = kwVALUES
	AS      = kwAS
	ON      = kwON
	USING   = kwUSING
	JOIN    = kwJOIN
	INNER   = kwINNER
	LEFT    = kwLEFT
	RIGHT   = kwRIGHT
	CROSS   = kwCROSS
	NATURAL = kwNATURAL
	ORDER   = kwORDER
	GROUP   = kwGROUP
	HAVING  = kwHAVING
	LIMIT   = kwLIMIT
	UNION   = kwUNION
	FOR     = kwFOR
	REPLACE = kwREPLACE
)
