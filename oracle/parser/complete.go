package parser

import (
	"sort"
	"strings"
)

// Exported token aliases used by completion callers and tests.
const (
	IDENT = tokIDENT

	SELECT    = kwSELECT
	INSERT    = kwINSERT
	UPDATE    = kwUPDATE
	DELETE    = kwDELETE
	MERGE     = kwMERGE
	CREATE    = kwCREATE
	ALTER     = kwALTER
	DROP      = kwDROP
	TRUNCATE  = kwTRUNCATE
	COMMENT   = kwCOMMENT
	GRANT     = kwGRANT
	REVOKE    = kwREVOKE
	COMMIT    = kwCOMMIT
	ROLLBACK  = kwROLLBACK
	SAVEPOINT = kwSAVEPOINT
	SET       = kwSET
	BEGIN     = kwBEGIN
	DECLARE   = kwDECLARE
	WITH      = kwWITH

	FROM      = kwFROM
	WHERE     = kwWHERE
	GROUP     = kwGROUP
	HAVING    = kwHAVING
	ORDER     = kwORDER
	BY        = kwBY
	AS        = kwAS
	JOIN      = kwJOIN
	ON        = kwON
	USING     = kwUSING
	UNION     = kwUNION
	INTERSECT = kwINTERSECT
	MINUS     = kwMINUS
	OFFSET    = kwOFFSET
	FETCH     = kwFETCH

	TABLE     = kwTABLE
	VIEW      = kwVIEW
	SEQUENCE  = kwSEQUENCE
	SYNONYM   = kwSYNONYM
	INDEX     = kwINDEX
	USER      = kwUSER
	ROLE      = kwROLE
	TYPE      = kwTYPE
	PACKAGE   = kwPACKAGE
	PROCEDURE = kwPROCEDURE
	FUNCTION  = kwFUNCTION
	TRIGGER   = kwTRIGGER
)

var reverseOracleKeywordMap map[int]string

func initReverseOracleKeywordMap() {
	reverseOracleKeywordMap = make(map[int]string, len(oracleKeywords))
	for word, tok := range oracleKeywords {
		reverseOracleKeywordMap[tok] = strings.ToUpper(word)
	}
}

// TokenName returns the SQL keyword string for a token type, or "" when the
// token is not a keyword useful for completion.
func TokenName(tokenType int) string {
	if tokenType > 0 && tokenType < 256 {
		return string(rune(tokenType))
	}
	if reverseOracleKeywordMap == nil {
		initReverseOracleKeywordMap()
	}
	if name, ok := reverseOracleKeywordMap[tokenType]; ok {
		return name
	}
	switch tokenType {
	case tokIDENT, tokQIDENT, tokICONST, tokFCONST, tokSCONST, tokNCHARLIT, tokBIND:
		return ""
	}
	return ""
}

// IsReservedKeyword reports whether word is an Oracle SQL reserved word that
// cannot be used as a nonquoted identifier.
func IsReservedKeyword(word string) bool {
	if reverseOracleKeywordMap == nil {
		initReverseOracleKeywordMap()
	}
	tok, ok := oracleKeywords[strings.ToUpper(word)]
	if !ok {
		return false
	}
	return isOracleSQLReservedKeyword(Token{Type: tok, Str: strings.ToUpper(word)})
}

// Tokenize runs the Oracle lexer and returns all non-EOF tokens.
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

// RuleCandidate represents a grammar rule candidate at the completion cursor.
type RuleCandidate struct {
	Rule string
}

// CandidateSet holds token and rule candidates collected for completion.
type CandidateSet struct {
	Tokens []int
	Rules  []RuleCandidate
	seen   map[int]bool
	seenR  map[string]bool
}

func newCandidateSet() *CandidateSet {
	return &CandidateSet{
		seen:  make(map[int]bool),
		seenR: make(map[string]bool),
	}
}

func (cs *CandidateSet) addToken(t int) {
	if cs == nil || cs.seen[t] {
		return
	}
	cs.seen[t] = true
	cs.Tokens = append(cs.Tokens, t)
}

func (cs *CandidateSet) addRule(rule string) {
	if cs == nil || cs.seenR[rule] {
		return
	}
	cs.seenR[rule] = true
	cs.Rules = append(cs.Rules, RuleCandidate{Rule: rule})
}

