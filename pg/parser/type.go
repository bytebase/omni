package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseTypename parses a Typename production.
//
// Ref: https://www.postgresql.org/docs/17/sql-expressions.html#SQL-SYNTAX-TYPE-CASTS
//
//	Typename:
//	    SimpleTypename opt_array_bounds
//	    | SETOF SimpleTypename opt_array_bounds
//	    | SimpleTypename ARRAY '[' Iconst ']'
//	    | SETOF SimpleTypename ARRAY '[' Iconst ']'
//	    | SimpleTypename ARRAY
//	    | SETOF SimpleTypename ARRAY
func (p *Parser) parseTypename() (*nodes.TypeName, error) {
	setof := false
	if _, ok := p.match(SETOF); ok {
		setof = true
	}

	tn, err := p.parseSimpleTypename()
	if err != nil {
		return nil, err
	}
	tn.Location = -1
	tn.Setof = setof

	// Check for ARRAY '[' Iconst ']' or ARRAY (without bounds)
	if p.cur.Type == ARRAY {
		p.advance()
		if p.cur.Type == '[' {
			p.advance()
			if p.cur.Type != ICONST {
				return nil, &ParseError{Message: "expected integer constant for array bound", Position: p.cur.Loc}
			}
			ival := p.cur.Ival
			p.advance()
			if _, err := p.expect(']'); err != nil {
				return nil, err
			}
			tn.ArrayBounds = &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: ival}}}
		} else {
			tn.ArrayBounds = &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: -1}}}
		}
	} else {
		// opt_array_bounds: zero or more '[ ]' or '[ Iconst ]'
		bounds := p.parseOptArrayBounds()
		if bounds != nil {
			tn.ArrayBounds = bounds
		}
	}

	return tn, nil
}

// parseOptArrayBounds parses zero or more array bound specifications.
//
//	opt_array_bounds:
//	    opt_array_bounds '[' ']'
//	    | opt_array_bounds '[' Iconst ']'
//	    | /* EMPTY */
func (p *Parser) parseOptArrayBounds() *nodes.List {
	var items []nodes.Node
	for p.cur.Type == '[' {
		p.advance()
		if p.cur.Type == ICONST {
			ival := p.cur.Ival
			p.advance()
			if _, err := p.expect(']'); err != nil {
				return nil
			}
			items = append(items, &nodes.Integer{Ival: ival})
		} else if p.cur.Type == ']' {
			p.advance()
			items = append(items, &nodes.Integer{Ival: -1})
		} else {
			// Not an array bound; caller will handle
			return nil
		}
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseSimpleTypename parses a SimpleTypename production.
//
//	SimpleTypename:
//	    GenericType
//	    | Numeric
//	    | Bit
//	    | Character
//	    | ConstDatetime
//	    | ConstInterval opt_interval
//	    | ConstInterval '(' Iconst ')'
//	    | BOOLEAN_P
//	    | JSON
func (p *Parser) parseSimpleTypename() (*nodes.TypeName, error) {
	switch p.cur.Type {
	case INT_P, INTEGER:
		p.advance()
		return makeTypeName("int4"), nil
	case SMALLINT:
		p.advance()
		return makeTypeName("int2"), nil
	case BIGINT:
		p.advance()
		return makeTypeName("int8"), nil
	case REAL:
		p.advance()
		return makeTypeName("float4"), nil
	case FLOAT_P:
		p.advance()
		return p.parseOptFloat()
	case DOUBLE_P:
		p.advance()
		if _, err := p.expect(PRECISION); err != nil {
			return nil, err
		}
		return makeTypeName("float8"), nil
	case DECIMAL_P:
		p.advance()
		tn := makeTypeName("numeric")
		tn.Typmods = p.parseOptTypeModifiers()
		return tn, nil
	case DEC:
		p.advance()
		tn := makeTypeName("numeric")
		tn.Typmods = p.parseOptTypeModifiers()
		return tn, nil
	case NUMERIC:
		p.advance()
		tn := makeTypeName("numeric")
		tn.Typmods = p.parseOptTypeModifiers()
		return tn, nil
	case BIT:
		return p.parseBit()
	case CHARACTER:
		return p.parseCharacterType()
	case CHAR_P:
		return p.parseCharType()
	case VARCHAR:
		return p.parseVarcharType()
	case NATIONAL:
		return p.parseNationalCharType()
	case NCHAR:
		return p.parseNcharType()
	case BOOLEAN_P:
		p.advance()
		return makeTypeName("bool"), nil
	case JSON:
		p.advance()
		return makeTypeName("json"), nil
	case TIMESTAMP:
		return p.parseTimestampType()
	case TIME:
		return p.parseTimeType()
	case INTERVAL:
		return p.parseIntervalType()
	default:
		// GenericType: type_function_name opt_type_modifiers
		//            | type_function_name '.' attr_name opt_type_modifiers
		return p.parseGenericType()
	}
}

// parseGenericType parses a GenericType production.
//
//	GenericType:
//	    type_function_name opt_type_modifiers
//	    | type_function_name '.' attr_name opt_type_modifiers
func (p *Parser) parseGenericType() (*nodes.TypeName, error) {
	name, err := p.parseTypeFunctionName()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == '.' {
		p.advance()
		attr, err := p.parseAttrName()
		if err != nil {
			return nil, err
		}
		typmods := p.parseOptTypeModifiers()
		return &nodes.TypeName{
			Names: &nodes.List{Items: []nodes.Node{
				&nodes.String{Str: name},
				&nodes.String{Str: attr},
			}},
			Typmods:  typmods,
			Location: -1,
		}, nil
	}

	typmods := p.parseOptTypeModifiers()
	return &nodes.TypeName{
		Names:    &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}},
		Typmods:  typmods,
		Location: -1,
	}, nil
}

