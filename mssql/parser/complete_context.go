package parser

import nodes "github.com/bytebase/omni/mssql/ast"

// CompletionContext is the parser-native context for SQL completion.
//
// Candidates contains grammar candidates at the cursor. Scope is a best-effort
// snapshot of range references visible to expression completion at the cursor.
// This is completion-only state; ordinary Parse remains strict.
type CompletionContext struct {
	Candidates *CandidateSet
	Scope      *ScopeSnapshot
	CTEs       []RangeReference
	Prefix     string
	Intent     *CompletionIntent
}

// CompletionIntent describes the object class and qualifier implied by the
// grammar position at the cursor.
type CompletionIntent struct {
	ObjectKinds []ObjectKind
	Qualifier   MultipartName
}

// ObjectKind classifies catalog object completion intent.
type ObjectKind int

const (
	ObjectKindUnknown ObjectKind = iota
	ObjectKindDatabase
	ObjectKindSchema
	ObjectKindTable
	ObjectKindView
	ObjectKindSequence
	ObjectKindProcedure
	ObjectKindFunction
	ObjectKindType
	ObjectKindColumn
)

// MultipartName is a partially typed MSSQL object name.
type MultipartName struct {
	Server   string
	Database string
	Schema   string
	Object   string
}

// ScopeSnapshot describes relation scope visible at a cursor.
type ScopeSnapshot struct {
	References      []RangeReference
	LocalReferences []RangeReference
	OuterReferences [][]RangeReference

	DMLTarget   *RangeReference
	MergeTarget *RangeReference
	MergeSource *RangeReference
}

// RangeReferenceKind classifies a range-table entry visible to completion.
type RangeReferenceKind int

const (
	RangeReferenceRelation RangeReferenceKind = iota
	RangeReferenceSubquery
	RangeReferenceFunction
	RangeReferenceJoinAlias
	RangeReferenceCTE
	RangeReferenceValues
	RangeReferenceTableVariable
	RangeReferenceDMLTarget
	RangeReferenceMergeTarget
	RangeReferenceMergeSource
)

// RangeReference is a syntax-level table expression reference. It deliberately
// avoids catalog metadata; callers resolve relation kinds and columns
// themselves.
type RangeReference struct {
	Kind RangeReferenceKind

	Server   string
	Database string
	Schema   string
	Object   string
	Alias    string

	Columns      []string
	AliasColumns []string

	Loc     nodes.Loc
	BodyLoc nodes.Loc
}

// CollectCompletion returns completion candidates plus a best-effort visible
// relation scope at cursorOffset.
func CollectCompletion(sql string, cursorOffset int) *CompletionContext {
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	prefix := completionPrefix(sql, cursorOffset)
	collectOffset := cursorOffset - len(prefix)
	candidates := Collect(sql, collectOffset)
	scope, ctes := collectCompletionScopeAndCTEs(sql, cursorOffset)
	return &CompletionContext{
		Candidates: candidates,
		Scope:      scope,
		CTEs:       ctes,
		Prefix:     prefix,
		Intent:     collectCompletionIntent(candidates, Tokenize(sql), collectOffset),
	}
}

type completionScopeToken struct {
	Token
	depth int
}

func collectCompletionScope(sql string, cursorOffset int) *ScopeSnapshot {
	scope, _ := collectCompletionScopeAndCTEs(sql, cursorOffset)
	return scope
}

func collectCompletionScopeAndCTEs(sql string, cursorOffset int) (*ScopeSnapshot, []RangeReference) {
	tokens := completionScopeTokens(Tokenize(sql))
	start, end := statementTokenBounds(tokens, cursorOffset)
	if start >= end {
		return &ScopeSnapshot{}, nil
	}

	stmtTokens := tokens[start:end]
	ctes, mainStart := collectCTEs(stmtTokens)
	cteMap := rangeReferenceMap(ctes)
	if scope := collectDMLScope(stmtTokens[mainStart:], cteMap); scope != nil {
		return scope, ctes
	}
	return collectSelectScope(stmtTokens[mainStart:], cursorOffset, cteMap), ctes
}

func completionScopeTokens(tokens []Token) []completionScopeToken {
	depth := 0
	result := make([]completionScopeToken, 0, len(tokens))
	for _, tok := range tokens {
		result = append(result, completionScopeToken{Token: tok, depth: depth})
		switch tok.Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		}
	}
	return result
}

