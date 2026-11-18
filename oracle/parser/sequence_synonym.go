package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateSequenceStmt parses a CREATE SEQUENCE statement.
// The CREATE keyword has already been consumed. The current token is SEQUENCE.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/19/sqlrf/CREATE-SEQUENCE.html
//
//	CREATE SEQUENCE [ IF NOT EXISTS ] [ schema. ] sequence
//	    [ SHARING = { METADATA | DATA | NONE } ]
//	    [ { INCREMENT BY | START WITH } integer
//	    | { MAXVALUE integer | NOMAXVALUE }
//	    | { MINVALUE integer | NOMINVALUE }
//	    | { CYCLE | NOCYCLE }
//	    | { CACHE integer | NOCACHE }
//	    | { ORDER | NOORDER }
//	    | { KEEP | NOKEEP }
//	    | { SCALE [ EXTEND | NOEXTEND ] | NOSCALE }
//	    | { SHARD [ EXTEND | NOEXTEND ] | NOSHARD }
//	    | { SESSION | GLOBAL }
//	    ]... ;
//
// The BNF file places IF NOT EXISTS after the name; Oracle 23ai rejects that
// position (ORA-03049) and accepts it before the name. SHARING must come
// first: after any other option Oracle raises ORA-03049.
func (p *Parser) parseCreateSequenceStmt(start int) (*nodes.CreateSequenceStmt, error) {
	stmt := &nodes.CreateSequenceStmt{
		Loc: nodes.Loc{Start: start},
	}

	// SEQUENCE keyword
	if p.cur.Type == kwSEQUENCE {
		p.advance()
	}

	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance() // consume IF
		p.advance() // consume NOT
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume EXISTS
		stmt.IfNotExists = true
	}
	var parseErr1098 error

	// Sequence name
	stmt.Name, parseErr1098 = p.parseReservedCheckedObjectName()
	if parseErr1098 != nil {
		return nil, parseErr1098
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if p.isKeywordStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type != '=' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume =
		if !p.isKeywordStr("METADATA") && !p.isKeywordStr("DATA") && !p.isKeywordStr("NONE") {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Sharing = p.cur.Str
		p.advance()
	}

	parseErr1099 := p.parseSequenceOptions(stmt)
	if parseErr1099 != nil {
		return nil, parseErr1099
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// sequenceOptionGroup names the CREATE SEQUENCE option group the current
// token starts, or "" when it starts none. Oracle accepts each group at most
// once and rejects a repeated or conflicting one (ORA-02278 through ORA-02285,
// ORA-03178, ORA-41415, ORA-64600).
func (p *Parser) sequenceOptionGroup() string {
	switch p.cur.Type {
	case kwINCREMENT:
		return "INCREMENT"
	case kwSTART:
		return "START"
	case kwMAXVALUE, kwNOMAXVALUE:
		return "MAXVALUE"
	case kwMINVALUE, kwNOMINVALUE:
		return "MINVALUE"
	case kwCYCLE, kwNOCYCLE:
		return "CYCLE"
	case kwCACHE, kwNOCACHE:
		return "CACHE"
	case kwORDER, kwNOORDER:
		return "ORDER"
	case kwKEEP:
		return "KEEP"
	case kwSESSION, kwGLOBAL:
		return "SESSION"
	}
	switch {
	case p.isKeywordStr("NOKEEP"):
		return "KEEP"
	case p.isKeywordStr("SCALE"), p.isKeywordStr("NOSCALE"):
		return "SCALE"
	case p.isKeywordStr("SHARD"), p.isKeywordStr("NOSHARD"):
		return "SHARD"
	}
	return ""
}

// parseSequenceOptions parses the various options for CREATE SEQUENCE.
func (p *Parser) parseSequenceOptions(stmt *nodes.CreateSequenceStmt) error {
	seen := make(map[string]bool)
	extendSeen := false
	for {
		group := p.sequenceOptionGroup()
		if group == "" {
			return nil
		}
		if seen[group] {
			return p.syntaxErrorAtCur()
		}
		seen[group] = true

		switch {
		case p.cur.Type == kwINCREMENT:
			p.advance()
			if err := p.expectCreateTableToken(kwBY); err != nil {
				return err
			}
			var parseErr1100 error
			stmt.IncrementBy, parseErr1100 = p.parseExpr()
			if parseErr1100 != nil {
				return parseErr1100
			}
			if stmt.IncrementBy == nil {
				return p.syntaxErrorAtCur()
			}
		case p.cur.Type == kwSTART:
			p.advance()
			if err := p.expectCreateTableToken(kwWITH); err != nil {
				return err
			}
			var parseErr1101 error
			stmt.StartWith, parseErr1101 = p.parseExpr()
			if parseErr1101 != nil {
				return parseErr1101
			}
			if stmt.StartWith == nil {
				return p.syntaxErrorAtCur()
			}
		case p.cur.Type == kwMAXVALUE:
			p.advance()
			var parseErr1102 error
			stmt.MaxValue, parseErr1102 = p.parseExpr()
			if parseErr1102 != nil {
				return parseErr1102
			}
			if stmt.MaxValue == nil {
				return p.syntaxErrorAtCur()
			}
		case p.cur.Type == kwNOMAXVALUE:
			stmt.NoMaxValue = true
			p.advance()
		case p.cur.Type == kwMINVALUE:
			p.advance()
			var parseErr1103 error
			stmt.MinValue, parseErr1103 = p.parseExpr()
			if parseErr1103 != nil {
				return parseErr1103
			}
			if stmt.MinValue == nil {
				return p.syntaxErrorAtCur()
			}
		case p.cur.Type == kwNOMINVALUE:
			stmt.NoMinValue = true
			p.advance()
		case p.cur.Type == kwCYCLE:
			stmt.Cycle = true
			p.advance()
		case p.cur.Type == kwNOCYCLE:
			stmt.NoCycle = true
			p.advance()
		case p.cur.Type == kwCACHE:
			p.advance()
			var parseErr1104 error
			stmt.Cache, parseErr1104 = p.parseExpr()
			if parseErr1104 != nil {
				return parseErr1104
			}
			if stmt.Cache == nil {
				return p.syntaxErrorAtCur()
			}
		case p.cur.Type == kwNOCACHE:
			stmt.NoCache = true
			p.advance()
		case p.cur.Type == kwORDER:
			stmt.Order = true
			p.advance()
		case p.cur.Type == kwNOORDER:
			stmt.NoOrder = true
			p.advance()
		case p.cur.Type == kwKEEP:
			stmt.Keep = true
			p.advance()
		case p.isKeywordStr("NOKEEP"):
			stmt.NoKeep = true
			p.advance()
		case p.isKeywordStr("SCALE"):
			stmt.Scale = true
			p.advance()
			var err error
			stmt.ScaleExtend, stmt.ScaleNoExtend, err = p.parseSequenceExtendModifier(&extendSeen)
			if err != nil {
				return err
			}
		case p.isKeywordStr("NOSCALE"):
			stmt.NoScale = true
			p.advance()
		case p.isKeywordStr("SHARD"):
			stmt.Shard = true
			p.advance()
			var err error
			stmt.ShardExtend, stmt.ShardNoExtend, err = p.parseSequenceExtendModifier(&extendSeen)
			if err != nil {
				return err
			}
		case p.isKeywordStr("NOSHARD"):
			stmt.NoShard = true
			p.advance()
		case p.cur.Type == kwSESSION:
			stmt.Session = true
			p.advance()
		case p.cur.Type == kwGLOBAL:
			stmt.Global = true
			p.advance()
		}
	}
}

// parseSequenceExtendModifier consumes the EXTEND or NOEXTEND that may follow
// SCALE or SHARD in CREATE and ALTER SEQUENCE and reports which it was. With
// both SCALE and SHARD, one modifier applies to both; writing it for each,
// with the same or a different value, is a parsing error ("duplicate or
// conflicting EXTEND clause"). seen records a modifier already consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-SEQUENCE.html
// (Sequence with SHARD and SCALE). A non-sharded Oracle raises ORA-02511 at
// SHARD before reaching this check, so the rule rests on the documentation.
func (p *Parser) parseSequenceExtendModifier(seen *bool) (extend, noExtend bool, err error) {
	if !p.isKeywordStr("EXTEND") && !p.isKeywordStr("NOEXTEND") {
		return false, false, nil
	}
	if *seen {
		return false, false, p.syntaxErrorAtCur()
	}
	*seen = true
	extend = p.cur.Str == "EXTEND"
	p.advance() // consume EXTEND or NOEXTEND
	return extend, !extend, nil
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

	// [ EDITIONABLE | NONEDITIONABLE ]
	if p.isKeywordStr("EDITIONABLE") || p.isKeywordStr("NONEDITIONABLE") {
		stmt.Edition = p.cur.Str
		p.advance()
	}
	// [ PUBLIC ]
	if p.cur.Type == kwPUBLIC {
		stmt.Public = true
		p.advance()
	}
	// SYNONYM keyword
	if p.cur.Type == kwSYNONYM {
		p.advance()
	}
	var err error
	if stmt.IfNotExists, err = p.parseIfNotExistsAfterKind(orReplace); err != nil {
		return nil, err
	}
	var parseErr1105 error

	// Synonym name
	stmt.Name, parseErr1105 = p.parseReservedCheckedObjectName()
	if parseErr1105 !=

		// FOR target
		nil {
		return nil, parseErr1105
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	// [ SHARING = NONE ]
	if p.isKeywordStr("SHARING") {
		toks := p.peekAhead(3)
		if len(toks) == 3 && toks[0].Type == '=' && p.isIdentLikeStrAt(toks[1], "NONE") {
			sharing, err := p.parseSharingClause()
			if err != nil {
				return nil, err
			}
			stmt.Sharing = sharing
		}
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

// parseDatabaseLinkName parses a database link name: a name, its domain
// after dots, and a connection qualifier after @.
func (p *Parser) parseDatabaseLinkName() (string, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", p.syntaxErrorAtCur()
	}
	for p.cur.Type == '.' || p.cur.Type == '@' {
		sep := string(rune(p.cur.Type))
		p.advance()
		part, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		if part == "" {
			return "", p.syntaxErrorAtCur()
		}
		name += sep + part
	}
	return name, nil
}

// parseDatabaseLinkUser parses user IDENTIFIED BY password. Oracle 23ai
// takes the password as a name, quoted or not, and not as a string
// (ORA-00988).
func (p *Parser) parseDatabaseLinkUser() (string, string, error) {
	user, err := p.parseIdentifier()
	if err != nil {
		return "", "", err
	}
	if user == "" {
		return "", "", p.syntaxErrorAtCur()
	}
	if p.cur.Type != kwIDENTIFIED {
		return "", "", p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != kwBY {
		return "", "", p.syntaxErrorAtCur()
	}
	p.advance()
	password, err := p.parseIdentifier()
	if err != nil {
		return "", "", err
	}
	if password == "" {
		return "", "", p.syntaxErrorAtCur()
	}
	return user, password, nil
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

	// DATABASE LINK
	if p.cur.Type != kwDATABASE {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != kwLINK {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// IF NOT EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance() // IF
		not := p.cur
		p.advance() // NOT
		if p.cur.Type != kwEXISTS {
			// Oracle reads IF as the name then (ORA-02010 at NOT).
			return nil, p.syntaxErrorAtTok(not)
		}
		p.advance() // EXISTS
		stmt.IfNotExists = true
	}

	// Link name
	var err error
	stmt.Name, err = p.parseDatabaseLinkName()
	if err != nil {
		return nil, err
	}

	// CONNECT { TO { CURRENT_USER | user IDENTIFIED BY password } | WITH credential }
	if p.cur.Type == kwCONNECT {
		p.advance()
		switch {
		case p.cur.Type == kwTO:
			p.advance()
			if p.isKeywordStr("CURRENT_USER") {
				stmt.ConnectTo = p.cur.Str
				p.advance()
				break
			}
			stmt.ConnectTo, stmt.Identified, err = p.parseDatabaseLinkUser()
			if err != nil {
				return nil, err
			}
		case p.cur.Type == kwWITH:
			p.advance()
			credential, err := p.parseObjectName()
			if err != nil {
				return nil, err
			}
			if credential.Name == "" {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.ConnectWith = credential
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	// USING 'connect_string'
	if p.cur.Type == kwUSING {
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Using = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
