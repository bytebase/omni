package parser

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/trino/ast"
)

// This file is the `types` DAG node: it implements Trino's `type` grammar rule
// (and its helpers `rowField`, `typeParameter`, `intervalField`) as a
// hand-written recursive-descent parser over the token stream, producing a
// DataType value tree.
//
// Trino's legacy ANTLR `type` rule (TrinoParser.g4):
//
//	type
//	    : ROW_ LPAREN_ rowField (COMMA_ rowField)* RPAREN_                  # rowType
//	    | INTERVAL_ from=intervalField (TO_ to=intervalField)?             # intervalType
//	    | base_=TIMESTAMP_ (LPAREN_ precision=typeParameter RPAREN_)? (WITHOUT_ TIME_ ZONE_)?  # dateTimeType
//	    | base_=TIMESTAMP_ (LPAREN_ precision=typeParameter RPAREN_)? WITH_ TIME_ ZONE_        # dateTimeType
//	    | base_=TIME_      (LPAREN_ precision=typeParameter RPAREN_)? (WITHOUT_ TIME_ ZONE_)?  # dateTimeType
//	    | base_=TIME_      (LPAREN_ precision=typeParameter RPAREN_)? WITH_ TIME_ ZONE_        # dateTimeType
//	    | DOUBLE_ PRECISION_                                               # doublePrecisionType
//	    | ARRAY_ LT_ type GT_                                              # legacyArrayType
//	    | MAP_ LT_ keyType=type COMMA_ valueType=type GT_                 # legacyMapType
//	    | type ARRAY_ (LSQUARE_ INTEGER_VALUE_ RSQUARE_)?                  # arrayType
//	    | identifier (LPAREN_ typeParameter (COMMA_ typeParameter)* RPAREN_)?  # genericType
//	    ;
//	rowField      : type | identifier type ;
//	typeParameter : INTEGER_VALUE_ | type ;
//	intervalField : YEAR_ | MONTH_ | DAY_ | HOUR_ | MINUTE_ | SECOND_ ;
//
// The implementation is adjudicated against the live Trino 481 oracle, not the
// literal legacy grammar. Two oracle-confirmed divergences from a naive reading
// of the legacy rule are baked in (see the migration divergence ledger):
//
//	D1 (interval qualifier set). The legacy `from (TO to)?` allows ANY field TO
//	   ANY field (e.g. SECOND TO YEAR). Trino 481 restricts `from TO to` to the
//	   two SQL-standard families with from strictly coarser than to:
//	   year-month {YEAR>MONTH} and day-time {DAY>HOUR>MINUTE>SECOND}; a TO
//	   crossing families (YEAR TO DAY) or reversed (MONTH TO YEAR) is a
//	   SYNTAX_ERROR. The single-field form `INTERVAL <field>` is accepted for
//	   every field (it fails later as NOT_SUPPORTED for year-month single
//	   fields, which is semantic, not syntactic).
//
//	D2 (INTERVAL not followed by a field is a generic type). Legacy requires a
//	   field after INTERVAL. Trino 481 instead falls back to the genericType
//	   path when INTERVAL is followed by `(`, `)`, `ARRAY`, or anything that is
//	   not an intervalField keyword: `INTERVAL`, `INTERVAL(5)`, `INTERVAL ARRAY`
//	   all parse (as the unknown type INTERVAL / INTERVAL(5) / ARRAY(INTERVAL)).
//	   So INTERVAL only enters the intervalType branch when the next token is one
//	   of YEAR/MONTH/DAY/HOUR/MINUTE/SECOND.
//
// Recursive-descent vs the ANTLR alternative ordering: the left-recursive
// `type ARRAY` postfix is handled iteratively (parse a base type, then consume
// zero or more trailing `ARRAY [n]` suffixes). The keyword-led alternatives
// (ROW, INTERVAL, TIMESTAMP/TIME, DOUBLE PRECISION, ARRAY<>, MAP<>) are
// dispatched on the leading keyword with lookahead; everything else, including
// the type-name keywords used WITHOUT their special syntax (e.g. `ARRAY` alone,
// `MAP` alone, `DOUBLE` alone, `ROW` alone), falls through to genericType,
// matching Trino which treats them as ordinary type-name identifiers there.

// TypeKind classifies the top-level shape of a parsed DataType. It mirrors the
// labeled alternatives of the legacy `type` rule so downstream consumers
// (expressions, DDL, deparse) can switch on the form without re-inspecting the
// token text.
type TypeKind int

