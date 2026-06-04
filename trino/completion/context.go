package completion

import (
	"strings"

	"github.com/bytebase/omni/trino/parser"
)

// completionKind enumerates the high-level context the caret sits in. It drives
// which candidate families candidatesFor offers.
type completionKind int

const (
	// kindStatementStart: start of input, or just after ';'. Offer the
	// statement-leading keywords (SELECT, INSERT, CREATE, ...).
	kindStatementStart completionKind = iota
	// kindRelation: a table/view name is expected (after FROM, JOIN, INTO,
	// UPDATE, TABLE, the relation position of DML/DDL). Offer relations in the
	// session schema plus schema and catalog names so the user can drill down.
	kindRelation
	// kindColumnOrExpr: a column or general expression is expected (after
	// SELECT, WHERE, ON, GROUP/ORDER BY, HAVING, AND/OR, a select-list comma,
	// comparison operators). Offer in-scope columns plus expression keywords.
	kindColumnOrExpr
	// kindColumnOnly: a bare column name is expected, NOT a general expression
	// (the assignment target after UPDATE ... SET). Offer in-scope columns only;
	// an expression keyword like CASE/CAST here is a syntax error (oracle 481:
	// "UPDATE t SET CASE" => SYNTAX_ERROR).
	kindColumnOnly
	// kindDottedName: the caret follows "<word> ." — a qualified-name drill-down.
	// The qualifier list (the dotted words before the caret) is carried in
	// context.qualifier so candidatesFor can resolve it against the catalog and
	// the FROM scope.
	kindDottedName
	// kindKeyword: a fallback where only keywords are sensible.
	kindKeyword
)

// completionContext is the result of classifying the token stream up to the
// caret. qualifier holds the normalized dotted-name parts preceding the caret
// when kind == kindDottedName (e.g. ["mycat", "myschema"] for
// "mycat.myschema." ); it is nil otherwise.
type completionContext struct {
	kind      completionKind
	qualifier []string
}

// scanTokens lexes sql[:limit] into the real (non-EOF) tokens. limit is clamped
// into range. The trailing EOF token is dropped so callers index the last
// meaningful token directly.
func scanTokens(sql string, limit int) []parser.Token {
	if limit > len(sql) {
		limit = len(sql)
	}
	if limit < 0 {
		limit = 0
	}
	toks, _ := parser.Tokenize(sql[:limit])
	// Tokenize always terminates with a tokEOF token (Str == "", a zero Kind);
	// drop it so the last element is the last real token.
	if n := len(toks); n > 0 && toks[n-1].Kind == 0 {
		toks = toks[:n-1]
	}
	return toks
}

// detectContext classifies the completion context at scanLimit (the caret with
// the partial word already stripped).
func detectContext(sql string, scanLimit int) completionContext {
	toks := scanTokens(sql, scanLimit)
	n := len(toks)
	if n == 0 {
		return completionContext{kind: kindStatementStart}
	}

	last := toks[n-1]

	// Caret immediately after a dot: "<chain> ." — a qualified-name drill-down.
	if last.Kind == int('.') {
		return completionContext{kind: kindDottedName, qualifier: collectQualifier(toks[:n-1])}
	}

	// Caret just after ';' — a fresh statement begins.
	if last.Kind == int(';') {
		return completionContext{kind: kindStatementStart}
	}

	// Keyword-driven context. We look at the most recent *keyword* that selects
	// a context, scanning back over intervening identifiers/operators so that
	// "SELECT a, " still resolves to a column context and "FROM a JOIN " still
	// resolves to a relation context.
	if k, ok := keywordContext(toks); ok {
		return k
	}

	// Default: keywords. (e.g. right after a complete table name, the next
	// useful tokens are clause keywords like WHERE/JOIN/GROUP.)
	return completionContext{kind: kindKeyword}
}

// keywordContext walks the token stream and returns the context implied by the
// nearest context-selecting keyword, if any.
func keywordContext(toks []parser.Token) (completionContext, bool) {
	n := len(toks)
	last := toks[n-1]

	// A relation is expected immediately after these keywords.
	switch {
	case isKeyword(last, "from"), isKeyword(last, "join"), isKeyword(last, "into"),
		isKeyword(last, "update"), isKeyword(last, "describe"), isKeyword(last, "desc"):
		return completionContext{kind: kindRelation}, true
	// "TABLE x" as a relation appears after TABLE in TABLE/INSERT/DROP/ALTER/
	// TRUNCATE/ANALYZE contexts; offering relations there is always safe.
	case isKeyword(last, "table"):
		return completionContext{kind: kindRelation}, true
	}

	// "UPDATE ... SET" introduces an assignment target, which must be a bare
	// column — not a general expression. (USING is intentionally NOT here: after
	// a JOIN's USING keyword Trino requires a parenthesized column list, so a
	// bare column candidate would produce invalid syntax; we leave USING to the
	// keyword fallback.)
	if isKeyword(last, "set") {
		return completionContext{kind: kindColumnOnly}, true
	}

	// A column / expression is expected immediately after these.
	switch {
	case isKeyword(last, "select"), isKeyword(last, "where"), isKeyword(last, "on"),
		isKeyword(last, "having"), isKeyword(last, "and"), isKeyword(last, "or"),
		isKeyword(last, "qualify"):
		return completionContext{kind: kindColumnOrExpr}, true
	// "GROUP BY" / "ORDER BY" / "PARTITION BY" / "CLUSTER BY" / "DISTRIBUTE BY"
	// — the column position follows the BY keyword.
	case isKeyword(last, "by"):
		return completionContext{kind: kindColumnOrExpr}, true
	}

	// A comma extends whatever the enclosing list is. The two common lists are
	// the SELECT list (columns) and the FROM list (relations). Resolve by the
	// nearest preceding SELECT vs FROM/JOIN keyword.
	if last.Kind == int(',') {
		switch nearestListKeyword(toks[:n-1]) {
		case "from":
			return completionContext{kind: kindRelation}, true
		case "select":
			return completionContext{kind: kindColumnOrExpr}, true
		}
		// Unknown enclosing list (e.g. a function-argument list); columns/expr
		// are the most useful default.
		return completionContext{kind: kindColumnOrExpr}, true
	}

	// A comparison / arithmetic operator implies an expression operand follows.
	switch last.Kind {
	case int('='), int('<'), int('>'), int('+'), int('-'), int('*'), int('/'), int('%'):
		return completionContext{kind: kindColumnOrExpr}, true
	}

	return completionContext{}, false
}

