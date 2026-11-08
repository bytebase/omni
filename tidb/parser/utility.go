package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

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
//	FLUSH [NO_WRITE_TO_BINLOG | LOCAL] flush_option
//
//	flush_option:
//	    PRIVILEGES | STATUS | HOSTS | CLIENT_ERRORS_SUMMARY
//	  | TIDB PLUGINS name [, name] ...
//	  | [BINARY | ENGINE | ERROR | GENERAL | SLOW] LOGS
//	  | {TABLE | TABLES} [tbl_name [, tbl_name] ...] [WITH READ LOCK]
//
// There is one option. MySQL's list of them, `FLUSH HOSTS, STATUS`, and its
// RELAY LOGS, OPTIMIZER_COSTS, USER_RESOURCES and FOR EXPORT are all 1064s on
// v8.5.5 (tidb_admin_forms_test.go measures them). PLUGINS is a list and each entry an
// Identifier, so `FLUSH TIDB PLUGINS a, status` names two plugins, and every option
// word has to be written bare.
//
// Ref: TiDB v8.5.5 parser.y:12147-12200 (FlushStmt, PluginNameList, FlushOption).
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

	switch {
	case p.cur.Type == kwTABLES || p.cur.Type == kwTABLE:
		stmt.Options = append(stmt.Options, "TABLES")
		p.advance()
		// TableNameListOpt: WITH is reserved, so it cannot be read as a name.
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
		if _, ok := p.match(kwWITH); ok {
			if _, err := p.expect(kwREAD); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwLOCK); err != nil {
				return nil, err
			}
			stmt.WithReadLock = true
		}

	case p.cur.Type == kwLOGS:
		p.advance()
		stmt.Options = append(stmt.Options, "LOGS")

	case p.isBareWord("TIDB"):
		p.advance()
		if _, err := p.expect(kwPLUGINS); err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, "TIDB PLUGINS")
		for {
			name, _, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			stmt.Plugins = append(stmt.Plugins, name)
			if _, ok := p.match(','); !ok {
				break
			}
		}

	default:
		word := p.brieWord()
		switch word {
		case "PRIVILEGES", "STATUS", "HOSTS", "CLIENT_ERRORS_SUMMARY":
			p.advance()
			stmt.Options = append(stmt.Options, word)
		case "BINARY", "ENGINE", "ERROR", "GENERAL", "SLOW":
			p.advance()
			if _, err := p.expect(kwLOGS); err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, word+" LOGS")
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// isWholeExpression reports whether what parseExpr returned is an Expression in the
// grammar's sense. parseExpr also returns forms only an enclosing production allows,
// each readable as a whole argument here: `*`, `t.*`, bare DEFAULT and an INTERVAL.
// `KILL foo(*)`, `KILL foo(t.*)`, `KILL foo(DEFAULT)` and `KILL foo(INTERVAL 1 DAY)`
// are 1064s, and `KILL foo(DEFAULT(a))` parses.
func isWholeExpression(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.StarExpr, *nodes.IntervalExpr:
		return false
	case *nodes.ColumnRef:
		return !e.Star
	case *nodes.DefaultExpr:
		return e.Column != ""
	}
	return true
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
//	KILL [TIDB] [CONNECTION | QUERY] number
//	KILL [TIDB] function_call
//
// There are no other expressions: `KILL 1 + 1`, `KILL a` and `KILL (1)` are 1064s.
// The number is the grammar's NUM, an integer literal up to the whole uint64 range.
// CONNECTION and QUERY take only a number, so `KILL QUERY CONNECTION_ID()` is a 1064
// too. TIDB is not an optional word in front of a name: `KILL tidb(1)` is a 1064
// where `KILL foo(1)` parses, and `KILL `tidb`(1)` calls a function of that name.
//
// Ref: TiDB v8.5.5 parser.y:15289-15331 (KillStmt, KillOrKillTiDB).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-kill/
func (p *Parser) parseKillStmt() (*nodes.KillStmt, error) {
	start := p.pos()
	p.advance() // consume KILL

	stmt := &nodes.KillStmt{Loc: nodes.Loc{Start: start}}
	if p.isBareWord("TIDB") {
		stmt.TiDB = true
		p.advance()
	}

	var numOnly bool
	switch p.cur.Type {
	case kwCONNECTION:
		p.advance()
		numOnly = true
	case kwQUERY:
		stmt.Query = true
		p.advance()
		numOnly = true
	}

	if p.cur.Type == tokICONST || numOnly {
		numLoc := p.cur.Loc
		ival := p.cur.Ival
		id, err := p.parseLengthNum()
		if err != nil {
			return nil, err
		}
		stmt.ID = id
		stmt.ConnectionID = &nodes.IntLit{Loc: nodes.Loc{Start: numLoc, End: p.prev.End}, Value: ival}
	} else {
		call, err := p.parseKillFunctionCall()
		if err != nil {
			return nil, err
		}
		stmt.ConnectionID = call
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseKillFunctionCall parses the grammar's BuiltinFunction, which is narrower than
// its name: a function named by `identifier`, or REPLACE with at least one argument,
// in any number of parentheses.
//
//	( function_call )
//	identifier ( [expr, ...] )
//	REPLACE ( expr, ... )
//
// An unreserved keyword is not an identifier here, which is the point of the table
// curIsGrammarIdentifier reads: `KILL ascii('a')`, `KILL now()` and `KILL user()` are
// 1064s. A name with a dot is one too, and nothing may follow the call.
//
// The parentheses are counted in a loop, so a run of them cannot overflow the stack.
// Each argument goes through parseExpr, which has its own depth check.
//
// Ref: TiDB v8.5.5 parser.y:4058-4084 (BuiltinFunction).
func (p *Parser) parseKillFunctionCall() (*nodes.FuncCallExpr, error) {
	parens := 0
	for p.cur.Type == '(' {
		parens++
		p.advance()
	}

	isReplace := p.cur.Type == kwREPLACE
	if !isReplace && !p.curIsGrammarIdentifier() {
		return nil, p.syntaxErrorAtCur()
	}
	name := p.advance()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	call := &nodes.FuncCallExpr{Loc: nodes.Loc{Start: name.Loc}, Name: name.Str, HasParens: true}
	if p.cur.Type != ')' || isReplace {
		for {
			argTok := p.cur
			arg, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if arg == nil || !isWholeExpression(arg) {
				return nil, &ParseError{Message: "expected an expression", Position: argTok.Loc}
			}
			call.Args = append(call.Args, arg)
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	call.Loc.End = p.prev.End

	for ; parens > 0; parens-- {
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	return call, nil
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

	// ProcedureCall reads `identifier` as FunctionCallGeneric does, with or without the
	// parentheses, so what is no function name is no bare procedure name either: v8.5.5
	// answers 1064 for CALL comment and CALL comment(1) and parses CALL `comment` and
	// CALL d.comment. MySQL reads all four.
	if p.cur.Type != tokIDENT && p.peekNext().Type != '.' && keywordsNotFunctionNames[strings.ToLower(p.cur.Str)] {
		return nil, p.syntaxErrorAtCur()
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

		// @host: the lexer folds `@'host'` into one tokIDENT (Str="@host"); a
		// bare `@` (a gap before the host) leaves the host as the next token.
		if p.cur.Type != tokIDENT || len(p.cur.Str) == 0 || p.cur.Str[0] != '@' {
			return nil, &ParseError{
				Message:  "expected @ after user",
				Position: p.cur.Loc,
			}
		}
		if len(p.cur.Str) > 1 {
			stmt.Host = p.cur.Str[1:]
			p.advance()
		} else {
			p.advance() // bare '@'
			if p.cur.Type != tokSCONST {
				return nil, &ParseError{
					Message:  "expected string literal for host",
					Position: p.cur.Loc,
				}
			}
			stmt.Host = p.cur.Str
			p.advance()
		}

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