const (
	// TypeGeneric is a named type with optional parenthesized parameters:
	// `identifier (LPAREN typeParameter (COMMA typeParameter)* RPAREN)?`.
	// Covers the overwhelming majority of Trino types — BIGINT, VARCHAR(n),
	// DECIMAL(p,s), JSON, UUID, IPADDRESS, HYPERLOGLOG, QDIGEST(t), VARIANT,
	// ARRAY(t)/MAP(k,v) in their parenthesized spelling, DOUBLE, etc. — because
	// Trino's type names are plain identifiers, not reserved keywords.
	TypeGeneric TypeKind = iota
	// TypeRow is `ROW ( rowField (, rowField)* )`.
	TypeRow
	// TypeInterval is `INTERVAL from (TO to)?`.
	TypeInterval
	// TypeDateTime is a TIMESTAMP or TIME type with optional precision and an
	// optional WITH/WITHOUT TIME ZONE clause.
	TypeDateTime
	// TypeArray is the legacy angle-bracket array `ARRAY < type >` OR the
	// postfix `type ARRAY [n]`. The element is in ElementType.
	TypeArray
	// TypeMap is the legacy angle-bracket map `MAP < keyType , valueType >`.
	// KeyType and ValueType hold the components.
	TypeMap
)

// String returns a human-readable tag for the kind (diagnostics only).
func (k TypeKind) String() string {
	switch k {
	case TypeGeneric:
		return "Generic"
	case TypeRow:
		return "Row"
	case TypeInterval:
		return "Interval"
	case TypeDateTime:
		return "DateTime"
	case TypeArray:
		return "Array"
	case TypeMap:
		return "Map"
	default:
		return "Unknown"
	}
}

// IntervalField enumerates the single-unit interval qualifiers
// (intervalField rule). Ordering is the standard coarse→fine sequence and is
// significant: a `from TO to` range is valid only when both fields share a
// family and from precedes to (see ValidIntervalRange).
type IntervalField int

const (
	IntervalYear IntervalField = iota
	IntervalMonth
	IntervalDay
	IntervalHour
	IntervalMinute
	IntervalSecond
)

// String returns the SQL keyword for the interval field.
func (f IntervalField) String() string {
	switch f {
	case IntervalYear:
		return "YEAR"
	case IntervalMonth:
		return "MONTH"
	case IntervalDay:
		return "DAY"
	case IntervalHour:
		return "HOUR"
	case IntervalMinute:
		return "MINUTE"
	case IntervalSecond:
		return "SECOND"
	default:
		return "?"
	}
}

// isYearMonth reports whether the field belongs to the year-month interval
// family ({YEAR, MONTH}); the remaining fields form the day-time family.
func (f IntervalField) isYearMonth() bool {
	return f == IntervalYear || f == IntervalMonth
}

// ValidIntervalRange reports whether `INTERVAL from TO to` is a syntactically
// valid range in Trino 481: both fields must be in the same family
// (year-month or day-time) and from must be strictly coarser than to
// (smaller enum value). This is the D1 divergence — the legacy grammar permits
// every combination; Trino accepts only YEAR TO MONTH and the day-time chain
// DAY/HOUR/MINUTE → a strictly finer day-time field.
func ValidIntervalRange(from, to IntervalField) bool {
	if from.isYearMonth() != to.isYearMonth() {
		return false
	}
	return from < to
}

// TypeParam is one entry of a genericType parameter list. Exactly one of
// IsInt/Type is meaningful: an INTEGER_VALUE_ parameter (e.g. the 10 in
// DECIMAL(10,2)) sets IsInt, IntVal, and IntText; a nested type parameter
// (e.g. the BIGINT in QDIGEST(BIGINT)) sets Type. typeParameter :=
// INTEGER_VALUE_ | type.
//
// IntText preserves the integer's exact source spelling so round-trip
// rendering is faithful even when the value overflows int64 (the lexer leaves
// IntVal at 0 for an out-of-range INTEGER_VALUE_ but keeps the literal text).
type TypeParam struct {
	IsInt   bool
	IntVal  int64
	IntText string
	Type    *DataType
	Loc     ast.Loc
}

// RowField is one field of a ROW type. Name is the optional field-name
// identifier (nil for an unnamed field); Type is the field type. rowField :=
// type | identifier type. Integer entries (ROW(1)) are not RowFields — they are
// parsed under the genericType "ROW" reading as TypeParams instead.
type RowField struct {
	Name *ast.Identifier
	Type *DataType
	Loc  ast.Loc
}

