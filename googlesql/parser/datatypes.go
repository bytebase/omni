package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
)

// This file is the `types` DAG node: it implements GoogleSQL's `type` grammar
// rule and its sub-rules (raw_type, type_name, array_type, struct_type,
// struct_field, range_type, map_type, function_type, opt_type_parameters,
// type_parameter, collate_clause) as a hand-written recursive-descent parser
// over the token stream, producing a *DataType value tree.
//
// Legacy ZetaSQL ANTLR `type` rule (GoogleSQLParser.g4):
//
//	type:        raw_type opt_type_parameters? collate_clause? ;
//	raw_type:    array_type | struct_type | type_name | range_type
//	             | function_type | map_type ;
//	type_name:   path_expression | INTERVAL ;
//	array_type:  ARRAY '<' type '>' ;
//	struct_type: STRUCT '<' '>' | STRUCT '<' struct_field (',' struct_field)* '>' ;
//	struct_field: identifier type | type ;
//	range_type:  RANGE '<' type '>' ;
//	map_type:    MAP '<' type ',' type '>' ;
//	function_type: FUNCTION '<' '(' ')' '->' type '>'
//	             | FUNCTION '<' type '->' type '>'
//	             | FUNCTION '<' '(' type (',' type)* ')' '->' type '>' ;
//	opt_type_parameters: '(' type_parameter (',' type_parameter)* ')'  (trailing comma -> error) ;
//	type_parameter: integer | boolean | string | bytes | float | MAX ;
//	collate_clause: COLLATE string_literal_or_parameter ;
//
// DataType is a PARSER-PACKAGE node, NOT an ast.Node — the googlesql ast tag
// set is closed to the ast-core File container (ast/parsenodes.go ships only
// *File; adding type tags belongs to ast-core, not this node's writes-glob).
// This follows the trino/parser.DataType precedent exactly: later DAG nodes
// (expressions: CAST/typed literals/struct constructors; parser-ddl: column
// schemas, function params/returns; parser-dml) embed *DataType in their ast
// results and render it via String. The node carries its own ast.Loc span.
//
// The implementation is adjudicated against the LIVE Cloud Spanner emulator
// oracle (oracle.md), not the literal legacy grammar. Two oracle-confirmed
// divergences from a naive reading of the ANTLR rule are baked in (see the
// migration divergence ledger):
//
//	D1 (template-close token splitting). The omni lexer greedily emits `>>`
//	   (tokShiftRight) and `>=` (tokGreaterEqual); the legacy grammar's
//	   `template_type_close: GT_OPERATOR` matches only a single `>`, so a
//	   literal reading would REJECT nested types like ARRAY<ARRAY<INT64>> and
//	   STRUCT<x ARRAY<STRUCT<y INT64>>>. The live oracle ACCEPTS them (semantic
//	   "Arrays of arrays are not supported" — it PARSED). So a `>>` (or `>=`)
//	   token at a template-close position is SPLIT: the first `>` closes the
//	   current template and the remaining `>`/`=` is pushed back as a fresh
//	   token. A leftover `>` where no template is open is left in place and
//	   surfaces as a syntax error (oracle: ARRAY<INT64>> rejects).
//
//	D2 (empty STRUCT<>). Written adjacently, `<>` lexes as one token
//	   (tokNotEqual2). The oracle ACCEPTS the empty struct STRUCT<> (and the
//	   spaced STRUCT< >) but REJECTS ARRAY<>/RANGE<> (those need content). So
//	   STRUCT treats a leading `<>` (or `< >`) as an empty template; ARRAY/RANGE/
//	   MAP/FUNCTION require a real `<` followed by content.
//
// Dispatch on the leading keyword (raw_type alternatives):
//   - ARRAY / STRUCT / RANGE are RESERVED keywords (keywords.go) and can ONLY
//     introduce their template form — they are NOT valid bare type_names
//     (oracle: bare ARRAY / STRUCT / RANGE all reject "Expected < ...").
//   - MAP / FUNCTION are NON-reserved: they take the template form when
//     followed by `<`, else fall through to a bare type_name (oracle: bare MAP /
//     FUNCTION accept as type names "Type not found: MAP").
//   - INTERVAL is the dedicated type_name alternative (no template, no path
//     continuation — oracle: INTERVAL<...> and INTERVAL.foo both reject).
//   - anything else (tokIdentifier, or a non-reserved keyword like
//     NUMERIC/JSON/DATE/TIMESTAMP) is a type_name path_expression.

