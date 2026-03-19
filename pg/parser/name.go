package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// Keyword category predicates for identifier parsing.

// isColId returns true if the current token can be used as a ColId.
//
//	ColId: IDENT | unreserved_keyword | col_name_keyword
func (p *Parser) isColId() bool {
	if p.collectMode() {
		p.addTokenCandidate(IDENT)
		p.addKeywordsByCategory(UnreservedKeyword, ColNameKeyword)
		return true
	}
	if p.cur.Type == IDENT {
		return true
	}
	if kw := LookupKeyword(p.cur.Str); kw != nil && kw.Token == p.cur.Type {
		return kw.Category == UnreservedKeyword || kw.Category == ColNameKeyword
	}
	return false
}

// isColLabel returns true if the current token can be used as a ColLabel.
//
//	ColLabel: IDENT | unreserved_keyword | col_name_keyword | type_func_name_keyword | reserved_keyword
func (p *Parser) isColLabel() bool {
	if p.collectMode() {
		p.addTokenCandidate(IDENT)
		p.addKeywordsByCategory(UnreservedKeyword, ColNameKeyword, TypeFuncNameKeyword, ReservedKeyword)
		return true
	}
	if p.cur.Type == IDENT {
		return true
	}
	if kw := LookupKeyword(p.cur.Str); kw != nil && kw.Token == p.cur.Type {
		return true
	}
	return false
}

// isTypeFunctionName returns true if the current token can be used as a type_function_name.
//
//	type_function_name: IDENT | unreserved_keyword | type_func_name_keyword
func (p *Parser) isTypeFunctionName() bool {
	if p.collectMode() {
		p.addTokenCandidate(IDENT)
		p.addKeywordsByCategory(UnreservedKeyword, TypeFuncNameKeyword)
		return true
	}
	if p.cur.Type == IDENT {
		return true
	}
	if kw := LookupKeyword(p.cur.Str); kw != nil && kw.Token == p.cur.Type {
		return kw.Category == UnreservedKeyword || kw.Category == TypeFuncNameKeyword
	}
	return false
}

// parseColId parses a ColId (identifier or unreserved/col_name keyword).
//
//	ColId: IDENT | unreserved_keyword | col_name_keyword
func (p *Parser) parseColId() (string, error) {
	if p.isColId() {
		tok := p.advance()
		return tok.Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseColLabel parses a ColLabel (any identifier or keyword).
//
//	ColLabel: IDENT | unreserved_keyword | col_name_keyword | type_func_name_keyword | reserved_keyword
func (p *Parser) parseColLabel() (string, error) {
	if p.isColLabel() {
		tok := p.advance()
		return tok.Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseName parses a name (same as ColId).
//
//	name: ColId
func (p *Parser) parseName() (string, error) {
	return p.parseColId()
}

// parseAttrName parses an attribute name (same as ColLabel).
//
//	attr_name: ColLabel
func (p *Parser) parseAttrName() (string, error) {
	return p.parseColLabel()
}

// parseTypeFunctionName parses a type_function_name.
//
//	type_function_name: IDENT | unreserved_keyword | type_func_name_keyword
func (p *Parser) parseTypeFunctionName() (string, error) {
	if p.isTypeFunctionName() {
		tok := p.advance()
		return tok.Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseQualifiedName parses a qualified name (up to 3 parts).
//
// Ref: gram.y qualified_name rule
//
//	qualified_name:
//	    ColId
//	    | ColId '.' attr_name
//	    | ColId '.' attr_name '.' attr_name
func (p *Parser) parseQualifiedName() (*nodes.List, error) {
	if p.collectMode() {
		p.addRuleCandidate("qualified_name")
		return nil, errCollecting
	}
	id, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{&nodes.String{Str: id}}

	if p.cur.Type == '.' {
		p.advance()
		// After "name.", we may have reached the cursor position.
		// Emit rule candidates so callers (relation_expr, etc.) are recognized.
		if p.collectMode() {
			p.addRuleCandidate("qualified_name")
		}
		attr, err := p.parseAttrName()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: attr})

		if p.cur.Type == '.' {
			p.advance()
			// Same for second dot in 3-part names.
			if p.collectMode() {
				p.addRuleCandidate("qualified_name")
			}
			attr2, err := p.parseAttrName()
			if err != nil {
				return nil, err
			}
			items = append(items, &nodes.String{Str: attr2})
		}
	}

	return &nodes.List{Items: items}, nil
}

// parseAnyName parses an any_name (dot-separated list of ColId).
//
//	any_name:
//	    ColId
//	    | ColId '.' any_name
func (p *Parser) parseAnyName() (*nodes.List, error) {
	id, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{&nodes.String{Str: id}}

	for p.cur.Type == '.' {
		p.advance()
		// The recursive any_name starts with ColId, but we need to handle
		// the case where the next token after '.' is '*' which could be
		// part of a different rule. Only continue if we see a ColId.
		if !p.isColId() {
			// Put back: this dot is not part of any_name.
			// Since we can't unread, this shouldn't happen in practice
			// as callers of any_name expect well-formed input.
			break
		}
		id2, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: id2})
	}

	return &nodes.List{Items: items}, nil
}

// parseAnyNameList parses a comma-separated list of any_name.
//
//	any_name_list:
//	    any_name
//	    | any_name_list ',' any_name
func (p *Parser) parseAnyNameList() (*nodes.List, error) {
	name, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	result := &nodes.List{Items: []nodes.Node{name}}

	for p.cur.Type == ',' {
		p.advance()
		name, err = p.parseAnyName()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, name)
	}

	return result, nil
}

// parseOptNameList parses an optional name list.
//
//	opt_name_list:
//	    name_list
//	    | /* EMPTY */
func (p *Parser) parseOptNameList() (*nodes.List, error) {
	if p.isColId() {
		return p.parseNameList()
	}
	return nil, nil
}

// parseNameList parses a comma-separated list of names.
//
//	name_list:
//	    name
//	    | name_list ',' name
func (p *Parser) parseNameList() (*nodes.List, error) {
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{&nodes.String{Str: name}}

	for p.cur.Type == ',' {
		p.advance()
		name, err = p.parseName()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: name})
	}

	return &nodes.List{Items: items}, nil
}

