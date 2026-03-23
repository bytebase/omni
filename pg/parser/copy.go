package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCopyStmt parses a COPY statement.
// The COPY keyword has already been consumed.
//
// CopyStmt:
//
//	COPY opt_binary relation_expr opt_column_list
//	    copy_from opt_program copy_file_name copy_delimiter copy_opt_with
//	    copy_options where_clause
//	| COPY '(' PreparableStmt ')' TO opt_program copy_file_name copy_opt_with copy_options
func (p *Parser) parseCopyStmt() (nodes.Node, error) {
	if p.cur.Type == '(' {
		return p.parseCopyQueryStmt()
	}
	optBin := p.parseCopyOptBinary()
	rel, _ := p.parseRelationExpr()
	attlist, _ := p.parseOptColumnList()
	isFrom := p.parseCopyFrom()
	isProgram := p.parseCopyOptProgram()
	filename := p.parseCopyFileName()
	delimOpt := p.parseCopyDelimiter()
	p.parseCopyOptWith()
	options := p.parseCopyOptions()
	whereClause, _ := p.parseWhereClause()
	stmt := &nodes.CopyStmt{
		Relation:    rel,
		Attlist:     attlist,
		IsFrom:      isFrom,
		IsProgram:   isProgram,
		Filename:    filename,
		WhereClause: whereClause,
	}
	var opts *nodes.List
	if optBin != nil {
		opts = appendToList(opts, optBin)
	}
	if delimOpt != nil {
		opts = appendToList(opts, delimOpt)
	}
	opts = concatNodeLists(opts, options)
	stmt.Options = opts
	return stmt, nil
}

// parseCopyQueryStmt parses COPY '(' PreparableStmt ')' TO ...
func (p *Parser) parseCopyQueryStmt() (nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	query, err := p.parsePreparableStmt()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	isProgram := p.parseCopyOptProgram()
	filename := p.parseCopyFileName()
	p.parseCopyOptWith()
	options := p.parseCopyOptions()
	return &nodes.CopyStmt{
		Query:     query,
		IsFrom:    false,
		IsProgram: isProgram,
		Filename:  filename,
		Options:   options,
	}, nil
}

// parsePreparableStmt parses SELECT, INSERT, UPDATE, or DELETE.
func (p *Parser) parsePreparableStmt() (nodes.Node, error) {
	switch p.cur.Type {
	case SELECT, VALUES, TABLE, WITH:
		return p.parseSelectNoParens()
	case INSERT:
		return p.parseInsertStmt(nil)
	case UPDATE:
		return p.parseUpdateStmt(nil)
	case DELETE_P:
		return p.parseDeleteStmt(nil)
	case MERGE:
		return p.parseMergeStmt(nil)
	default:
		return p.parseSelectNoParens()
	}
}

// parseCopyFrom: FROM -> true, TO -> false
func (p *Parser) parseCopyFrom() bool {
	if p.cur.Type == FROM {
		p.advance()
		return true
	}
	p.expect(TO)
	return false
}

// parseCopyOptProgram: PROGRAM -> true, EMPTY -> false
func (p *Parser) parseCopyOptProgram() bool {
	if p.cur.Type == PROGRAM {
		p.advance()
		return true
	}
	return false
}

// parseCopyFileName: Sconst | STDIN | STDOUT
func (p *Parser) parseCopyFileName() string {
	if p.cur.Type == SCONST {
		s := p.cur.Str
		p.advance()
		return s
	}
	if p.cur.Type == STDIN || p.cur.Type == STDOUT {
		p.advance()
		return ""
	}
	return ""
}

// parseCopyOptWith: WITH | EMPTY
func (p *Parser) parseCopyOptWith() {
	if p.cur.Type == WITH {
		p.advance()
	}
}

// parseCopyOptions: copy_opt_list | '(' utility_option_list ')'
func (p *Parser) parseCopyOptions() *nodes.List {
	if p.cur.Type == '(' {
		p.advance()
		opts := p.parseUtilityOptionList()
		p.expect(')')
		return opts
	}
	return p.parseCopyOptList()
}