// parseOptTypeModifiers parses optional type modifiers.
//
//	opt_type_modifiers:
//	    '(' expr_list ')'
//	    | /* EMPTY */
func (p *Parser) parseOptTypeModifiers() *nodes.List {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	exprs := p.parseExprList()
	if _, err := p.expect(')'); err != nil {
		return nil
	}
	return exprs
}

// parseExprList parses a comma-separated list of expressions.
// This is a minimal implementation for type modifiers (integer constants).
// The full implementation will be in batch 3 (expressions).
func (p *Parser) parseExprList() *nodes.List {
	first := p.parseAExprForTypmod()
	if first == nil {
		return nil
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		expr := p.parseAExprForTypmod()
		if expr == nil {
			break
		}
		items = append(items, expr)
	}
	return &nodes.List{Items: items}
}

// parseAExprForTypmod parses a simple expression suitable for type modifiers.
// For now, this only handles integer constants and signed integers.
// The full a_expr will be implemented in batch 3.
func (p *Parser) parseAExprForTypmod() nodes.Node {
	if p.cur.Type == ICONST {
		tok := p.advance()
		return &nodes.A_Const{
			Val:      &nodes.Integer{Ival: tok.Ival},
			Location: nodes.ParseLoc(tok.Loc),
		}
	}
	if p.cur.Type == FCONST {
		tok := p.advance()
		return &nodes.A_Const{
			Val:      &nodes.Float{Fval: tok.Str},
			Location: nodes.ParseLoc(tok.Loc),
		}
	}
	if p.cur.Type == SCONST {
		tok := p.advance()
		return &nodes.A_Const{
			Val:      &nodes.String{Str: tok.Str},
			Location: nodes.ParseLoc(tok.Loc),
		}
	}
	// Fall back to full a_expr for complex expressions
	return p.parseAExpr(0)
}

