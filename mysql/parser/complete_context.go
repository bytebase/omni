package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// CompletionContext is everything the MySQL parser knows about the cursor
// position that a completion caller cannot read off the AST, because the text
// is being typed and does not parse yet.
//
// It carries no catalog state: Candidates names the grammar positions the
// cursor accepts and the keywords they allow, Scope and CTEs name the relations
// and columns the parser can see without metadata, and Intent names the object
// class the qualifier before the cursor implies. A caller resolves those names
// against its own metadata and builds the suggestion list.
type CompletionContext struct {
	Candidates *CandidateSet

	// Scope is the best-effort relation scope of the query block holding the
	// cursor: LocalReferences for that block, OuterReferences for the blocks
	// around it, nearest first. References is both, in that order.
	Scope *ScopeSnapshot

	// CTEs are the WITH clause's definitions of the cursor statement,
	// whether or not the cursor is in one of their bodies.
	CTEs []RangeReference

	// SelectAliases are the aliases the statement's select lists declare
	// before the cursor. ORDER BY, GROUP BY, and HAVING may name them.
	SelectAliases []SelectAlias

	// Prefix is the partially typed identifier ending at the cursor.
	Prefix string

	// Intent is the object class and qualifier the cursor position implies,
	// or nil when the position accepts no catalog object.
	Intent *CompletionIntent

	// StatementLoc is the statement the cursor is in.
	StatementLoc nodes.Loc
}

// SelectAlias is an alias a select item declares, such as total in
// SELECT a + b AS total.
type SelectAlias struct {
	Name string
	Loc  nodes.Loc
}

// CompletionIntent describes the object class and the qualifier typed before
// the cursor.
type CompletionIntent struct {
	ObjectKinds []ObjectKind
	Qualifier   MultipartName
}

// ObjectKind classifies what the cursor position accepts. The set is MySQL's
// object classes, not the caller's suggestion categories: a caller maps the
// kinds it can resolve and ignores the rest.
type ObjectKind int

const (
	ObjectKindUnknown ObjectKind = iota
	ObjectKindDatabase
	ObjectKindTable
	ObjectKindView
	ObjectKindColumn
	ObjectKindIndex
	ObjectKindFunction
	ObjectKindProcedure
	ObjectKindTrigger
	ObjectKindEvent
	ObjectKindType
	ObjectKindCharset
	ObjectKindEngine
	ObjectKindVariable
	// ObjectKindSelectAlias is a select list alias the enclosing ORDER BY,
	// GROUP BY, or HAVING can name.
	ObjectKindSelectAlias
)

// MultipartName is the dot-separated name typed before the cursor.
//
// Parts holds the identifiers in the order typed, outermost first: db for
// "db.", {"db", "t"} for "db.t.", and {"t"} for "t.". MySQL reads "a.b" as
// table.column in a select list and as db.table in a FROM clause, so the parts
// stay flat rather than being assigned to fixed fields here.
type MultipartName struct {
	Parts []string
}

// Qualifier is the last part typed before the cursor: the name whose members
// the cursor is choosing among, or "" when nothing is qualified.
func (q MultipartName) Qualifier() string {
	if len(q.Parts) == 0 {
		return ""
	}
	return q.Parts[len(q.Parts)-1]
}

// ScopeSnapshot describes the relations visible at the cursor.
type ScopeSnapshot struct {
	// References is LocalReferences followed by each level of
	// OuterReferences, nearest first.
	References []RangeReference

	// LocalReferences is the cursor query block's own FROM scope.
	LocalReferences []RangeReference

	// OuterReferences holds the correlated outer query blocks, nearest first.
	OuterReferences [][]RangeReference

	// DMLTarget is the table an UPDATE, DELETE, INSERT, or REPLACE statement
	// names, wherever the cursor is in the statement.
	DMLTarget *RangeReference
}

// RangeReferenceKind classifies a relation the parser can see at the cursor.
type RangeReferenceKind int

const (
	RangeReferenceRelation RangeReferenceKind = iota
	RangeReferenceSubquery
	RangeReferenceFunction
	RangeReferenceJoinAlias
	RangeReferenceCTE
	RangeReferenceDMLTarget
)

// RangeReference is a syntax-level table expression. It carries no metadata:
// a caller resolves Table and the columns of a derived table or CTE against its
// own catalog.
type RangeReference struct {
	Kind RangeReferenceKind

	// Database and Table are as written; Database is "" when the reference
	// is not qualified.
	Database string
	Table    string

	// Alias is the name the rest of the query refers to this reference by,
	// or "" when it has none.
	Alias string

	// Columns are the column names this reference contributes that the
	// parser could read statically: a CTE's declared columns or its select
	// list, and a derived table's column list or select list. Empty when
	// they are not readable without metadata.
	Columns []string

	// AliasColumns is the column alias list written after the alias, as in
	// (SELECT ...) AS d (a, b). It repeats Columns when present.
	AliasColumns []string

	// Loc spans the reference including its alias; BodyLoc spans the
	// parenthesized body of a derived table or CTE.
	Loc     nodes.Loc
	BodyLoc nodes.Loc

	// Lateral marks a reference whose body may name the references before it.
	Lateral bool
}

// CollectCompletion returns completion candidates plus the parser-native
// scope, CTE, alias, and intent signals at cursorOffset. It is completion-only:
// Parse stays strict and returns no partial tree.
func CollectCompletion(sql string, cursorOffset int) *CompletionContext {
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	prefix, collectOffset := completionPrefix(sql, cursorOffset)
	ctx := &CompletionContext{
		Candidates:   Collect(sql, collectOffset),
		Prefix:       prefix,
		Scope:        collectCompletionScope(sql, cursorOffset),
		StatementLoc: completionStatementLoc(sql, cursorOffset),
	}
	ctx.CTEs = collectCompletionCTEs(sql, cursorOffset)
	ctx.SelectAliases = collectSelectAliases(sql, cursorOffset)
	ctx.Intent = collectCompletionIntent(sql, ctx.Candidates, collectOffset)
	return ctx
}

type completionScopeToken struct {
	Token
	depth int
	// close is the index of the ')' matching this '(' token, or -1. It is
	// filled once per token slice so a nested reader does not rescan the
	// parentheses at every level.
	close int
}

