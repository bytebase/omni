// Package parser provides a hand-written recursive-descent parser for
// PartiQL. This file holds the Parser struct, token buffer helpers,
// error construction, and the shared utility parsers
// (parseSymbolPrimitive, parseVarRef, parseType).
//
// Expression parsing is split across expr.go (precedence ladder) and
// exprprimary.go (primary-expression dispatch). Literal parsing lives
// in literals.go; path-step chains in path.go.
//
// Upstream: this parser consumes tokens produced by Lexer (lexer.go,
// token.go, keywords.go). The lexer uses first-error-and-stop via
// Lexer.Err; the parser surfaces lexer errors as ParseError.
package parser

import (
	"fmt"
	"strconv"

	"github.com/bytebase/omni/partiql/ast"
)

// ParseError represents a syntax error during parsing.
//
// Loc is the full byte range of the current (offending) token when the
// error was raised. Using the full token span rather than a zero-width
// point lets future error-rendering code highlight the problematic
// region rather than a bare position marker. The parser uses fail-fast
// semantics — the first ParseError aborts the parse and is returned to
// the caller.
type ParseError struct {
	Message string
	Loc     ast.Loc
}

// Error renders a human-readable message including the byte position.
func (e *ParseError) Error() string {
	return fmt.Sprintf("syntax error at position %d: %s", e.Loc.Start, e.Message)
}

// Parser is the recursive-descent parser for PartiQL.
//
// Fields:
//   - lexer:   the token source
//   - cur:     current lookahead token (the "peek" slot)
//   - prev:    most recently consumed token, used to compute end Locs
//   - nextBuf: one-token lookahead buffer for peekNext()
//   - hasNext: true when nextBuf holds a valid token
type Parser struct {
	lexer   *Lexer
	cur     Token
	prev    Token
	nextBuf Token
	hasNext bool
}

// NewParser creates a Parser that reads from the given input string.
// The first token is primed into cur before the constructor returns.
func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.advance() // prime cur with the first token
	return p
}

// advance moves cur to prev and reads the next token from the lexer.
// When the lexer's first-error-and-stop contract fires, subsequent
// calls produce tokEOF with a nil Str — callers must invoke
// checkLexerErr() at strategic points (function entry, after expect)
// to surface the lexer error. This function does NOT embed the error
// message in cur.Str.
func (p *Parser) advance() {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.Next()
	}
}

// peek returns the current token without consuming it.
func (p *Parser) peek() Token {
	return p.cur
}

// parserState is a snapshot of the full mutable parse position: the
// parser's token-buffer slots plus the lexer's scan cursor. It is the
// minimal state needed to rewind the parser to an earlier point so a
// production can be attempted speculatively and abandoned on failure.
//
// The lexer reads forward-only from an immutable input string, so its
// entire mutable state is (pos, start, Err); capturing those alongside
// the parser's (cur, prev, nextBuf, hasNext) makes save/restore exact.
type parserState struct {
	cur      Token
	prev     Token
	nextBuf  Token
	hasNext  bool
	lexPos   int
	lexStart int
	lexErr   error
}

// save captures the current parse position for later restore. Use it to
// bound a speculative parse: snapshot, attempt the optional production,
// and restore on failure to retry an alternative. The returned value is
// a by-value copy — no aliasing with the live parser/lexer.
func (p *Parser) save() parserState {
	return parserState{
		cur:      p.cur,
		prev:     p.prev,
		nextBuf:  p.nextBuf,
		hasNext:  p.hasNext,
		lexPos:   p.lexer.pos,
		lexStart: p.lexer.start,
		lexErr:   p.lexer.Err,
	}
}

// restore rewinds the parser and lexer to a position captured by save.
// After restore, the next token the parser sees is exactly the one it
// saw when save was called.
func (p *Parser) restore(s parserState) {
	p.cur = s.cur
	p.prev = s.prev
	p.nextBuf = s.nextBuf
	p.hasNext = s.hasNext
	p.lexer.pos = s.lexPos
	p.lexer.start = s.lexStart
	p.lexer.Err = s.lexErr
}

// peekNext returns the token AFTER cur without consuming either cur or
// the returned token. Uses the nextBuf slot; subsequent calls are
// idempotent until advance() is called.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.Next()
		p.hasNext = true
	}
	return p.nextBuf
}

