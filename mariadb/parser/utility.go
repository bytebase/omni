package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
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
	// `ANALYZE t` -- a 1064 on MariaDB 11.8.8 -- parsing as an ANALYZE of t.
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

	stmt.Loc.End = p.prev.End
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

	// table_or_tables is required and is TABLE or TABLES: MariaDB 11.8 runs
	// `OPTIMIZE TABLES t` and answers 1064 to `OPTIMIZE t` with neither.
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFlushStmt parses a FLUSH statement.
//
//	FLUSH [NO_WRITE_TO_BINLOG | LOCAL] {
//	    flush_option [, flush_option] ...
//	  | {TABLE | TABLES} [tbl_name [, tbl_name] ...] [WITH READ LOCK [AND DISABLE CHECKPOINT]]
//	  | {TABLE | TABLES} tbl_name [, tbl_name] ... FOR EXPORT
//	}
//
//	flush_option:
//	    {ENGINE | ERROR | GENERAL | SLOW} LOGS
//	  | BINARY LOGS [DELETE_DOMAIN_ID = ([domain_id [, domain_id] ...])]
//	  | RELAY LOGS ['connection_name'] [FOR CHANNEL 'connection_name']
//	  | {SLAVE | REPLICA} ['connection_name']
//	  | QUERY CACHE | [SESSION | GLOBAL] STATUS
//	  | HOSTS | LOGS | PRIVILEGES | MASTER | DES_KEY_FILE | USER_RESOURCES | SSL | THREADS
//	  | information_schema_table
//
// The option words are keywords, so a backquoted `hosts` is a 1064, and so is a word
// that is none of them: `FLUSH foo`, TiDB's `FLUSH TIDB PLUGINS` and MySQL's
// OPTIMIZER_COSTS. The last form names an information_schema table that has a reset
// function, which the server looks up while it parses and refuses with 1064 when there is
// none; flushResetTables lists them. Three more checks run while parsing: a domain id
// above 2^32-1 answers 1076, a connection name and FOR CHANNEL together answer 1210, and
// RELAY LOGS or SLAVE twice in one list answers 1221. FOR EXPORT needs a table, and
// `FLUSH TABLES FOR EXPORT` is a 1064 ("No tables used").
//
// Ref: flush, flush_options, flush_lock, flush_option, opt_delete_gtid_domain and
// for_channel in sql_yacc.yy (mariadb-11.8.8).
// Ref: https://mariadb.com/kb/en/flush/
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
				ref, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				stmt.Tables = append(stmt.Tables, ref)
				if _, ok := p.match(','); !ok {
					break
				}
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
			if _, ok := p.match(kwAND); ok {
				if _, err := p.expect(kwDISABLE); err != nil {
					return nil, err
				}
				if p.flushWord() != "CHECKPOINT" {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
			}
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
		optLoc := p.cur.Loc
		opt, err := p.parseFlushOption(stmt)
		if err != nil {
			return nil, err
		}
		if opt == "RELAY LOGS" || opt == "SLAVE" {
			for _, seen := range stmt.Options {
				if seen == opt {
					return nil, &ParseError{Message: "Incorrect usage of FLUSH and " + opt, Position: optLoc}
				}
			}
		}
		stmt.Options = append(stmt.Options, opt)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// flushResetTables are the information_schema tables FLUSH can name, those with a reset
// function: the userstat plugin's four, user_variables, and the thread pool's statistics
// and waits, all loaded on mariadb:11.8.8 by default, and query_response_time's three,
// which answer once the plugin is installed. The server looks the name up among the loaded
// plugins, which a parser cannot see, so these are the ones MariaDB ships.
//
// Ref: reset_table in plugin/userstat, plugin/user_variables, plugin/query_response_time
// and sql/thread_pool_info.cc (mariadb-11.8.8).
var flushResetTables = map[string]bool{
	"CLIENT_STATISTICS": true, "INDEX_STATISTICS": true, "TABLE_STATISTICS": true,
	"USER_STATISTICS": true, "USER_VARIABLES": true,
	"THREAD_POOL_STATS": true, "THREAD_POOL_WAITS": true,
	"QUERY_RESPONSE_TIME": true, "QUERY_RESPONSE_TIME_READ": true,
	"QUERY_RESPONSE_TIME_WRITE": true,
}

// parseFlushOption parses one flush_option and returns its words, upper-cased and joined
// by a space.
func (p *Parser) parseFlushOption(stmt *nodes.FlushStmt) (string, error) {
	word := p.flushWord()
	switch word {
	case "HOSTS", "LOGS", "PRIVILEGES", "STATUS", "MASTER", "DES_KEY_FILE", "USER_RESOURCES", "SSL", "THREADS":
		p.advance()
		return word, nil
	case "ENGINE", "ERROR", "GENERAL", "SLOW":
		p.advance()
		if _, err := p.expect(kwLOGS); err != nil {
			return "", err
		}
		return word + " LOGS", nil
	case "SESSION", "GLOBAL":
		p.advance()
		if _, err := p.expect(kwSTATUS); err != nil {
			return "", err
		}
		return word + " STATUS", nil
	case "QUERY":
		p.advance()
		if p.flushWord() != "CACHE" {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		return "QUERY CACHE", nil
	case "BINARY":
		p.advance()
		if _, err := p.expect(kwLOGS); err != nil {
			return "", err
		}
		if p.flushWord() == "DELETE_DOMAIN_ID" {
			p.advance()
			if err := p.parseDeleteDomainIDs(); err != nil {
				return "", err
			}
		}
		return "BINARY LOGS", nil
	case "RELAY":
		p.advance()
		if _, err := p.expect(kwLOGS); err != nil {
			return "", err
		}
		name, hasName := p.matchConnectionName()
		if _, ok := p.match(kwFOR); ok {
			if _, err := p.expect(kwCHANNEL); err != nil {
				return "", err
			}
			if hasName {
				return "", &ParseError{
					Message:  "Incorrect arguments to CONNECTION_NAME AND FOR CHANNEL CAN NOT BE SPECIFIED AT THE SAME TIME)",
					Position: p.prev.Loc,
				}
			}
			if name, hasName = p.matchConnectionName(); !hasName {
				return "", p.syntaxErrorAtCur()
			}
		}
		stmt.RelayChannel = name
		return "RELAY LOGS", nil
	case "SLAVE", "REPLICA":
		p.advance()
		p.matchConnectionName()
		return "SLAVE", nil
	}
	if p.cur.Type == tokIDENT && flushResetTables[strings.ToUpper(p.cur.Str)] {
		return strings.ToUpper(p.advance().Str), nil
	}
	return "", p.syntaxErrorAtCur()
}

// matchConnectionName consumes a connection name, the grammar's TEXT_STRING_sys: one
// quoted string, with no introducer.
func (p *Parser) matchConnectionName() (string, bool) {
	if p.cur.Type != tokSCONST {
		return "", false
	}
	return p.advance().Str, true
}

// parseDeleteDomainIDs parses the list after BINARY LOGS DELETE_DOMAIN_ID, which may be
// empty.
//
//	= ( [[+]number [, [+]number] ...] )
//
// Each number is any decimal literal, and the server keeps the integer digits it starts
// with: 1.5 deletes domain 1. One above 2^32-1 is refused while parsing with 1076.
func (p *Parser) parseDeleteDomainIDs() error {
	if _, err := p.expect('='); err != nil {
		return err
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	if _, ok := p.match(')'); ok {
		return nil
	}
	for {
		p.match('+')
		if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
			return p.syntaxErrorAtCur()
		}
		if !fitsUint32Prefix(p.cur.Str) {
			return &ParseError{
				Message: "The value of gtid domain being deleted ('" + p.cur.Str +
					"') exceeds its maximum size of 32 bit unsigned integer",
				Position: p.cur.Loc,
			}
		}
		p.advance()
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, err := p.expect(')')
	return err
}

// fitsUint32Prefix reports whether the run of digits num starts with is at most 2^32-1,
// as my_strtoll10 reads it.
func fitsUint32Prefix(num string) bool {
	digits := strings.TrimLeft(num[:len(num)-len(strings.TrimLeft(num, "0123456789"))], "0")
	return len(digits) < 10 || (len(digits) == 10 && digits <= "4294967295")
}

// flushWord returns the word at p.cur upper-cased when it is written bare, as every
// flush_option keyword must be, and "" otherwise.
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

	p.match(kwSONAME)

	soname := p.cur.Str
	p.expect(tokSCONST)

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
		comp := p.cur.Str
		p.expect(tokSCONST)
		stmt.Components = append(stmt.Components, comp)
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
		comp := p.cur.Str
		p.expect(tokSCONST)
		stmt.Components = append(stmt.Components, comp)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
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
// Ref: https://mariadb.com/kb/en/create-server/
//
//	CREATE [OR REPLACE] SERVER [IF NOT EXISTS] server_name
//	  FOREIGN DATA WRAPPER wrapper_name
//	  OPTIONS (option [, option] ...)
func (p *Parser) parseCreateServerStmt(start int, orReplace bool) (*nodes.CreateServerStmt, error) {
	p.advance() // consume SERVER

	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}

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
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
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

	name, _, err := p.parseIdentifier()
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
//	CREATE SPATIAL REFERENCE SYSTEM [IF NOT EXISTS] srid srs_attribute ...
//
//	srs_attribute: {
//	  NAME 'srs_name'
//	| DEFINITION 'definition'
//	| ORGANIZATION 'org_name' IDENTIFIED BY srid
//	| DESCRIPTION 'description'
//	}
//
// MariaDB 11.8 has no such statement (1064). MySQL's OR REPLACE form is
// rejected by parseCreateDispatch.
func (p *Parser) parseCreateSpatialRefSysStmt(start int) (*nodes.CreateSpatialRefSysStmt, error) {
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

	stmt := &nodes.CreateSpatialRefSysStmt{Loc: nodes.Loc{Start: start}}

	ifNotExists, err := p.parseCreateIfNotExists(false)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseLockInstanceStmt parses LOCK INSTANCE FOR BACKUP.
// LOCK already consumed. p.cur is INSTANCE.
func (p *Parser) parseLockInstanceStmt(start int) (*nodes.LockInstanceStmt, error) {
	p.advance() // consume INSTANCE
	p.advance() // consume FOR
	p.advance() // consume BACKUP

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

	if p.cur.Type == tokSCONST {
		stmt.Str = p.cur.Str
		p.advance()
	}

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

	if p.cur.Type == kwIF {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}

	if p.isLvalueIdentToken() {
		name, _, _ := p.parseLvalueIdent()
		stmt.Variable = name
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
