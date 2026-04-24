package parser

import (
	"sort"

	nodes "github.com/bytebase/omni/pg/ast"
)

// CompletionContext is the parser-native context for SQL completion.
//
// Candidates contains grammar candidates at the cursor. Scope is a
// best-effort snapshot of range references visible to the SELECT expression
// scope at the cursor. This is completion-only state; Parse remains strict and
// does not expose partial ASTs.
type CompletionContext struct {
	Candidates *CandidateSet
	Scope      *ScopeSnapshot

	// CTEs contains CTE definitions visible to the cursor statement. These are
	// definitions, not necessarily range references already present in FROM.
	CTEs []RangeReference
}

// ScopeSnapshot describes the SELECT-level relation scope visible at a cursor.
type ScopeSnapshot struct {
	// References is the full visible scope at the cursor: local references first,
	// then outer reference levels from nearest to farthest.
	References []RangeReference

	// LocalReferences is the cursor SELECT's own FROM/JOIN scope.
	LocalReferences []RangeReference

	// OuterReferences contains legal outer scopes, nearest first.
	OuterReferences [][]RangeReference

	Depth int
}

// RangeReferenceKind classifies a range-table entry visible to completion.
type RangeReferenceKind int

const (
	RangeReferenceRelation RangeReferenceKind = iota
	RangeReferenceSubquery
	RangeReferenceFunction
	RangeReferenceJoinAlias
	RangeReferenceCTE
)

// RangeReference is a syntax-level FROM/JOIN reference. It deliberately avoids
// catalog metadata; callers resolve relation kinds and columns themselves.
type RangeReference struct {
	Kind RangeReferenceKind

	Catalog string
	Schema  string
	Name    string
	Alias   string

	AliasColumns []string

	Loc     nodes.Loc
	BodyLoc nodes.Loc

	Lateral bool
}

// CollectCompletion runs completion collection and returns parser candidates
// plus a best-effort visible relation scope.
func CollectCompletion(sql string, cursorOffset int) *CompletionContext {
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	scope, ctes := collectCompletionScopeAndCTEs(sql, cursorOffset)
	return &CompletionContext{
		Candidates: Collect(sql, cursorOffset),
		Scope:      scope,
		CTEs:       ctes,
	}
}

type completionScopeToken struct {
	Token
	depth int
}

type completionScopeParser struct {
	tokens []completionScopeToken
	ctes   map[string]RangeReference
}

type completionSelectFrame struct {
	selectIdx int
	depth     int
	localRefs []RangeReference
}

func collectCompletionScope(sql string, cursorOffset int) *ScopeSnapshot {
	scope, _ := collectCompletionScopeAndCTEs(sql, cursorOffset)
	return scope
}

