package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseLoadDataStmt parses a LOAD DATA statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/load-data.html
//
//	LOAD DATA [LOCAL] INFILE 'file_name'
//	    [REPLACE | IGNORE]
//	    INTO TABLE tbl_name
//	    [FIELDS [TERMINATED BY 'string'] [[OPTIONALLY] ENCLOSED BY 'char'] [ESCAPED BY 'char']]
//	    [LINES [STARTING BY 'string'] [TERMINATED BY 'string']]
//	    [IGNORE number {LINES | ROWS}]
//	    [(col_name_or_user_var [, col_name_or_user_var] ...)]
//	    [SET col_name={expr | DEFAULT} [, col_name={expr | DEFAULT}] ...]
func (p *Parser) parseLoadDataStmt(start int) (*nodes.LoadDataStmt, error) {
	p.advance() // consume DATA or XML

	stmt := &nodes.LoadDataStmt{Loc: nodes.Loc{Start: start}}

	// [LOCAL]
	if p.cur.Type == kwLOCAL {
		stmt.Local = true
		p.advance()
	}

	// INFILE 'file_name'
	if _, err := p.expect(kwINFILE); err != nil {
		return nil, err
	}
	if p.cur.Type == tokSCONST {
		stmt.Infile = p.cur.Str
		p.advance()
	}

	// [REPLACE | IGNORE]
	if _, ok := p.match(kwREPLACE); ok {
		stmt.Replace = true
	} else if _, ok := p.match(kwIGNORE); ok {
		stmt.Ignore = true
	}

	// INTO TABLE tbl_name
	p.match(kwINTO)
	p.match(kwTABLE)
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	// [FIELDS ...]
	if p.cur.Type == kwFIELDS || (p.cur.Type == tokIDENT && eqFold(p.cur.Str, "columns")) {
		p.advance()
		p.parseFieldsClause(stmt)
	}

	// [LINES ...]
	if p.cur.Type == kwLINES {
		p.advance()
		p.parseLinesClause(stmt)
	}

	// [IGNORE number {LINES | ROWS}]
	if p.cur.Type == kwIGNORE {
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.IgnoreRows = int(p.cur.Ival)
			p.advance()
		}
		p.match(kwLINES, kwROWS)
	}

	// [(col_name_or_user_var, ...)]
	if p.cur.Type == '(' {
		p.advance()
		for {
			col, err := p.parseColumnRef()
			if err != nil {
				return nil, err
			}
			stmt.Columns = append(stmt.Columns, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		p.match(')')
	}

	// [SET ...]
	if p.cur.Type == kwSET {
		p.advance()
		for {
			col, err := p.parseColumnRef()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			val, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.SetList = append(stmt.SetList, &nodes.Assignment{
				Loc:    nodes.Loc{Start: col.Loc.Start, End: p.pos()},
				Column: col,
				Value:  val,
			})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseFieldsClause parses FIELDS/COLUMNS clause options.
func (p *Parser) parseFieldsClause(stmt *nodes.LoadDataStmt) {
	for {
		if p.cur.Type == kwTERMINATED {
			p.advance()
			p.match(kwBY)
			if p.cur.Type == tokSCONST {
				stmt.FieldsTerminatedBy = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwENCLOSED {
			p.advance()
			p.match(kwBY)
			if p.cur.Type == tokSCONST {
				stmt.FieldsEnclosedBy = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwOPTIONALLY {
			p.advance()
			if p.cur.Type == kwENCLOSED {
				p.advance()
				p.match(kwBY)
				if p.cur.Type == tokSCONST {
					stmt.FieldsEnclosedBy = p.cur.Str
					p.advance()
				}
			}
		} else if p.cur.Type == kwESCAPED {
			p.advance()
			p.match(kwBY)
			if p.cur.Type == tokSCONST {
				stmt.FieldsEscapedBy = p.cur.Str
				p.advance()
			}
		} else {
			break
		}
	}
}

// parseLinesClause parses LINES clause options.
func (p *Parser) parseLinesClause(stmt *nodes.LoadDataStmt) {
	for {
		if p.cur.Type == kwTERMINATED {
			p.advance()
			p.match(kwBY)
			if p.cur.Type == tokSCONST {
				stmt.LinesTerminatedBy = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwSTARTING {
			p.advance()
			p.match(kwBY)
			if p.cur.Type == tokSCONST {
				p.advance()
			}
		} else {
			break
		}
	}
}
