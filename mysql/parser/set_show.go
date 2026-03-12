package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseSetStmt parses a SET statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/set-variable.html
//
//	SET [GLOBAL | SESSION | LOCAL] var = expr [, var = expr] ...
//	SET NAMES charset [COLLATE collation]
//	SET CHARACTER SET charset
func (p *Parser) parseSetStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume SET

	stmt := &nodes.SetStmt{Loc: nodes.Loc{Start: start}}

	// Check for NAMES special form
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "names") {
		p.advance() // consume NAMES
		// Parse charset name
		charset, charsetLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		// Build assignment: NAMES = charset
		stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
			Loc:    nodes.Loc{Start: charsetLoc},
			Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: charsetLoc}, Column: "NAMES"},
			Value:  &nodes.StringLit{Loc: nodes.Loc{Start: charsetLoc}, Value: charset},
		})
		// Optional COLLATE
		if _, ok := p.match(kwCOLLATE); ok {
			collation, collLoc, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
				Loc:    nodes.Loc{Start: collLoc},
				Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: collLoc}, Column: "COLLATE"},
				Value:  &nodes.StringLit{Loc: nodes.Loc{Start: collLoc}, Value: collation},
			})
		}
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// Check for CHARACTER SET special form
	if p.cur.Type == kwCHARACTER {
		p.advance() // consume CHARACTER
		if _, ok := p.match(kwSET); ok {
			charset, charsetLoc, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Assignments = append(stmt.Assignments, &nodes.Assignment{
				Loc:    nodes.Loc{Start: charsetLoc},
				Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: charsetLoc}, Column: "CHARACTER SET"},
				Value:  &nodes.StringLit{Loc: nodes.Loc{Start: charsetLoc}, Value: charset},
			})
			stmt.Loc.End = p.pos()
			return stmt, nil
		}
	}

	// Check for SET DEFAULT ROLE
	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		if p.cur.Type == kwROLE {
			p.advance() // consume ROLE
			return p.parseSetDefaultRoleStmt(start)
		}
		return nil, &ParseError{Message: "expected ROLE after SET DEFAULT", Position: p.cur.Loc}
	}

	// Check for SET ROLE
	if p.cur.Type == kwROLE {
		p.advance() // consume ROLE
		return p.parseSetRoleStmt(start)
	}

	// Check for SET RESOURCE GROUP
	if p.cur.Type == kwRESOURCE {
		return p.parseSetResourceGroupStmt(start)
	}

	// Check for GLOBAL / SESSION / LOCAL scope
	scope := ""
	switch p.cur.Type {
	case kwGLOBAL:
		scope = "GLOBAL"
		p.advance()
	case kwSESSION:
		scope = "SESSION"
		p.advance()
	case kwLOCAL:
		scope = "LOCAL"
		p.advance()
	}

	// SET [GLOBAL|SESSION] TRANSACTION ...
	if p.cur.Type == kwTRANSACTION {
		return p.parseSetTransactionStmt(start, scope)
	}

	stmt.Scope = scope

	// Parse assignment list
	for {
		asgn, err := p.parseSetAssignment()
		if err != nil {
			return nil, err
		}
		stmt.Assignments = append(stmt.Assignments, asgn)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseSetAssignment parses a single SET assignment: var = expr
func (p *Parser) parseSetAssignment() (*nodes.Assignment, error) {
	start := p.pos()

	var col *nodes.ColumnRef

	// Handle @var and @@var references
	if p.isVariableRef() {
		vref, err := p.parseVariableRef()
		if err != nil {
			return nil, err
		}
		// Convert VariableRef to ColumnRef for the assignment
		prefix := "@"
		if vref.System {
			prefix = "@@"
			if vref.Scope != "" {
				prefix = "@@" + vref.Scope + "."
			}
		}
		col = &nodes.ColumnRef{
			Loc:    vref.Loc,
			Column: prefix + vref.Name,
		}
	} else {
		var err error
		col, err = p.parseColumnRef()
		if err != nil {
			return nil, err
		}
	}

	// Expect '='
	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// Parse value expression
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.pos()},
		Column: col,
		Value:  val,
	}, nil
}

