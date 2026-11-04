package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's bulk-load statements: LOAD DATA and IMPORT INTO.
//
// LOAD DATA began as MySQL's and TiDB has grown three clauses of its own on
// it -- FORMAT, FIELDS DEFINED NULL BY, and the WITH option list -- while
// dropping nothing MySQL's grammar carries in the positions below. IMPORT
// INTO is TiDB's alone and shares three of LOAD DATA's pieces: the
// column-and-user-variable list, the SET list, and that WITH list.
//
// Every verdict below was read from pingcap/tidb:v8.5.5 through PREPARE on
// 2026-11-04 and is held to that server by TestLoadDataFormsMatchTiDB.
//
// Ref: LoadDataStmt, ImportIntoStmt, Fields, FieldItem, FieldTerminator,
// Lines and LoadDataOption in pingcap/tidb pkg/parser/parser.y.
// https://docs.pingcap.com/tidb/stable/sql-statement-load-data/
// https://docs.pingcap.com/tidb/stable/sql-statement-import-into/

// parseLoadDataStmt parses a LOAD DATA or LOAD XML statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/load-data.html
// Ref: https://dev.mysql.com/doc/refman/8.0/en/load-xml.html
//
//	LOAD DATA
//	    [LOW_PRIORITY | CONCURRENT] [LOCAL]
//	    INFILE 'file_name'
//	    [FORMAT 'format_name']
//	    [REPLACE | IGNORE]
//	    INTO TABLE tbl_name
//	    [PARTITION (partition_name [, partition_name] ...)]
//	    [CHARACTER SET charset_name]
//	    [{FIELDS | COLUMNS} field_item ...]
//	    [LINES
//	        [STARTING BY terminator]
//	        [TERMINATED BY terminator]
//	    ]
//	    [IGNORE number LINES]
//	    [(col_name_or_user_var [, col_name_or_user_var] ...)]
//	    [SET col_name={expr | DEFAULT} [, col_name={expr | DEFAULT}] ...]
//	    [WITH load_option [, load_option] ...]
//
//	LOAD XML
//	    [LOW_PRIORITY | CONCURRENT] [LOCAL]
//	    INFILE 'file_name'
//	    [REPLACE | IGNORE]
//	    INTO TABLE [db_name.]tbl_name
//	    [CHARACTER SET charset_name]
//	    [ROWS IDENTIFIED BY '<tagname>']
//	    [IGNORE number LINES]
//	    [(field_name_or_user_var [, field_name_or_user_var] ...)]
//	    [SET col_name={expr | DEFAULT} [, col_name={expr | DEFAULT}] ...]
//
// CONCURRENT, PARTITION and the whole of LOAD XML are MySQL's and not
// TiDB's: v8.5.5 answers 1064 for each. They are still read here, which
// leaves this parser taking three forms the engine refuses; narrowing them
// means dropping AST the deparser and its tests hold and is its own change.
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

	// [FORMAT 'format_name'], which TiDB has and MySQL has not. It sits
	// between the file and REPLACE/IGNORE and nowhere else: v8.5.5 answers
	// 1064 for the same word after INTO TABLE, for a second FORMAT, and for
	// a hex literal in place of the string.
	if !isXML {
		format, has, err := p.parseLoadFormat()
		if err != nil {
			return nil, err
		}
		stmt.Format, stmt.HasFormat = format, has
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

	// [{FIELDS | COLUMNS} field_item ...]
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

	// [IGNORE number LINES]. The unit is LINES and only LINES: v8.5.5
	// answers 1064 for IGNORE 2 ROWS, where MySQL 8.0, MariaDB 11.8.8 and
	// OceanBase 4.3.5 all take it.
	if p.cur.Type == kwIGNORE {
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.IgnoreRows = int(p.cur.Ival)
			p.advance()
		}
		if _, err := p.expect(kwLINES); err != nil {
			return nil, err
		}
	}

	// [(col_name_or_user_var, ...)]
	cols, err := p.parseColumnOrUserVarList()
	if err != nil {
		return nil, err
	}
	stmt.Columns = cols

	// [SET ...]
	sets, err := p.parseLoadSetList()
	if err != nil {
		return nil, err
	}
	stmt.SetList = sets

	// [WITH load_option, ...], which TiDB has and MySQL has not.
	if !isXML {
		opts, err := p.parseLoadDataOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseImportIntoStmt parses TiDB's IMPORT INTO. IMPORT is consumed and
// p.cur is INTO.
//
//	IMPORT INTO tbl_name
//	    [(col_name_or_user_var [, col_name_or_user_var] ...)]
//	    [SET col_name={expr | DEFAULT} [, ...]]
//	    FROM 'file_path' [FORMAT 'format_name'] [WITH load_option, ...]
//	IMPORT INTO tbl_name
//	    [(col_name [, col_name] ...)]
//	    FROM select_statement [WITH load_option, ...]
//
// FROM is not optional -- `IMPORT INTO t` is a 1064 -- and the two sources
// take different tails: FORMAT follows a path and not a query, so v8.5.5
// answers 1064 for `FROM 'f' WITH detached FORMAT 'csv'` as it does for a
// FORMAT after a SELECT.
func (p *Parser) parseImportIntoStmt(start int) (*nodes.ImportIntoStmt, error) {
	p.advance() // consume INTO

	stmt := &nodes.ImportIntoStmt{Loc: nodes.Loc{Start: start}}

	// Completion: after IMPORT INTO, offer table_ref candidates.
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

	cols, err := p.parseColumnOrUserVarList()
	if err != nil {
		return nil, err
	}
	stmt.Columns = cols

	sets, err := p.parseLoadSetList()
	if err != nil {
		return nil, err
	}
	stmt.SetList = sets

	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	// The source: a file path, or one query expression. FORMAT follows a
	// path only.
	//
	// VALUES is a query primary only when it is spelled with ROW, here as
	// in every other position: v8.5.5 answers 1064 for `IMPORT INTO t FROM
	// VALUES (1)` and takes `FROM VALUES ROW(1)`.
	switch p.cur.Type {
	case tokSCONST:
		stmt.Path = p.cur.Str
		p.advance()

		format, has, err := p.parseLoadFormat()
		if err != nil {
			return nil, err
		}
		stmt.Format, stmt.HasFormat = format, has

	case kwSELECT, kwWITH, '(', kwTABLE, kwVALUES:
		// FROM takes one query expression, so a TABLE or VALUES target is
		// the query it begins rather than a leaf, and the set operation and
		// the clauses after it have to reach the tree. v8.5.5 accepts
		// `IMPORT INTO t FROM TABLE t1 UNION TABLE t2` and the same source
		// with INTERSECT, EXCEPT, UNION ALL, ORDER BY a, LIMIT 1 and
		// written parenthesized. None of that fits a bare TableStmt, which
		// is why the node keeps one Select rather than a field per primary;
		// a reader that took only the primary would accept the statement
		// and drop the rest, which no corpus row can see.
		sel, err := p.parseSelectStmt()
		if err != nil {
			return nil, err
		}
		stmt.Select = sel

	default:
		return nil, p.syntaxErrorAtCur()
	}

	// A query source carries neither a user variable in the column list nor
	// a SET list: v8.5.5 answers 1064 for `IMPORT INTO t (a,@1) FROM SELECT
	// a FROM s` and for `IMPORT INTO t SET a=1 FROM SELECT * FROM s`, both
	// of which it takes with a path in place of the query.
	if stmt.Select != nil {
		for _, c := range stmt.Columns {
			if v, ok := c.(*nodes.VariableRef); ok {
				return nil, &ParseError{
					Message:  "a user variable is no column of an IMPORT INTO that reads a query",
					Position: v.Loc.Start,
				}
			}
		}
		if len(stmt.SetList) > 0 {
			return nil, &ParseError{
				Message:  "SET is no part of an IMPORT INTO that reads a query",
				Position: stmt.SetList[0].Loc.Start,
			}
		}
	}

	opts, err := p.parseLoadDataOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseLoadFormat parses an optional `FORMAT 'format_name'`. The name is a
// string literal and nothing else: v8.5.5 answers 1064 for FORMAT X'6363'
// and for a backquoted `format`, which is a name rather than this word.
func (p *Parser) parseLoadFormat() (string, bool, error) {
	if p.cur.Type != kwFORMAT {
		return "", false, nil
	}
	p.advance()
	if p.cur.Type != tokSCONST {
		return "", false, p.syntaxErrorAtCur()
	}
	format := p.cur.Str
	p.advance()
	return format, true, nil
}

// parseColumnOrUserVarList parses an optional parenthesized
// column-and-user-variable list. A user variable is kept as a VariableRef and
// not as a column whose name begins with '@': `(a, @v)` names one column and
// one variable. The list may be empty -- v8.5.5 takes `IMPORT INTO t () FROM
// 'f'` -- so an empty pair of parentheses is a written list with no entries.
func (p *Parser) parseColumnOrUserVarList() ([]nodes.ExprNode, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance()
	var out []nodes.ExprNode
	if p.cur.Type == ')' {
		p.advance()
		return out, nil
	}
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
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
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

// parseLoadSetList parses an optional `SET col = expr [, ...]`, shared by
// LOAD DATA and IMPORT INTO.
func (p *Parser) parseLoadSetList() ([]*nodes.Assignment, error) {
	if p.cur.Type != kwSET {
		return nil, nil
	}
	p.advance()
	var out []*nodes.Assignment
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
		if val == nil {
			return nil, p.syntaxErrorAtCur()
		}
		out = append(out, &nodes.Assignment{
			Loc:    nodes.Loc{Start: col.Loc.Start, End: p.prev.End},
			Column: col,
			Value:  val,
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return out, nil
}

// parseLoadDataOptions parses TiDB's optional `WITH name [= literal] [, ...]`.
//
// A name is an identifier, so a reserved word is not one -- `WITH select=1`
// is a 1064 -- and a backquoted `WITH `detached“ is, because a backquoted
// word is an identifier. A value is a signed literal: a number, a string, a
// byte literal, TRUE/FALSE or NULL, optionally signed, and not an
// expression, so `WITH threads=@a` is a 1064. The list holds at least one
// option and no trailing comma, and one statement carries one WITH.
//
// None of these names is registered as a keyword. v8.5.5 spells DETACHED and
// THREAD as keywords of its own but they are ordinary names everywhere else,
// and registering one here would change every site that reads a name.
func (p *Parser) parseLoadDataOptions() ([]*nodes.LoadDataOption, error) {
	if p.cur.Type != kwWITH {
		return nil, nil
	}
	p.advance()
	var out []*nodes.LoadDataOption
	for {
		start := p.cur.Loc
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		opt := &nodes.LoadDataOption{Loc: nodes.Loc{Start: start}, Name: name}
		if _, ok := p.match('='); ok {
			// TiDB's SignedLiteral, which parseAdminSignedLiteral already
			// holds to v8.5.5 for ADMIN's own option values.
			val, err := p.parseAdminSignedLiteral()
			if err != nil {
				return nil, err
			}
			opt.Value = val
		}
		opt.Loc.End = p.prev.End
		out = append(out, opt)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return out, nil
}

// parseFieldsClause parses the one-or-more field items of a FIELDS/COLUMNS
// clause.
//
// The list is never empty -- `FIELDS` with nothing after it is a 1064 on
// v8.5.5, MySQL 8.0, MariaDB 11.8.8 and OceanBase 4.3.5 alike -- its items
// come in any order, and a repeat is allowed, the last one winning as it
// does in the engine's own AST.
//
//	field_item:
//	    TERMINATED BY terminator
//	  | [OPTIONALLY] ENCLOSED BY terminator
//	  | ESCAPED BY terminator
//	  | DEFINED NULL BY terminator [OPTIONALLY ENCLOSED]
//
// DEFINED NULL BY is TiDB's: MySQL 8.0, MariaDB 11.8.8 and OceanBase 4.3.5
// each answer 1064 for it. Its OPTIONALLY ENCLOSED tail is why OPTIONALLY
// needs ENCLOSED BY after it here: TiDB's lexer reads `OPTIONALLY ENCLOSED
// BY` as one token to tell the two apart, and `FIELDS OPTIONALLY ENCLOSED
// IGNORE 1 LINES` -- OPTIONALLY ENCLOSED with no BY and no DEFINED NULL BY
// before it -- is a 1064 on all four servers.
func (p *Parser) parseFieldsClause(stmt *nodes.LoadDataStmt) error {
	for n := 0; ; n++ {
		switch {
		case p.cur.Type == kwTERMINATED:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.FieldsTerminatedBy = v

		case p.cur.Type == kwOPTIONALLY:
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

		case p.cur.Type == kwENCLOSED:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.FieldsEnclosedBy = v
			// Each ENCLOSED item carries its own OPTIONALLY, so a plain
			// one after an OPTIONALLY one clears the flag. v8.5.5 writes
			// `"7"` for `INTO OUTFILE ... FIELDS OPTIONALLY ENCLOSED BY
			// '"' ENCLOSED BY '"'` and `7` with the OPTIONALLY one alone,
			// so the last item wins here as everywhere in this list.
			// MySQL 8.0.46 and MariaDB 11.8.8 do not clear it -- both
			// write `7` for that statement -- so their copy of this
			// clause keeps the flag (all read on 2026-11-04).
			stmt.FieldsOptionalEncl = false

		case p.cur.Type == kwESCAPED:
			p.advance()
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.FieldsEscapedBy = v

		case p.isBareWord("DEFINED"):
			p.advance()
			if _, err := p.expect(kwNULL); err != nil {
				return err
			}
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			v, err := p.parseFieldTerminator()
			if err != nil {
				return err
			}
			stmt.FieldsDefinedNullBy, stmt.HasDefinedNullBy = v, true
			// A repeat without the tail clears it, the last item winning:
			// v8.5.5 reads an enclosed `"a"` as NULL for `DEFINED NULL BY
			// 'a' OPTIONALLY ENCLOSED` and as the string "a" both for
			// `DEFINED NULL BY 'a'` and for `DEFINED NULL BY 'a'
			// OPTIONALLY ENCLOSED DEFINED NULL BY 'a'` (measured by
			// loading a file on 2026-11-04).
			stmt.DefinedNullOptEnclosed = false
			// [OPTIONALLY ENCLOSED], with no BY and no value after it.
			// What follows OPTIONALLY ENCLOSED here settles which of the
			// two it is: a BY makes it the next field item instead, so
			// v8.5.5 takes `DEFINED NULL BY 'a' OPTIONALLY ENCLOSED BY
			// 'b'` as this item and an ENCLOSED one.
			if p.cur.Type == kwOPTIONALLY {
				p.advance()
				if _, err := p.expect(kwENCLOSED); err != nil {
					return err
				}
				if _, ok := p.match(kwBY); ok {
					encl, err := p.parseFieldTerminator()
					if err != nil {
						return err
					}
					stmt.FieldsOptionalEncl = true
					stmt.FieldsEnclosedBy = encl
				} else {
					stmt.DefinedNullOptEnclosed = true
				}
			}

		default:
			if n == 0 {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
	}
}

// parseLinesClause parses a LINES clause.
//
// STARTING comes before TERMINATED and each appears at most once: v8.5.5
// answers 1064 for `LINES TERMINATED BY 'b' STARTING BY 'c'` and for a
// repeat of either, where MySQL 8.0, MariaDB 11.8.8 and OceanBase 4.3.5 read
// a free-order list. TiDB takes a bare LINES, which those three refuse.
func (p *Parser) parseLinesClause(stmt *nodes.LoadDataStmt) error {
	if p.cur.Type == kwSTARTING {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		v, err := p.parseFieldTerminator()
		if err != nil {
			return err
		}
		stmt.LinesStartingBy = v
	}
	if p.cur.Type == kwTERMINATED {
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		v, err := p.parseFieldTerminator()
		if err != nil {
			return err
		}
		stmt.LinesTerminatedBy = v
	}
	return nil
}

// parseFieldTerminator parses the value a FIELDS or LINES option takes: a
// string literal, a hex literal, or a bit literal, the last two decoded to
// the bytes they stand for, as the engine's own FieldTerminator rule does.
// So `TERMINATED BY X'6B6B'` and `TERMINATED BY B'110101101101011'` both
// read as "kk". MySQL 8.0, MariaDB 11.8.8 and OceanBase 4.3.5 take the same
// three spellings here.
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