// TypeKind classifies the top-level shape of a parsed DataType, mirroring the
// raw_type alternatives so downstream consumers (expressions, DDL, deparse) can
// switch on the form without re-inspecting tokens.
type TypeKind int

const (
	// TypeName is a type_name: a path_expression (INT64, STRING, NUMERIC,
	// foo.bar.Baz, `proto.Type`) or the bare INTERVAL keyword (IsInterval set).
	// Covers all scalars, parameterized scalars (STRING(N), NUMERIC(p,s)), and
	// proto/enum type names — GoogleSQL scalar type names are identifiers, not
	// reserved keywords.
	TypeName TypeKind = iota
	// TypeArray is ARRAY < type >. The element is in ElementType.
	TypeArray
	// TypeStruct is STRUCT < [field (, field)*] >. Fields holds the (possibly
	// empty) field list.
	TypeStruct
	// TypeRange is RANGE < type >. The element is in ElementType.
	TypeRange
	// TypeMap is MAP < keyType , valueType >. KeyType / ValueType hold the
	// components (a ZetaSQL/legacy raw_type alt; Spanner rejects it semantically).
	TypeMap
	// TypeFunction is FUNCTION < (args) -> returnType >. ArgTypes / ReturnType
	// hold the signature (a ZetaSQL/legacy raw_type alt; semantically
	// unsupported on Spanner but grammatically valid).
	TypeFunction
)

// String returns a human-readable tag for the kind (diagnostics only).
func (k TypeKind) String() string {
	switch k {
	case TypeName:
		return "Name"
	case TypeArray:
		return "Array"
	case TypeStruct:
		return "Struct"
	case TypeRange:
		return "Range"
	case TypeMap:
		return "Map"
	case TypeFunction:
		return "Function"
	default:
		return "Unknown"
	}
}

// TypeParam is one entry of an opt_type_parameters list. Exactly one of
// IsInt / IsMax / Text-bearing literal kinds is meaningful:
//
//	type_parameter: integer | boolean | string | bytes | float | MAX
//
// IsInt + IntVal capture an integer parameter (the 100 in STRING(100), the 38
// in NUMERIC(38,9)); IsMax marks the MAX keyword (STRING(MAX)); Text holds the
// source spelling for the remaining literal kinds (boolean / string / bytes /
// float) so a faithful round-trip is possible without re-classifying them.
// IntText preserves an integer's exact source spelling (faithful even on int64
// overflow, where the lexer leaves IntVal at 0 but keeps the text).
type TypeParam struct {
	IsInt   bool
	IsMax   bool
	IntVal  int64
	IntText string // exact source spelling of an integer parameter
	Text    string // source spelling of a non-integer literal parameter
	Loc     ast.Loc
}

func (p TypeParam) writeTo(b *strings.Builder) {
	switch {
	case p.IsMax:
		b.WriteString("MAX")
	case p.IsInt:
		if p.IntText != "" {
			b.WriteString(p.IntText)
		} else {
			b.WriteString(strconv.FormatInt(p.IntVal, 10))
		}
	default:
		b.WriteString(p.Text)
	}
}

// StructField is one field of a STRUCT type. Name is the optional field-name
// identifier ("" for an anonymous field); Type is the field type (which may
// itself carry a trailing COLLATE captured on Type.Collate).
//
//	struct_field: identifier type | type
type StructField struct {
	Name string // "" for an anonymous field
	Type *DataType
	Loc  ast.Loc
}

func (f StructField) writeTo(b *strings.Builder) {
	if f.Name != "" {
		b.WriteString(f.Name)
		b.WriteByte(' ')
	}
	if f.Type != nil {
		f.Type.writeTo(b)
	}
}

