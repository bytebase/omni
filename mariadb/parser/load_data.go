package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
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
//	    [{FIELDS | COLUMNS} field_item ...]
//	    [LINES line_item ...]
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
		if err := p.parseFieldsClause(stmt); err != nil {
			return nil, err
		}
	}

	// [LINES ...]
	if p.cur.Type == kwLINES {
		p.advance()
		if err := p.parseLinesClause(stmt); err != nil {
			return nil, err
		}
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
	cols, err := p.parseColumnOrUserVarList()
	if err != nil {
		return nil, err
	}
	stmt.Columns = cols

	// [SET ...]
	if p.cur.Type == kwSET {
		p.advance()
		for {
			col, err := p.parseNamedColumn()
			if err != nil {
				return nil, err
			}
			if err := p.expectEqual(); err != nil {
				return nil, err
			}
			val, ok := p.parseIgnoreValue()
			if !ok {
				var err error
				val, err = p.parseExpr()
				if err != nil {
					return nil, err
				}
				if val == nil {
					return nil, p.syntaxErrorAtCur()
				}
			}
			stmt.SetList = append(stmt.SetList, &nodes.Assignment{
				Loc:    nodes.Loc{Start: col.Loc.Start, End: p.prev.End},
				Column: col,
				Value:  val,
			})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFieldsClause parses the one-or-more items of a FIELDS/COLUMNS clause.
//
// The list is never empty -- `FIELDS` with nothing after it is a 1064 on
// mysql:8.0.46, mariadb:11.8.8, OceanBase 4.3.5 and pingcap/tidb:v8.5.5
// alike -- its items come in any order, and a repeat is allowed, the last one
// winning as it does in the server's own AST.
//
//	field_item:
//	    TERMINATED BY terminator
//	  | [OPTIONALLY] ENCLOSED BY terminator
//	  | ESCAPED BY terminator
//
// ENCLOSED BY is not optional after OPTIONALLY: `FIELDS OPTIONALLY ENCLOSED
// IGNORE 1 LINES` is a 1064 on all four servers, where this used to read
// OPTIONALLY on its own (read on 2026-11-04).
//
// Ref: field_term in sql_yacc.yy.
func (p *Parser) parseFieldsClause(stmt *nodes.LoadDataStmt) error {
	for n := 0; ; n++ {
		switch p.cur.Type {
		case kwTERMINATED:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.FieldsTerminatedBy = v

		case kwOPTIONALLY:
			p.advance()
			if _, err := p.expect(kwENCLOSED); err != nil {
				return err
			}
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.FieldsOptionalEncl = true
			stmt.FieldsEnclosedBy = v

		case kwENCLOSED:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.FieldsEnclosedBy = v

		case kwESCAPED:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.FieldsEscapedBy = v

		default:
			if n == 0 {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
	}
}

// parseLinesClause parses the one-or-more items of a LINES clause. As in a
// FIELDS clause the list is never empty -- a bare `LINES` is a 1064 on
// mysql:8.0.46, mariadb:11.8.8 and OceanBase 4.3.5, though TiDB v8.5.5 takes
// one -- and its two items come in any order and may repeat.
//
// Ref: line_term in sql_yacc.yy.
func (p *Parser) parseLinesClause(stmt *nodes.LoadDataStmt) error {
	for n := 0; ; n++ {
		switch p.cur.Type {
		case kwTERMINATED:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.LinesTerminatedBy = v

		case kwSTARTING:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.LinesStartingBy = v

		default:
			if n == 0 {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
	}
}

// parseFieldTerminator parses the value a FIELDS or LINES option takes: a
// string literal, a hex literal, or a bit literal, the last two decoded to
// the bytes they stand for. So `TERMINATED BY X'6B6B'` and `TERMINATED BY
// B'110101101101011'` both read as "kk". mysql:8.0.46, mariadb:11.8.8,
// OceanBase 4.3.5 and pingcap/tidb:v8.5.5 take the same three spellings
// (read on 2026-11-04).
//
// Ref: text_string in sql_yacc.yy, which field_term and line_term both take.
func (p *Parser) parseFieldTerminator() (string, error) {
	switch p.cur.Type {
	case tokSCONST:
		tok := p.advance()
		return tok.Str, nil
	case tokXCONST, tokBCONST:
		tok := p.cur
		v, err := p.decodeByteLiteral(tok)
		if err != nil {
			return "", err
		}
		p.advance()
		return v, nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}

// parseColumnOrUserVarList parses an optional parenthesized
// column-and-user-variable list. A user variable is kept as a VariableRef and
// not as a column whose name begins with '@': `(a, @v)` names one column and
// one variable.
func (p *Parser) parseColumnOrUserVarList() ([]nodes.ExprNode, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance()
	var out []nodes.ExprNode
	for {
		if p.isUserVariable() {
			v, err := p.parseVariableRef()
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		} else {
			col, err := p.parseColumnRef()
			if err != nil {
				return nil, err
			}
			out = append(out, col)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	p.match(')')
	return out, nil
}

// isUserVariable reports whether p.cur is a user variable (@name).
//
// The lexer emits one as tokIDENT carrying the '@', and it emits a
// backquoted name with its quotes stripped, so a column named `@v` arrives
// as the same token string. The source text tells them apart: a user
// variable's starts at the '@' and a quoted name's at the backquote. Without
// that check `(`@v`)` -- a column every one of these servers takes -- would
// read as a variable.
func (p *Parser) isUserVariable() bool {
	if p.cur.Type != tokIDENT || len(p.cur.Str) < 2 || p.cur.Str[0] != '@' {
		return false
	}
	return strings.HasPrefix(p.inputText(p.cur.Loc, p.cur.End), "@")
}
