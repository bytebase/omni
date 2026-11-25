// Package completion completes ClickHouse SQL with ClickHouse's own grammar,
// read from omni's ClickHouse parser.
//
//   - Keywords: the parser tries each keyword at the caret and keeps the ones
//     it reads as that keyword (parser.ExpectedKeywords).
//   - Names: the tokens around the caret say what a name there names: a
//     table after FROM, JOIN, INTO, TABLE and the like, a database after USE
//     or DATABASE, a column elsewhere, and after a qualifier what the
//     qualifier holds. Tables and their aliases in scope are read from the
//     statement's tokens, so a statement still being typed has them too.
package completion

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/clickhouse/parser"
)

// CandidateType classifies a completion candidate.
type CandidateType = analysis.CandidateType

const (
	CandidateKeyword          = analysis.CandidateKeyword          // SQL keyword
	CandidateDatabase         = analysis.CandidateDatabase         // database name
	CandidateTable            = analysis.CandidateTable            // table name
	CandidateView             = analysis.CandidateView             // view name
	CandidateMaterializedView = analysis.CandidateMaterializedView // materialized view name
	CandidateColumn           = analysis.CandidateColumn           // column name
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string        // the completion text
	Type       CandidateType // what kind of object this is
	Definition string        // optional definition
	Comment    string        // optional doc comment
}

// completionBudget caps the time keyword probing takes when ctx allows
// longer, so a long statement never blocks a keystroke; probing returns what
// it found by then. Tests raise it, so their answers do not depend on the
// machine's load.
var completionBudget = 200 * time.Millisecond

// caretWord is the word being typed at the caret: the bytes it spans and
// the text typed before the caret, unquoted.
type caretWord struct {
	start, end int
	typed      string
	quote      byte // the open quote the caret is inside, 0 when none
}