// DataType is a parsed GoogleSQL data type (the `type` rule). The active fields
// depend on Kind:
//
//	TypeName     — NameParts (path), or IsInterval for the INTERVAL alt
//	TypeArray    — ElementType
//	TypeStruct   — Fields (possibly empty)
//	TypeRange    — ElementType
//	TypeMap      — KeyType, ValueType
//	TypeFunction — ArgTypes, ReturnType
//
// Params (opt_type_parameters) and Collate (collate_clause) may decorate ANY
// kind (the grammar appends both to `type` after raw_type), though in practice
// they appear on scalar type_names (STRING(100), NUMERIC(p,s) COLLATE '…').
type DataType struct {
	Kind TypeKind

	// NameParts is the dotted path of a TypeName (e.g. ["INT64"],
	// ["foo","bar","Baz"]). Source case is preserved (GoogleSQL type names are
	// case-insensitive for resolution, but the parser keeps the spelling).
	NameParts []string
	// IsInterval marks the type_name INTERVAL alternative (Kind==TypeName,
	// NameParts nil). INTERVAL is a reserved keyword and a type_name in its own
	// right, distinct from a path_expression.
	IsInterval bool

	// ElementType is the element of an ARRAY (TypeArray) or RANGE (TypeRange).
	ElementType *DataType

	// Fields holds the STRUCT field list (TypeStruct); empty for STRUCT<>.
	Fields []StructField

	// KeyType / ValueType are the components of MAP<k,v> (TypeMap).
	KeyType   *DataType
	ValueType *DataType

	// ArgTypes / ReturnType describe a FUNCTION<(args) -> ret> (TypeFunction);
	// ArgTypes is empty for the no-argument form FUNCTION<() -> ret>.
	ArgTypes   []*DataType
	ReturnType *DataType

	// Params holds an opt_type_parameters list (nil when absent).
	Params []TypeParam

	// Collate is the collation name/spelling from a trailing collate_clause
	// ("" when absent). For a string literal it is the unquoted body; for a
	// parameter (@name / ? / @@sysvar) it is the source spelling.
	Collate string

	Loc ast.Loc
}

// String renders the DataType back to GoogleSQL source syntax. It is the
// round-trip rendering downstream deparse relies on: it is not byte-identical
// to the input (unquoted-name case is preserved, inter-token spacing is
// normalized), but re-parsing String() yields an equal DataType.
func (t *DataType) String() string {
	if t == nil {
		return ""
	}
	var b strings.Builder
	t.writeTo(&b)
	return b.String()
}

func (t *DataType) writeTo(b *strings.Builder) {
	switch t.Kind {
	case TypeName:
		if t.IsInterval {
			b.WriteString("INTERVAL")
		} else {
			b.WriteString(strings.Join(t.NameParts, "."))
		}
	case TypeArray:
		b.WriteString("ARRAY<")
		t.ElementType.writeTo(b)
		b.WriteByte('>')
	case TypeStruct:
		b.WriteString("STRUCT<")
		for i, f := range t.Fields {
			if i > 0 {
				b.WriteString(", ")
			}
			f.writeTo(b)
		}
		b.WriteByte('>')
	case TypeRange:
		b.WriteString("RANGE<")
		t.ElementType.writeTo(b)
		b.WriteByte('>')
	case TypeMap:
		b.WriteString("MAP<")
		t.KeyType.writeTo(b)
		b.WriteString(", ")
		t.ValueType.writeTo(b)
		b.WriteByte('>')
	case TypeFunction:
		b.WriteString("FUNCTION<(")
		for i, a := range t.ArgTypes {
			if i > 0 {
				b.WriteString(", ")
			}
			a.writeTo(b)
		}
		b.WriteString(") -> ")
		t.ReturnType.writeTo(b)
		b.WriteByte('>')
	}
	if len(t.Params) > 0 {
		b.WriteByte('(')
		for i, p := range t.Params {
			if i > 0 {
				b.WriteString(", ")
			}
			p.writeTo(b)
		}
		b.WriteByte(')')
	}
	if t.Collate != "" {
		b.WriteString(" COLLATE ")
		// Render string collations quoted (the common case); a parameter
		// spelling (@p / ? / @@v) is written verbatim.
		if isParameterSpelling(t.Collate) {
			b.WriteString(t.Collate)
		} else {
			b.WriteByte('\'')
			b.WriteString(t.Collate)
			b.WriteByte('\'')
		}
	}
}

// isParameterSpelling reports whether s looks like a collate parameter (a query
// parameter @name, positional ?, or system variable @@var) rather than a
// string-literal collation name. Used only by String for round-trip rendering.
func isParameterSpelling(s string) bool {
	return s == "?" || strings.HasPrefix(s, "@")
}

// ---------------------------------------------------------------------------
// parseType — the `type` rule entry point.
// ---------------------------------------------------------------------------

// parseType parses a GoogleSQL type (`type: raw_type opt_type_parameters?
// collate_clause?`) and returns it as a *DataType. It is the single entry point
// for every type position — CAST / SAFE_CAST targets, typed literals
// (ARRAY<T>[…], STRUCT<…>(…)), RANGE<T> literals, function parameter/return
// types, variable declarations, and (via the DDL nodes) column schemas.
// Returns a *ParseError if the current token does not begin a valid type.
func (p *Parser) parseType() (*DataType, error) {
	dt, err := p.parseRawType()
	if err != nil {
		return nil, err
	}

	// opt_type_parameters — `( type_parameter (, type_parameter)* )`.
	if p.cur.Type == int('(') {
		params, endLoc, err := p.parseTypeParameters()
		if err != nil {
			return nil, err
		}
		dt.Params = params
		dt.Loc.End = endLoc.End
	}

	// collate_clause — `COLLATE string_literal_or_parameter`.
	if p.cur.Type == kwCOLLATE {
		collate, endLoc, err := p.parseCollateClause()
		if err != nil {
			return nil, err
		}
		dt.Collate = collate
		dt.Loc.End = endLoc.End
	}

	return dt, nil
}

