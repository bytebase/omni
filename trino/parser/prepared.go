package parser

import (
	"sort"
	"strings"
	"unicode"

	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-dcl-tcl DAG node: it implements Trino's
// prepared-statement family — PREPARE, DEALLOCATE PREPARE, EXECUTE, EXECUTE
// IMMEDIATE, DESCRIBE INPUT, and DESCRIBE OUTPUT — as hand-written
// recursive-descent parsers over the token stream.
//
// The statement nodes are returned from parseStmt as ast.Node and carry an
// ast.NodeTag (trino/ast/nodetags.go); their concrete fields live here in
// package parser (Trino convention for parser-built node types).
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	statement
//	    : ...
//	    | PREPARE_ identifier FROM_ statement                                   # prepare
//	    | DEALLOCATE_ PREPARE_ identifier                                       # deallocate
//	    | EXECUTE_ identifier (USING_ expression (COMMA_ expression)*)?         # execute
//	    | EXECUTE_ IMMEDIATE_ string_ (USING_ expression (COMMA_ expression)*)? # executeImmediate
//	    | DESCRIBE_ INPUT_ identifier                                           # describeInput
//	    | DESCRIBE_ OUTPUT_ identifier                                          # describeOutput
//	    ;
//
// Trino 481 docs (truth1) confirm the same surface.
//
// Boundary decision (recorded as flagged-by-design in the migration divergence
// ledger, following the doris/snowflake/expressions B1 precedent): the inner
// `statement` of PREPARE is captured as RAW TEXT, not a parsed AST sub-tree.
// The omni statement-layer parsers (parser-select/ddl/dml) are not all merged
// when this node lands, and even once they are, an embedded statement stays a
// placeholder analogous to a SubqueryExpr — a later consumer (analysis,
// deparse) re-parses the RawText if it needs structure. We still require the
// inner statement to be non-empty so `PREPARE x FROM` (which Trino rejects as a
// syntax error) is rejected here too.
//
// The implementation is adjudicated against the live Trino 481 oracle.

// ---------------------------------------------------------------------------
// Statement AST
// ---------------------------------------------------------------------------

// PrepareStmt is a PREPARE <name> FROM <statement>. Name is the prepared
// statement name; Body is the raw source text of the prepared statement (see
// the file-header boundary decision).
type PrepareStmt struct {
	Name *ast.Identifier
	Body string
	Loc  ast.Loc
}

// Tag implements ast.Node.
func (n *PrepareStmt) Tag() ast.NodeTag { return ast.T_PrepareStmt }

// DeallocateStmt is a DEALLOCATE PREPARE <name>.
type DeallocateStmt struct {
	Name *ast.Identifier
	Loc  ast.Loc
}

// Tag implements ast.Node.
func (n *DeallocateStmt) Tag() ast.NodeTag { return ast.T_DeallocateStmt }

// ExecuteStmt is an EXECUTE <name> [USING <expr> [, <expr> ...]]. Using holds
// the positional parameter expressions (nil/empty when the USING clause is
// absent).
type ExecuteStmt struct {
	Name  *ast.Identifier
	Using []Expr
	Loc   ast.Loc
}

// Tag implements ast.Node.
func (n *ExecuteStmt) Tag() ast.NodeTag { return ast.T_ExecuteStmt }

// ExecuteImmediateStmt is an EXECUTE IMMEDIATE '<sql>' [USING <expr> [, ...]].
// SQL is the decoded statement-string literal; Using holds the positional
// parameter expressions.
type ExecuteImmediateStmt struct {
	SQL   string
	Using []Expr
	Loc   ast.Loc
}

// Tag implements ast.Node.
func (n *ExecuteImmediateStmt) Tag() ast.NodeTag { return ast.T_ExecuteImmediateStmt }

// DescribeInputStmt is a DESCRIBE INPUT <name>.
type DescribeInputStmt struct {
	Name *ast.Identifier
	Loc  ast.Loc
}

// Tag implements ast.Node.
func (n *DescribeInputStmt) Tag() ast.NodeTag { return ast.T_DescribeInputStmt }

// DescribeOutputStmt is a DESCRIBE OUTPUT <name>.
type DescribeOutputStmt struct {
	Name *ast.Identifier
	Loc  ast.Loc
}

// Tag implements ast.Node.
func (n *DescribeOutputStmt) Tag() ast.NodeTag { return ast.T_DescribeOutputStmt }

// Compile-time assertions that the statement types satisfy ast.Node.
var (
	_ ast.Node = (*PrepareStmt)(nil)
	_ ast.Node = (*DeallocateStmt)(nil)
	_ ast.Node = (*ExecuteStmt)(nil)
	_ ast.Node = (*ExecuteImmediateStmt)(nil)
	_ ast.Node = (*DescribeInputStmt)(nil)
	_ ast.Node = (*DescribeOutputStmt)(nil)
)

// ---------------------------------------------------------------------------
// Parsers
// ---------------------------------------------------------------------------

// parsePrepareStmt parses PREPARE <name> FROM <statement>. On entry cur is the
// PREPARE keyword. The inner statement is captured as raw text (see the
// file-header boundary decision) but is also validated as a real Trino
// statement so that PREPARE rejects an invalid inner statement the way Trino
// does (`PREPARE p FROM 1` is a SYNTAX_ERROR).
//
// The body runs to the end of the segment (the statement splitter has
// already isolated one statement per segment): it is the text from its first
// token to the segment's last, trimmed. It is validated by parsing it in
// place, as a statement of its own (see parsePrepareBody). Parsing its text
// again instead did the same for a PREPARE inside it, and every level
// scanned the rest of the segment, so n nested PREPAREs cost O(n²).
func (p *Parser) parsePrepareStmt() (ast.Node, error) {
	prepareTok := p.advance() // consume PREPARE
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	if p.cur.Kind == tokEOF {
		// PREPARE x FROM <nothing> — Trino rejects this as a syntax error.
		return nil, p.syntaxErrorAtCur()
	}
	bodyStart := p.cur.Loc.Start
	bodyLoc := p.cur.Loc
	bodyEnd := p.prepareBodyEnd()
	raw := p.sourceSlice(bodyStart, bodyEnd)
	body := strings.TrimSpace(raw)
	if body == "" {
		p.skipToEOF()
		return nil, &ParseError{Position: bodyStart, End: bodyEnd, Message: "expected a statement after PREPARE ... FROM"}
	}
	// lead is the length of the leading text the trim dropped. A position in
	// the body is reported as its offset in the trimmed text from bodyStart,
	// lead bytes short of where it is.
	lead := len(raw) - len(strings.TrimLeftFunc(raw, unicode.IsSpace))
	trimStart := bodyStart + lead

	// Validate the inner statement. Parsing the body through the full parser
	// reports a real error for non-statements (`PREPARE p FROM 1`,
	// `... FROM NOTASTATEMENT`). A "not yet supported" stub error means the
	// inner statement FORM is recognized (its concrete parser is a later DAG
	// node) — that counts as syntactically valid here, and the check tightens
	// automatically as those parsers land.
	hard := p.parsePrepareBody(trimStart, trimStart+len(body))
	p.skipToEOF()
	if hard != nil {
		if hard.Message == nestingTooDeep {
			// The body nests too deeply counting the PREPAREs around it: the
			// statement as a whole does, so it fails with the nesting error
			// itself.
			e := *hard
			e.Position -= lead
			if e.End >= 0 {
				e.End -= lead
			}
			p.nestingErr = &e
			return nil, &e
		}
		return nil, p.invalidPrepareBody(bodyLoc, hard)
	}

	return &PrepareStmt{
		Name: name,
		Body: body,
		Loc:  ast.Loc{Start: prepareTok.Loc.Start, End: bodyEnd},
	}, nil
}

// prepareInvalid prefixes the message of a PREPARE whose body is invalid.
const prepareInvalid = "invalid statement after PREPARE ... FROM: "

// isStubError reports whether msg is a "not yet supported" stub error, which
// a PREPARE body may carry: its statement form was recognized by the
// dispatch, and its body parser is a later DAG node.
func isStubError(msg string) bool {
	return strings.Contains(msg, "not yet supported")
}

// parsePrepareBody parses the PREPARE body spanning [start, end), absolute
// offsets, and returns the first of its errors that is not a stub error, or
// nil when it has none. It reads the body as parsing its text on its own did:
// in strict mode, from the PREPARE's nesting depth, with the lexer input
// ending at end; and the body's errors are the ones that parse reported, in
// its order: the statement's error, or its first leftover token, or the
// errors of the query bodies it captured raw; then the body's lex errors.
// The raw query bodies are validated here and dropped, as that parse did, so
// the statement does not validate them again.
//
// The outermost body is lexed by a lexer of its own: the trim can end the
// body inside a token, an unterminated string ending in spaces, which the
// statement's lexer reads whole. On return the parser is back at the body's
// first token, for the caller to read the segment from there with the
// statement's lexer, once. A body in another one is lexed by the enclosing
// body's lexer, its input cut at end while the body is parsed (the trim
// cannot end that body inside a token, as the text it is cut from ends in
// no space); the parser is left at the first token from end on.
func (p *Parser) parsePrepareBody(start, end int) *ParseError {
	lexer := p.lexer
	savedStrict := p.strictTrailing
	rawMark := len(p.rawQueries)
	outermost := p.prepareBodies == 0
	if outermost {
		from := p.checkpoint()
		p.lexer = &Lexer{input: lexer.input[:end-lexer.baseOffset], pos: start - lexer.baseOffset, baseOffset: lexer.baseOffset}
		p.hasNext = false
		p.cur = p.lexer.NextToken()
		defer func() {
			p.lexer = lexer
			p.replayTo(from)
		}()
	} else {
		// The tokens the trim dropped belong to the enclosing body, not this
		// one: their lex errors are that body's.
		for p.cur.Kind != tokEOF && p.cur.Loc.Start < start {
			p.advance()
		}
		full := lexer.input
		lexer.input = full[:end-lexer.baseOffset]
		// Lex the body from its first token on in the cut input: a token
		// buffered past end is not the body's.
		lexer.pos = min(p.cur.Loc.End, end) - lexer.baseOffset
		p.hasNext = false
		if p.cur.Loc.Start >= end {
			p.cur = lexer.NextToken()
		}
		defer func() {
			lexer.input = full
			lexer.pos = end - lexer.baseOffset
			p.hasNext = false
			p.cur = lexer.NextToken()
		}()
	}
	p.prepareBodies++
	defer func() {
		p.strictTrailing = savedStrict
		p.rawQueries = p.rawQueries[:rawMark]
		p.prepareBodies--
	}()
	if p.cur.Kind == tokEOF {
		// The trimmed body holds comments only: the parse of its text found
		// no statement, and no error.
		return nil
	}
	p.strictTrailing = true

	first := p.cur
	node, err := p.parseStmt()
	var hard *ParseError
	switch {
	case p.nestingErr != nil:
		hard = p.nestingErr
	case err == nil && p.treeTooHigh(node, end-start):
		hard = &ParseError{Position: first.Loc.Start, End: first.Loc.End, Message: nestingTooDeep}
	case err != nil:
		pe, ok := err.(*ParseError)
		if !ok {
			pe = &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: err.Error()}
		}
		if !isStubError(pe.Message) {
			hard = pe
		}
	case p.cur.Kind != tokEOF:
		hard = p.syntaxErrorAtCur()
	default:
		for _, e := range p.validateRawQueries(p.rawQueries[rawMark:]) {
			if !isStubError(e.Message) {
				hard = &e
				break
			}
		}
	}
	// Lex the rest of the body: its lex errors are the enclosing body's too.
	p.skipToEOF()
	if hard != nil {
		return hard
	}
	// No error but a stub's: the first lex error of the body is the first
	// one left.
	bl := p.lexer
	errs := bl.errors[:bl.errorCount()]
	i := sort.Search(len(errs), func(i int) bool { return errs[i].Loc.Start+bl.baseOffset >= start })
	if i < len(errs) {
		le := errs[i]
		return &ParseError{Position: le.Loc.Start + bl.baseOffset, End: le.Loc.End + bl.baseOffset, Message: le.Msg}
	}
	return nil
}

