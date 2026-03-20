package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// ---------------------------------------------------------------------------
// JSON expression and function parsing
// ---------------------------------------------------------------------------

// parseJsonValueExpr parses a json_value_expr production.
//
// Ref: https://www.postgresql.org/docs/17/functions-json.html
//
//	json_value_expr:
//	    a_expr json_format_clause_opt
func (p *Parser) parseJsonValueExpr() (*nodes.JsonValueExpr, error) {
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	// json_format_clause_opt: FORMAT JSON [ENCODING name] | empty
	// We parse but the yacc grammar discards the format in JsonValueExpr.RawExpr
	_, err = p.parseJsonFormatClauseOpt()
	if err != nil {
		return nil, err
	}
	return &nodes.JsonValueExpr{
		RawExpr: expr,
	}, nil
}

// parseJsonFormatClause parses FORMAT JSON [ENCODING name].
//
//	json_format_clause:
//	    FORMAT JSON
//	    | FORMAT JSON ENCODING name
func (p *Parser) parseJsonFormatClause() (*nodes.JsonFormat, error) {
	p.advance() // consume FORMAT
	if _, err := p.expect(JSON); err != nil {
		return nil, err
	}
	f := &nodes.JsonFormat{
		FormatType: nodes.JS_FORMAT_JSON,
		Loc:        nodes.NoLoc(),
	}
	if p.cur.Type == ENCODING {
		p.advance()
		p.parseColId() // consume encoding name, ignored
	}
	return f, nil
}

// parseJsonFormatClauseOpt parses an optional FORMAT clause.
func (p *Parser) parseJsonFormatClauseOpt() (*nodes.JsonFormat, error) {
	if p.cur.Type == FORMAT {
		return p.parseJsonFormatClause()
	}
	return nil, nil
}

// parseJsonReturningClauseOpt parses an optional RETURNING Typename [FORMAT ...].
//
//	json_returning_clause_opt:
//	    RETURNING Typename json_format_clause_opt
//	    | /* EMPTY */
func (p *Parser) parseJsonReturningClauseOpt() (*nodes.JsonOutput, error) {
	if p.cur.Type != RETURNING {
		return nil, nil
	}
	p.advance() // consume RETURNING
	tn, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.parseJsonFormatClauseOpt(); err != nil {
		return nil, err
	}
	return &nodes.JsonOutput{
		TypeName: tn,
	}, nil
}

// parseJsonBehavior parses a json_behavior production.
//
//	json_behavior:
//	    DEFAULT a_expr
//	    | json_behavior_type
func (p *Parser) parseJsonBehavior() *nodes.JsonBehavior {
	if p.cur.Type == DEFAULT {
		p.advance()
		expr, _ := p.parseAExpr(0)
		return &nodes.JsonBehavior{
			Btype: nodes.JSON_BEHAVIOR_DEFAULT,
			Expr:  expr,
			Loc:   nodes.NoLoc(),
		}
	}
	bt, ok := p.parseJsonBehaviorType()
	if !ok {
		return nil
	}
	return &nodes.JsonBehavior{
		Btype: bt,
		Loc:   nodes.NoLoc(),
	}
}

// parseJsonBehaviorType parses a json_behavior_type.
//
//	json_behavior_type:
//	    ERROR_P | NULL_P | TRUE_P | FALSE_P | UNKNOWN
//	    | EMPTY_P ARRAY | EMPTY_P OBJECT_P | EMPTY_P
func (p *Parser) parseJsonBehaviorType() (nodes.JsonBehaviorType, bool) {
	switch p.cur.Type {
	case ERROR_P:
		p.advance()
		return nodes.JSON_BEHAVIOR_ERROR, true
	case NULL_P:
		p.advance()
		return nodes.JSON_BEHAVIOR_NULL, true
	case TRUE_P:
		p.advance()
		return nodes.JSON_BEHAVIOR_TRUE, true
	case FALSE_P:
		p.advance()
		return nodes.JSON_BEHAVIOR_FALSE, true
	case UNKNOWN:
		p.advance()
		return nodes.JSON_BEHAVIOR_UNKNOWN, true
	case EMPTY_P:
		p.advance()
		if p.cur.Type == ARRAY {
			p.advance()
			return nodes.JSON_BEHAVIOR_EMPTY_ARRAY, true
		}
		if p.cur.Type == OBJECT_P {
			p.advance()
			return nodes.JSON_BEHAVIOR_EMPTY_OBJECT, true
		}
		// bare EMPTY => EMPTY_ARRAY (Oracle compat)
		return nodes.JSON_BEHAVIOR_EMPTY_ARRAY, true
	default:
		return 0, false
	}
}

