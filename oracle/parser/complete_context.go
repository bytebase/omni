package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// CompletionContext is Oracle parser-native completion state.
type CompletionContext struct {
	Candidates *CandidateSet
	Scope      *ScopeSnapshot
	CTEs       []RangeReference
	Prefix     string
	Intent     *CompletionIntent
}

// CompletionIntent describes the object class and qualifier implied by the
// cursor position.
type CompletionIntent struct {
	ObjectKinds []ObjectKind
	Qualifier   MultipartName
}

// ObjectKind classifies catalog object completion intent.
type ObjectKind int

const (
	ObjectKindUnknown ObjectKind = iota
	ObjectKindSchema
	ObjectKindTable
	ObjectKindView
	ObjectKindSequence
	ObjectKindProcedure
	ObjectKindFunction
	ObjectKindPackage
	ObjectKindType
	ObjectKindColumn
	ObjectKindIndex
	ObjectKindTrigger
	ObjectKindConstraint
	ObjectKindSynonym
	ObjectKindSequenceMember
	ObjectKindPackageMember
	ObjectKindDatabaseLink
	ObjectKindVariable
	ObjectKindUser
	ObjectKindRole
)

// MultipartName is a partially typed Oracle object name.
type MultipartName struct {
	Schema string
	Object string
}

// ScopeSnapshot describes parser-known relation scope visible at the cursor.
type ScopeSnapshot struct {
	// LocalReferences is the current query block's visible FROM/JOIN scope.
	LocalReferences []RangeReference
	// OuterReferences contains correlated outer query-block scopes, nearest first.
	OuterReferences [][]RangeReference
	// References contains all statement references for debugging/backward
	// compatibility. Callers should not use it for column completion because it
	// can include CTE body tables and other refs outside the cursor query block.
	References []RangeReference
}

// RangeReferenceKind classifies a range-table entry visible to completion.
type RangeReferenceKind int

const (
	RangeReferenceRelation RangeReferenceKind = iota
	RangeReferenceSubquery
	RangeReferenceJoinAlias
	RangeReferenceCTE
	RangeReferenceDMLTarget
	RangeReferenceMergeSource
)

// RangeReference is a syntax-level table expression reference. It deliberately
// avoids catalog metadata; Bytebase resolves metadata from its own store.
type RangeReference struct {
	Kind   RangeReferenceKind
	Schema string
	Name   string
	Alias  string

	Columns []string

	Loc     nodes.Loc
	BodyLoc nodes.Loc
}

// CollectCompletion returns parser candidates plus best-effort visible scope
// and object intent at cursorOffset.
func CollectCompletion(sql string, cursorOffset int) *CompletionContext {
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	prefix := completionPrefix(sql, cursorOffset)
	collectOffset := cursorOffset - len(prefix)
	tokens := Tokenize(sql)
	candidates, intent := collectOracleCandidates(sql, tokens, collectOffset, prefix)
	scope, ctes := collectOracleScopeAndCTEs(tokens, cursorOffset)
	return &CompletionContext{
		Candidates: candidates,
		Scope:      scope,
		CTEs:       ctes,
		Prefix:     prefix,
		Intent:     intent,
	}
}

func collectOracleScope(tokens []Token, cursorOffset int) *ScopeSnapshot {
	scope, _ := collectOracleScopeAndCTEs(tokens, cursorOffset)
	return scope
}