// HasToken reports whether the candidate set contains a token type.
func (cs *CandidateSet) HasToken(t int) bool {
	return cs != nil && cs.seen[t]
}

// HasRule reports whether the candidate set contains a grammar rule.
func (cs *CandidateSet) HasRule(rule string) bool {
	return cs != nil && cs.seenR[rule]
}

var oracleTopLevelCompletionTokens = []int{
	kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwMERGE, kwCREATE, kwALTER, kwDROP,
	kwTRUNCATE, kwCOMMENT, kwGRANT, kwREVOKE, kwCOMMIT, kwROLLBACK, kwSAVEPOINT,
	kwSET, kwBEGIN, kwDECLARE, kwWITH, kwCALL, kwEXPLAIN, kwLOCK, kwRENAME,
}

// Collect returns parser-native completion candidates for sql at cursorOffset.
//
// The initial implementation is deliberately conservative: it provides stable
// top-level token candidates and a few structural rule signals. Later phases add
// deeper clause and scope instrumentation.
func Collect(sql string, cursorOffset int) *CandidateSet {
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}

	prefix := completionPrefix(sql, cursorOffset)
	collectOffset := cursorOffset - len(prefix)
	cs, _ := collectOracleCandidates(sql, Tokenize(sql), collectOffset, prefix)
	return cs
}

func collectOracleCandidates(_ string, allTokens []Token, collectOffset int, prefix string) (*CandidateSet, *CompletionIntent) {
	cs := newCandidateSet()
	tokens := tokensBeforeOffset(allTokens, collectOffset)
	if atOracleStatementStart(tokens) {
		addTopLevelCandidates(cs)
		return cs, &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindUnknown}}
	}

	// Partial top-level keyword, e.g. SEL|.
	if len(tokens) == 0 && prefix != "" {
		addTopLevelCandidates(cs)
		return cs, &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindUnknown}}
	}
	if ddlIntent, handled := inferOracleDDLCompletionIntent(cs, allTokens, tokens, collectOffset); handled {
		addRulesForOracleIntent(cs, ddlIntent)
		addOracleStructuralTokenCandidates(cs, allTokens, tokens, collectOffset)
		return cs, ddlIntent
	}
	intent := inferOracleCompletionIntent(allTokens, tokens, collectOffset)
	addRulesForOracleIntent(cs, intent)
	addOracleStructuralTokenCandidates(cs, allTokens, tokens, collectOffset)
	return cs, intent
}

func addRulesForOracleIntent(cs *CandidateSet, intent *CompletionIntent) {
	if intent == nil {
		return
	}
	for _, kind := range intent.ObjectKinds {
		switch kind {
		case ObjectKindTable, ObjectKindView:
			cs.addRule("table_ref")
		case ObjectKindColumn:
			cs.addRule("columnref")
		case ObjectKindSchema:
			cs.addRule("schema_ref")
		case ObjectKindSequence:
			cs.addRule("sequence_ref")
		case ObjectKindSequenceMember:
			cs.addRule("sequence_member_ref")
		case ObjectKindProcedure:
			cs.addRule("proc_ref")
		case ObjectKindFunction:
			cs.addRule("func_name")
			addOracleKeywordCandidatesBy(cs, isOracleFunctionKeyword)
		case ObjectKindPackage:
			cs.addRule("package_ref")
		case ObjectKindType:
			cs.addRule("type_name")
			addOracleKeywordCandidatesBy(cs, isOracleTypeKeyword)
		case ObjectKindIndex:
			cs.addRule("index_ref")
		case ObjectKindTrigger:
			cs.addRule("trigger_ref")
		case ObjectKindConstraint:
			cs.addRule("constraint_ref")
		case ObjectKindSynonym:
			cs.addRule("synonym_ref")
		case ObjectKindPackageMember:
			cs.addRule("package_member_ref")
		case ObjectKindDatabaseLink:
			cs.addRule("database_link_ref")
		case ObjectKindVariable:
			cs.addRule("variable_ref")
		}
		if kind == ObjectKindColumn {
			addOracleKeywordCandidatesBy(cs, isOraclePseudoColumnKeyword)
		}
	}
}