// prepareBodyEnd returns where the PREPARE body starting at p.cur ends before
// it is trimmed: at the end of the last token of the text it is part of,
// which is the segment, or for a PREPARE in another one's body, that body.
// The segment's tokens from the outermost body on are listed once, and a
// PREPARE in a body ends at the last of them that the body holds, or at the
// body's end when the body ends inside a token (see parsePrepareBody): the
// token, unterminated, runs to the end of the text it is lexed in.
func (p *Parser) prepareBodyEnd() int {
	if p.prepareBodies == 0 {
		toks := append(p.prepareTokens[:0], p.cur.Loc)
		lx := &Lexer{input: p.lexer.input, pos: p.cur.Loc.End - p.lexer.baseOffset, baseOffset: p.lexer.baseOffset}
		for tok := lx.NextToken(); tok.Kind != tokEOF; tok = lx.NextToken() {
			toks = append(toks, tok.Loc)
		}
		p.prepareTokens = toks
		return toks[len(toks)-1].End
	}
	limit := p.lexer.baseOffset + len(p.lexer.input)
	// The last token starting before limit.
	last := p.prepareTokens[sort.Search(len(p.prepareTokens), func(i int) bool { return p.prepareTokens[i].Start >= limit })-1]
	return min(last.End, limit)
}

