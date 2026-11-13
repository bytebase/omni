package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCopyStmt parses a COPY statement.
// The COPY keyword has already been consumed.
//
// CopyStmt:
//
//	COPY opt_binary qualified_name opt_column_list
//	    copy_from opt_program copy_file_name copy_delimiter opt_with
//	    copy_options where_clause
//	| COPY '(' PreparableStmt ')' TO opt_program copy_file_name opt_with copy_options
//
// gram.y's actions raise two more syntax errors once the statement is read:
// STDIN or STDOUT after PROGRAM, and a WHERE clause on COPY TO.
func (p *Parser) parseCopyStmt() (nodes.Node, error) {
	loc := p.prev.Loc // COPY already consumed by caller
	if p.cur.Type == '(' {
		return p.parseCopyQueryStmt(loc)
	}
	optBin := p.parseCopyOptBinary()
	relLoc := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rel := makeRangeVarFromNames(names, nodes.Loc{Start: relLoc, End: p.prev.End})
	attlist, err := p.parseOptColumnList()
	if err != nil {
		return nil, err
	}
	isFrom, err := p.parseCopyFrom()
	if err != nil {
		return nil, err
	}
	isProgram := p.parseCopyOptProgram()
	fileTok := p.cur
	filename, err := p.parseCopyFileName()
	if err != nil {
		return nil, err
	}
	delimOpt, err := p.parseCopyDelimiter()
	if err != nil {
		return nil, err
	}
	p.parseCopyOptWith()
	options, err := p.parseCopyOptions()
	if err != nil {
		return nil, err
	}
	whereLoc := p.pos()
	whereClause, err := p.parseWhereClause()
	if err != nil {
		return nil, err
	}
	if isProgram && fileTok.Type != SCONST {
		return nil, &ParseError{Message: "STDIN/STDOUT not allowed with PROGRAM", Position: fileTok.Loc}
	}
	if !isFrom && whereClause != nil {
		return nil, &ParseError{Message: "WHERE clause not allowed with COPY TO", Position: whereLoc}
	}
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
	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// parseCopyQueryStmt parses COPY '(' PreparableStmt ')' TO ...
func (p *Parser) parseCopyQueryStmt(loc int) (nodes.Node, error) {
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
	toLoc := p.pos()
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	isProgram := p.parseCopyOptProgram()
	fileTok := p.cur
	filename, err := p.parseCopyFileName()
	if err != nil {
		return nil, err
	}
	p.parseCopyOptWith()
	options, err := p.parseCopyOptions()
	if err != nil {
		return nil, err
	}
	if isProgram && fileTok.Type != SCONST {
		// gram.y reports this one at TO (@5), not at the file name.
		return nil, &ParseError{Message: "STDIN/STDOUT not allowed with PROGRAM", Position: toLoc}
	}
	return &nodes.CopyStmt{
		Query:     query,
		IsFrom:    false,
		IsProgram: isProgram,
		Filename:  filename,
		Options:   options,
		Loc:       nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parsePreparableStmt parses PreparableStmt.
//
//	PreparableStmt: SelectStmt | InsertStmt | UpdateStmt | DeleteStmt | MergeStmt
//
// Each alternative is a whole statement, so a token that starts none of
// them is a syntax error.
func (p *Parser) parsePreparableStmt() (nodes.Node, error) {
	switch p.cur.Type {
	case INSERT:
		return p.parseInsertStmt(nil)
	case UPDATE:
		return p.parseUpdateStmt(nil)
	case DELETE_P:
		return p.parseDeleteStmt(nil)
	case MERGE:
		return p.parseMergeStmt(nil)
	case WITH:
		// opt_with_clause leads InsertStmt, UpdateStmt, DeleteStmt and
		// MergeStmt as well as a SelectStmt. A CTE's body is a
		// PreparableStmt, so this recursion cycle passes through here.
		if err := p.enterNesting(); err != nil {
			return nil, err
		}
		defer p.leaveNesting()
		stmt, err := p.parseWithStmt()
		if err != nil {
			return nil, err
		}
		if stmt == nil && !p.collectMode() {
			return nil, p.syntaxErrorAtCur()
		}
		return stmt, nil
	}
	stmt, err := p.parseRequiredSelectNoParens()
	if err != nil || stmt == nil {
		return nil, err
	}
	return stmt, nil
}

// parseRequiredSelectNoParens parses a SelectStmt where gram.y requires one
// (AS SelectStmt, FOR SelectStmt, PreparableStmt) and fails at the current
// token when no select_no_parens starts there. In collect mode it returns
// what parseSelectNoParens did, so completion goes on to offer what
// follows.
func (p *Parser) parseRequiredSelectNoParens() (nodes.Node, error) {
	stmt, err := p.parseSelectNoParens()
	if err != nil {
		return nil, err
	}
	if stmt == nil {
		// A nil *SelectStmt must not become a non-nil Node.
		if p.collectMode() {
			return nil, nil
		}
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// parseCopyFrom parses copy_from: FROM -> true, TO -> false.
func (p *Parser) parseCopyFrom() (bool, error) {
	if p.cur.Type == FROM {
		p.advance()
		return true, nil
	}
	if _, err := p.expect(TO); err != nil {
		return false, err
	}
	return false, nil
}

// parseCopyOptProgram: PROGRAM -> true, EMPTY -> false
func (p *Parser) parseCopyOptProgram() bool {
	if p.cur.Type == PROGRAM {
		p.advance()
		return true
	}
	return false
}

// parseCopyFileName parses copy_file_name. STDIN and STDOUT name no file
// and return "".
//
//	copy_file_name: Sconst | STDIN | STDOUT
func (p *Parser) parseCopyFileName() (string, error) {
	switch p.cur.Type {
	case SCONST:
		s := p.cur.Str
		p.advance()
		return s, nil
	case STDIN, STDOUT:
		p.advance()
		return "", nil
	}
	if p.collectMode() {
		p.addTokenCandidate(SCONST)
		p.addTokenCandidate(STDIN)
		p.addTokenCandidate(STDOUT)
		return "", errCollecting
	}
	return "", p.syntaxErrorAtCur()
}

// parseCopyOptWith: WITH | EMPTY
func (p *Parser) parseCopyOptWith() {
	if p.cur.Type == WITH {
		p.advance()
	}
}

// parseCopyOptions: copy_opt_list | '(' utility_option_list ')'
func (p *Parser) parseCopyOptions() (*nodes.List, error) {
	if p.cur.Type == '(' {
		p.advance()
		opts, err := p.parseCopyGenericOptList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return opts, nil
	}
	return p.parseCopyOptList()
}

// parseCopyGenericOptList parses the options of COPY's parenthesized
// syntax.
//
//	copy_generic_opt_list: copy_generic_opt_elem | copy_generic_opt_list ',' copy_generic_opt_elem
//	copy_generic_opt_elem: ColLabel copy_generic_opt_arg
func (p *Parser) parseCopyGenericOptList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		name, err := p.parseColLabel()
		if err != nil {
			return nil, err
		}
		arg, err := p.parseCopyGenericOptArg()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.DefElem{Defname: name, Arg: arg, Loc: nodes.NoLoc()})
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, nil
		}
		p.advance()
	}
}

// parseCopyGenericOptArg parses copy_generic_opt_arg.
//
//	copy_generic_opt_arg:
//	    opt_boolean_or_string | NumericOnly | '*' | DEFAULT
//	    | '(' copy_generic_opt_arg_list ')' | /* EMPTY */
//	copy_generic_opt_arg_list:
//	    opt_boolean_or_string | copy_generic_opt_arg_list ',' opt_boolean_or_string
func (p *Parser) parseCopyGenericOptArg() (nodes.Node, error) {
	switch p.cur.Type {
	case '*':
		starLoc := p.pos()
		p.advance()
		return &nodes.A_Star{Loc: nodes.Loc{Start: starLoc, End: p.prev.End}}, nil
	case DEFAULT:
		p.advance()
		return &nodes.String{Str: "default"}, nil
	case '(':
		p.advance()
		var items []nodes.Node
		for {
			s, err := p.parseBooleanOrString()
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: s})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.List{Items: items}, nil
	}
	return p.parseUtilityOptionArg()
}