// collectCompletionScope returns the relation scope of the query block holding
// the cursor. It reads the statement the cursor is in, so another statement
// still being typed cannot cost this one its scope.
func collectCompletionScope(sql string, cursorOffset int) *ScopeSnapshot {
	tokens := completionScopeTokens(Tokenize(sql))
	start, end := completionStatementBounds(tokens, cursorOffset)
	if start >= end {
		return &ScopeSnapshot{}
	}
	stmtTokens := tokens[start:end]
	_, mainStart := collectCTEs(stmtTokens)
	// The whole statement, not just its main clause: the WITH clause holds
	// query blocks the cursor may be in.
	cteMap := rangeReferenceMap(collectVisibleCTEs(stmtTokens, cursorOffset))
	target, dmlRefs := collectDMLScope(stmtTokens[mainStart:], cteMap)
	var enclosing []RangeReference
	if first := firstNonTerminator(stmtTokens[mainStart:]); first >= 0 {
		switch stmtTokens[mainStart+first].Type {
		case kwUPDATE, kwDELETE:
			// A subquery in an UPDATE or DELETE can name the statement's
			// own tables, so they are its outermost scope level. An INSERT
			// or REPLACE target is not readable from its SELECT.
			enclosing = dmlRefs
		}
	}
	scope := collectSelectScope(stmtTokens, cursorOffset, cteMap, enclosing)
	if scope == nil {
		// No query block at the cursor: a DML or DDL statement's own tables
		// are what its clauses name.
		refs := append(dmlRefs, collectDDLScope(stmtTokens)...)
		scope = &ScopeSnapshot{
			LocalReferences: refs,
			References:      append([]RangeReference{}, refs...),
		}
	}
	if target != nil {
		scope.DMLTarget = target
	}
	return scope
}

func completionScopeTokens(tokens []Token) []completionScopeToken {
	depth := 0
	var open []int
	result := make([]completionScopeToken, 0, len(tokens))
	for _, tok := range tokens {
		result = append(result, completionScopeToken{Token: tok, depth: depth, close: -1})
		switch tok.Type {
		case '(':
			open = append(open, len(result)-1)
			depth++
		case ')':
			if depth > 0 {
				depth--
				if n := len(open); n > 0 {
					result[open[n-1]].close = len(result) - 1
					open = open[:n-1]
				}
			}
		}
	}
	return result
}