func collectOracleScopeAndCTEs(tokens []Token, cursorOffset int) (*ScopeSnapshot, []RangeReference) {
	start, end := oracleStatementTokenBounds(tokens, cursorOffset)
	stmtTokens := tokens[start:end]
	ctes := collectOracleCTEs(stmtTokens)
	cteMap := oracleRangeReferenceMap(ctes)
	selectRefs := collectOracleSelectRefs(stmtTokens, cteMap)
	dmlRefs := collectOracleDMLRefs(stmtTokens)
	ddlRefs := collectOracleDDLRefs(stmtTokens)
	refs := append([]RangeReference{}, selectRefs...)
	refs = append(refs, dmlRefs...)
	refs = append(refs, ddlRefs...)
	refs = append(append([]RangeReference{}, ctes...), refs...)
	localRefs := refs
	outerRefs := [][]RangeReference(nil)
	if selectIdx, selectDepth, ok := innermostOracleSelectBeforeCursor(stmtTokens, cursorOffset); ok {
		armStart, armEnd := oracleSelectArmBounds(stmtTokens, selectIdx, selectDepth)
		localRefs = collectOracleSelectRefsInRange(stmtTokens, armStart, armEnd, selectDepth, cteMap)
		outerRefs = collectOracleOuterSelectRefs(stmtTokens, cursorOffset, selectDepth, cteMap)
	} else if len(dmlRefs) > 0 {
		localRefs = dmlRefs
	} else if len(ddlRefs) > 0 {
		localRefs = ddlRefs
	}
	return &ScopeSnapshot{
		LocalReferences: localRefs,
		OuterReferences: outerRefs,
		References:      refs,
	}, ctes
}

func collectOracleDDLRefs(tokens []Token) []RangeReference {
	first := firstOracleNonTerminator(tokens)
	if first < 0 {
		return nil
	}
	switch tokens[first].Type {
	case kwCREATE:
		var refs []RangeReference
		if first+2 < len(tokens) && tokens[first+1].Type == kwTABLE {
			_, ref, ok := parseOracleRangeReference(tokens, first+2)
			if ok {
				refs = append(refs, ref)
			}
		}
		for i := first + 1; i < len(tokens); i++ {
			if tokens[i].Type == kwON {
				_, ref, ok := parseOracleRangeReference(tokens, i+1)
				if ok {
					refs = append(refs, ref)
				}
			}
			if tokens[i].Type == kwREFERENCES {
				_, ref, ok := parseOracleRangeReference(tokens, i+1)
				if ok {
					refs = append(refs, ref)
				}
			}
		}
		return refs
	case kwALTER:
		for i := first + 1; i < len(tokens); i++ {
			if tokens[i].Type == kwTABLE {
				_, ref, ok := parseOracleRangeReference(tokens, i+1)
				if ok {
					return []RangeReference{ref}
				}
				break
			}
		}
	case kwCOMMENT:
		for i := first + 1; i+3 < len(tokens); i++ {
			if tokens[i].Type == kwON && tokens[i+1].Type == kwCOLUMN && isOracleCompletionIdent(tokens[i+2]) && tokens[i+3].Type == '.' {
				ref := RangeReference{
					Kind:  RangeReferenceRelation,
					Name:  tokens[i+2].Str,
					Alias: tokens[i+2].Str,
					Loc:   nodes.Loc{Start: tokens[i+2].Loc, End: tokens[i+2].End},
				}
				return []RangeReference{ref}
			}
		}
	}
	return nil
}

func collectOracleDMLRefs(tokens []Token) []RangeReference {
	first := firstOracleNonTerminator(tokens)
	if first < 0 {
		return nil
	}
	switch tokens[first].Type {
	case kwUPDATE:
		_, ref, ok := parseOracleRangeReference(tokens, first+1)
		if ok {
			ref.Kind = RangeReferenceDMLTarget
			return []RangeReference{ref}
		}
	case kwINSERT:
		for i := first + 1; i < len(tokens); i++ {
			if tokens[i].Type == kwINTO {
				_, ref, ok := parseOracleRangeReference(tokens, i+1)
				if ok {
					ref.Kind = RangeReferenceDMLTarget
					return []RangeReference{ref}
				}
				break
			}
		}
	case kwDELETE:
		for i := first + 1; i < len(tokens); i++ {
			if tokens[i].Type == kwFROM {
				_, ref, ok := parseOracleRangeReference(tokens, i+1)
				if ok {
					ref.Kind = RangeReferenceDMLTarget
					return []RangeReference{ref}
				}
				break
			}
		}
	case kwMERGE:
		var refs []RangeReference
		for i := first + 1; i < len(tokens); i++ {
			switch tokens[i].Type {
			case kwINTO:
				next, ref, ok := parseOracleRangeReference(tokens, i+1)
				if ok {
					ref.Kind = RangeReferenceDMLTarget
					refs = append(refs, ref)
					i = next - 1
				}
			case kwUSING:
				next, ref, ok := parseOracleRangeReference(tokens, i+1)
				if ok {
					ref.Kind = RangeReferenceMergeSource
					refs = append(refs, ref)
					i = next - 1
				}
			}
		}
		return refs
	}
	return nil
}

