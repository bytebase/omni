package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads the clauses of ALTER TABLE by the grammar of
// oracle/parser/bnf/ALTER-TABLE.bnf
// (https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-TABLE.html).
// Where Oracle's parser departs from the page, the comment beside the code
// says how, and compat.tsv (Oracle 23.26.3's verdicts) holds the cases.
//
// Oracle reads the optional clauses that follow a partition operation, or
// make up a partition's description, in any order: MOVE PARTITION p
// NOCOMPRESS TABLESPACE x and ... UPDATE INDEXES COMPRESS both parse. So
// parseATClauses reads them as a set (atClauses names the clauses an
// operation allows) rather than in the page's order.

// ---------------------------------------------------------------------------
// Token helpers
// ---------------------------------------------------------------------------

// atWord reports whether the current token is the unquoted word w.
func (p *Parser) atWord(w string) bool {
	return p.isKeywordStr(w)
}

// atNextWord reports whether the token after the current one is the
// unquoted word w.
func (p *Parser) atNextWord(w string) bool {
	next := p.peekNext()
	return next.Type != tokQIDENT && p.isIdentLikeStrAt(next, w)
}

// atExpectWord consumes the unquoted word w or fails at the current token.
func (p *Parser) atExpectWord(w string) error {
	if !p.atWord(w) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// atExpect consumes a token of type t or fails at the current token.
func (p *Parser) atExpect(t int) error {
	if p.cur.Type != t {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// isATNameStart reports whether the current token can be a name: an
// identifier, a quoted identifier or a keyword Oracle does not reserve.
func (p *Parser) isATNameStart() bool {
	return p.isIdentLike() && (p.cur.Type == tokQIDENT || !isOracleSQLReservedKeyword(p.cur))
}

// parseATNameString reads a name that is not a reserved word and returns
// it.
func (p *Parser) parseATNameString() (string, error) {
	if !p.isATNameStart() {
		return "", p.syntaxErrorAtCur()
	}
	return p.parseIdentifier()
}

// parseATName reads a name that is not a reserved word, where the AST
// keeps none.
func (p *Parser) parseATName() error {
	if !p.isATNameStart() {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseATQualifiedName reads [ schema. ] name.
func (p *Parser) parseATQualifiedName() error {
	if err := p.parseATName(); err != nil {
		return err
	}
	if p.cur.Type == '.' {
		p.advance()
		return p.parseATName()
	}
	return nil
}

// parseATInteger reads an integer literal.
func (p *Parser) parseATInteger() error {
	if p.cur.Type != tokICONST {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseATNameList reads ( name [, name ]... ).
func (p *Parser) parseATNameList() error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if err := p.parseATName(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATNameListStrings reads ( name [, name ]... ) and returns the names.
func (p *Parser) parseATNameListStrings() ([]string, error) {
	if err := p.atExpect('('); err != nil {
		return nil, err
	}
	var names []string
	for {
		name, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.atExpect(')'); err != nil {
		return nil, err
	}
	return names, nil
}

// parseATExprList reads ( expr [, expr ]... ).
func (p *Parser) parseATExprList() error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		e, err := p.parseExpr()
		if err != nil {
			return err
		}
		if e == nil {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// atOption records the tokens from start to the last one consumed as one
// option of an ALTER TABLE command: the first token is its key, the rest its
// value.
func (p *Parser) atOption(opts *nodes.List, start int) {
	if opts == nil || p.prev.End <= start {
		return
	}
	lexer := NewLexerRange(p.source, start, p.prev.End)
	var tokens []string
	for {
		tok := lexer.NextToken()
		if tok.Type == tokEOF || tok.Loc >= p.prev.End {
			break
		}
		tokens = append(tokens, p.alterTableTokenText(tok))
	}
	if len(tokens) == 0 {
		return
	}
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key:   tokens[0],
		Value: strings.Join(tokens[1:], " "),
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
	})
}

// ---------------------------------------------------------------------------
// Shared small clauses
// ---------------------------------------------------------------------------

// parseATSizeClause reads size_clause: integer [ K | M | G | T | P | E ].
func (p *Parser) parseATSizeClause() error {
	if err := p.parseATInteger(); err != nil {
		return err
	}
	if p.cur.Type != tokQIDENT && p.isIdentLike() {
		switch p.cur.Str {
		case "K", "M", "G", "T", "P", "E":
			p.advance()
		}
	}
	return nil
}

// parseATStorageClause reads storage_clause:
//
//	STORAGE ( { INITIAL size_clause | NEXT size_clause | MINEXTENTS integer
//	          | MAXEXTENTS { integer | UNLIMITED } | maxsize_clause
//	          | PCTINCREASE integer | FREELISTS integer
//	          | FREELIST GROUPS integer | OPTIMAL [ size_clause | NULL ]
//	          | BUFFER_POOL { KEEP | RECYCLE | DEFAULT }
//	          | FLASH_CACHE { KEEP | NONE | DEFAULT }
//	          | CELL_FLASH_CACHE { KEEP | NONE | DEFAULT } | ENCRYPT }... )
func (p *Parser) parseATStorageClause() error {
	p.advance() // STORAGE
	if err := p.atExpect('('); err != nil {
		return err
	}
	for n := 0; ; n++ {
		if p.cur.Type == ')' && n > 0 {
			p.advance()
			return nil
		}
		var err error
		switch {
		case p.atWord("INITIAL"), p.atWord("NEXT"):
			p.advance()
			err = p.parseATSizeClause()
		case p.atWord("MINEXTENTS"), p.atWord("PCTINCREASE"), p.atWord("FREELISTS"):
			p.advance()
			err = p.parseATInteger()
		case p.atWord("MAXEXTENTS"):
			p.advance()
			if p.atWord("UNLIMITED") {
				p.advance()
			} else {
				err = p.parseATInteger()
			}
		case p.atWord("MAXSIZE"):
			p.advance()
			if p.atWord("UNLIMITED") {
				p.advance()
			} else {
				err = p.parseATSizeClause()
			}
		case p.atWord("FREELIST"):
			p.advance()
			if err = p.atExpectWord("GROUPS"); err == nil {
				err = p.parseATInteger()
			}
		case p.atWord("OPTIMAL"):
			p.advance()
			if p.cur.Type == kwNULL {
				p.advance()
			} else if p.cur.Type == tokICONST {
				err = p.parseATSizeClause()
			}
		case p.atWord("BUFFER_POOL"):
			p.advance()
			err = p.atOneOf("KEEP", "RECYCLE", "DEFAULT")
		case p.atWord("FLASH_CACHE"), p.atWord("CELL_FLASH_CACHE"):
			p.advance()
			err = p.atOneOf("KEEP", "NONE", "DEFAULT")
		case p.atWord("ENCRYPT"):
			p.advance()
		default:
			return p.syntaxErrorAtCur()
		}
		if err != nil {
			return err
		}
	}
}

// atOneOf consumes one of the given words or fails at the current token.
func (p *Parser) atOneOf(words ...string) error {
	for _, w := range words {
		if p.atWord(w) {
			p.advance()
			return nil
		}
	}
	return p.syntaxErrorAtCur()
}

// atOptOneOf consumes one of the given words if it is next.
func (p *Parser) atOptOneOf(words ...string) bool {
	for _, w := range words {
		if p.atWord(w) {
			p.advance()
			return true
		}
	}
	return false
}

// isATPhysicalAttrStart reports a physical_attributes_clause item:
// PCTFREE, PCTUSED, INITRANS or storage_clause.
func (p *Parser) isATPhysicalAttrStart() bool {
	return p.atWord("PCTFREE") || p.atWord("PCTUSED") || p.atWord("INITRANS") || p.atWord("STORAGE")
}

// parseATPhysicalAttr reads one physical_attributes_clause item:
// { PCTFREE integer | PCTUSED integer | INITRANS integer | storage_clause }.
func (p *Parser) parseATPhysicalAttr() error {
	if p.atWord("STORAGE") {
		return p.parseATStorageClause()
	}
	p.advance()
	err := p.parseATInteger()
	return err
}

// isATLoggingStart reports logging_clause:
// { LOGGING | NOLOGGING | FILESYSTEM_LIKE_LOGGING }.
func (p *Parser) isATLoggingStart() bool {
	return p.cur.Type == kwLOGGING || p.cur.Type == kwNOLOGGING || p.atWord("FILESYSTEM_LIKE_LOGGING")
}

// isATSegmentAttrStart reports a segment_attributes_clause item:
// physical_attributes_clause, TABLESPACE tablespace or logging_clause.
func (p *Parser) isATSegmentAttrStart() bool {
	return p.isATPhysicalAttrStart() || p.cur.Type == kwTABLESPACE || p.isATLoggingStart()
}

// parseATSegmentAttr reads one segment_attributes_clause item.
func (p *Parser) parseATSegmentAttr() error {
	switch {
	case p.isATPhysicalAttrStart():
		return p.parseATPhysicalAttr()
	case p.cur.Type == kwTABLESPACE:
		p.advance()
		err := p.parseATName()
		return err
	case p.isATLoggingStart():
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseATSegmentAttrs reads segment_attributes_clause:
// { physical_attributes_clause | TABLESPACE tablespace | logging_clause }...
func (p *Parser) parseATSegmentAttrs(opts *nodes.List) error {
	if !p.isATSegmentAttrStart() {
		return p.syntaxErrorAtCur()
	}
	for p.isATSegmentAttrStart() {
		start := p.pos()
		if err := p.parseATSegmentAttr(); err != nil {
			return err
		}
		p.atOption(opts, start)
	}
	return nil
}

// isATTableCompressionStart reports table_compression or prefix_compression:
// COMPRESS, NOCOMPRESS, ROW STORE COMPRESS, COLUMN STORE COMPRESS.
func (p *Parser) isATTableCompressionStart() bool {
	switch {
	case p.cur.Type == kwCOMPRESS, p.cur.Type == kwNOCOMPRESS, p.cur.Type == kwROW:
		return true
	case p.cur.Type == kwCOLUMN:
		return p.atNextWord("STORE")
	}
	return false
}

// parseATTableCompression reads table_compression, or prefix_compression
// (COMPRESS integer) where prefix is set:
//
//	{ COMPRESS [ BASIC | ADVANCED ]
//	| ROW STORE COMPRESS [ BASIC | ADVANCED ]
//	| COLUMN STORE COMPRESS [ FOR { QUERY | ARCHIVE } [ LOW | HIGH ] ]
//	    [ [ NO ] ROW LEVEL LOCKING ]
//	| NOCOMPRESS }
func (p *Parser) parseATTableCompression(prefix bool) error {
	switch {
	case p.cur.Type == kwNOCOMPRESS:
		p.advance()
		return nil
	case p.cur.Type == kwCOMPRESS:
		p.advance()
		if prefix && p.cur.Type == tokICONST {
			p.advance()
			return nil
		}
		p.atOptOneOf("BASIC", "ADVANCED")
		return nil
	case p.cur.Type == kwROW:
		p.advance() // ROW
		if err := p.atExpectWord("STORE"); err != nil {
			return err
		}
		// Oracle takes ROW STORE and COLUMN STORE without COMPRESS, as it
		// does in an ILM policy.
		if p.cur.Type != kwCOMPRESS {
			return nil
		}
		p.advance()
		p.atOptOneOf("BASIC", "ADVANCED")
		return nil
	case p.cur.Type == kwCOLUMN:
		p.advance() // COLUMN
		p.advance() // STORE
		if p.cur.Type != kwCOMPRESS {
			return nil
		}
		p.advance()
		if p.cur.Type == kwFOR {
			p.advance()
			if err := p.atOneOf("QUERY", "ARCHIVE"); err != nil {
				return err
			}
			p.atOptOneOf("LOW", "HIGH")
		}
		return p.parseATRowLevelLocking()
	}
	return p.syntaxErrorAtCur()
}

// parseATRowLevelLocking reads [ [ NO ] ROW LEVEL LOCKING ].
func (p *Parser) parseATRowLevelLocking() error {
	if p.atWord("NO") && p.peekNext().Type == kwROW {
		p.advance()
	} else if !(p.cur.Type == kwROW && p.atNextWord("LEVEL")) {
		return nil
	}
	if err := p.atExpect(kwROW); err != nil {
		return err
	}
	if err := p.atExpectWord("LEVEL"); err != nil {
		return err
	}
	return p.atExpectWord("LOCKING")
}

// isATIndexCompressionStart reports index_compression:
// { COMPRESS [ integer ] | NOCOMPRESS }.
func (p *Parser) isATIndexCompressionStart() bool {
	return p.cur.Type == kwCOMPRESS || p.cur.Type == kwNOCOMPRESS
}

// parseATIndexCompression reads index_compression.
func (p *Parser) parseATIndexCompression() error {
	if p.cur.Type == kwCOMPRESS {
		// COMPRESS ADVANCED [ LOW | HIGH ] is accepted where the group ends
		// right after it (compat.tsv's ALTER MATERIALIZED VIEW rows); a word
		// after it is refused at ADVANCED, as it was before.
		toks := p.peekAhead(3)
		if p.alterMview && len(toks) == 3 && p.isIdentLikeStrAt(toks[0], "ADVANCED") &&
			(toks[1].Type == ')' ||
				((p.isIdentLikeStrAt(toks[1], "LOW") || p.isIdentLikeStrAt(toks[1], "HIGH")) && toks[2].Type == ')')) {
			p.advance() // COMPRESS
			p.advance() // ADVANCED
			p.atOptOneOf("LOW", "HIGH")
			return nil
		}
		p.advance()
		if p.cur.Type == tokICONST {
			p.advance()
		}
		return nil
	}
	p.advance() // NOCOMPRESS
	return nil
}

// isATInmemoryStart reports inmemory_clause: INMEMORY or NO INMEMORY.
func (p *Parser) isATInmemoryStart() bool {
	return p.atWord("INMEMORY") || (p.atWord("NO") && p.atNextWord("INMEMORY"))
}

// parseATInmemoryClause reads inmemory_clause:
// { INMEMORY [ inmemory_attributes ] | NO INMEMORY }.
func (p *Parser) parseATInmemoryClause() error {
	if p.atWord("NO") {
		p.advance() // NO
		p.advance() // INMEMORY
		return nil
	}
	p.advance() // INMEMORY
	return p.parseATInmemoryAttributes()
}

// parseATInmemoryAttributes reads inmemory_attributes:
// [ inmemory_memcompress ] [ inmemory_priority ] [ inmemory_distribute ]
// [ inmemory_duplicate ].
func (p *Parser) parseATInmemoryAttributes() error {
	if err := p.parseATInmemoryMemcompress(); err != nil {
		return err
	}
	return p.parseATInmemoryRest()
}

// parseATInmemoryRest reads [ inmemory_priority ] [ inmemory_distribute ]
// [ inmemory_duplicate ].
func (p *Parser) parseATInmemoryRest() error {
	if p.atWord("PRIORITY") {
		p.advance()
		if err := p.atOneOf("NONE", "LOW", "MEDIUM", "HIGH", "CRITICAL"); err != nil {
			return err
		}
	}
	if p.atWord("DISTRIBUTE") {
		p.advance()
		if p.atWord("AUTO") {
			p.advance()
		} else if p.cur.Type == kwBY {
			byTok := p.cur
			p.advance()
			switch {
			case p.cur.Type == kwROWID:
				p.advance()
				if !p.atWord("RANGE") {
					return p.syntaxErrorAtTok(byTok)
				}
				p.advance()
			case p.cur.Type == kwPARTITION, p.cur.Type == kwSUBPARTITION:
				p.advance()
			default:
				return p.syntaxErrorAtTok(byTok)
			}
		}
		if p.cur.Type == kwFOR {
			p.advance()
			if err := p.atExpectWord("SERVICE"); err != nil {
				return err
			}
			// A service is named by a string, not an identifier.
			if p.cur.Type == tokSCONST {
				p.advance()
			} else if err := p.atOneOf("DEFAULT", "ALL", "NONE"); err != nil {
				return err
			}
		}
	}
	switch {
	case p.atWord("DUPLICATE"):
		p.advance()
		if p.cur.Type == kwALL {
			p.advance()
		} else if p.cur.Type == kwFOR && p.alterMview {
			// On ALTER MATERIALIZED VIEW, Oracle consumes FOR after DUPLICATE
			// before refusing the unsupported FOR SERVICE clause at the next token.
			p.advance()
			return p.syntaxErrorAtCur()
		}
	case p.atWord("NO") && !p.atNextWord("INMEMORY") && p.peekNext().Type != kwFLASHBACK:
		p.advance()
		return p.atExpectWord("DUPLICATE")
	}
	return nil
}

// isATInmemoryMemcompressStart reports inmemory_memcompress: MEMCOMPRESS or
// NO MEMCOMPRESS.
func (p *Parser) isATInmemoryMemcompressStart() bool {
	return p.atWord("MEMCOMPRESS") || (p.atWord("NO") && p.atNextWord("MEMCOMPRESS"))
}

// parseATInmemoryMemcompress reads [ inmemory_memcompress ]:
//
//	MEMCOMPRESS FOR { DML | QUERY [ LOW | HIGH ] | CAPACITY [ LOW | HIGH ] }
//	| MEMCOMPRESS AUTO | NO MEMCOMPRESS
func (p *Parser) parseATInmemoryMemcompress() error {
	switch {
	case p.atWord("MEMCOMPRESS"):
		p.advance()
		if p.atWord("AUTO") {
			p.advance()
			return nil
		}
		if err := p.atExpect(kwFOR); err != nil {
			return err
		}
		switch {
		case p.atWord("DML"):
			p.advance()
		case p.atWord("QUERY"), p.atWord("CAPACITY"):
			p.advance()
			p.atOptOneOf("LOW", "HIGH")
		default:
			return p.syntaxErrorAtCur()
		}
	case p.atWord("NO") && p.atNextWord("MEMCOMPRESS"):
		p.advance()
		p.advance()
	}
	return nil
}

// parseATInmemoryColumnItem reads one item of inmemory_column_clause:
// { INMEMORY [ inmemory_memcompress ] | NO INMEMORY } ( column [, column ]... ).
func (p *Parser) parseATInmemoryColumnItem() error {
	var noTok Token
	hasNo := false
	if p.atWord("NO") {
		noTok = p.cur
		hasNo = true
		p.advance() // NO
		p.advance() // INMEMORY
	} else {
		p.advance() // INMEMORY
		if p.atWord("NO") && (p.peekNext().Type == ';' || p.peekNext().Type == tokEOF) {
			// Oracle takes INMEMORY NO ending the statement.
			p.advance()
			return nil
		}
		if p.atWord("NO") && p.atNextWord("MEMCOMPRESS") {
			noTok = p.cur
			hasNo = true
		}
		if err := p.parseATInmemoryMemcompress(); err != nil {
			return err
		}
	}
	if p.cur.Type != '(' && hasNo {
		return p.syntaxErrorAtTok(noTok)
	}
	return p.parseATNameList()
}

// parseATIlmClause reads ilm_clause:
//
//	ILM { ADD POLICY ilm_policy_clause
//	    | { DELETE | ENABLE | DISABLE } POLICY policy_name
//	    | DELETE_ALL | ENABLE_ALL | DISABLE_ALL }
func (p *Parser) parseATIlmClause() error {
	p.advance() // ILM
	switch {
	case p.cur.Type == kwADD:
		p.advance()
		if err := p.atExpectWord("POLICY"); err != nil {
			return err
		}
		return p.parseATIlmPolicy()
	case p.cur.Type == kwDELETE, p.cur.Type == kwENABLE, p.cur.Type == kwDISABLE:
		p.advance()
		if err := p.atExpectWord("POLICY"); err != nil {
			return err
		}
		err := p.parseATName()
		return err
	case p.atWord("DELETE_ALL"), p.atWord("ENABLE_ALL"), p.atWord("DISABLE_ALL"):
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseATIlmPolicy reads ilm_policy_clause:
//
//	ilm_compression_policy | ilm_tiering_policy | ilm_inmemory_policy
//
//	ilm_compression_policy:
//	    { table_compression { SEGMENT | GROUP }
//	    | { ROW STORE COMPRESS ADVANCED | COLUMN STORE COMPRESS FOR QUERY } ROW }
//	    AFTER ilm_time_period OF { NO ACCESS | NO MODIFICATION | CREATION }
//	ilm_tiering_policy:
//	    TIER TO tablespace [ READ ONLY ] [ SEGMENT | GROUP ] [ ON function ]
//	ilm_inmemory_policy:
//	    { SET INMEMORY [ inmemory_attributes ] | MODIFY INMEMORY inmemory_memcompress
//	    | NO INMEMORY } [ SEGMENT ] AFTER ilm_time_period OF ...
//
// Oracle reads ROW STORE and COLUMN STORE without COMPRESS here.
func (p *Parser) parseATIlmPolicy() error {
	switch {
	case p.atWord("TIER"):
		p.advance()
		if err := p.atExpect(kwTO); err != nil {
			return err
		}
		if err := p.parseATName(); err != nil {
			return err
		}
		if p.cur.Type == kwREAD {
			p.advance()
			if err := p.atExpectWord("ONLY"); err != nil {
				return err
			}
		}
		switch {
		case p.atWord("SEGMENT"), p.cur.Type == kwGROUP:
			p.advance()
		case p.cur.Type == kwROW:
			// Oracle reads ROW here and refuses what follows it, the end
			// of the statement included: a tiering policy has no row scope.
			p.advance()
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwON {
			p.advance()
			if err := p.parseATName(); err != nil {
				return err
			}
		}
		return nil
	case p.cur.Type == kwSET:
		p.advance()
		if err := p.atExpectWord("INMEMORY"); err != nil {
			return err
		}
		if err := p.parseATInmemoryAttributes(); err != nil {
			return err
		}
		return p.parseATIlmCondition(false, true)
	case p.cur.Type == kwMODIFY:
		p.advance()
		if err := p.atExpectWord("INMEMORY"); err != nil {
			return err
		}
		if p.atWord("NO") {
			p.advance()
			if err := p.atExpectWord("MEMCOMPRESS"); err != nil {
				return err
			}
			return p.parseATIlmCondition(false, true)
		}
		if !p.isATInmemoryMemcompressStart() {
			return p.syntaxErrorAtCur()
		}
		if err := p.parseATInmemoryMemcompress(); err != nil {
			return err
		}
		return p.parseATIlmCondition(false, true)
	case p.atWord("NO"):
		p.advance()
		if err := p.atExpectWord("INMEMORY"); err != nil {
			return err
		}
		return p.parseATIlmCondition(false, true)
	case p.cur.Type == kwROW, p.cur.Type == kwCOLUMN:
		column := p.cur.Type == kwCOLUMN
		rowLevel := false
		p.advance()
		if err := p.atExpectWord("STORE"); err != nil {
			return err
		}
		if p.cur.Type == kwCOMPRESS {
			p.advance()
			if column && p.cur.Type == kwFOR {
				p.advance()
				if err := p.atExpectWord("QUERY"); err != nil {
					return err
				}
				p.atOptOneOf("LOW", "HIGH")
				rowLevel = true
			} else if !column && p.atWord("ADVANCED") {
				p.advance()
				rowLevel = true
			}
		}
		return p.parseATIlmCondition(rowLevel, false)
	case p.cur.Type == kwCOMPRESS, p.cur.Type == kwNOCOMPRESS:
		if err := p.parseATTableCompression(false); err != nil {
			return err
		}
		return p.parseATIlmCondition(false, false)
	}
	return p.syntaxErrorAtCur()
}

// parseATIlmCondition reads { SEGMENT | GROUP | ROW } AFTER ilm_time_period
// OF { NO ACCESS | NO MODIFICATION | CREATION }, or ON function. ROW stands
// only after ROW STORE COMPRESS ADVANCED or COLUMN STORE COMPRESS FOR QUERY
// (rowLevel); an inmemory policy takes only SEGMENT, and may leave it out.
func (p *Parser) parseATIlmCondition(rowLevel, inmemory bool) error {
	switch {
	case p.atWord("SEGMENT"):
		p.advance()
	case !inmemory && (p.cur.Type == kwGROUP || rowLevel && p.cur.Type == kwROW):
		p.advance()
	case !inmemory && p.cur.Type == kwROW:
		// Oracle consumes ROW on a non-row-level compression policy before
		// refusing what follows it (matching TIER TO a ROW above).
		p.advance()
		return p.syntaxErrorAtCur()
	case !inmemory:
		return p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwON {
		p.advance()
		err := p.parseATName()
		return err
	}
	if err := p.atExpectWord("AFTER"); err != nil {
		return err
	}
	// ilm_time_period: integer { DAY | DAYS | MONTH | MONTHS | YEAR | YEARS };
	// Oracle takes the integer alone as days.
	if err := p.parseATInteger(); err != nil {
		return err
	}
	p.atOptOneOf("DAY", "DAYS", "MONTH", "MONTHS", "YEAR", "YEARS")
	if err := p.atExpect(kwOF); err != nil {
		return err
	}
	switch {
	case p.atWord("NO"):
		p.advance()
		return p.atOneOf("ACCESS", "MODIFICATION")
	case p.atWord("CREATION"):
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
}

// isATAllocateExtentStart reports allocate_extent_clause.
func (p *Parser) isATAllocateExtentStart() bool {
	return p.atWord("ALLOCATE")
}

// parseATAllocateExtent reads allocate_extent_clause:
//
//	ALLOCATE EXTENT [ ( { SIZE size_clause | DATAFILE 'filename'
//	                    | INSTANCE integer }... ) ]
func (p *Parser) parseATAllocateExtent() error {
	p.advance() // ALLOCATE
	if err := p.atExpectWord("EXTENT"); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	// Each option stands once (ORA-03278).
	seen := map[string]bool{}
	for n := 0; ; n++ {
		if p.cur.Type == ')' && n > 0 {
			p.advance()
			return nil
		}
		if seen[p.cur.Str] {
			return p.syntaxErrorAtCur()
		}
		seen[p.cur.Str] = true
		var err error
		switch {
		case p.cur.Type == kwSIZE:
			p.advance()
			err = p.parseATSizeClause()
		case p.atWord("DATAFILE"):
			p.advance()
			err = p.atExpect(tokSCONST)
		case p.atWord("INSTANCE"):
			p.advance()
			err = p.parseATInteger()
		default:
			return p.syntaxErrorAtCur()
		}
		if err != nil {
			return err
		}
	}
}

// isATDeallocateStart reports deallocate_unused_clause.
func (p *Parser) isATDeallocateStart() bool {
	return p.atWord("DEALLOCATE")
}

// parseATDeallocateUnused reads deallocate_unused_clause:
// DEALLOCATE UNUSED [ KEEP size_clause ].
func (p *Parser) parseATDeallocateUnused() error {
	p.advance() // DEALLOCATE
	if err := p.atExpectWord("UNUSED"); err != nil {
		return err
	}
	if p.atWord("KEEP") {
		p.advance()
		return p.parseATSizeClause()
	}
	return nil
}

// isATShrinkStart reports SHRINK SPACE [ CASCADE ].
func (p *Parser) isATShrinkStart() bool {
	return p.atWord("SHRINK")
}

// parseATShrink reads SHRINK SPACE [ COMPACT ] [ CASCADE ].
func (p *Parser) parseATShrink() error {
	p.advance() // SHRINK
	if err := p.atExpectWord("SPACE"); err != nil {
		return err
	}
	p.atOptOneOf("COMPACT")
	if p.cur.Type == kwCASCADE {
		p.advance()
	}
	return nil
}

// parseATEncryptionSpec reads encryption_spec after ENCRYPT:
//
//	[ USING 'encrypt_algorithm' ] [ IDENTIFIED BY password ]
//	[ 'integrity_algorithm' ] [ [ NO ] SALT ]
func (p *Parser) parseATEncryptionSpec() error {
	if p.cur.Type == kwUSING {
		p.advance()
		if err := p.atExpect(tokSCONST); err != nil {
			return err
		}
	}
	if p.cur.Type == kwIDENTIFIED {
		p.advance()
		if err := p.atExpect(kwBY); err != nil {
			return err
		}
		if err := p.parseATName(); err != nil {
			return err
		}
	}
	if p.cur.Type == tokSCONST {
		p.advance()
	}
	if p.atWord("NO") && p.atNextWord("SALT") {
		p.advance()
	}
	p.atOptOneOf("SALT")
	return nil
}

// parseATLobRetention reads LOB_retention_clause after RETENTION:
// [ MAX | MIN integer | AUTO | NONE ].
func (p *Parser) parseATLobRetention() error {
	switch {
	case p.atWord("MIN"):
		p.advance()
		err := p.parseATInteger()
		return err
	default:
		p.atOptOneOf("MAX", "AUTO", "NONE")
	}
	return nil
}

// parseATLobCacheLogging reads { CACHE | NOCACHE | CACHE READS }
// [ logging_clause ].
func (p *Parser) parseATLobCacheLogging() error {
	if p.cur.Type == kwCACHE {
		p.advance()
		p.atOptOneOf("READS")
	} else {
		p.advance() // NOCACHE
	}
	if p.isATLoggingStart() {
		p.advance()
	}
	return nil
}

// parseATModifyLobParameters reads ( modify_LOB_parameters ):
//
//	( { storage_clause | PCTVERSION integer | FREEPOOLS integer
//	  | REBUILD FREEPOOLS | LOB_retention_clause | LOB_deduplicate_clause
//	  | LOB_compression_clause | { ENCRYPT encryption_spec | DECRYPT }
//	  | { CACHE | NOCACHE | CACHE READS } [ logging_clause ]
//	  | allocate_extent_clause | SHRINK SPACE [ CASCADE ]
//	  | deallocate_unused_clause }... )
func (p *Parser) parseATModifyLobParameters() error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	// Each parameter stands once (ORA-32600 for RETENTION RETENTION).
	seen := map[string]bool{}
	for n := 0; ; n++ {
		if p.cur.Type == ')' && n > 0 {
			p.advance()
			return nil
		}
		if seen[p.cur.Str] {
			return p.syntaxErrorAtCur()
		}
		seen[p.cur.Str] = true
		var err error
		switch {
		case p.atWord("STORAGE"):
			err = p.parseATStorageClause()
		case p.atWord("PCTVERSION"), p.atWord("FREEPOOLS"):
			p.advance()
			err = p.parseATInteger()
		case p.atWord("REBUILD"):
			p.advance()
			err = p.atExpectWord("FREEPOOLS")
		case p.atWord("RETENTION"):
			p.advance()
			err = p.parseATLobRetention()
		case p.atWord("DEDUPLICATE"), p.atWord("KEEP_DUPLICATES"), p.atWord("DECRYPT"), p.cur.Type == kwNOCOMPRESS:
			p.advance()
		case p.cur.Type == kwCOMPRESS:
			p.advance()
			p.atOptOneOf("HIGH", "MEDIUM", "LOW")
		case p.atWord("ENCRYPT"):
			p.advance()
			err = p.parseATEncryptionSpec()
		case p.cur.Type == kwCACHE, p.atWord("NOCACHE"):
			err = p.parseATLobCacheLogging()
		case p.isATAllocateExtentStart():
			err = p.parseATAllocateExtent()
		case p.isATShrinkStart():
			err = p.parseATShrink()
		case p.isATDeallocateStart():
			err = p.parseATDeallocateUnused()
		default:
			return p.syntaxErrorAtCur()
		}
		if err != nil {
			return err
		}
	}
}

// parseATLobStorageParameters reads ( LOB_parameters ) of a STORE AS:
//
//	( { TABLESPACE tablespace | TABLESPACE SET tablespace_set
//	  | { ENABLE | DISABLE } STORAGE IN ROW [ 4000 | 8000 ]
//	  | CHUNK integer | PCTVERSION integer | FREEPOOLS integer
//	  | LOB_retention_clause | LOB_deduplicate_clause | LOB_compression_clause
//	  | { ENCRYPT encryption_spec | DECRYPT }
//	  | { CACHE | NOCACHE | CACHE READS } [ logging_clause ]
//	  | storage_clause }... )
//
// Oracle takes logging_clause only after CACHE or NOCACHE here: the page's
// standalone logging_clause is refused (ORA-22853 for STORE AS (LOGGING)).
func (p *Parser) parseATLobStorageParameters() error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for n := 0; ; n++ {
		if p.cur.Type == ')' && n > 0 {
			p.advance()
			return nil
		}
		var err error
		switch {
		case p.cur.Type == kwTABLESPACE:
			p.advance()
			if p.cur.Type == kwSET {
				p.advance()
			}
			err = p.parseATName()
		case p.cur.Type == kwENABLE, p.cur.Type == kwDISABLE:
			p.advance()
			if err = p.atExpectWord("STORAGE"); err == nil {
				if err = p.atExpect(kwIN); err == nil {
					err = p.atExpect(kwROW)
				}
			}
		case p.atWord("CHUNK"), p.atWord("PCTVERSION"), p.atWord("FREEPOOLS"):
			p.advance()
			err = p.parseATInteger()
		case p.atWord("RETENTION"):
			p.advance()
			err = p.parseATLobRetention()
		case p.atWord("DEDUPLICATE"), p.atWord("KEEP_DUPLICATES"), p.atWord("DECRYPT"), p.cur.Type == kwNOCOMPRESS:
			p.advance()
		case p.cur.Type == kwCOMPRESS:
			p.advance()
			p.atOptOneOf("HIGH", "MEDIUM", "LOW")
		case p.atWord("ENCRYPT"):
			p.advance()
			err = p.parseATEncryptionSpec()
		case p.cur.Type == kwCACHE, p.atWord("NOCACHE"):
			err = p.parseATLobCacheLogging()
		case p.atWord("STORAGE"):
			err = p.parseATStorageClause()
		default:
			return p.syntaxErrorAtCur()
		}
		if err != nil {
			return err
		}
	}
}

// parseATLobStoreAs reads the STORE AS part of LOB_storage_clause after
// LOB ( item [, item ]... ):
//
//	STORE AS { [ SECUREFILE | BASICFILE ] [ LOB_segname ] ( LOB_storage_parameters )
//	         | SECUREFILE | BASICFILE | LOB_segname }
func (p *Parser) parseATLobStoreAs() error {
	if err := p.atExpectWord("STORE"); err != nil {
		return err
	}
	if err := p.atExpect(kwAS); err != nil {
		return err
	}
	any := p.atOptOneOf("SECUREFILE", "BASICFILE")
	if p.cur.Type != '(' && p.isATNameStart() && (!any || !p.isATAfterLobStorage()) {
		p.advance()
		any = true
	}
	if p.cur.Type == '(' {
		return p.parseATLobStorageParameters()
	}
	if !any {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// isATAfterLobStorage reports a word that starts what may follow a LOB
// storage clause, which an optional LOB_segname before it does not take:
// another column property, a partition clause or an ALTER TABLE action.
func (p *Parser) isATAfterLobStorage() bool {
	return p.isATLobStorageStart() || p.cur.Type == kwNESTED || p.atWord("XMLTYPE") ||
		(p.cur.Type == kwJSON && p.peekNext().Type == '(') || p.isATAlterTableActionWord() ||
		p.atClauseAt(atcTablePartitionDescription|atcPartitioningStorage|atcIndexMaintenance|atcParallel) != 0
}

// parseATLobItemList reads ( lob_item [, lob_item ]... ).
func (p *Parser) parseATLobItemList() error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if err := p.parseATItemPath(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATItemPath reads a LOB or collection item: a column, or an
// attribute of an object column, column.attribute [ .attribute ]...
func (p *Parser) parseATItemPath() error {
	if err := p.parseATName(); err != nil {
		return err
	}
	for p.cur.Type == '.' {
		p.advance()
		if err := p.parseATName(); err != nil {
			return err
		}
	}
	return nil
}

// isATLobStorageStart reports LOB_storage_clause or a VARRAY's STORE AS:
// LOB ( ... or VARRAY varray.
func (p *Parser) isATLobStorageStart() bool {
	return (p.atWord("LOB") && p.peekNext().Type == '(') || p.atWord("VARRAY")
}

// parseATLobStorage reads one of
//
//	LOB ( lob_item [, lob_item ]... ) STORE AS ...
//	VARRAY varray_item [ substitutable_column_clause ]
//	    STORE AS [ SECUREFILE | BASICFILE ] LOB { LOB_segname | ( LOB_storage_parameters ) }
//	VARRAY varray_item substitutable_column_clause
//
// Oracle takes a LOB_segname after LOB only when neither SECUREFILE nor
// BASICFILE comes before it.
func (p *Parser) parseATLobStorage() error {
	if p.atWord("LOB") {
		p.advance()
		if err := p.parseATLobItemList(); err != nil {
			return err
		}
		return p.parseATLobStoreAs()
	}
	p.advance() // VARRAY
	if err := p.parseATName(); err != nil {
		return err
	}
	var sub bool
	if err := p.parseATSubstitutable(&sub); err != nil {
		return err
	}
	if sub && !p.atWord("STORE") {
		return nil
	}
	if err := p.atExpectWord("STORE"); err != nil {
		return err
	}
	if err := p.atExpect(kwAS); err != nil {
		return err
	}
	file := p.atOptOneOf("SECUREFILE", "BASICFILE")
	if err := p.atExpectWord("LOB"); err != nil {
		return err
	}
	if p.cur.Type == '(' {
		return p.parseATLobStorageParameters()
	}
	if file {
		return nil
	}
	return p.parseATName()
}

// parseATLobPartitionStorage reads ( LOB_partition_storage [, ...] ) after
// a LOB or VARRAY storage clause of column_properties:
//
//	PARTITION partition { LOB_storage_clause | varray_col_properties }...
//	    [ ( SUBPARTITION subpartition { LOB_partitioning_storage
//	      | varray_col_properties }... ) ]
func (p *Parser) parseATLobPartitionStorage() error {
	p.advance() // (
	for {
		if err := p.atExpect(kwPARTITION); err != nil {
			return err
		}
		if err := p.parseATName(); err != nil {
			return err
		}
		// Oracle parses repeated storage clauses inside a partition item and
		// rejects a second one at the token that closes or follows it.
		count := 0
		for p.isATLobStorageStart() || p.atWord("LOB") {
			if err := p.parseATLobStorage(); err != nil {
				return err
			}
			count++
		}
		if count == 0 && p.cur.Type == '(' {
			p.advance()
			return p.atExpect(kwSUBPARTITION)
		}
		if count == 0 || count > 1 {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type == '(' {
			p.advance()
			if err := p.atExpect(kwSUBPARTITION); err != nil {
				return err
			}
			if err := p.parseATName(); err != nil {
				return err
			}
			if !p.isATLobStorageStart() {
				return p.syntaxErrorAtCur()
			}
			if err := p.parseATLobStorage(); err != nil {
				return err
			}
			if err := p.atExpect(')'); err != nil {
				return err
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// isATLobModifyStart reports LOB ( lob_item ) ( modify_LOB_parameters ) of
// a partition or subpartition. Oracle refuses the page's VARRAY varray
// ( modify_LOB_parameters ) there (ORA-00905 at its parenthesis); ALTER
// TABLE MODIFY VARRAY takes it.
func (p *Parser) isATLobModifyStart() bool {
	return p.atWord("LOB") && p.peekNext().Type == '('
}

// parseATLobModify reads { LOB ( lob_item ) | VARRAY varray }
// ( modify_LOB_parameters ).
func (p *Parser) parseATLobModify() error {
	if p.atWord("LOB") {
		p.advance()
		if err := p.atExpect('('); err != nil {
			return err
		}
		if err := p.parseATItemPath(); err != nil {
			return err
		}
		if p.alterMview && p.cur.Type == ',' {
			p.advance()
			return p.syntaxErrorAtCur()
		}
		if err := p.atExpect(')'); err != nil {
			return err
		}
	} else {
		p.advance() // VARRAY
		if err := p.parseATItemPath(); err != nil {
			return err
		}
	}
	return p.parseATModifyLobParameters()
}

// ---------------------------------------------------------------------------
// Partition clauses
// ---------------------------------------------------------------------------

// atClauses names the optional clauses an ALTER TABLE operation or a
// partition description allows; parseATClauses reads them in any order.
type atClauses uint64

const (
	// atcSegment is segment_attributes_clause.
	atcSegment atClauses = 1 << iota
	// atcTablespace is TABLESPACE tablespace alone (partitioning_storage_clause).
	atcTablespace
	// atcPhysical is physical_attributes_clause and logging_clause.
	atcPhysical
	// atcCompression is table_compression or prefix_compression.
	atcCompression
	// atcInmemory is inmemory_clause.
	atcInmemory
	// atcIlm is ilm_clause.
	atcIlm
	// atcReadOnly is read_only_clause.
	atcReadOnly
	// atcIndexing is indexing_clause.
	atcIndexing
	// atcSegmentCreation is deferred_segment_creation.
	atcSegmentCreation
	// atcOverflow is OVERFLOW [ segment_attributes_clause ].
	atcOverflow
	// atcOverflowTablespace is OVERFLOW TABLESPACE tablespace.
	atcOverflowTablespace
	// atcLobStorage is LOB ( ... ) STORE AS ... and VARRAY ... STORE AS ....
	atcLobStorage
	// atcLobModify is { LOB ( item ) | VARRAY varray } ( modify_LOB_parameters ).
	atcLobModify
	// atcUpdateIndexes is update_index_clauses but INVALIDATE GLOBAL INDEXES.
	atcUpdateIndexes
	// atcInvalidate is INVALIDATE GLOBAL INDEXES.
	atcInvalidate
	// atcParallel is parallel_clause.
	atcParallel
	// atcClustering is allow_disallow_clustering.
	atcClustering
	// atcOnline is ONLINE.
	atcOnline
	// atcFilter is filter_condition.
	atcFilter
	// atcSplitFilter is filter_condition on SPLIT, where Oracle enters the
	// clause on INCLUDING even when ROWS is missing.
	atcSplitFilter
	// atcDependentTables is dependent_tables_clause.
	atcDependentTables
	// atcMappingTable is MAPPING TABLE.
	atcMappingTable
	// atcAllocate is allocate_extent_clause.
	atcAllocate
	// atcDeallocate is deallocate_unused_clause.
	atcDeallocate
	// atcShrink is SHRINK SPACE [ CASCADE ].
	atcShrink
	// atcCascade is CASCADE.
	atcCascade
	// atcExchange is the options of EXCHANGE PARTITION: { INCLUDING |
	// EXCLUDING } INDEXES, { WITH | WITHOUT } VALIDATION, exceptions_clause.
	atcExchange
	// atcNestedTable is nested_table_col_properties ... STORE AS.
	atcNestedTable
	// atcExternal is external_part_subpart_data_props: [ DEFAULT DIRECTORY
	// directory ] [ LOCATION ( ... ) ].
	atcExternal
)

// atcTablePartitionDescription is table_partition_description:
//
//	[ deferred_segment_creation ] [ read_only_clause ] [ indexing_clause ]
//	[ segment_attributes_clause ] [ table_compression | prefix_compression ]
//	[ inmemory_clause ] [ ilm_clause ] [ OVERFLOW [ segment_attributes_clause ] ]
//	[ { LOB_storage_clause | varray_col_properties | nested_table_col_properties }... ]
const atcTablePartitionDescription = atcSegmentCreation | atcReadOnly | atcIndexing | atcSegment |
	atcCompression | atcInmemory | atcIlm | atcOverflow | atcLobStorage | atcNestedTable

// atcPartitioningStorage is partitioning_storage_clause:
//
//	[ TABLESPACE tablespace | TABLESPACE SET tablespace_set ]
//	[ OVERFLOW [ TABLESPACE tablespace ] | TABLESPACE SET tablespace_set ]
//	[ table_compression | index_compression ] [ inmemory_clause ] [ ilm_clause ]
//	[ { LOB_partitioning_storage | VARRAY varray_item STORE AS ... }... ]
const atcPartitioningStorage = atcTablespace | atcOverflowTablespace | atcCompression |
	atcInmemory | atcIlm | atcLobStorage

// atcIndexMaintenance is update_index_clauses.
const atcIndexMaintenance = atcUpdateIndexes | atcInvalidate

// parseATClauses reads the clauses allowed lists, in any order, recording
// each as an option, and stops at the first token that starts none of them.
func (p *Parser) parseATClauses(allowed atClauses, opts *nodes.List) error {
	for {
		start := p.pos()
		ok, err := p.parseATClause(allowed, opts)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		p.atOption(opts, start)
	}
}

// atClauseAt returns the clause of allowed the current token starts, or 0.
func (p *Parser) atClauseAt(allowed atClauses) atClauses {
	has := func(c atClauses) bool { return allowed&c != 0 }
	switch {
	case has(atcSegment) && p.isATSegmentAttrStart():
		return atcSegment
	case has(atcPhysical) && (p.isATPhysicalAttrStart() || p.isATLoggingStart()):
		return atcPhysical
	case has(atcTablespace) && p.cur.Type == kwTABLESPACE:
		return atcTablespace
	case has(atcCompression) && p.isATTableCompressionStart():
		return atcCompression
	case has(atcInmemory) && p.isATInmemoryStart():
		return atcInmemory
	case has(atcIlm) && p.atWord("ILM"):
		return atcIlm
	case has(atcReadOnly) && p.cur.Type == kwREAD:
		return atcReadOnly
	case has(atcIndexing) && p.atWord("INDEXING"):
		return atcIndexing
	case has(atcSegmentCreation) && p.atWord("SEGMENT"):
		return atcSegmentCreation
	case has(atcOverflowTablespace) && p.atWord("OVERFLOW"):
		return atcOverflowTablespace
	case has(atcOverflow) && p.atWord("OVERFLOW"):
		return atcOverflow
	case has(atcLobModify) && p.atWord("LOB"):
		return atcLobModify
	case has(atcLobStorage) && p.isATLobStorageStart():
		return atcLobStorage
	case has(atcNestedTable) && p.cur.Type == kwNESTED:
		return atcNestedTable
	case has(atcExternal) && p.isATExternalDataPropStart():
		return atcExternal
	case has(atcUpdateIndexes) && p.cur.Type == kwUPDATE:
		return atcUpdateIndexes
	case has(atcInvalidate) && p.atWord("INVALIDATE"):
		return atcUpdateIndexes
	case has(atcParallel) && (p.cur.Type == kwPARALLEL || p.cur.Type == kwNOPARALLEL):
		return atcParallel
	case has(atcClustering) && (p.atWord("ALLOW") || p.atWord("DISALLOW")):
		// Oracle reads ALLOW or DISALLOW as the clause and wants CLUSTERING
		// after it (ORA-02000 at the end of MOVE PARTITION p ALLOW).
		return atcClustering
	case has(atcOnline) && p.cur.Type == kwONLINE:
		return atcOnline
	case has(atcFilter) && p.atWord("INCLUDING") && p.peekNext().Type == kwROWS:
		return atcFilter
	case has(atcSplitFilter) && p.atWord("INCLUDING"):
		return atcFilter
	case has(atcDependentTables) && p.atWord("DEPENDENT"):
		return atcDependentTables
	case has(atcMappingTable) && p.atWord("MAPPING"):
		return atcMappingTable
	case has(atcAllocate) && p.isATAllocateExtentStart():
		return atcAllocate
	case has(atcDeallocate) && p.isATDeallocateStart():
		return atcDeallocate
	case has(atcShrink) && p.isATShrinkStart():
		return atcShrink
	case has(atcCascade) && p.cur.Type == kwCASCADE:
		return atcCascade
	case has(atcExchange) && (p.atWord("INCLUDING") || p.atWord("EXCLUDING") ||
		p.cur.Type == kwWITH || p.atWord("WITHOUT") || p.atWord("EXCEPTIONS")):
		return atcExchange
	}
	return 0
}

// parseATClause reads one of the clauses allowed lists, or reports false
// when the current token starts none of them.
func (p *Parser) parseATClause(allowed atClauses, opts *nodes.List) (bool, error) {
	switch p.atClauseAt(allowed) {
	case 0:
		return false, nil
	case atcSegment:
		return true, p.parseATSegmentAttr()
	case atcPhysical:
		if p.isATLoggingStart() {
			p.advance()
			return true, nil
		}
		return true, p.parseATPhysicalAttr()
	case atcTablespace:
		p.advance()
		if p.cur.Type == kwSET {
			p.advance()
		}
		err := p.parseATName()
		return true, err
	case atcCompression:
		return true, p.parseATTableCompression(true)
	case atcInmemory:
		return true, p.parseATInmemoryClause()
	case atcIlm:
		return true, p.parseATIlmClause()
	case atcReadOnly:
		p.advance()
		if p.cur.Type == kwWRITE {
			p.advance()
			return true, nil
		}
		return true, p.atExpectWord("ONLY")
	case atcIndexing:
		p.advance()
		if p.cur.Type == kwON {
			p.advance()
			return true, nil
		}
		return true, p.atExpectWord("OFF")
	case atcSegmentCreation:
		p.advance()
		if err := p.atExpectWord("CREATION"); err != nil {
			return true, err
		}
		if p.cur.Type == kwIMMEDIATE || p.cur.Type == kwDEFERRED {
			p.advance()
			return true, nil
		}
		return true, p.syntaxErrorAtCur()
	case atcOverflowTablespace:
		p.advance()
		if p.cur.Type == kwTABLESPACE {
			p.advance()
			err := p.parseATName()
			return true, err
		}
		return true, nil
	case atcOverflow:
		p.advance()
		if p.isATSegmentAttrStart() {
			return true, p.parseATSegmentAttrs(nil)
		}
		return true, nil
	case atcLobStorage:
		return true, p.parseATLobStorage()
	case atcLobModify:
		return true, p.parseATLobModify()
	case atcNestedTable:
		return true, p.parseATNestedTableStorage()
	case atcExternal:
		return true, p.parseATExternalDataProp()
	case atcUpdateIndexes:
		return true, p.parseATUpdateIndexClauses()
	case atcParallel:
		return true, p.parseATParallel()
	case atcClustering:
		p.advance() // ALLOW / DISALLOW
		return true, p.atExpectWord("CLUSTERING")
	case atcOnline, atcCascade:
		p.advance()
		return true, nil
	case atcFilter:
		return true, p.parseATFilterCondition()
	case atcDependentTables:
		return true, p.parseATDependentTables()
	case atcMappingTable:
		p.advance()
		return true, p.atExpect(kwTABLE)
	case atcAllocate:
		return true, p.parseATAllocateExtent()
	case atcDeallocate:
		return true, p.parseATDeallocateUnused()
	case atcShrink:
		return true, p.parseATShrink()
	case atcExchange:
		switch {
		case p.atWord("INCLUDING"), p.atWord("EXCLUDING"):
			p.advance()
			return true, p.atExpectWord("INDEXES")
		case p.cur.Type == kwWITH, p.atWord("WITHOUT"):
			p.advance()
			return true, p.atExpectWord("VALIDATION")
		}
		p.advance() // EXCEPTIONS
		if err := p.atExpect(kwINTO); err != nil {
			return true, err
		}
		err := p.parseATQualifiedName()
		return true, err
	}
	return false, nil
}

// parseATParallel reads parallel_clause: { NOPARALLEL | PARALLEL [ integer ] }.
func (p *Parser) parseATParallel() error {
	if p.cur.Type == kwPARALLEL {
		p.advance()
		if p.cur.Type == tokICONST {
			p.advance()
		}
		return nil
	}
	p.advance() // NOPARALLEL
	return nil
}

// parseATNestedTableStorage reads nested_table_col_properties:
//
//	NESTED TABLE { nested_item | COLUMN_VALUE } [ substitutable_column_clause ]
//	    [ GLOBAL ] STORE AS storage_table
//	    [ ( { ( object_properties ) | physical_properties | column_properties }... ) ]
//	    [ RETURN [ AS ] { LOCATOR | VALUE } ]
//
// column_properties there name the storage of an inner collection level.
func (p *Parser) parseATNestedTableStorage() error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	p.advance() // NESTED
	if err := p.atExpect(kwTABLE); err != nil {
		return err
	}
	itemTok := p.cur
	if p.cur.Type == kwCOLUMN_VALUE {
		p.advance()
	} else if err := p.parseATItemPath(); err != nil {
		return err
	}
	if err := p.parseATSubstitutable(nil); err != nil {
		return err
	}
	p.atOptOneOf("GLOBAL")
	if err := p.atExpectWord("STORE"); err != nil {
		return err
	}
	if err := p.atExpect(kwAS); err != nil {
		return err
	}
	if err := p.parseATName(); err != nil {
		return err
	}
	if p.cur.Type == '(' && p.upgradeIncludingData {
		// Oracle answers ORA-01735 at byte 0 for UPGRADE [ NOT ] INCLUDING
		// DATA with a nested table's storage, on a table that holds the
		// column as on one that does not, and reads no item of the group:
		// ((a) ZZZ) and ((a) SEGMENT CREATION IMMEDIATE) answer alike, an
		// unbalanced group blames the nested item token, and a token after
		// the group ORA-00922. So the group is text, and the page's items
		// are held where an ordinary ALTER TABLE reads it.
		if err := p.skipATNestedStorageGroup(); err != nil {
			return p.syntaxErrorAtTok(itemTok)
		}
	} else if p.cur.Type == '(' {
		p.advance()
		var orgSeen, clusterSeen, objectSeen, nonObjectSeen bool
		phase := -1
		for n := 0; ; n++ {
			if p.cur.Type == ')' && n > 0 {
				p.advance()
				break
			}
			start := p.pos()
			if p.cur.Type == '(' {
				if objectSeen || nonObjectSeen {
					return p.syntaxErrorAtCur()
				}
				objectSeen = true
				if err := p.parseATObjectProperties(); err != nil {
					return err
				}
				continue
			}
			nonObjectSeen = true
			if p.atWord("ORGANIZATION") {
				if clusterSeen || orgSeen || phase > 4 {
					return p.syntaxErrorAtCur()
				}
				orgSeen = true
				phase = 5
				p.advance()
				if err := p.atOneOf("HEAP", "INDEX"); err != nil {
					return err
				}
				continue
			}
			// The rest of physical_properties that ALTER TABLE's clauses
			// do not share: an index-organized table's mapping table and
			// threshold, a heap table's cache hint, and the cluster its
			// storage table belongs to.
			if p.atWord("MAPPING") {
				if clusterSeen || orgSeen || phase > 4 {
					return p.syntaxErrorAtCur()
				}
				orgSeen = true
				phase = 5
				p.advance()
				if err := p.atExpect(kwTABLE); err != nil {
					return err
				}
				continue
			}
			if p.atWord("NOMAPPING") {
				if clusterSeen || orgSeen || phase > 4 {
					return p.syntaxErrorAtCur()
				}
				orgSeen = true
				phase = 5
				p.advance()
				continue
			}
			if p.atWord("PCTTHRESHOLD") {
				if clusterSeen || orgSeen || phase > 4 {
					return p.syntaxErrorAtCur()
				}
				orgSeen = true
				phase = 5
				p.advance()
				if err := p.parseATInteger(); err != nil {
					return err
				}
				continue
			}
			if p.atWord("CACHE") || p.atWord("NOCACHE") {
				if clusterSeen || orgSeen || phase > 4 {
					return p.syntaxErrorAtCur()
				}
				orgSeen = true
				phase = 5
				p.advance()
				continue
			}
			if p.atWord("INCLUDING") {
				ahead := p.peekAhead(2)
				if !isGrantWordToken(ahead[0]) || !p.isIdentLikeStrAt(ahead[1], "OVERFLOW") {
					return p.syntaxErrorAtCur()
				}
				if clusterSeen || orgSeen || phase > 4 {
					return p.syntaxErrorAtCur()
				}
				orgSeen = true
				phase = 5
				p.advance() // INCLUDING
				if err := p.parseATName(); err != nil {
					return err
				}
				if err := p.atExpectWord("OVERFLOW"); err != nil {
					return err
				}
				if p.isATSegmentAttrStart() {
					if err := p.parseATSegmentAttrs(nil); err != nil {
						return err
					}
				}
				continue
			}
			if p.cur.Type == kwCLUSTER {
				if clusterSeen || phase > 5 {
					return p.syntaxErrorAtCur()
				}
				clusterSeen = true
				phase = 6
				p.advance()
				if err := p.parseATName(); err != nil {
					return err
				}
				if err := p.parseATNameList(); err != nil {
					return err
				}
				continue
			}
			if clusterSeen && (p.isATSegmentAttrStart() || p.isATTableCompressionStart() ||
				p.isATInmemoryStart() || p.atWord("SEGMENT") || p.atWord("OVERFLOW")) {
				return p.syntaxErrorAtCur()
			}
			const physical = atcSegment | atcCompression | atcInmemory | atcSegmentCreation | atcOverflow
			for {
				rank := p.nestedPhysicalPhase(physical, phase)
				if rank < 0 {
					break
				}
				if rank < phase {
					return p.syntaxErrorAtCur()
				}
				switch rank {
				case 0, 1:
					phase = 1
				case 2:
					phase = 3
				case 3:
					phase = 4
				case 4:
					if orgSeen {
						return p.syntaxErrorAtCur()
					}
					orgSeen = true
					phase = 5
				}
				var parsed bool
				var err error
				if rank == 4 && (p.cur.Type == kwCOMPRESS || p.cur.Type == kwNOCOMPRESS) {
					err = p.parseATIndexCompression()
					parsed = true
				} else {
					parsed, err = p.parseATClause(physical, nil)
				}
				if err != nil {
					return err
				}
				if !parsed {
					return p.syntaxErrorAtCur()
				}
				if rank == 3 {
					for p.atWord("INMEMORY") || (p.atWord("NO") && p.atNextWord("INMEMORY")) {
						if err := p.parseATInmemoryColumnItem(); err != nil {
							return err
						}
					}
				}
			}
			columnStart := p.pos()
			if err := p.parseATColumnProperties(nil); err != nil {
				return err
			}
			if p.pos() == start {
				return p.syntaxErrorAtCur()
			}
			if p.pos() != columnStart && phase < 6 {
				phase = 6
			}
		}
	}
	if p.atWord("RETURN") {
		p.advance()
		if p.cur.Type == kwAS {
			p.advance()
		}
		return p.atOneOf("LOCATOR", "VALUE")
	}
	return nil
}

// nestedPhysicalPhase ranks a physical_properties clause in the page's
// order: deferred segment creation, segment attributes, compression,
// inmemory, then an organization clause.
func (p *Parser) nestedPhysicalPhase(allowed atClauses, phase int) int {
	if p.cur.Type == kwCOMPRESS || p.cur.Type == kwNOCOMPRESS {
		if p.cur.Type == kwCOMPRESS && p.peekNext().Type == tokICONST {
			// COMPRESS integer is an index-organized table's prefix
			// compression, the alternative to its mapping table and
			// threshold.
			return 4
		}
		if phase >= 3 {
			// After table compression or INMEMORY, COMPRESS or NOCOMPRESS
			// is the organization's prefix-compression alternative.
			return 4
		}
	}
	switch p.atClauseAt(allowed) {
	case atcSegmentCreation:
		return 0
	case atcSegment:
		return 1
	case atcCompression:
		return 2
	case atcInmemory:
		return 3
	case atcOverflow:
		return 4
	}
	return -1
}

// parseATObjectProperties reads ( object_properties ) of a nested table's
// storage table: ( { column [ DEFAULT expr ] [ inline_constraint ... ]
// [ inline_ref_constraint ] | out_of_line_constraint
// | out_of_line_ref_constraint | supplemental_logging_props } [, ...] ).
//
// An ordinary ALTER TABLE holds it to the page's items. Under
// upgrade_table_clause the whole group is text, its object_properties with
// it (parseATNestedTableStorage).
func (p *Parser) parseATObjectProperties() error {
	p.advance() // (
	for {
		switch {
		case p.isOutOfLineRefConstraintStart():
			tc, err := p.parseOutOfLineRefConstraint()
			if err != nil {
				return err
			}
			if tc == nil {
				return p.syntaxErrorAtCur()
			}
		case p.isTableConstraintStart():
			tc, err := p.parseTableConstraint()
			if err != nil {
				return err
			}
			if tc == nil {
				return p.syntaxErrorAtCur()
			}
		case p.atWord("SUPPLEMENTAL") && p.atNextWord("LOG"):
			p.advance()
			if err := p.parseATSupplementalLogging(true, true, nil); err != nil {
				return err
			}
		default:
			if err := p.parseATName(); err != nil {
				return err
			}
			if p.cur.Type == kwDEFAULT {
				p.advance()
				e, err := p.parseExpr()
				if err != nil {
					return err
				}
				if e == nil {
					return p.syntaxErrorAtCur()
				}
			}
			scopeSeen := false
			for p.isATObjectInlineConstraintStart() {
				if err := p.parseATObjectInlineConstraint(&scopeSeen); err != nil {
					return err
				}
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == ')' {
			return p.syntaxErrorAtCur()
		}
	}
	return p.atExpect(')')
}

// skipATNestedStorageGroup reads the balanced text of the parenthesized
// group a nested table's storage takes after STORE AS storage_table, which
// Oracle does not judge under upgrade_table_clause (parseATNestedTableStorage):
// the tokens up to the group's own ')'.
func (p *Parser) skipATNestedStorageGroup() error {
	p.advance() // (
	for depth := 1; depth > 0; {
		switch p.cur.Type {
		case tokEOF:
			return p.syntaxErrorAtCur()
		case ';':
			// A ';' ends the statement, so it cannot stand in the group.
			return p.syntaxErrorAtCur()
		case '(':
			depth++
		case ')':
			depth--
		}
		p.advance()
	}
	return nil
}

// isATObjectInlineConstraintStart reports an item of object_properties'
// inline_constraint or inline_ref_constraint after its column.
func (p *Parser) isATObjectInlineConstraintStart() bool {
	switch p.cur.Type {
	case kwCONSTRAINT, kwNOT, kwNULL, kwUNIQUE, kwPRIMARY, kwCHECK, kwREFERENCES, kwWITH:
		return true
	}
	return p.atWord("SCOPE")
}

// parseATObjectInlineConstraint reads one inline_constraint or
// inline_ref_constraint of object_properties. Its column has no declared
// type here, so the REF-only checks of parseInlineRefConstraint do not
// apply: SCOPE IS and WITH ROWID are read as the page gives them.
func (p *Parser) parseATObjectInlineConstraint(scopeSeen *bool) error {
	named := p.cur.Type == kwCONSTRAINT
	if named {
		p.advance()
		if !p.isIdentLike() {
			return p.syntaxErrorAtCur()
		}
		if err := p.parseATName(); err != nil {
			return err
		}
	}
	if named && (p.cur.Type == kwWITH || p.atWord("SCOPE")) {
		// SCOPE IS and WITH ROWID have no CONSTRAINT name in the page's
		// inline_ref_constraint.
		return p.syntaxErrorAtCur()
	}
	switch {
	case p.cur.Type == kwNOT:
		p.advance()
		if err := p.atExpect(kwNULL); err != nil {
			return err
		}
		var cs constraintState
		return p.parseConstraintState(&cs)
	case p.cur.Type == kwNULL:
		p.advance()
		var cs constraintState
		return p.parseConstraintState(&cs)
	case p.cur.Type == kwWITH:
		p.advance()
		return p.atExpect(kwROWID)
	case p.atWord("SCOPE"):
		if *scopeSeen {
			return p.syntaxErrorAtCur()
		}
		*scopeSeen = true
		p.advance()
		if err := p.atExpect(kwIS); err != nil {
			return err
		}
		return p.parseATQualifiedName()
	}
	cc, err := p.parseColumnConstraintInline()
	if err != nil {
		return err
	}
	if cc == nil {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// isATExternalDataPropStart reports an item of external_table_data_props or
// external_part_subpart_data_props: DEFAULT DIRECTORY or LOCATION.
func (p *Parser) isATExternalDataPropStart() bool {
	return (p.cur.Type == kwDEFAULT && p.atNextWord("DIRECTORY")) ||
		(p.atWord("LOCATION") && p.peekNext().Type == '(')
}

// parseATExternalDataProp reads one of
//
//	DEFAULT DIRECTORY directory
//	LOCATION ( [ directory: ] 'location_specifier' [, ...] )
func (p *Parser) parseATExternalDataProp() error {
	if p.cur.Type == kwDEFAULT {
		p.advance()
		if err := p.atExpectWord("DIRECTORY"); err != nil {
			return err
		}
		return p.parseATName()
	}
	p.advance() // LOCATION
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if p.cur.Type != tokSCONST {
			if err := p.parseATName(); err != nil {
				return err
			}
			if err := p.atExpect(':'); err != nil {
				return err
			}
		}
		if err := p.atExpect(tokSCONST); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATFilterCondition reads filter_condition: INCLUDING ROWS WHERE condition.
func (p *Parser) parseATFilterCondition() error {
	p.advance() // INCLUDING
	if err := p.atExpect(kwROWS); err != nil {
		return err
	}
	if err := p.atExpect(kwWHERE); err != nil {
		return err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return err
	}
	if cond == nil {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseATDependentTables reads dependent_tables_clause:
//
//	DEPENDENT TABLES ( table ( partition_spec [, partition_spec ]... )
//	    [, table ( partition_spec [, partition_spec ]... ) ]... )
func (p *Parser) parseATDependentTables() error {
	p.advance() // DEPENDENT
	if err := p.atExpectWord("TABLES"); err != nil {
		return err
	}
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if err := p.parseATName(); err != nil {
			return err
		}
		if err := p.atExpect('('); err != nil {
			return err
		}
		for {
			if err := p.parseATPartitionSpec(nil); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if err := p.atExpect(')'); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// isATUpdateIndexStart reports update_index_clauses: UPDATE GLOBAL INDEXES,
// INVALIDATE GLOBAL INDEXES or UPDATE INDEXES.
func (p *Parser) isATUpdateIndexStart() bool {
	return p.cur.Type == kwUPDATE || p.atWord("INVALIDATE")
}

// parseATUpdateIndexClauses reads update_index_clauses:
//
//	{ UPDATE | INVALIDATE } GLOBAL INDEXES
//	| UPDATE INDEXES [ ( index ( index_partition_description [, ...] )
//	                     [, index ( ... ) ]... ) ]
//
// Oracle's parser takes only index_partition_description in the list: an
// index_subpartition_clause there is refused with ORA-14004 (missing
// PARTITION keyword) wherever the clause stands.
func (p *Parser) parseATUpdateIndexClauses() error {
	invalidate := p.atWord("INVALIDATE")
	p.advance()
	if invalidate || p.atWord("GLOBAL") {
		if err := p.atExpectWord("GLOBAL"); err != nil {
			return err
		}
		return p.atExpectWord("INDEXES")
	}
	if err := p.atExpectWord("INDEXES"); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	for {
		if err := p.parseATName(); err != nil {
			return err
		}
		if err := p.atExpect('('); err != nil {
			return err
		}
		for {
			if err := p.parseATIndexPartitionDescription(); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if err := p.atExpect(')'); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATIndexPartitionDescription reads index_partition_description:
//
//	PARTITION [ partition [ { segment_attributes_clause | index_compression }...
//	                      | PARAMETERS ( 'ODCI_parameters' ) ] ]
//	    [ USABLE | UNUSABLE ]
func (p *Parser) parseATIndexPartitionDescription() error {
	if err := p.atExpect(kwPARTITION); err != nil {
		return err
	}
	var nameTok Token
	hasName := false
	if p.isATNameStart() && !p.isATSegmentAttrStart() && !p.isATIndexCompressionStart() &&
		!p.atWord("USABLE") && !p.atWord("UNUSABLE") {
		nameTok = p.cur
		hasName = true
		p.advance()
		if p.atWord("PARAMETERS") {
			p.advance()
			if err := p.atExpect('('); err != nil {
				return err
			}
			if err := p.atExpect(tokSCONST); err != nil {
				return err
			}
			if err := p.atExpect(')'); err != nil {
				return err
			}
		}
	}
	afterNamePos := p.pos()
	firstAttrTok := p.cur
	for {
		switch {
		case p.isATSegmentAttrStart():
			if err := p.parseATSegmentAttr(); err != nil {
				return err
			}
			continue
		case p.isATIndexCompressionStart():
			if err := p.parseATIndexCompression(); err != nil {
				return err
			}
			continue
		}
		break
	}
	p.atOptOneOf("USABLE", "UNUSABLE")
	if p.cur.Type != ',' && p.cur.Type != ')' {
		if hasName && p.pos() == afterNamePos {
			return p.syntaxErrorAtTok(nameTok)
		}
		if !hasName && (firstAttrTok.Type == kwCOMPRESS || firstAttrTok.Type == kwNOCOMPRESS) {
			return p.syntaxErrorAtTok(firstAttrTok)
		}
	}
	return nil
}

// parseATRangeValues reads range_values_clause:
// VALUES LESS THAN ( { value | MAXVALUE } [, { value | MAXVALUE } ]... ).
func (p *Parser) parseATRangeValues() error {
	if err := p.atExpect(kwVALUES); err != nil {
		return err
	}
	if err := p.atExpectWord("LESS"); err != nil {
		return err
	}
	if err := p.atExpectWord("THAN"); err != nil {
		return err
	}
	return p.parseATBoundValues(true)
}

// parseATBoundValues reads the values of a partition bound,
// ( value [, value ]... ), allowing MAXVALUE among them when maxvalue is
// set. Oracle takes a literal there, signed or not, NULL, a datetime or
// interval literal, or a function of literals such as TO_DATE: a name is
// refused (ORA-14019) and so is an operator (ORA-00907 at 1+1).
func (p *Parser) parseATBoundValues(maxvalue bool) error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if maxvalue && p.atWord("MAXVALUE") {
			p.advance()
		} else if err := p.parseATBoundValue(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATBoundValue reads one value of a partition bound.
func (p *Parser) parseATBoundValue() error {
	if p.cur.Type == '-' || p.cur.Type == '+' {
		p.advance()
	}
	e, err := p.parsePrimary()
	if err != nil {
		return err
	}
	if e == nil {
		return p.syntaxErrorAtCur()
	}
	if _, ok := e.(*nodes.ColumnRef); ok {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseATListValues reads list_values_clause:
// VALUES ( { value [, value ]... | DEFAULT } ).
func (p *Parser) parseATListValues() error {
	if err := p.atExpect(kwVALUES); err != nil {
		return err
	}
	if p.cur.Type == '(' && p.peekNext().Type == kwDEFAULT {
		p.advance()
		p.advance()
		return p.atExpect(')')
	}
	return p.parseATBoundValues(false)
}

// parseATValuesClause reads range_values_clause or list_values_clause,
// whichever follows VALUES.
func (p *Parser) parseATValuesClause() error {
	if p.atNextWord("LESS") {
		return p.parseATRangeValues()
	}
	return p.parseATListValues()
}

// atMethod is a partitioning method, where the statement says which: a
// subpartition description then takes that method's values clause.
type atMethod int

const (
	// atMethodAny is a method the statement does not name (ADD PARTITION,
	// SET SUBPARTITION TEMPLATE): any description parses.
	atMethodAny atMethod = iota
	atMethodRange
	atMethodList
	atMethodHash
)

// parseATMethodValues reads the values clause method needs, if any:
// range_values_clause for RANGE, list_values_clause for LIST, none for HASH
// (ORA-14007 for VALUES without LESS under RANGE, ORA-00926 for no VALUES).
func (p *Parser) parseATMethodValues(method atMethod) error {
	switch method {
	case atMethodRange:
		return p.parseATRangeValues()
	case atMethodList:
		if p.cur.Type == kwVALUES && p.atNextWord("LESS") {
			p.advance()
			return p.syntaxErrorAtCur()
		}
		return p.parseATListValues()
	case atMethodHash:
		return nil
	case atMethodAny:
		// A template names no method: Oracle reads a description's values
		// with or without VALUES (compat.tsv's SUBPARTITION LESS THAN(1)
		// row), and takes none for a hash subpartition.
		if p.cur.Type == kwVALUES {
			return p.parseATValuesClause()
		}
		if p.alterMview && p.atSubpartitionTemplate && p.atWord("LESS") {
			p.advance()
			if err := p.atExpectWord("THAN"); err != nil {
				return err
			}
			return p.parseATBoundValues(true)
		}
		return nil
	}
	return nil
}

// isATPartitionNameHere reports whether the token after PARTITION or
// SUBPARTITION is the partition's name rather than a clause the operation
// allows (or VALUES), or one that may follow the partition in the
// operation: dependent_tables_clause, update_index_clauses and the like.
func (p *Parser) isATPartitionNameHere(allowed atClauses) bool {
	if !p.isATNameStart() || p.cur.Type == kwVALUES {
		return false
	}
	if p.alterMview && p.atSubpartitionTemplate && p.atWord("LESS") {
		// A materialized view's template takes range values without VALUES
		// (compat.tsv's SUBPARTITION LESS THAN(1) row).
		return false
	}
	// A word that ends the partition is its name (ADD PARTITION INDEXING).
	switch p.peekNext().Type {
	case ';', tokEOF, ',', ')':
		return true
	}
	const follow = atcDependentTables | atcIndexMaintenance | atcParallel | atcFilter | atcClustering
	return p.atClauseAt(allowed|follow) == 0
}

// parseATPartitionSpec reads partition_spec:
// PARTITION [ partition ] [ table_partition_description ], or, where values
// is not nil, PARTITION [ partition ] [ range_values_clause |
// list_values_clause ] [ table_partition_description ], setting *values to
// whether the values came.
func (p *Parser) parseATPartitionSpec(values *bool) error {
	if err := p.atExpect(kwPARTITION); err != nil {
		return err
	}
	var nameTok Token
	hasName := false
	if p.isATPartitionNameHere(atcTablePartitionDescription) {
		nameTok = p.cur
		hasName = true
		p.advance()
	}
	afterNamePos := p.pos()
	if values != nil && p.cur.Type == kwVALUES {
		if err := p.parseATValuesClause(); err != nil {
			return err
		}
		*values = true
	}
	if err := p.parseATClauses(atcTablePartitionDescription, nil); err != nil {
		return err
	}
	const follow = atcSplitMerge | atcSplitFilter
	if hasName && p.pos() == afterNamePos && p.cur.Type != ',' && p.cur.Type != ')' &&
		p.cur.Type != ';' && p.cur.Type != tokEOF && p.atClauseAt(follow) == 0 {
		return p.syntaxErrorAtTok(nameTok)
	}
	return nil
}

// parseATSubpartitionSpec reads subpartition_spec:
// SUBPARTITION [ subpartition ] [ partitioning_storage_clause ], or, where
// values is not nil, with [ range_values_clause | list_values_clause ]
// after the name, setting *values to whether the values came.
func (p *Parser) parseATSubpartitionSpec(values *bool) error {
	if err := p.atExpect(kwSUBPARTITION); err != nil {
		return err
	}
	var nameTok Token
	hasName := false
	if p.isATPartitionNameHere(atcPartitioningStorage) {
		nameTok = p.cur
		hasName = true
		p.advance()
	}
	afterNamePos := p.pos()
	if values != nil && p.cur.Type == kwVALUES {
		if err := p.parseATValuesClause(); err != nil {
			return err
		}
		*values = true
	}
	if err := p.parseATClauses(atcPartitioningStorage, nil); err != nil {
		return err
	}
	const follow = atcSplitMerge | atcSplitFilter
	if hasName && p.pos() == afterNamePos && p.cur.Type != ',' && p.cur.Type != ')' &&
		p.cur.Type != ';' && p.cur.Type != tokEOF && p.atClauseAt(follow) == 0 {
		return p.syntaxErrorAtTok(nameTok)
	}
	return nil
}

// parseATSubpartitionDesc reads one of range_subpartition_desc,
// list_subpartition_desc or individual_hash_subparts:
//
//	SUBPARTITION [ subpartition ] [ range_values_clause | list_values_clause ]
//	    [ read_only_clause ] [ indexing_clause ] [ partitioning_storage_clause ]
//	    [ external_part_subpart_data_props ]
//
// Which one stands is the subpartitioning method's: method, where the
// statement names it, or the table's, decided after the parse.
func (p *Parser) parseATSubpartitionDesc(method atMethod) error {
	if err := p.atExpect(kwSUBPARTITION); err != nil {
		return err
	}
	const allowed = atcReadOnly | atcIndexing | atcPartitioningStorage | atcExternal
	var nameTok Token
	hasName := false
	if p.isATPartitionNameHere(allowed) {
		nameTok = p.cur
		hasName = true
		p.advance()
	}
	afterNamePos := p.pos()
	if p.alterMview && p.atSubpartitionTemplate && p.cur.Type == '(' && p.peekNext().Type != ')' {
		// Oracle reads a parenthesized list where a hash subpartition has
		// no values (compat.tsv's SUBPARTITION(1) row).
		depth := 0
		for {
			switch p.cur.Type {
			case tokEOF, ';':
				return p.syntaxErrorAtCur()
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					p.advance()
					return p.parseATClauses(allowed, nil)
				}
			}
			p.advance()
		}
	}
	p.atLastSubpartitionMethod = atMethodAny
	if method == atMethodAny && p.cur.Type == kwVALUES {
		if p.atNextWord("LESS") {
			p.atLastSubpartitionMethod = atMethodRange
		} else if p.peekNext().Type == '(' {
			p.atLastSubpartitionMethod = atMethodList
		}
	}
	if err := p.parseATMethodValues(method); err != nil {
		return err
	}
	if p.alterMview && p.atSubpartitionTemplate && p.isIdentLike() && (p.peekNext().Type == ',' || p.peekNext().Type == ')') {
		// Oracle reads a trailing word after the values as an acceptance
		// (compat.tsv's SUBPARTITION VALUES LESS THAN(1) zzz row).
		p.advance()
	}
	if err := p.parseATClauses(allowed, nil); err != nil {
		return err
	}
	const follow = atcDependentTables | atcIndexMaintenance | atcParallel
	if hasName && p.pos() == afterNamePos && p.cur.Type != ',' && p.cur.Type != ')' &&
		p.cur.Type != ';' && p.cur.Type != tokEOF && p.atClauseAt(follow) == 0 {
		return p.syntaxErrorAtTok(nameTok)
	}
	return nil
}

// parseATSubpartitionDescList reads ( subpartition_desc [, subpartition_desc ]... ).
func (p *Parser) parseATSubpartitionDescList(method atMethod) error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if err := p.parseATSubpartitionDesc(method); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATHashSubpartsByQuantity reads hash_subparts_by_quantity:
// SUBPARTITIONS integer [ STORE IN ( tablespace [, tablespace ]... ) ].
func (p *Parser) parseATHashSubpartsByQuantity() error {
	p.advance() // SUBPARTITIONS
	if err := p.parseATInteger(); err != nil {
		return err
	}
	return p.parseATStoreIn()
}

// parseATStoreIn reads [ STORE IN ( tablespace [, tablespace ]... ) ].
func (p *Parser) parseATStoreIn() error {
	if !p.atWord("STORE") {
		return nil
	}
	p.advance()
	if err := p.atExpect(kwIN); err != nil {
		return err
	}
	return p.parseATNameList()
}

// ---------------------------------------------------------------------------
// table_partitioning_clauses (modify_to_partitioned)
// ---------------------------------------------------------------------------

// parseATTablePartitioning reads table_partitioning_clauses from PARTITION
// BY:
//
//	{ range_partitions | list_partitions | hash_partitions
//	| composite_range_partitions | composite_list_partitions
//	| composite_hash_partitions | reference_partitioning
//	| system_partitioning | consistent_hash_partitions
//	| consistent_hash_with_subpartitions }
//
// partitionset_clauses and the consistent_hash_* clauses belong to
// composite-sharded and sharded tables: Oracle refuses them at PARTITIONSET
// (ORA-02577) and at what follows HASH (ORA-14151) before it reads the
// clause, so the switch reports those words. PARTITIONSET reaches the syntax
// error at the end of this switch.
//
// ALTER MATERIALIZED VIEW LOG and other statements that take these
// clauses can call it with the current token on PARTITION.
func (p *Parser) parseATTablePartitioning(opts *nodes.List) error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	start := p.pos()
	if err := p.atExpect(kwPARTITION); err != nil {
		return err
	}
	if err := p.atExpect(kwBY); err != nil {
		return err
	}
	switch {
	case p.atWord("RANGE"), p.atWord("LIST"):
		list := p.atWord("LIST")
		p.advance()
		if err := p.parseATPartitionColumnList(); err != nil {
			return err
		}
		if !list && p.cur.Type == kwINTERVAL {
			p.advance()
			if err := p.parseATExprList(); err != nil {
				return err
			}
			if err := p.parseATStoreIn(); err != nil {
				return err
			}
		}
		if list && p.cur.Type == kwAUTOMATIC {
			p.advance()
			if err := p.parseATStoreIn(); err != nil {
				return err
			}
		}
		p.atOption(opts, start)
		method, sub := atMethodRange, atMethodAny
		if list {
			method = atMethodList
		}
		if p.cur.Type == kwSUBPARTITION {
			s := p.pos()
			if err := p.parseATSubpartitionBy(&sub); err != nil {
				return err
			}
			p.atOption(opts, s)
		}
		s := p.pos()
		if err := p.parseATValuePartitionList(method, sub); err != nil {
			return err
		}
		p.atOption(opts, s)
		return nil
	case p.atWord("HASH"):
		p.advance()
		if err := p.parseATPartitionColumnList(); err != nil {
			return err
		}
		p.atOption(opts, start)
		if p.cur.Type == kwSUBPARTITION {
			s := p.pos()
			var sub atMethod
			if err := p.parseATSubpartitionBy(&sub); err != nil {
				return err
			}
			p.atOption(opts, s)
		}
		s := p.pos()
		if err := p.parseATHashPartitions(); err != nil {
			return err
		}
		p.atOption(opts, s)
		return nil
	case p.cur.Type == kwREFERENCE:
		p.advance()
		if err := p.atExpect('('); err != nil {
			return err
		}
		if err := p.parseATName(); err != nil {
			return err
		}
		if err := p.atExpect(')'); err != nil {
			return err
		}
		if p.cur.Type == '(' {
			if err := p.parseATPartitionSpecList(); err != nil {
				return err
			}
		}
		p.atOption(opts, start)
		return nil
	case p.atWord("SYSTEM"):
		p.advance()
		if p.atWord("PARTITIONS") {
			p.advance()
			if err := p.parseATInteger(); err != nil {
				return err
			}
		} else if p.cur.Type == '(' {
			if err := p.parseATPartitionSpecList(); err != nil {
				return err
			}
		}
		p.atOption(opts, start)
		return nil
	case p.atWord("CONSISTENT"):
		// consistent_hash_partitions and consistent_hash_with_subpartitions
		// belong to sharded tables: Oracle reads CONSISTENT HASH and refuses
		// the form elsewhere (ORA-14151) at what follows HASH, before it
		// reads the columns.
		p.advance() // CONSISTENT
		if err := p.atExpectWord("HASH"); err != nil {
			return err
		}
		return p.syntaxErrorAtCur()
	}
	return p.syntaxErrorAtCur()
}

// parseATPartitionColumnList reads ( column [, column ]... ) of a PARTITION BY
// or SUBPARTITION BY clause. When a name is followed by '(' (e.g. a missing ')'
// before "(PARTITION VALUES ...)"), parseExpr enters the parenthesized group to
// match Oracle's error token while rejecting completed function-call expressions.
func (p *Parser) parseATPartitionColumnList() error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if !p.isATNameStart() {
			return p.syntaxErrorAtCur()
		}
		if p.peekNext().Type == '(' {
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			if expr != nil {
				return p.syntaxErrorAtCur()
			}
			return p.syntaxErrorAtCur()
		}
		if err := p.parseATName(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATPartitionSpecList reads ( PARTITION [ partition ]
// [ table_partition_description ] [, ...] ), the partitions of reference and
// system partitioning.
func (p *Parser) parseATPartitionSpecList() error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if err := p.parseATPartitionSpec(nil); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATSubpartitionBy reads subpartition_by_range, subpartition_by_list
// or subpartition_by_hash:
//
//	SUBPARTITION BY { RANGE | LIST } ( column [, column ]... ) [ subpartition_template ]
//	SUBPARTITION BY HASH ( column [, column ]... )
//	    [ SUBPARTITIONS integer [ STORE IN ( tablespace [, tablespace ]... ) ]
//	    | subpartition_template ]
//
// It sets *method to the subpartitioning method.
func (p *Parser) parseATSubpartitionBy(method *atMethod) error {
	p.advance() // SUBPARTITION
	if err := p.atExpect(kwBY); err != nil {
		return err
	}
	switch {
	case p.atWord("HASH"):
		*method = atMethodHash
	case p.atWord("RANGE"):
		*method = atMethodRange
	case p.atWord("LIST"):
		*method = atMethodList
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if err := p.parseATPartitionColumnList(); err != nil {
		return err
	}
	if *method == atMethodHash && p.atWord("SUBPARTITIONS") {
		return p.parseATHashSubpartsByQuantity()
	}
	if p.cur.Type == kwSUBPARTITION {
		return p.parseATSubpartitionTemplate(*method)
	}
	return nil
}

// parseATSubpartitionTemplate reads subpartition_template:
//
//	SUBPARTITION TEMPLATE
//	    { ( { range_subpartition_desc | list_subpartition_desc
//	        | individual_hash_subparts } [, ...] )
//	    | hash_subpartition_quantity }
//
// A quantity stands for hash subpartitions only (ORA-14603).
func (p *Parser) parseATSubpartitionTemplate(method atMethod) error {
	defer func(saved bool) { p.atSubpartitionTemplate = saved }(p.atSubpartitionTemplate)
	p.atSubpartitionTemplate = true
	p.advance() // SUBPARTITION
	if err := p.atExpectWord("TEMPLATE"); err != nil {
		return err
	}
	if p.alterMview {
		// Oracle reads a bare SUBPARTITION here, and a word where
		// SUBPARTITIONS stands, as acceptances on ALTER MATERIALIZED VIEW.
		if p.cur.Type == kwSUBPARTITION && p.peekNext().Type == ')' {
			p.advance()
			p.advance()
			return nil
		}
		if p.isIdentLike() && p.cur.Type != kwSUBPARTITION && p.peekNext().Type == tokICONST {
			p.advance()
			p.advance()
			return nil
		}
	}
	if p.cur.Type == tokICONST {
		p.advance()
		if method != atMethodAny && method != atMethodHash {
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	if p.cur.Type == '(' && p.peekNext().Type == ')' {
		// An empty template drops the table's template.
		p.advance()
		p.advance()
		return nil
	}
	return p.parseATSubpartitionDescList(method)
}

// parseATValuePartitionList reads the partitions of range or list
// partitioning, simple or composite:
//
//	( PARTITION [ partition ] { range_values_clause | list_values_clause }
//	    table_partition_description
//	    [ ( subpartition_desc [, ...] ) | hash_subparts_by_quantity ] [, ...] )
//
// method is RANGE or LIST, and sub the subpartitioning method of a
// composite method, atMethodAny for none: a quantity of subpartitions
// stands for hash subpartitions only (ORA-14639).
func (p *Parser) parseATValuePartitionList(method, sub atMethod) error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		if err := p.atExpect(kwPARTITION); err != nil {
			return err
		}
		if p.cur.Type != kwVALUES {
			if err := p.parseATName(); err != nil {
				return err
			}
		}
		if err := p.parseATMethodValues(method); err != nil {
			return err
		}
		if err := p.parseATClauses(atcTablePartitionDescription, nil); err != nil {
			return err
		}
		if sub != atMethodAny {
			switch {
			case p.cur.Type == '(':
				if err := p.parseATSubpartitionDescList(sub); err != nil {
					return err
				}
			case p.atWord("SUBPARTITIONS") && sub == atMethodHash:
				if err := p.parseATHashSubpartsByQuantity(); err != nil {
					return err
				}
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseATHashPartitions reads individual_hash_partitions or
// hash_partitions_by_quantity:
//
//	( PARTITION [ partition ] [ read_only_clause ] [ indexing_clause ]
//	    [ partitioning_storage_clause ] [, ...] )
//	PARTITIONS integer [ STORE IN ( tablespace [, tablespace ]... ) ]
//	    [ table_compression ] [ OVERFLOW STORE IN ( tablespace [, ...] ) ]
func (p *Parser) parseATHashPartitions() error {
	if p.atWord("PARTITIONS") {
		p.advance()
		if err := p.parseATInteger(); err != nil {
			return err
		}
		if p.atWord("STORE") {
			p.advance()
			if err := p.atExpect(kwIN); err != nil {
				return err
			}
			// Oracle takes a trailing comma in this list.
			if err := p.atExpect('('); err != nil {
				return err
			}
			for {
				if err := p.parseATName(); err != nil {
					return err
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
				if p.cur.Type == ')' {
					break
				}
			}
			if err := p.atExpect(')'); err != nil {
				return err
			}
		}
		return nil
	}
	if err := p.atExpect('('); err != nil {
		return err
	}
	const allowed = atcReadOnly | atcIndexing | atcPartitioningStorage
	for {
		if err := p.atExpect(kwPARTITION); err != nil {
			return err
		}
		var nameTok Token
		hasName := false
		if p.isATPartitionNameHere(allowed) {
			nameTok = p.cur
			hasName = true
			p.advance()
		}
		afterNamePos := p.pos()
		if err := p.parseATClauses(allowed, nil); err != nil {
			return err
		}
		if hasName && p.pos() == afterNamePos && p.cur.Type != ',' && p.cur.Type != ')' {
			return p.syntaxErrorAtTok(nameTok)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}