func addOracleKeywordCandidatesBy(cs *CandidateSet, pred func(int) bool) {
	tokens := make([]int, 0)
	for _, tok := range oracleKeywords {
		if pred(tok) {
			tokens = append(tokens, tok)
		}
	}
	sort.Slice(tokens, func(i, j int) bool {
		return TokenName(tokens[i]) < TokenName(tokens[j])
	})
	for _, tok := range tokens {
		cs.addToken(tok)
	}
}

func completionPrefix(sql string, cursorOffset int) string {
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	i := cursorOffset
	for i > 0 {
		c := sql[i-1]
		if isCompletionIdentByte(c) {
			i--
			continue
		}
		break
	}
	return sql[i:cursorOffset]
}

func isCompletionIdentByte(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_' || c == '$' || c == '#'
}

func tokensBeforeOffset(tokens []Token, offset int) []Token {
	result := tokens[:0]
	for _, tok := range tokens {
		if tok.Loc >= offset {
			break
		}
		result = append(result, tok)
	}
	return result
}

func atOracleStatementStart(tokens []Token) bool {
	for i := len(tokens) - 1; i >= 0; i-- {
		switch tokens[i].Type {
		case ';':
			return true
		case tokEOF:
			continue
		default:
			return false
		}
	}
	return true
}

func addTopLevelCandidates(cs *CandidateSet) {
	for _, tok := range oracleTopLevelCompletionTokens {
		cs.addToken(tok)
	}
}

func inferOracleCompletionIntent(allTokens []Token, before []Token, collectOffset int) *CompletionIntent {
	intent := &CompletionIntent{}
	if len(before) == 0 {
		return intent
	}

	if before[len(before)-1].Type == '@' {
		intent.ObjectKinds = []ObjectKind{ObjectKindDatabaseLink}
		return intent
	}

	if dmlIntent := inferOracleDMLCompletionIntent(allTokens, before, collectOffset); len(dmlIntent.ObjectKinds) > 0 {
		return dmlIntent
	}

	if q, ok := oracleDottedQualifierBefore(before); ok {
		if oracleCursorInTableRef(allTokens, collectOffset) {
			intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
			intent.Qualifier.Schema = q
			return intent
		}
		intent.ObjectKinds = []ObjectKind{ObjectKindColumn, ObjectKindPackageMember, ObjectKindSequenceMember}
		intent.Qualifier.Object = q
		return intent
	}

	prev := previousNonPunctuation(before)
	if oracleCursorInPLSQLDeclarationType(allTokens, collectOffset) {
		intent.ObjectKinds = []ObjectKind{ObjectKindType}
		return intent
	}
	if prev.Type == kwBEGIN || prev.Type == kwTHEN || prev.Type == kwELSE {
		intent.ObjectKinds = []ObjectKind{ObjectKindUnknown}
		return intent
	}
	if prev.Type == kwFROM || prev.Type == kwJOIN {
		intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
		return intent
	}
	if prev.Type == kwWHERE || prev.Type == kwON || prev.Type == kwBY || prev.Type == kwUSING {
		intent.ObjectKinds = []ObjectKind{ObjectKindColumn}
		return intent
	}
	if prev.Type == kwINTO && oracleCursorInSelectInto(allTokens, collectOffset) {
		intent.ObjectKinds = []ObjectKind{ObjectKindVariable}
		return intent
	}
	if prev.Type == kwEXISTS {
		intent.ObjectKinds = []ObjectKind{ObjectKindUnknown}
		return intent
	}
	if prev.Type == kwSELECT {
		intent.ObjectKinds = []ObjectKind{ObjectKindColumn, ObjectKindFunction}
		return intent
	}
	if oracleTokenStartsExpression(prev.Type) {
		intent.ObjectKinds = []ObjectKind{ObjectKindColumn, ObjectKindFunction}
		return intent
	}
	if oracleCursorInSelectTarget(allTokens, collectOffset) ||
		oracleCursorInExpressionClause(allTokens, collectOffset) ||
		oracleCursorInOrderingExpression(allTokens, collectOffset) {
		intent.ObjectKinds = []ObjectKind{ObjectKindColumn, ObjectKindFunction}
		return intent
	}
	if oracleCursorInTableRef(allTokens, collectOffset) {
		intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
		return intent
	}
	return intent
}