func oracleStatementTokenBounds(tokens []Token, cursorOffset int) (int, int) {
	start := 0
	for i, tok := range tokens {
		if tok.Loc >= cursorOffset {
			break
		}
		if tok.Type == ';' {
			start = i + 1
		}
	}
	end := len(tokens)
	for i := start; i < len(tokens); i++ {
		if tokens[i].Loc < cursorOffset {
			continue
		}
		if tokens[i].Type == ';' {
			end = i
			break
		}
	}
	return start, end
}

func collectOracleSelectRefs(tokens []Token, ctes map[string]RangeReference) []RangeReference {
	return collectOracleSelectRefsAtDepth(tokens, -1, ctes)
}

func collectOracleSelectRefsAtDepth(tokens []Token, wantDepth int, ctes map[string]RangeReference) []RangeReference {
	var refs []RangeReference
	depth := 0
	for i, tok := range tokens {
		switch tok.Type {
		case '(':
			depth++
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			continue
		}
		if tok.Type != kwSELECT {
			continue
		}
		if wantDepth >= 0 && depth != wantDepth {
			continue
		}
		start, end := oracleSelectArmBounds(tokens, i, depth)
		refs = append(refs, collectOracleSelectRefsInRange(tokens, start, end, depth, ctes)...)
	}
	return refs
}

func collectOracleSelectRefsInRange(tokens []Token, start int, end int, wantDepth int, ctes map[string]RangeReference) []RangeReference {
	var refs []RangeReference
	depth := 0
	for i := 0; i < len(tokens); i++ {
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
		if wantDepth >= 0 && depth != wantDepth {
			continue
		}
		if i < start || i >= end {
			continue
		}
		switch tokens[i].Type {
		case kwFROM, kwJOIN:
			next, ref, ok := parseOracleRangeReference(tokens, i+1)
			if ok {
				ref = resolveOracleCTEReference(ref, ctes)
				refs = append(refs, ref)
				i = next - 1
			}
		}
	}
	return refs
}

func innermostOracleSelectBeforeCursor(tokens []Token, cursorOffset int) (int, int, bool) {
	depth := 0
	bestIdx := -1
	bestDepth := -1
	for i, tok := range tokens {
		if tok.Loc >= cursorOffset {
			break
		}
		switch tok.Type {
		case '(':
			depth++
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			continue
		case kwSELECT:
			bestIdx = i
			bestDepth = depth
		}
	}
	if bestIdx < 0 {
		return 0, 0, false
	}
	return bestIdx, bestDepth, true
}

func oracleSelectArmBounds(tokens []Token, selectIdx int, selectDepth int) (int, int) {
	start := selectIdx
	depth := 0
	for i := 0; i < len(tokens); i++ {
		if i >= selectIdx {
			break
		}
		switch tokens[i].Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case kwUNION, kwINTERSECT, kwMINUS:
			if depth == selectDepth {
				start = i + 1
			}
		}
	}
	end := len(tokens)
	depth = 0
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		if i > selectIdx && depth < selectDepth {
			end = i
			break
		}
		if i > selectIdx && depth == selectDepth {
			switch tok.Type {
			case kwUNION, kwINTERSECT, kwMINUS:
				end = i
				return start, end
			}
		}
		switch tok.Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		}
	}
	return start, end
}