func statementTokenBounds(tokens []completionScopeToken, cursorOffset int) (int, int) {
	start := 0
	for i, tok := range tokens {
		if tok.Loc >= cursorOffset {
			break
		}
		if tok.Type == ';' || tok.Type == kwGO {
			start = i + 1
		}
	}
	end := len(tokens)
	for i := start; i < len(tokens); i++ {
		if tokens[i].Loc < cursorOffset {
			continue
		}
		if tokens[i].Type == ';' || tokens[i].Type == kwGO {
			end = i
			break
		}
	}
	return start, end
}

func collectCTEs(tokens []completionScopeToken) ([]RangeReference, int) {
	first := firstNonTerminator(tokens)
	if first < 0 || tokens[first].Type != kwWITH {
		return nil, 0
	}
	var ctes []RangeReference
	idx := first + 1
	depth := tokens[first].depth
	for idx < len(tokens) {
		if tokens[idx].depth != depth || !IsIdentTokenType(tokens[idx].Type) {
			break
		}
		nameIdx := idx
		name := tokens[idx].Str
		idx++
		var columns []string
		if idx < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == '(' {
			closeIdx := matchingParen(tokens, idx)
			if closeIdx < 0 {
				break
			}
			for i := idx + 1; i < closeIdx; i++ {
				if IsIdentTokenType(tokens[i].Type) {
					columns = append(columns, tokens[i].Str)
				}
			}
			idx = closeIdx + 1
		}
		if idx >= len(tokens) || tokens[idx].Type != kwAS {
			break
		}
		idx++
		if idx >= len(tokens) || tokens[idx].Type != '(' {
			break
		}
		closeIdx := matchingParen(tokens, idx)
		if closeIdx < 0 {
			break
		}
		ctes = append(ctes, RangeReference{
			Kind:    RangeReferenceCTE,
			Object:  name,
			Columns: columns,
			Loc:     tokenLoc(tokens[nameIdx], tokens[closeIdx]),
			BodyLoc: tokenLoc(tokens[idx], tokens[closeIdx]),
		})
		idx = closeIdx + 1
		if idx < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == ',' {
			idx++
			continue
		}
		break
	}
	return ctes, idx
}

func rangeReferenceMap(refs []RangeReference) map[string]RangeReference {
	if len(refs) == 0 {
		return nil
	}
	result := make(map[string]RangeReference, len(refs))
	for _, ref := range refs {
		result[foldName(ref.Object)] = ref
	}
	return result
}

func collectDMLScope(tokens []completionScopeToken, ctes map[string]RangeReference) *ScopeSnapshot {
	first := firstNonTerminator(tokens)
	if first < 0 {
		return nil
	}
	switch tokens[first].Type {
	case kwUPDATE:
		return collectUpdateScope(tokens, first, ctes)
	case kwDELETE:
		return collectDeleteScope(tokens, first, ctes)
	case kwMERGE:
		return collectMergeScope(tokens, first, ctes)
	default:
		return nil
	}
}

func collectUpdateScope(tokens []completionScopeToken, updateIdx int, ctes map[string]RangeReference) *ScopeSnapshot {
	idx := skipTopClause(tokens, updateIdx+1, tokens[updateIdx].depth)
	target, next, ok := parseRangeReference(tokens, idx)
	if !ok {
		return nil
	}
	target.Kind = RangeReferenceDMLTarget
	target.Loc = tokenLoc(tokens[idx], tokens[next-1])
	target.Alias, target.AliasColumns, next = parseAliasAndColumns(tokens, next, tokens[idx].depth)
	local := []RangeReference{target}
	if fromIdx := findKeywordAtDepth(tokens, kwFROM, next, len(tokens), tokens[updateIdx].depth); fromIdx >= 0 {
		local = append(local, parseFromReferences(tokens, fromIdx+1, len(tokens), tokens[fromIdx].depth, ctes)...)
	}
	return scopeWithDMLTarget(local, target)
}

