package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// The SHOW statements Doris accepts, one pattern per form. parseShow builds
// the AST with readers that take the tail of most forms loosely, as raw text;
// in strict mode it first holds the whole statement to these patterns, so a
// SHOW cut short or with a word dropped is a syntax error, as in Doris.
//
// Doris parses a statement with its ANTLR grammar (testdata/DorisParser.g4)
// and hands one that matches an unsupportedStatement alternative to its
// legacy CUP parser (testdata/DorisLegacyParser.g4), so such a form is the
// part both grammars accept. Each pattern cites its DorisParser.g4 line;
// show_stmt and show_param are DorisLegacyParser.g4 lines 1367 and 1386.
//
// Pattern syntax: an upper-case word is a keyword, matched by its spelling;
// a lower-case word is one of the nonterminals of showShapeNonterminal;
// 'x' is a punctuation token; [ ... ] is optional; ( a | b ) is a choice.
var dorisShowForms = compileShowForms([]string{
	// statement (DorisParser.g4:45-48), not handed to the CUP parser.
	"SHOW ( PROCEDURE | FUNCTION ) STATUS [ ( LIKE expr | WHERE expr ) ]",
	"SHOW CREATE PROCEDURE name",
	"SHOW ( FRONTEND | BACKEND ) CONFIG [ LIKE expr ] [ FROM int ]",
	// materializedViewStatement (:112), constraintStatement (:131).
	"SHOW CREATE MATERIALIZED VIEW name",
	"SHOW CONSTRAINTS FROM name",
	// supportedShowStatement (:224-229).
	"SHOW VIEW ( FROM | IN ) name [ ( FROM | IN ) ident ]",
	"SHOW INVERTED INDEX ( ANALYZER | TOKENIZER | TOKEN_FILTER | CHAR_FILTER )",
	// unsupportedShowStatement (:266-374).
	"SHOW SQL_BLOCK_RULE [ FOR ident ]",
	"SHOW ROW POLICY [ FOR ( ROLE ident | user ) ]",
	"SHOW STORAGE POLICY [ USING [ FOR text ] ]",
	"SHOW STAGES",
	"SHOW STORAGE ( VAULT | VAULTS )",
	"SHOW CREATE REPOSITORY FOR ident",
	"SHOW WHITELIST",
	"SHOW [ ( GLOBAL | SESSION | LOCAL ) ] VARIABLES [ wild ]",
	// The CUP parser's opt_db names one database (:274, :281-282 and below).
	"SHOW OPEN TABLES [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW TABLE STATUS [ ( FROM | IN ) dbname ] [ wild ]",
	"SHOW [ FULL ] TABLES [ ( FROM | IN ) dbname ] [ wild ]",
	"SHOW [ FULL ] VIEWS [ ( FROM | IN ) dbname ] [ wild ]",
	"SHOW TABLE int",
	"SHOW [ FULL ] PROCESSLIST",
	"SHOW [ ( GLOBAL | SESSION | LOCAL ) ] STATUS [ wild ]",
	"SHOW [ FULL ] TRIGGERS [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW EVENTS [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW PLUGINS",
	"SHOW [ STORAGE ] ENGINES",
	"SHOW AUTHORS",
	"SHOW [ BRIEF ] CREATE TABLE name",
	"SHOW CREATE VIEW name",
	"SHOW CREATE ( DATABASE | SCHEMA ) name",
	"SHOW CREATE CATALOG ident",
	"SHOW CREATE [ ( GLOBAL | SESSION | LOCAL ) ] FUNCTION dbname '(' [ fargs ] ')' [ ( FROM | IN ) ident ]",
	"SHOW ( DATABASES | SCHEMAS ) [ FROM ident ] [ wild ]",
	"SHOW DATABASE int",
	"SHOW DATA TYPES",
	"SHOW CATALOGS [ wild ]",
	"SHOW CATALOG ident",
	"SHOW DYNAMIC PARTITION TABLES [ ( FROM | IN ) ident ]",
	"SHOW [ FULL ] ( COLUMNS | FIELDS ) ( FROM | IN ) name [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW COLLATION [ wild ]",
	"SHOW ( CHAR SET | CHARSET ) [ wild ]",
	"SHOW PROC str",
	"SHOW COUNT '(' '*' ')' ( WARNINGS | ERRORS )",
	"SHOW ( WARNINGS | ERRORS ) [ limit ]",
	"SHOW LOAD WARNINGS ( ON str | [ ( FROM | IN ) ident ] [ wild ] [ limit ] )",
	"SHOW [ STREAM ] LOAD [ ( FROM | IN ) ident ] [ wild ] [ sort ] [ limit ]",
	"SHOW EXPORT [ ( FROM | IN ) dbname ] [ wild ] [ sort ] [ limit ]",
	"SHOW DELETE [ ( FROM | IN ) ident ]",
	"SHOW ALTER TABLE ( ROLLUP | MATERIALIZED VIEW | COLUMN ) [ ( FROM | IN ) ident ] [ wild ] [ sort ] [ limit ]",
	"SHOW DATA SKEW FROM tableref",
	// :319 against the CUP parser's two show_param DATA alternatives.
	"SHOW DATA [ ALL ] FROM name [ sort ]",
	"SHOW DATA [ ALL ] [ sort ] [ props ]",
	"SHOW [ TEMPORARY ] PARTITIONS FROM name [ wild ] [ sort ] [ limit ]",
	"SHOW PARTITION int",
	"SHOW TABLET int",
	"SHOW TABLETS BELONG ints",
	"SHOW TABLETS FROM name [ part ] [ wild ] [ sort ] [ limit ]",
	"SHOW PROPERTY [ FOR text ] [ wild ]",
	"SHOW ALL PROPERTIES [ wild ]",
	"SHOW BACKUP [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW [ BRIEF ] RESTORE [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW BROKER",
	"SHOW RESOURCES [ wild ] [ sort ] [ limit ]",
	"SHOW WORKLOAD GROUPS [ wild ]",
	"SHOW BACKENDS",
	"SHOW TRASH [ ON str ]",
	"SHOW FRONTENDS [ ident ]",
	"SHOW REPOSITORIES",
	"SHOW SNAPSHOT ON ident [ wild ]",
	"SHOW [ ALL ] GRANTS",
	"SHOW GRANTS FOR user",
	"SHOW ROLES",
	"SHOW PRIVILEGES",
	"SHOW [ FULL ] [ BUILTIN ] FUNCTIONS [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW GLOBAL [ FULL ] FUNCTIONS [ wild ]",
	"SHOW TYPE_CAST [ ( FROM | IN ) ident ]",
	"SHOW FILE [ ( FROM | IN ) ident ]",
	"SHOW ( KEY | KEYS | INDEX | INDEXES ) ( FROM | IN ) name [ ( FROM | IN ) ident ]",
	"SHOW TRANSACTION [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW QUERY PROFILE [ str ] [ limit ]",
	"SHOW LOAD PROFILE [ str ] [ limit ]",
	"SHOW CACHE HOTSPOT str",
	"SHOW ENCRYPTKEYS [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW TABLE CREATION [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW LAST INSERT",
	"SHOW CREATE MATERIALIZED VIEW ident ON name",
	"SHOW CATALOG RECYCLE BIN [ wild ]",
	"SHOW QUERY STATS [ ( FOR ident | FROM name [ ALL [ VERBOSE ] ] ) ]",
	"SHOW BUILD INDEX [ ( FROM | IN ) ident ] [ wild ] [ sort ] [ limit ]",
	"SHOW ( CLUSTERS | COMPUTE GROUPS )",
	"SHOW CONVERT_LIGHT_SCHEMA_CHANGE_PROCESS [ ( FROM | IN ) ident ]",
	"SHOW REPLICA STATUS FROM tableref [ wild ]",
	"SHOW REPLICA DISTRIBUTION FROM tableref",
	"SHOW TABLET STORAGE FORMAT [ VERBOSE ]",
	"SHOW TABLET DIAGNOSIS int",
	"SHOW COPY [ ( FROM | IN ) ident ] [ WHERE expr ] [ sort ] [ limit ]",
	"SHOW WARM UP JOB [ wild ]",
	// unsupportedLoadStatement (:391-394); the CUP parser's job_label is
	// one or two names.
	"SHOW [ ALL ] ROUTINE LOAD [ ( FOR dbname | wild ) ]",
	"SHOW ROUTINE LOAD TASK [ ( FROM | IN ) ident ] [ wild ]",
	"SHOW [ ALL ] CREATE ROUTINE LOAD FOR dbname",
	"SHOW CREATE LOAD FOR dbname",
	// unsupportedStatsStatement (:767-778).
	"SHOW TABLE STATS name [ part ] [ idlist ]",
	"SHOW TABLE STATS int",
	"SHOW INDEX STATS name ident",
	"SHOW COLUMN [ CACHED ] STATS name [ idlist ] [ part ]",
	"SHOW COLUMN HISTOGRAM name idlist",
	"SHOW [ AUTO ] ANALYZE [ name ] [ wild ]",
	"SHOW ANALYZE int [ wild ]",
	"SHOW AUTO JOBS [ name ] [ wild ]",
	"SHOW ANALYZE TASK STATUS int",
})

