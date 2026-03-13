// Package parser implements a recursive descent SQL parser for T-SQL (SQL Server).
//
// This parser reuses the lexer and keyword definitions from this package
// and produces AST nodes from the mssql/ast package.
package parser

import (
	"strings"

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
		// Check for EXECUTE AS (context switching) vs EXEC proc
		if p.cur.Type == kwEXECUTE {
			next := p.peekNext()
			if next.Type == kwAS {
				return p.parseExecuteAsStmt()
			}
		}
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
		// Check for OPEN SYMMETRIC KEY / OPEN MASTER KEY vs OPEN cursor
		next := p.peekNext()
		if next.Str != "" && matchesKeywordCI(next.Str, "SYMMETRIC") {
			return p.parseOpenSymmetricKeyStmt()
		}
		if next.Str != "" && matchesKeywordCI(next.Str, "MASTER") {
			return p.parseOpenMasterKeyStmt()
		}
		return p.parseOpenCursorStmt()
	case kwFETCH:
		return p.parseFetchCursorStmt()
	case kwCLOSE:
		// Check for CLOSE SYMMETRIC KEY / CLOSE ALL SYMMETRIC KEYS / CLOSE MASTER KEY vs CLOSE cursor
		next := p.peekNext()
		if next.Str != "" && matchesKeywordCI(next.Str, "MASTER") {
			return p.parseCloseMasterKeyStmt()
		}
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
		// Check for BACKUP SERVICE MASTER KEY / BACKUP CERTIFICATE / BACKUP MASTER KEY vs BACKUP DATABASE/LOG
		next := p.peekNext()
		if next.Str != "" && matchesKeywordCI(next.Str, "SERVICE") {
			loc := p.pos()
			p.advance() // consume BACKUP
			p.advance() // consume SERVICE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MASTER") {
				p.advance() // consume MASTER
			}
			if p.cur.Type == kwKEY {
				p.advance() // consume KEY
			}
			stmt := p.parseBackupServiceMasterKeyStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		if (next.Str != "" && matchesKeywordCI(next.Str, "CERTIFICATE")) ||
			(next.Str != "" && matchesKeywordCI(next.Str, "MASTER")) ||
			(next.Str != "" && matchesKeywordCI(next.Str, "SYMMETRIC")) {
			return p.parseBackupCertificateStmt()
		}
		return p.parseBackupStmt()
	case kwRESTORE:
		// Check for RESTORE MASTER KEY / RESTORE SYMMETRIC KEY / RESTORE SERVICE MASTER KEY
		next := p.peekNext()
		if next.Str != "" && matchesKeywordCI(next.Str, "MASTER") {
			return p.parseRestoreMasterKeyStmt()
		}
		if next.Str != "" && matchesKeywordCI(next.Str, "SYMMETRIC") {
			return p.parseRestoreSymmetricKeyStmt()
		}
		if next.Str != "" && matchesKeywordCI(next.Str, "SERVICE") {
			loc := p.pos()
			p.advance() // consume RESTORE
			p.advance() // consume SERVICE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MASTER") {
				p.advance() // consume MASTER
			}
			if p.cur.Type == kwKEY {
				p.advance() // consume KEY
			}
			stmt := p.parseRestoreServiceMasterKeyStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		return p.parseRestoreStmt()
	case kwCHECKPOINT:
		return p.parseCheckpointStmt()
	case kwRECONFIGURE:
		return p.parseReconfigureStmt()
	case kwSHUTDOWN:
		return p.parseShutdownStmt()
	case kwKILL:
		return p.parseKillStmt()
	case kwSETUSER:
		return p.parseSetuserStmt()
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
			// GET CONVERSATION GROUP (service broker)
			if matchesKeywordCI(p.cur.Str, "GET") &&
				next.Str != "" && matchesKeywordCI(next.Str, "CONVERSATION") {
				return p.parseGetConversationGroupStmt()
			}
			// ENABLE TRIGGER
			if matchesKeywordCI(p.cur.Str, "ENABLE") &&
				next.Type == kwTRIGGER {
				return p.parseEnableDisableTriggerStmt(true)
			}
			// DISABLE TRIGGER
			if matchesKeywordCI(p.cur.Str, "DISABLE") &&
				next.Type == kwTRIGGER {
				return p.parseEnableDisableTriggerStmt(false)
			}
			// MOVE CONVERSATION (service broker)
			if matchesKeywordCI(p.cur.Str, "MOVE") &&
				next.Str != "" && matchesKeywordCI(next.Str, "CONVERSATION") {
				return p.parseMoveConversationStmt()
			}
			// COPY INTO (Azure Synapse/Fabric bulk load)
			if matchesKeywordCI(p.cur.Str, "COPY") &&
				next.Str != "" && matchesKeywordCI(next.Str, "INTO") {
				return p.parseCopyIntoStmt()
			}
			// RENAME OBJECT/DATABASE (Azure Synapse/PDW)
			if matchesKeywordCI(p.cur.Str, "RENAME") {
				return p.parseRenameStmt()
			}
			// REVERT (security context)
			if matchesKeywordCI(p.cur.Str, "REVERT") {
				return p.parseRevertStmt()
			}
			// PREDICT (ML scoring)
			if matchesKeywordCI(p.cur.Str, "PREDICT") {
				loc := p.pos()
				p.advance() // consume PREDICT
				stmt := p.parsePredictStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			// ADD SENSITIVITY CLASSIFICATION / ADD [COUNTER] SIGNATURE
			if matchesKeywordCI(p.cur.Str, "ADD") {
				return p.parseAddStmt()
			}
		}
		// kwADD is handled here too
		if p.cur.Type == kwADD {
			return p.parseAddStmt()
		}
		return nil
	}
}