// parseCopyOptList: copy_opt_list copy_opt_item | EMPTY
func (p *Parser) parseCopyOptList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		item, err := p.parseCopyOptItem()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseCopyOptItem parses old-style COPY option keywords.
func (p *Parser) parseCopyOptItem() (*nodes.DefElem, error) {
	switch p.cur.Type {
	case BINARY:
		p.advance()
		return &nodes.DefElem{Defname: "format", Arg: &nodes.String{Str: "binary"}, Loc: nodes.NoLoc()}, nil
	case FREEZE:
		p.advance()
		return &nodes.DefElem{Defname: "freeze", Arg: &nodes.Boolean{Boolval: true}, Loc: nodes.NoLoc()}, nil
	case DELIMITER:
		p.advance()
		p.parseOptAs()
		s := p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: "delimiter", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}, nil
	case NULL_P:
		p.advance()
		p.parseOptAs()
		s := p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: "null", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}, nil
	case CSV:
		p.advance()
		return &nodes.DefElem{Defname: "format", Arg: &nodes.String{Str: "csv"}, Loc: nodes.NoLoc()}, nil
	case HEADER_P:
		p.advance()
		return &nodes.DefElem{Defname: "header", Arg: &nodes.Boolean{Boolval: true}, Loc: nodes.NoLoc()}, nil
	case QUOTE:
		p.advance()
		p.parseOptAs()
		s := p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: "quote", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}, nil
	case ESCAPE:
		p.advance()
		p.parseOptAs()
		s := p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: "escape", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}, nil
	case FORCE:
		p.advance()
		if p.cur.Type == QUOTE {
			p.advance()
			if p.cur.Type == '*' {
				starLoc := p.pos()
				p.advance()
				return &nodes.DefElem{Defname: "force_quote", Arg: &nodes.A_Star{Loc: nodes.Loc{Start: starLoc, End: p.prev.End}}, Loc: nodes.NoLoc()}, nil
			}
			cols, err := p.parseColumnList()
			if err != nil {
				return nil, err
			}
			return &nodes.DefElem{Defname: "force_quote", Arg: cols, Loc: nodes.NoLoc()}, nil
		}
		if p.cur.Type == NOT {
			p.advance()
			if _, err := p.expect(NULL_P); err != nil {
				return nil, err
			}
			if p.cur.Type == '*' {
				starLoc := p.pos()
				p.advance()
				return &nodes.DefElem{Defname: "force_not_null", Arg: &nodes.A_Star{Loc: nodes.Loc{Start: starLoc, End: p.prev.End}}, Loc: nodes.NoLoc()}, nil
			}
			cols, err := p.parseColumnList()
			if err != nil {
				return nil, err
			}
			return &nodes.DefElem{Defname: "force_not_null", Arg: cols, Loc: nodes.NoLoc()}, nil
		}
		if p.cur.Type == NULL_P {
			p.advance()
			if p.cur.Type == '*' {
				starLoc := p.pos()
				p.advance()
				return &nodes.DefElem{Defname: "force_null", Arg: &nodes.A_Star{Loc: nodes.Loc{Start: starLoc, End: p.prev.End}}, Loc: nodes.NoLoc()}, nil
			}
			cols, err := p.parseColumnList()
			if err != nil {
				return nil, err
			}
			return &nodes.DefElem{Defname: "force_null", Arg: cols, Loc: nodes.NoLoc()}, nil
		}
		return nil, nil
	case ENCODING:
		p.advance()
		s := p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: "encoding", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}, nil
	// optional-probe: copy_opt_list terminator — parseCopyOptList breaks on nil
	default:
		return nil, nil
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

// parseCopyDelimiter parses copy_delimiter.
//
//	copy_delimiter: opt_using DELIMITERS Sconst | /* EMPTY */
func (p *Parser) parseCopyDelimiter() (*nodes.DefElem, error) {
	if p.cur.Type != USING && p.cur.Type != DELIMITERS {
		return nil, nil
	}
	if p.cur.Type == USING {
		p.advance()
	}
	if _, err := p.expect(DELIMITERS); err != nil {
		return nil, err
	}
	s := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	return &nodes.DefElem{Defname: "delimiter", Arg: &nodes.String{Str: s}, Loc: nodes.NoLoc()}, nil
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
