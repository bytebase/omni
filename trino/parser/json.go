package parser

import "github.com/bytebase/omni/trino/ast"

// This file is the `expr-json` DAG node: it implements Trino's SQL/JSON
// function family — JSON_EXISTS / JSON_VALUE / JSON_QUERY / JSON_OBJECT /
// JSON_ARRAY — and the shared json-path-invocation subsystem they embed
// (jsonPathInvocation, jsonValueExpression, jsonRepresentation, jsonArgument,
// and the per-function behavior clauses). It depends on the `expressions` node
// (expr.go / function.go / predicate.go): the dispatch in parsePrimaryAtom
// (expr.go) routes the five reserved JSON_* keywords here via parseJSONFunction,
// and the bodies reuse parseExpr / parseValueExpr / parseType / parseIdentifier
// from that node. JSON_TABLE is deliberately NOT here — it is a relationPrimary
// / primaryExpression form owned by parser-select (a P1 extension), not part of
// this node.
//
// The legacy ANTLR grammar (TrinoParser.g4) the implementation tracks, as the
// five primaryExpression alternatives plus their helper rules:
//
//	jsonExists : JSON_EXISTS_ LPAREN_ jsonPathInvocation
//	               (jsonExistsErrorBehavior ON_ ERROR_)? RPAREN_ ;
//	jsonValue  : JSON_VALUE_ LPAREN_ jsonPathInvocation (RETURNING_ type)?
//	               (emptyBehavior = jsonValueBehavior ON_ EMPTY_)?
//	               (errorBehavior = jsonValueBehavior ON_ ERROR_)? RPAREN_ ;
//	jsonQuery  : JSON_QUERY_ LPAREN_ jsonPathInvocation
//	               (RETURNING_ type (FORMAT_ jsonRepresentation)?)?
//	               (jsonQueryWrapperBehavior WRAPPER_)?
//	               ((KEEP_ | OMIT_) QUOTES_ (ON_ SCALAR_ TEXT_STRING_)?)?
//	               (emptyBehavior = jsonQueryBehavior ON_ EMPTY_)?
//	               (errorBehavior = jsonQueryBehavior ON_ ERROR_)? RPAREN_ ;
//	jsonObject : JSON_OBJECT_ LPAREN_
//	               ( jsonObjectMember (COMMA_ jsonObjectMember)*
//	                 (NULL_ ON_ NULL_ | ABSENT_ ON_ NULL_)?
//	                 (WITH_ UNIQUE_ KEYS_? | WITHOUT_ UNIQUE_ KEYS_?)? )?
//	               (RETURNING_ type (FORMAT_ jsonRepresentation)?)? RPAREN_ ;
//	jsonArray  : JSON_ARRAY_ LPAREN_
//	               ( jsonValueExpression (COMMA_ jsonValueExpression)*
//	                 (NULL_ ON_ NULL_ | ABSENT_ ON_ NULL_)? )?
//	               (RETURNING_ type (FORMAT_ jsonRepresentation)?)? RPAREN_ ;
//
//	jsonPathInvocation : jsonValueExpression COMMA_ path = string_
//	                       (PASSING_ jsonArgument (COMMA_ jsonArgument)*)? ;
//	jsonValueExpression : expression (FORMAT_ jsonRepresentation)? ;
//	jsonRepresentation  : JSON_ (ENCODING_ (UTF8_ | UTF16_ | UTF32_))? ;
//	jsonArgument        : jsonValueExpression AS_ identifier ;
//	jsonObjectMember    : KEY_? expression VALUE_ jsonValueExpression
//	                    | expression COLON_ jsonValueExpression ;
//	jsonExistsErrorBehavior : TRUE_ | FALSE_ | UNKNOWN_ | ERROR_ ;
//	jsonValueBehavior       : ERROR_ | NULL_ | DEFAULT_ expression ;
//	jsonQueryWrapperBehavior: WITHOUT_ ARRAY_? | WITH_ (CONDITIONAL_|UNCONDITIONAL_)? ARRAY_? ;
//	jsonQueryBehavior       : ERROR_ | NULL_ | EMPTY_ (ARRAY_ | OBJECT_) ;
//
// The implementation is adjudicated against the live Trino 481 oracle, not the
// literal legacy grammar (json_test.go's oracle differential is the gate).
// Oracle-confirmed facts baked into the parser:
//
//	J1 (clause ordering is FIXED, not free). Each function's optional clauses
//	   appear in exactly the grammar order. `JSON_QUERY(c,'$.x' KEEP QUOTES WITH
//	   ARRAY WRAPPER)` (QUOTES before WRAPPER) and `JSON_VALUE(c,'$.x' NULL ON
//	   ERROR ERROR ON EMPTY)` (ON ERROR before ON EMPTY) are SYNTAX_ERRORs in
//	   Trino 481. The parsers therefore read the clauses positionally.
//	J2 (path is a string LITERAL, not an expression). The path after the comma
//	   must be a string literal: `JSON_EXISTS(c, somecol)` and
//	   `JSON_VALUE(c,'$.x' || 'y')` are SYNTAX_ERRORs. parseJSONPathInvocation
//	   requires a string token for the path.
//	J3 (the five JSON_* names are RESERVED). Unlike SUBSTRING/POSITION, the
//	   JSON_* function keywords cannot be bare column references (`SELECT
//	   json_exists` is a SYNTAX_ERROR), so the dispatch from parsePrimaryAtom is
//	   unambiguous — no try-special-then-fall-back path is needed.
//	J4 (jsonObjectMember KEY is a backtracking-optional keyword). `KEY` is
//	   non-reserved, so `JSON_OBJECT(KEY VALUE 1)` parses KEY as the key
//	   *expression* (COLUMN key), while `JSON_OBJECT(KEY 'x' VALUE 1)` parses KEY
//	   as the leading keyword and 'x' as the expression. parseJSONObjectMember
//	   speculatively consumes a leading KEY and, if the result is not a valid
//	   `expr VALUE value`, rewinds and parses the member without it.
//	J5 (NULL is both a literal element and the start of NULL ON NULL). In
//	   JSON_ARRAY, `JSON_ARRAY(NULL)` is a single NULL element while
//	   `JSON_ARRAY(1 NULL ON NULL)` is element 1 with NULL-ON-NULL handling, and
//	   `JSON_ARRAY(NULL ON NULL)` is a SYNTAX_ERROR (the leading NULL is consumed
//	   as the element, leaving a dangling ON). The element list is parsed first
//	   (greedily taking a leading NULL/ABSENT-free value), then the null-handling
//	   clause is recognized only AFTER the comma-separated list.

