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

	StatementLoc nodes.Loc
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

	Unsupported       bool
	UnsupportedReason string
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
	tokens := multisetOperatorsAsWords(completionTokens(sql, cursorOffset))
	candidates, intent := collectOracleCandidates(sql, tokens, collectOffset, prefix)
	scope, ctes := collectOracleScopeAndCTEs(tokens, cursorOffset)
	stmtLoc := oracleCompletionStatementLoc(sql, tokens, cursorOffset)
	return &CompletionContext{
		Candidates:   candidates,
		Scope:        scope,
		CTEs:         ctes,
		Prefix:       prefix,
		Intent:       intent,
		StatementLoc: stmtLoc,
	}
}

func collectOracleScope(tokens []Token, cursorOffset int) *ScopeSnapshot {
	scope, _ := collectOracleScopeAndCTEs(tokens, cursorOffset)
	return scope
}

func collectOracleScopeAndCTEs(tokens []Token, cursorOffset int) (*ScopeSnapshot, []RangeReference) {
	start, end := oracleStatementTokenBounds(tokens, cursorOffset)
	stmtTokens := tokens[start:end]
	index := newOracleTokenIndex(stmtTokens)
	plain := index.reader(nil)
	ctes := plain.ctesOfQuery()
	cteMap := oracleRangeReferenceMap(ctes)
	r := index.reader(cteMap)
	selectRefs := r.selectRefs()
	dmlRefs := plain.dmlRefs()
	ddlRefs := plain.ddlRefs()
	refs := append([]RangeReference{}, selectRefs...)
	refs = append(refs, dmlRefs...)
	refs = append(refs, ddlRefs...)
	refs = append(append([]RangeReference{}, ctes...), refs...)
	localRefs := refs
	outerRefs := [][]RangeReference(nil)
	if selects := index.enclosingSelects(cursorOffset); len(selects) > 0 {
		inner := selects[len(selects)-1]
		armStart, armEnd := r.armBounds(inner, len(stmtTokens))
		localRefs = r.selectRefsInRange(armStart, armEnd, index.depth[inner], len(stmtTokens))
		outerRefs = r.outerSelectRefs(cursorOffset, selects[:len(selects)-1])
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

// oracleTokenIndex is what reading the range references of a statement
// looks up about its tokens, found in one pass. The readers below scan only
// the tokens a clause holds at its own depth, step over each parenthesized
// group in one move, and read the select list of a derived table once
// (oracleRefReader.columns). Reading the references of each SELECT from the
// start of the statement, finding each matching parenthesis by a scan, and
// expanding a star over each derived table again for every query around
// it made completion through 1,000 nested derived tables (18KB) take 35s,
// the time growing with the cube of the depth.
type oracleTokenIndex struct {
	toks []Token
	// depth[i] is the number of parentheses open before toks[i]; a ')'
	// with none open closes none.
	depth []int
	// match[i] is, for a parenthesis, the index of the one that matches
	// it, or -1 when none does.
	match []int
	// runs[i] is, for a matched '(', how deeply runs of '(' nest between it
	// and its ')', as oracleParenRunsNestAtLeast counts them.
	runs []int
	// setOp[i] is, for a SELECT, the index of the last set operator before
	// it inside the same parentheses, or -1. One at the same depth in an
	// earlier group is not it: in (SELECT a FROM t UNION SELECT b FROM u) x,
	// (SELECT c FROM v) y, the arm of SELECT c starts at SELECT c.
	setOp []int
}

func newOracleTokenIndex(toks []Token) *oracleTokenIndex {
	x := &oracleTokenIndex{
		toks:  toks,
		depth: make([]int, len(toks)),
		match: make([]int, len(toks)),
		runs:  make([]int, len(toks)),
		setOp: make([]int, len(toks)),
	}
	// open holds the '(' not yet closed, the depth counting them. For each
	// it keeps how deeply runs nest in what it has closed so far: nest, as
	// runs counts them, and inRun, the same but a run that the '(' itself
	// continues, one that opens right after it, not counted.
	type openParen struct{ idx, nest, inRun int }
	var open []openParen
	// lastSetOp holds, by depth, the last set operator inside each '(' open,
	// so a '(' forgets those of the groups closed before it at its depth.
	var lastSetOp []int
	for i, tok := range toks {
		x.depth[i], x.match[i], x.setOp[i] = len(open), -1, -1
		switch tok.Type {
		case '(':
			open = append(open, openParen{idx: i})
			lastSetOp = lastSetOp[:min(len(lastSetOp), len(open))]
		case ')':
			n := len(open)
			if n == 0 {
				continue
			}
			o := open[n-1]
			open = open[:n-1]
			x.match[o.idx], x.match[i] = i, o.idx
			x.runs[o.idx] = o.nest
			// How deeply runs nest from o on, o opening one.
			nest := 1 + o.inRun
			if n > 1 {
				p := &open[n-2]
				p.nest = max(p.nest, nest)
				if o.idx == p.idx+1 {
					nest--
				}
				p.inRun = max(p.inRun, nest)
			}
		case kwSELECT:
			if d := len(open); d < len(lastSetOp) {
				x.setOp[i] = lastSetOp[d]
			}
		case kwUNION, kwINTERSECT, kwMINUS, kwEXCEPT:
			for len(lastSetOp) <= len(open) {
				lastSetOp = append(lastSetOp, -1)
			}
			lastSetOp[len(open)] = i
		}
	}
	return x
}

// closeOf returns the index of the ')' that matches the '(' at open before
// hi, or -1: oracleMatchingParen(toks[:hi], open).
func (x *oracleTokenIndex) closeOf(open, hi int) int {
	if open < 0 || open >= hi || x.toks[open].Type != '(' {
		return -1
	}
	if c := x.match[open]; c >= 0 && c < hi {
		return c
	}
	return -1
}

// oracleRefReader reads range references from the tokens of an index, with
// the CTEs they may name.
type oracleRefReader struct {
	*oracleTokenIndex
	ctes map[string]RangeReference
	// columns holds, by the index of its '(', the select-list columns of
	// each derived table read so far: expanding a star reads the derived
	// tables it selects from, and each query around them reads it again.
	columns map[int][]string
}

func (x *oracleTokenIndex) reader(ctes map[string]RangeReference) *oracleRefReader {
	return &oracleRefReader{oracleTokenIndex: x, ctes: ctes}
}

func (r *oracleRefReader) ddlRefs() []RangeReference {
	tokens := r.toks
	first := firstOracleNonTerminator(tokens)
	if first < 0 {
		return nil
	}
	switch tokens[first].Type {
	case kwCREATE:
		var refs []RangeReference
		if first+2 < len(tokens) && tokens[first+1].Type == kwTABLE {
			_, ref, ok := r.rangeReference(first+2, len(tokens))
			if ok {
				refs = append(refs, ref)
			}
		}
		for i := first + 1; i < len(tokens); i++ {
			if tokens[i].Type == kwON {
				_, ref, ok := r.rangeReference(i+1, len(tokens))
				if ok {
					refs = append(refs, ref)
				}
			}
			if tokens[i].Type == kwREFERENCES {
				_, ref, ok := r.rangeReference(i+1, len(tokens))
				if ok {
					refs = append(refs, ref)
				}
			}
		}
		return refs
	case kwALTER:
		for i := first + 1; i < len(tokens); i++ {
			if tokens[i].Type == kwTABLE {
				_, ref, ok := r.rangeReference(i+1, len(tokens))
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

func (r *oracleRefReader) dmlRefs() []RangeReference {
	tokens := r.toks
	first := firstOracleNonTerminator(tokens)
	if first < 0 {
		return nil
	}
	switch tokens[first].Type {
	case kwUPDATE:
		_, ref, ok := r.rangeReference(first+1, len(tokens))
		if ok {
			ref.Kind = RangeReferenceDMLTarget
			return []RangeReference{ref}
		}
	case kwINSERT:
		for i := first + 1; i < len(tokens); i++ {
			if tokens[i].Type == kwINTO {
				_, ref, ok := r.rangeReference(i+1, len(tokens))
				if ok {
					ref.Kind = RangeReferenceDMLTarget
					return []RangeReference{ref}
				}
				break
			}
		}
	case kwDELETE:
		// The target follows DELETE or DELETE FROM, after any hint. A FROM
		// further on is a subquery's: DELETE t WHERE a IN (SELECT b FROM u).
		i := first + 1
		if i < len(tokens) && tokens[i].Type == tokHINT {
			i++
		}
		if i < len(tokens) && tokens[i].Type == kwFROM {
			i++
		}
		_, ref, ok := r.rangeReference(i, len(tokens))
		if ok {
			ref.Kind = RangeReferenceDMLTarget
			return []RangeReference{ref}
		}
	case kwMERGE:
		var refs []RangeReference
		for i := first + 1; i < len(tokens); i++ {
			if r.depth[i] != 0 {
				// An INTO or USING in parentheses is not MERGE's: USING
				// (SELECT a FROM t JOIN u USING (id)) s.
				continue
			}
			switch tokens[i].Type {
			case kwINTO:
				next, ref, ok := r.rangeReference(i+1, len(tokens))
				if ok {
					ref.Kind = RangeReferenceDMLTarget
					refs = append(refs, ref)
					i = next - 1
				}
			case kwUSING:
				next, ref, ok := r.rangeReference(i+1, len(tokens))
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

func oracleCompletionStatementLoc(sql string, tokens []Token, cursorOffset int) nodes.Loc {
	start, end := oracleStatementTokenBounds(tokens, cursorOffset)
	loc := nodes.Loc{Start: cursorOffset, End: cursorOffset}
	if start < len(tokens) {
		loc.Start = tokens[start].Loc
	} else {
		for loc.Start > 0 && sql[loc.Start-1] != ';' {
			loc.Start--
		}
	}
	if end < len(tokens) {
		loc.End = tokens[end].Loc
	} else {
		loc.End = len(sql)
	}
	for loc.Start < loc.End {
		switch sql[loc.Start] {
		case ' ', '\t', '\n', '\r':
			loc.Start++
		default:
			return loc
		}
	}
	return loc
}

// selectRefs returns the references of every SELECT in the statement.
func (r *oracleRefReader) selectRefs() []RangeReference {
	var refs []RangeReference
	for i, tok := range r.toks {
		if tok.Type != kwSELECT {
			continue
		}
		start, end := r.armBounds(i, len(r.toks))
		refs = append(refs, r.selectRefsInRange(start, end, r.depth[i], len(r.toks))...)
	}
	return refs
}

// selectRefsInRange returns the references of the FROM clauses at depth
// want among toks[start:end], hi bounding what a reference may read. A group
// that opens at depth want or deeper holds no FROM clause at want and is
// stepped over whole.
//
// The depth is the index's. This scan once counted it itself and missed
// the ')' a FROM clause read up to, so when a range ran past the end of the
// group it started in, it took each later group for a deeper one.
func (r *oracleRefReader) selectRefsInRange(start, end, want, hi int) []RangeReference {
	var refs []RangeReference
	for i := start; i < end; i++ {
		depth := r.depth[i]
		if r.toks[i].Type == '(' && depth >= want {
			if c := r.match[i]; c >= 0 && c < end {
				i = c
				continue
			}
			return refs
		}
		if depth != want {
			continue
		}
		if r.toks[i].Type == kwFROM {
			next, fromRefs := r.fromReferences(i+1, end, want, hi)
			refs = append(refs, fromRefs...)
			i = next - 1
		}
	}
	return refs
}

func (r *oracleRefReader) fromReferences(start, end, want, hi int) (int, []RangeReference) {
	var refs []RangeReference
	tokens, depths := r.toks, r.depth
	i := start
	expectingRef := true
	for i < end {
		if depths[i] < want {
			break
		}
		if depths[i] != want {
			i++
			continue
		}
		if oracleIsFromClauseBoundary(tokens[i].Type) {
			break
		}
		switch tokens[i].Type {
		case ',':
			expectingRef = true
			i++
			continue
		case kwJOIN:
			expectingRef = true
			i++
			continue
		case kwLEFT, kwRIGHT, kwFULL, kwCROSS, kwINNER, kwOUTER:
			i++
			continue
		case kwON, kwUSING:
			i = r.skipJoinPredicate(i+1, end, want)
			expectingRef = false
			continue
		}
		if expectingRef {
			next, ref, ok := r.rangeReference(i, hi)
			if ok && next <= end {
				ref = resolveOracleCTEReference(ref, r.ctes)
				ref = r.markUnsupportedTransform(ref, next, end, want)
				refs = append(refs, ref)
				i = next
				expectingRef = false
				continue
			}
		}
		if tokens[i].Type == '(' {
			// What it holds is deeper: past its ')'.
			if c := r.match[i]; c >= 0 && c < end {
				i = c + 1
				continue
			}
			i = end
			continue
		}
		i++
	}
	return i, refs
}

func (r *oracleRefReader) skipJoinPredicate(start, end, want int) int {
	tokens := r.toks
	for i := start; i < end; i++ {
		if r.depth[i] != want {
			continue
		}
		if oracleIsFromClauseBoundary(tokens[i].Type) || tokens[i].Type == ',' || tokens[i].Type == kwJOIN ||
			tokens[i].Type == kwLEFT || tokens[i].Type == kwRIGHT || tokens[i].Type == kwFULL ||
			tokens[i].Type == kwCROSS || tokens[i].Type == kwINNER {
			return i
		}
		if tokens[i].Type == '(' {
			if c := r.match[i]; c >= 0 && c < end {
				i = c
				continue
			}
			return end
		}
	}
	return end
}

func (r *oracleRefReader) markUnsupportedTransform(ref RangeReference, start, end, want int) RangeReference {
	if ref.Unsupported {
		return ref
	}
	for i := start; i < end; i++ {
		if r.depth[i] != want {
			continue
		}
		switch r.toks[i].Type {
		case kwPIVOT:
			ref.Unsupported = true
			ref.UnsupportedReason = "pivot"
			return ref
		case kwUNPIVOT:
			ref.Unsupported = true
			ref.UnsupportedReason = "unpivot"
			return ref
		case kwMODEL:
			ref.Unsupported = true
			ref.UnsupportedReason = "model"
			return ref
		case ',', kwJOIN, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwCONNECT, kwSTART,
			kwUNION, kwINTERSECT, kwMINUS, kwEXCEPT, kwOFFSET, kwFETCH:
			return ref
		case '(':
			if c := r.match[i]; c >= 0 && c < end {
				i = c
				continue
			}
			return ref
		}
	}
	return ref
}

// enclosingSelects returns the SELECTs whose query blocks hold the cursor,
// outermost first: at each depth up to the cursor's, the last SELECT before
// the cursor inside the parentheses still open there. A SELECT in a group
// closed before the cursor heads none of them: the scalar subquery in
// SELECT a, (SELECT max(b) FROM u), | FROM t, or the IN subquery in a IN
// (SELECT b FROM u) AND nvl((SELECT 1 FROM w WHERE |), 0) = 1.
func (x *oracleTokenIndex) enclosingSelects(cursorOffset int) []int {
	// last holds, by depth, the last SELECT inside each '(' open, or -1.
	var last []int
	depth := 0
	for i, tok := range x.toks {
		if tok.Loc >= cursorOffset {
			break
		}
		depth = x.depth[i]
		switch tok.Type {
		case '(':
			last = last[:min(len(last), depth+1)]
			depth++
		case ')':
			depth = max(depth-1, 0)
		case kwSELECT:
			for len(last) <= depth {
				last = append(last, -1)
			}
			last[depth] = i
		}
	}
	var selects []int
	for _, sel := range last[:min(len(last), depth+1)] {
		if sel >= 0 {
			selects = append(selects, sel)
		}
	}
	return selects
}

// armBounds returns the tokens toks[start:end] of the set-operation arm
// the SELECT at selectIdx heads, hi bounding them: from after the last set
// operator before it at its depth inside the same parentheses, to the first
// set operator after it at its depth or the end of the parentheses around
// it.
func (r *oracleRefReader) armBounds(selectIdx, hi int) (int, int) {
	start := selectIdx
	if p := r.setOp[selectIdx]; p >= 0 {
		start = p + 1
	}
	return start, r.armEnd(selectIdx, hi)
}

// armEnd returns where armBounds ends the arm of the SELECT at selectIdx.
func (r *oracleRefReader) armEnd(selectIdx, hi int) int {
	want := r.depth[selectIdx]
	for i := selectIdx + 1; i < hi; i++ {
		if r.depth[i] < want {
			return i
		}
		if r.depth[i] != want {
			continue
		}
		switch r.toks[i].Type {
		case kwUNION, kwINTERSECT, kwMINUS, kwEXCEPT:
			return i
		case '(':
			if c := r.match[i]; c >= 0 && c < hi {
				i = c
				continue
			}
			return hi
		}
	}
	return hi
}

// outerSelectRefs returns the references of each SELECT in selects, the
// query blocks around the cursor's own, nearest first.
func (r *oracleRefReader) outerSelectRefs(cursorOffset int, selects []int) [][]RangeReference {
	if len(selects) == 0 {
		return nil
	}
	tokens := r.toks
	insideLateral := r.cursorInsideLateral(cursorOffset)
	cursorIdx := len(tokens)
	if insideLateral {
		cursorIdx = oracleTokenIndexAtOrAfterOffset(tokens, cursorOffset)
	}
	var levels [][]RangeReference
	for k := len(selects) - 1; k >= 0; k-- {
		sel := selects[k]
		start, end := r.armBounds(sel, len(tokens))
		if insideLateral && cursorIdx < end {
			end = cursorIdx
		}
		refs := r.selectRefsInRange(start, end, r.depth[sel], len(tokens))
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

func parseOracleRangeReference(tokens []Token, i int, ctes map[string]RangeReference) (int, RangeReference, bool) {
	return newOracleTokenIndex(tokens).reader(ctes).rangeReference(i, len(tokens))
}

// rangeReference reads the range reference at token i, which reads no
// token at or past hi.
func (r *oracleRefReader) rangeReference(i, hi int) (int, RangeReference, bool) {
	tokens := r.toks
	for i < hi && tokens[i].Type == ',' {
		i++
	}
	if i >= hi {
		return i, RangeReference{}, false
	}
	if tokens[i].Type == kwLATERAL {
		return r.lateralReference(i, hi)
	}
	if tokens[i].Type == kwXMLTABLE || tokens[i].Type == kwJSON_TABLE {
		return r.structuredFunctionReference(i, hi)
	}
	if tokens[i].Type == kwTABLE && i+1 < hi && tokens[i+1].Type == '(' {
		return r.tableCollectionReference(i, hi)
	}
	if oracleTokenIsIdentLikeStr(tokens[i], "EXTERNAL") && i+1 < hi && tokens[i+1].Type == '(' {
		return r.inlineExternalReference(i, hi)
	}
	if (oracleTokenIsIdentLikeStr(tokens[i], "CONTAINERS") || oracleTokenIsIdentLikeStr(tokens[i], "SHARDS")) &&
		i+1 < hi && tokens[i+1].Type == '(' {
		return r.containersReference(i, hi)
	}
	if tokens[i].Type == '(' {
		return r.subqueryReference(i, hi)
	}
	if !isOracleCompletionIdent(tokens[i]) {
		return i, RangeReference{}, false
	}

	startTok := tokens[i]
	parts := []Token{tokens[i]}
	i++
	for i+1 < hi && tokens[i].Type == '.' && isOracleCompletionIdent(tokens[i+1]) {
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

	if i < hi && tokens[i].Type == kwAS {
		i++
	}
	if i < hi && isOracleCompletionIdent(tokens[i]) && !isOracleCompletionClauseBoundary(tokens[i].Type) {
		ref.Alias = tokens[i].Str
		ref.Loc.End = tokens[i].End
		i++
	}
	if ref.Alias == "" {
		ref.Alias = ref.Name
	}
	return i, ref, true
}

func (r *oracleRefReader) lateralReference(i, hi int) (int, RangeReference, bool) {
	tokens := r.toks
	if i+1 >= hi || tokens[i+1].Type != '(' {
		return i, RangeReference{}, false
	}
	closeIdx := r.closeOf(i+1, hi)
	if closeIdx < 0 {
		return i, RangeReference{}, false
	}
	next, alias, ok := parseOracleOptionalAlias(tokens[:hi], closeIdx+1)
	if !ok {
		return closeIdx + 1, RangeReference{}, false
	}
	ref := RangeReference{
		Kind:    RangeReferenceSubquery,
		Name:    alias.Str,
		Alias:   alias.Str,
		Columns: r.selectListColumns(i+2, closeIdx),
		Loc:     nodes.Loc{Start: tokens[i].Loc, End: alias.End},
		BodyLoc: nodes.Loc{Start: tokens[i+1].End, End: tokens[closeIdx].Loc},
	}
	return next, ref, true
}

func (r *oracleRefReader) structuredFunctionReference(i, hi int) (int, RangeReference, bool) {
	tokens := r.toks
	if i+1 >= hi || tokens[i+1].Type != '(' {
		return i, RangeReference{}, false
	}
	closeIdx := r.closeOf(i+1, hi)
	if closeIdx < 0 {
		return i, RangeReference{}, false
	}
	next, alias, ok := parseOracleOptionalAlias(tokens[:hi], closeIdx+1)
	if !ok {
		return closeIdx + 1, RangeReference{}, false
	}
	ref := RangeReference{
		Kind:    RangeReferenceSubquery,
		Name:    alias.Str,
		Alias:   alias.Str,
		Columns: extractOracleStructuredFunctionColumns(tokens[i+2 : closeIdx]),
		Loc:     nodes.Loc{Start: tokens[i].Loc, End: alias.End},
		BodyLoc: nodes.Loc{Start: tokens[i+1].End, End: tokens[closeIdx].Loc},
	}
	return next, ref, true
}

func (r *oracleRefReader) inlineExternalReference(i, hi int) (int, RangeReference, bool) {
	tokens := r.toks
	closeIdx := r.closeOf(i+1, hi)
	if closeIdx < 0 {
		return i, RangeReference{}, false
	}
	next, alias, ok := parseOracleOptionalAlias(tokens[:hi], closeIdx+1)
	if !ok {
		return closeIdx + 1, RangeReference{}, false
	}
	var columns []string
	if i+2 < closeIdx && tokens[i+2].Type == '(' {
		if colClose := r.closeOf(i+2, hi); colClose > i+2 && colClose < closeIdx {
			columns = extractOracleColumnDefNames(tokens[i+3 : colClose])
		}
	}
	ref := RangeReference{
		Kind:    RangeReferenceSubquery,
		Name:    alias.Str,
		Alias:   alias.Str,
		Columns: columns,
		Loc:     nodes.Loc{Start: tokens[i].Loc, End: alias.End},
		BodyLoc: nodes.Loc{Start: tokens[i+1].End, End: tokens[closeIdx].Loc},
	}
	return next, ref, true
}

func (r *oracleRefReader) containersReference(i, hi int) (int, RangeReference, bool) {
	tokens := r.toks
	closeIdx := r.closeOf(i+1, hi)
	if closeIdx < 0 {
		return i, RangeReference{}, false
	}
	parts := oracleObjectNameParts(tokens[i+2 : closeIdx])
	if len(parts) == 0 {
		return closeIdx + 1, RangeReference{}, false
	}
	next, alias, hasAlias := parseOracleOptionalAlias(tokens[:hi], closeIdx+1)
	name := parts[len(parts)-1].Str
	aliasName := name
	locEnd := tokens[closeIdx].End
	if hasAlias {
		aliasName = alias.Str
		locEnd = alias.End
	} else {
		next = closeIdx + 1
	}
	ref := RangeReference{
		Kind:    RangeReferenceRelation,
		Name:    name,
		Alias:   aliasName,
		Loc:     nodes.Loc{Start: tokens[i].Loc, End: locEnd},
		BodyLoc: nodes.Loc{Start: tokens[i+1].End, End: tokens[closeIdx].Loc},
	}
	if len(parts) > 1 {
		ref.Schema = parts[len(parts)-2].Str
	}
	return next, ref, true
}

func (r *oracleRefReader) tableCollectionReference(i, hi int) (int, RangeReference, bool) {
	tokens := r.toks
	closeIdx := r.closeOf(i+1, hi)
	if closeIdx < 0 {
		return i, RangeReference{}, false
	}
	nextIdx := closeIdx + 1
	if nextIdx+2 < hi && tokens[nextIdx].Type == '(' && tokens[nextIdx+1].Type == '+' && tokens[nextIdx+2].Type == ')' {
		nextIdx += 3
	}
	next, alias, ok := parseOracleOptionalAlias(tokens[:hi], nextIdx)
	if !ok {
		return nextIdx, RangeReference{}, false
	}
	ref := RangeReference{
		Kind:              RangeReferenceSubquery,
		Name:              alias.Str,
		Alias:             alias.Str,
		Unsupported:       true,
		UnsupportedReason: "table_collection",
		Loc:               nodes.Loc{Start: tokens[i].Loc, End: alias.End},
		BodyLoc:           nodes.Loc{Start: tokens[i+1].End, End: tokens[closeIdx].Loc},
	}
	return next, ref, true
}

func (r *oracleRefReader) subqueryReference(i, hi int) (int, RangeReference, bool) {
	tokens := r.toks
	closeIdx := r.closeOf(i, hi)
	if closeIdx < 0 {
		return i, RangeReference{}, false
	}
	aliasIdx := closeIdx + 1
	if aliasIdx < hi && tokens[aliasIdx].Type == kwAS {
		aliasIdx++
	}
	if aliasIdx >= hi || !isOracleCompletionIdent(tokens[aliasIdx]) {
		return closeIdx + 1, RangeReference{}, false
	}
	ref := RangeReference{
		Kind:    RangeReferenceSubquery,
		Name:    tokens[aliasIdx].Str,
		Alias:   tokens[aliasIdx].Str,
		Columns: r.selectListColumns(i+1, closeIdx),
		Loc:     nodes.Loc{Start: tokens[i].Loc, End: tokens[aliasIdx].End},
		BodyLoc: nodes.Loc{Start: tokens[i].End, End: tokens[closeIdx].Loc},
	}
	return aliasIdx + 1, ref, true
}

// ctesOfQuery returns the CTEs of the statement's WITH clause, each body
// read with the CTEs before it.
func (r *oracleRefReader) ctesOfQuery() []RangeReference {
	tokens := r.toks
	first := firstOracleNonTerminator(tokens)
	if first < 0 || tokens[first].Type != kwWITH {
		return nil
	}
	var refs []RangeReference
	seen := make(map[string]RangeReference)
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
			closeIdx := r.closeOf(i, len(tokens))
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
		closeIdx := r.closeOf(i, len(tokens))
		if closeIdx < 0 {
			break
		}
		if len(ref.Columns) == 0 {
			// seen grows after each CTE, so each body gets a reader, and
			// columns read, of its own.
			ref.Columns = r.reader(seen).selectListColumns(i+1, closeIdx)
		}
		ref.BodyLoc = nodes.Loc{Start: tokens[i].End, End: tokens[closeIdx].Loc}
		ref.Loc.End = tokens[closeIdx].End
		refs = append(refs, ref)
		seen[strings.ToLower(ref.Name)] = ref
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

// querySelect returns the SELECT whose list names the columns of the query
// toks[lo:hi] holds, and the toks[qlo:qhi] it stands in: that of the first
// operand, the first SELECT at the query's own depth or, when a
// parenthesized operand opens the query, ((SELECT c1 FROM t) UNION (SELECT
// c2 FROM u)), the first inside it, as deep as the wrappers go. It reads
// only the tokens at the depth it looks for a SELECT at.
func (r *oracleRefReader) querySelect(lo, hi int) (qlo, qhi, selectIdx int) {
	tokens := r.toks
	if lo < hi && tokens[lo].Type == kwWITH {
		// After its WITH clause a query may open with a parenthesized
		// operand, WITH w AS (SELECT c FROM t) (SELECT c FROM w): the first
		// '(' at its own depth after a ')' that ends a WITH element.
		for i := lo; i < hi; i++ {
			switch tokens[i].Type {
			case kwSELECT:
				return lo, hi, i
			case '(':
				if i > lo && tokens[i-1].Type == ')' {
					if qlo, qhi, selectIdx = r.querySelect(i, hi); selectIdx < 0 {
						return 0, 0, -1
					}
					return qlo, qhi, selectIdx
				}
				c := r.closeOf(i, hi)
				if c < 0 {
					return 0, 0, -1
				}
				i = c
			}
		}
		return 0, 0, -1
	}
	run := 0
	for lo+run < hi && tokens[lo+run].Type == '(' {
		run++
	}
	// The first operand names the columns, so a query that opens with '('
	// takes its query from inside those parentheses, past any later operand
	// at its own depth: ((SELECT a FROM t) UNION SELECT b FROM u) names A.
	qlo, qhi = lo, hi
	if run > 0 {
		qlo = lo + run
		if qhi = r.closeOf(qlo-1, hi); qhi < 0 {
			return 0, 0, -1
		}
	}
	for i := qlo; i < qhi; i++ {
		switch tokens[i].Type {
		case kwSELECT:
			return qlo, qhi, i
		case '(':
			c := r.closeOf(i, qhi)
			if c < 0 {
				return 0, 0, -1
			}
			i = c
		}
	}
	return 0, 0, -1
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

func oracleTokenIndexAtOrAfterOffset(tokens []Token, offset int) int {
	for i, tok := range tokens {
		if tok.Loc >= offset {
			return i
		}
	}
	return len(tokens)
}

func (r *oracleRefReader) cursorInsideLateral(cursorOffset int) bool {
	tokens := r.toks
	for i, tok := range tokens {
		if tok.Loc >= cursorOffset {
			break
		}
		if tok.Type != kwLATERAL || i+1 >= len(tokens) || tokens[i+1].Type != '(' {
			continue
		}
		closeIdx := r.closeOf(i+1, len(tokens))
		if closeIdx >= 0 && tokens[i+1].End <= cursorOffset && cursorOffset <= tokens[closeIdx].Loc {
			return true
		}
	}
	return false
}

func parseOracleOptionalAlias(tokens []Token, i int) (int, Token, bool) {
	if i < len(tokens) && tokens[i].Type == kwAS {
		i++
	}
	if i >= len(tokens) || !isOracleCompletionIdent(tokens[i]) || isOracleCompletionClauseBoundary(tokens[i].Type) {
		return i, Token{}, false
	}
	return i + 1, tokens[i], true
}

func extractOracleStructuredFunctionColumns(tokens []Token) []string {
	for i, tok := range tokens {
		if tok.Type != kwCOLUMNS {
			continue
		}
		if i+1 < len(tokens) && tokens[i+1].Type == '(' {
			closeIdx := oracleMatchingParen(tokens, i+1)
			if closeIdx > i+1 {
				return extractOracleColumnDefNames(tokens[i+2 : closeIdx])
			}
			return nil
		}
		return extractOracleColumnDefNames(tokens[i+1:])
	}
	return nil
}

func extractOracleColumnDefNames(tokens []Token) []string {
	var names []string
	depth := 0
	expectName := true
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
		case ',':
			if depth == 0 {
				expectName = true
			}
			continue
		}
		if depth != 0 || !expectName {
			continue
		}
		if tok.Type == kwNESTED {
			expectName = false
			continue
		}
		if isOracleCompletionIdent(tok) {
			names = append(names, tok.Str)
			expectName = false
		}
		if i == len(tokens)-1 {
			break
		}
	}
	return names
}

func oracleObjectNameParts(tokens []Token) []Token {
	var parts []Token
	expectName := true
	for _, tok := range tokens {
		if expectName {
			if !isOracleCompletionIdent(tok) {
				break
			}
			parts = append(parts, tok)
			expectName = false
			continue
		}
		if tok.Type != '.' {
			break
		}
		expectName = true
	}
	return parts
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

func extractOracleSelectListColumns(tokens []Token, ctes map[string]RangeReference) []string {
	return newOracleTokenIndex(tokens).reader(ctes).selectListColumns(0, len(tokens))
}

// selectListColumns returns the columns the select list of the query in
// toks[lo:hi] names, read once for a derived table however many queries
// around it expand a star over it. The caller gets a copy of its own.
func (r *oracleRefReader) selectListColumns(lo, hi int) []string {
	if lo == 0 || r.toks[lo-1].Type != '(' || r.match[lo-1] != hi {
		return r.readSelectListColumns(lo, hi)
	}
	cols, ok := r.columns[lo]
	if !ok {
		cols = r.readSelectListColumns(lo, hi)
		if r.columns == nil {
			r.columns = map[int][]string{}
		}
		r.columns[lo] = cols
	}
	return append([]string(nil), cols...)
}

// readSelectListColumns computes selectListColumns.
func (r *oracleRefReader) readSelectListColumns(lo, hi int) []string {
	// A SELECT * reads the columns of each derived table it selects from,
	// in a call of its own on the tokens inside its parentheses, so the
	// calls nest no deeper than the parentheses do, a run of them counting
	// once (oracleParenRunsNestAtLeast). Completion is best effort: a query
	// whose parentheses nest this deep contributes no columns, rather than
	// recursing until the stack overflows. Between a pair of parentheses
	// the index holds the count.
	if lo > 0 && r.toks[lo-1].Type == '(' && r.match[lo-1] == hi {
		if r.runs[lo-1] >= max(maxNestingDepth, 1) {
			return nil
		}
	} else if oracleParenRunsNestAtLeast(r.toks[lo:hi], maxNestingDepth) {
		return nil
	}
	qlo, qhi, selectIdx := r.querySelect(lo, hi)
	if selectIdx < 0 {
		return nil
	}
	tokens := r.toks
	var names []string
	targetStart := selectIdx + 1
	closedTargetList := false
	for i := selectIdx + 1; i < qhi && !closedTargetList; i++ {
		switch tokens[i].Type {
		case '(':
			// What it holds is deeper than the list's commas.
			if i = r.closeOf(i, qhi); i < 0 {
				i = qhi
			}
		case kwFROM, kwUNION, kwINTERSECT, kwMINUS, kwEXCEPT:
			names = append(names, r.targetColumns(targetStart, i, qlo, qhi, selectIdx)...)
			closedTargetList = true
		case ',':
			names = append(names, r.targetColumns(targetStart, i, qlo, qhi, selectIdx)...)
			targetStart = i + 1
		}
	}
	if !closedTargetList && targetStart < qhi {
		names = append(names, r.targetColumns(targetStart, qhi, qlo, qhi, selectIdx)...)
	}
	return names
}

// oracleParenRunsNestAtLeast reports whether parentheses nest depth deep in
// tokens, a run of '(' counting as one level: the wrappers of a query come
// off in one pass (oracleRefReader.querySelect), and only the parentheses
// that open a derived table, which never follow a '(', lead to a call of
// their own.
func oracleParenRunsNestAtLeast(tokens []Token, depth int) bool {
	// runStarts says, for each '(' open, whether it opened a run.
	var runStarts []bool
	open := 0
	for i, tok := range tokens {
		switch tok.Type {
		case '(':
			start := i == 0 || tokens[i-1].Type != '('
			runStarts = append(runStarts, start)
			if start {
				if open++; open >= depth {
					return true
				}
			}
		case ')':
			if n := len(runStarts); n > 0 {
				if runStarts[n-1] {
					open--
				}
				runStarts = runStarts[:n-1]
			}
		}
	}
	return false
}

// targetColumns returns the columns the select item toks[start:end] names,
// in the query toks[lo:hi] whose SELECT is at selectIdx.
func (r *oracleRefReader) targetColumns(start, end, lo, hi, selectIdx int) []string {
	tokens := r.toks
	for start < end && tokens[start].Type == ',' {
		start++
	}
	for end > start && tokens[end-1].Type == ',' {
		end--
	}
	target := tokens[start:end]
	if len(target) == 0 {
		return nil
	}
	if len(target) == 1 && target[0].Type == '*' {
		return r.starColumns(lo, hi, selectIdx, "")
	}
	if len(target) == 3 && isOracleCompletionIdent(target[0]) && target[1].Type == '.' && target[2].Type == '*' {
		return r.starColumns(lo, hi, selectIdx, target[0].Str)
	}
	if alias, ok := r.explicitSelectAlias(start, end); ok {
		return []string{alias}
	}
	if name, ok := oracleBareSelectColumnName(target); ok {
		return []string{name}
	}
	return nil
}

// explicitSelectAlias returns the alias that ends the select item
// toks[start:end]: its last token outside parentheses, read back past each
// parenthesized group in one move.
func (r *oracleRefReader) explicitSelectAlias(start, end int) (string, bool) {
	tokens := r.toks
	for i := end - 1; i >= start; i-- {
		switch tokens[i].Type {
		case ')':
			if open := r.match[i]; open >= start {
				i = open
				continue
			}
			// Its '(' is before the item, or it has none: read back from
			// here, the rest of the item is inside it.
			return "", false
		case '(':
			continue
		}
		if i > start && tokens[i-1].Type == kwAS && isOracleCompletionIdent(tokens[i]) {
			return tokens[i].Str, true
		}
		if i > start && isOracleCompletionIdent(tokens[i]) && tokens[i-1].Type != '.' && !oracleTargetEndsWithPlainQualifiedName(tokens[start:i+1]) {
			return tokens[i].Str, true
		}
		break
	}
	return "", false
}

func oracleBareSelectColumnName(tokens []Token) (string, bool) {
	if len(tokens) == 1 && isOracleCompletionIdent(tokens[0]) {
		return tokens[0].Str, true
	}
	if len(tokens) >= 3 && isOracleCompletionIdent(tokens[len(tokens)-1]) && tokens[len(tokens)-2].Type == '.' {
		return tokens[len(tokens)-1].Str, true
	}
	return "", false
}

func oracleTargetEndsWithPlainQualifiedName(tokens []Token) bool {
	if len(tokens) == 1 && isOracleCompletionIdent(tokens[0]) {
		return true
	}
	if len(tokens) >= 3 && isOracleCompletionIdent(tokens[len(tokens)-1]) && tokens[len(tokens)-2].Type == '.' {
		return true
	}
	return false
}

// starColumns returns the columns a star in the select list of the SELECT
// at selectIdx, in the query toks[lo:hi], stands for: those of the
// references of its FROM clause qualifier names, or of all of them.
func (r *oracleRefReader) starColumns(lo, hi, selectIdx int, qualifier string) []string {
	refs := r.selectRefsInRange(selectIdx, r.armEnd(selectIdx, hi), r.depth[selectIdx], hi)
	var names []string
	for _, ref := range refs {
		ref = resolveOracleCTEReference(ref, r.ctes)
		if qualifier != "" && !strings.EqualFold(ref.Alias, qualifier) && !strings.EqualFold(ref.Name, qualifier) {
			continue
		}
		names = append(names, ref.Columns...)
	}
	return names
}

func isOracleCompletionIdent(tok Token) bool {
	if tok.Type == tokIDENT || tok.Type == tokQIDENT {
		return true
	}
	return tok.Type >= 2000 && !isOracleSQLReservedKeyword(tok)
}

func oracleTokenIsIdentLikeStr(tok Token, value string) bool {
	return isOracleCompletionIdent(tok) && strings.EqualFold(tok.Str, value)
}

func oracleIsFromClauseBoundary(tokenType int) bool {
	switch tokenType {
	case kwWHERE, kwGROUP, kwHAVING, kwORDER, kwCONNECT, kwSTART, kwUNION,
		kwINTERSECT, kwMINUS, kwEXCEPT, kwOFFSET, kwFETCH, kwMODEL, kwPIVOT, kwUNPIVOT:
		return true
	default:
		return false
	}
}

func isOracleCompletionClauseBoundary(tokenType int) bool {
	switch tokenType {
	case kwWHERE, kwGROUP, kwHAVING, kwORDER, kwCONNECT, kwSTART, kwUNION,
		kwINTERSECT, kwMINUS, kwEXCEPT, kwOFFSET, kwFETCH, kwJOIN, kwON, kwUSING,
		kwLEFT, kwRIGHT, kwFULL, kwCROSS, kwINNER, kwOUTER, kwMODEL, kwPIVOT, kwUNPIVOT:
		return true
	default:
		return false
	}
}
