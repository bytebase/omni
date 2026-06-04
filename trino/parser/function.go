package parser

import "github.com/bytebase/omni/trino/ast"

// This file is part of the `expressions` DAG node (with expr.go and
// predicate.go): it implements the function-call, special-built-in, and
// keyword-led primaryExpression alternatives of Trino's grammar, plus the
// inline OVER / FILTER clauses that decorate a function call.
//
// Covered primaryExpression alternatives (TrinoParser.g4):
//
//	functionCall            : processingMode? qualifiedName LPAREN_
//	                            (label=identifier DOT_)? ASTERISK_ RPAREN_ filter? over?
//	                        | processingMode? qualifiedName LPAREN_
//	                            (setQuantifier? expression (COMMA_ expression)*)?
//	                            (ORDER_ BY_ sortItem (COMMA_ sortItem)*)? RPAREN_
//	                            filter? (nullTreatment? over)? ;
//	measure                 : identifier over ;
//	lambda                  : identifier RARROW_ expression
//	                        | LPAREN_ (identifier (COMMA_ identifier)*)? RPAREN_ RARROW_ expression ;
//	simpleCase / searchedCase, cast (CAST/TRY_CAST), exists,
//	position, substring, trim, normalize, extract, groupingOperation,
//	listagg, the specialDateTimeFunction / currentUser / currentCatalog /
//	currentSchema / currentPath set, typeConstructor (identifier string_ and
//	DOUBLE PRECISION string_).
//
// The SQL/JSON functions (JSON_EXISTS/JSON_VALUE/JSON_QUERY/JSON_OBJECT/
// JSON_ARRAY and the json-path subsystem) are the SEPARATE expr-json DAG node
// (json.go); parsePrimaryAtom routes them to parseJSONFunction, which lives in
// json.go. (During the expressions node this was a thin "not yet supported"
// stub here; the expr-json node replaced it with the real parser in json.go.)
//
// Oracle-confirmed reservation nuance baked into the dispatch: among the
// special built-ins, SUBSTRING and POSITION are NON-RESERVED keywords, so beyond
// their special syntax they also parse as a bare column reference (`substring`)
// and as ordinary comma-argument function calls (`substring('abc', 1, 2)`,
// `position('a', 'abc')`). TRIM and NORMALIZE are reserved but still accept an
// ordinary comma/single-argument call shape (`trim('abc')`, `normalize('abc')`),
// which their special parsers cover. EXTRACT, LISTAGG, and GROUPING are reserved
// and accept ONLY their special form (`extract('a','b')`, `listagg('a')`,
// `grouping('a')` are all SYNTAX_ERRORs). parseIdentifierOrFunction therefore
// routes SUBSTRING/POSITION through a try-special-then-fall-back path, while the
// reserved built-ins stay on their dedicated parsers.
//
// Node-scope decision B2 (file header of expr.go): a function's OVER frame may,
// per the legacy windowFrame rule, carry the full row-pattern subsystem
// (MEASURES / PATTERN / DEFINE / SUBSET / AFTER MATCH / INITIAL|SEEK). That
// subsystem is the parser-match-recognize node; here OVER parses the common
// window frame `(ROWS|RANGE|GROUPS) [BETWEEN bound AND] bound` only.

// ---------------------------------------------------------------------------
// Function-call & related node types
// ---------------------------------------------------------------------------

// FuncCall is a function invocation. Name is the (possibly qualified) function
// name. Distinct marks the `DISTINCT` setQuantifier; Star marks the `count(*)`
// form (with optional Label for `count(t.*)`). OrderBy carries an in-aggregate
// `ORDER BY`; Filter, NullTreatment, and Over carry the trailing decorators.
// ProcessingMode is "", "RUNNING", or "FINAL".
type FuncCall struct {
	Name           *ast.QualifiedName
	ProcessingMode string
	Distinct       bool
	Star           bool            // the count(*) form
	Label          *ast.Identifier // count(t.*) label, nil otherwise
	Args           []Expr          // ordinary argument list (nil for the star form)
	OrderBy        []SortItem      // in-aggregate ORDER BY, nil when absent
	Filter         Expr            // FILTER (WHERE …) condition, nil when absent
	NullTreatment  string          // "", "IGNORE NULLS", or "RESPECT NULLS"
	Over           *WindowSpec     // OVER clause, nil when absent
	OverName       *ast.Identifier // OVER windowName form, nil otherwise
	Loc            ast.Loc
}

func (n *FuncCall) Span() ast.Loc { return n.Loc }
func (*FuncCall) exprNode()       {}

// SortItem is one `expression [ASC|DESC] [NULLS FIRST|LAST]` of an ORDER BY
// (the sortItem rule), reused by in-aggregate ORDER BY and OVER's ORDER BY.
type SortItem struct {
	Expr      Expr
	Ordering  string // "", "ASC", or "DESC"
	NullOrder string // "", "FIRST", or "LAST"
	Loc       ast.Loc
}

// LambdaExpr is `param -> body` or `(param, …) -> body` (the lambda
// alternatives). Params is the (possibly empty) parameter identifier list.
type LambdaExpr struct {
	Params []*ast.Identifier
	Body   Expr
	Loc    ast.Loc
}

func (n *LambdaExpr) Span() ast.Loc { return n.Loc }
func (*LambdaExpr) exprNode()       {}

// CaseExpr is a simple or searched CASE (simpleCase / searchedCase). Operand is
// nil for searched CASE. Else is nil when no ELSE branch.
type CaseExpr struct {
	Operand Expr // nil for searched CASE
	Whens   []WhenClause
	Else    Expr // nil when absent
	Loc     ast.Loc
}