// DataType is a parsed Trino data type. It is a parser-package node (not an
// ast.Node — the Trino ast tag set is closed to the ast-core node) carrying its
// own source span; later DAG nodes (expressions, DDL) embed *DataType in their
// ast results and deparse renders it via String.
//
// The active fields depend on Kind:
//
//	TypeGeneric  — Name, Params
//	TypeRow      — Fields
//	TypeInterval — IntervalFrom, IntervalTo (IntervalTo==nil for single-field)
//	TypeDateTime — Name (TIMESTAMP|TIME), Precision, WithTimeZone, hasTimeZone
//	TypeArray    — ElementType (legacy ARRAY<t> or postfix `t ARRAY [n]`)
//	TypeMap      — KeyType, ValueType
type DataType struct {
	Kind TypeKind

	// Name is the source type name. For TypeGeneric it is the identifier
	// (source-faithful, case preserved). For TypeDateTime it is "TIMESTAMP" or
	// "TIME". For the legacy array/map and row/interval forms it is the leading
	// keyword ("ARRAY"/"MAP"/"ROW"/"INTERVAL") for round-tripping.
	Name string

	// Params holds a genericType parameter list (TypeGeneric only); nil when
	// the type had no parenthesized parameters.
	Params []TypeParam

	// Fields holds the ROW field list (TypeRow only).
	Fields []RowField

	// ElementType is the element of an array (TypeArray): the inner type of a
	// legacy ARRAY<t> or the operand of a postfix `t ARRAY`.
	ElementType *DataType

	// KeyType / ValueType are the components of a legacy MAP<k,v> (TypeMap).
	KeyType   *DataType
	ValueType *DataType

	// IntervalFrom / IntervalTo describe an INTERVAL type (TypeInterval).
	// IntervalTo is nil for the single-field form `INTERVAL <field>`.
	IntervalFrom IntervalField
	IntervalTo   *IntervalField

	// Precision is the optional precision typeParameter of a TIMESTAMP/TIME
	// (TypeDateTime); nil when absent. Trino allows precision to be an integer
	// or (oddly) a nested type; both are captured, with Type set for the latter.
	Precision *TypeParam

	// WithTimeZone is true for a `WITH TIME ZONE` date-time type, false for
	// `WITHOUT TIME ZONE` or no zone clause. HasTimeZoneClause distinguishes an
	// explicit clause from the default; only meaningful for TypeDateTime.
	WithTimeZone      bool
	HasTimeZoneClause bool

	// ArrayDim is the optional bracketed dimension of a postfix array
	// `t ARRAY[n]`; -1 when absent. (Trino accepts the bracket syntactically
	// but rejects it during analysis; we preserve it for fidelity.)
	ArrayDim int

	Loc ast.Loc
}

// String renders the DataType back to Trino source syntax. It is the
// round-trip rendering downstream deparse relies on; it is not required to be
// byte-identical to the input (case of unquoted names is preserved, but
// inter-token spacing is normalized).
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
	case TypeGeneric:
		b.WriteString(t.Name)
		if t.Params != nil {
			b.WriteByte('(')
			for i, p := range t.Params {
				if i > 0 {
					b.WriteString(", ")
				}
				p.writeTo(b)
			}
			b.WriteByte(')')
		}
	case TypeRow:
		b.WriteString("ROW(")
		for i, f := range t.Fields {
			if i > 0 {
				b.WriteString(", ")
			}
			f.writeTo(b)
		}
		b.WriteByte(')')
	case TypeInterval:
		b.WriteString("INTERVAL ")
		b.WriteString(t.IntervalFrom.String())
		if t.IntervalTo != nil {
			b.WriteString(" TO ")
			b.WriteString(t.IntervalTo.String())
		}
	case TypeDateTime:
		b.WriteString(t.Name)
		if t.Precision != nil {
			b.WriteByte('(')
			t.Precision.writeTo(b)
			b.WriteByte(')')
		}
		if t.HasTimeZoneClause {
			if t.WithTimeZone {
				b.WriteString(" WITH TIME ZONE")
			} else {
				b.WriteString(" WITHOUT TIME ZONE")
			}
		}
	case TypeArray:
		// Render in the legacy angle-bracket form for an unambiguous round-trip
		// (postfix `t ARRAY` and `ARRAY<t>` denote the same type).
		b.WriteString("ARRAY<")
		t.ElementType.writeTo(b)
		b.WriteByte('>')
	case TypeMap:
		b.WriteString("MAP<")
		t.KeyType.writeTo(b)
		b.WriteString(", ")
		t.ValueType.writeTo(b)
		b.WriteByte('>')
	}
}