// checkShowGrammar reports a syntax error unless the SHOW statement at the
// current token, through to the end of the segment, matches one of forms.
// It leaves the parser where it found it.
func (p *Parser) checkShowGrammar(forms []*showShape) error {
	start := p.save()
	defer p.restore(start)
	m := showMatcher{p: p, far: p.cur}
	for _, f := range forms {
		p.restore(start)
		if m.match(f, func() bool { return p.cur.Kind == tokEOF }) {
			return nil
		}
		m.fail()
	}
	return syntaxErrorAt(m.far)
}

// showShape is a compiled pattern node.
type showShape struct {
	kind  showShapeKind
	text  string // the word, the punctuation, or the nonterminal's name
	items []*showShape
}

type showShapeKind int

const (
	showWord showShapeKind = iota
	showPunct
	showNonterminal
	showSeq
	showChoice
	showOptional
)

// showMatcher matches patterns with full backtracking and remembers the
// furthest token any attempt stopped at, where the error is reported.
type showMatcher struct {
	p   *Parser
	far Token
}

func (m *showMatcher) fail() {
	if m.p.cur.Loc.Start > m.far.Loc.Start {
		m.far = m.p.cur
	}
}

// match matches n at the current token and then calls k, backtracking into
// n's choices and options while k fails. The recursion is bounded by the
// pattern's size: lists are read by the nonterminals' own loops.
func (m *showMatcher) match(n *showShape, k func() bool) bool {
	p := m.p
	switch n.kind {
	case showWord:
		if !p.atShowWord(n.text) {
			m.fail()
			return false
		}
		p.advance()
		return k()
	case showPunct:
		if p.cur.Kind != int(n.text[0]) {
			m.fail()
			return false
		}
		p.advance()
		return k()
	case showNonterminal:
		if !p.showShapeNonterminal(n.text) {
			m.fail()
			return false
		}
		return k()
	case showSeq:
		return m.matchSeq(n.items, k)
	case showChoice:
		cp := p.save()
		for _, alt := range n.items {
			p.restore(cp)
			if m.match(alt, k) {
				return true
			}
			m.fail()
		}
		return false
	case showOptional:
		cp := p.save()
		if m.match(n.items[0], k) {
			return true
		}
		m.fail()
		p.restore(cp)
		return k()
	}
	return false
}