// nearestListKeyword scans backward (respecting parentheses depth at the top
// level) for the nearest SELECT or FROM/JOIN keyword, returning "select",
// "from", or "". It is used to decide whether a top-level comma extends the
// select list or the from list.
func nearestListKeyword(toks []parser.Token) string {
	depth := 0
	for i := len(toks) - 1; i >= 0; i-- {
		switch toks[i].Kind {
		case int(')'), int(']'):
			depth++
			continue
		case int('('), int('['):
			if depth > 0 {
				depth--
			}
			continue
		}
		if depth != 0 {
			continue
		}
		switch {
		case isKeyword(toks[i], "from"), isKeyword(toks[i], "join"):
			return "from"
		case isKeyword(toks[i], "select"):
			return "select"
		// A clause keyword that ends the SELECT/FROM lists: stop here so a comma
		// inside e.g. GROUP BY a, b is not mis-attributed to the select list.
		case isKeyword(toks[i], "where"), isKeyword(toks[i], "group"),
			isKeyword(toks[i], "order"), isKeyword(toks[i], "having"),
			isKeyword(toks[i], "on"):
			return ""
		}
	}
	return ""
}

// collectQualifier returns the normalized dotted-name parts that form the
// qualifier chain ending at the caret. toks is the token slice WITHOUT the
// trailing '.'. It walks back over alternating name/'.' tokens:
// "a . b ." (caret) yields ["a", "b"]. A non-name, non-dot token ends the chain.
// Names are normalized via the lexer-decoded token text so a quoted qualifier
// resolves case-sensitively, matching the catalog's stored keys.
func collectQualifier(toks []parser.Token) []string {
	var parts []string
	i := len(toks) - 1
	for i >= 0 {
		if !isNameToken(toks[i]) {
			break
		}
		parts = append([]string{normalizeQualifierPart(toks[i])}, parts...)
		i--
		// Expect a separating '.' before the previous name; stop otherwise.
		if i < 0 || toks[i].Kind != int('.') {
			break
		}
		i--
	}
	return parts
}

// normalizeQualifierPart folds a name-token to the catalog's stored key form.
// For a quoted/back-quoted identifier the lexer already decoded the content
// (case preserved); for an unquoted identifier or a non-reserved keyword used
// as a name, the catalog folds to lower case. We replicate catalog.Normalize's
// rule using the token's decoded text and whether it was a quoted kind.
func normalizeQualifierPart(tok parser.Token) string {
	if isQuotedNameToken(tok) {
		// Quoted identifiers are case-sensitive; tok.Str is the decoded content.
		return tok.Str
	}
	return strings.ToLower(tok.Str)
}

// isKeyword reports whether tok is the given keyword (case-insensitive). It is
// true only for actual keyword tokens, not for a quoted identifier that happens
// to decode to the same letters.
func isKeyword(tok parser.Token, kw string) bool {
	if isQuotedNameToken(tok) {
		return false
	}
	k, ok := parser.KeywordToken(tok.Str)
	if !ok {
		return false
	}
	want, _ := parser.KeywordToken(kw)
	return k == want
}

// isNameToken reports whether tok can stand for a name in a qualified chain:
// any identifier kind, or a keyword (Trino allows non-reserved keywords as
// identifiers, and for completion a reserved word in a name slot is still a name
// the user is dotting through). String/number literals and operators are not
// names.
func isNameToken(tok parser.Token) bool {
	switch parser.TokenName(tok.Kind) {
	case "IDENTIFIER", "QUOTED_IDENTIFIER", "BACKQUOTED_IDENTIFIER", "DIGIT_IDENTIFIER":
		return true
	}
	// A keyword spelled as a word is name-eligible.
	_, ok := parser.KeywordToken(tok.Str)
	return ok && tok.Str != ""
}

// isQuotedNameToken reports whether tok is a quoted or back-quoted identifier
// (case-sensitive name), as opposed to an unquoted identifier or keyword.
func isQuotedNameToken(tok parser.Token) bool {
	switch parser.TokenName(tok.Kind) {
	case "QUOTED_IDENTIFIER", "BACKQUOTED_IDENTIFIER":
		return true
	}
	return false
}