// isJsonBehaviorStart returns true if the current token can start a json_behavior.
func (p *Parser) isJsonBehaviorStart() bool {
	switch p.cur.Type {
	case DEFAULT, ERROR_P, NULL_P, TRUE_P, FALSE_P, UNKNOWN, EMPTY_P:
		return true
	}
	return false
}

// parseJsonBehaviorClauseOpt parses optional ON EMPTY / ON ERROR behavior clauses.
//
//	json_behavior_clause_opt:
//	    json_behavior ON EMPTY_P json_behavior ON ERROR_P
//	    | json_behavior ON EMPTY_P
//	    | json_behavior ON ERROR_P
//	    | /* EMPTY */
func (p *Parser) parseJsonBehaviorClauseOpt() (onEmpty, onError *nodes.JsonBehavior) {
	if !p.isJsonBehaviorStart() {
		return nil, nil
	}

	first := p.parseJsonBehavior()
	if p.cur.Type != ON {
		return nil, nil
	}
	p.advance() // consume ON

	if p.cur.Type == EMPTY_P {
		p.advance() // consume EMPTY
		onEmpty = first

		// Check for second behavior ON ERROR
		if p.isJsonBehaviorStart() {
			second := p.parseJsonBehavior()
			if p.cur.Type == ON {
				p.advance()
				if p.cur.Type == ERROR_P {
					p.advance()
					onError = second
				}
			}
		}
		return onEmpty, onError
	}

	if p.cur.Type == ERROR_P {
		p.advance() // consume ERROR
		onError = first
		return nil, onError
	}

	return nil, nil
}

// parseJsonOnErrorClauseOpt parses optional json_behavior ON ERROR.
//
//	json_on_error_clause_opt:
//	    json_behavior ON ERROR_P
//	    | /* EMPTY */
func (p *Parser) parseJsonOnErrorClauseOpt() *nodes.JsonBehavior {
	if !p.isJsonBehaviorStart() {
		return nil
	}
	beh := p.parseJsonBehavior()
	if p.cur.Type == ON {
		p.advance()
		if p.cur.Type == ERROR_P {
			p.advance()
			return beh
		}
	}
	return nil
}

// parseJsonWrapperBehavior parses wrapper behavior options.
//
//	json_wrapper_behavior:
//	    WITHOUT [ARRAY] WRAPPER
//	    | WITH [UNCONDITIONAL|CONDITIONAL] [ARRAY] WRAPPER
//	    | /* EMPTY */
func (p *Parser) parseJsonWrapperBehavior() nodes.JsonWrapper {
	if p.cur.Type == WITHOUT {
		p.advance()
		if p.cur.Type == ARRAY {
			p.advance()
		}
		if p.cur.Type == WRAPPER {
			p.advance()
		}
		return nodes.JSW_NONE
	}
	if p.cur.Type == WITH {
		p.advance()
		wrapper := nodes.JSW_UNCONDITIONAL
		if p.cur.Type == CONDITIONAL {
			p.advance()
			wrapper = nodes.JSW_CONDITIONAL
		} else if p.cur.Type == UNCONDITIONAL {
			p.advance()
			// already JSW_UNCONDITIONAL
		}
		if p.cur.Type == ARRAY {
			p.advance()
		}
		if p.cur.Type == WRAPPER {
			p.advance()
		}
		return wrapper
	}
	return nodes.JSW_UNSPEC
}

// parseJsonQuotesClauseOpt parses KEEP/OMIT QUOTES options.
//
//	json_quotes_clause_opt:
//	    KEEP QUOTES [ON SCALAR STRING_P]
//	    | OMIT QUOTES [ON SCALAR STRING_P]
//	    | /* EMPTY */
func (p *Parser) parseJsonQuotesClauseOpt() nodes.JsonQuotes {
	if p.cur.Type == KEEP {
		p.advance()
		if p.cur.Type == QUOTES {
			p.advance()
		}
		if p.cur.Type == ON {
			p.advance()
			if p.cur.Type == SCALAR {
				p.advance()
			}
			if p.cur.Type == STRING_P {
				p.advance()
			}
		}
		return nodes.JS_QUOTES_KEEP
	}
	if p.cur.Type == OMIT {
		p.advance()
		if p.cur.Type == QUOTES {
			p.advance()
		}
		if p.cur.Type == ON {
			p.advance()
			if p.cur.Type == SCALAR {
				p.advance()
			}
			if p.cur.Type == STRING_P {
				p.advance()
			}
		}
		return nodes.JS_QUOTES_OMIT
	}
	return nodes.JS_QUOTES_UNSPEC
}