func (n *CaseExpr) Span() ast.Loc { return n.Loc }
func (*CaseExpr) exprNode()       {}

// WhenClause is one `WHEN condition THEN result` of a CASE expression.
type WhenClause struct {
	Cond   Expr
	Result Expr
	Loc    ast.Loc
}

// CastExpr is `CAST(expr AS type)` or `TRY_CAST(expr AS type)`. Try marks the
// TRY_CAST spelling; Type is the parsed target type (datatypes.go).
type CastExpr struct {
	Try  bool
	Expr Expr
	Type *DataType
	Loc  ast.Loc
}

func (n *CastExpr) Span() ast.Loc { return n.Loc }
func (*CastExpr) exprNode()       {}

// ExtractExpr is `EXTRACT(field FROM source)` (extract). Field is the unit
// identifier (YEAR, MONTH, …) kept as source text.
type ExtractExpr struct {
	Field  string
	Source Expr
	Loc    ast.Loc
}

func (n *ExtractExpr) Span() ast.Loc { return n.Loc }
func (*ExtractExpr) exprNode()       {}

// SubstringExpr is `SUBSTRING(source FROM start [FOR length])` (the SQL-standard
// substring spelling). Length is nil when no FOR clause.
type SubstringExpr struct {
	Source Expr
	From   Expr
	For    Expr // nil when absent
	Loc    ast.Loc
}

func (n *SubstringExpr) Span() ast.Loc { return n.Loc }
func (*SubstringExpr) exprNode()       {}

// TrimExpr is `TRIM([spec] [char] FROM source)` or `TRIM(source, char)` (the two
// trim spellings). Spec is "", "LEADING", "TRAILING", or "BOTH"; Char is the
// trim character (nil when absent).
type TrimExpr struct {
	Spec   string
	Char   Expr // nil when absent
	Source Expr
	Loc    ast.Loc
}

func (n *TrimExpr) Span() ast.Loc { return n.Loc }
func (*TrimExpr) exprNode()       {}

// NormalizeExpr is `NORMALIZE(source [, form])` (normalize). Form is "" or one
// of NFD/NFC/NFKD/NFKC.
type NormalizeExpr struct {
	Source Expr
	Form   string
	Loc    ast.Loc
}

func (n *NormalizeExpr) Span() ast.Loc { return n.Loc }
func (*NormalizeExpr) exprNode()       {}

// PositionExpr is `POSITION(needle IN haystack)` (position).
type PositionExpr struct {
	Needle   Expr
	Haystack Expr
	Loc      ast.Loc
}

func (n *PositionExpr) Span() ast.Loc { return n.Loc }
func (*PositionExpr) exprNode()       {}

// GroupingExpr is `GROUPING( qualifiedName , … )` (groupingOperation). The
// argument list may be empty.
type GroupingExpr struct {
	Args []*ast.QualifiedName
	Loc  ast.Loc
}

func (n *GroupingExpr) Span() ast.Loc { return n.Loc }
func (*GroupingExpr) exprNode()       {}

// ListaggExpr is `LISTAGG([DISTINCT] arg [, sep] [ON OVERFLOW …]) WITHIN GROUP
// (ORDER BY …) [FILTER (WHERE …)]` (listagg). Only the syntactic shape is
// captured; semantics are downstream.
type ListaggExpr struct {
	Distinct      bool
	Arg           Expr
	Separator     Expr   // nil when absent
	OnOverflow    string // "", "ERROR", or "TRUNCATE"
	WithinGroupBy []SortItem
	Filter        Expr // FILTER (WHERE …), nil when absent
	Loc           ast.Loc
}

func (n *ListaggExpr) Span() ast.Loc { return n.Loc }
func (*ListaggExpr) exprNode()       {}

// SpecialFuncExpr is one of the no-argument-or-precision special functions:
// CURRENT_DATE / CURRENT_TIME[(p)] / CURRENT_TIMESTAMP[(p)] / LOCALTIME[(p)] /
// LOCALTIMESTAMP[(p)] / CURRENT_USER / CURRENT_CATALOG / CURRENT_SCHEMA /
// CURRENT_PATH (specialDateTimeFunction + currentUser/Catalog/Schema/Path).
// Name is the keyword; Precision is the optional `(INTEGER)` argument (-1 when
// absent or not allowed).
type SpecialFuncExpr struct {
	Name      string
	Precision int // -1 when absent
	Loc       ast.Loc
}

func (n *SpecialFuncExpr) Span() ast.Loc { return n.Loc }
func (*SpecialFuncExpr) exprNode()       {}

// TypeConstructorExpr alias note: the `identifier string_` typeConstructor is
// modelled by TypeConstructor in expr.go (parseDoublePrecisionConstructor and
// the identifier path build it).

// WindowSpec is an inline OVER `( [name] [PARTITION BY …] [ORDER BY …] [frame] )`
// (windowSpecification). ExistingName is the referenced base window name (nil
// when absent).
type WindowSpec struct {
	ExistingName *ast.Identifier
	PartitionBy  []Expr
	OrderBy      []SortItem
	Frame        *WindowFrame // nil when absent
	Loc          ast.Loc
}

// WindowFrame is the common frame `frameType [BETWEEN start AND end] | frameType
// start` (frameExtent). FrameType is ROWS/RANGE/GROUPS. End is nil for the
// single-bound form. The pattern-recognition frame additions are deferred to
// parser-match-recognize (B2).
type WindowFrame struct {
	FrameType string
	Start     WindowBound
	End       *WindowBound // nil for the single-bound form
	Loc       ast.Loc
}

// WindowBound is one frame bound (frameBound). Kind is the bound shape; Value is
// the offset expression for the bounded forms (nil otherwise).
type WindowBound struct {
	Kind  WindowBoundKind
	Value Expr // for PRECEDING/FOLLOWING with an offset
	Loc   ast.Loc
}

