// Package parser - service_broker.go implements T-SQL Service Broker statement parsing.
//
// Every statement here follows ScriptDOM's TSql160.g (mssql/testdata/TSql160.g4
// is the conversion; each rule cites its TSql160.g line), checked against the
// pinned SQL Server 2022, whose verdict wins where the two differ.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// validationValues is the set of valid values for MESSAGE TYPE VALIDATION.
// BNF: VALIDATION = { NONE | EMPTY | WELL_FORMED_XML | VALID_XML WITH SCHEMA COLLECTION ... }
var validationValues = newOptionSet(kwNONE).withIdents("EMPTY", "WELL_FORMED_XML", "VALID_XML")

// contractSentByValues matches SqlScriptDOM's MessageBodyDirection enum
// used in CREATE CONTRACT (message SENT BY {INITIATOR | TARGET | ANY}).
var contractSentByValues = newOptionSet(kwANY).withIdents("INITIATOR", "TARGET", "ANY")

// brokerIsWord reports whether the current token is the word `word` (upper
// case): its keyword token when the lexer has one, else an undelimited
// identifier with that text.
func (p *Parser) brokerIsWord(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// brokerExpectWord consumes the word `word` or fails at the current token.
func (p *Parser) brokerExpectWord(word string) error {
	if !p.brokerIsWord(word) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// brokerExpectIdent consumes an identifier (ScriptDOM's `identifier`).
func (p *Parser) brokerExpectIdent() (string, error) {
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// brokerIsString reports whether the current token is a string literal.
func (p *Parser) brokerIsString() bool {
	return p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST
}

// brokerExpectString consumes a string literal (ScriptDOM's `stringLiteral`).
func (p *Parser) brokerExpectString() (string, error) {
	if !p.brokerIsString() {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// brokerExpectOnOff consumes `= { ON | OFF }` (ScriptDOM's EqualsSign optionOnOff).
func (p *Parser) brokerExpectOnOff() (string, error) {
	if _, err := p.expect('='); err != nil {
		return "", err
	}
	if p.cur.Type != kwON && p.cur.Type != kwOFF {
		return "", p.syntaxErrorAtCur()
	}
	return strings.ToUpper(p.advance().Str), nil
}

// brokerParseAuthorization parses ScriptDOM's authorizationOpt,
// `( Authorization identifier )?`, returning the owner or "".
func (p *Parser) brokerParseAuthorization() (string, error) {
	if p.cur.Type != kwAUTHORIZATION {
		return "", nil
	}
	p.advance()
	return p.brokerExpectIdent()
}

// brokerParseObjectName parses a dotted object name of at most maxParts
// parts. SQL Server refuses a longer one with Msg 117 ("contains more than the
// maximum number of prefixes").
func (p *Parser) brokerParseObjectName(maxParts int) (*nodes.TableRef, error) {
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, p.syntaxErrorAtCur()
	}
	parts := 1
	switch {
	case ref.Server != "":
		parts = 4
	case ref.Database != "":
		parts = 3
	case ref.Schema != "":
		parts = 2
	}
	if parts > maxParts {
		return nil, &ParseError{Message: "object name contains more than the maximum number of prefixes", Position: ref.Loc.Start, Source: p.source}
	}
	return ref, nil
}

// brokerQualifiedName joins the parts a TableRef holds.
func brokerQualifiedName(ref *nodes.TableRef) string {
	name := ref.Object
	if ref.Schema != "" || ref.Database != "" {
		name = ref.Schema + "." + name
	}
	if ref.Database != "" {
		name = ref.Database + "." + name
	}
	return name
}

// brokerParseNameOrVariable parses ScriptDOM's identifierOrVariable, which
// SQL Server also takes as a dotted name (BEGIN DIALOG ... FROM SERVICE a.b,
// ON CONTRACT a.b, SEND ... MESSAGE TYPE a.b all parse on the pinned server).
func (p *Parser) brokerParseNameOrVariable() (string, error) {
	if p.cur.Type == tokVARIABLE {
		return p.advance().Str, nil
	}
	ref, err := p.parseTableRef()
	if err != nil {
		return "", err
	}
	if ref == nil {
		return "", p.syntaxErrorAtCur()
	}
	if ref.Server != "" {
		return "", p.syntaxErrorAtCur()
	}
	return brokerQualifiedName(ref), nil
}

// brokerParseStringOrVariable parses ScriptDOM's stringOrVariable.
func (p *Parser) brokerParseStringOrVariable() (string, error) {
	if p.brokerIsString() || p.cur.Type == tokVARIABLE {
		return p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// brokerParseExprText parses an expression and returns its source text.
func (p *Parser) brokerParseExprText() (nodes.ExprNode, string, error) {
	start := p.pos()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, "", err
	}
	if expr == nil {
		return nil, "", p.syntaxErrorAtCur()
	}
	return expr, p.source[start:p.prevEnd()], nil
}

// parseCreateMessageTypeStmt parses CREATE MESSAGE TYPE.
//
// TSql160.g:9511 createMessageTypeStatement:
//
//	MESSAGE TYPE identifier authorizationOpt ( messageTypeValidation )?
func (p *Parser) parseCreateMessageTypeStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// TYPE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "MESSAGE TYPE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.brokerParseAuthorization(); err != nil {
		return nil, err
	}

	if p.brokerIsWord("VALIDATION") {
		opt, err := p.brokerParseMessageTypeValidation()
		if err != nil {
			return nil, err
		}
		stmt.Options = &nodes.List{Items: []nodes.Node{opt}}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// brokerParseMessageTypeValidation parses messageTypeValidation.
//
// TSql160.g:9539 and 9548: VALIDATION = { NONE | EMPTY | WELL_FORMED_XML |
// VALID_XML WITH SCHEMA COLLECTION name }; WITH follows VALID_XML only, and
// VALID_XML requires it. The collection name has at most two parts.
func (p *Parser) brokerParseMessageTypeValidation() (*nodes.ServiceBrokerOption, error) {
	valLoc := p.pos()
	if err := p.brokerExpectWord("VALIDATION"); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	valType, err := p.expectOption(validationValues)
	if err != nil {
		return nil, err
	}
	if valType != "VALID_XML" {
		if p.cur.Type == kwWITH {
			return nil, p.syntaxErrorAtCur()
		}
		return &nodes.ServiceBrokerOption{Name: "VALIDATION", Value: valType, Loc: nodes.Loc{Start: valLoc, End: p.prevEnd()}}, nil
	}
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSCHEMA); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCOLLECTION); err != nil {
		return nil, err
	}
	ref, err := p.brokerParseObjectName(2)
	if err != nil {
		return nil, err
	}
	schemaName := ref.Object
	if ref.Schema != "" {
		schemaName = ref.Schema + "." + ref.Object
	}
	return &nodes.ServiceBrokerOption{Name: "VALIDATION", Value: "VALID_XML WITH SCHEMA COLLECTION " + schemaName, Loc: nodes.Loc{Start: valLoc, End: p.prevEnd()}}, nil
}

// parseCreateContractStmt parses CREATE CONTRACT.
//
// TSql160.g:1559 createContractStatement and TSql160.g:1586 contractMessage:
//
//	CONTRACT identifier authorizationOpt
//	    ( identifier SENT BY { ANY | INITIATOR | TARGET } [ ,...n ] )
//
// [DEFAULT] is a delimited identifier; the bare keyword DEFAULT is refused.
func (p *Parser) parseCreateContractStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// CONTRACT keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "CONTRACT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.brokerParseAuthorization(); err != nil {
		return nil, err
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	for {
		entryLoc := p.pos()
		msgTypeName, err := p.brokerExpectIdent()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwSENT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		sentBy, err := p.expectOption(contractSentByValues)
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: msgTypeName, Value: "SENT BY " + sentBy, Loc: nodes.Loc{Start: entryLoc, End: p.prevEnd()}})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: opts}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateQueueStmt parses CREATE QUEUE.
//
// TSql160.g:10189 createQueueStatement:
//
//	QUEUE schemaObjectThreePartName ( queueOptionList )? ( On stringOrIdentifier )?
func (p *Parser) parseCreateQueueStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// QUEUE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "QUEUE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	ref, err := p.brokerParseObjectName(3)
	if err != nil {
		return nil, err
	}
	stmt.Name = ref.Object
	if ref.Schema != "" {
		stmt.Name = ref.Schema + "." + ref.Object
	}

	var opts []nodes.Node
	if p.cur.Type == kwWITH {
		opts, err = p.brokerParseQueueOptions(opts, true)
		if err != nil {
			return nil, err
		}
	}

	// ON { filegroup | 'filegroup' | [DEFAULT] }: stringOrIdentifier, so the
	// bare keyword DEFAULT is refused.
	if p.cur.Type == kwON {
		onLoc := p.pos()
		p.advance()
		var val string
		switch {
		case p.brokerIsString():
			val = p.advance().Str
		case p.isIdentLike():
			val = strings.ToUpper(p.advance().Str)
		default:
			return nil, p.syntaxErrorAtCur()
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: "ON", Value: val, Loc: nodes.Loc{Start: onLoc, End: p.prevEnd()}})
	}

	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// brokerParseQueueOptions parses queueOptionList.