// parseJsonPredicateTypeConstraint parses IS JSON type constraints.
//
//	json_predicate_type_constraint:
//	    JSON | JSON VALUE_P | JSON ARRAY | JSON OBJECT_P | JSON SCALAR
func (p *Parser) parseJsonPredicateTypeConstraint() (nodes.JsonValueType, bool) {
	if p.cur.Type != JSON {
		return 0, false
	}
	p.advance() // consume JSON
	switch p.cur.Type {
	case VALUE_P:
		p.advance()
		return nodes.JS_TYPE_ANY, true
	case ARRAY:
		p.advance()
		return nodes.JS_TYPE_ARRAY, true
	case OBJECT_P:
		p.advance()
		return nodes.JS_TYPE_OBJECT, true
	case SCALAR:
		p.advance()
		return nodes.JS_TYPE_SCALAR, true
	default:
		return nodes.JS_TYPE_ANY, true
	}
}

// parseJsonKeyUniquenessConstraintOpt parses WITH/WITHOUT UNIQUE [KEYS].
//
//	json_key_uniqueness_constraint_opt:
//	    WITH UNIQUE [KEYS] | WITHOUT UNIQUE [KEYS] | /* EMPTY */
func (p *Parser) parseJsonKeyUniquenessConstraintOpt() bool {
	if p.cur.Type == WITH {
		next := p.peekNext()
		if next.Type == UNIQUE {
			p.advance() // consume WITH
			p.advance() // consume UNIQUE
			if p.cur.Type == KEYS {
				p.advance()
			}
			return true
		}
	}
	if p.cur.Type == WITHOUT {
		next := p.peekNext()
		if next.Type == UNIQUE {
			p.advance() // consume WITHOUT
			p.advance() // consume UNIQUE
			if p.cur.Type == KEYS {
				p.advance()
			}
			return false
		}
	}
	return false
}

// parseJsonObjectConstructorNullClauseOpt parses NULL ON NULL / ABSENT ON NULL.
//
//	json_object_constructor_null_clause_opt:
//	    NULL_P ON NULL_P  => false (not absent)
//	    | ABSENT ON NULL_P => true (absent)
//	    | /* EMPTY */ => false (default: null on null)
func (p *Parser) parseJsonObjectConstructorNullClauseOpt() bool {
	if p.cur.Type == NULL_P {
		next := p.peekNext()
		if next.Type == ON {
			p.advance() // consume NULL
			p.advance() // consume ON
			p.expect(NULL_P)
			return false
		}
	}
	if p.cur.Type == ABSENT {
		p.advance() // consume ABSENT
		p.expect(ON)
		p.expect(NULL_P)
		return true
	}
	return false
}

// parseJsonArrayConstructorNullClauseOpt parses NULL ON NULL / ABSENT ON NULL for arrays.
// Default is ABSENT ON NULL (true) for arrays, unlike objects.
func (p *Parser) parseJsonArrayConstructorNullClauseOpt() bool {
	if p.cur.Type == NULL_P {
		next := p.peekNext()
		if next.Type == ON {
			p.advance() // consume NULL
			p.advance() // consume ON
			p.expect(NULL_P)
			return false
		}
	}
	if p.cur.Type == ABSENT {
		p.advance() // consume ABSENT
		p.expect(ON)
		p.expect(NULL_P)
		return true
	}
	return true // default for arrays is ABSENT ON NULL
}

// parseJsonNameAndValueList parses a comma-separated list of key-value pairs.
//
//	json_name_and_value_list:
//	    json_name_and_value (',' json_name_and_value)*
func (p *Parser) parseJsonNameAndValueList() (*nodes.List, error) {
	first, err := p.parseJsonNameAndValue()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		kv, err := p.parseJsonNameAndValue()
		if err != nil {
			return nil, err
		}
		items = append(items, kv)
	}
	return &nodes.List{Items: items}, nil
}

// parseJsonNameAndValue parses a single key-value pair.
//
//	json_name_and_value:
//	    c_expr VALUE_P json_value_expr
//	    | a_expr ':' json_value_expr
//
// Since we cannot easily distinguish c_expr VALUE_P from a_expr in an LL parser,
// we parse as a_expr and check for VALUE_P or ':'.
func (p *Parser) parseJsonNameAndValue() (*nodes.JsonKeyValue, error) {
	key, err := p.parseAExpr(0)

	if p.cur.Type == VALUE_P {
		p.advance() // consume VALUE
		val, err := p.parseJsonValueExpr()
		if err != nil {
			return nil, err
		}
		return &nodes.JsonKeyValue{Key: key, Value: val}, nil
	}

	if p.cur.Type == ':' {
		p.advance() // consume ':'
		val, err := p.parseJsonValueExpr()
		if err != nil {
			return nil, err
		}
		return &nodes.JsonKeyValue{Key: key, Value: val}, nil
	}

	// Fallback: treat as key VALUE expr
	val, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	return &nodes.JsonKeyValue{Key: key, Value: val}, nil
}