// WindowBoundKind enumerates the frame-bound shapes.
type WindowBoundKind int

const (
	// BoundUnboundedPreceding is UNBOUNDED PRECEDING.
	BoundUnboundedPreceding WindowBoundKind = iota
	// BoundUnboundedFollowing is UNBOUNDED FOLLOWING.
	BoundUnboundedFollowing
	// BoundCurrentRow is CURRENT ROW.
	BoundCurrentRow
	// BoundPreceding is `expr PRECEDING`.
	BoundPreceding
	// BoundFollowing is `expr FOLLOWING`.
	BoundFollowing
)

// ---------------------------------------------------------------------------
// identifier-or-function dispatch
// ---------------------------------------------------------------------------

// parseIdentifierOrFunction parses a primaryExpression that begins with an
// identifier (or a non-reserved keyword usable as one). The cases, in the order
// they are disambiguated:
//
//   - bare-identifier lambda `x -> body` (lambda);
//   - typeConstructor `name 'string'` (an identifier immediately followed by a
//     string literal);
//   - measure `identifier OVER (…)` (a bare identifier followed by OVER);
//   - function call `qualifiedName ( … )` (a name followed by '(');
//   - otherwise a column reference / dotted name (columnReference + dereference,
//     the dereference handled by parsePostfix on the leading ColumnRef).
//
// A multi-part name (a.b.c) is read as a qualifiedName so `a.b.c(1)` is a single
// function call; a trailing function-less name yields a ColumnRef for the first
// part with the remaining parts becoming Dereferences via parsePostfix — except
// the qualifiedName is consumed whole when a '(' follows.
func (p *Parser) parseIdentifierOrFunction() (Expr, error) {
	// Bare single-identifier lambda: `x -> body`. Detected before consuming the
	// name so the identifier becomes the sole lambda parameter.
	if p.peekNext().Kind == tokArrow {
		nameTok := p.advance() // the parameter identifier
		param := identFromToken(nameTok)
		p.advance() // consume ->
		body, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &LambdaExpr{
			Params: []*ast.Identifier{param},
			Body:   body,
			Loc:    ast.Loc{Start: param.Loc.Start, End: body.Span().End},
		}, nil
	}

	// typeConstructor: an identifier directly followed by a string literal, e.g.
	// `DATE '2020-01-01'`, `json '[1]'`, `bigint '7'`. A non-reserved keyword
	// name (DATE is non-reserved) is accepted as the constructor type name.
	if nk := p.peekNext().Kind; nk == tokString || nk == tokUnicodeString {
		nameTok := p.advance()
		name := identFromToken(nameTok)
		strTok := p.advance()
		return &TypeConstructor{
			Name:  name.String(),
			Value: strTok.Str,
			Loc:   ast.Loc{Start: nameTok.Loc.Start, End: strTok.Loc.End},
		}, nil
	}

	// measure: `identifier OVER (…)` — a bare identifier (single part) decorated
	// with an OVER window. Only applies when OVER directly follows the name.
	if p.peekNext().Kind == kwOVER {
		nameTok := p.advance()
		name := identFromToken(nameTok)
		over, overName, err := p.parseOver()
		if err != nil {
			return nil, err
		}
		end := name.Loc.End
		if over != nil {
			end = over.Loc.End
		} else if overName != nil {
			end = overName.Loc.End
		}
		// Model a measure as a degenerate FuncCall carrying only the OVER (no
		// args, no parens) so downstream consumers handle it uniformly.
		return &FuncCall{
			Name:     &ast.QualifiedName{Parts: []*ast.Identifier{name}, Loc: name.Loc},
			Over:     over,
			OverName: overName,
			Loc:      ast.Loc{Start: name.Loc.Start, End: end},
		}, nil
	}

	// SUBSTRING and POSITION are non-reserved keywords with BOTH a special syntax
	// (`SUBSTRING(s FROM a [FOR b])`, `POSITION(x IN y)`) AND ordinary readings —
	// as a bare column reference (`SELECT substring`) or a comma-argument function
	// call (`substring('abc', 1, 2)`, `position('a', 'abc')`), all oracle-accepted.
	// Try the special form speculatively when '(' follows; on failure fall back to
	// the ordinary identifier/function path.
	if (p.cur.Kind == kwSUBSTRING || p.cur.Kind == kwPOSITION) && p.peekNext().Kind == int('(') {
		isSubstring := p.cur.Kind == kwSUBSTRING
		cp := p.checkpoint()
		var (
			special Expr
			err     error
		)
		if isSubstring {
			special, err = p.parseSubstring()
		} else {
			special, err = p.parsePosition()
		}
		if err == nil {
			return special, nil
		}
		// Not the special form (e.g. a comma-argument call) — rewind and parse as
		// an ordinary function call.
		p.restore(cp)
	}

	// Read the (possibly qualified) name. A '(' immediately after makes it a
	// function call; otherwise it is a column reference / dotted name.
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if p.cur.Kind == int('(') {
		return p.parseFunctionCall(name, "")
	}
	return p.columnRefFromName(name), nil
}

// columnRefFromName turns a parsed qualifiedName into a ColumnRef (first part)
// wrapped in Dereferences for any trailing parts, so `a.b.c` (no call) becomes
// Dereference(Dereference(ColumnRef a, b), c). This keeps the dotted-name and
// the row-field-access readings structurally identical, matching Trino where the
// two are resolved the same way after parsing.
func (p *Parser) columnRefFromName(name *ast.QualifiedName) Expr {
	var expr Expr = &ColumnRef{Name: name.Parts[0], Loc: name.Parts[0].Loc}
	for _, part := range name.Parts[1:] {
		expr = &Dereference{
			Base:      expr,
			FieldName: part,
			Loc:       ast.Loc{Start: expr.Span().Start, End: part.Loc.End},
		}
	}
	return expr
}