func collectOracleOuterSelectRefs(tokens []Token, cursorOffset int, innerDepth int, ctes map[string]RangeReference) [][]RangeReference {
	if innerDepth <= 0 {
		return nil
	}
	type selectAtDepth struct {
		idx   int
		depth int
	}
	latest := make(map[int]selectAtDepth)
	depth := 0
	for i, tok := range tokens {
		if tok.Loc >= cursorOffset {
			break
		}
		switch tok.Type {
		case '(':
			depth++
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			continue
		case kwSELECT:
			if depth < innerDepth {
				latest[depth] = selectAtDepth{idx: i, depth: depth}
			}
		}
	}
	var levels [][]RangeReference
	for d := innerDepth - 1; d >= 0; d-- {
		sel, ok := latest[d]
		if !ok {
			continue
		}
		start, end := oracleSelectArmBounds(tokens, sel.idx, sel.depth)
		refs := collectOracleSelectRefsInRange(tokens, start, end, sel.depth, ctes)
		if len(refs) > 0 {
			levels = append(levels, refs)
		}
	}
	return levels
}

func oracleRangeReferenceMap(refs []RangeReference) map[string]RangeReference {
	if len(refs) == 0 {
		return nil
	}
	result := make(map[string]RangeReference, len(refs))
	for _, ref := range refs {
		if ref.Name != "" {
			result[strings.ToLower(ref.Name)] = ref
		}
	}
	return result
}

func resolveOracleCTEReference(ref RangeReference, ctes map[string]RangeReference) RangeReference {
	if ref.Kind != RangeReferenceRelation || ref.Schema != "" || ref.Name == "" {
		return ref
	}
	cte, ok := ctes[strings.ToLower(ref.Name)]
	if !ok {
		return ref
	}
	ref.Kind = RangeReferenceCTE
	ref.Columns = append([]string{}, cte.Columns...)
	ref.BodyLoc = cte.BodyLoc
	return ref
}

func parseOracleRangeReference(tokens []Token, i int) (int, RangeReference, bool) {
	for i < len(tokens) && tokens[i].Type == ',' {
		i++
	}
	if i >= len(tokens) {
		return i, RangeReference{}, false
	}
	if tokens[i].Type == '(' {
		return parseOracleSubqueryRangeReference(tokens, i)
	}
	if !isOracleCompletionIdent(tokens[i]) {
		return i, RangeReference{}, false
	}

	startTok := tokens[i]
	parts := []Token{tokens[i]}
	i++
	for i+1 < len(tokens) && tokens[i].Type == '.' && isOracleCompletionIdent(tokens[i+1]) {
		parts = append(parts, tokens[i+1])
		i += 2
	}

	ref := RangeReference{
		Kind: RangeReferenceRelation,
		Loc:  nodes.Loc{Start: startTok.Loc, End: parts[len(parts)-1].End},
	}
	switch len(parts) {
	case 1:
		ref.Name = parts[0].Str
	default:
		ref.Schema = parts[len(parts)-2].Str
		ref.Name = parts[len(parts)-1].Str
	}

	if i < len(tokens) && tokens[i].Type == kwAS {
		i++
	}
	if i < len(tokens) && isOracleCompletionIdent(tokens[i]) && !isOracleCompletionClauseBoundary(tokens[i].Type) {
		ref.Alias = tokens[i].Str
		ref.Loc.End = tokens[i].End
		i++
	}
	if ref.Alias == "" {
		ref.Alias = ref.Name
	}
	return i, ref, true
}

func parseOracleSubqueryRangeReference(tokens []Token, i int) (int, RangeReference, bool) {
	closeIdx := oracleMatchingParen(tokens, i)
	if closeIdx < 0 {
		return i, RangeReference{}, false
	}
	aliasIdx := closeIdx + 1
	if aliasIdx < len(tokens) && tokens[aliasIdx].Type == kwAS {
		aliasIdx++
	}
	if aliasIdx >= len(tokens) || !isOracleCompletionIdent(tokens[aliasIdx]) {
		return closeIdx + 1, RangeReference{}, false
	}
	ref := RangeReference{
		Kind:    RangeReferenceSubquery,
		Name:    tokens[aliasIdx].Str,
		Alias:   tokens[aliasIdx].Str,
		Columns: extractOracleSelectListColumns(tokens[i+1 : closeIdx]),
		Loc:     nodes.Loc{Start: tokens[i].Loc, End: tokens[aliasIdx].End},
		BodyLoc: nodes.Loc{Start: tokens[i].End, End: tokens[closeIdx].Loc},
	}
	return aliasIdx + 1, ref, true
}