func inferOracleDDLCompletionIntent(cs *CandidateSet, allTokens []Token, before []Token, collectOffset int) (*CompletionIntent, bool) {
	intent := &CompletionIntent{}
	stmtStart, _ := oracleStatementTokenBounds(allTokens, collectOffset)
	if stmtStart >= len(allTokens) {
		return intent, false
	}
	first := allTokens[stmtStart]
	prev := previousNonPunctuation(before)
	if q, ok := oracleDottedQualifierBefore(before); ok {
		context := oracleKeywordBeforeDottedQualifier(before)
		switch first.Type {
		case kwDROP:
			switch context.Type {
			case kwTABLE:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTable, ObjectKindView}, Qualifier: MultipartName{Schema: q}}, true
			case kwVIEW:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindView}, Qualifier: MultipartName{Schema: q}}, true
			case kwSEQUENCE:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindSequence}, Qualifier: MultipartName{Schema: q}}, true
			case kwPROCEDURE:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindProcedure}, Qualifier: MultipartName{Schema: q}}, true
			case kwFUNCTION:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindFunction}, Qualifier: MultipartName{Schema: q}}, true
			case kwINDEX:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindIndex}, Qualifier: MultipartName{Schema: q}}, true
			case kwSYNONYM:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindSynonym}, Qualifier: MultipartName{Schema: q}}, true
			case kwTRIGGER:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTrigger}, Qualifier: MultipartName{Schema: q}}, true
			}
		case kwCOMMENT:
			switch context.Type {
			case kwCOLUMN:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindColumn}, Qualifier: MultipartName{Object: q}}, true
			case kwTABLE:
				return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTable, ObjectKindView}, Qualifier: MultipartName{Schema: q}}, true
			}
		case kwGRANT, kwREVOKE:
			if context.Type == kwON {
				return &CompletionIntent{ObjectKinds: oracleGrantRevokeObjectKinds(before[stmtStart:]), Qualifier: MultipartName{Schema: q}}, true
			}
		}
	}
	switch first.Type {
	case kwCREATE:
		if prev.Type == kwCREATE {
			for _, tok := range []int{kwTABLE, kwVIEW, kwINDEX, kwSEQUENCE, kwSYNONYM, kwPROCEDURE, kwFUNCTION, kwPACKAGE, kwTRIGGER, kwUSER, kwROLE, kwTYPE} {
				cs.addToken(tok)
			}
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindUnknown}}, true
		}
		if oracleCreateIndexTableContext(allTokens[stmtStart:], collectOffset) {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTable, ObjectKindView}}, true
		}
		if q, ok := oracleCreateIndexColumnContext(allTokens[stmtStart:], collectOffset); ok {
			return &CompletionIntent{
				ObjectKinds: []ObjectKind{ObjectKindColumn},
				Qualifier:   MultipartName{Object: q},
			}, true
		}
		if prev.Type == kwREFERENCES {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTable, ObjectKindView}}, true
		}
		if q, ok := oracleCreateTableColumnReferenceContext(allTokens[stmtStart:], collectOffset); ok {
			return &CompletionIntent{
				ObjectKinds: []ObjectKind{ObjectKindColumn},
				Qualifier:   MultipartName{Object: q},
			}, true
		}
		if oracleCreateTableConstraintColumnContext(allTokens[stmtStart:], collectOffset) {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindColumn}}, true
		}
		if oracleCreateTableColumnStartContext(allTokens[stmtStart:], collectOffset) {
			for _, tok := range []int{kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN, kwCHECK} {
				cs.addToken(tok)
			}
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindUnknown}}, true
		}
		if oracleCreateTableColumnOptionContext(allTokens[stmtStart:], collectOffset) {
			for _, tok := range []int{kwNOT, kwNULL, kwDEFAULT, kwPRIMARY, kwUNIQUE, kwREFERENCES, kwCHECK} {
				cs.addToken(tok)
			}
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindUnknown}}, true
		}
		if oracleCreateTableTypeContext(allTokens[stmtStart:], collectOffset) {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindType}}, true
		}
	case kwALTER:
		switch prev.Type {
		case kwSEQUENCE:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindSequence}}, true
		case kwVIEW:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindView}}, true
		case kwPROCEDURE:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindProcedure}}, true
		}
		if prev.Type == kwTABLE {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTable}}, true
		}
		if oracleSeenToken(before[stmtStart:], kwTABLE) && prev.Type == kwADD {
			for _, tok := range []int{kwCOLUMN, kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN, kwCHECK} {
				cs.addToken(tok)
			}
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindUnknown}}, true
		}
		if oracleSeenToken(before[stmtStart:], kwTABLE) && oracleSeenToken(before[stmtStart:], kwDROP) && prev.Type == kwCOLUMN {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindColumn}}, true
		}
		if oracleSeenToken(before[stmtStart:], kwTABLE) && prev.Type == kwMODIFY {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindColumn}}, true
		}
		if oracleSeenToken(before[stmtStart:], kwTABLE) && oracleSeenToken(before[stmtStart:], kwDROP) && prev.Type == kwCONSTRAINT {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindConstraint}}, true
		}
	case kwDROP:
		switch prev.Type {
		case kwTABLE:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTable}}, true
		case kwVIEW:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindView}}, true
		case kwSEQUENCE:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindSequence}}, true
		case kwPROCEDURE:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindProcedure}}, true
		case kwFUNCTION:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindFunction}}, true
		case kwINDEX:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindIndex}}, true
		case kwSYNONYM:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindSynonym}}, true
		case kwTRIGGER:
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTrigger}}, true
		}
	case kwTRUNCATE:
		if prev.Type == kwTABLE {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTable}}, true
		}
	case kwCOMMENT:
		if q, ok := oracleDottedQualifierBefore(before); ok && oracleSeenToken(before[stmtStart:], kwCOLUMN) {
			return &CompletionIntent{
				ObjectKinds: []ObjectKind{ObjectKindColumn},
				Qualifier:   MultipartName{Object: q},
			}, true
		}
		if prev.Type == kwTABLE {
			return &CompletionIntent{ObjectKinds: []ObjectKind{ObjectKindTable}}, true
		}
	case kwGRANT, kwREVOKE:
		if prev.Type == kwON {
			return &CompletionIntent{ObjectKinds: oracleGrantRevokeObjectKinds(before[stmtStart:])}, true
		}
	}
	return intent, false
}