// parseProcessingModeFunction handles a `RUNNING`/`FINAL` processingMode prefix
// on a function call (functionCall's processingMode?). The mode keyword has been
// detected by the caller; the function name and arguments follow.
func (p *Parser) parseProcessingModeFunction() (Expr, error) {
	modeTok := p.advance() // RUNNING or FINAL
	mode := "RUNNING"
	if modeTok.Kind == kwFINAL {
		mode = "FINAL"
	}
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if p.cur.Kind != int('(') {
		return nil, p.exprErrorAt("expected '(' after " + mode + " function name")
	}
	return p.parseFunctionCall(name, mode)
}

// parseFunctionCall parses `qualifiedName ( … ) filter? (nullTreatment? over)?`
// (both functionCall alternatives). The opening '(' is the current token. The
// body is either the star form `(label.)? *` or an argument list optionally
// prefixed by a setQuantifier and suffixed by an in-aggregate `ORDER BY`.
// processingMode is "" for the common case.
func (p *Parser) parseFunctionCall(name *ast.QualifiedName, processingMode string) (Expr, error) {
	p.advance() // consume '('

	fc := &FuncCall{
		Name:           name,
		ProcessingMode: processingMode,
		Loc:            ast.Loc{Start: name.Loc.Start},
	}

	// Star form: `*` or `label.*`.
	if p.cur.Kind == int('*') {
		p.advance() // consume '*'
		fc.Star = true
	} else if isIdentifierStart(p.cur.Kind) && p.peekNext().Kind == int('.') {
		// Possible `label.*`. Speculate: a single identifier, a dot, then '*'.
		cp := p.checkpoint()
		labelTok := p.advance() // identifier
		p.advance()             // '.'
		if p.cur.Kind == int('*') {
			p.advance() // consume '*'
			fc.Star = true
			fc.Label = identFromToken(labelTok)
		} else {
			p.restore(cp)
		}
	}

	if !fc.Star {
		// setQuantifier? expression (, expression)*
		if p.cur.Kind == kwDISTINCT {
			p.advance()
			fc.Distinct = true
		} else if p.cur.Kind == kwALL {
			p.advance() // ALL is the default; recorded by absence of Distinct
		}
		if p.cur.Kind != int(')') {
			first, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			fc.Args = append(fc.Args, first)
			for p.cur.Kind == int(',') {
				p.advance() // consume ','
				next, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				fc.Args = append(fc.Args, next)
			}
		}
		// in-aggregate ORDER BY sortItem (, sortItem)*
		if p.cur.Kind == kwORDER {
			items, err := p.parseOrderByItems()
			if err != nil {
				return nil, err
			}
			fc.OrderBy = items
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	fc.Loc.End = closeTok.Loc.End

	// filter?
	if p.cur.Kind == kwFILTER {
		filter, end, err := p.parseFilter()
		if err != nil {
			return nil, err
		}
		fc.Filter = filter
		fc.Loc.End = end
	}

	// (nullTreatment? over)? — nullTreatment is valid only with an OVER.
	if nt, ok := p.peekNullTreatment(); ok {
		fc.NullTreatment = nt
		p.consumeNullTreatment()
	}
	if p.cur.Kind == kwOVER {
		over, overName, err := p.parseOver()
		if err != nil {
			return nil, err
		}
		fc.Over = over
		fc.OverName = overName
		if over != nil {
			fc.Loc.End = over.Loc.End
		} else if overName != nil {
			fc.Loc.End = overName.Loc.End
		}
	}

	return fc, nil
}

// peekNullTreatment reports whether the current two tokens are an `IGNORE NULLS`
// or `RESPECT NULLS` nullTreatment, returning its canonical spelling.
func (p *Parser) peekNullTreatment() (string, bool) {
	if (p.cur.Kind == kwIGNORE || p.cur.Kind == kwRESPECT) && p.peekNext().Kind == kwNULLS {
		if p.cur.Kind == kwIGNORE {
			return "IGNORE NULLS", true
		}
		return "RESPECT NULLS", true
	}
	return "", false
}

// consumeNullTreatment consumes a confirmed `(IGNORE|RESPECT) NULLS`.
func (p *Parser) consumeNullTreatment() {
	p.advance() // IGNORE or RESPECT
	p.advance() // NULLS
}

// parseFilter parses `FILTER ( WHERE booleanExpression )` (the filter rule),
// returning the condition and the closing-')' end offset. FILTER is the current
// token.
func (p *Parser) parseFilter() (Expr, int, error) {
	p.advance() // consume FILTER
	if _, err := p.expect(int('(')); err != nil {
		return nil, 0, err
	}
	if _, err := p.expect(kwWHERE); err != nil {
		return nil, 0, err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return nil, 0, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}
	return cond, closeTok.Loc.End, nil
}

// parseOrderByItems parses `ORDER BY sortItem (, sortItem)*` (ORDER is current).
func (p *Parser) parseOrderByItems() ([]SortItem, error) {
	p.advance() // consume ORDER
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	first, err := p.parseSortItem()
	if err != nil {
		return nil, err
	}
	items := []SortItem{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseSortItem()
		if err != nil {
			return nil, err
		}
		items = append(items, next)
	}
	return items, nil
}

// parseSortItem parses one `expression [ASC|DESC] [NULLS FIRST|LAST]`
// (the sortItem rule).
func (p *Parser) parseSortItem() (SortItem, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return SortItem{}, err
	}
	item := SortItem{Expr: expr, Loc: ast.Loc{Start: expr.Span().Start, End: expr.Span().End}}
	if p.cur.Kind == kwASC || p.cur.Kind == kwDESC {
		tok := p.advance()
		if tok.Kind == kwASC {
			item.Ordering = "ASC"
		} else {
			item.Ordering = "DESC"
		}
		item.Loc.End = tok.Loc.End
	}
	if p.cur.Kind == kwNULLS {
		p.advance() // consume NULLS
		switch p.cur.Kind {
		case kwFIRST:
			tok := p.advance()
			item.NullOrder = "FIRST"
			item.Loc.End = tok.Loc.End
		case kwLAST:
			tok := p.advance()
			item.NullOrder = "LAST"
			item.Loc.End = tok.Loc.End
		default:
			return SortItem{}, p.exprErrorAt("expected FIRST or LAST after NULLS")
		}
	}
	return item, nil
}

// ---------------------------------------------------------------------------
// OVER / window specification / window frame
// ---------------------------------------------------------------------------

// parseOver parses `OVER (windowName | ( windowSpecification ))` (the over rule).
// It returns either a *WindowSpec (the parenthesized form) or a windowName
// identifier (the named form); exactly one is non-nil. OVER is the current token.
func (p *Parser) parseOver() (*WindowSpec, *ast.Identifier, error) {
	overTok := p.advance() // consume OVER
	if p.cur.Kind == int('(') {
		spec, err := p.parseWindowSpecification(overTok.Loc.Start)
		if err != nil {
			return nil, nil, err
		}
		return spec, nil, nil
	}
	// Named window: OVER windowName.
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, nil, err
	}
	return nil, name, nil
}