// ---------------------------------------------------------------------------
// SQL/JSON node types (parser-package Expr values; see expr.go file header)
// ---------------------------------------------------------------------------

// JSONPathInvocation is the shared `jsonValueExpression , path [PASSING …]`
// head of every SQL/JSON function (jsonPathInvocation). Input is the JSON input
// value expression (with its optional FORMAT); Path is the decoded path string
// literal; Passing carries the named PASSING arguments (nil when absent).
type JSONPathInvocation struct {
	Input   *JSONValueExpr
	Path    string // decoded path string-literal content
	Passing []JSONArgument
	Loc     ast.Loc
}

// JSONValueExpr is `expression [FORMAT jsonRepresentation]` (jsonValueExpression):
// an ordinary expression optionally tagged with an input/output JSON format.
// Format is nil when no FORMAT clause is present.
type JSONValueExpr struct {
	Expr   Expr
	Format *JSONFormat // nil when absent
	Loc    ast.Loc
}

// JSONFormat is `FORMAT JSON [ENCODING (UTF8|UTF16|UTF32)]` (jsonRepresentation).
// Only the JSON representation keyword is accepted by Trino 481; Encoding is ""
// when no ENCODING clause is present, else "UTF8"/"UTF16"/"UTF32".
type JSONFormat struct {
	Encoding string // "", "UTF8", "UTF16", or "UTF32"
	Loc      ast.Loc
}

// JSONArgument is one `jsonValueExpression AS identifier` PASSING argument
// (jsonArgument).
type JSONArgument struct {
	Value *JSONValueExpr
	Name  *ast.Identifier
	Loc   ast.Loc
}

// JSONExistsExpr is `JSON_EXISTS(jsonPathInvocation [behavior ON ERROR])`
// (jsonExists). OnError is "", "TRUE", "FALSE", "UNKNOWN", or "ERROR".
type JSONExistsExpr struct {
	Path    *JSONPathInvocation
	OnError string // jsonExistsErrorBehavior, "" when absent
	Loc     ast.Loc
}

func (n *JSONExistsExpr) Span() ast.Loc { return n.Loc }
func (*JSONExistsExpr) exprNode()       {}

// JSONValueFunc is `JSON_VALUE(jsonPathInvocation [RETURNING type]
// [behavior ON EMPTY] [behavior ON ERROR])` (jsonValue). Returning is the
// optional RETURNING type (nil when absent). OnEmpty/OnError are the optional
// jsonValueBehavior clauses (nil when absent).
type JSONValueFunc struct {
	Path      *JSONPathInvocation
	Returning *DataType // nil when absent
	OnEmpty   *JSONValueBehavior
	OnError   *JSONValueBehavior
	Loc       ast.Loc
}

func (n *JSONValueFunc) Span() ast.Loc { return n.Loc }
func (*JSONValueFunc) exprNode()       {}

// JSONValueBehavior is one `ERROR | NULL | DEFAULT expression` clause of
// JSON_VALUE's ON EMPTY / ON ERROR (jsonValueBehavior). Kind is "ERROR",
// "NULL", or "DEFAULT"; Default is the expression for the DEFAULT form (nil
// otherwise).
type JSONValueBehavior struct {
	Kind    string // "ERROR", "NULL", or "DEFAULT"
	Default Expr   // non-nil only for the DEFAULT form
	Loc     ast.Loc
}