func (p TypeParam) writeTo(b *strings.Builder) {
	if p.IsInt {
		// Prefer the exact source spelling (faithful even on int64 overflow);
		// fall back to the parsed value if no text was captured.
		if p.IntText != "" {
			b.WriteString(p.IntText)
		} else {
			b.WriteString(strconv.FormatInt(p.IntVal, 10))
		}
		return
	}
	if p.Type != nil {
		p.Type.writeTo(b)
	}
}

func (f RowField) writeTo(b *strings.Builder) {
	if f.Name != nil {
		b.WriteString(f.Name.String())
		b.WriteByte(' ')
	}
	if f.Type != nil {
		f.Type.writeTo(b)
	}
}

// ---------------------------------------------------------------------------
// parseType — Trino `type` rule
// ---------------------------------------------------------------------------

// parseType parses a Trino type (the `type` grammar rule) and returns it as a
// *DataType. It is the single entry point used by every type position — CAST /
// TRY_CAST targets, column definitions, RETURNING clauses, routine parameter
// and RETURNS clauses, variable declarations. Returns a *ParseError if the
// current token does not begin a valid type.
//
// The postfix `type ARRAY [n]` is left-recursive in the grammar; it is parsed
// here by consuming a base type and then looping over trailing ARRAY suffixes.
func (p *Parser) parseType() (*DataType, error) {
	base, err := p.parseBaseType()
	if err != nil {
		return nil, err
	}
	return p.parseArraySuffixes(base)
}

// parseArraySuffixes consumes zero or more postfix `ARRAY (LSQUARE INTEGER
// RSQUARE)?` suffixes on an already-parsed base type, wrapping it in a
// TypeArray for each. `t ARRAY ARRAY` nests twice; `t ARRAY[5]` records the
// bracketed dimension.
//
// The bracket part is optional, but when a '[' is present the grammar requires
// EXACTLY `[ INTEGER ]` — Trino 481 rejects `t ARRAY[]`, `t ARRAY[abc]`,
// `t ARRAY[5` (no close), and `t ARRAY[5 6]` with a SYNTAX_ERROR, so a malformed
// bracket here returns a *ParseError rather than being silently tolerated.
// (Trino accepts `t ARRAY[n]` at parse time for any integer n; the bracketed
// dimension fails later during analysis, which is semantic, not syntactic.)
func (p *Parser) parseArraySuffixes(base *DataType) (*DataType, error) {
	for p.cur.Kind == kwARRAY {
		arrTok := p.advance() // consume ARRAY
		end := arrTok.Loc.End
		dim := -1
		if p.cur.Kind == int('[') {
			p.advance() // consume '['
			if p.cur.Kind != tokInteger {
				return nil, p.typeErrorAt("expected integer array dimension after '['")
			}
			dim = int(p.cur.Ival)
			p.advance()
			closeTok, err := p.expect(int(']'))
			if err != nil {
				return nil, err
			}
			end = closeTok.Loc.End
		}
		base = &DataType{
			Kind:        TypeArray,
			Name:        "ARRAY",
			ElementType: base,
			ArrayDim:    dim,
			Loc:         ast.Loc{Start: base.Loc.Start, End: end},
		}
	}
	return base, nil
}

// parseBaseType parses one `type` alternative WITHOUT the postfix-array suffix
// (handled by parseArraySuffixes). It dispatches on the leading keyword:
//
//   - ROW followed by '(' → rowType.
//   - INTERVAL followed by an intervalField keyword → intervalType (D2: any
//     other follow-token, or none, falls through to genericType "INTERVAL").
//   - TIMESTAMP / TIME → dateTimeType.
//   - DOUBLE followed by PRECISION → doublePrecisionType ("DOUBLE PRECISION").
//   - ARRAY followed by '<' → legacy ARRAY<type>.
//   - MAP followed by '<' → legacy MAP<key,value>.
//   - anything else (including ROW/ARRAY/MAP/DOUBLE/INTERVAL NOT followed by
//     their special syntax) → genericType.
func (p *Parser) parseBaseType() (*DataType, error) {
	switch p.cur.Kind {
	case kwROW:
		if p.peekNext().Kind == int('(') {
			return p.parseRowType()
		}
	case kwINTERVAL:
		if isIntervalFieldKind(p.peekNext().Kind) {
			return p.parseIntervalType()
		}
	case kwTIMESTAMP, kwTIME:
		return p.parseDateTimeType()
	case kwDOUBLE:
		if p.peekNext().Kind == kwPRECISION {
			return p.parseDoublePrecisionType()
		}
	case kwARRAY:
		if p.peekNext().Kind == int('<') {
			return p.parseLegacyArrayType()
		}
	case kwMAP:
		if p.peekNext().Kind == int('<') {
			return p.parseLegacyMapType()
		}
	}
	return p.parseGenericType()
}