func oracleGrantRevokeObjectKinds(tokens []Token) []ObjectKind {
	kinds := []ObjectKind{ObjectKindTable, ObjectKindView}
	for _, tok := range tokens {
		if tok.Type == kwON {
			break
		}
		switch tok.Type {
		case kwSELECT:
			return append(kinds, ObjectKindSequence)
		case kwEXECUTE:
			return []ObjectKind{ObjectKindProcedure, ObjectKindFunction, ObjectKindPackage, ObjectKindType}
		}
	}
	return kinds
}

func oracleCreateTableTypeContext(tokens []Token, offset int) bool {
	if len(tokens) < 2 || tokens[0].Type != kwCREATE || tokens[1].Type != kwTABLE {
		return false
	}
	for i := 2; i < len(tokens); i++ {
		if tokens[i].Loc >= offset {
			return false
		}
		if tokens[i].Type == '(' {
			closeIdx := oracleMatchingParen(tokens, i)
			return closeIdx < 0 || tokens[closeIdx].Loc >= offset
		}
	}
	return false
}

func oracleCreateTableColumnStartContext(tokens []Token, offset int) bool {
	if !oracleCreateTableOpenBefore(tokens, offset) {
		return false
	}
	return lastTokenBeforeOffset(tokens, offset).Type == ','
}

func oracleCreateTableColumnOptionContext(tokens []Token, offset int) bool {
	if !oracleCreateTableOpenBefore(tokens, offset) {
		return false
	}
	prev := previousNonPunctuation(tokensBeforeOffset(tokens, offset))
	return isOracleTypeKeyword(prev.Type) || prev.Type == tokICONST || prev.Type == ')'
}

func oracleCreateTableConstraintColumnContext(tokens []Token, offset int) bool {
	if !oracleCreateTableOpenBefore(tokens, offset) {
		return false
	}
	before := tokensBeforeOffset(tokens, offset)
	openIdx := lastUnclosedParenBeforeOffset(tokens, offset)
	if openIdx < 0 || openIdx >= len(before) {
		return false
	}
	for i := openIdx - 1; i >= 0; i-- {
		switch before[i].Type {
		case kwPRIMARY, kwFOREIGN:
			return true
		case kwREFERENCES, kwCHECK, ',':
			return false
		}
	}
	return false
}

