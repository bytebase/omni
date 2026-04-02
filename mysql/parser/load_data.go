package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseLoadDataStmt parses a LOAD DATA or LOAD XML statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/load-data.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/load-xml.html
//
//	LOAD DATA
//	    [LOW_PRIORITY | CONCURRENT] [LOCAL]
//	    INFILE 'file_name'
//	    [REPLACE | IGNORE]
//	    INTO TABLE tbl_name
//	    [PARTITION (partition_name [, partition_name] ...)]
//	    [CHARACTER SET charset_name]
//	    [{FIELDS | COLUMNS}
//	        [TERMINATED BY 'string']
//	        [[OPTIONALLY] ENCLOSED BY 'char']
//	        [ESCAPED BY 'char']
//	    ]
//	    [LINES
//	        [STARTING BY 'string']
//	        [TERMINATED BY 'string']
//	    ]
//	    [IGNORE number {LINES | ROWS}]
//	    [(col_name_or_user_var [, col_name_or_user_var] ...)]
//	    [SET col_name={expr | DEFAULT} [, col_name={expr | DEFAULT}] ...]
//
//	LOAD XML
//	    [LOW_PRIORITY | CONCURRENT] [LOCAL]
//	    INFILE 'file_name'
//	    [REPLACE | IGNORE]
//	    INTO TABLE [db_name.]tbl_name
//	    [CHARACTER SET charset_name]
//	    [ROWS IDENTIFIED BY '<tagname>']
//	    [IGNORE number {LINES | ROWS}]
//	    [(field_name_or_user_var [, field_name_or_user_var] ...)]
//	    [SET col_name={expr | DEFAULT} [, col_name={expr | DEFAULT}] ...]
func (p *Parser) parseLoadDataStmt(start int) (*nodes.LoadDataStmt, error) {
	isXML := p.cur.Type == kwXML
	p.advance() // consume DATA or XML

	stmt := &nodes.LoadDataStmt{Loc: nodes.Loc{Start: start}, IsXML: isXML}

	// [LOW_PRIORITY | CONCURRENT]
	if p.cur.Type == kwLOW_PRIORITY {
		stmt.LowPriority = true
		p.advance()
	} else if p.cur.Type == kwCONCURRENT {
		stmt.Concurrent = true
		p.advance()
	}

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

	// Completion: after INTO TABLE, offer table_ref candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		p.addRuleCandidate("database_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	// [PARTITION (partition_name, ...)] (LOAD DATA only)
	if !isXML && p.cur.Type == kwPARTITION {
		p.advance()
		parts, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Partitions = parts
	}

	// [CHARACTER SET charset_name]
	if p.cur.Type == kwCHARACTER {
		p.advance()
		p.match(kwSET)
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.CharacterSet = name
	} else if p.cur.Type == kwCHARSET {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.CharacterSet = name
	}

	// [ROWS IDENTIFIED BY '<tagname>'] (LOAD XML only)
	if isXML && p.cur.Type == kwROWS {
		p.advance()
		if p.cur.Type == kwIDENTIFIED {
			p.advance()
			p.match(kwBY)
			if p.cur.Type == tokSCONST {
				stmt.RowsIdentifiedBy = p.cur.Str
				p.advance()
			}
		}
	}

	// [{FIELDS | COLUMNS} ...]
	if p.cur.Type == kwFIELDS || p.cur.Type == kwCOLUMNS {
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
		} else if p.cur.Type == kwOPTIONALLY {
			p.advance()
			stmt.FieldsOptionalEncl = true
			if p.cur.Type == kwENCLOSED {
				p.advance()
				p.match(kwBY)
				if p.cur.Type == tokSCONST {
					stmt.FieldsEnclosedBy = p.cur.Str
					p.advance()
				}
			}
		} else if p.cur.Type == kwENCLOSED {
			p.advance()
			p.match(kwBY)
			if p.cur.Type == tokSCONST {
				stmt.FieldsEnclosedBy = p.cur.Str
				p.advance()
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
				stmt.LinesStartingBy = p.cur.Str
				p.advance()
			}
		} else {
			break
		}
	}
}