// parseRawType parses one raw_type alternative (without the trailing
// opt_type_parameters / collate_clause, which parseType appends).
func (p *Parser) parseRawType() (*DataType, error) {
	switch p.cur.Type {
	case kwARRAY:
		return p.parseArrayType()
	case kwSTRUCT:
		return p.parseStructType()
	case kwRANGE:
		return p.parseRangeType()
	case kwMAP:
		// MAP is non-reserved: template form only when the NEXT token opens a
		// template `<`. Otherwise it is a bare type_name (oracle: bare MAP
		// accepts as "Type not found: MAP").
		if p.peekNext().Type == int('<') {
			return p.parseMapType()
		}
		return p.parseTypeName()
	case kwFUNCTION:
		if p.peekNext().Type == int('<') {
			return p.parseFunctionType()
		}
		return p.parseTypeName()
	default:
		return p.parseTypeName()
	}
}

// parseTypeName parses a type_name: `path_expression | INTERVAL`.
//
//	path_expression: identifier (DOT identifier)*
//	identifier:      token_identifier | keyword_as_identifier
//
// The first component accepts a token identifier OR a non-reserved word-keyword
// (NUMERIC, JSON, DATE, MAP, FUNCTION, …); reserved keywords are rejected as a
// bare type name. Dotted continuations accept any keyword as a part (matching
// the foundation's permissive name-part rule). INTERVAL is its own alternative
// with no continuation.
func (p *Parser) parseTypeName() (*DataType, error) {
	if p.cur.Type == kwINTERVAL {
		tok := p.advance()
		return &DataType{Kind: TypeName, IsInterval: true, Loc: tok.Loc}, nil
	}

	if !isIdentifierStart(p.cur.Type) {
		return nil, p.typeError()
	}
	first := p.advance()
	part0, err := p.identifierText(first)
	if err != nil {
		return nil, err
	}
	dt := &DataType{
		Kind:      TypeName,
		NameParts: []string{part0},
		Loc:       first.Loc,
	}
	for p.cur.Type == int('.') {
		p.advance() // consume '.'
		// A dotted name part may be any keyword or identifier (path components
		// after the first are permissive — common for proto/enum type names).
		if !isAnyKeywordIdentifier(p.cur.Type) {
			return nil, p.typeErrorAt("expected name after '.' in type")
		}
		partTok := p.advance()
		part, err := p.identifierText(partTok)
		if err != nil {
			return nil, err
		}
		dt.NameParts = append(dt.NameParts, part)
		dt.Loc.End = partTok.Loc.End
	}
	return dt, nil
}

// identifierText returns the source text of an identifier-or-keyword token: the
// raw Str for a token identifier (the lexer already stripped backticks), or the
// keyword spelling for a keyword-as-identifier (whose Str is empty).
//
// An empty `tokIdentifier` (the body of an empty backtick pair “) is rejected
// as an "invalid empty identifier" — the lexer admits empty backticks without a
// lex error, but the GoogleSQL grammar rejects them (oracle: `CAST(NULL AS “)`
// → "Syntax error: Invalid empty identifier"). Keyword tokens legitimately
// carry an empty Str, so the substitution is keyed on the token TYPE, not on
// emptiness, to avoid turning an empty identifier into the literal "IDENTIFIER".
func (p *Parser) identifierText(tok Token) (string, error) {
	if tok.Type == tokIdentifier {
		if tok.Str == "" {
			return "", &ParseError{Loc: tok.Loc, Msg: "invalid empty identifier"}
		}
		return tok.Str, nil
	}
	// A keyword-as-identifier: its Str is empty; use the keyword spelling.
	return TokenName(tok.Type), nil
}