func collectDeleteScope(tokens []completionScopeToken, deleteIdx int, ctes map[string]RangeReference) *ScopeSnapshot {
	idx := skipTopClause(tokens, deleteIdx+1, tokens[deleteIdx].depth)
	if idx < len(tokens) && tokens[idx].Type == kwFROM {
		idx++
	}
	target, next, ok := parseRangeReference(tokens, idx)
	if !ok {
		return nil
	}
	target.Kind = RangeReferenceDMLTarget
	target.Loc = tokenLoc(tokens[idx], tokens[next-1])
	target.Alias, target.AliasColumns, next = parseAliasAndColumns(tokens, next, tokens[idx].depth)
	local := []RangeReference{target}
	// DELETE alias FROM ... has a second FROM that carries the real range scope.
	if fromIdx := findKeywordAtDepth(tokens, kwFROM, next, len(tokens), tokens[deleteIdx].depth); fromIdx >= 0 {
		local = append(local, parseFromReferences(tokens, fromIdx+1, len(tokens), tokens[fromIdx].depth, ctes)...)
	}
	return scopeWithDMLTarget(local, target)
}

func collectMergeScope(tokens []completionScopeToken, mergeIdx int, _ map[string]RangeReference) *ScopeSnapshot {
	idx := mergeIdx + 1
	if idx < len(tokens) && tokens[idx].Type == kwINTO {
		idx++
	}
	target, next, ok := parseRangeReference(tokens, idx)
	if !ok {
		return nil
	}
	target.Kind = RangeReferenceMergeTarget
	target.Alias, target.AliasColumns, next = parseAliasAndColumns(tokens, next, tokens[idx].depth)
	usingIdx := findKeywordAtDepth(tokens, kwUSING, next, len(tokens), tokens[mergeIdx].depth)
	local := []RangeReference{target}
	var source *RangeReference
	if usingIdx >= 0 {
		src, srcNext, ok := parseRangeReference(tokens, usingIdx+1)
		if ok {
			src.Kind = RangeReferenceMergeSource
			src.Alias, src.AliasColumns, srcNext = parseAliasAndColumns(tokens, srcNext, tokens[usingIdx].depth)
			src.Loc.End = tokens[srcNext-1].End
			source = &src
			local = append(local, src)
		}
	}
	scope := &ScopeSnapshot{
		References:      append([]RangeReference{}, local...),
		LocalReferences: local,
		MergeTarget:     &target,
		MergeSource:     source,
	}
	return scope
}

func collectSelectScope(tokens []completionScopeToken, cursorOffset int, ctes map[string]RangeReference) *ScopeSnapshot {
	frames := selectFramesForCursor(tokens, cursorOffset, ctes)
	if len(frames) == 0 {
		return &ScopeSnapshot{}
	}
	innermost := frames[len(frames)-1]
	var outerRefs [][]RangeReference
	for i := len(frames) - 2; i >= 0; i-- {
		if len(frames[i].localRefs) > 0 {
			outerRefs = append(outerRefs, frames[i].localRefs)
		}
	}
	refs := append([]RangeReference{}, innermost.localRefs...)
	for _, level := range outerRefs {
		refs = append(refs, level...)
	}
	return &ScopeSnapshot{
		References:      refs,
		LocalReferences: innermost.localRefs,
		OuterReferences: outerRefs,
	}
}

type completionSelectFrame struct {
	selectIdx int
	localRefs []RangeReference
}

func selectFramesForCursor(tokens []completionScopeToken, cursorOffset int, ctes map[string]RangeReference) []completionSelectFrame {
	cursorDepth := 0
	for _, tok := range tokens {
		if tok.Loc >= cursorOffset {
			break
		}
		switch tok.Type {
		case '(':
			cursorDepth++
		case ')':
			if cursorDepth > 0 {
				cursorDepth--
			}
		}
	}
	var frames []completionSelectFrame
	for i, tok := range tokens {
		if tok.Loc > cursorOffset {
			break
		}
		if tok.Type != kwSELECT || tok.depth > cursorDepth {
			continue
		}
		depth := tok.depth
		armEnd := selectArmEnd(tokens, i+1, depth)
		if armEnd < len(tokens) && tokens[armEnd].Loc < cursorOffset {
			continue
		}
		fromIdx := findKeywordAtDepth(tokens, kwFROM, i+1, armEnd, depth)
		var local []RangeReference
		if fromIdx >= 0 {
			local = parseFromReferences(tokens, fromIdx+1, armEnd, depth, ctes)
		}
		frames = append(frames, completionSelectFrame{selectIdx: i, localRefs: local})
	}
	return frames
}