func completionStatementBounds(tokens []completionScopeToken, cursorOffset int) (int, int) {
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

// completionStatementLoc spans the statement the cursor is in, from its first
// byte past the previous terminator to its terminator or the end of the text.
func completionStatementLoc(sql string, cursorOffset int) nodes.Loc {
	tokens := Tokenize(sql)
	loc := nodes.Loc{Start: 0, End: len(sql)}
	for _, tok := range tokens {
		if tok.Type != ';' {
			continue
		}
		if tok.End <= cursorOffset {
			loc.Start = tok.End
			continue
		}
		loc.End = tok.Loc
		break
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

// collectCompletionCTEs reads the CTE definitions the cursor can name: the
// statement's WITH clause and every enclosing query block's, whether or not
// the cursor is in one of their bodies.
func collectCompletionCTEs(sql string, cursorOffset int) []RangeReference {
	tokens := completionScopeTokens(Tokenize(sql))
	start, end := completionStatementBounds(tokens, cursorOffset)
	if start >= end {
		return nil
	}
	return collectVisibleCTEs(tokens[start:end], cursorOffset)
}

// collectVisibleCTEs reads every WITH clause whose query block holds the
// cursor, outer blocks first, and replaces a definition an inner clause
// shadows.
func collectVisibleCTEs(tokens []completionScopeToken, cursorOffset int) []RangeReference {
	var ctes []RangeReference
	index := make(map[string]int)
	for i, tok := range tokens {
		if tok.Type != kwWITH || tok.Loc > cursorOffset {
			continue
		}
		if end := completionBlockEnd(tokens, i, tok.depth); end < len(tokens) && tokens[end].Loc < cursorOffset {
			continue
		}
		defs, _ := collectCTEsAt(tokens, i)
		for _, def := range defs {
			name := foldName(def.Table)
			if at, ok := index[name]; ok {
				ctes[at] = def
				continue
			}
			index[name] = len(ctes)
			ctes = append(ctes, def)
		}
	}
	return ctes
}

// completionBlockEnd returns where the query block starting at idx ends: the
// first token shallower than the block, which is the one past the ')' closing
// the group holding it, or the end of the statement.
func completionBlockEnd(tokens []completionScopeToken, idx, depth int) int {
	for i := idx + 1; i < len(tokens); i++ {
		if tokens[i].depth < depth {
			return i
		}
	}
	return len(tokens)
}

// collectCTEs parses a WITH clause at the start of tokens. It returns the
// definitions and the index of the statement's main clause.
func collectCTEs(tokens []completionScopeToken) ([]RangeReference, int) {
	first := firstNonTerminator(tokens)
	if first < 0 || tokens[first].Type != kwWITH {
		return nil, 0
	}
	return collectCTEsAt(tokens, first)
}

// collectCTEsAt parses the WITH clause at withIdx and returns its definitions
// and the index after the clause, where the statement's main clause starts.
func collectCTEsAt(tokens []completionScopeToken, withIdx int) ([]RangeReference, int) {
	var ctes []RangeReference
	idx := withIdx + 1
	depth := tokens[withIdx].depth
	if idx < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == kwRECURSIVE {
		idx++
	}
	for idx < len(tokens) {
		if tokens[idx].depth != depth || !isCompletionIdent(tokens[idx]) {
			break
		}
		name := tokens[idx].Str
		idx++
		var columns []string
		if idx < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == '(' {
			closeIdx := matchingParen(tokens, idx)
			if closeIdx < 0 {
				break
			}
			columns = identNames(tokens, idx+1, closeIdx)
			idx = closeIdx + 1
		}
		if idx >= len(tokens) || tokens[idx].depth != depth || tokens[idx].Type != kwAS {
			break
		}
		idx++
		if idx >= len(tokens) || tokens[idx].depth != depth || tokens[idx].Type != '(' {
			break
		}
		closeIdx := matchingParen(tokens, idx)
		if closeIdx < 0 {
			break
		}
		if len(columns) == 0 {
			columns = selectListColumns(tokens, idx+1, closeIdx)
		}
		ctes = append(ctes, RangeReference{
			Kind:    RangeReferenceCTE,
			Table:   name,
			Columns: columns,
			Loc:     tokenLoc(tokens[withIdx], tokens[closeIdx]),
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
		result[foldName(ref.Table)] = ref
	}
	return result
}

// collectDMLScope returns the reference an UPDATE, DELETE, INSERT, or REPLACE
// statement writes and every reference the statement brings into scope. The
// target is also the first of those references.
func collectDMLScope(tokens []completionScopeToken, ctes map[string]RangeReference) (*RangeReference, []RangeReference) {
	first := firstNonTerminator(tokens)
	if first < 0 {
		return nil, nil
	}
	switch tokens[first].Type {
	case kwUPDATE:
		// UPDATE [LOW_PRIORITY] [IGNORE] t1 [AS a] [, t2] SET ... — the
		// table list reads like a FROM clause.
		idx := first + 1
		for idx < len(tokens) && tokens[idx].depth == tokens[first].depth {
			switch tokens[idx].Type {
			case kwLOW_PRIORITY, kwIGNORE:
				idx++
				continue
			}
			break
		}
		refs := parseFromReferences(tokens, idx, collectClauseEnd(tokens, idx, tokens[first].depth), tokens[first].depth, ctes)
		if len(refs) == 0 {
			return nil, nil
		}
		refs[0].Kind = RangeReferenceDMLTarget
		return &refs[0], refs
	case kwDELETE:
		target, refs, ok := collectDeleteScope(tokens, first, ctes)
		if !ok {
			return nil, nil
		}
		return target, refs
	case kwINSERT, kwREPLACE:
		idx := first + 1
		// INSERT modifiers: LOW_PRIORITY, DELAYED, HIGH_PRIORITY, IGNORE.
		for idx < len(tokens) && tokens[idx].depth == tokens[first].depth {
			switch tokens[idx].Type {
			case kwLOW_PRIORITY, kwDELAYED, kwHIGH_PRIORITY, kwIGNORE:
				idx++
				continue
			}
			break
		}
		if idx < len(tokens) && tokens[idx].Type == kwINTO {
			idx++
		}
		refs, next, ok := parseTableSource(tokens, idx, ctes)
		if !ok || len(refs) == 0 {
			return nil, nil
		}
		refs[0].Kind = RangeReferenceDMLTarget
		// An INSERT target takes a column list, which names no alias.
		refs[0].Loc = tokenLoc(tokens[idx], tokens[next-1])
		return &refs[0], refs
	}
	return nil, nil
}

// collectDDLScope returns the table an ALTER TABLE, TRUNCATE, or CREATE INDEX
// statement names. Its clauses take that table's columns.
func collectDDLScope(tokens []completionScopeToken) []RangeReference {
	first := firstNonTerminator(tokens)
	if first < 0 {
		return nil
	}
	depth := tokens[first].depth
	switch tokens[first].Type {
	case kwALTER, kwTRUNCATE:
		idx := first + 1
		if idx < len(tokens) && tokens[idx].depth == depth && (tokens[idx].Type == kwTABLE || tokens[idx].Type == kwIGNORE) {
			idx++
		}
		// No alias follows a DDL target: MODIFY, ADD, and the other
		// clause keywords are identifiers the parser accepts, so reading
		// one as an alias would bind the wrong name.
		ref, next, ok := parseRangeReference(tokens, idx)
		if !ok {
			return nil
		}
		ref.Alias = ref.Table
		ref.Loc = tokenLoc(tokens[idx], tokens[next-1])
		return []RangeReference{ref}
	case kwCREATE:
		// CREATE [UNIQUE | FULLTEXT | SPATIAL] INDEX name ON table
		idx := first + 1
		for idx < len(tokens) && tokens[idx].depth == depth {
			switch tokens[idx].Type {
			case kwUNIQUE, kwFULLTEXT, kwSPATIAL:
				idx++
				continue
			}
			break
		}
		if idx >= len(tokens) || tokens[idx].Type != kwINDEX {
			return nil
		}
		onIdx := findKeywordAtDepth(tokens, kwON, idx+1, len(tokens), depth)
		if onIdx < 0 {
			return nil
		}
		ref, next, ok := parseRangeReference(tokens, onIdx+1)
		if !ok {
			return nil
		}
		ref.Alias = ref.Table
		ref.Loc = tokenLoc(tokens[onIdx+1], tokens[next-1])
		return []RangeReference{ref}
	}
	return nil
}

// collectDeleteScope reads DELETE's target and, for the multi-table forms
// DELETE t1, t2 FROM ... and DELETE FROM t1 USING ..., the references the
// statement reads.
func collectDeleteScope(tokens []completionScopeToken, deleteIdx int, ctes map[string]RangeReference) (*RangeReference, []RangeReference, bool) {
	depth := tokens[deleteIdx].depth
	idx := deleteIdx + 1
	fromIdx := findKeywordAtDepth(tokens, kwFROM, idx, len(tokens), depth)
	if fromIdx < 0 {
		return nil, nil, false
	}
	// Before FROM: DELETE t1, t2 FROM — the target list, else DELETE FROM.
	// Multi-table syntax 1 has no USING clause; its FROM is the whole
	// table_references list, and a USING there is a join condition.
	if fromIdx > idx {
		names := parseFromReferences(tokens, idx, fromIdx, depth, ctes)
		if len(names) > 0 {
			var targets []RangeReference
			for i := range names {
				names[i].Kind = RangeReferenceDMLTarget
				targets = append(targets, names[i])
			}
			refs := append([]RangeReference{}, targets...)
			refs = append(refs, parseFromReferences(tokens, fromIdx+1, len(tokens), depth, ctes)...)
			return &targets[0], refs, true
		}
	}
	refs, next, ok := parseTableSource(tokens, fromIdx+1, ctes)
	if !ok || len(refs) == 0 {
		return nil, nil, false
	}
	refs[0].Kind = RangeReferenceDMLTarget
	refs[0].Loc = tokenLoc(tokens[fromIdx+1], tokens[next-1])
	// DELETE FROM t1 USING t1 JOIN t2 carries the rest of the scope in its
	// USING clause. The USING must follow the target directly: a USING after
	// a join keyword is that join's condition.
	if next < len(tokens) && tokens[next].depth == depth && tokens[next].Type == kwUSING {
		refs = append(refs, parseFromReferences(tokens, next+1, len(tokens), depth, ctes)...)
	}
	return &refs[0], refs, true
}

// collectSelectScope returns the scope of the query block holding the cursor,
// or nil when the cursor is in no query block. An empty block keeps its outer
// references: only a statement with no query block at all falls back to its
// DML or DDL tables. enclosing holds the references of an UPDATE or DELETE the
// cursor's block is nested in; they are visible only from a nested block.
func collectSelectScope(tokens []completionScopeToken, cursorOffset int, ctes map[string]RangeReference, enclosing []RangeReference) *ScopeSnapshot {
	frames := selectFramesForCursor(tokens, cursorOffset, ctes)
	if len(frames) == 0 {
		return nil
	}
	innermost := frames[len(frames)-1]
	var outerRefs [][]RangeReference
	for i := len(frames) - 2; i >= 0; i-- {
		if refs := frames[i].refsVisibleFrom(cursorOffset); len(refs) > 0 {
			outerRefs = append(outerRefs, refs)
		}
	}
	if len(enclosing) > 0 && innermost.depth > 0 {
		if refs := enclosingScopeRefs(enclosing); len(refs) > 0 {
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
	}
}

// enclosingScopeRefs names the references of the enclosing DML statement once:
// the write target travels in Scope.DMLTarget, not as an outer reference.
func enclosingScopeRefs(refs []RangeReference) []RangeReference {
	out := make([]RangeReference, 0, len(refs))
	seen := make(map[string]bool, len(refs))
	for _, ref := range refs {
		if ref.Kind == RangeReferenceDMLTarget {
			ref.Kind = RangeReferenceRelation
		}
		key := foldName(ref.Database) + "." + foldName(ref.Table) + "|" + foldName(ref.Alias)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, ref)
	}
	return out
}

type completionSelectFrame struct {
	selectIdx int
	depth     int
	localRefs []RangeReference
}

// refsVisibleFrom returns the frame's references a nested subquery may see. A
// subquery inside a derived table's body sees only the references of the same
// FROM clause that precede it, and none of them unless the derived table is
// LATERAL.
func (f completionSelectFrame) refsVisibleFrom(cursorOffset int) []RangeReference {
	for _, ref := range f.localRefs {
		if ref.BodyLoc.Start < 0 || cursorOffset < ref.BodyLoc.Start || cursorOffset > ref.BodyLoc.End {
			continue
		}
		if !ref.Lateral {
			return nil
		}
		var refs []RangeReference
		for _, candidate := range f.localRefs {
			if candidate.Loc.Start >= 0 && candidate.Loc.Start < ref.Loc.Start {
				refs = append(refs, candidate)
			}
		}
		return refs
	}
	return f.localRefs
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
		if closedBeforeCursor(tokens, i, depth, cursorOffset) {
			// Its query block ended before the cursor: the parentheses
			// around it did.
			continue
		}
		fromIdx := findKeywordAtDepth(tokens, kwFROM, i+1, armEnd, depth)
		var local []RangeReference
		if fromIdx >= 0 {
			local = parseFromReferences(tokens, fromIdx+1, armEnd, depth, ctes)
		}
		frames = append(frames, completionSelectFrame{selectIdx: i, depth: depth, localRefs: local})
	}
	return frames
}

// closedBeforeCursor reports whether the parenthesis holding the SELECT at
// selectIdx closes before the cursor. The ')' carries the depth of the group's
// contents, the same depth as the SELECT inside it.
func closedBeforeCursor(tokens []completionScopeToken, selectIdx, depth, cursorOffset int) bool {
	if depth == 0 {
		return false
	}
	for i := selectIdx + 1; i < len(tokens); i++ {
		if tokens[i].depth == depth && tokens[i].Type == ')' {
			return tokens[i].Loc < cursorOffset
		}
	}
	return false
}

// selectArmEnd returns where the SELECT starting at start is no longer its own
// query block: the set operator that ends its arm, or the ')' that closes the
// group holding it. Crossing the ')' would read a sibling block's FROM clause
// at the same depth.
func selectArmEnd(tokens []completionScopeToken, start, depth int) int {
	for i := start; i < len(tokens); i++ {
		if tokens[i].depth != depth {
			continue
		}
		switch tokens[i].Type {
		case kwUNION, kwINTERSECT, kwEXCEPT, ')':
			return i
		}
	}
	return len(tokens)
}

// collectClauseEnd returns the index of the clause that ends a table list
// starting at start, which SET, WHERE, or a statement terminator does.
func collectClauseEnd(tokens []completionScopeToken, start, depth int) int {
	for i := start; i < len(tokens); i++ {
		if tokens[i].depth != depth {
			continue
		}
		if tokens[i].Type == kwSET || tokens[i].Type == kwWHERE || tokens[i].Type == ';' {
			return i
		}
	}
	return len(tokens)
}

func parseFromReferences(tokens []completionScopeToken, start, end, depth int, ctes map[string]RangeReference) []RangeReference {
	var refs []RangeReference
	expectSource := true
	for i := start; i < end && i < len(tokens); {
		tok := tokens[i]
		if tok.depth == depth && isFromClauseTerminator(tok.Type) {
			break
		}
		if tok.depth != depth {
			i++
			continue
		}
		if expectSource {
			sources, next, ok := parseTableSource(tokens, i, ctes)
			if ok && next > i && len(sources) > 0 {
				refs = append(refs, sources...)
				i = next
				expectSource = false
				continue
			}
			// The next identifier in this clause is not a source: text that
			// does not parse, such as a statement still being typed after
			// the cursor's, must not turn a later name into a reference.
			expectSource = false
			if next > i {
				// A group the reader refuses whole, such as one past the
				// nesting limit: skip it rather than rescan its tokens.
				i = next
				continue
			}
		}
		switch tok.Type {
		case ',', kwJOIN, kwSTRAIGHT_JOIN:
			expectSource = true
		}
		i++
	}
	return refs
}

// parseTableSource reads the table source at token idx and returns every
// reference it contributes: a plain table or CTE contributes one, and a
// parenthesized join contributes a reference per member.
func parseTableSource(tokens []completionScopeToken, idx int, ctes map[string]RangeReference) ([]RangeReference, int, bool) {
	if idx >= len(tokens) {
		return nil, idx, false
	}
	lateral := false
	if tokens[idx].Type == kwLATERAL {
		lateral = true
		idx++
		if idx >= len(tokens) {
			return nil, idx, false
		}
	}
	if tokens[idx].Type == '(' {
		return parseParenthesizedTableSource(tokens, idx, lateral, ctes)
	}
	if tokens[idx].Type == kwJSON_TABLE {
		return parseJSONTableSource(tokens, idx, lateral)
	}
	ref, next, ok := parseRangeReference(tokens, idx)
	if !ok {
		return nil, idx, false
	}
	ref.Lateral = lateral
	if ref.Database == "" {
		if cte, ok := ctes[foldName(ref.Table)]; ok {
			ref.Kind = RangeReferenceCTE
			ref.Columns = append([]string{}, cte.Columns...)
		}
	}
	next = skipPartitionClause(tokens, next, tokens[idx].depth)
	ref.Alias, ref.AliasColumns, next = parseAliasAndColumns(tokens, next, tokens[idx].depth)
	if len(ref.AliasColumns) > 0 {
		ref.Columns = append([]string{}, ref.AliasColumns...)
	}
	next = skipIndexHints(tokens, next, tokens[idx].depth)
	if next > idx && next <= len(tokens) {
		ref.Loc = tokenLoc(tokens[idx], tokens[next-1])
	}
	return []RangeReference{ref}, next, true
}

// skipPartitionClause consumes the PARTITION (p0, ...) selection that may
// follow a table name, before its alias.
func skipPartitionClause(tokens []completionScopeToken, idx, depth int) int {
	if idx >= len(tokens) || tokens[idx].depth != depth || tokens[idx].Type != kwPARTITION || idx+1 >= len(tokens) || tokens[idx+1].Type != '(' {
		return idx
	}
	closeIdx := matchingParen(tokens, idx+1)
	if closeIdx < 0 {
		return idx
	}
	return closeIdx + 1
}

// skipIndexHints consumes the USE/FORCE/IGNORE index hints that may follow a
// table reference's alias. A hint's FOR JOIN must not end the FROM clause the
// way a locking clause's FOR does.
func skipIndexHints(tokens []completionScopeToken, idx, depth int) int {
	for idx < len(tokens) && tokens[idx].depth == depth {
		switch tokens[idx].Type {
		case kwUSE, kwFORCE, kwIGNORE:
		default:
			return idx
		}
		next := idx + 1
		if next >= len(tokens) || tokens[next].depth != depth || (tokens[next].Type != kwINDEX && tokens[next].Type != kwKEY) {
			return idx
		}
		next++
		if next+1 < len(tokens) && tokens[next].depth == depth && tokens[next].Type == kwFOR {
			switch tokens[next+1].Type {
			case kwJOIN, kwGROUP:
				next += 2
			case kwORDER:
				next += 2
				if next < len(tokens) && tokens[next].depth == depth && tokens[next].Type == kwBY {
					next++
				}
			}
		}
		if next >= len(tokens) || tokens[next].depth != depth || tokens[next].Type != '(' {
			return idx
		}
		closeIdx := matchingParen(tokens, next)
		if closeIdx < 0 {
			return idx
		}
		idx = closeIdx + 1
		// One hint list may hold several hints, comma-separated.
		if idx+1 < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == ',' && tokens[idx+1].depth == depth {
			switch tokens[idx+1].Type {
			case kwUSE, kwFORCE, kwIGNORE:
				idx++
			}
		}
	}
	return idx
}

// parseParenthesizedTableSource reads a parenthesized source at the '(' at idx:
// a derived table, or a parenthesized join. A derived table contributes one
// reference whether or not it has an alias; a parenthesized join contributes its
// members' references, or one reference when it has an alias of its own.
func parseParenthesizedTableSource(tokens []completionScopeToken, idx int, lateral bool, ctes map[string]RangeReference) ([]RangeReference, int, bool) {
	// Each nested group re-enters this function through parseFromReferences,
	// so the group's parenthesis depth is the recursion depth. Stop at the
	// depth Parse itself rejects: the statement is not valid input, and a
	// reader that kept descending would not return a scope either. The caller
	// skips the group.
	closeIdx := matchingParen(tokens, idx)
	if tokens[idx].depth >= maxNestingDepth {
		if closeIdx < 0 {
			return nil, idx + 1, false
		}
		return nil, closeIdx + 1, false
	}
	if closeIdx < 0 {
		// Unterminated group: its body is what is being typed, and the
		// cursor's own scope reader will read it.
		return nil, idx + 1, false
	}
	depth := tokens[idx].depth
	isQuery := idx+1 < closeIdx && (tokens[idx+1].Type == kwSELECT || tokens[idx+1].Type == kwWITH || tokens[idx+1].Type == kwVALUES)
	if !isQuery {
		inner := parseFromReferences(tokens, idx+1, closeIdx, depth+1, ctes)
		alias, columns, next := parseAliasAndColumns(tokens, closeIdx+1, depth)
		if alias == "" {
			return inner, closeIdx + 1, len(inner) > 0
		}
		ref := RangeReference{
			Kind:         RangeReferenceJoinAlias,
			Table:        alias,
			Alias:        alias,
			AliasColumns: columns,
			Columns:      joinAliasColumns(inner),
			Loc:          tokenLoc(tokens[idx], tokens[next-1]),
			BodyLoc:      tokenLoc(tokens[idx], tokens[closeIdx]),
			Lateral:      lateral,
		}
		if len(columns) > 0 {
			ref.Columns = append([]string{}, columns...)
		}
		return []RangeReference{ref}, next, true
	}
	ref := RangeReference{
		Kind:    RangeReferenceSubquery,
		Columns: selectListColumns(tokens, idx+1, closeIdx),
		BodyLoc: tokenLoc(tokens[idx], tokens[closeIdx]),
		Loc:     tokenLoc(tokens[idx], tokens[closeIdx]),
		Lateral: lateral,
	}
	alias, columns, next := parseAliasAndColumns(tokens, closeIdx+1, depth)
	if alias == "" {
		return []RangeReference{ref}, closeIdx + 1, false
	}
	ref.Alias = alias
	ref.Table = alias
	ref.AliasColumns = columns
	if len(columns) > 0 {
		ref.Columns = append([]string{}, columns...)
	}
	if next > closeIdx+1 {
		ref.Loc.End = tokens[next-1].End
	}
	return []RangeReference{ref}, next, true
}

// parseJSONTableSource reads MySQL's only built-in table function,
// JSON_TABLE(expr, path COLUMNS (...)) [AS] alias. Its COLUMNS clause names
// the columns in the text, so the reference carries them like a CTE's.
func parseJSONTableSource(tokens []completionScopeToken, idx int, lateral bool) ([]RangeReference, int, bool) {
	if idx+1 >= len(tokens) || tokens[idx+1].Type != '(' {
		return nil, idx, false
	}
	closeIdx := matchingParen(tokens, idx+1)
	if closeIdx < 0 {
		// Unterminated arguments: the cursor is inside them.
		return nil, idx + 1, false
	}
	ref := RangeReference{
		Kind:    RangeReferenceFunction,
		Table:   "JSON_TABLE",
		Columns: jsonTableColumns(tokens, idx+1, closeIdx),
		Loc:     tokenLoc(tokens[idx], tokens[closeIdx]),
		BodyLoc: tokenLoc(tokens[idx+1], tokens[closeIdx]),
		Lateral: lateral,
	}
	alias, columns, next := parseAliasAndColumns(tokens, closeIdx+1, tokens[idx].depth)
	if alias != "" {
		ref.Alias = alias
		ref.Table = alias
		ref.AliasColumns = columns
		if len(columns) > 0 {
			ref.Columns = append([]string{}, columns...)
		}
		if next > closeIdx+1 {
			ref.Loc.End = tokens[next-1].End
		}
	}
	return []RangeReference{ref}, next, true
}

// jsonTableColumns reads the column names a JSON_TABLE COLUMNS clause declares,
// including the columns of a nested COLUMNS clause. open is the '(' holding the
// function's arguments.
func jsonTableColumns(tokens []completionScopeToken, open, close int) []string {
	depth := tokens[open].depth + 1
	for i := open + 1; i < close; i++ {
		if tokens[i].depth != depth || tokens[i].Type != kwCOLUMNS {
			continue
		}
		if i+1 >= close || tokens[i+1].Type != '(' {
			return nil
		}
		columnsClose := matchingParen(tokens, i+1)
		if columnsClose < 0 || columnsClose > close {
			return nil
		}
		return jsonTableColumnItems(tokens, i+2, columnsClose, tokens[i+1].depth+1)
	}
	return nil
}

// jsonTableColumnItems reads the name each column definition in tokens[lo:hi]
// starts with, at itemDepth, and recurses into a NESTED COLUMNS clause.
func jsonTableColumnItems(tokens []completionScopeToken, lo, hi, itemDepth int) []string {
	var names []string
	for i := lo; i < hi; i++ {
		if tokens[i].depth != itemDepth {
			continue
		}
		if i > lo && tokens[i-1].depth == itemDepth && tokens[i-1].Type != ',' {
			continue
		}
		switch {
		case tokens[i].Type == kwNESTED:
			for j := i + 1; j < hi; j++ {
				if tokens[j].depth != itemDepth || tokens[j].Type != kwCOLUMNS {
					continue
				}
				if j+1 < hi && tokens[j+1].Type == '(' {
					if nested := matchingParen(tokens, j+1); nested > 0 && nested <= hi {
						names = append(names, jsonTableColumnItems(tokens, j+2, nested, tokens[j+1].depth+1)...)
					}
				}
				break
			}
		case isCompletionIdent(tokens[i]):
			names = append(names, tokens[i].Str)
		}
	}
	return names
}

func joinAliasColumns(refs []RangeReference) []string {
	var columns []string
	for _, ref := range refs {
		if len(ref.Columns) > 0 {
			columns = append(columns, ref.Columns...)
			continue
		}
		// A plain table's columns need metadata, so the parser cannot name
		// them; an empty list says so.
		return nil
	}
	return columns
}

func parseRangeReference(tokens []completionScopeToken, idx int) (RangeReference, int, bool) {
	if idx >= len(tokens) || !isCompletionIdent(tokens[idx]) {
		return RangeReference{}, idx, false
	}
	var parts []string
	i := idx
	for i < len(tokens) && isCompletionIdent(tokens[i]) {
		parts = append(parts, tokens[i].Str)
		if i+1 < len(tokens) && tokens[i+1].Type == '.' && (i+2 >= len(tokens) || !isCompletionIdent(tokens[i+2])) {
			// "db." with nothing after is a name being typed, not a
			// reference.
			return RangeReference{}, idx, false
		}
		if i+2 >= len(tokens) || tokens[i+1].Type != '.' || !isCompletionIdent(tokens[i+2]) {
			i++
			break
		}
		i += 2
	}
	ref := RangeReference{Kind: RangeReferenceRelation}
	switch len(parts) {
	case 1:
		ref.Table = parts[0]
	default:
		ref.Database = parts[len(parts)-2]
		ref.Table = parts[len(parts)-1]
	}
	return ref, i, true
}

func parseAliasAndColumns(tokens []completionScopeToken, idx, depth int) (string, []string, int) {
	if idx < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == kwAS {
		idx++
	}
	if idx >= len(tokens) || tokens[idx].depth != depth || !isCompletionIdent(tokens[idx]) {
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
			columns = identNames(tokens, idx+1, closeIdx)
			idx = closeIdx + 1
		}
	}
	return alias, columns, idx
}

// selectListColumns returns the columns the select list of the query in
// tokens[lo:hi] names: the explicit alias, the bare column name, or the star's
// references. A position that is not one of those contributes nothing.
func selectListColumns(tokens []completionScopeToken, lo, hi int) []string {
	if lo < 0 || hi > len(tokens) || lo >= hi {
		return nil
	}
	depth := tokens[lo].depth
	if tokens[lo].Type == '(' {
		depth = tokens[lo].depth + 1
	}
	selectIdx := -1
	for i := lo; i < hi; i++ {
		if tokens[i].Type == kwSELECT {
			selectIdx = i
			depth = tokens[i].depth
			break
		}
	}
	if selectIdx < 0 {
		return nil
	}
	fromIdx := selectListEnd(tokens, selectIdx+1, hi, depth)
	var names []string
	itemStart := selectIdx + 1
	for i := itemStart; i <= fromIdx; i++ {
		if i < fromIdx && !(tokens[i].depth == depth && tokens[i].Type == ',') {
			continue
		}
		if name, ok := selectItemColumn(tokens, itemStart, i, depth); ok {
			names = append(names, name)
		}
		itemStart = i + 1
	}
	return names
}

// selectListEnd returns where the select list starting at start ends: the next
// clause of the same query block — FROM, WHERE, GROUP BY, HAVING, WINDOW, ORDER
// BY, LIMIT, INTO, FOR, LOCK, PROCEDURE —, the ')' closing the block, or the set
// operator that ends its arm.
func selectListEnd(tokens []completionScopeToken, start, hi, depth int) int {
	for i := start; i < hi && i < len(tokens); i++ {
		if tokens[i].depth != depth {
			continue
		}
		switch tokens[i].Type {
		case kwFROM, kwUNION, kwINTERSECT, kwEXCEPT, ')', kwWHERE, kwGROUP,
			kwHAVING, kwWINDOW, kwORDER, kwLIMIT, kwINTO, kwFOR, kwLOCK, kwPROCEDURE:
			return i
		}
	}
	return hi
}

// selectItemColumn returns the name the select item tokens[start:end] gives to
// its column: its alias when it declares one, else the column it names.
func selectItemColumn(tokens []completionScopeToken, start, end, depth int) (string, bool) {
	start, end = trimSelectModifiers(tokens, start, end, depth)
	if start >= end {
		return "", false
	}
	if alias, ok := selectItemAlias(tokens, start, end, depth); ok {
		return alias.Name, true
	}
	if end-start == 1 && isCompletionIdent(tokens[start]) {
		return tokens[start].Str, true
	}
	last := tokens[end-1]
	if end-start >= 3 && isCompletionIdent(last) && tokens[end-2].Type == '.' {
		return last.Str, true
	}
	return "", false
}

// trimSelectModifiers drops the select options that may precede the first
// select item, such as DISTINCT: they are keywords, never a column's name.
func trimSelectModifiers(tokens []completionScopeToken, start, end, depth int) (int, int) {
	for start < end && tokens[start].depth == depth {
		switch tokens[start].Type {
		case kwDISTINCT, kwALL, kwHIGH_PRIORITY, kwSTRAIGHT_JOIN, kwSQL_SMALL_RESULT,
			kwSQL_BIG_RESULT, kwSQL_BUFFER_RESULT, kwSQL_NO_CACHE, kwSQL_CALC_FOUND_ROWS:
			start++
			continue
		}
		break
	}
	return start, end
}

// collectSelectAliases returns the aliases of the select lists the cursor's own
// query block declares. An alias of a query block around the cursor is not
// visible here: inside a derived table body, the outer query's aliases are not
// in scope, and outside it, the body's are not either.
func collectSelectAliases(sql string, cursorOffset int) []SelectAlias {
	tokens := completionScopeTokens(Tokenize(sql))
	start, end := completionStatementBounds(tokens, cursorOffset)
	if start >= end {
		return nil
	}
	frames := selectFramesForCursor(tokens[start:end], cursorOffset, nil)
	if len(frames) == 0 {
		return nil
	}
	selectIdx := frames[len(frames)-1].selectIdx + start
	depth := tokens[selectIdx].depth
	limit := selectArmEnd(tokens, selectIdx+1, depth)
	listEnd := selectListEnd(tokens, selectIdx+1, limit, depth)
	var aliases []SelectAlias
	seen := make(map[string]bool)
	itemStart := selectIdx + 1
	for j := itemStart; j <= listEnd; j++ {
		if j < listEnd && !(tokens[j].depth == depth && tokens[j].Type == ',') {
			continue
		}
		aliasStart, aliasEnd := trimSelectModifiers(tokens, itemStart, j, depth)
		if alias, ok := selectItemAlias(tokens, aliasStart, aliasEnd, depth); ok && !seen[foldName(alias.Name)] {
			seen[foldName(alias.Name)] = true
			aliases = append(aliases, alias)
		}
		itemStart = j + 1
	}
	return aliases
}

// selectItemAlias returns the alias the select item tokens[start:end] declares:
// the name after AS, a string literal, or a bare identifier that follows a
// complete expression. A qualified name's last part is the column, not an
// alias, and the trailing operand of an expression such as a + b is not one
// either.
func selectItemAlias(tokens []completionScopeToken, start, end, depth int) (SelectAlias, bool) {
	if end <= start || end > len(tokens) {
		return SelectAlias{}, false
	}
	last := tokens[end-1]
	if last.depth != depth || (last.Type != tokSCONST && !isCompletionIdent(last)) {
		return SelectAlias{}, false
	}
	if isAliasBoundary(last.Type) {
		return SelectAlias{}, false
	}
	if end-start >= 2 && tokens[end-2].Type == '.' {
		return SelectAlias{}, false
	}
	// A single token names the column itself, not an alias of it.
	if end-start < 2 {
		return SelectAlias{}, false
	}
	if prev := tokens[end-2]; prev.Type != kwAS && !endsSelectExpression(prev) {
		return SelectAlias{}, false
	}
	return SelectAlias{Name: last.Str, Loc: tokenLoc(last, last)}, true
}

// endsSelectExpression reports whether tok can end a select item's expression,
// so that the identifier after it is an implicit alias rather than the next
// operand.
func endsSelectExpression(tok completionScopeToken) bool {
	switch tok.Type {
	case tokIDENT, ')', '*', tokSCONST, tokICONST, tokFCONST, tokXCONST, tokBCONST:
		return true
	case kwNULL, kwTRUE, kwFALSE, kwEND:
		return true
	}
	return tok.Type >= 700 && !isReserved(tok.Type)
}

func identNames(tokens []completionScopeToken, lo, hi int) []string {
	var names []string
	for i := lo; i < hi && i < len(tokens); i++ {
		if isCompletionIdent(tokens[i]) {
			names = append(names, tokens[i].Str)
		}
	}
	return names
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
		if tok.Type != ';' {
			return i
		}
	}
	return -1
}

func matchingParen(tokens []completionScopeToken, openIdx int) int {
	if openIdx < 0 || openIdx >= len(tokens) || tokens[openIdx].Type != '(' {
		return -1
	}
	return tokens[openIdx].close
}

func tokenLoc(first, last completionScopeToken) nodes.Loc {
	return nodes.Loc{Start: first.Loc, End: last.End}
}

// isCompletionIdent reports whether a token can name an object here: a plain
// or backtick-quoted identifier, or a keyword MySQL accepts as one.
func isCompletionIdent(tok completionScopeToken) bool {
	return tok.Type == tokIDENT || (tok.Type >= 700 && !isReserved(tok.Type))
}

func isFromClauseTerminator(typ int) bool {
	switch typ {
	case kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwLIMIT, kwOFFSET,
		kwUNION, kwINTERSECT, kwEXCEPT, kwFOR, kwLOCK, kwINTO, kwPROCEDURE:
		return true
	default:
		return false
	}
}

// isAliasBoundary reports whether a token ends a table reference rather than
// naming it: a clause keyword, a join keyword, or one of MySQL's index hints.
func isAliasBoundary(typ int) bool {
	switch typ {
	case kwON, kwUSING, kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwORDER, kwLIMIT,
		kwUNION, kwINTERSECT, kwEXCEPT, kwFOR, kwLOCK, kwINTO, kwSET, kwVALUES,
		kwJOIN, kwSTRAIGHT_JOIN, kwINNER, kwCROSS, kwNATURAL, kwLEFT, kwRIGHT,
		kwOUTER, kwAS, kwUSE, kwFORCE, kwIGNORE, kwPARTITION:
		return true
	default:
		return false
	}
}

// completionPrefix returns the identifier the cursor is inside or at the end
// of, and the offset the caller should complete from: the identifier's first
// byte, so the grammar sees the position before the partial name. The lexer
// decides the token, so a backtick-quoted name keeps its spaces and a
// non-ASCII name is not cut short.
func completionPrefix(sql string, cursorOffset int) (string, int) {
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	for _, tok := range Tokenize(sql) {
		if tok.Loc > cursorOffset || tok.End < cursorOffset {
			continue
		}
		if tok.Type != tokIDENT && tok.Type < 700 {
			break
		}
		if tok.Loc == cursorOffset {
			break
		}
		if sql[tok.Loc] == '`' {
			if tok.End <= cursorOffset {
				return tok.Str, tok.Loc
			}
			return unquoteIdent(sql[tok.Loc+1 : cursorOffset]), tok.Loc
		}
		return sql[tok.Loc:cursorOffset], tok.Loc
	}
	return "", cursorOffset
}

// unquoteIdent drops the doubled backtick escaping from a partially typed
// quoted identifier, so the prefix matches the names the parser reports.
func unquoteIdent(name string) string {
	if !strings.Contains(name, "``") {
		return name
	}
	return strings.ReplaceAll(name, "``", "`")
}

// collectCompletionIntent reads the object class the cursor position accepts
// from the grammar's candidates and the qualifier typed before it.
func collectCompletionIntent(sql string, candidates *CandidateSet, collectOffset int) *CompletionIntent {
	if candidates == nil {
		return nil
	}
	tokens := completionScopeTokens(Tokenize(sql[:collectOffset]))
	qualifier := qualifierBeforeOffset(tokens)
	kinds := objectKindsForCandidates(candidates, len(qualifier.Parts))
	if len(kinds) == 0 && len(qualifier.Parts) == 0 {
		return nil
	}
	if qualifierAllowedSelectAlias(tokens) {
		kinds = appendObjectKind(kinds, ObjectKindSelectAlias)
	}
	return &CompletionIntent{
		ObjectKinds: kinds,
		Qualifier:   qualifier,
	}
}

// qualifierAllowedSelectAlias reports whether the cursor is in an ORDER BY,
// GROUP BY, or HAVING clause of the query block it is in, where a select list
// alias may be named. Only the cursor block's tokens decide: a completed
// subquery before the cursor ends its own clause, so its ORDER BY does not
// make an outer WHERE an alias position.
func qualifierAllowedSelectAlias(tokens []completionScopeToken) bool {
	cursorDepth := 0
	if len(tokens) > 0 {
		cursorDepth = tokens[len(tokens)-1].depth
	}
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].depth != cursorDepth {
			continue
		}
		switch tokens[i].Type {
		case kwORDER, kwGROUP, kwHAVING:
			return true
		case kwWHERE, kwON, kwFROM, kwSELECT, kwLIMIT, kwFOR, kwSET, kwVALUES:
			return false
		}
	}
	return false
}

