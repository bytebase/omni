package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseAnalyzeTableStmt parses an ANALYZE TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/analyze-table.html
//
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name [, tbl_name] ...
//
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name
//	    UPDATE HISTOGRAM ON col_name [, col_name] ...
//	        [WITH N BUCKETS]
//
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name
//	    UPDATE HISTOGRAM ON col_name [USING DATA 'json_data']
//
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name
//	    DROP HISTOGRAM ON col_name [, col_name] ...
func (p *Parser) parseAnalyzeTableStmt() (*nodes.AnalyzeTableStmt, error) {
	start := p.pos()
	p.advance() // consume ANALYZE

	stmt := &nodes.AnalyzeTableStmt{Loc: nodes.Loc{Start: start}}

	// Optional NO_WRITE_TO_BINLOG | LOCAL
	if p.cur.Type == kwLOCAL {
		stmt.NoWriteToBinlog = true
		p.advance()
	} else if p.cur.Type == kwNO_WRITE_TO_BINLOG {
		stmt.NoWriteToBinlog = true
		p.advance()
	}

	// TABLE
	p.match(kwTABLE)

	// Completion: after ANALYZE [LOCAL] TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Table list
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// UPDATE HISTOGRAM ON col [, col] [WITH N BUCKETS] | [USING DATA 'json_data']
	if p.cur.Type == kwUPDATE {
		p.advance()
		if p.cur.Type == kwHISTOGRAM {
			p.advance()
		}
		stmt.HistogramOp = "UPDATE"
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		for {
			col, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.HistogramColumns = append(stmt.HistogramColumns, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		// WITH N BUCKETS
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.Buckets = int(p.cur.Ival)
				p.advance()
			}
			// BUCKETS keyword
			if p.cur.Type == kwBUCKETS {
				p.advance()
			}
		}
		// USING DATA 'json_data'
		if p.cur.Type == kwUSING {
			p.advance() // consume USING
			if p.cur.Type == kwDATA {
				p.advance() // consume DATA
			}
			if p.cur.Type == tokSCONST {
				stmt.UsingData = p.cur.Str
				p.advance()
			}
		}
	}

	// DROP HISTOGRAM ON col [, col]
	if p.cur.Type == kwDROP {
		p.advance()
		if p.cur.Type == kwHISTOGRAM {
			p.advance()
		}
		stmt.HistogramOp = "DROP"
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		for {
			col, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.HistogramColumns = append(stmt.HistogramColumns, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseOptimizeTableStmt parses an OPTIMIZE TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/optimize-table.html
//
//	OPTIMIZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name [, tbl_name] ...
func (p *Parser) parseOptimizeTableStmt() (*nodes.OptimizeTableStmt, error) {
	start := p.pos()
	p.advance() // consume OPTIMIZE

	stmt := &nodes.OptimizeTableStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwLOCAL {
		stmt.NoWriteToBinlog = true
		p.advance()
	} else if p.cur.Type == kwNO_WRITE_TO_BINLOG {
		stmt.NoWriteToBinlog = true
		p.advance()
	}

	p.match(kwTABLE)

	// Completion: after OPTIMIZE TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCheckTableStmt parses a CHECK TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/check-table.html
//
//	CHECK TABLE tbl_name [, tbl_name] ... [option] ...
func (p *Parser) parseCheckTableStmt() (*nodes.CheckTableStmt, error) {
	start := p.pos()
	p.advance() // consume CHECK

	p.match(kwTABLE)

	// Completion: after CHECK TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.CheckTableStmt{Loc: nodes.Loc{Start: start}}

	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// Optional check options: FOR UPGRADE, QUICK, FAST, MEDIUM, EXTENDED, CHANGED
	for {
		if p.cur.Type == kwFOR {
			p.advance()
			if p.cur.Type == kwUPGRADE {
				stmt.Options = append(stmt.Options, "FOR UPGRADE")
				p.advance()
			}
		} else if p.cur.Type == kwQUICK {
			stmt.Options = append(stmt.Options, "QUICK")
			p.advance()
		} else if p.cur.Type == kwEXTENDED {
			stmt.Options = append(stmt.Options, "EXTENDED")
			p.advance()
		} else if p.cur.Type == kwFAST || p.cur.Type == kwMEDIUM || p.cur.Type == kwCHANGED {
			stmt.Options = append(stmt.Options, strings.ToUpper(p.cur.Str))
			p.advance()
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRepairTableStmt parses a REPAIR TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/repair-table.html
//
//	REPAIR [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name [, tbl_name] ... [QUICK] [EXTENDED] [USE_FRM]
func (p *Parser) parseRepairTableStmt() (*nodes.RepairTableStmt, error) {
	start := p.pos()
	p.advance() // consume REPAIR

	stmt := &nodes.RepairTableStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwLOCAL {
		stmt.NoWriteToBinlog = true
		p.advance()
	} else if p.cur.Type == kwNO_WRITE_TO_BINLOG {
		stmt.NoWriteToBinlog = true
		p.advance()
	}

	p.match(kwTABLE)

	// Completion: after REPAIR TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// Options: QUICK, EXTENDED, USE_FRM
	for {
		if p.cur.Type == kwQUICK {
			stmt.Quick = true
			p.advance()
		} else if p.cur.Type == kwEXTENDED {
			stmt.Extended = true
			p.advance()
		} else if p.cur.Type == kwUSE_FRM {
			stmt.UseFrm = true
			p.advance()
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseFlushStmt parses a FLUSH statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/flush.html
//
//	FLUSH [NO_WRITE_TO_BINLOG | LOCAL] {
//	    flush_option [, flush_option] ...
//	  | tables_option
//	}
//
//	flush_option: {
//	    BINARY LOGS
//	  | ENGINE LOGS
//	  | ERROR LOGS
//	  | GENERAL LOGS
//	  | HOSTS
//	  | LOGS
//	  | PRIVILEGES
//	  | OPTIMIZER_COSTS
//	  | RELAY LOGS [FOR CHANNEL channel]
//	  | SLOW LOGS
//	  | STATUS
//	  | USER_RESOURCES
//	}
//
//	tables_option: {
//	    table_synonym
//	  | table_synonym tbl_name [, tbl_name] ...
//	  | table_synonym WITH READ LOCK
//	  | table_synonym tbl_name [, tbl_name] ... WITH READ LOCK
//	  | table_synonym tbl_name [, tbl_name] ... FOR EXPORT
//	}
//
//	table_synonym: { TABLE | TABLES }
func (p *Parser) parseFlushStmt() (*nodes.FlushStmt, error) {
	start := p.pos()
	p.advance() // consume FLUSH

	// Completion: after FLUSH, offer flush option keywords.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwPRIVILEGES)
		p.addTokenCandidate(kwTABLES)
		p.addTokenCandidate(kwLOGS)
		p.addTokenCandidate(kwSTATUS)
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.FlushStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwLOCAL {
		stmt.NoWriteToBinlog = true
		p.advance()
	} else if p.cur.Type == kwNO_WRITE_TO_BINLOG {
		stmt.NoWriteToBinlog = true
		p.advance()
	}

	// FLUSH TABLES [tbl_name, ...] [WITH READ LOCK | FOR EXPORT]
	if p.cur.Type == kwTABLES || p.cur.Type == kwTABLE {
		stmt.Options = append(stmt.Options, "TABLES")
		p.advance()
		// Optional table list (but not WITH or FOR which are clause keywords)
		if (p.isIdentToken() || p.cur.Type == tokSCONST) && p.cur.Type != kwWITH && p.cur.Type != kwFOR {
			for {
				ref, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				stmt.Tables = append(stmt.Tables, ref)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
		}
		// WITH READ LOCK
		if p.cur.Type == kwWITH {
			p.advance()
			if p.cur.Type == kwREAD {
				p.advance()
				if p.cur.Type == kwLOCK {
					p.advance()
				}
				stmt.WithReadLock = true
			}
		}
		// FOR EXPORT
		if p.cur.Type == kwFOR {
			p.advance()
			if p.cur.Type == kwEXPORT {
				p.advance()
			}
			stmt.ForExport = true
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Flush options list
	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}
		opt := p.parseFlushOption(stmt)
		if opt == "" {
			break
		}
		stmt.Options = append(stmt.Options, opt)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseFlushOption parses a single flush_option and returns the normalized name.
// Multi-word options (ENGINE LOGS, ERROR LOGS, etc.) are combined into one string.
func (p *Parser) parseFlushOption(stmt *nodes.FlushStmt) string {
	// Handle keyword tokens that are also flush options
	switch p.cur.Type {
	case kwBINARY:
		p.advance()
		if p.cur.Type == kwLOGS {
			p.advance()
			return "BINARY LOGS"
		}
		return "BINARY"
	case kwSTATUS:
		p.advance()
		return "STATUS"
	case kwPRIVILEGES:
		p.advance()
		return "PRIVILEGES"
	case kwOPTIMIZER_COSTS:
		p.advance()
		return "OPTIMIZER_COSTS"
	}

	if !p.isIdentToken() {
		return ""
	}

	name := strings.ToUpper(p.cur.Str)
	p.advance()

	// Check for multi-word options: ENGINE LOGS, ERROR LOGS, GENERAL LOGS, RELAY LOGS, SLOW LOGS
	switch name {
	case "ENGINE", "ERROR", "GENERAL", "SLOW":
		if p.cur.Type == kwLOGS {
			p.advance()
			return name + " LOGS"
		}
		return name
	case "RELAY":
		if p.cur.Type == kwLOGS {
			p.advance()
			opt := "RELAY LOGS"
			// Optional FOR CHANNEL channel
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
				if p.cur.Type == kwCHANNEL {
					p.advance() // consume CHANNEL
					if p.isIdentToken() || p.cur.Type == tokSCONST {
						if p.cur.Type == tokSCONST {
							stmt.RelayChannel = p.cur.Str
						} else {
							stmt.RelayChannel = p.cur.Str
						}
						p.advance()
					}
				}
			}
			return opt
		}
		return name
	}

	return name
}

// parseResetStmt parses a RESET statement.
//
//	RESET reset_option [, reset_option] ...
func (p *Parser) parseResetStmt() (*nodes.FlushStmt, error) {
	start := p.pos()
	p.advance() // consume RESET

	stmt := &nodes.FlushStmt{Loc: nodes.Loc{Start: start}}

	// Require at least one reset option
	if p.cur.Type == tokEOF || p.cur.Type == ';' {
		return nil, p.syntaxErrorAtCur()
	}

	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}
		// Reset option words must accept any keyword (e.g., QUERY CACHE).
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			name, _, err := p.parseKeywordOrIdent()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, name)
		} else {
			break
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseKillStmt parses a KILL statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/kill.html
//
//	KILL [CONNECTION | QUERY] processlist_id
func (p *Parser) parseKillStmt() (*nodes.KillStmt, error) {
	start := p.pos()
	p.advance() // consume KILL

	stmt := &nodes.KillStmt{Loc: nodes.Loc{Start: start}}

	// Optional CONNECTION | QUERY
	if p.cur.Type == kwCONNECTION {
		p.advance()
	} else if p.cur.Type == kwQUERY {
		stmt.Query = true
		p.advance()
	}

	// Process ID
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.ConnectionID = expr

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCallStmt parses a CALL statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/call.html
//
//	CALL sp_name([parameter[,...]])
//	CALL sp_name[()]
func (p *Parser) parseCallStmt() (*nodes.CallStmt, error) {
	start := p.pos()
	p.advance() // consume CALL

	// Completion: after CALL keyword, offer procedure_ref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("procedure_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.CallStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Optional argument list in parentheses
	if p.cur.Type == '(' {
		p.advance() // consume '('
		if p.cur.Type != ')' {
			for {
				arg, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if arg == nil {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.Args = append(stmt.Args, arg)
				if p.cur.Type != ',' {
					break
				}
				p.advance() // consume ','
			}
		}
		p.expect(')')
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseHandlerOpenStmt parses a HANDLER ... OPEN statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/handler.html
//
//	HANDLER tbl_name OPEN [ [AS] alias]
func (p *Parser) parseHandlerOpenStmt(start int, table *nodes.TableRef) (*nodes.HandlerOpenStmt, error) {
	p.advance() // consume OPEN

	stmt := &nodes.HandlerOpenStmt{
		Loc:   nodes.Loc{Start: start},
		Table: table,
	}

	// Optional AS alias
	if p.cur.Type == kwAS {
		p.advance() // consume AS
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Alias = alias
	} else if p.isIdentToken() && p.cur.Type != ';' && p.cur.Type != tokEOF {
		alias, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Alias = alias
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseHandlerReadStmt parses a HANDLER ... READ statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/handler.html
//
//	HANDLER tbl_name READ index_name { = | <= | >= | < | > } (value1,value2,...)
//	    [ WHERE where_condition ] [LIMIT ... ]
//	HANDLER tbl_name READ index_name { FIRST | NEXT | PREV | LAST }
//	    [ WHERE where_condition ] [LIMIT ... ]
//	HANDLER tbl_name READ { FIRST | NEXT }
//	    [ WHERE where_condition ] [LIMIT ... ]
func (p *Parser) parseHandlerReadStmt(start int, table *nodes.TableRef) (*nodes.HandlerReadStmt, error) {
	p.advance() // consume READ

	stmt := &nodes.HandlerReadStmt{
		Loc:   nodes.Loc{Start: start},
		Table: table,
	}

	// Determine if this is a direction keyword (FIRST/NEXT/PREV/LAST) or an index name
	switch p.cur.Type {
	case kwFIRST:
		stmt.Direction = "FIRST"
		p.advance()
	case kwNEXT:
		stmt.Direction = "NEXT"
		p.advance()
	case kwPREV:
		stmt.Direction = "PREV"
		p.advance()
	case kwLAST:
		stmt.Direction = "LAST"
		p.advance()
	default:
		// Must be an index name followed by direction or comparison
		if p.isIdentToken() {
			idx, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Index = idx

			// Direction keyword after index name
			switch p.cur.Type {
			case kwFIRST:
				stmt.Direction = "FIRST"
				p.advance()
			case kwNEXT:
				stmt.Direction = "NEXT"
				p.advance()
			case kwPREV:
				stmt.Direction = "PREV"
				p.advance()
			case kwLAST:
				stmt.Direction = "LAST"
				p.advance()
			case '=', '<', '>':
				// Comparison operator with value list — skip operator and value list
				p.advance()
				if p.cur.Type == '=' || p.cur.Type == '>' {
					p.advance() // for <= or >=
				}
				if p.cur.Type == '(' {
					p.advance()
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						p.advance()
					}
					if p.cur.Type == ')' {
						p.advance()
					}
				}
				stmt.Direction = "NEXT" // default
			}
		}
	}

	// Optional WHERE clause
	if p.cur.Type == kwWHERE {
		p.advance() // consume WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if !p.collectMode() && where == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Where = where
	}

	// Optional LIMIT clause
	if _, ok := p.match(kwLIMIT); ok {
		lim, err := p.parseLimitClause()
		if err != nil {
			return nil, err
		}
		stmt.Limit = lim
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseHandlerCloseStmt parses a HANDLER ... CLOSE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/handler.html
//
//	HANDLER tbl_name CLOSE
func (p *Parser) parseHandlerCloseStmt(start int, table *nodes.TableRef) (*nodes.HandlerCloseStmt, error) {
	p.advance() // consume CLOSE

	stmt := &nodes.HandlerCloseStmt{
		Loc:   nodes.Loc{Start: start},
		Table: table,
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseHandlerStmt parses a HANDLER statement and dispatches to OPEN/READ/CLOSE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/handler.html
//
//	HANDLER tbl_name OPEN [ [AS] alias]
//	HANDLER tbl_name READ index_name { = | <= | >= | < | > } (value1,value2,...) [ WHERE ... ] [LIMIT ... ]
//	HANDLER tbl_name READ index_name { FIRST | NEXT | PREV | LAST } [ WHERE ... ] [LIMIT ... ]
//	HANDLER tbl_name READ { FIRST | NEXT } [ WHERE ... ] [LIMIT ... ]
//	HANDLER tbl_name CLOSE
func (p *Parser) parseHandlerStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume HANDLER

	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwOPEN:
		return p.parseHandlerOpenStmt(start, table)
	case kwREAD:
		return p.parseHandlerReadStmt(start, table)
	case kwCLOSE:
		return p.parseHandlerCloseStmt(start, table)
	default:
		return nil, &ParseError{
			Message:  "expected OPEN, READ, or CLOSE after HANDLER table_name",
			Position: p.cur.Loc,
		}
	}
}

// parseDoStmt parses a DO statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/do.html
//
//	DO expr [, expr] ...
func (p *Parser) parseDoStmt() (*nodes.DoStmt, error) {
	start := p.pos()
	p.advance() // consume DO

	stmt := &nodes.DoStmt{Loc: nodes.Loc{Start: start}}

	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Exprs = append(stmt.Exprs, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseChecksumTableStmt parses a CHECKSUM TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/checksum-table.html
//
//	CHECKSUM TABLE tbl_name [, tbl_name] ... [QUICK | EXTENDED]
func (p *Parser) parseChecksumTableStmt() (*nodes.ChecksumTableStmt, error) {
	start := p.pos()
	p.advance() // consume CHECKSUM

	p.match(kwTABLE)

	stmt := &nodes.ChecksumTableStmt{Loc: nodes.Loc{Start: start}}

	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// Optional QUICK | EXTENDED
	if p.cur.Type == kwQUICK {
		stmt.Quick = true
		p.advance()
	} else if p.cur.Type == kwEXTENDED {
		stmt.Extended = true
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseShutdownStmt parses a SHUTDOWN statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/shutdown.html
//
//	SHUTDOWN
func (p *Parser) parseShutdownStmt() (*nodes.ShutdownStmt, error) {
	start := p.pos()
	p.advance() // consume SHUTDOWN

	stmt := &nodes.ShutdownStmt{Loc: nodes.Loc{Start: start}}
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseRestartStmt parses a RESTART statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/restart.html
//
//	RESTART
func (p *Parser) parseRestartStmt() (*nodes.RestartStmt, error) {
	start := p.pos()
	p.advance() // consume RESTART

	stmt := &nodes.RestartStmt{Loc: nodes.Loc{Start: start}}
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCloneStmt parses a CLONE statement (MySQL 8.0.17+).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/clone.html
//
//	CLONE LOCAL DATA DIRECTORY [=] 'clone_dir'
//	CLONE INSTANCE FROM 'user'@'host':port
//	  IDENTIFIED BY 'password'
//	  [DATA DIRECTORY [=] 'clone_dir']
//	  [REQUIRE [NO] SSL]
func (p *Parser) parseCloneStmt() (*nodes.CloneStmt, error) {
	start := p.pos()
	p.advance() // consume CLONE

	stmt := &nodes.CloneStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwLOCAL {
		// CLONE LOCAL DATA DIRECTORY [=] 'clone_dir'
		stmt.Local = true
		p.advance() // consume LOCAL
		p.match(kwDATA)
		p.match(kwDIRECTORY)

		// Optional '='
		if p.cur.Type == '=' {
			p.advance()
		}

		// clone_dir string
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{
				Message:  "expected string literal for DATA DIRECTORY",
				Position: p.cur.Loc,
			}
		}
		stmt.Directory = p.cur.Str
		p.advance()
	} else if p.cur.Type == kwINSTANCE {
		// CLONE INSTANCE FROM 'user'@'host':port IDENTIFIED BY 'password' ...
		p.advance() // consume INSTANCE
		p.match(kwFROM)

		// 'user'
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{
				Message:  "expected string literal for user",
				Position: p.cur.Loc,
			}
		}
		stmt.User = p.cur.Str
		p.advance()

		// @ - lexer scans @ as tokIDENT with Str="@"
		if p.cur.Type == tokIDENT && p.cur.Str == "@" {
			p.advance()
		} else {
			return nil, &ParseError{
				Message:  "expected @ after user",
				Position: p.cur.Loc,
			}
		}

		// 'host'
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{
				Message:  "expected string literal for host",
				Position: p.cur.Loc,
			}
		}
		stmt.Host = p.cur.Str
		p.advance()

		// :port
		p.expect(':')
		if p.cur.Type != tokICONST {
			return nil, &ParseError{
				Message:  "expected integer for port",
				Position: p.cur.Loc,
			}
		}
		stmt.Port = p.cur.Ival
		p.advance()

		// IDENTIFIED BY 'password'
		p.match(kwIDENTIFIED)
		p.match(kwBY)
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{
				Message:  "expected string literal for password",
				Position: p.cur.Loc,
			}
		}
		stmt.Password = p.cur.Str
		p.advance()

		// Optional DATA DIRECTORY [=] 'clone_dir'
		if p.cur.Type == kwDATA {
			p.advance() // consume DATA
			p.match(kwDIRECTORY)
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for DATA DIRECTORY",
					Position: p.cur.Loc,
				}
			}
			stmt.Directory = p.cur.Str
			p.advance()
		}

		// Optional REQUIRE [NO] SSL
		if p.cur.Type == kwREQUIRE {
			p.advance() // consume REQUIRE
			if p.cur.Type == kwNO {
				p.advance() // consume NO
				p.match(kwSSL)
				f := false
				stmt.RequireSSL = &f
			} else {
				p.match(kwSSL)
				t := true
				stmt.RequireSSL = &t
			}
		}
	} else {
		return nil, &ParseError{
			Message:  "expected LOCAL or INSTANCE after CLONE",
			Position: p.cur.Loc,
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseInstallStmt parses INSTALL PLUGIN or INSTALL COMPONENT.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/install-plugin.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/install-component.html
//
//	INSTALL PLUGIN plugin_name SONAME 'shared_library_name'
//	INSTALL COMPONENT component_name [, component_name] ...
func (p *Parser) parseInstallStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume INSTALL

	switch p.cur.Type {
	case kwPLUGIN:
		return p.parseInstallPluginStmt(start)
	case kwCOMPONENT:
		return p.parseInstallComponentStmt(start)
	default:
		return nil, &ParseError{
			Message:  "expected PLUGIN or COMPONENT after INSTALL",
			Position: p.cur.Loc,
		}
	}
}

// parseInstallPluginStmt parses an INSTALL PLUGIN statement.
//
//	INSTALL PLUGIN plugin_name SONAME 'shared_library_name'
func (p *Parser) parseInstallPluginStmt(start int) (*nodes.InstallPluginStmt, error) {
	p.advance() // consume PLUGIN

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	p.match(kwSONAME)

	soname := p.cur.Str
	p.expect(tokSCONST)

	stmt := &nodes.InstallPluginStmt{
		Loc:        nodes.Loc{Start: start},
		PluginName: name,
		Soname:     soname,
	}
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseInstallComponentStmt parses an INSTALL COMPONENT statement.
//
//	INSTALL COMPONENT component_name [, component_name] ...
//	    [SET variable = expr [, variable = expr] ...]
//
//	variable: {
//	    {GLOBAL | @@GLOBAL.} [component_prefix.]system_var_name
//	  | {PERSIST | @@PERSIST.} [component_prefix.]system_var_name
//	}
func (p *Parser) parseInstallComponentStmt(start int) (*nodes.InstallComponentStmt, error) {
	p.advance() // consume COMPONENT

	stmt := &nodes.InstallComponentStmt{Loc: nodes.Loc{Start: start}}

	for {
		comp := p.cur.Str
		p.expect(tokSCONST)
		stmt.Components = append(stmt.Components, comp)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// Optional SET clause
	if p.cur.Type == kwSET {
		p.advance() // consume SET
		for {
			asgn, err := p.parseSetAssignment()
			if err != nil {
				return nil, err
			}
			stmt.SetVars = append(stmt.SetVars, asgn)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseUninstallStmt parses UNINSTALL PLUGIN or UNINSTALL COMPONENT.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/uninstall-plugin.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/uninstall-component.html
//
//	UNINSTALL PLUGIN plugin_name
//	UNINSTALL COMPONENT component_name [, component_name] ...
func (p *Parser) parseUninstallStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume UNINSTALL

	switch p.cur.Type {
	case kwPLUGIN:
		return p.parseUninstallPluginStmt(start)
	case kwCOMPONENT:
		return p.parseUninstallComponentStmt(start)
	default:
		return nil, &ParseError{
			Message:  "expected PLUGIN or COMPONENT after UNINSTALL",
			Position: p.cur.Loc,
		}
	}
}

// parseUninstallPluginStmt parses an UNINSTALL PLUGIN statement.
//
//	UNINSTALL PLUGIN plugin_name
func (p *Parser) parseUninstallPluginStmt(start int) (*nodes.UninstallPluginStmt, error) {
	p.advance() // consume PLUGIN

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.UninstallPluginStmt{
		Loc:        nodes.Loc{Start: start},
		PluginName: name,
	}
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseUninstallComponentStmt parses an UNINSTALL COMPONENT statement.
//
//	UNINSTALL COMPONENT component_name [, component_name] ...
func (p *Parser) parseUninstallComponentStmt(start int) (*nodes.UninstallComponentStmt, error) {
	p.advance() // consume COMPONENT

	stmt := &nodes.UninstallComponentStmt{Loc: nodes.Loc{Start: start}}

	for {
		comp := p.cur.Str
		p.expect(tokSCONST)
		stmt.Components = append(stmt.Components, comp)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCreateTablespaceStmt parses a CREATE [UNDO] TABLESPACE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-tablespace.html
//
//	CREATE [UNDO] TABLESPACE tablespace_name
//	    [ADD DATAFILE 'file_name']
//	    [AUTOEXTEND_SIZE [=] value]
//	    [FILE_BLOCK_SIZE = value]
//	    [ENCRYPTION [=] {'Y' | 'N'}]
//	    [USE LOGFILE GROUP logfile_group]
//	    [EXTENT_SIZE [=] extent_size]
//	    [INITIAL_SIZE [=] initial_size]
//	    [MAX_SIZE [=] max_size]
//	    [NODEGROUP [=] nodegroup_id]
//	    [WAIT]
//	    [COMMENT [=] 'string']
//	    [ENGINE [=] engine_name]
//	    [ENGINE_ATTRIBUTE [=] 'string']
func (p *Parser) parseCreateTablespaceStmt(start int, undo bool) (*nodes.CreateTablespaceStmt, error) {
	p.advance() // consume TABLESPACE

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.CreateTablespaceStmt{
		Loc:  nodes.Loc{Start: start},
		Undo: undo,
		Name: name,
	}

	// Parse optional clauses
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		switch {
		case p.cur.Type == kwADD:
			p.advance() // consume ADD
			p.match(kwDATAFILE)
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for DATAFILE",
					Position: p.cur.Loc,
				}
			}
			stmt.DataFile = p.cur.Str
			p.advance()
		case p.cur.Type == kwAUTOEXTEND_SIZE:
			p.advance()
			p.match('=')
			stmt.AutoextendSize = p.parseSizeValue()
		case p.cur.Type == kwFILE_BLOCK_SIZE:
			p.advance()
			p.match('=')
			stmt.FileBlockSize = p.cur.Str
			p.advance()
		case p.cur.Type == kwENCRYPTION:
			p.advance()
			p.match('=')
			if p.cur.Type == tokSCONST {
				stmt.Encryption = p.cur.Str
				p.advance()
			}
		case p.cur.Type == kwUSE:
			p.advance() // consume USE
			p.match(kwLOGFILE)
			p.match(kwGROUP)
			lg, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.UseLogfileGroup = lg
		case p.cur.Type == kwEXTENT_SIZE:
			p.advance()
			p.match('=')
			stmt.ExtentSize = p.parseSizeValue()
		case p.cur.Type == kwINITIAL_SIZE:
			p.advance()
			p.match('=')
			stmt.InitialSize = p.parseSizeValue()
		case p.cur.Type == kwMAX_SIZE:
			p.advance()
			p.match('=')
			stmt.MaxSize = p.parseSizeValue()
		case p.cur.Type == kwNODEGROUP:
			p.advance()
			p.match('=')
			stmt.NodeGroup = p.cur.Str
			p.advance()
		case p.cur.Type == kwWAIT:
			p.advance()
			stmt.Wait = true
		case p.cur.Type == kwCOMMENT:
			p.advance()
			p.match('=')
			if p.cur.Type == tokSCONST {
				stmt.Comment = p.cur.Str
				p.advance()
			}
		case p.cur.Type == kwENGINE:
			p.advance()
			p.match('=')
			ename, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Engine = ename
		case p.cur.Type == kwENGINE_ATTRIBUTE:
			p.advance()
			p.match('=')
			if p.cur.Type == tokSCONST {
				stmt.EngineAttribute = p.cur.Str
				p.advance()
			}
		default:
			goto done
		}
	}
done:
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterTablespaceStmt parses an ALTER [UNDO] TABLESPACE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-tablespace.html
//
//	ALTER [UNDO] TABLESPACE tablespace_name
//	    [ADD DATAFILE 'file_name']
//	    [DROP DATAFILE 'file_name']
//	    [INITIAL_SIZE [=] size]
//	    [WAIT]
//	    [RENAME TO tablespace_name]
//	    [AUTOEXTEND_SIZE [=] 'value']
//	    [SET {ACTIVE | INACTIVE}]
//	    [ENCRYPTION [=] {'Y' | 'N'}]
//	    [ENGINE [=] engine_name]
//	    [ENGINE_ATTRIBUTE [=] 'string']
func (p *Parser) parseAlterTablespaceStmt(start int, undo bool) (*nodes.AlterTablespaceStmt, error) {
	p.advance() // consume TABLESPACE

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.AlterTablespaceStmt{
		Loc:  nodes.Loc{Start: start},
		Undo: undo,
		Name: name,
	}

	// Parse optional clauses
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		switch {
		case p.cur.Type == kwADD:
			p.advance() // consume ADD
			p.match(kwDATAFILE)
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for DATAFILE",
					Position: p.cur.Loc,
				}
			}
			stmt.AddDataFile = p.cur.Str
			p.advance()
		case p.cur.Type == kwDROP:
			p.advance() // consume DROP
			p.match(kwDATAFILE)
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for DATAFILE",
					Position: p.cur.Loc,
				}
			}
			stmt.DropDataFile = p.cur.Str
			p.advance()
		case p.cur.Type == kwINITIAL_SIZE:
			p.advance()
			p.match('=')
			stmt.InitialSize = p.parseSizeValue()
		case p.cur.Type == kwWAIT:
			p.advance()
			stmt.Wait = true
		case p.cur.Type == kwRENAME:
			p.advance() // consume RENAME
			p.match(kwTO)
			rn, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.RenameTo = rn
		case p.cur.Type == kwAUTOEXTEND_SIZE:
			p.advance()
			p.match('=')
			stmt.AutoextendSize = p.parseSizeValue()
		case p.cur.Type == kwSET:
			p.advance() // consume SET
			if p.cur.Type == kwACTIVE {
				stmt.SetActive = true
				p.advance()
			} else if p.cur.Type == kwINACTIVE {
				stmt.SetInactive = true
				p.advance()
			}
		case p.cur.Type == kwENCRYPTION:
			p.advance()
			p.match('=')
			if p.cur.Type == tokSCONST {
				stmt.Encryption = p.cur.Str
				p.advance()
			}
		case p.cur.Type == kwENGINE:
			p.advance()
			p.match('=')
			ename, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Engine = ename
		case p.cur.Type == kwENGINE_ATTRIBUTE:
			p.advance()
			p.match('=')
			if p.cur.Type == tokSCONST {
				stmt.EngineAttribute = p.cur.Str
				p.advance()
			}
		default:
			goto done
		}
	}
done:
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropTablespaceStmt parses a DROP [UNDO] TABLESPACE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-tablespace.html
//
//	DROP [UNDO] TABLESPACE tablespace_name [ENGINE [=] engine_name]
func (p *Parser) parseDropTablespaceStmt(start int, undo bool) (*nodes.DropTablespaceStmt, error) {
	p.advance() // consume TABLESPACE

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.DropTablespaceStmt{
		Loc:  nodes.Loc{Start: start},
		Undo: undo,
		Name: name,
	}

	// Optional ENGINE [=] engine_name
	if p.cur.Type == kwENGINE {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		ename, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Engine = ename
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseServerOptions parses the OPTIONS (...) clause for SERVER statements.
// It consumes everything inside the parentheses as raw option strings.
func (p *Parser) parseServerOptions() ([]string, error) {
	p.match(kwOPTIONS)
	p.expect('(')

	var opts []string
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Each option is: keyword value
		// e.g., HOST 'host_name', DATABASE 'db_name', USER 'user_name', etc.
		// Must accept any keyword since option names (HOST, DATABASE, USER, etc.)
		// may be registered keywords.
		var optStr string
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			optName, _, _ := p.parseKeywordOrIdent()
			optStr = optName
		} else {
			optStr = p.cur.Str
			p.advance()
		}
		// value
		if p.cur.Type == tokSCONST {
			optStr += " " + p.cur.Str
			p.advance()
		} else if p.cur.Type == tokICONST {
			optStr += " " + p.cur.Str
			p.advance()
		}
		opts = append(opts, optStr)
		if p.cur.Type == ',' {
			p.advance()
		}
	}
	p.expect(')')
	return opts, nil
}

// parseCreateServerStmt parses a CREATE SERVER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-server.html
//
//	CREATE SERVER server_name
//	  FOREIGN DATA WRAPPER wrapper_name
//	  OPTIONS (option [, option] ...)
func (p *Parser) parseCreateServerStmt(start int) (*nodes.CreateServerStmt, error) {
	p.advance() // consume SERVER

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// FOREIGN DATA WRAPPER wrapper_name
	p.match(kwFOREIGN)
	p.match(kwDATA)
	p.match(kwWRAPPER)

	wrapperName, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.CreateServerStmt{
		Loc:         nodes.Loc{Start: start},
		Name:        name,
		WrapperName: wrapperName,
	}

	// OPTIONS (...)
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseServerOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterServerStmt parses an ALTER SERVER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-server.html
//
//	ALTER SERVER server_name OPTIONS (option [, option] ...)
func (p *Parser) parseAlterServerStmt(start int) (*nodes.AlterServerStmt, error) {
	p.advance() // consume SERVER

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.AlterServerStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// OPTIONS (...)
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseServerOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropServerStmt parses a DROP SERVER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-server.html
//
//	DROP SERVER [IF EXISTS] server_name
func (p *Parser) parseDropServerStmt(start int) (*nodes.DropServerStmt, error) {
	p.advance() // consume SERVER

	stmt := &nodes.DropServerStmt{Loc: nodes.Loc{Start: start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCreateSpatialRefSysStmt parses a CREATE SPATIAL REFERENCE SYSTEM statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-spatial-reference-system.html
//
//	CREATE OR REPLACE SPATIAL REFERENCE SYSTEM srid srs_attribute ...
//	CREATE SPATIAL REFERENCE SYSTEM [IF NOT EXISTS] srid srs_attribute ...
//
//	srs_attribute: {
//	  NAME 'srs_name'
//	| DEFINITION 'definition'
//	| ORGANIZATION 'org_name' IDENTIFIED BY srid
//	| DESCRIPTION 'description'
//	}
func (p *Parser) parseCreateSpatialRefSysStmt(start int, orReplace bool) (*nodes.CreateSpatialRefSysStmt, error) {
	// SPATIAL already consumed; consume REFERENCE SYSTEM
	// Current token should be REFERENCE (as ident)
	if p.cur.Type != kwREFERENCE {
		return nil, &ParseError{
			Message:  "expected REFERENCE after SPATIAL",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume REFERENCE

	if p.cur.Type != kwSYSTEM {
		return nil, &ParseError{
			Message:  "expected SYSTEM after REFERENCE",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume SYSTEM

	stmt := &nodes.CreateSpatialRefSysStmt{
		Loc:       nodes.Loc{Start: start},
		OrReplace: orReplace,
	}

	// Optional IF NOT EXISTS (only when not OR REPLACE)
	if !orReplace && p.cur.Type == kwIF {
		p.advance()
		p.match(kwNOT)
		p.match(kwEXISTS_KW)
		stmt.IfNotExists = true
	}

	// SRID (integer)
	if p.cur.Type != tokICONST {
		return nil, &ParseError{
			Message:  "expected integer SRID",
			Position: p.cur.Loc,
		}
	}
	stmt.SRID = p.cur.Ival
	p.advance()

	// SRS attributes
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		if !p.isIdentToken() {
			break
		}
		switch {
		case p.cur.Type == kwNAME:
			p.advance()
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for NAME",
					Position: p.cur.Loc,
				}
			}
			stmt.Name = p.cur.Str
			p.advance()
		case p.cur.Type == kwDEFINITION:
			p.advance()
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for DEFINITION",
					Position: p.cur.Loc,
				}
			}
			stmt.Definition = p.cur.Str
			p.advance()
		case p.cur.Type == kwORGANIZATION:
			p.advance()
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for ORGANIZATION",
					Position: p.cur.Loc,
				}
			}
			stmt.Organization = p.cur.Str
			p.advance()
			// IDENTIFIED BY srid
			p.match(kwIDENTIFIED)
			p.match(kwBY)
			if p.cur.Type != tokICONST {
				return nil, &ParseError{
					Message:  "expected integer for ORGANIZATION IDENTIFIED BY srid",
					Position: p.cur.Loc,
				}
			}
			stmt.OrgSRID = p.cur.Ival
			p.advance()
		case p.cur.Type == kwDESCRIPTION:
			p.advance()
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for DESCRIPTION",
					Position: p.cur.Loc,
				}
			}
			stmt.Description = p.cur.Str
			p.advance()
		default:
			goto done_srs
		}
	}
done_srs:
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropSpatialRefSysStmt parses a DROP SPATIAL REFERENCE SYSTEM statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-spatial-reference-system.html
//
//	DROP SPATIAL REFERENCE SYSTEM [IF EXISTS] srid
func (p *Parser) parseDropSpatialRefSysStmt(start int) (*nodes.DropSpatialRefSysStmt, error) {
	// SPATIAL already consumed; consume REFERENCE SYSTEM
	if p.cur.Type != kwREFERENCE {
		return nil, &ParseError{
			Message:  "expected REFERENCE after SPATIAL",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume REFERENCE

	if p.cur.Type != kwSYSTEM {
		return nil, &ParseError{
			Message:  "expected SYSTEM after REFERENCE",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume SYSTEM

	stmt := &nodes.DropSpatialRefSysStmt{Loc: nodes.Loc{Start: start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS_KW)
		stmt.IfExists = true
	}

	// SRID (integer)
	if p.cur.Type != tokICONST {
		return nil, &ParseError{
			Message:  "expected integer SRID",
			Position: p.cur.Loc,
		}
	}
	stmt.SRID = p.cur.Ival
	p.advance()

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseVCPUSpecs parses VCPU [=] vcpu_spec [, vcpu_spec] ...
// where vcpu_spec is N or N-M.
func (p *Parser) parseVCPUSpecs() ([]nodes.VCPUSpec, error) {
	// consume VCPU keyword (as ident)
	p.advance()

	// Optional '='
	if p.cur.Type == '=' {
		p.advance()
	}

	var specs []nodes.VCPUSpec
	for {
		if p.cur.Type != tokICONST {
			break
		}
		start := p.cur.Ival
		p.advance()

		spec := nodes.VCPUSpec{Start: start, End: -1}

		// Check for range: N-M (the '-' is parsed as a minus operator)
		if p.cur.Type == '-' {
			p.advance()
			if p.cur.Type != tokICONST {
				return nil, &ParseError{
					Message:  "expected integer after '-' in VCPU range",
					Position: p.cur.Loc,
				}
			}
			spec.End = p.cur.Ival
			p.advance()
		}

		specs = append(specs, spec)

		if p.cur.Type != ',' {
			break
		}
		// Peek ahead: if after comma we don't see a number, it's not a VCPU continuation
		// but could be a comma in some other context. For safety, always continue.
		p.advance()
	}

	return specs, nil
}

// parseResourceGroupOptions parses options common to CREATE/ALTER RESOURCE GROUP.
// It sets VCPU, THREAD_PRIORITY, ENABLE/DISABLE on the provided pointers.
func (p *Parser) parseResourceGroupOptions(vcpus *[]nodes.VCPUSpec, threadPriority **int64, enable **bool) error {
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		switch {
		case p.cur.Type == kwVCPU:
			specs, err := p.parseVCPUSpecs()
			if err != nil {
				return err
			}
			*vcpus = specs
		case p.cur.Type == kwTHREAD_PRIORITY:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			// Thread priority can be negative
			neg := false
			if p.cur.Type == '-' {
				neg = true
				p.advance()
			}
			if p.cur.Type != tokICONST {
				return &ParseError{
					Message:  "expected integer for THREAD_PRIORITY",
					Position: p.cur.Loc,
				}
			}
			val := p.cur.Ival
			if neg {
				val = -val
			}
			p.advance()
			*threadPriority = &val
		case p.cur.Type == kwENABLE:
			p.advance()
			t := true
			*enable = &t
		case p.cur.Type == kwDISABLE:
			p.advance()
			f := false
			*enable = &f
		default:
			return nil
		}
	}
	return nil
}

// parseCreateResourceGroupStmt parses a CREATE RESOURCE GROUP statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-resource-group.html
//
//	CREATE RESOURCE GROUP group_name
//	  TYPE = {SYSTEM | USER}
//	  [VCPU [=] vcpu_spec [, vcpu_spec] ...]
//	  [THREAD_PRIORITY [=] N]
//	  [ENABLE | DISABLE]
func (p *Parser) parseCreateResourceGroupStmt(start int) (*nodes.CreateResourceGroupStmt, error) {
	// RESOURCE already consumed by dispatch
	// Consume GROUP
	if p.cur.Type != kwGROUP {
		return nil, &ParseError{
			Message:  "expected GROUP after RESOURCE",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume GROUP

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.CreateResourceGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// TYPE = {SYSTEM | USER}
	if p.cur.Type == kwTYPE {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			typeName, _, _ := p.parseKeywordOrIdent()
			stmt.Type = typeName
		}
	}

	// Optional VCPU, THREAD_PRIORITY, ENABLE/DISABLE
	err = p.parseResourceGroupOptions(&stmt.VCPUs, &stmt.ThreadPriority, &stmt.Enable)
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterResourceGroupStmt parses an ALTER RESOURCE GROUP statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-resource-group.html
//
//	ALTER RESOURCE GROUP group_name
//	  [VCPU [=] vcpu_spec [, vcpu_spec] ...]
//	  [THREAD_PRIORITY [=] N]
//	  [ENABLE | DISABLE]
//	  [FORCE]
func (p *Parser) parseAlterResourceGroupStmt(start int) (*nodes.AlterResourceGroupStmt, error) {
	// RESOURCE already consumed by dispatch
	// Consume GROUP
	if p.cur.Type != kwGROUP {
		return nil, &ParseError{
			Message:  "expected GROUP after RESOURCE",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume GROUP

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.AlterResourceGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Optional VCPU, THREAD_PRIORITY, ENABLE/DISABLE
	err = p.parseResourceGroupOptions(&stmt.VCPUs, &stmt.ThreadPriority, &stmt.Enable)
	if err != nil {
		return nil, err
	}

	// Optional FORCE
	if p.cur.Type == kwFORCE {
		p.advance()
		stmt.Force = true
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropResourceGroupStmt parses a DROP RESOURCE GROUP statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-resource-group.html
//
//	DROP RESOURCE GROUP group_name [FORCE]
func (p *Parser) parseDropResourceGroupStmt(start int) (*nodes.DropResourceGroupStmt, error) {
	// RESOURCE already consumed by dispatch
	// Consume GROUP
	if p.cur.Type != kwGROUP {
		return nil, &ParseError{
			Message:  "expected GROUP after RESOURCE",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume GROUP

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.DropResourceGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// Optional FORCE
	if p.cur.Type == kwFORCE {
		p.advance()
		stmt.Force = true
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterInstanceStmt parses an ALTER INSTANCE statement.
// ALTER already consumed. p.cur is INSTANCE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-instance.html
//
//	ALTER INSTANCE instance_action
//
//	instance_action: {
//	    ENABLE INNODB REDO_LOG
//	  | DISABLE INNODB REDO_LOG
//	  | ROTATE INNODB MASTER KEY
//	  | ROTATE BINLOG MASTER KEY
//	  | RELOAD TLS [FOR CHANNEL {mysql_main | mysql_admin}] [NO ROLLBACK ON ERROR]
//	  | RELOAD KEYRING
//	}
func (p *Parser) parseAlterInstanceStmt(start int) (*nodes.AlterInstanceStmt, error) {
	p.advance() // consume INSTANCE

	stmt := &nodes.AlterInstanceStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Collect action words until end of statement or special clauses
	var words []string
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		// Check for FOR CHANNEL
		if p.cur.Type == kwFOR {
			p.advance() // consume FOR
			p.advance() // consume CHANNEL (identifier)
			name, _, _ := p.parseIdentifier()
			stmt.Channel = name
			continue
		}
		// Check for NO ROLLBACK ON ERROR
		if p.cur.Type == kwNO {
			p.advance() // consume NO
			p.advance() // consume ROLLBACK (identifier)
			p.advance() // consume ON
			p.advance() // consume ERROR (identifier)
			stmt.NoRollbackOnError = true
			continue
		}
		if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
			name, _, _ := p.parseKeywordOrIdent()
			words = append(words, name)
		} else {
			break
		}
	}

	stmt.Action = ""
	for i, w := range words {
		if i > 0 {
			stmt.Action += " "
		}
		stmt.Action += w
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseLockInstanceStmt parses LOCK INSTANCE FOR BACKUP.
// LOCK already consumed. p.cur is INSTANCE.
func (p *Parser) parseLockInstanceStmt(start int) (*nodes.LockInstanceStmt, error) {
	p.advance() // consume INSTANCE
	p.advance() // consume FOR
	p.advance() // consume BACKUP

	stmt := &nodes.LockInstanceStmt{
		Loc: nodes.Loc{Start: start, End: p.pos()},
	}
	return stmt, nil
}

// parseUnlockInstanceStmt parses UNLOCK INSTANCE.
// UNLOCK already consumed. p.cur is INSTANCE.
func (p *Parser) parseUnlockInstanceStmt(start int) (*nodes.UnlockInstanceStmt, error) {
	p.advance() // consume INSTANCE

	stmt := &nodes.UnlockInstanceStmt{
		Loc: nodes.Loc{Start: start, End: p.pos()},
	}
	return stmt, nil
}

// parseImportTableStmt parses IMPORT TABLE FROM sdi_file [, sdi_file] ...
// p.cur is IMPORT.
func (p *Parser) parseImportTableStmt() (*nodes.ImportTableStmt, error) {
	start := p.pos()
	p.advance() // consume IMPORT
	p.advance() // consume TABLE
	p.advance() // consume FROM

	stmt := &nodes.ImportTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	for {
		if p.cur.Type == tokSCONST {
			stmt.Files = append(stmt.Files, p.cur.Str)
			p.advance()
		} else {
			break
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseBinlogStmt parses BINLOG 'str'.
// p.cur is BINLOG.
func (p *Parser) parseBinlogStmt() (*nodes.BinlogStmt, error) {
	start := p.pos()
	p.advance() // consume BINLOG

	stmt := &nodes.BinlogStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == tokSCONST {
		stmt.Str = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseCacheIndexStmt parses CACHE INDEX statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/cache-index.html
//
//	CACHE INDEX tbl_name [[INDEX|KEY] (idx1, idx2, ...)]
//	    [PARTITION (p1, p2, ...) | PARTITION ALL]
//	    [, tbl_name ...] IN cache_name
//
// p.cur is CACHE.
func (p *Parser) parseCacheIndexStmt() (*nodes.CacheIndexStmt, error) {
	start := p.pos()
	p.advance() // consume CACHE
	p.advance() // consume INDEX

	stmt := &nodes.CacheIndexStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Parse table list
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)

		// Optional INDEX|KEY (idx1, idx2, ...)
		if p.cur.Type == kwINDEX || p.cur.Type == kwKEY {
			p.advance()
			if p.cur.Type == '(' {
				idxList, err := p.parseParenIdentList()
				if err != nil {
					return nil, err
				}
				stmt.Indexes = idxList
			}
		}

		// Optional PARTITION (p1, p2, ...) | PARTITION ALL
		if p.cur.Type == kwPARTITION {
			p.advance()
			if _, ok := p.match(kwALL); ok {
				stmt.Partitions = []string{"ALL"}
			} else {
				parts, err := p.parseParenIdentList()
				if err != nil {
					return nil, err
				}
				stmt.Partitions = parts
			}
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// IN cache_name
	if p.cur.Type == kwIN {
		p.advance()
		name, _, _ := p.parseIdentifier()
		stmt.CacheName = name
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseLoadIndexIntoCacheStmt parses LOAD INDEX INTO CACHE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/load-index.html
//
//	LOAD INDEX INTO CACHE
//	  tbl_index_list [, tbl_index_list] ...
//
//	tbl_index_list:
//	  tbl_name
//	    [PARTITION (partition_list)]
//	    [{INDEX | KEY} (index_name [, index_name] ...)]
//	    [IGNORE LEAVES]
//
//	partition_list: {
//	    partition_name [, partition_name] ...
//	  | ALL
//	}
//
// LOAD and INDEX already consumed. p.cur is INTO.
func (p *Parser) parseLoadIndexIntoCacheStmt(start int) (*nodes.LoadIndexIntoCacheStmt, error) {
	p.advance() // consume INTO
	p.advance() // consume CACHE

	stmt := &nodes.LoadIndexIntoCacheStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Parse table list
	for {
		tblStart := p.pos()
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		entry := &nodes.LoadIndexTable{
			Loc:   nodes.Loc{Start: tblStart},
			Table: ref,
		}

		// [PARTITION (partition_list | ALL)]
		if p.cur.Type == kwPARTITION {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				if p.cur.Type == kwALL {
					entry.Partitions = []string{"ALL"}
					p.advance()
				} else {
					for {
						name, _, _ := p.parseIdentifier()
						entry.Partitions = append(entry.Partitions, name)
						if p.cur.Type != ',' {
							break
						}
						p.advance()
					}
				}
				p.match(')')
			}
		}

		// [{INDEX | KEY} (index_name [, index_name] ...)]
		if p.cur.Type == kwINDEX || p.cur.Type == kwKEY {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				for {
					name, _, _ := p.parseIdentifier()
					entry.Indexes = append(entry.Indexes, name)
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				p.match(')')
			}
		}

		// [IGNORE LEAVES]
		if p.cur.Type == kwIGNORE {
			p.advance()
			if p.cur.Type == kwLEAVES {
				p.advance()
			}
			entry.IgnoreLeaves = true
		}

		entry.Loc.End = p.pos()
		stmt.Tables = append(stmt.Tables, entry)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseResetPersistStmt parses RESET PERSIST [IF EXISTS] [sys_var_name].
// RESET already consumed. p.cur is PERSIST.
func (p *Parser) parseResetPersistStmt(start int) (*nodes.ResetPersistStmt, error) {
	p.advance() // consume PERSIST

	stmt := &nodes.ResetPersistStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == kwIF {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}

	if p.isLvalueIdentToken() {
		name, _, _ := p.parseLvalueIdent()
		stmt.Variable = name
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseHelpStmt parses HELP 'topic'.
// p.cur is HELP.
func (p *Parser) parseHelpStmt() (*nodes.HelpStmt, error) {
	start := p.pos()
	p.advance() // consume HELP

	stmt := &nodes.HelpStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == tokSCONST {
		stmt.Topic = p.cur.Str
		p.advance()
	} else if p.cur.Type == tokIDENT || p.cur.Type >= 700 {
		// HELP topic should accept any keyword (e.g., HELP SELECT, HELP CREATE).
		name, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return nil, err
		}
		stmt.Topic = name
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSizeValue reads a size value which may be a number optionally followed
// by a suffix like K, M, G (e.g. "16M", "256K", "1G"). The lexer tokenizes
// "16M" as tokICONST("16") + tokIDENT("M"), so we concatenate them.
func (p *Parser) parseSizeValue() string {
	val := p.cur.Str
	p.advance()
	// Check for size suffix: K, M, G, T
	if p.cur.Type == tokIDENT && len(p.cur.Str) == 1 {
		ch := p.cur.Str[0]
		if ch == 'K' || ch == 'k' || ch == 'M' || ch == 'm' || ch == 'G' || ch == 'g' || ch == 'T' || ch == 't' {
			val += p.cur.Str
			p.advance()
		}
	}
	return val
}

// parseCreateLogfileGroupStmt parses a CREATE LOGFILE GROUP statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-logfile-group.html
//
//	CREATE LOGFILE GROUP logfile_group
//	    ADD UNDOFILE 'undo_file'
//	    [INITIAL_SIZE [=] initial_size]
//	    [UNDO_BUFFER_SIZE [=] undo_buffer_size]
//	    [REDO_BUFFER_SIZE [=] redo_buffer_size]
//	    [NODEGROUP [=] nodegroup_id]
//	    [WAIT]
//	    [COMMENT [=] 'string']
//	    ENGINE [=] engine_name
func (p *Parser) parseCreateLogfileGroupStmt(start int) (*nodes.CreateLogfileGroupStmt, error) {
	p.advance() // consume LOGFILE
	p.match(kwGROUP)

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.CreateLogfileGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// ADD UNDOFILE 'file_name'
	if p.cur.Type == kwADD {
		p.advance() // consume ADD
		if p.cur.Type == kwUNDOFILE {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.UndoFile = p.cur.Str
			p.advance()
		}
	}

	// Parse optional clauses
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		switch {
		case p.cur.Type == kwINITIAL_SIZE:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			stmt.InitialSize = p.parseSizeValue()
		case p.cur.Type == kwUNDO_BUFFER_SIZE:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			stmt.UndoBufferSize = p.parseSizeValue()
		case p.cur.Type == kwREDO_BUFFER_SIZE:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			stmt.RedoBufferSize = p.parseSizeValue()
		case p.cur.Type == kwNODEGROUP:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			stmt.NodeGroup = p.cur.Str
			p.advance()
		case p.cur.Type == kwWAIT:
			p.advance()
			stmt.Wait = true
		case p.cur.Type == kwCOMMENT:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokSCONST {
				stmt.Comment = p.cur.Str
				p.advance()
			}
		case p.cur.Type == kwENGINE:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			ename, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Engine = ename
		default:
			goto createLGDone
		}
	}
createLGDone:
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterLogfileGroupStmt parses an ALTER LOGFILE GROUP statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-logfile-group.html
//
//	ALTER LOGFILE GROUP logfile_group
//	    ADD UNDOFILE 'file_name'
//	    [INITIAL_SIZE [=] size]
//	    [WAIT]
//	    ENGINE [=] engine_name
func (p *Parser) parseAlterLogfileGroupStmt(start int) (*nodes.AlterLogfileGroupStmt, error) {
	p.advance() // consume LOGFILE
	p.match(kwGROUP)

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.AlterLogfileGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// ADD UNDOFILE 'file_name'
	if p.cur.Type == kwADD {
		p.advance() // consume ADD
		if p.cur.Type == kwUNDOFILE {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			stmt.UndoFile = p.cur.Str
			p.advance()
		}
	}

	// Parse optional clauses
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		switch {
		case p.cur.Type == kwINITIAL_SIZE:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			stmt.InitialSize = p.parseSizeValue()
		case p.cur.Type == kwWAIT:
			p.advance()
			stmt.Wait = true
		case p.cur.Type == kwENGINE:
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			ename, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Engine = ename
		default:
			goto alterLGDone
		}
	}
alterLGDone:
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropLogfileGroupStmt parses a DROP LOGFILE GROUP statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-logfile-group.html
//
//	DROP LOGFILE GROUP logfile_group
//	    ENGINE [=] engine_name
func (p *Parser) parseDropLogfileGroupStmt(start int) (*nodes.DropLogfileGroupStmt, error) {
	p.advance() // consume LOGFILE
	p.match(kwGROUP)

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.DropLogfileGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// ENGINE [=] engine_name
	if p.cur.Type == kwENGINE {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		ename, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Engine = ename
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
