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
// to the CURRENT top of the precedence ladder. Each task that extends
// the ladder updates this single function; internal callers that need
// a sub-expression (inside array literals, tuple pairs, parentheses,
// bracket-index steps) always call parseExprTop and never need updating.
//
// This indirection exists to prevent a footgun: if every internal
// caller called parsePrimary directly, Tasks 6-9 would need to touch
// every call site as the ladder grows. parseExprTop keeps the fix in
// one place.
//
//	Task 2: parseExprTop → parseLiteral
//	Task 5: parseExprTop → parsePrimary (this task)
//	Task 6: parseExprTop → parseMathOp00
//	Task 7: parseExprTop → parsePredicate
//	Task 8: parseExprTop → parseOr
//	Task 9: parseExprTop → parseBagOp
func (p *Parser) parseExprTop() (ast.ExprNode, error) {
	return p.parsePredicate()
}

// parseVarRef handles optional @-prefix plus symbolPrimitive. Matches
// grammar rule varRefExpr (lines 635-636).
//
// NOTE: Task 10 upgrades this function to detect IDENT followed by
// PAREN_LEFT (function call form) and return a deferred-feature stub
// error. The Task 1 version handles only bare varref.
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
	return &ast.VarRef{
		Name:          name,
		AtPrefixed:    atPrefixed,
		CaseSensitive: caseSensitive,
		Loc:           ast.Loc{Start: start, End: nameLoc.End},
	}, nil
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