// parseArrayType parses `ARRAY < type >` (caller confirmed the leading ARRAY).
func (p *Parser) parseArrayType() (*DataType, error) {
	arrTok := p.advance() // ARRAY
	if err := p.expectTemplateOpen(); err != nil {
		return nil, err
	}
	elem, err := p.parseType()
	if err != nil {
		return nil, err
	}
	endLoc, err := p.expectTemplateClose()
	if err != nil {
		return nil, err
	}
	return &DataType{
		Kind:        TypeArray,
		ElementType: elem,
		Loc:         ast.Loc{Start: arrTok.Loc.Start, End: endLoc.End},
	}, nil
}

// parseRangeType parses `RANGE < type >` (caller confirmed RANGE). The grammar
// admits ANY `type` inside RANGE (oracle: RANGE<INT64> parses, then fails
// semantically) — the DATE/DATETIME/TIMESTAMP restriction is a resolution
// concern, not syntax.
func (p *Parser) parseRangeType() (*DataType, error) {
	rangeTok := p.advance() // RANGE
	if err := p.expectTemplateOpen(); err != nil {
		return nil, err
	}
	elem, err := p.parseType()
	if err != nil {
		return nil, err
	}
	endLoc, err := p.expectTemplateClose()
	if err != nil {
		return nil, err
	}
	return &DataType{
		Kind:        TypeRange,
		ElementType: elem,
		Loc:         ast.Loc{Start: rangeTok.Loc.Start, End: endLoc.End},
	}, nil
}

// parseStructType parses a STRUCT type (caller confirmed STRUCT):
//
//	STRUCT '<' '>'                                  (empty)
//	STRUCT '<' struct_field (',' struct_field)* '>'
//
// D2: the empty form may arrive as the adjacent `<>` token (tokNotEqual2) or as
// a `< >` pair; both are accepted (oracle: STRUCT<> and STRUCT< > both parse).
func (p *Parser) parseStructType() (*DataType, error) {
	structTok := p.advance() // STRUCT

	// Empty struct via the adjacent `<>` token.
	if p.cur.Type == tokNotEqual2 {
		closeTok := p.advance()
		return &DataType{
			Kind: TypeStruct,
			Loc:  ast.Loc{Start: structTok.Loc.Start, End: closeTok.Loc.End},
		}, nil
	}

	if err := p.expectTemplateOpen(); err != nil {
		return nil, err
	}

	dt := &DataType{Kind: TypeStruct, Loc: structTok.Loc}

	// Empty struct via a `< >` pair: a template-close immediately after `<`.
	if p.atTemplateClose() {
		endLoc, err := p.expectTemplateClose()
		if err != nil {
			return nil, err
		}
		dt.Loc.End = endLoc.End
		return dt, nil
	}

	first, err := p.parseStructField()
	if err != nil {
		return nil, err
	}
	dt.Fields = append(dt.Fields, first)
	for p.cur.Type == int(',') {
		p.advance() // consume ','
		f, err := p.parseStructField()
		if err != nil {
			return nil, err
		}
		dt.Fields = append(dt.Fields, f)
	}
	endLoc, err := p.expectTemplateClose()
	if err != nil {
		return nil, err
	}
	dt.Loc.End = endLoc.End
	return dt, nil
}

// parseStructField parses one `struct_field: identifier type | type`.
//
// The two alternatives are ambiguous: a field name is itself an identifier, and
// a type begins with an identifier (a path_expression). Resolution follows the
// foundation's lookahead rule used elsewhere: a NAMED field is `name type`
// where `name` is a plain identifier-start token AND the token after it also
// begins a type (so the name is not itself the whole type). Concretely, the
// named form is taken when the current token is an identifier-start and the
// NEXT token begins a type and is not a field separator (',' / template-close)
// — otherwise the current token is the (single-component) type itself.
//
// This matches the oracle: STRUCT<x INT64> is named {x INT64}; STRUCT<INT64,
// STRING> is two anonymous types; STRUCT<x STRUCT<…>> names x; STRUCT<a
// ARRAY<…>> names a. A bare keyword type name (NUMERIC, DATE) in a field is the
// anonymous type, since the following token is a separator.
func (p *Parser) parseStructField() (StructField, error) {
	// Try the named form first: `identifier type`. The name must be a plain
	// identifier-start token, and what follows it must itself begin a type and
	// not be a field boundary.
	if isIdentifierStart(p.cur.Type) && beginsTypeFollowingName(p.peekNext().Type) {
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return StructField{}, err
		}
		fieldType, err := p.parseType()
		if err != nil {
			return StructField{}, err
		}
		return StructField{
			Name: name,
			Type: fieldType,
			Loc:  ast.Loc{Start: nameTok.Loc.Start, End: fieldType.Loc.End},
		}, nil
	}

	// Anonymous form: a bare `type`.
	fieldType, err := p.parseType()
	if err != nil {
		return StructField{}, err
	}
	return StructField{Type: fieldType, Loc: fieldType.Loc}, nil
}