// parseOptFloat parses opt_float after FLOAT keyword.
//
//	opt_float:
//	    '(' Iconst ')'
//	    | /* EMPTY */
func (p *Parser) parseOptFloat() (*nodes.TypeName, error) {
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ICONST {
			return nil, &ParseError{Message: "expected integer constant", Position: p.cur.Loc}
		}
		prec := p.cur.Ival
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if prec <= 24 {
			return makeTypeName("float4"), nil
		}
		return makeTypeName("float8"), nil
	}
	return makeTypeName("float8"), nil
}

// parseBit parses Bit type.
//
//	Bit:
//	    BIT opt_varying '(' expr_list ')'   -- BitWithLength
//	    | BIT opt_varying                    -- BitWithoutLength
func (p *Parser) parseBit() (*nodes.TypeName, error) {
	p.advance() // consume BIT
	varying := p.parseOptVarying()

	if p.cur.Type == '(' {
		// BitWithLength
		p.advance()
		exprs := p.parseExprList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if varying {
			tn := makeTypeName("varbit")
			tn.Typmods = exprs
			return tn, nil
		}
		tn := makeTypeName("bit")
		tn.Typmods = exprs
		return tn, nil
	}

	// BitWithoutLength
	if varying {
		return makeTypeName("varbit"), nil
	}
	tn := makeTypeName("bit")
	tn.Typmods = &nodes.List{Items: []nodes.Node{makeIntConst(1)}}
	return tn, nil
}

// parseOptVarying parses opt_varying.
//
//	opt_varying:
//	    VARYING
//	    | /* EMPTY */
func (p *Parser) parseOptVarying() bool {
	if _, ok := p.match(VARYING); ok {
		return true
	}
	return false
}

// parseCharacterType parses CHARACTER [VARYING] type.
func (p *Parser) parseCharacterType() (*nodes.TypeName, error) {
	p.advance() // consume CHARACTER
	varying := p.parseOptVarying()
	return p.finishCharType(varying)
}

// parseCharType parses CHAR [VARYING] type.
func (p *Parser) parseCharType() (*nodes.TypeName, error) {
	p.advance() // consume CHAR
	varying := p.parseOptVarying()
	return p.finishCharType(varying)
}

// parseVarcharType parses VARCHAR type.
func (p *Parser) parseVarcharType() (*nodes.TypeName, error) {
	p.advance() // consume VARCHAR
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ICONST {
			return nil, &ParseError{Message: "expected integer constant", Position: p.cur.Loc}
		}
		ival := p.cur.Ival
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		tn := makeTypeName("varchar")
		tn.Typmods = &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: ival}}}
		return tn, nil
	}
	return makeTypeName("varchar"), nil
}

// parseNationalCharType parses NATIONAL CHARACTER/CHAR type.
func (p *Parser) parseNationalCharType() (*nodes.TypeName, error) {
	p.advance() // consume NATIONAL
	if p.cur.Type == CHARACTER {
		p.advance()
	} else if p.cur.Type == CHAR_P {
		p.advance()
	} else {
		return nil, &ParseError{Message: "expected CHARACTER or CHAR after NATIONAL", Position: p.cur.Loc}
	}
	varying := p.parseOptVarying()
	return p.finishCharType(varying)
}

// parseNcharType parses NCHAR type.
func (p *Parser) parseNcharType() (*nodes.TypeName, error) {
	p.advance() // consume NCHAR
	varying := p.parseOptVarying()
	return p.finishCharType(varying)
}

// finishCharType completes parsing a character type with optional length.
func (p *Parser) finishCharType(varying bool) (*nodes.TypeName, error) {
	typName := "bpchar"
	if varying {
		typName = "varchar"
	}

	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ICONST {
			return nil, &ParseError{Message: "expected integer constant", Position: p.cur.Loc}
		}
		ival := p.cur.Ival
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		tn := makeTypeName(typName)
		tn.Typmods = &nodes.List{Items: []nodes.Node{&nodes.Integer{Ival: ival}}}
		return tn, nil
	}
	return makeTypeName(typName), nil
}

