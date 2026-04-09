package parser

import (
	"github.com/bytebase/omni/partiql/ast"
)

// parseLiteral dispatches on the current token to produce one of the
// literal AST nodes. Handles the 6 real forms (NULL/MISSING/TRUE/
// FALSE/string/number/Ion) directly. DATE and TIME literal forms
// are stubbed with a "deferred to parser-datetime-literals (DAG
// node 18)" error.
//
// Grammar: literal (PartiQLParser.g4 lines 661-672).
func (p *Parser) parseLiteral() (ast.ExprNode, error) {
	switch p.cur.Type {
	case tokNULL:
		loc := p.cur.Loc
		p.advance()
		return &ast.NullLit{Loc: loc}, nil

	case tokMISSING:
		loc := p.cur.Loc
		p.advance()
		return &ast.MissingLit{Loc: loc}, nil

	case tokTRUE:
		loc := p.cur.Loc
		p.advance()
		return &ast.BoolLit{Val: true, Loc: loc}, nil

	case tokFALSE:
		loc := p.cur.Loc
		p.advance()
		return &ast.BoolLit{Val: false, Loc: loc}, nil

	case tokSCONST:
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.StringLit{Val: val, Loc: loc}, nil

	case tokICONST, tokFCONST:
		// NumberLit.Val preserves the raw source text so callers can
		// distinguish integer/decimal/scientific forms. Token.Str is
		// already the raw slice (scanNumber in the lexer preserves it).
		val := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.NumberLit{Val: val, Loc: loc}, nil

	case tokION_LITERAL:
		// Lexer's scanIonLiteral delivers the verbatim inner content
		// (between backticks) in Token.Str. No further decoding at
		// this layer — that is DAG node 17's job.
		text := p.cur.Str
		loc := p.cur.Loc
		p.advance()
		return &ast.IonLit{Text: text, Loc: loc}, nil

	case tokDATE:
		return nil, p.deferredFeature("DATE literal", "parser-datetime-literals (DAG node 18)")

	case tokTIME:
		return nil, p.deferredFeature("TIME literal", "parser-datetime-literals (DAG node 18)")
	}

	return nil, &ParseError{
		Message: "expected literal",
		Loc:     p.cur.Loc,
	}
}