func oracleCreateTableColumnReferenceContext(tokens []Token, offset int) (string, bool) {
	if !oracleCreateTableOpenBefore(tokens, offset) {
		return "", false
	}
	openIdx := lastUnclosedParenBeforeOffset(tokens, offset)
	if openIdx < 2 {
		return "", false
	}
	tableTok := tokens[openIdx-1]
	if !isOracleCompletionIdent(tableTok) {
		return "", false
	}
	for i := openIdx - 2; i >= 0; i-- {
		if tokens[i].Type == kwREFERENCES {
			return tableTok.Str, true
		}
		if tokens[i].Type == ',' || tokens[i].Type == kwPRIMARY || tokens[i].Type == kwFOREIGN {
			return "", false
		}
	}
	return "", false
}

func oracleCreateIndexTableContext(tokens []Token, offset int) bool {
	if len(tokens) < 2 || tokens[0].Type != kwCREATE || !oracleSeenToken(tokensBeforeOffset(tokens, offset), kwINDEX) {
		return false
	}
	return previousNonPunctuation(tokensBeforeOffset(tokens, offset)).Type == kwON
}

func oracleCreateIndexColumnContext(tokens []Token, offset int) (string, bool) {
	if len(tokens) < 2 || tokens[0].Type != kwCREATE || !oracleSeenToken(tokensBeforeOffset(tokens, offset), kwINDEX) {
		return "", false
	}
	openIdx := lastUnclosedParenBeforeOffset(tokens, offset)
	if openIdx < 1 {
		return "", false
	}
	tableTok := tokens[openIdx-1]
	if !isOracleCompletionIdent(tableTok) {
		return "", false
	}
	for i := openIdx - 2; i >= 0; i-- {
		if tokens[i].Type == kwON {
			return tableTok.Str, true
		}
	}
	return "", false
}

func oracleCreateTableOpenBefore(tokens []Token, offset int) bool {
	if len(tokens) < 2 || tokens[0].Type != kwCREATE || tokens[1].Type != kwTABLE {
		return false
	}
	return lastUnclosedParenBeforeOffset(tokens, offset) >= 0
}

func inferOracleDMLCompletionIntent(allTokens []Token, before []Token, collectOffset int) *CompletionIntent {
	intent := &CompletionIntent{}
	stmtStart, _ := oracleStatementTokenBounds(allTokens, collectOffset)
	if stmtStart >= len(allTokens) {
		return intent
	}
	first := allTokens[stmtStart]
	prev := previousNonPunctuation(before)
	if q, ok := oracleDottedQualifierBefore(before); ok {
		context := oracleKeywordBeforeDottedQualifier(before)
		switch first.Type {
		case kwINSERT:
			if context.Type == kwINTO {
				intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
				intent.Qualifier = MultipartName{Schema: q}
				return intent
			}
		case kwUPDATE:
			if context.Type == kwUPDATE {
				intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
				intent.Qualifier = MultipartName{Schema: q}
				return intent
			}
		case kwDELETE:
			if context.Type == kwFROM {
				intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
				intent.Qualifier = MultipartName{Schema: q}
				return intent
			}
		case kwMERGE:
			if context.Type == kwINTO || context.Type == kwUSING {
				intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
				intent.Qualifier = MultipartName{Schema: q}
				return intent
			}
		}
	}
	switch first.Type {
	case kwINSERT:
		if prev.Type == kwINTO {
			intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
			return intent
		}
		if oracleInsertColumnListContext(allTokens[stmtStart:], collectOffset) {
			intent.ObjectKinds = []ObjectKind{ObjectKindColumn}
			return intent
		}
		if oracleSeenToken(before[stmtStart:], kwVALUES) {
			intent.ObjectKinds = []ObjectKind{ObjectKindColumn, ObjectKindFunction}
			return intent
		}
	case kwUPDATE:
		if prev.Type == kwUPDATE {
			intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
			return intent
		}
		if oracleSeenToken(before[stmtStart:], kwSET) || prev.Type == kwWHERE {
			intent.ObjectKinds = []ObjectKind{ObjectKindColumn, ObjectKindFunction}
			return intent
		}
	case kwDELETE:
		if prev.Type == kwFROM {
			intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
			return intent
		}
		if prev.Type == kwWHERE || oracleSeenToken(before[stmtStart:], kwWHERE) {
			intent.ObjectKinds = []ObjectKind{ObjectKindColumn, ObjectKindFunction}
			return intent
		}
	case kwMERGE:
		if prev.Type == kwINTO || prev.Type == kwUSING {
			intent.ObjectKinds = []ObjectKind{ObjectKindTable, ObjectKindView}
			return intent
		}
		if prev.Type == kwON || oracleSeenToken(before[stmtStart:], kwON) {
			intent.ObjectKinds = []ObjectKind{ObjectKindColumn, ObjectKindFunction}
			return intent
		}
	}
	return intent
}