// match consumes cur if its Type matches any of the given types.
// Returns true on match, false otherwise. Idiomatic for optional
// keywords like DISTINCT/ALL or NOT.
func (p *Parser) match(types ...int) bool {
	for _, t := range types {
		if p.cur.Type == t {
			p.advance()
			return true
		}
	}
	return false
}

// expect consumes cur if its Type matches tokenType, otherwise returns
// a *ParseError. Used for required tokens (PAREN_RIGHT, COMMA, etc.).
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		tok := p.cur
		p.advance()
		return tok, nil
	}
	return Token{}, &ParseError{
		Message: fmt.Sprintf("expected %s, got %q", tokenName(tokenType), p.cur.Str),
		Loc:     p.cur.Loc,
	}
}

// checkLexerErr returns a *ParseError wrapping the lexer's error if
// any. Parser methods call this at strategic points (function entry,
// after expect) to propagate lexer errors into the parse result.
func (p *Parser) checkLexerErr() error {
	if p.lexer.Err != nil {
		return &ParseError{
			Message: p.lexer.Err.Error(),
			Loc:     p.cur.Loc,
		}
	}
	return nil
}

// deferredFeature constructs a stub error pointing at the current
// token. Every non-foundation case in parsePrimaryBase (and a few
// other dispatch points) calls this helper so the error format is
// uniform across the whole parser.
//
// The error message format is the grep contract for future DAG node
// implementers: running
//
//	grep -rn "deferred to parser-builtins" partiql/parser/
//
// returns the full list of stub call sites node 15 needs to replace.
func (p *Parser) deferredFeature(feature, ownerNode string) error {
	return &ParseError{
		Message: fmt.Sprintf("%s is deferred to %s", feature, ownerNode),
		Loc:     p.cur.Loc,
	}
}