// isIntervalFieldKind reports whether kind is one of the six intervalField
// keywords (YEAR/MONTH/DAY/HOUR/MINUTE/SECOND).
func isIntervalFieldKind(kind TokenKind) bool {
	switch kind {
	case kwYEAR, kwMONTH, kwDAY, kwHOUR, kwMINUTE, kwSECOND:
		return true
	default:
		return false
	}
}

// intervalFieldFromKind maps an intervalField keyword token to its
// IntervalField; the bool is false for a non-field kind.
func intervalFieldFromKind(kind TokenKind) (IntervalField, bool) {
	switch kind {
	case kwYEAR:
		return IntervalYear, true
	case kwMONTH:
		return IntervalMonth, true
	case kwDAY:
		return IntervalDay, true
	case kwHOUR:
		return IntervalHour, true
	case kwMINUTE:
		return IntervalMinute, true
	case kwSECOND:
		return IntervalSecond, true
	default:
		return 0, false
	}
}

// parseGenericType parses `identifier (LPAREN typeParameter (COMMA
// typeParameter)* RPAREN)?`. The leading identifier accepts any non-reserved
// keyword as a type name (Trino's type names are identifiers), so ARRAY / MAP /
// ROW / DOUBLE / INTERVAL used without their special syntax land here.
//
// If a '(' opens the parameter list, at least one typeParameter is required —
// `varchar()` and `foo()` are SYNTAX_ERRORs in Trino 481.
func (p *Parser) parseGenericType() (*DataType, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		// Re-shape the "expected identifier" message as "expected type" so
		// diagnostics at a type position are actionable.
		return nil, p.typeError()
	}
	dt := &DataType{
		Kind: TypeGeneric,
		Name: name.String(),
		Loc:  name.Loc,
	}
	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		params, endLoc, err := p.parseTypeParamList()
		if err != nil {
			return nil, err
		}
		dt.Params = params
		dt.Loc.End = endLoc.End
	}
	return dt, nil
}