func oracleInsertColumnListContext(tokens []Token, offset int) bool {
	intoIdx := -1
	for i, tok := range tokens {
		if tok.Type == kwINTO {
			intoIdx = i
			break
		}
	}
	if intoIdx < 0 {
		return false
	}
	for i := intoIdx + 1; i < len(tokens); i++ {
		if tokens[i].Loc >= offset {
			return false
		}
		if tokens[i].Type == '(' {
			closeIdx := oracleMatchingParen(tokens, i)
			return closeIdx < 0 || tokens[closeIdx].Loc >= offset
		}
		if tokens[i].Type == kwVALUES || tokens[i].Type == kwSELECT {
			return false
		}
	}
	return false
}

func oracleSeenToken(tokens []Token, tokenType int) bool {
	for _, tok := range tokens {
		if tok.Type == tokenType {
			return true
		}
	}
	return false
}

func lastTokenBeforeOffset(tokens []Token, offset int) Token {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Loc < offset {
			return tokens[i]
		}
	}
	return Token{}
}

func previousNonPunctuation(tokens []Token) Token {
	for i := len(tokens) - 1; i >= 0; i-- {
		switch tokens[i].Type {
		case ',', '(', ')':
			continue
		default:
			return tokens[i]
		}
	}
	return Token{}
}

func oracleTokenStartsExpression(tokenType int) bool {
	switch tokenType {
	case kwAND, kwOR, kwBETWEEN, kwIN, kwTHEN, kwELSE, '+', '-', '*', '/', '=', '<', '>', '|':
		return true
	default:
		return false
	}
}

func oracleDottedQualifierBefore(tokens []Token) (string, bool) {
	if len(tokens) < 2 || tokens[len(tokens)-1].Type != '.' {
		return "", false
	}
	qualifier := tokens[len(tokens)-2]
	if !isOracleCompletionIdent(qualifier) {
		return "", false
	}
	return qualifier.Str, true
}

func oracleKeywordBeforeDottedQualifier(tokens []Token) Token {
	if _, ok := oracleDottedQualifierBefore(tokens); !ok {
		return Token{}
	}
	for i := len(tokens) - 3; i >= 0; i-- {
		switch tokens[i].Type {
		case ',', '(', ')', '.':
			continue
		default:
			return tokens[i]
		}
	}
	return Token{}
}

func oracleCursorInSelectTarget(tokens []Token, offset int) bool {
	stmtStart, stmtEnd := oracleStatementTokenBounds(tokens, offset)
	selectIdx := -1
	for i := stmtStart; i < stmtEnd; i++ {
		if tokens[i].Type == kwSELECT && tokens[i].Loc < offset {
			selectIdx = i
		}
	}
	if selectIdx < 0 {
		return false
	}
	for i := selectIdx + 1; i < stmtEnd; i++ {
		if tokens[i].Loc >= offset {
			break
		}
		if tokens[i].Type == kwFROM {
			return false
		}
	}
	for i := selectIdx + 1; i < stmtEnd; i++ {
		if tokens[i].Loc >= offset {
			return true
		}
		if tokens[i].Type == kwFROM {
			return true
		}
	}
	return false
}

func oracleCursorInExpressionClause(tokens []Token, offset int) bool {
	stmtStart, _ := oracleStatementTokenBounds(tokens, offset)
	for i := len(tokens) - 1; i >= stmtStart; i-- {
		if tokens[i].Loc >= offset {
			continue
		}
		switch tokens[i].Type {
		case kwWHERE, kwON, kwHAVING:
			return true
		case kwSTART, kwCONNECT:
			return true
		case kwFROM, kwJOIN, kwGROUP, kwORDER:
			return false
		}
	}
	return false
}