func collectCompletionScopeAndCTEs(sql string, cursorOffset int) (*ScopeSnapshot, []RangeReference) {
	p := &completionScopeParser{
		tokens: completionScopeTokens(Tokenize(sql)),
	}
	selects, cursorDepth := p.selectStackForCursor(cursorOffset)
	if len(selects) == 0 {
		return &ScopeSnapshot{Depth: cursorDepth}, nil
	}

	visibleCTEs := make(map[string]RangeReference)
	var frames []completionSelectFrame
	for _, sel := range selects {
		for name, ref := range p.collectCTEsForSelect(sel.selectIdx, sel.depth) {
			visibleCTEs[name] = ref
		}
		p.ctes = visibleCTEs
		frame := sel
		if fromIdx := p.findFromForSelect(sel.selectIdx, sel.depth); fromIdx >= 0 {
			frame.localRefs, _ = p.parseFromRefs(fromIdx+1, sel.depth, len(p.tokens))
		}
		frames = append(frames, frame)
	}

	innermost := frames[len(frames)-1]
	var outerRefs [][]RangeReference
	for i := len(frames) - 2; i >= 0; i-- {
		refs := p.refsVisibleFromOuterFrame(frames[i], cursorOffset)
		if len(refs) > 0 {
			outerRefs = append(outerRefs, refs)
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
		Depth:           innermost.depth,
	}, sortedRangeReferences(visibleCTEs)
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

func (p *completionScopeParser) selectStackForCursor(cursorOffset int) ([]completionSelectFrame, int) {
	cursorDepth := 0
	stmtStart := 0
	for i, tok := range p.tokens {
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
		case ';':
			stmtStart = i + 1
		}
	}

	bestByDepth := make(map[int]int)
	for i := stmtStart; i < len(p.tokens); i++ {
		tok := p.tokens[i]
		if tok.Loc > cursorOffset {
			break
		}
		if tok.Type == SELECT && tok.depth <= cursorDepth {
			bestByDepth[tok.depth] = i
		}
	}
	if len(bestByDepth) == 0 {
		return nil, cursorDepth
	}
	depths := make([]int, 0, len(bestByDepth))
	for depth := range bestByDepth {
		depths = append(depths, depth)
	}
	sort.Ints(depths)

	frames := make([]completionSelectFrame, 0, len(depths))
	for _, depth := range depths {
		frames = append(frames, completionSelectFrame{
			selectIdx: bestByDepth[depth],
			depth:     depth,
		})
	}
	return frames, cursorDepth
}

func (p *completionScopeParser) collectCTEsForSelect(selectIdx int, depth int) map[string]RangeReference {
	result := make(map[string]RangeReference)
	stmtStart := 0
	for i := selectIdx - 1; i >= 0; i-- {
		if p.tokens[i].depth == depth && p.tokens[i].Type == ';' {
			stmtStart = i + 1
			break
		}
	}
	stmtEnd := len(p.tokens)
	for i := selectIdx + 1; i < len(p.tokens); i++ {
		if p.tokens[i].Type == ';' {
			stmtEnd = i
			break
		}
	}
	for i := stmtStart; i < selectIdx; i++ {
		if p.tokens[i].depth <= depth && (p.tokens[i].Type == WITH || p.tokens[i].Type == WITH_LA) {
			p.parseCTEs(result, i+1, stmtEnd, p.tokens[i].depth)
		}
	}
	return result
}

func (p *completionScopeParser) parseCTEs(result map[string]RangeReference, i int, end int, depth int) {
	if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == RECURSIVE {
		i++
	}
	for i < end {
		if p.tokens[i].depth != depth || !isCompletionColID(p.tokens[i].Type, p.tokens[i].Str) {
			return
		}
		nameTok := p.tokens[i]
		ref := RangeReference{
			Kind: RangeReferenceCTE,
			Name: nameTok.Str,
			Loc:  nodes.Loc{Start: nameTok.Loc, End: nameTok.End},
		}
		i++
		if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == '(' {
			cols, next, loc := p.parseNameListAt(i)
			ref.AliasColumns = cols
			ref.Loc.End = loc.End
			i = next
		}
		if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == AS {
			i++
		}
		if i < end && p.tokens[i].depth == depth && (p.tokens[i].Type == MATERIALIZED || p.tokens[i].Type == NOT) {
			if p.tokens[i].Type == NOT {
				i++
			}
			if i < end && p.tokens[i].Type == MATERIALIZED {
				i++
			}
		}
		if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == '(' {
			closeIdx := p.matchingParen(i)
			if closeIdx < 0 || closeIdx >= end {
				return
			}
			if i+1 < closeIdx {
				ref.BodyLoc = nodes.Loc{Start: p.tokens[i+1].Loc, End: p.tokens[closeIdx].Loc}
			}
			ref.Loc.End = p.tokens[closeIdx].End
			i = closeIdx + 1
		}
		result[ref.Name] = ref
		if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == ',' {
			i++
			continue
		}
		return
	}
}

func (p *completionScopeParser) refsVisibleFromOuterFrame(frame completionSelectFrame, cursorOffset int) []RangeReference {
	for _, ref := range frame.localRefs {
		if ref.BodyLoc.Start >= 0 && ref.BodyLoc.Start <= cursorOffset && cursorOffset <= ref.BodyLoc.End {
			if !ref.Lateral {
				return nil
			}
			var refs []RangeReference
			for _, candidate := range frame.localRefs {
				if candidate.Loc.Start >= 0 && candidate.Loc.Start < ref.Loc.Start {
					refs = append(refs, candidate)
				}
			}
			return refs
		}
	}
	return frame.localRefs
}

func (p *completionScopeParser) findFromForSelect(selectIdx int, depth int) int {
	for i := selectIdx + 1; i < len(p.tokens); i++ {
		tok := p.tokens[i]
		if tok.depth < depth {
			return -1
		}
		if tok.depth != depth {
			continue
		}
		if tok.Type == FROM {
			return i
		}
		if tok.Type == ';' || isSelectBoundary(tok.Type) {
			return -1
		}
	}
	return -1
}

func (p *completionScopeParser) parseFromRefs(i int, depth int, end int) ([]RangeReference, int) {
	var refs []RangeReference
	for i < end {
		tok := p.tokens[i]
		if tok.depth < depth || (tok.depth == depth && (tok.Type == ';' || isFromTerminator(tok.Type))) {
			return refs, i
		}
		if tok.depth != depth {
			i++
			continue
		}
		switch {
		case tok.Type == ',':
			i++
		case tok.Type == ON:
			i = p.skipJoinON(i+1, depth, end)
		case tok.Type == USING:
			i = p.skipJoinUsing(i+1, depth, end)
		case isJoinStart(tok.Type):
			i = p.consumeJoinStart(i, depth, end)
		default:
			itemRefs, next, ok := p.parseFromPrimary(i, depth, end)
			if !ok {
				i++
				continue
			}
			refs = append(refs, itemRefs...)
			i = next
		}
	}
	return refs, i
}

func (p *completionScopeParser) parseFromPrimary(i int, depth int, end int) ([]RangeReference, int, bool) {
	if i >= end {
		return nil, i, false
	}
	lateral := false
	if p.tokens[i].depth == depth && p.tokens[i].Type == LATERAL_P {
		lateral = true
		i++
	}
	if i >= end || p.tokens[i].depth != depth {
		return nil, i, false
	}
	tok := p.tokens[i]
	if tok.Type == '(' {
		return p.parseParenFromPrimary(i, depth, end, lateral)
	}
	if !isCompletionColID(tok.Type, tok.Str) {
		return nil, i, false
	}
	return p.parseNamedFromPrimary(i, depth, end, lateral)
}

func (p *completionScopeParser) parseParenFromPrimary(i int, depth int, end int, lateral bool) ([]RangeReference, int, bool) {
	closeIdx := p.matchingParen(i)
	if closeIdx < 0 || closeIdx >= end {
		return nil, i + 1, false
	}
	if i+1 < closeIdx && isSelectLikeStart(p.tokens[i+1].Type) {
		alias, cols, next, aliasLoc := p.parseAlias(closeIdx+1, depth, end)
		loc := nodes.Loc{Start: p.tokens[i].Loc, End: p.tokens[closeIdx].End}
		if aliasLoc.Start >= 0 {
			loc.End = aliasLoc.End
		}
		ref := RangeReference{
			Kind:         RangeReferenceSubquery,
			Alias:        alias,
			AliasColumns: cols,
			Loc:          loc,
			BodyLoc:      nodes.Loc{Start: p.tokens[i+1].Loc, End: p.tokens[closeIdx].Loc},
			Lateral:      lateral,
		}
		return []RangeReference{ref}, next, true
	}

	inner, _ := p.parseFromRefs(i+1, depth+1, closeIdx)
	alias, cols, next, aliasLoc := p.parseAlias(closeIdx+1, depth, end)
	if alias != "" {
		loc := nodes.Loc{Start: p.tokens[i].Loc, End: p.tokens[closeIdx].End}
		if aliasLoc.Start >= 0 {
			loc.End = aliasLoc.End
		}
		return []RangeReference{{
			Kind:         RangeReferenceJoinAlias,
			Alias:        alias,
			AliasColumns: cols,
			Loc:          loc,
			Lateral:      lateral,
		}}, next, true
	}
	return inner, next, true
}

func (p *completionScopeParser) parseNamedFromPrimary(i int, depth int, end int, lateral bool) ([]RangeReference, int, bool) {
	start := p.tokens[i]
	var parts []string
	for i < end && p.tokens[i].depth == depth && isCompletionColID(p.tokens[i].Type, p.tokens[i].Str) {
		parts = append(parts, p.tokens[i].Str)
		i++
		if i+1 >= end || p.tokens[i].depth != depth || p.tokens[i].Type != '.' || p.tokens[i+1].depth != depth {
			break
		}
		i++
	}
	if len(parts) == 0 {
		return nil, i, false
	}
	isFunction := i < end && p.tokens[i].depth == depth && p.tokens[i].Type == '('
	if isFunction {
		closeIdx := p.matchingParen(i)
		if closeIdx < 0 || closeIdx >= end {
			return nil, i, false
		}
		alias, cols, next, aliasLoc := p.parseAlias(closeIdx+1, depth, end)
		loc := nodes.Loc{Start: start.Loc, End: p.tokens[closeIdx].End}
		if aliasLoc.Start >= 0 {
			loc.End = aliasLoc.End
		}
		ref := RangeReference{
			Kind:         RangeReferenceFunction,
			Name:         lastPart(parts),
			Alias:        alias,
			AliasColumns: cols,
			Loc:          loc,
			BodyLoc:      nodes.Loc{Start: start.Loc, End: p.tokens[closeIdx].End},
			Lateral:      lateral,
		}
		if len(parts) == 2 {
			ref.Schema = parts[0]
		} else if len(parts) >= 3 {
			ref.Catalog = parts[len(parts)-3]
			ref.Schema = parts[len(parts)-2]
		}
		return []RangeReference{ref}, next, true
	}

	alias, cols, next, aliasLoc := p.parseAlias(i, depth, end)
	ref := RangeReference{
		Kind:         RangeReferenceRelation,
		Name:         lastPart(parts),
		Alias:        alias,
		AliasColumns: cols,
		Loc:          nodes.Loc{Start: start.Loc, End: p.tokens[i-1].End},
		Lateral:      lateral,
	}
	if aliasLoc.Start >= 0 {
		ref.Loc.End = aliasLoc.End
	}
	if len(parts) == 1 {
		if cte, ok := p.ctes[parts[0]]; ok {
			ref.Kind = RangeReferenceCTE
			ref.AliasColumns = cte.AliasColumns
			ref.BodyLoc = cte.BodyLoc
			if len(cols) > 0 {
				ref.AliasColumns = cols
			}
		}
	} else if len(parts) == 2 {
		ref.Schema = parts[0]
	} else {
		ref.Catalog = parts[len(parts)-3]
		ref.Schema = parts[len(parts)-2]
	}
	return []RangeReference{ref}, next, true
}

func (p *completionScopeParser) parseAlias(i int, depth int, end int) (string, []string, int, nodes.Loc) {
	loc := nodes.NoLoc()
	if i >= end || p.tokens[i].depth != depth {
		return "", nil, i, loc
	}
	if p.tokens[i].Type == AS {
		if i+1 >= end || p.tokens[i+1].depth != depth || !isCompletionColID(p.tokens[i+1].Type, p.tokens[i+1].Str) {
			return "", nil, i, loc
		}
		start := p.tokens[i]
		name := p.tokens[i+1].Str
		loc = nodes.Loc{Start: start.Loc, End: p.tokens[i+1].End}
		i += 2
		if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == '(' {
			cols, next, colLoc := p.parseNameListAt(i)
			loc.End = colLoc.End
			return name, cols, next, loc
		}
		return name, nil, i, loc
	}
	if !isAliasToken(p.tokens[i]) {
		return "", nil, i, loc
	}
	name := p.tokens[i].Str
	loc = nodes.Loc{Start: p.tokens[i].Loc, End: p.tokens[i].End}
	i++
	if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == '(' {
		cols, next, colLoc := p.parseNameListAt(i)
		loc.End = colLoc.End
		return name, cols, next, loc
	}
	return name, nil, i, loc
}

func (p *completionScopeParser) parseNameListAt(i int) ([]string, int, nodes.Loc) {
	loc := nodes.Loc{Start: p.tokens[i].Loc, End: p.tokens[i].End}
	closeIdx := p.matchingParen(i)
	if closeIdx < 0 {
		return nil, i + 1, loc
	}
	var cols []string
	for j := i + 1; j < closeIdx; j++ {
		if isCompletionColID(p.tokens[j].Type, p.tokens[j].Str) {
			cols = append(cols, p.tokens[j].Str)
		}
	}
	loc.End = p.tokens[closeIdx].End
	return cols, closeIdx + 1, loc
}

func (p *completionScopeParser) matchingParen(openIdx int) int {
	if openIdx < 0 || openIdx >= len(p.tokens) || p.tokens[openIdx].Type != '(' {
		return -1
	}
	depth := 0
	for i := openIdx; i < len(p.tokens); i++ {
		switch p.tokens[i].Type {
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

func (p *completionScopeParser) consumeJoinStart(i int, depth int, end int) int {
	if i >= end || p.tokens[i].depth != depth {
		return i
	}
	switch p.tokens[i].Type {
	case JOIN:
		return i + 1
	case CROSS, INNER_P:
		if i+1 < end && p.tokens[i+1].depth == depth && p.tokens[i+1].Type == JOIN {
			return i + 2
		}
	case LEFT, RIGHT, FULL:
		i++
		if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == OUTER_P {
			i++
		}
		if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == JOIN {
			return i + 1
		}
	case NATURAL:
		i++
		if i < end && p.tokens[i].depth == depth {
			switch p.tokens[i].Type {
			case LEFT, RIGHT, FULL:
				i++
				if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == OUTER_P {
					i++
				}
			case INNER_P:
				i++
			}
		}
		if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == JOIN {
			return i + 1
		}
	}
	return i + 1
}

func (p *completionScopeParser) skipJoinON(i int, depth int, end int) int {
	for i < end {
		tok := p.tokens[i]
		if tok.depth == depth && (tok.Type == ',' || tok.Type == ';' || isFromTerminator(tok.Type) || isJoinStart(tok.Type)) {
			return i
		}
		i++
	}
	return i
}

func (p *completionScopeParser) skipJoinUsing(i int, depth int, end int) int {
	if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == '(' {
		closeIdx := p.matchingParen(i)
		if closeIdx >= 0 && closeIdx < end {
			i = closeIdx + 1
		}
	}
	if i < end && p.tokens[i].depth == depth && p.tokens[i].Type == AS {
		if i+1 < end && p.tokens[i+1].depth == depth && isCompletionColID(p.tokens[i+1].Type, p.tokens[i+1].Str) {
			return i + 2
		}
	}
	return i
}

func isCompletionColID(tokenType int, tokenText string) bool {
	if tokenType == IDENT {
		return true
	}
	if kw := LookupKeyword(tokenText); kw != nil && kw.Token == tokenType {
		return kw.Category == UnreservedKeyword || kw.Category == ColNameKeyword
	}
	return false
}

func isAliasToken(tok completionScopeToken) bool {
	if !isCompletionColID(tok.Type, tok.Str) {
		return false
	}
	if isFromTerminator(tok.Type) || isJoinStart(tok.Type) {
		return false
	}
	switch tok.Type {
	case ON, USING, AS, TABLESAMPLE:
		return false
	default:
		return true
	}
}

func isJoinStart(tokenType int) bool {
	switch tokenType {
	case JOIN, CROSS, LEFT, RIGHT, FULL, INNER_P, NATURAL:
		return true
	default:
		return false
	}
}

func isFromTerminator(tokenType int) bool {
	switch tokenType {
	case WHERE, GROUP_P, HAVING, WINDOW, ORDER, LIMIT, OFFSET, FETCH, FOR,
		UNION, EXCEPT, INTERSECT:
		return true
	default:
		return false
	}
}

func isSelectBoundary(tokenType int) bool {
	switch tokenType {
	case WHERE, GROUP_P, HAVING, WINDOW, ORDER, LIMIT, OFFSET, FETCH, FOR,
		UNION, EXCEPT, INTERSECT:
		return true
	default:
		return false
	}
}

func isSelectLikeStart(tokenType int) bool {
	switch tokenType {
	case SELECT, VALUES, WITH, WITH_LA, TABLE:
		return true
	default:
		return false
	}
}

func lastPart(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func sortedRangeReferences(in map[string]RangeReference) []RangeReference {
	refs := make([]RangeReference, 0, len(in))
	for _, ref := range in {
		refs = append(refs, ref)
	}
	sort.Slice(refs, func(i, j int) bool {
		if refs[i].Loc.Start != refs[j].Loc.Start {
			return refs[i].Loc.Start < refs[j].Loc.Start
		}
		return refs[i].Name < refs[j].Name
	})
	return refs
}