// parseSymbolPrimitive consumes an unquoted or double-quoted identifier
// and returns the name, case-sensitivity flag, and location. Rejects
// bare keywords. Matches grammar rule symbolPrimitive (line 742).
func (p *Parser) parseSymbolPrimitive() (name string, caseSensitive bool, loc ast.Loc, err error) {
	switch p.cur.Type {
	case tokIDENT:
		name = p.cur.Str
		caseSensitive = false
		loc = p.cur.Loc
		p.advance()
		return
	case tokIDENT_QUOTED:
		name = p.cur.Str
		caseSensitive = true
		loc = p.cur.Loc
		p.advance()
		return
	default:
		err = &ParseError{
			Message: fmt.Sprintf("expected identifier, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
		return
	}
}

// ParseExpr parses a single expression from the parser's input and
// asserts that the entire input was consumed. Trailing tokens after
// the expression produce an "unexpected token" error — callers cannot
// silently drop input.
//
// This is the foundation-level public entry point. Nodes 5-8 will
// add ParseStatement and ParseScript (SelectStmt-producing forms).
// Node 8 will add the top-level Parse(sql) function.
//
// ParseExpr delegates to parseExprTop for the actual parse (wrapping
// with a lexer-error check and an EOF check). Internal callers that
// parse an expression embedded in a larger construct (array items,
// tuple keys/values, parenthesized expressions, bracket-index steps)
// call parseExprTop directly because they must not assert EOF —
// their closing delimiter is the next token after the expression.
func (p *Parser) ParseExpr() (ast.ExprNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	expr, err := p.parseExprTop()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != tokEOF {
		return nil, &ParseError{
			Message: fmt.Sprintf("unexpected token %q after expression", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
	return expr, nil
}

// parseExprTop is the internal expression entry point. It dispatches
// to the top of the precedence ladder (parseBagOp) without the EOF
// check that ParseExpr adds.
//
// Internal callers that parse a sub-expression embedded inside a
// larger construct (array items, tuple keys/values, parenthesized
// expressions, bracket-index steps) call parseExprTop rather than
// ParseExpr because they must not assert EOF — their closing
// delimiter is the next token after the sub-expression.
//
// Future DAG nodes that extend the precedence ladder (e.g. node 5
// replacing parseBagOp/parseSelectExpr with real SELECT/UNION
// handling) only need to touch the bodies of those functions;
// parseExprTop's single-line dispatch stays the same.
func (p *Parser) parseExprTop() (ast.ExprNode, error) {
	return p.parseBagOp()
}

// parseVarRef handles optional @-prefix plus symbolPrimitive, and
// detects the function-call form `IDENT PAREN_LEFT ...`. Matches
// grammar rules varRefExpr (635-636) and functionCall#FunctionCallIdent
// (615).
//
// As of node 15a (parser-builtins-generic-call), generic IDENT(args)
// calls produce an *ast.FuncCall here. Typed keyword builtins (CAST,
// CASE, COALESCE, SUBSTRING, ...) keep their token-dispatched stubs in
// exprprimary.go until 15b lands.
func (p *Parser) parseVarRef() (ast.ExprNode, error) {
	start := p.cur.Loc.Start
	atPrefixed := false
	if p.cur.Type == tokAT_SIGN {
		atPrefixed = true
		p.advance()
	}
	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	// Function call: name=symbolPrimitive PAREN_LEFT (expr (COMMA expr)*)? PAREN_RIGHT
	// (PartiQLParser.g4:615 — FunctionCallIdent). Implements DAG node 15a.
	// @-prefix is NOT permitted before a function call: ANTLR's functionCall
	// rule names symbolPrimitive directly, only varRefExpr admits @.
	if p.cur.Type == tokPAREN_LEFT {
		if atPrefixed {
			return nil, &ParseError{
				Message: "@-prefix is not allowed before a function call",
				Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
			}
		}
		args, endOff, err := p.parseFuncCallArgs()
		if err != nil {
			return nil, err
		}
		return &ast.FuncCall{
			Name: name,
			Args: args,
			Loc:  ast.Loc{Start: start, End: endOff},
		}, nil
	}
	return &ast.VarRef{
		Name:          name,
		AtPrefixed:    atPrefixed,
		CaseSensitive: caseSensitive,
		Loc:           ast.Loc{Start: start, End: nameLoc.End},
	}, nil
}

// parseFuncCallArgs consumes the parenthesized argument list of a function
// call: PAREN_LEFT (expr (COMMA expr)*)? PAREN_RIGHT. The opening paren
// must be the current token. Returns the parsed args (nil if empty), the
// 0-based byte End offset of PAREN_RIGHT, and any parse error.
//
// Used by DAG node 15a (parser-builtins-generic-call) and will be reused
// by 15b for the typed builtins whose argument list is also comma-separated
// (DATE_ADD, DATE_DIFF, COALESCE, NULLIF, SIZE, EXISTS, etc.).
func (p *Parser) parseFuncCallArgs() ([]ast.ExprNode, int, error) {
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, 0, err
	}
	var args []ast.ExprNode
	if p.cur.Type != tokPAREN_RIGHT {
		for {
			arg, err := p.parseExprTop()
			if err != nil {
				return nil, 0, err
			}
			args = append(args, arg)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume COMMA
		}
	}
	endOff := p.cur.Loc.End
	if _, err := p.expect(tokPAREN_RIGHT); err != nil {
		return nil, 0, err
	}
	return args, endOff, nil
}

// parseType consumes one of the PartiQL type forms and returns a
// *ast.TypeRef. Handles:
//
//   - Atomic types: NULL, BOOL, BOOLEAN, SMALLINT, INT/INT2/INT4/INT8,
//     INTEGER/INTEGER2/INTEGER4/INTEGER8, BIGINT, REAL, TIMESTAMP,
//     CHAR, CHARACTER, MISSING, STRING, SYMBOL, BLOB, CLOB, DATE,
//     STRUCT, TUPLE, LIST, SEXP, BAG, ANY
//   - DOUBLE PRECISION (two-token form)
//   - Parameterized single-arg: CHAR(n), CHARACTER(n), FLOAT(p), VARCHAR(n)
//   - CHARACTER VARYING [(n)]
//   - Parameterized two-arg: DECIMAL(p,s), DEC(p,s), NUMERIC(p,s)
//   - TIME [(p)] [WITH TIME ZONE]
//   - Custom: any symbolPrimitive identifier (fallback)
//
// Grammar: type (PartiQLParser.g4 lines 674-686).
//
// Foundation ships parseType even though CAST is stubbed so that
// parser-ddl (DAG node 7) and parser-builtins (DAG node 15) can each
// consume it without coupling.
func (p *Parser) parseType() (*ast.TypeRef, error) {
	start := p.cur.Loc.Start

	// DOUBLE PRECISION is a two-token form.
	if p.cur.Type == tokDOUBLE {
		p.advance()
		if p.cur.Type != tokPRECISION {
			return nil, &ParseError{
				Message: fmt.Sprintf("expected PRECISION after DOUBLE, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		end := p.cur.Loc.End
		p.advance()
		return &ast.TypeRef{
			Name: "DOUBLE PRECISION",
			Loc:  ast.Loc{Start: start, End: end},
		}, nil
	}

	// CHARACTER VARYING is another two-token form that also takes an
	// optional (n) argument.
	if p.cur.Type == tokCHARACTER && p.peekNext().Type == tokVARYING {
		p.advance() // CHARACTER
		p.advance() // VARYING
		end := p.prev.Loc.End
		args, argsEnd, err := p.parseOptionalTypeArgs(1)
		if err != nil {
			return nil, err
		}
		if argsEnd > 0 {
			end = argsEnd
		}
		return &ast.TypeRef{
			Name: "CHARACTER VARYING",
			Args: args,
			Loc:  ast.Loc{Start: start, End: end},
		}, nil
	}

	// TIME [(p)] [WITH TIME ZONE] — TIME requires a tailored parse path
	// because it supports both a precision arg and the WITH TIME ZONE
	// trailing keywords.
	if p.cur.Type == tokTIME {
		p.advance()
		end := p.prev.Loc.End
		args, argsEnd, err := p.parseOptionalTypeArgs(1)
		if err != nil {
			return nil, err
		}
		if argsEnd > 0 {
			end = argsEnd
		}
		withTZ := false
		if p.cur.Type == tokWITH && p.peekNext().Type == tokTIME {
			// WITH TIME ZONE
			p.advance() // WITH
			p.advance() // TIME
			if p.cur.Type != tokZONE {
				return nil, &ParseError{
					Message: fmt.Sprintf("expected ZONE after WITH TIME, got %q", p.cur.Str),
					Loc:     p.cur.Loc,
				}
			}
			withTZ = true
			end = p.cur.Loc.End
			p.advance() // ZONE
		}
		return &ast.TypeRef{
			Name:         "TIME",
			Args:         args,
			WithTimeZone: withTZ,
			Loc:          ast.Loc{Start: start, End: end},
		}, nil
	}

	// Atomic types (no args). Map the keyword token to its canonical
	// uppercase name. This switch handles the "bare keyword" cases from
	// the grammar's TypeAtomic alternative (lines 675-680).
	atomicName := ""
	switch p.cur.Type {
	case tokNULL:
		atomicName = "NULL"
	case tokBOOL:
		atomicName = "BOOL"
	case tokBOOLEAN:
		atomicName = "BOOLEAN"
	case tokSMALLINT:
		atomicName = "SMALLINT"
	case tokINT2:
		atomicName = "INT2"
	case tokINTEGER2:
		atomicName = "INTEGER2"
	case tokINT4:
		atomicName = "INT4"
	case tokINTEGER4:
		atomicName = "INTEGER4"
	case tokINT8:
		atomicName = "INT8"
	case tokINTEGER8:
		atomicName = "INTEGER8"
	case tokINT:
		atomicName = "INT"
	case tokINTEGER:
		atomicName = "INTEGER"
	case tokBIGINT:
		atomicName = "BIGINT"
	case tokREAL:
		atomicName = "REAL"
	case tokTIMESTAMP:
		atomicName = "TIMESTAMP"
	case tokMISSING:
		atomicName = "MISSING"
	case tokSTRING:
		atomicName = "STRING"
	case tokSYMBOL:
		atomicName = "SYMBOL"
	case tokBLOB:
		atomicName = "BLOB"
	case tokCLOB:
		atomicName = "CLOB"
	case tokDATE:
		atomicName = "DATE"
	case tokSTRUCT:
		atomicName = "STRUCT"
	case tokTUPLE:
		atomicName = "TUPLE"
	case tokLIST:
		atomicName = "LIST"
	case tokSEXP:
		atomicName = "SEXP"
	case tokBAG:
		atomicName = "BAG"
	case tokANY:
		atomicName = "ANY"
	}
	if atomicName != "" {
		end := p.cur.Loc.End
		p.advance()
		return &ast.TypeRef{
			Name: atomicName,
			Loc:  ast.Loc{Start: start, End: end},
		}, nil
	}

	// Parameterized single-arg: CHAR(n), CHARACTER(n), FLOAT(p), VARCHAR(n).
	switch p.cur.Type {
	case tokCHAR:
		return p.parseTypeWithArgs(start, "CHAR", 1, 1)
	case tokCHARACTER:
		// Note: CHARACTER VARYING already handled above via peekNext.
		return p.parseTypeWithArgs(start, "CHARACTER", 1, 1)
	case tokFLOAT:
		return p.parseTypeWithArgs(start, "FLOAT", 1, 1)
	case tokVARCHAR:
		return p.parseTypeWithArgs(start, "VARCHAR", 1, 1)
	}

	// Parameterized two-arg: DECIMAL(p,s), DEC(p,s), NUMERIC(p,s).
	switch p.cur.Type {
	case tokDECIMAL:
		return p.parseTypeWithArgs(start, "DECIMAL", 1, 2)
	case tokDEC:
		return p.parseTypeWithArgs(start, "DEC", 1, 2)
	case tokNUMERIC:
		return p.parseTypeWithArgs(start, "NUMERIC", 1, 2)
	}

	// Custom type: any symbolPrimitive. The grammar calls this
	// TypeCustom (line 685).
	if p.cur.Type == tokIDENT || p.cur.Type == tokIDENT_QUOTED {
		name, _, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.TypeRef{
			Name: name,
			Loc:  ast.Loc{Start: start, End: nameLoc.End},
		}, nil
	}

	return nil, &ParseError{
		Message: fmt.Sprintf("expected type, got %q", p.cur.Str),
		Loc:     p.cur.Loc,
	}
}

// parseTypeWithArgs consumes a type keyword followed by an optional
// parenthesized argument list of integers. Used by CHAR(n), FLOAT(p),
// VARCHAR(n), DECIMAL(p,s), and related forms.
//
//   - name:    canonical type name (e.g. "DECIMAL")
//   - minArgs: minimum number of integer args when the parenthesized
//     list is present. Parens themselves are always optional (every
//     caller allows the bare form, e.g. CHAR or DECIMAL), so the
//     min-args enforcement is nested inside the paren-present branch.
//     If the grammar ever grows a type where parens are REQUIRED,
//     callers will need a separate enforcement path.
//   - maxArgs: maximum number of integer args when parens are present.
func (p *Parser) parseTypeWithArgs(start int, name string, minArgs, maxArgs int) (*ast.TypeRef, error) {
	// Consume the type keyword.
	end := p.cur.Loc.End
	p.advance()
	args, argsEnd, err := p.parseOptionalTypeArgs(maxArgs)
	if err != nil {
		return nil, err
	}
	if argsEnd > 0 {
		end = argsEnd
		if len(args) < minArgs {
			return nil, &ParseError{
				Message: fmt.Sprintf("%s: expected at least %d argument(s), got %d", name, minArgs, len(args)),
				Loc:     ast.Loc{Start: start, End: end},
			}
		}
	}
	return &ast.TypeRef{
		Name: name,
		Args: args,
		Loc:  ast.Loc{Start: start, End: end},
	}, nil
}

// parseOptionalTypeArgs consumes an optional (n[, n]*) argument list
// bounded by parentheses. If cur is not PAREN_LEFT, returns
// (nil, 0, nil) without consuming anything. Otherwise consumes the
// parenthesized comma-separated integer list and returns the parsed
// args plus the End position of the closing paren.
//
// maxArgs is the hard limit on the number of integers; exceeding it
// yields a ParseError.
func (p *Parser) parseOptionalTypeArgs(maxArgs int) (args []int, end int, err error) {
	if p.cur.Type != tokPAREN_LEFT {
		return nil, 0, nil
	}
	p.advance() // consume (
	for {
		// Targeted error for empty arg list (e.g. `DECIMAL()`) and
		// trailing comma (e.g. `DECIMAL(10,)`). Without these guards,
		// the integer-check below would fire "expected integer argument,
		// got ')'" which misattributes the error to the closing paren.
		if p.cur.Type == tokPAREN_RIGHT {
			if len(args) == 0 {
				return nil, 0, &ParseError{
					Message: "empty type argument list",
					Loc:     p.cur.Loc,
				}
			}
			return nil, 0, &ParseError{
				Message: "unexpected trailing comma in type arguments",
				Loc:     p.cur.Loc,
			}
		}
		if p.cur.Type != tokICONST {
			return nil, 0, &ParseError{
				Message: fmt.Sprintf("expected integer argument, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		n, perr := parseIntLiteral(p.cur.Str)
		if perr != nil {
			return nil, 0, &ParseError{
				Message: fmt.Sprintf("invalid integer argument %q: %v", p.cur.Str, perr),
				Loc:     p.cur.Loc,
			}
		}
		args = append(args, n)
		p.advance()
		if len(args) > maxArgs {
			return nil, 0, &ParseError{
				Message: fmt.Sprintf("too many type arguments (max %d)", maxArgs),
				Loc:     p.cur.Loc,
			}
		}
		if p.cur.Type == tokCOMMA {
			p.advance()
			continue
		}
		break
	}
	rp, perr := p.expect(tokPAREN_RIGHT)
	if perr != nil {
		return nil, 0, perr
	}
	return args, rp.Loc.End, nil
}

// parseIntLiteral converts a token's Str (raw source text for an
// integer literal) into an int. We use strconv.Atoi directly because
// PartiQL integer literals are plain decimal digits per the grammar
// LITERAL_INTEGER rule (no hex, no underscores, no sign — signs come
// from the unary-minus operator).
func parseIntLiteral(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseStatement parses a single top-level statement and asserts that
// the entire input was consumed. Supports DDL (CREATE/DROP TABLE/INDEX)
// and DML (INSERT/UPDATE/DELETE/REPLACE/UPSERT/REMOVE) directly; defers
// EXEC to parse-entry (DAG node 8). DQL (SELECT and set-ops) falls
// through to parseExprTop; if the result is already a StmtNode it is
// returned, otherwise the bare-expression form is deferred to
// parse-entry (DAG node 8).
func (p *Parser) ParseStatement() (ast.StmtNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	var stmt ast.StmtNode
	var err error
	switch p.cur.Type {
	case tokCREATE:
		stmt, err = p.parseCreateCommand()
	case tokDROP:
		stmt, err = p.parseDropCommand()
	case tokINSERT:
		stmt, err = p.parseInsertStmt()
	case tokUPDATE:
		stmt, err = p.parseUpdateStmt()
	case tokDELETE:
		stmt, err = p.parseDeleteStmt()
	case tokREPLACE:
		stmt, err = p.parseReplaceStmt()
	case tokUPSERT:
		stmt, err = p.parseUpsertStmt()
	case tokREMOVE:
		stmt, err = p.parseRemoveStmt()
	case tokFROM:
		// FROM-led DML: `FROM fromClause whereClause? dmlBaseCommand+
		// returningClause?` (dml#DmlBaseWrapper 2nd alt). A leading FROM at
		// statement position is always DML — SELECT begins with SELECT, and
		// a bare expression cannot start with FROM.
		stmt, err = p.parseFromLedDml()
	case tokEXEC, tokEXECUTE:
		p.advance() // consume EXEC/EXECUTE
		stmt, err = p.parseExecCommand()
	default:
		// DQL fallback: parse as expression. If the result implements
		// StmtNode (e.g. SelectStmt after node 5), return it. Otherwise
		// defer bare-expression-as-statement to parse-entry (node 8).
		var expr ast.ExprNode
		expr, err = p.parseExprTop()
		if err != nil {
			return nil, err
		}
		if sub, ok := expr.(*ast.SubLink); ok {
			stmt = sub.Stmt
		} else if sn, ok := expr.(ast.StmtNode); ok {
			stmt = sn
		} else {
			return nil, p.deferredFeature("bare expression as statement", "parse-entry (DAG node 8)")
		}
	}
	if err != nil {
		return nil, err
	}
	if p.cur.Type != tokEOF {
		return nil, &ParseError{
			Message: fmt.Sprintf("unexpected token %q after statement", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
	return stmt, nil
}
