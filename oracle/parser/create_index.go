package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateIndexStmt parses a CREATE [UNIQUE|BITMAP] INDEX statement.
// The CREATE keyword has already been consumed. The current token is
// UNIQUE, BITMAP, or INDEX.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-INDEX.html
//
//	CREATE [ UNIQUE | BITMAP ] INDEX [ schema. ] index_name
//	    ON [ schema. ] table_name ( col_expr [ ASC | DESC ] [, ...] )
//	    [ REVERSE ]
//	    [ TABLESPACE tablespace ]
//	    [ LOCAL | GLOBAL ... ]
//	    [ ONLINE ]
//	    [ PARALLEL n | NOPARALLEL ]
//	    [ COMPRESS n | NOCOMPRESS ]
func (p *Parser) parseCreateIndexStmt(start int) *nodes.CreateIndexStmt {
	stmt := &nodes.CreateIndexStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional UNIQUE or BITMAP
	if p.cur.Type == kwUNIQUE {
		stmt.Unique = true
		p.advance()
	} else if p.cur.Type == kwBITMAP {
		stmt.Bitmap = true
		p.advance()
	}

	// INDEX keyword
	if p.cur.Type == kwINDEX {
		p.advance()
	}

	// Index name
	stmt.Name = p.parseObjectName()

	// ON table_name
	if p.cur.Type == kwON {
		p.advance()
	}
	stmt.Table = p.parseObjectName()

	// ( column_list )
	if p.cur.Type == '(' {
		p.advance()
		stmt.Columns = p.parseIndexColumnList()
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Parse trailing options
	p.parseIndexOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseIndexColumnList parses a comma-separated list of index columns.
func (p *Parser) parseIndexColumnList() *nodes.List {
	list := &nodes.List{}
	for {
		col := p.parseIndexColumn()
		if col == nil {
			break
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list
}

// parseIndexColumn parses a single index column expression with optional ASC/DESC.
func (p *Parser) parseIndexColumn() *nodes.IndexColumn {
	start := p.pos()
	expr := p.parseExpr()
	if expr == nil {
		return nil
	}

	col := &nodes.IndexColumn{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// ASC | DESC
	switch p.cur.Type {
	case kwASC:
		col.Dir = nodes.SORTBY_ASC
		p.advance()
	case kwDESC:
		col.Dir = nodes.SORTBY_DESC
		p.advance()
	}

	// NULLS FIRST | NULLS LAST
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			col.NullOrder = nodes.SORTBY_NULLS_FIRST
			p.advance()
		case kwLAST:
			col.NullOrder = nodes.SORTBY_NULLS_LAST
			p.advance()
		}
	}

	col.Loc.End = p.pos()
	return col
}

// parseIndexOptions parses optional trailing clauses for CREATE INDEX.
func (p *Parser) parseIndexOptions(stmt *nodes.CreateIndexStmt) {
	for {
		switch p.cur.Type {
		case kwREVERSE:
			stmt.Reverse = true
			p.advance()
		case kwTABLESPACE:
			p.advance()
			stmt.Tablespace = p.parseIdentifier()
		case kwLOCAL:
			stmt.Local = true
			p.advance()
		case kwGLOBAL:
			stmt.Global = true
			p.advance()
		case kwONLINE:
			stmt.Online = true
			p.advance()
		case kwPARALLEL:
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			} else {
				stmt.Parallel = "PARALLEL"
			}
		case kwNOPARALLEL:
			stmt.Parallel = "NOPARALLEL"
			p.advance()
		case kwCOMPRESS:
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.Compress = p.cur.Str
				p.advance()
			} else {
				stmt.Compress = "COMPRESS"
			}
		case kwNOCOMPRESS:
			stmt.Compress = "NOCOMPRESS"
			p.advance()
		default:
			return
		}
	}
}