func selectArmEnd(tokens []completionScopeToken, start, depth int) int {
	for i := start; i < len(tokens); i++ {
		if tokens[i].depth != depth {
			continue
		}
		switch tokens[i].Type {
		case kwUNION, kwINTERSECT, kwEXCEPT:
			return i
		}
	}
	return len(tokens)
}

func parseFromReferences(tokens []completionScopeToken, start, end, depth int, ctes map[string]RangeReference) []RangeReference {
	var refs []RangeReference
	expectSource := true
	for i := start; i < end; {
		tok := tokens[i]
		if tok.depth == depth && isFromClauseTerminator(tok.Type) {
			break
		}
		if tok.depth != depth {
			i++
			continue
		}
		if expectSource {
			ref, next, ok := parseTableSource(tokens, i, ctes)
			if ok {
				refs = append(refs, ref)
				i = next
				expectSource = false
				continue
			}
		}
		if tok.Type == ',' || tok.Type == kwJOIN || tok.Type == kwAPPLY {
			expectSource = true
		}
		i++
	}
	return refs
}

func parseTableSource(tokens []completionScopeToken, idx int, ctes map[string]RangeReference) (RangeReference, int, bool) {
	if idx >= len(tokens) {
		return RangeReference{}, idx, false
	}
	if tokens[idx].Type == '(' {
		return parseParenthesizedTableSource(tokens, idx)
	}
	ref, next, ok := parseRangeReference(tokens, idx)
	if !ok {
		return RangeReference{}, idx, false
	}
	if ref.Schema == "" {
		if cte, ok := ctes[foldName(ref.Object)]; ok {
			ref.Kind = RangeReferenceCTE
			ref.Columns = append([]string{}, cte.Columns...)
			ref.AliasColumns = append([]string{}, cte.Columns...)
		}
	}
	ref.Loc = tokenLoc(tokens[idx], tokens[next-1])
	ref.Alias, ref.AliasColumns, next = parseAliasAndColumns(tokens, next, tokens[idx].depth)
	if len(ref.AliasColumns) > 0 {
		ref.Columns = append([]string{}, ref.AliasColumns...)
	}
	return ref, next, true
}

func parseParenthesizedTableSource(tokens []completionScopeToken, idx int) (RangeReference, int, bool) {
	closeIdx := matchingParen(tokens, idx)
	if closeIdx < 0 {
		return RangeReference{}, idx, false
	}
	ref := RangeReference{
		Kind:    RangeReferenceSubquery,
		BodyLoc: tokenLoc(tokens[idx], tokens[closeIdx]),
		Loc:     tokenLoc(tokens[idx], tokens[closeIdx]),
	}
	if idx+1 < len(tokens) && tokens[idx+1].Type == kwVALUES {
		ref.Kind = RangeReferenceValues
	}
	alias, columns, next := parseAliasAndColumns(tokens, closeIdx+1, tokens[idx].depth)
	ref.Alias = alias
	ref.Object = alias
	ref.Columns = columns
	ref.AliasColumns = columns
	if next > closeIdx+1 {
		ref.Loc.End = tokens[next-1].End
	}
	return ref, next, alias != ""
}

func parseRangeReference(tokens []completionScopeToken, idx int) (RangeReference, int, bool) {
	if idx >= len(tokens) {
		return RangeReference{}, idx, false
	}
	if tokens[idx].Type == tokVARIABLE {
		name := tokens[idx].Str
		return RangeReference{
			Kind:   RangeReferenceTableVariable,
			Object: name,
			Alias:  name,
			Loc:    tokenLoc(tokens[idx], tokens[idx]),
		}, idx + 1, true
	}
	if !IsIdentTokenType(tokens[idx].Type) {
		return RangeReference{}, idx, false
	}

	var parts []string
	i := idx
	for i < len(tokens) && IsIdentTokenType(tokens[i].Type) {
		parts = append(parts, tokens[i].Str)
		if i+2 >= len(tokens) || tokens[i+1].Type != '.' || !IsIdentTokenType(tokens[i+2].Type) {
			i++
			break
		}
		i += 2
	}
	ref := RangeReference{Kind: RangeReferenceRelation, Loc: tokenLoc(tokens[idx], tokens[i-1])}
	switch len(parts) {
	case 1:
		ref.Object = parts[0]
	case 2:
		ref.Schema = parts[0]
		ref.Object = parts[1]
	case 3:
		ref.Database = parts[0]
		ref.Schema = parts[1]
		ref.Object = parts[2]
	default:
		ref.Server = parts[len(parts)-4]
		ref.Database = parts[len(parts)-3]
		ref.Schema = parts[len(parts)-2]
		ref.Object = parts[len(parts)-1]
	}
	return ref, i, true
}