// parseJsonValueExprList parses a comma-separated list of json_value_expr.
func (p *Parser) parseJsonValueExprList() (*nodes.List, error) {
	first, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		v, err := p.parseJsonValueExpr()
		if err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	return &nodes.List{Items: items}, nil
}

// parseJsonPassingClauseOpt parses optional PASSING clause.
//
//	json_passing_clause_opt:
//	    PASSING json_arguments
//	    | /* EMPTY */
func (p *Parser) parseJsonPassingClauseOpt() (*nodes.List, error) {
	if p.cur.Type != PASSING {
		return nil, nil
	}
	p.advance() // consume PASSING
	return p.parseJsonArguments()
}

// parseJsonArguments parses comma-separated json_argument list.
func (p *Parser) parseJsonArguments() (*nodes.List, error) {
	first, err := p.parseJsonArgument()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		arg, err := p.parseJsonArgument()
		if err != nil {
			return nil, err
		}
		items = append(items, arg)
	}
	return &nodes.List{Items: items}, nil
}

// parseJsonArgument parses json_value_expr AS ColLabel.
func (p *Parser) parseJsonArgument() (*nodes.JsonArgument, error) {
	val, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	name, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}
	return &nodes.JsonArgument{
		Val:  val,
		Name: name,
	}, nil
}

// ---------------------------------------------------------------------------
// JSON constructor functions (dispatched from parseCExpr)
// ---------------------------------------------------------------------------

// parseJsonObjectExpr parses JSON_OBJECT(...).
//
// Ref: https://www.postgresql.org/docs/17/functions-json.html
//
//	JSON_OBJECT '(' func_arg_list ')'                          -- legacy
//	JSON_OBJECT '(' json_name_and_value_list ... ')'           -- SQL/JSON
//	JSON_OBJECT '(' json_returning_clause_opt ')'              -- empty
func (p *Parser) parseJsonObjectExpr() (nodes.Node, error) {
	p.advance() // consume JSON_OBJECT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Empty: JSON_OBJECT() or JSON_OBJECT(RETURNING ...)
	if p.cur.Type == ')' {
		p.advance()
		return &nodes.JsonObjectConstructor{Loc: nodes.NoLoc()}, nil
	}
	if p.cur.Type == RETURNING {
		output, err := p.parseJsonReturningClauseOpt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.JsonObjectConstructor{
			Output: output,
			Loc:    nodes.NoLoc(),
		}, nil
	}

	// We need to disambiguate:
	// 1. Legacy: JSON_OBJECT(func_arg_list) - args separated by commas, no VALUE or ':'
	// 2. SQL/JSON: JSON_OBJECT(key VALUE val, ...) or JSON_OBJECT(key : val, ...)
	//
	// Strategy: parse the first a_expr, then check what follows.
	// If VALUE_P or ':', it's SQL/JSON.
	// Otherwise, it's legacy (func_arg_list).

	// Save position for potential backtrack (we can't easily backtrack, so we parse first expr)
	firstExpr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}

	if p.cur.Type == VALUE_P {
		// SQL/JSON form: key VALUE val
		p.advance()
		firstVal, err := p.parseJsonValueExpr()
		if err != nil {
			return nil, err
		}
		firstKV := &nodes.JsonKeyValue{Key: firstExpr, Value: firstVal}
		items := []nodes.Node{firstKV}
		for p.cur.Type == ',' {
			p.advance()
			kv, err := p.parseJsonNameAndValue()
			if err != nil {
				return nil, err
			}
			items = append(items, kv)
		}
		exprs := &nodes.List{Items: items}
		absentOnNull := p.parseJsonObjectConstructorNullClauseOpt()
		uniqueKeys := p.parseJsonKeyUniquenessConstraintOpt()
		output, err := p.parseJsonReturningClauseOpt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.JsonObjectConstructor{
			Exprs:        exprs,
			Output:       output,
			AbsentOnNull: absentOnNull,
			UniqueKeys:   uniqueKeys,
			Loc:          nodes.NoLoc(),
		}, nil
	}

	if p.cur.Type == ':' {
		// SQL/JSON form: key : val
		p.advance()
		firstVal, err := p.parseJsonValueExpr()
		if err != nil {
			return nil, err
		}
		firstKV := &nodes.JsonKeyValue{Key: firstExpr, Value: firstVal}
		items := []nodes.Node{firstKV}
		for p.cur.Type == ',' {
			p.advance()
			kv, err := p.parseJsonNameAndValue()
			if err != nil {
				return nil, err
			}
			items = append(items, kv)
		}
		exprs := &nodes.List{Items: items}
		absentOnNull := p.parseJsonObjectConstructorNullClauseOpt()
		uniqueKeys := p.parseJsonKeyUniquenessConstraintOpt()
		output, err := p.parseJsonReturningClauseOpt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.JsonObjectConstructor{
			Exprs:        exprs,
			Output:       output,
			AbsentOnNull: absentOnNull,
			UniqueKeys:   uniqueKeys,
			Loc:          nodes.NoLoc(),
		}, nil
	}

	// Legacy form: JSON_OBJECT(func_arg_list)
	// We already parsed the first arg. Continue collecting comma-separated args.
	args := []nodes.Node{firstExpr}
	for p.cur.Type == ',' {
		p.advance()
		arg, err := p.parseFuncArgExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "json_object"),
		Args:       &nodes.List{Items: args},
		FuncFormat: int(nodes.COERCE_EXPLICIT_CALL),
		Loc:        nodes.NoLoc(),
	}, nil
}