// parseWindowSpecification parses `( [existingName] [PARTITION BY expr…]
// [ORDER BY sortItem…] [windowFrame] )` (windowSpecification). The opening '('
// is the current token; startOffset is the OVER keyword's offset so the span
// covers `OVER ( … )`.
func (p *Parser) parseWindowSpecification(startOffset int) (*WindowSpec, error) {
	p.advance() // consume '('
	spec := &WindowSpec{Loc: ast.Loc{Start: startOffset}}

	// Optional existing window name: a leading identifier that is NOT the start
	// of PARTITION/ORDER/a frame keyword or ')'.
	if isIdentifierStart(p.cur.Kind) && p.cur.Kind != kwPARTITION && p.cur.Kind != kwORDER &&
		!isFrameTypeKeyword(p.cur.Kind) {
		spec.ExistingName = identFromToken(p.advance())
	}

	if p.cur.Kind == kwPARTITION {
		p.advance() // consume PARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		first, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		spec.PartitionBy = append(spec.PartitionBy, first)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			spec.PartitionBy = append(spec.PartitionBy, next)
		}
	}

	if p.cur.Kind == kwORDER {
		items, err := p.parseOrderByItems()
		if err != nil {
			return nil, err
		}
		spec.OrderBy = items
	}

	if isFrameTypeKeyword(p.cur.Kind) {
		frame, err := p.parseWindowFrame()
		if err != nil {
			return nil, err
		}
		spec.Frame = frame
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	spec.Loc.End = closeTok.Loc.End
	return spec, nil
}

// isFrameTypeKeyword reports whether kind is a window-frame type keyword
// (ROWS / RANGE / GROUPS).
func isFrameTypeKeyword(kind TokenKind) bool {
	return kind == kwROWS || kind == kwRANGE || kind == kwGROUPS
}

// parseWindowFrame parses the common frame `frameType (BETWEEN start AND end |
// start)` (frameExtent + frameBound). The pattern-recognition frame additions
// (MEASURES/PATTERN/DEFINE/…) are deferred to parser-match-recognize (B2); this
// node implements only the standard frame.
func (p *Parser) parseWindowFrame() (*WindowFrame, error) {
	typeTok := p.advance() // ROWS / RANGE / GROUPS
	frame := &WindowFrame{
		FrameType: typeTok.Str,
		Loc:       ast.Loc{Start: typeTok.Loc.Start},
	}
	if p.cur.Kind == kwBETWEEN {
		p.advance() // consume BETWEEN
		start, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}
		end, err := p.parseFrameBound()
		if err != nil {
			return nil, err
		}
		frame.Start = start
		frame.End = &end
		frame.Loc.End = end.Loc.End
		return frame, nil
	}
	start, err := p.parseFrameBound()
	if err != nil {
		return nil, err
	}
	frame.Start = start
	frame.Loc.End = start.Loc.End
	return frame, nil
}

// parseFrameBound parses one `UNBOUNDED PRECEDING | UNBOUNDED FOLLOWING |
// CURRENT ROW | expr (PRECEDING|FOLLOWING)` (frameBound).
func (p *Parser) parseFrameBound() (WindowBound, error) {
	switch p.cur.Kind {
	case kwUNBOUNDED:
		unbTok := p.advance() // consume UNBOUNDED
		switch p.cur.Kind {
		case kwPRECEDING:
			tok := p.advance()
			return WindowBound{Kind: BoundUnboundedPreceding, Loc: ast.Loc{Start: unbTok.Loc.Start, End: tok.Loc.End}}, nil
		case kwFOLLOWING:
			tok := p.advance()
			return WindowBound{Kind: BoundUnboundedFollowing, Loc: ast.Loc{Start: unbTok.Loc.Start, End: tok.Loc.End}}, nil
		default:
			return WindowBound{}, p.exprErrorAt("expected PRECEDING or FOLLOWING after UNBOUNDED")
		}
	case kwCURRENT:
		curTok := p.advance() // consume CURRENT
		rowTok, err := p.expect(kwROW)
		if err != nil {
			return WindowBound{}, err
		}
		return WindowBound{Kind: BoundCurrentRow, Loc: ast.Loc{Start: curTok.Loc.Start, End: rowTok.Loc.End}}, nil
	default:
		expr, err := p.parseExpr()
		if err != nil {
			return WindowBound{}, err
		}
		switch p.cur.Kind {
		case kwPRECEDING:
			tok := p.advance()
			return WindowBound{Kind: BoundPreceding, Value: expr, Loc: ast.Loc{Start: expr.Span().Start, End: tok.Loc.End}}, nil
		case kwFOLLOWING:
			tok := p.advance()
			return WindowBound{Kind: BoundFollowing, Value: expr, Loc: ast.Loc{Start: expr.Span().Start, End: tok.Loc.End}}, nil
		default:
			return WindowBound{}, p.exprErrorAt("expected PRECEDING or FOLLOWING in frame bound")
		}
	}
}

