package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateSequenceStmt parses a CREATE SEQUENCE statement.
// The CREATE keyword has already been consumed. The current token is SEQUENCE.
//
// BNF: oracle/parser/bnf/CREATE-SEQUENCE.bnf
//
//	CREATE SEQUENCE [ schema. ] sequence
//	    [ IF NOT EXISTS ]
//	    [ SHARING = { METADATA | DATA | NONE } ]
//	    [ INCREMENT BY integer ]
//	    [ START WITH integer ]
//	    [ { MAXVALUE integer | NOMAXVALUE } ]
//	    [ { MINVALUE integer | NOMINVALUE } ]
//	    [ { CYCLE | NOCYCLE } ]
//	    [ { CACHE integer | NOCACHE } ]
//	    [ { ORDER | NOORDER } ]
//	    [ { KEEP | NOKEEP } ]
//	    [ { SCALE { EXTEND | NOEXTEND } | NOSCALE } ]
//	    [ { SESSION | GLOBAL } ]
//	    [ SHARD ] ;
func (p *Parser) parseCreateSequenceStmt(start int) (*nodes.CreateSequenceStmt, error) {
	stmt := &nodes.CreateSequenceStmt{
		Loc: nodes.Loc{Start: start},
	}

	// SEQUENCE keyword
	if p.cur.Type == kwSEQUENCE {
		p.advance()
	}
	var parseErr1098 error

	// Sequence name
	stmt.Name, parseErr1098 = p.parseObjectName()
	if parseErr1098 !=

		// Parse sequence options
		nil {
		return nil, parseErr1098
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}
	parseErr1099 := p.parseSequenceOptions(stmt)
	if parseErr1099 != nil {
		return nil, parseErr1099
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSequenceOptions parses the various options for CREATE SEQUENCE.
func (p *Parser) parseSequenceOptions(stmt *nodes.CreateSequenceStmt) error {
	for {
		switch p.cur.Type {
		case kwINCREMENT:
			p.advance()
			if p.cur.Type == kwBY {
				p.advance()
			}
			var parseErr1100 error
			stmt.IncrementBy, parseErr1100 = p.parseExpr()
			if parseErr1100 != nil {
				return parseErr1100
			}
		case kwSTART:
			p.advance()
			if p.cur.Type == kwWITH {
				p.advance()
			}
			var parseErr1101 error
			stmt.StartWith, parseErr1101 = p.parseExpr()
			if parseErr1101 != nil {
				return parseErr1101
			}
		case kwMAXVALUE:
			p.advance()
			var parseErr1102 error
			stmt.MaxValue, parseErr1102 = p.parseExpr()
			if parseErr1102 != nil {
				return parseErr1102
			}
		case kwNOMAXVALUE:
			stmt.NoMaxValue = true
			p.advance()
		case kwMINVALUE:
			p.advance()
			var parseErr1103 error
			stmt.MinValue, parseErr1103 = p.parseExpr()
			if parseErr1103 != nil {
				return parseErr1103
			}
		case kwNOMINVALUE:
			stmt.NoMinValue = true
			p.advance()
		case kwCYCLE:
			stmt.Cycle = true
			p.advance()
		case kwNOCYCLE:
			stmt.NoCycle = true
			p.advance()
		case kwCACHE:
			p.advance()
			var parseErr1104 error
			stmt.Cache, parseErr1104 = p.parseExpr()
			if parseErr1104 != nil {
				return parseErr1104
			}
		case kwNOCACHE:
			stmt.NoCache = true
			p.advance()
		case kwORDER:
			stmt.Order = true
			p.advance()
		case kwNOORDER:
			stmt.NoOrder = true
			p.advance()
		default:
			return nil
		}
	}
	return nil
}

// parseCreateSynonymStmt parses a CREATE [OR REPLACE] [PUBLIC] SYNONYM statement.
// The CREATE keyword has already been consumed. The caller has already parsed
// OR REPLACE and PUBLIC if present and passes them in.
//
// BNF: oracle/parser/bnf/CREATE-SYNONYM.bnf
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ]
//	    [ EDITIONABLE | NONEDITIONABLE ]
//	    [ PUBLIC ] SYNONYM [ schema. ] synonym
//	    [ SHARING = { METADATA | NONE } ]
//	    FOR [ schema. ] object [ @ dblink ] ;
func (p *Parser) parseCreateSynonymStmt(start int, orReplace, public bool) (*nodes.CreateSynonymStmt, error) {
	stmt := &nodes.CreateSynonymStmt{
		OrReplace: orReplace,
		Public:    public,
		Loc:       nodes.Loc{Start: start},
	}

	// SYNONYM keyword
	if p.cur.Type == kwSYNONYM {
		p.advance()
	}
	var parseErr1105 error

	// Synonym name
	stmt.Name, parseErr1105 = p.parseObjectName()
	if parseErr1105 !=

		// FOR target
		nil {
		return nil, parseErr1105
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type != kwFOR {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var parseErr1106 error
	stmt.Target, parseErr1106 = p.parseObjectName()
	if parseErr1106 != nil {
		return nil, parseErr1106
	}
	if stmt.Target == nil || stmt.Target.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateDatabaseLinkStmt parses a CREATE [PUBLIC] DATABASE LINK statement.
// The CREATE keyword has already been consumed. The caller has already parsed
// PUBLIC if present and passes it in.
//
// BNF: oracle/parser/bnf/CREATE-DATABASE-LINK.bnf
//
//	CREATE [ SHARED ] [ PUBLIC ] DATABASE LINK [ IF NOT EXISTS ] dblink
//	    [ CONNECT TO
//	        { CURRENT_USER
//	        | user IDENTIFIED BY password [ dblink_authentication ]
//	        }
//	    | CONNECT WITH credential
//	    ]
//	    [ dblink_authentication ]
//	    USING 'connect_string' ;
//
//	dblink_authentication:
//	    AUTHENTICATED BY user IDENTIFIED BY password
func (p *Parser) parseCreateDatabaseLinkStmt(start int, public bool) (*nodes.CreateDatabaseLinkStmt, error) {
	stmt := &nodes.CreateDatabaseLinkStmt{
		Public: public,
		Loc:    nodes.Loc{Start: start},
	}

	// DATABASE keyword
	if p.cur.Type == kwDATABASE {
		p.advance()
	}

	// LINK keyword
	if p.cur.Type == kwLINK {
		p.advance()
	}
	var parseErr1107 error

	// Link name
	stmt.Name, parseErr1107 = p.parseIdentifier()
	if parseErr1107 !=

		// CONNECT TO user IDENTIFIED BY password
		nil {
		return nil, parseErr1107
	}

	if p.cur.Type == kwCONNECT {
		p.advance()
		if p.cur.Type == kwTO {
			p.advance()
		}
		var parseErr1108 error
		stmt.ConnectTo, parseErr1108 = p.parseIdentifier()
		if parseErr1108 != nil {
			return nil, parseErr1108
		}

		if p.cur.Type == kwIDENTIFIED {
			p.advance()
			if p.cur.Type == kwBY {
				p.advance()
			}
			var parseErr1109 error
			stmt.Identified, parseErr1109 = p.parseIdentifier()
			if parseErr1109 !=

				// USING 'connect_string'
				nil {
				return nil, parseErr1109
			}
		}
	}

	if p.cur.Type == kwUSING {
		p.advance()
		if p.cur.Type == tokSCONST {
			stmt.Using = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