// parseCopyOptList: copy_opt_list copy_opt_item | EMPTY
func (p *Parser) parseCopyOptList() *nodes.List {
	var items []nodes.Node
	for {
		item := p.parseCopyOptItem()
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseCopyOptItem parses old-style COPY option keywords.
func (p *Parser) parseCopyOptItem() *nodes.DefElem {
	switch p.cur.Type {
	case BINARY:
		p.advance()
		return &nodes.DefElem{Defname: "format", Arg: &nodes.String{Str: "binary"}, Loc: nodes.NoLoc()}
	case FREEZE:
		p.advance()
		return &nodes.DefElem{Defname: "freeze", Arg: &nodes.Boolean{Boolval: true}, Loc: nodes.NoLoc()}
	case DELIMITER:
		p.advance()
		p.parseOptAs()
		s := p.cur.Str
		p.expect(SCONST)
		return &nodes.DefElem{Defname: "delimiter", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}
	case NULL_P:
		p.advance()
		p.parseOptAs()
		s := p.cur.Str
		p.expect(SCONST)
		return &nodes.DefElem{Defname: "null", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}
	case CSV:
		p.advance()
		return &nodes.DefElem{Defname: "format", Arg: &nodes.String{Str: "csv"}, Loc: nodes.NoLoc()}
	case HEADER_P:
		p.advance()
		return &nodes.DefElem{Defname: "header", Arg: &nodes.Boolean{Boolval: true}, Loc: nodes.NoLoc()}
	case QUOTE:
		p.advance()
		p.parseOptAs()
		s := p.cur.Str
		p.expect(SCONST)
		return &nodes.DefElem{Defname: "quote", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}
	case ESCAPE:
		p.advance()
		p.parseOptAs()
		s := p.cur.Str
		p.expect(SCONST)
		return &nodes.DefElem{Defname: "escape", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}
	case FORCE:
		p.advance()
		if p.cur.Type == QUOTE {
			p.advance()
			if p.cur.Type == '*' {
				starLoc := p.pos()
				p.advance()
				return &nodes.DefElem{Defname: "force_quote", Arg: &nodes.A_Star{Loc: nodes.Loc{Start: starLoc, End: p.prev.End}}, Loc: nodes.NoLoc()}
			}
			cols := p.parseColumnList()
			return &nodes.DefElem{Defname: "force_quote", Arg: cols, Loc: nodes.NoLoc()}
		}
		if p.cur.Type == NOT {
			p.advance()
			p.expect(NULL_P)
			if p.cur.Type == '*' {
				starLoc := p.pos()
				p.advance()
				return &nodes.DefElem{Defname: "force_not_null", Arg: &nodes.A_Star{Loc: nodes.Loc{Start: starLoc, End: p.prev.End}}, Loc: nodes.NoLoc()}
			}
			cols := p.parseColumnList()
			return &nodes.DefElem{Defname: "force_not_null", Arg: cols, Loc: nodes.NoLoc()}
		}
		if p.cur.Type == NULL_P {
			p.advance()
			if p.cur.Type == '*' {
				starLoc := p.pos()
				p.advance()
				return &nodes.DefElem{Defname: "force_null", Arg: &nodes.A_Star{Loc: nodes.Loc{Start: starLoc, End: p.prev.End}}, Loc: nodes.NoLoc()}
			}
			cols := p.parseColumnList()
			return &nodes.DefElem{Defname: "force_null", Arg: cols, Loc: nodes.NoLoc()}
		}
		return nil
	case ENCODING:
		p.advance()
		s := p.cur.Str
		p.expect(SCONST)
		return &nodes.DefElem{Defname: "encoding", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}
	default:
		return nil
	}
}

// parseCopyOptBinary: BINARY | EMPTY
func (p *Parser) parseCopyOptBinary() *nodes.DefElem {
	if p.cur.Type == BINARY {
		p.advance()
		return &nodes.DefElem{Defname: "format", Arg: &nodes.String{Str: "binary"}, Loc: nodes.NoLoc()}
	}
	return nil
}

// parseCopyDelimiter: opt_using DELIMITERS Sconst | EMPTY
func (p *Parser) parseCopyDelimiter() *nodes.DefElem {
	if p.cur.Type == USING {
		p.advance()
	}
	if p.cur.Type == DELIMITERS {
		p.advance()
		s := p.cur.Str
		p.expect(SCONST)
		return &nodes.DefElem{Defname: "delimiter", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}
	}
	return nil
}

// parseOptAs: AS | EMPTY
func (p *Parser) parseOptAs() {
	if p.cur.Type == AS {
		p.advance()
	}
}

// appendToList appends a node to a list, creating the list if nil.
func appendToList(l *nodes.List, n nodes.Node) *nodes.List {
	if l == nil {
		return &nodes.List{Items: []nodes.Node{n}}
	}
	l.Items = append(l.Items, n)
	return l
}

// concatNodeLists concatenates two lists. Returns nil if both are nil.
func concatNodeLists(a, b *nodes.List) *nodes.List {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	a.Items = append(a.Items, b.Items...)
	return a
}