// ---------------------------------------------------------------------------
// keyword-led primaryExpression atoms
// ---------------------------------------------------------------------------

// parseCaseExpr parses a simple `CASE operand WHEN … END` or searched
// `CASE WHEN … END` (simpleCase / searchedCase). At least one WHEN is required.
func (p *Parser) parseCaseExpr() (Expr, error) {
	caseTok := p.advance() // consume CASE
	ce := &CaseExpr{Loc: ast.Loc{Start: caseTok.Loc.Start}}

	// Simple CASE has an operand expression before the first WHEN.
	if p.cur.Kind != kwWHEN {
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Operand = operand
	}

	for p.cur.Kind == kwWHEN {
		whenTok := p.advance() // consume WHEN
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		result, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Whens = append(ce.Whens, WhenClause{
			Cond:   cond,
			Result: result,
			Loc:    ast.Loc{Start: whenTok.Loc.Start, End: result.Span().End},
		})
	}
	if len(ce.Whens) == 0 {
		return nil, p.exprErrorAt("expected WHEN in CASE expression")
	}

	if p.cur.Kind == kwELSE {
		p.advance() // consume ELSE
		elseExpr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		ce.Else = elseExpr
	}

	endTok, err := p.expect(kwEND)
	if err != nil {
		return nil, err
	}
	ce.Loc.End = endTok.Loc.End
	return ce, nil
}