// parseTimestampType parses TIMESTAMP type.
//
//	TIMESTAMP '(' Iconst ')' opt_timezone
//	| TIMESTAMP opt_timezone
func (p *Parser) parseTimestampType() (*nodes.TypeName, error) {
	p.advance() // consume TIMESTAMP
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ICONST {
			return nil, &ParseError{Message: "expected integer constant", Position: p.cur.Loc}
		}
		ival := p.cur.Ival
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		tz := p.parseOptTimezone()
		if tz {
			tn := makeTypeName("timestamptz")
			tn.Typmods = &nodes.List{Items: []nodes.Node{makeIntConst(ival)}}
			return tn, nil
		}
		tn := makeTypeName("timestamp")
		tn.Typmods = &nodes.List{Items: []nodes.Node{makeIntConst(ival)}}
		return tn, nil
	}
	tz := p.parseOptTimezone()
	if tz {
		return makeTypeName("timestamptz"), nil
	}
	return makeTypeName("timestamp"), nil
}

// parseTimeType parses TIME type.
//
//	TIME '(' Iconst ')' opt_timezone
//	| TIME opt_timezone
func (p *Parser) parseTimeType() (*nodes.TypeName, error) {
	p.advance() // consume TIME
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ICONST {
			return nil, &ParseError{Message: "expected integer constant", Position: p.cur.Loc}
		}
		ival := p.cur.Ival
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		tz := p.parseOptTimezone()
		if tz {
			tn := makeTypeName("timetz")
			tn.Typmods = &nodes.List{Items: []nodes.Node{makeIntConst(ival)}}
			return tn, nil
		}
		tn := makeTypeName("time")
		tn.Typmods = &nodes.List{Items: []nodes.Node{makeIntConst(ival)}}
		return tn, nil
	}
	tz := p.parseOptTimezone()
	if tz {
		return makeTypeName("timetz"), nil
	}
	return makeTypeName("time"), nil
}

// parseOptTimezone parses opt_timezone.
//
//	opt_timezone:
//	    WITH_LA TIME ZONE
//	    | WITHOUT_LA TIME ZONE
//	    | /* EMPTY */
func (p *Parser) parseOptTimezone() bool {
	if p.cur.Type == WITH_LA {
		p.advance()
		p.expect(TIME)
		p.expect(ZONE)
		return true
	}
	if p.cur.Type == WITHOUT_LA {
		p.advance()
		p.expect(TIME)
		p.expect(ZONE)
		return false
	}
	return false
}

// parseIntervalType parses INTERVAL type with optional qualifiers.
//
//	ConstInterval opt_interval
//	| ConstInterval '(' Iconst ')'
func (p *Parser) parseIntervalType() (*nodes.TypeName, error) {
	p.advance() // consume INTERVAL
	tn := makeTypeName("interval")

	if p.cur.Type == '(' {
		// ConstInterval '(' Iconst ')'
		p.advance()
		if p.cur.Type != ICONST {
			return nil, &ParseError{Message: "expected integer constant", Position: p.cur.Loc}
		}
		ival := p.cur.Ival
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		tn.Typmods = &nodes.List{Items: []nodes.Node{
			makeIntConst(int64(nodes.INTERVAL_FULL_RANGE)),
			makeIntConst(ival),
		}}
		return tn, nil
	}

	// opt_interval
	typmods := p.parseOptInterval()
	if typmods != nil {
		tn.Typmods = typmods
	}
	return tn, nil
}