// JSONQueryFunc is `JSON_QUERY(jsonPathInvocation [RETURNING type [FORMAT …]]
// [wrapper WRAPPER] [(KEEP|OMIT) QUOTES [ON SCALAR STRING]] [behavior ON EMPTY]
// [behavior ON ERROR])` (jsonQuery).
type JSONQueryFunc struct {
	Path            *JSONPathInvocation
	Returning       *DataType   // nil when absent
	ReturningFormat *JSONFormat // FORMAT after RETURNING type, nil when absent
	Wrapper         string      // "", "WITHOUT", "WITH", "WITH CONDITIONAL", "WITH UNCONDITIONAL"
	WrapperArray    bool        // the ARRAY keyword was present on the WRAPPER clause
	Quotes          string      // "", "KEEP", or "OMIT"
	QuotesOnScalar  bool        // the ON SCALAR STRING suffix was present on the QUOTES clause
	OnEmpty         *JSONQueryBehavior
	OnError         *JSONQueryBehavior
	Loc             ast.Loc
}

func (n *JSONQueryFunc) Span() ast.Loc { return n.Loc }
func (*JSONQueryFunc) exprNode()       {}

// JSONQueryBehavior is one `ERROR | NULL | EMPTY (ARRAY|OBJECT)` clause of
// JSON_QUERY's ON EMPTY / ON ERROR (jsonQueryBehavior). Kind is "ERROR",
// "NULL", "EMPTY ARRAY", or "EMPTY OBJECT".
type JSONQueryBehavior struct {
	Kind string // "ERROR", "NULL", "EMPTY ARRAY", or "EMPTY OBJECT"
	Loc  ast.Loc
}

// JSONObjectExpr is `JSON_OBJECT([members] [null-handling] [uniqueness]
// [RETURNING type [FORMAT …]])` (jsonObject). Members is the (possibly empty)
// key/value member list. NullHandling is "", "NULL ON NULL", or "ABSENT ON
// NULL". Uniqueness is "", "WITH UNIQUE", or "WITHOUT UNIQUE" (with the trailing
// optional KEYS captured by UniqueKeys).
type JSONObjectExpr struct {
	Members         []JSONObjectMember
	NullHandling    string // "", "NULL ON NULL", or "ABSENT ON NULL"
	Uniqueness      string // "", "WITH UNIQUE", or "WITHOUT UNIQUE"
	UniqueKeys      bool   // the trailing KEYS keyword was present
	Returning       *DataType
	ReturningFormat *JSONFormat
	Loc             ast.Loc
}

func (n *JSONObjectExpr) Span() ast.Loc { return n.Loc }
func (*JSONObjectExpr) exprNode()       {}

// JSONObjectMember is one key/value pair of JSON_OBJECT (jsonObjectMember).
// Key is the key expression; Value is the value (jsonValueExpression). KeyKeyword
// records whether the explicit leading `KEY` keyword was used; Separator is ":"
// (the colon form) or "VALUE" (the KEY?/VALUE form).
type JSONObjectMember struct {
	KeyKeyword bool   // explicit leading KEY keyword
	Key        Expr   // key expression
	Separator  string // ":" or "VALUE"
	Value      *JSONValueExpr
	Loc        ast.Loc
}

// JSONArrayExpr is `JSON_ARRAY([elements] [null-handling] [RETURNING type
// [FORMAT …]])` (jsonArray). Elements is the (possibly empty) element list.
// NullHandling is "", "NULL ON NULL", or "ABSENT ON NULL".
type JSONArrayExpr struct {
	Elements        []*JSONValueExpr
	NullHandling    string // "", "NULL ON NULL", or "ABSENT ON NULL"
	Returning       *DataType
	ReturningFormat *JSONFormat
	Loc             ast.Loc
}

func (n *JSONArrayExpr) Span() ast.Loc { return n.Loc }
func (*JSONArrayExpr) exprNode()       {}

// ---------------------------------------------------------------------------
// Dispatch
// ---------------------------------------------------------------------------

// parseJSONFunction is the entry point for the five SQL/JSON functions, invoked
// by parsePrimaryAtom (expr.go) on a leading JSON_EXISTS / JSON_VALUE /
// JSON_QUERY / JSON_OBJECT / JSON_ARRAY keyword. Because those names are
// reserved (J3), the leading keyword unambiguously selects the function.
func (p *Parser) parseJSONFunction() (Expr, error) {
	switch p.cur.Kind {
	case kwJSON_EXISTS:
		return p.parseJSONExists()
	case kwJSON_VALUE:
		return p.parseJSONValue()
	case kwJSON_QUERY:
		return p.parseJSONQuery()
	case kwJSON_OBJECT:
		return p.parseJSONObject()
	case kwJSON_ARRAY:
		return p.parseJSONArray()
	default:
		// Unreachable via parsePrimaryAtom's dispatch; defensive.
		return nil, p.exprErrorAt("expected a SQL/JSON function")
	}
}

// ---------------------------------------------------------------------------
// JSON_EXISTS
// ---------------------------------------------------------------------------