func (m *showMatcher) matchSeq(items []*showShape, k func() bool) bool {
	if len(items) == 0 {
		return k()
	}
	return m.match(items[0], func() bool { return m.matchSeq(items[1:], k) })
}

// atShowWord reports whether the current token is the keyword word, or an
// unquoted identifier spelled so where word is not one of omni's keywords.
func (p *Parser) atShowWord(word string) bool {
	if kind, ok := KeywordToken(word); ok {
		return p.cur.Kind == kind
	}
	return p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, word)
}

// showShapeNonterminal reads one nonterminal of a SHOW pattern and reports
// whether it matched. Where the two grammars differ it reads what both take.
func (p *Parser) showShapeNonterminal(name string) bool {
	switch name {
	case "ident":
		_, _, err := p.parseIdentifier()
		return err == nil
	case "name": // multipartIdentifier
		_, err := p.parseMultipartIdentifier()
		return err == nil
	case "dbname": // ident ('.' ident)?: the CUP parser's job_label, db.table
		if _, _, err := p.parseIdentifier(); err != nil {
			return false
		}
		if p.cur.Kind == int('.') {
			p.advance()
			_, _, err := p.parseIdentifierQualified()
			return err == nil
		}
		return true
	case "text": // identifierOrText
		_, _, err := p.parseIdentifierOrString()
		return err == nil
	case "str":
		return p.showShapeToken(tokString)
	case "int":
		return p.showShapeToken(tokInt)
	case "ints": // INTEGER_VALUE (',' INTEGER_VALUE)*
		for {
			if !p.showShapeToken(tokInt) {
				return false
			}
			if !p.showShapeToken(int(',')) {
				return true
			}
		}
	case "expr":
		_, err := p.parseExpr()
		return err == nil
	case "wild": // wildWhere (:521): LIKE STRING_LITERAL | WHERE expression
		if p.showShapeToken(kwLIKE) {
			return p.showShapeToken(tokString)
		}
		if p.showShapeToken(kwWHERE) {
			_, err := p.parseExpr()
			return err == nil
		}
		return false
	case "sort": // sortClause (:1316)
		if !p.showShapeToken(kwORDER) || !p.showShapeToken(kwBY) {
			return false
		}
		_, err := p.parseOrderByList()
		return err == nil
	case "limit": // limitClause (:1324)
		if !p.showShapeToken(kwLIMIT) || !p.showShapeToken(tokInt) {
			return false
		}
		if p.showShapeToken(kwOFFSET) || p.showShapeToken(int(',')) {
			return p.showShapeToken(tokInt)
		}
		return true
	case "props": // the CUP parser's properties: 'k' = 'v' pairs
		if !p.showShapeToken(kwPROPERTIES) || !p.showShapeToken(int('(')) {
			return false
		}
		for {
			if !p.showShapeToken(tokString) || !p.showShapeToken(int('=')) || !p.showShapeToken(tokString) {
				return false
			}
			if !p.showShapeToken(int(',')) {
				return p.showShapeToken(int(')'))
			}
		}
	case "user": // userIdentify (:1052) and the CUP parser's user_identity
		if _, _, err := p.parseIdentifierOrString(); err != nil {
			return false
		}
		if p.showShapeToken(int('@')) {
			_, _, err := p.parseIdentifierOrString()
			return err == nil
		}
		return true
	case "idlist": // identifierList (:1351)
		return p.showShapeIdentList()
	case "part": // partitionSpec (:962) and the CUP parser's partition_names
		p.showShapeToken(kwTEMPORARY)
		switch {
		case p.showShapeToken(kwPARTITION):
			if p.cur.Kind == int('(') {
				return p.showShapePartitionList()
			}
			_, _, err := p.parseIdentifier()
			return err == nil
		case p.showShapeToken(kwPARTITIONS):
			return p.showShapePartitionList()
		}
		return false
	case "tableref": // baseTableRef, as the ADMIN statements read it
		return p.parseAdminBaseTableRef(&ast.AdminStmt{}) == nil
	case "fargs": // functionArguments (:872)
		for {
			if !p.showShapeToken(tokDotDotDot) {
				if _, err := p.parseDataType(); err != nil {
					return false
				}
			}
			if !p.showShapeToken(int(',')) {
				return true
			}
		}
	}
	panic("show pattern names an unknown nonterminal: " + name)
}