// parseOptInterval parses opt_interval.
//
//	opt_interval:
//	    YEAR_P | MONTH_P | DAY_P | HOUR_P | MINUTE_P | interval_second
//	    | YEAR_P TO MONTH_P
//	    | DAY_P TO HOUR_P | DAY_P TO MINUTE_P | DAY_P TO interval_second
//	    | HOUR_P TO MINUTE_P | HOUR_P TO interval_second
//	    | MINUTE_P TO interval_second
//	    | /* EMPTY */
func (p *Parser) parseOptInterval() *nodes.List {
	switch p.cur.Type {
	case YEAR_P:
		p.advance()
		if _, ok := p.match(TO); ok {
			p.expect(MONTH_P)
			return &nodes.List{Items: []nodes.Node{
				makeIntConst(int64(nodes.INTERVAL_MASK_YEAR | nodes.INTERVAL_MASK_MONTH)),
			}}
		}
		return &nodes.List{Items: []nodes.Node{
			makeIntConst(int64(nodes.INTERVAL_MASK_YEAR)),
		}}
	case MONTH_P:
		p.advance()
		return &nodes.List{Items: []nodes.Node{
			makeIntConst(int64(nodes.INTERVAL_MASK_MONTH)),
		}}
	case DAY_P:
		p.advance()
		if _, ok := p.match(TO); ok {
			switch p.cur.Type {
			case HOUR_P:
				p.advance()
				return &nodes.List{Items: []nodes.Node{
					makeIntConst(int64(nodes.INTERVAL_MASK_DAY | nodes.INTERVAL_MASK_HOUR)),
				}}
			case MINUTE_P:
				p.advance()
				return &nodes.List{Items: []nodes.Node{
					makeIntConst(int64(nodes.INTERVAL_MASK_DAY | nodes.INTERVAL_MASK_HOUR | nodes.INTERVAL_MASK_MINUTE)),
				}}
			case SECOND_P:
				secs := p.parseIntervalSecond()
				secs.Items[0] = makeIntConst(int64(nodes.INTERVAL_MASK_DAY | nodes.INTERVAL_MASK_HOUR | nodes.INTERVAL_MASK_MINUTE | nodes.INTERVAL_MASK_SECOND))
				return secs
			}
		}
		return &nodes.List{Items: []nodes.Node{
			makeIntConst(int64(nodes.INTERVAL_MASK_DAY)),
		}}
	case HOUR_P:
		p.advance()
		if _, ok := p.match(TO); ok {
			switch p.cur.Type {
			case MINUTE_P:
				p.advance()
				return &nodes.List{Items: []nodes.Node{
					makeIntConst(int64(nodes.INTERVAL_MASK_HOUR | nodes.INTERVAL_MASK_MINUTE)),
				}}
			case SECOND_P:
				secs := p.parseIntervalSecond()
				secs.Items[0] = makeIntConst(int64(nodes.INTERVAL_MASK_HOUR | nodes.INTERVAL_MASK_MINUTE | nodes.INTERVAL_MASK_SECOND))
				return secs
			}
		}
		return &nodes.List{Items: []nodes.Node{
			makeIntConst(int64(nodes.INTERVAL_MASK_HOUR)),
		}}
	case MINUTE_P:
		p.advance()
		if _, ok := p.match(TO); ok {
			if p.cur.Type == SECOND_P {
				secs := p.parseIntervalSecond()
				secs.Items[0] = makeIntConst(int64(nodes.INTERVAL_MASK_MINUTE | nodes.INTERVAL_MASK_SECOND))
				return secs
			}
		}
		return &nodes.List{Items: []nodes.Node{
			makeIntConst(int64(nodes.INTERVAL_MASK_MINUTE)),
		}}
	case SECOND_P:
		return p.parseIntervalSecond()
	default:
		return nil
	}
}

// parseIntervalSecond parses interval_second.
//
//	interval_second:
//	    SECOND_P
//	    | SECOND_P '(' Iconst ')'
func (p *Parser) parseIntervalSecond() *nodes.List {
	p.advance() // consume SECOND_P
	if p.cur.Type == '(' {
		p.advance()
		ival := p.cur.Ival
		p.advance() // consume ICONST
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{
			makeIntConst(int64(nodes.INTERVAL_MASK_SECOND)),
			makeIntConst(ival),
		}}
	}
	return &nodes.List{Items: []nodes.Node{
		makeIntConst(int64(nodes.INTERVAL_MASK_SECOND)),
	}}
}