// complete returns the candidates for the statement stmt at byte offset
// caret.
func complete(ctx context.Context, stmt string, caret int, cat *catalog) []Candidate {
	caret = max(0, min(caret, len(stmt)))
	word, ok := caretWordAt(stmt, caret)
	if !ok {
		return nil
	}
	prefix := stmt[:word.start]
	before := significantTokens(prefix)
	qualifier := qualifierBefore(before)

	var out []Candidate
	if word.quote == 0 && len(qualifier) == 0 {
		ctx, cancel := context.WithTimeout(ctx, completionBudget)
		defer cancel()
		for _, kw := range parser.ExpectedKeywords(ctx, prefix) {
			out = append(out, Candidate{Text: kw, Type: CandidateKeyword})
		}
	}
	// The tables of the statement without the caret word, innermost query
	// first. Scope comes from the tokens, not a placeholder parse: the
	// statement is usually still being typed, and the ClickHouse tree has no
	// walker that would read its scopes.
	refs := refsInScope(significantTokens(stmt[:word.start]+" "+stmt[word.end:]), word.start)
	for _, c := range nameCandidates(before, qualifier, refs, cat) {
		if word.quote != 0 {
			q := string(word.quote)
			c.Text = strings.ReplaceAll(strings.ReplaceAll(c.Text, `\`, `\\`), q, q+q)
		} else {
			c.Text = QuoteIdentifierIfNeeded(c.Text)
		}
		out = append(out, c)
	}
	return filterByPrefix(out, word.typed)
}

// nameCandidates returns the names that fit after the tokens before, the
// last of which may be a qualifier.
func nameCandidates(before []parser.Token, qualifier []string, refs []tableRef, cat *catalog) []Candidate {
	if cat == nil {
		cat = &catalog{}
	}
	var out []Candidate
	switch len(qualifier) {
	case 0:
	case 1:
		// q. where a table is named is a database. Elsewhere it is a table
		// alias, a table in scope or of the current database, or a database.
		if nameKindAfter(before[:len(before)-2]) == kindRelation {
			return relationCandidates(cat.database(qualifier[0]))
		}
		for _, ref := range refs {
			if ref.cte && ref.addressedBy(qualifier[0]) {
				return nil
			}
			if rel := cat.relation(ref.database, ref.table); rel != nil && ref.addressedBy(qualifier[0]) {
				return columnCandidates(rel)
			}
		}
		if rel := cat.relation("", qualifier[0]); rel != nil {
			out = append(out, columnCandidates(rel)...)
		}
		return append(out, relationCandidates(cat.database(qualifier[0]))...)
	default:
		// db.table.
		return columnCandidates(cat.relation(qualifier[len(qualifier)-2], qualifier[len(qualifier)-1]))
	}

	switch nameKindAfter(before) {
	case kindRelation:
		out = append(out, relationCandidates(cat.database(""))...)
		return append(out, databaseCandidates(cat)...)
	case kindDatabase:
		return databaseCandidates(cat)
	}
	for _, ref := range refs {
		// A WITH query's columns are its select list's, which no metadata
		// holds; a table of the same name does not stand in for it.
		if !ref.cte {
			out = append(out, columnCandidates(cat.relation(ref.database, ref.table))...)
		}
	}
	return out
}

func columnCandidates(rel *relation) []Candidate {
	if rel == nil {
		return nil
	}
	out := make([]Candidate, 0, len(rel.columns))
	for _, c := range rel.columns {
		out = append(out, Candidate{Text: c.name, Type: CandidateColumn, Definition: rel.name + "." + c.name + " | " + c.typ, Comment: c.comment})
	}
	return out
}

func relationCandidates(db *database) []Candidate {
	if db == nil {
		return nil
	}
	out := make([]Candidate, 0, len(db.relations))
	for _, r := range db.relations {
		out = append(out, Candidate{Text: r.name, Type: r.kind, Definition: db.name + "." + r.name, Comment: r.comment})
	}
	return out
}

func databaseCandidates(cat *catalog) []Candidate {
	var out []Candidate
	for _, name := range cat.databaseNames() {
		out = append(out, Candidate{Text: name, Type: CandidateDatabase})
	}
	return out
}

type nameKind int

const (
	kindColumn nameKind = iota
	kindRelation
	kindDatabase
)

// nameKindAfter says what an unqualified name names after the tokens before.
func nameKindAfter(before []parser.Token) nameKind {
	if len(before) == 0 {
		return kindColumn
	}
	last := before[len(before)-1]
	switch last.Type {
	case parser.KwFrom:
		if relationFrom(before, len(before)-1) {
			return kindRelation
		}
	case parser.KwUpdate:
		// An UPDATE statement names a table; ALTER TABLE t UPDATE assigns
		// columns.
		if len(before) == 1 {
			return kindRelation
		}
	case parser.KwJoin:
		// ARRAY JOIN reads columns.
		if len(before) < 2 || before[len(before)-2].Type != parser.KwArray {
			return kindRelation
		}
	case parser.KwInto, parser.KwTable,
		parser.KwDescribe, parser.KwDesc, parser.KwTruncate, parser.KwExchange:
		return kindRelation
	case parser.KwUse, parser.KwDatabase:
		return kindDatabase
	case parser.TokComma:
		// FROM a, b lists tables.
		if inFromList(before[:len(before)-1]) {
			return kindRelation
		}
	}
	return kindColumn
}

// inFromList reports whether the tokens end with a FROM list's relation:
// FROM, then names, dots, aliases, table modifiers, commas and parenthesized
// table function arguments or subqueries only. A clause keyword such as
// ORDER or BY, which can be no alias, ends the list.
func inFromList(toks []parser.Token) bool {
	for i := len(toks) - 1; i >= 0; i-- {
		switch {
		case toks[i].Type == parser.KwFrom:
			return relationFrom(toks, i)
		case toks[i].Type == parser.TokRParen:
			if i = matchingParen(toks, i); i < 0 {
				return false
			}
		case parser.IsBareAlias(toks[i]), isRatioToken(toks[i]):
		case toks[i].Type == parser.KwAs, toks[i].Type == parser.TokDot, toks[i].Type == parser.TokComma:
		case toks[i].Type == parser.KwFinal, toks[i].Type == parser.KwSample, toks[i].Type == parser.KwOffset:
		default:
			return false
		}
	}
	return false
}

// qualifierBefore returns the names of the qualifier the tokens end with,
// db. or db.table., or nil when they end with no dot.
func qualifierBefore(toks []parser.Token) []string {
	var names []string
	for i := len(toks) - 1; i >= 1 && toks[i].Type == parser.TokDot && isName(toks[i-1]); i -= 2 {
		names = append([]string{toks[i-1].Val}, names...)
	}
	return names
}

// isName reports whether tok can be a name: ClickHouse reads most keywords
// as identifiers.
func isName(tok parser.Token) bool {
	return parser.IsIdentifier(tok)
}

// relationFrom reports whether the FROM at toks[i] names tables: it is not
// inside parentheses, or a SELECT comes before it inside them. The FROM of
// EXTRACT(DAY FROM d) or TRIM(BOTH 'x' FROM s) is part of an expression.
func relationFrom(toks []parser.Token, i int) bool {
	open := openParen(toks, i)
	if open < 0 {
		return true
	}
	depth := 0
	for j := open + 1; j < i; j++ {
		switch toks[j].Type {
		case parser.TokLParen:
			depth++
		case parser.TokRParen:
			depth--
		case parser.KwSelect:
			if depth == 0 {
				return true
			}
		}
	}
	return false
}

// matchingParen returns the index of the parenthesis toks[i] closes, or -1.
func matchingParen(toks []parser.Token, i int) int {
	depth := 0
	for j := i; j >= 0; j-- {
		switch toks[j].Type {
		case parser.TokRParen:
			depth++
		case parser.TokLParen:
			if depth--; depth == 0 {
				return j
			}
		}
	}
	return -1
}

// openParen returns the index of the innermost parenthesis open before
// toks[i], or -1 when none is.
func openParen(toks []parser.Token, i int) int {
	depth := 0
	for j := i - 1; j >= 0; j-- {
		switch toks[j].Type {
		case parser.TokRParen:
			depth++
		case parser.TokLParen:
			if depth == 0 {
				return j
			}
			depth--
		}
	}
	return -1
}

// tableRef is a table a statement reads or writes, with its alias, the
// query parenthesis its clause is inside (the index of the open
// parenthesis, -1 at the top level) and the UNION branch it is in there.
type tableRef struct {
	database, table, alias string
	group, branch          int
	cte                    bool // the name of a WITH query, whose columns are not synced
}

// queryParen returns the index of the innermost parenthesis open before
// toks[i] that holds a query, or -1 when none does. Other parentheses, a
// function's arguments or a parenthesized JOIN, open no scope.
func queryParen(toks []parser.Token, i int) int {
	g := openParen(toks, i)
	for g >= 0 && !startsQuery(toks[g+1:]) {
		g = openParen(toks, g)
	}
	return g
}

// branchAt returns which set-operation branch toks[i] is in, counting the
// UNION, INTERSECT and EXCEPT between queries inside its query parenthesis.
func branchAt(toks []parser.Token, i int) int {
	depth, branch := 0, 0
	for j := queryParen(toks, i) + 1; j < i; j++ {
		switch toks[j].Type {
		case parser.TokLParen:
			depth++
		case parser.TokRParen:
			depth--
		case parser.KwUnion, parser.KwIntersect, parser.KwExcept:
			// Each is also a name, and EXCEPT drops columns in SELECT *
			// EXCEPT (a): only one a query follows separates queries.
			if depth == 0 && startsQuery(toks[j+1:]) {
				branch++
			}
		}
	}
	return branch
}

// startsQuery reports whether toks start a set operation's next query.
func startsQuery(toks []parser.Token) bool {
	for len(toks) > 0 && (toks[0].Type == parser.KwAll || toks[0].Type == parser.KwDistinct || toks[0].Type == parser.TokLParen) {
		toks = toks[1:]
	}
	return len(toks) > 0 && (toks[0].Type == parser.KwSelect || toks[0].Type == parser.KwWith)
}

// refsInScope returns the tables in scope at byte offset caret: those whose
// clause is inside a query parenthesis that holds the caret, in the UNION
// branch that holds it, the innermost first.
func refsInScope(toks []parser.Token, caret int) []tableRef {
	at := 0
	for at < len(toks) && toks[at].Pos < caret {
		at++
	}
	// Each scope is a parenthesis holding the caret and the branch of it
	// the caret is in, measured where the caret's path leaves it.
	type scope struct{ group, branch int }
	var scopes []scope
	for {
		g := queryParen(toks, at)
		scopes = append(scopes, scope{g, branchAt(toks, at)})
		if g < 0 {
			break
		}
		at = g
	}
	all := tableRefs(toks)
	var out []tableRef
	for _, sc := range scopes {
		for _, ref := range all {
			if ref.group == sc.group && ref.branch == sc.branch {
				out = append(out, ref)
			}
		}
	}
	return out
}

// addressedBy reports whether a qualifier names the table: its alias, or its
// name when it has none.
func (r tableRef) addressedBy(q string) bool {
	if r.alias != "" {
		return r.alias == q
	}
	return r.table == q
}

// tableRefs reads the tables named after FROM, JOIN, INTO, TABLE and UPDATE,
// and in a FROM list, with their aliases and the parenthesis each is in. A
// name a WITH clause gives a query is marked as such.
func tableRefs(toks []parser.Token) []tableRef {
	ctes := cteNames(toks)
	var refs []tableRef
	for i := 0; i < len(toks); i++ {
		switch toks[i].Type {
		case parser.KwFrom:
			if !relationFrom(toks, i) {
				continue
			}
		case parser.TokComma:
			if !inFromList(toks[:i]) {
				continue
			}
		case parser.KwUpdate:
			if i != 0 {
				continue // ALTER TABLE t UPDATE
			}
		case parser.KwJoin:
			if i > 0 && toks[i-1].Type == parser.KwArray {
				continue // ARRAY JOIN reads columns
			}
		case parser.KwInto, parser.KwTable:
		default:
			continue
		}
		group, branch := queryParen(toks, i), branchAt(toks, i)
		// FROM (a JOIN b ...) groups a join without opening a scope.
		j := i + 1
		for j < len(toks) && toks[j].Type == parser.TokLParen && !startsQuery(toks[j:]) {
			j++
		}
		if ref, n := readTableRef(toks[j:]); n > 0 {
			ref.group, ref.branch = group, branch
			ref.cte = ref.database == "" && ctes[ref.table]
			refs = append(refs, ref)
		}
	}
	return refs
}

// cteNames returns the names WITH gives queries: name AS (SELECT ...).
func cteNames(toks []parser.Token) map[string]bool {
	names := map[string]bool{}
	for i := 0; i+2 < len(toks); i++ {
		if isName(toks[i]) && toks[i+1].Type == parser.KwAs && toks[i+2].Type == parser.TokLParen && startsQuery(toks[i+2:]) {
			names[toks[i].Val] = true
		}
	}
	return names
}

// readTableRef reads [db.]table [[AS] alias] from the start of toks and
// returns it with the number of tokens it spans, 0 when toks starts with no
// table name.
func readTableRef(toks []parser.Token) (tableRef, int) {
	if len(toks) == 0 || !isName(toks[0]) {
		return tableRef{}, 0
	}
	ref, n := tableRef{table: toks[0].Val}, 1
	if len(toks) > 2 && toks[1].Type == parser.TokDot && isName(toks[2]) {
		ref.database, ref.table, n = toks[0].Val, toks[2].Val, 3
	}
	if n < len(toks) && toks[n].Type == parser.KwAs && n+1 < len(toks) && isName(toks[n+1]) {
		ref.alias, n = toks[n+1].Val, n+2
	} else if n < len(toks) && parser.IsBareAlias(toks[n]) {
		ref.alias, n = toks[n].Val, n+1
	}
	return ref, n + tableModifiers(toks[n:])
}

// tableModifiers returns how many tokens at the start of toks are a
// table's FINAL and SAMPLE k [OFFSET m] modifiers.
func tableModifiers(toks []parser.Token) int {
	n := 0
	for n < len(toks) {
		switch toks[n].Type {
		case parser.KwFinal:
			n++
		case parser.KwSample, parser.KwOffset:
			n++
			for n < len(toks) && isRatioToken(toks[n]) {
				n++
			}
		default:
			return n
		}
	}
	return n
}

// isRatioToken reports whether tok can be part of a SAMPLE ratio: 0.1, 1/10.
func isRatioToken(tok parser.Token) bool {
	return tok.Type == parser.TokInt || tok.Type == parser.TokFloat || tok.Type == parser.TokSlash
}

// significantTokens lexes text up to its first lexical error, without its
// EOF.
func significantTokens(text string) []parser.Token {
	lex := parser.NewLexer(text)
	var out []parser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == parser.TokEOF || tok.Type == parser.TokIllegal {
			return out
		}
		out = append(out, tok)
	}
}

// caretWordAt finds the word at caret. It fails when the caret is inside a
// string or a comment, where nothing completes.
func caretWordAt(stmt string, caret int) (caretWord, bool) {
	prefix := stmt[:caret]
	lex := parser.NewLexer(prefix)
	var last parser.Token
	count := 0
	for {
		tok := lex.NextToken()
		if tok.Type == parser.TokEOF {
			break
		}
		last = tok
		count++
		if tok.Type == parser.TokIllegal {
			break
		}
	}
	if count > 0 && last.End == caret {
		switch {
		case last.Type == parser.TokIdent || last.IsKeyword():
			return caretWord{start: last.Pos, end: wordEnd(stmt, caret), typed: prefix[last.Pos:]}, true
		case last.Type == parser.TokIllegal && (prefix[last.Pos] == '`' || prefix[last.Pos] == '"'):
			// An open quoted identifier.
			return caretWord{start: last.Pos, end: quotedEnd(stmt, caret, prefix[last.Pos]), typed: prefix[last.Pos+1:], quote: prefix[last.Pos]}, true
		case last.Type == parser.TokIllegal:
			// An open string or comment.
			return caretWord{}, false
		}
	}
	// A word typed inside a comment adds no token.
	if len(significantTokens(prefix+" x")) != len(significantTokens(prefix))+1 {
		return caretWord{}, false
	}
	return caretWord{start: caret, end: wordEnd(stmt, caret)}, true
}

// wordEnd returns the end of the unquoted word continuing at i.
func wordEnd(stmt string, i int) int {
	for i < len(stmt) && isIdentPart(stmt[i]) {
		i++
	}
	return i
}

func isIdentPart(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_'
}

// quotedEnd returns the end of the quoted identifier the caret is inside:
// just past its closing quote, or the end of its line.
func quotedEnd(stmt string, caret int, quote byte) int {
	for i := caret; i < len(stmt); i++ {
		switch stmt[i] {
		case '\n':
			return i
		case '\\':
			i++
		case quote:
			return i + 1
		}
	}
	return len(stmt)
}

// filterByPrefix keeps the candidates whose text starts with typed, ignoring
// case and a leading quote.
func filterByPrefix(candidates []Candidate, typed string) []Candidate {
	if typed == "" {
		return candidates
	}
	upper := strings.ToUpper(typed)
	return slices.DeleteFunc(candidates, func(c Candidate) bool {
		return !strings.HasPrefix(strings.ToUpper(strings.TrimLeft(c.Text, "`")), upper)
	})
}
