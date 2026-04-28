package parser

import "strings"

// reverseKeywordMap maps token type → uppercase keyword name.
// Built once from keywordMap on first call to TokenName.
var reverseKeywordMap map[int]string

func initReverseKeywordMap() {
	reverseKeywordMap = make(map[int]string, len(keywordMap))
	for name, tok := range keywordMap {
		reverseKeywordMap[tok] = strings.ToUpper(name)
	}
}

// TokenName returns the SQL keyword string for a token type, or "" if not a keyword.
func TokenName(tokenType int) string {
	// Single-char tokens
	if tokenType > 0 && tokenType < 256 {
		return string(rune(tokenType))
	}
	// Build reverse map once
	if reverseKeywordMap == nil {
		initReverseKeywordMap()
	}
	if name, ok := reverseKeywordMap[tokenType]; ok {
		return name
	}
	// Non-keyword tokens (identifiers, literals, operators)
	switch tokenType {
	case tokIDENT:
		return ""
	case tokICONST, tokFCONST, tokSCONST, tokNSCONST, tokVARIABLE, tokSYSVARIABLE:
		return ""
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

// IsIdentTokenType reports whether a token type can be used as a T-SQL
// identifier (IDENT or context-sensitive keyword). Core keywords are not bare
// identifiers unless quoted, in which case the lexer emits tokIDENT.
func IsIdentTokenType(typ int) bool {
	return typ == tokIDENT || isContextKeyword(typ)
}

// Exported token type constants for completion callers.
const (
	ICONST      = tokICONST
	FCONST      = tokFCONST
	SCONST      = tokSCONST
	NSCONST     = tokNSCONST
	IDENT       = tokIDENT
	VARIABLE    = tokVARIABLE
	SYSVARIABLE = tokSYSVARIABLE

	SELECT    = kwSELECT
	INSERT    = kwINSERT
	UPDATE    = kwUPDATE
	DELETE    = kwDELETE
	MERGE     = kwMERGE
	FROM      = kwFROM
	WHERE     = kwWHERE
	SET       = kwSET
	INTO      = kwINTO
	VALUES    = kwVALUES
	AS        = kwAS
	ON        = kwON
	JOIN      = kwJOIN
	INNER     = kwINNER
	LEFT      = kwLEFT
	RIGHT     = kwRIGHT
	FULL      = kwFULL
	CROSS     = kwCROSS
	OUTER     = kwOUTER
	APPLY     = kwAPPLY
	ORDER     = kwORDER
	GROUP     = kwGROUP
	HAVING    = kwHAVING
	UNION     = kwUNION
	EXCEPT    = kwEXCEPT
	INTERSECT = kwINTERSECT
	FOR       = kwFOR
	WITH      = kwWITH
	OPTION    = kwOPTION
	TOP       = kwTOP
	DISTINCT  = kwDISTINCT
	ALL       = kwALL
	OUTPUT    = kwOUTPUT
	BY        = kwBY
	ASC       = kwASC
	DESC      = kwDESC
	OFFSET    = kwOFFSET
	FETCH     = kwFETCH
	NEXT      = kwNEXT
	FIRST     = kwFIRST
	ROW       = kwROW
	ROWS      = kwROWS
	ONLY      = kwONLY

	CREATE   = kwCREATE
	ALTER    = kwALTER
	DROP     = kwDROP
	TRUNCATE = kwTRUNCATE
	DECLARE  = kwDECLARE
	EXEC     = kwEXEC
	EXECUTE  = kwEXECUTE
	USE      = kwUSE
	GO       = kwGO

	DATABASE   = kwDATABASE
	SCHEMA     = kwSCHEMA
	TABLE      = kwTABLE
	VIEW       = kwVIEW
	SEQUENCE   = kwSEQUENCE
	COLUMN     = kwCOLUMN
	INDEX      = kwINDEX
	TRIGGER    = kwTRIGGER
	FUNCTION   = kwFUNCTION
	PROCEDURE  = kwPROCEDURE
	KEY        = kwKEY
	CONSTRAINT = kwCONSTRAINT
	PRIMARY    = kwPRIMARY
	FOREIGN    = kwFOREIGN
	REFERENCES = kwREFERENCES
	DEFAULT    = kwDEFAULT
	NULL       = kwNULL
	NOT        = kwNOT
	AND        = kwAND
	OR         = kwOR
	IS         = kwIS
	LIKE       = kwLIKE
	BETWEEN    = kwBETWEEN
	EXISTS     = kwEXISTS
	CASE       = kwCASE
	WHEN       = kwWHEN
	THEN       = kwTHEN
	ELSE       = kwELSE
	END        = kwEND
	OVER       = kwOVER
	PARTITION  = kwPARTITION
	XML        = kwXML
	JSON       = kwJSON
	BROWSE     = kwBROWSE
	PATH       = kwPATH
	AUTO       = kwAUTO
	RAW        = kwRAW
	RECOMPILE  = kwRECOMPILE
	OPTIMIZE   = kwOPTIMIZE
	MAXDOP     = kwMAXDOP
	NOLOCK     = kwNOLOCK
)

// topLevelKeywords are the statement-starting keywords offered when the cursor
// is at an empty input or after a semicolon.
var topLevelKeywords = []int{
	kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwALTER, kwDROP,
	kwMERGE, kwWITH, kwEXEC, kwEXECUTE, kwDECLARE, kwSET, kwIF, kwWHILE,
	kwBEGIN, kwPRINT, kwGRANT, kwREVOKE, kwDENY, kwTRUNCATE, kwUSE,
	kwBACKUP, kwRESTORE, kwRETURN, kwBREAK, kwCONTINUE, kwGOTO, kwWAITFOR,
	kwCOMMIT, kwROLLBACK, kwSAVE, kwOPEN, kwFETCH, kwCLOSE, kwDEALLOCATE,
	kwTHROW, kwRAISERROR, kwGO, kwDBCC, kwKILL, kwCHECKPOINT, kwRECONFIGURE,
	kwSHUTDOWN, kwBULK,
}

// RuleCandidate represents a grammar rule that is a completion candidate.
type RuleCandidate struct {
	Rule string
}

// CandidateSet holds the token and rule candidates collected during a
// completion-mode parse.
type CandidateSet struct {
	Tokens []int           // token type candidates
	Rules  []RuleCandidate // grammar rule candidates
	seen   map[int]bool    // dedup tokens
	seenR  map[string]bool // dedup rules

	// CTEPositions holds the byte offsets of WITH clause starts encountered
	// before the cursor. Used to re-parse CTE definitions and extract
	// virtual table names/columns for completion.
	CTEPositions []int

	// SelectAliasPositions holds the byte offsets of SELECT item alias
	// positions encountered before the cursor. Used to extract alias names
	// for ORDER BY / GROUP BY completion.
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
	if p.cur.Type == tokEOF && cursorOffset <= 0 {
		p.collecting = true
	}

	// Recover from panics in the parser — incomplete SQL can trigger nil
	// dereferences in grammar rules that expect more tokens.
	defer func() {
		recover() //nolint:errcheck
	}()

	// Parse statements in a loop (like the normal Parse path) so that
	// completion works after semicolons in multi-statement batches.
	for p.cur.Type != tokEOF {
		if p.cur.Type == ';' {
			p.advance()
			// After consuming semicolon, check if cursor is now reached.
			p.checkCursor()
			continue
		}
		prevPos := p.cur.Loc
		p.parseStmt() //nolint:errcheck
		// If the parser didn't advance (unrecognized token), skip it to avoid
		// an infinite loop.
		if p.cur.Loc == prevPos && p.cur.Type != tokEOF {
			p.advance()
		}
	}

	// If we reached EOF without having collected candidates (or if we were
	// already in collect mode but parseStmt was never called), add top-level
	// keywords as a fallback.
	if len(cs.Tokens) == 0 {
		for _, t := range topLevelKeywords {
			p.addTokenCandidate(t)
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

func (p *Parser) addExpressionCandidates() {
	p.addRuleCandidate("columnref")
	p.addRuleCandidate("func_name")
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