func parseAliasAndColumns(tokens []completionScopeToken, idx, depth int) (string, []string, int) {
	if idx < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == kwAS {
		idx++
	}
	if idx >= len(tokens) || tokens[idx].depth != depth || !IsIdentTokenType(tokens[idx].Type) {
		return "", nil, idx
	}
	if isAliasBoundary(tokens[idx].Type) {
		return "", nil, idx
	}
	alias := tokens[idx].Str
	idx++
	var columns []string
	if idx < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == '(' {
		closeIdx := matchingParen(tokens, idx)
		if closeIdx >= 0 {
			for i := idx + 1; i < closeIdx; i++ {
				if IsIdentTokenType(tokens[i].Type) {
					columns = append(columns, tokens[i].Str)
				}
			}
			idx = closeIdx + 1
		}
	}
	return alias, columns, idx
}

func scopeWithDMLTarget(local []RangeReference, target RangeReference) *ScopeSnapshot {
	return &ScopeSnapshot{
		References:      append([]RangeReference{}, local...),
		LocalReferences: local,
		DMLTarget:       &target,
	}
}

func skipTopClause(tokens []completionScopeToken, idx, depth int) int {
	if idx >= len(tokens) || tokens[idx].Type != kwTOP {
		return idx
	}
	idx++
	if idx < len(tokens) && tokens[idx].Type == '(' {
		if closeIdx := matchingParen(tokens, idx); closeIdx >= 0 {
			return closeIdx + 1
		}
	}
	if idx < len(tokens) && tokens[idx].depth == depth {
		return idx + 1
	}
	return idx
}

func findKeywordAtDepth(tokens []completionScopeToken, typ, start, end, depth int) int {
	for i := start; i < end && i < len(tokens); i++ {
		if tokens[i].depth == depth && tokens[i].Type == typ {
			return i
		}
	}
	return -1
}

func firstNonTerminator(tokens []completionScopeToken) int {
	for i, tok := range tokens {
		if tok.Type != ';' && tok.Type != kwGO {
			return i
		}
	}
	return -1
}

func matchingParen(tokens []completionScopeToken, openIdx int) int {
	if openIdx >= len(tokens) || tokens[openIdx].Type != '(' {
		return -1
	}
	balance := 0
	for i := openIdx; i < len(tokens); i++ {
		switch tokens[i].Type {
		case '(':
			balance++
		case ')':
			balance--
			if balance == 0 {
				return i
			}
		}
	}
	return -1
}

func tokenLoc(first, last completionScopeToken) nodes.Loc {
	return nodes.Loc{Start: first.Loc, End: last.End}
}

func isFromClauseTerminator(typ int) bool {
	switch typ {
	case kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION, kwINTERSECT, kwEXCEPT, kwOPTION, kwFOR:
		return true
	default:
		return false
	}
}

func isAliasBoundary(typ int) bool {
	switch typ {
	case kwON, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwUNION, kwINTERSECT, kwEXCEPT, kwJOIN, kwAPPLY, kwWITH, kwOPTION, kwFOR:
		return true
	default:
		return false
	}
}

func completionPrefix(sql string, cursorOffset int) string {
	i := cursorOffset
	for i > 0 {
		c := sql[i-1]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			i--
			continue
		}
		break
	}
	return sql[i:cursorOffset]
}

