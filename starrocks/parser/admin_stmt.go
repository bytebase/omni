package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// ADMIN statements
// ---------------------------------------------------------------------------

// parseAdminStmt parses an ADMIN statement. ADMIN has been consumed by the
// caller; cur is the word after it. The forms are StarRocks.g4's admin
// statements (StarRocks.g4:694-742), updateFailPointStatusStatement
// (StarRocks.g4:830) and executeScriptStatement (StarRocks.g4:2001); any
// other word after ADMIN is a syntax error.
func (p *Parser) parseAdminStmt(adminLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.AdminStmt{}
	var err error
	switch {
	case p.cur.Kind == kwSHOW:
		err = p.parseAdminShow(stmt)
	case p.cur.Kind == kwSET:
		err = p.parseAdminSet(stmt)
	case p.cur.Kind == kwREPAIR:
		// ADMIN REPAIR TABLE qualifiedName partitionNames?
		p.adminVerb(stmt, "TABLE")
		err = p.parseAdminTableAndPartitions(stmt, kwTABLE)
	case p.cur.Kind == kwCANCEL:
		// ADMIN CANCEL REPAIR TABLE qualifiedName partitionNames?
		p.adminVerb(stmt, "REPAIR TABLE")
		if _, err = p.expect(kwREPAIR); err == nil {
			err = p.parseAdminTableAndPartitions(stmt, kwTABLE)
		}
	case p.cur.Kind == kwCHECK:
		// ADMIN CHECK tabletList PROPERTIES '(' property ')': one property,
		// and the clause is required.
		p.adminVerb(stmt, "TABLET")
		if _, err = p.expect(kwTABLET); err == nil {
			start := p.cur.Loc.Start
			if err = p.parseAdminIntegerList(); err == nil {
				p.adminArgs(stmt, start)
				if _, err = p.expect(kwPROPERTIES); err == nil {
					err = p.parseAdminOneProperty(stmt)
				}
			}
		}
	case p.cur.Kind == kwENABLE || p.adminIsWord("DISABLE"):
		err = p.parseAdminFailPoint(stmt)
	case p.cur.Kind == kwEXECUTE:
		// ADMIN EXECUTE ON (FRONTEND | INTEGER_VALUE) string
		p.adminVerb(stmt, "")
		if _, err = p.expect(kwON); err != nil {
			break
		}
		switch p.cur.Kind {
		case kwFRONTEND:
			stmt.Object = "FRONTEND"
			p.advance()
		case tokInt:
			stmt.Object = "BACKEND"
			p.adminArgs(stmt, p.advance().Loc.Start)
		default:
			err = p.syntaxErrorAtCur()
		}
		if err == nil {
			err = p.parseAdminString(stmt)
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	stmt.Loc = ast.Loc{Start: adminLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseAdminShow reads the ADMIN SHOW forms. cur is SHOW.
func (p *Parser) parseAdminShow(stmt *ast.AdminStmt) error {
	p.adminVerb(stmt, "")
	switch p.cur.Kind {
	case kwFRONTEND:
		// ADMIN SHOW FRONTEND CONFIG (LIKE string)? (StarRocks.g4:700)
		stmt.Object = "FRONTEND CONFIG"
		p.advance() // consume FRONTEND
		if _, err := p.expect(kwCONFIG); err != nil {
			return err
		}
		if p.cur.Kind == kwLIKE {
			start := p.advance().Loc.Start // consume LIKE
			if _, err := p.expect(tokString); err != nil {
				return err
			}
			p.adminArgs(stmt, start)
		}
		return nil
	case kwREPLICA:
		p.advance() // consume REPLICA
		switch p.cur.Kind {
		case kwSTATUS:
			// ADMIN SHOW REPLICA STATUS FROM qualifiedName partitionNames?
			// (WHERE expression)? (StarRocks.g4:708)
			stmt.Object = "REPLICA STATUS"
			p.advance() // consume STATUS
			if err := p.parseAdminTableAndPartitions(stmt, kwFROM); err != nil {
				return err
			}
			if p.cur.Kind == kwWHERE {
				start := p.advance().Loc.Start // consume WHERE
				if _, err := p.parseExpr(); err != nil {
					return err
				}
				p.adminArgs(stmt, start)
			}
			return nil
		case kwDISTRIBUTION:
			// ADMIN SHOW REPLICA DISTRIBUTION FROM qualifiedName
			// partitionNames? (StarRocks.g4:704)
			stmt.Object = "REPLICA DISTRIBUTION"
			p.advance() // consume DISTRIBUTION
			return p.parseAdminTableAndPartitions(stmt, kwFROM)
		}
	}
	return p.syntaxErrorAtCur()
}

// parseAdminSet reads the ADMIN SET forms. cur is SET.
func (p *Parser) parseAdminSet(stmt *ast.AdminStmt) error {
	p.adminVerb(stmt, "")
	switch {
	case p.cur.Kind == kwFRONTEND:
		// ADMIN SET FRONTEND CONFIG '(' property ')' (StarRocks.g4:694): one
		// property, where Doris takes a list.
		stmt.Object = "FRONTEND CONFIG"
		p.advance() // consume FRONTEND
		if _, err := p.expect(kwCONFIG); err != nil {
			return err
		}
		return p.parseAdminOneProperty(stmt)
	case p.cur.Kind == kwREPLICA:
		// ADMIN SET REPLICA STATUS properties (StarRocks.g4:697)
		stmt.Object = "REPLICA STATUS"
		p.advance() // consume REPLICA
		if _, err := p.expect(kwSTATUS); err != nil {
			return err
		}
		if p.cur.Kind != kwPROPERTIES {
			return p.syntaxErrorAtCur()
		}
		props, err := p.parseProperties()
		if err != nil {
			return err
		}
		stmt.Properties = props
		return nil
	case p.cur.Kind == kwTABLE:
		// ADMIN SET TABLE qualifiedName PARTITION '(' (identifierOrString |
		// INTEGER_VALUE) ')' VERSION TO INTEGER_VALUE (StarRocks.g4:724)
		stmt.Object = "TABLE PARTITION VERSION"
		p.advance() // consume TABLE
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return err
		}
		stmt.Target = name
		start := p.cur.Loc.Start
		if err := p.expectAdminWords(kwPARTITION, int('(')); err != nil {
			return err
		}
		if p.cur.Kind == tokInt {
			p.advance()
		} else if _, _, err := p.parseIdentifierOrString(); err != nil {
			return err
		}
		if err := p.expectAdminWords(int(')'), kwVERSION, kwTO, tokInt); err != nil {
			return err
		}
		p.adminArgs(stmt, start)
		return nil
	case p.adminIsWord("AUTOMATED"):
		// ADMIN SET AUTOMATED CLUSTER SNAPSHOT ON (STORAGE VOLUME
		// identifier)? | OFF (StarRocks.g4:736-742)
		stmt.Object = "AUTOMATED CLUSTER SNAPSHOT"
		p.advance() // consume AUTOMATED
		if err := p.expectAdminWords(kwCLUSTER, kwSNAPSHOT); err != nil {
			return err
		}
		start := p.cur.Loc.Start
		switch p.cur.Kind {
		case kwOFF:
			p.advance()
		case kwON:
			p.advance()
			if p.cur.Kind == kwSTORAGE {
				p.advance() // consume STORAGE
				if !p.adminIsWord("VOLUME") {
					return p.syntaxErrorAtCur()
				}
				p.advance() // consume VOLUME
				if _, _, err := p.parseIdentifier(); err != nil {
					return err
				}
			}
		default:
			return p.syntaxErrorAtCur()
		}
		p.adminArgs(stmt, start)
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseAdminFailPoint reads updateFailPointStatusStatement
// (StarRocks.g4:830):
//
//	ADMIN DISABLE FAILPOINT string (ON BACKEND string)?
//	ADMIN ENABLE FAILPOINT string (WITH INTEGER_VALUE TIMES)? (ON BACKEND string)?
//	ADMIN ENABLE FAILPOINT string (WITH DECIMAL_VALUE PROBABILITY)? (ON BACKEND string)?
//
// cur is ENABLE or DISABLE.
func (p *Parser) parseAdminFailPoint(stmt *ast.AdminStmt) error {
	enable := p.cur.Kind == kwENABLE
	p.adminVerb(stmt, "FAILPOINT")
	if !p.adminIsWord("FAILPOINT") {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume FAILPOINT
	if err := p.parseAdminString(stmt); err != nil {
		return err
	}
	start := p.cur.Loc.Start
	if enable && p.cur.Kind == kwWITH {
		p.advance() // consume WITH
		switch p.cur.Kind {
		case tokInt:
			p.advance()
			if !p.adminIsWord("TIMES") {
				return p.syntaxErrorAtCur()
			}
			p.advance() // consume TIMES
		case tokFloat:
			// DECIMAL_VALUE has no exponent (StarRocksLex.g4:534). The AST
			// builder refuses a probability outside [0, 1]: StarRocks
			// rejects WITH 1.5 PROBABILITY (compat.tsv).
			tok := p.cur
			v, err := strconv.ParseFloat(tok.Str, 64)
			if strings.ContainsAny(tok.Str, "eE") || err != nil || v < 0 || v > 1 {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			if !p.adminIsWord("PROBABILITY") {
				return p.syntaxErrorAtCur()
			}
			p.advance() // consume PROBABILITY
		default:
			return p.syntaxErrorAtCur()
		}
	}
	if p.cur.Kind == kwON {
		p.advance() // consume ON
		if _, err := p.expect(kwBACKEND); err != nil {
			return err
		}
		if _, err := p.expect(tokString); err != nil {
			return err
		}
	}
	p.adminArgs(stmt, start)
	return nil
}

// adminVerb records the verb at cur (and the object, when the form fixes one)
// and consumes the verb.
func (p *Parser) adminVerb(stmt *ast.AdminStmt, object string) {
	stmt.Verb = strings.ToUpper(p.advance().Str)
	stmt.Object = object
}

// adminIsWord reports whether cur is the unquoted word w, for the words of
// the ADMIN grammar this lexer has no keyword for.
func (p *Parser) adminIsWord(w string) bool {
	return p.cur.Kind == tokIdent && strings.EqualFold(p.cur.Str, w)
}

// expectAdminWords consumes kinds in order.
func (p *Parser) expectAdminWords(kinds ...TokenKind) error {
	for _, k := range kinds {
		if _, err := p.expect(k); err != nil {
			return err
		}
	}
	return nil
}

// adminArgs appends the source text from start to the last consumed token to
// stmt.Args, which holds what a form takes beyond its verb, object, target
// and properties.
func (p *Parser) adminArgs(stmt *ast.AdminStmt, start int) {
	end := p.prev.Loc.End
	if end <= start {
		return
	}
	text := p.input[start-p.baseOffset : end-p.baseOffset]
	if stmt.Args != "" {
		stmt.Args += " "
	}
	stmt.Args += text
}

// parseAdminString reads the string operand of a form into stmt.Args.
func (p *Parser) parseAdminString(stmt *ast.AdminStmt) error {
	tok, err := p.expect(tokString)
	if err != nil {
		return err
	}
	p.adminArgs(stmt, tok.Loc.Start)
	return nil
}

// parseAdminOneProperty reads '(' property ')', a single string pair, into
// stmt.Properties.
func (p *Parser) parseAdminOneProperty(stmt *ast.AdminStmt) error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	prop, err := p.parseStringProperty()
	if err != nil {
		return err
	}
	if _, err := p.expect(int(')')); err != nil {
		return err
	}
	stmt.Properties = []*ast.Property{prop}
	return nil
}

// parseAdminIntegerList reads '(' INTEGER_VALUE (',' INTEGER_VALUE)* ')', the
// list of tabletList (StarRocks.g4:2278).
func (p *Parser) parseAdminIntegerList() error {
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	for {
		if _, err := p.expect(tokInt); err != nil {
			return err
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(int(')'))
	return err
}

// parseAdminTableAndPartitions reads lead qualifiedName partitionNames?: the
// name into stmt.Target, the partitions' text into stmt.Args.
//
//	partitionNames
//	    : TEMPORARY? (PARTITION | PARTITIONS) '(' identifierOrString (',' identifierOrString)* ')'
//	    | TEMPORARY? (PARTITION | PARTITIONS) identifierOrString
//	    | PARTITION '(' identifier '=' literalExpression (',' ...)* ')'
//
// (StarRocks.g4:2268).
func (p *Parser) parseAdminTableAndPartitions(stmt *ast.AdminStmt, lead TokenKind) error {
	if _, err := p.expect(lead); err != nil {
		return err
	}
	name, err := p.parseMultipartIdentifier()
	if err != nil {
		return err
	}
	stmt.Target = name

	start := p.cur.Loc.Start
	temporary := false
	if p.cur.Kind == kwTEMPORARY {
		if next := p.peekNext().Kind; next == kwPARTITION || next == kwPARTITIONS {
			p.advance() // consume TEMPORARY
			temporary = true
		}
	}
	if p.cur.Kind != kwPARTITION && p.cur.Kind != kwPARTITIONS {
		return nil
	}
	keyed := p.cur.Kind == kwPARTITION && !temporary
	p.advance() // consume PARTITION / PARTITIONS
	if p.cur.Kind != int('(') {
		if _, _, err := p.parseIdentifierOrString(); err != nil {
			return err
		}
		p.adminArgs(stmt, start)
		return nil
	}
	p.advance() // consume '('
	for {
		// A keyPartition (StarRocks.g4:2668) is identifier '=' literal.
		ident := isIdentifierToken(p.cur.Kind)
		if _, _, err := p.parseIdentifierOrString(); err != nil {
			return err
		}
		if keyed && ident && p.cur.Kind == int('=') {
			p.advance() // consume '='
			if err := p.parseAdminPartitionLiteral(); err != nil {
				return err
			}
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(int(')')); err != nil {
		return err
	}
	p.adminArgs(stmt, start)
	return nil
}

// parseAdminPartitionLiteral reads the value of a keyPartition, a
// literalExpression (StarRocks.g4:2428): NULL, a boolean, a number, a string
// or a DATE / DATETIME string. Its interval, binary and parameter forms are
// refused, failing closed.
func (p *Parser) parseAdminPartitionLiteral() error {
	switch p.cur.Kind {
	case kwNULL, kwTRUE, kwFALSE, tokInt, tokFloat, tokString:
		p.advance()
		return nil
	case int('-'):
		p.advance() // consume '-'
		if p.cur.Kind != tokInt && p.cur.Kind != tokFloat {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	case kwDATE, kwDATETIME:
		p.advance() // consume DATE / DATETIME
		_, err := p.expect(tokString)
		return err
	}
	return p.syntaxErrorAtCur()
}