// showShapePartitionList reads '(' '*' ')' or an identifierList.
func (p *Parser) showShapePartitionList() bool {
	if p.cur.Kind == int('(') && p.peekNext().Kind == int('*') {
		p.advance()
		p.advance()
		return p.showShapeToken(int(')'))
	}
	return p.showShapeIdentList()
}

func (p *Parser) showShapeIdentList() bool {
	if !p.showShapeToken(int('(')) {
		return false
	}
	for {
		if _, _, err := p.parseIdentifier(); err != nil {
			return false
		}
		if !p.showShapeToken(int(',')) {
			return p.showShapeToken(int(')'))
		}
	}
}

// showShapeToken consumes the current token if it is of kind.
func (p *Parser) showShapeToken(kind TokenKind) bool {
	if p.cur.Kind != kind {
		return false
	}
	p.advance()
	return true
}

// compileShowForms compiles patterns; a malformed one panics at init.
func compileShowForms(patterns []string) []*showShape {
	forms := make([]*showShape, len(patterns))
	for i, s := range patterns {
		c := showCompiler{src: s, toks: showPatternTokens(s)}
		n := c.seq()
		if c.pos != len(c.toks) {
			panic("malformed show pattern: " + s)
		}
		forms[i] = n
	}
	return forms
}

func showPatternTokens(s string) []string {
	return strings.Fields(s)
}