func collectCompletionIntent(candidates *CandidateSet, tokens []Token, collectOffset int) *CompletionIntent {
	if candidates == nil {
		return nil
	}
	kinds := objectKindsForCandidates(candidates)
	if contextKind, ok := contextObjectKind(tokens, collectOffset); ok {
		if shouldOverrideObjectKindsWithContext(kinds, contextKind) {
			kinds = appendObjectKind(nil, contextKind)
		}
	}
	if len(kinds) == 0 {
		return nil
	}
	return &CompletionIntent{
		ObjectKinds: kinds,
		Qualifier:   qualifierBeforeOffset(tokens, collectOffset),
	}
}

func shouldOverrideObjectKindsWithContext(kinds []ObjectKind, contextKind ObjectKind) bool {
	if contextKind == ObjectKindUnknown {
		return false
	}
	if len(kinds) == 0 {
		return true
	}
	return len(kinds) == 1 && kinds[0] == ObjectKindTable && contextKind != ObjectKindTable
}

func objectKindsForCandidates(candidates *CandidateSet) []ObjectKind {
	var kinds []ObjectKind
	for _, rc := range candidates.Rules {
		switch rc.Rule {
		case "database_ref":
			kinds = appendObjectKind(kinds, ObjectKindDatabase)
		case "schema_ref":
			kinds = appendObjectKind(kinds, ObjectKindSchema)
		case "table_ref":
			kinds = appendObjectKind(kinds, ObjectKindTable)
		case "view_name", "view_ref":
			kinds = appendObjectKind(kinds, ObjectKindView)
		case "sequence_ref":
			kinds = appendObjectKind(kinds, ObjectKindSequence)
		case "proc_name", "proc_ref":
			kinds = appendObjectKind(kinds, ObjectKindProcedure)
		case "func_name":
			kinds = appendObjectKind(kinds, ObjectKindFunction)
		case "type_name":
			kinds = appendObjectKind(kinds, ObjectKindType)
		case "columnref", "cte_column_name":
			kinds = appendObjectKind(kinds, ObjectKindColumn)
		}
	}
	return kinds
}

func appendObjectKind(kinds []ObjectKind, kind ObjectKind) []ObjectKind {
	for _, existing := range kinds {
		if existing == kind {
			return kinds
		}
	}
	return append(kinds, kind)
}

func contextObjectKind(tokens []Token, collectOffset int) (ObjectKind, bool) {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Loc >= collectOffset {
			continue
		}
		switch tokens[i].Type {
		case kwTABLE:
			return ObjectKindTable, true
		case kwVIEW:
			return ObjectKindView, true
		case kwSEQUENCE:
			return ObjectKindSequence, true
		case kwPROCEDURE, kwPROC:
			return ObjectKindProcedure, true
		case kwFUNCTION:
			return ObjectKindFunction, true
		case kwTYPE:
			return ObjectKindType, true
		case kwDATABASE:
			return ObjectKindDatabase, true
		case kwSCHEMA:
			return ObjectKindSchema, true
		case kwFROM, kwJOIN, kwINTO, kwUSING, kwFOR, kwREFERENCES, kwEXEC, kwEXECUTE, kwDROP, kwCREATE, kwALTER:
			return ObjectKindUnknown, false
		}
	}
	return ObjectKindUnknown, false
}

func qualifierBeforeOffset(tokens []Token, collectOffset int) MultipartName {
	idx := -1
	for i, tok := range tokens {
		if tok.Loc >= collectOffset {
			break
		}
		idx = i
	}
	var parts []string
	if idx >= 0 && tokens[idx].Type == '.' {
		idx--
		for idx >= 0 && IsIdentTokenType(tokens[idx].Type) {
			parts = append([]string{tokens[idx].Str}, parts...)
			if idx-1 < 0 || tokens[idx-1].Type != '.' {
				break
			}
			idx -= 2
		}
	}
	switch len(parts) {
	case 1:
		return MultipartName{Schema: parts[0]}
	case 2:
		return MultipartName{Database: parts[0], Schema: parts[1]}
	case 3:
		return MultipartName{Server: parts[0], Database: parts[1], Schema: parts[2]}
	default:
		if len(parts) > 3 {
			return MultipartName{Server: parts[len(parts)-3], Database: parts[len(parts)-2], Schema: parts[len(parts)-1]}
		}
		return MultipartName{}
	}
}

func foldName(name string) string {
	if name == "" {
		return ""
	}
	buf := make([]byte, len(name))
	for i := range name {
		c := name[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		buf[i] = c
	}
	return string(buf)
}