// parseTypeList parses a comma-separated list of typenames.
//
//	type_list:
//	    Typename
//	    | type_list ',' Typename
func (p *Parser) parseTypeList() (*nodes.List, error) {
	tn, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{tn}
	for p.cur.Type == ',' {
		p.advance()
		tn, err = p.parseTypename()
		if err != nil {
			return nil, err
		}
		items = append(items, tn)
	}
	return &nodes.List{Items: items}, nil
}

// parseFuncType parses a func_type production.
//
//	func_type:
//	    Typename
//	    | type_function_name attrs '%' TYPE_P
//	    | SETOF type_function_name attrs '%' TYPE_P
func (p *Parser) parseFuncType() (*nodes.TypeName, error) {
	if p.cur.Type == SETOF {
		// SETOF can be:
		// 1. SETOF type_function_name attrs '%' TYPE_P  (pct_type with SETOF)
		// 2. SETOF Typename (regular SETOF type)
		//
		// We check for the %TYPE pattern: SETOF name '.' ... '%' TYPE_P
		// If the token after SETOF is followed by '.', it might be case 1.
		// Otherwise, delegate to parseTypename which handles SETOF.

		// Check if this is the %TYPE pattern: need at least name '.' ... '%'
		next := p.peekNext()
		if p.isTypeFunctionNameToken(next) {
			// Could be either case. We can't easily distinguish without
			// more lookahead. Use heuristic: delegate to parseTypename
			// which handles SETOF Typename correctly.
			// The %TYPE case is very rare in practice.
		}

		// Delegate to parseTypename which handles SETOF
		return p.parseTypename()
	}

	// Non-SETOF: could be type_function_name attrs '%' TYPE_P or Typename.
	// The %TYPE case requires name '.' name '%' TYPE_P pattern.
	// Check if we see name '.' which could indicate qualified name with %TYPE.
	if p.isTypeFunctionName() {
		// Check if this looks like a %TYPE reference (name.name...%TYPE)
		next := p.peekNext()
		if next.Type == '.' {
			// Could be qualified type or %TYPE. We need to try parsing.
			// Save state for backtracking.
			savedCur := p.cur
			savedPrev := p.prev
			savedNext := p.nextBuf
			savedHasNext := p.hasNext

			name, _ := p.parseTypeFunctionName()
			if p.cur.Type == '.' {
				attrs, _ := p.parseAttrs()
				if p.cur.Type == '%' {
					p.advance()
					if p.cur.Type == TYPE_P {
						p.advance()
						nameItems := make([]nodes.Node, 0, 1+len(attrs.Items))
						nameItems = append(nameItems, &nodes.String{Str: name})
						nameItems = append(nameItems, attrs.Items...)
						return &nodes.TypeName{
							Names:   &nodes.List{Items: nameItems},
							PctType: true,
							Location: -1,
						}, nil
					}
				}
			}

			// Not %TYPE pattern - restore and parse as Typename
			p.cur = savedCur
			p.prev = savedPrev
			p.nextBuf = savedNext
			p.hasNext = savedHasNext
		}
	}

	return p.parseTypename()
}

// isTypeFunctionNameToken checks if a token could be a type_function_name.
func (p *Parser) isTypeFunctionNameToken(tok Token) bool {
	if tok.Type == IDENT {
		return true
	}
	if kw := LookupKeyword(tok.Str); kw != nil && kw.Token == tok.Type {
		return kw.Category == UnreservedKeyword || kw.Category == TypeFuncNameKeyword
	}
	return false
}

// makeTypeName creates a TypeName with pg_catalog schema prefix.
func makeTypeName(typName string) *nodes.TypeName {
	return &nodes.TypeName{
		Names: &nodes.List{Items: []nodes.Node{
			&nodes.String{Str: "pg_catalog"},
			&nodes.String{Str: typName},
		}},
		Location: -1,
	}
}

// makeIntConst creates an A_Const with an integer value.
func makeIntConst(val int64) nodes.Node {
	return &nodes.A_Const{Val: &nodes.Integer{Ival: val}}
}