// parseShowStmt parses a SHOW statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/show.html
func (p *Parser) parseShowStmt() (*nodes.ShowStmt, error) {
	start := p.pos()
	p.advance() // consume SHOW

	stmt := &nodes.ShowStmt{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwDATABASES:
		stmt.Type = "DATABASES"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwTABLES:
		stmt.Type = "TABLES"
		p.advance()
		// Optional FROM db
		if _, ok := p.match(kwFROM); ok {
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		}
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwFULL:
		p.advance() // consume FULL
		if p.cur.Type == kwCOLUMNS {
			stmt.Type = "FULL COLUMNS"
			p.advance()
			// FROM tbl
			if _, err := p.expect(kwFROM); err != nil {
				return nil, err
			}
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
			// Optional FROM db
			if _, ok := p.match(kwFROM); ok {
				dbRef, err := p.parseTableRef()
				if err != nil {
					return nil, err
				}
				// Merge: set schema on From
				stmt.From.Schema = dbRef.Name
			}
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwCOLUMNS:
		stmt.Type = "COLUMNS"
		p.advance()
		// FROM tbl
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.From = ref
		// Optional FROM db
		if _, ok := p.match(kwFROM); ok {
			dbRef, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From.Schema = dbRef.Name
		}
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwCREATE:
		p.advance() // consume CREATE
		switch p.cur.Type {
		case kwTABLE:
			stmt.Type = "CREATE TABLE"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwDATABASE:
			stmt.Type = "CREATE DATABASE"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwVIEW:
			stmt.Type = "CREATE VIEW"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwPROCEDURE:
			stmt.Type = "CREATE PROCEDURE"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwFUNCTION:
			stmt.Type = "CREATE FUNCTION"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwTRIGGER:
			stmt.Type = "CREATE TRIGGER"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		case kwEVENT:
			stmt.Type = "CREATE EVENT"
			p.advance()
			ref, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From = ref
		}

	case kwINDEX:
		stmt.Type = "INDEX"
		p.advance()
		// FROM tbl
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.From = ref
		// Optional FROM db
		if _, ok := p.match(kwFROM); ok {
			dbRef, err := p.parseTableRef()
			if err != nil {
				return nil, err
			}
			stmt.From.Schema = dbRef.Name
		}

	case kwGLOBAL, kwSESSION:
		scope := "GLOBAL"
		if p.cur.Type == kwSESSION {
			scope = "SESSION"
		}
		p.advance()
		if p.cur.Type == kwVARIABLES {
			stmt.Type = scope + " VARIABLES"
			p.advance()
		} else if p.cur.Type == kwSTATUS {
			stmt.Type = scope + " STATUS"
			p.advance()
		}
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwVARIABLES:
		stmt.Type = "VARIABLES"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwSTATUS:
		stmt.Type = "STATUS"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwWARNINGS:
		stmt.Type = "WARNINGS"
		p.advance()
		// Optional LIMIT
		if _, ok := p.match(kwLIMIT); ok {
			// Just skip the count for now
			p.advance()
		}

	case kwERRORS:
		stmt.Type = "ERRORS"
		p.advance()
		// Optional LIMIT
		if _, ok := p.match(kwLIMIT); ok {
			p.advance()
		}

	case kwENGINES:
		stmt.Type = "ENGINES"
		p.advance()

	case kwPLUGINS:
		stmt.Type = "PLUGINS"
		p.advance()

	case kwMASTER:
		p.advance() // consume MASTER
		if p.cur.Type == kwSTATUS {
			stmt.Type = "MASTER STATUS"
			p.advance()
		}

	case kwSLAVE:
		p.advance() // consume SLAVE
		if p.cur.Type == kwSTATUS {
			stmt.Type = "SLAVE STATUS"
			p.advance()
		}

	case kwREPLICA:
		p.advance() // consume REPLICA
		if p.cur.Type == kwSTATUS {
			stmt.Type = "REPLICA STATUS"
			p.advance()
		}

	case kwBINARY:
		p.advance() // consume BINARY
		if p.cur.Type == kwLOGS {
			stmt.Type = "BINARY LOGS"
			p.advance()
		}

	case kwBINLOG:
		p.advance() // consume BINLOG
		stmt.Type = "BINLOG EVENTS"
		// Expect EVENTS (as identifier since it may not be a keyword)
		if p.cur.Type == kwEVENT || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "events")) {
			p.advance()
		}
		// Optional IN 'log_name'
		if p.cur.Type == kwIN {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.Like = expr
		}
		// Optional FROM pos
		if _, ok := p.match(kwFROM); ok {
			// skip the position value
			p.advance()
		}
		// Optional LIMIT [offset,] count
		if _, ok := p.match(kwLIMIT); ok {
			// skip limit value(s)
			p.advance()
			if p.cur.Type == ',' {
				p.advance() // consume ','
				p.advance() // consume count
			}
		}

	case kwTABLE:
		p.advance() // consume TABLE
		if p.cur.Type == kwSTATUS {
			stmt.Type = "TABLE STATUS"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwTRIGGER:
		// SHOW TRIGGERS (kwTRIGGER won't match plural; handled in default)
		stmt.Type = "TRIGGERS"
		p.advance()
		if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwEVENT:
		// SHOW EVENTS (kwEVENT won't match plural; handled in default)
		stmt.Type = "EVENTS"
		p.advance()
		if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	case kwPROCEDURE:
		p.advance() // consume PROCEDURE
		if p.cur.Type == kwSTATUS {
			stmt.Type = "PROCEDURE STATUS"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwFUNCTION:
		p.advance() // consume FUNCTION
		if p.cur.Type == kwSTATUS {
			stmt.Type = "FUNCTION STATUS"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwOPEN:
		p.advance() // consume OPEN
		if p.cur.Type == kwTABLES {
			stmt.Type = "OPEN TABLES"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwPRIVILEGES:
		stmt.Type = "PRIVILEGES"
		p.advance()

	case kwPROFILES:
		stmt.Type = "PROFILES"
		p.advance()

	case kwRELAYLOG:
		p.advance() // consume RELAYLOG
		stmt.Type = "RELAYLOG EVENTS"
		if p.cur.Type == kwEVENT || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "events")) {
			p.advance()
		}
		// Optional IN 'log_name'
		if p.cur.Type == kwIN {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.Like = expr
		}
		// Optional FROM pos
		if _, ok := p.match(kwFROM); ok {
			p.advance()
		}
		// Optional LIMIT [offset,] count
		if _, ok := p.match(kwLIMIT); ok {
			p.advance()
			if p.cur.Type == ',' {
				p.advance()
				p.advance()
			}
		}

	case kwCHARACTER:
		p.advance() // consume CHARACTER
		if p.cur.Type == kwSET {
			stmt.Type = "CHARACTER SET"
			p.advance()
			if err := p.parseShowLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}

	case kwCOLLATION:
		stmt.Type = "COLLATION"
		p.advance()
		if err := p.parseShowLikeOrWhere(stmt); err != nil {
			return nil, err
		}

	default:
		// Handle GRANTS and PROCESSLIST as identifier-based keywords
		if p.cur.Type == kwGRANT || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "grants")) {
			stmt.Type = "GRANTS"
			p.advance()
			// Optional FOR user
			if _, ok := p.match(kwFOR); ok {
				name, nameLoc, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				stmt.From = &nodes.TableRef{
					Loc:  nodes.Loc{Start: nameLoc},
					Name: name,
				}
				stmt.From.Loc.End = p.pos()
			}
		} else if p.cur.Type == kwPROCESSLIST || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "processlist")) {
			stmt.Type = "PROCESSLIST"
			p.advance()
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "triggers") {
			stmt.Type = "TRIGGERS"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		} else if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "events") {
			stmt.Type = "EVENTS"
			p.advance()
			if err := p.parseShowFromLikeOrWhere(stmt); err != nil {
				return nil, err
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseShowLikeOrWhere parses optional LIKE or WHERE clause for SHOW statements.
func (p *Parser) parseShowLikeOrWhere(stmt *nodes.ShowStmt) error {
	if p.cur.Type == kwLIKE {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Like = expr
	} else if p.cur.Type == kwWHERE {
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return err
		}
		stmt.Where = expr
	}
	return nil
}

// parseShowFromLikeOrWhere parses optional FROM db, LIKE, or WHERE for SHOW statements.
func (p *Parser) parseShowFromLikeOrWhere(stmt *nodes.ShowStmt) error {
	if _, ok := p.match(kwFROM); ok {
		ref, err := p.parseTableRef()
		if err != nil {
			return err
		}
		stmt.From = ref
	}
	return p.parseShowLikeOrWhere(stmt)
}

// parseUseStmt parses a USE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/use.html
//
//	USE db_name
func (p *Parser) parseUseStmt() (*nodes.UseStmt, error) {
	start := p.pos()
	p.advance() // consume USE

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &nodes.UseStmt{
		Loc:      nodes.Loc{Start: start, End: p.pos()},
		Database: name,
	}, nil
}

// parseExplainStmt parses an EXPLAIN or DESCRIBE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/explain.html
//
//	EXPLAIN [EXTENDED | PARTITIONS] {SELECT|INSERT|UPDATE|DELETE|REPLACE} ...
//	EXPLAIN ANALYZE SELECT ...
//	EXPLAIN FORMAT = {TRADITIONAL|JSON|TREE} {SELECT|INSERT|UPDATE|DELETE|REPLACE} ...
//	DESCRIBE tbl_name [col_name | wild]
func (p *Parser) parseExplainStmt() (*nodes.ExplainStmt, error) {
	start := p.pos()
	isDescribe := p.cur.Type == kwDESCRIBE
	p.advance() // consume EXPLAIN or DESCRIBE

	stmt := &nodes.ExplainStmt{Loc: nodes.Loc{Start: start}}

	if isDescribe {
		// DESCRIBE tbl_name [col_name]
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		// Wrap as a ShowStmt for DESCRIBE (which is equivalent to SHOW COLUMNS FROM tbl)
		showStmt := &nodes.ShowStmt{
			Loc:  nodes.Loc{Start: start},
			Type: "COLUMNS",
			From: ref,
		}
		// Optional column name
		if p.cur.Type != tokEOF && p.cur.Type != ';' {
			colExpr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			showStmt.Like = colExpr
		}
		showStmt.Loc.End = p.pos()
		stmt.Stmt = showStmt
		stmt.Loc.End = p.pos()
		return stmt, nil
	}

	// EXPLAIN [ANALYZE] [EXTENDED] [PARTITIONS] [FORMAT = value] stmt

	// Check for ANALYZE
	if p.cur.Type == kwANALYZE {
		stmt.Analyze = true
		p.advance()
	}

	// Check for EXTENDED (deprecated in 8.0 but still parsed)
	if p.cur.Type == kwEXTENDED {
		stmt.Extended = true
		p.advance()
	}

	// Check for PARTITIONS (deprecated in 8.0 but still parsed)
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "partitions") {
		stmt.Partitions = true
		p.advance()
	}

	// Check for FORMAT = value
	if p.cur.Type == tokIDENT && eqFold(p.cur.Str, "format") {
		p.advance() // consume FORMAT
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		// Parse format value: TRADITIONAL, JSON, TREE
		formatName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Format = formatName
	}

	// Parse the explainable statement
	switch p.cur.Type {
	case kwSELECT:
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = sel
	case kwINSERT:
		ins, err := p.parseInsertStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = ins
	case kwUPDATE:
		upd, err := p.parseUpdateStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = upd
	case kwDELETE:
		del, err := p.parseDeleteStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = del
	case kwREPLACE:
		rep, err := p.parseReplaceStmt()
		if err != nil {
			return nil, err
		}
		stmt.Stmt = rep
	default:
		// For other tokens, try to parse as a table ref (EXPLAIN table_name)
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		showStmt := &nodes.ShowStmt{
			Loc:  nodes.Loc{Start: ref.Loc.Start, End: p.pos()},
			Type: "COLUMNS",
			From: ref,
		}
		stmt.Stmt = showStmt
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