// objectKindsForCandidates returns the object classes the cursor position
// accepts, given how many name parts are already typed. MySQL reads a dotted
// name left to right, so each part narrows what the next one may be: in a
// select list, the first part may start a database, a table, a view, or a
// column, the second a table or a column, and the third only a column.
func objectKindsForCandidates(candidates *CandidateSet, typedParts int) []ObjectKind {
	var kinds []ObjectKind
	for _, rc := range candidates.Rules {
		switch rc.Rule {
		case "database_ref":
			kinds = appendObjectKind(kinds, ObjectKindDatabase)
		case "table_ref":
			if typedParts == 0 {
				kinds = appendObjectKind(kinds, ObjectKindDatabase)
			}
			kinds = appendObjectKind(kinds, ObjectKindTable)
			kinds = appendObjectKind(kinds, ObjectKindView)
		case "view_ref":
			if typedParts == 0 {
				kinds = appendObjectKind(kinds, ObjectKindDatabase)
			}
			kinds = appendObjectKind(kinds, ObjectKindView)
		case "columnref":
			switch typedParts {
			case 0:
				kinds = appendObjectKind(kinds, ObjectKindDatabase)
				kinds = appendObjectKind(kinds, ObjectKindTable)
				kinds = appendObjectKind(kinds, ObjectKindView)
			case 1:
				kinds = appendObjectKind(kinds, ObjectKindTable)
				kinds = appendObjectKind(kinds, ObjectKindView)
			}
			kinds = appendObjectKind(kinds, ObjectKindColumn)
		case "index_ref":
			kinds = appendObjectKind(kinds, ObjectKindIndex)
		case "func_name", "function_ref":
			kinds = appendObjectKind(kinds, ObjectKindFunction)
		case "procedure_ref":
			kinds = appendObjectKind(kinds, ObjectKindProcedure)
		case "trigger_ref":
			kinds = appendObjectKind(kinds, ObjectKindTrigger)
		case "event_ref":
			kinds = appendObjectKind(kinds, ObjectKindEvent)
		case "type_name":
			kinds = appendObjectKind(kinds, ObjectKindType)
		case "charset":
			kinds = appendObjectKind(kinds, ObjectKindCharset)
		case "engine":
			kinds = appendObjectKind(kinds, ObjectKindEngine)
		case "variable":
			kinds = appendObjectKind(kinds, ObjectKindVariable)
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

// qualifierBeforeOffset reads the dotted name the cursor continues: "db.t" for
// "db.t.", and nothing while the cursor is in the first part.
func qualifierBeforeOffset(tokens []completionScopeToken) MultipartName {
	idx := -1
	for i, tok := range tokens {
		if tok.Type != ';' {
			idx = i
		}
	}
	if idx < 0 || tokens[idx].Type != '.' {
		return MultipartName{}
	}
	idx--
	var parts []string
	for idx >= 0 && isCompletionIdent(tokens[idx]) {
		parts = append([]string{tokens[idx].Str}, parts...)
		if idx-1 < 0 || tokens[idx-1].Type != '.' {
			break
		}
		idx -= 2
	}
	return MultipartName{Parts: parts}
}

func foldName(name string) string {
	if name == "" {
		return ""
	}
	buf := make([]byte, len(name))
	for i := range len(name) {
		c := name[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		buf[i] = c
	}
	return string(buf)
}