// parseJsonArrayExpr parses JSON_ARRAY(...).
//
// Ref: https://www.postgresql.org/docs/17/functions-json.html
//
//	JSON_ARRAY '(' json_value_expr_list ... ')'     -- with values
//	JSON_ARRAY '(' select_no_parens ... ')'         -- with subquery
//	JSON_ARRAY '(' json_returning_clause_opt ')'    -- empty
func (p *Parser) parseJsonArrayExpr() (nodes.Node, error) {
	p.advance() // consume JSON_ARRAY
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Empty: JSON_ARRAY() or JSON_ARRAY(RETURNING ...)
	if p.cur.Type == ')' {
		p.advance()
		return &nodes.JsonArrayConstructor{Loc: nodes.NoLoc()}, nil
	}
	if p.cur.Type == RETURNING {
		output, err := p.parseJsonReturningClauseOpt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.JsonArrayConstructor{
			Output: output,
			Loc:    nodes.NoLoc(),
		}, nil
	}

	// Check for subquery: SELECT, VALUES, WITH, TABLE
	if p.isSelectStart() {
		query, err := p.parseSelectStmtForExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.parseJsonFormatClauseOpt(); err != nil {
			return nil, err
		}
		output, err := p.parseJsonReturningClauseOpt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.JsonArrayQueryConstructor{
			Query:  query,
			Output: output,
			Loc:    nodes.NoLoc(),
		}, nil
	}

	// Value list
	exprs, err := p.parseJsonValueExprList()
	if err != nil {
		return nil, err
	}
	absentOnNull := p.parseJsonArrayConstructorNullClauseOpt()
	output, err := p.parseJsonReturningClauseOpt()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.JsonArrayConstructor{
		Exprs:        exprs,
		AbsentOnNull: absentOnNull,
		Output:       output,
		Loc:          nodes.NoLoc(),
	}, nil
}

// parseJsonParseExpr parses JSON(json_value_expr [WITH UNIQUE [KEYS]]).
//
// Ref: https://www.postgresql.org/docs/17/functions-json.html
//
//	JSON '(' json_value_expr json_key_uniqueness_constraint_opt ')'
func (p *Parser) parseJsonParseExpr() (nodes.Node, error) {
	p.advance() // consume JSON
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	val, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	uniqueKeys := p.parseJsonKeyUniquenessConstraintOpt()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.JsonParseExpr{
		Expr:       val,
		UniqueKeys: uniqueKeys,
		Loc:        nodes.NoLoc(),
	}, nil
}

// parseJsonScalarExpr parses JSON_SCALAR(a_expr).
//
//	JSON_SCALAR '(' a_expr ')'
func (p *Parser) parseJsonScalarExpr() (nodes.Node, error) {
	p.advance() // consume JSON_SCALAR
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.JsonScalarExpr{
		Expr: expr,
		Loc:  nodes.NoLoc(),
	}, nil
}

// parseJsonSerializeExpr parses JSON_SERIALIZE(json_value_expr [RETURNING ...]).
//
//	JSON_SERIALIZE '(' json_value_expr json_returning_clause_opt ')'
func (p *Parser) parseJsonSerializeExpr() (nodes.Node, error) {
	p.advance() // consume JSON_SERIALIZE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	val, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	output, err := p.parseJsonReturningClauseOpt()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.JsonSerializeExpr{
		Expr:   val,
		Output: output,
		Loc:    nodes.NoLoc(),
	}, nil
}

