package parser

import (
	"fmt"
	"strings"

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
	loc := p.pos()
	expr, err := p.parseRequiredAExpr(0)
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
		Loc:     nodes.Loc{Start: loc, End: p.prev.End},
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
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		// gram.y's action takes only these names, and reports no position.
		switch strings.ToLower(name) {
		case "utf8":
			f.Encoding = nodes.JS_ENC_UTF8
		case "utf16":
			f.Encoding = nodes.JS_ENC_UTF16
		case "utf32":
			f.Encoding = nodes.JS_ENC_UTF32
		default:
			return nil, &ParseError{Code: "22023", Message: fmt.Sprintf("unrecognized JSON encoding: %s", name), Position: -1}
		}
	}
	return f, nil
}

// parseJsonFormatClauseOpt parses an optional FORMAT clause.
func (p *Parser) parseJsonFormatClauseOpt() (*nodes.JsonFormat, error) {
	if p.isJsonFormatStart() {
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
	loc := p.pos()
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
		Loc:      nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseJsonBehavior parses a json_behavior production.
//
//	json_behavior:
//	    DEFAULT a_expr
//	    | json_behavior_type
func (p *Parser) parseJsonBehavior() (*nodes.JsonBehavior, error) {
	if p.cur.Type == DEFAULT {
		p.advance()
		expr, err := p.parseRequiredAExpr(0)
		if err != nil {
			return nil, err
		}
		return &nodes.JsonBehavior{
			Btype: nodes.JSON_BEHAVIOR_DEFAULT,
			Expr:  expr,
			Loc:   nodes.NoLoc(),
		}, nil
	}
	bt, ok := p.parseJsonBehaviorType()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.JsonBehavior{
		Btype: bt,
		Loc:   nodes.NoLoc(),
	}, nil
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
func (p *Parser) parseJsonBehaviorClauseOpt() (onEmpty, onError *nodes.JsonBehavior, err error) {
	if !p.isJsonBehaviorStart() {
		return nil, nil, nil
	}

	first, err := p.parseJsonBehavior()
	if err != nil {
		return nil, nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, nil, err
	}
	if p.cur.Type == ERROR_P {
		p.advance() // consume ERROR
		return nil, first, nil
	}
	if _, err := p.expect(EMPTY_P); err != nil {
		return nil, nil, err
	}
	onEmpty = first

	// json_behavior ON EMPTY_P json_behavior ON ERROR_P
	if p.isJsonBehaviorStart() {
		onError, err = p.parseJsonOnError()
		if err != nil {
			return nil, nil, err
		}
	}
	return onEmpty, onError, nil
}

// parseJsonOnErrorClauseOpt parses optional json_behavior ON ERROR.
//
//	json_on_error_clause_opt:
//	    json_behavior ON ERROR_P
//	    | /* EMPTY */
func (p *Parser) parseJsonOnErrorClauseOpt() (*nodes.JsonBehavior, error) {
	if !p.isJsonBehaviorStart() {
		return nil, nil
	}
	return p.parseJsonOnError()
}

// parseJsonOnError parses json_behavior ON ERROR_P.
func (p *Parser) parseJsonOnError() (*nodes.JsonBehavior, error) {
	beh, err := p.parseJsonBehavior()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	if _, err := p.expect(ERROR_P); err != nil {
		return nil, err
	}
	return beh, nil
}

// parseJsonWrapperBehavior parses wrapper behavior options.
//
//	json_wrapper_behavior:
//	    WITHOUT [ARRAY] WRAPPER
//	    | WITH [UNCONDITIONAL|CONDITIONAL] [ARRAY] WRAPPER
//	    | /* EMPTY */
func (p *Parser) parseJsonWrapperBehavior() (nodes.JsonWrapper, error) {
	var wrapper nodes.JsonWrapper
	switch p.cur.Type {
	case WITHOUT:
		p.advance()
		wrapper = nodes.JSW_NONE
	case WITH:
		p.advance()
		wrapper = nodes.JSW_UNCONDITIONAL
		if p.cur.Type == CONDITIONAL {
			p.advance()
			wrapper = nodes.JSW_CONDITIONAL
		} else if p.cur.Type == UNCONDITIONAL {
			p.advance()
		}
	default:
		return nodes.JSW_UNSPEC, nil
	}
	if p.cur.Type == ARRAY {
		p.advance()
	}
	if _, err := p.expect(WRAPPER); err != nil {
		return 0, err
	}
	return wrapper, nil
}

// parseJsonQuotesClauseOpt parses KEEP/OMIT QUOTES options.
//
//	json_quotes_clause_opt:
//	    KEEP QUOTES [ON SCALAR STRING_P]
//	    | OMIT QUOTES [ON SCALAR STRING_P]
//	    | /* EMPTY */
func (p *Parser) parseJsonQuotesClauseOpt() (nodes.JsonQuotes, error) {
	var quotes nodes.JsonQuotes
	switch p.cur.Type {
	case KEEP:
		quotes = nodes.JS_QUOTES_KEEP
	case OMIT:
		quotes = nodes.JS_QUOTES_OMIT
	default:
		return nodes.JS_QUOTES_UNSPEC, nil
	}
	p.advance() // consume KEEP or OMIT
	if _, err := p.expect(QUOTES); err != nil {
		return 0, err
	}
	if p.cur.Type == ON {
		p.advance()
		if _, err := p.expect(SCALAR); err != nil {
			return 0, err
		}
		if _, err := p.expect(STRING_P); err != nil {
			return 0, err
		}
	}
	return quotes, nil
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
//
// WITH and WITHOUT always start it: PostgreSQL shifts WITH there, so even
// CREATE TABLE AS SELECT a IS JSON WITH DATA is a syntax error at DATA.
func (p *Parser) parseJsonKeyUniquenessConstraintOpt() (bool, error) {
	if p.cur.Type == WITH {
		p.advance() // consume WITH
		if _, err := p.expect(UNIQUE); err != nil {
			return false, err
		}
		if p.cur.Type == KEYS {
			p.advance()
		}
		return true, nil
	}
	if p.cur.Type == WITHOUT {
		p.advance() // consume WITHOUT
		if _, err := p.expect(UNIQUE); err != nil {
			return false, err
		}
		if p.cur.Type == KEYS {
			p.advance()
		}
	}
	return false, nil
}

// parseJsonObjectConstructorNullClauseOpt parses NULL ON NULL / ABSENT ON NULL.
//
//	json_object_constructor_null_clause_opt:
//	    NULL_P ON NULL_P  => false (not absent)
//	    | ABSENT ON NULL_P => true (absent)
//	    | /* EMPTY */ => false (default: null on null)
func (p *Parser) parseJsonObjectConstructorNullClauseOpt() (bool, error) {
	return p.parseJsonNullClauseOpt(false)
}

// parseJsonNullClauseOpt parses NULL ON NULL (false) or ABSENT ON NULL
// (true), and returns dflt when neither starts here. Nothing else takes
// NULL or ABSENT after the values, so either one commits to the clause.
func (p *Parser) parseJsonNullClauseOpt(dflt bool) (bool, error) {
	var absent bool
	switch p.cur.Type {
	case NULL_P:
		absent = false
	case ABSENT:
		absent = true
	default:
		return dflt, nil
	}
	p.advance() // consume NULL or ABSENT
	if _, err := p.expect(ON); err != nil {
		return false, err
	}
	if _, err := p.expect(NULL_P); err != nil {
		return false, err
	}
	return absent, nil
}

// parseJsonArrayConstructorNullClauseOpt parses NULL ON NULL / ABSENT ON NULL for arrays.
// Default is ABSENT ON NULL (true) for arrays, unlike objects.
func (p *Parser) parseJsonArrayConstructorNullClauseOpt() (bool, error) {
	return p.parseJsonNullClauseOpt(true) // default for arrays is ABSENT ON NULL
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
	loc := p.pos()
	key, err := p.parseRequiredAExpr(0)
	if err != nil {
		return nil, err
	}

	if p.cur.Type == VALUE_P {
		p.advance() // consume VALUE
		val, err := p.parseJsonValueExpr()
		if err != nil {
			return nil, err
		}
		return &nodes.JsonKeyValue{Key: key, Value: val, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}

	if p.cur.Type == ':' {
		p.advance() // consume ':'
		val, err := p.parseJsonValueExpr()
		if err != nil {
			return nil, err
		}
		return &nodes.JsonKeyValue{Key: key, Value: val, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}

	return nil, p.syntaxErrorAtCur()
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
	loc := p.pos()
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
		Loc:  nodes.Loc{Start: loc, End: p.prev.End},
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

	// JSON_OBJECT '(' func_arg_list ')' with a named argument first.
	if p.startsNamedFuncArg() {
		return p.parseSpecialFuncArgCall("pg_catalog", "json_object")
	}

	// We need to disambiguate:
	// 1. Legacy: JSON_OBJECT(func_arg_list) - args separated by commas, no VALUE or ':'
	// 2. SQL/JSON: JSON_OBJECT(key VALUE val, ...) or JSON_OBJECT(key : val, ...)
	//
	// Strategy: parse the first a_expr, then check what follows.
	// If VALUE_P or ':', it's SQL/JSON.
	// Otherwise, it's legacy (func_arg_list).

	// Save position for potential backtrack (we can't easily backtrack, so we parse first expr)
	firstExprLoc := p.pos()
	firstExpr, err := p.parseRequiredAExpr(0)
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
		firstKV := &nodes.JsonKeyValue{Key: firstExpr, Value: firstVal, Loc: nodes.Loc{Start: firstExprLoc, End: p.prev.End}}
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
		absentOnNull, err := p.parseJsonObjectConstructorNullClauseOpt()
		if err != nil {
			return nil, err
		}
		uniqueKeys, err := p.parseJsonKeyUniquenessConstraintOpt()
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
		firstKV := &nodes.JsonKeyValue{Key: firstExpr, Value: firstVal, Loc: nodes.Loc{Start: firstExprLoc, End: p.prev.End}}
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
		absentOnNull, err := p.parseJsonObjectConstructorNullClauseOpt()
		if err != nil {
			return nil, err
		}
		uniqueKeys, err := p.parseJsonKeyUniquenessConstraintOpt()
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
	absentOnNull, err := p.parseJsonArrayConstructorNullClauseOpt()
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
	uniqueKeys, err := p.parseJsonKeyUniquenessConstraintOpt()
	if err != nil {
		return nil, err
	}
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
	expr, err := p.parseRequiredAExpr(0)
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
	pathspec, err := p.parseRequiredAExpr(0)
	if err != nil {
		return nil, err
	}
	passing, err := p.parseJsonPassingClauseOpt()
	if err != nil {
		return nil, err
	}
	output, err := p.parseJsonReturningClauseOpt()
	if err != nil {
		return nil, err
	}
	wrapper, err := p.parseJsonWrapperBehavior()
	if err != nil {
		return nil, err
	}
	quotes, err := p.parseJsonQuotesClauseOpt()
	if err != nil {
		return nil, err
	}
	onEmpty, onError, err := p.parseJsonBehaviorClauseOpt()
	if err != nil {
		return nil, err
	}
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
	pathspec, err := p.parseRequiredAExpr(0)
	if err != nil {
		return nil, err
	}
	passing, err := p.parseJsonPassingClauseOpt()
	if err != nil {
		return nil, err
	}
	onError, err := p.parseJsonOnErrorClauseOpt()
	if err != nil {
		return nil, err
	}
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
	pathspec, err := p.parseRequiredAExpr(0)
	if err != nil {
		return nil, err
	}
	passing, err := p.parseJsonPassingClauseOpt()
	if err != nil {
		return nil, err
	}
	output, err := p.parseJsonReturningClauseOpt()
	if err != nil {
		return nil, err
	}
	onEmpty, onError, err := p.parseJsonBehaviorClauseOpt()
	if err != nil {
		return nil, err
	}
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
	loc := p.pos()
	p.advance() // consume JSON_OBJECTAGG
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	kv, err := p.parseJsonNameAndValue()
	if err != nil {
		return nil, err
	}
	absentOnNull, err := p.parseJsonObjectConstructorNullClauseOpt()
	if err != nil {
		return nil, err
	}
	uniqueKeys, err := p.parseJsonKeyUniquenessConstraintOpt()
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

	agg.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return agg, nil
}

// parseJsonArrayAgg parses JSON_ARRAYAGG(...).
//
//	JSON_ARRAYAGG '(' json_value_expr
//	    json_array_aggregate_order_by_clause_opt
//	    json_array_constructor_null_clause_opt
//	    json_returning_clause_opt ')'
func (p *Parser) parseJsonArrayAgg() (nodes.Node, error) {
	loc := p.pos()
	p.advance() // consume JSON_ARRAYAGG
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	val, err := p.parseJsonValueExpr()
	if err != nil {
		return nil, err
	}
	aggOrder, err := p.parseOptSortClause()
	if err != nil {
		return nil, err
	}
	absentOnNull, err := p.parseJsonArrayConstructorNullClauseOpt()
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

	agg.Loc = nodes.Loc{Start: loc, End: p.prev.End}
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
		filter, err := p.parseRequiredAExpr(0)
		if err != nil {
			return err
		}
		c.Agg_filter = filter
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
	pathLoc := p.pos()
	pathExpr, err := p.parseRequiredAExpr(0)
	if err != nil {
		return nil, err
	}

	// json_table_path_name_opt: AS name | empty
	pathName := ""
	if p.cur.Type == AS {
		p.advance()
		pathName, err = p.parseColId()
		if err != nil {
			return nil, err
		}
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

	onError, err := p.parseJsonOnErrorClauseOpt()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// gram.y's json_table action takes only a string constant as the path,
	// and runs once the whole call is read.
	if c, ok := pathExpr.(*nodes.A_Const); !ok || c.Isnull {
		return nil, errJsonTablePath(pathLoc)
	} else if _, ok := c.Val.(*nodes.String); !ok {
		return nil, errJsonTablePath(pathLoc)
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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	// NESTED path_opt Sconst ...; before any token but PATH or an
	// Sconst, NESTED, an unreserved keyword, is the column's name.
	if p.isJsonNestedColumnStart() {
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
			var err error
			pathName, err = p.parseColId()
			if err != nil {
				return nil, err
			}
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
	colName, err := p.parseColId()
	if err != nil {
		return nil, err
	}

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

	// ColId Typename ... FORMAT_LA starts no Typename.
	if p.isJsonFormatStart() && !p.collectMode() {
		return nil, p.syntaxErrorAtCur()
	}
	tn, err := p.parseTypename()
	if err != nil {
		return nil, err
	}

	// Check for FORMAT JSON clause
	var format *nodes.JsonFormat
	if p.isJsonFormatStart() {
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
		// EXISTS column: ColId Typename EXISTS [PATH Sconst] [behavior]
		pathspec, err := p.parseJsonTableColumnPathClauseOpt()
		if err != nil {
			return nil, err
		}
		onEmpty, onError, err := p.parseJsonBehaviorClauseOpt()
		if err != nil {
			return nil, err
		}
		return &nodes.JsonTableColumn{
			Coltype:  nodes.JTC_EXISTS,
			Name:     colName,
			TypeName: tn,
			Pathspec: pathspec,
			OnEmpty:  onEmpty,
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
		wrapper, err := p.parseJsonWrapperBehavior()
		if err != nil {
			return nil, err
		}
		quotes, err := p.parseJsonQuotesClauseOpt()
		if err != nil {
			return nil, err
		}
		onEmpty, onError, err := p.parseJsonBehaviorClauseOpt()
		if err != nil {
			return nil, err
		}
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
	wrapper, err := p.parseJsonWrapperBehavior()
	if err != nil {
		return nil, err
	}
	quotes, err := p.parseJsonQuotesClauseOpt()
	if err != nil {
		return nil, err
	}
	onEmpty, onError, err := p.parseJsonBehaviorClauseOpt()
	if err != nil {
		return nil, err
	}
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
func (p *Parser) parseJsonIsPredicate(left nodes.Node, negated bool) (nodes.Node, error) {
	itemType, _ := p.parseJsonPredicateTypeConstraint()
	uniqueKeys, err := p.parseJsonKeyUniquenessConstraintOpt()
	if err != nil {
		return nil, err
	}
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
		}, nil
	}
	return pred, nil
}

// errJsonTablePath is the error gram.y's json_table action reports for a
// path specification that is not a string constant.
func errJsonTablePath(pos int) *ParseError {
	return &ParseError{Code: "0A000", Message: "only string constants are supported in JSON_TABLE path specification", Position: pos}
}