// parseColumnRef parses a column reference with optional indirection.
//
// Ref: gram.y columnref rule
//
//	columnref:
//	    ColId
//	    | ColId indirection
func (p *Parser) parseColumnRef() (nodes.Node, error) {
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}
	id, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	// Check for indirection (starts with '.' or '[')
	if p.cur.Type != '.' && p.cur.Type != '[' {
		return &nodes.ColumnRef{
			Fields: &nodes.List{Items: []nodes.Node{&nodes.String{Str: id}}},
		}, nil
	}

	// Parse indirection
	indir, err := p.parseIndirection()
	if err != nil {
		return nil, err
	}

	// Replicate PostgreSQL's makeColumnRef() logic from gram.y:
	// If indirection contains A_Indices (subscripts), split the list.
	cr := &nodes.ColumnRef{}

	// Find first A_Indices in indirection
	firstSubscript := -1
	if indir != nil {
		for i, item := range indir.Items {
			if _, ok := item.(*nodes.A_Indices); ok {
				firstSubscript = i
				break
			}
		}
	}

	if firstSubscript == -1 {
		// No subscripting: all indirection gets added to field list
		fieldItems := []nodes.Node{&nodes.String{Str: id}}
		if indir != nil {
			fieldItems = append(fieldItems, indir.Items...)
		}
		cr.Fields = &nodes.List{Items: fieldItems}
		return cr, nil
	}

	// Split: field selections before subscript go to ColumnRef.Fields,
	// subscript and after go to A_Indirection.Indirection
	ind := &nodes.A_Indirection{}
	if firstSubscript == 0 {
		cr.Fields = &nodes.List{Items: []nodes.Node{&nodes.String{Str: id}}}
		ind.Indirection = indir
	} else {
		fieldItems := make([]nodes.Node, 0, firstSubscript+1)
		fieldItems = append(fieldItems, &nodes.String{Str: id})
		fieldItems = append(fieldItems, indir.Items[:firstSubscript]...)
		cr.Fields = &nodes.List{Items: fieldItems}
		remaining := make([]nodes.Node, len(indir.Items)-firstSubscript)
		copy(remaining, indir.Items[firstSubscript:])
		ind.Indirection = &nodes.List{Items: remaining}
	}
	ind.Arg = cr
	return ind, nil
}

// parseIndirection parses one or more indirection elements.
//
//	indirection:
//	    indirection_el
//	    | indirection indirection_el
func (p *Parser) parseIndirection() (*nodes.List, error) {
	el, err := p.parseIndirectionEl()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{el}

	for p.cur.Type == '.' || p.cur.Type == '[' {
		el, err = p.parseIndirectionEl()
		if err != nil {
			return nil, err
		}
		items = append(items, el)
	}

	return &nodes.List{Items: items}, nil
}