// parseAddStmt dispatches ADD SENSITIVITY CLASSIFICATION and ADD [COUNTER] SIGNATURE.
func (p *Parser) parseAddStmt() nodes.StmtNode {
	loc := p.pos()
	p.advance() // consume ADD

	// ADD SENSITIVITY CLASSIFICATION
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SENSITIVITY") {
		p.advance() // consume SENSITIVITY
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLASSIFICATION") {
			p.advance() // consume CLASSIFICATION
		}
		stmt := p.parseAddSensitivityClassificationStmt()
		stmt.Loc.Start = loc
		return stmt
	}

	// ADD [COUNTER] SIGNATURE
	isCounter := false
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "COUNTER") {
		isCounter = true
		p.advance() // consume COUNTER
	}
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SIGNATURE") {
		p.advance() // consume SIGNATURE
		stmt := p.parseSignatureStmt("ADD")
		stmt.IsCounter = isCounter
		stmt.Loc.Start = loc
		return stmt
	}

	return nil
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
	case kwDEFAULT:
		p.advance() // consume DEFAULT
		stmt := p.parseCreateDefaultStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwRULE:
		p.advance() // consume RULE
		stmt := p.parseCreateRuleStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwTABLE:
		p.advance() // consume TABLE
		// Check if this is CREATE TABLE ... AS CLONE OF (Fabric) or CTAS (Synapse)
		// Save lexer state for potential lookahead detection
		savedLexerPos := p.lexer.pos
		savedCur := p.cur
		savedPrev := p.prev
		savedNextBuf := p.nextBuf
		savedHasNext := p.hasNext
		tableName := p.parseTableRef()

		// Check for AS CLONE OF
		if p.cur.Type == kwAS {
			next := p.peekNext()
			if next.Type == tokIDENT && strings.EqualFold(next.Str, "CLONE") {
				p.advance() // consume AS
				cloneStmt := p.parseCreateTableCloneStmt(tableName)
				cloneStmt.Loc.Start = loc
				return cloneStmt
			}
		}

		// Check for CTAS (CREATE TABLE ... AS SELECT / WITH ... AS SELECT)
		// Patterns after table name:
		//   1. WITH (...) AS SELECT ...        (WITH is table options, not CTE)
		//   2. (col_list) WITH (...) AS SELECT ...
		//   3. (col_list) AS SELECT ...
		//   4. AS SELECT ...                   (no WITH, no columns)
		if p.isCTAS(tableName) {
			// Restore and re-parse as CTAS
			p.lexer.pos = savedLexerPos
			p.cur = savedCur
			p.prev = savedPrev
			p.nextBuf = savedNextBuf
			p.hasNext = savedHasNext
			ctasStmt := p.parseCreateTableAsSelectStmt()
			ctasStmt.Loc.Start = loc
			return ctasStmt
		}

		// Not a clone or CTAS — restore lexer state and parse normally
		p.lexer.pos = savedLexerPos
		p.cur = savedCur
		p.prev = savedPrev
		p.nextBuf = savedNextBuf
		p.hasNext = savedHasNext
		stmt := p.parseCreateTableStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwINDEX, kwCLUSTERED, kwNONCLUSTERED, kwCOLUMNSTORE:
		stmt := p.parseCreateIndexStmt(unique)
		stmt.Loc.Start = loc
		return stmt
	case kwPRIMARY:
		// CREATE PRIMARY XML INDEX
		p.advance() // consume PRIMARY
		if p.cur.Type == kwXML {
			p.advance() // consume XML
			p.match(kwINDEX)
			stmt := p.parseCreateXmlIndexStmt(true)
			stmt.Loc.Start = loc
			return stmt
		}
		return nil
	case kwMATERIALIZED:
		// CREATE MATERIALIZED VIEW
		p.advance() // consume MATERIALIZED
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
			stmt := p.parseCreateMaterializedViewStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		return nil
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
	case kwCOLUMN:
		// CREATE COLUMN ENCRYPTION KEY / CREATE COLUMN MASTER KEY
		p.advance() // consume COLUMN
		stmt := p.parseSecurityKeyStmtColumn("CREATE")
		stmt.Loc.Start = loc
		return stmt
	case kwDATABASE:
		// Check for CREATE DATABASE AUDIT SPECIFICATION
		{
			next := p.peekNext()
			if next.Str != "" && matchesKeywordCI(next.Str, "AUDIT") {
				p.advance() // consume DATABASE
				p.advance() // consume AUDIT
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SPECIFICATION") {
					p.advance() // consume SPECIFICATION
				}
				stmt := p.parseCreateDatabaseAuditSpecStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			// CREATE DATABASE ENCRYPTION KEY / CREATE DATABASE SCOPED CREDENTIAL
			if next.Str != "" && (matchesKeywordCI(next.Str, "ENCRYPTION") || matchesKeywordCI(next.Str, "SCOPED")) {
				p.advance() // consume DATABASE
				stmt := p.parseSecurityKeyStmtDatabaseEncryption("CREATE")
				stmt.Loc.Start = loc
				return stmt
			}
		}
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
			matchesKeywordCI(p.cur.Str, "CREDENTIAL") ||
			matchesKeywordCI(p.cur.Str, "CRYPTOGRAPHIC")) {
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
		// CREATE FULLTEXT INDEX / FULLTEXT CATALOG / FULLTEXT STOPLIST
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
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STOPLIST") {
				p.advance() // consume STOPLIST
				stmt := p.parseCreateFulltextStoplistStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE XML SCHEMA COLLECTION / CREATE XML INDEX
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
			if p.cur.Type == kwINDEX {
				p.advance() // consume INDEX
				stmt := p.parseCreateXmlIndexStmt(false)
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE SELECTIVE XML INDEX
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SELECTIVE") {
			p.advance() // consume SELECTIVE
			if p.cur.Type == kwXML {
				p.advance() // consume XML
				p.match(kwINDEX)
				stmt := p.parseCreateSelectiveXmlIndexStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE JSON INDEX
		if p.cur.Type == kwJSON {
			p.advance() // consume JSON
			p.match(kwINDEX)
			stmt := p.parseCreateJsonIndexStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE VECTOR INDEX
		if p.cur.Type == kwVECTOR {
			p.advance() // consume VECTOR
			p.match(kwINDEX)
			stmt := p.parseCreateVectorIndexStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE SPATIAL INDEX
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SPATIAL") {
			p.advance() // consume SPATIAL
			p.match(kwINDEX)
			stmt := p.parseCreateSpatialIndexStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE AGGREGATE
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AGGREGATE") {
			p.advance() // consume AGGREGATE
			stmt := p.parseCreateAggregateStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE SECURITY POLICY
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SECURITY") {
			p.advance() // consume SECURITY
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POLICY") {
				p.advance() // consume POLICY
				stmt := p.parseCreateSecurityPolicyStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE SEARCH PROPERTY LIST
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SEARCH") {
			p.advance() // consume SEARCH
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PROPERTY") {
				p.advance() // consume PROPERTY
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LIST") {
					p.advance() // consume LIST
					stmt := p.parseCreateSearchPropertyListStmt()
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
		// CREATE ROUTE (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ROUTE") {
			p.advance() // consume ROUTE
			stmt := p.parseCreateRouteStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE ENDPOINT
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ENDPOINT") {
			p.advance() // consume ENDPOINT
			stmt := p.parseCreateEndpointStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE REMOTE SERVICE BINDING (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "REMOTE") {
			p.advance() // consume REMOTE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVICE") {
				p.advance() // consume SERVICE
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "BINDING") {
				p.advance() // consume BINDING
			}
			stmt := p.parseCreateRemoteServiceBindingStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE SERVER AUDIT [SPECIFICATION] / CREATE SERVER ROLE
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVER") {
			p.advance() // consume SERVER
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUDIT") {
				p.advance() // consume AUDIT
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SPECIFICATION") {
					p.advance() // consume SPECIFICATION
					stmt := p.parseCreateServerAuditSpecStmt()
					stmt.Loc.Start = loc
					return stmt
				}
				stmt := p.parseCreateServerAuditStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.cur.Type == kwROLE || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ROLE")) {
				p.advance() // consume ROLE
				stmt := p.parseCreateServerRoleStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE EVENT NOTIFICATION / EVENT SESSION
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "EVENT") {
			p.advance() // consume EVENT
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "NOTIFICATION") {
				p.advance() // consume NOTIFICATION
				stmt := p.parseCreateEventNotificationStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SESSION") {
				p.advance() // consume SESSION
				stmt := p.parseCreateEventSessionStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE EXTERNAL DATA SOURCE / EXTERNAL TABLE / EXTERNAL FILE FORMAT / EXTERNAL RESOURCE POOL
		if p.cur.Type == kwEXTERNAL {
			p.advance() // consume EXTERNAL
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DATA") {
				p.advance() // consume DATA
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SOURCE") {
					p.advance() // consume SOURCE
				}
				stmt := p.parseCreateExternalDataSourceStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.cur.Type == kwTABLE {
				p.advance() // consume TABLE
				// Could be CETAS: CREATE EXTERNAL TABLE ... WITH (...) AS SELECT
				// Try CETAS parser which falls back to regular external table if no AS SELECT
				cetasStmt := p.parseCreateExternalTableOrCETAS(loc)
				return cetasStmt
			}
			if p.cur.Type == kwFILE || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FILE")) {
				p.advance() // consume FILE
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FORMAT") {
					p.advance() // consume FORMAT
				}
				stmt := p.parseCreateExternalFileFormatStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RESOURCE") {
				p.advance() // consume RESOURCE
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POOL") {
					p.advance() // consume POOL
				}
				stmt := p.parseCreateExternalResourcePoolStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LIBRARY") {
				p.advance() // consume LIBRARY
				stmt := p.parseCreateExternalLibraryStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LANGUAGE") {
				p.advance() // consume LANGUAGE
				stmt := p.parseCreateExternalLanguageStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MODEL") {
				p.advance() // consume MODEL
				stmt := p.parseCreateExternalModelStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// CREATE AVAILABILITY GROUP
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AVAILABILITY") {
			p.advance() // consume AVAILABILITY
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GROUP") {
				p.advance() // consume GROUP
			}
			stmt := p.parseCreateAvailabilityGroupStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE BROKER PRIORITY (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "BROKER") {
			p.advance() // consume BROKER
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PRIORITY") {
				p.advance() // consume PRIORITY
			}
			stmt := p.parseCreateBrokerPriorityStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE WORKLOAD GROUP / CREATE WORKLOAD CLASSIFIER
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "WORKLOAD") {
			p.advance() // consume WORKLOAD
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLASSIFIER") {
				p.advance() // consume CLASSIFIER
				stmt := p.parseCreateWorkloadClassifierStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GROUP") {
				p.advance() // consume GROUP
			}
			stmt := p.parseCreateWorkloadGroupStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// CREATE RESOURCE POOL
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RESOURCE") {
			p.advance() // consume RESOURCE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POOL") {
				p.advance() // consume POOL
				stmt := p.parseCreateResourcePoolStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
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
	case kwCOLUMN:
		// ALTER COLUMN ENCRYPTION KEY / ALTER COLUMN MASTER KEY
		p.advance() // consume COLUMN
		stmt := p.parseSecurityKeyStmtColumn("ALTER")
		stmt.Loc.Start = loc
		return stmt
	case kwDATABASE:
		// Check for ALTER DATABASE AUDIT SPECIFICATION
		{
			next := p.peekNext()
			if next.Str != "" && matchesKeywordCI(next.Str, "AUDIT") {
				p.advance() // consume DATABASE
				p.advance() // consume AUDIT
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SPECIFICATION") {
					p.advance() // consume SPECIFICATION
				}
				stmt := p.parseAlterDatabaseAuditSpecStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			// ALTER DATABASE ENCRYPTION KEY / ALTER DATABASE SCOPED CREDENTIAL / ALTER DATABASE SCOPED CONFIGURATION
			if next.Str != "" && matchesKeywordCI(next.Str, "SCOPED") {
				// Peek further to distinguish SCOPED CREDENTIAL from SCOPED CONFIGURATION
				p.advance() // consume DATABASE
				p.advance() // consume SCOPED
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONFIGURATION") {
					p.advance() // consume CONFIGURATION
					stmt := p.parseAlterDatabaseScopedConfigStmt()
					stmt.Loc.Start = loc
					return stmt
				}
				// Otherwise it's DATABASE SCOPED CREDENTIAL - reparse via existing handler
				// At this point we've consumed DATABASE SCOPED, and cur is CREDENTIAL or similar
				if p.matchIdentCI("CREDENTIAL") {
					// Manually handle SCOPED CREDENTIAL
					stmtKey := &nodes.SecurityKeyStmt{
						Action:     "ALTER",
						ObjectType: "DATABASE SCOPED CREDENTIAL",
						Loc:        nodes.Loc{Start: loc},
					}
					name, _ := p.parseIdentifier()
					stmtKey.Name = name
					p.parseSecurityKeyOptions(stmtKey)
					stmtKey.Loc.End = p.pos()
					return stmtKey
				}
				return nil
			}
			if next.Str != "" && matchesKeywordCI(next.Str, "ENCRYPTION") {
				p.advance() // consume DATABASE
				stmt := p.parseSecurityKeyStmtDatabaseEncryption("ALTER")
				stmt.Loc.Start = loc
				return stmt
			}
		}
		p.advance() // consume DATABASE
		stmt := p.parseAlterDatabaseStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwINDEX:
		p.advance() // consume INDEX
		stmt := p.parseAlterIndexStmt()
		stmt.Loc.Start = loc
		return stmt
	case kwMATERIALIZED:
		// ALTER MATERIALIZED VIEW
		p.advance() // consume MATERIALIZED
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
			stmt := p.parseAlterMaterializedViewStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		return nil
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
	case kwTRIGGER:
		p.advance() // consume TRIGGER
		stmt := p.parseCreateTriggerStmt(true /* orAlter */)
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
			matchesKeywordCI(p.cur.Str, "CREDENTIAL") ||
			matchesKeywordCI(p.cur.Str, "CRYPTOGRAPHIC")) {
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
		// ALTER FULLTEXT INDEX / FULLTEXT CATALOG / FULLTEXT STOPLIST
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
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STOPLIST") {
				p.advance() // consume STOPLIST
				stmt := p.parseAlterFulltextStoplistStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// ALTER SECURITY POLICY
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SECURITY") {
			p.advance() // consume SECURITY
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POLICY") {
				p.advance() // consume POLICY
				stmt := p.parseAlterSecurityPolicyStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// ALTER SEARCH PROPERTY LIST
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SEARCH") {
			p.advance() // consume SEARCH
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PROPERTY") {
				p.advance() // consume PROPERTY
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LIST") {
					p.advance() // consume LIST
					stmt := p.parseAlterSearchPropertyListStmt()
					stmt.Loc.Start = loc
					return stmt
				}
			}
			return nil
		}
		// ALTER ENDPOINT
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ENDPOINT") {
			p.advance() // consume ENDPOINT
			stmt := p.parseAlterEndpointStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER SERVER AUDIT [SPECIFICATION] / ALTER SERVER ROLE / ALTER SERVER CONFIGURATION
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVER") {
			next := p.peekNext()
			if next.Str != "" && matchesKeywordCI(next.Str, "AUDIT") {
				p.advance() // consume SERVER
				p.advance() // consume AUDIT
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SPECIFICATION") {
					p.advance() // consume SPECIFICATION
					stmt := p.parseAlterServerAuditSpecStmt()
					stmt.Loc.Start = loc
					return stmt
				}
				stmt := p.parseAlterServerAuditStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if next.Type == kwROLE || (next.Str != "" && matchesKeywordCI(next.Str, "ROLE")) {
				p.advance() // consume SERVER
				p.advance() // consume ROLE
				stmt := p.parseAlterServerRoleStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if next.Str != "" && matchesKeywordCI(next.Str, "CONFIGURATION") {
				p.advance() // consume SERVER
				p.advance() // consume CONFIGURATION
				stmt := p.parseAlterServerConfigurationStmt()
				stmt.Loc.Start = loc
				return stmt
			}
		}
		// ALTER DATABASE AUDIT SPECIFICATION (handled via kwDATABASE case)
		// ALTER AUTHORIZATION
		if p.cur.Type == kwAUTHORIZATION {
			p.advance() // consume AUTHORIZATION
			stmt := p.parseAlterAuthorizationStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER QUEUE (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "QUEUE") {
			p.advance() // consume QUEUE
			stmt := p.parseAlterQueueStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER SERVICE (service broker) / ALTER SERVICE MASTER KEY
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVICE") {
			next := p.peekNext()
			if matchesKeywordCI(next.Str, "MASTER") {
				// ALTER SERVICE MASTER KEY
				p.advance() // consume SERVICE
				p.advance() // consume MASTER
				if p.cur.Type == kwKEY {
					p.advance() // consume KEY
				}
				stmt := p.parseAlterServiceMasterKeyStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			p.advance() // consume SERVICE
			stmt := p.parseAlterServiceStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER ROUTE (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ROUTE") {
			p.advance() // consume ROUTE
			stmt := p.parseAlterRouteStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER REMOTE SERVICE BINDING (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "REMOTE") {
			p.advance() // consume REMOTE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVICE") {
				p.advance() // consume SERVICE
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "BINDING") {
				p.advance() // consume BINDING
			}
			stmt := p.parseAlterRemoteServiceBindingStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER BROKER PRIORITY (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "BROKER") {
			p.advance() // consume BROKER
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PRIORITY") {
				p.advance() // consume PRIORITY
			}
			stmt := p.parseAlterBrokerPriorityStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER EVENT SESSION
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "EVENT") {
			p.advance() // consume EVENT
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SESSION") {
				p.advance() // consume SESSION
				stmt := p.parseAlterEventSessionStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// ALTER MESSAGE TYPE (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MESSAGE") {
			p.advance() // consume MESSAGE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TYPE") {
				p.advance() // consume TYPE
			}
			stmt := p.parseAlterMessageTypeStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER CONTRACT (service broker)
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CONTRACT") {
			p.advance() // consume CONTRACT
			stmt := p.parseAlterContractStmt()
			stmt.Loc.Start = loc
			return stmt
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
		// ALTER EXTERNAL DATA SOURCE / ALTER EXTERNAL RESOURCE POOL
		if p.cur.Type == kwEXTERNAL {
			p.advance() // consume EXTERNAL
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "DATA") {
				p.advance() // consume DATA
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SOURCE") {
					p.advance() // consume SOURCE
				}
				stmt := p.parseAlterExternalDataSourceStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RESOURCE") {
				p.advance() // consume RESOURCE
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POOL") {
					p.advance() // consume POOL
				}
				stmt := p.parseAlterExternalResourcePoolStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LIBRARY") {
				p.advance() // consume LIBRARY
				stmt := p.parseAlterExternalLibraryStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LANGUAGE") {
				p.advance() // consume LANGUAGE
				stmt := p.parseAlterExternalLanguageStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "MODEL") {
				p.advance() // consume MODEL
				stmt := p.parseAlterExternalModelStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// ALTER AVAILABILITY GROUP
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AVAILABILITY") {
			p.advance() // consume AVAILABILITY
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GROUP") {
				p.advance() // consume GROUP
			}
			stmt := p.parseAlterAvailabilityGroupStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER WORKLOAD GROUP / ALTER WORKLOAD CLASSIFIER
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "WORKLOAD") {
			p.advance() // consume WORKLOAD
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLASSIFIER") {
				p.advance() // consume CLASSIFIER
				stmt := p.parseAlterWorkloadClassifierStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GROUP") {
				p.advance() // consume GROUP
			}
			stmt := p.parseAlterWorkloadGroupStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// ALTER RESOURCE POOL / ALTER RESOURCE GOVERNOR
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RESOURCE") {
			p.advance() // consume RESOURCE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POOL") {
				p.advance() // consume POOL
				stmt := p.parseAlterResourcePoolStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GOVERNOR") {
				p.advance() // consume GOVERNOR
				stmt := p.parseAlterResourceGovernorStmt()
				stmt.Loc.Start = loc
				return stmt
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

	// Check for DROP DATABASE AUDIT SPECIFICATION — need to match DATABASE + AUDIT
	// We already have `next` pointing to the token after DROP.
	// For DATABASE AUDIT SPECIFICATION, we need to check if next is DATABASE and
	// the token after that is AUDIT. Since we can only peek one ahead, handle this
	// by falling through to parseDropStmt if it's just DROP DATABASE (no AUDIT).

	switch next.Type {
	case kwCOLUMN:
		// DROP COLUMN ENCRYPTION KEY / DROP COLUMN MASTER KEY
		p.advance() // consume DROP
		p.advance() // consume COLUMN
		stmt := p.parseSecurityKeyStmtColumn("DROP")
		stmt.Loc.Start = loc
		return stmt
	case kwDATABASE:
		// Could be DROP DATABASE, DROP DATABASE AUDIT SPECIFICATION,
		// DROP DATABASE ENCRYPTION KEY, or DROP DATABASE SCOPED CREDENTIAL
		p.advance() // consume DROP
		p.advance() // consume DATABASE
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUDIT") {
			p.advance() // consume AUDIT
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SPECIFICATION") {
				p.advance() // consume SPECIFICATION
			}
			stmt := p.parseDropDatabaseAuditSpecStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP DATABASE ENCRYPTION KEY / DROP DATABASE SCOPED CREDENTIAL
		if p.isIdentLike() && (matchesKeywordCI(p.cur.Str, "ENCRYPTION") || matchesKeywordCI(p.cur.Str, "SCOPED")) {
			stmt := p.parseSecurityKeyStmtDatabaseEncryption("DROP")
			stmt.Loc.Start = loc
			return stmt
		}
		// Not AUDIT/ENCRYPTION/SCOPED — it's a regular DROP DATABASE. We already consumed DROP + DATABASE.
		// We need to call parseDropStmt logic but we've already consumed DROP + DATABASE.
		// Let's just handle it inline.
		dropStmt := &nodes.DropStmt{
			ObjectType: nodes.DropDatabase,
			Loc:        nodes.Loc{Start: loc},
		}
		if p.cur.Type == kwIF {
			p.advance()
			p.match(kwEXISTS)
			dropStmt.IfExists = true
		}
		var nameItems []nodes.Node
		for {
			ref := p.parseTableRef()
			if ref != nil {
				nameItems = append(nameItems, ref)
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(nameItems) > 0 {
			dropStmt.Names = &nodes.List{Items: nameItems}
		}
		dropStmt.Loc.End = p.pos()
		return dropStmt
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
		// DROP ENDPOINT
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "ENDPOINT") {
			p.advance() // consume DROP
			p.advance() // consume ENDPOINT
			stmt := p.parseDropEndpointStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP SERVER AUDIT [SPECIFICATION] / DROP SERVER ROLE
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "SERVER") {
			p.advance() // consume DROP
			p.advance() // consume SERVER
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AUDIT") {
				p.advance() // consume AUDIT
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SPECIFICATION") {
					p.advance() // consume SPECIFICATION
					stmt := p.parseDropServerAuditSpecStmt()
					stmt.Loc.Start = loc
					return stmt
				}
				stmt := p.parseDropServerAuditStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.cur.Type == kwROLE || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ROLE")) {
				p.advance() // consume ROLE
				stmt := p.parseDropServerRoleStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// DROP SERVICE BROKER objects
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "MESSAGE") {
			p.advance() // consume DROP
			p.advance() // consume MESSAGE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TYPE") {
				p.advance() // consume TYPE
			}
			stmt := p.parseDropServiceBrokerStmt("MESSAGE TYPE")
			stmt.Loc.Start = loc
			return stmt
		}
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "CONTRACT") {
			p.advance() // consume DROP
			p.advance() // consume CONTRACT
			stmt := p.parseDropServiceBrokerStmt("CONTRACT")
			stmt.Loc.Start = loc
			return stmt
		}
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "QUEUE") {
			p.advance() // consume DROP
			p.advance() // consume QUEUE
			stmt := p.parseDropServiceBrokerStmt("QUEUE")
			stmt.Loc.Start = loc
			return stmt
		}
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "SERVICE") {
			p.advance() // consume DROP
			p.advance() // consume SERVICE
			stmt := p.parseDropServiceBrokerStmt("SERVICE")
			stmt.Loc.Start = loc
			return stmt
		}
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "ROUTE") {
			p.advance() // consume DROP
			p.advance() // consume ROUTE
			stmt := p.parseDropServiceBrokerStmt("ROUTE")
			stmt.Loc.Start = loc
			return stmt
		}
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "REMOTE") {
			p.advance() // consume DROP
			p.advance() // consume REMOTE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERVICE") {
				p.advance() // consume SERVICE
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "BINDING") {
				p.advance() // consume BINDING
			}
			stmt := p.parseDropServiceBrokerStmt("REMOTE SERVICE BINDING")
			stmt.Loc.Start = loc
			return stmt
		}
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "BROKER") {
			p.advance() // consume DROP
			p.advance() // consume BROKER
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PRIORITY") {
				p.advance() // consume PRIORITY
			}
			stmt := p.parseDropServiceBrokerStmt("BROKER PRIORITY")
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP EXTERNAL DATA SOURCE / EXTERNAL TABLE / EXTERNAL FILE FORMAT
		if next.Type == kwEXTERNAL {
			p.advance() // consume DROP
			p.advance() // consume EXTERNAL
			stmt := p.parseDropExternalStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP AVAILABILITY GROUP
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "AVAILABILITY") {
			p.advance() // consume DROP
			p.advance() // consume AVAILABILITY
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GROUP") {
				p.advance() // consume GROUP
			}
			stmt := p.parseDropAvailabilityGroupStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP EVENT NOTIFICATION / EVENT SESSION
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "EVENT") {
			p.advance() // consume DROP
			p.advance() // consume EVENT
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "NOTIFICATION") {
				p.advance() // consume NOTIFICATION
				stmt := p.parseDropEventNotificationStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SESSION") {
				p.advance() // consume SESSION
				stmt := p.parseDropEventSessionStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// DROP CRYPTOGRAPHIC PROVIDER
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "CRYPTOGRAPHIC") {
			p.advance() // consume DROP
			stmt := p.parseSecurityKeyStmt("DROP")
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP WORKLOAD GROUP / DROP WORKLOAD CLASSIFIER
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "WORKLOAD") {
			p.advance() // consume DROP
			p.advance() // consume WORKLOAD
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLASSIFIER") {
				p.advance() // consume CLASSIFIER
				stmt := p.parseDropWorkloadClassifierStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "GROUP") {
				p.advance() // consume GROUP
			}
			stmt := p.parseDropWorkloadGroupStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP RESOURCE POOL
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "RESOURCE") {
			p.advance() // consume DROP
			p.advance() // consume RESOURCE
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POOL") {
				p.advance() // consume POOL
				stmt := p.parseDropResourcePoolStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// DROP SECURITY POLICY
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "SECURITY") {
			p.advance() // consume DROP
			p.advance() // consume SECURITY
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "POLICY") {
				p.advance() // consume POLICY
				stmt := p.parseDropSecurityPolicyStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// DROP SENSITIVITY CLASSIFICATION
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "SENSITIVITY") {
			p.advance() // consume DROP
			p.advance() // consume SENSITIVITY
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CLASSIFICATION") {
				p.advance() // consume CLASSIFICATION
			}
			stmt := p.parseDropSensitivityClassificationStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP [COUNTER] SIGNATURE
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) &&
			(matchesKeywordCI(next.Str, "SIGNATURE") || matchesKeywordCI(next.Str, "COUNTER")) {
			p.advance() // consume DROP
			isCounter := false
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "COUNTER") {
				isCounter = true
				p.advance() // consume COUNTER
			}
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SIGNATURE") {
				p.advance() // consume SIGNATURE
				stmt := p.parseSignatureStmt("DROP")
				stmt.IsCounter = isCounter
				stmt.Loc.Start = loc
				return stmt
			}
			return nil
		}
		// DROP FULLTEXT STOPLIST
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "FULLTEXT") {
			// Peek further: need to distinguish FULLTEXT STOPLIST from FULLTEXT INDEX/CATALOG
			// FULLTEXT INDEX/CATALOG are handled in parseDropStmt, but STOPLIST is a separate stmt type
			p.advance() // consume DROP
			p.advance() // consume FULLTEXT
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "STOPLIST") {
				p.advance() // consume STOPLIST
				stmt := p.parseDropFulltextStoplistStmt()
				stmt.Loc.Start = loc
				return stmt
			}
			// Not STOPLIST — must be INDEX or CATALOG, handled by parseDropStmt.
			// We've already consumed DROP + FULLTEXT, so we need to handle inline.
			// Re-enter the generic drop logic for FULLTEXT INDEX/CATALOG
			dropStmt := &nodes.DropStmt{
				Loc: nodes.Loc{Start: loc},
			}
			if p.cur.Type == kwINDEX {
				dropStmt.ObjectType = nodes.DropFulltextIndex
				p.advance()
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "CATALOG") {
				dropStmt.ObjectType = nodes.DropFulltextCatalog
				p.advance()
			}
			if p.cur.Type == kwIF {
				p.advance()
				p.match(kwEXISTS)
				dropStmt.IfExists = true
			}
			var nameItems []nodes.Node
			if dropStmt.ObjectType == nodes.DropFulltextIndex {
				// DROP FULLTEXT INDEX ON table_name
				if _, ok := p.match(kwON); ok {
					if ref := p.parseTableRef(); ref != nil {
						nameItems = append(nameItems, ref)
					}
				}
			} else {
				for {
					ref := p.parseTableRef()
					if ref != nil {
						nameItems = append(nameItems, ref)
					}
					if _, ok := p.match(','); !ok {
						break
					}
				}
			}
			if len(nameItems) > 0 {
				dropStmt.Names = &nodes.List{Items: nameItems}
			}
			dropStmt.Loc.End = p.pos()
			return dropStmt
		}
		// DROP AGGREGATE
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "AGGREGATE") {
			p.advance() // consume DROP
			p.advance() // consume AGGREGATE
			stmt := p.parseDropAggregateStmt()
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP SEARCH PROPERTY LIST
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "SEARCH") {
			p.advance() // consume DROP
			p.advance() // consume SEARCH
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "PROPERTY") {
				p.advance() // consume PROPERTY
				if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LIST") {
					p.advance() // consume LIST
					stmt := p.parseDropSearchPropertyListStmt()
					stmt.Loc.Start = loc
					return stmt
				}
			}
			return nil
		}
		// DROP MASTER KEY
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "MASTER") {
			p.advance() // consume DROP
			stmt := p.parseSecurityKeyStmt("DROP")
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP SYMMETRIC KEY
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "SYMMETRIC") {
			p.advance() // consume DROP
			stmt := p.parseSecurityKeyStmt("DROP")
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP ASYMMETRIC KEY
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "ASYMMETRIC") {
			p.advance() // consume DROP
			stmt := p.parseSecurityKeyStmt("DROP")
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP CERTIFICATE
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "CERTIFICATE") {
			p.advance() // consume DROP
			stmt := p.parseSecurityKeyStmt("DROP")
			stmt.Loc.Start = loc
			return stmt
		}
		// DROP CREDENTIAL
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) && matchesKeywordCI(next.Str, "CREDENTIAL") {
			p.advance() // consume DROP
			stmt := p.parseSecurityKeyStmt("DROP")
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