type showCompiler struct {
	src  string
	toks []string
	pos  int
}

// seq reads items up to a ']', ')' or '|' it leaves for its caller.
func (c *showCompiler) seq() *showShape {
	n := &showShape{kind: showSeq}
	for c.pos < len(c.toks) {
		t := c.toks[c.pos]
		switch t {
		case "]", ")", "|":
			return n
		case "[":
			c.pos++
			inner := c.seq()
			c.expect("]")
			n.items = append(n.items, &showShape{kind: showOptional, items: []*showShape{inner}})
		case "(":
			c.pos++
			ch := &showShape{kind: showChoice, items: []*showShape{c.seq()}}
			for c.pos < len(c.toks) && c.toks[c.pos] == "|" {
				c.pos++
				ch.items = append(ch.items, c.seq())
			}
			c.expect(")")
			n.items = append(n.items, ch)
		default:
			c.pos++
			switch {
			case len(t) == 3 && t[0] == '\'' && t[2] == '\'':
				n.items = append(n.items, &showShape{kind: showPunct, text: t[1:2]})
			case strings.ToLower(t) == t && showNonterminals[t]:
				n.items = append(n.items, &showShape{kind: showNonterminal, text: t})
			case strings.ToUpper(t) == t:
				n.items = append(n.items, &showShape{kind: showWord, text: t})
			default:
				panic("malformed show pattern: " + c.src)
			}
		}
	}
	return n
}

// showNonterminals names the nonterminals showShapeNonterminal reads.
var showNonterminals = map[string]bool{
	"ident": true, "name": true, "dbname": true, "text": true, "str": true,
	"int": true, "ints": true, "expr": true, "wild": true, "sort": true,
	"limit": true, "props": true, "user": true, "idlist": true, "part": true,
	"tableref": true, "fargs": true,
}

func (c *showCompiler) expect(t string) {
	if c.pos >= len(c.toks) || c.toks[c.pos] != t {
		panic("malformed show pattern: " + c.src)
	}
	c.pos++
}