// parseIndirectionEl parses a single indirection element.
//
//	indirection_el:
//	    '.' attr_name
//	    | '.' '*'
//	    | '[' a_expr ']'
//	    | '[' opt_slice_bound ':' opt_slice_bound ']'
func (p *Parser) parseIndirectionEl() (nodes.Node, error) {
	if p.cur.Type == '.' {
		p.advance()
		if p.cur.Type == '*' {
			p.advance()
			return &nodes.A_Star{}, nil
		}
		name, err := p.parseAttrName()
		if err != nil {
			return nil, err
		}
		return &nodes.String{Str: name}, nil
	}

	if p.cur.Type == '[' {
		p.advance()

		// Check for empty lower bound followed by ':'
		if p.cur.Type == ':' {
			// Slice with empty lower bound: [ : opt_slice_bound ]
			p.advance()
			var uidx nodes.Node
			if p.cur.Type != ']' {
				uidx = p.parseAExpr(0)
			}
			if _, err := p.expect(']'); err != nil {
				return nil, err
			}
			return &nodes.A_Indices{IsSlice: true, Uidx: uidx}, nil
		}

		// Parse first expression
		expr := p.parseAExpr(0)

		if p.cur.Type == ':' {
			// Slice: [ expr : opt_slice_bound ]
			p.advance()
			var uidx nodes.Node
			if p.cur.Type != ']' {
				uidx = p.parseAExpr(0)
			}
			if _, err := p.expect(']'); err != nil {
				return nil, err
			}
			return &nodes.A_Indices{IsSlice: true, Lidx: expr, Uidx: uidx}, nil
		}

		// Simple subscript: [ expr ]
		if _, err := p.expect(']'); err != nil {
			return nil, err
		}
		return &nodes.A_Indices{Uidx: expr}, nil
	}

	return nil, p.syntaxErrorAtCur()
}

// parseOptIndirection parses optional indirection.
//
//	opt_indirection:
//	    indirection
//	    | /* EMPTY */
func (p *Parser) parseOptIndirection() (*nodes.List, error) {
	if p.cur.Type == '.' || p.cur.Type == '[' {
		return p.parseIndirection()
	}
	return nil, nil
}

// parseAttrs parses dot-separated attribute names.
//
//	attrs:
//	    '.' attr_name
//	    | attrs '.' attr_name
func (p *Parser) parseAttrs() (*nodes.List, error) {
	if _, err := p.expect('.'); err != nil {
		return nil, err
	}
	name, err := p.parseAttrName()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{&nodes.String{Str: name}}

	for p.cur.Type == '.' {
		p.advance()
		name, err = p.parseAttrName()
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: name})
	}

	return &nodes.List{Items: items}, nil
}

// parseFuncName parses a function name.
//
//	func_name:
//	    type_function_name
//	    | ColId indirection
func (p *Parser) parseFuncName() (*nodes.List, error) {
	if p.collectMode() {
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	// If token is a ColId and indirection follows, parse as ColId + indirection
	if p.isColId() {
		id, _ := p.parseColId()
		if p.cur.Type == '.' || p.cur.Type == '[' {
			indir, err := p.parseIndirection()
			if err != nil {
				return nil, err
			}
			items := []nodes.Node{&nodes.String{Str: id}}
			if indir != nil {
				items = append(items, indir.Items...)
			}
			return &nodes.List{Items: items}, nil
		}
		// No indirection - single name (ColId is a subset of type_function_name for
		// the IDENT and unreserved_keyword cases)
		return &nodes.List{Items: []nodes.Node{&nodes.String{Str: id}}}, nil
	}

	// Try type_function_name (handles type_func_name_keyword not in ColId)
	if p.isTypeFunctionName() {
		name, _ := p.parseTypeFunctionName()
		return &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}, nil
	}

	return nil, p.syntaxErrorAtCur()
}