// beginsTypeFollowingName reports whether tokenType can begin a struct field's
// TYPE when it follows a candidate field-name token. It is the disambiguator
// for `identifier type` vs a bare single-component `type`: the named form is
// only taken when the token after the name actually starts a type.
//
// A type begins with: ARRAY / STRUCT / RANGE / MAP / FUNCTION / INTERVAL
// keywords, a template-bearing name, a token identifier, or any non-reserved
// word-keyword (a scalar type_name). It does NOT begin with a field separator
// (',') or a template-close ('>', or the split '>>'/'>='), nor with '.' (which
// would be a path continuation of the name, handled inside parseTypeName).
func beginsTypeFollowingName(tokenType int) bool {
	switch tokenType {
	case kwARRAY, kwSTRUCT, kwRANGE, kwMAP, kwFUNCTION, kwINTERVAL:
		return true
	}
	return isIdentifierStart(tokenType)
}

// parseMapType parses `MAP < keyType , valueType >` (caller confirmed MAP '<').
func (p *Parser) parseMapType() (*DataType, error) {
	mapTok := p.advance() // MAP
	if err := p.expectTemplateOpen(); err != nil {
		return nil, err
	}
	keyType, err := p.parseType()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(',')); err != nil {
		return nil, err
	}
	valueType, err := p.parseType()
	if err != nil {
		return nil, err
	}
	endLoc, err := p.expectTemplateClose()
	if err != nil {
		return nil, err
	}
	return &DataType{
		Kind:      TypeMap,
		KeyType:   keyType,
		ValueType: valueType,
		Loc:       ast.Loc{Start: mapTok.Loc.Start, End: endLoc.End},
	}, nil
}

