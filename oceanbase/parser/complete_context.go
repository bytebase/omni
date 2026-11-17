// This file is OceanBase's copy of mysql/parser/complete_context.go. The design,
// the boundary, and the field meanings are documented there and in
// docs/specs/2026-11-16-completion-cut-design.md; this copy exists because a
// fork carries its own AST.
//
// The reader itself needs no dialect delta. OceanBase accepts UNION, INTERSECT,
// and EXCEPT (the set_type compat rows), which is the list selectArmEnd stops
// at, and its hints are comments, which the lexer skips. One engine form is a
// gap in this parser rather than in the reader: OceanBase accepts MINUS as a set
// operator (parser/testdata/compat.tsv records set_type/MINUS as omni=reject,
// engine=accept), so until the parser takes it the reader reads MINUS as the
// identifier it currently is. The corpus pins the resulting signals and names
// the row.
package parser

import (
	"strconv"
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
	// ObjectKindConstraint is a constraint name, as DROP FOREIGN KEY and
	// DROP CONSTRAINT take.
	ObjectKindConstraint
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
	// parenthesized body of a derived table or CTE. BodyLoc.End is -1 when
	// the group is unterminated and the body runs to the end of the text.
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
	// The reader works on the Split segment holding the cursor: a script with
	// its own DELIMITER separates its statements the way the parser does, and
	// the candidate collection already reads the same segment. The prefix is
	// read there too, so an identifier cannot start in another statement.
	text, base, local := completionSegment(sql, cursorOffset)
	prefix, localStart := completionPrefix(text, local)
	ctx := &CompletionContext{
		Candidates:   Collect(sql, localStart+base),
		Prefix:       prefix,
		Scope:        shiftScope(collectCompletionScope(text, local), base),
		StatementLoc: shiftLoc(completionStatementLoc(text, local), base),
	}
	ctx.CTEs = shiftReferences(collectCompletionCTEs(text, local), base)
	ctx.SelectAliases = shiftSelectAliases(collectSelectAliases(text, local), base)
	ctx.Intent = collectCompletionIntent(text, ctx.Candidates, localStart)
	return ctx
}

// completionSegment returns the text of the statement holding the cursor, the
// byte offset that text starts at in sql, and the cursor's offset within the
// text.
func completionSegment(sql string, cursorOffset int) (string, int, int) {
	segments := Split(sql)
	for _, seg := range segments {
		if seg.Empty() {
			continue
		}
		if cursorOffset >= seg.ByteStart && cursorOffset <= seg.ByteEnd {
			return seg.Text, seg.ByteStart, cursorOffset - seg.ByteStart
		}
	}
	if len(segments) > 0 {
		// Past the last statement: what follows its terminator is not part of
		// it, so the cursor completes a new statement with no scope from the
		// one before. A custom DELIMITER's own text sits there too.
		if last := segments[len(segments)-1]; cursorOffset >= last.ByteEnd {
			return sql[cursorOffset:], cursorOffset, 0
		}
	}
	return sql, 0, cursorOffset
}

func shiftLoc(loc nodes.Loc, base int) nodes.Loc {
	if base == 0 {
		return loc
	}
	// An unknown endpoint stays unknown: the convention for an open-ended
	// body is End == -1, and adding the segment base would invent a position.
	start, end := loc.Start, loc.End
	if start >= 0 {
		start += base
	}
	if end >= 0 {
		end += base
	}
	return nodes.Loc{Start: start, End: end}
}

func shiftReferences(refs []RangeReference, base int) []RangeReference {
	if base == 0 || len(refs) == 0 {
		return refs
	}
	out := make([]RangeReference, len(refs))
	copy(out, refs)
	for i := range out {
		out[i].Loc = shiftLoc(out[i].Loc, base)
		out[i].BodyLoc = shiftLoc(out[i].BodyLoc, base)
	}
	return out
}

func shiftScope(scope *ScopeSnapshot, base int) *ScopeSnapshot {
	if base == 0 || scope == nil {
		return scope
	}
	out := &ScopeSnapshot{
		References:      shiftReferences(scope.References, base),
		LocalReferences: shiftReferences(scope.LocalReferences, base),
	}
	for _, level := range scope.OuterReferences {
		out.OuterReferences = append(out.OuterReferences, shiftReferences(level, base))
	}
	if scope.DMLTarget != nil {
		target := *scope.DMLTarget
		target.Loc = shiftLoc(target.Loc, base)
		target.BodyLoc = shiftLoc(target.BodyLoc, base)
		out.DMLTarget = &target
	}
	return out
}