func collectOracleCTEs(tokens []Token) []RangeReference {
	first := firstOracleNonTerminator(tokens)
	if first < 0 || tokens[first].Type != kwWITH {
		return nil
	}
	var refs []RangeReference
	i := first + 1
	for i < len(tokens) {
		if !isOracleCompletionIdent(tokens[i]) {
			break
		}
		nameTok := tokens[i]
		ref := RangeReference{
			Kind:  RangeReferenceCTE,
			Name:  nameTok.Str,
			Alias: nameTok.Str,
			Loc:   nodes.Loc{Start: nameTok.Loc, End: nameTok.End},
		}
		i++
		if i < len(tokens) && tokens[i].Type == '(' {
			closeIdx := oracleMatchingParen(tokens, i)
			if closeIdx < 0 {
				break
			}
			ref.Columns = extractOracleNameList(tokens[i+1 : closeIdx])
			ref.Loc.End = tokens[closeIdx].End
			i = closeIdx + 1
		}
		if i >= len(tokens) || tokens[i].Type != kwAS {
			break
		}
		i++
		if i >= len(tokens) || tokens[i].Type != '(' {
			break
		}
		closeIdx := oracleMatchingParen(tokens, i)
		if closeIdx < 0 {
			break
		}
		if len(ref.Columns) == 0 {
			ref.Columns = extractOracleSelectListColumns(tokens[i+1 : closeIdx])
		}
		ref.BodyLoc = nodes.Loc{Start: tokens[i].End, End: tokens[closeIdx].Loc}
		ref.Loc.End = tokens[closeIdx].End
		refs = append(refs, ref)
		i = closeIdx + 1
		if i < len(tokens) && tokens[i].Type == ',' {
			i++
			continue
		}
		break
	}
	return refs
}

func firstOracleNonTerminator(tokens []Token) int {
	for i, tok := range tokens {
		if tok.Type != ';' {
			return i
		}
	}
	return -1
}

func oracleMatchingParen(tokens []Token, openIdx int) int {
	if openIdx < 0 || openIdx >= len(tokens) || tokens[openIdx].Type != '(' {
		return -1
	}
	depth := 0
	for i := openIdx; i < len(tokens); i++ {
		switch tokens[i].Type {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func extractOracleNameList(tokens []Token) []string {
	var names []string
	for _, tok := range tokens {
		if isOracleCompletionIdent(tok) {
			names = append(names, tok.Str)
		}
	}
	return names
}

func extractOracleSelectListColumns(tokens []Token) []string {
	selectIdx := -1
	for i, tok := range tokens {
		if tok.Type == kwSELECT {
			selectIdx = i
			break
		}
	}
	if selectIdx < 0 {
		return nil
	}
	var names []string
	depth := 0
	for i := selectIdx + 1; i < len(tokens); i++ {
		tok := tokens[i]
		switch tok.Type {
		case '(':
			depth++
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth == 0 && tok.Type == kwFROM {
			break
		}
		if depth == 0 && isOracleCompletionIdent(tok) {
			names = append(names, tok.Str)
			for i+1 < len(tokens) && tokens[i+1].Type != ',' && tokens[i+1].Type != kwFROM {
				i++
			}
		}
	}
	return names
}

func isOracleCompletionIdent(tok Token) bool {
	if tok.Type == tokIDENT || tok.Type == tokQIDENT {
		return true
	}
	return tok.Type >= 2000 && !isOracleSQLReservedKeyword(tok)
}

func isOracleCompletionClauseBoundary(tokenType int) bool {
	switch tokenType {
	case kwWHERE, kwGROUP, kwHAVING, kwORDER, kwCONNECT, kwSTART, kwUNION,
		kwINTERSECT, kwMINUS, kwOFFSET, kwFETCH, kwJOIN, kwON, kwUSING,
		kwLEFT, kwRIGHT, kwFULL, kwCROSS, kwINNER, kwOUTER:
		return true
	default:
		return false
	}
}