// parseFunctionType parses a FUNCTION type (caller confirmed FUNCTION '<'):
//
//	FUNCTION '<' '(' ')' '->' return '>'                 (no args)
//	FUNCTION '<' arg '->' return '>'                     (single bare arg)
//	FUNCTION '<' '(' type (',' type)* ')' '->' return '>'  (paren arg list)
func (p *Parser) parseFunctionType() (*DataType, error) {
	funcTok := p.advance() // FUNCTION
	if err := p.expectTemplateOpen(); err != nil {
		return nil, err
	}

	dt := &DataType{Kind: TypeFunction, Loc: funcTok.Loc}

	if p.cur.Type == int('(') {
		// Parenthesized argument list (possibly empty).
		p.advance() // consume '('
		if p.cur.Type != int(')') {
			arg, err := p.parseType()
			if err != nil {
				return nil, err
			}
			dt.ArgTypes = append(dt.ArgTypes, arg)
			for p.cur.Type == int(',') {
				p.advance() // consume ','
				a, err := p.parseType()
				if err != nil {
					return nil, err
				}
				dt.ArgTypes = append(dt.ArgTypes, a)
			}
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	} else {
		// Single bare argument type.
		arg, err := p.parseType()
		if err != nil {
			return nil, err
		}
		dt.ArgTypes = append(dt.ArgTypes, arg)
	}

	if _, err := p.expect(tokArrow); err != nil {
		return nil, err
	}
	ret, err := p.parseType()
	if err != nil {
		return nil, err
	}
	dt.ReturnType = ret

	endLoc, err := p.expectTemplateClose()
	if err != nil {
		return nil, err
	}
	dt.Loc.End = endLoc.End
	return dt, nil
}

// ---------------------------------------------------------------------------
// Template brackets — open `<`, close `>` (with the D1 `>>`/`>=` split).
// ---------------------------------------------------------------------------

// expectTemplateOpen consumes a template-open `<` (int('<')), returning a
// *ParseError otherwise. A `<>` (tokNotEqual2) at an ARRAY/RANGE/MAP/FUNCTION
// open is NOT a valid open (those require content) and is reported as a syntax
// error — matching the oracle reject of ARRAY<> / RANGE<>.
func (p *Parser) expectTemplateOpen() error {
	if p.cur.Type == int('<') {
		p.advance()
		return nil
	}
	return p.typeErrorAt("expected '<'")
}

// atTemplateClose reports whether the current token closes a template: a bare
// `>` (int('>')), or a `>>` / `>=` whose FIRST `>` closes this template (the
// D1 split). It does not consume anything.
func (p *Parser) atTemplateClose() bool {
	switch p.cur.Type {
	case int('>'), tokShiftRight, tokGreaterEqual:
		return true
	default:
		return false
	}
}

// expectTemplateClose consumes a single template-close `>` and returns its Loc.
//
// D1 split: the omni lexer greedily lexes `>>` (tokShiftRight) and `>=`
// (tokGreaterEqual). When the close is one of those compound tokens, only the
// leading `>` belongs to THIS template; the remainder is pushed back into the
// parser's lookahead as a fresh token (`>` for the second char of `>>`, `=` for
// `>=`) at the adjusted offset, so the enclosing template (or the trailing
// context) sees it next. The returned Loc covers just the consumed `>`.
//
// If the current token is not a close at all, a syntax error is returned. A
// leftover bare `>` with no enclosing template is left for the caller and
// surfaces downstream (oracle: ARRAY<INT64>> rejects with the stray `>`).
func (p *Parser) expectTemplateClose() (ast.Loc, error) {
	tok := p.cur
	switch tok.Type {
	case int('>'):
		p.advance()
		return tok.Loc, nil
	case tokShiftRight:
		// `>>` — consume one `>`, push back the second as a bare `>`.
		gtLoc := ast.Loc{Start: tok.Loc.Start, End: tok.Loc.Start + 1}
		p.replaceCurrent(Token{
			Type: int('>'),
			Loc:  ast.Loc{Start: tok.Loc.Start + 1, End: tok.Loc.End},
		})
		return gtLoc, nil
	case tokGreaterEqual:
		// `>=` — consume the `>`, push back the `=`.
		gtLoc := ast.Loc{Start: tok.Loc.Start, End: tok.Loc.Start + 1}
		p.replaceCurrent(Token{
			Type: int('='),
			Loc:  ast.Loc{Start: tok.Loc.Start + 1, End: tok.Loc.End},
		})
		return gtLoc, nil
	default:
		return ast.NoLoc(), p.typeErrorAt("expected '>'")
	}
}

// replaceCurrent overwrites the current token with tok WITHOUT touching the
// lexer or the lookahead buffer, used by the D1 template-close split to push
// back the second half of a compound `>>` / `>=` token. The lexer position is
// already past the original compound token, and any buffered lookahead
// (nextBuf/hasNext, set by a prior peekNext) is preserved unchanged — so the
// next advance() correctly yields the buffered token (or the next lexed one)
// after tok is consumed.
func (p *Parser) replaceCurrent(tok Token) {
	p.cur = tok
}

// ---------------------------------------------------------------------------
// opt_type_parameters — `( type_parameter (, type_parameter)* )`.
// ---------------------------------------------------------------------------

// parseTypeParameters parses an opt_type_parameters list. The opening '(' is
// the current token. At least one type_parameter is required (oracle: STRING()
// rejects "Unexpected )"). Returns the params and the Loc spanning through ')'.
//
//	type_parameter: integer | boolean | string | bytes | float | MAX
//
// A trailing comma is a syntax error (legacy grammar emits a dedicated message;
// oracle rejects equivalently).
func (p *Parser) parseTypeParameters() ([]TypeParam, ast.Loc, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, ast.NoLoc(), err
	}
	first, err := p.parseTypeParameter()
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	params := []TypeParam{first}
	for p.cur.Type == int(',') {
		p.advance() // consume ','
		// A `)` right after the comma is a trailing comma — reject.
		if p.cur.Type == int(')') {
			return nil, ast.NoLoc(), p.typeErrorAt("trailing comma in type parameters list is not allowed")
		}
		next, err := p.parseTypeParameter()
		if err != nil {
			return nil, ast.NoLoc(), err
		}
		params = append(params, next)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	return params, closeTok.Loc, nil
}

// parseTypeParameter parses one type_parameter. Accepts an integer literal, the
// MAX keyword, a boolean literal (TRUE/FALSE), a string literal, a bytes
// literal, or a floating-point literal. A leading sign is NOT part of an
// integer type_parameter (oracle: NUMERIC(-5) rejects "Unexpected -").
func (p *Parser) parseTypeParameter() (TypeParam, error) {
	tok := p.cur
	switch tok.Type {
	case tokInteger:
		p.advance()
		return TypeParam{IsInt: true, IntVal: tok.Ival, IntText: tok.Str, Loc: tok.Loc}, nil
	case kwMAX:
		p.advance()
		return TypeParam{IsMax: true, Loc: tok.Loc}, nil
	case kwTRUE, kwFALSE:
		p.advance()
		return TypeParam{Text: TokenName(tok.Type), Loc: tok.Loc}, nil
	case tokString:
		p.advance()
		return TypeParam{Text: tok.Str, Loc: tok.Loc}, nil
	case tokBytes:
		p.advance()
		return TypeParam{Text: tok.Str, Loc: tok.Loc}, nil
	case tokFloat:
		p.advance()
		return TypeParam{Text: tok.Str, Loc: tok.Loc}, nil
	default:
		return TypeParam{}, p.typeErrorAt("expected a type parameter (integer, boolean, string, bytes, float, or MAX)")
	}
}

// ---------------------------------------------------------------------------
// collate_clause — `COLLATE string_literal_or_parameter`.
// ---------------------------------------------------------------------------

// parseCollateClause parses `COLLATE string_literal_or_parameter`, where the
// argument is a string literal, a query parameter (@name | ?), or a system
// variable (@@path). The current token is COLLATE. Returns the collation
// spelling and the Loc through the argument.
func (p *Parser) parseCollateClause() (string, ast.Loc, error) {
	collateTok := p.advance() // COLLATE
	tok := p.cur
	switch tok.Type {
	case tokString:
		p.advance()
		return tok.Str, ast.Loc{Start: collateTok.Loc.Start, End: tok.Loc.End}, nil
	case int('?'):
		p.advance()
		return "?", ast.Loc{Start: collateTok.Loc.Start, End: tok.Loc.End}, nil
	case int('@'):
		// Named parameter `@identifier`.
		p.advance() // consume '@'
		if !isAnyKeywordIdentifier(p.cur.Type) {
			return "", ast.NoLoc(), p.typeErrorAt("expected parameter name after '@' in COLLATE")
		}
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return "", ast.NoLoc(), err
		}
		return "@" + name, ast.Loc{Start: collateTok.Loc.Start, End: nameTok.Loc.End}, nil
	case tokAtAt:
		// System variable `@@path`.
		p.advance() // consume '@@'
		if !isAnyKeywordIdentifier(p.cur.Type) {
			return "", ast.NoLoc(), p.typeErrorAt("expected system variable name after '@@' in COLLATE")
		}
		var parts []string
		nameTok := p.advance()
		name, err := p.identifierText(nameTok)
		if err != nil {
			return "", ast.NoLoc(), err
		}
		parts = append(parts, name)
		end := nameTok.Loc.End
		for p.cur.Type == int('.') {
			p.advance()
			if !isAnyKeywordIdentifier(p.cur.Type) {
				return "", ast.NoLoc(), p.typeErrorAt("expected name after '.' in system variable")
			}
			partTok := p.advance()
			part, err := p.identifierText(partTok)
			if err != nil {
				return "", ast.NoLoc(), err
			}
			parts = append(parts, part)
			end = partTok.Loc.End
		}
		return "@@" + strings.Join(parts, "."), ast.Loc{Start: collateTok.Loc.Start, End: end}, nil
	default:
		return "", ast.NoLoc(), p.typeErrorAt("expected a string or parameter after COLLATE")
	}
}

