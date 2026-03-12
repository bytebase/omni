package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseAnalyzeTableStmt parses an ANALYZE TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/analyze-table.html
//
//	ANALYZE [NO_WRITE_TO_BINLOG | LOCAL] TABLE tbl_name [, tbl_name] ...
func (p *Parser) parseAnalyzeTableStmt() (*nodes.AnalyzeTableStmt, error) {
	start := p.pos()
	p.advance() // consume ANALYZE

	// Optional NO_WRITE_TO_BINLOG | LOCAL
	if p.cur.Type == kwLOCAL {
		p.advance()
	} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "no_write_to_binlog") {
		p.advance()
	}

	// TABLE
	p.match(kwTABLE)

	stmt := &nodes.AnalyzeTableStmt{Loc: nodes.Loc{Start: start}}

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

	if p.cur.Type == kwLOCAL {
		p.advance()
	} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "no_write_to_binlog") {
		p.advance()
	}

	p.match(kwTABLE)

	stmt := &nodes.OptimizeTableStmt{Loc: nodes.Loc{Start: start}}

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
	for p.cur.Type == kwFOR || (p.cur.Type == tokIDENT &&
		(eqFold(p.cur.Str, "quick") || eqFold(p.cur.Str, "fast") ||
			eqFold(p.cur.Str, "medium") || eqFold(p.cur.Str, "extended") ||
			eqFold(p.cur.Str, "changed"))) {
		if p.cur.Type == kwFOR {
			p.advance()
			if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "upgrade") {
				stmt.Options = append(stmt.Options, "FOR UPGRADE")
				p.advance()
			}
		} else {
			stmt.Options = append(stmt.Options, p.cur.Str)
			p.advance()
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

	if p.cur.Type == kwLOCAL {
		p.advance()
	} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "no_write_to_binlog") {
		p.advance()
	}

	p.match(kwTABLE)

	stmt := &nodes.RepairTableStmt{Loc: nodes.Loc{Start: start}}

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

	// Options
	for p.cur.Type == tokIDENT {
		if eqFold(p.cur.Str, "quick") {
			stmt.Quick = true
			p.advance()
		} else if eqFold(p.cur.Str, "extended") {
			stmt.Extended = true
			p.advance()
		} else if eqFold(p.cur.Str, "use_frm") {
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
//	FLUSH [NO_WRITE_TO_BINLOG | LOCAL] flush_option [, flush_option] ...
func (p *Parser) parseFlushStmt() (*nodes.FlushStmt, error) {
	start := p.pos()
	p.advance() // consume FLUSH

	if p.cur.Type == kwLOCAL {
		p.advance()
	} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "no_write_to_binlog") {
		p.advance()
	}

	stmt := &nodes.FlushStmt{Loc: nodes.Loc{Start: start}}

	// Flush options
	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}
		if p.isIdentToken() {
			name, _, _ := p.parseIdentifier()
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

// parseResetStmt parses a RESET statement.
//
//	RESET reset_option [, reset_option] ...
func (p *Parser) parseResetStmt() (*nodes.FlushStmt, error) {
	start := p.pos()
	p.advance() // consume RESET

	stmt := &nodes.FlushStmt{Loc: nodes.Loc{Start: start}}

	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}
		if p.isIdentToken() {
			name, _, _ := p.parseIdentifier()
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
		stmt.Exprs = append(stmt.Exprs, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