func oracleCursorInOrderingExpression(tokens []Token, offset int) bool {
	stmtStart, _ := oracleStatementTokenBounds(tokens, offset)
	before := tokensBeforeOffset(tokens[stmtStart:], offset)
	for i := len(before) - 1; i >= 0; i-- {
		switch before[i].Type {
		case kwGROUP, kwORDER:
			return i+1 < len(before) && before[i+1].Type == kwBY
		case kwWHERE, kwFROM, kwJOIN, kwHAVING, kwON, kwUNION, kwINTERSECT, kwMINUS:
			return false
		}
	}
	return false
}

func oracleCursorInSelectInto(tokens []Token, offset int) bool {
	stmtStart, stmtEnd := oracleStatementTokenBounds(tokens, offset)
	seenSelect := false
	for i := stmtStart; i < stmtEnd; i++ {
		if tokens[i].Loc >= offset {
			break
		}
		switch tokens[i].Type {
		case kwSELECT:
			seenSelect = true
		case kwFROM:
			return false
		case kwINTO:
			return seenSelect
		}
	}
	return false
}

func oracleCursorInPLSQLDeclarationType(tokens []Token, offset int) bool {
	stmtStart, stmtEnd := oracleStatementTokenBounds(tokens, offset)
	if stmtStart >= stmtEnd || tokens[stmtStart].Type != kwDECLARE {
		return false
	}
	before := tokensBeforeOffset(tokens[stmtStart:stmtEnd], offset)
	for _, tok := range before {
		if tok.Type == kwBEGIN {
			return false
		}
	}
	prev := previousNonPunctuation(before)
	return isOracleCompletionIdent(prev)
}

func oracleCursorInTableRef(tokens []Token, offset int) bool {
	stmtStart, stmtEnd := oracleStatementTokenBounds(tokens, offset)
	depth := 0
	lastClause := 0
	for i := stmtStart; i < stmtEnd; i++ {
		if tokens[i].Loc >= offset {
			break
		}
		switch tokens[i].Type {
		case '(':
			depth++
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		switch tokens[i].Type {
		case kwFROM, kwJOIN:
			lastClause = tokens[i].Type
		case kwWHERE, kwON, kwGROUP, kwHAVING, kwORDER, kwCONNECT, kwSTART,
			kwUNION, kwINTERSECT, kwMINUS, kwOFFSET, kwFETCH:
			lastClause = tokens[i].Type
		}
	}
	return lastClause == kwFROM || lastClause == kwJOIN
}

func addOracleStructuralTokenCandidates(cs *CandidateSet, allTokens []Token, before []Token, offset int) {
	prev := previousNonPunctuation(before)
	switch prev.Type {
	case kwEXISTS:
		cs.addToken(kwSELECT)
	case kwBEGIN, kwTHEN, kwELSE:
		for _, tok := range []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwMERGE, kwNULL, kwBEGIN} {
			cs.addToken(tok)
		}
	}
	if oracleCursorAfterFromRelation(allTokens, offset) {
		for _, tok := range []int{kwWHERE, kwJOIN, kwGROUP, kwORDER, kwHAVING} {
			cs.addToken(tok)
		}
	}
}

func oracleCursorAfterFromRelation(tokens []Token, offset int) bool {
	stmtStart, stmtEnd := oracleStatementTokenBounds(tokens, offset)
	depth := 0
	lastFromOrJoin := -1
	lastBoundary := -1
	for i := stmtStart; i < stmtEnd; i++ {
		if tokens[i].Loc >= offset {
			break
		}
		switch tokens[i].Type {
		case '(':
			depth++
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		switch tokens[i].Type {
		case kwFROM, kwJOIN:
			lastFromOrJoin = i
			lastBoundary = -1
		case kwWHERE, kwON, kwUSING, kwGROUP, kwHAVING, kwORDER, kwUNION, kwINTERSECT, kwMINUS:
			lastBoundary = i
		}
	}
	if lastFromOrJoin < 0 || lastBoundary > lastFromOrJoin {
		return false
	}
	_, _, ok := parseOracleRangeReference(tokens, lastFromOrJoin+1, nil)
	return ok
}

func lastUnclosedParenBeforeOffset(tokens []Token, offset int) int {
	var stack []int
	for i, tok := range tokens {
		if tok.Loc >= offset {
			break
		}
		switch tok.Type {
		case '(':
			stack = append(stack, i)
		case ')':
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}
	if len(stack) == 0 {
		return -1
	}
	return stack[len(stack)-1]
}