// ---------------------------------------------------------------------------
// Errors & standalone entry point.
// ---------------------------------------------------------------------------

// typeError returns a *ParseError describing a missing type at the current
// token, distinct from a generic error so a type position reports "expected
// type".
func (p *Parser) typeError() *ParseError {
	if p.cur.Type == tokEOF {
		return &ParseError{Loc: p.cur.Loc, Msg: "expected type, found end of input"}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Type)
	}
	return &ParseError{Loc: p.cur.Loc, Msg: "expected type, found " + text}
}

// typeErrorAt returns a *ParseError with a custom message at the current token.
func (p *Parser) typeErrorAt(msg string) *ParseError {
	return &ParseError{Loc: p.cur.Loc, Msg: msg}
}

// ParseDataType parses a complete GoogleSQL type from a standalone string,
// returning the *DataType and any ParseErrors. Trailing tokens after the type
// are reported as an error (the type is still returned). It is the string-input
// counterpart of parseType, for tests and callers (catalog, deparse) that hold
// a type string rather than a token stream. Mirrors snowflake/parser and
// trino/parser ParseDataType.
func ParseDataType(input string) (*DataType, []ParseError) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()

	dt, err := p.parseType()
	if err != nil {
		if pe, ok := err.(*ParseError); ok {
			return nil, []ParseError{*pe}
		}
		return nil, []ParseError{{Msg: err.Error()}}
	}
	if p.cur.Type != tokEOF {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Type)
		}
		return dt, []ParseError{{Loc: p.cur.Loc, Msg: "unexpected token after type: " + text}}
	}
	return dt, nil
}