// parseJsonQueryExpr parses JSON_QUERY(...).
//
// Ref: https://www.postgresql.org/docs/17/functions-json.html
//
//	JSON_QUERY '(' json_value_expr ',' a_expr
//	    json_passing_clause_opt json_returning_clause_opt
//	    json_wrapper_behavior json_quotes_clause_opt
//	    json_behavior_clause_opt ')'
func (p *Parser) parseJsonQueryExpr() (nodes.Node, error) {
	p.advance() // consume JSON_QUERY
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	contextItem, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	pathspec, err := p.parseAExpr(0)
	passing, err := p.parseJsonPassingClauseOpt()
	if err != nil {
		return nil, err
	}
	output, err := p.parseJsonReturningClauseOpt()
	if err != nil {
		return nil, err
	}
	wrapper := p.parseJsonWrapperBehavior()
	quotes := p.parseJsonQuotesClauseOpt()
	onEmpty, onError := p.parseJsonBehaviorClauseOpt()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.JsonFuncExpr{
		Op:          nodes.JSON_QUERY_OP,
		ContextItem: contextItem,
		Pathspec:    pathspec,
		Passing:     passing,
		Output:      output,
		Wrapper:     wrapper,
		Quotes:      quotes,
		OnEmpty:     onEmpty,
		OnError:     onError,
		Loc:         nodes.NoLoc(),
	}, nil
}

// parseJsonExistsExpr parses JSON_EXISTS(...).
//
//	JSON_EXISTS '(' json_value_expr ',' a_expr
//	    json_passing_clause_opt json_on_error_clause_opt ')'
func (p *Parser) parseJsonExistsExpr() (nodes.Node, error) {
	p.advance() // consume JSON_EXISTS
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	contextItem, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	pathspec, err := p.parseAExpr(0)
	passing, err := p.parseJsonPassingClauseOpt()
	if err != nil {
		return nil, err
	}
	onError := p.parseJsonOnErrorClauseOpt()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.JsonFuncExpr{
		Op:          nodes.JSON_EXISTS_OP,
		ContextItem: contextItem,
		Pathspec:    pathspec,
		Passing:     passing,
		OnError:     onError,
		Loc:         nodes.NoLoc(),
	}, nil
}

// parseJsonValueFuncExpr parses JSON_VALUE(...).
//
//	JSON_VALUE '(' json_value_expr ',' a_expr
//	    json_passing_clause_opt json_returning_clause_opt
//	    json_behavior_clause_opt ')'
func (p *Parser) parseJsonValueFuncExpr() (nodes.Node, error) {
	p.advance() // consume JSON_VALUE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	contextItem, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	pathspec, err := p.parseAExpr(0)
	passing, err := p.parseJsonPassingClauseOpt()
	if err != nil {
		return nil, err
	}
	output, err := p.parseJsonReturningClauseOpt()
	if err != nil {
		return nil, err
	}
	onEmpty, onError := p.parseJsonBehaviorClauseOpt()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.JsonFuncExpr{
		Op:          nodes.JSON_VALUE_OP,
		ContextItem: contextItem,
		Pathspec:    pathspec,
		Passing:     passing,
		Output:      output,
		OnEmpty:     onEmpty,
		OnError:     onError,
		Loc:         nodes.NoLoc(),
	}, nil
}

// ---------------------------------------------------------------------------
// JSON aggregate functions
// ---------------------------------------------------------------------------

// parseJsonObjectAgg parses JSON_OBJECTAGG(...).
//
//	JSON_OBJECTAGG '(' json_name_and_value
//	    json_object_constructor_null_clause_opt
//	    json_key_uniqueness_constraint_opt
//	    json_returning_clause_opt ')'
//
// After parsing, FILTER and OVER clauses are handled by the caller.
func (p *Parser) parseJsonObjectAgg() (nodes.Node, error) {
	p.advance() // consume JSON_OBJECTAGG
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	kv, err := p.parseJsonNameAndValue()
	if err != nil {
		return nil, err
	}
	absentOnNull := p.parseJsonObjectConstructorNullClauseOpt()
	uniqueKeys := p.parseJsonKeyUniquenessConstraintOpt()
	output, err := p.parseJsonReturningClauseOpt()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	agg := &nodes.JsonObjectAgg{
		Constructor: &nodes.JsonAggConstructor{
			Output: output,
			Loc:    nodes.NoLoc(),
		},
		Arg:          kv,
		AbsentOnNull: absentOnNull,
		UniqueKeys:   uniqueKeys,
	}

	// Apply filter_clause and over_clause
	if err := p.applyJsonAggClauses(agg.Constructor); err != nil {
		return nil, err
	}

	return agg, nil
}

