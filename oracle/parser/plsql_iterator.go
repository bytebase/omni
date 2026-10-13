package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseIterator parses the iterator of a FOR loop or of an iterator choice,
// the current token after FOR, through its last iteration control.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/plsql-control-statements.html
//
//	iterand [, iterand ] IN qual_iteration_ctl [, qual_iteration_ctl ]...
//	iterand: name [ MUTABLE | IMMUTABLE ] [ datatype ]
//
// Oracle 23ai checks from the text that a PAIRS OF control has two
// iterands (PLS-00855) and that two iterands have PAIRS OF controls alone
// (PLS-00856).
func (p *Parser) parseIterator() (*nodes.PLSQLIterator, error) {
	it := &nodes.PLSQLIterator{
		Iterands: &nodes.List{},
		Controls: &nodes.List{},
		Loc:      nodes.Loc{Start: p.pos()},
	}
	for {
		iterand, err := p.parseIterand()
		if err != nil {
			return nil, err
		}
		it.Iterands.Items = append(it.Iterands.Items, iterand)
		if p.cur.Type != ',' || it.Iterands.Len() == 2 {
			break
		}
		p.advance() // consume ','
	}
	if p.cur.Type != kwIN {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume IN
	pairs := it.Iterands.Len() == 2
	for {
		ctl, err := p.parseIterationControl()
		if err != nil {
			return nil, err
		}
		if (ctl.Kind == "PAIRS OF") != pairs {
			return nil, p.syntaxErrorAtNode(ctl)
		}
		it.Controls.Items = append(it.Controls.Items, ctl)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	it.Loc.End = p.prev.End
	return it, nil
}

// parseIterand parses name [ MUTABLE | IMMUTABLE ] [ datatype ].
func (p *Parser) parseIterand() (*nodes.PLSQLIterand, error) {
	iterand := &nodes.PLSQLIterand{Loc: nodes.Loc{Start: p.pos()}}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, p.syntaxErrorAtCur()
	}
	iterand.Name = name
	if p.isKeywordStr("MUTABLE") || p.isKeywordStr("IMMUTABLE") {
		iterand.Mutability = p.cur.Str
		p.advance()
	}
	if p.cur.Type != kwIN && p.cur.Type != ',' {
		iterand.DataType, err = p.parsePLSQLTypeName()
		if err != nil {
			return nil, err
		}
		if iterand.DataType == nil || iterand.DataType.Names.Len() == 0 {
			return nil, p.syntaxErrorAtCur()
		}
	}
	iterand.Loc.End = p.prev.End
	return iterand, nil
}

// parseIterationControl parses one qualified iteration control. A single
// expression is also a cursor, a cursor call, or a parenthesized query.
// Oracle 23ai rejects REVERSE before a single expression, a cursor, or
// REPEAT (PLS-00857, PLS-00682), BY after anything but a range, and WHILE
// after WHEN (PLS-00103).
//
//	[ REVERSE ] { lower .. upper [ BY step ] | expr | REPEAT expr
//	            | { VALUES | INDICES | PAIRS } OF expr }
//	[ WHILE condition ] [ WHEN condition ]
func (p *Parser) parseIterationControl() (*nodes.PLSQLIterationControl, error) {
	ctl := &nodes.PLSQLIterationControl{Loc: nodes.Loc{Start: p.pos()}}
	if p.cur.Type == kwREVERSE {
		ctl.Reverse = true
		p.advance()
	}
	var err error
	switch {
	case p.isKeywordStr("REPEAT"):
		ctl.Kind = "REPEAT"
		p.advance()
		ctl.Expr, err = p.parseRequiredExpr()
	case (p.isKeywordStr("VALUES") || p.isKeywordStr("INDICES") || p.isKeywordStr("PAIRS")) && p.peekNext().Type == kwOF:
		ctl.Kind = p.cur.Str + " OF"
		p.advance() // consume VALUES, INDICES, or PAIRS
		p.advance() // consume OF
		ctl.Expr, err = p.parseRequiredExpr()
	default:
		var expr nodes.ExprNode
		expr, err = p.parseRequiredExpr()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != tokDOTDOT {
			ctl.Kind = "EXPR"
			ctl.Expr = expr
			break
		}
		ctl.Kind = "RANGE"
		ctl.Lower = expr
		p.advance() // consume ..
		if ctl.Upper, err = p.parseRequiredExpr(); err != nil {
			return nil, err
		}
		if p.cur.Type == kwBY {
			p.advance()
			ctl.Step, err = p.parseRequiredExpr()
		}
	}
	if err != nil {
		return nil, err
	}
	if ctl.Reverse && (ctl.Kind == "EXPR" || ctl.Kind == "REPEAT") {
		return nil, p.syntaxErrorAtTok(Token{Type: kwREVERSE, Loc: ctl.Loc.Start, End: ctl.Loc.Start + len("REVERSE")})
	}
	if p.isKeywordStr("WHILE") {
		p.advance()
		if ctl.While, err = p.parseRequiredExpr(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == kwWHEN {
		p.advance()
		if ctl.When, err = p.parseRequiredExpr(); err != nil {
			return nil, err
		}
	}
	ctl.Loc.End = p.prev.End
	return ctl, nil
}

// parseRequiredExpr parses an expression that must be there.
func (p *Parser) parseRequiredExpr() (nodes.ExprNode, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return expr, nil
}

// atQualifiedChoice reports whether a PL/SQL constructor argument opens
// with an iterator choice or the OTHERS choice.
func (p *Parser) atQualifiedChoice() bool {
	return p.cur.Type == kwFOR || p.isKeywordStr("OTHERS") && p.peekNext().Type == tokASSOC
}

// continuesQualifiedChoice reports whether the argument arg, just parsed,
// is the first choice of an explicit choice: a range or a choice list
// follows it, or => follows anything but one identifier, which names a
// parameter or a record field (NamedArgExpr).
func (p *Parser) continuesQualifiedChoice(arg nodes.ExprNode) bool {
	switch p.cur.Type {
	case '|', tokDOTDOT:
		return true
	case tokASSOC:
		ref, ok := arg.(*nodes.ColumnRef)
		return !ok || ref.Schema != "" || ref.Table != "" || ref.Column == ""
	}
	return false
}

// parseChoiceAssoc parses an explicit choice of a PL/SQL qualified
// expression (Oracle 18c, iterator choices 21c), first its first choice
// when the caller has parsed it. SQL takes none of them: Oracle 23ai
// rejects t(1 => 1) in a query with ORA-01748.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/qualified-expression.html
//
//	choice [ | choice ]... => expr
//	choice: expr | lower .. upper | OTHERS
//	FOR iterator [ INDEX expr | SEQUENCE ] => expr
func (p *Parser) parseChoiceAssoc(first nodes.ExprNode) (nodes.ExprNode, error) {
	assoc := &nodes.ChoiceAssocExpr{Loc: nodes.Loc{Start: p.pos()}}
	var err error
	switch {
	case first != nil:
		assoc.Loc.Start = nodes.NodeLoc(first).Start
		assoc.Choices = &nodes.List{}
		choice := first
		for {
			if p.cur.Type == tokDOTDOT {
				p.advance() // consume ..
				upper, err := p.parseRequiredExpr()
				if err != nil {
					return nil, err
				}
				choice = &nodes.ChoiceRange{
					Lower: choice,
					Upper: upper,
					Loc:   nodes.Loc{Start: nodes.NodeLoc(choice).Start, End: p.prev.End},
				}
			}
			assoc.Choices.Items = append(assoc.Choices.Items, choice)
			if p.cur.Type != '|' {
				break
			}
			p.advance() // consume |
			if choice, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
		}
	case p.cur.Type == kwFOR:
		p.advance() // consume FOR
		if assoc.Iterator, err = p.parseIterator(); err != nil {
			return nil, err
		}
		switch {
		case p.cur.Type == kwINDEX:
			p.advance()
			if assoc.Index, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
		case p.isKeywordStr("SEQUENCE"):
			p.advance()
			assoc.Sequence = true
		}
	default: // OTHERS
		p.advance()
		assoc.Others = true
	}
	if p.cur.Type != tokASSOC {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume =>
	if assoc.Value, err = p.parseRequiredExpr(); err != nil {
		return nil, err
	}
	assoc.Loc.End = p.prev.End
	return assoc, nil
}