// parseTypeParamList parses `typeParameter (COMMA typeParameter)* RPAREN` after
// the opening '(' has been consumed. At least one parameter is required.
// Returns the parameters and the Loc spanning through the closing ')'.
func (p *Parser) parseTypeParamList() ([]TypeParam, ast.Loc, error) {
	first, err := p.parseTypeParam()
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	params := []TypeParam{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseTypeParam()
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

// parseTypeParam parses one `typeParameter := INTEGER_VALUE_ | type`. An
// integer is taken literally (DECIMAL(10,2)); anything else recurses into a
// full type (QDIGEST(BIGINT), MAP(VARCHAR, ARRAY(INT))).
func (p *Parser) parseTypeParam() (TypeParam, error) {
	if p.cur.Kind == tokInteger {
		tok := p.advance()
		return TypeParam{IsInt: true, IntVal: tok.Ival, IntText: tok.Str, Loc: tok.Loc}, nil
	}
	inner, err := p.parseType()
	if err != nil {
		return TypeParam{}, err
	}
	return TypeParam{Type: inner, Loc: inner.Loc}, nil
}

// parseRowType parses a `ROW ( … )` type (the leading ROW and the '(' lookahead
// were confirmed by the caller). Trino accepts TWO disjoint readings of the
// parenthesized body and tries them in this order:
//
//  1. rowType — `rowField (, rowField)*`, where rowField is `type` or
//     `identifier type`. NO bare integers. This is ANTLR's `rowType` alt.
//  2. genericType "ROW" — `typeParameter (, typeParameter)*`, where a
//     typeParameter is `INTEGER | type`. This is the genericType fallback ANTLR
//     reaches when rowType fails, and it is the ONLY reading that admits integer
//     entries: ROW(1), ROW(1, 2), ROW(bigint, 1) all land here.
//
// The readings are mutually exclusive on integers and named fields, so a body
// that mixes them — ROW(a bigint, 1) (a named field plus an integer) — fits
// NEITHER and is a SYNTAX_ERROR, exactly as Trino rejects it. At least one entry
// is required (ROW() is a SYNTAX_ERROR under both readings).
func (p *Parser) parseRowType() (*DataType, error) {
	rowTok := p.advance() // consume ROW
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// Reading 1 — rowType (no integer fields).
	cp := p.checkpoint()
	if fields, closeTok, ok := p.tryRowFields(); ok {
		return &DataType{
			Kind:   TypeRow,
			Name:   "ROW",
			Fields: fields,
			Loc:    ast.Loc{Start: rowTok.Loc.Start, End: closeTok.Loc.End},
		}, nil
	}
	p.restore(cp)

	// Reading 2 — genericType "ROW" with a typeParameter list (admits integers).
	params, endLoc, err := p.parseTypeParamList()
	if err != nil {
		return nil, err
	}
	return &DataType{
		Kind:   TypeGeneric,
		Name:   "ROW",
		Params: params,
		Loc:    ast.Loc{Start: rowTok.Loc.Start, End: endLoc.End},
	}, nil
}

// tryRowFields parses `rowField (, rowField)* )` after the opening '(' has been
// consumed, returning the fields and the closing-')' token. ok is false (with
// the cursor left wherever the attempt stopped — the caller restores) if the
// body is not a valid rowField list, e.g. it contains a bare integer entry.
func (p *Parser) tryRowFields() (fields []RowField, closeTok Token, ok bool) {
	first, err := p.parseRowField()
	if err != nil {
		return nil, Token{}, false
	}
	fields = []RowField{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		f, err := p.parseRowField()
		if err != nil {
			return nil, Token{}, false
		}
		fields = append(fields, f)
	}
	closeTok, err = p.expect(int(')'))
	if err != nil {
		return nil, Token{}, false
	}
	return fields, closeTok, true
}

// parserCheckpoint captures enough parser+lexer state to rewind a speculative
// parse. The lexer is fully described by its byte position and pending error
// count (its input and baseOffset are immutable); the parser adds its current,
// previous, and buffered-lookahead tokens.
type parserCheckpoint struct {
	lexPos   int
	lexStart int
	lexErrs  int
	cur      Token
	prev     Token
	nextBuf  Token
	hasNext  bool
}

// checkpoint snapshots the parser/lexer state for a later restore.
func (p *Parser) checkpoint() parserCheckpoint {
	return parserCheckpoint{
		lexPos:   p.lexer.pos,
		lexStart: p.lexer.start,
		lexErrs:  len(p.lexer.errors),
		cur:      p.cur,
		prev:     p.prev,
		nextBuf:  p.nextBuf,
		hasNext:  p.hasNext,
	}
}

// restore rewinds the parser/lexer to a previously taken checkpoint, discarding
// any lex errors recorded after it (a speculative parse must not leak errors).
func (p *Parser) restore(c parserCheckpoint) {
	p.lexer.pos = c.lexPos
	p.lexer.start = c.lexStart
	if len(p.lexer.errors) > c.lexErrs {
		p.lexer.errors = p.lexer.errors[:c.lexErrs]
	}
	p.cur = c.cur
	p.prev = c.prev
	p.nextBuf = c.nextBuf
	p.hasNext = c.hasNext
}

// parseRowField parses one `rowField := type | identifier type` (the integer
// entries of ROW(1) / ROW(1, 2) are NOT handled here — they belong to the
// genericType "ROW" fallback the parseRowType dispatcher tries when this fails).
//
// The grammar is ambiguous between an unnamed `type` and a named
// `identifier type`: a field name is itself an identifier (and may even be a
// type keyword — `ROW(timestamp bigint)` names a field `timestamp`), while a
// type can be multi-token (DOUBLE PRECISION, TIME WITHOUT TIME ZONE, INTERVAL
// DAY TO SECOND, t ARRAY[5], ARRAY<…>, ARRAY(…)). Two tokens of lookahead are
// not enough to resolve every case (`ROW(a ARRAY<int>)` names a field `a` of
// type ARRAY<int>, but `ROW(bigint ARRAY)` is one unnamed postfix-array type —
// the deciding token comes after ARRAY).
//
// Resolution follows ANTLR's alternative order — the unnamed `type` alternative
// is listed first, so it is tried first: parse a single type and keep that
// reading if it lands exactly on a field boundary (',' or ')'). Only if the
// unnamed reading does not fit is the named form `identifier type` tried
// (consume a one-token name, then a full type, again requiring a clean
// boundary). Each attempt rewinds on failure via a parser/lexer checkpoint.
// This accepts every Trino-481 named/unnamed field form (ROW(a), ROW(a b),
// ROW(a, b), ROW(bigint), ROW(x bigint), ROW(timestamp bigint),
// ROW(DOUBLE PRECISION), ROW(TIME WITHOUT TIME ZONE), ROW(INTERVAL DAY TO
// SECOND), ROW(a ARRAY<int>), ROW(bigint ARRAY[5]); unknown type names fail
// later at resolution, not at parse time).
func (p *Parser) parseRowField() (RowField, error) {
	// Attempt 1 — unnamed `type` (ANTLR's first rowField alternative).
	cp := p.checkpoint()
	if fieldType, err := p.parseType(); err == nil && p.atRowFieldBoundary() {
		return RowField{Type: fieldType, Loc: fieldType.Loc}, nil
	}
	p.restore(cp)

	// Attempt 2 — named `identifier type`. The name is a single identifier
	// token; the type follows and must reach a field boundary.
	if isIdentifierStart(p.cur.Kind) {
		name := identFromToken(p.advance())
		if fieldType, err := p.parseType(); err == nil && p.atRowFieldBoundary() {
			return RowField{
				Name: name,
				Type: fieldType,
				Loc:  ast.Loc{Start: name.Loc.Start, End: fieldType.Loc.End},
			}, nil
		}
		p.restore(cp)
	}

	// Neither rowField reading fits. Return an error so the ROW dispatch can
	// fall back to the genericType "ROW" reading (which admits integer entries,
	// e.g. ROW(1)). The cursor is left for the caller's restore.
	return RowField{}, p.typeErrorAt("invalid ROW field")
}

// atRowFieldBoundary reports whether the cursor sits at the end of a ROW field
// (a ',' separating the next field or the closing ')').
func (p *Parser) atRowFieldBoundary() bool {
	return p.cur.Kind == int(',') || p.cur.Kind == int(')')
}

// parseIntervalType parses `INTERVAL from (TO to)?` (the caller confirmed the
// next token is an intervalField). The D1 divergence is enforced: a TO range is
// accepted only for valid SQL-standard qualifier pairs (ValidIntervalRange);
// an invalid pair (e.g. SECOND TO YEAR, YEAR TO DAY) is a SYNTAX_ERROR, exactly
// as Trino 481 rejects it.
func (p *Parser) parseIntervalType() (*DataType, error) {
	intervalTok := p.advance() // consume INTERVAL
	fromTok := p.advance()     // the intervalField (confirmed by caller)
	from, _ := intervalFieldFromKind(fromTok.Kind)

	dt := &DataType{
		Kind:         TypeInterval,
		Name:         "INTERVAL",
		IntervalFrom: from,
		Loc:          ast.Loc{Start: intervalTok.Loc.Start, End: fromTok.Loc.End},
	}

	if p.cur.Kind == kwTO {
		toTok := p.peekNext()
		to, ok := intervalFieldFromKind(toTok.Kind)
		if !ok {
			// `TO` not followed by an interval field — e.g. INTERVAL DAY TO foo.
			p.advance() // consume TO so the error points at the bad to-field
			return nil, p.typeErrorAt("expected interval field after TO")
		}
		if !ValidIntervalRange(from, to) {
			// Reversed or cross-family range — Trino rejects this at parse time.
			return nil, &ParseError{
				Loc: p.cur.Loc,
				Msg: "invalid interval qualifier: " + from.String() + " TO " + to.String(),
			}
		}
		p.advance() // consume TO
		toEnd := p.advance()
		toVal := to
		dt.IntervalTo = &toVal
		dt.Loc.End = toEnd.Loc.End
	}
	return dt, nil
}

// parseDateTimeType parses a TIMESTAMP or TIME type:
// `base (LPAREN precision RPAREN)? ((WITH | WITHOUT) TIME ZONE)?`. The base
// keyword (the caller confirmed it is TIMESTAMP or TIME) is retained verbatim
// for round-tripping. precision is a typeParameter (integer or, oddly, a
// nested type — Trino parses TIMESTAMP(bigint) and fails later).
func (p *Parser) parseDateTimeType() (*DataType, error) {
	baseTok := p.advance() // TIMESTAMP or TIME
	name := "TIMESTAMP"
	if baseTok.Kind == kwTIME {
		name = "TIME"
	}
	dt := &DataType{
		Kind: TypeDateTime,
		Name: name,
		Loc:  baseTok.Loc,
	}

	if p.cur.Kind == int('(') {
		p.advance() // consume '('
		prec, err := p.parseTypeParam()
		if err != nil {
			return nil, err
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		dt.Precision = &prec
		dt.Loc.End = closeTok.Loc.End
	}

	// Optional WITH/WITHOUT TIME ZONE.
	if p.cur.Kind == kwWITH || p.cur.Kind == kwWITHOUT {
		withTok := p.advance()
		if _, err := p.expect(kwTIME); err != nil {
			return nil, err
		}
		zoneTok, err := p.expect(kwZONE)
		if err != nil {
			return nil, err
		}
		dt.HasTimeZoneClause = true
		dt.WithTimeZone = withTok.Kind == kwWITH
		dt.Loc.End = zoneTok.Loc.End
	}
	return dt, nil
}

// parseDoublePrecisionType parses `DOUBLE PRECISION` (the caller confirmed both
// keywords). Plain `DOUBLE` is handled by genericType.
func (p *Parser) parseDoublePrecisionType() (*DataType, error) {
	doubleTok := p.advance()    // DOUBLE
	precisionTok := p.advance() // PRECISION
	return &DataType{
		Kind: TypeGeneric,
		Name: "DOUBLE PRECISION",
		Loc:  ast.Loc{Start: doubleTok.Loc.Start, End: precisionTok.Loc.End},
	}, nil
}

// parseLegacyArrayType parses `ARRAY < type >` (the caller confirmed ARRAY '<').
func (p *Parser) parseLegacyArrayType() (*DataType, error) {
	arrTok := p.advance() // ARRAY
	if _, err := p.expect(int('<')); err != nil {
		return nil, err
	}
	elem, err := p.parseType()
	if err != nil {
		return nil, err
	}
	gtTok, err := p.expect(int('>'))
	if err != nil {
		return nil, err
	}
	return &DataType{
		Kind:        TypeArray,
		Name:        "ARRAY",
		ElementType: elem,
		ArrayDim:    -1,
		Loc:         ast.Loc{Start: arrTok.Loc.Start, End: gtTok.Loc.End},
	}, nil
}

// parseLegacyMapType parses `MAP < keyType , valueType >` (the caller confirmed
// MAP '<').
func (p *Parser) parseLegacyMapType() (*DataType, error) {
	mapTok := p.advance() // MAP
	if _, err := p.expect(int('<')); err != nil {
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
	gtTok, err := p.expect(int('>'))
	if err != nil {
		return nil, err
	}
	return &DataType{
		Kind:      TypeMap,
		Name:      "MAP",
		KeyType:   keyType,
		ValueType: valueType,
		Loc:       ast.Loc{Start: mapTok.Loc.Start, End: gtTok.Loc.End},
	}, nil
}

// typeError returns a *ParseError describing a missing type at the current
// token, distinct from the generic identifier error so a type position reports
// "expected type".
func (p *Parser) typeError() *ParseError {
	if p.cur.Kind == tokEOF {
		return &ParseError{Loc: p.cur.Loc, Msg: "expected type, found end of input"}
	}
	text := p.cur.Str
	if text == "" {
		text = TokenName(p.cur.Kind)
	}
	return &ParseError{Loc: p.cur.Loc, Msg: "expected type, found " + text}
}

// typeErrorAt returns a *ParseError with a custom message at the current token.
func (p *Parser) typeErrorAt(msg string) *ParseError {
	return &ParseError{Loc: p.cur.Loc, Msg: msg}
}

// ParseDataType parses a complete Trino type from a standalone string,
// returning the *DataType and any ParseErrors. Trailing tokens after the type
// are reported as an error. It is the string-input counterpart of parseType,
// for tests and callers (catalog, deparse) that hold a type string rather than
// a token stream. Mirrors snowflake/parser.ParseDataType and
// ParseQualifiedName in identifiers.go.
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
	if p.cur.Kind != tokEOF {
		text := p.cur.Str
		if text == "" {
			text = TokenName(p.cur.Kind)
		}
		return dt, []ParseError{{Loc: p.cur.Loc, Msg: "unexpected token after type: " + text}}
	}
	return dt, nil
}
