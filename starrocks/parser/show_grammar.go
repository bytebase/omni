package parser

import (
	"strings"
)

// The SHOW statements StarRocks accepts (testdata/StarRocks.g4), one pattern
// per form, each citing its line. parseShow builds the AST with readers that
// take the tail of most forms loosely, as raw text; in strict mode it first
// holds the whole statement to these patterns, so a SHOW cut short or with a
// word dropped is a syntax error, as in StarRocks.
//
// Pattern syntax: an upper-case word is a keyword, matched by its spelling;
// a lower-case word is one of the nonterminals of showShapeNonterminal;
// 'x' is a punctuation token; [ ... ] is optional; ( a | b ) is a choice.
var starRocksShowForms = compileShowForms([]string{
	"SHOW DATABASES [ ( FROM | IN ) name ] [ wild ]",                                                           // :344
	"SHOW SCHEMAS [ wild ]",                                                                                    // :345
	"SHOW CREATE ( DATABASE | SCHEMA ) ident",                                                                  // :362
	"SHOW DATA [ FROM name ]",                                                                                  // :374
	"SHOW [ FULL ] TABLES [ ( FROM | IN ) name ] [ wild ]",                                                     // :530
	"SHOW TEMPORARY TABLES [ ( FROM | IN ) name ] [ wild ]",                                                    // :534
	"SHOW CREATE ( TABLE | VIEW | MATERIALIZED VIEW ) name",                                                    // :538
	"SHOW [ FULL ] ( COLUMNS | FIELDS ) ( FROM | IN ) name [ ( FROM | IN ) name ] [ wild ]",                    // :542
	"SHOW TABLE STATUS [ ( FROM | IN ) name ] [ wild ]",                                                        // :547
	"SHOW ALTER TABLE ( COLUMN | ROLLUP | OPTIMIZE ) [ ( FROM | IN ) name ] [ WHERE expr ] [ sort ] [ limit ]", // :555
	"SHOW ALTER MATERIALIZED VIEW [ ( FROM | IN ) name ] [ WHERE expr ] [ sort ] [ limit ]",                    // :557
	"SHOW ( INDEX | INDEXES | KEY | KEYS ) ( FROM | IN ) name [ ( FROM | IN ) name ]",                          // :574
	"SHOW [ TEMPORARY ] PARTITIONS FROM name [ WHERE expr ] [ sort ] [ limit ]",                                // :591
	"SHOW MATERIALIZED VIEWS [ ( FROM | IN ) name ] [ wild ]",                                                  // :668
	"SHOW COMPUTE NODES",                              // :755
	"SHOW CREATE CATALOG text",                        // :765
	"SHOW CATALOGS",                                   // :773
	"SHOW STORAGE VOLUMES [ LIKE str ]",               // :796
	"SHOW FAILPOINTS [ LIKE str ] [ ON BACKEND str ]", // :837
	"SHOW DICTIONARY [ name ]",                        // :857
	"SHOW [ ALL ] ROUTINE LOAD [ FOR name ] [ FROM name ] [ WHERE expr ] [ sort ] [ limit ]", // :1275
	"SHOW ROUTINE LOAD TASK [ FROM name ] WHERE expr",                                        // :1281
	"SHOW CREATE ROUTINE LOAD name",                                                          // :1287
	"SHOW [ ALL ] STREAM LOAD [ FOR name ] [ FROM name ] [ WHERE expr ] [ sort ] [ limit ]",  // :1291
	"SHOW ANALYZE [ ( JOB | STATUS ) ] [ WHERE expr ]",                                       // :1334
	"SHOW STATS META [ WHERE expr ]",                                                         // :1338
	"SHOW HISTOGRAM META [ WHERE expr ]",                                                     // :1342
	"SHOW [ VERBOSE ] RESOURCE GROUP ident",                                                  // :1375
	"SHOW [ VERBOSE ] RESOURCE GROUPS [ ALL ]",                                               // :1376
	"SHOW USAGE RESOURCE GROUP ident",                                                        // :1380
	"SHOW USAGE RESOURCE GROUPS",                                                             // :1381
	"SHOW RESOURCES",                                                                         // :1397
	"SHOW [ FULL ] [ ( BUILTIN | GLOBAL ) ] FUNCTIONS [ ( FROM | IN ) name ] [ wild ]",       // :1407
	"SHOW LOAD [ ALL ] [ FROM ident ] [ WHERE expr ] [ sort ] [ limit ]",                     // :1487
	"SHOW LOAD WARNINGS [ FROM ident ] [ WHERE expr ] [ limit ]",                             // :1491
	"SHOW LOAD WARNINGS ON str",                                                              // :1492
	"SHOW AUTHORS",                                                                           // :1513
	"SHOW BACKENDS",                                                                          // :1517
	"SHOW BROKER",                                                                            // :1521
	"SHOW ( CHAR SET | CHARSET | CHARACTER SET ) [ wild ]",                                   // :1525
	"SHOW COLLATION [ wild ]",                                                                // :1529
	"SHOW DELETE [ ( FROM | IN ) name ]",                                                     // :1533
	"SHOW DYNAMIC PARTITION TABLES [ ( FROM | IN ) name ]",                                   // :1537
	"SHOW EVENTS [ ( FROM | IN ) name ] [ wild ]",                                            // :1541
	"SHOW ENGINES",      // :1545
	"SHOW FRONTENDS",    // :1549
	"SHOW PLUGINS",      // :1553
	"SHOW REPOSITORIES", // :1557
	"SHOW OPEN TABLES",  // :1561
	"SHOW PRIVILEGES",   // :1564
	"SHOW ( PROCEDURE | FUNCTION ) STATUS [ wild ]", // :1568
	"SHOW PROC str",                                                   // :1572
	"SHOW [ FULL ] PROCESSLIST [ FOR str ]",                           // :1576
	"SHOW PROFILELIST [ LIMIT int ]",                                  // :1580
	"SHOW RUNNING QUERIES [ LIMIT int ]",                              // :1584
	"SHOW [ ( GLOBAL | LOCAL | SESSION | VERBOSE ) ] STATUS [ wild ]", // :1588, varType :2727
	"SHOW TABLET int",                                                 // :1592
	"SHOW ( TABLET | TABLETS ) FROM name [ part ] [ WHERE expr ] [ sort ] [ limit ]", // :1593
	"SHOW TRANSACTION [ ( FROM | IN ) name ] [ WHERE expr ]",                         // :1597
	"SHOW [ FULL ] TRIGGERS [ ( FROM | IN ) name ] [ wild ]",                         // :1601
	"SHOW ( PROPERTY | PROPERTIES ) [ FOR str ] [ LIKE str ]",                        // :1605
	"SHOW [ ( GLOBAL | LOCAL | SESSION | VERBOSE ) ] VARIABLES [ wild ]",             // :1610
	"SHOW ( WARNINGS | ERRORS ) [ limit ]",                                           // :1614
	"SHOW ( USER | USERS )",                                                          // :1638
	"SHOW ALL AUTHENTICATION",                                                        // :1642
	"SHOW AUTHENTICATION [ FOR user ]",                                               // :1643
	"SHOW ROLES",                                                                     // :1663
	"SHOW GRANTS [ FOR ( ROLE text | [ USER ] user ) ]",                              // :1717-1719
	"SHOW BACKUP [ ( FROM | IN ) ident ]",                                            // :1777
	"SHOW RESTORE [ ( FROM | IN ) ident ] [ WHERE expr ]",                            // :1794
	"SHOW SNAPSHOT ON ident [ WHERE expr ]",                                          // :1798
	"SHOW SQLBLACKLIST",                                                              // :1824
	"SHOW WHITELIST",                                                                 // :1828
	"SHOW BACKEND BLACKLIST",                                                         // :1842
	"SHOW DATACACHE RULES",                                                           // :1856
	"SHOW EXPORT [ ( FROM | IN ) name ] [ wild ] [ sort ] [ limit ]",                 // :1882
	"SHOW FILE [ ( FROM | IN ) name ]",                                               // :1908
	"SHOW PIPES [ ( wild | FROM name ) ] [ sort ] [ limit ]",                         // :1940
	"SHOW PLAN ADVISOR",                                                              // :2034
	"SHOW WAREHOUSES [ LIKE str ]",                                                   // :2060
	"SHOW CLUSTERS FROM WAREHOUSE ident",                                             // :2064
	"SHOW NODES FROM WAREHOUSES [ LIKE str ]",                                        // :2068
	"SHOW NODES FROM WAREHOUSE ident",                                                // :2069
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
	p.cur = m.far
	return p.syntaxErrorAtCur()
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
// whether it matched.
func (p *Parser) showShapeNonterminal(name string) bool {
	switch name {
	case "ident":
		return p.showShapeIdent()
	case "name": // qualifiedName (:2850)
		if !p.showShapeIdent() {
			return false
		}
		for p.cur.Kind == int('.') && isIdentifierTokenQualified(p.peekNext().Kind) {
			p.advance()
			p.advance()
		}
		return true
	case "text": // identifierOrString (:2877)
		return p.showShapeToken(tokString) || p.showShapeIdent()
	case "str":
		return p.showShapeToken(tokString)
	case "int":
		return p.showShapeToken(tokInt)
	case "expr":
		_, err := p.parseExpr()
		return err == nil
	case "wild": // (LIKE pattern=string) | (WHERE expression)
		if p.showShapeToken(kwLIKE) {
			return p.showShapeToken(tokString)
		}
		if p.showShapeToken(kwWHERE) {
			_, err := p.parseExpr()
			return err == nil
		}
		return false
	case "sort": // ORDER BY sortItem (',' sortItem)*
		if !p.showShapeToken(kwORDER) || !p.showShapeToken(kwBY) {
			return false
		}
		_, err := p.parseOrderByList()
		return err == nil
	case "limit": // limitElement (:2135)
		if !p.showShapeToken(kwLIMIT) || !p.showShapeLimitValue() {
			return false
		}
		if p.showShapeToken(kwOFFSET) || p.showShapeToken(int(',')) {
			return p.showShapeLimitValue()
		}
		return true
	case "user": // user (:2892)
		if !p.showShapeNonterminal("text") {
			return false
		}
		if !p.showShapeToken(int('@')) {
			return true
		}
		if p.showShapeToken(int('[')) {
			if !p.showShapeNonterminal("text") {
				return false
			}
			return p.showShapeToken(int(']'))
		}
		return p.showShapeNonterminal("text")
	case "part": // partitionNames (:2268)
		temporary := p.showShapeToken(kwTEMPORARY)
		switch {
		case p.cur.Kind == kwPARTITION && !temporary && p.peekNext().Kind == int('('):
			// keyPartitions (:2274) or a parenthesized identifierOrString list.
			p.advance()
			p.advance()
			for {
				if !p.showShapeNonterminal("text") {
					return false
				}
				if p.showShapeToken(int('=')) {
					if _, err := p.parseExpr(); err != nil {
						return false
					}
				}
				if !p.showShapeToken(int(',')) {
					return p.showShapeToken(int(')'))
				}
			}
		case p.showShapeToken(kwPARTITION), p.showShapeToken(kwPARTITIONS):
			if !p.showShapeToken(int('(')) {
				return p.showShapeNonterminal("text")
			}
			for {
				if !p.showShapeNonterminal("text") {
					return false
				}
				if !p.showShapeToken(int(',')) {
					return p.showShapeToken(int(')'))
				}
			}
		}
		return false
	}
	panic("show pattern names an unknown nonterminal: " + name)
}

// showShapeIdent reads an identifier: one omni's lexer allows, or a word
// StarRocks.g4 lists in nonReserved that omni reserves elsewhere.
func (p *Parser) showShapeIdent() bool {
	if starRocksNonReservedWords[p.cur.Kind] {
		p.advance()
		return true
	}
	_, _, err := p.parseIdentifier()
	return err == nil
}

// starRocksNonReservedWords are the words of StarRocks.g4's nonReserved that
// IsReserved reports reserved.
var starRocksNonReservedWords = map[TokenKind]bool{
	kwADMIN: true, kwANTI: true, kwBACKEND: true, kwBINARY: true, kwCANCEL: true,
	kwCAST: true, kwCLEAN: true, kwCURRENT: true, kwDECOMMISSION: true, kwDISK: true,
	kwDISTRIBUTED: true, kwDISTRIBUTION: true, kwDUPLICATE: true, kwEXECUTE: true,
	kwEXPORT: true, kwEXTRACT: true, kwFOLLOWER: true, kwFOLLOWING: true,
	kwFRONTEND: true, kwFUNCTIONS: true, kwGRANTS: true, kwINSTALL: true,
	kwINTERMEDIATE: true, kwINTERVAL: true, kwLIST: true, kwLOW_PRIORITY: true,
	kwMATCH: true, kwMINUS: true, kwOBSERVER: true, kwOVERWRITE: true,
	kwPRECEDING: true, kwREPAIR: true, kwREPLICA: true, kwROLE: true, kwROLES: true,
	kwROW: true, kwSEMI: true, kwSETS: true, kwSIGNED: true, kwSYNC: true,
	kwTABLET: true, kwTABLETS: true, kwUNBOUNDED: true, kwUNINSTALL: true,
	kwWHITELIST: true, kwWRITE: true,
}

// showShapeLimitValue reads INTEGER_VALUE or PARAMETER.
func (p *Parser) showShapeLimitValue() bool {
	return p.showShapeToken(tokInt) || p.showShapeToken(tokPlaceholder)
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
	"ident": true, "name": true, "text": true, "str": true, "int": true,
	"expr": true, "wild": true, "sort": true, "limit": true, "user": true,
	"part": true,
}

func (c *showCompiler) expect(t string) {
	if c.pos >= len(c.toks) || c.toks[c.pos] != t {
		panic("malformed show pattern: " + c.src)
	}
	c.pos++
}