// parseJSONExists parses `JSON_EXISTS( jsonPathInvocation (behavior ON ERROR)? )`
// (jsonExists). The leading JSON_EXISTS keyword is the current token.
func (p *Parser) parseJSONExists() (Expr, error) {
	startTok := p.advance() // consume JSON_EXISTS
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	path, err := p.parseJSONPathInvocation()
	if err != nil {
		return nil, err
	}
	je := &JSONExistsExpr{Path: path, Loc: ast.Loc{Start: startTok.Loc.Start}}

	// (jsonExistsErrorBehavior ON ERROR)?
	if behavior, ok := jsonExistsErrorBehaviorText(p.cur.Kind); ok && p.peekNext().Kind == kwON {
		p.advance() // behavior keyword
		p.advance() // ON
		if _, err := p.expect(kwERROR); err != nil {
			return nil, err
		}
		je.OnError = behavior
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	je.Loc.End = closeTok.Loc.End
	return je, nil
}

// jsonExistsErrorBehaviorText maps a jsonExistsErrorBehavior keyword to its
// canonical spelling (TRUE | FALSE | UNKNOWN | ERROR).
func jsonExistsErrorBehaviorText(kind TokenKind) (string, bool) {
	switch kind {
	case kwTRUE:
		return "TRUE", true
	case kwFALSE:
		return "FALSE", true
	case kwUNKNOWN:
		return "UNKNOWN", true
	case kwERROR:
		return "ERROR", true
	default:
		return "", false
	}
}

// ---------------------------------------------------------------------------
// JSON_VALUE
// ---------------------------------------------------------------------------

// parseJSONValue parses `JSON_VALUE( jsonPathInvocation (RETURNING type)?
// (behavior ON EMPTY)? (behavior ON ERROR)? )` (jsonValue). The clause order is
// fixed (J1). The leading JSON_VALUE keyword is the current token.
func (p *Parser) parseJSONValue() (Expr, error) {
	startTok := p.advance() // consume JSON_VALUE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	path, err := p.parseJSONPathInvocation()
	if err != nil {
		return nil, err
	}
	jv := &JSONValueFunc{Path: path, Loc: ast.Loc{Start: startTok.Loc.Start}}

	// (RETURNING type)?
	if p.cur.Kind == kwRETURNING {
		p.advance() // consume RETURNING
		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}
		jv.Returning = typ
	}

	// (emptyBehavior ON EMPTY)? (errorBehavior ON ERROR)?  — each clause is
	// self-identifying via its trailing ON EMPTY / ON ERROR. The grammar fixes
	// EMPTY before ERROR (J1): an ON ERROR clause cannot precede an ON EMPTY one.
	for p.startsJSONValueBehavior() {
		behavior, slot, err := p.parseJSONValueBehaviorClause()
		if err != nil {
			return nil, err
		}
		switch slot {
		case kwEMPTY:
			if jv.OnEmpty != nil || jv.OnError != nil {
				return nil, p.exprErrorAt("duplicate or misordered ON EMPTY in JSON_VALUE")
			}
			jv.OnEmpty = behavior
		case kwERROR:
			if jv.OnError != nil {
				return nil, p.exprErrorAt("duplicate ON ERROR in JSON_VALUE")
			}
			jv.OnError = behavior
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	jv.Loc.End = closeTok.Loc.End
	return jv, nil
}

// startsJSONValueBehavior reports whether the current token begins a
// jsonValueBehavior clause (ERROR | NULL | DEFAULT).
func (p *Parser) startsJSONValueBehavior() bool {
	switch p.cur.Kind {
	case kwERROR, kwNULL, kwDEFAULT:
		return true
	default:
		return false
	}
}

// parseJSONValueBehaviorClause parses one `(ERROR | NULL | DEFAULT expression)
// ON (EMPTY | ERROR)` clause (jsonValueBehavior + its ON tail), returning the
// behavior and which slot (kwEMPTY or kwERROR) it fills. The slot is read from
// the trailing keyword so the two clauses are disambiguated regardless of
// order; ordering itself is enforced by the caller (J1).
func (p *Parser) parseJSONValueBehaviorClause() (*JSONValueBehavior, TokenKind, error) {
	startTok := p.cur
	var b *JSONValueBehavior
	switch p.cur.Kind {
	case kwERROR:
		p.advance()
		b = &JSONValueBehavior{Kind: "ERROR", Loc: startTok.Loc}
	case kwNULL:
		p.advance()
		b = &JSONValueBehavior{Kind: "NULL", Loc: startTok.Loc}
	case kwDEFAULT:
		p.advance()
		def, err := p.parseExpr()
		if err != nil {
			return nil, 0, err
		}
		b = &JSONValueBehavior{Kind: "DEFAULT", Default: def, Loc: ast.Loc{Start: startTok.Loc.Start, End: def.Span().End}}
	default:
		return nil, 0, p.exprErrorAt("expected ERROR, NULL, or DEFAULT")
	}
	slot, err := p.expectOnSlot()
	if err != nil {
		return nil, 0, err
	}
	return b, slot, nil
}

// expectOnSlot consumes the required `ON (EMPTY | ERROR)` tail of a behavior
// clause and returns the slot keyword (kwEMPTY or kwERROR). It errors if ON is
// missing or the keyword is neither EMPTY nor ERROR.
func (p *Parser) expectOnSlot() (TokenKind, error) {
	if _, err := p.expect(kwON); err != nil {
		return 0, err
	}
	switch p.cur.Kind {
	case kwEMPTY:
		p.advance()
		return kwEMPTY, nil
	case kwERROR:
		p.advance()
		return kwERROR, nil
	default:
		return 0, p.exprErrorAt("expected EMPTY or ERROR after ON")
	}
}

// ---------------------------------------------------------------------------
// JSON_QUERY
// ---------------------------------------------------------------------------

// parseJSONQuery parses `JSON_QUERY( jsonPathInvocation
// (RETURNING type (FORMAT jsonRepresentation)?)? (wrapper WRAPPER)?
// ((KEEP|OMIT) QUOTES (ON SCALAR STRING)?)? (behavior ON EMPTY)?
// (behavior ON ERROR)? )` (jsonQuery). The clause order is fixed (J1). The
// leading JSON_QUERY keyword is the current token.
func (p *Parser) parseJSONQuery() (Expr, error) {
	startTok := p.advance() // consume JSON_QUERY
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	path, err := p.parseJSONPathInvocation()
	if err != nil {
		return nil, err
	}
	jq := &JSONQueryFunc{Path: path, Loc: ast.Loc{Start: startTok.Loc.Start}}

	// (RETURNING type (FORMAT jsonRepresentation)?)?
	if p.cur.Kind == kwRETURNING {
		p.advance() // consume RETURNING
		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}
		jq.Returning = typ
		if p.cur.Kind == kwFORMAT {
			format, err := p.parseJSONFormat()
			if err != nil {
				return nil, err
			}
			jq.ReturningFormat = format
		}
	}

	// (jsonQueryWrapperBehavior WRAPPER)?
	if wrapper, array, ok := p.tryJSONQueryWrapper(); ok {
		jq.Wrapper = wrapper
		jq.WrapperArray = array
		if _, err := p.expect(kwWRAPPER); err != nil {
			return nil, err
		}
	}

	// ((KEEP | OMIT) QUOTES (ON SCALAR STRING)?)?
	if p.cur.Kind == kwKEEP || p.cur.Kind == kwOMIT {
		if p.cur.Kind == kwKEEP {
			jq.Quotes = "KEEP"
		} else {
			jq.Quotes = "OMIT"
		}
		p.advance() // KEEP or OMIT
		if _, err := p.expect(kwQUOTES); err != nil {
			return nil, err
		}
		// (ON SCALAR STRING)?
		if p.cur.Kind == kwON && p.peekNext().Kind == kwSCALAR {
			p.advance() // ON
			p.advance() // SCALAR
			if _, err := p.expect(kwSTRING); err != nil {
				return nil, err
			}
			jq.QuotesOnScalar = true
		}
	}

	// (emptyBehavior ON EMPTY)? (errorBehavior ON ERROR)?  — each clause is
	// self-identifying via its trailing ON EMPTY / ON ERROR, with EMPTY fixed
	// before ERROR (J1).
	for p.startsJSONQueryBehavior() {
		behavior, slot, err := p.parseJSONQueryBehaviorClause()
		if err != nil {
			return nil, err
		}
		switch slot {
		case kwEMPTY:
			if jq.OnEmpty != nil || jq.OnError != nil {
				return nil, p.exprErrorAt("duplicate or misordered ON EMPTY in JSON_QUERY")
			}
			jq.OnEmpty = behavior
		case kwERROR:
			if jq.OnError != nil {
				return nil, p.exprErrorAt("duplicate ON ERROR in JSON_QUERY")
			}
			jq.OnError = behavior
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	jq.Loc.End = closeTok.Loc.End
	return jq, nil
}

// tryJSONQueryWrapper recognizes an optional jsonQueryWrapperBehavior prefix
// (`WITHOUT [ARRAY]` | `WITH [CONDITIONAL|UNCONDITIONAL] [ARRAY]`) WITHOUT
// consuming the trailing WRAPPER keyword (the caller consumes it). It returns
// ok=false when the current token is neither WITHOUT nor WITH. When ok, it has
// consumed the behavior tokens (WITHOUT/WITH, the optional conditionality, and
// the optional ARRAY) and returns the canonical wrapper spelling plus whether
// ARRAY was present.
//
// WITH/WITHOUT are reserved keywords that begin no other JSON_QUERY clause, so
// recognizing them here is unambiguous; the required WRAPPER keyword that must
// follow is enforced by the caller's expect.
func (p *Parser) tryJSONQueryWrapper() (wrapper string, array bool, ok bool) {
	switch p.cur.Kind {
	case kwWITHOUT:
		p.advance() // WITHOUT
		if p.cur.Kind == kwARRAY {
			p.advance()
			array = true
		}
		return "WITHOUT", array, true
	case kwWITH:
		p.advance() // WITH
		wrapper = "WITH"
		switch p.cur.Kind {
		case kwCONDITIONAL:
			p.advance()
			wrapper = "WITH CONDITIONAL"
		case kwUNCONDITIONAL:
			p.advance()
			wrapper = "WITH UNCONDITIONAL"
		}
		if p.cur.Kind == kwARRAY {
			p.advance()
			array = true
		}
		return wrapper, array, true
	default:
		return "", false, false
	}
}

// startsJSONQueryBehavior reports whether the current token begins a
// jsonQueryBehavior clause (ERROR | NULL | EMPTY …).
func (p *Parser) startsJSONQueryBehavior() bool {
	switch p.cur.Kind {
	case kwERROR, kwNULL, kwEMPTY:
		return true
	default:
		return false
	}
}

// parseJSONQueryBehaviorClause parses one `(ERROR | NULL | EMPTY (ARRAY|OBJECT))
// ON (EMPTY | ERROR)` clause (jsonQueryBehavior + its ON tail), returning the
// behavior and which slot (kwEMPTY or kwERROR) it fills. As with JSON_VALUE the
// slot is read from the trailing keyword so the two clauses disambiguate
// regardless of order; ordering is enforced by the caller (J1).
func (p *Parser) parseJSONQueryBehaviorClause() (*JSONQueryBehavior, TokenKind, error) {
	startTok := p.cur
	var b *JSONQueryBehavior
	switch p.cur.Kind {
	case kwERROR:
		p.advance()
		b = &JSONQueryBehavior{Kind: "ERROR", Loc: startTok.Loc}
	case kwNULL:
		p.advance()
		b = &JSONQueryBehavior{Kind: "NULL", Loc: startTok.Loc}
	case kwEMPTY:
		p.advance() // EMPTY
		var kind string
		switch p.cur.Kind {
		case kwARRAY:
			p.advance()
			kind = "EMPTY ARRAY"
		case kwOBJECT:
			p.advance()
			kind = "EMPTY OBJECT"
		default:
			return nil, 0, p.exprErrorAt("expected ARRAY or OBJECT after EMPTY")
		}
		b = &JSONQueryBehavior{Kind: kind, Loc: ast.Loc{Start: startTok.Loc.Start, End: p.prev.Loc.End}}
	default:
		return nil, 0, p.exprErrorAt("expected ERROR, NULL, or EMPTY")
	}
	slot, err := p.expectOnSlot()
	if err != nil {
		return nil, 0, err
	}
	return b, slot, nil
}

// ---------------------------------------------------------------------------
// JSON_OBJECT
// ---------------------------------------------------------------------------

// parseJSONObject parses `JSON_OBJECT( (member (, member)* null-handling?
// uniqueness?)? (RETURNING type (FORMAT …)?)? )` (jsonObject). The leading
// JSON_OBJECT keyword is the current token.
func (p *Parser) parseJSONObject() (Expr, error) {
	startTok := p.advance() // consume JSON_OBJECT
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	jo := &JSONObjectExpr{Loc: ast.Loc{Start: startTok.Loc.Start}}

	// Members are present unless the next token closes the call or begins the
	// trailing RETURNING clause. (Both NULL ON NULL / ABSENT ON NULL and the
	// uniqueness clause may only FOLLOW at least one member, so an empty member
	// list cannot start with them — confirmed by the oracle: `JSON_OBJECT(NULL
	// ON NULL)` and `JSON_OBJECT(ABSENT ON NULL)` are SYNTAX_ERRORs.)
	if p.cur.Kind != int(')') && p.cur.Kind != kwRETURNING {
		first, err := p.parseJSONObjectMember()
		if err != nil {
			return nil, err
		}
		jo.Members = append(jo.Members, first)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			member, err := p.parseJSONObjectMember()
			if err != nil {
				return nil, err
			}
			jo.Members = append(jo.Members, member)
		}

		// (NULL ON NULL | ABSENT ON NULL)?
		if nh, ok, err := p.tryJSONNullHandling(); err != nil {
			return nil, err
		} else if ok {
			jo.NullHandling = nh
		}

		// (WITH UNIQUE KEYS? | WITHOUT UNIQUE KEYS?)?
		if p.cur.Kind == kwWITH || p.cur.Kind == kwWITHOUT {
			if p.cur.Kind == kwWITH {
				jo.Uniqueness = "WITH UNIQUE"
			} else {
				jo.Uniqueness = "WITHOUT UNIQUE"
			}
			p.advance() // WITH or WITHOUT
			if _, err := p.expect(kwUNIQUE); err != nil {
				return nil, err
			}
			if p.cur.Kind == kwKEYS {
				p.advance()
				jo.UniqueKeys = true
			}
		}
	}

	// (RETURNING type (FORMAT jsonRepresentation)?)?
	if p.cur.Kind == kwRETURNING {
		typ, format, err := p.parseJSONReturning()
		if err != nil {
			return nil, err
		}
		jo.Returning = typ
		jo.ReturningFormat = format
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	jo.Loc.End = closeTok.Loc.End
	return jo, nil
}

// parseJSONObjectMember parses one `KEY? expression VALUE jsonValueExpression`
// or `expression COLON jsonValueExpression` member (jsonObjectMember).
//
// Per J4, a leading `KEY` keyword is ambiguous with `KEY` used as the key
// expression (KEY is non-reserved). The reading is resolved by speculation: when
// the current token is KEY, try consuming it as the keyword and parsing
// `expression VALUE value`; if that does not fit (e.g. `KEY VALUE 1`, where KEY
// is really the key column), rewind and parse the member with KEY as the
// expression's first token.
func (p *Parser) parseJSONObjectMember() (JSONObjectMember, error) {
	// Speculative KEY-as-keyword reading.
	if p.cur.Kind == kwKEY {
		cp := p.checkpoint()
		keyTok := p.advance() // tentatively consume KEY as the keyword
		if member, ok := p.tryJSONObjectMemberBody(keyTok.Loc.Start, true); ok {
			return member, nil
		}
		p.restore(cp)
	}

	// No leading KEY keyword: KEY (if present) is the key expression's first
	// token. The member must use `: value` or `VALUE value`.
	start := p.cur.Loc.Start
	if member, ok := p.tryJSONObjectMemberBody(start, false); ok {
		return member, nil
	}
	return JSONObjectMember{}, p.exprErrorAt("expected a JSON_OBJECT member (key VALUE value or key : value)")
}

// tryJSONObjectMemberBody parses `expression (VALUE | COLON) jsonValueExpression`
// after any leading KEY keyword has been handled, returning ok=false (leaving
// the cursor for the caller's restore) when the shape does not fit. keyKeyword
// records whether the caller consumed a leading KEY; start is the member's span
// start. A hard error inside the value (after the separator commits the reading)
// is surfaced via ok=false so the outer parser reports a single member error;
// the separators here are mandatory, so a missing separator simply means "not
// this reading."
func (p *Parser) tryJSONObjectMemberBody(start int, keyKeyword bool) (JSONObjectMember, bool) {
	key, err := p.parseExpr()
	if err != nil {
		return JSONObjectMember{}, false
	}
	switch p.cur.Kind {
	case kwVALUE:
		p.advance() // consume VALUE
		value, err := p.parseJSONValueExpr()
		if err != nil {
			return JSONObjectMember{}, false
		}
		return JSONObjectMember{
			KeyKeyword: keyKeyword,
			Key:        key,
			Separator:  "VALUE",
			Value:      value,
			Loc:        ast.Loc{Start: start, End: value.Loc.End},
		}, true
	case int(':'):
		// The colon form has no leading KEY keyword in the grammar; if the caller
		// consumed KEY this reading is invalid (forces a rewind to the non-keyword
		// path, which will re-parse KEY as part of the key expression).
		if keyKeyword {
			return JSONObjectMember{}, false
		}
		p.advance() // consume ':'
		value, err := p.parseJSONValueExpr()
		if err != nil {
			return JSONObjectMember{}, false
		}
		return JSONObjectMember{
			Key:       key,
			Separator: ":",
			Value:     value,
			Loc:       ast.Loc{Start: start, End: value.Loc.End},
		}, true
	default:
		return JSONObjectMember{}, false
	}
}

// ---------------------------------------------------------------------------
// JSON_ARRAY
// ---------------------------------------------------------------------------

// parseJSONArray parses `JSON_ARRAY( (jsonValueExpression (, jsonValueExpression)*
// null-handling?)? (RETURNING type (FORMAT …)?)? )` (jsonArray). Per J5 the
// element list is parsed first (a leading NULL is an element), then the
// null-handling clause is recognized only after the list. The leading JSON_ARRAY
// keyword is the current token.
func (p *Parser) parseJSONArray() (Expr, error) {
	startTok := p.advance() // consume JSON_ARRAY
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	ja := &JSONArrayExpr{Loc: ast.Loc{Start: startTok.Loc.Start}}

	// Elements are present unless the next token closes the call or begins the
	// trailing RETURNING clause.
	if p.cur.Kind != int(')') && p.cur.Kind != kwRETURNING {
		first, err := p.parseJSONValueExpr()
		if err != nil {
			return nil, err
		}
		ja.Elements = append(ja.Elements, first)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			elem, err := p.parseJSONValueExpr()
			if err != nil {
				return nil, err
			}
			ja.Elements = append(ja.Elements, elem)
		}

		// (NULL ON NULL | ABSENT ON NULL)?
		if nh, ok, err := p.tryJSONNullHandling(); err != nil {
			return nil, err
		} else if ok {
			ja.NullHandling = nh
		}
	}

	// (RETURNING type (FORMAT jsonRepresentation)?)?
	if p.cur.Kind == kwRETURNING {
		typ, format, err := p.parseJSONReturning()
		if err != nil {
			return nil, err
		}
		ja.Returning = typ
		ja.ReturningFormat = format
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	ja.Loc.End = closeTok.Loc.End
	return ja, nil
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// parseJSONPathInvocation parses `jsonValueExpression , path = string_
// (PASSING jsonArgument (, jsonArgument)*)?` (jsonPathInvocation): the JSON input
// value expression, a mandatory comma, the path STRING LITERAL (J2), and an
// optional PASSING argument list. The opening '(' has already been consumed.
func (p *Parser) parseJSONPathInvocation() (*JSONPathInvocation, error) {
	input, err := p.parseJSONValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(',')); err != nil {
		return nil, err
	}

	// path = string_ — a string LITERAL, not an arbitrary expression (J2).
	if p.cur.Kind != tokString && p.cur.Kind != tokUnicodeString {
		return nil, p.exprErrorAt("expected a string literal JSON path")
	}
	pathTok := p.advance()

	pi := &JSONPathInvocation{
		Input: input,
		Path:  pathTok.Str,
		Loc:   ast.Loc{Start: input.Loc.Start, End: pathTok.Loc.End},
	}

	// (PASSING jsonArgument (, jsonArgument)*)?
	if p.cur.Kind == kwPASSING {
		p.advance() // consume PASSING
		first, err := p.parseJSONArgument()
		if err != nil {
			return nil, err
		}
		pi.Passing = append(pi.Passing, first)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			arg, err := p.parseJSONArgument()
			if err != nil {
				return nil, err
			}
			pi.Passing = append(pi.Passing, arg)
		}
		pi.Loc.End = pi.Passing[len(pi.Passing)-1].Loc.End
	}

	return pi, nil
}

// parseJSONValueExpr parses `expression (FORMAT jsonRepresentation)?`
// (jsonValueExpression): a full expression optionally tagged with a JSON FORMAT.
func (p *Parser) parseJSONValueExpr() (*JSONValueExpr, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	jve := &JSONValueExpr{Expr: expr, Loc: ast.Loc{Start: expr.Span().Start, End: expr.Span().End}}
	if p.cur.Kind == kwFORMAT {
		format, err := p.parseJSONFormat()
		if err != nil {
			return nil, err
		}
		jve.Format = format
		jve.Loc.End = format.Loc.End
	}
	return jve, nil
}

// parseJSONFormat parses `FORMAT JSON (ENCODING (UTF8|UTF16|UTF32))?`
// (jsonRepresentation). Only the JSON representation keyword is accepted by
// Trino 481 (e.g. `FORMAT XML` is a SYNTAX_ERROR). FORMAT is the current token.
func (p *Parser) parseJSONFormat() (*JSONFormat, error) {
	formatTok := p.advance() // consume FORMAT
	// The representation keyword must be JSON (the non-reserved kwJSON token).
	if p.cur.Kind != kwJSON {
		return nil, p.exprErrorAt("expected JSON after FORMAT")
	}
	jsonTok := p.advance() // consume JSON
	jf := &JSONFormat{Loc: ast.Loc{Start: formatTok.Loc.Start, End: jsonTok.Loc.End}}

	if p.cur.Kind == kwENCODING {
		p.advance() // consume ENCODING
		switch p.cur.Kind {
		case kwUTF8:
			jf.Encoding = "UTF8"
		case kwUTF16:
			jf.Encoding = "UTF16"
		case kwUTF32:
			jf.Encoding = "UTF32"
		default:
			return nil, p.exprErrorAt("expected UTF8, UTF16, or UTF32 after ENCODING")
		}
		encTok := p.advance()
		jf.Loc.End = encTok.Loc.End
	}
	return jf, nil
}

// parseJSONArgument parses one `jsonValueExpression AS identifier` PASSING
// argument (jsonArgument).
func (p *Parser) parseJSONArgument() (JSONArgument, error) {
	value, err := p.parseJSONValueExpr()
	if err != nil {
		return JSONArgument{}, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return JSONArgument{}, err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return JSONArgument{}, err
	}
	return JSONArgument{
		Value: value,
		Name:  name,
		Loc:   ast.Loc{Start: value.Loc.Start, End: name.Loc.End},
	}, nil
}

// parseJSONReturning parses `RETURNING type (FORMAT jsonRepresentation)?`, the
// shared trailing clause of JSON_OBJECT and JSON_ARRAY. RETURNING is the current
// token. Returns the type and the optional FORMAT (nil when absent).
func (p *Parser) parseJSONReturning() (*DataType, *JSONFormat, error) {
	p.advance() // consume RETURNING
	typ, err := p.parseType()
	if err != nil {
		return nil, nil, err
	}
	var format *JSONFormat
	if p.cur.Kind == kwFORMAT {
		format, err = p.parseJSONFormat()
		if err != nil {
			return nil, nil, err
		}
	}
	return typ, format, nil
}

// tryJSONNullHandling recognizes an optional `NULL ON NULL | ABSENT ON NULL`
// clause shared by JSON_OBJECT and JSON_ARRAY, returning ok=false (consuming
// nothing) when the current tokens do not begin such a clause. The clause is
// detected by its two-token prefix (`NULL ON` / `ABSENT ON`); once that prefix
// is consumed the trailing `NULL` is mandatory — Trino 481 rejects `… NULL ON
// ERROR` / `… NULL ON` (oracle-confirmed), so a missing or wrong trailing token
// is a hard error here rather than a rewind.
func (p *Parser) tryJSONNullHandling() (string, bool, error) {
	var canonical string
	switch p.cur.Kind {
	case kwNULL:
		if p.peekNext().Kind != kwON {
			return "", false, nil
		}
		canonical = "NULL ON NULL"
	case kwABSENT:
		if p.peekNext().Kind != kwON {
			return "", false, nil
		}
		canonical = "ABSENT ON NULL"
	default:
		return "", false, nil
	}
	p.advance() // NULL or ABSENT
	p.advance() // ON
	if _, err := p.expect(kwNULL); err != nil {
		return "", false, err
	}
	return canonical, true, nil
}
