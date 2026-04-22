// Package parser implements a recursive descent SQL parser for T-SQL (SQL Server).
//
// This parser reuses the lexer and keyword definitions from this package
// and produces AST nodes from the mssql/ast package.
package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// Parser is a recursive descent parser for T-SQL.
type Parser struct {
	lexer   *Lexer
	source  string // original SQL input (for error context extraction)
	cur     Token  // current token
	prev    Token  // previous token (for error reporting)
	nextBuf Token  // buffered next token for 2-token lookahead
	hasNext bool   // whether nextBuf is valid

	// Completion mode fields.
	completing bool          // true when running in completion mode
	cursorOff  int           // byte offset of the cursor
	candidates *CandidateSet // collected candidates
	collecting bool          // true once cursor position has been reached
	maxCollect int           // max candidates before stopping
}

// Parse parses a T-SQL string into an AST list.
// Currently supports basic infrastructure; statement dispatch will be
// implemented incrementally across batches.
func Parse(sql string) (*nodes.List, error) {
	p := &Parser{
		lexer:  NewLexer(sql),
		source: sql,
	}
	p.advance()

	var stmts []nodes.Node
	for p.cur.Type != tokEOF {
		// Skip semicolons
		if p.cur.Type == ';' {
			p.advance()
			continue
		}
		// Check for lexer errors before parsing statement.
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		// Check for lexer errors after parsing statement.
		if p.lexer.Err != nil {
			return nil, p.lexerError()
		}
		if stmt == nil {
			if p.cur.Type != tokEOF {
				return nil, p.syntaxErrorAtCur()
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
func (p *Parser) parseStmt() (nodes.StmtNode, error) {
	// Completion: when in collect mode, offer all top-level statement keywords.
	if p.collectMode() {
		for _, t := range topLevelKeywords {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}
	p.checkCursor()

	switch p.cur.Type {
	case kwSELECT:
		return p.parseSelectStmt()
	case kwWITH:
		return p.parseSelectStmt()
	case kwINSERT:
		// Check for INSERT BULK
		{
			next := p.peekNext()
			if next.Type == kwBULK {
				loc := p.pos()
				p.advance() // consume INSERT
				p.advance() // consume BULK
				stmt, err := p.parseInsertBulkStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
		}
		return p.parseInsertStmt()
	case kwUPDATE:
		// Could be UPDATE table... or UPDATE STATISTICS
		next := p.peekNext()
		if next.Type == kwSTATISTICS {
			loc := p.pos()
			p.advance() // consume UPDATE
			p.advance() // consume STATISTICS
			stmt, err := p.parseUpdateStatisticsStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
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
		// Check for USE FEDERATION
		{
			next := p.peekNext()
			if next.Type == kwFEDERATION {
				loc := p.pos()
				p.advance() // consume USE
				p.advance() // consume FEDERATION
				stmt, err := p.parseUseFederationStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
		}
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
		if next.Type == kwSYMMETRIC {
			return p.parseOpenSymmetricKeyStmt()
		}
		if next.Type == kwMASTER {
			return p.parseOpenMasterKeyStmt()
		}
		return p.parseOpenCursorStmt()
	case kwFETCH:
		return p.parseFetchCursorStmt()
	case kwCLOSE:
		// Check for CLOSE SYMMETRIC KEY / CLOSE ALL SYMMETRIC KEYS / CLOSE MASTER KEY vs CLOSE cursor
		next := p.peekNext()
		if next.Type == kwMASTER {
			return p.parseCloseMasterKeyStmt()
		}
		if (next.Type == kwSYMMETRIC) ||
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
	case kwLINENO:
		return p.parseLinenoStmt()
	case kwBACKUP:
		// Check for BACKUP SERVICE MASTER KEY / BACKUP CERTIFICATE / BACKUP MASTER KEY vs BACKUP DATABASE/LOG
		next := p.peekNext()
		if next.Type == kwSERVICE {
			loc := p.pos()
			p.advance() // consume BACKUP
			p.advance() // consume SERVICE
			if p.cur.Type == kwMASTER {
				p.advance() // consume MASTER
			}
			if p.cur.Type == kwKEY {
				p.advance() // consume KEY
			}
			stmt, err := p.parseBackupServiceMasterKeyStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if (next.Type == kwCERTIFICATE) ||
			(next.Type == kwMASTER) ||
			(next.Type == kwSYMMETRIC) {
			return p.parseBackupCertificateStmt()
		}
		return p.parseBackupStmt()
	case kwRESTORE:
		// Check for RESTORE MASTER KEY / RESTORE SYMMETRIC KEY / RESTORE SERVICE MASTER KEY
		next := p.peekNext()
		if next.Type == kwMASTER {
			return p.parseRestoreMasterKeyStmt()
		}
		if next.Type == kwSYMMETRIC {
			return p.parseRestoreSymmetricKeyStmt()
		}
		if next.Type == kwSERVICE {
			loc := p.pos()
			p.advance() // consume RESTORE
			p.advance() // consume SERVICE
			if p.cur.Type == kwMASTER {
				p.advance() // consume MASTER
			}
			if p.cur.Type == kwKEY {
				p.advance() // consume KEY
			}
			stmt, err := p.parseRestoreServiceMasterKeyStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
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
	case kwEND:
		// END CONVERSATION (service broker)
		next := p.peekNext()
		if next.Type == kwCONVERSATION {
			return p.parseEndConversationStmt()
		}
		// END without matching BEGIN — return nil to let caller handle
		return nil, p.lexerError()
	case kwREVERT:
		return p.parseRevertStmt()
	default:
		// Check for label: identifier followed by ':'
		if p.isAnyKeywordIdent() {
			next := p.peekNext()
			if next.Type == ':' {
				// Labels must be identifiers; reject reserved keywords here.
				if !p.isIdentLike() {
					return nil, p.syntaxErrorAtCur()
				}
				return p.parseLabelStmt()
			}
			// SEND ON CONVERSATION ... (service broker)
			if p.cur.Type == kwSEND {
				return p.parseSendStmt()
			}
			// RECEIVE (service broker)
			if p.cur.Type == kwRECEIVE {
				return p.parseReceiveStmt()
			}
			// END CONVERSATION (service broker) - must check next token
			if p.cur.Type == kwEND &&
				next.Type == kwCONVERSATION {
				return p.parseEndConversationStmt()
			}
			// GET CONVERSATION GROUP (service broker)
			if p.cur.Type == kwGET &&
				next.Type == kwCONVERSATION {
				return p.parseGetConversationGroupStmt()
			}
			// ENABLE TRIGGER
			if p.cur.Type == kwENABLE &&
				next.Type == kwTRIGGER {
				return p.parseEnableDisableTriggerStmt(true)
			}
			// DISABLE TRIGGER
			if p.cur.Type == kwDISABLE &&
				next.Type == kwTRIGGER {
				return p.parseEnableDisableTriggerStmt(false)
			}
			// MOVE CONVERSATION (service broker)
			if p.cur.Type == kwMOVE &&
				next.Type == kwCONVERSATION {
				return p.parseMoveConversationStmt()
			}
			// COPY INTO (Azure Synapse/Fabric bulk load)
			if p.cur.Type == kwCOPY &&
				next.Type == kwINTO {
				return p.parseCopyIntoStmt()
			}
			// RENAME OBJECT/DATABASE (Azure Synapse/PDW)
			if p.cur.Type == kwRENAME {
				return p.parseRenameStmt()
			}
			// REVERT (security context)
			if p.cur.Type == kwREVERT {
				return p.parseRevertStmt()
			}
			// PREDICT (ML scoring)
			if p.cur.Type == kwPREDICT {
				loc := p.pos()
				p.advance() // consume PREDICT
				stmt, err := p.parsePredictStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			// ADD SENSITIVITY CLASSIFICATION / ADD [COUNTER] SIGNATURE
			if p.cur.Type == kwADD {
				return p.parseAddStmt()
			}
		}
		// kwADD is handled here too
		if p.cur.Type == kwADD {
			return p.parseAddStmt()
		}
		return nil, nil
	}
}

// parseAddStmt dispatches ADD SENSITIVITY CLASSIFICATION and ADD [COUNTER] SIGNATURE.
func (p *Parser) parseAddStmt() (nodes.StmtNode, error) {
	loc := p.pos()
	p.advance() // consume ADD

	// ADD SENSITIVITY CLASSIFICATION
	if p.cur.Type == kwSENSITIVITY {
		p.advance() // consume SENSITIVITY
		if p.cur.Type == kwCLASSIFICATION {
			p.advance() // consume CLASSIFICATION
		}
		stmt, err := p.parseAddSensitivityClassificationStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	}

	// ADD [COUNTER] SIGNATURE
	isCounter := false
	if p.cur.Type == kwCOUNTER {
		isCounter = true
		p.advance() // consume COUNTER
	}
	if p.cur.Type == kwSIGNATURE {
		p.advance() // consume SIGNATURE
		stmt, err := p.parseSignatureStmt("ADD")
		if err != nil {
			return nil, err
		}
		stmt.IsCounter = isCounter
		stmt.Loc.Start = loc
		return stmt, nil
	}

	return nil, nil
}

// parseCreateStmt dispatches CREATE to the appropriate sub-parser.
func (p *Parser) parseCreateStmt() (nodes.StmtNode, error) {
	loc := p.pos()
	p.advance() // consume CREATE

	// Completion: offer object type keywords after CREATE.
	if p.collectMode() {
		for _, kw := range []int{
			kwTABLE, kwINDEX, kwVIEW, kwPROCEDURE, kwFUNCTION, kwTRIGGER,
			kwDATABASE, kwSCHEMA, kwTYPE, kwUSER, kwLOGIN, kwROLE, kwSTATISTICS,
		} {
			p.addTokenCandidate(kw)
		}
		// Context-sensitive identifiers: SEQUENCE, SYNONYM
		p.addRuleCandidate("SEQUENCE")
		p.addRuleCandidate("SYNONYM")
		return nil, errCollecting
	}

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
		stmt, err := p.parseCreateDefaultStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwRULE:
		p.advance() // consume RULE
		stmt, err := p.parseCreateRuleStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwTABLE:
		p.advance() // consume TABLE
		// Check if this is CREATE TABLE ... AS CLONE OF (Fabric) or CTAS (Synapse)
		// Save lexer state for potential lookahead detection
		savedLexerPos := p.lexer.pos
		savedCur := p.cur
		savedPrev := p.prev
		savedNextBuf := p.nextBuf
		savedHasNext := p.hasNext
		savedCollecting := p.collecting
		tableName, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}

		// Check for AS CLONE OF
		if p.cur.Type == kwAS {
			next := p.peekNext()
			if next.Type == kwCLONE {
				p.advance() // consume AS
				cloneStmt, err := p.parseCreateTableCloneStmt(tableName)
				if err != nil {
					return nil, err
				}
				cloneStmt.Loc.Start = loc
				return cloneStmt, nil
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
			p.collecting = savedCollecting
			ctasStmt, err := p.parseCreateTableAsSelectStmt()
			if err != nil {
				return nil, err
			}
			ctasStmt.Loc.Start = loc
			return ctasStmt, nil
		}

		// Not a clone or CTAS — restore lexer state and parse normally
		p.lexer.pos = savedLexerPos
		p.cur = savedCur
		p.prev = savedPrev
		p.nextBuf = savedNextBuf
		p.hasNext = savedHasNext
		p.collecting = savedCollecting
		stmt, err := p.parseCreateTableStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwINDEX, kwCLUSTERED, kwNONCLUSTERED, kwCOLUMNSTORE:
		stmt, err := p.parseCreateIndexStmt(unique)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwPRIMARY:
		// CREATE PRIMARY XML INDEX
		p.advance() // consume PRIMARY
		if p.cur.Type == kwXML {
			p.advance() // consume XML
			p.match(kwINDEX)
			stmt, err := p.parseCreateXmlIndexStmt(true)
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		return nil, nil
	case kwMATERIALIZED:
		// CREATE MATERIALIZED VIEW
		p.advance() // consume MATERIALIZED
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
			stmt, err := p.parseCreateMaterializedViewStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		return nil, nil
	case kwVIEW:
		p.advance() // consume VIEW
		stmt, err := p.parseCreateViewStmt(orAlter)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwPROCEDURE, kwPROC:
		p.advance() // consume PROCEDURE/PROC
		stmt, err := p.parseCreateProcedureStmt(orAlter)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwFUNCTION:
		p.advance() // consume FUNCTION
		stmt, err := p.parseCreateFunctionStmt(orAlter)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwCOLUMN:
		// CREATE COLUMN ENCRYPTION KEY / CREATE COLUMN MASTER KEY
		p.advance() // consume COLUMN
		stmt, err := p.parseSecurityKeyStmtColumn("CREATE")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwDATABASE:
		// Check for CREATE DATABASE AUDIT SPECIFICATION
		{
			next := p.peekNext()
			if next.Type == kwAUDIT {
				p.advance() // consume DATABASE
				p.advance() // consume AUDIT
				if p.cur.Type == kwSPECIFICATION {
					p.advance() // consume SPECIFICATION
				}
				stmt, err := p.parseCreateDatabaseAuditSpecStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			// CREATE DATABASE ENCRYPTION KEY / CREATE DATABASE SCOPED CREDENTIAL
			if next.Str != "" && (next.Type == kwENCRYPTION || next.Type == kwSCOPED) {
				p.advance() // consume DATABASE
				stmt, err := p.parseSecurityKeyStmtDatabaseEncryption("CREATE")
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
		}
		p.advance() // consume DATABASE
		stmt, err := p.parseCreateDatabaseStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwTRIGGER:
		p.advance() // consume TRIGGER
		stmt, err := p.parseCreateTriggerStmt(orAlter)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwSCHEMA:
		p.advance() // consume SCHEMA
		stmt, err := p.parseCreateSchemaStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwTYPE:
		p.advance() // consume TYPE
		stmt, err := p.parseCreateTypeStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwUSER:
		p.advance() // consume USER
		stmt, err := p.parseSecurityUserStmt("CREATE")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwLOGIN:
		p.advance() // consume LOGIN
		stmt, err := p.parseSecurityLoginStmt("CREATE")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwROLE:
		p.advance() // consume ROLE
		stmt, err := p.parseSecurityRoleStmt("CREATE")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	default:
		if p.matchIdentCI("SEQUENCE") {
			stmt, err := p.parseCreateSequenceStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if p.matchIdentCI("SYNONYM") {
			stmt, err := p.parseCreateSynonymStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if p.matchIdentCI("ASSEMBLY") {
			stmt, err := p.parseCreateAssemblyStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if p.cur.Type == kwMASTER {
			stmt, err := p.parseSecurityKeyStmt("CREATE")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if p.isAnyKeywordIdent() && (p.cur.Type == kwSYMMETRIC ||
			p.cur.Type == kwASYMMETRIC ||
			p.cur.Type == kwCERTIFICATE ||
			p.cur.Type == kwCREDENTIAL ||
			p.cur.Type == kwCRYPTOGRAPHIC) {
			stmt, err := p.parseSecurityKeyStmt("CREATE")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// Check for APPLICATION ROLE (context-sensitive keyword)
		if p.cur.Type == kwAPPLICATION {
			next := p.peekNext()
			if next.Type == kwROLE {
				p.advance() // consume APPLICATION
				stmt, err := p.parseSecurityApplicationRoleStmt("CREATE")
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
		}
		// CREATE STATISTICS name ON table (col, ...)
		if p.cur.Type == kwSTATISTICS {
			p.advance() // consume STATISTICS
			stmt, err := p.parseCreateStatisticsStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE PARTITION FUNCTION / PARTITION SCHEME
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
			if p.cur.Type == kwFUNCTION {
				p.advance() // consume FUNCTION
				stmt, err := p.parseCreatePartitionFunctionStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSCHEME {
				p.advance() // consume SCHEME
				stmt, err := p.parseCreatePartitionSchemeStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE FULLTEXT INDEX / FULLTEXT CATALOG / FULLTEXT STOPLIST
		if p.cur.Type == kwFULLTEXT {
			p.advance() // consume FULLTEXT
			if p.cur.Type == kwINDEX {
				p.advance() // consume INDEX
				stmt, err := p.parseCreateFulltextIndexStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwCATALOG {
				p.advance() // consume CATALOG
				stmt, err := p.parseCreateFulltextCatalogStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSTOPLIST {
				p.advance() // consume STOPLIST
				stmt, err := p.parseCreateFulltextStoplistStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE XML SCHEMA COLLECTION / CREATE XML INDEX
		if p.cur.Type == kwXML {
			p.advance() // consume XML
			if p.cur.Type == kwSCHEMA {
				p.advance() // consume SCHEMA
				if p.cur.Type == kwCOLLECTION {
					p.advance() // consume COLLECTION
					stmt, err := p.parseCreateXmlSchemaCollectionStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
			}
			if p.cur.Type == kwINDEX {
				p.advance() // consume INDEX
				stmt, err := p.parseCreateXmlIndexStmt(false)
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE SELECTIVE XML INDEX
		if p.cur.Type == kwSELECTIVE {
			p.advance() // consume SELECTIVE
			if p.cur.Type == kwXML {
				p.advance() // consume XML
				p.match(kwINDEX)
				stmt, err := p.parseCreateSelectiveXmlIndexStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE JSON INDEX
		if p.cur.Type == kwJSON {
			p.advance() // consume JSON
			p.match(kwINDEX)
			stmt, err := p.parseCreateJsonIndexStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE VECTOR INDEX
		if p.cur.Type == kwVECTOR {
			p.advance() // consume VECTOR
			p.match(kwINDEX)
			stmt, err := p.parseCreateVectorIndexStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE SPATIAL INDEX
		if p.cur.Type == kwSPATIAL {
			p.advance() // consume SPATIAL
			p.match(kwINDEX)
			stmt, err := p.parseCreateSpatialIndexStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE AGGREGATE
		if p.cur.Type == kwAGGREGATE {
			p.advance() // consume AGGREGATE
			stmt, err := p.parseCreateAggregateStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE SECURITY POLICY
		if p.cur.Type == kwSECURITY {
			p.advance() // consume SECURITY
			if p.cur.Type == kwPOLICY {
				p.advance() // consume POLICY
				stmt, err := p.parseCreateSecurityPolicyStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE SEARCH PROPERTY LIST
		if p.cur.Type == kwSEARCH {
			p.advance() // consume SEARCH
			if p.cur.Type == kwPROPERTY {
				p.advance() // consume PROPERTY
				if p.cur.Type == kwLIST {
					p.advance() // consume LIST
					stmt, err := p.parseCreateSearchPropertyListStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
			}
			return nil, nil
		}
		// CREATE MESSAGE TYPE (service broker)
		if p.cur.Type == kwMESSAGE {
			p.advance() // consume MESSAGE
			if p.cur.Type == kwTYPE {
				p.advance() // consume TYPE
				stmt, err := p.parseCreateMessageTypeStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE CONTRACT (service broker)
		if p.cur.Type == kwCONTRACT {
			p.advance() // consume CONTRACT
			stmt, err := p.parseCreateContractStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE QUEUE (service broker)
		if p.cur.Type == kwQUEUE {
			p.advance() // consume QUEUE
			stmt, err := p.parseCreateQueueStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE SERVICE (service broker)
		if p.cur.Type == kwSERVICE {
			p.advance() // consume SERVICE
			stmt, err := p.parseCreateServiceStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE ROUTE (service broker)
		if p.cur.Type == kwROUTE {
			p.advance() // consume ROUTE
			stmt, err := p.parseCreateRouteStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE ENDPOINT
		if p.cur.Type == kwENDPOINT {
			p.advance() // consume ENDPOINT
			stmt, err := p.parseCreateEndpointStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE REMOTE TABLE AS SELECT (CRTAS) / CREATE REMOTE SERVICE BINDING (service broker)
		if p.cur.Type == kwREMOTE {
			p.advance() // consume REMOTE
			if p.cur.Type == kwTABLE {
				// CREATE REMOTE TABLE ... AT (...) AS SELECT
				p.advance() // consume TABLE
				stmt, err := p.parseCreateRemoteTableAsSelectStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSERVICE {
				p.advance() // consume SERVICE
			}
			if p.cur.Type == kwBINDING {
				p.advance() // consume BINDING
			}
			stmt, err := p.parseCreateRemoteServiceBindingStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE SERVER AUDIT [SPECIFICATION] / CREATE SERVER ROLE
		if p.cur.Type == kwSERVER {
			p.advance() // consume SERVER
			if p.cur.Type == kwAUDIT {
				p.advance() // consume AUDIT
				if p.cur.Type == kwSPECIFICATION {
					p.advance() // consume SPECIFICATION
					stmt, err := p.parseCreateServerAuditSpecStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
				stmt, err := p.parseCreateServerAuditStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwROLE {
				p.advance() // consume ROLE
				stmt, err := p.parseCreateServerRoleStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE EVENT NOTIFICATION / EVENT SESSION
		if p.cur.Type == kwEVENT {
			p.advance() // consume EVENT
			if p.cur.Type == kwNOTIFICATION {
				p.advance() // consume NOTIFICATION
				stmt, err := p.parseCreateEventNotificationStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSESSION {
				p.advance() // consume SESSION
				stmt, err := p.parseCreateEventSessionStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE EXTERNAL DATA SOURCE / EXTERNAL TABLE / EXTERNAL FILE FORMAT / EXTERNAL RESOURCE POOL
		if p.cur.Type == kwEXTERNAL {
			p.advance() // consume EXTERNAL
			if p.cur.Type == kwDATA {
				p.advance() // consume DATA
				if p.cur.Type == kwSOURCE {
					p.advance() // consume SOURCE
				}
				stmt, err := p.parseCreateExternalDataSourceStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwTABLE {
				p.advance() // consume TABLE
				// Could be CETAS: CREATE EXTERNAL TABLE ... WITH (...) AS SELECT
				// Try CETAS parser which falls back to regular external table if no AS SELECT
				cetasStmt, err := p.parseCreateExternalTableOrCETAS(loc)
				if err != nil {
					return nil, err
				}
				return cetasStmt, nil
			}
			if p.cur.Type == kwFILE {
				p.advance() // consume FILE
				if p.cur.Type == kwFORMAT {
					p.advance() // consume FORMAT
				}
				stmt, err := p.parseCreateExternalFileFormatStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwRESOURCE {
				p.advance() // consume RESOURCE
				if p.cur.Type == kwPOOL {
					p.advance() // consume POOL
				}
				stmt, err := p.parseCreateExternalResourcePoolStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwLIBRARY {
				p.advance() // consume LIBRARY
				stmt, err := p.parseCreateExternalLibraryStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwLANGUAGE {
				p.advance() // consume LANGUAGE
				stmt, err := p.parseCreateExternalLanguageStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwMODEL {
				p.advance() // consume MODEL
				stmt, err := p.parseCreateExternalModelStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSTREAM {
				p.advance() // consume STREAM
				stmt, err := p.parseCreateExternalStreamStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSTREAMING {
				p.advance() // consume STREAMING
				if p.cur.Type == kwJOB {
					p.advance() // consume JOB
				}
				stmt, err := p.parseCreateExternalStreamingJobStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE AVAILABILITY GROUP
		if p.cur.Type == kwAVAILABILITY {
			p.advance() // consume AVAILABILITY
			if p.cur.Type == kwGROUP {
				p.advance() // consume GROUP
			}
			stmt, err := p.parseCreateAvailabilityGroupStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE BROKER PRIORITY (service broker)
		if p.cur.Type == kwBROKER {
			p.advance() // consume BROKER
			if p.cur.Type == kwPRIORITY {
				p.advance() // consume PRIORITY
			}
			stmt, err := p.parseCreateBrokerPriorityStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE WORKLOAD GROUP / CREATE WORKLOAD CLASSIFIER
		if p.cur.Type == kwWORKLOAD {
			p.advance() // consume WORKLOAD
			if p.cur.Type == kwCLASSIFIER {
				p.advance() // consume CLASSIFIER
				stmt, err := p.parseCreateWorkloadClassifierStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwGROUP {
				p.advance() // consume GROUP
			}
			stmt, err := p.parseCreateWorkloadGroupStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// CREATE RESOURCE POOL
		if p.cur.Type == kwRESOURCE {
			p.advance() // consume RESOURCE
			if p.cur.Type == kwPOOL {
				p.advance() // consume POOL
				stmt, err := p.parseCreateResourcePoolStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// CREATE FEDERATION
		if p.cur.Type == kwFEDERATION {
			p.advance() // consume FEDERATION
			stmt, err := p.parseCreateFederationStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		return nil, nil
	}
}

// parseAlterStmt dispatches ALTER to the appropriate sub-parser.
func (p *Parser) parseAlterStmt() (nodes.StmtNode, error) {
	loc := p.pos()
	p.advance() // consume ALTER

	// Completion: offer object type keywords after ALTER.
	if p.collectMode() {
		for _, kw := range []int{
			kwTABLE, kwINDEX, kwVIEW, kwPROCEDURE, kwFUNCTION, kwTRIGGER,
			kwDATABASE, kwSCHEMA, kwUSER, kwLOGIN, kwROLE,
		} {
			p.addTokenCandidate(kw)
		}
		// Context-sensitive identifier: SEQUENCE
		p.addRuleCandidate("SEQUENCE")
		return nil, errCollecting
	}

	switch p.cur.Type {
	case kwTABLE:
		p.advance() // consume TABLE
		stmt, err := p.parseAlterTableStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwCOLUMN:
		// ALTER COLUMN ENCRYPTION KEY / ALTER COLUMN MASTER KEY
		p.advance() // consume COLUMN
		stmt, err := p.parseSecurityKeyStmtColumn("ALTER")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwDATABASE:
		// Check for ALTER DATABASE AUDIT SPECIFICATION
		{
			next := p.peekNext()
			if next.Type == kwAUDIT {
				p.advance() // consume DATABASE
				p.advance() // consume AUDIT
				if p.cur.Type == kwSPECIFICATION {
					p.advance() // consume SPECIFICATION
				}
				stmt, err := p.parseAlterDatabaseAuditSpecStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			// ALTER DATABASE ENCRYPTION KEY / ALTER DATABASE SCOPED CREDENTIAL / ALTER DATABASE SCOPED CONFIGURATION
			if next.Type == kwSCOPED {
				// Peek further to distinguish SCOPED CREDENTIAL from SCOPED CONFIGURATION
				p.advance() // consume DATABASE
				p.advance() // consume SCOPED
				if p.cur.Type == kwCONFIGURATION {
					p.advance() // consume CONFIGURATION
					stmt, err := p.parseAlterDatabaseScopedConfigStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
				// Otherwise it's DATABASE SCOPED CREDENTIAL - reparse via existing handler
				// At this point we've consumed DATABASE SCOPED, and cur is CREDENTIAL or similar
				if p.matchIdentCI("CREDENTIAL") {
					// Manually handle SCOPED CREDENTIAL
					stmtKey := &nodes.SecurityKeyStmt{
						Action:     "ALTER",
						ObjectType: "DATABASE SCOPED CREDENTIAL",
						Loc:        nodes.Loc{Start: loc, End: -1},
					}
					name, _ := p.parseIdentifier()
					stmtKey.Name = name
					p.parseSecurityKeyOptions(stmtKey)
					stmtKey.Loc.End = p.prevEnd()
					return stmtKey, nil
				}
				return nil, nil
			}
			if next.Type == kwENCRYPTION {
				p.advance() // consume DATABASE
				stmt, err := p.parseSecurityKeyStmtDatabaseEncryption("ALTER")
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
		}
		p.advance() // consume DATABASE
		stmt, err := p.parseAlterDatabaseStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwINDEX:
		p.advance() // consume INDEX
		stmt, err := p.parseAlterIndexStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwMATERIALIZED:
		// ALTER MATERIALIZED VIEW
		p.advance() // consume MATERIALIZED
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
			stmt, err := p.parseAlterMaterializedViewStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		return nil, nil
	case kwVIEW:
		p.advance() // consume VIEW
		stmt, err := p.parseCreateViewStmt(true /* orAlter */)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwPROCEDURE, kwPROC:
		p.advance() // consume PROCEDURE/PROC
		stmt, err := p.parseCreateProcedureStmt(true /* orAlter */)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwFUNCTION:
		p.advance() // consume FUNCTION
		stmt, err := p.parseCreateFunctionStmt(true /* orAlter */)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwTRIGGER:
		p.advance() // consume TRIGGER
		stmt, err := p.parseCreateTriggerStmt(true /* orAlter */)
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwSCHEMA:
		p.advance() // consume SCHEMA
		stmt, err := p.parseAlterSchemaStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwUSER:
		p.advance() // consume USER
		stmt, err := p.parseSecurityUserStmt("ALTER")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwLOGIN:
		p.advance() // consume LOGIN
		stmt, err := p.parseSecurityLoginStmt("ALTER")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwROLE:
		p.advance() // consume ROLE
		stmt, err := p.parseSecurityRoleStmt("ALTER")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	default:
		if p.matchIdentCI("SEQUENCE") {
			stmt, err := p.parseAlterSequenceStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if p.matchIdentCI("ASSEMBLY") {
			stmt, err := p.parseAlterAssemblyStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if p.isAnyKeywordIdent() && (p.cur.Type == kwMASTER ||
			p.cur.Type == kwSYMMETRIC ||
			p.cur.Type == kwASYMMETRIC ||
			p.cur.Type == kwCERTIFICATE ||
			p.cur.Type == kwCREDENTIAL ||
			p.cur.Type == kwCRYPTOGRAPHIC) {
			stmt, err := p.parseSecurityKeyStmt("ALTER")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// Check for APPLICATION ROLE (context-sensitive keyword)
		if p.cur.Type == kwAPPLICATION {
			next := p.peekNext()
			if next.Type == kwROLE {
				p.advance() // consume APPLICATION
				stmt, err := p.parseSecurityApplicationRoleStmt("ALTER")
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
		}
		// ALTER PARTITION FUNCTION / PARTITION SCHEME
		if p.cur.Type == kwPARTITION {
			p.advance() // consume PARTITION
			if p.cur.Type == kwFUNCTION {
				p.advance() // consume FUNCTION
				stmt, err := p.parseAlterPartitionFunctionStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSCHEME {
				p.advance() // consume SCHEME
				stmt, err := p.parseAlterPartitionSchemeStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// ALTER FULLTEXT INDEX / FULLTEXT CATALOG / FULLTEXT STOPLIST
		if p.cur.Type == kwFULLTEXT {
			p.advance() // consume FULLTEXT
			if p.cur.Type == kwINDEX {
				p.advance() // consume INDEX
				stmt, err := p.parseAlterFulltextIndexStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwCATALOG {
				p.advance() // consume CATALOG
				stmt, err := p.parseAlterFulltextCatalogStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSTOPLIST {
				p.advance() // consume STOPLIST
				stmt, err := p.parseAlterFulltextStoplistStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// ALTER SECURITY POLICY
		if p.cur.Type == kwSECURITY {
			p.advance() // consume SECURITY
			if p.cur.Type == kwPOLICY {
				p.advance() // consume POLICY
				stmt, err := p.parseAlterSecurityPolicyStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// ALTER SEARCH PROPERTY LIST
		if p.cur.Type == kwSEARCH {
			p.advance() // consume SEARCH
			if p.cur.Type == kwPROPERTY {
				p.advance() // consume PROPERTY
				if p.cur.Type == kwLIST {
					p.advance() // consume LIST
					stmt, err := p.parseAlterSearchPropertyListStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
			}
			return nil, nil
		}
		// ALTER ENDPOINT
		if p.cur.Type == kwENDPOINT {
			p.advance() // consume ENDPOINT
			stmt, err := p.parseAlterEndpointStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER SERVER AUDIT [SPECIFICATION] / ALTER SERVER ROLE / ALTER SERVER CONFIGURATION
		if p.cur.Type == kwSERVER {
			next := p.peekNext()
			if next.Type == kwAUDIT {
				p.advance() // consume SERVER
				p.advance() // consume AUDIT
				if p.cur.Type == kwSPECIFICATION {
					p.advance() // consume SPECIFICATION
					stmt, err := p.parseAlterServerAuditSpecStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
				stmt, err := p.parseAlterServerAuditStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if next.Type == kwROLE {
				p.advance() // consume SERVER
				p.advance() // consume ROLE
				stmt, err := p.parseAlterServerRoleStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if next.Type == kwCONFIGURATION {
				p.advance() // consume SERVER
				p.advance() // consume CONFIGURATION
				stmt, err := p.parseAlterServerConfigurationStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
		}
		// ALTER DATABASE AUDIT SPECIFICATION (handled via kwDATABASE case)
		// ALTER AUTHORIZATION
		if p.cur.Type == kwAUTHORIZATION {
			p.advance() // consume AUTHORIZATION
			stmt, err := p.parseAlterAuthorizationStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER QUEUE (service broker)
		if p.cur.Type == kwQUEUE {
			p.advance() // consume QUEUE
			stmt, err := p.parseAlterQueueStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER SERVICE (service broker) / ALTER SERVICE MASTER KEY
		if p.cur.Type == kwSERVICE {
			next := p.peekNext()
			if next.Type == kwMASTER {
				// ALTER SERVICE MASTER KEY
				p.advance() // consume SERVICE
				p.advance() // consume MASTER
				if p.cur.Type == kwKEY {
					p.advance() // consume KEY
				}
				stmt, err := p.parseAlterServiceMasterKeyStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			p.advance() // consume SERVICE
			stmt, err := p.parseAlterServiceStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER ROUTE (service broker)
		if p.cur.Type == kwROUTE {
			p.advance() // consume ROUTE
			stmt, err := p.parseAlterRouteStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER REMOTE SERVICE BINDING (service broker)
		if p.cur.Type == kwREMOTE {
			p.advance() // consume REMOTE
			if p.cur.Type == kwSERVICE {
				p.advance() // consume SERVICE
			}
			if p.cur.Type == kwBINDING {
				p.advance() // consume BINDING
			}
			stmt, err := p.parseAlterRemoteServiceBindingStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER BROKER PRIORITY (service broker)
		if p.cur.Type == kwBROKER {
			p.advance() // consume BROKER
			if p.cur.Type == kwPRIORITY {
				p.advance() // consume PRIORITY
			}
			stmt, err := p.parseAlterBrokerPriorityStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER EVENT SESSION
		if p.cur.Type == kwEVENT {
			p.advance() // consume EVENT
			if p.cur.Type == kwSESSION {
				p.advance() // consume SESSION
				stmt, err := p.parseAlterEventSessionStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// ALTER MESSAGE TYPE (service broker)
		if p.cur.Type == kwMESSAGE {
			p.advance() // consume MESSAGE
			if p.cur.Type == kwTYPE {
				p.advance() // consume TYPE
			}
			stmt, err := p.parseAlterMessageTypeStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER CONTRACT (service broker)
		if p.cur.Type == kwCONTRACT {
			p.advance() // consume CONTRACT
			stmt, err := p.parseAlterContractStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER XML SCHEMA COLLECTION
		if p.cur.Type == kwXML {
			p.advance() // consume XML
			if p.cur.Type == kwSCHEMA {
				p.advance() // consume SCHEMA
				if p.cur.Type == kwCOLLECTION {
					p.advance() // consume COLLECTION
					stmt, err := p.parseAlterXmlSchemaCollectionStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
			}
			return nil, nil
		}
		// ALTER EXTERNAL DATA SOURCE / ALTER EXTERNAL RESOURCE POOL
		if p.cur.Type == kwEXTERNAL {
			p.advance() // consume EXTERNAL
			if p.cur.Type == kwDATA {
				p.advance() // consume DATA
				if p.cur.Type == kwSOURCE {
					p.advance() // consume SOURCE
				}
				stmt, err := p.parseAlterExternalDataSourceStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwRESOURCE {
				p.advance() // consume RESOURCE
				if p.cur.Type == kwPOOL {
					p.advance() // consume POOL
				}
				stmt, err := p.parseAlterExternalResourcePoolStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwLIBRARY {
				p.advance() // consume LIBRARY
				stmt, err := p.parseAlterExternalLibraryStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwLANGUAGE {
				p.advance() // consume LANGUAGE
				stmt, err := p.parseAlterExternalLanguageStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwMODEL {
				p.advance() // consume MODEL
				stmt, err := p.parseAlterExternalModelStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// ALTER AVAILABILITY GROUP
		if p.cur.Type == kwAVAILABILITY {
			p.advance() // consume AVAILABILITY
			if p.cur.Type == kwGROUP {
				p.advance() // consume GROUP
			}
			stmt, err := p.parseAlterAvailabilityGroupStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER WORKLOAD GROUP / ALTER WORKLOAD CLASSIFIER
		if p.cur.Type == kwWORKLOAD {
			p.advance() // consume WORKLOAD
			if p.cur.Type == kwCLASSIFIER {
				p.advance() // consume CLASSIFIER
				stmt, err := p.parseAlterWorkloadClassifierStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwGROUP {
				p.advance() // consume GROUP
			}
			stmt, err := p.parseAlterWorkloadGroupStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// ALTER RESOURCE POOL / ALTER RESOURCE GOVERNOR
		if p.cur.Type == kwRESOURCE {
			p.advance() // consume RESOURCE
			if p.cur.Type == kwPOOL {
				p.advance() // consume POOL
				stmt, err := p.parseAlterResourcePoolStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwGOVERNOR {
				p.advance() // consume GOVERNOR
				stmt, err := p.parseAlterResourceGovernorStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// ALTER FEDERATION
		if p.cur.Type == kwFEDERATION {
			p.advance() // consume FEDERATION
			stmt, err := p.parseAlterFederationStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		return nil, nil
	}
}

// parseDropOrSecurityStmt dispatches DROP to either parseDropStmt (for tables,
// views, etc.) or the security principal parsers (for USER, LOGIN, ROLE,
// APPLICATION ROLE).
func (p *Parser) parseDropOrSecurityStmt() (nodes.StmtNode, error) {
	loc := p.pos()
	next := p.peekNext()

	// Completion: after DROP, offer object type keywords.
	// We must advance past DROP to trigger checkCursor, but only return
	// candidates if we're now in collect mode. Otherwise fall through to
	// the normal dispatch which expects DROP not yet consumed — so we
	// only enter this block when we know the cursor is right after DROP.
	if p.completing && !p.collecting {
		// Speculatively advance to trigger checkCursor.
		savedLexerPos := p.lexer.pos
		savedCur := p.cur
		savedPrev := p.prev
		savedNextBuf := p.nextBuf
		savedHasNext := p.hasNext
		savedCollecting := p.collecting

		p.advance() // consume DROP — triggers checkCursor
		if p.collectMode() {
			for _, kw := range []int{
				kwTABLE, kwINDEX, kwVIEW, kwPROCEDURE, kwFUNCTION, kwTRIGGER,
				kwDATABASE, kwSCHEMA, kwTYPE, kwUSER, kwLOGIN, kwROLE,
				kwSTATISTICS, kwIF,
			} {
				p.addTokenCandidate(kw)
			}
			// Context-sensitive identifiers: SEQUENCE, SYNONYM
			p.addRuleCandidate("SEQUENCE")
			p.addRuleCandidate("SYNONYM")
			return nil, errCollecting
		}
		// Not collecting yet — restore state and let normal dispatch handle it.
		p.lexer.pos = savedLexerPos
		p.cur = savedCur
		p.prev = savedPrev
		p.nextBuf = savedNextBuf
		p.hasNext = savedHasNext
		p.collecting = savedCollecting
	}

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
		stmt, err := p.parseSecurityKeyStmtColumn("DROP")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwDATABASE:
		// Could be DROP DATABASE, DROP DATABASE AUDIT SPECIFICATION,
		// DROP DATABASE ENCRYPTION KEY, or DROP DATABASE SCOPED CREDENTIAL
		p.advance() // consume DROP
		p.advance() // consume DATABASE
		if p.cur.Type == kwAUDIT {
			p.advance() // consume AUDIT
			if p.cur.Type == kwSPECIFICATION {
				p.advance() // consume SPECIFICATION
			}
			stmt, err := p.parseDropDatabaseAuditSpecStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP DATABASE ENCRYPTION KEY / DROP DATABASE SCOPED CREDENTIAL
		if p.isAnyKeywordIdent() && (p.cur.Type == kwENCRYPTION || p.cur.Type == kwSCOPED) {
			stmt, err := p.parseSecurityKeyStmtDatabaseEncryption("DROP")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// Not AUDIT/ENCRYPTION/SCOPED — it's a regular DROP DATABASE. We already consumed DROP + DATABASE.
		// We need to call parseDropStmt logic but we've already consumed DROP + DATABASE.
		// Let's just handle it inline.
		dropStmt := &nodes.DropStmt{
			ObjectType: nodes.DropDatabase,
			Loc:        nodes.Loc{Start: loc, End: -1},
		}
		if p.cur.Type == kwIF {
			p.advance()
			p.match(kwEXISTS)
			dropStmt.IfExists = true
		}
		// Completion: after DROP DATABASE [IF EXISTS] → database_ref
		if p.collectMode() {
			p.addRuleCandidate("database_ref")
			return nil, errCollecting
		}
		var nameItems []nodes.Node
		for {
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
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
		dropStmt.Loc.End = p.prevEnd()
		return dropStmt, nil
	case kwUSER:
		p.advance() // consume DROP
		p.advance() // consume USER
		stmt, err := p.parseSecurityUserStmt("DROP")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwLOGIN:
		p.advance() // consume DROP
		p.advance() // consume LOGIN
		stmt, err := p.parseSecurityLoginStmt("DROP")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwROLE:
		p.advance() // consume DROP
		p.advance() // consume ROLE
		stmt, err := p.parseSecurityRoleStmt("DROP")
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	case kwSTATISTICS:
		// DROP STATISTICS table.stats_name [, ...]
		p.advance() // consume DROP
		p.advance() // consume STATISTICS
		stmt, err := p.parseDropStatisticsStmt()
		if err != nil {
			return nil, err
		}
		stmt.Loc.Start = loc
		return stmt, nil
	default:
		// Check for DROP APPLICATION ROLE
		if next.Type == kwAPPLICATION {
			p.advance() // consume DROP
			p.advance() // consume APPLICATION
			stmt, err := p.parseSecurityApplicationRoleStmt("DROP")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP ENDPOINT
		if next.Type == kwENDPOINT {
			p.advance() // consume DROP
			p.advance() // consume ENDPOINT
			stmt, err := p.parseDropEndpointStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP SERVER AUDIT [SPECIFICATION] / DROP SERVER ROLE
		if next.Type == kwSERVER {
			p.advance() // consume DROP
			p.advance() // consume SERVER
			if p.cur.Type == kwAUDIT {
				p.advance() // consume AUDIT
				if p.cur.Type == kwSPECIFICATION {
					p.advance() // consume SPECIFICATION
					stmt, err := p.parseDropServerAuditSpecStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
				stmt, err := p.parseDropServerAuditStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwROLE {
				p.advance() // consume ROLE
				stmt, err := p.parseDropServerRoleStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// DROP SERVICE BROKER objects
		if next.Type == kwMESSAGE {
			p.advance() // consume DROP
			p.advance() // consume MESSAGE
			if p.cur.Type == kwTYPE {
				p.advance() // consume TYPE
			}
			stmt, err := p.parseDropServiceBrokerStmt("MESSAGE TYPE")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if next.Type == kwCONTRACT {
			p.advance() // consume DROP
			p.advance() // consume CONTRACT
			stmt, err := p.parseDropServiceBrokerStmt("CONTRACT")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if next.Type == kwQUEUE {
			p.advance() // consume DROP
			p.advance() // consume QUEUE
			stmt, err := p.parseDropServiceBrokerStmt("QUEUE")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if next.Type == kwSERVICE {
			p.advance() // consume DROP
			p.advance() // consume SERVICE
			stmt, err := p.parseDropServiceBrokerStmt("SERVICE")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if next.Type == kwROUTE {
			p.advance() // consume DROP
			p.advance() // consume ROUTE
			stmt, err := p.parseDropServiceBrokerStmt("ROUTE")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if next.Type == kwREMOTE {
			p.advance() // consume DROP
			p.advance() // consume REMOTE
			if p.cur.Type == kwSERVICE {
				p.advance() // consume SERVICE
			}
			if p.cur.Type == kwBINDING {
				p.advance() // consume BINDING
			}
			stmt, err := p.parseDropServiceBrokerStmt("REMOTE SERVICE BINDING")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		if next.Type == kwBROKER {
			p.advance() // consume DROP
			p.advance() // consume BROKER
			if p.cur.Type == kwPRIORITY {
				p.advance() // consume PRIORITY
			}
			stmt, err := p.parseDropServiceBrokerStmt("BROKER PRIORITY")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP EXTERNAL DATA SOURCE / EXTERNAL TABLE / EXTERNAL FILE FORMAT
		if next.Type == kwEXTERNAL {
			p.advance() // consume DROP
			p.advance() // consume EXTERNAL
			stmt, err := p.parseDropExternalStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP AVAILABILITY GROUP
		if next.Type == kwAVAILABILITY {
			p.advance() // consume DROP
			p.advance() // consume AVAILABILITY
			if p.cur.Type == kwGROUP {
				p.advance() // consume GROUP
			}
			stmt, err := p.parseDropAvailabilityGroupStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP EVENT NOTIFICATION / EVENT SESSION
		if next.Type == kwEVENT {
			p.advance() // consume DROP
			p.advance() // consume EVENT
			if p.cur.Type == kwNOTIFICATION {
				p.advance() // consume NOTIFICATION
				stmt, err := p.parseDropEventNotificationStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwSESSION {
				p.advance() // consume SESSION
				stmt, err := p.parseDropEventSessionStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// DROP CRYPTOGRAPHIC PROVIDER
		if next.Type == kwCRYPTOGRAPHIC {
			p.advance() // consume DROP
			stmt, err := p.parseSecurityKeyStmt("DROP")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP WORKLOAD GROUP / DROP WORKLOAD CLASSIFIER
		if next.Type == kwWORKLOAD {
			p.advance() // consume DROP
			p.advance() // consume WORKLOAD
			if p.cur.Type == kwCLASSIFIER {
				p.advance() // consume CLASSIFIER
				stmt, err := p.parseDropWorkloadClassifierStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			if p.cur.Type == kwGROUP {
				p.advance() // consume GROUP
			}
			stmt, err := p.parseDropWorkloadGroupStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP RESOURCE POOL
		if next.Type == kwRESOURCE {
			p.advance() // consume DROP
			p.advance() // consume RESOURCE
			if p.cur.Type == kwPOOL {
				p.advance() // consume POOL
				stmt, err := p.parseDropResourcePoolStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// DROP SECURITY POLICY
		if next.Type == kwSECURITY {
			p.advance() // consume DROP
			p.advance() // consume SECURITY
			if p.cur.Type == kwPOLICY {
				p.advance() // consume POLICY
				stmt, err := p.parseDropSecurityPolicyStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// DROP SENSITIVITY CLASSIFICATION
		if next.Type == kwSENSITIVITY {
			p.advance() // consume DROP
			p.advance() // consume SENSITIVITY
			if p.cur.Type == kwCLASSIFICATION {
				p.advance() // consume CLASSIFICATION
			}
			stmt, err := p.parseDropSensitivityClassificationStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP [COUNTER] SIGNATURE
		if (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) &&
			(next.Type == kwSIGNATURE || next.Type == kwCOUNTER) {
			p.advance() // consume DROP
			isCounter := false
			if p.cur.Type == kwCOUNTER {
				isCounter = true
				p.advance() // consume COUNTER
			}
			if p.cur.Type == kwSIGNATURE {
				p.advance() // consume SIGNATURE
				stmt, err := p.parseSignatureStmt("DROP")
				if err != nil {
					return nil, err
				}
				stmt.IsCounter = isCounter
				stmt.Loc.Start = loc
				return stmt, nil
			}
			return nil, nil
		}
		// DROP FULLTEXT STOPLIST
		if next.Type == kwFULLTEXT {
			// Peek further: need to distinguish FULLTEXT STOPLIST from FULLTEXT INDEX/CATALOG
			// FULLTEXT INDEX/CATALOG are handled in parseDropStmt, but STOPLIST is a separate stmt type
			p.advance() // consume DROP
			p.advance() // consume FULLTEXT
			if p.cur.Type == kwSTOPLIST {
				p.advance() // consume STOPLIST
				stmt, err := p.parseDropFulltextStoplistStmt()
				if err != nil {
					return nil, err
				}
				stmt.Loc.Start = loc
				return stmt, nil
			}
			// Not STOPLIST — must be INDEX or CATALOG, handled by parseDropStmt.
			// We've already consumed DROP + FULLTEXT, so we need to handle inline.
			// Re-enter the generic drop logic for FULLTEXT INDEX/CATALOG
			dropStmt := &nodes.DropStmt{
				Loc: nodes.Loc{Start: loc, End: -1},
			}
			if p.cur.Type == kwINDEX {
				dropStmt.ObjectType = nodes.DropFulltextIndex
				p.advance()
			} else if p.cur.Type == kwCATALOG {
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
					if ref , _ := p.parseTableRef(); ref != nil {
						nameItems = append(nameItems, ref)
					}
				}
			} else {
				for {
					ref, err := p.parseTableRef()
					if err != nil {
						return nil, err
					}
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
			dropStmt.Loc.End = p.prevEnd()
			return dropStmt, nil
		}
		// DROP AGGREGATE
		if next.Type == kwAGGREGATE {
			p.advance() // consume DROP
			p.advance() // consume AGGREGATE
			stmt, err := p.parseDropAggregateStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP SEARCH PROPERTY LIST
		if next.Type == kwSEARCH {
			p.advance() // consume DROP
			p.advance() // consume SEARCH
			if p.cur.Type == kwPROPERTY {
				p.advance() // consume PROPERTY
				if p.cur.Type == kwLIST {
					p.advance() // consume LIST
					stmt, err := p.parseDropSearchPropertyListStmt()
					if err != nil {
						return nil, err
					}
					stmt.Loc.Start = loc
					return stmt, nil
				}
			}
			return nil, nil
		}
		// DROP MASTER KEY
		if next.Type == kwMASTER {
			p.advance() // consume DROP
			stmt, err := p.parseSecurityKeyStmt("DROP")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP SYMMETRIC KEY
		if next.Type == kwSYMMETRIC {
			p.advance() // consume DROP
			stmt, err := p.parseSecurityKeyStmt("DROP")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP ASYMMETRIC KEY
		if next.Type == kwASYMMETRIC {
			p.advance() // consume DROP
			stmt, err := p.parseSecurityKeyStmt("DROP")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP CERTIFICATE
		if next.Type == kwCERTIFICATE {
			p.advance() // consume DROP
			stmt, err := p.parseSecurityKeyStmt("DROP")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP CREDENTIAL
		if next.Type == kwCREDENTIAL {
			p.advance() // consume DROP
			stmt, err := p.parseSecurityKeyStmt("DROP")
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		// DROP FEDERATION
		if next.Type == kwFEDERATION {
			p.advance() // consume DROP
			p.advance() // consume FEDERATION
			stmt, err := p.parseDropFederationStmt()
			if err != nil {
				return nil, err
			}
			stmt.Loc.Start = loc
			return stmt, nil
		}
		stmt, err := p.parseDropStmt()
		if err != nil {
			return nil, err
		}
		return stmt, nil
	}
}

// parseLabelStmt parses a label: statement.
func (p *Parser) parseLabelStmt() (*nodes.LabelStmt, error) {
	loc := p.pos()
	label := p.cur.Str
	p.advance() // consume identifier
	p.advance() // consume :
	return &nodes.LabelStmt{
		Label: label,
		Loc:   nodes.Loc{Start: loc, End: p.prevEnd()},
	}, nil
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
	if p.completing {
		p.checkCursor()
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
	return Token{}, p.syntaxErrorAtCur()
}

// pos returns the byte position of the current token.
func (p *Parser) pos() int {
	return p.cur.Loc
}

// prevEnd returns the exclusive end offset of the previously consumed token.
func (p *Parser) prevEnd() int {
	return p.prev.End
}

// unexpectedToken returns a ParseError for the current token position.
func (p *Parser) unexpectedToken() error {
	return p.syntaxErrorAtCur()
}

// tokenText extracts the source text for a token. It uses the token's Str
// field if available, otherwise falls back to the single-character representation
// for operator tokens. Returns "" for EOF or unknown tokens.
func (p *Parser) tokenText(tok Token) string {
	if tok.Type == tokEOF {
		return ""
	}
	if tok.Str != "" {
		return tok.Str
	}
	// For single-char operator tokens (e.g. '+', '-', '*', '/', etc.)
	if tok.Type > 0 && tok.Type < 256 {
		return string(rune(tok.Type))
	}
	return ""
}

// lexerError wraps the lexer's error (p.lexer.Err) into a ParseError with
// position information from the current token.
func (p *Parser) lexerError() *ParseError {
	msg := "lexer error"
	if p.lexer.Err != nil {
		msg = p.lexer.Err.Error()
	}
	return &ParseError{
		Message:  msg,
		Position: p.cur.Loc,
		Source:   p.source,
	}
}

// syntaxErrorAtCur returns a ParseError with "at or near" context for the
// current token, or "at end of input" if the current token is EOF.
func (p *Parser) syntaxErrorAtCur() *ParseError {
	text := p.tokenText(p.cur)
	var msg string
	if text == "" {
		msg = "syntax error at end of input"
	} else {
		msg = fmt.Sprintf("syntax error at or near \"%s\"", text)
	}
	return &ParseError{
		Message:  msg,
		Position: p.cur.Loc,
		Source:   p.source,
	}
}

// syntaxErrorAtEnd returns a ParseError for unexpected end-of-input.
func (p *Parser) syntaxErrorAtEnd() *ParseError {
	return &ParseError{
		Message:  "syntax error at end of input",
		Position: p.cur.Loc,
		Source:   p.source,
	}
}

// newParseError creates a ParseError with Source automatically set from the parser.
func (p *Parser) newParseError(pos int, msg string) *ParseError {
	return &ParseError{
		Message:  msg,
		Position: pos,
		Source:   p.source,
	}
}

// ParseError represents a parse error with position information.
type ParseError struct {
	Message  string
	Position int
	Source   string // original SQL input (for error context extraction)
}

func (e *ParseError) Error() string {
	// If Source is set, generate a contextual error message.
	if e.Source != "" {
		if e.Position >= len(e.Source) {
			return "syntax error at end of input"
		}
		// Extract the token text at Position from Source.
		text := extractNearText(e.Source, e.Position)
		if text == "" {
			return "syntax error at end of input"
		}
		return fmt.Sprintf("syntax error at or near \"%s\"", text)
	}
	return e.Message
}

// extractNearText extracts a meaningful token near the given byte position in src.
// It returns the word or operator at that position, or "" if at end/whitespace-only.
func extractNearText(src string, pos int) string {
	if pos >= len(src) {
		return ""
	}
	// Skip any leading whitespace at pos.
	i := pos
	for i < len(src) && (src[i] == ' ' || src[i] == '\t' || src[i] == '\n' || src[i] == '\r') {
		i++
	}
	if i >= len(src) {
		return ""
	}
	// If it's a letter or underscore, read an identifier/keyword.
	if isIdentStart(src[i]) {
		j := i
		for j < len(src) && isIdentCont(src[j]) {
			j++
		}
		return src[i:j]
	}
	// For quoted strings, return the first few chars.
	if src[i] == '\'' {
		j := i + 1
		for j < len(src) && src[j] != '\'' {
			j++
		}
		if j < len(src) {
			j++ // include closing quote
		}
		return src[i:j]
	}
	// For operators and punctuation, return single char (or two-char operators).
	if i+1 < len(src) {
		two := src[i : i+2]
		if two == "<>" || two == "<=" || two == ">=" || two == "!=" {
			return two
		}
	}
	return string(src[i])
}

// Note: isIdentStart and isIdentCont are defined in lexer.go and reused here.
