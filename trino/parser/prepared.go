package parser

import (
	"strings"

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
func (p *Parser) parsePrepareStmt() (ast.Node, error) {
	prepareTok := p.advance() // consume PREPARE
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// The inner statement is required and runs to the end of this segment (the
	// statement splitter has already isolated one statement per segment). Slice
	// the remaining source as raw text.
	if p.cur.Kind == tokEOF {
		// PREPARE x FROM <nothing> — Trino rejects this as a syntax error.
		return nil, p.syntaxErrorAtCur()
	}
	bodyStart := p.cur.Loc.Start
	bodyLoc := p.cur.Loc
	bodyEnd := p.cur.Loc.Start
	for p.cur.Kind != tokEOF {
		bodyEnd = p.cur.Loc.End
		p.advance()
	}
	body := strings.TrimSpace(p.sourceSlice(bodyStart, bodyEnd))
	if body == "" {
		return nil, &ParseError{Loc: ast.Loc{Start: bodyStart, End: bodyEnd}, Msg: "expected a statement after PREPARE ... FROM"}
	}

	// Validate the inner statement. Re-parsing the body through the full parser
	// reports a real error for non-statements (`PREPARE p FROM 1`,
	// `... FROM NOTASTATEMENT`). A "not yet supported" stub error means the
	// inner statement FORM is recognized (its concrete parser is a later DAG
	// node) — that counts as syntactically valid here, and the check tightens
	// automatically as those parsers land. See bodyHasHardError.
	if err := bodyHasHardError(body); err != nil {
		return nil, &ParseError{Loc: bodyLoc, Msg: "invalid statement after PREPARE ... FROM: " + err.Error()}
	}

	return &PrepareStmt{
		Name: name,
		Body: body,
		Loc:  ast.Loc{Start: prepareTok.Loc.Start, End: bodyEnd},
	}, nil
}

// bodyHasHardError re-parses a PREPARE inner-statement body and reports the
// first NON-STUB error, or nil if the body is a syntactically valid statement
// (possibly one whose concrete parser is not yet implemented). A "not yet
// supported" stub error is treated as success: it means the statement form was
// recognized by the dispatch but its body parser is a later DAG node. Any other
// error — an unknown leading keyword, a lex error, or (once the statement
// parsers land) a genuine inner syntax error — is a hard error.
func bodyHasHardError(body string) error {
	_, errs := Parse(body)
	for i := range errs {
		if !strings.Contains(errs[i].Msg, "not yet supported") {
			return &errs[i]
		}
	}
	return nil
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