//
// TSql160.g:12324 queueOptionList and TSql160.g:12332 queueOption:
//
//	WITH queueOption [ ,...n ]
//	queueOption ::= { STATUS | RETENTION } = { ON | OFF }
//	              | ACTIVATION ( ... )
//	              | POISON_MESSAGE_HANDLING ( STATUS = { ON | OFF } )
//
// SQL Server refuses an option given twice.
func (p *Parser) brokerParseQueueOptions(opts []nodes.Node, create bool) ([]nodes.Node, error) {
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for {
		optLoc := p.pos()
		var name string
		for _, w := range []string{"STATUS", "RETENTION", "ACTIVATION", "POISON_MESSAGE_HANDLING"} {
			if p.brokerIsWord(w) {
				name = w
				break
			}
		}
		if name == "" || seen[name] {
			return nil, p.syntaxErrorAtCur()
		}
		seen[name] = true
		p.advance()
		switch name {
		case "STATUS", "RETENTION":
			val, err := p.brokerExpectOnOff()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		case "ACTIVATION":
			var err error
			opts, err = p.brokerParseQueueActivation(opts, create)
			if err != nil {
				return nil, err
			}
		default: // POISON_MESSAGE_HANDLING, TSql160.g:12488
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			if err := p.brokerExpectWord("STATUS"); err != nil {
				return nil, err
			}
			val, err := p.brokerExpectOnOff()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "POISON_MESSAGE_HANDLING", Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return opts, nil
}

// brokerParseQueueActivation parses the parenthesized ACTIVATION arguments.
//
// TSql160.g:12369 queueActivationOption and TSql160.g:12401
// activationQueueOptionArgument:
//
//	ACTIVATION ( { STATUS = { ON | OFF }
//	             | PROCEDURE_NAME = schemaObjectThreePartName
//	             | MAX_QUEUE_READERS = integer
//	             | EXECUTE AS { SELF | OWNER | 'user_name' }
//	             | DROP } [ ,...n ] )
//
// DROP stands alone and only in ALTER QUEUE (TSql160.g:12469). CREATE QUEUE
// needs PROCEDURE_NAME, MAX_QUEUE_READERS and EXECUTE AS together, and SQL
// Server refuses an argument given twice and EXECUTE AS CALLER.
func (p *Parser) brokerParseQueueActivation(opts []nodes.Node, create bool) ([]nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	drop := false
	for {
		optLoc := p.pos()
		if drop {
			return nil, p.syntaxErrorAtCur()
		}
		switch {
		case p.cur.Type == kwDROP:
			if create || len(seen) > 0 {
				return nil, p.syntaxErrorAtCur()
			}
			drop = true
			p.advance()
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:DROP", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		case p.cur.Type == kwEXEC || p.cur.Type == kwEXECUTE:
			if seen["EXECUTE AS"] {
				return nil, p.syntaxErrorAtCur()
			}
			seen["EXECUTE AS"] = true
			p.advance()
			if _, err := p.expect(kwAS); err != nil {
				return nil, err
			}
			var val string
			switch {
			case p.cur.Type == kwSELF:
				val = "SELF"
			case p.cur.Type == kwOWNER:
				val = "OWNER"
			case p.brokerIsString():
				val = p.cur.Str
			default:
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:EXECUTE AS", Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		default:
			var name string
			for _, w := range []string{"STATUS", "PROCEDURE_NAME", "MAX_QUEUE_READERS"} {
				if p.brokerIsWord(w) {
					name = w
					break
				}
			}
			if name == "" || seen[name] {
				return nil, p.syntaxErrorAtCur()
			}
			seen[name] = true
			p.advance()
			var val string
			switch name {
			case "STATUS":
				v, err := p.brokerExpectOnOff()
				if err != nil {
					return nil, err
				}
				val = v
			case "PROCEDURE_NAME":
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				ref, err := p.brokerParseObjectName(3)
				if err != nil {
					return nil, err
				}
				val = brokerQualifiedName(ref)
			default: // MAX_QUEUE_READERS
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				if p.cur.Type != tokICONST {
					return nil, p.syntaxErrorAtCur()
				}
				val = p.advance().Str
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ACTIVATION:" + name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if create && !(seen["PROCEDURE_NAME"] && seen["MAX_QUEUE_READERS"] && seen["EXECUTE AS"]) {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseCreateServiceStmt parses CREATE SERVICE.
//
// TSql160.g:10327 createServiceStatement and TSql160.g:10357 onQueueClause:
//
//	SERVICE identifier authorizationOpt ON QUEUE queue_name
//	    ( ( identifier [ ,...n ] ) )?
//
// SQL Server takes a queue name of at most two parts here (Msg 117).
func (p *Parser) parseCreateServiceStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// SERVICE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "SERVICE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	var opts []nodes.Node
	authLoc := p.pos()
	owner, err := p.brokerParseAuthorization()
	if err != nil {
		return nil, err
	}
	if owner != "" {
		opts = append(opts, &nodes.ServiceBrokerOption{Name: "AUTHORIZATION", Value: owner, Loc: nodes.Loc{Start: authLoc, End: p.prevEnd()}})
	}

	queue, err := p.brokerParseOnQueue()
	if err != nil {
		return nil, err
	}
	opts = append(opts, queue)

	if p.cur.Type == '(' {
		p.advance()
		for {
			optLoc := p.pos()
			contract, err := p.brokerExpectIdent()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: strings.ToUpper(contract), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// brokerParseOnQueue parses onQueueClause (TSql160.g:10357), ON QUEUE name,
// with a queue name of at most two parts, as SQL Server takes it.
func (p *Parser) brokerParseOnQueue() (*nodes.ServiceBrokerOption, error) {
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	qLoc := p.pos()
	if _, err := p.expect(kwQUEUE); err != nil {
		return nil, err
	}
	ref, err := p.brokerParseObjectName(2)
	if err != nil {
		return nil, err
	}
	return &nodes.ServiceBrokerOption{Name: "QUEUE", Value: ref.Object, Loc: nodes.Loc{Start: qLoc, End: p.prevEnd()}}, nil
}

// parseSendStmt parses SEND ON CONVERSATION.
//
// TSql160.g:25611 sendStatement:
//
//	SEND ON CONVERSATION { ( expression [ ,...n ] ) | expression }
//	    ( MESSAGE TYPE identifierOrVariable )?
//	    ( ( expression ) )?
func (p *Parser) parseSendStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	p.advance() // consume SEND

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "SEND",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCONVERSATION); err != nil {
		return nil, err
	}
	if p.cur.Type == '(' {
		p.advance()
		for i := 0; ; i++ {
			_, text, err := p.brokerParseExprText()
			if err != nil {
				return nil, err
			}
			if i == 0 {
				stmt.Name = text
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else {
		_, text, err := p.brokerParseExprText()
		if err != nil {
			return nil, err
		}
		stmt.Name = text
	}

	if p.cur.Type == kwMESSAGE {
		p.advance()
		if _, err := p.expect(kwTYPE); err != nil {
			return nil, err
		}
		if _, err := p.brokerParseNameOrVariable(); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == '(' {
		p.advance()
		if _, _, err := p.brokerParseExprText(); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseReceiveStmt parses RECEIVE.
//
// TSql160.g:25536 receiveStatement, TSql160.g:25586 receiveSelectExpression
// and TSql160.g:25595 receiveColumnSelectExpression:
//
//	RECEIVE ( TOP ( expression ) )? receiveSelectExpression [ ,...n ]
//	    FROM schemaObjectThreePartName
//	    ( INTO variable )?
//	    ( WHERE { CONVERSATION_GROUP_ID | CONVERSATION_HANDLE } = stringOrVariable )?
//
//	receiveSelectExpression ::= @variable = expression
//	                          | *
//	                          | expression ( [ AS ] stringOrIdentifier )?
func (p *Parser) parseReceiveStmt() (*nodes.ReceiveStmt, error) {
	loc := p.pos()
	p.advance() // consume RECEIVE

	stmt := &nodes.ReceiveStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == kwTOP {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		top, _, err := p.brokerParseExprText()
		if err != nil {
			return nil, err
		}
		stmt.Top = top
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	var cols []nodes.Node
	for {
		col, err := p.brokerParseReceiveColumn()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(cols) == 1 {
		if _, ok := cols[0].(*nodes.ReceiveColumn).Expr.(*nodes.StarExpr); ok {
			stmt.AllColumns = true
			cols = nil
		}
	}
	if len(cols) > 0 {
		stmt.Columns = &nodes.List{Items: cols}
	}

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	queue, err := p.brokerParseObjectName(3)
	if err != nil {
		return nil, err
	}
	stmt.Queue = queue

	if p.cur.Type == kwINTO {
		p.advance()
		if p.cur.Type != tokVARIABLE {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.IntoVar = p.advance().Str
	}

	if p.cur.Type == kwWHERE {
		p.advance()
		whereLoc := p.pos()
		if !p.brokerIsWord("CONVERSATION_GROUP_ID") && !p.brokerIsWord("CONVERSATION_HANDLE") {
			return nil, p.syntaxErrorAtCur()
		}
		colTok := p.advance()
		left := &nodes.ColumnRef{Column: colTok.Str, Loc: nodes.Loc{Start: colTok.Loc, End: p.prevEnd()}}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if !p.brokerIsString() && p.cur.Type != tokVARIABLE {
			return nil, p.syntaxErrorAtCur()
		}
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		stmt.WhereClause = &nodes.BinaryExpr{Op: nodes.BinOpEq, Left: left, Right: right, Loc: nodes.Loc{Start: whereLoc, End: p.prevEnd()}}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// brokerParseReceiveColumn parses one receiveSelectExpression (TSql160.g:25586).
func (p *Parser) brokerParseReceiveColumn() (*nodes.ReceiveColumn, error) {
	col := &nodes.ReceiveColumn{Loc: nodes.Loc{Start: p.pos(), End: -1}}
	switch {
	case p.cur.Type == '*':
		tok := p.advance()
		col.Expr = &nodes.StarExpr{Loc: nodes.Loc{Start: tok.Loc, End: p.prevEnd()}}
	case p.cur.Type == tokVARIABLE && p.peekNext().Type == '=':
		// selectSetVariable: @variable = expression
		v := p.advance()
		p.advance() // consume '='
		val, _, err := p.brokerParseExprText()
		if err != nil {
			return nil, err
		}
		col.Expr = &nodes.SelectAssign{Variable: v.Str, Value: val, Loc: nodes.Loc{Start: v.Loc, End: p.prevEnd()}}
	default:
		expr, _, err := p.brokerParseExprText()
		if err != nil {
			return nil, err
		}
		col.Expr = expr
		// ( [ AS ] stringOrIdentifier )?
		hasAs := false
		if p.cur.Type == kwAS {
			p.advance()
			hasAs = true
		}
		switch {
		case p.brokerIsString() || p.isIdentLike():
			col.Alias = p.advance().Str
		case hasAs:
			return nil, p.syntaxErrorAtCur()
		}
	}
	col.Loc.End = p.prevEnd()
	return col, nil
}

// parseBeginConversationStmt parses BEGIN DIALOG CONVERSATION.
//
// TSql160.g:25678 beginDialogStatement and TSql160.g:25731 beginConversationArgument:
//
//	BEGIN DIALOG [ CONVERSATION ] @dialog_handle
//	    FROM SERVICE identifierOrVariable
//	    TO SERVICE stringOrVariable ( , stringOrVariable )?
//	    ( ON CONTRACT identifierOrVariable )?
//	    ( WITH beginConversationArgument [ ,...n ] )?
//
//	beginConversationArgument ::= { RELATED_CONVERSATION | RELATED_CONVERSATION_GROUP } = expression
//	                            | LIFETIME = expression
//	                            | ENCRYPTION = { ON | OFF }
//
// SQL Server refuses an argument given twice, and both RELATED_ arguments together.
func (p *Parser) parseBeginConversationStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// BEGIN DIALOG [CONVERSATION] already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "BEGIN",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type != tokVARIABLE {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = p.advance().Str

	var opts []nodes.Node

	optLoc := p.pos()
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSERVICE); err != nil {
		return nil, err
	}
	from, err := p.brokerParseNameOrVariable()
	if err != nil {
		return nil, err
	}
	opts = append(opts, &nodes.ServiceBrokerOption{Name: "FROM SERVICE", Value: from, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})

	optLoc = p.pos()
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSERVICE); err != nil {
		return nil, err
	}
	to, err := p.brokerParseStringOrVariable()
	if err != nil {
		return nil, err
	}
	opts = append(opts, &nodes.ServiceBrokerOption{Name: "TO SERVICE", Value: to, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	if _, ok := p.match(','); ok {
		biLoc := p.pos()
		instance, err := p.brokerParseStringOrVariable()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: "BROKER_INSTANCE", Value: instance, Loc: nodes.Loc{Start: biLoc, End: p.prevEnd()}})
	}

	if p.cur.Type == kwON {
		optLoc := p.pos()
		p.advance()
		if _, err := p.expect(kwCONTRACT); err != nil {
			return nil, err
		}
		contract, err := p.brokerParseNameOrVariable()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: "ON CONTRACT", Value: contract, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	}

	if p.cur.Type == kwWITH {
		p.advance()
		seen := map[string]bool{}
		for {
			optLoc := p.pos()
			var name string
			for _, w := range []string{"RELATED_CONVERSATION", "RELATED_CONVERSATION_GROUP", "LIFETIME", "ENCRYPTION"} {
				if p.brokerIsWord(w) {
					name = w
					break
				}
			}
			key := name
			if strings.HasPrefix(name, "RELATED_") {
				key = "RELATED"
			}
			if name == "" || seen[key] {
				return nil, p.syntaxErrorAtCur()
			}
			seen[key] = true
			p.advance()
			var val string
			if name == "ENCRYPTION" {
				v, err := p.brokerExpectOnOff()
				if err != nil {
					return nil, err
				}
				val = v
			} else {
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				_, text, err := p.brokerParseExprText()
				if err != nil {
					return nil, err
				}
				val = text
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseEndConversationStmt parses END CONVERSATION.
//
// TSql160.g:25471 endConversationStatement and TSql160.g:25483 endConversationArgumentsOpt:
//
//	END CONVERSATION expression
//	    ( WITH { CLEANUP | ERROR = integerOrVariable DESCRIPTION = stringOrVariable } )?
func (p *Parser) parseEndConversationStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	p.advance() // consume END
	if _, err := p.expect(kwCONVERSATION); err != nil {
		return nil, err
	}

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "END",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	_, handle, err := p.brokerParseExprText()
	if err != nil {
		return nil, err
	}
	stmt.Name = handle

	if p.cur.Type == kwWITH {
		p.advance()
		var opts []nodes.Node
		switch p.cur.Type {
		case kwERROR:
			errLoc := p.pos()
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			if p.cur.Type != tokICONST && p.cur.Type != tokVARIABLE {
				return nil, p.syntaxErrorAtCur()
			}
			errorCode := p.advance().Str
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "ERROR", Value: errorCode, Loc: nodes.Loc{Start: errLoc, End: p.prevEnd()}})

			descLoc := p.pos()
			if _, err := p.expect(kwDESCRIPTION); err != nil {
				return nil, err
			}
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			descVal, err := p.brokerParseStringOrVariable()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "DESCRIPTION", Value: descVal, Loc: nodes.Loc{Start: descLoc, End: p.prevEnd()}})
		case kwCLEANUP:
			cleanLoc := p.pos()
			p.advance()
			opts = append(opts, &nodes.ServiceBrokerOption{Name: "CLEANUP", Loc: nodes.Loc{Start: cleanLoc, End: p.prevEnd()}})
		default:
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Options = &nodes.List{Items: opts}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// brokerRouteOptions are the route options of CREATE/ALTER ROUTE.
//
// TSql160.g:12615 routeOption: LIFETIME = integer | Identifier = stringLiteral,
// the identifier one of RouteOptionKind's words (CodeGenerationSupporter.cs).
var brokerRouteOptions = []string{"SERVICE_NAME", "BROKER_INSTANCE", "LIFETIME", "ADDRESS", "MIRROR_ADDRESS"}

// brokerParseRouteOptions parses routeOptionList (TSql160.g:12596): WITH
// routeOption [ ,...n ]. SQL Server refuses an option given twice, and CREATE
// ROUTE without ADDRESS.
func (p *Parser) brokerParseRouteOptions(create bool) (*nodes.List, error) {
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	withLoc := p.prevEnd()
	seen := map[string]bool{}
	var opts []nodes.Node
	for {
		optLoc := p.pos()
		var name string
		for _, w := range brokerRouteOptions {
			if p.brokerIsWord(w) {
				name = w
				break
			}
		}
		if name == "" || seen[name] {
			return nil, p.syntaxErrorAtCur()
		}
		seen[name] = true
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		var val string
		if name == "LIFETIME" {
			if p.cur.Type != tokICONST {
				return nil, p.syntaxErrorAtCur()
			}
			val = p.advance().Str
		} else {
			v, err := p.brokerExpectString()
			if err != nil {
				return nil, err
			}
			val = v
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if create && !seen["ADDRESS"] {
		return nil, &ParseError{Message: "CREATE ROUTE requires ADDRESS", Position: withLoc, Source: p.source}
	}
	return &nodes.List{Items: opts}, nil
}

// parseCreateRouteStmt parses CREATE ROUTE.
//
// TSql160.g:10313 createRouteStatement:
//
//	ROUTE identifier authorizationOpt routeOptionList
func (p *Parser) parseCreateRouteStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// ROUTE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "ROUTE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.brokerParseAuthorization(); err != nil {
		return nil, err
	}
	stmt.Options, err = p.brokerParseRouteOptions(true)
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// brokerParseBindingOptions parses the WITH clause of CREATE/ALTER REMOTE
// SERVICE BINDING (TSql160.g:10214, 10264 bindingUserOption, 10275
// bindingAnonymousOption):
//
//	WITH { USER = identifier [ , ANONYMOUS = { ON | OFF } ]
//	     | ANONYMOUS = { ON | OFF } [ , USER = identifier ] }
//
// SQL Server refuses CREATE without USER.
func (p *Parser) brokerParseBindingOptions(create bool) (*nodes.List, error) {
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	withLoc := p.prevEnd()
	seen := map[string]bool{}
	var opts []nodes.Node
	for {
		optLoc := p.pos()
		var name, val string
		switch {
		case p.cur.Type == kwUSER && !seen["USER"]:
			name = "USER"
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			v, err := p.brokerExpectIdent()
			if err != nil {
				return nil, err
			}
			val = v
		case p.brokerIsWord("ANONYMOUS") && !seen["ANONYMOUS"]:
			name = "ANONYMOUS"
			p.advance()
			v, err := p.brokerExpectOnOff()
			if err != nil {
				return nil, err
			}
			val = v
		default:
			return nil, p.syntaxErrorAtCur()
		}
		seen[name] = true
		opts = append(opts, &nodes.ServiceBrokerOption{Name: name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		if len(seen) == 2 {
			break
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if create && !seen["USER"] {
		return nil, &ParseError{Message: "CREATE REMOTE SERVICE BINDING requires USER", Position: withLoc, Source: p.source}
	}
	return &nodes.List{Items: opts}, nil
}

// parseCreateRemoteServiceBindingStmt parses CREATE REMOTE SERVICE BINDING.
//
// TSql160.g:10214 createRemoteServiceBindingStatement:
//
//	REMOTE SERVICE BINDING identifier authorizationOpt TO SERVICE stringLiteral
//	    WITH { USER = identifier [ , ANONYMOUS = { ON | OFF } ] | ANONYMOUS = ... [ , USER = ... ] }
func (p *Parser) parseCreateRemoteServiceBindingStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// REMOTE SERVICE BINDING keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "CREATE",
		ObjectType: "REMOTE SERVICE BINDING",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.brokerParseAuthorization(); err != nil {
		return nil, err
	}

	optLoc := p.pos()
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSERVICE); err != nil {
		return nil, err
	}
	service, err := p.brokerExpectString()
	if err != nil {
		return nil, err
	}
	opts := []nodes.Node{&nodes.ServiceBrokerOption{Name: "SERVICE", Value: service, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}}}

	with, err := p.brokerParseBindingOptions(true)
	if err != nil {
		return nil, err
	}
	opts = append(opts, with.Items...)
	stmt.Options = &nodes.List{Items: opts}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseGetConversationGroupStmt parses GET CONVERSATION GROUP.
//
// TSql160.g:25521 getConversationGroupStatement:
//
//	GET CONVERSATION GROUP variable FROM schemaObjectThreePartName
func (p *Parser) parseGetConversationGroupStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	p.advance() // consume GET

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "GET",
		ObjectType: "CONVERSATION GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if _, err := p.expect(kwCONVERSATION); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}
	if p.cur.Type != tokVARIABLE {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = p.advance().Str

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	qLoc := p.pos()
	ref, err := p.brokerParseObjectName(3)
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.ServiceBrokerOption{Name: "QUEUE", Value: ref.Object, Loc: nodes.Loc{Start: qLoc, End: p.prevEnd()}}}}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterQueueStmt parses ALTER QUEUE.
//
// TSql160.g:12506 alterQueueStatement: QUEUE schemaObjectThreePartName
// queueOptionList. SQL Server also takes the queue actions of
// https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-queue-transact-sql,
// which ScriptDOM lacks:
//
//	REBUILD [ WITH ( MAXDOP = n ) ]
//	| REORGANIZE [ WITH ( LOB_COMPACTION = { ON | OFF } ) ]
//	| MOVE TO { file_group | [default] }
func (p *Parser) parseAlterQueueStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// QUEUE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "QUEUE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	ref, err := p.brokerParseObjectName(3)
	if err != nil {
		return nil, err
	}
	stmt.Name = ref.Object
	if ref.Schema != "" {
		stmt.Name = ref.Schema + "." + ref.Object
	}

	var opts []nodes.Node
	optLoc := p.pos()
	switch p.cur.Type {
	case kwREBUILD, kwREORGANIZE:
		action := strings.ToUpper(p.advance().Str)
		val := ""
		if p.cur.Type == kwWITH {
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			if action == "REBUILD" {
				if _, err := p.expect(kwMAXDOP); err != nil {
					return nil, err
				}
				if _, err := p.expect('='); err != nil {
					return nil, err
				}
				if p.cur.Type != tokICONST {
					return nil, p.syntaxErrorAtCur()
				}
				val = "MAXDOP=" + p.advance().Str
			} else {
				if _, err := p.expect(kwLOB_COMPACTION); err != nil {
					return nil, err
				}
				v, err := p.brokerExpectOnOff()
				if err != nil {
					return nil, err
				}
				val = "LOB_COMPACTION=" + v
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: action, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	case kwMOVE:
		p.advance()
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		fg, err := p.brokerExpectIdent()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: "MOVE TO", Value: fg, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	default:
		opts, err = p.brokerParseQueueOptions(opts, false)
		if err != nil {
			return nil, err
		}
	}

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterServiceStmt parses ALTER SERVICE.
//
// TSql160.g:10379 alterServiceStatement and TSql160.g:10417 addDropServiceContract:
//
//	SERVICE identifier { ON QUEUE queue_name [ ( alter_list ) ] | ( alter_list ) }
//	alter_list ::= { ADD | DROP } CONTRACT identifier [ ,...n ]
func (p *Parser) parseAlterServiceStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// SERVICE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "SERVICE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	var opts []nodes.Node
	if p.cur.Type == kwON {
		queue, err := p.brokerParseOnQueue()
		if err != nil {
			return nil, err
		}
		opts = append(opts, queue)
	} else if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type == '(' {
		p.advance()
		for {
			optLoc := p.pos()
			var action string
			switch p.cur.Type {
			case kwADD:
				action = "ADD CONTRACT"
			case kwDROP:
				action = "DROP CONTRACT"
			default:
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if _, err := p.expect(kwCONTRACT); err != nil {
				return nil, err
			}
			contract, err := p.brokerExpectIdent()
			if err != nil {
				return nil, err
			}
			opts = append(opts, &nodes.ServiceBrokerOption{Name: action, Value: contract, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterRouteStmt parses ALTER ROUTE.
//
// TSql160.g:12636 alterRouteStatement: ROUTE identifier routeOptionList
func (p *Parser) parseAlterRouteStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// ROUTE keyword already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "ROUTE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Options, err = p.brokerParseRouteOptions(false)
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterRemoteServiceBindingStmt parses ALTER REMOTE SERVICE BINDING.
//
// TSql160.g:11372 alterRemoteServiceBindingStatement:
//
//	REMOTE SERVICE BINDING identifier
//	    WITH { USER = identifier [ , ANONYMOUS = { ON | OFF } ] | ANONYMOUS = ... [ , USER = ... ] }
func (p *Parser) parseAlterRemoteServiceBindingStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// REMOTE SERVICE BINDING keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "REMOTE SERVICE BINDING",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Options, err = p.brokerParseBindingOptions(false)
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterMessageTypeStmt parses ALTER MESSAGE TYPE.
//
// TSql160.g:12152 alterMessageTypeStatement: MESSAGE TYPE identifier messageTypeValidation
func (p *Parser) parseAlterMessageTypeStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// MESSAGE TYPE keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "ALTER",
		ObjectType: "MESSAGE TYPE",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opt, err := p.brokerParseMessageTypeValidation()
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: []nodes.Node{opt}}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterContractStmt refuses ALTER CONTRACT: neither TSql160.g's
// alter2005Statements nor SQL Server 2022 has it (a contract's message types
// cannot change; ownership moves with ALTER AUTHORIZATION).
func (p *Parser) parseAlterContractStmt() (*nodes.ServiceBrokerStmt, error) {
	return nil, p.syntaxErrorAtCur()
}

// parseDropServiceBrokerStmt parses DROP for Service Broker objects.
//
// TSql160.g:15499-15555 and 22821: a queue takes schemaObjectThreePartName,
// every other object a one-part identifier (SQL Server: Msg 117, "maximum
// number of prefixes ... is 0").
//
//	DROP { MESSAGE TYPE | CONTRACT | QUEUE | SERVICE | ROUTE | REMOTE SERVICE BINDING | BROKER PRIORITY } name
func (p *Parser) parseDropServiceBrokerStmt(objectType string) (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "DROP",
		ObjectType: objectType,
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	maxParts := 1
	if objectType == "QUEUE" {
		maxParts = 3
	}
	ref, err := p.brokerParseObjectName(maxParts)
	if err != nil {
		return nil, err
	}
	stmt.Name = ref.Object
	if ref.Schema != "" {
		stmt.Name = ref.Schema + "." + ref.Object
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreateBrokerPriorityStmt parses CREATE BROKER PRIORITY.
//
// TSql160.g:22834 brokerPriorityStatementBody:
//
//	BROKER PRIORITY identifier FOR CONVERSATION
//	    ( SET ( brokerPriorityParameter [ ,...n ] ) )?
func (p *Parser) parseCreateBrokerPriorityStmt() (*nodes.ServiceBrokerStmt, error) {
	return p.brokerParsePriority("CREATE", false)
}

// parseAlterBrokerPriorityStmt parses ALTER BROKER PRIORITY, whose SET clause
// SQL Server requires (TSql160.g:22814 shares the CREATE body; the server
// refuses ALTER without SET).
func (p *Parser) parseAlterBrokerPriorityStmt() (*nodes.ServiceBrokerStmt, error) {
	return p.brokerParsePriority("ALTER", true)
}

func (p *Parser) brokerParsePriority(action string, needSet bool) (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	// BROKER PRIORITY keywords already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     action,
		ObjectType: "BROKER PRIORITY",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, err := p.brokerExpectIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwCONVERSATION); err != nil {
		return nil, err
	}

	if p.cur.Type == kwSET {
		stmt.Options, err = p.brokerParsePrioritySet()
		if err != nil {
			return nil, err
		}
	} else if needSet {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// brokerParsePrioritySet parses SET ( brokerPriorityParameter [ ,...n ] ).
//
// TSql160.g:22927 and 22937 brokerPriorityParameter (BrokerPriorityParameterType):
//
//	CONTRACT_NAME = { identifier | ANY }
//	LOCAL_SERVICE_NAME = { identifier | ANY }
//	REMOTE_SERVICE_NAME = { 'string' | ANY }
//	PRIORITY_LEVEL = { integer | DEFAULT }
//
// SQL Server refuses a parameter given twice.
func (p *Parser) brokerParsePrioritySet() (*nodes.List, error) {
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var opts []nodes.Node
	for {
		optLoc := p.pos()
		var name string
		for _, w := range []string{"CONTRACT_NAME", "LOCAL_SERVICE_NAME", "REMOTE_SERVICE_NAME", "PRIORITY_LEVEL"} {
			if p.brokerIsWord(w) {
				name = w
				break
			}
		}
		if name == "" || seen[name] {
			return nil, p.syntaxErrorAtCur()
		}
		seen[name] = true
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		var val string
		switch {
		case name == "PRIORITY_LEVEL" && p.cur.Type == tokICONST:
			val = p.advance().Str
		case name == "PRIORITY_LEVEL" && p.cur.Type == kwDEFAULT:
			p.advance()
			val = "DEFAULT"
		case name != "PRIORITY_LEVEL" && p.cur.Type == kwANY:
			p.advance()
			val = "ANY"
		case name == "REMOTE_SERVICE_NAME" && p.brokerIsString():
			val = p.advance().Str
		case (name == "CONTRACT_NAME" || name == "LOCAL_SERVICE_NAME") && p.isIdentLike():
			val = strings.ToUpper(p.advance().Str)
		default:
			return nil, p.syntaxErrorAtCur()
		}
		opts = append(opts, &nodes.ServiceBrokerOption{Name: name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: opts}, nil
}

// parseMoveConversationStmt parses MOVE CONVERSATION.
//
// TSql160.g:25507 moveConversationStatement: MOVE CONVERSATION expression TO expression
func (p *Parser) parseMoveConversationStmt() (*nodes.ServiceBrokerStmt, error) {
	loc := p.pos()
	p.advance() // consume MOVE

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "MOVE",
		ObjectType: "CONVERSATION",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if _, err := p.expect(kwCONVERSATION); err != nil {
		return nil, err
	}
	_, handle, err := p.brokerParseExprText()
	if err != nil {
		return nil, err
	}
	stmt.Name = handle

	toLoc := p.pos()
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	_, group, err := p.brokerParseExprText()
	if err != nil {
		return nil, err
	}
	stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.ServiceBrokerOption{Name: "TO", Value: group, Loc: nodes.Loc{Start: toLoc, End: p.prevEnd()}}}}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBeginConversationTimerStmt parses a BEGIN CONVERSATION TIMER statement.
//
// TSql160.g:25660 beginConversationTimerStatement:
//
//	BEGIN CONVERSATION TIMER ( expression ) TIMEOUT = expression
func (p *Parser) parseBeginConversationTimerStmt() (*nodes.ServiceBrokerStmt, error) {
	// BEGIN and CONVERSATION TIMER already consumed by caller

	stmt := &nodes.ServiceBrokerStmt{
		Action:     "BEGIN",
		ObjectType: "CONVERSATION TIMER",
	}

	hLoc := p.pos()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	_, handle, err := p.brokerParseExprText()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	opts := []nodes.Node{&nodes.ServiceBrokerOption{Name: "HANDLE", Value: handle, Loc: nodes.Loc{Start: hLoc, End: p.prevEnd()}}}

	tLoc := p.pos()
	if _, err := p.expect(kwTIMEOUT); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	_, timeout, err := p.brokerParseExprText()
	if err != nil {
		return nil, err
	}
	opts = append(opts, &nodes.ServiceBrokerOption{Name: "TIMEOUT", Value: timeout, Loc: nodes.Loc{Start: tLoc, End: p.prevEnd()}})

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