func shiftSelectAliases(aliases []SelectAlias, base int) []SelectAlias {
	if base == 0 || len(aliases) == 0 {
		return aliases
	}
	out := make([]SelectAlias, len(aliases))
	copy(out, aliases)
	for i := range out {
		out[i].Loc = shiftLoc(out[i].Loc, base)
	}
	return out
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
	stmtTokens := sliceScopeTokens(tokens, start, end)
	_, mainStart := collectCTEs(stmtTokens)
	// The whole statement, not just its main clause: the WITH clause holds
	// query blocks the cursor may be in.
	cteMap := mapCTELookup(rangeReferenceMap(collectVisibleCTEs(stmtTokens, cursorOffset)))
	// A routine body holds statements of its own: the cursor's is the one the
	// DML dispatch and its outer references read.
	if bodyStart := cursorStatementStart(stmtTokens, cursorOffset); bodyStart > mainStart {
		mainStart = bodyStart
	}
	mainTokens := sliceScopeTokens(stmtTokens, mainStart, len(stmtTokens))
	target, dmlRefs := collectDMLScope(mainTokens, cteMap)
	var enclosing []RangeReference
	enclosingFromNested := true
	if first := statementVerb(mainTokens); first >= 0 {
		switch mainTokens[first].Type {
		case kwUPDATE, kwDELETE:
			// A subquery in an UPDATE or DELETE can name the statement's
			// own tables, so they are its outermost scope level. An INSERT
			// or REPLACE target is not readable from its SELECT, but it is
			// from the ON DUPLICATE KEY UPDATE assignment, which also ends
			// the source query.
			enclosing = dmlRefs
		case kwINSERT, kwREPLACE:
			if cursorInDuplicateKeyUpdate(mainTokens, cursorOffset) {
				enclosing = dmlRefs
				enclosingFromNested = false
			}
		}
	}
	scope := collectSelectScope(stmtTokens, cursorOffset, enclosing, enclosingFromNested)
	if scope == nil {
		// No query block at the cursor: a DML or DDL statement's own tables
		// are what its clauses name.
		refs := append(append([]RangeReference{}, dmlRefs...), collectDDLScope(stmtTokens, cursorOffset)...)
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

// cursorStatementStart returns the index of the token that starts the statement
// holding the cursor. Inside a stored routine's body that is the statement after
// the last ';', BEGIN, THEN, ELSE, DO, LOOP, or REPEAT, not the CREATE. Those
// keywords only start a statement there: a CASE expression's THEN and ELSE, and
// a transaction's BEGIN, are part of the statement around them.
func cursorStatementStart(tokens []completionScopeToken, cursorOffset int) int {
	start := 0
	procedural := false
	if verb := firstNonTerminator(tokens); verb >= 0 {
		switch tokens[verb].Type {
		case kwCREATE, kwBEGIN:
			procedural = true
		}
	}
	var caseForms []bool // CASE nesting: true for a procedural CASE statement
	for i, tok := range tokens {
		if tok.Loc >= cursorOffset {
			break
		}
		if tok.Type == ';' {
			start = i + 1
			continue
		}
		if !procedural {
			continue
		}
		// A CASE expression's branches are expressions; a procedural CASE's
		// are statement lists, and only those THENs start a statement.
		switch tok.Type {
		case kwCASE:
			caseForms = append(caseForms, startsCaseStatement(tokens, i, start))
		case kwEND:
			if i+1 < len(tokens) && tokens[i+1].Type == kwCASE {
				// END CASE closes a procedural CASE.
				if len(caseForms) > 0 {
					caseForms = caseForms[:len(caseForms)-1]
				}
				continue
			}
			if len(caseForms) > 0 {
				caseForms = caseForms[:len(caseForms)-1]
			}
		case kwBEGIN, kwTHEN, kwELSE, kwDO, kwLOOP, kwREPEAT:
			if len(caseForms) > 0 && !caseForms[len(caseForms)-1] {
				continue
			}
			start = i + 1
		}
	}
	if start >= len(tokens) {
		return 0
	}
	return start
}

// startsCaseStatement reports whether the CASE at idx begins a statement rather
// than an expression: it follows the statement start, a ';', or a routine
// body's BEGIN, THEN, or ELSE.
func startsCaseStatement(tokens []completionScopeToken, idx, start int) bool {
	if idx == start {
		return true
	}
	prev := tokens[idx-1]
	if prev.depth != tokens[idx].depth {
		return false
	}
	switch prev.Type {
	case ';', kwBEGIN, kwTHEN, kwELSE:
		return true
	}
	return false
}

// sliceScopeTokens returns tokens[start:end] as a token slice of its own, with
// the cached closing-parenthesis indices rebased: the readers take those as
// indices into the slice they were given.
func sliceScopeTokens(tokens []completionScopeToken, start, end int) []completionScopeToken {
	if start < 0 {
		start = 0
	}
	if end > len(tokens) {
		end = len(tokens)
	}
	if start >= end {
		return nil
	}
	out := make([]completionScopeToken, end-start)
	copy(out, tokens[start:end])
	for i := range out {
		c := out[i].close
		switch {
		case c < 0:
		case c < start || c >= end:
			out[i].close = -1
		default:
			out[i].close = c - start
		}
	}
	return out
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
	// Every definition the cursor statement carries, whether or not it binds
	// a name at the cursor: a caller shows them as the statement's CTEs.
	return collectWithDefinitions(sliceScopeTokens(tokens, start, end), cursorOffset, false)
}

// collectVisibleCTEs reads every WITH clause whose query block holds the
// cursor, outer blocks first, and replaces a definition an inner clause
// shadows.
func collectVisibleCTEs(tokens []completionScopeToken, cursorOffset int) []RangeReference {
	return collectWithDefinitions(tokens, cursorOffset, true)
}

// completionCTEBlock is one WITH clause of a statement: the query block it
// governs and the definitions it carries, parsed once.
type completionCTEBlock struct {
	withIdx   int
	end       int
	recursive bool
	defs      []RangeReference
}

// collectCTEBlocks parses every WITH clause of a statement once. A frame's CTE
// environment is then a filter over these blocks rather than a rescan of the
// statement per frame.
func collectCTEBlocks(tokens []completionScopeToken) []completionCTEBlock {
	ends := completionBlockEnds(tokens)
	var blocks []completionCTEBlock
	for i, tok := range tokens {
		if tok.Type != kwWITH {
			continue
		}
		defs, _ := collectCTEsAt(tokens, i, 0)
		blocks = append(blocks, completionCTEBlock{
			withIdx:   i,
			end:       ends[i],
			recursive: withIsRecursive(tokens, i),
			defs:      defs,
		})
	}
	return blocks
}

// completionBlockEnds returns, for every token, the index of the first token to
// its right with a smaller depth: the token past the group holding it, or the
// end of the statement. One pass, with a monotonic stack.
func completionBlockEnds(tokens []completionScopeToken) []int {
	ends := make([]int, len(tokens))
	stack := make([]int, 0, len(tokens))
	for i := len(tokens) - 1; i >= 0; i-- {
		for len(stack) > 0 && tokens[stack[len(stack)-1]].depth >= tokens[i].depth {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			ends[i] = len(tokens)
		} else {
			ends[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	return ends
}

// ctesVisibleAt returns the definitions of the blocks holding token idx, outer
// blocks first, an inner clause shadowing an outer one. onlyVisible also drops
// a definition that does not bind a name at the cursor yet, which is what
// resolving a table source needs: a later CTE of the same clause is not in
// scope inside an earlier definition.
func ctesVisibleAt(blocks []completionCTEBlock, tokens []completionScopeToken, idx, cursorOffset int, onlyVisible bool) []RangeReference {
	var ctes []RangeReference
	index := make(map[string]int)
	for _, block := range blocks {
		if block.withIdx > idx || idx >= block.end {
			continue
		}
		for _, def := range block.defs {
			if onlyVisible && !cteDefinitionVisible(def, cursorOffset, block.recursive) {
				continue
			}
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

// tokenIndexAtOrBefore returns the index of the last token starting at or
// before offset, or -1 when the offset is before the first token.
func tokenIndexAtOrBefore(tokens []completionScopeToken, offset int) int {
	idx := -1
	for i, tok := range tokens {
		if tok.Loc > offset {
			break
		}
		idx = i
	}
	return idx
}

// collectWithDefinitions reads the definitions of every WITH clause whose query
// block holds the cursor.
func collectWithDefinitions(tokens []completionScopeToken, cursorOffset int, onlyVisible bool) []RangeReference {
	return ctesVisibleAt(collectCTEBlocks(tokens), tokens, tokenIndexAtOrBefore(tokens, cursorOffset), cursorOffset, onlyVisible)
}

// withIsRecursive reports whether the WITH clause at withIdx is RECURSIVE.
func withIsRecursive(tokens []completionScopeToken, withIdx int) bool {
	return withIdx+1 < len(tokens) && tokens[withIdx+1].depth == tokens[withIdx].depth && tokens[withIdx+1].Type == kwRECURSIVE
}

// cteDefinitionVisible reports whether a CTE definition binds names at offset
// at: a definition is visible once it is complete, and a RECURSIVE clause may
// name the definition being read.
func cteDefinitionVisible(def RangeReference, at int, recursive bool) bool {
	if def.BodyLoc.End >= 0 && def.BodyLoc.End <= at {
		return true
	}
	if !recursive || def.BodyLoc.Start < 0 || at < def.BodyLoc.Start {
		return false
	}
	// An unterminated body runs to the end of the fragment.
	return def.BodyLoc.End < 0 || at <= def.BodyLoc.End
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
	return collectCTEsAt(tokens, first, 0)
}

// collectCTEsAt parses the WITH clause at withIdx and returns its definitions
// and the index after the clause, where the statement's main clause starts.
// nesting bounds the recursion a body's output columns cause: a body may itself
// open with a WITH clause, so reading columns re-enters this reader.
func collectCTEsAt(tokens []completionScopeToken, withIdx, nesting int) ([]RangeReference, int) {
	if nesting >= maxNestingDepth {
		return nil, withIdx
	}
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
			// The body is still being typed: the definition is there, with
			// whatever columns the text already names, and its body runs to
			// the end of the fragment.
			if len(columns) == 0 {
				columns = selectListColumns(tokens, idx+1, len(tokens), nesting+1)
			}
			ctes = append(ctes, RangeReference{
				Kind:    RangeReferenceCTE,
				Table:   name,
				Columns: columns,
				Loc:     tokenLoc(tokens[withIdx], tokens[len(tokens)-1]),
				BodyLoc: nodes.Loc{Start: tokens[idx].Loc, End: -1},
			})
			return ctes, len(tokens)
		}
		if len(columns) == 0 {
			columns = selectListColumns(tokens, idx+1, closeIdx, nesting+1)
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

// cteLookup resolves a table name to the CTE definition in scope at the cursor
// or frame it was built for. It exists so a frame's environment can be resolved
// innermost-first without materializing a map of every enclosing definition.
type cteLookup func(name string) (RangeReference, bool)

// mapCTELookup returns a lookup over a finished map.
func mapCTELookup(refs map[string]RangeReference) cteLookup {
	if len(refs) == 0 {
		return nil
	}
	return func(name string) (RangeReference, bool) {
		ref, ok := refs[name]
		return ref, ok
	}
}

// cteFrameLookup returns the lookup for the frame at token idx: the definitions
// of the blocks holding it, innermost block first, later definitions of a name
// shadowing earlier ones, and a definition whose own body holds the frame
// excluded unless its clause is recursive.
func cteFrameLookup(blocks []completionCTEBlock, idx, loc int) cteLookup {
	return func(name string) (RangeReference, bool) {
		for i := len(blocks) - 1; i >= 0; i-- {
			block := blocks[i]
			if block.withIdx > idx {
				continue
			}
			if idx >= block.end {
				continue
			}
			for j := len(block.defs) - 1; j >= 0; j-- {
				def := block.defs[j]
				if !strings.EqualFold(def.Table, name) {
					continue
				}
				if !cteDefinitionVisible(def, loc, block.recursive) {
					continue
				}
				return def, true
			}
		}
		return RangeReference{}, false
	}
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
func collectDMLScope(tokens []completionScopeToken, ctes cteLookup) (*RangeReference, []RangeReference) {
	first := statementVerb(tokens)
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
		refs := parseFromReferences(tokens, idx, collectClauseEnd(tokens, idx, tokens[first].depth), tokens[first].depth, ctes, 0)
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
		if row, ok := insertRowAlias(tokens, tokens[first].depth, refs[0]); ok {
			// ON DUPLICATE KEY UPDATE names the incoming row by its alias.
			refs = append(refs, row)
		}
		return &refs[0], refs
	}
	return nil, nil
}

// collectDDLScope returns the table an ALTER TABLE, TRUNCATE, or CREATE INDEX
// statement names. Its clauses take that table's columns.
// triggerRowScope returns the trigger's table under the NEW and OLD row
// aliases: CREATE TRIGGER name {BEFORE|AFTER} event ON table FOR EACH ROW ...
func triggerRowScope(tokens []completionScopeToken, idx, depth int) []RangeReference {
	onIdx := findKeywordAtDepth(tokens, kwON, idx, len(tokens), depth)
	if onIdx < 0 {
		return nil
	}
	ref, next, ok := parseRangeReference(tokens, onIdx+1)
	if !ok {
		return nil
	}
	loc := tokenLoc(tokens[onIdx+1], tokens[next-1])
	// INSERT provides NEW only and DELETE OLD only, as the parser's
	// refuseTriggerRowRead enforces; UPDATE provides both.
	aliases := []string{"NEW", "OLD"}
	for i := idx; i < onIdx; i++ {
		switch tokens[i].Type {
		case kwINSERT:
			aliases = []string{"NEW"}
		case kwDELETE:
			aliases = []string{"OLD"}
		case kwUPDATE:
			aliases = []string{"NEW", "OLD"}
		}
	}
	var refs []RangeReference
	for _, alias := range aliases {
		row := ref
		row.Alias = alias
		row.Loc = loc
		refs = append(refs, row)
	}
	return refs
}

func collectDDLScope(tokens []completionScopeToken, cursorOffset int) []RangeReference {
	first := statementVerb(tokens)
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
		if parent, ok := foreignKeyParentAt(tokens, first, depth, cursorOffset); ok {
			// The column list of a REFERENCES clause names the parent's
			// columns, not the altered table's.
			return []RangeReference{parent}
		}
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
		if idx < len(tokens) && tokens[idx].Type == kwTRIGGER {
			return triggerRowScope(tokens, idx+1, depth)
		}
		if idx >= len(tokens) || tokens[idx].Type != kwINDEX {
			return nil
		}
		return indexTable(tokens, idx+1, depth)
	case kwDROP:
		// DROP INDEX name ON table — the index belongs to that table, as
		// ALTER TABLE ... DROP INDEX does.
		idx := first + 1
		for idx < len(tokens) && tokens[idx].depth == depth {
			switch tokens[idx].Type {
			case kwTEMPORARY:
				idx++
				continue
			}
			break
		}
		if idx >= len(tokens) || tokens[idx].Type != kwINDEX {
			return nil
		}
		return indexTable(tokens, idx+1, depth)
	}
	return nil
}

// foreignKeyParentAt returns the table a REFERENCES clause names when the
// cursor is in that table's column list, as in
// ALTER TABLE child ADD FOREIGN KEY (a) REFERENCES parent (|).
func foreignKeyParentAt(tokens []completionScopeToken, first, depth, cursorOffset int) (RangeReference, bool) {
	for i := first + 1; i < len(tokens); i++ {
		if tokens[i].depth != depth || tokens[i].Type != kwREFERENCES {
			continue
		}
		ref, next, ok := parseRangeReference(tokens, i+1)
		if !ok {
			continue
		}
		if next >= len(tokens) || tokens[next].Type != '(' {
			continue
		}
		closeIdx := matchingParen(tokens, next)
		if closeIdx < 0 {
			// The column list is still being typed and runs to the end.
			ref.Alias = ref.Table
			ref.Loc = tokenLoc(tokens[i+1], tokens[len(tokens)-1])
			return ref, true
		}
		if cursorOffset < tokens[next].Loc || cursorOffset > tokens[closeIdx].End {
			// Another constraint's reference: the cursor is not in it.
			continue
		}
		ref.Alias = ref.Table
		ref.Loc = tokenLoc(tokens[i+1], tokens[closeIdx])
		return ref, true
	}
	return RangeReference{}, false
}

// indexTable returns the table named by the ON clause of a CREATE INDEX or
// DROP INDEX statement, whose index name starts at idx.
func indexTable(tokens []completionScopeToken, idx, depth int) []RangeReference {
	onIdx := findKeywordAtDepth(tokens, kwON, idx, len(tokens), depth)
	if onIdx < 0 || onIdx+1 >= len(tokens) {
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

// collectDeleteScope reads DELETE's target and, for the multi-table forms
// DELETE t1, t2 FROM ... and DELETE FROM t1 USING ..., the references the
// statement reads.
func collectDeleteScope(tokens []completionScopeToken, deleteIdx int, ctes cteLookup) (*RangeReference, []RangeReference, bool) {
	depth := tokens[deleteIdx].depth
	// DELETE LOW_PRIORITY | QUICK | IGNORE precede the form's own keywords.
	idx := deleteIdx + 1
	for idx < len(tokens) && tokens[idx].depth == depth {
		switch tokens[idx].Type {
		case kwLOW_PRIORITY, kwQUICK, kwIGNORE, kwHIGH_PRIORITY, kwDELAYED:
			idx++
			continue
		}
		break
	}
	fromIdx := findKeywordAtDepth(tokens, kwFROM, idx, len(tokens), depth)
	if fromIdx < 0 {
		return nil, nil, false
	}
	// DELETE t1, t2 FROM table_references: the list before FROM names what is
	// deleted, and the FROM clause is the whole table list.
	if fromIdx > idx {
		targets, _ := parseDeleteTargets(tokens, idx, depth)
		if len(targets) == 0 {
			return nil, nil, false
		}
		for i := range targets {
			targets[i].Kind = RangeReferenceDMLTarget
		}
		refs := append([]RangeReference{}, targets...)
		refs = append(refs, parseFromReferences(tokens, fromIdx+1, len(tokens), depth, ctes, 0)...)
		return &refs[0], refs, true
	}
	// DELETE FROM t1, t2 USING table_references: the list after FROM names
	// what is deleted, and the USING clause is the table list. The list ends
	// at USING, or at the clause after it in the single-table form.
	targets, next := parseDeleteTargets(tokens, fromIdx+1, depth)
	if len(targets) == 0 {
		return nil, nil, false
	}
	for i := range targets {
		targets[i].Kind = RangeReferenceDMLTarget
	}
	if next < len(tokens) && tokens[next].depth == depth && tokens[next].Type == kwUSING {
		refs := append([]RangeReference{}, targets...)
		refs = append(refs, parseFromReferences(tokens, next+1, len(tokens), depth, ctes, 0)...)
		return &refs[0], refs, true
	}
	// The single-table form takes an alias, and names it in its clauses.
	if len(targets) == 1 {
		alias, columns, after := parseAliasAndColumns(tokens, next, depth)
		if alias != "" {
			targets[0].Alias = alias
			targets[0].AliasColumns = columns
			targets[0].Loc.End = tokens[after-1].End
			next = after
		}
	}
	return &targets[0], targets, true
}

// insertRowAlias returns the row alias an INSERT declares for ON DUPLICATE KEY
// UPDATE, as in INSERT ... VALUES (...) AS incoming(x). Only a statement that
// has an ON DUPLICATE KEY UPDATE clause can name it.
func insertRowAlias(tokens []completionScopeToken, depth int, target RangeReference) (RangeReference, bool) {
	onIdx := -1
	for i := 0; i+1 < len(tokens); i++ {
		if tokens[i].depth == depth && tokens[i].Type == kwON && tokens[i+1].depth == depth && tokens[i+1].Type == kwDUPLICATE {
			onIdx = i
			break
		}
	}
	if onIdx < 0 {
		return RangeReference{}, false
	}
	for i := onIdx - 1; i >= 0; i-- {
		if tokens[i].depth != depth || tokens[i].Type != kwAS {
			continue
		}
		if i+1 >= len(tokens) || tokens[i+1].depth != depth || !isCompletionIdent(tokens[i+1]) || isAliasBoundary(tokens[i+1].Type) {
			return RangeReference{}, false
		}
		alias := tokens[i+1].Str
		// The row is the destination's row, named by the alias: a caller
		// resolves its columns from the table, not from a table called
		// "incoming".
		ref := RangeReference{
			Kind:     RangeReferenceRelation,
			Database: target.Database,
			Table:    target.Table,
			Alias:    alias,
			Loc:      tokenLoc(tokens[i+1], tokens[i+1]),
		}
		if i+2 < len(tokens) && tokens[i+2].depth == depth && tokens[i+2].Type == '(' {
			if closeIdx := matchingParen(tokens, i+2); closeIdx > 0 {
				ref.Columns = identNames(tokens, i+3, closeIdx)
				ref.AliasColumns = ref.Columns
				ref.Loc.End = tokens[closeIdx].End
			}
		}
		return ref, true
	}
	return RangeReference{}, false
}

// parseDeleteTarget reads one name of a DELETE target list, which may end in
// .* as in DELETE FROM t1.* USING ...: the star names the table, so it is not
// the dangling dot of a name still being typed.
func parseDeleteTarget(tokens []completionScopeToken, idx, depth int) (RangeReference, int, bool) {
	var parts []string
	i := idx
	for i < len(tokens) && tokens[i].depth == depth && isCompletionIdent(tokens[i]) {
		parts = append(parts, tokens[i].Str)
		if i+1 >= len(tokens) || tokens[i+1].depth != depth || tokens[i+1].Type != '.' {
			i++
			break
		}
		if i+2 < len(tokens) && tokens[i+2].depth == depth && tokens[i+2].Type == '*' {
			i += 3
			break
		}
		if i+2 >= len(tokens) || tokens[i+2].depth != depth || !isCompletionIdent(tokens[i+2]) {
			return RangeReference{}, idx, false
		}
		i += 2
	}
	if len(parts) == 0 {
		return RangeReference{}, idx, false
	}
	ref := RangeReference{Kind: RangeReferenceRelation, Table: parts[len(parts)-1]}
	if len(parts) > 1 {
		ref.Database = parts[len(parts)-2]
	}
	return ref, i, true
}

// parseDeleteTargets reads the comma-separated target list of DELETE FROM,
// each name optionally followed by .*, and returns the index after it.
func parseDeleteTargets(tokens []completionScopeToken, idx, depth int) ([]RangeReference, int) {
	var refs []RangeReference
	for idx < len(tokens) && tokens[idx].depth == depth {
		ref, next, ok := parseDeleteTarget(tokens, idx, depth)
		if !ok {
			break
		}
		ref.Loc = tokenLoc(tokens[idx], tokens[next-1])
		refs = append(refs, ref)
		idx = next
		if idx < len(tokens) && tokens[idx].depth == depth && tokens[idx].Type == ',' {
			idx++
			continue
		}
		break
	}
	return refs, idx
}

// collectSelectScope returns the scope of the query block holding the cursor,
// or nil when the cursor is in no query block. An empty block keeps its outer
// references: only a statement with no query block at all falls back to its
// DML or DDL tables. enclosing holds the references of an UPDATE or DELETE the
// cursor's block is nested in; they are visible only from a nested block.
func collectSelectScope(tokens []completionScopeToken, cursorOffset int, enclosing []RangeReference, enclosingFromNested bool) *ScopeSnapshot {
	frames := selectFramesForCursor(tokens, cursorOffset)
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
	if len(enclosing) > 0 && (innermost.depth > 0 || !enclosingFromNested) {
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
		if ref.BodyLoc.Start < 0 || cursorOffset < ref.BodyLoc.Start {
			continue
		}
		if ref.BodyLoc.End >= 0 && cursorOffset > ref.BodyLoc.End {
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

// onBelongsToFrame reports whether the ON condition at onIdx belongs to the
// FROM clause of the frame whose SELECT is at selectIdx: a query block between
// them owns it instead.
func onBelongsToFrame(tokens []completionScopeToken, selectIdx, onIdx int) bool {
	depth := tokens[onIdx].depth
	for i := selectIdx + 1; i < onIdx; i++ {
		if tokens[i].Type == kwSELECT && tokens[i].depth == depth {
			return false
		}
	}
	return true
}

// joinOperandRefs returns the references the ON condition at onIdx can name:
// the join's two operands.
func joinOperandRefs(tokens []completionScopeToken, refs []RangeReference, onIdx, depth int) []RangeReference {
	joinIdx := -1
	for i := onIdx - 1; i >= 0; i-- {
		if tokens[i].depth != depth {
			continue
		}
		if tokens[i].Type == kwJOIN || tokens[i].Type == kwSTRAIGHT_JOIN {
			joinIdx = i
			break
		}
	}
	if joinIdx < 0 {
		return referencesBefore(refs, tokens[onIdx].Loc)
	}
	// The left operand is the whole join tree of this FROM item: everything
	// since the last comma at the join's depth, not just the table before the
	// keyword.
	listStart := 0
	for i := joinIdx - 1; i >= 0; i-- {
		if tokens[i].depth != depth {
			continue
		}
		if tokens[i].Type == ',' {
			listStart = tokens[i].End
			break
		}
	}
	if depth > 0 {
		// A parenthesized group bounds the operand too: the join's left side
		// cannot reach past the '(' that opened the group holding the join.
		for i := joinIdx; i >= 0; i-- {
			if tokens[i].Type == '(' && tokens[i].close > joinIdx {
				if tokens[i].End > listStart {
					listStart = tokens[i].End
				}
				break
			}
		}
	}
	var out []RangeReference
	for _, ref := range refs {
		if ref.Loc.Start >= listStart && ref.Loc.End <= tokens[joinIdx].Loc {
			out = append(out, ref)
		}
	}
	for _, ref := range refs {
		if ref.Loc.Start >= tokens[joinIdx].Loc && ref.Loc.End <= tokens[onIdx].Loc {
			out = append(out, ref)
		}
	}
	return out
}

// openParenBefore returns the index of the '(' matching the ')' at closeIdx.
func openParenBefore(tokens []completionScopeToken, closeIdx int) int {
	for i := closeIdx - 1; i >= 0; i-- {
		if tokens[i].Type == '(' && tokens[i].close == closeIdx {
			return i
		}
	}
	return -1
}

// selectFramesForCursor returns the query blocks holding the cursor, outermost
// first. Each frame resolves its table sources against the CTE definitions
// visible at *its own* position, not the cursor's: an inner WITH clause that
// shadows a name does not change what the block around it read.
func selectFramesForCursor(tokens []completionScopeToken, cursorOffset int) []completionSelectFrame {
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
	blocks := collectCTEBlocks(tokens)
	armEnds := selectArmEnds(tokens)
	for i, tok := range tokens {
		if tok.Loc > cursorOffset {
			break
		}
		if tok.depth > cursorDepth {
			continue
		}
		tablePrimary := false
		switch tok.Type {
		case kwSELECT:
		case kwTABLE:
			// TABLE t is a query primary: its relation is the table after
			// the keyword, and its clauses read like a SELECT's.
			if !startsQueryPrimary(tokens, i) {
				continue
			}
			tablePrimary = true
		default:
			continue
		}
		depth := tok.depth
		armEnd := armEnds[i]
		if armEnd < len(tokens) && tokens[armEnd].Loc < cursorOffset {
			continue
		}
		if closedBeforeCursor(tokens, i, depth, cursorOffset) {
			// Its query block ended before the cursor: the parentheses
			// around it did.
			continue
		}
		var local []RangeReference
		frameCTEs := cteFrameLookup(blocks, i, tok.Loc)
		if tablePrimary {
			if sources, _, ok := parseTableSource(tokens, i+1, frameCTEs); ok {
				local = sources
			}
		} else if fromIdx := findKeywordAtDepth(tokens, kwFROM, i+1, armEnd, depth); fromIdx >= 0 {
			local = parseFromReferences(tokens, fromIdx+1, armEnd, depth, frameCTEs, 0)
			if onIdx := onPredicateToken(tokens, cursorOffset, cursorDepth); onIdx >= 0 && onBelongsToFrame(tokens, i, onIdx) {
				// A join condition names its own join's operands: a comma
				// before the join starts an unrelated list item, and a
				// parenthesized operand contributes all its members. The
				// condition may sit in a group the frame's FROM opened.
				local = joinOperandRefs(tokens, local, onIdx, tokens[onIdx].depth)
			}
		}
		frames = append(frames, completionSelectFrame{selectIdx: i, depth: depth, localRefs: local})
	}
	return frames
}

// onPredicateToken returns the index of the ON or USING token whose condition
// holds the cursor, or -1 when the cursor is anywhere else in the query block.
func onPredicateToken(tokens []completionScopeToken, cursorOffset, depth int) int {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Loc >= cursorOffset {
			continue
		}
		if tokens[i].depth > depth {
			continue
		}
		if tokens[i].depth < depth {
			return -1
		}
		switch tokens[i].Type {
		case kwON, kwUSING:
			return i
		case ',', ';', kwFROM, kwWHERE, kwGROUP, kwHAVING, kwORDER, kwLIMIT,
			kwINTO, kwUNION, kwINTERSECT, kwEXCEPT, kwJOIN, kwSTRAIGHT_JOIN,
			kwLEFT, kwRIGHT, kwINNER, kwCROSS, kwNATURAL:
			return -1
		}
	}
	return -1
}

// matchingBrace returns the index of the '}' closing the '{' at openIdx, or -1.
func matchingBrace(tokens []completionScopeToken, openIdx int) int {
	if openIdx < 0 || openIdx >= len(tokens) || tokens[openIdx].Type != '{' {
		return -1
	}
	depth := 0
	for i := openIdx; i < len(tokens); i++ {
		switch tokens[i].Type {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// referencesBefore keeps the references that start before offset at.
func referencesBefore(refs []RangeReference, at int) []RangeReference {
	kept := make([]RangeReference, 0, len(refs))
	for _, ref := range refs {
		if ref.Loc.Start >= 0 && ref.Loc.Start < at {
			kept = append(kept, ref)
		}
	}
	return kept
}

// startsQueryPrimary reports whether the TABLE at idx begins a query
// expression, as in TABLE t or (SELECT a) UNION TABLE t, rather than a DDL
// clause, as in CREATE TABLE t.
func startsQueryPrimary(tokens []completionScopeToken, idx int) bool {
	if idx == 0 {
		return true
	}
	prev := tokens[idx-1]
	switch prev.Type {
	case '(', ';', kwUNION, kwINTERSECT, kwEXCEPT:
		return true
	case kwALL, kwDISTINCT:
		return idx >= 2 && tokens[idx-2].Type != kwSELECT
	}
	return false
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

// selectArmEnds returns, for every token, where the query block starting there
// ends: the next token at the same depth that ends an arm. One backward pass,
// so a statement with many frames does not rescan the tail per frame.
func selectArmEnds(tokens []completionScopeToken) []int {
	ends := make([]int, len(tokens))
	maxDepth := 0
	for _, tok := range tokens {
		if tok.depth > maxDepth {
			maxDepth = tok.depth
		}
	}
	stopAt := make([]int, maxDepth+1)
	for i := range stopAt {
		stopAt[i] = len(tokens)
	}
	for i := len(tokens) - 1; i >= 0; i-- {
		depth := tokens[i].depth
		ends[i] = stopAt[depth]
		switch tokens[i].Type {
		case kwUNION, kwINTERSECT, kwEXCEPT, kwSELECT, ')':
			stopAt[depth] = i
		}
	}
	return ends
}

// selectArmEnd returns where the SELECT starting at start is no longer its own
// query block: the set operator that ends its arm, the next query block, or the
// ')' that closes the group holding it. Crossing either would read a sibling
// block's FROM clause at the same depth.
func selectArmEnd(tokens []completionScopeToken, start, depth int) int {
	for i := start + 1; i < len(tokens); i++ {
		if tokens[i].depth != depth {
			continue
		}
		switch tokens[i].Type {
		case kwUNION, kwINTERSECT, kwEXCEPT, kwSELECT, ')':
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

func parseFromReferences(tokens []completionScopeToken, start, end, depth int, ctes cteLookup, nesting int) []RangeReference {
	if nesting >= maxNestingDepth {
		return nil
	}
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
		if expectSource && tok.Type == '{' {
			// ODBC's outer-join escape: { OJ table_references }.
			if closeIdx := matchingBrace(tokens, i); closeIdx > i {
				inner := i + 1
				if inner < closeIdx && tokens[inner].Type == kwOJ {
					inner++
				}
				refs = append(refs, parseFromReferences(tokens, inner, closeIdx, depth, ctes, nesting+1)...)
				i = closeIdx + 1
				expectSource = false
				continue
			}
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
func parseTableSource(tokens []completionScopeToken, idx int, ctes cteLookup) ([]RangeReference, int, bool) {
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
	if ref.Database == "" && ctes != nil {
		if cte, ok := ctes(foldName(ref.Table)); ok {
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
			case kwJOIN:
				next += 2
			case kwGROUP:
				next += 2
				if next < len(tokens) && tokens[next].depth == depth && tokens[next].Type == kwBY {
					next++
				}
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
func parseParenthesizedTableSource(tokens []completionScopeToken, idx int, lateral bool, ctes cteLookup) ([]RangeReference, int, bool) {
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
		// Unterminated group: its body is what is being typed. A query body
		// is still a derived table, so the enclosing FROM's references stay
		// invisible to it unless it is LATERAL.
		if parenOpensQuery(tokens, idx) {
			// An unterminated body runs to the end of the text: End is -1 so
			// the visibility rule sees the cursor inside it.
			return []RangeReference{{
				Kind:    RangeReferenceSubquery,
				Loc:     tokenLoc(tokens[idx], tokens[len(tokens)-1]),
				BodyLoc: nodes.Loc{Start: tokens[idx].Loc, End: -1},
				Lateral: lateral,
			}}, len(tokens), true
		}
		return nil, idx + 1, false
	}
	depth := tokens[idx].depth
	isQuery := parenOpensQuery(tokens, idx)
	if !isQuery {
		inner := parseFromReferences(tokens, idx+1, closeIdx, depth+1, ctes, 0)
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
		Columns: selectListColumns(tokens, idx+1, closeIdx, 0),
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

// parenOpensQuery reports whether the group opened at open holds a query
// rather than a table list, looking through parentheses a query may wrap
// itself in: ((SELECT ...) UNION (SELECT ...)) x and ((SELECT ...)) x are both
// derived tables.
func parenOpensQuery(tokens []completionScopeToken, open int) bool {
	for open+1 < len(tokens) && tokens[open+1].Type == '(' {
		open++
	}
	if open+1 >= len(tokens) || tokens[open+1].depth != tokens[open].depth+1 {
		return false
	}
	switch tokens[open+1].Type {
	case kwSELECT, kwWITH, kwVALUES, kwTABLE:
		return true
	}
	return false
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
		// A table function's arguments may name the sources before it
		// without LATERAL.
		Lateral: true,
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
		return jsonTableColumnItems(tokens, i+2, columnsClose, tokens[i+1].depth+1, 0)
	}
	return nil
}

// jsonTableColumnItems reads the name each column definition in tokens[lo:hi]
// starts with, at itemDepth, and recurses into a NESTED COLUMNS clause. nesting
// bounds the recursion the way Parse bounds its own: editor input is not
// guaranteed to be valid, and a reader that kept descending would not return a
// scope either.
func jsonTableColumnItems(tokens []completionScopeToken, lo, hi, itemDepth, nesting int) []string {
	if nesting >= maxNestingDepth {
		return nil
	}
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
						names = append(names, jsonTableColumnItems(tokens, j+2, nested, tokens[j+1].depth+1, nesting+1)...)
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
func selectListColumns(tokens []completionScopeToken, lo, hi, nesting int) []string {
	if nesting >= maxNestingDepth {
		return nil
	}
	if lo < 0 || hi > len(tokens) || lo >= hi {
		return nil
	}
	// A body may open with parentheses and WITH clauses: its output columns
	// are the query inside them, not a CTE's own select list.
	for lo < hi {
		if tokens[lo].Type == '(' {
			lo++
			continue
		}
		if tokens[lo].Type != kwWITH {
			break
		}
		_, next := collectCTEsAt(tokens, lo, nesting+1)
		if next <= lo || next >= hi {
			return nil
		}
		lo = next
	}
	if lo >= hi {
		return nil
	}
	if tokens[lo].Type == kwVALUES {
		return valuesColumnNames(tokens, lo, hi)
	}
	depth := tokens[lo].depth
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

// valuesColumnNames returns the names MySQL gives a VALUES query's columns:
// column_0 to column_{n-1}, where n is the first row's arity. The row's
// expressions are not the output names, and a scalar subquery inside one is
// not the body's query.
func valuesColumnNames(tokens []completionScopeToken, lo, hi int) []string {
	depth := tokens[lo].depth
	for i := lo + 1; i < hi; i++ {
		if tokens[i].depth != depth || tokens[i].Type != kwROW {
			continue
		}
		if i+1 >= hi || tokens[i+1].Type != '(' {
			return nil
		}
		closeIdx := matchingParen(tokens, i+1)
		if closeIdx < 0 || closeIdx > hi {
			closeIdx = hi
		}
		arity := 1
		inner := tokens[i+1].depth + 1
		for j := i + 2; j < closeIdx; j++ {
			if tokens[j].depth == inner && tokens[j].Type == ',' {
				arity++
			}
		}
		names := make([]string, arity)
		for k := range names {
			names[k] = "column_" + strconv.Itoa(k)
		}
		return names
	}
	return nil
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
	if isQualifiedNameChain(tokens, start, end, depth) {
		return tokens[end-1].Str, true
	}
	return "", false
}

// isQualifiedNameChain reports whether the select item is only a dotted name,
// as in a, t.a, or db.t.a. Any other expression names nothing the parser can
// read: its last operand is not its output name.
func isQualifiedNameChain(tokens []completionScopeToken, start, end, depth int) bool {
	if end-start < 3 || (end-start)%2 == 0 {
		return false
	}
	for i := start; i < end; i++ {
		if tokens[i].depth != depth {
			return false
		}
		if (i-start)%2 == 0 {
			if !isCompletionIdent(tokens[i]) {
				return false
			}
			continue
		}
		if tokens[i].Type != '.' {
			return false
		}
	}
	return true
}

// trimSelectModifiers drops the select options that may precede the first
// select item, such as DISTINCT: they are keywords, never a column's name.
func trimSelectModifiers(tokens []completionScopeToken, start, end, depth int) (int, int) {
	for start < end && tokens[start].depth == depth {
		switch tokens[start].Type {
		case kwDISTINCT, kwDISTINCTROW, kwALL, kwHIGH_PRIORITY, kwSTRAIGHT_JOIN, kwSQL_SMALL_RESULT,
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
	frames := selectFramesForCursor(sliceScopeTokens(tokens, start, end), cursorOffset)
	selectIdx := -1
	if len(frames) > 0 {
		selectIdx = frames[len(frames)-1].selectIdx + start
		if inSetOperationGlobalClause(tokens, cursorOffset, tokens[selectIdx].depth) {
			// The ORDER BY or LIMIT after a UNION belongs to the result, whose
			// column names come from the first arm.
			if first := firstArmSelect(tokens, selectIdx, tokens[selectIdx].depth); first >= 0 {
				selectIdx = first
			}
		}
	} else {
		// A set operation whose arms are parenthesized has no open frame at a
		// global ORDER BY, and the result's names still come from its first
		// arm.
		depth := 0
		if idx := tokenIndexAtOrBefore(tokens, cursorOffset); idx >= 0 {
			depth = tokens[idx].depth
		}
		if inSetOperationGlobalClause(tokens, cursorOffset, depth) {
			selectIdx = firstSetOperationArm(tokens, start, cursorOffset, depth)
		}
	}
	if selectIdx < 0 {
		return nil
	}
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

// inSetOperationGlobalClause reports whether the cursor is in the ORDER BY or
// LIMIT of a set operation rather than of one of its arms.
func inSetOperationGlobalClause(tokens []completionScopeToken, cursorOffset, depth int) bool {
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i].Loc >= cursorOffset {
			continue
		}
		if tokens[i].depth > depth {
			continue
		}
		if tokens[i].depth < depth {
			return false
		}
		switch tokens[i].Type {
		case kwORDER, kwLIMIT:
			return hasSetOperatorBefore(tokens, i, depth)
		case kwSELECT, kwFROM, kwWHERE, kwGROUP, kwHAVING, kwWINDOW, kwUNION, kwINTERSECT, kwEXCEPT, ';':
			return false
		}
	}
	return false
}

// hasSetOperatorBefore reports whether a set operator sits between idx and the
// start of its query block.
func hasSetOperatorBefore(tokens []completionScopeToken, idx, depth int) bool {
	for i := idx - 1; i >= 0; i-- {
		if tokens[i].depth < depth || (tokens[i].depth == depth && tokens[i].Type == ';') {
			return false
		}
		if tokens[i].depth != depth {
			continue
		}
		switch tokens[i].Type {
		case kwUNION, kwINTERSECT, kwEXCEPT:
			return true
		}
	}
	return false
}

// firstArmSelect returns the SELECT that names a set operation's result
// columns: the first arm's, which is the first SELECT of the query block.
// firstSetOperationArm returns the first arm's SELECT of the set operation at
// depth whose global clause holds the cursor, whether or not the arm is
// parenthesized.
func firstSetOperationArm(tokens []completionScopeToken, start, cursorOffset, depth int) int {
	for i := start; i < len(tokens); i++ {
		if tokens[i].Loc >= cursorOffset {
			break
		}
		switch {
		case tokens[i].Type == kwSELECT:
			return i
		case tokens[i].depth == depth && (tokens[i].Type == kwUNION || tokens[i].Type == kwINTERSECT || tokens[i].Type == kwEXCEPT):
			return -1
		}
	}
	return -1
}

func firstArmSelect(tokens []completionScopeToken, sel, depth int) int {
	start := 0
	for i := sel - 1; i >= 0; i-- {
		if tokens[i].depth < depth || (tokens[i].depth == depth && tokens[i].Type == ';') {
			start = i + 1
			break
		}
	}
	for i := start; i < sel; i++ {
		if tokens[i].depth == depth && tokens[i].Type == kwSELECT {
			return i
		}
	}
	return sel
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
	prev := tokens[end-2]
	if prev.Type != kwAS {
		if !endsSelectExpression(prev) {
			return SelectAlias{}, false
		}
		// Adjacent text literals concatenate: 'a' 'b' is one string, not a
		// string aliased by 'b'.
		if last.Type == tokSCONST && prev.Type == tokSCONST {
			return SelectAlias{}, false
		}
		// INTERVAL 1 DAY ends in its unit, not in an alias the item declares.
		if intervalUnitAt(tokens, start, end, depth, last) {
			return SelectAlias{}, false
		}
	}
	return SelectAlias{Name: last.Str, Loc: tokenLoc(last, last)}, true
}

// intervalUnitAt reports whether last, the trailing identifier of a select
// item, is the unit of an INTERVAL expression. A trailing alias is an alias
// even after an INTERVAL (`INTERVAL 1 DAY next_day`), so the token has to be a
// unit keyword, not merely follow INTERVAL.
func intervalUnitAt(tokens []completionScopeToken, start, end, depth int, last completionScopeToken) bool {
	if last.Type == tokSCONST || !isIntervalUnit(last.Type) {
		return false
	}
	for i := end - 2; i >= start; i-- {
		if tokens[i].depth != depth {
			continue
		}
		switch tokens[i].Type {
		case kwAS:
			return false
		case kwINTERVAL:
			return true
		}
	}
	return false
}

// isIntervalUnit reports whether typ is one of MySQL's INTERVAL unit keywords.
func isIntervalUnit(typ int) bool {
	switch typ {
	case kwMICROSECOND, kwSECOND, kwMINUTE, kwHOUR, kwDAY, kwWEEK, kwMONTH,
		kwQUARTER, kwYEAR, kwSECOND_MICROSECOND, kwMINUTE_MICROSECOND,
		kwMINUTE_SECOND, kwHOUR_MICROSECOND, kwHOUR_SECOND, kwHOUR_MINUTE,
		kwDAY_MICROSECOND, kwDAY_SECOND, kwDAY_MINUTE, kwDAY_HOUR,
		kwYEAR_MONTH:
		return true
	}
	return false
}

// endsSelectExpression reports whether tok can end a select item's expression,
// so that the identifier after it is an implicit alias rather than the next
// operand.
func endsSelectExpression(tok completionScopeToken) bool {
	switch tok.Type {
	case tokIDENT, ')', tokSCONST, tokICONST, tokFCONST, tokXCONST, tokBCONST:
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

// statementVerb returns the index of the statement's own verb, reading through
// an EXPLAIN wrapper and its options: EXPLAIN [ANALYZE] [FORMAT = x] UPDATE ...
// is an UPDATE whose columns the SET clause names.
// cursorInDuplicateKeyUpdate reports whether the cursor is inside an INSERT's or
// REPLACE's ON DUPLICATE KEY UPDATE assignment. That clause ends the source
// query block and makes the destination table readable.
func cursorInDuplicateKeyUpdate(tokens []completionScopeToken, cursorOffset int) bool {
	depth := -1
	for i := 0; i+2 < len(tokens); i++ {
		if tokens[i].Loc >= cursorOffset {
			break
		}
		if tokens[i].Type == kwON && tokens[i+1].Type == kwDUPLICATE && tokens[i+2].Type == kwKEY {
			depth = tokens[i].depth
		}
	}
	if depth < 0 {
		return false
	}
	for _, tok := range tokens {
		if tok.Loc >= cursorOffset && tok.depth == depth && tok.Type == ';' {
			return false
		}
	}
	return true
}

func statementVerb(tokens []completionScopeToken) int {
	first := firstNonTerminator(tokens)
	if first < 0 {
		return -1
	}
	depth := tokens[first].depth
	i := first
	for i < len(tokens) && tokens[i].depth == depth && tokens[i].Type == kwEXPLAIN {
		i++
		for i < len(tokens) && tokens[i].depth == depth {
			switch tokens[i].Type {
			case kwANALYZE, kwEXTENDED, kwPARTITIONS:
				i++
				continue
			case kwFORMAT:
				i++
				if i < len(tokens) && tokens[i].depth == depth && tokens[i].Type == '=' {
					i++
				}
				if i < len(tokens) && tokens[i].depth == depth {
					i++
				}
				continue
			}
			break
		}
	}
	if i >= len(tokens) {
		return -1
	}
	return i
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
		if !IsIdentTokenType(tok.Type) {
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
	if qualifierAllowedSelectAlias(tokens, collectOffset) {
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
// make an outer WHERE an alias position, and an argument list's parentheses do
// not hide the clause the block is in.
func qualifierAllowedSelectAlias(tokens []completionScopeToken, cursorOffset int) bool {
	cursorDepth := 0
	if len(tokens) > 0 {
		cursorDepth = tokens[len(tokens)-1].depth
	}
	if frames := selectFramesForCursor(tokens, cursorOffset); len(frames) > 0 {
		cursorDepth = frames[len(frames)-1].depth
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
		case "constraint_ref":
			kinds = appendObjectKind(kinds, ObjectKindConstraint)
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