// parseJsonArrayAgg parses JSON_ARRAYAGG(...).
//
//	JSON_ARRAYAGG '(' json_value_expr
//	    json_array_aggregate_order_by_clause_opt
//	    json_array_constructor_null_clause_opt
//	    json_returning_clause_opt ')'
func (p *Parser) parseJsonArrayAgg() (nodes.Node, error) {
	p.advance() // consume JSON_ARRAYAGG
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	val, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	aggOrder, err := p.parseOptSortClause()
	absentOnNull := p.parseJsonArrayConstructorNullClauseOpt()
	output, err := p.parseJsonReturningClauseOpt()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	agg := &nodes.JsonArrayAgg{
		Constructor: &nodes.JsonAggConstructor{
			Output:    output,
			Agg_order: aggOrder,
			Loc:       nodes.NoLoc(),
		},
		Arg:          val,
		AbsentOnNull: absentOnNull,
	}

	// Apply filter_clause and over_clause
	if err := p.applyJsonAggClauses(agg.Constructor); err != nil {
		return nil, err
	}

	return agg, nil
}

// applyJsonAggClauses parses and applies FILTER and OVER clauses to a JsonAggConstructor.
func (p *Parser) applyJsonAggClauses(c *nodes.JsonAggConstructor) error {
	// filter_clause: FILTER '(' WHERE a_expr ')'
	if p.cur.Type == FILTER {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return err
		}
		if _, err := p.expect(WHERE); err != nil {
			return err
		}
		c.Agg_filter, _ = p.parseAExpr(0)
		if _, err := p.expect(')'); err != nil {
			return err
		}
	}

	// over_clause: OVER window_specification | OVER ColId
	if p.cur.Type == OVER {
		over, err := p.parseOverClause()
		if err != nil {
			return err
		}
		if wd, ok := over.(*nodes.WindowDef); ok {
			c.Over = wd
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// JSON_TABLE (used in FROM clause)
// ---------------------------------------------------------------------------

// parseJsonTable parses a JSON_TABLE(...) expression.
//
// Ref: https://www.postgresql.org/docs/17/functions-json.html#FUNCTIONS-SQLJSON-TABLE
//
//	JSON_TABLE '('
//	    json_value_expr ',' a_expr json_table_path_name_opt
//	    json_passing_clause_opt
//	    COLUMNS '(' json_table_column_definition_list ')'
//	    json_on_error_clause_opt
//	')'
func (p *Parser) parseJsonTable() (*nodes.JsonTable, error) {
	p.advance() // consume JSON_TABLE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	contextItem, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	pathExpr, err := p.parseAExpr(0)

	// json_table_path_name_opt: AS name | empty
	pathName := ""
	if p.cur.Type == AS {
		p.advance()
		pathName, _ = p.parseColId()
	}

	passing, err := p.parseJsonPassingClauseOpt()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(COLUMNS); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	columns, err := p.parseJsonTableColumnDefinitionList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	onError := p.parseJsonOnErrorClauseOpt()

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.JsonTable{
		ContextItem: contextItem,
		Pathspec: &nodes.JsonTablePathSpec{
			String: pathExpr,
			Name:   pathName,
			Loc:    nodes.NoLoc(),
		},
		Passing: passing,
		Columns: columns,
		OnError: onError,
		Loc:     nodes.NoLoc(),
	}, nil
}

// parseJsonTableColumnDefinitionList parses comma-separated column definitions.
func (p *Parser) parseJsonTableColumnDefinitionList() (*nodes.List, error) {
	first, err := p.parseJsonTableColumnDefinition()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		col, err := p.parseJsonTableColumnDefinition()
		if err != nil {
			return nil, err
		}
		items = append(items, col)
	}
	return &nodes.List{Items: items}, nil
}

// parseJsonTableColumnDefinition parses a single JSON_TABLE column definition.
//
//	ColId FOR ORDINALITY
//	| ColId Typename [FORMAT JSON] [EXISTS] [PATH Sconst] [wrapper] [quotes] [behavior]
//	| NESTED [PATH] Sconst [AS name] COLUMNS '(' column_list ')'
func (p *Parser) parseJsonTableColumnDefinition() (*nodes.JsonTableColumn, error) {
	// NESTED path
	if p.cur.Type == NESTED {
		p.advance() // consume NESTED
		// path_opt: PATH | empty
		if p.cur.Type == PATH {
			p.advance()
		}
		// Sconst
		pathStr := p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}

		// json_table_path_name_opt
		pathName := ""
		if p.cur.Type == AS {
			p.advance()
			pathName, _ = p.parseColId()
		}

		if _, err := p.expect(COLUMNS); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		columns, err := p.parseJsonTableColumnDefinitionList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		return &nodes.JsonTableColumn{
			Coltype: nodes.JTC_NESTED,
			Pathspec: &nodes.JsonTablePathSpec{
				String: &nodes.A_Const{Val: &nodes.String{Str: pathStr}},
				Name:   pathName,
				Loc:    nodes.NoLoc(),
			},
			Columns: columns,
			Loc:     nodes.NoLoc(),
		}, nil
	}

	// ColId ...
	colName, _ := p.parseColId()

	// ColId FOR ORDINALITY
	if p.cur.Type == FOR {
		p.advance()
		if _, err := p.expect(ORDINALITY); err != nil {
			return nil, err
		}
		return &nodes.JsonTableColumn{
			Coltype: nodes.JTC_FOR_ORDINALITY,
			Name:    colName,
			Loc:     nodes.NoLoc(),
		}, nil
	}

	// ColId Typename ...
	tn, err := p.parseTypename()
	if err != nil {
		return &nodes.JsonTableColumn{
			Name: colName,
			Loc:  nodes.NoLoc(),
		}, nil
	}

	// Check for FORMAT JSON clause
	var format *nodes.JsonFormat
	if p.cur.Type == FORMAT {
		format, err = p.parseJsonFormatClause()
		if err != nil {
			return nil, err
		}
	}

	// Check for EXISTS
	isExists := false
	if p.cur.Type == EXISTS {
		p.advance()
		isExists = true
	}

	if isExists {
		// EXISTS column: ColId Typename EXISTS [PATH Sconst] [ON ERROR]
		pathspec, err := p.parseJsonTableColumnPathClauseOpt()
		if err != nil {
			return nil, err
		}
		onError := p.parseJsonOnErrorClauseOpt()
		return &nodes.JsonTableColumn{
			Coltype:  nodes.JTC_EXISTS,
			Name:     colName,
			TypeName: tn,
			Pathspec: pathspec,
			OnError:  onError,
			Loc:      nodes.NoLoc(),
		}, nil
	}

	if format != nil {
		// FORMATTED column: ColId Typename FORMAT JSON [PATH Sconst] [wrapper] [quotes] [behavior]
		pathspec, err := p.parseJsonTableColumnPathClauseOpt()
		if err != nil {
			return nil, err
		}
		wrapper := p.parseJsonWrapperBehavior()
		quotes := p.parseJsonQuotesClauseOpt()
		onEmpty, onError := p.parseJsonBehaviorClauseOpt()
		return &nodes.JsonTableColumn{
			Coltype:  nodes.JTC_FORMATTED,
			Name:     colName,
			TypeName: tn,
			Format:   format,
			Pathspec: pathspec,
			Wrapper:  wrapper,
			Quotes:   quotes,
			OnEmpty:  onEmpty,
			OnError:  onError,
			Loc:      nodes.NoLoc(),
		}, nil
	}

	// REGULAR column: ColId Typename [PATH Sconst] [wrapper] [quotes] [behavior]
	pathspec, err := p.parseJsonTableColumnPathClauseOpt()
	if err != nil {
		return nil, err
	}
	wrapper := p.parseJsonWrapperBehavior()
	quotes := p.parseJsonQuotesClauseOpt()
	onEmpty, onError := p.parseJsonBehaviorClauseOpt()
	return &nodes.JsonTableColumn{
		Coltype:  nodes.JTC_REGULAR,
		Name:     colName,
		TypeName: tn,
		Pathspec: pathspec,
		Wrapper:  wrapper,
		Quotes:   quotes,
		OnEmpty:  onEmpty,
		OnError:  onError,
		Loc:      nodes.NoLoc(),
	}, nil
}

