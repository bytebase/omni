package parser

import (
	"encoding/json"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
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

	// TABLE is required. It used to be matched leniently, which left
	// `ANALYZE t` -- a 1064 on MySQL 8.0 -- parsing as an ANALYZE of t.
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

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

	// opt_histogram: UPDATE HISTOGRAM ON ident_string_list opt_histogram_update_param
	// | DROP HISTOGRAM ON ident_string_list, one of the two, every word required.
	if p.cur.Type == kwUPDATE || p.cur.Type == kwDROP {
		stmt.HistogramOp = "UPDATE"
		if p.cur.Type == kwDROP {
			stmt.HistogramOp = "DROP"
		}
		p.advance()
		for _, kw := range []int{kwHISTOGRAM, kwON} {
			if _, err := p.expect(kw); err != nil {
				return nil, err
			}
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
		if stmt.HistogramOp == "UPDATE" {
			if err := p.parseHistogramUpdateParam(stmt); err != nil {
				return nil, err
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseHistogramUpdateParam parses what may follow UPDATE HISTOGRAM's columns: nothing,
// WITH n BUCKETS, or USING DATA 'json'.
//
// Ref: opt_histogram_update_param in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseHistogramUpdateParam(stmt *nodes.AnalyzeTableStmt) error {
	switch p.cur.Type {
	case kwWITH:
		p.advance()
		if p.cur.Type != tokICONST {
			return p.syntaxErrorAtCur()
		}
		// The count is NUM, held to 1..1024 while the statement parses (1690).
		if p.cur.Ival < 1 || p.cur.Ival > 1024 {
			return &ParseError{Message: "Number of buckets value is out of range in 'ANALYZE TABLE'", Position: p.cur.Loc}
		}
		stmt.Buckets = int(p.cur.Ival)
		p.advance()
		if _, err := p.expect(kwBUCKETS); err != nil {
			return err
		}
	case kwUSING:
		p.advance()
		if _, err := p.expect(kwDATA); err != nil {
			return err
		}
		if p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		stmt.UsingData = p.cur.Str
		p.advance()
	}
	return nil
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

	// table_or_tables is required and is TABLE or TABLES (sql_yacc.yy): MySQL 8.0
	// runs `OPTIMIZE TABLES t` and answers 1064 to `OPTIMIZE t` with neither.
	if _, ok := p.match(kwTABLE, kwTABLES); !ok {
		return nil, p.syntaxErrorAtCur()
	}

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

	stmt.Loc.End = p.prev.End
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

	// table_or_tables is required (check_table_stmt, sql_yacc.yy): `CHECK t` is 1064.
	if _, ok := p.match(kwTABLE, kwTABLES); !ok {
		return nil, p.syntaxErrorAtCur()
	}

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
			// mi_check_type: FOR_SYM UPGRADE_SYM, the one check type of two words.
			p.advance()
			if _, err := p.expect(kwUPGRADE); err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, "FOR UPGRADE")
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

	stmt.Loc.End = p.prev.End
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

	// table_or_tables is required (repair_table_stmt, sql_yacc.yy): `REPAIR t` is 1064.
	if _, ok := p.match(kwTABLE, kwTABLES); !ok {
		return nil, p.syntaxErrorAtCur()
	}

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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFlushStmt parses a FLUSH statement.
//
//	FLUSH [NO_WRITE_TO_BINLOG | LOCAL] {
//	    flush_option [, flush_option] ...
//	  | {TABLE | TABLES} [tbl_name [, tbl_name] ...] [WITH READ LOCK]
//	  | {TABLE | TABLES} tbl_name [, tbl_name] ... FOR EXPORT
//	}
//
//	flush_option:
//	    {BINARY | ENGINE | ERROR | GENERAL | SLOW} LOGS
//	  | RELAY LOGS [FOR CHANNEL 'channel']
//	  | HOSTS | LOGS | PRIVILEGES | OPTIMIZER_COSTS | STATUS | USER_RESOURCES
//
// The option words are keywords, so a backquoted `hosts` is a 1064, and so is any other
// word: `FLUSH foo` and TiDB's `FLUSH TIDB PLUGINS`. The channel is a plain string, not a
// name, and may not hold a line feed, which the server refuses while it parses with 1525.
// FOR EXPORT needs a table, and `FLUSH TABLES FOR EXPORT` is a 1064 ("No tables used").
//
// Ref: flush, flush_options, opt_flush_lock and flush_option in sql_yacc.yy, and
// TEXT_STRING_sys_nonewline (mysql-8.0.46).
// Ref: https://dev.mysql.com/doc/refman/8.0/en/flush.html
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

	if p.cur.Type == kwLOCAL || p.cur.Type == kwNO_WRITE_TO_BINLOG {
		stmt.NoWriteToBinlog = true
		p.advance()
	}

	if p.cur.Type == kwTABLES || p.cur.Type == kwTABLE {
		stmt.Options = append(stmt.Options, "TABLES")
		p.advance()
		// WITH and FOR are reserved, so neither starts a table name.
		if p.isIdentToken() {
			for {
				ref, err := p.parseNamedTableRef()
				if err != nil {
					return nil, err
				}
				stmt.Tables = append(stmt.Tables, ref)
				if _, ok := p.match(','); !ok {
					break
				}
			}
			if err := refuseRepeatedTables(stmt.Tables); err != nil {
				return nil, err
			}
		}
		switch p.cur.Type {
		case kwWITH:
			p.advance()
			if _, err := p.expect(kwREAD); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwLOCK); err != nil {
				return nil, err
			}
			stmt.WithReadLock = true
		case kwFOR:
			if len(stmt.Tables) == 0 {
				return nil, &ParseError{Message: "No tables used", Position: p.cur.Loc}
			}
			p.advance()
			if _, err := p.expect(kwEXPORT); err != nil {
				return nil, err
			}
			stmt.ForExport = true
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	for {
		opt, err := p.parseFlushOption(stmt)
		if err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFlushOption parses one flush_option and returns its words, upper-cased and joined
// by a space.
func (p *Parser) parseFlushOption(stmt *nodes.FlushStmt) (string, error) {
	word := p.flushWord()
	switch word {
	case "HOSTS", "LOGS", "PRIVILEGES", "OPTIMIZER_COSTS", "STATUS", "USER_RESOURCES":
		p.advance()
		return word, nil
	case "BINARY", "ENGINE", "ERROR", "GENERAL", "SLOW":
		p.advance()
		if _, err := p.expect(kwLOGS); err != nil {
			return "", err
		}
		return word + " LOGS", nil
	case "RELAY":
		p.advance()
		if _, err := p.expect(kwLOGS); err != nil {
			return "", err
		}
		if _, ok := p.match(kwFOR); ok {
			if _, err := p.expect(kwCHANNEL); err != nil {
				return "", err
			}
			if p.cur.Type != tokSCONST {
				return "", p.syntaxErrorAtCur()
			}
			if strings.Contains(p.cur.Str, "\n") {
				return "", &ParseError{
					Message:  "Incorrect argument contains not-allowed LF value: '" + p.cur.Str + "'",
					Position: p.cur.Loc,
				}
			}
			stmt.RelayChannel = p.advance().Str
		}
		return "RELAY LOGS", nil
	}
	return "", p.syntaxErrorAtCur()
}

// flushWord returns the word at p.cur upper-cased when it is written bare, as every
// flush_option word must be, and "" otherwise.
func (p *Parser) flushWord() string {
	if (p.cur.Type != tokIDENT && p.cur.Type < 700) || p.backtickQuotedAt(p.cur.Loc) {
		return ""
	}
	return strings.ToUpper(p.cur.Str)
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

	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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

	name, err := p.parseNamedRoutineRef()
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

	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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

	// handler_scan_function (FIRST or NEXT) with no index, or an index and
	// handler_rkey_function: FIRST, NEXT, PREV, LAST, or a comparison and its
	// parenthesized values.
	//
	// Ref: handler_stmt, handler_scan_function, handler_rkey_function and
	// handler_rkey_mode in MySQL 8.0 sql/sql_yacc.yy.
	switch p.cur.Type {
	case kwFIRST, kwNEXT:
		stmt.Direction = strings.ToUpper(p.advance().Str)
	default:
		if !p.isIdentToken() {
			return nil, p.syntaxErrorAtCur()
		}
		idx, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Index = idx
		switch p.cur.Type {
		case kwFIRST, kwNEXT, kwPREV, kwLAST:
			stmt.Direction = strings.ToUpper(p.advance().Str)
		case '=', '<', '>', tokLessEq, tokGreaterEq:
			stmt.KeyMode = p.advance().Str
			if stmt.KeyMode == "" {
				stmt.KeyMode = string(rune(p.prev.Type))
			}
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			for {
				v, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if v == nil {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.KeyValues = append(stmt.KeyValues, v)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
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

	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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
//
// The list is a select list's, so an item may be a wildcard, and the bare `*`
// its first item only: MySQL 8.0 parses `DO *` and `DO t.*`, refusing them past
// the parser with 1096 and 1051, and answers 1064 for `DO 1, *`. An expression
// takes a select item's alias too (`DO 1 a`, `DO 1 AS 's'`), which the server
// parses and drops.
//
// Ref: do_stmt, `DO_SYM select_item_list`, and select_item, `expr
// select_alias`, in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseDoStmt() (*nodes.DoStmt, error) {
	start := p.pos()
	p.advance() // consume DO

	stmt := &nodes.DoStmt{Loc: nodes.Loc{Start: start}}

	for {
		if len(stmt.Exprs) > 0 && p.cur.Type == '*' {
			return nil, p.syntaxErrorAtCur()
		}
		expr, err := p.parseSelectExpr()
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

	stmt.Loc.End = p.prev.End
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

	// table_or_tables is required (checksum, sql_yacc.yy): `CHECKSUM t` is 1064.
	if _, ok := p.match(kwTABLE, kwTABLES); !ok {
		return nil, p.syntaxErrorAtCur()
	}

	stmt := &nodes.ChecksumTableStmt{Loc: nodes.Loc{Start: start}}

	for {
		ref, err := p.parseNamedTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := refuseRepeatedTables(stmt.Tables); err != nil {
		return nil, err
	}

	// Optional QUICK | EXTENDED
	if p.cur.Type == kwQUICK {
		stmt.Quick = true
		p.advance()
	} else if p.cur.Type == kwEXTENDED {
		stmt.Extended = true
		p.advance()
	}

	stmt.Loc.End = p.prev.End
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
	stmt.Loc.End = p.prev.End
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
	stmt.Loc.End = p.prev.End
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
		// clone_stmt: CLONE_SYM LOCAL_SYM DATA_SYM DIRECTORY_SYM opt_equal, both words required.
		for _, kw := range []int{kwDATA, kwDIRECTORY} {
			if _, err := p.expect(kw); err != nil {
				return nil, err
			}
		}

		// Optional '='
		p.matchEqual()

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
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}

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
		for _, kw := range []int{kwIDENTIFIED, kwBY} {
			if _, err := p.expect(kw); err != nil {
				return nil, err
			}
		}
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
			if _, err := p.expect(kwDIRECTORY); err != nil {
				return nil, err
			}
			p.matchEqual()
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
			// opt_ssl: REQUIRE_SYM [NO_SYM] SSL_SYM, SSL required.
			required := true
			if p.cur.Type == kwNO {
				p.advance() // consume NO
				required = false
			}
			if _, err := p.expect(kwSSL); err != nil {
				return nil, err
			}
			stmt.RequireSSL = &required
		}
	} else {
		return nil, &ParseError{
			Message:  "expected LOCAL or INSTANCE after CLONE",
			Position: p.cur.Loc,
		}
	}

	stmt.Loc.End = p.prev.End
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

	// SONAME 'library', whole (install in MySQL 8.0 sql/sql_yacc.yy).
	if _, err := p.expect(kwSONAME); err != nil {
		return nil, err
	}
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	soname := p.advance().Str

	stmt := &nodes.InstallPluginStmt{
		Loc:        nodes.Loc{Start: start},
		PluginName: name,
		Soname:     soname,
	}
	stmt.Loc.End = p.prev.End
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
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Components = append(stmt.Components, p.advance().Str)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// Optional SET clause — INSTALL COMPONENT ... SET targets are
	// component-system variables, never stored-program locals.
	if p.cur.Type == kwSET {
		p.advance() // consume SET
		for {
			// install_option_type: an assignment may name GLOBAL or PERSIST.
			scope := ""
			if p.cur.Type == kwGLOBAL || p.cur.Type == kwPERSIST {
				scope = strings.ToUpper(p.advance().Str)
			}
			asgn, err := p.parseSetAssignment()
			if err != nil {
				return nil, err
			}
			stmt.SetVars = append(stmt.SetVars, asgn)
			stmt.SetScopes = append(stmt.SetScopes, scope)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
	}

	stmt.Loc.End = p.prev.End
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
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseUninstallComponentStmt parses an UNINSTALL COMPONENT statement.
//
//	UNINSTALL COMPONENT component_name [, component_name] ...
func (p *Parser) parseUninstallComponentStmt(start int) (*nodes.UninstallComponentStmt, error) {
	p.advance() // consume COMPONENT

	stmt := &nodes.UninstallComponentStmt{Loc: nodes.Loc{Start: start}}

	for {
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Components = append(stmt.Components, p.advance().Str)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// tsOption is one tablespace option the grammar names, read by
// parseTablespaceOption.
type tsOption struct {
	kw  int
	val string
}

// Which options each statement takes.
//
// Ref: tablespace_option, alter_tablespace_option, undo_tablespace_option and
// drop_ts_option in MySQL 8.0 sql/sql_yacc.yy.
var (
	tsCreateOptions = map[int]bool{kwINITIAL_SIZE: true, kwAUTOEXTEND_SIZE: true, kwMAX_SIZE: true,
		kwEXTENT_SIZE: true, kwNODEGROUP: true, kwENGINE: true, kwWAIT: true, kwNO_WAIT: true,
		kwCOMMENT: true, kwFILE_BLOCK_SIZE: true, kwENCRYPTION: true, kwENGINE_ATTRIBUTE: true}
	tsAlterOptions = map[int]bool{kwINITIAL_SIZE: true, kwAUTOEXTEND_SIZE: true, kwMAX_SIZE: true,
		kwENGINE: true, kwWAIT: true, kwNO_WAIT: true, kwENCRYPTION: true, kwENGINE_ATTRIBUTE: true}
	tsUndoOptions = map[int]bool{kwENGINE: true}
	tsDropOptions = map[int]bool{kwENGINE: true, kwWAIT: true, kwNO_WAIT: true}
)

// parseTablespaceOptions reads a tablespace option list: options from
// allowed, separated by an optional comma (opt_comma). With required the list
// may not be empty.
func (p *Parser) parseTablespaceOptions(allowed map[int]bool, required bool) ([]tsOption, error) {
	var opts []tsOption
	for {
		if !allowed[p.cur.Type] {
			if len(opts) > 0 && p.prev.Type == ',' || len(opts) == 0 && required {
				return nil, p.syntaxErrorAtCur()
			}
			return opts, nil
		}
		at := p.cur
		opt, err := p.parseTablespaceOption()
		if err != nil {
			return nil, err
		}
		// ENGINE comes once (ER_FILEGROUP_OPTION_ONLY_ONCE); WAIT may repeat.
		if opt.kw == kwENGINE {
			for _, o := range opts {
				if o.kw == kwENGINE {
					return nil, &ParseError{Message: "It is not allowed to specify STORAGE ENGINE more than once", Position: at.Loc}
				}
			}
		}
		opts = append(opts, opt)
		p.match(',')
	}
}

// parseTablespaceOption reads one option: a size, a number, a name or a
// string after an optional =, or WAIT or NO_WAIT alone. ENGINE_ATTRIBUTE is a
// JSON text, or empty (json_attribute's ER_INVALID_JSON_ATTRIBUTE).
//
// Ref: the ts_option_* rules and json_attribute in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseTablespaceOption() (tsOption, error) {
	opt := tsOption{kw: p.cur.Type}
	p.advance()
	switch opt.kw {
	case kwWAIT, kwNO_WAIT:
		return opt, nil
	}
	p.matchEqual()
	var err error
	switch opt.kw {
	case kwINITIAL_SIZE, kwAUTOEXTEND_SIZE, kwMAX_SIZE, kwEXTENT_SIZE, kwFILE_BLOCK_SIZE:
		opt.val, err = p.consumeSizeOptionValue()
	case kwNODEGROUP:
		opt.val, err = p.consumeIntegerOptionValue()
	case kwENGINE:
		opt.val, _, err = p.parseIdentOrText()
	case kwCOMMENT, kwENCRYPTION, kwENGINE_ATTRIBUTE:
		if p.cur.Type != tokSCONST {
			return opt, p.syntaxErrorAtCur()
		}
		if opt.kw == kwENGINE_ATTRIBUTE && p.cur.Str != "" && !json.Valid([]byte(p.cur.Str)) {
			return opt, &ParseError{Message: "Invalid json attribute", Position: p.cur.Loc}
		}
		opt.val = p.advance().Str
	}
	return opt, err
}

// parseTsDatafile reads ts_datafile, `DATAFILE 'name'`, after ADD or DROP.
func (p *Parser) parseTsDatafile() (string, error) {
	if _, err := p.expect(kwDATAFILE); err != nil {
		return "", err
	}
	if p.cur.Type != tokSCONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// parseCreateTablespaceStmt parses CREATE [UNDO] TABLESPACE. CREATE [UNDO]
// have been consumed; p.cur is TABLESPACE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-tablespace.html
// Ref: create's `CREATE TABLESPACE_SYM ident opt_ts_datafile_name
// opt_logfile_group_name opt_tablespace_options` and `CREATE UNDO_SYM
// TABLESPACE_SYM ident ADD ts_datafile opt_undo_tablespace_options` in MySQL
// 8.0 sql/sql_yacc.yy.
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

	if _, ok := p.match(kwADD); ok {
		if stmt.DataFile, err = p.parseTsDatafile(); err != nil {
			return nil, err
		}
	} else if undo {
		return nil, p.syntaxErrorAtCur()
	}
	allowed := tsUndoOptions
	if !undo {
		allowed = tsCreateOptions
		if _, ok := p.match(kwUSE); ok {
			for _, kw := range []int{kwLOGFILE, kwGROUP} {
				if _, err := p.expect(kw); err != nil {
					return nil, err
				}
			}
			if stmt.UseLogfileGroup, _, err = p.parseIdentifier(); err != nil {
				return nil, err
			}
		}
	}
	opts, err := p.parseTablespaceOptions(allowed, false)
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		switch o.kw {
		case kwINITIAL_SIZE:
			stmt.InitialSize = o.val
		case kwAUTOEXTEND_SIZE:
			stmt.AutoextendSize = o.val
		case kwMAX_SIZE:
			stmt.MaxSize = o.val
		case kwEXTENT_SIZE:
			stmt.ExtentSize = o.val
		case kwNODEGROUP:
			stmt.NodeGroup = o.val
		case kwENGINE:
			stmt.Engine = o.val
		case kwWAIT:
			stmt.Wait = true
		case kwNO_WAIT:
			stmt.NoWait = true
		case kwCOMMENT:
			stmt.Comment = o.val
		case kwFILE_BLOCK_SIZE:
			stmt.FileBlockSize = o.val
		case kwENCRYPTION:
			stmt.Encryption = o.val
		case kwENGINE_ATTRIBUTE:
			stmt.EngineAttribute = o.val
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterTablespaceStmt parses ALTER [UNDO] TABLESPACE. ALTER [UNDO] have
// been consumed; p.cur is TABLESPACE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-tablespace.html
// Ref: alter_tablespace_stmt, `ALTER TABLESPACE ident {ADD | DROP}
// ts_datafile opt_alter_tablespace_options | RENAME TO ident |
// alter_tablespace_option_list`, and alter_undo_tablespace_stmt, `ALTER UNDO
// TABLESPACE ident SET undo_tablespace_state opt_undo_tablespace_options`, in
// MySQL 8.0 sql/sql_yacc.yy.
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

	var opts []tsOption
	switch {
	case undo:
		if _, err := p.expect(kwSET); err != nil {
			return nil, err
		}
		switch p.cur.Type {
		case kwACTIVE:
			stmt.SetActive = true
		case kwINACTIVE:
			stmt.SetInactive = true
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		opts, err = p.parseTablespaceOptions(tsUndoOptions, false)
	case p.cur.Type == kwADD:
		p.advance()
		if stmt.AddDataFile, err = p.parseTsDatafile(); err != nil {
			return nil, err
		}
		opts, err = p.parseTablespaceOptions(tsAlterOptions, false)
	case p.cur.Type == kwDROP:
		p.advance()
		if stmt.DropDataFile, err = p.parseTsDatafile(); err != nil {
			return nil, err
		}
		opts, err = p.parseTablespaceOptions(tsAlterOptions, false)
	case p.cur.Type == kwRENAME:
		p.advance()
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		stmt.RenameTo, _, err = p.parseIdentifier()
	default:
		opts, err = p.parseTablespaceOptions(tsAlterOptions, true)
	}
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		switch o.kw {
		case kwINITIAL_SIZE:
			stmt.InitialSize = o.val
		case kwAUTOEXTEND_SIZE:
			stmt.AutoextendSize = o.val
		case kwMAX_SIZE:
			stmt.MaxSize = o.val
		case kwENGINE:
			stmt.Engine = o.val
		case kwWAIT:
			stmt.Wait = true
		case kwNO_WAIT:
			stmt.NoWait = true
		case kwENCRYPTION:
			stmt.Encryption = o.val
		case kwENGINE_ATTRIBUTE:
			stmt.EngineAttribute = o.val
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropTablespaceStmt parses DROP [UNDO] TABLESPACE. DROP [UNDO] have
// been consumed; p.cur is TABLESPACE.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/drop-tablespace.html
// Ref: drop_tablespace_stmt, `DROP TABLESPACE ident opt_drop_ts_options`,
// and drop_undo_tablespace_stmt, `DROP UNDO TABLESPACE ident
// opt_undo_tablespace_options`, in MySQL 8.0 sql/sql_yacc.yy.
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

	allowed := tsDropOptions
	if undo {
		allowed = tsUndoOptions
	}
	opts, err := p.parseTablespaceOptions(allowed, false)
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		switch o.kw {
		case kwENGINE:
			stmt.Engine = o.val
		case kwWAIT:
			stmt.Wait = true
		case kwNO_WAIT:
			stmt.NoWait = true
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseServerOptions reads `OPTIONS ( server_options_list )`, which CREATE
// SERVER and ALTER SERVER both require. Each option is a word and its value:
// a string for USER, HOST, DATABASE, OWNER, PASSWORD and SOCKET, and
// ulong_num for PORT. An option comes back as its word and its value, joined
// by a space.
//
// Ref: server_options_list and server_option in MySQL 8.0 sql/sql_yacc.yy.
func (p *Parser) parseServerOptions() ([]string, error) {
	if _, err := p.expect(kwOPTIONS); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []string
	for {
		word := ""
		switch p.cur.Type {
		case kwUSER, kwHOST, kwDATABASE, kwOWNER, kwPASSWORD, kwSOCKET, kwPORT:
			word = strings.ToUpper(p.cur.Str)
		default:
			return nil, p.syntaxErrorAtCur()
		}
		isPort := p.cur.Type == kwPORT
		p.advance()
		var val string
		switch {
		case isPort && (p.cur.Type == tokICONST || p.cur.Type == tokFCONST ||
			(p.cur.Type == tokXCONST && strings.HasPrefix(strings.ToLower(p.cur.Str), "0x"))):
			val = p.advance().Str
		case !isPort && p.cur.Type == tokSCONST:
			val = p.advance().Str
		default:
			return nil, p.syntaxErrorAtCur()
		}
		opts = append(opts, word+" "+val)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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

	name, _, err := p.parseIdentOrText()
	if err != nil {
		return nil, err
	}

	// FOREIGN DATA WRAPPER wrapper_name
	for _, kw := range []int{kwFOREIGN, kwDATA, kwWRAPPER} {
		if _, err := p.expect(kw); err != nil {
			return nil, err
		}
	}

	wrapperName, _, err := p.parseIdentOrText()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.CreateServerStmt{
		Loc:         nodes.Loc{Start: start},
		Name:        name,
		WrapperName: wrapperName,
	}

	// OPTIONS (...), which the grammar requires.
	opts, err := p.parseServerOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterServerStmt parses an ALTER SERVER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-server.html
//
//	ALTER SERVER server_name OPTIONS (option [, option] ...)
func (p *Parser) parseAlterServerStmt(start int) (*nodes.AlterServerStmt, error) {
	p.advance() // consume SERVER

	name, _, err := p.parseIdentOrText()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.AlterServerStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// OPTIONS (...), which the grammar requires.
	opts, err := p.parseServerOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.End
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

	name, _, err := p.parseIdentOrText()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.End
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
		for _, kw := range []int{kwNOT, kwEXISTS_KW} {
			if _, err := p.expect(kw); err != nil {
				return nil, err
			}
		}
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
			// IDENTIFIED BY srid, both words required (srs_attributes, sql_yacc.yy).
			for _, kw := range []int{kwIDENTIFIED, kwBY} {
				if _, err := p.expect(kw); err != nil {
					return nil, err
				}
			}
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
	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseVCPUSpecs parses VCPU [=] vcpu_spec [, vcpu_spec] ...
// where vcpu_spec is N or N-M.
func (p *Parser) parseVCPUSpecs() ([]nodes.VCPUSpec, error) {
	// consume VCPU keyword (as ident)
	p.advance()

	// Optional '='
	p.matchEqual()

	// A CPU number is NUM, as the thread priority is: mysql:8.0 answers 1064 for
	// VCPU = 2147483648 and for THREAD_PRIORITY = 2147483648.
	//
	// Ref: vcpu_num_or_range and signed_num in MySQL 8.0 sql/sql_yacc.yy.
	var specs []nodes.VCPUSpec
	for {
		if !p.cur.isNum() {
			if len(specs) == 0 {
				// vcpu_spec_list is not optional once VCPU is written:
				// mysql:8.0 answers 1064 for `VCPU =` with nothing after it.
				return nil, &ParseError{
					Message:  "expected a CPU number after VCPU",
					Position: p.cur.Loc,
				}
			}
			break
		}
		start := p.cur.Ival
		p.advance()

		spec := nodes.VCPUSpec{Start: start, End: -1}

		// Check for range: N-M (the '-' is parsed as a minus operator)
		if p.cur.Type == '-' {
			p.advance()
			if !p.cur.isNum() {
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
		switch p.cur.Type {
		case kwVCPU:
			specs, err := p.parseVCPUSpecs()
			if err != nil {
				return err
			}
			*vcpus = specs
		case kwTHREAD_PRIORITY:
			p.advance()
			p.matchEqual()
			// Thread priority can be negative
			neg := false
			if p.cur.Type == '-' {
				neg = true
				p.advance()
			}
			if !p.cur.isNum() {
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
		case kwENABLE, kwDISABLE:
			// One of the two, once: `ENABLE DISABLE` is 1064 on mysql:8.0,
			// because the grammar has a single opt_resource_group_enable_disable.
			if *enable != nil {
				return &ParseError{
					Message:  "ENABLE or DISABLE, not both",
					Position: p.cur.Loc,
				}
			}
			v := p.cur.Type == kwENABLE
			p.advance()
			*enable = &v
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

	// TYPE = {SYSTEM | USER}, which is required and takes those two words
	// only: mysql:8.0 answers 1064 for `CREATE RESOURCE GROUP x` and for
	// `TYPE = foo`.
	if p.cur.Type != kwTYPE {
		return nil, &ParseError{
			Message:  "expected TYPE after the resource group name",
			Position: p.cur.Loc,
		}
	}
	p.advance() // consume TYPE
	p.matchEqual()
	if p.cur.Type != kwSYSTEM && p.cur.Type != kwUSER {
		return nil, &ParseError{
			Message:  "expected SYSTEM or USER for TYPE",
			Position: p.cur.Loc,
		}
	}
	typeName, _, err := p.parseKeywordOrIdent()
	if err != nil {
		return nil, err
	}
	stmt.Type = typeName

	// Optional VCPU, THREAD_PRIORITY, ENABLE/DISABLE
	if err := p.parseResourceGroupOptions(&stmt.VCPUs, &stmt.ThreadPriority, &stmt.Enable); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
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
//
// Ref: alter_instance_action in MySQL 8.0 sql/sql_yacc.yy. ROTATE takes an
// ident_or_text, ENABLE and DISABLE two idents, which the actions hold to
// INNODB, BINLOG and REDO_LOG in any case (syntax_error_at otherwise), and
// FOR CHANNEL takes an ident.
func (p *Parser) parseAlterInstanceStmt(start int) (*nodes.AlterInstanceStmt, error) {
	p.advance() // consume INSTANCE

	stmt := &nodes.AlterInstanceStmt{Loc: nodes.Loc{Start: start}}

	// expectWord reads an ident (or, with text, an ident_or_text) that the
	// action holds to one of words.
	expectWord := func(text bool, words ...string) (string, error) {
		at := p.cur
		var name string
		var err error
		if text {
			name, _, err = p.parseIdentOrText()
		} else {
			name, _, err = p.parseIdent()
		}
		if err != nil {
			return "", err
		}
		for _, w := range words {
			if strings.EqualFold(name, w) {
				return w, nil
			}
		}
		return "", p.syntaxErrorAtTok(at)
	}

	switch p.cur.Type {
	case kwROTATE:
		p.advance()
		which, err := expectWord(true, "INNODB", "BINLOG")
		if err != nil {
			return nil, err
		}
		for _, kw := range []int{kwMASTER, kwKEY} {
			if _, err := p.expect(kw); err != nil {
				return nil, err
			}
		}
		stmt.Action = "ROTATE " + which + " MASTER KEY"
	case kwENABLE, kwDISABLE:
		verb := strings.ToUpper(p.advance().Str)
		if _, err := expectWord(false, "INNODB"); err != nil {
			return nil, err
		}
		if _, err := expectWord(false, "REDO_LOG"); err != nil {
			return nil, err
		}
		stmt.Action = verb + " INNODB REDO_LOG"
	case kwRELOAD:
		p.advance()
		switch p.cur.Type {
		case kwKEYRING:
			p.advance()
			stmt.Action = "RELOAD KEYRING"
		case kwTLS:
			p.advance()
			stmt.Action = "RELOAD TLS"
			if p.cur.Type == kwFOR {
				p.advance()
				if _, err := p.expect(kwCHANNEL); err != nil {
					return nil, err
				}
				name, _, err := p.parseIdent()
				if err != nil {
					return nil, err
				}
				stmt.Channel = name
			}
			if p.cur.Type == kwNO {
				p.advance()
				for _, kw := range []int{kwROLLBACK, kwON, kwERROR_KW} {
					if _, err := p.expect(kw); err != nil {
						return nil, err
					}
				}
				stmt.NoRollbackOnError = true
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseLockInstanceStmt parses LOCK INSTANCE FOR BACKUP.
// LOCK already consumed. p.cur is INSTANCE.
func (p *Parser) parseLockInstanceStmt(start int) (*nodes.LockInstanceStmt, error) {
	p.advance() // consume INSTANCE
	for _, kw := range []int{kwFOR, kwBACKUP} {
		if _, err := p.expect(kw); err != nil {
			return nil, err
		}
	}

	stmt := &nodes.LockInstanceStmt{
		Loc: nodes.Loc{Start: start, End: p.prev.End},
	}
	return stmt, nil
}

// parseUnlockInstanceStmt parses UNLOCK INSTANCE.
// UNLOCK already consumed. p.cur is INSTANCE.
func (p *Parser) parseUnlockInstanceStmt(start int) (*nodes.UnlockInstanceStmt, error) {
	p.advance() // consume INSTANCE

	stmt := &nodes.UnlockInstanceStmt{
		Loc: nodes.Loc{Start: start, End: p.prev.End},
	}
	return stmt, nil
}

// parseImportTableStmt parses IMPORT TABLE FROM sdi_file [, sdi_file] ...
// p.cur is IMPORT.
func (p *Parser) parseImportTableStmt() (*nodes.ImportTableStmt, error) {
	start := p.pos()
	p.advance() // consume IMPORT
	// import_stmt: IMPORT TABLE_SYM FROM TEXT_STRING_sys_list, a list of one string or more.
	for _, kw := range []int{kwTABLE, kwFROM} {
		if _, err := p.expect(kw); err != nil {
			return nil, err
		}
	}

	stmt := &nodes.ImportTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	for {
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Files = append(stmt.Files, p.cur.Str)
		p.advance()
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
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

	// binlog_base64_event: BINLOG_SYM TEXT_STRING_sys, the string required.
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Str = p.cur.Str
	p.advance()

	stmt.Loc.End = p.prev.End
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
	if _, err := p.expect(kwINDEX); err != nil {
		return nil, err
	}

	stmt := &nodes.CacheIndexStmt{
		Loc: nodes.Loc{Start: start},
	}

	// A list of tables each with an optional key list, or one table with
	// PARTITION (ALL | names) and a key list; then IN and the cache, a name or
	// DEFAULT. A key list is INDEX or KEY and parenthesized names, maybe none.
	//
	// Ref: keycache_stmt, adm_partition and opt_cache_key_list in MySQL 8.0
	// sql/sql_yacc.yy.
	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Tables = append(stmt.Tables, ref)

		if len(stmt.Tables) == 1 && p.cur.Type == kwPARTITION {
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			if _, ok := p.match(kwALL); ok {
				stmt.Partitions = []string{"ALL"}
			} else {
				for {
					name, _, err := p.parseIdentifier()
					if err != nil {
						return nil, err
					}
					stmt.Partitions = append(stmt.Partitions, name)
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}

		if p.cur.Type == kwINDEX || p.cur.Type == kwKEY {
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			stmt.Indexes = nil
			if p.cur.Type != ')' {
				for {
					if p.cur.Type == kwPRIMARY {
						stmt.Indexes = append(stmt.Indexes, "PRIMARY")
						p.advance()
					} else {
						name, _, err := p.parseIdentifier()
						if err != nil {
							return nil, err
						}
						stmt.Indexes = append(stmt.Indexes, name)
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}

		if p.cur.Type != ',' || stmt.Partitions != nil {
			break
		}
		p.advance()
	}

	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, ok := p.match(kwDEFAULT); ok {
		stmt.CacheName = "DEFAULT"
	} else {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.CacheName = name
	}

	stmt.Loc.End = p.prev.End
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
			// adm_partition: PARTITION '(' {ALL | name, ...} ')', whole.
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			if p.cur.Type == kwALL {
				entry.Partitions = []string{"ALL"}
				p.advance()
			} else {
				for {
					name, _, err := p.parseIdentifier()
					if err != nil {
						return nil, err
					}
					entry.Partitions = append(entry.Partitions, name)
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
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
			if _, err := p.expect(kwLEAVES); err != nil {
				return nil, err
			}
			entry.IgnoreLeaves = true
		}

		entry.Loc.End = p.prev.End
		stmt.Tables = append(stmt.Tables, entry)

		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseResetPersistStmt parses RESET PERSIST [IF EXISTS] [sys_var_name].
// RESET already consumed. p.cur is PERSIST.
func (p *Parser) parseResetPersistStmt(start int) (*nodes.ResetPersistStmt, error) {
	p.advance() // consume PERSIST

	stmt := &nodes.ResetPersistStmt{
		Loc: nodes.Loc{Start: start},
	}

	// RESET PERSIST [[IF EXISTS] persisted_variable_ident], the name ident,
	// ident.ident or DEFAULT.ident.
	//
	// Ref: reset and persisted_variable_ident in MySQL 8.0 sql/sql_yacc.yy.
	if p.cur.Type == kwIF {
		p.advance() // consume IF
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
		if !p.isLvalueIdentToken() && p.cur.Type != kwDEFAULT {
			return nil, p.syntaxErrorAtCur()
		}
	}

	if p.cur.Type == kwDEFAULT && p.peekNext().Type == '.' {
		p.advance()
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Variable = "DEFAULT." + name
	} else if p.isLvalueIdentToken() {
		name, _, err := p.parseLvalueIdent()
		if err != nil {
			return nil, err
		}
		stmt.Variable = name
		if p.cur.Type == '.' {
			p.advance()
			suffix, _, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Variable += "." + suffix
		}
	}

	stmt.Loc.End = p.prev.End
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.CreateLogfileGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// ADD UNDOFILE 'file_name', required and whole.
	//
	// Ref: create's `CREATE LOGFILE_SYM GROUP_SYM ident ADD lg_undofile` and
	// alter_logfile_stmt's `ALTER LOGFILE_SYM GROUP_SYM ident ADD lg_undofile`
	// in MySQL 8.0 sql/sql_yacc.yy.
	for _, kw := range []int{kwADD, kwUNDOFILE} {
		if _, err := p.expect(kw); err != nil {
			return nil, err
		}
	}
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.UndoFile = p.advance().Str

	// Parse optional clauses
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		switch {
		case p.cur.Type == kwINITIAL_SIZE:
			p.advance()
			p.matchEqual()
			size, err := p.consumeSizeOptionValue()
			if err != nil {
				return nil, err
			}
			stmt.InitialSize = size
		case p.cur.Type == kwUNDO_BUFFER_SIZE:
			p.advance()
			p.matchEqual()
			size, err := p.consumeSizeOptionValue()
			if err != nil {
				return nil, err
			}
			stmt.UndoBufferSize = size
		case p.cur.Type == kwREDO_BUFFER_SIZE:
			p.advance()
			p.matchEqual()
			size, err := p.consumeSizeOptionValue()
			if err != nil {
				return nil, err
			}
			stmt.RedoBufferSize = size
		case p.cur.Type == kwNODEGROUP:
			p.advance()
			p.matchEqual()
			group, err := p.consumeIntegerOptionValue()
			if err != nil {
				return nil, err
			}
			stmt.NodeGroup = group
		case p.cur.Type == kwWAIT:
			p.advance()
			stmt.Wait = true
		case p.cur.Type == kwCOMMENT:
			p.advance()
			p.matchEqual()
			// The comment is a quoted string and is required (1064 without it).
			if p.cur.Type != tokSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Comment = p.cur.Str
			p.advance()
		case p.cur.Type == kwENGINE:
			p.advance()
			p.matchEqual()
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
	stmt.Loc.End = p.prev.End
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
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	stmt := &nodes.AlterLogfileGroupStmt{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	// ADD UNDOFILE 'file_name', required and whole.
	//
	// Ref: create's `CREATE LOGFILE_SYM GROUP_SYM ident ADD lg_undofile` and
	// alter_logfile_stmt's `ALTER LOGFILE_SYM GROUP_SYM ident ADD lg_undofile`
	// in MySQL 8.0 sql/sql_yacc.yy.
	for _, kw := range []int{kwADD, kwUNDOFILE} {
		if _, err := p.expect(kw); err != nil {
			return nil, err
		}
	}
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.UndoFile = p.advance().Str

	// Parse optional clauses
	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		switch {
		case p.cur.Type == kwINITIAL_SIZE:
			p.advance()
			p.matchEqual()
			size, err := p.consumeSizeOptionValue()
			if err != nil {
				return nil, err
			}
			stmt.InitialSize = size
		case p.cur.Type == kwWAIT:
			p.advance()
			stmt.Wait = true
		case p.cur.Type == kwENGINE:
			p.advance()
			p.matchEqual()
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
	stmt.Loc.End = p.prev.End
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
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}

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
		p.matchEqual()
		ename, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Engine = ename
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