// parseFileName parses a file name (string constant).
//
//	file_name: Sconst
func (p *Parser) parseFileName() (string, error) {
	if p.cur.Type == SCONST {
		tok := p.advance()
		return tok.Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseCursorName parses a cursor name.
//
//	cursor_name: name
func (p *Parser) parseCursorName() (string, error) {
	return p.parseName()
}

// parseAnyOperator parses a possibly-qualified operator name.
//
//	any_operator:
//	    all_Op
//	    | ColId '.' any_operator
func (p *Parser) parseAnyOperator() (*nodes.List, error) {
	// Check if this could be ColId '.' any_operator
	// We need lookahead: if current is ColId and next is '.', it might be qualified.
	if p.isColId() && p.cur.Type != Op {
		// Save position to check for '.' after ColId
		// Since we can't easily backtrack, use a heuristic:
		// If the token is a ColId (but not an operator), check if next is '.'
		// Actually, all_Op only matches Op and single-char math ops.
		// ColId never matches Op. So if current is ColId, check for '.'.
		// But we need to also check if it could be an all_Op (which it can't if it's ColId).
		// So if current is ColId, we try the qualified path.
		// If current is Op or math op, we do all_Op.

		// Peek: is there a '.' next? We need two-token lookahead.
		// Save current state
		saved := p.cur
		id, _ := p.parseColId()
		if p.cur.Type == '.' {
			p.advance()
			rest, err := p.parseAnyOperator()
			if err != nil {
				return nil, err
			}
			items := make([]nodes.Node, 0, 1+len(rest.Items))
			items = append(items, &nodes.String{Str: id})
			items = append(items, rest.Items...)
			return &nodes.List{Items: items}, nil
		}
		// Not qualified - this is wrong, we consumed a ColId but it's not an operator.
		// Restore (we can't easily backtrack in this parser design).
		// This case shouldn't happen if the caller knows an operator is expected.
		_ = saved
		return nil, p.syntaxErrorAtCur()
	}

	// Parse all_Op
	op, err := p.parseAllOp()
	if err != nil {
		return nil, err
	}
	return &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}}, nil
}

// parseAllOp parses an operator token or math operator.
//
//	all_Op: Op | MathOp
//	MathOp: '+' | '-' | '*' | '/' | '%' | '^' | '<' | '>' | '='
//	        | LESS_EQUALS | GREATER_EQUALS | NOT_EQUALS
func (p *Parser) parseAllOp() (string, error) {
	if op, ok := p.parseMathOp(); ok {
		return op, nil
	}
	if p.cur.Type == Op {
		tok := p.advance()
		return tok.Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseMathOp tries to parse a math operator. Returns the operator string and true if matched.
func (p *Parser) parseMathOp() (string, bool) {
	switch p.cur.Type {
	case '+', '-', '*', '/', '%', '^', '<', '>', '=':
		tok := p.advance()
		return string(rune(tok.Type)), true
	case LESS_EQUALS:
		p.advance()
		return "<=", true
	case GREATER_EQUALS:
		p.advance()
		return ">=", true
	case NOT_EQUALS:
		p.advance()
		return "<>", true
	}
	return "", false
}

// parseQualifiedNameList parses a comma-separated list of qualified names as RangeVars.
//
//	qualified_name_list:
//	    qualified_name
//	    | qualified_name_list ',' qualified_name
func (p *Parser) parseQualifiedNameList() (*nodes.List, error) {
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	result := &nodes.List{Items: []nodes.Node{makeRangeVarFromNames(names)}}

	for p.cur.Type == ',' {
		p.advance()
		names, err = p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, makeRangeVarFromNames(names))
	}

	return result, nil
}

// makeRangeVarFromNames creates a RangeVar from a list of name String nodes.
func makeRangeVarFromNames(names *nodes.List) *nodes.RangeVar {
	rv := &nodes.RangeVar{Inh: true, Relpersistence: 'p', Loc: nodes.NoLoc()}
	if names != nil && len(names.Items) > 0 {
		switch len(names.Items) {
		case 1:
			rv.Relname = names.Items[0].(*nodes.String).Str
		case 2:
			rv.Schemaname = names.Items[0].(*nodes.String).Str
			rv.Relname = names.Items[1].(*nodes.String).Str
		case 3:
			rv.Catalogname = names.Items[0].(*nodes.String).Str
			rv.Schemaname = names.Items[1].(*nodes.String).Str
			rv.Relname = names.Items[2].(*nodes.String).Str
		}
	}
	return rv
}

// makeRangeVarFromAnyName creates a RangeVar from an any_name list with Loc = NoLoc().
func makeRangeVarFromAnyName(names *nodes.List) *nodes.RangeVar {
	rv := &nodes.RangeVar{
		Inh:            true,
		Relpersistence: 'p',
		Loc: nodes.NoLoc(),
	}
	if names == nil {
		return rv
	}
	switch len(names.Items) {
	case 1:
		rv.Relname = names.Items[0].(*nodes.String).Str
	case 2:
		rv.Schemaname = names.Items[0].(*nodes.String).Str
		rv.Relname = names.Items[1].(*nodes.String).Str
	case 3:
		rv.Catalogname = names.Items[0].(*nodes.String).Str
		rv.Schemaname = names.Items[1].(*nodes.String).Str
		rv.Relname = names.Items[2].(*nodes.String).Str
	}
	return rv
}