// parseJsonTableColumnPathClauseOpt parses optional PATH Sconst.
//
//	json_table_column_path_clause_opt:
//	    PATH Sconst | /* EMPTY */
func (p *Parser) parseJsonTableColumnPathClauseOpt() (*nodes.JsonTablePathSpec, error) {
	if p.cur.Type != PATH {
		return nil, nil
	}
	p.advance() // consume PATH
	pathStr := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	return &nodes.JsonTablePathSpec{
		String: &nodes.A_Const{Val: &nodes.String{Str: pathStr}},
		Loc:    nodes.NoLoc(),
	}, nil
}

// parseJsonIsPredicate parses IS [NOT] JSON [...] for the IS postfix handler.
// Called from parseIsPostfix when cur is JSON.
func (p *Parser) parseJsonIsPredicate(left nodes.Node, negated bool) nodes.Node {
	itemType, _ := p.parseJsonPredicateTypeConstraint()
	uniqueKeys := p.parseJsonKeyUniquenessConstraintOpt()
	pred := &nodes.JsonIsPredicate{
		Expr:       left,
		ItemType:   itemType,
		UniqueKeys: uniqueKeys,
		Loc:        nodes.NoLoc(),
	}
	if negated {
		return &nodes.BoolExpr{
			Boolop: nodes.NOT_EXPR,
			Args:   &nodes.List{Items: []nodes.Node{pred}},
			Loc:    nodes.NoLoc(),
		}
	}
	return pred
}