// skipToEOF consumes the rest of the segment.
func (p *Parser) skipToEOF() {
	for p.cur.Kind != tokEOF {
		p.advance()
	}
}

// invalidPrepareBody returns the error, at the body's first token, of a
// PREPARE whose body failed with hard: prepareInvalid and hard's message. A
// body that is a PREPARE failing that way adds the prefix once more, so the
// message of n nested PREPAREs holds it n times. A PREPARE in another one's
// body returns its error without the prefixes, counting them in
// prepareWraps, and the outermost writes them all at once: writing the
// message out at every level was quadratic.
func (p *Parser) invalidPrepareBody(bodyLoc ast.Loc, hard *ParseError) *ParseError {
	wraps := 1
	if n, ok := p.prepareWraps[hard]; ok {
		wraps += n
		delete(p.prepareWraps, hard)
	}
	e := &ParseError{Position: bodyLoc.Start, End: bodyLoc.End, Message: hard.Message}
	if p.prepareBodies > 0 {
		if p.prepareWraps == nil {
			p.prepareWraps = make(map[*ParseError]int)
		}
		p.prepareWraps[e] = wraps
		return e
	}
	e.Message = strings.Repeat(prepareInvalid, wraps) + hard.Message
	return e
}

// parseDeallocateStmt parses DEALLOCATE PREPARE <name>. On entry cur is the
// DEALLOCATE keyword.
func (p *Parser) parseDeallocateStmt() (ast.Node, error) {
	deallocTok := p.advance() // consume DEALLOCATE
	if _, err := p.expect(kwPREPARE); err != nil {
		return nil, err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &DeallocateStmt{
		Name: name,
		Loc:  ast.Loc{Start: deallocTok.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// parseExecuteStmt parses both EXECUTE forms, dispatching on a leading
// IMMEDIATE:
//
//	EXECUTE <name> [USING <expr> [, ...]]
//	EXECUTE IMMEDIATE '<sql>' [USING <expr> [, ...]]
//
// On entry cur is the EXECUTE keyword.
func (p *Parser) parseExecuteStmt() (ast.Node, error) {
	executeTok := p.advance() // consume EXECUTE
	start := executeTok.Loc.Start

	// EXECUTE IMMEDIATE '<sql>' is the immediate form ONLY when IMMEDIATE is
	// directly followed by a string literal. IMMEDIATE is a NON-RESERVED
	// keyword, so a bare `EXECUTE IMMEDIATE` (or `EXECUTE IMMEDIATE USING 1`)
	// is the named-execute form of a prepared statement named "immediate" —
	// Trino dispatches the same way (it reports a semantic
	// "Prepared statement not found: IMMEDIATE", i.e. accepts the syntax).
	if p.cur.Kind == kwIMMEDIATE {
		next := p.peekNext().Kind
		if next == tokString || next == tokUnicodeString {
			return p.parseExecuteImmediateTail(start)
		}
	}

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt := &ExecuteStmt{
		Name: name,
		Loc:  ast.Loc{Start: start, End: p.prev.Loc.End},
	}
	using, err := p.parseOptionalUsing()
	if err != nil {
		return nil, err
	}
	stmt.Using = using
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseExecuteImmediateTail parses the tail of EXECUTE IMMEDIATE:
//
//	IMMEDIATE '<sql>' [USING <expr> [, ...]]
//
// On entry cur is IMMEDIATE. start is the byte offset of the EXECUTE keyword.
func (p *Parser) parseExecuteImmediateTail(start int) (ast.Node, error) {
	p.advance() // consume IMMEDIATE

	// The grammar requires a string_ literal here (the SQL text). A non-string
	// token is a syntax error; in particular EXECUTE IMMEDIATE <ident> is
	// invalid (that is the named-EXECUTE form, which has no IMMEDIATE).
	if p.cur.Kind != tokString && p.cur.Kind != tokUnicodeString {
		return nil, p.syntaxErrorAtCur()
	}
	strTok := p.advance()

	stmt := &ExecuteImmediateStmt{
		SQL: strTok.Str,
		Loc: ast.Loc{Start: start, End: strTok.Loc.End},
	}
	using, err := p.parseOptionalUsing()
	if err != nil {
		return nil, err
	}
	stmt.Using = using
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseOptionalUsing parses an optional `USING <expr> (, <expr>)*` clause. It
// returns nil when no USING is present. A USING keyword with no following
// expression is a syntax error.
func (p *Parser) parseOptionalUsing() ([]Expr, error) {
	if _, ok := p.match(kwUSING); !ok {
		return nil, nil
	}
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	using := []Expr{first}
	for {
		if _, ok := p.match(int(',')); !ok {
			break
		}
		next, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		using = append(using, next)
	}
	return using, nil
}

// describeIntrospectionFollows reports whether the token two positions ahead of
// the current DESCRIBE keyword (i.e. the token after INPUT/OUTPUT) starts an
// identifier — the statement name of a DESCRIBE INPUT/OUTPUT introspection
// statement. It performs a bounded 3-token speculative scan via a checkpoint so
// the caller's cursor is unchanged; on entry cur is DESCRIBE and peekNext is
// INPUT/OUTPUT.
//
// When this returns false (INPUT/OUTPUT is the last token, or is followed by a
// '.' or any non-identifier), the statement is the SHOW COLUMNS alias
// `DESCRIBE qualifiedName` with INPUT/OUTPUT as the (non-reserved) table name,
// which the dispatcher leaves to parser-utility.
func (p *Parser) describeIntrospectionFollows() bool {
	cp := p.checkpoint()
	defer p.restore(cp)
	p.advance() // consume DESCRIBE
	p.advance() // consume INPUT / OUTPUT
	return isIdentifierStart(p.cur.Kind)
}

// parseDescribeInputStmt parses DESCRIBE INPUT <name>. On entry the DESCRIBE
// keyword and the INPUT keyword have BOTH already been consumed by the
// dispatcher (DESCRIBE is shared with the showColumns form, so the dispatcher
// peeks the second keyword); start is the DESCRIBE keyword's byte offset.
func (p *Parser) parseDescribeInputStmt(start int) (ast.Node, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &DescribeInputStmt{
		Name: name,
		Loc:  ast.Loc{Start: start, End: p.prev.Loc.End},
	}, nil
}

// parseDescribeOutputStmt parses DESCRIBE OUTPUT <name>. On entry the DESCRIBE
// and OUTPUT keywords have both already been consumed; start is the DESCRIBE
// keyword's byte offset.
func (p *Parser) parseDescribeOutputStmt(start int) (ast.Node, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &DescribeOutputStmt{
		Name: name,
		Loc:  ast.Loc{Start: start, End: p.prev.Loc.End},
	}, nil
}
