// Package parser implements a recursive descent SQL parser for T-SQL (SQL Server).
//
// This parser reuses the lexer and keyword definitions from this package
// and produces AST nodes from the mssql/ast package.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// Parser is a recursive descent parser for T-SQL.
type Parser struct {
	lexer   *Lexer
	cur     Token // current token
	prev    Token // previous token (for error reporting)
	nextBuf Token // buffered next token for 2-token lookahead
	hasNext bool  // whether nextBuf is valid
}

// Parse parses a T-SQL string into an AST list.
// Currently supports basic infrastructure; statement dispatch will be
// implemented incrementally across batches.
func Parse(sql string) (*nodes.List, error) {
	p := &Parser{
		lexer: NewLexer(sql),
	}
	p.advance()

	var stmts []nodes.Node
	for p.cur.Type != tokEOF {
		// Skip semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		stmt := p.parseStmt()
		if stmt == nil {
			if p.cur.Type != tokEOF {
				return nil, &ParseError{
					Message:  "unexpected token in statement",
					Position: p.cur.Loc,
				}
			}
			break
		}
		stmts = append(stmts, stmt)
	}

	if len(stmts) == 0 {
		return &nodes.List{}, nil
	}
	return &nodes.List{Items: stmts}, nil
}

// parseStmt dispatches to statement-specific parsers.
func (p *Parser) parseStmt() nodes.StmtNode {
	switch p.cur.Type {
	case kwSELECT:
		return p.parseSelectStmt()
	case kwWITH:
		return p.parseSelectStmt()
	case kwINSERT:
		return p.parseInsertStmt()
	case kwUPDATE:
		// Could be UPDATE table... or UPDATE STATISTICS
		next := p.peekNext()
		if next.Type == kwSTATISTICS {
			loc := p.pos()
			p.advance() // consume UPDATE
			p.advance() // consume STATISTICS
			stmt := p.parseUpdateStatisticsStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		return p.parseUpdateStmt()
	case kwDELETE:
		return p.parseDeleteStmt()
	case kwMERGE:
		return p.parseMergeStmt()
	case kwCREATE:
		return p.parseCreateStmt()
	case kwALTER:
		return p.parseAlterStmt()
	case kwDROP:
		return p.parseDropOrSecurityStmt()
	case kwTRUNCATE:
		return p.parseTruncateStmt()
	case kwIF:
		return p.parseIfStmt()
	case kwWHILE:
		return p.parseWhileStmt()
	case kwBEGIN:
		return p.parseBeginStmt()
	case kwRETURN:
		return p.parseReturnStmt()
	case kwBREAK:
		return p.parseBreakStmt()
	case kwCONTINUE:
		return p.parseContinueStmt()
	case kwGOTO:
		return p.parseGotoStmt()
	case kwWAITFOR:
		return p.parseWaitForStmt()
	case kwDECLARE:
		// Check if this is DECLARE cursor_name CURSOR (named cursor declaration).
		// Named cursors use a plain identifier (not @variable) after DECLARE.
		next := p.peekNext()
		if next.Type != tokVARIABLE {
			return p.parseDeclareCursorStmt()
		}
		return p.parseDeclareStmt()
	case kwSET:
		return p.parseSetStmt()
	case kwCOMMIT:
		return p.parseCommitStmt()
	case kwROLLBACK:
		return p.parseRollbackStmt()
	case kwSAVE:
		return p.parseSaveTransStmt()
	case kwEXEC, kwEXECUTE:
		return p.parseExecStmt()
	case kwGRANT:
		return p.parseGrantStmt()
	case kwREVOKE:
		return p.parseRevokeStmt()
	case kwDENY:
		return p.parseDenyStmt()
	case kwUSE:
		return p.parseUseStmt()
	case kwPRINT:
		return p.parsePrintStmt()
	case kwRAISERROR:
		return p.parseRaiseErrorStmt()
	case kwTHROW:
		return p.parseThrowStmt()
	case kwOPEN:
		// Check for OPEN SYMMETRIC KEY vs OPEN cursor
		next := p.peekNext()
		if next.Type >= kwADD && matchesKeywordCI(next.Str, "SYMMETRIC") {
			return p.parseOpenSymmetricKeyStmt()
		}
		return p.parseOpenCursorStmt()
	case kwFETCH:
		return p.parseFetchCursorStmt()
	case kwCLOSE:
		// Check for CLOSE SYMMETRIC KEY / CLOSE ALL SYMMETRIC KEYS vs CLOSE cursor
		next := p.peekNext()
		if (next.Type >= kwADD && matchesKeywordCI(next.Str, "SYMMETRIC")) ||
			next.Type == kwALL {
			return p.parseCloseSymmetricKeyStmt()
		}
		return p.parseCloseCursorStmt()
	case kwDEALLOCATE:
		return p.parseDeallocateCursorStmt()
	case kwGO:
		return p.parseGoStmt()
	case kwBULK:
		return p.parseBulkInsertStmt()
	case kwDBCC:
		return p.parseDbccStmt()
	case kwBACKUP:
		// Check for BACKUP CERTIFICATE / BACKUP MASTER KEY vs BACKUP DATABASE/LOG
		next := p.peekNext()
		if (next.Type >= kwADD && matchesKeywordCI(next.Str, "CERTIFICATE")) ||
			(next.Type >= kwADD && matchesKeywordCI(next.Str, "MASTER")) {
			return p.parseBackupCertificateStmt()
		}
		return p.parseBackupStmt()
	case kwRESTORE:
		return p.parseRestoreStmt()
	case kwCHECKPOINT:
		return p.parseCheckpointStmt()
	case kwRECONFIGURE:
		return p.parseReconfigureStmt()
	case kwSHUTDOWN:
		return p.parseShutdownStmt()
	case kwKILL:
		return p.parseKillStmt()
	case kwREADTEXT:
		return p.parseReadtextStmt()
	case kwWRITETEXT:
		return p.parseWritetextStmt()
	case kwUPDATETEXT:
		return p.parseUpdatetextStmt()
	default:
		// Check for label: identifier followed by ':'
		if p.isIdentLike() {
			next := p.peekNext()
			if next.Type == ':' {
				return p.parseLabelStmt()
			}
			// SEND ON CONVERSATION ... (service broker)
			if matchesKeywordCI(p.cur.Str, "SEND") {
				return p.parseSendStmt()
			}
			// RECEIVE (service broker)
			if matchesKeywordCI(p.cur.Str, "RECEIVE") {
				return p.parseReceiveStmt()
			}
			// END CONVERSATION (service broker) - must check next token
			if matchesKeywordCI(p.cur.Str, "END") &&
				next.Str != "" && matchesKeywordCI(next.Str, "CONVERSATION") {
				return p.parseEndConversationStmt()
			}
		}
		return nil
	}
}

// parseCreateStmt dispatches CREATE to the appropriate sub-parser.
func (p *Parser) parseCreateStmt() nodes.StmtNode {
	loc := p.pos()
	p.advance() // consume CREATE

	// OR ALTER
	orAlter := false
	if p.cur.Type == kwOR {
		next := p.peekNext()
		if next.Type == kwALTER {
			p.advance() // OR
			p.advance() // ALTER
			orAlter = true
		}
	}

	// UNIQUE for CREATE UNIQUE INDEX
	unique := false
	if p.cur.Type == kwUNIQUE {
		unique = true
		p.advance()
	}

	switch p.cur.Type {
	case kwTABLE:
		p.advance() // consume TABLE
		stmt := p.parseCreateTableStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwINDEX, kwCLUSTERED, kwNONCLUSTERED, kwCOLUMNSTORE:
		stmt := p.parseCreateIndexStmt(unique)
		stmt.Loc.Start = loc
		return stmt
	case kwVIEW:
		p.advance() // consume VIEW
		stmt := p.parseCreateViewStmt(orAlter)
		stmt.Loc.Start = loc
		return stmt
	case kwPROCEDURE, kwPROC:
		p.advance() // consume PROCEDURE/PROC
		stmt := p.parseCreateProcedureStmt(orAlter)
		stmt.Loc.Start = loc
		return stmt
	case kwFUNCTION:
		p.advance() // consume FUNCTION
		stmt := p.parseCreateFunctionStmt(orAlter)
		stmt.Loc.Start = loc
		return stmt
	case kwDATABASE:
		p.advance() // consume DATABASE
		stmt := p.parseCreateDatabaseStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwTRIGGER:
		p.advance() // consume TRIGGER
		stmt := p.parseCreateTriggerStmt(orAlter)
		stmt.Loc.Start = loc
		return stmt
	case kwSCHEMA:
		p.advance() // consume SCHEMA
		stmt := p.parseCreateSchemaStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwTYPE:
		p.advance() // consume TYPE
		stmt := p.parseCreateTypeStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwUSER:
		p.advance() // consume USER
		stmt := p.parseSecurityUserStmt("CREATE")
		stmt.Loc.Start = loc
		return stmt
	case kwLOGIN:
		p.advance() // consume LOGIN
		stmt := p.parseSecurityLoginStmt("CREATE")
		stmt.Loc.Start = loc
		return stmt
	case kwROLE:
		p.advance() // consume ROLE
		stmt := p.parseSecurityRoleStmt("CREATE")
		stmt.Loc.Start = loc
		return stmt
	default:
		if p.matchIdentCI("SEQUENCE") {
			stmt := p.parseCreateSequenceStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		if p.matchIdentCI("SYNONYM") {
			stmt := p.parseCreateSynonymStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		if p.matchIdentCI("ASSEMBLY") {
			stmt := p.parseCreateAssemblyStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MASTER") {
			stmt := p.parseSecurityKeyStmt("CREATE")
			stmt.Loc.Start = loc
			return stmt
		}
		if p.isIdentLike() && (matchesKeywordCI(p.cur.Str, "SYMMETRIC") ||
			matchesKeywordCI(p.cur.Str, "ASYMMETRIC") ||
			matchesKeywordCI(p.cur.Str, "CERTIFICATE") ||
			matchesKeywordCI(p.cur.Str, "CREDENTIAL")) {
			stmt := p.parseSecurityKeyStmt("CREATE")
			stmt.Loc.Start = loc
			return stmt
		}
		// Check for APPLICATION ROLE (context-sensitive keyword)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "APPLICATION") {
			next := p.peekNext()
			if next.Type == kwROLE || (next.Type >= kwADD && matchesKeywordCI(next.Str, "ROLE")) {
				p.advance() // consume APPLICATION
				stmt := p.parseSecurityApplicationRoleStmt("CREATE")
				stmt.Loc.Start = loc
				return stmt
			}
		}
		// CREATE STATISTICS name ON table (col, ...)
		if p.cur.Type == kwSTATISTICS {
			p.advance() // consume STATISTICS
			stmt := p.parseCreateStatisticsStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE PARTITION FUNCTION / PARTITION SCHEME
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FUNCTION") {
				p.advance() // consume FUNCTION
				stmt := p.parseCreatePartitionFunctionStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SCHEME") {
				p.advance() // consume SCHEME
				stmt := p.parseCreatePartitionSchemeStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE FULLTEXT INDEX / FULLTEXT CATALOG
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FULLTEXT") {
			p.advance() // consume FULLTEXT
			if p.cur.Type == kwINDEX {
				p.advance() // consume INDEX
				stmt := p.parseCreateFulltextIndexStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CATALOG") {
				p.advance() // consume CATALOG
				stmt := p.parseCreateFulltextCatalogStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE XML SCHEMA COLLECTION
		if p.cur.Type == kwXML {
			p.advance() // consume XML
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SCHEMA") {
				p.advance() // consume SCHEMA
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "COLLECTION") {
					p.advance() // consume COLLECTION
					stmt := p.parseCreateXmlSchemaCollectionStmt()
					stmt.Loc.Start = loc
					return stmt
				}
			}
			return nil
		}
		// CREATE MESSAGE TYPE (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MESSAGE") {
			p.advance() // consume MESSAGE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TYPE") {
				p.advance() // consume TYPE
				stmt := p.parseCreateMessageTypeStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE CONTRACT (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONTRACT") {
			p.advance() // consume CONTRACT
			stmt := p.parseCreateContractStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE QUEUE (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "QUEUE") {
			p.advance() // consume QUEUE
			stmt := p.parseCreateQueueStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE SERVICE (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVICE") {
			p.advance() // consume SERVICE
			stmt := p.parseCreateServiceStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		return nil
	}
}

// parseAlterStmt dispatches ALTER to the appropriate sub-parser.
func (p *Parser) parseAlterStmt() nodes.StmtNode {
	loc := p.pos()
	p.advance() // consume ALTER

	switch p.cur.Type {
	case kwTABLE:
		p.advance() // consume TABLE
		stmt := p.parseAlterTableStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwDATABASE:
		p.advance() // consume DATABASE
		stmt := p.parseAlterDatabaseStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwINDEX:
		p.advance() // consume INDEX
		stmt := p.parseAlterIndexStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwVIEW:
		p.advance() // consume VIEW
		stmt := p.parseCreateViewStmt(true /* orAlter */)
		stmt.Loc.Start = loc
		return stmt
	case kwPROCEDURE, kwPROC:
		p.advance() // consume PROCEDURE/PROC
		stmt := p.parseCreateProcedureStmt(true /* orAlter */)
		stmt.Loc.Start = loc
		return stmt
	case kwFUNCTION:
		p.advance() // consume FUNCTION
		stmt := p.parseCreateFunctionStmt(true /* orAlter */)
		stmt.Loc.Start = loc
		return stmt
	case kwSCHEMA:
		p.advance() // consume SCHEMA
		stmt := p.parseAlterSchemaStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwUSER:
		p.advance() // consume USER
		stmt := p.parseSecurityUserStmt("ALTER")
		stmt.Loc.Start = loc
		return stmt
	case kwLOGIN:
		p.advance() // consume LOGIN
		stmt := p.parseSecurityLoginStmt("ALTER")
		stmt.Loc.Start = loc
		return stmt
	case kwROLE:
		p.advance() // consume ROLE
		stmt := p.parseSecurityRoleStmt("ALTER")
		stmt.Loc.Start = loc
		return stmt
	default:
		if p.matchIdentCI("SEQUENCE") {
			stmt := p.parseAlterSequenceStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		if p.matchIdentCI("ASSEMBLY") {
			stmt := p.parseAlterAssemblyStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		if p.isIdentLike() && (matchesKeywordCI(p.cur.Str, "MASTER") ||
			matchesKeywordCI(p.cur.Str, "SYMMETRIC") ||
			matchesKeywordCI(p.cur.Str, "ASYMMETRIC") ||
			matchesKeywordCI(p.cur.Str, "CERTIFICATE") ||
			matchesKeywordCI(p.cur.Str, "CREDENTIAL")) {
			stmt := p.parseSecurityKeyStmt("ALTER")
			stmt.Loc.Start = loc
			return stmt
		}
		// Check for APPLICATION ROLE (context-sensitive keyword)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "APPLICATION") {
			next := p.peekNext()
			if next.Type == kwROLE || (next.Type >= kwADD && matchesKeywordCI(next.Str, "ROLE")) {
				p.advance() // consume APPLICATION
				stmt := p.parseSecurityApplicationRoleStmt("ALTER")
				stmt.Loc.Start = loc
				return stmt
			}
		}
		// ALTER PARTITION FUNCTION / PARTITION SCHEME
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FUNCTION") {
				p.advance() // consume FUNCTION
				stmt := p.parseAlterPartitionFunctionStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SCHEME") {
				p.advance() // consume SCHEME
				stmt := p.parseAlterPartitionSchemeStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// ALTER FULLTEXT INDEX / FULLTEXT CATALOG
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FULLTEXT") {
			p.advance() // consume FULLTEXT
			if p.cur.Type == kwINDEX {
				p.advance() // consume INDEX
				stmt := p.parseAlterFulltextIndexStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CATALOG") {
				p.advance() // consume CATALOG
				stmt := p.parseAlterFulltextCatalogStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// ALTER XML SCHEMA COLLECTION
		if p.cur.Type == kwXML {
			p.advance() // consume XML
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SCHEMA") {
				p.advance() // consume SCHEMA
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "COLLECTION") {
					p.advance() // consume COLLECTION
					stmt := p.parseAlterXmlSchemaCollectionStmt()
					stmt.Loc.Start = loc
					return stmt
				}
			}
			return nil
		}
		return nil
	}
}

// parseDropOrSecurityStmt dispatches DROP to either parseDropStmt (for tables,
// views, etc.) or the security principal parsers (for USER, LOGIN, ROLE,
// APPLICATION ROLE).
func (p *Parser) parseDropOrSecurityStmt() nodes.StmtNode {
	loc := p.pos()
	next := p.peekNext()

	switch next.Type {
	case kwUSER:
		p.advance() // consume DROP
		p.advance() // consume USER
		stmt := p.parseSecurityUserStmt("DROP")
		stmt.Loc.Start = loc
		return stmt
	case kwLOGIN:
		p.advance() // consume DROP
		p.advance() // consume LOGIN
		stmt := p.parseSecurityLoginStmt("DROP")
		stmt.Loc.Start = loc
		return stmt
	case kwROLE:
		p.advance() // consume DROP
		p.advance() // consume ROLE
		stmt := p.parseSecurityRoleStmt("DROP")
		stmt.Loc.Start = loc
		return stmt
	case kwSTATISTICS:
		// DROP STATISTICS table.stats_name [, ...]
		p.advance() // consume DROP
		p.advance() // consume STATISTICS
		stmt := p.parseDropStatisticsStmt()
		stmt.Loc.Start = loc
		return stmt
	default:
		// Check for DROP APPLICATION ROLE
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "APPLICATION") {
			p.advance() // consume DROP
			p.advance() // consume APPLICATION
			stmt := p.parseSecurityApplicationRoleStmt("DROP")
			stmt.Loc.Start = loc
			return stmt
		}
		return p.parseDropStmt()
	}
}

// parseLabelStmt parses a label: statement.
func (p *Parser) parseLabelStmt() *nodes.LabelStmt {
	loc := p.pos()
	label := p.cur.Str
	p.advance() // consume identifier
	p.advance() // consume :
	return &nodes.LabelStmt{
		Label: label,
		Loc:   nodes.Loc{Start: loc, End: p.pos()},
	}
}

// advance consumes the current token and moves to the next one.
func (p *Parser) advance() Token {
	p.prev = p.cur
	if p.hasNext {
		p.cur = p.nextBuf
		p.hasNext = false
	} else {
		p.cur = p.lexer.NextToken()
	}
	return p.prev
}

// peekNext returns the next token after cur without consuming it.
func (p *Parser) peekNext() Token {
	if !p.hasNext {
		p.nextBuf = p.lexer.NextToken()
		p.hasNext = true
	}
	return p.nextBuf
}

// match checks if the current token type matches any of the given types.
// If it matches, the token is consumed and returned with ok=true.
func (p *Parser) match(types ...int) (Token, bool) {
	for _, t := range types {
		if p.cur.Type == t {
			return p.advance(), true
		}
	}
	return Token{}, false
}

// expect consumes the current token if it matches the expected type.
// Returns an error if the token does not match.
func (p *Parser) expect(tokenType int) (Token, error) {
	if p.cur.Type == tokenType {
		return p.advance(), nil
	}
	return Token{}, &ParseError{
		Message:  "unexpected token",
		Position: p.cur.Loc,
	}
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int
}

func (e *ParseError) Error() string {
	return e.Message
}