// parseCast parses `CAST(expr AS type)` or `TRY_CAST(expr AS type)` (cast).
func (p *Parser) parseCast() (Expr, error) {
	castTok := p.advance() // CAST or TRY_CAST
	try := castTok.Kind == kwTRY_CAST
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	typ, err := p.parseType()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &CastExpr{
		Try:  try,
		Expr: expr,
		Type: typ,
		Loc:  ast.Loc{Start: castTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseExists parses `EXISTS ( query )` (exists). The subquery is a raw-text
// placeholder (B1). EXISTS strictly requires a query — `EXISTS ()`, `EXISTS (1)`,
// `EXISTS (1 + 2)` are SYNTAX_ERRORs in Trino 481 (EXISTS has no expression
// form), so the '(' must be followed by a query-starting token (SELECT / WITH /
// TABLE / VALUES); otherwise it is rejected here.
func (p *Parser) parseExists() (Expr, error) {
	existsTok := p.advance() // consume EXISTS
	openTok, err := p.expect(int('('))
	if err != nil {
		return nil, err
	}
	if !p.startsQuery() {
		return nil, p.exprErrorAt("EXISTS requires a subquery")
	}
	subq, err := p.parseSubqueryPlaceholder(openTok.Loc.Start, SubqueryExists)
	if err != nil {
		return nil, err
	}
	subq.Loc.Start = existsTok.Loc.Start
	return subq, nil
}

// parseExtract parses `EXTRACT( field FROM source )` (extract). field is an
// identifier (YEAR, MONTH, TIMEZONE_HOUR, …) kept as source text.
func (p *Parser) parseExtract() (Expr, error) {
	extractTok := p.advance() // consume EXTRACT
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	field, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	source, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ExtractExpr{
		Field:  field.String(),
		Source: source,
		Loc:    ast.Loc{Start: extractTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseSubstring parses `SUBSTRING( source FROM start (FOR length)? )` (the
// SQL-standard substring spelling). SUBSTRING is a NON-RESERVED keyword, so the
// comma-argument form `substring(s, a, b)` and the bare column reference
// `substring` are also valid (oracle-confirmed) and are handled by
// parseIdentifierOrFunction, which tries this special parser first and falls back
// to the ordinary identifier/function reading when the FROM form does not match.
func (p *Parser) parseSubstring() (Expr, error) {
	subTok := p.advance() // consume SUBSTRING
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	source, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	from, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	se := &SubstringExpr{Source: source, From: from, Loc: ast.Loc{Start: subTok.Loc.Start}}
	if p.cur.Kind == kwFOR {
		p.advance() // consume FOR
		length, err := p.parseValueExpr()
		if err != nil {
			return nil, err
		}
		se.For = length
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	se.Loc.End = closeTok.Loc.End
	return se, nil
}

// parseTrim parses the two TRIM spellings: `TRIM( [spec] [char] FROM source )`
// and `TRIM( source , char )` (trim). The first reading is tried, and on failure
// the parser rewinds and tries the comma form.
//
// Oracle-confirmed divergence (D-TRIM): a bare `TRIM(FROM source)` — FROM with
// neither a trims-specification nor a trim-character before it — is a
// SYNTAX_ERROR in Trino 481, even though the legacy grammar's
// `(trimsSpecification? trimChar? FROM)?` permits it. FROM is accepted only when
// preceded by a spec and/or a char (TRIM(BOTH FROM x), TRIM('y' FROM x),
// TRIM(BOTH 'y' FROM x)); TRIM(x) and TRIM(x, y) remain valid. See tryTrimFromForm.
func (p *Parser) parseTrim() (Expr, error) {
	trimTok := p.advance() // consume TRIM
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Reading 1 — `[spec] [char] FROM source`.
	cp := p.checkpoint()
	if te, ok := p.tryTrimFromForm(trimTok.Loc.Start); ok {
		return te, nil
	}
	p.restore(cp)

	// Reading 2 — `source , char`.
	source, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(',')); err != nil {
		return nil, err
	}
	char, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &TrimExpr{
		Char:   char,
		Source: source,
		Loc:    ast.Loc{Start: trimTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// tryTrimFromForm attempts the `[spec] [char] FROM source` TRIM reading after
// the opening '(' has been consumed, returning ok=false (cursor left for the
// caller's restore) when it does not fit. A bare `TRIM(source)` with no FROM and
// no comma also matches here (spec/char absent, the whole thing is the source)
// — Trino accepts `TRIM(x)`.
func (p *Parser) tryTrimFromForm(start int) (Expr, bool) {
	te := &TrimExpr{Loc: ast.Loc{Start: start}}

	// Optional trims specification.
	switch p.cur.Kind {
	case kwLEADING:
		p.advance()
		te.Spec = "LEADING"
	case kwTRAILING:
		p.advance()
		te.Spec = "TRAILING"
	case kwBOTH:
		p.advance()
		te.Spec = "BOTH"
	}

	// No spec present. A bare `FROM source` (no spec, no char) is a SYNTAX_ERROR
	// in Trino 481 — `TRIM(FROM x)` is rejected even though the legacy grammar's
	// `(trimsSpecification? trimChar? FROM)?` permits it (oracle-confirmed
	// divergence). So FROM here is valid only after a char: parse a value
	// expression first and require it to be followed by FROM or ')'.
	if te.Spec == "" {
		if p.cur.Kind == kwFROM {
			// Bare FROM with nothing before it — reject this reading (and, since
			// the comma form also cannot start with FROM, the whole TRIM is an
			// error, which the caller surfaces).
			return nil, false
		}
		first, err := p.parseValueExpr()
		if err != nil {
			return nil, false
		}
		switch p.cur.Kind {
		case int(')'):
			// TRIM(source)
			closeTok := p.advance()
			te.Source = first
			te.Loc.End = closeTok.Loc.End
			return te, true
		case kwFROM:
			// TRIM(char FROM source)
			p.advance() // consume FROM
			source, ok := p.tryValueExprThenClose()
			if !ok {
				return nil, false
			}
			te.Char = first
			te.Source = source.expr
			te.Loc.End = source.end
			return te, true
		default:
			// A comma (the second TRIM form) or anything else — not this reading.
			return nil, false
		}
	}

	// With a spec present, an optional char then a required FROM, then source.
	if p.cur.Kind != kwFROM {
		char, err := p.parseValueExpr()
		if err != nil {
			return nil, false
		}
		te.Char = char
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, false
	}
	source, ok := p.tryValueExprThenClose()
	if !ok {
		return nil, false
	}
	te.Source = source.expr
	te.Loc.End = source.end
	return te, true
}

// valueExprAndEnd bundles a parsed value expression with the closing-')' end
// offset, for the trim helper.
type valueExprAndEnd struct {
	expr Expr
	end  int
}

// tryValueExprThenClose parses a value expression followed by a required ')',
// returning ok=false (no error surfaced) when either step fails so the trim
// reading can be abandoned.
func (p *Parser) tryValueExprThenClose() (valueExprAndEnd, bool) {
	expr, err := p.parseValueExpr()
	if err != nil {
		return valueExprAndEnd{}, false
	}
	if p.cur.Kind != int(')') {
		return valueExprAndEnd{}, false
	}
	closeTok := p.advance()
	return valueExprAndEnd{expr: expr, end: closeTok.Loc.End}, true
}

// parseNormalize parses `NORMALIZE( source (, form)? )` (normalize). form is one
// of NFD/NFC/NFKD/NFKC.
func (p *Parser) parseNormalize() (Expr, error) {
	normTok := p.advance() // consume NORMALIZE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	source, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	ne := &NormalizeExpr{Source: source, Loc: ast.Loc{Start: normTok.Loc.Start}}
	if p.cur.Kind == int(',') {
		p.advance() // consume ','
		form, ok := normalFormText(p.cur.Kind)
		if !ok {
			return nil, p.exprErrorAt("expected normal form (NFD/NFC/NFKD/NFKC)")
		}
		p.advance()
		ne.Form = form
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	ne.Loc.End = closeTok.Loc.End
	return ne, nil
}

// normalFormText maps a normal-form keyword kind to its spelling.
func normalFormText(kind TokenKind) (string, bool) {
	switch kind {
	case kwNFD:
		return "NFD", true
	case kwNFC:
		return "NFC", true
	case kwNFKD:
		return "NFKD", true
	case kwNFKC:
		return "NFKC", true
	default:
		return "", false
	}
}

// parsePosition parses `POSITION( needle IN haystack )` (position). POSITION is
// a NON-RESERVED keyword, so the comma-argument form `position('a', 'abc')` and
// the bare column reference `position` are also valid (oracle-confirmed); like
// SUBSTRING, parseIdentifierOrFunction tries this special parser first and falls
// back to the ordinary reading when the IN form does not match.
func (p *Parser) parsePosition() (Expr, error) {
	posTok := p.advance() // consume POSITION
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	needle, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	haystack, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &PositionExpr{
		Needle:   needle,
		Haystack: haystack,
		Loc:      ast.Loc{Start: posTok.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseGrouping parses `GROUPING( (qualifiedName (, qualifiedName)*)? )`
// (groupingOperation). The argument list may be empty.
func (p *Parser) parseGrouping() (Expr, error) {
	groupingTok := p.advance() // consume GROUPING
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	ge := &GroupingExpr{Loc: ast.Loc{Start: groupingTok.Loc.Start}}
	if p.cur.Kind != int(')') {
		first, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		ge.Args = append(ge.Args, first)
		for p.cur.Kind == int(',') {
			p.advance() // consume ','
			next, err := p.parseQualifiedName()
			if err != nil {
				return nil, err
			}
			ge.Args = append(ge.Args, next)
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	ge.Loc.End = closeTok.Loc.End
	return ge, nil
}

// parseListagg parses `LISTAGG( [DISTINCT] arg (, sep)? (ON OVERFLOW
// listAggOverflowBehavior)? ) WITHIN GROUP ( ORDER BY sortItem (, sortItem)* )
// filter?` (listagg). The mandatory WITHIN GROUP (ORDER BY …) distinguishes
// LISTAGG from an ordinary function call.
func (p *Parser) parseListagg() (Expr, error) {
	listaggTok := p.advance() // consume LISTAGG
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	la := &ListaggExpr{Loc: ast.Loc{Start: listaggTok.Loc.Start}}

	if p.cur.Kind == kwDISTINCT {
		p.advance()
		la.Distinct = true
	} else if p.cur.Kind == kwALL {
		p.advance() // ALL is the default
	}

	arg, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	la.Arg = arg

	if p.cur.Kind == int(',') {
		p.advance() // consume ','
		sep, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		la.Separator = sep
	}

	if p.cur.Kind == kwON {
		p.advance() // consume ON
		if _, err := p.expect(kwOVERFLOW); err != nil {
			return nil, err
		}
		switch p.cur.Kind {
		case kwERROR:
			p.advance()
			la.OnOverflow = "ERROR"
		case kwTRUNCATE:
			p.advance()
			la.OnOverflow = "TRUNCATE"
			// Optional truncation filler string.
			if p.cur.Kind == tokString || p.cur.Kind == tokUnicodeString {
				p.advance()
			}
			// listaggCountIndication: (WITH | WITHOUT) COUNT.
			if p.cur.Kind == kwWITH || p.cur.Kind == kwWITHOUT {
				p.advance()
				if _, err := p.expect(kwCOUNT); err != nil {
					return nil, err
				}
			}
		default:
			return nil, p.exprErrorAt("expected ERROR or TRUNCATE after ON OVERFLOW")
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	// Mandatory WITHIN GROUP ( ORDER BY … ).
	if _, err := p.expect(kwWITHIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	items, err := p.parseOrderByItems()
	if err != nil {
		return nil, err
	}
	la.WithinGroupBy = items
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	la.Loc.End = closeTok.Loc.End

	// Optional FILTER.
	if p.cur.Kind == kwFILTER {
		filter, end, err := p.parseFilter()
		if err != nil {
			return nil, err
		}
		la.Filter = filter
		la.Loc.End = end
	}
	return la, nil
}

// parseSpecialDateTimeFunction parses the no-parenthesis / precision-only special
// functions: CURRENT_DATE, CURRENT_TIME[(p)], CURRENT_TIMESTAMP[(p)],
// LOCALTIME[(p)], LOCALTIMESTAMP[(p)], CURRENT_USER, CURRENT_CATALOG,
// CURRENT_SCHEMA, CURRENT_PATH (specialDateTimeFunction + the current* set).
// Only the four time functions accept an optional `(INTEGER)` precision.
func (p *Parser) parseSpecialDateTimeFunction() (Expr, error) {
	tok := p.advance()
	sf := &SpecialFuncExpr{Name: tok.Str, Precision: -1, Loc: tok.Loc}

	if allowsPrecision(tok.Kind) && p.cur.Kind == int('(') {
		p.advance() // consume '('
		if p.cur.Kind != tokInteger {
			return nil, p.exprErrorAt("expected integer precision")
		}
		precTok := p.advance()
		sf.Precision = int(precTok.Ival)
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		sf.Loc.End = closeTok.Loc.End
	}
	return sf, nil
}

// allowsPrecision reports whether a special date-time keyword accepts an optional
// `(precision)` argument (CURRENT_TIME / CURRENT_TIMESTAMP / LOCALTIME /
// LOCALTIMESTAMP). CURRENT_DATE and the current-user/catalog/schema/path set do
// not.
func allowsPrecision(kind TokenKind) bool {
	switch kind {
	case kwCURRENT_TIME, kwCURRENT_TIMESTAMP, kwLOCALTIME, kwLOCALTIMESTAMP:
		return true
	default:
		return false
	}
}

// parseDoublePrecisionConstructor parses `DOUBLE PRECISION 'string'`
// (the DOUBLE PRECISION typeConstructor). The caller confirmed DOUBLE then
// PRECISION; a string literal must follow.
func (p *Parser) parseDoublePrecisionConstructor() (Expr, error) {
	doubleTok := p.advance() // DOUBLE
	p.advance()              // PRECISION
	if p.cur.Kind != tokString && p.cur.Kind != tokUnicodeString {
		return nil, p.exprErrorAt("expected string literal after DOUBLE PRECISION")
	}
	strTok := p.advance()
	return &TypeConstructor{
		Name:  "DOUBLE PRECISION",
		Value: strTok.Str,
		Loc:   ast.Loc{Start: doubleTok.Loc.Start, End: strTok.Loc.End},
	}, nil
}

// parseJSONFunction (the SQL/JSON function dispatcher) lives in json.go (the
// expr-json DAG node). parsePrimaryAtom's JSON_* dispatch calls it from there.
